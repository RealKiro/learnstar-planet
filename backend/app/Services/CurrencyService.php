<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ExchangeLog;
use App\Models\ExchangeRate;
use App\Models\Score;
use App\Models\ScoreLog;
use App\Models\Student;
use App\Models\Wallet;
use Illuminate\Support\Facades\DB;

class CurrencyService
{
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
}
