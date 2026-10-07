// 管理端运维服务测试（一）：报表聚合、大屏登录日志、班级批量创建与教师分配、学年升级、LOGO 上传。
package services_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newAdminOps 构造运维服务（上传落盘指向临时目录，避免污染工作区）。
func newAdminOps(t *testing.T, db *gorm.DB) *services.AdminOps {
	t.Helper()
	ops := services.NewAdminOps(db)
	ops.UploadRoot = filepath.Join(t.TempDir(), "uploads")
	return ops
}

// makeClass 建班（带年级）。
func makeClass(t *testing.T, db *gorm.DB, schoolID uint, name, grade string) models.ClassRoom {
	t.Helper()
	class := models.ClassRoom{SchoolID: schoolID, Name: name, Grade: grade, Status: "active"}
	require.NoError(t, db.Create(&class).Error)
	return class
}

// makeStudent 建学生（指定状态）。
func makeStudent(t *testing.T, db *gorm.DB, classID uint, name, status string) models.Student {
	t.Helper()
	student := models.Student{ClassID: classID, Name: name, Status: status}
	require.NoError(t, db.Create(&student).Error)
	return student
}

// makeScore 建积分记录（可指定时间，用于本月/上月口径）。
func makeScore(t *testing.T, db *gorm.DB, classID, studentID uint, amount int, at time.Time) {
	t.Helper()
	require.NoError(t, db.Create(&models.Score{
		ClassID:   classID,
		StudentID: studentID,
		Amount:    amount,
		CreatedAt: at,
	}).Error)
}

// TestAdminOpsSchoolOverview 全校概览：班级/教师/学生数 + 本月/上月积分与环比。
func TestAdminOpsSchoolOverview(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	classA := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	classB := makeClass(t, db, school.ID, "一年级（2）班", "一年级")

	// 教师：2 启用 + 1 停用
	for _, tc := range []struct {
		username, name, status string
	}{
		{"t1", "张老师", "active"},
		{"t2", "李老师", "active"},
		{"t3", "王老师", "inactive"},
	} {
		require.NoError(t, db.Create(&models.User{
			SchoolID: school.ID, Role: "teacher", Username: tc.username,
			Name: tc.name, Status: tc.status,
		}).Error)
	}

	s1 := makeStudent(t, db, classA.ID, "甲", "active")
	s2 := makeStudent(t, db, classA.ID, "乙", "active")
	makeStudent(t, db, classA.ID, "丙", "graduated")
	s3 := makeStudent(t, db, classB.ID, "丁", "active")
	makeStudent(t, db, classB.ID, "戊", "active")

	now := util.Now()
	// 本月：+10（A）、-3（B）、+8（B）= 15；上月：+4（A）、+5（B）= 9
	makeScore(t, db, classA.ID, s1.ID, 10, now)
	makeScore(t, db, classB.ID, s2.ID, -3, now)
	makeScore(t, db, classB.ID, s3.ID, 8, now)
	lastMonth := now.AddDate(0, -1, 0)
	makeScore(t, db, classA.ID, s1.ID, 4, lastMonth)
	makeScore(t, db, classB.ID, s3.ID, 5, lastMonth)

	overview, err := ops.SchoolOverview(school.ID)
	require.NoError(t, err)

	assert.Equal(t, int64(2), overview.ClassCount)
	assert.Equal(t, int64(2), overview.TeacherCount, "只统计 status=active 的教师")
	assert.Equal(t, int64(5), overview.StudentCount, "统计班级下全部学生（含毕业，同 withCount）")
	assert.Equal(t, 15, overview.MonthlyScore)
	assert.Equal(t, 9, overview.LastMonthScore)
	assert.InDelta(t, 66.7, overview.ScoreTrendPercent, 0.05, "(15-9)/9*100 四舍五入到 1 位")
	assert.Equal(t, now.Format("2006-01"), overview.MonthLabel)
}

