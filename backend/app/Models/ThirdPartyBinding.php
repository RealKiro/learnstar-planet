<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class ThirdPartyBinding extends Model
{
    protected $fillable = [
        'user_id',
        'platform',         // wechat / wechat_work / qq / renren
        'platform_id',      // 该平台的用户唯一ID（openid / userid 等）
        'platform_union_id', // 联合ID（微信unionid，可跨小程序和开放平台）
        'platform_nick',
        'platform_avatar',
        'verified_at',
    ];

    protected $casts = [
        'verified_at' => 'datetime',
    ];

    // ========== 关系 ==========

    public function user()
    {
        return $this->belongsTo(User::class);
    }

    // ========== 查找：根据平台openid反查用户 ==========

    public static function findUserByPlatform(string $platform, string $platformId): ?User
    {
        $binding = static::where('platform', $platform)
            ->where('platform_id', $platformId)
            ->first();

        return $binding?->user;
    }

    public static function findUserByUnionId(string $unionId): ?User
    {
        $binding = static::where('platform_union_id', $unionId)
            ->first();

        return $binding?->user;
    }

    // ========== 平台标识 ==========

    /**
     * 全部支持的第三方平台（含新版办公平台 provider 的 dingtalk / feishu）
     * 账号密码绑定类：wechat / wechat_work / qq / renren
     * 办公平台扫码类：wechat_work / dingtalk / feishu（wechat_work 共用）
     */
    public static function platforms(): array
    {
        return ['wechat', 'wechat_work', 'qq', 'renren', 'dingtalk', 'feishu'];
    }

    public static function platformLabels(): array
    {
        return [
            'wechat'      => '微信',
            'wechat_work' => '企业微信',
            'qq'          => 'QQ',
            'renren'      => '人人通空间',
            'dingtalk'    => '钉钉',
            'feishu'      => '飞书',
        ];
    }

    public static function platformIcons(): array
    {
        return [
            'wechat'      => '💬',
            'wechat_work' => '🏢',
            'qq'          => '🐧',
            'renren'      => '🌐',
            'dingtalk'    => '🔷',
            'feishu'      => '🪶',
        ];
    }
}
