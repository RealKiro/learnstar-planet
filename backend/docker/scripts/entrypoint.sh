#!/bin/sh
# ============================================================
# 学趣星球 - Docker 入口脚本
# 初始化 Laravel 应用、运行迁移、启动服务
# ============================================================

set -e

echo "🚀 学趣星球 - 启动初始化..."

# 确保存储目录存在
mkdir -p storage/app/public \
    storage/framework/cache/data \
    storage/framework/sessions \
    storage/framework/views \
    storage/logs \
    bootstrap/cache

# 等待数据库就绪
echo "⏳ 等待数据库连接..."
MAX_RETRIES=30
RETRY_COUNT=0
until php artisan db:monitor 2>/dev/null || php artisan migrate:status 2>/dev/null; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        echo "⚠️  数据库连接超时（${MAX_RETRIES}次重试），跳过迁移继续启动..."
        break
    fi
    echo "  数据库未就绪，3秒后重试（${RETRY_COUNT}/${MAX_RETRIES}）..."
    sleep 3
done

if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo "✅ 数据库连接成功"

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

    # 创建默认管理员账号（首次部署）
    echo "👤 创建默认管理员..."
    php artisan db:seed --class=AdminUserSeeder --force
    echo "✅ 管理员账号已就绪"

    # 创建存储目录链接
    echo "🔗 创建存储链接..."
    php artisan storage:link --force 2>/dev/null || true
fi

# 清除并重建缓存
echo "🧹 清除缓存..."
php artisan config:cache 2>/dev/null || true
php artisan route:cache 2>/dev/null || true
php artisan view:cache 2>/dev/null || true
echo "✅ 缓存重建完成"

echo "🎉 学趣星球初始化完成！"
echo "  应用地址: http://localhost:${APP_PORT:-8080}"
echo "  API地址: http://localhost:${APP_PORT:-8080}/api"

# 启动 supervisor（管理 nginx + php-fpm + queue worker）
exec "$@"
