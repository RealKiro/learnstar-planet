// AI 配置与对话链路测试：管理端设置/开关/用量/模型列表/连通性测试，
// 教师端配置/对话/命令/用量，大屏端开关检查/对话。
//
// 所有上游调用都指向 httptest 假供应商（api_base 覆盖生效），**不访问真实外网**。
package services_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// 测试夹具
// ============================================================

// aiRequest 记录假供应商收到的一次请求。
type aiRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   map[string]any
}

// aiFakeProvider 假供应商服务器（可记录请求、按需应答）。
type aiFakeProvider struct {
	server *httptest.Server

	mu   sync.Mutex
	reqs []aiRequest
}

// aiNewFakeProvider 启动假供应商；respond 决定响应内容。
func aiNewFakeProvider(t *testing.T, respond func(w http.ResponseWriter, r *http.Request)) *aiFakeProvider {
	t.Helper()
	fake := &aiFakeProvider{}
	fake.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &body)
		}
		fake.mu.Lock()
		fake.reqs = append(fake.reqs, aiRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Header: r.Header.Clone(),
			Body:   body,
		})
		fake.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		respond(w, r)
	}))
	t.Cleanup(fake.server.Close)
	return fake
}

// aiRespondJSON 固定状态码 + 响应体。
func aiRespondJSON(status int, payload string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = io.WriteString(w, payload)
	}
}

// aiLast 取最后一次请求。
func (f *aiFakeProvider) aiLast(t *testing.T) aiRequest {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	require.NotEmpty(t, f.reqs, "假供应商未收到请求")
	return f.reqs[len(f.reqs)-1]
}

// aiReqs 请求次数。
func (f *aiFakeProvider) aiReqs() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.reqs)
}

// aiClient 返回该假服务器的 HTTP 客户端，用于注入 AIService。
func (f *aiFakeProvider) aiClient() *http.Client { return f.server.Client() }

// aiServiceFor 构建注入了假供应商客户端的 AIService。
func aiServiceFor(fake *aiFakeProvider) *services.AIService {
	ai := services.NewAIService()
	ai.SetHTTPClient(fake.aiClient())
	return ai
}

// aiErrorDoer 永远失败的 HTTP 客户端（模拟 DNS/连接失败的传输层异常）。
type aiErrorDoer struct{}

func (aiErrorDoer) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("dial tcp: connect: connection refused")
}

// aiSettingWith 写入一条 AI 设置。
func aiSettingWith(t *testing.T, db *gorm.DB, schoolID uint, enabled bool, providers string) models.AISetting {
	t.Helper()
	setting := models.AISetting{
		SchoolID:    schoolID,
		Enabled:     enabled,
		Provider:    "openai",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   2000,
		TokensLimit: 1000000,
		Providers:   providers,
	}
	require.NoError(t, db.Create(&setting).Error)
	return setting
}

// aiProviderEntry 组装一条供应商配置。
func aiProviderEntry(id, apiBase, model, key string) map[string]any {
	return map[string]any{
		"id":                 id,
		"label":              strings.ToUpper(id),
		"api_key":            key,
		"api_base":           apiBase,
		"model":              model,
		"models":             []any{model},
		"is_active":          true,
		"input_price_per_m":  1.0,
		"output_price_per_m": 2.0,
		"currency":           "USD",
	}
}

// aiProvidersJSON 把供应商列表序列化成 providers 列文本。
func aiProvidersJSON(t *testing.T, entries ...map[string]any) string {
	t.Helper()
	if len(entries) == 0 {
		return ""
	}
	raw, err := json.Marshal(entries)
	require.NoError(t, err)
	return string(raw)
}

// aiDecodeRaw 把请求体字面量解码成 RawMessage 字段表（模拟 handler 的解析）。
func aiDecodeRaw(t *testing.T, body string) map[string]json.RawMessage {
	t.Helper()
	raw := map[string]json.RawMessage{}
	require.NoError(t, json.Unmarshal([]byte(body), &raw))
	return raw
}

// aiReloadSetting 重新读取设置行。
func aiReloadSetting(t *testing.T, db *gorm.DB, schoolID uint) models.AISetting {
	t.Helper()
	var setting models.AISetting
	require.NoError(t, db.Where("school_id = ?", schoolID).First(&setting).Error)
	return setting
}

// aiSetTeacherActiveClass 写入教师的 active_class_id 设置并返回刷新后的用户。
func aiSetTeacherActiveClass(t *testing.T, db *gorm.DB, teacher models.User, classID uint) models.User {
	t.Helper()
	raw, err := teacher.WithSetting("active_class_id", classID)
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.User{}).Where("id = ?", teacher.ID).Update("settings", raw).Error)

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	return reloaded
}

// aiReadProviders 解析 providers 列。
func aiReadProviders(t *testing.T, raw string) []map[string]any {
	t.Helper()
	var list []map[string]any
	require.NoError(t, json.Unmarshal([]byte(raw), &list))
	return list
}

// ============================================================
// 1. 设置读写 / 开关
// ============================================================

