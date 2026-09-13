<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Broadcast;
use App\Models\User;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Support\Facades\Log;

/**
 * 班级广播（大屏横幅/弹窗/全屏）读写与实时推送。
 *
 * 只负责数据操作与事件推送；HTTP 校验与响应包装留在控制器。
 */
class BroadcastService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
        private readonly DisplayEventService $displayEvents,
    ) {
    }

    /** 最近 20 条广播 */
    public function recent(User $teacher): Collection
    {
        return Broadcast::whereIn('class_id', $this->scope->ids($teacher))
            ->orderBy('created_at', 'desc')->take(20)->get();
    }

    /** 按作用域查找广播（越权即 404） */
    public function findInScope(User $teacher, int $id): Broadcast
    {
        return Broadcast::whereIn('class_id', $this->scope->ids($teacher))->findOrFail($id);
    }

    /**
     * 向目标班级发送广播并推送到大屏。
     *
     * 目标班级 = 传入的 class_ids（若有）∩ 教师可访问班级；
     * 未指定 class_ids 时发送给全部可访问班级。
     *
     * @param  array<int, int>|null  $classIds
     * @return int  实际发送的班级数；返回 0 表示没有可发送的班级
     */
    public function send(
        User $teacher,
        string $content,
        string $type,
        bool $voice,
        bool $loop,
        int $duration,
        ?array $classIds,
    ): int {
        $accessibleClassIds = $this->scope->accessibleIds($teacher);

        $targetIds = $classIds ?? $accessibleClassIds;
        $targetIds = array_intersect($targetIds, $accessibleClassIds);

        if (empty($targetIds)) {
            return 0;
        }

        $sent = 0;
        foreach ($targetIds as $classId) {
            $broadcast = Broadcast::create([
                'class_id' => $classId,
                'content' => $content,
                'type' => $type,
                'voice_enabled' => $voice,
                'loop_enabled' => $loop,
                'display_seconds' => $duration,
                'status' => 'sent',
                'sent_at' => now(),
            ]);

            // 推送事件到班级大屏
            try {
                $this->displayEvents->publish($classId, 'broadcast', [
                    'id' => $broadcast->id,
                    'type' => $broadcast->type,
                    'content' => $broadcast->content,
                    'display_seconds' => $broadcast->display_seconds,
                    'voice_enabled' => $broadcast->voice_enabled,
                    'created_at' => $broadcast->created_at?->toIso8601String(),
                ]);
            } catch (\Throwable $e) {
                logger()->warning('Broadcast event publish failed for class ' . $classId . ': ' . $e->getMessage());
            }

            // 推送事件到班级大屏
            try {
                $this->displayEvents->publish($classId, 'broadcast', [
                    'id' => $broadcast->id,
                    'type' => $broadcast->type,
                    'content' => $broadcast->content,
                    'display_seconds' => $broadcast->display_seconds,
                    'voice_enabled' => $broadcast->voice_enabled ?? false,
                    'created_at' => $broadcast->created_at?->toIso8601String(),
                ]);
            } catch (\Throwable $e) {
                Log::warning('Broadcast event publish failed: ' . $e->getMessage());
            }

            $sent++;
        }

        return $sent;
    }
}
