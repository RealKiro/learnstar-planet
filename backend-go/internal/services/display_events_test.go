// 班级大屏事件链路测试：事件总线（每班 seq / 增量消费 / 200 条上限 / 10 分钟过期 / clear）、
// 各发布点确实产生事件、大屏加减分的限制与副作用。
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

// intPtr 取整数指针（consume 的 since 参数）。
func intPtr(n int) *int { return &n }

// countEvents 统计某班事件行数（验证惰性清理 / 裁剪）。
func countEvents(t *testing.T, db *gorm.DB, classID uint) int64 {
	t.Helper()
	var count int64
	require.NoError(t, db.Model(&models.DisplayEvent{}).Where("class_id = ?", classID).Count(&count).Error)
	return count
}

// 事件总线：seq 每班自增、增量消费语义、created_at/类型/载荷字段。
func TestDisplayEventBusSeqAndIncrementalConsume(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	events := services.NewDisplayEvents(db)

	// 空频道：consume 返回空列表（不是错误）。
	got, err := events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	assert.Empty(t, got)

	other := models.ClassRoom{SchoolID: f.School.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&other).Error)

	for i := 1; i <= 3; i++ {
		require.NoError(t, events.Publish(f.Class.ID, services.DisplayEventNotice, map[string]any{"n": i}))
	}
	require.NoError(t, events.Publish(other.ID, services.DisplayEventNotice, map[string]any{"n": 99}))

	all, err := events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, all, 3)
	assert.Equal(t, []int{1, 2, 3}, []int{all[0].ID, all[1].ID, all[2].ID}, "seq 单调自增")
	assert.Equal(t, services.DisplayEventNotice, all[0].Type)
	assert.JSONEq(t, `{"n":1}`, string(all[0].Data))
	_, err = time.Parse(time.RFC3339, all[0].CreatedAt)
	assert.NoError(t, err, "created_at 应为 ISO8601（RFC3339），实际 %q", all[0].CreatedAt)

	// seq 按班级独立：另一个班从 1 开始。
	otherEvents, err := events.Consume(other.ID, nil)
	require.NoError(t, err)
	require.Len(t, otherEvents, 1)
	assert.Equal(t, 1, otherEvents[0].ID)

	// 增量语义：since=1 → 只返回 id > 1。
	inc, err := events.Consume(f.Class.ID, intPtr(1))
	require.NoError(t, err)
	require.Len(t, inc, 2)
	assert.Equal(t, 2, inc[0].ID)
	assert.Equal(t, 3, inc[1].ID)

	// since 已是最新 → 空。
	none, err := events.Consume(f.Class.ID, intPtr(3))
	require.NoError(t, err)
	assert.Empty(t, none)

	// since 超过最新（客户端存了更大的值）→ 空，不报错。
	none, err = events.Consume(f.Class.ID, intPtr(99))
	require.NoError(t, err)
	assert.Empty(t, none)

	// 追加后 seq 继续自增。
	require.NoError(t, events.Publish(f.Class.ID, services.DisplayEventRefresh, map[string]any{"new_code": "LS11"}))
	all, err = events.Consume(f.Class.ID, intPtr(3))
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, 4, all[0].ID)

	// clear：只清本班。
	require.NoError(t, events.Clear(f.Class.ID))
	all, err = events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	assert.Empty(t, all)
	otherEvents, err = events.Consume(other.ID, nil)
	require.NoError(t, err)
	assert.Len(t, otherEvents, 1, "clear 不应影响其他班级")
}

