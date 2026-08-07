<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import ModalGlass from '@/components/common/ModalGlass.vue'
import type { ApiResponse, Student, ClassRoom } from '@/types'
import PetSprite from '@/components/pet/PetSprite.vue'
import * as XLSX from 'xlsx'

const students = ref<Student[]>([])
const classes = ref<ClassRoom[]>([])
const loading = ref(true)
const loadError = ref('')
const searchKeyword = ref('')
const meta = ref({ current_page: 1, last_page: 1, per_page: 50, total: 0 })

const pageLabel = computed(() => {
  const { current_page, last_page, total } = meta.value
  return last_page > 1 ? `第 ${current_page} / ${last_page} 页 · 共 ${total} 人` : `共 ${total} 人`
})

const classesById = computed(() => {
  const map = new Map<number, string>()
  for (const c of classes.value) map.set(c.id, c.name)
  return map
})
const classNameOf = (s: Student) =>
  (s as any).class_name || classesById.value.get(s.class_id as number) || (s as any).class_room?.name || '-'

// ===== 导入（CSV + 真 xlsx 解析） =====
const importError = ref('')
const importStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const importSuccessCount = ref(0)
const fileInput = ref<HTMLInputElement | null>(null)

function triggerImport() {
  importError.value = ''
  importStatus.value = 'idle'
  fileInput.value?.click()
}

async function onFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  await handleImportFile(file)
  input.value = '' // 允许再次选择同一文件
}

async function handleImportFile(file: File) {
  importStatus.value = 'loading'
  importError.value = ''
  try {
    const imported = await parseStudentFile(file)
    if (imported.length === 0) {
      importStatus.value = 'error'
      importError.value = '文件为空或格式不正确（表头：姓名,学号,班级）'
      setTimeout(() => { if (importStatus.value === 'error') importStatus.value = 'idle' }, 3000)
      return
    }
    await apiPost('/api/v1/teacher/students/import', { students: imported }, { skipToast: true })
    importSuccessCount.value = imported.length
    importStatus.value = 'success'
    setTimeout(() => { if (importStatus.value === 'success') importStatus.value = 'idle' }, 2000)
    await loadStudents(true)
  } catch (err: any) {
    importStatus.value = 'error'
    importError.value = err?.response?.data?.message || '导入失败，请检查文件格式（表头：姓名,学号,班级）'
    setTimeout(() => { if (importStatus.value === 'error') importStatus.value = 'idle' }, 3000)
  }
}

interface RowData { name: string; student_no: string; class_name: string }

function parseStudentFile(file: File): Promise<RowData[]> {
  return new Promise((resolve, reject) => {
    const ext = file.name.split('.').pop()?.toLowerCase() || ''
    const isCsv = ext === 'csv' || ext === 'txt'
    const reader = new FileReader()
    reader.onload = () => {
      try {
        let rows: RowData[] = []
        if (isCsv) {
          const text = String(reader.result).replace(/^﻿/, '')
          rows = parseCsv(text)
        } else {
          const wb = XLSX.read(reader.result, { type: 'array' })
          const sheet = wb.Sheets[wb.SheetNames[0]]
          const json = XLSX.utils.sheet_to_json<Record<string, unknown>>(sheet, { defval: '' })
          rows = json.map(r => ({
            name: String(r['姓名'] || r['name'] || '').trim(),
            student_no: String(r['学号'] || r['student_no'] || '').trim(),
            class_name: String(r['班级'] || r['class_name'] || '').trim(),
          }))
        }
        resolve(rows.filter(r => r.name))
      } catch {
        reject(new Error('文件解析失败'))
      }
    }
    reader.onerror = () => reject(new Error('文件读取失败'))
    if (isCsv) reader.readAsText(file)
    else reader.readAsArrayBuffer(file)
  })
}

function parseCsv(text: string): RowData[] {
  const lines = text.split(/\r?\n/).filter(l => l.trim())
  if (lines.length === 0) return []
  const hasHeader = /姓名|name/i.test(lines[0])
  const start = hasHeader ? 1 : 0
  const rows: RowData[] = []
  for (let i = start; i < lines.length; i++) {
    const cols = splitCsvLine(lines[i])
    rows.push({ name: (cols[0] || '').trim(), student_no: (cols[1] || '').trim(), class_name: (cols[2] || '').trim() })
  }
  return rows
}

function splitCsvLine(line: string): string[] {
  const out: string[] = []
  let cur = ''
  let inQuotes = false
  for (let i = 0; i < line.length; i++) {
    const ch = line[i]
    if (inQuotes) {
      if (ch === '"') {
        if (line[i + 1] === '"') { cur += '"'; i++ }
        else inQuotes = false
      } else cur += ch
    } else if (ch === '"') {
      inQuotes = true
    } else if (ch === ',') {
      out.push(cur); cur = ''
    } else {
      cur += ch
    }
  }
  out.push(cur)
  return out
}

// ===== 新增 / 编辑 / 删除 =====
const showEditModal = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const formError = ref('')
const submitStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const form = ref({ name: '', gender: '男', student_no: '', class_id: '' as number | '' })
const deleteStatusMap = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})

