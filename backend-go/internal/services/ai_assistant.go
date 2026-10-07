// 教师端 / 大屏端 AI 助教服务。
//
// 移植自 Laravel：
//   - App\Services\AiAssistantService（教师端 configFor / chat / usageFor / commands）；
//   - DisplayController::aiSettingsCheck / aiChat（大屏端，校验顺序、返回结构与教师端不同）。
//
// 供应商解析规则（多供应商优先，兼容旧版单供应商字段）逐条对齐：
//   - providers 数组中第一个 is_active 且 api_key 非空的供应商优先；
//   - 都没有时回退 settings.provider / api_key（旧版单供应商字段）。
package services

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// AIProvider 是一条供应商配置（JSON 对象；保留全部字段以便原样回写，等价 Laravel 的关联数组）。
type AIProvider = map[string]any

// AIAssistant 教师端 / 大屏端 AI 助教服务。
type AIAssistant struct {
	db *gorm.DB
	ai *AIService
}

// NewAIAssistant 创建助教服务。
func NewAIAssistant(db *gorm.DB, ai *AIService) *AIAssistant {
	return &AIAssistant{db: db, ai: ai}
}

// ============================================================
// providers JSON 与取值辅助（等价 PHP 的数组 / empty() / ?? 语义）
// ============================================================

// decodeAIProviders 解析 providers JSON 文本；空/非法/非数组时返回空切片（同 Laravel `$x ?: []`）。
func decodeAIProviders(raw string) []AIProvider {
	out := []AIProvider{}
	if strings.TrimSpace(raw) == "" {
		return out
	}
	var list []AIProvider
	if err := json.Unmarshal([]byte(raw), &list); err != nil || list == nil {
		return out
	}
	return list
}

// encodeAIProviders 序列化 providers（空切片固定输出 "[]"，避免落库成 JSON null）。
func encodeAIProviders(list []AIProvider) string {
	if list == nil {
		list = []AIProvider{}
	}
	buf, err := json.Marshal(list)
	if err != nil {
		return "[]"
	}
	return string(buf)
}

// aiTruthy 等价 PHP `!empty($v)`。
func aiTruthy(v any) bool {
	switch value := v.(type) {
	case nil:
		return false
	case bool:
		return value
	case string:
		return value != "" && value != "0"
	case float64:
		return value != 0
	case int:
		return value != 0
	case int64:
		return value != 0
	case json.Number:
		f, err := value.Float64()
		return err == nil && f != 0
	case []any:
		return len(value) > 0
	case map[string]any:
		return len(value) > 0
	}
	return true
}

// aiProviderString 读取供应商字符串字段；缺失/非字符串返回空串。
func aiProviderString(p AIProvider, key string) string {
	s, _ := p[key].(string)
	return s
}

// aiProviderStringOr 等价 PHP `$p[$key] ?? $def`：键不存在或值为 null 时取默认值（空串不触发默认）。
func aiProviderStringOr(p AIProvider, key, def string) string {
	v, ok := p[key]
	if !ok || v == nil {
		return def
	}
	if s, ok := v.(string); ok {
		return s
	}
	return def
}

// aiProviderNumber 读取供应商数值字段（数值或数字字符串），其余为 0。
func aiProviderNumber(p AIProvider, key string) float64 {
	return aiToFloat(p[key])
}

// aiToFloat 把 JSON 数值/数字字符串转 float64（其余为 0，等价 PHP 的 (float) 转换）。
func aiToFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0
		}
		return f
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		if err != nil {
			return 0
		}
		return f
	}
	return 0
}

// aiLocalCost 本地精确计费（同 AiBillingService::calculateCost：按每百万 token 单价，
// round(prompt/1e6*input + completion/1e6*output, 6)）。
//
// 说明：AiBilling 的 7 个「官方账单/余额」驱动不在本批移植范围（属 provider-official 另批），
// 但 recordUsage 本身就是这个本地计算，故在此内联实现，保证 ai_conversations.cost 口径一致。
func aiLocalCost(provider AIProvider, promptTokens, completionTokens int) float64 {
	inputPrice := aiProviderNumber(provider, "input_price_per_m")
	outputPrice := aiProviderNumber(provider, "output_price_per_m")

	value := float64(promptTokens)/1_000_000*inputPrice + float64(completionTokens)/1_000_000*outputPrice
	return math.Round(value*1_000_000) / 1_000_000
}

