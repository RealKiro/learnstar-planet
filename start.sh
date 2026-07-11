#!/bin/bash
# ============================================================
# 学趣星球 - 一键启动脚本
# ============================================================
set -e

cd "$(dirname "$0")"

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

set -a; source .env; set +a

echo "========================================"
echo "  学趣星球"
echo "  APP_URL:  ${APP_URL:-未设置}"
echo "  APP_PORT: ${APP_PORT:-8080}"
echo "  端口映射: ${APP_PORT:-8080}:8080"
echo "========================================"

docker-compose up -d "$@"