func TestAIAdminSettingsDefaultsAndRoundTrip(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAIAdmin(db, services.NewAIService())

	view, err := admin.Settings(school.ID)
	require.NoError(t, err)
	assert.False(t, view.Enabled)
	assert.Equal(t, 0, view.TokensUsed)
	assert.Equal(t, 1000000, view.TokensLimit)
	assert.Equal(t, 2000, view.MaxTokens)
	assert.Empty(t, view.Providers)

	// 默认值确实落库（provider=openai / model=gpt-3.5-turbo / enabled=false）
	row := aiReloadSetting(t, db, school.ID)
	assert.False(t, row.Enabled)
	assert.Equal(t, "openai", row.Provider)
	assert.Equal(t, "gpt-3.5-turbo", row.Model)
	assert.Equal(t, 1000000, row.TokensLimit)
	assert.Equal(t, 2000, row.MaxTokens)

	// 保存 → 再读一致（api_key 不做掩码/截断）
	fe, err := admin.SaveSettings(school.ID, aiDecodeRaw(t, `{
		"enabled": true,
		"max_tokens": 4096,
		"tokens_limit": 5000,
		"providers": [{
			"id": "openai", "label": "OpenAI", "api_key": "sk-abcdefghijklmnop",
			"api_base": "http://127.0.0.1:1/v1", "model": "gpt-4o",
			"models": ["gpt-4o"], "is_active": true,
			"input_price_per_m": 1.5, "output_price_per_m": 2.5, "currency": "USD"
		}]
	}`))
	require.NoError(t, err)
	assert.Empty(t, fe)

	view, err = admin.Settings(school.ID)
	require.NoError(t, err)
	assert.True(t, view.Enabled)
	assert.Equal(t, 4096, view.MaxTokens)
	assert.Equal(t, 5000, view.TokensLimit)
	require.Len(t, view.Providers, 1)
	assert.Equal(t, "sk-abcdefghijklmnop", view.Providers[0]["api_key"], "api_key 不得掩码")
	assert.Equal(t, "gpt-4o", view.Providers[0]["model"])
	assert.Equal(t, "USD", view.Providers[0]["currency"])

	// enabled=false 也必须正确落库（默认 false 的布尔列）
	fe, err = admin.SaveSettings(school.ID, aiDecodeRaw(t, `{"enabled": false}`))
	require.NoError(t, err)
	assert.Empty(t, fe)
	view, err = admin.Settings(school.ID)
	require.NoError(t, err)
	assert.False(t, view.Enabled)
	assert.False(t, aiReloadSetting(t, db, school.ID).Enabled)
	// 未带 providers 字段 → 按 Laravel `input('providers', [])` 语义被置空（详见下一个测试）
	assert.Empty(t, view.Providers)
}

func TestAIAdminSaveSettingsApiKeyRetentionCounterMergeAndWipe(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAIAdmin(db, services.NewAIService())

	fe, err := admin.SaveSettings(school.ID, aiDecodeRaw(t, `{"providers":[{
		"id":"openai","label":"OpenAI","api_key":"sk-first","is_active":true,
		"tokens_used":100,"total_calls":5,"estimated_cost":1.5,"currency":"USD"
	}]}`))
	require.NoError(t, err)
	assert.Empty(t, fe)

	// 第二次保存：api_key 传空串 + 前端旧缓存把计数器清零 → 应保留 Key 与较大计数
	fe, err = admin.SaveSettings(school.ID, aiDecodeRaw(t, `{"providers":[{
		"id":"openai","label":"OpenAI","api_key":"","is_active":true,
		"tokens_used":0,"total_calls":0,"estimated_cost":0
	}]}`))
	require.NoError(t, err)
	assert.Empty(t, fe)

	view, err := admin.Settings(school.ID)
	require.NoError(t, err)
	require.Len(t, view.Providers, 1)
	p := view.Providers[0]
	assert.Equal(t, "sk-first", p["api_key"], "空 api_key 必须保留原值")
	assert.Equal(t, 100.0, p["tokens_used"])
	assert.Equal(t, 5.0, p["total_calls"])
	assert.Equal(t, 1.5, p["estimated_cost"])
	assert.Equal(t, "USD", p["currency"], "currency 缺失时沿用旧值")

	// 未带某供应商 id 的新供应商：不被旧配置影响
	fe, err = admin.SaveSettings(school.ID, aiDecodeRaw(t, `{"providers":[
		{"id":"openai","label":"OpenAI","api_key":"","is_active":true},
		{"id":"claude","label":"Claude","api_key":"sk-claude","is_active":false}
	]}`))
	require.NoError(t, err)
	assert.Empty(t, fe)
	view, err = admin.Settings(school.ID)
	require.NoError(t, err)
	require.Len(t, view.Providers, 2)
	assert.Equal(t, "sk-first", view.Providers[0]["api_key"])
	assert.Equal(t, "sk-claude", view.Providers[1]["api_key"])
	assert.Equal(t, false, view.Providers[1]["is_active"])

	// 请求未带 providers 字段 → 按 Laravel `input('providers', [])` 语义置空
	fe, err = admin.SaveSettings(school.ID, aiDecodeRaw(t, `{"max_tokens":3000}`))
	require.NoError(t, err)
	assert.Empty(t, fe)
	view, err = admin.Settings(school.ID)
	require.NoError(t, err)
	assert.Empty(t, view.Providers, "未传 providers 时应清空（同 Laravel）")
	assert.Equal(t, 3000, view.MaxTokens)
}

