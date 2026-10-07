// 班级大屏（教室端）服务：班级码鉴权、只读接口与班级码管理。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\DisplayController（getDisplayCode /
// refreshDisplayCode / login / validateToken / initialData / timetable / exportCses /
// classSettings / classroomDashboard / classroomStudents / scoreRules /
// classroomPetsOverview / quickLeaderboard）与 SchoolAdminController 的班级码三个动作。
//
// 有意差异（本批统一约束）：
//  1. Token 存储：Laravel 用 Cache（display:token:<t> / class_token:<t>），Go 无 Redis/Cache
//     → 改为数据库表 display_tokens，过期即视为无效；disp_ 与 class_ 前缀共用该表。
//  2. 登录失败计数：Laravel 用 Cache 计数（15 分钟 TTL），Go 改为进程内并发安全 map
//     —— 仅单实例有效；多副本部署时各副本计数独立，需换成共享存储（Redis / DB 表）。
//     该机制与 `display/login` 上的 `throttle:10,1`（middleware.Throttle，同为本进程内实现）**并存**，
//     两者都会返回 429 —— 与 Laravel 的两套机制并存一致。
//  3. 事件总线（sse / poll / 各发布点）：存储由 Cache 改为数据库表 display_events，
//     详见 services/display_events.go 顶部说明。
//  4. 大屏加减分复用 ScoreService（含审计日志 + score_update 事件），与 Laravel
//     classroomGiveScore / BatchGiveScore 的手写 Score::create 路径有差异，见各方法注释。
//  5. 免费自选机会（pet_free_pick）：Laravel 存 Cache；Go 端改为 pet_free_picks 表
//     （models/pet_free_pick.go），各写路径已改用该表；classroomStudents 的 free_pick 字段
//     属只读接口批的遗留，仍未接该表（保持 false），待后续批次对齐。
//  6. 本批（写操作批）已补齐 quickShopItems 之外的写/查接口：quickRedeem / quickTransfer /
//     classroomSwitchSeries / classroomSwitchPet / classroomPKLeaderboard / 班级码登录。
//     它们分别在 services/{shop,pet_series,pet,pk,class_login}.go 与 display_writes.go 中实现。
//     仍未移植：reports/*、system/* 等（不提供空实现）。大屏 AI 已在 ai_assistant.go 实现。
package services

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// 班级码登录暴力破解防护阈值（同 Laravel DisplayController）。
const (
	// DisplayLoginMaxAttempts 同一班级码连续失败上限。
	DisplayLoginMaxAttempts = 5
	// DisplayLoginLockMinutes 触发上限后的锁定时长（分钟）。
	DisplayLoginLockMinutes = 15
	// displayRecentScoreHours 大屏最近积分的时间窗口（小时）。
	displayRecentScoreHours = 4
	// displayRecentScoreLimit 大屏最近积分条数上限。
	displayRecentScoreLimit = 20
	// displayBroadcastLimit 大屏生效广播条数上限。
	displayBroadcastLimit = 5
	// displayPeakPetLevel 班级总览「尖峰」宠物等级门槛（同 Laravel classroomDashboard）。
	displayPeakPetLevel = 10
)

// displayLoginLocks 班级码登录失败计数（键 = 大写后的班级码）。
//
// ⚠️ 仅单实例有效：本计数存在进程内存里，多副本部署时各副本各自计数，
// 需改为共享存储（Redis / 数据库表）才能保持「同码 5 次失败锁定 15 分钟」的语义。
var displayLoginLocks = struct {
	mu    sync.Mutex
	locks map[string]displayLoginLock
}{locks: map[string]displayLoginLock{}}

type displayLoginLock struct {
	attempts  int
	expiresAt time.Time
}

// displayLoginAttempts 当前失败次数（超过锁定窗口视为 0）。
func displayLoginAttempts(code string, now time.Time) int {
	displayLoginLocks.mu.Lock()
	defer displayLoginLocks.mu.Unlock()

	l, ok := displayLoginLocks.locks[code]
	if !ok || !now.Before(l.expiresAt) {
		return 0
	}
	return l.attempts
}

// displayLoginFail 记一次失败（TTL 从本次起算 15 分钟，同 Laravel 每失败即续期）。
func displayLoginFail(code string, now time.Time) {
	displayLoginLocks.mu.Lock()
	defer displayLoginLocks.mu.Unlock()

	attempts := 0
	if l, ok := displayLoginLocks.locks[code]; ok && now.Before(l.expiresAt) {
		attempts = l.attempts
	}
	displayLoginLocks.locks[code] = displayLoginLock{
		attempts:  attempts + 1,
		expiresAt: now.Add(DisplayLoginLockMinutes * time.Minute),
	}
}

// displayLoginReset 成功登录后清空计数。
func displayLoginReset(code string) {
	displayLoginLocks.mu.Lock()
	defer displayLoginLocks.mu.Unlock()
	delete(displayLoginLocks.locks, code)
}

// DisplayService 班级大屏服务。
type DisplayService struct {
	db        *gorm.DB
	scope     *Scope
	codes     *DisplayCodeService
	timetable *TimetableService
	events    *DisplayEvents
	scores    *ScoreService
}

// NewDisplayService 创建班级大屏服务（自带事件发布器与积分服务：大屏加减分复用 ScoreService）。
func NewDisplayService(db *gorm.DB) *DisplayService {
	return &DisplayService{
		db:        db,
		scope:     NewScope(db),
		codes:     NewDisplayCodeService(db),
		timetable: NewTimetableService(db),
		events:    NewDisplayEvents(db),
		scores:    NewScoreService(db),
	}
}

// ErrDisplayToken 构造大屏 token 失效错误（401，文案同 Laravel）。
func ErrDisplayToken() *AppError {
	return NewAppError(http.StatusUnauthorized, "Token 无效或已过期")
}

// ============================================================
// 显示端 Token
// ============================================================

// ValidateToken 校验显示端 token 并返回所属班级 ID（先 query 后 Bearer 的取值顺序在中间件完成）。
func (d *DisplayService) ValidateToken(token string) (uint, error) {
	classID, ok := models.ResolveDisplayToken(d.db, token)
	if !ok {
		return 0, ErrDisplayToken()
	}
	return classID, nil
}

