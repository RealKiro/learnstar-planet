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

function formatCode(e: Event) {
  const input = e.target as HTMLInputElement
  displayCode.value = input.value.toUpperCase().replace(/[^0-9A-Z-]/g, '')
  codeError.value = ''
}

function goToDisplay() {
  const code = displayCode.value.trim()
  if (code.length < 6) { codeError.value = '请输入完整的班级码'; return }
  router.push({ name: 'display-login', query: { code } })
}

function goLogin(role: string) {
  router.push({ name: 'login', query: { role } })
}
</script>

<template>
  <div class="home">
    <div class="home-bg">
      <div class="home-bg-orb home-bg-orb--1"></div>
      <div class="home-bg-orb home-bg-orb--2"></div>
      <div class="home-bg-orb home-bg-orb--3"></div>
    </div>

    <!-- 顶部导航：唯一登录入口，不再重复 -->
    <nav class="home-nav">
      <div class="nav-brand">
        <span class="nav-logo">🌌</span>
        <span class="nav-name">学趣星球</span>
      </div>
      <div class="nav-actions">
        <button class="nav-link" @click="goLogin('teacher')">教师登录</button>
        <button class="nav-link" @click="goLogin('admin')">管理员</button>
        <button class="nav-link" @click="goLogin('parent')">家长</button>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="nav-link nav-link--icon" title="GitHub">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
      </div>
    </nav>

    <!-- 一屏显示，不滚动 -->
    <main class="home-main">
      <div class="hero">
        <div class="hero-badge">
          <span class="hero-badge-dot"></span>
          MIT 开源 · 完全免费 · 自托管
        </div>
        <h1 class="hero-title">
          让每个孩子的努力<br>
          <span class="hero-title-grad">都被看见</span>
        </h1>
      </div>

      <!-- 唯一核心操作：大屏 -->
      <div class="display-card">
        <div class="display-card-glow"></div>
        <div class="display-card-bg"></div>
        <div class="display-card-content">
          <div class="display-icon">🖥️</div>
          <div class="display-text">
            <h2 class="display-title">班级大屏</h2>
            <p class="display-desc">教室触摸屏展示 · 实时积分 · 宠物矩阵</p>
          </div>
          <form class="display-form" @submit.prevent="goToDisplay">
            <div class="display-input-row">
              <input v-model="displayCode" type="text" class="display-input"
                placeholder="输入班级码" maxlength="12" autocomplete="off"
                @input="formatCode">
              <button type="submit" class="display-btn" :disabled="displayCode.length < 6">
                进入
              </button>
            </div>
            <p v-if="codeError" class="display-error">{{ codeError }}</p>
            <p class="display-hint">向班主任获取班级码</p>
          </form>
        </div>
      </div>

      <div class="footer-meta">
        <span>积分激励</span><span class="dot">·</span>
        <span>宠物进化</span><span class="dot">·</span>
        <span>排行榜</span><span class="dot">·</span>
        <span>实时广播</span><span class="dot">·</span>
        <span>AI 助教</span><span class="dot">·</span>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="footer-link">GitHub</a>
      </div>
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
}

