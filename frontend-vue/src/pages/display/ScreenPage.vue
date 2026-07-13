<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet } from '@/utils/api'
import { useDisplaySSE } from '@/composables/useDisplaySSE'
import type { ScoreUpdateData, BroadcastData, NoticeData, PetUpdateData } from '@/composables/useDisplaySSE'

const router = useRouter()

// ===== 认证状态 =====
const token = ref('')
const classInfo = ref<{ id: number; name: string; grade: string; student_count: number } | null>(null)
const authenticated = ref(false)
const loadError = ref('')
const loading = ref(true)

// ===== 大屏数据 =====
interface PetEntry {
  student_id: number; student_no: string; student_name: string; total_score: number
  has_pet: boolean; pet_name?: string; pet_type?: string
  level: number; experience: number; mood: number
  emoji: string; stage_name: string; stage_title: string; exp_max: number
}

interface DisplayData {
  class_name: string; grade: string; student_count: number
  pets: PetEntry[]; recent_scores: any[]; broadcasts: any[]
}

const data = ref<DisplayData | null>(null)
const scoreAnimations = ref<Record<number, { direction: 'up' | 'down'; amount: number }>>({})

// 8x8 矩阵
const gridSlots = computed(() => {
  if (!data.value) return []
  const slots: (PetEntry | null)[] = [...(data.value.pets || [])]
  while (slots.length < 64) slots.push(null)
  return slots.slice(0, 64)
})

// ===== 活跃广播 =====
const activeBroadcast = ref<BroadcastData | null>(null)
const broadcastTimer = ref<ReturnType<typeof setTimeout> | null>(null)

// ===== SSE 连接 =====
const {
  state: sseState,
  scoreUpdates,
  broadcasts,
  notices,
  petUpdates,
  connect: connectSSE,
} = useDisplaySSE()

// Token 过期监听
onMounted(() => {
  window.addEventListener('display:token-expired', () => {
    router.replace({ name: 'display-login' })
  })
})

// ===== 初始化 =====
onMounted(async () => {
  const storedToken = sessionStorage.getItem('display_token')
  const storedInfo = sessionStorage.getItem('display_class_info')

  if (!storedToken || !storedInfo) {
    router.replace({ name: 'display-login' })
    return
  }

  token.value = storedToken
  classInfo.value = JSON.parse(storedInfo)

  // 加载初始数据
  await loadInitialData()
})

async function loadInitialData() {
  loading.value = true
  loadError.value = ''
  try {
    const res = await apiGet<{ data: DisplayData }>('/api/v1/display/initial-data', {
      params: { token: token.value },
    })
    data.value = res.data
    authenticated.value = true
    loading.value = false

    // 加载完成后连接 SSE
    connectSSE(token.value, classInfo.value!.id)
  } catch (e: any) {
    if (e?.response?.status === 401) {
      sessionStorage.removeItem('display_token')
      router.replace({ name: 'display-login' })
      return
    }
    loadError.value = '加载数据失败，请检查网络连接'
    loading.value = false
  }
}

// ===== SSE 事件消费 =====

// 积分变化事件 → 触发动画
watch(scoreUpdates, (events) => {
  for (const ev of events) {
    triggerScoreAnimation(ev)
    updateLocalScore(ev)
  }
}, { deep: true })

// 广播事件 → 显示弹窗
watch(broadcasts, (events) => {
  const latest = events[events.length - 1]
  if (latest) {
    showBroadcast(latest)
  }
}, { deep: true })

// 通知事件 → 简要提示
watch(notices, (events) => {
  // 通知可用作底部滚动条
}, { deep: true })

// 宠物更新事件
watch(petUpdates, (events) => {
  for (const ev of events) {
    updateLocalPet(ev)
  }
}, { deep: true })

// ===== 动画触发 =====
function triggerScoreAnimation(ev: ScoreUpdateData) {
  const direction = ev.amount > 0 ? 'up' : 'down'
  scoreAnimations.value[ev.student_id] = { direction, amount: Math.abs(ev.amount) }
  setTimeout(() => {
    delete scoreAnimations.value[ev.student_id]
  }, 2000)
}

function updateLocalScore(ev: ScoreUpdateData) {
  if (!data.value?.pets) return
  const pet = data.value.pets.find(p => p.student_id === ev.student_id)
  if (pet) {
    pet.total_score = ev.total_score
    if (ev.pet_level !== undefined) pet.level = ev.pet_level
    if (ev.pet_experience !== undefined) pet.experience = ev.pet_experience
    if (ev.pet_mood !== undefined) pet.mood = ev.pet_mood
  }
}

