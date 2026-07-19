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
  <b>教室大屏 + 积分激励 + 宠物进化 + AI助教</b><br>
  一台普通电脑就能跑，一分钱不用花
</p>

<p align="center">
  <b>局域网就能用</b> · <b>学生不用注册</b> · <b>完全免费开源</b>
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

- [🎯 这是个什么系统？](#-这是个什么系统)
- [👩‍🏫 我在课堂上能用它做什么](#-我在课堂上能用它做什么)
- [🏫 为什么适合我们学校](#-为什么适合我们学校)
- [🚀 部署方式](#-部署方式)
- [🏗 系统架构](#-系统架构)
- [❓ 常见问题](#-常见问题)

---

## 🎯 这是个什么系统？

**学趣星球是一套教室大屏互动系统。** 老师办公室有台旧电脑、教室有台能投屏的电视/投影仪，就能用。

它的核心逻辑很简单：**学生做了好事→加分→宠物升级；做错了→扣分→宠物饿肚子。** 不用买服务、不用联网、学生连账号都不用注册。

学生只需要在教室大屏幕上输入 4-6 位**班级码**，就能看到自己的积分、宠物、排名。家长用手机打开浏览器也能看到孩子的在校表现。

---

## 👩‍🏫 我在课堂上能用它做什么

### 📡 实时广播——把消息推到教室大屏上

> 课堂里最常用的功能。不需要喊"安静"、不需要拍桌子，直接发个弹窗全班都看到。

| 方式 | 效果 |
|------|------|
| **顶部横幅** | 屏幕上方滚动显示通知文字，不打断课堂进度 |
| **居中弹窗** | 全屏弹窗，适合紧急通知或重要提醒 |
| **全屏覆盖** | 完全切换大屏画面，适合展示特定内容 |

教师手机上操作 → 教室大屏立刻显示（<200ms 延迟），不需要额外装任何软件。

### 🖥️ 教室大屏模式——投影仪/电视一接就能用

浏览器全屏打开，就是一间智能教室：

| 功能 | 上课时怎么用 |
|------|------------|
| **班级总览** | 一进教室投屏，全班积分之星、最近动态、TOP5 排名一目了然 |
| **课堂评价** | 点击学生头像 +/− 加分扣分，选行为原因（答对题、作业优秀、违反纪律），大屏实时更新 |
| **年级战场** | 同年级各班排行榜 PK，进度条 + 数据，学生集体荣誉感拉满 |
| **宠物图鉴** | 54 种宠物、12 级进化路线，学生随时查看自己的宠物和进化进度 |

> 加分时宠物会弹跳、扣分时会抖动，学生有即时反馈。

### 📝 课堂管理工具

| 功能 | 说明 |
|------|------|
| **智能考勤** | 大屏上一键点名，标记出勤/迟到/请假/缺席，自动记录到班级台账 |
| **扫码收作业** | 生成二维码投到大屏，学生扫码提交，谁交了谁没交一目了然 |
| **在线测验** | 题库里选题出卷，学生扫码作答，系统自动判分直接进成绩册 |
| **积分商城** | 老师设置兑换规则（比如"免一次作业"需要 100 分），学生自助兑换 |

### 🏆 积分与激励体系——学生最吃这套

| 机制 | 说明 |
|------|------|
| **积分（成长值）** | 加一分就是一点成长值，全班公开，每次变动都有实时推送 |
| **宠物进化** | 累计成长值喂养宠物，54 种宠物 × 8 大系列 × 12 级进化，每级有专属名字和诗文 |
| **跨班 PK** | 同年级班级之间 PK 总积分/平均等级/周增长，可发起挑战 |
| **排行榜** | 总积分榜 / 每周增长榜 / 宠物等级榜，大屏轮播 |

> 学生最在意的是"我的宠物什么时候进化"——这是最好的内驱力。

### 🤖 AI 助教（可选）

配置 AI 接口后，老师可以直接让 AI 辅助出题、批改、答疑。不想配也不影响任何核心功能。

### 📱 家长端

家长用微信小程序或手机浏览器就能看：孩子当前积分、宠物状态、班级排名、老师发的通知。

---

## 🏫 为什么适合我们学校

### 不花一分钱

- **MIT 开源协议**：没有任何隐藏收费、没有用户数限制、没有功能锁定
- **可以装在旧电脑上**：办公室淘汰的电脑就能跑，不需要买服务器
- **数据库也免费**：支持 SQLite，连 MySQL 都不用装

### 不需要联网

- 系统全部部署在**学校局域网**内，教室电脑浏览器打开就能用
- 家长如果想看，连上学校 WiFi 或回家用手机也能访问
- 数据全程在学校内部，不出校门

### 一台电脑就能带动

| 部署方式 | 需要的硬件 | 难度 |
|---------|-----------|:----:|
| 有 Docker 的环境 | 任何能装 Docker 的电脑（Windows/Mac/Linux） | ⭐ 简单 |
| 只有 PHP 的环境 | 只要有 PHP 8.3，Windows 电脑也行 | ⭐⭐ 中等 |
| SQLite 省到底 | 一台旧电脑，不需要 MySQL，不需要 Redis | ⭐⭐ 中等 |

更重要的是——**学生不需要账号**。给你所教的每个班生成一个班级码（4-6 位数字+字母），学生打开大屏输入班级码就能进去，没有注册流程、不会忘记密码。

---

## 🚀 部署方式

### 对比速览

| 方式 | 需要什么 | 适用于 |
|------|---------|--------|
| **① Docker 部署** | 能装 Docker 的电脑 | 有 Docker 环境，一步到位 |
| **② SQLite 极简** | PHP 8.3 即可 | 一台旧电脑，不想装 MySQL/Redis |
| **③ LAMP 标准部署** | PHP + Nginx + MySQL | 学校本身有服务器 |

### 方式一：Docker 一键部署（推荐）

```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

启动后访问：
- **教室大屏**：浏览器打开 `http://你的电脑IP:8080`，输入班级码（如 `LS301`）
- **教师管理**：`http://你的电脑IP:8080/login`，管理员或教师账号登录

> 全校师生都能通过局域网 IP 访问，不需要互联网。

### 方式二：SQLite 极简部署（零成本）

<details>
<summary><b>📋 展开步骤——不需要 Docker、MySQL、Redis</b></summary>

```bash
# 1. 安装后端
cd backend
composer install --no-dev --optimize-autoloader
cp .env.example .env
```

编辑 `backend/.env`：

```env
DB_CONNECTION=sqlite       # 用文件存数据，不需要装 MySQL
CACHE_DRIVER=file
SESSION_DRIVER=file
QUEUE_CONNECTION=database
BROADCAST_DRIVER=null
REDIS_HOST=                 # 禁掉 Redis
```

```bash
# 2. 初始化
touch storage/database.sqlite
php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder

# 3. 构建前端
cd ../frontend-vue
npm install && npm run build
cp -r dist/* ../backend/public/

# 4. 启动
cd ../backend
php artisan serve --host=0.0.0.0 --port=8080
```

> **一台 Windows 电脑也能跑**：装 PHP 8.3 和 Composer，然后执行上面的命令。

</details>

### 方式三：传统 LAMP/LEMP 部署

<details>
<summary><b>📋 展开步骤——学校有服务器的情况</b></summary>

```bash
cd backend
composer install --no-dev --optimize-autoloader
cp .env.example .env
# 编辑 .env：配置数据库连接
php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder

cd ../frontend-vue
npm install && npm run build
```

Nginx 配置参考：

```nginx
server {
    listen 80;
    server_name 你的IP或域名;
    root /var/www/learnstar-planet/backend/public;
    index index.php;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php8.3-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

</details>

### 三种数据库对比

| 对比项 | SQLite（推荐小规模） | MySQL | Redis |
|--------|:------------------:|:-----:|:-----:|
| 需要安装 | ❌ 不需要 | 需要 | 需要 |
| 建议学生数 | < 500 人 | 无限制 | — |
| 排行榜速度 | 常规 | 快 | 毫秒级 |
| 实时广播 | ❌ 不支持 | ✅ 支持 | ✅ 支持 |
| 适合场景 | 一间教室 / 单个年级 | 全校使用 | 配合 MySQL 加速 |

> 先用 SQLite 跑起来，以后人多了改两行配置就能升级到 MySQL。

---

## 🏗 系统架构

```
learnstar-planet/
├── frontend-vue/          # 前端界面（教室大屏 + 教师管理 + 家长端）
├── backend/               # 后端 API（17 个数据模型，~130 个接口）
├── mini-program/          # 微信小程序（家长用）
├── pwa/                   # 手机浏览器也能离线看
├── mcp-server/            # QQ/微信机器人（可选）
└── docker-compose.yml     # Docker 一键启动配置
```

<details>
<summary><b>📋 技术详情（展开查看）</b></summary>

| 层级 | 技术 |
|------|------|
| 后端框架 | Laravel 11（PHP 8.3） |
| 数据库 | SQLite / MySQL / MariaDB / PostgreSQL 任选 |
| 前端 | Vue 3 + TypeScript + Tailwind CSS |
| 小程序 | 微信原生 |
| 实时推送 | SSE 事件流（<200ms 延迟，断线自动重连） |
| AI | 通义千问 / DeepSeek / OpenAI（可选） |

</details>

---

## ❓ 常见问题

<details>
<summary><b>我没服务器，用自己的电脑行不行？</b></summary>
完全行。Windows 装 PHP 8.3 和 Composer，用 SQLite 模式跑起来，教室电脑浏览器访问你电脑的局域网 IP 就行。<b>推荐让学校电教老师帮忙部署。</b>
</details>

<details>
<summary><b>学生怎么登录？需要注册吗？</b></summary>
不需要注册。老师在管理后台给每个班生成一个班级码（比如 LS301），学生在教室大屏首页输入这 4-6 位码就能进去。
</details>

<details>
<summary><b>教室没联网能用吗？</b></summary>
可以。系统部署在学校局域网内，教室电脑、教师手机只要连同一个校园网就能访问。完全不需要上外网。
</details>

<details>
<summary><b>手机能操作吗？</b></summary>
可以。教师用手机浏览器打开管理页面，可以给学生加分、发广播、查看排行榜。手机上操作，大屏实时同步。
</details>

<details>
<summary><b>广播功能怎么用？</b></summary>
教师登录后台 → 实时广播，选择横幅/弹窗/全屏三种模式，输入内容发送。教室大屏会立刻弹出显示（<200ms 延迟）。适合：上课铃响发"准备上课"、自习课发通知、课间发活动提醒。
</details>

<details>
<summary><b>AI 功能必须配吗？</b></summary>
不是。AI 助教是可选的，不配任何 Key 也不影响加分、宠物、考勤、考试这些核心功能。
</details>

<details>
<summary><b>必须用 Redis 吗？</b></summary>
不是。小规模用 SQLite + 文件缓存就够了。等以后全校用了再加 Redis 提速。
</details>

<details>
<summary><b>有用户数限制吗？</b></summary>
没有。MIT 开源协议不限用户数，不限班级数，不限功能。唯一限制是服务器性能。
</details>

<details>
<summary><b>家长怎么看？</b></summary>
家长用手机浏览器打开学校给的访问地址，用家长账号登录即可查看：孩子积分、宠物、排名、通知。也可以部署微信小程序。
</details>

<details>
<summary><b>我想自己改代码，可以吗？</b></summary>
MIT 协议允许任意修改和二次发布。Fork 仓库改完自己用就行，不需要告诉我们。
</details>

---

## 📄 许可证

[MIT License](LICENSE) © 2024 RealKiro

> 全免费、开源、可自由使用和修改。如果你觉得有用，欢迎给个 Star ⭐
