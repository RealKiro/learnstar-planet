<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Score extends Model
{
    protected $fillable = [
        'student_id',
        'class_id',
        'score_rule_id',    // 关联规则
        'amount',           // 积分变化量（正为加分，负为减分）
        'reason',           // 手动填写的理由
        'given_by',         // 发放人（教师ID）
        'created_at',
    ];

    protected $casts = [
        'amount' => 'integer',
    ];

    public function student()
    {
        return $this->belongsTo(Student::class);
    }

    public function classRoom()
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function scoreRule()
    {
        return $this->belongsTo(ScoreRule::class);
    }

    public function giver()
    {
        return $this->belongsTo(User::class, 'given_by');
    }
}
