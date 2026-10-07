// 认证服务：登录、当前用户、修改密码。
package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务。
type AuthService struct {
	db     *gorm.DB
	jwtMgr *auth.Manager
}

// NewAuthService 创建认证服务。
func NewAuthService(db *gorm.DB, jwtMgr *auth.Manager) *AuthService {
	return &AuthService{db: db, jwtMgr: jwtMgr}
}

// Login 校验账号密码并签发 JWT。
func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, ErrUnprocessable("账号或密码错误")
	}
	if err != nil {
		return "", nil, err
	}
	if user.Status != "active" {
		return "", nil, ErrUnprocessable("账号已停用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrUnprocessable("账号或密码错误")
	}

	now := util.Now()
	if err := s.db.Model(&user).Update("last_login_at", now).Error; err != nil {
		return "", nil, err
	}
	user.LastLoginAt = &now

	token, err := s.jwtMgr.Generate(user.ID, user.Role, user.SchoolID)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

// credentialLoginErrorMessage 账号密码登录失败的统一文案（Laravel teacherLoginWithCredentials /
// adminLoginWithCredentials 对「账号不存在 / 角色不符 / 账号停用 / 密码错误」一律返回该 401 文案）。
const credentialLoginErrorMessage = "账号或密码错误，请核对后重试"

// TeacherLogin 教师账号密码登录：仅 role=teacher 且 status=active 的账号可通过。
// 移植自 AuthService::teacherLoginWithCredentials（账号不存在或角色不符同样返回 401 统一文案）。
func (s *AuthService) TeacherLogin(username, password string) (string, *models.User, error) {
	return s.credentialLogin(username, password, "teacher")
}

// AdminLogin 管理员账号密码登录：仅 role=school_admin 且 status=active 的账号可通过。
// 移植自 AuthService::adminLoginWithCredentials（管理员不支持第三方扫码，仅账号密码）。
func (s *AuthService) AdminLogin(username, password string) (string, *models.User, error) {
	return s.credentialLogin(username, password, "school_admin")
}

// credentialLogin 按角色限定账号密码登录并签发 JWT（同 Laravel 两条 xxxLoginWithCredentials）。
func (s *AuthService) credentialLogin(username, password, role string) (string, *models.User, error) {
	var user models.User
	err := s.db.Where("username = ? AND role = ? AND status = ?", username, role, "active").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, NewAppError(http.StatusUnauthorized, credentialLoginErrorMessage)
	}
	if err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, NewAppError(http.StatusUnauthorized, credentialLoginErrorMessage)
	}

	now := util.Now()
	if err := s.db.Model(&user).Update("last_login_at", now).Error; err != nil {
		return "", nil, err
	}
	user.LastLoginAt = &now

	token, err := s.jwtMgr.Generate(user.ID, user.Role, user.SchoolID)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

// Me 返回当前用户。
func (s *AuthService) Me(userID uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("账号不存在")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ChangePassword 校验原密码后更新为新密码。
func (s *AuthService) ChangePassword(u *models.User, oldPassword, newPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrUnprocessable("原密码错误")
	}
	if len(newPassword) < 6 {
		return ErrUnprocessable("新密码至少 6 位")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&models.User{}).Where("id = ?", u.ID).Updates(map[string]any{
		"password_hash":    string(hash),
		"plain_password":   newPassword,
		"password_changed": true,
	}).Error
}

// ============================================================
// 认证会话：登出 / 刷新 / 第三方绑定列表
// 移植自 Laravel AuthController::logout / refreshToken / getBindings。
// ============================================================

// BindingView 单个第三方平台的绑定状态（字段与顺序逐字同 Laravel getBindings 的 data 元素）。
type BindingView struct {
	Platform string  `json:"platform"`
	Label    string  `json:"label"`
	Icon     string  `json:"icon"`
	Bound    bool    `json:"bound"`
	Nick     *string `json:"nick"`
}

