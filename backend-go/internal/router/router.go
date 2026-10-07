// Package router 注册 HTTP 路由与中间件。
package router

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/handlers"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
)

// storageURLPrefix 是落库的上传文件 URL 前缀（与 Laravel 时代保持一致，
// 前端模板里已有该字面量）；目录部分由 UPLOAD_DIR 承载。
const storageURLPrefix = "/storage/app/uploads/"

// spaFallback 构建未命中路由的兜底处理器。
func spaFallback(cfg *config.Config) gin.HandlerFunc {
	publicDir := strings.TrimSpace(cfg.PublicDir)
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "storage/app/uploads"
	}
	publicFS := http.Dir(publicDir) // 目录不存在时 Open 直接报错，走 index/404 兜底

	return func(c *gin.Context) {
		isAPI := strings.HasPrefix(c.Request.URL.Path, "/api/")
		if isAPI || (c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead) {
			c.JSON(http.StatusNotFound, gin.H{"message": "资源不存在"})
			return
		}

		// 上传文件：/storage/app/uploads/<相对路径> → UPLOAD_DIR/<相对路径>
		if strings.HasPrefix(c.Request.URL.Path, storageURLPrefix) {
			rel := path.Clean("/" + strings.TrimPrefix(c.Request.URL.Path, storageURLPrefix))
			c.File(filepath.Join(uploadDir, rel))
			return
		}

		if publicDir == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "资源不存在"})
			return
		}

		// 静态文件存在且非目录 → 直接返回（http.Dir 自带目录穿越防护）。
		upath := path.Clean("/" + c.Request.URL.Path)
		if f, err := publicFS.Open(upath); err == nil {
			st, serr := f.Stat()
			f.Close()
			if serr == nil && !st.IsDir() {
				c.File(filepath.Join(publicDir, upath))
				return
			}
		}

		// SPA 前端路由兜底：未命中的页面路径回落 index.html。
		indexPath := filepath.Join(publicDir, "index.html")
		if st, err := os.Stat(indexPath); err == nil && !st.IsDir() {
			c.File(indexPath)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "资源不存在"})
	}
}

