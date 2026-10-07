// 教师端报表导出（CSV 版）测试：GET /api/v1/teacher/reports/export/{type}。
//
// 覆盖：三种报表（scores / pets / attendance）的表头列名与行数、数据映射（造 3 名学生 + 积分 + 宠物 + 考勤）、
// UTF-8 BOM、Content-Type / Content-Disposition（中文名 rawurlencode）、未知 type 文案、
// 跨班 / 跨校越权 403，以及 HTTP 层的响应头与文件名。
package services_test

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/handlers"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newExportEngine 构造带 JWT 鉴权的导出路由引擎（教师端报表导出 + 管理端课表导出）。
func newExportEngine(t *testing.T, db *gorm.DB) (*gin.Engine, *auth.Manager) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	jwtMgr := auth.New("test-secret", 72)
	h := handlers.New(db, jwtMgr)

	engine := gin.New()
	teacher := engine.Group("/api/v1/teacher", middleware.Auth(jwtMgr, db), middleware.RequireRole("teacher"))
	teacher.GET("/reports/export/:type", h.TeacherReportExport)

	admin := engine.Group("/api/v1/admin", middleware.Auth(jwtMgr, db), middleware.RequireRole("school_admin"))
	admin.GET("/classes/:id/timetable/export-excel", h.AdminClassTimetableExportExcel)
	admin.GET("/timetable/export-excel", h.AdminTimetableExportExcel)

	return engine, jwtMgr
}

// doExport 带 Bearer token 发一次 GET。
func doExport(t *testing.T, engine *gin.Engine, jwtMgr *auth.Manager, user *models.User, target string) *httptest.ResponseRecorder {
	t.Helper()
	token, err := jwtMgr.Generate(user.ID, user.Role, user.SchoolID)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

// parseCSV 解析带 BOM 的 CSV 内容（返回所有记录，列数不固定）。
func parseCSV(t *testing.T, content []byte) [][]string {
	t.Helper()
	require.True(t, strings.HasPrefix(string(content), "\ufeff"), "CSV 必须以 UTF-8 BOM 开头")

	reader := csv.NewReader(strings.NewReader(strings.TrimPrefix(string(content), "\ufeff")))
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	require.NoError(t, err)
	return rows
}

// exportFixture 报表导出夹具。
type exportFixture struct {
	School  models.School
	Teacher models.User
	Class   models.ClassRoom
	Ming    models.Student // 有宠物 + 正负积分
	Hong    models.Student // 有宠物 + 仅扣分
	Gang    models.Student // 无宠物 + 无积分
}

// newExportFixture 造学校 / 教师 / 班级（3 名 active 学生 + 1 名 inactive）+ 宠物 + 积分 + 考勤。
func newExportFixture(t *testing.T, db *gorm.DB) exportFixture {
	t.Helper()

	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "export-teacher")
	class, ming := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	hong := models.Student{ClassID: class.ID, Name: "小红", StudentNo: "002", Status: "active", TotalScore: 3}
	require.NoError(t, db.Create(&hong).Error)
	gang := models.Student{ClassID: class.ID, Name: "小刚", StudentNo: "003", Status: "active", TotalScore: 4}
	require.NoError(t, db.Create(&gang).Error)
	// 非 active 学生不出现在导出中。
	require.NoError(t, db.Create(&models.Student{ClassID: class.ID, Name: "休学", StudentNo: "004", Status: "inactive"}).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", ming.ID).Update("total_score", 10).Error)

	// 宠物：小明 Lv.5（成长期）/ 小红 Lv.1（新生之卵）；小刚无宠物。
	require.NoError(t, db.Create(&models.Pet{
		StudentID: ming.ID, ClassID: class.ID, Name: "阿龙", Species: "zhulong", Level: 5, Experience: 12,
	}).Error)
	require.NoError(t, db.Create(&models.Pet{
		StudentID: hong.ID, ClassID: class.ID, Name: "小鸽", Species: "pigeon", Level: 1, Experience: 0,
	}).Error)

	// 积分：小明 +5/+3/-2 → 获得 8 / 扣除 2；小红 -4 → 获得 0 / 扣除 4；小刚无积分。
	now := util.Now()
	makeScore(t, db, class.ID, ming.ID, 5, now)
	makeScore(t, db, class.ID, ming.ID, 3, now)
	makeScore(t, db, class.ID, ming.ID, -2, now)
	makeScore(t, db, class.ID, hong.ID, -4, now)

	// 考勤：昨天 小明=缺席；今天 小明=出勤(auto) / 小红=迟到(manual)；小刚无记录。
	yesterday := util.StartOfDay(now).AddDate(0, 0, -1).Add(8 * time.Hour)
	signIn := util.StartOfDay(now).Add(7*time.Hour + 50*time.Minute)
	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: ming.ID, Date: yesterday.Format("2006-01-02"),
		Status: "absent", Source: "auto", CreatedAt: yesterday,
	}).Error)
	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: ming.ID, Date: util.Today(), TeacherID: &teacher.ID,
		Status: "present", Source: "auto", SignInAt: &signIn,
	}).Error)
	require.NoError(t, db.Create(&models.Attendance{
		ClassID: class.ID, StudentID: hong.ID, Date: util.Today(), TeacherID: &teacher.ID,
		Status: "late", Source: "manual",
	}).Error)

	return exportFixture{School: school, Teacher: teacher, Class: class, Ming: ming, Hong: hong, Gang: gang}
}

