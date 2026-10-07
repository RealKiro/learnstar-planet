# backend-go — 学宠星球 Go 后端（核心 MVP）

学宠星球（LearnStar Planet）班级管理与学生激励系统的 Go 重写版后端。本目录与既有 Laravel 后端（`backend/`）**并存**，是独立的核心 MVP 实现，不破坏原项目。

## 范围（核心 MVP）

- 认证：管理员/教师登录（JWT + bcrypt）、`/me`、修改密码，以及**会话三件套** `auth/{logout,refresh,bindings}`（登出/刷新按 `jti` 撤销旧令牌，绑定列表按学校开关取交集）
- 学校、班级、学生管理（管理员）
- 教师管理（管理员创建/编辑/重置密码）
- 积分规则（默认规则幂等播种 + 班级自定义规则）
- 积分发放 / 批量发放 / 撤销 / 汇总 / 明细
- 宠物：班级总览、详情、喂食、重命名、切换物种
- 排行榜：总分 / 本周 / 宠物等级
- 公告（通知 CRUD + 发布/取消发布）
- 商城（学校级/班级级商品 CRUD + 兑换状态机 pending → approved/rejected → delivered）
- 考勤（今日考勤、开始考勤、单生状态、请假、缺勤、统计）
- 教室广播（banner/popup/fullscreen 发送与查询；发送时向大屏事件总线发布 `broadcast` 事件）
- 多币种（钱包、汇率 CRUD、积分兑换、钱包互兑、兑换流水、消费）
- 课表（科目/节次学校级共享、排课读写、教师提交修改申请 → 管理员审核、管理员直改、CSES YAML 导出）
- 管理端全校规则/商品（学校级共享 CRUD，带 `scope`/`class_name` 标注，班级级 ID 单条操作一律 404）
- 教师端我的班级 / 切换班级 / 模式（`active_class_id`、`display_mode` 存入用户 `settings`）
- **班级大屏（第三端）**：班级码（`LS` + 年级 + 班号，确定性生成）管理与重置；班级码登录（`auth/class/login` 发 `class_` token 与 `display/login` 发 `disp_` token，两者与 JWT 完全独立）；大屏只读接口（初始全量、课表、CSES 导出、总览看板、学生列表、积分规则、宠物总览、排行榜、班级设置）；**实时推送**（DB 事件总线 + `display/poll` 轮询降级 + `display/sse` 长连接）；大屏写操作（加减分、商品列表、兑换、学生间转赠、宠物切换、系列切换、PK 榜）
- **PK 与宠物系列**：教师端 PK 排行榜 / 挑战 / 我的战绩，教师端与大屏端的整班宠物系列切换（10 个系列物种池与 Laravel 逐项一致）
- **管理端运维**：报表（全校概览 / 按年级 / 按班级 / 班级码登录日志）、系统诊断 / 状态 / 日志 / 修复、批量账号（重置密码 / 批量删除）、教师批量创建与 CSV 导入（含模板下载）、学生导入与批量删除 / 转班、班级批量创建与教师分配（`class_room_teachers`）、学年升级（预览 → 执行，不可逆）、学校 LOGO 上传
- **AI 配置与对话**：管理端 AI 中心（设置读写 / 开关 / 用量 / 模型列表 / 连通性测试）、教师端 AI 助教（配置态 / 对话 / 预设命令 / 用量）、大屏端 AI（开关检查 / 对话）；多供应商分发（claude / google / qwen / mcp 专属 + 其余 OpenAI 兼容）、`api_base` 覆盖、本地精确计费与用量累加

- **登录别名 / 报表与课表导出 / AI 官方账单**：`auth/teacher/login`、`auth/admin/login` 两个**角色限定**的账号密码登录别名；教师报表导出（`reports/export/{scores|pets|attendance}`）与管理员单班 / 全校课表导出（**CSV 替代 xlsx**，见文末差异说明）；AI 供应商官方用量 / 余额直查（OpenAI / DeepSeek / Moonshot / SiliconFlow / OpenRouter / New API·One API 中转 / 本地精确计费共 7 个驱动）

- **认证会话 / 管理端人员 / 教师端学生与图鉴**：`auth/{logout,refresh,bindings}`（JWT `jti` 撤销名单表 `revoked_tokens` + 第三方绑定表 `third_party_bindings`）；管理端 `GET|POST /admin/students`、`GET /admin/classes/:id`、`PUT /admin/teachers/:id/classes`、`GET /admin/teachers/:id/password`（明文密码列 `users.plain_password`）；教师端 `GET /teacher/class`、`POST /teacher/students`、`POST /teacher/students/import`、`PUT|DELETE /teacher/students/:id`、`GET /teacher/pets/:studentId/collection`（图鉴表 `pet_collections`）、`POST /teacher/scores/give-by-rule/:ruleId`

> 上一批测试：`internal/services/new_routes_test.go`（服务层：会话撤销/绑定交集/列表筛选分页/班级详情/分配 replace/明文密码/学生 CRUD 与导入/图鉴/按规则加减分）与 `internal/router/session_routes_test.go`（15 条路由注册、无 token 401、登出/刷新后旧 token 401、以及各接口端到端行为）。
> 测试：本批新增 `internal/services/third_party_test.go`（服务层：企微 token 缓存/部门成员、钉钉与飞书 provider、平台选项、微信/QQ/人人通分支、免注册建号、bind-after-scan、通讯录导入）与 `internal/router/third_party_routes_test.go`（14 条路由注册、鉴权、端到端行为）；全部使用内存 SQLite + `httptest` 假上游，**不访问外网**。

- **第三方扫码登录 / 绑定 / 通讯录导入**：`auth/third-party/{auth-url,options,login}`、`auth/teacher/login/{wechat,wechat-work,qq,renren}`、`auth/teacher/bind-after-scan`、`auth/bind/{platform}`、`auth/unbind/{platform}`、`admin/wechat-work/{contacts,import}`、`admin/third-party/{contacts,import}` 共 **14 条**；企微 / 钉钉 / 飞书三个 provider + 企微通讯录服务（access_token 落表缓存、OAuth code 换 userid、部门/成员按 userid 去重）；扫码上下文落表 `temp_binding_contexts`（10 分钟有效、一次性消费）。

- **积分分类字典 + 企微回调 / 请假同步（最后一批）**：`GET common/score-categories`（公开字典，键值与 `ScoreRuleService::CATEGORY_LABELS` 同源）；`GET|POST wechat-work/callback`（企微回调验签 + AES-256-CBC 解密：GET 回显 echostr 纯文本、POST 接收 `sys_approval_change` 事件并同步请假 → 考勤）；`services.WechatWorkAttendanceService`（审批查询 → 家长企微绑定匹配学生 → 落 `wechat_work_leave_records` → 写考勤 `leave/wechat_work`，含点名时补同步）；配套 CLI `-sync-wechat-work-leave`。

- **宠物切换图鉴归档 / 考勤 `leave_record` / 登录限流（本批）**：教师端 `POST /teacher/pets/:student_id/switch` 与教室端 `POST /display/pets/switch` 都写 `pet_collections` —— 教师端**归档旧物种 + 恢复目标物种进度**（切回时等级/经验/心情全保留），教室端**归档旧物种 + 只换物种**（等级/经验/心情保留）；两条路径的差异是 Laravel 的原生行为，非移植差异。`GET /teacher/attendance/today` 每条考勤附带 `leave_record`（`{sp_no,leave_type,reason}`，无则 `null`），且 `check_in_time` 改为 Laravel `toDateTimeString` 形态的字符串。`auth/teacher/login` + `auth/admin/login`（`throttle:6,1`）与 `auth/class/login` + `display/login`（`throttle:10,1`）挂上**进程内固定窗口**限流（429 + `Retry-After` + `{"message":"Too Many Attempts."}`；**仅单实例有效**，多副本需换共享存储）。
未包含（不在当前范围）：作业、测验、题库；完整清单见文末「尚未移植清单」（**Laravel 有、Go 无的接口已清零**，剩余项全部是前端侧需适配的路径差异与前端死链）。

## 技术栈

- **Go 1.27** + **Gin**（HTTP 路由）+ **GORM**（ORM）
- 认证：JWT（`golang-jwt/jwt/v5`）+ bcrypt（`golang.org/x/crypto`）
- 数据库：SQLite（默认，纯 Go 驱动 `github.com/glebarez/sqlite`，无需 CGO）、MySQL、PostgreSQL
- 配置：环境变量 + `.env`（`godotenv`）

## 目录结构

```
backend-go/
├── main.go                       # 入口
├── internal/
│   ├── auth/                     # JWT 生成/解析
│   ├── config/                   # 环境变量加载
│   ├── database/                 # 连接/迁移/种子数据
│   ├── handlers/                 # HTTP 处理器（auth/admin/teacher）
│   ├── middleware/               # JWT 鉴权 + 角色隔离
│   ├── models/                   # GORM 模型 + 宠物等级/阶段逻辑
│   ├── router/                   # 路由注册
│   ├── services/                 # 业务服务（评分/宠物/排行榜/规则/认证/管理）
│   └── util/                     # 时区等工具
└── .env.example
```

## 快速开始

```bash
cd backend-go
cp .env.example .env   # 可选，不配置也能启动
go run .
```

默认监听 `http://localhost:8080`，首次启动自动创建 SQLite 数据库、默认学校与管理员：

- 管理员：`admin` / `admin123456`
- 默认学校代码：`learnstar`

> ⚠️ 上线前务必修改 `ADMIN_PASSWORD` 与 `JWT_SECRET`。

健康检查：

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## 数据库

默认 `DB_DRIVER=sqlite`，数据文件 `data/learnstar.db`。切换 MySQL/PostgreSQL 见 `.env.example` 中的注释示例。表结构由启动时的 `AutoMigrate` 自动创建。较早批次新增两张表：`temp_binding_contexts`（扫码临时上下文，替代 Laravel 的 Cache）与 `wechat_work_tokens`（企微 access_token 缓存，替代 Laravel 的 Cache）。最后一批新增 `wechat_work_leave_records`（企微请假记录，字段照 Laravel 迁移 `2026_07_13_000002`：`school_id/class_id/student_id/parent_wework_userid/student_name_from_wework/sp_no/leave_start_date/leave_end_date/leave_type/reason/approve_status/approved_at/raw_data/synced_at`，`sp_no` 唯一），并给 `students` 补上 `parent_id` 列（家长匹配用，见下文差异说明）。

## 构建与验证

```bash
go build ./...
go vet ./...
go test ./... -count=1
gofmt -l .          # 应无输出
```

CLI（企业微信请假同步，等价 Laravel `php artisan attendance:sync-wechat-leave`）：

```bash
go build -o learnstar-go .
./learnstar-go -sync-wechat-work-leave                    # 全部 status=active 学校，日期默认今天
./learnstar-go -sync-wechat-work-leave -school-id=1 -date=2026-09-20
# 输出：{"synced":0,"skipped":0,"failed":0}（执行后退出，不启动 HTTP 服务）
```

## API 概览

所有响应使用统一 JSON 信封：`{"data": ..., "message": "ok"}`；错误为 `{"message": "..."}` 并带对应 HTTP 状态码。认证使用 `Authorization: Bearer <token>`。角色隔离：`school_admin` 访问 `/admin/*`，`teacher` 访问 `/teacher/*`。

