<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\ClassRoom;
use App\Models\Student;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Support\Collection;

/**
 * 教师端学生管理服务：列表查询、批量导入（查重 + 跨班冲突防护）、增改删。
 */
class StudentService
{
    /**
     * 教师所带班级学生列表（分页，附加宠物字段）。
     *
     * @param Collection<int, int> $classIds
     */
    public function list(Collection $classIds, ?string $search): LengthAwarePaginator
    {
        $query = Student::whereIn('class_id', $classIds)
            ->with('classRoom:id,name,grade');

        if ($search !== null) {
            $query->where(function ($q) use ($search) {
                $q->where('name', 'like', "%{$search}%")
                  ->orWhere('student_no', 'like', "%{$search}%");
            });
        }

        return $query->with('pet')->orderBy('name')->paginate(50);
    }

    /**
     * 批量导入学生：同班学号/同名查重跳过，跨班同学号疑似转班跳过。
     * 与历史行为一致：逐行处理、跳过项不中断整体导入（未包事务）。
     *
     * @param Collection<int, int> $classIds
     * @param array<int, array<string, mixed>> $rows
     * @return array{message: string, imported_count: int, skipped: array<int, string>}
     */
    public function import(Collection $classIds, array $rows): array
    {
        $imported = 0;
        $skipped = [];

        foreach ($rows as $data) {
            if (empty($data['name']) || empty($data['class_name'])) {
                continue;
            }
            $classRoom = ClassRoom::whereIn('id', $classIds)
                ->where('name', $data['class_name'])
                ->first();
            if (!$classRoom) {
                continue;
            }
            $name = trim((string) $data['name']);
            $studentNo = trim((string) ($data['student_no'] ?? ''));

            // 查重：同班同学号/同班同名已存在则跳过（重复导入不再重复建人）
            $dupQuery = Student::where('class_id', $classRoom->id);
            $dup = $studentNo !== ''
                ? $dupQuery->where('student_no', $studentNo)->exists()
                : $dupQuery->where('name', $name)->exists();
            if ($dup) {
                $skipped[] = $name . ($studentNo !== '' ? "（学号 {$studentNo}）" : '') . '：' . $classRoom->name . ' 已存在';
                continue;
            }
            // 跨班冲突防护：同学号已在同校其他班级 → 疑似转班，跳过并提示
            if ($studentNo !== '') {
                $crossDup = Student::with('classRoom:id,name')
                    ->where('student_no', $studentNo)
                    ->where('class_id', '!=', $classRoom->id)
                    ->where('status', 'active')
                    ->whereHas('classRoom', function ($q) use ($classRoom) {
                        $q->where('school_id', $classRoom->school_id);
                    })
                    ->first();
                if ($crossDup) {
                    $skipped[] = $name . "（学号 {$studentNo}）：已存在于 "
                        . ($crossDup->classRoom->name ?? '其他班级')
                        . '，如为转班请联系管理员使用批量转班';
                    continue;
                }
            }
            Student::create([
                'class_id' => $classRoom->id,
                'name' => $name,
                'gender' => $data['gender'] ?? '未知',
                'student_no' => $studentNo !== '' ? $studentNo : null,
                'total_score' => 0,
                'status' => 'active',
            ]);
            $imported++;
        }

        return [
            'message' => "成功导入 {$imported} 名学生" . (count($skipped) > 0 ? "，跳过 " . count($skipped) . " 条重复/冲突记录" : ''),
            'imported_count' => $imported,
            'skipped' => $skipped,
        ];
    }

    /**
     * 添加学生（仅限教师所带班级；性别归一化为 男/女/未知）。
     * 返回 null 表示目标班级不在教师管理范围（控制器映射 403）。
     *
     * @param Collection<int, int> $classIds
     * @param array{name: string, class_id: int, gender?: ?string, student_no?: ?string} $data
     */
    public function create(Collection $classIds, array $data): ?Student
    {
        $class = ClassRoom::whereIn('id', $classIds)->find($data['class_id']);
        if (!$class) {
            return null;
        }
        $gender = $data['gender'] ?? null;
        if (in_array($gender, ['男生', '男'], true)) {
            $gender = '男';
        } elseif (in_array($gender, ['女生', '女'], true)) {
            $gender = '女';
        } else {
            $gender = '未知';
        }

        return Student::create([
            'class_id' => $class->id,
            'name' => $data['name'],
            'gender' => $gender,
            'student_no' => $data['student_no'] ?? null,
            'total_score' => 0,
            'status' => 'active',
        ]);
    }

    /**
     * 更新学生（仅限教师所带班级，找不到抛 ModelNotFoundException → 404）。
     *
     * @param Collection<int, int> $classIds
     * @param array{name?: string, gender?: string, student_no?: ?string} $fields
     */
    public function update(Collection $classIds, int $id, array $fields): Student
    {
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);
        $student->update($fields);

        return $student;
    }

    /**
     * 删除学生（仅限教师所带班级），返回被删学生（供响应文案）。
     *
     * @param Collection<int, int> $classIds
     */
    public function delete(Collection $classIds, int $id): Student
    {
        $student = Student::whereIn('class_id', $classIds)->findOrFail($id);
        $student->delete();

        return $student;
    }
}
