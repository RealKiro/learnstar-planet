// 管理端运维处理器：报表 / 大屏登录日志 / 系统运维 / 批量账号与导入 / 学年升级 / 学校 LOGO。
//
// 对应 Laravel App\Http\Controllers\Api\SchoolAdminController 的同名动作，
// 响应结构（data / meta / message）逐项对齐。
package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// maxTeacherImportBytes 教师导入文件上限（同 Laravel max:10240，单位 KB）。
const maxTeacherImportBytes = 10240 * 1024

// ============================================================
// 报表
// ============================================================

// AdminReportsOverview 全校概览。
func (h *Handlers) AdminReportsOverview(c *gin.Context) {
	u := middleware.CurrentUser(c)
	data, err := h.ops.SchoolOverview(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// AdminReportsByGrade 按年级汇总。
func (h *Handlers) AdminReportsByGrade(c *gin.Context) {
	u := middleware.CurrentUser(c)
	rows, err := h.ops.ReportsByGrade(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// AdminReportsByClass 按班级汇总。
func (h *Handlers) AdminReportsByClass(c *gin.Context) {
	u := middleware.CurrentUser(c)
	rows, err := h.ops.ReportsByClass(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// AdminDisplayLoginLogs 班级码大屏登录日志（分页）。
func (h *Handlers) AdminDisplayLoginLogs(c *gin.Context) {
	u := middleware.CurrentUser(c)
	res, err := h.ops.DisplayLoginLogs(u.SchoolID, services.DisplayLoginLogFilter{
		ClassID: queryID(c, "class_id"),
		IP:      c.Query("ip"),
		Date:    c.Query("date"),
		Page:    queryInt(c, "page", 1),
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res.Data, "meta": res.Meta})
}

// ============================================================
// 系统运维
// ============================================================

// AdminSystemDiagnose 系统诊断（表结构自检）。
func (h *Handlers) AdminSystemDiagnose(c *gin.Context) {
	res, err := h.ops.Diagnose()
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       res.Items,
		"has_issues": res.HasIssues,
		"message":    res.Message,
	})
}

// AdminSystemStatus 版本信息 + 库表清单。
func (h *Handlers) AdminSystemStatus(c *gin.Context) {
	data, err := h.ops.Status()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// AdminSystemLogs 读取日志（来源：环境变量 LOG_FILE，缺省返回空列表 + 说明）。
func (h *Handlers) AdminSystemLogs(c *gin.Context) {
	data := h.ops.SystemLogs(os.Getenv("LOG_FILE"), queryInt(c, "lines", 200), c.Query("level"))
	if data.Message != "" {
		c.JSON(http.StatusOK, gin.H{"data": data, "message": data.Message})
		return
	}
	ok(c, data)
}

// AdminSystemRepair 系统修复（AutoMigrate，幂等）。
func (h *Handlers) AdminSystemRepair(c *gin.Context) {
	res, err := h.ops.Repair()
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":                 res.Message,
		"output":                  res.Output,
		"pending_password_resets": res.PendingPasswordResets,
	})
}

// ============================================================
// 批量账号
// ============================================================

// AdminBatchResetPassword 批量重置教师密码。
func (h *Handlers) AdminBatchResetPassword(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Role     string `json:"role"`
		IDs      []uint `json:"ids"`
		Password string `json:"password"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Role != "teacher" {
		fail(c, services.ErrUnprocessable("参数错误"))
		return
	}

	count, err := h.ops.BatchResetPassword(u.SchoolID, req.IDs, req.Password)
	if err != nil {
		fail(c, err)
		return
	}
	password := req.Password
	if password == "" {
		password = services.DefaultTeacherPassword
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("已重置 %d 个账号的密码为「%s」", count, password),
		"data":    gin.H{"reset_count": count},
	})
}

// AdminDeleteTeacher 删除教师账号：DELETE /api/v1/admin/teachers/:id。
//
// 忠实移植 Laravel `SchoolAdminController::disableTeacher`：教师须在本校且 role=teacher（否则 404
// 「教师不存在」）；API 机器人账号 403「API 机器人账号不可删除。如需停用，请在 .env 中设置
// BOT_ENABLED=false 后重启」；成功解除班级关联并硬删除 → 「教师账号已删除」。
func (h *Handlers) AdminDeleteTeacher(c *gin.Context) {
	u := middleware.CurrentUser(c)
	teacherID, valid := paramID(c, "id")
	if !valid {
		return
	}

	if err := h.ops.DeleteTeacher(u.SchoolID, teacherID); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "教师账号已删除")
}

// AdminBatchDeleteAccounts 批量删除教师账号（API 机器人账号受保护）。
func (h *Handlers) AdminBatchDeleteAccounts(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Role string `json:"role"`
		IDs  []uint `json:"ids"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Role != "teacher" {
		fail(c, services.ErrUnprocessable("参数错误"))
		return
	}

	deleted, protected, err := h.ops.BatchDeleteAccounts(u.SchoolID, req.IDs)
	if err != nil {
		fail(c, err)
		return
	}
	message := fmt.Sprintf("已删除 %d 个账号", deleted)
	if protected > 0 {
		message += fmt.Sprintf("；%d 个 API 机器人账号不可删除，已跳过", protected)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    gin.H{"deleted_count": deleted, "protected_count": protected},
	})
}

// ============================================================
// 教师批量创建 / CSV 导入 / 模板下载
// ============================================================

// AdminBatchCreateTeachers 批量创建教师账号（可同时分配班级）。
func (h *Handlers) AdminBatchCreateTeachers(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Teachers []struct {
			Name        string `json:"name"`
			Nickname    string `json:"nickname"`
			Subject     string `json:"subject"`
			GradeTeam   string `json:"grade_team"`
			Phone       string `json:"phone"`
			Email       string `json:"email"`
			Username    string `json:"username"`
			Password    string `json:"password"`
			Assignments []struct {
				ClassID uint   `json:"class_id"`
				Role    string `json:"role"`
				Subject string `json:"subject"`
			} `json:"assignments"`
		} `json:"teachers"`
	}
	if !bindJSON(c, &req) {
		return
	}

	inputs := make([]services.TeacherInput, 0, len(req.Teachers))
	for _, t := range req.Teachers {
		in := services.TeacherInput{
			Name:      t.Name,
			Nickname:  t.Nickname,
			Subject:   t.Subject,
			GradeTeam: t.GradeTeam,
			Phone:     t.Phone,
			Email:     t.Email,
			Username:  t.Username,
			Password:  t.Password,
		}
		for _, a := range t.Assignments {
			in.Assignments = append(in.Assignments, services.TeacherAssignmentInput{
				ClassID: a.ClassID,
				Role:    a.Role,
				Subject: a.Subject,
			})
		}
		inputs = append(inputs, in)
	}

	created, err := h.ops.BatchCreateTeachers(u.SchoolID, inputs)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, created)
}

