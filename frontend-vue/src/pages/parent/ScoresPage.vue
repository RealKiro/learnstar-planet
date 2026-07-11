<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { formatTime } from '@/utils/constants'
import type { ApiResponse } from '@/types'

interface ChildInfo { id: number; name: string }

interface ScoreDetail {
  student_id: number
  student_name: string
  total_score: number
  today_score: number
  week_score: number
  class_rank: number
  class_name: string
}

interface ScoreHistory {
  id: number
  description: string
  balance_before: number
  balance_after: number
  change: number
  created_at: string
}

const children = ref<ChildInfo[]>([])
const selectedId = ref<number | null>(null)
const detail = ref<ScoreDetail | null>(null)
const history = ref<ScoreHistory[]>([])
const loading = ref(true)

async function fetchData(studentId: number) {
  loading.value = true
  try {
    const [dRes, hRes] = await Promise.all([
      apiGet<ApiResponse<ScoreDetail>>(`/api/v1/parent/scores/detail?student_id=${studentId}`),
      apiGet<ApiResponse<ScoreHistory[]>>(`/api/v1/parent/scores/history?student_id=${studentId}`),
    ])
    detail.value = dRes.data
    history.value = hRes.data || []
  } catch {
    detail.value = null
    history.value = []
  } finally {
    loading.value = false
  }
}

function onChildChange() {
  if (selectedId.value) fetchData(selectedId.value)
}

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<{ children: ChildInfo[] }>>('/api/v1/parent/home')
    children.value = res.data?.children || []
    if (children.value.length > 0) {
      selectedId.value = children.value[0].id
      await fetchData(selectedId.value)
      return
    }
  } catch { /* handled */ }
  loading.value = false
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">积分详情</h2>
      <select v-if="children.length > 1" v-model.number="selectedId" class="form-select" style="width:auto;" @change="onChildChange">
        <option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-else-if="detail">
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">⭐</span>
          <div class="stat-card__value">{{ detail.total_score.toLocaleString() }}</div>
          <div class="stat-card__label">总积分</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">📅</span>
          <div class="stat-card__value">{{ detail.today_score }}</div>
          <div class="stat-card__label">今日积分</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">📈</span>
          <div class="stat-card__value">{{ detail.week_score }}</div>
          <div class="stat-card__label">本周积分</div>
        </div>
        <div class="stat-card stat-card--info">
          <span class="stat-card__icon">🏆</span>
          <div class="stat-card__value">第{{ detail.class_rank }}名</div>
          <div class="stat-card__label">班级排名</div>
        </div>
      </div>

      <div class="card" style="margin-top:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">积分记录</h3>
        <div v-if="history.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无积分记录</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>描述</th><th style="width:100px;">变动</th><th style="width:100px;">余额</th><th style="width:160px;">时间</th></tr></thead>
            <tbody>
              <tr v-for="h in history" :key="h.id">
                <td>{{ h.description }}</td>
                <td :style="{ color: h.change >= 0 ? 'var(--color-accent)' : 'var(--color-danger)', fontWeight:600 }">
                  {{ h.change >= 0 ? '+' : '' }}{{ h.change }}
                </td>
                <td>{{ h.balance_after }}</td>
                <td style="color:var(--color-text-secondary);font-size:12px;">{{ formatTime(h.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <div v-else class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📊</div>
      <p>暂无积分数据</p>
    </div>
  </div>
</template>
