<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  <img src="https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/badge/PHP-8.3-777BB4?logo=php" alt="PHP 8.3">
  <img src="https://img.shields.io/badge/Laravel-11-F9322C?logo=laravel" alt="Laravel 11">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" alt="Docker">
</p>

<h1 align="center">学趣星球 / LearnStar Planet</h1>

<p align="center">
  开源、自托管的班级管理与学生激励系统。<br>
  积分评价 · 宠物进化 · 课堂互动 · 数据报表 · AI 助教
</p>

<p align="center">
  <a href="#-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-deployment">Deployment</a> •
  <a href="#-configuration">Configuration</a> •
  <a href="#-tech-stack">Tech Stack</a> •
  <a href="#-faq">FAQ</a>
</p>

---

## 📦 Features

- **积分评价系统** — 自定义加分/扣分规则，支持单人操作、批量操作、全班一键执行，每次变动实时推送至教室大屏
- **宠物进化系统** — 54 种宠物、8 大系列、12 级进化路线，积分即经验值，自动累计升级，每级独立名称和诗文
- **教室实时广播** — 横幅、弹窗、全屏三种模式，教师操作后 <200ms 到达教室大屏，支持多班级同时发送
- **跨班 PK 战场** — 同年级各班自动排行（总积分 / 平均等级 / 巅峰人数 / 周增长），支持发起班级挑战
- **积分商城** — 教师自定义奖品，学生自助兑换，完整的审核/发货/拒绝流程，支持多币种
- **排行榜** — 总积分榜、周增长榜、宠物等级榜，Redis ZSET 毫秒级排序，无 Redis 自动降级 MySQL
- **智能考勤** — 一键发起点名，出勤 / 迟到 / 请假 / 缺席自动汇总
- **在线成绩管理** — 按考试和科目录入成绩，统计平均分 / 最高分 / 最低分 / 分数段分布
- **数据报表** — 积分趋势（近 4 周）、宠物等级分布、学生进步追踪、分年级/班级统计
- **实时通知** — 班级通知发布后同步推送至家长端
- **学年升级** — 预览升级明细 → 事务性执行，自动处理毕业与班级迭代
- **学生/教师管理** — Excel 批量导入、批量转班/删除、智能去重账号创建
- **多数据库支持** — SQLite / MySQL / PostgreSQL / MariaDB 任选，从小规模起步可无缝升级
- **AI 中心** — 管理后台统一配置 AI，支持 30+ 供应商（OpenAI / Claude / Gemini / DeepSeek / 通义千问 / Kimi / 豆包 / MiniMax / 百川 / GLM / 星火等），MCP 通用接口可对接任意 OpenAI 兼容服务，班级码大屏可开启 AI 对话
- **实时广播** — 横幅、弹窗、全屏三种模式，教师操作后 <200ms 到达教室大屏
- **消息中心** — 实时广播 + 班级通知合并为统一入口
- **系统诊断与修复** — 一键检查数据库结构缺失并自动修复
- **教师端侧边栏重组** — 按课堂教学 / 成长激励 / 数据中心 / 沟通协作 / 系统管理分组
- **班级切换器** — 多班教师可在侧边栏一键切换当前班级，所有数据联动
- **演示数据管理** — 一键生成/清除演示教师账号
- **Excel 导出** — 积分报表、宠物报表、考勤报表一键导出 .xlsx
- **微信小程序** — 教师端看板、积分管理、宠物花园、排行榜及家长端全功能
- **MCP 机器人协议** — 标准 MCP 服务器，对接 QQ/微信群聊机器人实现自然语言查分、排行榜
- **家长端** — 手机浏览器或微信小程序查看孩子积分、宠物、排名、通知

## 🚀 Quick Start

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

启动后：
- **教室互动**：浏览器打开 `http://<服务器IP>:8080`，输入班级码即可进入
- **管理后台**：`http://<服务器IP>:8080/login`，使用管理员或教师账号登录

> 整个系统部署在学校局域网即可运行，无需互联网连接。手机连同一 WiFi 也可访问。

## 📋 Requirements

| 部署方式 | 前置依赖 |
|---------|---------|
| Docker | Docker Engine + Docker Compose |
| SQLite 部署 | PHP 8.3 + Composer + Node.js 22 |
| LAMP 部署 | Nginx + PHP 8.3-FPM + MySQL |

最低硬件：1 核 CPU、512MB 内存、5GB 磁盘（办公室淘汰 PC 即可）。

## 🛠 Deployment

### Docker（推荐）

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

内置 MySQL 8.0 + Redis 7，开箱即用。如需使用外部数据库，修改 `.env` 中 `DB_HOST` 后执行 `docker-compose up -d app --no-deps`。

### SQLite（零依赖部署）

<details>
<summary>无需 Docker、无需 MySQL、无需 Redis，适合 <500 学生</summary>

```bash
cd backend
composer install --no-dev --optimize-autoloader
cp .env.example .env
```

编辑 `.env`：

```env
DB_CONNECTION=sqlite
CACHE_DRIVER=file
SESSION_DRIVER=file
QUEUE_CONNECTION=database
BROADCAST_DRIVER=null
REDIS_HOST=
```

```bash
touch storage/database.sqlite
php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder

cd ../frontend-vue
npm install && npm run build
cp -r dist/* ../backend/public/

cd ../backend
php artisan serve --host=0.0.0.0 --port=8080
```

</details>

### 数据库方案选型

| | SQLite | MySQL + Redis |
|---|---|---|
| 安装 | 无需安装 | 需单独安装 |
| 推荐规模 | < 500 学生 | 无限制 |
| 排行榜 | SQL 查询 | Redis ZSET 毫秒级 |
| 实时广播 | 不支持 | 支持 |
| 适用场景 | 单间教室 / 单个年级 | 全校规模 |

