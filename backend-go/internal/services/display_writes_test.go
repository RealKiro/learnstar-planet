// 教室端（班级码）写操作测试：商品列表、快捷兑换、学生间转赠、整班切换系列、学生换宠。
package services_test

import (
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newDisplayShopService 造一个完整的商城服务（依赖与 handlers.New 一致）。
func newDisplayShopService(db *gorm.DB) *services.ShopService {
	return services.NewShopService(db, services.NewScope(db), services.NewScoreService(db), services.NewCurrencyService(db))
}

// seedClassShopItem 造班级级商品；active=false 时需创建后显式回写
// （is_active 带 default:true，GORM 会跳过 false 零值）。
func seedClassShopItem(t *testing.T, db *gorm.DB, classID, schoolID uint, name string, cost int, active bool) models.ShopItem {
	t.Helper()
	item := models.ShopItem{
		ClassID: &classID, SchoolID: schoolID, Name: name, Description: name + "（描述）",
		Category: "stationery", CostScore: cost, CurrencyType: "score", Stock: 5, IsActive: true,
	}
	require.NoError(t, db.Create(&item).Error)
	if !active {
		require.NoError(t, db.Model(&models.ShopItem{}).Where("id = ?", item.ID).
			Update("is_active", false).Error)
		item.IsActive = false
	}
	return item
}

// shop-items 只返回「本班级级 + 在售」；redeem 的余额/存在性守卫与结算字段。
func TestDisplayShopItemsAndRedeem(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	shop := newDisplayShopService(db)

	poor := models.Student{ClassID: f.Class.ID, Name: "小明", StudentNo: "1", Status: "active"}
	rich := models.Student{ClassID: f.Class.ID, Name: "有钱", StudentNo: "2", TotalScore: 100, Status: "active"}
	require.NoError(t, db.Create(&poor).Error)
	require.NoError(t, db.Create(&rich).Error)

	onSale := seedClassShopItem(t, db, f.Class.ID, f.School.ID, "铅笔", 30, true)
	seedClassShopItem(t, db, f.Class.ID, f.School.ID, "下架货", 10, false)
	// 学校级商品（class_id = null）与原实现的教室端列表口径无关，不应出现。
	require.NoError(t, db.Create(&models.ShopItem{
		SchoolID: f.School.ID, Name: "校级商品", Category: "stationery",
		CostScore: 20, CurrencyType: "score", Stock: 5, IsActive: true,
	}).Error)
	// 他班商品与他班学生（跨班兑换应 404）。
	otherClass := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	seedClassShopItem(t, db, otherClass.ID, f.School.ID, "他班商品", 10, true)
	outsider := models.Student{ClassID: otherClass.ID, Name: "外班同学", StudentNo: "1", TotalScore: 500, Status: "active"}
	require.NoError(t, db.Create(&outsider).Error)

	items, err := shop.DisplayItems(f.Class.ID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, onSale.ID, items[0].ID)
	assert.Equal(t, "铅笔", items[0].Name)
	assert.Equal(t, "铅笔（描述）", items[0].Description)
	assert.Equal(t, 30, items[0].CostScore)
	assert.Equal(t, 5, items[0].Stock)
	assert.Equal(t, "stationery", items[0].Category)

	// 积分不足 → 400，且不留兑换记录。
	_, err = shop.DisplayRedeem(f.Class.ID, poor.ID, onSale.ID)
	assert.Equal(t, "积分不足", appErrorOf(t, err, 400))
	var count int64
	require.NoError(t, db.Model(&models.ShopRedemption{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)

	// 学生或商品不存在（含跨班学生 / 已下架商品）→ 404。
	_, err = shop.DisplayRedeem(f.Class.ID, 999999, onSale.ID)
	assert.Equal(t, "学生或商品不存在", appErrorOf(t, err, 404))
	_, err = shop.DisplayRedeem(f.Class.ID, outsider.ID, onSale.ID)
	assert.Equal(t, "学生或商品不存在", appErrorOf(t, err, 404))

	// 成功：扣分 + 生成 approved 兑换记录 + 库存不变（同 Laravel quickRedeem）。
	result, err := shop.DisplayRedeem(f.Class.ID, rich.ID, onSale.ID)
	require.NoError(t, err)
	assert.Equal(t, "有钱", result.StudentName)
	assert.Equal(t, "铅笔", result.ItemName)
	assert.Equal(t, 30, result.Cost)
	assert.Equal(t, 70, result.TotalScore)

	var stored models.Student
	require.NoError(t, db.First(&stored, rich.ID).Error)
	assert.Equal(t, 70, stored.TotalScore)

	var redemption models.ShopRedemption
	require.NoError(t, db.Where("student_id = ?", rich.ID).First(&redemption).Error)
	assert.Equal(t, "approved", redemption.Status)
	assert.Equal(t, onSale.ID, redemption.ShopItemID)
	assert.Equal(t, f.Class.ID, redemption.ClassID)
	assert.Equal(t, 30, redemption.Cost)
	require.NotNil(t, redemption.ApprovedAt)
	require.NotNil(t, redemption.ApprovedBy)

	var itemAfter models.ShopItem
	require.NoError(t, db.First(&itemAfter, onSale.ID).Error)
	assert.Equal(t, 5, itemAfter.Stock, "Laravel quickRedeem 不校验也不扣库存")

	// 扣分记录（ScoreService 口径）+ score_update 事件（is_spend）。
	var score models.Score
	require.NoError(t, db.Where("student_id = ?", rich.ID).Order("id DESC").First(&score).Error)
	assert.Equal(t, -30, score.Amount)
	assert.Equal(t, "兑换消耗：兑换：铅笔", score.Reason)

	events, err := services.NewDisplayEvents(db).Consume(f.Class.ID, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)
	last := events[len(events)-1]
	assert.Equal(t, services.DisplayEventScoreUpdate, last.Type)
	assert.Contains(t, string(last.Data), `"is_spend":true`)
}

// transfer：成功（两生分数变化）、余额不足、给自己转、跨班学生拒绝、金额边界。
func TestDisplayTransfer(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	from := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 50, Status: "active"}
	to := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", TotalScore: 10, Status: "active"}
	require.NoError(t, db.Create(&from).Error)
	require.NoError(t, db.Create(&to).Error)

	otherClass := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	outsider := models.Student{ClassID: otherClass.ID, Name: "丙", StudentNo: "1", TotalScore: 100, Status: "active"}
	require.NoError(t, db.Create(&outsider).Error)

	// 给自己转 → 422（先于存在性判断）。
	_, err := f.Svc.QuickTransfer(f.Class.ID, from.ID, from.ID, 5)
	assert.Equal(t, "不能转赠给自己", appErrorOf(t, err, 422))

	// 金额越界 → 422。
	_, err = f.Svc.QuickTransfer(f.Class.ID, from.ID, to.ID, 0)
	assert.Equal(t, "转赠积分至少 1 分", appErrorOf(t, err, 422))
	_, err = f.Svc.QuickTransfer(f.Class.ID, from.ID, to.ID, 101)
	assert.Equal(t, "单次转赠不能超过 100 分", appErrorOf(t, err, 422))

	// 跨班学生 → 404。
	_, err = f.Svc.QuickTransfer(f.Class.ID, from.ID, outsider.ID, 5)
	assert.Equal(t, "学生不存在", appErrorOf(t, err, 404))

	// 余额不足 → 400。
	_, err = f.Svc.QuickTransfer(f.Class.ID, to.ID, from.ID, 40)
	assert.Equal(t, "积分不足", appErrorOf(t, err, 400))

	// 成功：转出 -30、转入 +30。
	result, err := f.Svc.QuickTransfer(f.Class.ID, from.ID, to.ID, 30)
	require.NoError(t, err)
	assert.Equal(t, "甲", result.FromName)
	assert.Equal(t, "乙", result.ToName)
	assert.Equal(t, 30, result.Amount)

	var a, b models.Student
	require.NoError(t, db.First(&a, from.ID).Error)
	require.NoError(t, db.First(&b, to.ID).Error)
	assert.Equal(t, 20, a.TotalScore)
	assert.Equal(t, 40, b.TotalScore)

	// 两条积分记录（同 Laravel 的两次 giveScore）。
	var scores []models.Score
	require.NoError(t, db.Where("class_id = ?", f.Class.ID).Order("id ASC").Find(&scores).Error)
	require.Len(t, scores, 2)
	assert.Equal(t, -30, scores[0].Amount)
	assert.Equal(t, "转赠给 乙", scores[0].Reason)
	assert.Equal(t, 30, scores[1].Amount)
	assert.Equal(t, "来自 甲 的转赠", scores[1].Reason)

	// 两条 score_update 事件。
	events, err := services.NewDisplayEvents(db).Consume(f.Class.ID, nil)
	require.NoError(t, err)
	assert.Len(t, events, 2)
}

// switch-series（教室端）：非法系列 422、积分不足 400、成功写设置 + 每人扣 20 + 发免费自选。
func TestClassroomSwitchSeries(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	petSeries := services.NewPetSeriesService(db, services.NewScope(db))

	rich := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 40, Status: "active"}
	require.NoError(t, db.Create(&rich).Error)
	// 已有宠物：整班切换不重抽（原实现同样不重抽）。
	require.NoError(t, db.Create(&models.Pet{
		StudentID: rich.ID, ClassID: f.Class.ID, Name: "甲的萌宠", Species: "zhulong", Level: 3, Mood: 70,
	}).Error)

	// 非法系列 → 422。
	_, err := petSeries.SwitchSeriesForClassroom(f.Class.ID, "nope")
	assert.Equal(t, "无效的系列ID", appErrorOf(t, err, 422))

	poor := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", TotalScore: 19, Status: "active"}
	require.NoError(t, db.Create(&poor).Error)

	// 有人积分不足 → 400（整批拒绝），且不写班级设置。
	_, err = petSeries.SwitchSeriesForClassroom(f.Class.ID, "pokemon")
	assert.Equal(t, "积分不足：乙 每人需要 20 积分", appErrorOf(t, err, 400))
	var untouched models.ClassRoom
	require.NoError(t, db.First(&untouched, f.Class.ID).Error)
	assert.Equal(t, "", untouched.SettingString("pet_series", ""))

	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", poor.ID).Update("total_score", 20).Error)

	result, err := petSeries.SwitchSeriesForClassroom(f.Class.ID, "pokemon")
	require.NoError(t, err)
	assert.Equal(t, "已切换至「pokemon」系列：全班 2 人各扣 20 积分，并各获一次免费自选该系列宠物的机会", result.Message)
	assert.Equal(t, "pokemon", result.Data.SeriesID)
	assert.Equal(t, 20, result.Data.CostPerStudent)
	assert.Equal(t, 2, result.Data.AffectedStudents)
	assert.True(t, result.Data.FreePickGranted)

	// 班级设置写入 pet_series。
	var class models.ClassRoom
	require.NoError(t, db.First(&class, f.Class.ID).Error)
	assert.Equal(t, "pokemon", class.SettingString("pet_series", ""))

	// 每人各扣 20。
	var afterRich, afterPoor models.Student
	require.NoError(t, db.First(&afterRich, rich.ID).Error)
	require.NoError(t, db.First(&afterPoor, poor.ID).Error)
	assert.Equal(t, 20, afterRich.TotalScore)
	assert.Equal(t, 0, afterPoor.TotalScore)

	// 免费自选机会：每人一条、3 天后过期。
	var picks []models.PetFreePick
	require.NoError(t, db.Order("student_id ASC").Find(&picks).Error)
	require.Len(t, picks, 2)
	assert.WithinDuration(t, time.Now().Add(3*24*time.Hour), picks[0].ExpiresAt, 2*time.Minute)

	// 已有宠物不被重抽，也没有积分/审计记录（Laravel 直接改 total_score）。
	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", rich.ID).First(&pet).Error)
	assert.Equal(t, "zhulong", pet.Species)
	assert.Equal(t, 3, pet.Level)
	var scoreCount int64
	require.NoError(t, db.Model(&models.Score{}).Count(&scoreCount).Error)
	assert.Equal(t, int64(0), scoreCount)

	// 班级不存在 → 404；没有活跃学生 → 400。
	_, err = petSeries.SwitchSeriesForClassroom(999999, "pokemon")
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))

	empty := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（3）班", Status: "active"}
	require.NoError(t, db.Create(&empty).Error)
	_, err = petSeries.SwitchSeriesForClassroom(empty.ID, "pokemon")
	assert.Equal(t, "班级没有活跃学生", appErrorOf(t, err, 400))
}

