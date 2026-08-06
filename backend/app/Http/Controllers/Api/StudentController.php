<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class StudentController extends Controller
{
    /**
     * 获取积分分类
     */
    public function scoreCategories(): JsonResponse
    {
        return response()->json([
            'data' => [
                ['id' => 'classroom', 'name' => '课堂表现', 'icon' => '📖', 'sort' => 1],
                ['id' => 'homework', 'name' => '作业完成', 'icon' => '📝', 'sort' => 2],
                ['id' => 'discipline', 'name' => '遵守纪律', 'icon' => '📏', 'sort' => 3],
                ['id' => 'teamwork', 'name' => '团结互助', 'icon' => '🤝', 'sort' => 4],
                ['id' => 'activity', 'name' => '活动参与', 'icon' => '🎯', 'sort' => 5],
                ['id' => 'cleanliness', 'name' => '卫生值日', 'icon' => '🧹', 'sort' => 6],
                ['id' => 'other', 'name' => '其他加分', 'icon' => '✨', 'sort' => 7],
            ],
        ]);
    }
}
