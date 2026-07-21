<script setup lang="ts">
import { useToastStore } from '@/stores/toast'
import { computed } from 'vue'

const store = useToastStore()

const iconMap: Record<string, string> = {
  success: '✅',
  error: '❌',
  info: 'ℹ️',
  warning: '⚠️',
}

const borderMap: Record<string, string> = {
  success: '#10B981',
  error: '#EF4444',
  info: '#3B82F6',
  warning: '#F59E0B',
}

const groupedToasts = computed(() => {
  const groups: Record<string, any[]> = {}
  for (const t of store.toasts) {
    if (!groups[t.position]) groups[t.position] = []
    groups[t.position].push(t)
  }
  return groups
})

function progressPercent(toast: any) {
  if (toast.duration <= 0) return 100
  const elapsed = Date.now() - toast.createdAt
  return Math.max(0, Math.min(100, (1 - elapsed / toast.duration) * 100))
}

let progressTimers: Record<number, number> = {}
function startProgress(id: number) {
  const t = store.toasts.find(x => x.id === id)
  if (!t || t.duration <= 0) return
  const update = () => {
    const p = progressPercent(t)
    const bar = document.getElementById(`tp-${id}`)
    if (bar) bar.style.width = p + '%'
    if (p > 0) progressTimers[id] = requestAnimationFrame(update)
  }
  progressTimers[id] = requestAnimationFrame(update)
}
function stopProgress(id: number) {
  if (progressTimers[id]) cancelAnimationFrame(progressTimers[id])
}
</script>

<template>
  <!-- 右上角 Toast（默认） -->
  <div v-if="groupedToasts['top-right']?.length" class="toast-region tr">
    <TransitionGroup name="tr">
      <div
        v-for="t in groupedToasts['top-right']" :key="t.id"
        class="toast-card"
        :style="{ borderLeftColor: borderMap[t.type] || borderMap.info }"
        @mouseenter="stopProgress(t.id)" @mouseleave="startProgress(t.id)"
      >
        <span class="toast-icon">{{ iconMap[t.type] || 'ℹ️' }}</span>
        <span class="toast-msg">{{ t.message }}</span>
        <button v-if="t.action" class="toast-action" @click="t.action.onClick; store.remove(t.id)">{{ t.action.label }}</button>
        <button class="toast-close" @click="store.remove(t.id)">✕</button>
        <div class="toast-progress"><div :id="`tp-${t.id}`" class="toast-progress-bar" :style="{ background: borderMap[t.type] }"></div></div>
      </div>
    </TransitionGroup>
  </div>

  <!-- 居中 Toast（重要操作反馈） -->
  <div v-if="groupedToasts['center']?.length" class="toast-region center">
    <div class="center-backdrop" @click="store.remove(groupedToasts['center'][0]?.id)"></div>
    <TransitionGroup name="center">
      <div
        v-for="t in groupedToasts['center']" :key="t.id"
        class="toast-card toast-card--center"
        @mouseenter="stopProgress(t.id)" @mouseleave="startProgress(t.id)"
      >
        <div class="center-icon" :style="{ background: borderMap[t.type] + '18' }">
          <span style="font-size:32px;">{{ iconMap[t.type] || 'ℹ️' }}</span>
        </div>
        <span class="toast-msg toast-msg--center">{{ t.message }}</span>
        <button v-if="t.action" class="toast-action toast-action--center" @click="t.action.onClick; store.remove(t.id)">{{ t.action.label }}</button>
        <div class="toast-progress"><div :id="`tp-${t.id}`" class="toast-progress-bar" :style="{ background: borderMap[t.type] }"></div></div>
      </div>
    </TransitionGroup>
  </div>

  <!-- 底部 Toast（全局通知） -->
  <div v-if="groupedToasts['bottom-center']?.length" class="toast-region bc">
    <TransitionGroup name="bc">
      <div
        v-for="t in groupedToasts['bottom-center']" :key="t.id"
        class="toast-card toast-card--bottom"
        :style="{ borderLeftColor: borderMap[t.type] }"
      >
        <span class="toast-icon">{{ iconMap[t.type] || 'ℹ️' }}</span>
        <span class="toast-msg">{{ t.message }}</span>
        <button v-if="t.action" class="toast-action" @click="t.action.onClick; store.remove(t.id)">{{ t.action.label }}</button>
        <button class="toast-close" @click="store.remove(t.id)">✕</button>
        <div class="toast-progress"><div :id="`tp-${t.id}`" class="toast-progress-bar" :style="{ background: borderMap[t.type] }"></div></div>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
