<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class DisplayLoginLog extends Model
{
    protected $table = 'display_login_logs';

    protected $fillable = [
        'class_id',
        'class_code',
        'ip_address',
        'user_agent',
    ];

    public function classRoom(): BelongsTo
    {
        return $this->belongsTo(ClassRoom::class, 'class_id');
    }
}
