<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

class HomeworkCollection extends Model
{
    protected $fillable = [
        'class_id', 'teacher_id', 'title',
        'deadline', 'description', 'status',
        'submit_types', 'qr_code_token',
    ];

    protected $casts = [
        'deadline' => 'datetime',
        'submit_types' => 'array',
    ];

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }

    public function teacher(): BelongsTo
    {
        return $this->belongsTo(User::class, 'teacher_id');
    }

    public function submissions(): HasMany
    {
        return $this->hasMany(HomeworkSubmission::class, 'homework_id');
    }

    /**
     * 提交状态
     */
    public static function getStatuses(): array
    {
        return [
            'active'   => '进行中',
            'closed'   => '已截止',
            'archived' => '已归档',
        ];
    }

    /**
     * 提交率
     */
    public function getSubmitRate(): array
    {
        $total = $this->classRoom->students()->where('status', 'active')->count();
        $submitted = $this->submissions()->count();

        return [
            'total'     => $total,
            'submitted' => $submitted,
            'rate'      => $total > 0 ? round($submitted / $total * 100, 1) : 0,
        ];
    }
}

class HomeworkSubmission extends Model
{
    protected $fillable = [
        'homework_id', 'student_id',
        'submitted_at', 'content', 'file_urls',
    ];

    protected $casts = [
        'submitted_at' => 'datetime',
        'file_urls' => 'array',
    ];

    public function homework(): BelongsTo
    {
        return $this->belongsTo(HomeworkCollection::class, 'homework_id');
    }

    public function student(): BelongsTo
    {
        return $this->belongsTo(Student::class);
    }
}