func TestAIAdminSaveSettingsValidation(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAIAdmin(db, services.NewAIService())

	cases := []struct {
		name string
		body string
		key  string
	}{
		{"enabled 非布尔", `{"enabled":"abc"}`, "enabled"},
		{"enabled 为 null", `{"enabled":null}`, "enabled"},
		{"max_tokens 过小", `{"max_tokens":50}`, "max_tokens"},
		{"max_tokens 非整数", `{"max_tokens":"x"}`, "max_tokens"},
		{"tokens_limit 为负", `{"tokens_limit":-1}`, "tokens_limit"},
		{"providers 非数组", `{"providers":"x"}`, "providers"},
		{"provider 缺 id", `{"providers":[{"label":"x"}]}`, "providers.0.id"},
		{"provider 缺 label", `{"providers":[{"id":"openai"}]}`, "providers.0.label"},
		{"is_active 非布尔", `{"providers":[{"id":"a","label":"b","is_active":"yes"}]}`, "providers.0.is_active"},
		{"单价为负", `{"providers":[{"id":"a","label":"b","input_price_per_m":-1}]}`, "providers.0.input_price_per_m"},
		{"models 元素非字符串", `{"providers":[{"id":"a","label":"b","models":[1]}]}`, "providers.0.models.0"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fe, err := admin.SaveSettings(school.ID, aiDecodeRaw(t, tc.body))
			require.NoError(t, err)
			require.NotEmpty(t, fe, "应返回校验错误")
			assert.NotEmpty(t, fe[tc.key], "缺少字段 %s 的报错：%v", tc.key, fe)
		})
	}

	// 校验失败不落库（validation 先于 firstOrCreate）
	var count int64
	require.NoError(t, db.Model(&models.AISetting{}).Where("school_id = ?", school.ID).Count(&count).Error)
	assert.Equal(t, int64(0), count)

	// null 的 max_tokens / tokens_limit 走 Laravel 的 has() 分支：(int) null = 0
	fe, err := admin.SaveSettings(school.ID, aiDecodeRaw(t, `{"max_tokens":null,"tokens_limit":null}`))
	require.NoError(t, err)
	assert.Empty(t, fe)
	setting := aiReloadSetting(t, db, school.ID)
	assert.Equal(t, 0, setting.MaxTokens)
	assert.Equal(t, 0, setting.TokensLimit)
}

func TestAIAdminToggle(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAIAdmin(db, services.NewAIService())

	enabled, err := admin.Toggle(school.ID, true)
	require.NoError(t, err)
	assert.True(t, enabled)
	assert.True(t, aiReloadSetting(t, db, school.ID).Enabled)

	enabled, err = admin.Toggle(school.ID, false)
	require.NoError(t, err)
	assert.False(t, enabled)
	assert.False(t, aiReloadSetting(t, db, school.ID).Enabled)

	// 只有一行设置（幂等 firstOrCreate）
	var count int64
	require.NoError(t, db.Model(&models.AISetting{}).Where("school_id = ?", school.ID).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

// ============================================================
// 2. fetch-models / test
// ============================================================

func TestAIAdminFetchModels(t *testing.T) {
	t.Run("openai 兼容", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		fake := aiNewFakeProvider(t, aiRespondJSON(200, `{"data":[{"id":"gpt-4o"},{"id":"gpt-4o-mini"},{"id":""}]}`))
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-x")))
		admin := services.NewAIAdmin(db, aiServiceFor(fake))

		view, err := admin.FetchModels(school.ID, "openai")
		require.NoError(t, err)
		assert.Equal(t, "openai", view.Provider)
		assert.Equal(t, []string{"gpt-4o", "gpt-4o-mini"}, view.Models, "空 id 被过滤")
		assert.NotEmpty(t, view.FetchedAt)

		req := fake.aiLast(t)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/models", req.Path)
		assert.Equal(t, "Bearer sk-x", req.Header.Get("Authorization"))
	})

	t.Run("claude", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		fake := aiNewFakeProvider(t, aiRespondJSON(200, `{"data":[{"id":"claude-3-5-sonnet"}]}`))
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("claude", fake.server.URL, "claude-3-5-sonnet", "key-claude")))
		admin := services.NewAIAdmin(db, aiServiceFor(fake))

		view, err := admin.FetchModels(school.ID, "claude")
		require.NoError(t, err)
		assert.Equal(t, []string{"claude-3-5-sonnet"}, view.Models)

		req := fake.aiLast(t)
		assert.Equal(t, "key-claude", req.Header.Get("x-api-key"))
		assert.Equal(t, "2023-06-01", req.Header.Get("anthropic-version"))
	})

	t.Run("google", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		fake := aiNewFakeProvider(t, aiRespondJSON(200, `{"models":[{"name":"models/gemini-1.5-pro"},{"name":"models/gemini-1.5-flash"}]}`))
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("google", fake.server.URL, "gemini-1.5-pro", "key-google")))
		admin := services.NewAIAdmin(db, aiServiceFor(fake))

		view, err := admin.FetchModels(school.ID, "google")
		require.NoError(t, err)
		assert.Equal(t, []string{"gemini-1.5-pro", "gemini-1.5-flash"}, view.Models, "models/ 前缀应剥离")

		req := fake.aiLast(t)
		assert.Equal(t, "key=key-google", req.Query)
	})

	t.Run("边界", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		admin := services.NewAIAdmin(db, services.NewAIService())

		_, err := admin.FetchModels(school.ID, "")
		assertAppStatus(t, err, 422, "缺少供应商应 422")
		assert.Contains(t, err.Error(), "缺少供应商")

		// 尚无设置行 → 视为未配置（Laravel 该分支会因读取 null 属性报 500，Go 端按业务语义返回 422）
		_, err = admin.FetchModels(school.ID, "openai")
		assertAppStatus(t, err, 422, "无设置行时应 422")
		assert.Equal(t, "该供应商未配置 API Key", err.Error())

		// 有供应商但没有 api_key
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t,
			map[string]any{"id": "openai", "label": "OpenAI", "api_key": "", "is_active": true}))
		_, err = admin.FetchModels(school.ID, "openai")
		assertAppStatus(t, err, 422, "无 api_key 应 422")
		assert.Equal(t, "该供应商未配置 API Key", err.Error())
	})
}

