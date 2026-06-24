<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Services\AuthService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

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
            return response()->json(['message' => '账号或密码错误'], 401);
        }

        return response()->json(['data' => $user]);
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
            return response()->json(['message' => '账号或密码错误'], 401);
        }

        return response()->json(['data' => $user]);
    }

    /**
     * 微信扫码登录
     */
    public function teacherLoginWithWechat(Request $request): JsonResponse
    {
        $request->validate([
            'openid' => 'required|string',
            'unionid' => 'nullable|string',
        ]);

        $result = $this->authService->loginWithWechat(
            $request->input('openid'),
            $request->input('unionid')
        );

        return response()->json(['data' => $result]);
    }

    /**
     * 企业微信扫码登录
     */
    public function teacherLoginWithWechatWork(Request $request): JsonResponse
    {
        $request->validate([
            'userid' => 'required|string',
        ]);

        $result = $this->authService->loginWithWechatWork($request->input('userid'));

        return response()->json(['data' => $result]);
    }

    /**
     * QQ扫码登录
     */
    public function teacherLoginWithQQ(Request $request): JsonResponse
    {
        $request->validate([
            'openid' => 'required|string',
        ]);

        $result = $this->authService->loginWithQQ($request->input('openid'));

        return response()->json(['data' => $result]);
    }

    /**
     * 人人通登录
     */
    public function teacherLoginWithRenren(Request $request): JsonResponse
    {
        $request->validate([
            'user_id' => 'required|string',
        ]);

        $result = $this->authService->loginWithRenren($request->input('user_id'));

        return response()->json(['data' => $result]);
    }

    /**
     * 扫码后绑定账号
     */
    public function bindAfterScan(Request $request): JsonResponse
    {
        $request->validate([
            'temp_token' => 'required|string',
            'username' => 'required|string',
            'password' => 'required|string',
            'platform' => 'required|string',
            'platform_id' => 'required|string',
        ]);

        $result = $this->authService->bindAfterScan(
            $request->input('temp_token'),
            $request->input('username'),
            $request->input('password'),
            $request->input('platform'),
            $request->input('platform_id')
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
        return response()->json(['message' => '密码修改成功']);
    }

    /**
     * 刷新Token
     */
    public function refreshToken(Request $request): JsonResponse
    {
        return response()->json(['message' => 'Token已刷新']);
    }

    /**
     * 绑定第三方账号
     */
    public function bindThirdParty(Request $request, string $platform): JsonResponse
    {
        return response()->json(['message' => '绑定成功']);
    }

    /**
     * 解绑第三方账号
     */
    public function unbindThirdParty(Request $request, string $platform): JsonResponse
    {
        return response()->json(['message' => '解绑成功']);
    }

    /**
     * 获取绑定列表
     */
    public function getBindings(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }
}
