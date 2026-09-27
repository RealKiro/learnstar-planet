<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\School;
use App\Models\ScoreRule;
use App\Models\User;

// 积分规则共享服务：默认规则模板统一（教师端/教室端同源），教室端据此连通后台规则
class ScoreRuleService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

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
        ['name' => '乱扔垃圾', 'amount' => -2, 'category' => 'literacy', 'is_positive' => false],
        // 📅 日常表现
        ['name' => '全勤表现', 'amount' => 3, 'category' => 'daily', 'is_positive' => true],
        ['name' => '按时到校', 'amount' => 2, 'category' => 'daily', 'is_positive' => true],
        ['name' => '迟到早退', 'amount' => -2, 'category' => 'daily', 'is_positive' => false],
        ['name' => '仪容整洁', 'amount' => 1, 'category' => 'daily', 'is_positive' => true],
        ['name' => '值日认真', 'amount' => 2, 'category' => 'daily', 'is_positive' => true],
        ['name' => '不戴红领巾/校牌', 'amount' => -1, 'category' => 'daily', 'is_positive' => false],
        // 📚 学业表现（前端 categoryLabels 一直有 academic 标签，此前无规则落在此分类）
        ['name' => '考试优秀', 'amount' => 10, 'category' => 'academic', 'is_positive' => true],
        ['name' => '考试进步', 'amount' => 8, 'category' => 'academic', 'is_positive' => true],
        ['name' => '考试作弊', 'amount' => -15, 'category' => 'academic', 'is_positive' => false],
    ];

    /**
     * 分类标签（后端唯一真源；键集必须与前端 utils/scoreRules.ts 的 categoryLabels 一致）。
     */
    public const CATEGORY_LABELS = [
        'classroom' => '课堂表现',
        'homework'  => '作业管理',
        'behavior'  => '行为习惯',
        'literacy'  => '综合素养',
        'daily'     => '日常表现',
        'academic'  => '学业表现',
        'custom'    => '自定义',
    ];

    /**
     * 增量补齐某校的学校级默认规则（幂等：缺哪条补哪条，已存在的同名规则一律不覆盖）。
     *
     * 「安装即可用」的关键：
     *   ① 判据是「该校是否播种过」而不是「该校规则是否为空」——旧实现用后者，
     *      导致只播种过惩罚规则的学校永远拿不到奖励规则（加分区分区空白）；
     *   ② 播种后写 schools.settings.score_rules_seeded 标记，之后教师对默认规则的
     *      增/删/改一律尊重，不会在下次访问时被"复活"。
     *
     * @return int 本次新建的规则条数
     */
    public static function ensureDefaultsForSchool(?int $schoolId, bool $force = false): int
    {
        if (!$schoolId) {
            return 0;
        }

        $school = School::find($schoolId);
        if (!$school) {
            return 0;
        }

        $settings = is_array($school->settings) ? $school->settings : [];
        if (!$force && !empty($settings['score_rules_seeded'])) {
            return 0;
        }

        $created = 0;
        foreach (self::DEFAULT_RULES as $i => $d) {
            $rule = ScoreRule::firstOrCreate(
                ['class_id' => null, 'school_id' => $schoolId, 'name' => $d['name']],
                [
                    'amount' => $d['amount'],
                    'category' => $d['category'],
                    'is_positive' => $d['is_positive'],
                    'is_active' => true,
                    'sort_order' => $i,
                ]
            );
            if ($rule->wasRecentlyCreated) {
                $created++;
            }
        }

        $settings['score_rules_seeded'] = true;
        $school->settings = $settings;
        $school->save();

        return $created;
    }

    /**
     * 按班级 + 学校取积分规则；首次访问时增量补齐学校级默认规则（幂等，见 ensureDefaultsForSchool）。
     *
     * @return \Illuminate\Database\Eloquent\Collection<int, ScoreRule>
     */
    public static function rulesForClass(int $classId, int $schoolId): \Illuminate\Database\Eloquent\Collection
    {
        if ($schoolId) {
            self::ensureDefaultsForSchool($schoolId);
        }

        return ScoreRule::where(function ($q) use ($classId, $schoolId) {
            $q->where('class_id', $classId)
              ->orWhere(function ($q2) use ($schoolId) {
                  $q2->whereNull('class_id')->where('school_id', $schoolId);
              });
        })->orderBy('sort_order')->get();
    }

    /**
     * 教师端规则列表：本班班级规则 + 本校学校级规则（class_id=null 且 school_id=本校），避免跨校泄漏。
     * 首次访问时增量补齐学校级默认规则（同校所有教师共享，见 ensureDefaultsForSchool）。
     *
     * @return \Illuminate\Database\Eloquent\Collection<int, ScoreRule>
     */
    public function listForTeacher(User $teacher): \Illuminate\Database\Eloquent\Collection
    {
        $classIds = $this->scope->ids($teacher);

        if ($teacher->school_id) {
            self::ensureDefaultsForSchool($teacher->school_id);
        }

        return ScoreRule::where(function ($q) use ($classIds, $teacher) {
            $q->whereIn('class_id', $classIds)
              ->orWhere(function ($q2) use ($teacher) {
                  $q2->whereNull('class_id')->where('school_id', $teacher->school_id);
              });
        })->orderBy('sort_order')->get();
    }

    /**
     * 教师端创建规则（默认落到教师首个可管理班级）。
     */
    public function createForTeacher(User $teacher, array $attributes): ScoreRule
    {
        $classIds = $this->scope->ids($teacher);

        return ScoreRule::create([
            'class_id' => $attributes['class_id'] ?? $classIds->first(),
            'name' => $attributes['name'],
            'amount' => (int) $attributes['amount'],
            'category' => $attributes['category'] ?? 'custom',
            'is_positive' => (bool) ($attributes['is_positive'] ?? true),
            'is_active' => true,
            'sort_order' => 0,
        ]);
    }

    /**
     * 教师可见范围内的单条规则（本班班级规则 + 本校学校级规则）。
     * 跨校/越权 ID 一律 404，与 listForTeacher 同一口径。
     */
    public function findScopedForTeacher(User $teacher, int $id): ScoreRule
    {
        $classIds = $this->scope->ids($teacher);

        return ScoreRule::where(function ($q) use ($classIds, $teacher) {
            $q->whereIn('class_id', $classIds)
              ->orWhere(function ($q2) use ($teacher) {
                  $q2->whereNull('class_id')->where('school_id', $teacher->school_id);
              });
        })->findOrFail($id);
    }

    /**
     * 教师端更新规则（可见范围同 listForTeacher，跨校不可见不可改）。
     */
    public function updateForTeacher(User $teacher, int $id, array $attributes): ScoreRule
    {
        $rule = $this->findScopedForTeacher($teacher, $id);

        $rule->update($attributes);

        return $rule;
    }

    /**
     * 教师端删除规则（可见范围同 listForTeacher，跨校不可见不可删）。
     */
    public function deleteForTeacher(User $teacher, int $id): void
    {
        $rule = $this->findScopedForTeacher($teacher, $id);
        $rule->delete();
    }
}
