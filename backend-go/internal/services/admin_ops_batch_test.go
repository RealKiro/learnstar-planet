// 管理端运维服务测试（二）：批量账号、教师批量创建 / CSV 导入、学生批量导入与转班、系统运维。
package services_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// makeTeacher 建教师账号（可标记为 API 机器人）。
func makeTeacher(t *testing.T, db *gorm.DB, schoolID uint, username, name string, isBot bool) models.User {
	t.Helper()
	user := models.User{
		SchoolID: schoolID, Role: "teacher", Username: username, Name: name,
		Nickname: name, Status: "active", IsAPIBot: isBot,
	}
	require.NoError(t, db.Create(&user).Error)
	return user
}

// ============================================================
// 批量账号
// ============================================================

// TestAdminOpsBatchResetPassword 批量重置密码：只影响选中账号、返回计数、守卫。
func TestAdminOpsBatchResetPassword(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	other := seedOtherSchool(t, db)
	ops := newAdminOps(t, db)

	t1 := makeTeacher(t, db, school.ID, "reset-1", "张老师", false)
	t2 := makeTeacher(t, db, school.ID, "reset-2", "李老师", false)
	untouched := makeTeacher(t, db, school.ID, "reset-3", "王老师", false)
	foreign := makeTeacher(t, db, other.ID, "reset-foreign", "他校老师", false)
	admin := models.User{SchoolID: school.ID, Role: "school_admin", Username: "reset-admin", Name: "管理员", Status: "active"}
	require.NoError(t, db.Create(&admin).Error)

	count, err := ops.BatchResetPassword(school.ID, []uint{t1.ID, t2.ID, foreign.ID, admin.ID}, "newpass123")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "只重置本校 role=teacher 的选中账号")

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, t1.ID).Error)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(reloaded.PasswordHash), []byte("newpass123")))
	assert.False(t, reloaded.PasswordChanged)

	var untouchedReload models.User
	require.NoError(t, db.First(&untouchedReload, untouched.ID).Error)
	assert.Empty(t, untouchedReload.PasswordHash, "未选中的账号不受影响")

	// 默认密码
	count, err = ops.BatchResetPassword(school.ID, []uint{t1.ID}, "")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	require.NoError(t, db.First(&reloaded, t1.ID).Error)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(reloaded.PasswordHash),
		[]byte(services.DefaultTeacherPassword)))

	// 守卫
	_, err = ops.BatchResetPassword(school.ID, nil, "newpass123")
	assertAppStatus(t, err, 422, "ids 至少 1 条")
	_, err = ops.BatchResetPassword(school.ID, []uint{t1.ID}, "123")
	assertAppStatus(t, err, 422, "密码长度不足")
	_, err = ops.BatchResetPassword(school.ID, []uint{99999}, "newpass123")
	assertAppStatus(t, err, 404, "未找到所选账号")
}

// TestAdminOpsBatchDeleteAccounts 批量删除账号：机器人受保护、班级关联解除。
func TestAdminOpsBatchDeleteAccounts(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	bot := makeTeacher(t, db, school.ID, "bot-teacher", "API 机器人", true)
	t1 := makeTeacher(t, db, school.ID, "delete-1", "张老师", false)
	t2 := makeTeacher(t, db, school.ID, "delete-2", "李老师", false)

	class := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Update("teacher_id", t1.ID).Error)
	require.NoError(t, db.Create(&models.ClassRoomTeacher{
		ClassRoomID: class.ID, UserID: t1.ID, Role: "head_teacher",
	}).Error)

	deleted, protected, err := ops.BatchDeleteAccounts(school.ID, []uint{t1.ID, t2.ID, bot.ID})
	require.NoError(t, err)
	assert.Equal(t, int64(2), deleted)
	assert.Equal(t, 1, protected)

	var remaining int64
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", bot.ID).Count(&remaining).Error)
	assert.Equal(t, int64(1), remaining, "机器人账号不可删除")
	require.NoError(t, db.Model(&models.User{}).Where("id IN ?", []uint{t1.ID, t2.ID}).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)

	var reloaded models.ClassRoom
	require.NoError(t, db.First(&reloaded, class.ID).Error)
	assert.Nil(t, reloaded.TeacherID, "删除教师时解除班级关联")
	require.NoError(t, db.Model(&models.ClassRoomTeacher{}).Where("user_id = ?", t1.ID).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)

	// 只选机器人 → 403
	_, _, err = ops.BatchDeleteAccounts(school.ID, []uint{bot.ID})
	assertAppStatus(t, err, 403, "机器人账号不可删除")

	// 无命中 → 404
	_, _, err = ops.BatchDeleteAccounts(school.ID, []uint{99999})
	assertAppStatus(t, err, 404, "未找到所选账号")

	// 空 ids → 422
	_, _, err = ops.BatchDeleteAccounts(school.ID, nil)
	assertAppStatus(t, err, 422, "ids 至少 1 条")
}

