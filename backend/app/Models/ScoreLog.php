<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

// 积分变化日志（用于家长端查看明细）
class ScoreLog extends Model
{
    protected $fillable = [
        'student_id',
        'score_id',
        'balance_before',
        'balance_after',
        'description',
    ];

    protected $casts = [
        'balance_before' => 'integer',
        'balance_after' => 'integer',
    ];

    public function student()
    {
        return $this->belongsTo(Student::class);
    }

    public function score()
    {
        return $this->belongsTo(Score::class);
    }
}
