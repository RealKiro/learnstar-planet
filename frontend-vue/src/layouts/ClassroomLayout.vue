<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { apiPost, apiGet } from '@/utils/api'
import { getAllSeries, getSeriesName } from '@/utils/petData'

const router = useRouter()
const route = useRoute()

const classInfo = ref<{ id: number; name: string; student_count?: number } | null>(null)
const showVoteModal = ref(false)
const voteSeries = ref('myth')
const voting = ref(false)
const voteDone = ref(false)
const aiEnabled = ref(false)

const activeNav = computed(() => String(route.name))

const baseNavItems = [
  { page: 'classroom-overview', label: '班级总览', icon: '🏠' },
  { page: 'classroom-timetable', label: '今日课表', icon: '🗓️' },
  { page: 'classroom-scores', label: '课堂评价', icon: '✏️' },
  { page: 'classroom-leaderboard', label: '排行榜单', icon: '🏆' },
  { page: 'classroom-pk', label: '年级战场', icon: '⚔️' },
  { page: 'classroom-pokedex', label: '宠物图鉴', icon: '📚' },
]

const navItems = computed(() => {
  const items = [...baseNavItems]
  if (aiEnabled.value) {
    items.push({ page: 'classroom-ai', label: 'AI 助手', icon: '🤖' })
  }
  return items
})

const allSeries = getAllSeries()
const lastEventId = ref(0)

function navigate(name: string) { router.push({ name }) }
function goToLogin() { sessionStorage.clear(); router.push({ name: 'landing' }) }

// ===== 广播/通知接收（SSE + 轮询降级） =====
const currentBroadcast = ref<{
  id: number; type: string; content: string; display_seconds: number; created_at: string
} | null>(null)

const currentNotice = ref<{
  id: number; title: string; content: string; type: string; published_at: string
} | null>(null)

let pollTimer: ReturnType<typeof setInterval> | null = null

function startPolling() {
  if (pollTimer) return
  pollTimer = setInterval(pollEvents, 5000)
  pollEvents()
}

async function pollEvents() {
  const token = sessionStorage.getItem('class_token')
  if (!token) return

  try {
    const res = await apiGet<{ data: { events: Array<{ id: number; type: string; data: any }>; last_event_id?: number } }>(
      '/api/v1/display/poll',
      { params: { token, last_event_id: lastEventId.value } }
    )
    const events = res.data?.events || []
    if (events.length === 0) return

    for (const ev of events) {
      if (ev.id && ev.id > lastEventId.value) lastEventId.value = ev.id
      if (ev.data) handleEvent(ev.type, ev.data)
    }
  } catch { /* polling fails silently */ }
}

function handleEvent(type: string, data: any) {
  if (type === 'notice') showNotice(data)
  else showBroadcast(data)
}

function showBroadcast(data: any) {
  currentBroadcast.value = {
    id: data.id,
    type: data.type || 'banner',
    content: data.content || '',
    display_seconds: data.display_seconds || 10,
    created_at: data.created_at || '',
  }
  // 自动关闭
  const sec = (data.display_seconds || 10) * 1000
  if (sec > 0) {
    setTimeout(() => {
      if (currentBroadcast.value?.id === data.id) {
        currentBroadcast.value = null
      }
    }, sec)
  }
}

function showNotice(data: any) {
  currentNotice.value = {
    id: data.id,
    title: data.title || '',
    content: data.content || '',
    type: data.type || 'info',
    published_at: data.published_at || '',
  }
  // 通知停留 15 秒后自动关闭
  setTimeout(() => {
    if (currentNotice.value?.id === data.id) {
      currentNotice.value = null
    }
  }, 15000)
}

function dismissBroadcast() { currentBroadcast.value = null }
function dismissNotice() { currentNotice.value = null }

async function confirmVote() {
  voting.value = true
  const token = sessionStorage.getItem('class_token') || ''
  try {
    await apiPost('/api/v1/display/switch-series', { token, series_id: voteSeries.value })
    voteDone.value = true
    sessionStorage.setItem('class_series', voteSeries.value)
    setTimeout(() => { showVoteModal.value = false }, 2000)
  } catch { /* ignore */ } finally { voting.value = false }
}

onMounted(() => {
  const ci = sessionStorage.getItem('class_info')
  if (ci) classInfo.value = JSON.parse(ci)
  // 如果已有系列配置则不弹投票
  if (sessionStorage.getItem('class_series')) {
    showVoteModal.value = false
  }

  // 检查 AI 功能状态
  checkAiStatus()

  // 使用轮询接收广播和通知（php artisan serve 不支持 SSE 长连接）
  lastEventId.value = parseInt(sessionStorage.getItem('last_event_id') || '0', 10)
  startPolling()
})