// resolveCodeToClassID 解析班级码 → 班级 ID。
// Laravel 先查缓存再查库；Go 无缓存，直接查 display_code，未命中时按确定性计算码全库匹配。
func (d *DisplayService) resolveCodeToClassID(code string) (uint, error) {
	if code == "" {
		return 0, nil
	}

	var room models.ClassRoom
	err := d.db.Where("display_code = ?", code).First(&room).Error
	if err == nil {
		return room.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	// 数据库未刷新时按确定性计算码匹配（同 Laravel：遍历全部班级调 generate()）。
	var rooms []models.ClassRoom
	if err := d.db.Order("id ASC").Find(&rooms).Error; err != nil {
		return 0, err
	}
	for i := range rooms {
		if d.codes.Generate(&rooms[i], nil) == code {
			return rooms[i].ID, nil
		}
	}
	return 0, nil
}

func (d *DisplayService) issueToken(class *models.ClassRoom, ip string) (string, error) {
	value, err := models.NewDisplayTokenValue()
	if err != nil {
		return "", err
	}

	// 过期行惰性清理：Laravel 由 Cache TTL 自动回收，Go 端借登录时机顺手清理。
	// ⚠️ 写入与比较统一 UTC 域（SQLite 对 time 列做文本比较，见 services/auth.go 的时区纪律）。
	_ = d.db.Where("expires_at < ?", time.Now().UTC()).Delete(&models.DisplayToken{}).Error

	row := models.DisplayToken{
		Token:     value,
		ClassID:   class.ID,
		ClassName: class.Name,
		Grade:     class.Grade,
		IPAddress: ip,
		ExpiresAt: time.Now().UTC().Add(models.DisplayTokenTTL * time.Second),
	}
	if err := d.db.Create(&row).Error; err != nil {
		return "", err
	}
	return value, nil
}

// ============================================================
// 显示端 — 班级码登录
// ============================================================

// DisplayClassInfo 登录返回的班级信息（字段名同 Laravel）。
type DisplayClassInfo struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Grade        string `json:"grade"`
	StudentCount int64  `json:"student_count"`
}

// DisplayLoginResult 班级码登录结果。
type DisplayLoginResult struct {
	Token     string           `json:"token"`
	ExpiresIn int              `json:"expires_in"`
	ClassInfo DisplayClassInfo `json:"class_info"`
}

// Login 班级码登录：校验失败计数、解析班级码、签发 disp_ token 并记录登录日志。
func (d *DisplayService) Login(code, ip, userAgent string) (*DisplayLoginResult, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	now := util.Now()

	// 暴力破解防护（同码 5 次失败 → 429，锁定 15 分钟）。
	if displayLoginAttempts(code, now) >= DisplayLoginMaxAttempts {
		return nil, NewAppError(http.StatusTooManyRequests,
			fmt.Sprintf("尝试次数过多，请 %d 分钟后再试", DisplayLoginLockMinutes))
	}

	classID, err := d.resolveCodeToClassID(code)
	if err != nil {
		return nil, err
	}
	if classID == 0 {
		displayLoginFail(code, now)
		return nil, ErrNotFound("班级码无效，请检查后重试")
	}

	// 命中后立即重置计数（同 Laravel Cache::forget）。
	displayLoginReset(code)

	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}

	token, err := d.issueToken(&class, ip)
	if err != nil {
		return nil, err
	}

	// 登录日志：写库失败不影响登录（Laravel 同样只记 warning）。
	_ = d.db.Create(&models.DisplayLoginLog{
		ClassID:   class.ID,
		ClassCode: code,
		IPAddress: ip,
		UserAgent: userAgent,
	}).Error

	studentCount, err := d.activeStudentCount(class.ID)
	if err != nil {
		return nil, err
	}

	return &DisplayLoginResult{
		Token:     token,
		ExpiresIn: models.DisplayTokenTTL,
		ClassInfo: DisplayClassInfo{
			ID:           class.ID,
			Name:         class.Name,
			Grade:        class.Grade,
			StudentCount: studentCount,
		},
	}, nil
}

// ============================================================
// 教师端 / 管理员端 — 班级码管理
// ============================================================

// DisplayCodeView 班级码信息（字段名同 Laravel，按动作裁剪可选字段）。
type DisplayCodeView struct {
	Code         string  `json:"code"`
	ClassName    string  `json:"class_name"`
	UpdatedAt    *string `json:"updated_at"`
	StudentCount *int64  `json:"student_count,omitempty"`
	Message      string  `json:"message,omitempty"`
}

// TeacherDisplayCode 获取（或创建）当前班级的大屏码。
func (d *DisplayService) TeacherDisplayCode(u *models.User, classID uint) (*DisplayCodeView, error) {
	if classID == 0 {
		classID = u.SettingUint(ActiveClassSettingKey)
	}
	if classID == 0 {
		return nil, ErrBadRequest("请先选择班级")
	}

	class, err := d.teacherClass(u, classID)
	if err != nil {
		return nil, err
	}

	// 班级码为确定性（LS + 年级 + 班号），实时计算返回，不依赖数据库存储。
	code := d.codes.Generate(class, nil)

	studentCount, err := d.activeStudentCount(class.ID)
	if err != nil {
		return nil, err
	}

	return &DisplayCodeView{
		Code:         code,
		ClassName:    class.Name,
		UpdatedAt:    isoTimePtrOf(class.DisplayCodeUpdatedAt),
		StudentCount: &studentCount,
	}, nil
}

