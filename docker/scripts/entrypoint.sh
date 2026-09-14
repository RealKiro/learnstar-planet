#!/bin/sh
# ============================================================
# 学宠星球 - Docker 入口脚本
# 初始化 Laravel 应用、运行迁移、启动 PHP 内置服务器
# ============================================================

set -e

echo "🚀 学宠星球 - 启动初始化..."

# 确保存储目录存在
mkdir -p storage/app/public \
    storage/framework/cache/data \
    storage/framework/sessions \
    storage/framework/views \
    storage/logs \
    bootstrap/cache

# 内置 SQLite：确保数据库文件存在（DB_CONNECTION=sqlite 时使用）
if [ "${DB_CONNECTION:-sqlite}" = "sqlite" ]; then
    # 数据卷若历史上被挂成目录（旧镜像未预建文件），SQLite 会彻底不可用 → 改名保留后重建文件
    if [ -d storage/database.sqlite ]; then
        SQLITE_BAK="storage/database.sqlite.bak.$(date +%Y%m%d%H%M%S)"
        echo "⚠️  storage/database.sqlite 是目录（历史数据卷残留），已改名保留：${SQLITE_BAK}"
        mv storage/database.sqlite "${SQLITE_BAK}"
    fi
    touch storage/database.sqlite
    chmod 666 storage/database.sqlite
fi

# ============================================================
# 从 Docker 环境变量创建 .env 文件
# 🔴 值必须加引号：phpdotenv 遇到未加引号的值内空格会抛
#    "Failed to parse dotenv file"，Laravel 随即以
#    "The environment file is invalid!" 退出 → 连 db:monitor 都跑不起来，
#    启动日志表现为「数据库未就绪」死循环 + 容器反复重启（真正的根因被掩盖）。
#    典型触发值：BOT_NAME=API 机器人（含空格）。
# ============================================================
ENV_KEY_PATTERN='^(APP_|DB_|REDIS_|CACHE_|SESSION_|QUEUE_|BROADCAST_|MAIL_|FILESYSTEM_|AI_|WECHAT_|DINGTALK_|FEISHU_|QQ_|RENREN_|ADMIN_|BOT_|UPLOAD_|GITHUB_)'

