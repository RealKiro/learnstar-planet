package router_test

// SPA 静态服务与前端路由兜底的端到端测试（不依赖任何业务表）。
//
// 覆盖：index.html 兜底、静态文件直出、/api/* JSON 404、上传文件映射
// （/storage/app/uploads/* → UPLOAD_DIR）、目录穿越防护、PUBLIC_DIR 未配置时
// 优雅 404。

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// newSPAEngine 构建带 PUBLIC_DIR 的完整引擎（内存库；NoRoute 路径不触库）。
func newSPAEngine(t *testing.T, publicDir string) *gin.Engine {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	return router.New(db, &config.Config{JWTSecret: "test-secret", JWTExpHours: 1, PublicDir: publicDir})
}

func TestSPAFallbackServesIndexAndAssets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "assets"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>learnstar-spa</html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("console.log(1)"), 0o644))

	e := newSPAEngine(t, dir)

	// 首页
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "learnstar-spa")

	// Vue Router history 模式路径回落 index.html
	w = httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/teacher/dashboard", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "learnstar-spa")

	// 存在的静态资源直出
	w = httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/assets/app.js", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "console.log(1)")

	// /api/* 未命中 → JSON 404
	w = httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/not-exist", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message":"资源不存在"}`, w.Body.String())

	// 非 GET 也走 JSON 404
	w = httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/teacher/dashboard", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSPAFallbackWithoutPublicDir(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := newSPAEngine(t, "") // 纯 API 部署：不配置 PUBLIC_DIR

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message":"资源不存在"}`, w.Body.String())
}

func TestSPAFallbackServesUploads(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uploadDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(uploadDir, "schools"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(uploadDir, "schools", "logo.png"), []byte("PNGDATA"), 0o644))
	t.Setenv("UPLOAD_DIR", uploadDir)

	// PUBLIC_DIR 指向另一个空目录（不存在 index.html）
	e := newSPAEngine(t, filepath.Join(uploadDir, "no-public"))

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/storage/app/uploads/schools/logo.png", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "PNGDATA")

	// 目录穿越不得逃出 UPLOAD_DIR
	w = httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/storage/app/uploads/../../etc/passwd", nil))
	assert.NotEqual(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), "root:")
}
