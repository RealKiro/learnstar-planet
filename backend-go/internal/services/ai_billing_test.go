// AI 供应商官方账单查询测试：POST /api/v1/admin/ai/provider-official。
//
// 覆盖：DeepSeek 余额、OpenAI Usage、Moonshot / SiliconFlow / OpenRouter / 中转（New API）/
// 本地精确计费（claude 回退）的解析结果与币种、账本落盘；上游 500 / 无法解析 / 传输异常 → 兜底形状；
// 未配置 API Key / 未知 provider / 未创建配置 / 缺 provider_id 的 422 文案。
//
// 所有出网一律走 httptest 假上游（DeepSeek 与 OpenRouter 的硬编码端点用 SetEndpointOverrides 指向假上游），
// **不访问真实外网**。
package services_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const billingFallbackMessage = "该平台未提供可用的官方用量/余额接口，当前数据为本地估算"

// newBillingFixture 建学校 + 一条 AI 设置，并返回注入假上游的计费服务。
func newBillingFixture(t *testing.T, db *gorm.DB, fake *aiFakeProvider, providers string) (models.School, *services.AIBilling) {
	t.Helper()
	school := seedSchool(t, db)
	aiSettingWith(t, db, school.ID, true, providers)
	return school, services.NewAIBilling(db, aiServiceFor(fake))
}

// officialDataOf 断言返回的是 official=true 的 data。
func officialDataOf(t *testing.T, data any) *services.AIProviderOfficialData {
	t.Helper()
	out, ok := data.(*services.AIProviderOfficialData)
	require.True(t, ok, "应为 AIProviderOfficialData，实际 %T", data)
	return out
}

// fallbackOf 断言返回的是 official=false 的兜底 data，并核对 JSON 形状（只有三个键）。
func fallbackOf(t *testing.T, data any, providerID string) {
	t.Helper()
	_, ok := data.(*services.AIProviderOfficialFallback)
	require.True(t, ok, "应为 AIProviderOfficialFallback，实际 %T", data)
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	assert.JSONEq(t, `{"provider_id":"`+providerID+`","official":false,"message":"`+billingFallbackMessage+`"}`, string(raw))
}

// queryOf 解析一次假上游请求的查询串。
func queryOf(t *testing.T, raw string) url.Values {
	t.Helper()
	values, err := url.ParseQuery(raw)
	require.NoError(t, err)
	return values
}

// aiRequestForPath 取假上游中指定路径的最后一次请求（按路径精确匹配）。
func aiRequestForPath(t *testing.T, fake *aiFakeProvider, path string) aiRequest {
	t.Helper()
	fake.mu.Lock()
	defer fake.mu.Unlock()
	for i := len(fake.reqs) - 1; i >= 0; i-- {
		if fake.reqs[i].Path == path {
			return fake.reqs[i]
		}
	}
	t.Fatalf("假上游未收到 %s 请求", path)
	return aiRequest{}
}

// DeepSeek：余额端点解析（币种 / 总额 / 赠送额 / is_available）与账本落盘；usage 恒为 null。
func TestAIBillingDeepSeekBalance(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK,
		`{"is_available":true,"balance_infos":[{"currency":"CNY","total_balance":"110.00","granted_balance":"10.00"}]}`))
	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("deepseek", "", "deepseek-chat", "sk-deepseek")))
	billing.SetEndpointOverrides(fake.server.URL+"/user/balance", fake.server.URL+"/auth/key")

	data, err := billing.AIProviderOfficial(school.ID, "deepseek")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	assert.Equal(t, "deepseek", out.ProviderID)
	assert.True(t, out.Official)
	assert.Nil(t, out.Usage, "DeepSeek 无官方用量接口 → usage 为 null")
	require.NotNil(t, out.Balance)
	assert.Equal(t, "CNY", out.Balance.Currency)
	assert.Equal(t, 110.0, out.Balance.TotalBalance)
	require.NotNil(t, out.Balance.GrantedBalance)
	assert.Equal(t, 10.0, *out.Balance.GrantedBalance)
	assert.True(t, out.Balance.IsAvailable)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, out.FetchedAt)

	// JSON 键形状（usage 保留为 null，供前端区分「无官方用量」与「未查询」）。
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	var shape map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &shape))
	assert.Equal(t, []string{"balance", "fetched_at", "official", "provider_id", "usage"},
		sortedKeysOf(shape))
	assert.Equal(t, "null", string(shape["usage"]))

	// 请求头与地址（Authorization: Bearer + Accept: application/json）。
	last := fake.aiLast(t)
	assert.Equal(t, http.MethodGet, last.Method)
	assert.Equal(t, "/user/balance", last.Path)
	assert.Equal(t, "Bearer sk-deepseek", last.Header.Get("Authorization"))
	assert.Equal(t, "application/json", last.Header.Get("Accept"))

	// 官方余额写回 providers[].balance。
	providers := aiReadProviders(t, aiReloadSetting(t, db, school.ID).Providers)
	require.Len(t, providers, 1)
	assert.Equal(t, 110.0, providers[0]["balance"])
}

