<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ClassRoom } from '@/types'

const toast = useToastStore()
const classes = ref<ClassRoom[]>([])
const loading = ref(true)
const displayCodeLoading = ref<Record<number, boolean>>({})

const gradeOptions = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']

// 默认全部折叠，创建后只展开对应年级
const expandedGrades = ref<Set<string>>(new Set())

// 弹窗状态
const showBatchClassModal = ref(false)
const showSingleClassModal = ref(false)
const showImportModal = ref(false)
const modalLoading = ref(false)

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

const petSeriesOptions = [
  { value: 'all',      label: '不限制', emoji: '🌐' },
  { value: 'cosmic',   label: '原创宇宙', emoji: '🌌' },
  { value: 'pokemon',  label: '宝可梦', emoji: '⚡' },
  { value: 'cute',     label: '萌宠', emoji: '🐱' },
  { value: 'treasure', label: '国宝', emoji: '🐼' },
  { value: 'mythic',   label: '神兽', emoji: '🐉' },
]

const petSeriesLabels: Record<string, string> = {
  all: '不限制', cosmic: '🌌 宇宙', pokemon: '⚡ 宝可梦', cute: '🐱 萌宠', treasure: '🐼 国宝', mythic: '🐉 神兽',
}

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
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
  } catch { classes.value = [] }
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

function expandGrade(grade: string) {
  if (!expandedGrades.value.has(grade)) {
    expandedGrades.value.add(grade)
    expandedGrades.value = new Set(expandedGrades.value)
  }
}

function fullyExpand() {
  expandedGrades.value = new Set(gradeOptions)
}

function collapseAll() {
  expandedGrades.value = new Set()
}

async function deleteClass(cls: ClassRoom) {
  if (!confirm(`确定删除班级「${cls.name}」？\n班级下所有学生记录也会一并删除。`)) return
  try {
    await apiDelete(`/api/v1/admin/classes/${cls.id}`)
    classes.value = classes.value.filter(c => c.id !== cls.id)
    toast.show('已删除班级：' + cls.name, 'success')
  } catch { /* handled */ }
}

// ===== 班级大屏码 =====

async function generateDisplayCode(cls: ClassRoom) {
  displayCodeLoading.value[cls.id] = true
  try {
    const res = await apiPost<{ data: { code: string } }>(`/api/v1/admin/classes/${cls.id}/display-code/refresh`)
    cls.display_code = res.data.code
    toast.show(`大屏码已生成：${res.data.code}`, 'success')
  } catch { /* handled */ }
  finally { displayCodeLoading.value[cls.id] = false }
}

async function copyDisplayCode(cls: ClassRoom) {
  if (!cls.display_code) return
  try {
    await navigator.clipboard.writeText(cls.display_code)
    toast.show('已复制：' + cls.display_code, 'success')
  } catch {
    const ta = document.createElement('textarea')
    ta.value = cls.display_code
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    toast.show('已复制', 'success')
  }
}

async function updatePetSeries(cls: ClassRoom, series: string) {
  try {
    await apiPut(`/api/v1/admin/classes/${cls.id}`, { pet_series: series })
    toast.show(`「${cls.name}」宠物系列已设为：${petSeriesLabels[series] || series}`, 'success')
  } catch { /* handled */ }
}

async function submitBatchClass() {
  if (batchCount.value < 1 || batchCount.value > 20) {
    toast.show('班级数量需在 1-20 之间', 'error')
    return
  }
  modalLoading.value = true
  try {
    await apiPost('/api/v1/admin/classes/batch-create', {
      grade: batchGrade.value,
      count: batchCount.value,
      year: batchYear.value,
    })
    toast.show(`已批量创建 ${batchCount.value} 个班级`, 'success')
    showBatchClassModal.value = false
    await reloadClasses(batchGrade.value)
  } catch { /* handled */ }
  finally { modalLoading.value = false }
}