// 10 分钟过期不返回且被惰性清理；每班最多保留 200 条（超出裁剪最早的）。
func TestDisplayEventExpiryAndPerClassTrim(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	events := services.NewDisplayEvents(db)

	// 过期事件（11 分钟前）不返回，并在 consume 时被惰性删除。
	old := models.DisplayEvent{
		ClassID: f.Class.ID, Seq: 1, Type: services.DisplayEventBroadcast, Data: `{}`,
		CreatedAt: util.Now().Add(-11 * time.Minute),
	}
	require.NoError(t, db.Create(&old).Error)

	got, err := events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	assert.Empty(t, got, "超过 10 分钟的事件不应返回")
	assert.Equal(t, int64(0), countEvents(t, db, f.Class.ID), "过期事件被惰性清理")

	// 刚写入（<10 分钟）的事件正常返回。
	require.NoError(t, events.Publish(f.Class.ID, services.DisplayEventNotice, map[string]any{"n": 1}))
	got, err = events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 1, got[0].ID)

	// 200 条上限：连发 205 条 → 只保留最近 200 条（seq 6..205），更早的被裁剪。
	capped := models.ClassRoom{SchoolID: f.School.ID, Name: "三班", Status: "active"}
	require.NoError(t, db.Create(&capped).Error)
	for i := 1; i <= 205; i++ {
		require.NoError(t, events.Publish(capped.ID, services.DisplayEventNotice, map[string]any{"n": i}))
	}

	assert.Equal(t, int64(models.DisplayEventMaxPerClass), countEvents(t, db, capped.ID),
		"每班最多保留 %d 条", models.DisplayEventMaxPerClass)

	cappedEvents, err := events.Consume(capped.ID, nil)
	require.NoError(t, err)
	require.Len(t, cappedEvents, models.DisplayEventMaxPerClass)
	assert.Equal(t, 6, cappedEvents[0].ID, "裁剪的是最早的 5 条（seq 1..5）")
	assert.Equal(t, 205, cappedEvents[len(cappedEvents)-1].ID)
}

