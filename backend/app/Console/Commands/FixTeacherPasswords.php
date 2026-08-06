<?php

declare(strict_types=1);

namespace App\Console\Commands;

use App\Models\User;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Hash;

/**
 * 修复第三方扫码自动注册教师账号的密码状态。
 *
 * 背景：早期 loginWithThirdParty / loginWithWechatWork 自动注册时未写 plain_password（NULL），
 * 管理员在密码弹窗「查看密码」时 getTeacherPassword 会因 plain_password 为空自动生成随机密码并覆盖，
 * 导致教师用自己以为正确的密码登录失败。
 *
 * 本命令将 plain_password 为空的教师账号统一重置为默认密码 ls123456
 * （password 与 plain_password 保持一致），管理员可在弹窗看到并告知教师。
 */
class FixTeacherPasswords extends Command
{
    protected $signature = 'teacher:fix-passwords {--password= : 指定的默认密码，默认 ls123456} {--dry-run : 只统计不修改}';

    protected $description = '将 plain_password 为空的教师账号重置为默认密码，修复登录失败';

    public function handle(): int
    {
        $newPassword = $this->option('password') ?: 'ls123456';
        $dryRun = (bool) $this->option('dry-run');

        $teachers = User::where('role', 'teacher')
            ->where(function ($q) {
                $q->whereNull('plain_password')->orWhere('plain_password', '');
            })
            ->get();

        $this->info("发现 plain_password 为空的教师账号：{$teachers->count()} 个");

        if ($dryRun) {
            $this->info('dry-run 模式，未做任何修改。');

            return self::SUCCESS;
        }

        $updated = 0;
        foreach ($teachers as $teacher) {
            $teacher->password = Hash::make($newPassword);
            $teacher->plain_password = $newPassword;
            $teacher->password_changed = false;
            $teacher->save();
            $updated++;
        }

        $this->info("已重置 {$updated} 个教师账号的密码为默认密码「{$newPassword}」。");
        $this->info('请告知这些教师用默认密码登录后自行修改；管理员也可在密码弹窗查看/重置。');

        return self::SUCCESS;
    }
}
