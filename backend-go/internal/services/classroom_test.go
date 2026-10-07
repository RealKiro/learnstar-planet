// 教师端「我的班级 / 切换班级 / 模式」测试。
package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newClassroomService(db *gorm.DB) *services.ClassroomService {
	return services.NewClassroomService(db, services.NewScope(db))
}

// seedBot 创建一个 API 机器人教师账号。
func seedBot(t *testing.T, db *gorm.DB, schoolID uint, username string) models.User {
	t.Helper()
	u := models.User{
		SchoolID: schoolID, Role: "teacher", Username: username,
		Name: "机器人", Status: "active", IsAPIBot: true,
	}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("create bot: %v", err)
	}
	return u
}

func TestMyClasses_Teacher(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	other := seedTeacher(t, db, school.ID, "t2")

	mine, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	seedTeacherClass(t, db, school.ID, other.ID, "二班") // 他人班级

	svc := newClassroomService(db)
	views, err := svc.MyClasses(&teacher)
	require.NoError(t, err)
	require.Len(t, views, 1)
	assert.Equal(t, mine.ID, views[0].ClassID)
	assert.Equal(t, "一班", views[0].ClassName)
	assert.Equal(t, "head_teacher", views[0].Role)
}

func TestMyClasses_ExcludesInactiveClass(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")

	active, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	archived := models.ClassRoom{SchoolID: school.ID, Name: "已归档班", TeacherID: &teacher.ID, Status: "archived"}
	require.NoError(t, db.Create(&archived).Error)

	svc := newClassroomService(db)
	views, err := svc.MyClasses(&teacher)
	require.NoError(t, err)
	require.Len(t, views, 1)
	assert.Equal(t, active.ID, views[0].ClassID)
}

func TestMyClasses_APIBotSeesWholeSchool(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	seedClass(t, db, school.ID) // 无班主任的班级

	other := seedOtherSchool(t, db)
	seedTeacherClass(t, db, other.ID, teacher.ID, "他校班")

	bot := seedBot(t, db, school.ID, "api-bot")
	svc := newClassroomService(db)

	views, err := svc.MyClasses(&bot)
	require.NoError(t, err)
	require.Len(t, views, 2, "机器人应看到本校全部启用班级，且不含他校")
	for _, v := range views {
		assert.Equal(t, "api_bot", v.Role)
	}
}

func TestSwitchTo_Success(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := newClassroomService(db)
	require.NoError(t, svc.SwitchTo(&teacher, class.ID))

	// 落库校验：重新读取用户，设置项应持久化。
	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.Equal(t, class.ID, reloaded.SettingUint(services.ActiveClassSettingKey))

	mode, err := svc.GetMode(&reloaded)
	require.NoError(t, err)
	require.NotNil(t, mode.ActiveClassID)
	assert.Equal(t, class.ID, *mode.ActiveClassID)
}

func TestSwitchTo_NotAssigned(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	otherClass := seedClass(t, db, school.ID)

	svc := newClassroomService(db)
	err := svc.SwitchTo(&teacher, otherClass.ID)
	assertAppStatus(t, err, 403, "未分配的班级应返回 403")

	_, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, "您未被分配到此班级", err.Error())

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.Equal(t, uint(0), reloaded.SettingUint(services.ActiveClassSettingKey), "失败不应写入设置")
}

func TestSwitchTo_APIBotCanUseAnySchoolClass(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	bot := seedBot(t, db, school.ID, "api-bot")
	anyClass := seedClass(t, db, school.ID)
	foreign := seedClass(t, db, school.ID+1)

	svc := newClassroomService(db)
	require.NoError(t, svc.SwitchTo(&bot, anyClass.ID))

	err := svc.SwitchTo(&bot, foreign.ID)
	assertAppStatus(t, err, 403, "他校班级仍应拒绝")
}

