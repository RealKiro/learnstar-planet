<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet } from '@/utils/api'
import { getStageEmoji } from '@/utils/constants'
import type { ApiResponse } from '@/types'
import axios from 'axios'

interface ScoreTrend { labels: string[]; datasets: { label: string; data: number[] }[] }
interface PetDist { stage_name: string; count: number; percentage: number }
interface StudentProgress { student_id: number; student_name: string; scores: number[]; trend: 'up'|'down'|'stable'; change: number }

interface ClassOption { class_id: number; class_name: string }

const loading = ref(true)
const loadError = ref('')
const exportStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const scoreTrend = ref<ScoreTrend | null>(null)
const petDist = ref<PetDist[]>([])
const studentProgress = ref<StudentProgress[]>([])
const myClasses = ref<ClassOption[]>([])
const selectedClassId = ref<number | null>(null)

const maxTrendValue = computed(() => {
  const all = scoreTrend.value?.datasets.flatMap(d => d.data) ?? []
  return Math.max(...all, 1)
})

const trendArrow = (t: string) => t === 'up' ? '↑' : t === 'down' ? '↓' : '→'
const trendColor = (t: string) => t === 'up' ? 'var(--color-accent)' : t === 'down' ? 'var(--color-danger)' : 'var(--color-text-secondary)'

const exportBtnClass = computed(() =>
  exportStatus.value === 'loading' ? 'btn-state-loading'
  : exportStatus.value === 'success' ? 'btn-state-success'
  : exportStatus.value === 'error' ? 'btn-state-error' : '',
)

