// 课表网格导出（CSV 版）测试：
//
//	GET /api/v1/admin/classes/:id/timetable/export-excel（单班）
//	GET /api/v1/admin/timetable/export-excel            （全校）
//
// 覆盖：网格行结构（标题行 / 空行 / 表头「节次 · 时间 · 星期X」/ 每节一行）、单元格文本
// （科目（单周）、「教师 · 教室」元信息行、多条排课换行）、无排课时仍可导出、
// 全校多班顺序拼接与班级间空行、归属校验 404、响应头与文件名（中文 rawurlencode）、BOM。
//
// 注意：encoding/csv 读回时会跳过空行，故「空行占位」用 CSV 原文断言（\n\n），
// 解析后的行号里不含空行。
package services_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// seedTimetableExport 建学校 + 两个班级（一班一年级 / 二班二年级）+ 两节课 + 三处排课。
//
// 课表内容（Laravel excelGrid 的期望网格）：
//
//	星期一 第1节：语文 + 「李老师 · 101」
//	星期三 第2节：数学（单周）
//	星期六 第1节：数学      → 上课日列数扩到 6
func seedTimetableExport(t *testing.T, db *gorm.DB) (models.School, models.ClassRoom, models.ClassRoom, *services.TimetableService) {
	t.Helper()

	school := seedSchool(t, db)
	classA := makeClass(t, db, school.ID, "一班", "一年级")
	classB := makeClass(t, db, school.ID, "二班", "二年级")

	svc := newTimetableService(db)
	_, err := svc.Save(classA.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}, {Name: "数学"}},
		Periods: []services.TimetablePeriodInput{
			{PeriodIndex: 1, Name: "第一节", StartTime: "08:00", EndTime: "08:45"},
			{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
		},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("李老师"), Room: strPtr("101")},
			{Weekday: 3, PeriodIndex: 2, SubjectName: "数学", WeekType: "odd"},
			{Weekday: 6, PeriodIndex: 1, SubjectName: "数学"},
		},
	})
	require.NoError(t, err)

	return school, classA, classB, svc
}

// 单班网格：标题行 / 空行 / 表头（按最大星期补到 6）/ 每节一行 + 单元格文本。
func TestTimetableExportClassGrid(t *testing.T) {
	db := setupDB(t)
	school, classA, _, svc := seedTimetableExport(t, db)

	grid, err := svc.ExcelGrid(classA.ID, school.ID)
	require.NoError(t, err)
	assert.Equal(t, "一班", grid.ClassName)

	require.Len(t, grid.Rows, 5, "标题 + 空行 + 表头 + 2 节")
	assert.Equal(t, []string{"一班 课表"}, grid.Rows[0])
	assert.Empty(t, grid.Rows[1], "第 2 行为空行")
	assert.Equal(t, []string{"节次", "时间", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}, grid.Rows[2])
	assert.Equal(t, []string{
		"第1节", "08:00-08:45", "语文\n李老师 · 101", "", "", "", "", "数学",
	}, grid.Rows[3])
	assert.Equal(t, []string{
		"第2节", "08:55-09:40", "", "", "数学（单周）", "", "", "",
	}, grid.Rows[4])

	// 导出为 CSV：BOM + 文件名 + 单元格换行按 CSV 引号规则封装后仍能原样读回。
	file, err := svc.ExportClassFile(classA.ID, school.ID)
	require.NoError(t, err)
	require.NotNil(t, file)
	assert.Equal(t, "一班-课表.csv", file.Filename)

	// 原文：BOM 开头 + 标题行后紧跟一个空行（空行占位，Excel 里表现为空行）。
	text := strings.TrimPrefix(string(file.Content), "\ufeff")
	assert.True(t, strings.HasPrefix(text, "一班 课表\n\n节次,时间,星期一"), "标题行后应有空行：%q", text)

	rows := parseCSV(t, file.Content)
	require.Len(t, rows, 4, "encoding/csv 读回时跳过空行")
	assert.Equal(t, grid.Rows[0], rows[0])
	assert.Equal(t, grid.Rows[2], rows[1])
	assert.Equal(t, "语文\n李老师 · 101", rows[2][2])
	assert.Equal(t, "数学（单周）", rows[3][4])
}

