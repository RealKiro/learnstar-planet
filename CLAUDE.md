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
| 框架 | Laravel 11 | PHP 8.3，RESTful API |
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

- Docker 多阶段构建（Node 22 + PHP 8.3-FPM + Nginx + Supervisor）
- Docker Compose 编排（app + MySQL 8.0 + Redis 7）
- GitHub Container Registry (GHCR) 镜像托管
- CI/CD: GitHub Actions + Gitee Go

### 后端规范（已实施）

- **Form Requests**: 登录/批量创建/班级/导入验证逻辑已从控制器中抽离到独立请求类
- **API Resources**: UserResource、ClassRoomResource、StudentResource、SchoolResource 提供一致的 JSON 输出
- **异常处理**: 所有异常统一返回 JSON（ModelNotFound→404、ValidationException→422、Auth→401/403、Throwable→500/503）
- **速率限制**: 登录端点限制 `throttle:6,1`，API 整体限制 `throttle:api`
- **API 版本控制**: `/api/v1/*` 前缀 + 向后兼容旧路由
- **Eloquent Scopes**: User/ClassRoom/Student 模型添加 `active()`、`byRole()`、`bySchool()`、`byClass()` 作用域

---

## 项目结构

```
learnstar-planet/                    # 真正的项目根（含 .git）
├── docker-compose.yml              # Docker Compose 编排
├── .env.example                    # 环境变量模板
├── README.md                       # 详尽文档（800+ 行）
├── serverless-analysis.md          # 无服务器部署可行性分析
├── LICENSE                         # MIT
│
├── frontend-vue/                   # Vue 3 前端 SPA
│   ├── vite.config.ts              # Vite 构建配置 + API 代理
│   ├── src/
│   │   ├── types/index.ts          # 完整 TypeScript 类型定义
│   │   ├── utils/
│   │   │   ├── api.ts              # Axios 封装（token 拦截、401 处理）
│   │   │   └── constants.ts        # 工具函数与常量
│   │   ├── stores/
│   │   │   ├── auth.ts             # 认证状态（Pinia）
│   │   │   ├── theme.ts            # 暗色模式
│   │   │   ├── toast.ts            # Toast 消息
│   │   │   └── crud.ts             # 通用 CRUD Store 工厂
│   │   ├── router/index.ts         # 路由树 + 角色守卫
│   │   ├── layouts/                # 布局组件
│   │   │   ├── TeacherLayout.vue   # 教师端侧边栏布局
│   │   │   ├── AdminLayout.vue     # 管理端侧边栏布局
│   │   │   └── ParentLayout.vue    # 家长端侧边栏布局
│   │   └── pages/                  # 页面组件（按角色分组）
│   │       ├── landing/            #   首页 Landing
│   │       ├── auth/               #   登录页
│   │       ├── teacher/            #   教师端 16 个页面
│   │       ├── admin/              #   管理端 8 个页面
│   │       └── parent/             #   家长端 5 个页面
│   └── package.json
│
├── backend/                        # Laravel 11 后端
│   ├── app/
│   │   ├── Models/                 # 17 个 Eloquent 模型
│   │   ├── Http/Controllers/Api/  # 5 个 API 控制器
│   │   ├── Services/              # 4 个业务服务
│   │   ├── Livewire/Teacher/      # 2 个 Livewire 组件
│   │   ├── Events/                # 2 个事件类
│   │   └── Http/Middleware/       # RoleMiddleware
│   ├── database/
│   │   ├── migrations/            # 5 个迁移（21 张表）
│   │   └── seeders/               # AdminUserSeeder, DatabaseSeeder
│   ├── routes/
│   │   ├── api.php                # ~130 个 API 端点
│   │   ├── web.php                # Web 路由
│   │   └── console.php            # 控制台路由
│   └── config/                    # Laravel 配置
│
├── mini-program/                   # 微信小程序
│   └── pages/                     # 14 个页面
│
├── pwa/                            # PWA 配置
│   ├── manifest.json
│   └── sw.js
│
├── .github/workflows/              # GitHub Actions
│   ├── ci.yml
│   ├── docker.yml
│   ├── deploy.yml
│   └── security.yml
│
└── .gitee/                         # Gitee CI/CD
    └── workflow.yml
```

---

## 数据库架构（21 张表）

### 基础表（迁移 1）

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `schools` | 学校 | name, code, address, logo, settings(json) |
| `users` | 用户 | role(enum: school_admin/teacher/parent), username, password(bcrypt), nickname, avatar |
| `third_party_bindings` | OAuth绑定 | user_id, platform, open_id, union_id (复合唯一索引) |
| `class_rooms` | 班级 | grade, year, teacher_id, max_students |
| `students` | 学生 | class_id, parent_id, total_score, student_no, gender |
| `pets` | 宠物 | student_id(一对一), level, exp, mood, species, last_fed_at |
| `score_rules` | 积分规则 | class_id/school_id, name, points, category, is_penalty |
| `scores` | 积分记录 | student_id, class_id, rule_id, points, reason, giver_id |
| `score_logs` | 审计日志 | score_id, balance_before, balance_after |
| `notices` | 公告 | class_id/school_id, title, content, publisher_id, is_published |
| `shop_items` | 商品 | class_id, name, cost, stock |
| `shop_redemptions` | 兑换 | student_id, item_id, status(enum: pending/approved/rejected/delivered) |

### 教室工具表（迁移 2）

