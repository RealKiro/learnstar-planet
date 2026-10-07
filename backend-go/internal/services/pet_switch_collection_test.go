// 宠物切换的图鉴归档 / 进度恢复单测（教师端 POST /teacher/pets/:student_id/switch 与
// 教室端 POST /display/pets/switch 两条语义**不同**的路径）。
//
// 权威来源：Laravel App\Services\PetService::switchPet（第 155-276 行）与
// App\Http\Controllers\Api\DisplayController::classroomSwitchPet（第 1182-1290 行）。
// 全部使用内存 SQLite，**不访问外网**。
package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// seedPetSwitchFixture 建教师 + 班级 + 学生（初始 100 分，便于验证扣分）。
func seedPetSwitchFixture(t *testing.T, db *gorm.DB) (models.User, models.ClassRoom, models.Student) {
	t.Helper()
	school := seedSchool(t, db)
	teacher, class := seedTeacherWithClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("total_score", 100).Error)
	student.TotalScore = 100
	return teacher, class, student
}

// petCollectionOf 读取某学生某物种的图鉴行；不存在返回 nil。
func petCollectionOf(t *testing.T, db *gorm.DB, studentID uint, species string) *models.PetCollection {
	t.Helper()
	var row models.PetCollection
	err := db.Where("student_id = ? AND species = ?", studentID, species).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	require.NoError(t, err)
	return &row
}

// countPetCollections 某学生的图鉴行数。
func countPetCollections(t *testing.T, db *gorm.DB, studentID uint) int64 {
	t.Helper()
	var count int64
	require.NoError(t, db.Model(&models.PetCollection{}).
		Where("student_id = ?", studentID).Count(&count).Error)
	return count
}

// ============================================================
// A1 · 教师端 POST /teacher/pets/:student_id/switch
// ============================================================

// TestTeacherSwitchPetNoPetCreatesPetAndCollection 无宠物 → 建宠物 + 图鉴行（is_active=true）。
func TestTeacherSwitchPetNoPetCreatesPetAndCollection(t *testing.T) {
	db := setupDB(t)
	teacher, class, student := seedPetSwitchFixture(t, db)

	svc := services.NewPetService(db, services.NewScope(db))
	res, err := svc.Switch(&teacher, student.ID, "zhulong", "小火龙")
	require.NoError(t, err)

	// 文案与字段逐字同 Laravel 第 267-275 行（无 cost / free_pick_used）。
	assert.Equal(t, "已为您分配宠物「小火龙」", res.Message)
	assert.Equal(t, "小火龙", res.Data.PetName)
	assert.Equal(t, "zhulong", res.Data.PetSpecies)
	assert.Equal(t, 1, res.Data.Level)
	assert.Equal(t, 0, res.Data.Experience)
	assert.Nil(t, res.Data.Cost)
	assert.Nil(t, res.Data.FreePickUsed)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, class.ID, pet.ClassID)
	assert.Equal(t, "小火龙", pet.Name)
	assert.Equal(t, 1, pet.Level)
	assert.Equal(t, 0, pet.Experience)
	assert.Equal(t, 80, pet.Mood)
	assert.Nil(t, pet.LastFedAt)
	// 首次分配不是「切换」，Laravel 不写 last_switched_at。
	assert.Nil(t, pet.LastSwitchedAt)

	row := petCollectionOf(t, db, student.ID, "zhulong")
	require.NotNil(t, row, "无宠物分支必须写入图鉴")
	assert.Equal(t, 1, row.Level)
	assert.Equal(t, 0, row.Experience)
	assert.Equal(t, 80, row.Mood)
	assert.True(t, row.IsActive)

	var after models.Student
	require.NoError(t, db.First(&after, student.ID).Error)
	assert.Equal(t, 100, after.TotalScore, "首次分配不扣积分")
}

