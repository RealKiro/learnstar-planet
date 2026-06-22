<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
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

    public function school()
    {
        return $this->belongsTo(School::class);
    }

    public function teacher()
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    public function students()
    {
        return $this->hasMany(Student::class);
    }

    public function pets()
    {
        return $this->hasMany(Pet::class);
    }

    public function notices()
    {
        return $this->hasMany(Notice::class);
    }

    public function scoreRules()
    {
        return $this->hasMany(ScoreRule::class);
    }

    public function shopItems()
    {
        return $this->hasMany(ShopItem::class);
    }
}