// TestAdminOpsReportsByGradeAndClass 按年级 / 按班级汇总的口径。
func TestAdminOpsReportsByGradeAndClass(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	teacher := seedTeacher(t, db, school.ID, "grade-teacher")

	classA := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", classA.ID).
		Update("teacher_id", teacher.ID).Error)
	classB := makeClass(t, db, school.ID, "一年级（2）班", "一年级")
	classC := makeClass(t, db, school.ID, "二年级（1）班", "二年级")
	classD := makeClass(t, db, school.ID, "未分年级班", "") // grade 为空 → 未分年级

	sA1 := makeStudent(t, db, classA.ID, "A1", "active")
	sA2 := makeStudent(t, db, classA.ID, "A2", "active")
	makeStudent(t, db, classA.ID, "A3", "active")
	sB1 := makeStudent(t, db, classB.ID, "B1", "active")
	makeStudent(t, db, classB.ID, "B2", "active")
	sC1 := makeStudent(t, db, classC.ID, "C1", "active")
	makeStudent(t, db, classD.ID, "D1", "active")

	now := util.Now()
	lastMonth := now.AddDate(0, -1, 0)
	makeScore(t, db, classA.ID, sA1.ID, 10, now)       // 本月 +10
	makeScore(t, db, classA.ID, sA2.ID, -3, lastMonth) // 上月 -3（只进总分，不进本月）
	makeScore(t, db, classB.ID, sB1.ID, 5, now)        // 本月 +5
	makeScore(t, db, classC.ID, sC1.ID, 8, now)        // 本月 +8

	gradeRows, err := ops.ReportsByGrade(school.ID)
	require.NoError(t, err)
	require.Len(t, gradeRows, 3, "一年级 / 二年级 / 未分年级")

	grade1 := gradeRows[0]
	assert.Equal(t, "一年级", grade1.Grade)
	assert.Equal(t, 2, grade1.ClassCount)
	assert.Equal(t, 5, grade1.StudentCount)
	assert.Equal(t, 12, grade1.TotalScore)
	assert.InDelta(t, 2.4, grade1.AvgScore, 0.001)

	grade2 := gradeRows[1]
	assert.Equal(t, "二年级", grade2.Grade)
	assert.Equal(t, 1, grade2.ClassCount)
	assert.Equal(t, 1, grade2.StudentCount)
	assert.Equal(t, 8, grade2.TotalScore)
	assert.InDelta(t, 8.0, grade2.AvgScore, 0.001)

	grade3 := gradeRows[2]
	assert.Equal(t, "未分年级", grade3.Grade, "空年级归入「未分年级」")
	assert.Equal(t, 0, grade3.TotalScore)
	assert.InDelta(t, 0.0, grade3.AvgScore, 0.001)

	classRows, err := ops.ReportsByClass(school.ID)
	require.NoError(t, err)
	require.Len(t, classRows, 4)

	rowA := classRows[0]
	assert.Equal(t, classA.ID, rowA.ClassID)
	assert.Equal(t, "一年级（1）班", rowA.ClassName)
	require.NotNil(t, rowA.TeacherName)
	assert.Equal(t, teacher.Name, *rowA.TeacherName)
	assert.Equal(t, 3, rowA.StudentCount)
	assert.Equal(t, 10, rowA.MonthlyScore, "本月积分只看当前月份（上月的 -3 不计入）")

	assert.Nil(t, classRows[1].TeacherName, "无班主任时为 null")
	assert.Equal(t, 5, classRows[1].MonthlyScore)
	assert.Equal(t, 8, classRows[2].MonthlyScore)
	assert.Equal(t, 0, classRows[3].MonthlyScore)
}

// TestAdminOpsReportsEmptySchool 空学校不报错，返回空切片而非 null。
func TestAdminOpsReportsEmptySchool(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	gradeRows, err := ops.ReportsByGrade(school.ID)
	require.NoError(t, err)
	assert.Empty(t, gradeRows)

	classRows, err := ops.ReportsByClass(school.ID)
	require.NoError(t, err)
	assert.Empty(t, classRows)
}