async function submitSingleClass() {
  const name = newClassName.value.trim()
  if (!name) {
    toast.show('请填写班级编号', 'error')
    return
  }
  const fullName = newClassGrade.value + '（' + name + '）班'
  if (classes.value.find(c => c.name === fullName)) {
    toast.show('班级「' + fullName + '」已存在', 'error')
    return
  }
  modalLoading.value = true
  try {
    await apiPost('/api/v1/admin/classes', {
      name: fullName,
      grade: newClassGrade.value,
      year: newClassYear.value,
    })
    toast.show('已创建班级：' + fullName, 'success')
    showSingleClassModal.value = false
    newClassName.value = ''
    await reloadClasses(newClassGrade.value)
  } catch { /* handled */ }
  finally { modalLoading.value = false }
}

async function submitImport() {
  const lines = importText.value.trim().split('\n').filter(l => l.trim())
  if (lines.length === 0) {
    toast.show('请输入至少一位学生信息', 'error')
    return
  }

  const students = lines.map(line => {
    const [name, class_name, gender, student_no] = line.split(',').map(s => s?.trim() || '')
    return { name, class_name: class_name || importClassName.value, gender, student_no }
  })

  modalLoading.value = true
  try {
    const res = await apiPost<ApiResponse<{ success: number; failed: number; errors: string[] }>>(
      '/api/v1/admin/students/import',
      { students },
    )
    const data = (res as unknown as { data: { success: number; failed: number; errors: string[] } }).data
    importResult.value = data || { success: students.length, failed: 0, errors: [] }
    toast.show(`导入完成：成功 ${importResult.value.success} 人`, 'success')
    await reloadClasses()
  } catch { /* handled */ }
  finally { modalLoading.value = false }
}

