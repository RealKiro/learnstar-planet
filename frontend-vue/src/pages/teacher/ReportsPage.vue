<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { getStageEmoji, getStageName } from '@/utils/constants'
import type { ApiResponse } from '@/types'
import axios from 'axios'

const toast = useToastStore()

interface ScoreTrend { labels: string[]; datasets: { label: string; data: number[] }[] }
interface PetDist { stage_name: string; count: number; percentage: number }
interface StudentProgress { student_id: number; student_name: string; scores: number[]; trend: 'up'|'down'|'stable'; change: number }

interface ClassOption { class_id: number; class_name: string }

const loading = ref(true)
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

async function exportFile(type: string) {
  if (!selectedClassId.value) { toast.show('请先选择班级', 'error', { position: 'top-right' }); return }
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

onMounted(async () => {
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
  } catch { /* handled */ }
  finally { loading.value = false }
})

function stageLevel(name: string): number {
  const idx = ['星尘','月芽','灵苗','青藤','慧树','蝶灵','鹰慧','狮睿','灵角','星耀','银河'].indexOf(name)
  return idx >= 0 ? idx : 0
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">数据报表</h2>
      <div style="display:flex;gap:12px;align-items:center;">
        <select v-model.number="selectedClassId" class="form-select" style="width:auto;padding:6px 12px;">
          <option v-for="c in myClasses" :key="c.class_id" :value="c.class_id">{{ c.class_name }}</option>
        </select>
        <button class="btn btn-sm btn-primary" :style="{ background: exportStatus === 'loading' ? '#f59e0b' : exportStatus === 'success' ? '#10b981' : exportStatus === 'error' ? '#ef4444' : '' }" :disabled="exportStatus === 'loading'" @click="exportFile('scores')">
          <template v-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">导出成功 ✓</template>
          <template v-else-if="exportStatus === 'error'">导出失败 ✗</template>
          <template v-else>📥 导出积分</template>
        </button>
        <button class="btn btn-sm btn-primary" :style="{ background: exportStatus === 'loading' ? '#f59e0b' : exportStatus === 'success' ? '#10b981' : exportStatus === 'error' ? '#ef4444' : '' }" :disabled="exportStatus === 'loading'" @click="exportFile('pets')">
          <template v-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">导出成功 ✓</template>
          <template v-else-if="exportStatus === 'error'">导出失败 ✗</template>
          <template v-else>📥 导出宠物</template>
        </button>
        <button class="btn btn-sm btn-primary" :style="{ background: exportStatus === 'loading' ? '#f59e0b' : exportStatus === 'success' ? '#10b981' : exportStatus === 'error' ? '#ef4444' : '' }" :disabled="exportStatus === 'loading'" @click="exportFile('attendance')">
          <template v-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">导出成功 ✓</template>
          <template v-else-if="exportStatus === 'error'">导出失败 ✗</template>
          <template v-else>📥 导出考勤</template>
        </button>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-else>
      <!-- 1. 积分趋势柱状图 -->
      <div class="card" style="margin-bottom:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:24px;">近7日积分趋势</h3>
        <div v-if="!scoreTrend || scoreTrend.labels.length === 0" style="text-align:center;padding:32px;color:var(--color-text-secondary);">
          暂无趋势数据
        </div>
        <div v-else>
          <div style="display:flex;align-items:flex-end;gap:12px;height:200px;padding:0 8px;">
            <div v-for="(label, i) in scoreTrend.labels" :key="i" style="flex:1;display:flex;flex-direction:column;align-items:center;gap:8px;height:100%;justify-content:flex-end;">
              <div style="display:flex;gap:4px;align-items:flex-end;height:100%;">
                <div v-for="ds in scoreTrend.datasets" :key="ds.label"
                  :style="{
                    width: '24px',
                    height: `${(ds.data[i] || 0) / maxTrendValue * 160}px`,
                    background: ds.label.includes('扣') || ds.label.includes('减') ? 'var(--color-danger)' : 'var(--color-primary)',
                    borderRadius: '4px 4px 0 0',
                    minHeight: '2px',
                  }"
                  :title="`${ds.label}: ${ds.data[i] || 0}`">
                </div>
              </div>
              <span style="font-size:12px;color:var(--color-text-secondary);">{{ label }}</span>
            </div>
          </div>
          <div style="display:flex;gap:16px;margin-top:16px;justify-content:center;">
            <span v-for="ds in scoreTrend.datasets" :key="ds.label" style="display:flex;align-items:center;gap:6px;font-size:12px;color:var(--color-text-secondary);">
              <span :style="{ width:'10px',height:'10px',borderRadius:'2px',background: ds.label.includes('扣')||ds.label.includes('减') ? 'var(--color-danger)' : 'var(--color-primary)' }"></span>
              {{ ds.label }}
            </span>
          </div>
        </div>
      </div>

      <!-- 2. 宠物进化阶段分布 -->
      <div class="card" style="margin-bottom:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:24px;">宠物进化阶段分布</h3>
        <div v-if="petDist.length === 0" style="text-align:center;padding:32px;color:var(--color-text-secondary);">
          暂无宠物数据
        </div>
        <div v-else style="display:flex;flex-direction:column;gap:12px;">
          <div v-for="pet in petDist" :key="pet.stage_name" style="display:flex;align-items:center;gap:12px;">
            <span style="font-size:20px;width:28px;text-align:center;">{{ getStageEmoji(stageLevel(pet.stage_name)) }}</span>
            <span style="width:60px;font-size:13px;font-weight:500;">{{ pet.stage_name }}</span>
            <div style="flex:1;height:24px;background:var(--color-bg);border-radius:var(--radius-sm);overflow:hidden;">
              <div :style="{
                width: `${pet.percentage}%`,
                height: '100%',
                background: 'var(--gradient-primary)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'flex-end',
                paddingRight: '8px',
                color: 'white',
                fontSize: '12px',
                fontWeight: '600',
                minWidth: pet.count > 0 ? '32px' : '0',
              }">{{ pet.count }}</div>
            </div>
            <span style="width:48px;text-align:right;font-size:12px;color:var(--color-text-secondary);">{{ pet.percentage }}%</span>
          </div>
        </div>
      </div>

      <!-- 3. 学生进步情况 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">学生进步情况</h3>
        <div v-if="studentProgress.length === 0" style="text-align:center;padding:32px;color:var(--color-text-secondary);">
          暂无学生进度数据
        </div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>学生</th><th>近期积分</th><th>变化</th><th>趋势</th></tr></thead>
            <tbody>
              <tr v-for="s in studentProgress" :key="s.student_id">
                <td style="font-weight:600;">{{ s.student_name }}</td>
                <td>
                  <span style="font-family:monospace;font-size:13px;color:var(--color-text-secondary);">
                    {{ s.scores.join(' · ') }}
                  </span>
                </td>
                <td :style="{ color: trendColor(s.trend), fontWeight: '600' }">
                  {{ s.change > 0 ? '+' : '' }}{{ s.change }}
                </td>
                <td :style="{ color: trendColor(s.trend), fontSize: '18px', fontWeight: '700' }">
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
