# LearnStar MCP Server

A standard MCP (Model Context Protocol) server that wraps the LearnStar Planet API, enabling AI-powered classroom management through QQ/WeChat bots via AstrBot or any MCP-compatible framework.

## Features

- **add_score**: 自然语言加分/减分 ("小明 +5 作业完成")
- **batch_add_score**: 批量加分 ("小明 小红 小刚 +5 作业完成")
- **query_score**: 查询学生积分和历史记录
- **get_leaderboard**: 排行榜（总积分/本周/宠物等级）
- **send_notice**: 发布班级通知
- **get_dashboard**: 查看班级概况
- **search_student**: 搜索学生

## Quick Start

```bash
# 1. Set environment variables
export LEARNSTAR_API_BASE=http://localhost:8080/api/v1
export LEARNSTAR_API_TOKEN=your_token_here

# 2. Run the server
python3 server.py
```

## Getting an API Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/teacher/login \
  -H "Content-Type: application/json" \
  -d '{"username":"teacher01","password":"your_password"}'
```

## Integration with AstrBot

1. Install the MCP client plugin in AstrBot
2. Add this server to the plugin config:

```yaml
mcp_servers:
  learnstar:
    command: python3
    args: ["/path/to/mcp-server/server.py"]
    env:
      LEARNSTAR_API_BASE: "http://your-server:8080/api/v1"
      LEARNSTAR_API_TOKEN: "your_token_here"
```

3. Restart AstrBot. The bot will now have these tools available.

## Usage via QQ/WeChat Bot

Once connected, users can interact naturally:

| User Message | Bot Action | Result |
|---|---|---|
| `小明 +5 作业完成` | add_score | ✅ 小明 +5分（作业完成）当前积分: 390分 |
| `小红 小刚 +3 主动发言` | batch_add_score | 批量加分完成 |
| `查分 小明` | query_score | 📊 小明 — 总积分: 390分 |
| `排行榜` | get_leaderboard | 🏆 总积分榜 TOP10 |
| `本周排行` | get_leaderboard(weekly) | 🏆 本周进步榜 |
| `通知 明天家长会提醒` | send_notice | 📢 通知已发布 |
| `班级概况` | get_dashboard | 📊 班级看板 |
| `搜学生 小` | search_student | 搜索「小」找到 6 个学生 |

## Demo Mode

If `LEARNSTAR_API_TOKEN` is not set, the server runs in demo mode with hardcoded sample data. This is useful for testing the MCP integration before connecting to a real backend.

## Architecture

```
QQ/WeChat Message
    ↓
AstrBot (message framework)
    ↓
MCP Client Plugin
    ↓ (JSON-RPC over stdio)
LearnStar MCP Server
    ↓ (REST API)
LearnStar Laravel Backend
```
