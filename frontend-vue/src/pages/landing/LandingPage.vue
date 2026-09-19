<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiPost } from '@/utils/api'

const router = useRouter()
const displayCode = ref('')
const loading = ref(false)
const codeError = ref('')

/** icon 为 emoji 字符 */
const slides = [
  { icon: '🌌', title: '让每个孩子的努力', highlight: '都被看见', desc: '积分激励 · 宠物养成 · AI 助教 · 多端同步' },
  { icon: '⚡', title: '覆盖班级管理', highlight: '全场景', desc: '积分规则 · 宠物进化 · 排行榜 · 通知公告' },
  { icon: '🌟', title: '积分变经验', highlight: '驱动成长', desc: '12 级宠物进化，从卵到传说' },
  { icon: '🐳', title: 'Docker 一键部署', highlight: '数据自主', desc: '4 种数据库支持，学校局域网离线可用' },
]

const features = [
  { icon: '⭐', title: '积分激励', desc: '自定义规则，实时到大屏' },
  { icon: '🐣', title: '宠物养成', desc: '126 种宠物，12 级进化' },
  { icon: '🤖', title: 'AI 助教', desc: '30+ 供应商，开箱即用' },
  { icon: '📊', title: '数据报表', desc: '趋势 / 分布 / 一键导出' },
]

const currentSlide = ref(0)
let slideTimer: ReturnType<typeof setInterval>

function onCodeInput(e: Event) {
  const input = e.target as HTMLInputElement
  displayCode.value = input.value.replace(/[^0-9A-Za-z]/g, '').toUpperCase().slice(0, 8)
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
    router.replace({ name: 'teacher-dashboard-basic' })
  } catch { codeError.value = '班级码无效，请检查后重试'; loading.value = false }
}

function goLogin() { router.push({ name: 'login' }) }

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
    <!-- 顶部导航 -->
    <header class="topbar">
      <div class="topbar-brand"><span class="topbar-brand-icon">🌌</span><span class="topbar-brand-text">学宠星球</span></div>
      <div class="topbar-links">
        <button class="topbar-btn" @click="goLogin">👨‍🏫 教师</button>
        <button class="topbar-btn" @click="goLogin">⚙️ 管理</button>
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="topbar-link" title="GitHub">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
      </div>
    </header>

    <main class="main">
      <!-- 左侧：产品介绍轮播 + 功能矩阵 -->
      <section class="left">
        <div class="badge-row">
          <span class="badge"><span class="badge-dot"></span> MIT 开源</span>
          <span class="badge">✨ 完全免费</span>
          <span class="badge">🏠 离线自托管</span>
        </div>

        <div class="slide-stage">
          <Transition name="slide" mode="out-in">
            <div :key="currentSlide" class="slide-card">
              <div class="slide-icon"><span>{{ slides[currentSlide].icon }}</span></div>
              <h1 class="slide-title">
                {{ slides[currentSlide].title }}<br>
                <span class="slide-highlight">{{ slides[currentSlide].highlight }}</span>
              </h1>
              <p class="slide-desc">{{ slides[currentSlide].desc }}</p>
            </div>
          </Transition>
          <div class="dots">
            <button
              v-for="(_, i) in slides" :key="i" class="dot"
              :class="{ active: currentSlide === i }"
              :aria-label="`切换到第 ${i + 1} 张`"
              @click="currentSlide = i"
            ></button>
          </div>
        </div>

        <div class="feature-grid">
          <div v-for="f in features" :key="f.title" class="feature-card">
            <span class="feature-card__icon">{{ f.icon }}</span>
            <div class="feature-card__body">
              <div class="feature-card__title">{{ f.title }}</div>
              <div class="feature-card__desc">{{ f.desc }}</div>
            </div>
          </div>
        </div>
      </section>

      <!-- 右侧：仅班级码登录 -->
      <section class="right">
        <div class="panel">
          <div class="panel-icon"><span>🐾</span></div>
          <h2 class="panel-title">进入我的班级</h2>
          <p class="panel-desc">输入班主任发的班级码，无需注册</p>

          <form class="panel-form" @submit.prevent="goToClassroom">
            <label class="panel-label" for="class-code">班级码</label>
            <input
              id="class-code" v-model="displayCode" type="text" class="panel-input"
              placeholder="如 LS11" maxlength="8" autocomplete="off" spellcheck="false"
              @input="onCodeInput"
            >
            <p class="panel-hint">班级码由管理员统一分配</p>
            <button type="submit" class="panel-btn" :disabled="loading || displayCode.length < 3">
              {{ loading ? '⏳ 验证中…' : '🚀 进入班级' }}
            </button>
            <p v-if="codeError" class="panel-error">{{ codeError }}</p>
          </form>

          <div class="panel-divider">或使用账号登录</div>
          <div class="panel-alt">
            <button class="panel-alt-btn" @click="goLogin">👨‍🏫 教师登录</button>
            <button class="panel-alt-btn" @click="goLogin">⚙️ 管理员登录</button>
          </div>
        </div>
      </section>
    </main>

    <!-- 页脚 -->
    <footer class="footer">
      <span>© 2026 学宠星球 · LearnStar Planet</span>
      <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" rel="noopener">GitHub</a>
      <span class="footer-sep">·</span>
      <span>MIT License</span>
      <span class="footer-sep">·</span>
      <span>局域网离线可用</span>
    </footer>
  </div>
