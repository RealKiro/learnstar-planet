// 班级大屏服务测试：token 生命周期、班级码登录与暴力破解防护、只读接口字段、班级码管理。
package services_test

import (
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

// displayFixture 大屏测试夹具：学校 + 一个确定班级码的班级 + 教师。
type displayFixture struct {
	School  models.School
	Class   models.ClassRoom
	Teacher models.User
	Svc     *services.DisplayService
}

// newDisplayFixture 建好学校 / 班级（display_code 显式写入，便于登录按码解析）。
func newDisplayFixture(t *testing.T, db *gorm.DB) displayFixture {
	t.Helper()
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher-display")

	class := models.ClassRoom{
		SchoolID: school.ID, Grade: "一年级", Name: "一年级（1）班",
		TeacherID: &teacher.ID, Status: "active", DisplayCode: "LS11",
	}
	require.NoError(t, db.Create(&class).Error)

	return displayFixture{School: school, Class: class, Teacher: teacher, Svc: services.NewDisplayService(db)}
}

// appErrorOf 要求 err 是带指定状态码的业务错误，并返回其文案。
func appErrorOf(t *testing.T, err error, status int) string {
	t.Helper()
	require.Error(t, err)
	ae, ok := services.AsAppError(err)
	require.True(t, ok, "应返回业务错误，实际: %v", err)
	require.Equal(t, status, ae.Status, "错误状态码不符，文案: %s", ae.Message)
	return ae.Message
}

// 登录成功 → 签发 disp_ token → 校验通过 → 过期 / 未知前缀失效。
func TestDisplayLoginAndTokenLifecycle(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	seedStudent(t, db, f.Class.ID)

	// 班级码大小写与首尾空格都会被规范化。
	result, err := f.Svc.Login("  ls11 ", "10.0.0.7", "unit-test-agent")
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result.Token, "disp_"), "token=%q", result.Token)
	assert.Len(t, result.Token, len("disp_")+32)
	assert.Equal(t, 86400, result.ExpiresIn)
	assert.Equal(t, f.Class.ID, result.ClassInfo.ID)
	assert.Equal(t, "一年级（1）班", result.ClassInfo.Name)
	assert.Equal(t, "一年级", result.ClassInfo.Grade)
	assert.Equal(t, int64(1), result.ClassInfo.StudentCount)

	// 登录日志落库（含 IP / UA）。
	var log models.DisplayLoginLog
	require.NoError(t, db.Where("class_id = ?", f.Class.ID).First(&log).Error)
	assert.Equal(t, "LS11", log.ClassCode)
	assert.Equal(t, "10.0.0.7", log.IPAddress)
	assert.Equal(t, "unit-test-agent", log.UserAgent)

	// disp_ token 可校验并解析出班级 ID。
	classID, err := f.Svc.ValidateToken(result.Token)
	require.NoError(t, err)
	assert.Equal(t, f.Class.ID, classID)

	// class_ 前缀 token 走同一张表（语义 = 取 class_id）。
	classToken := models.DisplayToken{
		Token: "class_manual-token", ClassID: f.Class.ID, ExpiresAt: time.Now().Add(time.Hour),
	}
	require.NoError(t, db.Create(&classToken).Error)
	classID, err = f.Svc.ValidateToken("class_manual-token")
	require.NoError(t, err)
	assert.Equal(t, f.Class.ID, classID)

	// 未知 token / 非法前缀 → 401「Token 无效或已过期」。
	for _, bad := range []string{"disp_unknown", "not-a-token", ""} {
		_, err := f.Svc.ValidateToken(bad)
		assert.Equal(t, "Token 无效或已过期", appErrorOf(t, err, 401), "token=%q", bad)
	}

	// 过期即无效。
	require.NoError(t, db.Model(&models.DisplayToken{}).Where("token = ?", result.Token).
		Update("expires_at", time.Now().Add(-time.Minute)).Error)
	_, err = f.Svc.ValidateToken(result.Token)
	assert.Equal(t, "Token 无效或已过期", appErrorOf(t, err, 401))
}

