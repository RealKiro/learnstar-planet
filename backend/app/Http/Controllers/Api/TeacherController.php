<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use App\Models\Score;
use App\Models\ScoreRule;
use App\Models\Student;
use App\Models\User;
use App\Services\LeaderboardService;
use App\Services\ScoreService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class TeacherController extends Controller
{
    public function __construct(
        private readonly ScoreService $scoreService,
        private readonly LeaderboardService $leaderboardService,
    ) {}

    // ============================================================
    // Dashboard
    // ============================================================

    public function dashboard(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        if ($classIds->isEmpty()) {
            return response()->json(['data' => [
                'student_count' => 0,
                'weekly_score' => 0,
                'avg_pet_level' => 0,
                'pending_redemptions' => 0,
                'recent_scores' => [],
                'top_students' => [],
            ]]);
        }

        $studentCount = Student::whereIn('class_id', $classIds)->where('status', 'active')->count();
        $weeklyScores = Score::whereIn('class_id', $classIds)
            ->where('created_at', '>=', now()->startOfWeek())
            ->sum('amount');
        $avgPetLevel = (int) Student::whereIn('class_id', $classIds)
            ->whereHas('pet')
            ->withAvg('pet', 'level')
            ->get()
            ->avg('pet_avg_level');
        $pendingRedemptions = \App\Models\ShopRedemption::whereIn('class_id', $classIds)
            ->where('status', 'pending')
            ->count();
        $recentScores = Score::whereIn('class_id', $classIds)
            ->with('student')
            ->orderBy('created_at', 'desc')
            ->take(10)
            ->get()
            ->map(fn (Score $s) => [
                'student_name' => $s->student?->name,
                'points' => $s->amount,
                'reason' => $s->reason,
                'created_at' => $s->created_at?->diffForHumans(),
            ]);
        $topStudents = Student::whereIn('class_id', $classIds)
            ->orderBy('total_score', 'desc')
            ->take(5)
            ->get()
            ->map(fn (Student $s) => [
                'student_name' => $s->name,
                'score' => $s->total_score,
            ]);

        return response()->json(['data' => [
            'student_count' => $studentCount,
            'weekly_score' => (int) $weeklyScores,
            'avg_pet_level' => $avgPetLevel,
            'pending_redemptions' => $pendingRedemptions,
            'recent_scores' => $recentScores,
            'top_students' => $topStudents,
        ]]);
    }

    // ============================================================
    // Student Management
    // ============================================================

    public function listStudents(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $query = Student::whereIn('class_id', $classIds)
            ->with('classRoom:id,name,grade');

        if ($request->has('search')) {
            $search = $request->input('search');
            $query->where(function ($q) use ($search) {
                $q->where('name', 'like', "%{$search}%")
                  ->orWhere('student_no', 'like', "%{$search}%");
            });
        }

        $students = $query->orderBy('name')->paginate(50);

        return response()->json([
            'data' => $students->items(),
            'meta' => [
                'current_page' => $students->currentPage(),
                'last_page' => $students->lastPage(),
                'per_page' => $students->perPage(),
                'total' => $students->total(),
            ],
        ]);
    }

    public function importStudents(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $students = $request->input('students', []);
        $imported = 0;

        foreach ($students as $data) {
            if (empty($data['name']) || empty($data['class_name'])) {
                continue;
            }
            $classRoom = ClassRoom::whereIn('id', $classIds)
                ->where('name', $data['class_name'])
                ->first();
            if (!$classRoom) {
                continue;
            }
            Student::create([
                'class_id' => $classRoom->id,
                'name' => $data['name'],
                'gender' => $data['gender'] ?? '未知',
                'student_no' => $data['student_no'] ?? null,
                'total_score' => 0,
                'status' => 'active',
            ]);
            $imported++;
        }

        return response()->json([
            'message' => "成功导入 {$imported} 名学生",
            'data' => ['imported_count' => $imported],
        ]);
    }

    public function updateStudent(Request $request, int $id): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);

        $student->update($request->only(['name', 'gender', 'student_no']));

        return response()->json(['message' => '更新成功', 'data' => $student]);
    }

    // ============================================================
    // Score Management
    // ============================================================

    public function scoreSummary(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $total = Score::whereIn('class_id', $classIds)->sum('amount');
        $today = Score::whereIn('class_id', $classIds)
            ->whereDate('created_at', today())
            ->sum('amount');
        $week = Score::whereIn('class_id', $classIds)
            ->where('created_at', '>=', now()->startOfWeek())
            ->sum('amount');

        return response()->json(['data' => [
            'total' => (int) $total,
            'today' => (int) $today,
            'this_week' => (int) $week,
        ]]);
    }

    public function giveScore(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'student_id' => 'required|integer',
            'points' => 'required|integer|not_in:0',
            'reason' => 'required|string|max:200',
        ]);

        $student = Student::whereIn('class_id', $classIds)
            ->with('pet')
            ->findOrFail($request->input('student_id'));

        $amount = (int) $request->input('points');
        $score = $this->scoreService->giveScore(
            $student,
            $amount,
            $request->input('reason'),
            $teacher->id,
        );

        if ($amount > 0 && $student->pet) {
            $this->leaderboardService->updateTotalScore($student->class_id, $student->id, $student->total_score);
        }

        return response()->json([
            'message' => ($amount > 0 ? '加分' : '减分') . '成功',
            'data' => [
                'student_id' => $student->id,
                'student_name' => $student->name,
                'points' => $amount,
                'new_score' => $student->fresh()->total_score,
            ],
        ]);
    }

    public function batchGiveScore(Request $request): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');

        $request->validate([
            'student_ids' => 'required|array|min:1|max:100',
            'student_ids.*' => 'integer',
            'points' => 'required|integer|not_in:0',
            'reason' => 'required|string|max:200',
        ]);

        $amount = (int) $request->input('points');
        $results = $this->scoreService->batchGiveScore(
            $request->input('student_ids'),
            $amount,
            $request->input('reason'),
            $teacher->id,
        );

        return response()->json([
            'message' => "批量操作完成，处理了 " . count($results) . " 名学生",
            'data' => ['count' => count($results)],
        ]);
    }

    public function giveScoreByRule(Request $request, int $ruleId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $rule = ScoreRule::findOrFail($ruleId);

        $request->validate([
            'student_id' => 'required|integer',
        ]);

        $student = Student::whereIn('class_id', $classIds)->findOrFail($request->input('student_id'));

        $this->scoreService->giveScoreByRule($student, $rule, $teacher->id);

        return response()->json([
            'message' => "已按规则「{$rule->name}」处理",
            'data' => [
                'rule' => $rule->name,
                'points' => $rule->amount,
                'student_name' => $student->name,
            ],
        ]);
    }

    public function scoreHistory(Request $request, int $studentId): JsonResponse
    {
        $teacher = $request->user();
        $classIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id');
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $history = $this->scoreService->getScoreHistory($student, 20);

        return response()->json([
            'data' => $history->items(),
            'meta' => [
                'current_page' => $history->currentPage(),
      