// ============================================================
// 教师批量创建 / CSV 导入 / 模板
// ============================================================

// TestAdminOpsBatchCreateTeachers 教师批量创建：重名去重 + 班级分配 + 返回结构。
func TestAdminOpsBatchCreateTeachers(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	classA := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	classB := makeClass(t, db, school.ID, "一年级（2）班", "一年级")

	// 预置同名教师，验证 username/nickname 追加 _2
	require.NoError(t, db.Create(&models.User{
		SchoolID: school.ID, Role: "teacher", Username: "张老师", Nickname: "张老师",
		Name: "张老师", Status: "active",
	}).Error)

	created, err := ops.BatchCreateTeachers(school.ID, []services.TeacherInput{
		{
			Name: "张老师", Subject: "语文", GradeTeam: "三年级团队", Password: "abc123456",
			Assignments: []services.TeacherAssignmentInput{
				{ClassID: classA.ID, Role: "head_teacher"},
				{ClassID: classB.ID, Role: "subject_teacher", Subject: "语文"},
			},
		},
		{Name: "李老师", Password: "abc123456"},
	})
	require.NoError(t, err)
	require.Len(t, created, 2)

	assert.Equal(t, "张老师_2", created[0].Username, "同名用户名追加 _2")
	assert.Equal(t, "张老师_2", created[0].Nickname)
	assert.Equal(t, "abc123456", created[0].InitialPassword, "仅创建时返回初始密码")
	assert.Equal(t, "李老师", created[1].Username)
	assert.NotZero(t, created[1].ID)

	// head_teacher 写班级 teacher_id；subject_teacher 不写
	var reloadedA, reloadedB models.ClassRoom
	require.NoError(t, db.First(&reloadedA, classA.ID).Error)
	require.NotNil(t, reloadedA.TeacherID)
	assert.Equal(t, created[0].ID, *reloadedA.TeacherID)
	require.NoError(t, db.First(&reloadedB, classB.ID).Error)
	assert.Nil(t, reloadedB.TeacherID)

	var links []models.ClassRoomTeacher
	require.NoError(t, db.Where("user_id = ?", created[0].ID).Order("class_room_id ASC").Find(&links).Error)
	require.Len(t, links, 2)
	assert.Equal(t, "head_teacher", links[0].Role)
	assert.Equal(t, "subject_teacher", links[1].Role)
	assert.Equal(t, "语文", links[1].Subject)

	// 校验守卫
	_, err = ops.BatchCreateTeachers(school.ID, nil)
	assertAppStatus(t, err, 422, "teachers 至少 1 条")
	_, err = ops.BatchCreateTeachers(school.ID, []services.TeacherInput{{Name: "", Password: "abc123456"}})
	assertAppStatus(t, err, 422, "姓名必填")
	_, err = ops.BatchCreateTeachers(school.ID, []services.TeacherInput{{Name: "无密码"}})
	assertAppStatus(t, err, 422, "密码必填")
	_, err = ops.BatchCreateTeachers(school.ID, []services.TeacherInput{{Name: "短密码", Password: "123"}})
	assertAppStatus(t, err, 422, "密码过短")
	_, err = ops.BatchCreateTeachers(school.ID, []services.TeacherInput{{
		Name: "坏角色", Password: "abc123456",
		Assignments: []services.TeacherAssignmentInput{{ClassID: classA.ID, Role: "boss"}},
	}})
	assertAppStatus(t, err, 422, "角色白名单")
	_, err = ops.BatchCreateTeachers(school.ID, []services.TeacherInput{{
		Name: "坏班级", Password: "abc123456",
		Assignments: []services.TeacherAssignmentInput{{ClassID: 99999, Role: "head_teacher"}},
	}})
	assertAppStatus(t, err, 422, "班级不存在")
}

