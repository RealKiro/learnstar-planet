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
        // 第1步：清理所有 demo_ 前缀的用户（可能遗留自不同 school_id 的历史记录）
        $demoUsernames = ['demo_admin', 'demo_t1', 'demo_t2', 'demo_t3', 'demo_t4'];
        foreach ($demoUsernames as $username) {
            $user = User::withTrashed()->where('username', $username)->first();
            if ($user) {
                Score::where('given_by', $user->id)->forceDelete();
                $user->forceDelete();
            }
        }

        // 第2步：清理所有 DEMO 前缀的班级码关联的班级和学生
        $demoClasses = ClassRoom::withTrashed()->where('display_code', 'LIKE', 'DEMO%')->get();
        foreach ($demoClasses as $class) {
            $studentIds = Student::withTrashed()->where('class_id', $class->id)->pluck('id');
            Score::whereIn('student_id', $studentIds)->forceDelete();
            Pet::whereIn('student_id', $studentIds)->forceDelete();
            Student::withTrashed()->where('class_id', $class->id)->forceDelete();
            $class->forceDelete();
        }

        // 第3步：按学校清理
        $school = School::withTrashed()->where('code', 'DEMO')->first();
        if (!$school) {
            $this->info('🗑️  历史残留演示数据已清除');

            return 0;
        }

        if (!$this->option('force') && !$this->confirm('确定要清除所有演示数据吗？此操作不可撤销。')) {
            $this->info('已取消');

            return 0;
        }

        $classIds = ClassRoom::withTrashed()->where('school_id', $school->id)->pluck('id');
        $studentIds = Student::withTrashed()->whereIn('class_id', $classIds)->pluck('id');

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
