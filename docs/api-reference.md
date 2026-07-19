# 学趣星球 (LearnStar Planet) API 接口文档

> 版本：v1 | 协议：RESTful | 格式：JSON
> 适用于第三方应用接入、小程序、移动端、AI 机器人等

---

## 目录

1. [快速入门](#1-快速入门)
2. [认证方式](#2-认证方式)
3. [API 路径结构](#3-api-路径结构)
4. [通用规范](#4-通用规范)
5. [接口清单](#5-接口清单)
6. [错误码](#6-错误码)
7. [SDK 与示例](#7-sdk-与示例)
8. [限速说明](#8-限速说明)
9. [更新日志](#9-更新日志)

---

## 1. 快速入门

### 1.1 基础 URL

```
生产环境：https://your-domain.com/api/v1
开发环境：http://localhost:8000/api/v1
```

### 1.2 一分钟接入

```bash
# 1. 教师登录获取 Token
curl -X POST https://your-domain.com/api/v1/auth/teacher/login \
  -H "Content-Type: application/json" \
  -d '{"username":"teacher001","password":"123456"}'

# 2. 使用 Token 调用 API
curl https://your-domain.com/api/v1/teacher/dashboard \
  -H "Authorization: Bearer <your-token>"

# 3. 学生使用班级码登录
curl -X POST https://your-domain.com/api/v1/auth/class/login \
  -H "Content-Type: application/json" \
  -d '{"class_code":"LS301"}'
```

---

## 2. 认证方式

### 2.1 Bearer Token（教师/管理员/家长）

```
Authorization: Bearer <token>
```

Token 通过登录接口获取，默认不过期（可主动调用 `logout` 销毁）。

### 2.2 班级码 Token（学生端/大屏）

```
Authorization: Bearer class_LS301_xxxxxxxx...
```

班级码 Token 24 小时有效，通过 `POST /auth/class/login` 获取。

### 2.3 角色权限

| 角色 | 标识 | 可访问前缀 |
|------|------|-----------|
| 教师 | `teacher` | `/teacher/*` |
| 管理员 | `school_admin` | `/admin/*` |
| 家长 | `parent` | `/parent/*` |
| 学生/大屏 | 无角色 | `/display/*`（班级码 Token） |
| 公开 | 无需认证 | `/common/*` |

---

## 3. API 路径结构

```
/api/v1/{module}/{resource}[/{action}|/{id}][/{sub-resource}]
```

### 3.1 模块划分

| 模块 | 前缀 | 认证 | 说明 |
|------|------|------|------|
| 认证 | `/auth` | 部分需要 | 登录、登出、第三方绑定 |
| 教师 | `/teacher` | Bearer Token | 班级管理、积分、宠物、课堂工具 |
| 管理员 | `/admin` | Bearer Token | 学校级管理、批量操作 |
| 家长 | `/parent` | Bearer Token | 查看孩子数据 |
| 大屏 | `/display` | 班级码 Token | 教室大屏展示 |
| 公共 | `/common` | 无 | 宠物类型、进化阶段等 |

### 3.2 命名规范

```
✅ 正确示例：
  GET    /teacher/students          → 学生列表
  GET    /teacher/students/{id}     → 单个学生
  POST   /teacher/scores/give       → 加分操作
  PUT    /teacher/scores/rules/{id}  → 更新规则
  DELETE /teacher/scores/rules/{id}  → 删除规则

❌ 避免：
  /teacher/getStudentList           → 不用动词
  /teacher/student_list             → 不用下划线
  /teacher/students/{id}/scores     → 嵌套过深时可拆分
```

---

## 4. 通用规范

### 4.1 请求头

```json
{
  "Content-Type": "application/json",
  "Accept": "application/json",
  "Authorization": "Bearer <token>"
}
```

### 4.2 响应格式

**成功响应：**
```json
{
  "data": { ... },
  "message": "操作成功"
}
```

**列表响应（分页）：**
```json
{
  "data": [ ... ],
  "meta": {
    "current_page": 1,
    "last_page": 5,
    "per_page": 20,
    "total": 100
  }
}
```

**错误响应：**
```json
{
  "message": "错误描述",
  "errors": {
    "field_name": ["验证错误信息"]
  }
}
```

### 4.3 HTTP 状态码

| 状态码 | 含义 | 说明 |
|--------|------|------|
| 200 | OK | 请求成功 |
| 201 | Created | 创建成功 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未登录或 Token 失效 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 422 | Unprocessable Entity | 验证失败 |
| 429 | Too Many Requests | 请求过于频繁 |
| 500 | Internal Server Error | 服务器错误 |

### 4.4 分页

使用 Query 参数：
```
GET /teacher/students?per_page=20&page=2
```

### 4.5 筛选

使用 Query 参数：
```
GET /teacher/students?search=张&status=active
GET /teacher/scores/history?student_id=1&page=1
```

---

## 5. 接口清单

### 5.1 认证模块

#### `POST /auth/teacher/login` 教师登录

```
速率限制: 6 次/分钟
```

**请求体：**
```json
{
  "username": "teacher001",
  "password": "123456"
}
```

**响应：**
```json
{
  "data": {
    "token": "1|xxxxxxxx...",
    "user": {
      "id": 1,
      "username": "teacher001",
      "name": "张老师",
      "role": "teacher",
      "school_id": 1
    }
  }
}
```

#### `POST /auth/admin/login` 管理员登录

同上，角色为 `school_admin`。

#### `POST /auth/parent/login` 家长登录

同上，角色为 `parent`。

#### `POST /auth/class/login` 班级码登录

```
速率限制: 10 次/分钟（学生端）
```

**请求体：**
```json
{
  "class_code": "LS301"
}
```

**响应：**
```json
{
  "data": {
    "token": "class_LS301_xxxxxxxx...",
    "class_id": 1,
    "class_name": "三年级（1）班",
    "grade": "三年级",
    "student_count": 42
  }
}
```

#### `POST /auth/logout` 登出

需要认证。销毁当前 Token。

#### `POST /auth/change-password` 修改密码

**请求体：**
```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

#### 第三方登录

| 端点 | 说明 |
|------|------|
| `POST /auth/teacher/login/wechat` | 微信扫码登录 |
| `POST /auth/teacher/login/wechat-work` | 企业微信扫码登录 |
| `POST /auth/teacher/login/qq` | QQ 扫码登录 |
| `POST /auth/teacher/login/renren` | 人人通登录 |

#### 第三方绑定

| 端点 | 方法 | 说明 |
|------|------|------|
| `/auth/bind/{platform}` | POST | 绑定第三方账号 |
| `/auth/unbind/{platform}` | DELETE | 解绑 |
| `/auth/bindings` | GET | 获取绑定列表 |

`platform` 取值范围：`wechat` / `wechat_work` / `qq` / `renren`

---

### 5.2 教师端

#### 5.2.1 班级与仪表盘

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/dashboard` | GET | 班级总览（统计 + 班级之星 + TOP5 + 最新动态） |
| `/teacher/my-classes` | GET | 教师所管理的班级列表 |
| `/teacher/switch-class` | POST | 切换当前操作班级 |
| `/teacher/mode` | GET | 获取当前模式（大屏/管理） |
| `/teacher/mode` | POST | 切换模式 |

#### 5.2.2 学生管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/students` | GET | 学生列表（支持 `?search=关键字` 搜索） |
| `/teacher/students/import` | POST | 批量导入学生 |
| `/teacher/students/{id}` | PUT | 更新学生信息 |

#### 5.2.3 积分管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/scores/summary` | GET | 积分汇总（累计/今日/本周） |
| `/teacher/scores/give` | POST | 单人加减分 |
| `/teacher/scores/batch-give` | POST | 批量加减分 |
| `/teacher/scores/give-by-rule/{ruleId}` | POST | 按规则加减分 |
| `/teacher/scores/history/{studentId}` | GET | 学生积分记录 |
| `/teacher/scores/rules` | GET | 积分规则列表 |
| `/teacher/scores/rules` | POST | 创建积分规则 |
| `/teacher/scores/rules/{id}` | PUT | 更新规则 |
| `/teacher/scores/rules/{id}` | DELETE | 删除规则 |

**`POST /teacher/scores/give` 请求体：**
```json
{
  "student_id": 1,
  "points": 5,
  "reason": "举手发言"
}
```

**`POST /teacher/scores/batch-give` 请求体：**
```json
{
  "student_ids": [1, 2, 3, 4, 5],
  "points": 3,
  "reason": "作业优秀"
}
```

#### 5.2.4 宠物系统

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/pets/types` | GET | 宠物类型列表（按系列分组） |
| `/teacher/pets/overview` | GET | 全班宠物概览 |
| `/teacher/pets/{studentId}` | GET | 单个学生宠物详情 |
| `/teacher/pets/{studentId}/switch` | POST | 切换学生宠物类型 |
| `/teacher/pets/{studentId}/feed` | POST | 投喂宠物 |
| `/teacher/pets/{studentId}/rename` | POST | 重命名宠物 |

#### 5.2.5 排行榜 & 年级战场

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/leaderboard/total` | GET | 总积分榜 |
| `/teacher/leaderboard/weekly` | GET | 进步最快榜（本周） |
| `/teacher/leaderboard/pet-level` | GET | 宠物等级榜 |
| `/teacher/pk/leaderboard` | GET | 同年级 PK 排行榜 |
| `/teacher/pk/my-stats` | GET | 本班 PK 战力分析 |
| `/teacher/pk/challenge` | POST | 向其他班级发起挑战 |

#### 5.2.6 班级配置

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/class/switch-series` | POST | 切换班级宠物系列 |

**`POST /teacher/class/switch-series` 请求体：**
```json
{
  "series_id": "myth"
}
```

`series_id` 取值范围：`myth` / `pokemon` / `national` / `mecha` / `magic` / `prehistoric` / `constellation` / `folklore`

#### 5.2.7 积分商城

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/shop/items` | GET | 商品列表 |
| `/teacher/shop/items` | POST | 创建商品 |
| `/teacher/shop/items/{id}` | PUT | 更新商品 |
| `/teacher/shop/items/{id}` | DELETE | 删除商品 |
| `/teacher/shop/redemptions` | GET | 兑换记录列表 |
| `/teacher/shop/redemptions` | POST | 创建兑换请求 |
| `/teacher/shop/redemptions/{id}/approve` | PUT | 批准兑换 |
| `/teacher/shop/redemptions/{id}/reject` | PUT | 拒绝兑换 |
| `/teacher/shop/redemptions/{id}/deliver` | PUT | 标记已发放 |

#### 5.2.8 班级工具

**通知：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/notices` | GET | 通知列表 |
| `/teacher/notices` | POST | 创建通知 |
| `/teacher/notices/{id}` | PUT | 更新通知 |
| `/teacher/notices/{id}/publish` | PUT | 发布通知 |
| `/teacher/notices/{id}` | DELETE | 删除通知 |

**广播（教室大屏）：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/broadcasts` | GET | 广播列表 |
| `/teacher/broadcasts` | POST | 发送广播 |
| `/teacher/broadcasts/{id}` | GET | 获取单个广播 |

**教室大屏数据：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/classroom/display` | GET | 班级大屏全量数据 |
| `/teacher/classroom/messages` | GET | 轮询新消息 |
| `/teacher/classroom/messages` | POST | 发送消息到大屏 |

**考勤：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/attendance/today` | GET | 今日考勤记录 |
| `/teacher/attendance/start` | POST | 开始考勤 |
| `/teacher/attendance/{studentId}` | PUT | 设置学生考勤状态 |
| `/teacher/attendance/{studentId}/mark-leave` | POST | 标记请假 |
| `/teacher/attendance/{studentId}/mark-absent` | POST | 标记缺勤 |
| `/teacher/attendance/summary` | GET | 考勤汇总 |

**作业：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/homework` | GET | 作业列表 |
| `/teacher/homework` | POST | 创建作业 |
| `/teacher/homework/{id}` | GET | 作业详情 |
| `/teacher/homework/{id}/close` | PUT | 关闭作业 |
| `/teacher/homework/{id}/submissions` | GET | 提交列表 |
| `/teacher/homework/{id}/qr-code` | GET | 获取二维码 |

**在线答题：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/quizzes` | GET | 测验列表 |
| `/teacher/quizzes` | POST | 创建测验 |
| `/teacher/quizzes/{id}/start` | POST | 开始测验 |
| `/teacher/quizzes/{id}/stop` | POST | 结束测验 |
| `/teacher/quizzes/{id}/stats` | GET | 统计数据 |

**题库：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/question-banks` | GET | 题库列表 |
| `/teacher/question-banks` | POST | 创建题库 |
| `/teacher/question-banks/{id}/questions` | POST | 添加题目 |
| `/teacher/question-banks/{id}/questions` | GET | 获取题目列表 |

**成绩管理：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/grades` | GET | 成绩列表 |
| `/teacher/grades` | POST | 录入成绩 |
| `/teacher/grades/stats` | GET | 成绩统计 |
| `/teacher/grades/distribution` | GET | 成绩分布 |

**AI 助教：**
| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/ai/chat` | POST | AI 对话 |
| `/teacher/ai/commands` | GET | 获取命令列表 |
| `/teacher/ai/usage` | GET | 用量统计 |

#### 5.2.9 多币种

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/currency/wallets` | GET | 钱包列表 |
| `/teacher/currency/exchange` | POST | 币种兑换 |
| `/teacher/currency/cross-exchange` | POST | 交叉汇率兑换 |

#### 5.2.10 报表

| 端点 | 方法 | 说明 |
|------|------|------|
| `/teacher/reports/score-trend` | GET | 积分趋势 |
| `/teacher/reports/pet-distribution` | GET | 宠物等级分布 |
| `/teacher/reports/student-progress` | GET | 学生进步情况 |
| `/teacher/reports/export/{type}` | GET | 导出报表 |

---

### 5.3 管理员端

#### 5.3.1 学校信息

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/school` | GET | 获取学校信息 |

#### 5.3.2 教师管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/teachers` | GET | 教师列表 |
| `/admin/teachers` | POST | 创建教师 |
| `/admin/teachers/batch-create` | POST | 批量创建 |
| `/admin/teachers/import` | POST | 导入教师 |
| `/admin/teachers/template-csv` | GET | 下载导入模板 |
| `/admin/teachers/{id}` | PUT | 更新教师信息 |
| `/admin/teachers/{id}/reset-password` | POST | 重置密码 |
| `/admin/teachers/{id}/classes` | PUT | 分配班级 |
| `/admin/teachers/{id}` | DELETE | 禁用教师 |

#### 5.3.3 班级管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/classes` | GET | 班级列表 |
| `/admin/classes` | POST | 创建班级 |
| `/admin/classes/batch-create` | POST | 批量创建 |
| `/admin/classes/{id}` | GET | 班级详情 |
| `/admin/classes/{id}` | PUT | 更新班级 |
| `/admin/classes/{id}` | DELETE | 删除班级 |
| `/admin/classes/{id}/assign-teacher` | POST | 指定班主任 |
| `/admin/classes/{id}/remove-teacher` | DELETE | 移除班主任 |
| `/admin/classes/{classId}/display-code` | GET | 获取班级码 |
| `/admin/classes/{classId}/display-code/refresh` | POST | 刷新班级码 |

#### 5.3.4 学生管理（学校级）

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/students` | GET | 学生列表 |
| `/admin/students` | POST | 创建学生 |
| `/admin/students/import` | POST | 导入学生 |
| `/admin/students/{id}` | PUT | 更新学生 |
| `/admin/students/{id}` | DELETE | 删除学生 |
| `/admin/students/batch-delete` | POST | 批量删除 |
| `/admin/students/batch-move` | POST | 批量转班 |

#### 5.3.5 学年升级

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/grade-upgrade/preview` | GET | 升级预览 |
| `/admin/grade-upgrade/execute` | POST | 执行升级 |

#### 5.3.6 报表

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/reports/overview` | GET | 学校总览 |
| `/admin/reports/by-grade` | GET | 按年级统计 |
| `/admin/reports/by-class` | GET | 按班级统计 |

#### 5.3.7 汇率管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/exchange-rates` | GET | 汇率列表 |
| `/admin/exchange-rates` | POST | 创建汇率 |
| `/admin/exchange-rates/{id}` | PUT | 更新汇率 |

---

### 5.4 家长端

| 端点 | 方法 | 说明 |
|------|------|------|
| `/parent/home` | GET | 首页概览（孩子列表 + 积分 + 宠物） |
| `/parent/scores/detail` | GET | 积分明细 |
| `/parent/scores/history` | GET | 积分历史 |
| `/parent/growth/log` | GET | 成长日志 |
| `/parent/growth/timeline` | GET | 成长时间线 |
| `/parent/pet` | GET | 宠物状态 |
| `/parent/pet/feed` | POST | 投喂宠物 |
| `/parent/ranking` | GET | 班级排名 |
| `/parent/notices` | GET | 通知列表 |
| `/parent/notices/{id}` | GET | 阅读通知 |

---

### 5.5 班级大屏

| 端点 | 方法 | 说明 |
|------|------|------|
| `/display/login` | POST | 大屏登录（班级码） |
| `/display/initial-data` | GET | 获取初始全量数据 |
| `/display/sse` | GET | SSE 实时推送 |
| `/display/poll` | GET | 轮询更新 |
| `/display/quick-score` | POST | 快速加减分 |
| `/display/leaderboard` | GET | 排行榜 |
| `/display/shop-items` | GET | 商品列表 |
| `/display/redeem` | POST | 兑换商品 |
| `/display/transfer` | POST | 转账 |

---

### 5.6 公共接口

| 端点 | 方法 | 说明 |
|------|------|------|
| `/common/pet-types` | GET | 获取全部宠物类型和系列 |
| `/common/evolution-stages` | GET | 获取进化阶段列表 |
| `/common/score-categories` | GET | 获取积分分类 |

---

## 6. 错误码

| code | message | HTTP 状态码 | 说明 |
|------|---------|-----------|------|
| `AUTH_FAILED` | 账号或密码错误 | 401 | 登录验证失败 |
| `TOKEN_EXPIRED` | 登录已过期 | 401 | Token 无效或过期 |
| `FORBIDDEN` | 无权限 | 403 | 角色权限不足 |
| `NOT_FOUND` | 资源不存在 | 404 | 请求的资源不存在 |
| `VALIDATION_ERROR` | 验证错误 | 422 | 请求参数不符合要求 |
| `RATE_LIMIT` | 请求过于频繁 | 429 | 超过速率限制 |
| `BALANCE_INSUFFICIENT` | 积分不足 | 400 | 余额不足以完成操作 |
| `DUPLICATE_ENTRY` | 数据重复 | 400 | 唯一约束冲突 |

---

## 7. SDK 与示例

### 7.1 JavaScript / TypeScript

```typescript
// 使用项目的 services/api.ts
import { authApi, scoreApi, petApi } from '@/services/api'

// 登录
const { data } = await authApi.login({ username: 'teacher01', password: '123456' })
const token = data.token

// 加分
await scoreApi.giveScore({ student_id: 1, points: 5, reason: '举手发言' })

// 获取宠物
const pets = await petApi.getClassOverview()
```

### 7.2 Python

```python
import requests

BASE = "https://your-domain.com/api/v1"
token = None

# 登录
resp = requests.post(f"{BASE}/auth/teacher/login", json={
    "username": "teacher01",
    "password": "123456"
})
token = resp.json()["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

# 获取仪表盘
resp = requests.get(f"{BASE}/teacher/dashboard", headers=headers)
print(resp.json())

# 加分
resp = requests.post(f"{BASE}/teacher/scores/give", headers=headers, json={
    "student_id": 1,
    "points": 5,
    "reason": "举手发言"
})
```

### 7.3 cURL

```bash
#!/bin/bash
BASE="https://your-domain.com/api/v1"

# 登录
TOKEN=$(curl -s -X POST "$BASE/auth/teacher/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"teacher01","password":"123456"}' | jq -r '.data.token')

# 获取数据
curl -s "$BASE/teacher/dashboard" -H "Authorization: Bearer $TOKEN" | jq

# 加分
curl -s -X POST "$BASE/teacher/scores/give" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"student_id":1,"points":5,"reason":"举手发言"}'
```

### 7.4 微信小程序

```javascript
// app.js
App({
  globalData: { baseUrl: 'https://your-domain.com/api/v1' },

  login(username, password) {
    return wx.request({
      url: `${this.globalData.baseUrl}/auth/teacher/login`,
      method: 'POST',
      data: { username, password },
    })
  }
})

// 使用
const app = getApp()
const res = await app.login('teacher01', '123456')
const token = res.data.data.token
```

---

## 8. 限速说明

| 端点 | 限速 | 适用范围 |
|------|------|---------|
| `POST /auth/teacher/login` | 6 次/分钟 | 教师登录 |
| `POST /auth/admin/login` | 6 次/分钟 | 管理员登录 |
| `POST /auth/parent/login` | 6 次/分钟 | 家长登录 |
| `POST /auth/class/login` | 10 次/分钟 | 班级码登录 |
| `POST /display/login` | 10 次/分钟 | 大屏登录 |
| 其余 API | 未限制 | 建议客户端自行控制 |

---

## 9. 更新日志

| 日期 | 版本 | 变更内容 |
|------|------|---------|
| 2026-07-19 | v1.2 | 新增 PK 战场、班级码登录、系列切换接口；重构路由结构 |
| 2026-07-13 | v1.0 | 初始版本 |

---

> **📎 相关文档**
> - [API 端点映射](./api-endpoints.md) — 前后端联调对照表
> - 后端源码：`backend/routes/api.php` — 路由注册
> - 前端服务层：`frontend-vue/src/services/api.ts` — 类型安全 API 调用
