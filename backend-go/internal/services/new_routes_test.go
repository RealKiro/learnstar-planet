package services_test

import (
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// A. 认证会话（logout / refresh / bindings）
// ============================================================

func TestAuthRefreshRevokesOldTokenAndIssuesNew(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	user := seedAdminUser(t, db, school.ID, "secret123")
	jwtMgr := auth.New("test-secret", 72)
	svc := services.NewAuthService(db, jwtMgr)

	token, _, err := svc.Login("admin", "secret123")
	require.NoError(t, err)
	claims, err := jwtMgr.Parse(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims.JTI(), "新签发的令牌必须带 jti")

	newToken, err := svc.Refresh(&user, claims.JTI(), claims.Expiry())
	require.NoError(t, err)
	require.NotEqual(t, token, newToken)

	var revoked int64
	require.NoError(t, db.Model(&models.RevokedToken{}).
		Where("jti = ? AND expires_at > ?", claims.JTI(), claims.Expiry().Add(-1)).Count(&revoked).Error)
	assert.Equal(t, int64(1), revoked, "旧令牌的 jti 应进入撤销名单")

	newClaims, err := jwtMgr.Parse(newToken)
	require.NoError(t, err)
	assert.NotEqual(t, claims.JTI(), newClaims.JTI())
	assert.Equal(t, claims.UserID, newClaims.UserID)

	// 登出：新令牌的 jti 也进撤销名单（累计 2 条）。
	require.NoError(t, svc.Logout(newClaims.JTI(), newClaims.Expiry()))
	require.NoError(t, db.Model(&models.RevokedToken{}).Count(&revoked).Error)
	assert.Equal(t, int64(2), revoked)
}

func TestAuthLogoutWithoutJTIAndPrunesExpired(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedAdminUser(t, db, school.ID, "secret123")
	svc := services.NewAuthService(db, auth.New("test-secret", 72))

	// 无 jti（历史令牌）→ 不写撤销名单，也不报错。
	require.NoError(t, svc.Logout("", nil))
	var count int64
	require.NoError(t, db.Model(&models.RevokedToken{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)

	// 已过期记录会在下一次撤销时被惰性清理。
	require.NoError(t, db.Create(&models.RevokedToken{JTI: "old", ExpiresAt: nowMinusHour()}).Error)
	require.NoError(t, svc.Logout("fresh", nil))
	require.NoError(t, db.Model(&models.RevokedToken{}).Count(&count).Error)
	assert.Equal(t, int64(1), count, "过期记录应被清理，只剩 fresh")
	var row models.RevokedToken
	require.NoError(t, db.Where("jti = ?", "fresh").First(&row).Error)
}

func TestAuthBindingsDefaultPlatforms(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	svc := services.NewAuthService(db, auth.New("test-secret", 72))

	// 未配置 enabled_third_party_platforms → 默认三平台（企业微信 / 微信 / QQ）。
	views, err := svc.Bindings(&teacher)
	require.NoError(t, err)
	require.Len(t, views, 3)
	assert.Equal(t, "wechat_work", views[0].Platform)
	assert.Equal(t, "企业微信", views[0].Label)
	assert.Equal(t, "🏢", views[0].Icon)
	assert.False(t, views[0].Bound)
	assert.Nil(t, views[0].Nick)
	assert.Equal(t, "wechat", views[1].Platform)
	assert.Equal(t, "qq", views[2].Platform)

	// 绑定企业微信后 bound=true 且带昵称（其余仍为 false）。
	nick := "王老师"
	require.NoError(t, db.Create(&models.ThirdPartyBinding{
		UserID: teacher.ID, Platform: "wechat_work", PlatformID: "ww-1", PlatformNick: nick,
	}).Error)
	views, err = svc.Bindings(&teacher)
	require.NoError(t, err)
	assert.True(t, views[0].Bound)
	require.NotNil(t, views[0].Nick)
	assert.Equal(t, nick, *views[0].Nick)
	assert.False(t, views[1].Bound)
}

func TestAuthBindingsSchoolIntersectionOrder(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	svc := services.NewAuthService(db, auth.New("test-secret", 72))

	// 学校只启用 qq / dingtalk → 交集并按 platforms() 顺序（qq 在 dingtalk 前）。
	settings, err := school.WithSetting("enabled_third_party_platforms", []string{"qq", "dingtalk"})
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).Update("settings", settings).Error)

	views, err := svc.Bindings(&teacher)
	require.NoError(t, err)
	require.Len(t, views, 2)
	assert.Equal(t, "qq", views[0].Platform)
	assert.Equal(t, "🐧", views[0].Icon)
	assert.Equal(t, "dingtalk", views[1].Platform)
	assert.Equal(t, "🔷", views[1].Icon)

	// 空数组同样回退默认三平台（Laravel `is_array && count > 0` 为假）。
	settings, err = school.WithSetting("enabled_third_party_platforms", []string{})
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).Update("settings", settings).Error)
	views, err = svc.Bindings(&teacher)
	require.NoError(t, err)
	assert.Len(t, views, 3)

	// 学校启用的平台全不在 platforms() 白名单内 → 空列表（不回落默认）。
	settings, err = school.WithSetting("enabled_third_party_platforms", []string{"unknown_platform"})
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).Update("settings", settings).Error)
	views, err = svc.Bindings(&teacher)
	require.NoError(t, err)
	assert.Empty(t, views)
}

