<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        // ===== 学校 =====
        Schema::create('schools', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('code')->unique();          // 学校唯一代码
            $table->string('address')->nullable();
            $table->string('contact_phone')->nullable();
            $table->string('contact_email')->nullable();
            $table->string('logo_path')->nullable();
            $table->json('settings')->nullable();
            $table->string('status')->default('active');
            $table->timestamps();
            $table->softDeletes();
        });

        // ===== 用户（教师/家长/管理员） =====
        Schema::create('users', function (Blueprint $table) {
            $table->id();
            $table->foreignId('school_id')->constrained()->cascadeOnDelete();
            $table->string('role');                     // school_admin / teacher / parent
            $table->string('username')->unique();       // 管理员分配的账号
            $table->string('password');
            $table->string('name');
            $table->string('avatar_path')->nullable();
            $table->string('phone')->nullable();
            $table->string('email')->nullable();
            $table->boolean('password_changed')->default(false);
            $table->string('status')->default('active');
            $table->timestamp('last_login_at')->nullable();
            $table->rememberToken();
            $table->timestamps();
            $table->softDeletes();
        });

        // ===== 第三方账号绑定 =====
        Schema::create('third_party_bindings', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->constrained()->cascadeOnDelete();
            $table->string('platform');                 // wechat / wechat_work / qq / renren
            $table->string('platform_id');              // openid / userid
            $table->string('platform_union_id')->nullable(); // unionid
            $table->string('platform_nick')->nullable();
            $table->string('platform_avatar')->nullable();
            $table->timestamp('verified_at')->nullable();
            $table->timestamps();

            $table->unique(['platform', 'platform_id']);
            $table->unique(['platform_union_id'], 'tpb_union_id_unique')->where('platform_union_id', '!=', null);
        });

        // ===== 班级 =====
        Schema::create('class_rooms', function (Blueprint $table) {
            $table->id();
            $table->foreignId('school_id')->constrained()->cascadeOnDelete();
            $table->string('name');
            $table->string('grade')->nullable();
            $table->string('year')->nullable();
            $table->foreignId('teacher_id')->nullable()->constrained('users')->nullOnDelete();
            $table->integer('max_students')->default(0); // 0 = 不限制（全免费）
            $table->json('settings')->nullable();
            $table->string('status')->default('active');
            $table->timestamps();
            $table->softDeletes();
        });

        // ===== 学生 =====
        Schema::create('students', function (Blueprint $table) {
            $table->id();
            $table->foreignId('class_id')->constrained('class_rooms')->cascadeOnDelete();
            $table->foreignId('parent_id')->nullable()->constrained('users')->nullOnDelete();
            $table->string('name');
            $table->string('student_no')->nullable();
            $table->string('avatar_path')->nullable();
            $table->integer('total_score')->default(0);
            $table->string('status')->default('active');
            $table->timestamps();
            $table->softDeletes();

            $table->index(['class_id', 'total_score']);
        });

        // ===== 宠物 =====
        Schema::create('pets', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->unique()->constrained()->cascadeOnDelete();
            $table->foreignId('class_id')->constrained('class_rooms')->cascadeOnDelete();
            $table->string('name')->default('星尘');
            $table->string('type')->default('stellar_cat');
            $table->integer('level')->default(0);
            $table->integer('experience')->default(0);
            $table->integer('mood')->default(80);
            $table->json('accessories')->nullable();
            $table->timestamp('last_fed_at')->nullable();
            $table->timestamps();
        });

        // ===== 积分规则 =====
        Schema::create('score_rules', function (Blueprint $table) {
            $table->id();
            $table->foreignId('class_id')->nullable()->constrained('class_rooms')->cascadeOnDelete();
            $table->foreignId('school_id')->nullable()->constrained()->cascadeOnDelete();
            $table->string('name');
            $table->integer('amount');
            $table->string('category')->default('academic');
            $table->boolean('is_positive')->default(true);
            $table->boolean('is_active')->default(true);
            $table->integer('sort_order')->default(0);
            $table->timestamps();
        });

        // ===== 积分记录 =====
        Schema::create('scores', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->constrained()->cascadeOnDelete();
            $table->foreignId('class_id')->constrained('class_rooms')->cascadeOnDelete();
            $table->foreignId('score_rule_id')->nullable()->constrained()->nullOnDelete();
            $table->integer('amount');
            $table->string('reason')->nullable();
            $table->foreignId('given_by')->constrained('users')->cascadeOnDelete();
            $table->timestamps();

            $table->index(['student_id', 'created_at']);
            $table->index(['class_id', 'created_at']);
        });

        // ===== 积分日志 =====
        Schema::create('score_logs', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->constrained()->cascadeOnDelete();
            $table->foreignId('score_id')->constrained()->cascadeOnDelete();
            $table->integer('balance_before');
            $table->integer('balance_after');
            $table->string('description')->nullable();
            $table->timestamps();
        });

        // ===== 班级通知 =====
        Schema::create('notices', function (Blueprint $table) {
            $table->id();
            $table->foreignId('class_id')->nullable()->constrained('class_rooms')->cascadeOnDelete();
            $table->foreignId('school_id')->nullable()->constrained()->cascadeOnDelete();
            $table->string('title');
            $table->text('content');
            $table->string('type')->default('info');
            $table->foreignId('published_by')->constrained('users')->cascadeOnDelete();
            $table->boolean('is_published')->default(false);
            $table->timestamp('published_at')->nullable();
            $table->timestamps();
        });

        // ===== 积分商城 =====
        Schema::create('shop_items', function (Blueprint $table) {
            $table->id();
            $table->foreignId('class_id')->constrained('class_rooms')->cascadeOnDelete();
            $table->string('name');
            $table->text('description')->nullable();
            $table->string('category')->default('physical');
            $table->integer('cost_score');
            $table->integer('stock')->default(0); // 0 = 不限制
            $table->string('image_path')->nullable();
            $table->boolean('is_active')->default(true);
            $table->timestamps();
        });

        // ===== 兑换记录 =====
        Schema::create('shop_redemptions', function (Blueprint $table) {
            $table->id();
            $table->foreignId('student_id')->constrained()->cascadeOnDelete();
            $table->foreignId('shop_item_id')->constrained()->cascadeOnDelete();
            $table->foreignId('class_id')->constrained('class_rooms')->cascadeOnDelete();
            $table->string('status')->default('pending');
            $table->foreignId('approved_by')->nullable()->constrained('users')->nullOnDelete();
            $table->timestamp('delivered_at')->nullable();
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('shop_redemptions');
        Schema::dropIfExists('shop_items');
        Schema::dropIfExists('notices');
        Schema::dropIfExists('score_logs');
        Schema::dropIfExists('scores');
        Schema::dropIfExists('score_rules');
        Schema::dropIfExists('pets');
        Schema::dropIfExists('students');
        Schema::dropIfExists('class_rooms');
        Schema::dropIfExists('third_party_bindings');
        Schema::dropIfExists('users');
        Schema::dropIfExists('schools');
    }
};