// TestAdminOpsTeacherTemplateCSV 模板内容与 Laravel 逐字一致（含 BOM）。
func TestAdminOpsTeacherTemplateCSV(t *testing.T) {
	db := setupDB(t)
	ops := newAdminOps(t, db)

	csv := ops.TeacherTemplateCSV()
	assert.True(t, strings.HasPrefix(csv, "\ufeff"), "带 UTF-8 BOM")
	assert.Equal(t,
		"\ufeff姓名,年级团队,科目,密码,手机号\n"+
			"name,grade_team,subject,password,phone\n"+
			"张老师,三年级团队,语文,star123456,13800138000\n"+
			"李老师,三年级团队,数学,,\n", csv)
}

// TestAdminOpsImportTeachersCSV 教师 CSV 导入：预览（dry_run 默认）与正式创建。
func TestAdminOpsImportTeachersCSV(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	csvContent := []byte("\ufeff姓名,年级团队,科目,密码,手机号\n" +
		"张老师,三年级团队,语文,pass123456,13800138000\n" +
		"李老师,三年级团队,数学,,\n" +
		",,,,\n") // 姓名为空的行被跳过

	preview, err := ops.ImportTeachersCSV(school.ID, "teachers.csv", csvContent, true)
	require.NoError(t, err)
	assert.Equal(t, 2, preview.Total)
	assert.Equal(t, "预览模式：共 2 条数据", preview.Message)
	require.Len(t, preview.Preview, 2)
	assert.Equal(t, "张老师", preview.Preview[0].Name)
	assert.Equal(t, "三年级团队", preview.Preview[0].GradeTeam)
	assert.Equal(t, "语文", preview.Preview[0].Subject)
	assert.Equal(t, "pass123456", preview.Preview[0].Password)
	assert.Equal(t, "13800138000", preview.Preview[0].Phone)
	assert.Empty(t, preview.Preview[1].Password)

	var count int64
	require.NoError(t, db.Model(&models.User{}).Where("school_id = ? AND role = ?", school.ID, "teacher").
		Count(&count).Error)
	assert.Equal(t, int64(0), count, "预览模式不建号")

	result, err := ops.ImportTeachersCSV(school.ID, "teachers.csv", csvContent, false)
	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, "已导入 2 名教师", result.Message)
	require.Len(t, result.Created, 2)
	assert.Equal(t, services.DefaultTeacherPassword, result.Created[1].InitialPassword, "空密码用默认 ls123456")

	require.NoError(t, db.Model(&models.User{}).Where("school_id = ? AND role = ?", school.ID, "teacher").
		Count(&count).Error)
	assert.Equal(t, int64(2), count)

	// 英文表头同样可用
	enCSV := []byte("name,grade_team,subject,password,phone\n王老师,四年级团队,英语,pwd123456,13900139000\n")
	enResult, err := ops.ImportTeachersCSV(school.ID, "teachers.csv", enCSV, true)
	require.NoError(t, err)
	require.Len(t, enResult.Preview, 1)
	assert.Equal(t, "王老师", enResult.Preview[0].Name)
	assert.Equal(t, "四年级团队", enResult.Preview[0].GradeTeam)

	// 制表符分隔
	tabCSV := []byte("姓名\t年级团队\t科目\t密码\t手机号\n赵老师\t五年级团队\t科学\tpwd123456\t\n")
	tabResult, err := ops.ImportTeachersCSV(school.ID, "teachers.txt", tabCSV, true)
	require.NoError(t, err)
	require.Len(t, tabResult.Preview, 1)
	assert.Equal(t, "赵老师", tabResult.Preview[0].Name)

	// Excel 不支持
	_, err = ops.ImportTeachersCSV(school.ID, "teachers.xlsx", csvContent, true)
	assertAppStatus(t, err, 422, "xlsx 不支持")

	// 非 UTF-8 内容
	_, err = ops.ImportTeachersCSV(school.ID, "teachers.csv", []byte{0xFF, 0xFE, 0x00, 0x41}, true)
	assertAppStatus(t, err, 422, "非 UTF-8")

	// 空文件 → 0 条预览（同 Laravel 解析空文件）
	empty, err := ops.ImportTeachersCSV(school.ID, "teachers.csv", nil, true)
	require.NoError(t, err)
	assert.Equal(t, 0, empty.Total)
}

