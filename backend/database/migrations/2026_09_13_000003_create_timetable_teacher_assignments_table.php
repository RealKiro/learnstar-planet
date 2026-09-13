<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 任课表：班级 × 科目 → 教师。
 *
 * 全校智能排课的依据：同一教师同一时段只能出现在一个班级。
 * 科目按名称存（与课表模块「科目名交互」口径一致）。
 */
return new class extends Migration
{
    public function up(): void
    {
        if (Schema::hasTable('timetable_teacher_assignments')) {
            return;
        }

        Schema::create('timetable_teacher_assignments', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('school_id')->index();
            $table->unsignedBigInteger('class_id')->index();
            $table->string('subject_name', 50);
            $table->string('teacher_name', 50);
            $table->timestamps();

            $table->unique(['class_id', 'subject_name']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('timetable_teacher_assignments');
    }
};
