<script setup lang="ts">
import { useToastStore } from '@/stores/toast'
import { computed } from 'vue'

const store = useToastStore()

/** 类型 → emoji 字形。全站统一 emoji 口径（2026-09-19 定），装饰性图标一律不用线性 SVG。 */
const iconMap: Record<string, string> = {
  success: '✅',
  error: '❌',
  info: 'ℹ️',
  warning: '⚠️',
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
        :class="'toast-card--' + (t.type || 'info')"
        @mouseenter="stopProgress(t.id)" @mouseleave="startProgress(t.id)"
      >
        <span class="toast-icon">{{ iconMap[t.type] || 'ℹ️' }}</span>
        <span class="toast-msg">{{ t.message }}</span>
        <button v-if="t.action" class="toast-action" @click="t.action.onClick; store.remove(t.id)">{{ t.action.label }}</button>
        <button class="toast-close" @click="store.remove(t.id)" aria-label="关闭">✕</button>
        <div class="toast-progress"><div :id="`tp-${t.id}`" class="toast-progress-bar"></div></div>
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
        :class="'toast-card--' + (t.type || 'info')"
        @mouseenter="stopProgress(t.id)" @mouseleave="startProgress(t.id)"
      >
        <div class="center-icon">
          <span class="tst-icon-32">{{ iconMap[t.type] || 'ℹ️' }}</span>
        </div>
        <span class="toast-msg toast-msg--center">{{ t.message }}</span>
        <button v-if="t.action" class="toast-action toast-action--center" @click="t.action.onClick; store.remove(t.id)">{{ t.action.label }}</button>
        <div class="toast-progress"><div :id="`tp-${t.id}`" class="toast-progress-bar"></div></div>
      </div>
    </TransitionGroup>
  </div>

  <!-- 底部 Toast（全局通知） -->
  <div v-if="groupedToasts['bottom-center']?.length" class="toast-region bc">
    <TransitionGroup name="bc">
      <div
        v-for="t in groupedToasts['bottom-center']" :key="t.id"
        class="toast-card toast-card--bottom"
        :class="'toast-card--' + (t.type || 'info')"
      >
        <span class="toast-icon">{{ iconMap[t.type] || 'ℹ️' }}</span>
        <span class="toast-msg">{{ t.message }}</span>
        <button v-if="t.action" class="toast-action" @click="t.action.onClick; store.remove(t.id)">{{ t.action.label }}</button>
        <button class="toast-close" @click="store.remove(t.id)" aria-label="关闭">✕</button>
        <div class="toast-progress"><div :id="`tp-${t.id}`" class="toast-progress-bar"></div></div>
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
  background: var(--color-modal-bg);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-radius: var(--ui-r-xl);
  border: 1px solid var(--ui-border);
  border-left: 3px solid var(--ui-border-strong);
  box-shadow: var(--shadow-lg);
  color: var(--ui-fg);
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
  border-left-width: 1px;
  border-radius: var(--ui-r-xl);
}
.toast-card--bottom {
  border-radius: var(--ui-r-lg);
  box-shadow: 0 -4px 24px rgba(0,0,0,0.15);
}

/* ===== 类型语义色：统一走 --c-* 变量，亮暗自动跟随 ===== */
.toast-card--success { border-left-color: var(--c-green); }
.toast-card--error { border-left-color: var(--c-red); }
.toast-card--info { border-left-color: var(--c-blue); }
.toast-card--warning { border-left-color: var(--c-amber); }
.toast-card--success .toast-progress-bar { background: var(--c-green); }
.toast-card--error .toast-progress-bar { background: var(--c-red); }
.toast-card--info .toast-progress-bar { background: var(--c-blue); }
.toast-card--warning .toast-progress-bar { background: var(--c-amber); }
/* 居中卡片的圆形底色：emoji 自带色彩，此处只做浅色衬底 */
.toast-card--success .center-icon { background: color-mix(in srgb, var(--c-green) 12%, transparent); }
.toast-card--error .center-icon { background: color-mix(in srgb, var(--c-red) 12%, transparent); }
.toast-card--info .center-icon { background: color-mix(in srgb, var(--c-blue) 12%, transparent); }
.toast-card--warning .center-icon { background: color-mix(in srgb, var(--c-amber) 12%, transparent); }

.center-icon {
  width: 64px; height: 64px;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
}

.toast-icon { font-size: 16px; flex-shrink: 0; line-height: 1; }
.toast-msg { flex: 1; line-height: 1.4; }
.toast-msg--center { font-size: 16px; font-weight: 600; margin: 4px 0; }
.toast-close { background:none; border:none; color:var(--ui-fg-subtle); cursor:pointer; font-size:14px; padding:2px; line-height:1; flex-shrink:0; transition: color 0.15s ease; }
.toast-close:hover { color: var(--ui-fg); }
.toast-action { background:var(--ui-muted); border:none; color:var(--ui-fg); padding:4px 12px; border-radius:var(--ui-r-sm); font-size:12px; cursor:pointer; white-space:nowrap; flex-shrink:0; transition: background 0.15s ease; }
.toast-action:hover { background:var(--ui-brand-soft); color: var(--ui-brand); }
.toast-action--center { padding:8px 24px; font-size:14px; margin-top:4px; border-radius:var(--ui-r-md); }

/* ===== 进度条 ===== */
.toast-progress { position:absolute; bottom:0; left:0; right:0; height:3px; background:var(--ui-muted); }
.toast-progress-bar { height:100%; border-radius:0 0 0 var(--ui-r-xl); transition:width 0.1s linear; }

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
.tst-icon-32 { font-size:32px; line-height: 1; }
</style>