write_env_from_environment() {
    : > .env
    env | grep -E "${ENV_KEY_PATTERN}" | while IFS= read -r _env_line; do
        _env_key=${_env_line%%=*}
        _env_val=${_env_line#*=}
        # 跳过非法变量名（多行值被 env 拆行后的残留等），否则会写出无法解析的 .env
        case "${_env_key}" in
            '' | *[!A-Za-z0-9_]*) continue ;;
        esac
        # 先转义 \ " $ `，再整体用双引号包裹：空格 / 中文 / 特殊字符均安全
        _env_esc=$(printf '%s' "${_env_val}" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g' -e 's/\$/\\$/g' -e 's/`/\\`/g')
        printf '%s="%s"\n' "${_env_key}" "${_env_esc}"
    done >> .env
}

if [ ! -f .env ]; then
    echo "📝 从环境变量创建 .env 文件..."
    write_env_from_environment
elif ! php artisan --version >/dev/null 2>&1; then
    # .env 存在但无法解析（如旧版本写入的未加引号含空格值）→ 保留 APP_KEY 重建
    echo "⚠️  检测到 .env 无法解析，正在重建（保留 APP_KEY）..."
    OLD_APP_KEY=$(grep -E '^APP_KEY=' .env | head -n 1 || true)
    write_env_from_environment
    if [ -n "${OLD_APP_KEY}" ]; then
        printf '%s\n' "${OLD_APP_KEY}" >> .env
        echo "   已保留原有 APP_KEY"
    fi
fi

# 环境文件自检：解析失败时打印真实错误，避免后续现象误导排查方向
if ! php artisan --version >/dev/null 2>&1; then
    echo "🔴 .env 仍然无法被 Laravel 解析，artisan 输出如下："
    php artisan --version 2>&1 | head -n 5 | sed 's/^/   /' || true
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
        # 打印真实错误：避免「未就绪」掩盖 .env 解析失败 / 驱动缺失 / 认证失败等根因
        echo "   最后一次连接错误（末 5 行）："
        php artisan db:monitor 2>&1 | tail -n 5 | sed 's/^/   /' || true
        break
    fi
    echo "  数据库未就绪，3秒后重试（${RETRY_COUNT}/${MAX_RETRIES}）..."
    sleep 3
done

if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo "✅ 数据库连接成功"

    # 生成应用密钥（持久化到数据卷：容器重建后 APP_KEY 保持不变，用户登录态不失效）
    if ! grep -q "^APP_KEY=base64:" .env 2>/dev/null; then
        echo "🔑 检查应用密钥..."
        KEYFILE="storage/app/uploads/.app_key"
        NEW_KEY=""
        if [ -f "$KEYFILE" ] && [ -s "$KEYFILE" ]; then
            # 从数据卷恢复上次生成的密钥
            NEW_KEY=$(cat "$KEYFILE")
            echo "  从数据卷恢复 APP_KEY"
        else
            NEW_KEY=$(php -r "echo 'base64:' . base64_encode(random_bytes(32));" 2>/dev/null)
            if [ -z "$NEW_KEY" ]; then
                NEW_KEY=$(php artisan key:generate --show 2>/dev/null)
            fi
            if [ -n "$NEW_KEY" ]; then
                echo "$NEW_KEY" > "$KEYFILE"
                chmod 600 "$KEYFILE"
                echo "  已生成并保存 APP_KEY"
            fi
        fi
        if [ -n "$NEW_KEY" ]; then
            sed -i "/^APP_KEY=/d" .env
            echo "APP_KEY=$NEW_KEY" >> .env
            export APP_KEY="$NEW_KEY"
            echo "✅ 密钥已就绪"
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

    # 创建 API 机器人账号（供外部系统 REST API 对接，.env 的 BOT_* 可配置/停用）
    echo "🤖 检查 API 机器人账号..."
    php artisan db:seed --class=BotTeacherSeeder --force 2>/dev/null || true
    echo "✅ 机器人账号已就绪"

    # 创建存储链接
    echo "🔗 创建存储链接..."
    php artisan storage:link --force 2>/dev/null || true

    fi

# 优化缓存
echo "🧹 优化缓存..."
php artisan config:cache 2>/dev/null || true
php artisan route:cache 2>/dev/null || true
# 视图缓存仅在存在 Blade 视图时执行：本项目前端是 Vue SPA，resources/views 不存在，
# 无条件执行会报 "The /app/resources/views directory does not exist"，
# 且该报错走 **stdout**（2>/dev/null 拦不住），会污染启动日志
if [ -d resources/views ]; then
    php artisan view:cache 2>/dev/null || true
fi
echo "✅ 缓存重建完成"

echo "🎉 学宠星球初始化完成！"

# ============================================================
# 安全自检：默认凭据告警（不阻断启动，仅打印醒目提示）
# ============================================================
SECURITY_WARN=0
if [ "${ADMIN_PASSWORD:-admin123456}" = "admin123456" ]; then
    SECURITY_WARN=1
fi
if [ "${BOT_ENABLED:-true}" = "true" ] && [ "${BOT_PASSWORD:-learnstar-bot-2026}" = "learnstar-bot-2026" ]; then
    SECURITY_WARN=1
fi

if [ "$SECURITY_WARN" = "1" ]; then
    echo ""
    echo "🔴 ==================== 安全告警 ===================="
    echo "🔴 检测到默认凭据，当前部署存在风险："
    if [ "${ADMIN_PASSWORD:-admin123456}" = "admin123456" ]; then
        echo "🔴   · 管理员密码仍为默认值 admin123456"
    fi
    if [ "${BOT_ENABLED:-true}" = "true" ] && [ "${BOT_PASSWORD:-learnstar-bot-2026}" = "learnstar-bot-2026" ]; then
        echo "🔴   · API 机器人密码仍为默认值，且该账号可操作全部班级"
    fi
    echo "🔴 修复：编辑 .env / docker-compose.yml 后重启容器"
    echo "🔴   ADMIN_PASSWORD=<强密码>"
    echo "🔴   BOT_PASSWORD=<强密码>   或   BOT_ENABLED=false"
    echo "🔴 详见 README「🔐 上线前安全检查」"
    echo "🔴 ================================================="
    echo ""
fi

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