// 无效班级码 → 404；同码连续 5 次失败后锁定 15 分钟 → 429。
func TestDisplayLoginInvalidCodeAndLockout(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	for i := 0; i < services.DisplayLoginMaxAttempts; i++ {
		_, err := f.Svc.Login("LS99", "10.0.0.1", "ua")
		assert.Equal(t, "班级码无效，请检查后重试", appErrorOf(t, err, 404), "第 %d 次", i+1)
	}

	// 第 6 次触发 429（计数在阈值检查时已满）。
	_, err := f.Svc.Login("LS99", "10.0.0.1", "ua")
	assert.Equal(t, "尝试次数过多，请 15 分钟后再试", appErrorOf(t, err, 429))

	// 锁定按班级码隔离：另一个有效码仍可正常登录。
	result, err := f.Svc.Login("LS11", "10.0.0.1", "ua")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Token)
}

// 成功登录后清空该码的失败计数（Laravel Cache::forget 语义）。
func TestDisplayLoginResetsAttempts(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	// 先失败一次。
	_, err := f.Svc.Login("LS12", "10.0.0.1", "ua")
	assert.Equal(t, "班级码无效，请检查后重试", appErrorOf(t, err, 404))

	// 让 LS12 变为有效码并成功登录 → 计数清空。
	second := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active", DisplayCode: "LS12"}
	require.NoError(t, db.Create(&second).Error)
	_, err = f.Svc.Login("LS12", "10.0.0.1", "ua")
	require.NoError(t, err)

	// 再让 LS12 失效。若成功登录没清空计数，5 次里有一次会命中 429；
	// 这里 5 次全为 404，即证明计数已被成功登录清空。
	require.NoError(t, db.Delete(&models.ClassRoom{}, second.ID).Error)
	for i := 0; i < services.DisplayLoginMaxAttempts; i++ {
		_, err := f.Svc.Login("LS12", "10.0.0.1", "ua")
		assert.Equal(t, "班级码无效，请检查后重试", appErrorOf(t, err, 404), "第 %d 次", i+1)
	}
	_, err = f.Svc.Login("LS12", "10.0.0.1", "ua")
	assert.Equal(t, "尝试次数过多，请 15 分钟后再试", appErrorOf(t, err, 429))
}

