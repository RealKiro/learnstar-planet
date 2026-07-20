<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * @property int $id
 * @property int $school_id
 * @property string $from_currency
 * @property string $to_currency
 * @property string $rate
 * @property bool $is_active
 */
class ExchangeRate extends Model
{
    protected $fillable = [
        'school_id',
        'from_currency',
        'to_currency',
        'rate',
        'is_active',
    ];

    protected $casts = [
        'is_active' => 'boolean',
        'rate' => 'decimal:2',
    ];

    public function school(): BelongsTo
    {
        return $this->belongsTo(School::class);
    }
}

