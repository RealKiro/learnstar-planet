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
    protected $signature = 'demo:clean';

    protected $description = '清除所有演示数据';

    public function handle(): int
    {
        $school = School::where('code', 'DEMO')->first();
        if (!$school) {
            $this->info('未找到演示数据');

            return 0;
        }

        if (!$this->confirm('确定要清除所有演示数据吗？此操作不可撤销。')) {
            $this->info('已取消');

            return 0;
        }

        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');
        $studentIds = Student::whereIn('class_id', $classIds)->pluck('id');

        Score::whereIn('student_id', $studentIds)->delete();
        Pet::whereIn('student_id', $studentIds)->delete();
        Student::whereIn('class_id', $classIds)->delete();
        ClassRoom::where('school_id', $school->id)->delete();
        User::where('school_id', $school->id)->delete();
        $school->delete();

        $this->info('🗑️  所有演示数据已清除');

        return 0;
    }
}
