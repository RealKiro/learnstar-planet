<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassPeriod;
use App\Models\Subject;
use App\Models\TimetableEntry;
use Illuminate\Support\Facades\DB;

/**
 * 课表读写与导出。
 *
 * - 科目 / 节次为学校级共享资源；排课按班级存
 * - API 边界上用「科目名」交互（前端新建科目时无 id），内部落库时才解析为 subject_id
 * - toCses() 产出 CSES（The Course Schedule Exchange Schema）YAML，
 *   ClassIsland 可直接「从 CSES 导入」，实现课表对接
 */
class TimetableService
{
    /** 星期中文名（CSES 的 enable_day 为 1-7，即周一到周日） */
    private const WEEKDAY_LABELS = [
        1 => '星期一', 2 => '星期二', 3 => '星期三', 4 => '星期四',
        5 => '星期五', 6 => '星期六', 7 => '星期日',
    ];

    private const WEEK_TYPES = ['all', 'odd', 'even'];

    /** 初始化数据：科目名列表 + 节次 + 该班排课（排课以科目名返回） */
    public function bootstrap(int $classId, int $schoolId): array
    {
        $subjects = Subject::where('school_id', $schoolId)
            ->orderBy('sort_order')->orderBy('id')
            ->get(['id', 'name', 'simplified_name', 'color']);

        $nameById = $subjects->pluck('name', 'id');

        $entries = TimetableEntry::where('class_id', $classId)
            ->orderBy('weekday')->orderBy('period_index')
            ->get(['weekday', 'period_index', 'week_type', 'subject_id', 'teacher_name', 'room'])
            ->map(fn (TimetableEntry $e) => [
                'weekday' => (int) $e->weekday,
                'period_index' => (int) $e->period_index,
                'week_type' => (string) ($e->week_type ?: 'all'),
                'subject_name' => $e->subject_id ? ($nameById[$e->subject_id] ?? null) : null,
                'teacher_name' => $e->teacher_name,
                'room' => $e->room,
            ])
            ->values();

        return [
            'subjects' => $subjects->map(fn (Subject $s) => [
                'name' => $s->name,
                'simplified_name' => $s->simplified_name,
                'color' => $s->color,
            ])->values(),
            'periods' => ClassPeriod::where('school_id', $schoolId)
                ->orderBy('period_index')
                ->get(['period_index', 'name', 'start_time', 'end_time']),
            'entries' => $entries,
        ];
    }

    /**
     * 整体保存：科目 upsert → 节次 upsert → 该班排课先清后插。
     *
     * @param  array{subjects?: array<int, array<string, mixed>>, periods?: array<int, array<string, mixed>>, entries?: array<int, array<string, mixed>>}  $payload
     */
    public function save(int $classId, int $schoolId, array $payload): void
    {
        DB::transaction(function () use ($classId, $schoolId, $payload): void {
            foreach (array_values($payload['subjects'] ?? []) as $i => $s) {
                $name = trim((string) ($s['name'] ?? ''));
                if ($name === '') {
                    continue;
                }
                Subject::updateOrCreate(
                    ['school_id' => $schoolId, 'name' => $name],
                    [
                        'simplified_name' => $s['simplified_name'] ?? null,
                        'color' => $s['color'] ?? null,
                        'sort_order' => $i,
                    ],
                );
            }

            foreach (array_values($payload['periods'] ?? []) as $p) {
                $idx = (int) ($p['period_index'] ?? 0);
                if ($idx < 1 || empty($p['start_time']) || empty($p['end_time'])) {
                    continue;
                }
                ClassPeriod::updateOrCreate(
                    ['school_id' => $schoolId, 'period_index' => $idx],
                    [
                        'name' => $p['name'] ?? ('第' . $idx . '节'),
                        'start_time' => (string) $p['start_time'],
                        'end_time' => (string) $p['end_time'],
                    ],
                );
            }

            $idByName = Subject::where('school_id', $schoolId)->pluck('id', 'name');

            TimetableEntry::where('class_id', $classId)->delete();

            foreach (array_values($payload['entries'] ?? []) as $e) {
                $name = trim((string) ($e['subject_name'] ?? ''));
                $subjectId = $idByName[$name] ?? null;
                $weekday = (int) ($e['weekday'] ?? 0);
                $periodIndex = (int) ($e['period_index'] ?? 0);

                if (!$subjectId || $weekday < 1 || $weekday > 7 || $periodIndex < 1) {
                    continue;
                }

                $weekType = (string) ($e['week_type'] ?? 'all');

                TimetableEntry::create([
                    'class_id' => $classId,
                    'weekday' => $weekday,
                    'period_index' => $periodIndex,
                    'week_type' => in_array($weekType, self::WEEK_TYPES, true) ? $weekType : 'all',
                    'subject_id' => $subjectId,
                    'teacher_name' => $e['teacher_name'] ?? null,
                    'room' => $e['room'] ?? null,
                ]);
            }
        });
    }

