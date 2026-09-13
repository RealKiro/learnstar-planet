<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import api, { apiGet, apiPost } from '@/utils/api'
import type {
  ApiResponse, TimetableData, TimetableSubject, TimetablePeriod, TimetableEntry, TimetableWeekType,
  TimetableChangeRequest,
} from '@/types'

const loading = ref(true)
const loadError = ref('')
const subjects = ref<TimetableSubject[]>([])
const periods = ref<TimetablePeriod[]>([])
const entries = ref<TimetableEntry[]>([])
const changes = ref<TimetableChangeRequest[]>([])

const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const saveMessage = ref('')
const exportStatus = ref<'idle' | 'loading' | 'error'>('idle')

const pendingCount = computed(() => changes.value.filter(c => c.status === 'pending').length)

const WEEKDAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const WEEK_TYPE_LABELS: Record<TimetableWeekType, string> = { all: '每周', odd: '单周', even: '双周' }

const DEFAULT_SUBJECTS = ['语文', '数学', '英语', '科学', '道德与法治', '体育', '音乐', '美术', '信息', '班会']
const DEFAULT_PERIODS: Array<{ start_time: string; end_time: string }> = [
  { start_time: '08:00', end_time: '08:40' },
  { start_time: '08:50', end_time: '09:30' },
  { start_time: '09:50', end_time: '10:30' },
  { start_time: '10:40', end_time: '11:20' },
  { start_time: '14:00', end_time: '14:40' },
  { start_time: '14:50', end_time: '15:30' },
  { start_time: '15:50', end_time: '16:30' },
  { start_time: '16:40', end_time: '17:20' },
]

const showWeekend = ref(false)
const visibleWeekdays = computed(() => (showWeekend.value ? [1, 2, 3, 4, 5, 6, 7] : [1, 2, 3, 4, 5]))

const newSubjectName = ref('')
const cellError = ref('')

// ===== 编辑格子模态 =====
const editingCell = ref<{ weekday: number; periodIndex: number } | null>(null)
const addSubjectName = ref('')
const addWeekType = ref<TimetableWeekType>('all')
const addTeacher = ref('')
const addRoom = ref('')

const cellEntries = computed(() => {
  const cell = editingCell.value
  if (!cell) return []
  return entries.value.filter(e => e.weekday === cell.weekday && e.period_index === cell.periodIndex)
})

onMounted(loadData)

async function loadData() {
  loading.value = true
  loadError.value = ''
  try {
    const res = await apiGet<ApiResponse<TimetableData>>('/api/v1/teacher/timetable')
    subjects.value = res.data?.subjects || []
    periods.value = res.data?.periods || []
    entries.value = res.data?.entries || []
    const chRes = await apiGet<ApiResponse<TimetableChangeRequest[]>>('/api/v1/teacher/timetable/changes')
    changes.value = chRes.data || []
  } catch {
    loadError.value = '课表加载失败，请刷新重试'
  } finally {
    loading.value = false
  }
}

function entriesAt(weekday: number, periodIndex: number): TimetableEntry[] {
  return entries.value.filter(e => e.weekday === weekday && e.period_index === periodIndex)
}

function subjectColor(name: string): string | null {
  return subjects.value.find(s => s.name === name)?.color || null
}

function openCell(weekday: number, periodIndex: number) {
  editingCell.value = { weekday, periodIndex }
  addSubjectName.value = subjects.value[0]?.name || ''
  addWeekType.value = 'all'
  addTeacher.value = ''
  addRoom.value = ''
  cellError.value = ''
}

function confirmAddEntry() {
  const cell = editingCell.value
  const name = addSubjectName.value.trim()
  if (!cell || !name) { cellError.value = '请选择科目'; return }
  if (cellEntries.value.some(e => e.week_type === addWeekType.value)) {
    cellError.value = `该时段已有${WEEK_TYPE_LABELS[addWeekType.value]}课程，请先删除或选择其他周次`
    return
  }
  entries.value.push({
    weekday: cell.weekday,
    period_index: cell.periodIndex,
    week_type: addWeekType.value,
    subject_name: name,
    teacher_name: addTeacher.value.trim() || null,
    room: addRoom.value.trim() || null,
  })
  addWeekType.value = 'all'
  addTeacher.value = ''
  addRoom.value = ''
  cellError.value = ''
}

