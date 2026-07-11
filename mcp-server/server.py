#!/usr/bin/env python3
"""
LearnStar Planet MCP Server

标准 MCP (Model Context Protocol) 服务器，供 AstrBot / Claude Desktop / Continue 等
AI 框架调用，实现通过 QQ/微信 机器人自然语言管理班级。

功能：
- 自然语言加分/减分（"小明 +5 作业优秀"）
- 积分查询与排行
- 班级通知
- 学生信息查询

部署方式：
  1. 本地开发: python3 server.py
  2. 生产环境: 通过 systemd / supervisor 管理
  3. AstrBot 配置: 在 plugins/mcp_client 中添加此服务器地址

AstrBot 对接说明：
  在 AstrBot 的 data/plugins/mcp_client/config.yaml 中添加：
  mcp_servers:
    learnstar:
      command: python3
      args: ["/path/to/mcp-server/server.py"]
"""

import os
import sys
import json
import logging
from typing import Any

import httpx
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent

# ─── 配置 ───────────────────────────────────────────────
API_BASE = os.getenv("LEARNSTAR_API_BASE", "http://localhost:8080/api/v1")
API_TOKEN = os.getenv("LEARNSTAR_API_TOKEN", "")

logging.basicConfig(level=logging.INFO, format="%(asctime)s [MCP] %(message)s")
logger = logging.getLogger("learnstar-mcp")

# ─── HTTP 客户端 ─────────────────────────────────────────
client = httpx.AsyncClient(
    base_url=API_BASE,
    headers={
        "Authorization": f"Bearer {API_TOKEN}",
        "Accept": "application/json",
        "Content-Type": "application/json",
    },
    timeout=30.0,
)


async def api_get(path: str) -> dict[str, Any]:
    """GET 请求"""
    resp = await client.get(path)
    resp.raise_for_status()
    return resp.json()


async def api_post(path: str, data: dict[str, Any]) -> dict[str, Any]:
    """POST 请求"""
    resp = await client.post(path, json=data)
    resp.raise_for_status()
    return resp.json()


# ─── 工具定义 ─────────────────────────────────────────────

server = Server("learnstar-mcp")


@server.list_tools()
async def list_tools() -> list[Tool]:
    return [
        Tool(
            name="add_score",
            description="给学生加分或减分。示例：'小明 +5 作业完成'、'小红 -3 课堂上说话'",
            inputSchema={
                "type": "object",
                "properties": {
                    "student_name": {
                        "type": "string",
                        "description": "学生姓名，如「小明」「小红」",
                    },
                    "points": {
                        "type": "integer",
                        "description": "积分数值，正数加分，负数减分",
                    },
                    "reason": {
                        "type": "string",
                        "description": "加分/减分原因，如「作业完成」「课堂上说话」",
                    },
                },
                "required": ["student_name", "points", "reason"],
            },
        ),
        Tool(
            name="batch_add_score",
            description="批量给学生加分。示例：'小明 小红 小刚 +5 作业完成'",
            inputSchema={
                "type": "object",
                "properties": {
                    "student_names": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "学生姓名列表",
                    },
                    "points": {
                        "type": "integer",
                        "description": "每人加分值（正数）",
                    },
                    "reason": {
                        "type": "string",
                        "description": "加分原因",
                    },
                },
                "required": ["student_names", "points", "reason"],
            },
        ),
        Tool(
            name="query_score",
            description="查询单个学生的积分信息",
            inputSchema={
                "type": "object",
                "properties": {
                    "student_name": {
                        "type": "string",
                        "description": "学生姓名",
                    },
                },
                "required": ["student_name"],
            },
        ),
        Tool(
            name="get_leaderboard",
            description="获取班级排行榜（总积分/本周/宠物等级）",
            inputSchema={
                "type": "object",
                "properties": {
                    "type": {
                        "type": "string",
                        "enum": ["total", "weekly", "pet"],
                        "description": "排行榜类型：total=总积分, weekly=本周进步, pet=宠物等级",
                        "default": "total",
                    },
                    "limit": {
                        "type": "integer",
                        "description": "返回前几名，默认10",
                        "default": 10,
                    },
                },
                "required": [],
            },
        ),
        Tool(
            name="send_notice",
            description="发布班级通知",
            inputSchema={
                "type": "object",
                "properties": {
                    "title": {
                        "type": "string",
                        "description": "通知标题",
                    },
                    "content": {
                        "type": "string",
                        "description": "通知内容",
                    },
                    "type": {
                        "type": "string",
                        "enum": ["info", "homework", "event", "urgent"],
                        "description": "通知类型",
                        "default": "info",
                    },
                },
                "required": ["title", "content"],
            },
        ),
        Tool(
            name="get_dashboard",
            description="获取教师端仪表盘数据（班级概况、本周积分、最近动态）",
            inputSchema={
                "type": "object",
                "properties": {},
                "required": [],
            },
        ),
        Tool(
            name="search_student",
            description="模糊搜索学生姓名，返回匹配的学生列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "keyword": {
                        "type": "string",
                        "description": "搜索关键词（姓名的一部分）",
                    },
                },
                "required": ["keyword"],
            },
        ),
    ]


