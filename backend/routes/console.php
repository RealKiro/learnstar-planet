<?php

declare(strict_types=1);

use Illuminate\Support\Facades\Schedule;

// ===== 企业微信请假同步 =====
// 工作日早上 7:50 自动从企微拉取当日请假数据
Schedule::command('attendance:sync-wechat-leave')
    ->t