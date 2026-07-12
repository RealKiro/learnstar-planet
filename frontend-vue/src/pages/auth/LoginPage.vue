<script setup lang="ts">
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { apiPost } from '@/utils/api'
import type { ApiResponse, User } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

// ===== 登录逻辑 =====
const loginType = ref<'teacher' | 'admin' | 'parent'>('teacher')
const teacherUsername = ref('')
const teacherPassword = ref('')
const adminUsername = ref('')
const adminPassword = ref('')
const parentUsername = ref('')
const parentPassword = ref('')
const loading = ref(false)

const teacherPwdRef = ref<HTMLInputElement>()
const adminPwdRef = ref<HTMLInputElement>()
const parentPwdRef = ref<HTMLInputElement>()

let teacherAttempts = 0
const MAX_ATTEMPTS = 3

function focusTeacherPwd() { if (teacherUsername.value.trim()) nextTick(() => teacherPwdRef.value?.focus()) }
function focusAdminPwd() { if (adminUsername.value.trim()) nextTick(() => adminPwdRef.value?.focus()) }
function focusParentPwd() { if (parentUsername.value.trim()) nextTick(() => parentPwdRef.value?.focus()) }

async function handleTeacherLogin() {
  if (!teacherUsername.value.trim() || !teacherPassword.value) { toast.show('请输入账号和密码', 'error'); return }
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/auth/teacher/login', {
      username: teacherUsername.value.trim(), password: teacherPassword.value,
    })
    teacherAttempts = 0
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'teacher-dashboard' })
  } catch (e: unknown) {
    teacherAttempts++
    const remaining = MAX_ATTEMPTS - teacherAttempts
    if (teacherAttempts >= MAX_ATTEMPTS) toast.show('密码错误次数过多，请联系管理员', 'error')
    else { const err = e as { response?: { data?: { message?: string } } }; toast.show((err.response?.data?.message || '账号或密码错误') + `，还剩 ${remaining} 次`, 'error') }
  } finally { loading.value = false }
}

async function handleAdminLogin() {
  if (!adminUsername.value.trim() || !adminPassword.value) { toast.show('请输入账号和密码', 'error'); return }
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/auth/admin/login', {
      username: adminUsername.value.trim(), password: adminPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'admin-dashboard' })
  } catch { /* handled */ } finally { loading.value = false }
}

async function handleParentLogin() {
  if (!parentUsername.value.trim() || !parentPassword.value) { toast.show('请输入账号和密码', 'error'); return }
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/auth/parent/login', {
      username: parentUsername.value.trim(), password: parentPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'parent-home' })
  } catch { /* handled */ } finally { loading.value = false }
}

// ===== 第三方登录（仅教师） =====
const platforms = [
  { key: 'wechat', label: '微信', icon: '💬', color: '#07C160' },
  { key: 'wechat_work', label: '企业微信', icon: '💼', color: '#2B7CE9' },
  { key: 'qq', label: 'QQ', icon: '🐧', color: '#12B7F5' },
  { key: 'renren', label: '人人通', icon: '🌐', color: '#FF6A00' },
]

function handleThirdPartyLogin(platform: string) {
  toast.show(`正在打开${platforms.find(p => p.key === platform)?.label || platform}扫码...`, 'success')
}

// ===== 左侧轮播 =====
const slides = [
  {
    badge: 'MIT 开源 · 完全免费 · 自托管',
    title: '让每个孩子的努力',
    highlight: '都被看见',
    desc: '积分激励 · 宠物养成 · AI 助教 · 多端同步\n开源班级管理系统，数据完全自主掌控',
    icon: '🌌',
  },
  {
    badge: '12 大功能模块',
    title: '覆盖班级管理',
    highlight: '全场景',
    desc: '积分规则 · 宠物进化 · 排行榜 · 通知公告\n考勤 · 作业 · 答题 · 成绩 · 商城 · 广播 · AI',
    icon: '⚡',
  },
  {
    badge: '11 阶宠物进化',
    title: '积分变经验',
    highlight: '驱动成长',
    desc: '星尘 → 月芽 → 灵苗 → 青藤 → 慧树 → 蝶灵\n→ 鹰慧 → 狮睿 → 灵角 → 星耀 → 银河',
    icon: '🌟',
  },
  {
    badge: 'Docker 一键部署',
    title: '4 种数据库',
    highlight: '自由选择',
    desc: 'MySQL · PostgreSQL · SQLite · MariaDB\n内置 Redis 缓存，支持多端登录',
    icon: '🐳',
  },
]
const currentSlide = ref(0)
let slideTimer: ReturnType<typeof setInterval>

onMounted(() => {
  slideTimer = setInterval(() => {
    currentSlide.value = (currentSlide.value + 1) % slides.length
  }, 5000)
})
onUnmounted(() => clearInterval(slideTimer))

