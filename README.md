<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  <img src="https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/badge/PHP-8.5-777BB4?logo=php" alt="PHP 8.5">
  <img src="https://img.shields.io/badge/Laravel-12-F9322C?logo=laravel" alt="Laravel 12">
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
  <a href="#-third-party-login">第三方登录</a> •
  <a href="#-tech-stack">Tech Stack</a> •
  <a href="#-faq">FAQ</a>
</p>

---

## 📦 Features

- **积分评价系统** — 自定义加分/扣分规则，支持单人操作、批量操作、全班一键执行，每次变动实时推送至教室大屏
- **宠物进化系统** — 125 种宠物、10 大系列（山海经 / 宝可梦 / 国宝守护 / 数码宝贝 / 魔法奇幻 / 史前生物 / 星座守护 / 传统节日 / 虹猫蓝兔七侠传 / 东方神话）、12 级进化路线，积分即经验值，自动累计升级，每级独立名称、专属诗文与进化台词；每只宠物以**程序化 SVG 艺术**呈现（每物种每阶段独立剪影与进化特效，可一键切换 emoji 展示，角色档案支撑未来 AI 生图）
- **教室实时广播** — 横幅、弹窗、全屏三种模式，教师操作后 <200ms 到达教室大屏，支持多班级同时发送
- **跨班 PK 战场** — 同年级各班自动排行（总积分 / 平均等级 / 巅峰人数 / 周增长），支持发起班级挑战
- **积分商城** — 教师自定义奖品，学生自助兑换，完整的审核/发货/拒绝流程，支持多币种
- **排行榜单** — 总积分榜、周增长榜、宠物等级榜，Redis ZSET 毫秒级排序，无 Redis 自动降级 MySQL
- **智能考勤** — 一键发起点名，出勤 / 迟到 / 请假 / 缺席自动汇总
- **在线成绩管理** — 按考试和科目录入成绩，统计平均分 / 最高分 / 最低分 / 分数段分布
- **数据报表** — 积分趋势（近 4 周）、宠物等级分布、学生进步追踪、分年级/班级统计
- **实时通知** — 班级通知发布后实时推送到各班级
- **学年升级** — 预览升级明细 → 事务性执行，自动处理毕业与班级迭代
- **学生/教师管理** — Excel 批量导入、批量转班/删除、智能去重账号创建
- **多数据库支持** — SQLite / MySQL / PostgreSQL / MariaDB 任选，从小规模起步可无缝升级
- **AI 中心** — 管理后台统一配置 AI，支持 30+ 供应商（OpenAI / Claude / Gemini / DeepSeek / 通义千问 / Kimi / 豆包 / MiniMax / 百川 / GLM / 星火等），MCP 通用接口可对接任意 OpenAI 兼容服务，班级码大屏可开启 AI 对话
- **消息中心** — 实时广播 + 班级通知合并为统一入口
- **第三方平台登录** — 支持企业微信 / 钉钉 / 飞书 / 微信 / QQ / 人人通空间，管理员在后台勾选启用，登录页自动展示带官方品牌图标的平台入口；企业微信/钉钉/飞书支持扫码免注册自动建号
- **第三方账号绑定** — 教师可在账号设置中绑定/解绑第三方平台，绑定后扫码一键登录
- **通讯录批量导入** — 从企业微信/钉钉/飞书拉取通讯录，一键导入教师与学生账号
- **系统诊断与修复** — 一键检查数据库结构缺失并自动修复
- **教师端侧边栏重组** — 按课堂教学 / 成长激励 / 数据中心 / 沟通协作 / 系统管理分组
- **班级切换器** — 多班教师可在侧边栏一键切换当前班级，所有数据联动
- **Excel 导出** — 积分报表、宠物报表、考勤报表一键导出 .xlsx
- **微信小程序** — 教师端看板、积分管理、宠物图鉴、排行榜单
- **MCP 机器人协议** — 标准 MCP 服务器，对接 QQ/微信群聊机器人实现自然语言查分、排行榜
- **三端口架构** — 管理员端 / 教师端 / 教室端（班级码进入，以班级为单元）

## 🚀 Quick Start

**三步部署，新手照抄即可。** 只需要一台能装 Docker 的电脑（Windows / macOS / Linux 均可），**不需要**单独安装 PHP、MySQL、Redis、Nginx。

**第 0 步：安装 Docker**

下载安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)（Linux 服务器装 Docker Engine），装完后打开终端（Windows 用 PowerShell），输入：

```bash
docker --version
```

