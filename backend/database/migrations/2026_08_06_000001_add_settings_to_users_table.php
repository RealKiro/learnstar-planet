<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 补充 users.settings 列（生产库可能缺失，导致 updateTeacher 保存 personal_role 时
 * SQLSTATE 42S22 Column not found: 'settings'）。
 *
 * 幂等：仅在列不存在时添加，避免在已修复的库上重复执行报错。
 */
return new class () extends Migration {
    public function up(): void
    {
        Schema::table('users', function (Blueprint $table) {
            if (!Schema::hasColumn('users', 'settings')) {
                $table->json('settings')->nullable()->comment('用户级配置 JSON');
            }
        });
    }

    public function down(): void
    {
        Schema::table('users', function (Blueprint $table) {
            if (Schema::hasColumn('users', 'settings')) {
                $table->dropColumn('settings');
            }
        });
    }
};
