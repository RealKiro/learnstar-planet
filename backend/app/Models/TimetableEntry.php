<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * 排课格子：某班级 · 星期 N · 第 M 节 · 单双周 → 科目
 */
class TimetableEntry extends Model
{
    /** @var list<string> */
    protected $fillable = [
        'class_id', 'weekday', 'period_index', 'week_type', 'subject_id', 'teacher_name', 'room',
    ];

    /** @var array<string, string> */
    protected $casts = [
        'weekday' => 'integer',
        'period_index' => 'integer',
        'subject_id' => 'integer',
    ];

    public function subject(): BelongsTo
    {
        return $this->belongsTo(Subject::class, 'subject_id');
    }
}
