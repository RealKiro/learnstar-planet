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
        'nickname',          // 昵称，默认 = name 的拼音
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
        'role' => 'string',
        'status' => 'string',
    ];

    // ========== 查询作用域 ==========

    /** @param \Illuminate\Database\Eloquent\Builder $query */
    public function scopeActive($query)
    {
        return $query->where('status', 'active');
    }

    /** @param \Illuminate\Database\Eloquent\Builder $query */
    public function scopeByRole($query, string $role)
    {
        return $query->where('role', $role);
    }

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

    public function classRoomAssignments(): HasMany
    {
        return $this->hasMany(ClassRoomTeacher::class);
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
        /** @var ThirdPartyBinding $binding */
        $binding = $this->thirdPartyBindings()->create([
            'platform' => 'wechat',
            'platform_id' => $openid,
            'platform_union_id' => $unionid,
            'platform_nick' => '',
            'platform_avatar' => '',
            'verified_at' => now(),
        ]);

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindWechatWork(string $userid): ThirdPartyBinding
    {
        /** @var ThirdPartyBinding $binding */
        $binding = $this->thirdPartyBindings()->create([
            'platform' => 'wechat_work',
            'platform_id' => $userid,
            'verified_at' => now(),
        ]);

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindQQ(string $openid): ThirdPartyBinding
    {
        /** @var ThirdPartyBinding $binding */
        $binding = $this->thirdPartyBindings()->create([
            'platform' => 'qq',
            'platform_id' => $openid,
            'verified_at' => now(),
        ]);

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding
     */
    public function bindRenren(string $userId): ThirdPartyBinding
    {
        /** @var ThirdPartyBinding $binding */
        $binding = $this->thirdPartyBindings()->create([
            'platform' => 'renren',
            'platform_id' => $userId,
            'verified_at' => now(),
        ]);

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getWechatBinding(): ?ThirdPartyBinding
    {
        /** @var ThirdPartyBinding|null $binding */
        $binding = $this->thirdPartyBindings()->where('platform', 'wechat')->first();

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getWechatWorkBinding(): ?ThirdPartyBinding
    {
        /** @var ThirdPartyBinding|null $binding */
        $binding = $this->thirdPartyBindings()->where('platform', 'wechat_work')->first();

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getQQBinding(): ?ThirdPartyBinding
    {
        /** @var ThirdPartyBinding|null $binding */
        $binding = $this->thirdPartyBindings()->where('platform', 'qq')->first();

        return $binding;
    }

    /**
     * @phpstan-return ThirdPartyBinding|null
     */
    public function getRenrenBinding(): ?ThirdPartyBinding
    {
        /** @var ThirdPartyBinding|null $binding */
        $binding = $this->thirdPartyBindings()->where('platform', 'renren')->first();

        return $binding;
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
