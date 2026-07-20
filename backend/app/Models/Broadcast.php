<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Broadcast extends Model
{
    protected $fillable = [
        'school_id', 'class_id', 'teacher_id',
        'content', 'type', 'voice_enabled', 'loop_enabled', 'display_seconds',
        'status', 'sent_at',
    ];

    protected $casts = [
        'voice_enabled' => 'boolean',
        'loop_enabled' => 'boolean',
        'sent_at' => 'datetime',
    ];

    public function school(): BelongsTo
    {
        return $this->belongsTo(School::class);
    }

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    /**
     * 广播类型
     */
    public static function getTypes(): array
    {
        return [
            'banner'     => '顶部横幅',
            'popup'      => '弹窗提示',
            'fullscreen' => '全屏展示',
        ];
    }

    /**
     * 广播状态
     */
    public static function getStatuses(): array
    {
        return [
            'pending' => '待发送',
            'sent'     => '已发送',
            'received' => '已接收',
            'expired'  => '已过期',
        ];
    }
}