// ============================================================
// 学生批量导入 / 删除 / 转班
// ============================================================

// TestAdminOpsImportStudents 学生导入：成功、行内错误、跨班学号拦截、同名提醒。
func TestAdminOpsImportStudents(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	other := seedOtherSchool(t, db)
	ops := newAdminOps(t, db)

	classA := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	classB := makeClass(t, db, school.ID, "一年级（2）班", "一年级")
	// 他校同号学生不应触发跨班拦截
	otherClass := makeClass(t, db, other.ID, "他校一年级", "一年级")
	require.NoError(t, db.Create(&models.Student{
		ClassID: otherClass.ID, Name: "他校生", StudentNo: "9001", Status: "active",
	}).Error)

	// 已存在：classA 的 1001、classB 的小明（用于同名提醒）
	require.NoError(t, db.Create(&models.Student{
		ClassID: classA.ID, Name: "占位", StudentNo: "1001", Status: "active",
	}).Error)
	require.NoError(t, db.Create(&models.Student{
		ClassID: classB.ID, Name: "小明", Status: "active",
	}).Error)

	result, err := ops.ImportStudents(school.ID, []services.StudentImportRow{
		{Name: "小明", ClassName: "一年级（1）班", Gender: "男生", StudentNo: "9001"}, // 他校同号不算跨班冲突 → 成功 + 同名提醒
		{Name: "小红", ClassName: "一年级（2）班", Gender: "女生", StudentNo: "2002"}, // 成功（同名提醒：无）
		{Name: "小刚", ClassName: "一年级（1）班", Gender: "男", StudentNo: "1001"},  // 同班同学号 → 错误
		{Name: "小丽", ClassName: "三年级（9）班", Gender: "女", StudentNo: "3003"},  // 班级不存在 → 错误
		{Name: "小美", ClassName: "一年级（2）班", Gender: "未知", StudentNo: "9002"}, // 本校无该学号 → 成功
	})
	require.NoError(t, err)

	assert.Len(t, result.Created, 3)
	assert.Len(t, result.Errors, 2)
	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Errors[0], "学号「1001」在班级「一年级（1）班」已存在")
	assert.Contains(t, result.Errors[1], "班级「三年级（9）班」不存在")
	assert.Contains(t, result.Warnings[0], "新导入的「小明」与「一年级（2）班」现有学生同名")

	assert.Equal(t, "男", result.Created[0].Gender, "男生 归一化为 男")
	assert.Equal(t, "女", result.Created[1].Gender)
	assert.Equal(t, "未知", result.Created[2].Gender)

	// 成功行自动分配默认宠物（同 Laravel assignDefaultPet）
	var pets int64
	require.NoError(t, db.Model(&models.Pet{}).
		Where("student_id = ?", result.Created[0].ID).Count(&pets).Error)
	assert.Equal(t, int64(1), pets)
	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", result.Created[0].ID).First(&pet).Error)
	assert.Equal(t, "小明的萌宠", pet.Name)
	assert.Equal(t, 1, pet.Level)
	assert.Equal(t, 80, pet.Mood)
	assert.NotEmpty(t, pet.Species)

	// 本校跨班学号拦截
	result2, err := ops.ImportStudents(school.ID, []services.StudentImportRow{
		{Name: "小强", ClassName: "一年级（2）班", Gender: "男", StudentNo: "1001"},
	})
	require.NoError(t, err)
	assert.Empty(t, result2.Created)
	require.Len(t, result2.Errors, 1)
	assert.Contains(t, result2.Errors[0], "已存在于「一年级（1）班」")
	assert.Contains(t, result2.Errors[0], "批量转班")

	// 空数组 → 422
	_, err = ops.ImportStudents(school.ID, nil)
	assertAppStatus(t, err, 422, "students 至少 1 条")
}

