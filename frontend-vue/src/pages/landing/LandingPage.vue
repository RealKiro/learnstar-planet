<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const displayCode = ref('')
const codeError = ref('')

onMounted(() => {
  if (authStore.isLoggedIn) {
    if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
    else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
    else if (authStore.isParent) router.replace({ name: 'parent-home' })
  }
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  if (code) {
    displayCode.value = code.toUpperCase()
    goToDisplay()
  }
})

function onCodeInput(e: Event) {
  const input = e.target as HTMLInputElement
  displayCode.value = input.value.toUpperCase().replace(/[^0-9A-Z-]/g, '')
  codeError.value = ''
}

function goToDisplay() {
  if (!displayCode.value.trim()) { codeError.value = '请输入班级码'; return }
  router.push({ name: 'display-login', query: { code: displayCode.value } })
}

function goLogin(role: string) {
  router.push({ name: 'login', query: { role } })
}
</script>

<template>
  <div class="page">
    <div class="bg">
      <div class="bg-orb bg-orb--1"></div>
      <div class="bg-orb bg-orb--2"></div>
      <div class="bg-orb bg-orb--3"></div>
    </div>

    <!-- 顶部 -->
    <header class="header">
      <div class="header-brand">
        <span class="header-logo">🌌</span>
        <span class="header-name">学趣星球</span>
      </div>
      <div class="header-links">
        <button class="header-link" @click="goLogin('teacher')">教师登录</button>
        <button class="header-link header-link--dim" @click="goLogin('admin')">管理员</button>
        <button class="header-link header-link--dim" @click="goLogin('parent')">家长</button>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="header-link header-link--icon" title="GitHub">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
      </div>
    </header>

    <!-- 主体 -->
    <main class="main">
      <div class="card">
        <div class="card-glow"></div>

        <!-- 顶部徽章 -->
        <div class="card-badge">
          <span class="card-badge-dot"></span>
          MIT 开源 · 完全免费 · 自托管
        </div>

        <!-- 标题 -->
        <h1 class="card-heading">
          让每个孩子的努力
          <br>
          <span class="card-heading-grad">都被看见</span>
        </h1>

        <!-- 分割 -->
        <div class="card-divider"></div>

        <!-- 大屏入口 -->
        <div class="card-body">
          <div class="card-body-icon">🖥️</div>
          <div class="card-body-title">班级大屏</div>
          <div class="card-body-desc">教室触摸屏展示 · 实时积分变化 · 宠物矩阵</div>

          <form class="card-form" @submit.prevent="goToDisplay">
            <div class="card-form-row">
              <input v-model="displayCode" type="text" class="card-input"
                placeholder="输入班级码" maxlength="12" autocomplete="off"
                @input="onCodeInput">
              <button type="submit" class="card-btn" :disabled="displayCode.length < 6">
                进入
              </button>
            </div>
            <p v-if="codeError" class="card-error">{{ codeError }}</p>
          </form>

          <div class="card-footnote">向班主任获取班级码</div>
        </div>
      </div>

      <!-- 底部功能标签 -->
      <div class="tags">
        <span>积分激励</span>
        <span>宠物进化</span>
        <span>排行榜</span>
        <span>实时广播</span>
        <span>AI 助教</span>
      </div>
    </main>
  </div>
</template>

<style scoped>
/* ============================================================
   学趣星球 — 首页
   ============================================================ */

.page {
  height: 100vh; overflow: hidden;
  background: #F5F5F7;
  color: #1D1D1F;
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display',
    'SF Pro Text', 'Helvetica Neue', 'Noto Sans SC', sans-serif;
  -webkit-font-smoothing: antialiased;
  display: flex; flex-direction: column;
  position: relative;
}

/* ===== 背景 ===== */
.bg { position: fixed; inset: 0; overflow: hidden; pointer-events: none; z-index: 0; }
.bg-orb {
  position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.15;
}
.bg-orb--1 { width: 500px; height: 500px; background: #C7D2FE; top: -150px; right: -100px; animation: d1 18s ease-in-out infinite; }
.bg-orb--2 { width: 400px; height: 400px; background: #A7F3D0; bottom: -120px; left: -80px; animation: d2 22s ease-in-out infinite; }
.bg-orb--3 { width: 250px; height: 250px; background: #FDE68A; top: 50%; left: 60%; animation: d3 15s ease-in-out infinite; }
@keyframes d1 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(40px,-30px)} }
@keyframes d2 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(-40px,40px)} }
@keyframes d3 { 0%,100%{transform:translate(0,0) scale(1)} 50%{transform:translate(30px,20px) scale(1.15)} }

