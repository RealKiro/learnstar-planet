<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\School;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class DemoTeacherSeeder extends Seeder
{
    public function run(): void
    {
        $school = School::first();
        if (!$school) {
            return;
        }

        // 检查是否已有 Demo 教师账号
        if (User::where('username', 'demo_t1')->exists()) {
            return;
        }

        $names = ['张老师', '李老师', '王老师', '赵老师'];
        foreach ($names as $i => $name) {
            User::create([
                'school_id' => $school->id,
                'role' => 'teacher',
                'username' => 'demo_t' . ($i + 1),
                'password' => Hash::make('demo123'),
                'name' => $name,
                'nickname' => 'Demo教师' . ($i + 1),
                'subject' => ['语文', '数学', '英语', '科学'][$i],
                'status' => 'active',
            ]);
        }

        $this->command?->info('✅ Demo 教师账号已创建（demo_t1~t4 / demo123，可在管理后台删除）');
    }
}
