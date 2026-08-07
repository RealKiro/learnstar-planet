<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import ModalGlass from '@/components/common/ModalGlass.vue'
import type { ApiResponse, ClassRoom } from '@/types'

const classes = ref<ClassRoom[]>([])
const loading = ref(true)
const loadError = ref('')
const deleteStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const displayCodeStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const copyCodeStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const seriesStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const batchStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const createStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const importStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const assignTeacherStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

const gradeOptions = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']

// 默认全部折叠，创建后只展开对应年级
const expandedGrades = ref<Set<string>>(new Set())

// 弹窗状态
const showBatchClassModal = ref(false)
const showSingleClassModal = ref(false)
const showImportModal = ref(false)

// 批量创建班级
const batchGrade = ref('一年级')
const batchCount = ref(3)
const batchYear = ref(new Date().getFullYear())

// 单个创建班级
const newClassName = ref('')
const newClassGrade = ref('一年级')
const newClassYear = ref(new Date().getFullYear())

// 导入学生
const importClassName = ref('')
const importText = ref('')
const importResult = ref<{ success: number; failed: number; errors: string[] } | null>(null)
const batchError = ref('')
const singleClassError = ref('')
const importError = ref('')

const petSeriesOptions = [
  { value: 'all',      label: '不限制', emoji: '🌐' },
  { value: 'cosmic',   label: '原创宇宙', emoji: '🌌' },
  { value: 'pokemon',  label: '宝可梦', emoji: '⚡' },
  { value: 'cute',     label: '萌宠', emoji: '🐱' },
  { value: 'treasure', label: '国宝', emoji: '🐼' },
  { value: 'mythic',   label: '神兽', emoji: '🐉' },
]

const gradeEmojis: Record<string, string> = {
  '一年级': '🌱', '二年级': '🌿', '三年级': '🌳',
  '四年级': '🍂', '五年级': '⭐', '六年级': '🎓',
}

const groupedClasses = computed(() => {
  const groups: { grade: string; classes: ClassRoom[]; total: number }[] = []
  const classByGrade = new Map<string, ClassRoom[]>()
  for (const g of gradeOptions) classByGrade.set(g, [])
  for (const c of classes.value) {
    const grade = c.grade || '未分年级'
    if (!classByGrade.has(grade)) classByGrade.set(grade, [])
    classByGrade.get(grade)!.push(c)
  }
  for (const grade of gradeOptions) {
    const list = classByGrade.get(grade)
    if (list && list.length > 0) {
      groups.push({
        grade,
        classes: list,
        total: list.reduce((s, c) => s + c.student_count, 0),
      })
    }
  }
  for (const [grade, list] of classByGrade) {
    if (!gradeOptions.includes(grade) && list.length > 0) {
      groups.push({
        grade,
        classes: list,
        total: list.reduce((s, c) => s + c.student_count, 0),
      })
    }
  }
  return groups
})

onMounted(async () => {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
    loadError.value = ''
  } catch {
    classes.value = []
    loadError.value = '班级列表加载失败'
  }
  finally { loading.value = false }
})

function toggleGrade(grade: string) {
  if (expandedGrades.value.has(grade)) expandedGrades.value.delete(grade)
  else expandedGrades.value.add(grade)
  expandedGrades.value = new Set(expandedGrades.value)
}

function isExpanded(grade: string) {
  return expandedGrades.value.has(grade)
}

function fullyExpand() {
  expandedGrades.value = new Set(gradeOptions)
}

function collapseAll() {
  expandedGrades.value = new Set()
}