// 三种报表：表头列名、行数与数据映射；BOM；文件名主体。
func TestReportExportCSVRows(t *testing.T) {
	db := setupDB(t)
	f := newExportFixture(t, db)

	t.Run("scores", func(t *testing.T) {
		file, err := f.Svc(t, db).ExportReport(&f.Teacher, "scores", "", false, "")
		require.NoError(t, err)
		require.NotNil(t, file)

		assert.Regexp(t, `^一班-\d{8}-\d{6}-积分报表\.csv$`, file.Filename)

		rows := parseCSV(t, file.Content)
		require.Len(t, rows, 4, "表头 + 3 名 active 学生")
		assert.Equal(t, []string{"姓名", "学号", "总积分", "获得积分", "扣除积分", "宠物名", "宠物等级"}, rows[0])
		assert.Equal(t, []string{"小明", "001", "10", "8", "2", "阿龙", "5"}, rows[1])
		assert.Equal(t, []string{"小红", "002", "3", "0", "4", "小鸽", "1"}, rows[2])
		assert.Equal(t, []string{"小刚", "003", "4", "0", "0", "", "0"}, rows[3], "无宠物时宠物名空串、等级 0")
	})

	t.Run("pets", func(t *testing.T) {
		file, err := f.Svc(t, db).ExportReport(&f.Teacher, "pets", "", false, "")
		require.NoError(t, err)
		require.NotNil(t, file)

		assert.Regexp(t, `^一班-\d{8}-\d{6}-宠物报表\.csv$`, file.Filename)

		rows := parseCSV(t, file.Content)
		require.Len(t, rows, 4)
		assert.Equal(t, []string{"学生姓名", "宠物名", "宠物系列", "等级", "进化阶段", "经验值"}, rows[0])
		assert.Equal(t, []string{"小明", "阿龙", "zhulong", "5", "成长期", "12"}, rows[1])
		assert.Equal(t, []string{"小红", "小鸽", "pigeon", "1", "新生之卵", "0"}, rows[2])
		assert.Equal(t, []string{"小刚", "无", "", "0", "未孵化", "0"}, rows[3], "无宠物时宠物名「无」、阶段「未孵化」")
	})

	t.Run("attendance", func(t *testing.T) {
		file, err := f.Svc(t, db).ExportReport(&f.Teacher, "attendance", "", false, util.Today())
		require.NoError(t, err)
		require.NotNil(t, file)

		assert.Regexp(t, `^一班-\d{8}-\d{6}-考勤报表\.csv$`, file.Filename)

		rows := parseCSV(t, file.Content)
		require.Len(t, rows, 4)
		assert.Equal(t, []string{"姓名", "状态", "签到时间", "来源"}, rows[0])
		// 签到时间列逐字照抄 Laravel 读取的不存在的属性 check_in_time → 恒为空串。
		assert.Equal(t, []string{"小明", "出勤", "", "auto"}, rows[1])
		assert.Equal(t, []string{"小红", "迟到", "", "manual"}, rows[2])
		assert.Equal(t, []string{"小刚", "未记录", "", ""}, rows[3], "无考勤记录 → 未记录 / 空来源")
	})

	t.Run("attendance date filter", func(t *testing.T) {
		svc := f.Svc(t, db)
		yesterday := util.StartOfDay(util.Now()).AddDate(0, 0, -1).Format("2006-01-02")
		file, err := svc.ExportReport(&f.Teacher, "attendance", "", false, yesterday)
		require.NoError(t, err)
		require.NotNil(t, file)

		rows := parseCSV(t, file.Content)
		require.Len(t, rows, 4)
		assert.Equal(t, []string{"小明", "缺席", "", "auto"}, rows[1], "按 date 过滤到昨天 → 缺席")
		assert.Equal(t, []string{"小红", "未记录", "", ""}, rows[2], "昨天小红无记录")

		// 不传 date 时取该生最后一条记录（id 最大者 = 今天的出勤）。
		file, err = svc.ExportReport(&f.Teacher, "attendance", "", false, "")
		require.NoError(t, err)
		rows = parseCSV(t, file.Content)
		assert.Equal(t, []string{"小明", "出勤", "", "auto"}, rows[1])
	})
}