function updateLocalPet(ev: PetUpdateData) {
  if (!data.value?.pets) return
  const pet = data.value.pets.find(p => p.student_id === ev.student_id)
  if (pet) {
    if (ev.mood !== undefined) pet.mood = ev.mood
    if (ev.level !== undefined) pet.level = ev.level
    if (ev.experience !== undefined) pet.experience = ev.experience
  }
}

// ===== 广播显示 =====
function showBroadcast(msg: BroadcastData) {
  if (broadcastTimer.value) clearTimeout(broadcastTimer.value)
  activeBroadcast.value = msg
  const duration = Math.max(3000, (msg.display_seconds || 8) * 1000)
  broadcastTimer.value = setTimeout(() => {
    activeBroadcast.value = null
  }, duration)
}

// ===== 退出（防误触：连续点 3 次才弹出确认） =====
const showExitConfirm = ref(false)
const exitTapCount = ref(0)
let exitTapTimer: ReturnType<typeof setTimeout> | null = null

function handleExitTap() {
  exitTapCount.value++
  if (exitTapCount.value >= 3) {
    showExitConfirm.value = true
    exitTapCount.value = 0
    if (exitTapTimer) clearTimeout(exitTapTimer)
  } else {
    if (exitTapTimer) clearTimeout(exitTapTimer)
    exitTapTimer = setTimeout(() => { exitTapCount.value = 0 }, 2000)
  }
}

function confirmExit() {
  showExitConfirm.value = false
  sessionStorage.removeItem('display_token')
  sessionStorage.removeItem('display_class_info')
  router.replace({ name: 'display-login' })
}

function cancelExit() {
  showExitConfirm.value = false
  exitTapCount.value = 0
}
</script>

