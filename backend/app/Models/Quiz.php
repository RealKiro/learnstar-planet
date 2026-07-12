<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

class Quiz extends Model
{
    protected $fillable = [
        'class_id', 'teacher_id', 'title',
        'question_bank_id', 'time_limit',  // 限时（分钟），0=不限时
        'auto_grade', 'realtime_stats', 'anonymous',
        'status', 'started_at', 'ended_at',
    ];

    protected $casts = [
        'auto_grade'      => 'boolean',
        'realtime_stats' => 'boolean',
        'anonymous'       => 'boolean',
        'time_limit'      => 'integer',
        'started_at'      => 'datetime',
        'ended_at'        => 'datetime',
    ];

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    public function questionBank(): BelongsTo
    {
        return $this->belongsTo(QuestionBank::class);
    }

    public function submissions(): HasMany
    {
        return $this->hasMany(QuizSubmission::class);
    }

    public static function getStatuses(): array
    {
        return [
            'draft'    => '草稿',
            'active'   => '进行中',
            'grading'  => '批改中',
            'finished' => '已结束',
        ];
    }
}

class QuizSubmission extends Model
{
    protected $fillable = [
        'quiz_id', 'student_id', 'answers', 'score', 'submitted_at',
    ];

    protected $casts = [
        'answers'      => 'array',
        'submitted_at' => 'datetime',
    ];

    public function quiz(): BelongsTo
    {
        return $this->belongsTo(Quiz::class);
    }

    public function student(): BelongsTo
    {
        return $this->belongsTo(Student::class);
    }
}

class QuestionBank extends Model
{
    protected $fillable = [
        'school_id', 'teacher_id', 'title',
        'subject', 'question_count', 'used_count', 'is_public',
    ];

    public function questions(): HasMany
    {
        return $this->hasMany(Question::class);
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    public static function getSubjects(): array
    {
        return ['语文', '数学', '英语', '物理', '化学', '生物', '历史', '地理', '政治'];
    }
}

class Question extends Model
{
    protected $fillable = [
        'question_bank_id', 'type', 'content', 'options', 'answer', 'score',
    ];

    protected $casts = [
        'options' => 'array',
    ];

    public function bank(): BelongsTo
    {
        return $this->belongsTo(QuestionBank::class, 'question_bank_id');
    }

    public static function getTypes(): array
    {
        return [
            'single'    => '单选题',
            'multiple'  => '多选题',
            'fill'      => '填空题',
            'judge'     => '判断题',
            'short'     => '简答题',
        ];
    }
}
