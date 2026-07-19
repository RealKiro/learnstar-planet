<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { apiPost } from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()

// 班级码
const displayCode = ref('')
const codeError = ref('')
const loading = ref(false)

// 幻灯片
const slides = [
  { icon: '🌌', title: '让每个孩子的努力', highlight: '都被看见', desc: '积分激励 · 宠物养成 · AI 助教\n开源班级管理系统，数据完全自主掌控' },
  { icon: '⚡', title: '覆盖班级管理', highlight: '全场景', desc: '积分规则 · 宠物进化 · 排行榜 · 通知\n考勤 · 作业 · 答题 · 成绩 · 商城 · 广播' },
  { icon: '🌟', title: '11 阶宠物进化', highlight: '驱动成长', desc: '星尘 → 月芽 → 灵苗 → 青藤 → 慧树\n→ 蝶灵 → 鹰慧 → 狮睿 → 灵角 → 星耀 → 银河' },
]
const currentSlide = ref(0)
let slideTimer: ReturnType<typeof setInterval>

onMounted(() => {
  if (authStore.isLoggedIn) {
    if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
    else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
    else if (authStore.isParent) router.replace({ name: 'parent-home' })
  }
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 4000)
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  if (code) { displayCode.value = code.toUpperCase(); goToDisplay() }
})
onUnmounted(() => clearInterval(slideTimer))

function onCodeInput(e: Event) {
  const input = e.target as HTMLInputElement
  displayCode.value = input.value.toUpperCase().replace(/[^0-9A-Z-]/g, '')
  codeError.value = ''
}

async function goToDisplay() {
  const code = displayCode.value.trim()
  if (!code) { codeError.value = '请输入班级码'; return }
  loading.value = true; codeError.value = ''
  try {
    const res = await apiPost<{ data: { token: string; class_id: number; class_name: string; student_count: number } }>(
      '/api/v1/auth/class/login', { class_code: code }
    )
    sessionStorage.setItem('class_token', res.data.token)
    sessionStorage.setItem('class_info', JSON.stringify({
      id: res.data.class_id,
      name: res.data.class_name,
      student_count: res.data.student_count,
    }))
    router.replace({ name: 'classroom-overview' })
  } catch { codeError.value = '班级码无效，请检查后重试'; loading.value = false }
}

function goLogin(role: string) { router.push({ name: 'login', query: { role } }) }
</script>

<template>
  <div class="home">
    <!-- 背景光晕 -->
    <div class="bg">
      <div class="bg-orb bg-orb--1"></div>
      <div class="bg-orb bg-orb--2"></div>
    </div>

    <header class="topbar">
      <div class="topbar-brand">🌌 学趣星球</div>
      <div class="topbar-links">
        <button class="topbar-btn" @click="goLogin('teacher')">教师登录</button>
        <button class="topbar-btn topbar-btn--dim" @click="goLogin('admin')">管理员</button>
        <button class="topbar-btn topbar-btn--dim" @click="goLogin('parent')">家长</button>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="topbar-link" title="GitHub">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
      </div>
    </header>

    <main class="main">
      <!-- 左侧：幻灯片介绍 -->
      <section class="left">
        <div class="slide-stage">
          <Transition name="slide" mode="out-in">
            <div :key="currentSlide" class="slide-card">
              <div class="slide-icon">{{ slides[currentSlide].icon }}</div>
              <h1 class="slide-title">
                {{ slides[currentSlide].title }}<br>
                <span class="slide-highlight">{{ slides[currentSlide].highlight }}</span>
              </h1>
              <p class="slide-desc">{{ slides[currentSlide].desc }}</p>
            </div>
          </Transition>
          <div class="dots">
            <button v-for="(_, i) in slides" :key="i" :class="['dot', { active: currentSlide === i }]" @click="currentSlide = i"></button>
          </div>
        </div>
      </section>

      <!-- 右侧：登录选择 -->
      <section class="right">
        <div class="panel">
          <div class="panel-icon">🌌</div>
          <h2 class="panel-title">进入学趣星球</h2>
          <p class="panel-desc">选择登录方式</p>

          <div class="panel-buttons">
            <button class="panel-btn panel-btn--primary" @click="goLogin('teacher')">
              👨‍🏫 教师登录
            </button>
            <button class="panel-btn panel-btn--outline" @click="goLogin('parent')">
              👨‍👩‍👧‍👦 家长登录
            </button>
            <button class="panel-btn panel-btn--ghost" @click="router.push({ name: 'login', query: { mode: 'code' } })">
              🔑 班级码登录
            </button>
          </div>

          <div class="panel-features">
            <span>⭐ 积分激励</span>
            <span>🌟 宠物进化</span>
            <span>🏆 排行榜</span>
            <span>📢 实时广播</span>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.home {
  height: 100vh; overflow: hidden;
  background: #FBFBFD; color: #1D1D1F;
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text',
    'Helvetica Neue', 'Noto Sans SC', sans-serif;
  -webkit-font-smoothing: antialiased;
  display: flex; flex-direction: column;
  position: relative;
}

/* ===== 背景 ===== */
.bg { position: fixed; inset: 0; overflow: hidden; pointer-events: none; z-index: 0; }
.bg-orb {
  position: absolute; border-radius: 50%; filter: blur(120px); opacity: 0.12;
}
.bg-orb--1 { width: 600px; height: 600px; background: #C7D2FE; top: -200px; right: -100px; animation: f1 20s ease infinite; }
.bg-orb--2 { width: 500px; height: 500px; background: #A7F3D0; bottom: -150px; left: -80px; animation: f2 25s ease infinite; }
@keyframes f1 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(50px,-40px)} }
@keyframes f2 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(-60px,50px)} }

