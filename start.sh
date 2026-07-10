#!/bin/bash
# ============================================================
# 学趣星球 - 一键启动脚本
# 从 APP_URL 自动解析端口，写入 .env，docker-compose 自动读取
# ============================================================
set -e

cd "$(dirname "$0")"

# 确保 .env 存在
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "已从 .env.example 创建 .env"
        echo "请修改 .env 中的配置后重新运行 bash start.sh"
        exit 1
    else
        echo "错误：.env 和 .env.example 都不存在"
        exit 1
    fi
fi

# 加载 .env
set -a; source .env; set +a

# 从 APP_URL 自动解析端口
if [ -n "$APP_URL" ]; then
    PORT=$(echo "$APP_URL" | sed -n 's|.*://[^/:]*:\([0-9]\{1,5\}\).*|\1|p')
    if [ -z "$PORT" ]; then
        case "$APP_URL" in
            https://*) PORT=443 ;;
            *)          PORT=80 ;;
        esac
    fi
else
    PORT=${APP_PORT:-8080}
fi

# 将 APP_PORT 写回 .env（这样面板/docker-compose 都能读到）
if grep -q '^APP_PORT=' .env; then
    CURRENT=$(grep '^APP_PORT=' .env | cut -d= -f2)
    if [ "$CURRENT" != "$PORT" ]; then
        sed -i "s/^APP_PORT=.*/APP_PORT=$PORT/" .env
        echo "已更新 APP_PORT=$PORT（从 APP_URL 自动解析）"
    fi
else
    echo "APP_PORT=$PORT" >> .env
    echo "已添加 APP_PORT=$PORT（从 APP_URL 自动解析）"
fi

export APP_PORT=$PORT

echo "========================================"
echo "  学趣星球"
echo "  APP_URL:  ${APP_URL:-未设置}"
echo "  端口映射: ${APP_PORT}:8080"
echo "========================================"

docker-compose up -d "$@"
