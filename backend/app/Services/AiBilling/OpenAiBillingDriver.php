<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use Carbon\CarbonInterface;
use Illuminate\Support\Facades\Http;

/**
 * OpenAI 官方 Usage API（需付费账号权限）。
 * GET {api_base}/usage?start_time&end_time&bucket_width=1d
 * 401/403/404/空数据返回 null → 调用方回退本地精确计费。
 */
class OpenAiBillingDriver implements AiBillingDriver
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

        $base = rtrim((string) ($provider['api_base'] ?? 'https://api.openai.com/v1'), '/');
        $from = $from ?: now()->subDays(30);
        $to = $to ?: now();

        try {
            $response = Http::withHeaders([
                'Authorization' => 'Bearer ' . $apiKey,
            ])->timeout(30)->get($base . '/usage', [
                'start_time' => $from->getTimestamp(),
                'end_time' => $to->getTimestamp(),
                'bucket_width' => '1d',
            ]);
        } catch (\Throwable) {
            return null;
        }

        if ($response->failed()) {
            return null;
        }

        $buckets = $response->json('data', []);
        if (!is_array($buckets) || count($buckets) === 0) {
            return null;
        }

        $prompt = 0;
        $completion = 0;
        $cost = 0.0;
        foreach ($buckets as $bucket) {
            $prompt += (int) ($bucket['n_context_tokens_total'] ?? 0);
            $completion += (int) ($bucket['n_generated_tokens_total'] ?? 0);
            $cost += (float) ($bucket['cost_in_usd'] ?? 0);
        }

        return new AiUsageSnapshot(
            promptTokens: $prompt,
            completionTokens: $completion,
            cost: $cost,
            currency: 'USD',
            source: 'official',
            raw: $response->json(),
        );
    }

    public function getBalance(array $provider): ?AiBalance
    {
        // OpenAI 无公开余额接口
        return null;
    }
}
