<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\AiConversation;
use App\Models\AiSetting;
use App\Models\User;
use App\Services\AiBilling\AiBillingService;

/**
 * 教师端 AI 助教服务：配置态查询、对话（含本地精确计费与用量回写）、用量汇总。
 *
 * 供应商解析规则（多供应商优先，兼容旧版单供应商字段）：
 * 任一启用的供应商配了 api_key 即可用；否则回退旧字段 provider/api_key。
 */
class AiAssistantService
{
    /**
     * AI 配置态（前端据此隐藏 AI 助教入口）。
     *
     * @return array{enabled: bool}
     */
    public function configFor(User $teacher): array
    {
        $setting = AiSetting::where('school_id', $teacher->school_id)->first();

        $enabled = false;
        if ($setting && $setting->enabled) {
            $providers = $setting->providers ?: [];
            // 任一启用的供应商配了 api_key 即视为已配置
            $hasProviderKey = collect($providers)->contains(
                fn ($p) => !empty($p['is_active']) && !empty($p['api_key'])
            );
            // 兼容旧字段（单 provider + api_key）
            $hasLegacyKey = !empty($setting->provider) && !empty($setting->api_key);
            $enabled = $hasProviderKey || $hasLegacyKey;
        }

        return ['enabled' => $enabled];
    }

    /**
     * 教师对话：解析供应商 → 落会话 → 调 AI → 本地计费 → 回写用量。
     * 未启用/未配置等软失败返回提示文案（HTTP 200，与历史行为一致）。
     *
     * @param  string|null  $model  教师指定的模型（New API 式映射的请求侧）。
     *                              必须在供应商模型白名单内（models 列表 / model_map 键 / 配置主模型），
     *                              不在白名单时回退默认模型；命中 model_map 则重定向到上游模型。
     */
    public function chat(User $teacher, string $message, ?string $model = null): string
    {
        $settings = AiSetting::where('school_id', $teacher->school_id)->first();
        if (!$settings || !$settings->enabled) {
            return 'AI 功能未启用，请联系管理员配置';
        }

        // 从多供应商配置中查找启用的供应商，兼容旧版单供应商配置
        $activeProvider = null;
        foreach ($settings->providers ?: [] as $p) {
            if (!empty($p['is_active']) && !empty($p['api_key'])) {
                $activeProvider = $p;
                break;
            }
        }
        if (!$activeProvider && !empty($settings->api_key)) {
            $activeProvider = [
                'id' => $settings->provider ?: 'openai',
                'api_key' => $settings->api_key,
                'api_base' => $settings->api_base,
                'model' => $settings->model ?: 'gpt-3.5-turbo',
            ];
        }
        if (!$activeProvider) {
            return '请先在 AI 中心配置并启用一个供应商';
        }

        // 模型解析（New API 式）：默认主模型 → 白名单校验 → model_map 重定向
        $defaultModel = $activeProvider['model'] ?: 'gpt-3.5-turbo';
        $allowed = $this->allowedModelsFor($activeProvider);
        $requested = (string) ($model ?: $defaultModel);
        if (!in_array($requested, $allowed, true)) {
            $requested = $defaultModel;
        }
        $modelMap = is_array($activeProvider['model_map'] ?? null) ? $activeProvider['model_map'] : [];
        $upstreamModel = (string) ($modelMap[$requested] ?? $requested);

        $classId = $teacher->getSetting('active_class_id') ?: null;

        $conversation = AiConversation::create([
            'school_id' => $teacher->school_id,
            'class_id' => $classId,
            'student_name' => '教师',
            'provider' => $activeProvider['id'],
            'question' => $message,
            'status' => 'pending',
        ]);

        try {
            $ai = new AiService();
            $result = $ai->chat(
                provider: $activeProvider['id'],
                apiKey: $activeProvider['api_key'],
                model: $upstreamModel,
                question: $message,
                apiBase: $activeProvider['api_base'] ?? null,
                maxTokens: $settings->max_tokens,
            );
            $reply = $result['answer'];
            $promptTokens = $result['prompt_tokens'] ?? 0;
            $completionTokens = $result['completion_tokens'] ?? 0;
            $tokensUsed = $result['tokens_used'] ?? ($promptTokens + $completionTokens);
        } catch (\Throwable $e) {
            $reply = 'AI 服务暂时不可用';
            $tokensUsed = 0;
            $promptTokens = 0;
            $completionTokens = 0;
        }

        // 本地精确计费：按供应商单价计算本次费用
        $cost = app(AiBillingService::class)->recordUsage($activeProvider, $promptTokens, $completionTokens);
        $currency = $activeProvider['currency'] ?? 'CNY';

        $conversation->update([
            'answer' => $reply,
            'tokens_used' => $tokensUsed,
            'prompt_tokens' => $promptTokens,
            'completion_tokens' => $completionTokens,
            'cost' => $cost,
            'currency' => $currency,
            'status' => 'completed',
        ]);

        if ($tokensUsed > 0) {
            $settings->increment('tokens_used', $tokensUsed);
            $providers = $settings->providers ?: [];
            foreach ($providers as &$p) {
                if (($p['id'] ?? '') === ($activeProvider['id'] ?? '')) {
                    $p['tokens_used'] = ($p['tokens_used'] ?? 0) + $tokensUsed;
                    $p['total_calls'] = ($p['total_calls'] ?? 0) + 1;
                    $p['estimated_cost'] = ($p['estimated_cost'] ?? 0) + $cost;
                    $p['currency'] = $currency;
                    break;
                }
            }
            $settings->providers = $providers;
            $settings->save();
        }

        return $reply;
    }

