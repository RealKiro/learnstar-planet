<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Attendance extends Model
{
    protected $fillable = [
        'class_id', 'student_id', 'teacher_id',
        'date', 'status', 'source', 'remark', 'leave_record_id',
        'sign_in_at',
    ];

    protected $casts = [
        'date' => 'date',
        'sign_in_at' => 'datetime',
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

    public function leaveRecord(): BelongsTo
    {
        return $this->belongsTo(WechatWorkLeaveRecord::class, 'leave_record_id');
    }

    public static function getStatuses(): array
    {
        return [
            'present' => '到课',
            'late'     => '迟到',
            'leave'    => '请假',
            'absent'   => '缺勤',
        ];
    }

    public static function getSources(): array
    {
        return [
            'auto'         => '自动',
            'wechat_work'  => '企业微信',
            'manual'       => '手动标记',
        ];
    }

    /**
     * 今日考勤统计
     */
    public static function todaySummary(int $classId): array
    {
        $today = now()->toDateString();
        $records = self::where('class_id', $classId)->where('date', $today)->get();

        return [
            'present' => $records->where('status', 'present')->count(),
            'late'    => $records->where('status', 'late')->count(),
            'leave'   => $records->where('status', 'leave')->count(),
            'absent'  => $records->where('status', 'absent')->count(),
        ];
    }
}