async function reloadClasses(expandGrade?: string) {
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
    if (expandGrade) {
      expandGradeAfterCreate(expandGrade)
    }
  } catch { /* handled */ }
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
const assignTeacherLoading = ref(false)

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
  assignTeacherLoading.value = true
  try {
    await apiPost(`/api/v1/admin/classes/${assigningClass.value.id}/assign-teacher`, {
      teacher_id: selectedTeacherId.value,
    })
    toast.show(`已为「${assigningClass.value.name}」分配班主任`, 'success')
    showAssignTeacherModal.value = false
    await reloadClasses()
  } catch { /* handled */ }
  finally { assignTeacherLoading.value = false }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">班级管理</p>
        <h2 style="font-size:24px;font-weight:700;">班级列表</h2>
      </div>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="openImportModal()">📥 导入学生</button>
        <button class="btn btn-sm btn-primary" @click="showBatchClassModal = true">+ 批量添加班级</button>
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showSingleClassModal = true">+ 添加班级</button>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="classes.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">🏫</div>
      <p>暂无班级</p>
    </div>

    <div v-else>
      <div style="display:flex;gap:8px;margin-bottom:16px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="fullyExpand">📂 全部展开</button>
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="collapseAll">📁 全部收起</button>
      </div>

      <div v-for="group in groupedClasses" :key="group.grade" style="margin-bottom:16px;">
        <div
          @click="toggleGrade(group.grade)"
          style="display:flex;align-items:center;justify-content:space-between;padding:16px 20px;background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:12px;cursor:pointer;transition:all 0.15s;"
          :style="isExpanded(group.grade) ? { borderBottomLeftRadius: '0', borderBottomRightRadius: '0', borderBottomColor: 'transparent' } : {}"
        >
          <div style="display:flex;align-items:center;gap:12px;">
            <span style="font-size:22px;">{{ gradeEmojis[group.grade] || '📚' }}</span>
            <div>
              <div style="font-size:15px;font-weight:700;">{{ group.grade }}</div>
              <div style="font-size:12px;color:var(--color-text-secondary);margin-top:2px;">
                {{ group.classes.length }} 个班级 · {{ group.total }} 名学生
              </div>
            </div>
          </div>
          <div style="display:flex;align-items:center;gap:8px;">
            <button
              class="btn btn-sm"
              style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);font-size:11px;"
              @click.stop="openImportModal()"
            >📥 导入</button>
            <span style="font-size:18px;color:var(--color-text-secondary);transition:transform 0.2s;display:inline-block;" :style="isExpanded(group.grade) ? { transform: 'rotate(180deg)' } : {}">▾</span>
          </div>
        </div>

        <div
          v-show="isExpanded(group.grade)"
          style="background:var(--color-bg-card);border:1px solid var(--color-border);border-top:none;border-radius:0 0 12px 12px;overflow:hidden;"
        >
          <div
            v-for="(c, i) in group.classes"
            :key="c.id"
            :style="{
              display: 'grid',
              gridTemplateColumns: '1fr auto auto auto auto auto',
              alignItems: 'center',
              padding: '12px 20px',
              borderBottom: i < group.classes.length - 1 ? '1px solid var(--color-border)' : 'none',
              transition: 'background 0.1s',
            }"
          >
            <div style="display:flex;align-items:center;gap:10px;">
              <div :style="{
                width: '36px', height: '36px', borderRadius: '10px',
                background: 'linear-gradient(135deg, rgba(79,70,229,0.12), rgba(99,102,241,0.06))',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                fontWeight: 700, fontSize: '18px', color: 'var(--color-primary)'
              }">
                {{ c.name.replace(group.grade, '').replace(/（(\d+)）班/, '$1') }}
              </div>
              <div>
                <div style="font-weight:600;font-size:14px;">{{ c.name }}</div>
                <div v-if="c.teacher_name" style="font-size:12px;color:var(--color-text-secondary);">{{ c.teacher_name }}</div>
              </div>
            </div>
            <!-- 大屏码 -->
            <div style="display:flex;align-items:center;gap:6px;margin-right:16px;min-width:120px;">
              <template v-if="c.display_code">
                <code style="font-size:12px;font-weight:600;color:var(--color-primary);background:rgba(79,70,229,0.06);padding:3px 8px;border-radius:6px;letter-spacing:0.05em;">{{ c.display_code }}</code>
                <button class="btn btn-sm" style="padding:2px 6px;font-size:10px;border:1px solid var(--color-border);background:transparent;color:var(--color-text-secondary);min-width:0;" @click.stop="copyDisplayCode(c)">📋</button>
                <button class="btn btn-sm" style="padding:2px 6px;font-size:10px;border:1px solid var(--color-border);background:transparent;color:var(--color-text-secondary);min-width:0;" @click.stop="generateDisplayCode(c)" :disabled="displayCodeLoading[c.id]">🔄</button>
              </template>
              <button v-else class="btn btn-sm" style="font-size:11px;padding:2px 10px;" @click.stop="generateDisplayCode(c)" :disabled="displayCodeLoading[c.id]">
                {{ displayCodeLoading[c.id] ? '...' : '🖥️ 生成大屏码' }}
              </button>
            </div>
            <select
              :value="(c as any).settings?.pet_series || 'all'"
              @change="updatePetSeries(c, ($event.target as HTMLSelectElement).value)"
              @click.stop
              style="margin-right:16px;padding:4px 8px;border-radius:8px;border:1px solid var(--color-border);background:var(--color-bg);color:var(--color-text);font-size:12px;cursor:pointer;"
            >
              <option v-for="opt in petSeriesOptions" :key="opt.value" :value="opt.value">{{ opt.emoji }} {{ opt.label }}</option>
            </select>
            <span style="margin-right:16px;font-weight:600;font-size:13px;color:var(--color-text-secondary);min-width:60px;text-align:right;">{{ c.student_count }} 人</span>
            <div style="display:flex;gap:6px;">
              <button class="btn btn-sm" style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);font-size:12px;" @click="openImportModal(c.name)">导入</button>
              <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);font-size:12px;" @click="openAssignTeacherModal(c)">👨‍🏫 班主任</button>
              <button class="btn btn-sm" style="background:#fee2e2;color:#dc2626;border:1px solid #fecaca;font-size:12px;" @click="deleteClass(c)">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 批量添加班级弹窗 -->
    <div v-if="showBatchClassModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.stop>
      <div class="card" style="width:90%;max-width:420px;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">批量添加班级</h3>
          <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showBatchClassModal = false">×</button>
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
        <div style="display:flex;gap:8px;justify-content:flex-end;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showBatchClassModal = false">取消</button>
          <button class="btn btn-sm btn-primary" :disabled="modalLoading" @click="submitBatchClass">{{ modalLoading ? '创建中...' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- 单个添加班级弹窗 -->
    <div v-if="showSingleClassModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.stop>
      <div class="card" style="width:90%;max-width:420px;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">添加班级</h3>
          <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showSingleClassModal = false">×</button>
        </div>
        <div class="form-group">
          <label>班级名称</label>
          <input v-model="newClassName" class="form-input" placeholder="如：3" @keydown.enter="submitSingleClass">
          <p style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">输入班级编号，自动生成"{{ newClassGrade }}（{{ newClassName || '...' }}）班"</p>
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
        <div style="display:flex;gap:8px;justify-content:flex-end;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showSingleClassModal = false">取消</button>
          <button class="btn btn-sm btn-primary" :disabled="modalLoading" @click="submitSingleClass">{{ modalLoading ? '创建中...' : '创建' }}</button>
        </div>
      </div>
    </div>



    <!-- 分配班主任弹窗 -->
    <div v-if="showAssignTeacherModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.stop>
      <div class="card" style="width:90%;max-width:420px;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">分配班主任</h3>
          <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showAssignTeacherModal = false">×</button>
        </div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
          为 <b>{{ assigningClass?.name }}</b> 分配班主任
        </p>
        <div class="form-group">
          <label>班主任</label>
          <select v-model="selectedTeacherId" class="form-select">
            <option value="">不分配（解除班主任）</option>
            <option v-for="t in teacherList" :key="t.id" :value="t.id">{{ t.name }}</option>
          </select>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showAssignTeacherModal = false">取消</button>
          <button class="btn btn-sm btn-primary" :disabled="assignTeacherLoading" @click="submitAssignTeacher">
            {{ assignTeacherLoading ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 导入学生弹窗 -->
    <div v-if="showImportModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.stop>
      <div class="card" style="width:90%;max-width:560px;max-height:85vh;overflow-y:auto;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">导入学生</h3>
          <div style="display:flex;align-items:center;gap:12px;">
            <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);font-size:12px;" @click="downloadStudentTemplate">📥 模板</button>
            <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showImportModal = false">×</button>
          </div>
        </div>

        <div v-if="importResult">
          <div style="text-align:center;padding:24px 0;">
            <div style="font-size:36px;font-weight:700;color:var(--color-accent);">{{ importResult.success }}</div>
            <div style="color:var(--color-text-secondary);font-size:13px;">导入成功</div>
            <div v-if="importResult.failed > 0" style="margin-top:8px;color:var(--color-danger);font-size:13px;">失败 {{ importResult.failed }} 人</div>
          </div>
          <div v-if="importResult.errors?.length" style="margin-bottom:16px;">
            <p style="font-size:12px;color:var(--color-danger);margin-bottom:4px;">错误详情：</p>
            <pre style="font-size:12px;color:var(--color-danger);max-height:120px;overflow-y:auto;background:rgba(239,68,68,0.05);padding:8px;border-radius:8px;">{{ importResult.errors.join('\n') }}</pre>
          </div>
          <button class="btn btn-primary" style="width:100%;" @click="showImportModal = false">完成</button>
        </div>

        <div v-else>
          <div class="form-group">
            <label>默认班级（学生未填班级时使用）</label>
            <input v-model="importClassName" class="form-input" placeholder="如：三年级（1）班">
          </div>
          <div class="form-group">
            <label>学生数据</label>
            <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:8px;">每行一位：姓名,班级,性别,学号（班级可留空使用上方默认）</p>
            <textarea
              v-model="importText"
              class="form-input"
              style="width:100%;min-height:160px;font-family:monospace;"
              placeholder="张小明,三年级（1）班,男,2026001&#10;李小红,三年级（1）班,女,2026002&#10;王刚,,男,2026003"
            ></textarea>
          </div>
          <div style="display:flex;gap:8px;justify-content:flex-end;">
            <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showImportModal = false">取消</button>
            <button class="btn btn-sm btn-primary" :disabled="modalLoading" @click="submitImport">{{ modalLoading ? '导入中...' : '开始导入' }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
