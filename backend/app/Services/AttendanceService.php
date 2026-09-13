<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Attendance;
use App\Models\Student;
use App\Models\User;
use Illuminate\Support\Collection;

/**
 * 考勤读写（含企业微信同步入口）。
 *
 * 只负责数据操作；HTTP 校验与响应包装留在控制器。
 * 所有查询都经 TeacherClassScope 限定班级作用域。
 */
class AttendanceService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
        private readonly WechatWorkAttendanceService $wechat,
    ) {
    }

    /** 今日考勤明细（含学生信息与请假记录） */
    public function today(User $teacher): Collection
    {
        return Attendance::whereIn('class_id', $this->scope->accessibleIds($teacher))
            ->whereDate('date', today())
            ->with(['student:id,name,student_no,class_id', 'leaveRecord:id,sp_no,leave_type,reason'])
            ->get()
            ->map(fn (Attendance $a) => [
                'id' => $a->id,
                'student_id' => $a->student_id,
                'student_name' => $a->student?->name,
                'student_no' => $a->student?->student_no,
                'status' => $a->status,
                'source' => $a->source,
                'remark' => $a->remark,
                'check_in_time' => $a->sign_in_at?->toDateTimeString(),
                'leave_record' => $a->leaveRecord ? [
                    'sp_no' => $a->leaveRecord->sp_no,
                    'leave_type' => $a->leaveRecord->leave_type,
                    'reason' => $a->leaveRecord->reason,
                ] : null,
            ]);
    }

    /**
     * 为教师所有可访问班级创建今日考勤记录（默认到课）。
     *
     * @return array{total: int, wechat_leave_count: int}
     */
    public function start(User $teacher): array
    {
        $classIds = $this->scope->accessibleIds($teacher);
        $count = 0;

        foreach ($classIds as $classId) {
            $count += $this->wechat->startAttendanceForClass($classId, $teacher->id, today()->toDateString());
        }

        $leaveCount = Attendance::whereIn('class_id', $classIds)
            ->whereDate('date', today())->where('source', 'wechat_work')->count();

        return ['total' => $count, 'wechat_leave_count' => $leaveCount];
    }

    /** 手动设置某学生今日考勤状态 */
    public function setStatus(User $teacher, int $studentId, string $status, ?string $remark): Attendance
    {
        $record = Attendance::whereIn('class_id', $this->scope->accessibleIds($teacher))
            ->where('student_id', $studentId)->whereDate('date', today())->firstOrFail();

        $record->update([
            'status' => $status,
            'source' => 'manual',
            'remark' => $remark,
            'sign_in_at' => $status === 'present' ? now() : $record->sign_in_at,
        ]);

        return $record;
    }

    /** 标记请假 */
    public function markLeave(User $teacher, int $studentId, string $remark): Attendance
    {
        $student = Student::whereIn('class_id', $this->scope->accessibleIds($teacher))->findOrFail($studentId);

        return $this->wechat->markManualLeave(
            $studentId,
            $student->class_id,
            $teacher->id,
            today()->toDateString(),
            $remark,
        );
    }

    /** 标记缺勤 */
    public function markAbsent(User $teacher, int $studentId, ?string $remark): Attendance
    {
        $student = Student::whereIn('class_id', $this->scope->accessibleIds($teacher))->findOrFail($studentId);

        return $this->wechat->markManualAbsent(
            $studentId,
            $student->class_id,
            $teacher->id,
            today()->toDateString(),
            $remark,
        );
    }

    /**
     * 今日考勤四态统计。
     *
     * @return array{present: int, late: int, leave: int, absent: int, rate: float, wechat_leave_count: int, manual_leave_count: int}
     */
    public function summary(User $teacher): array
    {
        $records = Attendance::whereIn('class_id', $this->scope->accessibleIds($teacher))
            ->whereDate('date', today())->get();

        $present = $records->where('status', 'present')->count();
        $late = $records->where('status', 'late')->count();
        $leave = $records->where('status', 'leave')->count();
        $absent = $records->where('status', 'absent')->count();
        $total = max($records->count(), 1);

        return [
            'present' => $present,
            'late' => $late,
            'leave' => $leave,
            'absent' => $absent,
            'rate' => round($present / $total * 100, 1),
            'wechat_leave_count' => $records->where('status', 'leave')->where('source', 'wechat_work')->count(),
            'manual_leave_count' => $records->where('status', 'leave')->where('source', 'manual')->count(),
        ];
    }
}
