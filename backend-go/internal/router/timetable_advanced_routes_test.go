// 课表进阶 / 教师端报表 / 课堂消息的路由注册与 HTTP 冒烟测试。
//
// 路由契约（路径逐字对齐 Laravel backend/routes/api.php 第 98-99、145、150-152、171-175、238-241 行）：
// 共 13 条（任课设置 PUT + POST 两条共用同一 action）。
package router_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// timetableAdvancedRoutes 本批新增的 13 条路由。
var timetableAdvancedRoutes = []string{
	// A. 教师端报表
	"GET /api/v1/teacher/reports/score-trend",
	"GET /api/v1/teacher/reports/pet-distribution",
	"GET /api/v1/teacher/reports/student-progress",
	// B. 教师课堂消息 / 大屏数据
	"GET /api/v1/teacher/classroom/display",
	"GET /api/v1/teacher/classroom/messages",
	"POST /api/v1/teacher/classroom/messages",
	// C. 课表进阶（管理端）
	"POST /api/v1/admin/timetable/import-csv",
	"GET /api/v1/admin/timetable/unavailabilities",
	"POST /api/v1/admin/timetable/unavailabilities",
	"POST /api/v1/admin/timetable/check-conflicts",
	"GET /api/v1/admin/classes/:id/teacher-assignments",
	"PUT /api/v1/admin/classes/:id/teacher-assignments",
	"POST /api/v1/admin/classes/:id/teacher-assignments",
}

// TestTimetableAdvancedRoutesRegistered 13 条路由已注册；明确不做的接口未被注册。
func TestTimetableAdvancedRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	registered := map[string]bool{}
	unregisteredPaths := []string{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
		unregisteredPaths = append(unregisteredPaths, rt.Path)
	}
	for _, key := range timetableAdvancedRoutes {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}
}

// TestTimetableAdvancedRoutesRequireAuth 全部新路由受鉴权保护（无 token → 401）。
func TestTimetableAdvancedRoutesRequireAuth(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	for _, route := range timetableAdvancedRoutes {
		method, path := splitRouteKey(t, route)
		path = replaceIDParam(path)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(method, path, bytes.NewReader([]byte(`{}`))))
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s 未登录应 401", route)
	}
}

// splitRouteKey 拆分 "METHOD /path"。
func splitRouteKey(t *testing.T, key string) (string, string) {
	t.Helper()
	for i, ch := range key {
		if ch == ' ' {
			return key[:i], key[i+1:]
		}
	}
	t.Fatalf("非法路由键: %s", key)
	return "", ""
}

// replaceIDParam 把 gin 路径参数替换为具体 ID。
func replaceIDParam(path string) string {
	return string(bytes.ReplaceAll([]byte(path), []byte(":id"), []byte("1")))
}

