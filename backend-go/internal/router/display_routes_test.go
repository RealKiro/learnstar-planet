// 班级大屏路由与中间件测试：路由注册（含静态/参数同级不冲突）、
// DisplayAuth 的 token 取值顺序、401 文案与 CSES 导出响应头。
package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newTestEngine 构建内存库 + 完整路由（router.New 在有路由冲突时会 panic）。
func newTestEngine(t *testing.T) (*gorm.DB, *gin.Engine, models.ClassRoom) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), database.GormConfig())
	require.NoError(t, err)
	require.NoError(t, database.Migrate(db))

	school := models.School{Name: "测试学校", Code: "router-test-school", Status: "active"}
	require.NoError(t, db.Create(&school).Error)
	class := models.ClassRoom{
		SchoolID: school.ID, Grade: "一年级", Name: "一年级（1）班",
		Status: "active", DisplayCode: "LS11",
	}
	require.NoError(t, db.Create(&class).Error)
	require.NoError(t, db.Create(&models.Student{ClassID: class.ID, Name: "小明", StudentNo: "1", Status: "active"}).Error)

	r := router.New(db, &config.Config{JWTSecret: "test-secret", JWTExpHours: 1})
	return db, r, class
}

func TestDisplayRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	want := []string{
		"POST /api/v1/display/login",
		"GET /api/v1/display/initial-data",
		"GET /api/v1/display/timetable",
		"GET /api/v1/display/export-cses",
		"GET /api/v1/display/leaderboard",
		"GET /api/v1/display/class-settings",
		"GET /api/v1/display/dashboard",
		"GET /api/v1/display/students",
		"GET /api/v1/display/scores/rules",
		"GET /api/v1/display/pets/overview",
		"GET /api/v1/display/sse",
		"GET /api/v1/display/poll",
		"POST /api/v1/display/quick-score",
		"POST /api/v1/display/scores/give",
		"POST /api/v1/display/scores/batch-give",
		"GET /api/v1/teacher/display-code",
		"POST /api/v1/teacher/display-code/refresh",
		"GET /api/v1/admin/classes/:id/display-code",
		"POST /api/v1/admin/classes/:id/display-code/refresh",
		"POST /api/v1/admin/classes/reset-display-codes",
		// 本批新增：班级码登录 + 大屏写操作/查询 + 教师端宠物系列与 PK。
		"POST /api/v1/auth/class/login",
		"GET /api/v1/display/shop-items",
		"POST /api/v1/display/redeem",
		"POST /api/v1/display/transfer",
		"POST /api/v1/display/switch-series",
		"POST /api/v1/display/pets/switch",
		"GET /api/v1/display/pk/leaderboard",
		"POST /api/v1/teacher/class/switch-series",
		"GET /api/v1/teacher/pk/leaderboard",
		"POST /api/v1/teacher/pk/challenge",
		"GET /api/v1/teacher/pk/my-stats",
		"GET /api/v1/display/ai/settings",
		"POST /api/v1/display/ai/chat",
	}
	// 已无未移植的接口需在此断言（provider-official 已注册，见 ai_routes_test.go）。
	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range want {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}
}