// DeepSeek：上游 500 / 无法解析 / 传输异常 → official=false 兜底形状，且不落盘 balance。
func TestAIBillingDeepSeekFallbacks(t *testing.T) {
	cases := []struct {
		name     string
		respond  func(w http.ResponseWriter, r *http.Request)
		erroring bool
	}{
		{name: "upstream 500", respond: aiRespondJSON(http.StatusInternalServerError, `{"error":"boom"}`)},
		{name: "empty balance_infos", respond: aiRespondJSON(http.StatusOK, `{"is_available":true,"balance_infos":[]}`)},
		{name: "unparsable body", respond: aiRespondJSON(http.StatusOK, `{"foo":1}`)},
		{name: "transport error", erroring: true},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			db := setupDB(t)

			var fake *aiFakeProvider
			var ai *services.AIService
			if item.erroring {
				ai = services.NewAIService()
				ai.SetHTTPClient(aiErrorDoer{})
			} else {
				fake = aiNewFakeProvider(t, item.respond)
				ai = aiServiceFor(fake)
			}

			school := seedSchool(t, db)
			aiSettingWith(t, db, school.ID, true,
				aiProvidersJSON(t, aiProviderEntry("deepseek", "", "deepseek-chat", "sk-deepseek")))
			billing := services.NewAIBilling(db, ai)
			billing.SetEndpointOverrides(fake2URL(fake)+"/user/balance", "")

			data, err := billing.AIProviderOfficial(school.ID, "deepseek")
			require.NoError(t, err)
			fallbackOf(t, data, "deepseek")

			// 未拿到余额 → 不写 providers[].balance。
			providers := aiReadProviders(t, aiReloadSetting(t, db, school.ID).Providers)
			require.Len(t, providers, 1)
			assert.NotContains(t, providers[0], "balance")
		})
	}
}

// OpenAI：官方 Usage API 汇总 token 与美元费用；无余额接口 → balance 为 null。
func TestAIBillingOpenAIUsage(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK,
		`{"data":[{"n_context_tokens_total":100,"n_generated_tokens_total":50,"cost_in_usd":0.25},`+
			`{"n_context_tokens_total":1,"n_generated_tokens_total":2,"cost_in_usd":0.01}]}`))
	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL+"/v1", "gpt-4o", "sk-openai")))

	data, err := billing.AIProviderOfficial(school.ID, "openai")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	require.NotNil(t, out.Usage)
	assert.Equal(t, 101, out.Usage.PromptTokens)
	assert.Equal(t, 52, out.Usage.CompletionTokens)
	assert.Equal(t, 153, out.Usage.TotalTokens)
	assert.InDelta(t, 0.26, out.Usage.Cost, 1e-9)
	assert.Equal(t, "USD", out.Usage.Currency)
	assert.Equal(t, "official", out.Usage.Source)
	assert.Nil(t, out.Balance, "OpenAI 无公开余额接口 → balance 为 null")

	// 首次查询即成功 → 只发一次请求；查询串与 Laravel 一致（bucket_width=1d + 秒级时间戳）。
	assert.Equal(t, 1, fake.aiReqs())
	last := fake.aiLast(t)
	assert.Equal(t, "/v1/usage", last.Path)
	query := queryOf(t, last.Query)
	assert.Equal(t, "1d", query.Get("bucket_width"))
	_, err = strconv.ParseInt(query.Get("start_time"), 10, 64)
	assert.NoError(t, err, "start_time 应为秒级时间戳")
	_, err = strconv.ParseInt(query.Get("end_time"), 10, 64)
	assert.NoError(t, err, "end_time 应为秒级时间戳")
	assert.Equal(t, "Bearer sk-openai", last.Header.Get("Authorization"))

	// 无 balance → 不写 providers[].balance。
	providers := aiReadProviders(t, aiReloadSetting(t, db, school.ID).Providers)
	assert.NotContains(t, providers[0], "balance")
}

