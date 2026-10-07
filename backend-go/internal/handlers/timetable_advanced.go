// 课表进阶处理器：CSV 批量导入、教师不可用时段、冲突检查、任课设置。
// 路由路径逐字对齐 Laravel backend/routes/api.php 第 98-99、145、150-152 行。
package handlers

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// maxTimetableImportBytes 课表 CSV 上传上限（同 Laravel max:10240，单位 KB）。
const maxTimetableImportBytes = 10240 * 1024

// AdminTimetableImportCsv 管理员批量导入课表（CSV，dry_run 预览；缺省 dry_run=true）。
func (h *Handlers) AdminTimetableImportCsv(c *gin.Context) {
	u := middleware.CurrentUser(c)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fail(c, services.ErrUnprocessable("参数错误：请上传文件（字段名 file）"))
		return
	}
	if fileHeader.Size > maxTimetableImportBytes {
		fail(c, services.ErrUnprocessable("参数错误：文件不能超过 10MB"))
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		fail(c, services.ErrUnprocessable("参数错误：文件读取失败"))
		return
	}
	defer func() { _ = f.Close() }()

	content, err := io.ReadAll(io.LimitReader(f, maxTimetableImportBytes+1))
	if err != nil {
		fail(c, services.ErrUnprocessable("参数错误：文件读取失败"))
		return
	}

	// Laravel `$request->boolean('dry_run', true)`：键缺失时为 true；空串解析为 false。
	dryRun := true
	if v, exists := c.GetPostForm("dry_run"); exists {
		dryRun = parseBoolLoose(v)
	}

	summary, err := h.timetable.ImportFromCsv(content, u.SchoolID, dryRun)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, summary)
}

// AdminTimetableUnavailabilities 全校教师不可用时段列表。
func (h *Handlers) AdminTimetableUnavailabilities(c *gin.Context) {
	u := middleware.CurrentUser(c)

	rows, err := h.timetable.ListUnavailabilities(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// adminUnavailabilitiesRequest 保存不可用时段的请求体。
type adminUnavailabilitiesRequest struct {
	TeacherName string                        `json:"teacher_name"`
	Cells       []services.UnavailabilityCell `json:"cells"`
}

// AdminTimetableSaveUnavailabilities 整体保存某教师的不可用时段（replace 语义）。
func (h *Handlers) AdminTimetableSaveUnavailabilities(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req adminUnavailabilitiesRequest
	if !bindJSON(c, &req) {
		return
	}
	req.TeacherName = trimSpace(req.TeacherName)
	if req.TeacherName == "" {
		fail(c, services.ErrUnprocessable("teacher_name 必填"))
		return
	}
	if runeLen(req.TeacherName) > 50 {
		fail(c, services.ErrUnprocessable("teacher_name 不能超过 50 字"))
		return
	}
	for _, cell := range req.Cells {
		if cell.Weekday < 1 || cell.Weekday > 7 {
			fail(c, services.ErrUnprocessable("cells.weekday 需在 1-7 之间"))
			return
		}
		if cell.PeriodIndex < 1 || cell.PeriodIndex > 30 {
			fail(c, services.ErrUnprocessable("cells.period_index 需在 1-30 之间"))
			return
		}
	}

	if err := h.timetable.SaveUnavailabilities(u.SchoolID, req.TeacherName, req.Cells); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "不可用时段已保存")
}

// adminConflictRequest 冲突检查请求体。
type adminConflictRequest struct {
	ClassID uint                              `json:"class_id"`
	Entries []services.TimetableConflictEntry `json:"entries"`
}

// AdminTimetableCheckConflicts 冲突检查（某班编辑中的排课 vs 其他班级与不可用时段）。
func (h *Handlers) AdminTimetableCheckConflicts(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req adminConflictRequest
	if !bindJSON(c, &req) {
		return
	}
	if req.ClassID == 0 {
		fail(c, services.ErrUnprocessable("class_id 必填"))
		return
	}
	exists, err := h.timetable.ClassInSchool(u.SchoolID, req.ClassID)
	if err != nil {
		fail(c, err)
		return
	}
	if !exists {
		fail(c, services.ErrNotFound("班级不存在"))
		return
	}
	for _, entry := range req.Entries {
		if entry.Weekday != 0 && (entry.Weekday < 1 || entry.Weekday > 7) {
			fail(c, services.ErrUnprocessable("entries.weekday 需在 1-7 之间"))
			return
		}
		if entry.PeriodIndex != 0 && (entry.PeriodIndex < 1 || entry.PeriodIndex > 30) {
			fail(c, services.ErrUnprocessable("entries.period_index 需在 1-30 之间"))
			return
		}
		if runeLen(entry.SubjectName) > 50 || runeLen(entry.TeacherName) > 50 {
			fail(c, services.ErrUnprocessable("entries.subject_name / teacher_name 不能超过 50 字"))
			return
		}
		switch entry.WeekType {
		case "", "all", "odd", "even":
		default:
			fail(c, services.ErrUnprocessable("entries.week_type 取值不合法"))
			return
		}
	}

	conflicts, err := h.timetable.CheckConflicts(u.SchoolID, req.ClassID, req.Entries)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, conflicts)
}

