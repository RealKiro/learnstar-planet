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
     *
     * 智能默认：
     * - username: 默认为 name（去重后加 _2/_3 后缀）
     * - nickname: 默认为 name 的拼音（去重后加 _2/_3 后缀）
     * - 第三方登录首次创建账号时,昵称/头像自动采用第三方平台的值（由调用方传入）
     */
    public function createTeacherAccounts(School $school, array $teachers): array
    {
        $created = [];

        foreach ($teachers as $teacher) {
            $name = $teacher['name'];

            // username 默认 = name
            $username = $teacher['username'] ?? $this->uniqueUsername($name, $school);
            // nickname 默认 = name 的拼音
            $nickname = $teacher['nickname'] ?? $this->uniqueNickname($name, $school);
            // 头像: 调用方若传入则使用(第三方首次登录场景),否则留空
            $avatar = $teacher['avatar_path'] ?? null;

            $initialPassword = $teacher['password'] ?? str()->random(10);

            $user = User::create([
                'school_id' => $school->id,
                'role' => 'teacher',
                'username' => $username,
                'password' => Hash::make($initialPassword),
                'name' => $name,
                'nickname' => $nickname,
                'avatar_path' => $avatar,
                'phone' => $teacher['phone'] ?? null,
                'email' => $teacher['email'] ?? null,
                'password_changed' => false,
                'status' => 'active',
            ]);

            $created[] = [
                'id' => $user->id,
                'username' => $username,
                'nickname' => $nickname,
                'initial_password' => $initialPassword, // 仅创建时返回，后续不可查
                'name' => $user->name,
            ];
        }

        return $created;
    }

    /**
     * 生成全校唯一的 username
     * 若 name 已被占用，则依次追加 _2/_3/_4...
     */
    public function uniqueUsername(string $base, School $school): string
    {
        return $this->makeUnique($base, function (string $candidate) use ($school) {
            return User::where('school_id', $school->id)
                ->where('username', $candidate)
                ->exists();
        });
    }

    /**
     * 生成全校唯一的 nickname（基于拼音）
     */
    public function uniqueNickname(string $chineseName, School $school): string
    {
        $base = PinyinService::toPinyin($chineseName);
        if ($base === '') {
            $base = $chineseName; // 拼音库未装时 fallback
        }

        return $this->makeUnique($base, function (string $candidate) use ($school) {
            return User::where('school_id', $school->id)
                ->where('nickname', $candidate)
                ->exists();
        });
    }

    /**
     * 通用去重：$exists($candidate) 为 true 时，追加 _2/_3/...
     */
    private function makeUnique(string $base, callable $exists): string
    {
        $candidate = $base;
        $i = 2;
        while ($exists($candidate)) {
            $candidate = $base . '_' . $i;
            $i++;
            if ($i > 999) {
                // 极端兜底：超过 999 次重名时使用随机后缀
                $candidate = $base . '_' . Str::lower(Str::random(4));
                break;
            }
        }

        return $candidate;
    }

    /**
     * 学校管理员创建家长账号（绑定学生）
     *
     * username 默认 = 家长姓名
     * nickname 默认 = 姓名拼音
     *
     * 返回数组（含明文初始密码），方便管理员一次性下发给家长
     */
    public function createParentAccount(School $school, array $parentData): array
    {
        $name = $parentData['name'];
        $username = $parentData['username'] ?? $this->uniqueUsername($name, $school);
        $nickname = $parentData['nickname'] ?? $this->uniqueNickname($name, $school);
        $initialPassword = $parentData['password'] ?? str()->random(10);

        $parent = User::create([
            'school_id' => $school->id,
            'role' => 'parent',
            'username' => $username,
            'password' => Hash::make($initialPassword),
            'name' => $name,
            'nickname' => $nickname,
            'avatar_path' => $parentData['avatar_path'] ?? null,
            'phone' => $parentData['phone'] ?? null,
            'password_changed' => false,
            'status' => 'active',
        ]);

        return [
            'id' => $parent->id,
            'username' => $parent->username,
            'nickname' => $parent->nickname,
            'initial_password' => $initialPassword,
            'name' => $parent->name,
        ];
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
     * nick/avatar 暂存到 temp_token 关联的缓存,绑定时取出同步到 user
     */
    public function loginWithWechat(string $openid, ?string $unionid = null, ?string $nick = null, ?string $avatar = null): array
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
        $tempToken = Str::uuid()->toString();
        $this->storeTempBindingContext($tempToken, [
            'platform' => 'wechat',
            'platform_id' => $openid,
            'unionid' => $unionid,
            'nick' => $nick,
            'avatar' => $avatar,
        ]);

        return [
            'status' => 'need_binding',
            'temp_token' => $tempToken,
            'openid' => $openid,
            'unionid' => $unionid,
        ];
    }

    /**
     * 企业微信扫码登录
     */
    public function loginWithWechatWork(string $userid, ?string $nick = null, ?string $avatar = null): array
    {
        $user = ThirdPartyBinding::findUserByPlatform('wechat_work', $userid);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        $tempToken = Str::uuid()->toString();
        $this->storeTempBindingContext($tempToken, [
            'platform' => 'wechat_work',
            'platform_id' => $userid,
            'nick' => $nick,
            'avatar' => $avatar,
        ]);

        return [
            'status' => 'need_binding',
            'temp_token' => $tempToken,
            'platform_id' => $userid,
        ];
    }

    /**
     * QQ扫码登录
     */
    public function loginWithQQ(string $openid, ?string $nick = null, ?string $avatar = null): array
    {
        $user = ThirdPartyBinding::findUserByPlatform('qq', $openid);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        $tempToken = Str::uuid()->toString();
        $this->storeTempBindingContext($tempToken, [
            'platform' => 'qq',
            'platform_id' => $openid,
            'nick' => $nick,
            'avatar' => $avatar,
        ]);

        return [
            'status' => 'need_binding',
            'temp_token' => $tempToken,
            'openid' => $openid,
        ];
    }

    /**
     * 人人通空间登录
     */
    public function loginWithRenren(string $userId, ?string $nick = null, ?string $avatar = null): array
    {
        $user = ThirdPartyBinding::findUserByPlatform('renren', $userId);

        if ($user) {
            $user->update(['last_login_at' => now()]);

            return ['status' => 'logged_in', 'user' => $user];
        }

        $tempToken = Str::uuid()->toString();
        $this->storeTempBindingContext($tempToken, [
            'platform' => 'renren',
            'platform_id' => $userId,
            'nick' => $nick,
            'avatar' => $avatar,
        ]);

        return [
            'status' => 'need_binding',
            'temp_token' => $tempToken,
            'platform_id' => $userId,
        ];
    }

    /**
     * 把扫码时的 nick/avatar 暂存到 cache（10 分钟），绑定时再取出同步
     * 避免前端在 need_binding 后再传一次
     */
    private function storeTempBindingContext(string $tempToken, array $context): void
    {
        \Illuminate\Support\Facades\Cache::put('wechat_scan_ctx:' . $tempToken, $context, now()->addMinutes(10));
    }

    private function getTempBindingContext(string $tempToken): ?array
    {
        $ctx = \Illuminate\Support\Facades\Cache::get('wechat_scan_ctx:' . $tempToken);

        return is_array($ctx) ? $ctx : null;
    }

    // ========== 绑定第三方账号 ==========

    /**
     * 用户主动绑定第三方平台（需先登录账号密码）
     * @phpstan-return ThirdPartyBinding
     */
    public function bindThirdParty(User $user, string $platform, string $platformId, ?string $unionId = null, ?string $nick = null, ?string $avatar = null): ThirdPartyBinding
    {
        // 检查是否已绑定同一平台
        /** @var ThirdPartyBinding|null $existing */
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

        /** @var ThirdPartyBinding $binding */
        $binding = $user->thirdPartyBindings()->create([
            'platform' => $platform,
            'platform_id' => $platformId,
            'platform_union_id' => $unionId,
            'platform_nick' => $nick,
            'platform_avatar' => $avatar,
            'verified_at' => now(),
        ]);

        return $binding;
    }

    /**
     * 第三方首次扫码后绑定已有教师账号
     * 仅允许绑定教师账号（role=teacher），管理员不支持第三方登录
     *
     * 第三方首次登录时,优先从 temp_token 对应的 cache 取出扫码时的 nick/avatar
     * 同步到 user 表（仅在 user 尚未设置时覆盖）
     */
    public function bindAfterScan(string $tempToken, string $username, string $password, string $platform, string $platformId, ?string $unionId = null, ?string $nick = null, ?string $avatar = null): array
    {
        $user = $this->teacherLoginWithCredentials($username, $password);

        if (!$user) {
            return ['status' => 'error', 'message' => '账号或密码错误，请核对后重试'];
        }

        // 优先从 cache 取（扫码时存的 nick/avatar）
        $ctx = $this->getTempBindingContext($tempToken) ?? [];
        $ctxPlatform = $ctx['platform'] ?? $platform;
        $ctxPlatformId = $ctx['platform_id'] ?? $platformId;
        $ctxUnionId = $ctx['unionid'] ?? $unionId;
        $ctxNick = $ctx['nick'] ?? $nick;
        $ctxAvatar = $ctx['avatar'] ?? $avatar;

        $this->bindThirdParty($user, $ctxPlatform, $ctxPlatformId, $ctxUnionId, $ctxNick, $ctxAvatar);

        // 同步昵称/头像到 user 表（仅在 user 尚未设置时覆盖）
        $updates = [];
        if ($ctxNick && (empty($user->nickname) || $user->nickname === PinyinService::toPinyin($user->name))) {
            $updates['nickname'] = $ctxNick;
        }
        if ($ctxAvatar && empty($user->avatar_path)) {
            $updates['avatar_path'] = $ctxAvatar;
        }
        if ($updates) {
            $user->update($updates);
        }

        // 清理临时 context
        \Illuminate\Support\Facades\Cache::forget('wechat_scan_ctx:' . $tempToken);

        return ['status' => 'bound', 'user' => $user->fresh()];
    }

    /**
     * 解绑第三方平台
     */
    public function unbindThirdParty(User $user, string $platform): bool
    {
        return $user->thirdPartyBindings()->where('platform', $platform)->delete() > 0;
    }
}
