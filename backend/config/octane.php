<?php

declare(strict_types=1);

return [
    'server' => 'roadrunner',
    'rr_server' => env('RR_SERVER', '/usr/local/bin/rr'),
    'https' => env('OCTANE_HTTPS', false),
    'host' => '0.0.0.0',
    'port' => (int) env('APP_PORT', 8080),
    'workers' => (int) env('OCTANE_WORKERS', 8),
    'max_requests' => (int) env('OCTANE_MAX_REQUESTS', 500),
    'tasks_enabled' => false,
    'max_execution_time' => 0,
    'swoole' => [
        'options' => [],
    ],
];
