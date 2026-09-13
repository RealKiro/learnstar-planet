<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse, TimetablePeriod, TimetableSubject, TimetableEntry, TimetableWeekType } from '@/types'

interface DisplayTimetableData {
  today_weekday: number
  class_name: string
  subjects: TimetableSubject[]
  periods: TimetablePeriod[]
  entries: TimetableEntry[]
}

const loading = ref(true)
const loadError = ref('')
const data = ref<DisplayTimetableData | null>(null)
const viewMode = ref<'today' | 'week'>('today')
const viewWeekday = ref(new Date().getDay() === 0 ? 7 : new Date().getDay())
const nowMinutes = ref(0)

const WEEKDAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const WEEK_TYPE_LABELS: Record<TimetableWeekType, string> = { all: '', odd: '单周', even: '双周' }

let clockTimer: ReturnType<typeof setInterval> | null = null

const todayWeekday = computed(() => data.value?.today_weekday ?? 1)
const shownWeekday = computed(() => (viewMode.value === 'today' ? todayWeekday.value : viewWeekday.value))

const dayEntries = computed(() => {
  if (!data.value) return []
  return data.value.entries
    .filter(e => e.weekday === shownWeekday.value)
    .sort((a, b) => a.period_index - b.period_index)
})

const periodByIdx = computed(() => {
  const map = new Map<number, TimetablePeriod>()
  for (const p of data.value?.periods || []) map.set(p.period_index, p)
  return map
})

function parseHm(time: string | undefined): number {
  if (!time) return -1
  const [h, m] = time.split(':').map(Number)
  return (h || 0) * 60 + (m || 0)
}

/** 当前节：按节次起止时间判断；返回 period_index，无匹配返回 0 */
const currentPeriodIdx = computed(() => {
  if (shownWeekday.value !== todayWeekday.value) return 0
  for (const p of data.value?.periods || []) {
    const start = parseHm(p.start_time)
    const end = parseHm(p.end_time)
    if (start >= 0 && end > start && nowMinutes.value >= start && nowMinutes.value <= end) {
      return p.period_index
    }
  }
  return 0
})

function subjectColor(name: string): string | undefined {
  return data.value?.subjects.find(s => s.name === name)?.color || undefined
}

function toTime(min: number): string {
  return `${String(Math.floor(min / 60)).padStart(2, '0')}:${String(min % 60).padStart(2, '0')}`
}

onMounted(async () => {
  const tick = () => {
    const now = new Date()
    nowMinutes.value = now.getHours() * 60 + now.getMinutes()
  }
  tick()
  clockTimer = setInterval(tick, 30_000)
  try {
    const res = await apiGet<ApiResponse<DisplayTimetableData>>('/api/v1/display/timetable', {
      params: { token: sessionStorage.getItem('class_token') || '' },
    })
    data.value = res.data || null
  } catch {
    loadError.value = '课表加载失败'
  } finally {
    loading.value = false
  }
})

onUnmounted(() => { if (clockTimer) clearInterval(clockTimer) })
</script>

