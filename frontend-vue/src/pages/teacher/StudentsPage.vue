<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { getStageEmoji, getStageName } from '@/utils/constants'
import type { ApiResponse, Student } from '@/types'

const students = ref<Student[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students')
    students.value = res.data || []
  } catch {
    students.value = []
  } finally {
    loading.value = false
  }
})

function getStatusLabel(status: string): string {
  const map: Record<string, string> = { active: '活跃', graduated: '已毕业', disabled: '停用' }
  return map[status] || status
}
</script>

<template>
  <div>
    <div class="page-header" style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">学生列表</h2>
      <button class="btn btn-sm btn-primary">批量导入</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="students.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📭</div>
      <p>暂无学生数据，请联系管理员导入</p>
    </div>

    <div v-else class="data-table">
      <table>
        <thead><tr><th>姓名</th><th>学号</th><th>班级</th><th>总积分</th><th>状态</th></tr></thead>
        <tbody>
          <tr v-for="s in students" :key="s.id">
            <td style="font-weight:600;">{{ s.name }}</td>
            <td style="color:var(--color-text-secondary);">{{ s.student_no || '-' }}</td>
            <td>{{ s.class_name || s.class_grade || '-' }}</td>
            <td style="font-weight:700;color:var(--color-primary);">{{ s.total_score }}</td>
            <td :style="{ color: s.status === 'active' ? 'var(--color-accent)' : 'var(--color-text-secondary)' }">
              {{ getStatusLabel(s.status) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
