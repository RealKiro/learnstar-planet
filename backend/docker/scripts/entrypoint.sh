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
    bootstrap/cache \
    /var/log/supervisor \
    /var/run

# ──────────────────────────────────────────
# 从 Docker 环境变量创建 .env 文件
# Laravel 的 artisan 命令（如 key:generate、migrate）需要读取 .env
# 但 Docker Compose 通过 environment: 注入变量，容器内没有 .env 文件
# ──────────────────────────────────────────
if [ ! -f /var/www/html/.env ]; then
    echo "📝 从环境变量创建 .env 文件..."
    env | grep -E "^(APP_|DB_|REDIS_|CACHE_|SESSION_|QUEUE_|MAIL_|FILESYSTEM_|AI_|WECHAT_|QQ_|RENREN_|ADMIN_|GITHUB_)" > /var/www/html/.env
fi

# ──────────────────────────────────────────
# 自动拼接端口到 APP_URL
# APP_URL 不带端口时，从 APP_PORT 自动补上（80/443 除外）
# 这样用户改端口只需改 APP_PORT，无需同步改 APP_URL
# ──────────────────────────────────────────
if [ -n "$APP_URL" ] && [ -n "$APP_PORT" ]; then
    # 检查 APP_URL 是否已带端口（URL 末尾的 :数字，可能后跟 /路径）
    if ! echo "$APP_URL" | grep -qE ':[0-9]+(/.*)?$'; then
        case "$APP_PORT" in
            80|443) ;;
            *)
                APP_URL="${APP_URL}:${APP_PORT}"
                export APP_URL
                if [ -f /var/www/html/.env ]; then
                    sed -i "s|^APP_URL=.*|APP_URL=${APP_URL}|" /var/www/html/.env
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

    # 生成应用密钥（首次部署）
    # 注意：不能用 php artisan key:generate（它会写 .env，但 Docker env var
    # APP_KEY="" 会覆盖 .env 的值）。改用 --show 输出密钥，直接写入 .env。
    if ! grep -q "^APP_KEY=base64:" /var/www/html/.env 2>/dev/null; then
        echo "🔑 生成应用密钥..."
        NEW_KEY=$(php -r "echo 'base64:' . base64_encode(random_bytes(32));" 2>/dev/null)
        if [ -z "$NEW_KEY" ]; then
            # fallback: 使用 artisan key:generate --show
            NEW_KEY=$(php artisan key:generate --show 2>/dev/null)
        fi
        if [ -n "$NEW_KEY" ]; then
            sed -i "/^APP_KEY=/d" /var/www/html/.env
            echo "APP_KEY=$NEW_KEY" >> /var/www/html/.env
            export APP_KEY="$NEW_KEY"
            echo "✅ 密钥已生成: ${NEW_KEY:0:16}..."
        else
            echo "⚠️  密钥生成失败，请手动设置 APP_KEY"
        fi
    else
        echo "✅ 密钥已存在，跳过生成"
    fi

    # 运行数据库迁移
    echo "📦 运行数据库迁移..."
    php artisan migrate --force
    echo "✅ 迁移完成"

    # 创建默认管理员账号（首次部署）
    echo "👤 创建默认管理员..."
    php artisan db:seed --class=AdminUserSeeder --force 2>/dev/null || true
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

# 从 APP_URL 解析协议和地址用于启动提示
APP_PROTO=$(echo "${APP_URL:-http://localhost}" | sed -n 's|^\([^:]*\)://.*|\1|p')
APP_HOST_PORT=$(echo "${APP_URL:-http://localhost}" | sed -n 's|^[^:]*://\(.*\)|\1|p')
echo "  访问地址: ${APP_PROTO:-http}://${APP_HOST_PORT:-localhost}"

# 启动 supervisor（管理 nginx + php-fpm + queue worker）
exec "$@"
