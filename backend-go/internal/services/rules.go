// 积分规则服务：默认规则模板 + 分类标签 + 学校级默认规则幂等播种。
package services

import (
	"sort"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// DefaultRule 是默认积分规则模板项。
type DefaultRule struct {
	Name       string
	Amount     int
	Category   string
	IsPositive bool
}

// DefaultRules 是学校级默认积分规则（唯一真源，忠实移植自 ScoreRuleService::DEFAULT_RULES，共 43 条）。
var DefaultRules = []DefaultRule{
	// 📖 课堂表现
	{"举手发言", 3, "classroom", true},
	{"认真听讲", 2, "classroom", true},
	{"积极互动", 3, "classroom", true},
	{"课堂专注", 2, "classroom", true},
	{"挑战难题", 5, "classroom", true},
	{"上课走神", -2, "classroom", false},
	{"打扰课堂", -3, "classroom", false},
	{"追逐打闹", -5, "classroom", false},
	{"课堂喧哗", -2, "classroom", false},
	{"趴桌睡觉", -2, "classroom", false},
	// 📝 作业管理
	{"作业优秀", 5, "homework", true},
	{"作业按时完成", 3, "homework", true},
	{"作业有进步", 3, "homework", true},
	{"书写工整", 2, "homework", true},
	{"作业缺交", -3, "homework", false},
	{"作业敷衍", -2, "homework", false},
	{"作业迟交", -1, "homework", false},
	// 🌟 行为习惯
	{"遵守纪律", 2, "behavior", true},
	{"帮助同学", 4, "behavior", true},
	{"诚实守信", 5, "behavior", true},
	{"拾金不昧", 5, "behavior", true},
	{"尊敬师长", 3, "behavior", true},
	{"说脏话", -3, "behavior", false},
	{"与同学冲突", -5, "behavior", false},
	{"撒谎欺骗", -5, "behavior", false},
	// 📊 综合素养
	{"科技创新", 8, "literacy", true},
	{"阅读之星", 5, "literacy", true},
	{"体育锻炼", 3, "literacy", true},
	{"艺术表现", 5, "literacy", true},
	{"劳动积极", 3, "literacy", true},
	{"节约环保", 2, "literacy", true},
	{"竞赛获奖", 10, "literacy", true},
	{"破坏公物", -8, "literacy", false},
	{"乱扔垃圾", -2, "literacy", false},
	// 📅 日常表现
	{"全勤表现", 3, "daily", true},
	{"按时到校", 2, "daily", true},
	{"迟到早退", -2, "daily", false},
	{"仪容整洁", 1, "daily", true},
	{"值日认真", 2, "daily", true},
	{"不戴红领巾/校牌", -1, "daily", false},
	// 📚 学业表现
	{"考试优秀", 10, "academic", true},
	{"考试进步", 8, "academic", true},
	{"考试作弊", -15, "academic", false},
}

// CategoryLabels 是分类标签（后端唯一真源，键集与前端 categoryLabels 一致）。
var CategoryLabels = map[string]string{
	"classroom": "课堂表现",
	"homework":  "作业管理",
	"behavior":  "行为习惯",
	"literacy":  "综合素养",
	"daily":     "日常表现",
	"academic":  "学业表现",
	"custom":    "自定义",
}

// categoryOrder 分类键序（逐字同 Laravel ScoreRuleService::CATEGORY_LABELS 的声明顺序）。
// Go 的 map 无序，故 `GET /common/score-categories` 的 `sort` 需要单独一份顺序表；
// 名称本身仍取 CategoryLabels（唯一真源），此处**只承载顺序、不承载文案**。
var categoryOrder = []string{"classroom", "homework", "behavior", "literacy", "daily", "academic", "custom"}

// categoryIcons 分类图标（逐字同 Laravel StudentController::scoreCategories 的 $icons）。
var categoryIcons = map[string]string{
	"classroom": "📖",
	"homework":  "📝",
	"behavior":  "🌟",
	"literacy":  "📊",
	"daily":     "📅",
	"academic":  "📚",
	"custom":    "✨",
}

// ScoreCategory 积分分类字典项（同 Laravel `/common/score-categories` 的元素形状）。
type ScoreCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int    `json:"sort"`
}

