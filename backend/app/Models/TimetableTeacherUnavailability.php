<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * 教师不可用时段（学校级）：weekday 1-7 + period_index。
 * 全校智能排课时视为该教师已占用。
 *
 * @property int $id
 * @property int $school_id
 * @property string $teacher_name
 * @property int $weekday
 * @property int $period_index
 */
class TimetableTeacherUnavailability extends Model
{
    protected $fillable = [
        'school_id',
        'teacher_name',
        'weekday',
        'period_index',
    ];
}
