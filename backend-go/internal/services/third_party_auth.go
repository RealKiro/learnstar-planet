// 第三方扫码登录 / 账号绑定服务。
//
// 忠实移植自 Laravel AuthService 的 loginWithWechat / loginWithWechatWork / loginWithThirdParty /
// loginWithQQ / loginWithRenren / storeTempBindingContext / getTempBindingContext /
// bindThirdParty / bindAfterScan / unbindThirdParty，以及 AuthController 的
// teacherLoginWithCredentials 复用点（bindAfterScan 的账号密码校验）。
//
// 与 Laravel 的有意差异（逐条）：
//  1. **暂存上下文落表**：Laravel `Cache::put('wechat_scan_ctx:<uuid>', $ctx, 10min)`；Go 端用表
//     `temp_binding_contexts(temp_token 主键, context JSON, expires_at)`，读取时校验未过期
//     （过期即视为不存在），绑定成功后删除（一次性），写入/读取时惰性清理过期行（无定时任务）。
//  2. **UUID 生成用 crypto/rand**（无新依赖），格式同 `Str::uuid()`（v4）。
//  3. **JWT 替代 Sanctum**：Laravel `createToken(...)->plainTextToken`；Go 端签发自签 JWT。
//     令牌形态不同（Laravel 是 `id|hash`，Go 是 JWT），但 `data.token` 的语义一致。
//  4. **无拼音库**：`uniqueNickname` 在 Laravel 里是「姓名拼音」；Go 端沿用本仓既有约定
//     （见 admin_ops_accounts.go 的说明）——昵称默认取**姓名本身**，等价于 Laravel 拼音库缺失时的
//     回退路径；`bindAfterScan` 的「本地昵称是否仍是默认值」判断因此比较 `user.Name`。
//  5. username 唯一性在 Laravel 里是「全校唯一」，Go 的 `users.username` 是全局唯一索引，
//     故去重按全局判重（与 AdminOps.CreateTeacherAccounts 一致，避免落库撞唯一索引）。
package services

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// tempBindingTTL 扫码上下文有效期（同 Laravel `now()->addMinutes(10)`）。
const tempBindingTTL = 10 * time.Minute

// ============================================================
// 结果类型（字段与键名逐字对齐 Laravel 返回数组）
// ============================================================

// AuthURLResult 扫码授权 URL 结果（GET /auth/third-party/auth-url）。
type AuthURLResult struct {
	Platform string `json:"platform"`
	AuthURL  string `json:"auth_url"`
}

// PlatformOption 登录页可选的第三方平台（GET /auth/third-party/options 的 data 元素）。
type PlatformOption struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// LoggedInResult 已绑定平台的直接登录结果（wechat / qq / renren 的 logged_in 分支，
// 同 Laravel：**不带 token**）。
type LoggedInResult struct {
	Status string       `json:"status"`
	User   *models.User `json:"user"`
}

// WechatNeedBindingResult 微信未绑定分支（`{status, temp_token, openid, unionid}`）。
type WechatNeedBindingResult struct {
	Status    string  `json:"status"`
	TempToken string  `json:"temp_token"`
	OpenID    string  `json:"openid"`
	UnionID   *string `json:"unionid"`
}

// QQNeedBindingResult QQ 未绑定分支（`{status, temp_token, openid}`）。
type QQNeedBindingResult struct {
	Status    string `json:"status"`
	TempToken string `json:"temp_token"`
	OpenID    string `json:"openid"`
}

// RenrenNeedBindingResult 人人通未绑定分支（`{status, temp_token, platform_id}`）。
type RenrenNeedBindingResult struct {
	Status     string `json:"status"`
	TempToken  string `json:"temp_token"`
	PlatformID string `json:"platform_id"`
}