// ScoreCategories 返回积分分类字典（Laravel StudentController::scoreCategories 的 Go 版）。
//
// id / name 一律来自 CategoryLabels（后端唯一真源，见项目决策 23），图标缺项兜底 📌，
// sort 从 1 递增。若将来 CategoryLabels 新增了键而 categoryOrder 未同步，这里会把这部分键
// 追加在末尾（按键名排序）而不是静默丢弃——只影响 `sort` 次序，不会丢分类。
func ScoreCategories() []ScoreCategory {
	out := make([]ScoreCategory, 0, len(CategoryLabels))
	seen := map[string]bool{}
	// sort 用自增计数器（同 Laravel `++$sort`）：顺序表里若有键不在 CategoryLabels 中，
	// 被 `continue` 跳过后也不会留下编号空档。
	sortNo := 0
	for _, id := range categoryOrder {
		name, ok := CategoryLabels[id]
		if !ok {
			continue
		}
		seen[id] = true
		sortNo++
		out = append(out, ScoreCategory{ID: id, Name: name, Icon: categoryIcon(id), Sort: sortNo})
	}

	// 防御性兜底：CategoryLabels 里不在 categoryOrder 中的键（顺序表漏更新时仍然输出）。
	rest := make([]string, 0, len(CategoryLabels))
	for id := range CategoryLabels {
		if !seen[id] {
			rest = append(rest, id)
		}
	}
	sort.Strings(rest)
	for _, id := range rest {
		sortNo++
		out = append(out, ScoreCategory{ID: id, Name: CategoryLabels[id], Icon: categoryIcon(id), Sort: sortNo})
	}
	return out
}

// categoryIcon 取分类图标，未登记时兜底 📌（同 Laravel `$icons[$id] ?? '📌'`）。
func categoryIcon(id string) string {
	if icon, ok := categoryIcons[id]; ok {
		return icon
	}
	return "📌"
}

// Rules 积分规则服务。
type Rules struct {
	db *gorm.DB
}

// NewRules 创建规则服务。
func NewRules(db *gorm.DB) *Rules {
	return &Rules{db: db}
}

// EnsureDefaultsForSchool 增量补齐某校的学校级默认规则（幂等，同名不覆盖）。
// 判据是「该校是否播种过」而非「规则是否为空」，与 Laravel 端决策 23 一致。
// 返回本次新建的规则条数。
func (r *Rules) EnsureDefaultsForSchool(schoolID uint) (int, error) {
	if schoolID == 0 {
		return 0, nil
	}

	var school models.School
	if err := r.db.First(&school, schoolID).Error; err != nil {
		return 0, err
	}
	if school.ScoreRulesSeeded {
		return 0, nil
	}

	created := 0
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i, d := range DefaultRules {
			var count int64
			if err := tx.Model(&models.ScoreRule{}).
				Where("class_id IS NULL AND school_id = ? AND name = ?", schoolID, d.Name).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				continue
			}

			rule := models.ScoreRule{
				SchoolID:   uintPtr(schoolID),
				ClassID:    nil,
				Name:       d.Name,
				Amount:     d.Amount,
				Category:   d.Category,
				IsPositive: d.IsPositive,
				IsActive:   true,
				SortOrder:  i,
			}
			if err := tx.Create(&rule).Error; err != nil {
				return err
			}
			// is_positive 带 default:true，惩罚规则（false）会被 GORM 跳过而落库成 true，需显式回写。
			if err := forceFalseBool(tx, &models.ScoreRule{}, rule.ID, "is_positive", d.IsPositive); err != nil {
				return err
			}
			rule.IsPositive = d.IsPositive
			created++
		}

		return tx.Model(&models.School{}).Where("id = ?", schoolID).
			Update("score_rules_seeded", true).Error
	})
	if err != nil {
		return 0, err
	}
	return created, nil
}

