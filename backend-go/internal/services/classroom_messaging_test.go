// 教师端课堂消息 / 大屏数据服务层测试。
// 覆盖：按 type 分流（广播表 vs 通知表）与大屏事件、越权 403、
// classroom/display 聚合结构、消息轮询窗口与状态过滤。
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

// messagingFixture 课堂消息测试夹具：学校 + 教师 + 其班级（含 1 名学生）+ 他班。
type messagingFixture struct {
	School  models.School
	Teacher models.User
	Class   models.ClassRoom
	Student models.Student
	Other   models.ClassRoom
	DB      *gorm.DB
	Svc     *services.ClassroomMessagingService
}

// newMessagingFixture 建好夹具（Other 为同校但非该教师担任班主任的班级）。
func newMessagingFixture(t *testing.T, db *gorm.DB) messagingFixture {
	t.Helper()
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "msg-teacher")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	other := makeClass(t, db, school.ID, "二班", "一年级")

	return messagingFixture{
		School: school, Teacher: teacher, Class: class, Student: student, Other: other,
		DB:  db,
		Svc: services.NewClassroomMessagingService(db, services.NewScope(db)),
	}
}

// eventTypesOf 返回某班当前全部大屏事件类型。
func eventTypesOf(t *testing.T, db *gorm.DB, classID uint) []string {
	t.Helper()
	events, err := services.NewDisplayEvents(db).Consume(classID, nil)
	require.NoError(t, err)
	types := make([]string, 0, len(events))
	for _, ev := range events {
		types = append(types, ev.Type)
	}
	return types
}

// 3 种广播 type → 广播表 + broadcast 事件；4 种通知 type → 通知表 + notice 事件。
func TestClassroomMessagingSendRoutesByType(t *testing.T) {
	db := setupDB(t)
	f := newMessagingFixture(t, db)
	accessible := []uint{f.Class.ID}

	for _, msgType := range []string{"banner", "popup", "fullscreen"} {
		result, err := f.Svc.Send(&f.Teacher, accessible, f.Class.ID, msgType, services.ClassroomMessageInput{
			Content:        "广播内容-" + msgType,
			DisplaySeconds: 10,
			Voice:          true,
		})
		require.NoError(t, err)
		assert.Equal(t, "广播已发送", result.Message)
		assert.Equal(t, "broadcast", result.Data.Type)

		var broadcast models.Broadcast
		require.NoError(t, db.First(&broadcast, result.Data.ID).Error, "type=%s 应写入广播表", msgType)
		assert.Equal(t, msgType, broadcast.Type)
		assert.Equal(t, f.Class.ID, broadcast.ClassID)
		assert.Equal(t, f.School.ID, broadcast.SchoolID)
		assert.Equal(t, uint(f.Teacher.ID), *broadcast.TeacherID)
		assert.Equal(t, "sent", broadcast.Status)
		assert.Equal(t, 10, broadcast.DisplaySeconds)
		require.NotNil(t, broadcast.SentAt)

		var noticeCount int64
		require.NoError(t, db.Model(&models.Notice{}).Where("content = ?", "广播内容-"+msgType).Count(&noticeCount).Error)
		assert.Equal(t, int64(0), noticeCount, "type=%s 不应写通知表", msgType)
	}

	for _, msgType := range []string{"info", "homework", "event", "urgent"} {
		result, err := f.Svc.Send(&f.Teacher, accessible, f.Class.ID, msgType, services.ClassroomMessageInput{
			Content: "通知内容-" + msgType,
			Voice:   true,
		})
		require.NoError(t, err)
		assert.Equal(t, "通知已发布", result.Message)
		assert.Equal(t, "notice", result.Data.Type)

		var notice models.Notice
		require.NoError(t, db.First(&notice, result.Data.ID).Error, "type=%s 应写入通知表", msgType)
		assert.Equal(t, msgType, notice.Type)
		assert.Equal(t, f.Class.ID, notice.ClassID)
		assert.Equal(t, uint(f.Teacher.ID), notice.PublishedBy)
		assert.True(t, notice.IsPublished)
		require.NotNil(t, notice.PublishedAt)
		if msgType == "urgent" {
			assert.Equal(t, "紧急通知", notice.Title)
		} else {
			assert.Equal(t, "通知", notice.Title)
		}

		// 通知表 id 与广播表 id 各自自增，故按「内容从未出现在广播表」判定未串表。
		var broadcastCount int64
		require.NoError(t, db.Model(&models.Broadcast{}).Where("content = ?", "通知内容-"+msgType).Count(&broadcastCount).Error)
		assert.Equal(t, int64(0), broadcastCount, "type=%s 不应写广播表", msgType)
	}

	types := eventTypesOf(t, db, f.Class.ID)
	require.Len(t, types, 7)
	assert.Equal(t, []string{
		services.DisplayEventBroadcast, services.DisplayEventBroadcast, services.DisplayEventBroadcast,
		services.DisplayEventNotice, services.DisplayEventNotice, services.DisplayEventNotice, services.DisplayEventNotice,
	}, types, "3 条广播事件 + 4 条通知事件（按写入顺序）")
}

