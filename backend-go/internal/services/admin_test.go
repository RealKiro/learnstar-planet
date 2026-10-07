// 管理端服务测试：学生创建/更新/删除（覆盖 IN (?) 子查询路径的回归）。
package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateStudent_WithStudentNo(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	admin := services.NewAdmin(db)

	// 回归：带 student_no 时校内查重走 class_id IN (?) 子查询，此前后者的写法会直接 500。
	student, err := admin.CreateStudent(school.ID, class.ID, "小明", "male", "2026001")
	require.NoError(t, err)
	assert.Equal(t, "小明", student.Name)
	assert.Equal(t, "2026001", student.StudentNo)
	assert.Equal(t, "active", student.Status)
	assert.Equal(t, class.ID, student.ClassID)
}

func TestCreateStudent_IdempotentOnSameName(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	admin := services.NewAdmin(db)

	first, err := admin.CreateStudent(school.ID, class.ID, "小明", "", "1")
	require.NoError(t, err)

	second, err := admin.CreateStudent(school.ID, class.ID, "小明", "", "999")
	require.NoError(t, err)
	assert.Equal(t, first.ID, second.ID, "同班同名应直接返回既有学生")

	var count int64
	require.NoError(t, db.Model(&models.Student{}).Where("class_id = ?", class.ID).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestCreateStudent_CrossClassStudentNoBlocked(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	classA := seedClass(t, db, school.ID)
	classB := models.ClassRoom{SchoolID: school.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&classB).Error)
	admin := services.NewAdmin(db)

	_, err := admin.CreateStudent(school.ID, classA.ID, "小明", "", "1001")
	require.NoError(t, err)

	_, err = admin.CreateStudent(school.ID, classB.ID, "小红", "", "1001")
	assertAppStatus(t, err, 422, "同学号跨班应被拦截并提示走批量转班")
	assert.Contains(t, err.Error(), "1001")
}

func TestCreateStudent_OtherSchoolClassRejected(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	foreign := seedOtherSchool(t, db)
	foreignClass := seedClass(t, db, foreign.ID)
	admin := services.NewAdmin(db)

	_, err := admin.CreateStudent(school.ID, foreignClass.ID, "小明", "", "1")
	assertAppStatus(t, err, 404, "跨校班级不可见")
}

func TestCreateStudent_EmptyNameRejected(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	admin := services.NewAdmin(db)

	_, err := admin.CreateStudent(school.ID, class.ID, "", "", "")
	assertAppStatus(t, err, 422, "姓名必填")
}

func TestUpdateAndDeleteStudent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	admin := services.NewAdmin(db)

	student, err := admin.CreateStudent(school.ID, class.ID, "小明", "", "1001")
	require.NoError(t, err)

	// 回归：更新/删除都经 studentInSchool 的 class_id IN (?) 子查询。
	updated, err := admin.UpdateStudent(school.ID, student.ID, map[string]any{"name": "小明明", "student_no": "1002"})
	require.NoError(t, err)
	assert.Equal(t, "小明明", updated.Name)
	assert.Equal(t, "1002", updated.StudentNo)

	require.NoError(t, admin.DeleteStudent(school.ID, student.ID))

	var count int64
	require.NoError(t, db.Model(&models.Student{}).Where("id = ? AND status = ?", student.ID, "active").Count(&count).Error)
	assert.Equal(t, int64(0), count, "删除后不应再有活跃学生")

	_, err = admin.UpdateStudent(school.ID, student.ID, map[string]any{"name": "ghost"})
	assertAppStatus(t, err, 404, "已删除学生不可更新")
}

func TestDeleteStudent_CrossSchoolRejected(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	foreign := seedOtherSchool(t, db)
	foreignClass := seedClass(t, db, foreign.ID)
	admin := services.NewAdmin(db)

	student := models.Student{ClassID: foreignClass.ID, Name: "他校学生", Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	err := admin.DeleteStudent(school.ID, student.ID)
	assertAppStatus(t, err, 404, "跨校学生不可见")

	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, student.ID).Error)
	assert.Equal(t, "active", reloaded.Status, "越权删除不应生效")
}
