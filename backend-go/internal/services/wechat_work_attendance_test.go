// 企微请假同步（matchStudent / applyLeave / HandleWebhookCallback / syncForSchool / SyncAll /
// 点名补同步）的服务层测试。
//
// 全部使用内存 SQLite + httptest 假上游（**不访问真实外网**）：覆盖学生匹配三段优先级与未绑定、
// applyLeave 三分支与 synced_at 回写、回调四个早退分支与成功建记录、syncForSchool 的未配置/已存在跳过/
// 成功同步、SyncAll 汇总（含单校失败计数）、以及 AttendanceService.Start 追加未同步请假记录。
package services_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// 固定装置
// ============================================================

// wechatWorkFixture 一校一班一师一学生 + 一名「家长用户 + 企微绑定」。
type wechatWorkFixture struct {
	DB      *gorm.DB
	School  models.School
	Teacher models.User
	Class   models.ClassRoom
	Student models.Student
	Parent  models.User
	Binding models.ThirdPartyBinding
	Service *services.WechatWorkAttendanceService
}

// httpDoerStub 可编程的 HTTP 客户端（按调用序号返回预置响应，或在第 errAt 次起报错），
// 用于 SyncAll 的「单校失败」分支（httptest 无法表达「第二次请求开始断网」的语义）。
type httpDoerStub struct {
	responses []string
	errAt     int
	calls     int
}

func (d *httpDoerStub) Do(req *http.Request) (*http.Response, error) {
	d.calls++
	if d.errAt > 0 && d.calls >= d.errAt {
		return nil, errors.New("boom: 传输层失败")
	}
	body := "{}"
	if len(d.responses) > 0 {
		idx := d.calls - 1
		if idx >= len(d.responses) {
			idx = len(d.responses) - 1
		}
		body = d.responses[idx]
	}
	return &http.Response{
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(strings.NewReader(body)),
		Header:        http.Header{},
		Request:       req,
		ContentLength: int64(len(body)),
	}, nil
}

// newWechatWorkFixture 建固定装置 + 指向假上游的同步服务。
func newWechatWorkFixture(t *testing.T, handler http.HandlerFunc) *wechatWorkFixture {
	t.Helper()

	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "ww-teacher")
	class := models.ClassRoom{SchoolID: school.ID, Name: "一班", TeacherID: &teacher.ID, Status: "active"}
	require.NoError(t, db.Create(&class).Error)
	student := models.Student{ClassID: class.ID, Name: "小明", Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	parent := models.User{SchoolID: school.ID, Role: "teacher", Username: "ww-parent", Name: "家长", Status: "active"}
	require.NoError(t, db.Create(&parent).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).Update("parent_id", parent.ID).Error)
	student.ParentID = &parent.ID

	binding := models.ThirdPartyBinding{UserID: parent.ID, Platform: "wechat_work", PlatformID: "parent-1"}
	require.NoError(t, db.Create(&binding).Error)

	srv, _ := recordingServer(t, handler)
	wework := services.NewWechatWorkService(db)
	wework.CorpID, wework.Secret, wework.APIBase = "corp-1", "sec-1", srv.URL
	wework.SetHTTPClient(srv.Client())

	return &wechatWorkFixture{
		DB: db, School: school, Teacher: teacher, Class: class, Student: student,
		Parent: parent, Binding: binding,
		Service: services.NewWechatWorkAttendanceService(db, wework),
	}
}

