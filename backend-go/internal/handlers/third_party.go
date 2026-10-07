// 第三方扫码登录 / 绑定 / 管理端通讯录处理器。
//
// 对应 Laravel：
//   - AuthController::thirdPartyAuthUrl / thirdPartyOptions / thirdPartyLogin /
//     teacherLoginWithWechat / teacherLoginWithWechatWork / teacherLoginWithQQ /
//     teacherLoginWithRenren / bindAfterScan / bindThirdParty / unbindThirdParty
//   - SchoolAdminController::wechatWorkContacts / importWechatWorkUsers /
//     thirdPartyContacts / importThirdPartyUsers
//
// 响应结构与文案逐项对齐；差异见 services 层各文件头部说明与 backend-go/README.md。
package handlers

import (
	"fmt"
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// ============================================================
// 公开：登录页平台选项 / 授权 URL / 扫码登录
// ============================================================

// ThirdPartyAuthURL 获取当前学校第三方平台的扫码授权 URL（GET /api/v1/auth/third-party/auth-url）。
//
// `school_id` 可选（未传取第一所学校）；无学校 → 400「系统尚未初始化」；
// 学校未配置平台或平台未知 → 400「未配置第三方平台，请在后台学校设置中选择」；
// `redirect_uri` 缺省用 `url('/login')`（Go 端读 APP_URL，见 services.LoginURL）。
func (h *Handlers) ThirdPartyAuthURL(c *gin.Context) {
	school, err := h.thirdParty.SchoolForAuthURL(queryID(c, "school_id"))
	if err != nil {
		fail(c, err)
		return
	}
	if school == nil {
		fail(c, services.ErrBadRequest("系统尚未初始化"))
		return
	}

	redirectURI, present := c.GetQuery("redirect_uri")
	if !present {
		redirectURI = services.LoginURL()
	}

	result, err := h.thirdParty.AuthURLOfSchool(school, redirectURI)
	if err != nil {
		fail(c, services.ErrBadRequest("未配置第三方平台，请在后台学校设置中选择"))
		return
	}
	ok(c, result)
}

// ThirdPartyOptions 登录页可选的第三方平台列表（GET /api/v1/auth/third-party/options）。
func (h *Handlers) ThirdPartyOptions(c *gin.Context) {
	options, err := h.thirdParty.Options()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, options)
}

// ThirdPartyLogin 第三方平台扫码登录回调（POST /api/v1/auth/third-party/login）。
//
// 顺序逐条同 Laravel：解析 state → 取学校（无学校 400「系统尚未初始化」）→ providerFor
// （失败 400 + 异常文案）→ getUserByCode（失败 400 + 异常文案）→ loginWithThirdParty
// （status != logged_in → 400 + result.message 或「登录失败」）。
func (h *Handlers) ThirdPartyLogin(c *gin.Context) {
	var req struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state"`
	}
	if !bindJSON(c, &req) {
		return
	}

	school, err := h.thirdParty.ResolveSchoolFromState(req.State)
	if err != nil {
		fail(c, err)
		return
	}
	if school == nil {
		fail(c, services.ErrBadRequest("系统尚未初始化"))
		return
	}

	provider, err := h.thirdParty.Manager().ProviderFor(school.SettingString("third_party_platform", ""))
	if err != nil {
		fail(c, services.ErrBadRequest(err.Error()))
		return
	}
	userInfo, err := provider.GetUserByCode(school.ID, req.Code)
	if err != nil {
		fail(c, services.ErrBadRequest(err.Error()))
		return
	}

	result, err := h.thirdParty.LoginWithThirdParty(provider.Key(), userInfo, school)
	if err != nil {
		fail(c, err)
		return
	}
	if result.Status != "logged_in" {
		message := result.Message
		if message == "" {
			message = "登录失败"
		}
		fail(c, services.ErrBadRequest(message))
		return
	}
	ok(c, result)
}