// TestAdminOpsParseStudentCSV 学生 CSV 解析：中文/英文表头 + 无表头兜底。
func TestAdminOpsParseStudentCSV(t *testing.T) {
	db := setupDB(t)
	ops := newAdminOps(t, db)

	rows, err := ops.ParseStudentCSV(strings.NewReader(
		"\ufeff姓名,班级,性别,学号\n" +
			"小明,一年级（1）班,男,1001\n" +
			"\n" +
			"小红,一年级（2）班,女,\n"))
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, services.StudentImportRow{Name: "小明", ClassName: "一年级（1）班", Gender: "男", StudentNo: "1001"}, rows[0])
	assert.Equal(t, services.StudentImportRow{Name: "小红", ClassName: "一年级（2）班", Gender: "女"}, rows[1])

	enRows, err := ops.ParseStudentCSV(strings.NewReader(
		"name,class_name,gender,student_no\n小明,一年级（1）班,男生,1001\n"))
	require.NoError(t, err)
	require.Len(t, enRows, 1)
	assert.Equal(t, "小明", enRows[0].Name)
	assert.Equal(t, "男生", enRows[0].Gender)

	// 无表头：按固定列序解析
	positional, err := ops.ParseStudentCSV(strings.NewReader("小明,一年级（1）班,男,1001\n"))
	require.NoError(t, err)
	require.Len(t, positional, 1)
	assert.Equal(t, "小明", positional[0].Name)
	assert.Equal(t, "一年级（1）班", positional[0].ClassName)
	assert.Equal(t, "1001", positional[0].StudentNo)
}

// TestAdminOpsBatchDeleteAndMoveStudents 学生批量删除与批量转班。
func TestAdminOpsBatchDeleteAndMoveStudents(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	other := seedOtherSchool(t, db)
	ops := newAdminOps(t, db)

	classA := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	classB := makeClass(t, db, school.ID, "一年级（2）班", "一年级")
	otherClass := makeClass(t, db, other.ID, "他校一班", "一年级")

	s1 := makeStudent(t, db, classA.ID, "学生一", "active")
	s2 := makeStudent(t, db, classA.ID, "学生二", "active")
	s3 := makeStudent(t, db, classB.ID, "学生三", "active")
	foreign := makeStudent(t, db, otherClass.ID, "他校生", "active")

	// 转班
	count, targetName, err := ops.BatchMoveStudents(school.ID, []uint{s1.ID, s2.ID, foreign.ID}, classB.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "他校学生不被移动")
	assert.Equal(t, "一年级（2）班", targetName)
	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, s1.ID).Error)
	assert.Equal(t, classB.ID, reloaded.ClassID)
	var foreignReload models.Student
	require.NoError(t, db.First(&foreignReload, foreign.ID).Error)
	assert.Equal(t, otherClass.ID, foreignReload.ClassID)

	// 目标班级跨校 → 404
	_, _, err = ops.BatchMoveStudents(school.ID, []uint{s1.ID}, otherClass.ID)
	assertAppStatus(t, err, 404, "跨校目标班级")

	// 空参数 → 422
	_, _, err = ops.BatchMoveStudents(school.ID, nil, classB.ID)
	assertAppStatus(t, err, 422, "ids 至少 1 条")

	// 删除（含宠物/积分级联）
	require.NoError(t, db.Create(&models.Pet{StudentID: s3.ID, ClassID: classB.ID, Name: "宠物", Species: "zhulong"}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: s3.ID, ClassID: classB.ID, Amount: 5, CreatedAt: util.Now()}).Error)
	require.NoError(t, db.Create(&models.ScoreLog{StudentID: s3.ID, ScoreID: 1, BalanceBefore: 0, BalanceAfter: 5}).Error)

	deleted, err := ops.BatchDeleteStudents(school.ID, []uint{s3.ID, foreign.ID})
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted, "他校学生不被删除")

	var remaining int64
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", s3.ID).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)
	require.NoError(t, db.Model(&models.Pet{}).Where("student_id = ?", s3.ID).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)
	require.NoError(t, db.Model(&models.Score{}).Where("student_id = ?", s3.ID).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)
	require.NoError(t, db.Model(&models.ScoreLog{}).Where("student_id = ?", s3.ID).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)

	// 无命中返回 0（不报错）
	deleted, err = ops.BatchDeleteStudents(school.ID, []uint{99999})
	require.NoError(t, err)
	assert.Equal(t, int64(0), deleted)

	_, err = ops.BatchDeleteStudents(school.ID, nil)
	assertAppStatus(t, err, 422, "ids 至少 1 条")
}

