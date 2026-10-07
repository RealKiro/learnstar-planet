// 年级战场（PK）测试：教室端/教师端排行口径、本班统计与挑战守卫。
package services_test

import (
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 教室端 PK 排行：同年级聚合（总积分/人数/均级/巅峰/本周）、按总积分降序、isOwn 标记、
// 排除 DEMO 学校、无年级返回空数组。
func TestDisplayPKLeaderboardByGrade(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	pk := services.NewPkService(db, services.NewScope(db))

	// 本班：30 + 10 = 40 分；宠物 8 级（巅峰）与 3 级 → 均级 5.5。
	top := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 30, Status: "active"}
	low := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", TotalScore: 10, Status: "active"}
	require.NoError(t, db.Create(&top).Error)
	require.NoError(t, db.Create(&low).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: top.ID, ClassID: f.Class.ID, Name: "甲宠", Species: "zhulong", Level: 8, Mood: 80}).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: low.ID, ClassID: f.Class.ID, Name: "乙宠", Species: "qilin", Level: 3, Mood: 80}).Error)

	// 本周 7 分（计入）+ 40 天前 5 分（不计入）。
	now := util.Now()
	require.NoError(t, db.Create(&models.Score{StudentID: top.ID, ClassID: f.Class.ID, Amount: 7, Reason: "举手", CreatedAt: now}).Error)
	require.NoError(t, db.Create(&models.Score{StudentID: top.ID, ClassID: f.Class.ID, Amount: 5, Reason: "上上月", CreatedAt: now.AddDate(0, 0, -40)}).Error)

	// 同校同年级另一班：100 分 → 排在本班之前。
	same := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&same).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: same.ID, Name: "丙", TotalScore: 100, Status: "active"}).Error)

	// 同校不同年级：不参战。
	otherGrade := models.ClassRoom{SchoolID: f.School.ID, Grade: "二年级", Name: "二年级（1）班", Status: "active"}
	require.NoError(t, db.Create(&otherGrade).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: otherGrade.ID, Name: "丁", TotalScore: 999, Status: "active"}).Error)

	// 非启用班级：不参战。
	inactive := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "停用班", Status: "archived"}
	require.NoError(t, db.Create(&inactive).Error)

	// 他校同年级：原实现的教室端口径不过滤学校 → 会出现在榜单里。
	otherSchool := seedOtherSchool(t, db)
	foreign := models.ClassRoom{SchoolID: otherSchool.ID, Grade: "一年级", Name: "他校一年级（1）班", Status: "active"}
	require.NoError(t, db.Create(&foreign).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: foreign.ID, Name: "戊", TotalScore: 5, Status: "active"}).Error)

	// DEMO 学校同年级：教室端排除（Laravel whereHas school code != 'DEMO'）。
	demo := models.School{Name: "演示学校", Code: "DEMO", Status: "active"}
	require.NoError(t, db.Create(&demo).Error)
	demoClass := models.ClassRoom{SchoolID: demo.ID, Grade: "一年级", Name: "演示班", Status: "active"}
	require.NoError(t, db.Create(&demoClass).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: demoClass.ID, Name: "己", TotalScore: 888, Status: "active"}).Error)

	rows, err := pk.LeaderboardForClass(f.Class.ID)
	require.NoError(t, err)
	require.Len(t, rows, 3)
	assert.Equal(t, "一年级（2）班", rows[0].Name)
	assert.Equal(t, 100, rows[0].TotalScore)
	assert.False(t, rows[0].IsOwn)

	own := rows[1]
	assert.Equal(t, "一年级（1）班", own.Name)
	assert.True(t, own.IsOwn)
	assert.Equal(t, 40, own.TotalScore)
	assert.Equal(t, 2, own.StudentCount)
	assert.Equal(t, 5.5, own.AvgLevel)
	assert.Equal(t, 1, own.PeakCount)
	assert.Equal(t, 7, own.WeekGrowth)

	assert.Equal(t, "他校一年级（1）班", rows[2].Name)
	for _, row := range rows {
		assert.NotEqual(t, "演示班", row.Name, "DEMO 学校应被排除")
		assert.NotEqual(t, "二年级（1）班", row.Name, "不同年级不参战")
	}

	// 无年级 → 空数组（不报错）。
	noGrade := models.ClassRoom{SchoolID: f.School.ID, Name: "无年级班", Status: "active"}
	require.NoError(t, db.Create(&noGrade).Error)
	rows, err = pk.LeaderboardForClass(noGrade.ID)
	require.NoError(t, err)
	assert.Empty(t, rows)

	// 班级不存在 → 空数组（不报错）。
	rows, err = pk.LeaderboardForClass(999999)
	require.NoError(t, err)
	assert.Empty(t, rows)
}