// initialData：学生按学号数字序 + 宠物卡片 + 广播 + 最近 4 小时积分。
func TestDisplayInitialData(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	// 学号故意乱序：验证「2 < 10」的数字序（等价 Laravel 的 CAST(student_no AS UNSIGNED)）。
	s2 := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", TotalScore: 10, Status: "active"}
	s10 := models.Student{ClassID: f.Class.ID, Name: "丙", StudentNo: "10", TotalScore: 0, Status: "active"}
	s1 := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 25, Status: "active"}
	inactive := models.Student{ClassID: f.Class.ID, Name: "丁", StudentNo: "3", Status: "inactive"}
	require.NoError(t, db.Create(&s2).Error)
	require.NoError(t, db.Create(&s10).Error)
	require.NoError(t, db.Create(&s1).Error)
	require.NoError(t, db.Create(&inactive).Error)

	// 只有一个学生已有宠物；其余会被自动分配。
	pet := models.Pet{StudentID: s1.ID, ClassID: f.Class.ID, Name: "甲的萌宠", Species: "zhulong", Level: 5, Experience: 3, Mood: 70}
	require.NoError(t, db.Create(&pet).Error)

	// 广播：生效中的（sent/pending）会被返回，其他状态不会。
	require.NoError(t, db.Create(&models.Broadcast{ClassID: f.Class.ID, SchoolID: f.School.ID, Content: "上课啦", Type: "banner", DisplaySeconds: 15, VoiceEnabled: true, Status: "sent"}).Error)
	require.NoError(t, db.Create(&models.Broadcast{ClassID: f.Class.ID, SchoolID: f.School.ID, Content: "已撤销", Type: "banner", Status: "cancelled"}).Error)

	// 积分：4 小时内 2 条 + 5 小时前 1 条（应被排除）。
	now := util.Now()
	require.NoError(t, db.Create(&models.Score{StudentID: s1.ID, ClassID: f.Class.ID, Amount: 3, Reason: "举手发言", CreatedAt: now.Add(-10 * time.Minute)}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: s2.ID, ClassID: f.Class.ID, Amount: -2, Reason: "上课走神", CreatedAt: now.Add(-2 * time.Hour)}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: s2.ID, ClassID: f.Class.ID, Amount: 5, Reason: "作业优秀", CreatedAt: now.Add(-5 * time.Hour)}).Error)

	data, err := f.Svc.InitialData(f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "一年级（1）班", data.ClassName)
	assert.Equal(t, "一年级", data.Grade)
	assert.Equal(t, 3, data.StudentCount)
	_, err = time.Parse(time.RFC3339, data.ServerTime)
	assert.NoError(t, err, "server_time 应为 RFC3339")

	require.Len(t, data.Pets, 3)
	assert.Equal(t, "甲", data.Pets[0].StudentName)
	assert.Equal(t, "乙", data.Pets[1].StudentName)
	assert.Equal(t, "丙", data.Pets[2].StudentName, "学号 10 应排在 2 之后")

	// 已有宠物：字段取自宠物实体。
	assert.True(t, data.Pets[0].HasPet)
	assert.Equal(t, 5, data.Pets[0].Level)
	assert.Equal(t, 3, data.Pets[0].Experience)
	assert.Equal(t, 70, data.Pets[0].Mood)
	assert.Equal(t, 25, data.Pets[0].TotalScore)
	// Lv.5 → 成长期（12 级制的阶段划分同 Laravel Pet::currentStage）。
	assert.Equal(t, "成长期", data.Pets[0].StageName)
	assert.Equal(t, "🌱", data.Pets[0].Emoji)
	assert.Equal(t, 60, data.Pets[0].ExpMax)
	assert.Nil(t, data.Pets[0].Image)
	require.NotNil(t, data.Pets[0].PetName)
	assert.Equal(t, "甲的萌宠", *data.Pets[0].PetName)

	// 无宠物学生：自动建宠（name = 姓名 + 的萌宠，level 1）并返回。
	var created models.Pet
	require.NoError(t, db.Where("student_id = ?", s10.ID).First(&created).Error)
	assert.Equal(t, "丙的萌宠", created.Name)
	assert.Equal(t, 1, created.Level)
	assert.Equal(t, 80, created.Mood)
	assert.NotEmpty(t, created.Species)

	// 广播：仅生效中的，字段名与实现一致。
	require.Len(t, data.Broadcasts, 1)
	assert.Equal(t, "上课啦", data.Broadcasts[0].Content)
	assert.Equal(t, 15, data.Broadcasts[0].DisplaySeconds)
	assert.True(t, data.Broadcasts[0].VoiceEnabled)
	assert.NotEmpty(t, data.Broadcasts[0].CreatedAt)

	// 最近积分：4 小时窗口，最新在前。
	require.Len(t, data.RecentScores, 2)
	assert.Equal(t, 3, data.RecentScores[0].Amount)
	require.NotNil(t, data.RecentScores[0].StudentName)
	assert.Equal(t, "甲", *data.RecentScores[0].StudentName)
	assert.Equal(t, "举手发言", data.RecentScores[0].Reason)
	assert.Contains(t, data.RecentScores[0].Time, "ago")
	assert.Equal(t, -2, data.RecentScores[1].Amount)

	// 班级不存在 → 404。
	_, err = f.Svc.InitialData(999999)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))
}

// 大屏课表：today_weekday + 班级名 + 科目/节次/排课（过滤空科目行）。
func TestDisplayTimetableAndExportCses(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	subject := models.Subject{SchoolID: f.School.ID, Name: "语文", SimplifiedName: "语", Color: "#123456", SortOrder: 0}
	require.NoError(t, db.Create(&subject).Error)
	require.NoError(t, db.Create(&models.ClassPeriod{SchoolID: f.School.ID, PeriodIndex: 1, Name: "第一节", StartTime: "08:00", EndTime: "08:45"}).Error)

	require.NoError(t, db.Create(&models.TimetableEntry{ClassID: f.Class.ID, Weekday: 1, PeriodIndex: 1, WeekType: "all", SubjectID: &subject.ID, TeacherName: "李老师", Room: "101"}).Error)
	// 空科目行应被过滤（Laravel forDisplay 的 ->filter(subject_name !== null)）。
	require.NoError(t, db.Create(&models.TimetableEntry{ClassID: f.Class.ID, Weekday: 2, PeriodIndex: 2, WeekType: "", SubjectID: nil}).Error)

	data, err := f.Svc.Timetable(f.Class.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, data.TodayWeekday, 1)
	assert.LessOrEqual(t, data.TodayWeekday, 7)
	assert.Equal(t, "一年级（1）班", data.ClassName)
	require.Len(t, data.Subjects, 1)
	assert.Equal(t, "语文", data.Subjects[0].Name)
	assert.Equal(t, "#123456", data.Subjects[0].Color)
	require.Len(t, data.Periods, 1)
	assert.Equal(t, "08:00", data.Periods[0].StartTime)
	require.Len(t, data.Entries, 1)
	assert.Equal(t, 1, data.Entries[0].Weekday)
	require.NotNil(t, data.Entries[0].SubjectName)
	assert.Equal(t, "语文", *data.Entries[0].SubjectName)
	assert.Equal(t, "all", data.Entries[0].WeekType)

	// CSES 导出（内容与教师端同源）。
	yaml, class, err := f.Svc.ExportCses(f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "一年级（1）班", class.Name)
	assert.True(t, strings.HasPrefix(yaml, "version: 1"))
	assert.Contains(t, yaml, "subject: \"语文\"")

	// 班级不存在 → 404。
	_, err = f.Svc.Timetable(999999)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))
	_, _, err = f.Svc.ExportCses(999999)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))
}

