// 登录别名（POST /auth/teacher/login、POST /auth/admin/login）测试。
//
// 覆盖：正确凭据 → 200 且结构与统一登录 POST /auth/login 一致；角色不符 / 密码错误 / 账号停用
// adminLoginWithCredentials 的统一文案）；入参缺失 → 422；签发的 JWT 可用。
package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/handlers"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const credentialLoginFailMessage = "账号或密码错误，请核对后重试"

// seedCredentialUser 创建带密码的账号（可指定角色与状态）。
func seedCredentialUser(t *testing.T, db *gorm.DB, schoolID uint, role, username, password, status string) models.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	user := models.User{
		SchoolID: schoolID, Role: role, Username: username,
		PasswordHash: string(hash), Name: username + "-姓名", Status: status,
	}
	require.NoError(t, db.Create(&user).Error)
	return user
}

// newLoginEngine 构造仅含两条角色登录路由的引擎（Laravel 无统一 /auth/login，故不再注册）。
func newLoginEngine(t *testing.T, db *gorm.DB) (*gin.Engine, *auth.Manager) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	jwtMgr := auth.New("test-secret", 72)
	h := handlers.New(db, jwtMgr)

	engine := gin.New()
	engine.POST("/api/v1/auth/teacher/login", h.TeacherLogin)
	engine.POST("/api/v1/auth/admin/login", h.AdminLogin)
	return engine, jwtMgr
}

