<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { formatTime } from '@/utils/constants'
import type { ApiResponse } from '@/types'

interface ChildInfo { id: number; name: string }

interface TimelineItem {
  time: string
  title: string
  description: string
  icon: string
}

interface GrowthLog {
  type: 'score' | 'pet'
  description: string
  change: string
  balance: string
  created_at: string
}

const children = ref<ChildInfo[]>([])
const selectedId = ref<number | null>(null)
const timeline = ref<TimelineItem[]>([])
const logs = ref<GrowthLog[]>([])
const loading = ref(true)

const typeLabels: Record<string, string> = { score: '积分', pet: '宠物' }

async function fetchData(studentId: number) {
  loading.value = true
  try {
    const [tRes, lRes] = await Promise.all([
      apiGet<ApiResponse<TimelineItem[]>>(`/api/v1/parent/growth/timeline?student_id=${studentId}`),
      apiGet<ApiResponse<GrowthLog[]>>(`/api/v1/parent/growth/log?student_id=${studentId}`),
    ])
    timeline.value = tRes.data || []
    logs.value = lRes.data || []
  } catch {
    timeline.value = []
    logs.value = []
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
      <h2 style="font-size:24px;font-weight:700;">成长轨迹</h2>
      <select v-if="children.length > 1" v-model.number="selectedId" class="form-select" style="width:auto;" @change="onChildChange">
        <option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-else>
      <!-- Timeline -->
      <div class="card" style="margin-bottom:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:20px;">成长时间线</h3>
        <div v-if="timeline.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无时间线数据</div>
        <div v-else>
          <div v-for="(t, i) in timeline" :key="i" style="position:relative;padding-left:32px;padding-bottom:20px;">
            <span style="position:absolute;left:0;top:0;width:24px;height:24px;border-radius:50%;background:var(--color-bg-card);border:2px solid var(--color-primary);display:flex;align-items:center;justify-content:center;font-size:12px;z-index:1;">
              {{ t.icon }}
            </span>
            <div v-if="i < timeline.length - 1" style="position:absolute;left:11px;top:24px;bottom:0;width:2px;background:var(--color-border);"></div>
            <div style="font-weight:600;font-size:14px;">{{ t.title }}</div>
            <div style="font-size:13px;color:var(--color-text-secondary);margin-top:4px;">{{ t.description }}</div>
            <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">{{ formatTime(t.time) }}</div>
          </div>
        </div>
      </div>

      <!-- Growth Log -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">成长日志</h3>
        <div v-if="logs.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无日志数据</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th style="width:80px;">类型</th><th>描述</th><th style="width:100px;">变动</th><th style="width:160px;">时间</th></tr></thead>
            <tbody>
              <tr v-for="(l, i) in logs" :key="i">
                <td>
                  <span :style="{
                    padding:'2px 8px', borderRadius:'4px', fontSize:'11px', fontWeight:600,
                    background: l.type === 'score' ? 'rgba(79,70,229,0.1)' : 'rgba(16,185,129,0.1)',
                    color: l.type === 'score' ? '#4F46E5' : '#10B981',
                  }">{{ typeLabels[l.type] || l.type }}</span>
                </td>
                <td>{{ l.description }}</td>
                <td style="font-weight:600;">{{ l.change }}</td>
                <td style="color:var(--color-text-secondary);font-size:12px;">{{ formatTime(l.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>