// switch-series（教师端）：没有可管理班级 400、非法系列 422（附可选值）、成功不扣分。
func TestTeacherSwitchSeries(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	petSeries := services.NewPetSeriesService(db, services.NewScope(db))

	student := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 40, Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	lonely := seedTeacher(t, db, f.School.ID, "teacher-no-class")
	_, err := petSeries.SwitchSeries(&lonely, "pokemon")
	assert.Equal(t, "没有可管理的班级", appErrorOf(t, err, 400))

	_, err = petSeries.SwitchSeries(&f.Teacher, "nope")
	assert.Equal(t, "无效的系列ID，可选值：myth, pokemon, national, digimon, magic, prehistoric, constellation, festival, qixia, dongfang",
		appErrorOf(t, err, 422))

	result, err := petSeries.SwitchSeries(&f.Teacher, "myth")
	require.NoError(t, err)
	assert.Equal(t, "已切换系列为「myth」，全班 1 人各获一次免费自选该系列宠物的机会", result.Message)
	assert.Equal(t, "myth", result.Data.SeriesID)
	assert.Equal(t, f.Class.ID, result.Data.ClassID)
	assert.True(t, result.Data.FreePickGranted)
	assert.Equal(t, 1, result.Data.GrantedStudents)

	// 不扣积分。
	var after models.Student
	require.NoError(t, db.First(&after, student.ID).Error)
	assert.Equal(t, 40, after.TotalScore)

	var class models.ClassRoom
	require.NoError(t, db.First(&class, f.Class.ID).Error)
	assert.Equal(t, "myth", class.SettingString("pet_series", ""))
}