// wechatApprovalUpstream 假企微上游：gettoken 恒定成功，getapprovalinfo 返回 spNos，
// getapprovaldetail 返回按单号区分的详情（未预置的单号返回 errcode 非 0 → 详情为 nil）。
func wechatApprovalUpstream(t *testing.T, spNos []string, details map[string]string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"tok-1","expires_in":7200}`)
		case "/cgi-bin/oa/getapprovalinfo":
			buf, _ := json.Marshal(map[string]any{"errcode": 0, "sp_no_list": spNos, "next_cursor": 0})
			jsonBody(w, string(buf))
		case "/cgi-bin/oa/getapprovaldetail":
			raw, _ := io.ReadAll(r.Body)
			req := map[string]any{}
			_ = json.Unmarshal(raw, &req)
			spNo, _ := req["sp_no"].(string)
			body, ok := details[spNo]
			if !ok {
				jsonBody(w, `{"errcode":40001,"errmsg":"sp_no 不存在"}`)
				return
			}
			jsonBody(w, body)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// approvalDetailJSON 拼一张「已通过」的详情报文（applyer.userid 固定为 parent-1）。
func approvalDetailJSON(spNo, studentName, reason string, leaveStart, leaveEnd int64) string {
	payload := map[string]any{
		"errcode": 0,
		"info": map[string]any{
			"sp_no":     spNo,
			"sp_status": 2,
			"applyer":   map[string]any{"userid": "parent-1"},
			"sp_record": []any{map[string]any{"sp_status": 2, "approve_time": 1700000000}},
			"apply_data": map[string]any{"contents": []any{
				map[string]any{"title": "学生姓名", "value": map[string]any{"text": studentName}},
				map[string]any{"title": "请假类型", "value": map[string]any{"text": "病假"}},
				map[string]any{"title": "请假事由", "value": map[string]any{"text": reason}},
				map[string]any{"title": "请假时间", "value": map[string]any{
					"new_begin": leaveStart, "new_end": leaveEnd,
				}},
			}},
		},
	}
	buf, _ := json.Marshal(payload)
	return string(buf)
}

// pendingApprovalDetailJSON 拼一张「待审批」（sp_status=1）的详情报文。
func pendingApprovalDetailJSON(spNo string, leaveStart, leaveEnd int64) string {
	return `{"errcode":0,"info":{"sp_no":"` + spNo + `","sp_status":1,"applyer":{"userid":"parent-1"},` +
		`"apply_data":{"contents":[{"title":"学生姓名","value":{"text":"小明"}},` +
		`{"title":"请假时间","value":{"new_begin":` + strconv.FormatInt(leaveStart, 10) +
		`,"new_end":` + strconv.FormatInt(leaveEnd, 10) + `}}]}}}`
}

// seedLeaveRecord 直接插入一条请假记录（返回落库后的实体，含 ID）。
func seedLeaveRecord(t *testing.T, db *gorm.DB, schoolID uint, rec models.WechatWorkLeaveRecord) models.WechatWorkLeaveRecord {
	t.Helper()
	rec.SchoolID = schoolID
	if rec.ApproveStatus == "" {
		rec.ApproveStatus = "approved"
	}
	if rec.ParentWeworkUserid == "" {
		rec.ParentWeworkUserid = "parent-1"
	}
	require.NoError(t, db.Create(&rec).Error)
	return rec
}

// ============================================================
// matchStudent
// ============================================================

// TestWechatWorkMatchStudentPriorities 精确 → 包含 → 第一条，以及未绑定/全停用返回 nil。
func TestWechatWorkMatchStudentPriorities(t *testing.T) {
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	db := f.DB

	// 家长名下再加两名学生（id 均大于「小明」）。
	like := models.Student{ClassID: f.Class.ID, Name: "小明的妹妹", Status: "active", ParentID: &f.Parent.ID}
	require.NoError(t, db.Create(&like).Error)
	exact := models.Student{ClassID: f.Class.ID, Name: "小明同学", Status: "active", ParentID: &f.Parent.ID}
	require.NoError(t, db.Create(&exact).Error)

	// 1) 精确命中：hint = 小明同学 → 命中该生（而不是 id 更小、姓名包含「小明」的其它学生）。
	hint := "小明同学"
	got, err := f.Service.MatchStudent("parent-1", &hint)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, exact.ID, got.ID, "应优先精确同名")

	// 2) 无精确匹配 → 退到 LIKE：hint = 妹妹。
	partialHint := "妹妹"
	got, err = f.Service.MatchStudent("parent-1", &partialHint)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, like.ID, got.ID, "无精确匹配时应命中姓名包含者")

	// 3) 都命中不到 → 取该家长名下第一条（按 id 升序，即小明）。
	missHint := "张三"
	got, err = f.Service.MatchStudent("parent-1", &missHint)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, f.Student.ID, got.ID, "都匹配不到时回退第一条")

	// 4) hint 为 nil → 直接第一条。
	got, err = f.Service.MatchStudent("parent-1", nil)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, f.Student.ID, got.ID)

	// 5) 未绑定该企微 userid → nil。
	got, err = f.Service.MatchStudent("unknown-user", &hint)
	require.NoError(t, err)
	assert.Nil(t, got)

	// 6) 家长名下学生全部停用 → nil（status=active 过滤）。
	require.NoError(t, db.Model(&models.Student{}).Where("parent_id = ?", f.Parent.ID).
		Update("status", "inactive").Error)
	got, err = f.Service.MatchStudent("parent-1", &hint)
	require.NoError(t, err)
	assert.Nil(t, got)
}

// ============================================================
// applyLeave
// ============================================================

// TestWechatWorkApplyLeaveCreates 无考勤记录 → 新建 leave/wechat_work + 回写 synced_at。
func TestWechatWorkApplyLeaveCreates(t *testing.T) {
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	date := util.Today()
	rec := seedLeaveRecord(t, f.DB, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-CREATE",
		LeaveStartDate: date, LeaveEndDate: date, Reason: "发烧",
	})

	require.NoError(t, f.Service.ApplyLeave(&f.Student, &rec, date))

	var att models.Attendance
	require.NoError(t, f.DB.Where("student_id = ? AND date = ?", f.Student.ID, date).First(&att).Error)
	assert.Equal(t, "leave", att.Status)
	assert.Equal(t, "wechat_work", att.Source)
	assert.Equal(t, "发烧", att.Remark)
	require.NotNil(t, att.TeacherID)
	assert.Equal(t, f.Teacher.ID, *att.TeacherID, "teacher_id 取该生班级的班主任")
	require.NotNil(t, att.LeaveRecordID)
	assert.Equal(t, rec.ID, *att.LeaveRecordID)

	var reloaded models.WechatWorkLeaveRecord
	require.NoError(t, f.DB.First(&reloaded, rec.ID).Error)
	require.NotNil(t, reloaded.SyncedAt, "新建考勤后应回写 synced_at")
	require.NotNil(t, rec.SyncedAt)
}

// TestWechatWorkApplyLeaveSkipsManual 已存在且 source=manual → 不改动、不回写 synced_at。
func TestWechatWorkApplyLeaveSkipsManual(t *testing.T) {
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	date := util.Today()
	manual := models.Attendance{
		ClassID: f.Class.ID, StudentID: f.Student.ID, TeacherID: &f.Teacher.ID,
		Date: date, Status: "late", Source: "manual", Remark: "老师手动标记",
	}
	require.NoError(t, f.DB.Create(&manual).Error)
	rec := seedLeaveRecord(t, f.DB, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-MANUAL",
		LeaveStartDate: date, LeaveEndDate: date, Reason: "发烧",
	})

	require.NoError(t, f.Service.ApplyLeave(&f.Student, &rec, date))

	var att models.Attendance
	require.NoError(t, f.DB.First(&att, manual.ID).Error)
	assert.Equal(t, "late", att.Status)
	assert.Equal(t, "manual", att.Source)
	assert.Equal(t, "老师手动标记", att.Remark)
	assert.Nil(t, att.LeaveRecordID)

	var reloaded models.WechatWorkLeaveRecord
	require.NoError(t, f.DB.First(&reloaded, rec.ID).Error)
	assert.Nil(t, reloaded.SyncedAt, "跳过手动记录时不回写 synced_at")
}

// TestWechatWorkApplyLeaveUpdatesExisting 已存在且非 manual → 更新为 leave/wechat_work + 回写 synced_at。
func TestWechatWorkApplyLeaveUpdatesExisting(t *testing.T) {
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	date := util.Today()
	auto := models.Attendance{
		ClassID: f.Class.ID, StudentID: f.Student.ID, TeacherID: &f.Teacher.ID,
		Date: date, Status: "present", Source: "auto",
	}
	require.NoError(t, f.DB.Create(&auto).Error)
	rec := seedLeaveRecord(t, f.DB, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-UPDATE",
		LeaveStartDate: date, LeaveEndDate: date, Reason: "事假",
	})

	require.NoError(t, f.Service.ApplyLeave(&f.Student, &rec, date))

	var att models.Attendance
	require.NoError(t, f.DB.First(&att, auto.ID).Error)
	assert.Equal(t, "leave", att.Status)
	assert.Equal(t, "wechat_work", att.Source)
	assert.Equal(t, "事假", att.Remark)
	require.NotNil(t, att.LeaveRecordID)
	assert.Equal(t, rec.ID, *att.LeaveRecordID)

	var reloaded models.WechatWorkLeaveRecord
	require.NoError(t, f.DB.First(&reloaded, rec.ID).Error)
	require.NotNil(t, reloaded.SyncedAt, "更新考勤后应回写 synced_at")
}

// ============================================================
// HandleWebhookCallback
// ============================================================

// TestWechatWorkHandleWebhookCallbackEarlyReturns sp_status != 2 / sp_no 为空 / 重复 sp_no / 详情 nil 四个早退分支。
func TestWechatWorkHandleWebhookCallbackEarlyReturns(t *testing.T) {
	date := util.Today()

	t.Run("sp_status 非 2 直接返回且不请求上游", func(t *testing.T) {
		detailCalls := 0
		f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/cgi-bin/oa/getapprovaldetail" {
				detailCalls++
			}
			w.WriteHeader(http.StatusNotFound)
		})
		require.NoError(t, f.Service.HandleWebhookCallback(f.School.ID, "SP-1", 1))
		assert.Equal(t, 0, detailCalls)

		var count int64
		require.NoError(t, f.DB.Model(&models.WechatWorkLeaveRecord{}).Count(&count).Error)
		assert.Equal(t, int64(0), count)
	})

	t.Run("sp_no 为空直接返回", func(t *testing.T) {
		f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
		require.NoError(t, f.Service.HandleWebhookCallback(f.School.ID, "", 2))

		var count int64
		require.NoError(t, f.DB.Model(&models.WechatWorkLeaveRecord{}).Count(&count).Error)
		assert.Equal(t, int64(0), count)
	})

	t.Run("重复 sp_no 直接返回且不查详情", func(t *testing.T) {
		detailCalls := 0
		f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/cgi-bin/oa/getapprovaldetail" {
				detailCalls++
			}
			w.WriteHeader(http.StatusNotFound)
		})
		seedLeaveRecord(t, f.DB, f.School.ID, models.WechatWorkLeaveRecord{
			ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-DUP",
			LeaveStartDate: date, LeaveEndDate: date,
		})

		require.NoError(t, f.Service.HandleWebhookCallback(f.School.ID, "SP-DUP", 2))
		assert.Equal(t, 0, detailCalls)

		var count int64
		require.NoError(t, f.DB.Model(&models.WechatWorkLeaveRecord{}).Count(&count).Error)
		assert.Equal(t, int64(1), count)
	})

	t.Run("详情为 nil 时不落库", func(t *testing.T) {
		f := newWechatWorkFixture(t, wechatApprovalUpstream(t, nil, nil))
		require.NoError(t, f.Service.HandleWebhookCallback(f.School.ID, "SP-MISSING", 2))

		var count int64
		require.NoError(t, f.DB.Model(&models.WechatWorkLeaveRecord{}).Count(&count).Error)
		assert.Equal(t, int64(0), count)
	})
}

// TestWechatWorkHandleWebhookCallbackSuccess 成功分支：建请假记录（approved）+ 落考勤 + 幂等。
func TestWechatWorkHandleWebhookCallbackSuccess(t *testing.T) {
	detail := approvalDetailJSON("SP-OK", "小明", "发烧", 0, 0)
	f := newWechatWorkFixture(t, wechatApprovalUpstream(t, []string{"SP-OK"}, map[string]string{"SP-OK": detail}))
	today := util.Today()

	require.NoError(t, f.Service.HandleWebhookCallback(f.School.ID, "SP-OK", 2))

	var rec models.WechatWorkLeaveRecord
	require.NoError(t, f.DB.Where("sp_no = ?", "SP-OK").First(&rec).Error)
	assert.Equal(t, "approved", rec.ApproveStatus)
	require.NotNil(t, rec.ApprovedAt, "回调路径 approved_at 取 now")
	assert.Equal(t, "parent-1", rec.ParentWeworkUserid)
	assert.Equal(t, "小明", rec.StudentNameFromWework)
	require.NotNil(t, rec.StudentID)
	assert.Equal(t, f.Student.ID, *rec.StudentID)
	require.NotNil(t, rec.ClassID)
	assert.Equal(t, f.Class.ID, *rec.ClassID)
	// 详情未带起止时间 → 回退今天。
	assert.Equal(t, today, rec.LeaveStartDate)
	assert.Equal(t, today, rec.LeaveEndDate)
	assert.Equal(t, "病假", rec.LeaveType)
	assert.Equal(t, "发烧", rec.Reason)
	assert.Contains(t, rec.RawData, `"sp_no":"SP-OK"`, "raw_data 存详情 JSON")
	require.NotNil(t, rec.SyncedAt, "匹配到学生时应已同步考勤")

	var att models.Attendance
	require.NoError(t, f.DB.Where("student_id = ? AND date = ?", f.Student.ID, today).First(&att).Error)
	assert.Equal(t, "leave", att.Status)
	assert.Equal(t, "wechat_work", att.Source)
	require.NotNil(t, att.LeaveRecordID)
	assert.Equal(t, rec.ID, *att.LeaveRecordID)

	// 幂等：再调一次不新增记录。
	require.NoError(t, f.Service.HandleWebhookCallback(f.School.ID, "SP-OK", 2))
	var count int64
	require.NoError(t, f.DB.Model(&models.WechatWorkLeaveRecord{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

// ============================================================
// syncForSchool / syncAll
// ============================================================

// TestWechatWorkSyncForSchoolUnconfigured corp_id 未配置 → 零值返回且不请求上游。
func TestWechatWorkSyncForSchoolUnconfigured(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	stub := &httpDoerStub{}
	wework := services.NewWechatWorkService(db)
	wework.SetHTTPClient(stub)
	svc := services.NewWechatWorkAttendanceService(db, wework)

	result, err := svc.SyncForSchool(school.ID, "2026-09-20")
	require.NoError(t, err)
	assert.Equal(t, services.WechatWorkLeaveSyncResult{}, result)
	assert.Equal(t, 0, stub.calls)
}

// TestWechatWorkSyncForSchoolSkipsExisting 已存在同 sp_no → skipped++（不查详情），第二次调用全跳过。
func TestWechatWorkSyncForSchoolSkipsExisting(t *testing.T) {
	date := util.Today()
	detailCalls := 0
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"tok-1"}`)
		case "/cgi-bin/oa/getapprovalinfo":
			jsonBody(w, `{"errcode":0,"sp_no_list":["SP-OLD","SP-NEW"],"next_cursor":0}`)
		case "/cgi-bin/oa/getapprovaldetail":
			detailCalls++
			jsonBody(w, approvalDetailJSON("SP-NEW", "小明", "感冒", 0, 0))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	syncedAt := util.Now()
	seedLeaveRecord(t, f.DB, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-OLD",
		LeaveStartDate: date, LeaveEndDate: date, SyncedAt: &syncedAt,
	})

	result, err := f.Service.SyncForSchool(f.School.ID, date)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Skipped)
	assert.Equal(t, 1, result.Synced)
	assert.Equal(t, 1, detailCalls, "已存在的单号不应再查详情")

	var rec models.WechatWorkLeaveRecord
	require.NoError(t, f.DB.Where("sp_no = ?", "SP-NEW").First(&rec).Error)
	assert.Equal(t, "approved", rec.ApproveStatus)
	assert.Equal(t, "感冒", rec.Reason)
	require.NotNil(t, rec.SyncedAt)

	// 幂等：第二次全部跳过。
	result, err = f.Service.SyncForSchool(f.School.ID, date)
	require.NoError(t, err)
	assert.Equal(t, services.WechatWorkLeaveSyncResult{Skipped: 2}, result)
}

