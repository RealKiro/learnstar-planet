<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * 任课表：班级 × 科目 → 教师（全校智能排课的依据）
 *
 * @property int $id
 * @property int $school_id
 * @property int $class_id
 * @property string $subject_name
 * @property string $teacher_name
 */
class TimetableTeacherAssignment extends Model
{
    protected $fillable = [
        'school_id',
        'class_id',
        'subject_name',
        'teacher_name',
    ];

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }
}
