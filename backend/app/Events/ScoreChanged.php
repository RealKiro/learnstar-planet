<?php

declare(strict_types=1);

namespace App\Events;

use Illuminate\Broadcasting\Channel;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Contracts\Broadcasting\ShouldBroadcastNow;
use Illuminate\Queue\SerializesModels;

class ScoreChanged implements ShouldBroadcastNow
{
    use InteractsWithSockets;
    use SerializesModels;

    public int $studentId;

    public int $amount;

    public string $reason;

    public function __construct(int $studentId, int $amount, string $reason)
    {
        $this->studentId = $studentId;
        $this->amount = $amount;
        $this->reason = $reason;
    }

    public function broadcastOn(): array
    {
        // 家长端频道：学生ID对应频道
        return [
            new Channel('student.' . $this->studentId),
        ];
    }

    public function broadcastWith(): array
    {
        return [
            'student_id' => $this->studentId,
            'amount' => $this->amount,
            'reason' => $this->reason,
            'timestamp' => now()->toIso8601String(),
        ];
    }
}

