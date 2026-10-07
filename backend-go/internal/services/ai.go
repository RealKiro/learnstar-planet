// AI 供应商调用服务（多供应商分发 + api_base 覆盖 + 各供应商请求体差异）。
//
// 忠实移植自 Laravel App\Services\AiService：
//   - chat 按供应商 id 分发到专属实现（claude / google / qwen / mcp），其余供应商走 OpenAI 兼容请求体；
//   - api_base 覆盖默认地址，未配置时按各供应商默认地址（qwen 的对话地址与模型列表地址不同）；
//   - HTTP 状态码 >= 400 时返回固定兜底文案（与 Laravel `$response->failed()` 分支逐字一致），
//     传输层异常（DNS/连接/超时）向上抛出 error，由调用方决定兜底行为（同 Laravel 的 ConnectionException）。
//
// 无第三方依赖：出网用标准库 net/http，JSON 用 encoding/json；HTTP 客户端可注入（SetHTTPClient），
// 便于测试用 httptest.NewServer 假造供应商响应（api_base 指向假服务器即生效）。
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 调用超时（同 Laravel）：chat 为 Http::timeout(30)，listModels 为 Http::timeout(15)。
const (
	aiChatTimeout       = 30 * time.Second
	aiListModelsTimeout = 15 * time.Second
)

// AIResult 一次 AI 对话结果（对应 Laravel AiService 返回的 answer/tokens_used/prompt_tokens/completion_tokens）。
type AIResult struct {
	Answer           string
	TokensUsed       int
	PromptTokens     int
	CompletionTokens int
}

// HTTPDoer 是 HTTP 客户端抽象（*http.Client 天然满足），测试可注入 httptest 的 client 或桩实现。
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// AIService 多供应商 AI 调用（无状态，可直接复用）。
type AIService struct {
	client HTTPDoer
}

// NewAIService 创建 AI 调用服务（默认标准库 http.Client，超时由每次请求的 context 控制）。
func NewAIService() *AIService {
	return &AIService{client: &http.Client{}}
}

// SetHTTPClient 注入 HTTP 客户端（测试用 httptest.NewServer 的 client）。
func (s *AIService) SetHTTPClient(c HTTPDoer) {
	if c != nil {
		s.client = c
	}
}

// aiProviderCaps 描述某供应商的调用方式与默认地址。
type aiProviderCaps struct {
	flavor   string // openai | claude | google | qwen | mcp
	base     string // 对话默认地址
	listBase string // 模型列表默认地址（空表示同 base）
}

// aiProviderTable 与 Laravel AiService::defaultBase + callXxx 的默认地址逐条对齐。
var aiProviderTable = map[string]aiProviderCaps{
	"openai":      {flavor: "openai", base: "https://api.openai.com/v1"},
	"claude":      {flavor: "claude", base: "https://api.anthropic.com/v1"},
	"google":      {flavor: "google", base: "https://generativelanguage.googleapis.com/v1"},
	"qwen":        {flavor: "qwen", base: "https://dashscope.aliyuncs.com/api/v1", listBase: "https://dashscope.aliyuncs.com/compatible-mode/v1"},
	"deepseek":    {flavor: "openai", base: "https://api.deepseek.com/v1"},
	"moonshot":    {flavor: "openai", base: "https://api.moonshot.cn/v1"},
	"grok":        {flavor: "openai", base: "https://api.x.ai/v1"},
	"siliconflow": {flavor: "openai", base: "https://api.siliconflow.cn/v1"},
	"nvidia":      {flavor: "openai", base: "https://integrate.api.nvidia.com/v1"},
	"openrouter":  {flavor: "openai", base: "https://openrouter.ai/api/v1"},
	"bytedance":   {flavor: "openai", base: "https://ark.cn-beijing.volces.com/api/v3"},
	"minimax":     {flavor: "openai", base: "https://api.minimax.chat/v1"},
	"baichuan":    {flavor: "openai", base: "https://api.baichuan-ai.com/v1"},
	"stepfun":     {flavor: "openai", base: "https://api.stepfun.com/v1"},
	"lingyi":      {flavor: "openai", base: "https://api.lingyiwanwu.com/v1"},
	"mistral":     {flavor: "openai", base: "https://api.mistral.ai/v1"},
	"cohere":      {flavor: "openai", base: "https://api.cohere.ai/v1"},
	"perplexity":  {flavor: "openai", base: "https://api.perplexity.ai/v1"},
	"ai21":        {flavor: "openai", base: "https://api.ai21.com/studio/v1"},
	"together":    {flavor: "openai", base: "https://api.together.xyz/v1"},
	"fireworks":   {flavor: "openai", base: "https://api.fireworks.ai/inference/v1"},
	"groq":        {flavor: "openai", base: "https://api.groq.com/openai/v1"},
	"replicate":   {flavor: "openai", base: "https://api.replicate.com/v1"},
	"anyscale":    {flavor: "openai", base: "https://api.endpoints.anyscale.com/v1"},
	"azure":       {flavor: "openai", base: "https://models.inference.ai.azure.com/v1"},
	"ollama":      {flavor: "openai", base: "http://localhost:11434/v1"},
	"vllm":        {flavor: "openai", base: "http://localhost:8000/v1"},
	// mcp 无默认地址：未配置 api_base 时直接在 chat 中返回提示（同 callMcp）。
	"mcp": {flavor: "mcp"},
}