### 公开

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| POST | `/api/v1/auth/class/login` | 班级码登录（换 `class_` token，24 小时；无效码 401） |
| POST | `/api/v1/auth/teacher/login` | 教师账号密码登录别名（仅 `role=teacher`；角色不符或密码错一律 401「账号或密码错误，请核对后重试」） |
| POST | `/api/v1/auth/admin/login` | 管理员账号密码登录别名（仅 `role=school_admin`，同上文案） |
| GET | `/api/v1/auth/third-party/auth-url` | 第三方扫码授权 URL（`?school_id=&redirect_uri=`；`redirect_uri` 缺省用 `APP_URL + /login`）。无学校 → 400「系统尚未初始化」；学校未配置平台 → 400「未配置第三方平台，请在后台学校设置中选择」；成功 `{platform, auth_url}` |
| GET | `/api/v1/auth/third-party/options` | 登录页可选平台（`{key,label,icon,color}`；学校开关与 `platforms()` 取交集，空则默认企业微信/微信/QQ；颜色表 `wechat_work #2B7CE9` / `dingtalk #0089FF` / `feishu #3370FF` / `wechat #07C160` / `qq #12B7F5` / `renren #FF6A00` / 兜底 `#7c3aed`） |
| POST | `/api/v1/auth/third-party/login` | 第三方平台扫码回调（`{code, state?}`，`state` 为 schoolId）。未配置平台 / 上游失败 → 400 + 异常文案；成功 `{status:'logged_in', token, user}`（已绑定分支无 token） |
| POST | `/api/v1/auth/teacher/login/wechat` | 微信扫码登录（`{openid, unionid?, nick?, avatar?}`）。已绑定 → `{status:'logged_in', user}`（**无 token**，Laravel 原样行为）；未绑定 → `{status:'need_binding', temp_token, openid, unionid}` |
| POST | `/api/v1/auth/teacher/login/wechat-work` | 企业微信扫码登录（`{userid?, code?, state?, nick?, avatar?}`）。无 `userid` 且有 `code` 时用 code 换 userid；都没有 → 400「无法获取企业微信用户身份」；未绑定**免注册建教师账号**并返回 token（默认密码 `ls123456`）；已绑定 → `{status, user}`（无 token） |
| POST | `/api/v1/auth/teacher/login/qq` | QQ 扫码登录（`{openid, nick?, avatar?}`）。未绑定 → `{status:'need_binding', temp_token, openid}`；已绑定 → 无 token |
| POST | `/api/v1/auth/teacher/login/renren` | 人人通空间登录（`{user_id, nick?, avatar?}`）。未绑定 → `{status:'need_binding', temp_token, platform_id}`；已绑定 → 无 token |
| POST | `/api/v1/auth/teacher/bind-after-scan` | 扫码后绑定已有教师账号（`{temp_token, username, password, platform, platform_id, unionid?, nick?, avatar?}`）。账号密码错 → **HTTP 200** + `{data:{status:'error', message:'账号或密码错误，请核对后重试'}}`；成功 → `{data:{status:'bound', user}}`；`temp_token` 一次性消费 |
| GET | `/api/v1/common/score-categories` | 积分分类字典（公开；`{data:[{id,name,icon,sort}]}`，**不带 message**）。id/name 取自 `ScoreRuleService::CATEGORY_LABELS` 的 Go 真源 `services.CategoryLabels`，`sort` 从 1 递增，图标缺项兜底 📌 |
| GET | `/api/v1/wechat-work/callback` | 企微回调 URL 验签（公开）。`msg_signature/timestamp/nonce/echostr` 任一为空、配置缺失、签名错、解密失败一律返回**空字符串**；成功返回解密出的 echostr **纯文本**（非 JSON） |
| POST | `/api/v1/wechat-work/callback` | 企微审批事件接收（公开，`?school_id=`）。恒返回 `{"errcode":0,"errmsg":"ok"}`（HTTP 200）；`Event`/`ChangeType` 命中 `sys_approval_change`、`ApprovalInfo.SpNo` 非空、`SpStatus=2`、`school_id>0` 时才同步请假并补考勤 |

### 认证（需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/change-password` | 修改密码（`{old_password, new_password}`） |
| POST | `/api/v1/auth/logout` | 登出（撤销当前令牌的 `jti` → 旧 token 立即 401；响应「已登出」） |
| POST | `/api/v1/auth/refresh` | 刷新令牌（撤销旧 token 并签发新 token；响应 `{data:{token}}`） |
| GET | `/api/v1/auth/bindings` | 第三方绑定列表（`{platform,label,icon,bound,nick}`；仅学校启用的平台，未配置时默认企业微信/微信/QQ） |
| POST | `/api/v1/auth/bind/:platform` | 登录后绑定第三方账号（`{platform_id, nick?, avatar?}`）。该 `platform + platform_id` 已被任意用户绑定 → 422「该第三方账号已被其他用户绑定」；成功「绑定成功」 |
| DELETE | `/api/v1/auth/unbind/:platform` | 解绑当前用户该平台的第三方账号（「解绑成功」，无绑定也成功） |

### 管理端 `/api/v1/admin/*`（`school_admin`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET / PUT / POST | `/school` | 查看 / 更新学校（`settings` 为**整体替换**语义，非法类型 422；成功「学校信息已更新」；`GET` 把 `settings` 作为**对象**返回） |
| GET / POST | `/classes` | 班级列表 / 创建 |
| PUT / DELETE | `/classes/:id` | 更新班级（`pet_series` 白名单 `all|cosmic|pokemon|cute|treasure|mythic`，**按键合并**进 `settings`；非法 → 422）/ 删除班级 |
| GET | `/classes/:id` | 班级详情（含 `settings` 对象、`teacher` 与 `students` 关联；非本校/不存在 → 404「班级不存在」） |
| GET / POST | `/students` | 全校学生分页列表（`?search=&class_id=&grade=&status=`，默认 `status=active`、`per_page=50`、`id desc`；每行附 `class_name/class_grade/pet_species/pet_level/pet_name`，响应含 `meta`）/ 单建学生（自动分配萌宠，201「学生「X」已添加，已自动分配萌宠」） |
| PUT / DELETE | `/students/:id` | 更新 / 删除学生 |
| GET / POST | `/teachers` | 教师列表 / 创建 |
| PUT | `/teachers/:id` | 更新教师 |
| DELETE | `/teachers/:id` | 删除教师（解除班级关联 + 硬删除；机器人账号 403；文案「教师账号已删除」） |
| POST | `/teachers/:id/reset-password` | 重置教师密码（同步写 `plain_password`，`password_changed=false`） |
| PUT | `/teachers/:id/classes` | 设置教师所带班级（**replace 语义**，`{assignments:[{class_id,role,subject?}]}`；`head_teacher` 同步 `class_rooms.teacher_id`，移除时置空） |
| GET | `/teachers/:id/password` | 查看教师明文密码（无记录时自动生成 8 位并写回，同 Laravel） |
| GET / POST | `/exchange-rates` | 汇率列表 / 创建 |
| PUT | `/exchange-rates/:id` | 更新汇率 |
| GET / POST | `/score-rules` | 全校规则列表（含班级级，带 `scope`/`class_name`）/ 创建学校级规则 |
| PUT / DELETE | `/score-rules/:id` | 更新 / 删除学校级规则（班级级 ID 一律 404） |
| GET / POST | `/shop-items` | 全校商品列表（`?currency_type=`）/ 创建学校级商品 |
| PUT / DELETE | `/shop-items/:id` | 更新 / 删除学校级商品（班级级 ID 一律 404） |
| GET | `/timetable/changes` | 全校课表修改申请列表（`?status=pending`） |
| POST | `/timetable/changes/:id/approve` | 通过申请（应用快照，其余 pending 自动作废） |
| POST | `/timetable/changes/:id/reject` | 驳回申请 |
| GET / PUT / POST | `/classes/:id/timetable` | 查看 / 直接保存某班课表（即时生效；PUT+POST 双注册，文案「课表已保存并即时生效」） |
| GET | `/classes/:id/display-code` | 查看某班班级大屏码 |
| POST | `/classes/:id/display-code/refresh` | 刷新某班班级码 |
| POST | `/classes/reset-display-codes` | 批量重生成全校班级码（可传 `prefix`） |
| POST | `/timetable/import-csv` | 课表 CSV 批量导入（multipart `file` + `dry_run`，默认仅预览） |
| GET / POST | `/timetable/unavailabilities` | 教师不可用时段（replace 语义） |
| POST | `/timetable/check-conflicts` | 排课冲突检查（教师跨班撞课 + 教师不可用时段） |
| GET / PUT / POST | `/classes/:id/teacher-assignments` | 任课设置（replace 语义，`assignments` 必填） |
| POST | `/timetable/generate` | 单班自动排课（纯计算不落库，返回 `{success, warnings, entries}`） |
| POST | `/timetable/generate-school` | 全校智能排课（依据任课表与教师不可用时段；`commit=true` 才事务落库并作废相关 pending 申请） |
| GET | `/reports/overview` | 全校概览（班级/教师/学生数 + 本月/上月积分与环比 `score_trend_percent`） |
| GET | `/reports/by-grade` | 按年级汇总（班级数 / 学生数 / 平均分 / 总分，空年级归入「未分年级」） |
| GET | `/reports/by-class` | 按班级汇总（班主任 / 学生数 / 本月积分） |
| GET | `/display-login-logs` | 班级码大屏登录日志（`?class_id=&ip=&date=&page=`，每页 50） |
| GET | `/system/diagnose` | 表结构自检（`data` + `has_issues` + `message`） |
| GET | `/system/status` | 版本信息 + 库表清单 + 上传目录等运行态 |
| GET | `/system/logs` | 日志（读环境变量 `LOG_FILE` 指定的文件；未配置 → 空列表 + 说明） |
| POST | `/system/repair` | 系统修复（GORM AutoMigrate，幂等） |
| POST | `/accounts/batch-reset-password` | 批量重置教师密码（`{role:"teacher", ids, password?}`，默认 `ls123456`） |
| POST | `/accounts/batch-delete` | 批量删除教师账号（API 机器人账号受保护） |
| GET | `/system/status` | 版本信息 + 库表清单（Go 无 `migrations` 表，`migrations` 恒为空） |
| POST | `/teachers/import` | CSV 导入教师（multipart `file` + `dry_run`，默认 true=预览；支持 `姓名/年级团队/科目/密码/手机号` 或英文字段名） |
| GET | `/teachers/template-csv` | 教师导入 CSV 模板（带 UTF-8 BOM，`Content-Type: text/csv; charset=UTF-8`） |
| POST | `/students/import` | 学生导入（JSON `{students:[{name,class_name,gender,student_no}]}`，**也支持** multipart CSV 上传） |
| POST | `/students/batch-delete` | 批量删除学生（`{student_ids}`，硬删除并级联清理宠物/积分/审计日志） |
| POST | `/students/batch-move` | 批量转班（`{student_ids, target_class_id}`） |
| POST | `/classes/batch-create` | 批量建班（`{grade, count(1-20), year}`，命名「年级（N）班」） |
| POST | `/classes/:id/assign-teacher` | 分配班级教师（`{teacher_id, role}`；`head_teacher` 同步 `class_rooms.teacher_id`） |
| DELETE | `/classes/:id/remove-teacher` | 移除班级教师（`{teacher_id}` 或 query） |
| GET | `/grade-upgrade/preview` | 学年升级预览（dry-run，不落库） |
| POST | `/grade-upgrade/execute` | 执行学年升级（**不可逆**） |
| POST | `/school/logo` | 上传学校 LOGO（multipart 字段名 `logo`；**不做缩放裁剪**） |
| GET | `/classes/:id/timetable/export-excel` | 单班课表导出 CSV（文件名 `<班级名>-课表.csv`；非本校班级 404「班级不存在」） |
| GET | `/timetable/export-excel` | 全校课表导出 CSV（各班网格按 grade/name 顺序拼接；无班级 404「暂无班级可导出」） |
| POST | `/ai/provider-official` | AI 供应商官方用量 / 余额查询（`{provider_id}`；缺参 / 未建配置 / 无 Key 均 422） |
| GET | `/wechat-work/contacts` | 拉取企业微信通讯录（`{departments, members}`；成员按 userid 去重、只遍历顶级部门 `fetch_child=1`、`department_names` 映射）。上游异常 → 400 + 文案；学校记录缺失 → 404「未找到学校」 |
| POST | `/wechat-work/import` | 从企业微信通讯录批量导入（`{teachers:[{name,mobile?,email?}], students:[{name,class_id,gender?}]}`）。教师按手机号（去空格/连字符）与「姓名=用户名」查重后跳过；学生同班同名跳过并自动分配萌宠；返回「已导入 N 名教师、M 名学生…」+ `{created_teachers, created_students, skipped_teachers, skipped_students, teacher_accounts}` |
| GET | `/third-party/contacts` | 拉取当前学校所选平台的通讯录（同企微结构 + 附带 `platform` 标识）。未配置平台 / 上游异常 → 400 + 文案 |
| POST | `/third-party/import` | 从当前学校所选平台导入教师与学生（与 `/wechat-work/import` 同一套逻辑，同 Laravel 委派） |

