<?php

declare(strict_types=1);

namespace App\Services\ThirdParty;

use App\Services\WechatWorkService;
use Illuminate\Support\Facades\Http;

class WeChatWorkProvider implements ThirdPartyProvider
{
    public function key(): string
    {
        return 'wechat_work';
    }

    public function authUrl(int $schoolId, string $redirectUri): string
    {
        return 'https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=' . urlencode((string) config('wechat-work.corp_id'))
            . '&agentid=' . (int) config('wechat-work.agent_id')
            . '&redirect_uri=' . urlencode($redirectUri)
            . '&state=' . $schoolId;
    }

    public function getUserByCode(int $schoolId, string $code): array
    {
        $svc = app(WechatWorkService::class);
        $userid = $svc->getUserIdByCode($schoolId, $code);
        if (!$userid) {
            throw new \RuntimeException('企业微信登录失败');
        }
        $detail = $this->userDetail($schoolId, $userid);

        return [
            'platform_id' => $userid,
            'name' => $detail['name'] ?: $userid,
            'mobile' => $detail['mobile'] ?? '',
            'avatar' => $detail['avatar'] ?? '',
        ];
    }

    public function fetchContacts(int $schoolId): array
    {
        return app(WechatWorkService::class)->fetchContacts($schoolId);
    }

    private function userDetail(int $schoolId, string $userid): array
    {
        $t = app(WechatWorkService::class)->getAccessToken($schoolId);
        $r = Http::timeout(10)->get('https://qyapi.weixin.qq.com/cgi-bin/user/get', [
            'access_token' => $t,
            'userid' => $userid,
        ])->json();
        if (($r['errcode'] ?? -1) !== 0) {
            return [];
        }

        return [
            'name' => (string) ($r['name'] ?? ''),
            'mobile' => (string) ($r['mobile'] ?? ''),
            'avatar' => (string) ($r['avatar'] ?? ''),
        ];
    }
}
