// 时区纪律守护测试（跨平台，不依赖宿主机时区）。
//
// 背景：glebarez/sqlite 把 time.Time 以**带偏移的 RFC3339 文本**落库，因此 SQL 里的
// `expires_at > ?` 实际是**文本比较**而非时刻比较。同一列的写入与比较必须落在同一
// 钟面域，否则比较结果与真实先后无关——UTC 域的 "08:16+00:00" 会被判为「小于」
// 上海域的 "14:16+08:00"，尽管它其实晚了 6 小时。
//
// 这正是 GH Actions（TZ=UTC）上 TestAuthLogoutAndRefreshInvalidateOldToken 失败的根因：
// JWT 的 NumericDate 由 time.Unix 解析，其 Location 随宿主机 / 部署环境变化（UTC 主机 =
// UTC 域，+08 开发机 = Local 域）；而 revoked_tokens 的惰性清理原先用 util.Now()
// （恒为 Asia/Shanghai）比较——在 UTC 主机上把「未来 72h 才过期」的刚写入行当成已过期删掉，
// 登出 / 刷新随即失效（安全回归，不只是测试问题）。
//
// 修法：本表全链路统一 UTC 域（写入 .UTC()、比较 time.Now().UTC()），见
// services/auth.go RevokeToken 的注释。display_tokens 同因一并对齐 UTC——该列有**两个
// 写入方**（display.go 的 issueToken 曾用 time.Now()，class_login.go 的 ClassLogin 用
// util.Now()），域不一致时文本比较同样错序。
//
// 本文件是「守护测试」：故意在任意宿主机上复现 UTC 域条件，一旦有人把域改回去就会立刻红。
package services_test

import (
	"strings"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// storedUTCOffset 读出某表某行 expires_at 的**裸落库文本**，解析后返回其时区偏移（秒）。
//
// 之所以绕开 time.Time 扫描：驱动会把文本还原成带原始 Location 的 time.Time，
// 时刻相同、域不同时扫描结果看起来一模一样，看不出问题；只有裸文本能暴露域。
func storedUTCOffset(t *testing.T, db *gorm.DB, table, keyCol, key string) int {
	t.Helper()

	var raw string
	require.NoError(t, db.Raw(
		"SELECT CAST(expires_at AS TEXT) FROM "+table+" WHERE "+keyCol+" = ?", key).Scan(&raw).Error)
	require.NotEmpty(t, raw, "%s 中未找到 %s = %s 的行", table, keyCol, key)

	// glebarez/sqlite 落库用的是「空格分隔」的类 RFC3339 文本
	// （如 "2026-10-08 15:13:57.8290404+08:00"），换掉首个空格即可交给 RFC3339Nano 解析。
	if parsed, err := time.Parse(time.RFC3339Nano, strings.Replace(raw, " ", "T", 1)); err == nil {
		_, offset := parsed.Zone()
		return offset
	}
	t.Fatalf("无法解析 %s.%s 的落库时区文本，裸值为 %q", table, keyCol, raw)
	return 0
}

// TestRevokeTokenExpiryUsesUTCDomain 撤销名单（revoked_tokens）全链路 UTC 域。
func TestRevokeTokenExpiryUsesUTCDomain(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedAdminUser(t, db, school.ID, "secret123")
	svc := services.NewAuthService(db, auth.New("test-secret", 72))

	// 显式 UTC 域：等价于 TZ=UTC 环境里 time.Unix 解析出的 JWT 过期时间（即 CI 条件）。
	//
	// 用「1 小时后」而不是「72 小时后」很关键：文本比较逐字符进行，过期时间一旦跨过一天，
	// 日期段先决出大小、时间段的错序被掩盖（72h 时本测试在 +08 宿主机上抓不到 bug，
	// 只有 UTC 宿主机才红）。取同一日历日内的 1 小时，任何宿主机都能复现 CI 的错序。
	exp := time.Now().Add(time.Hour).UTC()
	require.NoError(t, svc.RevokeToken("jti-utc", &exp))

	// ① 列内文本必须是 UTC 域——域一旦漂回 Local，下面的行为断言在非 UTC 宿主机上会失效。
	assert.Equal(t, 0, storedUTCOffset(t, db, "revoked_tokens", "jti", "jti-utc"),
		"revoked_tokens.expires_at 必须以 UTC 域落库，否则文本比较会错序")

	// ② 未过期的行不得被惰性清理误删（Logout 内部先写新行、再清理过期行）。
	require.NoError(t, svc.Logout("fresh", nil))
	var kept int64
	require.NoError(t, db.Model(&models.RevokedToken{}).Where("jti = ?", "jti-utc").Count(&kept).Error)
	assert.Equal(t, int64(1), kept, "未过期的撤销记录被惰性清理误删 → 登出/刷新将失效")

	// ③ 中间件口径（expires_at > time.Now().UTC()）必须命中该 jti。
	var hit int64
	require.NoError(t, db.Model(&models.RevokedToken{}).
		Where("jti = ? AND expires_at > ?", "jti-utc", time.Now().UTC()).Count(&hit).Error)
	assert.Equal(t, int64(1), hit, "中间件应把该 jti 判为仍在撤销名单内")

	// ④ 真过期的行仍要被清掉（别为了修 ② 把惰性清理废掉）。
	require.NoError(t, db.Create(&models.RevokedToken{
		JTI: "jti-old", ExpiresAt: time.Now().UTC().Add(-time.Hour),
	}).Error)
	require.NoError(t, svc.Logout("fresh-2", nil))
	var pruned int64
	require.NoError(t, db.Model(&models.RevokedToken{}).Where("jti = ?", "jti-old").Count(&pruned).Error)
	assert.Equal(t, int64(0), pruned, "已过期的撤销记录必须被惰性清理")
}

// TestDisplayTokenExpiryUsesUTCDomain 大屏/班级码 token（display_tokens）两个写入方同域。
func TestDisplayTokenExpiryUsesUTCDomain(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)
	seedStudent(t, db, f.Class.ID)

	// 先造一条已过期的行：登录时的惰性清理应删掉它，但绝不能碰刚签发的行。
	require.NoError(t, db.Create(&models.DisplayToken{
		Token: "disp_stale", ClassID: f.Class.ID,
		ExpiresAt: time.Now().UTC().Add(-time.Hour),
	}).Error)

	result, err := f.Svc.Login("LS11", "10.0.0.1", "ua")
	require.NoError(t, err)

	var stale int64
	require.NoError(t, db.Model(&models.DisplayToken{}).Where("token = ?", "disp_stale").Count(&stale).Error)
	assert.Equal(t, int64(0), stale, "已过期的 display_tokens 行应被惰性清理")

	assert.Equal(t, 0, storedUTCOffset(t, db, "display_tokens", "token", result.Token),
		"display.go 签发的 display_tokens 行必须以 UTC 域落库")

	classID, err := f.Svc.ValidateToken(result.Token)
	require.NoError(t, err, "刚签发的 disp_ token 被登录时的惰性清理误删")
	assert.Equal(t, f.Class.ID, classID)

	// ClassLogin 写的是同一张表，域必须与 display.go 一致。
	classResult, err := f.Svc.ClassLogin("LS11")
	require.NoError(t, err)
	assert.Equal(t, 0, storedUTCOffset(t, db, "display_tokens", "token", classResult.Token),
		"class_login 签发的 display_tokens 行也必须以 UTC 域落库")

	classID, err = f.Svc.ValidateToken(classResult.Token)
	require.NoError(t, err, "刚签发的 class_ token 被写入时的惰性清理误删")
	assert.Equal(t, f.Class.ID, classID)
}
