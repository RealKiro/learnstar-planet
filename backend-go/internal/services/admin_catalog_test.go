// 管理端全校规则/商品 + 用户设置读写测试。
package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ---------------------------------- 管理端：积分规则 ----------------------------------

func TestListForSchool_SeedsDefaults(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	rules := services.NewRules(db)

	views, err := rules.ListForSchool(school.ID)
	require.NoError(t, err)
	require.Len(t, views, len(services.DefaultRules), "首次访问应补齐全部默认规则")
	assert.Equal(t, 27, countPositive(views))
	assert.Equal(t, 16, len(views)-countPositive(views))

	for _, v := range views {
		assert.Equal(t, "school", v.Scope)
		assert.Nil(t, v.ClassName, "学校级规则不带班级名")
		require.NotNil(t, v.SchoolID)
		assert.Equal(t, school.ID, *v.SchoolID)
		assert.Nil(t, v.ClassID)
	}

	var reloaded models.School
	require.NoError(t, db.First(&reloaded, school.ID).Error)
	assert.True(t, reloaded.ScoreRulesSeeded, "播种后应写标记")

	// 幂等：再次访问不重复播种。
	again, err := rules.ListForSchool(school.ID)
	require.NoError(t, err)
	assert.Len(t, again, len(views))
}

func TestListForSchool_IncludesClassRules(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	rules := services.NewRules(db)

	// 教师创建的班级级规则（class_id 非空）。
	created, err := rules.Create(&class.ID, school.ID, "班级专属加分", 7, "custom", true)
	require.NoError(t, err)

	views, err := rules.ListForSchool(school.ID)
	require.NoError(t, err)

	var found bool
	for _, v := range views {
		if v.ID == created.ID {
			found = true
			assert.Equal(t, "class", v.Scope)
			require.NotNil(t, v.ClassName)
			assert.Equal(t, class.Name, *v.ClassName)
		}
	}
	assert.True(t, found, "教师创建的班级级规则应出现在管理员列表中")
}

func TestCreateForSchool(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	rules := services.NewRules(db)

	rule, err := rules.CreateForSchool(school.ID, "全校加分", 5, "", true, false)
	require.NoError(t, err)
	assert.Nil(t, rule.ClassID, "学校级规则 class_id 必须为 null")
	require.NotNil(t, rule.SchoolID)
	assert.Equal(t, school.ID, *rule.SchoolID)
	assert.Equal(t, "custom", rule.Category, "未传分类应回退 custom")

	// GORM 会跳过带 default 的零值字段，此处确认 is_active=false 真的落库。
	var reloaded models.ScoreRule
	require.NoError(t, db.First(&reloaded, rule.ID).Error)
	assert.False(t, reloaded.IsActive)
	assert.True(t, reloaded.IsPositive)
}

func TestFindSchoolLevel_RejectsClassAndCrossSchool(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	other := seedOtherSchool(t, db)
	class := seedClass(t, db, school.ID)
	rules := services.NewRules(db)

	schoolRule, err := rules.CreateForSchool(school.ID, "校级规则", 3, "classroom", true, true)
	require.NoError(t, err)

	classRule, err := rules.Create(&class.ID, school.ID, "班级规则", 2, "custom", true)
	require.NoError(t, err)

	otherRule, err := rules.CreateForSchool(other.ID, "他校规则", 4, "custom", true, true)
	require.NoError(t, err)

	got, err := rules.FindSchoolLevel(school.ID, schoolRule.ID)
	require.NoError(t, err)
	assert.Equal(t, schoolRule.ID, got.ID)

	_, err = rules.FindSchoolLevel(school.ID, classRule.ID)
	assertAppStatus(t, err, 404, "班级级规则不允许在管理员端点单条操作")

	_, err = rules.FindSchoolLevel(school.ID, otherRule.ID)
	assertAppStatus(t, err, 404, "跨校规则不可见")
}

