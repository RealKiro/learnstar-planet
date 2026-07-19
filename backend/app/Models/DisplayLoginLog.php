<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class DisplayLoginLog extends Model
{
    protected $table = 'display_login_logs';

    protected $fillable = [
        'class_id',
        'class_code',
        'ip_address',
        'user_agent',
    ];
}
