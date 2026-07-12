<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * @property int $id
 * @property int $student_id
 * @property string $currency_type
 * @property int $balance
 */
class Wallet extends Model
{
    protected $fillable = [
        'student_id',
        'currency_type',
        'balance',
    ];

    protected $casts = [
        'balance' => 'integer',
    ];

    /**
     * 支持的币种列表
     *
     * @return array<string, string>
     */
    public static function currencies(): array
    {
        return [
            'science'     => '科学币',
            'reading'     => '读书币',
            'class_point' => '班级积分',
        ];
    }

    /**
     * 获取或创建学生钱包
     */
    public static function getOrCreate(int $studentId, string $currencyType): self
    {
        return self::firstOrCreate(
            ['student_id' => $studentId, 'currency_type' => $currencyType],
            ['balance' => 0],
        );
    }

    public function student(): BelongsTo
    {
        return $this->belongsTo(Student::class);
    }
}
