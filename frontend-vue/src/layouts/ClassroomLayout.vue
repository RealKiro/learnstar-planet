<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { apiPost, apiGet } from '@/utils/api'
import { getAllSeries, getSeriesName } from '@/utils/petData'

const router = useRouter()
const route = useRoute()

const classInfo = ref<{ id: number; name: string; student_count?: number } | null>(null)
const showVoteModal = ref(true)
const voteSeries = ref('myth')
const voting = ref(false)
const voteDone = ref(false)

const activeNav = computed(() => String(route.name))

const navItems = [
  { page: 'classroom-overview', label: '班级总览', icon: '🏠' },
  { page: 'classroom-scores', label: '课堂评价', icon: '✏️' },
  { page: 'classroom-pk', label: '年级战场', icon: '🏆' },
  { page: 'classroom-pokedex', label: '宠物图鉴', icon: '📚' },
]

const allSeries = getAllSeries()
const lastEventId = ref(0)

function navigate(name: string) { router.push({ name }) }
function goToLogin() { sessionStorage.clear(); router.push({ name: 'login', query: { mode: 'code' } }) }

// ===== 广播/通知接收（SSE + 轮询降级） =====
const currentBroadcast = ref<{
  id: number; type: string; content: string; display_seconds: number; created_at: string
} | null>(null)

const currentNotice = ref<{
  id: number; title: string; content: string; type: string; published_at: string
} | null>(null)

let pollTimer: ReturnType<typeof setInterval> | null = null
    startPolling()
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

  // 使用轮询接收广播和通知（php artisan serve 不支持 SSE 长连接）
  lastEventId.value = parseInt(sessionStorage.getItem('last_event_id') || '0', 10)
  startPolling()
})

onUnmounted(() => {
  
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  sessionStorage.setItem('last_event_id', String(lastEventId.value))
})
</script>

