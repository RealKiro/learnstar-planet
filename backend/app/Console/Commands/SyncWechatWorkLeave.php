<?php
declare(strict_types=1);
namespace App\Console\Commands;
use App\Services\WechatWorkAttendanceService;
use Illuminate\Console\Command;
class SyncWechatWorkLeave extends Command
{
    protected $signature = 'attendance:sync-wechat-leave {--date=} {--school=}';
    protected $description = '同步企业微信请假到考勤';
    public function handle(WechatWorkAttendanceService $service): int {
        $date = $this->option('date') ?: now()->toDateString(); $this->info("同步: {$date}");
        if ($sid = $this->option('school')) { $r=$service->syncForSchool((int)$sid,$date); $this->info("同步{$r['synced']}跳过{$r['skipped']}"); return self::SUCCESS; }
        $r=$service->syncAll($date); $this->info("同步{$r['synced']}跳过{$r['skipped']}".($r['failed']>0?", {$r['failed']}校失败":"")); return self::SUCCESS;
    }
}
