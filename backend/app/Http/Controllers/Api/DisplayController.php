<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Broadcast;
use App\Models\ClassRoom;
use App\Models\ClassRoomTeacher;
use App\Models\DisplayLoginLog;
use App\Models\Pet;
use App\Models\Score;
use App\Models\ShopItem;
use App\Models\ShopRedemption;
use App\Models\Student;
use App\Models\User;
use App\Services\DisplayEventService;
use App\Services\ScoreService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Str;
use Symfony\Component\HttpFoundation\StreamedResponse;
use Throwable;

/**
 * DisplayController — 班级大屏单独入口
 *
 * 职责：
 * 1. 教师端：生成/刷新班级大屏码
 * 2. 显示端：班级码 → Token 认证
 * 3. 显示端：SSE 长连接 + 轮询降级
 * 4. 显示端：初始全量数据加载
 *
 * 与教师管理模式完全解耦，不依赖 auth:sanctum 中间件。
 */
class DisplayController extends Controller
{
    private const TOKEN_PREFIX = 'display:token:';
    private const TOKEN_TTL = 86400; // 24 小时
    private const CODE_LOCK_PREFIX = 'display:code_lock:';
    private const MAX_LOGIN_ATTEMPTS = 5;
    private const LOCKOUT_MINUTES = 15;
    private const SSE_HEARTBEAT_INTERVAL = 10; // 秒
    private const SSE_MAX_EXECUTION = 55; // 最大执行秒数（略小于 PHP/Nginx 超时）

    public function __construct(
        private readonly DisplayEventService $eventService,
        private readonly ScoreService $scoreService,
    ) {
    }

    // ============================================================
    // 教师端 API — 班级大屏码管理
    // ============================================================

    /**
     * 获取或创建当前班级的大屏码
     */
    public function getDisplayCode(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = (int) $request->input('class_id', $teacher->getSetting('active_class_id'));

        if (!$classId) {
            return response()->json(['message' => '请先选择班级'], 400);
        }

        $classRoom = ClassRoom::findOrFail($classId);

        // 首次访问时自动生成班级码
        if (empty($classRoom->display_code)) {
            $classRoom->display_code = $this->generateDisplayCode($classRoom);
            $classRoom->display_code_updated_at = now();
            $classRoom->save();

            // 缓存班级码 → 班级 ID 映射
            $this->cacheCodeMapping($classRoom->display_code, $classRoom->id);
        }

        return response()->json(['data' => [
            'code' => $classRoom->display_code,
            'class_name' => $classRoom->name,
            'updated_at' => $classRoom->display_code_updated_at?->toIso8601String(),
            'student_count' => Student::where('class_id', $classId)
                ->where('status', 'active')
                ->count(),
        ]]);
    }

    /**
     * 刷新班级大屏码（旧码立即失效）
     */
    public function refreshDisplayCode(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = (int) $request->input('class_id', $teacher->getSetting('active_class_id'));

        if (!$classId) {
            return response()->json(['message' => '请先选择班级'], 400);
        }

        $classRoom = ClassRoom::findOrFail($classId);

        // 清除旧码缓存
        if (!empty($classRoom->display_code)) {
            Cache::forget(DisplayEventService::codeCacheKey($classRoom->display_code));
        }

        // 生成新码
        $oldCode = $classRoom->display_code;
        $classRoom->display_code = $this->generateDisplayCode($classRoom);
        $classRoom->display_code_updated_at = now();
        $classRoom->save();

        // 缓存新码映射
        $this->cacheCodeMapping($classRoom->display_code, $classRoom->id);

        // 发送 refresh 事件通知当前大屏刷新
        $this->eventService->publish($classId, 'refresh', [
            'old_code' => $oldCode,
            'new_code' => $classRoom->display_code,
        ]);

        return response()->json(['data' => [
            'code' => $classRoom->display_code,
            'class_name' => $classRoom->name,
            'updated_at' => $classRoom->display_code_updated_at?->toIso8601String(),
            'message' => '班级码已刷新，旧码已失效',
        ]]);
    }

    // ============================================================
    // 显示端 API — 班级码登录
    // ============================================================

