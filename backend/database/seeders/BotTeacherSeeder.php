<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

/**
 * API 机器人教师账号 —— 供外部系统（QQ/微信机器人、课中等第三方项目）
 * 通过 REST API 对接学生积分管理。拥有本校全部班级的查询与操作权限
 * （权限在 TeacherController::teacherClassIds 中按 settings.is_api_bot 放行）。
 *
 * 配置（.env）：
 *   BOT_ENABLED=false          关闭后不创建/不恢复（已存在的账号会被禁用 status=disabled）
 *   BOT_USERNAME=api-bot       登录账号
 *   BOT_PASSWORD=              密码，以 .env 为唯一真相来源，每次启动同步
 */
class BotTeacherSeeder extends Seeder
{
    public function run(): void
    {
        // 注意：Laravel env() 会把 .env 中的 "false" 转为布尔 false，需兼容布尔与字符串
        $enabled = env('BOT_ENABLED', 'true');
        $isDisabled = $enabled === false || $enabled === null
            || $enabled === 'false' || $enabled === '0';

        if ($isDisabled) {
            // 显式停用：把已有机器人账号置为 disabled（不影响其他账号）
            User::where('role', 'teacher')
                ->get()
                ->filter(fn (User $u) => $u->isApiBot())
                ->each(fn (User $u) => $u->update(['status' => 'disabled']));

            return;
        }

        $username = (string) env('BOT_USERNAME', 'api-bot');
        $password = (string) env('BOT_PASSWORD', 'learnstar-bot-2026');
        $name = (string) env('BOT_NAME', 'API 机器人');
        $school = \App\Models\School::first();

        if (!$school) {
            $this->command?->warn('BotTeacherSeeder：学校不存在（请先运行 AdminUserSeeder），跳过');

            return;
        }

        $user = User::where('username', $username)->first();

        if ($user && !$user->isApiBot()) {
            // 同名账号存在但不是机器人：不劫持，换用 bot_ 前缀账号
            $username = 'bot_' . $username;
            $user = User::where('username', $username)->first();
        }

        if (!$user) {
            $user = User::create([
                'school_id' => $school->id,
                'role' => 'teacher',
                'username' => $username,
                'password' => Hash::make($password),
                'plain_password' => $password,
                'name' => $name,
                'nickname' => $name,
                'status' => 'active',
                'password_changed' => false,
                'settings' => ['is_api_bot' => true],
            ]);
            $this->command?->info("API 机器人账号已创建：{$username}");
        }

        // 每次启动同步：标记、密码（.env 为唯一真相来源）、状态
        $user->settings = array_merge($user->settings ?? [], ['is_api_bot' => true]);
        $user->password = Hash::make($password);
        $user->plain_password = $password;
        $user->status = 'active';
        $user->save();

        $this->command?->info("API 机器人账号已就绪：{$username}（密码已从 .env 同步）");
    }
}
