// 今日考勤列表的 leave_record / check_in_time 字段单测。
//
// 权威来源：Laravel App\Services\AttendanceService::today（第 27-48 行）——
// `with('leaveRecord:id,sp_no,leave_type,reason')` 与 `sign_in_at?->toDateTimeString()`。
// 全部使用内存 SQLite，**不访问外网**。
package services_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAttendanceTodayIncludesLeaveRecord 有 leave_record_id → 三字段取自请假记录。
func TestAttendanceTodayIncludesLeaveRecord(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	leave := models.WechatWorkLeaveRecord{
		SchoolID:              school.ID,
		ClassID:               &class.ID,
		StudentID:             &student.ID,
		ParentWeworkUserid:    "wework-user-1",
		StudentNameFromWework: "小明",
		SpNo:                  "SP2026092001",
		LeaveStartDate:        util.Today(),
		LeaveEndDate:          util.Today(),
		LeaveType:             "病假",
		Reason:                "发烧在家休息",
		ApproveStatus:         "approved",
	}
	require.NoError(t, db.Create(&leave).Error)

	signIn := time.Date(2026, 9, 20, 8, 30, 15, 0, util.Loc)
	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: student.ID, Date: util.Today(),
		Status: "leave", Source: "wechat_work", Remark: "企微请假同步", LeaveRecordID: &leave.ID,
		SignInAt: &signIn,
	}).Error)

	svc := services.NewAttendanceService(db, services.NewScope(db))
	items, err := svc.Today(&teacher)
	require.NoError(t, err)
	require.Len(t, items, 1)

	item := items[0]
	require.NotNil(t, item.LeaveRecord, "有 leave_record_id 时必须附带请假记录")
	assert.Equal(t, "SP2026092001", item.LeaveRecord.SpNo)
	assert.Equal(t, "病假", item.LeaveRecord.LeaveType)
	assert.Equal(t, "发烧在家休息", item.LeaveRecord.Reason)

	require.NotNil(t, item.CheckInTime)
	assert.Equal(t, "2026-09-20 08:30:15", *item.CheckInTime, "check_in_time 是 Laravel toDateTimeString 形态")
	assert.Equal(t, student.ID, item.StudentID)
	assert.Equal(t, "小明", item.StudentName)

	// JSON 形状：leave_record 是对象、check_in_time 是字符串（前端 types/index.ts 即如此声明）。
	raw, err := json.Marshal(item)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))

	leaveJSON, ok := decoded["leave_record"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "SP2026092001", leaveJSON["sp_no"])
	assert.Equal(t, "病假", leaveJSON["leave_type"])
	assert.Equal(t, "发烧在家休息", leaveJSON["reason"])
	assert.Len(t, leaveJSON, 3, "只暴露 sp_no / leave_type / reason 三个字段")
	assert.Equal(t, "2026-09-20 08:30:15", decoded["check_in_time"])
}

// TestAttendanceTodayLeaveRecordNull 无 leave_record_id → leave_record 为 null（键仍在）；
// 同时覆盖「sign_in_at 为空 → check_in_time 为 null」。
func TestAttendanceTodayLeaveRecordNull(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: student.ID, Date: util.Today(),
		Status: "present", Source: "auto", Remark: "",
	}).Error)

	svc := services.NewAttendanceService(db, services.NewScope(db))
	items, err := svc.Today(&teacher)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Nil(t, items[0].LeaveRecord)
	assert.Nil(t, items[0].CheckInTime)

	raw, err := json.Marshal(items[0])
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))
	assert.Contains(t, decoded, "leave_record", "键必须存在（Laravel 恒输出该键）")
	assert.Nil(t, decoded["leave_record"])
	assert.Contains(t, decoded, "check_in_time")
	assert.Nil(t, decoded["check_in_time"])
}

// TestAttendanceTodayMissingLeaveRecordFallsBackToNull 指向已删除的请假行 → leave_record 为 null（同 Eloquent 关联）。
func TestAttendanceTodayMissingLeaveRecordFallsBackToNull(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	dangling := uint(999999)
	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: student.ID, Date: util.Today(),
		Status: "leave", Source: "manual", Remark: "事假", LeaveRecordID: &dangling,
	}).Error)

	svc := services.NewAttendanceService(db, services.NewScope(db))
	items, err := svc.Today(&teacher)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Nil(t, items[0].LeaveRecord)

	// 其它字段不受影响（回归：加字段不删字段）。
	assert.Equal(t, student.ID, items[0].StudentID)
	assert.Equal(t, "leave", items[0].Status)
	assert.Equal(t, "manual", items[0].Source)
	assert.Equal(t, "事假", items[0].Remark)
}

// TestAttendanceTodayLeaveRecordOnlyForOwningRecord 请假记录只挂到引用它的那条考勤上。
func TestAttendanceTodayLeaveRecordOnlyForOwningRecord(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, s1 := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	s2 := seedExtraStudent(t, db, class.ID)

	leave := models.WechatWorkLeaveRecord{
		SchoolID: school.ID, ClassID: &class.ID, StudentID: &s1.ID,
		ParentWeworkUserid: "wework-user-1", SpNo: "SP2026092002",
		LeaveStartDate: util.Today(), LeaveEndDate: util.Today(),
		LeaveType: "事假", Reason: "家中有事", ApproveStatus: "approved",
	}
	require.NoError(t, db.Create(&leave).Error)

	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: s1.ID, Date: util.Today(),
		Status: "leave", Source: "wechat_work", LeaveRecordID: &leave.ID,
	}).Error)
	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: s2.ID, Date: util.Today(),
		Status: "present", Source: "auto",
	}).Error)

	svc := services.NewAttendanceService(db, services.NewScope(db))
	items, err := svc.Today(&teacher)
	require.NoError(t, err)
	require.Len(t, items, 2)

	byStudent := map[uint]services.AttendanceTodayItem{}
	for _, item := range items {
		byStudent[item.StudentID] = item
	}
	require.NotNil(t, byStudent[s1.ID].LeaveRecord)
	assert.Equal(t, "SP2026092002", byStudent[s1.ID].LeaveRecord.SpNo)
	assert.Nil(t, byStudent[s2.ID].LeaveRecord)
}
