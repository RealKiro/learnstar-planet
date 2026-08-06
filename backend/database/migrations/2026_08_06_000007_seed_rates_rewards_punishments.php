<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

/**
 * 兑换体系初始化：
 * 1. 所有汇率统一为 1:1（积分→科学币/读书币/体育币）
 * 2. 为每个学校播种默认奖励商品（行为奖励 + 实物奖励，school 级 → 自动同步到所有班级）
 * 3. 为每个学校播种默认惩罚规则（负分，school 级）
 */
return new class extends Migration
{
    public function up(): void
    {
        DB::transaction(function () {
            // ===== 1) 汇率统一 1:1 =====
            DB::table('exchange_rates')->update(['rate' => 1]);

            $schoolIds = DB::table('schools')->pluck('id');
            if ($schoolIds->isEmpty()) {
                return;
            }

            // ===== 2) 默认奖励商品（school 级，class_id=null 对所有班级可见） =====
            $rewards = [
                // 行为奖励
                ['name' => '集体观影', 'description' => '全班集体观影一次', 'category' => 'activity', 'cost_score' => 80],
                ['name' => '免作业一次', 'description' => '免交一次作业', 'category' => 'privilege', 'cost_score' => 30],
                ['name' => '免做卫生一次', 'description' => '免除一次值日卫生', 'category' => 'privilege', 'cost_score' => 20],
                ['name' => '免罚站一次', 'description' => '免除一次罚站', 'category' => 'privilege', 'cost_score' => 15],
                ['name' => '免罚跑步一次', 'description' => '免除一次罚跑步', 'category' => 'privilege', 'cost_score' => 15],
                // 实物奖励
                ['name' => '铅笔', 'description' => '标准 HB 铅笔一支', 'category' => 'stationery', 'cost_score' => 20],
                ['name' => '黑色圆珠笔', 'description' => '黑色圆珠笔一支', 'category' => 'stationery', 'cost_score' => 25],
                ['name' => '蓝色圆珠笔', 'description' => '蓝色圆珠笔一支', 'category' => 'stationery', 'cost_score' => 25],
                ['name' => '红色圆珠笔', 'description' => '红色圆珠笔一支', 'category' => 'stationery', 'cost_score' => 25],
                ['name' => '橡皮擦', 'description' => '白色橡皮擦一块', 'category' => 'stationery', 'cost_score' => 15],
                ['name' => '便利贴', 'description' => '彩色便利贴一本', 'category' => 'stationery', 'cost_score' => 18],
                ['name' => '草稿纸', 'description' => '草稿纸一本', 'category' => 'stationery', 'cost_score' => 10],
                ['name' => '练习本', 'description' => '练习本一本', 'category' => 'stationery', 'cost_score' => 30],
                ['name' => '饮料', 'description' => '瓶装饮料一瓶', 'category' => 'food', 'cost_score' => 40],
                ['name' => '水果', 'description' => '新鲜水果一份', 'category' => 'food', 'cost_score' => 35],
                ['name' => '3D打印作品', 'description' => '3D 打印小作品一件', 'category' => 'physical', 'cost_score' => 100],
            ];
            foreach ($schoolIds as $schoolId) {
                foreach ($rewards as $r) {
                    $exists = DB::table('shop_items')
                        ->where('school_id', $schoolId)
                        ->whereNull('class_id')
                        ->where('name', $r['name'])
                        ->exists();
                    if ($exists) {
                        continue;
                    }
                    DB::table('shop_items')->insert([
                        'school_id' => $schoolId,
                        'class_id' => null,
                        'name' => $r['name'],
                        'description' => $r['description'],
                        'category' => $r['category'],
                        'cost_score' => $r['cost_score'],
                        'currency_type' => 'score',
                        'stock' => 0,
                        'is_active' => true,
                        'created_at' => now(),
                        'updated_at' => now(),
                    ]);
                }
            }

            // ===== 3) 默认惩罚规则（负分，school 级） =====
            $punishments = [
                ['name' => '迟到', 'amount' => -2],
                ['name' => '上课讲话', 'amount' => -1],
                ['name' => '未交作业', 'amount' => -2],
                ['name' => '乱扔垃圾', 'amount' => -1],
                ['name' => '追逐打闹', 'amount' => -2],
                ['name' => '说脏话', 'amount' => -1],
                ['name' => '上课睡觉', 'amount' => -1],
                ['name' => '打架', 'amount' => -5],
                ['name' => '罚站', 'amount' => -1],
                ['name' => '罚跑步', 'amount' => -1],
            ];
            foreach ($schoolIds as $schoolId) {
                foreach ($punishments as $p) {
                    $exists = DB::table('score_rules')
                        ->where('school_id', $schoolId)
                        ->whereNull('class_id')
                        ->where('name', $p['name'])
                        ->exists();
                    if ($exists) {
                        continue;
                    }
                    DB::table('score_rules')->insert([
                        'school_id' => $schoolId,
                        'class_id' => null,
                        'name' => $p['name'],
                        'amount' => $p['amount'],
                        'category' => 'discipline',
                        'is_positive' => false,
                        'is_active' => true,
                        'sort_order' => 0,
                        'created_at' => now(),
                        'updated_at' => now(),
                    ]);
                }
            }
        });
    }

    public function down(): void
    {
        // 数据播种不可逆，down 不恢复
    }
};
