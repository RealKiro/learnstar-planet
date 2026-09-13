<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Pet;
use App\Models\Score;
use App\Models\ScoreRule;
use App\Models\Student;
use App\Services\AiAssistantService;
use App\Services\AttendanceService;
use App\Services\BroadcastService;
use App\Services\ClassroomMessagingService;
use App\Services\CurrencyService;
use App\Services\DashboardService;
use App\Services\LeaderboardService;
use App\Services\NoticeService;
use App\Services\PetSeriesService;
use App\Services\PetService;
use App\Services\PkService;
use App\Services\ReportService;
use App\Services\ScoreRuleService;
use App\Services\ScoreService;
use App\Services\ShopService;
use App\Services\StudentService;
use App\Services\TeacherClassScope;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class TeacherController extends Controller
{
    public function __construct(
        private readonly ScoreService $scoreService,
        private readonly LeaderboardService $leaderboardService,
        private readonly TeacherClassScope $classScope,
        private readonly NoticeService $noticeService,
        private readonly BroadcastService $broadcastService,
        private readonly AttendanceService $attendanceService,
        private readonly ShopService $shopService,
        private readonly ReportService $reportService,
        private readonly PkService $pkService,
        private readonly PetService $petService,
        private readonly PetSeriesService $petSeriesService,
        private readonly ScoreRuleService $scoreRuleService,
        private readonly AiAssistantService $aiAssistantService,
        private readonly ClassroomMessagingService $messagingService,
        private readonly DashboardService $dashboardService,
        private readonly StudentService $studentService,
    ) {
    }

    // ============================================================
    // 我的班级
    // ============================================================

    public function myClasses(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => $this->dashboardService->classesFor($teacher)]);
    }

    public function switchClass(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = (int) $request->input('class_id');

        if (!$this->dashboardService->switchTo($teacher, $classId)) {
            return response()->json(['message' => '您未被分配到此班级'], 403);
        }

        return response()->json(['message' => '已切换', 'data' => ['active_class_id' => $classId]]);
    }

    // ============================================================
    // Mode Management: classroom_display | teacher_manage
    // ============================================================

    public function getMode(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => $this->messagingService->getMode($teacher)]);
    }

    public function setMode(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $request->validate([
            'mode' => 'required|string|in:classroom_display,teacher_manage',
            'class_id' => 'nullable|integer',
            'password' => 'required_unless:mode,classroom_display|string|nullable',
        ]);

        $result = $this->messagingService->setMode(
            $teacher,
            (string) $request->input('mode'),
            $request->input('class_id'),
        );

        return response()->json($result);
    }

    // ============================================================
    // Classroom Display (大屏模式)
    // ============================================================

    public function classroomDisplay(Request $request): JsonResponse
    {
        $teacher = $request->user();

        try {
            $data = $this->messagingService->display(
                $teacher,
                $request->input('class_id', $teacher->getSetting('active_class_id')),
            );
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], $e->getCode());
        }

        return response()->json(['data' => $data]);
    }

    // ============================================================
    // Unified Classroom Messaging (merged broadcast + notice)
    // ============================================================

    /**
     * Send a classroom message (broadcast or notice)
     * type=banner|popup|fullscreen goes to broadcast table
     * type=info|homework|event|urgent goes to notice table
     */
    public function sendClassroomMessage(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $classId = (int) $request->input('class_id', $classIds[0] ?? 0);

        // 权限预检先于参数校验（保持历史响应顺序：越权 403 优先于 422）
        if (!in_array($classId, $classIds)) {
            return response()->json(['message' => '无权限'], 403);
        }

        $request->validate([
            'class_id' => 'required|integer',
            'content' => 'required|string|max:500',
            'title' => 'nullable|string|max:200',
            'type' => 'required|string|in:banner,popup,fullscreen,info,homework,event,urgent',
            'display_seconds' => 'nullable|integer|min:3|max:300',
            'voice' => 'nullable|boolean',
        ]);

        try {
            $result = $this->messagingService->send(
                $teacher,
                $classIds,
                $classId,
                (string) $request->input('type'),
                [
                    'content' => (string) $request->input('content'),
                    'title' => $request->input('title'),
                    'display_seconds' => $request->input('display_seconds', 10),
                    'voice' => $request->boolean('voice', true),
                ],
            );
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], $e->getCode());
        }

        return response()->json($result);
    }

    /**
     * Poll classroom messages for display mode
     */
    public function pollClassroomMessages(Request $request): JsonResponse
    {
        $teacher = $request->user();

        try {
            $data = $this->messagingService->poll(
                $teacher,
                $request->input('class_id', $teacher->getSetting('active_class_id')),
                $request->input('since'),
            );
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], $e->getCode());
        }

        return response()->json(['data' => $data]);
    }

    // ============================================================
    // Dashboard
    // ============================================================

    public function dashboard(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => $this->dashboardService->forTeacher($teacher, $this->teacherClassIds($teacher))]);
    }

    // ============================================================
    // Student Management
    // ============================================================

    public function listStudents(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $students = $this->studentService->list($classIds, $request->input('search'));

        return response()->json([
            // 附加宠物字段（保留原始字段，兼容所有调用方）
            'data' => collect($students->items())->map(fn ($s) => array_merge($s->toArray(), [
                'pet_species' => $s->pet->species ?? '',
                'pet_level' => $s->pet->level ?? 0,
                'pet_name' => $s->pet->name ?? '',
            ])),
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
        $classIds = $this->teacherClassIds($teacher);

        $result = $this->studentService->import($classIds, (array) $request->input('students', []));

        return response()->json([
            'message' => $result['message'],
            'data' => ['imported_count' => $result['imported_count'], 'skipped' => $result['skipped']],
        ]);
    }

    public function updateStudent(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $student = $this->studentService->update($classIds, $id, $request->only(['name', 'gender', 'student_no']));

        return response()->json(['message' => '更新成功', 'data' => $student]);
    }

    public function createStudent(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $validator = Validator::make($request->all(), [
            'name' => 'required|string|max:50',
            'class_id' => 'required|integer|exists:class_rooms,id',
            'gender' => 'nullable|string',
            'student_no' => 'nullable|string|max:50',
        ]);
        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $student = $this->studentService->create($classIds, $request->only(['name', 'class_id', 'gender', 'student_no']));
        if (!$student) {
            return response()->json(['message' => '只能在自己管理的班级添加学生'], 403);
        }

        return response()->json([
            'message' => '学生「' . $student->name . '」已添加',
            'data' => $student,
        ], 201);
    }

    public function deleteStudent(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $student = $this->studentService->delete($classIds, $id);

        return response()->json(['message' => '学生「' . $student->name . '」已删除']);
    }

    // ============================================================
    // Score Management
    // ============================================================

    public function scoreSummary(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        return response()->json(['data' => $this->scoreService->summaryFor($classIds)]);
    }

    public function giveScore(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $request->validate([
            'student_id' => 'required|integer',
            'points' => 'required|integer|not_in:0',
            'reason' => 'required|string|max:200',
        ]);

        $student = Student::whereIn('class_id', $classIds)
            ->with('pet')
            ->findOrFail($request->input('student_id'));

        if (!$student->pet) {
            // 自动造宠：随机神话系列物种（species 体系；旧 type 体系已废弃，level 从 1 起）
            $pool = Pet::speciesPoolForSeries('myth');
            $species = $pool[array_rand($pool)] ?? 'zhulong';
            $student->pet()->create([
                'class_id' => $student->class_id,
                'species' => $species,
                'level' => 1,
                'experience' => 0,
                'mood' => 80,
            ]);
            $student->refresh()->load('pet');
        }

        $amount = (int) $request->input('points');
        $score = $this->scoreService->giveScore(
            $student,
            $amount,
            $request->input('reason'),
            $teacher->id,
        );

        if ($student->pet) {
            if ($amount > 0) {
                $student->pet->addExperience(abs($amount));
            } else {
                $student->pet->removeExperience(abs($amount));
            }
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
        $classIds = $this->teacherClassIds($teacher);

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
        $classIds = $this->teacherClassIds($teacher);
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
        $classIds = $this->teacherClassIds($teacher);
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

    /**
     * 全班最近积分记录（课堂评价右侧"最近记录"）
     */
    public function recentScores(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        return response()->json(['data' => $this->scoreService->recentFor($classIds)]);
    }

    // ============================================================
    // Score Rules

    /**
     * 撤回一条积分记录
     */
    public function undoScore(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $score = Score::whereIn('class_id', $classIds)->findOrFail($id);
        $undo = app(\App\Services\ScoreService::class)->undoScore($score, $teacher->id);

        return response()->json([
            'message' => '已撤回',
            'data' => $undo,
        ]);
    }

    // ============================================================

    /**
     * 获取教师关联的所有班级 ID（含班主任 + 科任/副班等）
     */
    private function teacherClassIds(\App\Models\User $teacher): \Illuminate\Support\Collection
    {
        return $this->classScope->ids($teacher);
    }

    public function listScoreRules(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->scoreRuleService->listForTeacher($request->user())]);
    }

    public function createScoreRule(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $request->validate([
            'name' => 'required|string|max:50',
            'amount' => 'required|integer|not_in:0',
            'category' => 'nullable|string|max:50',
            'is_positive' => 'nullable|boolean',
            'class_id' => 'nullable|integer|in:' . $classIds->join(','),
        ]);

        $rule = $this->scoreRuleService->createForTeacher($teacher, $request->all([
            'name', 'amount', 'category', 'is_positive', 'class_id',
        ]));

        return response()->json(['message' => '规则创建成功', 'data' => $rule], 201);
    }

    public function updateScoreRule(Request $request, int $id): JsonResponse
    {
        $rule = $this->scoreRuleService->updateForTeacher($request->user(), $id, $request->only([
            'name', 'amount', 'category', 'is_positive', 'is_active', 'sort_order',
        ]));

        return response()->json(['message' => '规则更新成功', 'data' => $rule]);
    }

    public function deleteScoreRule(Request $request, int $id): JsonResponse
    {
        $this->scoreRuleService->deleteForTeacher($request->user(), $id);

        return response()->json(['message' => '规则已删除']);
    }

    // ============================================================
    // Pets
    // ============================================================

    public function classPetsOverview(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->petService->classPetsOverview($request->user())]);
    }

    public function getPet(Request $request, int $studentId): JsonResponse
    {
        try {
            return response()->json(['data' => $this->petService->petFor($request->user(), $studentId)]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 404);
        }
    }

    public function feedPet(Request $request, int $studentId): JsonResponse
    {
        try {
            $result = $this->petService->feed($request->user(), $studentId);

            return response()->json(['message' => $result['message'], 'data' => [
                'mood' => $result['mood'],
                'level' => $result['level'],
            ]]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 404);
        }
    }

    public function renamePet(Request $request, int $studentId): JsonResponse
    {
        $request->validate(['name' => 'required|string|max:20']);

        try {
            return response()->json(['message' => $this->petService->rename($request->user(), $studentId, (string) $request->input('name'))]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 404);
        }
    }

    // ============================================================
    // Leaderboard
    // ============================================================

    public function totalLeaderboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

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
        $classIds = $this->teacherClassIds($teacher);

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
        $classIds = $this->teacherClassIds($teacher);

        if ($classIds->isEmpty()) {
            return response()->json(['data' => []]);
        }

        $classId = $classIds->first();
        $limit = (int) $request->input('limit', 20);
        $data = $this->leaderboardService->getPetLevelLeaderboard($classId, $limit);

        return response()->json(['data' => $data]);
    }

    // ============================================================
    // 年级战场 (PK)
    // ============================================================

    /**
     * 获取同年级各班 PK 排行榜
     */
    public function pkLeaderboard(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->pkService->leaderboard($request->user())]);
    }

    /**
     * 获取本班 PK 统计数据
     */
    public function myPkStats(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->pkService->myStats($request->user())]);
    }

    /**
     * 发起 PK 挑战（记录挑战事件）
     */
    public function challengePk(Request $request): JsonResponse
    {
        $classIds = $this->teacherClassIds($request->user());

        $request->validate([
            'target_class_id' => 'required|integer',
        ]);

        $targetClassId = (int) $request->input('target_class_id');
        $myClassId = $classIds->first();

        if (!$myClassId || $myClassId === $targetClassId) {
            return response()->json(['message' => '无效的挑战目标'], 400);
        }

        $targetClass = \App\Models\ClassRoom::find($targetClassId);
        if (!$targetClass) {
            return response()->json(['message' => '目标班级不存在'], 404);
        }

        $data = $this->pkService->challenge($myClassId, $targetClass);

        return response()->json([
            'message' => '🚀 挑战已发起！',
            'data' => $data,
        ]);
    }

    /**
     * 获取当前班级信息（含总积分、设置、宠物系列等）
     */
    public function classInfo(Request $request): JsonResponse
    {
        try {
            return response()->json(['data' => $this->petSeriesService->classInfo($request->user())]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 400);
        }
    }

    /**
     * 切换班级宠物系列
     */
    public function switchSeries(Request $request): JsonResponse
    {
        $request->validate([
            'series_id' => 'required|string|max:50',
        ]);

        $seriesId = $request->input('series_id');
        $validSeries = ['myth', 'pokemon', 'national', 'digimon', 'magic', 'prehistoric', 'constellation', 'festival', 'qixia', 'dongfang'];

        if (!in_array($seriesId, $validSeries, true)) {
            return response()->json(['message' => '无效的系列ID，可选值：' . implode(', ', $validSeries)], 422);
        }

        try {
            $data = $this->petSeriesService->switchSeries($request->user(), $seriesId);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 400);
        }

        return response()->json([
            'message' => "已切换系列为「{$seriesId}」，全班 {$data['student_count']} 人各获一次免费自选该系列宠物的机会",
            'data' => [
                'series_id' => $data['series_id'],
                'class_id' => $data['class_id'],
                'free_pick_granted' => $data['free_pick_granted'],
                'granted_students' => $data['granted_students'],
            ],
        ]);
    }

    // ============================================================
    // Pet Switching & Auto-Assignment
    // ============================================================

    public function switchPet(Request $request, int $studentId): JsonResponse
    {
        $request->validate([
            'pet_species' => 'required|string|max:50',
            'pet_name' => 'nullable|string|max:20',
        ]);

        try {
            $result = $this->petService->switchPet(
                $request->user(),
                $studentId,
                (string) $request->input('pet_species'),
                (string) $request->input('pet_name', $request->input('pet_species')),
            );
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], $e->getCode() ?: 400);
        }

        return response()->json(['message' => $result['message'], 'data' => $result['data']]);
    }

    /**
     * 学生宠物图鉴收藏
     */
    public function petCollection(Request $request, int $studentId): JsonResponse
    {
        return response()->json(['data' => $this->petService->collection($request->user(), $studentId)]);
    }

    // ============================================================
    // Shop
    // ============================================================

    public function listShopItems(Request $request): JsonResponse
    {
        $items = $this->shopService->itemsFor($request->user(), $request->input('currency_type'));

        return response()->json(['data' => $items]);
    }

    public function createShopItem(Request $request): JsonResponse
    {
        $request->validate([
            'name' => 'required|string|max:100',
            'description' => 'nullable|string',
            'category' => 'nullable|string|max:50',
            'cost_score' => 'required|integer|min:1',
            'currency_type' => 'nullable|string|max:50',
            'event_tag' => 'nullable|string|max:50',
            'stock' => 'nullable|integer|min:0',
            'image_path' => 'nullable|string|max:255',
        ]);

        $item = $this->shopService->createItem($request->user(), $request->all([
            'name', 'description', 'category', 'cost_score',
            'currency_type', 'event_tag', 'stock', 'image_path',
        ]));

        return response()->json(['message' => '商品已添加', 'data' => $item], 201);
    }

    public function updateShopItem(Request $request, int $id): JsonResponse
    {
        $item = $this->shopService->findItem($request->user(), $id);
        $item = $this->shopService->updateItem($item, $request->only([
            'name',
            'description',
            'category',
            'cost_score',
            'currency_type',
            'event_tag',
            'stock',
            'image_path',
            'is_active',
        ]));

        return response()->json(['message' => '商品已更新', 'data' => $item]);
    }

    public function deleteShopItem(Request $request, int $id): JsonResponse
    {
        $item = $this->shopService->findItem($request->user(), $id);
        $this->shopService->deleteItem($item);

        return response()->json(['message' => '商品已删除']);
    }

    public function listRedemptions(Request $request): JsonResponse
    {
        $redemptions = $this->shopService->paginateRedemptions($request->user());

        return response()->json([
            'data' => $redemptions->items(),
            'meta' => [
                'current_page' => $redemptions->currentPage(),
                'last_page' => $redemptions->lastPage(),
                'total' => $redemptions->total(),
            ],
        ]);
    }

    /**
     * 创建兑换记录（教师代学生发起）
     */
    public function createRedemption(Request $request): JsonResponse
    {
        $request->validate([
            'student_id' => 'required|integer',
            'shop_item_id' => 'required|integer',
        ]);

        $item = $this->shopService->findItem($request->user(), (int) $request->input('shop_item_id'));

        if (!$item->is_active) {
            return response()->json(['message' => '该商品已下架'], 422);
        }

        $redemption = $this->shopService->createRedemption(
            $request->user(),
            (int) $request->input('student_id'),
            $item
        );

        return response()->json(['message' => '兑换请求已创建', 'data' => $redemption], 201);
    }

    public function approveRedemption(Request $request, int $id): JsonResponse
    {
        try {
            $result = $this->shopService->approveRedemption($request->user(), $id);

            return response()->json([
                'message' => $result['message'],
                'data' => [
                    'remaining_score' => $result['remaining_score'],
                    'pet_level' => $result['pet_level'],
                ],
            ]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 400);
        }
    }

    public function rejectRedemption(Request $request, int $id): JsonResponse
    {
        $this->shopService->rejectRedemption($request->user(), $id);

        return response()->json(['message' => '已拒绝兑换']);
    }

    public function deliverRedemption(Request $request, int $id): JsonResponse
    {
        $this->shopService->deliverRedemption($request->user(), $id);

        return response()->json(['message' => '已标记为已发放']);
    }

    // ============================================================
    // Notices
    // ============================================================

    public function listNotices(Request $request): JsonResponse
    {
        $notices = $this->noticeService->paginate($request->user());

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
        $request->validate([
            'title' => 'required|string|max:200',
            'content' => 'required|string',
            'type' => 'nullable|string|in:info,homework,event,urgent',
        ]);

        $notice = $this->noticeService->create($request->user(), [
            'title' => $request->input('title'),
            'content' => $request->input('content'),
            'type' => $request->input('type', 'info'),
        ]);

        return response()->json(['message' => '通知已创建', 'data' => $notice], 201);
    }

    public function updateNotice(Request $request, int $id): JsonResponse
    {
        $notice = $this->noticeService->findInScope($request->user(), $id);
        $this->noticeService->update($notice, $request->only(['title', 'content', 'type']));

        return response()->json(['message' => '通知已更新', 'data' => $notice]);
    }

    public function publishNotice(Request $request, int $id): JsonResponse
    {
        $notice = $this->noticeService->findInScope($request->user(), $id);
        $this->noticeService->publish($notice);

        return response()->json(['message' => '通知已发布']);
    }

    public function unpublishNotice(Request $request, int $id): JsonResponse
    {
        $notice = $this->noticeService->findInScope($request->user(), $id);
        $this->noticeService->unpublish($notice);

        return response()->json(['message' => '通知已撤回']);
    }

    public function deleteNotice(Request $request, int $id): JsonResponse
    {
        $notice = $this->noticeService->findInScope($request->user(), $id);
        $this->noticeService->delete($notice);

        return response()->json(['message' => '通知已删除']);
    }

    // ============================================================
    // Reports
    // ============================================================

    public function scoreTrend(Request $request): JsonResponse
    {
        $days = max(1, min((int) $request->input('days', 7), 365));

        return response()->json(['data' => $this->reportService->scoreTrend($request->user(), $days)]);
    }

    public function petDistribution(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->reportService->petDistribution($request->user())]);
    }

    public function studentProgress(Request $request): JsonResponse
    {
        $studentId = $request->input('student_id') ? (int) $request->input('student_id') : null;

        return response()->json(['data' => $this->reportService->studentProgress($request->user(), $studentId)]);
    }

    public function exportReport(Request $request, string $type)
    {
        $classIds = $this->teacherClassIds($request->user());
        $classId = (int) $request->input('class_id', $classIds->first() ?? 0);

        if (!in_array($classId, $classIds->toArray())) {
            return response()->json(['message' => '无权限'], 403);
        }

        $download = $this->reportService->export($type, $classId, $request->input('date'));

        if ($download === null) {
            return response()->json(['message' => "导出类型 {$type} 不支持，可选: scores, pets, attendance"]);
        }

        return $download;
    }

    // ============================================================
    // Broadcasts
    // ============================================================

    public function listBroadcasts(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->broadcastService->recent($request->user())]);
    }

    public function sendBroadcast(Request $request): JsonResponse
    {
        $request->validate([
            'content' => 'required|string|max:500',
            'type' => 'nullable|string|in:banner,popup,fullscreen',
            'class_ids' => 'nullable|array',
            'class_ids.*' => 'integer|exists:class_rooms,id',
            'voice' => 'nullable|boolean',
            'loop' => 'nullable|boolean',
            'duration' => 'nullable|integer|min:0|max:300',
        ]);

        $sent = $this->broadcastService->send(
            $request->user(),
            $request->input('content'),
            $request->input('type', 'banner'),
            $request->boolean('voice', true),
            $request->boolean('loop', false),
            (int) $request->input('duration', 10),
            $request->input('class_ids'),
        );

        if ($sent === 0) {
            return response()->json(['message' => '没有可发送的班级'], 400);
        }

        return response()->json([
            'message' => "广播已发送至 {$sent} 个班级",
            'data' => ['sent_count' => $sent],
        ]);
    }

    public function getBroadcast(Request $request, int $id): JsonResponse
    {
        return response()->json(['data' => $this->broadcastService->findInScope($request->user(), $id)]);
    }

    public function getTodayAttendance(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->attendanceService->today($request->user())]);
    }

    public function startAttendance(Request $request): JsonResponse
    {
        $result = $this->attendanceService->start($request->user());

        $msg = "已为 {$result['total']} 名学生创建考勤记录（默认到课）";
        if ($result['wechat_leave_count'] > 0) {
            $msg .= "，其中 {$result['wechat_leave_count']} 人已通过企业微信请假";
        }

        return response()->json(['message' => $msg, 'data' => $result]);
    }

    public function setAttendance(Request $request, int $studentId): JsonResponse
    {
        $request->validate([
            'status' => 'required|string|in:present,late,leave,absent',
            'remark' => 'nullable|string|max:500',
        ]);

        $record = $this->attendanceService->setStatus(
            $request->user(),
            $studentId,
            $request->input('status'),
            $request->input('remark'),
        );

        return response()->json(['message' => '考勤状态已更新', 'data' => $record]);
    }

    public function markManualLeave(Request $request, int $studentId): JsonResponse
    {
        $request->validate(['remark' => 'required|string|max:500']);

        $record = $this->attendanceService->markLeave(
            $request->user(),
            $studentId,
            $request->input('remark'),
        );

        return response()->json(['message' => '已标记为请假', 'data' => ['id' => $record->id, 'status' => $record->status, 'source' => $record->source, 'remark' => $record->remark]]);
    }

    public function markManualAbsent(Request $request, int $studentId): JsonResponse
    {
        $request->validate(['remark' => 'nullable|string|max:500']);

        $record = $this->attendanceService->markAbsent(
            $request->user(),
            $studentId,
            $request->input('remark'),
        );

        return response()->json(['message' => '已标记为缺勤，建议联系家长确认情况', 'data' => ['id' => $record->id, 'status' => $record->status, 'source' => $record->source, 'remark' => $record->remark]]);
    }

    public function attendanceSummary(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->attendanceService->summary($request->user())]);
    }

    private function getAccessibleClassIds($teacher): array
    {
        return $this->classScope->accessibleIds($teacher);
    }

    /**
     * 教师端 AI 配置状态：管理员未配置有效的 AI API Key 时返回 enabled=false，
     * 前端据此隐藏 AI 助教入口。
     */
    public function aiConfig(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => $this->aiAssistantService->configFor($teacher)]);
    }

    public function aiChat(Request $request): JsonResponse
    {
        $request->validate([
            'message' => 'required|string|max:2000',
            // 教师指定的模型（可选，须在供应商模型白名单内，否则回退默认）
            'model' => 'nullable|string|max:100',
        ]);

        $teacher = $request->user();
        $reply = $this->aiAssistantService->chat(
            $teacher,
            (string) $request->input('message'),
            $request->input('model'),
        );

        return response()->json(['data' => ['reply' => $reply]]);
    }

    /**
     * 教师端 AI 用量（前端 AIPage 每次发送后刷新）
     */
    public function getAiUsage(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => $this->aiAssistantService->usageFor($teacher)]);
    }

    /**
     * 教师端 AI 预设命令
     */
    public function getAiCommands(Request $request): JsonResponse
    {
        return response()->json(['data' => $this->aiAssistantService->commands()]);
    }

    // ============================================================
    // 汇率管理（教师端）
    // ============================================================

    public function listExchangeRates(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => app(CurrencyService::class)->ratesForSchool($teacher->school_id)]);
    }

    public function createExchangeRate(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $request->validate([
            'name' => 'required|string|max:100',
            'from_currency' => 'required|string|in:score,science,reading,class_point',
            'to_currency' => 'required|string|in:score,science,reading,class_point',
            'rate' => 'required|numeric|min:0.01',
        ]);

        $rate = app(CurrencyService::class)->createRate($teacher->school_id, $request->all([
            'name', 'from_currency', 'to_currency', 'rate',
        ]));

        return response()->json(['message' => '汇率已添加', 'data' => $rate], 201);
    }

    public function updateExchangeRate(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $request->validate([
            'rate' => 'sometimes|numeric|min:0.01',
            'is_active' => 'sometimes|boolean',
        ]);

        $rate = app(CurrencyService::class)->updateRate($teacher->school_id, $id, $request->only(['rate', 'is_active']));

        return response()->json(['message' => '汇率已更新', 'data' => $rate]);
    }

    // ============================================================
    // 兑换中心
    // ============================================================

    public function listWallets(Request $request): JsonResponse
    {
        return response()->json(['data' => app(CurrencyService::class)->walletsFor($request->user())]);
    }

    /**
     * 兑换记录（仅教师所带班级学生）
     */
    public function exchangeLogs(Request $request): JsonResponse
    {
        return response()->json(app(CurrencyService::class)->logsFor($request->user()));
    }

    /**
     * 积分兑换
     */
    public function exchangeCurrency(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $request->validate([
            'student_id' => 'required|integer',
            'to_currency' => 'required|string',
            'amount' => 'required|integer|min:1',
        ]);

        $student = Student::whereIn('class_id', $classIds)->findOrFail($request->input('student_id'));

        try {
            $result = app(CurrencyService::class)->exchange(
                $student->id,
                $request->input('to_currency'),
                (int) $request->input('amount'),
                $teacher->id,
            );

            return response()->json(['message' => '兑换成功', 'data' => $result]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 422);
        }
    }

    /**
     * 交叉兑换（币种间兑换）
     */
    public function crossExchangeCurrency(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $request->validate([
            'student_id' => 'required|integer',
            'from_currency' => 'required|string',
            'to_currency' => 'required|string',
            'amount' => 'required|integer|min:1',
        ]);

        $student = Student::whereIn('class_id', $classIds)->findOrFail($request->input('student_id'));

        try {
            $result = app(CurrencyService::class)->crossExchange(
                $student->id,
                $request->input('from_currency'),
                $request->input('to_currency'),
                (int) $request->input('amount'),
                $teacher->id,
            );

            return response()->json(['message' => '兑换成功', 'data' => $result]);
        } catch (\DomainException $e) {
            return response()->json(['message' => $e->getMessage()], 422);
        }
    }
}
