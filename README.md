<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  <img src="https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/badge/PHP-8.5-777BB4?logo=php" alt="PHP 8.5">
  <img src="https://img.shields.io/badge/Laravel-12-F9322C?logo=laravel" alt="Laravel 12">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" alt="Docker">
</p>

<h1 align="center">学宠星球 / LearnStar Planet</h1>

<p align="center">
  开源、自托管的班级管理与学生激励系统。<br>
  积分评价 · 宠物进化 · 课堂互动 · 数据报表 · AI 助教
</p>

<p align="center">
  <b>📖 文档地图</b><br><br>
  <b>🚀 上手</b>（第一次用，从这里开始）<br>
  <a href="#2-快速开始">2. 快速开始</a> • <a href="#3-部署指南">3. 部署指南</a> • <a href="#4-配置说明">4. 配置说明</a> • <a href="#5-数据备份与恢复">5. 数据备份与恢复</a><br><br>
  <b>🔎 了解</b>（选型与能力一览）<br>
  <a href="#1-功能特性">1. 功能特性</a> • <a href="#7-常见问题faq">7. 常见问题（FAQ）</a> • <a href="#8-技术架构">8. 技术架构</a><br><br>
  <b>🔗 进阶</b>（对接与扩展）<br>
  <a href="#6-集成与对接">6. 集成与对接</a>（REST API 机器人账号 · MCP · 第三方平台）
</p>

---