<template>
  <div class="app-shell">
    <!-- 首次使用投票弹窗 -->
    <Transition name="fade">
      <div v-if="showVoteModal" @click.self="() => {}"
        style="position:fixed;inset:0;z-index:999;background:rgba(5,2,20,0.9);backdrop-filter:blur(20px);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:linear-gradient(180deg,#1a1040,#0d1b2a);border:1px solid rgba(255,255,255,0.08);border-radius:24px;max-width:520px;width:100%;padding:36px 32px;text-align:center;box-shadow:0 20px 60px rgba(0,0,0,0.5);">
          <div v-if="!voteDone">
            <div style="font-size:48px;margin-bottom:12px;">🎉</div>
            <h2 style="font-size:24px;font-weight:700;margin-bottom:8px;">欢迎来到学趣星球！</h2>
            <p style="font-size:14px;color:var(--md-text-secondary);margin-bottom:20px;">
              请全班投票选择你们喜欢的宠物类别<br>
              <span style="font-size:12px;opacity:0.7;">选定后每人可免费选择一只心仪的宠物</span>
            </p>
            <div style="display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin-bottom:20px;max-height:320px;overflow-y:auto;">
              <button v-for="s in allSeries" :key="s.id" @click="voteSeries = s.id"
                :style="{
                  padding:'16px 12px', borderRadius:'16px', cursor:'pointer', transition:'0.2s', fontFamily:'inherit',
                  border: voteSeries === s.id ? '2px solid var(--md-primary)' : '1px solid rgba(255,255,255,0.06)',
                  background: voteSeries === s.id ? 'rgba(167,139,250,0.1)' : 'rgba(255,255,255,0.02)',
                  color: voteSeries === s.id ? 'var(--md-primary-light)' : 'var(--md-text)',
                }">
                <div style="font-size:32px;margin-bottom:6px;">{{ s.emoji }}</div>
                <div style="font-size:14px;font-weight:600;">{{ s.name }}</div>
                <div style="font-size:11px;color:var(--md-text-secondary);margin-top:2px;">{{ s.species.length }}种宠物</div>
              </button>
            </div>
            <button @click="confirmVote" :disabled="voting"
              style="width:100%;padding:14px;border-radius:14px;border:none;background:linear-gradient(135deg,var(--md-primary),var(--md-secondary));color:#fff;font-size:16px;font-weight:700;cursor:pointer;font-family:inherit;">
              {{ voting ? '投票中...' : '✅ 选择「' + getSeriesName(voteSeries) + '」系列' }}
            </button>
          </div>
          <div v-else>
            <div style="font-size:64px;margin-bottom:16px;">🎊</div>
            <h2 style="font-size:22px;font-weight:700;margin-bottom:8px;">选择成功！</h2>
            <p style="font-size:14px;color:var(--md-text-secondary);">
              已选定「{{ getSeriesName(voteSeries) }}」系列
            </p>
            <p style="font-size:13px;color:var(--md-gold);margin-top:8px;">
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
        @click="dismissBroadcast"
        style="position:fixed;inset:0;z-index:900;background:rgba(5,2,20,0.92);backdrop-filter:blur(24px);display:flex;align-items:center;justify-content:center;padding:40px;cursor:pointer;">
        <div style="max-width:700px;text-align:center;">
          <div style="font-size:64px;margin-bottom:20px;">📡</div>
          <div style="font-size:32px;font-weight:700;color:#fff;line-height:1.4;margin-bottom:12px;">{{ currentBroadcast.content }}</div>
          <div style="font-size:14px;color:var(--md-text-secondary);">点击任意位置关闭</div>
        </div>
      </div>
    </Transition>

    <!-- 弹窗广播 -->
    <Transition name="pop">
      <div v-if="currentBroadcast && currentBroadcast.type === 'popup'"
        @click.self="dismissBroadcast"
        style="position:fixed;inset:0;z-index:900;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:var(--md-surface-2);border:1px solid rgba(255,255,255,0.08);border-radius:24px;padding:32px 28px;max-width:460px;width:100%;box-shadow:0 20px 60px rgba(0,0,0,0.5);">
          <div style="font-size:40px;margin-bottom:12px;text-align:center;">📢</div>
          <div style="font-size:20px;font-weight:700;color:#fff;text-align:center;margin-bottom:16px;line-height:1.5;">{{ currentBroadcast.content }}</div>
          <button @click="dismissBroadcast"
            style="width:100%;padding:12px;border-radius:14px;border:1px solid rgba(255,255,255,0.06);background:rgba(255,255,255,0.05);color:var(--md-text-secondary);font-size:14px;cursor:pointer;font-family:inherit;">
            我知道了
          </button>
        </div>
      </div>
    </Transition>

    <!-- 横幅广播 -->
    <Transition name="slide-down">
      <div v-if="currentBroadcast && currentBroadcast.type === 'banner'"
        @click="dismissBroadcast"
        style="position:fixed;top:0;left:0;right:0;z-index:900;background:linear-gradient(135deg,var(--md-primary),var(--md-secondary));padding:12px 24px;text-align:center;cursor:pointer;box-shadow:0 4px 20px rgba(0,0,0,0.3);">
        <span style="color:#fff;font-size:16px;font-weight:600;">📢 {{ currentBroadcast.content }}</span>
        <span style="color:rgba(255,255,255,0.5);font-size:12px;margin-left:12px;">点击关闭</span>
      </div>
    </Transition>

    <!-- 通知提示 -->
    <Transition name="slide-down">
      <div v-if="currentNotice"
        style="position:fixed;top:60px;right:20px;z-index:900;max-width:380px;background:var(--md-surface-2);border:1px solid rgba(255,255,255,0.08);border-radius:16px;padding:16px 20px;box-shadow:0 8px 32px rgba(0,0,0,0.4);">
        <div style="display:flex;align-items:flex-start;gap:12px;">
          <span style="font-size:24px;">📋</span>
          <div style="flex:1;">
            <div style="font-size:14px;font-weight:700;color:#fff;margin-bottom:4px;">{{ currentNotice.title }}</div>
            <div style="font-size:13px;color:var(--md-text-secondary);line-height:1.5;">{{ currentNotice.content }}</div>
          </div>
          <button @click="dismissNotice"
            style="width:24px;height:24px;border-radius:50%;border:1px solid rgba(255,255,255,0.06);background:transparent;color:rgba(255,255,255,0.3);cursor:pointer;font-size:12px;display:flex;align-items:center;justify-content:center;flex-shrink:0;">✕</button>
        </div>
      </div>
    </Transition>

    <nav class="sidebar">
      <div class="logo">
        <div class="brand"><span>🌌</span> 学趣星球</div>
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
  width: var(--md-sidebar-width); background: var(--md-surface-2);
  border-right: 1px solid rgba(255,255,255,0.05);
  padding: 24px 16px 16px; display: flex; flex-direction: column;
  position: sticky; top: 0; height: 100vh; flex-shrink: 0;
  backdrop-filter: blur(12px); box-shadow: 4px 0 20px rgba(0,0,0,0.3); z-index: 10;
}
.logo { font-size: 22px; font-weight: 700; padding: 8px 12px; margin-bottom: 12px; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 12px; }
.brand { display: flex; align-items: center; gap: 10px; background: linear-gradient(135deg,var(--md-primary),var(--md-secondary)); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.brand span { font-size: 28px; -webkit-text-fill-color: initial; }
.class-badge { text-align: center; padding: 12px; margin-bottom: 12px; background: rgba(255,255,255,0.03); border-radius: var(--md-radius); border: 1px solid rgba(255,255,255,0.04); }
.class-name { font-size: 16px; font-weight: 700; display: block; }
.class-count { font-size: 12px; color: var(--md-text-secondary); }
.nav-list { display: flex; flex-direction: column; gap: 6px; flex: 1; }
.nav-item { display: flex; align-items: center; gap: 14px; padding: 12px 16px; border-radius: var(--md-radius); cursor: pointer; transition: 0.2s; color: var(--md-text-secondary); border: none; background: transparent; width: 100%; font-size: 16px; font-weight: 500; font-family: inherit; text-align: left; }
.nav-item:hover { background: rgba(255,255,255,0.05); color: #fff; }
.nav-item.active { background: rgba(167,139,250,0.15); color: var(--md-primary-light); box-shadow: inset 3px 0 0 var(--md-primary); }
.nav-item .icon { font-size: 22px; width: 28px; text-align: center; }

.sidebar-footer { border-top: 1px solid rgba(255,255,255,0.05); padding-top: 12px; }
.exit-btn { width: 100%; padding: 10px; border-radius: var(--md-radius); border: 1px solid rgba(255,100,100,0.15); background: rgba(255,100,100,0.08); color: #fca5a5; font-size: 14px; font-weight: 500; cursor: pointer; transition: 0.2s; font-family: inherit; }
.exit-btn:hover { background: rgba(255,100,100,0.15); }

.main-content { flex: 1; padding: 28px 32px 40px; max-width: calc(100% - var(--md-sidebar-width)); overflow-x: hidden; }

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
  .sidebar { width: 100%; height: auto; position: sticky; flex-direction: row; flex-wrap: wrap; padding: 12px 16px; border-right: none; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .logo { margin-bottom: 0; border-bottom: none; padding-bottom: 0; font-size: 18px; flex: 1; }
  .class-badge { display: none; }
  .nav-list { flex-direction: row; gap: 4px; flex: 2; justify-content: flex-end; }
  .nav-item { padding: 8px 12px; font-size: 14px; }
  .nav-item .icon { font-size: 18px; width: 24px; }
  .sidebar-footer { display: none; }
  .main-content { padding: 16px; max-width: 100%; }
}
</style>
