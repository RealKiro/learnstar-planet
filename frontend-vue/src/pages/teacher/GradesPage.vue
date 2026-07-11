<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse } from '@/types'

interface GradeEntry {
  student_id: number
  student_name: string
  subjects: Record<string, number>
  total: number
  average: number
  rank: number
}

const grades = ref<GradeEntry[]>([])
const loading = ref(true)
const selectedExam = ref('期中考试')

const examOptions = ['期中考试', '月考一', '月考二', '期末考试']
const subjects = ['语文', '数学', '英语']

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<GradeEntry[]>>('/api/v1/teacher/grades')
    grades.value = res.data || []
  } catch {
    grades.value = demoGrades()
  }
  finally { loading.value = false }
})

function demoGrades(): GradeEntry[] {
  const names = [
    { name: '小华', subjects: { '语文': 92, '数学': 98, '英语': 95 } },
    { name: '小丽', subjects: { '语文': 95, '数学': 94, '英语': 92 } },
    { name: '小明', subjects: { '语文': 88, '数学': 96, '英语': 90 } },
    { name: '小红', subjects: { '语文': 90, '数学': 85, '英语': 88 } },
    { name: '小刚', subjects: { '语文': 78, '数学': 82, '英语': 80 } },
    { name: '小强', subjects: { '语文': 55, '数学': 62, '英语': 58 } },
  ]
  return names.map((s, i) => ({
    student_id: i + 1,
    student_name: s.name,
    subjects: s.subjects,
    total: Object.values(s.subjects).reduce((a, b) => a + b, 0),
    average: Math.round(Object.values(s.subjects).reduce((a, b) => a + b, 0) / 3 * 10) / 10,
    rank: i + 1,
  }))
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">📊 成绩管理</h2>
      <div style="display:flex;gap:8px;">
        <select v-model="selectedExam" class="form-select" style="width:auto;">
          <option v-for="e in examOptions" :key="e" :value="e">{{ e }}</option>
        </select>
        <button class="btn btn-sm btn-primary">录入成绩</button>
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);">📤 导出</button>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="grades.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📊</div>
      <p>暂无成绩数据</p>
    </div>

    <template v-else>
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">📈</span>
          <div class="stat-card__value">{{ Math.round(grades.reduce((a, g) => a + g.average, 0) / grades.length * 10) / 10 }}</div>
          <div class="stat-card__label">班级平均分</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">🏆</span>
          <div class="stat-card__value">{{ Math.max(...grades.map(g => g.total)) }}</div>
          <div class="stat-card__label">最高分</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">📉</span>
          <div class="stat-card__value">{{ Math.min(...grades.map(g => g.total)) }}</div>
          <div class="stat-card__label">最低分</div>
        </div>
        <div class="stat-card stat-card--info">
          <span class="stat-card__icon">✅</span>
          <div class="stat-card__value">{{ Math.round(grades.filter(g => g.total >= 180).length / grades.length * 100) }}%</div>
          <div class="stat-card__label">及格率</div>
        </div>
      </div>

      <div class="data-table">
        <div class="data-table__header"><h3 style="font-size:16px;font-weight:600;">{{ selectedExam }} 成绩</h3></div>
        <table>
          <thead><tr><th>姓名</th><th v-for="s in subjects" :key="s">{{ s }}</th><th>总分</th><th>平均分</th><th>班名次</th></tr></thead>
          <tbody>
            <tr v-for="g in grades" :key="g.student_id">
              <td style="font-weight:600;">{{ g.student_name }}</td>
              <td v-for="s in subjects" :key="s" :style="{ color: (g.subjects[s] || 0) >= 90 ? 'var(--color-accent)' : (g.subjects[s] || 0) < 60 ? 'var(--color-danger)' : '' }">
                {{ g.subjects[s] || '-' }}
              </td>
              <td style="font-weight:700;">{{ g.total }}</td>
              <td style="color:var(--color-text-secondary);">{{ g.average }}</td>
              <td :style="{ color: g.rank <= 3 ? ['#EAB308','#94A3B8','#CD7F32'][g.rank-1] : '', fontWeight: 700 }">{{ g.rank }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>
