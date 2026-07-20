<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

// 班级通知/公告
class Notice extends Model
{
    protected $fillable = [
        'class_id',
        'school_id',
        'title',
        'content',
        'type',            // info / homework / event / urgent
        'published_by',
        'is_published',
        'published_at',
    ];

    protected $casts = [
        'is_published' => 'boolean',
        'published_at' => 'datetime',
    ];

    public static function typeLabels(): array
    {
        return [
            'info'    => '通知',
            'homework' => '作业',
            'event'   => '活动',
            'urgent'  => '紧急',
        ];
    }

    public function classRoom()
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function school()
    {
        return $this->belongsTo(School::class);
    }

    public function publisher()
    {
        return $this->belongsTo(User::class, 'published_by');
    }
}

