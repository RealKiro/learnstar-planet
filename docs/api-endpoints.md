# 学趣星球 API 端点映射

> 基于 `功能分类与版面设计.txt` 整理的前后端联调对照表
> 前端服务层：`frontend-vue/src/services/api.ts`

---

## 四大核心模块 API

### 🏠 班级总览

| 前端调用 | 方法 | 后端端点 | 说明 |
|---------|------|---------|------|
| `overviewApi.getOverview()` | GET | `/api/v1/teacher/dashboard` | 班级总览数据（班级之星、概况统计、TOP5、动态） |
| `overviewApi.getRecentNews()` | GET | `/api/v1/teacher/dashboard/news` | 仅获取最新动态列表 |

**响应结构** (`getOverview`):
```json
{
  "data": {
    "class_name": "三年级一班",
    "grade": "三年级",
    "student_count": 42,
    "total_score": 3840,
    "avg_pet_level": 6.2,
    "peak_count": 5,
    "weekly_score": 1260,
    "star_student": {
      "name": "张小明",
      "pet_name": "九尾天狐",
      "pet_species": "nine_tail_fox",
      "pet_level": 12,
      "score": 520
    },
    "top5": [
      { "name": "张小明", "score": 520, "pet_name": "九尾天狐", "pet_species": "nine_tail_fox", "pet_level": 12 }
    ],
    "recent_news": [
      { "icon": "🎉", "text": "孙七的【机械龙】进化到了 Lv.8！" }
    ]
  }
}
```

---

### ✏️ 课堂评价

| 前端调用 | 方法 | 后端端点 | 说明 |
|---------|------|---------|------|
| `scoreApi.getStudents()` | GET | `/api/v1/teacher/students?per_page=100` | 学生列表（含积分） |
| `scoreApi.getRules()` | GET | `/api/v1/teacher/scores/rules` | 积分规则列表 |
| `scoreApi.getSummary()` | GET | `/api/v1/teacher/scores/summary` | 积分汇总统计 |
| `scoreApi.giveScore()` | POST | `/api/v1/teacher/scores/give` | 单人加减分 |
| `scoreApi.batchGiveScore()` | POST | `/api/v1/teacher/scores/batch-give` | 批量加减分 |
| `scoreApi.getHistory()` | GET | `/api/v1/teacher/scores/history` | 积分记录 |

**`giveScore` 请求体**:
```json
{
  "student_id": 1,
  "points": 5,
  "reason": "📖 举手发言"
}
```

**`batchGiveScore` 请求体**:
```json
{
  "student_ids": [1, 2, 3, 4, 5],
  "points": 3,
  "reason": "✅ 作业优秀"
}
```

---

### 🏆 年级战场 (PK)

| 前端调用 | 方法 | 后端端点 | 说明 |
|---------|------|---------|------|
| `pkApi.getLeaderboard()` | GET | `/api/v1/teacher/leaderboard/pk` | 同年级各班PK数据 |
| `pkApi.getMyStats()` | GET | `/api/v1/teacher/leaderboard/my-stats` | 本班战力详情 |
| `pkApi.challenge()` | POST | `/api/v1/teacher/pk/challenge` | 发起PK挑战 |

**`getLeaderboard` 响应**:
```json
{
  "data": [
    {
      "name": "三年级二班",
      "totalScore": 4200,
      "studentCount": 45,
      "avgLevel": 7.2,
      "peakCount": 8,
      "weekGrowth": 156,
      "isOwn": false
    }
  ]
}
```

---

### 📚 宠物图鉴

| 前端调用 | 方法 | 后端端点 | 说明 |
|---------|------|---------|------|
| `petApi.getClassOverview()` | GET | `/api/v1/teacher/pets/class-overview` | 全班宠物概览 |
| `petApi.getPetDetail()` | GET | `/api/v1/teacher/pets/:id` | 单个宠物详情 |
| `petApi.feedPet()` | POST | `/api/v1/teacher/pets/feed` | 投喂宠物 |
| `petApi.renamePet()` | PUT | `/api/v1/teacher/pets/:id/rename` | 重命名宠物 |
| `petApi.getSeries()` | GET | `/api/v1/common/pet-types` | 获取全部系列 |
| `petApi.switchSeries()` | POST | `/api/v1/teacher/class/switch-series` | 切换班级系列 |

---

## 现有后端端点检查清单

> 标记 ✓ = 已实现，✗ = 需新增

| 端点 | 状态 | 备注 |
|------|------|------|
| `POST /api/v1/auth/teacher/login` | ✓ | 教师登录 |
| `POST /api/v1/auth/class/login` | ✗ | 班级码登录（演示功能） |
| `GET /api/v1/teacher/dashboard` | ✓ | 教师仪表盘 |
| `GET /api/v1/teacher/students` | ✓ | 学生列表 |
| `GET /api/v1/teacher/scores/rules` | ✓ | 积分规则 |
| `GET /api/v1/teacher/scores/summary` | ✓ | 积分汇总 |
| `POST /api/v1/teacher/scores/give` | ✓ | 加减分 |
| `POST /api/v1/teacher/scores/batch-give` | ✓ | 批量加减分 |
| `GET /api/v1/teacher/pets/class-overview` | ✓ | 宠物概览 |
| `POST /api/v1/teacher/pets/feed` | ✓ | 投喂宠物 |
| `GET /api/v1/teacher/leaderboard/total` | ✓ | 总积分榜 |
| `GET /api/v1/teacher/leaderboard/weekly` | ✓ | 周榜 |
| `GET /api/v1/teacher/leaderboard/pet` | ✓ | 宠物等级榜 |
| `GET /api/v1/teacher/leaderboard/pk` | ✗ | **年级PK榜（待实现）** |
| `GET /api/v1/teacher/leaderboard/my-stats` | ✗ | **本班PK战力（待实现）** |
| `POST /api/v1/teacher/pk/challenge` | ✗ | **发起挑战（待实现）** |
| `POST /api/v1/teacher/class/switch-series` | ✗ | **切换系列（待实现）** |
| `GET /api/v1/common/pet-types` | ✓ | 宠物类型 |
| `POST /api/v1/teacher/pets/feed` | ✓ | 投喂宠物 |

---

## 前端页面路由映射

| 设计文档模块 | 前端路由 | 组件文件 | 状态 |
|------------|---------|---------|------|
| 🏠 班级总览 | `/teacher/dashboard` | `pages/teacher/DashboardPage.vue` | ✅ 已重构 |
| ✏️ 课堂评价 | `/teacher/scores` | `pages/teacher/ScoresPage.vue` | ✅ 已重构 |
| 🏆 年级战场 | `/teacher/pk` | `pages/teacher/PKPage.vue` | ✅ 新建 |
| 📚 宠物图鉴 | `/teacher/pets` | `pages/teacher/PetsPage.vue` | ✅ 已重构 |
| 🌟 宠物详情(家长) | `/parent/pet` | `pages/parent/PetDetailPage.vue` | ✅ 新建 |
