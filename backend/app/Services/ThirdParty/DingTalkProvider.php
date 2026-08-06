<?php

declare(strict_types=1);

namespace App\Services\ThirdParty;

use Illuminate\Support\Facades\Http;

class DingTalkProvider implements ThirdPartyProvider
{
    public function key(): string
    {
        return 'dingtalk';
    }

    public function authUrl(int $schoolId, string $redirectUri): string
    {
        return 'https://login.dingtalk.com/oauth2/auth?redirect_uri=' . urlencode($redirectUri)
            . '&response_type=code&client_id=' . urlencode((string) config('dingtalk.app_key'))
            . '&scope=openid&state=' . $schoolId;
    }

    public function getUserByCode(int $schoolId, string $code): array
    {
        $r = Http::timeout(10)->post('https://api.dingtalk.com/v1.0/oauth2/userAccessToken', [
            'clientId' => config('dingtalk.app_key'),
            'clientSecret' => config('dingtalk.app_secret'),
            'code' => $code,
            'grantType' => 'authorization_code',
        ])->json();
        $token = (string) ($r['accessToken'] ?? '');
        if ($token === '') {
            throw new \RuntimeException('钉钉登录失败：' . (string) ($r['message'] ?? ''));
        }
        $me = Http::timeout(10)->withHeaders(['x-acs-dingtalk-access-token' => $token])
            ->get('https://api.dingtalk.com/v1.0/contact/users/me')->json();

        return [
            'platform_id' => (string) ($me['unionId'] ?? $me['userId'] ?? $code),
            'name' => (string) ($me['nick'] ?? $me['name'] ?? ''),
            'mobile' => (string) ($me['mobile'] ?? ''),
            'email' => (string) ($me['email'] ?? ''),
            'avatar' => (string) ($me['avatarUrl'] ?? ''),
        ];
    }

    public function fetchContacts(int $schoolId): array
    {
        $token = $this->internalToken();
        if ($token === '') {
            throw new \RuntimeException('钉钉未配置或获取企业内部 token 失败');
        }
        $departments = $this->allDepartments($token);
        $nameById = [];
        foreach ($departments as $d) {
            $nameById[(int) ($d['dept_id'] ?? 0)] = (string) ($d['name'] ?? '');
        }

        $members = [];
        $seen = [];
        foreach ($departments as $d) {
            $deptId = (int) ($d['dept_id'] ?? 0);
            if ($deptId <= 0) {
                continue;
            }
            foreach ($this->deptUsers($token, $deptId) as $u) {
                $uid = (string) ($u['userid'] ?? '');
                if ($uid === '' || isset($seen[$uid])) {
                    continue;
                }
                $seen[$uid] = true;

                $deptNames = [];
                foreach ($u['dept_id_list'] ?? [] as $did) {
                    $n = $nameById[(int) $did] ?? null;
                    if ($n !== null) {
                        $deptNames[] = $n;
                    }
                }

                $members[] = [
                    'userid' => $uid,
                    'name' => (string) ($u['name'] ?? ''),
                    'mobile' => (string) ($u['mobile'] ?? ''),
                    'email' => (string) ($u['email'] ?? ''),
                    'position' => (string) ($u['title'] ?? ''),
                    'department_names' => $deptNames,
                ];
            }
        }

        return ['departments' => $departments, 'members' => $members];
    }

    private function internalToken(): string
    {
        $r = Http::timeout(10)->get('https://oapi.dingtalk.com/gettoken', [
            'appkey' => config('dingtalk.app_key'),
            'appsecret' => config('dingtalk.app_secret'),
        ])->json();

        return (string) ($r['access_token'] ?? '');
    }

    private function allDepartments(string $token): array
    {
        $r = Http::timeout(10)->post("https://oapi.dingtalk.com/topapi/v2/department/listsub?access_token={$token}", [
            'dept_id' => 1,
        ])->json();

        return $r['result']['list'] ?? [];
    }

    private function deptUsers(string $token, int $deptId): array
    {
        $r = Http::timeout(10)->post("https://oapi.dingtalk.com/topapi/v2/user/list?access_token={$token}", [
            'dept_id' => $deptId,
            'cursor' => 0,
            'size' => 100,
        ])->json();

        return $r['result']['list'] ?? [];
    }
}
