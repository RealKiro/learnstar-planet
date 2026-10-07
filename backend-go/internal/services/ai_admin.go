// 管理端 AI 中心服务。
//
// 移植自 Laravel SchoolAdminController 的 AI 动作：
// getAiSettings / saveAiSettings / toggleAi / getAiUsage / fetchAiModels / testAiProvider
// （provider-official 及其 AiBilling 官方账单驱动不在本批范围）。
//
// 逐条对齐点：
//   - GET settings 返回 {enabled, tokens_used, tokens_limit, max_tokens, providers}，
//     **不含** 顶层 api_key/api_base/model/provider（Laravel 只回这几个字段），
//     但 providers[] 里的 api_key 是**原样**返回的（不做掩码/截断）；首次访问 firstOrCreate 建行。
//   - PUT settings：providers 数组「按 id 合并」——空 api_key 保留原值、计数器取新旧较大值、
//     currency 缺失时沿用旧值；**请求未带 providers 字段时按 Laravel 语义把 providers 置空**；
//     enabled/max_tokens/tokens_limit 仅在该键出现时更新（enabled 走 Laravel boolean 规则）。
//   - GET usage：时间范围 created_at >= now()-days、daily_usage 按日期分组升序、recent_logs 取 50 条降序、
//     by_provider 直接读 providers 里的计数器、total_conversations 为**全量**计数（不限时间范围）。
//   - POST fetch-models / POST test：都要求供应商已配 api_key（422 文案逐字对齐），test 不写库、不计费。
//   - POST toggle：enabled 必填且为 Laravel boolean，否则 422「参数错误」。
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// AIAdmin 管理端 AI 中心服务。
type AIAdmin struct {
	db *gorm.DB
	ai *AIService
}

// NewAIAdmin 创建管理端 AI 服务。
func NewAIAdmin(db *gorm.DB, ai *AIService) *AIAdmin {
	return &AIAdmin{db: db, ai: ai}
}

// FieldErrors 请求参数校验错误（形状对齐 Laravel Validator::errors()：字段 → 文案列表）。
type FieldErrors map[string][]string

// AISettingView 管理端 AI 设置视图（同 SchoolAdminController::getAiSettings 的 data）。
type AISettingView struct {
	Enabled     bool         `json:"enabled"`
	TokensUsed  int          `json:"tokens_used"`
	TokensLimit int          `json:"tokens_limit"`
	MaxTokens   int          `json:"max_tokens"`
	Providers   []AIProvider `json:"providers"`
}

// Settings 读取（必要时惰性创建）本校 AI 设置。
func (a *AIAdmin) Settings(schoolID uint) (*AISettingView, error) {
	setting, err := a.firstOrCreate(schoolID)
	if err != nil {
		return nil, err
	}
	return &AISettingView{
		Enabled:     setting.Enabled,
		TokensUsed:  setting.TokensUsed,
		TokensLimit: setting.TokensLimit,
		MaxTokens:   setting.MaxTokens,
		Providers:   decodeAIProviders(setting.Providers),
	}, nil
}

// firstOrCreate 取本校设置；不存在则按 Laravel 的默认值建行。
func (a *AIAdmin) firstOrCreate(schoolID uint) (*models.AISetting, error) {
	setting, err := a.loadExisting(schoolID)
	if err != nil {
		return nil, err
	}
	if setting != nil {
		return setting, nil
	}

	setting = &models.AISetting{
		SchoolID:    schoolID,
		Enabled:     false,
		Provider:    "openai",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   2000,
		TokensUsed:  0,
		TokensLimit: 1000000,
	}
	if err := a.db.Create(setting).Error; err != nil {
		return nil, err
	}
	// enabled 的库默认值为 false，GORM 会跳过硬编码的零值 false（库默认同样为 false，语义一致）；
	// 这里再回读一次，确保 enabled/provider/model 等默认值在内存与库中一致。
	var reloaded models.AISetting
	if err := a.db.First(&reloaded, setting.ID).Error; err != nil {
		return nil, err
	}
	return &reloaded, nil
}

