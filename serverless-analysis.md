# 学趣星球 - 数据库必要性 & 无服务器部署分析

> 分析日期: 2026-07-10

---

## 一、项目现状摘要

| 维度 | 现状 |
|---|---|
| **后端框架** | Laravel 11 (PHP 8.3) |
| **数据库表** | 23 张表，20+ 外键约束 |
| **API 端点** | ~130 个 RESTful 端点 |
| **前端** | 单文件 HTML (245KB) + 微信小程序 + PWA |
| **缓存/队列** | Redis（ZSET 排行榜、Session、Cache、Queue、Broadcast） |
| **认证** | Sanctum Token + bcrypt 密码哈希 |
| **当前部署** | Docker Compose (app + MySQL + Redis) |

---

## 二、数据库必要性分析：结论：**必须，不可省略**

### 2.1 数据关系复杂度

23 张表构成的实体关系图：

```
School (1) ──→ (N) ClassRoom (1) ──→ (N) Student (1) ──→ (1) Pet
  │                  │                      │
  │                  ├──→ ScoreRule          ├──→ Score (N)
  │                  ├──→ Notice             ├──→ ScoreLog (N)
  │                  ├──→ ShopItem           ├──→ Attendance (N)
  │                  ├──→ Broadcast          ├──→ Grade (N)
  │                  ├──→ HomeworkCollection ├──→ ShopRedemption (N)
  │                  └──→ Quiz               └──→ HomeworkSubmission (N)
  │
  └──→ User (teacher/parent/school_admin)
         └──→ ThirdPartyBinding (N)
```

### 2.2 需要数据库的核心场景

| 场景 | 涉及表 | 是否需要事务？ |
|---|---|---|
| 教师给学生加分 | scores, students, score_logs, pets | ✅ `DB::transaction` 保证原子性 |
| 批量创建教师账号 | users, schools | 需要唯一性检查（username/nickname） |
| 积分排行榜查询 | students, score_logs + Redis ZSET | 需要 JOIN 查询 + 聚合排序 |
| 学校管理后台统计 | 跨 10+ 张表的 COUNT/SUM/AVG 聚合 | LEFT JOIN + 子查询 |
| 第三方扫码登录绑定 | users, third_party_bindings | 需要唯一约束防重复绑定 |
| 宠物进化升级 | pets, students (total_score 作为经验来源) | ✅ 事务 |
| 考勤签到 | attendances | 唯一约束 (class_id, student_id, date) |

### 2.3 为什么不能用 KV/对象存储/Document DB 替代

- **外键约束**: 删除学校需级联处理所有班级和学生，关系型数据库 native 支持
- **复合唯一索引**: 如 `(class_id, student_id, date)` 防重复签到，`(platform, platform_id)` 防重复绑定
- **JOIN 查询**: SchoolAdminController 中大量 LEFT JOIN + 子查询（如教师列表 JOIN bindings 表）
- **聚合查询**: `COUNT(DISTINCT ...)`, `SUM()`, `AVG()` 跨表统计
- **事务一致性**: 积分变更同时修改 4 张表，必须原子操作

---

## 三、Cloudflare 无服务器部署可行性分析

### 3.1 Cloudflare Workers / Pages → ❌ 不适合

| 阻碍因素 | 详细说明 |
|---|---|
| **PHP 运行时** | Workers 仅支持 JavaScript/WebAssembly，不运行 PHP。Laravel 无法在 Workers 上运行 |
| **数据库** | Cloudflare D1 是 SQLite，最多支持 23 张表没问题，但缺少 MySQL 的一些特性（如锁粒度、并发写入能力） |
| **Redis** | 无 Redis 替代。排行榜 ZSET、会话、队列、广播均依赖 Redis。可用 KV 模拟部分功能但无法做 ZSET 排序 |
| **WebSocket/Livewire** | Livewire 需要持久连接，Workers 是无状态函数，不适合 |
| **队列/后台任务** | Laravel Horizon 需要长期运行的进程，Workers 有 30s CPU 限制 |
| **文件上传** | 需要 R2 替代本地存储 |
| **复杂度** | 将 19K+ 行 PHP 代码重写为 JS 不现实 |

### 3.2 Cloudflare 可以做什么 → ✅ 部分可用

