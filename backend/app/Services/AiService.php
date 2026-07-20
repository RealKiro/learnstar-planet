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
            return ['answer' => 'AI 服务暂时不可用', 'tokens_used' => 0];
        }

        return [
            'answer' => $response->json('choices.0.message.content') ?? '抱歉，无法回答。',
            'tokens_used' => $response->json('usage.total_tokens') ?? 0,
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
            return ['answer' => 'AI 服务不可用', 'tokens_used' => 0];
        }
        $answer = '';
        foreach ($response->json('content', []) as $block) {
            if (($block['type'] ?? '') === 'text') $answer .= $block['text'] ?? '';
        }
        return [
            'answer' => $answer ?: '抱歉，无法回答。',
            'tokens_used' => ($response->json('usage.input_tokens') ?? 0) + ($response->json('usage.output_tokens') ?? 0),
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
            return ['answer' => 'AI 服务不可用', 'tokens_used' => 0];
        }
        return [
            'answer' => $response->json('output.text') ?? '抱歉，无法回答。',
            'tokens_used' => $response->json('usage.total_tokens') ?? 0,
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
            return ['answer' => 'AI 服务不可用', 'tokens_used' => 0];
        }
        $answer = $response->json('candidates.0.content.parts.0.text') ?? '抱歉，无法回答。';
        return [
            'answer' => $answer,
            'tokens_used' => ($response->json('usageMetadata.promptTokenCount') ?? 0) + ($response->json('usageMetadata.candidatesTokenCount') ?? 0),
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
