<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Notice;
use App\Models\ScoreLog;
use App\Models\Student;
use App\Services\LeaderboardService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ParentController extends Controller
{
    /**
     * 家长首页 — 返回所有孩子的概览
     */
    public function home(Request $request): JsonResponse
    {
        $parent = $request->user();
        $children = $parent->children()->with(['classRoom', 'pet'])->get();

        if ($children->isEmpty()) {
            return response()->json(['data' => ['children' => [], 'message' => '尚未绑定孩子']]);
        }

        $data = $children->map(function (Student $child) {
            $pet = $child->pet;
            $classRank = null;
            if ($child->class_id) {
                $leaderboard = app(LeaderboardService::class)->getClassTotalLeaderboard($child->class_id, 100);
                $entry = collect($leaderboard)->firstWhere('student_id', $child->id);
                $classRank = $entry['rank'] ?? null;
            }

            return [
                'id' => $child->id,
                'name' => $child->name,
                'student_no' => $child->student_no,
                'class_name' => $child->classRoom?->name,
                'grade' => $child->classRoom?->grade,
                'total_score' => $child->total_score,
                'class_rank' => $classRank,
                'pet_name' => $pet?->name,
                'pet_level' => $pet?->level,
                'pet_stage' => $pet?->currentStage()['name'] ?? '星尘',
                'pet_emoji' => $pet?->currentStage()['emoji'] ?? '🌟',
                'pet_mood' => $pet?->mood,
            ];
        });

        return response()->json(['data' => ['children' => $data]]);
    }

    /**
     * 积分详情 — 返回孩子的积分统计
     */
    public function scoreDetail(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student) {
            return response()->json(['message' => '未找到该学生或无权查看'], 404);
        }

        $todayScores = $student->scores()->whereDate('created_at', today())->sum('amount');
        $weekScores = $student->scores()
            ->whereBetween('created_at', [now()->startOfWeek(), now()->endOfWeek()])
            ->sum('amount');

        $classRank = null;
        if ($student->class_id) {
            $leaderboard = app(LeaderboardService::class)->getClassTotalLeaderboard($student->class_id, 100);
            $entry = collect($leaderboard)->firstWhere('student_id', $student->id);
            $classRank = $entry['rank'] ?? null;
        }

        return response()->json([
            'data' => [
                'student_id' => $student->id,
                'student_name' => $student->name,
                'total_score' => $student->total_score,
                'today_score' => $todayScores,
                'week_score' => $weekScores,
                'class_rank' => $classRank,
                'class_name' => $student->classRoom?->name,
            ],
        ]);
    }

    /**
     * 积分历史 — 返回孩子的积分变动记录
     */
    public function scoreHistory(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student) {
            return response()->json(['message' => '未找到该学生或无权查看'], 404);
        }

        $logs = ScoreLog::where('student_id', $student->id)
            ->orderByDesc('created_at')
            ->limit(50)
            ->get();

        $data = $logs->map(fn (ScoreLog $log) => [
            'id' => $log->id,
            'description' => $log->description,
            'balance_before' => $log->balance_before,
            'balance_after' => $log->balance_after,
            'change' => $log->balance_after - $log->balance_before,
            'created_at' => $log->created_at?->toIso8601String(),
        ]);

        return response()->json(['data' => $data]);
    }

    /**
     * 成长日志 — 返回孩子的综合成长记录
     */
    public function growthLog(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student) {
            return response()->json(['message' => '未找到该学生或无权查看'], 404);
        }

        // 积分变动记录
        $scoreLogs = ScoreLog::where('student_id', $student->id)
            ->orderByDesc('created_at')
            ->limit(20)
            ->get()
            ->map(fn (ScoreLog $log) => [
                'type' => 'score',
                'description' => $log->description,
                'change' => $log->balance_after - $log->balance_before,
                'balance' => $log->balance_after,
                'created_at' => $log->created_at?->toIso8601String(),
            ]);

        // 宠物升级记录（通过 level 变化推断，简化为当前状态）
        $pet = $student->pet;
        $petEvents = collect();
        if ($pet && $pet->level > 0) {
            $petEvents->push([
                'type' => 'pet',
                'description' => '宠物进化为「' . ($pet->currentStage()['name'] ?? '未知') . '」',
                'change' => null,
                'balance' => $pet->level,
                'created_at' => $pet->updated_at?->toIso8601String(),
            ]);
        }

        $allEvents = $scoreLogs->concat($petEvents)
            ->sortByDesc('created_at')
            ->values()
            ->take(30);

        return response()->json(['data' => $allEvents]);
    }

    /**
     * 成长时间线 — 返回孩子的成长里程碑
     */
    public function growthTimeline(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student) {
            return response()->json(['message' => '未找到该学生或无权查看'], 404);
        }

        // 第一条积分记录（入学/加入）
        $firstScore = $student->scores()->oldest()->first();
        // 最高单次积分
        $topScore = $student->scores()->orderByDesc('amount')->first();
        // 宠物状态
        $pet = $student->pet;

        $timeline = collect();

        if ($firstScore) {
            $timeline->push([
                'time' => $firstScore->created_at?->toIso8601String(),
                'title' => '开始积分之旅',
                'description' => '获得第一次积分，开启学趣星球探索',
                'icon' => '🎯',
            ]);
        }

        if ($topScore && $topScore->amount > 0) {
            $timeline->push([
                'time' => $topScore->created_at?->toIso8601String(),
                'title' => '最高单次积分',
                'description' => "获得 {$topScore->amount} 积分 — {$topScore->reason}",
                'icon' => '🏆',
            ]);
        }

        if ($pet) {
            $timeline->push([
                'time' => $pet->created_at?->toIso8601String(),
                'title' => '宠物伙伴降临',
                'description' => "「{$pet->name}」成为了专属宠物伙伴",
                'icon' => '🥚',
            ]);

            if ($pet->level > 0) {
                $stage = $pet->currentStage();
                $timeline->push([
                    'time' => $pet->updated_at?->toIso8601String(),
                    'title' => "宠物进化为「{$stage['name']}」",
                    'description' => $stage['emoji'] . ' ' . $stage['name'] . ' · 等级 ' . $pet->level,
                    'icon' => $stage['emoji'],
                ]);
            }
        }

        if ($student->total_score > 100) {
            $timeline->push([
                'time' => null,
                'title' => '积分破百',
                'description' => "累计积分达到 {$student->total_score} 分",
                'icon' => '💯',
            ]);
        }

        $timeline = $timeline->sortByDesc('time')->values();

        return response()->json(['data' => $timeline]);
    }

    /**
     * 宠物状态 — 返回孩子的宠物详情
     */
    public function petStatus(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student) {
            return response()->json(['message' => '未找到该学生或无权查看'], 404);
        }

        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['data' => null, 'message' => '尚未拥有宠物']);
        }

        $stage = $pet->currentStage();

        return response()->json([
            'data' => [
                'id' => $pet->id,
                'name' => $pet->name,
                'type' => $pet->type,
                'level' => $pet->level,
                'experience' => $pet->experience,
                'mood' => $pet->mood,
                'stage_name' => $stage['name'] ?? '星尘',
                'stage_emoji' => $stage['emoji'] ?? '🌟',
                'last_fed_at' => $pet->last_fed_at?->toIso8601String(),
                'student_name' => $student->name,
            ],
        ]);
    }

    /**
     * 喂养宠物
     */
    public function feedPet(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student) {
            return response()->json(['message' => '未找到该学生或无权查看'], 404);
        }

        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['message' => '尚未拥有宠物，无法喂养'], 400);
        }

        try {
            $pet->feed();
            $stage = $pet->fresh()->currentStage();

            return response()->json([
                'message' => '喂养成功',
                'data' => [
                    'name' => $pet->name,
                    'level' => $pet->level,
                    'experience' => $pet->experience,
                    'mood' => $pet->mood,
                    'stage_name' => $stage['name'] ?? '星尘',
                    'stage_emoji' => $stage['emoji'] ?? '🌟',
                ],
            ]);
        } catch (\Throwable $e) {
            return response()->json(['message' => '喂养失败：' . $e->getMessage()], 500);
        }
    }

    /**
     * 排行榜 — 返回孩子所在班级的排行榜
     */
    public function ranking(Request $request): JsonResponse
    {
        $parent = $request->user();
        $studentId = (int) $request->query('student_id', 0);

        $student = $this->getChild($parent, $studentId);
        if (!$student || !$student->class_id) {
            return response()->json(['data' => []]);
        }

        $leaderboard = app(LeaderboardService::class)->getClassTotalLeaderboard($student->class_id, 50);

        // 标记当前孩子
        $data = collect($leaderboard)->map(function ($entry) use ($student) {
            $entry['is_mine'] = $entry['student_id'] === $student->id;

            return $entry;
        });

        return response()->json(['data' => $data]);
    }

    /**
     * 通知列表 — 返回孩子所在班级的已发布通知
     */
    public function listNotices(Request $request): JsonResponse
    {
        $parent = $request->user();
        $children = $parent->children;
        $classIds = $children->pluck('class_id')->unique()->filter()->toArray();

        if (empty($classIds)) {
            return response()->json(['data' => []]);
        }

        $notices = Notice::whereIn('class_id', $classIds)
            ->where('is_published', true)
            ->orderByDesc('published_at')
            ->limit(30)
            ->get();

        $data = $notices->map(fn (Notice $notice) => [
            'id' => $notice->id,
            'title' => $notice->title,
            'content' => $notice->content,
            'type' => $notice->type,
            'published_at' => $notice->published_at?->toIso8601String(),
            'class_name' => $notice->classRoom?->name,
        ]);

        return response()->json(['data' => $data]);
    }

    /**
     * 通知详情
     */
    public function readNotice(Request $request, int $id): JsonResponse
    {
        $parent = $request->user();
        $children = $parent->children;
        $classIds = $children->pluck('class_id')->unique()->filter()->toArray();

        $notice = Notice::whereIn('class_id', $classIds)
            ->where('is_published', true)
            ->find($id);

        if (!$notice) {
            return response()->json(['message' => '通知不存在或无权查看'], 404);
        }

        return response()->json([
            'data' => [
                'id' => $notice->id,
                'title' => $notice->title,
                'content' => $notice->content,
                'type' => $notice->type,
                'published_at' => $notice->published_at?->toIso8601String(),
                'class_name' => $notice->classRoom?->name,
                'publisher_name' => $notice->publisher?->name,
            ],
        ]);
    }

    /**
     * 获取家长的孩子（支持 student_id 选择，默认取第一个）
     */
    private function getChild($parent, int $studentId): ?Student
    {
        if ($studentId > 0) {
            return $parent->children()->where('id', $studentId)->first();
        }

        return $parent->children()->first();
    }
}
