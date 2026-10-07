// AI 供应商官方账单（用量 / 余额）查询：POST /api/v1/admin/ai/provider-official。
//
// 忠实移植自 Laravel：
//   - App\Http\Controllers\Api\SchoolAdminController::getAiProviderOfficial（第 2407-2480 行）；
//   - App\Services\AiBilling\AiBillingService（driverFor / syncOfficialUsage / getBalance）；
//   - AiBilling\AiBillingDriver / AiBalance / AiUsageSnapshot 与 7 个驱动
//     （OpenAi / DeepSeek / Moonshot / SiliconFlow / OpenRouter / Relay(newapi,oneapi) / LocalPrecise）。
//
// 逐条对齐点：
//   - provider_id 为空 → 422「缺少供应商」；本校无 ai_settings 行 → 422「尚未创建 AI 配置」；
//     providers[] 中无该 id → 422「未找到该供应商配置」；api_key 为空 → 422「该供应商未配置 API Key，请先填写并保存」；
//   - usage 与 balance 都为 null → {provider_id, official:false, message:"该平台未提供可用的官方用量/余额接口，当前数据为本地估算"}；
//   - 否则 official:true，usage（可为 null：Moonshot / SiliconFlow / DeepSeek 无用量接口）与 balance
//     （可为 null：OpenAI 无余额接口）按各驱动解析结果输出，另带 fetched_at（Y-m-d H:i:s）；
//   - 官方余额写回 providers[].balance（Laravel `$p['balance'] = $balance->totalBalance`，取较大值合并由保存接口负责）；
//   - 各驱动的取值路径 / 币种 / 默认值 / 失败兜底（传输异常、HTTP>=400、响应体非法一律视为失败 → null）逐条照抄。
//
// 有意差异（均为可测试性/无 net/http 表层差异）：
//  1. Laravel 用 Http Facade；Go 端复用上一批 AIService 的 HTTPDoer（SetHTTPClient 注入），
//     测试一律用 httptest 假上游，不访问真实外网。
//  2. DeepSeek 余额端点与 OpenRouter key 端点在 Laravel 里是硬编码 URL（忽略 api_base），
//     httptest 无法拦截，故提供 SetEndpointOverrides 覆盖口（生产恒为 Laravel 的原 URL；
//     仅测试指向假上游）。OpenAI / Moonshot / SiliconFlow 沿用 api_base 覆盖，与 Laravel 一致。
//  3. Laravel 的 providers 存的是关联数组，Go 端为 JSON 文本列，写回时整体重新序列化（语义等价）。
package services

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// aiBillingTimeout 各官方账单接口的调用超时（同 Laravel `Http::timeout(30)`）。
const aiBillingTimeout = 30 * time.Second

// Laravel 里硬编码的官方端点（生产值，测试可经 SetEndpointOverrides 覆盖）。
const (
	deepSeekBalanceURL    = "https://api.deepseek.com/user/balance"
	openRouterKeyInfoURL  = "https://openrouter.ai/api/v1/auth/key"
	localBillingCurrency  = "CNY"
	localBillingSource    = "local"
	officialBillingSource = "official"
)

// AIUsageSnapshot 统一的 AI 用量快照（对应 Laravel AiUsageSnapshot）。
type AIUsageSnapshot struct {
	PromptTokens     int
	CompletionTokens int
	Cost             float64
	Currency         string
	Source           string // official | local
}

// TotalTokens 总 token（同 AiUsageSnapshot::totalTokens）。
func (s AIUsageSnapshot) TotalTokens() int { return s.PromptTokens + s.CompletionTokens }

// AIBalance 供应商官方余额快照（对应 Laravel AiBalance）。
type AIBalance struct {
	Currency       string
	TotalBalance   float64
	GrantedBalance *float64
	IsAvailable    bool
}

// AIBillingDriver 计费适配器接口（对应 Laravel AiBillingDriver）。
//
// from / to 为零值时等价于 Laravel 的 null（各驱动按 `$from ?: 默认值` 处理）。
type AIBillingDriver interface {
	// SupportsUsage 是否提供官方 token 用量 API。
	SupportsUsage() bool
	// GetUsage 官方用量查询；返回 nil 表示不支持 / 查询失败（调用方回退本地计费）。
	GetUsage(provider AIProvider, schoolID uint, from, to time.Time) *AIUsageSnapshot
	// GetBalance 官方余额查询；返回 nil 表示不支持。
	GetBalance(provider AIProvider) *AIBalance
}

