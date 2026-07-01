<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\Score;
use App\Models\ScoreLog;
use App\Models\School;
use App\Models\Student;
use App\Models\User;
use App\Services\AuthService;
use Carbon\Carbon;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Validator;

class SchoolAdminController extends Controller
{
    public function __construct(
        private readonly AuthService $authService
    ) {
    }

    // ===== 学校管理 =====

    /**
     * 获取当前管理员所在学校信息
     */
    public function getSchool(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        if (!$school) {
            return response()->json(['message' => '学校不存在'], 404);
        }

        return response()->json([
            'data' => [
                'id' => $school->id,
                'name' => $school->name,
                'code' => $school->code,
                'address' => $school->address,
                'contact_phone' => $school->contact_phone,
                'contact_email' => $school->contact_email,
                'logo_path' => $school->logo_path,
                'settings' => $school->settings ?? [],
            ],
        ]);
    }

    /**
     * 更新学校信息
     */
    public function updateSchool(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'name' => 'sometimes|required|string|max:100',
            'address' => 'nullable|string|max:255',
            'contact_phone' => 'nullable|string|max:30',
            'contact_email' => 'nullable|email|max:100',
            'logo_path' => 'nullable|string|max:255',
            'settings' => 'nullable|array',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        if (!$school) {
            return response()->json(['message' => '学校不存在'], 404);
        }

        $data = $request->only(['name', 'address', 'contact_phone', 'contact_email', 'logo_path', 'settings']);
        $school->fill($data);
        $school->save();

        return response()->json(['message' => '学校信息已更新', 'data' => $school->fresh()]);
    }

    // ===== 教师管理 =====

    /**
     * 批量创建教师账号
     */
    public function batchCreateTeachers(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'teachers' => 'required|array|min:1',
            'teachers.*.name' => 'required|string|max:50',
            'teachers.*.phone' => 'nullable|string|max:30',
            'teachers.*.email' => 'nullable|email|max:100',
            'teachers.*.username' => 'nullable|string|max:50',
            'teachers.*.password' => 'nullable|string|min:6|max:50',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $created = $this->authService->createTeacherAccounts($school, $request->input('teachers'));

        return response()->json(['data' => $created]);
    }

    /**
     * 教师列表（含绑定情况）
     */
    public function listTeachers(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        $teachers = User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->with(['thirdPartyBindings', 'classesAsTeacher'])
            ->orderBy('created_at', 'desc')
            ->get()
            ->map(function (User $t) {
                $bindings = $t->thirdPartyBindings->pluck('platform')->toArray();
                $classes = $t->classesAsTeacher->pluck('name')->toArray();

                return [
                    'id' => $t->id,
                    'username' => $t->username,
                    'name' => $t->name,
                    'phone' => $t->phone,
                    'email' => $t->email,
                    'status' => $t->status,
                    'password_changed' => $t->password_changed,
                    'last_login_at' => $t->last_login_at?->toDateTimeString(),
                    'bindings' => $bindings,
                    'class_names' => $classes,
                    'created_at' => $t->created_at?->toDateTimeString(),
                ];
            });

        return response()->json(['data' => $teachers]);
    }

