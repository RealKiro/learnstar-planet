<?php

declare(strict_types=1);

namespace App\Services;

use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;

/**
 * AI 助教服务 — 对接大语言模型 API
 *
 * 支持 OpenAI 兼容接口的提供商：
 * - DeepSeek (默认, api.deepseek.com)
 * - 通义千问 (dashscope.aliyuncs.com)
 * - OpenAI (api.openai.com)
 * - 本地模型 (通过 ollama 等)
 */
class AIService
{
    private string $provider;

    private string $apiKey;

    private string $apiBase;

    private string $model;

    private int $maxTokens;

    private float $temperature;

    public function __construct()
    {
        $this->provider = config('ai.provider', '');
        $this->apiKey = config('ai.api_key', '');
        $this->apiBase = rtrim(config('ai.api_base', ''), '/');
        $this->model = config('ai.model', 'deepseek-chat');
        $this->maxTokens = (int) config('ai.max_tokens', 2000);
        $this->temperature = (float) config('ai.temperature', 0.7);
    }

    /**
     * 检查 AI 功能是否已配置
     */
    public function isConfigured(): bool
    {
        return !empty($this->provider) && !empty($this->apiKey);
    }

    /**
     * 发送聊天消息并获取回复
     */
    public function chat(string $message, array $history = []): string
    {
        if (!$this->isConfigured()) {
            return 'AI 助教功能未配置。请在 .env 文件中设置 AI_PROVIDER 和 AI_API_KEY。';
        }

        $messages = $this->buildMessages($message, $history);

        try {
            $response = Http::timeout(30)
                ->withHeaders([
                    'Authorization' => 'Bearer ' . $this->apiKey,
                    'Content-Type' => 'application/json',
                ])
                ->post($this->getApiUrl(), [
                    'model' => $this->model,
                    'messages' => $messages,
                    'max_tokens' => $this->maxTokens,
                    'temperature' => $this->temperature,
                    'stream' => false,
                ]);

            if ($response->successful()) {
                $data = $response->json();

                return $data['choices'][0]['message']['content'] ?? '抱歉，AI 返回了空回复。';
            }

            $status = $response->status();
            $body = $response->body();
            Log::warning("AI API 请求失败: status={$status}, body={$body}");

            if ($status === 401) {
                return 'AI API Key 无效，请在 .env 中检查 AI_API_KEY 配置。';
            }
            if ($status === 429) {
                return 'AI 请求过于频繁，请稍后再试。';
            }

            return "AI 服务暂时不可用 (HTTP {$status})，请稍后重试。";
        } catch (\Throwable $e) {
            Log::error('AI API 调用异常: ' . $e->getMessage());

            return 'AI 服务连接失败，请检查网络连接和 API 地址配置。';
        }
    }

    /**
     * 获取当前配置的提供商信息
     */
    public function getUsageInfo(): array
    {
        return [
            'configured' => $this->isConfigured(),
            'provider' => $this->provider ?: '未配置',
            'model' => $this->model,
        ];
    }

    /**
     * 构建消息列表
     */
    private function buildMessages(string $message, array $history): array
    {
        $systemPrompt = '你是一个班级管理系统的 AI 助教，名叫"学趣星球"。'
            . '你帮助教师处理日常班级管理事务，如：'
            . '积分管理、学生激励、课堂活动建议等。'
            . '请用中文回答，语气亲切专业，回答简洁实用。';

        $messages = [['role' => 'system', 'content' => $systemPrompt]];

        // 添加历史消息（最多保留最近 10 轮）
        $maxHistory = min(count($history), 10);
        for ($i = count($history) - $maxHistory; $i < count($history); $i++) {
            $messages[] = $history[$i];
        }

        $messages[] = ['role' => 'user', 'content' => $message];

        return $messages;
    }

    /**
     * 获取 API 地址
     */
    private function getApiUrl(): string
    {
        // 如果配置了完整 API 地址，直接使用
        if ($this->apiBase) {
            return $this->apiBase . '/chat/completions';
        }

        // 按提供商返回默认地址
        return match ($this->provider) {
            'deepseek' => 'https://api.deepseek.com/v1/chat/completions',
            'openai' => 'https://api.openai.com/v1/chat/completions',
            'qwen' => 'https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation',
            'moonshot' => 'https://api.moonshot.cn/v1/chat/completions',
            default => 'https://api.deepseek.com/v1/chat/completions',
        };
    }
}
