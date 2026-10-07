// 班级大屏事件链路 HTTP 层测试：SSE 响应头/事件帧/心跳/断开处理、
// 轮询降级返回结构、大屏加减分三个接口的限制与副作用。
//
// 说明：本文件位于 handlers 包内（而非 handlers_test），以便把 SSE 时序覆写成极短值，
// 保证 `go test` 不会长时间挂起。SSE 的最长运行被显式压到 200 毫秒级。
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// displayHTTPFixture 大屏 HTTP 测试夹具。
type displayHTTPFixture struct {
	DB     *gorm.DB
	Engine *gin.Engine
	Class  models.ClassRoom
	Token  string
	// ClassID/StudentIDs 供加减分用例使用。
	StudentIDs []uint
}

// newDisplayHTTPFixture 建内存库 + 大屏路由（含 DisplayAuth 中间件），并签一个 disp_ token。
func newDisplayHTTPFixture(t *testing.T) displayHTTPFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), database.GormConfig())
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	// 内存库必须单连接（否则每条连接各自一个空库）；请求会串行化，测试里据此放宽时限。
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.Migrate(db))

	school := models.School{Name: "测试学校", Code: "display-http-school", Status: "active"}
	require.NoError(t, db.Create(&school).Error)
	class := models.ClassRoom{
		SchoolID: school.ID, Grade: "一年级", Name: "一年级（1）班",
		Status: "active", DisplayCode: "LS11",
	}
	require.NoError(t, db.Create(&class).Error)

	studentIDs := make([]uint, 0, 2)
	names := []string{"小明", "小红"}
	for i, name := range names {
		student := models.Student{
			ClassID: class.ID, Name: name, StudentNo: string(rune('1' + i)), Status: "active",
		}
		require.NoError(t, db.Create(&student).Error)
		studentIDs = append(studentIDs, student.ID)
	}
	require.NoError(t, db.Create(&models.Pet{
		StudentID: studentIDs[0], ClassID: class.ID, Name: "小明的萌宠",
		Species: "zhulong", Level: 1, Experience: 0, Mood: 80,
	}).Error)

	const token = "disp_http-test-token"
	require.NoError(t, db.Create(&models.DisplayToken{
		Token: token, ClassID: class.ID, ExpiresAt: time.Now().Add(time.Hour),
	}).Error)

	h := New(db, jwtauth.New("test-secret", 1))
	engine := gin.New()
	group := engine.Group("/display", middleware.DisplayAuth(db))
	group.GET("/sse", h.DisplaySSE)
	group.GET("/poll", h.DisplayPoll)
	group.POST("/quick-score", h.DisplayQuickScore)
	group.POST("/scores/give", h.DisplayGiveScore)
	group.POST("/scores/batch-give", h.DisplayBatchGiveScore)

	return displayHTTPFixture{DB: db, Engine: engine, Class: class, Token: token, StudentIDs: studentIDs}
}

