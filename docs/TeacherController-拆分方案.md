# TeacherController 拆分方案

> 产出：2026-09-13 ｜ 对象：`backend/app/Http/Controllers/Api/TeacherController.php`
> 现状：**2707 行 / 84 个方法 / 20 个业务域**，是全项目最大的单点可维护性瓶颈。
> 前置约束：**本机无 PHP**，无法跑 PHPStan / PHPUnit。本方案刻意设计为「可分批、可回滚、每批独立验证」，禁止一次性大爆改。

---

## 一、核心策略：先抽服务，不动路由

**不要**把方法搬到新 Controller —— 那会连带改 `routes/api.php` 里 80+ 条路由与前端契约，风险面爆炸。

**要做**的是：把业务逻辑下沉到 Service，控制器方法退化为「校验 → 调服务 → 返回响应」的薄壳。这样：

- `routes/api.php` **一行都不用改**（方法名与签名保持不变）
- 每个 Service 可独立写单测，不再依赖 HTTP 层
- 每批拆完都能单独提交、单独回滚

这与项目现有模式一致（`ScoreService`、`LeaderboardService` 已通过构造函数注入）。

---

## 二、第 0 步（必做前置）：权限辅助必须先抽

`teacherClassIds()` 与 `getAccessibleClassIds()` 是**几乎所有方法的公共依赖**（班级作用域校验 + api-bot 全班级放行）。如果不动它，每个新 Service 都要重复注入 `User` 并复制这段逻辑，拆分只会更乱。

**动作**：抽为 `App\Support\TeacherClassScope`（或 `App\Services\TeacherScopeService`），提供：

| 方法 | 职责 |
|---|---|
| `accessibleClassIds(User $teacher): array` | 合并原 `teacherClassIds` + `getAccessibleClassIds`，含 `isApiBot()` 全班级分支 |
| `assertAssigned(User $teacher, int $classId): void` | 替代散落的 `isAssigned` 判断，抛 403 |
| `resolveActiveClass(Request $request): ClassRoom` | 替代 `switchClass` / 设置 active_class 的重复校验 |

> 注意：`User::isApiBot()` 走的是 `settings.is_api_bot`（**免迁移**标记），拆分时不要改成新字段，否则要动迁移。

---

## 三、目标 Service 映射

按「现有分段注释」为界，逐段下沉。**行号来自方案产出时的快照，实施时请以方法名为准。**

| # | 目标 Service | 下沉方法（原行号） | 行数规模 |
|---|---|---|---|
| 1 | `ClassroomMessagingService` | `getMode`(97) `setMode`(109) `classroomDisplay`(145) `sendClassroomMessage`(257) `pollClassroomMessages`(347) | 142–395 |
| 2 | `DashboardService` | `dashboard`(399) `myClasses`(40) `switchClass`(72) | 37–93 · 396–473 |
| 3 | `StudentRosterService` | `listStudents`(477) `importStudents`(511) `updateStudent`(575) `createStudent`(586) `deleteStudent`(627) | 474–638 |
| 4 | `ScoreRuleService`（已存在，补入） | `listScoreRules`(858) `createScoreRule`(891) `updateScoreRule`(917) `deleteScoreRule`(928) | 817–938 |
| 5 | `ScoreService`（已存在，补入） | `scoreSummary`(642) `scoreHistory`(769) `recentScores`(795) `undoScore`(822) | 639–816 |
| 6 | `PetService` | `classPetsOverview`(942) `getPet`(966) `feedPet`(992) `renamePet`(1026) `switchPet`(1353) `petCollection`(1489) | 939–1043 · 1350–1524 |
| 7 | `PetSeriesService` | `classInfo`(1272) `switchSeries`(1303)（含 `speciesPoolForSeries` 系列池逻辑） | 1350–1370 |
| 8 | `PkService` | `pkLeaderboard`(1102) `myPkStats`(1162) `challengePk`(1226) | 1096–1349 |
| 9 | `ShopService` | `listShopItems`(1528) `createShopItem`(1597) `updateShopItem`(1630) `deleteShopItem`(1654) `listRedemptions`(1667) `createRedemption`(1690) `approveRedemption`(1723) `rejectRedemption`(1786) `deliverRedemption`(1796) | 1525–1806 |
| 10 | `NoticeService` | `listNotices`(1810) `createNotice`(1829) `updateNotice`(1852) `publishNotice`(1863) `unpublishNotice`(1881) `deleteNotice`(1892) | 1807–1902 |
| 11 | `ReportService` | `scoreTrend`(1906) `petDistribution`(1942) `studentProgress`(1960) `exportReport`(1993) | 1903–2033 |
| 12 | `BroadcastService` | `listBroadcasts`(2037) `sendBroadcast`(2047) `getBroadcast`(2120) | 2034–2128 |
| 13 | `AttendanceService` | `getTodayAttendance`(2129) `startAttendance`(2157) `setAttendance`(2178) `markManualLeave`(2197) `markManualAbsent`(2209) `attendanceSummary`(2221) | 2129–2250 |
| 14 | `GradeService` | `listGrades`(2254) `inputGrades`(2281) `getGradeStats`(2319) `getGradeDistribution`(2338) | 2251–2360 |
| 15 | `AiAssistantService` | `aiConfig`(2361) `aiChat`(2381) `getAiUsage`(2479) `getAiCommands`(2517) | 2361–2527 |
| 16 | `CurrencyService`（已存在，补入） | `listExchangeRates`(2531) `createExchangeRate`(2559) `updateExchangeRate`(2581) `listWallets`(2599) `exchangeLogs`(2621) `exchangeCurrency`(2649) `crossExchangeCurrency`(2679) | 2528–2707 |

