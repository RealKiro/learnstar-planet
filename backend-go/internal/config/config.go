// Package config 从环境变量加载应用配置。
package config

import (
	"os"
	"strconv"
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

// Load 读取环境变量并返回配置。所有键均有合理默认值，便于零配置启动。
func Load() *Config {
	return &Config{
		Port:        env("PORT", "8080"),
		DBDriver:    env("DB_DRIVER", "sqlite"),
		DBDSN:       env("DB_DSN", "data/learnstar.db"),
		JWTSecret:   env("JWT_SECRET", "learnstar-planet-dev-secret-change-me"),
		JWTExpHours: envInt("JWT_EXP_HOURS", 72),
		BcryptCost:  envInt("BCRYPT_COST", 10),
		AutoMigrate: envBool("AUTO_MIGRATE", true),
		SeedAdmin:   envBool("SEED_ADMIN", true),
		AppTimezone: env("APP_TIMEZONE", "Asia/Shanghai"),

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