async function deleteClass(cls: ClassRoom) {
  const ok = await openConfirm({ title: '删除班级', message: `确定删除班级「${cls.name}」？\n班级下所有学生记录也会一并删除。`, danger: true, confirmText: '确认删除' })
  if (!ok) return
  deleteStatus.value[cls.id] = 'loading'
  try {
    await apiDelete(`/api/v1/admin/classes/${cls.id}`)
    classes.value = classes.value.filter(c => c.id !== cls.id)
    deleteStatus.value[cls.id] = 'success'
    setTimeout(() => { deleteStatus.value[cls.id] = 'idle' }, 1500)
  } catch {
    deleteStatus.value[cls.id] = 'error'
    setTimeout(() => { deleteStatus.value[cls.id] = 'idle' }, 3000)
  }
}

// ===== 班级码 =====

async function generateDisplayCode(cls: ClassRoom) {
  displayCodeStatus.value[cls.id] = 'loading'
  try {
    const res = await apiPost<{ data: { code: string } }>(`/api/v1/admin/classes/${cls.id}/display-code/refresh`)
    cls.display_code = res.data.code
    displayCodeStatus.value[cls.id] = 'success'
    setTimeout(() => { displayCodeStatus.value[cls.id] = 'idle' }, 2000)
  } catch {
    displayCodeStatus.value[cls.id] = 'error'
    setTimeout(() => { displayCodeStatus.value[cls.id] = 'idle' }, 3000)
  }
}

async function copyDisplayCode(cls: ClassRoom) {
  if (!cls.display_code) return
  try {
    await navigator.clipboard.writeText(cls.display_code)
    copyCodeStatus.value[cls.id] = 'success'
    setTimeout(() => { copyCodeStatus.value[cls.id] = 'idle' }, 1500)
  } catch {
    const ta = document.createElement('textarea')
    ta.value = cls.display_code
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    copyCodeStatus.value[cls.id] = 'success'
    setTimeout(() => { copyCodeStatus.value[cls.id] = 'idle' }, 1500)
  }
}

async function updatePetSeries(cls: ClassRoom, series: string) {
  seriesStatus.value[cls.id] = 'loading'
  try {
    await apiPut(`/api/v1/admin/classes/${cls.id}`, { pet_series: series })
    seriesStatus.value[cls.id] = 'success'
    setTimeout(() => { seriesStatus.value[cls.id] = 'idle' }, 2000)
  } catch {
    seriesStatus.value[cls.id] = 'error'
    setTimeout(() => { seriesStatus.value[cls.id] = 'idle' }, 3000)
  }
}

async function submitBatchClass() {
  if (batchCount.value < 1 || batchCount.value > 20) {
    batchError.value = '班级数量需在 1-20 之间'
    return
  }
  batchStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/classes/batch-create', {
      grade: batchGrade.value,
      count: batchCount.value,
      year: batchYear.value,
    })
    batchStatus.value = 'success'
    showBatchClassModal.value = false
    await reloadClasses(batchGrade.value)
    setTimeout(() => { batchStatus.value = 'idle' }, 2000)
  } catch {
    batchStatus.value = 'error'
    setTimeout(() => { batchStatus.value = 'idle' }, 3000)
  }
}

async function submitSingleClass() {
  const name = newClassName.value.trim()
  if (!name) {
    singleClassError.value = '请填写班级编号'
    return
  }
  const fullName = newClassGrade.value + '（' + name + '）班'
  if (classes.value.find(c => c.name === fullName)) {
    singleClassError.value = '班级「' + fullName + '」已存在'
    return
  }
  createStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/classes', {
      name: fullName,
      grade: newClassGrade.value,
      year: newClassYear.value,
    })
    createStatus.value = 'success'
    showSingleClassModal.value = false
    newClassName.value = ''
    await reloadClasses(newClassGrade.value)
    setTimeout(() => { createStatus.value = 'idle' }, 2000)
  } catch {
    createStatus.value = 'error'
    setTimeout(() => { createStatus.value = 'idle' }, 3000)
  }
}