// ListForTeacher 教师端规则列表：本班班级规则 + 本校学校级规则（class_id=NULL）。
// 首次访问时增量补齐学校级默认规则。
func (r *Rules) ListForTeacher(u *models.User, classIDs []uint) ([]models.ScoreRule, error) {
	if _, err := r.EnsureDefaultsForSchool(u.SchoolID); err != nil {
		return nil, err
	}

	var rules []models.ScoreRule
	q := r.db.Where("(class_id IN ?) OR (class_id IS NULL AND school_id = ?)", classIDs, u.SchoolID)
	if err := q.Order("sort_order ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// ListForClass 教室端（班级码）规则列表：本班班级规则 + 本校学校级规则。
// 忠实移植自 Laravel ScoreRuleService::rulesForClass（含首次访问补齐本校默认规则）。
// 排序在 Laravel 为 sort_order；Go 端追加 id 作为并列时的稳定次序。
func (r *Rules) ListForClass(classID, schoolID uint) ([]models.ScoreRule, error) {
	if schoolID != 0 {
		if _, err := r.EnsureDefaultsForSchool(schoolID); err != nil {
			return nil, err
		}
	}

	var rules []models.ScoreRule
	q := r.db.Where("class_id = ? OR (class_id IS NULL AND school_id = ?)", classID, schoolID)
	if err := q.Order("sort_order ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// Create 创建规则（默认落到教师首个可管理班级）。
func (r *Rules) Create(classID *uint, schoolID uint, name string, amount int, category string, isPositive bool) (*models.ScoreRule, error) {
	if category == "" {
		category = "custom"
	}
	rule := models.ScoreRule{
		SchoolID:   &schoolID,
		ClassID:    classID,
		Name:       name,
		Amount:     amount,
		Category:   category,
		IsPositive: isPositive,
		IsActive:   true,
		SortOrder:  0,
	}
	if err := r.db.Create(&rule).Error; err != nil {
		return nil, err
	}
	// is_positive 带 default:true，惩罚规则（false）需创建后显式回写。
	if err := forceFalseBool(r.db, &models.ScoreRule{}, rule.ID, "is_positive", isPositive); err != nil {
		return nil, err
	}
	rule.IsPositive = isPositive
	return &rule, nil
}

// FindScopedForTeacher 教师可见范围内的单条规则（跨校/越权 ID 返回 404）。
func (r *Rules) FindScopedForTeacher(u *models.User, classIDs []uint, id uint) (*models.ScoreRule, error) {
	var rule models.ScoreRule
	err := r.db.Where("id = ? AND ((class_id IN ?) OR (class_id IS NULL AND school_id = ?))",
		id, classIDs, u.SchoolID).First(&rule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound("规则不存在或不可见")
		}
		return nil, err
	}
	return &rule, nil
}

// Update 更新规则字段并返回最新实体。
func (r *Rules) Update(rule *models.ScoreRule, updates map[string]any) (*models.ScoreRule, error) {
	if err := r.db.Model(rule).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(rule, rule.ID).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

// Delete 删除规则。
func (r *Rules) Delete(rule *models.ScoreRule) error {
	return r.db.Delete(rule).Error
}

func uintPtr(v uint) *uint { return &v }

// ============================================================
// 管理员端：全校积分规则（学校级 school_id = 本校，class_id = null）
// 移植自 SchoolAdminController::adminListScoreRules 等四个动作。
// ============================================================

// AdminRuleView 管理员端规则视图（附班级名与 scope 标注）。
type AdminRuleView struct {
	ID         uint      `json:"id"`
	SchoolID   *uint     `json:"school_id"`
	ClassID    *uint     `json:"class_id"`
	Name       string    `json:"name"`
	Amount     int       `json:"amount"`
	Category   string    `json:"category"`
	IsPositive bool      `json:"is_positive"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
	ClassName  *string   `json:"class_name"`
	Scope      string    `json:"scope"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ListForSchool 管理员端规则列表：全校学校级规则 + 教师创建的班级级规则。
// 「装完即用」：管理员常是最先打开本页的人，同样触发默认规则补齐（幂等）。
func (r *Rules) ListForSchool(schoolID uint) ([]AdminRuleView, error) {
	if _, err := r.EnsureDefaultsForSchool(schoolID); err != nil {
		return nil, err
	}

	var rules []models.ScoreRule
	if err := r.db.Where("school_id = ?", schoolID).
		Order("sort_order ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}

	classNames, err := r.classNamesOfSchool(schoolID)
	if err != nil {
		return nil, err
	}

	views := []AdminRuleView{}
	for _, rule := range rules {
		v := AdminRuleView{
			ID: rule.ID, SchoolID: rule.SchoolID, ClassID: rule.ClassID,
			Name: rule.Name, Amount: rule.Amount, Category: rule.Category,
			IsPositive: rule.IsPositive, IsActive: rule.IsActive, SortOrder: rule.SortOrder,
			Scope: "school", CreatedAt: rule.CreatedAt, UpdatedAt: rule.UpdatedAt,
		}
		if rule.ClassID != nil {
			v.Scope = "class"
			if name, ok := classNames[*rule.ClassID]; ok {
				v.ClassName = &name
			}
		}
		views = append(views, v)
	}
	return views, nil
}

// CreateForSchool 创建学校级规则（class_id = null，全校共享）。
func (r *Rules) CreateForSchool(schoolID uint, name string, amount int, category string, isPositive, isActive bool) (*models.ScoreRule, error) {
	if category == "" {
		category = "custom"
	}
	rule := models.ScoreRule{
		SchoolID:   &schoolID,
		ClassID:    nil,
		Name:       name,
		Amount:     amount,
		Category:   category,
		IsPositive: isPositive,
		IsActive:   isActive,
		SortOrder:  0,
	}
	// is_active / is_positive 带 default:true，Create 时 GORM 会跳过零值 false，
	// 故创建后对 false 值显式回写（同 currency.go CreateRate 的做法）。
	if err := r.db.Create(&rule).Error; err != nil {
		return nil, err
	}
	if err := forceFalseBool(r.db, &models.ScoreRule{}, rule.ID, "is_active", isActive); err != nil {
		return nil, err
	}
	if err := forceFalseBool(r.db, &models.ScoreRule{}, rule.ID, "is_positive", isPositive); err != nil {
		return nil, err
	}
	rule.IsActive = isActive
	rule.IsPositive = isPositive
	return &rule, nil
}

// FindSchoolLevel 取本校学校级规则（班级级或跨校 ID 一律 404）。
func (r *Rules) FindSchoolLevel(schoolID, id uint) (*models.ScoreRule, error) {
	var rule models.ScoreRule
	err := r.db.Where("id = ? AND school_id = ? AND class_id IS NULL", id, schoolID).First(&rule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound("规则不存在或不可见")
		}
		return nil, err
	}
	return &rule, nil
}

// classNamesOfSchool 返回本校班级 id → 名称。
func (r *Rules) classNamesOfSchool(schoolID uint) (map[uint]string, error) {
	var classes []models.ClassRoom
	if err := r.db.Where("school_id = ?", schoolID).Find(&classes).Error; err != nil {
		return nil, err
	}
	names := map[uint]string{}
	for _, c := range classes {
		names[c.ID] = c.Name
	}
	return names, nil
}

// forceFalseBool 处理「带 default:true 的布尔列在 Create 时被 GORM 跳过零值、
// 由 DB 施加默认 true」的问题：仅当目标值为 false 时补一次显式更新。
// 与 currency.go 中 CreateRate 的既有做法一致（Update 单列会写入零值）。
func forceFalseBool(db *gorm.DB, model any, id uint, column string, value bool) error {
	if value {
		return nil
	}
	return db.Model(model).Where("id = ?", id).Update(column, false).Error
}
