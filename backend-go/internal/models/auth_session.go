// 认证会话相关模型：第三方绑定 + JWT 撤销名单。
//
// 移植自 Laravel：
//   - 表 third_party_bindings（迁移 2025_01_01_000001_create_all_tables.php 第 47 行 +
//     App\Models\ThirdPartyBinding）
//   - 会话撤销：Laravel 用 Sanctum 的 personal_access_tokens 删令牌实现「登出即失效」，
//     Go 端是自签 JWT（无服务端会话），故用 revoked_tokens 表按 jti 记录被撤销的令牌
//     （见 internal/middleware/auth.go 与 README「本批有意差异」）。
package models

import "time"

// ThirdPartyBinding 第三方账号绑定（一用户 × 一平台一条）。
// 字段照 Laravel 迁移：user_id / platform / platform_id / platform_union_id /
// platform_nick / platform_avatar / verified_at。
type ThirdPartyBinding struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	Platform        string     `gorm:"size:32;index;not null" json:"platform"`
	PlatformID      string     `gorm:"size:191;index" json:"platform_id"`
	PlatformUnionID string     `gorm:"size:191;index" json:"platform_union_id"`
	PlatformNick    string     `gorm:"size:100" json:"platform_nick"`
	PlatformAvatar  string     `gorm:"size:500" json:"platform_avatar"`
	VerifiedAt      *time.Time `json:"verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ThirdPartyPlatforms 全部支持的第三方平台（顺序逐字同 Laravel ThirdPartyBinding::platforms()）。
func ThirdPartyPlatforms() []string {
	return []string{"wechat", "wechat_work", "qq", "renren", "dingtalk", "feishu"}
}

// ThirdPartyPlatformLabels 平台显示名（同 ThirdPartyBinding::platformLabels()）。
func ThirdPartyPlatformLabels() map[string]string {
	return map[string]string{
		"wechat":      "微信",
		"wechat_work": "企业微信",
		"qq":          "QQ",
		"renren":      "人人通空间",
		"dingtalk":    "钉钉",
		"feishu":      "飞书",
	}
}

// ThirdPartyPlatformIcons 平台图标（同 ThirdPartyBinding::platformIcons()）。
func ThirdPartyPlatformIcons() map[string]string {
	return map[string]string{
		"wechat":      "💬",
		"wechat_work": "🏢",
		"qq":          "🐧",
		"renren":      "🌐",
		"dingtalk":    "🔷",
		"feishu":      "🪶",
	}
}

// DefaultThirdPartyPlatforms 学校未配置 enabled_third_party_platforms 时的默认三平台
// （同 Laravel AuthController::getBindings 的 ['wechat_work','wechat','qq']）。
var DefaultThirdPartyPlatforms = []string{"wechat_work", "wechat", "qq"}

// RevokedToken 被撤销的访问令牌（jti → 过期时间）。
// 命中该表且未过期的 jti 一律视为无效（中间件返回 401）。
type RevokedToken struct {
	JTI       string    `gorm:"primaryKey;size:64" json:"jti"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TempBindingContext 扫码登录的临时绑定上下文（temp_token → 上下文 JSON）。
//
// 替代 Laravel AuthService::storeTempBindingContext 的
// `Cache::put('wechat_scan_ctx:<uuid>', $ctx, now()->addMinutes(10))`：
// Go 端没有 Cache 层，改为落表（多实例共享，语义等价）。
// 读取时校验 expires_at（过期即视为不存在），绑定成功后删除记录（一次性），
// 每次写入/读取顺带惰性清理过期行（无定时任务）。
type TempBindingContext struct {
	// TempToken 一次性临时令牌（UUID v4 字符串，crypto/rand 生成，无新依赖）。
	TempToken string `gorm:"primaryKey;size:64" json:"temp_token"`
	// Context 上下文 JSON 文本（platform / platform_id / unionid / nick / avatar）。
	Context   string    `gorm:"type:text" json:"context"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WechatWorkToken 企业微信 access_token 缓存（每校一行）。
//
// 替代 Laravel WechatWorkService 的 `Cache::put("wecom_at:<schoolId>", $token, $ttl)`：
// Go 端没有 Cache 层，改为落表（多实例共享同一 token，避免各自请求 gettoken 触发限频）。
// TTL 口径与 Laravel 一致：`max(expires_in - 300, 60)` 秒。
type WechatWorkToken struct {
	SchoolID  uint      `gorm:"primaryKey" json:"school_id"`
	Token     string    `gorm:"size:512" json:"token"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
