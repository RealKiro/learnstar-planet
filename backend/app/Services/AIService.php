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

    private function callMcp(string $apiKey, string $model, string $question, ?string $apiBase, int $maxTokens): array
    {
        if (empty($apiBase)) {
            return ['answer' => 'MCP 接口需要配置 API 地址', 'tokens_used' => 0];
        }
        return $this->callOpenaiCompatible($apiKey, $model, $question, $apiBase, $maxTokens);
    }
}
