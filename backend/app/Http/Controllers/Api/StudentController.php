<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ClassRoom;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class StudentController extends Controller
{
    /**
     * 获取宠物类型列表
     * 不传 class_id：返回全部种类
     * 传 class_id：只返回该班级所绑定的系列种类
     */
    public function petTypes(Request $request): JsonResponse
    {
        $classId = (int) $request->query('class_id', 0);

        if ($classId > 0) {
            $class = ClassRoom::find($classId);
            if ($class) {
                $series = $class->settings['pet_series'] ?? null;
                if ($series && $series !== 'all') {
                    return response()->json([
                        'data' => \App\Models\Pet::petTypesBySeries($series),
                        'series' => $series,
                    ]);
                }
            }
        }

        // 默认返回全部分类和种类
        return response()->json([
            'data' => \App\Models\Pet::petTypes(),
            'series' => 'all',
            'categories' => \App\Models\Pet::petCategories(),
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
