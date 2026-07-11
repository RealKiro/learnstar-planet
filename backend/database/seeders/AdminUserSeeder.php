<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\School;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class AdminUserSeeder extends Seeder
{
    public function run(): void
    {
        // 读取环境变量，提供安全默认值
        // 兼容历史字段 ADMIN_* 与新统一字段 INITIAL_*（推荐用 INITIAL_*）
        $adminUsername = env('INITIAL_ADMIN_USERNAME', env('ADMIN_USERNAME', 'admin'));
        $adminPassword = env('INITIAL_ADMIN_PASSWORD', env('ADMIN_PASSWORD', 'ChangeThisImmediately'));
        $adminName     = env('INITIAL_ADMIN_NAME', env('ADMIN_NAME', '李校长'));
        $schoolName    = env('INITIAL_SCHOOL_NAME', env('ADMIN_SCHOOL_NAME', '阳光小学'));
        $schoolCode    = env('INITIAL_SCHOOL_CODE', env('ADMIN_SCHOOL_CODE', 'learnstar-default'));

        // 创建默认学校（如果不存在）
        $school = School::firstOrCreate(
            ['code' => $schoolCode],
            [
                'name'          => $schoolName,
                'status'        => 'active',
                'address'       => '',
                'contact_phone' => '',
                'contact_email' => '',
                'settings'      => [],
            ]
        );

        // 创建或更新超级管理员
        // 每次启动都从 .env 同步密码，确保 ADMIN_PASSWORD 是唯一真相来源
        $user = User::where('username', $adminUsername)->first();

        if (!$user) {
            // 首次创建
            User::create([
                'username'         => $adminUsername,
                'school_id'        => $school->id,
                'role'             => 'school_admin',
                'password'         => Hash::make($adminPassword),
                'name'             => $adminName,
                'status'           => 'active',
                'password_changed' => false,
            ]);
        } else {
            // 已存在：仅同步密码，不动其他字段（name/状态等可通过后台修改）
            $user->update(['password' => Hash::make($adminPassword)]);
        }

        $this->command?->info("超级管理员已就绪：{$adminUsername}（密码已从 .env 同步）");
    }
}
