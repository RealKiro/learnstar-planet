<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use Illuminate\Support\Facades\Cache;

/**
 * 班级码生成服务。
 * 规则（用户确认）：`LS` + 年级数字 + 班号数字，不补零。
 * 例：一年级1班 → LS11，二年级3班 → LS23。
 * 确定性、不随机、无需刷新。
 */
final class DisplayCodeService
{
    private const GRADE_MAP = [
        '一年级' => '1', '二年级' => '2', '三年级' => '3',
        '四年级' => '4', '五年级' => '5', '六年级' => '6',
        '七年级' => '7', '八年级' => '8', '九年级' => '9',
        '高一' => '10', '高二' => '11', '高三' => '12',
    ];

    private const CHINESE_NUMBERS = [
        '一' => '1', '二' => '2', '三' => '3', '四' => '4', '五' => '5',
        '六' => '6', '七' => '7', '八' => '8', '九' => '9', '十' => '10',
    ];

    /**
     * 年级 → 数字（'三年级'→'3'，不补零）。
     * 仅接受 1-9 年级（用户规则：暂不考虑 10+），无法解析或超出范围返回 ''
     */
    public static function gradeDigit(?string $grade): string
    {
        $digit = '';
        if (is_string($grade) && $grade !== '') {
            foreach (self::GRADE_MAP as $cn => $d) {
                if (str_contains($grade, $cn)) {
                    $digit = $d;
                    break;
                }
            }
            if ($digit === '' && preg_match('/(\d+)/', $grade, $m)) {
                $digit = $m[1];
            }
        }

        return (ctype_digit($digit) && (int) $digit >= 1 && (int) $digit <= 9) ? $digit : '';
    }

    /**
     * 班级名 → 班号数字（不补零），支持多种格式。
     * 仅接受 1-9 班号（用户规则：暂不考虑 10+），无法解析或超出范围返回 ''
     */
    public static function classNo(string $name): string
    {
        $num = '';
        // 1) 标准：X年级（N）班 / X年级(N)班
        if (preg_match('/[（(]\s*(\d+)\s*[)）]\s*班/', $name, $m)) {
            $num = $m[1];
        }
        // 2) 无括号数字：一年级1班 / 一年级 1 班
        elseif (preg_match('/(\d+)\s*班/', $name, $m)) {
            $num = $m[1];
        }
        // 3) 中文序数：一年级第一班
        elseif (preg_match('/第([一二三四五六七八九十]+)班/', $name, $m)) {
            $num = self::chineseNumberToArabic($m[1]);
        }
        // 4) 中文数字：一年级一班
        elseif (preg_match('/([一二三四五六七八九十]+)班/', $name, $m)) {
            $num = self::chineseNumberToArabic($m[1]);
        }

        return (ctype_digit($num) && (int) $num >= 1 && (int) $num <= 9) ? $num : '';
    }

    /**
     * 生成班级码（LS + 年级 + 班号，如 LS11）。
     * 年级无法解析返回 ''；班号无法解析时用同年级序号兜底。
     */
    public static function generate(ClassRoom $room): string
    {
        $grade = self::gradeDigit($room->grade);
        if ($grade === '') {
            return '';
        }

        $class = self::classNo($room->name ?? '');
        if ($class === '') {
            $class = self::fallbackClassNo($room);
        }

        return 'LS' . $grade . $class;
    }

    /**
     * 班号兜底：班级名无法解析班号时，用同年级内按 ID 排序的序号
     */
    private static function fallbackClassNo(ClassRoom $room): string
    {
        $ids = ClassRoom::where('school_id', $room->school_id)
            ->where('grade', $room->grade)
            ->orderBy('id')
            ->pluck('id');
        $index = $ids->search($room->id);

        return (string) ($index === false ? 1 : $index + 1);
    }

    private static function chineseNumberToArabic(string $cn): string
    {
        if ($cn === '十') {
            return '10';
        }
        if (str_contains($cn, '十')) {
            [$t, $o] = explode('十', $cn);
            $tens = $t === '' ? 1 : (int) (self::CHINESE_NUMBERS[$t] ?? 0);
            $ones = $o === '' ? 0 : (int) (self::CHINESE_NUMBERS[$o] ?? 0);

            return (string) ($tens * 10 + $ones);
        }

        return self::CHINESE_NUMBERS[$cn] ?? '0';
    }

    /**
     * 批量重生成所有班级码（确定性：同一班级每次结果一致）。
     * 旧码缓存映射一并清理，新码写入缓存。
     *
     * @return array{regenerated: int, skipped: int, conflicts: array<int, string>}
     */
    public static function regenerateAll(?int $schoolId = null): array
    {
        $query = ClassRoom::query();
        if ($schoolId !== null) {
            $query->where('school_id', $schoolId);
        }

        $regenerated = 0;
        $skipped = 0;
        $conflicts = [];
        /** @var array<int, array<string, int>> $used 校维度的已用码 */
        $used = [];

        foreach ($query->get() as $room) {
            $gradeDigit = self::gradeDigit($room->grade);
            if ($gradeDigit === '') {
                $skipped++;
                continue;
            }

            $classDigit = self::classNo($room->name ?? '');
            if ($classDigit === '') {
                $classDigit = self::fallbackClassNo($room);
            }

            $code = 'LS' . $gradeDigit . $classDigit;

            // 校内重复码检测（数据异常：同年级同班号存在多条）
            $schoolUsed = $used[$room->school_id] ?? [];
            if (isset($schoolUsed[$code]) && $schoolUsed[$code] !== $room->id) {
                $conflicts[$room->id] = $code;
                $skipped++;
                continue;
            }
            $used[$room->school_id][$code] = $room->id;

            if (!empty($room->display_code)) {
                Cache::forget(DisplayEventService::codeCacheKey($room->display_code));
            }

            $room->display_code = $code;
            $room->display_code_updated_at = now();
            $room->save();

            Cache::put(
                DisplayEventService::codeCacheKey($code),
                $room->id,
                now()->addDays(30),
            );

            $regenerated++;
        }

        return [
            'regenerated' => $regenerated,
            'skipped' => $skipped,
            'conflicts' => $conflicts,
        ];
    }
}