func TestAdminRuleUpdateAndDelete(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	rules := services.NewRules(db)

	rule, err := rules.CreateForSchool(school.ID, "待改规则", 3, "custom", true, true)
	require.NoError(t, err)

	updated, err := rules.Update(rule, map[string]any{"name": "已改规则", "amount": -6, "is_positive": false, "is_active": false})
	require.NoError(t, err)
	assert.Equal(t, "已改规则", updated.Name)
	assert.Equal(t, -6, updated.Amount)
	assert.False(t, updated.IsPositive)
	assert.False(t, updated.IsActive)

	require.NoError(t, rules.Delete(updated))
	var count int64
	require.NoError(t, db.Model(&models.ScoreRule{}).Where("id = ?", rule.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

// ---------------------------------- 管理端：商品 ----------------------------------

func TestItemsForSchool_ScopesAndImages(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	shop := newShopService(db, school.ID)

	schoolItem, err := shop.CreateSchoolItem(school.ID, services.SchoolItemInput{Name: "学校级奖品", CostScore: 50})
	require.NoError(t, err)

	classItem := models.ShopItem{
		ClassID: &class.ID, SchoolID: school.ID, Name: "班级级奖品",
		Category: "physical", CostScore: 10, CurrencyType: "score", Stock: 3, IsActive: true,
	}
	require.NoError(t, db.Create(&classItem).Error)

	views, err := shop.ItemsForSchool(school.ID, "")
	require.NoError(t, err)
	require.Len(t, views, 2)

	byID := map[uint]services.AdminItemView{}
	for _, v := range views {
		byID[v.ID] = v
	}

	v := byID[schoolItem.ID]
	assert.Equal(t, "school", v.Scope, "学校级商品同步所有班级")
	assert.Nil(t, v.ClassName)
	assert.True(t, v.IsActive)

	cv := byID[classItem.ID]
	assert.Equal(t, "class", cv.Scope)
	require.NotNil(t, cv.ClassName)
	assert.Equal(t, class.Name, *cv.ClassName)
}

func TestItemsForSchool_ClassFallbackAndCurrencyFilter(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	shop := newShopService(db, school.ID)

	missingClassID := uint(999)
	orphan := models.ShopItem{
		ClassID: &missingClassID, SchoolID: school.ID, Name: "孤儿商品",
		Category: "physical", CostScore: 5, CurrencyType: "science", IsActive: true,
	}
	require.NoError(t, db.Create(&orphan).Error)
	_, err := shop.CreateSchoolItem(school.ID, services.SchoolItemInput{Name: "积分商品", CostScore: 5})
	require.NoError(t, err)

	all, err := shop.ItemsForSchool(school.ID, "")
	require.NoError(t, err)
	require.Len(t, all, 2)
	for _, v := range all {
		if v.ID == orphan.ID {
			require.NotNil(t, v.ClassName)
			assert.Equal(t, "#999", *v.ClassName, "班级已不存在时用 #id 兜底")
		}
	}

	filtered, err := shop.ItemsForSchool(school.ID, "science")
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	assert.Equal(t, orphan.ID, filtered[0].ID)

	otherSchool, err := shop.ItemsForSchool(school.ID+1, "")
	require.NoError(t, err)
	assert.Empty(t, otherSchool, "其他学校不可见")
}

func TestCreateSchoolItem_DefaultsAndExplicitInactive(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	shop := newShopService(db, school.ID)

	def, err := shop.CreateSchoolItem(school.ID, services.SchoolItemInput{Name: "默认商品", CostScore: 20})
	require.NoError(t, err)
	assert.Nil(t, def.ClassID)
	assert.Equal(t, "physical", def.Category)
	assert.Equal(t, "score", def.CurrencyType)
	assert.True(t, def.IsActive, "未传 is_active 默认上架")

	inactive := false
	off, err := shop.CreateSchoolItem(school.ID, services.SchoolItemInput{Name: "下架商品", CostScore: 30, IsActive: &inactive})
	require.NoError(t, err)
	var reloaded models.ShopItem
	require.NoError(t, db.First(&reloaded, off.ID).Error)
	assert.False(t, reloaded.IsActive, "显式 false 必须落库（绕过 GORM 零值跳过）")
}

func TestFindSchoolItem_RejectsClassLevelAndCrossSchool(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	other := seedOtherSchool(t, db)
	class := seedClass(t, db, school.ID)
	shop := newShopService(db, school.ID)

	schoolItem, err := shop.CreateSchoolItem(school.ID, services.SchoolItemInput{Name: "校级商品", CostScore: 20})
	require.NoError(t, err)

	classItem := models.ShopItem{
		ClassID: &class.ID, SchoolID: school.ID, Name: "班级商品",
		Category: "physical", CostScore: 10, CurrencyType: "score", IsActive: true,
	}
	require.NoError(t, db.Create(&classItem).Error)

	otherItem, err := shop.CreateSchoolItem(other.ID, services.SchoolItemInput{Name: "他校商品", CostScore: 20})
	require.NoError(t, err)

	got, err := shop.FindSchoolItem(school.ID, schoolItem.ID)
	require.NoError(t, err)
	assert.Equal(t, schoolItem.ID, got.ID)

	_, err = shop.FindSchoolItem(school.ID, classItem.ID)
	assertAppStatus(t, err, 404, "班级级商品不允许在管理员端点单条操作")

	_, err = shop.FindSchoolItem(school.ID, otherItem.ID)
	assertAppStatus(t, err, 404, "跨校商品不可见")
}

func TestAdminItemUpdateAndDelete(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	shop := newShopService(db, school.ID)

	item, err := shop.CreateSchoolItem(school.ID, services.SchoolItemInput{Name: "待改商品", CostScore: 20})
	require.NoError(t, err)

	updated, err := shop.UpdateItem(item, map[string]any{"name": "已改商品", "cost_score": 35, "stock": 8, "is_active": false})
	require.NoError(t, err)
	assert.Equal(t, "已改商品", updated.Name)
	assert.Equal(t, 35, updated.CostScore)
	assert.Equal(t, 8, updated.Stock)
	assert.False(t, updated.IsActive)

	require.NoError(t, shop.DeleteItem(updated))
	var count int64
	require.NoError(t, db.Model(&models.ShopItem{}).Where("id = ?", item.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

// ---------------------------------- 辅助 ----------------------------------

// newShopService 构造商城服务（内部依赖评分与币种服务，测试中只需它们可构造）。
func newShopService(db *gorm.DB, _ uint) *services.ShopService {
	scope := services.NewScope(db)
	scores := services.NewScoreService(db)
	currency := services.NewCurrencyService(db)
	return services.NewShopService(db, scope, scores, currency)
}

func countPositive(views []services.AdminRuleView) int {
	n := 0
	for _, v := range views {
		if v.IsPositive {
			n++
		}
	}
	return n
}

// assertAppStatus 断言错误是带指定状态码的业务错误。
func assertAppStatus(t *testing.T, err error, status int, msg string) {
	t.Helper()
	require.Error(t, err, msg)
	ae, ok := services.AsAppError(err)
	require.True(t, ok, "应为业务错误：%v", err)
	assert.Equal(t, status, ae.Status, msg)
}

// seedOtherSchool 创建第二所学校：seedSchool 的 code 固定，同一测试内重复调用会撞唯一索引。
func seedOtherSchool(t *testing.T, db *gorm.DB) models.School {
	t.Helper()
	school := models.School{Name: "他校", Code: "other-school", Status: "active"}
	if err := db.Create(&school).Error; err != nil {
		t.Fatalf("create other school: %v", err)
	}
	return school
}