// TestAdminOpsDisplayLoginLogs 班级码登录日志：字段、筛选与分页。
func TestAdminOpsDisplayLoginLogs(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	classA := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	classB := makeClass(t, db, school.ID, "一年级（2）班", "一年级")

	// 固定取「昨天 12:00」作为基准（避免跨日/跨月边界导致日期筛选漂移）
	now := util.Now()
	base := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, util.Loc).AddDate(0, 0, -1)
	dayBefore := base.AddDate(0, 0, -1)

	require.NoError(t, db.Create(&models.DisplayLoginLog{
		ClassID: classA.ID, ClassCode: "LS11", IPAddress: "10.0.0.1",
		UserAgent: "UA-1", CreatedAt: base,
	}).Error)
	require.NoError(t, db.Create(&models.DisplayLoginLog{
		ClassID: classB.ID, ClassCode: "LS12", IPAddress: "192.168.1.7",
		UserAgent: "UA-2", CreatedAt: base.Add(time.Hour),
	}).Error)

	// 前天：52 条同班日志，用于验证分页（每页 50）与日期筛选隔离
	for i := 0; i < 52; i++ {
		require.NoError(t, db.Create(&models.DisplayLoginLog{
			ClassID: classA.ID, ClassCode: "LS11", IPAddress: "10.0.0.2",
			UserAgent: "UA-bulk", CreatedAt: dayBefore,
		}).Error)
	}

	res, err := ops.DisplayLoginLogs(school.ID, services.DisplayLoginLogFilter{})
	require.NoError(t, err)
	assert.Equal(t, 54, res.Meta.Total)
	assert.Equal(t, 1, res.Meta.CurrentPage)
	assert.Equal(t, 2, res.Meta.LastPage)
	assert.Len(t, res.Data, 50)

	// 首条为最新（created_at DESC）
	assert.Equal(t, "192.168.1.7", res.Data[0].IPAddress)
	assert.Equal(t, "LS12", res.Data[0].ClassCode)
	assert.Equal(t, "一年级（2）班", res.Data[0].ClassName)
	assert.Equal(t, "UA-2", res.Data[0].UserAgent)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, res.Data[0].LoginAt)

	// 按 IP 模糊筛选
	byIP, err := ops.DisplayLoginLogs(school.ID, services.DisplayLoginLogFilter{IP: "192.168"})
	require.NoError(t, err)
	require.Len(t, byIP.Data, 1)
	assert.Equal(t, "一年级（2）班", byIP.Data[0].ClassName)

	// 按班级筛选
	byClass, err := ops.DisplayLoginLogs(school.ID, services.DisplayLoginLogFilter{ClassID: classB.ID})
	require.NoError(t, err)
	assert.Equal(t, 1, byClass.Meta.Total)

	// 按日期筛选：昨天 2 条、前天 52 条（两种日期互不串味）
	byDate, err := ops.DisplayLoginLogs(school.ID, services.DisplayLoginLogFilter{
		Date: base.Format("2006-01-02"),
	})
	require.NoError(t, err)
	assert.Equal(t, 2, byDate.Meta.Total)

	byDayBefore, err := ops.DisplayLoginLogs(school.ID, services.DisplayLoginLogFilter{
		Date: dayBefore.Format("2006-01-02"),
	})
	require.NoError(t, err)
	assert.Equal(t, 52, byDayBefore.Meta.Total)

	// 他校日志不可见
	other := seedOtherSchool(t, db)
	otherClass := makeClass(t, db, other.ID, "他校一班", "一年级")
	require.NoError(t, db.Create(&models.DisplayLoginLog{
		ClassID: otherClass.ID, ClassCode: "XS11", IPAddress: "172.16.0.1", CreatedAt: base,
	}).Error)
	scoped, err := ops.DisplayLoginLogs(school.ID, services.DisplayLoginLogFilter{IP: "172.16"})
	require.NoError(t, err)
	assert.Equal(t, 0, scoped.Meta.Total)
}