// TestWechatWorkSyncForSchoolRequestsDayBounds 请求体里的起止时间是「当天 00:00:00 / 23:59:59」（业务时区）。
func TestWechatWorkSyncForSchoolRequestsDayBounds(t *testing.T) {
	date := "2026-09-20"
	bodies := []map[string]any{}
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"tok-1"}`)
		case "/cgi-bin/oa/getapprovalinfo":
			raw, _ := io.ReadAll(r.Body)
			payload := map[string]any{}
			_ = json.Unmarshal(raw, &payload)
			bodies = append(bodies, payload)
			jsonBody(w, `{"errcode":0,"sp_no_list":[],"next_cursor":0}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	result, err := f.Service.SyncForSchool(f.School.ID, date)
	require.NoError(t, err)
	assert.Equal(t, services.WechatWorkLeaveSyncResult{}, result)

	day, err := time.ParseInLocation("2006-01-02", date, util.Loc)
	require.NoError(t, err)
	require.Len(t, bodies, 1)
	assert.Equal(t, float64(day.Unix()), bodies[0]["starttime"])
	assert.Equal(t, float64(day.Add(24*time.Hour-time.Second).Unix()), bodies[0]["endtime"])
}

// TestWechatWorkSyncForSchoolDatesFromDetail 起止日期取详情（而非回退当日）；approve_status 照抄上游，
// 且**只要匹配到学生就补考勤**（Laravel syncForSchool 同样不校验审批状态）。
func TestWechatWorkSyncForSchoolDatesFromDetail(t *testing.T) {
	date := util.Today()
	start := time.Date(2026, 9, 18, 8, 0, 0, 0, util.Loc).Unix()
	end := time.Date(2026, 9, 22, 17, 0, 0, 0, util.Loc).Unix()
	f := newWechatWorkFixture(t, wechatApprovalUpstream(t,
		[]string{"SP-PENDING"},
		map[string]string{"SP-PENDING": pendingApprovalDetailJSON("SP-PENDING", start, end)},
	))

	result, err := f.Service.SyncForSchool(f.School.ID, date)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Synced)

	var rec models.WechatWorkLeaveRecord
	require.NoError(t, f.DB.Where("sp_no = ?", "SP-PENDING").First(&rec).Error)
	assert.Equal(t, "pending", rec.ApproveStatus, "上游 sp_status=1 → pending")
	assert.Equal(t, "2026-09-18", rec.LeaveStartDate)
	assert.Equal(t, "2026-09-22", rec.LeaveEndDate)
	require.NotNil(t, rec.SyncedAt)

	var att models.Attendance
	require.NoError(t, f.DB.Where("student_id = ? AND date = ?", f.Student.ID, date).First(&att).Error)
	assert.Equal(t, "leave", att.Status)
	assert.Equal(t, "wechat_work", att.Source)

	// pending 记录不会被「点名补同步」再次处理（synced_at 已非 NULL）。
	applied, err := f.Service.ApplyUnsyncedLeavesForClass(f.Class.ID, date)
	require.NoError(t, err)
	assert.Equal(t, 0, applied)
}

