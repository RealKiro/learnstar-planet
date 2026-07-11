# 学趣星球 (LearnStar Planet) — Web 前端

Vue 3 + Vite + TypeScript 构建的模块化 SPA，替代原有的单文件 `index.html`。

## 技术栈

- **Vue 3** — Composition API (`<script setup>`)
- **Vite 8** — 构建工具
- **TypeScript** — 类型安全
- **Vue Router 4** — SPA 路由（含角色守卫）
- **Pinia 3** — 状态管理
- **Axios** — HTTP 客户端（token 拦截 + 401 处理）
- **Tailwind CSS 4** — 原子化样式

## 本地开发

```bash
# 安装依赖
npm install

# 启动开发服务器（http://localhost:5173）
npm run dev

# TypeScript 类型检查
npm run typecheck

# 生产构建（输出到 dist/）
npm run build

# 构建并输出到后端 public/ 目录（用于 Docker 部署）
npm run build:deploy
```

## 项目结构

```
src/
├── main.ts               # 入口（Pinia + Router + Theme init）
├── App.vue               # 根组件
├── types/index.ts        # TypeScript 类型定义
├── router/index.ts       # 路由树 + 角色守卫
├── utils/
│   ├── api.ts            # Axios 封装
│   └── constants.ts      # 常量与工具函数
├── stores/
│   ├── auth.ts           # 认证状态
│   ├── theme.ts          # 暗色模式
│   ├── toast.ts          # Toast 消息
│   └── crud.ts           # 通用 CRUD Store
├── layouts/
│   ├── TeacherLayout.vue # 教师端布局
│   ├── AdminLayout.vue   # 管理端布局
│   └── ParentLayout.vue  # 家长端布局
├── pages/
│   ├── landing/          # 首页 Landing
│   ├── auth/             # 登录页
│   ├── teacher/          # 教师端页面（16个）
│   ├── admin/            # 管理端页面（8个）
│   └── parent/           # 家长端页面（5个）
└── components/
    └── common/           # 通用组件
```

## 部署

Vue 构建产物需部署到后端 `public/` 目录，由 Nginx 统一 serve（SPA 模式 `try_files $uri /index.html`）。两种方式：

1. **本地构建**：`npm run build:deploy` 直接输出到 `../backend/public/`
2. **Docker 部署**：`backend/Dockerfile` 已包含 Node 阶段自动构建前端并复制产物
