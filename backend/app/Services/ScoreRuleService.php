<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ScoreRule;

// 积分规则共享服务：默认规则模板统一（教师端/教室端同源），教室端据此连通后台规则
class ScoreRuleService
{
    /**
     * 默认积分规则（school 级创建，class_id=null，全校共享）
     * 教师端 listScoreRules 与教室端 display scoreRules 复用此模板，保证两端原因一致。
     */
    public const DEFAULT_RULES = [
        // 📖 课堂表现
        ['name' => '举手发言', 'amount' => 3, 'category' => 'classroom', 'is_positive' => true],
        ['name' => '认真听讲', 'amount' => 2, 'category' => 'classroom', 'is_positive' => true],
        ['name' => '积极互动', 'amount' => 3, 'category' => 'classroom', 'is_positive' => true],
        ['name' => '课堂专注', 'amount' => 2, 'category' => 'classroom', 'is_positive' => true],
        ['name' => '挑战难题', 'amount' => 5, 'category' => 'classroom', 'is_positive' => true],
        ['name' => '上课走神', 'amount' => -2, 'category' => 'classroom', 'is_positive' => false],
        ['name' => '打扰课堂', 'amount' => -3, 'category' => 'classroom', 'is_positive' => false],
        ['name' => '追逐打闹', 'amount' => -5, 'category' => 'classroom', 'is_positive' => false],
        ['name' => '课堂喧哗', 'amount' => -2, 'category' => 'classroom', 'is_positive' => false],
        ['name' => '趴桌睡觉', 'amount' => -2, 'category' => 'classroom', 'is_positive' => false],
        // 📝 作业管理
        ['name' => '作业优秀', 'amount' => 5, 'category' => 'homework', 'is_positive' => true],
        ['name' => '作业按时完成', 'amount' => 3, 'category' => 'homework', 'is_positive' => true],
        ['name' => '作业有进步', 'amount' => 3, 'category' => 'homework', 'is_positive' => true],
        ['name' => '书写工整', 'amount' => 2, 'category' => 'homework', 'is_positive' => true],
        ['name' => '作业缺交', 'amount' => -3, 'category' => 'homework', 'is_positive' => false],
        ['name' => '作业敷衍', 'amount' => -2, 'category' => 'homework', 'is_positive' => false],
        ['name' => '作业迟交', 'amount' => -1, 'category' => 'homework', 'is_positive' => false],
        // 🌟 行为习惯
        ['name' => '遵守纪律', 'amount' => 2, 'category' => 'behavior', 'is_positive' => true],
        ['name' => '帮助同学', 'amount' => 4, 'category' => 'behavior', 'is_positive' => true],
        ['name' => '诚实守信', 'amount' => 5, 'category' => 'behavior', 'is_positive' => true],
        ['name' => '拾金不昧', 'amount' => 5, 'category' => 'behavior', 'is_positive' => true],
        ['name' => '尊敬师长', 'amount' => 3, 'category' => 'behavior', 'is_positive' => true],
        ['name' => '说脏话', 'amount' => -3, 'category' => 'behavior', 'is_positive' => false],
        ['name' => '与同学冲突', 'amount' => -5, 'category' => 'behavior', 'is_positive' => false],
        ['name' => '撒谎欺骗', 'amount' => -5, 'category' => 'behavior', 'is_positive' => false],
        // 📊 综合素养
        ['name' => '科技创新', 'amount' => 8, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '阅读之星', 'amount' => 5, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '体育锻炼', 'amount' => 3, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '艺术表现', 'amount' => 5, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '劳动积极', 'amount' => 3, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '节约环保', 'amount' => 2, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '竞赛获奖', 'amount' => 10, 'category' => 'literacy', 'is_positive' => true],
        ['name' => '破坏公物', 'amount' => -8, 'category' => 'literacy', 'is_positive' => false],
        // 📅 日常表现
        ['name' => '全勤表现', 'amount' => 3, 'category' => 'daily', 'is_positive' => true],
        ['name' => '按时到校', 'amount' => 2, 'category' => 'daily', 'is_positive' => true],
        ['name' => '迟到早退', 'amount' => -2, 'category' => 'daily', 'is_positive' => false],
        ['name' => '仪容整洁', 'amount' => 1, 'category' => 'daily', 'is_positive' => true],
        ['name' => '值日认真', 'amount' => 2, 'category' => 'daily', 'is_positive' => true],
        ['name' => '不戴红领巾/校牌', 'amount' => -1, 'category' => 'daily', 'is_positive' => false],
    ];

    /**
     * 按班级 + 学校取积分规则；无规则时自动创建学校级默认规则
     *
     * @return \Illuminate\Database\Eloquent\Collection<int, ScoreRule>
     */
    public static function rulesForClass(int $classId, int $schoolId): \Illuminate\Database\Eloquent\Collection
    {
        $rules = ScoreRule::where(function ($q) use ($classId, $schoolId) {
            $q->where('class_id', $classId)
              ->orWhere(function ($q2) use ($schoolId) {
                  $q2->whereNull('class_id')->where('school_id', $schoolId);
              });
        })->orderBy('sort_order')->get();

        if ($rules->isEmpty() && $schoolId) {
            foreach (self::DEFAULT_RULES as $i => $d) {
                ScoreRule::create([
                    'class_id' => null,
                    'school_id' => $schoolId,
                    'name' => $d['name'], 'amount' => $d['amount'],
                    'category' => $d['category'], 'is_positive' => $d['is_positive'],
                    'is_active' => true, 'sort_order' => $i,
                ]);
            }
            $rules = ScoreRule::where('school_id', $schoolId)
                ->whereNull('class_id')->orderBy('sort_order')->get();
        }

        return $rules;
    }
}
