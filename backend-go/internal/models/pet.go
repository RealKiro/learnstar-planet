// 宠物模型与等级/阶段业务规则。
// 阈值与阶段划分忠实移植自 Laravel 端 App\Models\Pet。
package models

import "time"

// Pet 学生宠物（与学生一对一）。字段照 Laravel App\Models\Pet 的 $fillable。
type Pet struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	StudentID  uint       `gorm:"uniqueIndex;not null" json:"student_id"`
	ClassID    uint       `gorm:"index;not null" json:"class_id"`
	Name       string     `gorm:"size:100;not null" json:"name"`
	Species    string     `gorm:"size:64;default:zhulong" json:"species"`
	Level      int        `gorm:"default:1" json:"level"`
	Experience int        `gorm:"default:0" json:"experience"`
	Mood       int        `gorm:"default:80" json:"mood"`
	LastFedAt  *time.Time `json:"last_fed_at"`
	// LastSwitchedAt 上次切换物种时间（Laravel `pets.last_switched_at`，Pet 的 datetime cast）。
	LastSwitchedAt *time.Time `json:"last_switched_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName 表名（Laravel 迁移 2026_08_06_000005 的 pets 表）。
func (Pet) TableName() string { return "pets" }

// PetCollection 宠物图鉴收藏（每个学生 × 物种一条，记录该物种的养成进度）。
// 表名 pet_collections，字段照 Laravel App\Models\PetCollection 的 $fillable。
// 说明：教师端 / 教室端的「切换宠物」都会在此归档旧物种进度并在切回时恢复
// （同 Laravel PetService::switchPet / DisplayController::classroomSwitchPet）。
type PetCollection struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StudentID  uint      `gorm:"index;not null;uniqueIndex:idx_pet_collection_student_species" json:"student_id"`
	Species    string    `gorm:"size:64;not null;uniqueIndex:idx_pet_collection_student_species" json:"species"`
	Level      int       `gorm:"default:1" json:"level"`
	Experience int       `gorm:"default:0" json:"experience"`
	Mood       int       `gorm:"default:80" json:"mood"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UnlockSlotsForScore 已解锁的图鉴槽位数（每 100 积分 +1，初始 1），
// 逐字同 Laravel PetCollection::unlockSlotsForScore。
func UnlockSlotsForScore(totalScore int) int {
	if totalScore < 0 {
		totalScore = 0
	}
	return 1 + totalScore/100
}

// LevelThresholds 积分 → 等级阈值（Lv.0~12，索引即等级）。
// 阈值与前端 petData 统一：0/0/15/35/60/90/125/165/210/260/315/375/450。
var LevelThresholds = [13]int{0, 0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450}

// MaxLevel 宠物最高等级。
const MaxLevel = 12

