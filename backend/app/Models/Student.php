<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

class Student extends Model
{
    use SoftDeletes;

    protected $fillable = [
        'class_id',
        'parent_id',
        'name',
        'student_no',       // 学号
        'avatar_path',
        'total_score',      // 总积分
        'status',
    ];

    protected $casts = [
        'total_score' => 'integer',
    ];

    // ========== 关系 ==========

    public function classRoom()
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function parent()
    {
        return $this->belongsTo(User::class, 'parent_id');
    }

    public function pet()
    {
        return $this->hasOne(Pet::class);
    }

    public function scores()
    {
        return $this->hasMany(Score::class);
    }

    public function scoreLogs()
    {
        return $this->hasMany(ScoreLog::class);
    }

    public function shopRedemptions()
    {
        return $this->hasMany(ShopRedemption::class);
    }
}
