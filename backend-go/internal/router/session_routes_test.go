// 本批 15 条路由的路由契约与端到端 HTTP 行为：
// 认证会话（logout / refresh / bindings）、管理端 5 条、教师端 7 条。
//
// 覆盖要点：注册齐全、无 token 一律 401、登出后旧 token 401（另一 token 仍可用）、
// 刷新换新 token 且旧 token 失效、bindings 三种平台状态、学生列表筛选与分页 meta、
// 班级详情含 teacher+students、跨校 404、教师班级分配 replace 与 head_teacher 同步、
// 查看教师密码自动生成并写回、classInfo 成功与 400、教师建学生 403/201、导入统计、
// 更新/删除文案、宠物图鉴、按规则加减分文案与积分落账。
package router_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// sessionRoutes 本批新增的 15 条路由。
var sessionRoutes = []string{
	"POST /api/v1/auth/logout",
	"POST /api/v1/auth/refresh",
	"GET /api/v1/auth/bindings",
	"GET /api/v1/admin/students",
	"POST /api/v1/admin/students",
	"GET /api/v1/admin/classes/:id",
	"PUT /api/v1/admin/teachers/:id/classes",
	"GET /api/v1/admin/teachers/:id/password",
	"GET /api/v1/teacher/class",
	"POST /api/v1/teacher/students",
	"POST /api/v1/teacher/students/import",
	"PUT /api/v1/teacher/students/:id",
	"DELETE /api/v1/teacher/students/:id",
	"GET /api/v1/teacher/pets/:student_id/collection",
	"POST /api/v1/teacher/scores/give-by-rule/:ruleId",
}

// peopleFixture 本批测试夹具：一所学校 + 管理员 + 班主任（含所辖班级与学生）。
type peopleFixture struct {
	DB           *gorm.DB
	Client       *http.Client
	BaseURL      string
	AdminToken   string
	TeacherToken string
	School       models.School
	Class        models.ClassRoom
	Teacher      models.User
	Student      models.Student
}

// newPeopleFixture 自建内存库 + 完整路由（与 newTestEngine 同构，但保留 db 句柄）。
func newPeopleFixture(t *testing.T) *peopleFixture {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), database.GormConfig())
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1) // 内存库需单连接，否则每连接各自独立。
	require.NoError(t, database.Migrate(db))

	school := models.School{Name: "测试学校", Code: "people-school", Status: "active"}
	require.NoError(t, db.Create(&school).Error)

	admin := models.User{
		SchoolID: school.ID, Role: "school_admin", Username: "people-admin",
		Name: "管理员", Status: "active",
	}
	require.NoError(t, db.Create(&admin).Error)

	teacher := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "people-teacher",
		Name: "李老师", Status: "active",
	}
	require.NoError(t, db.Create(&teacher).Error)

	class := models.ClassRoom{
		SchoolID: school.ID, Grade: "一年级", Name: "一年级（1）班",
		Status: "active", TeacherID: &teacher.ID, DisplayCode: "LS11",
	}
	require.NoError(t, db.Create(&class).Error)
	student := models.Student{ClassID: class.ID, Name: "小明", StudentNo: "1", Status: "active"}
	require.NoError(t, db.Create(&student).Error)

	jwtMgr := jwtauth.New("test-secret", 1)
	adminToken, err := jwtMgr.Generate(admin.ID, admin.Role, admin.SchoolID)
	require.NoError(t, err)
	teacherToken, err := jwtMgr.Generate(teacher.ID, teacher.Role, teacher.SchoolID)
	require.NoError(t, err)

	engine := router.New(db, &config.Config{JWTSecret: "test-secret", JWTExpHours: 1})
	server := httptest.NewServer(engine)
	t.Cleanup(server.Close)

	return &peopleFixture{
		DB: db, Client: server.Client(), BaseURL: server.URL,
		AdminToken: adminToken, TeacherToken: teacherToken,
		School: school, Class: class, Teacher: teacher, Student: student,
	}
}

// do 发请求（token 为空表示不带 Authorization）。
func (f *peopleFixture) do(t *testing.T, token, method, path, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, f.BaseURL+path, bytes.NewReader([]byte(body)))
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := f.Client.Do(req)
	require.NoError(t, err)
	return resp
}

// itoa 十进制整数转字符串（拼接路径与请求体用）。
func itoa(n int) string { return strconv.Itoa(n) }

