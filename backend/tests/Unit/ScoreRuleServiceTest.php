<?php

declare(strict_types=1);

namespace Tests\Unit;

use App\Http\Controllers\Api\SchoolAdminController;
use App\Http\Controllers\Api\StudentController;
use App\Models\School;
use App\Models\ScoreRule;
use App\Models\User;
use App\Services\ScoreRuleService;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Http\Request;
use PHPUnit\Framework\Attributes\Test;
use Tests\TestCase;

/**
 * 「安装即可用」守护测试。
 *
 * 目标状态：新装/升级后默认就带一套完整、成对的奖励 + 惩罚规则，
 * 教师只需按需补特殊规则，不必先手工把基础规则加一遍，也不会出现
 * 「只有惩罚、没有奖励」或分类渲染成裸 slug 的情况。
 */
class ScoreRuleServiceTest extends TestCase
{
    use RefreshDatabase;

    private function makeSchool(string $code = 'TEST01'): School
    {
        return School::create(['name' => '测试学校', 'code' => $code, 'status' => 'active']);
    }

    private function makeTeacher(School $school): User
    {
        return User::create([
            'school_id' => $school->id,
            'role' => 'teacher',
            'username' => 't' . $school->id,
            'password' => bcrypt('secret'),
            'name' => '测试教师',
            'status' => 'active',
        ]);
    }

    #[Test]
    public function 模板规则的分类必须全部落在唯一真源键集内(): void
    {
        $canonical = array_keys(ScoreRuleService::CATEGORY_LABELS);

        foreach (ScoreRuleService::DEFAULT_RULES as $rule) {
            $this->assertContains(
                $rule['category'],
                $canonical,
                "规则「{$rule['name']}」使用了不在 CATEGORY_LABELS 内的分类：{$rule['category']}"
            );
        }
    }

    #[Test]
    public function 模板同时包含奖励与惩罚且数量不退化(): void
    {
        $positive = array_filter(ScoreRuleService::DEFAULT_RULES, fn ($r) => $r['is_positive'] === true);
        $negative = array_filter(ScoreRuleService::DEFAULT_RULES, fn ($r) => $r['is_positive'] === false);

        $this->assertGreaterThanOrEqual(5, count($positive), '奖励规则过少，装完即用不成立');
        $this->assertGreaterThanOrEqual(5, count($negative), '惩罚规则过少，装完即用不成立');
    }

    #[Test]
    public function 新学校会被补齐完整默认规则集(): void
    {
        $school = $this->makeSchool();

        $created = ScoreRuleService::ensureDefaultsForSchool($school->id);

        $this->assertSame(count(ScoreRuleService::DEFAULT_RULES), $created);

        $rules = ScoreRule::where('school_id', $school->id)->whereNull('class_id')->get();
        $this->assertSame(count(ScoreRuleService::DEFAULT_RULES), $rules->count());
        $this->assertGreaterThan(0, $rules->where('is_positive', true)->count(), '缺少奖励规则');
        $this->assertGreaterThan(0, $rules->where('is_positive', false)->count(), '缺少惩罚规则');

        // 兜底：不得出现前端认不出的分类（历史 discipline 那批即为此类）
        $canonical = array_keys(ScoreRuleService::CATEGORY_LABELS);
        $this->assertSame([], $rules->pluck('category')->unique()->diff($canonical)->values()->all());
    }

    #[Test]
    public function 只有惩罚规则的学校依然会被补上奖励规则(): void
    {
        // 复现历史缺陷：早期迁移只播种了 10 条 discipline 惩罚规则，
        // 旧逻辑（"规则为空才播种"）因此永不触发 → 该校加分区分空白。
        $school = $this->makeSchool();
        ScoreRule::create([
            'class_id' => null,
            'school_id' => $school->id,
            'name' => '迟到',
            'amount' => -2,
            'category' => 'discipline',
            'is_positive' => false,
            'is_active' => true,
            'sort_order' => 0,
        ]);

        $teacher = $this->makeTeacher($school);
        $rules = app(ScoreRuleService::class)->listForTeacher($teacher);

        $this->assertGreaterThan(
            0,
            $rules->where('is_positive', true)->count(),
            '只有惩罚规则的学校必须被补上奖励规则'
        );
    }