    /**
     * 更新教师信息
     */
    public function updateTeacher(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $teacher = User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->findOrFail($id);

        $validator = Validator::make($request->all(), [
            'name' => 'sometimes|required|string|max:50',
            'phone' => 'nullable|string|max:30',
            'email' => 'nullable|email|max:100',
            'status' => 'sometimes|required|in:active,disabled',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $teacher->fill($request->only(['name', 'phone', 'email', 'status']));
        $teacher->save();

        return response()->json(['message' => '更新成功', 'data' => $teacher->fresh()]);
    }

    /**
     * 重置教师密码
     */
    public function resetTeacherPassword(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $teacher = User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->findOrFail($id);

        $newPassword = $request->input('password', str()->random(8));
        $teacher->password = Hash::make($newPassword);
        $teacher->password_changed = false;
        $teacher->save();

        return response()->json([
            'message' => '密码已重置',
            'data' => ['new_password' => $newPassword],
        ]);
    }

    /**
     * 软删除教师账号（并解除班级关联）
     */
    public function disableTeacher(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $teacher = User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->findOrFail($id);

        // 解除该教师的所有班级关联
        ClassRoom::where('teacher_id', $teacher->id)->update(['teacher_id' => null]);
        // 软删除
        $teacher->delete();

        return response()->json(['message' => '教师账号已删除']);
    }

    // ===== 家长管理 =====

    /**
     * 批量创建家长账号
     */
    public function batchCreateParents(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'parents' => 'required|array|min:1',
            'parents.*.name' => 'required|string|max:50',
            'parents.*.phone' => 'nullable|string|max:30',
            'parents.*.student_id' => 'nullable|integer|exists:students,id',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $created = [];

        foreach ($request->input('parents') as $parentData) {
            $parent = $this->authService->createParentAccount($school, $parentData);
            if (!empty($parentData['student_id'])) {
                Student::where('id', $parentData['student_id'])->update(['parent_id' => $parent->id]);
            }
            $created[] = [
                'id' => $parent->id,
                'username' => $parent->username,
                'initial_password' => $parent->getRawOriginal('password')
                    ? '' // 已 hash，无法显示
                    : '',
                'name' => $parent->name,
            ];
        }

        return response()->json(['data' => $created]);
    }

    /**
     * 家长列表
     */
    public function listParents(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        $parents = User::where('school_id', $school->id)
            ->where('role', 'parent')
            ->with('children')
            ->orderBy('created_at', 'desc')
            ->get()
            ->map(function (User $p) {
                $childrenNames = $p->children->pluck('name')->toArray();
                return [
                    'id' => $p->id,
                    'username' => $p->username,
                    'name' => $p->name,
                    'phone' => $p->phone,
                    'status' => $p->status,
                    'children' => $childrenNames,
                    'created_at' => $p->created_at?->toDateTimeString(),
                ];
            });

        return response()->json(['data' => $parents]);
    }

    /**
     * 软删除家长账号
     */
    public function deleteParent(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $parent = User::where('school_id', $school->id)
            ->where('role', 'parent')
            ->findOrFail($id);
        // 解除学生绑定
        Student::where('parent_id', $parent->id)->update(['parent_id' => null]);
        $parent->delete();

        return response()->json(['message' => '家长账号已删除']);
    }

    // ===== 班级管理 =====

    public function index(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        $classes = ClassRoom::where('school_id', $school->id)
            ->with(['teacher', 'students'])
            ->orderBy('grade')
            ->orderBy('name')
            ->get()
            ->map(function (ClassRoom $c) {
                return [
                    'id' => $c->id,
                    'name' => $c->name,
                    'grade' => $c->grade,
                    'year' => $c->year,
                    'teacher_id' => $c->teacher_id,
                    'teacher_name' => $c->teacher?->name,
                    'student_count' => $c->students->count(),
                    'max_students' => $c->max_students,
                    'status' => $c->status,
                    'created_at' => $c->created_at?->toDateTimeString(),
                ];
            });

        return response()->json(['data' => $classes]);
    }

    public function store(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'name' => 'required|string|max:50',
            'grade' => 'nullable|string|max:50',
            'year' => 'nullable|string|max:20',
            'teacher_id' => 'nullable|integer|exists:users,id',
            'max_students' => 'nullable|integer|min:0',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $class = ClassRoom::create([
            'school_id' => $school->id,
            'name' => $request->input('name'),
            'grade' => $request->input('grade'),
            'year' => $request->input('year'),
            'teacher_id' => $request->input('teacher_id'),
            'max_students' => $request->input('max_students', 0),
            'status' => 'active',
        ]);

        return response()->json(['message' => '班级已创建', 'data' => $class], 201);
    }

    public function show(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)
            ->with(['teacher', 'students'])
            ->findOrFail($id);

        return response()->json(['data' => $class]);
    }

    public function update(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($id);

        $validator = Validator::make($request->all(), [
            'name' => 'sometimes|required|string|max:50',
            'grade' => 'nullable|string|max:50',
            'year' => 'nullable|string|max:20',
            'teacher_id' => 'nullable|integer|exists:users,id',
            'max_students' => 'nullable|integer|min:0',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $class->fill($request->only(['name', 'grade', 'year', 'teacher_id', 'max_students']));
        $class->save();

        return response()->json(['message' => '班级已更新', 'data' => $class->fresh()]);
    }

    public function destroy(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($id);
        $class->delete();

        return response()->json(['message' => '班级已删除']);
    }

    public function assignClassTeacher(Request $request, int $id): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'teacher_id' => 'required|integer|exists:users,id',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误'], 422);
        }

        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($id);
        $class->teacher_id = $request->input('teacher_id');
        $class->save();

        return response()->json(['message' => '班主任已分配', 'data' => $class->fresh()]);
    }

    /**
     * Excel 导入学生到指定班级
     * 请求体: {students: [{name, student_no?, phone?}, ...]}
     */
    public function importStudents(Request $request, int $id): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'students' => 'required|array|min:1',
            'students.*.name' => 'required|string|max:50',
            'students.*.student_no' => 'nullable|string|max:50',
            'students.*.phone' => 'nullable|string|max:30',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($id);

        $created = [];
        foreach ($request->input('students') as $row) {
            $student = Student::create([
                'class_id' => $class->id,
                'name' => $row['name'],
                'student_no' => $row['student_no'] ?? null,
                'total_score' => 0,
                'status' => 'active',
            ]);
            $created[] = ['id' => $student->id, 'name' => $student->name, 'student_no' => $student->student_no];
        }

        return response()->json([
            'message' => '导入完成',
            'data' => ['class_id' => $class->id, 'class_name' => $class->name, 'created' => $created],
        ]);
    }

