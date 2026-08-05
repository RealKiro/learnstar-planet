<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\School;
use App\Models\ThirdPartyBinding;
use App\Services\AuthService;
use App\Services\WechatWorkService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

class AuthController extends Controller
{
    public function __construct(
        private readonly AuthService $authService
    ) {
    }

    /**
     * 教师账号密码登录
     */
    public function teacherLoginWithCredentials(Request $request): JsonResponse
    {
        $request->validate([
            'username' => 'required|string',
            'password' => 'required|string',
        ]);

        $user = $this->authService->teacherLoginWithCredentials(
            $request->input('username'),
            $request->input('password')
        );

        if (!$user) {
            return response()->json(['message' => '账号或密码错误，请核对后重试'], 401);
        }

        // 生成 Sanctum token
        $token = $user->createToken('auth_token')->plainTextToken;

        return response()->json([
            'data' => [
                'token' => $token,
                'user' => [
                    'id' => $user->id,
                    'username' => $user->username,
                    'name' => $user->name,
                    'role' => $user->role,
                    'school_id' => $user->school_id,
                ],
            ],
        ]);
    }

    /**
     * 管理员账号密码登录
     */
    public function adminLoginWithCredentials(Request $request): JsonResponse
    {
        $request->validate([
            'username' => 'required|string',
            'password' => 'required|string',
        ]);

        $user = $this->authService->adminLoginWithCredentials(
            $request->input('username'),
            $request->input('password')
        );

        if (!$user) {
            return response()->json(['message' => '账号或密码错误，请核对后重试'], 401);
        }

        // 生成 Sanctum token
        $token = $user->createToken('auth_token')->plainTextToken;

        return response()->json([
            'data' => [
                'token' => $token,
                'user' => [
                    'id' => $user->id,
                    'username' => $user->username,
                    'name' => $user->name,
                    'role' => $user->role,
                    'school_id' => $user->school_id,
                ],
            ],
        ]);
    }

    /**
     * 家长账号密码登录
     */
    public function parentLoginWithCredentials(Request $request): JsonResponse
    {
        $request->validate([
            'username' => 'required|string',
            'password' => 'required|string',
        ]);

        $user = $this->authService->parentLoginWithCredentials(
            $request->input('username'),
            $request->input('password')
        );

        if (!$user) {
            return response()->json(['message' => '账号或密码错误，请核对后重试'], 401);
        }

        $token = $user->createToken('auth_token')->plainTextToken;

        return response()->json([
            'data' => [
                'token' => $token,
                'user' => [
                    'id' => $user->id,
                    'username' => $user->username,
                    'name' => $user->name,
                    'role' => $user->role,
                    'school_id' => $user->school_id,
                ],
            ],
        ]);
    }

    /**
     * 班级码登录（学生端/班级大屏统一入口）
     */
    public function classLogin(Request $request): JsonResponse
    {
        $request->validate([
            'class_code' => 'required|string|max:10',
        ]);

        $classCode = strtoupper(trim($request->input('class_code')));
        $classes = ClassRoom::where('display_code', $classCode)->get();

        // 班级码为确定性（LS + 年级 + 班号，如 LS11），数据库未刷新时按实时计算匹配
        if ($classes->isEmpty()) {
            $classes = ClassRoom::get()
                ->filter(fn ($c) => \App\Services\DisplayCodeService::generate($c) === $classCode)
                ->values();
        }

        if ($classes->isEmpty()) {
            return response()->json(['message' => '班级码无效，请核对后重试'], 401);
        }

        if ($classes->count() > 1) {
            return response()->json(['message' => '班级码在多所学校存在，请向班主任确认'], 401);
        }

        $class = $classes->first();

        $token = 'class_' . $classCode . '_' . Str::random(32);

        \Illuminate\Support\Facades\Cache::put('class_token:' . $token, $class->id, now()->addHours(24));

        return response()->json([
            'data' => [
                'token' => $token,
                'class_id' => $class->id,
                'class_name' => $class->name,
                'grade' => $class->grade,
                'student_count' => $class->students()->where('status', 'active')->count(),
            ],
        ]);
    }