// 各发布点确实产生事件：积分发放 / 广播 / 通知发布 / 宠物喂食 / 班级码刷新（教师端与管理员端）。
func TestDisplayEventsPublishedByServices(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	events := services.NewDisplayEvents(db)
	scope := services.NewScope(db)

	student := models.Student{ClassID: f.Class.ID, Name: "小明", StudentNo: "001", Status: "active"}
	require.NoError(t, db.Create(&student).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: f.Class.ID, Name: "小明的萌宠", Species: "zhulong",
		Level: 1, Experience: 0, Mood: 80,
	}).Error)

	// 每次动作后取该班全量事件，按类型记录「最新一条」。
	lastByType := map[string]services.DisplayEventView{}
	collect := func() {
		got, err := events.Consume(f.Class.ID, nil)
		require.NoError(t, err)
		for _, ev := range got {
			lastByType[ev.Type] = ev
		}
	}

	// ① 积分发放 → score_update
	_, err := services.NewScoreService(db).GiveScore(&student, 5, "举手发言", f.Teacher.ID, nil)
	require.NoError(t, err)
	collect()
	ev, ok := lastByType[services.DisplayEventScoreUpdate]
	require.True(t, ok, "积分发放应发布 score_update")
	assert.Equal(t, 1, ev.ID, "本班第一条事件 seq = 1")
	assert.JSONEq(t, `{
		"student_id": `+itoa(student.ID)+`,
		"student_name": "小明",
		"student_no": "001",
		"amount": 5,
		"reason": "举手发言",
		"total_score": 5,
		"pet_level": 1,
		"pet_experience": 5,
		"pet_mood": 80
	}`, string(ev.Data))

	// ② 批量发放 → 每个学生各一条 score_update
	extra := models.Student{ClassID: f.Class.ID, Name: "小红", StudentNo: "002", Status: "active"}
	require.NoError(t, db.Create(&extra).Error)
	_, err = services.NewScoreService(db).BatchGive([]*models.Student{&student, &extra}, 3, "全勤", f.Teacher.ID, nil)
	require.NoError(t, err)
	all, err := events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	scoreEvents := 0
	for _, e := range all {
		if e.Type == services.DisplayEventScoreUpdate {
			scoreEvents++
		}
	}
	assert.Equal(t, 3, scoreEvents, "单发 1 条 + 批量 2 条 = 3 条 score_update")

	// ③ 广播发送 → broadcast（载荷逐字对齐 Laravel：id/type/content/display_seconds/voice_enabled/created_at）
	sent, err := services.NewBroadcastService(db, scope).
		Send(&f.Teacher, "上课啦", "banner", true, false, 15, []uint{f.Class.ID})
	require.NoError(t, err)
	require.Equal(t, 1, sent)
	collect()
	ev, ok = lastByType[services.DisplayEventBroadcast]
	require.True(t, ok, "广播发送应发布 broadcast")
	assert.JSONEq(t, `{
		"id": 1,
		"type": "banner",
		"content": "上课啦",
		"display_seconds": 15,
		"voice_enabled": true,
		"created_at": "`+isoOf(broadcastCreatedAt(t, db))+`"
	}`, string(ev.Data))

	// ④ 通知发布 → notice（id/title/content/type/published_at）
	noticeSvc := services.NewNoticeService(db, scope)
	notice, err := noticeSvc.Create(&f.Teacher, "春游通知", "本周五春游", "event")
	require.NoError(t, err)
	_, err = noticeSvc.Publish(notice)
	require.NoError(t, err)
	collect()
	ev, ok = lastByType[services.DisplayEventNotice]
	require.True(t, ok, "通知发布应发布 notice")
	assert.JSONEq(t, `{
		"id": `+itoa(notice.ID)+`,
		"title": "春游通知",
		"content": "本周五春游",
		"type": "event",
		"published_at": "`+isoOf(*notice.PublishedAt)+`"
	}`, string(ev.Data))

	// ⑤ 宠物喂食 → pet_update（type=feed）
	_, err = services.NewPetService(db, scope).Feed(&f.Teacher, student.ID)
	require.NoError(t, err)
	collect()
	ev, ok = lastByType[services.DisplayEventPetUpdate]
	require.True(t, ok, "宠物喂食应发布 pet_update")
	assert.JSONEq(t, `{
		"student_id": `+itoa(student.ID)+`,
		"student_name": "小明",
		"type": "feed",
		"mood": 100,
		"level": 1,
		"experience": 8
	}`, string(ev.Data))

	// ⑥ 教师端班级码刷新 → refresh {old_code, new_code}
	_, err = f.Svc.TeacherRefreshDisplayCode(&f.Teacher, f.Class.ID)
	require.NoError(t, err)
	collect()
	ev, ok = lastByType[services.DisplayEventRefresh]
	require.True(t, ok, "教师端刷新班级码应发布 refresh")
	assert.JSONEq(t, `{"old_code":"LS11","new_code":"LS11"}`, string(ev.Data))

	// ⑦ 管理员端班级码刷新 → refresh {new_code}（无 old_code）
	admin := seedTeacher(t, db, f.School.ID, "admin-display-events")
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", admin.ID).Update("role", "school_admin").Error)
	require.NoError(t, db.First(&admin, admin.ID).Error)
	_, err = f.Svc.AdminRefreshDisplayCode(&admin, f.Class.ID)
	require.NoError(t, err)
	collect()
	ev, ok = lastByType[services.DisplayEventRefresh]
	require.True(t, ok)
	assert.JSONEq(t, `{"new_code":"LS11"}`, string(ev.Data))

	// 事件累计条数：1(单发) + 2(批量) + 1(广播) + 1(通知) + 1(喂食) + 2(刷新) = 8
	assert.Equal(t, int64(8), countEvents(t, db, f.Class.ID))
}

