// 教师端报表处理器（积分趋势 / 宠物等级分布 / 学生进度）。
// 路由路径逐字对齐 Laravel backend/routes/api.php 第 238-242 行。
package handlers

import (
	"fmt"
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// TeacherReportScoreTrend 近 N 天得分 / 扣分日趋势（days 默认 7，收敛到 1..365）。
func (h *Handlers) TeacherReportScoreTrend(c *gin.Context) {
	u := middleware.CurrentUser(c)

	data, err := h.reports.ScoreTrend(u, queryInt(c, "days", 7))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherReportPetDistribution 宠物等级分布（按 level 升序）。
func (h *Handlers) TeacherReportPetDistribution(c *gin.Context) {
	u := middleware.CurrentUser(c)

	data, err := h.reports.PetDistribution(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherReportStudentProgress 学生进度：带 student_id 返回单人近 50 条历史，否则全班列表。
func (h *Handlers) TeacherReportStudentProgress(c *gin.Context) {
	u := middleware.CurrentUser(c)

	data, err := h.reports.StudentProgress(u, queryID(c, "student_id"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherReportExport 报表导出（CSV 版）：GET /api/v1/teacher/reports/export/:type。
//
// 忠实移植 Laravel TeacherController::exportReport：作用域校验（class_id 缺省取管辖第一个班级）
// 不过则 403「无权限」；未知 type 则 200「导出类型 X 不支持，可选: scores, pets, attendance」；
// 成功则附件下载（Content-Type text/csv; charset=UTF-8，
// Content-Disposition attachment; filename="<rawurlencode(文件名)>"，中文编码同 timetable/export-cses）。
//
// 有意差异：Laravel 走 maatwebsite/excel 输出 xlsx，Go 端改用标准库 encoding/csv（后缀 .csv、内容带
// UTF-8 BOM）；「不支持」分支沿用本仓库的 {data:null, message} 信封（Laravel 只回 message）。
func (h *Handlers) TeacherReportExport(c *gin.Context) {
	u := middleware.CurrentUser(c)
	exportType := c.Param("type")

	classIDRaw, classIDProvided := c.GetQuery("class_id")
	file, err := h.reports.ExportReport(u, exportType, classIDRaw, classIDProvided, c.Query("date"))
	if err != nil {
		fail(c, err)
		return
	}
	if file == nil {
		okMessage(c, fmt.Sprintf("导出类型 %s 不支持，可选: scores, pets, attendance", exportType))
		return
	}

	c.Header("Content-Disposition", `attachment; filename="`+rawURLEncode(file.Filename)+`"`)
	c.Data(http.StatusOK, "text/csv; charset=UTF-8", file.Content)
}
