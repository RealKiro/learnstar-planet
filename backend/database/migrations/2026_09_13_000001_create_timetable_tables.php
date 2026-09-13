<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 课表模块：科目 / 节次 / 排课
 *
 * 设计要点：
 * - subjects 与 class_periods 为「学校级」共享资源（全校一套科目与作息时间）
 * - timetable_entries 按班级存排课格子，week_type 支持单双周（all/odd/even），
 *   与 CSES 交换格式的 weeks 字段一一对应，便于导出给 ClassIsland 导入
 */
return new class () extends Migration {
    public function up(): void
    {
        if (!Schema::hasTable('subjects')) {
            Schema::create('subjects', function (Blueprint $t) {
                $t->id();
                $t->unsignedBigInteger('school_id')->index();
                $t->string('name', 50);
                $t->string('simplified_name', 20)->nullable();
                $t->string('color', 20)->nullable();
                $t->unsignedInteger('sort_order')->default(0);
                $t->timestamps();
                $t->unique(['school_id', 'name'], 'subjects_school_name_unique');
            });
        }

        if (!Schema::hasTable('class_periods')) {
            Schema::create('class_periods', function (Blueprint $t) {
                $t->id();
                $t->unsignedBigInteger('school_id')->index();
                $t->unsignedInteger('period_index');
                $t->string('name', 30);
                $t->string('start_time', 8);
                $t->string('end_time', 8);
                $t->timestamps();
                $t->unique(['school_id', 'period_index'], 'class_periods_school_index_unique');
            });
        }

        if (!Schema::hasTable('timetable_entries')) {
            Schema::create('timetable_entries', function (Blueprint $t) {
                $t->id();
                $t->unsignedBigInteger('class_id')->index();
                $t->unsignedInteger('weekday');
                $t->unsignedInteger('period_index');
                $t->string('week_type', 10)->default('all');
                $t->unsignedBigInteger('subject_id')->nullable();
                $t->string('teacher_name', 50)->nullable();
                $t->string('room', 50)->nullable();
                $t->timestamps();
                $t->unique(['class_id', 'weekday', 'period_index', 'week_type'], 'timetable_slot_unique');
            });
        }
    }

    public function down(): void
    {
        Schema::dropIfExists('timetable_entries');
        Schema::dropIfExists('class_periods');
        Schema::dropIfExists('subjects');
    }
};
