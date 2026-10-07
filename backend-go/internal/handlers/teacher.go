// 教师端处理器。
package handlers

import (
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// TeacherDashboard 返回教师首页统计。
func (h *Handlers) TeacherDashboard(c *gin.Context) {
	u := middleware.CurrentUser(c)
	data, err := h.teacher.Dashboard(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherStudents 教师管辖班级内的学生分页列表（GET /teacher/students）。
//
// 忠实移植 Laravel TeacherController::listStudents + StudentService::list：
// **不过滤 status**（停用/毕业学生同样列出）、可选 `search`（姓名或学号 like）、
// 按 `name` 升序、每页固定 50（Laravel `paginate(50)`，`per_page` 参数不生效）；
// 每行在 Student 原字段之上追加 `pet_species`/`pet_level`/`pet_name`，并带
// `meta`（current_page/last_page/per_page/total）。
//
// 有意差异：额外带 `message:"ok"`（仓库统一信封，Laravel 只有 data/meta）。
func (h *Handlers) TeacherStudents(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	// `per_page` 刻意不读：Laravel `paginate(50)` 硬编码每页 50，
	// 前端即使传 per_page=100 也只返回 50 条（差异见函数注释）。
	page, err := h.students.List(classIDs, c.Query("search"), queryInt(c, "page", 1), 50)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": page.Data, "meta": page.Meta, "message": "ok"})
}

// TeacherScoreRules 返回教师可见积分规则（首次访问补齐默认规则）。
func (h *Handlers) TeacherScoreRules(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	rules, err := h.rules.ListForTeacher(u, classIDs)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rules)
}

// TeacherCreateScoreRule 创建积分规则。
func (h *Handlers) TeacherCreateScoreRule(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		ClassID  *uint  `json:"class_id"`
		Name     string `json:"name" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		Category string `json:"category"`
	}
	if !bindJSON(c, &req) {
		return
	}

	if req.ClassID != nil {
		if err := h.scope.ClassInScope(u, *req.ClassID); err != nil {
			fail(c, err)
			return
		}
	} else if len(classIDs) > 0 {
		// 未指定班级时落到教师首个可管理班级。
		req.ClassID = &classIDs[0]
	}

	rule, err := h.rules.Create(req.ClassID, u.SchoolID, req.Name, req.Amount, req.Category, req.Amount >= 0)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, rule)
}

// TeacherUpdateScoreRule 更新积分规则。
func (h *Handlers) TeacherUpdateScoreRule(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	ruleID, valid := paramID(c, "id")
	if !valid {
		return
	}

	rule, err := h.rules.FindScopedForTeacher(u, classIDs, ruleID)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Amount   *int    `json:"amount"`
		Category *string `json:"category"`
		IsActive *bool   `json:"is_active"`
	}
	if !bindJSON(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
		updates["is_positive"] = *req.Amount >= 0
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	updated, err := h.rules.Update(rule, updates)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, updated)
}

// TeacherDeleteScoreRule 删除积分规则。
func (h *Handlers) TeacherDeleteScoreRule(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	ruleID, valid := paramID(c, "id")
	if !valid {
		return
	}

	rule, err := h.rules.FindScopedForTeacher(u, classIDs, ruleID)
	if err != nil {
		fail(c, err)
		return
	}
	if err := h.rules.Delete(rule); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "规则已删除")
}

// TeacherScoreSummary 返回管辖班级积分汇总。
func (h *Handlers) TeacherScoreSummary(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	summary, err := h.scores.Summary(classIDs)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, summary)
}

// TeacherScoreRecent 返回最近积分记录。
func (h *Handlers) TeacherScoreRecent(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	items, err := h.scores.Recent(classIDs, queryInt(c, "limit", 20))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// TeacherScoreGive 给单个学生加减分。
func (h *Handlers) TeacherScoreGive(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		StudentID uint   `json:"student_id" binding:"required"`
		Points    int    `json:"points" binding:"required"`
		Reason    string `json:"reason"`
		RuleID    *uint  `json:"rule_id"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Points == 0 {
		fail(c, services.ErrUnprocessable("积分变动不能为 0"))
		return
	}

	student, err := h.scope.StudentInScope(u, req.StudentID)
	if err != nil {
		fail(c, err)
		return
	}

	score, err := h.scores.GiveScore(student, req.Points, req.Reason, u.ID, req.RuleID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"score": score, "balance": student.TotalScore})
}

// TeacherScoreBatchGive 批量加减分。
func (h *Handlers) TeacherScoreBatchGive(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		StudentIDs []uint `json:"student_ids" binding:"required"`
		Points     int    `json:"points" binding:"required"`
		Reason     string `json:"reason"`
		RuleID     *uint  `json:"rule_id"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Points == 0 {
		fail(c, services.ErrUnprocessable("积分变动不能为 0"))
		return
	}

	students := make([]*models.Student, 0, len(req.StudentIDs))
	for _, id := range req.StudentIDs {
		student, err := h.scope.StudentInScope(u, id)
		if err != nil {
			fail(c, err)
			return
		}
		students = append(students, student)
	}

	count, err := h.scores.BatchGive(students, req.Points, req.Reason, u.ID, req.RuleID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"count": count})
}

// TeacherScoreUndo 撤回一条积分记录。
func (h *Handlers) TeacherScoreUndo(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	scoreID, valid := paramID(c, "id")
	if !valid {
		return
	}

	score, err := h.scores.FindScoredForTeacher(classIDs, scoreID)
	if err != nil {
		fail(c, err)
		return
	}

	undo, err := h.scores.Undo(score, u.ID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, undo)
}

// TeacherScoreHistory 返回单个学生积分历史。
func (h *Handlers) TeacherScoreHistory(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}

	if _, err := h.scope.StudentInScope(u, studentID); err != nil {
		fail(c, err)
		return
	}

	history, err := h.scores.History(studentID, queryInt(c, "limit", 20))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, history)
}

// TeacherPetOverview 返回班级宠物总览。
func (h *Handlers) TeacherPetOverview(c *gin.Context) {
	u := middleware.CurrentUser(c)
	items, err := h.pets.ClassOverview(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// TeacherPetDetail 返回学生宠物详情。
func (h *Handlers) TeacherPetDetail(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}
	pet, err := h.pets.PetFor(u, studentID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, pet)
}

// TeacherPetFeed 喂养宠物。
func (h *Handlers) TeacherPetFeed(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}
	result, err := h.pets.Feed(u, studentID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// TeacherPetRename 宠物更名。
func (h *Handlers) TeacherPetRename(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	message, err := h.pets.Rename(u, studentID, req.Name)
	if err != nil {
		fail(c, err)
		return
	}
	okMessage(c, message)
}

// TeacherPetSwitch 切换/首次分配宠物。
//
// 入参与响应逐字同 Laravel `TeacherController::switchPet` / `PetService::switchPet`：
// `pet_species`（必填，≤50）、`pet_name`（可空，缺省回退 `pet_species`）；
// 响应为顶层 `{message, data}`（此前是 Go 端统一的 `{data:{...,message:"ok"}}` 信封，
// 且入参误用 `species`/`name` —— 前端 PetCollection.vue 发的是 `pet_species`，现按 Laravel 修正）。
func (h *Handlers) TeacherPetSwitch(c *gin.Context) {
	u := middleware.CurrentUser(c)
	studentID, valid := paramID(c, "student_id")
	if !valid {
		return
	}

	var req struct {
		PetSpecies string  `json:"pet_species" binding:"required,max=50"`
		PetName    *string `json:"pet_name"`
	}
	if !bindJSON(c, &req) {
		return
	}

	// Laravel：`(string) $request->input('pet_name', $request->input('pet_species'))`
	// —— pet_name 键缺失时用物种 id，显式传 null 时为空串。
	petName := req.PetSpecies
	if req.PetName != nil {
		petName = *req.PetName
	}

	result, err := h.pets.Switch(u, studentID, req.PetSpecies, petName)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result.Message, "data": result.Data})
}

// TeacherLeaderboardTotal 班级总分排行榜（同 Laravel `totalLeaderboard`）：
// 默认取教师可管理班级的第一个；没有任何班级 → `{data: []}`（200）。
func (h *Handlers) TeacherLeaderboardTotal(c *gin.Context) {
	classID, empty, classOK := h.leaderboardClass(c)
	if !classOK {
		return
	}
	if empty {
		ok(c, []services.LeaderboardEntry{})
		return
	}
	entries, err := h.leaderboard.Total(classID, queryInt(c, "limit", 20))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, entries)
}

// TeacherLeaderboardWeekly 班级本周积分排行榜（同 Laravel `weeklyLeaderboard`）。
func (h *Handlers) TeacherLeaderboardWeekly(c *gin.Context) {
	classID, empty, classOK := h.leaderboardClass(c)
	if !classOK {
		return
	}
	if empty {
		ok(c, []services.LeaderboardEntry{})
		return
	}
	entries, err := h.leaderboard.Weekly(classID, queryInt(c, "limit", 20))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, entries)
}

// TeacherLeaderboardPetLevel 宠物等级排行榜（同 Laravel `petLevelLeaderboard`）。
func (h *Handlers) TeacherLeaderboardPetLevel(c *gin.Context) {
	classID, empty, classOK := h.leaderboardClass(c)
	if !classOK {
		return
	}
	if empty {
		ok(c, []services.LeaderboardEntry{})
		return
	}
	entries, err := h.leaderboard.PetLevel(classID, queryInt(c, "limit", 20))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, entries)
}

// leaderboardClass 解析并校验排行榜所需的班级 ID。
// leaderboardClass 解析排行榜目标班级：**默认取教师可管理班级中的第一个**（同 Laravel
// `$classIds->first()`，前端 `LeaderboardPage.vue` 不传任何参数）；显式传 `class_id` 时按该班
// （越权/不存在 → 404，属 Go 端的可选扩展）。教师没有任何可管理班级时返回 `empty=true`，
// 由调用方按 Laravel 回 `{data: []}`（200），而不是报错。
func (h *Handlers) leaderboardClass(c *gin.Context) (classID uint, empty bool, ok bool) {
	u := middleware.CurrentUser(c)
	if requested := queryID(c, "class_id"); requested != 0 {
		if err := h.scope.ClassInScope(u, requested); err != nil {
			fail(c, err)
			return 0, false, false
		}
		return requested, false, true
	}

	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return 0, false, false
	}
	if len(classIDs) == 0 {
		return 0, true, true
	}
	return classIDs[0], false, true
}
