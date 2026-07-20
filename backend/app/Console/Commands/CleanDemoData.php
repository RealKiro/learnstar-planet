<?php

declare(strict_types=1);

namespace App\Console\Commands;

use App\Models\ClassRoom;
use App\Models\Pet;
use App\Models\School;
use App\Models\Score;
use App\Models\Student;
use App\Models\User;
use Illuminate\Console\Command;

class CleanDemoData extends Command
{
    protected $signature = 'demo:clean {--force : 跳过确认提示}';

    protected $description = '清除所有演示数据';

    public function handle(): int
    {
        $school = School::withTrashed()->where('code', 'DEMO')->first();
        if (!$school) {
            $this->info('未找到演示数据');

            return 0;
        }

        if (!$this->option('force') && !$this->confirm('确定要清除所有演示数据吗？此操作不可撤销。')) {
            $this->info('已取消');

            return 0;
        }

        $classIds = ClassRoom::withTrashed()->where('school_id', $school->id)->pluck('id');
        $studentIds = Student::withTrashed()->whereIn('class_id', $classIds)->pluck('id');

        // 使用 forceDelete 彻底删除，避免软删除导致的唯一索引冲突
        Score::whereIn('student_id', $studentIds)->forceDelete();
        Pet::whereIn('student_id', $studentIds)->forceDelete();
        Student::withTrashed()->whereIn('class_id', $classIds)->forceDelete();
        ClassRoom::withTrashed()->where('school_id', $school->id)->forceDelete();
        User::withTrashed()->where('school_id', $school->id)->forceDelete();
        $school->forceDelete();

        $this->info('🗑️  所有演示数据已清除');

        return 0;
    }
}
