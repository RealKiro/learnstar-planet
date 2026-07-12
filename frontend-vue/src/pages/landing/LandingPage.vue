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

if (authStore.isLoggedIn) {
  if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
  else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
  else if (authStore.isParent) router.replace({ name: 'parent-home' })
}

// ── 左侧轮播 ──
const slides = [
  {
    badge: 'MIT 开源 · 完全免费 · 自托管',
    title: '让每个孩子的努力', highlight: '都被看见',
    desc: '积分激励 · 宠物养成 · AI 助教 · 多端同步\n开源班级管理系统，数据完全自主掌控', icon: '🌌',
  },
  {
    badge: '12 大功能模块',
    title: '覆盖班级管理', highlight: '全场景',
    desc: '积分规则 · 宠物进化 · 排行榜 · 通知公告\n考勤 · 作业 · 答题 · 成绩 · 商城 · 广播 · AI', icon: '⚡',
  },
  {
    badge: '11 阶宠物进化',
    title: '积分变经验', highlight: '驱动成长',
    desc: '星尘 → 月芽 → 灵苗 → 青藤 → 慧树 → 蝶灵\n→ 鹰慧 → 狮睿 → 灵角 → 星耀 → 银河', icon: '🌟',
  },
  {
    badge: 'Docker 一键部署',
    title: '4 种数据库', highlight: '自由选择',
    desc: 'MySQL · PostgreSQL · SQLite · MariaDB\n内置 Redis 缓存，支持多端登录', icon: '🐳',
  },
]
const currentSlide = ref(0)
let slideTimer: ReturnType<typeof setInterval>

onMounted(() => {
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
})
onUnmounted(() => clearInterval(slideTimer))

function goToSlide(i: number) {
  currentSlide.value = i
  clearInterval(slideTimer)
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
}

// ── 登录逻辑 ──
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
    else { const err = e as { response?: { data?: { message?: string } } }; toast.show((err.response?.data?.message || '账号或密码错误') + '，还剩 ' + remaining + ' 次', 'error') }
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

const platforms = [
  { key: 'wechat', label: '微信', icon: '💬', color: '#07C160' },
  { key: 'wechat_work', label: '企业微信', icon: '💼', color: '#2B7CE9' },
  { key: 'qq', label: 'QQ', icon: '🐧', color: '#12B7F5' },
  { key: 'renren', label: '人人通', icon: '🌐', color: '#FF6A00' },
]

