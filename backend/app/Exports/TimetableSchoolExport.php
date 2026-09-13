<?php

declare(strict_types=1);

namespace App\Exports;

use Maatwebsite\Excel\Concerns\WithMultipleSheets;

/**
 * 全校课表导出：每个班级一个工作表（sheet 顺序即传入顺序）。
 *
 * @param  array<int, TimetableExport>  $sheets
 */
class TimetableSchoolExport implements WithMultipleSheets
{
    public function __construct(
        private readonly array $sheets,
    ) {
    }

    public function sheets(): array
    {
        return $this->sheets;
    }
}