    /**
     * 教师端 AI 用量（前端 AIPage 每次发送后刷新）。
     *
     * @return array{configured: bool, provider: ?string, model: ?string, models: array<int, string>, tokens_used: int, estimated_cost: float, currency: string}
     */
    public function usageFor(User $teacher): array
    {
        $settings = AiSetting::where('school_id', $teacher->school_id)->first();

        $active = null;
        if ($settings) {
            foreach ($settings->providers ?: [] as $p) {
                if (!empty($p['is_active']) && !empty($p['api_key'])) {
                    $active = $p;
                    break;
                }
            }
            if (!$active && !empty($settings->api_key)) {
                $active = ['id' => $settings->provider ?: 'openai', 'model' => $settings->model];
            }
        }

        $estimatedCost = 0.0;
        $currency = 'CNY';
        foreach ($settings->providers ?? [] as $p) {
            $estimatedCost += (float) ($p['estimated_cost'] ?? 0);
            $currency = (string) ($p['currency'] ?? $currency);
        }

        return [
            'configured' => $settings !== null && $settings->enabled && $active !== null,
            'provider' => $active['id'] ?? ($settings->provider ?? null),
            'model' => $active['model'] ?? ($settings->model ?? null),
            'models' => $active !== null ? $this->allowedModelsFor($active) : [],
            'tokens_used' => $settings ? (int) $settings->tokens_used : 0,
            'estimated_cost' => $estimatedCost,
            'currency' => $currency,
        ];
    }

    /**
     * 供应商可选模型白名单：主模型 + models 多选列表 + model_map 请求侧键名。
     *
     * @param  array<string, mixed>  $provider
     * @return array<int, string>
     */
    private function allowedModelsFor(array $provider): array
    {
        $models = [];
        $default = (string) ($provider['model'] ?? '');
        if ($default !== '') {
            $models[] = $default;
        }
        foreach ((array) ($provider['models'] ?? []) as $m) {
            if (is_string($m) && $m !== '') {
                $models[] = $m;
            }
        }
        $modelMap = $provider['model_map'] ?? null;
        if (is_array($modelMap)) {
            foreach (array_keys($modelMap) as $k) {
                if (is_string($k) && $k !== '') {
                    $models[] = $k;
                }
            }
        }

        return array_values(array_unique($models));
    }

    /**
     * 教师端 AI 预设命令。
     *
     * @return array<int, array{label: string, prompt: string}>
     */
    public function commands(): array
    {
        return [
            ['label' => '📝 本周教学总结', 'prompt' => '请帮我写一份本周教学总结，包含本周教学目标、课堂情况、学生表现和下周教学计划。'],
            ['label' => '🏅 积分规则建议', 'prompt' => '请根据班级日常情况，生成一套适合小学生的积分奖励规则建议。'],
            ['label' => '🎯 班会活动方案', 'prompt' => '请设计一个有趣的小学生班会活动方案，包含活动目标、流程和所需材料。'],
            ['label' => '📋 出练习题', 'prompt' => '请出一组适合本年级学生的练习题，包含题目和参考答案。'],
        ];
    }
}
