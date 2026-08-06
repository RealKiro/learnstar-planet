# CLAUDE.md — 学趣星球 (LearnStar Planet)

## 项目概述

学趣星球是一个开源的班级管理与学生激励系统，定位为商业产品"班宠星球"的免费替代品。教师使用它来管理学生积分、宠物养成、考勤、作业、测验、成绩等日常课堂事务。

- **许可证**: MIT
- **仓库**: https://github.com/RealKiro/learnstar-planet
- **定位**: 全免费、自托管、多端适配的班级管理平台

---

## 技术架构

### 后端

| 层面 | 技术 | 说明 |
|------|------|------|
| 框架 | Laravel 12（升级中） | PHP 8.5，RESTful API |
| 认证 | Laravel Sanctum 4 | API Token 认证，按角色隔离 |
| 权限 | spatie/laravel-permission 6 | 基于角色的权限（school_admin / teacher / parent） |
| 实时 | Livewire 3 + Flux 2 | 教师仪表盘与积分管理的动态 UI |
| 缓存 | Redis (Predis 2) | 排行榜用 ZSET，队列用 Horizon 5 |
| 数据库 | MySQL 8.0+ / MariaDB 10.3+ / PostgreSQL 14+ / SQLite 3.8+ | 四种数据库均支持 |
| Excel | maatwebsite/excel 3 | Excel 导入导出 |
| PDF | barryvdh/laravel-dompdf 2 | PDF 报表导出 |
| 图片 | intervention/image 3 | 图片处理 |
| 日志 | spatie/laravel-activitylog 4 | 活动日志审计 |
| 队列 | Laravel Horizon 5 | Redis 队列监控 |

### 前端

| 端 | 技术 | 说明 |
|------|------|------|
| Web | Vue 3 + Vite + TypeScript | `frontend-vue/` 目录，组件化 SPA，需 `npm run build` 构建 |
| 小程序 | 微信小程序原生 | `mini-program/` 目录，14 个页面，教师端 + 家长端 |
| PWA | 原生 Service Worker | `pwa/` 目录，离线缓存、推送通知、后台同步 |

### 基础设施

- Docker 多阶段构建（Node 22 + PHP 8.5-FPM + Nginx + Supervisor）
- Docker Compose 编排（app + MariaDB 10 + Redis 8，mariadb:10 / redis:alpine）
- GitHub Container Registry (GHCR) 镜像托管
- CI/CD: GitHub Actions + Gitee Go

### 后端规范（已实施）

- **Form Requests**: 登录/批量创建/班级/导入验证逻辑已从控制器中抽离到独立请求类
- **API Resources**: UserResource、ClassRoomResource、StudentResource、SchoolResource 提供一致的 JSON 输出
- **异常处理**: 所有异常统一返回 JSON（ModelNotFound→404、ValidationException→422、Auth→401/403、Throwable→500）
- **速率限制**: 登录端点限制 `throttle:6,1`，API 整体限制 `throttle:api`
- **API 版本控制**: `/api/v1/*` 前缀 + 向后兼容旧路由
- **Eloquent Scopes**: User/ClassRoom/Student 模型添加 `active()`、`byRole()`、`bySchool()`、`byClass()` 作用域

---

## 项目结构

