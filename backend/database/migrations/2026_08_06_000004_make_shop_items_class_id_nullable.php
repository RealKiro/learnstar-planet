<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 让 shop_items.class_id 可空，支持学校级商品（class_id = null, school_id = 当前校）。
 *
 * 原迁移 class_id 为 NOT NULL，但 createShopItem / adminCreateShopItem 都写 class_id = null
 * 创建学校级商品 → 触发 NOT NULL 约束错误。此处改为 nullable。
 *
 * 注意：改列在 SQLite 上通过重建表实现，Laravel 的 change() 在 SQLite 受支持（doctrine/dbal 或原生）。
 */
return new class extends Migration {
    public function up(): void
    {
        if (Schema::hasTable('shop_items') && Schema::hasColumn('shop_items', 'class_id')) {
            Schema::table('shop_items', function (Blueprint $table): void {
                $table->unsignedBigInteger('class_id')->nullable()->change();
            });
        }
    }

    public function down(): void
    {
        if (Schema::hasTable('shop_items') && Schema::hasColumn('shop_items', 'class_id')) {
            Schema::table('shop_items', function (Blueprint $table): void {
                $table->unsignedBigInteger('class_id')->nullable(false)->change();
            });
        }
    }
};
