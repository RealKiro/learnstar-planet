<?php

declare(strict_types=1);

use App\Http\Controllers\Api\AuthController;
use App\Http\Controllers\Api\DisplayController;
use App\Http\Controllers\Api\ParentController;
use App\Http\Controllers\Api\SchoolAdminController;
use App\Http\Controllers\Api\StudentController;
use App\Http\Controllers\Api\TeacherController;
use App\Http\Controllers\Api\WechatWorkWebhookController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| 学趣星球 API v1 路由
|--------------------------------------------------------------------------
|
| 路由设计原则：
| 1. 按角色分组：auth / admin / teacher / parent / display / common
| 2. 每个分组内按功能模块划分子组
| 3. RESTful 风格：资源用名词复数，操作用 HTTP 方法
| 4. 参数放在路径中（/{id}），过滤条件放在 Query 中（?name=xxx）
|
| 路径结构：
|   /api/v1/{role}/{module}[/{action}|/{id}][/{sub-module}]
|
| 示例：
|   GET  /api/v1/teacher/scores/rules        → 积分规则列表
|   POST /api/v1/teacher/scores/give         → 加分
|   GET  /api/v1/teacher/students            → 学生列表
|   GET  /api/v1/teacher/students/{id}       → 单个学生
|
*/

Route::prefix('v1')->group(function () {

    // ====================================================================
    // 1. 认证模块 — 无需登录
    // ====================================================================
    Route::prefix('auth')->group(function () {

        // --- 1.1 账号密码登录（不同角色）---
        Route::post('teacher/login', [AuthController::class, 'teacherLoginWithCredentials'])->middleware('throttle:6,1');
        Route::post('admin/login', [AuthController::class, 'adminLoginWithCredentials'])->middleware('throttle:6,1');
        Route::post('parent/login', [AuthController::class, 'parentLoginWithCredentials'])->middleware('throttle:6,1');

        // --- 1.2 非密码登录 ---
        Route::post('class/login', [AuthController::class, 'classLogin'])->middleware('throttle:10,1');

        // --- 1.3 第三方扫码登录（仅教师） ---
        Route::prefix('teacher')->group(function () {
            Route::post('login/wechat', [AuthController::class, 'teacherLoginWithWechat']);
            Route::post('login/wechat-work', [AuthController::class, 'teacherLoginWithWechatWork']);
            Route::post('login/qq', [AuthController::class, 'teacherLoginWithQQ']);
            Route::post('login/renren', [AuthController::class, 'teacherLoginWithRenren']);
            Route::post('bind-after-scan', [AuthController::class, 'bindAfterScan']);
        });

        // --- 1.4 需登录后的认证操作 ---
        Route::middleware('auth:sanctum')->group(function () {
            Route::post('logout', [AuthController::class, 'logout']);
            Route::post('change-password', [AuthController::class, 'changePassword']);
            Route::post('refresh', [AuthController::class, 'refreshToken']);
            Route::post('bind/{platform}', [AuthController::class, 'bindThirdParty']);
            Route::delete('unbind/{platform}', [AuthController::class, 'unbindThirdParty']);
            Route::get('bindings', [AuthController::class, 'getBindings']);
        });
    });

    // ====================================================================
    // 2. 学校管理员
    // ====================================================================
    Route::prefix('admin')->middleware(['auth:sanctum', 'role:school_admin'])->group(function () {

        // --- 2.1 学校信息 ---
        Route::get('school', [SchoolAdminController::class, 'getSchool']);

        // --- 2.2 教师管理 ---
        Route::prefix('teachers')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listTeachers']);
            Route::post('/', [SchoolAdminController::class, 'createTeacher']);
            Route::post('batch-create', [SchoolAdminController::class, 'batchCreateTeachers']);
            Route::post('import', [SchoolAdminController::class, 'importTeachers']);
            Route::get('template-csv', [SchoolAdminController::class, 'downloadTeacherTemplate']);
            Route::put('{id}', [SchoolAdminController::class, 'updateTeacher']);
            Route::post('{id}/reset-password', [SchoolAdminController::class, 'resetTeacherPassword']);
            Route::delete('{id}', [SchoolAdminController::class, 'disableTeacher']);
            Route::put('{id}/classes', [SchoolAdminController::class, 'assignTeacherClasses']);
        });

        // --- 2.3 家长管理 ---
        Route::prefix('parents')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listParents']);
            Route::post('batch-create', [SchoolAdminController::class, 'batchCreateParents']);
            Route::delete('{id}', [SchoolAdminController::class, 'deleteParent']);
        });

        // --- 2.4 班级管理 ---
        Route::prefix('classes')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'index']);
            Route::post('/', [SchoolAdminController::class, 'store']);
            Route::post('batch-create', [SchoolAdminController::class, 'batchCreateClasses']);
            Route::get('{id}', [SchoolAdminController::class, 'show']);
            Route::put('{id}', [SchoolAdminController::class, 'update']);
            Route::delete('{id}', [SchoolAdminController::class, 'destroy']);
            Route::post('{id}/assign-teacher', [SchoolAdminController::class, 'assignClassTeacher']);
            Route::delete('{id}/remove-teacher', [SchoolAdminController::class, 'removeClassTeacher']);
            Route::get('{classId}/display-code', [SchoolAdminController::class, 'getDisplayCode']);
            Route::post('{classId}/display-code/refresh', [SchoolAdminController::class, 'refreshDisplayCode']);
        });

        // --- 2.5 学生管理 ---
        Route::prefix('students')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listStudents']);
            Route::post('/', [SchoolAdminController::class, 'createStudent']);
            Route::post('import', [SchoolAdminController::class, 'importStudents']);
            Route::put('{id}', [SchoolAdminController::class, 'updateStudent']);
            Route::delete('{id}', [SchoolAdminController::class, 'deleteStudent']);
            Route::post('batch-delete', [SchoolAdminController::class, 'batchDeleteStudents']);
            Route::post('batch-move', [SchoolAdminController::class, 'batchMoveStudents']);
        });

        // --- 2.6 学年升级 ---
        Route::get('grade-upgrade/preview', [SchoolAdminController::class, 'previewGradeUpgrade']);
        Route::post('grade-upgrade/execute', [SchoolAdminController::class, 'executeGradeUpgrade']);

        // --- 2.7 报表 ---
        Route::prefix('reports')->group(function () {
            Route::get('overview', [SchoolAdminController::class, 'schoolOverview']);
            Route::get('by-grade', [SchoolAdminController::class, 'reportsByGrade']);
            Route::get('by-class', [SchoolAdminController::class, 'reportsByClass']);
        });

        // --- 2.8 汇率管理 ---
        Route::prefix('exchange-rates')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listExchangeRates']);
            Route::post('/', [SchoolAdminController::class, 'createExchangeRate']);
            Route::put('{id}', [SchoolAdminController::class, 'updateExchangeRate']);
        });
    });

    // ====================================================================
    // 3. 教师端
    // ====================================================================
    Route::prefix('teacher')->middleware(['auth:sanctum', 'role:teacher'])->group(function () {

        // --- 3.1 班级与模式 ---
        Route::get('mode', [TeacherController::class, 'getMode']);
        Route::post('mode', [TeacherController::class, 'setMode']);
        Route::get('my-classes', [TeacherController::class, 'myClasses']);
        Route::post('switch-class', [TeacherController::class, 'switchClass']);
        Route::get('dashboard', [TeacherController::class, 'dashboard']);

        // --- 3.2 学生管理 ---
        Route::prefix('students')->group(function () {
            Route::get('/', [TeacherController::class, 'listStudents']);
            Route::post('import', [TeacherController::class, 'importStudents']);
            Route::put('{id}', [TeacherController::class, 'updateStudent']);
        });

        // --- 3.3 教室大屏 ---
        Route::prefix('classroom')->group(function () {
            Route::get('display', [TeacherController::class, 'classroomDisplay']);
            Route::get('messages', [TeacherController::class, 'pollClassroomMessages']);
            Route::post('messages', [TeacherController::class, 'sendClassroomMessage']);
        });

        // --- 3.4 积分管理 ---
        Route::prefix('scores')->group(function () {
            Route::get('summary', [TeacherController::class, 'scoreSummary']);
            Route::post('give', [TeacherController::class, 'giveScore']);
            Route::post('batch-give', [TeacherController::class, 'batchGiveScore']);
            Route::post('give-by-rule/{ruleId}', [TeacherController::class, 'giveScoreByRule']);
            Route::get('history/{studentId}', [TeacherController::class, 'scoreHistory']);
            Route::get('rules', [TeacherController::class, 'listScoreRules']);
            Route::post('rules', [TeacherController::class, 'createScoreRule']);
            Route::put('rules/{id}', [TeacherController::class, 'updateScoreRule']);
            Route::delete('rules/{id}', [TeacherController::class, 'deleteScoreRule']);
        });

        // --- 3.5 宠物系统 ---
        Route::prefix('pets')->group(function () {
            Route::get('types', [TeacherController::class, 'getPetTypes']);
            Route::get('overview', [TeacherController::class, 'classPetsOverview']);
            Route::get('{studentId}', [TeacherController::class, 'getPet']);
            Route::post('{studentId}/switch', [TeacherController::class, 'switchPet']);
            Route::post('{studentId}/feed', [TeacherController::class, 'feedPet']);
            Route::post('{studentId}/rename', [TeacherController::class, 'renamePet']);
        });

        // --- 3.6 排行榜 & 年级战场 ---
        Route::prefix('leaderboard')->group(function () {
            Route::get('total', [TeacherController::class, 'totalLeaderboard']);
            Route::get('weekly', [TeacherController::class, 'weeklyLeaderboard']);
            Route::get('pet-level', [TeacherController::class, 'petLevelLeaderboard']);
        });

        Route::prefix('pk')->group(function () {
            Route::get('leaderboard', [TeacherController::class, 'pkLeaderboard']);
            Route::get('my-stats', [TeacherController::class, 'myPkStats']);
            Route::post('challenge', [TeacherController::class, 'challengePk']);
        });

        // --- 3.7 班级配置 ---
        Route::prefix('class')->group(function () {
            Route::post('switch-series', [TeacherController::class, 'switchSeries']);
        });

        // --- 3.8 积分商城 ---
        Route::prefix('shop')->group(function () {
            Route::get('items', [TeacherController::class, 'listShopItems']);
            Route::post('items', [TeacherController::class, 'createShopItem']);
            Route::put('items/{id}', [TeacherController::class, 'updateShopItem']);
            Route::delete('items/{id}', [TeacherController::class, 'deleteShopItem']);
            Route::get('redemptions', [TeacherController::class, 'listRedemptions']);
            Route::post('redemptions', [TeacherController::class, 'createRedemption']);
            Route::put('redemptions/{id}/approve', [TeacherController::class, 'approveRedemption']);
            Route::put('redemptions/{id}/reject', [TeacherController::class, 'rejectRedemption']);
            Route::put('redemptions/{id}/deliver', [TeacherController::class, 'deliverRedemption']);
        });

        // --- 3.9 班级通知 ---
        Route::prefix('notices')->group(function () {
            Route::get('/', [TeacherController::class, 'listNotices']);
            Route::post('/', [TeacherController::class, 'createNotice']);
            Route::put('{id}', [TeacherController::class, 'updateNotice']);
            Route::put('{id}/publish', [TeacherController::class, 'publishNotice']);
            Route::delete('{id}', [TeacherController::class, 'deleteNotice']);
        });

        // --- 3.10 数据报表 ---
        Route::prefix('reports')->group(function () {
            Route::get('score-trend', [TeacherController::class, 'scoreTrend']);
            Route::get('pet-distribution', [TeacherController::class, 'petDistribution']);
            Route::get('student-progress', [TeacherController::class, 'studentProgress']);
            Route::get('export/{type}', [TeacherController::class, 'exportReport']);
        });

        // --- 3.11 实时广播 ---
        Route::prefix('broadcasts')->group(function () {
            Route::get('/', [TeacherController::class, 'listBroadcasts']);
            Route::post('/', [TeacherController::class, 'sendBroadcast']);
            Route::get('{id}', [TeacherController::class, 'getBroadcast']);
        });

        // --- 3.12 智能考勤 ---
        Route::prefix('attendance')->group(function () {
            Route::get('today', [TeacherController::class, 'getTodayAttendance']);
            Route::post('start', [TeacherController::class, 'startAttendance']);
            Route::put('{studentId}', [TeacherController::class, 'setAttendance']);
            Route::post('{studentId}/mark-leave', [TeacherController::class, 'markManualLeave']);
            Route::post('{studentId}/mark-absent', [TeacherController::class, 'markManualAbsent']);
            Route::get('summary', [TeacherController::class, 'attendanceSummary']);
        });

        // --- 3.13 作业管理 ---
        Route::prefix('homework')->group(function () {
            Route::get('/', [TeacherController::class, 'listHomework']);
            Route::post('/', [TeacherController::class, 'createHomework']);
            Route::get('{id}', [TeacherController::class, 'getHomework']);
            Route::put('{id}/close', [TeacherController::class, 'closeHomework']);
            Route::get('{id}/submissions', [TeacherController::class, 'getHomeworkSubmissions']);
            Route::get('{id}/qr-code', [TeacherController::class, 'getHomeworkQrCode']);
        });

        // --- 3.14 在线答题 ---
        Route::prefix('quizzes')->group(function () {
            Route::get('/', [TeacherController::class, 'listQuizzes']);
            Route::post('/', [TeacherController::class, 'createQuiz']);
            Route::post('{id}/start', [TeacherController::class, 'startQuiz']);
            Route::post('{id}/stop', [TeacherController::class, 'stopQuiz']);
            Route::get('{id}/stats', [TeacherController::class, 'getQuizStats']);
        });

        // --- 3.15 题库 ---
        Route::prefix('question-banks')->group(function () {
            Route::get('/', [TeacherController::class, 'listQuestionBanks']);
            Route::post('/', [TeacherController::class, 'createQuestionBank']);
            Route::post('{id}/questions', [TeacherController::class, 'addQuestion']);
            Route::get('{id}/questions', [TeacherController::class, 'getQuestions']);
        });

        // --- 3.16 成绩管理 ---
        Route::prefix('grades')->group(function () {
            Route::get('/', [TeacherController::class, 'listGrades']);
            Route::post('/', [TeacherController::class, 'inputGrades']);
            Route::get('stats', [TeacherController::class, 'getGradeStats']);
            Route::get('distribution', [TeacherController::class, 'getGradeDistribution']);
        });

        // --- 3.17 AI 助教 ---
        Route::prefix('ai')->group(function () {
            Route::post('chat', [TeacherController::class, 'aiChat']);
            Route::get('commands', [TeacherController::class, 'getAiCommands']);
            Route::get('usage', [TeacherController::class, 'getAiUsage']);
        });

        // --- 3.18 多币种 ---
        Route::prefix('currency')->group(function () {
            Route::get('wallets', [TeacherController::class, 'listWallets']);
            Route::post('exchange', [TeacherController::class, 'exchangeCurrency']);
            Route::post('cross-exchange', [TeacherController::class, 'crossExchangeCurrency']);
        });

        // --- 3.19 教师端大屏管理 ---
        Route::get('display-code', [DisplayController::class, 'getDisplayCode']);
        Route::post('display-code/refresh', [DisplayController::class, 'refreshDisplayCode']);
    });

    // ====================================================================
    // 4. 家长端
    // ====================================================================
    Route::prefix('parent')->middleware(['auth:sanctum', 'role:parent'])->group(function () {
        Route::get('home', [ParentController::class, 'home']);

        Route::prefix('scores')->group(function () {
            Route::get('detail', [ParentController::class, 'scoreDetail']);
            Route::get('history', [ParentController::class, 'scoreHistory']);
        });

        Route::prefix('growth')->group(function () {
            Route::get('log', [ParentController::class, 'growthLog']);
            Route::get('timeline', [ParentController::class, 'growthTimeline']);
        });

        Route::prefix('pet')->group(function () {
            Route::get('/', [ParentController::class, 'petStatus']);
            Route::post('feed', [ParentController::class, 'feedPet']);
        });

        Route::get('ranking', [ParentController::class, 'ranking']);

        Route::prefix('notices')->group(function () {
            Route::get('/', [ParentController::class, 'listNotices']);
            Route::get('{id}', [ParentController::class, 'readNotice']);
        });
    });

    // ====================================================================
    // 5. 班级大屏（无需登录，使用班级码 Token）
    // ====================================================================
    Route::prefix('display')->group(function () {
        Route::post('login', [DisplayController::class, 'login'])->middleware('throttle:10,1');
        Route::get('initial-data', [DisplayController::class, 'initialData']);
        Route::get('sse', [DisplayController::class, 'sse']);
        Route::get('poll', [DisplayController::class, 'poll']);
        Route::post('quick-score', [DisplayController::class, 'quickScore']);
        Route::get('leaderboard', [DisplayController::class, 'quickLeaderboard']);
        Route::get('shop-items', [DisplayController::class, 'quickShopItems']);
        Route::post('redeem', [DisplayController::class, 'quickRedeem']);
        Route::post('transfer', [DisplayController::class, 'quickTransfer']);
    });

    // ====================================================================
    // 6. 公共接口（无需登录）
    // ====================================================================
    Route::prefix('common')->group(function () {
        Route::get('pet-types', [StudentController::class, 'petTypes']);
        Route::get('evolution-stages', [StudentController::class, 'evolutionStages']);
        Route::get('score-categories', [StudentController::class, 'scoreCategories']);
    });

    // ====================================================================
    // 7. 企业微信 Webhook
    // ====================================================================
    Route::prefix('wechat-work')->group(function () {
        Route::get('callback', [WechatWorkWebhookController::class, 'verify']);
        Route::post('callback', [WechatWorkWebhookController::class, 'receive']);
    });

}); // End API v1