/* ===== 顶栏 ===== */
.topbar {
  position: relative; z-index: 10;
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 40px; flex-shrink: 0;
}
.topbar-brand { font-size: 17px; font-weight: 700; letter-spacing: -.02em; }
.topbar-links { display: flex; align-items: center; gap: 4px; }
.topbar-btn {
  padding: 6px 14px; border: none; border-radius: 6px; background: transparent;
  color: #6E6E73; font-size: 13px; font-weight: 500; cursor: pointer;
  transition: all .2s; font-family: inherit;
}
.topbar-btn:hover { background: rgba(0,0,0,.04); color: #1D1D1F; }
.topbar-btn--dim { color: #AEAEB2; font-size: 12px; }
.topbar-btn--dim:hover { color: #6E6E73; }
.topbar-link { display: flex; align-items: center; padding: 6px; color: #AEAEB2; cursor: pointer; }
.topbar-link:hover { color: #6E6E73; }

/* ===== 左右主体 ===== */
.main {
  position: relative; z-index: 1; isolation: isolate;
  flex: 1; min-height: 0;
  display: flex; align-items: center; justify-content: center;
  padding: 0 40px 40px; gap: 60px;
}

/* ===== 左侧 ===== */
.left {
  flex: 1; max-width: 520px;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
}

.slide-stage { width: 100%; text-align: center; }
.slide-card { padding: 20px 0; }
.slide-icon { font-size: 64px; margin-bottom: 20px; }
.slide-title {
  font-size: 44px; font-weight: 800; line-height: 1.15;
  letter-spacing: -.03em; margin: 0 0 16px;
}
.slide-highlight {
  background: linear-gradient(135deg,#5E5CE6,#FF375F,#FF9F0A);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  background-clip: text;
}
.slide-desc { font-size: 16px; color: #86868B; line-height: 1.7; white-space: pre-line; }

.slide-enter-active, .slide-leave-active { transition: all .4s cubic-bezier(.4,0,.2,1); }
.slide-enter-from { opacity: 0; transform: translateX(30px); }
.slide-leave-to { opacity: 0; transform: translateX(-30px); }

.dots { display: flex; gap: 8px; justify-content: center; margin-top: 32px; }
.dot {
  width: 8px; height: 8px; border-radius: 50%; border: none; background: #D2D2D7;
  cursor: pointer; padding: 0; transition: all .3s;
}
.dot.active { width: 28px; border-radius: 4px; background: linear-gradient(135deg,#5E5CE6,#818CF8); }

/* ===== 右侧面板 ===== */
.right { flex: 0 0 380px; }

.panel {
  background: #FFFFFF;
  border-radius: 24px;
  box-shadow: 0 4px 30px rgba(0,0,0,.05), 0 1px 3px rgba(0,0,0,.03);
  padding: 36px 32px;
  text-align: center;
}
.panel-icon { font-size: 36px; margin-bottom: 10px; }
.panel-title { font-size: 20px; font-weight: 700; margin: 0 0 4px; }
.panel-desc { font-size: 13px; color: #8E8E93; margin: 0 0 24px; }

.panel-form { display: flex; flex-direction: column; gap: 10px; }
.panel-input {
  width: 100%; padding: 14px 16px; border-radius: 12px;
  border: 1.5px solid #E5E5EA;
  background: #F5F5F7; color: #1D1D1F;
  font-size: 24px; font-weight: 600; letter-spacing: .12em;
  outline: none; transition: all .2s; text-align: center;
  font-family: 'SF Mono', 'SF Pro Text', 'Noto Sans SC', monospace;
  box-sizing: border-box;
}
.panel-input::placeholder { color: #AEAEB2; font-weight: 400; font-size: 16px; letter-spacing: 0; }
.panel-input:focus { border-color: #5E5CE6; background: #fff; box-shadow: 0 0 0 3px rgba(94,92,230,.1); }
.panel-btn {
  width: 100%; padding: 14px; border-radius: 12px; border: none;
  background: linear-gradient(135deg,#5E5CE6,#818CF8); color: #fff;
  font-size: 16px; font-weight: 600; cursor: pointer;
  transition: all .2s; font-family: inherit;
}
.panel-btn:hover:not(:disabled) { box-shadow: 0 4px 14px rgba(94,92,230,.25); }
.panel-btn:disabled { opacity: .35; cursor: not-allowed; }
.panel-error { color: #EF4444; font-size: 13px; margin: 0; }
.panel-footnote { font-size: 12px; color: #AEAEB2; margin-top: 14px; }

.panel-features {
  margin-top: 24px; padding-top: 20px; border-top: 1px solid #F0F0F0;
  display: flex; gap: 10px; flex-wrap: wrap; justify-content: center;
  font-size: 11px; color: #AEAEB2;
}

/* ===== 响应式 ===== */
@media (max-width: 860px) {
  .main { flex-direction: column; padding: 0 20px 24px; gap: 24px; }
  .left { max-width: 100%; }
  .slide-title { font-size: 32px; }
  .slide-icon { font-size: 48px; margin-bottom: 12px; }
  .right { flex: none; width: 100%; max-width: 400px; }
  .panel { padding: 28px 24px; }
}
</style>
