<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

/**
 * 统一的 AI 用量快照（官方账单 API 或本地精确计费）
 */
final class AiUsageSnapshot
{
    public function __construct(
        public readonly int $promptTokens,
        public readonly int $completionTokens,
        public readonly float $cost,
        public readonly string $currency = 'CNY',
        public readonly string $source = 'local', // 'official' | 'local'
        public readonly ?array $raw = null,
    ) {
    }

    public function totalTokens(): int
    {
        return $this->promptTokens + $this->completionTokens;
    }
}
