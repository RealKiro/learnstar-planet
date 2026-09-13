<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\Notice;
use App\Models\Pet;
use App\Models\PetCollection;
use App\Models\Score;
use App\Models\ScoreRule;
use App\Models\ShopRedemption;
use App\Models\Student;
use App\Models\Wallet;
use App\Services\AttendanceService;
use App\Services\BroadcastService;
use App\Services\CurrencyService;
use App\Services\DisplayEventService;
use App\Services\LeaderboardService;
use App\Services\NoticeService;
use App\Services\PkService;
use App\Services\ReportService;
use App\Services\ScoreService;
use App\Services\ShopService;
use App\Services\TeacherClassScope;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
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
    ) {
    }

    // ============================================================
    // 我的班级
    // ============================================================

    public function myClasses(Request $request): JsonResponse
    {
        $teacher = $request->user();
        // API 机器人账号：返回本校全部启用中的班级（供外部系统枚举可用班级）
        if ($teacher->isApiBot()) {
            $assignments = ClassRoom::where('school_id', $teacher->school_id)
                ->where('status', 'active')
                ->orderBy('id')
                ->get(['id', 'name', 'grade'])
                ->map(fn (ClassRoom $c) => [
                    'class_id' => $c->id,
                    'class_name' => $c->name,
                    'grade' => $c->grade,
                    'role' => 'api_bot',
                ]);

            return response()->json(['data' => $assignments]);
        }

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

        $isAssigned = $teacher->isApiBot()
            ? ClassRoom::where('school_id', $teacher->school_id)->where('id', $classId)->where('status', 'active')->exists()
            : \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
                ->where('class_room_id', $classId)
                ->exists();

        if (!$isAssigned) {
            return response()->json(['message' => '您未被分配到此班级'], 403);
        }

        // 存储当前班级到用户设置，后续所有 API 都使用此班级
        $teacher->setSetting('active_class_id', $classId);

        return response()->json(['message' => '已切换', 'data' => ['active_class_id' => $classId]]);
    }

    // ============================================================
    // Mode Management: classroom_display | teacher_manage
    // ============================================================

    public function getMode(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $mode = $teacher->getSetting('display_mode', 'classroom_display');
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
            'password' => 'required_unless:mode,classroom_display|string|nullable',
        ]);

        $mode = $request->input('mode');
        $teacher->setSetting('display_mode', $mode);

        if ($classId = $request->input('class_id')) {
            $isAssigned = $teacher->isApiBot()
                ? ClassRoom::where('school_id', $teacher->school_id)->where('id', (int) $classId)->where('status', 'active')->exists()
                : \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
                    ->where('class_room_id', (int) $classId)
                    ->exists();
            if ($isAssigned) {
                $teacher->setSetting('active_class_id', (int) $classId);
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

        $isAssigned = $teacher->isApiBot()
            ? ClassRoom::where('school_id', $teacher->school_id)->where('id', (int) $classId)->where('status', 'active')->exists()
            : \App\Models\ClassRoomTeacher::where('user_id', $teacher->id)
                ->where('class_room_id', $classId)
                ->exists();
        if (!$isAssigned) {
            return response()->json(['message' => '您未被分配到此班级'], 403);
        }

        // Load class room with active students
        $classRoom = ClassRoom::findOrFail($classId);

        /** @var \Illuminate\Database\Eloquent\Collection<int, Student> $students */
        $students = Student::where('class_id', $classId)
            ->where('status', 'active')
            ->with('pet')
            ->orderByRaw('CAST(student_no AS UNSIGNED) ASC, id ASC')
            ->get();

        // Pet overview for all students
        $pets = $students->map(function (Student $s): array {
            $pet = $s->pet;
            $stage = $pet ? $pet->currentStage() : ['emoji' => '\ud83e\udd14', 'name' => '未孵化', 'title' => ''];

            return [
                'student_id' => $s->id,
                'student_name' => $s->name,
                'total_score' => $s->total_score,
                'has_pet' => $pet !== null,
                'pet_name' => $pet?->name,
                'pet_species' => $pet?->species,
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
            'student_count' => $students->count(),
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

            // 推送给班级大屏
            try {
                app(DisplayEventService::class)->publish($classId, 'broadcast', [
                    'id' => $broadcast->id,
                    'type' => $broadcast->type,
                    'content' => $broadcast->content,
                    'display_seconds' => $broadcast->display_seconds,
                    'voice_enabled' => $broadcast->voice_enabled,
                    'created_at' => $broadcast->created_at?->toIso8601String(),
                ]);
            } catch (\Throwable $e) {
                logger()->warning('Display broadcast publish failed: ' . $e->getMessage());
            }

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

        // 推送给班级大屏
        try {
            app(DisplayEventService::class)->publish($classId, 'notice', [
                'id' => $notice->id,
                'title' => $notice->title,
                'content' => $notice->content,
                'type' => $notice->type,
                'published_at' => $notice->published_at?->toIso8601String(),
            ]);
        } catch (\Throwable $e) {
            logger()->warning('Display notice publish failed: ' . $e->getMessage());
        }

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
        $classIds = $this->teacherClassIds($teacher);

        if ($classIds->isEmpty()) {
            return response()->json(['data' => [
                'class_name' => '', 'grade' => '', 'student_count' => 0,
                'total_score' => 0, 'avg_pet_level' => 0, 'peak_count' => 0, 'weekly_score' => 0,
                'pending_redemptions' => 0, 'star_student' => null, 'top5' => [], 'recent_news' => [],
            ]]);
        }

        // 教师当前激活班级（未设置时取第一个）
        $activeClassId = $teacher->getSetting('active_class_id') ?: $classIds->first();
        $class = \App\Models\ClassRoom::find($activeClassId) ?? \App\Models\ClassRoom::find($classIds->first());
        if (!$class) {
            return response()->json(['data' => [
                'class_name' => '', 'grade' => '', 'student_count' => 0,
                'total_score' => 0, 'avg_pet_level' => 0, 'peak_count' => 0, 'weekly_score' => 0,
                'pending_redemptions' => 0, 'star_student' => null, 'top5' => [], 'recent_news' => [],
            ]]);
        }

        $pendingRedemptions = \App\Models\ShopRedemption::whereIn('class_id', $classIds)
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
        $recentNews = \App\Models\Score::whereIn('student_id', $students->pluck('id'))
            ->with('student:id,name')->orderBy('created_at', 'desc')->take(20)->get()
            ->map(fn ($s) => [
                'icon' => $s->amount > 0 ? '🎉' : '📝',
                'text' => ($s->student->name ?? '同学') . ' ' . ($s->amount > 0 ? '+' . $s->amount : $s->amount) . '分 — ' . ($s->reason ?? ''),
            ])
            ->unique('text')->take(5)->values();

        return response()->json(['data' => [
            'class_name' => $class->name,
            'grade' => $class->grade,
            'student_count' => $count,
            'total_score' => (int) $totalScore,
            'avg_pet_level' => $avgLevel,
            'peak_count' => $peakCount,
            'weekly_score' => (int) \App\Models\Score::whereIn('student_id', $students->pluck('id'))
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
        ]]);
    }

    // ============================================================
    // Student Management
    // ============================================================

    public function listStudents(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $query = Student::whereIn('class_id', $classIds)
            ->with('classRoom:id,name,grade');

        if ($request->has('search')) {
            $search = $request->input('search');
            $query->where(function ($q) use ($search) {
                $q->where('name', 'like', "%{$search}%")
                  ->orWhere('student_no', 'like', "%{$search}%");
            });
        }

        $students = $query->with('pet')->orderBy('name')->paginate(50);

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
        $students = $request->input('students', []);
        $imported = 0;
        $skipped = [];

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
            $name = trim((string) $data['name']);
            $studentNo = trim((string) ($data['student_no'] ?? ''));

            // 查重：同班同学号/同班同名已存在则跳过（重复导入不再重复建人）
            $dupQuery = Student::where('class_id', $classRoom->id);
            $dup = $studentNo !== ''
                ? $dupQuery->where('student_no', $studentNo)->exists()
                : $dupQuery->where('name', $name)->exists();
            if ($dup) {
                $skipped[] = $name . ($studentNo !== '' ? "（学号 {$studentNo}）" : '') . '：' . $classRoom->name . ' 已存在';
                continue;
            }
            // 跨班冲突防护：同学号已在同校其他班级 → 疑似转班，跳过并提示
            if ($studentNo !== '') {
                $crossDup = Student::with('classRoom:id,name')
                    ->where('student_no', $studentNo)
                    ->where('class_id', '!=', $classRoom->id)
                    ->where('status', 'active')
                    ->whereHas('classRoom', function ($q) use ($classRoom) {
                        $q->where('school_id', $classRoom->school_id);
                    })
                    ->first();
                if ($crossDup) {
                    $skipped[] = $name . "（学号 {$studentNo}）：已存在于 "
                        . ($crossDup->classRoom->name ?? '其他班级')
                        . '，如为转班请联系管理员使用批量转班';
                    continue;
                }
            }
            Student::create([
                'class_id' => $classRoom->id,
                'name' => $name,
                'gender' => $data['gender'] ?? '未知',
                'student_no' => $studentNo !== '' ? $studentNo : null,
                'total_score' => 0,
                'status' => 'active',
            ]);
            $imported++;
        }

        return response()->json([
            'message' => "成功导入 {$imported} 名学生" . (count($skipped) > 0 ? "，跳过 " . count($skipped) . " 条重复/冲突记录" : ''),
            'data' => ['imported_count' => $imported, 'skipped' => $skipped],
        ]);
    }

    public function updateStudent(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);

        $student->update($request->only(['name', 'gender', 'student_no']));

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
        $class = ClassRoom::whereIn('id', $classIds)->find($request->input('class_id'));
        if (!$class) {
            return response()->json(['message' => '只能在自己管理的班级添加学生'], 403);
        }
        $gender = $request->input('gender');
        if (in_array($gender, ['男生', '男'], true)) {
            $gender = '男';
        } elseif (in_array($gender, ['女生', '女'], true)) {
            $gender = '女';
        } else {
            $gender = '未知';
        }
        $student = Student::create([
            'class_id' => $class->id,
            'name' => $request->input('name'),
            'gender' => $gender,
            'student_no' => $request->input('student_no'),
            'total_score' => 0,
            'status' => 'active',
        ]);

        return response()->json([
            'message' => '学生「' . $student->name . '」已添加',
            'data' => $student,
        ], 201);
    }

    public function deleteStudent(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);

        $student->delete();

        return response()->json(['message' => '学生「' . $student->name . '」已删除']);
    }

    // ============================================================
    // Score Management
    // ============================================================

    public function scoreSummary(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

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

        $recent = Score::whereIn('class_id', $classIds)
            ->with('student:id,name')
            ->orderBy('created_at', 'desc')
            ->limit(20)
            ->get()
            ->map(fn ($s) => [
                'id' => $s->id,
                'student_name' => $s->student?->name,
                'amount' => $s->amount,
                'reason' => $s->reason,
                'created_at' => $s->created_at?->toDateTimeString(),
            ]);

        return response()->json(['data' => $recent]);
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
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        // 本班班级规则 + 本校学校级规则（class_id=null 且 school_id=本校），避免跨校泄漏
        $rules = ScoreRule::where(function ($q) use ($classIds, $teacher) {
            $q->whereIn('class_id', $classIds)
              ->orWhere(function ($q2) use ($teacher) {
                  $q2->whereNull('class_id')->where('school_id', $teacher->school_id);
              });
        })->orderBy('sort_order')->get();

        // 无规则时自动创建默认规则（学校级别，同校所有教师共享）
        if ($rules->isEmpty() && $teacher->school_id) {
            $defaults = \App\Services\ScoreRuleService::DEFAULT_RULES;
            // 以 school_id 级别创建，class_id = null，全校共享
            foreach ($defaults as $i => $d) {
                ScoreRule::create([
                    'class_id' => null,
                    'school_id' => $teacher->school_id,
                    'name' => $d['name'], 'amount' => $d['amount'],
                    'category' => $d['category'], 'is_positive' => $d['is_positive'],
                    'is_active' => true, 'sort_order' => $i,
                ]);
            }
            $rules = ScoreRule::where('school_id', $teacher->school_id)
                ->whereNull('class_id')->orderBy('sort_order')->get();
        }

        return response()->json(['data' => $rules]);
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
        $classIds = $this->teacherClassIds($teacher);
        $rule = ScoreRule::whereIn('class_id', $classIds)->orWhereNull('class_id')->findOrFail($id);

        $rule->update($request->only(['name', 'amount', 'category', 'is_positive', 'is_active', 'sort_order']));

        return response()->json(['message' => '规则更新成功', 'data' => $rule]);
    }

    public function deleteScoreRule(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);
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
        $classIds = $this->teacherClassIds($teacher);

        $pets = Pet::whereIn('class_id', $classIds)
            ->with('student:id,name')
            ->get()
            /** @phpstan-ignore-next-line argument.unresolvableType */
            ->map(fn (\App\Models\Pet $p) => [
                'id' => $p->id,
                'student_id' => $p->student_id,
                'student_name' => $p->student?->name,
                'name' => $p->name,
                'species' => $p->species ?: 'zhulong',
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
        $classIds = $this->teacherClassIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['message' => '该学生还没有宠物'], 404);
        }

        $stage = $pet->currentStage();

        return response()->json(['data' => [
            'id' => $pet->id,
            'name' => $pet->name,
            'species' => $pet->species ?: 'zhulong',
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
        $classIds = $this->teacherClassIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            return response()->json(['message' => '该学生还没有宠物'], 404);
        }

        $pet->feed();
        $this->leaderboardService->updatePetLevel($student->class_id, $student->id, $pet->level);

        // 推送给班级大屏
        try {
            app(DisplayEventService::class)->publish($student->class_id, 'pet_update', [
                'student_id' => $student->id,
                'student_name' => $student->name,
                'type' => 'feed',
                'mood' => $pet->mood,
                'level' => $pet->level,
                'experience' => $pet->experience,
            ]);
        } catch (\Throwable $e) {
            logger()->warning('Display pet event failed: ' . $e->getMessage());
        }

        return response()->json(['message' => "已喂养「{$pet->name}」", 'data' => [
            'mood' => $pet->mood,
            'level' => $pet->level,
        ]]);
    }

    public function renamePet(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);
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
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        if ($classIds->isEmpty()) {
            return response()->json(['message' => '没有可管理的班级'], 400);
        }

        $classId = $classIds->first();
        $class = \App\Models\ClassRoom::findOrFail($classId);

        $totalScore = \App\Models\Student::where('class_id', $classId)
            ->where('status', 'active')
            ->sum('total_score');

        return response()->json(['data' => [
            'id' => $class->id,
            'name' => $class->name,
            'grade' => $class->grade,
            'student_count' => \App\Models\Student::where('class_id', $classId)->where('status', 'active')->count(),
            'total_score' => (int) $totalScore,
            'class_points' => (int) ($class->settings['class_points'] ?? 0),
            'settings' => $class->settings,
            'display_code' => \App\Services\DisplayCodeService::generate($class),
        ]]);
    }

    /**
     * 切换班级宠物系列
     */
    public function switchSeries(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $request->validate([
            'series_id' => 'required|string|max:50',
        ]);

        $seriesId = $request->input('series_id');
        $validSeries = ['myth', 'pokemon', 'national', 'digimon', 'magic', 'prehistoric', 'constellation', 'festival', 'qixia', 'dongfang'];

        if (!in_array($seriesId, $validSeries, true)) {
            return response()->json(['message' => '无效的系列ID，可选值：' . implode(', ', $validSeries)], 422);
        }

        if ($classIds->isEmpty()) {
            return response()->json(['message' => '没有可管理的班级'], 400);
        }

        $classId = $classIds->first();
        $class = \App\Models\ClassRoom::findOrFail($classId);
        $settings = $class->settings ?? [];
        $settings['pet_series'] = $seriesId;
        $class->settings = $settings;
        $class->save();

        // 发放「免费自选」机会：整班切换后 3 天内每人可免费切换一次当前类别的宠物，过期作废
        $students = \App\Models\Student::where('class_id', $classId)
            ->where('status', 'active')
            ->get();
        foreach ($students as $student) {
            Cache::put("pet_free_pick:{$student->id}", 1, now()->addDays(3));
        }

        return response()->json([
            'message' => "已切换系列为「{$seriesId}」，全班 {$students->count()} 人各获一次免费自选该系列宠物的机会",
            'data' => [
                'series_id' => $seriesId,
                'class_id' => $classId,
                'free_pick_granted' => true,
                'granted_students' => $students->count(),
            ],
        ]);
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
            'pet_species' => 'required|string|max:50',
            'pet_name' => 'nullable|string|max:20',
        ]);

        $petSpecies = $request->input('pet_species');
        $petName = $request->input('pet_name', $petSpecies);
        $now = now();

        $pet = $student->pet;

        // ===== 同物种切换守卫(不扣费、不重置冷却) =====
        if ($pet && $pet->species === $petSpecies) {
            return response()->json(['message' => '当前已经是这只宠物啦'], 422);
        }

        // ===== 类别限制：只能在本班当前类别内更换，不能跨类别领养 =====
        // 注：旧系列 id(cosmic/cute/all 等)在 speciesPoolForSeries 返回空池 → 视为不限制
        $classSeries = ClassRoom::find($student->class_id)?->settings['pet_series'] ?? null;
        $seriesPool = $classSeries ? Pet::speciesPoolForSeries($classSeries) : [];
        if ($classSeries && !empty($seriesPool) && !in_array($petSpecies, $seriesPool, true)) {
            return response()->json([
                'message' => '只能领养当前类别「' . $classSeries . '」的宠物，不能跨类别领养',
            ], 422);
        }

        // ===== 免费自选：整班切换后的一次机会（限当前类别，免费） =====
        $usedFreePick = false;
        $switchCost = $pet ? Pet::switchCost($pet->level) : 0;
        if ($pet && Cache::has("pet_free_pick:{$student->id}")) {
            $usedFreePick = true;
            $switchCost = 0;
        }

        // ===== 目标物种收藏进度（切回时恢复） =====
        $collection = $pet
            ? PetCollection::where('student_id', $student->id)->where('species', $petSpecies)->first()
            : null;

        if ($pet) {
            // 1) 先扣积分（等级越高越贵；免费自选不扣）——积分不足直接拒绝，不留脏数据
            if (!$usedFreePick) {
                if ($student->total_score < $switchCost) {
                    return response()->json(['message' => "积分不足，更换宠物需 {$switchCost} 积分"], 400);
                }
                $student->total_score -= $switchCost;
                $student->save();
            }

            // 2) 保存当前宠物进度到图鉴（进度全保留）
            PetCollection::updateOrCreate(
                ['student_id' => $student->id, 'species' => $pet->species],
                ['level' => $pet->level, 'experience' => $pet->experience, 'mood' => $pet->mood, 'is_active' => false]
            );

            // 3) 恢复目标物种进度（新物种为初始形态；进度全保留）
            if ($collection) {
                $pet->species = $petSpecies;
                $pet->level = $collection->level;
                $pet->experience = $collection->experience;
                $pet->mood = $collection->mood;
            } else {
                $pet->species = $petSpecies;
                $pet->level = 1;
                $pet->experience = 0;
                $pet->mood = 80;
            }
            $pet->name = $petName;
            $pet->last_switched_at = $now;
            $pet->save();

            // 4) 目标物种标记激活
            PetCollection::updateOrCreate(
                ['student_id' => $student->id, 'species' => $petSpecies],
                ['level' => $pet->level, 'experience' => $pet->experience, 'mood' => $pet->mood, 'is_active' => true]
            );

            // 免费自选机会使用即失效
            if ($usedFreePick) {
                Cache::forget("pet_free_pick:{$student->id}");
            }

            return response()->json([
                'message' => $usedFreePick
                    ? '✅ 已使用整班切换的免费自选机会！'
                    : '宠物已更换为「' . $petName . '」（扣除 ' . $switchCost . ' 积分）',
                'data' => [
                    'pet_name' => $pet->name,
                    'pet_species' => $pet->species,
                    'level' => $pet->level,
                    'experience' => $pet->experience,
                    'cost' => $switchCost,
                    'free_pick_used' => $usedFreePick,
                ],
            ]);
        }

        // 无宠物：创建新宠物并收入图鉴
        $pet = Pet::create([
            'student_id' => $student->id,
            'class_id' => $student->class_id,
            'name' => $petName,
            'species' => $petSpecies,
            'level' => 1,
            'experience' => 0,
            'mood' => 80,
        ]);
        PetCollection::create([
            'student_id' => $student->id,
            'species' => $petSpecies,
            'level' => 1,
            'experience' => 0,
            'mood' => 80,
            'is_active' => true,
        ]);

        return response()->json([
            'message' => '已为您分配宠物「' . $petName . '」',
            'data' => [
                'pet_name' => $pet->name,
                'pet_species' => $pet->species,
                'level' => $pet->level,
                'experience' => $pet->experience,
            ],
        ]);
    }

    /**
     * 学生宠物图鉴收藏
     */
    public function petCollection(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $activePet = $student->pet;

        // 确保当前激活宠物在收藏中
        if ($activePet && $activePet->species) {
            PetCollection::firstOrCreate(
                ['student_id' => $student->id, 'species' => $activePet->species],
                ['level' => $activePet->level, 'experience' => $activePet->experience, 'mood' => $activePet->mood, 'is_active' => true]
            );
        }

        $collections = PetCollection::where('student_id', $student->id)->orderBy('species')->get();

        return response()->json(['data' => [
            'student_id' => $student->id,
            'student_name' => $student->name,
            'total_score' => $student->total_score,
            'unlock_slots' => PetCollection::unlockSlotsForScore($student->total_score),
            'class_series' => ClassRoom::find($student->class_id)?->settings['pet_series'] ?? null,
            'active_species' => $activePet?->species,
            'collection' => $collections->map(fn (PetCollection $c) => [
                'species' => $c->species,
                'level' => $c->level,
                'experience' => $c->experience,
                'mood' => $c->mood,
                'is_active' => $c->is_active,
            ])->values(),
        ]]);
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
        $setting = \App\Models\AiSetting::where('school_id', $teacher->school_id)->first();

        $enabled = false;
        if ($setting && $setting->enabled) {
            $providers = $setting->providers ?: [];
            // 任一启用的供应商配了 api_key 即视为已配置
            $hasProviderKey = collect($providers)->contains(
                fn ($p) => !empty($p['is_active']) && !empty($p['api_key'])
            );
            // 兼容旧字段（单 provider + api_key）
            $hasLegacyKey = !empty($setting->provider) && !empty($setting->api_key);
            $enabled = $hasProviderKey || $hasLegacyKey;
        }

        return response()->json(['data' => ['enabled' => $enabled]]);
    }

    public function aiChat(Request $request): JsonResponse
    {
        $request->validate(['message' => 'required|string|max:2000']);

        $teacher = $request->user();
        $settings = \App\Models\AiSetting::where('school_id', $teacher->school_id)->first();
        if (!$settings || !$settings->enabled) {
            return response()->json(['data' => ['reply' => 'AI 功能未启用，请联系管理员配置']]);
        }

        // 从多供应商配置中查找启用的供应商，兼容旧版单供应商配置
        $activeProvider = null;
        foreach ($settings->providers ?: [] as $p) {
            if (!empty($p['is_active']) && !empty($p['api_key'])) {
                $activeProvider = $p;
                break;
            }
        }
        if (!$activeProvider && !empty($settings->api_key)) {
            $activeProvider = [
                'id' => $settings->provider ?: 'openai',
                'api_key' => $settings->api_key,
                'api_base' => $settings->api_base,
                'model' => $settings->model ?: 'gpt-3.5-turbo',
            ];
        }
        if (!$activeProvider) {
            return response()->json(['data' => ['reply' => '请先在 AI 中心配置并启用一个供应商']]);
        }

        $classId = $teacher->getSetting('active_class_id') ?: null;

        $conversation = \App\Models\AiConversation::create([
            'school_id' => $teacher->school_id,
            'class_id' => $classId,
            'student_name' => '教师',
            'provider' => $activeProvider['id'],
            'question' => $request->input('message'),
            'status' => 'pending',
        ]);

        try {
            $ai = new \App\Services\AiService();
            $result = $ai->chat(
                provider: $activeProvider['id'],
                apiKey: $activeProvider['api_key'],
                model: $activeProvider['model'] ?: 'gpt-3.5-turbo',
                question: $request->input('message'),
                apiBase: $activeProvider['api_base'] ?? null,
                maxTokens: $settings->max_tokens,
            );
            $reply = $result['answer'];
            $promptTokens = $result['prompt_tokens'] ?? 0;
            $completionTokens = $result['completion_tokens'] ?? 0;
            $tokensUsed = $result['tokens_used'] ?? ($promptTokens + $completionTokens);
        } catch (\Throwable $e) {
            $reply = 'AI 服务暂时不可用';
            $tokensUsed = 0;
            $promptTokens = 0;
            $completionTokens = 0;
        }

        // 本地精确计费：按供应商单价计算本次费用
        $cost = app(\App\Services\AiBilling\AiBillingService::class)->recordUsage($activeProvider, $promptTokens, $completionTokens);
        $currency = $activeProvider['currency'] ?? 'CNY';

        $conversation->update([
            'answer' => $reply,
            'tokens_used' => $tokensUsed,
            'prompt_tokens' => $promptTokens,
            'completion_tokens' => $completionTokens,
            'cost' => $cost,
            'currency' => $currency,
            'status' => 'completed',
        ]);

        if ($tokensUsed > 0) {
            $settings->increment('tokens_used', $tokensUsed);
            $providers = $settings->providers ?: [];
            foreach ($providers as &$p) {
                if (($p['id'] ?? '') === ($activeProvider['id'] ?? '')) {
                    $p['tokens_used'] = ($p['tokens_used'] ?? 0) + $tokensUsed;
                    $p['total_calls'] = ($p['total_calls'] ?? 0) + 1;
                    $p['estimated_cost'] = ($p['estimated_cost'] ?? 0) + $cost;
                    $p['currency'] = $currency;
                    break;
                }
            }
            $settings->providers = $providers;
            $settings->save();
        }

        return response()->json(['data' => ['reply' => $reply]]);
    }

    /**
     * 教师端 AI 用量（前端 AIPage 每次发送后刷新）
     */
    public function getAiUsage(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $settings = \App\Models\AiSetting::where('school_id', $teacher->school_id)->first();

        $active = null;
        if ($settings) {
            foreach ($settings->providers ?: [] as $p) {
                if (!empty($p['is_active']) && !empty($p['api_key'])) {
                    $active = $p;
                    break;
                }
            }
            if (!$active && !empty($settings->api_key)) {
                $active = ['id' => $settings->provider ?: 'openai', 'model' => $settings->model];
            }
        }

        $estimatedCost = 0.0;
        $currency = 'CNY';
        foreach ($settings->providers ?? [] as $p) {
            $estimatedCost += (float) ($p['estimated_cost'] ?? 0);
            $currency = (string) ($p['currency'] ?? $currency);
        }

        return response()->json(['data' => [
            'configured' => $settings !== null && $settings->enabled && $active !== null,
            'provider' => $active['id'] ?? ($settings->provider ?? null),
            'model' => $active['model'] ?? ($settings->model ?? null),
            'tokens_used' => $settings ? (int) $settings->tokens_used : 0,
            'estimated_cost' => $estimatedCost,
            'currency' => $currency,
        ]]);
    }

    /**
     * 教师端 AI 预设命令
     */
    public function getAiCommands(Request $request): JsonResponse
    {
        return response()->json(['data' => [
            ['label' => '📝 本周教学总结', 'prompt' => '请帮我写一份本周教学总结，包含本周教学目标、课堂情况、学生表现和下周教学计划。'],
            ['label' => '🏅 积分规则建议', 'prompt' => '请根据班级日常情况，生成一套适合小学生的积分奖励规则建议。'],
            ['label' => '🎯 班会活动方案', 'prompt' => '请设计一个有趣的小学生班会活动方案，包含活动目标、流程和所需材料。'],
            ['label' => '📋 出练习题', 'prompt' => '请出一组适合本年级学生的练习题，包含题目和参考答案。'],
        ]]);
    }

    // ============================================================
    // 汇率管理（教师端）
    // ============================================================

    public function listExchangeRates(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $schoolId = $teacher->school_id;

        // 首次访问惰性初始化默认汇率，保证积分充值类商品不配汇率也能结算
        $exists = \App\Models\ExchangeRate::where('school_id', $schoolId)->exists();
        if (!$exists) {
            $defaults = [
                // 2:1 防通胀：2 积分 = 1 币
                ['name' => '积分 → 科学币', 'from_currency' => 'score', 'to_currency' => 'science', 'rate' => 0.5],
                ['name' => '积分 → 读书币', 'from_currency' => 'score', 'to_currency' => 'reading', 'rate' => 0.5],
                ['name' => '积分 → 体育币', 'from_currency' => 'score', 'to_currency' => 'class_point', 'rate' => 0.5],
            ];
            foreach ($defaults as $d) {
                \App\Models\ExchangeRate::firstOrCreate(
                    ['school_id' => $schoolId, 'from_currency' => $d['from_currency'], 'to_currency' => $d['to_currency']],
                    ['name' => $d['name'], 'rate' => $d['rate'], 'is_active' => true],
                );
            }
        }

        $rates = \App\Models\ExchangeRate::where('school_id', $schoolId)
            ->orderBy('from_currency')->orderBy('to_currency')->get();

        return response()->json(['data' => $rates]);
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

        $rate = \App\Models\ExchangeRate::create([
            'school_id' => $teacher->school_id,
            'name' => $request->input('name'),
            'from_currency' => $request->input('from_currency'),
            'to_currency' => $request->input('to_currency'),
            'rate' => $request->input('rate'),
            'is_active' => true,
        ]);

        return response()->json(['message' => '汇率已添加', 'data' => $rate], 201);
    }

    public function updateExchangeRate(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $rate = \App\Models\ExchangeRate::where('school_id', $teacher->school_id)->findOrFail($id);
        $request->validate([
            'rate' => 'sometimes|numeric|min:0.01',
            'is_active' => 'sometimes|boolean',
        ]);

        $rate->update($request->only(['rate', 'is_active']));

        return response()->json(['message' => '汇率已更新', 'data' => $rate->fresh()]);
    }

    // ============================================================
    // 兑换中心
    // ============================================================

    public function listWallets(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);
        $students = Student::whereIn('class_id', $classIds)->where('status', 'active')->pluck('id');

        $wallets = Wallet::whereIn('student_id', $students)
            ->with('student:id,name')
            ->get()
            ->map(fn ($w) => [
                'student_id' => $w->student_id,
                'student_name' => $w->student?->name,
                'currency_type' => $w->currency_type,
                'balance' => (int) $w->balance,
            ]);

        return response()->json(['data' => $wallets]);
    }

    /**
     * 兑换记录（仅教师所带班级学生）
     */
    public function exchangeLogs(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = $this->teacherClassIds($teacher);

        $logs = \App\Models\ExchangeLog::with('student:id,name,student_no')
            ->whereIn('student_id', Student::whereIn('class_id', $classIds)->select('id'))
            ->orderBy('created_at', 'desc')
            ->paginate(20);

        return response()->json([
            'data' => collect($logs->items())->map(static function (\App\Models\ExchangeLog $log): array {
                return array_merge($log->toArray(), [
                    'student_name' => $log->student->name ?? '已删除学生',
                    'student_no' => $log->student->student_no ?? '',
                ]);
            }),
            'meta' => [
                'current_page' => $logs->currentPage(),
                'last_page' => $logs->lastPage(),
                'total' => $logs->total(),
            ],
        ]);
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
