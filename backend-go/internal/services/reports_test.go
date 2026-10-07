// 教师端报表服务层测试：积分趋势 / 宠物等级分布 / 学生进度口径。
package services_test

import (
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newReportFixture 报表测试夹具：学校 + 教师 + 其担任班主任的班级（含一名学生）。
func newReportFixture(t *testing.T, db *gorm.DB) (models.School, models.User, models.ClassRoom, models.Student) {
	t.Helper()
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "report-teacher")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	return school, teacher, class, student
}

// newReportService 构造报表服务。
func newReportService(db *gorm.DB) *services.ReportService {
	return services.NewReportService(db, services.NewScope(db))
}

// 近 N 天得分 / 扣分：按日分组、正负分数组分开、labels 为 m/d、days 收敛到 1..365。
func TestReportScoreTrend(t *testing.T) {
	db := setupDB(t)
	_, teacher, class, student := newReportFixture(t, db)

	now := util.Now()
	today := util.StartOfDay(now)
	makeScore(t, db, class.ID, student.ID, 5, today.Add(9*time.Hour))
	makeScore(t, db, class.ID, student.ID, 3, today.Add(10*time.Hour))
	makeScore(t, db, class.ID, student.ID, -2, today.Add(11*time.Hour))
	makeScore(t, db, class.ID, student.ID, 7, today.AddDate(0, 0, -2).Add(9*time.Hour))
	makeScore(t, db, class.ID, student.ID, -4, today.AddDate(0, 0, -2).Add(10*time.Hour))
	// 窗口外（8 天前）不计入 7 天趋势。
	makeScore(t, db, class.ID, student.ID, 99, today.AddDate(0, 0, -8).Add(9*time.Hour))

	svc := newReportService(db)
	data, err := svc.ScoreTrend(&teacher, 7)
	require.NoError(t, err)

	require.Len(t, data.Labels, 7)
	require.Len(t, data.Datasets, 2)
	assert.Equal(t, "得分", data.Datasets[0].Label)
	assert.Equal(t, "扣分", data.Datasets[1].Label)

	assert.Equal(t, now.Format("01/02"), data.Labels[6], "最后一天为今天")
	assert.Equal(t, 8, data.Datasets[0].Data[6], "今日得分 = 5 + 3 且不含扣分")
	assert.Equal(t, 2, data.Datasets[1].Data[6], "今日扣分取绝对值")
	assert.Equal(t, 7, data.Datasets[0].Data[4], "两天前得分")
	assert.Equal(t, 4, data.Datasets[1].Data[4], "两天前扣分")
	assert.Equal(t, 0, data.Datasets[0].Data[0], "窗口外记录不计入")
	assert.Len(t, data.Datasets[0].Data, 7)

	// days 收敛：< 1 → 1；> 365 → 365。
	one, err := svc.ScoreTrend(&teacher, 0)
	require.NoError(t, err)
	assert.Len(t, one.Labels, 1)

	capped, err := svc.ScoreTrend(&teacher, 999)
	require.NoError(t, err)
	assert.Len(t, capped.Labels, 365)
}

// 宠物等级分布：按 level 升序分组，count 与阶段名取该等级首只宠物。
func TestReportPetDistribution(t *testing.T) {
	db := setupDB(t)
	_, teacher, class, studentA := newReportFixture(t, db)
	studentB := seedExtraStudent(t, db, class.ID)
	studentC := makeStudent(t, db, class.ID, "小刚", "active")

	require.NoError(t, db.Create(&models.Pet{StudentID: studentA.ID, ClassID: class.ID, Name: "甲", Species: "zhulong", Level: 1}).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: studentB.ID, ClassID: class.ID, Name: "乙", Species: "zhulong", Level: 3}).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: studentC.ID, ClassID: class.ID, Name: "丙", Species: "zhulong", Level: 3}).Error)

	rows, err := newReportService(db).PetDistribution(&teacher)
	require.NoError(t, err)
	require.Len(t, rows, 2, "同 level 合并为一行")
	assert.Equal(t, 1, rows[0].Level)
	assert.Equal(t, 1, rows[0].Count)
	assert.Equal(t, "新生之卵", rows[0].StageName)
	assert.Equal(t, 3, rows[1].Level)
	assert.Equal(t, 2, rows[1].Count)
	assert.Equal(t, "幼年", rows[1].StageName)
}

