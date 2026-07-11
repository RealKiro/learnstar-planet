<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { apiPost } from '@/utils/api'
import type { ApiResponse, User } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

const loginType = ref<'teacher' | 'admin' | 'parent'>('teacher')
const teacherUsername = ref('')
const teacherPassword = ref('')
const adminUsername = ref('')
const adminPassword = ref('')
const parentUsername = ref('')
const parentPassword = ref('')
const loading = ref(false)

let teacherAttempts = 0
const MAX_ATTEMPTS = 3

async function handleTeacherLogin() {
  if (!teacherUsername.value.trim() || !teacherPassword.value) {
    toast.show('请输入教师账号和密码', 'error')
    return
  }
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/auth/teacher/login', {
      username: teacherUsername.value.trim(),
      password: teacherPassword.value,
    })
    teacherAttempts = 0
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'teacher-dashboard' })
  } catch (e: unknown) {
    teacherAttempts++
    const remaining = MAX_ATTEMPTS - teacherAttempts
    if (teacherAttempts >= MAX_ATTEMPTS) {
      toast.show('密码错误次数过多，请联系管理员修改密码', 'error')
    } else {
      const err = e as { response?: { data?: { message?: string } } }
      toast.show((err.response?.data?.message || '账号或密码错误') + `，请重试（还剩 ${remaining} 次）`, 'error')
    }
  } finally {
    loading.value = false
  }
}

async function handleAdminLogin() {
  if (!adminUsername.value.trim() || !adminPassword.value) {
    toast.show('请输入管理员账号和密码', 'error')
    return
  }
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/auth/admin/login', {
      username: adminUsername.value.trim(),
      password: adminPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('管理员登录成功', 'success')
    router.replace({ name: 'admin-dashboard' })
  } catch {
    // 错误已在拦截器统一处理
  } finally {
    loading.value = false
  }
}

async function handleParentLogin() {
  if (!parentUsername.value.trim() || !parentPassword.value) {
    toast.show('请输入家长账号和密码', 'error')
    return
  }
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/auth/parent/login', {
      username: parentUsername.value.trim(),
      password: parentPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('家长登录成功', 'success')
    router.replace({ name: 'parent-home' })
  } catch {
    // 错误已在拦截器统一处理
  } finally {
    loading.value = false
  }
}

const platforms = [
  { key: 'wechat', label: '微信', color: '#07C160', class: 'wechat' },
  { key: 'wechat_work', label: '企业微信', color: '#2B7CE9', class: 'wechat-work' },
  { key: 'qq', label: 'QQ', color: '#12B7F5', class: 'qq' },
  { key: 'renren', label: '人人通', color: '#FF6A00', class: 'renren' },
]

function handleThirdPartyLogin(platform: string) {
  toast.show(`正在打开${platforms.find(p => p.key === platform)?.label || platform}扫码...`, 'success')
}
</script>

