<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Services\AuthService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SchoolAdminController extends Controller
{
    public function __construct(
        private readonly AuthService $authService
    ) {}

    // ===== 学校管理 =====

    public function getSchool(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function updateSchool(Request $request): JsonResponse
    {
        return response()->json(['message' => '更新成功']);
    }

    // ===== 教师管理 =====

    public function batchCreateTeachers(Request $request): JsonResponse
    {
        $request->validate([
            'teachers' => 'required|array',
        ]);

        $school = $request->user()->school;

        $created = $this->authService->createTeacherAccounts($school, $request->input('teachers'));

        return response()->json(['data' => $created]);
    }

    public function listTeachers(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function updateTeacher(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '更新成功']);
    }

    public function resetTeacherPassword(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '密码已重置']);
    }

    public function disableTeacher(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '已禁用']);
    }

    // ===== 家长管理 =====

    public function batchCreateParents(Request $request): JsonResponse
    {
        return response()->json(['message' => '创建成功']);
    }

    public function listParents(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    // ===== 班级管理 (apiResource) =====

    public function index(): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function store(Request $request): JsonResponse
    {
        return response()->json(['message' => '创建成功'], 201);
    }

    public function show(int $id): JsonResponse
    {
        return response()->json(['data' => null]);
    }

    public function update(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '更新成功']);
    }

    public function destroy(int $id): JsonResponse
    {
        return response()->json(['message' => '删除成功']);
    }

    public function assignClassTeacher(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '分配成功']);
    }

    public function importStudents(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '导入成功']);
    }

    // ===== 报表 =====

    public function schoolOverview(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function reportsByGrade(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function reportsByClass(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }
}
