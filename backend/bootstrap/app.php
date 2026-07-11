<?php

declare(strict_types=1);

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__ . '/../routes/web.php',
        api: __DIR__ . '/../routes/api.php',
        commands: __DIR__ . '/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        // 信任反向代理头（Cloudflare Tunnel / Nginx / Traefik 等）
        // 服务器部署模式下默认信任所有代理，确保 HTTPS 和域名正确识别
        if (env('APP_ENV') === 'production' || env('TRUST_ALL_PROXIES', false)) {
            $middleware->trustProxies(at: '*');
        }

        $middleware->validateCsrfTokens(except: [
            'api/*',
        ]);

        $middleware->api(prepend: [
            \Laravel\Sanctum\Http\Middleware\EnsureFrontendRequestsAreStateful::class,
        ]);

        $middleware->alias([
            'role' => \App\Http\Middleware\RoleMiddleware::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        // API 统一异常响应（始