// ============================================================
// B. 管理端（列表 / 单建学生 / 班级详情 / 教师班级分配 / 查看密码）
// ============================================================

func seedStudentsForList(t *testing.T, db *gorm.DB, schoolID uint) (models.ClassRoom, models.ClassRoom) {
	t.Helper()
	classA := models.ClassRoom{SchoolID: schoolID, Name: "一年级（1）班", Grade: "一年级", Status: "active"}
	classB := models.ClassRoom{SchoolID: schoolID, Name: "二年级（1）班", Grade: "二年级", Status: "active"}
	require.NoError(t, db.Create(&classA).Error)
	require.NoError(t, db.Create(&classB).Error)

	students := []models.Student{
		{ClassID: classA.ID, Name: "小明", StudentNo: "001", Status: "active", TotalScore: 12},
		{ClassID: classA.ID, Name: "小红", StudentNo: "002", Status: "inactive"},
		{ClassID: classB.ID, Name: "小刚", StudentNo: "003", Status: "active"},
	}
	for i := range students {
		require.NoError(t, db.Create(&students[i]).Error)
	}
	require.NoError(t, db.Create(&models.Pet{
		StudentID: students[0].ID, ClassID: classA.ID, Name: "小明的萌宠",
		Species: "zhulong", Level: 3, Experience: 5, Mood: 80,
	}).Error)
	return classA, classB
}

