<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet, apiPost } from '@/utils/api'
import { useDisplaySSE, type ScoreUpdateData, type BroadcastData, type PetUpdateData } from '@/composables/useDisplaySSE'

const router = useRouter()

const token = ref('')
const classInfo = ref<{ id: number; name: string; grade: string; student_count: number } | null>(null)
const loadError = ref('')
const loading = ref(true)

interface PetEntry {
  student_id: number; student_no: string; student_name: string; total_score: number
  has_pet: boolean; pet_name?: string; level: number; experience: number; mood: number
  emoji: string; stage_name: string; exp_max: number
}

interface DisplayData {
  class_name: string; grade: string; student_count: number
  pets: PetEntry[]; recent_scores: any[]; broadcasts: any[]
}

const data = ref<DisplayData | null>(null)
const scoreAnimations = ref<Record<number, { dir: 'up' | 'down'; amt: number }>>({})
const scorePending = ref<Record<number, boolean>>({})

const gridSlots = computed(() => {
  if (!data.value) return []
  const slots: (PetEntry | null)[] = [...(data.value.pets || [])]
  while (slots.length < 64) slots.push(null)
  return slots.slice(0, 64)
})

// ===== 教师操作模式（防误触：连点 3 次激活） =====
const teacherMode = ref(false)
const tapCount = ref(0)
let tapTimer: ReturnType<typeof setTimeout> | null = null

function activateTeacherMode() {
  tapCount.value++
  if (tapCount.value >= 3) {
    teacherMode.value = !teacherMode.value
    tapCount.value = 0
    if (tapTimer) clearTimeout(tapTimer)
  } else {
    if (tapTimer) clearTimeout(tapTimer)
    tapTimer = setTimeout(() => { tapCount.value = 0 }, 2000)
  }
}

// ===== 广播 =====
const activeBroadcast = ref<BroadcastData | null>(null)
const broadcastTimer = ref<ReturnType<typeof setTimeout> | null>(null)

// ===== SSE =====
const { state: sseState, scoreUpdates, broadcasts, petUpdates, connect: connectSSE } = useDisplaySSE()

onMounted(() => {
  window.addEventListener('display:token-expired', () => router.replace({ name: 'display-login' }))
  const storedToken = sessionStorage.getItem('display_token')
  const storedInfo = sessionStorage.getItem('display_class_info')
  if (!storedToken || !storedInfo) { router.replace({ name: 'display-login' }); return }
  token.value = storedToken
  classInfo.value = JSON.parse(storedInfo)
  loadInitialData()
})

async function loadInitialData() {
  loading.value = true; loadError.value = ''
  try {
    const res = await apiGet<{ data: DisplayData }>('/api/v1/display/initial-data', { params: { token: token.value } })
    data.value = res.data; loading.value = false
    connectSSE(token.value, classInfo.value!.id)
  } catch (e: any) {
    if (e?.response?.status === 401) { sessionStorage.clear(); router.replace({ name: 'display-login' }); return }
    loadError.value = '加载失败'; loading.value = false
  }
}

// ===== 快捷加减分 =====
async function quickScore(studentId: number, amount: number) {
  if (scorePending.value[studentId]) return
  scorePending.value[studentId] = true
  try {
    await apiPost('/api/v1/display/quick-score', { token: token.value, student_id: studentId, amount })
    // 乐观更新
    const pet = data.value?.pets.find(p => p.student_id === studentId)
    if (pet) pet.total_score += amount
    triggerAnimation(studentId, amount)
  } catch { /* ignore */ }
  finally { scorePending.value[studentId] = false }
}

function triggerAnimation(studentId: number, amount: number) {
  scoreAnimations.value[studentId] = { dir: amount > 0 ? 'up' : 'down', amt: Math.abs(amount) }
  setTimeout(() => { delete scoreAnimations.value[studentId] }, 1800)
}

// SSE 事件
watch(scoreUpdates, (evts) => {
  for (const e of evts) {
    triggerAnimation(e.student_id, e.amount)
    const pet = data.value?.pets.find(p => p.student_id === e.student_id)
    if (pet) { pet.total_score = e.total_score; if (e.pet_level !== undefined) pet.level = e.pet_level; if (e.pet_experience !== undefined) pet.experience = e.pet_experience; if (e.pet_mood !== undefined) pet.mood = e.pet_mood }
  }
}, { deep: true })

