# 🚀 部署指南

## 环境要求

| 组件 | 版本 |
|------|------|
| Docker & Docker Compose | 最新版 |
| PHP | 8.3+（手动部署） |
| Composer | 2.x（手动部署） |
| Node.js | 22+（手动部署） |
| MySQL / PostgreSQL / SQLite | 任选其一 |

## 方式一：Docker 一键部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env，至少修改 APP_URL 和数据库密码

# 3. 启动（后台运行）
docker-compose up -d

# 4. 查看日志
docker-compose logs -f app

# 5. 停止
docker-compose down
```

### 端口说明

| 服务 | 端口 | 默认值 |
|------|------|--------|
| Web 应用 | APP_PORT | 8080 |
| MySQL | MYSQL_PORT | 3306 |
| Redis | REDIS_PORT | 6379 |

## 方式二：手动部署

### 后端

```bash
cd backend
composer install
cp .env.example .env
# 编辑 .env 配置数据库
php artisan key:generate
php artisan migrate
php artisan db:seed --class=AdminUserSeeder
php artisan serve
```

### 前端

```bash
cd frontend-vue
npm install
npm run dev       # 开发模式
npm run build     # 生产构建
```

## 测试模式

