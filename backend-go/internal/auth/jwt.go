// Package auth 提供 JWT 签发与校验。
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 是访问令牌携带的自定义声明。jti 走标准 RegisteredClaims.ID（序列化为 "jti"）。
type Claims struct {
	UserID   uint   `json:"uid"`
	Role     string `json:"role"`
	SchoolID uint   `json:"sid"`
	jwt.RegisteredClaims
}

// JTI 返回令牌的 jti；旧令牌（签发时无 jti）返回空串。
func (c *Claims) JTI() string {
	if c == nil {
		return ""
	}
	return c.ID
}

// Expiry 返回令牌过期时间；缺失返回 nil。
func (c *Claims) Expiry() *time.Time {
	if c == nil || c.ExpiresAt == nil {
		return nil
	}
	t := c.ExpiresAt.Time
	return &t
}

// Manager 负责签发与解析 JWT。
type Manager struct {
	secret  []byte
	expTime time.Duration
}

// New 创建 JWT 管理器。
func New(secret string, expHours int) *Manager {
	return &Manager{
		secret:  []byte(secret),
		expTime: time.Duration(expHours) * time.Hour,
	}
}

// Generate 为指定用户签发令牌（每次签发带唯一 jti，供登出/刷新撤销）。
func (m *Manager) Generate(userID uint, role string, schoolID uint) (string, error) {
	now := time.Now()
	jti, err := newJTI()
	if err != nil {
		return "", err
	}
	claims := Claims{
		UserID:   userID,
		Role:     role,
		SchoolID: schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    "learnstar-planet",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Exp 返回令牌有效期时长。
func (m *Manager) Exp() time.Duration { return m.expTime }

// newJTI 生成 16 字节随机 jti（hex 32 字符）。
func newJTI() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// Parse 校验并解析令牌，返回其中的声明。
func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("非预期的签名算法")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效令牌")
	}
	return claims, nil
}