// OpenAI：上游 401 → usage/balance 都为 null → 兜底；与 Laravel 一致会重试一次（共 2 次请求）。
func TestAIBillingOpenAIUnauthorizedFallback(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusUnauthorized, `{"error":"invalid api key"}`))

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL+"/v1", "gpt-4o", "sk-openai")))

	data, err := billing.AIProviderOfficial(school.ID, "openai")
	require.NoError(t, err)
	fallbackOf(t, data, "openai")
	assert.Equal(t, 2, fake.aiReqs(), "supportsUsage 为真但首次返回 null 时会再查一次（同 Laravel）")
}

// Moonshot：仅余额，CNY；available_balance 为 0 → is_available=false，无 voucher 时 granted 为 null。
func TestAIBillingMoonshotBalance(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK,
		`{"data":{"total_balance":88.5,"available_balance":80,"voucher_balance":8.5}}`))

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("moonshot", fake.server.URL+"/v1", "moonshot-v1-8k", "sk-moonshot")))

	data, err := billing.AIProviderOfficial(school.ID, "moonshot")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	assert.Nil(t, out.Usage)
	require.NotNil(t, out.Balance)
	assert.Equal(t, "CNY", out.Balance.Currency)
	assert.Equal(t, 88.5, out.Balance.TotalBalance)
	require.NotNil(t, out.Balance.GrantedBalance)
	assert.Equal(t, 8.5, *out.Balance.GrantedBalance)
	assert.True(t, out.Balance.IsAvailable)
	assert.Equal(t, "/v1/users/me/balance", fake.aiLast(t).Path)

	// 余额为 0（available_balance=0，无 voucher）→ is_available=false，granted_balance=null。
	zero := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{"data":{"total_balance":2,"available_balance":0}}`))
	school2 := seedExtraSchool(t, db, "billing-moonshot-zero")
	aiSettingWith(t, db, school2.ID, true,
		aiProvidersJSON(t, aiProviderEntry("moonshot", zero.server.URL+"/v1", "moonshot-v1-8k", "sk-m2")))
	zeroBilling := services.NewAIBilling(db, aiServiceFor(zero))

	data2, err := zeroBilling.AIProviderOfficial(school2.ID, "moonshot")
	require.NoError(t, err)
	out2 := officialDataOf(t, data2)
	require.NotNil(t, out2.Balance)
	assert.Equal(t, 2.0, out2.Balance.TotalBalance, "total_balance 优先于 available_balance")
	assert.Nil(t, out2.Balance.GrantedBalance)
	assert.False(t, out2.Balance.IsAvailable)

	raw, err := json.Marshal(data2)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"granted_balance":null`)
}