// firstActiveProvider 返回首个 is_active 且 api_key 非空的供应商；没有则 nil。
func firstActiveProvider(list []AIProvider) AIProvider {
	for _, p := range list {
		if aiTruthy(p["is_active"]) && aiTruthy(p["api_key"]) {
			return p
		}
	}
	return nil
}

// chatProvider 解析对话用供应商：多供应商优先；回退旧版单供应商字段时补 api_base 与默认模型。
func chatProvider(setting *models.AISetting) AIProvider {
	if setting == nil {
		return nil
	}
	if p := firstActiveProvider(decodeAIProviders(setting.Providers)); p != nil {
		return p
	}
	if setting.APIKey != "" {
		id := setting.Provider
		if id == "" {
			id = "openai"
		}
		model := setting.Model
		if model == "" {
			model = "gpt-3.5-turbo"
		}
		return AIProvider{"id": id, "api_key": setting.APIKey, "api_base": setting.APIBase, "model": model}
	}
	return nil
}

// usageProvider 解析用量展示用供应商（旧版回退只补 id 与 model，与 chat 的字段集不同）。
func usageProvider(setting *models.AISetting) AIProvider {
	if setting == nil {
		return nil
	}
	if p := firstActiveProvider(decodeAIProviders(setting.Providers)); p != nil {
		return p
	}
	if setting.APIKey != "" {
		id := setting.Provider
		if id == "" {
			id = "openai"
		}
		return AIProvider{"id": id, "model": setting.Model}
	}
	return nil
}