// TeacherLoginWithWechat 微信扫码登录（POST /api/v1/auth/teacher/login/wechat）。
//
// ⚠️ 已绑定时只返回 `{status:'logged_in', user}`，**不带 token**（Laravel 原样行为，照抄）。
func (h *Handlers) TeacherLoginWithWechat(c *gin.Context) {
	var req struct {
		OpenID  string  `json:"openid" binding:"required"`
		UnionID *string `json:"unionid"`
		Nick    *string `json:"nick"`
		Avatar  *string `json:"avatar"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !checkNullableLens(c, nullableLen{"nick", req.Nick, 80}, nullableLen{"avatar", req.Avatar, 500}) {
		return
	}

	result, err := h.thirdParty.LoginWithWechat(req.OpenID, req.UnionID, req.Nick, req.Avatar)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// TeacherLoginWithWechatWork 企业微信扫码登录（POST /api/v1/auth/teacher/login/wechat-work）。
//
// 无 `userid` 且有 `code` 时先用 code 换 userid（学校取自 state 解析结果）；
// 仍无 userid → 400「无法获取企业微信用户身份」；未绑定账号时免注册自动建教师账号并返回 token。
func (h *Handlers) TeacherLoginWithWechatWork(c *gin.Context) {
	var req struct {
		UserID string  `json:"userid"`
		Code   string  `json:"code"`
		State  string  `json:"state"`
		Nick   *string `json:"nick"`
		Avatar *string `json:"avatar"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !checkNullableLens(c, nullableLen{"nick", req.Nick, 80}, nullableLen{"avatar", req.Avatar, 500}) {
		return
	}

	school, err := h.thirdParty.ResolveSchoolFromState(req.State)
	if err != nil {
		fail(c, err)
		return
	}

	userid := req.UserID
	if userid == "" && req.Code != "" && school != nil {
		// Laravel 不捕获这里的异常（企微未配置 / token 失败 → 500）；Go 端按同一口径向上抛。
		userid, err = h.thirdParty.WechatWork().GetUserIDByCode(school.ID, req.Code)
		if err != nil {
			fail(c, err)
			return
		}
	}
	if userid == "" {
		fail(c, services.ErrBadRequest("无法获取企业微信用户身份"))
		return
	}

	result, err := h.thirdParty.LoginWithWechatWork(userid, req.Nick, req.Avatar, school)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// TeacherLoginWithQQ QQ 扫码登录（POST /api/v1/auth/teacher/login/qq）。
// ⚠️ 已绑定时只返回 `{status:'logged_in', user}`，**不带 token**（Laravel 原样行为）。
func (h *Handlers) TeacherLoginWithQQ(c *gin.Context) {
	var req struct {
		OpenID string  `json:"openid" binding:"required"`
		Nick   *string `json:"nick"`
		Avatar *string `json:"avatar"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !checkNullableLens(c, nullableLen{"nick", req.Nick, 80}, nullableLen{"avatar", req.Avatar, 500}) {
		return
	}

	result, err := h.thirdParty.LoginWithQQ(req.OpenID, req.Nick, req.Avatar)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// TeacherLoginWithRenren 人人通空间登录（POST /api/v1/auth/teacher/login/renren）。
// ⚠️ 已绑定时只返回 `{status:'logged_in', user}`，**不带 token**（Laravel 原样行为）。
func (h *Handlers) TeacherLoginWithRenren(c *gin.Context) {
	var req struct {
		UserID string  `json:"user_id" binding:"required"`
		Nick   *string `json:"nick"`
		Avatar *string `json:"avatar"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !checkNullableLens(c, nullableLen{"nick", req.Nick, 80}, nullableLen{"avatar", req.Avatar, 500}) {
		return
	}

	result, err := h.thirdParty.LoginWithRenren(req.UserID, req.Nick, req.Avatar)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// BindAfterScan 扫码后绑定已有教师账号（POST /api/v1/auth/teacher/bind-after-scan）。
//
// 账号密码错误 → **HTTP 200** + `{data:{status:'error', message:'账号或密码错误，请核对后重试'}}`
// （Laravel 原样行为）；成功 → `{data:{status:'bound', user}}`。
func (h *Handlers) BindAfterScan(c *gin.Context) {
	var req struct {
		TempToken  string  `json:"temp_token" binding:"required"`
		Username   string  `json:"username" binding:"required"`
		Password   string  `json:"password" binding:"required"`
		Platform   string  `json:"platform" binding:"required"`
		PlatformID string  `json:"platform_id" binding:"required"`
		UnionID    *string `json:"unionid"`
		Nick       *string `json:"nick"`
		Avatar     *string `json:"avatar"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !checkNullableLens(c, nullableLen{"nick", req.Nick, 80}, nullableLen{"avatar", req.Avatar, 500}) {
		return
	}

	result, err := h.thirdParty.BindAfterScan(
		req.TempToken, req.Username, req.Password, req.Platform, req.PlatformID,
		req.UnionID, req.Nick, req.Avatar,
	)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// ============================================================
// 需登录：绑定 / 解绑
// ============================================================

// BindThirdParty 登录后主动绑定第三方账号（POST /api/v1/auth/bind/:platform）。
// 该 platform + platform_id 已被任意用户绑定 → 422「该第三方账号已被其他用户绑定」；成功「绑定成功」。
func (h *Handlers) BindThirdParty(c *gin.Context) {
	var req struct {
		PlatformID string  `json:"platform_id" binding:"required"`
		Nick       *string `json:"nick"`
		Avatar     *string `json:"avatar"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !checkNullableLens(c, nullableLen{"nick", req.Nick, 80}, nullableLen{"avatar", req.Avatar, 500}) {
		return
	}

	user := middleware.CurrentUser(c)
	if err := h.thirdParty.BindForUser(user, c.Param("platform"), req.PlatformID, req.Nick, req.Avatar); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "绑定成功")
}

// UnbindThirdParty 解绑当前用户的第三方账号（DELETE /api/v1/auth/unbind/:platform）。
func (h *Handlers) UnbindThirdParty(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if err := h.thirdParty.UnbindForUser(user, c.Param("platform")); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "解绑成功")
}

// ============================================================
// 管理端：通讯录拉取 / 导入
// ============================================================

// AdminWechatWorkContacts 拉取企业微信通讯录（GET /api/v1/admin/wechat-work/contacts）。
func (h *Handlers) AdminWechatWorkContacts(c *gin.Context) {
	user := middleware.CurrentUser(c)
	contacts, err := h.thirdParty.WechatWorkContacts(user.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, contacts)
}

// AdminWechatWorkImport 从企业微信通讯录批量导入教师与学生（POST /api/v1/admin/wechat-work/import）。
func (h *Handlers) AdminWechatWorkImport(c *gin.Context) {
	h.adminImportContacts(c)
}

// AdminThirdPartyContacts 拉取当前学校所选第三方平台的通讯录（GET /api/v1/admin/third-party/contacts）。
func (h *Handlers) AdminThirdPartyContacts(c *gin.Context) {
	user := middleware.CurrentUser(c)
	contacts, err := h.thirdParty.ThirdPartyContacts(user.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, contacts)
}

// AdminThirdPartyImport 从当前学校所选第三方平台批量导入教师与学生（POST /api/v1/admin/third-party/import）。
func (h *Handlers) AdminThirdPartyImport(c *gin.Context) {
	h.adminImportContacts(c)
}

// adminImportContacts 两条导入路由共用的请求解析与响应（同 Laravel importThirdPartyUsers 委派）。
func (h *Handlers) adminImportContacts(c *gin.Context) {
	user := middleware.CurrentUser(c)

	var req services.ContactsImportRequest
	if !readJSONBody(c, &req) {
		return
	}

	result, err := h.thirdParty.ImportContacts(user.SchoolID, req)
	if err != nil {
		fail(c, err)
		return
	}
	// Laravel 顶层结构：{message, data}（不带 message:"ok" 信封）。
	c.JSON(http.StatusOK, result)
}

// nullableLen 一个可选字符串字段的长度上限（同 Laravel `nullable|string|max:N`）。
type nullableLen struct {
	Key   string
	Value *string
	Max   int
}

// checkNullableLens 校验可选字符串字段长度；超长时写 422 + errors 并返回 false。
func checkNullableLens(c *gin.Context, fields ...nullableLen) bool {
	errs := map[string][]string{}
	for _, f := range fields {
		if f.Value != nil && len([]rune(*f.Value)) > f.Max {
			errs[f.Key] = []string{fmt.Sprintf("%s 不能超过 %d 个字符", f.Key, f.Max)}
		}
	}
	if len(errs) > 0 {
		validationFailed(c, errs)
		return false
	}
	return true
}
