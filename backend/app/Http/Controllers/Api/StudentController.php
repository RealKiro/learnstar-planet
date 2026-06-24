<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;

class StudentController extends Controller
{
    /**
     * 获取宠物类型列表
     */
    public function petTypes(): JsonResponse
    {
        return response()->json([
            'data' => \App\Models\Pet::petTypes(),
        ]);
    }

    /**
     * 获取进化阶段列表
     */
    public function evolutionStages(): JsonResponse
    {
        return response()->json([
            'data' => \App\Models\Pet::evolutionStages(),
        ]);
    }

    /**
     * 获取积分分类
     */
    public function scoreCategories(): JsonResponse
    {
        return response()->json([
            'data' => [],
        ]);
    }
}
