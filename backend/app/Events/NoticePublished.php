<?php

declare(strict_types=1);

namespace App\Events;

use Illuminate\Broadcasting\Channel;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Contracts\Broadcasting\ShouldBroadcastNow;
use Illuminate\Queue\SerializesModels;

class NoticePublished implements ShouldBroadcastNow
{
    use InteractsWithSockets, SerializesModels;

    public int $classId;

    public int $noticeId;

    public string $title;

    public string $type;

    public function __construct(int $classId, int $noticeId, string $title, string $type)
    {
        $this->classId = $classId;
        $this->noticeId = $noticeId;
        $this->title = $title;
        $this->type = $type;
    }

    public function broadcastOn(): array
    {
        return [
            new Channel('class.' . $this->classId),
        ];
    }

    public function broadcastWith(): array
    {
        return [
            'notice_id' => $this->noticeId,
            'title' => $this->title,
            'type' => $this->type,
            'timestamp' => now()->toIso8601String(),
        ];
    }
}