// TokenLoginResult 带 token 的登录结果（企业微信 / 第三方平台）。
// 已绑定平台时 Laravel 返回 `{status, user}`（无 token），统一结构与 omitempty 表达。
type TokenLoginResult struct {
	Status  string       `json:"status"`
	Message string       `json:"message,omitempty"`
	Token   string       `json:"token,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

// BindAfterScanResult 扫码后绑定结果（POST /auth/teacher/bind-after-scan 的 data）。
type BindAfterScanResult struct {
	Status  string       `json:"status"`
	Message string       `json:"message,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

// tempBindingContext 暂存的扫码上下文（键名同 Laravel AuthService 存的数组）。
type tempBindingContext struct {
	Platform   string  `json:"platform"`
	PlatformID string  `json:"platform_id"`
	UnionID    *string `json:"unionid"`
	Nick       *string `json:"nick"`
	Avatar     *string `json:"avatar"`
}

// ============================================================
// 服务
// ============================================================

// ThirdParty 第三方登录 / 绑定服务。
type ThirdParty struct {
	db         *gorm.DB
	jwtMgr     *auth.Manager
	wechatWork *WechatWorkService
	manager    *ThirdPartyManager
	ops        *AdminOps
}

// NewThirdPartyService 创建第三方登录服务（内含企微服务与 provider 管理器）。
func NewThirdPartyService(db *gorm.DB, jwtMgr *auth.Manager) *ThirdParty {
	wechatWork := NewWechatWorkService(db)
	return &ThirdParty{
		db:         db,
		jwtMgr:     jwtMgr,
		wechatWork: wechatWork,
		manager:    NewThirdPartyManager(wechatWork),
		ops:        NewAdminOps(db),
	}
}

// Manager 返回 provider 管理器（供控制器按学校配置取 provider）。
func (s *ThirdParty) Manager() *ThirdPartyManager { return s.manager }

// WechatWork 返回企微服务（供 code 换 userid 等复用）。
func (s *ThirdParty) WechatWork() *WechatWorkService { return s.wechatWork }

// LoginURL 返回 `url('/login')` 的等价物（Laravel：APP_URL + /login）。
//
// 差异说明：Go 端没有 url() 助手，读环境变量 APP_URL（默认 http://localhost）
// —— 同 `/admin/system/status` 的 app_url 口径；未配置时与 Laravel 未设置 APP_URL 时的默认值一致。
func LoginURL() string {
	return envOrDefault("APP_URL", "http://localhost") + "/login"
}

// ============================================================
// 学校解析
// ============================================================
// SchoolForAuthURL 取第三方登录配置所属学校：schoolID > 0 时按其查找（找不到即返回 nil → 400
// 「系统尚未初始化」，同 Laravel `$schoolId > 0 ? School::find($schoolId) : School::first()`），
// 未传 school_id 时取第一所学校。
func (s *ThirdParty) SchoolForAuthURL(schoolID uint) (*models.School, error) {
	if schoolID > 0 {
		var school models.School
		err := s.db.First(&school, schoolID).Error
		if err == nil {
			return &school, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	return s.firstSchool()
}

// ResolveSchoolFromState 从 OAuth state 解析学校（provider 的 authUrl 以 schoolId 为 state）。
// 解析失败回退第一所学校（单校部署兼容，同 Laravel resolveSchoolFromState）。
func (s *ThirdParty) ResolveSchoolFromState(state string) (*models.School, error) {
	if id := phpIntCast(state); id > 0 {
		var school models.School
		err := s.db.First(&school, id).Error
		if err == nil {
			return &school, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return s.firstSchool()
}

// firstSchool 取第一所学校（`School::first()`：按主键升序取第一条）。
func (s *ThirdParty) firstSchool() (*models.School, error) {
	var school models.School
	err := s.db.Order("id ASC").First(&school).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &school, nil
}

// SchoolByID 按 ID 取学校；不存在返回 (nil, nil)（同 Laravel `$request->user()->school` 的 null 语义）。
func (s *ThirdParty) SchoolByID(schoolID uint) (*models.School, error) {
	var school models.School
	err := s.db.First(&school, schoolID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &school, nil
}

// phpIntCast 等价 PHP `(int) $value`：取前导可选符号 + 数字，无数字则为 0（不报错）。
func phpIntCast(value string) int {
	trimmed := strings.TrimSpace(value)
	sign := 1
	idx := 0
	if strings.HasPrefix(trimmed, "-") {
		sign = -1
		idx = 1
	} else if strings.HasPrefix(trimmed, "+") {
		idx = 1
	}
	digits := strings.Builder{}
	for i := idx; i < len(trimmed); i++ {
		ch := trimmed[i]
		if ch < '0' || ch > '9' {
			break
		}
		digits.WriteByte(ch)
	}
	if digits.Len() == 0 {
		return 0
	}
	n := 0
	for _, ch := range digits.String() {
		n = n*10 + int(ch-'0')
	}
	return sign * n
}

// PhpIntCast 供 handlers 层复用（企微回调的 school_id / SpStatus 都要按 PHP `(int)` 语义取值）。
func PhpIntCast(value string) int { return phpIntCast(value) }

// ============================================================
// 登录页平台选项 / 授权 URL
// ============================================================

// AuthURLOfSchool 组装授权 URL（平台未配置 / 未知 → 控制器返回 400 固定文案）。
func (s *ThirdParty) AuthURLOfSchool(school *models.School, redirectURI string) (*AuthURLResult, error) {
	provider, err := s.manager.ProviderFor(school.SettingString("third_party_platform", ""))
	if err != nil {
		return nil, err
	}
	return &AuthURLResult{Platform: provider.Key(), AuthURL: provider.AuthURL(school.ID, redirectURI)}, nil
}

// Options 登录页可选平台列表：`schools.settings.enabled_third_party_platforms` 与 platforms() 取交集，
// 未配置或交集为空时回退默认三平台 [企业微信, 微信, QQ]（同 Laravel thirdPartyOptions）。
func (s *ThirdParty) Options() ([]PlatformOption, error) {
	platforms := models.DefaultThirdPartyPlatforms

	school, err := s.firstSchool()
	if err != nil {
		return nil, err
	}
	if school != nil {
		if enabled, ok := school.SettingsMap()["enabled_third_party_platforms"].([]any); ok && len(enabled) > 0 {
			allowed := map[string]bool{}
			for _, raw := range enabled {
				if name, isStr := raw.(string); isStr {
					allowed[name] = true
				}
			}
			intersection := make([]string, 0, len(enabled))
			for _, p := range models.ThirdPartyPlatforms() {
				if allowed[p] {
					intersection = append(intersection, p)
				}
			}
			if len(intersection) > 0 {
				platforms = intersection
			}
		}
	}

	labels := models.ThirdPartyPlatformLabels()
	icons := models.ThirdPartyPlatformIcons()
	options := make([]PlatformOption, 0, len(platforms))
	for _, p := range platforms {
		label, ok := labels[p]
		if !ok {
			label = p
		}
		icon, ok := icons[p]
		if !ok {
			icon = "🔗"
		}
		options = append(options, PlatformOption{Key: p, Label: label, Icon: icon, Color: platformColor(p)})
	}
	return options, nil
}

// platformColor 平台品牌色（match 表逐字同 Laravel AuthController::thirdPartyOptions）。
func platformColor(platform string) string {
	switch platform {
	case "wechat_work":
		return "#2B7CE9"
	case "dingtalk":
		return "#0089FF"
	case "feishu":
		return "#3370FF"
	case "wechat":
		return "#07C160"
	case "qq":
		return "#12B7F5"
	case "renren":
		return "#FF6A00"
	default:
		return "#7c3aed"
	}
}

// ============================================================
// 第三方平台扫码登录（企微 / 钉钉 / 飞书）
// ============================================================

// LoginWithThirdParty 第三方平台登录：已绑定直接登录；未绑定先匹配本地已有账号（手机号 → 实名用户名），
// 仍无匹配则自动创建教师账号（姓名实名 + 默认密码 ls123456）并绑定。
//
// 返回值 status 为 logged_in 或 error（error 时带 message，同 Laravel 的 `status=error` 分支）。
func (s *ThirdParty) LoginWithThirdParty(platform string, info ThirdPartyUserInfo, school *models.School) (*TokenLoginResult, error) {
	if user, err := s.findUserByPlatform(platform, info.PlatformID); err != nil {
		return nil, err
	} else if user != nil {
		if err := s.touchLastLogin(user.ID); err != nil {
			return nil, err
		}
		return &TokenLoginResult{Status: "logged_in", User: user}, nil
	}

	if school == nil {
		var err error
		school, err = s.firstSchool()
		if err != nil {
			return nil, err
		}
	}
	if school == nil {
		return &TokenLoginResult{Status: "error", Message: "系统尚未初始化，请先联系管理员"}, nil
	}

	// 无绑定时先匹配本地已有账号（通讯录导入等场景），避免同一教师被重复建号。
	if existing, err := s.findLocalTeacherForThirdParty(school.ID, info.Name, info.Mobile); err != nil {
		return nil, err
	} else if existing != nil {
		sync := map[string]any{}
		if existing.Phone == "" && info.Mobile != "" {
			sync["phone"] = info.Mobile
		}
		if existing.AvatarPath == "" && info.Avatar != "" {
			sync["avatar_path"] = info.Avatar
		}
		if len(sync) > 0 {
			if err := s.db.Model(&models.User{}).Where("id = ?", existing.ID).Updates(sync).Error; err != nil {
				return nil, err
			}
		}
		if err := s.createBindingRecord(existing.ID, platform, info.PlatformID, nil, &info.Name, &info.Avatar); err != nil {
			return nil, err
		}
		token, err := s.jwtMgr.Generate(existing.ID, existing.Role, existing.SchoolID)
		if err != nil {
			return nil, err
		}
		if err := s.touchLastLogin(existing.ID); err != nil {
			return nil, err
		}
		fresh, err := s.userByID(existing.ID)
		if err != nil {
			return nil, err
		}
		return &TokenLoginResult{Status: "logged_in", Token: token, User: fresh}, nil
	}

	displayName := info.Name
	if displayName == "" {
		displayName = info.PlatformID
	}
	user, err := s.createThirdPartyTeacher(school.ID, displayName, info.Avatar, info.Mobile, info.Email)
	if err != nil {
		return nil, err
	}
	if err := s.createBindingRecord(user.ID, platform, info.PlatformID, nil, &info.Name, &info.Avatar); err != nil {
		return nil, err
	}
	token, err := s.jwtMgr.Generate(user.ID, user.Role, user.SchoolID)
	if err != nil {
		return nil, err
	}
	if err := s.touchLastLogin(user.ID); err != nil {
		return nil, err
	}
	fresh, err := s.userByID(user.ID)
	if err != nil {
		return nil, err
	}
	return &TokenLoginResult{Status: "logged_in", Token: token, User: fresh}, nil
}

// ============================================================
// 微信 / QQ / 人人通扫码登录
// ============================================================

// LoginWithWechat 微信扫码登录：优先 unionid（可跨小程序与开放平台）→ openid；
// 已绑定直接登录（**无 token**，Laravel 原样行为），未绑定返回 need_binding + temp_token。
func (s *ThirdParty) LoginWithWechat(openid string, unionID, nick, avatar *string) (any, error) {
	if unionID != nil && *unionID != "" {
		user, err := s.findUserByUnionID(*unionID)
		if err != nil {
			return nil, err
		}
		if user != nil {
			if err := s.touchLastLogin(user.ID); err != nil {
				return nil, err
			}
			fresh, err := s.userByID(user.ID)
			if err != nil {
				return nil, err
			}
			return LoggedInResult{Status: "logged_in", User: fresh}, nil
		}
	}

	user, err := s.findUserByPlatform("wechat", openid)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err := s.touchLastLogin(user.ID); err != nil {
			return nil, err
		}
		fresh, err := s.userByID(user.ID)
		if err != nil {
			return nil, err
		}
		return LoggedInResult{Status: "logged_in", User: fresh}, nil
	}

	tempToken, err := s.storeTempBindingContext(tempBindingContext{
		Platform:   "wechat",
		PlatformID: openid,
		UnionID:    unionID,
		Nick:       nick,
		Avatar:     avatar,
	})
	if err != nil {
		return nil, err
	}
	return WechatNeedBindingResult{Status: "need_binding", TempToken: tempToken, OpenID: openid, UnionID: unionID}, nil
}

// LoginWithWechatWork 企业微信扫码登录：已绑定直接登录（**无 token**）；
// 未绑定自动创建教师账号（免注册，`nick ?? userid` 为姓名、默认密码 ls123456）并返回 token。
func (s *ThirdParty) LoginWithWechatWork(userid string, nick, avatar *string, school *models.School) (*TokenLoginResult, error) {
	user, err := s.findUserByPlatform("wechat_work", userid)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err := s.touchLastLogin(user.ID); err != nil {
			return nil, err
		}
		fresh, err := s.userByID(user.ID)
		if err != nil {
			return nil, err
		}
		// Laravel 原样行为：该分支只返回 {status, user}，**不带 token**。
		return &TokenLoginResult{Status: "logged_in", User: fresh}, nil
	}

	if school == nil {
		school, err = s.firstSchool()
		if err != nil {
			return nil, err
		}
	}
	if school == nil {
		return &TokenLoginResult{Status: "error", Message: "系统尚未初始化，请先联系管理员"}, nil
	}

	name := userid
	if nick != nil {
		name = *nick
	}
	avatarPath := ""
	if avatar != nil {
		avatarPath = *avatar
	}

	created, err := s.createThirdPartyTeacher(school.ID, name, avatarPath, "", "")
	if err != nil {
		return nil, err
	}
	if err := s.createBindingRecord(created.ID, "wechat_work", userid, nil, nick, avatar); err != nil {
		return nil, err
	}
	token, err := s.jwtMgr.Generate(created.ID, created.Role, created.SchoolID)
	if err != nil {
		return nil, err
	}
	if err := s.touchLastLogin(created.ID); err != nil {
		return nil, err
	}
	fresh, err := s.userByID(created.ID)
	if err != nil {
		return nil, err
	}
	return &TokenLoginResult{Status: "logged_in", Token: token, User: fresh}, nil
}

// LoginWithQQ QQ 扫码登录（已绑定无 token；未绑定 need_binding）。
func (s *ThirdParty) LoginWithQQ(openid string, nick, avatar *string) (any, error) {
	user, err := s.findUserByPlatform("qq", openid)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err := s.touchLastLogin(user.ID); err != nil {
			return nil, err
		}
		fresh, err := s.userByID(user.ID)
		if err != nil {
			return nil, err
		}
		return LoggedInResult{Status: "logged_in", User: fresh}, nil
	}

	tempToken, err := s.storeTempBindingContext(tempBindingContext{
		Platform:   "qq",
		PlatformID: openid,
		Nick:       nick,
		Avatar:     avatar,
	})
	if err != nil {
		return nil, err
	}
	return QQNeedBindingResult{Status: "need_binding", TempToken: tempToken, OpenID: openid}, nil
}

// LoginWithRenren 人人通空间登录（已绑定无 token；未绑定 need_binding + platform_id）。
func (s *ThirdParty) LoginWithRenren(userID string, nick, avatar *string) (any, error) {
	user, err := s.findUserByPlatform("renren", userID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err := s.touchLastLogin(user.ID); err != nil {
			return nil, err
		}
		fresh, err := s.userByID(user.ID)
		if err != nil {
			return nil, err
		}
		return LoggedInResult{Status: "logged_in", User: fresh}, nil
	}

	tempToken, err := s.storeTempBindingContext(tempBindingContext{
		Platform:   "renren",
		PlatformID: userID,
		Nick:       nick,
		Avatar:     avatar,
	})
	if err != nil {
		return nil, err
	}
	return RenrenNeedBindingResult{Status: "need_binding", TempToken: tempToken, PlatformID: userID}, nil
}

// ============================================================
// 扫码后绑定已有账号 / 主动绑定 / 解绑
// ============================================================

// BindAfterScan 第三方首次扫码后绑定已有教师账号（仅 role=teacher）。
//
// 未命中账号密码时返回 status=error（HTTP 仍为 200，同 Laravel）；成功返回 status=bound。
// 语义细节照抄 Laravel：优先用 temp_token 对应的暂存上下文覆盖 platform/platform_id/unionid/nick/avatar；
// nickname 仅在「本地为空 或 等于默认昵称（本仓约定 = 姓名）」时覆盖；avatar_path 仅在本地为空时覆盖；
// **用完删除暂存上下文**（一次性）。
func (s *ThirdParty) BindAfterScan(tempToken, username, password, platform, platformID string, unionID, nick, avatar *string) (*BindAfterScanResult, error) {
	user, err := s.teacherLoginForBinding(username, password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &BindAfterScanResult{Status: "error", Message: "账号或密码错误，请核对后重试"}, nil
	}

	ctx, err := s.getTempBindingContext(tempToken)
	if err != nil {
		return nil, err
	}
	ctxPlatform, ctxPlatformID := platform, platformID
	ctxUnionID, ctxNick, ctxAvatar := unionID, nick, avatar
	if ctx != nil {
		// PHP `?? ` 只跳过 null；上下文里的 platform / platform_id 恒非空，等价于「非空即覆盖」。
		if ctx.Platform != "" {
			ctxPlatform = ctx.Platform
		}
		if ctx.PlatformID != "" {
			ctxPlatformID = ctx.PlatformID
		}
		if ctx.UnionID != nil {
			ctxUnionID = ctx.UnionID
		}
		if ctx.Nick != nil {
			ctxNick = ctx.Nick
		}
		if ctx.Avatar != nil {
			ctxAvatar = ctx.Avatar
		}
	}

	if err := s.bindThirdPartyForUser(user.ID, ctxPlatform, ctxPlatformID, ctxUnionID, ctxNick, ctxAvatar); err != nil {
		return nil, err
	}

	updates := map[string]any{}
	// Laravel：`$ctxNick && (empty($user->nickname) || $user->nickname === PinyinService::toPinyin($user->name))`。
	// Go 无拼音库，昵称约定为「姓名本身」（见 admin_ops_accounts.go），故比较对象为 Name。
	if ctxNick != nil && *ctxNick != "" && (user.Nickname == "" || user.Nickname == user.Name) {
		updates["nickname"] = *ctxNick
	}
	if ctxAvatar != nil && *ctxAvatar != "" && user.AvatarPath == "" {
		updates["avatar_path"] = *ctxAvatar
	}
	if len(updates) > 0 {
		if err := s.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := s.deleteTempBindingContext(tempToken); err != nil {
		return nil, err
	}

	fresh, err := s.userByID(user.ID)
	if err != nil {
		return nil, err
	}
	return &BindAfterScanResult{Status: "bound", User: fresh}, nil
}

// BindForUser 登录后主动绑定第三方账号：该 platform + platform_id 已被**任意用户**绑定时 422。
func (s *ThirdParty) BindForUser(user *models.User, platform, platformID string, nick, avatar *string) error {
	var existing models.ThirdPartyBinding
	err := s.db.Where("platform = ? AND platform_id = ?", platform, platformID).First(&existing).Error
	if err == nil {
		return ErrUnprocessable("该第三方账号已被其他用户绑定")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return s.createBindingRecord(user.ID, platform, platformID, nil, nick, avatar)
}

// UnbindForUser 解绑当前用户该平台的绑定（无绑定也返回成功，同 Laravel）。
func (s *ThirdParty) UnbindForUser(user *models.User, platform string) error {
	return s.db.Where("user_id = ? AND platform = ?", user.ID, platform).
		Delete(&models.ThirdPartyBinding{}).Error
}

// ============================================================
// 内部辅助
// ============================================================

// findUserByPlatform 按平台 + platform_id 反查用户（绑定不存在或对应用户已删除时返回 nil）。
func (s *ThirdParty) findUserByPlatform(platform, platformID string) (*models.User, error) {
	var binding models.ThirdPartyBinding
	err := s.db.Where("platform = ? AND platform_id = ?", platform, platformID).First(&binding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.userByID(binding.UserID)
}

// findUserByUnionID 按 unionid 反查用户（同 ThirdPartyBinding::findUserByUnionId，取第一条绑定）。
func (s *ThirdParty) findUserByUnionID(unionID string) (*models.User, error) {
	var binding models.ThirdPartyBinding
	err := s.db.Where("platform_union_id = ?", unionID).First(&binding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.userByID(binding.UserID)
}

// userByID 取用户；不存在返回 (nil, nil)。
func (s *ThirdParty) userByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// findLocalTeacherForThirdParty 无绑定时匹配本地已有教师账号：
// 手机号精确匹配优先，其次实名用户名（通讯录导入默认 username = 姓名）。
func (s *ThirdParty) findLocalTeacherForThirdParty(schoolID uint, name, mobile string) (*models.User, error) {
	if mobile != "" {
		var user models.User
		err := s.db.Where("school_id = ? AND role = ? AND phone = ?", schoolID, "teacher", mobile).
			First(&user).Error
		if err == nil {
			return &user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if name != "" {
		var user models.User
		err := s.db.Where("school_id = ? AND role = ? AND username = ?", schoolID, "teacher", name).
			First(&user).Error
		if err == nil {
			return &user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return nil, nil
}

// createThirdPartyTeacher 第三方免注册建号：username 默认 = 姓名（全局去重 _2/_3…）、
// nickname 默认 = 姓名（校内去重）、默认密码 ls123456（同时写明文，同 Laravel）。
func (s *ThirdParty) createThirdPartyTeacher(schoolID uint, name, avatar, phone, email string) (*models.User, error) {
	username, err := s.ops.uniqueValue(name, func(candidate string) (bool, error) {
		var count int64
		if err := s.db.Model(&models.User{}).Where("username = ?", candidate).Count(&count).Error; err != nil {
			return false, err
		}
		return count > 0, nil
	})
	if err != nil {
		return nil, err
	}
	nickname, err := s.ops.uniqueValue(name, func(candidate string) (bool, error) {
		var count int64
		if err := s.db.Model(&models.User{}).
			Where("school_id = ? AND nickname = ?", schoolID, candidate).
			Count(&count).Error; err != nil {
			return false, err
		}
		return count > 0, nil
	})
	if err != nil {
		return nil, err
	}

	password := DefaultTeacherPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		SchoolID:        schoolID,
		Role:            "teacher",
		Username:        username,
		PasswordHash:    string(hash),
		PlainPassword:   password,
		Name:            name,
		Nickname:        nickname,
		AvatarPath:      avatar,
		Phone:           phone,
		Email:           email,
		Status:          "active",
		PasswordChanged: false,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// teacherLoginForBinding 账号密码校验（仅 role=teacher 且 status=active）；
// 失败返回 (nil, nil)（同 Laravel teacherLoginWithCredentials 的 null 返回），成功时刷新 last_login_at。
func (s *ThirdParty) teacherLoginForBinding(username, password string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ? AND role = ? AND status = ?", username, "teacher", "active").
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, nil
	}
	if err := s.touchLastLogin(user.ID); err != nil {
		return nil, err
	}
	now := util.Now()
	user.LastLoginAt = &now
	return &user, nil
}

// touchLastLogin 刷新最后登录时间（同 Laravel `$user->update(['last_login_at' => now()])`）。
func (s *ThirdParty) touchLastLogin(userID uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login_at", util.Now()).Error
}

// createBindingRecord 建绑定记录（同 ThirdPartyBinding::create；空串昵称/头像落空串，
// 对外 JSON 由既有 Bindings() 统一转 null/字符串）。
func (s *ThirdParty) createBindingRecord(userID uint, platform, platformID string, unionID, nick, avatar *string) error {
	now := util.Now()
	binding := models.ThirdPartyBinding{
		UserID:     userID,
		Platform:   platform,
		PlatformID: platformID,
		VerifiedAt: &now,
	}
	if unionID != nil {
		binding.PlatformUnionID = *unionID
	}
	if nick != nil {
		binding.PlatformNick = *nick
	}
	if avatar != nil {
		binding.PlatformAvatar = *avatar
	}
	return s.db.Create(&binding).Error
}

// bindThirdPartyForUser 等价 Laravel AuthService::bindThirdParty：该用户已有同平台绑定则更新，否则新建。
func (s *ThirdParty) bindThirdPartyForUser(userID uint, platform, platformID string, unionID, nick, avatar *string) error {
	var existing models.ThirdPartyBinding
	err := s.db.Where("user_id = ? AND platform = ?", userID, platform).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.createBindingRecord(userID, platform, platformID, unionID, nick, avatar)
	}
	if err != nil {
		return err
	}
	return s.db.Model(&models.ThirdPartyBinding{}).Where("id = ?", existing.ID).Updates(map[string]any{
		"platform_id":       platformID,
		"platform_union_id": unionID,
		"platform_nick":     nick,
		"platform_avatar":   avatar,
		"verified_at":       util.Now(),
	}).Error
}

// ============================================================
// 暂存上下文（表替代 Cache）
// ============================================================

// newTempToken 生成 UUID v4 字符串（同 Laravel `Str::uuid()->toString()`，crypto/rand，无新依赖）。
func newTempToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	buf[6] = (buf[6] & 0x0f) | 0x40 // version 4
	buf[8] = (buf[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16]), nil
}

// storeTempBindingContext 存上下文（10 分钟有效），返回 temp_token；
// 顺带惰性清理过期行（无定时任务）。
func (s *ThirdParty) storeTempBindingContext(ctx tempBindingContext) (string, error) {
	tempToken, err := newTempToken()
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}
	if err := s.pruneTempContexts(); err != nil {
		return "", err
	}
	if err := s.db.Create(&models.TempBindingContext{
		TempToken: tempToken,
		Context:   string(payload),
		ExpiresAt: util.Now().Add(tempBindingTTL),
	}).Error; err != nil {
		return "", err
	}
	return tempToken, nil
}

// getTempBindingContext 读上下文；不存在或已过期返回 (nil, nil)（过期行顺带删除）。
func (s *ThirdParty) getTempBindingContext(tempToken string) (*tempBindingContext, error) {
	var row models.TempBindingContext
	err := s.db.First(&row, "temp_token = ?", tempToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !row.ExpiresAt.After(util.Now()) {
		if err := s.deleteTempBindingContext(tempToken); err != nil {
			return nil, err
		}
		return nil, nil
	}
	ctx := tempBindingContext{}
	if err := json.Unmarshal([]byte(row.Context), &ctx); err != nil {
		return nil, nil
	}
	return &ctx, nil
}

// deleteTempBindingContext 删除上下文（一次性消费）。
func (s *ThirdParty) deleteTempBindingContext(tempToken string) error {
	return s.db.Where("temp_token = ?", tempToken).Delete(&models.TempBindingContext{}).Error
}

// pruneTempContexts 惰性清理过期上下文。
func (s *ThirdParty) pruneTempContexts() error {
	return s.db.Where("expires_at <= ?", util.Now()).Delete(&models.TempBindingContext{}).Error
}
