// 第三方扫码登录 / 绑定 / 通讯录导入 14 条路由的路由契约与端到端 HTTP 行为。
//
// 覆盖要点：14 条路由注册齐全；公开 5 条无需 token 即可访问；bind/unbind 无 token → 401；
// admin/* 用教师 token → 403；options 三形态；auth-url 的 400 分支与 redirect_uri 透传；
// wechat 已绑定（无 token）/ 未绑定 need_binding；wechat-work 免注册建号；third-party/login
// 的 400 分支与成功分支；bind-after-scan 的 200+error 与成功；bind 重复 422 / 成功；unbind；
// 管理员两条 contacts（异常 → 400、成功）与两条 import（统计与去重）。
//
// 假上游：企微接口地址由 WECHAT_WORK_API_BASE 指向 httptest.NewServer（服务构造时读环境变量，
// 故 env 必须在 router.New 之前设置）；**不访问真实外网**。
package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// routerUpstream 启动假上游：按路径返回预置 JSON（未预置的路径 404）。返回服务器与请求日志。
func routerUpstream(t *testing.T, responses map[string]string) (*httptest.Server, func() []string) {
	t.Helper()
	seen := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		body, ok := responses[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, func() []string { return seen }
}

// stringReader 把请求体字符串包装成 io.Reader（空串即空请求体）。
func stringReader(body string) *strings.Reader { return strings.NewReader(body) }

// bcryptHash 生成 bcrypt 哈希（建带已知密码的账号用）。
func bcryptHash(t *testing.T, plain string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

// thirdPartyRoutes 本批新增的 14 条路由（路径逐字对齐 Laravel routes/api.php 第 24-32、38-39、133-136 行）。
var thirdPartyRoutes = []string{
	"GET /api/v1/auth/third-party/auth-url",
	"GET /api/v1/auth/third-party/options",
	"POST /api/v1/auth/third-party/login",
	"POST /api/v1/auth/teacher/login/wechat",
	"POST /api/v1/auth/teacher/login/wechat-work",
	"POST /api/v1/auth/teacher/login/qq",
	"POST /api/v1/auth/teacher/login/renren",
	"POST /api/v1/auth/teacher/bind-after-scan",
	"POST /api/v1/auth/bind/:platform",
	"DELETE /api/v1/auth/unbind/:platform",
	"GET /api/v1/admin/wechat-work/contacts",
	"POST /api/v1/admin/wechat-work/import",
	"GET /api/v1/admin/third-party/contacts",
	"POST /api/v1/admin/third-party/import",
}

// thirdPartyPublicRoutes 本批 8 条公开路由（无需 token）。
var thirdPartyPublicRoutes = []string{
	"GET /api/v1/auth/third-party/auth-url",
	"GET /api/v1/auth/third-party/options",
	"POST /api/v1/auth/third-party/login",
	"POST /api/v1/auth/teacher/login/wechat",
	"POST /api/v1/auth/teacher/login/wechat-work",
	"POST /api/v1/auth/teacher/login/qq",
	"POST /api/v1/auth/teacher/login/renren",
	"POST /api/v1/auth/teacher/bind-after-scan",
}

// TestThirdPartyRoutesRegistered 14 条路由全部注册。
func TestThirdPartyRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range thirdPartyRoutes {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}
}

// TestThirdPartyPublicRoutesNeedNoToken 公开的 8 条无需 token（不会被鉴权中间件拦成 401）。
func TestThirdPartyPublicRoutesNeedNoToken(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	for _, route := range thirdPartyPublicRoutes {
		method, path := splitRouteKey(t, route)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(method, path, nil))
		assert.NotEqual(t, http.StatusUnauthorized, w.Code, "%s 是公开路由，不应 401", route)
	}
}