// ============================================================
// 系统运维
// ============================================================

// TestAdminOpsDiagnose 诊断：结构、关键状态与 has_issues。
func TestAdminOpsDiagnose(t *testing.T) {
	db := setupDB(t)
	ops := newAdminOps(t, db)

	result, err := ops.Diagnose()
	require.NoError(t, err)
	require.NotEmpty(t, result.Items)
	assert.False(t, result.HasIssues)
	assert.Equal(t, "系统状态正常", result.Message)

	statusOf := map[string]string{}
	for _, item := range result.Items {
		statusOf[item.Item] = item.Status
	}
	assert.Equal(t, "ok", statusOf["users 表.nickname 字段"])
	assert.Equal(t, "ok", statusOf["users 表.subject 字段"])
	assert.Equal(t, "ok", statusOf["users 表.grade_team 字段"])
	assert.Equal(t, "ok", statusOf["class_rooms 表.display_code 字段"])
	assert.Equal(t, "ok", statusOf["class_rooms 表.display_code_updated_at 字段"])
	assert.Equal(t, "ok", statusOf["class_room_teachers 表"])
	assert.Equal(t, "ok", statusOf["class_room_teachers.subject 字段"])
	// 本批起 third_party_bindings 表已随 AutoMigrate 建表（支撑 /auth/bindings），故状态从 skipped 变 ok。
	assert.Equal(t, "ok", statusOf["third_party_bindings 表"], "本批已建 third_party_bindings 表")
	assert.Equal(t, "ok", statusOf["教师账号密码状态"])
}

// TestAdminOpsStatus 系统状态：版本字段 + 表清单（无 migrations 记录表）。
func TestAdminOpsStatus(t *testing.T) {
	db := setupDB(t)
	ops := newAdminOps(t, db)

	data, err := ops.Status()
	require.NoError(t, err)

	assert.Contains(t, data.Version.Go, "go")
	assert.Equal(t, "backend-go", data.Version.Backend)
	assert.NotEmpty(t, data.Version.AppEnv)
	assert.Equal(t, "Asia/Shanghai", data.Version.Timezone)
	assert.NotEmpty(t, data.Version.DBDriver)

	assert.Empty(t, data.Migrations)
	assert.Equal(t, 0, data.MigrationCount)
	assert.Contains(t, data.Tables, "students")
	assert.Contains(t, data.Tables, "class_room_teachers")
	assert.NotEmpty(t, data.Note)
}

