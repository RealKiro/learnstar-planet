<?php

declare(strict_types=1);

namespace App\Services\ThirdParty;

/**
 * 第三方办公平台适配接口（企业微信 / 钉钉 / 飞书）
 */
interface ThirdPartyProvider
{
    /**
     * 平台标识：wechat_work / dingtalk / feishu
     */
    public function key(): string;

    /**
     * 扫码授权 URL（前端生成二维码跳转）
     */
    public function authUrl(int $schoolId, string $redirectUri): string;

    /**
     * 用授权 code 换取用户身份
     *
     * @return array{platform_id:string, name:string, mobile:string, avatar:string}
     */
    public function getUserByCode(int $schoolId, string $code): array;

    /**
     * 拉取通讯录（部门 + 成员，按 userid 去重）
     *
     * @return array{departments: array, members: array}
     */
    public function fetchContacts(int $schoolId): array;
}
