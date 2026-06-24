<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\School;
use App\Models\ThirdPartyBinding;
use App\Models\User;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

class AuthService
{
    // ========== 管理员分配账号 ==========

    /**
     * 学校管理员批量创建教师账号
     */
    public function createTeacherAccounts(School $school, array $teachers): array
    {
        $created = [];

        foreach ($teachers as $teacher) {
            $username = $teacher['username'] ?? $school->code . '_T' . str_pad((string) ($school->teachers()->count() + 1), 3, '0', STR_PAD_LEFT);
            $initialPassword = $teacher['password'] ?? Str::random(8);

            $user = User::create([
                'school_id' => $school->id,
                'role' => 'teacher',
                'username' => $username,
                'password' => Hash::make($initialPassword),
                'name' => $teacher['name'],
                'phone' => $teacher['phone'] ?? null,
                'email' => $teacher['email'] ?? null,
                'password_changed' => false,
                'status' => 'active',
            ]);

            $created[] = [
                'id' => $user->id,
                'username' => $username,
                'initial_password' => $initialPassword, // 仅创建时返回，后续不可查
                'name' => $user->name,
            ];
        }

        return $created;
    }

    /**
     * 学校管理员创建家长账号（绑定学生）
     */
    public function createParentAccount(School $school, array $parentData): User
    {
        $username = $parentData['username'] ?? $school->code . '_P' . str_pad((string) User::where('school_id', $school->id)->where('role', 'parent')->count() + 1, 3, '0', STR_PAD_LEFT);
        $initialPassword = $parentData['password'] ?? Str::random(8);

        return User::create([
            'school_id' => $school->id,
            'role' => 'parent',
            'username' => $username,
            'password' => Hash::make($initialPassword),
            'name' => $parentData['name'],
            'phone' => $parentData['phone'] ?? null,
            'password_changed' => false,
            'status' => 'active',
        ]);
    }

    // ========== 账号密码登录（教师和管理员分开） ==========

    /**
     * 教师账号密码登录（默认入口）
     * 仅允许 role=teacher 的账号通过此接口登录
     */
    public function teacherLoginWithCredentials(string $username, string $password): ?User
    {
        $user = User::where('username', $username)
            ->where('role', 'teacher')
            ->where('status', 'active')
            ->first();

        if (!$user || !Hash::check($password, $user->password)) {
            return null;
        }

        $user->update(['last_login_at' => now()]);

        return $user;
    }

    /**
     * 管理员账号密码登录（独立入口）
     * 仅允许 role=school_admin 的账号通过此接口登录
     * 管理员不支持第三方扫码登录，仅限账号密码方式
     */
    public function adminLoginWithCredentials(string $username, string $password): ?User
    {
        $user = User::where('username', $username)
            ->where('role', 'school_admin')
            ->where('status', 'active')
            ->first();

        if (!$user || !Hash::check($password, $user->password)) {
            return null;
        }

        $user->update(['last_login_at' => now()]);

        return $user;
    }

    // ========== 第三方扫码登录 ==========

    /**
     * 微信扫码登录
     * 流程：前端获取openid → 调此服务 → 如果已绑定直接登录，否则返回需绑定提示
     */
    public function loginWithWechat(string $openid, ?string $unionid = null): array
    {
        // 优先用unionid查找（可跨小程序和开放平台）
        if ($unionid) {
            $user = ThirdPartyBinding::findUserByUnionId($unionid);

            if ($user) {
                $user->update(['last_login_at' => now()]);

                return ['status' => 'logged_in', 'user' => $user];
            }
        }

        // 用openid查找
        $user = ThirdPartyBinding::findUserByPlatform('wechat', $openid);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        // 未绑定，返回临时token用于后续绑定流程
        return [
            'status' => 'need_binding',
            'temp_token' => Str::uuid()->toString(),
            'openid' => $openid,
            'unionid' => $unionid,
        ];
    }

    /**
     * 企业微信扫码登录
     */
    public function loginWithWechatWork(string $userid): array
    {
        $user = ThirdPartyBinding::findUserByPlatform('wechat_work', $userid);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        return [
            'status' => 'need_binding',
            'temp_token' => Str::uuid()->toString(),
            'platform_id' => $userid,
        ];
    }

    /**
     * QQ扫码登录
     */
    public function loginWithQQ(string $openid): array
    {
        $user = ThirdPartyBinding::findUserByPlatform('qq', $openid);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        return [
            'status' => 'need_binding',
            'temp_token' => Str::uuid()->toString(),
            'openid' => $openid,
        ];
    }

    /**
     * 人人通空间登录
     */
    public function loginWithRenren(string $userId): array
    {
        $user = ThirdPartyBinding::findUserByPlatform('renren', $userId);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        return [
            'status' => 'need_binding',
            'temp_token' => Str::uuid()->toString(),
            'platform_id' => $userId,
        ];
    }

    // ========== 绑定第三方账号 ==========

    /**
     * 用户主动绑定第三方平台（需先登录账号密码）
     */
    public function bindThirdParty(User $user, string $platform, string $platformId, ?string $unionId = null, ?string $nick = null, ?string $avatar = null): ThirdPartyBinding
    {
        // 检查是否已绑定同一平台
        $existing = $user->thirdPartyBindings()->where('platform', $platform)->first();

        if ($existing) {
            // 更新绑定信息
            $existing->update([
                'platform_id' => $platformId,
                'platform_union_id' => $unionId,
                'platform_nick' => $nick,
                'platform_avatar' => $avatar,
                'verified_at' => now(),
            ]);

            return $existing;
        }

        return $user->thirdPartyBindings()->create([
            'platform' => $platform,
            'platform_id' => $platformId,
            'platform_union_id' => $unionId,
            'platform_nick' => $nick,
            'platform_avatar' => $avatar,
            'verified_at' => now(),
        ]);
    }

    /**
     * 第三方首次扫码后绑定已有教师账号
     * 仅允许绑定教师账号（role=teacher），管理员不支持第三方登录
     */
    public function bindAfterScan(string $tempToken, string $username, string $password, string $platform, string $platformId, ?string $unionId = null): array
    {
        $user = $this->teacherLoginWithCredentials($username, $password);

        if (!$user) {
            return ['status' => 'error', 'message' => '教师账号或密码错误，仅支持绑定教师账号'];
        }

        $this->bindThirdParty($user, $platform, $platformId, $unionId);

        return ['status' => 'bound', 'user' => $user];
    }

    /**
     * 解绑第三方平台
     */
    public function unbindThirdParty(User $user, string $platform): bool
    {
        return $user->thirdPartyBindings()->where('platform', $platform)->delete() > 0;
    }
}
