<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

use App\Models\AiConversation;
use Carbon\CarbonInterface;

/**
 * 通用回退：本地精确计费。
 * 从 ai_conversations 聚合本供应商在时间区间内的真实 token 与 cost（逐笔按单价累加）。
 */
class LocalPreciseBillingDriver implements AiBillingDriver
{
    public function supportsUsage(): bool
    {
        return true;
    }

    public function getUsage(array $provider, ?int $schoolId = null, ?CarbonInterface $from = null, ?CarbonInterface $to = null): ?AiUsageSnapshot
    {
        if ($schoolId === null) {
            return null;
        }

        $query = AiConversation::where('school_id', $schoolId)
            ->where('provider', (string) ($provider['id'] ?? ''))
            ->where('status', 'completed');

        if ($from !== null) {
            $query->where('created_at', '>=', $from);
        }
        if ($to !== null) {
            $query->where('created_at', '<=', $to);
        }

        $row = (clone $query)
            ->selectRaw('COALESCE(SUM(prompt_tokens),0) as prompt, COALESCE(SUM(completion_tokens),0) as completion, COALESCE(SUM(cost),0) as cost')
            ->first();

        $prompt = (int) ($row->prompt ?? 0);
        $completion = (int) ($row->completion ?? 0);
        $cost = (float) ($row->cost ?? 0);

        return new AiUsageSnapshot(
            promptTokens: $prompt,
            completionTokens: $completion,
            cost: $cost,
            currency: (string) ($provider['currency'] ?? 'CNY'),
            source: 'local',
        );
    }

    public function getBalance(array $provider): ?AiBalance
    {
        return null;
    }
}
