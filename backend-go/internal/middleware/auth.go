// Package middleware 提供 HTTP 中间件。
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContextKey 用于在 gin.Context 中存取当前登录用户。
const ContextKey = "auth_user"

// ClaimsKey 用于在 gin.Context 中存取当前访问令牌的声明（供 logout/refresh 取 jti）。
const ClaimsKey = "auth_claims"

// Auth 校验 Bearer Token，并把当前用户写入上下文。
func Auth(jwtMgr *auth.Manager, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "认证格式错误"})
			return
		}

		claims, err := jwtMgr.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录已过期或无效"})
			return
		}

		// 令牌撤销检查：命中 revoked_tokens（未过期）→ 401；无 jti 的历史令牌跳过（向后兼容）。
		if jti := claims.JTI(); jti != "" {
			var revoked int64
			lookupErr := db.Model(&models.RevokedToken{}).
				Where("jti = ? AND expires_at > ?", jti, time.Now().UTC()).
				Count(&revoked).Error
			if lookupErr != nil || revoked > 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录已过期或无效"})
				return
			}
		}

		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "账号不存在"})
			return
		}
		if user.Status != "active" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "账号已停用"})
			return
		}

		c.Set(ClaimsKey, claims)

		c.Set(ContextKey, &user)
		c.Next()
	}
}

// RequireRole 限制访问角色（school_admin / teacher）。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		user, ok := c.Get(ContextKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
			return
		}
		u := user.(*models.User)
		if !allowed[u.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无权访问"})
			return
		}
		c.Next()
	}
}

// CurrentUser 从上下文取出当前登录用户；未登录返回 nil。
func CurrentUser(c *gin.Context) *models.User {
	if v, ok := c.Get(ContextKey); ok {
		if u, ok := v.(*models.User); ok {
			return u
		}
	}
	return nil
}

// CurrentClaims 从上下文取出当前访问令牌的声明；不存在返回 nil。
func CurrentClaims(c *gin.Context) *auth.Claims {
	if v, ok := c.Get(ClaimsKey); ok {
		if claims, ok := v.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}
