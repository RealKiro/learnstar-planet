<?php

declare(strict_types=1);

use Illuminate\Support\Facades\Route;

Route::get('/health', function () {
    return response()->json(['status' => 'ok', 'timestamp' => now()->toISOString()]);
});

Route::get('/up', function () {
    return response('OK');
});

// 调试诊断路由（只在非生产环境可用）
Route::get('/debug', function (\Illuminate\Http\Request $req) {
    $results = [];

    // 1. PHP 基本信息
    $results[] = ['test' => 'PHP 版本', 'result' => phpversion()];
    $results[] = ['test' => 'Laravel 版本', 'result' => app()->version()];
    $results[] = ['test' => 'APP_ENV', 'result' => env('APP_ENV', '未设置')];
    $results[] = ['test' => 'APP_DEBUG', 'result' => env('APP_DEBUG', '未设置') ? 'true' : 'false'];

    // 2. 数据库连接测试
    try {
        \Illuminate\Support\Facades\DB::connection()->getPdo();
        $results[] = ['test' => '数据库连接', 'result' => '✅ 成功'];
    } catch (\Throwable $e) {
        $results[] = ['test' => '数据库连接', 'result' => '❌ 失败: ' . $e->getMessage()];
    }

    // 3. 检查 users 表
    try {
        $hasUsersTable = \Illuminate\Support\Facades\Schema::hasTable('users');
        $results[] = ['test' => 'users 表存在', 'result' => $hasUsersTable ? '✅ 是' : '❌ 否'];
        if ($hasUsersTable) {
            $cols = ['nickname', 'subject', 'grade_team'];
            foreach ($cols as $col) {
                $has = \Illuminate\Support\Facades\Schema::hasColumn('users', $col);
                $results[] = ['test' => "users 表.{$col} 字段", 'result' => $has ? '✅ 存在' : '❌ 缺失'];
            }
        }
    } catch (\Throwable $e) {
        $results[] = ['test' => 'users 表检查', 'result' => '❌ 错误: ' . $e->getMessage()];
    }

    // 4. 检查 class_rooms 表
    try {
        $hasClassTable = \Illuminate\Support\Facades\Schema::hasTable('class_rooms');
        $results[] = ['test' => 'class_rooms 表存在', 'result' => $hasClassTable ? '✅ 是' : '❌ 否'];
        if ($hasClassTable) {
            $cols = ['display_code', 'display_code_updated_at'];
            foreach ($cols as $col) {
                $has = \Illuminate\Support\Facades\Schema::hasColumn('class_rooms', $col);
                $results[] = ['test' => "class_rooms 表.{$col} 字段", 'result' => $has ? '✅ 存在' : '❌ 缺失'];
            }
        }
    } catch (\Throwable $e) {
        $results[] = ['test' => 'class_rooms 表检查', 'result' => '❌ 错误: ' . $e->getMessage()];
    }

    // 5. 检查其他表
    foreach (['class_room_teachers', 'third_party_bindings', 'pets', 'scores'] as $table) {
        try {
            $has = \Illuminate\Support\Facades\Schema::hasTable($table);
            $results[] = ['test' => "{$table} 表存在", 'result' => $has ? '✅ 是' : '❌ 否'];
        } catch (\Throwable $e) {
            $results[] = ['test' => "{$table} 表检查", 'result' => '❌ 错误: ' . $e->getMessage()];
        }
    }

    // 6. 测试 Artisan 命令调用
    try {
        $code = \Illuminate\Support\Facades\Artisan::call('demo:clean', ['--force' => true]);
        $results[] = ['test' => 'Artisan demo:clean 调用', 'result' => "✅ 返回码: {$code}"];
    } catch (\Throwable $e) {
        $results[] = ['test' => 'Artisan demo:clean 调用', 'result' => '❌ 错误: ' . $e->getMessage()];
    }

    try {
        $code = \Illuminate\Support\Facades\Artisan::call('demo:seed', ['--force' => true]);
        $results[] = ['test' => 'Artisan demo:seed 调用', 'result' => "✅ 返回码: {$code}"];
    } catch (\Throwable $e) {
        $results[] = ['test' => 'Artisan demo:seed 调用', 'result' => '❌ 错误: ' . $e->getMessage()];
    }

    // 7. 模拟前端调用方式：带 Bearer token（读取前端实际使用的 localStorage key）
    try {
        $token = \Illuminate\Support\Facades\Auth::guard('sanctum')->user()?->currentAccessToken()?->token ?? 'none';
        $results[] = ['test' => '当前认证方式', 'result' => $token !== 'none' ? '✅ 已认证' : 'ℹ️ 未登录（debug 页面无需登录）'];
    } catch (\Throwable $e) {
        $results[] = ['test' => '认证检查', 'result' => 'ℹ️ 未登录: ' . $e->getMessage()];
    }

    return response()->json(['results' => $results]);
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