// TestSessionRoutesRegistered 15 条路由全部注册（路径逐字对齐 Laravel）。
func TestSessionRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range sessionRoutes {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}
}

// TestSessionRoutesRequireAuth 全部新路由受鉴权保护（无 token → 401）。
func TestSessionRoutesRequireAuth(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	for _, route := range sessionRoutes {
		method, path := splitRouteKey(t, route)
		path = replaceIDParam(path)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(method, path, bytes.NewReader([]byte(`{}`))))
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s 未登录应 401", route)
	}
}

// TestAuthLogoutAndRefreshInvalidateOldToken 登出/刷新真的让旧 token 失效（另一 token 不受影响）。
func TestAuthLogoutAndRefreshInvalidateOldToken(t *testing.T) {
	f := newPeopleFixture(t)

	// ① 两个独立 token 均可用。
	assert.Equal(t, http.StatusOK, f.do(t, f.AdminToken, http.MethodGet, "/api/v1/auth/bindings", "").StatusCode)
	assert.Equal(t, http.StatusOK, f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/auth/bindings", "").StatusCode)

	// ② 登出管理员 token → 该 token 立即 401，教师 token 仍 200。
	logout := f.do(t, f.AdminToken, http.MethodPost, "/api/v1/auth/logout", "")
	require.Equal(t, http.StatusOK, logout.StatusCode)
	assert.Equal(t, "已登出", decodeBody(t, logout)["message"])

	afterLogout := f.do(t, f.AdminToken, http.MethodGet, "/api/v1/auth/bindings", "")
	require.Equal(t, http.StatusUnauthorized, afterLogout.StatusCode, "登出后旧 token 必须被拒")
	assert.Equal(t, "登录已过期或无效", decodeBody(t, afterLogout)["message"])
	assert.Equal(t, http.StatusOK, f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/auth/bindings", "").StatusCode,
		"另一 token 不应受影响")

	// ③ 刷新：拿到新 token，旧 token 失效，新 token 可用。
	refresh := f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/auth/refresh", "")
	require.Equal(t, http.StatusOK, refresh.StatusCode)
	inner, ok := decodeBody(t, refresh)["data"].(map[string]any)
	require.True(t, ok)
	newToken, ok := inner["token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, newToken)
	require.NotEqual(t, f.TeacherToken, newToken)

	assert.Equal(t, http.StatusUnauthorized,
		f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/auth/bindings", "").StatusCode, "刷新后旧 token 必须被拒")
	assert.Equal(t, http.StatusOK,
		f.do(t, newToken, http.MethodGet, "/api/v1/auth/bindings", "").StatusCode, "新 token 必须可用")
}

// TestAuthBindingsHTTP 三种平台状态：未配置默认三平台、学校交集、已绑定带昵称。
func TestAuthBindingsHTTP(t *testing.T) {
	f := newPeopleFixture(t)

	resp := f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/auth/bindings", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	list, ok := decodeBody(t, resp)["data"].([]any)
	require.True(t, ok)
	require.Len(t, list, 3, "未配置学校开关 → 默认 [企业微信, 微信, QQ]")
	first := list[0].(map[string]any)
	assert.Equal(t, "wechat_work", first["platform"])
	assert.Equal(t, "企业微信", first["label"])
	assert.Equal(t, "🏢", first["icon"])
	assert.Equal(t, false, first["bound"])
	assert.Nil(t, first["nick"])

	// 学校只启用 wechat → 仅 1 条；绑定后 bound=true 且带昵称。
	settings, err := f.School.WithSetting("enabled_third_party_platforms", []string{"wechat"})
	require.NoError(t, err)
	require.NoError(t, f.DB.Model(&models.School{}).Where("id = ?", f.School.ID).
		Update("settings", settings).Error)
	require.NoError(t, f.DB.Create(&models.ThirdPartyBinding{
		UserID: f.Teacher.ID, Platform: "wechat", PlatformID: "wx-1", PlatformNick: "李老师",
	}).Error)

	resp = f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/auth/bindings", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	list, ok = decodeBody(t, resp)["data"].([]any)
	require.True(t, ok)
	require.Len(t, list, 1)
	only := list[0].(map[string]any)
	assert.Equal(t, "wechat", only["platform"])
	assert.Equal(t, true, only["bound"])
	assert.Equal(t, "李老师", only["nick"])
}

// TestAdminStudentsHTTP 学生列表筛选/分页 meta 与单建学生（422 / 201）。
func TestAdminStudentsHTTP(t *testing.T) {
	f := newPeopleFixture(t)
	require.NoError(t, f.DB.Create(&models.Student{
		ClassID: f.Class.ID, Name: "小红", StudentNo: "2", Status: "active",
	}).Error)

	list := f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/students?per_page=1&page=2", "")
	require.Equal(t, http.StatusOK, list.StatusCode)
	body := decodeBody(t, list)
	rows, ok := body["data"].([]any)
	require.True(t, ok)
	require.Len(t, rows, 1)
	meta, ok := body["meta"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(2), meta["total"])
	assert.Equal(t, float64(2), meta["last_page"])
	assert.Equal(t, float64(1), meta["per_page"])
	assert.Equal(t, float64(2), meta["current_page"])
	row := rows[0].(map[string]any)
	assert.Equal(t, f.Class.Name, row["class_name"])
	assert.Equal(t, "一年级", row["class_grade"])
	assert.Contains(t, row, "pet_species")

	// search 命中。
	search := f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/students?search=小红", "")
	require.Equal(t, http.StatusOK, search.StatusCode)
	assert.Equal(t, float64(1), decodeBody(t, search)["meta"].(map[string]any)["total"])

	// 校验失败 → 422「参数错误」+ errors。
	bad := f.do(t, f.AdminToken, http.MethodPost, "/api/v1/admin/students", `{"name":"","class_id":0}`)
	require.Equal(t, http.StatusUnprocessableEntity, bad.StatusCode)
	badBody := decodeBody(t, bad)
	assert.Equal(t, "参数错误", badBody["message"])
	errs, ok := badBody["errors"].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, errs, "name")
	assert.Contains(t, errs, "class_id")

	// 班级不存在（exists 规则）→ 422。
	missing := f.do(t, f.AdminToken, http.MethodPost, "/api/v1/admin/students",
		`{"name":"新同学","class_id":999999}`)
	require.Equal(t, http.StatusUnprocessableEntity, missing.StatusCode)
	assert.Contains(t, decodeBody(t, missing)["errors"].(map[string]any), "class_id")

	// 成功 → 201 + 文案 + 自动分配萌宠。
	created := f.do(t, f.AdminToken, http.MethodPost, "/api/v1/admin/students",
		`{"name":"新同学","class_id":`+itoa(int(f.Class.ID))+`,"gender":"男生","student_no":"9"}`)
	require.Equal(t, http.StatusCreated, created.StatusCode)
	cBody := decodeBody(t, created)
	assert.Equal(t, "学生「新同学」已添加，已自动分配萌宠", cBody["message"])
	data, ok := cBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "男", data["gender"])

	var pet models.Pet
	require.NoError(t, f.DB.Where("student_id = ?", uint(data["id"].(float64))).First(&pet).Error)
	assert.Equal(t, 1, pet.Level)
}

// TestAdminClassShowHTTP 班级详情含 teacher + students；跨校 404。
func TestAdminClassShowHTTP(t *testing.T) {
	f := newPeopleFixture(t)

	resp := f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/classes/"+itoa(int(f.Class.ID)), "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, ok := decodeBody(t, resp)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, f.Class.Name, data["name"])
	teacher, ok := data["teacher"].(map[string]any)
	require.True(t, ok, "应含 teacher 关联")
	assert.Equal(t, f.Teacher.ID, uint(teacher["id"].(float64)))
	students, ok := data["students"].([]any)
	require.True(t, ok, "应含 students 关联")
	require.Len(t, students, 1)
	assert.Equal(t, "小明", students[0].(map[string]any)["name"])

	// 跨校班级 → 404「班级不存在」。
	other := models.School{Name: "别校", Code: "people-other-school", Status: "active"}
	require.NoError(t, f.DB.Create(&other).Error)
	foreign := models.ClassRoom{SchoolID: other.ID, Name: "别校一班", Status: "active"}
	require.NoError(t, f.DB.Create(&foreign).Error)

	cross := f.do(t, f.AdminToken, http.MethodGet, "/api/v1/admin/classes/"+itoa(int(foreign.ID)), "")
	require.Equal(t, http.StatusNotFound, cross.StatusCode)
	assert.Equal(t, "班级不存在", decodeBody(t, cross)["message"])
}

