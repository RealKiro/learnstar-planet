<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;
use Illuminate\Support\Facades\Http;

/**
 * DeepSeek：无官方 token 用量 API，仅提供余额查询。
 * GET https://api.deepseek.com/user/balance
 */
class DeepSeekBillingDriver implements AiBillingDriver
{
    public function supportsUsage(): bool
    {
        return false;
    }

    public function getUsage(array $provider, ?int $schoolId = null, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot
    {
        // DeepSeek 无官方用量明细 API → 走本地精确计费
        return null;
    }

    public function getBalance(array $provider): ?AiBalance
    {
        $apiKey = $provider['api_key'] ?? '';
        if ($apiKey === '') {
            return null;
        }

        try {
            $response = Http::withHeaders([
                'Authorization' => 'Bearer ' . $apiKey,
                'Accept' => 'application/json',
            ])->timeout(30)->get('https://api.deepseek.com/user/balance');
        } catch (\Throwable) {
            return null;
        }

        if ($response->failed()) {
            return null;
        }

        $info = $response->json('balance_infos.0');
        if (!is_array($info)) {
            return null;
        }

        return new AiBalance(
            currency: (string) ($info['currency'] ?? 'CNY'),
            totalBalance: (float) ($info['total_balance'] ?? 0),
            grantedBalance: isset($info['granted_balance']) ? (float) $info['granted_balance'] : null,
            isAvailable: (bool) ($response->json('is_available') ?? true),
        );
    }
}
