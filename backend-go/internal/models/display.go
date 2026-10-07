// 班级大屏（教室端）鉴权持久化模型：显示端 Token 与班级码登录日志。
//
// 有意差异：Laravel 端把 disp_ token 与 class_token: 映射放在 Cache（Redis/文件）里，
// Go 端没有 Redis / Cache，改为数据库表持久化——语义等价（过期即视为无效），
// 且 disp_ 与 class_ 两种前缀共用同一张表（class_ token 解析出来同样是 class_id）。
package models

import (
	"crypto/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 显示端 token 前缀与有效期（同 Laravel DisplayController::TOKEN_PREFIX / TOKEN_TTL）。
const (
	// DisplayTokenPrefix disp_ 前缀（本控制器签发）。
	DisplayTokenPrefix = "disp_"
	// ClassTokenPrefix class_ 前缀（Laravel AuthController::classLogin 签发，同一张表兼容）。
	ClassTokenPrefix = "class_"
	// DisplayTokenTTL token 有效期（秒）。
	DisplayTokenTTL = 86400
)

// DisplayToken 显示端登录 Token。
type DisplayToken struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Token   string `gorm:"size:64;uniqueIndex;not null" json:"token"`
	ClassID uint   `gorm:"index;not null" json:"class_id"`
	// 以下字段等价 Laravel 缓存中的 {class_id, class_name, grade, ip, created_at} 载荷。
	ClassName string    `gorm:"size:100" json:"class_name"`
	Grade     string    `gorm:"size:50" json:"grade"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DisplayLoginLog 班级码登录日志（含 IP 与 UA）。
type DisplayLoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClassID   uint      `gorm:"index;not null" json:"class_id"`
	ClassCode string    `gorm:"size:20" json:"class_code"`
	IPAddress string    `gorm:"size:45;index" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// randomAlphaNumeric 生成 n 位随机字母数字串（字符集同 Laravel Str::random）。
func randomAlphaNumeric(n int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i, b := range buf {
		buf[i] = charset[int(b)%len(charset)]
	}
	return string(buf), nil
}

// NewDisplayTokenValue 生成 disp_ + 32 位随机字母数字串（字符集同 Laravel Str::random(32)）。
func NewDisplayTokenValue() (string, error) {
	suffix, err := randomAlphaNumeric(32)
	if err != nil {
		return "", err
	}
	return DisplayTokenPrefix + suffix, nil
}

// NewClassTokenValue 生成 class_<班级码>_<32 位随机串>。
// 同 Laravel AuthController::classLogin 的 'class_' . $classCode . '_' . Str::random(32)；
// 班级码上限 10 位，故 token 最长 49 字符（token 列 64 位足够）。
func NewClassTokenValue(classCode string) (string, error) {
	suffix, err := randomAlphaNumeric(32)
	if err != nil {
		return "", err
	}
	return ClassTokenPrefix + classCode + "_" + suffix, nil
}

// ResolveDisplayToken 校验显示端 token 并返回所属班级 ID。
// 仅接受 disp_ / class_ 前缀；token 不存在或已过期一律返回 false。
// 过期判定在 Go 侧完成（不依赖各数据库的时间比较语义）。
func ResolveDisplayToken(db *gorm.DB, token string) (uint, bool) {
	if token == "" {
		return 0, false
	}
	if !strings.HasPrefix(token, DisplayTokenPrefix) && !strings.HasPrefix(token, ClassTokenPrefix) {
		return 0, false
	}

	var row DisplayToken
	if err := db.Where("token = ?", token).First(&row).Error; err != nil {
		return 0, false
	}
	if !row.ExpiresAt.After(time.Now()) {
		return 0, false
	}
	return row.ClassID, true
}

// DisplayEvent 班级大屏事件（事件总线的持久化载体）。
//
// 有意差异：Laravel DisplayEventService 把事件列表与「每班自增计数器」放在 Cache
// （display:events:<class_id> / :counter，TTL 600 秒，每班最多 200 条）；Go 端没有 Cache，
// 改为数据库表 display_events，用 (class_id, seq) 表达同一语义——seq 即客户端 last_event_id
// 对应的 id 值，同为「每班自增」。TTL 与条数上限在写入/读取时惰性清理（无定时任务）。
type DisplayEvent struct {
	ID      uint   `gorm:"primaryKey" json:"-"`
	ClassID uint   `gorm:"index;not null" json:"-"`
	Seq     int    `gorm:"index;not null" json:"id"`
	Type    string `gorm:"size:32;not null" json:"type"`
	// Data 事件载荷（JSON 文本，写入时序列化，读取时原样作为嵌套对象输出，同 Laravel poll 结构）。
	Data      string    `gorm:"type:text" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// DisplayEventTTL 事件存活时长（同 Laravel DisplayEventService::CACHE_TTL = 600 秒）。
const DisplayEventTTL = 600 * time.Second

// DisplayEventMaxPerClass 每班最多保留的事件条数（同 Laravel MAX_EVENTS = 200）。
const DisplayEventMaxPerClass = 200