// postJSON 发一个带 token 的 JSON POST 请求。
func postJSON(t *testing.T, engine *gin.Engine, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path+"?token="+token, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

// shortSSETiming 把 SSE 时序压成毫秒级，返回恢复函数（避免测试挂起）。
func shortSSETiming(t *testing.T) func() {
	t.Helper()
	old := sseTiming
	sseTiming = struct {
		Heartbeat    time.Duration
		MaxExecution time.Duration
		IdleSleep    time.Duration
		BusySleep    time.Duration
	}{
		Heartbeat:    30 * time.Millisecond,
		MaxExecution: 200 * time.Millisecond,
		IdleSleep:    10 * time.Millisecond,
		BusySleep:    5 * time.Millisecond,
	}
	return func() { sseTiming = old }
}

// SSE：响应头、事件帧、心跳、结束帧（reconnect）逐项同 Laravel。
func TestDisplaySSEHeadersAndFrames(t *testing.T) {
	f := newDisplayHTTPFixture(t)
	defer shortSSETiming(t)()

	// 先放一条事件，验证事件帧格式：`id: N\nevent: <type>\ndata: <json>\n\n`。
	require.NoError(t, services.NewDisplayEvents(f.DB).Publish(f.Class.ID, services.DisplayEventBroadcast,
		services.DisplayBroadcastData{
			ID: 3, Type: "banner", Content: "上课啦", DisplaySeconds: 15,
			VoiceEnabled: true, CreatedAt: "2026-09-20T08:00:00+08:00",
		}))

	w := httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/sse?token="+f.Token, nil))

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	assert.Equal(t, "keep-alive", w.Header().Get("Connection"))
	assert.Equal(t, "no", w.Header().Get("X-Accel-Buffering"))

	body := w.Body.String()

	// 事件帧（id / event / data 三行 + 空行）。
	require.Contains(t, body, "id: 1\nevent: broadcast\ndata: ")
	frame := body[strings.Index(body, "id: 1\n"):]
	frame = frame[:strings.Index(frame, "\n\n")]
	dataLine := strings.TrimPrefix(frame[strings.LastIndex(frame, "\ndata: "):], "\ndata: ")
	assert.JSONEq(t, `{
		"id": 3,
		"type": "banner",
		"content": "上课啦",
		"display_seconds": 15,
		"voice_enabled": true,
		"created_at": "2026-09-20T08:00:00+08:00"
	}`, dataLine)

	// 心跳帧：`event: heartbeat` + data {"time": <iso8601>}。
	require.Contains(t, body, "event: heartbeat\ndata: {\"time\":\"")
	hb := body[strings.Index(body, "event: heartbeat\ndata: "):]
	hb = hb[len("event: heartbeat\ndata: "):]
	hb = hb[:strings.Index(hb, "\n")]
	var heartbeat struct {
		Time string `json:"time"`
	}
	require.NoError(t, json.Unmarshal([]byte(hb), &heartbeat))
	_, err := time.Parse(time.RFC3339, heartbeat.Time)
	assert.NoError(t, err, "心跳时间应为 ISO8601，实际 %q", heartbeat.Time)

	// 优雅关闭：`event: reconnect` + data {}。
	assert.True(t, strings.HasSuffix(body, "event: reconnect\ndata: {}\n\n"), "尾部应为 reconnect 帧：%q", tail(body, 40))

	// 无 token / 非法 token → 401（走 DisplayAuth）。
	for _, q := range []string{"", "?token=not-a-token"} {
		w = httptest.NewRecorder()
		f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/sse"+q, nil))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"message":"Token 无效或已过期"}`, w.Body.String())
	}
}

// SSE：客户端断开（Context 取消）后立即结束循环，不空转到最长运行时间。
func TestDisplaySSEDisconnectStopsLoop(t *testing.T) {
	f := newDisplayHTTPFixture(t)

	old := sseTiming
	sseTiming = struct {
		Heartbeat    time.Duration
		MaxExecution time.Duration
		IdleSleep    time.Duration
		BusySleep    time.Duration
	}{
		Heartbeat:    20 * time.Millisecond,
		MaxExecution: 10 * time.Second, // 故意设长：若未处理断开，测试会挂 10 秒
		IdleSleep:    20 * time.Millisecond,
		BusySleep:    5 * time.Millisecond,
	}
	defer func() { sseTiming = old }()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(60 * time.Millisecond)
		cancel()
	}()

	w := httptest.NewRecorder()
	start := time.Now()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/sse?token="+f.Token, nil).WithContext(ctx))
	elapsed := time.Since(start)
	cancel()

	assert.Less(t, elapsed, 2*time.Second, "客户端断开后应立即结束（实际 %v）", elapsed)
	assert.Contains(t, w.Body.String(), "event: heartbeat")
	assert.NotContains(t, w.Body.String(), "event: reconnect", "断开路径不应发送优雅关闭帧")
}

// SSE：长连接期间其他请求不被阻塞（轮询仍可正常返回）。
func TestDisplaySSEDoesNotBlockOtherRequests(t *testing.T) {
	f := newDisplayHTTPFixture(t)
	defer shortSSETiming(t)()

	done := make(chan struct{})
	go func() {
		defer close(done)
		w := httptest.NewRecorder()
		f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/sse?token="+f.Token, nil))
	}()

	time.Sleep(20 * time.Millisecond) // 确保 SSE 已在运行
	start := time.Now()
	w := httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/poll?token="+f.Token, nil))
	elapsed := time.Since(start)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Less(t, elapsed, time.Second, "SSE 运行期间轮询应立即返回（实际 %v）", elapsed)
	<-done
}