// New 构建 gin 引擎并注册全部路由。
func New(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	jwtMgr := jwtauth.New(cfg.JWTSecret, cfg.JWTExpHours)
	h := handlers.New(db, jwtMgr)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// 认证（无需登录）

	// 班级码登录（学生端 / 班级大屏统一入口，无需登录，签发 class_ token）
	// throttle:10,1（Laravel routes/api.php 第 23 行）。
	api.POST("/auth/class/login", middleware.Throttle(10, time.Minute), h.ClassLogin)

	// 账号密码登录别名（公开，对应 Laravel routes/api.php 第 21-22 行；均为角色限定入口）
	// throttle:6,1（同 Laravel）。
	api.POST("/auth/teacher/login", middleware.Throttle(6, time.Minute), h.TeacherLogin)
	api.POST("/auth/admin/login", middleware.Throttle(6, time.Minute), h.AdminLogin)

	// 第三方扫码登录（公开，对应 Laravel routes/api.php 第 24-32 行）
	api.GET("/auth/third-party/auth-url", h.ThirdPartyAuthURL)
	api.GET("/auth/third-party/options", h.ThirdPartyOptions)
	api.POST("/auth/third-party/login", h.ThirdPartyLogin)
	api.POST("/auth/teacher/login/wechat", h.TeacherLoginWithWechat)
	api.POST("/auth/teacher/login/wechat-work", h.TeacherLoginWithWechatWork)
	api.POST("/auth/teacher/login/qq", h.TeacherLoginWithQQ)
	api.POST("/auth/teacher/login/renren", h.TeacherLoginWithRenren)
	api.POST("/auth/teacher/bind-after-scan", h.BindAfterScan)

	// 认证（需登录）
	auth := api.Group("/auth", middleware.Auth(jwtMgr, db))
	auth.POST("/change-password", h.ChangePassword)
	// 登录后绑定 / 解绑第三方账号（Laravel routes/api.php 第 38-39 行）
	auth.POST("/bind/:platform", h.BindThirdParty)
	auth.DELETE("/unbind/:platform", h.UnbindThirdParty)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.RefreshToken)
	auth.GET("/bindings", h.GetBindings)

	// 管理端（学校管理员）
	admin := api.Group("/admin", middleware.Auth(jwtMgr, db), middleware.RequireRole("school_admin"))
	admin.GET("/school", h.AdminSchool)
	// Laravel 为 match(['put','post'], 'school') → 两个方法都注册（前端用 PUT）。
	admin.PUT("/school", h.AdminUpdateSchool)
	admin.POST("/school", h.AdminUpdateSchool)

	admin.GET("/classes", h.AdminListClasses)
	admin.POST("/classes", h.AdminCreateClass)
	admin.PUT("/classes/:id", h.AdminUpdateClass)
	admin.DELETE("/classes/:id", h.AdminDeleteClass)

	admin.PUT("/students/:id", h.AdminUpdateStudent)
	admin.DELETE("/students/:id", h.AdminDeleteStudent)

	admin.GET("/teachers", h.AdminListTeachers)
	admin.POST("/teachers", h.AdminCreateTeacher)
	admin.PUT("/teachers/:id", h.AdminUpdateTeacher)
	admin.POST("/teachers/:id/reset-password", h.AdminResetPassword)
	// Laravel：DELETE teachers/{id} → disableTeacher（解除班级关联 + 硬删除；机器人账号 403）
	admin.DELETE("/teachers/:id", h.AdminDeleteTeacher)

	// 管理端：学生分页列表 / 单建学生 / 班级详情 / 教师班级分配 / 查看教师密码
	//（Laravel routes/api.php 第 57、59、86、102-103 行）
	admin.GET("/students", h.AdminStudentsList)
	admin.POST("/students", h.AdminStudentCreate)
	admin.GET("/classes/:id", h.AdminClassShow)
	admin.PUT("/teachers/:id/classes", h.AdminTeacherAssignClasses)
	admin.GET("/teachers/:id/password", h.AdminTeacherPassword)

	// 教师端
	teacher := api.Group("/teacher",
		middleware.Auth(jwtMgr, db),
		middleware.RequireRole("teacher"),
		// 按主体限流（非 Laravel 原样，属安全加固）：600 次/分钟/教师 ≈ 10 req/s。
		// 阈值刻意远高于正常课堂节奏（批量给分只算 1 次请求），只用于挡住脚本刷与重试风暴；
		// 挂在分组上而非逐路由，是为了避免「按 IP 计数让全校共享额度」并把 GET 也一并兜住。
		middleware.ThrottlePrincipal(600, time.Minute),
	)
	teacher.GET("/dashboard", h.TeacherDashboard)
	teacher.GET("/students", h.TeacherStudents)

	// 教师端：班级信息卡 / 学生 CRUD / 学生导入 / 宠物图鉴 / 按规则加减分
	//（Laravel routes/api.php 第 163-168、181、194、213 行）
	teacher.GET("/class", h.TeacherClassInfo)
	teacher.POST("/students", h.TeacherStudentCreate)
	teacher.POST("/students/import", h.TeacherStudentImport)
	teacher.PUT("/students/:id", h.TeacherStudentUpdate)
	teacher.DELETE("/students/:id", h.TeacherStudentDelete)
	teacher.GET("/pets/:student_id/collection", h.TeacherPetCollection)
	teacher.POST("/scores/give-by-rule/:ruleId", h.TeacherScoreGiveByRule)

	// 积分规则：Laravel routes/api.php 第 185-188 行为 `scores` 前缀组内的 `rules`
	// → 实际路径 /teacher/scores/rules（前端 frontend-vue 与 mini-program 均按此调用）
	teacher.GET("/scores/rules", h.TeacherScoreRules)
	teacher.POST("/scores/rules", h.TeacherCreateScoreRule)
	teacher.PUT("/scores/rules/:id", h.TeacherUpdateScoreRule)
	teacher.DELETE("/scores/rules/:id", h.TeacherDeleteScoreRule)

	teacher.GET("/scores/summary", h.TeacherScoreSummary)
	teacher.GET("/scores/recent", h.TeacherScoreRecent)
	teacher.POST("/scores/give", h.TeacherScoreGive)
	teacher.POST("/scores/batch-give", h.TeacherScoreBatchGive)
	teacher.POST("/scores/:id/undo", h.TeacherScoreUndo)
	teacher.GET("/scores/history/:student_id", h.TeacherScoreHistory)

	teacher.GET("/pets/overview", h.TeacherPetOverview)
	teacher.GET("/pets/:student_id", h.TeacherPetDetail)
	teacher.POST("/pets/:student_id/feed", h.TeacherPetFeed)
	teacher.POST("/pets/:student_id/rename", h.TeacherPetRename)
	teacher.POST("/pets/:student_id/switch", h.TeacherPetSwitch)

	teacher.GET("/leaderboard/total", h.TeacherLeaderboardTotal)
	teacher.GET("/leaderboard/weekly", h.TeacherLeaderboardWeekly)
	teacher.GET("/leaderboard/pet-level", h.TeacherLeaderboardPetLevel)

	// 教师端：班级宠物系列 + 年级战场（PK）
	teacher.POST("/class/switch-series", h.TeacherSwitchSeries)
	teacher.GET("/pk/leaderboard", h.TeacherPKLeaderboard)
	teacher.GET("/pk/my-stats", h.TeacherPKMyStats)
	teacher.POST("/pk/challenge", h.TeacherPKChallenge)

	teacher.GET("/notices", h.TeacherNoticesList)
	teacher.POST("/notices", h.TeacherNoticeCreate)
	teacher.PUT("/notices/:id", h.TeacherNoticeUpdate)
	teacher.PUT("/notices/:id/publish", h.TeacherNoticePublish)
	teacher.PUT("/notices/:id/unpublish", h.TeacherNoticeUnpublish)
	teacher.DELETE("/notices/:id", h.TeacherNoticeDelete)

	teacher.GET("/shop/items", h.TeacherShopItems)
	teacher.POST("/shop/items", h.TeacherShopItemCreate)
	teacher.PUT("/shop/items/:id", h.TeacherShopItemUpdate)
	teacher.DELETE("/shop/items/:id", h.TeacherShopItemDelete)
	teacher.GET("/shop/redemptions", h.TeacherShopRedemptions)
	teacher.POST("/shop/redemptions", h.TeacherShopRedemptionCreate)
	teacher.PUT("/shop/redemptions/:id/approve", h.TeacherShopRedemptionApprove)
	teacher.PUT("/shop/redemptions/:id/reject", h.TeacherShopRedemptionReject)
	teacher.PUT("/shop/redemptions/:id/deliver", h.TeacherShopRedemptionDeliver)

	teacher.GET("/attendance/today", h.TeacherAttendanceToday)
	teacher.POST("/attendance/start", h.TeacherAttendanceStart)
	teacher.PUT("/attendance/:student_id", h.TeacherAttendanceSet)
	teacher.POST("/attendance/:student_id/mark-leave", h.TeacherAttendanceMarkLeave)
	teacher.POST("/attendance/:student_id/mark-absent", h.TeacherAttendanceMarkAbsent)
	teacher.GET("/attendance/summary", h.TeacherAttendanceSummary)

	teacher.GET("/broadcasts", h.TeacherBroadcastsList)
	teacher.POST("/broadcasts", h.TeacherBroadcastSend)
	teacher.GET("/broadcasts/:id", h.TeacherBroadcastGet)

	teacher.GET("/currency/wallets", h.TeacherCurrencyWallets)
	teacher.GET("/currency/exchange-logs", h.TeacherCurrencyLogs)
	teacher.POST("/currency/exchange", h.TeacherCurrencyExchange)
	teacher.POST("/currency/cross-exchange", h.TeacherCurrencyCrossExchange)
	teacher.GET("/exchange-rates", h.TeacherCurrencyRates)
	teacher.POST("/exchange-rates", h.TeacherCurrencyRateCreate)
	teacher.PUT("/exchange-rates/:id", h.TeacherCurrencyRateUpdate)

	admin.GET("/exchange-rates", h.AdminExchangeRates)
	admin.POST("/exchange-rates", h.AdminExchangeRateCreate)
	admin.PUT("/exchange-rates/:id", h.AdminExchangeRateUpdate)

	// 管理端：报表（全校概览 / 按年级 / 按班级 / 班级码登录日志）
	admin.GET("/reports/overview", h.AdminReportsOverview)
	admin.GET("/reports/by-grade", h.AdminReportsByGrade)
	admin.GET("/reports/by-class", h.AdminReportsByClass)
	admin.GET("/display-login-logs", h.AdminDisplayLoginLogs)

	// 管理端：系统运维（诊断 / 状态 / 日志 / 修复）
	admin.GET("/system/diagnose", h.AdminSystemDiagnose)
	admin.GET("/system/status", h.AdminSystemStatus)
	admin.GET("/system/logs", h.AdminSystemLogs)
	admin.POST("/system/repair", h.AdminSystemRepair)

	// 管理端：批量账号（重置密码 / 删除）
	admin.POST("/accounts/batch-reset-password", h.AdminBatchResetPassword)
	admin.POST("/accounts/batch-delete", h.AdminBatchDeleteAccounts)

	// 管理端：教师批量创建 / CSV 导入 / 模板下载
	admin.POST("/teachers/batch-create", h.AdminBatchCreateTeachers)
	admin.POST("/teachers/import", h.AdminImportTeachers)
	admin.GET("/teachers/template-csv", h.AdminTeacherTemplateCsv)

	// 管理端：学生批量导入 / 删除 / 转班
	admin.POST("/students/import", h.AdminImportStudents)
	admin.POST("/students/batch-delete", h.AdminBatchDeleteStudents)
	admin.POST("/students/batch-move", h.AdminBatchMoveStudents)

	// 管理端：班级批量创建 + 班级教师分配
	admin.POST("/classes/batch-create", h.AdminBatchCreateClasses)
	admin.POST("/classes/:id/assign-teacher", h.AdminAssignClassTeacher)
	admin.DELETE("/classes/:id/remove-teacher", h.AdminRemoveClassTeacher)

	// 管理端：学年升级（不可逆；预览 → 确认 → 执行）
	// 学年升级：Laravel 原路径为 /admin/grade-upgrade/*（routes/api.php 第 110-111 行，位于 students 分组之外）。
	admin.GET("/grade-upgrade/preview", h.AdminGradeUpgradePreview)
	admin.POST("/grade-upgrade/execute", h.AdminGradeUpgradeExecute)

	// 管理端：学校 LOGO 上传
	admin.POST("/school/logo", h.AdminUploadSchoolLogo)

	teacher.GET("/timetable", h.TeacherTimetableShow)
	teacher.POST("/timetable", h.TeacherTimetableSave)
	teacher.GET("/timetable/changes", h.TeacherTimetableChanges)
	teacher.GET("/timetable/my-schedule", h.TeacherTimetableMySchedule)
	teacher.GET("/timetable/export-cses", h.TeacherTimetableExportCses)
	// 教师端：报表（积分趋势 / 宠物等级分布 / 学生进度）
	teacher.GET("/reports/score-trend", h.TeacherReportScoreTrend)
	teacher.GET("/reports/pet-distribution", h.TeacherReportPetDistribution)
	teacher.GET("/reports/student-progress", h.TeacherReportStudentProgress)
	// 教师端报表导出（CSV 版：scores / pets / attendance；Laravel routes/api.php 第 242 行）
	teacher.GET("/reports/export/:type", h.TeacherReportExport)

	// 教师端：班级大屏数据与课堂消息（banner/popup/fullscreen → 广播；其余 → 通知）
	teacher.GET("/classroom/display", h.TeacherClassroomDisplay)
	teacher.GET("/classroom/messages", h.TeacherClassroomMessages)
	teacher.POST("/classroom/messages", h.TeacherClassroomMessageSend)

	admin.GET("/timetable/changes", h.AdminTimetableChanges)
	admin.POST("/timetable/changes/:id/approve", h.AdminTimetableApprove)
	admin.POST("/timetable/changes/:id/reject", h.AdminTimetableReject)

	// 管理端：课表批量导入（CSV，dry_run 预览）/ 教师不可用时段 / 冲突检查 / 任课设置
	//（Laravel routes/api.php 第 98-99、145、150-152 行；任课路由沿用既有 :id 以匹配 gin 参数名）
	admin.POST("/timetable/import-csv", h.AdminTimetableImportCsv)
	admin.GET("/timetable/unavailabilities", h.AdminTimetableUnavailabilities)
	admin.POST("/timetable/unavailabilities", h.AdminTimetableSaveUnavailabilities)
	admin.POST("/timetable/check-conflicts", h.AdminTimetableCheckConflicts)
	// 管理端：规则自动排课（单班纯计算预览；全校智能排课 commit=true 落库）
	admin.POST("/timetable/generate", h.AdminTimetableGenerate)
	admin.POST("/timetable/generate-school", h.AdminTimetableGenerateSchool)
	admin.GET("/classes/:id/teacher-assignments", h.AdminClassTeacherAssignmentsList)
	admin.PUT("/classes/:id/teacher-assignments", h.AdminClassTeacherAssignmentsSave)
	admin.POST("/classes/:id/teacher-assignments", h.AdminClassTeacherAssignmentsSave)
	admin.GET("/classes/:id/timetable", h.AdminClassTimetableShow)
	// Laravel 为 match(['put','post'], '{classId}/timetable') → 两个方法都注册
	// （教师端 assignments 同款双注册；前端班级课表页用 POST）。
	admin.PUT("/classes/:id/timetable", h.AdminClassTimetableSave)
	admin.POST("/classes/:id/timetable", h.AdminClassTimetableSave)
	// 管理端：课表单班 / 全校网格导出（CSV 版；Laravel routes/api.php 第 97、146 行）
	admin.GET("/classes/:id/timetable/export-excel", h.AdminClassTimetableExportExcel)
	admin.GET("/timetable/export-excel", h.AdminTimetableExportExcel)

	// 管理端：全校积分规则与商品（学校级共享）
	admin.GET("/score-rules", h.AdminScoreRules)
	admin.POST("/score-rules", h.AdminScoreRuleCreate)
	admin.PUT("/score-rules/:id", h.AdminScoreRuleUpdate)
	admin.DELETE("/score-rules/:id", h.AdminScoreRuleDelete)

	admin.GET("/shop-items", h.AdminShopItems)
	admin.POST("/shop-items", h.AdminShopItemCreate)
	admin.PUT("/shop-items/:id", h.AdminShopItemUpdate)

	// 管理端：AI 中心（设置 / 开关 / 用量 / 模型列表 / 连通性测试）
	admin.GET("/ai/settings", h.AdminAISettings)
	admin.PUT("/ai/settings", h.AdminAISaveSettings)
	admin.POST("/ai/toggle", h.AdminAIToggle)
	admin.GET("/ai/usage", h.AdminAIUsage)
	admin.POST("/ai/fetch-models", h.AdminAIFetchModels)
	admin.POST("/ai/test", h.AdminAITest)
	// 管理端：AI 供应商官方账单直查（余额 / 用量；Laravel routes/api.php 第 131 行）
	admin.POST("/ai/provider-official", h.AdminAIProviderOfficial)

	// 管理端：第三方通讯录拉取与导入（Laravel routes/api.php 第 133-136 行）
	admin.GET("/wechat-work/contacts", h.AdminWechatWorkContacts)
	admin.POST("/wechat-work/import", h.AdminWechatWorkImport)
	admin.GET("/third-party/contacts", h.AdminThirdPartyContacts)
	admin.POST("/third-party/import", h.AdminThirdPartyImport)
	admin.DELETE("/shop-items/:id", h.AdminShopItemDelete)

	// 教师端：我的班级 / 切换班级 / 模式
	teacher.GET("/my-classes", h.TeacherMyClasses)
	teacher.POST("/switch-class", h.TeacherSwitchClass)
	teacher.GET("/mode", h.TeacherGetMode)
	teacher.POST("/mode", h.TeacherSetMode)

	// 教师端：AI 助教（配置态 / 对话 / 预设命令 / 用量）
	teacher.GET("/ai/config", h.TeacherAIConfig)
	teacher.POST("/ai/chat", h.TeacherAIChat)
	teacher.GET("/ai/commands", h.TeacherAICommands)
	teacher.GET("/ai/usage", h.TeacherAIUsage)
	// 班级大屏 / 教室端（班级码 token，disp_ / class_ 前缀；不加 JWT 中间件）
	// display/login 另有 throttle:10,1（Laravel routes/api.php 第 295 行）；它与
	// DisplayService 内部「同班级码连续失败 5 次锁 15 分钟」是两套并存的独立机制。
	displayGroup := api.Group("/display")
	displayGroup.POST("/login", middleware.Throttle(10, time.Minute), h.DisplayLogin)

	displayAuth := api.Group("/display",
		middleware.DisplayAuth(db),
		// 按主体限流（同上）：1200 次/分钟/班级 ≈ 20 req/s。教室端可能多台设备共用同一个
		// 班级码 token，且本组含 `/poll` 轮询与 `/sse` 长连接，故阈值比教师端更宽松。
		middleware.ThrottlePrincipal(1200, time.Minute),
	)
	displayAuth.GET("/initial-data", h.DisplayInitialData)
	displayAuth.GET("/timetable", h.DisplayTimetable)
	displayAuth.GET("/export-cses", h.DisplayExportCses)
	displayAuth.GET("/leaderboard", h.DisplayLeaderboard)
	displayAuth.GET("/class-settings", h.DisplayClassSettings)
	displayAuth.GET("/dashboard", h.DisplayDashboard)
	displayAuth.GET("/students", h.DisplayStudents)
	displayAuth.GET("/scores/rules", h.DisplayScoreRules)
	displayAuth.GET("/pets/overview", h.DisplayPetsOverview)
	// 大屏端 AI（开关检查 + 对话，均走班级码 token）
	displayAuth.GET("/ai/settings", h.DisplayAISettings)
	displayAuth.POST("/ai/chat", h.DisplayAIChat)

	// 事件链路：SSE 长连接 + 轮询降级
	displayAuth.GET("/sse", h.DisplaySSE)
	displayAuth.GET("/poll", h.DisplayPoll)

	// 大屏加减分
	displayAuth.POST("/quick-score", h.DisplayQuickScore)
	displayAuth.POST("/scores/give", h.DisplayGiveScore)
	displayAuth.POST("/scores/batch-give", h.DisplayBatchGiveScore)

	// 大屏写操作 / 查询（班级码 token）
	displayAuth.GET("/shop-items", h.DisplayShopItems)
	displayAuth.POST("/redeem", h.DisplayRedeem)
	displayAuth.POST("/transfer", h.DisplayTransfer)
	displayAuth.POST("/switch-series", h.DisplaySwitchSeries)
	displayAuth.POST("/pets/switch", h.DisplaySwitchPet)
	displayAuth.GET("/pk/leaderboard", h.DisplayPKLeaderboard)

	// 教师端班级码
	teacher.GET("/display-code", h.TeacherGetDisplayCode)
	teacher.POST("/display-code/refresh", h.TeacherRefreshDisplayCode)

	// 管理端班级码。
	// Gin 的路径参数名在同一层级必须一致，故沿用既有路由的 :id（Laravel 为 {classId}；
	// 路径形状 /admin/classes/:id/display-code 与原实现一致）。
	admin.POST("/classes/reset-display-codes", h.AdminResetDisplayCodes)
	admin.GET("/classes/:id/display-code", h.AdminGetDisplayCode)
	admin.POST("/classes/:id/display-code/refresh", h.AdminRefreshDisplayCode)

	// 公开字典与 webhook（Laravel routes/api.php 第 326、331-332 行）：
	//   GET  common/score-categories —— 积分分类字典（无需登录；响应裸 {data:[…]}，无 message）
	//   GET  wechat-work/callback    —— 企微回调验签，成功回显 echostr 纯文本
	//   POST wechat-work/callback    —— 企微审批事件接收，恒返回 {"errcode":0,"errmsg":"ok"}
	api.GET("/common/score-categories", h.ScoreCategories)
	api.GET("/wechat-work/callback", h.WechatWorkCallbackVerify)
	api.POST("/wechat-work/callback", h.WechatWorkCallbackReceive)

	// SPA 静态资源与前端路由兜底（等价旧 Laravel routes/web.php 的 SPA 兜底行为）：
	//   - /api/* 未命中 → 统一 JSON 404；
	//   - /storage/app/uploads/* → UPLOAD_DIR 下的上传文件（学校 LOGO 等）；
	//   - 其余 GET/HEAD → PUBLIC_DIR 静态文件；文件不存在回落 index.html
	//     （Vue Router history 模式刷新不 404）；
	//   - PUBLIC_DIR 未配置或不存在（纯 API 部署）→ 一律 JSON 404，不报错。
	r.NoRoute(spaFallback(cfg))

	return r
}
