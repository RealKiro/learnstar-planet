// 教师端班级大屏消息服务（广播 + 通知统一收口）。
// 忠实移植自 Laravel App\Services\ClassroomMessagingService（display / send / poll）：
//   - type = banner|popup|fullscreen → 写 broadcasts 表并发布 broadcast 大屏事件；
//   - type = info|homework|event|urgent → 写 notices 表并发布 notice 大屏事件。
//
// 有意差异：Go 端 schema 重设计后只保留 class_rooms.teacher_id 这一条教师-班级关联
// （无 class_room_teachers 多角色关联表），故「是否被分配到此班级」改用 Scope.ClassIDs
// （班主任 ∪ API 机器人本校全部班级）判定，与教师端其余接口口径一致。
package services

import (
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// 课堂消息类型（同 Laravel 控制器 in: 枚举）。
const (
	// ClassroomMessageBanner 横幅广播。
	ClassroomMessageBanner = "banner"
	// ClassroomMessagePopup 弹窗广播。
	ClassroomMessagePopup = "popup"
	// ClassroomMessageFullscreen 全屏广播。
	ClassroomMessageFullscreen = "fullscreen"
	// ClassroomMessageInfo 普通通知。
	ClassroomMessageInfo = "info"
	// ClassroomMessageHomework 作业通知。
	ClassroomMessageHomework = "homework"
	// ClassroomMessageEvent 活动通知。
	ClassroomMessageEvent = "event"
	// ClassroomMessageUrgent 紧急通知。
	ClassroomMessageUrgent = "urgent"
)

// 课堂消息条数与时间窗口（同 Laravel）。
const (
	classroomDisplayBroadcastLimit = 5
	classroomDisplayNoticeLimit    = 3
	classroomDisplayNoticeDays     = 7
	classroomDisplayScoreLimit     = 20
	classroomDisplayScoreHours     = 4
	classroomPollBroadcastLimit    = 5
	classroomPollNoticeLimit       = 3
	classroomPollDefaultMinutes    = 5
)

// broadcastMessageTypes 走广播表的类型。
var broadcastMessageTypes = map[string]bool{
	ClassroomMessageBanner:     true,
	ClassroomMessagePopup:      true,
	ClassroomMessageFullscreen: true,
}

// noticeMessageTypes 走通知表的类型。
var noticeMessageTypes = map[string]bool{
	ClassroomMessageInfo:     true,
	ClassroomMessageHomework: true,
	ClassroomMessageEvent:    true,
	ClassroomMessageUrgent:   true,
}

// IsClassroomMessageType 判断是否为合法的课堂消息类型（banner/popup/fullscreen/info/homework/event/urgent）。
func IsClassroomMessageType(msgType string) bool {
	return broadcastMessageTypes[msgType] || noticeMessageTypes[msgType]
}

// ClassroomMessagingService 班级大屏聚合数据与消息读写服务。
type ClassroomMessagingService struct {
	db     *gorm.DB
	scope  *Scope
	events *DisplayEvents
}

// NewClassroomMessagingService 创建课堂消息服务（内部自带大屏事件发布器）。
func NewClassroomMessagingService(db *gorm.DB, scope *Scope) *ClassroomMessagingService {
	return &ClassroomMessagingService{db: db, scope: scope, events: NewDisplayEvents(db)}
}

// AccessibleClassIDs 教师可访问的班级 ID 列表（发送前的越权预检依据）。
func (s *ClassroomMessagingService) AccessibleClassIDs(u *models.User) ([]uint, error) {
	return s.scope.ClassIDs(u)
}

// ============================================================
// 班级大屏聚合数据（GET /teacher/classroom/display）
// ============================================================

// ClassroomDisplayPet 大屏学生宠物总览的一行。
type ClassroomDisplayPet struct {
	StudentID   uint    `json:"student_id"`
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
}

// ClassroomDisplayBroadcast 大屏生效广播（created_at 为 diffForHumans 展示文案）。
type ClassroomDisplayBroadcast struct {
	ID             uint   `json:"id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	DisplaySeconds int    `json:"display_seconds"`
	VoiceEnabled   bool   `json:"voice_enabled"`
	CreatedAt      string `json:"created_at"`
}

// ClassroomDisplayNotice 大屏近期通知（published_at 为 diffForHumans 展示文案）。
type ClassroomDisplayNotice struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	PublishedAt string `json:"published_at"`
}

// ClassroomDisplayScore 大屏最近积分（time 为 diffForHumans 展示文案）。
type ClassroomDisplayScore struct {
	StudentName *string `json:"student_name"`
	Amount      int     `json:"amount"`
	Reason      string  `json:"reason"`
	Time        string  `json:"time"`
}

// ClassroomDisplayData 班级大屏聚合数据。
type ClassroomDisplayData struct {
	ClassName    string                      `json:"class_name"`
	Grade        string                      `json:"grade"`
	StudentCount int                         `json:"student_count"`
	Pets         []ClassroomDisplayPet       `json:"pets"`
	Broadcasts   []ClassroomDisplayBroadcast `json:"broadcasts"`
	Notices      []ClassroomDisplayNotice    `json:"notices"`
	RecentScores []ClassroomDisplayScore     `json:"recent_scores"`
}

// Display 班级大屏聚合数据（学生宠物总览 + 生效广播 + 近期通知 + 最近积分）。
//
// classID 为空 → 400「请先选择班级」；未分配该班级 → 403「您未被分配到此班级」。
func (s *ClassroomMessagingService) Display(u *models.User, classID uint) (*ClassroomDisplayData, error) {
	if classID == 0 {
		return nil, ErrBadRequest("请先选择班级")
	}

	accessible, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if !containsUintValue(accessible, classID) {
		return nil, ErrForbidden("您未被分配到此班级")
	}

	var classRoom models.ClassRoom
	if err := s.db.First(&classRoom, classID).Error; err != nil {
		return nil, ErrNotFound("班级不存在")
	}

	var students []models.Student
	if err := s.db.Where("class_id = ? AND status = ?", classID, "active").
		Order("id ASC").Find(&students).Error; err != nil {
		return nil, err
	}
	sortStudentsByNo(students)

	petByStudent := map[uint]models.Pet{}
	if len(students) > 0 {
		ids := make([]uint, 0, len(students))
		for _, st := range students {
			ids = append(ids, st.ID)
		}
		var pets []models.Pet
		if err := s.db.Where("student_id IN ?", ids).Find(&pets).Error; err != nil {
			return nil, err
		}
		for _, pet := range pets {
			petByStudent[pet.StudentID] = pet
		}
	}

	// 无宠物时 Laravel 回退为 {emoji:'🤔', name:'未孵化', title:''}。
	petRows := make([]ClassroomDisplayPet, 0, len(students))
	for _, st := range students {
		row := ClassroomDisplayPet{
			StudentID:   st.ID,
			StudentName: st.Name,
			TotalScore:  st.TotalScore,
			Emoji:       "🤔",
			StageName:   "未孵化",
		}
		if pet, ok := petByStudent[st.ID]; ok {
			stage := pet.CurrentStage()
			name := pet.Name
			species := pet.Species
			row.HasPet = true
			row.PetName = &name
			row.PetSpecies = &species
			row.Level = pet.Level
			row.Experience = pet.Experience
			row.Mood = pet.Mood
			row.Emoji = stage.Emoji
			row.StageName = stage.Name
		}
		petRows = append(petRows, row)
	}

	now := util.Now()

	var broadcasts []models.Broadcast
	if err := s.db.Where("class_id = ? AND status IN ?", classID, []string{"pending", "sent"}).
		Order("created_at DESC, id DESC").Limit(classroomDisplayBroadcastLimit).
		Find(&broadcasts).Error; err != nil {
		return nil, err
	}
	broadcastRows := make([]ClassroomDisplayBroadcast, 0, len(broadcasts))
	for _, b := range broadcasts {
		broadcastRows = append(broadcastRows, ClassroomDisplayBroadcast{
			ID:             b.ID,
			Content:        b.Content,
			Type:           b.Type,
			DisplaySeconds: b.DisplaySeconds,
			VoiceEnabled:   b.VoiceEnabled,
			CreatedAt:      relativeTime(b.CreatedAt, now),
		})
	}

	var notices []models.Notice
	if err := s.db.Where("class_id = ? AND is_published = ? AND published_at >= ?",
		classID, true, now.AddDate(0, 0, -classroomDisplayNoticeDays)).
		Order("published_at DESC, id DESC").Limit(classroomDisplayNoticeLimit).
		Find(&notices).Error; err != nil {
		return nil, err
	}
	noticeRows := make([]ClassroomDisplayNotice, 0, len(notices))
	for _, n := range notices {
		noticeRows = append(noticeRows, ClassroomDisplayNotice{
			ID:          n.ID,
			Title:       n.Title,
			Content:     n.Content,
			Type:        n.Type,
			PublishedAt: relativeTimeOrEmpty(n.PublishedAt, now),
		})
	}

	var scores []models.Score
	if err := s.db.Where("class_id = ? AND created_at >= ?", classID, now.Add(-classroomDisplayScoreHours*time.Hour)).
		Order("created_at DESC, id DESC").Limit(classroomDisplayScoreLimit).
		Find(&scores).Error; err != nil {
		return nil, err
	}
	names, err := s.studentNamesOf(scores)
	if err != nil {
		return nil, err
	}
	scoreRows := make([]ClassroomDisplayScore, 0, len(scores))
	for _, sc := range scores {
		item := ClassroomDisplayScore{
			Amount: sc.Amount,
			Reason: sc.Reason,
			Time:   relativeTime(sc.CreatedAt, now),
		}
		if name, ok := names[sc.StudentID]; ok {
			n := name
			item.StudentName = &n
		}
		scoreRows = append(scoreRows, item)
	}

	return &ClassroomDisplayData{
		ClassName:    classRoom.Name,
		Grade:        classRoom.Grade,
		StudentCount: len(students),
		Pets:         petRows,
		Broadcasts:   broadcastRows,
		Notices:      noticeRows,
		RecentScores: scoreRows,
	}, nil
}

// studentNamesOf 取一批积分记录对应的学生姓名表。
func (s *ClassroomMessagingService) studentNamesOf(scores []models.Score) (map[uint]string, error) {
	names := map[uint]string{}
	if len(scores) == 0 {
		return names, nil
	}
	ids := make([]uint, 0, len(scores))
	seen := map[uint]bool{}
	for _, sc := range scores {
		if seen[sc.StudentID] {
			continue
		}
		seen[sc.StudentID] = true
		ids = append(ids, sc.StudentID)
	}
	var students []models.Student
	if err := s.db.Where("id IN ?", ids).Find(&students).Error; err != nil {
		return nil, err
	}
	for _, st := range students {
		names[st.ID] = st.Name
	}
	return names, nil
}

// ============================================================
// 发送班级消息（POST /teacher/classroom/messages）
// ============================================================

// ClassroomMessageInput 发送参数（控制器已完成 422 校验与默认值填充）。
type ClassroomMessageInput struct {
	Content        string
	Title          *string
	DisplaySeconds int
	Voice          bool
}

// ClassroomMessageData 发送结果的 data 段（type 为落库表名 broadcast / notice）。
type ClassroomMessageData struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

// ClassroomMessageResult 发送结果（Laravel 直接返回 {message, data}）。
type ClassroomMessageResult struct {
	Message string               `json:"message"`
	Data    ClassroomMessageData `json:"data"`
}

// Send 发送班级消息：按类型分流到广播表 / 通知表，并推送班级大屏。
// accessible 为教师可访问班级；classID 不在其中 → 403「无权限」。
func (s *ClassroomMessagingService) Send(
	u *models.User,
	accessible []uint,
	classID uint,
	msgType string,
	in ClassroomMessageInput,
) (*ClassroomMessageResult, error) {
	if !containsUintValue(accessible, classID) {
		return nil, ErrForbidden("无权限")
	}

	if broadcastMessageTypes[msgType] {
		now := util.Now()
		broadcast := models.Broadcast{
			SchoolID:       u.SchoolID,
			ClassID:        classID,
			TeacherID:      &u.ID,
			Content:        in.Content,
			Type:           msgType,
			VoiceEnabled:   in.Voice,
			DisplaySeconds: in.DisplaySeconds,
			Status:         "sent",
			SentAt:         &now,
		}
		if err := s.db.Create(&broadcast).Error; err != nil {
			return nil, err
		}
		// voice_enabled 带 default:true：GORM 在 Create 时跳过 false 零值，需显式回写。
		if err := forceFalseBool(s.db, &models.Broadcast{}, broadcast.ID, "voice_enabled", in.Voice); err != nil {
			return nil, err
		}

		publishDisplayEvent(s.events, classID, DisplayEventBroadcast, DisplayBroadcastData{
			ID:             broadcast.ID,
			Type:           broadcast.Type,
			Content:        broadcast.Content,
			DisplaySeconds: broadcast.DisplaySeconds,
			VoiceEnabled:   in.Voice,
			CreatedAt:      isoTimeOrEmpty(&broadcast.CreatedAt),
		})

		return &ClassroomMessageResult{
			Message: "广播已发送",
			Data:    ClassroomMessageData{ID: broadcast.ID, Type: "broadcast"},
		}, nil
	}

	// 通知类型：info / homework / event / urgent
	title := "通知"
	if msgType == ClassroomMessageUrgent {
		title = "紧急通知"
	}
	if in.Title != nil {
		title = *in.Title
	}

	now := util.Now()
	notice := models.Notice{
		ClassID:     classID,
		SchoolID:    u.SchoolID,
		Title:       title,
		Content:     in.Content,
		Type:        msgType,
		PublishedBy: u.ID,
		IsPublished: true,
		PublishedAt: &now,
	}
	if err := s.db.Create(&notice).Error; err != nil {
		return nil, err
	}

	publishDisplayEvent(s.events, classID, DisplayEventNotice, DisplayNoticeData{
		ID:          notice.ID,
		Title:       notice.Title,
		Content:     notice.Content,
		Type:        notice.Type,
		PublishedAt: isoTimeOrEmpty(notice.PublishedAt),
	})

	return &ClassroomMessageResult{
		Message: "通知已发布",
		Data:    ClassroomMessageData{ID: notice.ID, Type: "notice"},
	}, nil
}

// ============================================================
// 大屏轮询增量消息（GET /teacher/classroom/messages）
// ============================================================

// ClassroomPollBroadcast 轮询返回的广播（created_at 为 ISO8601）。
type ClassroomPollBroadcast struct {
	ID             uint   `json:"id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	DisplaySeconds int    `json:"display_seconds"`
	VoiceEnabled   bool   `json:"voice_enabled"`
	CreatedAt      string `json:"created_at"`
}

// ClassroomPollNotice 轮询返回的通知（published_at 为 ISO8601）。
type ClassroomPollNotice struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	PublishedAt string `json:"published_at"`
}