能打印出版本号（如 `Docker version 27.x`）就说明装好了。

**第 1 步：下载代码**

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
```

> 不会用 git？打开 GitHub 仓库页面 → 绿色 `Code` 按钮 → `Download ZIP`，解压后进入目录一样可以用。

**第 2 步：生成配置文件**

```bash
cp .env.example .env
```

> Windows PowerShell 请用：`copy .env.example .env`
>
> **默认配置 = 内置 SQLite 数据库，什么都不用装、什么都不用改**，直接进第 3 步。全校规模（>500 学生）再按下文 [数据库切换](#-deployment) 外接你已有的 MySQL/PostgreSQL 服务器。

**第 3 步：启动**

```bash
docker-compose up -d
```

首次启动会从 GitHub 拉取镜像并自动建表、创建默认学校和管理员账号，大约 1~3 分钟。

**✅ 验证部署成功（逐项核对）：**

```bash
docker-compose ps
```

- `learnstar-app` 状态为 `Up (healthy)` 即为成功
- 浏览器打开 `http://localhost:8080` → 能看到「学趣星球」首页
- 用默认管理员登录：**账号 `admin`，密码 `admin123456`**（来自 `.env` 的 `ADMIN_USERNAME` / `ADMIN_PASSWORD`，⚠️ 上线前务必修改，改完 `docker-compose up -d` 重启生效）
- 在管理后台「班级列表」创建班级后，系统会自动生成 4 位班级码，学生在首页输入班级码即可进入教室端

**日常运维三件套：**

```bash
docker-compose stop      # 停止（数据不会丢）
docker-compose start     # 再次启动
docker-compose pull && docker-compose up -d   # 升级到最新版
```

> 所有数据都存在 Docker 数据卷里，停止、重启、升级容器都不会丢数据。

## 📋 Requirements

| 部署方式 | 前置依赖 |
|---------|---------|
| Docker（推荐） | Docker Engine / Docker Desktop，仅此而已 |
| 裸机 SQLite | PHP 8.5 + Composer + Node.js 22（不推荐新手） |
| 裸机 LAMP | Nginx + PHP 8.5-FPM + MySQL（不推荐新手） |

最低硬件：1 核 CPU、512MB 内存、5GB 磁盘（办公室淘汰 PC 即可）。全校规模建议 2 核 4GB。

## 🛠 Deployment

### 方案一：SQLite（默认 · 零依赖 · 推荐起步）