<template>
  <div class="display-screen"
    :class="{ 'has-broadcast': activeBroadcast }"
  >
    <!-- ===== 加载中 ===== -->
    <div v-if="loading" class="screen-loading">
      <div class="loading-spinner">🌌</div>
      <p>正在连接班级大屏...</p>
    </div>

    <!-- ===== 加载错误 ===== -->
    <div v-else-if="loadError" class="screen-error">
      <div class="error-icon">⚠️</div>
      <p>{{ loadError }}</p>
      <button class="retry-btn" @click="loadInitialData">重试</button>
      <button class="exit-btn" @click="confirmExit">返回登录</button>
    </div>

    <!-- ===== 大屏主体 ===== -->
    <template v-else-if="data">
      <!-- 广播覆盖层 -->
      <Transition name="broadcast-pop">
        <div v-if="activeBroadcast" class="broadcast-overlay"
          :class="activeBroadcast.type"
          @click="activeBroadcast = null"
        >
          <div class="broadcast-content" @click.stop>
            <div class="broadcast-icon">
              {{ activeBroadcast.type === 'fullscreen' ? '📢' : activeBroadcast.type === 'banner' ? '💬' : '🔔' }}
            </div>
            <div class="broadcast-text">{{ activeBroadcast.content }}</div>
            <div class="broadcast-timer">
              <div class="timer-bar" :style="{ animationDuration: (activeBroadcast.display_seconds || 8) + 's' }"></div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- SSE 连接状态指示 -->
      <div class="connection-status" :class="{ connected: sseState.connected, polling: sseState.polling }">
        <span class="status-dot"></span>
        <span class="status-text">{{ sseState.connected ? '实时连接' : sseState.polling ? '轮询中' : '重连中...' }}</span>
      </div>

      <!-- 头部 -->
      <div class="screen-header">
        <div class="header-left">
          <h1 class="class-name">{{ data.class_name }}</h1>
          <span class="class-meta">{{ data.grade }} · {{ data.student_count }}人</span>
        </div>
        <button class="exit-screen-btn" @click="handleExitTap" title="连续点击3次退出(防误触)">
          ✕
        </button>
      </div>

      <!-- 退出确认弹窗 -->
      <Transition name="fade">
        <div v-if="showExitConfirm" class="exit-overlay" @click.self="cancelExit">
          <div class="exit-modal">
            <div class="exit-modal-icon">🚪</div>
            <h3 class="exit-modal-title">退出大屏？</h3>
            <p class="exit-modal-desc">退出后需要重新输入班级码</p>
            <div class="exit-modal-actions">
              <button class="exit-modal-btn exit-modal-btn--cancel" @click="cancelExit">取消</button>
              <button class="exit-modal-btn exit-modal-btn--confirm" @click="confirmExit">确认退出</button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 8x8 宠物矩阵 -->
      <div v-if="data.pets.length === 0" class="empty-state">
        <div class="empty-icon">📭</div>
        <p>暂无学生数据</p>
      </div>

      <div v-else class="pet-grid">
        <div
          v-for="(slot, idx) in gridSlots"
          :key="idx"
          class="grid-cell"
          :class="{
            empty: !slot,
            'score-up': slot && scoreAnimations[slot.student_id]?.direction === 'up',
            'score-down': slot && scoreAnimations[slot.student_id]?.direction === 'down',
            'has-pet': slot?.has_pet,
          }"
        >
          <template v-if="slot">
            <!-- 宠物 emoji -->
            <div
              class="pet-display"
              :class="{
                bounce: scoreAnimations[slot.student_id]?.direction === 'up',
                shake: scoreAnimations[slot.student_id]?.direction === 'down',
                energetic: !scoreAnimations[slot.student_id] && slot.mood > 50,
                normal: !scoreAnimations[slot.student_id] && slot.mood >= 30 && slot.mood <= 50,
                weak: !scoreAnimations[slot.student_id] && slot.mood < 30,
              }"
            >
              <span class="pet-emoji">{{ slot.has_pet ? slot.emoji : '🥚' }}</span>

              <!-- 积分变化浮动数字 -->
              <Transition name="float-up">
                <span v-if="scoreAnimations[slot.student_id]" class="score-float" :class="scoreAnimations[slot.student_id].direction">
                  {{ scoreAnimations[slot.student_id].direction === 'up' ? '+' : '-' }}{{ scoreAnimations[slot.student_id].amount }}
                </span>
              </Transition>
            </div>

            <!-- 经验值条 -->
            <div class="exp-bar">
              <div class="exp-fill" :style="{ width: Math.min(100, (slot.experience / Math.max(1, slot.exp_max)) * 100) + '%' }"></div>
            </div>

            <!-- 信息 -->
            <div class="cell-info">
              <span class="cell-name">{{ slot.student_name }}</span>
              <span class="cell-no">{{ slot.student_no || '--' }}</span>
              <span class="cell-score" :class="{ positive: slot.total_score >= 0, negative: slot.total_score < 0 }">
                {{ slot.total_score }}分
              </span>
            </div>

            <!-- 宠物等级标签 -->
            <div class="level-badge" v-if="slot.has_pet">Lv.{{ slot.level }}</div>
          </template>

          <template v-else>
            <div class="empty-slot"></div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.display-screen {
  min-height: 100vh;
  background: linear-gradient(135deg, #0c0a20 0%, #1a1040 30%, #0d1b2a 70%, #0a1628 100%);
  color: #e8e6f0;
  font-family: "PingFang SC", "Noto Sans SC", system-ui, sans-serif;
  padding: 12px 16px 20px;
  position: relative;
  overflow: hidden;
  user-select: none;
}

/* ===== 背景装饰 ===== */
.display-screen::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse at 20% 20%, rgba(120, 80, 255, 0.06) 0%, transparent 60%),
    radial-gradient(ellipse at 80% 80%, rgba(30, 140, 220, 0.04) 0%, transparent 60%);
  pointer-events: none;
}

/* ===== 加载 & 错误 ===== */
.screen-loading,
.screen-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  gap: 12px;
  color: rgba(200, 190, 240, 0.6);
}

.loading-spinner {
  font-size: 48px;
  animation: spin 2s linear infinite;
}

.error-icon { font-size: 48px; }
.retry-btn, .exit-btn {
  padding: 8px 24px;
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.04);
  color: rgba(200,190,240,0.7);
  cursor: pointer;
  font-size: 14px;
}
.retry-btn:hover, .exit-btn:hover {
  background: rgba(255,255,255,0.1);
}

/* ===== 连接状态 ===== */
.connection-status {
  position: fixed;
  top: 10px;
  right: 50px;
  z-index: 50;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  color: rgba(200, 190, 240, 0.3);
  padding: 3px 10px;
  border-radius: 20px;
  background: rgba(0,0,0,0.3);
}
.connection-status.connected .status-dot {
  background: #4ade80;
  box-shadow: 0 0 6px rgba(74,222,128,0.5);
}
.connection-status.polling .status-dot {
  background: #f59e0b;
  box-shadow: 0 0 6px rgba(245,158,11,0.5);
}
.status-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  background: #64748b;
  transition: all 0.3s;
}

