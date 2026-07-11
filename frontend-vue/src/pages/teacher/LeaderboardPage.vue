<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse, LeaderboardEntry } from '@/types'

type LbType = 'total' | 'weekly' | 'pet'

const activeTab = ref<LbType>('total')
const entries = ref<LeaderboardEntry[]>([])
const loading = ref(false)

const tabs: Array<{ key: LbType; label: string }> = [
  { key: 'total', label: '总积分榜' },
  { key: 'weekly', label: '进步最快榜' },
  { key: 'pet', label: '宠物等级榜' },
]

const rankMedals = (i: number) => ['🥇', '🥈', '🥉'][i] || (i + 1).toString()

async function fetchData(type: LbType) {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<LeaderboardEntry[]>>(`/api/v1/teacher/leaderboard/${type}`)
    entries.value = res.data || []
  } catch {
    entries.value = []
  } finally {
    loading.value = false
  }
}

async function switchTab(type: LbType) {
  activeTab.value = type
  await fetchData(type)
}

onMounted(() => fetchData('total'))
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">排行榜</h2>
    </div>

    <div class="card">
      <div style="display:flex;gap:8px;margin-bottom:24px;">
        <button
          v-for="t in tabs" :key="t.key"
          :style="activeTab === t.key ? { background:'var(--color-primary)', color:'white', borderColor:'var(--color-primary)' } : {}"
          style="padding:8px 16px;border-radius:var(--radius-sm);font-size:13px;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);"
          @click="switchTab(t.key)"
        >{{ t.label }}</button>
      </div>

      <div v-if="loading" style="text-align:center;padding:24px;color:var(--color-text-secondary);">加载中...</div>
      <div v-else-if="entries.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无排行数据</div>
      <div v-else>
        <div v-for="(e, i) in entries" :key="e.student_id"
          style="display:flex;align-items:center;gap:16px;padding:12px 0;border-bottom:1px solid var(--color-border);">
          <span style="font-size:20px;width:32px;text-align:center;">{{ rankMedals(i) }}</span>
          <div style="width:36px;height:36px;border-radius:var(--radius-sm);background:var(--gradient-primary);color:white;display:flex;align-items:center;justify-content:center;font-size:14px;font-weight:600;">
            {{ e.student_name.charAt(0) }}
          </div>
          <span style="font-weight:500;flex:1;">{{ e.student_name }}</span>
          <span style="font-weight:700;font-size:16px;color:var(--color-primary);">{{ e.score }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