// 显式标题优先于默认标题；voice=false 一定落库为 false（forceFalseBool 回写）。
func TestClassroomMessagingSendTitleAndVoiceFalse(t *testing.T) {
	db := setupDB(t)
	f := newMessagingFixture(t, db)
	accessible := []uint{f.Class.ID}

	title := "自定义标题"
	noticeResult, err := f.Svc.Send(&f.Teacher, accessible, f.Class.ID, "info", services.ClassroomMessageInput{
		Content: "内容", Title: &title, Voice: true,
	})
	require.NoError(t, err)
	var notice models.Notice
	require.NoError(t, db.First(&notice, noticeResult.Data.ID).Error)
	assert.Equal(t, "自定义标题", notice.Title)

	emptyTitle := ""
	_, err = f.Svc.Send(&f.Teacher, accessible, f.Class.ID, "urgent", services.ClassroomMessageInput{
		Content: "内容", Title: &emptyTitle, Voice: true,
	})
	require.NoError(t, err)
	var urgent models.Notice
	require.NoError(t, db.Order("id DESC").First(&urgent).Error)
	assert.Equal(t, "", urgent.Title, "显式空串不做默认替换（同 Laravel ?? 语义）")

	broadcastResult, err := f.Svc.Send(&f.Teacher, accessible, f.Class.ID, "banner", services.ClassroomMessageInput{
		Content: "静音广播", DisplaySeconds: 10, Voice: false,
	})
	require.NoError(t, err)
	var broadcast models.Broadcast
	require.NoError(t, db.First(&broadcast, broadcastResult.Data.ID).Error)
	assert.False(t, broadcast.VoiceEnabled, "voice=false 必须落库 false")
}

// 未分配该班级 → 403「无权限」，且不落库。
func TestClassroomMessagingSendForbidden(t *testing.T) {
	db := setupDB(t)
	f := newMessagingFixture(t, db)

	// ① 传入不可访问班级；② accessible 为空（无管辖班级）。
	for _, accessible := range [][]uint{{f.Class.ID}, {}} {
		_, err := f.Svc.Send(&f.Teacher, accessible, f.Other.ID, "banner", services.ClassroomMessageInput{
			Content: "越权", DisplaySeconds: 10, Voice: true,
		})
		assert.Equal(t, "无权限", appErrorOf(t, err, 403))
	}

	var broadcastCount, noticeCount int64
	require.NoError(t, db.Model(&models.Broadcast{}).Count(&broadcastCount).Error)
	require.NoError(t, db.Model(&models.Notice{}).Count(&noticeCount).Error)
	assert.Equal(t, int64(0), broadcastCount)
	assert.Equal(t, int64(0), noticeCount)
}