// 教室端班级总览：汇总、平均等级、尖峰、本周积分、榜单与动态。
func TestDisplayClassroomDashboard(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	top := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 100, Status: "active"}
	second := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", TotalScore: 40, Status: "active"}
	require.NoError(t, db.Create(&top).Error)
	require.NoError(t, db.Create(&second).Error)

	// 宠物等级：12（尖峰）与 2 → 平均 7.0。
	require.NoError(t, db.Create(&models.Pet{StudentID: top.ID, ClassID: f.Class.ID, Name: "甲的萌宠", Species: "zhulong", Level: 12, Mood: 80}).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: second.ID, ClassID: f.Class.ID, Name: "乙的萌宠", Species: "qilin", Level: 2, Mood: 80}).Error)

	// 本周内 1 条、7 天前 1 条（只统计本周）。
	now := util.Now()
	require.NoError(t, db.Create(&models.Score{StudentID: top.ID, ClassID: f.Class.ID, Amount: 8, Reason: "考试优秀", CreatedAt: now}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: top.ID, ClassID: f.Class.ID, Amount: 3, Reason: "考试优秀", CreatedAt: now}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: second.ID, ClassID: f.Class.ID, Amount: 5, Reason: "大扫除", CreatedAt: now.AddDate(0, 0, -20)}).Error)

	data, err := f.Svc.ClassroomDashboard(f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "一年级（1）班", data.ClassName)
	assert.Equal(t, 2, data.StudentCount)
	assert.Equal(t, 140, data.TotalScore)
	assert.Equal(t, 7.0, data.AvgPetLevel)
	assert.Equal(t, 1, data.PeakCount)
	assert.Equal(t, 11, data.WeeklyScore)

	require.NotNil(t, data.StarStudent)
	assert.Equal(t, "甲", data.StarStudent.Name)
	assert.Equal(t, 100, data.StarStudent.Score)
	assert.Equal(t, 12, data.StarStudent.PetLevel)
	assert.Equal(t, "zhulong", data.StarStudent.PetSpecies)

	require.Len(t, data.Top5, 2)
	assert.Equal(t, "甲", data.Top5[0].Name)
	assert.Equal(t, "乙", data.Top5[1].Name)

	// recent_news：最近 20 条积分（不限时间）按文案去重、最多 5 条；
	// 两条 CreatedAt 相同的记录按 id DESC 排序（后写入的在前）。
	require.Len(t, data.RecentNews, 3)
	assert.Equal(t, "🎉", data.RecentNews[0].Icon)
	assert.Equal(t, "甲 +3分 — 考试优秀", data.RecentNews[0].Text)
	assert.Equal(t, "🎉", data.RecentNews[1].Icon)
	assert.Equal(t, "甲 +8分 — 考试优秀", data.RecentNews[1].Text, "同文案才去重，金额不同各保留一条")
	assert.Equal(t, "🎉", data.RecentNews[2].Icon)
	assert.Equal(t, "乙 +5分 — 大扫除", data.RecentNews[2].Text)

	// 班级不存在 → 404。
	_, err = f.Svc.ClassroomDashboard(999999)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))
}