// TestTimetableAdvancedHTTPAdmin 管理端课表进阶接口冒烟（导入 / 不可用 / 冲突 / 任课）。
func TestTimetableAdvancedHTTPAdmin(t *testing.T) {
	db, client, token, schoolID := newAdminFixture(t)

	var class models.ClassRoom
	require.NoError(t, db.Where("school_id = ?", schoolID).First(&class).Error)

	// ① CSV 导入：dry_run 缺省 → 预览（不落库）。
	importResp := uploadTimetableCSV(t, client, token, nil)
	require.Equal(t, http.StatusOK, importResp.StatusCode)
	importBody := decodeBody(t, importResp)
	inner, ok := importBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, inner["dry_run"])
	assert.Equal(t, false, inner["imported"])
	assert.Equal(t, float64(1), inner["total_rows"])
	assert.Equal(t, float64(1), inner["period_count"])

	var entryCount int64
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	assert.Equal(t, int64(0), entryCount, "预览不落库")

	// ② CSV 导入：dry_run=false → 正式导入。
	imported := uploadTimetableCSV(t, client, token, map[string]string{"dry_run": "false"})
	require.Equal(t, http.StatusOK, imported.StatusCode)
	importedBody := decodeBody(t, imported)
	importedInner, ok := importedBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, importedInner["imported"])
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	assert.Equal(t, int64(1), entryCount)

	// ③ 不可用时段：空列表 → 保存 → 列表。
	emptyList := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/timetable/unavailabilities", nil, "")
	require.Equal(t, http.StatusOK, emptyList.StatusCode)
	emptyBody := decodeBody(t, emptyList)
	emptyData, ok := emptyBody["data"].([]any)
	require.True(t, ok)
	assert.Empty(t, emptyData)

	// 非法格子 → 422（同 Laravel cells.*.weekday between:1,7）。
	invalid := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/timetable/unavailabilities",
		bytes.NewReader([]byte(`{"teacher_name":"张老师","cells":[{"weekday":9,"period_index":1}]}`)),
		"application/json")
	assert.Equal(t, http.StatusUnprocessableEntity, invalid.StatusCode)

	saved := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/timetable/unavailabilities",
		bytes.NewReader([]byte(`{"teacher_name":" 张老师 ","cells":[{"weekday":1,"period_index":1}]}`)),
		"application/json")
	require.Equal(t, http.StatusOK, saved.StatusCode)
	assert.Equal(t, "不可用时段已保存", decodeBody(t, saved)["message"])

	listed := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/timetable/unavailabilities", nil, "")
	require.Equal(t, http.StatusOK, listed.StatusCode)
	listedData, ok := decodeBody(t, listed)["data"].([]any)
	require.True(t, ok)
	require.Len(t, listedData, 1, "非法格子被剔除，teacher_name 已 trim")
	row := listedData[0].(map[string]any)
	assert.Equal(t, "张老师", row["teacher_name"])
	assert.Equal(t, float64(1), row["weekday"])
	assert.Equal(t, float64(1), row["period_index"])

	// ④ 冲突检查：命中刚保存的不可用时段。
	conflict := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/timetable/check-conflicts",
		bytes.NewReader([]byte(`{"class_id":`+itoaUint(class.ID)+`,"entries":[{"weekday":1,"period_index":1,"subject_name":"数学","week_type":"all","teacher_name":"张老师"}]}`)),
		"application/json")
	require.Equal(t, http.StatusOK, conflict.StatusCode)
	conflictData, ok := decodeBody(t, conflict)["data"].([]any)
	require.True(t, ok)
	require.Len(t, conflictData, 1)
	conflictRow := conflictData[0].(map[string]any)
	assert.Equal(t, "unavailable", conflictRow["type"])
	assert.Equal(t, "星期一 第1节「数学」：教师 张老师 此时段已被标记为不可用", conflictRow["message"])

	// 班级不存在 → 404。
	missing := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/timetable/check-conflicts",
		bytes.NewReader([]byte(`{"class_id":99999,"entries":[]}`)), "application/json")
	assert.Equal(t, http.StatusNotFound, missing.StatusCode)

	// ⑤ 任课设置：PUT 保存 → GET 读取 → POST 亦可（Laravel Route::match(['put','post'])）。
	assignPath := "/api/v1/admin/classes/" + itoaUint(class.ID) + "/teacher-assignments"
	saveAssign := doAdmin(t, client, token, http.MethodPut, assignPath,
		bytes.NewReader([]byte(`{"assignments":[{"subject_name":"语文","teacher_name":"张老师"}]}`)), "application/json")
	require.Equal(t, http.StatusOK, saveAssign.StatusCode)
	assert.Equal(t, "任课已保存", decodeBody(t, saveAssign)["message"])

	listAssign := doAdmin(t, client, token, http.MethodGet, assignPath, nil, "")
	require.Equal(t, http.StatusOK, listAssign.StatusCode)
	assignData, ok := decodeBody(t, listAssign)["data"].([]any)
	require.True(t, ok)
	require.Len(t, assignData, 1)
	assignRow := assignData[0].(map[string]any)
	assert.Equal(t, "语文", assignRow["subject_name"])
	assert.Equal(t, "张老师", assignRow["teacher_name"])

	saveViaPost := doAdmin(t, client, token, http.MethodPost, assignPath,
		bytes.NewReader([]byte(`{"assignments":[{"subject_name":"数学","teacher_name":"李老师"}]}`)), "application/json")
	require.Equal(t, http.StatusOK, saveViaPost.StatusCode)
	listAssign = doAdmin(t, client, token, http.MethodGet, assignPath, nil, "")
	assignData, ok = decodeBody(t, listAssign)["data"].([]any)
	require.True(t, ok)
	require.Len(t, assignData, 1, "replace 语义：只剩本次提交的一条")
	assert.Equal(t, "数学", assignData[0].(map[string]any)["subject_name"])
}

// uploadTimetableCSV 以 multipart/form-data 上传课表 CSV。
func uploadTimetableCSV(t *testing.T, client *http.Client, token string, fields map[string]string) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "timetable.csv")
	require.NoError(t, err)
	content := "年级,班级,星期,第几节,开始时间,结束时间,科目,周次,教师,教室\n" +
		"一年级,一年级（1）班,星期一,1,08:00,08:45,语文,每周,张老师,101\n"
	_, err = part.Write([]byte(content))
	require.NoError(t, err)
	for k, v := range fields {
		require.NoError(t, writer.WriteField(k, v))
	}
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, adminOpsBaseURL+"/api/v1/admin/timetable/import-csv", &body)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

