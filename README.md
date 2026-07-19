<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  <img src="https://img.shields.io/github/stars/RealKiro/learnstar-planet?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/badge/PHP-8.3-777BB4?logo=php" alt="PHP 8.3">
  <img src="https://img.shields.io/badge/Laravel-11-F9322C?logo=laravel" alt="Laravel 11">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" alt="Docker">
</p>

# 🌌 学趣星球 (LearnStar Planet)

<h3 align="center">开源班级管理系统 · 积分激励 · 宠物养成 · AI 助教</h3>

<p align="center">
  <b>全免费</b> · <b>自托管</b> · <b>多端适配</b> · <b>MIT 开源</b>
</p>

---

## ✨ 快速体验

```bash
# 一分钟部署
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

浏览器打开 `http://localhost:8080`，输入班级码进入教室端，或使用管理员账号登录。

---

## 📋 目录

- [核心功能](#-核心功能)
- [快速部署](#-快速部署)
- [系统架构](#-系统架构)
- [前后端联调](#-前后端联调)
- [开发指南](#-开发指南)
- [FAQ](#-faq)
- [贡献指南](#-贡献指南)

---

## 🎯 核心功能

### 🏠 班级总览
| 功能 | 说明 |
|------|------|
| 班级之星 | 展示积分最高的学生及其宠物 |
| 实时动态 | 全班最新加分/进化消息滚动 |
| TOP5 排行 | 前五名学生积分排名与进度条 |

### ✏️ 课堂评价
| 功能 | 说明 |
|------|------|
| 积分加减 | 卡片式学生网格，点击 +/− 选择行为原因 |
| 自定义步长 | 点击中间数字可编辑加减分值 |
| 快捷规则 | 一键全班按规则加减分 |
| 宠物切换 | 学生自选宠物（首次免费，后续扣20积分） |

### 🏆 年级战场
| 功能 | 说明 |
|------|------|
| 跨班 PK | 同年级各班积分排行榜（进度条 + 详细数据） |
| 本班战力 | 总积分、平均等级、巅峰人数、周增长 |
| 发起挑战 | 班级间 PK 挑战机制 |

### 📚 宠物图鉴
| 功能 | 说明 |
|------|------|
| 54 物种 | 8 大系列、54 种宠物、12 级进化 |
| 专属诗文 | 每物种 6 阶段专属七言律诗 + 进化台词 |
| 系列切换 | 全班投票选择系列（每人扣 20 积分） |

### 🔧 高级功能（教师账号）
| 功能 | 说明 |
|------|------|
| 实时广播 | 顶部横幅/弹窗/全屏三种广播模式 |
| 智能考勤 | 一键点名，4 种考勤状态 |
| 扫码收作业 | 生成二维码，学生扫码提交 |
| 在线答题 | 题库管理，自动判分 |
| 成绩管理 | 成绩录入、统计分析 |
| AI 助教 | 集成通义千问/Qwen/DeepSeek |
| 积分商城 | 学生兑换奖品 |
| 多币种 | 积分/金币等多币种支持 |

---

## 🚀 快速部署

### 方式一：Docker 一键部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env，修改 APP_URL 和数据库密码（可选）

# 3. 启动服务
docker-compose up -d

# 4. 查看运行状态
docker-compose ps
```

启动后访问：
- **教室端**：`http://localhost:8080` — 输入班级码（如 `LS301`）
- **教师端**：`http://localhost:8080/login` — 管理员/教师登录

### 方式二：手动部署

**后端（Laravel 11）：**
```bash
cd backend
composer install
php artisan migrate
php artisan db:seed --class=AdminUserSeeder
php artisan serve
```

**前端（Vue 3）：**
```bash
cd frontend-vue
npm install
npm run dev       # 开发服务器 http://localhost:5173
npm run build     # 生产构建
```

---

## 🏗 系统架构

```
learnstar-planet/
├── frontend-vue/          # Vue 3 + Vite + TypeScript SPA
│   └── src/
│       ├── pages/         # 页面组件
│       │   ├── teacher/   # 教师端（完整功能）
│       │   ├── classroom/ # 教室端（班级码登录，4大核心模块）
│       │   ├── parent/    # 家长端
│       │   └── landing/   # 首页/登录页
│       ├── layouts/       # 布局组件
│       ├── components/    # 通用组件
│       ├── services/      # API 服务层
│       ├── stores/        # Pinia 状态管理
│       └── utils/         # 工具函数
├── backend/               # Laravel 11 API
│   ├── app/
│   │   ├── Models/        # 17 个 Eloquent 模型
│   │   ├── Http/Controllers/Api/  # API 控制器
│   │   ├── Services/      # 业务服务层
│   │   └── Http/Requests/ # 表单验证
│   └── routes/api.php     # ~130 个 API 端点
├── mini-program/          # 微信小程序
├── pwa/                   # PWA 配置
├── mcp-server/            # MCP AI 机器人服务器
└── docker-compose.yml     # Docker 编排
```

---

## 🔗 前后端联调

### API 路径规范

```
/api/v1/{role}/{resource}/{action}
```

| 角色 | 前缀 | 认证 | 适用场景 |
|------|------|------|---------|
| 认证 | `/auth` | 部分需要 | 登录/登出/第三方绑定 |
| 教师 | `/teacher` | Bearer Token | 全部管理功能 |
| 管理员 | `/admin` | Bearer Token | 学校级管理 |
| 家长 | `/parent` | Bearer Token | 查看孩子数据 |
| 大屏 | `/display` | 班级码 Token | 教室端 4 大模块 |
| 公共 | `/common` | 无 | 宠物类型等公开数据 |

### 权限分级

| 登录方式 | 访问范围 | 侧边栏 |
|---------|---------|--------|
| **班级码** | 4 大核心模块 | 教室端侧边栏 |
| **教师账号** | 全部功能 | 教师端侧边栏 |

> 📖 完整 API 文档见 [`docs/api-reference.md`](docs/api-reference.md)

---

## 🛠 开发指南

### 常用命令

```bash
# 前端
cd frontend-vue
npm run dev         # 开发模式
npm run typecheck   # TypeScript 类型检查
npm run build       # 生产构建

# 后端
cd backend
php artisan test                    # 运行测试
vendor/bin/phpstan analyse          # PHPStan 静态分析
vendor/bin/php-cs-fixer fix         # CS Fixer
```

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `APP_DEBUG` | `false` | 调试模式 |
| `INITIALIZATION_SETTINGS` | `false` | 测试模式（与 DEBUG 配合重置积分） |
| `DB_CONNECTION` | `mysql` | 数据库类型 |
| `AI_PROVIDER` | 空 | AI 服务商 |

---

## ❓ FAQ

**Q: 如何获取班级码？**  
A: 管理员登录后台 → 班级管理 → 查看/刷新班级码。

**Q: 学生如何登录？**  
A: 首页输入班级码即可进入教室端，无需账号密码。

**Q: 如何切换数据库类型？**  
A: 修改 `.env` 中的 `DB_CONNECTION`，参考 `.env.example` 中的注释。

**Q: AI 功能如何使用？**  
A: 配置 `AI_PROVIDER` 和 `AI_API_KEY`，支持 qwen/openai/deepseek/moonshot。

---

## 🤝 贡献指南

欢迎贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支: `git checkout -b feature/my-feature`
3. 提交变更: `git commit -am 'feat: add my feature'`
4. 推送: `git push origin feature/my-feature`
5. 提交 Pull Request

### 代码规范
- PHP: PSR-12
- TypeScript: ESLint
- 提交信息: Conventional Commits

---

## 📄 许可证

[MIT License](LICENSE) © 2024 RealKiro

> 全免费、开源、可自由使用和修改
