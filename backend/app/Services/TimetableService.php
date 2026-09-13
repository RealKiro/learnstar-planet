<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassPeriod;
use App\Models\Subject;
use App\Models\TimetableChangeRequest;
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

    /** 大屏只读数据：节次 + 该班排课（科目名形式）+ 科目配色。week_type 原样返回，单双周归属由展示层标注 */
    public function forDisplay(int $classId, int $schoolId): array
    {
        $subjects = Subject::where('school_id', $schoolId)
            ->orderBy('sort_order')->orderBy('id')
            ->get(['name', 'simplified_name', 'color']);

        $nameById = Subject::where('school_id', $schoolId)->pluck('name', 'id');

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
            ->filter(fn (array $e) => $e['subject_name'] !== null)
            ->values();

        return [
            'subjects' => $subjects->values(),
            'periods' => ClassPeriod::where('school_id', $schoolId)
                ->orderBy('period_index')
                ->get(['period_index', 'name', 'start_time', 'end_time'])
                ->values(),
            'entries' => $entries,
        ];
    }

    /**
     * 提交课表修改申请（不直接生效，待管理员审核）。
     *
     * @param  array{subjects?: array<int, array<string, mixed>>, periods?: array<int, array<string, mixed>>, entries?: array<int, array<string, mixed>>}  $payload
     */
    public function submitChange(int $classId, int $schoolId, int $teacherId, array $payload): TimetableChangeRequest
    {
        $entryCount = count(array_values($payload['entries'] ?? []));

        return TimetableChangeRequest::create([
            'school_id' => $schoolId,
            'class_id' => $classId,
            'requested_by' => $teacherId,
            'payload' => $payload,
            'entry_count' => $entryCount,
            'status' => TimetableChangeRequest::STATUS_PENDING,
        ]);
    }

    /** 审核通过并应用申请的课表快照（幂等：仅 pending 可通过） */
    public function approveChange(int $requestId, int $reviewerId, ?string $note = null): bool
    {
        return DB::transaction(function () use ($requestId, $reviewerId, $note): bool {
            $request = TimetableChangeRequest::whereKey($requestId)->lockForUpdate()->first();

            if (!$request || $request->status !== TimetableChangeRequest::STATUS_PENDING) {
                return false;
            }

            // 应用快照（save 内部也是事务，嵌套事务由 Laravel savepoint 保证）
            $this->save((int) $request->class_id, (int) $request->school_id, $request->payload ?? []);

            $request->update([
                'status' => TimetableChangeRequest::STATUS_APPROVED,
                'reviewed_by' => $reviewerId,
                'review_note' => $note,
                'reviewed_at' => now(),
            ]);

            // 该班其余 pending 申请自动作废（课表已被本次覆盖）
            TimetableChangeRequest::where('class_id', $request->class_id)
                ->where('status', TimetableChangeRequest::STATUS_PENDING)
                ->update(['status' => TimetableChangeRequest::STATUS_REJECTED, 'review_note' => '已由更新的申请取代', 'reviewed_at' => now()]);

            return true;
        });
    }

    /** 驳回申请（幂等：仅 pending 可驳回） */
    public function rejectChange(int $requestId, int $reviewerId, ?string $note = null): bool
    {
        $request = TimetableChangeRequest::find($requestId);

        if (!$request || $request->status !== TimetableChangeRequest::STATUS_PENDING) {
            return false;
        }

        $request->update([
            'status' => TimetableChangeRequest::STATUS_REJECTED,
            'reviewed_by' => $reviewerId,
            'review_note' => $note,
            'reviewed_at' => now(),
        ]);

        return true;
    }

    /** 某班申请历史（最新在前） */
    public function listChangesForClass(int $classId): array
    {
        return TimetableChangeRequest::where('class_id', $classId)
            ->with(['requester:id,name', 'reviewer:id,name'])
            ->orderByDesc('id')
            ->limit(20)
            ->get()
            ->map(fn (TimetableChangeRequest $r) => [
                'id' => $r->id,
                'class_id' => (int) $r->class_id,
                'entry_count' => (int) $r->entry_count,
                'status' => (string) $r->status,
                'review_note' => $r->review_note,
                'created_at' => $r->created_at?->toIso8601String(),
                'reviewed_at' => $r->reviewed_at?->toIso8601String(),
                'requester_name' => $r->requester?->name,
                'reviewer_name' => $r->reviewer?->name,
            ])
            ->values()
            ->all();
    }

    /** 管理员待办列表（全部状态，可按状态过滤，最新在前） */
    public function listChangesForSchool(int $schoolId, ?string $status = null): array
    {
        return TimetableChangeRequest::where('school_id', $schoolId)
            ->when($status !== null && $status !== '', fn ($q) => $q->where('status', $status))
            ->with(['requester:id,name', 'reviewer:id,name', 'classRoom:id,name,grade'])
            ->orderByDesc('id')
            ->limit(100)
            ->get()
            ->map(fn (TimetableChangeRequest $r) => [
                'id' => $r->id,
                'class_id' => (int) $r->class_id,
                'class_name' => $r->classRoom?->name,
                'grade' => $r->classRoom?->grade,
                'entry_count' => (int) $r->entry_count,
                'status' => (string) $r->status,
                'review_note' => $r->review_note,
                'created_at' => $r->created_at?->toIso8601String(),
                'reviewed_at' => $r->reviewed_at?->toIso8601String(),
                'requester_name' => $r->requester?->name,
                'reviewer_name' => $r->reviewer?->name,
            ])
            ->values()
            ->all();
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
