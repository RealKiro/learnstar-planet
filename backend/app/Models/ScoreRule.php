<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

// 积分规则库（教师可自定义）
class ScoreRule extends Model
{
    protected $fillable = [
        'class_id',
        'school_id',
        'name',             // 如：作业完成、课堂发言
        'amount',           // 积分值
        'category',         // classroom / homework / behavior / literacy / daily / academic / custom
        'is_positive',      // true=加分 false=减分
        'is_active',
        'sort_order',
    ];

    protected $casts = [
        'amount' => 'integer',
        'is_positive' => 'boolean',
        'is_active' => 'boolean',
        'sort_order' => 'integer',
    ];
    // ========== 分类标签与默认规则 ==========
    // 唯一真源在 App\Services\ScoreRuleService（CATEGORY_LABELS / DEFAULT_RULES），
    // 本模型不再自带第二份分类表与默认规则表（历史遗留的两份已删除，避免三套词汇并存）。

    public function classRoom(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function school()
    {
        return $this->belongsTo(School::class);
    }

    /**
     * 获取所有活跃的规则（用于 Score::getRules() 调用）
     *
     * @return \Illuminate\Database\Eloquent\Collection<int, self>
     */
    public static function getActiveRules(): \Illuminate\Database\Eloquent\Collection
    {
        return self::where('is_active', true)
            ->orderBy('sort_order')
            ->get();
    }
}