// TestAdminOpsBatchCreateClasses 班级批量创建：序号自现有班级数递增 + 同名跳过。
func TestAdminOpsBatchCreateClasses(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	_, err := ops.BatchCreateClasses(school.ID, "", 2, "2026")
	assertAppStatus(t, err, 422, "年级必填")
	_, err = ops.BatchCreateClasses(school.ID, "一年级", 0, "")
	assertAppStatus(t, err, 422, "数量下限")
	_, err = ops.BatchCreateClasses(school.ID, "一年级", 21, "")
	assertAppStatus(t, err, 422, "数量上限")

	created, err := ops.BatchCreateClasses(school.ID, "一年级", 2, "2026")
	require.NoError(t, err)
	require.Len(t, created, 2)
	assert.Equal(t, "一年级（1）班", created[0].Name)
	assert.Equal(t, "一年级（2）班", created[1].Name)

	// 再建 2 个：序号从现有 2 个之后开始
	created2, err := ops.BatchCreateClasses(school.ID, "一年级", 2, "2026")
	require.NoError(t, err)
	require.Len(t, created2, 2)
	assert.Equal(t, "一年级（3）班", created2[0].Name)
	assert.Equal(t, "一年级（4）班", created2[1].Name)

	// 同名 active 班级不重复创建
	var count int64
	require.NoError(t, db.Model(&models.ClassRoom{}).
		Where("school_id = ? AND name = ?", school.ID, "一年级（1）班").Count(&count).Error)
	assert.Equal(t, int64(1), count)

	// 其他学校互不影响
	other := seedOtherSchool(t, db)
	otherCreated, err := ops.BatchCreateClasses(other.ID, "一年级", 1, "")
	require.NoError(t, err)
	require.Len(t, otherCreated, 1)
	assert.Equal(t, "一年级（1）班", otherCreated[0].Name)
}

// TestAdminOpsAssignAndRemoveTeacher 班级教师分配 / 移除对 class_rooms.teacher_id 的影响。
func TestAdminOpsAssignAndRemoveTeacher(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	class := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	teacher := seedTeacher(t, db, school.ID, "assign-teacher")
	other := seedOtherSchool(t, db)
	foreignClass := makeClass(t, db, other.ID, "他校一班", "一年级")

	// subject_teacher（默认）不改 teacher_id
	class2 := makeClass(t, db, school.ID, "一年级（2）班", "一年级")
	updated, err := ops.AssignClassTeacher(school.ID, class2.ID, teacher.ID, "")
	require.NoError(t, err)
	assert.Nil(t, updated.TeacherID, "非班主任角色不写 teacher_id")

	// head_teacher 写 teacher_id 并落 class_room_teachers
	updated, err = ops.AssignClassTeacher(school.ID, class.ID, teacher.ID, "head_teacher")
	require.NoError(t, err)
	require.NotNil(t, updated.TeacherID)
	assert.Equal(t, teacher.ID, *updated.TeacherID)

	var link models.ClassRoomTeacher
	require.NoError(t, db.Where("class_room_id = ? AND user_id = ?", class.ID, teacher.ID).First(&link).Error)
	assert.Equal(t, "head_teacher", link.Role)

	// 重复分配为 co_teacher：更新同一行角色，但 teacher_id 保持不变（Laravel 只改关联）
	_, err = ops.AssignClassTeacher(school.ID, class.ID, teacher.ID, "co_teacher")
	require.NoError(t, err)
	var linkCount int64
	require.NoError(t, db.Model(&models.ClassRoomTeacher{}).
		Where("class_room_id = ? AND user_id = ?", class.ID, teacher.ID).Count(&linkCount).Error)
	assert.Equal(t, int64(1), linkCount)
	require.NoError(t, db.First(&link, link.ID).Error)
	assert.Equal(t, "co_teacher", link.Role)

	// 守卫
	_, err = ops.AssignClassTeacher(school.ID, class.ID, teacher.ID, "not_a_role")
	assertAppStatus(t, err, 422, "角色白名单")
	_, err = ops.AssignClassTeacher(school.ID, foreignClass.ID, teacher.ID, "head_teacher")
	assertAppStatus(t, err, 404, "跨校班级不可见")
	_, err = ops.AssignClassTeacher(school.ID, class.ID, 99999, "head_teacher")
	assertAppStatus(t, err, 404, "教师不存在")

	// 移除：清关联 + 清 teacher_id
	require.NoError(t, ops.RemoveClassTeacher(school.ID, class.ID, teacher.ID))
	var reloaded models.ClassRoom
	require.NoError(t, db.First(&reloaded, class.ID).Error)
	assert.Nil(t, reloaded.TeacherID, "被移除的教师是班主任时 teacher_id 置空")
	var remaining int64
	require.NoError(t, db.Model(&models.ClassRoomTeacher{}).
		Where("class_room_id = ? AND user_id = ?", class.ID, teacher.ID).Count(&remaining).Error)
	assert.Equal(t, int64(0), remaining)
}

