// 课表模型：科目 / 节次（学校级共享）+ 排课格子 + 课表修改申请。
// 忠实移植自 Laravel 端 subjects / class_periods / timetable_entries /
// timetable_change_requests 四张表及其业务口径。
package models

import "time"

// Subject 科目（学校级，全校共享一套，同名同色）。
type Subject struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SchoolID       uint      `gorm:"index;not null;uniqueIndex:idx_subject_school_name" json:"school_id"`
	Name           string    `gorm:"size:50;not null;uniqueIndex:idx_subject_school_name" json:"name"`
	SimplifiedName string    `gorm:"size:20" json:"simplified_name"`
	Color          string    `gorm:"size:20" json:"color"`
	SortOrder      int       `gorm:"default:0" json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ClassPeriod 节次作息（学校级，全校共享一套）。
type ClassPeriod struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"index;not null;uniqueIndex:idx_class_period_school_index" json:"school_id"`
	PeriodIndex int       `gorm:"not null;uniqueIndex:idx_class_period_school_index" json:"period_index"`
	Name        string    `gorm:"size:30;not null" json:"name"`
	StartTime   string    `gorm:"size:8;not null" json:"start_time"`
	EndTime     string    `gorm:"size:8;not null" json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TimetableEntry 排课格子：某班级 · 星期 · 第几节 · 单双周 → 科目。
type TimetableEntry struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClassID     uint      `gorm:"index;not null;uniqueIndex:idx_timetable_slot" json:"class_id"`
	Weekday     int       `gorm:"not null;uniqueIndex:idx_timetable_slot" json:"weekday"`
	PeriodIndex int       `gorm:"not null;uniqueIndex:idx_timetable_slot" json:"period_index"`
	WeekType    string    `gorm:"size:10;default:all;uniqueIndex:idx_timetable_slot" json:"week_type"`
	SubjectID   *uint     `gorm:"index" json:"subject_id"`
	TeacherName string    `gorm:"size:50" json:"teacher_name"`
	Room        string    `gorm:"size:50" json:"room"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 课表修改申请状态。
const (
	TimetableChangePending  = "pending"
	TimetableChangeApproved = "approved"
	TimetableChangeRejected = "rejected"
)

// TimetableChangeRequest 课表修改申请（教师提交 → 管理员审核通过后生效）。
//
// Payload 存提交时的完整课表快照（JSON 文本：subjects / periods / entries），
// 批准时一次性应用，不产生部分生效。
type TimetableChangeRequest struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SchoolID    uint       `gorm:"index;not null" json:"school_id"`
	ClassID     uint       `gorm:"index;not null" json:"class_id"`
	RequestedBy uint       `gorm:"index;not null" json:"requested_by"`
	Payload     string     `gorm:"type:text" json:"payload"`
	EntryCount  int        `gorm:"default:0" json:"entry_count"`
	Status      string     `gorm:"size:20;default:pending;index" json:"status"`
	ReviewedBy  *uint      `json:"reviewed_by"`
	ReviewNote  *string    `gorm:"size:200" json:"review_note"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TimetableTeacherAssignment 任课表：班级 × 科目 → 教师（全校智能排课的依据）。
// 忠实移植自 Laravel timetable_teacher_assignments（唯一键 class_id + subject_name）。
type TimetableTeacherAssignment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"index;not null" json:"school_id"`
	ClassID     uint      `gorm:"index;not null;uniqueIndex:idx_tt_assignment" json:"class_id"`
	SubjectName string    `gorm:"size:50;not null;uniqueIndex:idx_tt_assignment" json:"subject_name"`
	TeacherName string    `gorm:"size:50;not null" json:"teacher_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TimetableTeacherUnavailability 教师不可用时段（学校级，1-7 星期 × 节次）。
// 忠实移植自 Laravel timetable_teacher_unavailabilities
// （唯一键 school_id + teacher_name + weekday + period_index）。
type TimetableTeacherUnavailability struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"index;not null;uniqueIndex:idx_tt_unavail" json:"school_id"`
	TeacherName string    `gorm:"size:50;not null;uniqueIndex:idx_tt_unavail" json:"teacher_name"`
	Weekday     int       `gorm:"not null;uniqueIndex:idx_tt_unavail" json:"weekday"`
	PeriodIndex int       `gorm:"not null;uniqueIndex:idx_tt_unavail" json:"period_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
