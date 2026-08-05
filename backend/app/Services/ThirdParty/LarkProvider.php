<?php

declare(strict_types=1);

namespace App\Services\ThirdParty;

use Illuminate\Support\Facades\Http;

class LarkProvider implements ThirdPartyProvider
{
    public function key(): string
    {
        return 'feishu';
    }

    public function authUrl(int $schoolId, string $redirectUri): string
    {
        return 'https://accounts.feishu.cn/open-apis/authen/v1/authorize?client_id=' . urlencode((string) config('feishu.app_id'))
            . '&redirect_uri=' . urlencode($redirectUri)
            . '&response_type=code&state=' . $schoolId;
    }

    public function getUserByCode(int $schoolId, string $code): array
    {
        $r = Http::timeout(10)->post('https://accounts.feishu.cn/oauth/v3/token', [
            'grant_type' => 'authorization_code',
            'client_id' => config('feishu.app_id'),
            'client_secret' => config('feishu.app_secret'),
            'code' => $code,
        ])->json();
        $token = (string) ($r['access_token'] ?? '');
        if ($token === '') {
            throw new \RuntimeException('飞书登录失败：' . (string) ($r['error_description'] ?? ''));
        }
        $me = Http::timeout(10)->withHeaders(['Authorization' => 'Bearer ' . $token])
            ->get('https://open.feishu.cn/open-apis/authen/v1/user_info')->json();
        $data = $me['data'] ?? [];

        return [
            'platform_id' => (string) ($data['union_id'] ?? $data['open_id'] ?? $code),
            'name' => (string) ($data['name'] ?? ''),
            'mobile' => (string) ($data['mobile'] ?? ''),
            'avatar' => (string) ($data['avatar_url'] ?? ''),
        ];
    }

    public function fetchContacts(int $schoolId): array
    {
        $token = $this->tenantToken();
        if ($token === '') {
            throw new \RuntimeException('飞书未配置或获取 tenant_access_token 失败');
        }
        $departments = $this->allDepartments($token);
        $members = [];
        $seen = [];
        foreach ($departments as $d) {
            $deptId = (string) ($d['department_id'] ?? '');
            if ($deptId === '') {
                continue;
            }
            foreach ($this->deptUsers($token, $deptId) as $u) {
                $uid = (string) ($u['union_id'] ?? $u['user_id'] ?? '');
                if ($uid === '' || isset($seen[$uid])) {
                    continue;
                }
                $seen[$uid] = true;
                $members[] = [
                    'userid' => $uid,
                    'name' => (string) ($u['name'] ?? ''),
                    'mobile' => (string) ($u['mobile'] ?? ''),
                    'email' => (string) ($u['email'] ?? ''),
                    'position' => (string) ($u['job_title'] ?? ''),
                    'department_names' => $u['department_ids'] ?? [],
                ];
            }
        }

        return ['departments' => $departments, 'members' => $members];
    }

    private function tenantToken(): string
    {
        $r = Http::timeout(10)->post('https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal', [
            'app_id' => config('feishu.app_id'),
            'app_secret' => config('feishu.app_secret'),
        ])->json();

        return (string) ($r['tenant_access_token'] ?? '');
    }

    private function allDepartments(string $token): array
    {
        $r = Http::timeout(10)->withHeaders(['Authorization' => 'Bearer ' . $token])
            ->get('https://open.feishu.cn/open-apis/contact/v3/departments', ['fetch_child' => 'true', 'page_size' => 100])->json();

        return $r['data']['items'] ?? [];
    }

    private function deptUsers(string $token, string $deptId): array
    {
        $r = Http::timeout(10)->withHeaders(['Authorization' => 'Bearer ' . $token])
            ->get('https://open.feishu.cn/open-apis/contact/v3/users/find_by_department', [
                'department_id' => $deptId,
                'page_size' => 100,
            ])->json();

        return $r['data']['items'] ?? [];
    }
}