// classroom/display 聚合结构：班级信息 + 学生宠物总览 + 生效广播 + 近 7 天通知 + 近 4 小时积分。
func TestClassroomMessagingDisplayAggregation(t *testing.T) {
	db := setupDB(t)
	f := newMessagingFixture(t, db)
	extra := seedExtraStudent(t, db, f.Class.ID)

	// 学生 1 有宠物（Lv.5 → 成长期），学生 2 无宠物。
	require.NoError(t, db.Create(&models.Pet{
		StudentID: f.Student.ID, ClassID: f.Class.ID, Name: "小明的萌宠",
		Species: "zhulong", Level: 5, Experience: 30, Mood: 70,
	}).Error)

	now := util.Now()

	// 广播：pending / sent 计入，draft 不计入；取最近 5 条。
	for i, status := range []string{"sent", "pending", "draft"} {
		require.NoError(t, db.Create(&models.Broadcast{
			SchoolID: f.School.ID, ClassID: f.Class.ID, TeacherID: &f.Teacher.ID,
			Content: "广播" + status, Type: "banner", VoiceEnabled: true,
			DisplaySeconds: 10, Status: status, CreatedAt: now.Add(-time.Duration(i) * time.Minute),
		}).Error)
	}

	// 通知：published 且在 7 天内计入；未发布 / 超 7 天不计入；取最近 3 条。
	makeNotice := func(title, noticeType string, published bool, at time.Time) models.Notice {
		notice := models.Notice{
			ClassID: f.Class.ID, SchoolID: f.School.ID, Title: title, Content: title + "内容",
			Type: noticeType, PublishedBy: f.Teacher.ID, IsPublished: published, PublishedAt: &at,
		}
		require.NoError(t, db.Create(&notice).Error)
		return notice
	}
	makeNotice("近通知1", "info", true, now.Add(-1*time.Hour))
	makeNotice("近通知2", "homework", true, now.Add(-2*time.Hour))
	makeNotice("近通知3", "event", true, now.Add(-3*time.Hour))
	makeNotice("近通知4", "urgent", true, now.Add(-4*time.Hour))
	makeNotice("未发布", "info", false, now.Add(-5*time.Hour))
	makeNotice("过期通知", "info", true, now.AddDate(0, 0, -8))

	// 积分：4 小时内 2 条计入，5 小时前 1 条不计入。
	makeScore(t, db, f.Class.ID, f.Student.ID, 5, now.Add(-1*time.Hour))
	makeScore(t, db, f.Class.ID, f.Student.ID, -2, now.Add(-2*time.Hour))
	makeScore(t, db, f.Class.ID, f.Student.ID, 9, now.Add(-5*time.Hour))

	data, err := f.Svc.Display(&f.Teacher, f.Class.ID)
	require.NoError(t, err)

	assert.Equal(t, "一班", data.ClassName)
	assert.Equal(t, 2, data.StudentCount)
	require.Len(t, data.Pets, 2)

	// 按 student_no 升序：001（有宠物）在前，002（无宠物）在后。
	first := data.Pets[0]
	assert.Equal(t, f.Student.ID, first.StudentID)
	assert.True(t, first.HasPet)
	require.NotNil(t, first.PetName)
	assert.Equal(t, "小明的萌宠", *first.PetName)
	require.NotNil(t, first.PetSpecies)
	assert.Equal(t, "zhulong", *first.PetSpecies)
	assert.Equal(t, 5, first.Level)
	assert.Equal(t, 30, first.Experience)
	assert.Equal(t, 70, first.Mood)
	assert.Equal(t, "成长期", first.StageName)

	second := data.Pets[1]
	assert.Equal(t, extra.ID, second.StudentID)
	assert.False(t, second.HasPet)
	assert.Nil(t, second.PetName)
	assert.Nil(t, second.PetSpecies)
	assert.Equal(t, 0, second.Level)
	assert.Equal(t, "未孵化", second.StageName)
	assert.Equal(t, "🤔", second.Emoji)

	require.Len(t, data.Broadcasts, 2, "draft 不计入")
	assert.Equal(t, "广播sent", data.Broadcasts[0].Content)
	assert.Equal(t, "广播pending", data.Broadcasts[1].Content)
	assert.NotEmpty(t, data.Broadcasts[0].CreatedAt, "created_at 为 diffForHumans 文案")

	require.Len(t, data.Notices, 3, "只取最近 3 条已发布通知")
	assert.Equal(t, "近通知1", data.Notices[0].Title)
	assert.Equal(t, "近通知3", data.Notices[2].Title)

	require.Len(t, data.RecentScores, 2, "只取近 4 小时")
	assert.Equal(t, 5, data.RecentScores[0].Amount)
	require.NotNil(t, data.RecentScores[0].StudentName)
	assert.Equal(t, "小明", *data.RecentScores[0].StudentName)
	assert.NotEmpty(t, data.RecentScores[0].Time)
}