// aiCapsFor 返回供应商调用参数；未登记的供应商按 OpenAI 兼容兜底（同 PHP 的 default 分支）。
func aiCapsFor(provider string) aiProviderCaps {
	if caps, ok := aiProviderTable[provider]; ok {
		return caps
	}
	return aiProviderCaps{flavor: "openai", base: "https://api.openai.com/v1"}
}

// aiBaseOr 等价 PHP `$apiBase ?: $def`（空串视为未配置）。
func aiBaseOr(apiBase, def string) string {
	if apiBase != "" {
		return apiBase
	}
	return def
}

// Chat 发起一次对话。maxTokens 为单次最大输出 token（Laravel 传 ai_settings.max_tokens）。
func (s *AIService) Chat(provider, apiKey, model, question, apiBase string, maxTokens int) (AIResult, error) {
	caps := aiCapsFor(provider)

	switch caps.flavor {
	case "claude":
		return s.callClaude(apiKey, model, question, aiBaseOr(apiBase, caps.base), maxTokens)
	case "qwen":
		return s.callQwen(apiKey, model, question, aiBaseOr(apiBase, caps.base), maxTokens)
	case "google":
		return s.callGoogle(apiKey, model, question, aiBaseOr(apiBase, caps.base), maxTokens)
	case "mcp":
		if apiBase == "" {
			return AIResult{Answer: "MCP 接口需要配置 API 地址"}, nil
		}
		return s.callOpenAICompatible(apiKey, model, question, apiBase, maxTokens)
	default:
		return s.callOpenAICompatible(apiKey, model, question, aiBaseOr(apiBase, caps.base), maxTokens)
	}
}

// ListModels 从官方 API 拉取供应商模型列表（CC Switch 风格）。
// 任何异常（含未配置 api_key）一律返回空列表，与 Laravel try/catch 后的行为一致。
func (s *AIService) ListModels(provider, apiKey, apiBase string) []string {
	ids := []string{}
	if apiKey == "" {
		return ids
	}

	caps := aiCapsFor(provider)
	base := strings.TrimRight(aiBaseOr(apiBase, aiBaseOr(caps.listBase, caps.base)), "/")

	var body map[string]any
	var status int
	var err error

	switch provider {
	case "claude":
		body, status, err = s.doJSON(http.MethodGet, base+"/models", map[string]string{
			"x-api-key":         apiKey,
			"anthropic-version": "2023-06-01",
		}, nil, aiListModelsTimeout)
	case "google":
		body, status, err = s.doJSON(http.MethodGet, base+"/models?key="+apiKey, nil, nil, aiListModelsTimeout)
	default:
		body, status, err = s.doJSON(http.MethodGet, base+"/models", map[string]string{
			"Authorization": "Bearer " + apiKey,
		}, nil, aiListModelsTimeout)
	}
	if err != nil || status >= 400 || body == nil {
		return ids
	}

	switch provider {
	case "google":
		if list, ok := body["models"].([]any); ok {
			for _, item := range list {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				name, _ := m["name"].(string)
				if id := strings.Replace(name, "models/", "", 1); id != "" {
					ids = append(ids, id)
				}
			}
		}
	default:
		if list, ok := body["data"].([]any); ok {
			for _, item := range list {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if id, ok := m["id"].(string); ok && id != "" {
					ids = append(ids, id)
				}
			}
		}
	}

	return ids
}

