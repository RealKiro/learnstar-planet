// 课表服务：科目 / 节次（学校级共享）、排课读写、修改申请审核流、CSES 导出。
// 忠实移植自 Laravel App\Services\TimetableService（save/bootstrap/审批流/forTeacher/toCses）。
package services

import (
	"encoding/json"
	"hash/crc32"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// TimetableService 课表业务服务。
type TimetableService struct {
	db *gorm.DB
	// seedSource 排课随机种子源（可注入，便于测试）；nil 时按时间种子。
	seedSource func() int64
}

// NewTimetableService 创建课表服务。
func NewTimetableService(db *gorm.DB) *TimetableService {
	return &TimetableService{db: db}
}

// SetSeedSource 注入排课随机种子源（传 nil 恢复按时间种子）。
// 仅用于测试：注入后调用方需自行保证并发安全（生产路径不注入、每次按时间取种子）。
func (s *TimetableService) SetSeedSource(src func() int64) {
	s.seedSource = src
}

// timetableWeekTypes 合法周次（单双周）。
var timetableWeekTypes = map[string]bool{"all": true, "odd": true, "even": true}

// weekdayLabels 星期中文名（CSES 的 enable_day 为 1-7）。
var weekdayLabels = map[int]string{
	1: "星期一", 2: "星期二", 3: "星期三", 4: "星期四",
	5: "星期五", 6: "星期六", 7: "星期日",
}

// subjectPalette 科目默认色板：新科目无颜色时按名称稳定分配（同名同色）。
var subjectPalette = []string{
	"#5B8FF9", "#5AD8A6", "#F6BD16", "#E8684A", "#6DC8EC",
	"#9270CA", "#FF9D4D", "#269A99", "#FF99C3", "#A0D911",
	"#722ED1", "#13C2C2", "#FA8C16", "#EB2F96", "#52C41A",
}

// PaletteColor 按科目名从色板稳定取色（跨班级 / 跨导入一致）。
func (s *TimetableService) PaletteColor(name string) string {
	return subjectPalette[int(crc32.ChecksumIEEE([]byte(name)))%len(subjectPalette)]
}

// TimetableSubjectInput 科目输入（color/simplified_name 缺省时保留库中已有值）。
type TimetableSubjectInput struct {
	Name           string  `json:"name"`
	SimplifiedName *string `json:"simplified_name"`
	Color          *string `json:"color"`
}

// TimetablePeriodInput 节次输入。
type TimetablePeriodInput struct {
	PeriodIndex int    `json:"period_index"`
	Name        string `json:"name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// TimetableEntryInput 排课输入（API 边界上用科目名交互）。
type TimetableEntryInput struct {
	Weekday     int     `json:"weekday"`
	PeriodIndex int     `json:"period_index"`
	SubjectName string  `json:"subject_name"`
	WeekType    string  `json:"week_type"`
	TeacherName *string `json:"teacher_name"`
	Room        *string `json:"room"`
}

// TimetablePayload 一次保存 / 提交申请的完整课表快照。
type TimetablePayload struct {
	Subjects []TimetableSubjectInput `json:"subjects"`
	Periods  []TimetablePeriodInput  `json:"periods"`
	Entries  []TimetableEntryInput   `json:"entries"`
}

// TimetableSubjectView 科目视图。
type TimetableSubjectView struct {
	Name           string `json:"name"`
	SimplifiedName string `json:"simplified_name"`
	Color          string `json:"color"`
}

// TimetablePeriodView 节次视图。
type TimetablePeriodView struct {
	PeriodIndex int    `json:"period_index"`
	Name        string `json:"name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// TimetableEntryView 排课格子视图（科目名形式）。
type TimetableEntryView struct {
	Weekday     int     `json:"weekday"`
	PeriodIndex int     `json:"period_index"`
	WeekType    string  `json:"week_type"`
	SubjectName *string `json:"subject_name"`
	TeacherName *string `json:"teacher_name"`
	Room        *string `json:"room"`
}

// TimetableBootstrap 课表初始化数据。
type TimetableBootstrap struct {
	Subjects []TimetableSubjectView `json:"subjects"`
	Periods  []TimetablePeriodView  `json:"periods"`
	Entries  []TimetableEntryView   `json:"entries"`
}

// TimetableTeacherEntryView 教师周课表条目。
type TimetableTeacherEntryView struct {
	ClassID     uint    `json:"class_id"`
	ClassName   *string `json:"class_name"`
	Weekday     int     `json:"weekday"`
	PeriodIndex int     `json:"period_index"`
	WeekType    string  `json:"week_type"`
	SubjectName *string `json:"subject_name"`
	Room        *string `json:"room"`
}

// TimetableTeacherSchedule 教师周课表。
type TimetableTeacherSchedule struct {
	TeacherName string                      `json:"teacher_name"`
	Subjects    []TimetableSubjectView      `json:"subjects"`
	Periods     []TimetablePeriodView       `json:"periods"`
	Entries     []TimetableTeacherEntryView `json:"entries"`
}

// TimetableChangeView 修改申请视图（列表用）。
type TimetableChangeView struct {
	ID            uint    `json:"id"`
	ClassID       uint    `json:"class_id"`
	EntryCount    int     `json:"entry_count"`
	Status        string  `json:"status"`
	ReviewNote    *string `json:"review_note"`
	CreatedAt     string  `json:"created_at"`
	ReviewedAt    string  `json:"reviewed_at"`
	RequesterName *string `json:"requester_name"`
	ReviewerName  *string `json:"reviewer_name"`
}

// AdminTimetableChangeView 管理员待办视图（带班级信息）。
type AdminTimetableChangeView struct {
	ID            uint    `json:"id"`
	ClassID       uint    `json:"class_id"`
	ClassName     *string `json:"class_name"`
	Grade         *string `json:"grade"`
	EntryCount    int     `json:"entry_count"`
	Status        string  `json:"status"`
	ReviewNote    *string `json:"review_note"`
	CreatedAt     string  `json:"created_at"`
	ReviewedAt    string  `json:"reviewed_at"`
	RequesterName *string `json:"requester_name"`
	ReviewerName  *string `json:"reviewer_name"`
}

// Save 整体保存课表：科目 upsert → 节次 upsert → 该班排课先清后插。
// 科目未传 color / simplified_name 时保留库中已有值，两者皆无时按色板自动配色。
// 返回本次自动配色的科目数。
func (s *TimetableService) Save(classID, schoolID uint, payload TimetablePayload) (int, error) {
	autoColored := 0

	err := s.db.Transaction(func(tx *gorm.DB) error {
		for i, sub := range payload.Subjects {
			name := strings.TrimSpace(sub.Name)
			if name == "" {
				continue
			}
			var existing models.Subject
			_ = tx.Where("school_id = ? AND name = ?", schoolID, name).First(&existing).Error
			color := ""
			if sub.Color != nil {
				color = *sub.Color
			} else if existing.ID != 0 {
				color = existing.Color
			}
			if color == "" {
				color = s.PaletteColor(name)
				autoColored++
			}
			simplified := ""
			if sub.SimplifiedName != nil {
				simplified = *sub.SimplifiedName
			} else if existing.ID != 0 {
				simplified = existing.SimplifiedName
			}
			if existing.ID == 0 {
				if err := tx.Create(&models.Subject{
					SchoolID: schoolID, Name: name, SimplifiedName: simplified, Color: color, SortOrder: i,
				}).Error; err != nil {
					return err
				}
			} else {
				existing.SimplifiedName = simplified
				existing.Color = color
				existing.SortOrder = i
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			}
		}

		for _, p := range payload.Periods {
			if p.PeriodIndex < 1 || p.StartTime == "" || p.EndTime == "" {
				continue
			}
			name := p.Name
			if name == "" {
				name = "第" + itoa(p.PeriodIndex) + "节"
			}
			var period models.ClassPeriod
			if err := tx.Where("school_id = ? AND period_index = ?", schoolID, p.PeriodIndex).First(&period).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
				period = models.ClassPeriod{SchoolID: schoolID, PeriodIndex: p.PeriodIndex}
			}
			period.Name = name
			period.StartTime = p.StartTime
			period.EndTime = p.EndTime
			if err := tx.Save(&period).Error; err != nil {
				return err
			}
		}

		idByName := map[string]uint{}
		var subjects []models.Subject
		if err := tx.Where("school_id = ?", schoolID).Find(&subjects).Error; err != nil {
			return err
		}
		for _, sub := range subjects {
			idByName[sub.Name] = sub.ID
		}

		if err := tx.Where("class_id = ?", classID).Delete(&models.TimetableEntry{}).Error; err != nil {
			return err
		}
		for _, e := range payload.Entries {
			name := strings.TrimSpace(e.SubjectName)
			subjectID, ok := idByName[name]
			weekday := e.Weekday
			periodIndex := e.PeriodIndex
			if !ok || weekday < 1 || weekday > 7 || periodIndex < 1 {
				continue
			}
			weekType := e.WeekType
			if !timetableWeekTypes[weekType] {
				weekType = "all"
			}
			var teacher, room string
			if e.TeacherName != nil {
				teacher = *e.TeacherName
			}
			if e.Room != nil {
				room = *e.Room
			}
			entry := models.TimetableEntry{
				ClassID: classID, Weekday: weekday, PeriodIndex: periodIndex,
				WeekType: weekType, SubjectID: &subjectID,
				TeacherName: teacher, Room: room,
			}
			if err := tx.Create(&entry).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return autoColored, nil
}

// Bootstrap 初始化数据：科目 + 节次 + 该班排课（以科目名返回）。
func (s *TimetableService) Bootstrap(classID, schoolID uint) (*TimetableBootstrap, error) {
	var subjects []models.Subject
	if err := s.db.Where("school_id = ?", schoolID).Order("sort_order ASC, id ASC").Find(&subjects).Error; err != nil {
		return nil, err
	}
	nameById := map[uint]string{}
	for _, sub := range subjects {
		nameById[sub.ID] = sub.Name
	}

	var entries []models.TimetableEntry
	if err := s.db.Where("class_id = ?", classID).Order("weekday ASC, period_index ASC").Find(&entries).Error; err != nil {
		return nil, err
	}

	result := &TimetableBootstrap{Subjects: []TimetableSubjectView{}, Periods: []TimetablePeriodView{}, Entries: []TimetableEntryView{}}
	for _, sub := range subjects {
		result.Subjects = append(result.Subjects, TimetableSubjectView{
			Name: sub.Name, SimplifiedName: sub.SimplifiedName, Color: sub.Color,
		})
	}
	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).Order("period_index ASC").Find(&periods).Error; err != nil {
		return nil, err
	}
	for _, p := range periods {
		result.Periods = append(result.Periods, TimetablePeriodView{
			PeriodIndex: p.PeriodIndex, Name: p.Name, StartTime: p.StartTime, EndTime: p.EndTime,
		})
	}
	for _, e := range entries {
		weekType := e.WeekType
		if weekType == "" {
			weekType = "all"
		}
		var namePtr *string
		if e.SubjectID != nil {
			if n, ok := nameById[*e.SubjectID]; ok {
				namePtr = &n
			}
		}
		var teacher, room *string
		if e.TeacherName != "" {
			t := e.TeacherName
			teacher = &t
		}
		if e.Room != "" {
			r := e.Room
			room = &r
		}
		result.Entries = append(result.Entries, TimetableEntryView{
			Weekday: e.Weekday, PeriodIndex: e.PeriodIndex, WeekType: weekType,
			SubjectName: namePtr, TeacherName: teacher, Room: room,
		})
	}
	return result, nil
}

// AdminSave 管理员直接保存课表（即时生效，不走审核流）。
// 保存后该班所有待审申请自动作废，防止旧快照获批后覆盖管理员的直接修改。
func (s *TimetableService) AdminSave(classID, schoolID uint, payload TimetablePayload) error {
	if _, err := s.Save(classID, schoolID, payload); err != nil {
		return err
	}
	return s.db.Model(&models.TimetableChangeRequest{}).
		Where("class_id = ? AND status = ?", classID, models.TimetableChangePending).
		Updates(map[string]any{
			"status":      models.TimetableChangeRejected,
			"review_note": "管理员已直接修改课表，本申请自动作废",
			"reviewed_at": time.Now(),
		}).Error
}

// SubmitChange 提交课表修改申请（不直接生效，待管理员审核）。
func (s *TimetableService) SubmitChange(classID, schoolID, teacherID uint, payload TimetablePayload) (*models.TimetableChangeRequest, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req := models.TimetableChangeRequest{
		SchoolID: schoolID, ClassID: classID, RequestedBy: teacherID,
		Payload: string(raw), EntryCount: len(payload.Entries),
		Status: models.TimetableChangePending,
	}
	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// ApproveChange 审核通过并应用申请的课表快照（幂等：仅 pending 可通过）。
func (s *TimetableService) ApproveChange(requestID, reviewerID uint, note *string) (bool, error) {
	applied := false
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var req models.TimetableChangeRequest
		if err := tx.Where("id = ?", requestID).First(&req).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		if req.Status != models.TimetableChangePending {
			return nil
		}
		svc := TimetableService{db: tx}
		var payload TimetablePayload
		if err := json.Unmarshal([]byte(req.Payload), &payload); err != nil {
			return err
		}
		if _, err := svc.Save(req.ClassID, req.SchoolID, payload); err != nil {
			return err
		}

		now := time.Now()
		updates := map[string]any{
			"status":      models.TimetableChangeApproved,
			"reviewed_by": reviewerID,
			"reviewed_at": now,
		}
		if note != nil {
			updates["review_note"] = *note
		}
		if err := tx.Model(&req).Updates(updates).Error; err != nil {
			return err
		}

		// 该班其余 pending 申请自动作废（课表已被本次覆盖）
		if err := tx.Model(&models.TimetableChangeRequest{}).
			Where("class_id = ? AND status = ?", req.ClassID, models.TimetableChangePending).
			Updates(map[string]any{
				"status":      models.TimetableChangeRejected,
				"review_note": "已由更新的申请取代",
				"reviewed_at": now,
			}).Error; err != nil {
			return err
		}
		applied = true
		return nil
	})
	return applied, err
}

// RejectChange 驳回申请（幂等：仅 pending 可驳回）。
func (s *TimetableService) RejectChange(requestID, reviewerID uint, note *string) (bool, error) {
	var req models.TimetableChangeRequest
	if err := s.db.Where("id = ?", requestID).First(&req).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if req.Status != models.TimetableChangePending {
		return false, nil
	}
	updates := map[string]any{
		"status":      models.TimetableChangeRejected,
		"reviewed_by": reviewerID,
		"reviewed_at": time.Now(),
	}
	if note != nil {
		updates["review_note"] = *note
	}
	return true, s.db.Model(&req).Updates(updates).Error
}

// ListChangesForClass 某班申请历史（最新在前，最多 20 条）。
func (s *TimetableService) ListChangesForClass(classID uint) ([]TimetableChangeView, error) {
	var reqs []models.TimetableChangeRequest
	if err := s.db.Where("class_id = ?", classID).Order("id DESC").Limit(20).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return s.changeViews(reqs), nil
}

// ListChangesForSchool 管理员待办列表（全部状态，可按状态过滤，最新在前，最多 100 条）。
func (s *TimetableService) ListChangesForSchool(schoolID uint, status string) ([]AdminTimetableChangeView, error) {
	q := s.db.Where("school_id = ?", schoolID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var reqs []models.TimetableChangeRequest
	if err := q.Order("id DESC").Limit(100).Find(&reqs).Error; err != nil {
		return nil, err
	}

	classIDs := map[uint]bool{}
	for _, r := range reqs {
		classIDs[r.ClassID] = true
	}
	className := map[uint]models.ClassRoom{}
	if len(classIDs) > 0 {
		var classes []models.ClassRoom
		if err := s.db.Where("id IN ?", keysOf(classIDs)).Find(&classes).Error; err != nil {
			return nil, err
		}
		for _, cl := range classes {
			className[cl.ID] = cl
		}
	}

	views := []AdminTimetableChangeView{}
	for _, r := range reqs {
		v := AdminTimetableChangeView{
			ID: r.ID, ClassID: r.ClassID, EntryCount: r.EntryCount, Status: r.Status,
			ReviewNote: r.ReviewNote, CreatedAt: isoTime(r.CreatedAt),
			ReviewedAt: isoTimePtr(r.ReviewedAt),
		}
		if name := s.userName(r.RequestedBy); name != "" {
			v.RequesterName = &name
		}
		if r.ReviewedBy != nil {
			if name := s.userName(*r.ReviewedBy); name != "" {
				v.ReviewerName = &name
			}
		}
		if cl, ok := className[r.ClassID]; ok {
			name, grade := cl.Name, cl.Grade
			v.ClassName = &name
			v.Grade = &grade
		}
		views = append(views, v)
	}
	return views, nil
}

func (s *TimetableService) changeViews(reqs []models.TimetableChangeRequest) []TimetableChangeView {
	views := []TimetableChangeView{}
	for _, r := range reqs {
		v := TimetableChangeView{
			ID: r.ID, ClassID: r.ClassID, EntryCount: r.EntryCount, Status: r.Status,
			ReviewNote: r.ReviewNote, CreatedAt: isoTime(r.CreatedAt),
			ReviewedAt: isoTimePtr(r.ReviewedAt),
		}
		if name := s.userName(r.RequestedBy); name != "" {
			v.RequesterName = &name
		}
		if r.ReviewedBy != nil {
			if name := s.userName(*r.ReviewedBy); name != "" {
				v.ReviewerName = &name
			}
		}
		views = append(views, v)
	}
	return views
}

func (s *TimetableService) userName(userID uint) string {
	var u models.User
	if err := s.db.Select("name").Where("id = ?", userID).First(&u).Error; err != nil {
		return ""
	}
	return u.Name
}

// ForTeacher 某教师的周课表（全校各班聚合，按排课里的教师姓名匹配）。
func (s *TimetableService) ForTeacher(schoolID uint, teacherName string) (*TimetableTeacherSchedule, error) {
	var subjects []models.Subject
	if err := s.db.Where("school_id = ?", schoolID).Order("sort_order ASC, id ASC").Find(&subjects).Error; err != nil {
		return nil, err
	}
	nameById := map[uint]string{}
	for _, sub := range subjects {
		nameById[sub.ID] = sub.Name
	}

	var classes []models.ClassRoom
	if err := s.db.Where("school_id = ?", schoolID).Find(&classes).Error; err != nil {
		return nil, err
	}
	className := map[uint]string{}
	for _, cl := range classes {
		className[cl.ID] = cl.Name
	}

	var entries []models.TimetableEntry
	if err := s.db.Where("class_id IN ? AND teacher_name = ?", keysOfClass(className), teacherName).
		Order("weekday ASC, period_index ASC").Find(&entries).Error; err != nil {
		return nil, err
	}

	result := &TimetableTeacherSchedule{
		TeacherName: teacherName,
		Subjects:    []TimetableSubjectView{},
		Periods:     []TimetablePeriodView{},
		Entries:     []TimetableTeacherEntryView{},
	}
	for _, sub := range subjects {
		result.Subjects = append(result.Subjects, TimetableSubjectView{
			Name: sub.Name, SimplifiedName: sub.SimplifiedName, Color: sub.Color,
		})
	}
	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).Order("period_index ASC").Find(&periods).Error; err != nil {
		return nil, err
	}
	for _, p := range periods {
		result.Periods = append(result.Periods, TimetablePeriodView{
			PeriodIndex: p.PeriodIndex, Name: p.Name, StartTime: p.StartTime, EndTime: p.EndTime,
		})
	}
	for _, e := range entries {
		if e.SubjectID == nil {
			continue
		}
		name, ok := nameById[*e.SubjectID]
		if !ok {
			continue
		}
		weekType := e.WeekType
		if weekType == "" {
			weekType = "all"
		}
		entry := TimetableTeacherEntryView{
			ClassID: e.ClassID, Weekday: e.Weekday, PeriodIndex: e.PeriodIndex,
			WeekType: weekType, SubjectName: &name,
		}
		if cn, ok := className[e.ClassID]; ok {
			c := cn
			entry.ClassName = &c
		}
		if e.Room != "" {
			r := e.Room
			entry.Room = &r
		}
		result.Entries = append(result.Entries, entry)
	}
	return result, nil
}

// ForDisplay 大屏只读数据：科目（含配色）+ 节次 + 该班排课（科目名形式）。
// 忠实移植自 Laravel TimetableService::forDisplay：week_type 原样返回（空值补 all），
// 且**过滤掉 subject_id 为空或指向不存在科目的排课行**（前端不展示空课格）。
func (s *TimetableService) ForDisplay(classID, schoolID uint) (*TimetableBootstrap, error) {
	var subjects []models.Subject
	if err := s.db.Where("school_id = ?", schoolID).Order("sort_order ASC, id ASC").Find(&subjects).Error; err != nil {
		return nil, err
	}
	nameById := map[uint]string{}
	for _, sub := range subjects {
		nameById[sub.ID] = sub.Name
	}

	result := &TimetableBootstrap{
		Subjects: []TimetableSubjectView{},
		Periods:  []TimetablePeriodView{},
		Entries:  []TimetableEntryView{},
	}
	for _, sub := range subjects {
		result.Subjects = append(result.Subjects, TimetableSubjectView{
			Name: sub.Name, SimplifiedName: sub.SimplifiedName, Color: sub.Color,
		})
	}

	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).Order("period_index ASC").Find(&periods).Error; err != nil {
		return nil, err
	}
	for _, p := range periods {
		result.Periods = append(result.Periods, TimetablePeriodView{
			PeriodIndex: p.PeriodIndex, Name: p.Name, StartTime: p.StartTime, EndTime: p.EndTime,
		})
	}

	var entries []models.TimetableEntry
	if err := s.db.Where("class_id = ?", classID).
		Order("weekday ASC, period_index ASC").Find(&entries).Error; err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.SubjectID == nil {
			continue
		}
		name, ok := nameById[*e.SubjectID]
		if !ok {
			continue
		}
		weekType := e.WeekType
		if weekType == "" {
			weekType = "all"
		}
		var teacher, room *string
		if e.TeacherName != "" {
			t := e.TeacherName
			teacher = &t
		}
		if e.Room != "" {
			r := e.Room
			room = &r
		}
		result.Entries = append(result.Entries, TimetableEntryView{
			Weekday: e.Weekday, PeriodIndex: e.PeriodIndex, WeekType: weekType,
			SubjectName: &name, TeacherName: teacher, Room: room,
		})
	}

	return result, nil
}

