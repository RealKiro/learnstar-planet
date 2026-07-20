<?php

declare(strict_types=1);

return [
    'provider' => env('AI_PROVIDER', ''),
    'api_key' => env('AI_API_KEY', ''),
    'api_base' => env('AI_API_BASE', ''),
    'model' => env('AI_MODEL', 'deepseek-chat'),
    'max_tokens' => (int) env('AI_MAX_TOKENS', 2000),
    'temperature' => (float) env('AI_TEMPERATURE', 0.7),
];
