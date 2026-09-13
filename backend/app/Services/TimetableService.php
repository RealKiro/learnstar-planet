<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassPeriod;
use App\Models\ClassRoom;
use App\Models\Subject;
use App\Models\TimetableChangeRequest;
use App\Models\TimetableEntry;
use App\Models\TimetableTeacherAssignment;
use App\Models\TimetableTeacherUnavailability;
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
     * 管理员直接保存课表（即时生效，不走审核流）。
     *
     * 保存后该班所有待审申请自动作废，防止管理员改完后旧的申请快照
     * 再被批准而覆盖管理员的直接修改。
     *
     * @param  array{subjects?: array<int, array<string, mixed>>, periods?: array<int, array<string, mixed>>, entries?: array<int, array<string, mixed>>}  $payload
     */
    public function adminSave(int $classId, int $schoolId, array $payload): void
    {
        $this->save($classId, $schoolId, $payload);

        TimetableChangeRequest::where('class_id', $classId)
            ->where('status', TimetableChangeRequest::STATUS_PENDING)
            ->update([
                'status' => TimetableChangeRequest::STATUS_REJECTED,
                'review_note' => '管理员已直接修改课表，本申请自动作废',
                'reviewed_at' => now(),
            ]);
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

    // ============================================================
    // 批量导入（CSV）：一次导入全校各班课表
    // ============================================================

    /**
     * 从 CSV 文本批量导入课表。
     *
     * 模板列：年级,班级,星期,第几节,开始时间,结束时间,科目,周次,教师,教室
     * - 按「年级 + 班级」定位 ClassRoom，找不到的行报错跳过（不自动建班）
     * - 科目 / 节次为学校级 upsert（同节次时间冲突时取首次出现并告警）
     * - 某班出现在 CSV 中即整体覆盖该班课表（先清后插），不影响未出现的班级
     * - dryRun=true 只做解析校验不落库
     *
     * @return array{total_rows: int, classes: array<int, array<string, mixed>>, errors: array<int, string>, dry_run: bool, imported: bool}
     */
    public function importFromCsv(string $content, int $schoolId, bool $dryRun): array
    {
        [$header, $lines] = $this->parseCsv($content);

        $col = [];
        $aliases = [
            'grade' => ['年级', 'grade'],
            'class' => ['班级', '班级名称', 'class', 'class_name'],
            'weekday' => ['星期', 'weekday'],
            'period' => ['第几节', '节次', 'period', 'period_index'],
            'start' => ['开始时间', 'start', 'start_time'],
            'end' => ['结束时间', 'end', 'end_time'],
            'subject' => ['科目', 'subject'],
            'week_type' => ['周次', 'week_type'],
            'teacher' => ['教师', '老师', 'teacher'],
            'room' => ['教室', 'room'],
        ];
        foreach ($aliases as $key => $names) {
            foreach ($header as $i => $h) {
                if (in_array($h, $names, true)) {
                    $col[$key] = $i;
                    break;
                }
            }
        }

        $errors = [];
        if (!isset($col['grade'], $col['class'], $col['weekday'], $col['period'], $col['subject'])) {
            $errors[] = 'CSV 缺少必需列（年级 / 班级 / 星期 / 第几节 / 科目），请使用模板';
        }

        $classes = [];  // key: grade|name => ['grade','name','class_id'=>?int,'entries'=>[],'subjects'=>[]]
        $periods = [];  // period_index => ['start_time','end_time']
        $totalRows = 0;

        if (isset($col['grade'], $col['class'], $col['weekday'], $col['period'], $col['subject'])) {
            foreach ($lines as $lineNo => $line) {
                $cells = str_getcsv($line);
                $get = fn (string $k): string => isset($col[$k]) ? trim((string) ($cells[$col[$k]] ?? '')) : '';
                $rowNo = $lineNo + 2; // +1 表头，+1 从 1 计数

                $grade = $get('grade');
                $className = $get('class');
                if ($grade === '' && $className === '') {
                    continue; // 空行
                }
                $totalRows++;

                $subjectName = $get('subject');
                $weekday = $this->parseWeekday($get('weekday'));
                $periodIndex = ctype_digit($get('period')) ? (int) $get('period') : 0;

                if ($subjectName === '' || $weekday === null || $periodIndex < 1) {
                    $errors[] = sprintf('第%d行：星期 / 第几节 / 科目无效（星期=%s 节次=%s 科目=%s），已跳过', $rowNo, $get('weekday') ?: '空', $get('period') ?: '空', $subjectName ?: '空');
                    continue;
                }

                $key = $grade . '|' . $className;
                if (!isset($classes[$key])) {
                    $classes[$key] = [
                        'grade' => $grade,
                        'name' => $className,
                        'class_id' => null,
                        'entries' => [],
                        'subjects' => [],
                        'row_label' => sprintf('%s/%s', $grade ?: '-', $className ?: '-'),
                    ];
                }
                $classes[$key]['subjects'][] = $subjectName;
                $classes[$key]['entries'][] = [
                    'weekday' => $weekday,
                    'period_index' => $periodIndex,
                    'week_type' => $this->parseWeekType($get('week_type')),
                    'subject_name' => $subjectName,
                    'teacher_name' => $get('teacher') ?: null,
                    'room' => $get('room') ?: null,
                ];

                // 节次时间（学校级共享，冲突时首次出现为准）
                $idx = $periodIndex;
                $start = $this->normalizeTime($get('start'));
                $end = $this->normalizeTime($get('end'));
                if ($start && $end) {
                    if (isset($periods[$idx]) && ($periods[$idx]['start_time'] !== $start || $periods[$idx]['end_time'] !== $end)) {
                        $errors[] = sprintf('第%d行：第%d节时间（%s–%s）与之前出现的不一致，以首次为准', $rowNo, $idx, $start, $end);
                    } else {
                        $periods[$idx] = ['start_time' => $start, 'end_time' => $end];
                    }
                } elseif (!isset($periods[$idx])) {
                    $errors[] = sprintf('第%d行：第%d节缺少开始 / 结束时间，导入后该节将无作息时间', $rowNo, $idx);
                }
            }

            // 解析班级 id
            $grades = array_values(array_unique(array_map(fn ($c) => $c['grade'], $classes)));
            $names = array_values(array_unique(array_map(fn ($c) => $c['name'], $classes)));
            $matched = ClassRoom::where('school_id', $schoolId)
                ->whereIn('grade', $grades)
                ->whereIn('name', $names)
                ->get(['id', 'grade', 'name']);
            foreach ($classes as &$c) {
                $found = $matched->first(fn ($m) => $m->grade === $c['grade'] && $m->name === $c['name']);
                $c['class_id'] = $found ? (int) $found->id : null;
                if (!$found) {
                    $errors[] = sprintf('班级「%s」不存在，该班 %d 行未导入（请先在班级列表创建）', $c['row_label'], count($c['entries']));
                }
            }
            unset($c);
        }

        $summary = [
            'total_rows' => $totalRows,
            'classes' => array_map(fn (array $c) => [
                'grade' => $c['grade'],
                'name' => $c['name'],
                'found' => $c['class_id'] !== null,
                'entry_count' => count($c['entries']),
            ], array_values($classes)),
            'period_count' => count($periods),
            'errors' => $errors,
            'dry_run' => $dryRun,
            'imported' => false,
        ];

        // 预览模式只返回解析结果；正式导入时匹配到的班级照常导入（错误行与未匹配班级已在 errors 中列明）
        if ($dryRun) {
            return $summary;
        }

        $importable = array_values(array_filter($classes, fn ($c) => $c['class_id'] !== null));

        DB::transaction(function () use ($importable, $periods, $schoolId): void {
            foreach ($importable as $c) {
                $this->save((int) $c['class_id'], $schoolId, [
                    'subjects' => array_map(fn (string $n) => ['name' => $n], array_values(array_unique($c['subjects']))),
                    'periods' => collect($periods)
                        ->map(fn (array $p, int $idx) => ['period_index' => $idx, 'start_time' => $p['start_time'], 'end_time' => $p['end_time']])
                        ->values()->all(),
                    'entries' => $c['entries'],
                ]);
            }

            // 受影响班级的待审申请自动作废（导入即覆盖）
            TimetableChangeRequest::whereIn('class_id', array_map(fn ($c) => $c['class_id'], $importable))
                ->where('status', TimetableChangeRequest::STATUS_PENDING)
                ->update([
                    'status' => TimetableChangeRequest::STATUS_REJECTED,
                    'review_note' => '管理员批量导入课表，本申请自动作废',
                    'reviewed_at' => now(),
                ]);
        });

        $summary['imported'] = true;

        return $summary;
    }

    /** 解析 CSV：去 BOM，首行为表头，返回 [表头, 数据行] */
    private function parseCsv(string $content): array
    {
        $content = preg_replace('/^\xEF\xBB\xBF/', '', $content) ?? $content;
        $lines = preg_split('/\r\n|\r|\n/', trim($content)) ?: [];

        $header = [];
        while ($header === [] && $lines !== []) {
            $header = array_map(fn ($h) => trim((string) $h), str_getcsv((string) array_shift($lines)));
        }

        return [$header, $lines];
    }

    /** 星期解析：支持 1-7 / 周一~周日 / 星期一~星期日（含「天」） */
    private function parseWeekday(string $raw): ?int
    {
        $raw = trim($raw);
        if ($raw === '') {
            return null;
        }
        if (ctype_digit($raw)) {
            $n = (int) $raw;

            return ($n >= 1 && $n <= 7) ? $n : null;
        }
        $map = ['一' => 1, '二' => 2, '三' => 3, '四' => 4, '五' => 5, '六' => 6, '日' => 7, '天' => 7];
        foreach ($map as $ch => $n) {
            if (mb_strpos($raw, $ch) !== false) {
                return $n;
            }
        }

        return null;
    }

    /** 周次解析：每周/all（默认）、单周/单/odd、双周/双/even */
    private function parseWeekType(string $raw): string
    {
        if (preg_match('/双|even/u', $raw)) {
            return 'even';
        }
        if (preg_match('/单|odd/u', $raw)) {
            return 'odd';
        }

        return 'all';
    }

    /** 时间归一化：8:00 / 08:00 / 0800 → HH:MM，非法返回 null */
    private function normalizeTime(string $raw): ?string
    {
        $raw = trim($raw);
        if (preg_match('/^(\d{1,2}):(\d{2})$/', $raw, $m)) {
            return sprintf('%02d:%s', (int) $m[1], $m[2]);
        }
        if (preg_match('/^(\d{2})(\d{2})$/', $raw, $m)) {
            return sprintf('%02d:%s', (int) $m[1], $m[2]);
        }

        return null;
    }

    // ============================================================
    // 任课设置（班级 × 科目 → 教师）
    // ============================================================

    /** 某班任课列表（subject_name => teacher_name） */
    public function listAssignments(int $classId): array
    {
        return TimetableTeacherAssignment::where('class_id', $classId)
            ->orderBy('subject_name')
            ->get(['subject_name', 'teacher_name'])
            ->map(fn (TimetableTeacherAssignment $a) => [
                'subject_name' => (string) $a->subject_name,
                'teacher_name' => (string) $a->teacher_name,
            ])
            ->values()
            ->all();
    }

    /** 整体保存某班任课（replace 语义：以提交列表为准） */
    public function saveAssignments(int $classId, int $schoolId, array $rows): void
    {
        DB::transaction(function () use ($classId, $schoolId, $rows): void {
            TimetableTeacherAssignment::where('class_id', $classId)->delete();

            foreach (array_values($rows) as $row) {
                $subject = trim((string) ($row['subject_name'] ?? ''));
                $teacher = trim((string) ($row['teacher_name'] ?? ''));
                if ($subject === '' || $teacher === '') {
                    continue;
                }

                TimetableTeacherAssignment::create([
                    'school_id' => $schoolId,
                    'class_id' => $classId,
                    'subject_name' => $subject,
                    'teacher_name' => $teacher,
                ]);
            }
        });
    }

    // ============================================================
    // 自动生成（规则排课）
    // ============================================================

    /**
     * 按规则自动生成课表（纯计算，不落库；前端预览后调 adminSave 保存）。
     *
     * 规则：rules.days 上课日；rules.subjects[] 每科目 weekly（每周节数）/
     * double（允许连堂）/ session（any|am|pm）/ max_per_day（同一天该科目上限）。
     * 算法：科目课时拆块（连堂优先成 2 节块）→ 随机贪心放置（分散优先）→
     * 失败重试最多 60 次。节次上午/下午按 start_time 是否早于 12:00 划分。
     *
     * @param  array{days?: array<int, int>, subjects?: array<int, array<string, mixed>>}  $rules
     * @return array{success: bool, warnings: array<int, string>, entries: array<int, array<string, mixed>>}
     */
    public function generate(int $schoolId, array $rules): array
    {
        $warnings = [];

        $periods = ClassPeriod::where('school_id', $schoolId)->orderBy('period_index')->get();
        if ($periods->isEmpty()) {
            return ['success' => false, 'warnings' => ['请先在「节次作息」中设置每日节次'], 'entries' => []];
        }

        $days = array_values(array_unique(array_filter(array_map('intval', $rules['days'] ?? [1, 2, 3, 4, 5]), fn ($d) => $d >= 1 && $d <= 7)));
        if ($days === []) {
            return ['success' => false, 'warnings' => ['请至少选择一个上课日'], 'entries' => []];
        }

        $periodIndexes = $periods->pluck('period_index')->map(fn ($i) => (int) $i)->all();
        $isAm = [];
        foreach ($periods as $p) {
            $ts = strtotime('1970-01-01 ' . (string) $p->start_time);
            $isAm[(int) $p->period_index] = $ts !== false && $ts < strtotime('1970-01-01 12:00');
        }

        // 拆课时块
        $blocks = [];
        foreach (array_values($rules['subjects'] ?? []) as $s) {
            $name = trim((string) ($s['name'] ?? ''));
            $weekly = (int) ($s['weekly'] ?? 0);
            if ($name === '' || $weekly < 1) {
                continue;
            }
            $double = (bool) ($s['double'] ?? false);
            $session = in_array($s['session'] ?? 'any', ['any', 'am', 'pm'], true) ? (string) $s['session'] : 'any';
            $maxPerDay = max(1, (int) ($s['max_per_day'] ?? 1));
            if ($double && $maxPerDay < 2) {
                $maxPerDay = 2; // 连堂本身占同日 2 节
            }

            $forbid = array_values(array_unique(array_map('intval', $s['forbid_periods'] ?? [])));

            $remaining = $weekly;
            if ($double) {
                while ($remaining >= 2) {
                    $blocks[] = ['subject' => $name, 'size' => 2, 'session' => $session, 'max_per_day' => $maxPerDay, 'forbid' => $forbid, 'teacher' => null];
                    $remaining -= 2;
                }
            }
            for ($i = 0; $i < $remaining; $i++) {
                $blocks[] = ['subject' => $name, 'size' => 1, 'session' => $session, 'max_per_day' => $maxPerDay, 'forbid' => $forbid, 'teacher' => null];
            }

            if ($weekly > $maxPerDay * count($days)) {
                $warnings[] = sprintf('科目「%s」每周 %d 节，超出每日上限 %d × %d 天，无法排下', $name, $weekly, $maxPerDay, count($days));
            }
        }

        $capacity = count($days) * count($periodIndexes);
        $needed = array_sum(array_map(fn (array $b) => $b['size'], $blocks));
        if ($needed > $capacity) {
            return ['success' => false, 'warnings' => array_merge($warnings, [sprintf('所需节数 %d 超过可用格数 %d（%d 天 × %d 节），请调整规则', $needed, $capacity, count($days), count($periodIndexes))]), 'entries' => []];
        }
        if ($needed === 0) {
            return ['success' => false, 'warnings' => array_merge($warnings, ['请为至少一个科目设置每周节数']), 'entries' => []];
        }

        for ($attempt = 0; $attempt < 60; $attempt++) {
            $entries = $this->tryGenerate($blocks, $days, $periodIndexes, $isAm);
            if ($entries !== null) {
                return ['success' => true, 'warnings' => $warnings, 'entries' => $entries];
            }
        }

        return ['success' => false, 'warnings' => array_merge($warnings, ['尝试多次均无法排出满足全部规则的课表，请减少节数、放宽连堂或时段限制']), 'entries' => []];
    }

    /** 单轮贪心：随机块序 + 候选位置按「同科目当日记数少者优先」放置，失败返回 null */
    private function tryGenerate(array $blocks, array $days, array $periodIndexes, array $isAm): ?array
    {
        shuffle($blocks);

        $grid = [];   // day => [period_index => subject]
        $perDay = []; // day => subject => count
        foreach ($days as $d) {
            $grid[$d] = [];
            $perDay[$d] = [];
        }

        foreach ($blocks as $block) {
            if (!$this->placeBlock($block, $days, $periodIndexes, $isAm, $grid, $perDay)) {
                return null;
            }
        }

        $entries = [];
        foreach ($grid as $day => $cells) {
            ksort($cells);
            foreach ($cells as $period => $subject) {
                $entries[] = [
                    'weekday' => (int) $day,
                    'period_index' => (int) $period,
                    'week_type' => 'all',
                    'subject_name' => $subject,
                    'teacher_name' => null,
                    'room' => null,
                ];
            }
        }

        return $entries;
    }

    /**
     * 放置一个课时块（贪心 + 候选评分）。
     *
     * $teacherBusy / $teacher 非空时（全校排课）需保证教师该时段全校无冲突。
     *
     * @param  array<int, int>  $periodIndexes
     * @param  array<int, bool>  $isAm
     * @param  array<int, array<int, string>>  $grid
     * @param  array<int, array<string, int>>  $perDay
     * @param  array<string, array<int, array<int, true>>>|null  $teacherBusy
     */
    private function placeBlock(array $block, array $days, array $periodIndexes, array $isAm, array &$grid, array &$perDay, ?array &$teacherBusy = null, ?string $teacher = null): bool
    {
        $size = (int) $block['size'];
        $subject = (string) $block['subject'];
        $session = (string) $block['session'];
        $maxPerDay = (int) $block['max_per_day'];
        $forbid = $block['forbid'] ?? [];

        // 枚举候选起点（size=2 需同日相邻节且同半天）
        $candidates = [];
        foreach ($days as $day) {
            foreach ($periodIndexes as $i => $period) {
                if ($i + $size > count($periodIndexes)) {
                    continue;
                }
                $cells = array_slice($periodIndexes, $i, $size);

                $ok = true;
                foreach ($cells as $cell) {
                    if (isset($grid[$day][$cell])) {
                        $ok = false;
                        break;
                    }
                    if (in_array((int) $cell, $forbid, true)) {
                        $ok = false;
                        break;
                    }
                    if ($teacher !== null && isset($teacherBusy[$teacher][$day][$cell])) {
                        $ok = false; // 教师冲突：该教师此时段已在其他班上课
                        break;
                    }
                    if (!$this->sessionMatches($session, $isAm[$cell] ?? null)) {
                        $ok = false;
                        break;
                    }
                }
                // 连堂不允许跨上 / 下午（如第 4 节上午 + 第 5 节下午）
                if ($ok && $size === 2 && ($isAm[$cells[0]] ?? null) !== null && ($isAm[$cells[1]] ?? null) !== null && $isAm[$cells[0]] !== $isAm[$cells[1]]) {
                    $ok = false;
                }
                if ($ok && ($perDay[$day][$subject] ?? 0) + $size > $maxPerDay) {
                    $ok = false;
                }
                if ($ok) {
                    $candidates[] = ['day' => $day, 'cells' => $cells];
                }
            }
        }

        if ($candidates === []) {
            return false;
        }

        // 评分：该日该科目已有越少越好（分散），其余随机打散
        shuffle($candidates);
        usort($candidates, fn ($a, $b) => ($perDay[$a['day']][$subject] ?? 0) <=> ($perDay[$b['day']][$subject] ?? 0));

        $best = $candidates[0];
        foreach ($best['cells'] as $cell) {
            $grid[$best['day']][$cell] = $subject;
            $perDay[$best['day']][$subject] = ($perDay[$best['day']][$subject] ?? 0) + 1;

            if ($teacher !== null) {
                $teacherBusy[$teacher][$best['day']][$cell] = true;
            }
        }

        return true;
    }

    /** 时段限定匹配：any 恒真；am/pm 需已知该节上/下午，未知（无时间）则放行 */
    private function sessionMatches(string $session, ?bool $isAm): bool
    {
        if ($session === 'any' || $isAm === null) {
            return true;
        }

        return $session === 'am' ? $isAm : !$isAm;
    }

    /**
     * 全校智能排课（参考水晶排课核心规则）。
     *
     * 依据「任课表」为全校各班同时排课，硬约束：同一教师同一时段全校仅一处；
     * 软规则与单班生成一致（每周节数 / 连堂 / 每日上限 / 限上下午 / 禁排节次 / 分散分布）。
     * 科目配置按科目名全局生效；某班某科目未登记任课教师则该班不排该科目。
     *
     * $commit=false 只验证可行性并返回摘要；true 则事务内落库全部班级并作废相关待审申请。
     *
     * @param  array{days?: array<int, int>, subjects?: array<int, array<string, mixed>>}  $rules
     * @return array{success: bool, warnings: array<int, string>, classes: array<int, array<string, mixed>>}
     */
    public function generateSchool(int $schoolId, array $rules, bool $commit): array
    {
        $periods = ClassPeriod::where('school_id', $schoolId)->orderBy('period_index')->get();
        if ($periods->isEmpty()) {
            return ['success' => false, 'warnings' => ['请先在「节次作息」中设置每日节次'], 'classes' => []];
        }

        $days = array_values(array_unique(array_filter(array_map('intval', $rules['days'] ?? [1, 2, 3, 4, 5]), fn ($d) => $d >= 1 && $d <= 7)));
        if ($days === []) {
            return ['success' => false, 'warnings' => ['请至少选择一个上课日'], 'classes' => []];
        }

        $periodIndexes = $periods->pluck('period_index')->map(fn ($i) => (int) $i)->all();
        $isAm = [];
        foreach ($periods as $p) {
            $ts = strtotime('1970-01-01 ' . (string) $p->start_time);
            $isAm[(int) $p->period_index] = $ts !== false && $ts < strtotime('1970-01-01 12:00');
        }

        // 任课表
        $assignments = TimetableTeacherAssignment::where('school_id', $schoolId)->get();
        if ($assignments->isEmpty()) {
            return ['success' => false, 'warnings' => ['尚未登记任课：请先在各班「任课设置」中填写科目 → 教师'], 'classes' => []];
        }
        $teacherOf = []; // classId => subjectName => teacher
        foreach ($assignments as $a) {
            $teacherOf[(int) $a->class_id][(string) $a->subject_name] = (string) $a->teacher_name;
        }

        // 科目配置（按名全局生效）
        $config = [];
        $warnings = [];
        foreach (array_values($rules['subjects'] ?? []) as $s) {
            $name = trim((string) ($s['name'] ?? ''));
            $weekly = (int) ($s['weekly'] ?? 0);
            if ($name === '' || $weekly < 1) {
                continue;
            }
            $double = (bool) ($s['double'] ?? false);
            $session = in_array($s['session'] ?? 'any', ['any', 'am', 'pm'], true) ? (string) $s['session'] : 'any';
            $maxPerDay = max(1, (int) ($s['max_per_day'] ?? 1));
            if ($double && $maxPerDay < 2) {
                $maxPerDay = 2;
            }
            $config[$name] = [
                'weekly' => $weekly,
                'double' => $double,
                'session' => $session,
                'max_per_day' => $maxPerDay,
                'forbid' => array_values(array_unique(array_map('intval', $s['forbid_periods'] ?? []))),
            ];
        }
        if ($config === []) {
            return ['success' => false, 'warnings' => ['请为至少一个科目设置每周节数'], 'classes' => []];
        }

        // 各班课时块
        $classes = ClassRoom::where('school_id', $schoolId)->orderBy('grade')->orderBy('name')->get(['id', 'name']);
        $classBlocks = []; // classId => ['name' =>, 'blocks' => []]
        foreach ($classes as $class) {
            $cid = (int) $class->id;
            if (empty($teacherOf[$cid])) {
                continue;
            }
            $blocks = [];
            foreach ($config as $name => $cfg) {
                $teacher = $teacherOf[$cid][$name] ?? null;
                if ($teacher === null) {
                    continue; // 该班此科目未任课
                }

                if ($cfg['weekly'] > $cfg['max_per_day'] * count($days)) {
                    $warnings[] = sprintf('%s「%s」每周 %d 节超出每日上限 × %d 天', $class->name, $name, $cfg['weekly'], count($days));
                }

                $remaining = $cfg['weekly'];
                if ($cfg['double']) {
                    while ($remaining >= 2) {
                        $blocks[] = ['subject' => $name, 'size' => 2, 'session' => $cfg['session'], 'max_per_day' => $cfg['max_per_day'], 'forbid' => $cfg['forbid'], 'teacher' => $teacher];
                        $remaining -= 2;
                    }
                }
                for ($i = 0; $i < $remaining; $i++) {
                    $blocks[] = ['subject' => $name, 'size' => 1, 'session' => $cfg['session'], 'max_per_day' => $cfg['max_per_day'], 'forbid' => $cfg['forbid'], 'teacher' => $teacher];
                }
            }
            if ($blocks !== []) {
                $classBlocks[$cid] = ['name' => (string) $class->name, 'blocks' => $blocks];
            }
        }

        if ($classBlocks === []) {
            return ['success' => false, 'warnings' => array_merge($warnings, ['任课表中的科目与规则配置无交集，请核对科目名称']), 'classes' => []];
        }

        // 教师不可用时段（学校级）：排课时视为该教师已被占用
        $unavailMap = [];
        foreach (TimetableTeacherUnavailability::where('school_id', $schoolId)->get() as $u) {
            $unavailMap[(string) $u->teacher_name][(int) $u->weekday][(int) $u->period_index] = true;
        }

        // 多轮重试：每次尝试全校统一放置，教师冲突为硬约束
        for ($attempt = 0; $attempt < 80; $attempt++) {
            $teacherBusy = $unavailMap;
            $results = []; // classId => entries
            $ok = true;

            foreach ($classBlocks as $classId => $info) {
                $grid = [];
                $perDay = [];
                foreach ($days as $d) {
                    $grid[$d] = [];
                    $perDay[$d] = [];
                }

                foreach ($info['blocks'] as $block) {
                    if (!$this->placeBlock($block, $days, $periodIndexes, $isAm, $grid, $perDay, $teacherBusy, (string) $block['teacher'])) {
                        $ok = false;
                        break 2;
                    }
                }

                $entries = [];
                foreach ($grid as $day => $cells) {
                    ksort($cells);
                    foreach ($cells as $period => $subject) {
                        $entries[] = [
                            'weekday' => (int) $day,
                            'period_index' => (int) $period,
                            'week_type' => 'all',
                            'subject_name' => $subject,
                            'teacher_name' => $teacherOf[$classId][$subject] ?? null,
                            'room' => null,
                        ];
                    }
                }
                $results[$classId] = $entries;
            }

            if ($ok) {
                if ($commit) {
                    DB::transaction(function () use ($results, $schoolId): void {
                        foreach ($results as $classId => $entries) {
                            $subjectNames = array_values(array_unique(array_map(fn (array $e) => $e['subject_name'], $entries)));
                            $this->save((int) $classId, $schoolId, [
                                'subjects' => array_map(fn (string $n) => ['name' => $n], $subjectNames),
                                'periods' => [],
                                'entries' => $entries,
                            ]);
                        }

                        TimetableChangeRequest::whereIn('class_id', array_keys($results))
                            ->where('status', TimetableChangeRequest::STATUS_PENDING)
                            ->update([
                                'status' => TimetableChangeRequest::STATUS_REJECTED,
                                'review_note' => '管理员全校自动排课，本申请自动作废',
                                'reviewed_at' => now(),
                            ]);
                    });
                }

                return [
                    'success' => true,
                    'warnings' => $warnings,
                    'classes' => array_map(fn (int $cid) => [
                        'class_id' => $cid,
                        'class_name' => $classBlocks[$cid]['name'],
                        'entry_count' => count($results[$cid]),
                    ], array_keys($results)),
                ];
            }
        }

        return ['success' => false, 'warnings' => array_merge($warnings, ['尝试多次均无法排出满足全部规则的全校课表（教师冲突难以避免），请减少节数、放宽连堂或调整任课']), 'classes' => []];
    }

    // ============================================================
    // 教师视角：我的课表（跨班聚合）/ 不可用时段 / 冲突检查
    // ============================================================

    /**
     * 某教师的周课表（全校各班聚合，按排课里的教师姓名匹配）。
     *
     * @return array{subjects: array<int, array<string, mixed>>, periods: array<int, mixed>, entries: array<int, array<string, mixed>>}
     */
    public function forTeacher(int $schoolId, string $teacherName): array
    {
        $subjects = Subject::where('school_id', $schoolId)
            ->orderBy('sort_order')->orderBy('id')
            ->get(['name', 'simplified_name', 'color']);
        $nameById = Subject::where('school_id', $schoolId)->pluck('name', 'id');

        $classNames = ClassRoom::where('school_id', $schoolId)->pluck('name', 'id');

        $entries = TimetableEntry::whereHas('classRoom', fn ($q) => $q->where('school_id', $schoolId))
            ->where('teacher_name', $teacherName)
            ->orderBy('weekday')->orderBy('period_index')
            ->get(['class_id', 'weekday', 'period_index', 'week_type', 'subject_id', 'room'])
            ->map(fn (TimetableEntry $e) => [
                'class_id' => (int) $e->class_id,
                'class_name' => $classNames[(int) $e->class_id] ?? null,
                'weekday' => (int) $e->weekday,
                'period_index' => (int) $e->period_index,
                'week_type' => (string) ($e->week_type ?: 'all'),
                'subject_name' => $e->subject_id ? ($nameById[$e->subject_id] ?? null) : null,
                'room' => $e->room,
            ])
            ->filter(fn (array $e) => $e['subject_name'] !== null)
            ->values();

        return [
            'teacher_name' => $teacherName,
            'subjects' => $subjects->values(),
            'periods' => ClassPeriod::where('school_id', $schoolId)
                ->orderBy('period_index')
                ->get(['period_index', 'name', 'start_time', 'end_time'])
                ->values(),
            'entries' => $entries,
        ];
    }

    /** 全校教师不可用时段（平铺列表） */
    public function listUnavailabilities(int $schoolId): array
    {
        return TimetableTeacherUnavailability::where('school_id', $schoolId)
            ->orderBy('teacher_name')->orderBy('weekday')->orderBy('period_index')
            ->get(['teacher_name', 'weekday', 'period_index'])
            ->map(fn (TimetableTeacherUnavailability $u) => [
                'teacher_name' => (string) $u->teacher_name,
                'weekday' => (int) $u->weekday,
                'period_index' => (int) $u->period_index,
            ])
            ->values()
            ->all();
    }

    /** 整体保存某教师的不可用时段（replace 语义，cells = [{weekday, period_index}]） */
    public function saveUnavailabilities(int $schoolId, string $teacherName, array $cells): void
    {
        DB::transaction(function () use ($schoolId, $teacherName, $cells): void {
            TimetableTeacherUnavailability::where('school_id', $schoolId)
                ->where('teacher_name', $teacherName)
                ->delete();

            foreach (array_values($cells) as $c) {
                $weekday = (int) ($c['weekday'] ?? 0);
                $periodIndex = (int) ($c['period_index'] ?? 0);
                if ($weekday < 1 || $weekday > 7 || $periodIndex < 1) {
                    continue;
                }

                TimetableTeacherUnavailability::create([
                    'school_id' => $schoolId,
                    'teacher_name' => $teacherName,
                    'weekday' => $weekday,
                    'period_index' => $periodIndex,
                ]);
            }
        });
    }

    /**
     * 冲突检查：给定某班当前编辑中的排课，检查
     * ① 教师冲突——同一教师在其他班级同一时段（周次重叠）已有课；
     * ② 不可用冲突——教师被标记了不可用时段。
     * 返回冲突明细列表（空数组 = 无冲突）。
     *
     * @param  array<int, array<string, mixed>>  $entries
     * @return array<int, array<string, mixed>>
     */
    public function checkConflicts(int $schoolId, int $classId, array $entries): array
    {
        $teachers = array_values(array_unique(array_filter(
            array_map(fn ($e) => trim((string) ($e['teacher_name'] ?? '')), $entries),
            fn (string $t) => $t !== '',
        )));
        if ($teachers === []) {
            return [];
        }

        // 其他班级相关教师的全部排课
        $others = TimetableEntry::whereHas('classRoom', fn ($q) => $q->where('school_id', $schoolId))
            ->where('class_id', '!=', $classId)
            ->whereIn('teacher_name', $teachers)
            ->get(['class_id', 'weekday', 'period_index', 'week_type', 'teacher_name']);

        $classNames = ClassRoom::where('school_id', $schoolId)->pluck('name', 'id');
        $nameById = Subject::where('school_id', $schoolId)->pluck('name', 'id');

        // 不可用时段
        $unavail = [];
        foreach (TimetableTeacherUnavailability::where('school_id', $schoolId)->whereIn('teacher_name', $teachers)->get() as $u) {
            $unavail[(string) $u->teacher_name][(int) $u->weekday][(int) $u->period_index] = true;
        }

        $conflicts = [];
        $seen = [];

        foreach ($entries as $e) {
            $teacher = trim((string) ($e['teacher_name'] ?? ''));
            $weekday = (int) ($e['weekday'] ?? 0);
            $periodIndex = (int) ($e['period_index'] ?? 0);
            $weekType = (string) ($e['week_type'] ?? 'all');
            $subjectName = (string) ($e['subject_name'] ?? '');

            if ($teacher === '' || $weekday < 1 || $weekday > 7 || $periodIndex < 1) {
                continue;
            }

            // ① 与其他班的教师冲突（周次重叠才算：all 与任何都重叠；odd/even 仅同类重叠）
            foreach ($others as $o) {
                if ((string) $o->teacher_name !== $teacher
                    || (int) $o->weekday !== $weekday
                    || (int) $o->period_index !== $periodIndex
                    || !$this->weekTypeOverlaps($weekType, (string) ($o->week_type ?: 'all'))) {
                    continue;
                }

                $key = $teacher . '|' . $weekday . '|' . $periodIndex . '|' . $weekType . '|' . $o->class_id;
                if (isset($seen[$key])) {
                    continue;
                }
                $seen[$key] = true;

                $conflicts[] = [
                    'type' => 'teacher',
                    'teacher_name' => $teacher,
                    'weekday' => $weekday,
                    'period_index' => $periodIndex,
                    'week_type' => $weekType,
                    'subject_name' => $subjectName,
                    'other_class_id' => (int) $o->class_id,
                    'other_class_name' => $classNames[(int) $o->class_id] ?? null,
                    'other_subject_name' => $o->subject_id ? ($nameById[$o->subject_id] ?? null) : null,
                    'message' => sprintf(
                        '%s%s第%d节「%s」与 %s「%s」冲突（同一教师）',
                        self::WEEKDAY_LABELS[$weekday] ?? "第{$weekday}天",
                        $weekType === 'all' ? '' : ($weekType === 'odd' ? '单周' : '双周'),
                        $periodIndex,
                        $subjectName ?: '课程',
                        $classNames[(int) $o->class_id] ?? '其他班级',
                        $o->subject_id ? ($nameById[$o->subject_id] ?? '') : '',
                    ),
                ];
            }

            // ② 教师不可用时段
            if (isset($unavail[$teacher][$weekday][$periodIndex])) {
                $key = 'u|' . $teacher . '|' . $weekday . '|' . $periodIndex . '|' . $weekType;
                if (!isset($seen[$key])) {
                    $seen[$key] = true;
                    $conflicts[] = [
                        'type' => 'unavailable',
                        'teacher_name' => $teacher,
                        'weekday' => $weekday,
                        'period_index' => $periodIndex,
                        'week_type' => $weekType,
                        'subject_name' => $subjectName,
                        'other_class_id' => null,
                        'other_class_name' => null,
                        'other_subject_name' => null,
                        'message' => sprintf(
                            '%s 第%d节「%s」：教师 %s 此时段已被标记为不可用',
                            self::WEEKDAY_LABELS[$weekday] ?? "第{$weekday}天",
                            $periodIndex,
                            $subjectName ?: '课程',
                            $teacher,
                        ),
                    ];
                }
            }
        }

        return $conflicts;
    }

    /** 周次是否重叠：all 与任何都重叠；odd/even 仅同类重叠 */
    private function weekTypeOverlaps(string $a, string $b): bool
    {
        if ($a === 'all' || $b === 'all') {
            return true;
        }

        return $a === $b;
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