// Stage 宠物成长阶段。
type Stage struct {
	Stage  string `json:"stage"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Emoji  string `json:"emoji"`
	Color  string `json:"color"`
	Level  int    `json:"level"`
	ExpMax int    `json:"exp_max"`
}

// CurrentStage 返回指定等级对应的成长阶段（Lv.1-2 蛋 / 3-4 幼年 / 5-6 成长 / 7-8 成熟 / 9-10 传说 / 11-12 归真）。
func (p *Pet) CurrentStage() Stage {
	stage, name, title, emoji, color := "transcendent", "归真级", "返璞归真", "👑", "#F472B6"
	switch {
	case p.Level <= 2:
		stage, name, title, emoji, color = "egg", "新生之卵", "破壳新生", "🥚", "#F59E0B"
	case p.Level <= 4:
		stage, name, title, emoji, color = "baby", "幼年", "蹒跚学步", "🐣", "#10B981"
	case p.Level <= 6:
		stage, name, title, emoji, color = "growing", "成长期", "茁壮成长", "🌱", "#3B82F6"
	case p.Level <= 8:
		stage, name, title, emoji, color = "mature", "成熟期", "英姿勃发", "🌿", "#8B5CF6"
	case p.Level <= 10:
		stage, name, title, emoji, color = "legendary", "传说级", "超凡入圣", "🌟", "#F59E0B"
	}
	return Stage{
		Stage:  stage,
		Name:   name,
		Title:  title,
		Emoji:  emoji,
		Color:  color,
		Level:  p.Level,
		ExpMax: maxInt(1, (p.Level+1)*10),
	}
}

// ExperienceForNextLevel 当前等级升下一级所需经验（(level+1)*10）。
func (p *Pet) ExperienceForNextLevel() int { return (p.Level + 1) * 10 }

// LevelFromScore 根据总积分计算等级（与 syncLevelWithScore 阈值一致）。
func LevelFromScore(score int) int {
	newLevel := 0
	for level, threshold := range LevelThresholds {
		if score >= threshold {
			newLevel = level
		}
	}
	return newLevel
}

// SyncLevelWithScore 根据学生总积分同步宠物等级（积分即成长值）。
func (p *Pet) SyncLevelWithScore(score int) {
	p.Level = LevelFromScore(score)
}

// SwitchCost 更换宠物所需积分：等级越高越贵（Lv.1=5 ... Lv.12=60）。
func SwitchCost(level int) int {
	return 5 * maxInt(1, level)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// seriesSpeciesPools 各系列物种池（仅物种 id），逐项复制自 Laravel Pet::speciesPoolForSeries，
// 用于「整班切换系列」与「无宠物学生自动分配」时的取种。
// 各系列物种数不等：myth 18 / pokemon 12 / national 12 / digimon 6 / magic 12 /
// prehistoric 12 / constellation 12 / festival 12 / qixia 10 / dongfang 20。
var seriesSpeciesPools = map[string][]string{
	"myth":          {"zhulong", "yinglong", "nine_tail_fox", "kunpeng", "fenghuang", "qilin", "qinglong", "baihu", "zhuque", "xuanwu", "taotie", "baize", "qiongqi", "bifang", "pixiu", "jingwei", "xiangliu", "xiezhi"},
	"pokemon":       {"charmander", "bulbasaur", "squirtle", "eevee", "pikachu", "riolu", "ice_fox", "rock_rhino", "wind_falcon", "light_deer", "dark_panther", "steel_armadillo"},
	"national":      {"panda", "golden_monkey", "red_crowned_crane", "south_china_tiger", "chinese_alligator", "crested_ibis", "tibetan_antelope", "snow_leopard", "milu_deer", "siberian_tiger", "red_panda", "finless_porpoise"},
	"digimon":       {"mecha_dragon", "cyber_cat", "space_mecha", "quantum_beast", "digital_phoenix", "mecha_shark"},
	"magic":         {"unicorn", "wyvern", "fairy", "treant", "griffin", "mermaid", "grey_wizard", "wand_cat", "dragon_knight", "alchemy_golem", "nightmare_horse", "lamp_spirit"},
	"prehistoric":   {"t_rex", "triceratops", "pterosaur", "mammoth", "sabertooth", "mosasaur", "spinosaurus", "ankylosaurus", "diplodocus", "megalodon", "ground_sloth", "woolly_rhino"},
	"constellation": {"aries", "taurus", "gemini", "cancer", "leo", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces"},
	"festival":      {"zongzi", "tangyuan", "mooncake", "qingtuan", "chongyang_cake", "niangao", "laba_porridge", "spring_pancake", "tanghulu", "osmanthus_cake", "wonton", "festival_lantern"},
	"qixia":         {"hongmao", "lantu", "doudou", "dabeng", "tiaotiao", "shali", "dada", "heixinhu", "heixiaohu", "ma_sanniang"},
	"dongfang":      {"jiang_ziya", "nezha", "yang_jian", "lei_zhenzi", "huang_tianhua", "tu_xingsun", "yang_ren", "wei_hu", "daji", "shen_gongbao", "sun_wukong", "lv_dongbin", "he_xiangu", "zhang_guolao", "tie_guaili", "han_zhongli", "lan_caihe", "cao_guojiu", "taishang_laojun", "zhong_kui"},
}

// seriesLabels 系列中文名（同 Laravel Pet::seriesLabel；未知系列原样返回系列 id）。
var seriesLabels = map[string]string{
	"myth":          "山海经",
	"pokemon":       "宝可梦",
	"national":      "国宝守护",
	"digimon":       "数码宝贝",
	"magic":         "魔法奇幻",
	"prehistoric":   "史前生物",
	"constellation": "星座守护",
	"festival":      "传统节日",
	"qixia":         "虹猫蓝兔七侠传",
	"dongfang":      "东方神话",
}

// SeriesIDs 合法系列 id 白名单（逐字同 Laravel 两个控制器内联的 $validSeries：
// DisplayController::classroomSwitchSeries 与 TeacherController::switchSeries）。
var SeriesIDs = []string{
	"myth", "pokemon", "national", "digimon", "magic",
	"prehistoric", "constellation", "festival", "qixia", "dongfang",
}

// IsValidSeries 判断系列 id 是否在白名单内。
func IsValidSeries(seriesID string) bool {
	for _, id := range SeriesIDs {
		if id == seriesID {
			return true
		}
	}
	return false
}

// SeriesLabel 系列中文名（同 Laravel Pet::seriesLabel；未知系列原样返回）。
func SeriesLabel(seriesID string) string {
	if label, ok := seriesLabels[seriesID]; ok {
		return label
	}
	return seriesID
}

// SpeciesPoolForSeries 返回某系列的可选物种 id（未知系列返回空表，同 Laravel）。
func SpeciesPoolForSeries(seriesID string) []string {
	return seriesSpeciesPools[seriesID]
}

// AddExperience 增加经验并处理连续升级（与 Laravel Pet::addExperience 一致）。
func (p *Pet) AddExperience(amount int) {
	p.Experience += amount
	for p.Level < MaxLevel && p.Experience >= p.ExperienceForNextLevel() {
		p.Experience -= p.ExperienceForNextLevel()
		p.Level++
	}
}

// RemoveExperience 扣除经验并处理连续降级（与 Laravel Pet::removeExperience 一致）。
func (p *Pet) RemoveExperience(amount int) {
	p.Experience -= amount
	for p.Level > 0 && p.Experience < 0 {
		p.Level--
		p.Experience += p.ExperienceForNextLevel()
	}
	if p.Level == 0 && p.Experience < 0 {
		p.Experience = 0
	}
}

// Feed 喂养宠物：心情 +20，上限 100，并记录喂养时间。
func (p *Pet) Feed(now time.Time) {
	p.Mood = minInt(100, p.Mood+20)
	t := now
	p.LastFedAt = &t
}

// DecayMood 距上次喂养满 24 小时则心情 -10，下限 0（与 Laravel Pet::decayMood 一致）。
func (p *Pet) DecayMood(now time.Time) {
	if p.LastFedAt != nil && now.Sub(*p.LastFedAt) >= 24*time.Hour {
		p.Mood = maxInt(0, p.Mood-10)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
