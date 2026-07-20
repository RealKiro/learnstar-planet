<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

class School extends Model
{
    use SoftDeletes;

    protected $fillable = [
        'name',
        'code',           // 学校唯一代码，用于分配账号前缀
        'address',
        'contact_phone',
        'contact_email',
        'logo_path',
        'settings',        // JSON: 学校级配置（积分规则、宠物类型等）
        'status',          // active / suspended
    ];

    protected $casts = [
        'settings' => 'array',
    ];

    public function classes()
    {
        return $this->hasMany(ClassRoom::class);
    }

    public function admins()
    {
        return $this->hasMany(User::class, 'school_id')->where('role', 'school_admin');
    }

    public function teachers()
    {
        return $this->hasMany(User::class, 'school_id')->where('role', 'teacher');
    }

    public function students()
    {
        return $this->hasManyThrough(Student::class, ClassRoom::class);
    }
}