func TestAIAdminTestProvider(t *testing.T) {
	t.Run("成功", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		fake := aiNewFakeProvider(t, aiRespondJSON(200,
			`{"choices":[{"message":{"content":"你好呀，同学"}}],"usage":{"prompt_tokens":3,"completion_tokens":4}}`))
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-x")))
		admin := services.NewAIAdmin(db, aiServiceFor(fake))

		result, err := admin.Test(school.ID, "openai")
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, result.LatencyMs, 0)
		require.NotNil(t, result.Reply)
		assert.Equal(t, "你好呀，同学", *result.Reply)
		assert.Nil(t, result.Error)

		req := fake.aiLast(t)
		assert.Equal(t, "/chat/completions", req.Path)
		assert.Equal(t, "gpt-4o", req.Body["model"])
		assert.Equal(t, 16.0, req.Body["max_tokens"], "连通性测试固定 max_tokens=16")
		messages := req.Body["messages"].([]any)
		assert.Equal(t, "你好", messages[1].(map[string]any)["content"])

		// 不写库、不计费
		var count int64
		require.NoError(t, db.Model(&models.AIConversation{}).Count(&count).Error)
		assert.Equal(t, int64(0), count)
		assert.Equal(t, 0, aiReloadSetting(t, db, school.ID).TokensUsed)
	})

	t.Run("供应商返回错误", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		fake := aiNewFakeProvider(t, aiRespondJSON(500, `{"error":{"message":"bad key"}}`))
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-x")))
		admin := services.NewAIAdmin(db, aiServiceFor(fake))

		result, err := admin.Test(school.ID, "openai")
		require.NoError(t, err)
		assert.False(t, result.Success)
		// HTTP 失败分支（callXxx 的固定兜底文案）仍带 reply，只有 catch（传输层异常）分支不带
		require.NotNil(t, result.Reply)
		assert.Equal(t, "AI 服务暂时不可用", *result.Reply)
		assert.Equal(t, "供应商返回异常，请检查 API Key / 模型 / 地址", result.Error)
	})

	t.Run("传输层异常", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", "http://127.0.0.1:1/v1", "gpt-4o", "sk-x")))
		ai := services.NewAIService()
		ai.SetHTTPClient(aiErrorDoer{})
		admin := services.NewAIAdmin(db, ai)

		result, err := admin.Test(school.ID, "openai")
		require.NoError(t, err)
		assert.False(t, result.Success)
		assert.Nil(t, result.Reply)
		assert.Contains(t, result.Error, "connection refused")
	})

	t.Run("边界", func(t *testing.T) {
		db := setupDB(t)
		school := seedSchool(t, db)
		admin := services.NewAIAdmin(db, services.NewAIService())

		_, err := admin.Test(school.ID, "")
		assertAppStatus(t, err, 422, "缺少供应商")
		assert.Equal(t, "缺少供应商", err.Error())

		aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t,
			map[string]any{"id": "openai", "label": "OpenAI", "api_key": "", "model": "gpt-4o"}))
		_, err = admin.Test(school.ID, "openai")
		assertAppStatus(t, err, 422, "无 api_key")
		assert.Equal(t, "该供应商未配置 API Key，请先填写并保存", err.Error())

		// 有 Key 但没有模型
		aiSettingWith(t, db, schoolsSecond(t, db).ID, true, aiProvidersJSON(t,
			map[string]any{"id": "openai", "label": "OpenAI", "api_key": "sk-y", "model": ""}))
		other := schoolsSecond(t, db)
		_, err = admin.Test(other.ID, "openai")
		assertAppStatus(t, err, 422, "无模型")
		assert.Equal(t, "该供应商未配置模型，请先填写并保存", err.Error())
	})
}

// ============================================================
// 3. 教师端对话
// ============================================================