// 大屏快捷加减分：白名单校验、跨班 404、成功副作用（总分 / 宠物经验 / 审计日志 / 事件）。
func TestDisplayQuickScoreGuardsAndSideEffects(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	events := services.NewDisplayEvents(db)

	student := models.Student{ClassID: f.Class.ID, Name: "小明", StudentNo: "001", Status: "active"}
	require.NoError(t, db.Create(&student).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: f.Class.ID, Name: "小明的萌宠", Species: "zhulong",
		Level: 1, Experience: 0, Mood: 80,
	}).Error)

	// 0 分 → 422；不在白名单的 6 / 30 → 422（Laravel quickScore 的 in:-5,-3,-1,1,3,5）。
	for _, amount := range []int{0, 6, 30, -30} {
		_, err := f.Svc.QuickScore(f.Class.ID, student.ID, amount)
		assert.Equal(t, "无效的分值", appErrorOf(t, err, 422), "amount=%d", amount)
	}

	// 跨班学生 → 404「学生不存在」。
	otherClass := models.ClassRoom{SchoolID: f.School.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	foreign := models.Student{ClassID: otherClass.ID, Name: "他班学生", Status: "active"}
	require.NoError(t, db.Create(&foreign).Error)
	_, err := f.Svc.QuickScore(f.Class.ID, foreign.ID, 5)
	assert.Equal(t, "学生不存在", appErrorOf(t, err, 404))
	// 不存在的学生 → 404。
	_, err = f.Svc.QuickScore(f.Class.ID, 999999, 5)
	assert.Equal(t, "学生不存在", appErrorOf(t, err, 404))

	// 成功：+5（原因固定「课堂表现」）。
	result, err := f.Svc.QuickScore(f.Class.ID, student.ID, 5)
	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalScore)

	// 副作用 1：学生总分落库。
	var stored models.Student
	require.NoError(t, db.First(&stored, student.ID).Error)
	assert.Equal(t, 5, stored.TotalScore)

	// 副作用 2：积分记录（reason = 课堂表现，given_by = 班级教师）。
	var score models.Score
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&score).Error)
	assert.Equal(t, 5, score.Amount)
	assert.Equal(t, "课堂表现", score.Reason)
	assert.Equal(t, f.Teacher.ID, score.GivenBy)
	assert.Equal(t, f.Class.ID, score.ClassID)

	// 副作用 3：审计日志（余额前后）。
	var log models.ScoreLog
	require.NoError(t, db.Where("score_id = ?", score.ID).First(&log).Error)
	assert.Equal(t, 0, log.BalanceBefore)
	assert.Equal(t, 5, log.BalanceAfter)

	// 副作用 4：宠物经验 +5。
	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	assert.Equal(t, 5, pet.Experience)

	// 副作用 5：score_update 事件已发布。
	got, err := events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, services.DisplayEventScoreUpdate, got[0].Type)
	assert.JSONEq(t, `{
		"student_id": `+itoa(student.ID)+`,
		"student_name": "小明",
		"student_no": "001",
		"amount": 5,
		"reason": "课堂表现",
		"total_score": 5,
		"pet_level": 1,
		"pet_experience": 5,
		"pet_mood": 80
	}`, string(got[0].Data))

	// 减分：-3 → 总分 2（不为负）。
	result, err = f.Svc.QuickScore(f.Class.ID, student.ID, -3)
	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalScore)

	// 分组约束：班级无教师时兜底 given_by = 1（Laravel `$teacherId ?: 1`）。
	noTeacher := models.ClassRoom{SchoolID: f.School.ID, Name: "无教师班", Status: "active"}
	require.NoError(t, db.Create(&noTeacher).Error)
	loner := models.Student{ClassID: noTeacher.ID, Name: "孤单", Status: "active"}
	require.NoError(t, db.Create(&loner).Error)
	_, err = f.Svc.QuickScore(noTeacher.ID, loner.ID, 1)
	require.NoError(t, err)
	var lonerScore models.Score
	require.NoError(t, db.Where("student_id = ?", loner.ID).First(&lonerScore).Error)
	assert.Equal(t, f.Teacher.ID, lonerScore.GivenBy, "平台首位教师兜底（本测试里此教师就是首位）")
}