// 教师端 PK：排行（含 class_id）、本班统计（含名次）与挑战守卫（未选班级 → 400）。
func TestTeacherPKLeaderboardStatsAndChallenge(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	pk := services.NewPkService(db, services.NewScope(db))

	top := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 30, Status: "active"}
	require.NoError(t, db.Create(&top).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: top.ID, ClassID: f.Class.ID, Name: "甲宠", Species: "zhulong", Level: 9, Mood: 80}).Error)

	second := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&second).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: second.ID, Name: "乙", TotalScore: 90, Status: "active"}).Error)

	// 没有任何可管辖班级的教师：排行空、统计全 0（rank 0）、挑战 400。
	lonely := seedTeacher(t, db, f.School.ID, "teacher-lonely")
	rows, err := pk.Leaderboard(&lonely)
	require.NoError(t, err)
	assert.Empty(t, rows)

	stats, err := pk.MyStats(&lonely)
	require.NoError(t, err)
	assert.Equal(t, 0, stats.TotalScore)
	assert.Equal(t, 0, stats.StudentCount)
	assert.Equal(t, 0, stats.Rank)

	_, err = pk.Challenge(&lonely, second.ID)
	assert.Equal(t, "无效的挑战目标", appErrorOf(t, err, 400))

	// 本班班主任：排行两行（同年级），按总积分降序，isOwn 标记本班。
	rows, err = pk.Leaderboard(&f.Teacher)
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, second.ID, rows[0].ClassID)
	assert.Equal(t, "一年级（2）班", rows[0].Name)
	assert.False(t, rows[0].IsOwn)
	assert.Equal(t, 90, rows[0].TotalScore)

	own := rows[1]
	assert.Equal(t, f.Class.ID, own.ClassID)
	assert.True(t, own.IsOwn)
	assert.Equal(t, 30, own.TotalScore)
	assert.Equal(t, 1, own.StudentCount)
	assert.Equal(t, 9.0, own.AvgLevel)
	assert.Equal(t, 1, own.PeakCount)

	// 本班统计：名次 = 同年级内按总积分降序的第 2 名。
	stats, err = pk.MyStats(&f.Teacher)
	require.NoError(t, err)
	assert.Equal(t, 30, stats.TotalScore)
	assert.Equal(t, 1, stats.StudentCount)
	assert.Equal(t, 9.0, stats.AvgLevel)
	assert.Equal(t, 1, stats.PeakCount)
	assert.Equal(t, 0, stats.WeekGrowth)
	assert.Equal(t, 2, stats.Rank)

	// 挑战守卫：目标是本班 → 400；目标不存在 → 404；目标合法 → 成功。
	_, err = pk.Challenge(&f.Teacher, f.Class.ID)
	assert.Equal(t, "无效的挑战目标", appErrorOf(t, err, 400))

	_, err = pk.Challenge(&f.Teacher, 999999)
	assert.Equal(t, "目标班级不存在", appErrorOf(t, err, 404))

	result, err := pk.Challenge(&f.Teacher, second.ID)
	require.NoError(t, err)
	assert.Equal(t, "🚀 挑战已发起！", result.Message)
	assert.Equal(t, "一年级（2）班", result.Data.TargetClass)

	expires, err := time.ParseInLocation("2006-01-02 15:04:05", result.Data.ExpiresAt, util.Loc)
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().AddDate(0, 0, 7), expires, 2*time.Minute)
}

// 同分并列时的名次与排序稳定性（PHP arsort / sortByDesc 在 PHP 8 均为稳定排序）。
func TestPKTieBreakKeepsClassOrder(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	pk := services.NewPkService(db, services.NewScope(db))

	require.NoError(t, db.Create(&models.Student{ClassID: f.Class.ID, Name: "甲", TotalScore: 50, Status: "active"}).Error)
	tie := models.ClassRoom{SchoolID: f.School.ID, Grade: "一年级", Name: "一年级（9）班", Status: "active"}
	require.NoError(t, db.Create(&tie).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: tie.ID, Name: "乙", TotalScore: 50, Status: "active"}).Error)

	rows, err := pk.LeaderboardForClass(f.Class.ID)
	require.NoError(t, err)
	require.Len(t, rows, 2)
	// 同分保持班级 id 升序（本班先建，id 更小）。
	assert.Equal(t, "一年级（1）班", rows[0].Name)
	assert.True(t, rows[0].IsOwn)
	assert.Equal(t, "一年级（9）班", rows[1].Name)

	stats, err := pk.MyStats(&f.Teacher)
	require.NoError(t, err)
	assert.Equal(t, 1, stats.Rank, "同分时本班（id 更小）排在前")
}
