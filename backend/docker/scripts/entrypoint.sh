#!/bin/sh
# ============================================================
# 学趣星球 - Docker 入口脚本
# 初始化 Laravel 应用、运行迁移、启动 PHP 内置服务器
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

# 从 Docker 环境变量创建 .env 文件
if [ ! -f .env ]; then
    echo "📝 从环境变量创建 .env 文件..."
    env | grep -E "^(APP_|DB_|REDIS_|CACHE_|SESSION_|QUEUE_|MAIL_|FILESYSTEM_|AI_|WECHAT_|QQ_|RENREN_|ADMIN_|GITHUB_)" > .env
fi

# 自动拼接端口到 APP_URL
if [ -n "$APP_URL" ] && [ -n "$APP_PORT" ]; then
    if ! echo "$APP_URL" | grep -qE ':[0-9]+(/.*)?$'; then
        case "$APP_PORT" in
            80|443) ;;
            *)
                APP_URL="${APP_URL}:${APP_PORT}"
                export APP_URL
                if [ -f .env ]; then
                    sed -i "s|^APP_URL=.*|APP_URL=${APP_URL}|" .env
                fi
                echo "  APP_URL 自动拼接端口: ${APP_URL}"
                ;;
        esac
    fi
fi

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

    # 生成应用密钥
    if ! grep -q "^APP_KEY=base64:" .env 2>/dev/null; then
        echo "🔑 生成应用密钥..."
        NEW_KEY=$(php -r "echo 'base64:' . base64_encode(random_bytes(32));" 2>/dev/null)
        if [ -z "$NEW_KEY" ]; then
            NEW_KEY=$(php artisan key:generate --show 2>/dev/null)
        fi
        if [ -n "$NEW_KEY" ]; then
            sed -i "/^APP_KEY=/d" .env
            echo "APP_KEY=$NEW_KEY" >> .env
            export APP_KEY="$NEW_KEY"
            echo "✅ 密钥已生成"
        else
            echo "⚠️  密钥生成失败，请手动设置 APP_KEY"
        fi
    else
        echo "✅ 密钥已存在"
    fi

    # 运行迁移
    echo "📦 运行数据库迁移..."
    php artisan migrate --force
    echo "✅ 迁移完成"

    # 创建管理员账号（首次部署）
    echo "👤 创建默认管理员..."
    php artisan db:seed --class=AdminUserSeeder --force 2>/dev/null || true
    echo "✅ 管理员账号已就绪"

    # 创建存储链接
    echo "🔗 创建存储链接..."
    php artisan storage:link --force 2>/dev/null || true

    fi

# 优化缓存
echo "🧹 优化缓存..."
php artisan config:cache 2>/dev/null || true
php artisan route:cache 2>/dev/null || true
php artisan view:cache 2>/dev/null || true
echo "✅ 缓存重建完成"

echo "🎉 学趣星球初始化完成！"

# 如果使用 Redis 队列，启动后台 queue worker
if [ "${QUEUE_CONNECTION:-redis}" = "redis" ] || [ "${QUEUE_CONNECTION:-redis}" = "database" ]; then
    echo "⏳ 启动队列 worker（${QUEUE_CONNECTION:-redis}）..."
    php artisan queue:work --queue=default --sleep=3 --tries=3 --timeout=60 &
    echo "✅ 队列 worker 已启动"
else
    echo "⚡ 队列模式: sync（同步处理）"
fi

# 启动 PHP 内置服务器
echo "  启动 PHP 内置服务器（http://0.0.0.0:8080）..."
exec php artisan serve --host=0.0.0.0 --port=8080
