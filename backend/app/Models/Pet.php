<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Pet extends Model
{
    protected $fillable = [
        'student_id',
        'class_id',
        'name',
        'type',             // 宠物类型标识
        'level',            // 0-10 级进化
        'experience',       // 经验值
        'mood',             // 心情指数 (0-100)
        'accessories',      // JSON: 装饰配件列表
        'last_fed_at',      // 最后喂养时间
    ];

    protected $casts = [
        'level' => 'integer',
        'experience' => 'integer',
        'mood' => 'integer',
        'accessories' => 'array',
        'last_fed_at' => 'datetime',
    ];

    public function student(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(Student::class);
    }

    public function classRoom(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    // ========== 进化系统 ==========

    // 完全原创的10阶进化体系（灵感来自宇宙探索）
    public static function evolutionStages(): array
    {
        return [
            0  => ['emoji' => '🌟', 'name' => '星尘',   'title' => '初生的光芒'],
            1  => ['emoji' => '🌙', 'name' => '月芽',   'title' => '萌芽的灵光'],
            2  => ['emoji' => '🌱', 'name' => '灵苗',   'title' => '生长的力量'],
            3  => ['emoji' => '🌿', 'name' => '青藤',   'title' => '蓬勃向上'],
            4  => ['emoji' => '🌳', 'name' => '慧树',   'title' => '枝繁叶茂'],
            5  => ['emoji' => '🦋', 'name' => '蝶灵',   'title' => '翩然蜕变'],
            6  => ['emoji' => '🦅', 'name' => '鹰慧',   'title' => '翱翔天际'],
            7  => ['emoji' => '🦁', 'name' => '狮睿',   'title' => '王者风范'],
            8  => ['emoji' => '🦄', 'name' => '灵角',   'title' => '奇幻之力'],
            9  => ['emoji' => '✨', 'name' => '星耀',   'title' => '璀璨巅峰'],
            10 => ['emoji' => '🌌', 'name' => '银河',   'title' => '宇宙之灵'],
        ];
    }

    public static function petTypes(): array
    {
        return [
            'stellar_cat'  => '星猫',
            'moon_rabbit'  => '月兔',
            'sun_deer'     => '日鹿',
            'cloud_fox'    => '云狐',
            'sky_penguin'  => '天企鹅',
            'light_horse'  => '光马',
            'dream_bear'   => '梦熊',
            'hope_dragon'  => '望龙',
        ];
    }

    public function currentStage(): array
    {
        return static::evolutionStages()[$this->level] ?? static::evolutionStages()[0];
    }

    // ========== 升级逻辑 ==========

    public function experienceForNextLevel(): int
    {
        // 每级需要的经验递增
        return ($this->level + 1) * 100;
    }

    public function canLevelUp(): bool
    {
        return $this->level < 10 && $this->experience >= $this->experienceForNextLevel();
    }

    public function levelUp(): bool
    {
        if (!$this->canLevelUp()) {
            return false;
        }

        $this->experience -= $this->experienceForNextLevel();
        $this->level += 1;
        $this->save();

        return true;
    }

    public function addExperience(int $amount): void
    {
        $this->experience += $amount;
        $this->save();
        // 自动检查升级
        while ($this->canLevelUp()) {
            $this->levelUp();
        }
    }

    // ========== 心情系统 ==========

    public function feed(): void
    {
        $this->mood = min(100, $this->mood + 20);
        $this->last_fed_at = now();
        $this->save();
    }

    public function decayMood(): void
    {
        // 每24小时未喂养，心情下降10点
        if ($this->last_fed_at && $this->last_fed_at->diffInHours(now()) >= 24) {
            $this->mood = max(0, $this->mood - 10);
            $this->save();
        }
    }
}
