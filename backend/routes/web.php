<?php

declare(strict_types=1);

use Illuminate\Support\Facades\Route;

Route::get('/health', function () {
    return response()->json(['status' => 'ok', 'timestamp' => now()->toISOString()]);
});

Route::get('/up', function () {
    return response('OK');
});

// SPA 兜底路由：所有非 API、非静态文件的 GET 请求返回 Vue 前端入口
Route::get('/{any?}', function () {
    $path = public_path('index.html');
    if (!file_exists($path)) {
        return response('Frontend not built. Run `npm run build` in frontend-vue/.', 500);
    }

    return response(file_get_contents($path), 200, [
        'Content-Type' => 'text/html; charset=UTF-8',
    ]);
})->where('any', '^(?!api/|_debugbar|telescope|horizon|nova|storage).*$');