**保留在控制器内**（不抽）：

- `giveScore`(662) `batchGiveScore`(719) `giveScoreByRule`(745) —— 已委派 `ScoreService`，只需保留校验与响应组装
- `__construct`、`teacherClassIds` / `getAccessibleClassIds`（改为委派第 0 步的 scope 服务）

**预期结果**：控制器由 2707 行降至约 **400–500 行**（纯校验 + 响应组装）。

---

## 四、实施顺序（按「耦合度从低到高」）

| 批次 | 范围 | 为什么放这个位置 |
|---|---|---|
| **P0** | 第 0 步 `TeacherClassScope` | 所有后续批次的前置依赖 |
| **P1** | #10 Notice · #12 Broadcast · #13 Attendance · #14 Grade | 无跨域调用、无 SSE、无事务，先拿这 4 块建立拆分范式 |
| **P2** | #9 Shop · #11 Report · #8 Pk | 有状态流转（兑换审核）与聚合查询，中等复杂度 |
| **P3** | #6 Pet · #7 PetSeries · #4 ScoreRule · #16 Currency | 与宠物/积分数据层强耦合，且宠物有自动分配逻辑 |
| **P4** | #5 Score（只抽查询侧）· #15 AiAssistant | 计费与用量统计有外部依赖 |
| **P5** | #3 Student（含导入查重 + DB 事务）· #1 ClassroomMessaging（SSE + 缓存）· #2 Dashboard | 最重、最易回归，放最后单独做 |

> 每批**单独一个提交**，不要混批。

---

## 五、每批的验证清单（无 PHP 也照此执行）

在 CI 或具备 PHP 的环境执行：

```bash
cd backend

# 1. 语法
php -l app/Http/Controllers/Api/TeacherController.php
find app/Services -name '*.php' -exec php -l {} \;

# 2. 路由零变化（关键回归点）
php artisan route:list --path=api --json > /tmp/routes-after.json
git stash && php artisan route:list --path=api --json > /tmp/routes-before.json && git stash pop
diff /tmp/routes-before.json /tmp/routes-after.json   # 必须为空

# 3. 静态分析 + 测试
vendor/bin/phpstan analyse --level=5
php artisan test
```

**人工回归重点**（响应体形状不能变，前端依赖）：

- 积分：`scores/summary` · `scores/history/{id}` · `scores/give` · `scores/batch-give` · `scores/{id}/undo`
- 学生：`students` 列表的 `meta` 分页字段、`students/import` 的 `skipped_*` / `teacher_accounts`
- 商城：兑换状态机 `pending → approved/rejected → delivered`
- 考勤：`attendance/summary` 的四态统计（present/late/leave/absent）
- 宠物：`pets/class-overview`、`switch-pet` 后 `class.active_series` 落库

---

## 六、红线（拆错就是生产事故）

1. **响应结构不许变**：`{ data, message, meta }` 包装与字段名（含 snake_case）是前端契约，Service 只做数据，包装留在控制器
2. **权限校验不许绕过**：所有按班级过滤的查询必须继续经 `TeacherClassScope`，尤其 `isApiBot()` 的全班级放行分支
3. **事务边界不许碎**：`importStudents`、积分批量、兑换审批必须保持原 `DB::transaction` 粒度，不要把事务拆进两个 Service
4. **SSE 与缓存不许丢**：`DisplayEventService` 的广播调用与班级码缓存（`Cache::remember`）在拆分中必须原样保留
5. **别顺手改逻辑**：这是纯结构重构。发现疑似 bug 单独记录、单独提交，不夹在拆分里

---

## 七、若暂无 PHP 环境

可以只做**第一批（P0 + P1）**：`TeacherClassScope` + Notice/Broadcast/Attendance/Grade 四个 Service。这四块：

- 无事务、无 SSE、无外部 API
- 方法数 20 / 84，行数约 470 / 2707，收益约 17%
- 出错时影响面可控（通知 / 广播 / 考勤 / 成绩），不影响积分与宠物主链路

**除 P0+P1 之外，建议先补齐 PHP 环境或依赖 CI 通过再继续。**
