// 企业微信请假记录（企微审批 → 考勤请假的中转表）。
//
// 逐字段移植自 Laravel 迁移 2026_07_13_000002_create_wechat_work_leave_records.php 与
// App\Models\WechatWorkLeaveRecord（含 getUnsyncedForDate 的过滤口径）。
//
// 有意差异：`leave_start_date` / `leave_end_date` 在 Laravel 是 `date` 列并 cast 为 Carbon，
// 过滤用 `whereDate(...)`；Go 端按本仓既有约定（Attendance.Date）存 `YYYY-MM-DD` 字符串、
// 直接做字符串比较——三种数据库（SQLite/MySQL/PostgreSQL）下排序与比较语义都正确。
package models

import (
	"time"

	"gorm.io/gorm"
)

// WechatWorkLeaveRecord 一条企微请假（审批）记录。
type WechatWorkLeaveRecord struct {
	ID       uint  `gorm:"primaryKey" json:"id"`
	SchoolID uint  `gorm:"index;not null" json:"school_id"`
	ClassID  *uint `json:"class_id"`
	// StudentID 匹配到的学生；未匹配到为 NULL（Laravel 允许为空，后续 startAttendanceForClass 会跳过）。
	StudentID             *uint      `gorm:"index" json:"student_id"`
	ParentWeworkUserid    string     `gorm:"size:128;not null" json:"parent_wework_userid"`
	StudentNameFromWework string     `gorm:"size:100" json:"student_name_from_wework"`
	SpNo                  string     `gorm:"size:128;uniqueIndex;not null" json:"sp_no"`
	LeaveStartDate        string     `gorm:"size:10;index:idx_wwlr_school_start,priority:2;index:idx_wwlr_student_start,priority:2;not null" json:"leave_start_date"`
	LeaveEndDate          string     `gorm:"size:10;not null" json:"leave_end_date"`
	LeaveType             string     `gorm:"size:50" json:"leave_type"`
	Reason                string     `gorm:"type:text" json:"reason"`
	ApproveStatus         string     `gorm:"size:20;default:approved" json:"approve_status"`
	ApprovedAt            *time.Time `json:"approved_at"`
	// RawData 详情 JSON 文本（Laravel `json` 列 + array cast；Go 端直接存 parse 结果的 JSON 文本）。
	RawData string `gorm:"type:text" json:"raw_data"`
	// SyncedAt 已写入考勤的时间；NULL 表示尚未同步（startAttendanceForClass 的后半段会补写）。
	SyncedAt  *time.Time `gorm:"index" json:"synced_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName 表名（与 Laravel 迁移一致）。
func (WechatWorkLeaveRecord) TableName() string { return "wechat_work_leave_records" }

// GetUnsyncedLeaveRecordsForDate 返回指定日期仍未同步到考勤的「已通过」请假记录。
//
// 同 Laravel `WechatWorkLeaveRecord::getUnsyncedForDate($date)`：
// `approve_status = 'approved' AND synced_at IS NULL`
// `AND leave_start_date <= date AND leave_end_date >= date`。
func GetUnsyncedLeaveRecordsForDate(db *gorm.DB, date string) ([]WechatWorkLeaveRecord, error) {
	var rows []WechatWorkLeaveRecord
	if err := db.Where(
		"approve_status = ? AND synced_at IS NULL AND leave_start_date <= ? AND leave_end_date >= ?",
		"approved", date, date,
	).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
