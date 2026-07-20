<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class () extends Migration {
    public function up(): void
    {
        // ===== 钱包（多币种余额） =====
        Schema::create('wallets', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->constrained()->cascadeOnDelete();
            $table->string('currency_type'); // science, reading, class_point
            $table->integer('balance')->default(0);
            $table->timestamps();

            $table->unique(['student_id', 'currency_type']);
        });

        // ===== 汇率配置（学校级） =====
        Schema::create('exchange_rates', function (Blueprint $table) {
            $table->id();
            $table->foreignId('school_id')->constrained()->cascadeOnDelete();
            $table->string('from_currency'); // score
            $table->string('to_currency');   // science, reading, class_point
            $table->decimal('rate', 8, 2)->default(1); // 10 score = 1 coin → rate=0.1
            $table->boolean('is_active')->default(true);
            $table->timestamps();

            $table->unique(['school_id', 'from_currency', 'to_currency']);
        });

        // ===== 兑换日志 =====
        Schema::create('exchange_logs', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->constrained()->cascadeOnDelete();
            $table->string('from_currency');
            $table->string('to_currency');
            $table->integer('from_amount');
            $table->integer('to_amount');
            $table->foreignId('operated_by')->nullable()->constrained('users')->nullOnDelete();
            $table->timestamps();

            $table->index(['student_id', 'created_at']);
        });

        // ===== 给 shop_items 添加币种和活动标签字段 =====
        Schema::table('shop_items', function (Blueprint $table) {
            $table->string('currency_type')->default('score')->after('cost_score');
            $table->string('event_tag')->nullable()->after('category');
        });
    }

    public function down(): void
    {
        Schema::table('shop_items', function (Blueprint $table) {
            $table->dropColumn(['currency_type', 'event_tag']);
        });

        Schema::dropIfExists('exchange_logs');
        Schema::dropIfExists('exchange_rates');
        Schema::dropIfExists('wallets');
    }
};
