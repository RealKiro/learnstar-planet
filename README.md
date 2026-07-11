# 🌌 学趣星球 (LearnStar Planet)

> 开源版班级管理系统 - 积分激励 · 宠物养成 · AI助教 · 全免费

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social)](https://github.com/RealKiro/learnstar-planet)

---

## 📖 目录

- [项目介绍](#项目介绍)
- [功能特性](#功能特性)
- [项目结构](#项目结构)
- [快速部署](#快速部署)
- [配置说明](#配置说明)
- [FAQ](#faq)
- [贡献指南](#贡献指南)

---

## 项目介绍

**学趣星球**是一款开源的班级管理与学生激励系统，旨在为教师提供一个有趣、高效、免费的工具，通过**积分激励**、**宠物养成**、**AI助教**等创新功能，提升学生学习积极性。

### 🎯 设计理念

- **全免费**：无任何付费功能，所有功能完全开放
- **易部署**：支持 Docker 一键部署，预构建镜像
- **多端适配**：Web、微信小程序、PWA
- **多数据库**：MySQL、PostgreSQL、SQLite、MariaDB
- **AI 赋能**：集成通义千问等 AI 平台（可选）

### 🆚 与原版对比

| 功能 | 原版班宠星球 | 学趣星球 |
|------|--------------|----------|
| **费用** | VIP 收费 | 💯 全免费 |
| **源码** | 闭源 | 💯 开源 |
| **部署** | 仅云端 | 💯 自托管+云端 |
| **AI 功能** | 无 | 💯 有（可选） |
| **多数据库** | 仅 MySQL | 💯 4 种数据库 |
| **版权** | 商业版权 | 💯 MIT 许可证 |

---

## 功能特性

### 🎓 教师端

#### 核心功能
- **积分管理系统**：自定义积分规则、批量/单个打分、积分排行榜
- **宠物养成系统**：班级宠物、学生领养、宠物进化、互动玩法
- **通知公告系统**：发布通知、已读统计、紧急通知弹窗

#### 教室小喇叭（新增）
- 📢 **实时广播**：顶部横幅/弹窗/全屏三种模式
- ✅ **智能考勤**：一键点名、4种状态（到/迟/缺/事）
- 📷 **扫码收作业**：生成二维码、进度跟踪
- 📝 **在线答题**：题库管理、自动判分
- 📊 **成绩管理**：成绩录入、统计分析、分布图
- 🤖 **AI 助教**：智能问答、学习建议（需配置 API Key）

#### 管理功能
- **学校/班级管理**：创建学校、分配班级、管理教师
- **账号管理**：批量创建账号、重置密码（无自注册）
- **数据导出**：积分记录、考勤记录、成绩报表

### 👨👩👧 家长端

- **积分查看**：实时查看孩子积分变动
- **宠物查看**：查看孩子领养的宠物状态
- **通知接收**：接收教师发布的通知
- **成绩查看**：查看孩子考试成绩

### 🔐 第三方登录

教师账号支持绑定以下平台，绑定后可快捷登录：
- 个人微信扫码登录
- 企业微信扫码登录
- QQ 扫码登录
- 人人通空间（账号密码 / 客户端扫码登录）

---

## 技术架构

### 后端
- **框架**：Laravel 11
- **API**：RESTful API
- **缓存**：Redis（排行榜用 ZSET 实现）
- **数据库**：MySQL / PostgreSQL / SQLite / MariaDB

### 前端
- **Web 端**：Vue 3 + Vite + TypeScript（SPA，需构建步骤）
- **小程序**：微信小程序原生开发
- **PWA**：支持离线访问

### 部署
- **容器化**：Docker + Docker Compose
- **CI/CD**：GitHub Actions + Gitee Go
- **镜像**：GitHub Container Registry (GHCR)

---

## 项目结构

```
learnstar-planet/
├── frontend-vue/                     # Web 前端 - Vue 3 + Vite + TypeScript SPA
│   ├── src/                          #   源代码（组件/路由/状态/类型）
│   ├── vite.config.ts                #   Vite 构建配置
│   └── package.json                  #   Node.js 依赖
├── docker-compose.yml              # Docker Compose 编排（应用 + MySQL + Redis）
├── .env.example                    # 环境变量模板（复制为 .env 后修改）
├── CLAUDE.md                       # 开发文档与规范
├── LICENSE                         # MIT 开源许可证
├── README.md                       # 本文档
│
├── backend/                        # 后端 - Laravel 11 API 服务
│   ├── Dockerfile                  # 生产环境多阶段构建（Node + PHP 8.3 + Nginx + Supervisor）
│   ├── Dockerfile.dev              # 开发环境构建
│   ├── .dockerignore               # Docker 构建忽略文件
│   ├── composer.json               # PHP 依赖声明（Laravel/Livewire/Flux 等）
│   ├── .env.example                # 后端环境变量模板
│   ├── .env.staging                # 预发布环境配置
│   ├── .php-cs-fixer.php           # PHP CS Fixer 代码格式化规则（PSR-12）
│   ├── phpstan.neon                # PHPStan 静态分析配置（Level 5）
│   ├── phpunit.xml                 # PHPUnit 测试配置
│   │
│   ├── app/                        # 应用核心代码
│   │   ├── Models/                 # Eloquent 数据模型（17 个）
│   │   │   ├── School.php          #   学校
│   │   │   ├── ClassRoom.php       #   班级
│   │   │   ├── User.php            #   用户（教师/家长/管理员）
│   │   │   ├── Student.php         #   学生
│   │   │   ├── ThirdPartyBinding.php #  第三方账号绑定（微信/企微/QQ/人人通）
│   │   │   ├── Pet.php             #   宠物
│   │   │   ├── Score.php           #   积分
│   │   │   ├── ScoreRule.php       #   积分规则
│   │   │   ├── ScoreLog.php        #   积分变动日志
│   │   │   ├── ShopItem.php        #   商城物品
│   │   │   ├── ShopRedemption.php  #   兑换记录
│   │   │   ├── Notice.php          #   通知公告
│   │   │   ├── Broadcast.php       #   实时广播（教室小喇叭）
│   │   │   ├── Attendance.php      #   考勤记录
│   │   │   ├── HomeworkCollection.php # 作业收集
│   │   │   ├── Quiz.php            #   在线答题/题库
│   │   │   └── Grade.php           #   成绩
│   │   │
│   │   ├── Services/               # 业务逻辑服务层
│   │   │   ├── AuthService.php     #   认证服务（教师/管理员登录拆分）
│   │   │   ├── ScoreService.php    #   积分管理服务
│   │   │   └── LeaderboardService.php # 排行榜服务（Redis ZSET）
│   │   │
│   │   ├── Livewire/               # Livewire 组件
│   │   │   └── Teacher/
│   │   │       ├── Dashboard.php   #   教师仪表盘组件
│   │   │       └── ScoreManager.php #  积分管理组件
│   │   │
│   │   └── Events/                 # 事件类
│   │       ├── NoticePublished.php #   通知发布事件（SSE 推送）
│   │       └── ScoreChanged.php    #   积分变动事件
│   │
│   ├── config/                     # Laravel 配置文件
│   │   ├── database.php            #   多数据库配置（MySQL/PgSQL/SQLite）
│   │   └── cache.php               #   缓存配置（Redis）
│   │
│   ├── database/
│   │   ├── migrations/             # 数据库迁移
│   │   │   ├── 2025_01_01_000001_create_all_tables.php  # 基础表（12 张）
│   │   │   └── 2025_01_01_000002_create_classroom_tools_tables.php # 教室工具表（8 张）
│   │   └── seeders/                # 数据填充（预留）
│   │
│   ├── routes/
│   │   └── api.php                 # API 路由（按角色分组：auth/teacher/parent/admin）
│   │
│   ├── resources/views/            # Blade 视图模板
│   │   ├── layouts/app.blade.php   #   基础布局
│   │   └── livewire/teacher/       #   Livewire 组件视图
│   │
│   ├── tests/                      # 测试用例
│   │   ├── Feature/TeacherApiTest.php  # 功能测试
│   │   └── Unit/ScoreServiceTest.php   # 单元测试
│   │
│   └── docker/                     # Docker 相关配置
│       ├── nginx/default.conf      #   Nginx 配置（单容器模式，监听 8080）
│       ├── supervisor/supervisord.conf # Supervisor 配置（管理 Nginx+PHP-FPM+Queue）
│       └── scripts/entrypoint.sh   #   容器启动脚本（迁移/缓存/启动）
│
├── mini-program/                   # 微信小程序端
│   ├── app.js                      # 小程序入口
│   ├── app.json                    # 全局配置（页面路由、窗口样式）
│   ├── app.wxss                    # 全局样式
│   ├── project.config.json         # 项目配置
│   ├── sitemap.json                # 小程序索引配置
│   └── pages/                      # 页面
│       ├── login/                  #   登录页（教师/管理员双面板 + 微信登录）
│       ├── index/                  #   首页
│       └── dashboard/              #   仪表盘
│
├── pwa/                            # PWA 渐进式 Web 应用
│   ├── manifest.json               # PWA 清单（应用名/图标/主题色）
│   └── sw.js                       # Service Worker（离线缓存）
│
├── .github/workflows/              # GitHub Actions CI/CD 工作流
│   ├── ci.yml                      #   持续集成（PHP 测试 + 代码质量 + 前端检查）
│   ├── docker.yml                  #   Docker 镜像构建并推送到 GHCR
│   ├── deploy.yml                  #   自动部署到服务器（SSH + Docker）
│   └── security.yml                #   安全审计（依赖审查 + Composer 审计）
│
└── .gitee/                         # Gitee（码云）CI/CD
    └── workflow.yml                #   Gitee Go 流水线配置
```

---

## 快速部署

> 💡 **推荐方式**：Fork 本仓库后自行构建镜像（方式一），完全掌控代码和部署流程，免费使用 GitHub Actions 构建镜像并推送到你的 GHCR 命名空间。

---

### 🔧 方式一：Fork 仓库后自行构建镜像（推荐 ✅）

适合需要在原项目基础上做定制开发，或希望完全掌控部署流程的用户。Fork 后 GitHub Actions 会自动构建 Docker 镜像并推送到你的 GHCR，无需本地配置构建环境。

#### 1. Fork 仓库

1. 访问 [learnstar-planet 仓库](https://github.com/RealKiro/learnstar-planet)
2. 点击右上角 **Fork** 按钮，将仓库 Fork 到你的 GitHub 账号下

#### 2. 启用 GitHub Actions

1. 进入你 Fork 后的仓库页面
2. 点击 **Actions** 标签页
3. 如果提示需要启用 Workflows，点击 **I understand my workflows, go ahead and enable them**
4. 确认以下工作流已就绪：
   - `Build & Push Docker Images` — 自动构建并推送 Docker 镜像
   - `CI - Test & Lint` — 代码质量检查

#### 3. 确认仓库权限

GitHub Actions 构建镜像需要 `packages: write` 权限，本仓库的工作流已内置配置：

1. 进入 **Settings → Actions → General**
2. 在 **Workflow permissions** 下选择 **Read and write permissions**
3. 勾选 **Allow GitHub Actions to create and approve pull requests**（可选）
4. 点击 **Save**

#### 4. 触发镜像构建

**方式 A — 推送代码自动构建**（修改代码后推送即可）：

```bash
# 克隆你 Fork 的仓库
git clone https://github.com/YOUR_USERNAME/learnstar-planet.git
cd learnstar-planet

# 做你的修改后提交推送
git add .
git commit -m "feat: my custom changes"
git push origin main
```

推送后 GitHub Actions 会自动触发构建，镜像将推送到：
```
ghcr.io/YOUR_USERNAME/learnstar-planet/backend:latest
```

**方式 B — 手动触发构建**：

1. 进入 **Actions → Build & Push Docker Images**
2. 点击 **Run workflow** → 选择 `main` 分支 → 点击 **Run workflow**
3. 等待构建完成（通常 5-10 分钟），可在 Actions 页面查看进度

#### 5. （可选）设置镜像为公开

默认情况下，GHCR 镜像是私有仓库。如需免登录拉取：

1. 访问 [GitHub Packages](https://github.com/YOUR_USERNAME?tab=packages)
2. 点击 `learnstar-planet/backend` 包
3. 进入 **Package settings**
4. 在 **Danger Zone → Change visibility** 中选择 **Public**

#### 6. 部署

**Using Docker Run（快速启动，需自备数据库）**

```bash
# 拉取镜像
docker pull ghcr.io/YOUR_USERNAME/learnstar-planet/backend:latest

# 创建 .env 文件（参考 .env.example，配置你的数据库地址）
# 启动容器（映射端口 8080）
docker run -d -p 8080:8080 --name learnstar-planet \
  --env-file .env \
  ghcr.io/YOUR_USERNAME/learnstar-planet/backend:latest

# 初始化数据库（自动创建管理员账号：admin / admin123456）
docker exec learnstar-planet php artisan migrate --force
docker exec learnstar-planet php artisan db:seed --class=AdminUserSeeder --force
```

**Using Docker Compose（推荐，全自动，含数据库）**

```bash
# 下载配置文件
curl -O https://raw.githubusercontent.com/YOUR_USERNAME/learnstar-planet/main/docker-compose.yml
curl -O https://raw.githubusercontent.com/YOUR_USERNAME/learnstar-planet/main/.env.example
cp .env.example .env

# 编辑 .env
nano .env
```

**必改项**：
```env
# 改为你的 GitHub 用户名（用于拉取你自行构建的镜像）
GITHUB_USERNAME=YOUR_USERNAME

# 改为你访问该服务的 IP 或域名（不要用 Docker 内部 IP）
# 只改 APP_HOST，APP_URL 会自动拼接 APP_PORT
# 公网服务器：your-server-ip
# 局域网 NAS（群晖/飞牛）：192.168.1.100
# 有域名/反向代理：learnstar.yourdomain.com
APP_HOST=your-server-ip

# 宿主机访问端口（容器内部固定为 8080，不要改 docker-compose 里的 8080）
# 如果 8080 被占用，改成 8081/8090/9000 等未被占用的端口
APP_PORT=8080

# APP_URL 由 APP_HOST 和 APP_PORT 自动拼接，一般无需手动修改
APP_URL=http://${APP_HOST}:${APP_PORT}

# 应用数据库账号密码（Laravel 应用用 learnstar 账号读写业务数据）
DB_PASSWORD=your_secure_password

# MySQL root 超级管理员密码（数据库维护/备份/恢复用，建议和 DB_PASSWORD 不同）
MYSQL_ROOT_PASSWORD=your_root_password

# 后台超级管理员账号（首次部署自动创建，生产环境请务必修改密码）
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123456
```

> ⚠️ **关于数据库配置**（重要）：
> - `DB_PASSWORD` 是 **Laravel 应用**连接数据库时用的密码（对应 `DB_USERNAME=learnstar`）
> - `MYSQL_ROOT_PASSWORD` 是 **MySQL root 管理员**密码，用于数据库维护、备份、紧急恢复
> - 这两个密码**建议不同**，不要把 root 密码泄露给应用
> - `DB_PORT=3306` 是**容器内部**端口，一般不需要改
> - `MYSQL_PORT` 是**宿主机**暴露给外部工具（如 Navicat/DBeaver）的端口，如果 3306 被占用，可改成 `3307`/`3308`
>
> 示例：宿主机 MySQL 端口改成 3307 后，Navicat 连接时主机填 `你的IP`，端口填 `3307`，但 `DB_PORT` 仍然保持 `3306`

```bash
# 启动服务（自动拉取镜像、启动应用+数据库+Redis、并自动创建管理员）
docker-compose up -d

# 查看运行状态
docker-compose ps

# 如果需要手动初始化数据库（自动创建管理员账号：admin / admin123456）
docker-compose exec app php artisan migrate --force
docker-compose exec app php artisan db:seed --class=AdminUserSeeder --force
```

#### 7. 后续更新

当你修改代码并推送后，GitHub Actions 会自动重新构建镜像：

```bash
# 拉取最新镜像
docker-compose pull

# 重启服务
docker-compose up -d

# 运行数据库迁移（如有）
docker-compose exec app php artisan migrate --force
```

> 💡 **提示**：如果上游仓库有更新，可以通过 GitHub 的 **Sync fork** 功能同步最新代码，同步后推送即可触发重新构建。

---

### 📦 方式二：使用预构建镜像（快速试用）

适合快速体验项目，或暂时不想 Fork 仓库的用户。直接用官方预构建镜像，5 分钟内完成部署。

#### 1. 下载配置文件

```bash
# 下载 docker-compose.yml 和 .env.example
curl -O https://raw.githubusercontent.com/RealKiro/learnstar-planet/main/docker-compose.yml
curl -O https://raw.githubusercontent.com/RealKiro/learnstar-planet/main/.env.example
```

#### 2. 修改配置

```bash
# 复制 .env.example 为 .env
cp .env.example .env

# 编辑 .env，修改以下配置
nano .env
```

**必改项**：
```env
# 修改为你访问该服务的 IP 或域名（不要用 Docker 内部 IP）
# 只改 APP_HOST，APP_URL 会自动拼接 APP_PORT
# 公网服务器：your-server-ip
# 局域网 NAS（群晖/飞牛）：192.168.1.100
# 有域名/反向代理：learnstar.yourdomain.com
APP_HOST=your-server-ip

# 宿主机访问端口（容器内部固定为 8080，不要改 docker-compose.yml 里的 8080）
# 如果 8080 被占用，改成 8081/8090/9000 等未被占用的端口
APP_PORT=8080

# APP_URL 由 APP_HOST 和 APP_PORT 自动拼接，一般无需手动修改
# 如果你使用 https 或自定义路径，可以手动覆盖
APP_URL=http://${APP_HOST}:${APP_PORT}

# 应用数据库账号密码（Laravel 应用用 learnstar 账号读写业务数据）
DB_PASSWORD=your_secure_password

# MySQL root 超级管理员密码（数据库维护/备份/恢复用，建议和 DB_PASSWORD 不同）
MYSQL_ROOT_PASSWORD=your_root_password

# 后台超级管理员账号（首次部署自动创建，生产环境请务必修改密码）
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123456
```

> 💡 **提示**：`GITHUB_USERNAME=realkiro` 已预置为源作者用户名（GHCR 镜像路径要求全小写），用于拉取官方预构建镜像，无需修改。

> ⚠️ **关于端口映射**（重要）：
> - `docker-compose.yml` 中的端口映射格式是：`宿主机端口:容器端口`
> - **容器内部端口固定为 8080**（Nginx 监听这个端口），通常不需要改它
> - **APP_HOST 只填 IP 或域名**，例如 `192.168.1.100`、`localhost`、`learnstar.yourdomain.com`
> - **APP_PORT 是宿主机端口**，你可以改成任意未被占用的端口（如 8081、8090、9000）
> - **APP_URL 由 APP_HOST 和 APP_PORT 自动拼接**（`http://${APP_HOST}:${APP_PORT}`），一般无需手动改
> - 如果 `APP_PORT=80`，可以省略端口写 `http://your-server-ip`，此时需要手动覆盖 APP_URL

> ⚠️ **关于数据库配置**（重要）：
> - `DB_PASSWORD` 是 **Laravel 应用**连接数据库时用的密码（对应 `DB_USERNAME=learnstar`）
> - `MYSQL_ROOT_PASSWORD` 是 **MySQL root 管理员**密码，用于数据库维护、备份、紧急恢复
> - 这两个密码**建议不同**，不要把 root 密码泄露给应用
> - `DB_PORT=3306` 是**容器内部**端口，一般不需要改
> - `MYSQL_PORT` 是**宿主机**暴露给外部工具（如 Navicat/DBeaver）的端口，如果 3306 被占用，可改成 `3307`/`3308`
>
> 示例：宿主机 MySQL 端口改成 3307 后，Navicat 连接时主机填 `你的IP`，端口填 `3307`，但 `DB_PORT` 仍然保持 `3306`

> ⚠️ **关于 APP_URL**：
> - 这是 Laravel 生成绝对 URL 用的（邮件链接、API 返回的 URL 等），必须填**其他设备访问你时用的地址**
> - **不要**填 Docker 容器内部 IP（如 `172.17.0.x`），局域网其他设备无法访问该地址
> - 群晖/飞牛 NAS 部署：填 NAS 的局域网 IP + APP_PORT，如 `http://192.168.1.100:8080`

**可选项**（让 AI 助教可用）：
```env
AI_PROVIDER=qwen
AI_API_KEY=your_qwen_api_key
```

#### 3. 启动服务

```bash
# 启动（包含应用+数据库+Redis，会自动运行迁移并创建管理员账号）
docker-compose up -d

# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f app
```

> 💡 **默认管理员**：`admin` / `admin123456`，可在 `.env` 中通过 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD` 修改。

#### 4. 初始化数据库

```bash
# 运行数据库迁移（启动脚本已自动执行，也可手动执行）
docker-compose exec app php artisan migrate --force

# 创建默认管理员账号（默认：admin / admin123456）
docker-compose exec app php artisan db:seed --class=AdminUserSeeder --force
```

> 💡 **提示**：管理员账号可在 `.env` 中通过 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD` 自定义。如果账号已存在，重复执行不会修改密码。

#### 5. 访问应用

打开浏览器，访问：`http://your-server-ip`

---

### 🗄️ 方式三：使用外部数据库（适合已有数据库的用户）

如果你已经有外部 MySQL/MariaDB/PostgreSQL 或 Redis，可以跳过内置容器：

#### 1. 修改 `.env`

```env
# 修改为你的外部数据库地址
DB_HOST=your_mysql_host
DB_PORT=3306
DB_DATABASE=learnstar
DB_USERNAME=learnstar_user
DB_PASSWORD=your_password

# 修改为你的外部 Redis 地址
REDIS_HOST=your_redis_host
REDIS_PORT=6379
REDIS_PASSWORD=
```

#### 2. 启动服务（只启动应用，跳过内置 MySQL/Redis）

```bash
# --no-deps 跳过内置 mysql/redis 依赖，直接启动 app
docker-compose up -d app --no-deps
```

---

### 🐳 方式四：本地构建镜像（适合开发者）

```bash
# 1. Fork 本仓库
# 2. 克隆到本地
git clone https://github.com/YOUR_USERNAME/learnstar-planet.git
cd learnstar-planet

# 3. 构建镜像
docker build -t learnstar-planet:latest ./backend

# 4. 修改 docker-compose.yml 中的 image 为本地镜像
# image: learnstar-planet:latest

# 5. 启动
docker-compose up -d
```

---

## 配置说明

### 🤖 AI 功能配置（可选）

AI 助教功能需要配置 API Key，支持以下平台：

| 平台 | 获取 API Key | 推荐模型 | 价格 |
|------|-------------|----------|------|
| **通义千问**（推荐） | [阿里云百炼](https://dashscope.aliyuncs.com/) | `qwen-turbo`（¥0.002/1k tokens） | 便宜 |
| **OpenAI** | [OpenAI Platform](https://platform.openai.com/) | `gpt-4o-mini` | 需翻墙 |
| **DeepSeek** | [DeepSeek 平台](https://platform.deepseek.com/) | `deepseek-chat` | 便宜 |
| **月之暗面** | [Moonshot 平台](https://platform.moonshot.cn/) | `moonshot-v1-8k` | 中等 |

**配置步骤**（以通义千问为例）：

1. 访问 [阿里云百炼控制台](https://dashscope.aliyuncs.com/)
2. 注册/登录（需实名认证）
3. 点击 **API Key 管理** → **创建 API Key**
4. 复制 API Key，修改 `.env`：
   ```env
   AI_PROVIDER=qwen
   AI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxx
   AI_API_BASE=https://dashscope.aliyuncs.com/api/v1
   AI_MODEL=qwen-turbo
   ```
5. 重启服务：`docker-compose restart app`

**没有 API Key？** 没关系！所有核心功能都可以正常使用，AI 功能会在配置后自动启用。

---

### 🔐 第三方登录配置（可选）

#### 个人微信扫码登录

1. 访问 [微信开放平台](https://open.weixin.qq.com/)
2. 注册开发者账号（需企业资质并通过审核）
3. 创建**网站应用**，获取 **AppID** 和 **AppSecret**
4. 配置**授权回调域名**（填写你的域名）
5. 修改 `.env`：
   ```env
   WECHAT_OPEN_APPID=wx1234567890abcdef
   WECHAT_OPEN_SECRET=xxxxxxxxxxxxxxxxxxxxxx
   ```

#### 企业微信扫码登录

1. 访问 [企业微信管理后台](https://work.weixin.qq.com/)
2. 进入 **应用管理** → **创建应用**
3. 获取 **CorpID**、**AgentID**、**Secret**
4. 配置**网页授权及JS-SDK**中的可信域名
5. 修改 `.env`：
   ```env
   WECHAT_WORK_CORPID=ww1234567890abcdef
   WECHAT_WORK_AGENTID=1000002
   WECHAT_WORK_SECRET=xxxxxxxxxxxxxxxxxxxxxx
   ```

#### QQ 扫码登录

1. 访问 [QQ 互联平台](https://connect.qq.com/)
2. 注册开发者账号并创建**网站应用**
3. 获取 **App ID** 和 **App Key**
4. 配置**回调地址**（填写 `https://你的域名/auth/qq/callback`）
5. 修改 `.env`：
   ```env
   QQ_APP_ID=1012345678
   QQ_APP_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   QQ_REDIRECT_URI=https://your-domain.com/auth/qq/callback
   ```

#### 人人通空间登录

人人通空间支持两种登录方式：**账号密码登录**与**客户端扫码登录**。

1. 联系当地教育主管部门或人人通空间服务方，申请接入权限
2. 获取 **应用 ID（AppKey）** 和 **应用密钥（AppSecret）**
3. 配置**回调地址**（填写 `https://你的域名/auth/rrt/callback`）
4. 修改 `.env`：
   ```env
   RRT_APP_KEY=rrt1234567890abcdef
   RRT_APP_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   RRT_REDIRECT_URI=https://your-domain.com/auth/rrt/callback
   ```

---

### 🗄️ 数据库选择

| 数据库 | 适用场景 | 配置方法 |
|--------|----------|----------|
| **MySQL 8.0+**（推荐） | 中大型部署（1000+ 学生） | `DB_CONNECTION=mysql` |
| **MariaDB 10.3+** | MySQL 替代品，完全兼容 | `DB_CONNECTION=mysql` |
| **PostgreSQL 14+** | 需要高级功能（JSONB、全文搜索） | `DB_CONNECTION=pgsql` |
| **SQLite 3.8+** | 小型部署（< 500 学生） | `DB_CONNECTION=sqlite` |

#### 使用 SQLite（最简单）

```bash
# 1. 修改 .env
DB_CONNECTION=sqlite
# 同时将缓存改为 file（SQLite 部署通常无 Redis）
CACHE_DRIVER=file
SESSION_DRIVER=file
QUEUE_CONNECTION=database

# 2. 启动（跳过内置 MySQL/Redis）
docker-compose up -d app --no-deps

# 3. 运行迁移
docker-compose exec app php artisan migrate --force
```

---

## FAQ

### Q1：我没有服务器，能部署吗？

**A**：可以！以下免费/低成本方案供参考：

| 平台 | 类型 | 价格 | 适用场景 |
|------|------|------|----------|
| **Fly.io** | 容器托管 | 免费额度 3 个小型 VM | 小型部署 |
| **Railway** | PaaS | $5/月 免费额度 | 快速部署 |
| **Render** | PaaS | 免费版（休眠） | 测试环境 |
| **阿里云轻量** | VPS | ¥60/月 | 生产环境（推荐） |
| **腾讯云轻量** | VPS | ¥50/月 | 生产环境（推荐） |

---

### Q2：AI 助教功能一定要配置吗？

**A**：不是必须的。所有核心功能（积分、宠物、考勤、广播等）都可以独立使用。AI 功能只是增强体验，会在配置 API Key 后自动启用。

---

### Q3：第三方登录一定要配置吗？

**A**：不是必须的。你可以使用账号密码登录。第三方登录只是提供更便捷的登录方式。

---

### Q4：数据库和 Redis 一定要单独部署吗？

**A**：不是必须的。Docker Compose 会启动数据库和 Redis 容器，适合小型部署。当学生数超过 1000 人时，建议将它们单独部署到更高性能的服务器。

---

### Q5：如何升级到新版本？

**A**：

```bash
# 1. 拉取最新镜像
docker-compose pull

# 2. 重新启动
docker-compose up -d

# 3. 运行数据库迁移（如果有）
docker-compose exec app php artisan migrate --force
```

---

### Q6：忘记管理员密码怎么办？

**A**：

```bash
# 通过命令行重置
docker-compose exec app php artisan admin:reset-password --username=admin
# 会提示输入新密码
```

---

### Q7：如何备份数据？

**A**：

```bash
# 备份 MySQL 数据库
docker-compose exec mysql mysqldump -u root -p learnstar > backup_$(date +%Y%m%d).sql

# 备份上传的文件
docker-compose exec app tar -czf /tmp/uploads_backup.tar.gz /var/www/html/storage/app/uploads
docker cp learnstar-app:/tmp/uploads_backup.tar.gz ./
```

---

### Q8：出现 500 错误怎么办？

**A**：

```bash
# 1. 查看日志
docker-compose logs app

# 2. 查看 Laravel 日志
docker-compose exec app tail -100 storage/logs/laravel.log

# 3. 常见原因：
#    - .env 配置错误（检查数据库/Redis 连接）
#    - 数据库迁移未运行
#    - 文件权限问题（运行 chmod -R 777 storage）
```

---

### Q9：可以商用吗？需要付费吗？

**A**：本项目采用 **MIT 许可证**，你可以自由使用、修改、商用，无需支付任何费用。但请保留原作者版权声明。

---

## 贡献指南

我们欢迎任何形式的贡献！

### 🐛 报告 Bug

请在 [GitHub Issues](https://github.com/RealKiro/learnstar-planet/issues) 中提交 Bug 报告，包含以下信息：
- 问题描述
- 复现步骤
- 预期行为
- 实际行为
- 截图（如有）

### 💡 功能建议

请在 [GitHub Discussions](https://github.com/RealKiro/learnstar-planet/discussions) 中发起功能讨论。

### 🔧 提交代码

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

### 📝 代码规范

- PHP：遵循 PSR-12 规范
- JavaScript：使用 ESLint 检查
- 提交信息：使用约定式提交（Conventional Commits）

---

## 许可证

本项目采用 **MIT 许可证** - 查看 [LICENSE](LICENSE) 文件了解详情。

---

## 联系我们

- 📖 **文档**：[查看完整文档](https://github.com/RealKiro/learnstar-planet/wiki)
- 🐛 **Bug 反馈**：[提交 Issue](https://github.com/RealKiro/learnst