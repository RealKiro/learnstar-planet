// 登录端点限流的**路由级**回归：只有 Laravel 带 `throttle:N,1` 的 4 条路由受限流。
//
// Laravel 权威来源 routes/api.php：
//
//	POST /api/v1/auth/teacher/login → throttle:6,1（第 21 行）
//	POST /api/v1/auth/admin/login   → throttle:6,1（第 22 行）
//	POST /api/v1/auth/class/login   → throttle:10,1（第 23 行）
//	POST /api/v1/display/login      → throttle:10,1（第 295 行）
//
// 全部为进程内 httptest 调用，**不访问外网**。
package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// throttleCall 一次调用的结果（状态码 / 响应头 / 解析后的 JSON 体）。
type throttleCall struct {
	code   int
	header http.Header
	body   map[string]any
}

// callJSON 以指定客户端 IP 发一次 POST JSON（RemoteAddr 即限流键里的客户端 IP）。
func callJSON(t *testing.T, engine *gin.Engine, path, body, ip string) throttleCall {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = ip + ":54321"
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	out := throttleCall{code: w.Code, header: w.Header(), body: map[string]any{}}
	if err := json.Unmarshal(w.Body.Bytes(), &out.body); err != nil {
		out.body = map[string]any{}
	}
	return out
}

// assertThrottled 断言响应是限流拦截：429 + Retry-After + Laravel 框架默认文案。
func assertThrottled(t *testing.T, call throttleCall) {
	t.Helper()
	require.Equal(t, http.StatusTooManyRequests, call.code)
	assert.Equal(t, "Too Many Attempts.", call.body["message"])

	retry := call.header.Get("Retry-After")
	require.NotEmpty(t, retry, "429 必须带 Retry-After 头")
	seconds, err := strconv.Atoi(retry)
	require.NoError(t, err)
	assert.Greater(t, seconds, 0)
	assert.LessOrEqual(t, seconds, 60)
}

// TestTeacherLoginThrottle 第 6 次放行、第 7 次 429（throttle:6,1）。
func TestTeacherLoginThrottle(t *testing.T) {
	_, engine, _ := newTestEngine(t)
	const path = "/api/v1/auth/teacher/login"
	const body = `{"username":"nobody","password":"wrong"}`

	for i := 1; i <= 6; i++ {
		call := callJSON(t, engine, path, body, "10.1.0.1")
		require.NotEqual(t, http.StatusTooManyRequests, call.code, "第 %d 次（≤ 6）不应被限流", i)
	}
	assertThrottled(t, callJSON(t, engine, path, body, "10.1.0.1"))
}

// TestAdminLoginThrottle 第 6 次放行、第 7 次 429（throttle:6,1）。
func TestAdminLoginThrottle(t *testing.T) {
	_, engine, _ := newTestEngine(t)
	const path = "/api/v1/auth/admin/login"
	const body = `{"username":"nobody","password":"wrong"}`

	for i := 1; i <= 6; i++ {
		call := callJSON(t, engine, path, body, "10.1.0.2")
		require.NotEqual(t, http.StatusTooManyRequests, call.code, "第 %d 次（≤ 6）不应被限流", i)
	}
	assertThrottled(t, callJSON(t, engine, path, body, "10.1.0.2"))
}

// TestClassLoginThrottle 第 10 次放行、第 11 次 429（throttle:10,1）。
func TestClassLoginThrottle(t *testing.T) {
	_, engine, _ := newTestEngine(t)
	const path = "/api/v1/auth/class/login"

	for i := 1; i <= 10; i++ {
		call := callJSON(t, engine, path, `{"class_code":"LX`+strconv.Itoa(i)+`"}`, "10.1.0.3")
		require.Equal(t, http.StatusUnauthorized, call.code,
			"无效班级码 → 401「班级码无效，请核对后重试」；第 %d 次不应被限流", i)
	}
	assertThrottled(t, callJSON(t, engine, path, `{"class_code":"LX99"}`, "10.1.0.3"))
}

// TestDisplayLoginThrottle 第 10 次放行、第 11 次 429（throttle:10,1）。
//
// 必须每次换班级码：`display/login` 另有「同班级码连续失败 5 次锁 15 分钟」的**独立**机制
// （DisplayService.displayLoginFail），换码才能把两者区分开——两条机制在 Laravel 里本就并存。
func TestDisplayLoginThrottle(t *testing.T) {
	_, engine, _ := newTestEngine(t)
	const path = "/api/v1/display/login"

	for i := 1; i <= 10; i++ {
		call := callJSON(t, engine, path, `{"code":"LX`+strconv.Itoa(i)+`"}`, "10.1.0.4")
		require.Equal(t, http.StatusNotFound, call.code,
			"无效班级码 → 404；第 %d 次不应被限流", i)
	}
	assertThrottled(t, callJSON(t, engine, path, `{"code":"LX99"}`, "10.1.0.4"))
}

// TestLoginThrottleIsPerIP 不同客户端 IP 各自计数（限流键含 IP）。
func TestLoginThrottleIsPerIP(t *testing.T) {
	_, engine, _ := newTestEngine(t)
	const path = "/api/v1/auth/teacher/login"
	const body = `{"username":"nobody","password":"wrong"}`

	for i := 1; i <= 6; i++ {
		require.NotEqual(t, http.StatusTooManyRequests, callJSON(t, engine, path, body, "10.1.0.5").code)
	}
	require.Equal(t, http.StatusTooManyRequests, callJSON(t, engine, path, body, "10.1.0.5").code)

	// 换一个 IP：不受前一个 IP 的计数影响。
	assert.NotEqual(t, http.StatusTooManyRequests,
		callJSON(t, engine, path, body, "10.1.0.6").code, "另一个 IP 的计数独立")
}

// TestUnthrottledRouteNotAffected 未挂限流的路由不受影响（`auth/third-party/login` 在 Laravel 同样没有 throttle）。
func TestUnthrottledRouteNotAffected(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	for i := 0; i < 12; i++ {
		call := callJSON(t, engine, "/api/v1/auth/third-party/login",
			`{"platform":"wechat_work","code":"x"}`, "10.1.0.7")
		require.NotEqual(t, http.StatusTooManyRequests, call.code, "auth/third-party/login 未挂限流")
	}
}

// TestRouteCountUnchangedWithThrottle 新中间件不改变路由条数（与 Laravel 对齐后为 211）。
func TestRouteCountUnchangedWithThrottle(t *testing.T) {
	_, engine, _ := newTestEngine(t)
	assert.Equal(t, 211, len(engine.Routes()), "限流中间件不新增路由")
}
