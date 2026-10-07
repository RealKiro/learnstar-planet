// 企业微信请假同步服务：审批回调 / 定时同步 → 考勤请假。
//
// 忠实移植自 Laravel App\Services\WechatWorkAttendanceService（handleWebhookCallback /
// syncForSchool / syncAll / applyLeave / startAttendanceForClass / matchStudent），
// commit 的判定、日期回退、跳过手动修改等口径逐条对齐。
//
// 与 Laravel 的有意差异（逐条见 README「企微回调与请假同步」小节）：
//  1. `matchStudent` 按**规格语义**实现三段优先级（精确 → LIKE → 第一条）。Laravel 原实现里
//     `$q` 是同一个 Builder，`$q->where('name', $hint)` 会**留在**后续 `LIKE` 查询上，导致
//     LIKE 与「第一条」两段分支实际不可达（等价于「只认精确匹配」）。此处按业务意图实现。
//  2. 日志走标准库 `log`（stdout），Laravel 走 `Log::error/info`（storage/logs）。
//  3. `students.parent_id` 在 Laravel 已于 2026_08_06_000006 迁移清空、家长角色被移除；
//     Go 端保留该列语义，但同样只有外部写入家长数据时才会命中匹配。
package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// WechatWorkLeaveSyncResult 同步结果（同 Laravel syncForSchool / syncAll 返回的数组，
// 多带一个 `failed` 键 —— 仅 syncAll 会累加，syncForSchool 恒为 0）。
type WechatWorkLeaveSyncResult struct {
	Synced  int `json:"synced"`
	Skipped int `json:"skipped"`
	Failed  int `json:"failed"`
}

// WechatWorkAttendanceService 企微请假 → 考勤的同步服务。
type WechatWorkAttendanceService struct {
	db     *gorm.DB
	wework *WechatWorkService
}

// NewWechatWorkAttendanceService 创建同步服务；wework 为 nil 时按环境变量自建一个。
func NewWechatWorkAttendanceService(db *gorm.DB, wework *WechatWorkService) *WechatWorkAttendanceService {
	if wework == nil {
		wework = NewWechatWorkService(db)
	}
	return &WechatWorkAttendanceService{db: db, wework: wework}
}

// WechatWork 返回内部企微调用服务（供 HTTP 层读取 token / encoding_aes_key 配置）。
func (s *WechatWorkAttendanceService) WechatWork() *WechatWorkService { return s.wework }

// SyncAll 遍历全部 `status = active` 的学校做同步（单校失败计入 failed 并继续，同 Laravel syncAll）。
// 学校列表查询失败直接返回错误（Laravel 该查询在 try/catch 之外，同样会向上抛）。
func (s *WechatWorkAttendanceService) SyncAll(date string) (WechatWorkLeaveSyncResult, error) {
	if date == "" {
		date = util.Today()
	}

	var schools []models.School
	if err := s.db.Where("status = ?", "active").Find(&schools).Error; err != nil {
		return WechatWorkLeaveSyncResult{}, err
	}

	result := WechatWorkLeaveSyncResult{}
	for _, school := range schools {
		one, err := s.SyncForSchool(school.ID, date)
		if err != nil {
			log.Printf("企微请假同步失败 学校id=%d 日期=%s：%v", school.ID, date, err)
			result.Failed++
			continue
		}
		result.Synced += one.Synced
		result.Skipped += one.Skipped
	}
	return result, nil
}

// SyncForSchool 按天内区间拉取该校「已通过」的请假审批并落库（幂等：同 sp_no 已存在即跳过）。
// 企微未配置（corp_id 为空）时原样返回零值（同 Laravel）。
func (s *WechatWorkAttendanceService) SyncForSchool(schoolID uint, date string) (WechatWorkLeaveSyncResult, error) {
	result := WechatWorkLeaveSyncResult{}
	if s.wework.CorpID == "" {
		return result, nil
	}

	dayStart, dayEnd, err := wechatWorkDayBounds(date)
	if err != nil {
		return result, err
	}

	spNos, err := s.wework.GetLeaveApprovalSpNos(schoolID, dayStart, dayEnd)
	if err != nil {
		return result, err
	}

	for _, spNo := range spNos {
		exists, err := s.leaveRecordExists(spNo)
		if err != nil {
			return result, err
		}
		if exists {
			result.Skipped++
			continue
		}

		detail, err := s.wework.GetApprovalDetail(schoolID, spNo)
		if err != nil {
			return result, err
		}
		if detail == nil {
			continue
		}

		student, err := s.MatchStudent(detail.SubmitterUserID, detail.StudentName)
		if err != nil {
			return result, err
		}

		rec := models.WechatWorkLeaveRecord{
			SchoolID:              schoolID,
			ClassID:               classIDOf(student),
			StudentID:             studentIDOf(student),
			ParentWeworkUserid:    detail.SubmitterUserID,
			StudentNameFromWework: stringOrEmpty(detail.StudentName),
			SpNo:                  detail.SpNo,
			LeaveStartDate:        leaveBoundary(detail.LeaveStart, date),
			LeaveEndDate:          leaveBoundary(detail.LeaveEnd, date),
			LeaveType:             stringOrEmpty(detail.LeaveType),
			Reason:                stringOrEmpty(detail.Reason),
			ApproveStatus:         approveStatusOf(detail.ApproveStatus),
			ApprovedAt:            approvedAtOf(detail.ApprovedAt),
			RawData:               rawDataOf(detail),
		}
		if err := s.db.Create(&rec).Error; err != nil {
			return result, err
		}
		if student != nil {
			if err := s.ApplyLeave(student, &rec, date); err != nil {
				return result, err
			}
		}
		result.Synced++
	}

	return result, nil
}