// ToCses 导出为 CSES YAML（ClassIsland 可直接「从 CSES 导入」）。
func (s *TimetableService) ToCses(classID, schoolID uint) (string, error) {
	var subjects []models.Subject
	if err := s.db.Where("school_id = ?", schoolID).Order("sort_order ASC, id ASC").Find(&subjects).Error; err != nil {
		return "", err
	}
	nameById := map[uint]string{}
	for _, sub := range subjects {
		nameById[sub.ID] = sub.Name
	}

	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).Find(&periods).Error; err != nil {
		return "", err
	}
	periodByIndex := map[int]models.ClassPeriod{}
	for _, p := range periods {
		periodByIndex[p.PeriodIndex] = p
	}

	var entries []models.TimetableEntry
	if err := s.db.Where("class_id = ?", classID).Order("weekday ASC, period_index ASC").Find(&entries).Error; err != nil {
		return "", err
	}

	lines := []string{"version: 1"}
	lines = append(lines, "subjects:")
	if len(subjects) == 0 {
		lines = append(lines, "  []")
	}
	for _, sub := range subjects {
		lines = append(lines, "  - name: "+yamlString(sub.Name))
		if sub.SimplifiedName != "" {
			lines = append(lines, "    simplified_name: "+yamlString(sub.SimplifiedName))
		}
	}

	// 按 星期 → 单双周 分组（Go map 无序，遍历前显式排序键）
	grouped := map[int]map[string][]models.TimetableEntry{}
	for _, e := range entries {
		if e.SubjectID == nil {
			continue
		}
		weekType := e.WeekType
		if weekType == "" {
			weekType = "all"
		}
		if grouped[e.Weekday] == nil {
			grouped[e.Weekday] = map[string][]models.TimetableEntry{}
		}
		grouped[e.Weekday][weekType] = append(grouped[e.Weekday][weekType], e)
	}
	weekdays := make([]int, 0, len(grouped))
	for d := range grouped {
		weekdays = append(weekdays, d)
	}
	for i := 0; i < len(weekdays); i++ {
		for j := i + 1; j < len(weekdays); j++ {
			if weekdays[j] < weekdays[i] {
				weekdays[i], weekdays[j] = weekdays[j], weekdays[i]
			}
		}
	}

	lines = append(lines, "schedules:")
	if len(grouped) == 0 {
		lines = append(lines, "  []")
	}
	for _, weekday := range weekdays {
		byWeek := grouped[weekday]
		for _, weekType := range []string{"all", "odd", "even"} {
			if len(byWeek[weekType]) == 0 {
				continue
			}
			label := weekdayLabels[weekday]
			if label == "" {
				label = "第" + itoa(weekday) + "天"
			}
			if weekType == "odd" {
				label += "·单周"
			} else if weekType == "even" {
				label += "·双周"
			}
			lines = append(lines, "  - name: "+yamlString(label))
			lines = append(lines, "    enable_day: "+itoa(weekday))
			lines = append(lines, "    weeks: "+weekType)
			lines = append(lines, "    classes:")
			for _, e := range byWeek[weekType] {
				if e.SubjectID == nil {
					continue
				}
				subjectName, ok := nameById[*e.SubjectID]
				if !ok {
					continue
				}
				lines = append(lines, "      - subject: "+yamlString(subjectName))
				if p, ok := periodByIndex[e.PeriodIndex]; ok {
					lines = append(lines, "        start_time: \""+p.StartTime+"\"")
					lines = append(lines, "        end_time: \""+p.EndTime+"\"")
				}
				if e.TeacherName != "" {
					lines = append(lines, "        teacher: "+yamlString(e.TeacherName))
				}
				if e.Room != "" {
					lines = append(lines, "        room: "+yamlString(e.Room))
				}
			}
		}
	}

	return strings.Join(lines, "\n") + "\n", nil
}

// yamlString YAML 双引号标量（转义反斜杠与双引号）。
func yamlString(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	return "\"" + replacer.Replace(value) + "\""
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}

func keysOf(m map[uint]bool) []uint {
	ks := make([]uint, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func keysOfClass(m map[uint]string) []uint {
	ks := make([]uint, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func isoTime(t time.Time) string { return t.Format(time.RFC3339) }

func isoTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
