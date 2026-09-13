<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use App\Models\Student;
use App\Models\User;
use Illuminate\Support\Facades\Cache;

/**
 * 教师端班级宠物系列：班级信息卡（含展示码）与整班系列切换（发放免费自选机会）。
 *
 * 只负责数据操作；HTTP 校验、响应包装留在控制器。
 */
class PetSeriesService
{
    /** 整班切换系列后，每人可免费自选一次的有效期（天） */
    private const FREE_PICK_DAYS = 3;

    public function __construct(
        private readonly TeacherClassScope $scope,
    ) {
    }

    /**
     * 当前班级信息卡（系列设置页顶部）。
     *
     * @return array<string, mixed>
     *
     * @throws \DomainException 没有可管理的班级（400）
     */
    public function classInfo(User $teacher): array
    {
        $classIds = $this->scope->ids($teacher);

        if ($classIds->isEmpty()) {
            throw new \DomainException('没有可管理的班级', 400);
        }

        $classId = $classIds->first();
        $class = ClassRoom::findOrFail($classId);

        $totalScore = Student::where('class_id', $classId)
            ->where('status', 'active')
            ->sum('total_score');

        return [
            'id' => $class->id,
            'name' => $class->name,
            'grade' => $class->grade,
            'student_count' => Student::where('class_id', $classId)->where('status', 'active')->count(),
            'total_score' => (int) $totalScore,
            'class_points' => (int) ($class->settings['class_points'] ?? 0),
            'settings' => $class->settings,
            'display_code' => DisplayCodeService::generate($class),
        ];
    }

    /**
     * 整班切换宠物系列：写班级设置，并给全班学生各发放一次 3 天内有效的免费自选机会。
     *
     * @return array<string, mixed>
     *
     * @throws \DomainException 没有可管理的班级（400）
     */
    public function switchSeries(User $teacher, string $seriesId): array
    {
        $classIds = $this->scope->ids($teacher);

        if ($classIds->isEmpty()) {
            throw new \DomainException('没有可管理的班级', 400);
        }

        $classId = $classIds->first();
        $class = ClassRoom::findOrFail($classId);
        $settings = $class->settings ?? [];
        $settings['pet_series'] = $seriesId;
        $class->settings = $settings;
        $class->save();

        // 发放「免费自选」机会：整班切换后 3 天内每人可免费切换一次当前类别的宠物，过期作废
        $students = Student::where('class_id', $classId)
            ->where('status', 'active')
            ->get();
        foreach ($students as $student) {
            Cache::put("pet_free_pick:{$student->id}", 1, now()->addDays(self::FREE_PICK_DAYS));
        }

        return [
            'series_id' => $seriesId,
            'class_id' => $classId,
            'free_pick_granted' => true,
            'granted_students' => $students->count(),
            'student_count' => $students->count(),
        ];
    }
}
