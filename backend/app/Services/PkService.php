<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use App\Models\Score;
use App\Models\Student;
use App\Models\User;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Cache;

/**
 * 年级战场（PK）：同年级班级排行、本班统计、挑战记录。
 *
 * 只负责聚合查询与缓存写入；HTTP 校验与响应包装留在控制器。
 */
class PkService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

    /**
     * 计算单个班级的 PK 指标（总积分/人数/均级/巅峰数/本周增长）
     *
     * @return array{totalScore: int, studentCount: int, avgLevel: float, peakCount: int, weekGrowth: int}
     */
    private function metricsFor(int $classId): array
    {
        $students = Student::where('class_id', $classId)
            ->where('status', 'active')
            ->with('pet')
            ->get();

        $count = $students->count();
        $avgLevel = $count > 0 ? $students->avg(function ($s) {
            return $s->pet ? $s->pet->level : 0;
        }) : 0;
        $peakCount = $students->filter(function ($s) {
            return $s->pet && $s->pet->level >= 8;
        })->count();

        $weeklyScore = Score::whereIn('student_id', $students->pluck('id'))
            ->where('created_at', '>=', now()->startOfWeek())
            ->sum('amount');

        return [
            'totalScore' => (int) $students->sum('total_score'),
            'studentCount' => $count,
            'avgLevel' => round($avgLevel, 1),
            'peakCount' => $peakCount,
            'weekGrowth' => (int) $weeklyScore,
        ];
    }

    /**
     * 同年级各班 PK 排行榜（按总积分降序，isOwn 标记本班）。
     *
     * @return Collection<int, array{class_id: int, name: string, isOwn: bool, totalScore: int, studentCount: int, avgLevel: float, peakCount: int, weekGrowth: int}>
     */
    public function leaderboard(User $teacher): Collection
    {
        $classIds = $this->scope->ids($teacher);

        if ($classIds->isEmpty()) {
            return collect();
        }

        $myClass = ClassRoom::find($classIds->first());
        if (!$myClass) {
            return collect();
        }

        $gradeClasses = ClassRoom::where('grade', $myClass->grade)
            ->where('status', 'active')
            ->get();

        return $gradeClasses->map(function (ClassRoom $class) use ($myClass) {
            return array_merge([
                'class_id' => (int) $class->id,
                'name' => (string) $class->name,
                'isOwn' => $class->id === $myClass->id,
            ], $this->metricsFor((int) $class->id));
        })->sortByDesc('totalScore')->values();
    }

    /**
     * 本班 PK 统计（含同年级内排名）。
     *
     * @return array<string, int|float>
     */
    public function myStats(User $teacher): array
    {
        $classIds = $this->scope->ids($teacher);

        if ($classIds->isEmpty()) {
            return ['totalScore' => 0, 'avgLevel' => 0, 'peakCount' => 0, 'weekGrowth' => 0, 'rank' => 0];
        }

        $classId = $classIds->first();
        $metrics = $this->metricsFor($classId);

        // 计算排名（同年级内，按总积分）
        $rank = 0;
        $myClass = ClassRoom::find($classId);
        if ($myClass) {
            $allClasses = ClassRoom::where('grade', $myClass->grade)
                ->where('status', 'active')
                ->get();

            $classScores = [];
            foreach ($allClasses as $c) {
                $cStudents = Student::where('class_id', $c->id)
                    ->where('status', 'active')->get();
                $classScores[$c->id] = $cStudents->sum('total_score');
            }
            arsort($classScores);
            $rank = array_search($classId, array_keys($classScores), true);
            if ($rank !== false) {
                $rank++;
            }
        }

        return array_merge($metrics, ['rank' => $rank]);
    }

    /**
     * 记录 PK 挑战事件（缓存 7 天）。
     *
     * @return array{target_class: string, expires_at: string}
     */
    public function challenge(int $myClassId, ClassRoom $targetClass): array
    {
        $challengeKey = 'pk_challenge:' . $myClassId . ':' . $targetClass->id;
        $challengerClass = ClassRoom::find($myClassId);
        $expiresAt = now()->addDays(7);

        Cache::put($challengeKey, [
            'challenger_class_id' => $myClassId,
            'target_class_id' => $targetClass->id,
            'challenger_name' => ($challengerClass ? $challengerClass->name : '未知'),
            'target_name' => $targetClass->name,
            'challenged_at' => now()->toDateTimeString(),
            'expires_at' => $expiresAt->toDateTimeString(),
            'status' => 'active',
        ], $expiresAt);

        return [
            'target_class' => $targetClass->name,
            'expires_at' => $expiresAt->toDateTimeString(),
        ];
    }
}