// AIBilling AI 计费服务（对应 Laravel AiBillingService，driverFor 结果同样按 provider 缓存）。
type AIBilling struct {
	db      *gorm.DB
	ai      *AIService
	drivers map[string]AIBillingDriver

	deepSeekBalanceEndpoint string
	openRouterKeyEndpoint   string
}

// NewAIBilling 创建 AI 计费服务（默认使用 Laravel 的硬编码官方端点）。
func NewAIBilling(db *gorm.DB, ai *AIService) *AIBilling {
	return &AIBilling{
		db:                      db,
		ai:                      ai,
		drivers:                 map[string]AIBillingDriver{},
		deepSeekBalanceEndpoint: deepSeekBalanceURL,
		openRouterKeyEndpoint:   openRouterKeyInfoURL,
	}
}

// SetEndpointOverrides 覆盖 Laravel 里硬编码的两个官方端点（仅测试用：指向 httptest 假上游）。
// 传空串表示保持原值。
func (b *AIBilling) SetEndpointOverrides(deepSeekBalance, openRouterKey string) {
	if deepSeekBalance != "" {
		b.deepSeekBalanceEndpoint = deepSeekBalance
	}
	if openRouterKey != "" {
		b.openRouterKeyEndpoint = openRouterKey
	}
}

// driverFor 返回供应商对应的计费适配器（同 AiBillingService::driverFor）。
func (b *AIBilling) driverFor(providerID string) AIBillingDriver {
	if driver, ok := b.drivers[providerID]; ok {
		return driver
	}

	var driver AIBillingDriver
	switch providerID {
	case "openai":
		driver = &openAiBillingDriver{billing: b}
	case "deepseek":
		driver = &deepSeekBillingDriver{billing: b}
	case "moonshot":
		driver = &moonshotBillingDriver{billing: b}
	case "siliconflow":
		driver = &siliconFlowBillingDriver{billing: b}
	case "openrouter":
		driver = &openRouterBillingDriver{billing: b}
	case "newapi", "oneapi":
		driver = &relayBillingDriver{billing: b}
	default:
		driver = &localPreciseBillingDriver{billing: b}
	}

	b.drivers[providerID] = driver
	return driver
}

// SyncOfficialUsage 同步官方用量：官方 API 可用则返回官方快照，否则回退本地聚合。
//
// 注意：与 Laravel 一致——supportsUsage() 为真但首次查询返回 null 时会再查一次（原实现如此）。
func (b *AIBilling) SyncOfficialUsage(setting *models.AISetting, providerID string, from, to time.Time) *AIUsageSnapshot {
	provider := findProviderByID(setting, providerID)
	if provider == nil {
		return nil
	}

	driver := b.driverFor(providerID)
	if driver.SupportsUsage() {
		if official := driver.GetUsage(provider, setting.SchoolID, from, to); official != nil {
			return official
		}
	}

	return driver.GetUsage(provider, setting.SchoolID, from, to)
}

// GetBalance 官方余额查询（只有支持的驱动才有非 nil 结果）。
func (b *AIBilling) GetBalance(setting *models.AISetting, providerID string) *AIBalance {
	provider := findProviderByID(setting, providerID)
	if provider == nil {
		return nil
	}

	return b.driverFor(providerID).GetBalance(provider)
}

// ============================================================
// 控制器视图（SchoolAdminController::getAiProviderOfficial 的 data）
// ============================================================

// AIProviderOfficialUsage 官方用量视图（键名逐字同 Laravel）。
type AIProviderOfficialUsage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	Cost             float64 `json:"cost"`
	Currency         string  `json:"currency"`
	Source           string  `json:"source"`
}

// AIProviderOfficialBalance 官方余额视图（键名逐字同 Laravel）。
type AIProviderOfficialBalance struct {
	Currency       string   `json:"currency"`
	TotalBalance   float64  `json:"total_balance"`
	GrantedBalance *float64 `json:"granted_balance"`
	IsAvailable    bool     `json:"is_available"`
}

// AIProviderOfficialFallback 官方接口不可用时的 data（只有三个键，无 usage/balance/fetched_at）。
type AIProviderOfficialFallback struct {
	ProviderID string `json:"provider_id"`
	Official   bool   `json:"official"`
	Message    string `json:"message"`
}