| 表名 | 说明 |
|------|------|
| `broadcasts` | 实时广播（type: banner/popup/fullscreen, voice, loop, duration）|
| `attendances` | 考勤（status: present/late/leave/absent, 唯一约束: class_id+student_id+date）|
| `homework_collections` | 作业（title, deadline, submission_types, qr_token）|
| `homework_submissions` | 作业提交（student_id, submitted_at, content, files）|
| `question_banks` | 题库（teacher_id, subject, is_public, question_count）|
| `questions` | 题目（type: single/multi/fill/truefalse/short, options(json), answer, points）|
| `quizzes` | 测验（class_id, bank_id, time_limit, auto_grade, realtime_stats）|
| `quiz_submissions` | 测验提交（answers(json), score）|
| `grades` | 成绩（exam_name, subject, 唯一约束: class_id+student_id+exam_name+subject, rank）|

---

## API 架构（~130 个端点）

路由按角色分组在 `backend/routes/api.php`：

### `/api/auth/*` — 认证（公开 + 需认证）
- `POST teacher/login` — 教师账号密码登录
- `POST admin/login` — 管理员账号密码登录
- `POST teacher/login/{wechat|wechat-work|qq|renren}` — 第三方扫码登录
- `POST teacher/bind-after-scan` — 扫码后绑定
- 需认证: logout, change-password, refresh, bind/{platform}, unbind/{platform}, GET bindings

### `/api/admin/*` — 学校管理员（role: school_admin）
- 学校 CRUD、教师/家长批量创建与管理
- 班级 CRUD + 批量创建 + 班主任分配
- 学生导入/CRUD/批量删除/批量转班
- 学年升级预览与执行
- 报表: overview, by-grade, by-class

### `/api/teacher/*` — 教师（role: teacher）
- dashboard, students (list + import + update)
- **scores/**: summary, give, batch-give, give-by-rule, history, rules CRUD
- **pets/**: class-overview, get pet, feed, rename
- **leaderboard/**: total, weekly, pet-level
- **shop/**: items CRUD, redemptions (approve/reject/deliver)
- **notices/**: CRUD + publish
- **reports/**: score-trend, pet-distribution, student-progress, export
- **broadcasts/**: list, send, get
- **attendance/**: today, start, set, summary
- **homework/**: CRUD + close + submissions + qr-code
- **quizzes/**: CRUD + start/stop + stats
- **question-banks/**: CRUD + add/get questions
- **grades/**: CRUD + stats + distribution
- **ai/**: chat, commands, usage

### `/api/parent/*` — 家长（role: parent）
- home, scores (detail + history), growth (log + timeline), pet (status + feed), ranking, notices

### `/api/common/*` — 公开
- pet-types, evolution-stages, score-categories

---

## 服务层

| 服务 | 职责 |
|------|------|
| `AuthService` | 教师账号创建（智能去重用户名/昵称）、管理员登录、密码管理 |
| `ScoreService` | DB 事务内积分操作（创建积分记录 + 更新余额 + 审计日志 + 宠物经验）；批量操作 |
| `LeaderboardService` | 基于 Redis ZSET 的排行榜（总分/本周/宠物等级）；MySQL 回退 |
| `PinyinService` | 中文名称拼音转换（昵称默认用）|

---

## 开发规范

### 代码风格
- **PHP**: PSR-12，使用 `.php-cs-fixer.php` 配置，Laravel 专用规则集
- **PHP 静态分析**: PHPStan Level 5（`phpstan.neon`），使用 Larastan
- **测试**: PHPUnit 11，Feature + Unit 测试
- **Git 提交**: 约定式提交（Conventional Commits）— `feat:`, `fix:`, `docs:`, `chore:` 等
- **JavaScript**: 使用 ESLint 检查

### 分支策略
- `main` — 生产就绪分支
- `feature/*` — 新功能开发
- `fix/*` — Bug 修复

### API 设计约定
- 所有 API 返回 JSON，格式: `{ data, message, meta }`
- 认证使用 Bearer Token（Sanctum）
- 角色中间件 `role:school_admin|teacher|parent` 控制访问
- 401 时前端自动清除 token 并跳转登录

### 数据库约定
- 迁移文件命名: `YYYY_MM_DD_HHMMSS_descriptive_name.php`
- 使用 Laravel 的 `Schema` builder，避免原始 SQL
- 外键使用 `constrained()` 和 `cascadeOnDelete()` 明确声明
- 所有表都有 `created_at` 和 `updated_at` 时间戳

---

## 常用命令

### 环境搭建

```bash
# 克隆项目
git clone https://github.com/YOUR_USERNAME/learnstar-planet.git
cd learnstar-planet

# 复制环境配置
cp .env.example .env
# 编辑 .env 修改数据库、Redis、管理员账号等配置

# Docker Compose 启动（全栈: app + MySQL + Redis）
docker-compose up -d

# 查看运行状态
docker-compose ps
```

### 后端开发

```bash
# 进入后端目录
cd backend

# 安装 PHP 依赖
composer install

# 运行迁移
php artisan migrate

# 创建管理员（首次部署）
php artisan db:seed --class=AdminUserSeeder

# 启动开发服务器
php artisan serve

# 运行测试
php artisan test

# 代码格式化
vendor/bin/php-cs-fixer fix

# 静态分析
vendor/bin/phpstan analyse

# 查看路由列表
php artisan route:list

# 队列处理
php artisan queue:work

# Horizon 面板
php artisan horizon
```

### Docker 常用操作

```bash
# 查看日志
docker-compose logs -f app

# 