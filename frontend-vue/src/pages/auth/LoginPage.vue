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

// 登录类型对应的渐变色
const typeColors: Record<string, string> = {
  teacher: 'linear-gradient(135deg, #6366F1, #818CF8)',
  admin: 'linear-gradient(135deg, #F59E0B, #FBBF24)',
  parent: 'linear-gradient(135deg, #10B981, #34D399)',
}
</script>

<template>
  <div class="login-page">
    <!-- 左侧：项目介绍轮播 -->
    <div class="intro">
      <!-- 背景光效 -->
      <div class="intro-orb intro-orb--top" />
      <div class="intro-orb intro-orb--bottom" />

      <div class="intro-content">
        <transition name="slide-fade" mode="out-in">
          <div :key="currentSlide" class="intro-slide">
            <div class="intro-badge">
              <span class="intro-badge-dot" />
              {{ slides[currentSlide].badge }}
            </div>

            <div class="intro-icon">{{ slides[currentSlide].icon }}</div>

            <h1 class="intro-title">
              {{ slides[currentSlide].title }}<br />
              <span class="intro-title-grad">{{ slides[currentSlide].highlight }}</span>
            </h1>

            <p class="intro-desc">{{ slides[currentSlide].desc }}</p>
          </div>
        </transition>

        <!-- 轮播指示器 -->
        <div class="intro-dots">
          <button
            v-for="(slide, i) in slides"
            :key="i"
            :class="['intro-dot', { 'intro-dot--active': currentSlide === i }]"
            @click="goToSlide(i)"
          />
        </div>

        <!-- 底部链接 -->
        <div class="intro-footer">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="intro-footer-link">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            GitHub
          </a>
          <span class="intro-footer-sep">MIT 开源许可证</span>
        </div>
      </div>
    </div>

    <!-- 右侧：登录框 -->
    <div class="login-panel">
      <div class="login-card">
        <div class="login-card-header">
          <span class="login-card-icon">🌌</span>
          <h1 class="login-card-title">学趣星球</h1>
        </div>

        <!-- 登录类型切换 -->
        <div class="login-tabs">
          <button
            v-for="t in (['teacher', 'admin', 'parent'] as const)"
            :key="t"
            :class="['login-tab', { 'login-tab--active': loginType === t }]"
            @click="loginType = t"
          >
            {{ t === 'teacher' ? '👩‍🏫 教师' : t === 'admin' ? '🔧 管理员' : '👨‍👩‍👧 家长' }}
          </button>
        </div>

        <!-- 教师登录 -->
        <div v-if="loginType === 'teacher'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="teacherUsername" class="form-input" placeholder="教师账号" @keydown.enter="focusTeacherPwd" />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input ref="teacherPwdRef" v-model="teacherPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleTeacherLogin" />
          </div>
          <button class="login-submit" :disabled="loading" @click="handleTeacherLogin">
            {{ loading ? '登录中...' : '登录' }}
          </button>

          <!-- 扫码快捷登录 -->
          <div class="login-social">
            <div class="login-social-label">
              <span class="login-social-line" /> 扫码快捷登录 <span class="login-social-line" />
            </div>
            <div class="login-social-grid">
              <button
                v-for="p in platforms" :key="p.key"
                class="login-social-btn"
                @click="handleThirdPartyLogin(p.key)"
              >
                <span class="login-social-icon" :style="{ background: p.color }">{{ p.icon }}</span>
                {{ p.label }}
              </button>
            </div>
          </div>
        </div>

        <!-- 管理员登录 -->
        <div v-if="loginType === 'admin'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="adminUsername" class="form-input" placeholder="管理员账号" @keydown.enter="focusAdminPwd" />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input ref="adminPwdRef" v-model="adminPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleAdminLogin" />
          </div>
          <button class="login-submit login-submit--amber" :disabled="loading" @click="handleAdminLogin">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>

        <!-- 家长登录 -->
        <div v-if="loginType === 'parent'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="parentUsername" class="form-input" placeholder="家长账号" @keydown.enter="focusParentPwd" />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input ref="parentPwdRef" v-model="parentPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleParentLogin" />
          </div>
          <button class="login-submit login-submit--green" :disabled="loading" @click="handleParentLogin">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ================================================================
   LOGIN PAGE — Apple 浅色主题
   左侧轮播介绍 + 右侧登录卡片
   ================================================================ */

