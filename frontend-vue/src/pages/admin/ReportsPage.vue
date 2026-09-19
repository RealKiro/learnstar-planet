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
    overview.value = ovRes.data
    byGrade.value = gradeRes.data || []
    byClass.value = classRes.data || []
  } catch { overview.value = null; byGrade.value = []; byClass.value = [] }
  finally { loading.value = false }
})
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <p class="page-eyebrow">数据报表</p>
        <h2 class="page-title">学校报表</h2>
      </div>
      <span v-if="overview?.month_label" class="text-muted-13">{{ overview.month_label }} · 数据快照</span>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>

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

      <div v-if="overview.score_trend_percent !== undefined" class="card rep-bar-row">
        <span class="text-muted-13">积分环比：</span>
        <span :class="['trend-num', overview.score_trend_percent >= 0 ? 'num-up' : 'num-down']">
          {{ overview.score_trend_percent >= 0 ? '▲' : '▼' }} {{ Math.abs(overview.score_trend_percent) }}%
        </span>
      </div>

      <!-- 按年级汇总 -->
      <div class="card rep-mt-24">
        <h3 class="section-title">按年级汇总</h3>
        <div v-if="byGrade.length === 0" class="muted-center">暂无数据</div>
        <template v-else>
          <div class="rep-mb-20">
            <div v-for="g in byGrade" :key="g.grade" class="rep-line">
              <span class="rep-line-name">{{ g.grade }}</span>
              <div class="rep-bar-track">
                <div :style="{ width: (g.total_score / maxGradeScore * 100) + '%', height:'100%', background:'var(--ui-brand)', borderRadius:'6px', transition:'width 0.4s' }"></div>
              </div>
              <span class="rep-line-count">{{ g.total_score.toLocaleString() }} 分</span>
            </div>
          </div>
          <div class="data-table">
            <table>
              <thead><tr><th>年级</th><th>班级数</th><th>学生数</th><th>平均分</th><th>总积分</th></tr></thead>
              <tbody>
                <tr v-for="g in byGrade" :key="g.grade">
                  <td class="rep-fw-600">{{ g.grade }}</td>
                  <td>{{ g.class_count }}</td>
                  <td>{{ g.student_count }}</td>
                  <td>{{ g.avg_score.toFixed(1) }}</td>
                  <td class="rep-fw-accent">{{ g.total_score.toLocaleString() }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
      </div>

      <!-- 按班级明细 -->
      <div class="card rep-mt-24">
        <h3 class="section-title">按班级明细</h3>
        <div v-if="byClass.length === 0" class="muted-center">暂无数据</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>班级</th><th>年级</th><th>班主任</th><th>学生数</th><th>平均分</th><th>总积分</th></tr></thead>
            <tbody>
              <tr v-for="c in byClass" :key="c.class_name">
                <td class="rep-fw-600">{{ c.class_name }}</td>
                <td><span class="rep-pill">{{ c.grade || '-' }}</span></td>
                <td>{{ c.teacher_name || '-' }}</td>
                <td>{{ c.student_count }}</td>
                <td>{{ c.avg_score.toFixed(1) }}</td>
                <td class="rep-fw-accent">{{ c.total_score.toLocaleString() }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <div v-else class="card empty-state">
      <div class="empty-state__icon">📊</div>
      <p>暂无报表数据</p>
    </div>
  </div>
</template>

<style scoped>
/* ===== P1 内联样式收口（声明逐字保留以保渲染等价） ===== */
.rep-mt-24 { margin-top:24px; }
.rep-fw-600 { font-weight:600; }
.rep-fw-accent { font-weight:600;color:var(--ui-brand); }
.rep-bar-row { margin-top:16px;padding:12px 20px;display:flex;align-items:center;gap:8px; }
.rep-mb-20 { margin-bottom:20px; }
.rep-line { display:flex;align-items:center;gap:12px;margin-bottom:8px; }
.rep-line-name { width:60px;font-size:13px;font-weight:600;flex-shrink:0; }
.rep-bar-track { flex:1;height:24px;background:var(--color-bg);border-radius:6px;overflow:hidden; }
.rep-line-count { width:80px;font-size:12px;color:var(--color-text-secondary);text-align:right;flex-shrink:0; }
.rep-pill { display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;background:var(--ui-brand-soft);color:var(--color-primary); }
</style>
