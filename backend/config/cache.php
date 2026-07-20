<?php

declare(strict_types=1);

/**
 * 班宠星球 - 缓存配置
 * 使用独立部署的 Redis，不内置
 */

use Illuminate\Support\Env;

return [
    'default' => Env::get('CACHE_DRIVER', 'redis'),

    'stores' => [
        'apc' => [
            'driver' => 'apc',
        ],

        'array' => [
            'driver' => 'array',
            'serialize' => false,
        ],

        'database' => [
            'driver' => 'database',
            'table' => 'cache',
            'connection' => null,
            'lock_connection' => null,
        ],

        'file' => [
            'driver' => 'file',
            'path' => storage_path('framework/cache/data'),
            'lock_path' => storage_path('framework/cache/data'),
        ],

        'memcached' => [
            'driver' => 'memcached',
            'persistent_id' => Env::get('MEMCACHED_PERSISTENT_ID'),
            'sasl' => [
                Env::get('MEMCACHED_USERNAME'),
                Env::get('MEMCACHED_PASSWORD'),
            ],
            'options' => [
                // Memcached::OPT_CONNECT_TIMEOUT => 2000,
            ],
            'servers' => [
                [
                    'host' => Env::get('MEMCACHED_HOST', '127.0.0.1'),
                    'port' => Env::get('MEMCACHED_PORT', 11211),
                    'weight' => 100,
                ],
            ],
        ],

        'redis' => [
            'driver' => 'redis',
            'connection' => 'cache',
            'lock_connection' => 'cache',
        ],
    ],

    'prefix' => Env::get('CACHE_PREFIX', 'bancxq_cache_'),
];