function goToSlide(i: number) {
  currentSlide.value = i
  clearInterval(slideTimer)
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
}
</script>

<template>
  <div style="min-height:100vh;display:flex;background:#0F172A;">
    <!-- 左侧：项目介绍轮播 -->
    <div class="intro-panel" style="flex:1;position:relative;overflow:hidden;display:flex;align-items:center;justify-content:center;padding:48px 64px;">
      <!-- 背景光效 -->
      <div style="position:absolute;width:600px;height:600px;background:radial-gradient(circle,rgba(79,70,229,0.12) 0%,transparent 70%);top:-200px;right:-100px;animation:floatOrb 12s ease-in-out infinite;"></div>
      <div style="position:absolute;width:400px;height:400px;background:radial-gradient(circle,rgba(16,185,129,0.08) 0%,transparent 70%);bottom:-100px;left:-50px;animation:floatOrb 15s ease-in-out infinite reverse;"></div>

      <!-- 轮播内容 -->
      <div style="position:relative;z-index:1;max-width:480px;width:100%;">
        <transition name="slide-fade" mode="out-in">
          <div :key="currentSlide" style="text-align:left;">
            <!-- 开源徽章 -->
            <div style="display:inline-flex;align-items:center;gap:8px;background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.12);border-radius:100px;padding:6px 16px;font-size:13px;color:#94A3B8;margin-bottom:24px;">
              <span style="color:#10B981;">●</span>
              {{ slides[currentSlide].badge }}
            </div>

            <div style="font-size:56px;margin-bottom:16px;">{{ slides[currentSlide].icon }}</div>

            <h1 style="font-size:40px;font-weight:900;color:#F1F5F9;margin-bottom:16px;line-height:1.2;letter-spacing:-0.5px;">
              {{ slides[currentSlide].title }}<br>
              <span style="background:linear-gradient(135deg,#818CF8,#6EE7B7);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text;">{{ slides[currentSlide].highlight }}</span>
            </h1>

            <p style="font-size:16px;color:#64748B;line-height:1.7;white-space:pre-line;">{{ slides[currentSlide].desc }}</p>
          </div>
        </transition>

        <!-- 轮播指示器 -->
        <div style="display:flex;gap:8px;margin-top:40px;">
          <button
            v-for="(slide, i) in slides" :key="i"
            @click="goToSlide(i)"
            :style="currentSlide === i ? { width:'32px', background:'linear-gradient(135deg,#6366F1,#818CF8)' } : { width:'8px', background:'rgba(255,255,255,0.2)' }"
            style="height:8px;border-radius:4px;border:none;cursor:pointer;transition:all 0.3s;"
          />
        </div>

        <!-- 底部开源链接 -->
        <div style="position:absolute;bottom:-80px;left:0;display:flex;gap:16px;align-items:center;">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="color:#475569;font-size:13px;text-decoration:none;display:flex;align-items:center;gap:6px;transition:color 0.2s;" onmouseover="this.style.color='#94A3B8'" onmouseout="this.style.color='#475569'">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            GitHub
          </a>
          <span style="color:#334155;font-size:13px;">MIT 开源许可证</span>
        </div>
      </div>
    </div>

    <!-- 右侧：登录框 -->
    <div style="width:400px;flex-shrink:0;display:flex;align-items:center;justify-content:center;padding:32px;background:rgba(255,255,255,0.03);border-left:1px solid rgba(255,255,255,0.06);position:relative;z-index:1;">
      <div style="width:100%;max-width:360px;">
        <div style="font-size:40px;text-align:center;margin-bottom:12px;">🌌</div>
        <h1 style="font-size:24px;font-weight:800;color:#F1F5F9;text-align:center;margin-bottom:24px;">学趣星球</h1>

        <!-- 登录类型标签 -->
        <div style="display:flex;gap:0;margin-bottom:24px;background:rgba(255,255,255,0.06);border-radius:10px;border:1px solid rgba(255,255,255,0.12);padding:4px;">
          <button
            :style="loginType === 'teacher' ? { background:'linear-gradient(135deg,#6366F1,#818CF8)', color:'#F1F5F9' } : { color:'#94A3B8' }"
            style="flex:1;padding:10px 8px;text-align:center;font-size:14px;font-weight:500;border-radius:7px;cursor:pointer;border:none;background:transparent;transition:all 0.2s;"
            @click="loginType = 'teacher'"
          >👩‍🏫 教师</button>
          <button
            :style="loginType === 'admin' ? { background:'linear-gradient(135deg,#F59E0B,#D97706)', color:'#F1F5F9' } : { color:'#94A3B8' }"
            style="flex:1;padding:10px 8px;text-align:center;font-size:14px;font-weight:500;border-radius:7px;cursor:pointer;border:none;background:transparent;transition:all 0.2s;"
            @click="loginType = 'admin'"
          >🔧 管理员</button>
          <button
            :style="loginType === 'parent' ? { background:'linear-gradient(135deg,#10B981,#059669)', color:'#F1F5F9' } : { color:'#94A3B8' }"
            style="flex:1;padding:10px 8px;text-align:center;font-size:14px;font-weight:500;border-radius:7px;cursor:pointer;border:none;background:transparent;transition:all 0.2s;"
            @click="loginType = 'parent'"
          >👨‍👩‍👧 家长</button>
        </div>

        <!-- 教师登录 -->
        <div v-if="loginType === 'teacher'">
          <div class="form-group">
            <label style="color:#CBD5E1;font-size:13px;">账号</label>
            <input v-model="teacherUsername" class="form-input" placeholder="教师账号" @keydown.enter="focusTeacherPwd">
          </div>
          <div class="form-group">
            <label style="color:#CBD5E1;font-size:13px;">密码</label>
            <input ref="teacherPwdRef" v-model="teacherPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleTeacherLogin">
          </div>
          <button class="btn btn-primary" style="width:100%;background:linear-gradient(135deg,#6366F1,#818CF8);border:none;padding:12px;" :disabled="loading" @click="handleTeacherLogin">
            {{ loading ? '登录中...' : '登录' }}
          </button>

          <!-- 扫码快捷登录（仅教师） -->
          <div style="margin-top:20px;text-align:center;">
            <div style="color:#475569;font-size:12px;margin-bottom:12px;display:flex;align-items:center;gap:8px;">
              <span style="flex:1;height:1px;background:rgba(255,255,255,0.08);"></span>
              扫码快捷登录
              <span style="flex:1;height:1px;background:rgba(255,255,255,0.08);"></span>
            </div>
            <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:6px;">
              <button
                v-for="p in platforms" :key="p.key"
                style="display:flex;flex-direction:column;align-items:center;gap:4px;padding:10px 4px;background:rgba(255,255,255,0.04);border:1px solid rgba(255,255,255,0.08);border-radius:10px;color:#94A3B8;font-size:11px;cursor:pointer;transition:all 0.2s;"
                onmouseover="this.style.background='rgba(255,255,255,0.08)';this.style.color='#CBD5E1'"
                onmouseout="this.style.background='rgba(255,255,255,0.04)';this.style.color='#94A3B8'"
                @click="handleThirdPartyLogin(p.key)"
              >
                <span :style="{ display:'flex', alignItems:'center', justifyContent:'center', width:'28px', height:'28px', background:p.color, borderRadius:'8px', fontSize:'14px' }">{{ p.icon }}</span>
                {{ p.label }}
              </button>
            </div>
          </div>
        </div>

        <!-- 管理员登录 -->
        <div v-if="loginType === 'admin'">
          <div class="form-group">
            <label style="color:#CBD5E1;font-size:13px;">账号</label>
            <input v-model="adminUsername" class="form-input" placeholder="管理员账号" @keydown.enter="focusAdminPwd">
          </div>
          <div class="form-group">
            <label style="color:#CBD5E1;font-size:13px;">密码</label>
            <input ref="adminPwdRef" v-model="adminPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleAdminLogin">
          </div>
          <button class="btn" style="width:100%;background:linear-gradient(135deg,#F59E0B,#D97706);color:#F1F5F9;border:none;padding:12px;" :disabled="loading" @click="handleAdminLogin">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>

        <!-- 家长登录 -->
        <div v-if="loginType === 'parent'">
          <div class="form-group">
            <label style="color:#CBD5E1;font-size:13px;">账号</label>
            <input v-model="parentUsername" class="form-input" placeholder="家长账号" @keydown.enter="focusParentPwd">
          </div>
          <div class="form-group">
            <label style="color:#CBD5E1;font-size:13px;">密码</label>
            <input ref="parentPwdRef" v-model="parentPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleParentLogin">
          </div>
          <button class="btn" style="width:100%;background:linear-gradient(135deg,#10B981,#059669);color:#F1F5F9;border:none;padding:12px;" :disabled="loading" @click="handleParentLogin">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 移动端隐藏左侧（@media 通过内联无法实现，用 JS 判断） -->
  </div>
</template>

<style scoped>
@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(40px, -30px); }
}
.slide-fade-enter-active { transition: all 0.4s ease; }
.slide-fade-leave-active { transition: all 0.3s ease; }
.slide-fade-enter-from { opacity: 0; transform: translateX(20px); }
.slide-fade-leave-to { opacity: 0; transform: translateX(-20px); }

/* 移动端：隐藏左侧介绍区 */
@media (max-width: 768px) {
  .intro-panel { display: none !important; }
}
</style>
