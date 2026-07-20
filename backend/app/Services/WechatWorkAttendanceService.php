<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Attendance;
use App\Models\School;
use App\Models\Student;
use App\Models\ThirdPartyBinding;
use App\Models\WechatWorkLeaveRecord;
use Illuminate\Support\Facades\Log;

class WechatWorkAttendanceService
{
    public function __construct(private readonly WechatWorkService $wework)
    {
    }

    public function syncAll(?string $date = null): array
    {
        $date ??= now()->toDateString();
        $r = ['synced' => 0,'skipped' => 0,'failed' => 0];
        foreach (School::where('status', 'active')->get() as $s) {
            try {
                $x = $this->syncForSchool($s->id, $date);
                $r['synced'] += $x['synced'];
                $r['skipped'] += $x['skipped'];
            } catch (\Throwable $e) {
                Log::error('同步失败', ['sid' => $s->id,'err' => $e->getMessage()]);
                $r['failed']++;
            }
        }

        return $r;
    }

    public function syncForSchool(int $schoolId, string $date): array
    {
        if (empty(config('wechat-work.corp_id'))) {
            return ['synced' => 0,'skipped' => 0];
        }
        $ds = strtotime($date . ' 00:00:00');
        $de = strtotime($date . ' 23:59:59');
        $sns = $this->wework->getLeaveApprovalSpNos($schoolId, $ds, $de);
        $syn = 0;
        $skp = 0;
        foreach ($sns as $spNo) {
            if (WechatWorkLeaveRecord::where('sp_no', $spNo)->exists()) {
                $skp++;
                continue;
            }
            $d = $this->wework->getApprovalDetail($schoolId, $spNo);
            if ($d === null) {
                continue;
            }
            $stu = $this->matchStudent($d['submitter_userid'], $d['student_name']);
            $rec = WechatWorkLeaveRecord::create([
                'school_id' => $schoolId,'class_id' => $stu?->class_id,'student_id' => $stu?->id,'parent_wework_userid' => $d['submitter_userid'],'student_name_from_wework' => $d['student_name'],'sp_no' => $d['sp_no'],'leave_start_date' => date('Y-m-d', $d['leave_start'] ?: $ds),'leave_end_date' => date('Y-m-d', $d['leave_end'] ?: $de),'leave_type' => $d['leave_type'],'reason' => $d['reason'],'approve_status' => $d['approve_status'] === 2 ? 'approved' : 'pending','approved_at' => $d['approved_at'] ? date('Y-m-d H:i:s', $d['approved_at']) : null,'raw_data' => $d,'synced_at' => null,
            ]);
            if ($stu !== null) {
                $this->applyLeave($stu, $rec, $date);
            }
            $syn++;
        }

        return ['synced' => $syn,'skipped' => $skp];
    }

    public function handleWebhookCallback(int $schoolId, array $data): void
    {
        $spNo = $data['sp_no'] ?? '';
        if (empty($spNo) || ($data['sp_status'] ?? 1) !== 2) {
            return;
        }
        if (WechatWorkLeaveRecord::where('sp_no', $spNo)->exists()) {
            return;
        }
        $d = $this->wework->getApprovalDetail($schoolId, $spNo);
        if ($d === null) {
            return;
        }
        $stu = $this->matchStudent($d['submitter_userid'], $d['student_name']);
        $t = now()->toDateString();
        $rec = WechatWorkLeaveRecord::create([
            'school_id' => $schoolId,'class_id' => $stu?->class_id,'student_id' => $stu?->id,'parent_wework_userid' => $d['submitter_userid'],'student_name_from_wework' => $d['student_name'],'sp_no' => $d['sp_no'],'leave_start_date' => date('Y-m-d', $d['leave_start'] ?: strtotime($t)),'leave_end_date' => date('Y-m-d', $d['leave_end'] ?: strtotime($t)),'leave_type' => $d['leave_type'],'reason' => $d['reason'],'approve_status' => 'approved','approved_at' => now(),'raw_data' => $d,'synced_at' => null,
        ]);
        if ($stu !== null) {
            $this->applyLeave($stu, $rec, $t);
        }
    }

    public function applyLeave(Student $stu, WechatWorkLeaveRecord $rec, string $date): void
    {
        $ex = Attendance::where('student_id', $stu->id)->whereDate('date', $date)->first();
        if ($ex === null) {
            Attendance::create(['class_id' => $stu->class_id,'student_id' => $stu->id,'teacher_id' => $stu->classRoom?->teacher_id,'date' => $date,'status' => 'leave','source' => 'wechat_work','remark' => $rec->reason,'leave_record_id' => $rec->id]);
            $rec->update(['synced_at' => now()]);

            return;
        }
        if ($ex->source === 'manual') {
            Log::info('跳过手动修改', ['aid' => $ex->id,'sid' => $stu->id,'sp' => $rec->sp_no]);

            return;
        }
        $ex->update(['status' => 'leave','source' => 'wechat_work','remark' => $rec->reason,'leave_record_id' => $rec->id]);
        $rec->update(['synced_at' => now()]);
    }

    public function startAttendanceForClass(int $classId, int $teacherId, string $date): int
    {
        $stus = Student::where('class_id', $classId)->where('status', 'active')->get();
        $c = 0;
        foreach ($stus as $s) {
            $a = Attendance::firstOrCreate(['class_id' => $classId,'student_id' => $s->id,'date' => $date], ['teacher_id' => $teacherId,'status' => 'present','source' => 'auto','remark' => null]);
            if ($a->wasRecentlyCreated === false && $a->source !== 'manual') {
                $a->update(['status' => 'present','source' => 'auto','remark' => null,'leave_record_id' => null]);
            }
            $c++;
        }
        foreach (WechatWorkLeaveRecord::getUnsyncedForDate($date) as $rec) {
            if ($rec->student_id === null) {
                continue;
            }
            $s = Student::find($rec->student_id);
            if ($s !== null && $s->class_id === $classId) {
                $this->applyLeave($s, $rec, $date);
            }
        }

        return $c;
    }

    public function markManualLeave(int $studentId, int $classId, int $teacherId, string $date, string $remark): Attendance
    {
        $a = Attendance::firstOrCreate(['class_id' => $classId,'student_id' => $studentId,'date' => $date], ['teacher_id' => $teacherId,'status' => 'present','source' => 'auto']);
        $a->update(['status' => 'leave','source' => 'manual','remark' => $remark,'leave_record_id' => null]);

        return $a->fresh();
    }

    public function markManualAbsent(int $studentId, int $classId, int $teacherId, string $date, ?string $remark): Attendance
    {
        $a = Attendance::firstOrCreate(['class_id' => $classId,'student_id' => $studentId,'date' => $date], ['teacher_id' => $teacherId,'status' => 'present','source' => 'auto']);
        $a->update(['status' => 'absent','source' => 'manual','remark' => $remark ?? '未联系到家长，建议后续跟进']);

        return $a->fresh();
    }

    private function matchStudent(string $pwid, ?string $hint): ?Student
    {
        $b = ThirdPartyBinding::where('platform', 'wechat_work')->where('platform_id', $pwid)->first();
        if ($b === null) {
            return null;
        }
        $q = Student::where('parent_id', $b->user_id)->where('status', 'active');
        if ($hint !== null) {
            $s = $q->where('name', $hint)->first();
            if ($s !== null) {
                return $s;
            }
            $s = $q->where('name', 'like', "%{$hint}%")->first();
            if ($s !== null) {
                return $s;
            }
        }

        return $q->first();
    }
}

