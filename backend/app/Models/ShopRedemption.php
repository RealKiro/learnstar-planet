<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

// 商城兑换记录
class ShopRedemption extends Model
{
    protected $fillable = [
        'student_id',
        'shop_item_id',
        'class_id',
        'status',         // pending / approved / rejected / delivered
        'approved_by',
        'delivered_at',
    ];

    protected $casts = [
        'delivered_at' => 'datetime',
    ];

    public static function statusLabels(): array
    {
        return [
            'pending'   => '待审批',
            'approved'  => '已批准',
            'rejected'  => '已拒绝',
            'delivered' => '已发放',
        ];
    }

    public function student()
    {
        return $this->belongsTo(Student::class);
    }

    public function shopItem()
    {
        return $this->belongsTo(ShopItem::class);
    }

    public function classRoom()
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function approver()
    {
        return $this->belongsTo(User::class, 'approved_by');
    }
}
