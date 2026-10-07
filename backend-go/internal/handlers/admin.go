// 管理端处理器。
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// AdminSchool 返回当前管理员所属学校。
//
// `settings` 以**对象**输出（见 `services.SchoolView`）——Go 的模型把 settings 标了
// `json:"-"`，不经视图输出会让前端「学校设置」页读不到第三方平台开关等设置。

// classRoomPayload 班级对外载荷（settings 以对象输出，见 services.ClassRoomView）。
func classRoomPayload(class *models.ClassRoom) services.ClassRoomView {
	return services.ClassRoomViewOf(class)
}
func (h *Handlers) AdminSchool(c *gin.Context) {
	u := middleware.CurrentUser(c)
	school, err := h.admin.School(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, services.SchoolViewOf(school))
}

// AdminUpdateSchool 更新学校信息（Laravel 路由为 `match(['put','post'], 'school')`）。
//
// 与 Laravel `SchoolAdminController::updateSchool` 对齐：
//   - `settings` 为 `nullable|array`，**提供时整体替换** settings 列（不做按键合并）；
//   - 非法类型 → 422「参数错误」+ `errors.settings`；
//   - 成功文案「学校信息已更新」，`data` 为最新学校视图（Laravel 返回 `$school->fresh()`）。
func (h *Handlers) AdminUpdateSchool(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Name         *string         `json:"name"`
		Address      *string         `json:"address"`
		ContactPhone *string         `json:"contact_phone"`
		ContactEmail *string         `json:"contact_email"`
		LogoPath     *string         `json:"logo_path"`
		Settings     json.RawMessage `json:"settings"`
	}
	if !bindJSON(c, &req) {
		return
	}

	settings, errs := parseSchoolSettings(req.Settings)
	if len(errs) > 0 {
		validationFailed(c, errs)
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		updates["contact_email"] = *req.ContactEmail
	}
	if req.LogoPath != nil {
		updates["logo_path"] = *req.LogoPath
	}

	school, err := h.admin.UpdateSchool(u.SchoolID, updates, settings)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "学校信息已更新", "data": services.SchoolViewOf(school)})
}

// parseSchoolSettings 解析 `/admin/school` 的 `settings`（Laravel `nullable|array`）：
//
//	缺失      → (nil, nil)      保持原值
//	null      → (&空对象, nil)  清空（Laravel `fill(['settings' => null])`）
//	JSON 对象 → (&对象, nil)    **整体替换**
//	其他类型  → (nil, 422 errors)
//
// 用 `json.RawMessage` 而不是 `*map[string]any`，是为了区分「未提供」与「显式 null」。
func parseSchoolSettings(raw json.RawMessage) (*map[string]any, map[string][]string) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, nil
	}
	if string(trimmed) == "null" {
		empty := map[string]any{}
		return &empty, nil
	}
	var parsed map[string]any
	if err := json.Unmarshal(trimmed, &parsed); err != nil {
		return nil, map[string][]string{"settings": {"settings 必须是对象"}}
	}
	return &parsed, nil
}

// AdminListClasses 返回学校班级列表。
func (h *Handlers) AdminListClasses(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classes, err := h.admin.ListClasses(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	views := make([]services.ClassRoomView, 0, len(classes))
	for i := range classes {
		views = append(views, services.ClassRoomViewOf(&classes[i]))
	}
	ok(c, views)
}

// AdminCreateClass 创建班级。
func (h *Handlers) AdminCreateClass(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Grade       string `json:"grade"`
		Year        string `json:"year"`
		TeacherID   *uint  `json:"teacher_id"`
		MaxStudents int    `json:"max_students"`
	}
	if !bindJSON(c, &req) {
		return
	}

	class, err := h.admin.CreateClass(u.SchoolID, req.Name, req.Grade, req.Year, req.TeacherID, req.MaxStudents)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, classRoomPayload(class))
}

// AdminUpdateClass 更新班级。
func (h *Handlers) AdminUpdateClass(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Grade       *string `json:"grade"`
		Year        *string `json:"year"`
		TeacherID   *uint   `json:"teacher_id"`
		MaxStudents *int    `json:"max_students"`
		Status      *string `json:"status"`
		PetSeries   *string `json:"pet_series"`
	}
	if !bindJSON(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Grade != nil {
		updates["grade"] = *req.Grade
	}
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.TeacherID != nil {
		updates["teacher_id"] = *req.TeacherID
	}
	if req.MaxStudents != nil {
		updates["max_students"] = *req.MaxStudents
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// `pet_series` 为 `nullable|string|in:...`（逐字同 Laravel update）：合法则**合并**进 settings，不合法 422。
	settingsPatch := map[string]any{}
	if req.PetSeries != nil {
		if !services.ValidClassRoomSeries(*req.PetSeries) {
			validationFailed(c, map[string][]string{"pet_series": {"班级宠物系列取值不合法"}})
			return
		}
		settingsPatch["pet_series"] = *req.PetSeries
	}

	class, err := h.admin.UpdateClassWithSettings(u.SchoolID, classID, updates, settingsPatch)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, classRoomPayload(class))
}

// AdminDeleteClass 删除班级。
func (h *Handlers) AdminDeleteClass(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}

	if err := h.admin.DeleteClass(u.SchoolID, classID); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "班级已删除")
}

// AdminUpdateStudent 更新学生。
func (h *Handlers) AdminUpdateStudent(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		Name      *string `json:"name"`
		Gender    *string `json:"gender"`
		StudentNo *string `json:"student_no"`
		Status    *string `json:"status"`
	}
	if !bindJSON(c, &req) {
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
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	student, err := h.admin.UpdateStudent(u.SchoolID, studentID, updates)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, student)
}

// AdminDeleteStudent 删除学生。
func (h *Handlers) AdminDeleteStudent(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "id")
	if !valid {
		return
	}

	if err := h.admin.DeleteStudent(u.SchoolID, studentID); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "学生已删除")
}

// AdminListTeachers 返回学校教师列表。
func (h *Handlers) AdminListTeachers(c *gin.Context) {
	u := middleware.CurrentUser(c)
	teachers, err := h.admin.ListTeachers(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, teachers)
}

// AdminCreateTeacher 创建教师。
func (h *Handlers) AdminCreateTeacher(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone"`
	}
	if !bindJSON(c, &req) {
		return
	}

	teacher, err := h.admin.CreateTeacher(u.SchoolID, req.Username, req.Password, req.Name, req.Phone)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, teacher)
}

// AdminUpdateTeacher 更新教师。
func (h *Handlers) AdminUpdateTeacher(c *gin.Context) {
	u := middleware.CurrentUser(c)
	teacherID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Nickname *string `json:"nickname"`
		Subject  *string `json:"subject"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		Status   *string `json:"status"`
	}
	if !bindJSON(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	teacher, err := h.admin.UpdateTeacher(u.SchoolID, teacherID, updates)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, teacher)
}

// AdminResetPassword 重置教师密码。
func (h *Handlers) AdminResetPassword(c *gin.Context) {
	u := middleware.CurrentUser(c)
	teacherID, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	if err := h.admin.ResetPassword(u.SchoolID, teacherID, req.Password); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "密码已重置")
}
