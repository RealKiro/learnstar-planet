<?php

declare(strict_types=1);

return [
    'corp_id' => env('WECHAT_WORK_CORPID', ''),
    'agent_id' => (int) env('WECHAT_WORK_AGENTID', 0),
    'secret' => env('WECHAT_WORK_SECRET', ''),
    'token' => env('WECHAT_WORK_TOKEN', ''),
    'encoding_aes_key' => env('WECHAT_WORK_ENCODING_AES_KEY', ''),
    'leave_template_id' => env('WECHAT_WORK_LEAVE_TEMPLATE_ID', ''),
];