    /**
     * 班级码登录，返回 SSE Token
     */
    public function login(Request $request): JsonResponse
    {
        $request->validate([
            'code' => 'required|string|max:12',
        ]);

        $code = strtoupper(trim($request->input('code')));

        // 暴力破解防护
        $lockKey = self::CODE_LOCK_PREFIX . $code;
        $attempts = (int) Cache::get($lockKey, 0);

        if ($attempts >= self::MAX_LOGIN_ATTEMPTS) {
            return response()->json([
                'message' => '尝试次数过多，请 ' . self::LOCKOUT_MINUTES . ' 分钟后再试',
            ], 429);
        }

        // 查找班级码（先查缓存，再查 DB）
        $classId = $this->resolveCodeToClassId($code);

        if (!$classId) {
            Cache::put($lockKey, $attempts + 1, now()->addMinutes(self::LOCKOUT_MINUTES));

            return response()->json(['message' => '班级码无效，请检查后重试'], 404);
        }

        // 重置尝试计数
        Cache::forget($lockKey);

        $classRoom = ClassRoom::find($classId);
        if (!$classRoom) {
            return response()->json(['message' => '班级不存在'], 404);
        }

        // 生成显示端 Token
        $token = 'disp_' . Str::random(32);
        $tokenKey = self::TOKEN_PREFIX . $token;

        Cache::put($tokenKey, [
            'class_id' => $classId,
            'class_name' => $classRoom->name,
            'grade' => $classRoom->grade,
            'ip' => $request->ip(),
            'created_at' => now()->toIso8601String(),
        ], now()->addSeconds(self::TOKEN_TTL));

        // 记录登录日志（含 IP 地址）
        try {
            DisplayLoginLog::create([
                'class_id' => $classId,
                'class_code' => $code,
                'ip_address' => $request->ip(),
                'user_agent' => $request->userAgent(),
            ]);
        } catch (\Throwable $e) {
            logger()->warning('Failed to record display login log: ' . $e->getMessage());
        }

        // 获取班级统计数据
        $studentCount = Student::where('class_id', $classId)
            ->where('status', 'active')
            ->count();

        return response()->json(['data' => [
            'token' => $token,
            'expires_in' => self::TOKEN_TTL,
            'class_info' => [
                'id' => $classId,
                'name' => $classRoom->name,
                'grade' => $classRoom->grade,
                'student_count' => $studentCount,
            ],
        ]]);
    }

    // ============================================================
    // 显示端 API — 初始全量数据
    // ============================================================

    /**
     * 获取大屏初始全量数据（连接 SSE 前先调用一次）
     */
    public function initialData(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token 无效或已过期'], 401);
        }

        $classId = (int) $classInfo['class_id'];
        $classRoom = ClassRoom::find($classId);

        if (!$classRoom) {
            return response()->json(['message' => '班级不存在'], 404);
        }

        $students = Student::where('class_id', $classId)
            ->where('status', 'active')
            ->with('pet')
            ->orderByRaw('CAST(student_no AS UNSIGNED) ASC, id ASC')
            ->get();

        // 没有宠物的学生自动分配一只（兼容旧数据）
        $cuteTypes = ['orange_cat', 'husky', 'shiba', 'guinea_pig', 'hamster', 'bunny', 'parrot', 'hedgehog', 'chinchilla', 'teacup_pig', 'sugar_glider', 'alpaca'];
        foreach ($students as $s) {
            if (!$s->pet) {
                $type = $cuteTypes[array_rand($cuteTypes)];
                $s->pet()->create([
                    'student_id' => $s->id,
                    'class_id' => $s->class_id,
                    'name' => $s->name . '的萌宠',
                    'type' => $type,
                    'level' => 0,
                    'experience' => 0,
                    'mood' => 80,
                ]);
                $s->load('pet');
            }
        }

        $pets = $this->formatPets($students);

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

        // 活跃广播（尚未过期的）
        $broadcasts = Broadcast::where('class_id', $classId)
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