// 教室端学生列表 / 宠物概览 / 排行榜 / 班级设置。
func TestDisplayClassroomReadOnlyViews(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	pet := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "007", TotalScore: 30, Status: "active"}
	noPet := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "008", TotalScore: 10, Status: "active"}
	require.NoError(t, db.Create(&pet).Error)
	require.NoError(t, db.Create(&noPet).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: pet.ID, ClassID: f.Class.ID, Name: "甲的萌宠", Species: "zhulong", Level: 3, Experience: 5, Mood: 60}).Error)

	students, err := f.Svc.ClassroomStudents(f.Class.ID)
	require.NoError(t, err)
	require.Len(t, students, 2)
	assert.Equal(t, "007", students[0].StudentNo)
	assert.Equal(t, "甲的萌宠", students[0].PetName)
	assert.Equal(t, 3, students[0].PetLevel)
	assert.Equal(t, "🐣", students[0].PetEmoji)
	assert.False(t, students[0].FreePick, "无 Cache：free_pick 恒为 false")
	assert.Equal(t, "", students[1].PetName)
	assert.Equal(t, 0, students[1].PetLevel)
	assert.Equal(t, "🥚", students[1].PetEmoji)

	pets, err := f.Svc.ClassroomPetsOverview(f.Class.ID)
	require.NoError(t, err)
	require.Len(t, pets, 1)
	assert.Equal(t, "甲", pets[0].StudentName)
	assert.Equal(t, 5, pets[0].Exp)
	assert.Equal(t, "幼年", pets[0].StageName)

	board, err := f.Svc.QuickLeaderboard(f.Class.ID)
	require.NoError(t, err)
	require.Len(t, board, 2)
	assert.Equal(t, 1, board[0].Rank)
	assert.Equal(t, "甲", board[0].Name)
	assert.Equal(t, 30, board[0].Score)
	assert.Equal(t, "007", board[0].No)

	// 班级设置：未配置 → null；配置后原样返回。
	settings, err := f.Svc.ClassSettings(f.Class.ID)
	require.NoError(t, err)
	assert.Nil(t, settings.PetSeries)
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("id = ?", f.Class.ID).
		Update("settings", `{"pet_series":"myth"}`).Error)
	settings, err = f.Svc.ClassSettings(f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "myth", settings.PetSeries)

	// 班级不存在时同样返回 pet_series = null（不报错，同 Laravel）。
	settings, err = f.Svc.ClassSettings(999999)
	require.NoError(t, err)
	assert.Nil(t, settings.PetSeries)
}

// 教室端积分规则：本班班级规则 + 本校默认（学校级）规则。
func TestDisplayScoreRules(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	rules, err := f.Svc.ScoreRules(f.Class.ID)
	require.NoError(t, err)
	assert.Len(t, rules, len(services.DefaultRules), "首次访问应补齐本校默认规则")

	// 班级级自定义规则同样可见。
	require.NoError(t, db.Create(&models.ScoreRule{
		SchoolID: &f.School.ID, ClassID: &f.Class.ID, Name: "本班加分", Amount: 2,
		Category: "custom", IsPositive: true, IsActive: true, SortOrder: 0,
	}).Error)
	rules, err = f.Svc.ScoreRules(f.Class.ID)
	require.NoError(t, err)
	assert.Len(t, rules, len(services.DefaultRules)+1)

	// 跨校规则不可见。
	other := models.School{Name: "别的学校", Code: "other-school", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	require.NoError(t, db.Create(&models.ScoreRule{
		SchoolID: &other.ID, Name: "他校规则", Amount: 1, Category: "custom",
		IsPositive: true, IsActive: true, SortOrder: 0,
	}).Error)
	rules, err = f.Svc.ScoreRules(f.Class.ID)
	require.NoError(t, err)
	for _, rule := range rules {
		assert.NotEqual(t, "他校规则", rule.Name)
	}

	_, err = f.Svc.ScoreRules(999999)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))
}