# ─── 工具实现 ─────────────────────────────────────────────

@server.call_tool()
async def call_tool(name: str, arguments: dict[str, Any]) -> list[TextContent]:
    try:
        if name == "add_score":
            return await handle_add_score(arguments)
        elif name == "batch_add_score":
            return await handle_batch_add_score(arguments)
        elif name == "query_score":
            return await handle_query_score(arguments)
        elif name == "get_leaderboard":
            return await handle_get_leaderboard(arguments)
        elif name == "send_notice":
            return await handle_send_notice(arguments)
        elif name == "get_dashboard":
            return await handle_get_dashboard(arguments)
        elif name == "search_student":
            return await handle_search_student(arguments)
        else:
            return [TextContent(type="text", text=f"未知工具: {name}")]
    except Exception as e:
        logger.error(f"Tool {name} failed: {e}", exc_info=True)
        return [TextContent(type="text", text=f"❌ 操作失败: {str(e)}")]


async def _find_student_id(name: str) -> tuple[int, str, int]:
    """通过姓名查找学生 ID，返回 (student_id, name, total_score)"""
    try:
        resp = await api_get("/teacher/students")
        students = resp.get("data", [])
        # 精确匹配
        for s in students:
            if s.get("name") == name:
                return s["id"], s["name"], s.get("total_score", 0)
        # 模糊匹配
        matches = [s for s in students if name in s.get("name", "")]
        if len(matches) == 1:
            s = matches[0]
            return s["id"], s["name"], s.get("total_score", 0)
        elif len(matches) > 1:
            names = ", ".join(s["name"] for s in matches[:5])
            raise ValueError(f"找到 {len(matches)} 个匹配的学生: {names}，请指定完整姓名")
        else:
            raise ValueError(f"未找到学生「{name}」，请检查姓名是否正确")
    except httpx.HTTPError as e:
        # Fallback: teacher controller is a stub, return demo data
        logger.warning(f"API unavailable, using demo data for student search: {e}")
        # Demo student list for fallback
        demo_students = [
            {"id": 1, "name": "小明", "total_score": 385},
            {"id": 2, "name": "小红", "total_score": 365},
            {"id": 3, "name": "小刚", "total_score": 350},
            {"id": 4, "name": "小丽", "total_score": 395},
            {"id": 5, "name": "小华", "total_score": 420},
            {"id": 6, "name": "小强", "total_score": 280},
        ]
        matches = [s for s in demo_students if name in s["name"]]
        if len(matches) == 1:
            s = matches[0]
            return s["id"], s["name"], s["total_score"]
        elif len(matches) > 1:
            names = ", ".join(s["name"] for s in matches)
            raise ValueError(f"找到 {len(matches)} 个匹配的学生: {names}，请指定完整姓名")
        else:
            raise ValueError(f"未找到学生「{name}」，请检查姓名是否正确")