### 教师端 `/api/v1/teacher/*`（`teacher`）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/dashboard` | 概览 |
| GET | `/students` | 学生列表（`?search=` 命中姓名或学号、`?page=` 分页每页 **50**（同 Laravel `paginate(50)`，`per_page` 不生效）、**不过滤 status**、按 `name` 升序；行内含 `pet_species`/`pet_level`/`pet_name`，响应 `{data, meta}`） |
| POST | `/students` | 添加学生（`{name,class_id,gender?,student_no?}`；不在自己管辖班级 → 403「只能在自己管理的班级添加学生」；成功 201「学生「X」已添加」） |
| POST | `/students/import` | 批量导入（`{students:[{name,class_name,gender?,student_no?}]}`；按班级名在管辖范围内定位，同班学号/同名与跨班同学号跳过，返回 `{message, data:{imported_count, skipped}}`） |
| PUT | `/students/:id` | 更新学生（仅 `name`/`gender`/`student_no`；越权/不存在 404） |
| DELETE | `/students/:id` | 删除学生（硬删除并级联清理宠物/积分/审计日志；「学生「X」已删除」） |
| GET | `/class` | 班级信息卡（首个可管辖班级：`id/name/grade/student_count/total_score/class_points/settings/display_code`；无可管辖班级 → 400「没有可管理的班级」） |
| GET | `/pets/:student_id/collection` | 学生宠物图鉴（`student_id/student_name/total_score/unlock_slots/class_series/active_species/collection`；当前激活宠物自动补录；越权 404） |
| POST | `/scores/give-by-rule/:ruleId` | 按规则加减分（`{student_id}`；规则越权/不存在 404；「已按规则「X」处理」，走与 `scores/give` 同一套事务/审计/宠物经验/大屏事件） |
| GET / POST | `/scores/rules` | 积分规则列表 / 创建（Laravel `scores` 前缀组内的 `rules`） |
| PUT / DELETE | `/scores/rules/:id` | 更新 / 删除规则 |
| GET | `/scores/summary` | 积分汇总 |
| GET | `/scores/recent` | 最近积分记录 |
| POST | `/scores/give` | 发放积分（`{student_id, points, reason, rule_id?}`，请求字段名 `points` 与 Laravel `input('points')` 对齐；`points` 为 0 → 422「积分变动不能为 0」） |
| POST | `/scores/batch-give` | 批量发放（`{student_ids, points, reason, rule_id?}`，同上） |
| POST | `/scores/:id/undo` | 撤销 |
| GET | `/scores/history/:student_id` | 学生积分明细 |
| GET | `/pets/overview` | 宠物班级总览 |
| GET | `/pets/:student_id` | 宠物详情 |
| POST | `/pets/:student_id/feed` | 喂食 |
| POST | `/pets/:student_id/rename` | 重命名 |
| POST | `/pets/:student_id/switch` | 切换物种 |
| GET | `/leaderboard/total` | 总分排行榜（默认取教师可管理班级的第一个；`?class_id=` 可选，`?limit=` 默认 20；无班级 → `{data:[]}`） |
| GET | `/leaderboard/weekly` | 本周排行榜（同上） |
| GET | `/leaderboard/pet-level` | 宠物等级排行榜（同上） |
| GET / POST | `/notices` | 公告列表 / 创建 |
| PUT | `/notices/:id` | 更新公告 |
| PUT | `/notices/:id/publish` `/unpublish` | 发布 / 取消发布 |
| DELETE | `/notices/:id` | 删除公告 |
| GET / POST | `/shop/items` | 商品列表 / 创建 |
| PUT / DELETE | `/shop/items/:id` | 更新 / 删除商品 |
| GET / POST | `/shop/redemptions` | 兑换列表 / 发起兑换 |
| PUT | `/shop/redemptions/:id/approve` `/reject` `/deliver` | 兑换审核 / 驳回 / 发放 |
| GET | `/attendance/today` | 今日考勤 |
| POST | `/attendance/start` | 开始今日考勤 |
| PUT | `/attendance/:student_id` | 设置学生考勤状态 |
| POST | `/attendance/:student_id/mark-leave` `/mark-absent` | 标记请假 / 缺勤 |
| GET | `/attendance/summary` | 考勤统计 |
| GET / POST | `/broadcasts` | 广播列表 / 发送 |
| GET | `/broadcasts/:id` | 单条广播 |
| GET | `/currency/wallets` | 学生钱包列表 |
| GET | `/currency/exchange-logs` | 兑换流水（分页） |
| POST | `/currency/exchange` | 积分换币 |
| POST | `/currency/cross-exchange` | 钱包互兑 |
| GET / POST | `/exchange-rates` | 汇率列表 / 创建 |
| PUT | `/exchange-rates/:id` | 更新汇率 |
| GET | `/timetable` | 本班课表（科目 + 节次 + 排课） |
| POST | `/timetable` | 提交课表修改申请（不直接生效） |
| GET | `/timetable/changes` | 本班申请历史 |
| GET | `/timetable/my-schedule` | 我的周课表（`?teacher_name=`） |
| GET | `/timetable/export-cses` | 导出 CSES YAML（ClassIsland 可导入） |
| GET | `/my-classes` | 我的班级（机器人返回本校全部启用班级，`role=api_bot`） |
| POST | `/switch-class` | 切换当前激活班级（未分配 → 403） |
| GET / POST | `/mode` | 当前模式（默认 `classroom_display`）/ 切换模式 |
| GET | `/display-code` | 查看本班班级大屏码（`?class_id=` 或当前激活班级） |
| POST | `/display-code/refresh` | 刷新本班班级码 |
| POST | `/class/switch-series` | 整班切换宠物系列（**不扣分**，给全班发放 3 天免费自选） |
| GET | `/pk/leaderboard` | 同年级班级 PK 榜 |
| POST | `/pk/challenge` | 发起班级 PK 挑战（`target_class_id`） |
| GET | `/pk/my-stats` | 我的班级战绩（含名次） |
| GET | `/reports/score-trend` | 积分趋势（`?days=`，返回 `labels` + 「得分/扣分」datasets） |
| GET | `/reports/pet-distribution` | 宠物等级分布（`level`/`count`/`stage_name`） |
| GET | `/reports/student-progress` | 学生进度（`?student_id=` 单人近 50 条；否则全班近 10 条 + `trend`/`change`） |
| GET | `/reports/export/:type` | 报表导出 CSV（`scores|pets|attendance`，`?class_id=` 缺省取管辖首个班级；越权 403；未知 type → 200 文案） |
| GET | `/classroom/display` | 班级大屏数据（`?class_id=`；未选班级 → 400） |
| GET / POST | `/classroom/messages` | 课堂消息（GET 拉取广播/通知；POST 按 `type` 分流见下） |

