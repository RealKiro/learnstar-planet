<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  <img src="https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/badge/PHP-8.3-777BB4?logo=php" alt="PHP 8.3">
  <img src="https://img.shields.io/badge/Laravel-11-F9322C?logo=laravel" alt="Laravel 11">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" alt="Docker">
</p>

# 🌌 学趣星球 (LearnStar Planet)

<p align="center">
  <b>教室大屏互动 + 积分激励 + 宠物进化</b><br>
  局域网就能用 · 学生不用注册 · 完全免费开源
</p>

```bash
# 有 Docker 的环境，一分钟跑起来
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet && cp .env.example .env
docker-compose up -d
# 浏览器打开 http://localhost:8080
```

---

## 📖 目录

- [这系统能干什么](#-这系统能干什么)
- [功能清单和实现状态](#-功能清单和实现状态)
- [为什么适合学校用](#-为什么适合学校用)
- [部署方式](#-部署方式)
- [系统架构](#-系统架构)
- [常见问题](#-常见问题)

---

## 🎯 这系统能干什么

给教室接一台能投屏的电脑/电视，老师用手机或另一台电脑操作，大屏上实时显示。核心逻辑就是：**学生表现好→加分→宠物升级；表现不好→扣分。** 简单直接。

**不需要买服务器、不需要联网、学生不用注册账号**——输入老师给的班级码就能进。

适合的场景：
- 上课加分，全班大屏实时看到
- 老师发个通知，大屏上立刻弹出来
- 学期末看哪个小组积分高，宠物进化到哪一级了
- 家长想看看孩子在校表现，手机打开就能看

---

## ✅ 功能清单和实现状态

以下功能都是**已经实现了的**（后端有真实逻辑、前端有对应页面），不是画饼。

### 全部实现的功能

| 功能 | 说明 | 后端 | 前端 |
|------|------|------|------|
| **实时广播** | 教师发通知 → 教室大屏弹出（弹窗/横幅/全屏三种模式），适合上课喊话、自习课通知 | ✅ | ✅ |
| **班级大屏** | 投屏后显示全班积分之星、最新动态、TOP5排名，加分时宠物弹跳/抖动的反馈动画 | ✅ | ✅ |
| **课堂加减分** | 点击学生头像选行为原因（答对题/作业优秀/违反纪律），可自定义步长和规则 | ✅ | ✅ |
| **宠物系统** | 54种宠物×8大系列×12级进化，积分喂养升级，每级有专属名字和诗文 | ✅ | ✅ |
| **跨班PK** | 同年级各班积分排行榜、战力分析、发起挑战 | ✅ | ✅ |
| **积分商城** | 老师设置兑奖规则，学生自助兑换，老师审核/发货 | ✅ | ✅ |
| **排行榜** | 总积分榜/周榜/宠物等级榜（支持Redis加速，没有Redis也能用） | ✅ | ✅ |
| **智能考勤** | 大屏上一键点名，标记出勤/迟到/请假/缺席 | ✅ | ✅ |
| **班级通知** | 发通知到班级，家长端同步可见 | ✅ | ✅ |
| **作业管理** | 布置作业、生成二维码学生扫码提交、查看谁交了谁没交 | ✅ | ✅ |
| **在线测验** | 创建题库→出卷→学生作答→自动判分 | ✅ | ✅ |
| **成绩管理** | 按考试/科目录入成绩，统计分析和分布查看 | ✅ | ✅ |
| **学年升级** | 预览升级情况→确认执行，自动更新班级年级 | ✅ | ✅ |
| **学生管理** | Excel批量导入/创建/转班/删除 | ✅ | ✅ |
| **教师管理** | 批量创建账号、分配班级 | ✅ | ✅ |
| **数据报表** | 积分趋势、宠物分布、学生进步情况、分年级/班级统计 | ✅ | ✅ |
| **家长端** | 手机浏览器查看孩子积分、宠物、排名、通知 | ✅ | ✅ |
| **第三方登录** | 教师可用微信/QQ/人人扫码登录（需要对应平台申请） | ✅ | ⚠️ 简单 |
| **微信小程序** | 家长在微信里打开使用 | ✅ | ⚠️ 部分页面 |
| **SSE实时推送** | 后台变动 <200ms 推送到大屏，断线自动重连，彻底失败降级轮询 | ✅ | — |
| **MCP机器人** | 对接QQ/微信群聊机器人，自然语言查积分、排行榜（需要额外部署AstrBot） | ✅ | — |
| **Excel导出** | 报表导出Excel/PDF | ✅ | — |

### ⚠️ 需要注意的功能

| 功能 | 实际状态 |
|------|---------|
| **AI助教** | 后端有请求接口和前端聊天界面，但**需要自己配AI接口**（通义千问/DeepSeek/OpenAI的API Key），不配的话AI回复固定提示 |
| **微信小程序** | 14个页面已注册，目前实际构建了其中2个核心页面（首页+看板），其余为框架已搭但内容待完善 |
| **在线测验学生端** | 老师端出卷和自动判分OK，学生扫码作答的流程还需要完善 |

### 没做的事情
- ✖ 没有手机App（用手机浏览器就可以访问）
- ✖ 没有视频直播/网课功能（这是班级管理系统，不是在线教育平台）

---

## 🏫 为什么适合学校用

| 特点 | 说明 |
|------|------|
| **一分钱不花** | MIT协议开源，没有隐藏收费、没有用户数限制、没有功能锁定 |
| **不用联网** | 全部部署在学校局域网内，教室电脑、教师手机连同一个校园网就能用 |
| **旧电脑能跑** | 办公室淘汰的电脑就可以当服务器，装PHP或Docker都行 |
| **学生不用注册** | 老师给每个班生成一个班级码，学生输入4-6位码就能进入，没有忘记密码这种事 |
| **数据在校内** | 所有数据存在学校自己的机器上，不出校门 |

---

## 🚀 部署方式

### 快速选型

| 方式 | 需要什么 | 适合场景 |
|------|---------|---------|
| ① Docker部署 | 能装Docker的电脑 | 一步到位，推荐 |
| ② SQLite部署 | PHP 8.3即可 | 旧电脑，不想装数据库 |
| ③ LAMP部署 | Nginx+PHP+MySQL | 学校有服务器 |

### 方式一：Docker 一键部署

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

启动后访问：
- **教室大屏**：浏览器打开 `http://你的电脑IP:8080`，输入班级码
- **教师管理**：`http://你的电脑IP:8080/login`，管理员或教师账号登录

> 全校师生通过局域网IP就能访问，不需要互联网。手机连同一个WiFi也能打开。

### 方式二：SQLite 部署（不用Docker、不用MySQL、不用Redis）

<details>
<summary><b>📋 展开步骤</b></summary>

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

### SQLite 和 MySQL 对比

| | SQLite | MySQL + Redis |
|---|---|---|
| 需要安装 | 不需要 | 需要安装 |
| 建议规模 | < 500 学生 | 无限制 |
| 排行榜速度 | 常规SQL查询 | Redis加速（毫秒级） |
| 实时广播 | 不支持 | 支持 |
| 适用场景 | 一间教室/单个年级 | 全校使用 |

> 先用 SQLite 跑起来，以后人多了改两行配置就能无缝升级到 MySQL。

---

## 🏗 系统架构

```
learnstar-planet/
├── frontend-vue/          # 前端（Vue 3 + TypeScript）
│   ├── pages/teacher/     # 教师管理端（17个页面）
│   ├── pages/classroom/   # 教室大屏端（4个页面）
│   ├── pages/parent/      # 家长端（6个页面）
│   └── pages/admin/       # 学校管理端（9个页面）
├── backend/               # 后端（Laravel 11，17个数据模型，~130个接口）
├── mini-program/          # 微信小程序（14页面注册，2个已构建）
├── pwa/                   # PWA离线缓存配置
├── mcp-server/            # QQ/微信机器人（Python，可选）
└── docker-compose.yml     # Docker编排
```

<details>
<summary><b>📋 技术栈详情</b></summary>

| 层级 | 技术 |
|------|------|
| 后端框架 | Laravel 11（PHP 8.3） |
| 数据库 | SQLite / MySQL / MariaDB / PostgreSQL 任选 |
| 前端 | Vue 3 + TypeScript + Tailwind CSS |
| 实时推送 | SSE（服务端推送事件） |
| AI | 通义千问 / DeepSeek / OpenAI（需自配Key） |

</details>

---

## ❓ 常见问题

<details>
<summary><b>我没服务器，用自己办公室的旧电脑行不行？</b></summary>
行。Windows电脑装PHP 8.3和Composer，用SQLite模式跑起来就行。推荐让学校电教老师帮忙搞。
</details>

<details>
<summary><b>学生怎么登录？要注册吗？</b></summary>
不用注册。老师在后台给每个班生成一个4-6位的班级码，学生在首页输入这个码就能进去。
</details>

<details>
<summary><b>教室没外网能用吗？</b></summary>
能。系统部署在局域网内，教室电脑、教师手机连同一个校园网就能访问，完全不需要上外网。
</details>

<details>
<summary><b>手机能操作吗？</b></summary>
能。教师用手机浏览器打开管理页面，可以加分、发广播、看排行榜。手机上操作，大屏实时同步。
</details>

<details>
<summary><b>广播功能具体怎么用？</b></summary>
教师登录后台→实时广播，选横幅/弹窗/全屏，输入内容发送，教室大屏立刻弹出。适合：上课发通知、自习课提醒、课间活动通知。
</details>

<details>
<summary><b>AI助教要额外花钱吗？</b></summary>
AI需要你自己去通义千问/DeepSeek/OpenAI申请API Key，这些服务商有免费额度也有付费方案。不配AI也不影响任何核心功能。
</details>

<details>
<summary><b>家长怎么看孩子的表现？</b></summary>
家长用手机浏览器打开学校给的访问地址，用家长账号登录，就能看到孩子的积分、宠物、排名、通知。
</details>

<details>
<summary><b>有用户数限制吗？能支持全校用吗？</b></summary>
MIT开源协议不限用户数和班级数。全校用建议用MySQL方式部署，配一台好一点的服务器就行。
</details>

<details>
<summary><b>你们和商业班宠系统有什么区别？</b></summary>
完全免费、开源、数据在校内。商业系统每年收几千到几万，这个一分钱不用花。功能上该有的都有，界面可能没那么精致。
</details>

---

## 📄 许可证

[MIT License](LICENSE) © 2024 RealKiro

> 全免费、开源、可自由使用和修改。