function handleThirdPartyLogin(platform: string) {
  toast.show('正在打开' + (platforms.find(p => p.key === platform)?.label || platform) + '扫码...', 'success')
}
</script>
<template>
<div class="layout">
<div class="left">
<div class="left-glow left-glow--top"></div>
<div class="left-glow left-glow--bottom"></div>
<div class="left-content">
<transition name="fade" mode="out-in">
<div :key="currentSlide" class="slide">
<div class="badge"><span class="badge-dot"></span>{{ slides[currentSlide].badge }}</div>
<div class="slide-icon">{{ slides[currentSlide].icon }}</div>
<h1 class="slide-title">{{ slides[currentSlide].title }}<br><span class="slide-gradient">{{ slides[currentSlide].highlight }}</span></h1>
<p class="slide-desc">{{ slides[currentSlide].desc }}</p>
</div>
</transition>
<div class="dots">
<button v-for="(s, i) in slides" :key="i" :class="['dot', { 'dot--active': currentSlide === i }]" @click="goToSlide(i)"></button>
</div>
<div class="left-footer"><a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="left-link">GitHub</a><span class="left-sep">MIT 开源许可证</span></div>
</div>
<div class="features">
<h2 class="features-title">12 大功能模块</h2>
<p class="features-sub">覆盖班级管理全场景，全部免费</p>
<div class="feature-grid">
<div v-for="f in features" :key="f.title" class="feat">
<span class="feat-icon">{{ f.icon }}</span>
<div><div class="feat-name">{{ f.title }}</div><div class="feat-desc">{{ f.desc }}</div></div>
</div>
</div>
<h2 class="features-title" style="margin-top:48px">11 阶宠物进化</h2>
<p class="features-sub">积分变经验，从星尘到银河</p>
<div class="evo">
<div v-for="(s, i) in stages" :key="s.name" style="display:flex;align-items:center;gap:2px">
<span v-if="i > 0" class="evo-arrow">→</span>
<div class="evo-item"><span class="evo-emoji">{{ s.emoji }}</span><span class="evo-name">{{ s.name }}</span></div>
</div>
</div>
<div class="cta-box">
<div class="cta-icon">🚀</div>
<h3>加入开源社区</h3>
<p>项目完全开源，欢迎 Star / Fork / Issue / PR</p>
<div class="cta-links">
<a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="cta-btn cta-btn--dark">⭐ 给个 Star</a>
<a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" class="cta-btn cta-btn--ghost">提交 Issue</a>
</div>
</div>
<footer class="foot"><a href="https://github.com/RealKiro/learnstar-planet" target="_blank">学趣星球</a> · MIT 开源许可证 · 自托管 · 完全免费</footer>
</div>
</div>
<div class="right">
<div class="login-card">
<div class="login-header"><span class="login-icon">🌌</span><h1 class="login-brand">学趣星球</h1></div>
<div class="login-tabs">
<button v-for="t in (['teacher', 'admin', 'parent'] as const)" :key="t" :class="['login-tab', { 'login-tab--active': loginType === t }]" @click="loginType = t">{{ t === 'teacher' ? '👩‍🏫 教师' : t === 'admin' ? '🔧 管理员' : '👨‍👩‍👧 家长' }}</button>
</div>
<div v-if="loginType === 'teacher'" class="login-form">
<div class="field"><label>账号</label><input v-model="teacherUsername" placeholder="教师账号" @keydown.enter="focusTeacherPwd"></div>
<div class="field"><label>密码</label><input ref="teacherPwdRef" v-model="teacherPassword" type="password" placeholder="输入密码" @keydown.enter="handleTeacherLogin"></div>
<button class="btn-login" :disabled="loading" @click="handleTeacherLogin">{{ loading ? '登录中...' : '登录' }}</button>
<div class="social"><div class="social-label"><span></span> 扫码快捷登录 <span></span></div>
<div class="social-btns">
<button v-for="p in platforms" :key="p.key" @click="handleThirdPartyLogin(p.key)"><span :style="{ background: p.color }">{{ p.icon }}</span>{{ p.label }}</button>
</div></div>
</div>
<div v-if="loginType === 'admin'" class="login-form">
<div class="field"><label>账号</label><input v-model="adminUsername" placeholder="管理员账号" @keydown.enter="focusAdminPwd"></div>
<div class="field"><label>密码</label><input ref="adminPwdRef" v-model="adminPassword" type="password" placeholder="输入密码" @keydown.enter="handleAdminLogin"></div>
<button class="btn-login btn-login--amber" :disabled="loading" @click="handleAdminLogin">{{ loading ? '登录中...' : '登录' }}</button>
</div>
<div v-if="loginType === 'parent'" class="login-form">
<div class="field"><label>账号</label><input v-model="parentUsername" placeholder="家长账号" @keydown.enter="focusParentPwd"></div>
<div class="field"><label>密码</label><input ref="parentPwdRef" v-model="parentPassword" type="password" placeholder="输入密码" @keydown.enter="handleParentLogin"></div>
<button class="btn-login btn-login--green" :disabled="loading" @click="handleParentLogin">{{ loading ? '登录中...' : '登录' }}</button>
</div>
</div>
</div>
</div>
</template>
<style scoped>
/* ================================================================
   LANDING PAGE — 左右分栏
   左侧：可滚动介绍（轮播 + 功能 + 进化 + CTA）
   右侧：固定登录面板
   ================================================================ */

* { box-sizing: border-box; }

.layout {
  display: flex;
  min-height: 100vh;
  background: #FFFFFF;
  color: #1D1D1F;
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'SF Pro Display', sans-serif;
  -webkit-font-smoothing: antialiased;
}

