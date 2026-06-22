#!/bin/bash
# ============================================================
# 班宠星球 - Docker 入口脚本
# 初始化 Laravel 应用、运行迁移、启动服务
# ============================================================

set -e

echo "🚀 班宠星球 - 启动初始化..."

# 等待数据库就绪
echo "⏳ 等待数据库连接..."
until php artisan db:monitor --database=${DB_CONNECTION} 2>/dev/null || \
      php artisan migrate:status 2>/dev/null; do
    echo "  数据库未就绪，3秒后重试..."
    sleep 3
done
echo "✅ 数据库连接成功"

# 等待 Redis 就绪
echo "⏳ 等待 Redis 连接..."
until php artisan redis:ping 2>/dev/null || \
      redis-cli -h ${REDIS_HOST} -p ${REDIS_PORT} ping 2>/dev/null; do
    echo "  Redis 未就绪，3秒后重试..."
    sleep 3
done
echo "✅ Redis 连接成功"

# 生成应用密钥（首次部署）
if [ -z "$APP_KEY" ]; then
    echo "🔑 生成应用密钥..."
    php artisan key:generate --force
    echo "✅ 密钥已生成"
fi

# 运行数据库迁移
echo "📦 运行数据库迁移..."
php artisan migrate --force
echo "✅ 迁移完成"

# 同步排行榜数据到 Redis
echo "📊 同步排行榜数据..."
php artisan leaderboard:sync --all
echo "✅ 排行榜同步完成"

# 创建存储目录链接
echo "🔗 创建存储链接..."
php artisan storage:link --force 2>/dev/null || true

# 清除并重建缓存
echo "🧹 清除缓存..."
php artisan config:cache
php artisan route:cache
php artisan view:cache
echo "✅ 缓存重建完成"

echo "🎉 班宠星球初始化完成！"
echo "  应用地址: http://localhost:${APP_PORT:-80}"
echo "  API地址: http://localhost:${APP_PORT:-80}/api"
echo "  Horizon面板: http://localhost:${APP_PORT:-80}/horizon"

# 执行传入的 CMD
exec "$@"
