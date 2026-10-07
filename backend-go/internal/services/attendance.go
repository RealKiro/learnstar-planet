// 考勤服务：今日明细 / 开始点名 / 手动设状态 / 请假 / 缺勤 / 统计。
// 忠实移植自 Laravel App\Services\AttendanceService；企业微信请假同步为
// startForClass 末尾的补充步骤（WechatWorkAttendanceService.ApplyUnsyncedLeavesForClass）。
// date 以 "2006-01-02" 形式存储便于跨库比较。
package services

import (
	"errors"
	"math"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// AttendanceService 考勤读写服务。
type AttendanceService struct {
	db     *gorm.DB
	scope  *Scope
	wechat *WechatWorkAttendanceService
}

// NewAttendanceService 创建考勤服务。
//
// 可选参数 wechat 用于注入企微请假同步服务（HTTP 层复用同一个实例，避免重复读环境变量）；
// 不传时自建一个（服务层既有测试与默认路径无需感知该依赖）。
func NewAttendanceService(db *gorm.DB, scope *Scope, wechat ...*WechatWorkAttendanceService) *AttendanceService {
	s := &AttendanceService{db: db, scope: scope}
	if len(wechat) > 0 && wechat[0] != nil {
		s.wechat = wechat[0]
	} else {
		s.wechat = NewWechatWorkAttendanceService(db, nil)
	}
	return s
}

// AttendanceLeaveRecord 今日考勤条目附带的企微请假记录（同 Laravel `with('leaveRecord:id,sp_no,leave_type,reason')`）。
type AttendanceLeaveRecord struct {
	SpNo      string `json:"sp_no"`
	LeaveType string `json:"leave_type"`
	Reason    string `json:"reason"`
}

// AttendanceTodayItem 今日考勤明细条目（含学生信息与请假记录）。
// 字段与顺序逐字同 Laravel `AttendanceService::today` 的元素：
// `check_in_time` 是 `sign_in_at?->toDateTimeString()`（`Y-m-d H:i:s` 字符串，空值为 null），
// `leave_record` 无请假记录时为 null。
type AttendanceTodayItem struct {
	ID          uint                   `json:"id"`
	StudentID   uint                   `json:"student_id"`
	StudentName string                 `json:"student_name"`
	StudentNo   string                 `json:"student_no"`
	Status      string                 `json:"status"`
	Source      string                 `json:"source"`
	Remark      string                 `json:"remark"`
	CheckInTime *string                `json:"check_in_time"`
	LeaveRecord *AttendanceLeaveRecord `json:"leave_record"`
}

// Today 返回教师管辖班级今日考勤明细（按创建时间倒序，含学生信息与企微请假记录）。
func (s *AttendanceService) Today(u *models.User) ([]AttendanceTodayItem, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return []AttendanceTodayItem{}, nil
	}

	var records []models.Attendance
	if err := s.db.Where("class_id IN ? AND date = ?", classIDs, util.Today()).
		Order("created_at DESC, id DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []AttendanceTodayItem{}, nil
	}

	studentIDs := make([]uint, 0, len(records))
	leaveIDs := make([]uint, 0, len(records))
	seenLeave := map[uint]bool{}
	for _, r := range records {
		studentIDs = append(studentIDs, r.StudentID)
		if r.LeaveRecordID != nil && !seenLeave[*r.LeaveRecordID] {
			seenLeave[*r.LeaveRecordID] = true
			leaveIDs = append(leaveIDs, *r.LeaveRecordID)
		}
	}
	var students []models.Student
	if err := s.db.Where("id IN ?", studentIDs).Find(&students).Error; err != nil {
		return nil, err
	}
	studentMap := make(map[uint]models.Student, len(students))
	for _, st := range students {
		studentMap[st.ID] = st
	}

	// 请假记录：Laravel 是 Eloquent 关联（`with('leaveRecord:id,sp_no,leave_type,reason')`），
	// 指向不存在的行时关联为 null，故此处只按 id 取需要的三列，缺失即视为 null。
	leaveMap := map[uint]AttendanceLeaveRecord{}
	if len(leaveIDs) > 0 {
		var leaves []models.WechatWorkLeaveRecord
		if err := s.db.Select("id, sp_no, leave_type, reason").
			Where("id IN ?", leaveIDs).Find(&leaves).Error; err != nil {
			return nil, err
		}
		for _, lv := range leaves {
			leaveMap[lv.ID] = AttendanceLeaveRecord{
				SpNo:      lv.SpNo,
				LeaveType: lv.LeaveType,
				Reason:    lv.Reason,
			}
		}
	}

	out := make([]AttendanceTodayItem, 0, len(records))
	for _, r := range records {
		item := AttendanceTodayItem{
			ID:          r.ID,
			StudentID:   r.StudentID,
			Status:      r.Status,
			Source:      r.Source,
			Remark:      r.Remark,
			CheckInTime: toDateTimeString(r.SignInAt),
		}
		if st, ok := studentMap[r.StudentID]; ok {
			item.StudentName = st.Name
			item.StudentNo = st.StudentNo
		}
		if r.LeaveRecordID != nil {
			if lv, ok := leaveMap[*r.LeaveRecordID]; ok {
				leave := lv
				item.LeaveRecord = &leave
			}
		}
		out = append(out, item)
	}
	return out, nil
}

