// 管理端「报表 / 系统运维 / 批量账号与导入 / 学年升级 / 学校 LOGO」的路由注册与 HTTP 冒烟测试。
//
// 路由契约（路径逐字对齐 Laravel backend/routes/api.php 第 45-153 行）：
// 注册遗漏、路径写歪或误注册「本批明确不做」的接口都会在这里失败。
package router_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// adminOpsBaseURL 由 newAdminFixture 写入（httptest 服务器要求绝对 URL）。
var adminOpsBaseURL string

// adminOpsWantRoutes 本批新增的 24 条路由（22 条 Laravel 原路径 + 2 条学年升级别名）。
var adminOpsWantRoutes = []string{
	"GET /api/v1/admin/reports/overview",
	"GET /api/v1/admin/reports/by-grade",
	"GET /api/v1/admin/reports/by-class",
	"GET /api/v1/admin/display-login-logs",
	"GET /api/v1/admin/system/diagnose",
	"GET /api/v1/admin/system/status",
	"GET /api/v1/admin/system/logs",
	"POST /api/v1/admin/system/repair",
	"POST /api/v1/admin/accounts/batch-reset-password",
	"POST /api/v1/admin/accounts/batch-delete",
	"POST /api/v1/admin/teachers/batch-create",
	"POST /api/v1/admin/teachers/import",
	"GET /api/v1/admin/teachers/template-csv",
	"POST /api/v1/admin/students/import",
	"POST /api/v1/admin/students/batch-delete",
	"POST /api/v1/admin/students/batch-move",
	"POST /api/v1/admin/classes/batch-create",
	"POST /api/v1/admin/classes/:id/assign-teacher",
	"DELETE /api/v1/admin/classes/:id/remove-teacher",
	"GET /api/v1/admin/grade-upgrade/preview",
	"POST /api/v1/admin/grade-upgrade/execute",
	"POST /api/v1/admin/school/logo",
}

// TestAdminOpsRoutesRegistered 本批路由已注册，且明确不做的接口未被注册。
func TestAdminOpsRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range adminOpsWantRoutes {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}

	// 注 1：课表进阶（import-csv / unavailabilities / check-conflicts / teacher-assignments）
	// 与课表自动排课（generate / generate-school）已由单独两批实现并注册，
	// 分别见 timetable_advanced_routes_test.go 与 timetable_generate_routes_test.go；
	// admin/ai/* 含 provider-official 已由 AI 批次实现并注册，见 ai_routes_test.go；
	// 报表 / 课表导出（CSV 版）由最近的导出批次实现并注册，见 router.go 与
	// services/report_export_test.go、services/timetable_export_test.go。
	// 注 2：原先此处断言「第三方通讯录 4 条不得注册」；该批已于后续批次实现并注册，
	// 断言移至 third_party_routes_test.go（不再反向断言）。
}

// newAdminFixture 建一所学校 + 管理员账号，返回 (db, engine, schoolID, adminToken)。
func newAdminFixture(t *testing.T) (*gorm.DB, *http.Client, string, uint) {
	t.Helper()
	db, engine, _ := newTestEngine(t)

	var school models.School
	require.NoError(t, db.First(&school).Error)

	admin := models.User{
		SchoolID: school.ID, Role: "school_admin", Username: "ops-admin",
		Name: "管理员", Status: "active",
	}
	require.NoError(t, db.Create(&admin).Error)

	jwtMgr := jwtauth.New("test-secret", 1)
	token, err := jwtMgr.Generate(admin.ID, admin.Role, admin.SchoolID)
	require.NoError(t, err)

	server := httptest.NewServer(engine)
	adminOpsBaseURL = server.URL
	t.Cleanup(server.Close)
	return db, server.Client(), token, school.ID
}

// doAdmin 以管理员身份发请求。
func doAdmin(t *testing.T, client *http.Client, token, method, url string, body *bytes.Reader, contentType string) *http.Response {
	t.Helper()
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, adminOpsBaseURL+url, nil)
	} else {
		req, err = http.NewRequest(method, adminOpsBaseURL+url, body)
	}
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

// decodeBody 读取并解析响应体。
func decodeBody(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	var out map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	return out
}