// 轮询降级：返回结构与 Laravel 一致（events / last_event_id / server_time）+ 增量语义。
func TestDisplayPollStructureAndIncremental(t *testing.T) {
	f := newDisplayHTTPFixture(t)

	require.NoError(t, services.NewDisplayEvents(f.DB).Publish(f.Class.ID, services.DisplayEventNotice,
		services.DisplayNoticeData{
			ID: 7, Title: "春游通知", Content: "本周五春游", Type: "info",
			PublishedAt: "2026-09-20T08:00:00+08:00",
		}))

	type pollPayload struct {
		Data struct {
			Events []struct {
				ID        int             `json:"id"`
				Type      string          `json:"type"`
				Data      json.RawMessage `json:"data"`
				CreatedAt string          `json:"created_at"`
			} `json:"events"`
			LastEventID int    `json:"last_event_id"`
			ServerTime  string `json:"server_time"`
		} `json:"data"`
	}

	// 不传 last_event_id → 返回全部事件。
	w := httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/poll?token="+f.Token, nil))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var payload pollPayload
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload.Data.Events, 1)
	assert.Equal(t, 1, payload.Data.Events[0].ID)
	assert.Equal(t, services.DisplayEventNotice, payload.Data.Events[0].Type)
	assert.JSONEq(t, `{
		"id": 7,
		"title": "春游通知",
		"content": "本周五春游",
		"type": "info",
		"published_at": "2026-09-20T08:00:00+08:00"
	}`, string(payload.Data.Events[0].Data))
	assert.Equal(t, 1, payload.Data.LastEventID, "last_event_id = 本次返回的最大事件 id")
	_, err := time.Parse(time.RFC3339, payload.Data.ServerTime)
	assert.NoError(t, err, "server_time 应为 ISO8601")

	// last_event_id=1（已消费到最新）→ 空事件列表，last_event_id 保持不变。
	w = httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/poll?token="+f.Token+"&last_event_id=1", nil))
	require.Equal(t, http.StatusOK, w.Code)
	payload = pollPayload{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	assert.Empty(t, payload.Data.Events)
	assert.Equal(t, 1, payload.Data.LastEventID)

	// 新事件 → 增量返回。
	require.NoError(t, services.NewDisplayEvents(f.DB).Publish(f.Class.ID, services.DisplayEventRefresh,
		services.DisplayRefreshData{NewCode: "LS11"}))
	w = httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/poll?token="+f.Token+"&last_event_id=1", nil))
	require.Equal(t, http.StatusOK, w.Code)
	payload = pollPayload{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload.Data.Events, 1)
	assert.Equal(t, 2, payload.Data.Events[0].ID)
	assert.Equal(t, services.DisplayEventRefresh, payload.Data.Events[0].Type)
	assert.Equal(t, 2, payload.Data.LastEventID)

	// 无 token → 401。
	w = httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/poll", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"message":"Token 无效或已过期"}`, w.Body.String())
}