// HandleWebhookCallback 处理企微审批状态变更回调（同 Laravel handleWebhookCallback）：
// 只处理「已通过」（sp_status == 2）且非重复的审批单；详情取不到就放弃；匹配到学生才写考勤。
// 记录一律落库（未匹配到学生时 student_id / class_id 为 NULL），供后续点名时补同步。
func (s *WechatWorkAttendanceService) HandleWebhookCallback(schoolID uint, spNo string, spStatus int) error {
	if spNo == "" || spStatus != 2 {
		return nil
	}

	exists, err := s.leaveRecordExists(spNo)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	detail, err := s.wework.GetApprovalDetail(schoolID, spNo)
	if err != nil {
		return err
	}
	if detail == nil {
		return nil
	}

	student, err := s.MatchStudent(detail.SubmitterUserID, detail.StudentName)
	if err != nil {
		return err
	}

	today := util.Today()
	now := util.Now()
	rec := models.WechatWorkLeaveRecord{
		SchoolID:              schoolID,
		ClassID:               classIDOf(student),
		StudentID:             studentIDOf(student),
		ParentWeworkUserid:    detail.SubmitterUserID,
		StudentNameFromWework: stringOrEmpty(detail.StudentName),
		SpNo:                  detail.SpNo,
		LeaveStartDate:        leaveBoundary(detail.LeaveStart, today),
		LeaveEndDate:          leaveBoundary(detail.LeaveEnd, today),
		LeaveType:             stringOrEmpty(detail.LeaveType),
		Reason:                stringOrEmpty(detail.Reason),
		ApproveStatus:         "approved",
		ApprovedAt:            &now,
		RawData:               rawDataOf(detail),
	}
	if err := s.db.Create(&rec).Error; err != nil {
		return err
	}

	if student != nil {
		return s.ApplyLeave(student, &rec, today)
	}
	return nil
}

// ApplyLeave 把一条请假记录落到该生当天的考勤上（同 Laravel applyLeave）：
// 无记录 → 新建 leave/wechat_work；已存在且 source=manual → **不覆盖**；否则更新为 leave/wechat_work。
// 两条写入路径都会回写记录的 synced_at。
func (s *WechatWorkAttendanceService) ApplyLeave(student *models.Student, rec *models.WechatWorkLeaveRecord, date string) error {
	var existing models.Attendance
	err := s.db.Where("student_id = ? AND date = ?", student.ID, date).First(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		teacherID, err := s.classTeacherID(student.ClassID)
		if err != nil {
			return err
		}
		created := models.Attendance{
			ClassID:       student.ClassID,
			StudentID:     student.ID,
			TeacherID:     teacherID,
			Date:          date,
			Status:        "leave",
			Source:        "wechat_work",
			Remark:        rec.Reason,
			LeaveRecordID: &rec.ID,
		}
		if err := s.db.Create(&created).Error; err != nil {
			return err
		}
		return s.markLeaveRecordSynced(rec)
	case err != nil:
		return err
	}

	if existing.Source == "manual" {
		log.Printf("企微请假同步：跳过手动修改的考勤 考勤id=%d 学生id=%d 审批单=%s", existing.ID, student.ID, rec.SpNo)
		return nil
	}

	if err := s.db.Model(&existing).Updates(map[string]any{
		"status":          "leave",
		"source":          "wechat_work",
		"remark":          rec.Reason,
		"leave_record_id": rec.ID,
	}).Error; err != nil {
		return err
	}
	return s.markLeaveRecordSynced(rec)
}

// ApplyUnsyncedLeavesForClass 把指定日期仍未同步（synced_at IS NULL）且属于该班的请假记录补写到考勤上。
//
// 这是 Laravel `startAttendanceForClass($classId, $teacherId, $date)` 的**后半段**；
// 前半段（为全班活跃学生 firstOrCreate present/auto、非 manual 重置）由
// AttendanceService.startForClass 承担，两者拼起来与 Laravel 一致。
// 返回值是本次实际补写的条数（Laravel 该段不返回值，仅为便于测试与观测而返回）。
func (s *WechatWorkAttendanceService) ApplyUnsyncedLeavesForClass(classID uint, date string) (int, error) {
	records, err := models.GetUnsyncedLeaveRecordsForDate(s.db, date)
	if err != nil {
		return 0, err
	}

	applied := 0
	for i := range records {
		rec := records[i]
		if rec.StudentID == nil {
			continue
		}
		var student models.Student
		err := s.db.First(&student, *rec.StudentID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return applied, err
		}
		if student.ClassID != classID {
			continue
		}
		if err := s.ApplyLeave(&student, &rec, date); err != nil {
			return applied, err
		}
		applied++
	}
	return applied, nil
}

