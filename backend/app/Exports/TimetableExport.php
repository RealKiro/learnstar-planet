<?php

declare(strict_types=1);

namespace App\Exports;

use Maatwebsite\Excel\Concerns\FromArray;
use Maatwebsite\Excel\Concerns\WithTitle;
use Maatwebsite\Excel\Concerns\WithStyles;
use PhpOffice\PhpSpreadsheet\Style\Alignment;
use PhpOffice\PhpSpreadsheet\Worksheet\Worksheet;

/**
 * 课表网格 Excel 导出。
 *
 * 单班导出 = 单工作表；全校导出时每班一个工作表（同名班级自动加序号去重）。
 * 行结构由 TimetableService::excelGrid() 生成：标题 / 空行 / 表头 / 节次行。
 */
class TimetableExport implements FromArray, WithTitle, WithStyles
{
    /** 已用过的 sheet 名（全校导出按顺序去重） */
    public static array $usedTitles = [];

    public function __construct(
        private readonly string $sheetTitle,
        private readonly array $rows,
    ) {
    }

    public function array(): array
    {
        return $this->rows;
    }

    public function title(): string
    {
        // Excel 工作表名 ≤31 字符且不可重复
        $base = mb_substr($this->sheetTitle, 0, 28) ?: '班级';
        $title = $base;
        $n = 2;
        while (in_array($title, self::$usedTitles, true)) {
            $title = $base . '-' . $n++;
        }
        self::$usedTitles[] = $title;

        return $title;
    }

    public function styles(Worksheet $sheet)
    {
        $lastCol = $sheet->getHighestColumn();
        $lastRow = $sheet->getHighestRow();

        // 全表：自动换行 + 垂直居中
        $sheet->getStyle('A1:' . $lastCol . $lastRow)->getAlignment()
            ->setWrapText(true)
            ->setVertical(Alignment::VERTICAL_CENTER);

        // 标题行加粗放大，表头行（第 3 行）加粗
        $sheet->getStyle('A1')->getFont()->setBold(true)->setSize(14);
        $sheet->getStyle('A3:' . $lastCol . '3')->getFont()->setBold(true);

        $sheet->getColumnDimension('A')->setWidth(10);
        $sheet->getColumnDimension('B')->setWidth(13);
    }
}