    #[Test]
    public function 重复补齐是幂等的(): void
    {
        $school = $this->makeSchool();

        ScoreRuleService::ensureDefaultsForSchool($school->id);
        $afterFirst = ScoreRule::where('school_id', $school->id)->count();

        $created = ScoreRuleService::ensureDefaultsForSchool($school->id);
        $afterSecond = ScoreRule::where('school_id', $school->id)->count();

        $this->assertSame(0, $created, '第二次补齐不应再创建任何规则');
        $this->assertSame($afterFirst, $afterSecond, '重复补齐不得产生重复规则');
    }

    #[Test]
    public function 教师删除的默认规则不会被复活(): void
    {
        $school = $this->makeSchool();
        ScoreRuleService::ensureDefaultsForSchool($school->id);

        $victim = ScoreRule::where('school_id', $school->id)->where('name', '举手发言')->firstOrFail();
        $victim->delete();

        // 再次访问（等同教师端重新进入规则页）
        ScoreRuleService::ensureDefaultsForSchool($school->id);

        $this->assertNull(
            ScoreRule::where('school_id', $school->id)->where('name', '举手发言')->first(),
            '教师删除的规则不应被再次补齐（播种标记已写入 schools.settings）'
        );
    }

    #[Test]
    public function 分类端点与唯一真源保持一致(): void
    {
        $payload = app(StudentController::class)->scoreCategories()->getData(true);

        $this->assertSame(
            array_keys(ScoreRuleService::CATEGORY_LABELS),
            array_column($payload['data'], 'id')
        );
    }

    #[Test]
    public function 升级迁移会清理历史discipline分类并补齐完整规则集(): void
    {
        $school = $this->makeSchool();

        // 复现历史现场：早期迁移为「已存在的学校」播种了 10 条 discipline 惩罚规则
        $legacy = ['迟到', '上课讲话', '未交作业', '乱扔垃圾', '追逐打闹', '说脏话', '上课睡觉', '打架', '罚站', '罚跑步'];
        foreach ($legacy as $i => $name) {
            ScoreRule::create([
                'class_id' => null,
                'school_id' => $school->id,
                'name' => $name,
                'amount' => -1,
                'category' => 'discipline',
                'is_positive' => false,
                'is_active' => true,
                'sort_order' => 0,
            ]);
        }

        // 执行「升级迁移」本体（匿名类实例，与 artisan migrate 加载方式一致）
        $migration = require database_path('migrations/2026_09_20_000001_ensure_default_score_rules.php');
        $migration->up();

        $rules = ScoreRule::where('school_id', $school->id)->get();

        $this->assertSame(
            0,
            $rules->where('category', 'discipline')->count(),
            '历史 discipline 分类必须被清理（前端标签表里没有该分类）'
        );
        $this->assertGreaterThan(0, $rules->where('is_positive', true)->count(), '升级后必须补齐奖励规则');
        $this->assertSame(
            count(ScoreRuleService::DEFAULT_RULES),
            $rules->count(),
            '升级后规则总数应等于唯一真源模板条数'
        );
    }

    #[Test]
    public function 管理员端规则列表同样会触发默认规则补齐(): void
    {
        $school = $this->makeSchool();
        $admin = User::create([
            'school_id' => $school->id,
            'role' => 'school_admin',
            'username' => 'admin' . $school->id,
            'password' => bcrypt('secret'),
            'name' => '测试管理员',
            'status' => 'active',
        ]);

        $request = Request::create('/api/v1/admin/score-rules', 'GET');
        $request->setUserResolver(fn () => $admin);

        app(SchoolAdminController::class)->adminListScoreRules($request);

        $this->assertSame(
            count(ScoreRuleService::DEFAULT_RULES),
            ScoreRule::where('school_id', $school->id)->count(),
            '管理员先打开规则页时也应拿到完整默认规则集'
        );
    }
}
