// Package handlers 将 HTTP 请求转换为服务调用并统一 JSON 输出。
package handlers

import (
	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"gorm.io/gorm"
)

// Handlers 聚合全部业务服务，供路由注册使用。
type Handlers struct {
	auth        *services.AuthService
	admin       *services.Admin
	teacher     *services.Teacher
	students    *services.StudentService
	rules       *services.Rules
	scores      *services.ScoreService
	pets        *services.PetService
	leaderboard *services.LeaderboardService
	scope       *services.Scope
	notice      *services.NoticeService
	shop        *services.ShopService
	attendance  *services.AttendanceService
	broadcast   *services.BroadcastService
	currency    *services.CurrencyService
	timetable   *services.TimetableService
	classroom   *services.ClassroomService
	messaging   *services.ClassroomMessagingService
	reports     *services.ReportService
	display     *services.DisplayService
	petSeries   *services.PetSeriesService
	pk          *services.PkService
	ai          *services.AIAssistant
	aiAdmin     *services.AIAdmin
	aiBilling   *services.AIBilling
	ops         *services.AdminOps
	thirdParty  *services.ThirdParty
	// wechatWork / wechatAttendance 企微回调（验签解密）与请假同步（审批 → 考勤）。
	wechatWork       *services.WechatWorkService
	wechatAttendance *services.WechatWorkAttendanceService
}

// New 构造 Handler 聚合，并注入 JWT 管理器。
func New(db *gorm.DB, jwtMgr *jwtauth.Manager) *Handlers {
	scope := services.NewScope(db)
	scores := services.NewScoreService(db)
	currency := services.NewCurrencyService(db)
	aiService := services.NewAIService()
	// 企微服务只建一次：回调处理器读 token / encoding_aes_key，考勤服务走同一实例（同 Laravel
	// 容器单例语义；服务构造时读环境变量，故测试须在 router.New 之前注入 WECHAT_WORK_* 变量）。
	wechatWork := services.NewWechatWorkService(db)
	wechatAttendance := services.NewWechatWorkAttendanceService(db, wechatWork)

	return &Handlers{
		auth:             services.NewAuthService(db, jwtMgr),
		admin:            services.NewAdmin(db),
		teacher:          services.NewTeacher(db, scope),
		students:         services.NewStudentService(db),
		rules:            services.NewRules(db),
		scores:           scores,
		pets:             services.NewPetService(db, scope),
		leaderboard:      services.NewLeaderboardService(db),
		scope:            scope,
		notice:           services.NewNoticeService(db, scope),
		shop:             services.NewShopService(db, scope, scores, currency),
		attendance:       services.NewAttendanceService(db, scope, wechatAttendance),
		broadcast:        services.NewBroadcastService(db, scope),
		currency:         currency,
		timetable:        services.NewTimetableService(db),
		classroom:        services.NewClassroomService(db, scope),
		messaging:        services.NewClassroomMessagingService(db, scope),
		reports:          services.NewReportService(db, scope),
		display:          services.NewDisplayService(db),
		petSeries:        services.NewPetSeriesService(db, scope),
		ai:               services.NewAIAssistant(db, aiService),
		aiAdmin:          services.NewAIAdmin(db, aiService),
		pk:               services.NewPkService(db, scope),
		aiBilling:        services.NewAIBilling(db, aiService),
		ops:              services.NewAdminOps(db),
		thirdParty:       services.NewThirdPartyService(db, jwtMgr),
		wechatWork:       wechatWork,
		wechatAttendance: wechatAttendance,
	}
}
