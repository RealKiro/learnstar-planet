<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\ScoreRule;
use Illuminate\Database\Seeder;

class ScoreRulesSeeder extends Seeder
{
    public function run(): void
    {
        $defaultRules = [
            // Positive behaviors (+score, +pet exp)
            ['name' => '课堂积极发言', 'amount' => 2, 'category' => '课堂表现', 'is_positive' => true],
            ['name' => '作业完成优秀', 'amount' => 3, 'category' => '作业', 'is_positive' => true],
            ['name' => '帮助同学', 'amount' => 3, 'category' => '品德', 'is_positive' => true],
            ['name' => '课堂纪律好', 'amount' => 1, 'category' => '课堂表现', 'is_positive' => true],
            ['name' => '按时完成作业', 'amount' => 2, 'category' => '作业', 'is_positive' => true],
            ['name' => '卫生打扫认真', 'amount' => 2, 'category' => '劳动', 'is_positive' => true],
            ['name' => '礼貌问候', 'amount' => 1, 'category' => '品德', 'is_positive' => true],
            ['name' => '早读认真', 'amount' => 1, 'category' => '日常', 'is_positive' => true],
            ['name' => '考试进步明显', 'amount' => 5, 'category' => '学业', 'is_positive' => true],
            ['name' => '获得校级奖项', 'amount' => 10, 'category' => '荣誉', 'is_positive' => true],

            // Negative behaviors (-score, -pet exp)
            ['name' => '课堂讲小话', 'amount' => -1, 'category' => '课堂表现', 'is_positive' => false],
            ['name' => '未完成作业', 'amount' => -2, 'category' => '作业', 'is_positive' => false],
            ['name' => '迟到', 'amount' => -2, 'category' => '考勤', 'is_positive' => false],
            ['name' => '课堂玩手机', 'amount' => -3, 'category' => '课堂表现', 'is_positive' => false],
            ['name' => '乱扔垃圾', 'amount' => -1, 'category' => '卫生', 'is_positive' => false],
            ['name' => '说脏话', 'amount' => -2, 'category' => '品德', 'is_positive' => false],
            ['name' => '打架斗殴', 'amount' => -10, 'category' => '品德', 'is_positive' => false],
            ['name' => '破坏公物', 'amount' => -5, 'category' => '品德', 'is_positive' => false],
            ['name' => '考试作弊', 'amount' => -15, 'category' => '学业', 'is_positive' => false],
            ['name' => '早退', 'amount' => -3, 'category' => '考勤', 'is_positive' => false],
        ];

        foreach ($defaultRules as $i => $rule) {
            ScoreRule::firstOrCreate(
                ['name' => $rule['name'], 'class_id' => null, 'school_id' => null],
                [
                    'amount' => $rule['amount'],
                    'category' => $rule['category'],
                    'is_positive' => $rule['is_positive'],
                    'is_active' => true,
                    'sort_order' => $i,
                ]
            );
        }
    }
}

