// Package database 负责数据库连接、自动迁移与初始种子数据。
package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 按配置连接数据库并执行自动迁移。
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "sqlite":
		dsn := cfg.DBDSN
		if dsn != ":memory:" {
			if err := os.MkdirAll(filepath.Dir(dsn), 0o755); err != nil {
				return nil, fmt.Errorf("创建数据库目录失败: %w", err)
			}
		}
		dialector = sqlite.Open(dsn)
	case "mysql":
		dialector = mysql.Open(cfg.DBDSN)
	case "postgres":
		dialector = postgres.Open(cfg.DBDSN)
	default:
		return nil, fmt.Errorf("不支持的 DB_DRIVER: %q（可选 sqlite/mysql/postgres）", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	if cfg.AutoMigrate {
		if err := Migrate(db); err != nil {
			return nil, err
		}
	}

	if cfg.SeedAdmin {
		if err := Seed(db, cfg); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Migrate 自动迁移全部模型。
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.School{},
		&models.User{},
		&models.ClassRoom{},
		&models.Student{},
		&models.Pet{},
		&models.ScoreRule{},
		&models.Score{},
		&models.ScoreLog{},
		&models.Notice{},
		&models.ShopItem{},
		&models.ShopRedemption{},
		&models.TimetableTeacherAssignment{},
		&models.TimetableTeacherUnavailability{},
		&models.Attendance{},
		&models.Wallet{},
		&models.ExchangeRate{},
		&models.ExchangeLog{},
		&models.ClassRoomTeacher{},
		&models.Broadcast{},
		&models.Subject{},
		&models.ClassPeriod{},
		&models.TimetableEntry{},
		&models.TimetableChangeRequest{},
		&models.DisplayToken{},
		&models.DisplayLoginLog{},
		&models.DisplayEvent{},
		&models.PetFreePick{},
		&models.PetCollection{},
		&models.ThirdPartyBinding{},
		&models.RevokedToken{},
		&models.TempBindingContext{},
		&models.WechatWorkToken{},
		&models.WechatWorkLeaveRecord{},
		&models.WechatWorkToken{},
		&models.AISetting{},
		&models.AIConversation{},
	)
}

// Seed 幂等地创建默认学校、管理员与（可选）API 机器人账号。
func Seed(db *gorm.DB, cfg *config.Config) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1) 默认学校
		var school models.School
		err := tx.Where("code = ?", cfg.SchoolCode).First(&school).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			school = models.School{
				Name: cfg.SchoolName,
				Code: cfg.SchoolCode,
			}
			if err := tx.Create(&school).Error; err != nil {
				return fmt.Errorf("创建默认学校失败: %w", err)
			}
		} else if err != nil {
			return err
		}

		// 2) 默认管理员
		if err := seedAdmin(tx, cfg, school.ID); err != nil {
			return err
		}

		// 3) 可选 API 机器人教师
		if cfg.BotEnabled {
			if err := seedBot(tx, cfg, school.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

func seedAdmin(tx *gorm.DB, cfg *config.Config, schoolID uint) error {
	// 与 Laravel AdminUserSeeder 同口径：按用户名查找；已存在则**只同步密码**
	//（.env 是唯一真相来源，改 ADMIN_PASSWORD 重启即生效——「忘了密码」的官方
	// 恢复路径依赖它）；name/状态等其余字段不覆盖，可通过后台修改。
	var admin models.User
	err := tx.Where("username = ?", cfg.AdminUsername).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return createUser(tx, schoolID, "school_admin", cfg.AdminUsername, cfg.AdminPassword, cfg.AdminName, false)
	}
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), cfg.BcryptCost)
	if err != nil {
		return err
	}
	admin.PasswordHash = string(hash)
	return tx.Save(&admin).Error
}

func seedBot(tx *gorm.DB, cfg *config.Config, schoolID uint) error {
	var bot models.User
	err := tx.Where("username = ?", cfg.BotUsername).First(&bot).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return createUser(tx, schoolID, "teacher", cfg.BotUsername, cfg.BotPassword, cfg.BotName, true)
	}
	if err != nil {
		return err
	}

	// 已存在则同步密码与机器人标记（.env 是唯一真相来源，同 Laravel BotTeacherSeeder）。
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.BotPassword), cfg.BcryptCost)
	if err != nil {
		return err
	}
	bot.IsAPIBot = true
	bot.Name = cfg.BotName
	bot.PasswordHash = string(hash)
	bot.PlainPassword = cfg.BotPassword
	bot.PasswordChanged = true
	return tx.Save(&bot).Error
}

func createUser(tx *gorm.DB, schoolID uint, role, username, password, name string, isBot bool) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}

	u := models.User{
		SchoolID:        schoolID,
		Role:            role,
		Username:        username,
		PasswordHash:    string(hash),
		Name:            name,
		Nickname:        name,
		Status:          "active",
		IsAPIBot:        isBot,
		PasswordChanged: isBot,
	}
	// 明文密码：同 Laravel，只有 BotTeacherSeeder 写 plain_password（AdminUserSeeder 不写）。
	if isBot {
		u.PlainPassword = password
	}
	if err := tx.Create(&u).Error; err != nil {
		return fmt.Errorf("创建账号 %s 失败: %w", username, err)
	}
	return nil
}
