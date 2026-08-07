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
        'species',
        'type', // 兼容旧数据，新数据使用 species
        'level',
        'experience',
        'mood',
        'accessories',
        'last_fed_at',
        'last_switched_at',
    ];

    protected $casts = [
        'level' => 'integer',
        'experience' => 'integer',
        'mood' => 'integer',
        'accessories' => 'array',
        'last_fed_at' => 'datetime',
        'last_switched_at' => 'datetime',
    ];

    public function student(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(Student::class);
    }

    public function classRoom(): \Illuminate\Database\Eloquent\Relations\BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    /** 兼容前端 exp 字段（对应 experience 列） */
    public function getExpAttribute(): int
    {
        return (int) $this->experience;
    }

    /**
     * 当前阶段（按 12 级制划分，与前端 petData 阶段一致）
     * Lv.1-2 卵生 / 3-5 幼年 / 6-8 成长 / 9-10 成熟 / 11-12 传说
     */
    public function currentStage(): array
    {
        [$stage, $name, $title] = match (true) {
            $this->level <= 2 => ['egg', '新生之卵', '破壳新生'],
            $this->level <= 5 => ['baby', '幼年', '蹒跚学步'],
            $this->level <= 8 => ['growing', '成长期', '茁壮成长'],
            $this->level <= 10 => ['mature', '成熟期', '英姿勃发'],
            default => ['legendary', '传说级', '不朽传奇'],
        };

        return [
            'stage' => $stage,
            'name' => $name,
            'title' => $title,
            'emoji' => $this->currentStageEmoji(),
            'color' => match ($stage) {
                'egg' => '#F59E0B',
                'baby' => '#10B981',
                'growing' => '#3B82F6',
                'mature' => '#8B5CF6',
                default => '#F59E0B',
            },
            'level' => $this->level,
            'exp_max' => max(1, ($this->level + 1) * 10),
        ];
    }

    /** 阶段占位 emoji（SVG 化后仅作兜底） */
    public function currentStageEmoji(): string
    {
        return match (true) {
            $this->level <= 2 => '🥚',
            $this->level <= 5 => '🐣',
            $this->level <= 8 => '🌱',
            $this->level <= 10 => '🌟',
            default => '👑',
        };
    }

    /**
     * 获取某类宠物的专属进化链
     */
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
        $this->syncLevelWithScore();
    }

    /**
     * 根据学生总积分同步宠物等级（积分=成长值，积分够了等级自动升）
     * 阈值与前端 petData 统一: 0/0/15/35/60/90/125/165/210/260/315/375/450 (Lv.0~12)
     * 旧表有 50 个阈值会把等级顶到 49，前端只支持 12 级，必须对齐。
     */
    public function syncLevelWithScore(): void
    {
        $score = optional($this->student)->total_score ?? 0;
        $thresholds = [0, 0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450];
        $newLevel = 0;
        foreach ($thresholds as $level => $threshold) {
            if ($score >= $threshold) {
                $newLevel = $level;
            }
        }
        if ($newLevel !== $this->level) {
            $this->level = $newLevel;
            $this->save();
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
        $this->syncLevelWithScore();
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

    /**
     * 各系列 12 物种池（仅物种 id，与前端 petData 对齐；用于整班随机分配）
     */
    public static function speciesPoolForSeries(string $seriesId): array
    {
        return match ($seriesId) {
            'myth' => ['zhulong', 'yinglong', 'nine_tail_fox', 'kunpeng', 'fenghuang', 'qilin', 'qinglong', 'baihu', 'zhuque', 'xuanwu', 'taotie', 'baize', 'qiongqi', 'bifang', 'pixiu', 'jingwei', 'xiangliu', 'xiezhi'],
            'pokemon' => ['charmander', 'bulbasaur', 'squirtle', 'eevee', 'pikachu', 'riolu', 'ice_fox', 'rock_rhino', 'wind_falcon', 'light_deer', 'dark_panther', 'steel_armadillo'],
            'national' => ['panda', 'golden_monkey', 'red_crowned_crane', 'south_china_tiger', 'chinese_alligator', 'crested_ibis', 'tibetan_antelope', 'snow_leopard', 'milu_deer', 'siberian_tiger', 'red_panda', 'finless_porpoise'],
            'digimon' => ['mecha_dragon', 'cyber_cat', 'space_mecha', 'quantum_beast', 'digital_phoenix', 'mecha_shark', 'lightsaber_warrior', 'fission_giant', 'nano_swarm', 'storm_jet', 'bio_armor', 'starship_core'],
            'magic' => ['unicorn', 'wyvern', 'fairy', 'treant', 'griffin', 'mermaid', 'grey_wizard', 'wand_cat', 'dragon_knight', 'alchemy_golem', 'nightmare_horse', 'lamp_spirit'],
            'prehistoric' => ['t_rex', 'triceratops', 'pterosaur', 'mammoth', 'sabertooth', 'mosasaur', 'spinosaurus', 'ankylosaurus', 'diplodocus', 'megalodon', 'ground_sloth', 'woolly_rhino'],
            'constellation' => ['aries', 'taurus', 'gemini', 'cancer', 'leo', 'virgo', 'libra', 'scorpio', 'sagittarius', 'capricorn', 'aquarius', 'pisces'],
            'festival' => ['zongzi', 'tangyuan', 'mooncake', 'qingtuan', 'chongyang_cake', 'niangao', 'laba_porridge', 'spring_pancake', 'tanghulu', 'osmanthus_cake', 'wonton', 'festival_lantern'],
            'qixia' => ['hongmao', 'lantu', 'doudou', 'dabeng', 'tiaotiao', 'shali', 'dada', 'heixinhu', 'heixiaohu'],
            'fengshen' => ['jiang_ziya', 'nezha', 'yang_jian', 'lei_zhenzi', 'huang_tianhua', 'tu_xingsun', 'yang_ren', 'wei_hu', 'daji', 'shen_gongbao'],
            default => [],
        };
    }

    /**
     * 系列中文名（用于报错/提示文案，与前端 petData 系列名对齐）
     */
    public static function seriesLabel(string $seriesId): string
    {
        return match ($seriesId) {
            'myth' => '山海经',
            'pokemon' => '宝可梦风格',
            'national' => '国宝守护',
            'digimon' => '数码宝贝',
            'magic' => '魔法奇幻',
            'prehistoric' => '史前生物',
            'constellation' => '星座守护',
            'festival' => '传统节日',
            'qixia' => '七侠剑客',
            'fengshen' => '封神演义',
            default => $seriesId,
        };
    }

    /**
     * 更换宠物所需积分：等级越高越贵（Lv.1=5 ... Lv.12=60）
     */
    public static function switchCost(int $level): int
    {
        return 5 * max(1, $level);
    }
}