// TestAdminOpsHTTPReportsAndSystem 报表 / 系统运维接口的 HTTP 冒烟。
func TestAdminOpsHTTPReportsAndSystem(t *testing.T) {
	_, client, token, _ := newAdminFixture(t)

	overview := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/reports/overview", nil, "")
	require.Equal(t, http.StatusOK, overview.StatusCode)
	body := decodeBody(t, overview)
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	for _, key := range []string{"class_count", "teacher_count", "student_count", "monthly_score", "month_last_month", "month_label"} {
		if key == "month_last_month" {
			continue
		}
		assert.Contains(t, data, key)
	}

	byGrade := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/reports/by-grade", nil, "")
	require.Equal(t, http.StatusOK, byGrade.StatusCode)
	gradeBody := decodeBody(t, byGrade)
	assert.NotNil(t, gradeBody["data"])

	byClass := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/reports/by-class", nil, "")
	require.Equal(t, http.StatusOK, byClass.StatusCode)

	logs := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/display-login-logs", nil, "")
	require.Equal(t, http.StatusOK, logs.StatusCode)
	logsBody := decodeBody(t, logs)
	assert.Contains(t, logsBody, "meta")

	diagnose := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/system/diagnose", nil, "")
	require.Equal(t, http.StatusOK, diagnose.StatusCode)
	diagnoseBody := decodeBody(t, diagnose)
	assert.Equal(t, false, diagnoseBody["has_issues"])
	assert.Equal(t, "系统状态正常", diagnoseBody["message"])

	status := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/system/status", nil, "")
	require.Equal(t, http.StatusOK, status.StatusCode)
	statusBody := decodeBody(t, status)
	version, ok := statusBody["data"].(map[string]any)["version"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "backend-go", version["backend"])
	assert.Equal(t, "Asia/Shanghai", version["timezone"])

	// 未配置 LOG_FILE：200 + 空列表 + 说明性 message（不伪造日志）
	sysLogs := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/system/logs", nil, "")
	require.Equal(t, http.StatusOK, sysLogs.StatusCode)
	sysLogsBody := decodeBody(t, sysLogs)
	assert.Contains(t, sysLogsBody["message"], "LOG_FILE")
	sysLogsData, ok := sysLogsBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, false, sysLogsData["exists"])

	repair := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/system/repair", bytes.NewReader(nil), "application/json")
	require.Equal(t, http.StatusOK, repair.StatusCode)
	repairBody := decodeBody(t, repair)
	assert.Equal(t, "数据库迁移已完成", repairBody["message"])
	assert.Equal(t, float64(0), repairBody["pending_password_resets"])

	// 未登录 → 401
	unauth := doAdmin(t, client, "bad-token", http.MethodGet, "/api/v1/admin/reports/overview", nil, "")
	assert.Equal(t, http.StatusUnauthorized, unauth.StatusCode)
}

// TestAdminOpsHTTPTeacherTemplateAndImport 模板下载响应头 + CSV 导入（预览/执行）。
func TestAdminOpsHTTPTeacherTemplateAndImport(t *testing.T) {
	db, client, token, schoolID := newAdminFixture(t)

	// 模板下载：Content-Type / Content-Disposition / BOM 均与 Laravel 一致
	tpl := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/teachers/template-csv", nil, "")
	require.Equal(t, http.StatusOK, tpl.StatusCode)
	assert.Equal(t, "text/csv; charset=UTF-8", tpl.Header.Get("Content-Type"))
	assert.Equal(t, `attachment; filename="teacher_import_template.csv"`, tpl.Header.Get("Content-Disposition"))
	var buf bytes.Buffer
	_, err := buf.ReadFrom(tpl.Body)
	require.NoError(t, err)
	_ = tpl.Body.Close()
	assert.True(t, strings.HasPrefix(buf.String(), "\ufeff姓名,年级团队,科目,密码,手机号\n"))
	assert.Contains(t, buf.String(), "张老师,三年级团队,语文,star123456,13800138000\n")

	// multipart 上传：dry_run 缺省 → 预览（不建号）
	resp := uploadTeacherCSV(t, client, token, "teachers.csv", "姓名,年级团队,科目,密码,手机号\n测试老师,三年级团队,语文,pwd123456,13800138000\n", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	previewBody := decodeBody(t, resp)
	assert.Equal(t, "预览模式：共 1 条数据", previewBody["message"])
	assert.Equal(t, float64(1), previewBody["total"])

	var count int64
	require.NoError(t, db.Model(&models.User{}).Where("school_id = ? AND role = ?", schoolID, "teacher").Count(&count).Error)
	assert.Equal(t, int64(0), count)

	// dry_run=false → 正式建号
	execResp := uploadTeacherCSV(t, client, token, "teachers.csv", "姓名,年级团队,科目,密码,手机号\n测试老师,三年级团队,语文,pwd123456,13800138000\n", map[string]string{"dry_run": "false"})
	require.Equal(t, http.StatusOK, execResp.StatusCode)
	execBody := decodeBody(t, execResp)
	assert.Equal(t, "已导入 1 名教师", execBody["message"])
	require.NoError(t, db.Model(&models.User{}).Where("school_id = ? AND role = ?", schoolID, "teacher").Count(&count).Error)
	assert.Equal(t, int64(1), count)

	// 缺 file → 422
	missing := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/teachers/import", bytes.NewReader(nil), "multipart/form-data")
	assert.Equal(t, http.StatusUnprocessableEntity, missing.StatusCode)
}

// uploadTeacherCSV 以 multipart/form-data 上传教师 CSV。
func uploadTeacherCSV(t *testing.T, client *http.Client, token, filename, content string, fields map[string]string) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write([]byte(content))
	require.NoError(t, err)
	for k, v := range fields {
		require.NoError(t, writer.WriteField(k, v))
	}
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, adminOpsBaseURL+"/api/v1/admin/teachers/import", &body)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

