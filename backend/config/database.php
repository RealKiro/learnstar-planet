<?php

declare(strict_types=1);

/**
 * 班宠星球 - 数据库配置
 * 支持 MySQL / MariaDB / PostgreSQL / SQLite
 * 通过 .env 中 DB_CONNECTION 切换
 */

use Illuminate\Support\Env;

return [
    'default' => Env::get('DB_CONNECTION', 'mysql'),

    'connections' => [
        // MySQL 配置
        'mysql' => [
            'driver' => 'mysql',
            'url' => Env::get('DATABASE_URL'),
            'host' => Env::get('DB_HOST', '127.0.0.1'),
            'port' => Env::get('DB_PORT', '3306'),
            'database' => Env::get('DB_DATABASE', 'bancxq_planet'),
            'username' => Env::get('DB_USERNAME', 'bancxq'),
            'password' => Env::get('DB_PASSWORD', ''),
            'unix_socket' => Env::get('DB_SOCKET', ''),
            'charset' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
            'prefix' => '',
            'prefix_indexes' => true,
            'strict' => true,
            'engine' => null,
            'options' => extension_loaded('pdo_mysql') ? array_filter([
                PDO::MYSQL_ATTR_SSL_CA => Env::get('MYSQL_ATTR_SSL_CA'),
                PDO::MYSQL_ATTR_SSL_VERIFY_SERVER_CERT => false,
                PDO::MYSQL_ATTR_MULTI_STATEMENTS => false,
            ]) : [],
        ],

        // MariaDB 配置（与 MySQL 共用驱动，端口可能不同）
        'mariadb' => [
            'driver' => 'mysql',
            'url' => Env::get('DATABASE_URL'),
            'host' => Env::get('DB_HOST', '127.0.0.1'),
            'port' => Env::get('DB_PORT', '3306'),
            'database' => Env::get('DB_DATABASE', 'bancxq_planet'),
            'username' => Env::get('DB_USERNAME', 'bancxq'),
            'password' => Env::get('DB_PASSWORD', ''),
            'unix_socket' => Env::get('DB_SOCKET', ''),
            'charset' => 'utf8mb4',
            'collation' => 'utf8mb4_unicode_ci',
            'prefix' => '',
            'prefix_indexes' => true,
            'strict' => true,
            'engine' => null,
            'options' => extension_loaded('pdo_mysql') ? array_filter([
                PDO::MYSQL_ATTR_SSL_CA => Env::get('MYSQL_ATTR_SSL_CA'),
            ]) : [],
        ],

        // PostgreSQL 配置
        'pgsql' => [
            'driver' => 'pgsql',
            'url' => Env::get('DATABASE_URL'),
            'host' => Env::get('DB_HOST', '127.0.0.1'),
            'port' => Env::get('DB_PORT', '5432'),
            'database' => Env::get('DB_DATABASE', 'bancxq_planet'),
            'username' => Env::get('DB_USERNAME', 'bancxq'),
            'password' => Env::get('DB_PASSWORD', ''),
            'charset' => 'utf8',
            'prefix' => '',
            'prefix_indexes' => true,
            'search_path' => 'public',
            'sslmode' => 'prefer',
        ],

        // SQLite 配置（轻量开发/测试）
        'sqlite' => [
            'driver' => 'sqlite',
            'url' => Env::get('DATABASE_URL'),
            'database' => Env::get('DB_DATABASE', database_path('database.sqlite')),
            'prefix' => '',
            'prefix_indexes' => true,
        ],
    ],

    'migrations' => [
        'table' => 'migrations',
        'update_date_on_publish' => true,
    ],

    'redis' => [
        'client' => Env::get('REDIS_CLIENT', 'predis'),

        'default' => [
            'url' => Env::get('REDIS_URL'),
            'host' => Env::get('REDIS_HOST', 'redis'),
            'username' => Env::get('REDIS_USERNAME'),
            'password' => Env::get('REDIS_PASSWORD'),
            'port' => Env::get('REDIS_PORT', '6379'),
            'database' => Env::get('REDIS_DB', '0'),
        ],

        // 排行榜专用 Redis 连接
        'leaderboard' => [
            'url' => Env::get('REDIS_URL'),
            'host' => Env::get('REDIS_HOST', 'redis'),
            'username' => Env::get('REDIS_USERNAME'),
            'password' => Env::get('REDIS_PASSWORD'),
            'port' => Env::get('REDIS_PORT', '6379'),
            'database' => Env::get('REDIS_LEADERBOARD_DB', '1'),
        ],

        // 缓存专用 Redis 连接
        'cache' => [
            'url' => Env::get('REDIS_URL'),
            'host' => Env::get('REDIS_HOST', 'redis'),
            'username' => Env::get('REDIS_USERNAME'),
            'password' => Env::get('REDIS_PASSWORD'),
            'port' => Env::get('REDIS_PORT', '6379'),
            'database' => Env::get('REDIS_CACHE_DB', '2'),
        ],

        // 会话专用 Redis 连接
        'session' => [
            'url' => Env::get('REDIS_URL'),
            'host' => Env::get('REDIS_HOST', 'redis'),
            'username' => Env::get('REDIS_USERNAME'),
            'password' => Env::get('REDIS_PASSWORD'),
            'port' => Env::get('REDIS_PORT', '6379'),
            'database' => Env::get('REDIS_SESSION_DB', '3'),
        ],
    ],
];
