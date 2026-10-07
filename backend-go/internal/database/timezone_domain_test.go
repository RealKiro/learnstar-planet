// GORM 自动时间戳的时区守护测试（跨平台，不依赖宿主机时区）。
//
// 背景：GORM 的默认 NowFunc 是 time.Now().Local()（gorm.io/gorm/gorm.go 的注释明说
// "defaults to time.Now().Local()"），随宿主机与容器 TZ 漂移；而本仓全部业务时间窗口
// （今日积分 / 本周榜 / 日报表 / 月报 / 大屏最近积分…）都由 util.Now()、util.StartOfDay、
// util.StartOfWeek 生成，域恒为 Asia/Shanghai。SQLite 对 time 列做**文本比较**
// （见 services/auth.go RevokeToken 的时区纪律说明），两者不同域时跨日的比较会错序：
// alpine 运行时容器（无 TZ → Local=UTC）里「今天的积分」会被判成不在今天而少算一整天。
// 故 GormConfig() 把自动时间戳钉到业务时区，本文件守护这条约束。
//
// ⚠️ 本测试**主动把 time.Local 换成 UTC**，从而在任何开发机上复现「CI / alpine 容器」的
// 错配条件（宿主域 = UTC，业务域 = Asia/Shanghai）。之所以不用 TZ 环境变量：
// **Windows 上的 Go 运行时不读 TZ**（实测 TZ=UTC 时 time.Now() 仍为 +08:00），
// 只靠环境变量无法在本机覆盖这类缺陷——这也正是它历史上只在 CI 上红的原因。
package database

import (
	"strings"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// rawOffset 读 schools 表某行的 time 列**裸落库文本**并返回其时区偏移（秒）。
//
// 绕开 time.Time 扫描：驱动会把文本还原成带原始 Location 的 time.Time，
// 时刻相同、域不同时扫描结果看起来一模一样，看不出域差异；只有裸文本能暴露它
// （与 services 包 timezone_domain_test.go 同一手法）。
func rawOffset(t *testing.T, db *gorm.DB, column string, id uint) int {
	t.Helper()

	var raw string
	require.NoError(t, db.Raw("SELECT CAST("+column+" AS TEXT) FROM schools WHERE id = ?", id).
		Scan(&raw).Error)
	require.NotEmpty(t, raw, "schools.%s 缺失", column)

	// 落库文本是「空格分隔」的类 RFC3339（如 "2026-10-08 15:13:57.8290404+08:00"）。
	parsed, err := time.Parse(time.RFC3339Nano, strings.Replace(raw, " ", "T", 1))
	require.NoError(t, err, "无法解析 schools.%s 的落库文本 %q", column, raw)

	_, offset := parsed.Zone()
	return offset
}

// TestAutoTimestampsUseBusinessTimezone 自动写入的 created_at / updated_at 必须与 util.Now() 同域。
func TestAutoTimestampsUseBusinessTimezone(t *testing.T) {
	// 把宿主时区强制为 UTC：GORM 的旧默认 NowFunc 正是 time.Now().Local()，
	// 于是「旧默认」在本机也会写出 UTC 域时间戳，断言就能在任何机器上把它抓住。
	// util.Now() 走独立的 util.Loc，不受本行影响——错配正是这样构成的。
	// 本测试不并行且退出时恢复，故不会污染同包其他测试。
	origLocal := time.Local
	time.Local = time.UTC
	t.Cleanup(func() { time.Local = origLocal })

	db, err := gorm.Open(sqlite.Open(":memory:"), GormConfig())
	require.NoError(t, err)
	require.NoError(t, Migrate(db))

	// 自检：错配条件确实已成立（宿主域必须与业务域不同），否则下面的断言没有意义。
	_, hostOffset := time.Now().Zone()
	_, wantOffset := util.Now().Zone()
	require.Equal(t, 0, hostOffset, "本测试需要宿主时区为 UTC 才能复现 CI 条件")
	require.Equal(t, 8*3600, wantOffset, "业务时区应为 UTC+8（Asia/Shanghai），util.Loc 配置异常")

	// 不显式给时间：让 GORM 走 NowFunc 自动写入 created_at / updated_at。
	school := models.School{Name: "时区测试学校", Code: "tz-probe"}
	require.NoError(t, db.Create(&school).Error)

	assert.Equal(t, wantOffset, rawOffset(t, db, "created_at", school.ID),
		"GORM 自动写入的 created_at 必须与业务时区同域，否则跨日的文本比较会错序")
	assert.Equal(t, wantOffset, rawOffset(t, db, "updated_at", school.ID),
		"GORM 自动写入的 updated_at 必须与业务时区同域")

	// 行为面复核：刚写入的行必须落在「今日」窗口内——各报表/周榜用的就是这种 util 窗口。
	var today int64
	require.NoError(t, db.Model(&models.School{}).
		Where("id = ? AND created_at >= ?", school.ID, util.StartOfDay(util.Now())).
		Count(&today).Error)
	assert.Equal(t, int64(1), today, "刚写入的行被「今日」窗口漏掉 → 今日积分/日报表会少算")
}
