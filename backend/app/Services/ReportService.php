<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use App\Models\Pet;
use App\Models\Score;
use App\Models\Student;
use App\Models\User;
use Illuminate\Support\Collection;

/**
 * 教师端报表：积分趋势、宠物等级分布、学生进度、Excel 导出。
 *
 * 只负责聚合查询与导出构建；HTTP 校验与响应包装留在控制器。
 * 所有查询都经 TeacherClassScope 限定班级作用域。
 */
class ReportService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

    /**
     * 近 N 天得分/扣分日趋势（双序列，供前端折线图）。
     *
     * @return array{labels: string[], datasets: array{label: string, data: int[]}[]}
     */
    public function scoreTrend(User $teacher, int $days): array
    {
        $days = max(1, min($days, 365));
        $start = now()->startOfDay()->subDays($days - 1);

        $rows = Score::whereIn('class_id', $this->scope->ids($teacher))
            ->where('created_at', '>=', $start)
            ->get(['amount', 'created_at'])
            ->groupBy(fn ($s) => $s->created_at->format('Y-m-d'));

        $labels = [];
        $positive = [];
        $negative = [];
        for ($d = 0; $d < $days; $d++) {
            $date = $start->copy()->addDays($d);
            $key = $date->format('Y-m-d');
            $dayRows = $rows->get($key, collect());
            $labels[] = $date->format('m/d');
            $positive[] = (int) $dayRows->where('amount', '>', 0)->sum('amount');
            $negative[] = (int) abs($dayRows->where('amount', '<', 0)->sum('amount'));
        }

        return [
            'labels' => $labels,
            'datasets' => [
                ['label' => '得分', 'data' => $positive],
                ['label' => '扣分', 'data' => $negative],
            ],
        ];
    }

    /**
     * 宠物等级分布（按 level 分组，含阶段名；多行共用阶段名，key 以 level 为准）。
     *
     * @return Collection<int, array{level: int, count: int, stage_name: string}>
     */
    public function petDistribution(User $teacher): Collection
    {
        $pets = Pet::whereIn('class_id', $this->scope->ids($teacher))->get();

        return $pets->groupBy('level')
            ->map(fn ($group, $level) => [
                'level' => (int) $level,
                'count' => (int) $group->count(),
                'stage_name' => (string) $group->first()->currentStage()['name'],
            ])
            ->sortKeys()
            ->values();
    }

    /**
     * 学生进度：有 studentId 返回单人近 50 条历史；否则返回全班每人近 10 条与涨跌趋势。
     *
     * @return array<string, mixed>|list<array<string, mixed>>
     */
    public function studentProgress(User $teacher, ?int $studentId): array
    {
        $classIds = $this->scope->ids($teacher);

        if ($studentId) {
            $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);
            $history = Score::where('student_id', $student->id)
                ->orderBy('created_at', 'desc')->take(50)->get();

            return [
                'student' => ['id' => $student->id, 'name' => $student->name, 'total_score' => $student->total_score],
                'history' => $history,
            ];
        }

        // 无 student_id 时返回班级所有学生进度列表
        $students = Student::whereIn('class_id', $classIds)->where('status', 'active')->get();

        return $students->map(function ($student) {
            $scores = Score::where('student_id', $student->id)
                ->orderBy('created_at', 'desc')->take(10)->pluck('amount');
            $change = $scores->take(5)->sum() - $scores->slice(5)->sum();

            return [
                'student_id' => $student->id,
                'student_name' => $student->name,
                'scores' => $scores->values()->toArray(),
                'trend' => $change > 5 ? 'up' : ($change < -5 ? 'down' : 'stable'),
                'change' => $change,
            ];
        })->values()->all();
    }

    /**
     * 构建报表导出下载（scores / pets / attendance）。
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse|null 不支持的类型返回 null
     */
    public function export(string $type, int $classId, ?string $date)
    {
        $class = ClassRoom::find($classId);
        $className = optional($class)->name ?? '未知班级';
        $fileName = $className . '-' . now()->format('Ymd-His');

        switch ($type) {
            case 'scores':
                return \Maatwebsite\Excel\Facades\Excel::download(
                    new \App\Exports\ScoresExport($classId, $className),
                    $fileName . '-积分报表.xlsx'
                );

            case 'pets':
                return \Maatwebsite\Excel\Facades\Excel::download(
                    new \App\Exports\PetsExport($classId, $className),
                    $fileName . '-宠物报表.xlsx'
                );

            case 'attendance':
                return \Maatwebsite\Excel\Facades\Excel::download(
                    new \App\Exports\AttendanceExport($classId, $className, $date),
                    $fileName . '-考勤报表.xlsx'
                );

            default:
                return null;
        }
    }
}
