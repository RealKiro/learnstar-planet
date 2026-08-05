<?php

declare(strict_types=1);

namespace App\Services\AiBilling;

/**
 * 供应商官方余额快照
 */
final class AiBalance
{
    public function __construct(
        public readonly string $currency,
        public readonly float $totalBalance,
        public readonly ?float $grantedBalance = null,
        public readonly bool $isAvailable = true,
    ) {
    }
}
