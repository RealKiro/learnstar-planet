package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"

	// 内嵌 IANA 时区数据库：保证在无 tzdata 的最小镜像/系统里
	// util 的 Asia/Shanghai（考勤、今日积分、周榜的「今天」边界）仍然可用。
	_ "time/tzdata"
)

// devDefaultJWTSecret 与 config.Load 的默认值保持一致，
// 用于识别「用户没有显式配置 JWT_SECRET」的情形。
const devDefaultJWTSecret = "learnstar-planet-dev-secret-change-me"

func main() {
	// 加载 .env（若存在）；不存在时忽略，走环境变量默认值。
	_ = godotenv.Load()

	// CLI：企微请假同步（原 `php artisan attendance:sync-wechat-leave` 的等价物）。
	// 命中时不启动 HTTP 服务，执行完打印结果并退出。
	syncWechatLeave := flag.Bool("sync-wechat-work-leave", false,
		"同步企业微信请假到考勤后退出（等价 Laravel `php artisan attendance:sync-wechat-leave`）")
	syncSchoolID := flag.Int("school-id", 0, "只同步该校（0 = 全部 status=active 学校）")
	syncDate := flag.String("date", "", "同步日期 YYYY-MM-DD（缺省今天）")
	flag.Parse()

	cfg := config.Load()
	applyPersistentJWTSecret(cfg)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if *syncWechatLeave {
		runWechatWorkLeaveSync(db, *syncSchoolID, *syncDate)
		return
	}

	r := router.New(db, cfg)
	log.Printf("学宠星球 Go 后端已启动，监听端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// applyPersistentJWTSecret 在用户未显式配置 JWT_SECRET（未设置或仍为开发默认值）时，
// 自动生成 64 位 hex 随机密钥并持久化到数据目录的 .jwt_secret 文件——等价 Laravel
// 把 APP_KEY 持久化到数据卷的做法，避免「容器重建全员登出」与「公开默认密钥可伪造
// 令牌」两个问题。显式配置了非默认 JWT_SECRET 时本函数不做任何事。
func applyPersistentJWTSecret(cfg *config.Config) {
	if v := strings.TrimSpace(os.Getenv("JWT_SECRET")); v != "" && v != devDefaultJWTSecret {
		return // 用户显式配置，尊重之
	}

	path := filepath.Join(jwtSecretDir(cfg), ".jwt_secret")
	if b, err := os.ReadFile(path); err == nil {
		if secret := strings.TrimSpace(string(b)); len(secret) >= 32 {
			cfg.JWTSecret = secret
			return
		}
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		// 无法生成随机数属于极端环境异常；退回默认密钥并大声提醒。
		log.Printf("⚠️ 无法生成随机 JWT_SECRET：%v（继续使用开发默认值，请手动配置 JWT_SECRET）", err)
		return
	}
	secret := hex.EncodeToString(buf)
	cfg.JWTSecret = secret

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err == nil {
		if err := os.WriteFile(path, []byte(secret), 0o600); err == nil {
			log.Printf("已自动生成 JWT_SECRET 并保存到 %s（容器请把该目录挂到数据卷）", path)
		} else {
			log.Printf("⚠️ JWT_SECRET 生成成功但写入 %s 失败：%v（容器重建后全员需重新登录）", path, err)
		}
	}
}

// jwtSecretDir 决定 .jwt_secret 的落盘位置：SQLite 落数据库同目录（即数据卷）；
// 外置数据库落在 UPLOAD_DIR；都不可用时退回当前目录。
func jwtSecretDir(cfg *config.Config) string {
	if cfg.DBDriver == "sqlite" && cfg.DBDSN != "" && cfg.DBDSN != ":memory:" {
		if dir := filepath.Dir(cfg.DBDSN); dir != "" && dir != "." {
			return dir
		}
	}
	if dir := strings.TrimSpace(os.Getenv("UPLOAD_DIR")); dir != "" {
		return dir
	}
	return "."
}

// runWechatWorkLeaveSync 执行企微请假同步 CLI，并按 `{synced, skipped, failed}` 输出结果。
//
// 单项同步（-school-id=N）时 failed 恒为 0（单校失败直接以错误退出，同 artisan 的异常路径）。
func runWechatWorkLeaveSync(db *gorm.DB, schoolID int, date string) {
	if date == "" {
		date = util.Today()
	}

	svc := services.NewWechatWorkAttendanceService(db, nil)

	result := services.WechatWorkLeaveSyncResult{}
	if schoolID > 0 {
		one, err := svc.SyncForSchool(uint(schoolID), date)
		if err != nil {
			log.Fatalf("企微请假同步失败 学校id=%d 日期=%s: %v", schoolID, date, err)
		}
		result = one
	} else {
		all, err := svc.SyncAll(date)
		if err != nil {
			log.Fatalf("企微请假同步失败 日期=%s: %v", date, err)
		}
		result = all
	}

	buf, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("结果序列化失败: %v", err)
	}
	fmt.Printf("%s\n", buf)
}
