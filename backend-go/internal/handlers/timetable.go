// 课表处理器：教师端（查看/提交申请/我的课表/CSES 导出）+ 管理员端（直改/审批）。
package handlers

import (
	"net/http"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// timetableRequest 课表保存 / 提交申请的请求体。
type timetableRequest struct {
	ClassID  uint                             `json:"class_id"`
	Subjects []services.TimetableSubjectInput `json:"subjects"`
	Periods  []services.TimetablePeriodInput  `json:"periods"`
	Entries  []services.TimetableEntryInput   `json:"entries"`
}

// resolveTimetableClass 解析目标班级：优先显式 class_id，否则取管辖的第一个班级。
func (h *Handlers) resolveTimetableClass(c *gin.Context, explicit uint) (uint, bool) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return 0, false
	}
	if len(classIDs) == 0 {
		fail(c, services.ErrNotFound("您当前没有管辖的班级"))
		return 0, false
	}
	if explicit > 0 {
		for _, id := range classIDs {
			if id == explicit {
				return explicit, true
			}
		}
		fail(c, services.ErrNotFound("班级不存在或不在管辖范围"))
		return 0, false
	}
	return classIDs[0], true
}

// TeacherTimetableShow 教师端查看本班课表（科目 + 节次 + 排课）。
func (h *Handlers) TeacherTimetableShow(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, bok := h.resolveTimetableClass(c, queryID(c, "class_id"))
	if !bok {
		return
	}
	data, err := h.timetable.Bootstrap(classID, u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherTimetableSave 教师端提交课表修改申请（不直接生效，待管理员审核）。
func (h *Handlers) TeacherTimetableSave(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req timetableRequest
	if !bindJSON(c, &req) {
		return
	}
	classID, bok := h.resolveTimetableClass(c, req.ClassID)
	if !bok {
		return
	}
	change, err := h.timetable.SubmitChange(classID, u.SchoolID, u.ID, req.payload())
	if err != nil {
		fail(c, err)
		return
	}
	// Laravel `TimetableController::save` 逐字文案：'修改申请已提交，待管理员审核'（201）
	c.JSON(http.StatusCreated, gin.H{"data": change, "message": "修改申请已提交，待管理员审核"})
}

// TeacherTimetableChanges 本班申请历史（最新在前）。
func (h *Handlers) TeacherTimetableChanges(c *gin.Context) {
	classID, bok := h.resolveTimetableClass(c, queryID(c, "class_id"))
	if !bok {
		return
	}
	views, err := h.timetable.ListChangesForClass(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, views)
}

// TeacherTimetableMySchedule 某教师的自己的周课表（默认本人，可指定教师姓名）。
func (h *Handlers) TeacherTimetableMySchedule(c *gin.Context) {
	u := middleware.CurrentUser(c)
	name := strings.TrimSpace(c.Query("teacher_name"))
	if name == "" {
		name = u.Name
	}
	data, err := h.timetable.ForTeacher(u.SchoolID, name)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherTimetableExportCses 导出本班课表为 CSES YAML（ClassIsland 可直接导入）。
func (h *Handlers) TeacherTimetableExportCses(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, bok := h.resolveTimetableClass(c, queryID(c, "class_id"))
	if !bok {
		return
	}
	yaml, err := h.timetable.ToCses(classID, u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=timetable-cses.yaml")
	c.String(http.StatusOK, yaml)
}

// AdminClassTimetableShow 管理员查看某班课表。
func (h *Handlers) AdminClassTimetableShow(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, bok := paramID(c, "id")
	if !bok {
		return
	}
	data, err := h.timetable.Bootstrap(id, u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// AdminClassTimetableSave 管理员直接保存某班课表（即时生效，待审申请自动作废）。
func (h *Handlers) AdminClassTimetableSave(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, bok := paramID(c, "id")
	if !bok {
		return
	}
	var req timetableRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.timetable.AdminSave(id, u.SchoolID, req.payload()); err != nil {
		fail(c, err)
		return
	}
	// Laravel `adminSave` 逐字文案：'课表已保存并即时生效'
	okMessage(c, "课表已保存并即时生效")
}

// AdminTimetableChanges 管理员查看全校申请列表（可按状态过滤）。
func (h *Handlers) AdminTimetableChanges(c *gin.Context) {
	u := middleware.CurrentUser(c)
	views, err := h.timetable.ListChangesForSchool(u.SchoolID, c.Query("status"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, views)
}

// AdminTimetableChangeReview 审核请求体。
type timetableReviewRequest struct {
	Note *string `json:"note"`
}

// AdminTimetableApprove 审核通过并应用申请的课表快照。
func (h *Handlers) AdminTimetableApprove(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, bok := paramID(c, "id")
	if !bok {
		return
	}
	var req timetableReviewRequest
	_ = c.ShouldBindJSON(&req)
	applied, err := h.timetable.ApproveChange(id, u.ID, req.Note)
	if err != nil {
		fail(c, err)
		return
	}
	if !applied {
		fail(c, services.ErrBadRequest("申请不存在或已被处理"))
		return
	}
	okMessage(c, "已通过并应用课表")
}

// AdminTimetableReject 驳回课表修改申请。
func (h *Handlers) AdminTimetableReject(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, bok := paramID(c, "id")
	if !bok {
		return
	}
	var req timetableReviewRequest
	_ = c.ShouldBindJSON(&req)
	applied, err := h.timetable.RejectChange(id, u.ID, req.Note)
	if err != nil {
		fail(c, err)
		return
	}
	if !applied {
		fail(c, services.ErrBadRequest("申请不存在或已被处理"))
		return
	}
	okMessage(c, "已驳回")
}

// AdminClassTimetableExportExcel 导出某班课表网格：GET /api/v1/admin/classes/:id/timetable/export-excel。
//
// 忠实移植 Laravel TimetableController::adminExportClassExcel（含 resolveAdminClass 的归属校验：
// 班级不属于本校 → 404「班级不存在」）。网格行结构见 services.TimetableService.ExcelGrid。
//
// 有意差异：Laravel 输出 .xlsx（文件名 <班级名>-课表.xlsx），Go 端零新依赖改用 CSV
// （文件名 <班级名>-课表.csv，带 UTF-8 BOM），响应头 text/csv; charset=UTF-8 +
// rawurlencode 后的 Content-Disposition（中文编码同教师端 timetable/export-cses）。
func (h *Handlers) AdminClassTimetableExportExcel(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, bok := paramID(c, "id")
	if !bok {
		return
	}

	file, err := h.timetable.ExportClassFile(id, u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	writeCSVAttachment(c, file)
}

// AdminTimetableExportExcel 导出全校课表网格：GET /api/v1/admin/timetable/export-excel。
//
// 忠实移植 Laravel TimetableController::adminExportSchoolExcel：本校无班级 → 404「暂无班级可导出」；
// 否则按 grade ASC, name ASC 逐班生成网格（xlsx 为每班一个工作表；CSV 为顺序拼接，班级间空行分隔）。
func (h *Handlers) AdminTimetableExportExcel(c *gin.Context) {
	u := middleware.CurrentUser(c)

	file, err := h.timetable.ExportSchoolFile(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	if file == nil {
		fail(c, services.ErrNotFound("暂无班级可导出"))
		return
	}
	writeCSVAttachment(c, file)
}

// writeCSVAttachment 输出 CSV 附件（响应头与中文文件名编码方式，同 timetable/export-cses）。
func writeCSVAttachment(c *gin.Context, file *services.TimetableExportFile) {
	c.Header("Content-Disposition", `attachment; filename="`+rawURLEncode(file.Filename)+`"`)
	c.Data(http.StatusOK, "text/csv; charset=UTF-8", file.Content)
}

// payload 将请求体转换为服务层载荷。
func (r timetableRequest) payload() services.TimetablePayload {
	return services.TimetablePayload{
		Subjects: r.Subjects,
		Periods:  r.Periods,
		Entries:  r.Entries,
	}
}
