<?php

declare(strict_types=1);

namespace App\Services;

use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;

class WechatWorkService
{
    public function getAccessToken(int $schoolId): string
    {
        $k = "wecom_at:{$schoolId}";
        $t = Cache::get($k);
        if (is_string($t) && $t !== '') {
            return $t;
        }
        $cid = config('wechat-work.corp_id');
        $sec = config('wechat-work.secret');
        if (empty($cid) || empty($sec)) {
            throw new \RuntimeException('企微未配置');
        }
        $r = Http::timeout(10)->get('https://qyapi.weixin.qq.com/cgi-bin/gettoken', ['corpid' => $cid,'corpsecret' => $sec])->json();
        if (!isset($r['access_token'])) {
            Log::error('企微token失败', $r);

            throw new \RuntimeException('token失败');
        }
        Cache::put($k, $r['access_token'], max(($r['expires_in'] ?? 7200) - 300, 60));

        return $r['access_token'];
    }

    public function getLeaveApprovalSpNos(int $schoolId, int $start, int $end): array
    {
        $t = $this->getAccessToken($schoolId);
        $sp = [];
        $c = 0;
        do {
            $r = Http::timeout(10)->post("https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovalinfo?access_token={$t}", ['starttime' => $start,'endtime' => $end,'cursor' => $c,'size' => 100,'filters' => ['key' => 'sp_status','value' => '2']])->json();
            if (($r['errcode'] ?? -1) !== 0) {
                break;
            }
            foreach ($r['sp_no_list'] ?? [] as $n) {
                if (is_string($n) && $n !== '') {
                    $sp[] = $n;
                }
            }
            $c = $r['next_cursor'] ?? 0;
        } while ($c > 0);

        return $sp;
    }

    public function getApprovalDetail(int $schoolId, string $spNo): ?array
    {
        $t = $this->getAccessToken($schoolId);
        $r = Http::timeout(10)->post("https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovaldetail?access_token={$t}", ['sp_no' => $spNo])->json();
        if (($r['errcode'] ?? -1) !== 0) {
            return null;
        }
        $info = $r['info'] ?? [];
        if (!is_array($info) || empty($info)) {
            return null;
        }

        return $this->parse($info);
    }

    private function parse(array $info): array
    {
        $spNo = (string) ($info['sp_no'] ?? '');
        $uid = (string) ($info['applyer']['userid'] ?? '');
        $st = (int) ($info['sp_status'] ?? 1);
        $at = null;
        foreach ($info['sp_record'] ?? [] as $rc) {
            if (($rc['sp_status'] ?? 0) === 2 && isset($rc['approve_time'])) {
                $at = (int) $rc['approve_time'];
            }
        }
        $sn = null;
        $lt = null;
        $rs = null;
        $ls = 0;
        $le = 0;
        foreach ($info['apply_data']['contents'] ?? [] as $c) {
            $tl = $c['title'] ?? '';
            $v = $c['value'] ?? '';
            $vt = is_array($v) ? (is_string($v['text'] ?? null) ? $v['text'] : '') : (string) $v;
            if ($this->mf($tl, ['学生姓名','学生','姓名'])) {
                $sn = $vt ?: ((is_array($v) && isset($v['controls'][0]['text'])) ? $v['controls'][0]['text'] : null);
            } elseif ($this->mf($tl, ['请假类型','类型'])) {
                $lt = $vt;
            } elseif ($this->mf($tl, ['请假事由','事由','原因','说明'])) {
                $rs = $vt;
            } elseif ($this->mf($tl, ['请假时间','起止时间','开始时间'])) {
                if (is_array($v)) {
                    $ls = (int) ($v['new_begin'] ?? $v['new_start'] ?? 0);
                    $le = (int) ($v['new_end'] ?? 0);
                }
            }
        }

        return ['sp_no' => $spNo,'submitter_userid' => $uid,'student_name' => $sn,'leave_type' => $lt,'leave_start' => $ls,'leave_end' => $le,'reason' => $rs,'approve_status' => $st,'approved_at' => $at];
    }

    /**
     * 通过 OAuth code 换取企业微信 UserId
     */
    public function getUserIdByCode(int $schoolId, string $code): ?string
    {
        $token = $this->getAccessToken($schoolId);
        $r = Http::timeout(10)->get('https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo', [
            'access_token' => $token,
            'code' => $code,
        ])->json();

        if (isset($r['errcode']) && $r['errcode'] !== 0) {
            Log::warning('企微 code 换取 userid 失败', $r);

            return null;
        }

        return $r['UserId'] ?? $r['userid'] ?? null;
    }

    private function mf(string $t, array $ks): bool
    {
        foreach ($ks as $k) {
            if (mb_strpos($t, $k) !== false) {
                return true;
            }
        }

        return false;
    }
}