// TestWechatWorkSyncAllAggregates 遍历 active 学校累加；单校传输层失败计入 failed 并继续。
func TestWechatWorkSyncAllAggregates(t *testing.T) {
	db := setupDB(t)
	first := seedSchool(t, db)
	second := models.School{Name: "第二学校", Code: "ww-school-2", Status: "active"}
	require.NoError(t, db.Create(&second).Error)
	inactive := models.School{Name: "停用学校", Code: "ww-school-3", Status: "inactive"}
	require.NoError(t, db.Create(&inactive).Error)

	// 调用序列：学校1 gettoken(1) → getapprovalinfo(2) → getapprovaldetail(3) → 学校2 gettoken(4) 失败。
	stub := &httpDoerStub{
		responses: []string{
			`{"errcode":0,"access_token":"tok-first"}`,
			`{"errcode":0,"sp_no_list":["SP-A"],"next_cursor":0}`,
			approvalDetailJSON("SP-A", "小明", "病假", 0, 0),
		},
		errAt: 4,
	}
	wework := services.NewWechatWorkService(db)
	wework.CorpID, wework.Secret = "corp-1", "sec-1"
	wework.SetHTTPClient(stub)
	svc := services.NewWechatWorkAttendanceService(db, wework)

	result, err := svc.SyncAll(util.Today())
	require.NoError(t, err)
	assert.Equal(t, 1, result.Synced, "第一所学校成功同步 1 条")
	assert.Equal(t, 1, result.Failed, "第二所学校上游失败计入 failed")
	assert.Equal(t, 0, result.Skipped)

	var rec models.WechatWorkLeaveRecord
	require.NoError(t, db.Where("sp_no = ?", "SP-A").First(&rec).Error)
	assert.Equal(t, first.ID, rec.SchoolID)
}