// 无排课也仍能导出：同校仅节次无排课 → 表头 + 空单元格；全校无任何节次 → 仅标题 + 空行 + 表头。
func TestTimetableExportEmptyTimetableStillExports(t *testing.T) {
	db := setupDB(t)
	school, _, emptyClass, svc := seedTimetableExport(t, db)

	// 二班没有任何排课，但学校有节次 → 仍导出，节次行单元格全空。
	grid, err := svc.ExcelGrid(emptyClass.ID, school.ID)
	require.NoError(t, err)
	assert.Equal(t, "二班", grid.ClassName)
	require.Len(t, grid.Rows, 5, "标题 + 空行 + 表头 + 2 节（无排课）")
	assert.Equal(t, []string{"节次", "时间", "星期一", "星期二", "星期三", "星期四", "星期五"}, grid.Rows[2],
		"无排课时上课日列数为 5")
	assert.Equal(t, []string{"第1节", "08:00-08:45", "", "", "", "", ""}, grid.Rows[3])

	file, err := svc.ExportClassFile(emptyClass.ID, school.ID)
	require.NoError(t, err)
	require.NotNil(t, file)
	assert.Equal(t, "二班-课表.csv", file.Filename)
	rows := parseCSV(t, file.Content)
	require.Len(t, rows, 4)
	assert.Equal(t, "二班 课表", rows[0][0])
	assert.Equal(t, "第2节", rows[3][0])

	// 完全没有任何节次的学校：只有标题 + 空行 + 表头。
	freshSchool := seedOtherSchool(t, db)
	freshClass := makeClass(t, db, freshSchool.ID, "三班", "一年级")
	bare, err := svc.ExportClassFile(freshClass.ID, freshSchool.ID)
	require.NoError(t, err)
	require.NotNil(t, bare)

	bareRows := parseCSV(t, bare.Content)
	require.Len(t, bareRows, 2, "标题 + 表头（空行被 csv.Reader 跳过）")
	assert.Equal(t, []string{"三班 课表"}, bareRows[0])
	assert.Equal(t, []string{"节次", "时间", "星期一", "星期二", "星期三", "星期四", "星期五"}, bareRows[1])
}

// 全校导出：按 grade ASC, name ASC 顺序拼接各班网格（班级之间一个空行），文件名 全校课表.csv。
func TestTimetableExportSchoolCSV(t *testing.T) {
	db := setupDB(t)
	school, _, _, svc := seedTimetableExport(t, db)

	file, err := svc.ExportSchoolFile(school.ID)
	require.NoError(t, err)
	require.NotNil(t, file)
	assert.Equal(t, "全校课表.csv", file.Filename)

	text := strings.TrimPrefix(string(file.Content), "\ufeff")
	assert.True(t, strings.Contains(text, "\n\n二班 课表"), "两班之间应有一个空行分隔")

	rows := parseCSV(t, file.Content)
	require.Len(t, rows, 8, "一班 4 行 + 二班 4 行（空行被 csv.Reader 跳过）")
	assert.Equal(t, []string{"一班 课表"}, rows[0])
	assert.Equal(t, []string{"第1节", "08:00-08:45", "语文\n李老师 · 101", "", "", "", "", "数学"}, rows[2])
	assert.Equal(t, []string{"二班 课表"}, rows[4], "第 5 行起为二班（顺序：一年级 → 二年级）")
	assert.Equal(t, []string{"第1节", "08:00-08:45", "", "", "", "", ""}, rows[6], "二班无排课")

	// 无班级的学校 → (nil, nil)，控制器 404「暂无班级可导出」。
	emptySchool := seedOtherSchool(t, db)
	none, err := svc.ExportSchoolFile(emptySchool.ID)
	require.NoError(t, err)
	assert.Nil(t, none)

	// 有班级但未建节次的学校也能导出（本校自己的一套网格）。
	freshSchool := seedExtraSchool(t, db, "fresh-export-school")
	makeClass(t, db, freshSchool.ID, "三班", "一年级")
	freshFile, err := svc.ExportSchoolFile(freshSchool.ID)
	require.NoError(t, err)
	require.NotNil(t, freshFile)
	assert.Equal(t, "全校课表.csv", freshFile.Filename)
	assert.Len(t, parseCSV(t, freshFile.Content), 2)
}