// RevokeToken 撤销一个 jti（登出 / 刷新旧令牌）。
// Laravel 是删除 Sanctum 令牌行（删了即失效）；Go 端把 jti 写入 revoked_tokens，
// 中间件命中且未过期即 401。空 jti（历史令牌，签发时无 jti）不做处理，保持向后兼容。
//
// ⚠️ 时区纪律：SQLite 对 time 列做**文本比较**，同一列的写入与比较必须同一钟面域。
// JWT 的过期时间是 Local 域（随部署环境变化），而 util.Now() 恒为 Asia/Shanghai——
// 两者混比在 TZ≠UTC 的环境会错序（未来 72h 才过期的行被惰性清理当成已过期删除，
// 登出/刷新随即失效，CI 上曾因此挂测试）。故本表全链路统一 UTC 域：
// 写入 .UTC()、清理与中间件比较一律 time.Now().UTC()。
func (s *AuthService) RevokeToken(jti string, expiresAt *time.Time) error {
	if jti == "" {
		return nil
	}
	exp := time.Now().UTC().Add(s.jwtMgr.Exp())
	if expiresAt != nil {
		exp = expiresAt.UTC()
	}
	if err := s.db.Where("jti = ?", jti).Delete(&models.RevokedToken{}).Error; err != nil {
		return err
	}
	if err := s.db.Create(&models.RevokedToken{JTI: jti, ExpiresAt: exp}).Error; err != nil {
		return err
	}
	// 顺带清理已过期记录（无定时任务，惰性清理）。
	return s.db.Where("expires_at <= ?", time.Now().UTC()).Delete(&models.RevokedToken{}).Error
}

// Logout 撤销当前访问令牌（Laravel `$user->currentAccessToken()->delete()`）。
func (s *AuthService) Logout(jti string, expiresAt *time.Time) error {
	return s.RevokeToken(jti, expiresAt)
}

// Refresh 撤销旧令牌并签发新令牌（Laravel refreshToken：删旧发新）。
func (s *AuthService) Refresh(u *models.User, jti string, expiresAt *time.Time) (string, error) {
	if err := s.RevokeToken(jti, expiresAt); err != nil {
		return "", err
	}
	return s.jwtMgr.Generate(u.ID, u.Role, u.SchoolID)
}

// Bindings 当前用户的第三方平台绑定列表。
//
// 口径逐条照抄 Laravel getBindings：仅展示学校后台启用的平台
// （`schools.settings.enabled_third_party_platforms` 与 platforms() 取交集、顺序按 platforms()）；
// 未配置或为空数组时回退默认三平台 ['wechat_work','wechat','qq']。
func (s *AuthService) Bindings(u *models.User) ([]BindingView, error) {
	var bindings []models.ThirdPartyBinding
	if err := s.db.Where("user_id = ?", u.ID).Order("id ASC").Find(&bindings).Error; err != nil {
		return nil, err
	}

	platforms := models.DefaultThirdPartyPlatforms
	var school models.School
	if err := s.db.First(&school, u.SchoolID).Error; err == nil {
		if enabled, ok := school.SettingsMap()["enabled_third_party_platforms"].([]any); ok && len(enabled) > 0 {
			allowed := map[string]bool{}
			for _, raw := range enabled {
				if name, isStr := raw.(string); isStr {
					allowed[name] = true
				}
			}
			platforms = []string{}
			for _, p := range models.ThirdPartyPlatforms() {
				if allowed[p] {
					platforms = append(platforms, p)
				}
			}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	labels := models.ThirdPartyPlatformLabels()
	icons := models.ThirdPartyPlatformIcons()
	out := make([]BindingView, 0, len(platforms))
	for _, p := range platforms {
		label, ok := labels[p]
		if !ok {
			label = p
		}
		icon, ok := icons[p]
		if !ok {
			icon = "🔗"
		}
		view := BindingView{Platform: p, Label: label, Icon: icon}
		for i := range bindings {
			if bindings[i].Platform != p {
				continue
			}
			view.Bound = true
			if bindings[i].PlatformNick != "" {
				nick := bindings[i].PlatformNick
				view.Nick = &nick
			}
			break
		}
		out = append(out, view)
	}
	return out, nil
}