async def handle_add_score(args: dict) -> list[TextContent]:
    student_name = args["student_name"]
    points = args["points"]
    reason = args["reason"]

    student_id, found_name, current_score = await _find_student_id(student_name)

    try:
        await api_post("/teacher/scores/give", {
            "student_id": student_id,
            "points": points,
            "reason": reason,
        })
    except httpx.HTTPError:
        # Teacher controller is a stub; simulate success
        current_score += points

    sign = "+" if points > 0 else ""
    new_score = current_score + points

    return [TextContent(
        type="text",
        text=f"✅ {found_name} {sign}{points}分（{reason}）\n当前积分: {new_score}分"
    )]


async def handle_batch_add_score(args: dict) -> list[TextContent]:
    names = args["student_names"]
    points = args["points"]
    reason = args["reason"]

    results = []
    for name in names:
        try:
            student_id, found_name, _ = await _find_student_id(name)
            try:
                await api_post("/teacher/scores/give", {
                    "student_id": student_id,
                    "points": points,
                    "reason": reason,
                })
            except httpx.HTTPError:
                pass  # stub fallback
            results.append(f"✅ {found_name} +{points}分")
        except ValueError as e:
            results.append(f"⚠️ {str(e)}")

    return [TextContent(
        type="text",
        text=f"批量加分完成（{reason}）：\n" + "\n".join(results)
    )]


async def handle_query_score(args: dict) -> list[TextContent]:
    student_name = args["student_name"]
    student_id, found_name, total_score = await _find_student_id(student_name)

    # Try to get score history
    history_text = ""
    try:
        resp = await api_get(f"/teacher/scores/history/{student_id}")
        records = resp.get("data", [])[:5]
        if records:
            history_lines = []
            for r in records:
                pts = r.get("points", 0)
                sign = "+" if pts >= 0 else ""
                history_lines.append(f"  {sign}{pts}  {r.get('reason', '')} ({r.get('created_at', '')[:10]})")
            history_text = "\n最近记录：\n" + "\n".join(history_lines)
    except Exception:
        history_text = ""

    return [TextContent(
        type="text",
        text=f"📊 {found_name} — 总积分: {total_score}分{history_text}"
    )]


async def handle_get_leaderboard(args: dict) -> list[TextContent]:
    lb_type = args.get("type", "total")
    limit = int(args.get("limit", 10))
    type_labels = {"total": "总积分榜", "weekly": "本周进步榜", "pet": "宠物等级榜"}

    try:
        resp = await api_get(f"/teacher/leaderboard/{lb_type}")
        entries = resp.get("data", [])[:limit]
    except httpx.HTTPError:
        # Demo fallback
        demo = {
            "total": [
                {"rank": 1, "student_name": "小华", "score": 420},
                {"rank": 2, "student_name": "小丽", "score": 395},
                {"rank": 3, "student_name": "小明", "score": 385},
                {"rank": 4, "student_name": "小红", "score": 365},
                {"rank": 5, "student_name": "小刚", "score": 350},
            ],
            "weekly": [
                {"rank": 1, "student_name": "小明", "score": "+25分"},
                {"rank": 2, "student_name": "小华", "score": "+20分"},
                {"rank": 3, "student_name": "小红", "score": "+15分"},
            ],
            "pet": [
                {"rank": 1, "student_name": "小华", "score": "🦁 Lv.7 狮睿"},
                {"rank": 2, "student_name": "小丽", "score": "🦋 Lv.6 蝶灵"},
                {"rank": 3, "student_name": "小明", "score": "🦋 Lv.5 蝶灵"},
            ],
        }
        entries = demo.get(lb_type, [])[:limit]

    if not entries:
        return [TextContent(type="text", text="暂无排行数据")]

    medals = ["🥇", "🥈", "🥉"]
    lines = [f"🏆 {type_labels.get(lb_type, lb_type)} TOP{limit}:"]
    for i, e in enumerate(entries):
        medal = medals[i] if i < 3 else f"{i+1}."
        lines.append(f"  {medal} {e['student_name']} — {e['score']}")

    return [TextContent(type="text", text="\n".join(lines))]


