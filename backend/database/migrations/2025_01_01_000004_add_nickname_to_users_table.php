<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    /**
     * 给 users 表加 nickname 和 avatar_path 字段（教师/家长有昵称和头像）
     * 昵称默认为中文转拼音（如果装了 overtrue/pinyin）
     * 头像默认为 null（教师/家长可上传或留空）
     */
    public function up(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->string('nickname', 80)->nullable()->after('name')->comment('昵称，默认 = 名字的拼音');
            $table->string('avatar_path')->nullable()->change();
        });
    }

    public function down(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn('nickname');
        });
    }
};
