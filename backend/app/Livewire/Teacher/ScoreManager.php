<?php

declare(strict_types=1);

namespace App\Livewire\Teacher;

use App\Models\ClassRoom;
use App\Models\Score;
use App\Models\Student;
use Livewire\Component;
use Livewire\WithFileUploads;

/**
 * 积分管理 Livewire 组件
 * 全免费：所有积分规则、批量操作均可使用
 */
class ScoreManager extends Component
{
    use WithFileUploads;

    public $classRoomId;

    public $selectedStudentId;

    public $scoreType = 'add';

    public $amount = 5;

    public $ruleCode = 'homework_complete';

    public $comment = '';

    public $selectedStudents = [];

    public $batchAmount = 5;

    public $batchRuleCode = 'homework_complete';

    public $batchComment = '';

    public $showBatchModal = false;

    public $showQuickPanel = true;

    protected $rules = [
        'amount' => 'required|integer|min:1|max:100',
        'ruleCode' => 'required|string',
        'comment' => 'nullable|string|max:200',
    ];

    public function mount($classRoomId)
    {
        $this->classRoomId = $classRoomId;
    }

    /**
     * 快捷加分
     */
    public function quickAdd(int $studentId, string $ruleCode)
    {
        $rules = Score::getRules();
        $rule = $rules[$ruleCode];

        try {
            app(\App\Services\ScoreService::class)->quickScore(
                $studentId,
                auth()->user()->teacher->id,
                $ruleCode
            );

            $student = Student::find($studentId);
            flux()->toast("{$student->name} {$rule['name']} +{$rule['amount']}分", 'success');
        } catch (\Exception $e) {
            flux()->toast('操作失败：' . $e->getMessage(), 'error');
        }
    }

    /**
     * 自定义加减分
     */
    public function customScore()
    {
        $this->validate();

        try {
            app(\App\Services\ScoreService::class)->adjustScore(
                $this->selectedStudentId,
                auth()->user()->teacher->id,
                $this->scoreType,
                $this->amount,
                $this->ruleCode,
                $this->comment
            );

            $student = Student::find($this->selectedStudentId);
            $typeLabel = $this->scoreType === 'add' ? '加' : '减';
            flux()->toast("{$student->name} {$typeLabel}{$this->amount}分", 'success');

            $this->reset(['amount', 'comment']);
        } catch (\Exception $e) {
            flux()->toast('操作失败：' . $e->getMessage(), 'error');
        }
    }

    /**
     * 批量加分
     */
    public function batchAdd()
    {
        if (empty($this->selectedStudents)) {
            flux()->toast('请先选择学生', 'warning');

            return;
        }

        try {
            app(\App\Services\ScoreService::class)->batchAddScore(
                $this->selectedStudents,
                auth()->user()->teacher->id,
                $this->batchAmount,
                $this->batchRuleCode,
                $this->batchComment
            );

            flux()->toast('已为 ' . count($this->selectedStudents) . ' 名学生加分', 'success');
            $this->selectedStudents = [];
            $this->showBatchModal = false;
        } catch (\Exception $e) {
            flux()->toast('批量操作失败：' . $e->getMessage(), 'error');
        }
    }

    public function render()
    {
        $students = Student::where('class_room_id', $this->classRoomId)
            ->where('is_active', true)
            ->orderBy('total_score', 'desc')
            ->get();

        $rules = Score::getRules();

        return view('livewire.teacher.score-manager', compact('students', 'rules'));
    }
}