// TestAdminOpsGradeUpgrade 学年升级：预览 → 执行（不可逆）的全部语义。
func TestAdminOpsGradeUpgrade(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	grade1 := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	makeStudent(t, db, grade1.ID, "小一", "active")
	makeStudent(t, db, grade1.ID, "小二", "active")
	makeStudent(t, db, grade1.ID, "已毕业", "graduated")

	grade5 := makeClass(t, db, school.ID, "五年级（2）班", "五年级")
	makeStudent(t, db, grade5.ID, "小五", "active")

	grade6 := makeClass(t, db, school.ID, "六年级（1）班", "六年级")
	makeStudent(t, db, grade6.ID, "小六", "active")
	makeStudent(t, db, grade6.ID, "小六二", "active")

	// 非 active 班级不参与升级
	archived := models.ClassRoom{SchoolID: school.ID, Name: "二年级（9）班", Grade: "二年级", Status: "archived"}
	require.NoError(t, db.Create(&archived).Error)

	// 他校班级不受影响
	other := seedOtherSchool(t, db)
	otherClass := makeClass(t, db, other.ID, "一年级（1）班", "一年级")

	preview, err := ops.PreviewGradeUpgrade(school.ID)
	require.NoError(t, err)

	require.Len(t, preview.UpgradeClasses, 2)
	assert.Equal(t, grade1.ID, preview.UpgradeClasses[0].ClassID)
	assert.Equal(t, "一年级", preview.UpgradeClasses[0].OldGrade)
	assert.Equal(t, "二年级", preview.UpgradeClasses[0].NewGrade)
	assert.Equal(t, "二年级（1）班", preview.UpgradeClasses[0].NewName)
	assert.Equal(t, int64(3), preview.UpgradeClasses[0].StudentCount, "withCount 含非 active 学生")
	assert.Equal(t, "六年级（2）班", preview.UpgradeClasses[1].NewName)

	require.Len(t, preview.GraduateClasses, 1)
	assert.Equal(t, grade6.ID, preview.GraduateClasses[0].ClassID)
	assert.Equal(t, int64(2), preview.GraduateClasses[0].StudentCount)

	assert.Equal(t, 2, preview.Summary.UpgradeClassCount)
	assert.Equal(t, 1, preview.Summary.GraduateClassCount)
	assert.Equal(t, int64(4), preview.Summary.UpgradeStudentCount)
	assert.Equal(t, int64(2), preview.Summary.GraduateStudentCount)
	assert.Contains(t, preview.Summary.Note, "一年级新生需在升级后手动创建班级并导入")

	// 预览不写库
	var stillGrade1 models.ClassRoom
	require.NoError(t, db.First(&stillGrade1, grade1.ID).Error)
	assert.Equal(t, "一年级", stillGrade1.Grade)

	result, err := ops.ExecuteGradeUpgrade(school.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.UpgradedClasses)
	assert.Equal(t, 1, result.ArchivedClasses)
	assert.Equal(t, int64(2), result.GraduatedStudents)
	assert.Contains(t, result.Note, "六年级学生已标记为毕业")

	var upgraded models.ClassRoom
	require.NoError(t, db.First(&upgraded, grade1.ID).Error)
	assert.Equal(t, "二年级", upgraded.Grade)
	assert.Equal(t, "二年级（1）班", upgraded.Name)

	var graduatedClass models.ClassRoom
	require.NoError(t, db.First(&graduatedClass, grade6.ID).Error)
	assert.Equal(t, "archived", graduatedClass.Status)
	assert.Equal(t, "六年级", graduatedClass.Grade, "归档班级不再升级")

	var graduatedCount, activeCount int64
	require.NoError(t, db.Model(&models.Student{}).
		Where("class_id = ? AND status = ?", grade6.ID, "graduated").Count(&graduatedCount).Error)
	require.NoError(t, db.Model(&models.Student{}).
		Where("class_id = ? AND status = ?", grade6.ID, "active").Count(&activeCount).Error)
	assert.Equal(t, int64(2), graduatedCount)
	assert.Equal(t, int64(0), activeCount)

	// 他校班级未被改动
	var otherReload models.ClassRoom
	require.NoError(t, db.First(&otherReload, otherClass.ID).Error)
	assert.Equal(t, "一年级", otherReload.Grade)

	// 不可逆：再次执行会继续升级（归档班级因非 active 被排除）
	second, err := ops.ExecuteGradeUpgrade(school.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, second.UpgradedClasses, "已升级为二年级的班级会被再升一级")
	assert.Equal(t, 1, second.ArchivedClasses, "升级产生的六年级班级被归档")
	assert.Equal(t, int64(1), second.GraduatedStudents)
}