// TestAdminOpsSystemLogs 日志：未配置路径 / 文件不存在 / 正常读取与等级过滤。
func TestAdminOpsSystemLogs(t *testing.T) {
	db := setupDB(t)
	ops := newAdminOps(t, db)

	// 未配置 LOG_FILE → 空列表 + 说明，不伪造日志
	empty := ops.SystemLogs("", 200, "")
	assert.Empty(t, empty.Entries)
	assert.False(t, empty.Exists)
	assert.Contains(t, empty.Message, "LOG_FILE")

	// 文件不存在
	missing := ops.SystemLogs(filepath.Join(t.TempDir(), "nope.log"), 200, "")
	assert.Empty(t, missing.Entries)
	assert.False(t, missing.Exists)
	assert.Contains(t, missing.Message, "日志文件不存在")

	// 正常读取：末 N 行、等级识别、时间逆序、等级过滤
	logPath := filepath.Join(t.TempDir(), "app.log")
	content := strings.Join([]string{
		"[2026-01-01 10:00:00] production.INFO: 第一条",
		"[2026-01-01 10:00:01] production.ERROR: 第二条",
		"",
		"[2026-01-01 10:00:02] production.DEBUG: 第三条",
		"[2026-01-01 10:00:03] 无等级前缀的行",
	}, "\n")
	require.NoError(t, os.WriteFile(logPath, []byte(content), 0o644))

	all := ops.SystemLogs(logPath, 200, "")
	require.True(t, all.Exists)
	assert.Equal(t, 4, all.Total)
	assert.NotZero(t, all.Size)
	assert.NotEmpty(t, all.ModifiedAt)
	assert.Equal(t, logPath, all.Path)
	// 逆序：最后一行最先出现，无前缀行默认 INFO
	assert.Equal(t, "INFO", all.Entries[0].Level)
	assert.Contains(t, all.Entries[0].Line, "无等级前缀的行")
	assert.Equal(t, "DEBUG", all.Entries[1].Level)
	assert.Equal(t, "ERROR", all.Entries[2].Level)
	assert.Equal(t, "INFO", all.Entries[3].Level)

	// 等级过滤（大小写不敏感）
	errorsOnly := ops.SystemLogs(logPath, 200, "error")
	require.Len(t, errorsOnly.Entries, 1)
	assert.Equal(t, "ERROR", errorsOnly.Entries[0].Level)

	// 末 N 行
	lastTwo := ops.SystemLogs(logPath, 2, "")
	require.Len(t, lastTwo.Entries, 2)
	assert.Equal(t, "DEBUG", lastTwo.Entries[1].Level, "末 2 行 = DEBUG 行 + 无前缀行，逆序后前者在后")
	assert.Contains(t, lastTwo.Entries[0].Line, "无等级前缀的行")

	// 默认行数（lines <= 0 → 200）
	defaultLines := ops.SystemLogs(logPath, 0, "")
	assert.Equal(t, 4, defaultLines.Total)
}

// TestAdminOpsRepairIdempotent 系统修复幂等：重复执行均成功且不改变数据。
func TestAdminOpsRepairIdempotent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	teacher := makeTeacher(t, db, school.ID, "repair-teacher", "张老师", false)
	class := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	student := makeStudent(t, db, class.ID, "小明", "active")

	first, err := ops.Repair()
	require.NoError(t, err)
	assert.Contains(t, first.Message, "数据库迁移已完成")
	// makeTeacher 直接建库、未写明文密码 → 计入「待重置」（同 Laravel systemRepair 口径）。
	assert.Contains(t, first.Message,
		"检测到 1 个教师账号缺少明文密码记录（第三方自动注册历史问题），请在「教师管理 → 密码」中逐个重置为默认密码")
	assert.NotEmpty(t, first.Output)
	assert.Equal(t, int64(1), first.PendingPasswordResets)

	second, err := ops.Repair()
	require.NoError(t, err)
	assert.Equal(t, first, second, "幂等：重复执行结果一致")

	// 补齐明文密码后：提示消失、计数归零。
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", teacher.ID).
		Update("plain_password", "ls123456").Error)
	afterFix, err := ops.Repair()
	require.NoError(t, err)
	assert.Equal(t, "数据库迁移已完成", afterFix.Message)
	assert.Equal(t, int64(0), afterFix.PendingPasswordResets)

	var userCount, classCount, studentCount int64
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", teacher.ID).Count(&userCount).Error)
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).Count(&classCount).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).Count(&studentCount).Error)
	assert.Equal(t, int64(1), userCount)
	assert.Equal(t, int64(1), classCount)
	assert.Equal(t, int64(1), studentCount)
}

// TestAdminOpsMigrationSource 守卫：诊断/状态/修复三者读的是同一套 schema。
func TestAdminOpsMigrationSource(t *testing.T) {
	db := setupDB(t)
	ops := newAdminOps(t, db)

	status, err := ops.Status()
	require.NoError(t, err)
	diag, err := ops.Diagnose()
	require.NoError(t, err)

	// 诊断里 checks 的表必须真实存在于 AutoMigrate 结果中
	for _, item := range diag.Items {
		if strings.HasSuffix(item.Item, " 表") && item.Status == "ok" {
			table := strings.TrimSuffix(item.Item, " 表")
			assert.Contains(t, status.Tables, table, fmt.Sprintf("诊断项 %q 声称表存在", item.Item))
		}
	}
}