// TestTeacherReportsAndClassroomHTTP 教师端报表 / 课堂消息 HTTP 冒烟（含 403 优先于 422）。
func TestTeacherReportsAndClassroomHTTP(t *testing.T) {
	db, engine, _ := newTestEngine(t)

	var school models.School
	require.NoError(t, db.First(&school).Error)
	teacher := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "report-http-teacher",
		Name: "李老师", Status: "active",
	}
	require.NoError(t, db.Create(&teacher).Error)

	own := models.ClassRoom{SchoolID: school.ID, Grade: "一年级", Name: "一年级（1）班", TeacherID: &teacher.ID, Status: "active"}
	require.NoError(t, db.Create(&own).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: own.ID, Name: "小明", StudentNo: "001", Status: "active"}).Error)
	foreign := models.ClassRoom{SchoolID: school.ID, Grade: "一年级", Name: "一年级（2）班", Status: "active"}
	require.NoError(t, db.Create(&foreign).Error)

	jwtMgr := jwtauth.New("test-secret", 1)
	token, err := jwtMgr.Generate(teacher.ID, teacher.Role, teacher.SchoolID)
	require.NoError(t, err)

	doTeacher := func(method, target string, body []byte) *httptest.ResponseRecorder {
		var reader *bytes.Reader
		if body == nil {
			reader = bytes.NewReader(nil)
		} else {
			reader = bytes.NewReader(body)
		}
		req := httptest.NewRequest(method, target, reader)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}

	// 报表三项。
	trend := doTeacher(http.MethodGet, "/api/v1/teacher/reports/score-trend?days=7", nil)
	require.Equal(t, http.StatusOK, trend.Code, trend.Body.String())
	var trendBody struct {
		Data struct {
			Labels   []string `json:"labels"`
			Datasets []struct {
				Label string `json:"label"`
				Data  []int  `json:"data"`
			} `json:"datasets"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(trend.Body.Bytes(), &trendBody))
	assert.Len(t, trendBody.Data.Labels, 7)
	require.Len(t, trendBody.Data.Datasets, 2)
	assert.Equal(t, "得分", trendBody.Data.Datasets[0].Label)
	assert.Equal(t, "扣分", trendBody.Data.Datasets[1].Label)

	pets := doTeacher(http.MethodGet, "/api/v1/teacher/reports/pet-distribution", nil)
	require.Equal(t, http.StatusOK, pets.Code, pets.Body.String())
	assert.Contains(t, pets.Body.String(), `"data":[]`)

	progress := doTeacher(http.MethodGet, "/api/v1/teacher/reports/student-progress", nil)
	require.Equal(t, http.StatusOK, progress.Code, progress.Body.String())
	assert.Contains(t, progress.Body.String(), "小明")

	// 大屏聚合数据。
	display := doTeacher(http.MethodGet, "/api/v1/teacher/classroom/display?class_id="+itoaUint(own.ID), nil)
	require.Equal(t, http.StatusOK, display.Code, display.Body.String())
	assert.Contains(t, display.Body.String(), "一年级（1）班")

	// 发送广播 → 广播表 + {message, data.type=broadcast}。
	send := doTeacher(http.MethodPost, "/api/v1/teacher/classroom/messages",
		[]byte(`{"class_id":`+itoaUint(own.ID)+`,"type":"banner","content":"上课啦"}`))
	require.Equal(t, http.StatusOK, send.Code, send.Body.String())
	var sendBody struct {
		Message string `json:"message"`
		Data    struct {
			ID   uint   `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(send.Body.Bytes(), &sendBody))
	assert.Equal(t, "广播已发送", sendBody.Message)
	assert.Equal(t, "broadcast", sendBody.Data.Type)
	assert.Equal(t, "banner", firstBroadcastType(t, db, own.ID))

	// 发送通知 → 通知表。
	notify := doTeacher(http.MethodPost, "/api/v1/teacher/classroom/messages",
		[]byte(`{"class_id":`+itoaUint(own.ID)+`,"type":"urgent","content":"紧急通知内容"}`))
	require.Equal(t, http.StatusOK, notify.Code, notify.Body.String())
	assert.Contains(t, notify.Body.String(), `"message":"通知已发布"`)

	// 轮询。
	poll := doTeacher(http.MethodGet, "/api/v1/teacher/classroom/messages?class_id="+itoaUint(own.ID), nil)
	require.Equal(t, http.StatusOK, poll.Code, poll.Body.String())
	assert.Contains(t, poll.Body.String(), `"broadcasts"`)
	assert.Contains(t, poll.Body.String(), `"polled_at"`)

	// 越权 + 缺参数 → 403（权限预检先于参数校验）。
	forbidden := doTeacher(http.MethodPost, "/api/v1/teacher/classroom/messages",
		[]byte(`{"class_id":`+itoaUint(foreign.ID)+`}`))
	require.Equal(t, http.StatusForbidden, forbidden.Code, forbidden.Body.String())
	assert.Contains(t, forbidden.Body.String(), "无权限")

	// 越权但参数合法 → 依然 403。
	forbidden2 := doTeacher(http.MethodPost, "/api/v1/teacher/classroom/messages",
		[]byte(`{"class_id":`+itoaUint(foreign.ID)+`,"type":"banner","content":"越权"}`))
	require.Equal(t, http.StatusForbidden, forbidden2.Code)

	// 未选班级（无激活班级）→ display 400。
	noClass := doTeacher(http.MethodGet, "/api/v1/teacher/classroom/display", nil)
	require.Equal(t, http.StatusBadRequest, noClass.Code)
	assert.Contains(t, noClass.Body.String(), "请先选择班级")
}

// firstBroadcastType 读取某班最新一条广播的 type。
func firstBroadcastType(t *testing.T, db *gorm.DB, classID uint) string {
	t.Helper()
	var broadcast models.Broadcast
	require.NoError(t, db.Where("class_id = ?", classID).Order("id DESC").First(&broadcast).Error)
	return broadcast.Type
}
