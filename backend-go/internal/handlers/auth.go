// 认证相关处理器。
package handlers

import (
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TeacherLogin 教师账号密码登录：POST /api/v1/auth/teacher/login。
//
// 忠实移植 Laravel AuthController::teacherLoginWithCredentials（入参 username/password，
// 仅 role=teacher 且 status=active 可通过；失败一律 401「账号或密码错误，请核对后重试」）。
// 返回结构与既有统一登录 POST /auth/login 一致（data = {token, user}）。
//
// 有意差异：
//  1. Laravel 该路由带 throttle:6,1；Go 端已按同参数挂 `middleware.Throttle`（见 internal/router/router.go）。
//  2. 入参校验失败：Laravel 走 validate() 返回 422 + errors；Go 端沿用 bindJSON 的 422「请求参数格式错误」。
func (h *Handlers) TeacherLogin(c *gin.Context) {
	var req loginReq
	if !bindJSON(c, &req) {
		return
	}

	token, user, err := h.auth.TeacherLogin(req.Username, req.Password)
	if err != nil {
		fail(c, err)
		return
	}

	ok(c, gin.H{"token": token, "user": user})
}

// AdminLogin 管理员账号密码登录：POST /api/v1/auth/admin/login。
//
// 忠实移植 Laravel AuthController::adminLoginWithCredentials（仅 role=school_admin 且 status=active
// 可通过；教师账号走此接口一定 401，反之亦然）。限流已按 Laravel 的 throttle:6,1 挂载；
// 校验失败文案为 Go 统一的 422「请求参数格式错误」。
func (h *Handlers) AdminLogin(c *gin.Context) {
	var req loginReq
	if !bindJSON(c, &req) {
		return
	}

	token, user, err := h.auth.AdminLogin(req.Username, req.Password)
	if err != nil {
		fail(c, err)
		return
	}

	ok(c, gin.H{"token": token, "user": user})
}

// ChangePassword 修改当前登录用户密码。
func (h *Handlers) ChangePassword(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	if err := h.auth.ChangePassword(u, req.OldPassword, req.NewPassword); err != nil {
		fail(c, err)
		return
	}

	okMessage(c, "密码已修改")
}

// ClassLogin 班级码登录（学生端 / 班级大屏统一入口）：POST /api/v1/auth/class/login。
// 参数 class_code（required|string|max:10）；响应逐字同 Laravel —— 只有 data、没有 message 字段。
func (h *Handlers) ClassLogin(c *gin.Context) {
	var req struct {
		ClassCode string `json:"class_code" binding:"required,max=10"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.display.ClassLogin(req.ClassCode)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Logout 登出：撤销当前访问令牌（POST /api/v1/auth/logout）。
// 移植自 Laravel AuthController::logout（`$user->currentAccessToken()->delete()`）。
// 差异说明：Laravel 是删除 Sanctum 令牌行；Go 端把该令牌的 jti 写入 revoked_tokens，
// 中间件随即对旧 token 返回 401「登录已过期或无效」（无 jti 的历史令牌不受影响）。
func (h *Handlers) Logout(c *gin.Context) {
	claims := middleware.CurrentClaims(c)
	jti, exp := "", claims.Expiry()
	if claims != nil {
		jti = claims.JTI()
	}
	if err := h.auth.Logout(jti, exp); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "已登出")
}

// RefreshToken 刷新令牌：撤销旧令牌并签发新令牌（POST /api/v1/auth/refresh）。
// 响应 `{data:{token}}`（同 Laravel refreshToken）。
func (h *Handlers) RefreshToken(c *gin.Context) {
	u := middleware.CurrentUser(c)
	claims := middleware.CurrentClaims(c)
	jti, exp := "", claims.Expiry()
	if claims != nil {
		jti = claims.JTI()
	}

	token, err := h.auth.Refresh(u, jti, exp)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"token": token})
}

// GetBindings 第三方绑定列表（GET /api/v1/auth/bindings）。
// 仅返回学校后台启用的平台；未配置/空数组时默认 [企业微信, 微信, QQ]（同 Laravel getBindings）。
func (h *Handlers) GetBindings(c *gin.Context) {
	u := middleware.CurrentUser(c)
	views, err := h.auth.Bindings(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, views)
}
