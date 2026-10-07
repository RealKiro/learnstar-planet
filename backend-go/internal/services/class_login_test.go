// 班级码登录测试（Laravel AuthController::classLogin）。
package services_test

import (
	"strings"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 成功：签发 class_<班级码>_<32 位> token，落库 display_tokens（24 小时），
// 且该 token 可直接用于 display 接口（共用同一张表 / 同一套校验）。
func TestClassLoginSuccessAndTokenUsableForDisplay(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	seedStudent(t, db, f.Class.ID)
	require.NoError(t, db.Create(&models.Student{ClassID: f.Class.ID, Name: "休学", Status: "inactive"}).Error)

	// 大小写与首尾空格都会被规范化。
	result, err := f.Svc.ClassLogin("  ls11 ")
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result.Token, "class_LS11_"), "token=%q", result.Token)
	assert.Len(t, result.Token, len("class_LS11_")+32)
	assert.Equal(t, f.Class.ID, result.ClassID)
	assert.Equal(t, "一年级（1）班", result.ClassName)
	assert.Equal(t, "一年级", result.Grade)
	assert.Equal(t, int64(1), result.StudentCount, "只统计 active 学生")

	// 复用 display_tokens 表：24 小时有效。
	var row models.DisplayToken
	require.NoError(t, db.Where("token = ?", result.Token).First(&row).Error)
	assert.Equal(t, f.Class.ID, row.ClassID)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), row.ExpiresAt, 2*time.Minute)

	// 当大屏 token 用：校验 + 调 display 接口。
	classID, err := f.Svc.ValidateToken(result.Token)
	require.NoError(t, err)
	assert.Equal(t, f.Class.ID, classID)

	// 当大屏 token 用：调 display 读接口（班级设置）。
	settings, err := f.Svc.ClassSettings(classID)
	require.NoError(t, err)
	assert.Nil(t, settings.PetSeries)

	// 教室端学生列表只含 active 学生（与大屏一致）。
	students, err := f.Svc.ClassroomStudents(classID)
	require.NoError(t, err)
	require.Len(t, students, 1)
	assert.Equal(t, "小明", students[0].Name)
	// 同一班级码第二次登录会签发另一个 token（Laravel 每次 Str::random(32)）。
	again, err := f.Svc.ClassLogin("LS11")
	require.NoError(t, err)
	assert.NotEqual(t, result.Token, again.Token)
}

// 无效班级码 → 401，文案与大屏登录（「请检查后重试」）逐字不同。
func TestClassLoginInvalidCode(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	for _, code := range []string{"LS99", "", "   "} {
		_, err := f.Svc.ClassLogin(code)
		assert.Equal(t, "班级码无效，请核对后重试", appErrorOf(t, err, 401), "code=%q", code)
	}
}

// 多校同码 → 401（构造成两所不同学校都有同一班级码）。
func TestClassLoginDuplicateCodeAcrossSchools(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	other := seedOtherSchool(t, db)
	dup := models.ClassRoom{
		SchoolID: other.ID, Grade: "一年级", Name: "一年级（1）班", Status: "active", DisplayCode: "LS11",
	}
	require.NoError(t, db.Create(&dup).Error)

	_, err := f.Svc.ClassLogin("LS11")
	assert.Equal(t, "班级码在多所学校存在，请向班主任确认", appErrorOf(t, err, 401))

	// 同校同码两条也会命中「多所」守卫（Laravel 只按 count>1 判断）。
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", dup.ID).
		Update("school_id", f.School.ID).Error)
	_, err = f.Svc.ClassLogin("LS11")
	assert.Equal(t, "班级码在多所学校存在，请向班主任确认", appErrorOf(t, err, 401))
}

// display_code 未刷新时按确定性计算码登录（一年级（2）班 → LS12）。
func TestClassLoginByComputedCode(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	fresh := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&fresh).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: fresh.ID, Name: "小红", Status: "active"}).Error)

	result, err := f.Svc.ClassLogin("ls12")
	require.NoError(t, err)
	assert.Equal(t, fresh.ID, result.ClassID)
	assert.Equal(t, "一年级（2）班", result.ClassName)
	assert.Equal(t, int64(1), result.StudentCount)
}
