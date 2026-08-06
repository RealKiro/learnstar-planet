<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 剥离教师外键：删除教师账号时不级联删除学生数据。
 *
 * - scores.given_by：教师删除后置空（积分记录保留，属于学生）
 * - notices.published_by：教师删除后置空（班级通知保留）
 *
 * 原为 cascadeOnDelete，改为 nullable + nullOnDelete。
 */
return new class extends Migration {
    public function up(): void
    {
        // SQLite 不支持按名 dropForeign，且默认不强制外键约束，跳过；
        // MySQL / PostgreSQL 等才需要改外键为 nullOnDelete。
        if (\Illuminate\Support\Facades\DB::connection()->getDriverName() === 'sqlite') {
            return;
        }
        $this->detach('scores', 'given_by');
        $this->detach('notices', 'published_by');
    }

    public function down(): void
    {
        if (\Illuminate\Support\Facades\DB::connection()->getDriverName() === 'sqlite') {
            return;
        }
        // 回滚到 cascade（不常用，谨慎）
        $this->reattach('scores', 'given_by');
        $this->reattach('notices', 'published_by');
    }

    private function detach(string $table, string $column): void
    {
        if (!Schema::hasTable($table) || !Schema::hasColumn($table, $column)) {
            return;
        }
        $fk = "{$table}_{$column}_foreign";
        Schema::table($table, function (Blueprint $t) use ($fk, $column): void {
            $t->dropForeign($fk);
            $t->unsignedBigInteger($column)->nullable()->change();
        });
        Schema::table($table, function (Blueprint $t) use ($column): void {
            $t->foreign($column)->references('id')->on('users')->nullOnDelete();
        });
    }

    private function reattach(string $table, string $column): void
    {
        if (!Schema::hasTable($table) || !Schema::hasColumn($table, $column)) {
            return;
        }
        $fk = "{$table}_{$column}_foreign";
        Schema::table($table, function (Blueprint $t) use ($fk, $column): void {
            $t->dropForeign($fk);
            $t->unsignedBigInteger($column)->nullable(false)->change();
        });
        Schema::table($table, function (Blueprint $t) use ($column): void {
            $t->foreign($column)->references('id')->on('users')->cascadeOnDelete();
        });
    }
};