.login-page {
  display: flex;
  justify-content: center;
  min-height: 100vh;
  background: #FFFFFF;
  color: #1D1D1F;
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'SF Pro Display', sans-serif;
  -webkit-font-smoothing: antialiased;
}

/* ── 左侧介绍轮播 ── */
.intro {
  flex: 1;
  max-width: 680px;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 64px;
  background: #FBFBFD;
}
.intro-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.25;
  pointer-events: none;
}
.intro-orb--top {
  width: 500px; height: 500px;
  background: #C7D2FE;
  top: -200px; right: -100px;
  animation: orbFloat 12s ease-in-out infinite;
}
.intro-orb--bottom {
  width: 350px; height: 350px;
  background: #A7F3D0;
  bottom: -100px; left: -50px;
  animation: orbFloat 15s ease-in-out infinite reverse;
}
@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(40px, -30px); }
}

.intro-content {
  position: relative;
  z-index: 1;
  max-width: 480px;
  width: 100%;
}

.intro-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid #E5E5EA;
  border-radius: 9999px;
  padding: 6px 16px;
  font-size: 13px;
  color: #6E6E73;
  margin-bottom: 28px;
}
.intro-badge-dot {
  width: 6px; height: 6px;
  background: #34C759;
  border-radius: 50%;
  display: inline-block;
}

.intro-icon { font-size: 56px; margin-bottom: 20px; }

.intro-title {
  font-size: 40px;
  font-weight: 900;
  color: #1D1D1F;
  line-height: 1.2;
  letter-spacing: -1px;
  margin-bottom: 20px;
}
.intro-title-grad {
  background: linear-gradient(135deg, #5E5CE6, #FF375F, #FF9F0A);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.intro-desc {
  font-size: 16px;
  color: #86868B;
  line-height: 1.7;
  white-space: pre-line;
}

/* 轮播指示器 */
.intro-dots {
  display: flex;
  gap: 8px;
  margin-top: 44px;
}
.intro-dot {
  height: 8px;
  width: 8px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  background: #D2D2D7;
  transition: all 0.3s ease;
}
.intro-dot--active {
  width: 32px;
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
}

/* 底部链接 */
.intro-footer {
  position: absolute;
  bottom: -80px;
  left: 0;
  display: flex;
  gap: 16px;
  align-items: center;
}
.intro-footer-link {
  color: #86868B;
  font-size: 13px;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: color 0.2s;
}
.intro-footer-link:hover { color: #1D1D1F; }
.intro-footer-sep { color: #AEAEB2; font-size: 13px; }

/* 轮播动画 */
.slide-fade-enter-active { transition: all 0.4s ease; }
.slide-fade-leave-active { transition: all 0.3s ease; }
.slide-fade-enter-from { opacity: 0; transform: translateX(20px); }
.slide-fade-leave-to { opacity: 0; transform: translateX(-20px); }

/* ── 右侧登录面板 ── */
.login-panel {
  flex: 0 0 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  background: #FFFFFF;
  border-left: 1px solid #F0F0F3;
  position: relative;
  z-index: 1;
}
.login-card { width: 100%; max-width: 360px; }
.login-card-header { text-align: center; margin-bottom: 28px; }
.login-card-icon  { font-size: 40px; margin-bottom: 8px; display: block; }
.login-card-title { font-size: 24px; font-weight: 800; color: #1D1D1F; }

/* 登录类型切换 */
.login-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 24px;
  background: #F5F5F7;
  border-radius: 12px;
  border: 1px solid #E5E5EA;
  padding: 4px;
}
.login-tab {
  flex: 1;
  padding: 10px 8px;
  text-align: center;
  font-size: 14px;
  font-weight: 500;
  border-radius: 9px;
  cursor: pointer;
  border: none;
  background: transparent;
  color: #86868B;
  transition: all 0.25s ease;
}
.login-tab--active {
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
  color: #FFFFFF;
  box-shadow: 0 2px 8px rgba(94, 92, 230, 0.2);
}

/* 表单 */
.login-form { display: flex; flex-direction: column; gap: 16px; }
.form-group { display: flex; fl