// SiliconFlow：data.balance / data.totalBalance，CNY。
func TestAIBillingSiliconFlowBalance(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{"data":{"balance":"1.5","totalBalance":"9"}}`))

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("siliconflow", fake.server.URL+"/v1", "Qwen/Qwen2.5-7B", "sk-sf")))

	data, err := billing.AIProviderOfficial(school.ID, "siliconflow")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	assert.Nil(t, out.Usage)
	require.NotNil(t, out.Balance)
	assert.Equal(t, "CNY", out.Balance.Currency)
	assert.Equal(t, 9.0, out.Balance.TotalBalance, "totalBalance 优先")
	require.NotNil(t, out.Balance.GrantedBalance)
	assert.Equal(t, 1.5, *out.Balance.GrantedBalance)
	assert.True(t, out.Balance.IsAvailable)
	assert.Equal(t, "/v1/user/info", fake.aiLast(t).Path)

	// 两个键都没有 → 视为失败（is_array + isset 判断）→ 兜底。
	empty := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{"data":{"foo":1}}`))
	school2 := seedExtraSchool(t, db, "billing-sf-empty")
	aiSettingWith(t, db, school2.ID, true,
		aiProvidersJSON(t, aiProviderEntry("siliconflow", empty.server.URL+"/v1", "m", "sk-sf2")))
	data2, err := services.NewAIBilling(db, aiServiceFor(empty)).AIProviderOfficial(school2.ID, "siliconflow")
	require.NoError(t, err)
	fallbackOf(t, data2, "siliconflow")

	// 键存在但值为显式 null：PHP `isset` 视为缺失 → 同样兜底
	// （不能用「键是否存在」判断，否则会给出 official=true + 余额 0 的假官方数据）。
	nulls := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{"data":{"balance":null,"totalBalance":null}}`))
	school3 := seedExtraSchool(t, db, "billing-sf-nulls")
	aiSettingWith(t, db, school3.ID, true,
		aiProvidersJSON(t, aiProviderEntry("siliconflow", nulls.server.URL+"/v1", "m", "sk-sf3")))
	data3, err := services.NewAIBilling(db, aiServiceFor(nulls)).AIProviderOfficial(school3.ID, "siliconflow")
	require.NoError(t, err)
	fallbackOf(t, data3, "siliconflow")
}

// OpenRouter：usage 为已用美元；limit 存在时余额 = limit - usage，limit 为 null → 无余额。
func TestAIBillingOpenRouter(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{"data":{"usage":12.5,"limit":25}}`))

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("openrouter", "", "openai/gpt-4o", "sk-or")))
	billing.SetEndpointOverrides("", fake.server.URL+"/auth/key")

	data, err := billing.AIProviderOfficial(school.ID, "openrouter")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	require.NotNil(t, out.Usage)
	assert.Equal(t, 0, out.Usage.PromptTokens)
	assert.Equal(t, 0, out.Usage.CompletionTokens)
	assert.Equal(t, 12.5, out.Usage.Cost)
	assert.Equal(t, "USD", out.Usage.Currency)
	assert.Equal(t, "official", out.Usage.Source)
	require.NotNil(t, out.Balance)
	assert.Equal(t, 12.5, out.Balance.TotalBalance, "limit 25 - usage 12.5")
	assert.True(t, out.Balance.IsAvailable)

	// limit 为 null（不限额度的 Key）→ 只有 usage，无余额。
	unlimited := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{"data":{"usage":3,"limit":null}}`))
	school2 := seedExtraSchool(t, db, "billing-or-unlimited")
	aiSettingWith(t, db, school2.ID, true, aiProvidersJSON(t, aiProviderEntry("openrouter", "", "m", "sk-or2")))
	billing2 := services.NewAIBilling(db, aiServiceFor(unlimited))
	billing2.SetEndpointOverrides("", unlimited.server.URL+"/auth/key")

	data2, err := billing2.AIProviderOfficial(school2.ID, "openrouter")
	require.NoError(t, err)
	out2 := officialDataOf(t, data2)
	require.NotNil(t, out2.Usage)
	assert.Nil(t, out2.Balance)
}

// 中转（New API / One API）：billing 兼容端点，total_usage 以美分计。
func TestAIBillingRelay(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/dashboard/billing/usage":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"total_usage":1234}`))
		case "/v1/dashboard/billing/subscription":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"hard_limit_usd":10}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("newapi", fake.server.URL+"/v1", "gpt-4o", "sk-relay")))

	data, err := billing.AIProviderOfficial(school.ID, "newapi")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	require.NotNil(t, out.Usage)
	assert.InDelta(t, 12.34, out.Usage.Cost, 1e-9, "total_usage 1234 美分 = 12.34 USD")
	assert.Equal(t, "USD", out.Usage.Currency)
	assert.Equal(t, "official", out.Usage.Source)
	require.NotNil(t, out.Balance)
	assert.Equal(t, 10.0, out.Balance.TotalBalance)
	assert.True(t, out.Balance.IsAvailable)

	// usage 走 /v1/dashboard/billing/usage 且带 start_date / end_date（Y-m-d）；余额走 /subscription。
	usageReq := aiRequestForPath(t, fake, "/v1/dashboard/billing/usage")
	query := queryOf(t, usageReq.Query)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, query.Get("start_date"))
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, query.Get("end_date"))
	assert.Equal(t, "Bearer sk-relay", usageReq.Header.Get("Authorization"))
	assert.Equal(t, "/v1/dashboard/billing/subscription", fake.aiLast(t).Path)

	// oneapi 走同一驱动；无 api_base 时无法构造端点 → 兜底（api_key 已配置）。
	school2 := seedExtraSchool(t, db, "billing-oneapi")
	aiSettingWith(t, db, school2.ID, true, aiProvidersJSON(t, aiProviderEntry("oneapi", "", "gpt-4o", "sk-one")))
	data2, err := services.NewAIBilling(db, aiServiceFor(fake)).AIProviderOfficial(school2.ID, "oneapi")
	require.NoError(t, err)
	fallbackOf(t, data2, "oneapi")
}

