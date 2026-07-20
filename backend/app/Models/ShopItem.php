<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;

// 商城兑换项目（教师自定义奖励）
class ShopItem extends Model
{
    protected $fillable = [
        'class_id',
        'name',
        'description',
        'category',       // physical / privilege / activity / experience
        'cost_score',     // 兑换所需积分
        'currency_type',  // score / science / reading / class_point
        'event_tag',      // 活动标签，用于活动专属商城
        'stock',          // 库存数量，0=不限
        'image_path',
        'is_active',
    ];

    protected $casts = [
        'cost_score' => 'integer',
        'stock' => 'integer',
        'is_active' => 'boolean',
        'currency_type' => 'string',
    ];

    public static function categories(): array
    {
        return [
            'physical'   => '实物奖品',
            'privilege'  => '特权奖励',
            'activity'   => '活动参与',
            'experience' => '体验活动',
        ];
    }

    /**
     * 按币种筛选
     *
     * @param  Builder $query
     */
    public function scopeByCurrency($query, string $currency)
    {
        return $query->where('currency_type', $currency);
    }

    public function classRoom()
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function redemptions()
    {
        return $this->hasMany(ShopRedemption::class);
    }
}
