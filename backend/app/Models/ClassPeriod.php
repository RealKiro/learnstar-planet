<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * 节次（学校级作息时间：第 N 节 + 起止时间）
 */
class ClassPeriod extends Model
{
    /** @var list<string> */
    protected $fillable = ['school_id', 'period_index', 'name', 'start_time', 'end_time'];

    /** @var array<string, string> */
    protected $casts = ['period_index' => 'integer'];
}