async function checkAiStatus() {
  const token = sessionStorage.getItem('class_token') || ''
  if (!token) return
  try {
    const res = await apiGet<{ data: { enabled: boolean } }>('/api/v1/display/ai/settings', { params: { token } })
    aiEnabled.value = res.data?.enabled || false
  } catch { /* ignore */ }
}

onUnmounted(() => {
  
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  sessionStorage.setItem('last_event_id', String(lastEventId.value))
})
</script>

<template>
  <div class="app-shell">
    <!-- 首次使用投票弹窗 -->
    <Transition name="fade">
      <div v-if="showVoteModal" @click.self="() => {}" class="vote-mask">
        <div class="vote-panel">
          <div v-if="!voteDone">
            <div class="vote-emoji">🎉</div>
            <h2 class="vote-title">欢迎来到学宠星球！</h2>
            <p class="vote-desc">
              请全班投票选择你们喜欢的宠物类别<br>
              <span class="vote-hint">选定后每人可免费选择一只心仪的宠物</span>
            </p>
            <div class="vote-grid">
              <button v-for="s in allSeries" :key="s.id" @click="voteSeries = s.id"
                :style="{
                  padding:'16px 12px', borderRadius:'16px', cursor:'pointer', transition:'0.2s', fontFamily:'inherit',
                  border: voteSeries === s.id ? '2px solid var(--color-primary)' : '1px solid var(--tint-3)',
                  background: voteSeries === s.id ? 'rgba(167,139,250,0.1)' : 'var(--tint-1)',
                  color: voteSeries === s.id ? 'var(--color-primary)' : 'var(--color-text)',
                }">
                <div class="vote-item-emoji">{{ s.emoji }}</div>
                <div class="vote-item-name">{{ s.name }}</div>
                <div class="vote-item-meta">{{ s.species.length }}种宠物</div>
              </button>
            </div>
            <button @click="confirmVote" :disabled="voting" class="vote-submit">
              {{ voting ? '投票中...' : '✅ 选择「' + getSeriesName(voteSeries) + '」系列' }}
            </button>
          </div>
          <div v-else>
            <div class="vote-done-emoji">🎊</div>
            <h2 class="vote-done-title">选择成功！</h2>
            <p class="vote-done-desc">
              已选定「{{ getSeriesName(voteSeries) }}」系列
            </p>
            <p class="vote-done-hint">
              现在去为每位同学免费选择宠物吧！
            </p>
          </div>
        </div>
      </div>
    </Transition>

    <!-- ===== 广播覆盖层 ===== -->

    <!-- 全屏广播 -->
    <Transition name="fade">
      <div v-if="currentBroadcast && currentBroadcast.type === 'fullscreen'"
        @click="dismissBroadcast" class="bc-full-mask">
        <div class="bc-full-inner">
          <div class="bc-full-emoji">📡</div>
          <div class="bc-full-text">{{ currentBroadcast.content }}</div>
          <div class="bc-full-hint">点击任意位置关闭</div>
        </div>
      </div>
    </Transition>

    <!-- 弹窗广播 -->
    <Transition name="pop">
      <div v-if="currentBroadcast && currentBroadcast.type === 'popup'"
        @click.self="dismissBroadcast" class="bc-pop-mask">
        <div class="bc-pop-panel">
          <div class="bc-pop-emoji">📢</div>
          <div class="bc-pop-text">{{ currentBroadcast.content }}</div>
          <button @click="dismissBroadcast" class="bc-pop-btn">
            我知道了
          </button>
        </div>
      </div>
    </Transition>

    <!-- 横幅广播 -->
    <Transition name="slide-down">
      <div v-if="currentBroadcast && currentBroadcast.type === 'banner'"
        @click="dismissBroadcast" class="bc-banner">
        <span class="bc-banner-text">📢 {{ currentBroadcast.content }}</span>
        <span class="bc-banner-hint">点击关闭</span>
      </div>
    </Transition>

    <!-- 通知提示 -->
    <Transition name="slide-down">
      <div v-if="currentNotice" class="notice-toast">
        <div class="notice-body">
          <span class="notice-icon">📋</span>
          <div class="flex-1">
            <div class="notice-title">{{ currentNotice.title }}</div>
            <div class="notice-text">{{ currentNotice.content }}</div>
          </div>
          <button @click="dismissNotice" class="notice-close">✕</button>
        </div>
      </div>
    </Transition>

    <nav class="sidebar">
      <div class="logo">
        <div class="brand"><span>🌌</span> 学宠星球</div>
      </div>

      <div class="class-badge" v-if="classInfo">
        <span class="class-name">{{ classInfo.name }}</span>
        <span class="class-count">{{ classInfo.student_count || '--' }} 人</span>
      </div>

      <div class="nav-list">
        <button v-for="item in navItems" :key="item.page"
          :class="['nav-item', { active: activeNav === item.page }]"
          @click="navigate(item.page)">
          <span class="icon">{{ item.icon }}</span> {{ item.label }}
        </button>
      </div>

      <!-- 底部：退出按钮 -->
      <div class="sidebar-footer">
        <button class="exit-btn" @click="goToLogin">✕ 退出班级</button>
      </div>
    </nav>

    <main class="main-content">
      <router-view v-slot="{ Component }">
        <component :is="Component" :key="$route.fullPath" />
      </router-view>
    </main>
  </div>
