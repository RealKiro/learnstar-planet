<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  <img src="https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/badge/PHP-8.3-777BB4?logo=php" alt="PHP 8.3">
  <img src="https://img.shields.io/badge/Laravel-11-F9322C?logo=laravel" alt="Laravel 11">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" alt="Docker">
</p>

<h1 align="center">🌌 学趣星球 (LearnStar Planet)</h1>

<p align="center">
  <b>开源班级管理平台 · 积分激励 · 宠物养成 · 课堂互动</b><br>
  局域网部署 · 零经费 · 学生免注册 · 数据在校内
</p>

<p align="center">
  <i>面向中小学教师的自托管课堂管理系统，一台普通电脑即可运行。</i>
</p>

```bash
# Docker 部署（一分钟上线）
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet && cp .env.example .env
docker-compose up -d
```

---

## 📖 目录

- [系统角色与功能清单](#-系统角色与功能清单)
- [为什么适合学校用](#-为什么适合学校用)
- [部署方式](#-部署方式)
- [常见问题](#-常见问题)

---

## 🎯 系统角色与功能清单

系统分三个角色端，其中教师端有两种进入方式，权限不同。

### 教师端——输入班级码（免登录）

教室电脑打开浏览器，输入班级码即可看到以下界面，适合投影仪/电视展示：

| 功能 | 说明 |
|------|------|
| **班级总览** | 显示当前积分之星、TOP5 排名、全班最新积分动态滚动 |
| **课堂评价** | 点击学生头像进行加减分，选择行为原因（答对题、作业优秀、违反纪律等），全班实时可见。加分宠物弹跳、扣分抖动 |
| **年级战场** | 同年级各班积分排行，支持发起挑战 |
| **宠物图鉴** | 全班宠物展示，54种、8大系列、12级进化路线 |

> 此模式不需要账号密码，适合上课时直接投屏使用。

### 教师端——登录教师账号（完整管理）

输入教师账号密码后，可进行全部课堂管理操作：

| 功能 | 说明 |
|------|------|
| **实时广播** | 向教室大屏发送横幅/弹窗/全屏三种模式通知，＜200ms 到达 |
| **积分规则管理** | 自定义加减分规则（如"答对一题 +2 分"），支持批量操作和全班一键执行 |
| **宠物系统管理** | 查看全班宠物状态、喂食、改名、切换系列 |
| **排行榜** | 总积分榜、周增长榜、宠物等级榜（Redis 加速，无 Redis 自动降级 MySQL） |
| **智能考勤** | 一键点名，标记出勤/迟到/请假/缺席，自动汇总 |
| **积分商城** | 自定义奖品，学生自助兑换，教师审核/发货 |
| **班级通知** | 发布通知，家长端同步可见 |
| **成绩管理** | 按考试/科目录入成绩，统计平均分/最高分/最低分/分布 |
| **数据报表** | 积分趋势、宠物分布、学生进步追踪 |

### 管理员端

| 功能 | 说明 |
|------|------|
| **学校/班级管理** | 创建班级、生成班级码、分配班主任 |
| **学生管理** | 单个/批量创建、Excel 导入、批量转班/删除 |
| **教师/家长管理** | 批量创建账号（智能去重）、分配权限、重置密码 |
| **学年升级** | 预览 → 确认 → 事务执行，自动处理毕业与年级迭代 |
| **数据报表** | 分年级/分班级统计 |
| **第三方登录配置** ⚠️ | 后端微信/QQ/企业微信/人人接口已实现，前端 OAuth 跳转未对接 |

### 家长端

| 功能 | 说明 |
|------|------|
| **查看孩子数据** | 当前积分、积分明细、宠物状态与喂食、班级排名、学校通知 |

手机浏览器访问即可，无需安装 App。

### 扩展集成

| 功能 | 说明 |
|------|------|
| **MCP 机器人** ✅ | 标准 MCP 协议 Python 服务器，7 个工具。对接 AstrBot 实现 QQ/微信群聊中查积分、排行榜 |

---

## 📋 TODO（未达到可用状态的功能）

以下功能或缺少关键环节，或仅有前端/后端框架，需要进一步开发方可使用。

| 功能 | 现状 | 待完成工作 |
|------|------|-----------|
| **微信小程序** | 已注册 14 个页面，目前仅 3 个核心页面（首页/登录/教师看板）可运行 | 补充其余 11 个页面内容 |
| **第三方登录（前端）** | 后端微信/企业微信/QQ/人人 OAuth 接口已实现 | 前端扫码按钮需对接真实 OAuth 跳转流程 |
| **Excel 导出** | Excel 导入已实现 | 实现导出接口（积分报表、成绩报表等） |
| **SSE 实时推送（前端）** | 后端 SSE 实现完整（心跳/超时/断线重连） | 前端改用 EventSource 代替轮询 |
| **AI 助教** | 前端聊天界面和环境变量已预留 | 后端对接 AI API（通义千问/DeepSeek/OpenAI）实现实际对话 |

---

## 🏫 为什么适合学校用

| 考量 | 说明 |
|------|------|
| **零采购成本** | MIT 协议开源，无隐藏收费、无用户数限制、无功能锁定 |
| **局域网即可运行** | 全部部署于校内网络，无需互联网出口 |
| **硬件门槛低** | 办公室淘汰的 PC 即可作为服务器 |
| **学生零注册** | 教师生成 4-6 位班级码，学生输入即进入，无需账号密码 |
| **数据不出校** | 所有数据存储在校内服务器，不经第三方平台 |
| **多数据库适配** | SQLite / MySQL / PostgreSQL 任选，可从 SQLite 起步无缝升级 |

---

## 🚀 部署方式

### 部署方案对比

| 方式 | 前置依赖 | 适用规模 | 操作难度 |
|------|---------|---------|:--------:|
| Docker 部署 | Docker 环境 | 任意规模 | ⭐ 低 |
| SQLite 部署 | PHP 8.3 + Composer | < 500 学生 | ⭐⭐ 中 |
| LAMP 部署 | Nginx + PHP + MySQL | 全校规模 | ⭐⭐ 中 |

### 方式一：Docker 部署（推荐）

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

启动后访问：
- **教室互动（免登录）**：`http://<服务器IP>:8080` → 输入班级码
- **教师/管理员登录**：`http://<服务器IP>:8080/login` → 输入账号密码

> 局域网内所有设备均可通过服务器 IP 访问，手机连同一 WiFi 也能操作。

### 方式二：SQLite 部署（不用 Docker、不用 MySQL、不用 Redis）

<details>
<summary><b>📋 详细步骤</b></summary>

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

### 数据库方案对比

| 维度 | SQLite | MySQL + Redis |
|------|--------|---------------|
| 安装配置 | 无需安装 | 需额外安装维护 |
| 推荐规模 | < 500 学生 | 无限制 |
| 排行榜性能 | 常规 SQL 查询 | Redis ZSET（毫秒级） |
| 实时广播 | 不支持 | 支持 |
| 适用场景 | 单间教室 / 单个年级 | 全校规模 |

> 先用 SQLite 跑起来，后续改配置即可切换到 MySQL，代码无需改动。

<details>
<summary><b>📋 LAMP 部署（点击展开）</b></summary>

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

---

<details>
<summary><b>🏗 系统架构（点击展开）</b></summary>

```
learnstar-planet/
├── frontend-vue/          # Vue 3 + TypeScript SPA（三个角色端共用）
│   └── pages/
│       ├── teacher/       # 教师端（班级码模式 + 账号模式）
│       ├── admin/         # 管理员端
│       └── parent/        # 家长端
├── backend/               # Laravel 11 API（22 模型，~130 端点）
├── mini-program/          # 微信小程序（14 页面注册，3 个已实现）
├── pwa/                   # PWA 离线配置
├── mcp-server/            # MCP 机器人服务（Python）
└── docker-compose.yml     # Docker 编排
```

### 技术栈

| 层级 | 技术 |
|------|------|
| 后端框架 | Laravel 11（PHP 8.3） |
| 数据库 | SQLite / MySQL / MariaDB / PostgreSQL |
| 缓存与队列 | Redis 7（可降级为 file/database） |
| 前端 | Vue 3 + TypeScript + Tailwind CSS |
| 实时推送 | SSE 协议（后端已实现，前端使用轮询） |
| 机器人 | MCP 协议（Python） |
| 部署 | Docker 多阶段构建 + Compose 编排 |

</details>

---

## ❓ 常见问题

<details>
<summary><b>没有专用服务器，用办公室旧电脑行不行？</b></summary>
可以。Windows 装 PHP 8.3 和 Composer，用 SQLite 模式运行。建议由学校信息技术教师协助部署。
</details>

<details>
<summary><b>学生需要注册账号吗？</b></summary>
不需要。教师后台为每个班生成班级码（4-6 位），学生在首页输入班级码即可进入，无需任何注册流程。
</details>

<details>
<summary><b>教室没有外网能用吗？</b></summary>
能。全部部署在局域网内，教室电脑和教师设备连同一个校园网即可使用，完全不需要互联网。
</details>

<details>
<summary><b>教师怎么用手机操作？</b></summary>
教师用手机浏览器打开管理页面，登录教师账号即可加分、发广播、看排行，教室大屏实时同步。
</details>

<details>
<summary><b>广播功能怎么用？</b></summary>
教师登录账号 → 实时广播 → 选横幅/弹窗/全屏 → 输入内容发送，教室大屏即时弹出。
</details>

<details>
<summary><b>家长怎么看孩子表现？</b></summary>
家长用手机浏览器打开学校访问地址，用家长账号登录即可查看孩子积分、宠物、排名、通知。
</details>

<details>
<summary><b>有用户数/班级数限制吗？</b></summary>
无。MIT 协议不限制用户数和班级数。全校使用建议 MySQL + Redis 方案。
</details>

<details>
<summary><b>和商业班宠系统的区别？</b></summary>
完全免费开源，数据在校内，不经过第三方。商业系统年费数千至数万元。功能覆盖度相近，界面精致度有差距。
</details>

---

## 📄 许可证

<p align="center">
  <a href="LICENSE">MIT License</a> © 2024 RealKiro<br>
  自由使用、修改和再发布。
</p>