async function submitImport() {
  const lines = importText.value.trim().split('\n').filter(l => l.trim())
  if (lines.length === 0) {
    importError.value = '请输入至少一位学生信息'
    return
  }

  const students = lines.map(line => {
    const [name, class_name, gender, student_no] = line.split(',').map(s => s?.trim() || '')
    return { name, class_name: class_name || importClassName.value, gender, student_no }
  })

  importStatus.value = 'loading'
  try {
    const res = await apiPost<ApiResponse<{ success: number; failed: number; errors: string[] }>>(
      '/api/v1/admin/students/import',
      { students },
    )
    const data = (res as unknown as { data: { success: number; failed: number; errors: string[] } }).data
    importResult.value = data || { success: students.length, failed: 0, errors: [] }
    importStatus.value = 'success'
    setTimeout(() => { importStatus.value = 'idle' }, 2000)
    await reloadClasses()
  } catch {
    importStatus.value = 'error'
    setTimeout(() => { importStatus.value = 'idle' }, 3000)
  }
}

async function reloadClasses(expandGrade?: string) {
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
    loadError.value = ''
    if (expandGrade) {
      expandGradeAfterCreate(expandGrade)
    }
  } catch {
    loadError.value = '班级列表加载失败'
  }
}

function expandGradeAfterCreate(grade: string) {
  expandedGrades.value.add(grade)
  expandedGrades.value = new Set(expandedGrades.value)
}

function openImportModal(className?: string) {
  importClassName.value = className || ''
  importText.value = ''
  importResult.value = null
  showImportModal.value = true
}

