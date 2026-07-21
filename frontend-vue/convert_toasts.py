import os

base = r"C:\Users\ZUOQIHE\Documents\GitHub\learnstar-planet\learnstar-planet\frontend-vue"

# File 1: LoginPage.vue
path = os.path.join(base, "src", "pages", "auth", "LoginPage.vue")
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace(
    "else toast.show((e?.response?.data?.message || '账号或密码错误') + `，还剩 ${remaining} 次`, 'error')",
    "else toast.show((e?.response?.data?.message || '账号或密码错误') + `，还剩 ${remaining} 次`, 'error', { position: 'center', duration: 3000 })"
)
content = content.replace(
    "toast.show(`正在打开${label}扫码...`, 'success')",
    "toast.show(`正在打开${label}扫码...`, 'success', { position: 'center', duration: 1500 })"
)
content = content.replace(
    "toast.show('人人通登录需要管理员在后台配置', 'info')",
    "toast.show('人人通登录需要管理员在后台配置', 'info', { position: 'center', duration: 2000 })"
)
content = content.replace(
    "toast.show('登录成功', 'success')",
    "toast.show('登录成功', 'success', { position: 'center', duration: 1500 })"
)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)
print("LoginPage.vue done")

# File 2: LoginPanel.vue
path = os.path.join(base, "src", "pages", "landing", "LoginPanel.vue")
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace(
    "toast.show('正在打开' + (p ? p.label : key) + '扫码...', 'success')",
    "toast.show('正在打开' + (p ? p.label : key) + '扫码...', 'success', { position: 'center', duration: 1500 })"
)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)
print("LoginPanel.vue done")

# File 3: api.ts
path = os.path.join(base, "src", "utils", "api.ts")
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace(
    "toast.show('登录已过期，请重新登录', 'error')",
    "toast.show('登录已过期，请重新登录', 'error', { position: 'center', duration: 3000 })"
)
content = content.replace(
    "toast.show(error.response.data.message, 'error')",
    "toast.show(error.response.data.message, 'error', { position: 'center', duration: 3000 })"
)
content = content.replace(
    "toast.show('网络错误，请稍后重试', 'error')",
    "toast.show('网络错误，请稍后重试', 'error', { position: 'center', duration: 3000 })"
)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)
print("api.ts done")

# File 4: TeacherLayout.vue - skip, already has bottom-center

# File 5: crud.ts
path = os.path.join(base, "src", "stores", "crud.ts")
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace(
    "toast.show('创建成功', 'success')\n        await fetch()",
    "toast.show('创建成功', 'success', { position: 'center', duration: 1500 })\n        await fetch()"
)
content = content.replace(
    "toast.show('更新成功', 'success')\n        await fetch()",
    "toast.show('更新成功', 'success', { position: 'center', duration: 1500 })\n        await fetch()"
)
content = content.replace(
    "toast.show('删除成功', 'success')\n        await fetch()",
    "toast.show('删除成功', 'success', { position: 'center', duration: 1500 })\n        await fetch()"
)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)
print("crud.ts done")

print("\nAll conversions complete!")