    /**
     * 导出为 CSES YAML（ClassIsland 可直接导入）。
     *
     * 结构：version / subjects[] / schedules[]（每个「星期 + 单双周」一个 schedule，
     * 其 classes[] 内每节带 subject + start_time + end_time）。
     */
    public function toCses(int $classId, int $schoolId): string
    {
        $subjects = Subject::where('school_id', $schoolId)->orderBy('sort_order')->orderBy('id')->get();
        $periods = ClassPeriod::where('school_id', $schoolId)->get()->keyBy('period_index');
        $entries = TimetableEntry::where('class_id', $classId)
            ->orderBy('weekday')->orderBy('period_index')->get();

        $lines = [];
        $lines[] = 'version: 1';

        $lines[] = 'subjects:';
        if ($subjects->isEmpty()) {
            $lines[] = '  []';
        }
        foreach ($subjects as $s) {
            $lines[] = '  - name: ' . $this->yamlString((string) $s->name);
            if (!empty($s->simplified_name)) {
                $lines[] = '    simplified_name: ' . $this->yamlString((string) $s->simplified_name);
            }
        }

        // 按 星期 → 单双周 分组
        $grouped = [];
        foreach ($entries as $e) {
            if (!$e->subject_id) {
                continue;
            }
            $grouped[(int) $e->weekday][(string) ($e->week_type ?: 'all')][] = $e;
        }
        ksort($grouped);

        $lines[] = 'schedules:';
        if (empty($grouped)) {
            $lines[] = '  []';
        }
        foreach ($grouped as $weekday => $byWeek) {
            foreach (self::WEEK_TYPES as $weekType) {
                if (empty($byWeek[$weekType])) {
                    continue;
                }
                $label = (self::WEEKDAY_LABELS[$weekday] ?? ('第' . $weekday . '天'))
                    . ($weekType === 'odd' ? '·单周' : ($weekType === 'even' ? '·双周' : ''));
                $lines[] = '  - name: ' . $this->yamlString($label);
                $lines[] = '    enable_day: ' . $weekday;
                $lines[] = '    weeks: ' . $weekType;
                $lines[] = '    classes:';
                foreach ($byWeek[$weekType] as $e) {
                    $subjectName = optional($subjects->firstWhere('id', $e->subject_id))->name;
                    if (!$subjectName) {
                        continue;
                    }
                    $lines[] = '      - subject: ' . $this->yamlString((string) $subjectName);
                    $period = $periods->get($e->period_index);
                    if ($period) {
                        $lines[] = '        start_time: "' . $period->start_time . '"';
                        $lines[] = '        end_time: "' . $period->end_time . '"';
                    }
                    if (!empty($e->teacher_name)) {
                        $lines[] = '        teacher: ' . $this->yamlString((string) $e->teacher_name);
                    }
                    if (!empty($e->room)) {
                        $lines[] = '        room: ' . $this->yamlString((string) $e->room);
                    }
                }
            }
        }

        return implode("\n", $lines) . "\n";
    }

    /** YAML 双引号标量（转义反斜杠与双引号） */
    private function yamlString(string $value): string
    {
        return '"' . str_replace(['\\', '"'], ['\\\\', '\\"'], $value) . '"';
    }
}