function downloadStudentTemplate() {
  const csv = '姓名,班级,性别,学号\n张小明,三年级（1）班,男,2026001\n李小红,三年级（1）班,女,2026002\n'
  const blob = new Blob(['﻿' + csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'students_template.csv'
  a.click()
  URL.revokeObjectURL(url)
}

// 分配班主任弹窗
const showAssignTeacherModal = ref(false)
const assigningClass = ref<ClassRoom | null>(null)
const teacherList = ref<Array<{ id: number; name: string }>>([])
const selectedTeacherId = ref<number | ''>('')

async function openAssignTeacherModal(cls: ClassRoom) {
  assigningClass.value = cls
  selectedTeacherId.value = cls.teacher_id || ''
  try {
    const res = await apiGet<ApiResponse<Array<{ id: number; name: string }>>>('/api/v1/admin/teachers')
    teacherList.value = res.data || []
  } catch { teacherList.value = [] }
  showAssignTeacherModal.value = true
}

async function submitAssignTeacher() {
  if (!assigningClass.value || selectedTeacherId.value === '') return
  assignTeacherStatus.value = 'loading'
  try {
    await apiPost(`/api/v1/admin/classes/${assigningClass.value.id}/assign-teacher`, {
      teacher_id: selectedTeacherId.value,
    })
    assignTeacherStatus.value = 'success'
    showAssignTeacherModal.value = false
    await reloadClasses()
    setTimeout(() => { assignTeacherStatus.value = 'idle' }, 2000)
  } catch {
    assignTeacherStatus.value = 'error'
    setTimeout(() => { assignTeacherStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div class="page-header">
      <div>
        <p class="page-subtitle">班级管理</p>
        <h2 class="page-title">班级列表</h2>
      </div>
      <div class="header-actions">
        <button class="btn btn-sm btn-ghost-card" @click="openImportModal()">📥 导入学生</button>
        <button class="btn btn-sm btn-primary" @click="showBatchClassModal = true">+ 批量添加班级</button>
        <button class="btn btn-sm btn-ghost-card" @click="showSingleClassModal = true">+ 添加班级</button>
      </div>
    </div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="reloadClasses()">重试</button>
    </div>

    <div v-else-if="classes.length === 0" class="card empty-state">
      <div class="empty-state__icon">🏫</div>
      <p>暂无班级</p>
    </div>

    <div v-else>
      <div class="toolbar">
        <button class="btn btn-sm btn-ghost-card" @click="fullyExpand">📂 全部展开</button>
        <button class="btn btn-sm btn-ghost-card" @click="collapseAll">📁 全部收起</button>
      </div>

      <div v-for="group in groupedClasses" :key="group.grade" class="grade-group">
        <div
          class="grade-header"
          :class="{ 'grade-header--expanded': isExpanded(group.grade) }"
          @click="toggleGrade(group.grade)"
        >
          <div class="grade-header__left">
            <span class="grade-header__emoji">{{ gradeEmojis[group.grade] || '📚' }}</span>
            <div>
              <div class="grade-header__name">{{ group.grade }}</div>
              <div class="grade-header__meta">{{ group.classes.length }} 个班级 · {{ group.total }} 名学生</div>
            </div>
          </div>
          <div class="grade-header__right">
            <button
              class="btn btn-sm btn-tool"
              @click.stop="openImportModal()"
            >📥 导入</button>
            <span class="grade-chevron" :class="{ 'grade-chevron--open': isExpanded(group.grade) }">▾</span>
          </div>
        </div>

        <div
          v-show="isExpanded(group.grade)"
          class="grade-body"
        >
          <div
            v-for="(c, i) in group.classes"
            :key="c.id"
            class="class-row"
            :class="{ 'class-row--last': i === group.classes.length - 1 }"
          >
            <div class="class-main">
              <div class="class-avatar">{{ c.name.replace(group.grade, '').replace(/（(\d+)）班/, '$1') }}</div>
              <div>
                <div class="class-name">{{ c.name }}</div>
                <div v-if="c.teacher_name" class="class-teacher">{{ c.teacher_name }}</div>
              </div>
            </div>
            <!-- 班级码 -->
            <div class="code-chip">
              <span class="code-chip__icon">🖥️</span>
              <code class="code-chip__value">{{ c.display_code || '--' }}</code>
              <div class="code-chip__actions">
                <button v-if="c.display_code" class="btn btn-sm code-btn" :class="copyCodeStatus[c.id] === 'loading' ? 'btn-state-loading' : copyCodeStatus[c.id] === 'success' ? 'btn-state-success' : copyCodeStatus[c.id] === 'error' ? 'btn-state-error' : ''" @click.stop="copyDisplayCode(c)" title="复制班级码">{{ copyCodeStatus[c.id] === 'success' ? '✅' : copyCodeStatus[c.id] === 'error' ? '❌' : '📋' }}</button>
                <button class="btn btn-sm code-btn" :class="displayCodeStatus[c.id] === 'loading' ? 'btn-state-loading' : displayCodeStatus[c.id] === 'success' ? 'btn-state-success' : displayCodeStatus[c.id] === 'error' ? 'btn-state-error' : ''" @click.stop="generateDisplayCode(c)" :disabled="displayCodeStatus[c.id] === 'loading'" title="刷新班级码旧码失效">{{ displayCodeStatus[c.id] === 'loading' ? '⏳' : displayCodeStatus[c.id] === 'success' ? '✅' : displayCodeStatus[c.id] === 'error' ? '❌' : '🔄' }}</button>
              </div>
            </div>
            <select
              class="pet-select"
              :value="(c as any).settings?.pet_series || 'all'"
              @change="updatePetSeries(c, ($event.target as HTMLSelectElement).value)"
              @click.stop
            >
              <option v-for="opt in petSeriesOptions" :key="opt.value" :value="opt.value">{{ opt.emoji }} {{ opt.label }}</option>
            </select>
            <span v-if="seriesStatus[c.id] === 'loading'" class="status-icon">⏳</span>
            <span v-else-if="seriesStatus[c.id] === 'success'" class="status-icon status-icon--success">✓</span>
            <span v-else-if="seriesStatus[c.id] === 'error'" class="status-icon status-icon--error">✗</span>
            <span class="student-count">{{ c.student_count }} 人</span>
            <div class="row-actions">
              <button class="btn btn-sm btn-mini" @click="openImportModal(c.name)">导入</button>
              <button class="btn btn-sm btn-ghost-card" @click="openAssignTeacherModal(c)">👨‍🏫 班主任</button>
              <button class="btn btn-sm" :class="deleteStatus[c.id] === 'loading' ? 'btn-state-loading' : deleteStatus[c.id] === 'success' ? 'btn-state-success' : deleteStatus[c.id] === 'error' ? 'btn-state-error' : 'btn-danger'" :disabled="deleteStatus[c.id] === 'loading'" @click="deleteClass(c)">{{ deleteStatus[c.id] === 'loading' ? '删除中...' : deleteStatus[c.id] === 'error' ? '删除失败' : '删除' }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 批量添加班级弹窗 -->
    <ModalGlass :visible="showBatchClassModal" @update:visible="showBatchClassModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">批量添加班级</h3>
        <button class="modal-close" @click="showBatchClassModal = false">×</button>
      </div>
      <div class="form-group">
        <label>年级</label>
        <select v-model="batchGrade" class="form-select">
          <option v-for="g in gradeOptions" :key="g" :value="g">{{ g }}</option>
        </select>
      </div>
      <div class="form-group">
        <label>班级数量</label>
        <input v-model.number="batchCount" type="number" min="1" max="20" class="form-input">
      </div>
      <div class="form-group">
        <label>学年</label>
        <input v-model.number="batchYear" type="number" class="form-input">
      </div>
      <div class="modal-actions">
        <div v-if="batchError" class="error-banner">{{ batchError }}</div>
        <button class="btn btn-sm btn-ghost-card" @click="showBatchClassModal = false">取消</button>
        <button class="btn btn-sm" :class="batchStatus === 'loading' ? 'btn-state-loading' : batchStatus === 'success' ? 'btn-state-success' : batchStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="batchStatus === 'loading'" @click="submitBatchClass">{{ batchStatus === 'loading' ? '创建中...' : batchStatus === 'success' ? '已创建 ✓' : batchStatus === 'error' ? '创建失败' : '创建' }}</button>
      </div>
    </ModalGlass>

    <!-- 单个添加班级弹窗 -->
    <ModalGlass :visible="showSingleClassModal" @update:visible="showSingleClassModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">添加班级</h3>
        <button class="modal-close" @click="showSingleClassModal = false">×</button>
      </div>
      <div class="form-group">
        <label>班级名称</label>
        <input v-model="newClassName" class="form-input" placeholder="如：3" @keydown.enter="submitSingleClass">
        <p class="modal-hint">输入班级编号，自动生成"{{ newClassGrade }}（{{ newClassName || '...' }}）班"</p>
      </div>
      <div class="form-group">
        <label>年级</label>
        <select v-model="newClassGrade" class="form-select">
          <option v-for="g in gradeOptions" :key="g" :value="g">{{ g }}</option>
        </select>
      </div>
      <div class="form-group">
        <label>学年</label>
        <input v-model.number="newClassYear" type="number" class="form-input">
      </div>
      <div class="modal-actions">
        <div v-if="singleClassError" class="error-banner">{{ singleClassError }}</div>
        <button class="btn btn-sm btn-ghost-card" @click="showSingleClassModal = false">取消</button>
        <button class="btn btn-sm" :class="createStatus === 'loading' ? 'btn-state-loading' : createStatus === 'success' ? 'btn-state-success' : createStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="createStatus === 'loading'" @click="submitSingleClass">{{ createStatus === 'loading' ? '创建中...' : createStatus === 'success' ? '已创建 ✓' : createStatus === 'error' ? '创建失败' : '创建' }}</button>
      </div>
    </ModalGlass>

    <!-- 分配班主任弹窗 -->
    <ModalGlass :visible="showAssignTeacherModal" @update:visible="showAssignTeacherModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">分配班主任</h3>
        <button class="modal-close" @click="showAssignTeacherModal = false">×</button>
      </div>
      <p class="assign-hint">为 <b>{{ assigningClass?.name }}</b> 分配班主任</p>
      <div class="form-group">
        <label>班主任</label>
        <select v-model="selectedTeacherId" class="form-select">
          <option value="">不分配（解除班主任）</option>
          <option v-for="t in teacherList" :key="t.id" :value="t.id">{{ t.name }}</option>
        </select>
      </div>
      <div class="modal-actions">
        <button class="btn btn-sm btn-ghost-card" @click="showAssignTeacherModal = false">取消</button>
        <button class="btn btn-sm" :class="assignTeacherStatus === 'loading' ? 'btn-state-loading' : assignTeacherStatus === 'success' ? 'btn-state-success' : assignTeacherStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="assignTeacherStatus === 'loading'" @click="submitAssignTeacher">
          {{ assignTeacherStatus === 'loading' ? '保存中...' : assignTeacherStatus === 'success' ? '已保存 ✓' : assignTeacherStatus === 'error' ? '保存失败' : '保存' }}
        </button>
      </div>
    </ModalGlass>

    <!-- 导入学生弹窗 -->
    <ModalGlass :visible="showImportModal" @update:visible="showImportModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">导入学生</h3>
        <div class="modal-header__right">
          <button class="btn btn-sm btn-ghost-card" @click="downloadStudentTemplate">📥 模板</button>
          <button class="modal-close" @click="showImportModal = false">×</button>
        </div>
      </div>

      <div v-if="importResult">
        <div class="import-result">
          <div class="import-result__count">{{ importResult.success }}</div>
          <div class="import-result__label">导入成功</div>
          <div v-if="importResult.failed > 0" class="import-result__failed">失败 {{ importResult.failed }} 人</div>
        </div>
        <div v-if="importResult.errors?.length" class="error-details">
          <p class="error-details__title">错误详情：</p>
          <pre class="error-details__list">{{ importResult.errors.join('\n') }}</pre>
        </div>
        <button class="btn btn-primary btn-block" @click="showImportModal = false">完成</button>
      </div>

      <div v-else>
        <div class="form-group">
          <label>默认班级（学生未填班级时使用）</label>
          <input v-model="importClassName" class="form-input" placeholder="如：三年级（1）班">
        </div>
        <div class="form-group">
          <label>学生数据</label>
          <p class="modal-hint">每行一位：姓名,班级,性别,学号（班级可留空使用上方默认）</p>
          <textarea
            v-model="importText"
            class="form-input textarea-monospace"
            placeholder="张小明,三年级（1）班,男,2026001&#10;李小红,三年级（1）班,女,2026002&#10;王刚,,男,2026003"
          ></textarea>
        </div>
        <div class="modal-actions">
          <div v-if="importError" class="error-banner">{{ importError }}</div>
          <button class="btn btn-sm btn-ghost-card" @click="showImportModal = false">取消</button>
          <button class="btn btn-sm" :class="importStatus === 'loading' ? 'btn-state-loading' : importStatus === 'success' ? 'btn-state-success' : importStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="importStatus === 'loading'" @click="submitImport">{{ importStatus === 'loading' ? '导入中...' : importStatus === 'success' ? '已导入 ✓' : importStatus === 'error' ? '导入失败' : '开始导入' }}</button>
        </div>
      </div>
    </ModalGlass>
  </div>
</template>

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); margin-bottom: 4px; }
.page-title { font-size: 24px; font-weight: 700; }
.header-actions { display: flex; gap: 8px; }

/* ===== 次要按钮（btn-ghost-card 已全局化） ===== */
.btn-mini {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  font-size: 12px;
}
.btn-tool {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  font-size: 11px;
}

/* ===== 工具栏 ===== */
.toolbar { display: flex; gap: 8px; margin-bottom: 16px; }

/* ===== 年级分组 ===== */
.grade-group { margin-bottom: 16px; }
.grade-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.15s;
}
.grade-header:hover { border-color: var(--tint-4); }
.grade-header--expanded {
  border-bottom-left-radius: 0;
  border-bottom-right-radius: 0;
  border-bottom-color: transparent;
}
.grade-header__left { display: flex; align-items: center; gap: 12px; }
.grade-header__emoji { font-size: 22px; }
.grade-header__name { font-size: 15px; font-weight: 700; }
.grade-header__meta { font-size: 12px; color: var(--color-text-secondary); margin-top: 2px; }
.grade-header__right { display: flex; align-items: center; gap: 8px; }
.grade-chevron {
  font-size: 18px;
  color: var(--color-text-secondary);
  transition: transform 0.2s;
  display: inline-block;
}
.grade-chevron--open { transform: rotate(180deg); }
.grade-body {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-top: none;
  border-radius: 0 0 12px 12px;
  overflow: hidden;
}

/* ===== 班级行 ===== */
.class-row {
  display: grid;
  grid-template-columns: 1fr auto auto auto auto auto;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid var(--color-border);
  transition: background 0.1s;
}
.class-row:hover { background: var(--tint-1); }
.class-row--last { border-bottom: none; }
.class-main { display: flex; align-items: center; gap: 10px; }
.class-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(79,70,229,0.12), rgba(99,102,241,0.06));
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 18px;
  color: var(--color-primary);
}
.class-name { font-weight: 600; font-size: 14px; }
.class-teacher { font-size: 12px; color: var(--color-text-secondary); }

/* ===== 班级码 ===== */
.code-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 16px;
  padding: 4px 10px 4px 8px;
  border: 1px solid rgba(79,70,229,0.12);
  border-radius: 8px;
  background: rgba(79,70,229,0.02);
  min-width: 140px;
}
.code-chip__icon { font-size: 11px; color: var(--color-text-secondary); white-space: nowrap; }
.code-chip__value {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-primary);
  letter-spacing: 0.08em;
  font-family: 'SF Mono', monospace;
}
.code-chip__actions { display: flex; gap: 2px; margin-left: auto; }

