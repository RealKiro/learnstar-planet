<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        Schema::create('ai_settings', function (Blueprint $table) {
            $table->id();
            $table->foreignId('school_id')->constrained()->cascadeOnDelete();
            $table->boolean('enabled')->default(false);
            $table->string('provider')->default('openai');
            $table->string('api_key')->nullable();
            $table->string('api_base')->nullable();
            $table->string('model')->default('gpt-3.5-turbo');
            $table->integer('max_tokens')->default(2000);
            $table->integer('tokens_used')->default(0);
            $table->integer('tokens_limit')->default(1000000);
            $table->timestamps();
        });

        Schema::create('ai_conversations', function (Blueprint $table) {
            $table->id();
            $table->foreignId('school_id')->constrained()->cascadeOnDelete();
            $table->foreignId('class_id')->nullable()->constrained('class_rooms')->nullOnDelete();
            $table->foreignId('student_id')->nullable()->constrained()->nullOnDelete();
            $table->string('student_name')->nullable();
            $table->text('question');
            $table->text('answer')->nullable();
            $table->integer('tokens_used')->default(0);
            $table->string('status')->default('completed');
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('ai_conversations');
        Schema::dropIfExists('ai_settings');
    }
};
