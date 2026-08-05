<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use Illuminate\Support\Facades\Cache;

/**
 * 班级码生成服务。
 * 规则（用户确认）：纯 4 位数字，前 2 位年级（补零，三年级→03）、后 2 位班号（补零，2班→02）。
 * 例：三年级2班 → 0302。确定性、不随机、无需刷新。
 */
final class DisplayCodeService
{
    private const GRADE_MAP = [
        '一年级' => '1', '二年级' => '2', '三年级' => '3',
        '四年级' => '4', '五年级' => '5', '六年级' => '6',
        '七年级' => '7', '八年级' => '8', '九年级' => '9',
        '高一' => '10', '高二' => '11', '高三' => '12',
    ];

    /**
     * 年级 → 2 位数字（'三年级'→'03'，'高一'→'10'），无法解析返回 '00'
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

        return $digit === '' ? '00' : str_pad($digit, 2, '0', STR_PAD_LEFT);
    }

    /**
     * 班级名 → 2 位班号（'三年级（2）班'→'02'，'三年级(12)班'→'12'），无法解析返回 '00'
     */
    public static function classNo(string $name): string
    {
        $num = '';
        if (preg_match('/[（(]\s*(\d+)\s*[)）]\s*班/', $name, $m)) {
            $num = $m[1];
        }

        return $num === '' ? '00' : str_pad($num, 2, '0', STR_PAD_LEFT);
    }

    /**
     * 生成班级码（4 位数字）
     */
    public static function generate(ClassRoom $room): string
    {
        return self::gradeDigit($room->grade) . self::classNo($room->name ?? '');
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
            $classDigit = self::classNo($room->name ?? '');

            // 年级或班号无法解析（00 不可能为真实值）→ 跳过并报告
            if ($gradeDigit === '00' || $classDigit === '00') {
                $skipped++;
                continue;
            }

            $code = $gradeDigit . $classDigit;

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
