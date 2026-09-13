<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 课表修改申请：教师提交的课表变更需经学校管理员审核通过后才生效
 *
 * - payload 存提交时的完整课表快照（subjects / periods / entries），批准时一次性应用
 * - 审核通过 / 驳回均记录审核人与备注，不产生部分生效
 */
return new class () extends Migration {
    public function up(): void
    {
        if (Schema::hasTable('timetable_change_requests')) {
            return;
        }

        Schema::create('timetable_change_requests', function (Blueprint $t) {
            $t->id();
            $t->unsignedBigInteger('school_id')->index();
            $t->unsignedBigInteger('class_id')->index();
            $t->unsignedBigInteger('requested_by')->comment('提交教师 user id');
            $t->json('payload')->comment('完整课表快照：{subjects,periods,entries}');
            $t->unsignedInteger('entry_count')->default(0)->comment('申请的排课条数（列表摘要用）');
            $t->string('status', 20)->default('pending')->comment('pending/approved/rejected');
            $t->unsignedBigInteger('reviewed_by')->nullable()->comment('审核管理员 user id');
            $t->string('review_note', 200)->nullable()->comment('审核备注（驳回原因等）');
            $t->timestamp('reviewed_at')->nullable();
            $t->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('timetable_change_requests');
    }
};