func TestGetMode_Default(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")

	svc := newClassroomService(db)
	mode, err := svc.GetMode(&teacher)
	require.NoError(t, err)
	assert.Equal(t, services.ModeClassroomDisplay, mode.Mode, "默认应为班级大屏模式")
	assert.Nil(t, mode.ActiveClassID, "未选择班级时应为 null")
}

func TestSetMode_Validation(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	svc := newClassroomService(db)

	_, err := svc.SetMode(&teacher, "unknown_mode", 0, "pw")
	assertAppStatus(t, err, 422, "非法模式应 422")

	_, err = svc.SetMode(&teacher, services.ModeTeacherManage, 0, "")
	assertAppStatus(t, err, 422, "切到教师管理模式必须带密码")

	// 班级大屏模式无需密码。
	result, err := svc.SetMode(&teacher, services.ModeClassroomDisplay, 0, "")
	require.NoError(t, err)
	assert.Equal(t, "已切换为班级大屏模式", result.Message)
	assert.Equal(t, services.ModeClassroomDisplay, result.Data.Mode)
}

func TestSetMode_SwitchesActiveClassWhenAssigned(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	foreign := seedClass(t, db, school.ID)

	svc := newClassroomService(db)

	result, err := svc.SetMode(&teacher, services.ModeTeacherManage, class.ID, "any-password")
	require.NoError(t, err)
	assert.Equal(t, "已切换为教师管理模式", result.Message)
	assert.Equal(t, services.ModeTeacherManage, result.Data.Mode)
	require.NotNil(t, result.Data.ActiveClassID)
	assert.Equal(t, class.ID, *result.Data.ActiveClassID)

	// 切到未分配的班级：模式生效但激活班级保持原值（同 Laravel 静默忽略）。
	result, err = svc.SetMode(&teacher, services.ModeClassroomDisplay, foreign.ID, "")
	require.NoError(t, err)
	require.NotNil(t, result.Data.ActiveClassID)
	assert.Equal(t, class.ID, *result.Data.ActiveClassID, "未分配班级不应改变激活班级")
}

func TestSetMode_PersistsAndMergesSettings(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	svc := newClassroomService(db)

	require.NoError(t, svc.SwitchTo(&teacher, class.ID))
	_, err := svc.SetMode(&teacher, services.ModeTeacherManage, 0, "pw")
	require.NoError(t, err)

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.Equal(t, services.ModeTeacherManage, reloaded.SettingString(services.DisplayModeSettingKey, ""))
	assert.Equal(t, class.ID, reloaded.SettingUint(services.ActiveClassSettingKey),
		"切换模式不得清掉已有的 active_class_id")

	mode, err := svc.GetMode(&reloaded)
	require.NoError(t, err)
	assert.Equal(t, services.ModeTeacherManage, mode.Mode)
	require.NotNil(t, mode.ActiveClassID)
	assert.Equal(t, class.ID, *mode.ActiveClassID)
}

func TestUserSettingsHelpers(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")

	// 空设置：默认值生效。
	assert.Equal(t, "classroom_display", teacher.SettingString("display_mode", "classroom_display"))
	assert.Equal(t, uint(0), teacher.SettingUint("active_class_id"))

	raw, err := teacher.WithSetting("active_class_id", 42)
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", teacher.ID).Update("settings", raw).Error)

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.Equal(t, uint(42), reloaded.SettingUint("active_class_id"))

	// 合并写入不得丢失其他键。
	raw2, err := reloaded.WithSetting("display_mode", "teacher_manage")
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", reloaded.ID).Update("settings", raw2).Error)

	var again models.User
	require.NoError(t, db.First(&again, reloaded.ID).Error)
	assert.Equal(t, uint(42), again.SettingUint("active_class_id"))
	assert.Equal(t, "teacher_manage", again.SettingString("display_mode", ""))

	// 非法 JSON 不应 panic，按空表处理。
	broken := models.User{Settings: "{not-json"}
	assert.Equal(t, uint(0), broken.SettingUint("active_class_id"))
	assert.Equal(t, "def", broken.SettingString("missing", "def"))
}
