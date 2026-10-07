// 班级大屏（教室端）处理器：班级码管理 + 班级码登录 + 只读接口 + CSES 导出。
//
// 仅包含班级码管理 / 班级码登录 / 只读接口 / CSES 导出；大屏写操作与查询见 display_writes.go，
// 大屏端 AI（ai/settings 开关检查、ai/chat 对话）见 ai.go。
package handlers

import (
	"net/http"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// displayClassID 取大屏中间件注入的班级 ID（理论上不会为 0；防御性兜底 401）。

// displayClassIDParam 解析教师端班级码接口的 class_id：优先 query，其次 JSON body，
// 都没有时返回 0（由服务层回落到用户的 active_class_id），语义同 Laravel `$request->input()`。
func displayClassIDParam(c *gin.Context) uint {
	if id := queryID(c, "class_id"); id != 0 {
		return id
	}
	var req struct {
		ClassID *uint `json:"class_id"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.ClassID != nil {
		return *req.ClassID
	}
	return 0
}
func displayClassID(c *gin.Context) (uint, bool) {
	classID := middleware.DisplayClassID(c)
	if classID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token 无效或已过期"})
		return 0, false
	}
	return classID, true
}

// ============================================================
// 显示端 — 班级码登录
// ============================================================

// displayLoginRequest 班级码登录请求体。
type displayLoginRequest struct {
	Code string `json:"code" binding:"required,max=12"`
}

// DisplayLogin 班级码登录，返回 disp_ token（有效期 86400 秒）。
func (h *Handlers) DisplayLogin(c *gin.Context) {
	var req displayLoginRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.display.Login(req.Code, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// ============================================================
// 显示端 — 只读接口
// ============================================================

// DisplayInitialData 大屏初始全量数据。
func (h *Handlers) DisplayInitialData(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.InitialData(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayTimetable 大屏课表（今天星期几 + 科目/节次/排课）。
func (h *Handlers) DisplayTimetable(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.Timetable(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayExportCses 大屏免登录导出 CSES YAML（ClassIsland 可直接拉取该 URL）。
func (h *Handlers) DisplayExportCses(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	yaml, class, err := h.display.ExportCses(classID)
	if err != nil {
		fail(c, err)
		return
	}

	// 响应头与 Laravel 一致：application/x-yaml; charset=UTF-8 +
	// attachment; filename="<rawurlencode(班级名 + "-课表.cses.yaml")>" + Cache-Control: no-store。
	c.Header("Content-Disposition", `attachment; filename="`+rawURLEncode(class.Name+"-课表.cses.yaml")+`"`)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/x-yaml; charset=UTF-8", []byte(yaml))
}

// DisplayClassSettings 教室端班级设置（pet_series）。
func (h *Handlers) DisplayClassSettings(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.ClassSettings(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayDashboard 教室端班级总览。
func (h *Handlers) DisplayDashboard(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.ClassroomDashboard(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayStudents 教室端学生列表（含宠物信息）。
func (h *Handlers) DisplayStudents(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.ClassroomStudents(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayScoreRules 教室端积分规则。
func (h *Handlers) DisplayScoreRules(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.ScoreRules(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayPetsOverview 教室端宠物概览。
func (h *Handlers) DisplayPetsOverview(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.ClassroomPetsOverview(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// DisplayLeaderboard 大屏排行榜（本班前 20）。
func (h *Handlers) DisplayLeaderboard(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	data, err := h.display.QuickLeaderboard(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// ============================================================
// 教师端 / 管理员端 — 班级码
// ============================================================

// TeacherGetDisplayCode 获取当前班级的大屏码（未指定 class_id 时取当前激活班级）。
func (h *Handlers) TeacherGetDisplayCode(c *gin.Context) {
	u := middleware.CurrentUser(c)
	data, err := h.display.TeacherDisplayCode(u, queryID(c, "class_id"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherRefreshDisplayCode 刷新当前班级的大屏码。
// class_id 取值顺序同 Laravel `$request->input('class_id', ...)`：先 query，再 JSON body，
// 都没有时回落到用户的 active_class_id。
func (h *Handlers) TeacherRefreshDisplayCode(c *gin.Context) {
	u := middleware.CurrentUser(c)
	data, err := h.display.TeacherRefreshDisplayCode(u, displayClassIDParam(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// AdminGetDisplayCode 管理员获取某班大屏码。
func (h *Handlers) AdminGetDisplayCode(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}
	data, err := h.display.AdminDisplayCode(u, classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// AdminRefreshDisplayCode 管理员刷新某班大屏码。
func (h *Handlers) AdminRefreshDisplayCode(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classID, valid := paramID(c, "id")
	if !valid {
		return
	}
	data, err := h.display.AdminRefreshDisplayCode(u, classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// AdminResetDisplayCodes 批量重置本校全部班级码（可传统一字母前缀）。
// prefix 取值顺序同 Laravel `$request->input('prefix', ...)`：先 query，再 JSON body。
func (h *Handlers) AdminResetDisplayCodes(c *gin.Context) {
	u := middleware.CurrentUser(c)

	prefix := strings.TrimSpace(c.Query("prefix"))
	if prefix == "" {
		var req struct {
			Prefix *string `json:"prefix"`
		}
		_ = c.ShouldBindJSON(&req)
		if req.Prefix != nil {
			prefix = strings.TrimSpace(*req.Prefix)
		}
	}

	data, err := h.display.AdminResetDisplayCodes(u, prefix)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// rawURLEncode 等价 PHP rawurlencode（RFC 3986：仅 A-Za-z0-9-_.~ 不转义，
// 其余字节按 UTF-8 逐字节大写十六进制转义；Go 的 url.QueryEscape/PathEscape 语义不同）。
func rawURLEncode(s string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '_' || ch == '.' || ch == '~' {
			b.WriteByte(ch)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(hex[ch>>4])
		b.WriteByte(hex[ch&0x0f])
	}
	return b.String()
}
