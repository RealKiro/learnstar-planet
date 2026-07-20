<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class AiSetting extends Model
{
    protected $table = 'ai_settings';

    protected $fillable = [
        'school_id', 'enabled', 'provider', 'api_key', 'api_base',
        'model', 'max_tokens', 'tokens_used', 'tokens_limit',
    ];

    protected $casts = [
        'enabled' => 'boolean',
        'max_tokens' => 'integer',
        'tokens_used' => 'integer',
        'tokens_limit' => 'integer',
    ];

    public function school()
    {
        return $this->belongsTo(School::class);
    }
}
