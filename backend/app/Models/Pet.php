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
            0  => ['emoji' => '🥚', 'name' => '宠物蛋',   'title' => '等待破壳',  'color' => '#94a3b8', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=egg&backgroundColor=b8c5d6'],
            1  => ['emoji' => '🐣', 'name' => '绒绒雏鸟', 'title' => '破壳新生',  'color' => '#fcd34d', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=chick&backgroundColor=fef3c7'],
            2  => ['emoji' => '🐥', 'name' => '黄毛小鸭', 'title' => '蹒跚学步',  'color' => '#fbbf24', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=duck&backgroundColor=fef3c7'],
            3  => ['emoji' => '🐰', 'name' => '绒耳萌兔', 'title' => '活泼好动',  'color' => '#f472b6', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=rabbit&backgroundColor=fce7f3'],
            4  => ['emoji' => '🦊', 'name' => '机灵小狐', 'title' => '机智灵敏',  'color' => '#fb923c', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=fox&backgroundColor=fff7ed'],
            5  => ['emoji' => '🐱', 'name' => '优雅萌猫', 'title' => '优雅从容',  'color' => '#a78bfa', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=cat&backgroundColor=f5f3ff'],
            6  => ['emoji' => '🐶', 'name' => '忠诚幼犬', 'title' => '忠诚守护',  'color' => '#60a5fa', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=dog&backgroundColor=eff6ff'],
            7  => ['emoji' => '🦁', 'name' => '威风小狮', 'title' => '王者气度',  'color' => '#f59e0b', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=lion&backgroundColor=fef3c7'],
            8  => ['emoji' => '🐯', 'name' => '勇猛虎崽', 'title' => '勇往直前',  'color' => '#f97316', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=tiger&backgroundColor=fff7ed'],
            9  => ['emoji' => '🦄', 'name' => '神圣灵兽', 'title' => '超凡脱俗',  'color' => '#c084fc', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=unicorn&backgroundColor=f5f3ff'],
            10 => ['emoji' => '🐉', 'name' => '东方神龙', 'title' => '至尊无上',  'color' => '#ef4444', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=dragon&backgroundColor=fef2f2'],
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
        return static::evolutionStagesForType($this->type)[$this->level] ?? static::evolutionStages()[0];
    }

    /**
     * 获取某类宠物的专属进化链
     */
    public static function evolutionStagesForType(?string $type): array
    {
        $category = $type ? self::getCategoryByType($type) : 'cosmic';

        switch ($category) {
            case 'cute':
                return [
                    0  => ['emoji' => '🥚', 'name' => '宠物蛋',  'title' => '等待破壳',  'color' => '#94a3b8', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=egg&backgroundColor=b8c5d6'],
                    1  => ['emoji' => '🐹', 'name' => '小仓鼠',  'title' => '圆滚滚',    'color' => '#fcd34d', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=hamster&backgroundColor=fef3c7'],
                    2  => ['emoji' => '🐰', 'name' => '绒耳兔',  'title' => '蹦蹦跳',    'color' => '#f472b6', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=rabbit&backgroundColor=fce7f3'],
                    3  => ['emoji' => '🐹', 'name' => '荷兰猪',  'title' => '胖嘟嘟',    'color' => '#fb923c', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=guinea&backgroundColor=fff7ed'],
                    4  => ['emoji' => '🦔', 'name' => '小刺猬',  'title' => '竖起来',    'color' => '#a78bfa', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=hedgehog&backgroundColor=f5f3ff'],
                    5  => ['emoji' => '🐱', 'name' => '小橘猫',  'title' => '懒洋洋',    'color' => '#f97316', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=cat&backgroundColor=fff7ed'],
                    6  => ['emoji' => '🐶', 'name' => '哈士奇',  'title' => '二哈',      'color' => '#60a5fa', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=husky&backgroundColor=eff6ff'],
                    7  => ['emoji' => '🐕', 'name' => '柴犬',    'title' => '微笑天使',  'color' => '#f59e0b', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=shiba&backgroundColor=fef3c7'],
                    8  => ['emoji' => '🐑', 'name' => '小羊驼',  'title' => '毛茸茸',    'color' => '#c084fc', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=alpaca&backgroundColor=f5f3ff'],
                    9  => ['emoji' => '🦄', 'name' => '彩虹兽',  'title' => '梦幻',      'color' => '#ec4899', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=unicorn&backgroundColor=fce7f3'],
                    10 => ['emoji' => '🌟', 'name' => '精灵王',  'title' => '至尊',      'color' => '#fbbf24', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=fairy&backgroundColor=fef3c7'],
                ];

            case 'mythic':
                return [
                    0  => ['emoji' => '🥚', 'name' => '灵兽蛋',  'title' => '等待觉醒',   'color' => '#94a3b8', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=egg&backgroundColor=1e1b4b'],
                    1  => ['emoji' => '🐉', 'name' => '小龙仔',  'title' => '初生牛犊',   'color' => '#818cf8', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=small-dragon&backgroundColor=1e1b4b'],
                    2  => ['emoji' => '🦅', 'name' => '朱雀雏',  'title' => '浴火重生',   'color' => '#f87171', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=phoenix&backgroundColor=451a03'],
                    3  => ['emoji' => '🐯', 'name' => '白虎崽',  'title' => '王者初现',   'color' => '#e2e8f0', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=white-tiger&backgroundColor=1e293b'],
                    4  => ['emoji' => '🐢', 'name' => '玄武幼',  'title' => '沉稳如山',   'color' => '#4ade80', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=turtle&backgroundColor=052e16'],
                    5  => ['emoji' => '🦊', 'name' => '灵狐仙',  'title' => '九尾初显',   'color' => '#fb923c', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=nine-tail&backgroundColor=431407'],
                    6  => ['emoji' => '🐉', 'name' => '应龙',    'title' => '呼风唤雨',   'color' => '#fbbf24', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=dragon&backgroundColor=422006'],
                    7  => ['emoji' => '🦄', 'name' => '麒麟',    'title' => '祥瑞之兆',   'color' => '#c084fc', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=qilin&backgroundColor=2e1065'],
                    8  => ['emoji' => '🐲', 'name' => '神龙',    'title' => '腾云驾雾',   'color' => '#ef4444', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=god-dragon&backgroundColor=450a0a'],
                    9  => ['emoji' => '🌟', 'name' => '圣兽王',  'title' => '统御八方',   'color' => '#fbbf24', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=beast-king&backgroundColor=422006'],
                    10 => ['emoji' => '✨', 'name' => '混沌',    'title' => '开天辟地',   'color' => '#a855f7', 'image' => 'https://api.dicebear.com/9.x/icons/svg?seed=chaos&backgroundColor=3b0764'],
                ];

            default:
                return self::evolutionStages();
        }
    }

    public function experienceForNextLevel(): int
    {
        return ($this->level + 1) * 10;
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

