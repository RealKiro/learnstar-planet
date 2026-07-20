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
use Illuminate\Support\Facades\Hash;

class SeedDemoData extends Command
{
    protected $signature = 'demo:seed {--force : 跳过确认提示}';

    protected $description = '生成演示数据（仅用于测试/试用，不影响正式数据）';

    public function handle(): int
    {
        if ($this->hasDemoData()) {
            if ($this->option('force')) {
                $this->cleanDemoData();
            } elseif (!$this->confirm('检测到已有演示数据，是否重新生成？这会先删除现有演示数据。')) {
                $this->info('已取消');

                return 0;
            } else {
                $this->cleanDemoData();
            }
        }

        $this->info('正在生成演示数据...');

        // 1. 创建演示学校
        $school = School::create([
            'name' => '阳光小学（演示）',
            'code' => 'DEMO',
            'address' => '演示地址',
            'settings' => json_encode(['is_demo' => true]),
        ]);
        $this->line('  ✅ 演示学校已创建');

        // 2. 创建管理员账号（演示用）
        User::create([
            'school_id' => $school->id,
            'role' => 'school_admin',
            'username' => 'demo_admin',
            'password' => Hash::make('demo123456'),
            'name' => '演示管理员',
            'nickname' => 'demo_admin',
            'status' => 'active',
        ]);
        $this->line('  ✅ 演示管理员已创建（账号: demo_admin, 密码: demo123456）');

        // 3. 创建教师账号
        $teacherNames = ['张老师', '李老师', '王老师', '赵老师'];
        $teachers = [];
        foreach ($teacherNames as $i => $name) {
            $teacher = User::create([
                'school_id' => $school->id,
                'role' => 'teacher',
                'username' => 'demo_teacher_' . ($i + 1),
                'password' => Hash::make('demo123456'),
                'name' => $name,
                'nickname' => $name,
                'subject' => ['语文', '数学', '英语', '科学'][$i],
                'status' => 'active',
            ]);
            $teachers[] = $teacher;
        }
        $this->line('  ✅ 4 位演示教师已创建');

        // 4. 创建班级和学生
        $classNames = ['一年级（1）班', '一年级（2）班', '二年级（1）班', '三年级（1）班'];
        $studentNames = [
            '张小明','李小红','王小红','赵小刚','刘小美','陈小龙','杨小丽','周小杰',
            '吴小芳','郑小强','孙小艺','黄小婷','林小峰','郭小雪','唐小涛','曹小敏',
            '马小亮','胡小欣','朱小磊','秦小雪','许小飞','何小雅','罗小辉','梁小静',
        ];

        foreach ($classNames as $ci => $className) {
            $class = ClassRoom::create([
                'school_id' => $school->id,
                'name' => $className,
                'grade' => '一年级',
                'teacher_id' => $teachers[$ci % count($teachers)]->id,
                'status' => 'active',
                'display_code' => 'DEMO-' . str_pad((string) ($ci + 1), 2, '0', STR_PAD_LEFT) . '-' . strtoupper(\Illuminate\Support\Str::random(4)),
            ]);

            for ($si = 0; $si < 6; $si++) {
                $name = $studentNames[($ci * 6 + $si) % count($studentNames)];
                $student = Student::create([
                    'class_id' => $class->id,
                    'name' => $name,
                    'student_no' => 'DEMO' . str_pad((string) ($ci * 6 + $si + 1), 4, '0', STR_PAD_LEFT),
                    'total_score' => rand(20, 400),
                    'status' => 'active',
                ]);

                // 创建宠物
                $types = array_keys(Pet::petTypes());
                Pet::create([
                    'student_id' => $student->id,
                    'class_id' => $class->id,
                    'type' => $types[array_rand($types)],
                    'name' => $student->name . '的宠物',
                    'level' => rand(0, 8),
                    'experience' => rand(0, 50),
                    'mood' => rand(50, 100),
                ]);

                // 创建积分记录
                for ($s = 0; $s < rand(3, 8); $s++) {
                    Score::create([
                        'student_id' => $student->id,
                        'class_id' => $class->id,
                        'amount' => rand(-5, 10),
                        'reason' => ['举手发言', '作业优秀', '遵守纪律', '帮助同学', '上课走神'][rand(0, 4)],
                        'given_by' => $teachers[$ci % count($teachers)]->id,
                        'created_at' => now()->subDays(rand(0, 14)),
                    ]);
                }
            }
            $this->line("  ✅ {$className} 已创建（6 名学生）- 班级码: {$class->display_code}");
        }

        $this->info('');
        $this->info('🎉 演示数据生成完成！');
        $this->info('━━━━━━━━━━━━━━━━━━━━');
        $this->info('  管理员账号: demo_admin');
        $this->info('  教师账号:   demo_teacher_1 ~ demo_teacher_4');
        $this->info('  密码:       demo123456');
        $this->info('  班级码:     创建成功后在控制台可见，也可以在后台班级列表中查看');
        $this->info('━━━━━━━━━━━━━━━━━━━━');
        $this->info('提示：运行 php artisan demo:clean 可清除所有演示数据');

        return 0;
    }

    private function hasDemoData(): bool
    {
        return School::where('code', 'DEMO')->exists();
    }

    private function cleanDemoData(): void
    {
        $school = School::where('code', 'DEMO')->first();
        if (!$school) {
            return;
        }

        $classIds = ClassRoom::where('school_id', $school->id)->pluck('id');
        $studentIds = Student::whereIn('class_id', $classIds)->pluck('id');

        Score::whereIn('student_id', $studentIds)->delete();
        Pet::whereIn('student_id', $studentIds)->delete();
        Student::whereIn('class_id', $classIds)->delete();
        ClassRoom::where('school_id', $school->id)->delete();
        User::where('school_id', $school->id)->delete();
        $school->delete();

        $this->line('  🗑️  演示数据已清除');
    }
}
