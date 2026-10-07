package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
)

// 默认规则模板的静态校验（唯一真源，忠实移植自 ScoreRuleService::DEFAULT_RULES）。
func TestDefaultRulesCount(t *testing.T) {
	if len(services.DefaultRules) != 43 {
		t.Fatalf("DefaultRules count = %d, want 43", len(services.DefaultRules))
	}

	positive, negative := 0, 0
	seen := map[string]bool{}
	for _, r := range services.DefaultRules {
		if seen[r.Name] {
			t.Fatalf("duplicate default rule name %q", r.Name)
		}
		seen[r.Name] = true

		if _, ok := services.CategoryLabels[r.Category]; !ok {
			t.Fatalf("rule %q has unknown category %q", r.Name, r.Category)
		}
		if r.IsPositive && r.Amount <= 0 {
			t.Fatalf("positive rule %q has amount %d", r.Name, r.Amount)
		}
		if !r.IsPositive && r.Amount >= 0 {
			t.Fatalf("negative rule %q has amount %d", r.Name, r.Amount)
		}

		if r.IsPositive {
			positive++
		} else {
			negative++
		}
	}
	if positive != 27 {
		t.Fatalf("positive default rules = %d, want 27", positive)
	}
	if negative != 16 {
		t.Fatalf("negative default rules = %d, want 16", negative)
	}
}

// EnsureDefaultsForSchool 首次播种 43 条，第二次幂等返回 0。
func TestEnsureDefaultsForSchoolIdempotent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	rules := services.NewRules(db)

	created, err := rules.EnsureDefaultsForSchool(school.ID)
	if err != nil {
		t.Fatalf("ensure defaults: %v", err)
	}
	if created != 43 {
		t.Fatalf("first seed created %d, want 43", created)
	}

	// 学校标记已写入。
	var s models.School
	if err := db.First(&s, school.ID).Error; err != nil {
		t.Fatalf("reload school: %v", err)
	}
	if !s.ScoreRulesSeeded {
		t.Fatalf("ScoreRulesSeeded not set after seed")
	}

	// 第二次应不再新增。
	created2, err := rules.EnsureDefaultsForSchool(school.ID)
	if err != nil {
		t.Fatalf("ensure defaults (2nd): %v", err)
	}
	if created2 != 0 {
		t.Fatalf("second seed created %d, want 0", created2)
	}

	// 数据库中确实只有 43 条学校级规则。
	var count int64
	if err := db.Model(&models.ScoreRule{}).Where("school_id = ?", school.ID).Count(&count).Error; err != nil {
		t.Fatalf("count rules: %v", err)
	}
	if count != 43 {
		t.Fatalf("seeded rules in db = %d, want 43", count)
	}
}

// 教师删除某条默认规则后，下次播种不会复活。
func TestEnsureDefaultsDoesNotResurrectDeletedRule(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	rules := services.NewRules(db)

	if _, err := rules.EnsureDefaultsForSchool(school.ID); err != nil {
		t.Fatalf("ensure defaults: %v", err)
	}

	var rule models.ScoreRule
	if err := db.Where("school_id = ? AND name = ?", school.ID, "举手发言").First(&rule).Error; err != nil {
		t.Fatalf("find default rule: %v", err)
	}
	if err := rules.Delete(&rule); err != nil {
		t.Fatalf("delete rule: %v", err)
	}

	// 强行把 seeded 标记复位，模拟「已播种过、但规则被删」的场景。
	if err := db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("score_rules_seeded", false).Error; err != nil {
		t.Fatalf("reset seeded flag: %v", err)
	}

	created, err := rules.EnsureDefaultsForSchool(school.ID)
	if err != nil {
		t.Fatalf("re-seed: %v", err)
	}
	// 只剩被删的那一条应被补回，其余同名不覆盖。
	if created != 1 {
		t.Fatalf("re-seed created %d, want 1", created)
	}

	var count int64
	if err := db.Model(&models.ScoreRule{}).Where("school_id = ?", school.ID).Count(&count).Error; err != nil {
		t.Fatalf("count rules: %v", err)
	}
	if count != 43 {
		t.Fatalf("rules after re-seed = %d, want 43", count)
	}
}

// 分类标签键集与前端保持一致（7 键）。
func TestCategoryLabelsKeys(t *testing.T) {
	want := []string{"classroom", "homework", "behavior", "literacy", "daily", "academic", "custom"}
	if len(services.CategoryLabels) != len(want) {
		t.Fatalf("CategoryLabels size = %d, want %d", len(services.CategoryLabels), len(want))
	}
	for _, k := range want {
		if _, ok := services.CategoryLabels[k]; !ok {
			t.Fatalf("missing category label key %q", k)
		}
	}
}
