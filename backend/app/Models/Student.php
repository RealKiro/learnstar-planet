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
        'gender',          // 性别
        'student_no',       // 学号
        'avatar_path',
        'total_score',      // 总积分
        'status',
    ];

    protected $casts = [
        'total_score' => 'integer',
    ];

    // ========== 关系 ==========

    /**
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo<ClassRoom, $this>
     */
    public function classRoom(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo<User, $this>
     */
    public function parent(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(User::class, 'parent_id');
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\HasOne<Pet, $this>
     */
    public function pet(): \Illuminate\Database\Eloquent\Relations\HasOne
    {
        return $this->hasOne(Pet::class);
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\HasMany<Score, $this>
     */
    public function scores(): \Illuminate\Database\Eloquent\Relations\HasMany
    {
        return $this->hasMany(Score::class);
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\HasMany<ScoreLog, $this>
     */
    public function scoreLogs(): \Illuminate\Database\Eloquent\Relations\HasMany
    {
        return $this->hasMany(ScoreLog::class);
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\HasMany<ShopRedemption, $this>
     */
    public function shopRedemptions(): \Illuminate\Database\Eloquent\Relations\HasMany
    {
        return $this->hasMany(ShopRedemption::class);
    }
}
