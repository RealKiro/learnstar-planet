<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'

interface PetEntry {
  student_id: number; student_no: string; student_name: string; total_score: number
  has_pet: boolean; pet_name?: string; pet_type?: string; level: number
  experience: number; mood: number; emoji: string; stage_name: string
}
interface DisplayData {
  class_name: string; grade: string; student_count: number
  pets: PetEntry[]; broadcasts: any[]; notices: any[]; recent_scores: any[]
}
interface Message {
  id: number; content: string; type: string; display_seconds: number; created_at: string
}

const toast = useToastStore()
const data = ref<DisplayData | null>(null)
const loading = ref(true)
const currentMessage = ref<Message | null>(null)
const messageTimer = ref<NodeJS.Timeout | null>(null)
const showUnlock = ref(false)
const unlockPassword = ref('')
const unlockError = ref('')
const scoreFlash = ref<Record<number, 'up' | 'down' | null>>({})

const unlockStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
let pollTimer: NodeJS.Timeout | null = null
let lastScoreMap: Record<number, number> = {}

// 8x8 matrix padding
const gridSlots = computed(() => {
  if (!data.value) return []
  const slots: (PetEntry | null)[] = [...data.value.pets]
  while (slots.length < 64) slots.push(null)
  return slots.slice(0, 64)
})

async function loadDisplay() {
  try {
    const res = await apiGet<{ data: DisplayData }>('/api/v1/teacher/classroom/display')
    const prev = data.value
    data.value = res.data

    // Detect score changes for animation
    if (prev && res.data.pets) {
      const newMap: Record<number, number> = {}
      for (const p of res.data.pets) newMap[p.student_id] = p.total_score
      for (const pid in lastScoreMap) {
        const old = lastScoreMap[pid]
        const nu = newMap[Number(pid)]
        if (nu !== undefined && nu !== old) {
          scoreFlash.value[Number(pid)] = nu > old ? 'up' : 'down'
          setTimeout(() => { scoreFlash.value[Number(pid)] = null }, 2000)
        }
      }
      lastScoreMap = newMap
    } else if (res.data.pets) {
      for (const p of res.data.pets) lastScoreMap[p.student_id] = p.total_score
    }
  } catch (e: any) {
    if (e?.response?.status === 400) {
      toast.show('请先在教师管理中指定班级', 'info')
    }
  } finally {
    loading.value = false
  }
}

async function pollMessages() {
  try {
    const res = await apiGet<{ data: { broadcasts: any[]; notices: any[] } }>('/api/v1/teacher/classroom/messages')
    const msgs = [...(res.data?.broadcasts || []), ...(res.data?.notices || [])]
    if (msgs.length > 0) {
      const msg = msgs[0]
      showMessage({ id: msg.id, content: msg.content || msg.title, type: msg.type, display_seconds: msg.display_seconds || 8, created_at: msg.created_at || msg.published_at })
    }
  } catch {}
}

function showMessage(msg: Message) {
  currentMessage.value = msg
  if (messageTimer.value) clearTimeout(messageTimer.value)
  messageTimer.value = setTimeout(() => {
    currentMessage.value = null
  }, Math.max(3000, (msg.display_seconds || 8) * 1000))
}

function startUnlock() {
  showUnlock.value = true
  unlockPassword.value = ''
  unlockError.value = ''
}
async function submitUnlock() {
  unlockStatus.value = 'loading'
  try {
    await apiPost('/api/v1/teacher/mode', { mode: 'teacher_manage', password: unlockPassword.value })
    unlockStatus.value = 'success'
    setTimeout(() => {
      showUnlock.value = false
      window.location.reload()
    }, 800)
  } catch {
    unlockStatus.value = 'error'
    unlockError.value = '密码错误，请重试'
    setTimeout(() => { unlockStatus.value = 'idle' }, 3000)
  }
}

onMounted(() => {
  loadDisplay()
  pollTimer = setInterval(() => {
    loadDisplay()
    pollMessages()
  }, 5000)
})
onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
  if (messageTimer.value) clearTimeout(messageTimer.value)
})
</script>

