<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse, Student } from '@/types'

const students = ref<Student[]>([])
const loading = ref(true)
const importError = ref('')
const importStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students?per_page=200')
    students.value = res.data || []
  } catch { students.value = [] } finally { loading.value = false }
})

function getStatusLabel(status: string): string {
  const map: Record<string, string> = { active: '活跃', graduated: '已毕业', disabled: '停用' }
  return map[status] || status
}

function importStudents() {
  importError.value = ''
  importStatus.value = 'idle'
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.xlsx,.xls,.csv'
  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file) return
    importStatus.value = 'loading'
    try {
      const text = await file.text()
      // 简单 CSV 解析：每行格式为 "姓名,学号,班级"
      const lines = text.split('\n').filter(l => l.trim())
      const importedStudents = lines.slice(1).map(l => {
        const [name, student_no, class_name] = l.split(',').map(s => s.trim())
        return { name, student_no, class_name }
      }).filter(s => s.name)
      if (importedStudents.length === 0) {
        importStatus.value = 'error'
        importError.value = '文件为空或格式不正确（CSV: 姓名,学号,班级）'
        setTimeout(() => { if (importStatus.value === 'error') importStatus.value = 'idle' }, 3000)
        return
      }
      await apiPost('/api/v1/teacher/students', { students: importedStudents }, { skipToast: true })
      const res = await apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students?per_page=200')
      students.value = res.data || []
      importStatus.value = 'success'
      setTimeout(() => { if (importStatus.value === 'success') importStatus.value = 'idle' }, 1500)
    } catch (e: any) {
      importStatus.value = 'error'
      importError.value = e?.response?.data?.message || '导入失败，请检查文件格式（CSV: 姓名,学号,班级）'
      setTimeout(() => { if (importStatus.value === 'error') importStatus.value = 'idle' }, 3000)
    }
  }
  input.click()
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">🎒 学生列表</h2>
      <div style="display:flex;gap:8px;align-items:center;">
        <button
          class="btn btn-sm"
          :disabled="importStatus === 'loading'"
          :style="{
            background: importStatus === 'loading' ? '#f59e0b' : importStatus === 'success' ? '#10b981' : importStatus === 'error' ? '#ef4444' : 'var(--color-primary)',
            color: '#fff', border: '1px solid transparent'
          }"
          @click="importStudents"
        >
          {{ importStatus === 'loading' ? '导入中...' : importStatus === 'success' ? '已导入 ✓' : importStatus === 'error' ? '导入失败' : '📥 批量导入' }}
        </button>
      </div>
    </div>

    <div v-if="importError" style="margin:-12px 0 16px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color:#fca5a5;font-size:12px;">{{ importError }}</div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="students.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📭</div>
      <p>暂无学生数据</p>
    </div>

    <div v-else class="data-table">
      <table>
        <thead><tr><th>姓名</th><th>学号</th><th>积分</th><th>状态</th></tr></thead>
        <tbody>
          <tr v-for="s in students" :key="s.id">
            <td style="font-weight:600;">{{ s.name }}</td>
            <td style="color:var(--color-text-secondary);">{{ s.student_no || '-' }}</td>
            <td style="font-weight:700;color:var(--color-primary);">{{ s.total_score || 0 }}</td>
            <td :style="{ color: s.status === 'active' ? 'var(--color-accent)' : 'var(--color-text-secondary)' }">
              {{ getStatusLabel(s.status) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
