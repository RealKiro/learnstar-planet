<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;
use Illuminate\Support\Facades\Http;

/**
 * 中转平台（New API / One API 等自建网关）：
 * 沿用 OpenAI billing 兼容端点（one-api / new-api 均实现）：
 *   GET {base}/v1/dashboard/billing/subscription → hard_limit_usd（剩余额度 USD）
 *   GET {base}/v1/dashboard/billing/usage?start_date&end_date → total_usage（美分，已用 USD）
 * 部分部署不含 /v1 前缀，自动回退无前缀路径。
 */
class RelayBillingDriver implements AiBillingDriver
{
    public function supportsUsage(): bool
    {
        return true;
    }

    public function getUsage(array $provider, ?int $schoolId = null, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot
    {
        $apiKey = $provider['api_key'] ?? '';
        if ($apiKey === '') {
            return null;
        }

        $usedUsd = $this->fetchUsedUsd($provider, $from, $to);
        if ($usedUsd === null) {
            return null;
        }

        return new AiUsageSnapshot(
            promptTokens: 0,
            completionTokens: 0,
            cost: $usedUsd,
            currency: 'USD',
            source: 'official',
        );
    }

    public function getBalance(array $provider): ?AiBalance
    {
        $apiKey = $provider['api_key'] ?? '';
        if ($apiKey === '') {
            return null;
        }

        $json = $this->billingGet($provider, '/subscription');
        if ($json === null) {
            return null;
        }

        $limit = $json['hard_limit_usd'] ?? $json['system_hard_limit_usd'] ?? null;
        if ($limit === null) {
            return null;
        }

        return new AiBalance(
            currency: 'USD',
            totalBalance: round((float) $limit, 4),
            isAvailable: (float) $limit > 0,
        );
    }

    /**
     * 查询区间已用金额（USD）。total_usage 单位为美分。
     */
    private function fetchUsedUsd(array $provider, ?CarbonInterface $from, ?CarbonInterface $to): ?float
    {
        $start = ($from ?: now()->subDays(30))->format('Y-m-d');
        $end = ($to ?: now()->addDay())->format('Y-m-d');

        $json = $this->billingGet($provider, '/usage', ['start_date' => $start, 'end_date' => $end]);
        if ($json === null || !isset($json['total_usage'])) {
            return null;
        }

        return round(((float) $json['total_usage']) / 100, 6);
    }

    /**
     * 计费端点 GET，自动尝试 /v1 前缀与无前缀两种部署布局。
     *
     * @param  array<string, string>  $query
     * @return array<string, mixed>|null
     */
    private function billingGet(array $provider, string $path, array $query = []): ?array
    {
        $apiKey = $provider['api_key'] ?? '';
        $base = rtrim((string) ($provider['api_base'] ?? ''), '/');
        if ($apiKey === '' || $base === '') {
            return null;
        }

        // OpenAI 兼容 base 以 /v1 结尾 → 计费端点挂在同一 /v1 下；剥掉后两种布局都试
        $origin = preg_replace('/\/v1$/', '', $base) ?? $base;
        $candidates = [$origin . '/v1/dashboard/billing' . $path, $origin . '/dashboard/billing' . $path];

        foreach ($candidates as $url) {
            try {
                $response = Http::withHeaders([
                    'Authorization' => 'Bearer ' . $apiKey,
                    'Accept' => 'application/json',
                ])->timeout(30)->get($url, $query);
            } catch (\Throwable) {
                continue;
            }

            if ($response->failed()) {
                continue;
            }

            $json = $response->json();

            return is_array($json) ? $json : null;
        }

        return null;
    }
}
