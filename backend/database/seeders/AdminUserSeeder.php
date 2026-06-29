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
        $adminUsername = env('ADMIN_USERNAME', 'admin');
        $adminPassword = env('ADMIN_PASSWORD', 'admin123456');
        $adminName     = env('ADMIN_NAME', '系统管理员');
        $schoolName    = env('ADMIN_SCHOOL_NAME', '学趣星球默认学校');
        $schoolCode    = env('ADMIN_SCHOOL_CODE', 'learnstar-default');

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

        // 创建超级管理员（如果不存在）
        User::firstOrCreate(
            ['username' => $adminUsername],
            [
                'school_id'       => $school->id,
                'role'            => 'school_admin',
                'password'        => Hash::make($adminPassword),
                'name'            => $adminName,
                'status'          => 'active',
                'password_changed' => false,
            ]
        );

        $this->command?->info("超级管理员已就绪：{$adminUsername}");
        $this->command?->warn('生产环境请务必在 .env 中修改 ADMIN_PASSWORD！');
    }
}