func TestAITeacherChatOpenAICompatible(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	teacher := aiSetTeacherActiveClass(t, db, seedTeacher(t, db, school.ID, "t1"), class.ID)

	fake := aiNewFakeProvider(t, aiRespondJSON(200,
		`{"choices":[{"message":{"content":"答案是 2"}}],"usage":{"prompt_tokens":10,"completion_tokens":20}}`))
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-test")))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	reply, err := assistant.Chat(&teacher, "1+1等于几？", "")
	require.NoError(t, err)
	assert.Equal(t, "答案是 2", reply)

	// 请求体关键字段
	req := fake.aiLast(t)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, "/chat/completions", req.Path)
	assert.Equal(t, "Bearer sk-test", req.Header.Get("Authorization"))
	assert.Equal(t, "gpt-4o", req.Body["model"])
	assert.Equal(t, 2000.0, req.Body["max_tokens"], "max_tokens 取 ai_settings.max_tokens")
	assert.Equal(t, 0.7, req.Body["temperature"])
	messages := req.Body["messages"].([]any)
	require.Len(t, messages, 2)
	assert.Equal(t, "system", messages[0].(map[string]any)["role"])
	assert.Equal(t, "你是一个学习助手，请用中文回答。", messages[0].(map[string]any)["content"])
	assert.Equal(t, "1+1等于几？", messages[1].(map[string]any)["content"])

	// 会话落库
	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, "教师", *conv.StudentName)
	assert.Equal(t, "openai", *conv.Provider)
	require.NotNil(t, conv.ClassID)
	assert.Equal(t, class.ID, *conv.ClassID)
	assert.Nil(t, conv.StudentID)
	assert.Equal(t, "completed", conv.Status)
	assert.Equal(t, "1+1等于几？", conv.Question)
	assert.Equal(t, "答案是 2", *conv.Answer)
	assert.Equal(t, 30, conv.TokensUsed)
	assert.Equal(t, 10, conv.PromptTokens)
	assert.Equal(t, 20, conv.CompletionTokens)
	assert.InDelta(t, 0.00005, conv.Cost, 1e-9, "cost = 10/1e6*1 + 20/1e6*2")
	assert.Equal(t, "USD", conv.Currency)

	// 学校 tokens_used 与供应商计数器累加
	setting := aiReloadSetting(t, db, school.ID)
	assert.Equal(t, 30, setting.TokensUsed)
	providers := aiReadProviders(t, setting.Providers)
	require.Len(t, providers, 1)
	assert.Equal(t, 30.0, providers[0]["tokens_used"])
	assert.Equal(t, 1.0, providers[0]["total_calls"])
	assert.InDelta(t, 0.00005, providers[0]["estimated_cost"].(float64), 1e-9)
	assert.Equal(t, "USD", providers[0]["currency"])
}

func TestAITeacherChatClaude(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-claude")

	fake := aiNewFakeProvider(t, aiRespondJSON(200, `{
		"content": [{"type":"text","text":"你好"},{"type":"tool_use","text":"忽略"},{"type":"text","text":"，世界"}],
		"usage": {"input_tokens":5,"output_tokens":7}
	}`))
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("claude", fake.server.URL, "claude-3-5-sonnet", "key-claude")))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	reply, err := assistant.Chat(&teacher, "打个招呼", "")
	require.NoError(t, err)
	assert.Equal(t, "你好，世界", reply, "只拼接 type=text 的块")

	req := fake.aiLast(t)
	assert.Equal(t, "/messages", req.Path)
	assert.Equal(t, "key-claude", req.Header.Get("x-api-key"))
	assert.Equal(t, "2023-06-01", req.Header.Get("anthropic-version"))
	assert.Equal(t, "claude-3-5-sonnet", req.Body["model"])
	assert.Equal(t, 2000.0, req.Body["max_tokens"])
	messages := req.Body["messages"].([]any)
	require.Len(t, messages, 1, "claude 请求体没有 system 消息")
	assert.Equal(t, "user", messages[0].(map[string]any)["role"])
	assert.Equal(t, "打个招呼", messages[0].(map[string]any)["content"])

	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, "claude", *conv.Provider)
	assert.Equal(t, 12, conv.TokensUsed)
	assert.Equal(t, 5, conv.PromptTokens)
	assert.Equal(t, 7, conv.CompletionTokens)
	assert.InDelta(t, 0.000019, conv.Cost, 1e-9) // 5/1e6*1 + 7/1e6*2
}

func TestAITeacherChatGemini(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-google")

	fake := aiNewFakeProvider(t, aiRespondJSON(200, `{
		"candidates": [{"content":{"parts":[{"text":"来自 Gemini"}]}}],
		"usageMetadata": {"promptTokenCount":3,"candidatesTokenCount":4}
	}`))
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("google", fake.server.URL, "gemini-1.5-pro", "key-google")))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	reply, err := assistant.Chat(&teacher, "讲个笑话", "")
	require.NoError(t, err)
	assert.Equal(t, "来自 Gemini", reply)

	req := fake.aiLast(t)
	assert.Equal(t, "/models/gemini-1.5-pro:generateContent", req.Path)
	assert.Equal(t, "key=key-google", req.Query, "api_key 走 query（同 Laravel）")
	assert.Equal(t, 2000.0, req.Body["generationConfig"].(map[string]any)["maxOutputTokens"])
	contents := req.Body["contents"].([]any)
	assert.Equal(t, "讲个笑话", contents[0].(map[string]any)["parts"].([]any)[0].(map[string]any)["text"])

	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, 7, conv.TokensUsed)
	assert.Equal(t, 3, conv.PromptTokens)
	assert.Equal(t, 4, conv.CompletionTokens)
}

func TestAITeacherChatQwen(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-qwen")

	fake := aiNewFakeProvider(t, aiRespondJSON(200, `{
		"output": {"text":"来自通义"},
		"usage": {"input_tokens":2,"output_tokens":6}
	}`))
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("qwen", fake.server.URL, "qwen-turbo", "key-qwen")))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	reply, err := assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "来自通义", reply)

	req := fake.aiLast(t)
	assert.Equal(t, "/services/aigc/text-generation/generation", req.Path)
	assert.Equal(t, "Bearer key-qwen", req.Header.Get("Authorization"))
	input := req.Body["input"].(map[string]any)
	messages := input["messages"].([]any)
	assert.Equal(t, "你是一个学习助手。", messages[0].(map[string]any)["content"], "qwen 的 system 文案与其他供应商不同")
	assert.Equal(t, 2000.0, req.Body["parameters"].(map[string]any)["max_tokens"])

	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, 8, conv.TokensUsed)
}