// TestWechatWorkSyncAllDefaultDate 空日期回退今天，非 active 学校被跳过。
func TestWechatWorkSyncAllDefaultDate(t *testing.T) {
	db := setupDB(t)
	seedSchool(t, db)
	stub := &httpDoerStub{responses: []string{
		`{"errcode":0,"access_token":"tok"}`,
		`{"errcode":0,"sp_no_list":[],"next_cursor":0}`,
	}}
	wework := services.NewWechatWorkService(db)
	wework.CorpID, wework.Secret = "corp-1", "sec-1"
	wework.SetHTTPClient(stub)
	svc := services.NewWechatWorkAttendanceService(db, wework)

	result, err := svc.SyncAll("")
	require.NoError(t, err)
	assert.Equal(t, services.WechatWorkLeaveSyncResult{}, result)
}

// ============================================================
// 点名补同步（Laravel startAttendanceForClass 后半段）
// ============================================================

// TestWechatWorkApplyUnsyncedLeavesForClass 未同步记录生效；学生为空 / 他班 / pending / 日期外一律跳过。
func TestWechatWorkApplyUnsyncedLeavesForClass(t *testing.T) {
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	db, date := f.DB, util.Today()

	hit := seedLeaveRecord(t, db, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-HIT",
		LeaveStartDate: date, LeaveEndDate: date, Reason: "发烧",
	})
	// 学生为空。
	seedLeaveRecord(t, db, f.School.ID, models.WechatWorkLeaveRecord{
		SpNo: "SP-NOSTUDENT", LeaveStartDate: date, LeaveEndDate: date,
	})
	// 日期不覆盖今天。
	seedLeaveRecord(t, db, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-OLD-DATE",
		LeaveStartDate: "2000-01-01", LeaveEndDate: "2000-01-02",
	})
	// 未通过。
	seedLeaveRecord(t, db, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-PENDING",
		LeaveStartDate: date, LeaveEndDate: date, ApproveStatus: "pending",
	})
	// 他班学生。
	otherClass := models.ClassRoom{SchoolID: f.School.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&otherClass).Error)
	otherStudent := models.Student{ClassID: otherClass.ID, Name: "小红", Status: "active"}
	require.NoError(t, db.Create(&otherStudent).Error)
	seedLeaveRecord(t, db, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &otherClass.ID, StudentID: &otherStudent.ID, SpNo: "SP-OTHER-CLASS",
		LeaveStartDate: date, LeaveEndDate: date,
	})

	applied, err := f.Service.ApplyUnsyncedLeavesForClass(f.Class.ID, date)
	require.NoError(t, err)
	assert.Equal(t, 1, applied)

	var att models.Attendance
	require.NoError(t, db.Where("student_id = ? AND date = ?", f.Student.ID, date).First(&att).Error)
	assert.Equal(t, "leave", att.Status)
	assert.Equal(t, "wechat_work", att.Source)

	var reloaded models.WechatWorkLeaveRecord
	require.NoError(t, db.First(&reloaded, hit.ID).Error)
	require.NotNil(t, reloaded.SyncedAt)

	// 他班学生未被动过。
	var otherCount int64
	require.NoError(t, db.Model(&models.Attendance{}).Where("student_id = ?", otherStudent.ID).Count(&otherCount).Error)
	assert.Equal(t, int64(0), otherCount)
}