// 教室端加减分：±30 → 403（逐字文案）、0 → 422、成功副作用。
func TestDisplayClassroomGiveScoreGuards(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	student := models.Student{ClassID: f.Class.ID, Name: "小明", StudentNo: "001", Status: "active"}
	require.NoError(t, db.Create(&student).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: student.ID, ClassID: f.Class.ID, Name: "小明的猛宠", Species: "zhulong",
		Level: 1, Experience: 0, Mood: 80,
	}).Error)

	const overLimit = "单次加减分超过 30 分，请使用教师账号登录操作"

	// 0 分 → 422。
	_, err := f.Svc.ClassroomGiveScore(f.Class.ID, student.ID, 0, "课堂表现")
	assert.Equal(t, "分值不能为 0", appErrorOf(t, err, 422))

	// ±30 → 403（含边界 30 与 -30）。
	for _, points := range []int{30, -30, 31} {
		_, err := f.Svc.ClassroomGiveScore(f.Class.ID, student.ID, points, "课堂表现")
		assert.Equal(t, overLimit, appErrorOf(t, err, 403), "points=%d", points)
	}
	// 29 分放行（边界内），随后清理掉该条记录影响。
	_, err = f.Svc.ClassroomGiveScore(f.Class.ID, student.ID, 29, "课堂表现")
	require.NoError(t, err)

	// 原因超 200 字 → 422。
	longReason := ""
	for i := 0; i < 201; i++ {
		longReason += "字"
	}
	_, err = f.Svc.ClassroomGiveScore(f.Class.ID, student.ID, 1, longReason)
	assert.Equal(t, "原因不能超过 200 字", appErrorOf(t, err, 422))

	// 跨班学生 → 404。
	otherClass := models.ClassRoom{SchoolID: f.School.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	foreign := models.Student{ClassID: otherClass.ID, Name: "他班学生", Status: "active"}
	require.NoError(t, db.Create(&foreign).Error)
	_, err = f.Svc.ClassroomGiveScore(f.Class.ID, foreign.ID, 5, "课堂表现")
	assert.Equal(t, "学生不存在", appErrorOf(t, err, 404))

	// 成功：reason 传空 → 回落「课堂评价」（Laravel input('reason', '课堂评价')）。
	result, err := f.Svc.ClassroomGiveScore(f.Class.ID, student.ID, 5, "")
	require.NoError(t, err)
	assert.Equal(t, "小明", result.StudentName)
	assert.Equal(t, 5, result.Points)
	assert.Equal(t, 34, result.NewScore, "29 + 5 = 34")

	var stored models.Student
	require.NoError(t, db.First(&stored, student.ID).Error)
	assert.Equal(t, 34, stored.TotalScore)

	var score models.Score
	require.NoError(t, db.Where("student_id = ? AND amount = ?", student.ID, 5).First(&score).Error)
	assert.Equal(t, "课堂评价", score.Reason)
	assert.Equal(t, f.Teacher.ID, score.GivenBy)

	// 审计日志 + 事件（复用 ScoreService 带来的额外副作用，见方法注释）。
	var logCount int64
	require.NoError(t, db.Model(&models.ScoreLog{}).Where("student_id = ?", student.ID).Count(&logCount).Error)
	assert.Equal(t, int64(2), logCount, "两次成功加减分各写一条审计日志")

	var pet models.Pet
	require.NoError(t, db.Where("student_id = ?", student.ID).First(&pet).Error)
	// 宠物经验：+29 先把 Lv.1 推到 Lv.2（20 点）并结转 9 点，再 +5 → 14；等级由总分折算仍为 Lv.2。
	assert.Equal(t, 14, pet.Experience)
	assert.Equal(t, 2, pet.Level)

	got, err := services.NewDisplayEvents(db).Consume(f.Class.ID, nil)
	require.NoError(t, err)
	assert.Len(t, got, 2, "每次成功加减分各发布一条 score_update")
	assert.Equal(t, services.DisplayEventScoreUpdate, got[0].Type)
}

