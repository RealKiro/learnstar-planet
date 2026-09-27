<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Services\ScoreRuleService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class StudentController extends Controller
{
    /**
     * 获取积分分类。
     *
     * id / 名称一律取自 {@see ScoreRuleService::CATEGORY_LABELS}（后端唯一真源，键集与前端
     * utils/scoreRules.ts 的 categoryLabels 一致）——本端点曾自带第三套词汇
     * （discipline / teamwork / cleanliness / other，且 homework 叫「作业完成」），
     * 与界面标签互相矛盾，已收口。
     */
    public function scoreCategories(): JsonResponse
    {
        /** @var array<string, string> $icons 分类图标（缺项由 loop 内兜底，新增分类不会渲染为空） */
        $icons = [
            'classroom' => '📖',
            'homework' => '📝',
            'behavior' => '🌟',
            'literacy' => '📊',
            'daily' => '📅',
            'academic' => '📚',
            'custom' => '✨',
        ];

        $data = [];
        $sort = 0;
        foreach (ScoreRuleService::CATEGORY_LABELS as $id => $name) {
            $data[] = [
                'id' => $id,
                'name' => $name,
                'icon' => $icons[$id] ?? '📌',
                'sort' => ++$sort,
            ];
        }

        return response()->json(['data' => $data]);
    }
}
