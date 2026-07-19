# 🌌 学趣星球 Wiki

欢迎来到学趣星球 Wiki！这里包含了项目的完整文档。

## 📖 快速导航

| 页面 | 说明 |
|------|------|
| [🏠 首页](Home) | 你在这里 |
| [🚀 部署指南](Deployment-Guide) | Docker 一键部署 & 手动部署 |
| [🎯 核心功能](Core-Features) | 班级总览、课堂评价、年级战场、宠物图鉴 |
| [🔧 教师端](Teacher-Guide) | 高级管理功能介绍 |
| [📚 宠物系统](Pet-System) | 宠物数据、进化路线、诗文 |
| [🔗 API 参考](API-Reference) | 第三方应用接入文档 |
| [⚙️ 配置说明](Configuration) | 环境变量与数据库配置 |
| [❓ FAQ](FAQ) | 常见问题 |

## 🎯 快速开始

### 学生/班级码登录
1. 向班主任获取班级码
2. 打开首页，输入班级码
3. 进入教室端：**班级总览 → 课堂评价 → 年级战场 → 宠物图鉴**

### 教师/管理员登录
1. 点击右上角「教师登录」或「管理员」
2. 输入账号密码（管理员默认：`admin / admin123456`）
3. 进入教师端：全部管理功能

### Docker 部署
```bash
git clone https://github.com/RealKiro/learnstar-planet.git
cd learnstar-planet
cp .env.example .env
docker-compose up -d
```

---

> 📝 更多内容请参考左侧导航栏