    /**
     * 微信扫码登录
     * 第三方首次登录时,前端应同时传 nick (昵称) 和 avatar (头像URL)
     */
    public function teacherLoginWithWechat(Request $request): JsonResponse
    {
        $request->validate([
            'openid' => 'required|string',
            'unionid' => 'nullable|string',
            'nick' => 'nullable|string|max:80',
            'avatar' => 'nullable|string|max:500',
        ]);

        $result = $this->authService->loginWithWechat(
            $request->input('openid'),
            $request->input('unionid'),
            $request->input('nick'),
            $request->input('avatar')
        );

        return response()->json(['data' => $result]);
    }

    /**
     * 企业微信扫码登录
     *
     * 两种调用方式：
     * 1. 直接传入 userid（前端已通过 JSAPI 获取）
     * 2. 传入 code（后端自动换取 userid）
     */
    public function teacherLoginWithWechatWork(Request $request): JsonResponse
    {
        $request->validate([
            'userid' => 'nullable|string',
            'code' => 'nullable|string',
            'nick' => 'nullable|string|max:80',
            'avatar' => 'nullable|string|max:500',
        ]);

        $userid = $request->input('userid');

        // 如果传的是 code，先换取 userid
        if (!$userid && $code = $request->input('code')) {
            $school = School::first();
            if ($school) {
                $wechatWork = app(WechatWorkService::class);
                $userid = $wechatWork->getUserIdByCode($school->id, $code);
            }
        }

        if (!$userid) {
            return response()->json(['message' => '无法获取企业微信用户身份'], 400);
        }

        $result = $this->authService->loginWithWechatWork(
            $userid,
            $request->input('nick'),
            $request->input('avatar')
        );

        return response()->json(['data' => $result]);
    }

    /**
     * 获取当前学校第三方平台的扫码授权 URL（前端生成二维码）
     */
    public function thirdPartyAuthUrl(Request $request): JsonResponse
    {
        $school = School::first();
        if (!$school) {
            return response()->json(['message' => '系统尚未初始化'], 400);
        }
        $platform = $school->settings['third_party_platform'] ?? null;
        try {
            $provider = app(\App\Services\ThirdParty\ThirdPartyManager::class)->providerFor($platform);
        } catch (\Throwable) {
            return response()->json(['message' => '未配置第三方平台，请在后台学校设置中选择'], 400);
        }
        $redirectUri = (string) $request->input('redirect_uri', url('/login'));

        return response()->json(['data' => [
            'platform' => $provider->key(),
            'auth_url' => $provider->authUrl($school->id, $redirectUri),
        ]]);
    }

    /**
     * 第三方平台扫码登录回调（企业微信 / 钉钉 / 飞书，按学校配置的平台分发）
     */
    public function thirdPartyLogin(Request $request): JsonResponse
    {
        $request->validate([
            'code' => 'required|string',
            'state' => 'nullable|string',
        ]);

        $school = School::first();
        if (!$school) {
            return response()->json(['message' => '系统尚未初始化'], 400);
        }
        $platform = $school->settings['third_party_platform'] ?? null;
        try {
            $provider = app(\App\Services\ThirdParty\ThirdPartyManager::class)->providerFor($platform);
            $userInfo = $provider->getUserByCode($school->id, $request->input('code'));
        } catch (\Throwable $e) {
            return response()->json(['message' => $e->getMessage()], 400);
        }

        $result = $this->authService->loginWithThirdParty($provider->key(), $userInfo);
        if (($result['status'] ?? '') !== 'logged_in') {
            return response()->json(['message' => $result['message'] ?? '登录失败'], 400);
        }

        return response()->json(['data' => $result]);
    }

    /**
     * QQ扫码登录
     */
    public function teacherLoginWithQQ(Request $request): JsonResponse
    {
        $request->validate([
            'openid' => 'required|string',
            'nick' => 'nullable|string|max:80',
            'avatar' => 'nullable|string|max:500',
        ]);

        $result = $this->authService->loginWithQQ(
            $request->input('openid'),
            $request->input('nick'),
            $request->input('avatar')
        );

        return response()->json(['data' => $result]);
    }