// TeacherRefreshDisplayCode 刷新班级大屏码（确定性码：刷新仅更新时间戳）。
func (d *DisplayService) TeacherRefreshDisplayCode(u *models.User, classID uint) (*DisplayCodeView, error) {
	if classID == 0 {
		classID = u.SettingUint(ActiveClassSettingKey)
	}
	if classID == 0 {
		return nil, ErrBadRequest("请先选择班级")
	}

	class, err := d.teacherClass(u, classID)
	if err != nil {
		return nil, err
	}

	code := d.codes.Generate(class, nil)
	now := util.Now()
	class.DisplayCode = code
	class.DisplayCodeUpdatedAt = &now
	if err := d.db.Save(class).Error; err != nil {
		return nil, err
	}

	// 发送 refresh 事件通知当前大屏刷新（载荷同 Laravel：{old_code, new_code}）。
	codeCopy := code
	publishDisplayEvent(d.events, class.ID, DisplayEventRefresh, DisplayRefreshData{
		OldCode: &codeCopy,
		NewCode: code,
	})

	return &DisplayCodeView{
		Code:      code,
		ClassName: class.Name,
		UpdatedAt: isoTimePtrOf(class.DisplayCodeUpdatedAt),
		Message:   "班级码已刷新，旧码已失效",
	}, nil
}

// AdminDisplayCode 管理员获取某班的大屏码（限本校班级）。
func (d *DisplayService) AdminDisplayCode(u *models.User, classID uint) (*DisplayCodeView, error) {
	class, err := d.adminClass(u, classID)
	if err != nil {
		return nil, err
	}

	code := d.codes.Generate(class, nil)
	studentCount, err := d.activeStudentCount(class.ID)
	if err != nil {
		return nil, err
	}

	return &DisplayCodeView{
		Code:         code,
		ClassName:    class.Name,
		UpdatedAt:    isoTimePtrOf(class.DisplayCodeUpdatedAt),
		StudentCount: &studentCount,
	}, nil
}

// AdminRefreshDisplayCode 管理员刷新某班的大屏码。
func (d *DisplayService) AdminRefreshDisplayCode(u *models.User, classID uint) (*DisplayCodeView, error) {
	class, err := d.adminClass(u, classID)
	if err != nil {
		return nil, err
	}

	code := d.codes.Generate(class, nil)
	now := util.Now()
	class.DisplayCode = code
	class.DisplayCodeUpdatedAt = &now
	if err := d.db.Save(class).Error; err != nil {
		return nil, err
	}

	// 通知大屏刷新（管理员端载荷只有 {new_code}，同 Laravel SchoolAdminController）。
	publishDisplayEvent(d.events, class.ID, DisplayEventRefresh, DisplayRefreshData{NewCode: code})
	return &DisplayCodeView{
		Code:      code,
		ClassName: class.Name,
		UpdatedAt: isoTimePtrOf(class.DisplayCodeUpdatedAt),
	}, nil
}

// AdminResetDisplayCodes 批量重置本校全部班级码（自定义统一字母前缀，如 LS / BJ / LE）。
// 前缀持久化到学校设置，后续新建班级沿用同一前缀；空前缀按学校当前设置（缺失即 LS）处理。
//
// 说明：Laravel 在此额外返回一句中文 message（成功/跳过/冲突计数）；Go 端统一走
// {"data":...,"message":"ok"} 信封，计数信息保留在 data 内。
func (d *DisplayService) AdminResetDisplayCodes(u *models.User, rawPrefix string) (*DisplayCodeRegenerateResult, error) {
	if rawPrefix != "" && !IsValidDisplayPrefixInput(rawPrefix) {
		return nil, ErrUnprocessable("字母前缀需为 2-4 个英文字母（如 LS / BJ / LE）")
	}
	prefix := strings.ToUpper(rawPrefix)

	var school models.School
	if err := d.db.First(&school, u.SchoolID).Error; err != nil {
		return nil, ErrNotFound("学校不存在")
	}
	settings, err := school.WithSetting("display_code_prefix", prefix)
	if err != nil {
		return nil, err
	}
	if err := d.db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", settings).Error; err != nil {
		return nil, err
	}

	return d.codes.RegenerateAll(&u.SchoolID, prefix)
}

// teacherClass 取教师管辖范围内的班级（越权/不存在 → 404）。
// 差异说明：Laravel 两个班级码动作直接用 ClassRoom::findOrFail（不校验管辖范围），
// Go 端沿袭本项目其他接口的做法收敛到教师可管辖班级。
func (d *DisplayService) teacherClass(u *models.User, classID uint) (*models.ClassRoom, error) {
	if err := d.scope.ClassInScope(u, classID); err != nil {
		return nil, ErrNotFound("班级不存在或不在管辖范围")
	}
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}
	return &class, nil
}

// adminClass 取本校班级（跨校/不存在 → 404）。
func (d *DisplayService) adminClass(u *models.User, classID uint) (*models.ClassRoom, error) {
	var class models.ClassRoom
	err := d.db.Where("id = ? AND school_id = ?", classID, u.SchoolID).First(&class).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("班级不存在")
	}
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// ============================================================
// 显示端 — 初始全量数据
// ============================================================

// DisplayPetItem 大屏宠物卡片（字段名同 Laravel formatPets）。
type DisplayPetItem struct {
	StudentID   uint    `json:"student_id"`
	StudentNo   string  `json:"student_no"`
	StudentName string  `json:"student_name"`
	TotalScore  int     `json:"total_score"`
	HasPet      bool    `json:"has_pet"`
	PetName     *string `json:"pet_name"`
	PetSpecies  *string `json:"pet_species"`
	Level       int     `json:"level"`
	Experience  int     `json:"experience"`
	Mood        int     `json:"mood"`
	Emoji       string  `json:"emoji"`
	StageName   string  `json:"stage_name"`
	StageTitle  string  `json:"stage_title"`
	ExpMax      int     `json:"exp_max"`
	Color       string  `json:"color"`
	Image       *string `json:"image"`
}

// DisplayRecentScore 大屏最近积分条目。
type DisplayRecentScore struct {
	StudentName *string `json:"student_name"`
	Amount      int     `json:"amount"`
	Reason      string  `json:"reason"`
	Time        string  `json:"time"`
}

