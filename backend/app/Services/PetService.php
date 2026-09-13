<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use App\Models\Pet;
use App\Models\PetCollection;
use App\Models\Student;
use App\Models\User;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Cache;

/**
 * 教师端宠物域：班级宠物总览、详情、喂养、更名、切换（含免费自选与图鉴恢复）、收藏图鉴。
 *
 * 只负责数据操作与业务守卫；HTTP 校验、响应包装与状态码映射留在控制器。
 * 守卫失败抛 DomainException，code 即 HTTP 状态码（422 业务守卫 / 400 资源不足）。
 * 学生归属查询经 TeacherClassScope 限定；switchPet 沿用原 getAccessibleClassIds 全放行口径。
 */
class PetService
{
    public function __construct(
        private readonly TeacherClassScope $scope,
        private readonly LeaderboardService $leaderboardService,
        private readonly DisplayEventService $displayEventService,
    ) {
    }

    /**
     * 班级宠物总览（按教师可管理班级聚合）。
     */
    public function classPetsOverview(User $teacher): Collection
    {
        $classIds = $this->scope->ids($teacher);

        return Pet::whereIn('class_id', $classIds)
            ->with('student:id,name')
            ->get()
            /** @phpstan-ignore-next-line argument.unresolvableType */
            ->map(fn (\App\Models\Pet $p) => [
                'id' => $p->id,
                'student_id' => $p->student_id,
                'student_name' => $p->student?->name,
                'name' => $p->name,
                'species' => $p->species ?: 'zhulong',
                'level' => $p->level,
                'exp' => $p->exp,
                'mood' => $p->mood,
                'stage_name' => $p->currentStage()['name'],
            ]);
    }

    /**
     * 学生宠物详情。
     *
     * @return array<string, mixed>
     *
     * @throws \DomainException 学生不在管辖班级（404）或尚无宠物（404）
     */
    public function petFor(User $teacher, int $studentId): array
    {
        $classIds = $this->scope->ids($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            throw new \DomainException('该学生还没有宠物', 404);
        }

        $stage = $pet->currentStage();

        return [
            'id' => $pet->id,
            'name' => $pet->name,
            'species' => $pet->species ?: 'zhulong',
            'level' => $pet->level,
            'exp' => $pet->exp,
            'mood' => $pet->mood,
            'emoji' => $stage['emoji'],
            'stage_name' => $stage['name'],
            'last_fed_at' => $pet->last_fed_at?->toDateTimeString(),
        ];
    }

    /**
     * 喂养宠物：更新状态、刷新排行榜、推送班级大屏。
     *
     * @return array{message: string, mood: int, level: int}
     *
     * @throws \DomainException 学生不在管辖班级（404）或尚无宠物（404）
     */
    public function feed(User $teacher, int $studentId): array
    {
        $classIds = $this->scope->ids($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            throw new \DomainException('该学生还没有宠物', 404);
        }

        $pet->feed();
        $this->leaderboardService->updatePetLevel($student->class_id, $student->id, $pet->level);

        // 推送给班级大屏
        try {
            $this->displayEventService->publish($student->class_id, 'pet_update', [
                'student_id' => $student->id,
                'student_name' => $student->name,
                'type' => 'feed',
                'mood' => $pet->mood,
                'level' => $pet->level,
                'experience' => $pet->experience,
            ]);
        } catch (\Throwable $e) {
            logger()->warning('Display pet event failed: ' . $e->getMessage());
        }

        return [
            'message' => "已喂养「{$pet->name}」",
            'mood' => $pet->mood,
            'level' => $pet->level,
        ];
    }

