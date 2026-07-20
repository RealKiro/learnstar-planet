<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('class_rooms', function (Blueprint $table) {
            $table->string('display_code', 12)
                ->nullable()
                ->unique()
                ->after('max_students')
                ->comment('班级大屏码，用于触摸屏免账号登录');
            $table->timestamp('display_code_updated_at')
                ->nullable()
                ->after('display_code')
                ->comment('班级码刷新时间');
        });
    }

    public function down(): void
    {
        Schema::table('class_rooms', function (Blueprint $table) {
            $table->dropColumn(['display_code', 'display_code_updated_at']);
        });
    }
};
