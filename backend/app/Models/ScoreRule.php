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
        'category',         // academic / behavior / activity / custom
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

    // ========== 分类标签 ==========

    public static function categories(): array
    {
        return [
            'academic'  => '学业表现',
            'behavior'  => '行为习惯',
            'activity'  => '活动参与',
            'custom'    => '自定义',
        ];
    }

    // ========== 默认规则模板 ==========

    public static function defaultRules(): array
    {
        return [
            // 加分
            ['name' => '作业完成', 'amount' => 5,  'category' => 'academic', 'is_positive' => true],
            ['name' => '作业优秀', 'amount' => 10, 'category' => 'academic', 'is_positive' => true],
            ['name' => '主动发言', 'amount' => 3,  'category' => 'academic', 'is_positive' => true],
            ['name' => '帮助同学', 'amount' => 5,  'category' => 'behavior', 'is_positive' => true],
            ['name' => '课堂专注', 'amount' => 3,  'category' => 'behavior', 'is_positive' => true],
            ['name' => '活动参与', 'amount' => 5,  'category' => 'activity', 'is_positive' => true],
            ['name' => '考试进步', 'amount' => 8,  'category' => 'academic', 'is_positive' => true],
            ['name' => '考试优秀', 'amount' => 15, 'category' => 'academic', 'is_positive' => true],
            // 减分
            ['name' => '课堂违纪', 'amount' => -3, 'category' => 'behavior', 'is_positive' => false],
            ['name' => '作业未交', 'amount' => -5, 'category' => 'academic', 'is_positive' => false],
            ['name' => '迟到早退', 'amount' => -2, 'category' => 'behavior', 'is_positive' => false],
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

    /**
     * 获取所有活跃的规则（用于 Score::getRules() 调用）
     *
     * @return \Illuminate\Database\Eloquent\Collection<int, static>
     */
    public static function getActiveRules(): \Illuminate\Database\Eloquent\Collection
    {
        return static::where('is_active', true)
            ->orderBy('sort_order')
            ->get();
    }
}
