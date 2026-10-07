// AI 配置与对话处理器。
//
// 覆盖 12 个动作（管理端 6 / 教师端 4 / 大屏端 2），逐一对应 Laravel：
// SchoolAdminController 的 getAiSettings/saveAiSettings/toggleAi/getAiUsage/fetchAiModels/testAiProvider、
// TeacherController 的 aiConfig/aiChat/getAiCommands/getAiUsage、
// DisplayController 的 aiSettingsCheck/aiChat。
//
// 有意差异：响应统一走本仓库的 {data, message:"ok"} 信封（同其余已移植接口），
// 文案与校验状态码逐字对齐；saveAiSettings 的「参数错误」额外带 errors 字段（形状同 Laravel Validator）。
package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// readRawJSONBody 读取请求体原始字段表（空体视为空表，非法 JSON → 422）。
// 用 RawMessage 保留「键是否出现」与原始字面量，才能逐条对齐 Laravel 的 has()/nullable/boolean 语义。
func readRawJSONBody(c *gin.Context) (map[string]json.RawMessage, bool) {
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "请求参数格式错误"})
		return nil, false
	}
	raw := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(buf)) == 0 {
		return raw, true
	}
	if err := json.Unmarshal(buf, &raw); err != nil || raw == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "请求参数格式错误"})
		return nil, false
	}
	return raw, true
}

// ============================================================
// 管理端 /api/v1/admin/ai/*
// ============================================================

// AdminAISettings GET /admin/ai/settings（首次访问惰性建行）。
func (h *Handlers) AdminAISettings(c *gin.Context) {
	u := middleware.CurrentUser(c)
	view, err := h.aiAdmin.Settings(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// AdminAISaveSettings PUT /admin/ai/settings。
func (h *Handlers) AdminAISaveSettings(c *gin.Context) {
	u := middleware.CurrentUser(c)
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	fieldErrs, err := h.aiAdmin.SaveSettings(u.SchoolID, raw)
	if err != nil {
		fail(c, err)
		return
	}
	if len(fieldErrs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "参数错误", "errors": fieldErrs})
		return
	}
	okMessage(c, "AI 设置已保存")
}

// AdminAIToggle POST /admin/ai/toggle（enabled 必填且为 Laravel boolean，否则 422「参数错误」）。
func (h *Handlers) AdminAIToggle(c *gin.Context) {
	u := middleware.CurrentUser(c)
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	rawEnabled, present := raw["enabled"]
	var enabled bool
	if present {
		enabled, valid = services.ParseLaravelBool(rawEnabled)
	}
	if !present || !valid {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "参数错误"})
		return
	}

	if _, err := h.aiAdmin.Toggle(u.SchoolID, enabled); err != nil {
		fail(c, err)
		return
	}
	message := "AI 功能已关闭"
	if enabled {
		message = "AI 功能已开启"
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "data": gin.H{"enabled": enabled}})
}

