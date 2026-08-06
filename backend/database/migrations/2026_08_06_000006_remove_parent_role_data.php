<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

/**
 * 清理家长角色数据（家长端已从代码移除）：
 * 1. 解除 students.parent_id 对学生家长的绑定（置 NULL）
 * 2. 删除 role='parent' 的 user 记录
 *
 * 外键联动（均为 nullOnDelete / cascadeOnDelete，无需手工）：
 * - students.parent_id            → nullOnDelete（自动置 NULL）
 * - third_party_bindings.user_id  → cascadeOnDelete（自动删绑定）
 * - personal_access_tokens        → cascadeOnDelete（自动删 token）
 */
return new class extends Migration
{
    public function up(): void
    {
        DB::transaction(function () {
            $parentIds = DB::table('users')->where('role', 'parent')->pluck('id');

            // 防御性先解除学生绑定（虽然 nullOnDelete 也会处理，显式执行更稳）
            if ($parentIds->isNotEmpty()) {
                DB::table('students')->whereIn('parent_id', $parentIds)->update(['parent_id' => null]);
            }

            DB::table('users')->where('role', 'parent')->delete();
        });
    }

    public function down(): void
    {
        // 数据删除不可逆，down 不恢复
    }
};
