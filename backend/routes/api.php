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
| API v1 — 学趣星球
| 规范：/api/v1/{role}/{resource}[/{id}][/{action}]
*/
Route::prefix('v1')->group(function () {
    // ===== 1. 认证 (无需登录) =====
    Route::prefix('auth')->group(function () {
        Route::post('teacher/login', [AuthController::class, 'teacherLoginWithCredentials'])->middleware('throttle:6,1');
        Route::post('admin/login', [AuthController::class, 'adminLoginWithCredentials'])->middleware('throttle:6,1');
        Route::post('parent/login', [AuthController::class, 'parentLoginWithCredentials'])->middleware('throttle:6,1');
        Route::post('class/login', [AuthController::class, 'classLogin'])->middleware('throttle:10,1');
        Route::prefix('teacher')->group(function () {
            Route::post('login/wechat', [AuthController::class, 'teacherLoginWithWechat']);
            Route::post('login/wechat-work', [AuthController::class, 'teacherLoginWithWechatWork']);
            Route::post('login/qq', [AuthController::class, 'teacherLoginWithQQ']);
            Route::post('login/renren', [AuthController::class, 'teacherLoginWithRenren']);
            Route::post('bind-after-scan', [AuthController::class, 'bindAfterScan']);
        });
        Route::middleware('auth:sanctum')->group(function () {
            Route::post('logout', [AuthController::class, 'logout']);
            Route::post('change-password', [AuthController::class, 'changePassword']);
            Route::post('refresh', [AuthController::class, 'refreshToken']);
            Route::post('bind/{platform}', [AuthController::class, 'bindThirdParty']);
            Route::delete('unbind/{platform}', [AuthController::class, 'unbindThirdParty']);
            Route::get('bindings', [AuthController::class, 'getBindings']);
        });
    });

    // ===== 2. 学校管理员 =====
    Route::prefix('admin')->middleware(['auth:sanctum', 'role:school_admin'])->group(function () {
        Route::get('school', [SchoolAdminController::class, 'getSchool']);
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
        Route::prefix('parents')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listParents']);
            Route::post('batch-create', [SchoolAdminController::class, 'batchCreateParents']);
            Route::delete('{id}', [SchoolAdminController::class, 'deleteParent']);
        });
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
        Route::prefix('students')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listStudents']);
            Route::post('/', [SchoolAdminController::class, 'createStudent']);
            Route::post('import', [SchoolAdminController::class, 'importStudents']);
            Route::put('{id}', [SchoolAdminController::class, 'updateStudent']);
            Route::delete('{id}', [SchoolAdminController::class, 'deleteStudent']);
            Route::post('batch-delete', [SchoolAdminController::class, 'batchDeleteStudents']);
            Route::post('batch-move', [SchoolAdminController::class, 'batchMoveStudents']);
        });
        Route::get('grade-upgrade/preview', [SchoolAdminController::class, 'previewGradeUpgrade']);
        Route::post('grade-upgrade/execute', [SchoolAdminController::class, 'executeGradeUpgrade']);
        Route::prefix('reports')->group(function () {
            Route::get('overview', [SchoolAdminController::class, 'schoolOverview']);
            Route::get('by-grade', [SchoolAdminController::class, 'reportsByGrade']);
            Route::get('by-class', [SchoolAdminController::class, 'reportsByClass']);
        });
        Route::prefix('exchange-rates')->group(function () {
            Route::get('/', [SchoolAdminController::class, 'listExchangeRates']);
            Route::post('/', [SchoolAdminController::class, 'createExchangeRate']);
            Route::put('{id}', [SchoolAdminController::class, 'updateExchangeRate']);
        });
    });

    // ===== 3. 教师端 =====
    Route::prefix('teacher')->middleware(['auth:sanctum', 'role:teacher'])->group(function () {
        Route::get('mode', [TeacherController::class, 'getMode']);
        Route::post('mode', [TeacherController::class, 'setMode']);
        Route::get('my-classes', [TeacherController::class, 'myClasses']);
        Route::post('switch-class', [TeacherController::class, 'switchClass']);
        Route::get('dashboard', [TeacherController::class, 'dashboard']);

        Route::prefix('students')->group(function () {
            Route::get('/', [TeacherController::class, 'listStudents']);
            Route::post('import', [TeacherController::class, 'importStudents']);
            Route::put('{id}', [TeacherController::class, 'updateStudent']);
        });

        Route::prefix('classroom')->group(function () {
            Route::get('display', [TeacherController::class, 'classroomDisplay']);
            Route::get('messages', [TeacherController::class, 'pollClassroomMessages']);
            Route::post('messages', [TeacherController::class, 'sendClassroomMessage']);
        });

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

        Route::prefix('pets')->group(function () {
            Route::get('types', [TeacherController::class, 'getPetTypes']);
            Route::get('overview', [TeacherController::class, 'classPetsOverview']);
            Route::get('{studentId}', [TeacherController::class, 'getPet']);
            Route::post('{studentId}/switch', [TeacherController::class, 'switchPet']);
            Route::post('{studentId}/feed', [TeacherController::class, 'feedPet']);
            Route::post('{studentId}/rename', [TeacherController::class, 'renamePet']);
        });

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

        Route::prefix('class')->group(function () {
            Route::get('/', [TeacherController::class, 'classInfo']);
            Route::post('switch-series', [TeacherController::class, 'switchSeries']);
        });

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

        Route::prefix('notices')->group(function () {
            Route::get('/', [TeacherController::class, 'listNotices']);
            Route::post('/', [TeacherController::class, 'createNotice']);
            Route::put('{id}', [TeacherController::class, 'updateNotice']);
            Route::put('{id}/publish', [TeacherController::class, 'publishNotice']);
            Route::delete('{id}', [TeacherController::class, 'deleteNotice']);
        });

        Route::prefix('reports')->group(function () {
            Route::get('score-trend', [TeacherController::class, 'scoreTrend']);
            Route::get('pet-distribution', [TeacherController::class, 'petDistribution']);
            Route::get('student-progress', [TeacherController::class, 'studentProgress']);
            Route::get('export/{type}', [TeacherController::class, 'exportReport']);
        });

        Route::prefix('broadcasts')->group(function () {
            Route::get('/', [TeacherController::class, 'listBroadcasts']);
            Route::post('/', [TeacherController::class, 'sendBroadcast']);
            Route::get('{id}', [TeacherController::class, 'getBroadcast']);
        });

        Route::prefix('attendance')->group(function () {
            Route::get('today', [TeacherController::class, 'getTodayAttendance']);
            Route::post('start', [TeacherController::class, 'startAttendance']);
            Route::put('{studentId}', [TeacherController::class, 'setAttendance']);
            Route::post('{studentId}/mark-leave', [TeacherController::class, 'markManualLeave']);
            Route::post('{studentId}/mark-absent', [TeacherController::class, 'markManualAbsent']);
            Route::get('summary', [TeacherController::class, 'attendanceSummary']);
        });

        Route::prefix('homework')->group(function () {
            Route::get('/', [TeacherController::class, 'listHomework']);
            Route::post('/', [TeacherController::class, 'createHomework']);
            Route::get('{id}', [TeacherController::class, 'getHomework']);
            Route::put('{id}/close', [TeacherController::class, 'closeHomework']);
            Route::get('{id}/submissions', [TeacherController::class, 'getHomeworkSubmissions']);
            Route::get('{id}/qr-code', [TeacherController::class, 'getHomeworkQrCode']);
        });

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

        Route::prefix('grades')->group(function () {
            Route::get('/', [TeacherController::class, 'listGrades']);
            Route::post('/', [TeacherController::class, 'inputGrades']);
            Route::get('stats', [TeacherController::class, 'getGradeStats']);
            Route::get('distribution', [TeacherController::class, 'getGradeDistribution']);
        });

        Route::prefix('ai')->group(function () {
            Route::post('chat', [TeacherController::class, 'aiChat']);
            Route::get('commands', [TeacherController::class, 'getAiCommands']);
            Route::get('usage', [TeacherController::class, 'getAiUsage']);
        });

        Route::prefix('currency')->group(function () {
            Route::get('wallets', [TeacherController::class, 'listWallets']);
            Route::post('exchange', [TeacherController::class, 'exchangeCurrency']);
            Route::post('cross-exchange', [TeacherController::class, 'crossExchangeCurrency']);
        });

        Route::get('display-code', [DisplayController::class, 'getDisplayCode']);
        Route::post('display-code/refresh', [DisplayController::class, 'refreshDisplayCode']);
    });

    // ===== 4. 家长端 =====
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

    // ===== 5. 班级大屏 (班级码 Token) =====
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

    // ===== 6. 公共接口 =====
    Route::prefix('common')->group(function () {
        Route::get('pet-types', [StudentController::class, 'petTypes']);
        Route::get('evolution-stages', [StudentController::class, 'evolutionStages']);
        Route::get('score-categories', [StudentController::class, 'scoreCategories']);
    });

    // ===== 7. 企业微信 Webhook =====
    Route::prefix('wechat-work')->group(function () {
        Route::get('callback', [WechatWorkWebhookController::class, 'verify']);
        Route::post('callback', [WechatWorkWebhookController::class, 'receive']);
    });
}); // End API v1
