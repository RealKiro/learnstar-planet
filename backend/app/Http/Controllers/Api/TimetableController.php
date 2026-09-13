<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\User;
use App\Services\TeacherClassScope;
use App\Services\TimetableService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

/**
 * 课表（教师端）：读取 / 保存 / 导出 CSES
 *
 * 一切操作都经 TeacherClassScope 限定在教师可管理的班级内。
 * 导出格式 CSES 可被 ClassIsland「从 CSES 导入」直接消费。
 */
class TimetableController extends Controller
{
    public function __construct(
        private readonly TimetableService $timetableService,
        private readonly TeacherClassScope $classScope,
    ) {
    }

    /** 课表初始化数据：科目 + 节次 + 该班排课 */
    public function show(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = $this->resolveClassId($request, $teacher);

        if (!$classId) {
            return response()->json(['message' => '当前账号没有可管理的班级'], 400);
        }

        return response()->json([
            'data' => [
                'class_id' => $classId,
                ...$this->timetableService->bootstrap($classId, (int) $teacher->school_id),
            ],
        ]);
    }

    /** 教师提交课表修改申请（不直接生效，待学校管理员审核） */
    public function save(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = $this->resolveClassId($request, $teacher);

        if (!$classId) {
            return response()->json(['message' => '当前账号没有可管理的班级'], 400);
        }

        $payload = $this->validatedPayload($request);

        $requestModel = $this->timetableService->submitChange($classId, (int) $teacher->school_id, (int) $teacher->id, $payload);

        return response()->json([
            'message' => '修改申请已提交，待管理员审核',
            'data' => ['request_id' => $requestModel->id],
        ], 201);
    }

    /** 我的班级申请历史 */
    public function myChanges(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classId = $this->resolveClassId($request, $teacher);

        if (!$classId) {
            return response()->json(['message' => '当前账号没有可管理的班级'], 400);
        }

        return response()->json(['data' => $this->timetableService->listChangesForClass($classId)]);
    }

    // ============================================================
    // 管理员审批（/api/v1/admin/timetable/*，role:school_admin）
    // ============================================================

    /** 全校申请列表（可按 status 过滤） */
    public function adminChanges(Request $request): JsonResponse
    {
        $status = $request->query('status');

        return response()->json([
            'data' => $this->timetableService->listChangesForSchool((int) $request->user()->school_id, $status),
        ]);
    }

    /** 通过申请并应用 */
    public function approve(Request $request, int $id): JsonResponse
    {
        $request->validate(['note' => 'nullable|string|max:200']);

        $ok = $this->timetableService->approveChange($id, (int) $request->user()->id, $request->input('note'));

        return $ok
            ? response()->json(['message' => '已通过并应用课表'])
            : response()->json(['message' => '申请不存在或已处理'], 409);
    }

    /** 驳回申请 */
    public function reject(Request $request, int $id): JsonResponse
    {
        $request->validate(['note' => 'nullable|string|max:200']);

        $ok = $this->timetableService->rejectChange($id, (int) $request->user()->id, $request->input('note'));

        return $ok
            ? response()->json(['message' => '已驳回'])
            : response()->json(['message' => '申请不存在或已处理'], 409);
    }

    /** 提取并校验课表快照 payload */
    private function validatedPayload(Request $request): array
    {
        $request->validate([
            'subjects' => 'nullable|array',
            'subjects.*.name' => 'required_with:subjects|string|max:50',
            'periods' => 'nullable|array',
            'periods.*.period_index' => 'required_with:periods|integer|min:1|max:30',
            'periods.*.start_time' => 'required_with:periods|string|max:8',
            'periods.*.end_time' => 'required_with:periods|string|max:8',
            'entries' => 'nullable|array',
            'entries.*.weekday' => 'required_with:entries|integer|min:1|max:7',
            'entries.*.period_index' => 'required_with:entries|integer|min:1|max:30',
            'entries.*.subject_name' => 'required_with:entries|string|max:50',
            'entries.*.week_type' => 'nullable|string|in:all,odd,even',
        ]);

        return [
            'subjects' => $request->input('subjects', []),
            'periods' => $request->input('periods', []),
            'entries' => $request->input('entries', []),
        ];
    }

    /** 导出 CSES（.yaml），可直接在 ClassIsland「从 CSES 导入」 */
    public function exportCses(Request $request): Response
    {
        $teacher = $request->user();
        $classId = $this->resolveClassId($request, $teacher);

        if (!$classId) {
            return response()->json(['message' => '当前账号没有可管理的班级'], 400);
        }

        $yaml = $this->timetableService->toCses($classId, (int) $teacher->school_id);
        $className = optional(ClassRoom::find($classId))->name ?? 'class';
        $fileName = $className . '-课表.cses.yaml';

        return response($yaml, 200, [
            'Content-Type' => 'application/x-yaml; charset=UTF-8',
            'Content-Disposition' => 'attachment; filename="' . rawurlencode($fileName) . '"',
            'Cache-Control' => 'no-store',
        ]);
    }

    /** 解析并校验班级归属（越权返回 null） */
    private function resolveClassId(Request $request, User $teacher): ?int
    {
        $ids = $this->classScope->ids($teacher);
        $classId = (int) $request->input('class_id', $ids->first() ?? 0);

        return $ids->contains($classId) ? $classId : null;
    }
}
