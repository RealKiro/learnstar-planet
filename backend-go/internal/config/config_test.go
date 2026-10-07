package config

// DB 配置解析测试：驱动别名归一化 + Laravel 风格分项键组装 DSN + DB_DSN 优先。

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDriver(t *testing.T) {
	assert.Equal(t, "postgres", normalizeDriver("pgsql"))
	assert.Equal(t, "postgres", normalizeDriver("PostgreSQL"))
	assert.Equal(t, "mysql", normalizeDriver("MySQL"))
	assert.Equal(t, "sqlite", normalizeDriver(""))
	assert.Equal(t, "sqlite", normalizeDriver("sqlite"))
}

func TestDSNFor(t *testing.T) {
	t.Run("sqlite 默认", func(t *testing.T) {
		assert.Equal(t, "data/learnstar.db", dsnFor("sqlite", "Asia/Shanghai"))
	})

	t.Run("DB_DSN 显式优先", func(t *testing.T) {
		t.Setenv("DB_DSN", "file:test.db?cache=shared")
		assert.Equal(t, "file:test.db?cache=shared", dsnFor("sqlite", "Asia/Shanghai"))
	})

	t.Run("mysql 分项键组装（密码特殊字符转义）", func(t *testing.T) {
		t.Setenv("DB_HOST", "192.168.1.50")
		t.Setenv("DB_PORT", "3307")
		t.Setenv("DB_DATABASE", "learnstar")
		t.Setenv("DB_USERNAME", "ls_user")
		t.Setenv("DB_PASSWORD", "p@ss:word/1")
		got := dsnFor("mysql", "Asia/Shanghai")
		assert.Equal(t,
			"ls_user:p%40ss%3Aword%2F1@tcp(192.168.1.50:3307)/learnstar?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
			got)
	})

	t.Run("postgres 分项键组装", func(t *testing.T) {
		t.Setenv("DB_HOST", "10.0.0.8")
		t.Setenv("DB_PORT", "5432")
		t.Setenv("DB_DATABASE", "ls")
		t.Setenv("DB_USERNAME", "pg_user")
		t.Setenv("DB_PASSWORD", "secret")
		got := dsnFor("postgres", "Asia/Shanghai")
		assert.Contains(t, got, "host=10.0.0.8")
		assert.Contains(t, got, "dbname=ls")
		assert.Contains(t, got, "TimeZone=Asia/Shanghai")
	})
}

func TestLoadAliasDBConnection(t *testing.T) {
	t.Setenv("DB_CONNECTION", "pgsql") // 旧 .env 键名
	t.Setenv("DB_HOST", "127.0.0.1")
	cfg := Load()
	assert.Equal(t, "postgres", cfg.DBDriver)
	assert.Contains(t, cfg.DBDSN, "host=127.0.0.1")
}
