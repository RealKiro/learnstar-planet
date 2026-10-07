package database

// Seed 行为测试：管理员「.env 唯一真相来源」语义（改 ADMIN_PASSWORD 重启即生效）。

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func newSeedCfg(adminUser, adminPass string) *config.Config {
	return &config.Config{
		SchoolName: "种子学校", SchoolCode: "seed-school",
		AdminUsername: adminUser, AdminPassword: adminPass, AdminName: "管理员",
		BcryptCost: bcrypt.MinCost,
		SeedAdmin:  true,
	}
}

func TestSeedAdminSyncsPasswordFromEnv(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), GormConfig())
	require.NoError(t, err)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, Migrate(db))

	// 首次播种
	require.NoError(t, Seed(db, newSeedCfg("boss", "first-pass")))

	var admin models.User
	require.NoError(t, db.Where("username = ?", "boss").First(&admin).Error)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte("first-pass")))

	// 用户改过名字（模拟后台修改）
	require.NoError(t, db.Model(&admin).Update("name", "李校长").Error)

	// 「改 .env 重启」：同用户名新密码 → 只同步密码，不覆盖 name
	cfg := newSeedCfg("boss", "second-pass")
	cfg.SchoolCode = "seed-school" // 学校已存在，走 FirstOrCreate 已有分支
	require.NoError(t, Seed(db, cfg))

	var again models.User
	require.NoError(t, db.Where("username = ?", "boss").First(&again).Error)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(again.PasswordHash), []byte("second-pass")))
	assert.Equal(t, "李校长", again.Name, "seedAdmin 只同步密码，不得覆盖 name")

	// 换用户名 = 新建管理员，不动旧账号
	cfg2 := newSeedCfg("boss2", "another-pass")
	cfg2.SchoolCode = "seed-school"
	require.NoError(t, Seed(db, cfg2))
	var count int64
	require.NoError(t, db.Model(&models.User{}).Where("role = ?", "school_admin").Count(&count).Error)
	assert.EqualValues(t, 2, count)
}