// TestTeacherSwitchPetRestoresTargetProgress 目标物种有图鉴进度 → 恢复 level/exp/mood；旧物种归档为 is_active=false。
func TestTeacherSwitchPetRestoresTargetProgress(t *testing.T) {
	db := setupDB(t)
	teacher, _, student := seedPetSwitchFixture(t, db)

	// 当前宠物：火狐 Lv.4 / exp 7 / mood 55；目标物种烛龙在图鉴里已有进度 Lv.9 / exp 12 / mood 30。
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: student.ClassID, Name: "火狐", Species: "huohu",
		Level: 4, Experience: 7, Mood: 55,
	}).Error)
	require.NoError(t, db.Create(&models.PetCollection{
		StudentID: student.ID, Species: "zhulong", Level: 9, Experience: 12, Mood: 30, IsActive: false,
	}).Error)

	svc := services.NewPetService(db, services.NewScope(db))
	res, err := svc.Switch(&teacher, student.ID, "zhulong", "烛龙归来")
	require.NoError(t, err)

	assert.Equal(t, "宠物已更换为「烛龙归来」（扣除 20 积分）", res.Message) // SwitchCost(4) = 20
	assert.Equal(t, 9, res.Data.Level, "恢复图鉴里的等级")
	assert.Equal(t, 12, res.Data.Experience)
	require.NotNil(t, res.Data.Cost)
	assert.Equal(t, 20, *res.Data.Cost)
	require.NotNil(t, res.Data.FreePickUsed)
	assert.False(t, *res.Data.FreePickUsed)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, "zhulong", pet.Species)
	assert.Equal(t, "烛龙归来", pet.Name)
	assert.Equal(t, 9, pet.Level)
	assert.Equal(t, 12, pet.Experience)
	assert.Equal(t, 30, pet.Mood)
	require.NotNil(t, pet.LastSwitchedAt, "切换应写 last_switched_at")

	// 旧物种归档（进度取切换前的 pet 值，is_active=false）。
	archived := petCollectionOf(t, db, student.ID, "huohu")
	require.NotNil(t, archived)
	assert.Equal(t, 4, archived.Level)
	assert.Equal(t, 7, archived.Experience)
	assert.Equal(t, 55, archived.Mood)
	assert.False(t, archived.IsActive)

	// 目标物种标为激活，且进度为恢复后的值。
	target := petCollectionOf(t, db, student.ID, "zhulong")
	require.NotNil(t, target)
	assert.Equal(t, 9, target.Level)
	assert.Equal(t, 12, target.Experience)
	assert.True(t, target.IsActive)
	assert.Equal(t, int64(2), countPetCollections(t, db, student.ID))

	var after models.Student
	require.NoError(t, db.First(&after, student.ID).Error)
	assert.Equal(t, 80, after.TotalScore)
}

// TestTeacherSwitchPetInitialisesUnknownSpecies 目标物种无图鉴 → 初始化为 1/0/80。
func TestTeacherSwitchPetInitialisesUnknownSpecies(t *testing.T) {
	db := setupDB(t)
	teacher, _, student := seedPetSwitchFixture(t, db)

	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: student.ClassID, Name: "烛龙", Species: "zhulong",
		Level: 5, Experience: 3, Mood: 90,
	}).Error)

	svc := services.NewPetService(db, services.NewScope(db))
	res, err := svc.Switch(&teacher, student.ID, "qilin", "麒麟")
	require.NoError(t, err)

	assert.Equal(t, 1, res.Data.Level)
	assert.Equal(t, 0, res.Data.Experience)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, 1, pet.Level)
	assert.Equal(t, 0, pet.Experience)
	assert.Equal(t, 80, pet.Mood)

	newRow := petCollectionOf(t, db, student.ID, "qilin")
	require.NotNil(t, newRow)
	assert.Equal(t, 1, newRow.Level)
	assert.Equal(t, 0, newRow.Experience)
	assert.Equal(t, 80, newRow.Mood)
	assert.True(t, newRow.IsActive)
}

// TestTeacherSwitchPetSameSpecies 同物种 → 422，且不写图鉴、不扣分。
func TestTeacherSwitchPetSameSpecies(t *testing.T) {
	db := setupDB(t)
	teacher, _, student := seedPetSwitchFixture(t, db)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: student.ClassID, Name: "烛龙", Species: "zhulong", Level: 3, Mood: 80,
	}).Error)

	svc := services.NewPetService(db, services.NewScope(db))
	_, err := svc.Switch(&teacher, student.ID, "zhulong", "别的名字")
	assert.Equal(t, "当前已经是这只宠物啦", appErrorOf(t, err, 422))

	assert.Equal(t, int64(0), countPetCollections(t, db, student.ID), "同物种守卫不应写图鉴")
	var after models.Student
	require.NoError(t, db.First(&after, student.ID).Error)
	assert.Equal(t, 100, after.TotalScore)
}