func TestAdminListStudentsPagedDefaultsAndFilters(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	classA, classB := seedStudentsForList(t, db, school.ID)
	admin := services.NewAdmin(db)

	// 默认 status=active、per_page=50、id desc，并带 class_name/pet_* 追加字段。
	page, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{})
	require.NoError(t, err)
	assert.Equal(t, int64(2), page.Meta.Total, "默认只统计 active")
	assert.Equal(t, 1, page.Meta.CurrentPage)
	assert.Equal(t, 1, page.Meta.LastPage)
	assert.Equal(t, 50, page.Meta.PerPage)
	require.Len(t, page.Data, 2)
	assert.Equal(t, "小刚", page.Data[0].Name, "id desc：后建的小刚在前")
	require.NotNil(t, page.Data[0].ClassName)
	assert.Equal(t, "二年级（1）班", *page.Data[0].ClassName)
	require.NotNil(t, page.Data[0].ClassGrade)
	assert.Equal(t, "二年级", *page.Data[0].ClassGrade)
	assert.Equal(t, "", page.Data[0].PetSpecies, "无宠物 → 空物种/0 级/空名")

	row := page.Data[1]
	assert.Equal(t, "小明", row.Name)
	assert.Equal(t, "zhulong", row.PetSpecies)
	assert.Equal(t, 3, row.PetLevel)
	assert.Equal(t, "小明的萌宠", row.PetName)
	require.NotNil(t, row.StudentNo)
	assert.Equal(t, "001", *row.StudentNo)

	// status=all 放开全部。
	all, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{Status: "all"})
	require.NoError(t, err)
	assert.Equal(t, int64(3), all.Meta.Total)

	// search 命中姓名或学号。
	found, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{Search: "小"})
	require.NoError(t, err)
	assert.Equal(t, int64(2), found.Meta.Total)
	byNo, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{Search: "003"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), byNo.Meta.Total)
	assert.Equal(t, "小刚", byNo.Data[0].Name)

	// class_id / grade 筛选。
	byClass, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{ClassID: classA.ID})
	require.NoError(t, err)
	assert.Equal(t, int64(1), byClass.Meta.Total)
	byGrade, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{Grade: "二年级"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), byGrade.Meta.Total)
	assert.Equal(t, classB.ID, byGrade.Data[0].ClassID)

	// 分页元信息：per_page=1 时 last_page=2、第 2 页只 1 条。
	p2, err := admin.ListStudentsPaged(school.ID, services.AdminStudentQuery{PerPage: 1, Page: 2})
	require.NoError(t, err)
	assert.Equal(t, 2, p2.Meta.LastPage)
	assert.Equal(t, 1, p2.Meta.PerPage)
	require.Len(t, p2.Data, 1)

	// 跨校不可见。
	other := models.School{Name: "别校", Code: "other-school", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	empty, err := admin.ListStudentsPaged(other.ID, services.AdminStudentQuery{})
	require.NoError(t, err)
	assert.Equal(t, int64(0), empty.Meta.Total)
	assert.Equal(t, 1, empty.Meta.LastPage, "total=0 时 last_page 为 1（同 Laravel paginator）")
}

func TestAdminCreateStudentInClassAndValidation(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, _ := seedStudentsForList(t, db, school.ID)
	admin := services.NewAdmin(db)

	// 班级不存在（exists 规则失败）→ ValidationError（422 + errors）。
	_, err := admin.CreateStudentInClass(school.ID, 999999, "新同学", "男生", "009")
	ve, ok := services.AsValidationError(err)
	require.True(t, ok, "err = %v", err)
	assert.Contains(t, ve.Errors, "class_id")

	// 跨校班级 → 404。
	other := models.School{Name: "别校", Code: "other-school-2", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	foreignClass := models.ClassRoom{SchoolID: other.ID, Name: "别校一班", Status: "active"}
	require.NoError(t, db.Create(&foreignClass).Error)
	_, err = admin.CreateStudentInClass(school.ID, foreignClass.ID, "新同学", "", "")
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)

	// 成功：性别归一化 + 自动分配萌宠。
	student, err := admin.CreateStudentInClass(school.ID, class.ID, "新同学", "女生", "009")
	require.NoError(t, err)
	assert.Equal(t, "女", student.Gender)
	assert.Equal(t, "active", student.Status)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, 1, pet.Level)
	assert.Equal(t, "新同学的萌宠", pet.Name)
	assert.NotEmpty(t, pet.Species)

	// 性别无法识别 → 未知（同 Laravel else 分支）。
	other2, err := admin.CreateStudentInClass(school.ID, class.ID, "第三个", "外星人", "")
	require.NoError(t, err)
	assert.Equal(t, "未知", other2.Gender)
}

func TestAdminClassDetailIncludesTeacherAndStudents(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, _ := seedStudentsForList(t, db, school.ID)
	admin := services.NewAdmin(db)

	teacher := seedTeacher(t, db, school.ID, "head1")
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Update("teacher_id", teacher.ID).Error)

	detail, err := admin.ClassDetail(school.ID, class.ID)
	require.NoError(t, err)
	require.NotNil(t, detail.Teacher)
	assert.Equal(t, teacher.ID, detail.Teacher.ID)
	assert.Equal(t, class.Name, detail.Name)
	// students 关联：该班全部学生（含已停用，同 Eloquent 关联语义）。
	require.Len(t, detail.Students, 2)
	assert.Equal(t, "小明", detail.Students[0].Name)

	// 跨校 / 不存在 → 404「班级不存在」。
	other := models.School{Name: "别校", Code: "other-school-3", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	_, err = admin.ClassDetail(other.ID, class.ID)
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)
	assert.Equal(t, "班级不存在", ae.Message)
}