<template>
  <div class="classroom-display">
    <!-- Message overlay -->
    <Transition name="message-pop">
      <div v-if="currentMessage" class="message-overlay" :class="currentMessage.type">
        <div class="message-content">
          <div class="message-icon">
            {{ currentMessage.type === 'urgent' ? '⚠️' : currentMessage.type === 'fullscreen' ? '📢' : '💬' }}
          </div>
          <div class="message-text">{{ currentMessage.content }}</div>
          <div class="message-timer">
            <div class="timer-bar" :style="{ animationDuration: (currentMessage.display_seconds || 8) + 's' }"></div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Lock overlay (double-tap to unlock) -->
    <Transition name="fade">
      <div v-if="showUnlock" class="unlock-overlay" @click.self="showUnlock = false">
        <div class="unlock-modal">
          <div class="unlock-icon">🔒</div>
          <h3>解锁教师管理模式</h3>
          <input v-model="unlockPassword" type="password" class="unlock-input" placeholder="请输入登录密码" @keyup.enter="submitUnlock" />
          <div v-if="unlockError" class="unlock-error">{{ unlockError }}</div>
          <div class="unlock-buttons">
            <button class="unlock-btn cancel" @click="showUnlock = false">取消</button>
            <button class="unlock-btn confirm" :style="{ background: unlockStatus === 'loading' ? '#f59e0b' : unlockStatus === 'success' ? '#10b981' : unlockStatus === 'error' ? '#ef4444' : '#7c3aed' }" :disabled="unlockStatus === 'loading'" @click="submitUnlock">
              <template v-if="unlockStatus === 'idle'">确认解锁</template>
              <template v-else-if="unlockStatus === 'loading'">解锁中...</template>
              <template v-else-if="unlockStatus === 'success'">已解锁 ✓</template>
              <template v-else>解锁失败 ✗</template>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Header -->
    <div class="classroom-header">
      <div class="header-left">
        <h1 class="class-name">{{ data?.class_name || '--' }}</h1>
        <span class="class-meta">{{ data?.grade }} · {{ data?.student_count }}人</span>
      </div>
      <div class="header-right">
        <div class="mode-badge">🔒 班级大屏</div>
        <button class="unlock-trigger" @dblclick="startUnlock" title="双击解锁">🔑</button>
      </div>
    </div>

    <!-- Grid -->
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="!data || data.pets.length === 0" class="empty">
      <div class="empty-icon">📭</div>
      <p>该班级暂无学生</p>
    </div>
    <div v-else class="pet-grid">
      <div
        v-for="(slot, idx) in gridSlots"
        :key="idx"
        class="grid-cell"
        :class="{ empty: !slot, flash: slot && scoreFlash[slot.student_id] }"
      >
        <template v-if="slot">
          <div class="pet-emoji" :class="{ bounce: scoreFlash[slot.student_id] === 'up', shake: scoreFlash[slot.student_id] === 'down' }">
            {{ slot.emoji }}
          </div>
          <div class="pet-exp">
            <div class="exp-fill" :style="{ width: Math.min(100, (slot.experience / ((slot.level + 1) * 100)) * 100) + '%' }"></div>
          </div>
          <div class="cell-info">
            <span class="cell-name">{{ slot.student_name }}</span>
            <span class="cell-no">{{ slot.student_no }}</span>
          </div>
        </template>
        <template v-else>
          <div class="empty-slot"></div>
        </template>
      </div>
    </div>

    <!-- Recent scores ticker -->
    <div v-if="data?.recent_scores?.length" class="score-ticker">
      <div class="ticker-label">实时加分</div>
      <div class="ticker-items">
        <span v-for="s in data.recent_scores.slice(0, 5)" :key="s.student_name + s.time" class="ticker-item">
          {{ s.student_name }} <strong :class="s.amount > 0 ? 'positive' : 'negative'">{{ s.amount > 0 ? '+' + s.amount : s.amount }}</strong>
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.classroom-display {
  min-height: 100vh;
  background: linear-gradient(135deg, #0c0a20 0%, #1a1040 30%, #0d1b2a 70%, #0a1628 100%);
  color: #e8e6f0;
  font-family: "PingFang SC", "Noto Sans SC", system-ui, sans-serif;
  padding: 12px 16px 20px;
  position: relative;
  overflow: hidden;
  user-select: none;
}
.classroom-display::before {
  content: ''; position: absolute; inset: 0;
  background: radial-gradient(ellipse at 20% 20%, rgba(120, 80, 255, 0.06) 0%, transparent 60%),
              radial-gradient(ellipse at 80% 80%, rgba(30, 140, 220, 0.04) 0%, transparent 60%);
  pointer-events: none;
}

/* Header */
.classroom-header {
  display: flex; justify-content: space-between; align-items: center;
  padding-bottom: 8px; border-bottom: 1px solid rgba(255,255,255,0.06); margin-bottom: 10px;
}
.header-left { display: flex; align-items: baseline; gap: 12px; }
.class-name { font-size: 20px; font-weight: 700; margin: 0; letter-spacing: 0.03em; color: #f0ecff; }
.class-meta { font-size: 13px; color: rgba(200, 190, 240, 0.6); }
.header-right { display: flex; align-items: center; gap: 10px; }
.mode-badge {
  font-size: 11px; padding: 3px 10px; border-radius: 20px;
  background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.08);
  color: rgba(200,190,240,0.7); font-weight: 500;
}
.unlock-trigger {
  background: none; border: 1px solid rgba(255,255,255,0.08); color: rgba(200,190,240,0.5);
  width: 32px; height: 32px; border-radius: 10px; font-size: 16px;
  cursor: pointer; display: flex; align-items: center; justify-content: center;
  transition: all 0.2s; opacity: 0.4;
}
.unlock-trigger:hover { opacity: 1; border-color: rgba(255,255,255,0.2); color: #fff; }

/* Grid */
.pet-grid {
  display: grid; grid-template-columns: repeat(8, 1fr);
  gap: 6px; max-width: 1100px; margin: 0 auto;
}
.grid-cell {
  aspect-ratio: 1; border-radius: 14px;
  background: rgba(255,255,255,0.025); border: 1px solid rgba(255,255,255,0.04);
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 4px; position: relative; overflow: hidden;
  transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
}
.grid-cell:not(.empty):hover { transform: scale(1.04); border-color: rgba(255,255,255,0.12); box-shadow: 0 4px 20px rgba(100,60,220,0.15); }
.grid-cell.flash { border-color: rgba(255,200,60,0.25); }
.grid-cell.empty { background: transparent; border-style: dashed; border-color: rgba(255,255,255,0.03); }
.empty-slot { width: 100%; height: 100%; }

/* Pet emoji */
.pet-emoji {
  font-size: 28px; line-height: 1; margin-top: 2px;
  filter: drop-shadow(0 0 6px rgba(180,140,255,0.3));
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.pet-emoji.bounce {
  animation: bounceAnim 0.6s cubic-bezier(0.28, 1.33, 0.64, 1) 2;
}
.pet-emoji.shake {
  animation: shakeAnim 0.5s ease-in-out 3;
}

/* EXP bar */
.pet-exp {
  width: 80%; height: 3px; background: rgba(255,255,255,0.06);
  border-radius: 2px; margin: 3px 0; overflow: hidden;
}
.exp-fill {
  height: 100%; background: linear-gradient(90deg, #7c3aed, #a78bfa);
  border-radius: 2px; transition: width 0.6s ease;
}

/* Cell info */
.cell-info {
  display: flex; flex-direction: column; align-items: center; gap: 0;
  margin-top: 1px;
}
.cell-name {
  font-size: 10px; font-weight: 600; color: rgba(220,210,250,0.8);
  max-width: 80px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  line-height: 1.2;
}
.cell-no {
  font-size: 8px; color: rgba(200,190,240,0.4); font-weight: 500;
}

/* Score ticker */
.score-ticker {
  position: fixed; bottom: 12px; left: 16px; right: 16px;
  background: rgba(20,16,50,0.85); backdrop-filter: blur(12px);
  border-radius: 12px; border: 1px solid rgba(255,255,255,0.06);
  padding: 6px 14px; display: flex; align-items: center; gap: 10px;
  max-width: 1100px; margin: 0 auto;
}
.ticker-label {
  font-size: 10px; font-weight: 600; color: rgba(200,190,240,0.5);
  white-space: nowrap; background: rgba(255,255,255,0.04); padding: 2px 8px; border-radius: 6px;
}
.ticker-items { display: flex; gap: 16px; overflow: hidden; }
.ticker-item { font-size: 11px; white-space: nowrap; color: rgba(200,190,240,0.7); }
.ticker-item .positive { color: #4ade80; }
.ticker-item .negative { color: #f87171; }

/* Message overlay */
.message-overlay {
  position: fixed; inset: 0; z-index: 100;
  background: rgba(10,5,30,0.85); backdrop-filter: blur(20px);
  display: flex; align-items: center; justify-content: center;
}
.message-overlay.banner { align-items: flex-start; padding-top: 40px; background: rgba(10,5,30,0.6); }
.message-overlay.fullscreen, .message-overlay.urgent { background: rgba(10,5,30,0.92); }
.message-content {
  text-align: center; max-width: 600px; padding: 40px;
  animation: messageIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.message-icon { font-size: 48px; margin-bottom: 16px; }
.message-text { font-size: 26px; font-weight: 700; line-height: 1.5; color: #f0ecff; }
.message-timer {
  margin-top: 20px; width: 200px; height: 3px; background: rgba(255,255,255,0.08);
  border-radius: 2px; margin-left: auto; margin-right: auto; overflow: hidden;
}
.timer-bar {
  height: 100%; background: linear-gradient(90deg, #7c3aed, #a78bfa);
  border-radius: 2px; animation: timerShrink linear forwards;
}
@keyframes timerShrink { from { width: 100%; } to { width: 0%; } }

/* Unlock modal */
.unlock-overlay {
  position: fixed; inset: 0; z-index: 200;
  background: rgba(5,2,20,0.9); backdrop-filter: blur(16px);
  display: flex; align-items: center; justify-content: center;
}
.unlock-modal {
  background: linear-gradient(135deg, #1a1040, #0d1b2a);
  border: 1px solid rgba(255,255,255,0.08); border-radius: 20px;
  padding: 36px 40px; text-align: center; max-width: 380px; width: 90%;
  box-shadow: 0 20px 60px rgba(0,0,0,0.5);
}
.unlock-icon { font-size: 40px; margin-bottom: 12px; }
.unlock-modal h3 { font-size: 16px; font-weight: 600; margin: 0 0 16px; color: #e8e0f8; }
.unlock-input {
  width: 100%; padding: 12px 16px; border-radius: 12px; border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.04); color: #f0ecff; font-size: 15px;
  text-align: center; outline: none; box-sizing: border-box;
}
.unlock-input:focus { border-color: #7c3aed; box-shadow: 0 0 0 3px rgba(124,58,237,0.15); }
.unlock-error { color: #f87171; font-size: 12px; margin-top: 8px; }
.unlock-buttons { display: flex; gap: 10px; margin-top: 16px; }
.unlock-btn {
  flex: 1; padding: 10px; border-radius: 10px; border: none; font-size: 14px;
  font-weight: 600; cursor: pointer; transition: all 0.15s;
}
.unlock-btn.cancel { background: rgba(255,255,255,0.06); color: rgba(200,190,240,0.7); }
.unlock-btn.confirm { background: #7c3aed; color: white; }
.unlock-btn.confirm:hover { background: #6d28d9; }

/* Animations */
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

.message-pop-enter-active { animation: messageIn 0.4s ease-out; }
.message-pop-leave-active { animation: messageIn 0.3s ease-in reverse; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.loading, .empty { text-align: center; padding: 80px 20px; color: rgba(200,190,240,0.5); }
.empty-icon { font-size: 48px; margin-bottom: 12px; }
</style>
