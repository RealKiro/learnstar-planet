<?php

declare(strict_types=1);

namespace App\Livewire\Teacher;

use App\Models\ClassRoom;
use App\Services\LeaderboardService;
use App\Services\ScoreService;
use Livewire\Component;

/**
 * 教师端数据看板 Livewire 组件
 * 全免费：所有数据、报表、功能均可使用
 */
class Dashboard extends Component
{
    public $classRoomId;

    public $classRooms = [];

    public $stats = [];

    public $recentScores = [];

    public $topStudents = [];

    public $scoreTrend = [];

    public $petDistribution = [];

    public function mount()
    {
        $teacher = auth()->user()->teacher;
        $this->classRooms = $teacher->classRooms()->where('is_active', true)->get();

        if ($this->classRooms->isNotEmpty()) {
            $this->classRoomId = $this->classRooms->first()->id;
            $this->loadDashboard();
        }
    }

    public function loadDashboard()
    {
        $scoreService = app(ScoreService::class);
        $leaderboardService = app(LeaderboardService::class);

        $this->stats = $scoreService->getClassStats($this->classRoomId);
        $this->topStudents = $leaderboardService->getClassLeaderboard($this->classRoomId, 10);
        $this->scoreTrend = $scoreService->getScoreTrend($this->classRoomId, 7);

        $classRoom = ClassRoom::find($this->classRoomId);
        $this->petDistribution = $classRoom?->getPetLevelStats() ?? [];
    }

    public function updatedClassRoomId()
    {
        $this->loadDashboard();
    }

    public function render()
    {
        return view('livewire.teacher.dashboard');
    }
}

