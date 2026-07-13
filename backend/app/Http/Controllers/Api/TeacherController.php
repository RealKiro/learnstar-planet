<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\Notice;
use App\Models\Pet;
use App\Models\Score;
use App\Models\ScoreRule;
use App\Models\ShopItem;
use App\Models\ShopRedemption;
use App\Models\Student;
use App\Models\Wallet;
use App\Services\CurrencyService;
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
    // 我的班级
    // ============================================================

    public function myClasses(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $assignments = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->with('classRoom:id,name,grade')
            ->get()
            ->map(fn ($a) => [
                'class_id' => $a->class_room_id,
                'class_name' => $a->classRoom?->name,
                'grade' => $a->classRoom?->grade,
                'role' => $a->role,
            ]);

        return response()->json(['data' => $assignments]);
    }

    public function switchClass(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = (int) $request->input('class_id');

        $isAssigned = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->where('class_room_id', $classId)
            ->exists();

        if (!$isAssigned) {
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
        $mode = $teacher->getSetting('display_mode', 'teacher_manage');
        $activeClassId = $teacher->getSetting('active_class_id', null);

        return response()->json(['data' => [
            'mode' => $mode,
            'active_class_id' => $activeClassId,
        ]]);
    }

    public function setMode(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $request->validate([
            'mode' => 'required|string|in:classroom_display,teacher_manage',
            'class_id' => 'nullable|integer',
        ]);

        $mode = $request->input('mode');
        $teacher->setSetting('display_mode', $mode);

        if ($classId = $request->input('class_id')) {
            $isAssigned = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
                ->where('class_room_id', $classId)
                ->exists();
            if ($isAssigned) {
                $teacher->setSetting('active_class_id', $classId);
            }
        }

        return response()->json([
            'message' => '已切换为' . ($mode === 'classroom_display' ? '班级大屏' : '教师管理') . '模式',
            'data' => [
                'mode' => $mode,
                'active_class_id' => $teacher->getSetting('active_class_id'),
            ],
        ]);
    }

    // ============================================================
    // Classroom Display (大屏模式)
    // ============================================================

    public function classroomDisplay(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = $request->input('class_id', $teacher->getSetting('active_class_id'));

        if (!$classId) {
            return response()->json(['message' => '请先选择班级'], 400);
        }

        $isAssigned = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->where('class_room_id', $classId)
            ->exists();
        if (!$isAssigned) {
            return response()->json(['message' => '您未被分配到此班级'], 403);
        }

        // Load class room
        $classRoom = ClassRoom::with(['students' => function ($q) {
            $q->where('status', 'active')->with('pet');
        }])->findOrFail($classId);

        // Pet overview for all students
        $pets = $classRoom->students->map(function (Student $s): array {
            $pet = $s->pet;
            $stage = $pet ? $pet->currentStage() : ['emoji' => '\ud83e\udd14', 'name' => '未孵化', 'title' => ''];

            return [
                'student_id' => $s->id,
                'student_name' => $s->name,
                'total_score' => $s->total_score,
                'has_pet' => $pet !== null,
                'pet_name' => $pet?->name,
                'pet_type' => $pet?->type,
                'level' => $pet->level ?? 0,
                'experience' => $pet->experience ?? 0,
                'mood' => $pet->mood ?? 0,
                'emoji' => $stage['emoji'],
                'stage_name' => $stage['name'],
            ];
        })->values();

        // Active broadcasts for this class (not expired)
        $broadcasts = \App\Models\Broadcast::where('class_id', $classId)
            ->whereIn('status', ['pending', 'sent'])
            ->orderBy('created_at', 'desc')
            ->take(5)
            ->get()
            ->map(fn ($b) => [
                'id' => $b->id,
                'content' => $b->content,
                'type' => $b->type,
                'display_seconds' => $b->display_seconds,
                'voice_enabled' => $b->voice_enabled,
                'created_at' => $b->created_at?->diffForHumans(),
            ]);

        // Recent notices (last 7 days, published)
        $notices = Notice::where('class_id', $classId)
            ->where('is_published', true)
            ->where('published_at', '>=', now()->subDays(7))
            ->orderBy('published_at', 'desc')
            ->take(3)
            ->get()
            ->map(fn ($n) => [
                'id' => $n->id,
                'title' => $n->title,
                'content' => $n->content,
                'type' => $n->type,
                'published_at' => $n->published_at?->diffForHumans(),
            ]);

        // Recent scores feed
        $recentScores = Score::where('class_id', $classId)
            ->where('created_at', '>=', now()->subHours(4))
            ->with('student:id,name')
            ->orderBy('created_at', 'desc')
            ->take(20)
            ->get()
            ->map(fn (Score $s) => [
                'student_name' => $s->student?->name,
                'amount' => $s->amount,
                'reason' => $s->reason,
                'time' => $s->created_at?->diffForHumans(),
            ]);

        return response()->json(['data' => [
            'class_name' => $classRoom->name,
            'grade' => $classRoom->grade,
            'student_count' => $classRoom->students->count(),
            'pets' => $pets,
            'broadcasts' => $broadcasts,
            'notices' => $notices,
            'recent_scores' => $recentScores,
        ]]);
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

        $type = $request->input('type');
        $content = $request->input('content');

        // Broadcast types: banner, popup, fullscreen
        if (in_array($type, ['banner', 'popup', 'fullscreen'])) {
            $broadcast = \App\Models\Broadcast::create([
                'school_id' => $teacher->school_id,
                'class_id' => $classId,
                'teacher_id' => $teacher->id,
                'content' => $content,
                'type' => $type,
                'voice_enabled' => $request->boolean('voice', true),
                'display_seconds' => (int) $request->input('display_seconds', 10),
                'status' => 'sent',
                'sent_at' => now(),
            ]);

            return response()->json([
                'message' => '广播已发送',
                'data' => ['id' => $broadcast->id, 'type' => 'broadcast'],
            ]);
        }

        // Notice types: info, homework, event, urgent
        $notice = Notice::create([
            'class_id' => $classId,
            'school_id' => $teacher->school_id,
            'title' => $request->input('title', $type === 'urgent' ? '\u7d27\u6025\u901a\u77e5' : '\u901a\u77e5'),
            'content' => $content,
            'type' => $type,
            'published_by' => $teacher->id,
            'is_published' => true,
            'published_at' => now(),
        ]);

        return response()->json([
            'message' => '通知已发布',
            'data' => ['id' => $notice->id, 'type' => 'notice'],
        ]);
    }

    /**
     * Poll classroom messages for display mode
     */
    public function pollClassroomMessages(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = $request->input('class_id', $teacher->getSetting('active_class_id'));

        if (!$classId) {
            return response()->json(['message' => '\u8bf7\u5148\u9009\u62e9\u73ed\u7ea7'], 400);
        }

        $since = $request->input('since');
        $sinceTime = $since ? \Carbon\Carbon::parse($since) : now()->subMinutes(5);

        $broadcasts = \App\Models\Broadcast::where('class_id', $classId)
            ->where('created_at', '>=', $sinceTime)
            ->whereIn('status', ['sent'])
            ->orderBy('created_at', 'desc')
            ->take(5)
            ->get()
            ->map(fn ($b) => [
                'id' => $b->id,
                'content' => $b->content,
                'type' => $b->type,
                'display_seconds' => $b->display_seconds,
                'voice_enabled' => $b->voice_enabled,
                'created_at' => $b->created_at?->toIso8601String(),
            ]);

        $notices = Notice::where('class_id', $classId)
            ->where('is_published', true)
            ->where('published_at', '>=', $sinceTime)
            ->orderBy('published_at', 'desc')
            ->take(3)
            ->get()
            ->map(fn ($n) => [
                'id' => $n->id,
                'title' => $n->title,
                'content' => $n->content,
                'type' => $n->type,
                'published_at' => $n->published_at?->toIso8601String(),
            ]);

        return response()->json(['data' => [
            'broadcasts' => $broadcasts,
            'notices' => $notices,
            'polled_at' => now()->toIso8601String(),
        ]]);
    }

    // ============================================================
    // Dashboard
    // ============================================================

    public function dashboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->pluck('class_room_id');

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
        $classIds = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->pluck('class_room_id');

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

        if (!$student->pet) {
            $types = array_keys(Pet::petTypes());
            $randomType = $types[array_rand($types)];
            $randomName = Pet::petTypes()[$randomType];
            $student->pet()->create([
                'class_id' => $student->class_id,
                'type' => $randomType,
                'name' => $randomName,
                'level' => 0,
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

        $pets = Pet::whereIn('class_id', $classIds)
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
        $classIds = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->pluck('class_room_id');

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
        $classIds = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->pluck('class_room_id');

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
        $classIds = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
            ->pluck('class_room_id');

        if ($classIds->isEmpty()) {
            return response()->json(['data' => []]);
        }

        $classId = $classIds->first();
        $limit = (int) $request->input('limit', 20);
        $data = $this->leaderboardService->getPetLevelLeaderboard($classId, $limit);

        return response()->json(['data' => $data]);
    }

    // ============================================================
    // Pet Switching & Auto-Assignment
    // ============================================================

    public function switchPet(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $request->validate([
            'pet_type' => 'required|string|max:50',
            'pet_name' => 'nullable|string|max:20',
        ]);

        $petType = $request->input('pet_type');
        $petName = $request->input('pet_name', $petType);

        if (!in_array($petType, array_keys(Pet::petTypes()))) {
            return response()->json(['message' => '无效的宠物类型'], 422);
        }

        $pet = $student->pet;
        $switchCost = 5;

        if ($pet) {
            $pet->experience = max(0, $pet->experience - $switchCost);
            $pet->type = $petType;
            $pet->name = $petName;
            $pet->save();

            return response()->json([
                'message' => '宠物已切换为「' . $petName . '」（扣除 ' . $switchCost . ' 成长值）',
                'data' => [
                    'pet_name' => $pet->name,
                    'pet_type' => $pet->type,
                    'level' => $pet->level,
                    'experience' => $pet->experience,
                    'cost' => $switchCost,
                ],
            ]);
        }

        $pet = Pet::create([
            'student_id' => $student->id,
            'class_id' => $student->class_id,
            'name' => $petName,
            'type' => $petType,
            'level' => 0,
            'experience' => 0,
            'mood' => 80,
        ]);

        return response()->json([
            'message' => '已为您分配宠物「' . $petName . '」',
            'data' => [
                'pet_name' => $pet->name,
                'pet_type' => $pet->type,
                'level' => $pet->level,
                'experience' => $pet->experience,
            ],
        ]);
    }

    public function getPetTypes(): JsonResponse
    {
        $categories = Pet::petCategories();
        $allTypes = Pet::petTypes();
        $result = [];

        foreach ($categories as $catKey => $catName) {
            $types = [];
            foreach (Pet::petTypesBySeries($catKey) as $typeKey => $typeName) {
                $types[] = ['key' => $typeKey, 'name' => $typeName];
            }
            $result[] = ['category' => $catKey, 'category_name' => $catName, 'pets' => $types];
        }

        return response()->json(['data' => $result]);
    }

    // ============================================================
    // Shop
    // ============================================================

    public function listShopItems(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $query = ShopItem::whereIn('class_id', $classIds);

        if ($currency = $request->input('currency_type')) {
            $query->byCurrency($currency);
        }

        $items = $query->orderBy('created_at', 'desc')->get();

        return response()->json(['data' => $items]);
    }

    public function createShopItem(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

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

        $item = ShopItem::create([
            'class_id' => $classIds->first(),
            'name' => $request->input('name'),
            'description' => $request->input('description'),
            'category' => $request->input('category', 'physical'),
            'cost_score' => (int) $request->input('cost_score'),
            'currency_type' => $request->input('currency_type', 'score'),
            'event_tag' => $request->input('event_tag'),
            'stock' => (int) $request->input('stock', 0),
            'image_path' => $request->input('image_path'),
            'is_active' => true,
        ]);

        return response()->json(['message' => '商品已添加', 'data' => $item], 201);
    }

    public function updateShopItem(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $item = ShopItem::whereIn('class_id', $classIds)->findOrFail($id);

        $item->update($request->only([
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
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $item = ShopItem::whereIn('class_id', $classIds)->findOrFail($id);
        $item->delete();

        return response()->json(['message' => '商品已删除']);
    }

    public function listRedemptions(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $redemptions = ShopRedemption::whereIn('class_id', $classIds)
            ->with(['student:id,name', 'shopItem:id,name,cost_score,currency_type,event_tag,category'])
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

    /**
     * 创建兑换记录（教师代学生发起）
     */
    public function createRedemption(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'student_id' => 'required|integer',
            'shop_item_id' => 'required|integer',
        ]);

        $item = ShopItem::whereIn('class_id', $classIds)->findOrFail($request->input('shop_item_id'));
        $student = Student::whereIn('class_id', $classIds)->findOrFail($request->input('student_id'));

        $redemption = ShopRedemption::create([
            'student_id' => $student->id,
            'shop_item_id' => $item->id,
            'class_id' => $item->class_id,
            'cost' => $item->cost_score,
            'status' => 'pending',
        ]);

        return response()->json(['message' => '兑换请求已创建', 'data' => $redemption], 201);
    }

    public function approveRedemption(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $redemption = ShopRedemption::with(['student.pet', 'shopItem'])
            ->whereIn('class_id', $classIds)->findOrFail($id);

        if ($redemption->status !== 'pending') {
            return response()->json(['message' => '该兑换已处理'], 400);
        }

        /** @var Student $student */
        $student = $redemption->student;
        $item = $redemption->shopItem;
        $itemName = $item->name ?? '未知物品';
        $cost = $redemption->cost;
        $currency = $item->currency_type ?? 'score';

        try {
            // 根据币种选择扣款方式
            if ($currency === 'score') {
                // 积分兑换：扣积分 + 扣宠物经验
                $this->scoreService->spendScore($student, $cost, '兑换：' . $itemName, $teacher->id);
            } else {
                // 钱包币兑换：只扣钱包余额
                app(CurrencyService::class)->spend($student->id, $currency, $cost, '兑换：' . $itemName);
            }

            $redemption->update([
                'status' => 'approved',
                'approved_by' => $teacher->id,
                'approved_at' => now(),
            ]);

            // P4: 特权奖励自动发班级通知
            if ($item && $item->category === 'privilege') {
                Notice::create([
                    'class_id' => $student->class_id,
                    'school_id' => $teacher->school_id,
                    'title' => '特权奖励：' . $itemName,
                    'content' => $student->name . ' 使用 ' . $cost . ' ' . ($currency === 'score' ? '积分' : (Wallet::currencies()[$currency] ?? $currency)) . ' 兑换了「' . $itemName . '」',
                    'type' => 'event',
                    'published_by' => $teacher->id,
                    'is_published' => true,
                    'published_at' => now(),
                ]);
            }

            return response()->json([
                'message' => '已批准兑换，扣除 ' . $cost . ' ' . ($currency === 'score' ? '积分' : (Wallet::currencies()[$currency] ?? $currency)),
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

        $pets = Pet::whereIn('class_id', $classIds)->get();
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
        $classIds = $this->getAccessibleClassIds($teacher);

        $records = \App\Models\Attendance::whereIn('class_id', $classIds)
            ->whereDate('date', today())
            ->with(['student:id,name,student_no,class_id', 'leaveRecord:id,sp_no,leave_type,reason'])
            ->get()
            ->map(fn (\App\Models\Attendance $a) => [
                'id' => $a->id,
                'student_id' => $a->student_id,
                'student_name' => $a->student?->name,
                'student_no' => $a->student?->student_no,
                'status' => $a->status,
                'source' => $a->source,
                'remark' => $a->remark,
                'check_in_time' => $a->sign_in_at?->toDateTimeString(),
                'leave_record' => $a->leaveRecord ? [
                    'sp_no' => $a->leaveRecord->sp_no,
                    'leave_type' => $a->leaveRecord->leave_type,
                    'reason' => $a->leaveRecord->reason,
                ] : null,
            ]);

        return response()->json(['data' => $records]);
    }

    public function startAttendance(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $count = 0;
        $service = app(\App\Services\WechatWorkAttendanceService::class);

        foreach ($classIds as $classId) {
            $count += $service->startAttendanceForClass($classId, $teacher->id, today()->toDateString());
        }

        $leaveCount = \App\Models\Attendance::whereIn('class_id', $classIds)
            ->whereDate('date', today())->where('source', 'wechat_work')->count();
        $msg = "已为 {$count} 名学生创建考勤记录（默认到课）";
        if ($leaveCount > 0) {
            $msg .= "，其中 {$leaveCount} 人已通过企业微信请假";
        }

        return response()->json(['message' => $msg, 'data' => ['total' => $count, 'wechat_leave_count' => $leaveCount]]);
    }

    public function setAttendance(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $request->validate(['status' => 'required|string|in:present,late,leave,absent', 'remark' => 'nullable|string|max:500']);
        $status = $request->input('status');

        $record = \App\Models\Attendance::whereIn('class_id', $classIds)
            ->where('student_id', $studentId)->whereDate('date', today())->firstOrFail();
        $record->update([
            'status' => $status,
            'source' => 'manual',
            'remark' => $request->input('remark'),
            'sign_in_at' => $status === 'present' ? now() : $record->sign_in_at,
        ]);

        return response()->json(['message' => '考勤状态已更新', 'data' => $record]);
    }

    public function markManualLeave(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $request->validate(['remark' => 'required|string|max:500']);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);
        $service = app(\App\Services\WechatWorkAttendanceService::class);
        $record = $service->markManualLeave($studentId, $student->class_id, $teacher->id, today()->toDateString(), $request->input('remark'));

        return response()->json(['message' => '已标记为请假', 'data' => ['id' => $record->id, 'status' => $record->status, 'source' => $record->source, 'remark' => $record->remark]]);
    }

    public function markManualAbsent(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $request->validate(['remark' => 'nullable|string|max:500']);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);
        $service = app(\App\Services\WechatWorkAttendanceService::class);
        $record = $service->markManualAbsent($studentId, $student->class_id, $teacher->id, today()->toDateString(), $request->input('remark'));

        return response()->json(['message' => '已标记为缺勤，建议联系家长确认情况', 'data' => ['id' => $record->id, 'status' => $record->status, 'source' => $record->source, 'remark' => $record->remark]]);
    }

    public function attendanceSummary(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->getAccessibleClassIds($teacher);
        $records = \App\Models\Attendance::whereIn('class_id', $classIds)->whereDate('date', today())->get();
        $present = $records->where('status', 'present')->count();
        $late = $records->where('status', 'late')->count();
        $leave = $records->where('status', 'leave')->count();
        $absent = $records->where('status', 'absent')->count();
        $total = max($records->count(), 1);
        $wechatLeave = $records->where('status', 'leave')->where('source', 'wechat_work')->count();
        $manualLeave = $records->where('status', 'leave')->where('source', 'manual')->count();

        return response()->json(['data' => [
            'present' => $present, 'late' => $late, 'leave' => $leave, 'absent' => $absent,
            'rate' => round($present / $total * 100, 1),
            'wechat_leave_count' => $wechatLeave,
            'manual_leave_count' => $manualLeave,
        ]]);
    }

    private function getAccessibleClassIds($teacher): array
    {
        $ownClassIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id')->toArray();
        $relatedClassIds = \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)->pluck('class_room_id')->toArray();

        return array_values(array_unique(array_merge($ownClassIds, $relatedClassIds)));
    }

    // ============================================================
    // Homework
    // ============================================================

    public function listHomework(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $homework = \App\Models\HomeworkCollection::whereIn('class_id', $classIds)
            ->withCount('submissions')
            ->orderBy('created_at', 'desc')
            ->get()
            ->map(fn ($h) => [
                'id' => $h->id,
                'title' => $h->title,
                'deadline' => $h->deadline?->toDateTimeString(),
                'status' => $h->deadline && $h->deadline->isPast() ? 'closed' : 'active',
                'submission_count' => $h->submissions_count,
                'total_students' => Student::where('class_id', $h->class_id)->count(),
            ]);

        return response()->json(['data' => $homework]);
    }

    public function createHomework(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'title' => 'required|string|max:200',
            'deadline' => 'nullable|date',
            'description' => 'nullable|string',
        ]);

        $hw = \App\Models\HomeworkCollection::create([
            'class_id' => $classIds->first(),
            'teacher_id' => $teacher->id,
            'title' => $request->input('title'),
            'deadline' => $request->input('deadline'),
            'description' => $request->input('description'),
            'status' => 'active',
            'qr_code_token' => \Illuminate\Support\Str::random(16),
        ]);

        return response()->json(['message' => '作业已创建', 'data' => ['id' => $hw->id, 'qr_token' => $hw->qr_code_token]], 201);
    }

    public function getHomework(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $hw = \App\Models\HomeworkCollection::whereIn('class_id', $classIds)->with('submissions.student')->findOrFail($id);

        return response()->json(['data' => [
            'id' => $hw->id,
            'title' => $hw->title,
            'deadline' => $hw->deadline?->toDateTimeString(),
            'description' => $hw->description,
            'status' => $hw->deadline && $hw->deadline->isPast() ? 'closed' : 'active',
            'submissions' => $hw->submissions->map(fn ($s) => [
                'student_name' => $s->student?->name,
                'submitted_at' => $s->submitted_at?->toDateTimeString(),
            ]),
        ]]);
    }

    public function closeHomework(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        \App\Models\HomeworkCollection::whereIn('class_id', $classIds)->findOrFail($id)->update(['status' => 'closed', 'deadline' => now()]);

        return response()->json(['message' => '作业已关闭']);
    }

    public function getHomeworkSubmissions(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $hw = \App\Models\HomeworkCollection::whereIn('class_id', $classIds)->findOrFail($id);

        return response()->json(['data' => $hw->submissions()->with('student')->get()->map(fn ($s) => [
            'student_name' => $s->student?->name,
            'student_no' => $s->student?->student_no,
            'submitted_at' => $s->submitted_at?->toDateTimeString(),
            'content' => $s->content,
            'file_urls' => $s->file_urls,
        ])]);
    }

    public function getHomeworkQrCode(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $hw = \App\Models\HomeworkCollection::whereIn('class_id', $classIds)->findOrFail($id);

        return response()->json(['data' => [
            'qr_url' => url("/api/v1/common/homework-submit/{$hw->qr_code_token}"),
            'token' => $hw->qr_code_token,
        ]]);
    }

    // ============================================================
    // Quizzes
    // ============================================================

    public function listQuizzes(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        return response()->json(['data' => \App\Models\Quiz::whereIn('class_id', $classIds)
            ->withCount('submissions')
            ->orderBy('created_at', 'desc')
            ->get()
            ->map(fn ($q) => [
                'id' => $q->id,
                'title' => $q->title,
                'time_limit' => $q->time_limit,
                'status' => $q->status,
                'is_active' => $q->status === 'active',
                'submission_count' => $q->submissions_count,
            ])]);
    }

    public function createQuiz(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'title' => 'required|string|max:200',
            'question_bank_id' => 'nullable|integer',
            'time_limit' => 'nullable|integer|min:0',
            'subject' => 'nullable|string|max:50',
            'auto_grade' => 'nullable|boolean',
        ]);

        $quiz = \App\Models\Quiz::create([
            'class_id' => $classIds->first(),
            'teacher_id' => $teacher->id,
            'title' => $request->input('title'),
            'question_bank_id' => $request->input('question_bank_id'),
            'time_limit' => $request->input('time_limit', 0),
            'subject' => $request->input('subject', ''),
            'auto_grade' => $request->boolean('auto_grade', true),
            'status' => 'draft',
        ]);

        return response()->json(['message' => '测验已创建', 'data' => ['id' => $quiz->id]], 201);
    }

    public function startQuiz(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        \App\Models\Quiz::whereIn('class_id', $classIds)->findOrFail($id)->update([
            'status' => 'active',
            'started_at' => now(),
        ]);

        return response()->json(['message' => '测验已开始']);
    }

    public function stopQuiz(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $quiz = \App\Models\Quiz::whereIn('class_id', $classIds)->findOrFail($id);
        $quiz->update([
            'status' => 'finished',
            'ended_at' => now(),
        ]);

        return response()->json(['message' => '测验已结束', 'data' => [
            'submissions' => $quiz->submissions()->count(),
            'avg_score' => round((float) $quiz->submissions()->avg('score'), 1),
        ]]);
    }

    public function getQuizStats(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $quiz = \App\Models\Quiz::whereIn('class_id', $classIds)->findOrFail($id);
        $subs = $quiz->submissions()->with('student')->get();

        return response()->json(['data' => [
            'title' => $quiz->title,
            'total' => $subs->count(),
            'avg_score' => round((float) $subs->avg('score'), 1),
            'max_score' => $subs->max('score'),
            'min_score' => $subs->min('score'),
            'submissions' => $subs->map(fn ($s) => [
                'student_name' => $s->student?->name,
                'score' => $s->score,
            ]),
        ]]);
    }

    // ============================================================
    // Question Banks
    // ============================================================

    public function listQuestionBanks(Request $request): JsonResponse
    {
        $teacher = $request->user();

        return response()->json(['data' => \App\Models\QuestionBank::where('teacher_id', $teacher->id)
            ->orWhere('is_public', true)
            ->withCount('questions')
            ->orderBy('created_at', 'desc')
            ->get()
            ->map(fn ($b) => [
                'id' => $b->id,
                'title' => $b->title,
                'subject' => $b->subject,
                'is_public' => $b->is_public,
                'question_count' => $b->questions_count,
            ])]);
    }

    public function createQuestionBank(Request $request): JsonResponse
    {
        $request->validate([
            'title' => 'required|string|max:200',
            'subject' => 'nullable|string|max:50',
            'is_public' => 'nullable|boolean',
        ]);

        $bank = \App\Models\QuestionBank::create([
            'teacher_id' => $request->user()->id,
            'title' => $request->input('title'),
            'subject' => $request->input('subject', ''),
            'is_public' => $request->boolean('is_public', false),
        ]);

        return response()->json(['message' => '题库已创建', 'data' => ['id' => $bank->id]], 201);
    }

    public function addQuestion(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $bank = \App\Models\QuestionBank::where('teacher_id', $teacher->id)->findOrFail($id);

        $request->validate([
            'type' => 'required|string|in:single,multiple,fill,truefalse,short',
            'content' => 'required|string',
            'options' => 'nullable|array',
            'answer' => 'required|string',
            'points' => 'nullable|integer|min:1',
        ]);

        $q = \App\Models\Question::create([
            'question_bank_id' => $bank->id,
            'type' => $request->input('type'),
            'content' => $request->input('content'),
            'options' => $request->input('options'),
            'answer' => $request->input('answer'),
            'points' => $request->input('points', 5),
        ]);

        return response()->json(['message' => '题目已添加', 'data' => ['id' => $q->id]], 201);
    }

    public function getQuestions(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $bank = \App\Models\QuestionBank::where('teacher_id', $teacher->id)->orWhere('is_public', true)->findOrFail($id);

        return response()->json(['data' => $bank->questions()->get()->map(fn ($q) => [
            'id' => $q->id,
            'type' => $q->type,
            'content' => $q->content,
            'options' => $q->options,
            'points' => $q->points,
        ])]);
    }

    // ============================================================
    // Grades
    // ============================================================

    public function listGrades(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $query = \App\Models\Grade::whereIn('class_id', $classIds)->with('student:id,name');

        if ($examName = $request->input('exam_name')) {
            $query->where('exam_name', $examName);
        }
        if ($subject = $request->input('subject')) {
            $query->where('subject', $subject);
        }

        $grades = $query->orderBy('score', 'desc')->paginate(50);

        return response()->json([
            'data' => $grades->items(),
            'meta' => [
                'current_page' => $grades->currentPage(),
                'last_page' => $grades->lastPage(),
                'per_page' => $grades->perPage(),
                'total' => $grades->total(),
            ],
        ]);
    }

    public function inputGrades(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'exam_name' => 'required|string|max:50',
            'subject' => 'required|string|max:50',
            'grades' => 'required|array|min:1',
            'grades.*.student_id' => 'required|integer',
            'grades.*.score' => 'required|numeric|min:0',
        ]);

        $count = 0;
        foreach ($request->input('grades') as $g) {
            $student = Student::whereIn('class_id', $classIds)->find($g['student_id']);
            if (!$student) {
                continue;
            }

            \App\Models\Grade::updateOrCreate(
                [
                    'class_id' => $student->class_id,
                    'student_id' => $g['student_id'],
                    'exam_name' => $request->input('exam_name'),
                    'subject' => $request->input('subject'),
                ],
                [
                    'teacher_id' => $teacher->id,
                    'score' => $g['score'],
                ]
            );
            $count++;
        }

        return response()->json(['message' => "已录入 {$count} 条成绩"]);
    }

    public function getGradeStats(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'exam_name' => 'required|string',
            'subject' => 'required|string',
        ]);

        $stats = \App\Models\Grade::classStats(
            $classIds->first(),
            $request->input('exam_name'),
            $request->input('subject'),
        );

        return response()->json(['data' => $stats]);
    }

    public function getGradeDistribution(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'exam_name' => 'required|string',
            'subject' => 'required|string',
        ]);

        $distribution = \App\Models\Grade::scoreDistribution(
            $classIds->first(),
            $request->input('exam_name'),
            $request->input('subject'),
        );

        return response()->json(['data' => $distribution]);
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
            ['label' => '重点关注学生', 'prompt' => '请分析班上学情'],
        ]]);
    }
}
