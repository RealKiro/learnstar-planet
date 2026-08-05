<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use App\Models\AiSetting;
use Carbon\CarbonInterface;

/**
 * AI 计费核心服务。
 * - 官方账单 API 优先（OpenAI Usage / DeepSeek Balance）
 * - 不支持的供应商回退本地精确计费（响应 usage 拆 prompt/completion × 配置单价）
 */
class AiBillingService
{
    /** @var array<string, AiBillingDriver> */
    private array $drivers = [];

    /**
     * 返回供应商对应的计费适配器
     */
    public function driverFor(string $providerId): AiBillingDriver
    {
        return $this->drivers[$providerId] ??= match ($providerId) {
            'openai' => new OpenAiBillingDriver(),
            'deepseek' => new DeepSeekBillingDriver(),
            default => new LocalPreciseBillingDriver(),
        };
    }

    /**
     * 按配置单价计算一次对话的费用。
     * 单价存于 providers[].input_price_per_m / output_price_per_m（每百万 token）。
     *
     * @param  array<string, mixed>  $provider
     */
    public function calculateCost(array $provider, int $promptTokens, int $completionTokens): float
    {
        $inputPrice = (float) ($provider['input_price_per_m'] ?? 0);
        $outputPrice = (float) ($provider['output_price_per_m'] ?? 0);

        return round(
            ($promptTokens / 1_000_000) * $inputPrice
            + ($completionTokens / 1_000_000) * $outputPrice,
            6,
        );
    }

    /**
     * 记录一次对话用量并返回 cost（调用方负责累加到 providers[].estimated_cost）。
     *
     * @param  array<string, mixed>  $provider
     */
    public function recordUsage(array $provider, int $promptTokens, int $completionTokens): float
    {
        return $this->calculateCost($provider, $promptTokens, $completionTokens);
    }

    /**
     * 同步官方用量：官方 API 可用则返回官方快照，否则回退本地聚合。
     */
    public function syncOfficialUsage(AiSetting $setting, string $providerId, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot
    {
        $provider = $this->findProvider($setting, $providerId);
        if ($provider === null) {
            return null;
        }

        $driver = $this->driverFor($providerId);

        if ($driver->supportsUsage()) {
            $official = $driver->getUsage($provider, $setting->school_id, $from, $to);
            if ($official !== null) {
                return $official;
            }
        }

        return $driver->getUsage($provider, $setting->school_id, $from, $to);
    }

    /**
     * 官方余额查询（DeepSeek 等支持），不支持返回 null。
     */
    public function getBalance(AiSetting $setting, string $providerId): ?AiBalance
    {
        $provider = $this->findProvider($setting, $providerId);
        if ($provider === null) {
            return null;
        }

        return $this->driverFor($providerId)->getBalance($provider);
    }

    /**
     * @return array<string, mixed>|null
     */
    private function findProvider(AiSetting $setting, string $providerId): ?array
    {
        foreach ($setting->providers ?: [] as $p) {
            if (($p['id'] ?? '') === $providerId) {
                return $p;
            }
        }

        return null;
    }
}