</template>

<style scoped>
/* ===================================================================
 * 入口页 — 视觉对标重绘（2026-09-19）
 * 去掉：彩虹渐变顶条、光晕 orb、星空层、渐变文字裁剪、玻璃拟态。
 * 保留：两栏职责分工、班级码的居中大字输入（课堂大屏的功能性需求，不削平）。
 * 全量走 --ui-* 令牌 → 暗色无需逐条手写覆盖层。
 * =================================================================== */
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--ui-bg);
  color: var(--ui-fg);
  position: relative;
}

/* 唯一保留的氛围：静态品牌色晕染（无动画、无 blur 圆斑） */
.home::before {
  content: '';
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  background:
    radial-gradient(680px 420px at 88% -8%, var(--ui-brand-soft), transparent 70%),
    radial-gradient(520px 360px at -6% 96%, var(--ui-muted), transparent 72%);
}

/* 顶栏 */
.topbar {
  position: relative; z-index: 1;
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 28px; max-width: 1120px; margin: 0 auto; width: 100%; box-sizing: border-box;
}
.topbar-brand { display: flex; align-items: center; gap: 9px; font-size: 15px; font-weight: 650; letter-spacing: -0.01em; }
.topbar-mark {
  width: 26px; height: 26px; border-radius: var(--ui-r-sm);
  background: var(--ui-brand); color: var(--ui-brand-fg);
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.topbar-links { display: flex; align-items: center; gap: 6px; }
.topbar-btn {
  display: inline-flex; align-items: center; gap: 6px;
  height: 32px; padding: 0 12px; border: 1px solid var(--ui-border); border-radius: var(--ui-r-md);
  background: var(--ui-card); color: var(--ui-fg-muted); font-size: 13px; font-weight: 500;
  cursor: pointer; transition: 0.15s; font-family: inherit;
}
.topbar-btn:hover { background: var(--ui-muted); color: var(--ui-fg); border-color: var(--ui-border-strong); }
.topbar-icon {
  width: 32px; height: 32px; border: 1px solid var(--ui-border); border-radius: var(--ui-r-md);
  display: flex; align-items: center; justify-content: center;
  color: var(--ui-fg-subtle); transition: 0.15s;
}
.topbar-icon:hover { background: var(--ui-muted); color: var(--ui-fg); border-color: var(--ui-border-strong); }

/* 主内容 */
.main {
  flex: 1; display: flex; align-items: center; gap: 48px;
  max-width: 1120px; margin: 0 auto; width: 100%;
  padding: 24px 28px 40px; position: relative; z-index: 1; box-sizing: border-box;
}

/* 左侧 */
.left { flex: 1; min-width: 0; display: flex; flex-direction: column; justify-content: center; }
.badge-row { display: flex; gap: 6px; margin-bottom: 26px; flex-wrap: wrap; }
.badge {
  display: inline-flex; align-items: center; gap: 6px;
  height: 24px; padding: 0 10px; border-radius: var(--ui-r-sm); font-size: 12px; font-weight: 500;
  background: var(--ui-card); border: 1px solid var(--ui-border); color: var(--ui-fg-muted);
}
.badge svg { color: var(--ui-fg-subtle); }
.badge-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--c-green); flex-shrink: 0; }

.slide-stage { max-width: 520px; width: 100%; }
.slide-icon {
  width: 44px; height: 44px; border-radius: var(--ui-r-lg); margin-bottom: 18px;
  display: flex; align-items: center; justify-content: center;
  background: var(--ui-brand-soft); color: var(--ui-brand);
  border: 1px solid var(--ui-brand-border);
}
.slide-title { font-size: 36px; font-weight: 700; line-height: 1.22; letter-spacing: -0.025em; margin: 0 0 14px; }
/* 高亮词用纯色品牌，不再做渐变文字裁剪（低 DPI 下笔画会发灰） */
.slide-highlight { color: var(--ui-brand); }
.slide-desc { font-size: 14.5px; color: var(--ui-fg-muted); line-height: 1.7; }
.dots { display: flex; gap: 6px; margin-top: 30px; }
.dot {
  height: 6px; width: 6px; border-radius: 3px; border: none; cursor: pointer;
  background: var(--ui-border-strong); transition: all 0.25s var(--ease-smooth); padding: 0;
}
.dot:hover { background: var(--ui-fg-subtle); }
.dot.active { width: 24px; background: var(--ui-brand); }

/* 功能矩阵 */
.feature-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; max-width: 520px; margin-top: 28px; }
.feature-card {
  display: flex; align-items: center; gap: 10px; padding: 12px 14px;
  background: var(--ui-card); border: 1px solid var(--ui-border); border-radius: var(--ui-r-lg);
  transition: border-color 0.18s var(--ease-smooth), box-shadow 0.18s var(--ease-smooth);
}
.feature-card:hover { border-color: var(--ui-border-strong); box-shadow: var(--ui-shadow-xs); }
.feature-card__icon {
  width: 30px; height: 30px; border-radius: var(--ui-r-sm); flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  background: var(--ui-muted); color: var(--ui-fg-muted);
}
.feature-card__body { min-width: 0; }
.feature-card__title { font-size: 13px; font-weight: 600; color: var(--ui-fg); }
.feature-card__desc { font-size: 11.5px; color: var(--ui-fg-muted); margin-top: 1px; }

