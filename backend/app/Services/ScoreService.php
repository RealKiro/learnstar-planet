<?php

declare(strict_types=1);

namespace App\Services;

use App\Events\ScoreChanged;
use App\Models\Score;
use App\Models\ScoreLog;
use App\Models\ScoreRule;
use App\Models\Student;
use Illuminate\Support\Facades\DB;

class ScoreService
{
    // ========== 给单个学生加减分 ==========

    public function giveScore(Student $student, int $amount, string $reason, int $givenBy, ?int $scoreRuleId = null): Score
    {
        return DB::transaction(function () use ($student, $amount, $reason, $givenBy, $scoreRuleId) {
            $balanceBefore = $student->total_score;

            // 创建积分记录
            $score = Score::create([
                'student_id' => $student->id,
                'class_id' => $student->class_id,
                'score_rule_id' => $scoreRuleId,
                'amount' => $amount,
                'reason' => $reason,
                'given_by' => $givenBy,
            ]);

            // 更新学生总积分
            $student->update(['total_score' => $balanceBefore + $amount]);
            $balanceAfter = $student->total_score;

            // 记录积分日志
            ScoreLog::create([
                'student_id' => $student->id,
                'score_id' => $score->id,
                'balance_before' => $balanceBefore,
                'balance_after' => $balanceAfter,
                'description' => $reason,
            ]);

            // 如果是加分，宠物也获得经验
            if ($amount > 0 && $student->pet) {
                $student->pet->addExperience($amount);
            }

            // 广播积分变化事件（实时推送给家长端）
            event(new ScoreChanged($student->id, $amount, $reason));

            return $score;
        });
    }

    // ========== 批量加分 ==========

    public function batchGiveScore(array $studentIds, int $amount, string $reason, int $givenBy, ?int $scoreRuleId = null): array
    {
        $results = [];

        foreach ($studentIds as $studentId) {
            $student = Student::find($studentId);

            if ($student) {
                $results[] = $this->giveScore($student, $amount, $reason, $givenBy, $scoreRuleId);
            }
        }

        return $results;
    }

    // ========== 使用规则加分 ==========

    public function giveScoreByRule(Student $student, ScoreRule $rule, int $givenBy): Score
    {
        return $this->giveScore(
            $student,
            $rule->amount,
            $rule->name,
            $givenBy,
            $rule->id
        );
    }

    // ========== 学生积分历史 ==========

    public function getScoreHistory(Student $student, int $perPage = 20): \Illuminate\Pagination\LengthAwarePaginator
    {
        return Score::where('student_id', $student->id)
            ->with(['scoreRule', 'giver'])
            ->orderBy('created_at', 'desc')
            ->paginate($perPage);
    }

    // ========== 班级积分统计 ==========

    public function getClassScoreSummary(int $classId): array
    {
        $students = Student::where('class_id', $classId)
            ->orderBy('total_score', 'desc')
            ->get();

        return [
            'total_students' => $students->count(),
            'total_score' => $students->sum('total_score'),
            'average_score' => $students->avg('total_score'),
            'top_scorers' => $students->take(5)->map(fn ($s) => [
                'id' => $s->id,
                'name' => $s->name,
                'score' => $s->total_score,
            ]),
            'recent_scores' => Score::where('class_id', $classId)
                ->with(['student', 'giver'])
                ->orderBy('created_at', 'desc')
                ->take(10)
                ->get(),
        ];
    }
}