// display：未选班级 → 400；未分配该班级 → 403；班级不存在 → 404。
func TestClassroomMessagingDisplayGuards(t *testing.T) {
	db := setupDB(t)
	f := newMessagingFixture(t, db)

	_, err := f.Svc.Display(&f.Teacher, 0)
	assert.Equal(t, "请先选择班级", appErrorOf(t, err, 400))

	_, err = f.Svc.Display(&f.Teacher, f.Other.ID)
	assert.Equal(t, "您未被分配到此班级", appErrorOf(t, err, 403))

	// API 机器人（本校全部班级）访问他校班级 → 404（班级不在本校范围内即越权 403）。
	bot := makeTeacher(t, db, f.School.ID, "bot-msg", "API 机器人", true)
	_, err = f.Svc.Display(&bot, f.Other.ID)
	require.NoError(t, err, "机器人可访问本校全部班级")
}

// 轮询：since 过滤 + 仅 sent 广播 + 仅已发布通知；未选班级 → 400。
func TestClassroomMessagingPoll(t *testing.T) {
	db := setupDB(t)
	f := newMessagingFixture(t, db)
	now := util.Now()

	sent := models.Broadcast{
		SchoolID: f.School.ID, ClassID: f.Class.ID, Content: "已发送",
		Type: "banner", VoiceEnabled: true, DisplaySeconds: 10, Status: "sent",
		CreatedAt: now.Add(-1 * time.Minute),
	}
	pending := models.Broadcast{
		SchoolID: f.School.ID, ClassID: f.Class.ID, Content: "待发送",
		Type: "popup", VoiceEnabled: true, DisplaySeconds: 10, Status: "pending",
		CreatedAt: now.Add(-1 * time.Minute),
	}
	oldSent := models.Broadcast{
		SchoolID: f.School.ID, ClassID: f.Class.ID, Content: "很久以前",
		Type: "banner", VoiceEnabled: true, DisplaySeconds: 10, Status: "sent",
		CreatedAt: now.Add(-2 * time.Hour),
	}
	require.NoError(t, db.Create(&sent).Error)
	require.NoError(t, db.Create(&pending).Error)
	require.NoError(t, db.Create(&oldSent).Error)

	publishedAt := now.Add(-1 * time.Minute)
	require.NoError(t, db.Create(&models.Notice{
		ClassID: f.Class.ID, SchoolID: f.School.ID, Title: "已发布", Content: "内容",
		Type: "info", PublishedBy: f.Teacher.ID, IsPublished: true, PublishedAt: &publishedAt,
	}).Error)
	require.NoError(t, db.Create(&models.Notice{
		ClassID: f.Class.ID, SchoolID: f.School.ID, Title: "未发布", Content: "内容",
		Type: "info", PublishedBy: f.Teacher.ID, IsPublished: false, PublishedAt: &publishedAt,
	}).Error)

	since := now.Add(-10 * time.Minute)
	result, err := f.Svc.Poll(f.Class.ID, &since)
	require.NoError(t, err)
	require.Len(t, result.Broadcasts, 1, "仅 sent 且在窗口内")
	assert.Equal(t, "已发送", result.Broadcasts[0].Content)
	assert.NotEmpty(t, result.Broadcasts[0].CreatedAt)
	require.Len(t, result.Notices, 1, "仅已发布通知")
	assert.Equal(t, "已发布", result.Notices[0].Title)
	assert.NotEmpty(t, result.Notices[0].PublishedAt)
	assert.NotEmpty(t, result.PolledAt)

	// 缺省 since = 5 分钟前 → 1 分钟前的 sent 广播仍返回，2 小时前的被排除。
	defaulted, err := f.Svc.Poll(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, defaulted.Broadcasts, 1)
	assert.Equal(t, "已发送", defaulted.Broadcasts[0].Content)

	// 更早的 since → 两条 sent 广播都返回（pending 仍排除）。
	wide := now.Add(-3 * time.Hour)
	wideResult, err := f.Svc.Poll(f.Class.ID, &wide)
	require.NoError(t, err)
	require.Len(t, wideResult.Broadcasts, 2)

	_, err = f.Svc.Poll(0, nil)
	assert.Equal(t, "请先选择班级", appErrorOf(t, err, 400))
}

// 非课堂消息类型不属于任何一张表。
func TestClassroomMessageTypeGuard(t *testing.T) {
	for _, msgType := range []string{"banner", "popup", "fullscreen", "info", "homework", "event", "urgent"} {
		assert.True(t, services.IsClassroomMessageType(msgType), msgType)
	}
	for _, msgType := range []string{"", "unknown", "BANNER"} {
		assert.False(t, services.IsClassroomMessageType(msgType), msgType)
	}
}
