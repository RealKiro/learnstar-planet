<?php

declare(strict_types=1);

namespace App\Services\ThirdParty;

/**
 * 按学校配置的平台返回对应 provider（每校一平台）
 */
class ThirdPartyManager
{
    public function providerFor(?string $platform): ThirdPartyProvider
    {
        return match ($platform) {
            'wechat_work' => new WeChatWorkProvider(),
            'dingtalk' => new DingTalkProvider(),
            'feishu' => new LarkProvider(),
            default => throw new \RuntimeException('未配置或未知的第三方平台'),
        };
    }

    public function availablePlatforms(): array
    {
        return [
            'wechat_work' => '企业微信',
            'dingtalk' => '钉钉',
            'feishu' => '飞书',
        ];
    }
}