        return response()->json(['data' => [
            'class_name' => $classRoom->name,
            'grade' => $classRoom->grade,
            'student_count' => $students->count(),
            'pets' => $pets,
            'recent_scores' => $recentScores,
            'broadcasts' => $broadcasts,
            'server_time' => now()->toIso8601String(),
        ]]);
    }

    // ============================================================
    // 显示端 API — SSE 长连接（核心）
    // ============================================================

    /**
     * SSE 端点：浏览器建立长连接，实时接收事件推送
     *
     * 实现原理:
     * 1. 在 PHP 中循环，每次休眠 1-2 秒后检查 Cache 中是否有新事件
     * 2. 有新事件则立即推送
     * 3. 每 10 秒发一次心跳保持连接
     * 4. 最大运行 55 秒后优雅断开，浏览器自动重连
     * 5. EventSource 内置重连机制
     */
    public function sse(Request $request): StreamedResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            // Token 无效，返回 401 事件让前端跳转到登录页
            $response = new StreamedResponse(function () {
                echo "event: error\n";
                echo "data: {\"code\":401,\"message\":\"Token无效或已过期\"}\n\n";
                ob_flush();
                flush();
            });
            $response->headers->set('Content-Type', 'text/event-stream');
            $response->headers->set('Cache-Control', 'no-cache');
            $response->headers->set('X-Accel-Buffering', 'no');

            return $response;
        }

        $classId = (int) $classInfo['class_id'];
        $lastEventId = (int) ($request->input('last_event_id', 0));

        $response = new StreamedResponse(function () use ($classId, &$lastEventId) {
            // 关闭输出缓冲
            if (ob_get_level()) {
                ob_end_clean();
            }

            $startTime = time();
            $lastHeartbeat = 0;

            while (true) {
                // 检查是否超时
                $elapsed = time() - $startTime;
                if ($elapsed >= self::SSE_MAX_EXECUTION) {
                    break;
                }

                // 消费新事件
                $events = $this->eventService->consume($classId, $lastEventId);

                foreach ($events as $event) {
                    $eventId = $event['id'] ?? 0;
                    $eventType = $event['type'] ?? 'unknown';
                    $eventData = json_encode($event['data'] ?? [], JSON_UNESCAPED_UNICODE);

                    echo "id: {$eventId}\n";
                    echo "event: {$eventType}\n";
                    echo "data: {$eventData}\n\n";

                    $lastEventId = max($lastEventId, $eventId);
                }

                // 心跳（保持连接活跃）
                if (time() - $lastHeartbeat >= self::SSE_HEARTBEAT_INTERVAL) {
                    echo "event: heartbeat\n";
                    echo 'data: {"time":"' . now()->toIso8601String() . "\"}\n\n";
                    $lastHeartbeat = time();
                }

                ob_flush();
                flush();

                // 休眠一小段时间以减少 CPU 占用
                if (empty($events)) {
                    usleep(2_000_000); // 2 秒
                } else {
                    usleep(500_000); // 有事件时加快响应
                }
            }

            // 优雅关闭：发送重连指示
            echo "event: reconnect\n";
            echo "data: {}\n\n";
            ob_flush();
            flush();
        });

        // SSE 必需的头
        $response->headers->set('Content-Type', 'text/event-stream');
        $response->headers->set('Cache-Control', 'no-cache');
        $response->headers->set('Connection', 'keep-alive');
        $response->headers->set('X-Accel-Buffering', 'no'); // 禁用 Nginx 缓冲

        return $response;
    }

    // ============================================================
    // 显示端 API — 轮询降级（SSE 不可用时的备选）
    // ============================================================

    /**
     * 轮询降级端点：当 SSE 连接失败 3 次后使用
     */
    public function poll(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token 无效或已过期'], 401);
        }

        $classId = (int) $classInfo['class_id'];
        $lastEventId = (int) ($request->input('last_event_id', 0));

        $events = [];

        try {
            $events = $this->eventService->consume($classId, $lastEventId ?: null);
        } catch (\Throwable $e) {
            Log::warning('Display poll consume failed: ' . $e->getMessage());
        }

        $maxId = $lastEventId;
        foreach ($events as $e) {
            if (($e['id'] ?? 0) > $maxId) {
                $maxId = $e['id'];
            }
        }

        return response()->json([
            'data' => [
                'events' => $events,
                'last_event_id' => $maxId,
                'server_time' => now()->toIso8601String(),
            ],
        ]);
    }

    // ============================================================
    // 内部辅助
    // ============================================================

    /**
     * 生成班级大屏码
     * 格式: {年级缩写}-{班号}-{4位随机码}
     * 例如: 3-1-A7K2
     */
    private function generateDisplayCode(ClassRoom $classRoom): string
    {
        $grade = $classRoom->grade ?? '0';
        $name = $classRoom->name ?? '';

        // 从班级名提取班号: "三年级（1）班" → "1"
        $classNo = '0';
        if (preg_match('/（(\d+)）班/', $name, $m)) {
            $classNo = $m[1];
        }

        $random = strtoupper(Str::random(4));

        // 确保不重复
        $code = "{$grade}-{$classNo}-{$random}";
        $existing = ClassRoom::where('display_code', $code)->where('id', '!=', $classRoom->id)->exists();
        if ($existing) {
            $code = "{$grade}-{$classNo}-" . strtoupper(Str::random(4));
        }

        return $code;
    }

    /**
     * 缓存班级码 → 班级 ID 映射
     */
    private function cacheCodeMapping(string $code, int $classId): void
    {
        Cache::put(
            DisplayEventService::codeCacheKey($code),
            $classId,
            now()->addDays(30)
        );
    }

    /**
     * 解析班级码 → 班级 ID
     * 先查缓存，再查 DB
     */
    private function resolveCodeToClassId(string $code): ?int
    {
        $cacheKey = DisplayEventService::codeCacheKey($code);
        $classId = Cache::get($cacheKey);

        if ($classId) {
            return (int) $classId;
        }

        // 缓存未命中，查数据库
        $classRoom = ClassRoom::where('display_code', $code)->first();
        if ($classRoom) {
            $this->cacheCodeMapping($code, $classRoom->id);

            return $classRoom->id;
        }

        return null;
    }

    /**
     * 验证请求中的 SSE Token
     */
    private function validateToken(Request $request): ?array
    {
        // 优先从 query 参数取，其次从 Authorization header 取
        $token = $request->input('token', '');
        if (empty($token)) {
            $bearer = $request->bearerToken();
            if ($bearer && (str_starts_with($bearer, 'disp_') || str_starts_with($bearer, 'class_'))) {
                $token = $bearer;
            }
        }

        if (empty($token)) {
            return null;
        }

        // 兼容两种 token 格式：
        // 1. disp_ 开头 → DisplayController 传统 token
        // 2. class_ 开头 → AuthController::classLogin 统一 token
        if (str_starts_with($token, 'disp_')) {
            $data = Cache::get(self::TOKEN_PREFIX . $token);

            return $data ?: null;
        }

        if (str_starts_with($token, 'class_')) {
            $classId = Cache::get('class_token:' . $token);
            if ($classId) {
                return ['class_id' => $classId];
            }

            return null;
        }

        return null;
    }

    /**
     * 格式化学生宠物数据（复用 ClassroomDisplay 逻辑）
     */
    private function formatPets($students): array
    {
        return $students->map(function (Student $s): array {
            $pet = $s->pet;
            $stage = $this->getPetStageInfo($pet->level ?? 0, $pet->type ?? null);

            return [
                'student_id' => $s->id,
                'student_no' => $s->student_no ?? '',
                'student_name' => $s->name,
                'total_score' => $s->total_score,
                'has_pet' => $pet !== null,
                'pet_name' => $pet?->name,
                'pet_type' => $pet?->type,
                'level' => $pet->level ?? 0,
                'experience' => $pet->experience ?? 0,
                'mood' => $pet->mood ?? 50,
                'emoji' => $stage['emoji'] ?? '🌟',
                'stage_name' => $stage['name'] ?? '未孵化',
                'stage_title' => $stage['title'] ?? '',
                'exp_max' => $stage['exp_max'] ?? 100,
                'color' => $stage['color'] ?? '#94a3b8',
                'image' => $stage['image'] ?? null,
            ];
        })->values()->toArray();
    }

    /**
     * 获取宠物阶段信息（复用 Pet 模型定义）
     */
    private function getPetStageInfo(int $level, ?string $petType = null): array
    {
        if ($petType) {
            $stages = Pet::evolutionStagesForType($petType);
        } else {
            $stages = Pet::evolutionStages();
        }
        $stage = $stages[min($level, 10)] ?? $stages[0];
        $stage['exp_max'] = ($level + 1) * 10;

        return $stage;
    }

    // ============================================================
    // 大屏快捷加减分
    // ============================================================

    public function quickScore(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token 无效或已过期'], 401);
        }

        $request->validate([
            'student_id' => 'required|integer',
            'amount' => 'required|integer|in:-5,-3,-1,1,3,5',
        ]);

        $student = Student::where('class_id', $classInfo['class_id'])->find((int) $request->input('student_id'));
        if (!$student) {
            return response()->json(['message' => '学生不存在'], 404);
        }

        $amount = (int) $request->input('amount');
        $teacherId = $this->getClassTeacherId($classInfo['class_id']);

        try {
            $this->scoreService->giveScore($student, $amount, '课堂表现', $teacherId ?: 1);

            return response()->json(['data' => ['total_score' => $student->fresh()->total_score]]);
        } catch (\Throwable $e) {
            return response()->json(['message' => '操作失败'], 500);
        }
    }

    // ============================================================
    // 大屏排行榜
    // ============================================================

    public function quickLeaderboard(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效'], 401);
        }

        $students = Student::where('class_id', $classInfo['class_id'])
            ->where('status', 'active')
            ->orderBy('total_score', 'desc')
            ->take(20)
            ->get(['id', 'name', 'total_score', 'student_no']);

        $data = [];
        foreach ($students as $i => $s) {
            $data[] = ['rank' => $i + 1, 'id' => $s->id, 'name' => $s->name, 'score' => $s->total_score, 'no' => $s->student_no];
        }

        return response()->json(['data' => $data]);
    }

    // ============================================================
    // 大屏商品列表
    // ============================================================

    public function quickShopItems(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效'], 401);
        }

        $items = ShopItem::where('class_id', $classInfo['class_id'])
            ->where('is_active', true)
            ->get(['id', 'name', 'description', 'cost_score', 'stock', 'category']);

        return response()->json(['data' => $items]);
    }

    // ============================================================
    // 大屏快捷兑换（教师操作，学生扣分）
    // ============================================================

    public function quickRedeem(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效'], 401);
        }

        $request->validate(['student_id' => 'required|integer', 'item_id' => 'required|integer']);

        $classId = (int) $classInfo['class_id'];
        $student = Student::where('class_id', $classId)->find((int) $request->input('student_id'));
        $item = ShopItem::where('class_id', $classId)->where('is_active', true)->find((int) $request->input('item_id'));
        if (!$student || !$item) {
            return response()->json(['message' => '学生或商品不存在'], 404);
        }
        if ($student->total_score < $item->cost_score) {
            return response()->json(['message' => '积分不足'], 400);
        }

        $teacherId = $this->getClassTeacherId($classId);

        try {
            $this->scoreService->spendScore($student, $item->cost_score, '兑换：' . $item->name, $teacherId ?: 1);
            ShopRedemption::create([
                'student_id' => $student->id, 'shop_item_id' => $item->id, 'class_id' => $classId,
                'cost' => $item->cost_score, 'status' => 'approved', 'approved_by' => $teacherId, 'approved_at' => now(),
            ]);

            return response()->json(['data' => ['student_name' => $student->name, 'item_name' => $item->name, 'cost' => $item->cost_score, 'total_score' => $student->fresh()->total_score]]);
        } catch (\Throwable $e) {
            return response()->json(['message' => '兑换失败'], 500);
        }
    }

    // ============================================================
    // 学生间积分转赠（鼓励互助）
    // ============================================================

    public function quickTransfer(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效'], 401);
        }

        $request->validate(['from_id' => 'required|integer', 'to_id' => 'required|integer|different:from_id', 'amount' => 'required|integer|min:1|max:100']);

        $classId = (int) $classInfo['class_id'];
        $from = Student::where('class_id', $classId)->find((int) $request->input('from_id'));
        $to = Student::where('class_id', $classId)->find((int) $request->input('to_id'));
        $amount = (int) $request->input('amount');

        if (!$from || !$to) {
            return response()->json(['message' => '学生不存在'], 404);
        }
        if ($from->total_score < $amount) {
            return response()->json(['message' => '积分不足'], 400);
        }

        $teacherId = $this->getClassTeacherId($classId);

        try {
            $this->scoreService->giveScore($from, -$amount, '转赠给 ' . $to->name, $teacherId ?: 1);
            $this->scoreService->giveScore($to, $amount, '来自 ' . $from->name . ' 的转赠', $teacherId ?: 1);

            return response()->json(['data' => ['from_name' => $from->name, 'to_name' => $to->name, 'amount' => $amount]]);
        } catch (\Throwable $e) {
            return response()->json(['message' => '转赠失败'], 500);
        }
    }

    // ============================================================
    // 教室端 API（设计文档4大模块，班级码Token）
    // ============================================================

    /**
     * 教室端 · 班级总览
     */
    
    public function classSettings(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效'], 401);
        }
        $class = ClassRoom::find($classInfo['class_id']);
        if (!$class) {
            return response()->json(['data' => ['pet_series' => null]]);
        }
        return response()->json(['data' => ['pet_series' => $class->settings['pet_series'] ?? null]]);
    }

    public function classroomDashboard(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $classId = (int) $classInfo['class_id'];
        $class = ClassRoom::find($classId);
        if (!$class) {
            return response()->json(['message' => '班级不存在'], 404);
        }

        $students = Student::where('class_id', $classId)->where('status', 'active')->with('pet')->get();
        $totalScore = $students->sum('total_score');
        $count = $students->count();
        $avgLevel = $count > 0 ? round($students->avg(fn ($s) => $s->pet->level ?? 0), 1) : 0;
        $peakCount = $students->filter(fn ($s) => $s->pet && $s->pet->level >= 10)->count();

        $sorted = $students->sortByDesc('total_score')->values();
        $top5 = $sorted->take(5)->map(fn ($s) => [
            'name' => $s->name,
            'score' => $s->total_score,
            'pet_name' => $s->pet->name ?? '',
            'pet_species' => $s->pet->type ?? '',
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
            'star_student' => $starStudent ? [
                'name' => $starStudent->name,
                'pet_name' => $starStudent->pet->name ?? '',
                'pet_species' => $starStudent->pet->type ?? '',
                'pet_level' => $starStudent->pet->level ?? 0,
                'score' => $starStudent->total_score,
            ] : null,
            'top5' => $top5,
            'recent_news' => $recentNews,
        ]]);
    }

    /**
     * 教室端 · 学生列表（含宠物信息）
     */
    public function classroomStudents(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $students = Student::where('class_id', $classInfo['class_id'])
            ->where('status', 'active')
            ->with('pet')
            ->get()
            ->map(fn (Student $s) => [
                'id' => $s->id,
                'name' => $s->name,
                'student_no' => $s->student_no,
                'total_score' => $s->total_score,
                'pet_name' => $s->pet->name ?? '',
                'pet_species' => $s->pet->type ?? '',
                'pet_level' => $s->pet->level ?? 0,
                'pet_emoji' => $s->pet ? ($s->pet->currentStage()['emoji'] ?? '🥚') : '🥚',
            ]);

        return response()->json(['data' => $students]);
    }

    /**
     * 教室端 · 加减分
     */
    public function classroomGiveScore(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $request->validate([
            'student_id' => 'required|integer',
            'points' => 'required|integer|not_in:0',
            'reason' => 'nullable|string|max:200',
        ]);

        $classId = (int) $classInfo['class_id'];
        $student = Student::where('class_id', $classId)->findOrFail((int) $request->input('student_id'));
        $amount = (int) $request->input('points');
        $reason = $request->input('reason', '课堂评价');

        try {
// 班级码模式限制：单次加减分不得超过 ±30
        if (abs($amount) >= 30) {
            return response()->json(['message' => '单次加减分超过 30 分，请使用教师账号登录操作'], 403);
        }

        $student->total_score = max(0, $student->total_score + $amount);
        $student->save();

        if ($student->pet) {
            if ($amount > 0) {
                $student->pet->addExperience($amount);
            } else {
                $student->pet->removeExperience(abs($amount));
            }
        }

        $teacherId = $this->getClassTeacherId($classId);
        if (!$teacherId) {
            $teacherId = \App\Models\User::whereIn('role', ['school_admin', 'admin'])->value('id') ?? 1;
        }
        \App\Models\Score::create([
            'student_id' => $student->id,
            'class_id' => $classId,
            'amount' => $amount,
            'reason' => $reason,
            'given_by' => $teacherId,
        ]);
        } catch (Throwable $e) {
            Log::error("classroomGiveScore failed: " . $e->getMessage());
            return response()->json(["message" => "操作失败: " . $e->getMessage()], 500);
        }

        return response()->json([
            'message' => ($amount > 0 ? '加' : '减') . '分成功',
            'data' => [
                'student_name' => $student->name,
                'points' => $amount,
                'new_score' => $student->fresh()->total_score,
            ],
        ]);
    }

    /**
     * 教室端 · 宠物概览
     */
    public function classroomPetsOverview(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $pets = Pet::where('class_id', $classInfo['class_id'])
            ->with('student:id,name')
            ->get()
            ->map(fn (Pet $p) => [
                'id' => $p->id,
                'student_id' => $p->student_id,
                'student_name' => $p->student->name,
                'name' => $p->name,
                'species' => $p->species,
                'level' => $p->level,
                'exp' => $p->exp,
                'mood' => $p->mood,
                'stage_name' => $p->currentStage()['name'] ?? '',
                'emoji' => $p->currentStage()['emoji'] ?? '🐾',
            ]);

        return response()->json(['data' => $pets]);
    }

    /**
     * 教室端 · PK 排行榜（同年级各班）
     */
    public function classroomPKLeaderboard(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $classId = (int) $classInfo['class_id'];
        $myClass = ClassRoom::find($classId);
        if (!$myClass || !$myClass->grade) {
            return response()->json(['data' => []]);
        }

        $gradeClasses = ClassRoom::where('grade', $myClass->grade)
            ->where('status', 'active')->get();

        if ($gradeClasses->isEmpty()) {
            return response()->json(['data' => []]);
        }

        $pkData = $gradeClasses->map(function ($class) use ($myClass) {
            $students = Student::where('class_id', $class->id)
                ->where('status', 'active')->with('pet')->get();
            $totalScore = $students->sum('total_score');
            $count = $students->count();
            $avgLevel = $count > 0 ? $students->avg(fn ($s) => $s->pet->level ?? 0) : 0;
            $peakCount = $students->filter(fn ($s) => $s->pet && $s->pet->level >= 8)->count();
            $weekStart = now()->startOfWeek();
            $weeklyScore = \App\Models\Score::whereIn('student_id', $students->pluck('id'))
                ->where('created_at', '>=', $weekStart)->sum('amount');

            return [
                'name' => $class->name,
                'totalScore' => (int) $totalScore,
                'studentCount' => $count,
                'avgLevel' => round($avgLevel, 1),
                'peakCount' => $peakCount,
                'weekGrowth' => (int) $weeklyScore,
                'isOwn' => $class->id === $myClass->id,
            ];
        })->sortByDesc('totalScore')->values();

        return response()->json(['data' => $pkData]);
    }

    /**
     * 教室端 · 切换宠物系列（每人扣除20积分）
     */
    public function classroomSwitchSeries(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $request->validate(['series_id' => 'required|string|max:50']);
        $seriesId = $request->input('series_id');
        $validSeries = ['myth', 'pokemon', 'national', 'mecha', 'magic', 'prehistoric', 'constellation', 'folklore'];

        if (!in_array($seriesId, $validSeries, true)) {
            return response()->json(['message' => '无效的系列ID'], 422);
        }

        $classId = (int) $classInfo['class_id'];
        $class = ClassRoom::findOrFail($classId);

        // 计算扣除总额：每人20积分
        $costPerStudent = 20;
        $students = Student::where('class_id', $classId)->where('status', 'active')->get();
        $activeCount = $students->count();

        if ($activeCount === 0) {
            return response()->json(['message' => '班级没有活跃学生'], 400);
        }

        // 检查每个学生积分是否足够
        $insufficient = $students->filter(fn ($s) => $s->total_score < $costPerStudent);
        if ($insufficient->isNotEmpty()) {
            $names = $insufficient->take(3)->pluck('name')->implode('、');
            $more = $insufficient->count() > 3 ? '等' . $insufficient->count() . '人' : '';

            return response()->json([
                'message' => "积分不足：{$names}{$more} 每人需要 {$costPerStudent} 积分",
            ], 400);
        }

        // 扣除每个学生20积分
        foreach ($students as $student) {
            $student->total_score -= $costPerStudent;
            $student->save();
        }

        // 更新班级配置
        $settings = $class->settings ?? [];
        $settings['pet_series'] = $seriesId;
        $class->settings = $settings;
        $class->save();

        return response()->json([
            'message' => "已切换至「{$seriesId}」系列，全班 {$activeCount} 人各扣除 {$costPerStudent} 积分",
            'data' => [
                'series_id' => $seriesId,
                'cost_per_student' => $costPerStudent,
                'affected_students' => $activeCount,
            ],
        ]);
    }

    /**
     * 教室端 · 学生自行切换宠物（首次免费，后续扣20积分）
     */
    public function classroomSwitchPet(Request $request): JsonResponse
    {
        $classInfo = $this->validateToken($request);
        if (!$classInfo) {
            return response()->json(['message' => 'Token无效或已过期'], 401);
        }

        $request->validate([
            'student_id' => 'required|integer',
            'pet_species' => 'required|string|max:50',
        ]);

        $classId = (int) $classInfo['class_id'];
        $student = Student::where('class_id', $classId)->findOrFail((int) $request->input('student_id'));
        $newSpecies = $request->input('pet_species');

        $pet = $student->pet;

        if ($pet) {
            // 检查是否首次切换
            $switchCount = (int) Cache::get("pet_switch_count:{$pet->id}", 0);
            $cost = 0;

            if ($switchCount > 0) {
                // 后续切换扣20积分
                $cost = 20;
                if ($student->total_score < $cost) {
                    return response()->json(['message' => "积分不足，切换需要 {$cost} 积分"], 400);
                }
                $student->total_score -= $cost;
                $student->save();
            }

            // 保留等级和经验，只换种类
            $pet->type = $newSpecies;
            $pet->save();

            Cache::put("pet_switch_count:{$pet->id}", $switchCount + 1, now()->addYears(1));

            $stage = $pet->currentStage();

            return response()->json([
                'message' => $switchCount === 0 ? '✅ 首次切换免费！' : "✅ 已切换，扣除 {$cost} 积分",
                'data' => [
                    'pet_emoji' => $stage['emoji'] ?? '🐾',
                    'pet_name' => $pet->name,
                    'pet_species' => $newSpecies,
                    'total_score' => $student->fresh()->total_score,
                    'switch_count' => $switchCount + 1,
                    'cost' => $switchCount > 0 ? $cost : 0,
                ],
            ]);
        }

        // 无宠物：创建新宠物
        $pet = Pet::create([
            'student_id' => $student->id,
            'class_id' => $classId,
            'name' => $student->name . '的伙伴',
            'type' => $newSpecies,
            'level' => 0,
            'experience' => 0,
            'mood' => 80,
        ]);

        $stage = $pet->currentStage();

        return response()->json([
            'message' => '🎉 新宠物已诞生！',
            'data' => [
                'pet_emoji' => $stage['emoji'] ?? '🐾',
                'pet_name' => $pet->name,
                'pet_species' => $newSpecies,
                'total_score' => $student->total_score,
                'switch_count' => 0,
                'cost' => 0,
            ],
        ]);
    }

    // ============================================================
    // 辅助
    // ============================================================

    private function getClassTeacherId(int $classId): ?int
    {
        $id = ClassRoom::where('id', $classId)->value('teacher_id');
        if (!$id) {
            $id = ClassRoomTeacher::where('class_room_id', $classId)->value('user_id');
        }
        if (!$id) {
            $id = User::where('role', 'teacher')->value('id');
        }

        return $id;
    }
}
