<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Broadcast;
use App\Models\ClassRoom;
use App\Models\ClassRoomTeacher;
use App\Models\Notice;
use App\Models\Score;
use App\Models\ShopRedemption;
use App\Models\Student;
use App\Models\User;
use Illuminate\Support\Collection;

/**
 * 教师端仪表盘/我的班级服务。
 */
class DashboardService
{
    /**
     * 教师关联班级列表（API 机器人返回本校全部启用班级）。
     *
     * @return array<int, array<string, mixed>>
     */
    public function classesFor(User $teacher): array
    {
        // API 机器人账号：返回本校全部启用中的班级（供外部系统枚举可用班级）
        if ($teacher->isApiBot()) {
            return ClassRoom::where('school_id', $teacher->school_id)
                ->where('status', 'active')
                ->orderBy('id')
                ->get(['id', 'name', 'grade'])
                ->map(fn (ClassRoom $c) => [
                    'class_id' => $c->id,
                    'class_name' => $c->name,
                    'grade' => $c->grade,
                    'role' => 'api_bot',
                ])
                ->all();
        }

        return ClassRoomTeacher::where('user_id', $teacher->id)
            ->with('classRoom:id,name,grade')
            ->get()
            ->map(fn ($a) => [
                'class_id' => $a->class_room_id,
                'class_name' => $a->classRoom?->name,
                'grade' => $a->classRoom?->grade,
                'role' => $a->role,
            ])
            ->all();
    }

    /**
     * 切换当前激活班级。返回 false 表示未分配该班级（控制器映射 403）。
     */
    public function switchTo(User $teacher, int $classId): bool
    {
        $isAssigned = $teacher->isApiBot()
            ? ClassRoom::where('school_id', $teacher->school_id)->where('id', $classId)->where('status', 'active')->exists()
            : ClassRoomTeacher::where('user_id', $teacher->id)
                ->where('class_room_id', $classId)
                ->exists();

        if (!$isAssigned) {
            return false;
        }

        // 存储当前班级到用户设置，后续所有 API 都使用此班级
        $teacher->setSetting('active_class_id', $classId);

        return true;
    }

    /**
     * 空仪表盘（无可管理班级或班级不存在时）。
     *
     * @return array<string, mixed>
     */
    private function emptyPayload(): array
    {
        return [
            'class_name' => '', 'grade' => '', 'student_count' => 0,
            'total_score' => 0, 'avg_pet_level' => 0, 'peak_count' => 0, 'weekly_score' => 0,
            'pending_redemptions' => 0, 'star_student' => null, 'top5' => [], 'recent_news' => [],
        ];
    }

    /**
     * 教师仪表盘聚合（当前激活班级统计 + Top5 + 星光学生 + 最新动态）。
     *
     * @param Collection<int, int> $classIds
     * @return array<string, mixed>
     */
    public function forTeacher(User $teacher, Collection $classIds): array
    {
        if ($classIds->isEmpty()) {
            return $this->emptyPayload();
        }

        // 教师当前激活班级（未设置时取第一个）
        $activeClassId = $teacher->getSetting('active_class_id') ?: $classIds->first();
        $class = ClassRoom::find($activeClassId) ?? ClassRoom::find($classIds->first());
        if (!$class) {
            return $this->emptyPayload();
        }

        $pendingRedemptions = ShopRedemption::whereIn('class_id', $classIds)
            ->where('status', 'pending')
            ->count();

        $students = Student::where('class_id', $class->id)->where('status', 'active')->with('pet')->get();
        $totalScore = $students->sum('total_score');
        $count = $students->count();
        $avgLevel = $count > 0 ? round($students->avg(fn ($s) => $s->pet->level ?? 0), 1) : 0;
        $peakCount = $students->filter(fn ($s) => $s->pet && $s->pet->level >= 10)->count();
        $sorted = $students->sortByDesc('total_score')->values();
        $top5 = $sorted->take(5)->map(fn ($s) => [
            'name' => $s->name,
            'student_no' => $s->student_no,
            'score' => $s->total_score,
            'pet_name' => $s->pet->name ?? '',
            'pet_species' => $s->pet->species ?? '',
            'pet_level' => $s->pet->level ?? 0,
        ]);
        $starStudent = $sorted->first();
        $recentNews = Score::whereIn('student_id', $students->pluck('id'))
            ->with('student:id,name')->orderBy('created_at', 'desc')->take(20)->get()
            ->map(fn ($s) => [
                'icon' => $s->amount > 0 ? '🎉' : '📝',
                'text' => ($s->student->name ?? '同学') . ' ' . ($s->amount > 0 ? '+' . $s->amount : $s->amount) . '分 — ' . ($s->reason ?? ''),
            ])
            ->unique('text')->take(5)->values();

        return [
            'class_name' => $class->name,
            'grade' => $class->grade,
            'student_count' => $count,
            'total_score' => (int) $totalScore,
            'avg_pet_level' => $avgLevel,
            'peak_count' => $peakCount,
            'weekly_score' => (int) Score::whereIn('student_id', $students->pluck('id'))
                ->where('created_at', '>=', now()->startOfWeek())->sum('amount'),
            'pending_redemptions' => $pendingRedemptions,
            'star_student' => $starStudent ? [
                'name' => $starStudent->name,
                'student_no' => $starStudent->student_no,
                'pet_name' => $starStudent->pet->name ?? '',
                'pet_species' => $starStudent->pet->species ?? '',
                'pet_level' => $starStudent->pet->level ?? 0,
                'score' => $starStudent->total_score,
            ] : null,
            'top5' => $top5,
            'recent_news' => $recentNews,
        ];
    }
}