watch(broadcasts, (evts) => {
  const m = evts[evts.length - 1]
  if (m) { if (broadcastTimer.value) clearTimeout(broadcastTimer.value); activeBroadcast.value = m; broadcastTimer.value = setTimeout(() => { activeBroadcast.value = null }, Math.max(3000, (m.display_seconds || 8) * 1000)) }
}, { deep: true })

watch(petUpdates, (evts) => {
  for (const e of evts) {
    const pet = data.value?.pets.find(p => p.student_id === e.student_id)
    if (pet) { if (e.mood !== undefined) pet.mood = e.mood; if (e.level !== undefined) pet.level = e.level; if (e.experience !== undefined) pet.experience = e.experience }
  }
}, { deep: true })

function confirmExit() { sessionStorage.clear(); router.replace({ name: 'display-login' }) }
</script>

<template>
  <div class="screen" :class="{ 'teacher-active': teacherMode }">
    <!-- 广播 -->
    <Transition name="pop">
      <div v-if="activeBroadcast" class="broadcast" :class="activeBroadcast.type" @click="activeBroadcast = null">
        <div class="broadcast-inner" @click.stop>
          <div class="broadcast-icon">{{ activeBroadcast.type === 'fullscreen' ? '📢' : '💬' }}</div>
          <div class="broadcast-text">{{ activeBroadcast.content }}</div>
          <div class="broadcast-bar"><div class="broadcast-fill" :style="{ animationDuration: (activeBroadcast.display_seconds || 8) + 's' }"></div></div>
        </div>
      </div>
    </Transition>

    <!-- 连接状态 -->
    <div class="status" :class="{ on: sseState.connected, pol: sseState.polling }">
      <span class="status-dot"></span>{{ sseState.connected ? '实时' : sseState.polling ? '轮询' : '重连' }}
    </div>

    <!-- 头部 -->
    <header class="header" @dblclick="activateTeacherMode">
      <div class="header-left">
        <h1 class="header-title">{{ data?.class_name || '--' }}</h1>
        <span class="header-meta">{{ data?.grade }} · {{ data?.student_count }}人</span>
      </div>
      <button class="header-exit" @click="confirmExit" title="退出大屏">✕</button>
    </header>

    <!-- 教师模式工具栏 -->
    <div v-if="teacherMode" class="toolbar">
      <span class="toolbar-label">👨‍🏫 课堂模式 · 点击学生快速加减分</span>
      <button class="toolbar-btn" @click="teacherMode = false">关闭</button>
    </div>

    <!-- 加载 -->
    <div v-if="loading" class="state"><div class="state-spin">🌌</div><p>连接中…</p></div>
    <div v-else-if="loadError" class="state"><p>{{ loadError }}</p><button class="state-btn" @click="loadInitialData">重试</button></div>

    <!-- 矩阵 -->
    <div v-else-if="data" class="grid">
      <div v-for="(s, i) in gridSlots" :key="i" class="cell"
        :class="{ empty: !s, 'score-up': s && scoreAnimations[s.student_id]?.dir === 'up', 'score-down': s && scoreAnimations[s.student_id]?.dir === 'down' }">

        <template v-if="s">
          <!-- 宠物 -->
          <div class="cell-pet" :class="{ bnc: scoreAnimations[s.student_id]?.dir === 'up', shk: scoreAnimations[s.student_id]?.dir === 'down' }">
            <span class="cell-emoji">{{ s.has_pet ? s.emoji : '🥚' }}</span>
            <Transition name="flt">
              <span v-if="scoreAnimations[s.student_id]" class="cell-float" :class="scoreAnimations[s.student_id].dir">
                {{ scoreAnimations[s.student_id].dir === 'up' ? '+' : '-' }}{{ scoreAnimations[s.student_id].amt }}
              </span>
            </Transition>
          </div>

          <!-- 信息 -->
          <div class="cell-info">
            <span class="cell-name">{{ s.student_name }}</span>
            <span class="cell-no">{{ s.student_no || '--' }}</span>
          </div>
          <div class="cell-score">{{ s.total_score }}分</div>

          <!-- 经验条 -->
          <div class="cell-exp"><div class="cell-exp-fill" :style="{ width: Math.min(100, (s.experience / Math.max(1, s.exp_max)) * 100) + '%' }"></div></div>

          <!-- 教师模式：快捷加减 -->
          <div v-if="teacherMode" class="cell-actions">
            <button class="ca-btn ca-minus" @click.stop="quickScore(s.student_id, -1)" :disabled="scorePending[s.student_id]">−1</button>
            <button class="ca-btn ca-minus" @click.stop="quickScore(s.student_id, -3)" :disabled="scorePending[s.student_id]">−3</button>
            <button class="ca-btn ca-plus" @click.stop="quickScore(s.student_id, 1)" :disabled="scorePending[s.student_id]">+1</button>
            <button class="ca-btn ca-plus" @click.stop="quickScore(s.student_id, 3)" :disabled="scorePending[s.student_id]">+3</button>
          </div>
        </template>
        <div v-else class="cell-empty"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.screen {
  min-height: 100vh;
  background: linear-gradient(135deg,#0c0a20 0%,#1a1040 30%,#0d1b2a 70%,#0a1628 100%);
  color: #e8e6f0;
  font-family: "PingFang SC","Noto Sans SC",-apple-system,sans-serif;
  padding: 10px 14px 16px;
  user-select: none;
  position: relative;
}

/* 背景 */
.screen::before {
  content: ''; position: absolute; inset: 0;
  background: radial-gradient(ellipse at 20% 20%,rgba(120,80,255,.06),transparent 60%),
              radial-gradient(ellipse at 80% 80%,rgba(30,140,220,.04),transparent 60%);
  pointer-events: none;
}

/* 连接状态 */
.status {
  position: fixed; top: 8px; right: 44px; z-index: 50;
  display: flex; align-items: center; gap: 5px;
  font-size: 10px; color: rgba(200,190,240,.25);
  padding: 2px 8px; border-radius: 20px; background: rgba(0,0,0,.3);
}
.status-dot { width: 5px; height: 5px; border-radius: 50%; background: #64748b; }
.status.on .status-dot { background: #4ade80; box-shadow: 0 0 6px rgba(74,222,128,.5); }
.status.pol .status-dot { background: #f59e0b; box-shadow: 0 0 6px rgba(245,158,11,.5); }

/* 头部 */
.header {
  display: flex; justify-content: space-between; align-items: center;
  padding-bottom: 8px; border-bottom: 1px solid rgba(255,255,255,.06); margin-bottom: 8px;
  position: relative; z-index: 1;
}
.header-left { display: flex; align-items: baseline; gap: 10px; }
.header-title { font-size: 18px; font-weight: 700; margin: 0; color: #f0ecff; }
.header-meta { font-size: 12px; color: rgba(200,190,240,.5); }
.header-exit {
  background: none; border: 1px solid rgba(255,255,255,.06); color: rgba(200,190,240,.3);
  width: 28px; height: 28px; border-radius: 8px; font-size: 12px; cursor: pointer;
  display: flex; align-items: center; justify-content: center; opacity: .4;
}
.header-exit:hover { opacity: 1; color: #f87171; }

/* 工具栏 */
.toolbar {
  position: relative; z-index: 1;
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 14px; margin-bottom: 8px;
  background: rgba(124,58,237,.15); border: 1px solid rgba(124,58,237,.2);
  border-radius: 10px; font-size: 13px;
}
.toolbar-label { color: #c4b5fd; font-weight: 500; }
.toolbar-btn {
  padding: 4px 14px; border-radius: 8px; border: none;
  background: rgba(255,255,255,.1); color: #e8e6f0; font-size: 12px;
  cursor: pointer; font-family: inherit;
}
.toolbar-btn:hover { background: rgba(255,255,255,.15); }

/* 加载 */
.state { text-align: center; padding: 80px 20px; color: rgba(200,190,240,.5); }
.state-spin { font-size: 48px; animation: spin 2s linear infinite; }
.state-btn { padding: 8px 20px; border-radius: 10px; border: 1px solid rgba(255,255,255,.1); background: rgba(255,255,255,.04); color: rgba(200,190,240,.7); cursor: pointer; font-size: 14px; font-family: inherit; }

/* 8x8 网格 */
.grid {
  display: grid; grid-template-columns: repeat(8, 1fr); gap: 5px;
  max-width: 1120px; margin: 0 auto; position: relative; z-index: 1;
}

.cell {
  aspect-ratio: 1; border-radius: 12px;
  background: rgba(255,255,255,.025); border: 1px solid rgba(255,255,255,.04);
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 3px; position: relative; overflow: hidden;
  transition: transform .15s, border-color .15s;
  min-width: 0;
}
.cell.empty { background: transparent; border-style: dashed; border-color: rgba(255,255,255,.03); }
.cell-empty { width: 100%; height: 100%; }
.cell.score-up { border-color: rgba(74,222,128,.2); box-shadow: inset 0 0 20px rgba(74,222,128,.04); }
.cell.score-down { border-color: rgba(248,113,113,.2); box-shadow: inset 0 0 20px rgba(248,113,113,.04); }

/* 宠物 */
.cell-pet { position: relative; display: flex; align-items: center; justify-content: center; margin-top: 2px; }
.cell-emoji { font-size: 26px; line-height: 1; filter: drop-shadow(0 0 5px rgba(180,140,255,.2)); transition: all .3s; }
.cell-float {
  position: absolute; top: -10px; right: -14px; font-size: 10px; font-weight: 700;
  padding: 1px 4px; border-radius: 4px; pointer-events: none;
}
.cell-float.up { color: #4ade80; background: rgba(74,222,128,.15); }
.cell-float.down { color: #f87171; background: rgba(248,113,113,.15); }
.cell-pet.bnc .cell-emoji { animation: bounce .5s cubic-bezier(.28,1.33,.64,1) 2; }
.cell-pet.shk .cell-emoji { animation: shake .4s ease-in-out 2; }

/* 信息 */
.cell-info { display: flex; flex-direction: column; align-items: center; line-height: 1.1; margin-top: 1px; }
.cell-name { font-size: 10px; font-weight: 600; color: rgba(220,210,250,.8); max-width: 70px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cell-no { font-size: 8px; color: rgba(200,190,240,.35); }
.cell-score { font-size: 11px; font-weight: 700; color: rgba(200,190,240,.6); margin-top: 1px; }

/* 经验条 */
.cell-exp { width: 70%; height: 2px; background: rgba(255,255,255,.05); border-radius: 2px; margin-top: 2px; overflow: hidden; }
.cell-exp-fill { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); border-radius: 2px; transition: width .5s; }

/* 快捷操作（教师模式） */
.cell-actions {
  display: none; gap: 2px; margin-top: auto; padding: 2px 0;
}
.teacher-active .cell-actions { display: flex; }
.ca-btn {
  flex: 1; padding: 2px 0; border: none; border-radius: 4px;
  font-size: 9px; font-weight: 700; cursor: pointer;
  transition: all .1s; font-family: inherit; line-height: 1;
}
.ca-btn:disabled { opacity: .3; cursor: not-allowed; }
.ca-plus { background: rgba(74,222,128,.15); color: #4ade80; }
.ca-plus:hover:not(:disabled) { background: rgba(74,222,128,.25); }
.ca-minus { background: rgba(248,113,113,.15); color: #f87171; }
.ca-minus:hover:not(:disabled) { background: rgba(248,113,113,.25); }

/* 广播 */
.broadcast {
  position: fixed; inset: 0; z-index: 100;
  display: flex; align-items: center; justify-content: center;
}
.broadcast.banner { align-items: flex-start; padding-top: 30px; background: rgba(10,5,30,.6); backdrop-filter: blur(8px); }
.broadcast.popup { background: rgba(10,5,30,.75); backdrop-filter: blur(12px); }
.broadcast.fullscreen { background: rgba(10,5,30,.92); backdrop-filter: blur(20px); }
.broadcast-inner { text-align: center; max-width: 560px; padding: 30px; animation: popIn .35s cubic-bezier(.34,1.56,.64,1); }
.broadcast-icon { font-size: 44px; margin-bottom: 12px; }
.broadcast-text { font-size: 24px; font-weight: 700; color: #f0ecff; line-height: 1.4; }
.broadcast-bar { margin: 16px auto 0; width: 180px; height: 3px; background: rgba(255,255,255,.06); border-radius: 2px; overflow: hidden; }
.broadcast-fill { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); border-radius: 2px; animation: shrink linear forwards; }

@keyframes shrink { from { width: 100%; } to { width: 0%; } }
@keyframes bounce { 0%,100%{transform:translateY(0)} 30%{transform:translateY(-6px) scale(1.1)} 60%{transform:translateY(-2px) scale(1.03)} }
@keyframes shake { 0%,100%{transform:translateX(0)} 25%{transform:translateX(-2px)} 75%{transform:translateX(2px)} }
@keyframes popIn { from { opacity: 0; transform: scale(.9) translateY(20px); } to { opacity: 1; transform: scale(1) translateY(0); } }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.pop-enter-active { animation: popIn .35s ease-out; }
.pop-leave-active { animation: popIn .25s ease-in reverse; }
.flt-enter-active { transition: all .25s ease-out; }
.flt-leave-active { transition: all .15s ease-in; }
.flt-enter-from { opacity: 0; transform: translateY(6px); }
.flt-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
