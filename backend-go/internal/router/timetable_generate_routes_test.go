// 课表自动排课路由契约（路径逐字对齐 Laravel backend/routes/api.php 第 147-148 行）：
// POST /admin/timetable/generate（单班纯计算）、POST /admin/timetable/generate-school（全校，commit 落库）。
package router_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// timetableGenerateRoutes 本批新增的 2 条路由。
var timetableGenerateRoutes = []string{
	"POST /api/v1/admin/timetable/generate",
	"POST /api/v1/admin/timetable/generate-school",
}

// TestTimetableGenerateRoutesRegistered 2 条路由已注册。
func TestTimetableGenerateRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range timetableGenerateRoutes {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}
}

// TestTimetableGenerateRoutesRequireAuth 新路由受鉴权保护（无 token → 401）。
func TestTimetableGenerateRoutesRequireAuth(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	for _, route := range timetableGenerateRoutes {
		method, path := splitRouteKey(t, route)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(method, path, bytes.NewReader([]byte(`{}`))))
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s 未登录应 401", route)
	}
}

// TestTimetableGenerateHTTPAdmin 管理端规则排课冒烟：入参校验 422、单班纯计算、全校 commit 落库。
func TestTimetableGenerateHTTPAdmin(t *testing.T) {
	db, client, token, schoolID := newAdminFixture(t)

	var class models.ClassRoom
	require.NoError(t, db.Where("school_id = ?", schoolID).First(&class).Error)

	require.NoError(t, db.Create(&models.ClassPeriod{
		SchoolID: schoolID, PeriodIndex: 1, Name: "第1节", StartTime: "08:00", EndTime: "08:45",
	}).Error)
	require.NoError(t, db.Create(&models.ClassPeriod{
		SchoolID: schoolID, PeriodIndex: 2, Name: "第2节", StartTime: "08:55", EndTime: "09:40",
	}).Error)
	require.NoError(t, db.Create(&models.TimetableTeacherAssignment{
		SchoolID: schoolID, ClassID: class.ID, SubjectName: "语文", TeacherName: "张老师",
	}).Error)

	generatePath := "/api/v1/admin/timetable/generate"
	generateSchoolPath := "/api/v1/admin/timetable/generate-school"

	// ① 入参校验：缺 rules.days / weekly 越界 → 422（同 Laravel generateRules()）
	noDays := doAdmin(t, client, token, http.MethodPost, generatePath,
		bytes.NewReader([]byte(`{"rules":{"subjects":[{"name":"语文","weekly":1}]}}`)), "application/json")
	assert.Equal(t, http.StatusUnprocessableEntity, noDays.StatusCode)

	badWeekly := doAdmin(t, client, token, http.MethodPost, generatePath,
		bytes.NewReader([]byte(`{"rules":{"days":[1,2],"subjects":[{"name":"语文","weekly":36}]}}`)), "application/json")
	assert.Equal(t, http.StatusUnprocessableEntity, badWeekly.StatusCode)

	// ② 单班：纯计算，不落库
	generated := doAdmin(t, client, token, http.MethodPost, generatePath,
		bytes.NewReader([]byte(`{"rules":{"days":[1,2],"subjects":[{"name":"语文","weekly":2,"max_per_day":1}]}}`)), "application/json")
	require.Equal(t, http.StatusOK, generated.StatusCode)
	generatedData, ok := decodeBody(t, generated)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, generatedData["success"])
	entries, ok := generatedData["entries"].([]any)
	require.True(t, ok)
	require.Len(t, entries, 2)
	entry := entries[0].(map[string]any)
	assert.Equal(t, "all", entry["week_type"])
	assert.Nil(t, entry["teacher_name"])
	assert.Nil(t, entry["room"])

	var entryCount int64
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	assert.Equal(t, int64(0), entryCount, "单班排课不落库")

	// ③ 全校：commit 缺省 false → 只验证可行性，返回 classes 而不是 entries
	preview := doAdmin(t, client, token, http.MethodPost, generateSchoolPath,
		bytes.NewReader([]byte(`{"rules":{"days":[1,2],"subjects":[{"name":"语文","weekly":2,"max_per_day":1}]}}`)), "application/json")
	require.Equal(t, http.StatusOK, preview.StatusCode)
	previewData, ok := decodeBody(t, preview)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, previewData["success"])
	assert.Nil(t, previewData["entries"], "全校排课返回 classes")
	previewClasses, ok := previewData["classes"].([]any)
	require.True(t, ok)
	require.Len(t, previewClasses, 1)
	assert.Equal(t, class.Name, previewClasses[0].(map[string]any)["class_name"])
	assert.Equal(t, float64(2), previewClasses[0].(map[string]any)["entry_count"])
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	assert.Equal(t, int64(0), entryCount, "commit 缺省不落库")

	// ④ 全校：commit=true → 事务内落库
	committed := doAdmin(t, client, token, http.MethodPost, generateSchoolPath,
		bytes.NewReader([]byte(`{"rules":{"days":[1,2],"subjects":[{"name":"语文","weekly":2,"max_per_day":1}]},"commit":true}`)), "application/json")
	require.Equal(t, http.StatusOK, committed.StatusCode)
	committedData, ok := decodeBody(t, committed)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, committedData["success"])

	require.NoError(t, db.Model(&models.TimetableEntry{}).Where("class_id = ?", class.ID).Count(&entryCount).Error)
	assert.Equal(t, int64(2), entryCount, "commit=true 落库")
}
