<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Broadcast;
use App\Models\ClassRoom;
use App\Models\ClassRoomTeacher;
use App\Models\Notice;
use App\Models\Student;
use App\Models\User;
use Carbon\Carbon;

/**
 * 教师端模式切换与班级大屏消息服务（广播 + 通知统一收口）。
 *
 * 权限失败以 \DomainException 携带 HTTP 状态码抛出，由控制器原样映射。
 */
class ClassroomMessagingService
{
    /**
     * 当前模式与激活班级。
     *
     * @return array{mode: mixed, active_class_id: mixed}
     */
    public function getMode(User $teacher): array
    {
        return [
            'mode' => $teacher->getSetting('display_mode', 'classroom_display'),
            'active_class_id' => $teacher->getSetting('active_class_id', null),
        ];
    }

    /**
     * 切换模式；附带 class_id 时校验分配关系后同步切换激活班级。
     *
     * @return array{message: string, data: array{mode: string, active_class_id: mixed}}
     */
    public function setMode(User $teacher, string $mode, mixed $classId): array
    {
        $teacher->setSetting('display_mode', $mode);

        if ($classId) {
            $isAssigned = $teacher->isApiBot()
                ? ClassRoom::where('school_id', $teacher->school_id)->where('id', (int) $classId)->where('status', 'active')->exists()
                : ClassRoomTeacher::where('user_id', $teacher->id)
                    ->where('class_room_id', (int) $classId)
                    ->exists();
            if ($isAssigned) {
                $teacher->setSetting('active_class_id', (int) $classId);
            }
        }

        return [
            'message' => '已切换为' . ($mode === 'classroom_display' ? '班级大屏' : '教师管理') . '模式',
            'data' => [
                'mode' => $mode,
                'active_class_id' => $teacher->getSetting('active_class_id'),
            ],
        ];
    }

    /**
     * 班级大屏聚合数据（学生宠物总览 + 广播 + 通知 + 最近积分）。
     *
     * @return array<string, mixed>
     * @throws \DomainException 400 未选班级 / 403 未分配该班级
     */
    public function display(User $teacher, mixed $classId): array
    {
        if (!$classId) {
            throw new \DomainException('请先选择班级', 400);
        }

        $isAssigned = $teacher->isApiBot()
            ? ClassRoom::where('school_id', $teacher->school_id)->where('id', (int) $classId)->where('status', 'active')->exists()
            : ClassRoomTeacher::where('user_id', $teacher->id)
                ->where('class_room_id', $classId)
                ->exists();
        if (!$isAssigned) {
            throw new \DomainException('您未被分配到此班级', 403);
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
            $stage = $pet ? $pet->currentStage() : ['emoji' => '🤔', 'name' => '未孵化', 'title' => ''];

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
        $recentScores = \App\Models\Score::where('class_id', $classId)
            ->where('created_at', '>=', now()->subHours(4))
            ->with('student:id,name')
            ->orderBy('created_at', 'desc')
            ->take(20)
            ->get()
            ->map(fn (\App\Models\Score $s) => [
                'student_name' => $s->student?->name,
                'amount' => $s->amount,
                'reason' => $s->reason,
                'time' => $s->created_at?->diffForHumans(),
            ]);

        return [
            'class_name' => $classRoom->name,
            'grade' => $classRoom->grade,
            'student_count' => $students->count(),
            'pets' => $pets,
            'broadcasts' => $broadcasts,
            'notices' => $notices,
            'recent_scores' => $recentScores,
        ];
    }

    /**
     * 发送班级消息：banner/popup/fullscreen 入广播表，info/homework/event/urgent 入通知表，
     * 均推送班级大屏。
     *
     * @param array<int, int> $accessibleIds
     * @param array{content: string, title?: ?string, display_seconds?: mixed, voice?: mixed} $input
     * @return array{message: string, data: array{id: int, type: string}}
     * @throws \DomainException 403 无权限
     */
    public function send(User $teacher, array $accessibleIds, int $classId, string $type, array $input): array
    {
        if (!in_array($classId, $accessibleIds)) {
            throw new \DomainException('无权限', 403);
        }

        $content = $input['content'];

        // Broadcast types: banner, popup, fullscreen
        if (in_array($type, ['banner', 'popup', 'fullscreen'])) {
            $broadcast = Broadcast::create([
                'school_id' => $teacher->school_id,
                'class_id' => $classId,
                'teacher_id' => $teacher->id,
                'content' => $content,
                'type' => $type,
                'voice_enabled' => (bool) ($input['voice'] ?? true),
                'display_seconds' => (int) ($input['display_seconds'] ?? 10),
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

            return [
                'message' => '广播已发送',
                'data' => ['id' => $broadcast->id, 'type' => 'broadcast'],
            ];
        }

        // Notice types: info, homework, event, urgent
        $notice = Notice::create([
            'class_id' => $classId,
            'school_id' => $teacher->school_id,
            'title' => $input['title'] ?? ($type === 'urgent' ? '紧急通知' : '通知'),
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

        return [
            'message' => '通知已发布',
            'data' => ['id' => $notice->id, 'type' => 'notice'],
        ];
    }

    /**
     * 大屏轮询增量消息（since 之后的广播与通知）。
     *
     * @return array{broadcasts: mixed, notices: mixed, polled_at: string}
     * @throws \DomainException 400 未选班级
     */
    public function poll(User $teacher, mixed $classId, mixed $since): array
    {
        if (!$classId) {
            throw new \DomainException('请先选择班级', 400);
        }

        $sinceTime = $since ? Carbon::parse($since) : now()->subMinutes(5);

        $broadcasts = Broadcast::where('class_id', $classId)
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

        return [
            'broadcasts' => $broadcasts,
            'notices' => $notices,
            'polled_at' => now()->toIso8601String(),
        ];
    }
}
