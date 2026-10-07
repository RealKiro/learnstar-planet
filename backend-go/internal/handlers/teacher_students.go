// 教师端处理器（本批 7 条）：
//   - GET    /teacher/class                              班级信息卡（宠物系列页顶）
//   - POST   /teacher/students                           添加学生
//   - POST   /teacher/students/import                    批量导入学生
//   - PUT    /teacher/students/:id                       更新学生
//   - DELETE /teacher/students/:id                       删除学生
//   - GET    /teacher/pets/:student_id/collection        学生宠物图鉴
//   - POST   /teacher/scores/give-by-rule/:ruleId        按规则加减分
//
// 逐条对齐 Laravel TeacherController / StudentService / PetService / ScoreService
// / PetSeriesService；校验失败统一 422「参数错误」+ errors。
package handlers

import (
	"net/http"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// TeacherClassInfo 当前班级信息卡（GET /teacher/class）。
// 没有可管理的班级 → 400「没有可管理的班级」（Laravel DomainException 映射 400）。
func (h *Handlers) TeacherClassInfo(c *gin.Context) {
	u := middleware.CurrentUser(c)
	view, err := h.petSeries.ClassInfo(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// TeacherStudentCreate 添加学生（POST /teacher/students）。
// 目标班级不在自己管辖范围 → 403「只能在自己管理的班级添加学生」；成功 201。
func (h *Handlers) TeacherStudentCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Name      string  `json:"name"`
		ClassID   *uint   `json:"class_id"`
		Gender    *string `json:"gender"`
		StudentNo *string `json:"student_no"`
	}
	if !readJSONBody(c, &req) {
		return
	}

	errs := map[string][]string{}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		errs["name"] = []string{"姓名必填"}
	} else if len([]rune(name)) > 50 {
		errs["name"] = []string{"姓名不能超过 50 字"}
	}
	if req.ClassID == nil || *req.ClassID == 0 {
		errs["class_id"] = []string{"班级必填"}
	}
	if req.StudentNo != nil && len([]rune(*req.StudentNo)) > 50 {
		errs["student_no"] = []string{"学号不能超过 50 字"}
	}
	if len(errs) > 0 {
		validationFailed(c, errs)
		return
	}

	gender, studentNo := "", ""
	if req.Gender != nil {
		gender = *req.Gender
	}
	if req.StudentNo != nil {
		studentNo = *req.StudentNo
	}

	student, err := h.students.Create(classIDs, name, *req.ClassID, gender, studentNo)
	if err != nil {
		fail(c, err)
		return
	}
	if student == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "只能在自己管理的班级添加学生"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "学生「" + student.Name + "」已添加",
		"data":    student,
	})
}

// TeacherStudentImport 批量导入学生（POST /teacher/students/import）。
func (h *Handlers) TeacherStudentImport(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Students []struct {
			Name      string `json:"name"`
			ClassName string `json:"class_name"`
			Gender    string `json:"gender"`
			StudentNo string `json:"student_no"`
		} `json:"students"`
	}
	if !readJSONBody(c, &req) {
		return
	}

	rows := make([]services.TeacherStudentImportRow, 0, len(req.Students))
	for _, s := range req.Students {
		rows = append(rows, services.TeacherStudentImportRow{
			Name:      s.Name,
			ClassName: s.ClassName,
			Gender:    s.Gender,
			StudentNo: s.StudentNo,
		})
	}

	result, err := h.students.Import(classIDs, rows)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result.Message,
		"data": gin.H{
			"imported_count": result.ImportedCount,
			"skipped":        result.Skipped,
		},
	})
}

// TeacherStudentUpdate 更新学生（PUT /teacher/students/:id，仅 name/gender/student_no）。
func (h *Handlers) TeacherStudentUpdate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "id")
	if !valid {
		return
	}
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Name      *string `json:"name"`
		Gender    *string `json:"gender"`
		StudentNo *string `json:"student_no"`
	}
	if !readJSONBody(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.StudentNo != nil {
		updates["student_no"] = *req.StudentNo
	}

	student, err := h.students.Update(classIDs, studentID, updates)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": student})
}

// TeacherStudentDelete 删除学生（DELETE /teacher/students/:id）。
// 与 Laravel 的差异：Laravel 是软删除（SoftDeletes）；Go 端与既有 AdminDeleteStudent
// 一致——硬删除并级联清理宠物 / 积分 / 审计日志（见 services/student.go 注释）。
func (h *Handlers) TeacherStudentDelete(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "id")
	if !valid {
		return
	}
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	student, err := h.students.Delete(classIDs, studentID)
	if err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "学生「"+student.Name+"」已删除")
}

// TeacherPetCollection 学生宠物图鉴（GET /teacher/pets/:student_id/collection）。
func (h *Handlers) TeacherPetCollection(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}
	view, err := h.pets.Collection(u, studentID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// TeacherScoreGiveByRule 按规则给学生加减分（POST /teacher/scores/give-by-rule/:ruleId）。
// 顺序同 Laravel：先取可见规则（越权/不存在 → 404），再校验 student_id，最后取管辖内的学生。
func (h *Handlers) TeacherScoreGiveByRule(c *gin.Context) {
	u := middleware.CurrentUser(c)
	ruleID, valid := paramID(c, "ruleId")
	if !valid {
		return
	}

	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	rule, err := h.rules.FindScopedForTeacher(u, classIDs, ruleID)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		StudentID *uint `json:"student_id"`
	}
	if !readJSONBody(c, &req) {
		return
	}
	if req.StudentID == nil || *req.StudentID == 0 {
		validationFailed(c, map[string][]string{"student_id": {"学生必填"}})
		return
	}

	student, err := h.scope.StudentInScope(u, *req.StudentID)
	if err != nil {
		fail(c, err)
		return
	}

	if _, err := h.scores.GiveScoreByRule(student, rule, u.ID); err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "已按规则「" + rule.Name + "」处理",
		"data": gin.H{
			"rule":         rule.Name,
			"points":       rule.Amount,
			"student_name": student.Name,
		},
	})
}