func TestAdminAssignTeacherClassesReplaceSemantics(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	classA, classB := seedStudentsForList(t, db, school.ID)
	admin := services.NewAdmin(db)
	teacher := seedTeacher(t, db, school.ID, "assignee")

	subject := "数学"
	synced, err := admin.AssignTeacherClasses(school.ID, teacher.ID, []services.TeacherClassAssignment{
		{ClassID: classA.ID, Role: "head_teacher", Subject: &subject},
		{ClassID: classB.ID, Role: "co_teacher"},
	})
	require.NoError(t, err)
	require.Len(t, synced, 2)
	assert.Equal(t, classA.ID, synced[0].ClassID)
	assert.Equal(t, "一年级（1）班", synced[0].ClassName)
	assert.Equal(t, "head_teacher", synced[0].Role)
	require.NotNil(t, synced[0].Subject)
	assert.Equal(t, "数学", *synced[0].Subject)
	assert.Nil(t, synced[1].Subject, "未传 subject → null")

	// head_teacher 同步 class_rooms.teacher_id。
	var refreshedA models.ClassRoom
	require.NoError(t, db.First(&refreshedA, classA.ID).Error)
	require.NotNil(t, refreshedA.TeacherID)
	assert.Equal(t, teacher.ID, *refreshedA.TeacherID)

	var rows []models.ClassRoomTeacher
	require.NoError(t, db.Where("user_id = ?", teacher.ID).Find(&rows).Error)
	assert.Len(t, rows, 2)

	// replace：只提交 classB 的 grade_lead → classA 的 head_teacher 被删且 teacher_id 置空。
	synced, err = admin.AssignTeacherClasses(school.ID, teacher.ID, []services.TeacherClassAssignment{
		{ClassID: classB.ID, Role: "grade_lead"},
	})
	require.NoError(t, err)
	require.Len(t, synced, 1)
	assert.Equal(t, "grade_lead", synced[0].Role)

	require.NoError(t, db.Where("user_id = ?", teacher.ID).Find(&rows).Error)
	require.Len(t, rows, 1)
	assert.Equal(t, classB.ID, rows[0].ClassRoomID)
	assert.Equal(t, "grade_lead", rows[0].Role)

	require.NoError(t, db.First(&refreshedA, classA.ID).Error)
	assert.Nil(t, refreshedA.TeacherID, "被移除的 head_teacher 班级应清空 teacher_id")

	// 空 assignments → 清空全部（classB 原 role 为 grade_lead，不清 teacher_id）。
	synced, err = admin.AssignTeacherClasses(school.ID, teacher.ID, nil)
	require.NoError(t, err)
	assert.Empty(t, synced)
	require.NoError(t, db.Where("user_id = ?", teacher.ID).Find(&rows).Error)
	assert.Empty(t, rows)
}

