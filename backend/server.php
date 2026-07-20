<?php

/**
 * Laravel - PHP Framework
 *
 * 此文件用于 PHP 内置开发服务器的路由规则。
 * 作用：SPA 路由回退——Vue Router history 模式的 URL
 * （如 /teacher/dashboard）不会找到静态文件，需要交给 index.php 处理。
 */

$uri = urldecode(parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH));

// 如果请求的是真实文件或目录，直接返回
if ($uri !== '/' && file_exists(__DIR__ . '/public' . $uri)) {
    return false;
}

// 所有其他请求交给 Laravel 处理
require_once __DIR__ . '/public/index.php';