> 先用 SQLite 跑起来，后续改两行配置即可切换到 MySQL，代码无需改动。

<details>
<summary>LAMP 部署方式</summary>

```bash
cd backend
composer install --no-dev --optimize-autoloader
cp .env.example .env
# 编辑 .env 配置数据库连接
php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder

cd ../frontend-vue
npm install && npm run build
# 将 dist/ 部署至 Nginx 根目录或 backend/public/
```

Nginx 参考配置：

```nginx
server {
    listen 80;
    server_name your-domain.local;
    root /var/www/learnstar-planet/backend/public;
    index index.php;
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php8.3-fpm.sock;
        include fastcgi_params;
    }
}
```

</details>

## ⚙️ Configuration

### 核心环境变量（`.env`）

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_CONNECTION` | `mysql` | 数据库类型：`mysql` / `pgsql` / `sqlite` |
| `CACHE_DRIVER` | `redis` | 缓存驱动：`redis` / `file` |
| `QUEUE_CONNECTION` | `redis` | 队列驱动：`redis` / `database` |
| `AI_PROVIDER` | — | AI 服务商：`deepseek` / `openai` / `qwen` / `moonshot` |
| `AI_API_KEY` | — | AI API Key |
| `APP_URL` | `http://localhost` | 系统访问地址，用于生成链接 |

## 🏗 Tech Stack

| 层级 | 技术 |
|------|------|
| 后端 | Laravel 11（PHP 8.3） |
| 数据库 | SQLite / MySQL / MariaDB / PostgreSQL |
| 缓存与队列 | Redis 7（可降级 file / database） |
| 前端 | Vue 3 + TypeScript + Tailwind CSS |
| 实时推送 | SSE 协议（后端实现，前端 EventSource 优先 / 轮询降级） |
| 小程序 | 微信原生（14 页面，教师端 + 家长端） |
| 机器人 | MCP 协议 Python 服务器 |
| 部署 | Docker 多阶段构建 + Compose 编排 / 裸 PHP 环境 |

## 🗂 Project Structure

```
learnstar-planet/
├── frontend-vue/          # Vue 3 SPA（教室互动 + 管理后台 + 家长端）
├── backend/               # Laravel 11 API（22 模型，130+ 端点）
├── mini-program/          # 微信小程序（14 页面）
├── pwa/                   # PWA 离线配置
├── mcp-server/            # MCP 机器人服务
└── docker-compose.yml     # Docker 编排
```

## 🧪 安全调试模式

## 📋 TODO

| 优先级 | 功能 | 备注 |
|:------:|------|------|
| 🔵 | 宠物系列扩展 | 追加机甲、星座、魔法 3 个系列（含 36 个物种 × 11 阶进化数据） |


## ❓ FAQ

管理员后台一键生成演示数据（不影响真实数据）：

| 项目 | 内容 |
|------|------|
| 管理员账号 | `demo_admin` |
| 教师账号 | `demo_t1` ~ `demo_t4` |
| 统一密码 | `demo123` |
| 班级码 | `DEMO00`（通用演示码）|
| 班级码 | 生成后在管理后台查看 |
| 清除数据 | 后台一键清除，不留痕迹 |

生成路径：管理员登录 → 学校设置 → 演示数据管理


<details>
<summary>学生需要注册账号吗？</summary>
不需要。教师在后台为每个班级生成班级码，学生在首页输入班级码即可进入，无需账号密码。
</details>

<details>
<summary>需要联网吗？</summary>
不需要。系统部署在校园局域网内即可运行，教室电脑与教师手机连同一个网络即可访问。
</details>

<details>
<summary>对硬件有什么要求？</summary>
最低 1 核 CPU、512MB 内存。办公室淘汰的 PC 安装 Linux + Docker 即可运行。全校规模建议 2 核 4GB 以上。
</details>

<details>
<summary>AI 助教怎么用？</summary>
在 `.env` 中配置 `AI_PROVIDER` 和 `AI_API_KEY`，支持 DeepSeek / OpenAI / 通义千问 / Moonshot。不配置不影响其他功能。
</details>

<details>
<summary>家长怎么看孩子数据？</summary>
家长通过手机浏览器或微信小程序访问部署地址，使用家长账号登录即可查看。
</details>

<details>
<summary>有用户数限制吗？</summary>
无。MIT 开源协议不限制用户数和班级数。全校使用建议 MySQL + Redis 方案。
</details>

## 🔗 企业微信集成

部分高级功能需要管理员配置企业微信（WeChat Work）后方可使用：

| 功能 | 说明 | 是否必须配置 |
|------|------|:-----------:|
| **教师免注册登录** | 教师扫码企业微信后自动登录，未绑定账号时根据企业微信通讯录信息自动创建账号，无需管理员后台手动分配 | 可选 |
| **学生请假自动同步** | 学生在企业微信提交请假申请后，审批通过自动同步到学趣星球考勤系统，教师端直接看到请假记录及来源标记 | 可选 |
| **企业微信消息推送** | 通知、广播等可通过企业微信应用推送到教师/家长手机 | 可选 |

配置方式：在 `.env` 中填写 `WECHAT_WORK_CORPID`、`WECHAT_WORK_AGENTID`、`WECHAT_WORK_SECRET`。不配置不影响积分、宠物、考勤等核心功能。

## 🤝 Contributing

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/my-feature`
3. 提交变更：`git commit -am 'feat: add my feature'`
4. 推送：`git push origin feature/my-feature`
5. 提交 Pull Request

代码规范：PHP PSR-12（PHP CS Fixer）· TypeScript ESLint · Conventional Commits

## 📄 License

[MIT License](LICENSE) © 2024 RealKiro