// AdminImportTeachers CSV 批量导入教师（multipart：file + dry_run，dry_run 默认 true=预览）。
func (h *Handlers) AdminImportTeachers(c *gin.Context) {
	u := middleware.CurrentUser(c)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fail(c, services.ErrUnprocessable("文件上传失败"))
		return
	}
	if fileHeader.Size > maxTeacherImportBytes {
		fail(c, services.ErrUnprocessable("文件上传失败：文件不能超过 10MB"))
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		fail(c, services.ErrUnprocessable("文件上传失败"))
		return
	}
	defer func() { _ = f.Close() }()

	content, err := io.ReadAll(io.LimitReader(f, maxTeacherImportBytes+1))
	if err != nil {
		fail(c, services.ErrUnprocessable("文件上传失败"))
		return
	}

	// Laravel `$request->boolean('dry_run', true)`：键缺失时为 true；空串解析为 false。
	dryRun := true
	if v, exists := c.GetPostForm("dry_run"); exists {
		dryRun = parseBoolLoose(v)
	}

	result, err := h.ops.ImportTeachersCSV(u.SchoolID, fileHeader.Filename, content, dryRun)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// AdminTeacherTemplateCsv 下载教师导入 CSV 模板（带 BOM，响应头同 Laravel）。
func (h *Handlers) AdminTeacherTemplateCsv(c *gin.Context) {
	c.Header("Content-Disposition", `attachment; filename="teacher_import_template.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=UTF-8", []byte(h.ops.TeacherTemplateCSV()))
}

// ============================================================
// 学生批量导入 / 删除 / 转班
// ============================================================

// AdminImportStudents 学生批量导入。
//
// 与 Laravel 一致：JSON 体 `{"students":[{name,class_name,gender,student_no}]}`；
// 额外支持 multipart 文件上传（字段名 `file`，CSV）。
func (h *Handlers) AdminImportStudents(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var rows []services.StudentImportRow

	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			fail(c, services.ErrUnprocessable("参数错误：请上传文件（字段名 file）"))
			return
		}
		f, err := fileHeader.Open()
		if err != nil {
			fail(c, services.ErrUnprocessable("参数错误：文件读取失败"))
			return
		}
		defer func() { _ = f.Close() }()

		parsed, perr := h.ops.ParseStudentCSV(f)
		if perr != nil {
			fail(c, perr)
			return
		}
		rows = parsed
	} else {
		var req struct {
			Students []struct {
				Name      string `json:"name"`
				ClassName string `json:"class_name"`
				Gender    string `json:"gender"`
				StudentNo string `json:"student_no"`
			} `json:"students"`
		}
		if !bindJSON(c, &req) {
			return
		}
		if len(req.Students) == 0 {
			fail(c, services.ErrUnprocessable("参数错误：students 字段必填"))
			return
		}
		for i, s := range req.Students {
			if strings.TrimSpace(s.Name) == "" {
				fail(c, services.ErrUnprocessable(fmt.Sprintf("参数错误：第 %d 行的姓名必填", i+1)))
				return
			}
			if strings.TrimSpace(s.ClassName) == "" {
				fail(c, services.ErrUnprocessable(fmt.Sprintf("参数错误：第 %d 行的班级必填", i+1)))
				return
			}
			if !validImportGender(s.Gender) {
				fail(c, services.ErrUnprocessable(fmt.Sprintf("参数错误：第 %d 行的性别不合法", i+1)))
				return
			}
			rows = append(rows, services.StudentImportRow{
				Name:      s.Name,
				ClassName: s.ClassName,
				Gender:    s.Gender,
				StudentNo: s.StudentNo,
			})
		}
	}

	result, err := h.ops.ImportStudents(u.SchoolID, rows)
	if err != nil {
		fail(c, err)
		return
	}

	if len(result.Errors) > 0 {
		// Laravel：部分失败返回 422，但已成功导入的行仍然回传（且已落库）。
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message":  "部分导入失败",
			"errors":   result.Errors,
			"warnings": result.Warnings,
			"data": gin.H{
				"created":       result.Created,
				"created_count": len(result.Created),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "导入完成",
		"data": gin.H{
			"created_count": len(result.Created),
			"created":       result.Created,
			"warnings":      result.Warnings,
		},
	})
}

// AdminBatchDeleteStudents 批量删除学生。
func (h *Handlers) AdminBatchDeleteStudents(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		StudentIDs []uint `json:"student_ids"`
	}
	if !bindJSON(c, &req) {
		return
	}

	count, err := h.ops.BatchDeleteStudents(u.SchoolID, req.StudentIDs)
	if err != nil {
		fail(c, err)
		return
	}
	okMessage(c, fmt.Sprintf("已删除 %d 名学生", count))
}

// AdminBatchMoveStudents 批量转班。
func (h *Handlers) AdminBatchMoveStudents(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		StudentIDs    []uint `json:"student_ids"`
		TargetClassID uint   `json:"target_class_id"`
	}
	if !bindJSON(c, &req) {
		return
	}

	count, targetName, err := h.ops.BatchMoveStudents(u.SchoolID, req.StudentIDs, req.TargetClassID)
	if err != nil {
		fail(c, err)
		return
	}
	okMessage(c, fmt.Sprintf("已将 %d 名学生移动到「%s」", count, targetName))
}

// validImportGender 校验性别取值（同 Laravel in:男,女,男生,女生,未知）。
func validImportGender(g string) bool {
	switch strings.TrimSpace(g) {
	case "男", "女", "男生", "女生", "未知":
		return true
	}
	return false
}

// ============================================================
// 班级批量创建 + 教师分配
// ============================================================

// AdminBatchCreateClasses 批量创建班级（grade + count(1-20) + year）。
func (h *Handlers) AdminBatchCreateClasses(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Grade string `json:"grade"`
		Count int    `json:"count"`
		Year  string `json:"year"`
	}
	if !bindJSON(c, &req) {
		return
	}

	created, err := h.ops.BatchCreateClasses(u.SchoolID, req.Grade, req.Count, req.Year)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("已批量创建 %d 个班级", req.Count),
		"data":    created,
	})
}

// AdminAssignClassTeacher 分配班级教师（role=head_teacher 时同步 class_rooms.teacher_id）。
func (h *Handlers) AdminAssignClassTeacher(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		TeacherID uint   `json:"teacher_id"`
		Role      string `json:"role"`
	}
	if !bindJSON(c, &req) {
		return
	}

	class, err := h.ops.AssignClassTeacher(u.SchoolID, classID, req.TeacherID, req.Role)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "教师已分配", "data": class})
}

// AdminRemoveClassTeacher 移除班级教师。
func (h *Handlers) AdminRemoveClassTeacher(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		TeacherID uint `json:"teacher_id"`
	}
	// DELETE 允许无 body：解析失败时回退 query 参数（Laravel 同时接受 body/query）。
	if err := c.ShouldBindJSON(&req); err != nil {
		req.TeacherID = queryID(c, "teacher_id")
	}
	if req.TeacherID == 0 {
		req.TeacherID = queryID(c, "teacher_id")
	}

	if err := h.ops.RemoveClassTeacher(u.SchoolID, classID, req.TeacherID); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "教师已从班级移除")
}

// ============================================================
// 学年升级（不可逆）
// ============================================================

// AdminGradeUpgradePreview 学年升级预览（dry-run）。
func (h *Handlers) AdminGradeUpgradePreview(c *gin.Context) {
	u := middleware.CurrentUser(c)
	preview, err := h.ops.PreviewGradeUpgrade(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, preview)
}

// AdminGradeUpgradeExecute 执行学年升级（不可逆）。
func (h *Handlers) AdminGradeUpgradeExecute(c *gin.Context) {
	u := middleware.CurrentUser(c)
	result, err := h.ops.ExecuteGradeUpgrade(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "学年升级完成", "data": result})
}

// ============================================================
// 学校 LOGO 上传
// ============================================================

// AdminUploadSchoolLogo 上传学校 LOGO（multipart 字段名 logo）。
// 有意差异：Go 端不缩放/裁剪，直接保存原图。
func (h *Handlers) AdminUploadSchoolLogo(c *gin.Context) {
	u := middleware.CurrentUser(c)

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		fail(c, services.ErrUnprocessable("参数错误：logo 字段必填"))
		return
	}
	if fileHeader.Size > 2048*1024 {
		fail(c, services.ErrUnprocessable("参数错误：LOGO 不能超过 2MB"))
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		fail(c, services.ErrUnprocessable("参数错误：LOGO 读取失败"))
		return
	}
	defer func() { _ = f.Close() }()

	content, err := io.ReadAll(io.LimitReader(f, 2048*1024+1))
	if err != nil {
		fail(c, services.ErrUnprocessable("参数错误：LOGO 读取失败"))
		return
	}

	logoPath, err := h.ops.SaveSchoolLogo(u.SchoolID, fileHeader.Filename, content)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "LOGO 已上传", "data": gin.H{"logo_path": logoPath}})
}

// parseBoolLoose 解析 PHP FILTER_VALIDATE_BOOLEAN 口径的布尔值
// （"1"/"true"/"on"/"yes" → true，其余含空串 → false）。
func parseBoolLoose(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}
