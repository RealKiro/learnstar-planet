<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        // ===== broadcasts 广播表 =====
        Schema::create('broadcasts', function ($table) {
            $table->id();
            $table->foreignId('school_id')->constrained('schools')->onDelete('cascade');
            $table->foreignId('class_id')->constrained('class_rooms')->onDelete('cascade');
            $table->foreignId('teacher_id')->constrained('users')->onDelete('cascade');
            $table->text('content');
            $table->enum('type', ['banner', 'popup', 'fullscreen'])->default('banner');
            $table->boolean('voice_enabled')->default(true);
            $table->boolean('loop_enabled')->default(false);
            $table->integer('display_seconds')->default(10);
            $table->enum('status', ['pending', 'sent', 'received', 'expired'])->default('pending');
            $table->timestamp('sent_at')->nullable();
            $table->timestamps();
            $table->index(['class_id', 'status']);
        });

        // ===== attendances 考勤表 =====
        Schema::create('attendances', function ($table) {
            $table->id();
            $table->foreignId('class_id')->constrained('class_rooms')->onDelete('cascade');
            $table->foreignId('student_id')->constrained('students')->onDelete('cascade');
            $table->foreignId('teacher_id')->constrained('users')->onDelete('cascade');
            $table->date('date');
            $table->enum('status', ['present', 'late', 'leave', 'absent'])->default('present');
            $table->timestamp('sign_in_at')->nullable();
            $table->string('remark')->nullable();
            $table->timestamps();
            $table->unique(['class_id', 'student_id', 'date']);
            $table->index(['class_id', 'date']);
        });

        // ===== homework_collections 作业收取表 =====
        Schema::create('homework_collections', function ($table) {
            $table->id();
            $table->foreignId('class_id')->constrained('class_rooms')->onDelete('cascade');
            $table->foreignId('teacher_id')->constrained('users')->onDelete('cascade');
            $table->string('title', 200);
            $table->timestamp('deadline')->nullable();
            $table->text('description')->nullable();
            $table->json('submit_types')->nullable();  // ["qrcode","upload","online"]
            $table->enum('status', ['active', 'closed', 'archived'])->default('active');
            $table->string('qr_code_token', 100)->nullable()->unique();
            $table->timestamps();
            $table->index(['class_id', 'status']);
        });

        // ===== homework_submissions 作业提交表 =====
        Schema::create('homework_submissions', function ($table) {
            $table->id();
            $table->foreignId('homework_id')->constrained('homework_collections')->onDelete('cascade');
            $table->foreignId('student_id')->constrained('students')->onDelete('cascade');
            $table->timestamp('submitted_at')->nullable();
            $table->text('content')->nullable();
            $table->json('file_urls')->nullable();
            $table->timestamps();
            $table->unique(['homework_id', 'student_id']);
        });

        // ===== question_banks 题库表 (must be created before quizzes) =====
        Schema::create('question_banks', function ($table) {
            $table->id();
            $table->foreignId('school_id')->nullable()->constrained('schools')->onDelete('cascade');
            $table->foreignId('teacher_id')->constrained('users')->onDelete('cascade');
            $table->string('title', 200);
            $table->string('subject', 50);
            $table->integer('question_count')->default(0);
            $table->integer('used_count')->default(0);
            $table->boolean('is_public')->default(false);  // 是否全校共享
            $table->timestamps();
            $table->index(['teacher_id', 'subject']);
        });

        // ===== questions 题目表 (depends on question_banks) =====
        Schema::create('questions', function ($table) {
            $table->id();
            $table->foreignId('question_bank_id')->constrained('question_banks')->onDelete('cascade');
            $table->enum('type', ['single', 'multiple', 'fill', 'judge', 'short']);
            $table->text('content');
            $table->json('options')->nullable();  // 选择题选项
            $table->text('answer');  // 正确答案
            $table->decimal('score', 4, 1)->default(1.0);
            $table->text('explanation')->nullable();
            $table->integer('sort_order')->default(0);
            $table->timestamps();
        });

        // ===== quizzes 答题表 (depends on question_banks) =====
        Schema::create('quizzes', function ($table) {
            $table->id();
            $table->foreignId('class_id')->constrained('class_rooms')->onDelete('cascade');
            $table->foreignId('teacher_id')->constrained('users')->onDelete('cascade');
            $table->foreignId('question_bank_id')->nullable()->constrained('question_banks')->onDelete('set null');
            $table->string('title', 200);
            $table->integer('time_limit')->default(0);  // 0=不限时
            $table->boolean('auto_grade')->default(true);
            $table->boolean('realtime_stats')->default(true);
            $table->boolean('anonymous')->default(false);
            $table->enum('status', ['draft', 'active', 'grading', 'finished'])->default('draft');
            $table->timestamp('started_at')->nullable();
            $table->timestamp('ended_at')->nullable();
            $table->timestamps();
            $table->index(['class_id', 'status']);
        });

        // ===== quiz_submissions 答题提交表 (depends on quizzes) =====
        Schema::create('quiz_submissions', function ($table) {
            $table->id();
            $table->foreignId('quiz_id')->constrained('quizzes')->onDelete('cascade');
            $table->foreignId('student_id')->constrained('students')->onDelete('cascade');
            $table->json('answers')->nullable();
            $table->decimal('score', 5, 1)->nullable();
            $table->timestamp('submitted_at')->nullable();
            $table->timestamps();
            $table->unique(['quiz_id', 'student_id']);
        });

        // ===== grades 成绩表 =====
        Schema::create('grades', function ($table) {
            $table->id();
            $table->foreignId('class_id')->constrained('class_rooms')->onDelete('cascade');
            $table->foreignId('student_id')->constrained('students')->onDelete('cascade');
            $table->foreignId('teacher_id')->constrained('users')->onDelete('cascade');
            $table->string('exam_name', 100);
            $table->string('subject', 50);
            $table->decimal('score', 5, 1);
            $table->integer('rank_in_class')->nullable();
            $table->integer('rank_change')->default(0);
            $table->timestamps();
            $table->unique(['class_id', 'student_id', 'exam_name', 'subject']);
            $table->index(['class_id', 'exam_name', 'subject']);
        });
    }

    public function down(): void
    {
        $tables = [
            'grades',
            'quiz_submissions', 'quizzes',
            'questions', 'question_banks',
            'homework_submissions', 'homework_collections',
            'attendances', 'broadcasts',
        ];
        foreach ($tables as $table) {
            Schema::dropIfExists($table);
        }
    }
};

