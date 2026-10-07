// Package util 提供业务时间工具。全站业务时间统一使用亚洲/上海时区，
// 与旧 Laravel 端 APP_TIMEZONE=Asia/Shanghai 的口径一致。
package util

import (
	"time"
)

// Loc 是业务时区（Asia/Shanghai）。加载失败时回退到本地时区。
var Loc = mustLoadLocation("Asia/Shanghai")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Local
	}
	return loc
}

// Now 返回业务时区下的当前时间。
func Now() time.Time {
	return time.Now().In(Loc)
}

// StartOfDay 返回 t 所在业务日的零点。
func StartOfDay(t time.Time) time.Time {
	tt := t.In(Loc)
	return time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, Loc)
}

// StartOfWeek 返回 t 所在自然周的周一零点（与 Laravel startOfWeek 默认周一一致）。
func StartOfWeek(t time.Time) time.Time {
	tt := StartOfDay(t)
	// time.Weekday: Sunday=0 ... Saturday=6，换算为周一=0。
	offset := (int(tt.Weekday()) + 6) % 7
	return tt.AddDate(0, 0, -offset)
}

// Today 返回业务时区下的今日日期字符串（"2006-01-02"）。
func Today() string {
	return Now().Format("2006-01-02")
}
