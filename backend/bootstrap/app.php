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
        // API 统一异常响应（始终返回 JSON）
        $exceptions->shouldRenderJsonWhen(function () {
            return true;
        });

        // 404 — 模型未找到
        $exceptions->renderable(function (\Illuminate\Database\Eloquent\ModelNotFoundException $e) {
            return response()->json([
                'message' => '资源不存在',
            ], 404);
        });

        // 422 — 表单验证失败（统一格式）
        $exceptions->renderable(function (\Illuminate\Validation\ValidationException $e) {
            return response()->json([
                'message' => '参数错误',
                'errors' => $e->errors(),
            ], 422);
        });

        // 401 — 未认证
        $exceptions->renderable(function (\Illuminate\Auth\AuthenticationException $e) {
            return response()->json([
                'message' => '未登录或登录已过期',
            ], 401);
        });

        // 403 — 无权限
        $exceptions->renderable(function (\Illuminate\Auth\Access\AuthorizationException $e) {
            return response()->json([
                'message' => $e->getMessage() ?: '无权执行此操作',
            ], 403);
        });

        // 500 — 其他所有异常
        $exceptions->renderable(function (\Throwable $e) {
            $code = method_exists($e, 'getStatusCode') ? $e->getStatusCode() : 500;
            if ($code < 100 || $code > 599) {
                $code = 500;
            }

            $payload = [
                'message' => app()->environment('production') ? '服务器内部错误' : $e->getMessage(),
            ];

            if (!app()->environment('production')) {
                $payload['file'] = $e->getFile();
                $payload['line'] = $e->getLine();
                $payload['trace'] = explode("\n", $e->getTraceAsString());
            }

            return response()->json($payload, $code);
        });
    })
    ->create();