function openAddModal() {
  isEditing.value = false
  editingId.value = null
  formError.value = ''
  form.value = { name: '', gender: '男', student_no: '', class_id: classes.value[0]?.id || '' }
  showEditModal.value = true
}

function openEditModal(s: Student) {
  isEditing.value = true
  editingId.value = s.id
  formError.value = ''
  form.value = {
    name: s.name,
    gender: s.gender || '男',
    student_no: s.student_no || '',
    class_id: s.class_id as number | '',
  }
  showEditModal.value = true
}

async function submitForm() {
  if (!form.value.name.trim()) { formError.value = '请填写学生姓名'; return }
  if (!isEditing.value && !form.value.class_id) { formError.value = '请选择班级'; return }
  submitStatus.value = 'loading'
  const payload = {
    name: form.value.name.trim(),
    gender: form.value.gender,
    student_no: form.value.student_no.trim(),
  }
  try {
    if (isEditing.value && editingId.value !== null) {
      await apiPut(`/api/v1/teacher/students/${editingId.value}`, payload)
    } else {
      await apiPost('/api/v1/teacher/students', { ...payload, class_id: form.value.class_id })
    }
    submitStatus.value = 'success'
    showEditModal.value = false
    await loadStudents(true)
    setTimeout(() => { submitStatus.value = 'idle' }, 1500)
  } catch {
    submitStatus.value = 'error'
    setTimeout(() => { submitStatus.value = 'idle' }, 3000)
  }
}

async function deleteStudent(s: Student) {
  const ok = await openConfirm({ title: '删除学生', message: `确定删除学生「${s.name}」？`, danger: true, confirmText: '确认删除' })
  if (!ok) return
  deleteStatusMap.value[s.id] = 'loading'
  try {
    await apiDelete(`/api/v1/teacher/students/${s.id}`)
    deleteStatusMap.value[s.id] = 'success'
    await loadStudents()
    setTimeout(() => { delete deleteStatusMap.value[s.id] }, 1500)
  } catch {
    deleteStatusMap.value[s.id] = 'error'
    setTimeout(() => { delete deleteStatusMap.value[s.id] }, 3000)
  }
}

function getStatusLabel(status: string): string {
  const map: Record<string, string> = { active: '活跃', graduated: '已毕业', disabled: '停用' }
  return map[status] || status
}

// ===== 加载 =====
async function loadStudents(resetPage = false) {
  if (resetPage) meta.value.current_page = 1
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.set('per_page', String(meta.value.per_page))
    if (meta.value.current_page > 1) params.set('page', String(meta.value.current_page))
    if (searchKeyword.value.trim()) params.set('search', searchKeyword.value.trim())
    const res = await apiGet<ApiResponse<Student[]>>(`/api/v1/teacher/students?${params.toString()}`)
    students.value = res.data || []
    loadError.value = ''
    if (res.meta) meta.value = res.meta
  } catch {
    students.value = []
    loadError.value = '学生数据加载失败'
  } finally { loading.value = false }
}

async function loadClasses() {
  try {
    const res = await apiGet<{ data: { class_id: number; class_name: string }[] }>('/api/v1/teacher/my-classes')
    classes.value = (res.data || []).map(c => ({ id: c.class_id, name: c.class_name } as ClassRoom))
  } catch { classes.value = [] }
}

function changePage(page: number) {
  const { current_page, last_page } = meta.value
  if (page < 1 || page > last_page || page === current_page) return
  meta.value.current_page = page
  loadStudents()
}

onMounted(async () => {
  await Promise.all([loadClasses(), loadStudents()])
})
</script>

