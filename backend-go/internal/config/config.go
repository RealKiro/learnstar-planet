// Package config 从环境变量加载应用配置。
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config 保存应用运行所需的环境配置。
type Config struct {
	Port        string
	DBDriver    string // sqlite | mysql | postgres
	DBDSN       string
	JWTSecret   string
	JWTExpHours int
	BcryptCost  int
	AutoMigrate bool
	SeedAdmin   bool
	AppTimezone string
	PublicDir   string // SPA 静态资源目录（构建后的前端产物）

	// 种子数据（首次启动时创建默认学校/管理员）。
	SchoolName    string
	SchoolCode    string
	AdminUsername string
	AdminPassword string
	AdminName     string
	BotEnabled    bool
	BotUsername   string
	BotPassword   string
	BotName       string
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// normalizeDriver 归一化数据库驱动别名（pgsql/postgresql → postgres）。
func normalizeDriver(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	switch v {
	case "pgsql", "postgresql":
		return "postgres"
	case "":
		return "sqlite"
	default:
		return v
	}
}

// dsnFor 解析数据库连接串：
//   - DB_DSN 显式配置时优先（一行搞定，支持任意特殊字符）；
//   - 否则按 Laravel 风格分项键组装（DB_HOST/DB_PORT/DB_DATABASE/DB_USERNAME/
//     DB_PASSWORD），让旧部署的 .env 平滑迁移；
//   - SQLite 缺省落到 data/learnstar.db（Docker 里 data 目录挂了数据卷）。
//
// 注意：MySQL 用户名/密码按 DSN 规范做百分号转义；PostgreSQL 为 key=value 形式，
// 密码含空格/等号的极端场景请直接配置 DB_DSN。
func dsnFor(driver, tz string) string {
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		return dsn
	}
	switch driver {
	case "mysql":
		user := url.QueryEscape(env("DB_USERNAME", "learnstar"))
		pass := url.QueryEscape(env("DB_PASSWORD", ""))
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
			user, pass,
			env("DB_HOST", "127.0.0.1"), env("DB_PORT", "3306"),
			env("DB_DATABASE", "learnstar"),
			url.QueryEscape(tz))
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
			env("DB_HOST", "127.0.0.1"), env("DB_PORT", "5432"),
			env("DB_USERNAME", "postgres"), env("DB_PASSWORD", ""),
			env("DB_DATABASE", "learnstar"), tz)
	default: // sqlite
		return "data/learnstar.db"
	}
}

// Load 读取环境变量并返回配置。所有键均有合理默认值，便于零配置启动。
func Load() *Config {
	driver := normalizeDriver(env("DB_DRIVER", env("DB_CONNECTION", "sqlite")))
	tz := env("APP_TIMEZONE", "Asia/Shanghai")
	return &Config{
		Port:        env("PORT", "8080"),
		DBDriver:    driver,
		DBDSN:       dsnFor(driver, tz),
		JWTSecret:   env("JWT_SECRET", "learnstar-planet-dev-secret-change-me"),
		JWTExpHours: envInt("JWT_EXP_HOURS", 72),
		BcryptCost:  envInt("BCRYPT_COST", 10),
		AutoMigrate: envBool("AUTO_MIGRATE", true),
		SeedAdmin:   envBool("SEED_ADMIN", true),
		AppTimezone: tz,
		PublicDir:   env("PUBLIC_DIR", "public"),

		SchoolName:    env("SCHOOL_NAME", "学宠星球"),
		SchoolCode:    env("SCHOOL_CODE", "learnstar"),
		AdminUsername: env("ADMIN_USERNAME", "admin"),
		AdminPassword: env("ADMIN_PASSWORD", "admin123456"),
		AdminName:     env("ADMIN_NAME", "学校管理员"),
		BotEnabled:    envBool("BOT_ENABLED", false),
		BotUsername:   env("BOT_USERNAME", "api-bot"),
		BotPassword:   env("BOT_PASSWORD", "learnstar-bot-2026"),
		BotName:       env("BOT_NAME", "API 机器人"),
	}
}