    /**
     * 人人通登录
     */
    public function teacherLoginWithRenren(Request $request): JsonResponse
    {
        $request->validate([
            'user_id' => 'required|string',
            'nick' => 'nullable|string|max:80',
            'avatar' => 'nullable|string|max:500',
        ]);

        $result = $this->authService->loginWithRenren(
            $request->input('user_id'),
            $request->input('nick'),
            $request->input('avatar')
        );

        return response()->json(['data' => $result]);
    }

    /**
     * 扫码后绑定账号
     * 第三方平台昵称/头像会同步给绑定的 user
     */
    public function bindAfterScan(Request $request): JsonResponse
    {
        $request->validate([
            'temp_token' => 'required|string',
            'username' => 'required|string',
            'password' => 'required|string',
            'platform' => 'required|string',
            'platform_id' => 'required|string',
            'unionid' => 'nullable|string',
            'nick' => 'nullable|string|max:80',
            'avatar' => 'nullable|string|max:500',
        ]);

        $result = $this->authService->bindAfterScan(
            $request->input('temp_token'),
            $request->input('username'),
            $request->input('password'),
            $request->input('platform'),
            $request->input('platform_id'),
            $request->input('unionid'),
            $request->input('nick'),
            $request->input('avatar')
        );

        return response()->json(['data' => $result]);
    }

    /**
     * 登出
     */
    public function logout(Request $request): JsonResponse
    {
        $request->user()->currentAccessToken()->delete();

        return response()->json(['message' => '已登出']);
    }

    /**
     * 修改密码
     */
    public function changePassword(Request $request): JsonResponse
    {
        $request->validate([
            'current_password' => 'required|string',
            'new_password' => 'required|string|min:6|max:50',
        ]);

        $user = $request->user();

        if (!Hash::check($request->input('current_password'), $user->password)) {
            return response()->json(['message' => '当前密码不正确'], 422);
        }

        $user->password = Hash::make($request->input('new_password'));
        $user->plain_password = $request->input('new_password');
        $user->password_changed = true;
        $user->save();

        return response()->json(['message' => '密码修改成功']);
    }

    /**
     * 刷新Token
     */
    public function refreshToken(Request $request): JsonResponse
    {
        $user = $request->user();
        $user->currentAccessToken()->delete();
        $token = $user->createToken('auth-token')->plainTextToken;

        return response()->json(['data' => ['token' => $token]]);
    }

    /**
     * 绑定第三方账号
     */
    public function bindThirdParty(Request $request, string $platform): JsonResponse
    {
        $request->validate([
            'platform_id' => 'required|string',
            'nick' => 'nullable|string|max:80',
            'avatar' => 'nullable|string|max:500',
        ]);

        $user = $request->user();

        $existing = ThirdPartyBinding::where('platform', $platform)
            ->where('platform_id', $request->input('platform_id'))
            ->first();

        if ($existing) {
            return response()->json(['message' => '该第三方账号已被其他用户绑定'], 422);
        }

        ThirdPartyBinding::create([
            'user_id' => $user->id,
            'platform' => $platform,
            'platform_id' => $request->input('platform_id'),
            'platform_nick' => $request->input('nick'),
            'platform_avatar' => $request->input('avatar'),
            'verified_at' => now(),
        ]);

        return response()->json(['message' => '绑定成功']);
    }

    /**
     * 解绑第三方账号
     */
    public function unbindThirdParty(Request $request, string $platform): JsonResponse
    {
        $user = $request->user();

        ThirdPartyBinding::where('user_id', $user->id)
            ->where('platform', $platform)
            ->delete();

        return response()->json(['message' => '解绑成功']);
    }

    /**
     * 获取绑定列表
     */
    public function getBindings(Request $request): JsonResponse
    {
        $user = $request->user();
        $bindings = ThirdPartyBinding::where('user_id', $user->id)->get();

        $platforms = ['wechat', 'wechat_work', 'qq', 'renren'];
        $data = collect($platforms)->map(fn ($p) => [
            'platform' => $p,
            'label' => ThirdPartyBinding::platformLabels()[$p] ?? $p,
            'bound' => $bindings->firstWhere('platform', $p) !== null,
            'nick' => $bindings->firstWhere('platform', $p)?->platform_nick,
        ]);

        return response()->json(['data' => $data]);
    }
}
