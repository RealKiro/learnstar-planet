<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::table('class_rooms', function (Blueprint $table) {
            // 班级码（LS + 年级 + 班号）跨校可能撞码，唯一索引改为校维度
            if (Schema::hasIndex('class_rooms', 'class_rooms_display_code_unique')) {
                $table->dropUnique('class_rooms_display_code_unique');
            }
            $table->unique(['school_id', 'display_code']);
        });
    }

    public function down(): void
    {
        Schema::table('class_rooms', function (Blueprint $table) {
            if (Schema::hasIndex('class_rooms', 'class_rooms_school_id_display_code_unique')) {
                $table->dropUnique('class_rooms_school_id_display_code_unique');
            }
            $table->unique('display_code');
        });
    }
};
