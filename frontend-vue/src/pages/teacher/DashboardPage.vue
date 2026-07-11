<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { apiGet } from '@/utils/api'
import type { ApiResponse, Student } from '@/types'

const authStore = useAuthStore()

interface DashboardData {
  student_count: number
  weekly_score: number
  avg_pet_level: number
  pending_redemptions: number
  recent_scores: Array<{
    student_name: string
    points: number
    reason: string
    created_at: string
  }>
  top_students: Array<{
    student_name: string
    score: number
  }>
}

const data = ref<DashboardData | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<DashboardData>>('/api/teacher/dashboard')
    data.value = res.data
  } catch {
    // Fallback demo data
    data.value = {
      student_count: 42,
      weekly_score: 1260,
      avg_pet_level: 5,
      pending_redemptions: 8,
      recent_scores: [
        { student_name: '小明', points: 5, reason: '作业完成', created_at: '10分钟前' },
        { student_name: '小红', points: 10, reason: '作业优秀', created_at: '15分钟前' },
        { student_name: '小刚', points: -3, reason: '课堂违纪', created_at: '20分钟前' },
        { student_name: '小丽', points: 3, reason: '主动发言', created_at: '30分钟前' },
      ],
      top_students: [
        { student_name: '小明', score: 85 },
        { student_name: '小丽', score: 78 },
        { student_name: '小华', score: 72 },
        { student_name: '小红', score: 65 },
        { student_name: '小强', score: 58 },
      ],
    }
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">数据看板</h2>
      <button class="btn btn-sm btn-ghost">刷新数据</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-if="data">
      <!-- 统计卡片 -->
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">👨‍🎓</span>
          <div class="stat-card__value">{{ data.student_count }}</div>
          <div class="stat-card__label">班级学生</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">⭐</span>
          <div class="stat-card__value">{{ data.weekly_score.toLocaleString() }}</div>
          <div class="stat-card__label">本周积分发放</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">🌟</span>
          <div class="stat-card__value">Lv.{{ data.avg_pet_level }}</div>
          <div class="stat-card__label">全班平均宠物等级</div>
        </div>
        <div class="stat-card stat-card--info">
          <span class="stat-card__icon">🛍️</span>
          <div class="stat-card__value">{{ data.pending_redemptions }}</div>
          <div class="stat-card__label">待审批兑换</div>
        </div>
      </div>

      <!-- 最近积分动态 -->
      <div class="data-table" style="margin-bottom:24px;">
        <div class="data-table__header"><h3 style="font-size:16px;font-weight:600;">最近积分动态</h3></div>
        <table>
          <thead><tr><th>学生</th><th>积分变化</th><th>原因</th><th>时间</th></tr></thead>
          <tbody>
            <tr v-for="s in data.recent_scores" :key="s.student_name + s.created_at">
              <td>{{ s.student_name }}</td>
              <td :style="{ color: s.points >= 0 ? 'var(--color-accent)' : 'var(--color-danger)', fontWeight: 600 }">
                {{ s.points >= 0 ? '+' : '' }}{{ s.points }}
              </td>
              <td>{{ s.reason }}</td>
              <td>{{ s.created_at }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- TOP 5排行 -->
      <div class="card" style="margin-bottom:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">本周 TOP 5</h3>
        <div
          v-for="(s, i) in data.top_students"
          :key="s.student_name"
          style="display:flex;align-items:center;gap:16px;padding:12px 0;border-bottom:1px solid var(--color-border);"
        >
          <span style="font-size:20px;">{{ ['🥇','🥈','🥉'][i] || i + 1 }}</span>
          <div
            style="width:36px;height:36px;border-radius:var(--radius-sm);background:var(--gradient-primary);color:white;display:flex;align-items:center;justify-content:center;font-size:14px;font-weight:600;"
          >{{ s.student_name.charAt(0) }}</div>
          <span style="font-weight:500;flex:1;">{{ s.student_name }}</span>
          <span style="font-weight:700;font-size:16px;color:var(--color-primary);">{{ s.score }}分</span>
        </div>
      </div>
    </template>
  </div>
</template>
