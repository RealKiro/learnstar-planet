<?php

declare(strict_types=1);

namespace App\Exports;

use App\Models\Student;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\FromCollection;
use Maatwebsite\Excel\Concerns\WithHeadings;
use Maatwebsite\Excel\Concerns\WithTitle;

class PetsExport implements FromCollection, WithHeadings, WithTitle
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
            $pet = $student->pet;

            return [
                'student_name' => $student->name,
                // @phpstan-ignore-next-line (pet can be null)
                'pet_name' => $pet->name ?? '无',
                // @phpstan-ignore-next-line
                'pet_type' => $pet->type ?? '',
                // @phpstan-ignore-next-line
                'level' => $pet->level ?? 0,
                'stage' => $pet ? ($pet->currentStage()['name'] ?? '未孵化') : '未孵化',
                // @phpstan-ignore-next-line
                'experience' => $pet->experience ?? 0,
            ];
        });
    }

    public function headings(): array
    {
        return ['学生姓名', '宠物名', '宠物系列', '等级', '进化阶段', '经验值'];
    }

    public function title(): string
    {
        return $this->className;
    }
}