// AIProviderOfficialData 官方接口可用时的 data（usage / balance 可为 null）。
type AIProviderOfficialData struct {
	ProviderID string                     `json:"provider_id"`
	Official   bool                       `json:"official"`
	Usage      *AIProviderOfficialUsage   `json:"usage"`
	Balance    *AIProviderOfficialBalance `json:"balance"`
	FetchedAt  string                     `json:"fetched_at"`
}

// AIProviderOfficial 查询某供应商的官方用量与余额（对应 getAiProviderOfficial）。
// 返回 AIProviderOfficialFallback 或 AIProviderOfficialData，由 handler 放进 data。
func (b *AIBilling) AIProviderOfficial(schoolID uint, providerID string) (any, error) {
	if providerID == "" {
		return nil, ErrUnprocessable("缺少供应商")
	}

	var setting models.AISetting
	err := b.db.Where("school_id = ?", schoolID).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUnprocessable("尚未创建 AI 配置")
	}
	if err != nil {
		return nil, err
	}

	provider := findProviderByID(&setting, providerID)
	if provider == nil {
		return nil, ErrUnprocessable("未找到该供应商配置")
	}
	if !aiTruthy(provider["api_key"]) {
		return nil, ErrUnprocessable("该供应商未配置 API Key，请先填写并保存")
	}

	now := util.Now()
	snapshot := b.SyncOfficialUsage(&setting, providerID, now.AddDate(0, 0, -30), now)
	balance := b.GetBalance(&setting, providerID)

	if snapshot == nil && balance == nil {
		return &AIProviderOfficialFallback{
			ProviderID: providerID,
			Official:   false,
			Message:    "该平台未提供可用的官方用量/余额接口，当前数据为本地估算",
		}, nil
	}

	// 官方余额落盘，供列表展示与保存合并（取较大值防清零由 saveAiSettings 负责）。
	if balance != nil {
		providers := decodeAIProviders(setting.Providers)
		for i := range providers {
			if aiProviderString(providers[i], "id") == providerID {
				providers[i]["balance"] = balance.TotalBalance
				break
			}
		}
		if err := b.db.Model(&models.AISetting{}).Where("id = ?", setting.ID).
			Update("providers", encodeAIProviders(providers)).Error; err != nil {
			return nil, err
		}
	}

	data := &AIProviderOfficialData{
		ProviderID: providerID,
		Official:   true,
		FetchedAt:  util.Now().Format("2006-01-02 15:04:05"),
	}
	if snapshot != nil {
		data.Usage = &AIProviderOfficialUsage{
			PromptTokens:     snapshot.PromptTokens,
			CompletionTokens: snapshot.CompletionTokens,
			TotalTokens:      snapshot.TotalTokens(),
			Cost:             snapshot.Cost,
			Currency:         snapshot.Currency,
			Source:           snapshot.Source,
		}
	}
	if balance != nil {
		data.Balance = &AIProviderOfficialBalance{
			Currency:       balance.Currency,
			TotalBalance:   balance.TotalBalance,
			GrantedBalance: balance.GrantedBalance,
			IsAvailable:    balance.IsAvailable,
		}
	}
	return data, nil
}

// ============================================================
// HTTP 辅助（复用 AIService 注入的客户端）
// ============================================================

// getJSON 发 GET 并按 Laravel `$response->failed()` 语义判定成功（status >= 400 或传输异常视为失败）。
func (b *AIBilling) getJSON(rawURL string, query map[string]string, headers map[string]string) (map[string]any, bool) {
	if len(query) > 0 {
		values := url.Values{}
		for key, value := range query {
			values.Set(key, value)
		}
		if encoded := values.Encode(); encoded != "" {
			rawURL += "?" + encoded
		}
	}

	body, status, err := b.ai.doJSON(http.MethodGet, rawURL, headers, nil, aiBillingTimeout)
	if err != nil || status >= 400 {
		return nil, false
	}
	return body, true
}

// bearerHeaders 官方接口通用请求头（Authorization: Bearer + Accept: application/json）。
func bearerHeaders(apiKey string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + apiKey,
		"Accept":        "application/json",
	}
}