    /**
     * 宠物更名。
     *
     * @throws \DomainException 学生不在管辖班级（404）或尚无宠物（404）
     */
    public function rename(User $teacher, int $studentId, string $name): string
    {
        $classIds = $this->scope->ids($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $pet = $student->pet;
        if (!$pet) {
            throw new \DomainException('该学生还没有宠物', 404);
        }

        $pet->update(['name' => $name]);

        return "宠物已更名为「{$pet->name}」";
    }

    /**
     * 切换宠物：同物种守卫、类别限制、免费自选、积分扣减、图鉴进度保存与恢复、新宠创建。
     *
     * @return array<string, mixed>
     *
     * @throws \DomainException 学生不在管辖班级（404）、同物种（422）、跨类别（422）或积分不足（400）
     */
    public function switchPet(User $teacher, int $studentId, string $petSpecies, string $petName): array
    {
        // 沿用原 switchPet 口径：getAccessibleClassIds（api-bot 全校放行）
        $classIds = $this->scope->accessibleIds($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);
        $now = now();

        $pet = $student->pet;

        // ===== 同物种切换守卫(不扣费、不重置冷却) =====
        if ($pet && $pet->species === $petSpecies) {
            throw new \DomainException('当前已经是这只宠物啦', 422);
        }

        // ===== 类别限制：只能在本班当前类别内更换，不能跨类别领养 =====
        // 注：旧系列 id(cosmic/cute/all 等)在 speciesPoolForSeries 返回空池 → 视为不限制
        $classSeries = ClassRoom::find($student->class_id)?->settings['pet_series'] ?? null;
        $seriesPool = $classSeries ? Pet::speciesPoolForSeries($classSeries) : [];
        if ($classSeries && !empty($seriesPool) && !in_array($petSpecies, $seriesPool, true)) {
            throw new \DomainException('只能领养当前类别「' . $classSeries . '」的宠物，不能跨类别领养', 422);
        }

        // ===== 免费自选：整班切换后的一次机会（限当前类别，免费） =====
        $usedFreePick = false;
        $switchCost = $pet ? Pet::switchCost($pet->level) : 0;
        if ($pet && Cache::has("pet_free_pick:{$student->id}")) {
            $usedFreePick = true;
            $switchCost = 0;
        }

        // ===== 目标物种收藏进度（切回时恢复） =====
        $collection = $pet
            ? PetCollection::where('student_id', $student->id)->where('species', $petSpecies)->first()
            : null;

        if ($pet) {
            // 1) 先扣积分（等级越高越贵；免费自选不扣）——积分不足直接拒绝，不留脏数据
            if (!$usedFreePick) {
                if ($student->total_score < $switchCost) {
                    throw new \DomainException("积分不足，更换宠物需 {$switchCost} 积分", 400);
                }
                $student->total_score -= $switchCost;
                $student->save();
            }

            // 2) 保存当前宠物进度到图鉴（进度全保留）
            PetCollection::updateOrCreate(
                ['student_id' => $student->id, 'species' => $pet->species],
                ['level' => $pet->level, 'experience' => $pet->experience, 'mood' => $pet->mood, 'is_active' => false]
            );

            // 3) 恢复目标物种进度（新物种为初始形态；进度全保留）
            if ($collection) {
                $pet->species = $petSpecies;
                $pet->level = $collection->level;
                $pet->experience = $collection->experience;
                $pet->mood = $collection->mood;
            } else {
                $pet->species = $petSpecies;
                $pet->level = 1;
                $pet->experience = 0;
                $pet->mood = 80;
            }
            $pet->name = $petName;
            $pet->last_switched_at = $now;
            $pet->save();

            // 4) 目标物种标记激活
            PetCollection::updateOrCreate(
                ['student_id' => $student->id, 'species' => $petSpecies],
                ['level' => $pet->level, 'experience' => $pet->experience, 'mood' => $pet->mood, 'is_active' => true]
            );

            // 免费自选机会使用即失效
            if ($usedFreePick) {
                Cache::forget("pet_free_pick:{$student->id}");
            }

            return [
                'message' => $usedFreePick
                    ? '✅ 已使用整班切换的免费自选机会！'
                    : '宠物已更换为「' . $petName . '」（扣除 ' . $switchCost . ' 积分）',
                'data' => [
                    'pet_name' => $pet->name,
                    'pet_species' => $pet->species,
                    'level' => $pet->level,
                    'experience' => $pet->experience,
                    'cost' => $switchCost,
                    'free_pick_used' => $usedFreePick,
                ],
            ];
        }

        // 无宠物：创建新宠物并收入图鉴
        $pet = Pet::create([
            'student_id' => $student->id,
            'class_id' => $student->class_id,
            'name' => $petName,
            'species' => $petSpecies,
            'level' => 1,
            'experience' => 0,
            'mood' => 80,
        ]);
        PetCollection::create([
            'student_id' => $student->id,
            'species' => $petSpecies,
            'level' => 1,
            'experience' => 0,
            'mood' => 80,
            'is_active' => true,
        ]);

        return [
            'message' => '已为您分配宠物「' . $petName . '」',
            'data' => [
                'pet_name' => $pet->name,
                'pet_species' => $pet->species,
                'level' => $pet->level,
                'experience' => $pet->experience,
            ],
        ];
    }

    /**
     * 学生宠物图鉴收藏。
     *
     * @return array<string, mixed>
     *
     * @throws \DomainException 学生不在管辖班级（404）
     */
    public function collection(User $teacher, int $studentId): array
    {
        $classIds = $this->scope->ids($teacher);
        $student = Student::whereIn('class_id', $classIds)->findOrFail($studentId);

        $activePet = $student->pet;

        // 确保当前激活宠物在收藏中
        if ($activePet && $activePet->species) {
            PetCollection::firstOrCreate(
                ['student_id' => $student->id, 'species' => $activePet->species],
                ['level' => $activePet->level, 'experience' => $activePet->experience, 'mood' => $activePet->mood, 'is_active' => true]
            );
        }

        $collections = PetCollection::where('student_id', $student->id)->orderBy('species')->get();

        return [
            'student_id' => $student->id,
            'student_name' => $student->name,
            'total_score' => $student->total_score,
            'unlock_slots' => PetCollection::unlockSlotsForScore($student->total_score),
            'class_series' => ClassRoom::find($student->class_id)?->settings['pet_series'] ?? null,
            'active_species' => $activePet?->species,
            'collection' => $collections->map(fn (PetCollection $c) => [
                'species' => $c->species,
                'level' => $c->level,
                'experience' => $c->experience,
                'mood' => $c->mood,
                'is_active' => $c->is_active,
            ])->values(),
        ];
    }
}