### 班级大屏 `/api/v1/display/*`（班级码 token，非 JWT）
登录用班级码换取 `disp_` token（TTL 24 小时）；其余接口通过 `DisplayAuth` 中间件鉴权，token 取自 query `token` 或 `Authorization: Bearer`（仅接受 `disp_`/`class_` 前缀）。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/login` | 班级码登录（无效码 404；同码连续失败 5 次锁定 15 分钟 → 429；另有 `throttle:10,1` 限流 → 429） |
| GET | `/initial-data` | 大屏初始全量（学生 + 宠物 + 近 4 小时积分 + 生效广播；无宠物学生自动分配一只） |
| GET | `/timetable` | 本班课表（含 `today_weekday`） |
| GET | `/export-cses` | 免登录导出 CSES YAML（教室机课表软件直接拉取） |
| GET | `/dashboard` | 班级总览（总分 / 平均宠物等级 / 尖峰数 / 本周积分 / Top5 / 星光学生 / 最新动态） |
| GET | `/students` | 学生列表（含宠物名/物种/等级/emoji） |
| GET | `/scores/rules` | 积分规则（与教师端同源） |
| GET | `/pets/overview` | 宠物总览 |
| GET | `/leaderboard` | 班级总分榜 Top20 |
| GET | `/class-settings` | 班级设置（`pet_series`） |
| GET | `/shop-items` | 本班在售商品（**仅班级级商品**） |
| POST | `/redeem` | 学生兑换（扣分 + 生成已通过兑换记录；**不扣库存**，同 Laravel） |
| POST | `/transfer` | 学生间转赠（`from_id`/`to_id`/`amount` 1–100，双方须同班） |
| POST | `/switch-series` | 整班切换宠物系列（每人扣 20 分并发放 3 天免费自选） |
| POST | `/pets/switch` | 学生切换宠物（跨系列 422；免费自选/按等级扣分；旧物种归档进图鉴、目标物种标记激活，等级/经验/心情保留） |
| GET | `/pk/leaderboard` | 同年级班级 PK 榜 |
| GET | `/poll` | 轮询降级拉取事件（`?last_event_id=`，含 `events` / `last_event_id` / `server_time`） |
| GET | `/sse` | SSE 长连接（`?last_event_id=`；心跳 10s、最长 55s 后发 `reconnect` 帧） |
| POST | `/quick-score` | 大屏快捷加减分（`amount` 限 `-5/-3/-1/1/3/5`；返回 `total_score`） |
| POST | `/scores/give` | 大屏加减分（单次绝对值 ≥ 30 → 403「单次加减分超过 30 分，请使用教师账号登录操作」） |
| POST | `/scores/batch-give` | 大屏批量加减分（`student_ids` + `points`） |

## 关键业务规则（忠实移植自 Laravel 端）
- **宠物等级**：由学生 `TotalScore` 经等级阈值表映射（`internal/models/pet.go` 的 `LevelThresholds`），1–12 级；升级所需经验为 `(等级+1)*10`。
- **用户设置**：`active_class_id` / `display_mode` 存入 `users.settings`（JSON 文本列），读写语义同 Laravel `getSetting/setSetting`（按键合并、缺失取默认值）。切换模式时若带 `class_id` 且已分配则同步切换激活班级，未分配则静默忽略；`switch-class` 未分配返回 403。
- **惩罚规则标记（Go 端坑位）**：`score_rules.is_positive` / `is_active`、`shop_items.is_active` 等带 `gorm:"default:true"` 的布尔列，GORM 在 `Create` 时会跳过 `false` 零值而由 DB 施加默认 `true`。所有写入 `false` 的路径（默认规则播种、教师/管理员创建规则、管理员创建商品）都在创建后补一次显式 `Update` 回写，避免「惩罚规则被标记成奖励」「显式下架变上架」。
- **积分事务**：积分变动在事务内完成——创建 `Score` → 余额钳制为 `max(0, TotalScore)` → 写 `ScoreLog` → 同步宠物经验/等级。
- **默认积分规则**：`services.DefaultRules` 共 43 条（奖励 27 / 惩罚 16），按学校幂等播种（判据是「是否播种过」而非「集合是否为空」，同名不覆盖、教师删除不复活）。
- **排行榜**：Go MVP 直接查库（不依赖 Redis ZSET）。
- **班级可见性**：普通教师只看到 `teacher_id = 本人` 的班级；API 机器人教师可见本校全部班级。
- **课表修改两步流**：教师端 `POST /timetable` 只落 `timetable_change_requests`（存完整快照 JSON），**不改动生效课表**；管理员 `approve` 时才应用快照，并把该班其余 `pending` 申请自动置为 `rejected`（「已由更新的申请取代」）；`reject`/`approve` 均幂等，仅 `pending` 可被处理。管理员 `PUT /classes/:id/timetable` 直改即时生效，同时作废该班全部 `pending` 申请（防旧快照覆盖）。
- **课表保存语义**：科目/节次为**学校级** upsert（科目按 `school_id + name` 去重，颜色未传时保留库中原值、两者皆无才按名称从 15 色色板稳定取色；节次按 `school_id + period_index` 去重，未传 name 生成「第N节」）；某班排课为**先清后插**（整体覆盖该班，不影响其他班），只接受科目名已存在、`weekday` 1–7、`period_index` ≥ 1 的行，非法周次回退 `all`。
- **多币种**：汇率默认仅播种 `score → 科学币/读书币`；`Exchange` 为积分换币（`to_amount = floor(from_amount × rate)`，积分不足 422），`CrossExchange` 为钱包互兑（同币种 400、余额不足 422）。GORM 的 `Update` 会回写内存字段，故余额更新显式赋值避免重复累加。

- **班级码（大屏）**：确定性生成，不落库也无需刷新 —— `前缀 + 年级数字 + 班号数字`（如 `LS11`）。前缀取 `schools.settings.display_code_prefix`（须匹配 `^[A-Z]{2,4}$`，否则回退 `LS`）；年级支持「一~九年级 / 高一~高三」及纯数字，班号支持「(1)班 / 1班 / 第一班 / 一班」，无法解析或班号 > 9 时生成空串并跳过；班号兜底取同年级按 id 排序的序号。`regenerateAll` 返回 `{regenerated, skipped, conflicts}`，同校内重复码计入 `conflicts` 并跳过。
- **大屏鉴权**：班级码换 `disp_` token（TTL 86400 秒），token 与班级码映射存 **DB 表 `display_tokens`**（差异：Laravel 用 Cache + Redis，Go 端无缓存层）；`class_` 前缀同样接受并解析为 class_id。失败计数（5 次 / 15 分钟）为进程内并发安全 map —— **仅单实例有效，多副本部署需换共享存储**（`display/login` 的 `throttle:10,1` 同样是进程内固定窗口，见文末「宠物切换图鉴 / 考勤 `leave_record` / 登录限流」）。`display/export-cses` 与登录同属免 JWT 通道。
- **大屏初始数据**：`initial-data` 会为**没有宠物的活跃学生自动创建宠物**（物种从 `myth` 系列 18 种随机取，名称 `{学生名}的萌宠`，等级 1 / 经验 0 / 心情 80），因此该接口非纯只读。
- **大屏事件总线**：教师侧操作 → 写库 → 追加事件到该班列表 → 大屏 `poll`（或 `sse`）按 `last_event_id` 增量消费。事件 `{id, type, data, created_at}`，`type` ∈ `score_update / broadcast / notice / pet_update / refresh`；**每班 10 分钟 TTL、最多保留 200 条**，过期与超限在发布/消费时惰性清理（无定时任务）。差异：Laravel 存 Cache（Redis List），Go 存 **DB 表 `display_events`**，事件 id 用「该班未过期事件的 `seq`」表达同义语义；`clear()` 后 seq 从 1 重开（Laravel 的计数器键存活到 TTL 到期，边界同类）。
- **发布点接线**：积分发放/批量发放（事务提交后发）、广播发送、通知发布、宠物喂食、班级码刷新（教师端带 `old_code`/管理员端不带）各发对应事件。`quick-score`、`scores/give`、`scores/batch-give` 复用既有 `ScoreService`，因此相对 Laravel 的手写实现**多出**审计日志 `score_logs`、`score_update` 事件与宠物等级校正（已在代码注释写明）。
- **大屏加减分口径差异**：`quick-score` 的 `amount` 是**白名单** `-5/-3/-1/1/3/5`（其余 422），**没有** ±30 限制；`scores/give` 与 `scores/batch-give` 才有「单次绝对值 ≥ 30 → 403」。
- **SSE 与鉴权**：`display/sse` 帧格式（`id/event/data` + 空行）、心跳 10s、最长运行 55s 后发 `reconnect` 帧、响应头（`text/event-stream` / `no-cache` / `keep-alive` / `X-Accel-Buffering: no`）逐项对齐；token 无效时 Go 走 `DisplayAuth` 返回 **401 JSON**，而 Laravel 返回 200 + `event: error` 帧（当前前端课堂页只用 `poll`，未使用 `sse`）。
- **大屏写操作口径**：`display/shop-items` 与 `display/redeem` **只看班级级商品**（`class_id = 本班 AND is_active`），学校级商品（`class_id = NULL`）在大屏不可见也不可兑换；兑换**不校验也不扣减库存**（Laravel 原实现即如此），扣分走 `spendScore`、生成 `status = approved` 的兑换记录。`display/transfer` 为学生间转赠（`amount` 1–100、不得转给自己、双方必须同班，`from` 分数不足 → 400「积分不足」）。
- **宠物系列与切换**：`display/switch-series` 每人扣 **20 分**（任一活跃学生不足 → 400「积分不足：{前 3 名} 每人需要 20 积分」，无活跃学生 → 400「班级没有活跃学生」，非法系列 → 422），并把 `pet_series` 写入班级设置 + 给全班发放 **3 天**免费自选；教师端 `teacher/class/switch-series` 同流程但**不扣分**。`display/pets/switch`（参数名 **`pet_species`**）跨系列 → 422「只能领养当前类别「X」的宠物，不能跨类别领养」；同物种 → 422「当前已经是这只宠物啦」；无宠物学生**直接创建**（「🎉 新宠物已诞生！」）；有宠物时优先消耗免费自选，否则按 `SwitchCost(等级)` 扣分（不足 → 400「积分不足，更换宠物需 N 积分」），随后把旧物种归档进 `pet_collections`（`is_active=false`）、只换 `pets.species`（等级/经验/心情保留）、目标物种写图鉴并标 `is_active=true`。教师端 `teacher/pets/:student_id/switch` 语义不同：先查目标物种的图鉴进度，**切回时恢复等级/经验/心情**（无记录则初始 1/0/80），并在扣分后写 `pets.last_switched_at`；两端各自的语义都与 Laravel 对应实现逐条一致。
- **本批差异**：免费自选机会由 Laravel 的 Cache 键改为落库表 `pet_free_picks`（TTL 3 天）；`class_` token 与 `disp_` 共用 `display_tokens` 表；`pk/challenge` 的挑战记录 Laravel 只写不读，Go 未落库；校验类错误文案统一为中文（状态码一致）。（原先的「切换宠物不归档旧物种」「`auth/class/login` 的 `throttle:10,1` 未移植」两条已在本批补齐，见文末「宠物切换图鉴 / 考勤 leave_record / 登录限流」。）
- **教师报表口径**：`score-trend`（`days` 钳制 1–365，标签 `m/d`，datasets 为「得分」/「扣分」，扣分取绝对值）与 `pet-distribution`（按等级分组，带 `stage_name`）照 Laravel 逐字对齐；`student-progress` 有 `student_id` 返回单人近 50 条历史，否则返回全班活跃学生每人近 10 条 + `change = 前 5 条和 − 其余和`、`trend` 按 `>5 / <-5` 判 up/down/stable。
- **课堂消息分流**：`POST /teacher/classroom/messages` 按 `type` 分流 —— `banner|popup|fullscreen` 写广播表并发布 `broadcast` 事件（`voice` 默认 true、`display_seconds` 默认 10、`status=sent`），返回「广播已发送」；`info|homework|event|urgent` 写通知表并发布 `notice` 事件（`urgent` 默认标题「紧急通知」，其余「通知」），返回「通知已发布」。**越权预检返回 403 先于参数校验 422**（同 Laravel 顺序）。`classroom/display` 与 `classroom/messages` 未指定 `class_id` 且无激活班级 → 400「请先选择班级」。
- **自动排课口径**：单班 `timetable/generate` 与全校 `timetable/generate-school` 都是**纯计算**（前端预览后再保存），入参为 `{rules:{days:[…], subjects:[{name, weekly, double, session, max_per_day, forbid_periods}]}}`。硬约束：格子不重复、`forbid_periods` 禁排、`session` 限上/下午（`start_time` < 12:00 视为上午）、连堂不跨上/下午、`max_per_day` 上限；软偏好为「同科目分散到不同天」。重试上限：单班 **60** 次、全校 **80** 次；不可行时 `success=false` 并返回逐字告警（如「尝试多次均无法排出满足全部规则的课表，请减少节数、放宽连堂或时段限制」）。
- **全校排课与 `commit`**：`generate-school` 依据各班「任课设置」把科目→教师映射（某班未登记任课则不排该科目），并把**教师不可用时段**预先标记为占用；同一教师同一时段全校只允许一处（跨班零冲突，已验证）。`commit` **默认 false**（只验证可行性、返回 `classes:[{class_id,class_name,entry_count}]`，**不落库**）；`commit=true` 才在**单个事务**内逐班落库，并把这些班级的 `pending` 课表修改申请自动作废（note「管理员全校自动排课，本申请自动作废」）。
- **课表进阶口径**：`import-csv`（multipart `file` + `dry_run`，**默认仅预览**）按「年级 + 班级」定位班级，找不到的行报错跳过；某班出现在 CSV 即整体覆盖该班课表（先清后插）并自动为科目标色；`unavailabilities` 与 `teacher-assignments` 均为 **replace 语义**（`assignments` 为 `required|array`，缺失或空数组 → 422；同科目重复行去重）；`check-conflicts` 同时返回 `teacher`（跨班同格子且周次重叠：`all` 与任何重叠、`odd`/`even` 仅同类重叠）与 `unavailable`（命中教师不可用时段）两类冲突，冲突文案与 Laravel 逐字一致。
- **相对时间文案**：Laravel 未配置 locale（`config/app.php` 只覆盖 timezone），框架默认 `en`，故 `diffForHumans()` 输出英文；Go 端按同口径输出英文近似文案（`3 minutes ago` / `just now`）。若将来把 Laravel locale 改为中文，这里需同步替换。

### 管理端运维的**有意差异**（逐条）

以下都是本批移植时确认与 Laravel 不同、且有意为之的点（其余口径逐条对齐）：

- **学校 LOGO 不缩放/裁剪**：Laravel 用 intervention/image 处理后再存；Go 端**直接保存原文件**到 `UPLOAD_DIR/schools/`（默认 `storage/app/uploads/schools/`），并把 `/storage/app/uploads/schools/<文件名>` 写入 `schools.logo_path`。校验一致：jpeg/png/gif/webp、≤ 2MB、内容必须是图片（按文件头嗅探）。
- **`system/logs` 的日志来源**：Laravel 固定读 `storage/logs/laravel.log`；Go 端日志走 stdout、**没有日志文件**，因此改读环境变量 `LOG_FILE` 指定的文件——文件存在时返回其末 N 行（等级识别与 `?level=` 过滤、时间逆序与 Laravel 一致）；未配置或文件不存在时返回空列表 + 一条说明性 `message`，**不伪造日志**。
- **`system/status` 无迁移记录表**：Go 用 `AutoMigrate` 建表，没有 `migrations` 表 → `migrations` 恒为 `[]`、`migration_count` 恒为 0，改以 `tables`（真实表清单）承载信息；`version` 用 `go`/`backend` 替换 Laravel 的 `php`/`laravel`，并追加 `db_driver`/`timezone`。
- **`system/diagnose`**：`third_party_bindings 表` **已随 `AutoMigrate` 建表**（支撑 `GET /auth/bindings`），故该项状态由 `skipped` 变为 `ok`（自检按表是否存在动态判定）；`教师账号密码状态` 已改为**与 Laravel 同口径**——统计 `role=teacher` 且 `plain_password` 为 NULL/空串的账号数，>0 时状态为 `fixable`，detail 为「N 个教师账号缺明文密码（第三方自动注册历史问题），请在教师列表逐个重置」。
- **`system/repair`**：以 GORM `AutoMigrate` 代替 `artisan migrate --force`（幂等）；`pending_password_resets` **与提示文案同 Laravel**——统计缺 `plain_password` 的教师数，>0 时 message 追加「；检测到 N 个教师账号缺少明文密码记录（第三方自动注册历史问题），请在「教师管理 → 密码」中逐个重置为默认密码」。
- **明文密码（本批已从「不存」改为「存并受控暴露」）**：Go `users` 表现有 `plain_password` 列，创建教师（单条 + 批量 + CSV 导入）、重置密码（单条 + 批量）、教师自助改密、`API 机器人` 播种时都会同步写入明文；`GET /admin/teachers/:id/password` 无明文时自动生成 8 位并写回（同 Laravel）。差异：`models.User.PlainPassword` 带 `json:"-"`（**不随 User JSON 输出**），Laravel 的 `User` 未把该字段放进 `$hidden`，其 user JSON 会带上它——Go 端刻意不扩大泄漏面，明文只经 `GET /admin/teachers/:id/password` 单点返回。
- **批量删除账号的机器人守卫**：Laravel 的 `batchDeleteAccounts` **没有**机器人守卫（只有单删 `disableTeacher` 有）；Go 端按项目决策 13 统一保护 API 机器人账号（跳过并在 `data.protected_count` 中回报），同时在删除前解除 `class_rooms.teacher_id` 与 `class_room_teachers` 关联（Laravel 会留下指向已删用户的 `teacher_id` 脏引用）。
- **教师 CSV 导入**：仅接受 UTF-8（Laravel 会用 `mb_convert_encoding` 自动转 GBK）；**不会解析 xlsx/xls**（Laravel 走 PhpSpreadsheet，Go 无该依赖，直接 422 提示另存为 CSV）。表头/分隔符/别名口径（`姓名|年级团队|所属年级团队|科目|密码|手机号` 与英文键、tab/逗号/分号自动识别、姓名为空的行跳过）与 Laravel 一致。
- **学生导入 `students/import`**：Laravel 只接受 JSON 数组；Go 端在 JSON 之外**额外**支持 multipart CSV 上传（字段名 `file`，表头 `姓名/班级/性别/学号` 或英文字段名；无表头时按列序解析）。查重、跨班学号拦截、同名提醒（非阻塞 warning）、422 部分失败仍回传已成功行等语义逐条对齐；姓名入库前做 TrimSpace（Laravel 用原值）。
- **教师的昵称默认值**：Laravel 用拼音（`PinyinService`）并去重；Go 端无拼音库，昵称默认取**姓名本身**（等价于 Laravel 拼音缺失时的回退路径），`_2/_3` 去重规则不变。
- **`class_room_teachers` 已落库**：本批新增该关联表（`role`/`subject`），`assign-teacher`/`remove-teacher`/批量创建教师均按 Laravel 的 `updateOrCreate` 语义写入；但**教师可见班级仍只看 `class_rooms.teacher_id`**（协同教师不扩大权限范围，Laravel 的 `TeacherClassScope` 亦如此）。
- **学生批量删除是硬删除**：Laravel 的 `Student` 有 SoftDeletes（只打 `deleted_at`）；Go schema 的 `students.deleted_at` 不触发软删除，故与 Go 端既有单条删除一致：**硬删除并级联清理**宠物 / 积分 / 审计日志。

### AI 配置与对话（本批）

管理端 6 个（`admin/ai/settings` GET/PUT、`toggle`、`usage`、`fetch-models`、`test`）+ 教师端 4 个（`teacher/ai/config`、`chat`、`commands`、`usage`）+ 大屏端 2 个（`display/ai/settings`、`chat`，班级码 token）共 **12 条路由**；设置/会话落在 `ai_settings` / `ai_conversations`（`AutoMigrate` 建表）。
### AI 配置与对话（本批）的**有意差异**（逐条）

- **两张表 + 一套服务**：`ai_settings`（每校一行，`enabled/provider/api_key/api_base/model/max_tokens/tokens_used/tokens_limit/providers`）与 `ai_conversations`（`provider/question/answer/tokens_used/prompt_tokens/completion_tokens/cost/currency/status`）按 Laravel 迁移逐字段建表，走 `AutoMigrate`（Go 无迁移文件）。`providers` 与 Laravel 一样是 **JSON 文本列**，内部用 `map[string]any` 原样读写（未知字段不丢失）。
- **供应商调用**：`services.AIService` 逐条对齐 `App\Services\AiService`——`claude`/`google`/`qwen`/`mcp` 四种专属实现 + 其余 20 余种 OpenAI 兼容；`api_base` 覆盖默认地址真正生效（qwen 的对话地址 `dashscope…/api/v1` 与模型列表地址 `…/compatible-mode/v1` 不同，与 Laravel 一致）；HTTP 状态码 ≥400 返回固定兜底文案（`AI 服务暂时不可用` / `AI 服务不可用`），传输层异常向上抛出。**无新依赖**（`net/http` + `encoding/json`），HTTP 客户端可注入（`SetHTTPClient`），测试全部走 `httptest` 假供应商。
- **计费只在本地算**：`ai_conversations.cost` = `round(prompt/1e6*input_price_per_m + completion/1e6*output_price_per_m, 6)`（同 `AiBillingService::calculateCost`）。**官方账单/余额查询已在本批补齐**（`POST /admin/ai/provider-official` + 7 个驱动，详见文末「登录别名 / CSV 导出 / AI 官方账单」）。
- **tokens_limit 只展示不拦截**：Laravel 全仓只读该值做界面展示，`chat` 不做任何超限拒绝；Go 端照做（不新增拦截）。
- **api_key 语义**：`GET /admin/ai/settings` 的顶层**不返回** `api_key/api_base/model/provider`（Laravel 只回 `enabled/tokens_used/tokens_limit/max_tokens/providers`），但 `providers[].api_key` 原样返回（**不掩码/不截断**）；`PUT` 时同一供应商传空 `api_key` 保留原值、计数器（`tokens_used`/`total_calls`/`estimated_cost`/`balance`）取新旧较大值、`currency` 缺失时沿用旧值；**请求未带 `providers` 字段时按 `input('providers', [])` 语义把 providers 置空**（Laravel 行为，已在测试中固化）。
- **无效请求的校验文案**：状态码与 Laravel 一致（`参数错误` 422 / `缺少供应商` 422 / `该供应商未配置 API Key` 422 / `请输入问题` 422 / `AI 功能未开启` 403 / `班级不存在` 404），`saveAiSettings` 的 422 额外带 `errors` 字段（字段键用 Laravel 的 `providers.0.id` 形式），但错误明细文案统一为中文（同其余批次约定，Laravel 为英文默认文案）。
- **空设置行不再 500**：Laravel 在 `fetch-models` / `test` / 教师端 `usage` / `config` 读取 `null` 属性时会触发 ErrorException（500）；Go 端按业务语义返回 422「该供应商未配置 API Key」或空态（`configured=false`），不伪造 500。
- **大屏端学生信息恒为「匿名」**：Laravel 的 `disp_`/`class_` token 载荷里没有学生信息（`validateToken` 只给 `class_id`），`student_id` 恒为 `NULL`、`student_name` 恒为 `匿名`；Go 端同。大屏端不校验 `question == "0"` 之外的边界（PHP `empty("0")` 为真，Go 端同样按 422 处理）。
- **display 未授权文案走既有中间件**：`display/ai/*` 无 token 时返回 `DisplayAuth` 的 401 `Token 无效或已过期`（Laravel 该动作自己校验并回 `无效的班级码`；与 Go 端其余 display 接口口径一致）。
- **响应信封**：统一 `{data, message:"ok"}`（同其余已移植接口）；`toggle` 为逐字对齐的 `{message:"AI 功能已开启|已关闭", data:{enabled}}`，`saveAiSettings` 为 `{data:null, message:"AI 设置已保存"}`。
- **测试**：`internal/services/ai_test.go`（设置读写/掩码与空值语义/校验、toggle、fetch-models 三种供应商、test 成功失败与传输层异常、教师端 4 种供应商请求体与 token/cost/落库/累加、上游错误与未配置边界、白名单与 model_map、教师 config/usage/commands、大屏开关检查与对话校验顺序、管理端用量分组）+ `internal/services/ai_billing_test.go`（7 个官方账单驱动、守卫文案、上游 500/空数组/不可解析/传输异常兜底）+ `internal/router/ai_routes_test.go`（AI 路由注册、display 鉴权与守卫）。所有上游调用指向 `httptest`，不访问外网。
- **学年升级路径**：与 Laravel 一致为 `/admin/grade-upgrade/preview|execute`（位于 `students` 路由分组之外）。升级语义逐条对齐且**不可逆**（六年级学生标记毕业 + 班级归档；一~五年级年级与班名升级；重复调用会继续升级）。
### 登录别名 / CSV 导出 / AI 官方账单（本批）的有意差异

- **导出用 CSV 代替 xlsx**（零新依赖）：Laravel 走 `maatwebsite/excel` 输出 `.xlsx`；Go 端用标准库 `encoding/csv`。文件名保留 Laravel 主体、后缀改 `.csv`（`<班级名>-<Ymd-His>-积分报表|宠物报表|考勤报表.csv`、`<班级名>-课表.csv`、`全校课表.csv`），统一 `Content-Type: text/csv; charset=UTF-8` + `Content-Disposition: attachment; filename="<percentEncode>"`（同既有 `timetable/export-cses`），并写 UTF-8 BOM（Excel 打开中文不乱码）。xlsx 的**样式与工作表名**丢弃；**全校课表**由「每班一个工作表」改为单 CSV 顺序拼接（班级之间空行分隔）。
- **报表列与行语义逐字照抄**：`ScoresExport`（姓名/学号/总积分/获得积分/扣除积分/宠物名/宠物等级；`negative_score` 取绝对值；无宠物 → 宠物名空串 + 等级 0）、`PetsExport`（无宠物 → 宠物名「无」、阶段「未孵化」）、`AttendanceExport`（状态中文映射，未知状态原样输出；**「签到时间」列恒为空串**——Laravel 读的是并不存在的属性 `check_in_time`，Go 端照抄，不改为 `sign_in_at`；`pluck('id','name')` 的**同名去重**语义同样保留）。
- **考勤按日期过滤**：Laravel `whereDate('created_at', $date)` 走 DB 端日期函数；Go 端改为**业务时区（Asia/Shanghai）当日 `[00:00, 次日 00:00)` 区间**（纯 Go sqlite 驱动写入的是带时区偏移的时间串，`date()` 会先换算成 UTC 导致跨日错位）；非法日期不匹配任何行。查询补 `Order("id ASC")`（Laravel 无显式排序，隐式插入顺序等价）。
- **「不支持的类型」分支**：Laravel 返回 200 + `{message:"导出类型 X 不支持，可选: scores, pets, attendance"}`（无 `data` 字段）；Go 端沿用仓库统一信封，额外带 `data: null`。
- **登录别名的校验失败文案**：Laravel 为 422 + `errors`；Go 端为 422 +「请求参数格式错误」。速率限制本批已补齐（`throttle:6,1`，见文末「宠物切换图鉴 / 考勤 `leave_record` / 登录限流」）。
- **AI 官方账单**：7 个驱动（OpenAI usage、DeepSeek 余额、Moonshot 余额、SiliconFlow 余额、OpenRouter key、New API/One API 中转 billing、本地精确计费回退）逐条对齐取值路径、默认值与失败兜底。`official=true` 仅表示「用量或余额至少有一项非 null」，**本地回退时 `usage.source` 仍为 `local`**（Laravel 原样行为）。DeepSeek 余额端点与 OpenRouter key 端点在 Laravel 中为硬编码 URL（忽略 `api_base`），Go 端生产同样硬编码，仅保留测试用 `SetEndpointOverrides`（httptest 假上游，不访问外网）；`syncOfficialUsage` 的「`supportsUsage()` 为真但首次返回 null 时会再查一次」原样保留。查到余额时写回 `providers[].balance`（「取较大值合并」仍由保存接口负责）。`providers[].balance` 写回在 Go 端为整体重序列化 JSON 文本列（语义等价）。

### 认证会话 / 管理端人员 / 教师端学生与图鉴（本批）的有意差异（逐条）

本批 15 条路由 + 相关服务函数的差异点（其余口径逐条对齐）：

- **JWT 撤销机制（登出/刷新真的失效）**：Laravel 用 Sanctum 的 `personal_access_tokens`——登出/刷新即删除令牌行，旧 token 立即 401。Go 端是自签 JWT（无服务端会话），因此给令牌加标准 `jti` 声明（`RegisteredClaims.ID`，32 位 hex），新增表 **`revoked_tokens(jti 主键, expires_at)`**：`auth/logout` 与 `auth/refresh` 把**旧令牌的 jti** 写入该表，`middleware.Auth` 命中未过期的 jti → 401「登录已过期或无效」（复用既有文案）。**无 `jti` 的历史令牌视为有效**（向后兼容，不做强制下线）。撤销记录惰性清理（每次撤销顺带删除过期行，无定时任务）。已验证：登出后旧 token 401、另一 token 不受影响、刷新后旧 token 401 而新 token 可用（`internal/router/session_routes_test.go`）。
- **响应信封差异**：Laravel `logout` 返回 `{"message":"已登出"}`、`refresh` 返回 `{"data":{"token":...}}`、`bindings` 返回 `{"data":[...]}`；Go 端沿用仓库统一信封 —— `logout` 为 `{data:null, message:"已登出"}`、`refresh` 为 `{data:{token}, message:"ok"}`、`bindings` 为 `{data:[...], message:"ok"}`。`GET /admin/students` 例外，与 Laravel 一致地返回**裸** `{data, meta}`（无 `message`）。
- **`bindings` 的平台集合与字段**：字段与顺序逐字照抄（`platform/label/icon/bound/nick`），`platforms()/platformLabels()/platformIcons()` 三张映射表也逐字搬（含 `renren/dingtalk/feishu`）；学校开关 `schools.settings.enabled_third_party_platforms` 为空数组或未配置时回退**默认三平台** `[wechat_work, wechat, qq]`，非空则与白名单取交集并保持白名单顺序（交为空 → 返回空列表，**不回落默认**，同 Laravel `array_intersect` 行为）。`nick` 为 `null` 时输出 null（Go 用 `*string`）。
- **`plain_password`（明文密码）**：`users` 表新增 `plain_password` 列，写入点为创建教师（`POST /admin/teachers`、`batch-create`、CSV 导入）、重置密码（`POST /admin/teachers/:id/reset-password`、`accounts/batch-reset-password`）、教师自助改密、API 机器人播种。Laravel 只在「创建教师 / 单条 reset / changePassword / BotTeacherSeeder」写明文，**`batchResetPassword` 不写**；Go 端两条重置路径都写（更贴近前端「查看/重置后能直接看到密码」的用法）。`GET /admin/teachers/:id/password` 无明文时生成 8 位随机密码并写回（同时重置 `password_hash` 与 `password_changed=false`，同 Laravel）。**唯一收紧**：`models.User.PlainPassword` 带 `json:"-"`，不在任何 User JSON 里出现（Laravel 的 user JSON 会带），明文只经该单点接口返回。
- **`GET /admin/students` 的 DTO**：Laravel 直接 `$s->toArray()`，故其元素含 `deleted_at` 与嵌套的 `class_room`/`pet` 关联对象；Go 端输出**显式字段**（学生基础字段 + 追加的 `class_name/class_grade/pet_species/pet_level/pet_name`），空串的 `student_no`/`avatar_path` 序列化为 `null`（对齐 Laravel 可空列）。筛选/排序/分页口径一致（`status` 默认 `active`、`all` 不过滤，`per_page` 默认 50，`id desc`，`meta` 四项同 LengthAwarePaginator；`total=0` 时 `last_page=1`）。
- **`POST /admin/students` 的校验分层**：`exists:class_rooms,id` 不通过（班级整体不存在）→ 422「参数错误」+ `errors.class_id`；班级存在但非本校 → 404「班级不存在」（`findOrFail`）。Go 端新增 `services.ValidationError`（带字段错误的 422）来承载这一分层，`fail()` 统一渲染成 Laravel 风格 `{message, errors}`。
- **`GET /admin/classes/:id`**：返回班级 + `teacher`（无班主任为 `null`）+ `students`（该班全部学生，含已停用，同 Eloquent 关联语义）；班级字段沿用 Go `ClassRoom` 模型，`settings` 仍是 `json:"-"`（Laravel 会带 `settings`），需要班级设置请用 `teacher/class` 或大屏 `class-settings`。
- **`PUT /admin/teachers/:id/classes` 的 replace 语义**：逐条 upsert `class_room_teachers`，`head_teacher` 同步 `class_rooms.teacher_id`；随后删除「本次未提交」的旧分配，被删且原 role 为 `head_teacher` 的班级把 `teacher_id` 置空。`assignments` 缺省/空数组 = **清空该教师全部分配**（同 Laravel `whereNotIn(..., [0])`）。`subject` 未传 → `null`。校验分层同 Laravel：`class_id` 指向不存在的班级 → 422 + `errors.assignments.N.class_id`；班级存在但非本校、或教师非本校 → 404。整批在**单个事务**内完成（Laravel 未包事务；失败时 Go 端不会留下半套分配）。
- **教师可见班级范围**：沿用 Go 端既有 `Scope.ClassIDs`（本校 + `class_rooms.teacher_id = 本人`，API 机器人放行全校；`status=active`）。Laravel `TeacherClassScope` 还会并入 `class_room_teachers` 里的副班/科任关联——**Go 端仍不认协同教师**（沿袭上一批的既定差异），因此本批的教师端学生/图鉴/按规则加分接口同样只看班主任班级。
- **教师端学生删除是硬删除**：Laravel `Student` 有 SoftDeletes（只打 `deleted_at`）；Go schema 的 `students.deleted_at` 不触发软删除，故与既有 `DELETE /admin/students/:id`、`students/batch-delete` 保持一致：**硬删除并级联清理**宠物 / 积分 / 审计日志。其余口径（同班学号/同名查重、跨班同学号拦截文案、导入不分配宠物、`create` 性别归一化、`import` 性别仅空值补「未知」）逐条对齐。
- **404 文案**：Laravel 的 `findOrFail` 统一返回「资源不存在」；Go 端沿用既有更具体的中文文案（`学生不存在` / `班级不存在` / `教师不存在` / `规则不存在或不可见`），**状态码一致**。
- **宠物图鉴表（本批已补齐「切换即归档」）**：`pet_collections` 表（`student_id/species/level/experience/mood/is_active`）支撑 `GET /teacher/pets/:studentId/collection`，读取时把当前激活宠物**补录**进图鉴（`firstOrCreate` 语义）；本批起教师端与教室端的换宠都会写图鉴（见上文「宠物系列与切换」），图鉴内容因此来自真实切换历史 + 补录两条来源。`unlock_slots = 1 + 总积分/100`（同 `PetCollection::unlockSlotsForScore`）。
- **`GET /teacher/class` 的 `settings`**：空设置输出 `null`（对齐 Laravel 的 `$class->settings`），非空输出解析后的对象；`class_points` 取 `settings.class_points ?? 0`；`display_code` 复用确定性班级码生成（同 Laravel `DisplayCodeService::generate`）。
- **`POST /teacher/scores/give-by-rule/:ruleId`**：校验顺序同 Laravel（**先取可见规则 404，再校验 `student_id` 422**，最后取管辖内学生 404）；落账复用既有 `ScoreService.GiveScore`，因此事务、`score_logs` 审计、宠物经验/等级校正、`score_update` 大屏事件、余额钳制口径与 `POST /teacher/scores/give` 完全一致（Laravel 亦如此：`giveScoreByRule` 内部调 `giveScore`）。
- **速率限制**：Laravel 的这些路由在 `auth:sanctum` 组内，无单独 throttle；Go 端同样不加。

### 第三方登录与通讯录（本批）的有意差异（逐条）

本批 14 条路由 + 相关服务函数的差异点（其余口径逐条对齐）：

- **扫码上下文由 Cache 改为表**：Laravel 用 `Cache::put('wechat_scan_ctx:<uuid>', $ctx, 10min)`；Go 端新增表 **`temp_binding_contexts(temp_token 主键, context JSON, expires_at)`**——读取时校验未过期（过期视为不存在并顺带删除），`bind-after-scan` 成功后删除记录（一次性消费），写入路径惰性清理过期行（无定时任务）。UUID 用 `crypto/rand` 生成 v4（**零新依赖**，未引入 `google/uuid`）。已验证：落库含 `platform/platform_id`、`expires_at ≈ now+10min`、过期后绑定改走请求体参数、消费后行数 0。
- **JWT 替代 Sanctum**：Laravel `createToken(...)->plainTextToken`；Go 端签发自签 JWT（`data.token` 语义一致，形态不同）。**已绑定分支不带 token** 的 Laravel 原样行为在 wechat / qq / renren / wechat-work / third-party-login 五处均照抄（详见各路由表格），Go 端用 `omitempty` 表达「无 token 键」。
- **企微 access_token 缓存落表**：Laravel `Cache::put("wecom_at:<schoolId>", $token, max(expires_in-300, 60))`；Go 端改用表 **`wechat_work_tokens(school_id 主键, token, expires_at)`**（多实例共享，避免各实例重复 gettoken 触发限频；**不用进程内全局缓存**）。原样保留：`corp_id/secret` 缺失 → 「企微未配置」；响应无 `access_token`（含 null）→ 「token失败」；`errcode != 0` 但带 token 时 Laravel 仍返回该 token（只判 `!isset`），Go 端同。
- **`WECHAT_WORK_API_BASE` 注入点（仅测试 / 自建代理）**：Laravel 把 `https://qyapi.weixin.qq.com` 硬编码；Go 端默认值即官方地址，允许环境变量覆盖——**生产恒为官方地址**，不配置即官方。钉钉 / 飞书的端点（`api.dingtalk.com` / `oapi.dingtalk.com` / `accounts.feishu.cn` / `open.feishu.cn`）在 Laravel 同样是硬编码，Go 端保留官方默认值并把 `ThirdPartyManager` 的四个端点字段作为**测试用注入点**（同 AI 官方账单批次保留 `SetEndpointOverrides` 的做法）。所有 provider 的 HTTP 客户端可注入（`SetHTTPClient`），测试全部指向 `httptest` 假上游，**不访问外网**。
- **应用凭证走环境变量**：键名与 Laravel 配置一致——`WECHAT_WORK_CORPID/AGENTID/SECRET/TOKEN/ENCODING_AES_KEY`（`TOKEN` / `ENCODING_AES_KEY` 供 `GET|POST /wechat-work/callback` 验签解密，最后一批已实现）、`DINGTALK_APP_KEY/SECRET`、`FEISHU_APP_ID/SECRET`；服务构造时读取（同 `AdminOps` 读 `UPLOAD_DIR` 的做法），不占用 `config.Config`。
- **`redirect_uri` 默认值**：Laravel `url('/login')`（= `APP_URL + /login`）；Go 端无 `url()` 助手，读环境变量 `APP_URL`（默认 `http://localhost`，同 `/admin/system/status` 的 `app_url` 口径）。参数存在即按原值使用（含空串），缺省才回落默认（同 `$request->input($key, $default)`）。
- **`school_id` 不存在的语义**：`GET /auth/third-party/auth-url` 在 `school_id > 0` 时**只按该 ID 查找**，找不到即 400「系统尚未初始化」，**不回退第一所学校**（同 Laravel `$schoolId > 0 ? School::find($schoolId) : School::first()`）；`state` 解析（`resolveSchoolFromState`）则按 Laravel 回退第一所学校，`(int)` 转换用 PHP 语义（前导数字生效，如 `"2abc"` → 2）。
- **无拼音库**：`uniqueNickname` 在 Laravel 是「姓名拼音」，Go 端沿用本仓约定——昵称默认取**姓名本身**（等价于 Laravel 拼音库缺失时的回退路径），`_2/_3` 去重不变；`bind-after-scan` 的「本地昵称仍是默认值才覆盖」判断因此与 `user.name` 比较（Laravel 比 `PinyinService::toPinyin($user->name)`）。
- **建号去重的范围**：Laravel `uniqueUsername` 是**全校唯一**；Go 的 `users.username` 是全局唯一索引，故按全局判重（与 `AdminOps.CreateTeacherAccounts` 一致，避免落库撞唯一索引）。`nickname` 仍按校内去重。
- **第三方登录复用本地账号**：`loginWithThirdParty` 无绑定时先按**手机号**、再按**实名用户名**匹配本地教师（同 Laravel），命中则补建绑定、本地缺手机号 / 头像时顺手同步，并返回 token。
- **通讯录导入的校验与错误码**：企微 / 第三方两条 import 的校验规则与 Laravel 一致（`teachers.*.name` 必填 ≤50、`mobile` ≤30、`email` ≤120、`students.*.name` ≤50、`class_id` 必填、`gender` ≤10），但错误明细文案为中文，键名保持 Laravel 的 `teachers.0.name` 点号形式；校验失败 422「参数错误」+ `errors`。`contacts` 类接口的上游异常统一包成 400 + 原始中文文案（同 Laravel `catch (Throwable)`）。
- **学校记录缺失**：Laravel 的 `wechatWorkContacts` 直接读 `$request->user()->school->id`（null 会 500），import 显式 404「未找到学校」；Go 端两条路径统一 404「未找到学校」，**不伪造 500**。
- **可空列落空串**：`phone/email/gender` 等可空字段 Go 端落空串（本仓既有约定）而非 NULL，查重与展示口径不受影响；`students.*.class_id` 允许指向不存在的班级（Laravel 亦只校验 `integer`，Go 端照抄，不额外加 `exists`）。
- **响应信封**：统一 `{data, message:"ok"}`（同其余已移植接口）；两条 import 例外，与 Laravel 一致返回 `{message, data}`（不带 `message:"ok"`）。
- **速率限制**：Laravel 的 `auth/teacher/login`、`auth/admin/login`（`throttle:6,1`）与 `auth/class/login`、`display/login`（`throttle:10,1`）在 Laravel 里带 throttle，**本批已在 Go 端补齐**（进程内固定窗口，见文末「宠物切换图鉴 / 考勤 `leave_record` / 登录限流」）；本批 8 条公开第三方路由在 Laravel 里**没有** throttle，Go 端同样不加。
- **测试**：`internal/services/third_party_test.go`（企微 token 成功/失败与缓存命中、code 换 userid、部门+成员去重与 `department_names`、钉钉/飞书 `getUserByCode`/`fetchContacts` 与 authUrl 逐字断言、options 三形态、auth-url、微信/QQ/人人通分支、免注册建号、第三方登录四分支、bind-after-scan 覆盖条件、bind/unbind、导入统计与去重）+ `internal/router/third_party_routes_test.go`（14 条路由注册、公开无需 token、bind/unbind 401、admin 403、端到端行为、路由总数）+ `internal/services/{school_settings_test.go,delete_teacher_test.go}`（学校 settings 替换语义、`DELETE /admin/teachers/:id` 的解除关联/机器人守卫/跨校 404）。全部使用内存 SQLite + `httptest`，**不访问外网**。

### 第三方批之后顺带补齐的三处（学校设置 / 方法补齐 / 课表文案）

- **学校 `settings` 真正可读写（此前只能读代码里的内部键）**：`GET /admin/school` 现在把 `settings` 作为**对象**返回（此前 `models.School.Settings` 带 `json:"-"`，前端「学校设置」页读不到 `third_party_platform` / `enabled_third_party_platforms`，保存的开关也永远不会生效）；`PUT|POST /admin/school` 接受 `settings`（Laravel `nullable|array`）并**整体替换**该列（`$school->fill($request->only([...,'settings']))` 的语义，不做按键合并），非法类型 → 422「参数错误」+ `errors.settings`，成功文案改为 Laravel 原文「学校信息已更新」。新增 `services.SchoolView` / `SchoolViewOf` 供管理端统一输出设置对象。
- **补两处 Laravel 有、Go 漏注册的方法**：`POST /admin/school`（Laravel `match(['put','post'],'school')`）与 `POST /admin/classes/:id/timetable`（Laravel `match(['put','post'],'{classId}/timetable')`，前端班级课表页正是用 POST）——两条路径此前只有 `PUT`，前端调用会 405。`DELETE /admin/teachers/:id`（Laravel `disableTeacher`）也一并补齐：解除 `class_rooms.teacher_id` 与 `class_room_teachers` 关联后**硬删除**，机器人账号 403「API 机器人账号不可删除。如需停用，请在 .env 中设置 BOT_ENABLED=false 后重启」，成功「教师账号已删除」（`internal/services/admin_ops_accounts.go::DeleteTeacher`）。
- **课表三处文案改为 Laravel 逐字**：`PUT|POST /admin/classes/:id/timetable` → 「课表已保存并即时生效」（原为「课表已保存」）；教师端提交申请 → 「修改申请已提交，待管理员审核」（原为「课表修改申请已提交，等待管理员审核」）；审批 → 「已通过并应用课表」/「已驳回」（原为「已通过该课表修改申请」/「已驳回该课表修改申请」）。

### 积分分类字典 / 企微回调与请假同步（最后一批，3 条路由）

`GET /api/v1/common/score-categories`（公开）与 `GET|POST /api/v1/wechat-work/callback`（公开）；服务层新增
`services.ScoreCategories`（分类字典）与 `services.WechatWorkAttendanceService`（企微审批 → 请假 → 考勤），
`services.WechatWorkService` 补上 `GetLeaveApprovalSpNos` / `GetApprovalDetail` / `parse` 与 `postJSON`。
配套 CLI：`./learnstar-go -sync-wechat-work-leave [-school-id=N] [-date=YYYY-MM-DD]`（打印 `{synced, skipped, failed}` 后退出，不启动 HTTP 服务）——
Laravel 的等价物是 `php artisan attendance:sync-wechat-leave {--date=} {--school=}`。

**有意差异（逐条）**：

- **XML 解析用 `encoding/xml` 而非 SimpleXML**：定义最小结构体（不带根元素名，语义同 PHP `$xml->Encrypt` 的「根元素直接子元素」），`<![CDATA[…]]>` 由 encoding/xml 自动并入文本（等价 `LIBXML_NOCDATA`）；`SpStatus` 用指针区分「元素缺失」（PHP `?? 1` → 1）与「元素为空」（`(int) ''` → 0）。**唯一放宽**：`sys_approval_change` 的判定除顶层 `Event` / `ChangeType`（Laravel 只判这两个）外，还兜底判 `ApprovalInfo` 内层的 `ChangeType`。
- **PKCS7 剥离沿用 Laravel 的 `strip()` 口径**：末字节 1..32 时按该长度截尾（`pad >= len` 时返回空串，同 PHP `substr`）；末字节为 0 时**原样返回**（零填充场景）。AES 解密为 `key = base64(aes_key + "=")`（长度 %4 == 1 判非法，同 PHP 严格模式；%4 == 2/3 自动补 `=`，PHP 亦如此）、IV = `key[:16]`、AES-256-CBC、**不做填充校验**且要求密文长度为 16 的倍数（同 `OPENSSL_RAW_DATA | OPENSSL_ZERO_PADDING`）。密文用**严格** base64 解码（PHP 用非严格模式，忽略非法字符）——仅在报文被篡改时有差别，两边都判为「解密失败」。
- **GET 失败一律返回空字符串、POST 一律返回 ok**：与 Laravel 相同（`Content-Type` 为 `text/plain`；POST 为 JSON `{"errcode":0,"errmsg":"ok"}`）。
- **日志走 stdout**：Laravel 记 `Log::error('企微验证失败'/'企微回调异常')` 到 `storage/logs/laravel.log`；Go 端用标准库 `log`（同仓库既有约定，`system/logs` 接口需自行把日志重定向到文件才能查看）。
- **`matchStudent` 按规格语义实现三段优先级**（精确同名 → 姓名包含 → 该家长名下第一条）。Laravel 原实现里 `$q` 是同一个 Builder，`$q->where('name', $hint)` 会**留在**后续 `LIKE` 查询上，使 LIKE 与「第一条」两段实际不可达（等价于「只认精确匹配」）——Go 端按业务意图实现，是**有意行为放宽**。
- **`students.parent_id` 的现状**：Laravel 已于迁移 `2026_08_06_000006` 移除家长角色并清空该列，Laravel 与 Go 的匹配路径都只在「家长数据被外部写入」时才会命中；Go 端为对齐该语义新增了 `students.parent_id` 列（`AutoMigrate` 建列）。
- **日期列存字符串**：`wechat_work_leave_records.leave_start_date/leave_end_date` 按本仓既有约定存 `YYYY-MM-DD` 字符串（Laravel 是 `date` 列 + Carbon cast），因此 `getUnsyncedForDate` 用字符串比较；`synced_at IS NULL` 的过滤口径不变。
- **`syncForSchool` 的日期区间**是「当天 00:00:00 / 23:59:59」的业务时区 Unix 秒（Laravel `strtotime($date.' 00:00:00')`）；日期字符串非法时 Go 端直接报错（PHP `strtotime` 返回 false 后会落到 `date('Y-m-d', 0)` 的怪结果）。
- **`syncForSchool` / `handleWebhookCallback` 只要匹配到学生就补考勤**（Laravel 同样不校验 `approve_status`，`pending` 记录也会落考勤）；`syncAll` 遍历 `status=active` 学校、单校异常计 `failed` 且继续。
- **CLI 取代 artisan**：`-sync-wechat-work-leave` 在 `main.go` 里分支实现（命中时不启动 HTTP 服务），输出为 JSON 单行 `{"synced":N,"skipped":N,"failed":N}`；`-school-id=N` 单校同步失败时以非 0 退出（对应 artisan 的异常路径）。
- **`wechat_leave_count` 不再是常量 0**：`AttendanceService.Start` / `Summary` 现在统计当日 `source = wechat_work` 的考勤记录（同 Laravel），点名（`startForClass`）末尾会接上 Laravel `startAttendanceForClass` 的后半段：把当日仍未同步的企微请假补写到考勤。
- **测试**：`internal/services/wechat_work_test.go`（getapprovalinfo 分页/errcode 中断、getapprovaldetail 的 nil 与 parse 字段映射）+ `internal/services/wechat_work_attendance_test.go`（匹配优先级、applyLeave 三分支与 `synced_at` 回写、回调四个早退分支与成功、syncForSchool 未配置/跳过/日期区间、SyncAll 汇总与单校失败、点名补同步）+ `internal/router/wechat_work_routes_test.go`（3 条路由注册与公开性、分类字典字段、**标准库自造签名与 AES 密文的加密链路**、GET 四种结果、POST 事件解析与四类边界）。全部使用内存 SQLite + `httptest` 假上游，**不访问外网**。


### 宠物切换图鉴 / 考勤 `leave_record` / 登录限流（本批）的有意差异（逐条）

本批收敛了 4 项此前记录在案的「Laravel 有、Go 缺」的产品可见行为（该批路由总数 216 条；随后清理别名路由后为 211，见「别名路由清理」）：

1. **教师端换宠（`POST /teacher/pets/:student_id/switch`）**：逐条对齐 `PetService::switchPet` —— 同物种 422 → 类别限制 422 → 免费自选（`pet_free_picks`，用掉即删）→ 按等级扣分（不足 400）→ 旧物种进度归档（`is_active=false`）→ **恢复目标物种进度**（无记录则 1/0/80）→ 目标物种标记激活 → 写 `pets.last_switched_at`。响应为顶层 `{message, data}`，`cost`/`free_pick_used` 仅在「原有宠物」分支出现。
   - **顺带修正**：入参由 Go 端自造的 `species`/`name` 改回 Laravel/前端的 `pet_species`（必填 ≤50）+ `pet_name`（可空，缺省回退 `pet_species`）——此前前端 `PetCollection.vue` 的 `{pet_species}` 调用在 Go 端恒 422；响应也不再套 Go 的统一信封。
   - **同时按 Laravel 补齐了教师端的类别限制**：Go 端此前教师端换宠**完全不校验** `pet_series`（只有教室端校验）。文案照抄 Laravel——**教师端回显原始系列 id**（`只能领养当前类别「pokemon」的宠物…`），教室端回显 `seriesLabel` 中文标签（`…「宝可梦」…`），两端本就不同字。
   - **不再按积分重算等级**：原实现切换后调 `SyncLevelWithScore(total_score)`，Laravel 没有这一步（等级取图鉴或 1），已删除。
   - **新增列**：`pets.last_switched_at`（`AutoMigrate` 建列，Laravel 迁移 `2026_08_06_000005` 同名同义）。
2. **教室端换宠（`POST /display/pets/switch`）**：对齐 `DisplayController::classroomSwitchPet` —— 归档旧物种（`is_active=false`）→ 只换 `pets.species`（**保留**等级/经验/心情，不做恢复）→ 目标物种写图鉴并标 `is_active=true`（用当前 pet 的值）。无宠物分支照抄 Laravel：只建 `pets`、**不写图鉴**，文案「🎉 新宠物已诞生！」。
3. **考勤今日列表（`GET /teacher/attendance/today`）**：每条记录新增 `leave_record`（`{sp_no, leave_type, reason}`，无 `leave_record_id` 或关联行已删时为 `null`），取自 `wechat_work_leave_records`（同 Laravel 的 `with('leaveRecord:id,sp_no,leave_type,reason')`，只暴露三个字段）。**同时对齐 `check_in_time` 的形态**：原 Go 端直出 `time.Time`（RFC3339），现改为 Laravel `Carbon::toDateTimeString()` 的 `Y-m-d H:i:s` 字符串（业务时区）或 `null`——前端 `types/index.ts` 本就按字符串声明。
4. **登录端点限流**（`internal/middleware/throttle.go`）：`Throttle(max, window)` 按「**客户端 IP + 路由路径**」做**固定窗口**计数，仅挂在 Laravel 带 `throttle` 的 4 条路由上（`auth/teacher/login` 6/分钟、`auth/admin/login` 6/分钟、`auth/class/login` 10/分钟、`display/login` 10/分钟）；超限返回 **429 + `Retry-After: <剩余秒>` + `{"message":"Too Many Attempts."}`**（Laravel 框架默认响应体）。**零新依赖**（`sync` + `time`）。
   - ⚠️ **进程内计数，仅单实例有效**（与既有「大屏同码 5 次失败锁定 15 分钟」同一约定）；多副本部署需换共享存储。
   - ⚠️ **固定窗口 vs Laravel 的滑动窗口**：Laravel `RateLimiter` 逐次记录命中时间，窗口随最后一次命中顺延；Go 端首次命中开窗、到点整体重置。极端情形（第 59 秒打满 N 次）Go 端会在第 60 秒立刻放行，Laravel 需等最后一次命中 + 60 秒。
   - **客户端 IP 口径**：`c.ClientIP()`（与 `display/login` 的登录日志同一取值），沿用 gin 默认的可信代理配置——**默认信任全部代理**，因此会采用 `X-Forwarded-For` / `X-Real-IP` 的首段（反向代理后即真实客户端 IP；直连时可被伪造）。与 Laravel 受 `TrustProxies` 配置影响同理，未额外收窄。
   - 不附带 Laravel `ThrottleRequests` 的 `X-RateLimit-*` 响应头（前端未使用）；过期键惰性清理（每过一个窗口扫描一次），无定时任务。
   - `display/login` 既有的「同班级码连续失败 5 次锁 15 分钟」机制**保持独立并存**（Laravel 亦然），两者都会返回 429。

**测试**：`internal/services/pet_switch_collection_test.go`（教师端 6 例：无宠物建宠物+图鉴、恢复目标进度、未知物种初始化、同物种 422、跨类别原始 id 文案、免费自选；教室端 4 例：归档+保留等级、跨类别中文标签、无宠物不写图鉴、积分不足/免费自选）、`internal/services/pet_test.go`（既有两例改造为图鉴断言）、`internal/services/attendance_leave_record_test.go`（有/无/悬空 `leave_record_id` + `check_in_time` 形态 + JSON 字段）、`internal/middleware/throttle_test.go`（第 N+1 次 429 与 `Retry-After`、窗口过期恢复、按 IP 隔离、未挂路由不受影响、同实例多路由键含路径、并发放行数恰为 max）、`internal/router/throttle_routes_test.go`（4 条路由的路由级限流 + 未挂路由 + 路由总数 211）、`internal/router/pet_switch_route_test.go`（教师端换宠 HTTP 契约：入参/出参/图鉴副作用/422 分支）。全部使用内存 SQLite + `httptest`，**不访问外网**。

### 尚未移植清单（对照 Laravel `routes/api.php` 与前端实际调用）

统计口径（**方法级**）：把 `frontend-vue/src/**` 中出现的 `/api/v1/...` 字面量路径（模板变量归一为 `{p}`、去掉 query、去重）连同**推断出的 HTTP 方法**（`apiGet/apiPost/apiPut/apiDelete` 或 `fetch({method})`）与 Go 路由表逐条比对。路由总数 196 → 216 → **211**（第三方批 14 条 + 3 条方法补齐 + 最后一批的 `common/score-categories` 与 `wechat-work/callback` GET|POST，随后删除 4 条别名路由；`engine.Routes()` 实测，见 `internal/router/third_party_routes_test.go::TestThirdPartyRouteCount`）。

**结论：Laravel 有、Go 无的接口已清零；前端在用的接口也已无 Go 侧缺口，且三条前端侧不匹配已在前端修好（本轮）。**

前端侧修正（`frontend-vue/`，本轮；`npm run typecheck` 与 `npm run build` 均通过）：

- **积分规则路径（自查纠错）**：Laravel 的真实路径是 `/teacher/scores/rules*`（`routes/api.php` 第 185-188 行，`rules` 位于 `scores` 前缀组内），Vue 前端原版与 `mini-program/pages/scores/scores.js` 也都按此调用；Go 侧此前误注册为 `/teacher/score-rules*`，已改回 `scores/rules`。
- **按规则加减分**：`ScoresPage.vue` 由 `/teacher/scores/by-rule/{id}` 改为 Laravel 的真实路径 **`/teacher/scores/give-by-rule/{ruleId}`**（原路径在 Laravel 上也是 404）。
- **宠物接口**：`services/api.ts` 的 `renamePet` 由 `PUT .../pets/{petId}/rename` 改为 **`POST .../pets/{studentId}/rename`**（Laravel/Go 只有 POST，且路径参数是**学生 ID**）；同文件的 `feedPet` 修掉了「单引号字符串里写 `${studentId}` 导致未插值」的缺陷，`getPetDetail` 参数同样按学生 ID 命名。
- **删除死端点调用**：`services/api.ts` 的 `getRecentNews()`（`/teacher/dashboard/news`）在 Laravel 与 Go 都不存在，已删除（无调用方；`DashboardPage.vue` 的动态列表来自 `/teacher/dashboard` 的 `recent_news`）。
- **排行榜**：`LeaderboardPage.vue` 的 UI 键 `'pet'` 现在映射为后端片段 **`pet-level`**（此前直接拼接会打到 `/leaderboard/pet` → 404）。
- **班级宠物系列**：`ClassesPage.vue` 读 `settings.pet_series`、写 `PUT /admin/classes/:id {pet_series}`——这两条**后端此前不生效**，已在 Go 端补齐（见下）。

Go 端随之补齐的两处（本轮，均为 Laravel 有、Go 无）：

- **`PUT|POST /admin/classes/:id` 接受 `pet_series`**：取值白名单逐字同 Laravel `in:cosmic,pokemon,cute,treasure,mythic,all`（不合法 → 422「参数错误」+ `errors.pet_series`），并**按键合并**进 `class_rooms.settings`（保留其它键）；新增 `services.ValidClassRoomSeries` / `Admin.UpdateClassWithSettings`。
- **已清零的 Laravel-only 路由**（最后一批补齐）：`GET /api/v1/common/score-categories`（`StudentController::scoreCategories`，分类字典，已与 `ScoreRuleService::CATEGORY_LABELS` 唯一真源对齐）与 `GET|POST /api/v1/wechat-work/callback`（`WechatWorkWebhookController::verify/receive`，企微回调验签 / AES 解密 / 审批事件接收）。
- **`GET /teacher/students` 口径修正（早前批次顺带）**：原 Go 实现只返回管辖班级内 `status=active` 的学生、按 `class_id/student_no` 排序且无 `search`/分页；现按 Laravel `StudentService::list` 改为：**不过滤 status**、支持 `?search=`（姓名或学号，条件显式加括号以免 `AND` 抢在 `OR` 前导致越权行被带出）、`?page=` 分页每页固定 **50**（Laravel `paginate(50)`，`per_page` 不生效）、按 `name` 升序、行内追加 `pet_species/pet_level/pet_name`、响应带 `meta`（`message:"ok"` 为仓库信封额外字段）。

- **班级响应输出 `settings` 对象**：新增 `services.ClassRoomView` / `ClassRoomViewOf`，`GET /admin/classes`、`POST /admin/classes`、`PUT|POST /admin/classes/:id`、`GET /admin/classes/:id`（`AdminClassDetail.Settings`）统一输出对象形式的 settings——此前模型层 `json:"-"` 导致前端班级页永远显示「不限制」。
- **排行榜默认班级**：`/teacher/leaderboard/{total|weekly|pet-level}` 原先**必填** `class_id`（缺参 400），前端与 Laravel 都不传该参数；现按 Laravel 改为**默认取教师可管理班级的第一个**（`$classIds->first()`），无任何班级时返回 `{data: []}`；显式传 `class_id` 仍支持（Go 端可选扩展，越权/不存在 → 404）。

统计口径（**方法级**）：把 `frontend-vue/src/**` 中出现的 `/api/v1/...` 字面量路径（模板变量归一为 `{p}`、去掉 query、去重）连同**推断出的 HTTP 方法**（`apiGet/apiPost/apiPut/apiDelete` 或 `fetch({method})`）与 Go 路由表逐条比对。路由总数 196 → 216 → **211**（随后删除 4 条别名路由；`engine.Routes()` 实测，见 `internal/router/third_party_routes_test.go::TestThirdPartyRouteCount`）。前端修复后仅剩 2 条**正则碎片**：`POST /api/v1/auth/teacher`（`baseUrl` 字符串拼接被误判）与 `GET /api/v1/teacher/leaderboard/{p}`（模板变量，运行时只会取 `total|weekly|pet-level`）。
### 别名路由清理（与 Laravel 逐条对齐）

第三次全量比对（`engine.Routes()` 实测导出 vs. 解析 `backend/routes/api.php` 得到的 210 条）后，删除了 4 条 **Laravel 中不存在**的 Go 端别名路由，路由总数 216 → **211**：

| 已删除 | 原用途 | 改为（Laravel 真实路径） |
|--------|--------|--------------------------|
| `POST /api/v1/auth/login` | 早期「统一登录」 | `POST /auth/teacher/login`、`POST /auth/admin/login`、`POST /auth/class/login`（三者均带角色校验） |
| `GET /api/v1/auth/me` | 当前用户信息 | Laravel 无此路由；Go 测试的探针改用 `GET /api/v1/auth/bindings` |
| `GET /api/v1/teacher/classes` | 教师班级列表 | `GET /api/v1/teacher/my-classes`（`AdminCatalog` 口径，机器人返回全校班级） |
| `GET|POST /api/v1/admin/classes/:id/students` | 班级学生列表 / 建学生 | `GET|POST /api/v1/admin/students`（`?class_id=`） |

随路由一并删除的死代码：`handlers.Login`、`handlers.Me`、`handlers.TeacherClasses`、`handlers.AdminListStudents`、`handlers.AdminCreateStudent`、`services.Admin.ListStudents`。

保留 `GET /health`：Laravel `routes/web.php` 亦有同名路由（根级，不属于 `/api`），且 `docker/Dockerfile` 的 `HEALTHCHECK` 依赖它。

客户端影响：`mini-program/`（10 页小程序）历史代码里仍写着 `/auth/login`、`/teacher/classes`、`/teacher/scores/stats/{id}` 等路径，但它的 `globalData.apiBase` 是占位域名 `https://your-domain.com/api` 且与 `/api/v1/...` 字面量叠用（会拼成 `/api/api/v1/...`），对 Laravel 同样不可用，属未完成的客户端；Vue 前端与 `mini-program/pages/scores/scores.js` 真正在用的接口不受本次清理影响。
## 与原 Laravel 后端的关系

- 新后端位于独立的 `backend-go/` 模块（`github.com/RealKiro/learnstar-planet/backend-go`），**不影响** `backend/`。
- 数据模型**重新设计**，不要求与旧库兼容；API 契约可调整（前端将适配）。