// postLogin 发一次登录请求并返回响应。
func postLogin(t *testing.T, engine *gin.Engine, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

// 正确凭据：200，data = {token, user{id, username, name, role, school_id}}，信封 {data, message}。
func TestLoginAliasesSuccessShape(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedCredentialUser(t, db, school.ID, "teacher", "t-1", "secret123", "active")
	admin := seedCredentialUser(t, db, school.ID, "school_admin", "a-1", "secret123", "active")

	engine, jwtMgr := newLoginEngine(t, db)

	type payload struct {
		Data struct {
			Token string `json:"token"`
			User  struct {
				ID       uint   `json:"id"`
				Username string `json:"username"`
				Name     string `json:"name"`
				Role     string `json:"role"`
				SchoolID uint   `json:"school_id"`
			} `json:"user"`
		} `json:"data"`
		Message string `json:"message"`
	}
	keysOf := func(t *testing.T, raw json.RawMessage) []string {
		t.Helper()
		out := map[string]json.RawMessage{}
		require.NoError(t, json.Unmarshal(raw, &out))
		keys := make([]string, 0, len(out))
		for key := range out {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys
	}
	dataRawOf := func(t *testing.T, raw json.RawMessage) json.RawMessage {
		t.Helper()
		body := map[string]json.RawMessage{}
		require.NoError(t, json.Unmarshal(raw, &body))
		return body["data"]
	}
	userKeysOf := func(t *testing.T, raw json.RawMessage) []string {
		t.Helper()
		data := map[string]json.RawMessage{}
		require.NoError(t, json.Unmarshal(dataRawOf(t, raw), &data))
		return keysOf(t, data["user"])
	}

	teacherResp := postLogin(t, engine, "/api/v1/auth/teacher/login", `{"username":"t-1","password":"secret123"}`)
	require.Equal(t, http.StatusOK, teacherResp.Code, teacherResp.Body.String())

	adminResp := postLogin(t, engine, "/api/v1/auth/admin/login", `{"username":"a-1","password":"secret123"}`)
	require.Equal(t, http.StatusOK, adminResp.Code, adminResp.Body.String())

	// 顶层 / data / data.user 的键集合固定为 {data, message} 与 {token, user}。
	for _, resp := range []*httptest.ResponseRecorder{teacherResp, adminResp} {
		assert.Equal(t, []string{"data", "message"}, keysOf(t, resp.Body.Bytes()))
		assert.Equal(t, []string{"token", "user"}, keysOf(t, dataRawOf(t, resp.Body.Bytes())))
		userKeys := userKeysOf(t, resp.Body.Bytes())
		assert.Subset(t, userKeys, []string{"id", "username", "name", "role", "school_id"})
		assert.NotContains(t, userKeys, "password_hash", "密码哈希不对外输出")
	}

	var got payload
	require.NoError(t, json.Unmarshal(teacherResp.Body.Bytes(), &got))
	assert.Equal(t, "ok", got.Message)
	assert.NotEmpty(t, got.Data.Token)
	assert.Equal(t, teacher.ID, got.Data.User.ID)
	assert.Equal(t, "t-1", got.Data.User.Username)
	assert.Equal(t, "t-1-姓名", got.Data.User.Name)
	assert.Equal(t, "teacher", got.Data.User.Role)
	assert.Equal(t, school.ID, got.Data.User.SchoolID)

	var gotAdmin payload
	require.NoError(t, json.Unmarshal(adminResp.Body.Bytes(), &gotAdmin))
	assert.Equal(t, admin.ID, gotAdmin.Data.User.ID)
	assert.Equal(t, "school_admin", gotAdmin.Data.User.Role)

	// 签发的 token 可用且 claims 正确。
	claims, err := jwtMgr.Parse(got.Data.Token)
	require.NoError(t, err)
	assert.Equal(t, teacher.ID, claims.UserID)
	assert.Equal(t, "teacher", claims.Role)
	assert.Equal(t, school.ID, claims.SchoolID)

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.NotNil(t, reloaded.LastLoginAt, "登录成功写入 last_login_at")
}

// 角色不符：教师账号不能走 /auth/admin/login，管理员账号也不能走 /auth/teacher/login。
func TestLoginAliasesRejectWrongRole(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedCredentialUser(t, db, school.ID, "teacher", "t-2", "secret123", "active")
	seedCredentialUser(t, db, school.ID, "school_admin", "a-2", "secret123", "active")
	engine, _ := newLoginEngine(t, db)

	adminPathAsTeacher := postLogin(t, engine, "/api/v1/auth/admin/login", `{"username":"t-2","password":"secret123"}`)
	assert.Equal(t, http.StatusUnauthorized, adminPathAsTeacher.Code)
	assert.JSONEq(t, `{"message":"`+credentialLoginFailMessage+`"}`, adminPathAsTeacher.Body.String())

	teacherPathAsAdmin := postLogin(t, engine, "/api/v1/auth/teacher/login", `{"username":"a-2","password":"secret123"}`)
	assert.Equal(t, http.StatusUnauthorized, teacherPathAsAdmin.Code)
	assert.JSONEq(t, `{"message":"`+credentialLoginFailMessage+`"}`, teacherPathAsAdmin.Body.String())
}

// 密码错误 / 账号不存在 / 账号停用：一律 401 同一文案（不区分具体原因）。
func TestLoginAliasesFailureMessages(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedCredentialUser(t, db, school.ID, "teacher", "t-3", "secret123", "active")
	seedCredentialUser(t, db, school.ID, "teacher", "t-off", "secret123", "inactive")
	seedCredentialUser(t, db, school.ID, "school_admin", "a-3", "secret123", "active")
	seedCredentialUser(t, db, school.ID, "school_admin", "a-off", "secret123", "inactive")
	engine, _ := newLoginEngine(t, db)

	cases := []struct {
		path string
		body string
	}{
		{"/api/v1/auth/teacher/login", `{"username":"t-3","password":"wrong"}`},
		{"/api/v1/auth/teacher/login", `{"username":"nobody","password":"secret123"}`},
		{"/api/v1/auth/teacher/login", `{"username":"t-off","password":"secret123"}`},
		{"/api/v1/auth/admin/login", `{"username":"a-3","password":"wrong"}`},
		{"/api/v1/auth/admin/login", `{"username":"nobody","password":"secret123"}`},
		{"/api/v1/auth/admin/login", `{"username":"a-off","password":"secret123"}`},
	}
	for _, item := range cases {
		w := postLogin(t, engine, item.path, item.body)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s %s", item.path, item.body)
		assert.JSONEq(t, `{"message":"`+credentialLoginFailMessage+`"}`, w.Body.String())
	}
}

// 入参非法（缺 username / 缺 password / 空串）→ 422（Go 侧沿用统一登录的 bindJSON 文案）。
func TestLoginAliasesValidation(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedCredentialUser(t, db, school.ID, "teacher", "t-4", "secret123", "active")
	seedCredentialUser(t, db, school.ID, "school_admin", "a-4", "secret123", "active")
	engine, _ := newLoginEngine(t, db)

	for _, path := range []string{"/api/v1/auth/teacher/login", "/api/v1/auth/admin/login"} {
		for _, body := range []string{`{}`, `{"username":"t-4"}`, `{"password":"secret123"}`, `{"username":"","password":"secret123"}`} {
			w := postLogin(t, engine, path, body)
			assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "%s %s", path, body)
			assert.JSONEq(t, `{"message":"请求参数格式错误"}`, w.Body.String())
		}
	}
}