// TestAdminOpsSaveSchoolLogo LOGO 上传：落盘 + 写库 + 校验。
func TestAdminOpsSaveSchoolLogo(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	ops := newAdminOps(t, db)

	pngBytes := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0, 0, 0, 0x0D, 'I', 'H', 'D', 'R'}

	logoPath, err := ops.SaveSchoolLogo(school.ID, "my-logo.png", pngBytes)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(logoPath, services.DefaultUploadURLPrefix+"/schools/"), logoPath)
	assert.True(t, strings.HasSuffix(logoPath, ".png"))

	fileName := filepath.Base(logoPath)
	saved, err := os.ReadFile(filepath.Join(ops.UploadRoot, "schools", fileName))
	require.NoError(t, err)
	assert.Equal(t, pngBytes, saved, "Go 端保存原文件（不缩放/裁剪）")

	var reloaded models.School
	require.NoError(t, db.First(&reloaded, school.ID).Error)
	assert.Equal(t, logoPath, reloaded.LogoPath)

	// 类型不支持
	_, err = ops.SaveSchoolLogo(school.ID, "note.txt", pngBytes)
	assertAppStatus(t, err, 422, "扩展名白名单")
	// 伪装成 png 的文本
	_, err = ops.SaveSchoolLogo(school.ID, "fake.png", []byte("not an image at all"))
	assertAppStatus(t, err, 422, "内容不是图片")
	// 空文件
	_, err = ops.SaveSchoolLogo(school.ID, "empty.png", []byte{})
	assertAppStatus(t, err, 422, "空文件")
	// 超过 2MB
	big := make([]byte, 2048*1024+1)
	copy(big, pngBytes)
	_, err = ops.SaveSchoolLogo(school.ID, "big.png", big)
	assertAppStatus(t, err, 422, "超过 2MB")
	// 学校不存在
	_, err = ops.SaveSchoolLogo(99999, "logo.png", pngBytes)
	assertAppStatus(t, err, 404, "学校不存在")
}

// TestAdminOpsReportsOtherSchoolIsolated 报表只统计本校数据。
func TestAdminOpsReportsOtherSchoolIsolated(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	other := seedOtherSchool(t, db)
	ops := newAdminOps(t, db)

	otherClass := makeClass(t, db, other.ID, "他校一年级", "一年级")
	otherStudent := makeStudent(t, db, otherClass.ID, "他校生", "active")
	makeScore(t, db, otherClass.ID, otherStudent.ID, 100, util.Now())

	overview, err := ops.SchoolOverview(school.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), overview.ClassCount)
	assert.Equal(t, int64(0), overview.StudentCount)
	assert.Equal(t, 0, overview.MonthlyScore)

	rows, err := ops.ReportsByGrade(school.ID)
	require.NoError(t, err)
	assert.Empty(t, rows)

	classRows, err := ops.ReportsByClass(school.ID)
	require.NoError(t, err)
	assert.Empty(t, classRows)
	assert.NotEmpty(t, overview.MonthLabel)
}