\`\`\`
learnstar-planet/
├── frontend-vue/                   # Vue 3 + Vite + TypeScript SPA
│   ├── src/                        # 源代码（组件/路由/状态/类型）
│   ├── vite.config.ts              # Vite 构建配置
│   └── package.json                # Node.js 依赖
├── mcp-server/                     # MCP 服务器（AI 机器人集成）
│   ├── server.py                   #   MCP 协议服务器
│   └── README.md                   #   部署说明
├── docker-compose.yml              # Docker Compose 编排（应用 + MySQL + Redis）
├── .env.example                    # 环境变量模板
├── CLAUDE.md                       # 开发文档与规范
├── LICENSE                         # MIT 开源许可证
├── README.md                       # 项目说明
│
├── backend/                        # Laravel 11 API
│   ├── Dockerfile                  # 生产环境多阶段构建（Node + PHP + Nginx）
│   ├── Dockerfile.dev              # 开发环境构建
│   ├── app/
│   │   ├── Models/                 # 23 个 Eloquent 模型
│   │   ├── Http/Controllers/Api/  # 7 个 API 控制器
│   │   ├── Services/              # 9 个业务服务
│   │   ├── Http/Requests/         # Form Request 验证类
│   │   ├── Http/Resources/        # JsonResource 响应类
│   │   └── Livewire/              # 2 个 Livewire 组件
│   ├── database/migrations/       # 21 个迁移（含 2026_08_05 计费/班级码/汇率）
│   └── routes/api.php             # ~210 个 API 端点
│
├── mini-program/                   # 微信小程序
│   └── pages/                     # 14 个页面
│
├── pwa/                            # PWA 配置
│
├── .github/workflows/              # GitHub Actions CI/CD
│   ├── ci.yml
│   ├── docker.yml
│   ├── deploy.yml
│   └── security.yml
│
└── .gitee/                         # Gitee CI/CD
    └── workflow.yml
\`\`\`

---

## 数据库架构（21 张表）

### 基础表
| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `schools` | 学校 | name, code, address, logo, settings(json) |
| `users` | 用户 | role, username, password(bcrypt), nickname, avatar |
| `third_party_bindings` | OAuth绑定 | user_id, platform, open_id |
| `class_rooms` | 班级 | grade, year, teacher_id, max_students |
| `students` | 学生 | class_id, parent_id, total_score, student_no |
| `pets` | 宠物 | student_id(一对一), level, exp, mood, species |
| `score_rules` | 积分规则 | class_id/school_id, name, points, category |
| `scores` | 积分记录 | student_id, class_id, rule_id, points, reason |
| `score_logs` | 审计日志 | score_id, balance_before, balance_after |
| `notices` | 公告 | class_id/school_id, title, content, publisher_id |
| `shop_items` | 商品 | class_id, name, cost, stock |
| `shop_redemptions` | 兑换 | student_id, item_id, status |

### 教室工具表
| 表名 | 说明 |
|------|------|
| `broadcasts` | 实时广播（banner/popup/fullscreen）|
| `attendances` | 考勤（present/late/leave/absent）|
| `homework_collections` | 作业布置 |
| `homework_submissions` | 作业提交 |
| `question_banks` | 题库 |
| `questions` | 题目 |
| `quizzes` | 测验 |
| `quiz_submissions` | 测验提交 |
| `grades` | 成绩 |

---

## API 架构（~210 个端点）

### `/api/v1/auth/*` — 认证
- POST teacher/login, admin/login, teacher/login/{platform}
- POST bind-after-scan
- 需认证: logout, change-password, refresh, bind/{platform}, unbind/{platform}, GET bindings

### `/api/v1/admin/*` — 学校管理员
- 学校 CRUD、教师/家长批量创建与管理
- 班级 CRUD + 批量创建 + 班主任分配
- 学生导入/CRUD/批量删除/批量转班
- 学年升级预览与执行
- 报表: overview, by-grade, by-class

### `/api/v1/teacher/*` — 教师
- dashboard, students
- scores/ (summary, give, batch-give, history, rules CRUD)
- pets/ (class-overview, pet, feed, rename)
- leaderboard/ (total, weekly, pet-level)
- shop/ (items CRUD + redemptions)
- notices/ (CRUD + publish)
- reports/ (trend, distribution, progress, export)
- broadcasts/, attendance/, homework/, quizzes/, question-banks/, grades/, ai/

### `/api/v1/parent/*` — 家长
- home, scores, growth, pet, ranking, notices

### `/api/v1/common/*` — 公开
- pet-types, evolution-stages, score-categories

---

## 服务层

| 服务 | 职责 |
|------|------|
| `AuthService` | 教师账号创建（智能去重用户名/昵称）、管理员登录、密码管理 |
| `ScoreService` | DB 事务内积分操作（创建积分记录 + 更新余额 + 审计日志 + 宠物经验）；批量操作 |
| `LeaderboardService` | 基于 Redis ZSET 的排行榜（总分/本周/宠物等级）；MySQL 回退 |
| `PinyinService` | 中文名称拼音转换（昵称默认用）|
| `AiService` | AI 对话调用（30+ 供应商 OpenAI 兼容分发），返回 answer + token 拆分 |
| `AiBilling\AiBillingService` | Token 计费：OpenAI/DeepSeek 官方账单查询 + 本地精确计价兜底 |
| `CurrencyService` | 多币种体系（科学币/读书币/班级积分）：积分兑换、钱包互兑、消费 |
| `DisplayEventService` | 教室大屏 SSE 实时广播（积分变动/横幅/通知），班级码缓存 |
| `ThirdParty\ThirdPartyManager` + `WeChatWorkProvider`/`DingTalkProvider`/`LarkProvider` | 多平台第三方扫码（每校选一平台）：authUrl/getUserByCode/fetchContacts |
| `WechatWorkService` / `WechatWorkAttendanceService` | 企业微信集成：免注册登录、请假同步 |

---

## 开发规范

### 代码风格
- **PHP**: PSR-12，使用 `.php-cs-fixer.php` 配置
- **PHP 静态分析**: PHPStan Level 5
- **测试**: PHPUnit 11
- **Git 提交**: 约定式提交（Conventional Commits）
- **JavaScript**: ESLint

### 分支策略
- `main` — 生产就绪分支
- `feature/*` — 新功能开发
- `fix/*` — Bug 修复

### API 设计约定
- 所有 API 返回 JSON，格式: `{ data, message, meta }`
- 认证使用 Bearer Token（Sanctum）
- 角色中间件 `role:school_admin|teacher|parent` 控制访问
- 401 时前端自动清除 token 并跳转登录
- API 版本化：前端统一调用 `/api/v1/*`，后端保留向后兼容旧路由
- 登录端点速率限制：`throttle:6,1`
- 异常处理统一返回 JSON
- Docker 多阶段构建：Node 22 → Vue 3 SPA，PHP 8.5 → Laravel API

---

## 常用命令

### 环境搭建
\`\`\`bash
git clone https://github.com/YOUR_USERNAME/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
\`\`\`

### 后端开发
\`\`\`bash
cd backend
composer install
php artisan migrate
php artisan db:seed --class=AdminUserSeeder
php artisan serve
php artisan test
vendor/bin/php-cs-fixer fix
vendor/bin/phpstan analyse
\`\`\`

### Docker 常用操作
\`\`\`bash
docker-compose logs -f app
docker-compose exec app bash
docker-compose restart app
docker-compose pull && docker-compose up -d
docker-compose exec mysql mysqldump -u root -p learnstar > backup.sql
\`\`\`

### 前端开发
\`\`\`bash
cd frontend-vue
npm install
npm run dev         # 开发服务器 http://localhost:5173
npm run typecheck   # TypeScript 类型检查
npm run build       # 生产构建
npm run build:deploy # 输出到 ../backend/public/
\`\`\`

---

## 设计系统（CSS 变量）

项目采用"星空探索"主题，定义在 `frontend-vue/src/assets/style.css`：

| 类别 | 变量 | 值 |
|------|------|-----|
| 品牌主色 | `--primary` | `#4F46E5` 星空靛蓝 |
| 品牌辅色 | `--secondary` | `#F59E0B` 星光琥珀 |
| 功能色 | `--accent` / `--danger` / `--info` | `#10B981` / `#EF4444` / `#3B82F6` |
| 背景 | `--bg` / `--bg-card` | `#F8FAFC` / `#FFFFFF` |
| 文字 | `--text` / `--text-secondary` | `#1E293B` / `#64748B` |
| 圆角 | `--radius-sm/md/lg/xl` | `8px` / `12px` / `20px` / `28px` |
| 阴影 | `--shadow-sm/md/lg/glow` | 4 级阴影深度 |

支持暗色模式，通过给 `<html>` 添加 `.dark` 类切换。

---

## 关键设计决策与注意事项

1. **前端已完成 Vue 3 重构**: `frontend-vue/` 采用 Vue 3 + Vite + TypeScript + Pinia + Vue Router
2. **无自注册**: 所有账号由管理员在后台创建分配
3. **角色严格隔离**: 管理员/教师/家长界面和 API 完全不同
4. **第三方登录仅限教师**: 管理员不支持第三方扫码；后台学校设置可勾选启用平台（企业微信/钉钉/飞书/微信/QQ/人人通空间），登录页按配置动态展示，存储于 `schools.settings.enabled_third_party_platforms`
5. **AI 功能可选**: 不配置 AI API Key 不影响核心功能
6. **排行榜使用 Redis ZSET**: 有 MySQL 回退方案
7. **学年升级不可逆**: 预览→确认→执行，
8. **家长功能目前简约**: 以查看为主（积分、宠物、通知、排名）
9. **PWA 主要用于离线缓存**: Service Worker 缓存静态资源，网络优先策略
10. **MCP Server 支持 AI 机器人**: `mcp-server/` 提供标准 MCP 协议服务器，可对接 AstrBot + NapCatQQ/Lagrange 实现 QQ/微信 自然语言积分管理
11. **Laravel 版本（11 → 12 升级中，2026-08-06 启动）**:

    **现状**：原锁定 `laravel/framework: ^11.0`（PHP 8.3），代码为 L11 新结构（`bootstrap/app.php` 的 `Application::configure()`，无 `app/Http/Kernel.php`）。2026-08-06 升级 **Laravel 12 + PHP 8.5**：`composer.json` 声明 `laravel/framework ^12.0`、`php ^8.5`，Dockerfile 基础镜像 `php:8.5-cli-alpine`。

    **为什么 12 不是 13**：先尝试 13，但生态未跟上（`laravel/tinker` 2.x 最高 illuminate ^12、`nunomaduro/collision` 无 v9、dusk/larastan 未确认 L13）→ **降回 L12**（发布一年多，生态全兼容，PHP `^8.2` 约束允许 8.5）。等生态跟进后可再升 13。

    **升级前置（已完成 2026-08-06）**：
    - ✅ **消除 bootstrap env() 直调**：`bootstrap/app.php` 的 `env('APP_ENV')`/`env('TRUST_ALL_PROXIES')` 改为 `app()->environment('production')` + `config('proxy.trust_all')`（新增 `config/proxy.php`）。**关键理解**：config 文件内的 `env()` 在 `config:cache` 时会被正确解析并序列化，真正失效的是**非 config 位置（bootstrap/控制器/服务）**的 env() 直调。
    - ✅ **依赖约束已按 L12 调整**：dompdf ^3.1（v2 只支持 ≤11）、permission ^7.0、collision ^8.0、tinker ^2.0 已加回；horizon/octane/activitylog/dusk 用宽松 ^5.0/^2.0/^4.0/^8.0。**仍待 PHP 环境 composer 解析验证**。

    **升级手册（需在 PHP/Docker 环境执行）**：
    ```bash
    cd backend
    composer update --with-all-dependencies
    php artisan migrate --force
    php artisan test
    php artisan config:cache   # 验证 env() 修复后正常
    php artisan route:list --path=api
    php artisan optimize
    ```

    **L11→L12 破坏性变更核对清单**：
    | 项 | 说明 |
    |---|---|
    | 中间件/异常 | `withMiddleware`/`withExceptions` 结构 L11→L12 延续，基本兼容 |
    | config:cache env() | bootstrap 已消除；config 文件内 env() 缓存时正常 |
    | PHP | L12 的 `php ^8.2` 约束允许 8.5，无碍 |
    | 依赖包 | dompdf 需 ^3.1（v2 不支持 L12）、permission ^7.0，其余 ^4.0/^5.0 兼容 |
    | Collection | L12 部分 Collection 方法更严格，需回归 |
    | 广播/事件 | `ShouldBroadcastNow` + Redis 广播，L12 兼容 |

    **结论**：升级前必须确保 bootstrap 无 env() 直调（已完成），composer update 在 PHP 环境跑通并 `php artisan test` 全绿后再推送（否则 CI 的 composer install 会挂）。
