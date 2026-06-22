<?php

namespace App\Services;

use App\Models\Student;
use Illuminate\Support\Facades\Redis;

class LeaderboardService
{
    // Redis Key 格式
    private function classTotalKey(int $classId): string
    {
        return "leaderboard:class:{$classId}:total";
    }

    private function classWeeklyKey(int $classId): string
    {
        return "leaderboard:class:{$classId}:weekly";
    }

    private function classPetLevelKey(int $classId): string
    {
        return "leaderboard:class:{$classId}:pet_level";
    }

    private function schoolTotalKey(int $schoolId): string
    {
        return "leaderboard:school:{$schoolId}:total";
    }

    // ========== 更新排行榜缓存 ==========

    public function updateTotalScore(int $classId, int $studentId, int $score): void
    {
        Redis::zadd($this->classTotalKey($classId), $score, $studentId);
    }

    public function updateWeeklyScore(int $classId, int $studentId, int $weeklyScore): void
    {
        Redis::zadd($this->classWeeklyKey($classId), $weeklyScore, $studentId);
    }

    public function updatePetLevel(int $classId, int $studentId, int $level): void
    {
        Redis::zadd($this->classPetLevelKey($classId), $level, $studentId);
    }

    // ========== 获取排行榜 ==========

    public function getClassTotalLeaderboard(int $classId, int $limit = 20): array
    {
        $raw = Redis::zrevrange($this->classTotalKey($classId), 0, $limit - 1, 'WITHSCORES');

        if (empty($raw)) {
            // Redis无数据时从数据库回填
            return $this->rebuildFromClassDB($classId, 'total_score', $limit);
        }

        return $this->enrichLeaderboard($classId, $raw);
    }

    public function getClassWeeklyLeaderboard(int $classId, int $limit = 20): array
    {
        $raw = Redis::zrevrange($this->classWeeklyKey($classId), 0, $limit - 1, 'WITHSCORES');

        if (empty($raw)) {
            // 从DB计算本周积分增量
            return $this->buildWeeklyFromDB($classId, $limit);
        }

        return $this->enrichLeaderboard($classId, $raw);
    }

    public function getPetLevelLeaderboard(int $classId, int $limit = 20): array
    {
        $raw = Redis::zrevrange($this->classPetLevelKey($classId), 0, $limit - 1, 'WITHSCORES');

        if (empty($raw)) {
            return $this->rebuildPetLevelFromDB($classId, $limit);
        }

        return $this->enrichLeaderboard($classId, $raw, true);
    }

    // ========== 内部方法 ==========

    private function enrichLeaderboard(int $classId, array $raw, bool $isPet = false): array
    {
        $studentIds = array_keys($raw);
        $students = Student::whereIn('id', $studentIds)->with('pet')->get()->keyBy('id');

        $result = [];
        $rank = 1;
        foreach ($raw as $studentId => $score) {
            $student = $students->get($studentId);
            if (!$student) continue;

            $entry = [
                'rank' => $rank,
                'student_id' => $studentId,
                'name' => $student->name,
                'avatar' => $student->avatar_path,
                'score' => (int) $score,
            ];

            if ($isPet && $student->pet) {
                $stage = $student->pet->currentStage();
                $entry['pet_emoji'] = $stage['emoji'];
                $entry['pet_name'] = $stage['name'];
                $entry['pet_level'] = $student->pet->level;
            }

            $result[] = $entry;
            $rank++;
        }

        return $result;
    }

    private function rebuildFromClassDB(int $classId, string $field, int $limit): array
    {
        $students = Student::where('class_id', $classId)
            ->orderBy($field, 'desc')
            ->take($limit)
            ->get();

        $result = [];
        $rank = 1;
        foreach ($students as $student) {
            // 回填Redis
            Redis::zadd($this->classTotalKey($classId), $student->total_score, $student->id);

            $result[] = [
                'rank' => $rank,
                'student_id' => $student->id,
                'name' => $student->name,
                'avatar' => $student->avatar_path,
                'score' => $student->total_score,
            ];
            $rank++;
        }

        return $result;
    }

    private function buildWeeklyFromDB(int $classId, int $limit): array
    {
        $weekStart = now()->startOfWeek();

        $students = Student::where('class_id', $classId)
            ->with(['scores' => function ($q) use ($weekStart) {
                $q->where('created_at', '>=', $weekStart)->selectRaw('student_id, SUM(amount) as weekly_total');
            }])
            ->take($limit)
            ->get();

        $result = [];
        $rank = 1;
        $sorted = $students->sortByDesc(function ($s) {
            return $s->scores->sum('amount');
        });

        foreach ($sorted as $student) {
            $weeklyTotal = $student->scores->sum('amount');
            Redis::zadd($this->classWeeklyKey($classId), $weeklyTotal, $student->id);

            $result[] = [
                'rank' => $rank,
                'student_id' => $student->id,
                'name' => $student->name,
                'score' => $weeklyTotal,
            ];
            $rank++;
        }

        return $result;
    }

    private function rebuildPetLevelFromDB(int $classId, int $limit): array
    {
        $students = Student::where('class_id', $classId)
            ->whereHas('pet')
            ->with('pet')
            ->take($limit)
            ->get()
            ->sortByDesc(fn($s) => $s->pet->level);

        $result = [];
        $rank = 1;
        foreach ($students as $student) {
            $stage = $student->pet->currentStage();
            Redis::zadd($this->classPetLevelKey($classId), $student->pet->level, $student->id);

            $result[] = [
                'rank' => $rank,
                'student_id' => $student->id,
                'name' => $student->name,
                'pet_emoji' => $stage['emoji'],
                'pet_name' => $stage['name'],
                'pet_level' => $student->pet->level,
                'score' => $student->pet->level,
            ];
            $rank++;
        }

        return $result;
    }
}
