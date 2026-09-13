<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Grade;
use App\Models\Student;
use App\Models\User;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;

/**
 * 成绩读写与统计。
 *
 * 只负责数据操作；HTTP 校验与响应包装留在控制器。
 * 所有查询都经 TeacherClassScope 限定班级作用域。
 */
class GradeService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

    /** 成绩列表（可按考试名/科目过滤，每页 50 条，按分数降序） */
    public function paginate(User $teacher, ?string $examName, ?string $subject): LengthAwarePaginator
    {
        $query = Grade::whereIn('class_id', $this->scope->ids($teacher))->with('student:id,name');

        if ($examName) {
            $query->where('exam_name', $examName);
        }
        if ($subject) {
            $query->where('subject', $subject);
        }

        return $query->orderBy('score', 'desc')->paginate(50);
    }

    /**
     * 批量录入成绩（缺失/越权的学生自动跳过）。
     *
     * @param  array<int, array{student_id: int, score: numeric-string|int|float}>  $grades
     * @return int  实际录入条数
     */
    public function input(User $teacher, array $grades, string $examName, string $subject): int
    {
        $classIds = $this->scope->ids($teacher);
        $count = 0;

        foreach ($grades as $g) {
            $student = Student::whereIn('class_id', $classIds)->find($g['student_id']);
            if (!$student) {
                continue;
            }

            Grade::updateOrCreate(
                [
                    'class_id' => $student->class_id,
                    'student_id' => $g['student_id'],
                    'exam_name' => $examName,
                    'subject' => $subject,
                ],
                [
                    'teacher_id' => $teacher->id,
                    'score' => $g['score'],
                ]
            );
            $count++;
        }

        return $count;
    }

    /**
     * 单个考试/科目的成绩统计（针对教师第一个班级）。
     *
     * @return array<string, mixed>
     */
    public function stats(User $teacher, string $examName, string $subject): array
    {
        return Grade::classStats($this->scope->ids($teacher)->first(), $examName, $subject);
    }

    /**
     * 单个考试/科目的分数段分布（针对教师第一个班级）。
     *
     * @return array<string, mixed>
     */
    public function distribution(User $teacher, string $examName, string $subject): array
    {
        return Grade::scoreDistribution($this->scope->ids($teacher)->first(), $examName, $subject);
    }
}
