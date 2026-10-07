package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

func main() {
	// 加载 .env（若存在）；不存在时忽略，走环境变量默认值。
	_ = godotenv.Load()

	// CLI：企微请假同步（Laravel `php artisan attendance:sync-wechat-leave` 的等价物）。
	// 命中时不启动 HTTP 服务，执行完打印结果并退出。
	syncWechatLeave := flag.Bool("sync-wechat-work-leave", false,
		"同步企业微信请假到考勤后退出（等价 Laravel `php artisan attendance:sync-wechat-leave`）")
	syncSchoolID := flag.Int("school-id", 0, "只同步该校（0 = 全部 status=active 学校）")
	syncDate := flag.String("date", "", "同步日期 YYYY-MM-DD（缺省今天）")
	flag.Parse()

	cfg := config.Load()

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
