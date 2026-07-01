<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\SoftDeletes;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Laravel\Sanctum\HasApiTokens;

class User extends Authenticatable
{
    use HasApiTokens;
    use SoftDeletes;

    protected $fillable = [
        'school_id',
        'role',              // school_admin / teacher / parent
        'username',          // 管理员分配的账号
        'password',          // 管理员分配的初始密码
        'name',
        'avatar_path',
        'phone',
        'email',
        'status',            // active / disabled
        'last_login_at',
        'password_changed',  // 是否已修改初始密码
    ];

    protected $hidden = [
        'password',
        'remember_token',
    ];

    protected $casts = [
        'password_changed' => 'boolean',
        'last_login_at' => 'datetime',
    ];

    // ========== 关系 ==========

    public function school(): BelongsTo
    {
        return $this->belongsTo(School::class);
    }

    public function thirdPartyBindings(): HasMany
    {
        return $this->hasMany(ThirdPartyBinding::class);
    }

    public function classesAsTeacher(): HasMany
    {
        return $this->hasMany(ClassRoom::class, 'teacher_id');
    }

    public function children(): HasMany
    {
        return $this->hasMany(Student::class, 'parent_id');
    }

    // ========== 第三方登录快捷方法 ==========

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindWechat(string $openid, string $unionid = null): ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->create([
            'platform' => 'wechat',
            'platform_id' => $openid,
            'platform_union_id' => $unionid,
            'platform_nick' => '',
            'platform_avatar' => '',
            'verified_at' => now(),
        ]);
    }

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindWechatWork(string $userid): ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->create([
            'platform' => 'wechat_work',
            'platform_id' => $userid,
            'verified_at' => now(),
        ]);
    }

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindQQ(string $openid): ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->create([
            'platform' => 'qq',
            'platform_id' => $openid,
            'verified_at' => now(),
        ]);
    }

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindRenren(string $userId): ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->create([
            'platform' => 'renren',
            'platform_id' => $userId,
            'verified_at' => now(),
        ]);
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getWechatBinding(): ?ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->where('platform', 'wechat')->first();
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getWechatWorkBinding(): ?ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->where('platform', 'wechat_work')->first();
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getQQBinding(): ?ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->where('platform', 'qq')->first();
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getRenrenBinding(): ?ThirdPartyBinding
    {
        return $this->thirdPartyBindings()->where('platform', 'renren')->first();
    }

    public function hasAnyBinding(): bool
    {
        return $this->thirdPartyBindings()->exists();
    }

    // ========== 辅助 ==========

    public function isSchoolAdmin(): bool
    {
        return $this->role === 'school_admin';
    }

    public function isTeacher(): bool
    {
        return $this->role === 'teacher';
    }

    public function isParent(): bool
    {
        return $this->role === 'parent';
    }
}