<template>
  <div>
    <div class="page-header">
      <h2 class="page-title">🎒 学生列表</h2>
      <div class="header-actions">
        <input
          v-model="searchKeyword"
          class="form-input search-input"
          placeholder="🔍 搜索姓名 / 学号"
          @keyup.enter="loadStudents(true)"
          @keyup.esc="searchKeyword = ''"
        >
        <button class="btn btn-sm btn-primary" @click="openAddModal">+ 添加学生</button>
        <button
          class="btn btn-sm"
          :class="importStatus === 'loading' ? 'btn-state-loading' : importStatus === 'success' ? 'btn-state-success' : importStatus === 'error' ? 'btn-state-error' : 'btn-solid'"
          :disabled="importStatus === 'loading'"
          @click="triggerImport"
        >
          {{ importStatus === 'loading' ? '导入中...' : importStatus === 'success' ? `✅ 已导入 ${importSuccessCount} 人` : importStatus === 'error' ? '导入失败' : '📥 批量导入' }}
        </button>
        <input ref="fileInput" type="file" accept=".xlsx,.xls,.csv" class="hidden-input" @change="onFileSelected">
      </div>
    </div>

    <div v-if="importError" class="error-banner error-banner--top">{{ importError }}</div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="loadStudents(true)">重试</button>
    </div>

    <div v-else-if="students.length === 0" class="card empty-state">
      <div class="empty-state__icon">📭</div>
      <p>暂无学生数据</p>
    </div>

    <div v-else class="data-table">
      <table>
        <thead><tr><th>姓名</th><th>学号</th><th>班级</th><th>宠物</th><th>积分</th><th>状态</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="s in students" :key="s.id">
            <td class="student-name">{{ s.name }}</td>
            <td class="student-no">{{ s.student_no || '-' }}</td>
            <td class="student-no">{{ classNameOf(s) }}</td>
            <td>
              <div class="pet-cell">
                <div class="pet-avatar">
                  <PetSprite v-if="(s as any).pet_species" :species-id="(s as any).pet_species" :level="(s as any).pet_level || 1" :animate="true" />
                  <span v-else class="pet-egg">🥚</span>
                </div>
                <span v-if="(s as any).pet_level" class="pet-level">Lv.{{ (s as any).pet_level }}</span>
              </div>
            </td>
            <td class="score-text">{{ s.total_score || 0 }}</td>
            <td :class="s.status === 'active' ? 'status-active' : 'status-other'">
              {{ getStatusLabel(s.status) }}
            </td>
            <td>
              <div class="cell-actions">
                <button class="btn btn-sm btn-mini" @click="openEditModal(s)">编辑</button>
                <button class="btn btn-sm" :class="deleteStatusMap[s.id] === 'loading' ? 'btn-state-loading' : deleteStatusMap[s.id] === 'success' ? 'btn-state-success' : deleteStatusMap[s.id] === 'error' ? 'btn-state-error' : 'btn-danger'" :disabled="deleteStatusMap[s.id] === 'loading'" @click="deleteStudent(s)">{{ deleteStatusMap[s.id] === 'loading' ? '删除中...' : deleteStatusMap[s.id] === 'success' ? '已删除' : deleteStatusMap[s.id] === 'error' ? '删除失败' : '删除' }}</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="meta.last_page > 1" class="pagination">
        <button class="btn btn-sm btn-ghost-card" :disabled="meta.current_page <= 1" @click="changePage(meta.current_page - 1)">← 上一页</button>
        <span class="pagination__info">{{ pageLabel }}</span>
        <button class="btn btn-sm btn-ghost-card" :disabled="meta.current_page >= meta.last_page" @click="changePage(meta.current_page + 1)">下一页 →</button>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <ModalGlass :visible="showEditModal" @update:visible="showEditModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">{{ isEditing ? '编辑学生' : '添加学生' }}</h3>
        <button class="modal-close" @click="showEditModal = false">×</button>
      </div>
      <div v-if="!isEditing" class="form-group">
        <label>班级</label>
        <select v-model.number="form.class_id" class="form-select">
          <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
      </div>
      <div class="form-group">
        <label>姓名</label>
        <input v-model="form.name" class="form-input" placeholder="如：张小明" @keydown.enter="submitForm">
      </div>
      <div class="form-group">
        <label>性别</label>
        <select v-model="form.gender" class="form-select">
          <option value="男">男</option>
          <option value="女">女</option>
        </select>
      </div>
      <div class="form-group">
        <label>学号</label>
        <input v-model="form.student_no" class="form-input" placeholder="可选" @keydown.enter="submitForm">
      </div>
      <div class="modal-actions">
        <div v-if="formError" class="error-banner">{{ formError }}</div>
        <button class="btn btn-sm btn-ghost-card" @click="showEditModal = false">取消</button>
        <button class="btn btn-sm" :class="submitStatus === 'loading' ? 'btn-state-loading' : submitStatus === 'success' ? 'btn-state-success' : submitStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="submitStatus === 'loading'" @click="submitForm">{{ submitStatus === 'loading' ? '保存中...' : submitStatus === 'success' ? '已保存 ✓' : submitStatus === 'error' ? '保存失败' : '保存' }}</button>
      </div>
    </ModalGlass>
  </div>
</template>

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px; }
.page-title { font-size: 24px; font-weight: 700; }
.header-actions { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
.search-input { width: 200px; }
.hidden-input { display: none; }

/* ===== 错误横幅（顶置） ===== */
.error-banner--top { margin: -12px 0 16px; }

/* ===== 表格 ===== */
.student-name { font-weight: 600; }
.student-no { color: var(--color-text-secondary); }
.pet-cell { display: flex; align-items: center; gap: 6px; }
.pet-avatar { width: 32px; height: 32px; flex-shrink: 0; }
.pet-egg { font-size: 18px; }
.pet-level { font-size: 11px; color: var(--color-text-secondary); }
.score-text { font-weight: 700; color: var(--color-primary); }
.status-active { color: var(--color-accent); }
.status-other { color: var(--color-text-secondary); }
.cell-actions { display: flex; gap: 4px; }

/* ===== 次要按钮 ===== */
.btn-mini {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  font-size: 12px;
}

/* ===== 分页 ===== */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 14px;
  border-top: 1px solid var(--color-border);
}
.pagination__info { font-size: 13px; color: var(--color-text-secondary); }
.pagination .btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