// AdminClassTeacherAssignmentsList 某班任课列表（subject_name => teacher_name）。
func (h *Handlers) AdminClassTeacherAssignmentsList(c *gin.Context) {
	u := middleware.CurrentUser(c)

	classID, valid := paramID(c, "id")
	if !valid {
		return
	}
	exists, err := h.timetable.ClassInSchool(u.SchoolID, classID)
	if err != nil {
		fail(c, err)
		return
	}
	if !exists {
		fail(c, services.ErrNotFound("班级不存在"))
		return
	}

	rows, err := h.timetable.ListAssignments(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// adminAssignmentsRequest 任课保存请求体。
type adminAssignmentsRequest struct {
	Assignments []services.TimetableAssignmentRow `json:"assignments"`
}

// AdminClassTeacherAssignmentsSave 整体保存某班任课（replace 语义）。
func (h *Handlers) AdminClassTeacherAssignmentsSave(c *gin.Context) {
	u := middleware.CurrentUser(c)

	classID, valid := paramID(c, "id")
	if !valid {
		return
	}
	exists, err := h.timetable.ClassInSchool(u.SchoolID, classID)
	if err != nil {
		fail(c, err)
		return
	}
	if !exists {
		fail(c, services.ErrNotFound("班级不存在"))
		return
	}

	var req adminAssignmentsRequest
	if !bindJSON(c, &req) {
		return
	}
	// 与 Laravel 的 `assignments => required|array` 对齐：字段缺失或空数组一律 422。
	if len(req.Assignments) == 0 {
		fail(c, services.ErrUnprocessable("assignments 必填"))
		return
	}
	for _, row := range req.Assignments {
		if trimSpace(row.SubjectName) == "" {
			fail(c, services.ErrUnprocessable("assignments.subject_name 必填"))
			return
		}
		if trimSpace(row.TeacherName) == "" {
			fail(c, services.ErrUnprocessable("assignments.teacher_name 必填"))
			return
		}
		if runeLen(row.SubjectName) > 50 || runeLen(row.TeacherName) > 50 {
			fail(c, services.ErrUnprocessable("assignments 名称不能超过 50 字"))
			return
		}
	}

	if err := h.timetable.SaveAssignments(classID, u.SchoolID, req.Assignments); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "任课已保存")
}

// trimSpace 去除首尾空格。
func trimSpace(value string) string {
	return strings.TrimSpace(value)
}

// runeLen 字符串的字符数（Laravel max 校验按字符计）。
func runeLen(value string) int {
	return utf8.RuneCountInString(value)
}
