<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;
use Illuminate\Support\Facades\Http;

/**
 * Moonshot（月之暗面 Kimi）：无官方 token 用量 API，仅余额查询。
 * GET https://api.moonshot.cn/v1/users/me/balance
 */
class MoonshotBillingDriver implements AiBillingDriver
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

        $base = rtrim((string) ($provider['api_base'] ?? 'https://api.moonshot.cn/v1'), '/');

        try {
            $response = Http::withHeaders([
                'Authorization' => 'Bearer ' . $apiKey,
                'Accept' => 'application/json',
            ])->timeout(30)->get($base . '/users/me/balance');
        } catch (\Throwable) {
            return null;
        }

        if ($response->failed()) {
            return null;
        }

        $data = $response->json('data');
        if (!is_array($data)) {
            return null;
        }

        $total = (float) ($data['total_balance'] ?? $data['available_balance'] ?? 0);
        $available = isset($data['available_balance']) ? (float) $data['available_balance'] : $total;

        return new AiBalance(
            currency: 'CNY',
            totalBalance: $total,
            grantedBalance: isset($data['voucher_balance']) ? (float) $data['voucher_balance'] : null,
            isAvailable: $available > 0,
        );
    }
}