// allowedModelsFor 供应商可选模型白名单：主模型 + models 多选 + model_map 请求侧键名（去重保序）。
func allowedModelsFor(p AIProvider) []string {
	candidates := []string{}
	if def := aiProviderString(p, "model"); def != "" {
		candidates = append(candidates, def)
	}
	if list, ok := p["models"].([]any); ok {
		for _, item := range list {
			if m, ok := item.(string); ok && m != "" {
				candidates = append(candidates, m)
			}
		}
	}
	if mm, ok := p["model_map"].(map[string]any); ok {
		for key := range mm {
			if key != "" {
				candidates = append(candidates, key)
			}
		}
	}

	seen := map[string]bool{}
	out := []string{}
	for _, m := range candidates {
		if seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	return out
}

// aiStringPtr 返回字符串指针（用于可空列的写入与视图）。
func aiStringPtr(s string) *string { return &s }

// aiUintPtr 返回 uint 指针。
func aiUintPtr(v uint) *uint { return &v }

// ============================================================
// 教师端
// ============================================================

// AIConfigView 教师端 AI 配置态（前端据此隐藏 AI 助教入口）。
type AIConfigView struct {
	Enabled bool `json:"enabled"`
}

// loadSetting 读取本校 AI 配置；不存在时返回 (nil, nil)（同 Laravel `where(...)->first()`）。
func (s *AIAssistant) loadSetting(schoolID uint) (*models.AISetting, error) {
	var setting models.AISetting
	err := s.db.Where("school_id = ?", schoolID).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// ConfigFor 计算本校 AI 是否可用（开关已开 且 至少一个启用供应商配了 Key，或旧版单供应商字段齐全）。
func (s *AIAssistant) ConfigFor(u *models.User) (*AIConfigView, error) {
	setting, err := s.loadSetting(u.SchoolID)
	if err != nil {
		return nil, err
	}

	enabled := false
	if setting != nil && setting.Enabled {
		hasProviderKey := firstActiveProvider(decodeAIProviders(setting.Providers)) != nil
		hasLegacyKey := setting.Provider != "" && setting.APIKey != ""
		enabled = hasProviderKey || hasLegacyKey
	}
	return &AIConfigView{Enabled: enabled}, nil
}

// Chat 教师对话：解析供应商 → 落会话（status=pending）→ 调 AI → 本地计费 → 回写用量/计数器。
// 未启用、未配置供应商等软失败返回提示文案（HTTP 200，与 Laravel 行为一致）。
// model 为教师指定模型（须在供应商白名单内，命中 model_map 时重定向到上游模型）。
func (s *AIAssistant) Chat(u *models.User, message, model string) (string, error) {
	setting, err := s.loadSetting(u.SchoolID)
	if err != nil {
		return "", err
	}
	if setting == nil || !setting.Enabled {
		return "AI 功能未启用，请联系管理员配置", nil
	}

	provider := chatProvider(setting)
	if provider == nil {
		return "请先在 AI 中心配置并启用一个供应商", nil
	}

	// 模型解析（New API 式）：默认主模型 → 白名单校验 → model_map 重定向
	defaultModel := aiProviderString(provider, "model")
	if defaultModel == "" {
		defaultModel = "gpt-3.5-turbo"
	}
	allowed := allowedModelsFor(provider)
	requested := model
	if requested == "" {
		requested = defaultModel
	}
	if !containsString(allowed, requested) {
		requested = defaultModel
	}
	upstreamModel := requested
	if mm, ok := provider["model_map"].(map[string]any); ok {
		if v, exists := mm[requested]; exists && v != nil {
			if s, ok := v.(string); ok {
				upstreamModel = s
			}
		}
	}

	var classID *uint
	if id := u.SettingUint(ActiveClassSettingKey); id != 0 {
		classID = aiUintPtr(id)
	}

	providerID := aiProviderString(provider, "id")
	conversation := models.AIConversation{
		SchoolID:    u.SchoolID,
		ClassID:     classID,
		StudentName: aiStringPtr("教师"),
		Provider:    aiStringPtr(providerID),
		Question:    message,
		Status:      "pending",
	}
	if err := s.db.Create(&conversation).Error; err != nil {
		return "", err
	}

	result, chatErr := s.ai.Chat(
		providerID,
		aiProviderString(provider, "api_key"),
		upstreamModel,
		message,
		aiProviderString(provider, "api_base"),
		setting.MaxTokens,
	)

	reply := result.Answer
	if chatErr != nil {
		reply = "AI 服务暂时不可用"
		result = AIResult{}
	}

	if err := s.finishConversation(setting, provider, &conversation, reply, result); err != nil {
		return "", err
	}
	return reply, nil
}

// AIUsageView 教师端 AI 用量（前端 AIPage 每次发送后刷新）。
type AIUsageView struct {
	Configured    bool     `json:"configured"`
	Provider      *string  `json:"provider"`
	Model         *string  `json:"model"`
	Models        []string `json:"models"`
	TokensUsed    int      `json:"tokens_used"`
	EstimatedCost float64  `json:"estimated_cost"`
	Currency      string   `json:"currency"`
}

// UsageFor 教师端用量汇总（口径同 AiAssistantService::usageFor）。
func (s *AIAssistant) UsageFor(u *models.User) (*AIUsageView, error) {
	setting, err := s.loadSetting(u.SchoolID)
	if err != nil {
		return nil, err
	}

	active := usageProvider(setting)

	estimatedCost := 0.0
	currency := "CNY"
	if setting != nil {
		for _, p := range decodeAIProviders(setting.Providers) {
			estimatedCost += aiProviderNumber(p, "estimated_cost")
			currency = aiProviderStringOr(p, "currency", currency)
		}
	}

	view := &AIUsageView{
		Configured:    setting != nil && setting.Enabled && active != nil,
		Models:        []string{},
		EstimatedCost: estimatedCost,
		Currency:      currency,
	}
	if setting != nil {
		view.TokensUsed = setting.TokensUsed
	}
	if active != nil {
		view.Models = allowedModelsFor(active)
		if v, ok := active["id"]; ok && v != nil {
			view.Provider = aiStringPtr(aiProviderString(active, "id"))
		}
		if v, ok := active["model"]; ok && v != nil {
			view.Model = aiStringPtr(aiProviderString(active, "model"))
		}
	}
	if view.Provider == nil && setting != nil {
		view.Provider = aiStringPtr(setting.Provider)
	}
	if view.Model == nil && setting != nil {
		view.Model = aiStringPtr(setting.Model)
	}
	return view, nil
}

// AICommand 教师端 AI 预设命令。
type AICommand struct {
	Label  string `json:"label"`
	Prompt string `json:"prompt"`
}

// Commands 教师端预设命令清单（逐字对齐 AiAssistantService::commands）。
func (s *AIAssistant) Commands() []AICommand {
	return []AICommand{
		{Label: "📝 本周教学总结", Prompt: "请帮我写一份本周教学总结，包含本周教学目标、课堂情况、学生表现和下周教学计划。"},
		{Label: "🏅 积分规则建议", Prompt: "请根据班级日常情况，生成一套适合小学生的积分奖励规则建议。"},
		{Label: "🎯 班会活动方案", Prompt: "请设计一个有趣的小学生班会活动方案，包含活动目标、流程和所需材料。"},
		{Label: "📋 出练习题", Prompt: "请出一组适合本年级学生的练习题，包含题目和参考答案。"},
	}
}

// ============================================================
// 大屏端（班级码 token）
// ============================================================

// AIDisplaySettingsView 大屏端 AI 开关检查结果。
type AIDisplaySettingsView struct {
	Enabled bool `json:"enabled"`
}

// DisplaySettingsCheck 大屏端 AI 开关检查（同 DisplayController::aiSettingsCheck）：
// 班级或学校不存在时返回 enabled=false（不报错）。
func (s *AIAssistant) DisplaySettingsCheck(classID uint) (*AIDisplaySettingsView, error) {
	var class models.ClassRoom
	err := s.db.First(&class, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &AIDisplaySettingsView{Enabled: false}, nil
	}
	if err != nil {
		return nil, err
	}
	if !s.schoolExists(class.SchoolID) {
		return &AIDisplaySettingsView{Enabled: false}, nil
	}

	setting, err := s.loadSetting(class.SchoolID)
	if err != nil {
		return nil, err
	}
	return &AIDisplaySettingsView{Enabled: setting != nil && setting.Enabled}, nil
}

// AIDisplayChatResult 大屏端对话返回（只有 answer 与 tokens_used，同 Laravel）。
type AIDisplayChatResult struct {
	Answer     string `json:"answer"`
	TokensUsed int    `json:"tokens_used"`
}

// DisplayChat 大屏端 AI 对话（同 DisplayController::aiChat）。
// 校验顺序与 Laravel 一致：班级/学校（404）→ 开关（403）→ 供应商（403）→ 问题非空（422）。
// 与教师端的差异：不做模型白名单/model_map 解析、不写 class_id 之外的学生信息（学生恒为「匿名」）、
// 传输层异常时把错误信息拼进回答（「AI 服务调用失败：…」）。
func (s *AIAssistant) DisplayChat(classID uint, question string) (*AIDisplayChatResult, error) {
	var class models.ClassRoom
	err := s.db.First(&class, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("班级不存在")
	}
	if err != nil {
		return nil, err
	}
	if !s.schoolExists(class.SchoolID) {
		return nil, ErrNotFound("班级不存在")
	}

	setting, err := s.loadSetting(class.SchoolID)
	if err != nil {
		return nil, err
	}
	if setting == nil || !setting.Enabled {
		return nil, ErrForbidden("AI 功能未开启")
	}

	provider := chatProvider(setting)
	if provider == nil {
		return nil, ErrForbidden("请先在 AI 中心配置并启用一个供应商")
	}

	if phpEmptyString(question) {
		return nil, ErrUnprocessable("请输入问题")
	}

	providerID := aiProviderString(provider, "id")
	conversation := models.AIConversation{
		SchoolID:    class.SchoolID,
		ClassID:     aiUintPtr(classID),
		StudentName: aiStringPtr("匿名"),
		Provider:    aiStringPtr(providerID),
		Question:    question,
		Status:      "pending",
	}
	if err := s.db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	model := aiProviderString(provider, "model")
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	result, chatErr := s.ai.Chat(
		providerID,
		aiProviderString(provider, "api_key"),
		model,
		question,
		aiProviderString(provider, "api_base"),
		setting.MaxTokens,
	)

	answer := result.Answer
	if chatErr != nil {
		answer = "AI 服务调用失败：" + chatErr.Error()
		result = AIResult{}
	}

	if err := s.finishConversation(setting, provider, &conversation, answer, result); err != nil {
		return nil, err
	}
	return &AIDisplayChatResult{Answer: answer, TokensUsed: result.TokensUsed}, nil
}

// ============================================================
// 内部辅助
// ============================================================

// schoolExists 判断学校行是否存在（Laravel 用 with('school') 判断关联）。
func (s *AIAssistant) schoolExists(schoolID uint) bool {
	var count int64
	if err := s.db.Model(&models.School{}).Where("id = ?", schoolID).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

// finishConversation 回写会话结果（answer/tokens/cost/currency/status=completed），
// 并在有 token 消耗时累加学校 tokens_used 与该供应商的计数器（与 Laravel 顺序一致）。
func (s *AIAssistant) finishConversation(setting *models.AISetting, provider AIProvider, conversation *models.AIConversation, reply string, result AIResult) error {
	cost := aiLocalCost(provider, result.PromptTokens, result.CompletionTokens)
	currency := aiProviderStringOr(provider, "currency", "CNY")

	updates := map[string]any{
		"answer":            reply,
		"tokens_used":       result.TokensUsed,
		"prompt_tokens":     result.PromptTokens,
		"completion_tokens": result.CompletionTokens,
		"cost":              cost,
		"currency":          currency,
		"status":            "completed",
	}
	if err := s.db.Model(&models.AIConversation{}).Where("id = ?", conversation.ID).Updates(updates).Error; err != nil {
		return err
	}

	conversation.Answer = aiStringPtr(reply)
	conversation.TokensUsed = result.TokensUsed
	conversation.PromptTokens = result.PromptTokens
	conversation.CompletionTokens = result.CompletionTokens
	conversation.Cost = cost
	conversation.Currency = currency
	conversation.Status = "completed"

	if result.TokensUsed <= 0 {
		return nil
	}
	return s.accumulateUsage(setting, aiProviderString(provider, "id"), result.TokensUsed, cost, currency)
}

// accumulateUsage 累计学校 tokens_used 与 providers[] 中该供应商的计数器（只命中第一条）。
func (s *AIAssistant) accumulateUsage(setting *models.AISetting, providerID string, tokensUsed int, cost float64, currency string) error {
	providers := decodeAIProviders(setting.Providers)
	for _, p := range providers {
		if aiProviderString(p, "id") != providerID {
			continue
		}
		p["tokens_used"] = aiProviderNumber(p, "tokens_used") + float64(tokensUsed)
		p["total_calls"] = aiProviderNumber(p, "total_calls") + 1
		p["estimated_cost"] = aiProviderNumber(p, "estimated_cost") + cost
		p["currency"] = currency
		break
	}

	updates := map[string]any{
		"tokens_used": gorm.Expr("tokens_used + ?", tokensUsed),
		"providers":   encodeAIProviders(providers),
	}
	return s.db.Model(&models.AISetting{}).Where("id = ?", setting.ID).Updates(updates).Error
}

// phpEmptyString 等价 PHP `empty($str)`：空串与 "0" 均为真。
func phpEmptyString(s string) bool {
	return s == "" || s == "0"
}

// truncateRunes 按字符（而非字节）截断，同 PHP mb_substr。
func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}