// jsonDig 按路径取 JSON 子节点（对应 Laravel data_get 的 `a.b.0.c` 语义）。
func jsonDig(root map[string]any, path ...string) (any, bool) {
	var current any = root
	for _, key := range path {
		switch node := current.(type) {
		case map[string]any:
			value, ok := node[key]
			if !ok {
				return nil, false
			}
			current = value
		case []any:
			index, err := strconv.Atoi(key)
			if err != nil || index < 0 || index >= len(node) {
				return nil, false
			}
			current = node[index]
		default:
			return nil, false
		}
	}
	return current, true
}

// jsonObject 取 JSON 子对象（对应 Laravel `is_array($x)`）。
func jsonObject(root map[string]any, path ...string) (map[string]any, bool) {
	value, ok := jsonDig(root, path...)
	if !ok {
		return nil, false
	}
	object, ok := value.(map[string]any)
	return object, ok
}

// jsonFloat 取数值子节点（数值或数字字符串，等价 PHP `(float)` 转换）。
func jsonFloat(root map[string]any, key string) (float64, bool) {
	value, ok := root[key]
	if !ok || value == nil {
		return 0, false
	}
	switch current := value.(type) {
	case float64:
		return current, true
	case string:
		parsed, err := strconv.ParseFloat(current, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	}
	return 0, false
}

// jsonIsset 等价 PHP `isset($root[$key])`：键不存在或值为 null 都为 false
// （JSON 里的显式 null 与缺失键在 PHP isset / `??` 语义下等价）。
func jsonIsset(root map[string]any, key string) bool {
	value, ok := root[key]
	return ok && value != nil
}

// roundN 等价 PHP round($value, $precision)。
func roundN(value float64, precision int) float64 {
	factor := 1.0
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	scaled := value * factor
	if scaled < 0 {
		return float64(int64(scaled-0.5)) / factor
	}
	return float64(int64(scaled+0.5)) / factor
}

// fromOr 等价 PHP `$from ?: 默认值`。
func fromOr(from time.Time, def time.Time) time.Time {
	if from.IsZero() {
		return def
	}
	return from
}

// ============================================================
// 各驱动
// ============================================================

// openAiBillingDriver OpenAI 官方 Usage API（对应 OpenAiBillingDriver）。
type openAiBillingDriver struct{ billing *AIBilling }

// SupportsUsage OpenAI 提供官方用量接口。
func (d *openAiBillingDriver) SupportsUsage() bool { return true }

// GetUsage GET {api_base}/usage?start_time&end_time&bucket_width=1d。
func (d *openAiBillingDriver) GetUsage(provider AIProvider, _ uint, from, to time.Time) *AIUsageSnapshot {
	apiKey := aiProviderString(provider, "api_key")
	if apiKey == "" {
		return nil
	}

	base := trimTrailingSlash(aiBaseOr(aiProviderString(provider, "api_base"), "https://api.openai.com/v1"))
	from = fromOr(from, util.Now().AddDate(0, 0, -30))
	to = fromOr(to, util.Now())

	body, ok := d.billing.getJSON(base+"/usage", map[string]string{
		"start_time":   strconv.FormatInt(from.Unix(), 10),
		"end_time":     strconv.FormatInt(to.Unix(), 10),
		"bucket_width": "1d",
	}, map[string]string{"Authorization": "Bearer " + apiKey})
	if !ok {
		return nil
	}

	buckets, _ := body["data"].([]any)
	if len(buckets) == 0 {
		return nil
	}

	prompt := 0
	completion := 0
	cost := 0.0
	for _, item := range buckets {
		bucket, isObject := item.(map[string]any)
		if !isObject {
			continue
		}
		prompt += int(aiToFloat(bucket["n_context_tokens_total"]))
		completion += int(aiToFloat(bucket["n_generated_tokens_total"]))
		cost += aiToFloat(bucket["cost_in_usd"])
	}

	return &AIUsageSnapshot{
		PromptTokens:     prompt,
		CompletionTokens: completion,
		Cost:             cost,
		Currency:         "USD",
		Source:           officialBillingSource,
	}
}

// GetBalance OpenAI 无公开余额接口。
func (d *openAiBillingDriver) GetBalance(AIProvider) *AIBalance { return nil }

// deepSeekBillingDriver DeepSeek：无官方 token 用量 API，仅余额查询（对应 DeepSeekBillingDriver）。
type deepSeekBillingDriver struct{ billing *AIBilling }

// SupportsUsage DeepSeek 无官方用量明细接口。
func (d *deepSeekBillingDriver) SupportsUsage() bool { return false }

// GetUsage 走本地精确计费。
func (d *deepSeekBillingDriver) GetUsage(AIProvider, uint, time.Time, time.Time) *AIUsageSnapshot {
	return nil
}

// GetBalance GET <balanceEndpoint>，解析 balance_infos.0 与 is_available。
func (d *deepSeekBillingDriver) GetBalance(provider AIProvider) *AIBalance {
	apiKey := aiProviderString(provider, "api_key")
	if apiKey == "" {
		return nil
	}

	body, ok := d.billing.getJSON(d.billing.deepSeekBalanceEndpoint, nil, bearerHeaders(apiKey))
	if !ok {
		return nil
	}

	info, ok := jsonObject(body, "balance_infos", "0")
	if !ok {
		return nil
	}

	// 同 Laravel：`(string) ($info['currency'] ?? 'CNY')` —— 键缺失/null 才取默认值。
	currency := "CNY"
	if value, found := info["currency"]; found && value != nil {
		if text, isString := value.(string); isString {
			currency = text
		}
	}

	var granted *float64
	if value, found := jsonFloat(info, "granted_balance"); found {
		granted = &value
	}

	// 同 Laravel：`(bool) ($response->json('is_available') ?? true)`（键缺失/null 视为 true）。
	isAvailable := true
	if value, found := body["is_available"]; found && value != nil {
		isAvailable = aiTruthy(value)
	}

	return &AIBalance{
		Currency:       currency,
		TotalBalance:   aiToFloat(info["total_balance"]),
		GrantedBalance: granted,
		IsAvailable:    isAvailable,
	}
}

// moonshotBillingDriver Moonshot（Kimi）：仅余额查询（对应 MoonshotBillingDriver）。
type moonshotBillingDriver struct{ billing *AIBilling }

// SupportsUsage Moonshot 无官方用量 API。
func (d *moonshotBillingDriver) SupportsUsage() bool { return false }

// GetUsage 走本地精确计费。
func (d *moonshotBillingDriver) GetUsage(AIProvider, uint, time.Time, time.Time) *AIUsageSnapshot {
	return nil
}

// GetBalance GET {api_base}/users/me/balance。
func (d *moonshotBillingDriver) GetBalance(provider AIProvider) *AIBalance {
	apiKey := aiProviderString(provider, "api_key")
	if apiKey == "" {
		return nil
	}

	base := trimTrailingSlash(aiBaseOr(aiProviderString(provider, "api_base"), "https://api.moonshot.cn/v1"))
	body, ok := d.billing.getJSON(base+"/users/me/balance", nil, bearerHeaders(apiKey))
	if !ok {
		return nil
	}

	data, ok := jsonObject(body, "data")
	if !ok {
		return nil
	}

	total := 0.0
	if value, found := jsonFloat(data, "total_balance"); found {
		total = value
	} else if value, found := jsonFloat(data, "available_balance"); found {
		total = value
	}
	available := total
	if value, found := jsonFloat(data, "available_balance"); found {
		available = value
	}

	var granted *float64
	if value, found := jsonFloat(data, "voucher_balance"); found {
		granted = &value
	}

	return &AIBalance{
		Currency:       "CNY",
		TotalBalance:   total,
		GrantedBalance: granted,
		IsAvailable:    available > 0,
	}
}

// siliconFlowBillingDriver SiliconFlow（硅基流动）：仅余额查询（对应 SiliconFlowBillingDriver）。
type siliconFlowBillingDriver struct{ billing *AIBilling }

// SupportsUsage SiliconFlow 无官方用量 API。
func (d *siliconFlowBillingDriver) SupportsUsage() bool { return false }

// GetUsage 走本地精确计费。
func (d *siliconFlowBillingDriver) GetUsage(AIProvider, uint, time.Time, time.Time) *AIUsageSnapshot {
	return nil
}

// GetBalance GET {api_base}/user/info → data.balance / data.totalBalance。
func (d *siliconFlowBillingDriver) GetBalance(provider AIProvider) *AIBalance {
	apiKey := aiProviderString(provider, "api_key")
	if apiKey == "" {
		return nil
	}

	base := trimTrailingSlash(aiBaseOr(aiProviderString(provider, "api_base"), "https://api.siliconflow.cn/v1"))
	body, ok := d.billing.getJSON(base+"/user/info", nil, bearerHeaders(apiKey))
	if !ok {
		return nil
	}

	data, ok := jsonObject(body, "data")
	if !ok {
		return nil
	}
	hasBalance := jsonIsset(data, "balance")
	hasTotalBalance := jsonIsset(data, "totalBalance")
	if !hasBalance && !hasTotalBalance {
		// 同 Laravel：`!isset($data['balance']) && !isset($data['totalBalance'])` → null
		return nil
	}

	total, found := jsonFloat(data, "totalBalance")
	if !found {
		total, _ = jsonFloat(data, "balance")
	}
	balanceValue, _ := jsonFloat(data, "balance")

	var granted *float64
	if hasBalance {
		value := balanceValue
		granted = &value
	}

	return &AIBalance{
		Currency:       "CNY",
		TotalBalance:   total,
		GrantedBalance: granted,
		IsAvailable:    balanceValue > 0,
	}
}

// openRouterBillingDriver OpenRouter：credits 已用金额与额度（对应 OpenRouterBillingDriver）。
type openRouterBillingDriver struct{ billing *AIBilling }

// SupportsUsage OpenRouter 提供官方 usage。
func (d *openRouterBillingDriver) SupportsUsage() bool { return true }

// GetUsage 已用金额（USD），无 token 拆分。
func (d *openRouterBillingDriver) GetUsage(provider AIProvider, _ uint, _, _ time.Time) *AIUsageSnapshot {
	info := d.keyInfo(provider)
	if info == nil {
		return nil
	}
	usage, ok := jsonFloat(info, "usage")
	if !ok {
		return nil
	}

	return &AIUsageSnapshot{
		PromptTokens:     0,
		CompletionTokens: 0,
		Cost:             usage,
		Currency:         "USD",
		Source:           officialBillingSource,
	}
}

// GetBalance limit - usage（limit 为 null / 非数字时视为不限额度的 Key，无余额概念）。
func (d *openRouterBillingDriver) GetBalance(provider AIProvider) *AIBalance {
	info := d.keyInfo(provider)
	if info == nil {
		return nil
	}
	limit, ok := jsonFloat(info, "limit")
	if !ok {
		return nil
	}

	used, _ := jsonFloat(info, "usage")
	remaining := roundN(limit-used, 4)

	return &AIBalance{
		Currency:     "USD",
		TotalBalance: remaining,
		IsAvailable:  limit-used > 0,
	}
}

// keyInfo GET <keyEndpoint> → data。
func (d *openRouterBillingDriver) keyInfo(provider AIProvider) map[string]any {
	apiKey := aiProviderString(provider, "api_key")
	if apiKey == "" {
		return nil
	}

	body, ok := d.billing.getJSON(d.billing.openRouterKeyEndpoint, nil, bearerHeaders(apiKey))
	if !ok {
		return nil
	}

	info, _ := jsonObject(body, "data")
	return info
}

// relayBillingDriver 中转平台（New API / One API）：OpenAI billing 兼容端点（对应 RelayBillingDriver）。
type relayBillingDriver struct{ billing *AIBilling }

// SupportsUsage 中转平台提供兼容用量端点。
func (d *relayBillingDriver) SupportsUsage() bool { return true }

// GetUsage 区间已用金额（USD）。
func (d *relayBillingDriver) GetUsage(provider AIProvider, _ uint, from, to time.Time) *AIUsageSnapshot {
	if aiProviderString(provider, "api_key") == "" {
		return nil
	}

	used := d.fetchUsedUSD(provider, from, to)
	if used == nil {
		return nil
	}

	return &AIUsageSnapshot{
		PromptTokens:     0,
		CompletionTokens: 0,
		Cost:             *used,
		Currency:         "USD",
		Source:           officialBillingSource,
	}
}

// GetBalance hard_limit_usd / system_hard_limit_usd 作为剩余额度（USD）。
func (d *relayBillingDriver) GetBalance(provider AIProvider) *AIBalance {
	if aiProviderString(provider, "api_key") == "" {
		return nil
	}

	json, ok := d.billingGet(provider, "/subscription", nil)
	if !ok {
		return nil
	}

	limit, found := jsonFloat(json, "hard_limit_usd")
	if !found {
		limit, found = jsonFloat(json, "system_hard_limit_usd")
	}
	if !found {
		return nil
	}

	return &AIBalance{
		Currency:     "USD",
		TotalBalance: roundN(limit, 4),
		IsAvailable:  limit > 0,
	}
}

// fetchUsedUSD 查询区间已用金额（USD）；total_usage 单位为美分。
func (d *relayBillingDriver) fetchUsedUSD(provider AIProvider, from, to time.Time) *float64 {
	start := fromOr(from, util.Now().AddDate(0, 0, -30)).Format("2006-01-02")
	end := fromOr(to, util.Now().AddDate(0, 0, 1)).Format("2006-01-02")

	json, ok := d.billingGet(provider, "/usage", map[string]string{"start_date": start, "end_date": end})
	if !ok {
		return nil
	}

	total, found := jsonFloat(json, "total_usage")
	if !found {
		return nil
	}

	used := roundN(total/100, 6)
	return &used
}

// billingGet 计费端点 GET，自动尝试 /v1 前缀与无前缀两种部署布局（同 RelayBillingDriver::billingGet）。
func (d *relayBillingDriver) billingGet(provider AIProvider, path string, query map[string]string) (map[string]any, bool) {
	apiKey := aiProviderString(provider, "api_key")
	base := trimTrailingSlash(aiProviderString(provider, "api_base"))
	if apiKey == "" || base == "" {
		return nil, false
	}

	origin := trimTrailingSlash(trimSuffix(base, "/v1"))
	candidates := []string{
		origin + "/v1/dashboard/billing" + path,
		origin + "/dashboard/billing" + path,
	}

	for _, candidate := range candidates {
		body, ok := d.billing.getJSON(candidate, query, bearerHeaders(apiKey))
		if !ok {
			continue
		}
		// 同 Laravel：首个成功响应即返回（不再回退第二个地址）。
		return body, true
	}
	return nil, false
}

// localPreciseBillingDriver 通用回退：本地精确计费（对应 LocalPreciseBillingDriver）。
type localPreciseBillingDriver struct{ billing *AIBilling }

// SupportsUsage 本地聚合始终可用。
func (d *localPreciseBillingDriver) SupportsUsage() bool { return true }

// GetUsage 从 ai_conversations 聚合本供应商的 token 与 cost。
func (d *localPreciseBillingDriver) GetUsage(provider AIProvider, schoolID uint, from, to time.Time) *AIUsageSnapshot {
	if schoolID == 0 {
		return nil
	}

	query := d.billing.db.Model(&models.AIConversation{}).
		Where("school_id = ? AND provider = ? AND status = ?", schoolID, aiProviderString(provider, "id"), "completed")
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	var row struct {
		Prompt     int64   `gorm:"column:prompt"`
		Completion int64   `gorm:"column:completion"`
		Cost       float64 `gorm:"column:cost"`
	}
	err := query.Select("COALESCE(SUM(prompt_tokens),0) AS prompt, " +
		"COALESCE(SUM(completion_tokens),0) AS completion, COALESCE(SUM(cost),0) AS cost").
		Scan(&row).Error
	if err != nil {
		return nil
	}

	// 同 Laravel：`(string) ($provider['currency'] ?? 'CNY')` —— 键缺失/null 才取默认值。
	currency := localBillingCurrency
	if value, found := provider["currency"]; found && value != nil {
		if text, isString := value.(string); isString {
			currency = text
		}
	}

	return &AIUsageSnapshot{
		PromptTokens:     int(row.Prompt),
		CompletionTokens: int(row.Completion),
		Cost:             row.Cost,
		Currency:         currency,
		Source:           localBillingSource,
	}
}

// GetBalance 本地计费无余额概念。
func (d *localPreciseBillingDriver) GetBalance(AIProvider) *AIBalance { return nil }

// trimTrailingSlash 等价 PHP rtrim($value, '/').
func trimTrailingSlash(value string) string {
	for len(value) > 0 && value[len(value)-1] == '/' {
		value = value[:len(value)-1]
	}
	return value
}

// trimSuffix 去掉后缀（等价 PHP preg_replace('/\/v1$/', ”, $base)）。
func trimSuffix(value, suffix string) string {
	if len(value) >= len(suffix) && value[len(value)-len(suffix):] == suffix {
		return value[:len(value)-len(suffix)]
	}
	return value
}