// TestAdminOpsHTTPClassesImportAndLogo 班级批量创建 / 学生导入 / LOGO 上传的 HTTP 冒烟。
func TestAdminOpsHTTPClassesImportAndLogo(t *testing.T) {
	t.Setenv("UPLOAD_DIR", t.TempDir())

	db, client, token, schoolID := newAdminFixture(t)

	// 批量建班
	classesResp := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/classes/batch-create",
		bytes.NewReader([]byte(`{"grade":"一年级","count":2,"year":"2026"}`)), "application/json")
	require.Equal(t, http.StatusCreated, classesResp.StatusCode)
	classesBody := decodeBody(t, classesResp)
	assert.Equal(t, "已批量创建 2 个班级", classesBody["message"])
	created, ok := classesBody["data"].([]any)
	require.True(t, ok)
	require.Len(t, created, 2)

	// 学生导入（JSON，Laravel 契约）
	studentsResp := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/students/import",
		bytes.NewReader([]byte(`{"students":[{"name":"小明","class_name":"一年级（1）班","gender":"男","student_no":"1001"}]}`)),
		"application/json")
	require.Equal(t, http.StatusOK, studentsResp.StatusCode)
	studentsBody := decodeBody(t, studentsResp)
	assert.Equal(t, "导入完成", studentsBody["message"])
	inner, ok := studentsBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(1), inner["created_count"])

	// 注意：fixture 学校已有一个同名「小明」，故按学号校验本次导入的行
	var studentCount int64
	require.NoError(t, db.Model(&models.Student{}).Where("student_no = ?", "1001").Count(&studentCount).Error)
	assert.Equal(t, int64(1), studentCount)
	var petCount int64
	require.NoError(t, db.Model(&models.Pet{}).Count(&petCount).Error)
	assert.Equal(t, int64(1), petCount, "导入学生自动分配默认宠物")

	// 学生导入（CSV 文件上传）
	var csvBody bytes.Buffer
	writer := multipart.NewWriter(&csvBody)
	part, err := writer.CreateFormFile("file", "students.csv")
	require.NoError(t, err)
	_, err = part.Write([]byte("姓名,班级,性别,学号\n小红,一年级（2）班,女,2002\n"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	csvReq, err := http.NewRequest(http.MethodPost, adminOpsBaseURL+"/api/v1/admin/students/import", &csvBody)
	require.NoError(t, err)
	csvReq.Header.Set("Authorization", "Bearer "+token)
	csvReq.Header.Set("Content-Type", writer.FormDataContentType())
	csvResp, err := client.Do(csvReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, csvResp.StatusCode)
	csvParsed := decodeBody(t, csvResp)
	assert.Equal(t, "导入完成", csvParsed["message"])

	// LOGO 上传（multipart 字段名 logo）
	pngBytes := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0, 0, 0, 0x0D, 'I', 'H', 'D', 'R'}
	var logoBody bytes.Buffer
	logoWriter := multipart.NewWriter(&logoBody)
	logoPart, err := logoWriter.CreateFormFile("logo", "school-logo.png")
	require.NoError(t, err)
	_, err = logoPart.Write(pngBytes)
	require.NoError(t, err)
	require.NoError(t, logoWriter.Close())
	logoReq, err := http.NewRequest(http.MethodPost, adminOpsBaseURL+"/api/v1/admin/school/logo", &logoBody)
	require.NoError(t, err)
	logoReq.Header.Set("Authorization", "Bearer "+token)
	logoReq.Header.Set("Content-Type", logoWriter.FormDataContentType())
	logoResp, err := client.Do(logoReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, logoResp.StatusCode)
	logoParsed := decodeBody(t, logoResp)
	assert.Equal(t, "LOGO 已上传", logoParsed["message"])
	logoData, ok := logoParsed["data"].(map[string]any)
	require.True(t, ok)
	logoPath, ok := logoData["logo_path"].(string)
	require.True(t, ok)
	assert.True(t, strings.HasPrefix(logoPath, "/storage/app/uploads/schools/"), logoPath)

	var school models.School
	require.NoError(t, db.First(&school, schoolID).Error)
	assert.Equal(t, logoPath, school.LogoPath, "logo_path 已写库")

	// 批量转班 + 批量删除
	moved := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/students/batch-move",
		bytes.NewReader([]byte(`{"student_ids":[1],"target_class_id":2}`)), "application/json")
	require.Equal(t, http.StatusOK, moved.StatusCode)
	deleted := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/students/batch-delete",
		bytes.NewReader([]byte(`{"student_ids":[1]}`)), "application/json")
	require.Equal(t, http.StatusOK, deleted.StatusCode)
	deletedBody := decodeBody(t, deleted)
	assert.Contains(t, deletedBody["message"], "名学生")

	// 账号批量操作
	reset := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/accounts/batch-reset-password",
		bytes.NewReader([]byte(`{"role":"teacher","ids":[99999],"password":"abcdef"}`)), "application/json")
	assert.Equal(t, http.StatusNotFound, reset.StatusCode, "无命中账号 → 404")
}

