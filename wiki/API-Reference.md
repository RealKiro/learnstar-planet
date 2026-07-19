# 🔗 API 参考

> 完整文档：[docs/api-reference.md](https://github.com/RealKiro/learnstar-planet/blob/main/docs/api-reference.md)

## 基础 URL

```
生产环境：https://your-domain.com/api/v1
开发环境：http://localhost:8000/api/v1
```

## 认证方式

| 方式 | Token 格式 | 有效期 | 适用场景 |
|------|-----------|--------|---------|
| Bearer Token | Sanctum Token | 永久 | 教师/管理员/家长 |
| 班级码 Token | `class_*` | 24 小时 | 学生端/教室大屏 |

## 快速示例

### 班级码登录（学生端）
```bash
curl -X POST /api/v1/auth/class/login \
  -d '{"class_code":"LS301"}'
```

### 教师登录
```bash
curl -X POST /api/v1/auth/teacher/login \
  -d '{"username":"teacher01","password":"123456"}'
```

### 获取班级总览
```bash
curl /api/v1/display/dashboard?token=<class_token>
```

### 加减分
```bash
curl -X POST /api/v1/display/scores/give \
  -d '{"token":"<token>","student_id":1,"points":5,"reason":"举手发言"}'
```

### 切换宠物（首次免费）
```bash
curl -X POST /api/v1/display/pets/switch \
  -d '{"token":"<token>","student_id":1,"pet_species":"panda"}'
```
