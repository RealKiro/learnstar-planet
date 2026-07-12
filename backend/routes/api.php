<?php

declare(strict_types=1);

use App\Http\Controllers\Api\AuthController;
use App\Http\Controllers\Api\ParentController;
use App\Http\Controllers\Api\SchoolAdminController;
use App\Http\Controllers\Api\StudentController;
use App\Http\Controllers\Api\TeacherController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| 学趣星球 API 路由
|--------------------------------------------------------------------------
|
| 路由结构按角色清晰分组：
| - /api/v1/auth/*          → 登录认证（账号密码 + 第三方扫码）
| - /api/v1/admin/*         → 学校管理员（分配账号、管理班级）
| - /api/v1/teacher/*       → 教师（积分管理、宠物、通知、教室小喇叭）
| - /api/v1/parent/*        → 家长（查看积分、宠物状态）
| - /api/v1/common/*        → 公共接口（排行榜等）
|
*/

// API v1
Route::prefix('v1')->group(function () {
    // ===== 认证 =====
    Route::prefix('auth')->group(function () {
        // 教师账号密码登录（速率限制：每分钟 6 次）
        Route::post('teacher/login', [AuthController::class, 'teacherLoginWithCredentials'])->middleware('throttle:6,1');

        // 管理员账号密码登录
        Route::post('admin/login', [AuthController::class, 'adminLoginWithCredentials'])->middleware('throttle:6,1');

        // 家长账号密码登录
        Route::post('parent/login', [AuthController::class, 'parentLoginWithCredentials'])->middleware('throttle:6,1');

        // 第三方扫码登录（仅教师可用）
        Route::post('teacher/login/wechat', [AuthController::class, 'teacherLoginWithWechat']);
        Route::post('teacher/login/wechat-work', [AuthController::class, 'teacherLoginWithWechatWork']);
        Route::post('teacher/login/qq', [AuthController::class, 'teacherLoginWithQQ']);
        Route::post('teacher/login/renren', [AuthController::class, 'teacherLoginWithRenren']);

        // 第三方扫码后绑定已有教师账号
        Route::post('teacher/bind-after-scan', [AuthController::class, 'bindAfterScan']);

        // 需登录后操作
        Route::middleware('auth:sanctum')->group(function () {
            Route::post('logout', [AuthController::class, 'logout']);
            Route::post('change-password', [AuthController::class, 'changePassword']);
            Route::post('refresh', [AuthController::class, 'refreshToken']);

            // 第三方账号绑定/解绑
            Route::post('bind/{platform}', [AuthController::class, 'bindThirdParty']);
            Route::delete('unbind/{platform}', [AuthController::class, 'unbindThirdParty']);
            Route::get('bindings', [AuthController::class, 'getBindings']);
        });
    });

    // ===== 学校管理员 =====
    Route::prefix('admin')->middleware(['auth:sanctum', 'role:school_admin'])->group(function () {
        Route::get('school', [SchoolAdminController::class, 'getSchool']);
        Route::put('school', [SchoolAdminController::class, 'updateSchool']);
        Route::post('teachers/batch-create', [SchoolAdminController::class, 'batchCreateTeachers']);
        Route::get('teachers', [SchoolAdminController::class, 'listTeachers']);
        Route::put('teachers/{id}', [SchoolAdminController::class, 'updateTeacher']);
        Route::post('teachers/{id}/reset-password', [SchoolAdminController::class, 'resetTeacherPassword']);
        Route::delete('teachers/{id}', [SchoolAdminController::class, 'disableTeacher']);
        Route::post('parents/batch-create', [SchoolAdminController::class, 'batchCreateParents']);
        Route::get('parents', [SchoolAdminController::class, 'listParents']);
        Route::delete('parents/{id}', [SchoolAdminController::class, 'deleteParent']);
        Route::get('classes', [SchoolAdminController::class, 'index']);
        Route::post('classes', [SchoolAdminController::class, 'store']);
        Route::post('classes/batch-create', [SchoolAdminController::class, 'batchCreateClasses']);
        Route::get('classes/{id}', [SchoolAdminController::class, 'show']);
        Route::put('classes/{id}', [SchoolAdminController::class, 'update']);
        Route::delete('classes/{id}', [SchoolAdminController::class, 'destroy']);
        Route::post('classes/{id}/assign-teacher', [SchoolAdminController::class, 'assignClassTeacher']);
        // 学生管理
        Route::post('students/import', [SchoolAdminController::class, 'importStudents']);
        Route::get('students', [SchoolAdminController::class, 'listStudents']);
        Route::post('students', [SchoolAdminController::class, 'createStudent']);
        Route::put('students/{id}', [SchoolAdminController::class, 'updateStudent']);
        Route::delete('students/{id}', [SchoolAdminController::class, 'deleteStudent']);
        Route::post('students/batch-delete', [SchoolAdminController::class, 'batchDeleteStudents']);
        Route::post('students/batch-move', [SchoolAdminController::class, 'batchMoveStudents']);
        // 学年升级
        Route::get('grade-upgrade/preview', [SchoolAdminController::class, 'previewGradeUpgrade']);
        Route::post('grade-upgrade/execute', [SchoolAdminController::class, 'executeGradeUpgrade']);
        // 报表
        Route::get('reports/overview', [SchoolAdminController::class, 'schoolOverview']);
        Route::get('reports/by-grade', [SchoolAdminController::class, 'reportsByGrade']);
        Route::get('reports/by-class', [SchoolAdminController::class, 'reportsByClass']);
        // 汇率管理
        Route::get('exchange-rates', [SchoolAdminController::class, 'listExchangeRates']);
        Route::post('exchange-rates', [SchoolAdminController::class, 'createExchangeRate']);
        Route::put('exchange-rates/{id}', [SchoolAdminController::class, 'updateExchangeRate']);
    });

    // ===== 教师 =====
    Route::prefix('teacher')->middleware(['auth:sanctum', 'role:teacher'])->group(function () {
        Route::get('dashboard', [TeacherController::class, 'dashboard']);
        Route::get('students', [TeacherController::class, 'listStudents']);
        Route::post('students/import', [TeacherController::class, 'importStudents']);
        Route::put('students/{id}', [TeacherController::class, 'updateStudent']);

        // ===== 积分管理 =====
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

        // ===== 宠物园 =====
        Route::prefix('pets')->group(function () {
            Route::get('class-overview', [TeacherController::class, 'classPetsOverview']);
            Route::get('{studentId}', [TeacherController::class, 'getPet']);
            Route::post('{studentId}/feed', [TeacherController::class, 'feedPet']);
            Route::post('{studentId}/rename', [TeacherController::class, 'renamePet']);
        });

        // ===== 排行榜 =====
        Route::prefix('leaderboard')->group(function () {
            Route::get('total', [TeacherController::class, 'totalLeaderboard']);
            Route::get('weekly', [TeacherController::class, 'weeklyLeaderboard']);
            Route::get('pet-level', [TeacherController::class, 'petLevelLeaderboard']);
        });

        // ===== 积分商城 =====
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

        // ===== 班级通知 =====
        Route::prefix('notices')->group(function () {
            Route::get('/', [TeacherController::class, 'listNotices']);
            Route::post('/', [TeacherController::class, 'createNotice']);
            Route::put('{id}', [TeacherController::class, 'updateNotice']);
            Route::put('{id}/publish', [TeacherController::class, 'publishNotice']);
            Route::delete('{id}', [TeacherController::class, 'deleteNotice']);
        });

        // ===== 数据报表 =====
        Route::prefix('reports')->group(function () {
            Route::get('score-trend', [TeacherController::class, 'scoreTrend']);
            Route::get('pet-distribution', [TeacherController::class, 'petDistribution']);
            Route::get('student-progress', [TeacherController::class, 'studentProgress']);
            Route::get('export/{type}', [TeacherController::class, 'exportReport']);
        });

        // ===== 教室小喇叭新功能 =====
        // -- 实时广播
        Route::prefix('broadcasts')->group(function () {
            Route::get('/', [TeacherController::class, 'listBroadcasts']);
            Route::post('/', [TeacherController::class, 'sendBroadcast']);
            Route::get('{id}', [TeacherController::class, 'getBroadcast']);
        });

        // -- 智能考勤
        Route::prefix('attendance')->group(function () {
            Route::get('today', [TeacherController::class, 'getTodayAttendance']);
            Route::post('start', [TeacherController::class, 'startAttendance']);
            Route::put('{studentId}', [TeacherController::class, 'setAttendance']);
            Route::get('summary', [TeacherController::class, 'attendanceSummary']);
        });

        // -- 扫码收作业
        Route::prefix('homework')->group(function () {
            Route::get('/', [TeacherController::class, 'listHomework']);
            Route::post('/', [TeacherController::class, 'createHomework']);
            Route::get('{id}', [TeacherController::class, 'getHomework']);
            Route::put('{id}/close', [TeacherController::class, 'closeHomework']);
            Route::get('{id}/submissions', [TeacherController::class, 'getHomeworkSubmissions']);
            Route::get('{id}/qr-code', [TeacherController::class, 'getHomeworkQrCode']);
        });

        // -- 在线答题 & 题库
        Route::prefix('quizzes')->group(function () {
            Route::get('/', [TeacherController::class, 'listQuizzes']);
            Route::post('/', [TeacherController::class, 'createQuiz']);
            Route::post('{id}/start', [TeacherController::class, 'startQuiz']);
            Route::post('{id}/stop', [TeacherController::class, 'stopQuiz']);
            Route::get('{id}/stats', [TeacherController::class, 'getQuizStats']);
        });

        Route::prefix('question-banks')->group(function () {
            Route::get('/', [TeacherController::class, 'listQuestionBanks']);
            Route::post('/', [TeacherController::class, 'createQuestionBank']);
            Route::post('{id}/questions', [TeacherController::class, 'addQuestion']);
            Route::get('{id}/questions', [TeacherController::class, 'getQuestions']);
        });

        // -- 成绩管理
        Route::prefix('grades')->group(function () {
            Route::get('/', [TeacherController::class, 'listGrades']);
            Route::post('/', [TeacherController::class, 'inputGrades']);
            Route::get('stats', [TeacherController::class, 'getGradeStats']);
            Route::get('distribution', [TeacherController::class, 'getGradeDistribution']);
        });

        // -- AI助教
        Route::prefix('ai')->group(function () {
            Route::post('chat', [TeacherController::class, 'aiChat']);
            Route::get('commands', [TeacherController::class, 'getAiCommands']);
            Route::get('usage', [TeacherController::class, 'getAiUsage']);
        });

        // ===== 多币种系统 =====
        Route::prefix('currency')->group(function () {
            Route::get('wallets', [TeacherController::class, 'listWallets']);
            Route::post('exchange', [TeacherController::class, 'exchangeCurrency']);
            Route::post('cross-exchange', [TeacherController::class, 'crossExchangeCurrency']);
        });
    });

    // ===== 家长 =====
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

    // ===== 公共接口 =====
    Route::prefix('common')->group(function () {
        Route::get('pet-types', [StudentController::class, 'petTypes']);
        Route::get('evolution-stages', [StudentController::class, 'evolutionStages']);
        Route::get('score-categories', [StudentController::class, 'scoreCategories']);
    });
}); // End API v1 prefix

// 向后兼容：v1 之前的路由也映射到 v1
Route::prefix('auth')->group(function () {
    Route::post('teacher/login', [AuthController::class, 'teacherLoginWithCredentials'])->middleware('throttle:6,1');
    Route::post('admin/login', [AuthController::class, 'adminLoginWithCredentials'])->middleware('throttle:6,1');
    Route::post('parent/login', [AuthController::class, 'parentLoginWithCredentials'])->middleware('throttle:6,1');
    Route::post('teacher/login/{platform}', [AuthController::class, 'teacherLoginWithWechat']);
});
Route::prefix('admin')->middleware(['auth:sanctum', 'role:school_admin'])->group(function () {
    Route::get('school', [SchoolAdminController::class, 'getSchool']);
    Route::put('school', [SchoolAdminController::class, 'updateSchool']);
    Route::get('teachers', [SchoolAdminController::class, 'listTeachers']);
    Route::get('classes', [SchoolAdminController::class, 'index']);
    Route::get('students', [SchoolAdminController::class, 'listStudents']);
    Route::get('reports/overview', [SchoolAdminController::class, 'schoolOverview']);
});