// TestAdminTeacherClassesAndPasswordHTTP 班级分配 replace/head 同步 + 查看教师密码自动生成。
func TestAdminTeacherClassesAndPasswordHTTP(t *testing.T) {
	f := newPeopleFixture(t)
	second := models.ClassRoom{SchoolID: f.School.ID, Name: "一年级（2）班", Grade: "一年级", Status: "active"}
	require.NoError(t, f.DB.Create(&second).Error)

	path := "/api/v1/admin/teachers/" + itoa(int(f.Teacher.ID)) + "/classes"
	assign := f.do(t, f.AdminToken, http.MethodPut, path,
		`{"assignments":[{"class_id":`+itoa(int(f.Class.ID))+`,"role":"head_teacher","subject":"数学"},`+
			`{"class_id":`+itoa(int(second.ID))+`,"role":"co_teacher"}]}`)
	require.Equal(t, http.StatusOK, assign.StatusCode)
	body := decodeBody(t, assign)
	assert.Equal(t, "已为教师「李老师」分配 2 个班级", body["message"])
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, f.Teacher.ID, uint(data["teacher_id"].(float64)))
	assert.Equal(t, "李老师", data["teacher_name"])
	assignments, ok := data["assignments"].([]any)
	require.True(t, ok)
	require.Len(t, assignments, 2)
	first := assignments[0].(map[string]any)
	assert.Equal(t, f.Class.Name, first["class_name"])
	assert.Equal(t, "head_teacher", first["role"])
	assert.Equal(t, "数学", first["subject"])
	assert.Nil(t, assignments[1].(map[string]any)["subject"])

	var refreshed models.ClassRoom
	require.NoError(t, f.DB.First(&refreshed, f.Class.ID).Error)
	require.NotNil(t, refreshed.TeacherID)
	assert.Equal(t, f.Teacher.ID, *refreshed.TeacherID)

	// replace：只提交第二个班 → 第一个班关联被删、teacher_id 置空。
	replace := f.do(t, f.AdminToken, http.MethodPut, path,
		`{"assignments":[{"class_id":`+itoa(int(second.ID))+`,"role":"grade_lead"}]}`)
	require.Equal(t, http.StatusOK, replace.StatusCode)
	assert.Equal(t, "已为教师「李老师」分配 1 个班级", decodeBody(t, replace)["message"])

	var rows []models.ClassRoomTeacher
	require.NoError(t, f.DB.Where("user_id = ?", f.Teacher.ID).Find(&rows).Error)
	require.Len(t, rows, 1)
	assert.Equal(t, second.ID, rows[0].ClassRoomID)

	require.NoError(t, f.DB.First(&refreshed, f.Class.ID).Error)
	assert.Nil(t, refreshed.TeacherID, "被移除的 head_teacher 班级应清空 teacher_id")

	// 非法角色 → 422「参数错误」。
	bad := f.do(t, f.AdminToken, http.MethodPut, path, `{"assignments":[{"class_id":1,"role":"nope"}]}`)
	require.Equal(t, http.StatusUnprocessableEntity, bad.StatusCode)
	assert.Equal(t, "参数错误", decodeBody(t, bad)["message"])

	// 查看密码：无明文 → 自动生成 8 位并写回。
	pw := f.do(t, f.AdminToken, http.MethodGet,
		"/api/v1/admin/teachers/"+itoa(int(f.Teacher.ID))+"/password", "")
	require.Equal(t, http.StatusOK, pw.StatusCode)
	pwData, ok := decodeBody(t, pw)["data"].(map[string]any)
	require.True(t, ok)
	password, ok := pwData["password"].(string)

	// 班级不存在（exists 规则）→ 422 + errors。
	missingClass := f.do(t, f.AdminToken, http.MethodPut, path,
		`{"assignments":[{"class_id":999999,"role":"co_teacher"}]}`)
	require.Equal(t, http.StatusUnprocessableEntity, missingClass.StatusCode)
	assert.Contains(t, decodeBody(t, missingClass)["errors"].(map[string]any), "assignments.0.class_id")
	require.True(t, ok)
	assert.Len(t, password, 8)

	var reloaded models.User
	require.NoError(t, f.DB.First(&reloaded, f.Teacher.ID).Error)
	assert.Equal(t, password, reloaded.PlainPassword)

	// 教师身份访问管理端 → 403（角色隔离）。
	assert.Equal(t, http.StatusForbidden,
		f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/admin/students", "").StatusCode)
}

