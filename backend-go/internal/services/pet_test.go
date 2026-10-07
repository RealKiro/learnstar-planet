package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"gorm.io/gorm"
)

// seedTeacherWithClass 创建一个教师及其担任班主任的班级，并返回教师。
func seedTeacherWithClass(t *testing.T, db *gorm.DB, schoolID uint) (models.User, models.ClassRoom) {
	t.Helper()
	teacher := models.User{
		SchoolID: schoolID, Role: "teacher", Username: "teacher-1",
		Name: "王老师", Status: "active",
	}
	if err := db.Create(&teacher).Error; err != nil {
		t.Fatalf("create teacher: %v", err)
	}
	class := models.ClassRoom{
		SchoolID: schoolID, Name: "一班", TeacherID: &teacher.ID, Status: "active",
	}
	if err := db.Create(&class).Error; err != nil {
		t.Fatalf("create class: %v", err)
	}
	return teacher, class
}

func TestPetClassOverview(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher, class := seedTeacherWithClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "宠", Species: "zhulong", Level: 3, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewPetService(db, services.NewScope(db))
	items, err := svc.ClassOverview(&teacher)
	if err != nil {
		t.Fatalf("ClassOverview: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].StudentName != "小明" || items[0].StageName != "幼年" || items[0].Emoji != "🐣" {
		t.Fatalf("item = %+v", items[0])
	}
}

func TestPetFeed(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher, class := seedTeacherWithClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "宠", Species: "zhulong", Level: 1, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewPetService(db, services.NewScope(db))
	res, err := svc.Feed(&teacher, student.ID)
	if err != nil {
		t.Fatalf("Feed: %v", err)
	}
	if res["mood"] != 100 {
		t.Fatalf("mood = %v, want 100", res["mood"])
	}

	var got models.Pet
	db.First(&got, pet.ID)
	if got.Mood != 100 {
		t.Fatalf("pet.Mood = %d, want 100", got.Mood)
	}
}

func TestPetRename(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher, class := seedTeacherWithClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "旧名", Species: "zhulong", Level: 1, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewPetService(db, services.NewScope(db))
	if _, err := svc.Rename(&teacher, student.ID, "新名"); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	var got models.Pet
	db.First(&got, pet.ID)
	if got.Name != "新名" {
		t.Fatalf("pet.Name = %s, want 新名", got.Name)
	}
}

// Switch：有宠物 → 按等级扣分、旧物种归档、目标物种无图鉴记录时初始化为 1/0/80。
func TestPetSwitchDeductsScore(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher, class := seedTeacherWithClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	student.TotalScore = 20
	db.Save(&student)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "烛龙", Species: "zhulong", Level: 3, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewPetService(db, services.NewScope(db))
	res, err := svc.Switch(&teacher, student.ID, "huohu", "火狐")
	if err != nil {
		t.Fatalf("Switch: %v", err)
	}
	if res.Data.Cost == nil || *res.Data.Cost != 15 { // SwitchCost(3) = 15
		t.Fatalf("cost = %v, want 15", res.Data.Cost)
	}
	if res.Message != "宠物已更换为「火狐」（扣除 15 积分）" {
		t.Fatalf("message = %q", res.Message)
	}
	if res.Data.FreePickUsed == nil || *res.Data.FreePickUsed {
		t.Fatalf("free_pick_used = %v, want false", res.Data.FreePickUsed)
	}

	var gotStudent models.Student
	db.First(&gotStudent, student.ID)
	if gotStudent.TotalScore != 5 {
		t.Fatalf("student.TotalScore = %d, want 5", gotStudent.TotalScore)
	}

	var got models.Pet
	db.First(&got, pet.ID)
	if got.Species != "huohu" || got.Name != "火狐" || got.Level != 1 || got.Mood != 80 {
		t.Fatalf("pet = %+v", got)
	}
	if got.LastSwitchedAt == nil {
		t.Fatalf("last_switched_at 应被写入")
	}

	// 旧物种被归档（is_active=false），目标物种标记激活。
	var archived models.PetCollection
	if err := db.Where("student_id = ? AND species = ?", student.ID, "zhulong").First(&archived).Error; err != nil {
		t.Fatalf("旧物种应写入图鉴: %v", err)
	}
	if archived.Level != 3 || archived.Experience != 0 || archived.Mood != 80 || archived.IsActive {
		t.Fatalf("归档行 = %+v", archived)
	}
	var active models.PetCollection
	if err := db.Where("student_id = ? AND species = ?", student.ID, "huohu").First(&active).Error; err != nil {
		t.Fatalf("新物种应写入图鉴: %v", err)
	}
	if active.Level != 1 || !active.IsActive {
		t.Fatalf("激活行 = %+v", active)
	}
}

// Switch：积分不足 → 400，且不写图鉴、不换物种。
func TestPetSwitchInsufficientScore(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher, class := seedTeacherWithClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	student.TotalScore = 3
	db.Save(&student)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "烛龙", Species: "zhulong", Level: 3, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewPetService(db, services.NewScope(db))
	_, err := svc.Switch(&teacher, student.ID, "huohu", "火狐")
	ae, ok := services.AsAppError(err)
	if !ok || ae.Status != 400 || ae.Message != "积分不足，更换宠物需 15 积分" {
		t.Fatalf("Switch err = %v, want 400 AppError", err)
	}

	var got models.Pet
	db.First(&got, pet.ID)
	if got.Species != "zhulong" {
		t.Fatalf("species = %s, want unchanged zhulong", got.Species)
	}
	var count int64
	db.Model(&models.PetCollection{}).Count(&count)
	if count != 0 {
		t.Fatalf("积分不足不应写图鉴，got %d 行", count)
	}
}
