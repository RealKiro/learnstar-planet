package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupDB 打开内存 SQLite 并迁移全部模型，供服务层测试复用。
func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1) // 内存库需单连接，否则每连接各自独立。

	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seedSchool(t *testing.T, db *gorm.DB) models.School {
	t.Helper()
	school := models.School{Name: "测试学校", Code: "test-school", Status: "active"}
	if err := db.Create(&school).Error; err != nil {
		t.Fatalf("create school: %v", err)
	}
	return school
}

func seedClass(t *testing.T, db *gorm.DB, schoolID uint) models.ClassRoom {
	t.Helper()
	class := models.ClassRoom{SchoolID: schoolID, Name: "一班", Status: "active"}
	if err := db.Create(&class).Error; err != nil {
		t.Fatalf("create class: %v", err)
	}
	return class
}

func seedStudent(t *testing.T, db *gorm.DB, classID uint) models.Student {
	t.Helper()
	student := models.Student{ClassID: classID, Name: "小明", Status: "active"}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("create student: %v", err)
	}
	return student
}

// seedTeacher 创建一个教师用户。
func seedTeacher(t *testing.T, db *gorm.DB, schoolID uint, username string) models.User {
	t.Helper()
	u := models.User{SchoolID: schoolID, Role: "teacher", Username: username, Name: "李老师", Status: "active"}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("create teacher: %v", err)
	}
	return u
}

// seedTeacherClass 创建归属指定教师的可管辖班级，并放入一个活跃学生。
func seedTeacherClass(t *testing.T, db *gorm.DB, schoolID, teacherID uint, className string) (models.ClassRoom, models.Student) {
	t.Helper()
	class := models.ClassRoom{SchoolID: schoolID, Name: className, TeacherID: &teacherID, Status: "active"}
	if err := db.Create(&class).Error; err != nil {
		t.Fatalf("create class: %v", err)
	}
	student := models.Student{ClassID: class.ID, Name: "小明", StudentNo: "001", Status: "active"}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("create student: %v", err)
	}
	return class, student
}

// seedExtraStudent 在同班追加一名活跃学生。
func seedExtraStudent(t *testing.T, db *gorm.DB, classID uint) models.Student {
	t.Helper()
	s := models.Student{ClassID: classID, Name: "小红", StudentNo: "002", Status: "active"}
	if err := db.Create(&s).Error; err != nil {
		t.Fatalf("create extra student: %v", err)
	}
	return s
}
