<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 教师不可用时段（学校级）：
 * 自动排课（全校智能排课）时视为该教师已占用，不会把课排进这些格子。
 */
return new class extends Migration
{
    public function up(): void
    {
        if (Schema::hasTable('timetable_teacher_unavailabilities')) {
            return;
        }

        Schema::create('timetable_teacher_unavailabilities', function (Blueprint $table): void {
            $table->id();
            $table->unsignedBigInteger('school_id')->index();
            $table->string('teacher_name', 50);
            $table->unsignedTinyInteger('weekday');       // 1-7（周一到周日）
            $table->unsignedTinyInteger('period_index');  // 节次
            $table->timestamps();

            $table->unique(['school_id', 'teacher_name', 'weekday', 'period_index'], 'tt_unavail_unique');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('timetable_teacher_unavailabilities');
    }
};