// MatchStudent 依据「家长企微 userid」找学生（同 Laravel matchStudent 的**规格语义**）：
// 先用 third_party_bindings（platform=wechat_work, platform_id=userid）找到家长用户，
// 再在该用户 parent_id 名下且 status=active 的学生中依次尝试：精确同名 → 姓名包含 → 第一条。
// 未绑定或名下无活跃学生返回 (nil, nil)。
func (s *WechatWorkAttendanceService) MatchStudent(parentWeworkUserID string, nameHint *string) (*models.Student, error) {
	var binding models.ThirdPartyBinding
	err := s.db.Where("platform = ? AND platform_id = ?", "wechat_work", parentWeworkUserID).
		First(&binding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	base := func() *gorm.DB {
		return s.db.Where("parent_id = ? AND status = ?", binding.UserID, "active")
	}

	if nameHint != nil {
		var exact models.Student
		err = base().Where("name = ?", *nameHint).Order("id ASC").First(&exact).Error
		if err == nil {
			return &exact, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var partial models.Student
		err = base().Where("name LIKE ?", "%"+*nameHint+"%").Order("id ASC").First(&partial).Error
		if err == nil {
			return &partial, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	var fallback models.Student
	err = base().Order("id ASC").First(&fallback).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fallback, nil
}

// leaveRecordExists 判断该审批单是否已落库（unique sp_no）。
func (s *WechatWorkAttendanceService) leaveRecordExists(spNo string) (bool, error) {
	var count int64
	if err := s.db.Model(&models.WechatWorkLeaveRecord{}).Where("sp_no = ?", spNo).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// markLeaveRecordSynced 回写 synced_at（同 Laravel `$rec->update(['synced_at' => now()])`）。
func (s *WechatWorkAttendanceService) markLeaveRecordSynced(rec *models.WechatWorkLeaveRecord) error {
	now := util.Now()
	if err := s.db.Model(&models.WechatWorkLeaveRecord{}).Where("id = ?", rec.ID).
		Update("synced_at", now).Error; err != nil {
		return err
	}
	rec.SyncedAt = &now
	return nil
}

// classTeacherID 取班级班主任 ID（等价 Laravel `$stu->classRoom?->teacher_id`）。
func (s *WechatWorkAttendanceService) classTeacherID(classID uint) (*uint, error) {
	var class models.ClassRoom
	err := s.db.Select("id", "teacher_id").First(&class, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return class.TeacherID, nil
}

// wechatWorkDayBounds 返回「当天 00:00:00 / 23:59:59」的 Unix 秒（业务时区，同 Laravel strtotime）。
func wechatWorkDayBounds(date string) (int64, int64, error) {
	day, err := time.ParseInLocation("2006-01-02", date, util.Loc)
	if err != nil {
		return 0, 0, err
	}
	return day.Unix(), day.Add(24*time.Hour - time.Second).Unix(), nil
}

// leaveBoundary 把 Unix 秒格式化为日期；为 0（详情缺失）时回退到 fallback（同 Laravel `?:`）。
func leaveBoundary(unixSec int64, fallback string) string {
	if unixSec <= 0 {
		return fallback
	}
	return time.Unix(unixSec, 0).In(util.Loc).Format("2006-01-02")
}

// approveStatusOf 同 Laravel `$d['approve_status'] === 2 ? 'approved' : 'pending'`。
func approveStatusOf(status int) string {
	if status == 2 {
		return "approved"
	}
	return "pending"
}

// approvedAtOf 同 Laravel `$d['approved_at'] ? date(...) : null`（0 视为无）。
func approvedAtOf(unixSec *int64) *time.Time {
	if unixSec == nil || *unixSec == 0 {
		return nil
	}
	t := time.Unix(*unixSec, 0).In(util.Loc)
	return &t
}

// rawDataOf 把详情序列化为 JSON 文本（Laravel 存的是 parse 后的关联数组）。
func rawDataOf(detail *ApprovalDetail) string {
	buf, err := json.Marshal(detail)
	if err != nil {
		return ""
	}
	return string(buf)
}

// stringOrEmpty 解引用可空字符串（Go 端空串代替 NULL，同本仓既有约定）。
func stringOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// classIDOf / studentIDOf 把「可能未匹配到学生」的两种情况统一成可空外键。
func classIDOf(student *models.Student) *uint {
	if student == nil {
		return nil
	}
	return &student.ClassID
}

func studentIDOf(student *models.Student) *uint {
	if student == nil {
		return nil
	}
	return &student.ID
}
