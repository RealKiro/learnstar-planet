<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ClassRoom } from '@/types'

const toast = useToastStore()
const classes = ref<ClassRoom[]>([])
const loading = ref(true)

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

const gradeOptions = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']

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

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
  } catch { classes.value = [] }
  finally { loading.value = false }
})

async function deleteClass(cls: ClassRoom) {
  if (!confirm(`确定删除班级「${cls.name}」？\n班级下所有学生记录也会一并删除。`)) return
  try {
    await apiDelete(`/api/v1/admin/classes/${cls.id}`)
    classes.value = classes.value.filter(c => c.id !== cls.id)
    toast.show('已删除班级：' + cls.name, 'success')
  } catch { /* handled */ }
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
    await reloadClasses()
  } catch { /* handled */ }
  finally { modalLoading.value = false }
}

async function submitSingleClass() {
  if (!newClassName.value.trim()) {
    toast.show('请填写班级名称', 'error')
    return
  }
  modalLoading.value = true
  try {
    await apiPost('/api/v1/admin/classes', {
      name: newClassName.value.trim(),
      grade: newClassGrade.value,
      year: newClassYear.value,
    })
    toast.show('已创建班级：' + newClassName.value, 'success')
    showSingleClassModal.value = false
    newClassName.value = ''
    await reloadClasses()
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

async function reloadClasses() {
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
  } catch { /* handled */ }
}

function openImportModal(className?: string) {
  importClassName.value = className || ''
  importText.value = ''
  importResult.value = null
  showImportModal.value = true
}

function downloadStudentTemplate() {
  const csv = '姓名,班级,性别,学号\n张小明,三年级（1）班,男,2026001\n李小红,三年级（1）班,女,2026002\n'
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'students_template.csv'
  a.click()
  URL.revokeObjectURL(url)
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

    <div v-else class="data-table">
      <table>
        <thead><tr><th>班级名称</th><th>年级</th><th>宠物系列</th><th>班主任</th><th>学生数</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="c in classes" :key="c.id">
            <td style="font-weight:600;">{{ c.name }}</td>
            <td><span style="display:inline-block;padding:4px 12px;border-radius:20px;font-size:12px;font-weight:600;background:rgba(79,70,229,0.08);color:var(--color-primary);">{{ c.grade || '未分' }}</span></td>
            <td>
              <select
                :value="(c as any).settings?.pet_series || 'all'"
                @change="updatePetSeries(c, ($event.target as HTMLSelectElement).value)"
                style="padding:4px 8px;border-radius:8px;border:1px solid var(--color-border);background:var(--color-bg);color:var(--color-text);font-size:13px;cursor:pointer;"
              >
                <option v-for="opt in petSeriesOptions" :key="opt.value" :value="opt.value">{{ opt.emoji }} {{ opt.label }}</option>
              </select>
            </td>
            <td>{{ c.teacher_name || '-' }}</td>
            <td style="font-weight:600;">{{ c.student_count }} 人</td>
            <td style="display:flex;gap:4px;">
              <button class="btn btn-sm" style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);" @click="openImportModal(c.name)">导入学生</button>
              <button class="btn btn-sm" style="background:#fee2e2;color:#dc2626;border:1px solid #fecaca;" @click="deleteClass(c)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
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
          <input v-model="newClassName" class="form-input" placeholder="如：三年级（3）班" @keydown.enter="submitSingleClass">
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

    <!-- 导入学生弹窗 -->
    <div v-if="showImportModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.self="showImportModal = false">
      <div class="card" style="width:90%;max-width:560px;max-height:85vh;overflow-y:auto;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">导入学生</h3>
          <div style="display:flex;align-items:center;gap:12px;">
            <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);font-size:12px;" @click="downloadStudentTemplate">📥 模板</button>
            <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showImportModal = false">×</button>
          </div>
        </div>

        <!-- 导入结果 -->
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

        <!-- 输入表单 -->
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
                                                                                     