// toDateTimeString 把可空时间格式化为 Laravel `Carbon::toDateTimeString()` 的
// `Y-m-d H:i:s`（业务时区 Asia/Shanghai）；nil → null。
// 注：Go 端此前直接序列化 `time.Time`（RFC3339），与 Laravel 的字符串形态不一致。
func toDateTimeString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.In(util.Loc).Format("2006-01-02 15:04:05")
	return &s
}

// StartResult 开始点名结果。
type StartResult struct {
	Total            int `json:"total"`
	WechatLeaveCount int `json:"wechat_leave_count"`
}

// Start 为教师所有可访问班级的活跃学生创建/重置今日考勤记录（默认到课），
// 并把当日未同步的企微请假补写到考勤；wechat_leave_count 统计这些
// `source = wechat_work` 的今日记录（同 Laravel `AttendanceService::start`）。
func (s *AttendanceService) Start(u *models.User) (StartResult, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return StartResult{}, err
	}

	date := util.Today()
	count := 0
	for _, classID := range classIDs {
		n, err := s.startForClass(classID, u.ID, date)
		if err != nil {
			return StartResult{}, err
		}
		count += n
	}

	wechatLeave := 0
	if len(classIDs) > 0 {
		var n int64
		if err := s.db.Model(&models.Attendance{}).
			Where("class_id IN ? AND date = ? AND source = ?", classIDs, date, "wechat_work").
			Count(&n).Error; err != nil {
			return StartResult{}, err
		}
		wechatLeave = int(n)
	}

	return StartResult{Total: count, WechatLeaveCount: wechatLeave}, nil
}

// startForClass 为该班级所有活跃学生建立今日考勤记录（默认到课）。
// 已存在且非手动标记的自动记录会被重置为 present/auto。
// 末尾接上 Laravel `startAttendanceForClass` 的后半段：把当日仍未同步的企微请假
// （approve_status=approved 且 synced_at IS NULL）落到考勤上。
func (s *AttendanceService) startForClass(classID, teacherID uint, date string) (int, error) {
	var students []models.Student
	if err := s.db.Where("class_id = ? AND status = ?", classID, "active").Find(&students).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, stu := range students {
		var rec models.Attendance
		err := s.db.Where("class_id = ? AND student_id = ? AND date = ?", classID, stu.ID, date).
			First(&rec).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			rec = models.Attendance{
				ClassID:   classID,
				StudentID: stu.ID,
				TeacherID: &teacherID,
				Date:      date,
				Status:    "present",
				Source:    "auto",
				Remark:    "",
			}
			if err := s.db.Create(&rec).Error; err != nil {
				return 0, err
			}
		case err != nil:
			return 0, err
		default:
			if rec.Source != "manual" {
				if err := s.db.Model(&rec).Updates(map[string]any{
					"status":          "present",
					"source":          "auto",
					"remark":          "",
					"leave_record_id": nil,
				}).Error; err != nil {
					return 0, err
				}
			}
		}
		count++
	}
	if _, err := s.wechat.ApplyUnsyncedLeavesForClass(classID, date); err != nil {
		return 0, err
	}
	return count, nil
}

