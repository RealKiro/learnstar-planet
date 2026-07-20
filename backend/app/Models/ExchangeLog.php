<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * @property int $id
 * @property int $student_id
 * @property string $from_currency
 * @property string $to_currency
 * @property int $from_amount
 * @property int $to_amount
 * @property int|null $operated_by
 */
class ExchangeLog extends Model
{
    protected $fillable = [
        'student_id',
        'from_currency',
        'to_currency',
        'from_amount',
        'to_amount',
        'operated_by',
    ];

    public function student(): BelongsTo
    {
        return $this->belongsTo(Student::class);
    }

    public function operator(): BelongsTo
    {
        return $this->belongsTo(User::class, 'operated_by');
    }
}

