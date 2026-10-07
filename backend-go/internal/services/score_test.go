package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
)

func TestGiveScoreCreatesScoreAndLog(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	svc := services.NewScoreService(db)
	score, err := svc.GiveScore(&student, 5, "举手发言", 1, nil)
	if err != nil {
		t.Fatalf("GiveScore: %v", err)
	}
	if score.Amount != 5 || score.StudentID != student.ID || score.ClassID != class.ID {
		t.Fatalf("score = %+v", score)
	}
	if student.TotalScore != 5 {
		t.Fatalf("student.TotalScore = %d, want 5", student.TotalScore)
	}

	var log models.ScoreLog
	if err := db.Where("score_id = ?", score.ID).First(&log).Error; err != nil {
		t.Fatalf("find score log: %v", err)
	}
	if log.BalanceBefore != 0 || log.BalanceAfter != 5 {
		t.Fatalf("log balance = %d -> %d, want 0 -> 5", log.BalanceBefore, log.BalanceAfter)
	}
}

func TestGiveScoreNegativeClampsToZero(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	svc := services.NewScoreService(db)
	if _, err := svc.GiveScore(&student, -10, "扣分", 1, nil); err != nil {
		t.Fatalf("GiveScore: %v", err)
	}
	if student.TotalScore != 0 {
		t.Fatalf("student.TotalScore = %d, want 0", student.TotalScore)
	}
}

func TestGiveScoreSyncsPetExperienceAndLevel(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "宠", Species: "zhulong", Level: 1, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewScoreService(db)
	if _, err := svc.GiveScore(&student, 15, "挑战难题", 1, nil); err != nil {
		t.Fatalf("GiveScore: %v", err)
	}

	var got models.Pet
	if err := db.First(&got, pet.ID).Error; err != nil {
		t.Fatalf("reload pet: %v", err)
	}
	// +15 经验：Lv1 需 20，故经验 15；积分 15 对应等级 2。
	if got.Level != 2 || got.Experience != 15 {
		t.Fatalf("pet = level %d exp %d, want level 2 exp 15", got.Level, got.Experience)
	}
}

func TestSpendScore(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	pet := models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "宠", Species: "zhulong", Level: 1, Mood: 80}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("create pet: %v", err)
	}

	svc := services.NewScoreService(db)
	if _, err := svc.GiveScore(&student, 20, "加分", 1, nil); err != nil {
		t.Fatalf("GiveScore: %v", err)
	}
	if _, err := svc.SpendScore(&student, 10, "兑换橡皮", 1); err != nil {
		t.Fatalf("SpendScore: %v", err)
	}
	if student.TotalScore != 10 {
		t.Fatalf("student.TotalScore = %d, want 10", student.TotalScore)
	}

	var got models.Pet
	if err := db.First(&got, pet.ID).Error; err != nil {
		t.Fatalf("reload pet: %v", err)
	}
	// 20 分升到 Lv2 exp 0；扣 10 → Lv1 exp 10；积分 10 对应等级 1。
	if got.Level != 1 || got.Experience != 10 {
		t.Fatalf("pet = level %d exp %d, want level 1 exp 10", got.Level, got.Experience)
	}
}

func TestSpendScoreInsufficient(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	svc := services.NewScoreService(db)
	if _, err := svc.GiveScore(&student, 5, "加分", 1, nil); err != nil {
		t.Fatalf("GiveScore: %v", err)
	}
	_, err := svc.SpendScore(&student, 10, "太贵", 1)
	ae, ok := services.AsAppError(err)
	if !ok || ae.Status != 400 {
		t.Fatalf("SpendScore err = %v, want 400 AppError", err)
	}
	if student.TotalScore != 5 {
		t.Fatalf("student.TotalScore = %d, want 5 (unchanged)", student.TotalScore)
	}
}

func TestBatchGive(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	s1 := seedStudent(t, db, class.ID)
	s2 := seedStudent(t, db, class.ID)

	svc := services.NewScoreService(db)
	count, err := svc.BatchGive([]*models.Student{&s1, &s2}, 3, "全勤", 1, nil)
	if err != nil {
		t.Fatalf("BatchGive: %v", err)
	}
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}

	var r1, r2 models.Student
	db.First(&r1, s1.ID)
	db.First(&r2, s2.ID)
	if r1.TotalScore != 3 || r2.TotalScore != 3 {
		t.Fatalf("totals = %d, %d, want 3, 3", r1.TotalScore, r2.TotalScore)
	}
}

func TestUndoScore(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	svc := services.NewScoreService(db)
	score, err := svc.GiveScore(&student, 5, "加分", 1, nil)
	if err != nil {
		t.Fatalf("GiveScore: %v", err)
	}

	undo, err := svc.Undo(score, 1)
	if err != nil {
		t.Fatalf("Undo: %v", err)
	}
	if undo.Amount != -5 {
		t.Fatalf("undo.Amount = %d, want -5", undo.Amount)
	}

	var reloaded models.Student
	db.First(&reloaded, student.ID)
	if reloaded.TotalScore != 0 {
		t.Fatalf("student.TotalScore = %d, want 0", reloaded.TotalScore)
	}
}