/* 右侧 */
.right { flex: 0 0 372px; display: flex; align-items: center; }
.panel {
  width: 100%; padding: 24px; background: var(--ui-card);
  border: 1px solid var(--ui-border); border-radius: var(--ui-r-xl);
  box-shadow: var(--ui-shadow-xs);
}
.panel-head { display: flex; align-items: center; gap: 11px; margin-bottom: 18px; }
.panel-icon {
  width: 36px; height: 36px; border-radius: var(--ui-r-md); flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  background: var(--ui-brand-soft); color: var(--ui-brand);
  border: 1px solid var(--ui-brand-border);
}
.panel-head-text { min-width: 0; }
.panel-title { font-size: 15px; font-weight: 650; letter-spacing: -0.01em; margin: 0; }
.panel-desc { font-size: 12.5px; color: var(--ui-fg-muted); margin-top: 1px; }

.panel-form { display: flex; flex-direction: column; }
.panel-label { display: block; font-size: 12.5px; font-weight: 550; color: var(--ui-fg); margin-bottom: 6px; }
/* 班级码保留「居中 + 字距」——课堂大屏输入的功能性需求，与常规输入框有意区分 */
.panel-input {
  width: 100%; height: 48px; box-sizing: border-box;
  background: var(--ui-card); border: 1px solid var(--ui-border); border-radius: var(--ui-r-md);
  color: var(--ui-fg); font-size: 19px; font-weight: 650;
  text-align: center; letter-spacing: 0.14em; text-indent: 0.14em;
  outline: none; transition: border-color 0.15s, box-shadow 0.15s; font-family: inherit;
}
.panel-input:focus { border-color: var(--ui-brand); box-shadow: 0 0 0 3px var(--ui-ring); }
.panel-input::placeholder { color: var(--ui-fg-subtle); font-weight: 400; font-size: 14px; letter-spacing: 0; text-indent: 0; }

.panel-hint { font-size: 11.5px; color: var(--ui-fg-subtle); margin-top: 6px; }

.panel-btn {
  width: 100%; height: 42px; margin-top: 14px; border: 1px solid var(--ui-brand); border-radius: var(--ui-r-md);
  background: var(--ui-brand); color: var(--ui-brand-fg);
  font-size: 14.5px; font-weight: 600; cursor: pointer; font-family: inherit;
  display: inline-flex; align-items: center; justify-content: center; gap: 7px;
  transition: background 0.15s, transform 0.08s;
}
.panel-btn:hover:not(:disabled) { background: var(--ui-brand-hover); }
.panel-btn:active:not(:disabled) { transform: scale(0.99); }
.panel-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.panel-error {
  display: flex; align-items: center; gap: 6px;
  color: var(--color-danger-text); font-size: 12.5px; margin: 10px 0 0;
  padding: 8px 10px; background: var(--c-red-bg); border: 1px solid var(--c-red-border);
  border-radius: var(--ui-r-md);
}

.panel-divider { display: flex; align-items: center; gap: 10px; margin: 18px 0 12px; color: var(--ui-fg-subtle); font-size: 11.5px; }
.panel-divider::before, .panel-divider::after { content: ''; flex: 1; height: 1px; background: var(--ui-border); }

.panel-alt { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; }
.panel-alt-btn {
  height: 36px; padding: 0 8px; border: 1px solid var(--ui-border); border-radius: var(--ui-r-md);
  background: var(--ui-card); color: var(--ui-fg); font-size: 13px; font-weight: 500;
  cursor: pointer; transition: 0.15s; font-family: inherit;
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
}
.panel-alt-btn:hover { background: var(--ui-muted); border-color: var(--ui-border-strong); }
.panel-alt-btn svg { color: var(--ui-fg-subtle); }

/* 页脚 */
.footer {
  position: relative; z-index: 1;
  display: flex; align-items: center; justify-content: center; gap: 10px;
  padding: 18px 28px 22px; font-size: 12px; color: var(--ui-fg-subtle);
}
.footer a { color: var(--ui-fg-muted); transition: color 0.15s; }
.footer a:hover { color: var(--ui-fg); }
.footer-sep { opacity: 0.5; }

.slide-enter-active { transition: all 0.35s ease; }
.slide-leave-active { transition: all 0.25s ease; }
.slide-enter-from { opacity: 0; transform: translateX(16px); }
.slide-leave-to { opacity: 0; transform: translateX(-16px); }

@media (max-width: 860px) {
  .main { flex-direction: column; gap: 24px; padding: 16px; }
  .left { display: none; }
  .right { flex: 1; width: 100%; }
  .panel { box-shadow: none; }
  .footer { flex-wrap: wrap; gap: 6px; }
}
</style>