// TestThirdPartyProtectedRoutes 需登录的两条：无 token → 401；管理员 4 条用教师 token → 403。
func TestThirdPartyProtectedRoutes(t *testing.T) {
	f := newThirdPartyFixture(t, nil)

	// bind / unbind：无 token → 401。
	for _, route := range []string{"POST /api/v1/auth/bind/:platform", "DELETE /api/v1/auth/unbind/:platform"} {
		method, path := splitRouteKey(t, route)
		path = replaceIDParam(path)
		path = replacePlatformParam(path)
		w := httptest.NewRecorder()
		engine := f.Engine
		engine.ServeHTTP(w, httptest.NewRequest(method, path, nil))
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s 未登录应 401", route)
	}

	// 管理员 4 条：教师 token → 403。
	for _, route := range thirdPartyRoutes[10:] {
		method, path := splitRouteKey(t, route)
		resp := f.do(t, f.TeacherToken, method, path, "{}")
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "%s 教师身份应 403", route)
		resp.Body.Close()
	}

	// 教师 token 访问 admin 通讯录 → 403。
	resp := f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/admin/wechat-work/contacts", "")
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()
}

// thirdPartyFixture 一校 + 管理员 + 教师 + 一个班级 + 完整路由。
type thirdPartyFixture struct {
	DB           *gorm.DB
	Engine       *gin.Engine
	Client       *http.Client
	BaseURL      string
	AdminToken   string
	TeacherToken string
	School       models.School
	Teacher      models.User
	Class        models.ClassRoom
}

// newThirdPartyFixture 建内存库 + 完整路由；upstream 非空时把企微接口指向假上游（env 注入）。
func newThirdPartyFixture(t *testing.T, upstream *httptest.Server) *thirdPartyFixture {
	t.Helper()

	// 服务在 router.New 时读环境变量，故 env 必须先设置。
	t.Setenv("APP_URL", "https://learnstar.example.com")
	if upstream != nil {
		t.Setenv("WECHAT_WORK_API_BASE", upstream.URL)
		t.Setenv("WECHAT_WORK_CORPID", "corp-1")
		t.Setenv("WECHAT_WORK_AGENTID", "1000014")
		t.Setenv("WECHAT_WORK_SECRET", "sec-1")
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), database.GormConfig())
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.Migrate(db))

	school := models.School{Name: "测试学校", Code: "tp-school", Status: "active"}
	require.NoError(t, db.Create(&school).Error)

	admin := models.User{
		SchoolID: school.ID, Role: "school_admin", Username: "tp-admin", Name: "管理员", Status: "active",
	}
	require.NoError(t, db.Create(&admin).Error)
	teacher := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "tp-teacher", Name: "李老师", Status: "active",
	}
	require.NoError(t, db.Create(&teacher).Error)
	class := models.ClassRoom{SchoolID: school.ID, Name: "一年级（1）班", Grade: "一年级", Status: "active"}
	require.NoError(t, db.Create(&class).Error)

	jwtMgr := jwtauth.New("test-secret", 1)
	adminToken, err := jwtMgr.Generate(admin.ID, admin.Role, admin.SchoolID)
	require.NoError(t, err)
	teacherToken, err := jwtMgr.Generate(teacher.ID, teacher.Role, teacher.SchoolID)
	require.NoError(t, err)

	engine := router.New(db, &config.Config{JWTSecret: "test-secret", JWTExpHours: 1})
	server := httptest.NewServer(engine)
	t.Cleanup(server.Close)

	return &thirdPartyFixture{
		DB: db, Engine: engine, Client: server.Client(), BaseURL: server.URL,
		AdminToken: adminToken, TeacherToken: teacherToken,
		School: school, Teacher: teacher, Class: class,
	}
}

// do 发请求（token 为空表示不带 Authorization）。
func (f *thirdPartyFixture) do(t *testing.T, token, method, path, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, f.BaseURL+path, stringReader(body))
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := f.Client.Do(req)
	require.NoError(t, err)
	return resp
}

// replacePlatformParam 把 gin 路径参数 :platform 替换为具体平台。
func replacePlatformParam(path string) string {
	return strings.ReplaceAll(path, ":platform", "wechat")
}

// tpJSON 读取响应体为 map 并关闭。
func tpJSON(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	return decodeBody(t, resp)
}