// TestTeacherSwitchPetCrossSeriesUsesRawSeriesID 教师端跨类别 422 的文案用**原始系列 id**
// （Laravel `PetService::switchPet` 第 174 行），与教室端用 seriesLabel 中文标签不同。
func TestTeacherSwitchPetCrossSeriesUsesRawSeriesID(t *testing.T) {
	db := setupDB(t)
	teacher, class, student := seedPetSwitchFixture(t, db)
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Update("settings", `{"pet_series":"pokemon"}`).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: class.ID, Name: "皮卡丘", Species: "pikachu", Level: 2, Mood: 80,
	}).Error)

	svc := services.NewPetService(db, services.NewScope(db))

	_, err := svc.Switch(&teacher, student.ID, "qilin", "麒麟")
	assert.Equal(t, "只能领养当前类别「pokemon」的宠物，不能跨类别领养", appErrorOf(t, err, 422))

	// 同系列内可换（Lv.2 → 10 分）。
	res, err := svc.Switch(&teacher, student.ID, "riolu", "利欧路")
	require.NoError(t, err)
	assert.Equal(t, "riolu", res.Data.PetSpecies)
}

// TestTeacherSwitchPetUsesFreePick 免费自选 → 不扣分、free_pick_used=true、pet_free_picks 记录被删除。
func TestTeacherSwitchPetUsesFreePick(t *testing.T) {
	db := setupDB(t)
	teacher, _, student := seedPetSwitchFixture(t, db)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: student.ClassID, Name: "烛龙", Species: "zhulong", Level: 12, Mood: 80,
	}).Error)
	require.NoError(t, models.GrantPetFreePick(db, student.ID, time.Now()))

	svc := services.NewPetService(db, services.NewScope(db))
	res, err := svc.Switch(&teacher, student.ID, "qilin", "麒麟")
	require.NoError(t, err)

	assert.Equal(t, "✅ 已使用整班切换的免费自选机会！", res.Message)
	require.NotNil(t, res.Data.Cost)
	assert.Equal(t, 0, *res.Data.Cost, "免费自选不扣分")
	require.NotNil(t, res.Data.FreePickUsed)
	assert.True(t, *res.Data.FreePickUsed)

	var after models.Student
	require.NoError(t, db.First(&after, student.ID).Error)
	assert.Equal(t, 100, after.TotalScore)

	has, err := models.HasPetFreePick(db, student.ID, time.Now())
	require.NoError(t, err)
	assert.False(t, has, "免费自选机会用掉即失效（pet_free_picks 行被删除）")

	// 目标物种仍是初始形态（首次进入图鉴）。
	target := petCollectionOf(t, db, student.ID, "qilin")
	require.NotNil(t, target)
	assert.Equal(t, 1, target.Level)
	assert.True(t, target.IsActive)
}

// ============================================================
// A2 · 教室端 POST /display/pets/switch
// ============================================================

// TestClassroomSwitchPetKeepsLevelAndArchives 归档旧物种 + 目标激活 + 等级/经验保留（不恢复）。
func TestClassroomSwitchPetKeepsLevelAndArchives(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	svc := services.NewPetService(db, services.NewScope(db))

	student := models.Student{ClassID: f.Class.ID, Name: "小明", StudentNo: "1", TotalScore: 100, Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	// 图鉴里已有目标物种的旧进度（教室端**不**恢复它，只覆盖为当前值）。
	require.NoError(t, db.Create(&models.PetCollection{
		StudentID: student.ID, Species: "qilin", Level: 1, Experience: 0, Mood: 80, IsActive: false,
	}).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: f.Class.ID, Name: "小明的萌宠", Species: "zhulong",
		Level: 6, Experience: 9, Mood: 44,
	}).Error)

	res, err := svc.SwitchForClassroom(f.Class.ID, student.ID, "qilin")
	require.NoError(t, err)
	assert.Equal(t, "✅ 已更换为「小明的萌宠」，扣除 30 积分", res.Message) // SwitchCost(6) = 30
	require.NotNil(t, res.Data.FreePickUsed)
	assert.False(t, *res.Data.FreePickUsed)
	assert.Equal(t, 70, res.Data.TotalScore)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, "qilin", pet.Species)
	assert.Equal(t, 6, pet.Level, "教室端保留等级（不恢复图鉴里的 Lv.1）")
	assert.Equal(t, 9, pet.Experience)
	assert.Equal(t, 44, pet.Mood)

	archived := petCollectionOf(t, db, student.ID, "zhulong")
	require.NotNil(t, archived, "旧物种必须归档进图鉴")
	assert.Equal(t, 6, archived.Level)
	assert.Equal(t, 9, archived.Experience)
	assert.Equal(t, 44, archived.Mood)
	assert.False(t, archived.IsActive)

	active := petCollectionOf(t, db, student.ID, "qilin")
	require.NotNil(t, active)
	assert.Equal(t, 6, active.Level, "目标物种用当前 pet 的等级覆盖")
	assert.Equal(t, 9, active.Experience)
	assert.Equal(t, 44, active.Mood)
	assert.True(t, active.IsActive)
	assert.Equal(t, int64(2), countPetCollections(t, db, student.ID))
}

