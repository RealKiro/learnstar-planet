<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

/**
 * @property int $id
 * @property int $class_id
 * @property int|null $parent_id
 * @property string $name
 * @property string|null $student_no
 * @property string|null $avatar_path
 * @property int $total_score
 * @property string $status
 * @property \Illuminate\Database\Eloquent\Collection<int, Pet> $pet
 * @property \Illuminate\Database\Eloquent\Collection<int, Score> $scores
 * @property \Illuminate\Database\Eloquent\Collection<int, ScoreLog> $scoreLogs
 * @property \Illuminate\Database\Eloquent\Collection<int, ShopRedemption> $shopRedemptions
 */
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