<template>
  <div style="min-height:100vh;display:flex;align-items:center;justify-content:center;background:var(--gradient-hero);position:relative;overflow:hidden;">
    <!-- 装饰光球 -->
    <div style="position:absolute;width:600px;height:600px;background:radial-gradient(circle,rgba(79,70,229,0.15) 0%,transparent 70%);top:-200px;right:-200px;animation:floatOrb 8s ease-in-out infinite;"></div>
    <div style="position:absolute;width:400px;height:400px;background:radial-gradient(circle,rgba(245,158,11,0.1) 0%,transparent 70%);bottom:-100px;left:-100px;animation:floatOrb 10s ease-in-out infinite reverse;"></div>

    <div style="background:rgba(255,255,255,0.08);backdrop-filter:blur(40px) saturate(180%);border:1px solid rgba(255,255,255,0.12);border-radius:var(--radius-xl);padding:32px;width:420px;max-width:90vw;position:relative;z-index:1;box-shadow:var(--shadow-lg);">
      <div style="font-size:48px;text-align:center;margin-bottom:16px;">🌌</div>
      <h1 style="font-size:28px;font-weight:900;color:#F1F5F9;text-align:center;margin-bottom:8px;">学趣星球</h1>

      <!-- 登录类型标签 -->
      <div style="display:flex;gap:0;margin-bottom:24px;background:rgba(255,255,255,0.06);border-radius:var(--radius-md);border:1px solid rgba(255,255,255,0.12);padding:4px;">
        <button
          :style="loginType === 'teacher' ? { background:'var(--gradient-primary)', color:'#F1F5F9' } : { color:'#94A3B8' }"
          style="flex:1;padding:10px 16px;text-align:center;font-size:14px;font-weight:500;border-radius:var(--radius-sm);cursor:pointer;border:none;background:transparent;transition:all 0.2s;"
          @click="loginType = 'teacher'"
        >👩‍🏫 教师登录</button>
        <button
          :style="loginType === 'admin' ? { background:'linear-gradient(135deg,#F59E0B,#D97706)', color:'#F1F5F9' } : { color:'#94A3B8' }"
          style="flex:1;padding:10px 16px;text-align:center;font-size:14px;font-weight:500;border-radius:var(--radius-sm);cursor:pointer;border:none;background:transparent;transition:all 0.2s;"
          @click="loginType = 'admin'"
        >🔧 管理员登录</button>
        <button
          :style="loginType === 'parent' ? { background:'linear-gradient(135deg,#10B981,#059669)', color:'#F1F5F9' } : { color:'#94A3B8' }"
          style="flex:1;padding:10px 16px;text-align:center;font-size:14px;font-weight:500;border-radius:var(--radius-sm);cursor:pointer;border:none;background:transparent;transition:all 0.2s;"
          @click="loginType = 'parent'"
        >👨‍👩‍👧 家长登录</button>
      </div>

      <!-- 教师登录表单 -->
      <div v-if="loginType === 'teacher'">
        <p style="color:#94A3B8;text-align:center;margin-bottom:24px;font-size:14px;">账号密码由学校管理员统一分配，直接登录即可</p>
        <div class="form-group">
          <label style="color:#CBD5E1;">教师账号</label>
          <input v-model="teacherUsername" class="form-input" placeholder="输入管理员分配的教师账号" @keydown.enter="handleTeacherLogin">
        </div>
        <div class="form-group">
          <label style="color:#CBD5E1;">密码</label>
          <input v-model="teacherPassword" type="password" class="form-input" placeholder="输入管理员分配的密码" @keydown.enter="handleTeacherLogin">
        </div>
        <button class="btn btn-primary" style="width:100%;" :disabled="loading" @click="handleTeacherLogin">
          {{ loading ? '登录中...' : '教师登录' }}
        </button>

        <!-- 快捷登录 -->
        <div style="margin-top:24px;text-align:center;">
          <div style="color:#64748B;font-size:13px;margin-bottom:16px;display:flex;align-items:center;gap:12px;">
            <span style="flex:1;height:1px;background:rgba(255,255,255,0.1);"></span>
            扫码快捷登录
            <span style="flex:1;height:1px;background:rgba(255,255,255,0.1);"></span>
          </div>
          <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:8px;">
            <button
              v-for="p in platforms"
              :key="p.key"
              class="btn"
              style="flex-direction:column;gap:4px;padding:12px 8px;background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.1);color:#CBD5E1;font-size:12px;"
              @click="handleThirdPartyLogin(p.key)"
            >
              <div :style="{ width:'28px', height:'28px', background:p.color, borderRadius:'6px', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'14px' }">
                {{ p.key === 'wechat' ? '💬' : p.key === 'wechat_work' ? '💼' : p.key === 'qq' ? '🐧' : '🌐' }}
              </div>
              <span>{{ p.label }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 管理员登录表单 -->
      <div v-if="loginType === 'admin'">
        <div style="background:rgba(245,158,11,0.08);border:1px solid rgba(245,158,11,0.15);border-radius:var(--radius-md);padding:12px 16px;margin-bottom:16px;color:#FCD34D;font-size:13px;line-height:1.5;">
          ⚠️ 管理员账号仅用于学校后台管理（创建教师账号、班级管理、全校数据查看），不可用于日常教学操作
        </div>
        <div class="form-group">
          <label style="color:#CBD5E1;">管理员账号</label>
          <input v-model="adminUsername" class="form-input" placeholder="学校注册时分配的管理员账号" @keydown.enter="handleAdminLogin">
        </div>
        <div class="form-group">
          <label style="color:#CBD5E1;">管理员密码</label>
          <input v-model="adminPassword" type="password" class="form-input" placeholder="输入管理员密码" @keydown.enter="handleAdminLogin">
        </div>
        <button
          class="btn"
          style="background:linear-gradient(135deg,#F59E0B,#D97706);color:#F1F5F9;width:100%;box-shadow:0 4px 12px rgba(245,158,11,0.3);"
          :disabled="loading"
          @click="handleAdminLogin"
        >
          {{ loading ? '登录中...' : '管理员登录' }}
        </button>
        <p style="margin-top:24px;text-align:center;color:#64748B;font-size:12px;">
          管理员账号不支持第三方扫码登录，仅限账号密码方式
        </p>
      </div>

      <!-- 家长登录表单 -->
      <div v-if="loginType === 'parent'">
        <p style="color:#94A3B8;text-align:center;margin-bottom:24px;font-size:14px;">家长账号由学校管理员创建，可查看孩子的积分、宠物和通知</p>
        <div class="form-group">
          <label style="color:#CBD5E1;">家长账号</label>
          <input v-model="parentUsername" class="form-input" placeholder="输入管理员分配的家长账号" @keydown.enter="handleParentLogin">
        </div>
        <div class="form-group">
          <label style="color:#CBD5E1;">密码</label>
          <input v-model="parentPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleParentLogin">
        </div>
        <button
          class="btn"
          style="background:linear-gradient(135deg,#10B981,#059669);color:#F1F5F9;width:100%;box-shadow:0 4px 12px rgba(16,185,129,0.3);"
          :disabled="loading"
          @click="handleParentLogin"
        >
          {{ loading ? '登录中...' : '家长登录' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -20px); }
}
</style>
