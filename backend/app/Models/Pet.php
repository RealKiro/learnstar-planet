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
        'type',
        'level',
        'experience',
        'mood',
        'accessories',
        'last_fed_at',
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
            'star_fish'    => '星鱼',
            'moon_moth'    => '月蝶',
            'sun_sparrow'  => '日雀',
            'cloud_whale'  => '云鲸',
            'pikachu'      => '皮卡丘',
            'eevee'        => '伊布',
            'charmander'   => '小火龙',
            'squirtle'     => '杰尼龟',
            'bulbasaur'    => '妙蛙种子',
            'snorlax'      => '卡比兽',
            'mewtwo'       => '超梦',
            'dragonite'    => '快龙',
            'venusaur'     => '妙蛙花',
            'blastoise'    => '水箭龟',
            'charizard'    => '喷火龙',
            'gengar'       => '耿鬼',
            'orange_cat'   => '橘猫',
            'husky'        => '哈士奇',
            'shiba'        => '柴犬',
            'guinea_pig'   => '荷兰猪',
            'hamster'      => '仓鼠',
            'bunny'        => '兔子',
            'parrot'       => '鹦鹉',
            'hedgehog'     => '刺猬',
            'chinchilla'   => '龙猫',
            'teacup_pig'   => '小香猪',
            'sugar_glider' => '蜜袋鼯',
            'alpaca'       => '羊驼',
            'panda'           => '大熊猫',
            'golden_monkey'   => '金丝猴',
            'crested_ibis'    => '朱鹮',
            'snow_leopard'    => '雪豹',
            'tibetan_antelope' => '藏羚羊',
            'red_crowned_crane' => '丹顶鹤',
            'milu_deer'       => '麋鹿',
            'chinese_alligator' => '扬子鳄',
            'siberian_tiger'  => '东北虎',
            'pangolin'        => '穿山甲',
            'red_panda'       => '小熊猫',
            'finless_porpoise' => '江豚',
            'qilin'    => '麒麟',
            'fenghuang' => '凤凰',
            'baihu'    => '白虎',
            'xuanwu'   => '玄武',
            'zhuque'   => '朱雀',
            'taotie'   => '饕餮',
            'pixiu'    => '貔貅',
            'kunpeng'  => '鲲鹏',
            'qinglong' => '青龙',
            'qiongqi'  => '穷奇',
            'hundun'   => '混沌',
            'zhulong'  => '烛龙',
        ];
    }

    public static function petCategories(): array
    {
        return [
            'cosmic'  => '原创宇宙系列',
            'pokemon' => '宝可梦系列',
            'cute'    => '萌宠系列',
            'treasure' => '国宝系列',
            'mythic'  => '神兽系列',
        ];
    }

    public static function petTypesBySeries(string $series): array
    {
        $all = self::petTypes();
        $filtered = [];
        foreach ($all as $type => $name) {
            if (self::getCategoryByType($type) === $series) {
                $filtered[$type] = $name;
            }
        }

        return $filtered;
    }

    public static function getCategoryByType(string $type): string
    {
        $categories = [
            'pokemon' => ['pikachu', 'eevee', 'charmander', 'squirtle', 'bulbasaur', 'snorlax', 'mewtwo', 'dragonite', 'venusaur', 'blastoise', 'charizard', 'gengar'],
            'cute'    => ['orange_cat', 'husky', 'shiba', 'guinea_pig', 'hamster', 'bunny', 'parrot', 'hedgehog', 'chinchilla', 'teacup_pig', 'sugar_glider', 'alpaca'],
            'treasure' => ['panda', 'golden_monkey', 'crested_ibis', 'snow_leopard', 'tibetan_antelope', 'red_crowned_crane', 'milu_deer', 'chinese_alligator', 'siberian_tiger', 'pangolin', 'red_panda', 'finless_porpoise'],
            'mythic'  => ['qilin', 'fenghuang', 'baihu', 'xuanwu', 'zhuque', 'taotie', 'pixiu', 'kunpeng', 'qinglong', 'qiongqi', 'hundun', 'zhulong'],
        ];

        foreach ($categories as $category => $types) {
            if (in_array($type, $types, true)) {
                return $category;
            }
        }

        return 'cosmic';
    }

    public function currentStage(): array
    {
        return static::evolutionStages()[$this->level] ?? static::evolutionStages()[0];
    }

    public function experienceForNextLevel(): int
    {
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

        while ($this->canLevelUp()) {
            $this->levelUp();
        }
    }

    public function removeExperience(int $amount): void
    {
        $this->experience -= $amount;
        $this->save();

        while ($this->canLevelDown()) {
            $this->levelDown();
        }

        if ($this->level === 0 && $this->experience < 0) {
            $this->experience = 0;
            $this->save();
        }
    }

    public function canLevelDown(): bool
    {
        return $this->level > 0 && $this->experience < 0;
    }

    public function levelDown(): bool
    {
        if (!$this->canLevelDown()) {
            return false;
        }

        $this->level -= 1;
        $this->experience += $this->experienceForNextLevel();
        $this->save();

        return true;
    }

    public function feed(): void
    {
        $this->mood = min(100, $this->mood + 20);
        $this->last_fed_at = now();
        $this->save();
    }

    public function decayMood(): void
    {
        if ($this->last_fed_at && $this->last_fed_at->diffInHours(now()) >= 24) {
            $this->mood = max(0, $this->mood - 10);
            $this->save();
        }
    }
}