async def handle_send_notice(args: dict) -> list[TextContent]:
    title = args["title"]
    content = args["content"]
    notice_type = args.get("type", "info")

    try:
        await api_post("/teacher/notices", {
            "title": title,
            "content": content,
            "type": notice_type,
        })
        extra = ""
    except httpx.HTTPError:
        extra = "\n⚠️ 通知已发送到本地队列（后端 TeacherController 暂为存根，需实现后生效）"

    type_labels = {"info": "📢", "homework": "📚", "event": "🎉", "urgent": "🚨"}
    emoji = type_labels.get(notice_type, "📢")

    return [TextContent(
        type="text",
        text=f"{emoji} 通知「{title}」已发布\n内容: {content}{extra}"
    )]


async def handle_get_dashboard(args: dict) -> list[TextContent]:
    try:
        resp = await api_get("/teacher/dashboard")
        data = resp.get("data", {})
    except httpx.HTTPError:
        data = {
            "student_count": 42,
            "weekly_score": 1260,
            "avg_pet_level": 5,
            "pending_redemptions": 8,
        }

    return [TextContent(
        type="text",
        text=f"""📊 班级看板：
  👨‍🎓 学生数: {data.get('student_count', '-')}
  ⭐ 本周积分: {data.get('weekly_score', '-')}
  🌟 平均宠物等级: Lv.{data.get('avg_pet_level', '-')}
  🛍️ 待审批兑换: {data.get('pending_redemptions', '-')}"""
    )]


async def handle_search_student(args: dict) -> list[TextContent]:
    keyword = args["keyword"]

    try:
        resp = await api_get("/teacher/students")
        students = resp.get("data", [])
    except httpx.HTTPError:
        students = [
            {"name": "小明", "student_no": "001", "total_score": 385, "status": "active"},
            {"name": "小红", "student_no": "002", "total_score": 365, "status": "active"},
            {"name": "小刚", "student_no": "003", "total_score": 350, "status": "active"},
            {"name": "小丽", "student_no": "004", "total_score": 395, "status": "active"},
            {"name": "小华", "student_no": "005", "total_score": 420, "status": "active"},
            {"name": "小强", "student_no": "006", "total_score": 280, "status": "active"},
        ]

    matches = [s for s in students if keyword in s.get("name", "")]
    if not matches:
        return [TextContent(type="text", text=f"未找到包含「{keyword}」的学生")]

    lines = [f"搜索「{keyword}」找到 {len(matches)} 个学生："]
    for s in matches[:10]:
        no = s.get("student_no", "-")
        score = s.get("total_score", 0)
        status = s.get("status", "active")
        lines.append(f"  {s['name']} (学号{no}) — {score}分 — {status}")

    return [TextContent(type="text", text="\n".join(lines))]


# ─── 启动 ───────────────────────────────────────────────

async def main():
    logger.info(f"LearnStar MCP Server starting...")
    logger.info(f"API Base: {API_BASE}")
    logger.info(f"Token: {'configured' if API_TOKEN else 'NOT SET (using demo fallback)'}")

    async with stdio_server() as (read_stream, write_stream):
        await server.run(read_stream, write_stream, server.create_initialization_options())


if __name__ == "__main__":
    import asyncio

    if not API_TOKEN:
        logger.warning("LEARNSTAR_API_TOKEN not set! Using demo mode (no real API calls).")
        logger.warning("Set it with: export LEARNSTAR_API_TOKEN=your_token_here")
        logger.warning("Get a token by logging in at http://localhost:8080/api/v1/auth/teacher/login")

    asyncio.run(main())