// 单人进度：返回近 50 条历史（倒序）；越权学生 404。
func TestReportStudentProgressDetail(t *testing.T) {
	db := setupDB(t)
	_, teacher, class, student := newReportFixture(t, db)

	now := util.Now()
	for i := 0; i < 55; i++ {
		makeScore(t, db, class.ID, student.ID, i+1, now.Add(-time.Duration(i)*time.Minute))
	}

	svc := newReportService(db)
	data, err := svc.StudentProgress(&teacher, student.ID)
	require.NoError(t, err)
	detail, ok := data.(*services.StudentProgressDetail)
	require.True(t, ok, "带 student_id 时返回单人结构")
	assert.Equal(t, student.ID, detail.Student.ID)
	assert.Equal(t, "小明", detail.Student.Name)
	require.Len(t, detail.History, 50, "最多 50 条")
	assert.Equal(t, 1, detail.History[0].Amount, "最新一条在最前（i=0 → amount 1）")
	assert.Equal(t, 50, detail.History[49].Amount, "第 50 条（超出部分被截断）")

	// 越权 / 不存在学生 → 404。
	assert.Equal(t, "学生不存在或不在管辖范围", appErrorOf(t, err2(svc.StudentProgress(&teacher, student.ID+999)), 404))
}

// err2 把 (any, error) 归一化为 error，便于单行断言。
func err2(_ any, err error) error { return err }

// 全班进度：每人近 10 条 + 前后 5 条对比的涨跌趋势；inactive 学生不计入。
func TestReportStudentProgressClassTrend(t *testing.T) {
	db := setupDB(t)
	_, teacher, class, studentA := newReportFixture(t, db)
	studentB := seedExtraStudent(t, db, class.ID)
	studentC := makeStudent(t, db, class.ID, "小刚", "active")
	makeStudent(t, db, class.ID, "已转出", "inactive")

	now := util.Now()
	// A：10 条各 +1 → change = 5 - 5 = 0 → stable。
	for i := 0; i < 10; i++ {
		makeScore(t, db, class.ID, studentA.ID, 1, now.Add(-time.Duration(i)*time.Minute))
	}
	// B：近 5 条 +3（前 5 条），更早 5 条 -3 → change = 15 + 15 = 30 → up。
	for i := 0; i < 5; i++ {
		makeScore(t, db, class.ID, studentB.ID, 3, now.Add(-time.Duration(i)*time.Minute))
		makeScore(t, db, class.ID, studentB.ID, -3, now.Add(-time.Duration(i+50)*time.Minute))
	}
	// C：近 5 条 -3，更早 5 条 +3 → change = -15 - 15 = -30 → down。
	for i := 0; i < 5; i++ {
		makeScore(t, db, class.ID, studentC.ID, -3, now.Add(-time.Duration(i)*time.Minute))
		makeScore(t, db, class.ID, studentC.ID, 3, now.Add(-time.Duration(i+50)*time.Minute))
	}

	progress, err := newReportService(db).StudentProgress(&teacher, 0)
	require.NoError(t, err)
	rows, isList := progress.([]services.StudentProgressRow)
	require.True(t, isList, "无 student_id 时返回全班列表")
	require.Len(t, rows, 3, "inactive 学生不计入")

	byID := map[uint]services.StudentProgressRow{}
	for _, row := range rows {
		byID[row.StudentID] = row
	}
	assert.Len(t, byID[studentA.ID].Scores, 10)
	assert.Equal(t, "stable", byID[studentA.ID].Trend)
	assert.Equal(t, 0, byID[studentA.ID].Change)
	assert.Equal(t, "up", byID[studentB.ID].Trend)
	assert.Equal(t, 30, byID[studentB.ID].Change)
	assert.Equal(t, "down", byID[studentC.ID].Trend)
	assert.Equal(t, -30, byID[studentC.ID].Change)
}

// 其它学校的班级/学生不进本教师报表范围。
func TestReportScopeIsolation(t *testing.T) {
	db := setupDB(t)
	_, teacher, _, _ := newReportFixture(t, db)

	otherSchool := seedOtherSchool(t, db)
	otherTeacher := seedTeacher(t, db, otherSchool.ID, "other-teacher")
	otherClass, otherStudent := seedTeacherClass(t, db, otherSchool.ID, otherTeacher.ID, "二班")
	makeScore(t, db, otherClass.ID, otherStudent.ID, 42, util.Now())
	require.NoError(t, db.Create(&models.Pet{StudentID: otherStudent.ID, ClassID: otherClass.ID, Name: "他班宠物", Species: "zhulong", Level: 6}).Error)

	svc := newReportService(db)

	rows, err := svc.PetDistribution(&teacher)
	require.NoError(t, err)
	assert.Empty(t, rows, "他校宠物不计入")

	trend, err := svc.ScoreTrend(&teacher, 7)
	require.NoError(t, err)
	assert.Equal(t, 0, trend.Datasets[0].Data[6], "他校积分不计入")

	// 他校学生不在管辖范围：单人进度 → 404「学生不存在或不在管辖范围」。
	_, err = svc.StudentProgress(&teacher, otherStudent.ID)
	assert.Equal(t, "学生不存在或不在管辖范围", appErrorOf(t, err, 404))
}