async function exportFile(type: string) {
  exportStatus.value = 'loading'
  try {
    const token = localStorage.getItem('auth_token')
    const url = `/api/v1/teacher/reports/export/${type}?class_id=${selectedClassId.value}`
    const res = await axios.get(url, {
      responseType: 'blob',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
    // 下载文件
    const disposition = res.headers['content-disposition'] || ''
    const match = disposition.match(/filename\*?=(?:UTF-8'')?([^;]+)/)
    const filename = match ? decodeURIComponent(match[1]) : `${type}.xlsx`
    const blob = new Blob([res.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const urlObj = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = urlObj; a.download = filename
    document.body.appendChild(a); a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(urlObj)
    exportStatus.value = 'success'
    setTimeout(() => { exportStatus.value = 'idle' }, 1500)
  } catch {
    exportStatus.value = 'error'
    setTimeout(() => { exportStatus.value = 'idle' }, 3000)
  } finally {
    if (exportStatus.value === 'loading') exportStatus.value = 'idle'
  }
}

async function loadAll() {
  loading.value = true
  try {
    const [clsRes, trendRes, petRes, progressRes] = await Promise.all([
      apiGet<{ data: ClassOption[] }>('/api/v1/teacher/my-classes'),
      apiGet<ApiResponse<ScoreTrend>>('/api/v1/teacher/reports/score-trend'),
      apiGet<ApiResponse<PetDist[]>>('/api/v1/teacher/reports/pet-distribution'),
      apiGet<ApiResponse<StudentProgress[]>>('/api/v1/teacher/reports/student-progress'),
    ])
    myClasses.value = clsRes.data || []
    if (myClasses.value.length > 0) selectedClassId.value = myClasses.value[0].class_id
    scoreTrend.value = trendRes.data || null
    petDist.value = petRes.data || []
    studentProgress.value = progressRes.data || []
    loadError.value = ''
  } catch {
    loadError.value = '报表数据加载失败'
  }
  finally { loading.value = false }
}

onMounted(loadAll)

function stageLevel(name: string): number {
  const idx = ['星尘','月芽','灵苗','青藤','慧树','蝶灵','鹰慧','狮睿','灵角','星耀','银河'].indexOf(name)
  return idx >= 0 ? idx : 0
}
</script>

<template>
  <div>
    <div class="page-header">
      <h2 class="page-title">数据报表</h2>
      <div class="header-actions">
        <select v-model.number="selectedClassId" class="form-select class-select">
          <option v-for="c in myClasses" :key="c.class_id" :value="c.class_id">{{ c.class_name }}</option>
        </select>
        <button class="btn btn-sm btn-primary" :class="exportBtnClass" :disabled="exportStatus === 'loading'" @click="exportFile('scores')">
          <template v-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">导出成功 ✓</template>
          <template v-else-if="exportStatus === 'error'">导出失败 ✗</template>
          <template v-else>📥 导出积分</template>
        </button>
        <button class="btn btn-sm btn-primary" :class="exportBtnClass" :disabled="exportStatus === 'loading'" @click="exportFile('pets')">
          <template v-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">导出成功 ✓</template>
          <template v-else-if="exportStatus === 'error'">导出失败 ✗</template>
          <template v-else>📥 导出宠物</template>
        </button>
        <button class="btn btn-sm btn-primary" :class="exportBtnClass" :disabled="exportStatus === 'loading'" @click="exportFile('attendance')">
          <template v-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">导出成功 ✓</template>
          <template v-else-if="exportStatus === 'error'">导出失败 ✗</template>
          <template v-else>📥 导出考勤</template>
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="loadAll">重试</button>
    </div>

    <template v-else>
      <!-- 1. 积分趋势柱状图 -->
      <div class="card report-card">
        <h3 class="card-title">近7日积分趋势</h3>
        <div v-if="!scoreTrend || scoreTrend.labels.length === 0" class="empty-chart">
          暂无趋势数据
        </div>
        <div v-else>
          <div class="chart-wrap">
            <div v-for="(label, i) in scoreTrend.labels" :key="i" class="chart-col">
              <div class="chart-bars">
                <div v-for="ds in scoreTrend.datasets" :key="ds.label"
                  class="chart-bar"
                  :class="ds.label.includes('扣') || ds.label.includes('减') ? 'chart-bar--deduct' : ''"
                  :style="{ height: `${(ds.data[i] || 0) / maxTrendValue * 160}px` }"
                  :title="`${ds.label}: ${ds.data[i] || 0}`">
                </div>
              </div>
              <span class="chart-label">{{ label }}</span>
            </div>
          </div>
          <div class="chart-legend">
            <span v-for="ds in scoreTrend.datasets" :key="ds.label" class="legend-item">
              <span class="legend-swatch" :class="ds.label.includes('扣') || ds.label.includes('减') ? 'legend-swatch--deduct' : ''"></span>
              {{ ds.label }}
            </span>
          </div>
        </div>
      </div>

      <!-- 2. 宠物进化阶段分布 -->
      <div class="card report-card">
        <h3 class="card-title">宠物进化阶段分布</h3>
        <div v-if="petDist.length === 0" class="empty-chart">
          暂无宠物数据
        </div>
        <div v-else class="dist-list">
          <div v-for="pet in petDist" :key="pet.stage_name" class="dist-row">
            <span class="dist-emoji">{{ getStageEmoji(stageLevel(pet.stage_name)) }}</span>
            <span class="dist-name">{{ pet.stage_name }}</span>
            <div class="dist-track">
              <div class="dist-fill"
                :style="{ width: `${pet.percentage}%`, minWidth: pet.count > 0 ? '32px' : '0' }"
              >{{ pet.count }}</div>
            </div>
            <span class="dist-pct">{{ pet.percentage }}%</span>
          </div>
        </div>
      </div>

      <!-- 3. 学生进步情况 -->
      <div class="card">
        <h3 class="card-title card-title--sm">学生进步情况</h3>
        <div v-if="studentProgress.length === 0" class="empty-chart">
          暂无学生进度数据
        </div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>学生</th><th>近期积分</th><th>变化</th><th>趋势</th></tr></thead>
            <tbody>
              <tr v-for="s in studentProgress" :key="s.student_id">
                <td class="student-name">{{ s.student_name }}</td>
                <td>
                  <span class="mono-scores">
                    {{ s.scores.join(' · ') }}
                  </span>
                </td>
                <td class="trend-change" :style="{ color: trendColor(s.trend) }">
                  {{ s.change > 0 ? '+' : '' }}{{ s.change }}
                </td>
                <td class="trend-arrow" :style="{ color: trendColor(s.trend) }">
                  {{ trendArrow(s.trend) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-title { font-size: 24px; font-weight: 700; }
.header-actions { display: flex; gap: 12px; align-items: center; }
.class-select { width: auto; padding: 6px 12px; }

/* ===== 报表卡片 ===== */
.report-card { margin-bottom: 24px; }
.card-title { font-size: 16px; font-weight: 600; margin-bottom: 24px; }
.card-title--sm { margin-bottom: 16px; }
.empty-chart { text-align: center; padding: 32px; color: var(--color-text-secondary); }

/* ===== 积分趋势柱状图 ===== */
.chart-wrap { display: flex; align-items: flex-end; gap: 12px; height: 200px; padding: 0 8px; }
.chart-col { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 8px; height: 100%; justify-content: flex-end; }
.chart-bars { display: flex; gap: 4px; align-items: flex-end; height: 100%; }
.chart-bar { width: 24px; border-radius: 4px 4px 0 0; min-height: 2px; background: var(--color-primary); }
.chart-bar--deduct { background: var(--color-danger); }
.chart-label { font-size: 12px; color: var(--color-text-secondary); }
.chart-legend { display: flex; gap: 16px; margin-top: 16px; justify-content: center; }
.legend-item { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--color-text-secondary); }
.legend-swatch { width: 10px; height: 10px; border-radius: 2px; background: var(--color-primary); }
.legend-swatch--deduct { background: var(--color-danger); }

/* ===== 宠物分布 ===== */
.dist-list { display: flex; flex-direction: column; gap: 12px; }
.dist-row { display: flex; align-items: center; gap: 12px; }
.dist-emoji { font-size: 20px; width: 28px; text-align: center; }
.dist-name { width: 60px; font-size: 13px; font-weight: 500; }
.dist-track { flex: 1; height: 24px; background: var(--color-bg); border-radius: var(--radius-sm); overflow: hidden; }
.dist-fill {
  height: 100%;
  background: var(--gradient-primary);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 8px;
  color: #fff;
  font-size: 12px;
  font-weight: 600;
}
.dist-pct { width: 48px; text-align: right; font-size: 12px; color: var(--color-text-secondary); }

/* ===== 学生进步表 ===== */
.student-name { font-weight: 600; }
.mono-scores { font-family: monospace; font-size: 13px; color: var(--color-text-secondary); }
.trend-change { font-weight: 600; }
.trend-arrow { font-size: 18px; font-weight: 700; }
</style>