> **💡 按场景阅读**
> - 5 分钟跑起来 → 只看 [2. 快速开始](#2-快速开始)
> - 全校推广 → [3. 部署指南](#3-部署指南)（数据库切换）+ [4. 配置说明](#4-配置说明)
> - 日常维护 → [5. 数据备份与恢复](#5-数据备份与恢复)（含忘密码 / 重置 / 旧版数据抢救）
> - 让机器人或别的系统管积分 → [6. 集成与对接](#6-集成与对接)

---

## 1. 功能特性

### 1.1 课堂管理

| 功能 | 说明 |
|------|------|
| **积分评价系统** | 自定义加分/扣分规则，支持单人操作、批量操作、全班一键执行，每次变动实时推送至教室大屏 |
| **教室实时广播** | 横幅、弹窗、全屏三种模式，教师操作后 <200ms 到达教室大屏，支持多班级同时发送 |
| **智能考勤** | 一键发起点名，出勤 / 迟到 / 请假 / 缺席自动汇总 |
| **在线成绩管理** | 按考试和科目录入成绩，统计平均分 / 最高分 / 最低分 / 分数段分布 |
| **消息中心** | 实时广播 + 班级通知合并为统一入口 |
| **班级切换器** | 多班教师可在侧边栏一键切换当前班级，所有数据联动 |

### 1.2 成长激励

| 功能 | 说明 |
|------|------|
| **宠物进化系统** | 125 种宠物、10 大系列（山海经 / 宝可梦 / 国宝守护 / 数码宝贝 / 魔法奇幻 / 史前生物 / 星座守护 / 传统节日 / 虹猫蓝兔七侠传 / 东方神话）、12 级进化路线，积分即经验值；每只宠物以程序化 SVG 艺术呈现（可一键切换 emoji 展示），角色档案支撑未来 AI 生图 |
| **跨班 PK 战场** | 同年级各班自动排行（总积分 / 平均等级 / 巅峰人数 / 周增长），支持发起班级挑战 |
| **积分商城** | 教师自定义奖品，学生自助兑换，完整的审核/发货/拒绝流程，支持多币种 |
| **排行榜单** | 总积分榜、周增长榜、宠物等级榜，Redis ZSET 毫秒级排序，无 Redis 自动降级 SQL |
| **学年升级** | 预览升级明细 → 事务性执行，自动处理毕业与班级迭代 |

### 1.3 数据与运营

| 功能 | 说明 |
|------|------|
| **数据报表** | 积分趋势（近 4 周）、宠物等级分布、学生进步追踪、分年级/班级统计 |
| **学生 / 教师管理** | Excel 批量导入、批量转班/删除、智能去重账号创建 |
| **Excel 导出** | 积分报表、宠物报表、考勤报表一键导出 .xlsx |
| **系统诊断与修复** | 一键检查数据库结构缺失并自动修复 |
| **多数据库支持** | SQLite / MySQL / PostgreSQL / MariaDB 任选，从小规模起步可无缝升级 |

### 1.4 AI 能力

| 功能 | 说明 |
|------|------|
| **AI 中心** | 管理后台统一配置 AI，支持 30+ 供应商（OpenAI / Claude / Gemini / DeepSeek / 通义千问 / Kimi / 豆包 / MiniMax / 百川 / GLM / 星火等） |
| **AI 助教** | 班级码大屏可开启 AI 对话，学生直接向 AI 提问 |
| **MCP 通用接口** | 可对接任意 OpenAI 兼容服务（自建 vLLM、本地大模型等） |

### 1.5 平台与集成

| 功能 | 说明 |
|------|------|
| **三端架构** | 管理员端 / 教师端 / 教室端（班级码进入，以班级为单元），学生无需账号 |
| **REST API 机器人账号** | 内置机器人教师账号（全部班级权限），供外部系统调用积分管理能力 |
| **MCP 机器人协议** | 标准 MCP 服务器，对接 QQ/微信群聊机器人实现自然语言加减分、查分、排行榜 |
| **第三方平台登录** | 企业微信 / 钉钉 / 飞书 / 微信 / QQ / 人人通空间，扫码免注册自动建号，详见 [6. 集成与对接](#6-集成与对接) |
| **通讯录批量导入** | 从企业微信/钉钉/飞书拉取通讯录，一键导入教师与学生账号 |
| **微信小程序** | 教师端看板、积分管理、宠物图鉴、排行榜单 |

## 2. 快速开始

**三步部署（下载 → 配置 → 启动），新手照抄即可。** 只需要一台能装 Docker 的电脑（Windows / macOS / Linux 均可），**不需要**单独安装 PHP、MySQL、Redis、Nginx。

### 2.1 硬件要求

最低 1 核 CPU、512MB 内存、5GB 磁盘（办公室淘汰 PC 即可）。全校规模建议 2 核 4GB。

### 2.2 安装 Docker

下载安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)（Linux 服务器装 Docker Engine），装完后打开终端（Windows 用 PowerShell），输入：

```bash
docker --version
```

能打印出版本号（如 `Docker version 27.x`）就说明装好了。

### 2.3 下载代码

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
```

> 不会用 git？打开 GitHub 仓库页面 → 绿色 `Code` 按钮 → `Download ZIP`，解压后进入目录一样可以用。

### 2.4 生成配置文件

```bash
cp .env.example .env
```

> Windows PowerShell 请用：`copy .env.example .env`
>
> **默认配置 = 内置 SQLite 数据库，什么都不用装、什么都不用改**，直接进 [2.5 启动](#25-启动)。全校规模（>500 学生）再按下文 [3. 部署指南](#3-部署指南) 外接你已有的 MySQL/PostgreSQL 服务器。

### 2.5 启动

```bash
docker-compose up -d
```

首次启动会从 GitHub 拉取镜像并自动建表、创建默认学校和管理员账号，大约 1~3 分钟。

### 2.6 验证部署成功（逐项核对）

```bash
docker-compose ps
```

- `learnstar-app` 状态为 `Up (healthy)` 即为成功
- 浏览器打开 `http://localhost:8080` → 能看到「学宠星球」首页
- 用默认管理员登录：**账号 `admin`，密码 `admin123456`**（来自 `.env` 的 `ADMIN_USERNAME` / `ADMIN_PASSWORD`，⚠️ 上线前务必修改，改完 `docker-compose up -d` 重启生效）
- 在管理后台「班级列表」创建班级后，系统会自动生成 4 位班级码，学生在首页输入班级码即可进入教室端

### 2.7 日常运维

```bash
docker-compose stop      # 停止（数据不会丢）
docker-compose start     # 再次启动
docker-compose pull && docker-compose up -d   # 升级到最新版
```

> 所有数据都存在 Docker 数据卷里，停止、重启、升级容器都不会丢数据。

## 3. 部署指南

### 3.1 SQLite（默认 · 零依赖 · 推荐起步）

就是 [2. 快速开始](#2-快速开始) 的方式：不装数据库、不装 Redis，数据保存在 Docker 数据卷里的一个 SQLite 文件（容器内路径 `storage/database.sqlite`）。

- 适合：单机、500 学生以内的学校
- 备份/恢复方法见下文 [5. 数据备份与恢复](#5-数据备份与恢复)

### 3.2 外置数据库（MySQL / MariaDB / PostgreSQL · 全校规模）

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

### 3.3 数据库切换与数据迁移（重要，先看再切）

- **切换 = 改 `.env` + 重启**：首次以新数据库启动时会自动建表，无需手工执行 SQL
- ⚠️ **原数据库里的数据不会自动搬家**。SQLite → MySQL 切换后是全新的空库，需要重新导入教师/学生（Excel 或第三方通讯录导入均可）；积分历史、宠物等级等运行数据无法自动迁移
- **选型建议**：预计一个班试用 → SQLite；打算全校推广 → 一开始就用外置数据库（见 3.2），避免后期迁移

| 对比项 | SQLite（默认） | 外置 MySQL / PostgreSQL |
|---|---|---|
| 额外容器 | 无 | 无（连你自己的数据库服务器） |
| 推荐规模 | < 500 学生 | 无限制 |
| 排行榜性能 | SQL 查询 | 可外接 Redis，毫秒级 |
| 实时广播 | 轮询 | SSE 实时推送 |

## 4. 配置说明

`.env` 逐项白话说明。改完任何配置，执行 `docker-compose up -d` 重建容器后生效（数据不丢）。**不要改的**：`GITHUB_USERNAME`（除非你 fork 后自行构建镜像，改成你的 GitHub 用户名并小写）。

### 4.1 基础配置

| 变量 | 默认值 | 白话说明 |
|------|--------|---------|
| `APP_PORT` | `8080` | 浏览器访问的端口。被占用就改成 `8081` 等 |
| `APP_URL` | `http://localhost` | **填别人浏览器里实际访问的地址**（不带端口）。本机试玩用默认；局域网访问改成 `http://你的电脑IP`，手机连同一 WiFi 才能打开 |
| `APP_DEBUG` | `false` | 报错时显示详细信息，仅供排障，平时保持 false |

### 4.2 管理员账号（首次启动自动创建）

| 变量 | 默认值 | 白话说明 |
|------|--------|---------|
| `ADMIN_USERNAME` | `admin` | 管理员登录账号 |
| `ADMIN_PASSWORD` | `admin123456` | ⚠️ **务必修改**。且每次重启容器都会以此值为准同步密码（忘了密码 = 改这里重启） |
| `ADMIN_NAME` / `ADMIN_SCHOOL_NAME` | 见 .env | 显示用的姓名 / 校名 |

### 4.3 数据库与缓存（二选一，详见「3. 部署指南」）

| 变量 | SQLite 模式 | 外置 MySQL 模式 | 白话说明 |
|------|------------|---------------|---------|
| `DB_CONNECTION` | `sqlite` | `mysql` / `pgsql` | 数据库类型 |
| `DB_HOST` | 留空 | 你的数据库服务器地址 | 数据库在哪 |
| `DB_DATABASE` 等 | 留空 | 见 .env.example | 库名 / 账号 / 密码 |
| `CACHE_DRIVER` / `SESSION_DRIVER` | `file` | `file`（有外置 Redis 则 `redis`） | 缓存与会话存哪 |
| `QUEUE_CONNECTION` | `database` | `database`（有外置 Redis 则 `redis`） | 后台任务队列 |
| `REDIS_HOST` | 留空 | 外部 Redis 地址（可选） | 有 Redis 排行榜才走毫秒级 |

### 4.4 AI 助教（可选，不配不影响任何核心功能）

| 变量 | 白话说明 |
|------|---------|
| `AI_PROVIDER` / `AI_API_KEY` / `AI_MODEL` | 供应商（deepseek/openai/qwen/moonshot 等）+ 密钥 + 模型，在管理后台「AI 中心」也可视化配置（推荐） |

### 4.5 第三方平台对接（可选）

企业微信扫码登录 / 通讯录导入：填 `WECHAT_WORK_CORPID` / `WECHAT_WORK_AGENTID` / `WECHAT_WORK_SECRET`（在企业微信管理后台创建自建应用获取）；钉钉/飞书在后端 `config/dingtalk.php`、`config/feishu.php` 填凭证。详细步骤见 [6. 集成与对接](#6-集成与对接)。

## 5. 数据备份与恢复

> **⚡ 场景速查**
> - 我要备份数据 → [5.2 备份（SQLite 模式）](#52-备份sqlite-模式) / [5.4 备份 / 恢复（外置 MySQL 模式）](#54-备份--恢复外置-mysql-模式)
> - 我要恢复数据 → [5.3 恢复（SQLite 模式）](#53-恢复sqlite-模式) / [5.4 备份 / 恢复（外置 MySQL 模式）](#54-备份--恢复外置-mysql-模式)
> - 系统坏了想重来 → 先 [5.2 备份](#52-备份sqlite-模式)，再看 [7. 常见问题（FAQ）](#7-常见问题faq) 的「彻底重置」

### 5.1 数据都存在哪里？

| 数据 | 位置（容器内） | Docker 数据卷 |
|------|--------------|--------------|
| SQLite 数据库（积分/学生/教师/商城全部业务数据） | `/app/storage/database.sqlite` | `app-db` |
| 上传的附件（Logo 等） | `/app/storage/app/uploads` | `app-uploads` |
| 运行日志 | `/app/storage/logs` | `app-logs` |

> 外置 MySQL/PostgreSQL 模式下，业务数据在你自己的数据库服务器上，随你的数据库备份策略走。

各数据卷由 Docker 管理，`stop` / `start` / `restart` / `up -d` / 升级镜像都**不会丢数据**。只有 `docker-compose down -v` 会连数据卷一起删除。

### 5.2 备份（SQLite 模式）

先停应用再拷贝，避免拷到写一半的文件：

```bash
docker-compose stop app
docker cp learnstar-app:/app/storage/database.sqlite ./backup.sqlite
docker-compose start app
```

> Windows PowerShell 把 `./backup.sqlite` 换成 `.\backup.sqlite` 即可。
> 建议每周备份一次，重大操作（升班、批量导入）前手动备份一次。

### 5.3 恢复（SQLite 模式）

```bash
docker-compose stop app
docker cp ./backup.sqlite learnstar-app:/app/storage/database.sqlite
docker-compose start app
```

恢复 = 用备份文件覆盖数据库文件，所有数据回到备份那一刻。

### 5.4 备份 / 恢复（外置 MySQL 模式）

外置模式下业务数据在你自己的 MySQL 服务器上，用标准的 mysqldump 流程（在你的数据库服务器或任何能连上它的机器上执行）：

```bash
# 备份
mysqldump -h <数据库地址> -u learnstar -p learnstar > backup.sql

# 恢复
mysql -h <数据库地址> -u learnstar -p learnstar < backup.sql
```

`backup.sql` 是纯文本 SQL，请妥善保管（含学生与积分数据）。

### 5.5 旧版本升级注意（仅 2026-09 之前部署的需要看一次）

<details>
<summary>展开：如何把旧版本容器里的数据抢救出来</summary>

早期版本的 docker-compose 把数据卷挂错了路径，**数据库实际存在容器内部，`down` 或升级镜像会丢数据**。升级到本版本前，先从旧容器把数据库抢救出来：

```bash
docker cp learnstar-app:/app/storage/database.sqlite ./backup-before-upgrade.sqlite
```

然后正常 `git pull && docker-compose pull && docker-compose up -d`（会重建容器），最后用上面的「恢复」命令把数据灌回去。**只需做这一次**，之后的升级都安全了。

</details>

## 6. 集成与对接

本项目的三类对外能力都在这一章：REST API 机器人账号（通用对接）、MCP 服务器（AI 机器人插件）、第三方平台登录与同步（企业微信/钉钉/飞书）。

| 小节 | 适合谁 |
|--------|--------|
| [6.1 REST API 机器人账号（外部系统对接）](#61-rest-api-机器人账号外部系统对接) | 想让**自己的项目/脚本**调用积分管理能力 |
| [6.2 MCP 服务器（AI 机器人插件）](#62-mcp-服务器ai-机器人插件) | 想接 **AstrBot / Claude Desktop** 等 AI 宿主，自然语言管积分 |
| [6.3 第三方平台登录与通讯录同步](#63-第三方平台登录与通讯录同步) | 学校用**企业微信 / 钉钉 / 飞书**，要扫码登录与名单导入 |

### 6.1 REST API 机器人账号（外部系统对接）

系统内置一个专用的 **API 机器人教师账号**，供 QQ/微信群机器人、课中工具等外部项目通过 REST API 管理学生积分。该账号**拥有本校全部班级的查询与操作权限**（含未来新建的班级，无需维护关联）。

三步接入：

```bash
# 1. 用机器人账号换取 token（账号密码来自 .env 的 BOT_USERNAME / BOT_PASSWORD）
TOKEN=$(curl -s -X POST http://<服务器>:8080/api/v1/auth/teacher/login \
  -H "Content-Type: application/json" \
  -d '{"username":"api-bot","password":"learnstar-bot-2026"}' | jq -r .data.token)

# 2. 查询可用班级（机器人返回本校全部班级）
curl -H "Authorization: Bearer $TOKEN" http://<服务器>:8080/api/v1/teacher/my-classes

# 3. 给学生加分（student_id 来自 /api/v1/teacher/students?per_page=100）
curl -X POST http://<服务器>:8080/api/v1/teacher/scores/give \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"student_id":1,"points":5,"reason":"作业优秀"}'
```

常用端点速查：

| 能力 | 端点 |
|------|------|
| 班级列表 | `GET /api/v1/teacher/my-classes` |
| 学生列表 | `GET /api/v1/teacher/students?per_page=100` |
| 单个加减分 | `POST /api/v1/teacher/scores/give` |
| 批量加减分 | `POST /api/v1/teacher/scores/batch-give` |
| 按规则加减分 | `POST /api/v1/teacher/scores/give-by-rule/{ruleId}` |
| 撤销 | `POST /api/v1/teacher/scores/{id}/undo` |
| 历史 / 汇总 | `GET /api/v1/teacher/scores/history/{studentId}` · `GET /api/v1/teacher/scores/summary` |

完整清单见 [docs/api-reference.md](docs/api-reference.md)。

配置与安全：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `BOT_ENABLED` | `true` | 设为 `false` 并重启，账号会被置为 disabled（接口 401） |
| `BOT_USERNAME` / `BOT_PASSWORD` | `api-bot` / `learnstar-bot-2026` | ⚠️ 密码以 .env 为唯一真相来源，每次重启同步，**上线前务必修改** |
| `BOT_NAME` | `API 机器人` | 在教师列表中的显示名 |

安全须知：机器人账号可操作**全部班级**，请只在内网/可信环境使用并修改默认密码；积分记录的审计字段 `given_by` 会归到该账号；账号无法从管理后台删除（如需停用走 `BOT_ENABLED=false`）。

### 6.2 MCP 服务器（AI 机器人插件）

[mcp-server/](mcp-server/) 提供标准 MCP 协议服务器（工具：add_score / batch_add_score / query_score / search_student / get_leaderboard / get_dashboard / send_notice），AstrBot、Claude Desktop 等 MCP 宿主配 `LEARNSTAR_API_BASE` + `LEARNSTAR_API_TOKEN` 两个环境变量即可接入——把 6.1 中机器人账号的 token 填进去，QQ/微信群即可用自然语言加减分。部署详见 [mcp-server/README.md](mcp-server/README.md)。

### 6.3 第三方平台登录与通讯录同步

学宠星球支持多平台第三方扫码登录。**管理员在后台勾选哪些平台启用**，勾选后登录页才会显示对应平台入口，并自动展示各平台官方品牌图标。

#### 6.3.1 后台启用第三方平台

1. 管理员登录 → **学校设置** → 「第三方登录平台」
2. 勾选要启用的平台（可多选）：**企业微信 / 钉钉 / 飞书 / 人人通空间 / 微信 / QQ**
3. 保存后，登录页教师端自动显示勾选的平台入口

> 未勾选任何平台时，默认启用 **企业微信 / 微信 / QQ**。
> 勾选平台的展示与登录无需额外配置；但**扫码登录真正可用**需在对应平台开放平台注册应用并配置凭证（见下表）。

#### 6.3.2 各平台配置凭证

| 平台 | 是否需要凭证 | 凭证位置 | 支持能力 |
|------|:-----------:|----------|----------|
| **企业微信** | ✅ | `.env`：`WECHAT_WORK_CORPID` / `WECHAT_WORK_AGENTID` / `WECHAT_WORK_SECRET` | 扫码免注册登录 + 通讯录导入 + 请假同步 + 消息推送 |
| **钉钉** | ✅ | `backend/config/dingtalk.php`：`app_key` / `app_secret` | 扫码登录 + 通讯录导入 |
| **飞书** | ✅ | `backend/config/feishu.php`：`app_id` / `app_secret` | 扫码登录 + 通讯录导入 |
| **微信** | ⚠️ | 需前端接入微信开放平台 JSAPI 取 openid | 当前仅展示入口，扫码流程待接入 |
| **QQ** | ⚠️ | 需前端接入 QQ 互联取 openid | 当前仅展示入口，扫码流程待接入 |
| **人人通空间** | ⚠️ | 需对接区域人人通开放平台 | 当前仅展示入口，扫码流程待接入 |

#### 6.3.3 教师绑定与免注册登录

- **企业微信 / 钉钉 / 飞书**：教师扫码后，若系统已有绑定账号则直接登录；否则**按手机号/实名用户名匹配本地已有账号并自动绑定**（不会重复建号）；确无匹配时自动创建教师账号（实名用户名 + 默认密码 `ls123456`）
- **账号设置 → 第三方账号绑定**：教师可查看/解绑已绑定的第三方平台，或扫码绑定新平台

#### 6.3.4 通讯录批量导入

管理员 → 教师管理 → **🏢 第三方导入**：从学校配置的第三方平台（企业微信/钉钉/飞书）拉取通讯录，勾选成员后批量创建教师与学生账号。教师按手机号/实名账号自动去重（不会生成"张老师_2"冗余账号）；学生按部门名自动匹配班级（支持"六年级1班 ↔ 六年级（1）班"模糊匹配），可搜索、批量设置班级，导入完成后展示新增/跳过明细与新教师初始密码。

#### 6.3.5 数据同步与冲突处理（升班 / 名单更新必读）

系统以**本地数据库为准**，第三方平台（企业微信/钉钉/飞书）是数据来源之一。平台侧的调整不会自动写入本地，按以下规则协同：

| 场景 | 系统的处理 |
|------|-----------|
| 通讯录重复导入 | 教师按手机号/实名账号自动跳过已有账号，**不会**生成"张老师_2"这类冗余账号；学生按"同班同名"跳过 |
| 学生转班（平台名单已调整） | Excel/通讯录导入时，**同学号已在其他班级会被拦截**并提示用「批量转班」处理；跨班同名会给出提醒供管理员判断（同名不同人可直接忽略） |
| 教师未导入过、直接扫码登录 | 系统按手机号 → 实名用户名匹配本地已有账号并自动绑定，**不会**重复建号 |
| 学年升班 | 在本系统「学年升级」执行（六年级毕业、班级整体升级）；**第三方平台的部门/名单需要平台管理员同步调整**，之后再回来导入通讯录并核对班级映射 |

> 推荐顺序：每个学年开始时，先在第三方平台完成升班与名单调整 → 本系统执行「学年升级」→ 建新一年级班级 → 通讯录导入新生（同学号冲突会被自动拦截）。

#### 6.3.6 工作原理（可选了解）

- 登录页平台列表来自接口 `GET /api/v1/auth/third-party/options`（返回管理员勾选的平台 + 品牌图标）
- 学校配置存储在 `schools.settings.enabled_third_party_platforms`（JSON 数组）
- 扫码回调按学校配置的平台分发（`App\Services\ThirdParty\ThirdPartyManager`），多校部署时通过 OAuth `state` 参数区分学校

## 7. 常见问题（FAQ）

### 7.1 使用问题

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
在管理后台「AI 中心」可视化配置，或 `.env` 配置 `AI_PROVIDER` 和 `AI_API_KEY`，支持 DeepSeek / OpenAI / 通义千问 / Moonshot 等 30+ 供应商。不配置不影响其他功能。
</details>

<details>
<summary>学生怎么查看自己的数据？</summary>
学生无需单独登录，通过班级码进入教室端（/classroom）即可查看本班积分、排行榜、宠物与图鉴。
</details>

<details>
<summary>有用户数限制吗？</summary>
无。MIT 开源协议不限制用户数和班级数。全校使用建议外接你已有的 MySQL 服务器。
</details>

<details>
<summary>别的系统/机器人能对接积分管理吗？</summary>
能。系统内置拥有全部班级权限的 API 机器人账号，通过 REST API 即可加减分/查分/排行；也提供标准 MCP 服务器对接 AstrBot、Claude Desktop 等宿主。见 <a href="#6-集成与对接">集成与对接</a>。
</details>

### 7.2 部署故障排查（新手向）

| 现象 | 原因与解法 |
|------|-----------|
| `up` 时报端口被占用 | `APP_PORT` 改成 `8081` 等未占用端口，`APP_URL` 同步改，重启 |
| 容器一直 `restarting` 或 `unhealthy` | `docker-compose logs app` 看最后 50 行日志；多为 `.env` 数据库段改错（两段同时生效或格式错误） |
| 镜像拉取超时（国内网络） | 给 Docker 配置镜像加速器；或 fork 仓库用 Actions 自行构建，`GITHUB_USERNAME` 改成你的用户名（小写） |
| 手机打不开系统 | `APP_URL` 改成部署电脑的局域网 IP（如 `http://192.168.1.100`），防火墙放行 `APP_PORT` |
| 忘记管理员密码 | 改 `.env` 的 `ADMIN_PASSWORD` → `docker-compose up -d` 重启即同步 |
| 想彻底重置 | `docker-compose down -v`（⚠️ **删除全部数据**，包括学生积分），再 `up -d` 从零开始 |

## 8. 技术架构

### 8.1 技术栈

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

### 8.2 目录结构

```
learnstar-planet/
├── frontend-vue/          # Vue 3 SPA（管理员端 / 教师端 / 教室端）
├── backend/               # Laravel 12 API（24 模型，220 端点）
├── mini-program/          # 微信小程序（10 页面，教师端）
├── pwa/                   # PWA 离线配置
├── mcp-server/            # MCP 机器人服务
└── docker-compose.yml     # Docker 编排（默认仅 app 单容器）
```

## 9. 参与贡献

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/my-feature`
3. 提交变更：`git commit -am 'feat: add my feature'`
4. 推送：`git push origin feature/my-feature`
5. 提交 Pull Request

代码规范：PHP PSR-12（PHP CS Fixer）· TypeScript ESLint · Conventional Commits

## 10. 许可证（License）

[MIT License](LICENSE) © 2024 RealKiro
