<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ParentController extends Controller
{
    /**
     * 家长首页
     */
    public function home(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    // ===== 积分 =====

    public function scoreDetail(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function scoreHistory(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    // ===== 成长 =====

    public function growthLog(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function growthTimeline(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    // ===== 宠物 =====

    public function petStatus(Request $request): JsonResponse
    {
        return response()->json(['data' => null]);
    }

    public function feedPet(Request $request): JsonResponse
    {
        return response()->json(['message' => '喂养成功']);
    }

    // ===== 排行榜 =====

    public function ranking(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    // ===== 通知 =====

    public function listNotices(Request $request): JsonResponse
    {
        return response()->json(['data' => []]);
    }

    public function readNotice(Request $request, int $id): JsonResponse
    {
        return response()->json(['message' => '已读']);
    }
}
