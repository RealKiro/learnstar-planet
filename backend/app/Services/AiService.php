<?php

declare(strict_types=1);

namespace App\Services;

use Illuminate\Support\Facades\Http;

class AiService
{
    public function chat(string $provider, string $apiKey, string $model, string $question, ?string $apiBase = null, int $maxTokens = 2000): array
    {
        $method = 'call' . str_replace(' ', '', ucwords(str_replace(['-', '_'], ' ', $provider)));

        if (method_exists($this, $method)) {
            return $this->$method($apiKey, $model, $question, $apiBase, $maxTokens);
        }

        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase, $maxTokens);
    }

    /**
     * 各供应商默认 API 地址（OpenAI 兼容类）
     */
    private function defaultBase(string $provider): string
    {
        return match ($provider) {
            'openai' => 'https://api.openai.com/v1',
            'claude' => 'https://api.anthropic.com/v1',
            'google' => 'https://generativelanguage.googleapis.com/v1',
            'qwen' => 'https://dashscope.aliyuncs.com/compatible-mode/v1',
            'deepseek' => 'https://api.deepseek.com/v1',
            'moonshot' => 'https://api.moonshot.cn/v1',
            'grok' => 'https://api.x.ai/v1',
            'siliconflow' => 'https://api.siliconflow.cn/v1',
            'nvidia' => 'https://integrate.api.nvidia.com/v1',
            'openrouter' => 'https://openrouter.ai/api/v1',
            'bytedance' => 'https://ark.cn-beijing.volces.com/api/v3',
            'minimax' => 'https://api.minimax.chat/v1',
            'baichuan' => 'https://api.baichuan-ai.com/v1',
            'stepfun' => 'https://api.stepfun.com/v1',
            'lingyi' => 'https://api.lingyiwanwu.com/v1',
            'mistral' => 'https://api.mistral.ai/v1',
            'cohere' => 'https://api.cohere.ai/v1',
            'perplexity' => 'https://api.perplexity.ai/v1',
            'ai21' => 'https://api.ai21.com/studio/v1',
            'together' => 'https://api.together.xyz/v1',
            'fireworks' => 'https://api.fireworks.ai/inference/v1',
            'groq' => 'https://api.groq.com/openai/v1',
            'replicate' => 'https://api.replicate.com/v1',
            'anyscale' => 'https://api.endpoints.anyscale.com/v1',
            'azure' => 'https://models.inference.ai.azure.com/v1',
            'ollama' => 'http://localhost:11434/v1',
            'vllm' => 'http://localhost:8000/v1',
            default => 'https://api.openai.com/v1',
        };
    }

    /**
     * 从官方 API 拉取供应商模型列表（CC Switch 风格，避免模型列表过时）。
     * OpenAI 兼容类用 GET /models；Claude 用 GET /models（x-api-key）；Gemini 用 GET /models?key=
     */
    public function listModels(string $provider, string $apiKey, ?string $apiBase = null): array
    {
        if ($apiKey === '') {
            return [];
        }

        $base = rtrim($apiBase ?: $this->defaultBase($provider), '/');

        try {
            if ($provider === 'claude') {
                $response = Http::withHeaders([
                    'x-api-key' => $apiKey,
                    'anthropic-version' => '2023-06-01',
                ])->timeout(15)->get($base . '/models');
                $ids = array_map(static fn ($m) => $m['id'] ?? '', $response->json('data', []));
            } elseif ($provider === 'google') {
                $response = Http::timeout(15)->get($base . '/models?key=' . $apiKey);
                $ids = array_map(
                    static fn ($m) => str_replace('models/', '', (string) ($m['name'] ?? '')),
                    $response->json('models', [])
                );
            } else {
                $response = Http::withHeaders([
                    'Authorization' => 'Bearer ' . $apiKey,
                ])->timeout(15)->get($base . '/models');
                $ids = array_map(static fn ($m) => $m['id'] ?? '', $response->json('data', []));
            }
        } catch (\Throwable) {
            return [];
        }

        return array_values(array_filter($ids, static fn ($id) => $id !== ''));
    }

    private function callOpenaiCompatible(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        $base = $apiBase ?: 'https://api.openai.com/v1';
        $response = Http::withHeaders([
            'Authorization' => 'Bearer ' . $apiKey,
            'Content-Type' => 'application/json',
        ])->timeout(30)->post($base . '/chat/completions', [
            'model' => $model,
            'messages' => [
                ['role' => 'system', 'content' => '你是一个学习助手，请用中文回答。'],
                ['role' => 'user', 'content' => $question],
            ],
            'max_tokens' => $maxTokens,
            'temperature' => 0.7,
        ]);

        if ($response->failed()) {
            return ['answer' => 'AI 服务暂时不可用', 'tokens_used' => 0, 'prompt_tokens' => 0, 'completion_tokens' => 0];
        }

        $usage = $response->json('usage', []);
        $promptTokens = $usage['prompt_tokens'] ?? $usage['total_tokens'] ?? 0;
        $completionTokens = $usage['completion_tokens'] ?? 0;

        return [
            'answer' => $response->json('choices.0.message.content') ?? '抱歉，无法回答。',
            'tokens_used' => $promptTokens + $completionTokens,
            'prompt_tokens' => $promptTokens,
            'completion_tokens' => $completionTokens,
        ];
    }

    private function callClaude(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        $base = $apiBase ?: 'https://api.anthropic.com/v1';
        $response = Http::withHeaders([
            'x-api-key' => $apiKey, 'anthropic-version' => '2023-06-01',
            'Content-Type' => 'application/json',
        ])->timeout(30)->post($base . '/messages', [
            'model' => $model, 'max_tokens' => $maxTokens,
            'messages' => [['role' => 'user', 'content' => $question]],
        ]);

        if ($response->failed()) {
            return ['answer' => 'AI 服务不可用', 'tokens_used' => 0, 'prompt_tokens' => 0, 'completion_tokens' => 0];
        }
        $answer = '';
        foreach ($response->json('content', []) as $block) {
            if (($block['type'] ?? '') === 'text') $answer .= $block['text'] ?? '';
        }
        $promptTokens = $response->json('usage.input_tokens') ?? 0;
        $completionTokens = $response->json('usage.output_tokens') ?? 0;
        return [
            'answer' => $answer ?: '抱歉，无法回答。',
            'tokens_used' => $promptTokens + $completionTokens,
            'prompt_tokens' => $promptTokens,
            'completion_tokens' => $completionTokens,
        ];
    }

    private function callQwen(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        $base = $apiBase ?: 'https://dashscope.aliyuncs.com/api/v1';
        $response = Http::withHeaders([
            'Authorization' => 'Bearer ' . $apiKey, 'Content-Type' => 'application/json',
        ])->timeout(30)->post($base . '/services/aigc/text-generation/generation', [
            'model' => $model,
            'input' => ['messages' => [
                ['role' => 'system', 'content' => '你是一个学习助手。'],
                ['role' => 'user', 'content' => $question],
            ]],
            'parameters' => ['max_tokens' => $maxTokens, 'temperature' => 0.7],
        ]);
        if ($response->failed()) {
            return ['answer' => 'AI 服务不可用', 'tokens_used' => 0, 'prompt_tokens' => 0, 'completion_tokens' => 0];
        }
        $promptTokens = $response->json('usage.input_tokens') ?? 0;
        $completionTokens = $response->json('usage.output_tokens') ?? 0;
        return [
            'answer' => $response->json('output.text') ?? '抱歉，无法回答。',
            'tokens_used' => $promptTokens + $completionTokens,
            'prompt_tokens' => $promptTokens,
            'completion_tokens' => $completionTokens,
        ];
    }

    private function callDeepseek(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.deepseek.com/v1', $maxTokens);
    }

    /**
     * Google Gemini
     */
    private function callGoogle(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        $base = $apiBase ?: 'https://generativelanguage.googleapis.com/v1';
        $response = Http::timeout(30)->post($base . '/models/' . $model . ':generateContent?key=' . $apiKey, [
            'contents' => [['parts' => [['text' => $question]]]],
            'generationConfig' => ['maxOutputTokens' => $maxTokens, 'temperature' => 0.7],
        ]);

        if ($response->failed()) {
            return ['answer' => 'AI 服务不可用', 'tokens_used' => 0, 'prompt_tokens' => 0, 'completion_tokens' => 0];
        }
        $answer = $response->json('candidates.0.content.parts.0.text') ?? '抱歉，无法回答。';
        $promptTokens = $response->json('usageMetadata.promptTokenCount') ?? 0;
        $completionTokens = $response->json('usageMetadata.candidatesTokenCount') ?? 0;
        return [
            'answer' => $answer,
            'tokens_used' => $promptTokens + $completionTokens,
            'prompt_tokens' => $promptTokens,
            'completion_tokens' => $completionTokens,
        ];
    }

    /**
     * Moonshot / Kimi（月之暗面）— OpenAI 兼容
     */
    private function callMoonshot(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.moonshot.cn/v1', $maxTokens);
    }

    /**
     * Grok（xAI）— OpenAI 兼容
     */
    private function callGrok(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.x.ai/v1', $maxTokens);
    }

    /**
     * SiliconFlow（硅基流动）— OpenAI 兼容
     */
    private function callSiliconflow(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.siliconflow.cn/v1', $maxTokens);
    }

    /**
     * NVIDIA NIM — OpenAI 兼容
     */
    private function callNvidia(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://integrate.api.nvidia.com/v1', $maxTokens);
    }

    /**
     * OpenRouter — OpenAI 兼容
     */
    private function callOpenrouter(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://openrouter.ai/api/v1', $maxTokens);
    }

    /**
     * ByteDance Doubao（豆包）— OpenAI 兼容（火山引擎）
     */
    private function callBytedance(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://ark.cn-beijing.volces.com/api/v3', $maxTokens);
    }

    // ===== 新增国内产商（OpenAI 兼容） =====
    private function callMinimax(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.minimax.chat/v1', $maxTokens);
    }
    private function callBaichuan(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.baichuan-ai.com/v1', $maxTokens);
    }
    private function callStepfun(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.stepfun.com/v1', $maxTokens);
    }
    private function callLingyi(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.lingyiwanwu.com/v1', $maxTokens);
    }

    // ===== 国际产商（OpenAI 兼容） =====
    private function callMistral(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.mistral.ai/v1', $maxTokens);
    }
    private function callCohere(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.cohere.ai/v1', $maxTokens);
    }
    private function callPerplexity(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.perplexity.ai/v1', $maxTokens);
    }
    private function callAi21(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.ai21.com/studio/v1', $maxTokens);
    }

    // ===== 聚合平台（OpenAI 兼容） =====
    private function callTogether(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.together.xyz/v1', $maxTokens);
    }
    private function callFireworks(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.fireworks.ai/inference/v1', $maxTokens);
    }
    private function callGroq(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.groq.com/openai/v1', $maxTokens);
    }
    private function callReplicate(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.replicate.com/v1', $maxTokens);
    }
    private function callAnyscale(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://api.endpoints.anyscale.com/v1', $maxTokens);
    }
    private function callAzure(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'https://models.inference.ai.azure.com/v1', $maxTokens);
    }

    // ===== 本地自托管 =====
    private function callOllama(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'http://localhost:11434/v1', $maxTokens);
    }
    private function callVllm(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase ?: 'http://localhost:8000/v1', $maxTokens);
    }

    private function callMcp(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        if (empty($apiBase)) {
            return ['answer' => 'MCP 接口需要配置 API 地址', 'tokens_used' => 0];
        }
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase, $maxTokens);
    }
}
