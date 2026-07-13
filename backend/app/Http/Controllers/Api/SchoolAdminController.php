<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\ClassRoomTeacher;
use App\Models\Score;
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
     * 单独创建教师账号（可选同时分配班级和角色）
     */
    public function createTeacher(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'name' => 'required|string|max:50',
            'nickname' => 'nullable|string|max:80',
            'subject' => 'nullable|string|max:50',
            'grade_team' => 'nullable|string|max:50',
            'phone' => 'nullable|string|max:30',
            'email' => 'nullable|email|max:100',
            'username' => 'nullable|string|max:50',
            'password' => 'required|string|min:6|max:50',
            'assignments' => 'nullable|array',
            'assignments.*.class_id' => 'required|integer|exists:class_rooms,id',
            'assignments.*.role' => 'required|string|in:head_teacher,co_teacher,subject_teacher,grade_lead,admin_director',
            'assignments.*.subject' => 'nullable|string|max:50',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        if (!$school instanceof \App\Models\School) {
            return response()->json(['message' => '未找到学校'], 404);
        }

        $teacherData = $request->only([
            'name', 'nickname', 'subject', 'grade_team', 'phone', 'email', 'username', 'password',
        ]);
        $teacherData['password'] = $teacherData['password'] ?? \Illuminate\Support\Str::random(10);

        $created = $this->authService->createTeacherAccounts($school, [$teacherData]);
        $teacher = $created[0] ?? null;

        if ($teacher) {
            $assignmentsInput = $request->input('assignments', []);
            $assigned = [];
            foreach ($assignmentsInput as $assignment) {
                $classId = (int) $assignment['class_id'];
                $role = $assignment['role'];

                ClassRoomTeacher::updateOrCreate(
                    ['class_room_id' => $classId, 'user_id' => $teacher['id']],
                    ['role' => $role, 'subject' => $assignment['subject'] ?? null],
                );

                if ($role === 'head_teacher') {
                    ClassRoom::where('id', $classId)
                        ->where('school_id', $school->id)
                        ->update(['teacher_id' => $teacher['id']]);
                }

                $classRoom = ClassRoom::find($classId);
                $assigned[] = [
                    'class_id' => $classId,
                    'class_name' => $classRoom?->name,
                    'role' => $role,
                    'subject' => $assignment['subject'] ?? null,
                ];
            }
            $teacher['assignments'] = $assigned;
        }

        return response()->json([
            'message' => '已创建教师「' . ($teacher['name'] ?? '') . '」',
            'data' => $teacher,
        ], 201);
    }

    /**
     * 批量创建教师账号
     */
    public function batchCreateTeachers(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'teachers' => 'required|array|min:1',
            'teachers.*.name' => 'required|string|max:50',
            'teachers.*.nickname' => 'nullable|string|max:80',
            'teachers.*.subject' => 'nullable|string|max:50',
            'teachers.*.grade_team' => 'nullable|string|max:50',
            'teachers.*.phone' => 'nullable|string|max:30',
            'teachers.*.email' => 'nullable|email|max:100',
            'teachers.*.username' => 'nullable|string|max:50',
            'teachers.*.password' => 'required|string|min:6|max:50',
            'teachers.*.assignments' => 'nullable|array',
            'teachers.*.assignments.*.class_id' => 'required|integer|exists:class_rooms,id',
            'teachers.*.assignments.*.role' => 'required|string|in:head_teacher,co_teacher,subject_teacher,grade_lead,admin_director',
            'teachers.*.assignments.*.subject' => 'nullable|string|max:50',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        if (!$school instanceof \App\Models\School) {
            return response()->json(['message' => '未找到学校'], 404);
        }

        $teachers = $request->input('teachers');
        $created = $this->authService->createTeacherAccounts($school, $teachers);

        foreach ($created as $teacher) {
            if (empty($teacher['id'])) {
                continue;
            }

            $teacherInput = collect($teachers)->firstWhere('name', $teacher['name']);
            if (!$teacherInput || empty($teacherInput['assignments'])) {
                continue;
            }

            foreach ($teacherInput['assignments'] as $assignment) {
                ClassRoomTeacher::updateOrCreate(
                    ['class_room_id' => (int) $assignment['class_id'], 'user_id' => $teacher['id']],
                    [
                        'role' => $assignment['role'],
                        'subject' => $assignment['subject'] ?? null,
                    ],
                );

                if ($assignment['role'] === 'head_teacher') {
                    ClassRoom::where('id', (int) $assignment['class_id'])
                        ->where('school_id', $school->id)
                        ->update(['teacher_id' => $teacher['id']]);
                }
            }
        }

        return response()->json(['data' => $created]);
    }

    /**
     * 教师列表（含绑定情况）
     */
    public function listTeachers(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        /** @var \Illuminate\Database\Eloquent\Collection<int, User> $teacherUsers */
        $teacherUsers = User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->with(['thirdPartyBindings', 'classRoomAssignments.classRoom'])
            ->orderBy('created_at', 'desc')
            ->get();

        $teachers = $teacherUsers->map(function (User $t) {
            $bindings = $t->thirdPartyBindings->pluck('platform')->toArray();
            $assignments = $t->classRoomAssignments->map(fn ($a) => [
                'class_id' => (int) $a->class_room_id,
                'class_name' => $a->classRoom?->name,
                'grade' => $a->classRoom?->grade,
                'role' => $a->role,
                'subject' => $a->subject,
            ])->values()->toArray();

            return [
                'id' => $t->id,
                'username' => $t->username,
                'name' => $t->name,
                'nickname' => $t->nickname,
                'subject' => $t->subject,
                'grade_team' => $t->grade_team,
                'avatar_path' => $t->avatar_path,
                'phone' => $t->phone,
                'email' => $t->email,
                'status' => $t->status,
                'password_changed' => $t->password_changed,
                'last_login_at' => $t->last_login_at?->toDateTimeString(),
                'bindings' => $bindings,
                'assignments' => $assignments,
                'class_names' => $assignments ? array_column($assignments, 'class_name') : [],
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

        $teacher->fill($request->only(['name', 'nickname', 'subject', 'phone', 'email', 'status']));
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
            'parents.*.password' => 'required|string|min:6|max:50',
            'parents.*.phone' => 'nullable|string|max:30',
            'parents.*.student_id' => 'nullable|integer|exists:students,id',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        if (!$school instanceof \App\Models\School) {
            return response()->json(['message' => '未找到学校'], 404);
        }
        $created = [];

        foreach ($request->input('parents') as $parentData) {
            $result = $this->authService->createParentAccount($school, $parentData);
            if (!empty($parentData['student_id'])) {
                Student::where('id', $parentData['student_id'])->update(['parent_id' => $result['id']]);
            }
            $created[] = $result;
        }

        return response()->json(['data' => $created]);
    }

    /**
     * 家长列表
     */
    public function listParents(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        /** @var \Illuminate\Database\Eloquent\Collection<int, User> $parentUsers */
        $parentUsers = User::where('school_id', $school->id)
            ->where('role', 'parent')
            ->with('children')
            ->orderBy('created_at', 'desc')
            ->get();

        $parents = $parentUsers->map(function (User $p) {
            $childrenNames = $p->children->pluck('name')->toArray();

            return [
                'id' => $p->id,
                'username' => $p->username,
                'name' => $p->name,
                'nickname' => $p->nickname,
                'avatar_path' => $p->avatar_path,
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

        /** @var \Illuminate\Database\Eloquent\Collection<int, ClassRoom> $allClasses */
        $allClasses = ClassRoom::where('school_id', $school->id)
            ->with(['teacher', 'students'])
            ->orderBy('grade')
            ->orderBy('name')
            ->get();

        $classes = [];
        foreach ($allClasses as $c) {
            $classes[] = [
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
        }

        return response()->json(['data' => $classes]);
    }

    public function store(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'name' => 'required|string|max:50',
            'grade' => 'nullable|string|max:50',
            'year' => 'nullable|max:20',
            'teacher_id' => 'nullable|integer|exists:users,id',
            'max_students' => 'nullable|integer|min:0',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;

        // 检查班级名称是否已存在
        $name = $request->input('name');
        if (ClassRoom::where('school_id', $school->id)->where('name', $name)->where('status', 'active')->exists()) {
            return response()->json(['message' => '班级「' . $name . '」已存在'], 409);
        }

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

    /**
     * 批量创建班级
     */
    public function batchCreateClasses(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'grade' => 'required|string|max:50',
            'count' => 'required|integer|min:1|max:20',
            'year' => 'nullable|max:20',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $grade = $request->input('grade');
        $count = $request->input('count');
        $year = $request->input('year');

        // 获取该年级现有班级数量，用于确定起始序号
        $existingCount = ClassRoom::where('school_id', $school->id)
            ->where('grade', $grade)
            ->count();

        $created = [];
        for ($i = 1; $i <= $count; $i++) {
            $classNum = $existingCount + $i;
            $name = $grade . '（' . $classNum . '）班';

            // 跳过已存在的同名班级
            if (ClassRoom::where('school_id', $school->id)->where('name', $name)->where('status', 'active')->exists()) {
                continue;
            }

            $class = ClassRoom::create([
                'school_id' => $school->id,
                'name' => $name,
                'grade' => $grade,
                'year' => $year,
                'status' => 'active',
            ]);
            $created[] = ['id' => $class->id, 'name' => $class->name];
        }

        return response()->json(['message' => '已批量创建 ' . $count . ' 个班级', 'data' => $created], 201);
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
            'pet_series' => 'nullable|string|in:cosmic,pokemon,cute,treasure,mythic,all',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $class->fill($request->only(['name', 'grade', 'year', 'teacher_id', 'max_students']));

        // 宠物系列设置（存入 settings JSON）
        if ($request->has('pet_series')) {
            $settings = $class->settings ?? [];
            $settings['pet_series'] = $request->input('pet_series');
            $class->settings = $settings;
        }

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
            'role' => 'nullable|string|in:head_teacher,co_teacher,subject_teacher,grade_lead,admin_director',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误'], 422);
        }

        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($id);
        $teacherId = (int) $request->input('teacher_id');
        $role = $request->input('role', 'subject_teacher');

        User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->findOrFail($teacherId);

        if ($role === 'head_teacher') {
            $class->teacher_id = $teacherId;
            $class->save();
        }

        ClassRoomTeacher::updateOrCreate(
            ['class_room_id' => $class->id, 'user_id' => $teacherId],
            ['role' => $role],
        );

        return response()->json(['message' => '教师已分配', 'data' => $class->fresh()]);
    }

    /**
     * 移除班级教师关联
     */
    public function removeClassTeacher(Request $request, int $id): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'teacher_id' => 'required|integer|exists:users,id',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误'], 422);
        }

        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($id);

        ClassRoomTeacher::where('class_room_id', $class->id)
            ->where('user_id', (int) $request->input('teacher_id'))
            ->delete();

        if ($class->teacher_id === (int) $request->input('teacher_id')) {
            $class->teacher_id = null;
            $class->save();
        }

        return response()->json(['message' => '教师已从班级移除']);
    }

    /**
     * 学生批量导入（按模板班级名称匹配）
     * 模板字段：姓名（必填）、班级（必填）、性别（必填）、学号（选填）、手机号（选填）
     */
    public function importStudents(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'students' => 'required|array|min:1',
            'students.*.name' => 'required|string|max:50',
            'students.*.class_name' => 'required|string|max:50',
            'students.*.gender' => 'required|string|in:男,女,男生,女生,未知',
            'students.*.student_no' => 'nullable|string|max:50',
            'students.*.phone' => 'nullable|string|max:30',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误：' . $validator->errors()->first()], 422);
        }

        $school = $request->user()->school;

        $classes = ClassRoom::where('school_id', $school->id)
            ->where('status', 'active')
            ->pluck('id', 'name')
            ->all();

        $created = [];
        $errors = [];
        foreach ($request->input('students') as $idx => $row) {
            $className = trim($row['class_name']);
            if (!isset($classes[$className])) {
                $errors[] = '第 ' . ($idx + 1) . ' 行：班级「' . $className . '」不存在';
                continue;
            }

            $gender = trim($row['gender']);
            if (in_array($gender, ['男生', '男'], true)) {
                $gender = '男';
            } elseif (in_array($gender, ['女生', '女'], true)) {
                $gender = '女';
            } else {
                $gender = '未知';
            }

            $student = Student::create([
                'class_id' => $classes[$className],
                'name' => $row['name'],
                'gender' => $gender,
                'student_no' => $row['student_no'] ?? null,
                'total_score' => 0,
                'status' => 'active',
            ]);
            $created[] = ['id' => $student->id, 'name' => $student->name, 'class_name' => $className, 'gender' => $student->gender];
        }

        if ($errors) {
            return response()->json([
                'message' => '部分导入失败',
                'errors' => $errors,
                'data' => ['created' => $created, 'created_count' => count($created)],
            ], 422);
        }

        return response()->json([
            'message' => '导入完成',
            'data' => ['created_count' => count($created), 'created' => $created],
        ]);
    }

    // ===== 学生管理 =====

    /**
     * 学生列表（支持搜索、按班级/年级筛选）
     */
    public function listStudents(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        $query = Student::whereHas('classRoom', function ($q) use ($school) {
            $q->where('school_id', $school->id);
        })->with('classRoom:id,name,grade');

        // 搜索（姓名或学号）
        if ($search = $request->input('search')) {
            $query->where(function ($q) use ($search) {
                $q->where('name', 'like', "%{$search}%")
                  ->orWhere('student_no', 'like', "%{$search}%");
            });
        }

        // 按班级筛选
        if ($classId = $request->input('class_id')) {
            $query->where('class_id', $classId);
        }

        // 按年级筛选
        if ($grade = $request->input('grade')) {
            $query->whereHas('classRoom', function ($q) use ($grade) {
                $q->where('grade', $grade);
            });
        }

        // 状态筛选（默认不显示已删除的）
        $status = $request->input('status', 'active');
        if ($status !== 'all') {
            $query->where('status', $status);
        }

        $perPage = (int) ($request->input('per_page') ?? 50);
        $students = $query->orderBy('id', 'desc')->paginate($perPage);

        $data = $students->map(function ($s) {
            $arr = $s->toArray();
            $arr['class_name'] = $s->classRoom->name ?? null;
            $arr['class_grade'] = $s->classRoom->grade ?? null;

            return $arr;
        })->values();

        return response()->json([
            'data' => $data,
            'meta' => [
                'current_page' => $students->currentPage(),
                'last_page' => $students->lastPage(),
                'total' => $students->total(),
                'per_page' => $students->perPage(),
            ],
        ]);
    }

    /**
     * 单个添加学生
     */
    public function createStudent(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'name' => 'required|string|max:50',
            'class_id' => 'required|integer|exists:class_rooms,id',
            'gender' => 'nullable|string|in:男,女,男生,女生,未知',
            'student_no' => 'nullable|string|max:50',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $class = ClassRoom::where('school_id', $school->id)->findOrFail($request->input('class_id'));

        $gender = $request->input('gender');
        if (in_array($gender, ['男生', '男'], true)) {
            $gender = '男';
        } elseif (in_array($gender, ['女生', '女'], true)) {
            $gender = '女';
        } else {
            $gender = '未知';
        }

        $student = Student::create([
            'class_id' => $class->id,
            'name' => $request->input('name'),
            'gender' => $gender,
            'student_no' => $request->input('student_no'),
            'total_score' => 0,
            'status' => 'active',
        ]);

        return response()->json([
            'message' => '学生「' . $student->name . '」已添加',
            'data' => $student,
        ], 201);
    }

    /**
     * 更新学生信息
     */
    public function updateStudent(Request $request, int $id): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'name' => 'sometimes|required|string|max:50',
            'class_id' => 'sometimes|required|integer|exists:class_rooms,id',
            'gender' => 'nullable|string|in:男,女,男生,女生,未知',
            'student_no' => 'nullable|string|max:50',
            'status' => 'sometimes|in:active,graduated,inactive',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);

        $data = $request->only(['name', 'class_id', 'gender', 'student_no', 'status']);
        // 验证新班级也属于该学校
        if (isset($data['class_id'])) {
            ClassRoom::where('school_id', $school->id)->findOrFail($data['class_id']);
        }
        if (isset($data['gender'])) {
            if (in_array($data['gender'], ['男生', '男'], true)) {
                $data['gender'] = '男';
            } elseif (in_array($data['gender'], ['女生', '女'], true)) {
                $data['gender'] = '女';
            } else {
                $data['gender'] = '未知';
            }
        }

        $student->fill($data);
        $student->save();

        return response()->json(['message' => '学生信息已更新', 'data' => $student->fresh()]);
    }

    /**
     * 删除学生（软删除）
     */
    public function deleteStudent(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);

        $student->delete();

        return response()->json(['message' => '学生「' . $student->name . '」已删除']);
    }

    /**
     * 批量删除学生
     */
    public function batchDeleteStudents(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'student_ids' => 'required|array|min:1',
            'student_ids.*' => 'required|integer',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');
        $count = Student::whereIn('id', $request->input('student_ids'))
            ->whereIn('class_id', $classIds)
            ->delete();

        return response()->json(['message' => '已删除 ' . $count . ' 名学生']);
    }

    /**
     * 批量移动学生到其他班级
     */
    public function batchMoveStudents(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'student_ids' => 'required|array|min:1',
            'student_ids.*' => 'required|integer',
            'target_class_id' => 'required|integer|exists:class_rooms,id',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;
        $targetClass = ClassRoom::where('school_id', $school->id)->findOrFail($request->input('target_class_id'));
        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');

        $count = Student::whereIn('id', $request->input('student_ids'))
            ->whereIn('class_id', $classIds)
            ->update(['class_id' => $targetClass->id]);

        return response()->json([
            'message' => '已将 ' . $count . ' 名学生移动到「' . $targetClass->name . '」',
        ]);
    }

    // ===== 学年升级 =====

    /**
     * 预览学年升级（dry-run）
     */
    public function previewGradeUpgrade(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        $gradeMap = [
            '一年级' => '二年级',
            '二年级' => '三年级',
            '三年级' => '四年级',
            '四年级' => '五年级',
            '五年级' => '六年级',
        ];

        $classes = ClassRoom::where('school_id', $school->id)
            ->where('status', 'active')
            ->withCount('students')
            ->get();

        $upgrade = [];
        $graduate = [];

        foreach ($classes as $class) {
            if ($class->grade === '六年级') {
                $graduate[] = [
                    'class_id' => $class->id,
                    'class_name' => $class->name,
                    'student_count' => $class->students_count,
                ];
            } elseif (isset($gradeMap[$class->grade])) {
                $newGrade = $gradeMap[$class->grade];
                // 生成新班级名：替换年级前缀
                $newName = str_replace($class->grade, $newGrade, $class->name);
                $upgrade[] = [
                    'class_id' => $class->id,
                    'class_name' => $class->name,
                    'new_name' => $newName,
                    'old_grade' => $class->grade,
                    'new_grade' => $newGrade,
                    'student_count' => $class->students_count,
                ];
            }
        }

        $graduateStudentCount = array_sum(array_column($graduate, 'student_count'));
        $upgradeStudentCount = array_sum(array_column($upgrade, 'student_count'));

        return response()->json([
            'data' => [
                'upgrade_classes' => $upgrade,
                'graduate_classes' => $graduate,
                'summary' => [
                    'upgrade_class_count' => count($upgrade),
                    'graduate_class_count' => count($graduate),
                    'upgrade_student_count' => $upgradeStudentCount,
                    'graduate_student_count' => $graduateStudentCount,
                    'note' => '六年级学生将标记为毕业，二至五年级学生随班级升级到下一年级。一年级新生需在升级后手动创建班级并导入。',
                ],
            ],
        ]);
    }

    /**
     * 执行学年升级
     */
    public function executeGradeUpgrade(Request $request): JsonResponse
    {
        $school = $request->user()->school;
        $gradeMap = [
            '一年级' => '二年级',
            '二年级' => '三年级',
            '三年级' => '四年级',
            '四年级' => '五年级',
            '五年级' => '六年级',
        ];

        $classes = ClassRoom::where('school_id', $school->id)
            ->where('status', 'active')
            ->get();

        $upgraded = 0;
        $graduatedStudents = 0;
        $archivedClasses = 0;

        \DB::transaction(function () use ($classes, $gradeMap, &$upgraded, &$graduatedStudents, &$archivedClasses) {
            // 1. 六年级：标记学生毕业，归档班级
            foreach ($classes as $class) {
                if ($class->grade === '六年级') {
                    // 标记学生为毕业
                    $count = Student::where('class_id', $class->id)
                        ->where('status', 'active')
                        ->update(['status' => 'graduated']);
                    $graduatedStudents += $count;

                    // 归档班级
                    $class->status = 'archived';
                    $class->save();
                    $archivedClasses++;
                }
            }

            // 2. 其他年级：升级年级并重命名
            foreach ($classes as $class) {
                if (isset($gradeMap[$class->grade])) {
                    $newGrade = $gradeMap[$class->grade];
                    $newName = str_replace($class->grade, $newGrade, $class->name);

                    $class->grade = $newGrade;
                    $class->name = $newName;
                    $class->save();
                    $upgraded++;
                }
            }
        });

        return response()->json([
            'message' => '学年升级完成',
            'data' => [
                'upgraded_classes' => $upgraded,
                'archived_classes' => $archivedClasses,
                'graduated_students' => $graduatedStudents,
                'note' => '六年级学生已标记为毕业。请创建新一年级班级并导入新生名单。',
            ],
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
                $classIds = $classes->pluck('id');
                $totalScore = (int) \App\Models\Score::whereIn('class_id', $classIds)->sum('amount');
                $studentCount = $classes->sum('students_count');

                return [
                    'grade' => $grade ?: '未分年级',
                    'class_count' => $classes->count(),
                    'student_count' => $studentCount,
                    'avg_score' => $studentCount > 0 ? round($totalScore / $studentCount, 1) : 0.0,
                    'total_score' => $totalScore,
                ];
            })->values();

        return response()->json(['data' => $rows]);
    }

    public function reportsByClass(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        /** @var \Illuminate\Database\Eloquent\Collection<int, ClassRoom> $classes */
        $classes = ClassRoom::where('school_id', $school->id)
            ->with(['teacher'])
            ->withCount('students')
            ->get();

        $rows = [];
        foreach ($classes as $c) {
            $classScore = (int) Score::where('class_id', $c->id)
                ->whereMonth('created_at', Carbon::now()->month)
                ->sum('amount');

            $rows[] = [
                'class_id' => $c->id,
                'class_name' => $c->name,
                'grade' => $c->grade,
                'teacher_name' => $c->teacher?->name,
                'student_count' => $c->students_count,
                'monthly_score' => $classScore,
            ];
        }

        return response()->json(['data' => $rows]);
    }

    // ============================================================
    // 汇率管理
    // ============================================================

    public function listExchangeRates(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        $rates = \App\Models\ExchangeRate::where('school_id', $school->id)
            ->orderBy('from_currency')->orderBy('to_currency')->get();

        return response()->json(['data' => $rates]);
    }

    public function createExchangeRate(Request $request): JsonResponse
    {
        $school = $request->user()->school;

        $request->validate([
            'from_currency' => 'required|string|in:score,science,reading,class_point',
            'to_currency' => 'required|string|in:score,science,reading,class_point',
            'rate' => 'required|numeric|min:0.01',
            'is_active' => 'boolean',
        ]);

        $rate = \App\Models\ExchangeRate::create([
            'school_id' => $school->id,
            'from_currency' => $request->input('from_currency'),
            'to_currency' => $request->input('to_currency'),
            'rate' => $request->input('rate'),
            'is_active' => $request->boolean('is_active', true),
        ]);

        return response()->json(['message' => '汇率已创建', 'data' => $rate], 201);
    }

    public function updateExchangeRate(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;
        $rate = \App\Models\ExchangeRate::where('school_id', $school->id)->findOrFail($id);

        $request->validate([
            'rate' => 'sometimes|numeric|min:0.01',
            'is_active' => 'sometimes|boolean',
        ]);

        $rate->update($request->only(['rate', 'is_active']));

        return response()->json(['message' => '汇率已更新', 'data' => $rate->fresh()]);
    }

    // ===== 教师导入模板下载 =====

    /**
     * 下载教师批量导入的 CSV 模板
     */
    public function downloadTeacherTemplate(): \Symfony\Component\HttpFoundation\StreamedResponse
    {
        $headers = [
            'Content-Type' => 'text/csv; charset=UTF-8',
            'Content-Disposition' => 'attachment; filename="teacher_import_template.csv"',
        ];

        return response()->streamDownload(function () {
            $fp = fopen('php://output', 'w');

            // BOM for Excel UTF-8 compatibility
            fwrite($fp, chr(0xEF) . chr(0xBB) . chr(0xBF));

            fputcsv($fp, ['姓名', '年级团队', '科目', '密码', '手机号']);
            fputcsv($fp, ['name', 'grade_team', 'subject', 'password', 'phone']);
            // 密码默认为 star123456，不填亦可
            fputcsv($fp, ['张老师', '三年级团队', '语文', 'star123456', '13800138000']);
            fputcsv($fp, ['李老师', '三年级团队', '数学', '', '']);

            fclose($fp);
        }, 'teacher_import_template.csv', $headers);
    }

    // ===== CSV 批量导入教师 =====

    /**
     * CSV/Excel 批量导入教师（预览 + 创建）
     */
    public function importTeachers(Request $request): JsonResponse
    {
        $validator = Validator::make($request->all(), [
            'file' => 'required|file|mimes:csv,txt,xlsx,xls|max:10240',
            'dry_run' => 'boolean',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '文件上传失败', 'errors' => $validator->errors()], 422);
        }

        $school = $request->user()->school;

        if (!$school instanceof \App\Models\School) {
            return response()->json(['message' => '未找到学校'], 404);
        }

        $file = $request->file('file');
        $dryRun = $request->boolean('dry_run', true);

        $rows = $this->parseTeacherFile($file);

        $preview = [];
        $errors = [];

        foreach ($rows as $idx => $row) {
            $name = trim($row['name'] ?? $row['姓名'] ?? '');
            if ($name === '') {
                continue;
            }

            $teacherData = [
                'name'       => $name,
                'grade_team' => trim($row['grade_team'] ?? $row['年级团队'] ?? $row['所属年级团队'] ?? ''),
                'subject'    => trim($row['subject'] ?? $row['科目'] ?? ''),
                'password'   => trim($row['password'] ?? $row['密码'] ?? ''),
                'phone'      => trim($row['phone'] ?? $row['手机号'] ?? ''),
            ];

            $preview[] = $teacherData;
        }

        if (!$dryRun) {
            $created = [];
            foreach ($preview as $teacherData) {
                $teacherData['password'] = $teacherData['password'] ?: 'star123456';
                $result = $this->authService->createTeacherAccounts($school, [$teacherData]);
                $created[] = $result[0] ?? null;
            }

            return response()->json([
                'message' => '已导入 ' . count($created) . ' 名教师',
                'total' => count($created),
                'created' => $created,
            ]);
        }

        return response()->json([
            'message' => '预览模式：共 ' . count($preview) . ' 条数据',
            'total' => count($preview),
            'preview' => $preview,
        ]);
    }

    private function parseTeacherFile(mixed $file): array
    {
        $ext = strtolower($file->getClientOriginalExtension());

        if (in_array($ext, ['xlsx', 'xls'])) {
            return $this->parseExcelFile($file->getPathname());
        }

        $path = $file->getPathname();
        $raw = file_get_contents($path);

        if ($raw === false || $raw === '') {
            return [];
        }

        $enc = mb_detect_encoding($raw, ['UTF-8', 'GBK', 'GB2312', 'BIG5'], true);
        if ($enc && $enc !== 'UTF-8') {
            $raw = mb_convert_encoding($raw, 'UTF-8', $enc);
        }

        $bom = chr(0xEF) . chr(0xBB) . chr(0xBF);
        if (str_starts_with($raw, $bom)) {
            $raw = substr($raw, 3);
        }

        $lines = [];
        $fh = @fopen($path, 'r');
        if ($fh !== false) {
            while (($line = fgets($fh)) !== false) {
                $line = trim($line);
                if ($line !== '') {
                    $lines[] = $line;
                }
            }
            fclose($fh);
        }

        if (empty($lines)) {
            return [];
        }

        $delimiter = $this->detectCsvDelimiter($lines[0]);
        $rows = [];
        foreach ($lines as $line) {
            $cols = str_getcsv($line, $delimiter);
            $rows[] = $cols;
        }

        $header = array_shift($rows);
        $result = [];
        foreach ($rows as $cols) {
            $row = [];
            foreach ($header as $i => $key) {
                $row[$key] = $cols[$i] ?? '';
            }
            $result[] = $row;
        }

        return $result;
    }

    private function detectCsvDelimiter(string $line): string
    {
        $candidates = [chr(9), ',', ';'];
        $best = ',';
        $bestCount = 0;

        foreach ($candidates as $d) {
            $count = substr_count($line, $d);
            if ($count > $bestCount) {
                $bestCount = $count;
                $best = $d;
            }
        }

        return $best;
    }

    private function parseExcelFile(string $path): array
    {
        if (!class_exists('PhpOffice\PhpSpreadsheet\IOFactory')) {
            return [];
        }

        $spreadsheet = \PhpOffice\PhpSpreadsheet\IOFactory::load($path);
        $worksheet = $spreadsheet->getActiveSheet();
        $data = $worksheet->toArray();

        if (empty($data)) {
            return [];
        }

        $header = array_shift($data);
        $result = [];
        foreach ($data as $row) {
            $item = [];
            foreach ($header as $i => $key) {
                $item[$key] = $row[$i] ?? '';
            }
            $result[] = $item;
        }

        return $result;
    }

    // ===== 批量教师班级分配 =====

    /**
     * 批量将教师分配到多个班级
     */
    public function assignTeacherClasses(Request $request, int $id): JsonResponse
    {
        $school = $request->user()->school;

        $teacher = User::where('school_id', $school->id)
            ->where('role', 'teacher')
            ->findOrFail($id);

        $validator = Validator::make($request->all(), [
            'assignments' => 'required|array|min:1',
            'assignments.*.class_id' => 'required|integer|exists:class_rooms,id',
            'assignments.*.role' => 'required|string|in:head_teacher,co_teacher,subject_teacher,grade_lead,admin_director',
            'assignments.*.subject' => 'nullable|string|max:50',
        ]);

        if ($validator->fails()) {
            return response()->json(['message' => '参数错误', 'errors' => $validator->errors()], 422);
        }

        $reqAssignments = $request->input('assignments');
        $synced = [];

        foreach ($reqAssignments as $item) {
            $classId = (int) $item['class_id'];
            $role = $item['role'];

            $classRoom = ClassRoom::where('school_id', $school->id)->findOrFail($classId);

            $updateData = ['role' => $role];
            if (isset($item['subject'])) {
                $updateData['subject'] = $item['subject'];
            }

            $assignment = ClassRoomTeacher::updateOrCreate(
                ['class_room_id' => $classId, 'user_id' => $teacher->id],
                $updateData,
            );

            if ($role === 'head_teacher') {
                $classRoom->teacher_id = $teacher->id;
                $classRoom->save();
            }

            $synced[] = [
                'class_id' => $classId,
                'class_name' => $classRoom->name,
                'role' => $assignment->role,
                'subject' => $assignment->subject,
            ];
        }

        return response()->json([
            'message' => '已为教师「' . $teacher->name . '」分配 ' . count($synced) . ' 个班级',
            'data' => [
                'teacher_id' => $teacher->id,
                'teacher_name' => $teacher->name,
                'assignments' => $synced,
            ],
        ]);
    }
}
