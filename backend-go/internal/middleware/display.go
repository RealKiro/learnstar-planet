package middleware

import (
	"net/http"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DisplayContextKey 用于在 gin.Context 中存取当前大屏所属班级 ID。
const DisplayContextKey = "display_class_id"

// DisplayAuth 校验班级大屏（教室端）token，并把 class_id 写入上下文。
//
// token 取值顺序与 Laravel DisplayController::validateToken 一致：
//  1. 先取 query 参数 token；
//  2. 为空时再取 Authorization: Bearer，且只接受 disp_ / class_ 前缀。
//
// 前缀与有效期校验见 models.ResolveDisplayToken；失效统一返回
// 401 {"message":"Token 无效或已过期"}。
func DisplayAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			bearer := bearerToken(c)
			if strings.HasPrefix(bearer, models.DisplayTokenPrefix) ||
				strings.HasPrefix(bearer, models.ClassTokenPrefix) {
				token = bearer
			}
		}

		classID, ok := models.ResolveDisplayToken(db, token)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token 无效或已过期"})
			return
		}

		c.Set(DisplayContextKey, classID)
		c.Next()
	}
}

// DisplayClassID 从上下文取出大屏当前班级 ID；未注入时返回 0。
func DisplayClassID(c *gin.Context) uint {
	if v, ok := c.Get(DisplayContextKey); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// bearerToken 取出 Authorization: Bearer <token> 中的 token（无则返回空串）。
func bearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
