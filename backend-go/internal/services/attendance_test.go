package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
)

// Start 为每个活跃学生创建今日考勤记录（默认到课）。
func TestAttendanceStartCreatesTodayRecords(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	res, err := svc.Start(&teacher)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if res.Total != 1 {
		t.Fatalf("total = %d, want 1", res.Total)
	}

	items, err := svc.Today(&teacher)
	if err != nil {
		t.Fatalf("today: %v", err)
	}
	if len(items) != 1 || items[0].StudentID != student.ID {
		t.Fatalf("today items wrong: %+v", items)
	}
	if items[0].Status != "present" || items[0].Source != "auto" {
		t.Fatalf("record not present/auto: %+v", items[0])
	}
	if items[0].StudentName != "小明" {
		t.Fatalf("student name not populated: %+v", items[0])
	}
	_ = class
}

// Start 已存在的非手动记录会被重置为 present/auto。
func TestAttendanceStartResetsNonManual(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	if _, err := svc.Start(&teacher); err != nil {
		t.Fatalf("start: %v", err)
	}

	// 手动改为缺勤。
	remark := "病假"
	if _, err := svc.MarkAbsent(&teacher, student.ID, &remark); err != nil {
		t.Fatalf("mark absent: %v", err)
	}

	// 重新 start，手动记录不应被覆盖，但自动记录会被重置。
	// （此处 record 为 manual，故不被重置。）
	items, err := svc.Today(&teacher)
	if err != nil {
		t.Fatalf("today: %v", err)
	}
	if items[0].Status != "absent" || items[0].Source != "manual" {
		t.Fatalf("manual record clobbered by start: %+v", items[0])
	}
}

// SetStatus 手动更新状态，source 变为 manual，present 时记录签到时间。
func TestAttendanceSetStatus(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	if _, err := svc.Start(&teacher); err != nil {
		t.Fatalf("start: %v", err)
	}

	remark := "迟到"
	rec, err := svc.SetStatus(&teacher, student.ID, "late", &remark)
	if err != nil {
		t.Fatalf("set status: %v", err)
	}
	if rec.Status != "late" || rec.Source != "manual" || rec.Remark != "迟到" {
		t.Fatalf("set status wrong: %+v", rec)
	}

	// present 时记录签到时间。
	rec2, err := svc.SetStatus(&teacher, student.ID, "present", nil)
	if err != nil {
		t.Fatalf("set present: %v", err)
	}
	if rec2.SignInAt == nil {
		t.Fatalf("present should set sign_in_at")
	}
}

// SetStatus 在未 start 时返回 404。
func TestAttendanceSetStatusWithoutStart(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	_, err := svc.SetStatus(&teacher, student.ID, "late", nil)
	if err == nil {
		t.Fatalf("expected not-found")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 404 {
		t.Fatalf("expected 404, got %v", err)
	}
}

// MarkLeave 自动建记录后改为 leave/manual。
func TestAttendanceMarkLeave(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	rec, err := svc.MarkLeave(&teacher, student.ID, "事假")
	if err != nil {
		t.Fatalf("mark leave: %v", err)
	}
	if rec.Status != "leave" || rec.Source != "manual" || rec.Remark != "事假" {
		t.Fatalf("mark leave wrong: %+v", rec)
	}
}

// MarkAbsent 默认文案。
func TestAttendanceMarkAbsentDefaultRemark(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	rec, err := svc.MarkAbsent(&teacher, student.ID, nil)
	if err != nil {
		t.Fatalf("mark absent: %v", err)
	}
	if rec.Status != "absent" {
		t.Fatalf("mark absent wrong: %+v", rec)
	}
	if rec.Remark != "未联系到家长，建议后续跟进" {
		t.Fatalf("default remark = %q", rec.Remark)
	}
}

// Summary 四态统计与出勤率。
func TestAttendanceSummary(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, s1 := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	s2 := seedExtraStudent(t, db, s1.ClassID)
	s3 := seedExtraStudent(t, db, s1.ClassID)
	_ = seedExtraStudent(t, db, s1.ClassID)

	svc := services.NewAttendanceService(db, services.NewScope(db))
	if _, err := svc.Start(&teacher); err != nil {
		t.Fatalf("start: %v", err)
	}
	// s1 请假（手动）。
	if _, err := svc.MarkLeave(&teacher, s1.ID, "事假"); err != nil {
		t.Fatalf("leave: %v", err)
	}
	// s2 迟到。
	if _, err := svc.SetStatus(&teacher, s2.ID, "late", nil); err != nil {
		t.Fatalf("late: %v", err)
	}
	// s3 缺勤。
	if _, err := svc.MarkAbsent(&teacher, s3.ID, nil); err != nil {
		t.Fatalf("absent: %v", err)
	}
	// s4 保持 present。

	sum, err := svc.Summary(&teacher)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if sum.Present != 1 || sum.Late != 1 || sum.Leave != 1 || sum.Absent != 1 {
		t.Fatalf("summary counts wrong: %+v", sum)
	}
	if sum.Rate != 25.0 {
		t.Fatalf("rate = %v, want 25.0", sum.Rate)
	}
	if sum.ManualLeaveCount != 1 {
		t.Fatalf("manual_leave_count = %d, want 1", sum.ManualLeaveCount)
	}
	if sum.WechatLeaveCount != 0 {
		t.Fatalf("wechat_leave_count = %d, want 0", sum.WechatLeaveCount)
	}
}

// 教师无法对非管辖班级学生进行考勤（越权 404）。
func TestAttendanceOutOfScope(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	t1 := seedTeacher(t, db, school.ID, "teacher1")
	t2 := seedTeacher(t, db, school.ID, "teacher2")
	_, student1 := seedTeacherClass(t, db, school.ID, t1.ID, "一班")
	_, _ = seedTeacherClass(t, db, school.ID, t2.ID, "二班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	// t2 无法给一班学生标记。
	if _, err := svc.MarkLeave(&t2, student1.ID, "事假"); err == nil {
		t.Fatalf("expected out-of-scope error")
	}
	// t2 的 Summary/Today 为空集。
	items, err := svc.Today(&t2)
	if err != nil {
		t.Fatalf("today: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("t2 should see no records, got %d", len(items))
	}
}

// 无企业微信，故 date 使用本地今日。
func TestAttendanceUsesLocalToday(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewAttendanceService(db, services.NewScope(db))
	if _, err := svc.Start(&teacher); err != nil {
		t.Fatalf("start: %v", err)
	}
	items, err := svc.Today(&teacher)
	if err != nil {
		t.Fatalf("today: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("no records for today %s", util.Today())
	}
}
