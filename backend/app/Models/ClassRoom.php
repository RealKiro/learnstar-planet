<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\SoftDeletes;

class ClassRoom extends Model
{
    use SoftDeletes;

    protected $fillable = [
        'school_id',
        'name',            // 如：三年二班
        'grade',           // 年级
        'year',            // 学年
        'teacher_id',      // 班主任
        'max_students',    // 0 = 不限制（全免费）
        'settings',        // JSON: 班级级配置
        'status',
    ];

    protected $casts = [
        'settings' => 'array',
    ];

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
