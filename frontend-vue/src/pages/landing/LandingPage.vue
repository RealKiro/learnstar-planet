<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const displayCode = ref('')
const codeError = ref('')

// 已登录直接跳转
onMounted(() => {
  if (authStore.isLoggedIn) {
    if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
    else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
    else if (authStore.isParent) router.replace({ name: 'parent-home' })
  }
  // URL 带 code 参数？自动跳转大屏
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
    <!-- 动态渐变背景 -->
    <div class="home-bg">
      <div class="home-bg-orb home-bg-orb--1"></div>
      <div class="home-bg-orb home-bg-orb--2"></div>
      <div class="home-bg-orb home-bg-orb--3"></div>
    </div>

    <!-- 顶部导航 -->
    <nav class="home-nav">
      <div class="nav-brand">
        <span class="nav-logo">🌌</span>
        <span class="nav-name">学趣星球</span>
      </div>
      <div class="nav-actions">
        <button class="nav-link" @click="goLogin('teacher')">教师登录</button>
        <button class="nav-link nav-link--muted" @click="goLogin('admin')">管理员</button>
        <button class="nav-link nav-link--muted" @click="goLogin('parent')">家长</button>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="nav-link nav-link--icon" title="GitHub">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
      </div>
    </nav>

    <!-- 主内容 -->
    <main class="home-main">
      <!-- Hero -->
      <section class="hero">
        <div class="hero-badge">
          <span class="hero-badge-dot"></span>
          MIT 开源 · 完全免费 · 自托管
        </div>
        <h1 class="hero-title">
          让每个孩子的努力<br>
          <span class="hero-title-grad">都被看见</span>
        </h1>
        <p class="hero-desc">积分激励 · 宠物养成 · AI 助教 · 开源班级管理系统</p>
      </section>

      <!-- 核心操作区 -->
      <section class="actions">
        <!-- 主：班级大屏 -->
        <div class="action-card action-card--primary">
          <div class="action-card-glow"></div>
          <div class="action-card-content">
            <div class="action-icon">🖥️</div>
            <h2 class="action-title">班级大屏</h2>
            <p class="action-desc">教室触摸屏展示 · 实时积分变化 · 宠物矩阵</p>
            <form class="code-form" @submit.prevent="goToDisplay">
              <div class="code-field">
                <input v-model="displayCode" type="text" class="code-input"
                  placeholder="输入班级码，如 3-1-A7K2" maxlength="12" autocomplete="off"
                  @input="formatCode">
                <button type="submit" class="code-btn" :disabled="displayCode.length < 6">
                  进入 →
                </button>
              </div>
              <p v-if="codeError" class="code-error">{{ codeError }}</p>
            </form>
            <p class="action-hint">向班主任获取班级码，在触摸屏上打开即可使用</p>
          </div>
        </div>

        <!-- 辅助：登录入口 -->
        <div class="actions-secondary">
          <div class="action-card action-card--small" @click="goLogin('teacher')">
            <div class="action-card-content action-card-content--row">
              <div class="action-icon action-icon--sm">👨‍🏫</div>
              <div class="action-text">
                <h3 class="action-title action-title--sm">教师管理</h3>
                <p class="action-desc action-desc--sm">积分、宠物、课堂工具</p>
              </div>
              <span class="action-arrow">→</span>
            </div>
          </div>
          <div class="action-card action-card--small" @click="goLogin('admin')">
            <div class="action-card-content action-card-content--row">
              <div class="action-icon action-icon--sm">🏫</div>
              <div class="action-text">
                <h3 class="action-title action-title--sm">学校管理</h3>
                <p class="action-desc action-desc--sm">班级、账号、数据报表</p>
              </div>
              <span class="action-arrow">→</span>
            </div>
          </div>
          <div class="action-card action-card--small" @click="goLogin('parent')">
            <div class="action-card-content action-card-content--row">
              <div class="action-icon action-icon--sm">👪</div>
              <div class="action-text">
                <h3 class="action-title action-title--sm">家长端</h3>
                <p class="action-desc action-desc--sm">查看积分、宠物、通知</p>
              </div>
              <span class="action-arrow">→</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 功能标签 -->
      <section class="features">
        <span v-for="f in ['⭐积分激励','🌟宠物进化','🏆排行榜','📢实时广播','🤖AI助教','🐳Docker部署']" :key="f" class="feature-item">{{ f }}</span>
      </section>
    </main>

    <!-- 页脚 -->
    <footer class="home-footer">
      <span>学趣星球</span><span class="footer-sep">·</span>
      <span>MIT 开源许可证</span><span class="footer-sep">·</span>
      <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="footer-link">GitHub</a>
    </footer>
  </div>
</template>

<style scoped>
.home {
  position: relative; min-height: 100vh;
  background: #FBFBFD; color: #1D1D1F;
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text',
    'Helvetica Neue', 'Noto Sans SC', sans-serif;
  -webkit-font-smoothing: antialiased;
  overflow-x: hidden;
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
  padding: 16px 32px; max-width: 1100px; margin: 0 auto; width: 100%; box-sizing: border-box;
}
.nav-brand { display: flex; align-items: center; gap: 10px; }
.nav-logo { font-size: 28px; }
.nav-name { font-size: 17px; font-weight: 700; letter-spacing: -.02em; }
.nav-actions { display: flex; align-items: center; gap: 4px; }
.nav-link {
  padding: 8px 16px; border: none; border-radius: 9999px; background: transparent;
  color: #6E6E73; font-size: 13px; font-weight: 500; cursor: pointer;
  transition: all .2s; text-decoration: none; font-family: inherit;
}
.nav-link:hover { background: rgba(0,0,0,.04); color: #1D1D1F; }
.nav-link--muted { color: #AEAEB2; font-size: 12px; }
.nav-link--muted:hover { color: #6E6E73; }
.nav-link--icon { display: flex; align-items: center; padding: 8px; color: #AEAEB2; }
.nav-link--icon:hover { color: #1D1D1F; }

/* 主内容 */
.home-main {
  position: relative; z-index: 1; flex: 1;
  max-width: 920px; width: 100%; margin: 0 auto;
  padding: 40px 32px 60px; box-sizing: border-box;
  display: flex; flex-direction: column; align-items: center;
}

/* Hero */
.hero { text-align: center; margin-bottom: 48px; max-width: 640px; }
.hero-badge {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 16px; background: rgba(255,255,255,.8);
  backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(0,0,0,.06); border-radius: 9999px;
  font-size: 13px; color: #6E6E73; margin-bottom: 28px;
}
.hero-badge-dot { width: 6px; height: 6px; background: #34C759; border-radius: 50%; box-shadow: 0 0 6px rgba(52,199,89,.4); }
.hero-title { font-size: 52px; font-weight: 800; line-height: 1.1; letter-spacing: -.03em; margin: 0 0 20px; }
.hero-title-grad { background: linear-gradient(135deg,#5E5CE6,#FF375F,#FF9F0A); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.hero-desc { font-size: 18px; color: #86868B; margin: 0; }

/* 操作区 */
.actions { width: 100%; display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 48px; }
.actions-secondary { display: flex; flex-direction: column; gap: 12px; }

/* 大屏卡片（深色主题） */
.action-card--primary {
  background: linear-gradient(145deg,#1a1a2e,#16213e); border: none !important;
  position: relative; overflow: hidden; cursor: default;
}
.action-card--primary::before {
  content: ''; position: absolute; inset: 0;
  background: radial-gradient(ellipse at 20% 20%,rgba(120,80,255,.3),transparent 60%),
              radial-gradient(ellipse at 80% 80%,rgba(30,140,220,.2),transparent 60%);
  pointer-events: none;
}
.action-card--primary .action-title,
.action-card--primary .action-desc,
.action-card--primary .action-hint { color: rgba(255,255,255,.9) !important; }
.action-card--primary .action-desc,
.action-card--primary .action-hint { color: rgba(255,255,255,.6) !important; }
.action-card--primary .action-hint { color: rgba(255,255,255,.35) !important; }
.action-card--primary .action-icon { font-size: 40px; }
.action-card-glow {
  position: absolute; width: 200px; height: 200px;
  background: radial-gradient(circle,rgba(120,80,255,.15),transparent);
  top: -60px; right: -60px; border-radius: 50%; pointer-events: none;
}

/* 通用卡片 */
.action-card {
  background: rgba(255,255,255,.75); backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);
  border: 1px solid rgba(0,0,0,.06); border-radius: 20px;
  transition: all .3s cubic-bezier(.16,1,.3,1); position: relative; cursor: pointer;
}
.action-card:hover { transform: translateY(-2px); box-shadow: 0 12px 40px rgba(0,0,0,.08); }
.action-card--primary:hover { transform: translateY(-2px); box-shadow: 0 16px 48px rgba(26,26,46,.3); }
.action-card-content { padding: 32px; position: relative; z-index: 1; }
.action-icon { font-size: 32px; margin-bottom: 12px; }
.action-icon--sm { font-size: 24px; margin-bottom: 0; flex-shrink: 0; }
.action-title { font-size: 22px; font-weight: 700; margin: 0 0 4px; letter-spacing: -.02em; }
.action-title--sm { font-size: 16px; margin: 0 0 2px; }
.action-desc { font-size: 14px; color: #86868B; margin: 0 0 16px; }
.action-desc--sm { font-size: 12px; margin: 0; color: #AEAEB2; }

/* 班级码输入 */
.code-form { margin-top: 4px; }
.code-field { display: flex; gap: 8px; }
.code-input {
  flex: 1; padding: 12px 16px; border-radius: 12px;
  border: 1.5px solid rgba(255,255,255,.15);
  background: rgba(255,255,255,.08); color: #fff;
  font-size: 16px; font-weight: 600; letter-spacing: .1em;
  outline: none; transition: all .2s;
  font-family: 'SF Mono', 'SF Pro Text', 'Noto Sans SC', monospace;
}
.code-input::placeholder { color: rgba(255,255,255,.3); font-weight: 400; letter-spacing: 0; }
.code-input:focus { border-color: rgba(120,80,255,.6); box-shadow: 0 0 0 4px rgba(120,80,255,.15); }
.code-btn {
  padding: 12px 24px; border-radius: 12px; border: none;
  background: linear-gradient(135deg,#7C3AED,#6D28D9); color: #fff;
  font-size: 15px; font-weight: 600; cursor: pointer;
  transition: all .2s; white-space: nowrap; font-family: inherit;
}
.code-btn:hover:not(:disabled) { background: linear-gradient(135deg,#8B5CF6,#7C3AED); transform: scale(1.02); box-shadow: 0 8px 24px rgba(124,58,237,.3); }
.code-btn:disabled { opacity: .4; cursor: not-allowed; }
.code-error { color: #FCA5A5; font-size: 12px; margin: 8px 0 0; }
.action-hint { font-size: 12px; margin: 12px 0 0; }

/* 小卡片 */
.action-card--small .action-card-content { padding: 20px 24px; }
.action-card-content--row { display: flex; align-items: center; gap: 16px; }
.action-text { flex: 1; }
.action-arrow { font-size: 18px; color: #C7C7CC; transition: all .2s; }
.action-card--small:hover .action-arrow { color: #5E5CE6; transform: translateX(4px); }

/* 功能标签 */
.features { display: flex; gap: 10px; flex-wrap: wrap; justify-content: center; margin-bottom: 20px; }
.feature-item {
  padding: 6px 14px; background: rgba(255,255,255,.6);
  backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(0,0,0,.04); border-radius: 9999px;
  font-size: 13px; color: #6E6E73; font-weight: 500;
}

/* 页脚 */
.home-footer {
  position: relative; z-index: 1;
  display: flex; align-items: center; justify-content: center; gap: 8px;
  padding: 24px 32px; font-size: 12px; color: #AEAEB2;
}
.footer-sep { color: #D2D2D7; }
.footer-link { color: #6E6E73; text-decoration: none; transition: color .2s; }
.footer-link:hover { color: #1D1D1F; }

@media (max-width: 768px) {
  .home-nav { padding: 12px 20px; }
  .nav-link--muted { display: none; }
  .home-main { padding: 24px 20px 40px; }
  .hero-title { font-size: 36px; }
  .actions { grid-template-columns: 1fr; gap: 16px; }
  .action-card--primary .action-card-content { padding: 24px; }
}
</style>