/* ===== 顶栏 ===== */
.header {
  position: relative; z-index: 10;
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 32px; max-width: 520px; margin: 0 auto; width: 100%;
  box-sizing: border-box; flex-shrink: 0;
}
.header-brand { display: flex; align-items: center; gap: 8px; }
.header-logo { font-size: 24px; line-height: 1; }
.header-name { font-size: 16px; font-weight: 700; letter-spacing: -.02em; }
.header-links { display: flex; align-items: center; gap: 2px; }
.header-link {
  padding: 6px 12px; border: none; border-radius: 6px; background: transparent;
  color: #6E6E73; font-size: 13px; font-weight: 500; cursor: pointer;
  transition: all .2s; text-decoration: none; font-family: inherit;
}
.header-link:hover { background: rgba(0,0,0,.04); color: #1D1D1F; }
.header-link--dim { color: #AEAEB2; font-size: 12px; }
.header-link--dim:hover { color: #6E6E73; }
.header-link--icon { display: flex; align-items: center; padding: 6px; color: #AEAEB2; }
.header-link--icon:hover { color: #6E6E73; }

/* ===== 主体 ===== */
.main {
  position: relative; z-index: 1; flex: 1;
  max-width: 420px; width: 100%; margin: 0 auto;
  padding: 0 24px 24px; box-sizing: border-box;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  min-height: 0;
}

/* ===== 卡片 ===== */
.card {
  width: 100%; position: relative;
  background: #FFFFFF;
  border-radius: 24px;
  box-shadow: 0 2px 20px rgba(0,0,0,.04), 0 1px 3px rgba(0,0,0,.02);
  padding: 36px 32px 32px;
  overflow: hidden;
}

.card-glow {
  position: absolute; width: 240px; height: 240px;
  background: radial-gradient(circle at 30% 20%, rgba(120,80,255,.06), transparent 70%);
  top: -60px; right: -60px; border-radius: 50%; pointer-events: none;
}

/* 徽章 */
.card-badge {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 4px 12px;
  background: #F5F5F7; border-radius: 9999px;
  font-size: 11px; color: #8E8E93; margin-bottom: 20px;
}
.card-badge-dot { width: 5px; height: 5px; background: #34C759; border-radius: 50%; flex-shrink: 0; }

/* 标题 */
.card-heading {
  font-size: 34px; font-weight: 800; line-height: 1.15;
  letter-spacing: -.025em; margin: 0 0 24px;
}
.card-heading-grad {
  background: linear-gradient(135deg,#5E5CE6,#FF375F,#FF9F0A);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 分割线 */
.card-divider {
  height: 1px; background: #F0F0F0; margin-bottom: 24px;
}

/* 大屏入口 */
.card-body { text-align: center; }
.card-body-icon { font-size: 32px; margin-bottom: 8px; }
.card-body-title { font-size: 18px; font-weight: 700; margin-bottom: 4px; }
.card-body-desc { font-size: 13px; color: #8E8E93; margin-bottom: 20px; }

/* 输入 */
.card-form-row { display: flex; gap: 8px; }
.card-input {
  flex: 1; padding: 12px 16px; border-radius: 12px;
  border: 1.5px solid #E5E5EA;
  background: #F5F5F7; color: #1D1D1F;
  font-size: 20px; font-weight: 600; letter-spacing: .1em;
  outline: none; transition: all .2s;
  font-family: 'SF Mono', 'SF Pro Text', 'Noto Sans SC', monospace;
  text-align: center;
}
.card-input::placeholder { color: #AEAEB2; font-weight: 400; letter-spacing: 0; font-size: 16px; }
.card-input:focus { border-color: #5E5CE6; background: #fff; box-shadow: 0 0 0 3px rgba(94,92,230,.1); }
.card-btn {
  padding: 12px 24px; border-radius: 12px; border: none;
  background: linear-gradient(135deg,#5E5CE6,#818CF8); color: #fff;
  font-size: 15px; font-weight: 600; cursor: pointer;
  transition: all .2s; white-space: nowrap; font-family: inherit;
}
.card-btn:hover:not(:disabled) { transform: scale(1.02); box-shadow: 0 4px 14px rgba(94,92,230,.25); }
.card-btn:disabled { opacity: .35; cursor: not-allowed; }
.card-error { color: #EF4444; font-size: 12px; margin: 8px 0 0; text-align: left; }
.card-footnote { font-size: 12px; color: #AEAEB2; margin-top: 14px; }

/* ===== 底部标签 ===== */
.tags {
  margin-top: 24px;
  display: flex; gap: 10px; flex-wrap: wrap; justify-content: center;
  font-size: 11px; color: #AEAEB2;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .header { padding: 14px 20px; }
  .main { padding: 0 16px 20px; }
  .card { padding: 28px 24px 28px; border-radius: 20px; }
  .card-heading { font-size: 28px; }
  .card-input { font-size: 18px; }
}
</style>
