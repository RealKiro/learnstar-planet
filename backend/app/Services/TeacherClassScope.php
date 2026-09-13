<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use App\Models\ClassRoomTeacher;
use App\Models\User;
use Illuminate\Support\Collection;

/**
 * 教师班级作用域。
 *
 * 集中承载「这名教师能看到哪些班级」这一权限判定，供控制器与各业务 Service 复用。
 *
 * 注意：API 机器人账号走 User::isApiBot()（基于 settings.is_api_bot 的免迁移标记），
 * 不要改成新的数据库字段，否则需要额外迁移。
 */
class TeacherClassScope
{
    /**
     * 教师关联的所有班级 ID（班主任 + 科任/副班等）。
     *
     * API 机器人账号返回本校全部启用中的班级（含未来新建，无需维护关联表）。
     */
    public function ids(User $teacher): Collection
    {
        // API 机器人账号：本校全部启用中的班级（含未来新建，无需维护关联表）
        if ($teacher->isApiBot()) {
            return ClassRoom::where('school_id', $teacher->school_id)
                ->where('status', 'active')
                ->pluck('id')
                ->merge(ClassRoom::where('teacher_id', $teacher->id)->pluck('id'))
                ->unique();
        }

        return ClassRoomTeacher::where('user_id', $teacher->id)
            ->pluck('class_room_id')
            ->merge(ClassRoom::where('teacher_id', $teacher->id)->pluck('id'))
            ->unique();
    }

    /**
     * 可访问班级 ID 数组（自任班级 ∪ 关联班级，去重）。
     *
     * @return array<int, int>
     */
    public function accessibleIds(User $teacher): array
    {
        $ownClassIds = ClassRoom::where('teacher_id', $teacher->id)->pluck('id')->toArray();
        $relatedClassIds = $this->ids($teacher)->toArray();

        return array_values(array_unique(array_merge($ownClassIds, $relatedClassIds)));
    }
}