// DisplayBroadcast 大屏生效中的广播。
type DisplayBroadcast struct {
	ID             uint   `json:"id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	DisplaySeconds int    `json:"display_seconds"`
	VoiceEnabled   bool   `json:"voice_enabled"`
	CreatedAt      string `json:"created_at"`
}

// DisplayInitialData 大屏初始全量数据。
type DisplayInitialData struct {
	ClassName    string               `json:"class_name"`
	Grade        string               `json:"grade"`
	StudentCount int                  `json:"student_count"`
	Pets         []DisplayPetItem     `json:"pets"`
	RecentScores []DisplayRecentScore `json:"recent_scores"`
	Broadcasts   []DisplayBroadcast   `json:"broadcasts"`
	ServerTime   string               `json:"server_time"`
}

// InitialData 大屏初始全量：班级信息 + 学生（含宠物卡）+ 生效广播 + 最近 4 小时积分。

func (d *DisplayService) InitialData(classID uint) (*DisplayInitialData, error) {
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}

	students, err := d.activeStudents(classID)
	if err != nil {
		return nil, err
	}
	sortStudentsByNo(students)

	pets, err := d.ensureAndFormatPets(students)
	if err != nil {
		return nil, err
	}

	now := util.Now()

	var scores []models.Score
	if err := d.db.Where("class_id = ? AND created_at >= ?", classID, now.Add(-displayRecentScoreHours*time.Hour)).
		Order("created_at DESC, id DESC").Limit(displayRecentScoreLimit).Find(&scores).Error; err != nil {
		return nil, err
	}
	names, err := d.studentNamesOf(scores)
	if err != nil {
		return nil, err
	}
	recentScores := make([]DisplayRecentScore, 0, len(scores))
	for _, sc := range scores {
		item := DisplayRecentScore{
			Amount: sc.Amount,
			Reason: sc.Reason,
			Time:   relativeTime(sc.CreatedAt, now),
		}
		if name, ok := names[sc.StudentID]; ok {
			n := name
			item.StudentName = &n
		}
		recentScores = append(recentScores, item)
	}

	var broadcasts []models.Broadcast
	if err := d.db.Where("class_id = ? AND status IN ?", classID, []string{"pending", "sent"}).
		Order("created_at DESC, id DESC").Limit(displayBroadcastLimit).Find(&broadcasts).Error; err != nil {
		return nil, err
	}
	broadcastItems := make([]DisplayBroadcast, 0, len(broadcasts))
	for _, b := range broadcasts {
		broadcastItems = append(broadcastItems, DisplayBroadcast{
			ID:             b.ID,
			Content:        b.Content,
			Type:           b.Type,
			DisplaySeconds: b.DisplaySeconds,
			VoiceEnabled:   b.VoiceEnabled,
			CreatedAt:      relativeTime(b.CreatedAt, now),
		})
	}

	return &DisplayInitialData{
		ClassName:    class.Name,
		Grade:        class.Grade,
		StudentCount: len(students),
		Pets:         pets,
		RecentScores: recentScores,
		Broadcasts:   broadcastItems,
		ServerTime:   now.Format(time.RFC3339),
	}, nil
}

// ensureAndFormatPets 为学生补齐宠物（无宠物者随机分配本系列物种）并格式化为大屏卡片。
func (d *DisplayService) ensureAndFormatPets(students []models.Student) ([]DisplayPetItem, error) {
	petByStudent := map[uint]models.Pet{}
	if len(students) > 0 {
		ids := make([]uint, 0, len(students))
		for _, s := range students {
			ids = append(ids, s.ID)
		}
		var pets []models.Pet
		if err := d.db.Where("student_id IN ?", ids).Find(&pets).Error; err != nil {
			return nil, err
		}
		for _, pet := range pets {
			petByStudent[pet.StudentID] = pet
		}
	}

	pool := models.SpeciesPoolForSeries("myth")
	items := make([]DisplayPetItem, 0, len(students))
	for _, s := range students {
		pet, ok := petByStudent[s.ID]
		if !ok {
			// 没有宠物的学生自动分配一只（species 体系；旧 type 体系已废弃）。
			species := "zhulong"
			if len(pool) > 0 {
				species = pool[rand.IntN(len(pool))]
			}
			pet = models.Pet{
				StudentID:  s.ID,
				ClassID:    s.ClassID,
				Name:       s.Name + "的萌宠",
				Species:    species,
				Level:      1,
				Experience: 0,
				Mood:       80,
			}
			if err := d.db.Create(&pet).Error; err != nil {
				return nil, err
			}
			petByStudent[s.ID] = pet
		}

		stage := pet.CurrentStage()
		item := DisplayPetItem{
			StudentID:   s.ID,
			StudentNo:   s.StudentNo,
			StudentName: s.Name,
			TotalScore:  s.TotalScore,
			HasPet:      true,
			PetName:     strPtr(pet.Name),
			PetSpecies:  strPtr(pet.Species),
			Level:       pet.Level,
			Experience:  pet.Experience,
			Mood:        pet.Mood,
			Emoji:       stage.Emoji,
			StageName:   stage.Name,
			StageTitle:  stage.Title,
			ExpMax:      stage.ExpMax,
			Color:       stage.Color,
			// Laravel 的 currentStage() 不返回 image，故恒为 null。
			Image: nil,
		}
		items = append(items, item)
	}
	return items, nil
}

// ============================================================
// 显示端 — 课表
// ============================================================

// DisplayTimetable 大屏课表（顶层展开 forDisplay 的结果，同 Laravel）。
type DisplayTimetable struct {
	TodayWeekday int                    `json:"today_weekday"`
	ClassName    string                 `json:"class_name"`
	Subjects     []TimetableSubjectView `json:"subjects"`
	Periods      []TimetablePeriodView  `json:"periods"`
	Entries      []TimetableEntryView   `json:"entries"`
}

// Timetable 大屏课表：今天星期几 + 班级名 + 科目/节次/排课。
func (d *DisplayService) Timetable(classID uint) (*DisplayTimetable, error) {
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}

	data, err := d.timetable.ForDisplay(class.ID, class.SchoolID)
	if err != nil {
		return nil, err
	}

	return &DisplayTimetable{
		// Laravel now()->dayOfWeekIso：周一=1 … 周日=7。
		TodayWeekday: isoWeekday(util.Now()),
		ClassName:    class.Name,
		Subjects:     data.Subjects,
		Periods:      data.Periods,
		Entries:      data.Entries,
	}, nil
}

// ExportCses 大屏免登录导出 CSES YAML（响应头在 handler 设置）。
func (d *DisplayService) ExportCses(classID uint) (string, *models.ClassRoom, error) {
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return "", nil, ErrNotFound("班级不存在")
	}

	yaml, err := d.timetable.ToCses(class.ID, class.SchoolID)
	if err != nil {
		return "", nil, err
	}
	return yaml, &class, nil
}

// ============================================================
// 教室端 — 只读接口
// ============================================================

// DisplayClassSettings 教室端班级设置。
type DisplayClassSettings struct {
	PetSeries any `json:"pet_series"`
}

// ClassSettings 教室端班级设置（班级不存在时只返回 pet_series = null，同 Laravel）。
func (d *DisplayService) ClassSettings(classID uint) (*DisplayClassSettings, error) {
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &DisplayClassSettings{PetSeries: nil}, nil
		}
		return nil, err
	}
	return &DisplayClassSettings{PetSeries: classPetSeries(&class)}, nil
}

// DisplayTopStudent 大屏总览榜单条目。
type DisplayTopStudent struct {
	Name       string `json:"name"`
	StudentNo  string `json:"student_no"`
	Score      int    `json:"score"`
	PetName    string `json:"pet_name"`
	PetSpecies string `json:"pet_species"`
	PetLevel   int    `json:"pet_level"`
}

// DisplayStarStudent 今日之星。
type DisplayStarStudent struct {
	Name       string `json:"name"`
	StudentNo  string `json:"student_no"`
	PetName    string `json:"pet_name"`
	PetSpecies string `json:"pet_species"`
	PetLevel   int    `json:"pet_level"`
	Score      int    `json:"score"`
}

// DisplayNewsItem 班级动态条目。
type DisplayNewsItem struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
}

// DisplayDashboard 教室端班级总览。
type DisplayDashboard struct {
	ClassName    string              `json:"class_name"`
	Grade        string              `json:"grade"`
	StudentCount int                 `json:"student_count"`
	TotalScore   int                 `json:"total_score"`
	AvgPetLevel  float64             `json:"avg_pet_level"`
	PeakCount    int                 `json:"peak_count"`
	WeeklyScore  int                 `json:"weekly_score"`
	StarStudent  *DisplayStarStudent `json:"star_student"`
	Top5         []DisplayTopStudent `json:"top5"`
	RecentNews   []DisplayNewsItem   `json:"recent_news"`
}

// ClassroomDashboard 教室端班级总览（积分/等级/尖峰/本周/榜单/动态）。
func (d *DisplayService) ClassroomDashboard(classID uint) (*DisplayDashboard, error) {
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}

	students, err := d.activeStudents(classID)
	if err != nil {
		return nil, err
	}
	sortStudentsByNo(students)

	// 宠物按学生一次性取回（Laravel 用 with('pet')）。
	petByStudent, err := d.petsByStudent(students)
	if err != nil {
		return nil, err
	}

	totalScore := 0
	levelSum := 0
	peakCount := 0
	studentIDs := make([]uint, 0, len(students))
	for _, s := range students {
		totalScore += s.TotalScore
		studentIDs = append(studentIDs, s.ID)
		if pet, ok := petByStudent[s.ID]; ok {
			levelSum += pet.Level
			if pet.Level >= displayPeakPetLevel {
				peakCount++
			}
		}
	}
	count := len(students)
	avgLevel := 0.0
	if count > 0 {
		avgLevel = math.Round(float64(levelSum)/float64(count)*10) / 10
	}

	sorted := make([]models.Student, len(students))
	copy(sorted, students)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].TotalScore > sorted[j].TotalScore })

	top5 := make([]DisplayTopStudent, 0, 5)
	for i, s := range sorted {
		if i >= 5 {
			break
		}
		pet := petByStudent[s.ID]
		top5 = append(top5, DisplayTopStudent{
			Name:       s.Name,
			StudentNo:  s.StudentNo,
			Score:      s.TotalScore,
			PetName:    pet.Name,
			PetSpecies: pet.Species,
			PetLevel:   pet.Level,
		})
	}

	var star *DisplayStarStudent
	if len(sorted) > 0 {
		s := sorted[0]
		pet := petByStudent[s.ID]
		star = &DisplayStarStudent{
			Name:       s.Name,
			StudentNo:  s.StudentNo,
			PetName:    pet.Name,
			PetSpecies: pet.Species,
			PetLevel:   pet.Level,
			Score:      s.TotalScore,
		}
	}

	weeklyScore, err := d.scoreSum(studentIDs, util.StartOfWeek(util.Now()))
	if err != nil {
		return nil, err
	}

	recentNews, err := d.recentNews(studentIDs)
	if err != nil {
		return nil, err
	}

	return &DisplayDashboard{
		ClassName:    class.Name,
		Grade:        class.Grade,
		StudentCount: count,
		TotalScore:   totalScore,
		AvgPetLevel:  avgLevel,
		PeakCount:    peakCount,
		WeeklyScore:  weeklyScore,
		StarStudent:  star,
		Top5:         top5,
		RecentNews:   recentNews,
	}, nil
}

// DisplayStudentItem 教室端学生条目（含宠物）。
type DisplayStudentItem struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	StudentNo  string `json:"student_no"`
	TotalScore int    `json:"total_score"`
	PetName    string `json:"pet_name"`
	PetSpecies string `json:"pet_species"`
	PetLevel   int    `json:"pet_level"`
	PetEmoji   string `json:"pet_emoji"`
	FreePick   bool   `json:"free_pick"`
}

// ClassroomStudents 教室端学生列表（含宠物信息）。
// free_pick 在 Laravel 是「Cache::has("pet_free_pick:<id>")」，Go 端无缓存 → 恒为 false。
func (d *DisplayService) ClassroomStudents(classID uint) ([]DisplayStudentItem, error) {
	students, err := d.activeStudents(classID)
	if err != nil {
		return nil, err
	}
	sortStudentsByNo(students)

	petByStudent, err := d.petsByStudent(students)
	if err != nil {
		return nil, err
	}

	items := make([]DisplayStudentItem, 0, len(students))
	for _, s := range students {
		item := DisplayStudentItem{
			ID:         s.ID,
			Name:       s.Name,
			StudentNo:  s.StudentNo,
			TotalScore: s.TotalScore,
			PetEmoji:   "🥚",
			FreePick:   false,
		}
		if pet, ok := petByStudent[s.ID]; ok {
			stage := pet.CurrentStage()
			item.PetName = pet.Name
			item.PetSpecies = pet.Species
			item.PetLevel = pet.Level
			if stage.Emoji != "" {
				item.PetEmoji = stage.Emoji
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// ScoreRules 教室端积分规则（与教师端同源，首次访问补齐本校默认规则）。
func (d *DisplayService) ScoreRules(classID uint) ([]models.ScoreRule, error) {
	var class models.ClassRoom
	if err := d.db.First(&class, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}
	return NewRules(d.db).ListForClass(class.ID, class.SchoolID)
}

// DisplayPetOverviewItem 教室端宠物概览条目。
type DisplayPetOverviewItem struct {
	ID          uint   `json:"id"`
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	Name        string `json:"name"`
	Species     string `json:"species"`
	Level       int    `json:"level"`
	Exp         int    `json:"exp"`
	Mood        int    `json:"mood"`
	StageName   string `json:"stage_name"`
	Emoji       string `json:"emoji"`
}

// ClassroomPetsOverview 教室端宠物概览（含学生姓名与阶段名/emoji）。
func (d *DisplayService) ClassroomPetsOverview(classID uint) ([]DisplayPetOverviewItem, error) {
	var pets []models.Pet
	if err := d.db.Where("class_id = ?", classID).Order("id ASC").Find(&pets).Error; err != nil {
		return nil, err
	}

	names := map[uint]string{}
	if len(pets) > 0 {
		ids := make([]uint, 0, len(pets))
		for _, p := range pets {
			ids = append(ids, p.StudentID)
		}
		var students []models.Student
		if err := d.db.Where("id IN ?", ids).Find(&students).Error; err != nil {
			return nil, err
		}
		for _, s := range students {
			names[s.ID] = s.Name
		}
	}

	items := make([]DisplayPetOverviewItem, 0, len(pets))
	for _, p := range pets {
		stage := p.CurrentStage()
		items = append(items, DisplayPetOverviewItem{
			ID:          p.ID,
			StudentID:   p.StudentID,
			StudentName: names[p.StudentID],
			Name:        p.Name,
			Species:     p.Species,
			Level:       p.Level,
			Exp:         p.Experience,
			Mood:        p.Mood,
			StageName:   stage.Name,
			Emoji:       stage.Emoji,
		})
	}
	return items, nil
}

// DisplayLeaderboardItem 大屏排行榜条目（字段名同 Laravel quickLeaderboard：no / id / score）。
type DisplayLeaderboardItem struct {
	Rank  int    `json:"rank"`
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
	No    string `json:"no"`
}

// QuickLeaderboard 大屏排行榜（本班前 20 名学生，按总积分降序）。
func (d *DisplayService) QuickLeaderboard(classID uint) ([]DisplayLeaderboardItem, error) {
	students, err := d.activeStudents(classID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(students, func(i, j int) bool {
		if students[i].TotalScore != students[j].TotalScore {
			return students[i].TotalScore > students[j].TotalScore
		}
		return students[i].ID < students[j].ID
	})

	items := make([]DisplayLeaderboardItem, 0, len(students))
	for i, s := range students {
		if i >= 20 {
			break
		}
		items = append(items, DisplayLeaderboardItem{
			Rank:  i + 1,
			ID:    s.ID,
			Name:  s.Name,
			Score: s.TotalScore,
			No:    s.StudentNo,
		})
	}
	return items, nil
}

// ============================================================
// 内部辅助
// ============================================================

func (d *DisplayService) activeStudentCount(classID uint) (int64, error) {
	var count int64
	err := d.db.Model(&models.Student{}).Where("class_id = ? AND status = ?", classID, "active").
		Count(&count).Error
	return count, err
}

func (d *DisplayService) activeStudents(classID uint) ([]models.Student, error) {
	var students []models.Student
	if err := d.db.Where("class_id = ? AND status = ?", classID, "active").Order("id ASC").
		Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (d *DisplayService) petsByStudent(students []models.Student) (map[uint]models.Pet, error) {
	out := map[uint]models.Pet{}
	if len(students) == 0 {
		return out, nil
	}

	ids := make([]uint, 0, len(students))
	for _, s := range students {
		ids = append(ids, s.ID)
	}

	var pets []models.Pet
	if err := d.db.Where("student_id IN ?", ids).Find(&pets).Error; err != nil {
		return nil, err
	}
	for _, pet := range pets {
		out[pet.StudentID] = pet
	}
	return out, nil
}

func (d *DisplayService) studentNamesOf(scores []models.Score) (map[uint]string, error) {
	names := map[uint]string{}
	if len(scores) == 0 {
		return names, nil
	}

	ids := make([]uint, 0, len(scores))
	for _, sc := range scores {
		ids = append(ids, sc.StudentID)
	}

	var students []models.Student
	if err := d.db.Where("id IN ?", ids).Find(&students).Error; err != nil {
		return nil, err
	}
	for _, s := range students {
		names[s.ID] = s.Name
	}
	return names, nil
}

func (d *DisplayService) scoreSum(studentIDs []uint, from time.Time) (int, error) {
	if len(studentIDs) == 0 {
		return 0, nil
	}

	var sum int64
	if err := d.db.Model(&models.Score{}).
		Where("student_id IN ? AND created_at >= ?", studentIDs, from).
		Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error; err != nil {
		return 0, err
	}
	return int(sum), nil
}

// recentNews 最近 20 条积分记录按文案去重后取前 5 条（同 Laravel unique('text')->take(5)）。
func (d *DisplayService) recentNews(studentIDs []uint) ([]DisplayNewsItem, error) {
	out := []DisplayNewsItem{}
	if len(studentIDs) == 0 {
		return out, nil
	}

	var scores []models.Score
	if err := d.db.Where("student_id IN ?", studentIDs).
		Order("created_at DESC, id DESC").Limit(20).Find(&scores).Error; err != nil {
		return nil, err
	}

	names, err := d.studentNamesOf(scores)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	for _, sc := range scores {
		name := names[sc.StudentID]
		if name == "" {
			name = "同学"
		}
		amount := fmt.Sprintf("%d", sc.Amount)
		if sc.Amount > 0 {
			amount = "+" + amount
		}
		icon := "📝"
		if sc.Amount > 0 {
			icon = "🎉"
		}
		text := name + " " + amount + "分 — " + sc.Reason
		if seen[text] {
			continue
		}
		seen[text] = true
		out = append(out, DisplayNewsItem{Icon: icon, Text: text})
		if len(out) >= 5 {
			break
		}
	}
	return out, nil
}

// classPetSeries 读取班级 settings 的 pet_series（键不存在时为 null，同 Laravel）。
func classPetSeries(class *models.ClassRoom) any {
	if v, ok := class.SettingsMap()["pet_series"]; ok {
		return v
	}
	return nil
}

// isoTimePtrOf 兼容 *time.Time → *string（RFC3339，等价 Laravel toIso8601String）。
func isoTimePtrOf(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// isoWeekday 今天星期几（周一=1 … 周日=7，等价 Carbon dayOfWeekIso）。
func isoWeekday(t time.Time) int {
	wd := int(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}

// sortStudentsByNo 按 student_no 的数字序 + id 升序排序。
// 等价 Laravel 的 ORDER BY CAST(student_no AS UNSIGNED) ASC, id ASC；数字转换在 Go 侧完成，
// 以避免 MySQL(SIGNED) 与 SQLite/PostgreSQL(INTEGER) 的方言差异（非数字前缀按 0 处理）。
func sortStudentsByNo(students []models.Student) {
	sort.SliceStable(students, func(i, j int) bool {
		ni, nj := leadingNumber(students[i].StudentNo), leadingNumber(students[j].StudentNo)
		if ni != nj {
			return ni < nj
		}
		return students[i].ID < students[j].ID
	})
}

// leadingNumber 取字符串的前导数字（无前导数字为 0，同 MySQL CAST('x' AS UNSIGNED)）。
func leadingNumber(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			break
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}

// relativeTime 近似 Carbon diffForHumans 的输出（Laravel 默认 locale = en）。
// 差异说明：此处按「整秒/分/时/天/周/月/年」向下取整拼接英文串，不做 Carbon 的复数与
// 分词精细化处理；仅用于大屏「xg 前」这类展示文案。
// 差异说明：按整单位向下取整拼接英文串（"5 minutes ago"），不做 Carbon 的
// "a few seconds"/"a minute ago" 等分词精细化处理；仅用于大屏展示文案。
func relativeTime(t, now time.Time) string {
	seconds := int64(now.Sub(t).Seconds())
	if seconds < 1 {
		return "just now"
	}

	switch {
	case seconds < 60:
		return fmt.Sprintf("%d second%s ago", seconds, plural(seconds))
	case seconds < 3600:
		m := seconds / 60
		return fmt.Sprintf("%d minute%s ago", m, plural(m))
	case seconds < 86400:
		h := seconds / 3600
		return fmt.Sprintf("%d hour%s ago", h, plural(h))
	case seconds < 604800:
		d := seconds / 86400
		return fmt.Sprintf("%d day%s ago", d, plural(d))
	case seconds < 2592000:
		w := seconds / 604800
		return fmt.Sprintf("%d week%s ago", w, plural(w))
	case seconds < 31536000:
		mo := seconds / 2592000
		return fmt.Sprintf("%d month%s ago", mo, plural(mo))
	default:
		y := seconds / 31536000
		return fmt.Sprintf("%d year%s ago", y, plural(y))
	}
}

func plural(n int64) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func strPtr(s string) *string { return &s }

// ============================================================
// 大屏事件消费（sse / poll 的数据源）
// ============================================================

// ConsumeDisplayEvents 消费本班自 since 之后的大屏事件（since 为 nil 时返回全部）。
func (d *DisplayService) ConsumeDisplayEvents(classID uint, since *int) ([]DisplayEventView, error) {
	return d.events.Consume(classID, since)
}

// ============================================================
// 大屏加减分
// ============================================================

// displayQuickScoreAmounts 大屏快捷加减分允许的分值（Laravel quickScore 的 in:-5,-3,-1,1,3,5）。
var displayQuickScoreAmounts = []int{-5, -3, -1, 1, 3, 5}

// DisplayQuickScoreResult 大屏快捷加减分结果（字段名同 Laravel：data.total_score）。
type DisplayQuickScoreResult struct {
	TotalScore int `json:"total_score"`
}

// DisplayGiveScoreResult 教室端加减分结果（字段名同 Laravel classroomGiveScore 的 data）。
type DisplayGiveScoreResult struct {
	StudentName string `json:"student_name"`
	Points      int    `json:"points"`
	NewScore    int    `json:"new_score"`
}

// DisplayBatchGiveResult 教室端批量加减分结果（字段名同 Laravel classroomBatchGiveScore 的 data）。
type DisplayBatchGiveResult struct {
	Count  int `json:"count"`
	Points int `json:"points"`
}

// QuickScore 大屏快捷加减分：分值限于 ±1/±3/±5，原因固定「课堂表现」，学生必须在本班。
func (d *DisplayService) QuickScore(classID, studentID uint, amount int) (*DisplayQuickScoreResult, error) {
	if !displayQuickScoreAmountAllowed(amount) {
		return nil, ErrUnprocessable("无效的分值")
	}

	student, err := d.classStudent(classID, studentID)
	if err != nil {
		return nil, err
	}

	// 操作人：班级教师，缺失时兜底 user_id = 1（Laravel `$teacherId ?: 1`）。
	operator := classTeacherID(d.db, classID)
	if operator == 0 {
		operator = 1
	}

	if _, err := d.scores.GiveScore(student, amount, "课堂表现", operator, nil); err != nil {
		return nil, err
	}
	return &DisplayQuickScoreResult{TotalScore: student.TotalScore}, nil
}

// ClassroomGiveScore 教室端单个学生加减分。
//
// 复用 Go 端 ScoreService（等价 Laravel ScoreService::giveScore 的副作用：积分记录 + 学生总分
// max(0) 兜底 + 宠物经验增减 + 审计日志 + score_update 事件）。
// 有意差异：Laravel classroomGiveScore 手写 Score::create，且不写 score_logs 审计、不发布大屏
// 事件、不校正宠物等级；Go 端统一走 ScoreService，故这些副作用会一并产生（与教师端给分同源）。
func (d *DisplayService) ClassroomGiveScore(classID, studentID uint, points int, reason string) (*DisplayGiveScoreResult, error) {
	if points == 0 {
		return nil, ErrUnprocessable("分值不能为 0")
	}
	if strings.TrimSpace(reason) == "" {
		// Laravel `$request->input('reason', '课堂评价')`：缺失时回落「课堂评价」。
		reason = "课堂评价"
	}
	if len([]rune(reason)) > 200 {
		return nil, ErrUnprocessable("原因不能超过 200 字")
	}

	student, err := d.classStudent(classID, studentID)
	if err != nil {
		return nil, err
	}

	// 班级码模式限制：单次加减分不得超过 ±30（|amount| >= 30 即拒绝，同 Laravel abs() >= 30）。
	if absInt(points) >= displayClassroomMaxPoints {
		return nil, ErrForbidden(displayClassroomMaxPointsMessage)
	}

	if _, err := d.scores.GiveScore(student, points, reason, d.scoreOperatorID(classID), nil); err != nil {
		return nil, err
	}
	return &DisplayGiveScoreResult{
		StudentName: student.Name,
		Points:      points,
		NewScore:    student.TotalScore,
	}, nil
}

// ClassroomBatchGiveScore 教室端批量加减分（多选学生，仅处理属于本班的学生）。
//
// 复用 Go 端 ScoreService.BatchGive（单事务提交，逐学生产生积分记录 + 总分 + 宠物经验 +
// 审计日志 + score_update 事件）。有意差异同 ClassroomGiveScore。
func (d *DisplayService) ClassroomBatchGiveScore(
	classID uint, studentIDs []uint, points int, reason string,
) (*DisplayBatchGiveResult, error) {
	if len(studentIDs) == 0 {
		return nil, ErrUnprocessable("请至少选择一名学生")
	}
	if len(studentIDs) > displayBatchMaxStudents {
		return nil, ErrUnprocessable(fmt.Sprintf("单次最多操作 %d 名学生", displayBatchMaxStudents))
	}
	if points == 0 {
		return nil, ErrUnprocessable("分值不能为 0")
	}
	if strings.TrimSpace(reason) == "" {
		return nil, ErrUnprocessable("请填写加减分原因")
	}
	if len([]rune(reason)) > 200 {
		return nil, ErrUnprocessable("原因不能超过 200 字")
	}

	// 班级码模式限制先于学生查询判定（同 Laravel classroomBatchGiveScore 的顺序）。
	if absInt(points) >= displayClassroomMaxPoints {
		return nil, ErrForbidden(displayClassroomMaxPointsMessage)
	}

	var students []*models.Student
	if err := d.db.Where("class_id = ? AND id IN ?", classID, studentIDs).Find(&students).Error; err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return nil, ErrUnprocessable("未找到可操作的学生")
	}

	count, err := d.scores.BatchGive(students, points, reason, d.scoreOperatorID(classID), nil)
	if err != nil {
		return nil, err
	}
	return &DisplayBatchGiveResult{Count: count, Points: points}, nil
}

// displayClassroomMaxPoints 班级码模式单次加减分上限（Laravel `abs($amount) >= 30` → 403）。
const displayClassroomMaxPoints = 30

// displayClassroomMaxPointsMessage 超限文案（逐字同 Laravel）。
const displayClassroomMaxPointsMessage = "单次加减分超过 30 分，请使用教师账号登录操作"

// displayBatchMaxStudents 批量加减分单次学生数上限（Laravel max:50）。
const displayBatchMaxStudents = 50

// classStudent 取本班学生（不在本班或不存在 → 404「学生不存在」）。
// 差异说明：Laravel classroomGiveScore 用 findOrFail（404 文案由异常处理器决定），
// 这里统一为 quickScore 的「学生不存在」文案。
func (d *DisplayService) classStudent(classID, studentID uint) (*models.Student, error) {
	var student models.Student
	err := d.db.Where("class_id = ? AND id = ?", classID, studentID).First(&student).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("学生不存在")
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// classTeacherID 班级授课教师 ID（同 Laravel getClassTeacherId）：班级 teacher_id → 平台首位教师。
// 差异说明：Laravel 在两者之间还会查 class_room_teachers 关联表，Go 端未建该表（多教师关联
// 未纳入重设计 schema），故跳过该级。
func classTeacherID(db *gorm.DB, classID uint) uint {
	var class models.ClassRoom
	if err := db.First(&class, classID).Error; err == nil && class.TeacherID != nil && *class.TeacherID != 0 {
		return *class.TeacherID
	}

	var teacher models.User
	if err := db.Where("role = ?", "teacher").Order("id ASC").First(&teacher).Error; err == nil {
		return teacher.ID
	}
	return 0
}

// scoreOperatorID 大屏加减分的操作人：班级教师 → 管理员 → 1（同 Laravel 的 `?: 1` 兜底）。
// 差异说明：Laravel 的管理员兜底查询不加学校过滤（User::whereIn('role', [...])->value('id')），
// 此处保持一致（只按 id ASC 取首位），以免跨校部署时行为漂移。
func (d *DisplayService) scoreOperatorID(classID uint) uint {
	if id := classTeacherID(d.db, classID); id != 0 {
		return id
	}

	var admin models.User
	if err := d.db.Where("role IN ?", []string{"school_admin", "admin"}).
		Order("id ASC").First(&admin).Error; err == nil {
		return admin.ID
	}
	return 1
}

// displayQuickScoreAmountAllowed 判断分值是否在大屏快捷加减分白名单内。
func displayQuickScoreAmountAllowed(amount int) bool {
	for _, allowed := range displayQuickScoreAmounts {
		if amount == allowed {
			return true
		}
	}
	return false
}

// absInt 求整数绝对值。
func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