func TestAITeacherChatUpstreamError(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-err")

	fake := aiNewFakeProvider(t, aiRespondJSON(429, `{"error":{"message":"rate limited"}}`))
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-x")))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	reply, err := assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "AI 服务暂时不可用", reply, "上游失败文案与 Laravel 逐字一致")

	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, "completed", conv.Status)
	assert.Equal(t, "AI 服务暂时不可用", *conv.Answer)
	assert.Equal(t, 0, conv.TokensUsed)
	assert.Equal(t, 0.0, conv.Cost)
	assert.Equal(t, 0, aiReloadSetting(t, db, school.ID).TokensUsed, "无 token 消耗时不累加")
}

func TestAITeacherChatTransportError(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-net")

	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", "http://127.0.0.1:1/v1", "gpt-4o", "sk-x")))
	ai := services.NewAIService()
	ai.SetHTTPClient(aiErrorDoer{})
	assistant := services.NewAIAssistant(db, ai)

	reply, err := assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "AI 服务暂时不可用", reply)

	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, 0, conv.TokensUsed)
}

func TestAITeacherChatNotConfigured(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-none")
	assistant := services.NewAIAssistant(db, services.NewAIService())

	// 完全没有设置行
	reply, err := assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "AI 功能未启用，请联系管理员配置", reply)

	// 有设置行但未启用（供应商只填了 is_active、没有 api_key，避免任何真实出网）
	keyless := map[string]any{"id": "openai", "label": "OpenAI", "api_key": "", "model": "gpt-4o", "is_active": true}
	aiSettingWith(t, db, school.ID, false, aiProvidersJSON(t, keyless))
	reply, err = assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "AI 功能未启用，请联系管理员配置", reply)

	// 已启用但没有可用供应商（api_key 为空）
	require.NoError(t, db.Model(&models.AISetting{}).Where("school_id = ?", school.ID).Update("enabled", true).Error)
	reply, err = assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "请先在 AI 中心配置并启用一个供应商", reply)

	// 已启用 + 供应商未勾选启用（is_active=false）→ 同样视为未配置
	inactive := map[string]any{"id": "openai", "label": "OpenAI", "api_key": "sk-x", "is_active": false}
	require.NoError(t, db.Model(&models.AISetting{}).Where("school_id = ?", school.ID).
		Update("providers", aiProvidersJSON(t, inactive)).Error)
	reply, err = assistant.Chat(&teacher, "你好", "")
	require.NoError(t, err)
	assert.Equal(t, "请先在 AI 中心配置并启用一个供应商", reply)

	// 全程不落会话
	var count int64
	require.NoError(t, db.Model(&models.AIConversation{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestAITeacherModelWhitelistAndMap(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-model")

	fake := aiNewFakeProvider(t, aiRespondJSON(200, `{"choices":[{"message":{"content":"ok"}}]}`))
	entry := aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-x")
	entry["models"] = []any{"gpt-4o-mini"}
	entry["model_map"] = map[string]any{
		"gpt-4o":      "gpt-4o-2024-08-06",
		"gpt-4o-mini": "gpt-4o-mini-2024-07-18",
	}
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, entry))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	// 白名单内 + 命中 model_map → 重定向到上游模型
	_, err := assistant.Chat(&teacher, "你好", "gpt-4o-mini")
	require.NoError(t, err)
	assert.Equal(t, "gpt-4o-mini-2024-07-18", fake.aiLast(t).Body["model"])

	// 白名单外 → 回退默认主模型，再按其 model_map 重定向
	_, err = assistant.Chat(&teacher, "你好", "not-in-whitelist")
	require.NoError(t, err)
	assert.Equal(t, "gpt-4o-2024-08-06", fake.aiLast(t).Body["model"])

	assert.Equal(t, 2, fake.aiReqs())
}

func TestAITeacherConfigAndUsage(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t-usage")
	assistant := services.NewAIAssistant(db, services.NewAIService())

	// 无设置行
	config, err := assistant.ConfigFor(&teacher)
	require.NoError(t, err)
	assert.False(t, config.Enabled)

	usage, err := assistant.UsageFor(&teacher)
	require.NoError(t, err)
	assert.False(t, usage.Configured)
	assert.Nil(t, usage.Provider)
	assert.Nil(t, usage.Model)
	assert.Empty(t, usage.Models)
	assert.Equal(t, 0, usage.TokensUsed)
	assert.Equal(t, 0.0, usage.EstimatedCost)
	assert.Equal(t, "CNY", usage.Currency)

	// 启用但无 key → enabled=false
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t,
		map[string]any{"id": "openai", "label": "OpenAI", "api_key": "", "model": "gpt-4o", "is_active": true}))
	config, err = assistant.ConfigFor(&teacher)
	require.NoError(t, err)
	assert.False(t, config.Enabled)

	// 启用 + 有 key → enabled=true
	require.NoError(t, db.Model(&models.AISetting{}).Where("school_id = ?", school.ID).Updates(map[string]any{
		"providers": aiProvidersJSON(t, aiProviderEntry("openai", "http://127.0.0.1:1/v1", "gpt-4o", "sk-x")),
	}).Error)
	config, err = assistant.ConfigFor(&teacher)
	require.NoError(t, err)
	assert.True(t, config.Enabled)

	require.NoError(t, db.Model(&models.AISetting{}).Where("school_id = ?", school.ID).Update("tokens_used", 20).Error)
	usage, err = assistant.UsageFor(&teacher)
	require.NoError(t, err)
	assert.True(t, usage.Configured)
	require.NotNil(t, usage.Provider)
	assert.Equal(t, "openai", *usage.Provider)
	require.NotNil(t, usage.Model)
	assert.Equal(t, "gpt-4o", *usage.Model)
	assert.Equal(t, []string{"gpt-4o"}, usage.Models)
	assert.Equal(t, 20, usage.TokensUsed)
	assert.Equal(t, 0.0, usage.EstimatedCost)
	assert.Equal(t, "USD", usage.Currency)
}

