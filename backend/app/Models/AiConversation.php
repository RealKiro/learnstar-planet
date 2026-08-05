<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class AiConversation extends Model
{
    protected $table = 'ai_conversations';

    protected $fillable = [
        'school_id', 'class_id', 'student_id', 'student_name', 'provider',
        'question', 'answer', 'tokens_used', 'prompt_tokens', 'completion_tokens',
        'cost', 'currency', 'status',
    ];

    protected $casts = [
        'tokens_used' => 'integer',
        'prompt_tokens' => 'integer',
        'completion_tokens' => 'integer',
        'cost' => 'decimal:6',
    ];

    public function school()
    {
        return $this->belongsTo(School::class);
    }

    public function classRoom()
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function student()
    {
        return $this->belongsTo(Student::class);
    }
}
