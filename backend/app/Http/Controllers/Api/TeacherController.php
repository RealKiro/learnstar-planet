<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\Score;
use App\Models\ScoreRule;
use App\Models\Student;
use App\Services\LeaderboardService;
use App\Services\ScoreService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class TeacherController extends Controller
{
    public function __construct(
        private readonly ScoreService $scoreService,
        private readonly LeaderboardService $leaderboardService,
    ) {
    }

    // ============================================================
    // Dashboard
    // ============================================================

    public function dashboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        if ($classIds->isEmpty()) {
            return response()->json(['data' => [
                'student_count' => 0,
                'weekly_score' => 0,
                'avg_pet_level' => 0,
                'pending_redemptions' => 0,
                'recent_scores' => [],
                'top_students' => [],
            ]]);
        }

        $studentCount = Student::whereIn('class_id', $classIds)->where('status', 'active')->count();
        $weeklyScores = Score::whereIn('class_id', $classIds)
            ->where('created_at', '>=', now()->startOfWeek())
            ->sum('amount');
        $avgPetLevel = (int) Student::whereIn('class_id', $classIds)
            ->whereHas('pet')
            ->withAvg('pet', 'level')
            ->get()
            ->avg('pet_avg_level');
        $pendingRedemptions = \App\Models\ShopRedemption::whereIn('class_id', $classIds)
            ->where('status', 'pending')
            ->count();
        $recentScores = Score::whereIn('class_id', $classIds)
            ->with('student')
            ->orderBy('created_at', 'desc')
            ->take(10)
            ->get()
            ->map(fn (Score $s) => [
                'student_name' => $s->student?->name,
                'points' => $s->amount,
                'reason' => $s->reason,
                'created_at' => $s->created_at?->diffForHumans(),
            ]);
        $topStudents = Student::whereIn('class_id', $classIds)
            ->orderBy('total_score', 'desc')
            ->take(5)
            ->get()
            ->map(fn (Student $s) => [
                'student_name' => $s->name,
                'score' => $s->total_score,
            ]);

        return response()->json(['data' => [
            'student_count' => $studentCount,
            'weekly_score' => (int) $weeklyScores,
            'avg_pet_level' => $avgPetLevel,
            'pending_redemptions' => $pendingRedemptions,
            'recent_scores' => $recentScores,
            'top_students' => $topStudents,
        ]]);
    }

    // ============================================================
    // Student Management
    // ============================================================

    public function listStudents(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $query = Student::whereIn('class_id', $classIds)
            ->with('classRoom:id,name,grade');

        if ($request->has('search')) {
            $search = $request->input('search');
            $query->where(function ($q) use ($search) {
                $q->where('name', 'like', "%{$search}%")
                  ->orWhere('student_no', 'like', "%{$search}%");
            });
        }

        $students = $query->orderBy('name')->paginate(50);

        return response()->json([
            'data' => $students->items(),
            'meta' => [
                'current_page' => $students->currentPage(),
                'last_page' => $students->lastPage(),
                'per_page' => $students->perPage(),
                'total' => $students->total(),
            ],
        ]);
    }

    public function importStudents(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $students = $request->input('students', []);
        $imported = 0;

        foreach ($students as $data) {
            if (empty($data['name']) || empty($data['class_name'])) {
                continue;
            }
            $classRoom = ClassRoom::whereIn('id', $classIds)
                ->where('name', $data['class_name'])
                ->first();
            if (!$classRoom) {
                continue;
            }
            Student::create([
                'class_id' => $classRoom->id,
                'name' => $data['name'],
                'gender' => $data['gender'] ?? '未知',
                'student_no' => $data['student_no'] ?? null,
                'total_score' => 0,
                'status' => 'active',
            ]);
            $imported++;
        }

        return response()->json([
            'message' => "成功导入 {$imported} 名学生",
            'data' => ['imported_count' => $imported],
        ]);
    }

    public function updateStudent(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);

        $student->update($request->only(['name', 'gender', 'student_no']));

        return response()->json(['message' => '更新成功', 'data' => $student]);
    }

    // ============================================================
    // Score Management
    // ============================================================

    public function scoreSummary(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $total = Score::whereIn('class_id', $classIds)->sum('amount');
        $today = Score::whereIn('class_id', $classIds)
            ->whereDate('created_at', today())
            ->sum('amount');
        $week = Score::whereIn('class_id', $classIds)
            ->where('created_at', '>=', now()->startOfWeek())
            ->sum('amount');

        return response()->json(['data' => [
            'total' => (int) $total,
            'today' => (int) $today,
            'this_week' => (int) $week,
        ]]);
    }

    public function giveScore(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'student_id' => 'required|integer',
            'points' => 'required|integer|not_in:0',
            'reason' => 'required|string|max:200',
        ]);

        $student = Student::whereIn('class_id', $classIds)
            ->with('pet')
            ->findOrFail($request->input('student_id'));

        $amount = (int) $request->input('points');
        $score = $this->scoreService->giveScore(
            $student,
            $amount,
            $request->input('reason'),
            $teacher->id,
        );

        if ($amount > 0 && $student->pet) {
            $this->leaderboardService->updateTotalScore($student->class_id, $student->id, $student->total_score);
        }

        return response()->json([
            'message' => ($amount > 0 ? '加分' : '减分') . '成功',
            'data' => [
                'student_id' => $student->id,
                'student_name' => $student->name,
                'points' => $amount,
                'new_score' => $student->fresh()->total_score,
            ],
        ]);
    }

    public function batchGiveScore(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'student_ids' => 'required|array|min:1|max:100',
            'student_ids.*' => 'integer',
            'points' => 'required|integer|not_in:0',
            'reason' => 'required|string|max:200',
        ]);

        $amount = (int) $request->input('points');
        $results = $this->scoreService->batchGiveScore(
            $request->input('student_ids'),
            $amount,
            $request->input('reason'),
            $teacher->id,
        );

        return response()->json([
            'message' => '批量操作完成，处理了 ' . count($results) . ' 名学生',
            'data' => ['count' => count($results)],
        ]);
    }

    public function giveScoreByRule(Request $request, int $ruleId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $rule = ScoreRule::findOrFail($ruleId);

        $request->validate([
            'student_id' => 'required|integer',
        ]);

        $student = Student::whereIn('class_id', $classIds)->findOrFail($request->input('student_id'));

        $this->scoreService->giveScoreByRule($student, $rule, $teacher->id);

        return response()->json([
            'message' => "已按规则「{$rule->name}」处理",
            'data' => [
                'rule' => $rule->name,
                'points' => $rule->amount,
                'student_name' => $student->name,
            ],
        ]);
    }

    public function scoreHistory(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $history = $this->scoreService->getScoreHistory($student, 20);

        return response()->json([
            'data' => $history->items(),
            'meta' => [
                'current_page' => $history->currentPage(),
                'last_page' => $history->lastPage(),
                'total' => $history->total(),
            ],
            'student' => [
                'id' => $student->id,
                'name' => $student->name,
                'total_score' => $student->total_score,
            ],
        ]);
    }

    // ============================================================
    // Score Rules
    // ============================================================

    public function listScoreRules(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $rules = ScoreRule::whereIn('class_id', $classIds)
            ->orWhereNull('class_id')
            ->orderBy('sort_order')
            ->get();

        return response()->json(['data' => $rules]);
    }

    public function createScoreRule(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'name' => 'required|string|max:50',
            'amount' => 'required|integer|not_in:0',
            'category' => 'nullable|string|max:50',
            'is_positive' => 'nullable|boolean',
            'class_id' => 'nullable|integer|in:' . $classIds->join(','),
        ]);

        $rule = ScoreRule::create([
            'class_id' => $request->input('class_id', $classIds->first()),
            'name' => $request->input('name'),
            'amount' => (int) $request->input('amount'),
            'category' => $request->input('category', 'custom'),
            'is_positive' => $request->boolean('is_positive', true),
            'is_active' => true,
            'sort_order' => 0,
        ]);

        return response()->json(['message' => '规则创建成功', 'data' => $rule], 201);
    }

    public function updateScoreRule(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $rule = ScoreRule::whereIn('class_id', $classIds)->orWhereNull('class_id')->findOrFail($id);

        $rule->update($request->only(['name', 'amount', 'category', 'is_positive', 'is_active', 'sort_order']));

        return response()->json(['message' => '规则更新成功', 'data' => $rule]);
    }

    public function deleteScoreRule(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $rule = ScoreRule::whereIn('class_id', $classIds)->orWhereNull('class_id')->findOrFail($id);
        $rule->delete();

        return response()->json(['message' => '规则已删除']);
    }

    // ============================================================
    // Pets
    // ============================================================

    public function classPetsOverview(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $pets = \App\Models\Pet::whereIn('class_id', $classIds)
            ->with('student:id,name')
            ->get()
            /** @phpstan-ignore-next-line argument.unresolvableType */
            ->map(fn (\App\Models\Pet $p) => [
                'id' => $p->id,
                'student_id' => $p->student_id,
                'student_name' => $p->student?->name,
                'name' => $p->name,
                'species' => $p->species,
                'level' => $p->level,
                'exp' => $p->exp,
                'mood' => $p->mood,
                'stage_name' => $p->currentStage()['name'],
            ]);

        return response()->json(['data' => $pets]);
    }

    public function getPet(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['message' => '该学生还没有宠物'], 404);
        }

        $stage = $pet->currentStage();

        return response()->json(['data' => [
            'id' => $pet->id,
            'name' => $pet->name,
            'species' => $pet->species,
            'level' => $pet->level,
            'exp' => $pet->exp,
            'mood' => $pet->mood,
            'emoji' => $stage['emoji'],
            'stage_name' => $stage['name'],
            'last_fed_at' => $pet->last_fed_at?->toDateTimeString(),
        ]]);
    }

    public function feedPet(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['message' => '该学生还没有宠物'], 404);
        }

        $pet->feed();
        $this->leaderboardService->updatePetLevel($student->class_id, $student->id, $pet->level);

        return response()->json(['message' => "已喂养「{$pet->name}」", 'data' => [
            'mood' => $pet->mood,
            'level' => $pet->level,
        ]]);
    }

    public function renamePet(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $request->validate(['name' => 'required|string|max:20']);
        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['message' => '该学生还没有宠物'], 404);
        }

        $pet->update(['name' => $request->input('name')]);

        return response()->json(['message' => "宠物已更名为「{$pet->name}」"]);
    }

    // ============================================================
    // Leaderboard
    // ============================================================

    public function totalLeaderboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        if ($classIds->isEmpty()) {
            return response()->json(['data' => []]);
        }

        $classId = $classIds->first();
        $limit = (int) $request->input('limit', 20);
        $data = $this->leaderboardService->getClassTotalLeaderboard($classId, $limit);

        return response()->json(['data' => $data]);
    }

    public function weeklyLeaderboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        if ($classIds->isEmpty()) {
            return response()->json(['data' => []]);
        }

        $classId = $classIds->first();
        $limit = (int) $request->input('limit', 20);
        $data = $this->leaderboardService->getClassWeeklyLeaderboard($classId, $limit);

        return response()->json(['data' => $data]);
    }

    public function petLevelLeaderboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        if ($classIds->isEmpty()) {
            return response()->json(['data' => []]);
        }

        $classId = $classIds->first();
        $limit = (int) $request->input('limit', 20);
        $data = $this->leaderboardService->getPetLevelLeaderboard($classId, $limit);

        return response()->json(['data' => $data]);
    }

    // ============================================================
    // Shop
    // ============================================================

    public function listShopItems(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $items = \App\Models\ShopItem::whereIn('class_id', $classIds)
            ->orderBy('created_at', 'desc')
            ->get();

        return response()->json(['data' => $items]);
    }

    public function createShopItem(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'name' => 'required|string|max:100',
            'cost' => 'required|integer|min:1',
            'stock' => 'nullable|integer|min:0',
        ]);

        $item = \App\Models\ShopItem::create([
            'class_id' => $classIds->first(),
            'name' => $request->input('name'),
            'cost' => (int) $request->input('cost'),
            'stock' => (int) $request->input('stock', 0),
            'image' => $request->input('image'),
        ]);

        return response()->json(['message' => '商品已添加', 'data' => $item], 201);
    }

    public function updateShopItem(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $item = \App\Models\ShopItem::whereIn('class_id', $classIds)->findOrFail($id);

        $item->update($request->only(['name', 'cost', 'stock', 'image']));

        return response()->json(['message' => '商品已更新', 'data' => $item]);
    }

    public function deleteShopItem(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $item = \App\Models\ShopItem::whereIn('class_id', $classIds)->findOrFail($id);
        $item->delete();

        return response()->json(['message' => '商品已删除']);
    }

    public function listRedemptions(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $redemptions = \App\Models\ShopRedemption::whereIn('class_id', $classIds)
            ->with(['student:id,name', 'shopItem:id,name,cost'])
            ->orderBy('created_at', 'desc')
            ->paginate(20);

        return response()->json([
            'data' => $redemptions->items(),
            'meta' => [
                'current_page' => $redemptions->currentPage(),
                'last_page' => $redemptions->lastPage(),
                'total' => $redemptions->total(),
            ],
        ]);
    }

    public function approveRedemption(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $redemption = \App\Models\ShopRedemption::with(['student', 'shopItem'])
            ->whereIn('class_id', $classIds)->findOrFail($id);

        if ($redemption->status !== 'pending') {
            return response()->json(['message' => '该兑换已处理'], 400);
        }

        $student = $redemption->student;
        $itemName = $redemption->shopItem?->name ?? '未知物品';
        $cost = $redemption->cost;

        try {
            // 扣积分 + 扣宠物经验（1:1 比例）
            $this->scoreService->spendScore(
                $student,
                $cost,
                '兑换：' . $itemName,
                $teacher->id,
            );

            $redemption->update([
                'status' => 'approved',
                'approved_by' => $teacher->id,
                'approved_at' => now(),
            ]);

            return response()->json([
                'message' => '已批准兑换，扣除 ' . $cost . ' 积分',
                'data' => [
                    'remaining_score' => $student->fresh()->total_score,
                    'pet_level' => $student->pet?->fresh()->level,
                ],
            ]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 400);
        }
    }

    public function rejectRedemption(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $redemption = \App\Models\ShopRedemption::whereIn('class_id', $classIds)->findOrFail($id);
        $redemption->update(['status' => 'rejected']);

        return response()->json(['message' => '已拒绝兑换']);
    }

    public function deliverRedemption(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $redemption = \App\Models\ShopRedemption::whereIn('class_id', $classIds)->findOrFail($id);
        $redemption->update(['status' => 'delivered']);

        return response()->json(['message' => '已标记为已发放']);
    }

    // ============================================================
    // Notices
    // ============================================================

    public function listNotices(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $notices = \App\Models\Notice::whereIn('class_id', $classIds)
            ->orderBy('created_at', 'desc')
            ->paginate(20);

        return response()->json([
            'data' => $notices->items(),
            'meta' => [
                'current_page' => $notices->currentPage(),
                'last_page' => $notices->lastPage(),
                'total' => $notices->total(),
            ],
        ]);
    }

    public function createNotice(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'title' => 'required|string|max:200',
            'content' => 'required|string',
            'type' => 'nullable|string|in:info,homework,event,urgent',
        ]);

        $notice = \App\Models\Notice::create([
            'class_id' => $classIds->first(),
            'title' => $request->input('title'),
            'content' => $request->input('content'),
            'type' => $request->input('type', 'info'),
            'publisher_id' => $teacher->id,
            'is_published' => false,
        ]);

        return response()->json(['message' => '通知已创建', 'data' => $notice], 201);
    }

    public function updateNotice(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $notice = \App\Models\Notice::whereIn('class_id', $classIds)->findOrFail($id);

        $notice->update($request->only(['title', 'content', 'type']));

        return response()->json(['message' => '通知已更新', 'data' => $notice]);
    }

    public function publishNotice(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $notice = \App\Models\Notice::whereIn('class_id', $classIds)->findOrFail($id);

        $notice->update(['is_published' => true]);

        event(new \App\Events\NoticePublished(
            $notice->class_id,
            $notice->id,
            $notice->title,
            $notice->type,
        ));

        return response()->json(['message' => '通知已发布']);
    }

    public function deleteNotice(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $notice = \App\Models\Notice::whereIn('class_id', $classIds)->findOrFail($id);
        $notice->delete();

        return response()->json(['message' => '通知已删除']);
    }

    // ============================================================
    // Reports
    // ============================================================

    public function scoreTrend(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $weeks = 4;

        $trend = [];
        for ($i = $weeks - 1; $i >= 0; $i--) {
            $start = now()->subWeeks($i)->startOfWeek();
            $end = now()->subWeeks($i)->endOfWeek();
            $total = Score::whereIn('class_id', $classIds)
                ->whereBetween('created_at', [$start, $end])
                ->sum('amount');
            $trend[] = [
                'week' => $start->format('m/d'),
                'label' => "第{$start->weekOfYear}周",
                'total' => (int) $total,
            ];
        }

        return response()->json(['data' => $trend]);
    }

    public function petDistribution(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $pets = \App\Models\Pet::whereIn('class_id', $classIds)->get();
        $distribution = $pets->groupBy('level')
            ->map(fn ($group, $level) => [
                'level' => (int) $level,
                'count' => $group->count(),
                'stage_name' => $group->first()->currentStage()['name'],
            ])
            ->sortKeys()
            ->values();

        return response()->json(['data' => $distribution]);
    }

    public function studentProgress(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $studentId = $request->input('student_id');

        if (!$studentId) {
            return response()->json(['message' => '请指定学生'], 422);
        }

        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);
        $history = Score::where('student_id', $student->id)
            ->orderBy('created_at', 'desc')
            ->take(50)
            ->get();

        return response()->json(['data' => [
            'student' => ['id' => $student->id, 'name' => $student->name, 'total_score' => $student->total_score],
            'history' => $history,
        ]]);
    }

    public function exportReport(Request $request, string $type): JsonResponse
    {
        return response()->json(['message' => "导出类型 {$type} 暂不支持，请使用前端导出功能"]);
    }

    // ============================================================
    // Remaining stubs (broadcasts, attendance, homework, quizzes, grades, AI)
    // ============================================================

    public function listBroadcasts(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $broadcasts = \App\Models\Broadcast::whereIn('class_id', $classIds)
            ->orderBy('created_at', 'desc')->take(20)->get();

        return response()->json(['data' => $broadcasts]);
    }

    public function sendBroadcast(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $request->validate([
            'content' => 'required|string|max:500',
            'type' => 'nullable|string|in:banner,popup,fullscreen',
            'voice' => 'nullable|boolean',
            'loop' => 'nullable|boolean',
            'duration' => 'nullable|integer|min:0|max:300',
        ]);
        $broadcast = \App\Models\Broadcast::create([
            'class_id' => $classIds->first(),
            'content' => $request->input('content'),
            'type' => $request->input('type', 'banner'),
            'voice' => $request->boolean('voice', true),
            'loop' => $request->boolean('loop', false),
            'display_seconds' => (int) $request->input('duration', 10),
        ]);

        return response()->json(['message' => '广播已发送', 'data' => $broadcast]);
    }

    public function getBroadcast(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $broadcast = \App\Models\Broadcast::whereIn('class_id', $classIds)->findOrFail($id);

        return response()->json(['data' => $broadcast]);
    }

    public function getTodayAttendance(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $today = \App\Models\Attendance::whereIn('class_id', $classIds)
            ->whereDate('date', today())->with('student:id,name,student_no')->get();

        return response()->json(['data' => $today]);
    }

    public function startAttendance(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $students = Student::whereIn('class_id', $classIds)->where('status', 'active')->get();
        $count = 0;
        foreach ($students as $s) {
            \App\Models\Attendance::firstOrCreate(
                ['class_id' => $s->class_id, 'student_id' => $s->id, 'date' => today()],
                ['status' => 'absent']
            );
            $count++;
        }

        return response()->json(['message' => "已为 {$count} 名学生创建签到记录"]);
    }

    public function setAttendance(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $request->validate(['status' => 'required|string|in:present,late,leave,absent']);
        $record = \App\Models\Attendance::whereIn('class_id', $classIds)
            ->where('student_id', $studentId)->whereDate('date', today())->firstOrFail();
        $record->update([
            'status' => $request->input('status'),
            'check_in_time' => $request->input('status') === 'present' ? now() : null,
        ]);

        return response()->json(['message' => '考勤状态已更新']);
    }

    public function attendanceSummary(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $records = \App\Models\Attendance::whereIn('class_id', $classIds)
            ->whereDate('date', today())->get();
        $present = $records->where('status', 'present')->count();
        $late = $records->where('status', 'late')->count();
        $leave = $records->where('status', 'leave')->count();
        $absent = $records->where('status', 'absent')->count();
        $total = max($records->count(), 1);

        return response()->json(['data' => [
            'present' => $present, 'late' => $late, 'leave' => $leave, 'absent' => $absent,
            'rate' => round($present / $total * 100, 1),
        ]]);
    }

    // Simple stubs for remaining endpoints
    public function listHomework(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function createHomework(Request $request): JsonResponse
    {
        return response()->json(['message' => '作业已创建'], 201);
    }

    public function getHomework(Request $request, int $id): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function closeHomework(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '作业已关闭']);
    }

    public function getHomeworkSubmissions(Request $request, int $id): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function getHomeworkQrCode(Request $request, int $id): JsonResponse
    {
        return response()->json(['data' => ['qr_url' => '']]);
    }

    public function listQuizzes(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function createQuiz(Request $request): JsonResponse
    {
        return response()->json(['message' => '测验已创建'], 201);
    }

    public function startQuiz(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '测验已开始']);
    }

    public function stopQuiz(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '测验已结束']);
    }

    public function getQuizStats(Request $request, int $id): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function listQuestionBanks(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function createQuestionBank(Request $request): JsonResponse
    {
        return response()->json(['message' => '题库已创建'], 201);
    }

    public function addQuestion(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '题目已添加']);
    }

    public function getQuestions(Request $request, int $id): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function listGrades(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function inputGrades(Request $request): JsonResponse
    {
        return response()->json(['message' => '成绩已录入']);
    }

    public function getGradeStats(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function getGradeDistribution(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function aiChat(Request $request): JsonResponse
    {
        $request->validate(['message' => 'required|string|max:2000']);

        return response()->json(['data' => ['reply' => 'AI 助教功能需要配置 AI_PROVIDER 和 AI_API_KEY 环境变量。']]);
    }

    public function getAiCommands(Request $request): JsonResponse
    {
        return response()->json(['data' => [
            ['label' => '生成班级反馈', 'prompt' => '请根据本周课堂情况生成一段班级反馈'],
            ['label' => '重点关注学生', 'prompt' => '请分析班上需要重点关注的学生，并给出建议'],
            ['label' => '家校沟通建议', 'prompt' => '请生成本周家校沟通建议'],
            ['label' => '自动出题', 'prompt' => '请根据第三单元知识点出10道选择题'],
        ]]);
    }

    public function getAiUsage(Request $request): JsonResponse
    {
        return response()->json(['data' => ['count' => 0, 'tokens' => 0]]);
    }
}