// TestClassroomSwitchPetCrossSeriesUsesLabel 跨类别 422 用**系列中文标签**（seriesLabel）。
func TestClassroomSwitchPetCrossSeriesUsesLabel(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	svc := services.NewPetService(db, services.NewScope(db))

	student := models.Student{ClassID: f.Class.ID, Name: "小明", TotalScore: 100, Status: "active"}
	require.NoError(t, db.Create(&student).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: f.Class.ID, Name: "小明的萌宠", Species: "zhulong", Level: 2, Mood: 80,
	}).Error)
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", f.Class.ID).
		Update("settings", `{"pet_series":"pokemon"}`).Error)

	_, err := svc.SwitchForClassroom(f.Class.ID, student.ID, "qilin")
	assert.Equal(t, "只能领养当前类别「宝可梦」的宠物，不能跨类别领养", appErrorOf(t, err, 422))
}

// TestClassroomSwitchPetNoPetDoesNotWriteCollection 无宠物分支：建 pets、**不**写图鉴（同 Laravel else 分支）。
func TestClassroomSwitchPetNoPetDoesNotWriteCollection(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	svc := services.NewPetService(db, services.NewScope(db))

	student := models.Student{ClassID: f.Class.ID, Name: "小明", TotalScore: 40, Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	res, err := svc.SwitchForClassroom(f.Class.ID, student.ID, "zhulong")
	require.NoError(t, err)
	assert.Equal(t, "🎉 新宠物已诞生！", res.Message)
	assert.Equal(t, "小明的伙伴", res.Data.PetName)
	assert.Equal(t, "zhulong", res.Data.PetSpecies)
	assert.Equal(t, 0, res.Data.Cost)
	assert.Equal(t, 40, res.Data.TotalScore)
	assert.Nil(t, res.Data.FreePickUsed, "无宠物分支无 free_pick_used 字段")

	assert.Equal(t, int64(0), countPetCollections(t, db, student.ID), "Laravel 的无宠物分支不写图鉴")
}

// TestClassroomSwitchPetFreePickAndInsufficient 积分不足分支与免费自选分支。
func TestClassroomSwitchPetFreePickAndInsufficient(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	svc := services.NewPetService(db, services.NewScope(db))

	student := models.Student{ClassID: f.Class.ID, Name: "小明", TotalScore: 10, Status: "active"}
	require.NoError(t, db.Create(&student).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: f.Class.ID, Name: "小明的萌宠", Species: "zhulong", Level: 12, Mood: 80,
	}).Error)

	// 积分不足（Lv.12 → 60 分 > 10 分）。
	_, err := svc.SwitchForClassroom(f.Class.ID, student.ID, "qilin")
	assert.Equal(t, "积分不足，更换宠物需 60 积分", appErrorOf(t, err, 400))
	assert.Equal(t, int64(0), countPetCollections(t, db, student.ID), "积分不足不应写图鉴")

	// 免费自选 → cost 0、机会消费、图鉴照常归档/激活。
	require.NoError(t, models.GrantPetFreePick(db, student.ID, time.Now()))
	free, err := svc.SwitchForClassroom(f.Class.ID, student.ID, "qilin")
	require.NoError(t, err)
	assert.Equal(t, "✅ 已使用整班切换的免费自选机会！", free.Message)
	assert.Equal(t, 0, free.Data.Cost)
	assert.Equal(t, 10, free.Data.TotalScore)

	has, err := models.HasPetFreePick(db, student.ID, time.Now())
	require.NoError(t, err)
	assert.False(t, has)

	archived := petCollectionOf(t, db, student.ID, "zhulong")
	require.NotNil(t, archived)
	assert.False(t, archived.IsActive)
	active := petCollectionOf(t, db, student.ID, "qilin")
	require.NotNil(t, active)
	assert.True(t, active.IsActive)
}