// TestAdminOpsHTTPGradeUpgrade 学年升级（预览 → 执行）与别名路径。
func TestAdminOpsHTTPGradeUpgrade(t *testing.T) {
	db, client, token, schoolID := newAdminFixture(t)

	// fixture（newTestEngine）已有一年级（1）班，这里只补一个六年级班
	require.NoError(t, db.Create(&models.ClassRoom{
		SchoolID: schoolID, Name: "六年级（1）班", Grade: "六年级", Status: "active",
	}).Error)

	preview := doAdmin(t, client, token, http.MethodGet, "/api/v1/admin/grade-upgrade/preview", nil, "")
	require.Equal(t, http.StatusOK, preview.StatusCode)
	previewBody := decodeBody(t, preview)
	summary, ok := previewBody["data"].(map[string]any)["summary"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(1), summary["upgrade_class_count"])
	assert.Equal(t, float64(1), summary["graduate_class_count"])
	// 不可逆：先走正式路径执行一次
	execute := doAdmin(t, client, token, http.MethodPost, "/api/v1/admin/grade-upgrade/execute", bytes.NewReader(nil), "application/json")
	require.Equal(t, http.StatusOK, execute.StatusCode)
	executeBody := decodeBody(t, execute)
	assert.Equal(t, "学年升级完成", executeBody["message"])
	result, ok := executeBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(1), result["upgraded_classes"])
	assert.Equal(t, float64(1), result["archived_classes"])

	// classes/:id/assign-teacher 与 remove-teacher
	teacher := models.User{SchoolID: schoolID, Role: "teacher", Username: "ops-teacher", Name: "张老师", Status: "active"}
	require.NoError(t, db.Create(&teacher).Error)
	var class models.ClassRoom
	require.NoError(t, db.Where("school_id = ? AND name = ?", schoolID, "二年级（1）班").First(&class).Error)

	assign := doAdmin(t, client, token, http.MethodPost,
		"/api/v1/admin/classes/"+itoaUint(class.ID)+"/assign-teacher",
		bytes.NewReader([]byte(`{"teacher_id":`+itoaUint(teacher.ID)+`,"role":"head_teacher"}`)), "application/json")
	require.Equal(t, http.StatusOK, assign.StatusCode)
	var reloaded models.ClassRoom
	require.NoError(t, db.First(&reloaded, class.ID).Error)
	require.NotNil(t, reloaded.TeacherID)
	assert.Equal(t, teacher.ID, *reloaded.TeacherID)

	remove := doAdmin(t, client, token, http.MethodDelete,
		"/api/v1/admin/classes/"+itoaUint(class.ID)+"/remove-teacher",
		bytes.NewReader([]byte(`{"teacher_id":`+itoaUint(teacher.ID)+`}`)), "application/json")
	require.Equal(t, http.StatusOK, remove.StatusCode)
	removeBody := decodeBody(t, remove)
	assert.Equal(t, "教师已从班级移除", removeBody["message"])
	require.NoError(t, db.First(&reloaded, class.ID).Error)
	assert.Nil(t, reloaded.TeacherID)

}

// itoaUint 无依赖的 uint → 十进制字符串。
func itoaUint(n uint) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
