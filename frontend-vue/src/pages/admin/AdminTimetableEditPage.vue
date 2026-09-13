<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import api, { apiGet, apiPost, apiPut } from '@/utils/api'
import type {
  ApiResponse, ClassRoom, AdminTimetableData, TimetableSubject, TimetablePeriod, TimetableEntry,
  TimetableWeekType, TimetableTeacherAssignment, TimetableGenerateRules, TimetableSubjectRule,
  TimetableGenerateResult, TimetableSchoolGenerateResult, TimetableImportSummary,
} from '@/types'

// ===== 班级选择 =====
const classes = ref<ClassRoom[]>([])
const currentClassId = ref<number>(0)
const currentClass = computed(() => classes.value.find(c => c.id === currentClassId.value))

const loading = ref(true)
const loadError = ref('')
const subjects = ref<TimetableSubject[]>([])
const periods = ref<TimetablePeriod[]>([])
const entries = ref<TimetableEntry[]>([])

const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const saveMessage = ref('')

const WEEKDAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const WEEK_TYPE_LABELS: Record<TimetableWeekType, string> = { all: '每周', odd: '单周', even: '双周' }
const DEFAULT_SUBJECTS = ['语文', '数学', '英语', '科学', '道德与法治', '体育', '音乐', '美术', '信息', '班会']
const DEFAULT_PERIODS: Array<{ start_time: string; end_time: string }> = [
  { start_time: '08:00', end_time: '08:40' }, { start_time: '08:50', end_time: '09:30' },
  { start_time: '09:50', end_time: '10:30' }, { start_time: '10:40', end_time: '11:20' },
  { start_time: '14:00', end_time: '14:40' }, { start_time: '14:50', end_time: '15:30' },
  { start_time: '15:50', end_time: '16:30' }, { start_time: '16:40', end_time: '17:20' },
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

// ===== 任课设置 =====
const assignments = ref<TimetableTeacherAssignment[]>([])
const assignmentStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

// ===== 自动生成 =====
type GenRule = TimetableSubjectRule & { forbidText: string }

function toGenRule(name: string, existing?: TimetableSubjectRule): GenRule {
  const base: GenRule = { name, weekly: 0, double: false, session: 'any', max_per_day: 1, forbid_periods: [], forbidText: '' }
  if (existing) {
    return { ...existing, forbidText: (existing.forbid_periods || []).join(',') }
  }
  return base
}

function parseForbid(r: GenRule) {
  r.forbid_periods = r.forbidText
    .split(/[,，\s]+/)
    .map(s => parseInt(s, 10))
    .filter(n => Number.isFinite(n) && n >= 1 && n <= 30)
}

const showGenModal = ref(false)
const genSchoolWide = ref(false)
const genDays = ref<number[]>([1, 2, 3, 4, 5])
const genRules = ref<GenRule[]>([])
const genLoading = ref(false)
const genErrors = ref<string[]>([])
const genNotice = ref('')

// ===== 批量导入 =====
const importFile = ref<File | null>(null)
const importSummary = ref<TimetableImportSummary | null>(null)
const importLoading = ref(false)
const importError = ref('')

onMounted(async () => {
  loading.value = true
  try {
    const clsRes = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = clsRes.data || []
    if (classes.value.length > 0) {
      currentClassId.value = classes.value[0].id
      await loadTimetable()
    }
  } catch {
    loadError.value = '班级列表加载失败，请刷新重试'
  } finally {
    loading.value = false
  }
})

async function loadTimetable() {
  if (!currentClassId.value) return
  loadError.value = ''
  try {
    const res = await apiGet<ApiResponse<AdminTimetableData>>(`/api/v1/admin/classes/${currentClassId.value}/timetable`)
    subjects.value = res.data?.subjects || []
    periods.value = res.data?.periods || []
    entries.value = res.data?.entries || []
    const aRes = await apiGet<ApiResponse<TimetableTeacherAssignment[]>>(`/api/v1/admin/classes/${currentClassId.value}/teacher-assignments`)
    assignments.value = aRes.data || []
  } catch {
    loadError.value = '课表加载失败，请刷新重试'
  }
}

function onClassChange() {
  entries.value = []
  loadTimetable()
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
  addTeacher.value = assignments.value.find(a => a.subject_name === addSubjectName.value)?.teacher_name || ''
  addRoom.value = ''
  cellError.value = ''
}

function onModalSubjectChange() {
  const known = assignments.value.find(a => a.subject_name === addSubjectName.value)
  if (known) addTeacher.value = known.teacher_name
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

// ===== 保存（即时生效） =====
async function saveTimetable() {
  if (!currentClassId.value) return
  saveStatus.value = 'loading'
  saveMessage.value = ''
  try {
    const res = await apiPost<{ message: string }>(`/api/v1/admin/classes/${currentClassId.value}/timetable`, {
      subjects: subjects.value,
      periods: periods.value,
      entries: entries.value,
    })
    saveStatus.value = 'success'
    saveMessage.value = res.message || '课表已保存并即时生效'
    setTimeout(() => { saveStatus.value = 'idle'; saveMessage.value = '' }, 2500)
    await loadTimetable()
  } catch {
    saveStatus.value = 'error'
    setTimeout(() => { saveStatus.value = 'idle' }, 3000)
  }
}

// ===== 任课设置 =====
function addAssignment() {
  assignments.value.push({ subject_name: '', teacher_name: '' })
}

function removeAssignment(index: number) {
  assignments.value.splice(index, 1)
}

async function saveAssignments() {
  if (!currentClassId.value) return
  assignmentStatus.value = 'loading'
  try {
    await apiPut(`/api/v1/admin/classes/${currentClassId.value}/teacher-assignments`, {
      assignments: assignments.value.filter(a => a.subject_name.trim() && a.teacher_name.trim()),
    })
    assignmentStatus.value = 'success'
    setTimeout(() => { assignmentStatus.value = 'idle' }, 1500)
    await loadTimetable()
  } catch {
    assignmentStatus.value = 'error'
    setTimeout(() => { assignmentStatus.value = 'idle' }, 3000)
  }
}

// ===== 自动生成 =====
function openGenModal() {
  genErrors.value = []
  genNotice.value = ''
  genRules.value = subjects.value.map(s => toGenRule(s.name, genRules.value.find(r => r.name === s.name)))
  showGenModal.value = true
}

function toggleGenDay(day: number) {
  const i = genDays.value.indexOf(day)
  if (i >= 0) genDays.value.splice(i, 1)
  else genDays.value.push(day)
}

async function runGenerate() {
  genErrors.value = []
  genNotice.value = ''
  const active = genRules.value.filter(r => r.weekly > 0)
  if (active.length === 0) { genErrors.value = ['请为至少一个科目设置每周节数']; return }
  if (genDays.value.length === 0) { genErrors.value = ['请至少选择一个上课日']; return }

  genLoading.value = true
  try {
    const payload: TimetableGenerateRules = { days: [...genDays.value].sort(), subjects: active }

    if (genSchoolWide.value) {
      const res = await apiPost<ApiResponse<TimetableSchoolGenerateResult>>('/api/v1/admin/timetable/generate-school', { rules: payload, commit: true })
      const data = res.data
      if (!data?.success) {
        genErrors.value = data?.warnings || ['全校排课失败']
        return
      }
      genNotice.value = `全校排课完成：${data.classes.length} 个班级共 ${data.classes.reduce((sum, c) => sum + c.entry_count, 0)} 节课${data.warnings.length ? '；提示：' + data.warnings.join('；') : ''}`
    } else {
      const res = await apiPost<ApiResponse<TimetableGenerateResult>>('/api/v1/admin/timetable/generate', { rules: payload })
      const data = res.data
      if (!data?.success) {
        genErrors.value = data?.warnings || ['生成失败']
        return
      }
      entries.value = data.entries
      genNotice.value = '已生成课表预览，请检查后点击「保存课表」生效'
    }
    showGenModal.value = false
    await loadTimetable()
  } catch {
    genErrors.value = ['请求失败，请重试']
  } finally {
    genLoading.value = false
  }
}

// ===== 批量导入 =====
function onImportFile(e: Event) {
  const files = (e.target as HTMLInputElement).files
  importFile.value = files && files.length > 0 ? files[0] : null
  importSummary.value = null
  importError.value = ''
}

function downloadTemplate() {
  const header = '年级,班级,星期,第几节,开始时间,结束时间,科目,周次,教师,教室'
  const samples = [
    '一年级,1班,1,1,08:00,08:40,语文,每周,张老师,101',
    '一年级,1班,周一,2,08:50,09:30,数学,单周,李老师,101',
    '二年级,2班,星期三,5,14:00,14:40,体育,每周,刘老师,操场',
  ].join('\n')
  // BOM 保证 Excel 打开中文不乱码
  const blob = new Blob(['\uFEFF' + header + '\n' + samples + '\n'], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '课表批量导入模板.csv'
  link.click()
  URL.revokeObjectURL(url)
}

async function submitImport(dryRun: boolean) {
  if (!importFile.value) { importError.value = '请先选择 CSV 文件'; return }
  importLoading.value = true
  importError.value = ''
  try {
    const form = new FormData()
    form.append('file', importFile.value)
    form.append('dry_run', String(dryRun))
    const res = await api.post<TimetableImportSummary>('/api/v1/admin/timetable/import-csv', form)
    importSummary.value = res.data
    if (!dryRun) await loadTimetable()
  } catch {
    importError.value = '导入失败，请检查文件格式'
  } finally {
    importLoading.value = false
  }
}
</script>

<template>
  <div>
    <div class="page-head">
      <h2 class="page-title">课表管理</h2>
      <div class="head-actions">
        <select v-model="currentClassId" class="form-select class-select" @change="onClassChange">
          <option v-for="c in classes" :key="c.id" :value="c.id">{{ [c.grade, c.name].filter(Boolean).join(' ') }}</option>
        </select>
        <button class="btn btn-sm btn-ghost" @click="openGenModal">⚙️ 自动排课</button>
        <button class="btn btn-sm btn-primary" :class="{ 'btn-state-loading': saveStatus === 'loading', 'btn-state-success': saveStatus === 'success', 'btn-state-error': saveStatus === 'error' }" :disabled="saveStatus === 'loading'" @click="saveTimetable">
          {{ { idle: '保存课表', loading: '保存中...', success: '已保存 ✓', error: '保存失败' }[saveStatus] }}
        </button>
      </div>
    </div>
    <div v-if="saveMessage" class="save-tip">{{ saveMessage }}</div>
    <div v-else-if="saveStatus === 'idle'" class="save-tip save-tip--pending">管理员修改保存后即时生效；教师提交的修改仍需在此审核。自动排课见「自动排课」按钮与下方批量导入。</div>

    <div v-if="loading" class="empty-state">加载中...</div>
    <div v-else-if="loadError" class="error-banner">{{ loadError }}</div>

    <template v-else-if="currentClassId">
      <div class="config-grid">
        <div class="card">
          <div class="card-head">
            <h3 class="card-title">科目（全校共享）</h3>
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
            <h3 class="card-title">节次作息（全校共享）</h3>
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
          <h3 class="card-title">{{ currentClass?.name ? currentClass.name + ' 排课' : '排课' }}</h3>
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
        <p class="muted-tip">点击格子添加 / 删除课程；同一节可为单周 / 双周排不同科目。保存后即时生效，并可导出 CSES 供 ClassIsland 同步。</p>
      </div>

      <div class="card">
        <div class="card-head">
          <h3 class="card-title">任课设置（全校智能排课的依据）</h3>
          <div class="head-actions">
            <button class="btn btn-sm btn-ghost" @click="addAssignment">加一行</button>
            <button class="btn btn-sm btn-primary" :class="{ 'btn-state-loading': assignmentStatus === 'loading', 'btn-state-success': assignmentStatus === 'success', 'btn-state-error': assignmentStatus === 'error' }" :disabled="assignmentStatus === 'loading'" @click="saveAssignments">
              {{ { idle: '保存任课', loading: '保存中...', success: '已保存 ✓', error: '保存失败' }[assignmentStatus] }}
            </button>
          </div>
        </div>
        <div class="assign-list">
          <div v-for="(a, i) in assignments" :key="i" class="assign-row">
            <select v-model="a.subject_name" class="form-select assign-subject">
              <option value="" disabled>科目</option>
              <option v-for="s in subjects" :key="s.name" :value="s.name">{{ s.name }}</option>
            </select>
            <span class="assign-arrow">→</span>
            <input v-model="a.teacher_name" class="form-input" placeholder="教师姓名" maxlength="50" />
            <button class="chip-remove" title="删除" @click="removeAssignment(i)">&times;</button>
          </div>
          <div v-if="assignments.length === 0" class="text-muted-13">尚未登记任课。填写后可在「自动排课」中使用全校智能排课（教师不冲突）。</div>
        </div>
      </div>

      <div class="card">
        <div class="card-head">
          <h3 class="card-title">批量导入（CSV）</h3>
          <div class="head-actions">
            <button class="btn btn-sm btn-ghost" @click="downloadTemplate">下载模板</button>
          </div>
        </div>
        <p class="muted-tip">模板列：年级,班级,星期,第几节,开始时间,结束时间,科目,周次,教师,教室。星期支持 1-7 / 周一 / 星期一；周次支持 每周 / 单周 / 双周。某班出现在文件中即整体覆盖该班课表。</p>
        <div class="import-row">
          <input type="file" accept=".csv,.txt" class="form-input import-file" @change="onImportFile" />
          <button class="btn btn-sm btn-ghost" :disabled="importLoading || !importFile" @click="submitImport(true)">预览检查</button>
          <button class="btn btn-sm btn-primary" :disabled="importLoading || !importFile" @click="submitImport(false)">确认导入</button>
        </div>
        <div v-if="importError" class="field-error">{{ importError }}</div>
        <div v-if="importSummary" class="import-summary">
          <div class="text-muted-13">共 {{ importSummary.total_rows }} 行 · 覆盖 {{ importSummary.classes.filter(c => c.found).length }} 个班级 · {{ importSummary.period_count }} 个节次{{ importSummary.imported ? ' · 已导入' : ' · 预览（未落库）' }}</div>
          <div v-for="(c, i) in importSummary.classes" :key="i" class="import-class-row">
            <span>{{ c.grade }} {{ c.name }}</span>
            <span :class="c.found ? 'fw-600' : 'text-muted-13'">{{ c.found ? `${c.entry_count} 节` : '班级不存在，跳过' }}</span>
          </div>
          <div v-for="(err, i) in importSummary.errors" :key="'e' + i" class="field-error">{{ err }}</div>
        </div>
      </div>
    </template>

    <!-- 自动排课弹窗 -->
    <div v-if="showGenModal" class="modal-overlay" @click.self="showGenModal = false">
      <div class="modal-card modal-card--wide">
        <div class="modal-header">
          <h3>自动排课规则</h3>
          <button class="modal-close" @click="showGenModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <label class="schoolwide-toggle">
            <input v-model="genSchoolWide" type="checkbox" />
            应用到全校所有班级（依据各班任课设置，教师同一时段不冲突）
          </label>
          <div class="gen-days-row">
            <span class="form-label">上课日：</span>
            <label v-for="(label, i) in WEEKDAY_LABELS" :key="i" class="day-check">
              <input type="checkbox" :checked="genDays.includes(i + 1)" @change="toggleGenDay(i + 1)" />
              {{ label }}
            </label>
          </div>
          <table class="gen-table">
            <thead>
              <tr>
                <th>科目</th><th>每周节数</th><th>连堂</th><th>每日上限</th><th>时段</th><th>禁排节次</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in genRules" :key="r.name">
                <td class="fw-600">{{ r.name }}</td>
                <td><input v-model.number="r.weekly" type="number" min="0" max="35" class="form-input num-input" /></td>
                <td><input v-model="r.double" type="checkbox" /></td>
                <td>
                  <select v-model.number="r.max_per_day" class="form-select num-select">
                    <option :value="1">1</option><option :value="2">2</option><option :value="3">3</option>
                  </select>
                </td>
                <td>
                  <select v-model="r.session" class="form-select num-select">
                    <option value="any">不限</option><option value="am">仅上午</option><option value="pm">仅下午</option>
                  </select>
                </td>
                <td><input v-model="r.forbidText" placeholder="如 1,8" class="form-input num-input" @input="parseForbid(r)" /></td>
              </tr>
            </tbody>
          </table>
          <div v-if="genErrors.length > 0">
            <div v-for="(e, i) in genErrors" :key="i" class="field-error">{{ e }}</div>
          </div>
        </div>
        <div class="modal-footer">
          <span class="text-muted-13">{{ genSchoolWide ? '全校排课确认后直接生效' : '生成后先预览，点「保存课表」生效' }}</span>
          <div class="head-actions">
            <button class="btn btn-sm btn-ghost" @click="showGenModal = false">取消</button>
            <button class="btn btn-sm btn-primary" :disabled="genLoading" @click="runGenerate">{{ genLoading ? '排课中...' : '开始排课' }}</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑格子模态 -->
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
                <select v-model="addSubjectName" class="form-select" @change="onModalSubjectChange">
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
.class-select { width: 180px }
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
.assign-list { display: flex; flex-direction: column; gap: 8px }
.assign-row { display: flex; align-items: center; gap: 8px }
.assign-subject { width: 160px }
.assign-arrow { color: #86868b }
.assign-row .form-input { width: 160px }
.assign-row .chip-remove { margin-left: 8px; font-size: 17px }
.import-row { display: flex; gap: 8px; align-items: center; margin-top: 10px }
.import-file { max-width: 320px }
.import-summary { margin-top: 12px; display: flex; flex-direction: column; gap: 4px; font-size: 13px }
.import-class-row { display: flex; justify-content: space-between; max-width: 420px; padding: 3px 0; border-bottom: 1px dashed var(--color-border, #f0f0f3) }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,.4); display: flex; align-items: center; justify-content: center; z-index: 1000 }
.modal-card { background: var(--color-bg-card, #fff); border-radius: 16px; width: 100%; max-width: 480px; box-shadow: 0 20px 60px rgba(0,0,0,.15); overflow: hidden }
.modal-card--wide { max-width: 720px }
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
.schoolwide-toggle { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 600; margin-bottom: 12px; cursor: pointer }
.gen-days-row { display: flex; gap: 12px; align-items: center; margin-bottom: 12px; flex-wrap: wrap }
.day-check { display: flex; align-items: center; gap: 4px; font-size: 13px; cursor: pointer }
.gen-table { width: 100%; border-collapse: collapse; font-size: 13px }
.gen-table th { text-align: left; font-weight: 600; color: var(--color-text-secondary, #86868b); padding: 6px 8px; border-bottom: 1px solid var(--color-border, #e5e5ea) }
.gen-table td { padding: 6px 8px; border-bottom: 1px dashed var(--color-border, #f0f0f3) }
.num-input { width: 72px; padding: 5px 8px }
.num-select { width: 92px; padding: 5px 8px }
.field-error { color: #ef4444; font-size: 12px; margin-top: 6px }
</style>