// 未知 type → (nil, nil)；控制器返回 200 + 「导出类型 X 不支持，可选: ...」文案。
func TestReportExportUnsupportedType(t *testing.T) {
	db := setupDB(t)
	f := newExportFixture(t, db)

	file, err := f.Svc(t, db).ExportReport(&f.Teacher, "unknown", "", false, "")
	require.NoError(t, err)
	assert.Nil(t, file)

	engine, jwtMgr := newExportEngine(t, db)
	w := doExport(t, engine, jwtMgr, &f.Teacher, "/api/v1/teacher/reports/export/unknown")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data":null,"message":"导出类型 unknown 不支持，可选: scores, pets, attendance"}`, w.Body.String())
}

// 越权：跨班（同校非班主任）/ 跨校 / class_id 非数字 → 403「无权限」；无班级 → 403。
func TestReportExportScopeGuards(t *testing.T) {
	db := setupDB(t)
	f := newExportFixture(t, db)
	svc := f.Svc(t, db)

	foreignSameSchool := models.ClassRoom{SchoolID: f.School.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&foreignSameSchool).Error)

	otherSchool := seedOtherSchool(t, db)
	foreignOtherSchool := models.ClassRoom{SchoolID: otherSchool.ID, Name: "他校一班", Status: "active"}
	require.NoError(t, db.Create(&foreignOtherSchool).Error)

	_, err := svc.ExportReport(&f.Teacher, "scores", strconv.Itoa(int(foreignSameSchool.ID)), true, "")
	assert.Equal(t, "无权限", appErrorOf(t, err, http.StatusForbidden))

	_, err = svc.ExportReport(&f.Teacher, "scores", strconv.Itoa(int(foreignOtherSchool.ID)), true, "")
	assert.Equal(t, "无权限", appErrorOf(t, err, http.StatusForbidden))

	// PHP `(int) "abc" === 0` → 不在管辖范围 → 403
	_, err = svc.ExportReport(&f.Teacher, "scores", "abc", true, "")
	assert.Equal(t, "无权限", appErrorOf(t, err, http.StatusForbidden))

	// 无管辖班级且未传 class_id → classId = 0 → 403
	lonely := seedTeacher(t, db, f.School.ID, "lonely-teacher")
	_, err = svc.ExportReport(&lonely, "scores", "", false, "")
	assert.Equal(t, "无权限", appErrorOf(t, err, http.StatusForbidden))

	// 显式指定本班 → 放行
	file, err := svc.ExportReport(&f.Teacher, "scores", strconv.Itoa(int(f.Class.ID)), true, "")
	require.NoError(t, err)
	require.NotNil(t, file)

	// HTTP 层：越权 403 文案
	engine, jwtMgr := newExportEngine(t, db)
	w := doExport(t, engine, jwtMgr, &f.Teacher,
		"/api/v1/teacher/reports/export/scores?class_id="+strconv.Itoa(int(foreignSameSchool.ID)))
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{"message":"无权限"}`, w.Body.String())
}

// HTTP 层：响应头（Content-Type / Content-Disposition 中文 rawurlencode）、BOM 与内容。
func TestReportExportHTTPHeaders(t *testing.T) {
	db := setupDB(t)
	f := newExportFixture(t, db)
	engine, jwtMgr := newExportEngine(t, db)

	w := doExport(t, engine, jwtMgr, &f.Teacher,
		"/api/v1/teacher/reports/export/scores?class_id="+strconv.Itoa(int(f.Class.ID)))
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	assert.Equal(t, "text/csv; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Regexp(t,
		`^attachment; filename="%E4%B8%80%E7%8F%AD-\d{8}-\d{6}-%E7%A7%AF%E5%88%86%E6%8A%A5%E8%A1%A8\.csv"$`,
		w.Header().Get("Content-Disposition"))

	rows := parseCSV(t, w.Body.Bytes())
	require.Len(t, rows, 4)
	assert.Equal(t, "小明", rows[1][0])

	// 未登录 → 401（路由受 JWT 保护）
	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/reports/export/scores", nil)
	plain := httptest.NewRecorder()
	engine.ServeHTTP(plain, req)
	assert.Equal(t, http.StatusUnauthorized, plain.Code)
}

// Svc 报表服务（复用既有 newReportService）。
func (f exportFixture) Svc(_ *testing.T, db *gorm.DB) *services.ReportService {
	return newReportService(db)
}