// SetStatus 手动设置某学生今日考勤状态（source 变为 manual）。
func (s *AttendanceService) SetStatus(u *models.User, studentID uint, status string, remark *string) (*models.Attendance, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var rec models.Attendance
	err = s.db.Where("class_id IN ? AND student_id = ? AND date = ?", classIDs, studentID, util.Today()).
		First(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("今日考勤记录不存在")
	}
	if err != nil {
		return nil, err
	}

	remarkVal := ""
	if remark != nil {
		remarkVal = *remark
	}

	updates := map[string]any{
		"status": status,
		"source": "manual",
		"remark": remarkVal,
	}
	if status == "present" {
		now := util.Now()
		updates["sign_in_at"] = now
	}
	if err := s.db.Model(&rec).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&rec, rec.ID).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

// MarkLeave 手动标记请假（自动建缺省到课记录后改为 leave/manual）。
func (s *AttendanceService) MarkLeave(u *models.User, studentID uint, remark string) (*models.Attendance, error) {
	student, err := s.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}
	return s.markManual(student, "leave", remark)
}

// MarkAbsent 手动标记缺勤（未填 remark 时用默认文案）。
func (s *AttendanceService) MarkAbsent(u *models.User, studentID uint, remark *string) (*models.Attendance, error) {
	student, err := s.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}
	remarkVal := "未联系到家长，建议后续跟进"
	if remark != nil && *remark != "" {
		remarkVal = *remark
	}
	return s.markManual(student, "absent", remarkVal)
}

// markManual 用 firstOrCreate 建立今日记录（默认到课）后改为指定手动状态。
func (s *AttendanceService) markManual(student *models.Student, status, remark string) (*models.Attendance, error) {
	date := util.Today()
	var rec models.Attendance
	err := s.db.Where("class_id = ? AND student_id = ? AND date = ?", student.ClassID, student.ID, date).
		First(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rec = models.Attendance{
			ClassID:   student.ClassID,
			StudentID: student.ID,
			Date:      date,
			Status:    "present",
			Source:    "auto",
			Remark:    "",
		}
		if err := s.db.Create(&rec).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// 更新为手动状态（markManualLeave 的 Laravel 实现会清空 leave_record_id）。
	if err := s.db.Model(&rec).Updates(map[string]any{
		"status":          status,
		"source":          "manual",
		"remark":          remark,
		"leave_record_id": nil,
	}).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&rec, rec.ID).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

// AttendanceSummary 今日考勤四态统计。
type AttendanceSummary struct {
	Present          int     `json:"present"`
	Late             int     `json:"late"`
	Leave            int     `json:"leave"`
	Absent           int     `json:"absent"`
	Rate             float64 `json:"rate"`
	WechatLeaveCount int     `json:"wechat_leave_count"`
	ManualLeaveCount int     `json:"manual_leave_count"`
}

// Summary 今日考勤四态统计；wechat_leave_count / manual_leave_count 分别统计
// `status = leave` 且来源为 wechat_work / manual 的记录（同 Laravel）。
func (s *AttendanceService) Summary(u *models.User) (AttendanceSummary, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return AttendanceSummary{}, err
	}

	var records []models.Attendance
	if len(classIDs) > 0 {
		if err := s.db.Where("class_id IN ? AND date = ?", classIDs, util.Today()).
			Find(&records).Error; err != nil {
			return AttendanceSummary{}, err
		}
	}

	present, late, leave, absent := 0, 0, 0, 0
	manualLeave := 0
	wechatLeave := 0
	for _, r := range records {
		switch r.Status {
		case "present":
			present++
		case "late":
			late++
		case "leave":
			leave++
			if r.Source == "manual" {
				manualLeave++
			}
			if r.Source == "wechat_work" {
				wechatLeave++
			}
		case "absent":
			absent++
		}
	}

	total := len(records)
	if total < 1 {
		total = 1
	}
	rate := math.Round(float64(present)/float64(total)*100*10) / 10

	return AttendanceSummary{
		Present:          present,
		Late:             late,
		Leave:            leave,
		Absent:           absent,
		Rate:             rate,
		WechatLeaveCount: wechatLeave,
		ManualLeaveCount: manualLeave,
	}, nil
}