// TestTeacherStudentsHTTP 教师建学生（403 / 422 / 201）、导入统计、更新与删除文案。
func TestTeacherStudentsHTTP(t *testing.T) {
	f := newPeopleFixture(t)
	otherClass := models.ClassRoom{SchoolID: f.School.ID, Name: "别的班", Grade: "一年级", Status: "active"}
	require.NoError(t, f.DB.Create(&otherClass).Error)

	// 不在自己管辖的班级 → 403「只能在自己管理的班级添加学生」。
	forbidden := f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/teacher/students",
		`{"name":"越权同学","class_id":`+itoa(int(otherClass.ID))+`}`)
	require.Equal(t, http.StatusForbidden, forbidden.StatusCode)
	assert.Equal(t, "只能在自己管理的班级添加学生", decodeBody(t, forbidden)["message"])

	// 校验失败 → 422 + errors。
	bad := f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/teacher/students", `{"name":""}`)
	require.Equal(t, http.StatusUnprocessableEntity, bad.StatusCode)
	badBody := decodeBody(t, bad)
	assert.Equal(t, "参数错误", badBody["message"])
	assert.Contains(t, badBody["errors"].(map[string]any), "class_id")

	// 成功 → 201 + 「学生「x」已添加」。
	created := f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/teacher/students",
		`{"name":"新同学","class_id":`+itoa(int(f.Class.ID))+`,"gender":"女生","student_no":"3"}`)
	require.Equal(t, http.StatusCreated, created.StatusCode)
	cBody := decodeBody(t, created)
	assert.Equal(t, "学生「新同学」已添加", cBody["message"])
	newID := uint(cBody["data"].(map[string]any)["id"].(float64))

	// 导入：1 条成功 + 1 条重复跳过（沿用既有学生学号 1）。
	imp := f.do(t, f.TeacherToken, http.MethodPost, "/api/v1/teacher/students/import",
		`{"students":[{"name":"导入同学","class_name":"`+f.Class.Name+`","gender":"男生","student_no":"4"},`+
			`{"name":"小明","class_name":"`+f.Class.Name+`","student_no":"1"}]}`)
	require.Equal(t, http.StatusOK, imp.StatusCode)
	impBody := decodeBody(t, imp)
	assert.Equal(t, "成功导入 1 名学生，跳过 1 条重复/冲突记录", impBody["message"])
	impData, ok := impBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(1), impData["imported_count"])
	skipped, ok := impData["skipped"].([]any)
	require.True(t, ok)
	require.Len(t, skipped, 1)
	assert.Contains(t, skipped[0].(string), "已存在")

	// 更新 → 「更新成功」。
	updated := f.do(t, f.TeacherToken, http.MethodPut, "/api/v1/teacher/students/"+itoa(int(newID)),
		`{"name":"改名同学","student_no":"5"}`)
	require.Equal(t, http.StatusOK, updated.StatusCode)
	uBody := decodeBody(t, updated)
	assert.Equal(t, "更新成功", uBody["message"])
	assert.Equal(t, "改名同学", uBody["data"].(map[string]any)["name"])
	assert.Equal(t, "5", uBody["data"].(map[string]any)["student_no"])

	// 越权更新 → 404。
	cross := f.do(t, f.TeacherToken, http.MethodPut, "/api/v1/teacher/students/999999", `{"name":"x"}`)
	require.Equal(t, http.StatusNotFound, cross.StatusCode)

	// 删除 → 「学生「x」已删除」（硬删除）。
	deleted := f.do(t, f.TeacherToken, http.MethodDelete, "/api/v1/teacher/students/"+itoa(int(newID)), "")
	require.Equal(t, http.StatusOK, deleted.StatusCode)
	assert.Equal(t, "学生「改名同学」已删除", decodeBody(t, deleted)["message"])
	var count int64
	require.NoError(t, f.DB.Model(&models.Student{}).Where("id = ?", newID).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

// TestTeacherClassInfoCollectionAndGiveByRuleHTTP 班级信息卡 / 宠物图鉴 / 按规则加减分。
func TestTeacherClassInfoCollectionAndGiveByRuleHTTP(t *testing.T) {
	f := newPeopleFixture(t)

	info := f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/teacher/class", "")
	require.Equal(t, http.StatusOK, info.StatusCode)
	data, ok := decodeBody(t, info)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, f.Class.ID, uint(data["id"].(float64)))
	assert.Equal(t, f.Class.Name, data["name"])
	assert.Equal(t, float64(1), data["student_count"])
	assert.Equal(t, "LS11", data["display_code"])
	assert.Nil(t, data["settings"], "空 settings → null（同 Laravel）")

	// 没有可管辖班级的教师 → 400「没有可管理的班级」。
	lonely := models.User{
		SchoolID: f.School.ID, Role: "teacher", Username: "lonely-teacher", Name: "无班老师", Status: "active",
	}
	require.NoError(t, f.DB.Create(&lonely).Error)
	lonelyToken, err := jwtauth.New("test-secret", 1).Generate(lonely.ID, lonely.Role, lonely.SchoolID)
	require.NoError(t, err)
	noClass := f.do(t, lonelyToken, http.MethodGet, "/api/v1/teacher/class", "")
	require.Equal(t, http.StatusBadRequest, noClass.StatusCode)
	assert.Equal(t, "没有可管理的班级", decodeBody(t, noClass)["message"])

	// 宠物图鉴：先给小明一只宠物 → 被补录进图鉴。
	require.NoError(t, f.DB.Create(&models.Pet{
		StudentID: f.Student.ID, ClassID: f.Class.ID, Name: "小明的萌宠",
		Species: "zhulong", Level: 2, Experience: 4, Mood: 80,
	}).Error)
	collection := f.do(t, f.TeacherToken, http.MethodGet,
		"/api/v1/teacher/pets/"+itoa(int(f.Student.ID))+"/collection", "")
	require.Equal(t, http.StatusOK, collection.StatusCode)
	cData, ok := decodeBody(t, collection)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, f.Student.ID, uint(cData["student_id"].(float64)))
	assert.Equal(t, "小明", cData["student_name"])
	assert.Equal(t, float64(1), cData["unlock_slots"])
	assert.Equal(t, "zhulong", cData["active_species"])
	entries, ok := cData["collection"].([]any)
	require.True(t, ok)
	require.Len(t, entries, 1)
	assert.Equal(t, "zhulong", entries[0].(map[string]any)["species"])

	// 越权学生 → 404。
	missing := f.do(t, f.TeacherToken, http.MethodGet, "/api/v1/teacher/pets/999999/collection", "")
	require.Equal(t, http.StatusNotFound, missing.StatusCode)

	// 按规则加减分：文案 + 落账 + score_rule_id。
	rule := models.ScoreRule{
		SchoolID: &f.School.ID, Name: "举手发言", Amount: 3,
		Category: "classroom", IsPositive: true, IsActive: true,
	}
	require.NoError(t, f.DB.Create(&rule).Error)

	give := f.do(t, f.TeacherToken, http.MethodPost,
		"/api/v1/teacher/scores/give-by-rule/"+itoa(int(rule.ID)),
		`{"student_id":`+itoa(int(f.Student.ID))+`}`)
	require.Equal(t, http.StatusOK, give.StatusCode)
	gBody := decodeBody(t, give)
	assert.Equal(t, "已按规则「举手发言」处理", gBody["message"])
	gData, ok := gBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "举手发言", gData["rule"])
	assert.Equal(t, float64(3), gData["points"])
	assert.Equal(t, "小明", gData["student_name"])

	var reloaded models.Student
	require.NoError(t, f.DB.First(&reloaded, f.Student.ID).Error)
	assert.Equal(t, 3, reloaded.TotalScore)
	var score models.Score
	require.NoError(t, f.DB.Where("student_id = ?", f.Student.ID).Order("id DESC").First(&score).Error)
	assert.Equal(t, 3, score.Amount)
	assert.Equal(t, "举手发言", score.Reason)
	require.NotNil(t, score.ScoreRuleID)
	assert.Equal(t, rule.ID, *score.ScoreRuleID)

	// 缺 student_id → 422「参数错误」+ errors；不存在规则 → 404。
	noStudent := f.do(t, f.TeacherToken, http.MethodPost,
		"/api/v1/teacher/scores/give-by-rule/"+itoa(int(rule.ID)), `{}`)
	require.Equal(t, http.StatusUnprocessableEntity, noStudent.StatusCode)
	assert.Contains(t, decodeBody(t, noStudent)["errors"].(map[string]any), "student_id")

	unknownRule := f.do(t, f.TeacherToken, http.MethodPost,
		"/api/v1/teacher/scores/give-by-rule/999999", `{"student_id":`+itoa(int(f.Student.ID))+`}`)
	require.Equal(t, http.StatusNotFound, unknownRule.StatusCode)
}
