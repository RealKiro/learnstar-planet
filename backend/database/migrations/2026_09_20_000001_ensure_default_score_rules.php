<?php

declare(strict_types=1);

use App\Services\ScoreRuleService;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

/**
 * 「安装即可用」收口：默认积分规则补齐 + 历史 discipline 分类清理。
 *
 * 背景：早期迁移 2026_08_06_000007 曾为每个已存在的学校播种 10 条惩罚规则，分类写死 `discipline`。
 *   - 该分类不在前端标签表内 → 规则卡片渲染成裸 slug「discipline」；
 *   - 那批规则只含惩罚、不含奖励 → 学校的规则集合非空，使当时「规则为空才播种默认」的
 *     逻辑永不触发 → 该校**有惩罚、零奖励**，教师端加分区分空白。
 *
 * 本迁移做两件事：
 *   1) 清理上述历史播种痕迹（仅限 category='discipline' 且名称命中当年那 10 条的记录）；
 *   2) 为每个学校增量补齐完整默认规则集（幂等，已存在的同名规则一律不覆盖）。
 *
 * 由 entrypoint.sh 的 `php artisan migrate --force` 在启动时自动执行，
 * 使**既有部署**升级后同样达到「装完即用」，无需人工补规则。
 */
return new class extends Migration
{
    /**
     * 2026_08_06_000007 当年播种的 10 条惩罚规则名。
     * 只有「分类为 discipline」且「名称在其中」的记录才会被删除 —— 精确匹配那批自动播种的行，
     * 不触碰教师在界面上自建的任何规则（界面提供的分类里没有 discipline）。
     */
    private const LEGACY_PUNISHMENT_NAMES = [
        '迟到', '上课讲话', '未交作业', '乱扔垃圾', '追逐打闹',
        '说脏话', '上课睡觉', '打架', '罚站', '罚跑步',
    ];

    public function up(): void
    {
        DB::table('score_rules')
            ->where('category', 'discipline')
            ->whereNull('class_id')
            ->whereIn('name', self::LEGACY_PUNISHMENT_NAMES)
            ->delete();

        foreach (DB::table('schools')->pluck('id') as $schoolId) {
            ScoreRuleService::ensureDefaultsForSchool((int) $schoolId);
        }
    }

    public function down(): void
    {
        // 数据播种不可逆，down 不恢复
    }
};
