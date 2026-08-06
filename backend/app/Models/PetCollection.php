<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * 宠物图鉴收藏：每个学生 × 物种 一条记录，记录该物种的最高养成进度，
 * 用于「切换宠物进度全保留」与「图鉴收集」玩法。
 */
class PetCollection extends Model
{
    protected $table = 'pet_collections';

    protected $fillable = [
        'student_id',
        'species',
        'level',
        'experience',
        'mood',
        'is_active',
    ];

    protected $casts = [
        'level' => 'integer',
        'experience' => 'integer',
        'mood' => 'integer',
        'is_active' => 'boolean',
    ];

    public function student(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(Student::class);
    }

    /** 当前学生已解锁的图鉴槽位数（每 100 积分 +1，初始 1） */
    public static function unlockSlotsForScore(mixed $totalScore): int
    {
        return 1 + intdiv(max(0, (int) $totalScore), 100);
    }
}