// 登录 → query token / Bearer token 均可通过；无 token 或非法 token → 401。
func TestDisplayAuthMiddleware(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	body := strings.NewReader(`{"code":"ls11"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/display/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var payload struct {
		Data struct {
			Token     string `json:"token"`
			ExpiresIn int    `json:"expires_in"`
			ClassInfo struct {
				ID           uint   `json:"id"`
				Grade        string `json:"grade"`
				StudentCount int    `json:"student_count"`
			} `json:"class_info"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	token := payload.Data.Token
	require.True(t, strings.HasPrefix(token, "disp_"))
	assert.Equal(t, 86400, payload.Data.ExpiresIn)
	assert.Equal(t, "一年级", payload.Data.ClassInfo.Grade)
	assert.Equal(t, 1, payload.Data.ClassInfo.StudentCount)

	// query 参数 token。
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/class-settings?token="+token, nil))
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{"data":{"pet_series":null},"message":"ok"}`, w.Body.String())

	// Authorization: Bearer token。
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/display/initial-data", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// 缺失 / 非法 token → 401（文案同 Laravel）。
	for _, target := range []string{"/api/v1/display/class-settings", "/api/v1/display/class-settings?token=not-a-token"} {
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"message":"Token 无效或已过期"}`, w.Body.String())
	}

	// 错误 token 前缀不接受。
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/display/class-settings", nil)
	req.Header.Set("Authorization", "Bearer "+strings.TrimPrefix(token, "disp_"))
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 无效班级码登录 → 404。
	req = httptest.NewRequest(http.MethodPost, "/api/v1/display/login", strings.NewReader(`{"code":"LS99"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message":"班级码无效，请检查后重试"}`, w.Body.String())

	// 管理端班级码路由受 JWT 保护（未登录 401）。
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/admin/classes/reset-display-codes", strings.NewReader(`{}`)))
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/admin/classes/1/display-code", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/teacher/display-code", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// CSES 导出响应头：Content-Type / Content-Disposition(rawurlencode) / Cache-Control。
func TestDisplayExportCsesHeaders(t *testing.T) {
	_, engine, class := newTestEngine(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/display/login", strings.NewReader(`{"code":"LS11"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var payload struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))

	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/export-cses?token="+payload.Data.Token, nil))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "application/x-yaml; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	// 文件名 = rawurlencode(班级名 + "-课表.cses.yaml")（UTF-8 逐字节 %XX，空格为 %20 而非 +）。
	assert.Equal(t,
		`attachment; filename="%E4%B8%80%E5%B9%B4%E7%BA%A7%EF%BC%881%EF%BC%89%E7%8F%AD-%E8%AF%BE%E8%A1%A8.cses.yaml"`,
		w.Header().Get("Content-Disposition"))
	assert.True(t, strings.HasPrefix(w.Body.String(), "version: 1"))
	assert.NotEmpty(t, class.Name)
}

// 班级码登录 → class_ token 可直接用于大屏接口；响应无 message 字段。
func TestClassLoginRouteAndTokenReuse(t *testing.T) {
	_, engine, class := newTestEngine(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/class/login", strings.NewReader(`{"class_code":"ls11"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var payload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	_, hasMessage := payload["message"]
	assert.False(t, hasMessage, "Laravel 的 classLogin 响应没有 message 字段")

	var data struct {
		Token        string `json:"token"`
		ClassID      uint   `json:"class_id"`
		ClassName    string `json:"class_name"`
		Grade        string `json:"grade"`
		StudentCount int    `json:"student_count"`
	}
	require.NoError(t, json.Unmarshal(payload["data"], &data))
	require.True(t, strings.HasPrefix(data.Token, "class_LS11_"))
	assert.Equal(t, class.ID, data.ClassID)
	assert.Equal(t, "一年级（1）班", data.ClassName)
	assert.Equal(t, "一年级", data.Grade)
	assert.Equal(t, 1, data.StudentCount)

	// class_ token 直接用于大屏读接口（同一套中间件与 token 表）。
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/class-settings?token="+data.Token, nil))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{"data":{"pet_series":null},"message":"ok"}`, w.Body.String())

	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/shop-items?token="+data.Token, nil))
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/pk/leaderboard?token="+data.Token, nil))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var pkBody struct {
		Data []struct {
			Name  string `json:"name"`
			IsOwn bool   `json:"isOwn"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pkBody))
	require.Len(t, pkBody.Data, 1)
	assert.Equal(t, "一年级（1）班", pkBody.Data[0].Name)
	assert.True(t, pkBody.Data[0].IsOwn)

	// 无效班级码 → 401（文案与大屏登录的「请检查后重试」不同字）。
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/class/login", strings.NewReader(`{"class_code":"LS99"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"message":"班级码无效，请核对后重试"}`, w.Body.String())

	// 缺少 class_code → 422。
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/class/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	// 大屏写接口无 token → 401。
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/display/switch-series", strings.NewReader(`{}`)))
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 教师端 PK / 宠物系列受 JWT 保护。
	for _, target := range []string{"/api/v1/teacher/pk/leaderboard", "/api/v1/teacher/pk/my-stats"} {
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
		assert.Equal(t, http.StatusUnauthorized, w.Code, target)
	}
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/teacher/pk/challenge", strings.NewReader(`{}`)))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/teacher/class/switch-series", strings.NewReader(`{}`)))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