// TestThirdPartyOptionsAndAuthURLHTTP options 三形态 + auth-url 的 400 / 成功 / redirect_uri 透传。
func TestThirdPartyOptionsAndAuthURLHTTP(t *testing.T) {
	// 企微 corp_id / agentid 参与 auth-url 拼接（无需 HTTP，故只要 env）。
	t.Setenv("WECHAT_WORK_CORPID", "corp-1")
	t.Setenv("WECHAT_WORK_AGENTID", "1000014")
	t.Setenv("WECHAT_WORK_SECRET", "sec-1")
	f := newThirdPartyFixture(t, nil)

	// options：未配置 → 默认三平台。
	resp := f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/options", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	options, ok := tpJSON(t, resp)["data"].([]any)
	require.True(t, ok)
	require.Len(t, options, 3)
	first := options[0].(map[string]any)
	assert.Equal(t, "wechat_work", first["key"])
	assert.Equal(t, "企业微信", first["label"])
	assert.Equal(t, "🏢", first["icon"])
	assert.Equal(t, "#2B7CE9", first["color"])

	// options：学校开关交集（顺序按 platforms()）。
	settings, err := f.School.WithSetting("enabled_third_party_platforms", []string{"wechat", "qq"})
	require.NoError(t, err)
	require.NoError(t, f.DB.Model(&models.School{}).Where("id = ?", f.School.ID).Update("settings", settings).Error)
	resp = f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/options", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	options, ok = tpJSON(t, resp)["data"].([]any)
	require.True(t, ok)
	require.Len(t, options, 2)
	assert.Equal(t, "wechat", options[0].(map[string]any)["key"])
	assert.Equal(t, "#07C160", options[0].(map[string]any)["color"])
	assert.Equal(t, "qq", options[1].(map[string]any)["key"])
	assert.Equal(t, "#12B7F5", options[1].(map[string]any)["color"])

	// auth-url：学校未配置平台 → 400 固定文案。
	resp = f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/auth-url", "")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "未配置第三方平台，请在后台学校设置中选择", tpJSON(t, resp)["message"])

	// 配置为企微 → 200，默认 redirect_uri 取 APP_URL + /login。
	platformSettings, err := f.School.WithSetting("third_party_platform", "wechat_work")
	require.NoError(t, err)
	require.NoError(t, f.DB.Model(&models.School{}).Where("id = ?", f.School.ID).Update("settings", platformSettings).Error)

	resp = f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/auth-url", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "wechat_work", data["platform"])
	authURL, ok := data["auth_url"].(string)
	require.True(t, ok)
	assert.Contains(t, authURL, "appid=corp-1")
	assert.Contains(t, authURL, "agentid=1000014")
	assert.Contains(t, authURL, "redirect_uri=https%3A%2F%2Flearnstar.example.com%2Flogin")
	assert.Contains(t, authURL, "state="+itoa(int(f.School.ID)))

	// redirect_uri 透传。
	resp = f.do(t, "", http.MethodGet,
		"/api/v1/auth/third-party/auth-url?redirect_uri=https%3A%2F%2Ffoo.example%2Fcb", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, data["auth_url"].(string), "redirect_uri=https%3A%2F%2Ffoo.example%2Fcb")

	// school_id 指定不存在的学校 → 400「系统尚未初始化」。
	resp = f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/auth-url?school_id=999999", "")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "系统尚未初始化", tpJSON(t, resp)["message"])
}

// TestAuthURLWithoutSchool 库里没有学校 → 400「系统尚未初始化」。
func TestAuthURLWithoutSchool(t *testing.T) {
	f := newThirdPartyFixture(t, nil)
	require.NoError(t, f.DB.Where("1 = 1").Delete(&models.School{}).Error)

	resp := f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/auth-url", "")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "系统尚未初始化", tpJSON(t, resp)["message"])

	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/third-party/login", `{"code":"c"}`)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "系统尚未初始化", tpJSON(t, resp)["message"])

	resp = f.do(t, "", http.MethodGet, "/api/v1/auth/third-party/options", "")
	require.Equal(t, http.StatusOK, resp.StatusCode, "options 无学校时仍回默认三平台")
}