func TestAITeacherCommands(t *testing.T) {
	db := setupDB(t)
	assistant := services.NewAIAssistant(db, services.NewAIService())

	commands := assistant.Commands()
	require.Len(t, commands, 4)
	assert.Equal(t, "📝 本周教学总结", commands[0].Label)
	assert.Equal(t, "请帮我写一份本周教学总结，包含本周教学目标、课堂情况、学生表现和下周教学计划。", commands[0].Prompt)
	assert.Equal(t, "🏅 积分规则建议", commands[1].Label)
	assert.Equal(t, "🎯 班会活动方案", commands[2].Label)
	assert.Equal(t, "📋 出练习题", commands[3].Label)
	assert.Equal(t, "请出一组适合本年级学生的练习题，包含题目和参考答案。", commands[3].Prompt)
}

// ============================================================
// 4. 大屏端
// ============================================================

func TestAIDisplaySettingsCheck(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	assistant := services.NewAIAssistant(db, services.NewAIService())

	view, err := assistant.DisplaySettingsCheck(class.ID)
	require.NoError(t, err)
	assert.False(t, view.Enabled)

	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", "http://127.0.0.1:1/v1", "gpt-4o", "sk-x")))
	view, err = assistant.DisplaySettingsCheck(class.ID)
	require.NoError(t, err)
	assert.True(t, view.Enabled)

	// 班级不存在 → enabled=false（不报错，同 Laravel）
	view, err = assistant.DisplaySettingsCheck(99999)
	require.NoError(t, err)
	assert.False(t, view.Enabled)
}

func TestAIDisplayChat(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)

	fake := aiNewFakeProvider(t, aiRespondJSON(200,
		`{"choices":[{"message":{"content":"大屏回答"}}],"usage":{"prompt_tokens":4,"completion_tokens":6}}`))
	aiSettingWith(t, db, school.ID, true, aiProvidersJSON(t, aiProviderEntry("openai", fake.server.URL, "gpt-4o", "sk-x")))
	assistant := services.NewAIAssistant(db, aiServiceFor(fake))

	result, err := assistant.DisplayChat(class.ID, "为什么天是蓝的？")
	require.NoError(t, err)
	assert.Equal(t, "大屏回答", result.Answer)
	assert.Equal(t, 10, result.TokensUsed)

	req := fake.aiLast(t)
	assert.Equal(t, "/chat/completions", req.Path)
	assert.Equal(t, "gpt-4o", req.Body["model"], "大屏端不做白名单/model_map 解析")

	var conv models.AIConversation
	require.NoError(t, db.First(&conv).Error)
	assert.Equal(t, "匿名", *conv.StudentName)
	assert.Nil(t, conv.StudentID)
	require.NotNil(t, conv.ClassID)
	assert.Equal(t, class.ID, *conv.ClassID)
	assert.Equal(t, "openai", *conv.Provider)
	assert.Equal(t, "completed", conv.Status)
	assert.Equal(t, 10, conv.TokensUsed)
	assert.Equal(t, 0, aiReloadSetting(t, db, school.ID).TokensUsed-10, "大屏端同样累加学校 tokens_used")
}