function removeEntry(target: TimetableEntry) {
  entries.value = entries.value.filter(e => e !== target)
}

function formatTime(iso?: string | null): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return isNaN(d.getTime()) ? '-' : d.toLocaleString('zh-CN', { month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

// ===== 科目管理 =====
function addSubject() {
  const name = newSubjectName.value.trim()
  if (!name) return
  if (subjects.value.some(s => s.name === name)) { newSubjectName.value = ''; return }
  subjects.value.push({ name })
  newSubjectName.value = ''
}

function removeSubject(index: number) {
  const removed = subjects.value[index]?.name
  subjects.value.splice(index, 1)
  if (removed) entries.value = entries.value.filter(e => e.subject_name !== removed)
}

function fillDefaultSubjects() {
  for (const name of DEFAULT_SUBJECTS) {
    if (!subjects.value.some(s => s.name === name)) subjects.value.push({ name })
  }
}

// ===== 节次管理 =====
function addPeriod() {
  const nextIndex = (periods.value[periods.value.length - 1]?.period_index || 0) + 1
  periods.value.push({ period_index: nextIndex, start_time: '08:00', end_time: '08:40' })
}

function removePeriod(index: number) {
  const removed = periods.value[index]?.period_index
  periods.value.splice(index, 1)
  if (removed) entries.value = entries.value.filter(e => e.period_index !== removed)
}

function fillDefaultPeriods() {
  for (let i = 0; i < DEFAULT_PERIODS.length; i++) {
    const index = i + 1
    if (!periods.value.some(p => p.period_index === index)) {
      periods.value.push({ period_index: index, start_time: DEFAULT_PERIODS[i].start_time, end_time: DEFAULT_PERIODS[i].end_time })
    }
  }
  periods.value.sort((a, b) => a.period_index - b.period_index)
}

// ===== 提交审核 / 导出 =====
async function saveTimetable() {
  saveStatus.value = 'loading'
  saveMessage.value = ''
  try {
    const res = await apiPost<{ message: string }>('/api/v1/teacher/timetable', {
      subjects: subjects.value,
      periods: periods.value,
      entries: entries.value,
    })
    saveStatus.value = 'success'
    saveMessage.value = res.message || '修改申请已提交，待管理员审核'
    setTimeout(() => { saveStatus.value = 'idle'; saveMessage.value = '' }, 2500)
    await loadData()
  } catch {
    saveStatus.value = 'error'
    setTimeout(() => { saveStatus.value = 'idle' }, 3000)
  }
}

async function exportCses() {
  exportStatus.value = 'loading'
  try {
    const res = await api.get<Blob>('/api/v1/teacher/timetable/export-cses', { responseType: 'blob' })
    const url = URL.createObjectURL(res.data)
    const link = document.createElement('a')
    link.href = url
    link.download = '课表.cses.yaml'
    link.click()
    URL.revokeObjectURL(url)
    exportStatus.value = 'idle'
  } catch {
    exportStatus.value = 'error'
    setTimeout(() => { exportStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div class="page-head">
      <h2 class="page-title">课表管理</h2>
      <div class="head-actions">
        <router-link :to="{ name: 'teacher-my-timetable' }" class="btn btn-sm btn-ghost">📆 我的课表</router-link>
        <button class="btn btn-sm btn-ghost" :disabled="exportStatus === 'loading'" @click="exportCses">
          {{ exportStatus === 'loading' ? '导出中...' : exportStatus === 'error' ? '导出失败' : '导出 CSES (ClassIsland)' }}
        </button>
        <button class="btn btn-sm btn-primary" :class="{ 'btn-state-loading': saveStatus === 'loading', 'btn-state-success': saveStatus === 'success', 'btn-state-error': saveStatus === 'error' }" :disabled="saveStatus === 'loading'" @click="saveTimetable">
          {{ { idle: '提交审核', loading: '提交中...', success: '已提交 ✓', error: '提交失败' }[saveStatus] }}
        </button>
      </div>
    </div>
    <div v-if="saveMessage" class="save-tip">{{ saveMessage }}</div>
    <div v-else-if="pendingCount > 0" class="save-tip save-tip--pending">有 {{ pendingCount }} 条修改申请待管理员审核，审核通过前课表保持现状</div>

    <div v-if="loading" class="empty-state">加载中...</div>
    <div v-else-if="loadError" class="error-banner">{{ loadError }}</div>

    <template v-else>
      <div class="config-grid">
        <div class="card">
          <div class="card-head">
            <h3 class="card-title">科目</h3>
            <div class="head-actions">
              <button class="btn btn-sm btn-ghost" @click="fillDefaultSubjects">填入常用科目</button>
            </div>
          </div>
          <div class="chip-row">
            <span v-for="(s, i) in subjects" :key="s.name" class="chip" :style="subjectColor(s.name) ? { borderColor: subjectColor(s.name) || undefined } : {}">
              {{ s.simplified_name || s.name }}
              <button class="chip-remove" :title="`删除 ${s.name}`" @click="removeSubject(i)">&times;</button>
            </span>
            <span v-if="subjects.length === 0" class="text-muted-13">暂无科目，可一键填入</span>
          </div>
          <div class="subject-add-row">
            <input v-model="newSubjectName" class="form-input" placeholder="新科目名称，如：书法" maxlength="50" @keyup.enter="addSubject" />
            <button class="btn btn-sm btn-ghost" :disabled="!newSubjectName.trim()" @click="addSubject">添加</button>
          </div>
        </div>

        <div class="card">
          <div class="card-head">
            <h3 class="card-title">节次作息</h3>
            <div class="head-actions">
              <button class="btn btn-sm btn-ghost" @click="fillDefaultPeriods">填入默认作息</button>
              <button class="btn btn-sm btn-ghost" @click="addPeriod">加一节</button>
            </div>
          </div>
          <div class="period-list">
            <div v-for="(p, i) in periods" :key="p.period_index" class="period-row">
              <span class="period-index">第{{ p.period_index }}节</span>
              <input v-model="p.start_time" type="time" class="form-input time-input" />
              <span class="period-dash">–</span>
              <input v-model="p.end_time" type="time" class="form-input time-input" />
              <button class="chip-remove period-remove" title="删除该节" @click="removePeriod(i)">&times;</button>
            </div>
            <div v-if="periods.length === 0" class="text-muted-13">暂无节次，可一键填入默认作息</div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-head">
          <h3 class="card-title">排课</h3>
          <label class="weekend-toggle">
            <input v-model="showWeekend" type="checkbox" />
            显示周末
          </label>
        </div>
        <div v-if="periods.length === 0" class="empty-state">先在上方添加节次，再进行排课</div>
        <div v-else class="table-scroll">
          <table class="grid-table">
            <thead>
              <tr>
                <th class="corner-th">节次</th>
                <th v-for="d in visibleWeekdays" :key="d">{{ WEEKDAY_LABELS[d - 1] }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in periods" :key="p.period_index">
                <td class="period-cell">
                  <div class="fw-600">第{{ p.period_index }}节</div>
                  <div class="time-hint">{{ p.start_time }}–{{ p.end_time }}</div>
                </td>
                <td v-for="d in visibleWeekdays" :key="d" class="slot-cell" @click="openCell(d, p.period_index)">
                  <div v-for="e in entriesAt(d, p.period_index)" :key="e.week_type" class="slot-entry" :style="subjectColor(e.subject_name) ? { borderColor: subjectColor(e.subject_name) || undefined } : {}">
                    <span class="slot-subject">{{ e.subject_name }}</span>
                    <span v-if="e.week_type !== 'all'" class="slot-week">{{ WEEK_TYPE_LABELS[e.week_type] }}</span>
                    <span v-if="e.teacher_name || e.room" class="slot-meta">{{ [e.teacher_name, e.room].filter(Boolean).join(' · ') }}</span>
                  </div>
                  <span v-if="entriesAt(d, p.period_index).length === 0" class="slot-plus">+</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <p class="muted-tip">点击格子添加课程；同一节可为单周 / 双周排不同科目。保存需提交审核，管理员通过后课表才会更新。导出的 CSES 文件可在 ClassIsland 中「从 CSES 导入」同步课表。</p>
      </div>

      <div v-if="changes.length > 0" class="card">
        <div class="card-head"><h3 class="card-title">修改申请记录</h3></div>
        <div class="change-list">
          <div v-for="c in changes" :key="c.id" class="change-row">
            <span :class="['status-badge', c.status]">{{ { pending: '待审核', approved: '已通过', rejected: '已驳回' }[c.status] }}</span>
            <span class="text-muted-13">{{ c.entry_count }} 节排课</span>
            <span class="text-muted-13">{{ formatTime(c.created_at) }} 提交</span>
            <span v-if="c.status !== 'pending' && c.reviewer_name" class="text-muted-13">{{ c.reviewer_name }} 审核</span>
            <span v-if="c.review_note" class="change-note" :title="c.review_note">{{ c.review_note }}</span>
          </div>
        </div>
      </div>
    </template>

    <div v-if="editingCell" class="modal-overlay" @click.self="editingCell = null">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ WEEKDAY_LABELS[editingCell.weekday - 1] }} · 第{{ editingCell.periodIndex }}节</h3>
          <button class="modal-close" @click="editingCell = null">&times;</button>
        </div>
        <div class="modal-body">
          <div v-if="cellEntries.length > 0" class="cell-entry-list">
            <div v-for="e in cellEntries" :key="e.week_type" class="cell-entry">
              <span class="fw-600">{{ e.subject_name }}</span>
              <span class="week-badge">{{ WEEK_TYPE_LABELS[e.week_type] }}</span>
              <span v-if="e.teacher_name || e.room" class="text-muted-13">{{ [e.teacher_name, e.room].filter(Boolean).join(' · ') }}</span>
              <button class="chip-remove" title="删除" @click="removeEntry(e)">&times;</button>
            </div>
          </div>
          <div v-else class="text-muted-13 modal-hint-line">该时段暂无课程</div>

          <div class="add-entry-form">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">科目 <span class="req-star-red">*</span></label>
                <select v-model="addSubjectName" class="form-select">
                  <option value="" disabled>选择科目</option>
                  <option v-for="s in subjects" :key="s.name" :value="s.name">{{ s.name }}</option>
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">周次</label>
                <select v-model="addWeekType" class="form-select">
                  <option value="all">每周</option>
                  <option value="odd">单周</option>
                  <option value="even">双周</option>
                </select>
              </div>
            </div>
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">教师（可选）</label>
                <input v-model="addTeacher" class="form-input" placeholder="如：张老师" maxlength="50" />
              </div>
              <div class="form-group">
                <label class="form-label">教室（可选）</label>
                <input v-model="addRoom" class="form-input" placeholder="如：301" maxlength="50" />
              </div>
            </div>
            <div v-if="cellError" class="field-error">{{ cellError }}</div>
            <button class="btn btn-sm btn-primary" @click="confirmAddEntry">添加到该时段</button>
          </div>
        </div>
        <div class="modal-footer">
          <span class="text-muted-13">添加 / 删除后记得点击「保存课表」</span>
          <button class="btn btn-sm btn-ghost" @click="editingCell = null">完成</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.head-actions { display: flex; gap: 8px; align-items: center }
.config-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 16px }
.card-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px }
.card-title { font-size: 16px; font-weight: 600; margin: 0 }
.chip-row { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 12px }
.chip { display: inline-flex; align-items: center; gap: 4px; padding: 4px 10px; border: 1.5px solid var(--color-border); border-radius: 9999px; font-size: 13px; font-weight: 500 }
.chip-remove { background: none; border: none; color: #86868b; cursor: pointer; font-size: 15px; line-height: 1; padding: 0 }
.chip-remove:hover { color: var(--color-danger, #ef4444) }
.subject-add-row { display: flex; gap: 8px }
.subject-add-row .form-input { flex: 1 }
.period-list { display: flex; flex-direction: column; gap: 6px; max-height: 260px; overflow-y: auto }
.period-row { display: flex; align-items: center; gap: 8px }
.period-index { width: 56px; font-size: 13px; font-weight: 600; color: var(--color-text-secondary, #86868b) }
.time-input { width: 100px; padding: 5px 8px }
.period-dash { color: #86868b }
.period-remove { font-size: 16px; margin-left: auto }
.weekend-toggle { display: flex; align-items: center; gap: 6px; font-size: 13px; color: var(--color-text-secondary, #86868b); cursor: pointer }
.table-scroll { overflow-x: auto }
.grid-table { width: 100%; border-collapse: separate; border-spacing: 4px; table-layout: fixed }
.grid-table th { font-size: 13px; font-weight: 600; color: var(--color-text-secondary, #86868b); padding: 4px }
.corner-th { width: 92px; text-align: left }
.period-cell { vertical-align: middle; padding: 4px }
.time-hint { font-size: 11px; color: #86868b }
.slot-cell { position: relative; min-height: 56px; height: 56px; background: var(--color-bg, #f7f7f9); border-radius: 8px; cursor: pointer; padding: 3px; vertical-align: top }
.slot-cell:hover { background: var(--color-border, #ececf1) }
.slot-entry { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 48px; border: 1.5px solid var(--color-border); border-radius: 6px; font-size: 12px; line-height: 1.3; overflow: hidden }
.slot-subject { font-weight: 600 }
.slot-week { font-size: 10px; color: #86868b }
.slot-meta { font-size: 10px; color: #86868b; max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap }
.slot-plus { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); color: #c7c7cc; font-size: 16px }
.muted-tip { font-size: 12px; color: #86868b; margin: 10px 2px 0 }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,.4); display: flex; align-items: center; justify-content: center; z-index: 1000 }
.modal-card { background: var(--color-bg-card, #fff); border-radius: 16px; width: 100%; max-width: 480px; box-shadow: 0 20px 60px rgba(0,0,0,.15); overflow: hidden }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0f0f3 }
.modal-header h3 { font-size: 16px; font-weight: 700; margin: 0 }
.modal-close { background: none; border: none; font-size: 22px; color: #86868b; cursor: pointer; padding: 0; line-height: 1 }
.modal-close:hover { color: var(--color-text) }
.modal-body { padding: 20px }
.modal-footer { display: flex; gap: 8px; align-items: center; justify-content: space-between; padding: 16px 20px; border-top: 1px solid #f0f0f3 }
.cell-entry-list { display: flex; flex-direction: column; gap: 6px; margin-bottom: 16px }
.cell-entry { display: flex; align-items: center; gap: 8px; padding: 8px 10px; background: var(--color-bg, #f7f7f9); border-radius: 8px; font-size: 13px }
.cell-entry .chip-remove { margin-left: auto }
.week-badge { font-size: 11px; padding: 1px 6px; border-radius: 9999px; background: #ede9fe; color: #7c3aed }
.modal-hint-line { margin-bottom: 16px }
.add-entry-form { display: flex; flex-direction: column; gap: 10px; border-top: 1px dashed var(--color-border, #e5e5ea); padding-top: 14px }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px }
.form-label { font-size: 13px; font-weight: 600; display: block; margin-bottom: 4px }
.add-entry-form .btn { align-self: flex-start }
.req-star-red { color: #ef4444 }
.fw-600 { font-weight: 600 }
.save-tip { background: #fffbeb; border: 1px solid #fde68a; color: #92400e; font-size: 13px; border-radius: 10px; padding: 8px 14px; margin-bottom: 12px }
.save-tip--pending { background: var(--color-bg, #f7f7f9); border-color: var(--color-border, #e5e5ea); color: var(--color-text-secondary, #86868b) }
.change-list { display: flex; flex-direction: column; gap: 6px }
.change-row { display: flex; align-items: center; gap: 12px; font-size: 13px; padding: 6px 0; border-bottom: 1px dashed var(--color-border, #f0f0f3); flex-wrap: wrap }
.change-row:last-child { border-bottom: none }
.status-badge { font-size: 12px; font-weight: 600; padding: 2px 10px; border-radius: 9999px }
.status-badge.pending { background: #fef3c7; color: #92400e }
.status-badge.approved { background: #d1fae5; color: #065f46 }
.status-badge.rejected { background: #fee2e2; color: #991b1b }
.change-note { color: var(--color-text-secondary, #86868b); max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap }
</style>
