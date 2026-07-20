<?php

declare(strict_types=1);

namespace App\Exports;

use App\Models\Score;
use App\Models\Student;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\FromCollection;
use Maatwebsite\Excel\Concerns\WithHeadings;
use Maatwebsite\Excel\Concerns\WithTitle;

class ScoresExport implements FromCollection, WithHeadings, WithTitle
{
    public function __construct(
        private readonly int $classId,
        private readonly string $className,
    ) {
    }

    public function collection(): Collection
    {
        $students = Student::where('class_id', $this->classId)
            ->where('status', 'active')
            ->with('pet')
            ->get();

        return $students->map(function ($student) {
            $positive = Score::where('student_id', $student->id)->where('amount', '>', 0)->sum('amount');
            $negative = Score::where('student_id', $student->id)->where('amount', '<', 0)->sum('amount');

            return [
                'name' => $student->name,
                'student_no' => $student->student_no ?? '',
                'total_score' => $student->total_score,
                'positive_score' => (int) $positive,
                'negative_score' => abs((int) $negative),
                // @phpstan-ignore-next-line (pet can be null)
                'pet_name' => $student->pet->name ?? '',
                // @phpstan-ignore-next-line
                'pet_level' => $student->pet->level ?? 0,
            ];
        });
    }

    public function headings(): array
    {
        return ['姓名', '学号', '总积分', '获得积分', '扣除积分', '宠物名', '宠物等级'];
    }

    public function title(): string
    {
        return $this->className;
    }
}