| 组件 | 迁移到 Cloudflare | 可行性 |
|---|---|---|
| **前端 index.html** | → Cloudflare Pages | ✅ 完美！免费 CDN，全球加速 |
| **PWA + Service Worker** | → Cloudflare Pages | ✅ |
| **静态资源** (CSS/JS/图片) | → Cloudflare Pages/R2 | ✅ |
| **域名 & DNS** | → Cloudflare DNS | ✅ 免费的 DDoS 防护 |
| **API 网关/代理** | → Cloudflare Workers 做反向代理 | ✅ 可在 Workers 上做简单路由/鉴权转发到后端 |
| **文件上传存储** | → Cloudflare R2 (S3 兼容) | ✅ 无出口流量费 |

### 3.3 可实现的混合架构

```
                    ┌────────────────────┐
  用户浏览器  ─────→│ Cloudflare Pages   │ ← 托管 index.html (免费)
                    │ (前端静态托管)      │
                    └────────┬───────────┘
                             │ API 请求
                    ┌────────▼───────────┐
                    │ Cloudflare Workers │ ← API 反向代理/DNS (免费)
                    │ (轻量级代理层)      │
                    └────────┬───────────┘
                             │ 转发
                    ┌────────▼───────────┐
                    │ 后端 API 服务器     │ ← 廉价 VPS / Fly.io
                    │ Laravel + SQLite   │
                    │ + File Cache       │
                    └────────────────────┘
```

---

## 四、实际推荐方案

### 方案 A：SQLite 极简部署（最省钱，推荐小型使用）

```
成本: ¥0/月（Fly.io 免费额度 或 最便宜 VPS ¥30/月）
```

| 组件 | 方案 |
|---|---|
| 前端 | Cloudflare Pages (免费) |
| 后端 | 单容器运行 Laravel + SQLite (无需 MySQL/Redis) |
| 数据库 | SQLite 文件数据库 (storage/database.sqlite) |
| 缓存 | file driver |
| 会话 | file driver |
| 队列 | database driver |
| 排行榜 | MySQL/Redis → 改为直接从 scores 表实时查询 |
| 限制 | <500 学生，排行榜查询需实时 SQL |

**具体做法**: 修改 `.env`:
```env
DB_CONNECTION=sqlite
CACHE_DRIVER=file
SESSION_DRIVER=file
QUEUE_CONNECTION=database
BROADCAST_DRIVER=null
```

**代码改动**：LeaderboardService 需要从 Redis ZSET 改为 SQL 实时查询（已有 MySQL 回退逻辑 `rebuildFromClassDB`）。

### 方案 B：Cloudflare Pages + Fly.io（零成本起步）

```
成本: ¥0/月（Fly.io 免费额度：3个共享 VM，3GB 持久卷）
```

- Fly.io 原生支持 Docker 容器部署
- SQLite 配合 Fly.io 持久卷 (Volume)，数据不丢失
- 前端放 Cloudflare Pages，后端放 Fly.io

### 方案 C：Cloudflare Pages + 廉价 VPS（小规模稳定）

```
成本: ¥30-60/月
```

- 前端：Cloudflare Pages（免费）
- 后端：阿里云/腾讯云轻量服务器
- 数据库：同 VPS 上运行 MySQL 容器 或直接用 SQLite
- 这也是 README 中已有的推荐

### 方案 D：全量 Docker Compose（当前方案，功能最全）

```
成本: ¥50-60/月
```

- Docker Compose 一键部署 app + MySQL + Redis
- Redis 排行榜、队列、广播全部可用
- 适合 >500 学生或需要全部功能

---

## 五、结论

### 数据库必要性

**必须使用关系型数据库。** 23 张表、复杂外键关系、事务性操作、聚合查询、唯一约束 —— 这些都不是 KV 存储或文档数据库能替代的。SQLite 可以作为小型部署的替代品（<500 学生）。

### 无服务器可行性

**无法完全迁移到 Cloudflare Workers/Pages。** Laravel/PHP 生态与 Cloudflare 的 JS 运行时格格不入。但可以采用混合架构：前端放 Cloudflare Pages，后端放廉价 VPS/Fly.io 用 SQLite，实现接近零成本的运行。

### 最推荐的"类无服务器"方案

| 层级 | 方案 | 月费 |
|---|---|---|
| 前端 | Cloudflare Pages | ¥0 |
| 后端 | Fly.io 免费额度 (Docker + SQLite) | ¥0 |
| 域名 | Cloudflare DNS | ¥0 |
| **总计** | | **¥0/月** |

这是在不重写代码的前提下，最接近"无服务器"的部署方式。
