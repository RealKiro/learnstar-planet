<?php

declare(strict_types=1);

namespace App\Exports;

use App\Models\Attendance;
use App\Models\Student;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\FromCollection;
use Maatwebsite\Excel\Concerns\WithHeadings;
use Maatwebsite\Excel\Concerns\WithTitle;

class AttendanceExport implements FromCollection, WithHeadings, WithTitle
{
    public function __construct(
        private readonly int $classId,
        private readonly string $className,
        private readonly ?string $date = null,
    ) {
    }

    public function collection(): Collection
    {
        $students = Student::where('class_id', $this->classId)
            ->where('status', 'active')
            ->pluck('id', 'name');

        $query = Attendance::where('class_id', $this->classId);
        if ($this->date) {
            $query->whereDate('created_at', $this->date);
        }

        $records = $query->get()->keyBy('student_id');

        return $students->map(function ($studentId, $studentName) use ($records) {
            $record = $records->get($studentId);

            return [
                'name' => $studentName,
                'status' => $record ? $this->statusLabel($record->status) : '未记录',
                'check_in' => $record?->check_in_time ? (string) $record->check_in_time : '',
                'source' => $record?->source ?? '',
            ];
        })->values();
    }

    public function headings(): array
    {
        return ['姓名', '状态', '签到时间', '来源'];
    }

    public function title(): string
    {
        return $this->className . '考勤';
    }

    private function statusLabel(string $status): string
    {
        return match ($status) {
            'present' => '出勤',
            'late' => '迟到',
            'leave' => '请假',
            'absent' => '缺席',
            default => $status,
        };
    }
}
