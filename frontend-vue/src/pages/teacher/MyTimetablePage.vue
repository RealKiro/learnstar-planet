<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import type {
  ApiResponse, TimetableTeacherSchedule, TimetableWeekType,
} from '@/types'

const loading = ref(true)
const loadError = ref('')
const schedule = ref<TimetableTeacherSchedule | null>(null)
const showWeekend = ref(false)

const WEEKDAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const WEEK_TYPE_LABELS: Record<TimetableWeekType, string> = { all: '每周', odd: '单周', even: '双周' }

const visibleWeekdays = computed(() => (showWeekend.value ? [1, 2, 3, 4, 5, 6, 7] : [1, 2, 3, 4, 5]))
const todayWeekday = computed(() => {
  const d = new Date().getDay()
  return d === 0 ? 7 : d
})

onMounted(loadData)

async function loadData() {
  loading.value = true
  loadError.value = ''
  try {
    const res = await apiGet<ApiResponse<TimetableTeacherSchedule>>('/api/v1/teacher/timetable/my-schedule')
    schedule.value = res.data || null
  } catch {
    loadError.value = '课表加载失败'
  } finally {
    loading.value = false
  }
}

function entriesAt(weekday: number, periodIndex: number) {
  return schedule.value?.entries.filter(e => e.weekday === weekday && e.period_index === periodIndex) || []
}

function subjectColor(name: string): string | null {
  return schedule.value?.subjects.find(s => s.name === name)?.color || null
}
</script>

<template>
  <div>
    <div class="page-head">
      <h2 class="page-title">我的课表</h2>
      <div class="head-actions">
        <span v-if="schedule" class="text-muted-13">{{ schedule.teacher_name }} · 共 {{ schedule.entries.length }} 节课</span>
        <label class="weekend-toggle">
          <input v-model="showWeekend" type="checkbox" />
          显示周末
        </label>
      </div>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>
    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__title">{{ loadError }}</p>
      <p class="error-state__desc">请稍后重试</p>
      <button class="btn btn-sm btn-primary" @click="loadData">重试</button>
    </div>
    <div v-else-if="!schedule || schedule.entries.length === 0" class="empty-state">
      暂未找到您的排课记录（按教师姓名匹配各班排课，可让管理员核对任课设置与排课中的教师姓名）
    </div>

    <div v-else class="card">
      <div class="table-scroll">
        <table class="grid-table">
          <thead>
            <tr>
              <th class="corner-th">节次</th>
              <th v-for="d in visibleWeekdays" :key="d" :class="{ 'today-th': d === todayWeekday }">{{ WEEKDAY_LABELS[d - 1] }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in schedule.periods" :key="p.period_index">
              <td class="period-cell">
                <div class="fw-600">第{{ p.period_index }}节</div>
                <div class="time-hint">{{ p.start_time }}–{{ p.end_time }}</div>
              </td>
              <td v-for="d in visibleWeekdays" :key="d" class="slot-cell" :class="{ 'today-col': d === todayWeekday }">
                <div v-for="e in entriesAt(d, p.period_index)" :key="e.week_type + e.class_id" class="slot-entry" :style="subjectColor(e.subject_name) ? { borderColor: subjectColor(e.subject_name) || undefined } : {}">
                  <span class="slot-subject">{{ e.subject_name }}</span>
                  <span class="slot-class">{{ e.class_name }}</span>
                  <span v-if="e.week_type !== 'all'" class="slot-week">{{ WEEK_TYPE_LABELS[e.week_type] }}</span>
                  <span v-if="e.room" class="slot-meta">{{ e.room }}</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <p class="muted-tip">按教师姓名聚合您在全校各班的排课；同一时段如出现多班即为冲突，请联系管理员调整。编辑班级课表请前往「课表管理」。</p>
    </div>
  </div>
</template>

<style scoped>
.head-actions { display: flex; gap: 10px; align-items: center }
.card-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px }
.weekend-toggle { display: flex; align-items: center; gap: 6px; font-size: 13px; color: var(--color-text-secondary, #86868b); cursor: pointer }
.table-scroll { overflow-x: auto }
.grid-table { width: 100%; border-collapse: separate; border-spacing: 4px; table-layout: fixed }
.grid-table th { font-size: 13px; font-weight: 600; color: var(--color-text-secondary, #86868b); padding: 4px }
.corner-th { width: 92px; text-align: left }
.today-th { color: #007aff }
.period-cell { vertical-align: middle; padding: 4px }
.time-hint { font-size: 11px; color: var(--c-gray-apple) }
.slot-cell { position: relative; min-height: 56px; height: 56px; background: var(--color-bg, #f7f7f9); border-radius: 8px; padding: 3px; vertical-align: top }
.today-col { background: var(--c-blue-tint) }
.slot-entry { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 48px; border: 1.5px solid var(--color-border); border-radius: 6px; font-size: 12px; line-height: 1.3; overflow: hidden; background: var(--color-bg-card, #fff) }
.slot-subject { font-weight: 600 }
.slot-week { font-size: 10px; color: var(--c-gray-apple) }
.slot-meta { font-size: 10px; color: var(--c-gray-apple); max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap }
.slot-class { font-size: 10px; font-weight: 600; color: #007aff }
.muted-tip { font-size: 12px; color: var(--c-gray-apple); margin: 10px 2px 0 }
.fw-600 { font-weight: 600 }
</style>
