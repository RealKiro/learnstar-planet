<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 幂等补齐关键列（生产库可能缺失，导致创建/更新教师报 Column not found）。
 *
 * 覆盖：
 * - users.settings（之前 2026_08_06 已加过，此处幂等兜底）
 * - users.subject / users.grade_team（主教科目/年级团队）
 * - class_room_teachers.subject（班级任教科目）
 *
 * 仅在列不存在时添加，已修复的库重复执行无副作用。
 */
return new class extends Migration {
    public function up(): void
    {
        $this->ensureColumn('users', 'settings', fn (Blueprint $t) => $t->json('settings')->nullable()->comment('用户级配置 JSON'));
        $this->ensureColumn('users', 'subject', fn (Blueprint $t) => $t->string('subject', 50)->nullable()->comment('科目'));
        $this->ensureColumn('users', 'grade_team', fn (Blueprint $t) => $t->string('grade_team', 50)->nullable()->comment('所属年级团队'));

        if (Schema::hasTable('class_room_teachers')) {
            $this->ensureColumn('class_room_teachers', 'subject', fn (Blueprint $t) => $t->string('subject', 50)->nullable()->comment('任教科目'));
        }
    }

    public function down(): void
    {
        // 不动，幂等补列不回滚（历史数据修复性质）
    }

    private function ensureColumn(string $table, string $column, \Closure $definition): void
    {
        if (Schema::hasTable($table) && !Schema::hasColumn($table, $column)) {
            Schema::table($table, function (Blueprint $t) use ($column, $definition): void {
                $definition($t);
            });
        }
    }
};