// ============================================================
// 各供应商请求体与响应解析
// ============================================================

// aiMessage OpenAI 风格的对话消息（claude / qwen 复用同一形状）。
type aiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatRequest struct {
	Model       string      `json:"model"`
	Messages    []aiMessage `json:"messages"`
	MaxTokens   int         `json:"max_tokens"`
	Temperature float64     `json:"temperature"`
}

// callOpenAICompatible OpenAI 兼容供应商（含 openai 自身与大多国内厂商）。
func (s *AIService) callOpenAICompatible(apiKey, model, question, base string, maxTokens int) (AIResult, error) {
	payload := openAIChatRequest{
		Model: model,
		Messages: []aiMessage{
			{Role: "system", Content: "你是一个学习助手，请用中文回答。"},
			{Role: "user", Content: question},
		},
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	}

	body, status, err := s.doJSON(http.MethodPost, base+"/chat/completions", map[string]string{
		"Authorization": "Bearer " + apiKey,
		"Content-Type":  "application/json",
	}, payload, aiChatTimeout)
	if err != nil {
		return AIResult{}, err
	}
	if status >= 400 {
		return AIResult{Answer: "AI 服务暂时不可用"}, nil
	}

	usage, _ := body["usage"].(map[string]any)
	prompt := aiFirstInt(usage, 0, "prompt_tokens", "total_tokens")
	completion := aiFirstInt(usage, 0, "completion_tokens")

	return AIResult{
		Answer:           aiDigString(body, "抱歉，无法回答。", "choices", "0", "message", "content"),
		TokensUsed:       prompt + completion,
		PromptTokens:     prompt,
		CompletionTokens: completion,
	}, nil
}

type claudeChatRequest struct {
	Model     string      `json:"model"`
	MaxTokens int         `json:"max_tokens"`
	Messages  []aiMessage `json:"messages"`
}

func (s *AIService) callClaude(apiKey, model, question, base string, maxTokens int) (AIResult, error) {
	payload := claudeChatRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  []aiMessage{{Role: "user", Content: question}},
	}

	body, status, err := s.doJSON(http.MethodPost, base+"/messages", map[string]string{
		"x-api-key":         apiKey,
		"anthropic-version": "2023-06-01",
		"Content-Type":      "application/json",
	}, payload, aiChatTimeout)
	if err != nil {
		return AIResult{}, err
	}
	if status >= 400 {
		return AIResult{Answer: "AI 服务不可用"}, nil
	}

	answer := ""
	if blocks, ok := body["content"].([]any); ok {
		for _, item := range blocks {
			block, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if block["type"] == "text" {
				if text, ok := block["text"].(string); ok {
					answer += text
				}
			}
		}
	}
	if answer == "" {
		answer = "抱歉，无法回答。"
	}

	usage, _ := body["usage"].(map[string]any)
	prompt := aiFirstInt(usage, 0, "input_tokens")
	completion := aiFirstInt(usage, 0, "output_tokens")

	return AIResult{
		Answer:           answer,
		TokensUsed:       prompt + completion,
		PromptTokens:     prompt,
		CompletionTokens: completion,
	}, nil
}

type qwenChatRequest struct {
	Model      string        `json:"model"`
	Input      qwenInput     `json:"input"`
	Parameters qwenParameter `json:"parameters"`
}

type qwenInput struct {
	Messages []aiMessage `json:"messages"`
}

