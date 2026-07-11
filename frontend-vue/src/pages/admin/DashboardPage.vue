<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, SchoolOverview } from '@/types'

const toast = useToastStore()
const overview = ref<SchoolOverview | null>(null)
const byGrade = ref<Array<{ grade: string; class_count: number; student_count: number }>>([])
const loading = ref(true)
const schoolName = ref('')

onMounted(async () => {
  try {
    const [overviewRes, gradeRes, schoolRes] = await Promise.all([
      apiGet<ApiResponse<SchoolOverview>>('/api/admin/reports/overview'),
      apiGet<ApiResponse<Array<{ grade: string; class_count: number; student_count: number }>>>('/api/admin/reports/by-grade'),
      apiGet<ApiResponse<{ name: string }>>('/api/admin/school'),
    ])
    overview.value = overviewRes.data
    byGrade.value = gradeRes.data || []
    schoolName.value = schoolRes.data?.name || ''
  } catch {
    // Demo数据
    overview.value = {
      class_count: 24,
      teacher_count: 48,
      student_count: 1200,
      monthly_score: 15600,
      parent_count: 850,
    }
    byGrade.value = [
      { grade: '一年级', class_count: 4, student_count: 200 },
      { grade: '二年级', class_count: 4, student_count: 195 },
      { grade: '三年级', class_count: 5, student_count: 220 },
      { grade: '四年级', class_count: 4, student_count: 210 },
      { grade: '五年级', class_count: 4, student_count: 205 },
      { grade: '六年级', class_count: 3, student_count: 170 },
    ]
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">欢迎回来</p>
        <h2 style="font-size:24px;font-weight:700;">全校看板</h2>
      </div>
      <span style="font-size:13px;color:var(--color-text-secondary);">{{ schoolName || '学校' }} · 今日数据概览</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-if="overview">
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">🏫</span>
          <div class="stat-card__value">{{ overview.class_count }}</div>
          <div class="stat-card__label">班级总数</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">👨‍🏫</span>
          <div class="stat-card__value">{{ overview.teacher_count }}</div>
          <div class="stat-card__label">教师账号</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">👨‍🎓</span>
          <div class="stat-card__value">{{ overview.student_count.toLocaleString() }}</div>
          <div class="stat-card__label">学生总数</div>
        </div>
        <div class="stat-card stat-card--info">
          <span class="stat-card__icon">⭐</span>
          <div class="stat-card__value">{{ overview.monthly_score.toLocaleString() }}</div>
          <div class="stat-card__label">本月积分发放</div>
        </div>
      </div>

      <!-- 按年级分布 -->
      <div class="card" style="margin-top:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">按年级分布</h3>
        <div class="data-table">
          <table>
            <thead><tr><th style="width:60px;">年级</th><th style="width:60px;">班级数</th><th style="width:60px;">学生数</th></tr></thead>
            <tbody>
              <tr v-for="g in byGrade" :key="g.grade">
                <td style="font-weight:600;">{{ g.grade }}</td>
                <td>{{ g.class_count }}</td>
                <td>{{ g.student_count }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>