// AdminAIUsage GET /admin/ai/usage?days=7。
func (h *Handlers) AdminAIUsage(c *gin.Context) {
	u := middleware.CurrentUser(c)
	view, err := h.aiAdmin.Usage(u.SchoolID, aiDaysParam(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// AdminAIFetchModels POST /admin/ai/fetch-models。
func (h *Handlers) AdminAIFetchModels(c *gin.Context) {
	u := middleware.CurrentUser(c)
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	view, err := h.aiAdmin.FetchModels(u.SchoolID, services.RawScalarString(raw["provider_id"]))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// AdminAITest POST /admin/ai/test。
func (h *Handlers) AdminAITest(c *gin.Context) {
	u := middleware.CurrentUser(c)
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	view, err := h.aiAdmin.Test(u.SchoolID, services.RawScalarString(raw["provider_id"]))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// AdminAIProviderOfficial POST /admin/ai/provider-official（官方账单直查）。
//
// 忠实移植 SchoolAdminController::getAiProviderOfficial：provider_id 必传（缺 → 422「缺少供应商」）；
// 未创建 AI 配置 / 供应商不存在 / 未配置 API Key 各自 422；官方接口不可用时 data 里 official=false
// 与说明文案（形状见 services.AIProviderOfficialFallback / AIProviderOfficialData）。
func (h *Handlers) AdminAIProviderOfficial(c *gin.Context) {
	u := middleware.CurrentUser(c)
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	data, err := h.aiBilling.AIProviderOfficial(u.SchoolID, services.RawScalarString(raw["provider_id"]))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// ============================================================
// 教师端 /api/v1/teacher/ai/*
// ============================================================

// TeacherAIConfig GET /teacher/ai/config。
func (h *Handlers) TeacherAIConfig(c *gin.Context) {
	u := middleware.CurrentUser(c)
	view, err := h.ai.ConfigFor(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// TeacherAIChat POST /teacher/ai/chat（message 必填且 ≤2000 字，model 可选 ≤100 字）。
func (h *Handlers) TeacherAIChat(c *gin.Context) {
	u := middleware.CurrentUser(c)
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	rawMessage, present := raw["message"]
	if !present {
		fail(c, services.ErrUnprocessable("message 必填"))
		return
	}
	message, isString := rawJSONString(rawMessage)
	if !isString {
		fail(c, services.ErrUnprocessable("message 必须是字符串"))
		return
	}
	if !validAILength(message, 2000) {
		fail(c, services.ErrUnprocessable("message 不能超过 2000 字"))
		return
	}

	model := ""
	if rawModel, ok := raw["model"]; ok && !rawJSONNull(rawModel) {
		value, isString := rawJSONString(rawModel)
		if !isString {
			fail(c, services.ErrUnprocessable("model 必须是字符串"))
			return
		}
		if !validAILength(value, 100) {
			fail(c, services.ErrUnprocessable("model 不能超过 100 字"))
			return
		}
		model = value
	}

	reply, err := h.ai.Chat(u, message, model)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"reply": reply})
}

// TeacherAIUsage GET /teacher/ai/usage。
func (h *Handlers) TeacherAIUsage(c *gin.Context) {
	u := middleware.CurrentUser(c)
	view, err := h.ai.UsageFor(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// TeacherAICommands GET /teacher/ai/commands。
func (h *Handlers) TeacherAICommands(c *gin.Context) {
	ok(c, h.ai.Commands())
}

// ============================================================
// 大屏端 /api/v1/display/ai/*（班级码 token，DisplayAuth 中间件）
// ============================================================

// DisplayAISettings GET /display/ai/settings。
func (h *Handlers) DisplayAISettings(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	view, err := h.ai.DisplaySettingsCheck(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, view)
}

// DisplayAIChat POST /display/ai/chat（question 必填，空/\"0\" → 422「请输入问题」）。
func (h *Handlers) DisplayAIChat(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	raw, valid := readRawJSONBody(c)
	if !valid {
		return
	}

	result, err := h.ai.DisplayChat(classID, services.RawScalarString(raw["question"]))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// ============================================================
// 局部 JSON 辅助（仅本文件使用）
// ============================================================

// rawJSONString 取 JSON 字符串值。
func rawJSONString(raw json.RawMessage) (string, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}

// rawJSONNull 判断原始值是否为 null。
func rawJSONNull(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) == 0 || string(trimmed) == "null"
}

// validAILength 判断字符串是否不超过 max 个字符（Laravel max 按字符计）。
func validAILength(s string, max int) bool {
	return len([]rune(s)) <= max
}

// aiDaysParam 解析 ?days=，语义同 Laravel `(int) $request->input('days', 7)`：
// 参数缺失取 7；存在但非整数（含空串）按 PHP 强转取 0。
func aiDaysParam(c *gin.Context) int {
	values, ok := c.Request.URL.Query()["days"]
	if !ok || len(values) == 0 {
		return 7
	}
	n, err := strconv.Atoi(values[0])
	if err != nil {
		return 0
	}
	return n
}
