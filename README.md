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
  <b>开源班级管理平台 · 教室大屏互动 · 积分激励 · 宠物养成</b><br>
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

- [功能清单](#-功能清单哪些真能用哪些还在完善)
- [为什么适合学校用](#-为什么适合学校用)
- [部署方式](#-部署方式)
- [常见问题](#-常见问题)

---

## 🎯 功能清单——哪些真能用，哪些还在完善

### 实时广播 ✅
教师可向教室大屏发送三种模式的实时通知：顶部横幅（不打断课堂）、居中弹窗（紧急通知）、全屏覆盖（专注展示）。适合上课铃响提醒、自习课通告、课间活动通知。前后端均已实现。

### 教室大屏 ✅
浏览器全屏打开即进入教室模式，显示班级之星、实时动态滚动、TOP5排名进度条。加分时宠物弹跳、扣分时抖动，学生获得即时视觉反馈。前后端均已实现。

### 课堂加减分 ✅
以学生网格卡片为交互单元，点击头像选择行为原因（答对问题、作业优秀、违反纪律等），支持自定义分值步长和积分规则。可一键按规则操作全班。前后端均已实现。

### 宠物系统 ✅
54种宠物，分属8大系列（神话、宝可梦、机甲、国宝、魔法、史前、星座、民间），每级进化有独立名称、阶段描述和专属诗文。积分即宠物经验值，累计升级。本项目最完善的功能模块。前后端均已实现。

### 跨班PK ✅
同年级各班积分排行榜，展示总积分、平均等级、巅峰人数、周增长四项指标，支持发起班级间PK挑战。前后端均已实现。

### 积分商城 ✅
教师自定义奖品（减负卡、免作业、实物等），学生自助兑换，教师审核/发货/拒绝。支持多币种。前后端均已实现。

### 排行榜 ✅
总积分榜、周增长榜、宠物等级榜三种排行。优先使用Redis ZSET实现毫秒级排序，Redis不可用时自动降级为MySQL查询。前后端均已实现。

### 智能考勤 ✅
一键发起点名，学生标记为出勤/迟到/请假/缺席，自动汇总到班级台账。前后端均已实现。

### 班级通知 ✅
向班级或全校发布通知，家长端实时同步可见。前后端均已实现。

### 作业管理 ⚠️
布置作业、生成二维码token、查看提交列表——核心管理功能已实现。但**二维码扫码提交的入口路由尚未注册**，学生端暂不能实际扫码交作业。

### 在线测验 ⚠️
教师端题库管理（单选/多选/填空/判断/简答）、组卷、发布、自动判分已实现。但**学生端作答流程尚未完成**，仍需完善。

### 成绩管理 ✅
按考试名称和科目录入成绩（支持批量）、统计平均分/最高分/最低分、查看分数段分布。前后端均已实现。

### 学年升级 ✅
预览升级明细（哪些班升到哪个年级、哪些学生毕业）→ 确认后以数据库事务执行，自动处理班级迭代和毕业生状态。前后端均已实现。

### 学生管理 ✅
单个创建、Excel批量导入（自动匹配班级）、批量转班、批量删除、按年级/班级筛选。前后端均已实现。

### 教师/家长管理 ✅
批量创建账号（智能去重用户名和昵称）、分配班级权限、重置密码。前后端均已实现。

### 数据报表 ✅
积分趋势（近4周走势）、宠物等级分布、学生进步追踪、分年级/分班级统计报表。前后端均已实现。

### 家长端 ✅
家长通过手机浏览器访问，可查看绑定孩子的当前积分、积分明细、宠物状态与喂食、班级排名、学校通知。前后端均已实现。

### 第三方登录 ⚠️
后端已完成微信、企业微信、QQ、人人四种OAuth认证与绑定接口。但**前端扫码按钮目前仅弹出提示信息，尚未对接真实OAuth跳转流程**。账号密码登录完全正常。

### 微信小程序 ⚠️
已注册14个页面，目前完成3个核心页面（首页、登录、教师看板），其余页面待建设。

### SSE 实时推送 ⚠️
后端实现完整的Server-Sent Events协议（<200ms推送、每10秒心跳检测、连接超时管理、断线自动重连协议）。但**前端当前使用HTTP轮询消费事件**，未接入EventSource客户端。

### MCP 机器人 ✅
基于标准MCP协议的Python服务器，提供7个工具（加分、批量加分、查分、排行榜、通知、看板、搜索学生）。可对接AstrBot + NapCatQQ/Lagrange 实现QQ/微信群聊中自然语言查询积分。API请求失败时自动降级为演示数据。

### Excel 导出 ❌
Excel导入功能已实现，但导出接口目前返回"暂不支持"。建议通过数据库管理工具自行导出。

### AI 助教 ❌
前端有聊天界面，环境变量也预留了AI相关配置项，但**后端实际未对接任何AI API**。需要先实现AI调用逻辑方可使用。

---

## 🏫 为什么适合学校用

| 考量 | 说明 |
|------|------|
| **零采购成本** | MIT协议开源，无隐藏收费、无用户数限制、无功能锁定 |
| **局域网即可运行** | 全部部署于校内网络，教室电脑和教师设备接入同一校园网即可访问，无需互联网出口 |
| **硬件门槛低** | 办公室淘汰的PC即可作为服务器，支持Docker或裸PHP环境 |
| **学生零注册** | 教师后台生成班级码（4-6位），学生首页输入即进入，无需账号密码 |
| **数据不出校** | 所有学生数据存储在校内服务器，不经过任何第三方平台 |
| **多数据库适配** | SQLite / MySQL / PostgreSQL 任选，从小规模到全校使用均可扩展 |

---

## 🚀 部署方式

### 部署方案对比

| 方式 | 前置依赖 | 适用规模 | 操作复杂度 |
|------|---------|---------|:----------:|
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
- **教室大屏**：`http://<服务器IP>:8080` → 输入班级码
- **教师管理**：`http://<服务器IP>:8080/login` → 管理员/教师登录

> 局域网内所有设备均可通过服务器IP访问。手机连接同一WiFi后同样可操作。

### 方式二：SQLite 部署（零数据库依赖）

<details>
<summary><b>📋 详细步骤</b></summary>

```bash
# 安装后端
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
# 初始化
touch storage/database.sqlite
php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder

# 构建前端
cd ../frontend-vue
npm install && npm run build
cp -r dist/* ../backend/public/

# 启动服务
cd ../backend
php artisan serve --host=0.0.0.0 --port=8080
```

</details>

### 数据库方案对比

| 维度 | SQLite | MySQL + Redis |
|------|--------|---------------|
| 安装配置 | 无需安装 | 需额外安装维护 |
| 推荐规模 | < 500 学生 | 无限制 |
| 排行榜性能 | 常规SQL查询 | Redis ZSET（毫秒级） |
| 实时广播 | 不支持 | 支持 |
| 应用场景 | 单间教室 / 单个年级 | 全校规模 |

> 从小规模起步，后续可直接修改配置切换到 MySQL，代码无需改动。

<details>
<summary><b>📋 LAMP 部署方式（点击展开）</b></summary>

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
├── frontend-vue/          # Vue 3 + TypeScript SPA
│   └── pages/
│       ├── teacher/       # 教师管理（17页面）
│       ├── classroom/     # 教室大屏（4页面）
│       ├── parent/        # 家长端（6页面）
│       └── admin/         # 学校管理（9页面）
├── backend/               # Laravel 11 API（22模型，~130端点）
├── mini-program/          # 微信小程序（14页面注册，3个已实现）
├── pwa/                   # PWA离线配置
├── mcp-server/            # MCP 机器人服务（Python）
└── docker-compose.yml     # Docker 编排
```

### 技术栈

| 层级 | 技术 |
|------|------|
| 后端框架 | Laravel 11（PHP 8.3） |
| 数据库 | SQLite / MySQL / MariaDB / PostgreSQL |
| 缓存与队列 | Redis 7（可降级为 file/database 驱动） |
| 前端框架 | Vue 3 + TypeScript + Tailwind CSS |
| 实时推送 | SSE 协议（后端已实现，前端使用轮询） |
| 小程序 | 微信原生 |
| 机器人 | MCP 协议（Python，对接 QQ/微信） |
| 部署 | Docker 多阶段构建 + Docker Compose 编排 |

</details>

---

## ❓ 常见问题

<details>
<summary><b>没有专用服务器，能否用办公室电脑部署？</b></summary>
可以。Windows 系统安装 PHP 8.3 和 Composer 后，使用 SQLite 模式运行即可。建议由学校信息技术教师协助部署。
</details>

<details>
<summary><b>学生需要注册账号吗？</b></summary>
不需要。教师在后台为每个班级生成班级码（4-6位），学生在登录页输入班级码即可进入教室端。
</details>

<details>
<summary><b>教室没有外网能否正常使用？</b></summary>
可以。系统部署在校园局域网内，教室电脑与教师设备接入同一网络即可访问，完全无需互联网连接。
</details>

<details>
<summary><b>教师可以用手机操作吗？</b></summary>
可以。教师用手机浏览器访问管理页面，可执行加分、发广播、查看排行榜等操作，教室大屏实时同步。
</details>

<details>
<summary><b>广播功能如何使用？</b></summary>
教师登录后台 → 实时广播 → 选择横幅/弹窗/全屏模式 → 输入内容发送，教室大屏即时弹出显示。
</details>

<details>
<summary><b>AI 助教需要额外配置吗？</b></summary>
AI 助教功能当前尚未对接实际 AI 接口，如需使用需要自行开发后端调用逻辑。不配置不影响积分、宠物、考勤等核心功能。
</details>

<details>
<summary><b>家长如何查看学生数据？</b></summary>
家长通过手机浏览器访问学校部署地址，使用家长账号登录即可查看绑定孩子的积分、宠物、排名和通知。
</details>

<details>
<summary><b>有用户数或班级数限制吗？</b></summary>
无限制。MIT 协议不限制用户数和班级数。大规模部署建议使用 MySQL + Redis 方案，并按需配置服务器资源。
</details>

<details>
<summary><b>本项目与商业班宠系统的区别？</b></summary>
本项目完全免费开源，数据存储在校内服务器，不经过第三方平台。商业系统年费数千至数万元不等。功能覆盖度相近，但界面精致度有差距。
</details>

---

## 📄 许可证

<p align="center">
  <a href="LICENSE">MIT License</a> © 2024 RealKiro<br>
  自由使用、修改和再发布。
</p>