/* ===== 状态小按钮（复制/刷新班级码，状态色用全局 btn-state-*） ===== */
.code-btn {
  padding: 2px 6px;
  font-size: 11px;
  min-width: 0;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
}

/* ===== 宠物系列下拉 ===== */
.pet-select {
  margin-right: 16px;
  padding: 4px 8px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  color: var(--color-text);
  font-size: 12px;
  cursor: pointer;
}

/* ===== 状态图标 & 人数 ===== */
.status-icon { font-size: 11px; margin-right: 16px; }
.status-icon--success { color: var(--color-success-text); }
.status-icon--error { color: var(--color-danger-text); }
.student-count {
  margin-right: 16px;
  font-weight: 600;
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 60px;
  text-align: right;
}
.row-actions { display: flex; gap: 6px; }

/* ===== 弹窗（modal-header/title/close/actions/hint + error-banner 已全局化） ===== */
.modal-header__right { display: flex; align-items: center; gap: 12px; }
.assign-hint { font-size: 13px; color: var(--color-text-secondary); margin-bottom: 16px; }

/* ===== 导入结果 ===== */
.import-result { text-align: center; padding: 24px 0; }
.import-result__count { font-size: 36px; font-weight: 700; color: var(--color-accent); }
.import-result__label { color: var(--color-text-secondary); font-size: 13px; }
.import-result__failed { margin-top: 8px; color: var(--color-danger); font-size: 13px; }
.error-details { margin-bottom: 16px; }
.error-details__title { font-size: 12px; color: var(--color-danger); margin-bottom: 4px; }
.error-details__list {
  font-size: 12px;
  color: var(--color-danger);
  max-height: 120px;
  overflow-y: auto;
  background: rgba(239,68,68,0.05);
  padding: 8px;
  border-radius: 8px;
  white-space: pre-wrap;
}
.btn-block { width: 100%; }
.textarea-monospace { width: 100%; min-height: 160px; font-family: monospace; }
</style>
