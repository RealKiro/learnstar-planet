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

    public function school(): BelongsTo
    {
        return $this->belongsTo(School::class);
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    public function students(): HasMany
   