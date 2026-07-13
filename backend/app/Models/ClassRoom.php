<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use App\Models\ClassRoomTeacher;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\SoftDeletes;

class ClassRoom extends Model
{
    use SoftDeletes;

    // 中文序数 → 数字映射
    private const CHINESE_ORDINALS = [
        '第一' => '1', '第二' => '2', '第三' => '3', '第四' => '4', '第五' => '5',
        '第六' => '6', '第七' => '7', '第八' => '8', '第九' => '9', '第十' => '10',
    ];

    // 中文数字 → 数字映射
    private const CHINESE_NUMBERS = [
        '一' => '1', '二' => '2', '三' => '3', '四' => '4', '五' => '5',
        '六' => '6', '七' => '7', '八' => '8', '九' => '9', '十' => '10',
    ];

    protected static function booted(): void
    {
        static::saving(function (self $class) {
            $class->name = self::normalizeClassName($class->name);
        });
    }

    /**
     * 将各种班级名称格式统一为：X年级（Y）班
     * 支持输入：一年级第一班、一年级1班、一年级一 班 等
     */
    public static function normalizeClassName(string $name): string
    {
        // 如果已经是标准格式（含全角括号），直接返回
        if (preg_match('/^.+（\d+）班$/', $name)) {
            return $name;
        }

        // 尝试匹配: X年级 + 序数/数字 + (可选空格) + 班
        if (preg_match('/^(.+年级)(.+?)班$/u', $name, $m)) {
            $grade = $m[1];
            $num = trim($m[2]);

            // 转换中文序数: 第一 → 1
            if (isset(self::CHINESE_ORDINALS[$num])) {
                return $grade . '（' . self::CHINESE_ORDINALS[$num] . '）班';
            }

            // 转换中文数字: 一 → 1
            if (isset(self::CHINESE_NUMBERS[$num])) {
                return $grade . '（' . self::CHINESE_NUMBERS[$num] . '）班';
            }

            // 已经是阿拉伯数字: 1 → 1（加括号）
            if (is_numeric($num)) {
                return $grade . '（' . (int) $num . '）班';
            }
        }

        return $name;
    }

    protected $fillable = [
        'school_id',
        'name',            // 如：三年级（2）班
        'grade',           // 年级
        'year',            // 学年
        'teacher_id',      // 班主任
        'max_students',    // 0 = 不限制（全免费）
        'settings',        // JSON: 班级级配置
        'status',
    ];

    protected $casts = [
        'settings' => 'array',
        'max_students' => 'integer',
        'status' => 'string',
    ];

    // ========== 查询作用域 ==========

    /** @param \Illuminate\Database\Eloquent\Builder $query */
    public function scopeActive($query)
    {
        return $query->where('status', 'active');
    }

    /** @param \Illuminate\Database\Eloquent\Builder $query */
    public function scopeBySchool($query, int $schoolId)
    {
        return $query->where('school_id', $schoolId);
    }

    public function classRoomTeachers(): HasMany
    {
        return $this->hasMany(ClassRoomTeacher::class);
    }

    public function school(): BelongsTo
    {
        return $this->belongsTo(School::class);
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    public function students(): HasMany
    {
        return $this->hasMany(Student::class, 'class_id');
    }

    public function pets(): HasMany
    {
        return $this->hasMany(Pet::class);
    }

    public function notices(): HasMany
    {
        return $this->hasMany(Notice::class);
    }

    public function scoreRules(): HasMany
    {
        return $this->hasMany(ScoreRule::class);
    }

    public function shopItems(): HasMany
    {
        return $this->hasMany(ShopItem::class);
    }
}