// 大屏加减分三个接口：分值校验、±30 文案、副作用、事件发布。
func TestDisplayScoreEndpointsOverHTTP(t *testing.T) {
	f := newDisplayHTTPFixture(t)
	studentID := f.StudentIDs[0]

	// quick-score：+5 成功（原因固定「课堂表现」）。
	w := postJSON(t, f.Engine, "/display/quick-score", f.Token,
		`{"student_id":`+utoa(studentID)+`,"amount":5}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{"data":{"total_score":5},"message":"ok"}`, w.Body.String())

	// quick-score：0 分与白名单外的 30 → 422（Laravel quickScore 的 in:-5,-3,-1,1,3,5）。
	for _, amount := range []string{"0", "30"} {
		w = postJSON(t, f.Engine, "/display/quick-score", f.Token,
			`{"student_id":`+utoa(studentID)+`,"amount":`+amount+`}`)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "amount=%s body=%s", amount, w.Body.String())
		assert.JSONEq(t, `{"message":"无效的分值"}`, w.Body.String())
	}

	// quick-score：缺参 → 422。
	w = postJSON(t, f.Engine, "/display/quick-score", f.Token, `{"amount":5}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	// scores/give：±30 → 403（逐字文案）。
	for _, points := range []string{"30", "-30"} {
		w = postJSON(t, f.Engine, "/display/scores/give", f.Token,
			`{"student_id":`+utoa(studentID)+`,"points":`+points+`,"reason":"课堂表现"}`)
		assert.Equal(t, http.StatusForbidden, w.Code, "points=%s", points)
		assert.JSONEq(t, `{"message":"单次加减分超过 30 分，请使用教师账号登录操作"}`, w.Body.String())
	}

	// scores/give：0 分 → 422。
	w = postJSON(t, f.Engine, "/display/scores/give", f.Token,
		`{"student_id":`+utoa(studentID)+`,"points":0,"reason":"课堂表现"}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"message":"分值不能为 0"}`, w.Body.String())

	// scores/give：+5 成功（文案「加分成功」，data 字段同 Laravel）。
	w = postJSON(t, f.Engine, "/display/scores/give", f.Token,
		`{"student_id":`+utoa(studentID)+`,"points":5,"reason":"举手发言"}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{
		"message": "加分成功",
		"data": {"student_name": "小明", "points": 5, "new_score": 10}
	}`, w.Body.String())

	// scores/give：-3 → 「减分成功」。
	w = postJSON(t, f.Engine, "/display/scores/give", f.Token,
		`{"student_id":`+utoa(studentID)+`,"points":-3,"reason":"课堂纪律"}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{
		"message": "减分成功",
		"data": {"student_name": "小明", "points": -3, "new_score": 7}
	}`, w.Body.String())

	// scores/batch-give：±30 → 403。
	for _, points := range []string{"30", "-30"} {
		w = postJSON(t, f.Engine, "/display/scores/batch-give", f.Token,
			`{"student_ids":[`+utoa(studentID)+`],"points":`+points+`,"reason":"课堂表现"}`)
		assert.Equal(t, http.StatusForbidden, w.Code, "points=%s", points)
		assert.JSONEq(t, `{"message":"单次加减分超过 30 分，请使用教师账号登录操作"}`, w.Body.String())
	}

	// scores/batch-give：空名单 → 422。
	w = postJSON(t, f.Engine, "/display/scores/batch-give", f.Token,
		`{"student_ids":[],"points":3,"reason":"课堂表现"}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"message":"请至少选择一名学生"}`, w.Body.String())

	// scores/batch-give：2 名学生 +3 → 文案逐字同 Laravel。
	w = postJSON(t, f.Engine, "/display/scores/batch-give", f.Token,
		`{"student_ids":[`+utoa(f.StudentIDs[0])+`,`+utoa(f.StudentIDs[1])+`],"points":3,"reason":"全勤"}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{
		"message": "批量操作完成，处理了 2 名学生",
		"data": {"count": 2, "points": 3}
	}`, w.Body.String())

	// 副作用：学生总分（5+5-3+3 = 10；第二个学生 3）。
	var first models.Student
	require.NoError(t, f.DB.First(&first, f.StudentIDs[0]).Error)
	assert.Equal(t, 10, first.TotalScore)

	// 事件：quick-score 1 条 + give 2 条 + batch 2 条 = 5 条 score_update，可由 poll 读到。
	w = httptest.NewRecorder()
	f.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/display/poll?token="+f.Token, nil))
	require.Equal(t, http.StatusOK, w.Code)
	var payload struct {
		Data struct {
			Events []struct {
				ID        int    `json:"id"`
				Type      string `json:"type"`
				CreatedAt string `json:"created_at"`
			} `json:"events"`
			LastEventID int `json:"last_event_id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload.Data.Events, 5)
	for _, ev := range payload.Data.Events {
		assert.Equal(t, services.DisplayEventScoreUpdate, ev.Type)
		_, err := time.Parse(time.RFC3339, ev.CreatedAt)
		assert.NoError(t, err)
	}
	assert.Equal(t, 5, payload.Data.LastEventID)

	// 无 token → 401。
	for _, path := range []string{"/display/quick-score", "/display/scores/give", "/display/scores/batch-give"} {
		w = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		f.Engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "path=%s", path)
	}
}

// 事件 created_at 为业务时区（Asia/Shanghai，同 Laravel APP_TIMEZONE）的 ISO8601。
func TestDisplayEventCreatedAtUsesBusinessTimezone(t *testing.T) {
	f := newDisplayHTTPFixture(t)

	require.NoError(t, services.NewDisplayEvents(f.DB).Publish(f.Class.ID, services.DisplayEventNotice,
		map[string]any{"n": 1}))

	got, err := services.NewDisplayEvents(f.DB).Consume(f.Class.ID, nil)
	require.NoError(t, err)
	require.Len(t, got, 1)

	parsed, err := time.Parse(time.RFC3339, got[0].CreatedAt)
	require.NoError(t, err)
	_, offset := parsed.In(util.Loc).Zone()
	assert.Equal(t, offset, 8*3600, "created_at 应为 +08:00 时区（实际 %q）", got[0].CreatedAt)
}

// utoa 无符号整数转十进制字符串。
func utoa(n uint) string {
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

// tail 取字符串末尾 n 个字符（断言失败时便于定位）。
func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
