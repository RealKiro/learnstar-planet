<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiPost } from '@/utils/api'

const router = useRouter()
const displayCode = ref('')
const loading = ref(false)
const codeError = ref('')

const slides = [
  { icon: '🌌', title: '让每个孩子的努力', highlight: '都被看见', desc: '积分激励 宠物养成 AI 助教 多端同步' },
  { icon: '⚡', title: '覆盖班级管理', highlight: '全场景', desc: '积分规则 宠物进化 排行榜 通知公告' },
  { icon: '🌟', title: '积分变经验', highlight: '驱动成长', desc: '12级宠物进化，从卵到传说' },
  { icon: '🐳', title: 'Docker 一键部署', highlight: '自由选择', desc: '4种数据库支持，数据完全自主掌控' },
]

const currentSlide = ref(0)
let slideTimer: ReturnType<typeof setInterval>

function onCodeInput(e: Event) {
  const input = e.target as HTMLInputElement
  displayCode.value = input.value.toUpperCase().replace(/[^0-9A-Z-]/g, '')
  codeError.value = ''
}

async function goToClassroom() {
  const code = displayCode.value.trim()
  if (!code) { codeError.value = '请输入班级码'; return }
  loading.value = true; codeError.value = ''
  try {
    const res = await apiPost<{ data: { token: string; class_id: number; class_name: string; grade: string; student_count: number } }>(
      '/api/v1/auth/class/login', { class_code: code }
    )
    sessionStorage.setItem('class_token', res.data.token)
    sessionStorage.setItem('class_info', JSON.stringify({
      id: res.data.class_id, name: res.data.class_name, grade: res.data.grade, student_count: res.data.student_count,
    }))
    router.replace({ name: 'classroom-overview' })
  } catch { codeError.value = '班级码无效，请检查后重试'; loading.value = false }
}

function goLogin(role: string) { router.push({ name: 'login', query: { role } }) }

onMounted(() => {
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  if (code) { displayCode.value = code.toUpperCase(); goToClassroom() }
})
onUnmounted(() => clearInterval(slideTimer))
</script>

<template>
  <div class="home">
    <div class="bg">
      <div class="bg-orb bg-orb--1"></div>
      <div class="bg-orb bg-orb--2"></div>
    </div>

    <!-- 顶部导航 -->
    <header class="topbar">
      <div class="topbar-brand">🌌 学趣星球</div>
      <div class="topbar-links">
        <button class="topbar-btn topbar-btn--dim" @click="goLogin('teacher')">教师登录</button>
        <button class="topbar-btn topbar-btn--dim" @click="goLogin('admin')">管理员</button>
        <button class="topbar-btn topbar-btn--dim" @click="goLogin('parent')">家长</button>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="topbar-link" title="GitHub">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
      </div>
    </header>

    <main class="main">
      <!-- 左侧：产品介绍轮播 -->
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

      <!-- 右侧：班级码登录卡片 -->
      <section class="right">
        <div class="panel">
          <div class="panel-icon">🐾</div>
          <h2 class="panel-title">班级宠物星球</h2>
          <p class="panel-desc">输入班级码，进入你的班级</p>

          <form class="panel-form" @submit.prevent="goToClassroom">
            <input v-model="displayCode" type="text" class="panel-input"
              placeholder="输入班级码" maxlength="12" autocomplete="off" @input="onCodeInput">
            <button type="submit" class="panel-btn" :disabled="loading || displayCode.length < 3">
              {{ loading ? '⏳ 验证中…' : '🚀 进入班级' }}
            </button>
            <p v-if="codeError" class="panel-error">{{ codeError }}</p>
          </form>

          <p class="panel-footnote">班级码由班主任统一分配</p>

          <div class="panel-divider"><span>其他登录方式</span></div>

          <div class="panel-alt">
            <button class="panel-alt-btn" @click="goLogin('teacher')">👨‍🏫 教师</button>
            <button class="panel-alt-btn" @click="goLogin('parent')">👨‍👩‍👧‍👦 家长</button>
            <button class="panel-alt-btn" @click="goLogin('admin')">⚙️ 管理</button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.home {
  min-height: 100vh; display: flex; flex-direction: column;
  background: #FBFBFD; color: #1D1D1F; position: relative;
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', 'Noto Sans SC', sans-serif;
  -webkit-font-smoothing: antialiased;
}