// loadExisting 读取本校设置；不存在返回 (nil, nil)。
func (a *AIAdmin) loadExisting(schoolID uint) (*models.AISetting, error) {
	var setting models.AISetting
	err := a.db.Where("school_id = ?", schoolID).First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// SaveSettings 保存 AI 设置（raw 为请求体原始 JSON 字段表，nil 值表示该键未出现）。
// 返回非空 FieldErrors 表示校验失败（HTTP 422「参数错误」+ errors），此时不落库。
func (a *AIAdmin) SaveSettings(schoolID uint, raw map[string]json.RawMessage) (FieldErrors, error) {
	errs := FieldErrors{}

	var enabled *bool
	if v, ok := raw["enabled"]; ok {
		if b, valid := ParseLaravelBool(v); valid {
			enabled = &b
		} else {
			errs["enabled"] = []string{"enabled 必须是布尔值"}
		}
	}

	var maxTokens *int
	if v, ok := raw["max_tokens"]; ok {
		if isNullJSON(v) {
			zero := 0 // Laravel：nullable 通过后仍走 has() 分支，(int) null = 0
			maxTokens = &zero
		} else if n, valid := aiJSONInt(v); !valid {
			errs["max_tokens"] = []string{"max_tokens 必须是整数"}
		} else if n < 100 || n > 32000 {
			errs["max_tokens"] = []string{"max_tokens 必须在 100-32000 之间"}
		} else {
			maxTokens = &n
		}
	}

	var tokensLimit *int
	if v, ok := raw["tokens_limit"]; ok {
		if isNullJSON(v) {
			zero := 0
			tokensLimit = &zero
		} else if n, valid := aiJSONInt(v); !valid {
			errs["tokens_limit"] = []string{"tokens_limit 必须是整数"}
		} else if n < 0 {
			errs["tokens_limit"] = []string{"tokens_limit 不能小于 0"}
		} else {
			tokensLimit = &n
		}
	}

	providers := []AIProvider{}
	if v, ok := raw["providers"]; ok && !isNullJSON(v) {
		elems, isArray := aiJSONArray(v)
		if !isArray {
			errs["providers"] = []string{"providers 必须是数组"}
		} else {
			for i, elem := range elems {
				obj, isObject := aiJSONObject(elem)
				if !isObject {
					errs[fmt.Sprintf("providers.%d", i)] = []string{"必须是对象"}
					continue
				}
				validateAIProviderEntry(i, obj, errs)

				var decoded AIProvider
				if err := json.Unmarshal(elem, &decoded); err != nil {
					errs[fmt.Sprintf("providers.%d", i)] = []string{"必须是对象"}
					continue
				}
				providers = append(providers, decoded)
			}
		}
	}

	if len(errs) > 0 {
		return errs, nil
	}

	setting, err := a.firstOrCreate(schoolID)
	if err != nil {
		return nil, err
	}

	// 按 id 合并旧配置：保留空 api_key、计数器取较大值、currency 缺失时沿用。
	existing := decodeAIProviders(setting.Providers)
	for _, p := range providers {
		for _, old := range existing {
			if aiProviderString(old, "id") != aiProviderString(p, "id") {
				continue
			}
			if !aiTruthy(p["api_key"]) && aiTruthy(old["api_key"]) {
				p["api_key"] = old["api_key"]
			}
			merges := []string{"tokens_used", "total_calls", "estimated_cost", "balance"}
			for _, counter := range merges {
				oldValue, hasOld := old[counter]
				if !hasOld {
					continue
				}
				if current, hasCurrent := p[counter]; hasCurrent && aiToFloat(current) > aiToFloat(oldValue) {
					p[counter] = current
				} else {
					p[counter] = oldValue
				}
			}
			if !aiTruthy(p["currency"]) && aiTruthy(old["currency"]) {
				p["currency"] = old["currency"]
			}
			break
		}
	}

	updates := map[string]any{"providers": encodeAIProviders(providers)}
	if enabled != nil {
		updates["enabled"] = *enabled
	}
	if maxTokens != nil {
		updates["max_tokens"] = *maxTokens
	}
	if tokensLimit != nil {
		updates["tokens_limit"] = *tokensLimit
	}
	if err := a.db.Model(&models.AISetting{}).Where("id = ?", setting.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

// Toggle 切换 AI 开关（返回落库后的 enabled，响应文案由 handler 拼装）。
func (a *AIAdmin) Toggle(schoolID uint, enabled bool) (bool, error) {
	setting, err := a.firstOrCreate(schoolID)
	if err != nil {
		return false, err
	}
	if err := a.db.Model(&models.AISetting{}).Where("id = ?", setting.ID).
		Update("enabled", enabled).Error; err != nil {
		return false, err
	}
	return enabled, nil
}

// ============================================================
// 用量统计
// ============================================================

// AIDailyUsageView daily_usage 的一行（按日期分组的 token / 次数 / 费用）。
type AIDailyUsageView struct {
	Date   string  `json:"date"`
	Tokens int64   `json:"tokens"`
	Count  int64   `json:"count"`
	Cost   float64 `json:"cost"`
}

// AIProviderUsageView by_provider 的一项（直接读 providers[].计数器，缺失取 0 / CNY）。
type AIProviderUsageView struct {
	Tokens        any    `json:"tokens"`
	TotalCalls    any    `json:"total_calls"`
	EstimatedCost any    `json:"estimated_cost"`
	CostPerToken  any    `json:"cost_per_token"`
	Currency      string `json:"currency"`
}

// AIConversationLogView recent_logs 的一条。
type AIConversationLogView struct {
	ID               uint    `json:"id"`
	StudentName      *string `json:"student_name"`
	ClassID          *uint   `json:"class_id"`
	Provider         *string `json:"provider"`
	Question         string  `json:"question"`
	Answer           *string `json:"answer"`
	TokensUsed       int     `json:"tokens_used"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	Cost             float64 `json:"cost"`
	Currency         string  `json:"currency"`
	CreatedAt        string  `json:"created_at"`
}

// AIAdminUsageView 管理端用量视图（同 SchoolAdminController::getAiUsage 的 data）。
type AIAdminUsageView struct {
	Enabled            bool                    `json:"enabled"`
	TokensUsed         int                     `json:"tokens_used"`
	TokensLimit        int                     `json:"tokens_limit"`
	EstimatedCost      float64                 `json:"estimated_cost"`
	TotalConversations int64                   `json:"total_conversations"`
	DailyUsage         []AIDailyUsageView      `json:"daily_usage"`
	ByProvider         AIProviderUsageMap      `json:"by_provider"`
	RecentLogs         []AIConversationLogView `json:"recent_logs"`
}

// AIProviderUsageMap by_provider 视图：非空时序列化为对象，空时序列化为 []（同 PHP 空数组的 JSON 形状）。
type AIProviderUsageMap map[string]AIProviderUsageView

// MarshalJSON 空表输出 "[]"，非空输出对象。
func (m AIProviderUsageMap) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("[]"), nil
	}
	type plain AIProviderUsageMap
	return json.Marshal(plain(m))
}

// Usage 管理端用量统计（days 默认为 7，与 Laravel `input('days', 7)` 一致）。
func (a *AIAdmin) Usage(schoolID uint, days int) (*AIAdminUsageView, error) {
	setting, err := a.loadExisting(schoolID)
	if err != nil {
		return nil, err
	}

	from := util.Now().AddDate(0, 0, -days)

	var inRange []models.AIConversation
	if err := a.db.Select("id", "tokens_used", "cost", "created_at").
		Where("school_id = ? AND created_at >= ?", schoolID, from).
		Find(&inRange).Error; err != nil {
		return nil, err
	}

	daily := map[string]*AIDailyUsageView{}
	dates := []string{}
	for i := range inRange {
		key := inRange[i].CreatedAt.In(util.Loc).Format("2006-01-02")
		row, ok := daily[key]
		if !ok {
			row = &AIDailyUsageView{Date: key}
			daily[key] = row
			dates = append(dates, key)
		}
		row.Tokens += int64(inRange[i].TokensUsed)
		row.Count++
		row.Cost += inRange[i].Cost
	}
	sort.Strings(dates)

	dailyUsage := []AIDailyUsageView{}
	for _, key := range dates {
		dailyUsage = append(dailyUsage, *daily[key])
	}

	var recent []models.AIConversation
	if err := a.db.Where("school_id = ? AND created_at >= ?", schoolID, from).
		Order("created_at DESC").Limit(50).Find(&recent).Error; err != nil {
		return nil, err
	}

	recentLogs := []AIConversationLogView{}
	for i := range recent {
		log := &recent[i]
		recentLogs = append(recentLogs, AIConversationLogView{
			ID:               log.ID,
			StudentName:      log.StudentName,
			ClassID:          log.ClassID,
			Provider:         log.Provider,
			Question:         log.Question,
			Answer:           log.Answer,
			TokensUsed:       log.TokensUsed,
			PromptTokens:     log.PromptTokens,
			CompletionTokens: log.CompletionTokens,
			Cost:             log.Cost,
			Currency:         log.Currency,
			CreatedAt:        log.CreatedAt.In(util.Loc).Format("2006-01-02 15:04:05"),
		})
	}

	byProvider := AIProviderUsageMap{}
	estimatedCost := 0.0
	if setting != nil {
		for _, p := range decodeAIProviders(setting.Providers) {
			id := aiProviderString(p, "id")
			if id == "" {
				continue
			}
			item := AIProviderUsageView{
				Tokens:        aiCounterOr(p, "tokens_used", 0),
				TotalCalls:    aiCounterOr(p, "total_calls", 0),
				EstimatedCost: aiCounterOr(p, "estimated_cost", 0),
				CostPerToken:  aiCounterOr(p, "cost_per_token", 0),
				Currency:      aiProviderStringOr(p, "currency", "CNY"),
			}
			byProvider[id] = item
			estimatedCost += aiToFloat(item.EstimatedCost)
		}
	}

	var total int64
	if err := a.db.Model(&models.AIConversation{}).Where("school_id = ?", schoolID).Count(&total).Error; err != nil {
		return nil, err
	}

	view := &AIAdminUsageView{
		Enabled:            setting != nil && setting.Enabled,
		EstimatedCost:      estimatedCost,
		TotalConversations: total,
		DailyUsage:         dailyUsage,
		ByProvider:         byProvider,
		RecentLogs:         recentLogs,
	}
	if setting != nil {
		view.TokensUsed = setting.TokensUsed
		view.TokensLimit = setting.TokensLimit
	}
	return view, nil
}

// aiCounterOr 读取计数器原值（缺失/为 null 时返回 def；保留原类型以便原样输出）。
func aiCounterOr(p AIProvider, key string, def any) any {
	v, ok := p[key]
	if !ok || v == nil {
		return def
	}
	return v
}

// ============================================================
// 模型列表 / 连通性测试
// ============================================================

// AIModelsView fetch-models 的 data。
type AIModelsView struct {
	Provider  string   `json:"provider"`
	Models    []string `json:"models"`
	FetchedAt string   `json:"fetched_at"`
}

// FetchModels 从官方 API 拉取某供应商的模型列表。
func (a *AIAdmin) FetchModels(schoolID uint, providerID string) (*AIModelsView, error) {
	if providerID == "" {
		return nil, ErrUnprocessable("缺少供应商")
	}

	setting, err := a.loadExisting(schoolID)
	if err != nil {
		return nil, err
	}
	provider := findProviderByID(setting, providerID)
	if provider == nil || !aiTruthy(provider["api_key"]) {
		return nil, ErrUnprocessable("该供应商未配置 API Key")
	}

	models := a.ai.ListModels(providerID, aiProviderString(provider, "api_key"), aiProviderString(provider, "api_base"))
	return &AIModelsView{
		Provider:  providerID,
		Models:    models,
		FetchedAt: util.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// AITestResult test 的 data。
// Reply 用指针 + omitempty：传输层异常分支 Laravel 不返回 reply 键；
// Error 为 any：成功分支编码为 null（Laravel 显式给 null），失败分支为文案。
type AITestResult struct {
	Success   bool    `json:"success"`
	LatencyMs int     `json:"latency_ms"`
	Reply     *string `json:"reply,omitempty"`
	Error     any     `json:"error"`
}

// aiTestFailedPhrases 与 Laravel testAiProvider 的失败兜底文案判定一致。
var aiTestFailedPhrases = []string{"AI 服务暂时不可用", "AI 服务不可用", "抱歉，无法回答"}

// Test 供应商连通性测试（向该供应商发一条最小对话；不写库、不计费）。
func (a *AIAdmin) Test(schoolID uint, providerID string) (*AITestResult, error) {
	if providerID == "" {
		return nil, ErrUnprocessable("缺少供应商")
	}

	setting, err := a.loadExisting(schoolID)
	if err != nil {
		return nil, err
	}
	provider := findProviderByID(setting, providerID)
	if provider == nil || !aiTruthy(provider["api_key"]) {
		return nil, ErrUnprocessable("该供应商未配置 API Key，请先填写并保存")
	}
	model := aiProviderString(provider, "model")
	if model == "" {
		return nil, ErrUnprocessable("该供应商未配置模型，请先填写并保存")
	}

	startedAt := time.Now()
	result, chatErr := a.ai.Chat(
		providerID,
		aiProviderString(provider, "api_key"),
		model,
		"你好",
		aiProviderString(provider, "api_base"),
		16,
	)
	latencyMs := int(time.Since(startedAt).Milliseconds())

	if chatErr != nil {
		return &AITestResult{Success: false, LatencyMs: latencyMs, Error: chatErr.Error()}, nil
	}

	answer := strings.TrimSpace(result.Answer)
	success := answer != "" && !containsString(aiTestFailedPhrases, answer)
	reply := truncateRunes(answer, 50)

	out := &AITestResult{Success: success, LatencyMs: latencyMs, Reply: &reply}
	if !success {
		out.Error = "供应商返回异常，请检查 API Key / 模型 / 地址"
	}
	return out, nil
}

// findProviderByID 在设置里按 id 找供应商（设置不存在时返回 nil）。
func findProviderByID(setting *models.AISetting, providerID string) AIProvider {
	if setting == nil {
		return nil
	}
	for _, p := range decodeAIProviders(setting.Providers) {
		if aiProviderString(p, "id") == providerID {
			return p
		}
	}
	return nil
}

// ============================================================
// 原始 JSON 取值与 Laravel 规则校验（供 handler 与 SaveSettings 复用）
// ============================================================

// ParseLaravelBool 解析 Laravel `boolean` 规则接受的取值：true/false/1/0/"1"/"0"（其余不合法）。
func ParseLaravelBool(raw json.RawMessage) (bool, bool) {
	switch string(bytes.TrimSpace(raw)) {
	case "true", "1", `"1"`:
		return true, true
	case "false", "0", `"0"`:
		return false, true
	}
	return false, false
}

// RawScalarString 等价 PHP `(string) $request->input($key, ”)`：字符串原样（null 视为空串）。
func RawScalarString(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return ""
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return ""
		}
		return s
	}
	switch string(trimmed) {
	case "null":
		return ""
	case "true":
		return "1"
	case "false":
		return ""
	}
	return string(trimmed)
}

// isNullJSON 判断原始 JSON 值是否为 null（或为空）。
func isNullJSON(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) == 0 || string(trimmed) == "null"
}

// aiJSONInt 解析整数（接受数字与数字字符串，要求为整数值）。
func aiJSONInt(raw json.RawMessage) (int, bool) {
	s := bytes.TrimSpace(raw)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	f, err := strconv.ParseFloat(string(s), 64)
	if err != nil || f != math.Trunc(f) {
		return 0, false
	}
	return int(f), true
}

// aiJSONNumeric 判断是否数值（Laravel numeric 规则：数字或数字字符串）。
func aiJSONNumeric(raw json.RawMessage) (float64, bool) {
	s := bytes.TrimSpace(raw)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	if len(s) == 0 {
		return 0, false
	}
	f, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

// aiJSONString 取 JSON 字符串值。
func aiJSONString(raw json.RawMessage) (string, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}

// aiJSONNonEmptyString 判断是否非空 JSON 字符串（Laravel required|string）。
func aiJSONNonEmptyString(raw json.RawMessage) bool {
	s, ok := aiJSONString(raw)
	return ok && s != ""
}

// aiJSONObject 解析 JSON 对象。
func aiJSONObject(raw json.RawMessage) (map[string]json.RawMessage, bool) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil || m == nil {
		return nil, false
	}
	return m, true
}

// aiJSONArray 解析 JSON 数组（null 视为不合法，由调用方先过滤）。
func aiJSONArray(raw json.RawMessage) ([]json.RawMessage, bool) {
	var list []json.RawMessage
	if err := json.Unmarshal(raw, &list); err != nil || list == nil {
		return nil, false
	}
	return list, true
}

// validateAIProviderEntry 按 Laravel saveAiSettings 的 providers.* 规则逐条校验。
func validateAIProviderEntry(idx int, obj map[string]json.RawMessage, errs FieldErrors) {
	prefix := fmt.Sprintf("providers.%d", idx)
	add := func(rule string, msg string) {
		errs[prefix+"."+rule] = append(errs[prefix+"."+rule], msg)
	}

	// required|string
	for _, key := range []string{"id", "label"} {
		if v, ok := obj[key]; !ok || !aiJSONNonEmptyString(v) {
			add(key, key+" 必填且为字符串")
		}
	}

	// nullable|string|max:N
	optionalString := func(key string, max int) {
		v, ok := obj[key]
		if !ok || isNullJSON(v) {
			return
		}
		s, isString := aiJSONString(v)
		if !isString {
			add(key, key+" 必须是字符串")
			return
		}
		if utf8.RuneCountInString(s) > max {
			add(key, fmt.Sprintf("%s 不能超过 %d 字", key, max))
		}
	}
	optionalString("api_key", 2000)
	optionalString("api_base", 500)
	optionalString("model", 100)
	optionalString("currency", 8)

	// nullable|array + 元素 string|max:100
	if v, ok := obj["models"]; ok && !isNullJSON(v) {
		list, isArray := aiJSONArray(v)
		if !isArray {
			add("models", "models 必须是数组")
		} else {
			for i, elem := range list {
				s, isString := aiJSONString(elem)
				if !isString {
					errs[fmt.Sprintf("%s.models.%d", prefix, i)] = []string{"必须是字符串"}
					continue
				}
				if utf8.RuneCountInString(s) > 100 {
					errs[fmt.Sprintf("%s.models.%d", prefix, i)] = []string{"不能超过 100 字"}
				}
			}
		}
	}

	// nullable|array|max:50 + nullable|string|max:100
	if v, ok := obj["model_map"]; ok && !isNullJSON(v) {
		mapping, isObject := aiJSONObject(v)
		if !isObject {
			add("model_map", "model_map 必须是对象")
		} else {
			if len(mapping) > 50 {
				add("model_map", "model_map 最多 50 项")
			}
			for key, elem := range mapping {
				if isNullJSON(elem) {
					continue
				}
				s, isString := aiJSONString(elem)
				if !isString {
					add("model_map", fmt.Sprintf("model_map.%s 必须是字符串", key))
					continue
				}
				if utf8.RuneCountInString(s) > 100 {
					add("model_map", fmt.Sprintf("model_map.%s 不能超过 100 字", key))
				}
			}
		}
	}

	// boolean（不可为 null）
	for _, key := range []string{"is_active", "billing_enabled"} {
		if v, ok := obj[key]; ok {
			if _, valid := ParseLaravelBool(v); !valid {
				add(key, key+" 必须是布尔值")
			}
		}
	}

	// nullable|numeric|min:0
	for _, key := range []string{"input_price_per_m", "output_price_per_m"} {
		v, ok := obj[key]
		if !ok || isNullJSON(v) {
			continue
		}
		n, isNumeric := aiJSONNumeric(v)
		if !isNumeric {
			add(key, key+" 必须是数字")
			continue
		}
		if n < 0 {
			add(key, key+" 不能小于 0")
		}
	}
}