// 归属校验：跨校班级 → 404「班级不存在」。
func TestTimetableExportClassOwnership(t *testing.T) {
	db := setupDB(t)
	school, _, _, svc := seedTimetableExport(t, db)

	otherSchool := seedOtherSchool(t, db)
	foreign := makeClass(t, db, otherSchool.ID, "他校一班", "一年级")

	file, err := svc.ExportClassFile(foreign.ID, school.ID)
	require.Nil(t, file)
	assert.Equal(t, "班级不存在", appErrorOf(t, err, http.StatusNotFound))
}

// HTTP 层：管理员令牌 → 200 + 响应头 + 文件名；跨校 404；全校导出；无班级 404。
func TestTimetableExportHTTP(t *testing.T) {
	db := setupDB(t)
	school, classA, _, _ := seedTimetableExport(t, db)

	admin := models.User{SchoolID: school.ID, Role: "school_admin", Username: "timetable-admin", Name: "管理员", Status: "active"}
	require.NoError(t, db.Create(&admin).Error)

	engine, jwtMgr := newExportEngine(t, db)

	w := doExport(t, engine, jwtMgr, &admin,
		"/api/v1/admin/classes/"+strconv.Itoa(int(classA.ID))+"/timetable/export-excel")
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "text/csv; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Equal(t, `attachment; filename="%E4%B8%80%E7%8F%AD-%E8%AF%BE%E8%A1%A8.csv"`,
		w.Header().Get("Content-Disposition"))
	rows := parseCSV(t, w.Body.Bytes())
	require.Len(t, rows, 4)
	assert.Equal(t, "语文\n李老师 · 101", rows[2][2])

	// 跨校班级 → 404「班级不存在」
	otherSchool := seedOtherSchool(t, db)
	foreign := makeClass(t, db, otherSchool.ID, "他校一班", "一年级")
	w = doExport(t, engine, jwtMgr, &admin,
		"/api/v1/admin/classes/"+strconv.Itoa(int(foreign.ID))+"/timetable/export-excel")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message":"班级不存在"}`, w.Body.String())

	// 全校导出
	w = doExport(t, engine, jwtMgr, &admin, "/api/v1/admin/timetable/export-excel")
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, `attachment; filename="%E5%85%A8%E6%A0%A1%E8%AF%BE%E8%A1%A8.csv"`,
		w.Header().Get("Content-Disposition"))
	assert.Len(t, parseCSV(t, w.Body.Bytes()), 8)

	// 无班级的管理员 → 404「暂无班级可导出」
	emptySchool := seedExtraSchool(t, db, "empty-export-school")
	lonelyAdmin := models.User{SchoolID: emptySchool.ID, Role: "school_admin", Username: "empty-admin", Name: "管理员2", Status: "active"}
	require.NoError(t, db.Create(&lonelyAdmin).Error)
	w = doExport(t, engine, jwtMgr, &lonelyAdmin, "/api/v1/admin/timetable/export-excel")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message":"暂无班级可导出"}`, w.Body.String())

	// 教师令牌访问管理端导出 → 403（RequireRole）
	teacher := seedTeacher(t, db, school.ID, "timetable-teacher")
	w = doExport(t, engine, jwtMgr, &teacher, "/api/v1/admin/timetable/export-excel")
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// seedExtraSchool 建一所 code 唯一的额外学校（seedOtherSchool 的 code 固定，不可重复调用）。
func seedExtraSchool(t *testing.T, db *gorm.DB, code string) models.School {
	t.Helper()
	school := models.School{Name: "测试学校-" + code, Code: code, Status: "active"}
	require.NoError(t, db.Create(&school).Error)
	return school
}