type qwenParameter struct {
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

func (s *AIService) callQwen(apiKey, model, question, base string, maxTokens int) (AIResult, error) {
	payload := qwenChatRequest{
		Model: model,
		Input: qwenInput{Messages: []aiMessage{
			{Role: "system", Content: "你是一个学习助手。"},
			{Role: "user", Content: question},
		}},
		Parameters: qwenParameter{MaxTokens: maxTokens, Temperature: 0.7},
	}

	body, status, err := s.doJSON(http.MethodPost, base+"/services/aigc/text-generation/generation", map[string]string{
		"Authorization": "Bearer " + apiKey,
		"Content-Type":  "application/json",
	}, payload, aiChatTimeout)
	if err != nil {
		return AIResult{}, err
	}
	if status >= 400 {
		return AIResult{Answer: "AI 服务不可用"}, nil
	}

	usage, _ := body["usage"].(map[string]any)
	prompt := aiFirstInt(usage, 0, "input_tokens")
	completion := aiFirstInt(usage, 0, "output_tokens")

	return AIResult{
		Answer:           aiDigString(body, "抱歉，无法回答。", "output", "text"),
		TokensUsed:       prompt + completion,
		PromptTokens:     prompt,
		CompletionTokens: completion,
	}, nil
}

type googleChatRequest struct {
	Contents         []googleContent        `json:"contents"`
	GenerationConfig googleGenerationConfig `json:"generationConfig"`
}

type googleContent struct {
	Parts []googlePart `json:"parts"`
}

type googlePart struct {
	Text string `json:"text"`
}

type googleGenerationConfig struct {
	MaxOutputTokens int     `json:"maxOutputTokens"`
	Temperature     float64 `json:"temperature"`
}

func (s *AIService) callGoogle(apiKey, model, question, base string, maxTokens int) (AIResult, error) {
	payload := googleChatRequest{
		Contents: []googleContent{{Parts: []googlePart{{Text: question}}}},
		GenerationConfig: googleGenerationConfig{
			MaxOutputTokens: maxTokens,
			Temperature:     0.7,
		},
	}

	body, status, err := s.doJSON(http.MethodPost,
		base+"/models/"+model+":generateContent?key="+apiKey, nil, payload, aiChatTimeout)
	if err != nil {
		return AIResult{}, err
	}
	if status >= 400 {
		return AIResult{Answer: "AI 服务不可用"}, nil
	}

	answer := aiDigString(body, "抱歉，无法回答。", "candidates", "0", "content", "parts", "0", "text")
	meta, _ := body["usageMetadata"].(map[string]any)
	prompt := aiFirstInt(meta, 0, "promptTokenCount")
	completion := aiFirstInt(meta, 0, "candidatesTokenCount")

	return AIResult{
		Answer:           answer,
		TokensUsed:       prompt + completion,
		PromptTokens:     prompt,
		CompletionTokens: completion,
	}, nil
}

// ============================================================
// 内部辅助
// ============================================================

// doJSON 发送请求（payload 为 nil 时不带请求体、不设 Content-Type），返回解析后的 JSON 与状态码。
// 传输层异常返回 error；响应体非法 JSON 时 body 为 nil（同 Laravel `$response->json()` 返回 null）。
func (s *AIService) doJSON(method, url string, headers map[string]string, payload any, timeout time.Duration) (map[string]any, int, error) {
	var reader io.Reader
	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, err
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if payload != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	resp, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, _ := io.ReadAll(resp.Body)

	var body map[string]any
	if len(data) > 0 {
		_ = json.Unmarshal(data, &body)
	}
	return body, resp.StatusCode, nil
}

// aiDigString 按路径取值并转字符串；缺失/为 null/非字符串时返回 def（同 PHP `?? 兜底`）。
func aiDigString(root map[string]any, def string, path ...string) string {
	var cur any = root
	for _, key := range path {
		switch node := cur.(type) {
		case map[string]any:
			cur = node[key]
		case []any:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(node) {
				return def
			}
			cur = node[idx]
		default:
			return def
		}
	}
	if s, ok := cur.(string); ok {
		return s
	}
	return def
}

// aiFirstInt 取首个存在且非 null 的整数键，全缺时返回 def（同 PHP `?? ?? def`）。
func aiFirstInt(m map[string]any, def int, keys ...string) int {
	for _, key := range keys {
		if v, ok := m[key]; ok && v != nil {
			if n, ok := aiToInt(v); ok {
				return n
			}
		}
	}
	return def
}

// aiToInt 把 JSON 解出的数值/数字字符串转成整数。
func aiToInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return int(i), true
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(n))
		if err != nil {
			return 0, false
		}
		return i, true
	}
	return 0, false
}
