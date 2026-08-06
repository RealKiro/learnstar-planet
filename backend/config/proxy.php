<?php

declare(strict_types=1);

// 反向代理信任配置（独立于框架默认 app.php，避免 bootstrap 直调 env()）
return [
    'trust_all' => (bool) env('TRUST_ALL_PROXIES', false),
];
