// 登录端点限流：按「客户端 IP + 路由路径」的**固定窗口**计数。
//
// 对齐 Laravel `routes/api.php` 的 `throttle:N,1`（第 21-23 行与第 295 行）对外可见行为：
// 首次命中开窗，60 秒窗口内第 N+1 次请求 → 429 + `Retry-After: <剩余秒>` +
// `{"message":"Too Many Attempts."}`（Laravel 框架默认的 429 响应体与响应头）。
//
// 有意差异（逐条）：
//  1. 计数落在**进程内** map（与既有大屏登录失败计数同一约定），**仅单实例有效**；
//     多副本部署需换共享存储（Redis / 数据库表）才能保持全局限流语义。
//  2. Laravel 的 RateLimiter 是**滑动窗口**（逐次命中都写时间戳，窗口随最后命中顺延）；
//     Go 端为**固定窗口**（首次命中开窗，到点整体重置）。极端情形：第 59 秒打满 N 次时，
//     Go 端在开窗后第 60 秒立刻放行，Laravel 还要等到「最后一次命中 + 60 秒」。
//  3. 客户端 IP 取 `c.ClientIP()`（同 `display/login` 的登录日志口径），沿用 gin 默认的
//     可信代理配置（默认信任全部代理 → 会采用 `X-Forwarded-For` / `X-Real-IP` 的首段）。
//     Laravel 侧同样受 TrustProxies 配置影响，此处不额外收窄。
//  4. 不带 Laravel `ThrottleRequests` 附带的 `X-RateLimit-Limit/Remaining/Reset` 头（前端未使用）。
package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// throttleMessage Laravel 框架 429 的默认响应体文案。
const throttleMessage = "Too Many Attempts."

// throttleEntry 某个「IP + 路由」在当前窗口内的命中计数。
type throttleEntry struct {
	count int
	start time.Time
}

// throttleLimiter 固定窗口计数器（进程内）。
//
// 惰性清理：每经过一个窗口宽度就把已过期的键删掉，避免 map 无界增长
// （无定时任务，与既有 display 事件 / display_tokens 的惰性清理约定一致）。
type throttleLimiter struct {
	mu        sync.Mutex
	entries   map[string]*throttleEntry
	max       int
	window    time.Duration
	lastSweep time.Time
	now       func() time.Time
}

func newThrottleLimiter(max int, window time.Duration) *throttleLimiter {
	return &throttleLimiter{
		entries: map[string]*throttleEntry{},
		max:     max,
		window:  window,
		now:     time.Now,
	}
}

// allow 记一次命中；返回是否放行与「窗口剩余秒数」（仅在拦截时有意义）。
func (l *throttleLimiter) allow(key string) (bool, int) {
	// 非正参数视为不限流（防御性：避免 max<=0 把一切请求都拦掉）。
	if l.max <= 0 || l.window <= 0 {
		return true, 0
	}

	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.lastSweep) >= l.window {
		for k, e := range l.entries {
			if now.Sub(e.start) >= l.window {
				delete(l.entries, k)
			}
		}
		l.lastSweep = now
	}

	entry, ok := l.entries[key]
	if !ok || now.Sub(entry.start) >= l.window {
		l.entries[key] = &throttleEntry{count: 1, start: now}
		return true, 0
	}

	entry.count++
	if entry.count <= l.max {
		return true, 0
	}

	remaining := int(math.Ceil(entry.start.Add(l.window).Sub(now).Seconds()))
	if remaining < 1 {
		remaining = 1
	}
	return false, remaining
}

// Throttle 生成限流中间件：同一「客户端 IP + 路由路径」在 window 内最多放行 max 次请求。
//
// 与 Laravel `throttle:max,windowMinutes` 的用法一一对应（挂载点见 internal/router/router.go）：
// 每次调用 Throttle 都会新建一个独立计数器，故不同路由的计数互不影响；
// 键里再带上路由路径，同一个中间件实例挂到多条路由时也不会互相干扰。
func Throttle(max int, window time.Duration) gin.HandlerFunc {
	limiter := newThrottleLimiter(max, window)

	return func(c *gin.Context) {
		allowed, retryAfter := limiter.allow(c.ClientIP() + "|" + throttleRoutePath(c))
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": throttleMessage})
			return
		}
		c.Next()
	}
}

// ThrottlePrincipal 与 Throttle 同构，但计数主体优先取**已认证身份**而不是 IP：
// 教师/管理员 → `user:<id>`；教室端 → `class:<班级 id>`；两者都没有才退回 `ip:<ClientIP>`；
// 最后统一拼上路由路径。
//
// 为什么按主体：教室端设备通常共用一个出口 IP（学校 NAT），按 IP 计数会让一个班的操作吃掉
// 全校额度；按主体则「一个教师/一个班一个额度」，既能挡住脚本刷与重试风暴，又不影响
// 正常课堂节奏（阈值见 router.go 的 scoreWrite / moneyWrite）。
//
// 与 Throttle 相同的既有约束：计数在进程内，**仅单实例有效**；多副本部署需换共享存储。
// 挂载点必须位于认证中间件**之后**（否则取不到主体，退化为 IP 计数）。
func ThrottlePrincipal(max int, window time.Duration) gin.HandlerFunc {
	limiter := newThrottleLimiter(max, window)

	return func(c *gin.Context) {
		allowed, retryAfter := limiter.allow(throttlePrincipalKey(c) + "|" + throttleRoutePath(c))
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": throttleMessage})
			return
		}
		c.Next()
	}
}

// throttlePrincipalKey 取限流主体：已认证用户 > 班级码 token 所属班级 > 客户端 IP。
func throttlePrincipalKey(c *gin.Context) string {
	if u := CurrentUser(c); u != nil && u.ID != 0 {
		return "user:" + strconv.FormatUint(uint64(u.ID), 10)
	}
	if classID := DisplayClassID(c); classID != 0 {
		return "class:" + strconv.FormatUint(uint64(classID), 10)
	}
	return "ip:" + c.ClientIP()
}

// throttleRoutePath 取限流键里的「路由路径」：优先用命中路由的注册模式
// （`/api/v1/auth/teacher/login`，不含 query），未命中路由时退回原始 URL 路径。
func throttleRoutePath(c *gin.Context) string {
	if path := c.FullPath(); path != "" {
		return path
	}
	return c.Request.URL.Path
}