/* 背景 */
.bg { position: fixed; inset: 0; overflow: hidden; pointer-events: none; z-index: 0; }
.bg-orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.25; }
.bg-orb--1 { width: 500px; height: 500px; background: #C7D2FE; top: -200px; right: -100px; animation: orbFloat 12s ease-in-out infinite; }
.bg-orb--2 { width: 350px; height: 350px; background: #A7F3D0; bottom: -100px; left: -50px; animation: orbFloat 15s ease-in-out infinite reverse; }
@keyframes orbFloat { 0%,100% { transform: translate(0,0); } 50% { transform: translate(40px,-30px); } }

/* 顶栏 */
.topbar {
  position: relative; z-index: 1;
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 32px; max-width: 1100px; margin: 0 auto; width: 100%; box-sizing: border-box;
}
.topbar-brand { font-size: 20px; font-weight: 700; background: linear-gradient(135deg,#5E5CE6,#FF375F); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.topbar-links { display: flex; align-items: center; gap: 8px; }
.topbar-btn {
  padding: 6px 14px; border: 1px solid #E5E5EA; border-radius: 9999px;
  background: #FFFFFF; color: #6E6E73; font-size: 13px; font-weight: 500;
  cursor: pointer; transition: all .2s; font-family: inherit;
}
.topbar-btn:hover { background: #F5F5F7; color: #1D1D1F; }
.topbar-btn--dim { background: transparent; border-color: transparent; color: #86868B; }
.topbar-btn--dim:hover { background: rgba(0,0,0,.04); }
.topbar-link { color: #AEAEB2; padding: 4px; transition: color .2s; display: flex; }
.topbar-link:hover { color: #6E6E73; }

/* 主内容 */
.main {
  flex: 1; display: flex; max-width: 1100px; margin: 0 auto; width: 100%;
  padding: 20px 32px 40px; position: relative; z-index: 1;
}

/* 左侧 */
.left { flex: 1; display: flex; align-items: center; justify-content: center; padding: 48px 64px 48px 0; }
.slide-stage { max-width: 480px; width: 100%; }
.slide-card { }
.slide-icon { font-size: 56px; margin-bottom: 20px; }
.slide-title { font-size: 40px; font-weight: 900; line-height: 1.2; letter-spacing: -1px; margin-bottom: 20px; }
.slide-highlight { background: linear-gradient(135deg,#5E5CE6,#FF375F,#FF9F0A); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.slide-desc { font-size: 16px; color: #86868B; line-height: 1.7; white-space: pre-line; }
.dots { display: flex; gap: 8px; margin-top: 44px; }
.dot { height: 8px; width: 8px; border-radius: 4px; border: none; cursor: pointer; background: #D2D2D7; transition: all 0.3s; }
.dot.active { width: 32px; background: linear-gradient(135deg,#5E5CE6,#818CF8); }

/* 右侧 */
.right { flex: 0 0 380px; display: flex; align-items: center; }
.panel { width: 100%; padding: 36px 32px; background: #FFFFFF; border-radius: 20px; border: 1px solid #F0F0F3; box-shadow: 0 8px 30px rgba(0,0,0,.04); text-align: center; }
.panel-icon { font-size: 48px; margin-bottom: 8px; }
.panel-title { font-size: 22px; font-weight: 800; margin-bottom: 4px; }
.panel-desc { font-size: 14px; color: #86868B; margin-bottom: 24px; }

.panel-form { display: flex; flex-direction: column; gap: 12px; }
.panel-input {
  width: 100%; padding: 14px 18px;
  background: #F5F5F7; border: 1px solid #E5E5EA; border-radius: 12px;
  color: #1D1D1F; font-size: 24px; font-weight: 700;
  text-align: center; letter-spacing: 0.1em;
  outline: none; transition: all 0.2s;
}
.panel-input:focus { border-color: #5E5CE6; background: #FFFFFF; box-shadow: 0 0 0 3px rgba(94,92,230,0.12); }
.panel-input::placeholder { color: #AEAEB2; font-weight: 400; font-size: 16px; letter-spacing: 0; }

.panel-btn {
  width: 100%; padding: 14px; border: none; border-radius: 12px;
  background: linear-gradient(135deg,#5E5CE6,#818CF8); color: #FFFFFF;
  font-size: 16px; font-weight: 600; cursor: pointer;
  box-shadow: 0 4px 14px rgba(94,92,230,0.2); transition: all 0.3s;
}
.panel-btn:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 20px rgba(94,92,230,0.3); }
.panel-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.panel-error { color: #EF4444; font-size: 13px; padding: 8px; background: rgba(239,68,68,0.06); border-radius: 8px; }
.panel-footnote { font-size: 12px; color: #AEAEB2; margin-top: 16px; }

.panel-divider { display: flex; align-items: center; gap: 12px; margin: 20px 0 14px; color: #AEAEB2; font-size: 12px; }
.panel-divider::before, .panel-divider::after { content: ''; flex: 1; height: 1px; background: #E5E5EA; }

.panel-alt { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.panel-alt-btn {
  padding: 10px 8px; border: 1px solid #E5E5EA; border-radius: 10px;
  background: #F5F5F7; color: #6E6E73; font-size: 13px; font-weight: 500;
  cursor: pointer; transition: all 0.2s; font-family: inherit;
}
.panel-alt-btn:hover { background: #FFFFFF; border-color: #D2D2D7; color: #1D1D1F; }

.slide-enter-active { transition: all 0.4s ease; }
.slide-leave-active { transition: all 0.3s ease; }
.slide-enter-from { opacity: 0; transform: translateX(20px); }
.slide-leave-to { opacity: 0; transform: translateX(-20px); }

@media (max-width: 768px) {
  .main { flex-direction: column; padding: 16px; }
  .left { display: none; }
  .right { flex: 1; }
  .panel { box-shadow: none; border-color: transparent; }
}
</style>
