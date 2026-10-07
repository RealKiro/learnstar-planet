package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
)

func TestLeaderboardTotal(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)

	s1 := models.Student{ClassID: class.ID, Name: "甲", TotalScore: 10, Status: "active"}
	s2 := models.Student{ClassID: class.ID, Name: "乙", TotalScore: 30, Status: "active"}
	s3 := models.Student{ClassID: class.ID, Name: "丙", TotalScore: 20, Status: "active"}
	db.Create(&s1)
	db.Create(&s2)
	db.Create(&s3)

	svc := services.NewLeaderboardService(db)
	entries, err := svc.Total(class.ID, 10)
	if err != nil {
		t.Fatalf("Total: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("len = %d, want 3", len(entries))
	}
	wantOrder := []struct {
		id    uint
		score int
		rank  int
	}{
		{s2.ID, 30, 1}, {s3.ID, 20, 2}, {s1.ID, 10, 3},
	}
	for i, w := range wantOrder {
		if entries[i].StudentID != w.id || entries[i].Score != w.score || entries[i].Rank != w.rank {
			t.Fatalf("entries[%d] = %+v, want id %d score %d rank %d", i, entries[i], w.id, w.score, w.rank)
		}
	}
}

func TestLeaderboardWeekly(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	s1 := seedStudent(t, db, class.ID)
	s2 := seedStudent(t, db, class.ID)

	now := util.Now()
	seedScore := func(studentID uint, amount int) {
		if err := db.Create(&models.Score{
			StudentID: studentID, ClassID: class.ID, Amount: amount,
			CreatedAt: now, UpdatedAt: now,
		}).Error; err != nil {
			t.Fatalf("create score: %v", err)
		}
	}
	seedScore(s1.ID, 10)
	seedScore(s1.ID, 5)
	seedScore(s2.ID, 20)

	svc := services.NewLeaderboardService(db)
	entries, err := svc.Weekly(class.ID, 10)
	if err != nil {
		t.Fatalf("Weekly: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}
	if entries[0].StudentID != s2.ID || entries[0].Score != 20 {
		t.Fatalf("entries[0] = %+v, want s2 20", entries[0])
	}
	if entries[1].StudentID != s1.ID || entries[1].Score != 15 {
		t.Fatalf("entries[1] = %+v, want s1 15", entries[1])
	}
}

func TestLeaderboardPetLevel(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)

	s1 := seedStudent(t, db, class.ID)
	s2 := seedStudent(t, db, class.ID)
	s3 := seedStudent(t, db, class.ID)
	db.Create(&models.Pet{StudentID: s1.ID, ClassID: class.ID, Name: "a", Level: 1})
	db.Create(&models.Pet{StudentID: s2.ID, ClassID: class.ID, Name: "b", Level: 5})
	db.Create(&models.Pet{StudentID: s3.ID, ClassID: class.ID, Name: "c", Level: 3})

	svc := services.NewLeaderboardService(db)
	entries, err := svc.PetLevel(class.ID, 10)
	if err != nil {
		t.Fatalf("PetLevel: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("len = %d, want 3", len(entries))
	}
	wantOrder := []struct {
		id    uint
		level int
	}{
		{s2.ID, 5}, {s3.ID, 3}, {s1.ID, 1},
	}
	for i, w := range wantOrder {
		if entries[i].StudentID != w.id || entries[i].PetLevel != w.level {
			t.Fatalf("entries[%d] = %+v, want id %d level %d", i, entries[i], w.id, w.level)
		}
	}
}

func TestLeaderboardLimit(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	for i := 0; i < 5; i++ {
		db.Create(&models.Student{ClassID: class.ID, Name: "s", TotalScore: i, Status: "active"})
	}

	svc := services.NewLeaderboardService(db)
	entries, err := svc.Total(class.ID, 3)
	if err != nil {
		t.Fatalf("Total: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("len = %d, want 3", len(entries))
	}
}
