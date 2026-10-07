package services_test

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcastRecentEmpty(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	results, err := svc.Recent(&teacher)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestBroadcastRecentReturnsOrdered(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t2")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewBroadcastService(db, services.NewScope(db))

	b1 := seedBroadcastRaw(t, db, school.ID, class.ID, teacher.ID, "first")
	b2 := seedBroadcastRaw(t, db, school.ID, class.ID, teacher.ID, "second")

	results, err := svc.Recent(&teacher)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, b2.ID, results[0].ID)
	assert.Equal(t, b1.ID, results[1].ID)
}

func TestBroadcastSendToAll(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t3")
	_, _ = seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	_, _ = seedTeacherClass(t, db, school.ID, teacher.ID, "二班")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	sent, err := svc.Send(&teacher, "hello all", "banner", true, false, 10, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, sent)

	var count int64
	db.Model(&models.Broadcast{}).Count(&count)
	assert.EqualValues(t, 2, count)
}

func TestBroadcastSendSpecificClasses(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t4")
	class1, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	seedTeacherClass(t, db, school.ID, teacher.ID, "二班")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	sent, err := svc.Send(&teacher, "specific", "popup", false, true, 15, []uint{class1.ID})
	require.NoError(t, err)
	assert.Equal(t, 1, sent)

	var broadcasts []models.Broadcast
	db.Find(&broadcasts)
	require.Len(t, broadcasts, 1)
	assert.Equal(t, class1.ID, broadcasts[0].ClassID)
	assert.Equal(t, "popup", broadcasts[0].Type)
	assert.Equal(t, 15, broadcasts[0].DisplaySeconds)
}

func TestBroadcastSendNoTargets(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t5")
	// Teacher has no classes
	otherTeacher := seedTeacher(t, db, school.ID, "other")
	otherClass, _ := seedTeacherClass(t, db, school.ID, otherTeacher.ID, "其他班")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	sent, err := svc.Send(&teacher, "no targets", "banner", true, false, 10, []uint{otherClass.ID})
	require.NoError(t, err)
	assert.Equal(t, 0, sent)
}

func TestBroadcastSendWithInaccessibleFiltered(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t6")
	class1, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	// Other teacher's class — inaccessible
	otherTeacher := seedTeacher(t, db, school.ID, "other6")
	otherClass, _ := seedTeacherClass(t, db, school.ID, otherTeacher.ID, "其他班")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	// Send to both — only accessible one gets it
	sent, err := svc.Send(&teacher, "mixed", "fullscreen", true, false, 5, []uint{class1.ID, otherClass.ID})
	require.NoError(t, err)
	assert.Equal(t, 1, sent)

	var broadcasts []models.Broadcast
	db.Where("class_id = ?", otherClass.ID).Find(&broadcasts)
	assert.Empty(t, broadcasts)
}

func TestBroadcastFindInScope(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t7")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	b := seedBroadcastRaw(t, db, school.ID, class.ID, teacher.ID, "found it")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	got, err := svc.FindInScope(&teacher, b.ID)
	require.NoError(t, err)
	assert.Equal(t, b.ID, got.ID)
}

func TestBroadcastFindInScopeNotFound(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t8")

	otherTeacher := seedTeacher(t, db, school.ID, "other8")
	otherClass, _ := seedTeacherClass(t, db, school.ID, otherTeacher.ID, "其他班")
	b := seedBroadcastRaw(t, db, school.ID, otherClass.ID, otherTeacher.ID, "invisible")

	svc := services.NewBroadcastService(db, services.NewScope(db))
	_, err := svc.FindInScope(&teacher, b.ID)
	require.Error(t, err)
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)
}

// seedBroadcastRaw 直接插入一条广播记录。
func seedBroadcastRaw(t *testing.T, db *gorm.DB, schoolID, classID, teacherID uint, content string) models.Broadcast {
	t.Helper()
	b := models.Broadcast{
		SchoolID:  schoolID,
		ClassID:   classID,
		TeacherID: &teacherID,
		Content:   content,
		Type:      "banner",
		Status:    "sent",
	}
	if err := db.Create(&b).Error; err != nil {
		t.Fatalf("create broadcast: %v", err)
	}
	return b
}
