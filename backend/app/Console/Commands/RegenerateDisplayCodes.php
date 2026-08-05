<?php

declare(strict_types=1);

namespace App\Console\Commands;

use App\Services\DisplayCodeService;
use Illuminate\Console\Command;

/**
 * 批量重生成班级码（LS + 年级 + 班号，如一年级1班 → LS11）
 */
class RegenerateDisplayCodes extends Command
{
    protected $signature = 'class:regenerate-display-codes {--school= : 只重生成指定学校}';

    protected $description = '批量重生成所有班级码（LS + 年级 + 班号）';

    public function handle(): int
    {
        $schoolId = $this->option('school') !== null ? (int) $this->option('school') : null;

        $result = DisplayCodeService::regenerateAll($schoolId);

        $this->info(sprintf(
            '班级码重生成完成：成功 %d，跳过 %d，校内冲突 %d',
            $result['regenerated'],
            $result['skipped'],
            count($result['conflicts']),
        ));

        foreach ($result['conflicts'] as $classId => $code) {
            $this->warn("班级 #{$classId} 与同校其他班级冲突（{$code}），已跳过，请人工修正班级名/年级后重跑");
        }

        if ($result['skipped'] > 0) {
            $this->warn('部分班级因年级或班号无法解析被跳过，请检查这些班级的「年级」与「班级名」格式（应为 X年级（N）班）。');
        }

        return self::SUCCESS;
    }
}
