package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DELETE /admin/teachers/:id（Laravel disableTeacher）：解除班级关联后硬删除；
// API 机器人账号 403；不在本校 / 不是教师 / 不存在 → 404。
func TestAdminOpsDeleteTeacher(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	teacher := makeTeacher(t, db, school.ID, "del-teacher", "待删老师", false)
	class := makeClass(t, db, school.ID, "一年级（1）班", "一年级")

	// 班主任关联 + 任课关联都要被解除。
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Update("teacher_id", teacher.ID).Error)
	require.NoError(t, db.Create(&models.ClassRoomTeacher{
		ClassRoomID: class.ID, UserID: teacher.ID, Role: "head_teacher",
	}).Error)

	require.NoError(t, ops.DeleteTeacher(school.ID, teacher.ID))

	var userCount int64
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", teacher.ID).Count(&userCount).Error)
	assert.Equal(t, int64(0), userCount, "用户应被硬删除")

	var classTeacherID *uint
	var reloaded models.ClassRoom
	require.NoError(t, db.First(&reloaded, class.ID).Error)
	classTeacherID = reloaded.TeacherID
	assert.Nil(t, classTeacherID, "class_rooms.teacher_id 应置空")

	var linkCount int64
	require.NoError(t, db.Model(&models.ClassRoomTeacher{}).Where("user_id = ?", teacher.ID).
		Count(&linkCount).Error)
	assert.Equal(t, int64(0), linkCount, "class_room_teachers 关联应删除")

	// 幂等性：再次删除 → 404 教师不存在。
	err := ops.DeleteTeacher(school.ID, teacher.ID)
	require.Error(t, err)
	appErr, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Status)

	// 机器人账号受保护：403 且文案同 Laravel。
	bot := makeTeacher(t, db, school.ID, "bot-teacher", "机器人", true)
	err = ops.DeleteTeacher(school.ID, bot.ID)
	require.Error(t, err)
	appErr, ok = services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 403, appErr.Status)
	assert.Equal(t, "API 机器人账号不可删除。如需停用，请在 .env 中设置 BOT_ENABLED=false 后重启", appErr.Message)

	// 其他学校的教师不可删（越权 → 404）。
	other := seedExtraSchool(t, db, "delete-teacher-other")
	foreign := makeTeacher(t, db, other.ID, "foreign-teacher", "外校老师", false)
	err = ops.DeleteTeacher(school.ID, foreign.ID)
	require.Error(t, err)
	_, ok = services.AsAppError(err)
	require.True(t, ok)
}
