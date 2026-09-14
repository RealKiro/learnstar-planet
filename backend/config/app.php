<?php

declare(strict_types=1);

/**
 * 应用配置（只写需要覆盖的项，其余沿用 Laravel 框架默认值）
 *
 * 本项目 config/ 目录刻意保持精简：Laravel 11+ 会把应用 config 与框架基础配置
 * 做 array_merge（见 Illuminate\Foundation\Bootstrap\LoadConfiguration），
 * 因此这里不写 name/providers/aliases 等键也不会丢失。
 *
 * 注意：框架基础配置把 timezone 硬编码为 UTC（不读环境变量），
 * 只能由本文件覆盖。
 */

return [
    /*
     * 时区：默认东八区（北京时间）
     *
     * 为什么必须改：框架默认 UTC 会让「今天」的边界错到北京时间 08:00。
     * 本项目大量使用 today() / whereDate('created_at', today()) / 报表按日聚合
     * （考勤、今日积分、日榜、日报表），UTC 时区下北京时间 00:00–08:00 产生的记录
     * 会被算进前一天。框架引导阶段会用本值调用 date_default_timezone_set()，
     * 故 Web / 队列 worker / artisan 命令三处行为一致。
     */
    'timezone' => env('APP_TIMEZONE', 'Asia/Shanghai'),
];
