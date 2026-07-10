#!/bin/bash
# ============================================================
# 学趣星球 - 一键启动脚本
# 自动从 APP_URL 解析端口，无需手动配置 APP_PORT
# ============================================================
set -e

cd "$(dirname "$0")"

# 加载 .env
if [ -f .env ]; then
    set -a; source .env; set +a
fi

# 从 APP_URL 自动解析端口
if [ -z "$APP_PORT" ] && [ -n "$APP_URL" ]; then
    PORT=$(echo "$APP_URL" | sed -n 's|.*://[^/:]*:\([0-9]\{1,5\}\).*|\1|p')
    if [ -z "$PORT" ]; then
        case "$APP_URL" in
            https://*) PORT=443 ;;
            *)          PORT=80 ;;
        esac
    fi
    export APP_PORT=$PORT
fi
export APP_PORT=${APP_PORT:-8080}

echo "========================================"
echo "  学趣星球"
echo "  APP_URL:  ${APP_URL:-未设置}"
echo "  端口映射: ${APP_PORT}:8080"
echo "========================================"

docker-compose up -d "$@"