// switch-pet（教室端）：新宠创建、同物种 422、跨班 404、系列限制 422、按等级扣分/不足、免费自选。
func TestClassroomSwitchPet(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	pets := services.NewPetService(db, services.NewScope(db))

	student := models.Student{ClassID: f.Class.ID, Name: "小明", StudentNo: "1", TotalScore: 100, Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	// 无宠物：直接创建（cost = 0，无 free_pick_used 字段）。
	created, err := pets.SwitchForClassroom(f.Class.ID, student.ID, "zhulong")
	require.NoError(t, err)
	assert.Equal(t, "🎉 新宠物已诞生！", created.Message)
	assert.Equal(t, "小明的伙伴", created.Data.PetName)
	assert.Equal(t, "zhulong", created.Data.PetSpecies)
	assert.Equal(t, "🥚", created.Data.PetEmoji)
	assert.Equal(t, 100, created.Data.TotalScore)
	assert.Equal(t, 0, created.Data.Cost)
	assert.Nil(t, created.Data.FreePickUsed)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, f.Class.ID, pet.ClassID)
	assert.Equal(t, 1, pet.Level)
	assert.Equal(t, 0, pet.Experience)
	assert.Equal(t, 80, pet.Mood)

	// 已是该物种 → 422，不扣分。
	_, err = pets.SwitchForClassroom(f.Class.ID, student.ID, "zhulong")
	assert.Equal(t, "当前已经是这只宠物啦", appErrorOf(t, err, 422))

	// 班级不存在 → 404；非本班学生 → 404。
	_, err = pets.SwitchForClassroom(999999, student.ID, "qilin")
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))

	otherClass := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	outsider := models.Student{ClassID: otherClass.ID, Name: "外人", TotalScore: 500, Status: "active"}
	require.NoError(t, db.Create(&outsider).Error)
	_, err = pets.SwitchForClassroom(f.Class.ID, outsider.ID, "qilin")
	assert.Equal(t, "学生不存在", appErrorOf(t, err, 404))

	// 班级设了系列 → 只能在本系列内换（跨系列 422，文案带中文系列名）。
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", f.Class.ID).
		Update("settings", `{"pet_series":"pokemon"}`).Error)
	_, err = pets.SwitchForClassroom(f.Class.ID, student.ID, "qilin")
	assert.Equal(t, "只能领养当前类别「宝可梦」的宠物，不能跨类别领养", appErrorOf(t, err, 422))

	// 换宠扣分：Lv.1 → 5 分（SwitchCost = 5 × 等级），等级/经验/心情保留。
	switched, err := pets.SwitchForClassroom(f.Class.ID, student.ID, "pikachu")
	require.NoError(t, err)
	assert.Equal(t, "✅ 已更换为「小明的伙伴」，扣除 5 积分", switched.Message)
	assert.Equal(t, "pikachu", switched.Data.PetSpecies)
	assert.Equal(t, 5, switched.Data.Cost)
	assert.Equal(t, 95, switched.Data.TotalScore)
	require.NotNil(t, switched.Data.FreePickUsed)
	assert.False(t, *switched.Data.FreePickUsed)

	var after models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&after).Error)
	assert.Equal(t, "pikachu", after.Species)
	assert.Equal(t, 1, after.Level)

	// 不写积分记录 / 不发事件（Laravel classroomSwitchPet 直接改 total_score）。
	var scoreCount int64
	require.NoError(t, db.Model(&models.Score{}).Count(&scoreCount).Error)
	assert.Equal(t, int64(0), scoreCount)
	events, err := services.NewDisplayEvents(db).Consume(f.Class.ID, nil)
	require.NoError(t, err)
	assert.Empty(t, events)

	// 积分不足 → 400（Lv.12 → 60 分）。
	require.NoError(t, db.Model(&models.Pet{}).Where("id = ?", after.ID).Update("level", 12).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).Update("total_score", 10).Error)
	_, err = pets.SwitchForClassroom(f.Class.ID, student.ID, "riolu")
	assert.Equal(t, "积分不足，更换宠物需 60 积分", appErrorOf(t, err, 400))

	// 持有免费自选机会 → 本次免费（用掉即失效）。
	require.NoError(t, models.GrantPetFreePick(db, student.ID, time.Now()))
	free, err := pets.SwitchForClassroom(f.Class.ID, student.ID, "riolu")
	require.NoError(t, err)
	assert.Equal(t, "✅ 已使用整班切换的免费自选机会！", free.Message)
	assert.Equal(t, 0, free.Data.Cost)
	assert.Equal(t, 10, free.Data.TotalScore)
	require.NotNil(t, free.Data.FreePickUsed)
	assert.True(t, *free.Data.FreePickUsed)

	has, err := models.HasPetFreePick(db, student.ID, time.Now())
	require.NoError(t, err)
	assert.False(t, has, "机会使用后即失效")

	// 机会已消费 → 再换需扣分（此时 10 分 < 60 分 → 400）。
	_, err = pets.SwitchForClassroom(f.Class.ID, student.ID, "pikachu")
	assert.Equal(t, "积分不足，更换宠物需 60 积分", appErrorOf(t, err, 400))
}
