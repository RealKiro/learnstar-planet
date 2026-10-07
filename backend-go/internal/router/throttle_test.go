// 写端点限流的守护测试：确认「按主体计数」的限流真的会触发 429，且不误伤鉴权。
//
// 背景：登录以外的大量写端点此前**完全没有限流**（Laravel 也只在 login / class-login 挂了 throttle），
// 重试风暴与脚本刷没有任何兜底。现在 teacher / display 两个分组各挂一个按主体
// （教师 user:<id> / 班级 class:<id>）的宽松限流，见 middleware.ThrottlePrincipal 与该处的注释。
//
// 为什么按主体而不是 IP：教室端设备通常共用一个出口 IP（学校 NAT），按 IP 计数会让一个班的
// 操作吃掉全校额度。
package router_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTeacherGroupThrottleBurst teacher 分组：同一教师的高频请求最终会被 429 拦住。
//
// 用一个几乎无副作用的写路由（不存在的兑换 id → 404）来打，只关心是否会被限流；
// 循环上限远大于阈值，且不假设此前是否已有计数，故不会因限流窗口滚动而偶发失败。
func TestTeacherGroupThrottleBurst(t *testing.T) {
	f := newPeopleFixture(t)

	const path = "/api/v1/teacher/shop/redemptions/999999/approve"
	got429 := false
	for i := 0; i < 900 && !got429; i++ {
		resp := f.do(t, f.TeacherToken, http.MethodPut, path, "")
		require.NotNil(t, resp)
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			got429 = true
			assert.NotEmpty(t, resp.Header.Get("Retry-After"), "429 应带 Retry-After")
		case http.StatusUnauthorized, http.StatusForbidden:
			t.Fatalf("限流不应影响鉴权，第 %d 次得到 %d", i+1, resp.StatusCode)
		}
		_ = resp.Body.Close()
	}
	assert.True(t, got429, "同一教师的高频写请求应被限流（阈值 600/分钟）")
}

// TestDisplayGroupThrottleByClass 教室端按「班级」计数：同一班级码打满后 429。
func TestDisplayGroupThrottleByClass(t *testing.T) {
	f := newPeopleFixture(t)

	login := f.do(t, "", http.MethodPost, "/api/v1/auth/class/login", `{"class_code":"LS11"}`)
	require.Equal(t, http.StatusOK, login.StatusCode, "班级码登录应成功")
	inner, ok := decodeBody(t, login)["data"].(map[string]any)
	require.True(t, ok)
	token, ok := inner["token"].(string)
	require.True(t, ok, "班级码登录应返回 token")
	require.NotEmpty(t, token)

	got429 := false
	for i := 0; i < 1600 && !got429; i++ {
		resp := f.do(t, token, http.MethodGet, "/api/v1/display/students", "")
		require.NotNil(t, resp)
		if resp.StatusCode == http.StatusTooManyRequests {
			got429 = true
			break
		}
		require.Equal(t, http.StatusOK, resp.StatusCode, "限流前应正常返回（第 %d 次）", i+1)
		_ = resp.Body.Close()
	}
	assert.True(t, got429, "同一班级码的高频请求应被限流（阈值 1200/分钟）")
}
