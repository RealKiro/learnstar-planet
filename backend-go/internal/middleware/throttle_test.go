// 限流中间件单测：第 N+1 次 429（带 Retry-After 与 Laravel 默认文案）、窗口内第 N 次仍 200、
// 窗口过期自动恢复、不同 IP 互不影响、未挂限流的路由不受影响、同一实例挂多条路由时键含路径。
//
// 全部为进程内 httptest 调用，**不访问外网**。
package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newThrottleEngine 建一个只挂限流的测试引擎：/login 受限流（max/window），/other 不受限。
func newThrottleEngine(max int, window time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", middleware.Throttle(max, window), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	r.GET("/other", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return r
}

// postLoginFrom 从指定 IP 调一次 /login。
func postLoginFrom(r *gin.Engine, ip string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = ip + ":12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestThrottleBlocksAfterMax 第 N 次放行、第 N+1 次 429（Retry-After + Laravel 默认文案）。
func TestThrottleBlocksAfterMax(t *testing.T) {
	r := newThrottleEngine(3, time.Minute)

	for i := 1; i <= 3; i++ {
		w := postLoginFrom(r, "10.0.0.1")
		require.Equal(t, http.StatusOK, w.Code, "第 %d 次（≤ max）应放行", i)
		assert.Empty(t, w.Header().Get("Retry-After"))
	}

	blocked := postLoginFrom(r, "10.0.0.1")
	require.Equal(t, http.StatusTooManyRequests, blocked.Code)

	retryAfter := blocked.Header().Get("Retry-After")
	require.NotEmpty(t, retryAfter, "429 必须带 Retry-After 头")
	seconds, err := strconv.Atoi(retryAfter)
	require.NoError(t, err)
	assert.Greater(t, seconds, 0)
	assert.LessOrEqual(t, seconds, 60)

	var body map[string]any
	require.NoError(t, json.Unmarshal(blocked.Body.Bytes(), &body))
	assert.Equal(t, "Too Many Attempts.", body["message"])
	assert.Len(t, body, 1, "响应体只有 message（同 Laravel 框架默认 429）")

	// 继续请求仍然 429（窗口未过期）。
	assert.Equal(t, http.StatusTooManyRequests, postLoginFrom(r, "10.0.0.1").Code)
}

// TestThrottleWindowExpires 窗口过期后同一 IP 恢复放行。
func TestThrottleWindowExpires(t *testing.T) {
	r := newThrottleEngine(1, 60*time.Millisecond)

	require.Equal(t, http.StatusOK, postLoginFrom(r, "10.0.0.2").Code)
	require.Equal(t, http.StatusTooManyRequests, postLoginFrom(r, "10.0.0.2").Code)

	time.Sleep(80 * time.Millisecond)

	assert.Equal(t, http.StatusOK, postLoginFrom(r, "10.0.0.2").Code, "窗口过期后计数重置")
}

// TestThrottleSeparatesIPs 不同客户端 IP 各自计数。
func TestThrottleSeparatesIPs(t *testing.T) {
	r := newThrottleEngine(1, time.Minute)

	require.Equal(t, http.StatusOK, postLoginFrom(r, "10.0.0.3").Code)
	require.Equal(t, http.StatusTooManyRequests, postLoginFrom(r, "10.0.0.3").Code)

	assert.Equal(t, http.StatusOK, postLoginFrom(r, "10.0.0.4").Code, "另一个 IP 不受影响")
	assert.Equal(t, http.StatusOK, postLoginFrom(r, "10.0.0.5").Code)
}

// TestThrottleIgnoresUnthrottledRoutes 未挂限流的路由不受影响。
func TestThrottleIgnoresUnthrottledRoutes(t *testing.T) {
	r := newThrottleEngine(1, time.Minute)

	require.Equal(t, http.StatusTooManyRequests, func() int {
		postLoginFrom(r, "10.0.0.6")
		return postLoginFrom(r, "10.0.0.6").Code
	}())

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/other", nil)
		req.RemoteAddr = "10.0.0.6:12345"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code, "未挂限流的路由恒放行")
	}
}

// TestThrottleKeyIncludesRoutePath 同一个中间件实例挂到两条路由时，两条路由的计数互不影响。
func TestThrottleKeyIncludesRoutePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw := middleware.Throttle(1, time.Minute)
	handler := func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "ok"}) }
	r.POST("/a", mw, handler)
	r.POST("/b", mw, handler)

	call := func(path string) int {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.RemoteAddr = "10.0.0.7:12345"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	require.Equal(t, http.StatusOK, call("/a"))
	require.Equal(t, http.StatusTooManyRequests, call("/a"), "/a 已用尽额度")
	assert.Equal(t, http.StatusOK, call("/b"), "不同路由路径独立计数")
}

// TestThrottleConcurrent 并发下不 panic 且放行次数恰为 max（-race 下校验无数据竞争）。
func TestThrottleConcurrent(t *testing.T) {
	r := newThrottleEngine(5, time.Minute)

	const attempts = 40
	codes := make(chan int, attempts)
	for i := 0; i < attempts; i++ {
		go func() {
			codes <- postLoginFrom(r, "10.0.0.8").Code
		}()
	}

	allowed := 0
	for i := 0; i < attempts; i++ {
		if <-codes == http.StatusOK {
			allowed++
		}
	}
	assert.Equal(t, 5, allowed, "并发下放行次数仍等于 max")
}