/* 动态渐变背景 */
.home-bg { position: fixed; inset: 0; overflow: hidden; pointer-events: none; z-index: 0; }
.home-bg-orb {
  position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.18;
}
.home-bg-orb--1 { width: 600px; height: 600px; background: #C7D2FE; top: -200px; right: -150px; animation: orb1 20s ease-in-out infinite; }
.home-bg-orb--2 { width: 450px; height: 450px; background: #A7F3D0; bottom: -150px; left: -100px; animation: orb2 25s ease-in-out infinite; }
.home-bg-orb--3 { width: 300px; height: 300px; background: #FDE68A; top: 40%; left: 50%; animation: orb3 18s ease-in-out infinite; }
@keyframes orb1 { 0%,100%{transform:translate(0,0)} 33%{transform:translate(50px,-40px)} 66%{transform:translate(-20px,30px)} }
@keyframes orb2 { 0%,100%{transform:translate(0,0)} 50%{transform:translate(-60px,-50px)} }
@keyframes orb3 { 0%,100%{transform:translate(0,0) scale(1)} 50%{transform:translate(40px,30px) scale(1.1)} }

/* 导航 */
.home-nav {
  position: relative; z-index: 10;
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 32px; max-width: 680px; margin: 0 auto; width: 100%; box-sizing: border-box;
  flex-shrink: 0;
}
.nav-brand { display: flex; align-items: center; gap: 10px; }
.nav-logo { font-size: 28px; }
.nav-name { font-size: 17px; font-weight: 700; color: #1D1D1F; }
.nav-actions { display: flex; align-items: center; gap: 2px; }
.nav-link {
  padding: 8px 14px; border: none; border-radius: 9999px; background: transparent;
  color: #6E6E73; font-size: 13px; font-weight: 500; cursor: pointer;
  transition: all .2s; text-decoration: none; font-family: inherit;
}
.nav-link:hover { background: rgba(0,0,0,.04); color: #1D1D1F; }
.nav-link--icon { display: flex; align-items: center; padding: 8px; color: #AEAEB2; }
.nav-link--icon:hover { color: #1D1D1F; }

/* 主内容 - 垂直居中，一屏显示 */
.home-main {
  position: relative; z-index: 1; flex: 1;
  max-width: 480px; width: 100%; margin: 0 auto;
  padding: 0 32px 20px; box-sizing: border-box;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  min-height: 0;
}

/* Hero - 紧凑 */
.hero { text-align: center; margin-bottom: 32px; }
.hero-badge {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 5px 14px; background: rgba(255,255,255,.8);
  backdrop-filter: blur(20px); border: 1px solid rgba(0,0,0,.06);
  border-radius: 9999px; font-size: 12px; color: #6E6E73; margin-bottom: 20px;
}
.hero-badge-dot { width: 6px; height: 6px; background: #34C759; border-radius: 50%; }
.hero-title { font-size: 40px; font-weight: 800; line-height: 1.15; letter-spacing: -.03em; margin: 0; }
.hero-title-grad { background: linear-gradient(135deg,#5E5CE6,#FF375F,#FF9F0A); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }

/* 大屏卡片 - 深色，占据主体 */
.display-card {
  width: 100%; position: relative; overflow: hidden;
  border-radius: 20px; cursor: default;
  background: linear-gradient(145deg,#1a1a2e,#16213e);
}
.display-card::before {
  content: ''; position: absolute; inset: 0;
  background: radial-gradient(ellipse at 30% 20%,rgba(120,80,255,.25),transparent 60%),
              radial-gradient(ellipse at 70% 80%,rgba(30,140,220,.15),transparent 60%);
  pointer-events: none;
}
.display-card-glow {
  position: absolute; width: 200px; height: 200px;
  background: radial-gradient(circle,rgba(120,80,255,.1),transparent);
  top: -60px; right: -60px; border-radius: 50%; pointer-events: none;
}
.display-card-content { padding: 28px; position: relative; z-index: 1; }
.display-icon { font-size: 36px; margin-bottom: 8px; }
.display-text { margin-bottom: 20px; }
.display-title { font-size: 20px; font-weight: 700; color: #fff; margin: 0 0 4px; }
.display-desc { font-size: 13px; color: rgba(255,255,255,.55); margin: 0; }

/* 输入 */
.display-input-row { display: flex; gap: 8px; }
.display-input {
  flex: 1; padding: 12px 16px; border-radius: 12px;
  border: 1.5px solid rgba(255,255,255,.12);
  background: rgba(255,255,255,.06); color: #fff;
  font-size: 18px; font-weight: 600; letter-spacing: .12em;
  outline: none; transition: all .2s;
  font-family: 'SF Mono', 'SF Pro Text', 'Noto Sans SC', monospace;
}
.display-input::placeholder { color: rgba(255,255,255,.25); font-weight: 400; letter-spacing: 0; font-size: 15px; }
.display-input:focus { border-color: rgba(120,80,255,.5); box-shadow: 0 0 0 3px rgba(120,80,255,.12); }
.display-btn {
  padding: 12px 24px; border-radius: 12px; border: none;
  background: linear-gradient(135deg,#7C3AED,#6D28D9); color: #fff;
  font-size: 15px; font-weight: 600; cursor: pointer;
  transition: all .2s; white-space: nowrap; font-family: inherit;
}
.display-btn:hover:not(:disabled) { background: linear-gradient(135deg,#8B5CF6,#7C3AED); transform: scale(1.02); }
.display-btn:disabled { opacity: .4; cursor: not-allowed; }
.display-error { color: #FCA5A5; font-size: 12px; margin: 8px 0 0; }
.display-hint { font-size: 12px; color: rgba(255,255,255,.3); margin: 10px 0 0; }

/* 底部 */
.footer-meta {
  margin-top: 28px; display: flex; align-items: center; gap: 6px;
  font-size: 12px; color: #8E8E93; flex-wrap: wrap; justify-content: center;
}
.dot { color: #C7C7CC; }
.footer-link { color: #8E8E93; text-decoration: none; transition: color .2s; }
.footer-link:hover { color: #1D1D1F; }

@media (max-width: 768px) {
  .home-nav { padding: 14px 20px; }
  .home-main { padding: 0 20px 16px; max-width: 100%; }
  .hero-title { font-size: 32px; }
  .hero { margin-bottom: 24px; }
  .display-card-content { padding: 22px; }
  .display-input { font-size: 16px; }
}
</style>