/* ===== 定位 ===== */
.toast-region { position: fixed; z-index: 9999; pointer-events: none; }
.toast-region.tr { top: 16px; right: 16px; display:flex; flex-direction:column; gap:8px; max-width:380px; }
.toast-region.center { inset:0; display:flex; align-items:center; justify-content:center; }
.toast-region.bc { bottom:32px; left:50%; transform:translateX(-50%); display:flex; flex-direction:column; gap:8px; max-width:420px; width:90%; }

.center-backdrop { position:fixed; inset:0; background:rgba(0,0,0,0.2); backdrop-filter:blur(4px); pointer-events:auto; }

/* ===== 卡片 ===== */
.toast-card {
  pointer-events: auto;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: rgba(30, 30, 46, 0.92);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-radius: 14px;
  border: 1px solid rgba(255,255,255,0.08);
  border-left: 4px solid #10B981;
  box-shadow: 0 12px 40px rgba(0,0,0,0.25);
  color: #fff;
  font-size: 13px;
  position: relative;
  overflow: hidden;
  min-width: 280px;
}
.toast-card--center {
  flex-direction: column;
  padding: 32px 40px;
  min-width: 200px;
  max-width: 360px;
  text-align: center;
  border-left: none;
  border-radius: 20px;
}
.toast-card--bottom {
  border-radius: 16px;
  box-shadow: 0 -4px 24px rgba(0,0,0,0.15);
}

.center-icon {
  width: 64px; height: 64px;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
}

.toast-icon { font-size: 16px; flex-shrink: 0; line-height: 1; }
.toast-msg { flex: 1; line-height: 1.4; }
.toast-msg--center { font-size: 16px; font-weight: 600; margin: 4px 0; }
.toast-close { background:none; border:none; color:rgba(255,255,255,0.4); cursor:pointer; font-size:14px; padding:2px; line-height:1; flex-shrink:0; }
.toast-close:hover { color:rgba(255,255,255,0.8); }
.toast-action { background:rgba(255,255,255,0.1); border:none; color:#fff; padding:4px 12px; border-radius:8px; font-size:12px; cursor:pointer; white-space:nowrap; flex-shrink:0; }
.toast-action:hover { background:rgba(255,255,255,0.2); }
.toast-action--center { padding:8px 24px; font-size:14px; margin-top:4px; border-radius:10px; }

/* ===== 进度条 ===== */
.toast-progress { position:absolute; bottom:0; left:0; right:0; height:3px; background:rgba(255,255,255,0.05); }
.toast-progress-bar { height:100%; border-radius:0 0 0 14px; transition:width 0.1s linear; }

/* ===== 动画（右上角） ===== */
.tr-enter-active { animation: slideInRight 0.35s cubic-bezier(0.34,1.56,0.64,1) forwards; }
.tr-leave-active { animation: fadeOut 0.25s ease forwards; }
@keyframes slideInRight { 0% { opacity:0; transform:translateX(60px) scale(0.95); } 100% { opacity:1; transform:translateX(0) scale(1); } }
@keyframes fadeOut { 0% { opacity:1; transform:scale(1); } 100% { opacity:0; transform:scale(0.9); } }

/* ===== 动画（居中） ===== */
.center-enter-active { animation: popIn 0.4s cubic-bezier(0.34,1.56,0.64,1) forwards; }
.center-leave-active { animation: popOut 0.25s ease forwards; }
@keyframes popIn { 0% { opacity:0; transform:scale(0.8) translateY(20px); } 100% { opacity:1; transform:scale(1) translateY(0); } }
@keyframes popOut { 0% { opacity:1; transform:scale(1); } 100% { opacity:0; transform:scale(0.85); } }

/* ===== 动画（底部） ===== */
.bc-enter-active { animation: slideUp 0.35s cubic-bezier(0.34,1.56,0.64,1) forwards; }
.bc-leave-active { animation: fadeOut 0.25s ease forwards; }
@keyframes slideUp { 0% { opacity:0; transform:translateY(40px); } 100% { opacity:1; transform:translateY(0); } }
</style>