// 教室端批量加减分：±30 → 403、空名单/超 50/全是外班学生 → 422、成功逐学生副作用。
func TestDisplayClassroomBatchGiveScoreGuardsAndSideEffects(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	events := services.NewDisplayEvents(db)

	s1 := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", Status: "active"}
	s2 := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", Status: "active"}
	require.NoError(t, db.Create(&s1).Error)
	require.NoError(t, db.Create(&s2).Error)
	require.NoError(t, db.Create(&models.Pet{StudentID: s1.ID, ClassID: f.Class.ID, Name: "甲的萌宠", Species: "zhulong", Level: 1, Mood: 80}).Error)

	const overLimit = "单次加减分超过 30 分，请使用教师账号登录操作"

	// 空名单 → 422。
	_, err := f.Svc.ClassroomBatchGiveScore(f.Class.ID, nil, 3, "课堂表现")
	assert.Equal(t, "请至少选择一名学生", appErrorOf(t, err, 422))

	// 超过 50 人 → 422。
	tooMany := make([]uint, 51)
	for i := range tooMany {
		tooMany[i] = uint(i + 1)
	}
	_, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, tooMany, 3, "课堂表现")
	assert.Equal(t, "单次最多操作 50 名学生", appErrorOf(t, err, 422))

	// 0 分 → 422；原因为空 → 422。
	_, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{s1.ID}, 0, "课堂表现")
	assert.Equal(t, "分值不能为 0", appErrorOf(t, err, 422))
	_, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{s1.ID}, 3, "  ")
	assert.Equal(t, "请填写加减分原因", appErrorOf(t, err, 422))

	// ±30 → 403（此判定先于学生查询，同 Laravel）。
	_, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{s1.ID}, -30, "课堂表现")
	assert.Equal(t, overLimit, appErrorOf(t, err, 403))
	_, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{s1.ID}, 30, "课堂表现")
	assert.Equal(t, overLimit, appErrorOf(t, err, 403))

	// 全是外班学生 → 422「未找到可操作的学生」。
	otherClass := models.ClassRoom{SchoolID: f.School.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	foreign := models.Student{ClassID: otherClass.ID, Name: "他班学生", Status: "active"}
	require.NoError(t, db.Create(&foreign).Error)
	_, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{foreign.ID}, 3, "课堂表现")
	assert.Equal(t, "未找到可操作的学生", appErrorOf(t, err, 422))

	// 成功：本班 2 人 + 外班 1 人 → 只处理 2 人（外班被忽略）。
	result, err := f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{s1.ID, s2.ID, foreign.ID}, 3, "课堂表现")
	require.NoError(t, err)
	assert.Equal(t, 2, result.Count)
	assert.Equal(t, 3, result.Points)

	for _, id := range []uint{s1.ID, s2.ID} {
		var stored models.Student
		require.NoError(t, db.First(&stored, id).Error)
		assert.Equal(t, 3, stored.TotalScore, "学生 %d 总分 +3", id)

		var logs int64
		require.NoError(t, db.Model(&models.ScoreLog{}).Where("student_id = ?", id).Count(&logs).Error)
		assert.Equal(t, int64(1), logs)
	}
	var foreignStored models.Student
	require.NoError(t, db.First(&foreignStored, foreign.ID).Error)
	assert.Equal(t, 0, foreignStored.TotalScore, "外班学生不受影响")

	// 逐学生发布 score_update 事件。
	got, err := events.Consume(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, got, 2)
	for _, ev := range got {
		assert.Equal(t, services.DisplayEventScoreUpdate, ev.Type)
	}

	// 减分不为负：-5 后总分为 0。
	result, err = f.Svc.ClassroomBatchGiveScore(f.Class.ID, []uint{s1.ID, s2.ID}, -5, "课堂纪律")
	require.NoError(t, err)
	assert.Equal(t, 2, result.Count)
	var clamped models.Student
	require.NoError(t, db.First(&clamped, s1.ID).Error)
	assert.Equal(t, 0, clamped.TotalScore)
}

// 消费路径与 SSE/poll 同源：ConsumeDisplayEvents 暴露服务层入口。
func TestDisplayConsumeDisplayEvents(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	require.NoError(t, services.NewDisplayEvents(db).Publish(f.Class.ID, services.DisplayEventBroadcast, map[string]any{"content": "上课啦"}))

	got, err := f.Svc.ConsumeDisplayEvents(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, services.DisplayEventBroadcast, got[0].Type)

	got, err = f.Svc.ConsumeDisplayEvents(f.Class.ID, intPtr(1))
	require.NoError(t, err)
	assert.Empty(t, got)
}

// itoa 小工具（避免在多处引入 strconv 的格式化差异）。
func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 12)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

// isoOf 把业务时区时间格式化为 RFC3339（与事件 created_at 的格式一致）。
func isoOf(t time.Time) string { return t.In(util.Loc).Format(time.RFC3339) }

// broadcastCreatedAt 取最新一条广播的创建时间（用于逐字比对事件载荷里的 created_at）。
func broadcastCreatedAt(t *testing.T, db *gorm.DB) time.Time {
	t.Helper()
	var broadcast models.Broadcast
	require.NoError(t, db.Order("id DESC").First(&broadcast).Error)
	return broadcast.CreatedAt
}