// 教师端班级码：取当前激活班级、刷新写库与作用域校验。
func TestTeacherDisplayCode(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	view, err := f.Svc.TeacherDisplayCode(&f.Teacher, 0)
	assert.Equal(t, "请先选择班级", appErrorOf(t, err, 400))

	// 走用户的 active_class_id 设置。
	raw, err := f.Teacher.WithSetting(services.ActiveClassSettingKey, f.Class.ID)
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", f.Teacher.ID).Update("settings", raw).Error)
	require.NoError(t, db.First(&f.Teacher, f.Teacher.ID).Error)

	view, err = f.Svc.TeacherDisplayCode(&f.Teacher, 0)
	require.NoError(t, err)
	assert.Equal(t, "LS11", view.Code)
	assert.Equal(t, "一年级（1）班", view.ClassName)
	assert.Nil(t, view.UpdatedAt, "尚未刷新过时 updated_at 为 null")
	require.NotNil(t, view.StudentCount)
	assert.Equal(t, int64(0), *view.StudentCount)
	assert.Empty(t, view.Message)

	// 刷新：写 display_code + updated_at，并返回提示文案。
	refreshed, err := f.Svc.TeacherRefreshDisplayCode(&f.Teacher, f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "LS11", refreshed.Code)
	assert.Equal(t, "班级码已刷新，旧码已失效", refreshed.Message)
	require.NotNil(t, refreshed.UpdatedAt)

	var stored models.ClassRoom
	require.NoError(t, db.First(&stored, f.Class.ID).Error)
	assert.Equal(t, "LS11", stored.DisplayCode)
	require.NotNil(t, stored.DisplayCodeUpdatedAt)

	// 越权班级 → 404。
	_, err = f.Svc.TeacherDisplayCode(&f.Teacher, 999999)
	assert.Equal(t, "班级不存在或不在管辖范围", appErrorOf(t, err, 404))
}

// 管理员班级码：本校可读写、跨校 404、批量重置前缀校验与落库。
func TestAdminDisplayCodeAndReset(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	admin := seedTeacher(t, db, f.School.ID, "admin-display")
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", admin.ID).Update("role", "school_admin").Error)
	require.NoError(t, db.First(&admin, admin.ID).Error)

	view, err := f.Svc.AdminDisplayCode(&admin, f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "LS11", view.Code)
	require.NotNil(t, view.StudentCount)

	refreshed, err := f.Svc.AdminRefreshDisplayCode(&admin, f.Class.ID)
	require.NoError(t, err)
	assert.Equal(t, "LS11", refreshed.Code)
	assert.Empty(t, refreshed.Message, "管理员刷新不返回提示文案（同 Laravel）")
	require.NotNil(t, refreshed.UpdatedAt)

	// 跨校班级 → 404。
	other := models.School{Name: "别的学校", Code: "other-school-2", Status: "active"}
	require.NoError(t, db.Create(&other).Error)
	foreign := models.ClassRoom{SchoolID: other.ID, Grade: "一年级", Name: "一年级（1）班", Status: "active", DisplayCode: "LS11"}
	require.NoError(t, db.Create(&foreign).Error)
	_, err = f.Svc.AdminDisplayCode(&admin, foreign.ID)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))
	_, err = f.Svc.AdminRefreshDisplayCode(&admin, foreign.ID)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, 404))

	// 非法前缀 → 422（Laravel 原文文案）。
	_, err = f.Svc.AdminResetDisplayCodes(&admin, "b1")
	assert.Equal(t, "字母前缀需为 2-4 个英文字母（如 LS / BJ / LE）", appErrorOf(t, err, 422))

	// 合法前缀：持久化到学校设置，并批量重生成本校班级码（不触碰他校）。
	result, err := f.Svc.AdminResetDisplayCodes(&admin, "bj")
	require.NoError(t, err)
	assert.Equal(t, 1, result.Regenerated)
	assert.Equal(t, 0, result.Skipped)
	assert.Empty(t, result.Conflicts)

	var school models.School
	require.NoError(t, db.First(&school, f.School.ID).Error)
	assert.Equal(t, "BJ", school.SettingString("display_code_prefix", ""))

	var stored models.ClassRoom
	require.NoError(t, db.First(&stored, f.Class.ID).Error)
	assert.Equal(t, "BJ11", stored.DisplayCode)

	var foreignStored models.ClassRoom
	require.NoError(t, db.First(&foreignStored, foreign.ID).Error)
	assert.Equal(t, "LS11", foreignStored.DisplayCode, "他校班级码不应被改动")
}

// 未指定学校 + 空前缀时按学校设置解析（回归覆盖：resolved 前缀写回班级码）。
func TestDisplayCodeServiceWithSchoolSettings(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	require.NoError(t, db.Model(&models.School{}).Where("id = ?", f.School.ID).
		Update("settings", `{"display_code_prefix":"LE"}`).Error)

	// 登录解析：库里 display_code 未刷新（LS11），但计算码为 LE11 → 计算码也能登录。
	result, err := f.Svc.Login("LE11", "10.0.0.1", "ua")
	require.NoError(t, err)
	assert.Equal(t, f.Class.ID, result.ClassInfo.ID)
}
