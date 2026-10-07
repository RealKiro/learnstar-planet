// AI 链路路由与鉴权测试：12 条 AI 路由注册、display/ai/* 走班级码 token、
// teacher/admin 走 JWT + 角色隔离；官方账单查询 provider-official 已注册（见 want 列表）。
package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	want := []string{
		"GET /api/v1/admin/ai/settings",
		"PUT /api/v1/admin/ai/settings",
		"POST /api/v1/admin/ai/toggle",
		"GET /api/v1/admin/ai/usage",
		"POST /api/v1/admin/ai/fetch-models",
		"POST /api/v1/admin/ai/test",
		"POST /api/v1/admin/ai/provider-official",
		"GET /api/v1/teacher/ai/config",
		"POST /api/v1/teacher/ai/chat",
		"GET /api/v1/teacher/ai/commands",
		"GET /api/v1/teacher/ai/usage",
		"GET /api/v1/display/ai/settings",
		"POST /api/v1/display/ai/chat",
	}

	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range want {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}
}

// 大屏端 AI：无 token → 401；有班级码 token → 按业务分支返回（未启用 403 / 开关检查 200）。
func TestDisplayAIRoutesAuthAndGuards(t *testing.T) {
	db, engine, class := newTestEngine(t)

	// 无 token → 401（DisplayAuth 文案）
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/ai/settings", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"message":"Token 无效或已过期"}`, w.Body.String())

	// 班级码登录拿 token
	req := httptest.NewRequest(http.MethodPost, "/api/v1/display/login", strings.NewReader(`{"code":"LS11"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var payload struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	token := payload.Data.Token
	require.NotEmpty(t, token)

	// 开关检查：尚未配置 → enabled=false（200）
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/ai/settings?token="+token, nil))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{"data":{"enabled":false},"message":"ok"}`, w.Body.String())

	// 对话：未启用 → 403「AI 功能未开启」
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/display/ai/chat?token="+token,
		strings.NewReader(`{"question":"为什么天是蓝的？"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code, w.Body.String())
	assert.JSONEq(t, `{"message":"AI 功能未开启"}`, w.Body.String())

	// 启用开关（无供应商）→ 403「请先在 AI 中心配置并启用一个供应商」
	setting := models.AISetting{SchoolID: class.SchoolID, Enabled: true, Provider: "openai", Model: "gpt-3.5-turbo", MaxTokens: 2000, TokensLimit: 1000000}
	require.NoError(t, db.Create(&setting).Error)

	w = httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/display/ai/settings?token="+token, nil))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.JSONEq(t, `{"data":{"enabled":true},"message":"ok"}`, w.Body.String())

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/display/ai/chat?token="+token,
		strings.NewReader(`{"question":"为什么天是蓝的？"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code, w.Body.String())
	assert.JSONEq(t, `{"message":"请先在 AI 中心配置并启用一个供应商"}`, w.Body.String())

	// 管理端 / 教师端 AI 路由受 JWT 保护（未登录 401）
	for _, target := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/admin/ai/settings"},
		{http.MethodGet, "/api/v1/admin/ai/usage"},
		{http.MethodPost, "/api/v1/admin/ai/toggle"},
		{http.MethodGet, "/api/v1/teacher/ai/config"},
		{http.MethodPost, "/api/v1/teacher/ai/chat"},
		{http.MethodGet, "/api/v1/teacher/ai/commands"},
		{http.MethodGet, "/api/v1/teacher/ai/usage"},
	} {
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(target.method, target.path, strings.NewReader(`{}`)))
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s %s 应要求登录", target.method, target.path)
	}
}
