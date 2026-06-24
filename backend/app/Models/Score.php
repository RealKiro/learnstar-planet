<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * @property int $id
 * @property int $student_id
 * @property int $class_id
 * @property int $score_rule_id
 * @property int $amount
 * @property string|null $reason
 * @property int $given_by
 * @property Student $student
 * @property ClassRoom $classRoom
 * @property ScoreRule $scoreRule
 * @property User $giver
 */
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

    /**
     * 获取所有活跃的积分规则（用于前端展示和快速操作）
     *
     * @return array<string, array{name: string, amount: int, description: string}>
     */
    public static function getRules(): array
    {
        return ScoreRule::getActiveRules()
            ->mapWithKeys(fn ($rule) => [
                $rule->code => [
                    'name' => $rule->name,
                    'amount' => $rule->default_amount,
                    'description' => $rule->description ?? '',
                ],
            ])
            ->toArray();
    }

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
