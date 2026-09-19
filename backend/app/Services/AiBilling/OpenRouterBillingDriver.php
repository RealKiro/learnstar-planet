<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;
use Illuminate\Support\Facades\Http;

/**
 * OpenRouter：credits 余额与已用金额查询。
 * GET https://openrouter.ai/api/v1/auth/key（Bearer Key）
 * → data.usage（已用 USD）/ data.limit（额度上限，null 为不限额）
 */
class OpenRouterBillingDriver implements AiBillingDriver
{
    public function supportsUsage(): bool
    {
        return true;
    }

    public function getUsage(array $provider, ?int $schoolId = null, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot
    {
        $info = $this->keyInfo($provider);
        if ($info === null || !isset($info['usage'])) {
            return null;
        }

        return new AiUsageSnapshot(
            promptTokens: 0,
            completionTokens: 0,
            cost: (float) $info['usage'],
            currency: 'USD',
            source: 'official',
        );
    }

    public function getBalance(array $provider): ?AiBalance
    {
        $info = $this->keyInfo($provider);
        if ($info === null || !isset($info['limit']) || $info['limit'] === null) {
            // 不限额度的 Key 无「余额」概念
            return null;
        }

        $limit = (float) $info['limit'];
        $used = (float) ($info['usage'] ?? 0);

        return new AiBalance(
            currency: 'USD',
            totalBalance: round($limit - $used, 4),
            isAvailable: $limit - $used > 0,
        );
    }

    /**
     * @return array<string, mixed>|null
     */
    private function keyInfo(array $provider): ?array
    {
        $apiKey = $provider['api_key'] ?? '';
        if ($apiKey === '') {
            return null;
        }

        try {
            $response = Http::withHeaders([
                'Authorization' => 'Bearer ' . $apiKey,
                'Accept' => 'application/json',
            ])->timeout(30)->get('https://openrouter.ai/api/v1/auth/key');
        } catch (\Throwable) {
            return null;
        }

        if ($response->failed()) {
            return null;
        }

        $data = $response->json('data');

        return is_array($data) ? $data : null;
    }
}