// ClassroomPollResult 轮询结果。
type ClassroomPollResult struct {
	Broadcasts []ClassroomPollBroadcast `json:"broadcasts"`
	Notices    []ClassroomPollNotice    `json:"notices"`
	PolledAt   string                   `json:"polled_at"`
}

// Poll 大屏轮询增量消息（since 之后的广播与通知；since 缺省为 5 分钟前）。
// classID 为空 → 400「请先选择班级」（同 Laravel，此处不校验班级分配）。
func (s *ClassroomMessagingService) Poll(classID uint, since *time.Time) (*ClassroomPollResult, error) {
	if classID == 0 {
		return nil, ErrBadRequest("请先选择班级")
	}

	now := util.Now()
	sinceTime := now.Add(-classroomPollDefaultMinutes * time.Minute)
	if since != nil {
		sinceTime = *since
	}

	var broadcasts []models.Broadcast
	if err := s.db.Where("class_id = ? AND created_at >= ? AND status IN ?", classID, sinceTime, []string{"sent"}).
		Order("created_at DESC, id DESC").Limit(classroomPollBroadcastLimit).
		Find(&broadcasts).Error; err != nil {
		return nil, err
	}
	broadcastRows := make([]ClassroomPollBroadcast, 0, len(broadcasts))
	for _, b := range broadcasts {
		broadcastRows = append(broadcastRows, ClassroomPollBroadcast{
			ID:             b.ID,
			Content:        b.Content,
			Type:           b.Type,
			DisplaySeconds: b.DisplaySeconds,
			VoiceEnabled:   b.VoiceEnabled,
			CreatedAt:      isoTimeOrEmpty(&b.CreatedAt),
		})
	}

	var notices []models.Notice
	if err := s.db.Where("class_id = ? AND is_published = ? AND published_at >= ?", classID, true, sinceTime).
		Order("published_at DESC, id DESC").Limit(classroomPollNoticeLimit).
		Find(&notices).Error; err != nil {
		return nil, err
	}
	noticeRows := make([]ClassroomPollNotice, 0, len(notices))
	for _, n := range notices {
		noticeRows = append(noticeRows, ClassroomPollNotice{
			ID:          n.ID,
			Title:       n.Title,
			Content:     n.Content,
			Type:        n.Type,
			PublishedAt: isoTimeOrEmpty(n.PublishedAt),
		})
	}

	return &ClassroomPollResult{
		Broadcasts: broadcastRows,
		Notices:    noticeRows,
		PolledAt:   now.Format(time.RFC3339),
	}, nil
}

// relativeTimeOrEmpty 可空时间的 diffForHumans 展示文案（nil 输出空串）。
func relativeTimeOrEmpty(t *time.Time, now time.Time) string {
	if t == nil {
		return ""
	}
	return relativeTime(*t, now)
}

// containsUintValue 判断 uint 切片是否包含目标值。
func containsUintValue(list []uint, value uint) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
