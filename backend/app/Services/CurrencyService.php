<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ExchangeLog;
use App\Models\ExchangeRate;
use App\Models\Score;
use App\Models\ScoreLog;
use App\Models\Student;
use App\Models\User;
use App\Models\Wallet;
use Illuminate\Support\Facades\DB;

class CurrencyService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

    /**
     * 积分 → 币种兑换（扣减学生 total_score，增加钱包余额，扣减宠物经验）
     *
     * @return array{remaining_score: int, wallet_balance: int}
     *
     * @throws \DomainException
     */
    public function exchange(int $studentId, string $toCurrency, int $scoreAmount, int $operatedBy): array
    {
        if ($scoreAmount <= 0) {
            throw new \DomainException('兑换积分必须大于 0');
        }

        $student = Student::with(['pet', 'classRoom'])->findOrFail($studentId);

        if ($student->total_score < $scoreAmount) {
            throw new \DomainException('积分不足，当前余额：' . $student->total_score);
        }

        // 获取学校汇率配置
        $schoolId = $student->classRoom?->school_id;
        if (!$schoolId) {
            throw new \DomainException('未找到学生所属学校');
        }

        $rate = ExchangeRate::where('school_id', $schoolId)
            ->where('from_currency', 'score')
            ->where('to_currency', $toCurrency)
            ->where('is_active', true)
            ->first();

        if (!$rate) {
            throw new \DomainException('未找到可用的汇率配置');
        }

        $toAmount = (int) round($scoreAmount * (float) $rate->rate);
        if ($toAmount <= 0) {
            throw new \DomainException('兑换金额过小，无法兑换');
        }

        return DB::transaction(function () use ($student, $scoreAmount, $toCurrency, $toAmount, $operatedBy) {
            $balanceBefore = $student->total_score;

            // 1. 扣减学生积分
            $student->update(['total_score' => $balanceBefore - $scoreAmount]);
            $balanceAfter = $student->total_score;

            // 2. 创建扣分记录
            $currencyLabel = Wallet::currencies()[$toCurrency] ?? $toCurrency;
            $reason = '兑换' . $currencyLabel;
            $score = Score::create([
                'student_id' => $student->id,
                'class_id' => $student->class_id,
                'amount' => -$scoreAmount,
                'reason' => $reason,
                'given_by' => $operatedBy,
            ]);

            // 3. 记录积分日志
            ScoreLog::create([
                'student_id' => $student->id,
                'score_id' => $score->id,
                'balance_before' => $balanceBefore,
                'balance_after' => $balanceAfter,
                'description' => $reason,
            ]);

            // 4. 扣减宠物经验（1:1 比例）
            if ($student->pet) {
                $student->pet->removeExperience($scoreAmount);
            }

            // 5. 增加钱包余额
            $wallet = Wallet::getOrCreate($student->id, $toCurrency);
            $wallet->increment('balance', $toAmount);

            // 6. 记录兑换日志
            ExchangeLog::create([
                'student_id' => $student->id,
                'from_currency' => 'score',
                'to_currency' => $toCurrency,
                'from_amount' => $scoreAmount,
                'to_amount' => $toAmount,
                'operated_by' => $operatedBy,
            ]);

            return [
                'remaining_score' => $balanceAfter,
                'wallet_balance' => $wallet->fresh()->balance,
            ];
        });
    }

    /**
     * 跨币种兑换（科学币 → 班级积分等，操作钱包余额）
     *
     * @return array{from_balance: int, to_balance: int}
     *
     * @throws \DomainException
     */
    public function crossExchange(int $studentId, string $from, string $to, int $amount, int $operatedBy): array
    {
        if ($amount <= 0) {
            throw new \DomainException('兑换数量必须大于 0');
        }

        if ($from === $to) {
            throw new \DomainException('源币种和目标币种不能相同');
        }

        $student = Student::with('classRoom')->findOrFail($studentId);
        $schoolId = $student->classRoom?->school_id;
        if (!$schoolId) {
            throw new \DomainException('未找到学生所属学校');
        }

        $rate = ExchangeRate::where('school_id', $schoolId)
            ->where('from_currency', $from)
            ->where('to_currency', $to)
            ->where('is_active', true)
            ->first();

        if (!$rate) {
            throw new \DomainException('未找到可用的汇率配置');
        }

        $fromWallet = Wallet::getOrCreate($studentId, $from);
        if ($fromWallet->balance < $amount) {
            $fromLabel = Wallet::currencies()[$from] ?? $from;

            throw new \DomainException($fromLabel . '余额不足，当前余额：' . $fromWallet->balance);
        }

        $toAmount = (int) round($amount * (float) $rate->rate);
        if ($toAmount <= 0) {
            throw new \DomainException('兑换金额过小，无法兑换');
        }

        return DB::transaction(function () use ($student, $fromWallet, $from, $to, $amount, $toAmount, $operatedBy) {
            // 扣减源币种
            $fromWallet->decrement('balance', $amount);
            $fromBalance = $fromWallet->fresh()->balance;

            // 增加目标币种
            $toWallet = Wallet::getOrCreate($student->id, $to);
            $toWallet->increment('balance', $toAmount);
            $toBalance = $toWallet->fresh()->balance;

            // 记录兑换日志
            ExchangeLog::create([
                'student_id' => $student->id,
                'from_currency' => $from,
                'to_currency' => $to,
                'from_amount' => $amount,
                'to_amount' => $toAmount,
                'operated_by' => $operatedBy,
            ]);

            return [
                'from_balance' => $fromBalance,
                'to_balance' => $toBalance,
            ];
        });
    }

    /**
     * 消费钱包币种（用于活动专属商城）
     *
     * @throws \DomainException
     */
    public function spend(int $studentId, string $currency, int $amount, string $reason): void
    {
        if ($amount <= 0) {
            throw new \DomainException('消费数量必须大于 0');
        }

        $wallet = Wallet::getOrCreate($studentId, $currency);

        if ($wallet->balance < $amount) {
            $currencyLabel = Wallet::currencies()[$currency] ?? $currency;

            throw new \DomainException($currencyLabel . '余额不足，当前余额：' . $wallet->balance);
        }

        DB::transaction(function () use ($wallet, $amount) {
            $wallet->decrement('balance', $amount);

            ExchangeLog::create([
                'student_id' => $wallet->student_id,
                'from_currency' => $wallet->currency_type,
                'to_currency' => $wallet->currency_type,
                'from_amount' => $amount,
                'to_amount' => $amount,
                'operated_by' => null,
            ]);
        });
    }

    // ============================================================
    // 教师端：汇率配置与钱包/兑换记录查询
    // ============================================================

    /** 默认汇率（2:1 防通胀：2 积分 = 1 币），首次访问惰性播种 */
    private const DEFAULT_RATES = [
        ['name' => '积分 → 科学币', 'from_currency' => 'score', 'to_currency' => 'science', 'rate' => 0.5],
        ['name' => '积分 → 读书币', 'from_currency' => 'score', 'to_currency' => 'reading', 'rate' => 0.5],
        ['name' => '积分 → 体育币', 'from_currency' => 'score', 'to_currency' => 'class_point', 'rate' => 0.5],
    ];

    /**
     * 学校汇率列表；首次访问（为空）惰性初始化默认汇率，保证积分充值类商品不配汇率也能结算。
     *
     * @return \Illuminate\Database\Eloquent\Collection<int, ExchangeRate>
     */
    public function ratesForSchool(int $schoolId): \Illuminate\Database\Eloquent\Collection
    {
        $exists = ExchangeRate::where('school_id', $schoolId)->exists();
        if (!$exists) {
            foreach (self::DEFAULT_RATES as $d) {
                ExchangeRate::firstOrCreate(
                    ['school_id' => $schoolId, 'from_currency' => $d['from_currency'], 'to_currency' => $d['to_currency']],
                    ['name' => $d['name'], 'rate' => $d['rate'], 'is_active' => true],
                );
            }
        }

        return ExchangeRate::where('school_id', $schoolId)
            ->orderBy('from_currency')->orderBy('to_currency')->get();
    }

    /**
     * 新增学校级汇率。
     */
    public function createRate(int $schoolId, array $attributes): ExchangeRate
    {
        return ExchangeRate::create([
            'school_id' => $schoolId,
            'name' => $attributes['name'],
            'from_currency' => $attributes['from_currency'],
            'to_currency' => $attributes['to_currency'],
            'rate' => $attributes['rate'],
            'is_active' => true,
        ]);
    }

    /**
     * 更新学校级汇率。
     */
    public function updateRate(int $schoolId, int $id, array $attributes): ExchangeRate
    {
        $rate = ExchangeRate::where('school_id', $schoolId)->findOrFail($id);
        $rate->update($attributes);

        return $rate->fresh();
    }

    /**
     * 教师所带班级学生的全部钱包。
     */
    public function walletsFor(User $teacher): \Illuminate\Support\Collection
    {
        $classIds = $this->scope->ids($teacher);
        $students = Student::whereIn('class_id', $classIds)->where('status', 'active')->pluck('id');

        return Wallet::whereIn('student_id', $students)
            ->with('student:id,name')
            ->get()
            ->map(fn ($w) => [
                'student_id' => $w->student_id,
                'student_name' => $w->student?->name,
                'currency_type' => $w->currency_type,
                'balance' => (int) $w->balance,
            ]);
    }

    /**
     * 兑换记录（仅教师所带班级学生，分页）。
     *
     * @return array{data: \Illuminate\Support\Collection, meta: array{current_page: int, last_page: int, total: int}}
     */
    public function logsFor(User $teacher): array
    {
        $classIds = $this->scope->ids($teacher);

        $logs = ExchangeLog::with('student:id,name,student_no')
            ->whereIn('student_id', Student::whereIn('class_id', $classIds)->select('id'))
            ->orderBy('created_at', 'desc')
            ->paginate(20);

        return [
            'data' => collect($logs->items())->map(static function (ExchangeLog $log): array {
                return array_merge($log->toArray(), [
                    'student_name' => $log->student->name ?? '已删除学生',
                    'student_no' => $log->student->student_no ?? '',
                ]);
            }),
            'meta' => [
                'current_page' => $logs->currentPage(),
                'last_page' => $logs->lastPage(),
                'total' => $logs->total(),
            ],
        ];
    }
}
