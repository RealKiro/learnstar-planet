<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

/**
 * 兑换体系按实际价值定价(防通胀):
 * 1. 汇率改为 2:1(2 积分 = 1 科学币/体育币/读书币)
 * 2. 商品积分按实际价值重定: 小商品 ≈100(1周, 日最高20×5天)、大商品 ≈200(2周)
 */
return new class extends Migration
{
    public function up(): void
    {
        DB::transaction(function () {
            // ===== 1) 汇率 2:1 =====
            DB::table('exchange_rates')->where('from_currency', 'score')->update(['rate' => 0.5]);

            // ===== 2) 商品按实际价值重定积分 =====
            $prices = [
                // 小商品(≈100, 1 周可攒够)
                '铅笔' => 100,
                '橡皮擦' => 100,
                '草稿纸' => 100,
                '免罚站一次' => 100,
                '免罚跑步一次' => 100,
                // 中商品(120~150)
                '便利贴' => 120,
                '黑色圆珠笔' => 150,
                '蓝色圆珠笔' => 150,
                '红色圆珠笔' => 150,
                '香蕉' => 150,
                '免做卫生一次' => 150,
                // 大商品(≈200, 2 周攒够)
                '练习本' => 180,
                '苹果' => 180,
                '饮料' => 200,
                '牛奶' => 200,
                '水果' => 200,
                '3D打印作品' => 200,
                '集体观影' => 200,
                '免作业一次' => 200,
            ];
            foreach ($prices as $name => $cost) {
                DB::table('shop_items')->where('name', $name)->update(['cost_score' => $cost]);
            }

            // ===== 3) 确保行为奖励存在(school 级,新价格) =====
            $rewards = [
                ['name' => '集体观影', 'description' => '全班集体观影一次', 'category' => 'activity', 'cost_score' => 200],
                ['name' => '免作业一次', 'description' => '免交一次作业', 'category' => 'privilege', 'cost_score' => 200],
                ['name' => '免做卫生一次', 'description' => '免除一次值日卫生', 'category' => 'privilege', 'cost_score' => 150],
                ['name' => '免罚站一次', 'description' => '免除一次罚站', 'category' => 'privilege', 'cost_score' => 100],
                ['name' => '免罚跑步一次', 'description' => '免除一次罚跑步', 'category' => 'privilege', 'cost_score' => 100],
            ];
            $schoolIds = DB::table('schools')->pluck('id');
            foreach ($schoolIds as $schoolId) {
                foreach ($rewards as $r) {
                    DB::table('shop_items')->updateOrInsert(
                        ['school_id' => $schoolId, 'class_id' => null, 'name' => $r['name']],
                        [
                            'description' => $r['description'],
                            'category' => $r['category'],
                            'cost_score' => $r['cost_score'],
                            'currency_type' => 'score',
                            'stock' => 0,
                            'is_active' => true,
                            'updated_at' => now(),
                        ]
                    );
                }
            }
        });
    }

    public function down(): void
    {
        // 数据调整不可逆，down 不恢复
    }
};
