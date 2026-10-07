// 管理端处理器（本批 5 条）：
//   - GET  /admin/students              学生分页列表
//   - POST /admin/students              单建学生
//   - GET  /admin/classes/:id           班级详情（含 teacher / students）
//   - PUT  /admin/teachers/:id/classes   教师班级分配（replace 语义）
//   - GET  /admin/teachers/:id/password  查看教师密码（无明文则自动生成并写回）
//
// 逐条对齐 Laravel SchoolAdminController：listStudents / createStudent / show /
// assignTeacherClasses / getTeacherPassword。校验失败统一 422「参数错误」+ errors。
package handlers

import (
	"net/http"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// AdminStudentsList 学生分页列表（GET /admin/students）。
func (h *Handlers) AdminStudentsList(c *gin.Context) {
	u := middleware.CurrentUser(c)
	page, err := h.admin.ListStudentsPaged(u.SchoolID, services.AdminStudentQuery{
		Search:  c.Query("search"),
		ClassID: queryID(c, "class_id"),
		Grade:   c.Query("grade"),
		Status:  c.Query("status"),
		PerPage: queryInt(c, "per_page", 50),
		Page:    queryInt(c, "page", 1),
	})
	if err != nil {
		fail(c, err)
		return
	}
	// Laravel 直接返回 {data, meta}（不带 message）。
	c.JSON(http.StatusOK, gin.H{"data": page.Data, "meta": page.Meta})
}

// AdminStudentCreate 单建学生（POST /admin/students）。
func (h *Handlers) AdminStudentCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)

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
	if req.Gender != nil && !validStudentGender(*req.Gender) {
		errs["gender"] = []string{"性别不合法"}
	}
	if req.StudentNo != nil && len([]rune(*req.StudentNo)) > 50 {
		errs["student_no"] = []string{"学号不能超过 50 字"}
	}
	if len(errs) > 0 {
		validationFailed(c, errs)
		return
	}

	classID := *req.ClassID

	gender, studentNo := "", ""
	if req.Gender != nil {
		gender = *req.Gender
	}
	if req.StudentNo != nil {
		studentNo = *req.StudentNo
	}

	// 班级整体不存在 → 422（服务层抛 ValidationError）；存在但非本校 → 404「班级不存在」。
	student, err := h.admin.CreateStudentInClass(u.SchoolID, classID, name, gender, studentNo)
	if err != nil {
		fail(c, err)
		return
	}
	// 同 Laravel：201 + 「学生「X」已添加，已自动分配萌宠」。
	c.JSON(http.StatusCreated, gin.H{
		"message": "学生「" + student.Name + "」已添加，已自动分配萌宠",
		"data":    student,
	})
}

// validStudentGender 同 Laravel `in:男,女,男生,女生,未知`。
func validStudentGender(g string) bool {
	switch g {
	case "男", "女", "男生", "女生", "未知":
		return true
	}
	return false
}

// AdminClassShow 班级详情（GET /admin/classes/:id）：含 teacher 与 students 关联。
func (h *Handlers) AdminClassShow(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}
	detail, err := h.admin.ClassDetail(u.SchoolID, classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, detail)
}

// AdminTeacherAssignClasses 教师班级分配（PUT /admin/teachers/:id/classes，replace 语义）。
func (h *Handlers) AdminTeacherAssignClasses(c *gin.Context) {
	u := middleware.CurrentUser(c)
	teacherID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		Assignments []struct {
			ClassID *uint   `json:"class_id"`
			Role    *string `json:"role"`
			Subject *string `json:"subject"`
		} `json:"assignments"`
	}
	if !readJSONBody(c, &req) {
		return
	}

	errs := map[string][]string{}
	assignments := make([]services.TeacherClassAssignment, 0, len(req.Assignments))
	for i, item := range req.Assignments {
		prefix := "assignments." + itoa(i) + "."
		if item.ClassID == nil || *item.ClassID == 0 {
			errs[prefix+"class_id"] = []string{"班级必填"}
		} else if ok, err := h.admin.ClassExistsAnywhere(*item.ClassID); err != nil {
			fail(c, err)
			return
		} else if !ok {
			// 同 Laravel `exists:class_rooms,id`：班级整体不存在 → 422（存在但非本校由服务层回 404）。
			errs[prefix+"class_id"] = []string{"所选班级不存在"}
		}
		if item.Role == nil || *item.Role == "" {
			errs[prefix+"role"] = []string{"角色必填"}
		} else if !validAssignmentRole(*item.Role) {
			errs[prefix+"role"] = []string{"角色不合法"}
		}
		if item.Subject != nil && len([]rune(*item.Subject)) > 50 {
			errs[prefix+"subject"] = []string{"科目不能超过 50 字"}
		}
		entry := services.TeacherClassAssignment{Subject: item.Subject}
		if item.ClassID != nil {
			entry.ClassID = *item.ClassID
		}
		if item.Role != nil {
			entry.Role = *item.Role
		}
		assignments = append(assignments, entry)
	}
	if len(errs) > 0 {
		validationFailed(c, errs)
		return
	}

	synced, err := h.admin.AssignTeacherClasses(u.SchoolID, teacherID, assignments)
	if err != nil {
		fail(c, err)
		return
	}

	// 教师名从教师本体读取（同 Laravel 的响应文案）。
	teacher, err := h.admin.ListTeachers(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	teacherName := ""
	for _, t := range teacher {
		if t.ID == teacherID {
			teacherName = t.Name
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "已为教师「" + teacherName + "」分配 " + itoa(len(synced)) + " 个班级",
		"data": gin.H{
			"teacher_id":   teacherID,
			"teacher_name": teacherName,
			"assignments":  synced,
		},
	})
}

// validAssignmentRole 同 Laravel `in:head_teacher,co_teacher,subject_teacher,grade_lead,admin_director`。
func validAssignmentRole(role string) bool {
	switch role {
	case "head_teacher", "co_teacher", "subject_teacher", "grade_lead", "admin_director":
		return true
	}
	return false
}

// AdminTeacherPassword 查看教师密码（GET /admin/teachers/:id/password）。
// 无明文记录时自动生成 8 位随机密码并写回（同 Laravel，会重置该教师密码）。
func (h *Handlers) AdminTeacherPassword(c *gin.Context) {
	u := middleware.CurrentUser(c)
	teacherID, valid := paramID(c, "id")
	if !valid {
		return
	}
	password, _, err := h.admin.TeacherPassword(u.SchoolID, teacherID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"password": password})
}