func TestAdminAssignTeacherClassesGuards(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, _ := seedStudentsForList(t, db, school.ID)
	admin := services.NewAdmin(db)
	teacher := seedTeacher(t, db, school.ID, "assignee2")

	// 非本校教师 → 404；跨校班级 → 404；非法角色 → 422。
	other := models.School{Name: "别校", Code: "other-school-4", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	_, err := admin.AssignTeacherClasses(other.ID, teacher.ID, nil)
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)

	foreignClass := models.ClassRoom{SchoolID: other.ID, Name: "别校一班", Status: "active"}
	require.NoError(t, db.Create(&foreignClass).Error)
	_, err = admin.AssignTeacherClasses(school.ID, teacher.ID, []services.TeacherClassAssignment{
		{ClassID: foreignClass.ID, Role: "co_teacher"},
	})
	ae, ok = services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)
	assert.Equal(t, "班级不存在", ae.Message)

	_, err = admin.AssignTeacherClasses(school.ID, teacher.ID, []services.TeacherClassAssignment{
		{ClassID: class.ID, Role: "not_a_role"},
	})
	ae, ok = services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 422, ae.Status)

	// 角色非法时不应残留任何关联。
	var count int64
	require.NoError(t, db.Model(&models.ClassRoomTeacher{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestAdminTeacherPasswordAutoGenerateAndReuse(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAdmin(db)
	teacher := seedTeacher(t, db, school.ID, "pw-teacher")

	// 无明文记录 → 自动生成 8 位并写回，password_changed 置回 false。
	password, generated, err := admin.TeacherPassword(school.ID, teacher.ID)
	require.NoError(t, err)
	assert.True(t, generated)
	assert.Len(t, password, 8)

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.Equal(t, password, reloaded.PlainPassword)
	assert.False(t, reloaded.PasswordChanged)

	// 已有明文 → 直接返回同一值，不再生成。
	again, generated, err := admin.TeacherPassword(school.ID, teacher.ID)
	require.NoError(t, err)
	assert.False(t, generated)
	assert.Equal(t, password, again)

	// 跨校 → 404。
	other := models.School{Name: "别校", Code: "other-school-5", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	_, _, err = admin.TeacherPassword(other.ID, teacher.ID)
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)
}

func TestAdminResetPasswordAndCreateTeacherWritePlainPassword(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAdmin(db)

	created, err := admin.CreateTeacher(school.ID, "plain1", "secret123", "明文老师", "")
	require.NoError(t, err)
	assert.Equal(t, "secret123", created.PlainPassword)

	require.NoError(t, admin.ResetPassword(school.ID, created.ID, "newpass123"))
	var reloaded models.User
	require.NoError(t, db.First(&reloaded, created.ID).Error)
	assert.Equal(t, "newpass123", reloaded.PlainPassword)
	assert.False(t, reloaded.PasswordChanged)
}

// ============================================================
// C. 教师端（班级信息卡 / 学生 CRUD 与导入 / 图鉴 / 按规则加减分）
// ============================================================

func TestTeacherClassInfoAndDomainError(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "info-teacher")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "三年级（2）班")
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Updates(map[string]any{"grade": "三年级", "settings": `{"class_points":7,"pet_series":"myth"}`}).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("total_score", 30).Error)

	svc := services.NewPetSeriesService(db, services.NewScope(db))
	view, err := svc.ClassInfo(&teacher)
	require.NoError(t, err)
	assert.Equal(t, class.ID, view.ID)
	assert.Equal(t, "三年级（2）班", view.Name)
	assert.Equal(t, "三年级", view.Grade)
	assert.Equal(t, int64(1), view.StudentCount)
	assert.Equal(t, 30, view.TotalScore)
	assert.Equal(t, 7, view.ClassPoints)
	assert.Equal(t, "LS32", view.DisplayCode)
	require.NotNil(t, view.Settings)
	assert.Equal(t, "myth", view.Settings["pet_series"])

	// 没有可管理的班级 → 400「没有可管理的班级」（Laravel DomainException → 400）。
	lonely := seedTeacher(t, db, school.ID, "lonely-teacher")
	_, err = svc.ClassInfo(&lonely)
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 400, ae.Status)
	assert.Equal(t, "没有可管理的班级", ae.Message)
}