/* ═══════════ 左侧 ═══════════ */
.left {
  flex: 1;
  max-height: 100vh;
  overflow-y: auto;
  position: relative;
  background: #FBFBFD;
}
.left::-webkit-scrollbar { width: 4px; }
.left::-webkit-scrollbar-thumb { background: #D2D2D7; border-radius: 4px; }

.left-glow {
  position: fixed; border-radius: 50%; filter: blur(100px);
  opacity: 0.2; pointer-events: none; z-index: 0;
}
.left-glow--top    { width: 500px; height: 500px; background: #C7D2FE; top: -200px; right: -100px; animation: glow 12s ease-in-out infinite; }
.left-glow--bottom { width: 350px; height: 350px; background: #A7F3D0; bottom: -100px; left: -50px; animation: glow 15s ease-in-out infinite reverse; }
@keyframes glow {
  0%, 100% { transform: translate(0, 0); }
  50%      { transform: translate(40px, -30px); }
}

.left-content {
  position: relative; z-index: 1;
  min-height: 100vh;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 60px 64px 80px;
}

/* 轮播 */
.slide { max-width: 440px; text-align: left; }
.badge {
  display: inline-flex; align-items: center; gap: 8px;
  background: #FFFFFF; border: 1px solid #E5E5EA;
  border-radius: 9999px; padding: 6px 16px;
  font-size: 13px; color: #6E6E73; margin-bottom: 28px;
}
.badge-dot { width: 6px; height: 6px; background: #34C759; border-radius: 50%; display: inline-block; box-shadow: 0 0 6px rgba(52, 199, 89, 0.3); }
.slide-icon { font-size: 56px; margin-bottom: 20px; }
.slide-title {
  font-size: 40px; font-weight: 900; line-height: 1.2;
  letter-spacing: -1px; color: #1D1D1F; margin-bottom: 20px;
}
.slide-gradient {
  background: linear-gradient(135deg, #5E5CE6, #FF375F, #FF9F0A);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
}
.slide-desc { font-size: 16px; color: #86868B; line-height: 1.7; white-space: pre-line; }

/* 轮播点 */
.dots { display: flex; gap: 8px; margin-top: 44px; }
.dot {
  height: 8px; border-radius: 4px; border: none; cursor: pointer;
  background: #D2D2D7; transition: all 0.3s;
  width: 8px;
}
.dot--active { width: 32px; background: linear-gradient(135deg, #5E5CE6, #818CF8); }

/* 底部 */
.left-footer { position: absolute; bottom: 40px; left: 64px; display: flex; gap: 16px; align-items: center; z-index: 1; }
.left-link { color: #86868B; font-size: 13px; text-decoration: none; transition: color 0.2s; }
.left-link:hover { color: #1D1D1F; }
.left-sep { color: #AEAEB2; font-size: 13px; }

/* 轮播动画 */
.fade-enter-active { transition: all 0.4s ease; }
.fade-leave-active { transition: all 0.3s ease; }
.fade-enter-from { opacity: 0; transform: translateX(20px); }
.fade-leave-to   { opacity: 0; transform: translateX(-20px); }

/* 功能区域 */
.features { padding: 0 64px 60px; position: relative; z-index: 1; }
.features-title { font-size: 24px; font-weight: 800; color: #1D1D1F; margin-bottom: 8px; }
.features-sub   { font-size: 15px; color: #86868B; margin-bottom: 32px; }
.feature-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}
.feat {
  display: flex; align-items: flex-start; gap: 12px;
  padding: 16px; background: #FFFFFF; border: 1px solid #F0F0F3;
  border-radius: 14px; transition: all 0.2s;
}
.feat:hover { border-color: #E5E5EA; box-shadow: 0 4px 16px rgba(0,0,0,0.04); transform: translateY(-2px); }
.feat-icon { font-size: 24px; flex-shrink: 0; width: 36px; text-align: center; }
.feat-name { font-size: 14px; font-weight: 700; color: #1D1D1F; margin-bottom: 3px; }
.feat-desc { font-size: 12px; color: #86868B; line-height: 1.5; }

/* 进化 */
.evo { display: flex; align-items: center; flex-wrap: wrap; gap: 2px; padding: 16px 0; }
.evo-arrow { color: #C7C7CC; font-size: 13px; user-select: none; }
.evo-item {
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  padding: 8px 12px; border-radius: 14px; transition: all 0.2s; cursor: default;
}
.evo-item:hover { background: #F0F0F3; transform: scale(1.08); }
.evo-emoji { font-size: 26px; }
.evo-name  { font-size: 11px; color: #86868B; font-weight: 500; }

/* CTA */
.cta-box {
  margin-top: 48px; text-align: center; padding: 40px;
  background: #FFFFFF; border: 1px solid #F0F0F3;
  border-radius: 20px;
}
.cta-icon { font-size: 40px; margin-bottom: 12px; }
.cta-box h3 { font-size: 22px; font-weight: 800; color: #1D1D1F; margin-bottom: 8px; }
.cta-box p  { font-size: 14px; color: #86868B; margin-bottom: 24px; }
.cta-links { display: flex; gap: 10px; justify-content: center; }
.cta-btn {
  padding: 10px 24px; border-radius: 9999px; font-size: 14px; font-weight: 600;
  text-decoration: none; transition: all 0.2s;
}
.cta-btn--dark  { background: #1D1D1F; color: #FFFFFF; }
.cta-btn--dark:hover  { background: #333; }
.cta-btn--ghost { border: 1px solid #D2D2D7; color: #1D1D1F; }
.cta-btn--ghost:hover { background: #F5F5F7; }

/* Footer */
.foot {
  margin-top: 40px; padding-bottom: 40px;
  text-align: center; font-size: 13px; color: #AEAEB2;
}
.foot a { color: #86868B; text-decoration: none; }
.foot a:hover { color: #1D1D1F; }

/* ═══════════ 右侧：固定登录 ═══════════ */
.right {
  flex: 0 0 400px;
  display: flex; align-items: center; justify-content: center;
  padding: 32px;
  background: #FFFFFF;
  border-left: 1px solid #F0F0F3;
  position: sticky; top: 0; height: 100vh;
}
.login-card { width: 100%; max-width: 360px; }
.login-header { text-align: center; margin-bottom: 28px; }
.login-icon  { font-size: 40px; margin-bottom: 8px; display: block; }
.login-brand { font-size: 24px; font-weight: 800; color: #1D1D1F; }

.login-tabs {
  display: flex; gap: 0; margin-bottom: 24px;
  background: #F5F5F7; border