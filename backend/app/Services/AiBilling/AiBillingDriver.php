<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;

/**
 * 计费适配器接口：官方账单 API 优先，不支持的供应商由 LocalPreciseBillingDriver 回退本地计费
 */
interface AiBillingDriver
{
    /**
     * 是否提供官方 token 用量 API
     */
    public function supportsUsage(): bool;

    /**
     * 官方用量查询。返回 null 表示不支持 / 查询失败（调用方回退本地计费）。
     *
     * @param  array<string, mixed>  $provider  ai_settings.providers 中的单个供应商配置
     */
    public function getUsage(array $provider, ?int $schoolId = null, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot;

    /**
     * 官方余额查询。返回 null 表示不支持。
     *
     * @param  array<string, mixed>  $provider  ai_settings.providers 中的单个供应商配置
     */
    public function getBalance(array $provider): ?AiBalance;
}