</template>

<style scoped>
.app-shell { display: flex; min-height: 100vh; }
.sidebar {
  width: var(--md-sidebar-width); background: var(--ui-bg-subtle);
  border-right: none;
  padding: 20px 12px 16px; display: flex; flex-direction: column;
  position: sticky; top: 0; height: 100vh; flex-shrink: 0;
  z-index: 10;
}
/* 品牌标识：纯色 + 字重（去掉渐变裁剪，小字号下渐变会让笔画发灰） */
.logo { display: flex; align-items: center; gap: 9px; padding: 2px 8px 14px; font-size: 16px; font-weight: 650; letter-spacing: -0.01em; color: var(--ui-fg); }
.brand { display: flex; align-items: center; gap: 9px; }
.brand-mark { width: 24px; height: 24px; border-radius: var(--ui-r-sm); background: var(--ui-brand); color: var(--ui-brand-fg); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.class-badge { text-align: center; padding: 10px; margin-bottom: 12px; background: var(--ui-card); border-radius: var(--ui-r-lg); border: 1px solid var(--ui-border); }
.class-name { font-size: 15px; font-weight: 650; display: block; color: var(--ui-fg); }
.class-count { font-size: 12px; color: var(--ui-fg-muted); }
.nav-list { display: flex; flex-direction: column; gap: 2px; flex: 1; }
/* 大屏场景保留较大的行高与字号（40px / 15px），但层级语言与教师/管理端一致 */
.nav-item { display: flex; align-items: center; gap: 12px; height: 40px; padding: 0 10px; border-radius: var(--ui-r-md); cursor: pointer; transition: 0.15s; color: var(--ui-fg-muted); border: none; background: transparent; width: 100%; font-size: 15px; font-weight: 500; font-family: inherit; text-align: left; }
.nav-item:hover { background: var(--ui-muted); color: var(--ui-fg); }
.nav-item.active { background: var(--ui-card); color: var(--ui-fg); font-weight: 600; box-shadow: var(--ui-shadow-xs); }
.nav-item .icon { color: var(--ui-fg-subtle); flex-shrink: 0; transition: color 0.15s; }
.nav-item:hover .icon { color: var(--ui-fg-muted); }
.nav-item.active .icon { color: var(--ui-brand); }

.sidebar-footer { border-top: 1px solid var(--ui-border); padding-top: 10px; }
.exit-btn { width: 100%; height: 36px; padding: 0 12px; border-radius: var(--ui-r-md); border: 1px solid transparent; background: var(--c-red-bg); color: var(--color-danger-text); font-size: 14px; font-weight: 500; cursor: pointer; transition: 0.15s; font-family: inherit; display: flex; align-items: center; justify-content: center; gap: 7px; }
.exit-btn:hover { background: rgba(255,100,100,0.15); }

.main-content { flex: 1; padding: 24px 28px 40px; max-width: calc(100% - var(--md-sidebar-width)); overflow-x: hidden; background: var(--ui-bg); }

/* 过渡动画 */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.pop-enter-active { transition: all 0.25s ease-out; }
.pop-leave-active { transition: all 0.15s ease-in; }
.pop-enter-from { opacity: 0; transform: scale(0.9); }
.pop-leave-to { opacity: 0; transform: scale(0.95); }

.slide-down-enter-active { transition: all 0.3s ease-out; }
.slide-down-leave-active { transition: all 0.2s ease-in; }
.slide-down-enter-from { opacity: 0; transform: translateY(-100%); }
.slide-down-leave-to { opacity: 0; transform: translateY(-100%); }

@media (max-width: 768px) {
  .sidebar { width: 100%; height: auto; position: sticky; flex-direction: row; flex-wrap: wrap; align-items: center; padding: 10px 14px; border-bottom: 1px solid var(--ui-border); }
  .logo { padding: 0; font-size: 16px; flex: 1; }
  .class-badge { display: none; }
  .nav-list { flex-direction: row; gap: 4px; flex: 2; justify-content: flex-end; }
  .nav-item { height: 34px; padding: 0 10px; font-size: 14px; }
  .nav-item.active { box-shadow: none; background: var(--ui-muted); }
  .sidebar-footer { display: none; }
  .main-content { padding: 16px; max-width: 100%; }
}
/* ===== 覆盖层语义类（P1 内联样式收口；声明逐字保留以保渲染等价） ===== */
/* 首次使用投票弹窗 */
.vote-mask { position:fixed;inset:0;z-index:999;background:rgba(5,2,20,0.9);backdrop-filter:blur(20px);display:flex;align-items:center;justify-content:center;padding:20px; }
.vote-panel { background:linear-gradient(180deg,var(--color-bg-card),var(--color-bg));border:1px solid var(--tint-3);border-radius:24px;max-width:520px;width:100%;padding:36px 32px;text-align:center;box-shadow:0 20px 60px rgba(0,0,0,0.5); }
.vote-emoji { font-size:48px;margin-bottom:12px; }
.vote-title { font-size:24px;font-weight:700;margin-bottom:8px; }
.vote-desc { font-size:14px;color:var(--md-text-secondary);margin-bottom:20px; }
.vote-hint { font-size:12px;opacity:0.7; }
.vote-grid { display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin-bottom:20px;max-height:320px;overflow-y:auto; }
.vote-item-emoji { font-size:32px;margin-bottom:6px; }
.vote-item-name { font-size:14px;font-weight:600; }
.vote-item-meta { font-size:11px;color:var(--md-text-secondary);margin-top:2px; }
.vote-submit { width:100%;padding:14px;border-radius:14px;border:none;background:linear-gradient(135deg,var(--md-primary),var(--md-secondary));color:#fff;font-size:16px;font-weight:700;cursor:pointer;font-family:inherit; }
.vote-done-emoji { font-size:64px;margin-bottom:16px; }
.vote-done-title { font-size:22px;font-weight:700;margin-bottom:8px; }
.vote-done-desc { font-size:14px;color:var(--md-text-secondary); }
.vote-done-hint { font-size:13px;color:var(--md-gold);margin-top:8px; }

/* 全屏广播 */
.bc-full-mask { position:fixed;inset:0;z-index:900;background:rgba(5,2,20,0.92);backdrop-filter:blur(24px);display:flex;align-items:center;justify-content:center;padding:40px;cursor:pointer; }
.bc-full-inner { max-width:700px;text-align:center; }
.bc-full-emoji { font-size:64px;margin-bottom:20px; }
.bc-full-text { font-size:32px;font-weight:700;color:#fff;line-height:1.4;margin-bottom:12px; }
.bc-full-hint { font-size:14px;color:rgba(255,255,255,0.6); }

/* 弹窗广播 */
.bc-pop-mask { position:fixed;inset:0;z-index:900;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;padding:20px; }
.bc-pop-panel { background:var(--md-surface-2);border:1px solid var(--tint-3);border-radius:24px;padding:32px 28px;max-width:460px;width:100%;box-shadow:0 20px 60px rgba(0,0,0,0.5); }
.bc-pop-emoji { font-size:40px;margin-bottom:12px;text-align:center; }
.bc-pop-text { font-size:20px;font-weight:700;color:var(--color-text);text-align:center;margin-bottom:16px;line-height:1.5; }
.bc-pop-btn { width:100%;padding:12px;border-radius:14px;border:1px solid var(--tint-3);background:var(--tint-2);color:var(--color-text-secondary);font-size:14px;cursor:pointer;font-family:inherit; }

/* 横幅广播 */
.bc-banner { position:fixed;top:0;left:0;right:0;z-index:900;background:linear-gradient(135deg,var(--md-primary),var(--md-secondary));padding:12px 24px;text-align:center;cursor:pointer;box-shadow:0 4px 20px rgba(0,0,0,0.3); }
.bc-banner-text { color:#fff;font-size:16px;font-weight:600; }
.bc-banner-hint { color:rgba(255,255,255,0.5);font-size:12px;margin-left:12px; }

/* 通知 toast */
.notice-toast { position:fixed;top:60px;right:20px;z-index:900;max-width:380px;background:var(--md-surface-2);border:1px solid var(--tint-3);border-radius:16px;padding:16px 20px;box-shadow:0 8px 32px rgba(0,0,0,0.4); }
.notice-body { display:flex;align-items:flex-start;gap:12px; }
.notice-icon { display: inline-flex; align-items: center; justify-content: center; font-size:24px; }
.notice-title { font-size:14px;font-weight:700;color:var(--color-text);margin-bottom:4px; }
.notice-text { font-size:13px;color:var(--color-text-secondary);line-height:1.5; }
.notice-close { width:24px;height:24px;border-radius:50%;border:1px solid var(--tint-3);background:transparent;color:var(--color-text-secondary);cursor:pointer;font-size:12px;display:flex;align-items:center;justify-content:center;flex-shrink:0; }
</style>
