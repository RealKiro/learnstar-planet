<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
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

// ── 滚动揭示动画 ──
const observedEls = ref<Set<HTMLElement>>(new Set())
let observer: IntersectionObserver | null = null

function setRevealRef(el: HTMLElement | null) {
  if (el && !observedEls.value.has(el) && observer) {
    observer.observe(el)
    observedEls.value.add(el)
  }
}

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          entry.target.classList.add('revealed')
          observer!.unobserve(entry.target)
          observedEls.value.delete(entry.target as HTMLElement)
        }
      }
    },
    { threshold: 0.12, rootMargin: '0px 0px -30px 0px' },
  )
})

onUnmounted(() => observer?.disconnect())

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

function scrollToLogin() {
  document.getElementById('login-section')?.scrollIntoView({ behavior: 'smooth' })
}
</script>

<template>
  <div class="page">
    <!-- ===== 导航栏 ===== -->
    <header class="nav">
      <div class="nav-inner">
        <a href="/" class="nav-brand">🌌&nbsp;学趣星球</a>
        <nav class="nav-links">
          <a href="#features">功能</a>
          <a href="#evolution">宠物进化</a>
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="nav-github">GitHub</a>
          <button class="nav-cta" @click="scrollToLogin">登录</button>
        </nav>
      </div>
    </header>

    <!-- ===== Hero ===== -->
    <section class="hero">
      <div class="hero-inner reveal" :ref="setRevealRef">
        <span class="badge">
          <span class="badge-dot"></span> MIT 开源 · 完全免费 · 自托管
        </span>
        <h1 class="hero-title">
          让每个孩子的努力<br>
          <span class="hero-title-grad">都被看见</span>
        </h1>
        <p class="hero-desc">
          积分激励 · 宠物养成 · AI 助教 · 多端同步<br>
          开源班级管理系统，Docker 一键部署，数据完全自主掌控
        </p>
        <div class="hero-actions">
          <button class="btn-primary" @click="scrollToLogin">立即使用</button>
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="btn-secondary">源码</a>
        </div>
        <div class="hero-tags">
          <div v-for="item in highlights" :key="item.label" class="tag">
            <span class="tag-icon">{{ item.icon }}</span>
            <span class="tag-label">{{ item.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 登录区 ===== -->
    <section id="login-section" class="section-gray">
      <div class="container">
        <div class="login-card reveal" :ref="setRevealRef">
          <div class="login-card-header">
            <span class="login-card-icon">🌌</span>
            <h2 class="login-card-title">学趣星球</h2>
          </div>

          <div class="login-tabs">
            <button
              v-for="t in (['teacher', 'admin', 'parent'] as const)" :key="t"
              :class="['login-tab', { 'login-tab--active': loginType === t }]"
              @click="loginType = t"
            >{{ t === 'teacher' ? '👩‍🏫 教师' : t === 'admin' ? '🔧 管理员' : '👨‍👩‍👧 家长' }}</button>
          </div>

          <div v-if="loginType === 'teacher'" class="login-form">
            <div class="form-group">
              <label>账号</label>
              <input v-model="teacherUsername" class="form-input" placeholder="教师账号" @keydown.enter="focusTeacherPwd">
            </div>
            <div class="form-group">
              <label>密码</label>
              <input ref="teacherPwdRef" v-model="teacherPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleTeacherLogin">
            </div>
            <button class="login-submit" :disabled="loading" @click="handleTeacherLogin">{{ loading ? '登录中...' : '登录' }}</button>
            <div class="login-social">
              <div class="login-social-label"><span class="login-social-line"></span> 扫码快捷登录 <span class="login-social-line"></span></div>
              <div class="login-social-grid">
                <button v-for="p in platforms" :key="p.key" class="login-social-btn" @click="handleThirdPartyLogin(p.key)">
                  <span class="login-social-icon" :style="{ background: p.color }">{{ p.icon }}</span>
                  {{ p.label }}
                </button>
              </div>
            </div>
          </div>

          <div v-if="loginType === 'admin'" class="login-form">
            <div class="form-group">
              <label>账号</label>
              <input v-model="adminUsername" class="form-input" placeholder="管理员账号" @keydown.enter="focusAdminPwd">
            </div>
            <div class="form-group">
              <label>密码</label>
              <input ref="adminPwdRef" v-model="adminPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleAdminLogin">
            </div>
            <button class="login-submit login-submit--amber" :disabled="loading" @click="handleAdminLogin">{{ loading ? '登录中...' : '登录' }}</button>
          </div>

          <div v-if="loginType === 'parent'" class="login-form">
            <div class="form-group">
              <label>账号</label>
              <input v-model="parentUsername" class="form-input" placeholder="家长账号" @keydown.enter="focusParentPwd">
            </div>
            <div class="form-group">
              <label>密码</label>
              <input ref="parentPwdRef" v-model="parentPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleParentLogin">
            </div>
            <button class="login-submit login-submit--green" :disabled="loading" @click="handleParentLogin">{{ loading ? '登录中...' : '登录' }}</button>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 功能模块 ===== -->
    <section id="features" class="section-white">
      <div class="container">
        <div class="section-head reveal" :ref="setRevealRef">
          <span class="section-eyebrow">功能模块</span>
          <h2 class="section-title">12 大模块，覆盖全场景</h2>
          <p class="section-desc">全部免费，无任何付费功能，Docker 一条命令即可部署</p>
        </div>
        <div class="bento">
          <div v-for="feat in features" :key="feat.title" class="bento-card reveal" :class="feat.size" :ref="setRevealRef">
            <span class="bento-icon">{{ feat.icon }}</span>
            <h3 class="bento-title">{{ feat.title }}</h3>
            <p class="bento-desc">{{ feat.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 宠物进化 ===== -->
    <section id="evolution" class="section-gray">
      <div class="container">
        <div class="section-head reveal" :ref="setRevealRef">
          <span class="section-eyebrow">宠物系统</span>
          <h2 class="section-title">11 阶进化，积分驱动成长</h2>
          <p class="section-desc">从星尘到银河，每一次进步都有回响</p>
        </div>
        <div class="evo-track reveal" :ref="setRevealRef">
          <template v-for="(stage, i) in evolutionStages" :key="stage.name">
            <span v-if="i > 0" class="evo-arrow">→</span>
            <div class="evo-node">
              <span class="evo-emoji">{{ stage.emoji }}</span>
              <span class="evo-name">{{ stage.name }}</span>
            </div>
          </template>
        </div>
      </div>
    </section>

    <!-- ===== CTA ===== -->
    <section class="section-white">
      <div class="container">
        <div class="cta reveal" :ref="setRevealRef">
          <div class="cta-icon">🚀</div>
          <h2 class="cta-title">加入开源社区</h2>
          <p class="cta-desc">
            项目完全开源，欢迎 Star ⭐、Fork、提交 Issue 和 PR<br>
            如果你正在用它管理班级，你就是这个社区的一员
          </p>
          <div class="cta-btns">
            <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="btn-outline">⭐ 给个 Star</a>
            <a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" class="btn-ghost">提交 Issue</a>
          </div>
        </div>
      </div>
    </section>

    <footer class="footer">
      <p class="footer-line">
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="footer-link">学趣星球</a>
        &nbsp;· MIT 开源许可证 · 自托管 · 完全免费
      </p>
      <p class="footer-made">Made with 💜 by open source community</p>
    </footer>
  </div>
</template>

<script lang="ts">
const highlights = [
  { icon: '🆓', label: '完全免费' },
  { icon: '🔓', label: 'MIT 开源' },
  { icon: '🏠', label: '自托管部署' },
  { icon: '🗄️', label: '4 种数据库' },
  { icon: '🐳', label: 'Docker 一键部署' },
  { icon: '📱', label: '多端适配' },
]

const features = [
  { icon: '⭐', title: '积分激励', desc: '自定义规则，实时加减分，进步看得见', size: 'bento-lg' },
  { icon: '🌟', title: '宠物进化', desc: '11 阶宇宙进化体系，积分变经验驱动成长', size: '' },
  { icon: '🏆', title: '排行竞技', desc: '总积分 / 周进步 / 宠物等级三大排行', size: '' },
  { icon: '📢', title: '班级通知', desc: '一键发布，实时推送家长端', size: 'bento-wide' },
  { icon: '📊', title: '成绩管理', desc: '成绩录入分析，班级对比，趋势可视化', size: '' },
  { icon: '🤖', title: 'AI 助教', desc: '生成班级反馈、学生分析、家校沟通建议', size: 'bento-tall' },
  { icon: '✅', title: '智能考勤', desc: '一键点名签到，到课 / 请假 / 迟到状态统计', size: '' },
  { icon: '📱', title: '扫码收作业', desc: '生成专属二维码，学生扫码提交，自动汇总', size: 'bento-wide' },
  { icon: '🛍️', title: '积分商城', desc: '学生兑换实物 / 特权奖励，教师审批发放', size: '' },
  { icon: '📡', title: '实时广播', desc: '消息直达教室桌面端，支持文字 / 语音 / 横幅 / 全屏', size: '' },
  { icon: '📝', title: '在线答题', desc: '题库管理 + 课堂即时检测，自动判分统计', size: '' },
  { icon: '🔗', title: '多端登录', desc: '微信 / 企微 / QQ / 人人通扫码，账号密码双通道', size: 'bento-lg' },
]

const evolutionStages = [
  { emoji: '🌟', name: '星尘' }, { emoji: '🌙', name: '月芽' },
  { emoji: '🌱', name: '灵苗' }, { emoji: '🌿', name: '青藤' },
  { emoji: '🌳', name: '慧树' }, { emoji: '🦋', name: '蝶灵' },
  { emoji: '🦅', name: '鹰慧' }, { emoji: '🦁', name: '狮睿' },
  { emoji: '🦄', name: '灵角' }, { emoji: '✨', name: '星耀' },
  { emoji: '🌌', name: '银河' },
]
</script>

<style scoped>
/* ================================================================
   LANDING + LOGIN — 单页 Apple 浅色风格
   Hero → 登录 → 功能卡片 → 进化 → CTA
   ================================================================ */

.page {
  min-height: 100vh;
  background: #FFFFFF;
  color: #1D1D1F;
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'SF Pro Display', sans-serif;
  -webkit-font-smoothing: antialiased;
}

/* ── 导航栏 ── */
.nav {
  position: fixed; top: 0; left: 0; right: 0; z-index: 100;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.nav-inner {
  max-width: 1100px; margin: 0 auto; padding: 0 24px;
  height: 52px;
  display: flex; align-items: center; justify-content: space-between;
}
.nav-brand {
  font-size: 18px; font-weight: 700; color: #1D1D1F;
  text-decoration: none; letter-spacing: -0.3px;
}
.nav-links { display: flex; align-items: center; gap: 20px; }
.nav-links a { font-size: 13px; font-weight: 500; color: #6E6E73; text-decoration: none; transition: color 0.2s; }
.nav-links a:hover { color: #1D1D1F; }
.nav-github { display: inline-flex; align-items: center; gap: 5px; }
.nav-cta {
  display: inline-block; padding: 7px 18px; border-radius: 9999px;
  background: #1D1D1F; color: #FFFFFF;
  font-size: 13px; font-weight: 600;
  border: none; cursor: pointer;
  transition: all 0.2s;
}
.nav-cta:hover { background: #333; transform: scale(1.03); }

/* ── Hero ── */
.hero {
  padding: 140px 20px 80px;
  background: linear-gradient(180deg, #FBFBFD 0%, #FFFFFF 100%);
}
.hero-inner { max-width: 660px; margin: 0 auto; text-align: center; }
.badge {
  display: inline-flex; align-items: center; gap: 8px;
  background: #F5F5F7; border: 1px solid #E5E5EA;
  border-radius: 9999px; padding: 6px 18px;
  font-size: 13px; color: #6E6E73; margin-bottom: 32px;
}
.badge-dot {
  width: 7px; height: 7px; background: #34C759; border-radius: 50%;
  display: inline-block; box-shadow: 0 0 6px rgba(52, 199, 89, 0.35);
  animation: pulse-dot 2.5s ease-in-out infinite;
}
@keyframes pulse-dot {
  0%, 100% { box-shadow: 0 0 4px rgba(52, 199, 89, 0.3); }
  50%      { box-shadow: 0 0 12px rgba(52, 199, 89, 0.55); }
}
.hero-title {
  font-size: clamp(40px, 6vw, 60px); font-weight: 800;
  line-height: 1.12; letter-spacing: -1.5px;
  color: #1D1D1F; margin-bottom: 24px;
}
.hero-title-grad {
  background: linear-gradient(135deg, #5E5CE6 0%, #FF375F 50%, #FF9F0A 100%);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
}
.hero-desc { font-size: 17px; color: #6E6E73; line-height: 1.7; margin-bottom: 36px; }

.hero-actions { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; margin-bottom: 56px; }
.btn-primary {
  display: inline-flex; align-items: center; justify-content: center;
  padding: 14px 36px; border-radius: 9999px;
  background: linear-gradient(135deg, #5E5CE6, #7D7AFF);
  color: #FFFFFF; font-size: 16px; font-weight: 600;
  border: none; cursor: pointer;
  box-shadow: 0 4px 20px rgba(94, 92, 230, 0.25);
  transition: all 0.3s ease;
}
.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 8px 28px rgba(94, 92, 230, 0.35); }
.btn-secondary {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 14px 36px; border-radius: 9999px;
  background: #FFFFFF; color: #1D1D1F; font-size: 16px; font-weight: 600;
  border: 1px solid #D2D2D7; text-decoration: none;
  transition: all 0.3s ease;
}
.btn-secondary:hover { background: #F5F5F7; border-color: #C7C7CC; }
.btn-outline {
  display: inline-flex; align-items: center; padding: 12px 28px;
  border-radius: 9999px; background: #FFFFFF; color: #1D1D1F;
  font-size: 14px; font-weight: 600; border: 1px solid #D2D2D7;
  text-decoration: none; transition: all 0.2s;
}
.btn-outline:hover { background: #F5F5F7; border-color: #C7C7CC; }
.btn-ghost {
  display: inline-flex; align-items: center; padding: 12px 28px;
  border-radius: 9999px; background: transparent; color: #6E6E73;
  font-size: 14px; font-weight: 500; border: 1px solid #E5E5EA;
  text-decoration: none; transition: all 0.2s;
}
.btn-ghost:hover { color: #1D1D1F; border-color: #D2D2D7; }

.hero-tags { display: flex; gap: 24px; justify-content: center; flex-wrap: wrap; }
.tag {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 18px; background: #FFFFFF;
  border: 1px solid #F0F0F3; border-radius: 14px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  transition: all 0.3s ease;
}
.tag:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.06); }
.tag-icon { font-size: 18px; }
.tag-label { font-size: 13px; color: #86868B; font-weight: 500; }

/* ── 分区 ── */
.section-gray  { background: #F5F5F7; padding: 80px 20px; }
.section-white { background: #FFFFFF; padding: 80px 20px; }
.container { max-width: 1100px; margin: 0 auto; }
.section-head { text-align: center; margin-bottom: 52px; }
.section-eyebrow {
  font-size: 12px; font-weight: 700; text-transform: uppercase;
  letter-spacing: 2.5px; color: #5E5CE6; margin-bottom: 12px;
}
.section-title {
  font-size: clamp(26px, 4vw, 34px); font-weight: 800;
  color: #1D1D1F; margin-bottom: 12px; letter-spacing: -0.4px;
}
.section-desc { font-size: 16px; color: #86868B; }

/* ── 登录卡片 ── */
.login-card {
  max-width: 420px; margin: 0 auto;
  background: #FFFFFF; border: 1px solid #F0F0F3;
  border-radius: 24px; padding: 44px 36px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.login-card-header { text-align: center; margin-bottom: 28px; }
.login-card-icon  { font-size: 40px; margin-bottom: 8px; display: block; }
.login-card-title { font-size: 24px; font-weight: 800; color: #1D1D1F; }

.login-tabs {
  display: flex; gap: 0; margin-bottom: 24px;
  background: #F5F5F7; border-radius: 12px;
  border: 1px solid #E5E5EA; padding: 4px;
}
.login-tab {
  flex: 1; padding: 10px 8px; text-align: center;
  font-size: 14px; font-weight: 500; border-radius: 9px;
  cursor: pointer; border: none; background: transparent;
  color: #86868B; transition: all 0.25s ease;
}
.login-tab--active {
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
  color: #FFFFFF; box-shadow: 0 2px 8px rgba(94, 92, 230, 0.2);
}

.login-form { display: flex; flex-direction: column; gap: 16px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-group label { font-size: 13px; font-weight: 500; color: #6E6E73; }
.form-input {
  width: 100%; padding: 11px 14px;
  background: #F5F5F7; border: 1px solid #E5E5EA;
  border-radius: 10px; color: #1D1D1F; font-size: 14px;
  outline: none; transition: all 0.2s;
}
.form-input:focus {
  border-color: #5E5CE6; background: #FFFFFF;
  box-shadow: 0 0 0 3px rgba(94, 92, 230, 0.12);
}
.form-input::placeholder { color: #AEAEB2; }

.login-submit {
  width: 100%; padding: 12px; border-radius: 10px;
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
  color: #FFFFFF; font-size: 15px; font-weight: 600;
  border: none; cursor: pointer;
  box-shadow: 0 4px 14px rgba(94, 92, 230, 0.2);
  transition: all 0.3s ease;
}
.login-submit:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 20px rgba(94, 92, 230, 0.3); }
.login-submit:disabled { opacity: 0.6; cursor: not-allowed; }
.login-submit--amber { background: linear-gradient(135deg, #F59E0B, #FBBF24); box-shadow: 0 4px 14px rgba(245, 158, 11, 0.2); }
.login-submit--green { background: linear-gradient(135deg, #10B981, #34D399); box-shadow: 0 4px 14px rgba(16, 185, 129, 0.2); }

.login-social { margin-top: 24px; text-align: center; }
.login-social-label { display: flex; align-items: center; gap: 8px; color: #AEAEB2; font-size: 12px; margin-bottom: 14px; }
.login-social-line { flex: 1; height: 1px; background: #E5E5EA; }
.login-social-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.login-social-btn {
  display: flex; flex-direction: column; align-items: center; gap: 5px;
  padding: 10px 4px; background: #F5F5F7; border: 1px solid #E5E5EA;
  border-radius: 10px; color: #6E6E73; font-size: 11px;
  cursor: pointer; transition: all 0.2s;
}
.login-social-btn:hover { background: #FFFFFF; border-color: #D2D2D7; color: #1D1D1F; }
.login-social-icon {
  display: flex; align-items: center; justify-content: center;
  width: 28px; height: 28px; border-radius: 8px;
  font-size: 14px; color: #FFFFFF;
}

/* ── Bento 卡片网格 ── */
.bento { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; }
.bento-card {
  position: relative; background: #FFFFFF;
  border: 1px solid #F0F0F3; border-radius: 20px; padding: 28px;
  transition: all 0.35s ease; cursor: default;
  box-shadow: 0 1px 3px rgba(0,0,0,0.03);
}
.bento-card:hover { transform: translateY(-4px); box-shadow: 0 12px 36px rgba(0,0,0,0.08); border-color: #E5E5EA; }
.bento-icon  { display: block; font-size: 32px; margin-bottom: 14px; }
.bento-title { font-size: 15px; font-weight: 700; color: #1D1D1F; margin-bottom: 6px; }
.bento-desc  { font-size: 13px; color: #86868B; line-height: 1.55; }
.bento-wide { grid-column: span 2; }
.bento-lg   { grid-column: span 2; }
.bento-tall { grid-row: span 2; }

/* ── 宠物进化 ── */
.evo-track { display: flex; align-items: center; justify-content: center; gap: 2px; flex-wrap: wrap; padding: 20px 0; }
.evo-arrow { color: #C7C7CC; font-size: 14px; font-weight: 300; user-select: none; }
.evo-node {
  display: flex; flex-direction: column; align-items: center; gap: 6px;
  padding: 10px 14px; border-radius: 16px;
  transition: all 0.25s ease; cursor: default;
}
.evo-node:hover { background: #F5F5F7; transform: scale(1.1); }
.evo-emoji { font-size: 30px; }
.evo-name  { font-size: 12px; color: #86868B; font-weight: 500; }

/* ── CTA ── */
.cta {
  max-width: 600px; margin: 0 auto; text-align: center;
  background: #FFFFFF; border: 1px solid #F0F0F3;
  border-radius: 24px; padding: 52px 40px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.cta-icon  { font-size: 44px; margin-bottom: 18px; }
.cta-title { font-size: 28px; font-weight: 800; color: #1D1D1F; margin-bottom: 16px; }
.cta-desc  { font-size: 15px; color: #86868B; line-height: 1.7; margin-bottom: 32px; }
.cta-btns  { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; }

/* ── Footer ── */
.footer {
  padding: 32px 20px; text-align: center;
  background: #F5F5F7; border-top: 1px solid #E5E5EA;
}
.footer-line { color: #86868B; font-size: 13px; margin-bottom: 6px; }
.footer-link { color: #515154; text-decoration: none; transition: color 0.2s; }
.footer-link:hover { color: #1D1D1F; }
.footer-made { color: #AEAEB2; font-size: 12px; }

/* ── 滚动揭示 ── */
.reveal {
  opacity: 0; transform: translateY(36px);
  transition: opacity 0.7s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}
.revealed { opacity: 1; transform: translateY(0); }

/* Bento 交错延迟 */
.bento-card:nth-child(1)  { transition-delay: 0.00s, 0.00s; }
.bento-card:nth-child(2)  { transition-delay: 0.05s, 0.05s; }
.bento-card:nth-child(3)  { transition-delay: 0.10s, 0.10s; }
.bento-card:nth-child(4)  { transition-delay: 0.15s, 0.15s; }
.bento-card:nth-child(5)  { transition-delay: 0.20s, 0.20s; }
.bento-card:nth-child(6)  { transition-delay: 0.05s, 0.05s; }
.bento-card:nth-child(7)  { transition-delay: 0.10s, 0.10s; }
.bento-card:nth-child(8)  { transition-delay: 0.15s, 0.15s; }
.bento-card:nth-child(9)  { transition-delay: 0.20s, 0.20s; }
.bento-card:nth-child(10) { transition-delay: 0.05s, 0.05s; }
.bento-card:nth-child(11) { transition-delay: 0.10s, 0.10s; }
.bento-card:nth-child(12) { transition-delay: 0.15s, 0.15s; }

/* ── 响应式 ── */
@media (max-width: 1024px) {
  .bento { grid-template-columns: repeat(3, 1fr); }
  .bento-wide, .bento-lg { grid-column: span 2; }
}
@media (max-width: 768px) {
  .nav-inner { padding: 0 16px; height: 48px; }
  .nav-links a { display: none; }
  .hero { padding: 100px 16px 60px; }
  .hero-title { font-size: clamp(28px, 7vw, 42px); }
  .hero-actions { flex-direction: column; align-items: center; }
  .hero-tags { gap: 10px; }
  .bento { grid-template-columns: 1fr; }
  .bento-wide, .bento-lg, .bento-tall { grid-column: span 1; grid-row: span 1; }
  .section-gray, .section-white { padding: 48px 16px; }
  .section-head { margin-bottom: 32px; }
  .login-card { padding: 32px 24px; }
  .evo-track { gap: 1px; }
  .evo-node  { padding: 6px 8px; }
  .evo-emoji { font-size: 22px; }
  .evo-name  { font-size: 10px; }
  .cta { padding: 32px 20px; }
}
@media (max-width: 480px) {
  .evo-arrow { font-size: 10px; }
}
</style>
