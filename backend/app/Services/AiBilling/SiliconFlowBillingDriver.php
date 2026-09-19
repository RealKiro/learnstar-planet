<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;
use Illuminate\Support\Facades\Http;

/**
 * SiliconFlow（硅基流动）：无官方 token 用量 API，仅余额查询。
 * GET https://api.siliconflow.cn/v1/user/info → data.balance / data.totalBalance（CNY）
 */
class SiliconFlowBillingDriver implements AiBillingDriver
{
    public function supportsUsage(): bool
    {
        return false;
    }

    public function getUsage(array $provider, ?int $schoolId = null, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot
    {
        return null;
    }

    public function getBalance(array $provider): ?AiBalance
    {
        $apiKey = $provider['api_key'] ?? '';
        if ($apiKey === '') {
            return null;
        }

        $base = rtrim((string) ($provider['api_base'] ?? 'https://api.siliconflow.cn/v1'), '/');

        try {
            $response = Http::withHeaders([
                'Authorization' => 'Bearer ' . $apiKey,
                'Accept' => 'application/json',
            ])->timeout(30)->get($base . '/user/info');
        } catch (\Throwable) {
            return null;
        }

        if ($response->failed()) {
            return null;
        }

        $data = $response->json('data');
        if (!is_array($data) || !isset($data['balance']) && !isset($data['totalBalance'])) {
            return null;
        }

        $total = (float) ($data['totalBalance'] ?? $data['balance'] ?? 0);

        return new AiBalance(
            currency: 'CNY',
            totalBalance: $total,
            grantedBalance: isset($data['balance']) ? (float) $data['balance'] : null,
            isAvailable: (float) ($data['balance'] ?? 0) > 0,
        );
    }
}
