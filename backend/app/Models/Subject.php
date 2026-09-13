<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * 科目（学校级，全校共享一套）
 */
class Subject extends Model
{
    /** @var list<string> */
    protected $fillable = ['school_id', 'name', 'simplified_name', 'color', 'sort_order'];

    /** @var array<string, string> */
    protected $casts = ['sort_order' => 'integer'];
}