// TestWechatLoginAndBindAfterScanHTTP 微信两条分支 + bind-after-scan（错密码 200+error、成功、一次性）。
func TestWechatLoginAndBindAfterScanHTTP(t *testing.T) {
	f := newThirdPartyFixture(t, nil)

	// 缺 openid → 422。
	resp := f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/wechat", `{}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	tpJSON(t, resp)

	// 未绑定 → need_binding + temp_token。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/wechat",
		`{"openid":"wx-1","nick":"微信昵称","avatar":"https://a/1.png"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "need_binding", data["status"])
	assert.Equal(t, "wx-1", data["openid"])
	assert.Nil(t, data["unionid"])
	tempToken, ok := data["temp_token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, tempToken)
	_, hasToken := data["token"]
	assert.False(t, hasToken, "need_binding 不带 token")

	// bind-after-scan：密码错 → HTTP 200 + status=error。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/bind-after-scan",
		`{"temp_token":"`+tempToken+`","username":"tp-teacher","password":"wrong","platform":"wechat","platform_id":"wx-1"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body := tpJSON(t, resp)
	data, ok = body["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "error", data["status"])
	assert.Equal(t, "账号或密码错误，请核对后重试", data["message"])

	// 缺字段 → 422。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/bind-after-scan", `{"temp_token":"x"}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	tpJSON(t, resp)

	// 正确密码 → 绑定成功（tp-teacher 无密码哈希，故先给它写一个已知密码）。
	require.NoError(t, f.DB.Model(&models.User{}).Where("id = ?", f.Teacher.ID).
		Update("password_hash", bcryptHash(t, "ls123456")).Error)

	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/bind-after-scan",
		`{"temp_token":"`+tempToken+`","username":"tp-teacher","password":"ls123456","platform":"wechat","platform_id":"ignored"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bound", data["status"])
	user, ok := data["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "李老师", user["name"])
	assert.Equal(t, "微信昵称", user["nickname"], "默认昵称（= 姓名）被第三方昵称覆盖")

	var binding models.ThirdPartyBinding
	require.NoError(t, f.DB.Where("user_id = ? AND platform = ?", f.Teacher.ID, "wechat").First(&binding).Error)
	assert.Equal(t, "wx-1", binding.PlatformID, "平台 ID 取自 temp_token 上下文")

	var leftover int64
	require.NoError(t, f.DB.Model(&models.TempBindingContext{}).
		Where("temp_token = ?", tempToken).Count(&leftover).Error)
	assert.Equal(t, int64(0), leftover, "temp_token 一次性消费")

	// 已绑定后再登录 → logged_in 且**不带 token**（Laravel 原样行为）。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/wechat", `{"openid":"wx-1"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "logged_in", data["status"])
	_, hasToken = data["token"]
	assert.False(t, hasToken, "已绑定分支不带 token")
	_, hasUser := data["user"]
	assert.True(t, hasUser)
}

// TestQQAndRenrenLoginHTTP QQ / 人人通未绑定分支的键集。
func TestQQAndRenrenLoginHTTP(t *testing.T) {
	f := newThirdPartyFixture(t, nil)

	resp := f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/qq", `{"openid":"qq-1"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "need_binding", data["status"])
	assert.Equal(t, "qq-1", data["openid"])
	require.NotEmpty(t, data["temp_token"])
	_, hasPlatformID := data["platform_id"]
	assert.False(t, hasPlatformID, "QQ 分支不带 platform_id")

	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/renren", `{"user_id":"rr-1"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "need_binding", data["status"])
	assert.Equal(t, "rr-1", data["platform_id"])

	// 必填缺失 → 422。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/renren", `{}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	tpJSON(t, resp)
}

// TestWechatWorkLoginHTTP 免注册建号（含 nick 缺省用 userid、明文密码、绑定行、token）。
func TestWechatWorkLoginHTTP(t *testing.T) {
	srv, _ := routerUpstream(t, map[string]string{
		"/cgi-bin/gettoken":         `{"errcode":0,"access_token":"AT-1","expires_in":7200}`,
		"/cgi-bin/user/getuserinfo": `{"errcode":0,"UserId":"ww-1"}`,
		"/cgi-bin/user/get":         `{"errcode":0,"name":"企微老师","mobile":"137","email":"ww@x.c","avatar":"https://a/ww.png"}`,
	})
	f := newThirdPartyFixture(t, srv)

	// 无 userid 且无 code → 400「无法获取企业微信用户身份」。
	resp := f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/wechat-work", `{}`)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "无法获取企业微信用户身份", tpJSON(t, resp)["message"])

	// 传 code → 走 getuserinfo 换 userid → 免注册建号。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/wechat-work", `{"code":"code-1"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "logged_in", data["status"])
	token, ok := data["token"].(string)
	require.True(t, ok, "未绑定分支必须返回 token")
	require.NotEmpty(t, token)
	user, ok := data["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "ww-1", user["name"], "nick 缺省 → 姓名取 userid")
	assert.Equal(t, "ww-1", user["username"])
	newUserID := uint(user["id"].(float64))

	var stored models.User
	require.NoError(t, f.DB.First(&stored, newUserID).Error)
	assert.Equal(t, "ls123456", stored.PlainPassword, "写明文密码")
	var binding models.ThirdPartyBinding
	require.NoError(t, f.DB.Where("platform = ? AND platform_id = ?", "wechat_work", "ww-1").First(&binding).Error)
	assert.Equal(t, newUserID, binding.UserID)

	// 新 token 可用（用 Laravel 存在的受保护端点探测；非 Laravel 的 `/auth/me` 已移除）。
	authProbe := f.do(t, token, http.MethodGet, "/api/v1/auth/bindings", "")
	require.Equal(t, http.StatusOK, authProbe.StatusCode)
	probeBody := tpJSON(t, authProbe)
	_, hasData := probeBody["data"]
	require.True(t, hasData, "受保护端点应返回 data")

	// 已绑定后（同一 userid，传 nick）→ 无 token。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/teacher/login/wechat-work",
		`{"userid":"ww-1","nick":"昵称老师"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "logged_in", data["status"])
	_, hasToken := data["token"]
	assert.False(t, hasToken, "已绑定分支不带 token")
}

// TestThirdPartyLoginHTTP code 必填 422 / 未配置平台 400 / 上游失败 400 / 成功分支。
func TestThirdPartyLoginHTTP(t *testing.T) {
	// 假上游：code=bad-code 时 getuserinfo 报错（触发「企业微信登录失败」）。
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_, _ = w.Write([]byte(`{"errcode":0,"access_token":"AT-1","expires_in":7200}`))
		case "/cgi-bin/user/getuserinfo":
			if r.URL.Query().Get("code") == "bad-code" {
				_, _ = w.Write([]byte(`{"errcode":40029,"errmsg":"invalid code"}`))
				return
			}
			_, _ = w.Write([]byte(`{"errcode":0,"UserId":"ww-9"}`))
		case "/cgi-bin/user/get":
			_, _ = w.Write([]byte(`{"errcode":0,"name":"第三方老师","mobile":"139","email":"tp@x.c","avatar":""}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	f := newThirdPartyFixture(t, srv)

	// 缺 code → 422。
	resp := f.do(t, "", http.MethodPost, "/api/v1/auth/third-party/login", `{}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	tpJSON(t, resp)

	// 未配置平台 → 400 + 异常文案。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/third-party/login", `{"code":"c"}`)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "未配置或未知的第三方平台", tpJSON(t, resp)["message"])

	// 配置为企微：上游 getUserByCode 失败（getuserinfo 报错）→ 400「企业微信登录失败」。
	settings, err := f.School.WithSetting("third_party_platform", "wechat_work")
	require.NoError(t, err)
	require.NoError(t, f.DB.Model(&models.School{}).Where("id = ?", f.School.ID).Update("settings", settings).Error)

	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/third-party/login", `{"code":"bad-code"}`)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "企业微信登录失败", tpJSON(t, resp)["message"])

	// 成功分支：新建教师账号 + token（state 解析到本校）。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/third-party/login",
		`{"code":"code-1","state":"`+itoa(int(f.School.ID))+`"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "logged_in", data["status"])
	require.NotEmpty(t, data["token"])
	user, ok := data["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "第三方老师", user["name"])

	var binding models.ThirdPartyBinding
	require.NoError(t, f.DB.Where("platform = ? AND platform_id = ?", "wechat_work", "ww-9").First(&binding).Error)
	assert.Equal(t, uint(user["id"].(float64)), binding.UserID)

	// 再次登录（同 platform_id）→ 已绑定，直接登录（无 token）。
	resp = f.do(t, "", http.MethodPost, "/api/v1/auth/third-party/login", `{"code":"code-1"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "logged_in", data["status"])
	_, hasToken := data["token"]
	assert.False(t, hasToken, "已绑定分支不带 token")
}

// TestBindAndUnbindHTTP 绑定成功 / 重复 422 / 解绑成功 / 缺参 422。
func TestBindAndUnbindHTTP(t *testing.T) {
	f := newThirdPartyFixture(t, nil)

	// 缺 platform_id → 422。
	resp := f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/auth/bind/wechat", `{}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Equal(t, "请求参数格式错误", tpJSON(t, resp)["message"])

	// 成功 → 「绑定成功」。
	resp = f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/auth/bind/wechat",
		`{"platform_id":"wx-bind-1","nick":"昵称","avatar":"https://a/1.png"}`)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "绑定成功", tpJSON(t, resp)["message"])

	var binding models.ThirdPartyBinding
	require.NoError(t, f.DB.Where("user_id = ? AND platform = ?", f.Teacher.ID, "wechat").First(&binding).Error)
	assert.Equal(t, "wx-bind-1", binding.PlatformID)
	assert.Equal(t, "昵称", binding.PlatformNick)
	require.NotNil(t, binding.VerifiedAt)

	// 同一 platform_id 再绑（任何用户）→ 422。
	resp = f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/auth/bind/wechat", `{"platform_id":"wx-bind-1"}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Equal(t, "该第三方账号已被其他用户绑定", tpJSON(t, resp)["message"])

	// 解绑 → 「解绑成功」；再次解绑仍成功。
	resp = f.do(t, f.TeacherToken, http.MethodDelete, "/api/v1/auth/unbind/wechat", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "解绑成功", tpJSON(t, resp)["message"])

	var count int64
	require.NoError(t, f.DB.Model(&models.ThirdPartyBinding{}).
		Where("user_id = ? AND platform = ?", f.Teacher.ID, "wechat").Count(&count).Error)
	assert.Equal(t, int64(0), count)

	resp = f.do(t, f.TeacherToken, http.MethodDelete, "/api/v1/auth/unbind/wechat", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	tpJSON(t, resp)

	// 解绑后 bindings 列表显示未绑定。
	resp = f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/auth/bindings", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	tpJSON(t, resp)
}

// TestAdminContactsHTTP 两条 contacts：成功（企微 / 第三方平台）与失败（未配置 → 400）。
func TestAdminContactsHTTP(t *testing.T) {
	srv, _ := routerUpstream(t, map[string]string{
		"/cgi-bin/gettoken":        `{"errcode":0,"access_token":"AT-1","expires_in":7200}`,
		"/cgi-bin/department/list": `{"errcode":0,"department":[{"id":1,"name":"总部","parentid":0},{"id":2,"name":"一年级组","parentid":1}]}`,
		"/cgi-bin/user/list":       `{"errcode":0,"userlist":[{"userid":"u1","name":"张三","mobile":"138","email":"a@b.c","position":"教师","department":[1,2]},{"userid":"u1","name":"重复"}]}`,
	})
	f := newThirdPartyFixture(t, srv)

	// 企微通讯录成功。
	resp := f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/wechat-work/contacts", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	departments, ok := data["departments"].([]any)
	require.True(t, ok)
	require.Len(t, departments, 2)
	members, ok := data["members"].([]any)
	require.True(t, ok)
	require.Len(t, members, 1, "按 userid 去重")
	member := members[0].(map[string]any)
	assert.Equal(t, "张三", member["name"])
	assert.Equal(t, []any{"总部", "一年级组"}, member["department_names"])
	_, hasPlatform := data["platform"]
	assert.False(t, hasPlatform, "企微通讯录不带 platform 字段")

	// 第三方通讯录：未配置平台 → 400。
	resp = f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/third-party/contacts", "")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "未配置或未知的第三方平台", tpJSON(t, resp)["message"])

	// 配置为企微 → 200 且附带 platform=wechat_work。
	settings, err := f.School.WithSetting("third_party_platform", "wechat_work")
	require.NoError(t, err)
	require.NoError(t, f.DB.Model(&models.School{}).Where("id = ?", f.School.ID).Update("settings", settings).Error)

	resp = f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/third-party/contacts", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok = tpJSON(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "wechat_work", data["platform"])
	require.Len(t, data["members"].([]any), 1)

	// 上游异常（企微未配置凭证）→ 400。「企微未配置」：清掉 corp_id 后重新建引擎。
	t.Setenv("WECHAT_WORK_CORPID", "")
	t.Setenv("WECHAT_WORK_SECRET", "")
	noCreds := newThirdPartyFixture(t, nil)
	require.NoError(t, noCreds.DB.Model(&models.WechatWorkToken{}).Where("1 = 1").Delete(&models.WechatWorkToken{}).Error)
	resp = noCreds.do(t, noCreds.AdminToken, http.MethodGet, "/api/v1/admin/wechat-work/contacts", "")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "企微未配置", tpJSON(t, resp)["message"])
}

// TestAdminImportHTTP 两条 import：教师 + 学生统计与去重（同 Laravel 文案）。
func TestAdminImportHTTP(t *testing.T) {
	f := newThirdPartyFixture(t, nil)

	// 先放一个同班同名学生，触发跳过。
	require.NoError(t, f.DB.Create(&models.Student{
		ClassID: f.Class.ID, Name: "已有学生", Status: "active",
	}).Error)

	body := `{"teachers":[{"name":"导入老师","mobile":"138-0000-0001","email":"t@x.c"},` +
		`{"name":"导入老师重复","mobile":"13800000001"}],` +
		`"students":[{"name":"导入同学","class_id":` + itoa(int(f.Class.ID)) + `,"gender":"女生"},` +
		`{"name":"已有学生","class_id":` + itoa(int(f.Class.ID)) + `}]}`

	for _, path := range []string{"/api/v1/admin/wechat-work/import", "/api/v1/admin/third-party/import"} {
		resp := f.do(t, f.AdminToken, http.MethodPost, path, body)
		require.Equal(t, http.StatusOK, resp.StatusCode, path)
		respBody := tpJSON(t, resp)
		assert.Contains(t, respBody["message"], "已导入", path)
		data, ok := respBody["data"].(map[string]any)
		require.True(t, ok, path)
		assert.NotNil(t, data["skipped_teachers"])
		assert.NotNil(t, data["skipped_students"])
		assert.NotNil(t, data["teacher_accounts"])
	}

	// 第二次调用（企微那条）：教师已存在 → 全部跳过；学生同名 → 跳过。
	resp := f.do(t, f.AdminToken, http.MethodPost, "/api/v1/admin/wechat-work/import", body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	respBody := tpJSON(t, resp)
	assert.Equal(t, "已导入 0 名教师、0 名学生，跳过已存在教师 1 名，跳过同名学生 2 名", respBody["message"])
	data, ok := respBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(0), data["created_teachers"])
	assert.Equal(t, float64(0), data["created_students"])
	require.Len(t, data["skipped_teachers"].([]any), 1)
	assert.Equal(t, "手机号已存在", data["skipped_teachers"].([]any)[0].(map[string]any)["reason"])

	// 校验失败 → 422 + errors（点号键）。
	resp = f.do(t, f.AdminToken, http.MethodPost, "/api/v1/admin/wechat-work/import",
		`{"teachers":[{"name":""}],"students":[{"name":"缺班级"}]}`)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	badBody := tpJSON(t, resp)
	assert.Equal(t, "参数错误", badBody["message"])
	errs, ok := badBody["errors"].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, errs, "teachers.0.name")
	assert.Contains(t, errs, "students.0.class_id")
}

// TestThirdPartyRouteCount 路由总数（engine.Routes()）回归断言：第三方批 14 条 + 后续补齐的
// `POST /admin/school`（Laravel `match(['put','post'],'school')`）、`DELETE /admin/teachers/:id`
// （Laravel `disableTeacher`）、`POST /admin/classes/:id/timetable`（Laravel match），
// 以及最后一批的 `common/score-categories` 与 `wechat-work/callback`（GET|POST）；
// 随后又删除了 4 条 Laravel 不存在的别名路由（`/auth/login`、`/auth/me`、`/teacher/classes`、
// `/admin/classes/:id/students`），当前共 211 条。
func TestThirdPartyRouteCount(t *testing.T) {
	db, engine, _ := newTestEngine(t)
	assert.Equal(t, 211, len(engine.Routes()), "与 Laravel 逐条对齐后共 211 条")

	// 本批两张新表按预期命名（README 中记录的表名）。
	assert.True(t, db.Migrator().HasTable("temp_binding_contexts"))
	assert.True(t, db.Migrator().HasTable("wechat_work_tokens"))
}
