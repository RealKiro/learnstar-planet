#!/bin/sh
# 学宠星球 Go 后端入口：安全自检 → 启动。
#
# 运行时目录由 Go 程序自建（database.Connect 会 MkdirAll），这里提前建一遍
# 是为了让「先挂卷再启动」的路径一定存在（部分 Docker 版本挂卷时要求目录已建）。

set -e

mkdir -p /app/data /app/storage/app/uploads /app/logs

# ============================================================
# 启动安全自检：检测默认凭据，命中则在日志首屏打印告警横幅。
# 刻意不阻断启动——为保住「三步部署、零配置」体验，只把风险暴露出来
# （同旧 Laravel 版 entrypoint 的口径，README「上线前安全检查」章节对应）。
# ============================================================
warn=""
if [ "${ADMIN_PASSWORD:-admin123456}" = "admin123456" ]; then
    warn="${warn}
  - 管理员密码仍为默认值 admin123456：任何人都可登录后台管理全校数据。
    修改方法：.env 中设置 ADMIN_PASSWORD=你的强密码，然后 docker-compose up -d 重启"
fi
if [ "${BOT_ENABLED:-false}" = "true" ] && [ "${BOT_PASSWORD:-learnstar-bot-2026}" = "learnstar-bot-2026" ]; then
    warn="${warn}
  - API 机器人已启用且密码仍为默认值：该账号拥有本校全部班级的加减分权限。
    修改方法：.env 中设置 BOT_PASSWORD=你的强密码（或 BOT_ENABLED=false 停用），重启生效"
fi
if [ -n "$warn" ]; then
    echo ""
    echo "============================================================"
    echo "🔴 安全告警：检测到默认凭据，接入校园网前请务必修改！"
    echo "$warn"
    echo "============================================================"
    echo ""
fi

# JWT_SECRET 未显式配置时，Go 程序会自动生成随机密钥并持久化到
# /app/data/.jwt_secret（该目录挂了数据卷，容器重建不再全员登出）。

exec ./learnstar-go
