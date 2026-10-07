// 教师端考勤处理器。
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
)

// TeacherAttendanceToday 今日考勤记录列表。
func (h *Handlers) TeacherAttendanceToday(c *gin.Context) {
	u := middleware.CurrentUser(c)
	items, err := h.attendance.Today(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// TeacherAttendanceStart 开始点名（为管辖班级活跃学生创建今日考勤记录）。
func (h *Handlers) TeacherAttendanceStart(c *gin.Context) {
	u := middleware.CurrentUser(c)
	result, err := h.attendance.Start(u)
	if err != nil {
		fail(c, err)
		return
	}

	msg := "已为 " + strconv.Itoa(result.Total) + " 名学生创建考勤记录（默认到课）"
	ok(c, gin.H{"message": msg, "data": result})
}

// TeacherAttendanceSet 手动设置某学生今日考勤状态。
func (h *Handlers) TeacherAttendanceSet(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}

	var req struct {
		Status string  `json:"status" binding:"required"`
		Remark *string `json:"remark"`
	}
	if !bindJSON(c, &req) {
		return
	}
	// 校验 status 合法枚举（Laravel: in:present,late,leave,absent）。
	switch req.Status {
	case "present", "late", "leave", "absent":
	default:
		fail(c, services.ErrBadRequest("考勤状态不合法"))
		return
	}

	record, err := h.attendance.SetStatus(u, studentID, req.Status, req.Remark)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"message": "考勤状态已更新", "data": record})
}

// TeacherAttendanceMarkLeave 手动标记请假。
func (h *Handlers) TeacherAttendanceMarkLeave(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}

	var req struct {
		Remark string `json:"remark" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	record, err := h.attendance.MarkLeave(u, studentID, req.Remark)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{
		"message": "已标记为请假",
		"data": gin.H{
			"id":     record.ID,
			"status": record.Status,
			"source": record.Source,
			"remark": record.Remark,
		},
	})
}

// TeacherAttendanceMarkAbsent 手动标记缺勤。
func (h *Handlers) TeacherAttendanceMarkAbsent(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}

	var req struct {
		Remark *string `json:"remark"`
	}
	if !bindJSON(c, &req) {
		return
	}

	record, err := h.attendance.MarkAbsent(u, studentID, req.Remark)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{
		"message": "已标记为缺勤，建议联系家长确认情况",
		"data": gin.H{
			"id":     record.ID,
			"status": record.Status,
			"source": record.Source,
			"remark": record.Remark,
		},
	})
}

// TeacherAttendanceSummary 今日考勤四态统计。
func (h *Handlers) TeacherAttendanceSummary(c *gin.Context) {
	u := middleware.CurrentUser(c)
	summary, err := h.attendance.Summary(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, summary)
}