// ====================================================================
// 8. 向后兼容路由（逐步弃用，仅供旧版前端使用）
// ====================================================================

// --- 8.1 旧版登录路径 ---
Route::prefix('auth')->group(function () {
    Route::post('teacher/login', [AuthController::class, 'teacherLoginWithCredentials'])->middleware('throttle:6,1');
    Route::post('admin/login', [AuthController::class, 'adminLoginWithCredentials'])->middleware('throttle:6,1');
    Route::post('parent/login', [AuthController::class, 'parentLoginWithCredentials'])->middleware('throttle:6,1');
});

// --- 8.2 旧版教师端点别名（v1 命名调整前的兼容） ---
Route::prefix('v1/teacher')->middleware(['auth:sanctum', 'role:teacher'])->group(function () {
    // classroom-display → classroom/display
    Route::get('classroom-display', [TeacherController::class, 'classroomDisplay']);
    Route::get('classroom-messages', [TeacherController::class, 'pollClassroomMessages']);
    Route::post('classroom-messages', [TeacherController::class, 'sendClassroomMessage']);
    // pets/class-overview → pets/overview
    Route::get('pets/class-overview', [TeacherController::class, 'classPetsOverview']);
    // leaderboard/pk → pk/leaderboard
    Route::get('leaderboard/pk', [TeacherController::class, 'pkLeaderboard']);
    Route::get('leaderboard/my-stats', [TeacherController::class, 'myPkStats']);
});