就是 [Quick Start](#-quick-start三步部署新手照抄即可) 的方式：不装数据库、不装 Redis，数据保存在 Docker 数据卷里的一个 SQLite 文件（容器内路径 `storage/database.sqlite`）。

- 适合：单机、500 学生以内的学校
- 备份/恢复方法见下文 [数据备份与恢复](#-数据备份与恢复)

### 方案二：外置数据库（MySQL / MariaDB / PostgreSQL · 全校规模）

学生多、并发高时，把数据库换成你学校已有的 MySQL/MariaDB/PostgreSQL 服务器（**Docker 刻意不内置数据库容器，以保持默认部署体积最小**）：

1. 在你的数据库服务器上**手动创建**一个空库（如 `learnstar`）和账号
2. 编辑 `.env`：**注释掉**「一、SQLite」段的 10 行配置，**取消注释**「二、外置 MySQL/MariaDB」段（`.env.example` 里有分段标注，两段只能有一段生效）

   ```env
   # DB_CONNECTION=sqlite          ← 行首加 # 注释掉
   # CACHE_DRIVER=file
   # ...(SQLite 段全部注释)

   DB_CONNECTION=mysql             ← MySQL 段去掉行首 #
   DB_HOST=192.168.1.50            ← 你的数据库服务器地址
   DB_DATABASE=learnstar
   DB_USERNAME=learnstar
   DB_PASSWORD=你的数据库密码
   ```

3. 只启动应用容器（连接外部数据库）：

   ```bash
   docker-compose up -d app --no-deps
   ```

4. 验证：`docker-compose ps` 中 app 为 `Up (healthy)`，首次启动会自动在空库里建表

> 缓存可选外置 Redis（`.env` 填 `REDIS_HOST` 并把 `CACHE_DRIVER` 等改为 `redis`），排行榜走毫秒级排序；没有 Redis 就保持 `file`，功能完全可用。

### 数据库切换与数据迁移（重要，先看再切）

- **切换 = 改 `.env` + 重启**：首次以新数据库启动时会自动建表，无需手工执行 SQL
- ⚠️ **原数据库里的数据不会自动搬家**。SQLite → MySQL 切换后是全新的空库，需要重新导入教师/学生（Excel 或第三方通讯录导入均可）；积分历史、宠物等级等运行数据无法自动迁移
- **选型建议**：预计一个班试用 → SQLite；打算全校推广 → 一开始就用方案二，避免后期迁移

| | SQLite（默认） | 外置 MySQL / PostgreSQL |
|---|---|---|
| 额外容器 | 无 | 无（连你自己的数据库服务器） |
| 推荐规模 | < 500 学生 | 无限制 |
| 排行榜性能 | SQL 查询 | 可外接 Redis，毫秒级 |
| 实时广播 | 轮询 | SSE 实时推送 |

## ⚙️ Configuration

`.env` 逐项白话说明。改完任何配置，执行 `docker-compose up -d` 重建容器后生效（数据不丢）。**不要改的**：`GITHUB_USERNAME`（除非你 fork 后自行构建镜像，改成你的 GitHub 用户名并小写）。

### 基础配置

| 变量 | 默认值 | 白话说明 |
|------|--------|---------|
| `APP_PORT` | `8080` | 浏览器访问的端口。被占用就改成 `8081` 等 |
| `APP_URL` | `http://localhost` | **填别人浏览器里实际访问的地址**（不带端口）。本机试玩用默认；局域网访问改成 `http://你的电脑IP`，手机连同一 WiFi 才能打开 |
| `APP_DEBUG` | `false` | 报错时显示详细信息，仅供排障，平时保持 false |

### 管理员账号（首次启动自动创建）

| 变量 | 默认值 | 白话说明 |
|------|--------|---------|
| `ADMIN_USERNAME` | `admin` | 管理员登录账号 |
| `ADMIN_PASSWORD` | `admin123456` | ⚠️ **务必修改**。且每次重启容器都会以此值为准同步密码（忘了密码 = 改这里重启） |
| `ADMIN_NAME` / `ADMIN_SCHOOL_NAME` | 见 .env | 显示用的姓名 / 校名 |

### 数据库与缓存（三选一，详见上方 Deployment）

| 变量 | SQLite 模式 | 外置 MySQL 模式 | 白话说明 |
|------|------------|---------------|---------|
| `DB_CONNECTION` | `sqlite` | `mysql` / `pgsql` | 数据库类型 |
| `DB_HOST` | 留空 | 你的数据库服务器地址 | 数据库在哪 |
| `DB_DATABASE` 等 | 留空 | 见 .env.example | 库名 / 账号 / 密码 |
| `CACHE_DRIVER` / `SESSION_DRIVER` | `file` | `file`（有外置 Redis 则 `redis`） | 缓存与会话存哪 |
| `QUEUE_CONNECTION` | `database` | `database`（有外置 Redis 则 `redis`） | 后台任务队列 |
| `REDIS_HOST` | 留空 | 外部 Redis 地址（可选） | 有 Redis 排行榜才走毫秒级 |

### AI 助教（可选，不配不影响任何核心功能）

| 变量 | 白话说明 |
|------|---------|
| `AI_PROVIDER` / `AI_API_KEY` / `AI_MODEL` | 供应商（deepseek/openai/qwen/moonshot 等）+ 密钥 + 模型，在管理后台「AI 中心」也可视化配置（推荐） |

### 第三方平台对接（可选）

企业微信扫码登录 / 通讯录导入：填 `WECHAT_WORK_CORPID` / `WECHAT_WORK_AGENTID` / `WECHAT_WORK_SECRET`（在企业微信管理后台创建自建应用获取）；钉钉/飞书在后端 `config/dingtalk.php`、`config/feishu.php` 填凭证。详细步骤见下文 [第三方登录与配置指引](#-第三方登录与配置指引)。

### 常见问题排查（新手向）

| 现象 | 原因与解法 |
|------|-----------|
| `up` 时报端口被占用 | `APP_PORT` 改成 `8081` 等未占用端口，`APP_URL` 同步改，重启 |
| 容器一直 `restarting` 或 `unhealthy` | `docker-compose logs app` 看最后 50 行日志；多为 `.env` 数据库段改错（两段同时生效或格式错误） |
| 镜像拉取超时（国内网络） | 给 Docker 配置镜像加速器；或 fork 仓库用 Actions 自行构建，`GITHUB_USERNAME` 改成你的用户名（小写） |
| 手机打不开系统 | `APP_URL` 改成部署电脑的局域网 IP（如 `http://192.168.1.100`），防火墙放行 `APP_PORT` |
| 忘记管理员密码 | 改 `.env` 的 `ADMIN_PASSWORD` → `docker-compose up -d` 重启即同步 |
| 想彻底重置 | `docker-compose down -v`（⚠️ **删除全部数据**，包括学生积分），再 `up -d` 从零开始 |

## 💾 数据备份与恢复

### 数据都存在哪里？

| 数据 | 位置（容器内） | Docker 数据卷 |
|------|--------------|--------------|
| SQLite 数据库（积分/学生/教师/商城全部业务数据） | `/app/storage/database.sqlite` | `app-db` |
| 上传的附件（Logo 等） | `/app/storage/app/uploads` | `app-uploads` |
| 运行日志 | `/app/storage/logs` | `app-logs` |

> 外置 MySQL/PostgreSQL 模式下，业务数据在你自己的数据库服务器上，随你的数据库备份策略走。

各数据卷由 Docker 管理，`stop` / `start` / `restart` / `up -d` / 升级镜像都**不会丢数据**。只有 `docker-compose down -v` 会连数据卷一起删除。

### 备份（SQLite 模式）

先停应用再拷贝，避免拷到写一半的文件：

```bash
docker-compose stop app
docker cp learnstar-app:/app/storage/database.sqlite ./backup.sqlite
docker-compose start app
```

> Windows PowerShell 把 `./backup.sqlite` 换成 `.\backup.sqlite` 即可。
> 建议每周备份一次，重大操作（升班、批量导入）前手动备份一次。

### 恢复（SQLite 模式）

```bash
docker-compose stop app
docker cp ./backup.sqlite learnstar-app:/app/storage/database.sqlite
docker-compose start app
```

恢复 = 用备份文件覆盖数据库文件，所有数据回到备份那一刻。

### 备份 / 恢复（外置 MySQL 模式）

外置模式下业务数据在你自己的 MySQL 服务器上，用标准的 mysqldump 流程（在你的数据库服务器或任何能连上它的机器上执行）：

```bash
# 备份
mysqldump -h <数据库地址> -u learnstar -p learnstar > backup.sql

# 恢复
mysql -h <数据库地址> -u learnstar -p learnstar < backup.sql
```

`backup.sql` 是纯文本 SQL，请妥善保管（含学生与积分数据）。

### ⚠️ 旧版本升级注意（2026-09 之前部署的）

早期版本的 docker-compose 把数据卷挂错了路径，**数据库实际存在容器内部，`down` 或升级镜像会丢数据**。升级到本版本前，先从旧容器把数据库抢救出来：

```bash
docker cp learnstar-app:/app/storage/database.sqlite ./backup-before-upgrade.sqlite
```

然后正常 `git pull && docker-compose pull && docker-compose up -d`（会重建容器），最后用上面的「恢复」命令把数据灌回去。**只需做这一次**，之后的升级都安全了。

## 🏗 Tech Stack

| 层级 | 技术 |
|------|------|
| 后端 | Laravel 12（PHP 8.5） |
| 数据库 | SQLite / MySQL / MariaDB / PostgreSQL |
| 缓存与队列 | Redis 8（可降级 file / database） |
| 前端 | Vue 3 + TypeScript + Tailwind CSS |
| 实时推送 | SSE 协议（后端实现，前端 EventSource 优先 / 轮询降级） |
| 小程序 | 微信原生（教师端 / 教室端） |
| 机器人 | MCP 协议 Python 服务器 |
| 部署 | Docker 多阶段构建 + Compose 编排 / 裸 PHP 环境 |

## 🗂 Project Structure

```
learnstar-planet/
├── frontend-vue/          # Vue 3 SPA（管理员端 / 教师端 / 教室端）
├── backend/               # Laravel 12 API（24 模型，214 端点）
├── mini-program/          # 微信小程序（10 页面，教师端）
├── pwa/                   # PWA 离线配置
├── mcp-server/            # MCP 机器人服务
└── docker-compose.yml     # Docker 编排
```

## ❓ FAQ

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
<summary>学生怎么查看自己的数据？</summary>
学生无需单独登录，通过班级码进入教室端（/classroom）即可查看本班积分、排行榜、宠物与图鉴。
</details>

<details>
<summary>有用户数限制吗？</summary>
无。MIT 开源协议不限制用户数和班级数。全校使用建议外接你已有的 MySQL 服务器。
</details>

## 🔐 第三方登录与配置指引

学趣星球支持多平台第三方扫码登录。**管理员在后台勾选哪些平台启用**，勾选后登录页才会显示对应平台入口，并自动展示各平台官方品牌图标。

### 一、后台启用第三方平台

1. 管理员登录 → **学校设置** → 「第三方登录平台」
2. 勾选要启用的平台（可多选）：**企业微信 / 钉钉 / 飞书 / 人人通空间 / 微信 / QQ**
3. 保存后，登录页教师端自动显示勾选的平台入口

> 未勾选任何平台时，默认启用 **企业微信 / 微信 / QQ**。
> 勾选平台的展示与登录无需额外配置；但**扫码登录真正可用**需在对应平台开放平台注册应用并配置凭证（见下表）。

### 二、各平台配置凭证

| 平台 | 是否需要凭证 | 凭证位置 | 支持能力 |
|------|:-----------:|----------|----------|
| **企业微信** | ✅ | `.env`：`WECHAT_WORK_CORPID` / `WECHAT_WORK_AGENTID` / `WECHAT_WORK_SECRET` | 扫码免注册登录 + 通讯录导入 + 请假同步 + 消息推送 |
| **钉钉** | ✅ | `backend/config/dingtalk.php`：`app_key` / `app_secret` | 扫码登录 + 通讯录导入 |
| **飞书** | ✅ | `backend/config/feishu.php`：`app_id` / `app_secret` | 扫码登录 + 通讯录导入 |
| **微信** | ⚠️ | 需前端接入微信开放平台 JSAPI 取 openid | 当前仅展示入口，扫码流程待接入 |
| **QQ** | ⚠️ | 需前端接入 QQ 互联取 openid | 当前仅展示入口，扫码流程待接入 |
| **人人通空间** | ⚠️ | 需对接区域人人通开放平台 | 当前仅展示入口，扫码流程待接入 |

### 三、教师绑定与免注册登录

- **企业微信 / 钉钉 / 飞书**：教师扫码后，若系统已有绑定账号则直接登录；否则**根据平台通讯录自动创建教师账号**（实名用户名 + 默认密码 `ls123456`）并自动绑定，无需管理员手动分配
- **账号设置 → 第三方账号绑定**：教师可查看/解绑已绑定的第三方平台，或扫码绑定新平台

### 四、通讯录批量导入

管理员 → 教师管理 → **🏢 第三方导入**：从学校配置的第三方平台（企业微信/钉钉/飞书）拉取通讯录，勾选成员后批量创建教师与学生账号。教师默认以姓名作为登录账号（实名），学生需选择目标班级（可按部门名自动匹配）。

### 五、工作原理（可选了解）

- 登录页平台列表来自接口 `GET /api/v1/auth/third-party/options`（返回管理员勾选的平台 + 品牌图标）
- 学校配置存储在 `schools.settings.enabled_third_party_platforms`（JSON 数组）
- 扫码回调按学校配置的平台分发（`App\Services\ThirdParty\ThirdPartyManager`），多校部署时通过 OAuth `state` 参数区分学校

### 五、数据同步与冲突处理（升班 / 名单更新必读）

系统以**本地数据库为准**，第三方平台（企业微信/钉钉/飞书）是数据来源之一。平台侧的调整不会自动写入本地，按以下规则协同：

| 场景 | 系统的处理 |
|------|-----------|
| 通讯录重复导入 | 教师按手机号/实名账号自动跳过已有账号，**不会**生成"张老师_2"这类冗余账号；学生按"同班同名"跳过 |
| 学生转班（平台名单已调整） | Excel/通讯录导入时，**同学号已在其他班级会被拦截**并提示用「批量转班」处理；跨班同名会给出提醒供管理员判断（同名不同人可直接忽略） |
| 教师未导入过、直接扫码登录 | 系统按手机号 → 实名用户名匹配本地已有账号并自动绑定，**不会**重复建号 |
| 学年升班 | 在本系统「学年升级」执行（六年级毕业、班级整体升级）；**第三方平台的部门/名单需要平台管理员同步调整**，之后再回来导入通讯录并核对班级映射 |

> 推荐顺序：每个学年开始时，先在第三方平台完成升班与名单调整 → 本系统执行「学年升级」→ 建新一年级班级 → 通讯录导入新生（同学号冲突会被自动拦截）。

## 🤝 Contributing

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/my-feature`
3. 提交变更：`git commit -am 'feat: add my feature'`
4. 推送：`git push origin feature/my-feature`
5. 提交 Pull Request

代码规范：PHP PSR-12（PHP CS Fixer）· TypeScript ESLint · Conventional Commits

## 📄 License

[MIT License](LICENSE) © 2024 RealKiro