func TestStudentServiceCreateScopeAndNormalization(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "svc-teacher")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	otherClass := seedClass(t, db, school.ID)
	svc := services.NewStudentService(db)

	// 不在管辖班级 → nil（控制器映射 403）。
	created, err := svc.Create([]uint{class.ID}, "越权同学", otherClass.ID, "男", "")
	require.NoError(t, err)
	assert.Nil(t, created)

	student, err := svc.Create([]uint{class.ID}, "新同学", class.ID, "女生", "100")
	require.NoError(t, err)
	require.NotNil(t, student)
	assert.Equal(t, "女", student.Gender)
	assert.Equal(t, "100", student.StudentNo)
	assert.Equal(t, "active", student.Status)

	// 更新：仅改传入字段。
	updated, err := svc.Update([]uint{class.ID}, student.ID, map[string]any{"name": "改名同学"})
	require.NoError(t, err)
	assert.Equal(t, "改名同学", updated.Name)
	assert.Equal(t, "100", updated.StudentNo)

	// 越权更新 → 404。
	_, err = svc.Update([]uint{otherClass.ID}, student.ID, map[string]any{"name": "x"})
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)

	// 删除：硬删除 + 级联清理宠物/积分/审计日志（与 Laravel SoftDeletes 不同）。
	require.NoError(t, db.Create(&models.Pet{StudentID: student.ID, ClassID: class.ID, Name: "p", Species: "zhulong", Level: 1}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: student.ID, ClassID: class.ID, Amount: 3, Reason: "r", GivenBy: teacher.ID}).Error)
	var score models.Score
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&score).Error)
	require.NoError(t, db.Create(&models.ScoreLog{StudentID: student.ID, ScoreID: score.ID, BalanceBefore: 0, BalanceAfter: 3}).Error)

	deleted, err := svc.Delete([]uint{class.ID}, student.ID)
	require.NoError(t, err)
	assert.Equal(t, "改名同学", deleted.Name)

	var count int64
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
	require.NoError(t, db.Model(&models.Pet{}).Where("student_id = ?", student.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
	require.NoError(t, db.Model(&models.Score{}).Where("student_id = ?", student.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
	require.NoError(t, db.Model(&models.ScoreLog{}).Where("student_id = ?", student.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestStudentServiceImportDedupAndConflicts(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "import-teacher")
	class, existing := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	otherClass := models.ClassRoom{SchoolID: school.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	// 二班已有一个学号 777 的学生（跨班冲突源）。
	require.NoError(t, db.Create(&models.Student{ClassID: otherClass.ID, Name: "他班同学", StudentNo: "777", Status: "active"}).Error)

	svc := services.NewStudentService(db)
	result, err := svc.Import([]uint{class.ID, otherClass.ID}, []services.TeacherStudentImportRow{
		{Name: "新同学", ClassName: "一班", Gender: "男生", StudentNo: "200"},        // 成功
		{Name: existing.Name, ClassName: "一班", StudentNo: existing.StudentNo}, // 同班同学号重复 → 跳过
		{Name: "重名", ClassName: "一班", StudentNo: "777"},                       // 跨班同学号 → 跳过
		{Name: "", ClassName: "一班"},                                           // 缺姓名 → 静默忽略
		{Name: "无班级名", ClassName: ""},                                         // 缺班级名 → 静默忽略
		{Name: "外班", ClassName: "不存在的班"},                                      // 班级名不在管辖 → 静默忽略
	})
	require.NoError(t, err)
	assert.Equal(t, 1, result.ImportedCount)
	require.Len(t, result.Skipped, 2)
	assert.Contains(t, result.Skipped[0], "小明")
	assert.Contains(t, result.Skipped[0], "一班 已存在")
	assert.Contains(t, result.Skipped[1], "学号 777")
	assert.Contains(t, result.Skipped[1], "二班")
	assert.Contains(t, result.Skipped[1], "批量转班")
	assert.Equal(t, "成功导入 1 名学生，跳过 2 条重复/冲突记录", result.Message)

	// 导入不分配宠物（同 Laravel StudentService::import）。
	var student models.Student
	require.NoError(t, db.Where("student_no = ?", "200").First(&student).Error)
	assert.Equal(t, "男生", student.Gender, "导入路径不做性别归一化（同 Laravel）")
	var count int64
	require.NoError(t, db.Model(&models.Pet{}).Where("student_id = ?", student.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)

	// 空列表：0 导入、无跳过文案。
	empty, err := svc.Import([]uint{class.ID}, nil)
	require.NoError(t, err)
	assert.Equal(t, 0, empty.ImportedCount)
	assert.Empty(t, empty.Skipped)
	assert.Equal(t, "成功导入 0 名学生", empty.Message)
}

func TestStudentServiceImportGenderDefaultsToUnknown(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "import-teacher2")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewStudentService(db)
	result, err := svc.Import([]uint{class.ID}, []services.TeacherStudentImportRow{
		{Name: "无性别", ClassName: "一班"},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, result.ImportedCount)

	var student models.Student
	require.NoError(t, db.Where("name = ?", "无性别").First(&student).Error)
	assert.Equal(t, "未知", student.Gender)
	assert.Equal(t, "active", student.Status)
}

func TestPetServiceCollection(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "pet-teacher")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Update("settings", `{"pet_series":"myth"}`).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("total_score", 250).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: class.ID, Name: "萌宠", Species: "zhulong",
		Level: 4, Experience: 8, Mood: 70,
	}).Error)

	svc := services.NewPetService(db, services.NewScope(db))
	view, err := svc.Collection(&teacher, student.ID)
	require.NoError(t, err)
	assert.Equal(t, student.ID, view.StudentID)
	assert.Equal(t, "小明", view.StudentName)
	assert.Equal(t, 250, view.TotalScore)
	assert.Equal(t, 3, view.UnlockSlots, "1 + 250/100 = 3")
	require.NotNil(t, view.ClassSeries)
	assert.Equal(t, "myth", *view.ClassSeries)
	require.NotNil(t, view.ActiveSpecies)
	assert.Equal(t, "zhulong", *view.ActiveSpecies)
	require.Len(t, view.Collection, 1, "当前激活宠物被补录进图鉴")
	assert.Equal(t, "zhulong", view.Collection[0].Species)
	assert.Equal(t, 4, view.Collection[0].Level)
	assert.True(t, view.Collection[0].IsActive)

	// 再次调用不重复补录（firstOrCreate 语义）。
	again, err := svc.Collection(&teacher, student.ID)
	require.NoError(t, err)
	assert.Len(t, again.Collection, 1)

	// 无宠物学生：集合为空，active_species 为 null。
	lonely := seedExtraStudent(t, db, class.ID)
	view, err = svc.Collection(&teacher, lonely.ID)
	require.NoError(t, err)
	assert.Nil(t, view.ActiveSpecies)
	assert.Empty(t, view.Collection)
	assert.Equal(t, 1, view.UnlockSlots)

	// 越权 / 不存在 → 404。
	_, err = svc.Collection(&teacher, 999999)
	ae, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, 404, ae.Status)
}

func TestScoreServiceGiveScoreByRule(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "rule-teacher")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: student.ClassID, Name: "萌宠", Species: "zhulong", Level: 1, Mood: 80,
	}).Error)

	rule := models.ScoreRule{SchoolID: &school.ID, Name: "举手发言", Amount: 3, Category: "classroom", IsPositive: true, IsActive: true}
	require.NoError(t, db.Create(&rule).Error)

	svc := services.NewScoreService(db)
	score, err := svc.GiveScoreByRule(&student, &rule, teacher.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, score.Amount)
	assert.Equal(t, "举手发言", score.Reason)
	require.NotNil(t, score.ScoreRuleID)
	assert.Equal(t, rule.ID, *score.ScoreRuleID)
	assert.Equal(t, 3, student.TotalScore)

	// 审计日志 + 宠物经验同步（与既有 GiveScore 同一实现）。
	var log models.ScoreLog
	require.NoError(t, db.Where("score_id = ?", score.ID).First(&log).Error)
	assert.Equal(t, 0, log.BalanceBefore)
	assert.Equal(t, 3, log.BalanceAfter)

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, 3, pet.Experience)
}

// nowMinusHour 一小时前（用于制造「已过期」的撤销记录）。
func nowMinusHour() time.Time { return time.Now().Add(-time.Hour) }
