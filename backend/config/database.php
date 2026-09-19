<?php

declare(strict_types=1);

/**
 * 班宠星球 - 数据库配置
 * 支持 MySQL / MariaDB / PostgreSQL / SQLite
 * 通过 .env 中 DB_CONNECTION 切换
 */

use Illuminate\Support\Env;

/*
 * MySQL 连接属性的常量名兼容层
 *
 * PHP 8.4 起这些属性迁到 Pdo\Mysql::ATTR_*，旧的 PDO::MYSQL_ATTR_* 自 PHP 8.5 起弃用
 * （PHP 9 将移除）。启动日志里成排出现的
 *   "Deprecated: Constant PDO::MYSQL_ATTR_SSL_CA is deprecated since 8.5"
 * 就是这里发出的。按运行时可用的常量名取值，新老 PHP 均无告警。
 * 注意：未装 pdo_mysql 时两组常量都不存在，故取 null —— 使用处仍被
 * extension_loaded('pdo_mysql') 保护，null 不会被写进 options。
 */
$mysqlSslCa = defined('Pdo\Mysql::ATTR_SSL_CA')
    ? Pdo\Mysql::ATTR_SSL_CA
    : (defined('PDO::MYSQL_ATTR_SSL_CA') ? PDO::MYSQL_ATTR_SSL_CA : null);

$mysqlMultiStatements = defined('Pdo\Mysql::ATTR_MULTI_STATEMENTS')
    ? Pdo\Mysql::ATTR_MULTI_STATEMENTS
    : (defined('PDO::MYSQL_ATTR_MULTI_STATEMENTS') ? PDO::MYSQL_ATTR_MULTI_STATEMENTS : null);

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
            // ⚠️ array_filter 的默认行为会连 false 一起丢掉，必须显式只过滤空值，
            //    否则 MULTI_STATEMENTS=false 这条加固会被静默吃掉（历史上就没生效过）。
            'options' => extension_loaded('pdo_mysql') ? array_filter([
                $mysqlSslCa => Env::get('MYSQL_ATTR_SSL_CA'),
                // 关掉叠加查询（防 SQL 注入串联多语句）；本项目无多语句用法（仅一条单句 UPDATE）
                $mysqlMultiStatements => false,
            ], static fn ($value) => $value !== null && $value !== '') : [],
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
                $mysqlSslCa => Env::get('MYSQL_ATTR_SSL_CA'),
            ], static fn ($value) => $value !== null && $value !== '') : [],
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
            'database' => Env::get('DB_DATABASE', storage_path('database.sqlite')),
            'prefix' => '',
            'prefix_indexes' => true,
            // 并发写保护（毫秒）：Octane 多 worker 同时写 SQLite 时，写锁冲突先等待而非立即抛
            // "database is locked"。逐连接生效、不落盘，与 compose 的单文件数据卷兼容（勿改用 WAL：
            // WAL 的 -wal/-shm 副文件不在卷内，容器重建会丢最近写入）
            'busy_timeout' => Env::get('DB_BUSY_TIMEOUT', 5000),
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
