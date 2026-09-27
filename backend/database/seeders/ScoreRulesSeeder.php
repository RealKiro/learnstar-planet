<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\School;
use App\Services\ScoreRuleService;
use Illuminate\Database\Seeder;

/**
 * 默认积分规则播种。
 *
 * 与运行时（教师端首次访问时的惰性补齐）**同源**：都取
 * {@see ScoreRuleService::DEFAULT_RULES}，因此 db:seed 与界面首访的结果完全一致。
 *
 * 旧版本此处另维护了一份 20 条、分类为中文（课堂表现/作业/品德…）的独立清单，
 * 与前端分类标签体系冲突，已删除改用唯一真源。
 */
class ScoreRulesSeeder extends Seeder
{
    public function run(): void
    {
        foreach (School::pluck('id') as $schoolId) {
            ScoreRuleService::ensureDefaultsForSchool((int) $schoolId);
        }
    }
}
