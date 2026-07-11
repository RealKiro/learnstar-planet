<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse, SchoolOverview } from '@/types'

const overview = ref<SchoolOverview | null>(null)
const byGrade = ref<Array<{ grade: string; class_count: number; student_count: number; avg_score: number; total_score: number }>>([])
const byClass = ref<Array<{ class_name: string; grade: string; student_count: number; avg_score: number; total_score: number; teacher_name: string }>>([])
const loading = ref(true)

const maxGradeScore = computed(() => Math.max(1, ...byGrade.value.map(g => g.total_score)))

onMounted(async () => {
  try {
    const [ovRes, gradeRes, classRes] = await Promise.all([
      apiGet<ApiResponse<SchoolOverview>>('/api/v1/admin/reports/overview'),
      apiGet<ApiResponse<Array<{ grade: string; class_count: number; student_count: number; avg_score: number; total_score: number }>>>('/api/v1/admin/reports/by-grade'),
      apiGet<ApiResponse<Array<{ class_name: string; grade: string; student_count: number; avg_score: number; total_score: number; teacher_name: string }>>>('/api/v1/admin/reports/by-class'),
    ])
    overview.value = (ovRes as unknown as { data: SchoolOverview }).data
    byGrade.value = (gradeRes as unknown as { data: typeof byGrade.value }).data || []
    byClass.value = (classRes as unknown as { data: typeof byClass.value }).data || []
  } catch { overview.value = null; byGrade.value = []; byClass.value = [] }
  finally { loading.value = false }
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">数据报表</p>
        <h2 style="font-size:24px;font-weight:700;">学校报表</h2>
      </div>
      <span v-if="overview?.month_label" style="font-size:13px;color:var(--color-text-secondary);">{{ overview.month_label }} · 数据快照</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-else-if="overview">
      <!-- 顶部统计卡片 -->
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

      <div v-if="overview.score_trend_percent !== undefined" class="card" style="margin-top:16px;padding:12px 20px;display:flex;align-items:center;gap:8px;">
        <span style="font-size:13px;color:var(--color-text-secondary);">积分环比：</span>
        <span :style="{ fontSize:'13px', fontWeight:600, color: overview.score_trend_percent >= 0 ? '#10B981' : '#EF4444' }">
          {{ overview.score_trend_percent >= 0 ? '▲' : '▼' }} {{ Math.abs(overview.score_trend_percent) }}%
        </span>
      </div>

      <!-- 按年级汇总 + 条形图 -->
      <div class="card" style="margin-top:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">按年级汇总</h3>
        <div v-if="byGrade.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无数据</div>
        <template v-else>
          <!-- 条形图 -->
          <div style="margin-bottom:20px;">
            <div v-for="g in byGrade" :key="g.grade" style="display:flex;align-items:center;gap:12px;margin-bottom:8px;">
              <span style="width:60px;font-size:13px;font-weight:600;flex-shrink:0;">{{ g.grade }}</span>
              <div style="flex:1;height:24px;background:var(--color-bg);border-radius:6px;overflow:hidden;">
                <div :style="{ width: (g.total_score / maxGradeScore * 100) + '%', height:'100%', background:'linear-gradient(90deg,#6366F1,#818CF8)', borderRadius:'6px', transition:'width 0.4s' }"></div>
              </div>
              <span style="width:80px;font-size:12px;color:var(--color-text-secondary);text-align:right;flex-shrink:0;">{{ g.total_score.toLocaleString() }} 分</span>
            </div>
          </div>
          <!-- 表格 -->
          <div class="data-table">
            <table>
              <thead><tr><th>年级</th><th>班级数</th><th>学生数</th><th>平均分</th><th>总积分</th></tr></thead>
              <tbody>
                <tr v-for="g in byGrade" :key="g.grade">
                  <td style="font-weight:600;">{{ g.grade }}</td>
                  <td>{{ g.class_count }}</td>
                  <td>{{ g.student_count }}</td>
                  <td>{{ g.avg_score.toFixed(1) }}</td>
                  <td style="font-weight:600;color:var(--color-accent);">{{ g.total_score.toLocaleString() }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
      </div>

      <!-- 按班级明细 -->
      <div class="card" style="margin-top:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">按班级明细</h3>
        <div v-if="byClass.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无数据</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>班级</th><th>年级</th><th>班主任</th><th>学生数</th><th>平均分</th><th>总积分</th></tr></thead>
            <tbody>
              <tr v-for="c in byClass" :key="c.class_name">
                <td style="font-weight:600;">{{ c.class_name }}</td>
                <td><span style="display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;background:rgba(79,70,229,0.08);color:var(--color-primary);">{{ c.grade || '-' }}</span></td>
                <td>{{ c.teacher_name || '-' }}</td>
                <td>{{ c.student_count }}</td>
                <td>{{ c.avg_score.toFixed(1) }}</td>
                <td style="font-weight:600;color:var(--color-accent);">{{ c.total_score.toLocaleString() }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <div v-else class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📊</div>
      <p>暂无报表数据</p>
    </div>
  </div>
</template>