<template>
  <div class="tt-page">
    <div class="tt-head">
      <div>
        <h2 class="tt-title">今日课表</h2>
        <p class="tt-sub">{{ data?.class_name || '' }} · {{ WEEKDAY_LABELS[todayWeekday - 1] }}<template v-if="data"> · {{ toTime(nowMinutes) }}</template></p>
      </div>
      <div class="tt-switch">
        <button :class="['tt-switch-btn', { active: viewMode === 'today' }]" @click="viewMode = 'today'">今天</button>
        <button :class="['tt-switch-btn', { active: viewMode === 'week' }]" @click="viewMode = 'week'">整周</button>
      </div>
    </div>

    <div v-if="viewMode === 'week'" class="tt-week-tabs">
      <button v-for="d in 7" :key="d" :class="['tt-tab', { active: viewWeekday === d, today: d === todayWeekday }]" @click="viewWeekday = d">
        {{ WEEKDAY_LABELS[d - 1] }}
      </button>
    </div>

    <div v-if="loading" class="tt-empty">加载中...</div>
    <div v-else-if="loadError" class="tt-empty tt-error">{{ loadError }}</div>
    <div v-else-if="dayEntries.length === 0" class="tt-empty">
      <span class="tt-empty-icon">🗓️</span>
      {{ WEEKDAY_LABELS[shownWeekday - 1] }}暂无课程安排
    </div>

    <div v-else class="tt-list">
      <div
        v-for="e in dayEntries"
        :key="`${e.period_index}-${e.week_type}`"
        :class="['tt-row', { current: currentPeriodIdx === e.period_index }]"
      >
        <div class="tt-time">
          <template v-if="periodByIdx.get(e.period_index)">
            <span class="tt-time-range">{{ periodByIdx.get(e.period_index)!.start_time }}–{{ periodByIdx.get(e.period_index)!.end_time }}</span>
            <span class="tt-period-name">{{ periodByIdx.get(e.period_index)!.name || `第${e.period_index}节` }}</span>
          </template>
          <template v-else>
            <span class="tt-period-name">第{{ e.period_index }}节</span>
          </template>
        </div>
        <div class="tt-main">
          <span class="tt-subject" :style="subjectColor(e.subject_name) ? { color: subjectColor(e.subject_name) } : {}">{{ e.subject_name }}</span>
          <span v-if="WEEK_TYPE_LABELS[e.week_type]" class="tt-week">{{ WEEK_TYPE_LABELS[e.week_type] }}</span>
        </div>
        <div class="tt-meta">{{ [e.teacher_name, e.room].filter(Boolean).join(' · ') || '&nbsp;' }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tt-page { color: var(--color-text); }
.tt-head { display: flex; align-items: flex-end; justify-content: space-between; margin-bottom: 16px; }
.tt-title { font-size: 28px; font-weight: 800; margin: 0; }
.tt-sub { font-size: 14px; color: var(--color-text-secondary); margin: 4px 0 0; }
.tt-switch { display: flex; gap: 4px; background: var(--color-bg); border-radius: 10px; padding: 4px; }
.tt-switch-btn { border: none; background: transparent; padding: 6px 16px; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; color: var(--color-text-secondary); font-family: inherit; }
.tt-switch-btn.active { background: var(--color-bg-card, #fff); color: var(--color-text); box-shadow: 0 1px 4px rgba(0,0,0,.08); }
.tt-week-tabs { display: flex; gap: 6px; margin-bottom: 14px; flex-wrap: wrap; }
.tt-tab { border: 1.5px solid var(--color-border); background: transparent; border-radius: 9999px; padding: 5px 14px; font-size: 13px; font-weight: 600; cursor: pointer; color: var(--color-text-secondary); font-family: inherit; }
.tt-tab.active { background: var(--md-primary, #7c3aed); border-color: var(--md-primary, #7c3aed); color: #fff; }
.tt-tab.today { border-color: var(--md-primary, #7c3aed); color: var(--md-primary, #7c3aed); }
.tt-tab.today.active { color: #fff; }
.tt-list { display: flex; flex-direction: column; gap: 10px; }
.tt-row { display: flex; align-items: center; gap: 18px; background: var(--color-bg-card, #fff); border: 2px solid transparent; border-radius: 14px; padding: 14px 22px; box-shadow: 0 1px 3px rgba(0,0,0,.05); }
.tt-row.current { border-color: var(--md-primary, #7c3aed); box-shadow: 0 4px 14px rgba(124,58,237,.18); }
.tt-time { display: flex; flex-direction: column; width: 96px; flex-shrink: 0; }
.tt-time-range { font-size: 15px; font-weight: 700; font-variant-numeric: tabular-nums; }
.tt-period-name { font-size: 12px; color: var(--color-text-secondary); }
.tt-main { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 0; }
.tt-subject { font-size: 22px; font-weight: 800; }
.tt-week { font-size: 12px; font-weight: 600; padding: 2px 8px; border-radius: 9999px; background: var(--c-violet-chip); color: var(--color-primary); }
.tt-meta { font-size: 14px; color: var(--color-text-secondary); flex-shrink: 0; }
.tt-empty { text-align: center; padding: 60px 0; color: var(--color-text-secondary); font-size: 15px; }
.tt-empty-icon { display: block; font-size: 40px; margin-bottom: 10px; }
.tt-error { color: var(--color-danger, #ef4444); }
</style>