// 未列入驱动表的供应商（claude 等）→ 本地精确计费：聚合 ai_conversations，source=local，币种取配置。
func TestAIBillingLocalPreciseFallback(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{}`))

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("claude", "", "claude-3-sonnet", "sk-claude")))

	now := util.Now()
	seedAIConversation(t, db, school.ID, now.Add(-2*time.Hour), 300, 0.5, "claude", "小明", "问题1", "答案1")
	seedAIConversation(t, db, school.ID, now.Add(-1*time.Hour), 600, 1.25, "claude", "小红", "问题2", "答案2")
	// 其他供应商 / 非 completed 的流水不计入。
	seedAIConversation(t, db, school.ID, now.Add(-1*time.Hour), 999, 9.9, "openai", "小明", "问题3", "答案3")

	data, err := billing.AIProviderOfficial(school.ID, "claude")
	require.NoError(t, err)

	out := officialDataOf(t, data)
	require.NotNil(t, out.Usage)
	assert.Equal(t, 100+200, out.Usage.PromptTokens, "prompt = 300/3 + 600/3")
	assert.Equal(t, 200+400, out.Usage.CompletionTokens)
	assert.Equal(t, 900, out.Usage.TotalTokens)
	assert.InDelta(t, 1.75, out.Usage.Cost, 1e-9)
	assert.Equal(t, "local", out.Usage.Source)
	assert.Equal(t, "USD", out.Usage.Currency, "币种取供应商配置 currency")
	assert.Nil(t, out.Balance)
	assert.Equal(t, 0, fake.aiReqs(), "本地精确计费不出网")
}

// 守卫：缺 provider_id / 未创建配置 / 未知供应商 / 未配置 API Key 的 422 文案。
func TestAIBillingGuards(t *testing.T) {
	db := setupDB(t)
	fake := aiNewFakeProvider(t, aiRespondJSON(http.StatusOK, `{}`))

	school, billing := newBillingFixture(t, db, fake,
		aiProvidersJSON(t, aiProviderEntry("deepseek", "", "deepseek-chat", ""), // 无 api_key
			aiProviderEntry("openai", "", "gpt-4o", "sk-openai")))

	_, err := billing.AIProviderOfficial(school.ID, "")
	assert.Equal(t, "缺少供应商", appErrorOf(t, err, http.StatusUnprocessableEntity))

	_, err = billing.AIProviderOfficial(school.ID, "moonshot")
	assert.Equal(t, "未找到该供应商配置", appErrorOf(t, err, http.StatusUnprocessableEntity))

	_, err = billing.AIProviderOfficial(school.ID, "deepseek")
	assert.Equal(t, "该供应商未配置 API Key，请先填写并保存", appErrorOf(t, err, http.StatusUnprocessableEntity))

	// 本校尚未创建 AI 配置 → 422（注意 Laravel 不做 firstOrCreate）。
	noSetting := seedExtraSchool(t, db, "billing-no-setting")
	_, err = billing.AIProviderOfficial(noSetting.ID, "openai")
	assert.Equal(t, "尚未创建 AI 配置", appErrorOf(t, err, http.StatusUnprocessableEntity))
	assert.Equal(t, 0, fake.aiReqs(), "守卫分支不出网")
}

// fake2URL 假上游地址（erroring 场景没有假服务器时返回空串）。
func fake2URL(fake *aiFakeProvider) string {
	if fake == nil {
		return ""
	}
	return fake.server.URL
}

// sortedKeysOf 排序后的 JSON 键列表。
func sortedKeysOf(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sortStrings(keys)
	return keys
}

// sortStrings 极简插入排序（避免为一个断言引入更多依赖）。
func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
