<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\Notice;
use App\Models\ShopItem;
use App\Models\ShopRedemption;
use App\Models\Student;
use App\Models\User;
use App\Models\Wallet;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Support\Collection;

/**
 * 班级商城：商品管理与兑换状态机（pending → approved/rejected → delivered）。
 *
 * 只负责数据操作与结算；HTTP 校验、422/400 映射与响应包装留在控制器。
 * 所有查询都经 TeacherClassScope 限定班级作用域（学校级商品按 school_id 放行）。
 */
class ShopService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
        private readonly ScoreService $scoreService,
        private readonly CurrencyService $currencyService,
    ) {
    }

    /**
     * 商品列表：学校级（class_id=null）+ 教师可管理班级级。
     *
     * 首次访问（列表为空且学校已注册）自动播种学校级默认商品。
     *
     * @return Collection<int, ShopItem>
     */
    public function itemsFor(User $teacher, ?string $currencyType): Collection
    {
        $classIds = $this->scope->ids($teacher);

        $query = ShopItem::where(function ($q) use ($teacher, $classIds) {
            $q->where('school_id', $teacher->school_id)->whereNull('class_id')
              ->orWhereIn('class_id', $classIds);
        });

        if ($currencyType) {
            $query->byCurrency($currencyType);
        }

        $items = $query->orderBy('category')->orderBy('cost_score')->get();

        if ($items->isEmpty() && $teacher->school_id) {
            $this->seedDefaultItems($teacher);

            return ShopItem::where('school_id', $teacher->school_id)->whereNull('class_id')
                ->orderBy('category')->orderBy('cost_score')->get();
        }

        return $items;
    }

    /** 按作用域查找商品：学校级或本班班级级（越权即 404） */
    public function findItem(User $teacher, int $id): ShopItem
    {
        $classIds = $this->scope->ids($teacher);

        return ShopItem::where(function ($q) use ($teacher, $classIds) {
            $q->where('school_id', $teacher->school_id)->whereNull('class_id')
                ->orWhereIn('class_id', $classIds);
        })->findOrFail($id);
    }

    /**
     * @param  array<string, mixed>  $attributes
     */
    public function createItem(User $teacher, array $attributes): ShopItem
    {
        return ShopItem::create([
            'class_id' => null,
            'school_id' => $teacher->school_id,
            'name' => $attributes['name'],
            'description' => $attributes['description'],
            'category' => $attributes['category'] ?? 'physical',
            'cost_score' => (int) $attributes['cost_score'],
            'currency_type' => $attributes['currency_type'] ?? 'score',
            'event_tag' => $attributes['event_tag'],
            'stock' => (int) ($attributes['stock'] ?? 0),
            'image_path' => $attributes['image_path'],
            'is_active' => true,
        ]);
    }

    /**
     * @param  array<string, mixed>  $attributes
     */
    public function updateItem(ShopItem $item, array $attributes): ShopItem
    {
        $item->update($attributes);

        return $item;
    }

    public function deleteItem(ShopItem $item): void
    {
        $item->delete();
    }

    /** 兑换记录列表（按教师可管理班级过滤，每页 20 条） */
    public function paginateRedemptions(User $teacher): LengthAwarePaginator
    {
        return ShopRedemption::whereIn('class_id', $this->scope->ids($teacher))
            ->with(['student:id,name', 'shopItem:id,name,cost_score,currency_type,event_tag,category'])
            ->orderBy('created_at', 'desc')
            ->paginate(20);
    }

    /**
     * 创建兑换记录（教师代学生发起）。
     *
     * class_id 存学生所在班级：学校级商品（class_id=null）的兑换单也能被该班教师看到/审批。
     */
    public function createRedemption(User $teacher, int $studentId, ShopItem $item): ShopRedemption
    {
        $student = Student::whereIn('class_id', $this->scope->ids($teacher))->findOrFail($studentId);

        return ShopRedemption::create([
            'student_id' => $student->id,
            'shop_item_id' => $item->id,
            'class_id' => $student->class_id,
            'cost' => $item->cost_score,
            'status' => 'pending',
        ]);
    }

    /**
     * 审批通过并结算。
     *
     * 结算规则：积分充值类商品走 CurrencyService::exchange（扣积分+扣经验→按汇率发钱包币）；
     * 积分商品走 ScoreService::spendScore（扣积分+扣宠物经验）；钱包币商品只扣余额。
     * 特权类商品自动发一条已发布的班级事件通知。
     *
     * @return array{message: string, remaining_score: int, pet_level: int|null}
     *
     * @throws \DomainException 兑换单不在待审状态
     */
    public function approveRedemption(User $teacher, int $id): array
    {
        $redemption = ShopRedemption::with(['student.pet', 'shopItem'])
            ->whereIn('class_id', $this->scope->ids($teacher))->findOrFail($id);

        if ($redemption->status !== 'pending') {
            throw new \DomainException('该兑换已处理');
        }

        /** @var Student $student */
        $student = $redemption->student;
        $item = $redemption->shopItem;
        $itemName = $item->name ?? '未知物品';
        $cost = $redemption->cost;
        $currency = $item->currency_type ?? 'score';

        // 根据商品类型选择结算方式
        if ($item && $item->category === 'points' && in_array($currency, ['science', 'reading', 'class_point'], true)) {
            // 积分充值类（如"科学币+5"）：扣积分 + 扣宠物经验 → 按汇率发放钱包币
            $this->currencyService->exchange($student->id, $currency, $cost, $teacher->id);
        } elseif ($currency === 'score') {
            // 积分兑换：扣积分 + 扣宠物经验
            $this->scoreService->spendScore($student, $cost, '兑换：' . $itemName, $teacher->id);
        } else {
            // 钱包币兑换：只扣钱包余额
            $this->currencyService->spend($student->id, $currency, $cost, '兑换：' . $itemName);
        }

        $redemption->update([
            'status' => 'approved',
            'approved_by' => $teacher->id,
            'approved_at' => now(),
        ]);

        // 特权奖励自动发班级通知
        if ($item && $item->category === 'privilege') {
            Notice::create([
                'class_id' => $student->class_id,
                'school_id' => $teacher->school_id,
                'title' => '特权奖励：' . $itemName,
                'content' => $student->name . ' 使用 ' . $cost . ' ' . ($currency === 'score' ? '积分' : (Wallet::currencies()[$currency] ?? $currency)) . ' 兑换了「' . $itemName . '」',
                'type' => 'event',
                'published_by' => $teacher->id,
                'is_published' => true,
                'published_at' => now(),
            ]);
        }

        return [
            'message' => '已批准兑换，扣除 ' . $cost . ' ' . ($currency === 'score' ? '积分' : (Wallet::currencies()[$currency] ?? $currency)),
            'remaining_score' => $student->fresh()->total_score,
            'pet_level' => $student->pet?->fresh()->level,
        ];
    }

    /** 拒绝兑换（不结算） */
    public function rejectRedemption(User $teacher, int $id): void
    {
        $redemption = ShopRedemption::whereIn('class_id', $this->scope->ids($teacher))->findOrFail($id);
        $redemption->update(['status' => 'rejected']);
    }

    /** 标记为已发放 */
    public function deliverRedemption(User $teacher, int $id): void
    {
        $redemption = ShopRedemption::whereIn('class_id', $this->scope->ids($teacher))->findOrFail($id);
        $redemption->update(['status' => 'delivered']);
    }

    /** 首次访问播种学校级默认商品（价格按日最高加分可攒周期分层） */
    private function seedDefaultItems(User $teacher): void
    {
        $defaults = [
            // 积分充值类（按 2:1 汇率）
            ['name' => '班级积分 +10', 'description' => '兑换 10 班级积分', 'category' => 'points', 'cost_score' => 20, 'currency_type' => 'class_point'],
            ['name' => '科学币 +5', 'description' => '兑换 5 科学币', 'category' => 'points', 'cost_score' => 10, 'currency_type' => 'science'],
            ['name' => '读书币 +5', 'description' => '兑换 5 读书币', 'category' => 'points', 'cost_score' => 10, 'currency_type' => 'reading'],
            ['name' => '体育币 +5', 'description' => '兑换 5 体育币', 'category' => 'points', 'cost_score' => 10, 'currency_type' => 'class_point'],
            // 小商品 ≈100（日最高 20 分 × 5 天 = 一周可攒）
            ['name' => '铅笔', 'description' => '标准 HB 铅笔一支', 'category' => 'stationery', 'cost_score' => 100, 'currency_type' => 'score'],
            ['name' => '橡皮擦', 'description' => '4B 橡皮擦一块', 'category' => 'stationery', 'cost_score' => 100, 'currency_type' => 'score'],
            ['name' => '草稿纸', 'description' => 'A4 草稿纸 10 张', 'category' => 'stationery', 'cost_score' => 100, 'currency_type' => 'score'],
            ['name' => '免罚站一次', 'description' => '免除一次罚站', 'category' => 'privilege', 'cost_score' => 100, 'currency_type' => 'score'],
            ['name' => '免罚跑步一次', 'description' => '免除一次罚跑步', 'category' => 'privilege', 'cost_score' => 100, 'currency_type' => 'score'],
            // 中商品 120~150
            ['name' => '便利贴', 'description' => '彩色便利贴一本', 'category' => 'stationery', 'cost_score' => 120, 'currency_type' => 'score'],
            ['name' => '黑色圆珠笔', 'description' => '0.5mm 黑色圆珠笔一支', 'category' => 'stationery', 'cost_score' => 150, 'currency_type' => 'score'],
            ['name' => '蓝色圆珠笔', 'description' => '0.5mm 蓝色圆珠笔一支', 'category' => 'stationery', 'cost_score' => 150, 'currency_type' => 'score'],
            ['name' => '红色圆珠笔', 'description' => '红色批改用笔一支', 'category' => 'stationery', 'cost_score' => 150, 'currency_type' => 'score'],
            ['name' => '香蕉', 'description' => '新鲜香蕉一根 🍌（请勿乱扔果皮）', 'category' => 'food', 'cost_score' => 150, 'currency_type' => 'score'],
            ['name' => '免做卫生一次', 'description' => '免除一次值日卫生', 'category' => 'privilege', 'cost_score' => 150, 'currency_type' => 'score'],
            // 大商品 ≈200（两周可攒）
            ['name' => '练习本', 'description' => '方格练习本一本', 'category' => 'stationery', 'cost_score' => 180, 'currency_type' => 'score'],
            ['name' => '苹果', 'description' => '新鲜苹果一个 🍎（请勿乱扔果皮）', 'category' => 'food', 'cost_score' => 180, 'currency_type' => 'score'],
            ['name' => '饮料', 'description' => '矿泉水/饮料一瓶 🧃', 'category' => 'food', 'cost_score' => 200, 'currency_type' => 'score'],
            ['name' => '牛奶', 'description' => '纯牛奶一盒 🥛', 'category' => 'food', 'cost_score' => 200, 'currency_type' => 'score'],
            ['name' => '集体观影', 'description' => '全班集体观影一次', 'category' => 'activity', 'cost_score' => 200, 'currency_type' => 'score'],
            ['name' => '免作业一次', 'description' => '免交一次作业', 'category' => 'privilege', 'cost_score' => 200, 'currency_type' => 'score'],
            ['name' => '3D打印作品', 'description' => '3D 打印小作品一件', 'category' => 'physical', 'cost_score' => 200, 'currency_type' => 'score'],
        ];

        foreach ($defaults as $d) {
            ShopItem::create([
                'class_id' => null,
                'school_id' => $teacher->school_id,
                'name' => $d['name'],
                'description' => $d['description'],
                'category' => $d['category'],
                'cost_score' => $d['cost_score'],
                'currency_type' => $d['currency_type'],
                'stock' => 0,
                'is_active' => true,
            ]);
        }
    }
}
