<?php

declare(strict_types=1);

namespace App\Services;

use Illuminate\Support\Facades\Cache;

/**
 * DisplayEventService — 班级大屏事件发布/消费
 *
 * 架构: 教师操作 → 写入 DB → 发布事件到 Cache → SSE 长连接轮询消费
 * 存储: Redis List (或 File Cache 降级)
 * Key: display:events:{class_id}
 * TTL: 10 分钟（事件超过 10 分钟自动清理）
 *
 * 无需 Pusher/Reverb/Node.js，兼容所有部署环境。
 */
class DisplayEventService
{
    private const CACHE_PREFIX = 'display:events:';
    private const CACHE_TTL = 600; // 10 分钟
    private const MAX_EVENTS = 200; // 每个班级最多保留事件数

    /**
     * 发布事件到班级频道
     *
     * @param int    $classId 班级 ID
     * @param string $type    事件类型: score_update|broadcast|notice|pet_update|refresh
     * @param array  $data    事件载荷
     */
    public function publish(int $classId, string $type, array $data): void
    {
        $key = self::CACHE_PREFIX . $classId;

        $event = [
            'id' => $this->nextEventId($classId),
            'type' => $type,
            'data' => $data,
            'created_at' => now()->toIso8601String(),
        ];

        // 获取现有事件列表，追加新事件
        $events = Cache::get($key, []);
        $events[] = $event;

        // 只保留最近 N 条
        if (count($events) > self::MAX_EVENTS) {
            $events = array_slice($events, -self::MAX_EVENTS);
        }

        Cache::put($key, $events, now()->addSeconds(self::CACHE_TTL));
    }

    /**
     * 消费自 $sinceId 之后的新事件
     *
     * @param int $classId
     * @param int|null $sinceId 上次收到的最新事件 ID，null 则返回全部
     * @return array 新事件列表
     */
    public function consume(int $classId, ?int $sinceId = null): array
    {
        $key = self::CACHE_PREFIX . $classId;
        $events = Cache::get($key, []);

        if ($sinceId === null) {
            return $events;
        }

        return array_values(
            array_filter($events, fn (array $e) => ($e['id'] ?? 0) > $sinceId)
        );
    }

    /**
     * 清除班级所有事件
     */
    public function clear(int $classId): void
    {
        Cache::forget(self::CACHE_PREFIX . $classId);
    }

    /**
     * 获取当前最大事件 ID
     */
    private function nextEventId(int $classId): int
    {
        $key = self::CACHE_PREFIX . $classId . ':counter';
        $id = Cache::get($key, 0);
        $id++;
        Cache::put($key, $id, now()->addSeconds(self::CACHE_TTL));

        return $id;
    }

    /**
     * 获取班级码缓存键
     */
    public static function codeCacheKey(string $code): string
    {
        return 'display:code:' . $code;
    }
}