func TestAIDisplayChatGuards(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	assistant := services.NewAIAssistant(db, services.NewAIService())

	// 未启用 → 403「AI 功能未开启」
	_, err := assistant.DisplayChat(class.ID, "你好")
	assertAppStatus(t, err, 403, "未启用应 403")
	assert.Equal(t, "AI 功能未开启", err.Error())

	// 启用但无供应商 → 403
	aiSettingWith(t, db, school.ID, true, "")
	_, err = assistant.DisplayChat(class.ID, "你好")
	assertAppStatus(t, err, 403, "无供应商应 403")
	assert.Equal(t, "请先在 AI 中心配置并启用一个供应商", err.Error())

	// 有供应商 + 空问题 → 422「请输入问题」（note：PHP empty("0") 亦为真）
	aiSettingWith(t, db, schoolsSecond(t, db).ID, true, aiProvidersJSON(t, aiProviderEntry("openai", "http://127.0.0.1:1/v1", "gpt-4o", "sk-x")))
	other := schoolsSecond(t, db)
	otherClass := seedClass(t, db, other.ID)
	for _, question := range []string{"", "0"} {
		_, err = assistant.DisplayChat(otherClass.ID, question)
		assertAppStatus(t, err, 422, "空问题应 422")
		assert.Equal(t, "请输入问题", err.Error())
	}

	// 班级不存在 → 404「班级不存在」
	_, err = assistant.DisplayChat(99999, "你好")
	assertAppStatus(t, err, 404, "班级不存在应 404")
	assert.Equal(t, "班级不存在", err.Error())

	// 全程不落会话
	var count int64
	require.NoError(t, db.Model(&models.AIConversation{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

// ============================================================
// 5. 管理端用量统计
// ============================================================

func TestAIAdminUsage(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)

	providers := aiProvidersJSON(t, map[string]any{
		"id": "openai", "label": "OpenAI", "api_key": "sk-x", "is_active": true,
		"tokens_used": 120, "total_calls": 3, "estimated_cost": 0.5,
		"cost_per_token": 0.004, "currency": "USD",
	}, map[string]any{"label": "无 id 的供应商", "api_key": "sk-y"})
	aiSettingWith(t, db, school.ID, true, providers)

	day := util.StartOfDay(util.Now())
	seedAIConversation(t, db, school.ID, day.Add(-48*time.Hour).Add(12*time.Hour), 30, 0.25, "openai", "教师", "问题一", "回答一")
	seedAIConversation(t, db, school.ID, day.Add(-24*time.Hour).Add(9*time.Hour), 20, 0.15, "openai", "匿名", "问题二", "回答二")
	seedAIConversation(t, db, school.ID, day.Add(8*time.Hour), 10, 0.05, "claude", "教师", "问题三", "")
	seedAIConversation(t, db, school.ID, day.AddDate(0, 0, -30).Add(10*time.Hour), 99, 9.9, "openai", "教师", "范围外", "范围外")

	admin := services.NewAIAdmin(db, services.NewAIService())
	usage, err := admin.Usage(school.ID, 7)
	require.NoError(t, err)

	assert.True(t, usage.Enabled)
	assert.Equal(t, 0, usage.TokensUsed)
	assert.Equal(t, 1000000, usage.TokensLimit)
	assert.InDelta(t, 0.5, usage.EstimatedCost, 1e-9)
	assert.Equal(t, int64(4), usage.TotalConversations, "总数为全量计数（含范围外）")

	require.Len(t, usage.DailyUsage, 3, "只有窗口内的 3 天")
	assert.Equal(t, day.Add(-48*time.Hour).Format("2006-01-02"), usage.DailyUsage[0].Date)
	assert.Equal(t, int64(30), usage.DailyUsage[0].Tokens)
	assert.Equal(t, int64(1), usage.DailyUsage[0].Count)
	assert.InDelta(t, 0.25, usage.DailyUsage[0].Cost, 1e-9)
	assert.Equal(t, day.Format("2006-01-02"), usage.DailyUsage[2].Date, "按日期升序")
	assert.Equal(t, int64(10), usage.DailyUsage[2].Tokens)

	require.Len(t, usage.RecentLogs, 3)
	assert.Equal(t, "问题三", usage.RecentLogs[0].Question, "最近一条在最前")
	assert.Equal(t, day.Add(8*time.Hour).In(util.Loc).Format("2006-01-02 15:04:05"), usage.RecentLogs[0].CreatedAt)
	assert.InDelta(t, 0.05, usage.RecentLogs[0].Cost, 1e-9)
	assert.Equal(t, "claude", *usage.RecentLogs[0].Provider)
	assert.Equal(t, "教师", *usage.RecentLogs[0].StudentName)

	require.Len(t, usage.ByProvider, 1, "无 id 的供应商被跳过")
	openai := usage.ByProvider["openai"]
	assert.Equal(t, 120.0, openai.Tokens)
	assert.Equal(t, 3.0, openai.TotalCalls)
	assert.Equal(t, 0.5, openai.EstimatedCost)
	assert.Equal(t, 0.004, openai.CostPerToken)
	assert.Equal(t, "USD", openai.Currency)

	// 无设置行时不报错，全部取默认值（Laravel 该分支会因读取 null 属性报 500，Go 端按业务语义返回空态）
	other := schoolsSecond(t, db)
	empty, err := admin.Usage(other.ID, 7)
	require.NoError(t, err)
	assert.False(t, empty.Enabled)
	assert.Equal(t, 0, empty.TokensUsed)
	assert.Equal(t, 0, empty.TokensLimit)
	assert.Empty(t, empty.DailyUsage)
	assert.Empty(t, empty.ByProvider)
	assert.Empty(t, empty.RecentLogs)

	// by_provider 的 JSON 形状：空表输出 []（PHP 空数组），非空输出对象
	emptyJSON, err := json.Marshal(empty)
	require.NoError(t, err)
	assert.Contains(t, string(emptyJSON), `"by_provider":[]`)
	fullJSON, err := json.Marshal(usage)
	require.NoError(t, err)
	assert.Contains(t, string(fullJSON), `"by_provider":{"openai":`)
}

// seedAIConversation 写入一条对话流水（可控 created_at）。
func seedAIConversation(t *testing.T, db *gorm.DB, schoolID uint, createdAt time.Time, tokens int, cost float64, provider, studentName, question, answer string) {
	t.Helper()
	conv := models.AIConversation{
		SchoolID:         schoolID,
		StudentName:      &studentName,
		Provider:         &provider,
		Question:         question,
		TokensUsed:       tokens,
		PromptTokens:     tokens / 3,
		CompletionTokens: tokens - tokens/3,
		Cost:             cost,
		Currency:         "USD",
		Status:           "completed",
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
	}
	if answer != "" {
		conv.Answer = &answer
	}
	require.NoError(t, db.Create(&conv).Error)
}

// schoolsSecond 第二所学校（避免唯一 code 冲突）。
func schoolsSecond(t *testing.T, db *gorm.DB) models.School {
	t.Helper()
	var school models.School
	err := db.Where("code = ?", "ai-second-school").First(&school).Error
	if err == nil {
		return school
	}
	school = models.School{Name: "第二学校", Code: "ai-second-school", Status: "active"}
	require.NoError(t, db.Create(&school).Error)
	return school
}
