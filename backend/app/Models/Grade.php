<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Grade extends Model
{
    protected $fillable = [
        'class_id', 'student_id', 'teacher_id',
        'exam_name', 'subject', 'score', 'rank_in_class',
        'rank_change',  // 排名变化：正数上升，负数下降
    ];

    protected $casts = [
        'score' => 'decimal:1',
        'rank_change' => 'integer',
    ];

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function student(): BelongsTo
    {
        return $this->belongsTo(Student::class);
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    /**
     * 考试名称列表
     */
    public static function getExamNames(): array
    {
        return ['期中考试', '期末考试', '月考一', '月考二', '月考三', '单元测试'];
    }

    /**
     * 班级成绩统计
     */
    public static function classStats(int $classId, string $examName, string $subject): array
    {
        $grades = self::where('class_id', $classId)
            ->where('exam_name', $examName)
            ->where('subject', $subject)
            ->get();

        if ($grades->isEmpty()) {
            return ['avg' => 0, 'max' => 0, 'min' => 0, 'pass_rate' => 0];
        }

        $scores = $grades->pluck('score')->map(fn ($s) => (float) $s);
        $passCount = $scores->filter(fn ($s) => $s >= 60)->count();

        return [
            'avg'       => round($scores->avg(), 1),
            'max'       => $scores->max(),
            'min'       => $scores->min(),
            'pass_rate' => round($passCount / $scores->count() * 100, 1),
            'count'     => $grades->count(),
        ];
    }

    /**
     * 分数段分布
     */
    public static function scoreDistribution(int $classId, string $examName, string $subject): array
    {
        $grades = self::where('class_id', $classId)
            ->where('exam_name', $examName)
            ->where('subject', $subject)
            ->get();

        return [
            '90-100' => $grades->whereBetween('score', [90, 100])->count(),
            '80-89'  => $grades->whereBetween('score', [80, 89])->count(),
            '70-79'  => $grades->whereBetween('score', [70, 79])->count(),
            '60-69'  => $grades->whereBetween('score', [60, 69])->count(),
            '0-59'   => $grades->where('score', '<', 60)->count(),
        ];
    }
}
