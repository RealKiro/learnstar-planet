<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    /**
     * 给 users 表加 subject 字段（教师所教科目）
     */
    public function up(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->string('subject', 50)->nullable()->after('nickname')->comment('科目，如：语文、数学、英语');
        });
    }

    public function down(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn('subject');
        });
    }
};