/* ===== 头部 ===== */
.screen-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  margin-bottom: 10px;
}
.header-left { display: flex; align-items: baseline; gap: 12px; }
.class-name { font-size: 20px; font-weight: 700; margin: 0; letter-spacing: 0.03em; color: #f0ecff; }
.class-meta { font-size: 13px; color: rgba(200, 190, 240, 0.6); }

.exit-screen-btn {
  background: none;
  border: 1px solid rgba(255,255,255,0.08);
  color: rgba(200,190,240,0.3);
  width: 32px; height: 32px;
  border-radius: 10px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}
.exit-screen-btn:hover { color: #f87171; border-color: rgba(248,113,113,0.3); }

/* ===== 空状态 ===== */
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: rgba(200,190,240,0.5);
}
.empty-icon { font-size: 48px; margin-bottom: 12px; }

/* ===== 8x8 网格 ===== */
.pet-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 6px;
  max-width: 1100px;
  margin: 0 auto;
}

.grid-cell {
  aspect-ratio: 1;
  border-radius: 14px;
  background: rgba(255,255,255,0.025);
  border: 1px solid rgba(255,255,255,0.04);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4px;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s, border-color 0.2s;
}

.grid-cell:not(.empty):hover {
  transform: scale(1.03);
  border-color: rgba(255,255,255,0.1);
}

.grid-cell.score-up {
  border-color: rgba(74, 222, 128, 0.25);
  box-shadow: 0 0 16px rgba(74, 222, 128, 0.08);
}

.grid-cell.score-down {
  border-color: rgba(248, 113, 113, 0.25);
  box-shadow: 0 0 16px rgba(248, 113, 113, 0.08);
}

.grid-cell.empty {
  background: transparent;
  border-style: dashed;
  border-color: rgba(255,255,255,0.03);
}

.empty-slot { width: 100%; height: 100%; }

/* ===== 宠物显示 ===== */
.pet-display {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 4px;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.pet-emoji {
  font-size: 28px;
  line-height: 1;
  filter: drop-shadow(0 0 6px rgba(180, 140, 255, 0.3));
  transition: all 0.3s;
}

/* 动态状态 */
.pet-display.energetic .pet-emoji {
  animation: float 2s ease-in-out infinite;
  filter: drop-shadow(0 0 10px rgba(74, 222, 128, 0.4));
}

.pet-display.normal .pet-emoji {
  animation: breathe-pet 3s ease-in-out infinite;
}

.pet-display.weak .pet-emoji {
  opacity: 0.5;
  filter: grayscale(0.6) drop-shadow(0 0 3px rgba(100, 100, 100, 0.2));
  animation: weak-shake 4s ease-in-out infinite;
}

/* 积分变化动画 */
.pet-display.bounce .pet-emoji {
  animation: bounceAnim 0.6s cubic-bezier(0.28, 1.33, 0.64, 1) 2;
}

.pet-display.shake .pet-emoji {
  animation: shakeAnim 0.5s ease-in-out 3;
}

/* 浮动数字 */
.score-float {
  position: absolute;
  top: -8px;
  right: -12px;
  font-size: 11px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 6px;
  pointer-events: none;
}
.score-float.up {
  color: #4ade80;
  background: rgba(74,222,128,0.15);
}
.score-float.down {
  color: #f87171;
  background: rgba(248,113,113,0.15);
}

/* ===== 等级标签 ===== */
.level-badge {
  position: absolute;
  top: 3px;
  left: 3px;
  font-size: 7px;
  font-weight: 700;
  color: rgba(200, 190, 240, 0.5);
  background: rgba(255,255,255,0.04);
  padding: 1px 4px;
  border-radius: 4px;
}

/* ===== 经验条 ===== */
.exp-bar {
  width: 80%;
  height: 3px;
  background: rgba(255,255,255,0.06);
  border-radius: 2px;
  margin: 3px 0;
  overflow: hidden;
}

.exp-fill {
  height: 100%;
  background: linear-gradient(90deg, #7c3aed, #a78bfa);
  border-radius: 2px;
  transition: width 0.6s ease;
}

/* ===== 格子底部信息 ===== */
.cell-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0;
  margin-top: 1px;
}

.cell-name {
  font-size: 10px;
  font-weight: 600;
  color: rgba(220, 210, 250, 0.8);
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.2;
}

.cell-no {
  font-size: 8px;
  color: rgba(200, 190, 240, 0.4);
  font-weight: 500;
}

.cell-score {
  font-size: 9px;
  font-weight: 700;
  margin-top: 1px;
}
.cell-score.positive { color: rgba(74, 222, 128, 0.7); }
.cell-score.negative { color: rgba(248, 113, 113, 0.7); }

/* ===== 广播覆盖层 ===== */
.broadcast-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
}