// TestAttendanceStartAppliesUnsyncedWechatLeave 教师开始点名时补写企微请假，并把 wechat_leave_count 计入统计。
func TestAttendanceStartAppliesUnsyncedWechatLeave(t *testing.T) {
	f := newWechatWorkFixture(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	db, date := f.DB, util.Today()
	seedLeaveRecord(t, db, f.School.ID, models.WechatWorkLeaveRecord{
		ClassID: &f.Class.ID, StudentID: &f.Student.ID, SpNo: "SP-START",
		LeaveStartDate: date, LeaveEndDate: date, Reason: "病假",
	})

	attendance := services.NewAttendanceService(db, services.NewScope(db), f.Service)
	result, err := attendance.Start(&f.Teacher)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 1, result.WechatLeaveCount, "当日 source=wechat_work 的记录应计入")

	var att models.Attendance
	require.NoError(t, db.Where("student_id = ? AND date = ?", f.Student.ID, date).First(&att).Error)
	assert.Equal(t, "leave", att.Status)
	assert.Equal(t, "wechat_work", att.Source)
	assert.Equal(t, "病假", att.Remark)

	summary, err := attendance.Summary(&f.Teacher)
	require.NoError(t, err)
	assert.Equal(t, 1, summary.Leave)
	assert.Equal(t, 1, summary.WechatLeaveCount)
	assert.Equal(t, 0, summary.ManualLeaveCount)

	// 教师手动改判后：手动记录受保护（不重置），且已同步的记录不再改写。
	if _, err := attendance.MarkLeave(&f.Teacher, f.Student.ID, "老师手填"); err != nil {
		t.Fatalf("mark leave: %v", err)
	}
	result, err = attendance.Start(&f.Teacher)
	require.NoError(t, err)
	var after models.Attendance
	require.NoError(t, db.First(&after, att.ID).Error)
	assert.Equal(t, "manual", after.Source, "manual 记录不应被点名重置")
	assert.Equal(t, 0, result.WechatLeaveCount)
}