    // ===== 报表 =====

    /**
     * 全校概览
     */
    public function schoolOverview(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        $classCount = ClassRoom::where('school_id', $school->id)->count();
        $teacherCount = User::where('school_id', $school->id)->where('role', 'teacher')->where('status', 'active')->count();
        $parentCount = User::where('school_id', $school->id)->where('role', 'parent')->where('status', 'active')->count();
        $studentCount = Student::whereIn('class_id', ClassRoom::where('school_id', $school->id)->pluck('id'))->count();

        $startOfMonth = Carbon::now()->startOfMonth();
        $endOfMonth = Carbon::now()->endOfMonth();
        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');

        $monthlyScore = (int) Score::whereIn('class_id', $classIds)
            ->whereBetween('created_at', [$startOfMonth, $endOfMonth])
            ->sum('amount');

        $lastMonthStart = Carbon::now()->subMonth()->startOfMonth();
        $lastMonthEnd = Carbon::now()->subMonth()->endOfMonth();
        $lastMonthScore = (int) Score::whereIn('class_id', $classIds)
            ->whereBetween('created_at', [$lastMonthStart, $lastMonthEnd])
            ->sum('amount');

        $scoreTrend = $lastMonthScore > 0
            ? round(($monthlyScore - $lastMonthScore) / max($lastMonthScore, 1) * 100, 1)
            : ($monthlyScore > 0 ? 100 : 0);

        return response()->json([
            'data' => [
                'class_count' => $classCount,
                'teacher_count' => $teacherCount,
                'parent_count' => $parentCount,
                'student_count' => $studentCount,
                'monthly_score' => $monthlyScore,
                'last_month_score' => $lastMonthScore,
                'score_trend_percent' => $scoreTrend,
                'month_label' => Carbon::now()->format('Y-m'),
            ],
        ]);
    }

    public function reportsByGrade(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        $rows = ClassRoom::where('school_id', $school->id)
            ->withCount('students')
            ->get()
            ->groupBy('grade')
            ->map(function ($classes, $grade) {
                return [
                    'grade' => $grade ?: '未分年级',
                    'class_count' => $classes->count(),
                    'student_count' => $classes->sum('students_count'),
                ];
            })->values();

        return response()->json(['data' => $rows]);
    }

    public function reportsByClass(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        $rows = ClassRoom::where('school_id', $school->id)
            ->with(['teacher'])
            ->withCount('students')
            ->get()
            ->map(function (ClassRoom $c) {
                $classScore = (int) Score::where('class_id', $c->id)
                    ->whereMonth('created_at', Carbon::now()->month)
                    ->sum('amount');
                return [
                    'class_id' => $c->id,
                    'class_name' => $c->name,
                    'grade' => $c->grade,
                    'teacher_name' => $c->teacher?->name,
                    'student_count' => $c->students_count,
                    'monthly_score' => $classScore,
                ];
            });

        return response()->json(['data' => $rows]);
    }
}