.broadcast-overlay.banner {
  align-items: flex-start;
  padding-top: 40px;
  background: rgba(10, 5, 30, 0.6);
  backdrop-filter: blur(8px);
}

.broadcast-overlay.popup {
  background: rgba(10, 5, 30, 0.75);
  backdrop-filter: blur(12px);
}

.broadcast-overlay.fullscreen {
  background: rgba(10, 5, 30, 0.92);
  backdrop-filter: blur(20px);
}

.broadcast-content {
  text-align: center;
  max-width: 600px;
  padding: 40px;
  animation: messageIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.broadcast-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.broadcast-text {
  font-size: 26px;
  font-weight: 700;
  line-height: 1.5;
  color: #f0ecff;
}

.broadcast-timer {
  margin-top: 20px;
  width: 200px;
  height: 3px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
  margin-left: auto;
  margin-right: auto;
  overflow: hidden;
}

.timer-bar {
  height: 100%;
  background: linear-gradient(90deg, #7c3aed, #a78bfa);
  border-radius: 2px;
  animation: timerShrink linear forwards;
}

@keyframes timerShrink {
  from { width: 100%; }
  to { width: 0%; }
}

/* ===== 动画关键帧 ===== */
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

@keyframes breathe-pet {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.05); opacity: 0.85; }
}

@keyframes weak-shake {
  0%, 100% { transform: translateX(0); }
  10% { transform: translateX(-1px); }
  20% { transform: translateX(1px); }
  30% { transform: translateX(0); }
}

@keyframes bounceAnim {
  0%, 100% { transform: translateY(0) scale(1); }
  30% { transform: translateY(-8px) scale(1.15); }
  60% { transform: translateY(-3px) scale(1.05); }
}

@keyframes shakeAnim {
  0%, 100% { transform: translateX(0); }
  20% { transform: translateX(-3px); }
  40% { transform: translateX(3px); }
  60% { transform: translateX(-2px); }
  80% { transform: translateX(2px); }
}

@keyframes messageIn {
  from { opacity: 0; transform: scale(0.9) translateY(20px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ===== 过渡 ===== */
.broadcast-pop-enter-active { animation: messageIn 0.4s ease-out; }
.broadcast-pop-leave-active { animation: messageIn 0.3s ease-in reverse; }

.float-up-enter-active { transition: all 0.3s ease-out; }
.float-up-leave-active { transition: all 0.2s ease-in; }
.float-up-enter-from { opacity: 0; transform: translateY(8px); }
.float-up-leave-to { opacity: 0; transform: translateY(-12px); }

/* ===== 退出确认弹窗 ===== */
.exit-overlay {
  position: fixed; inset: 0; z-index: 200;
  background: rgba(5,2,20,.85); backdrop-filter: blur(16px);
  display: flex; align-items: center; justify-content: center;
}
.exit-modal {
  background: linear-gradient(135deg,#1a1040,#0d1b2a);
  border: 1px solid rgba(255,255,255,.08); border-radius: 20px;
  padding: 36px 40px; text-align: center; max-width: 320px; width: 85%;
}
.exit-modal-icon { font-size: 40px; margin-bottom: 12px; }
.exit-modal-title { font-size: 18px; font-weight: 700; color: #e8e0f8; margin: 0 0 8px; }
.exit-modal-desc { font-size: 13px; color: rgba(200,190,240,.5); margin: 0 0 24px; }
.exit-modal-actions { display: flex; gap: 10px; }
.exit-modal-btn { flex:1; padding:12px; border-radius:12px; border:none; font-size:14px; font-weight:600; cursor:pointer; transition:all .15s; font-family:inherit; }
.exit-modal-btn--cancel { background:rgba(255,255,255,.06); color:rgba(200,190,240,.7); }
.exit-modal-btn--cancel:hover { background:rgba(255,255,255,.1); }
.exit-modal-btn--confirm { background:linear-gradient(135deg,#ef4444,#dc2626); color:#fff; }
.exit-modal-btn--confirm:hover { background:linear-gradient(135deg,#f87171,#ef4444); }
.fade-enter-active,.fade-leave-active { transition:opacity .25s ease; }
.fade-enter-from,.fade-leave-to { opacity:0; }
</style>
