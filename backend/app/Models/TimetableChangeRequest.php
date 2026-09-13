<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * 课表修改申请（教师提交 → 管理员审核 → 通过后应用）
 *
 * @property array{subjects?: array, periods?: array, entries?: array} $payload
 */
class TimetableChangeRequest extends Model
{
    public const STATUS_PENDING = 'pending';
    public const STATUS_APPROVED = 'approved';
    public const STATUS_REJECTED = 'rejected';

    /** @var list<string> */
    protected $fillable = [
        'school_id', 'class_id', 'requested_by', 'payload', 'entry_count',
        'status', 'reviewed_by', 'review_note', 'reviewed_at',
    ];

    /** @var array<string, string> */
    protected $casts = [
        'payload' => 'array',
        'entry_count' => 'integer',
        'reviewed_at' => 'datetime',
    ];

    public function requester(): BelongsTo
    {
        return $this->belongsTo(User::class, 'requested_by');
    }

    public function reviewer(): BelongsTo
    {
        return $this->belongsTo(User::class, 'reviewed_by');
    }

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }
}
