<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { avatarGradient } from '@/utils/constants'
import type { ApiResponse } from '@/types'

interface ChildInfo { id: number; name: string }

interface RankEntry {
  rank: number
  student_id: number
  student_name: string
  score: number
  is_mine: boolean
}

const children = ref<ChildInfo[]>([])
const selectedId = ref<number | null>(null)
const rankings = ref<RankEntry[]>([])
const loading = ref(true)

function rankMedal(rank: number): string {
  return ['🥇', '🥈', '🥉'][rank - 1] || rank.toString()
}

async function fetchData(studentId: number) {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<RankEntry[]>>(`/api/v1/parent/ranking?student_id=${studentId}`)
    rankings.value = res.data || []
  } catch {
    rankings.value = []
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
      <h2 style="font-size:24px;font-weight:700;">班级排行</h2>
      <select v-if="children.length > 1" v-model.number="selectedId" class="form-select" style="width:auto;" @change="onChildChange">
        <option v-for="c in children" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
    </div>

    <div class="card">
      <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>
      <div v-else-if="rankings.length === 0" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
        <div style="font-size:48px;margin-bottom:8px;">🏆</div>
        <p>暂无排行数据</p>
      </div>
      <div v-else>
        <div v-for="r in rankings" :key="r.student_id"
          :style="{
            display:'flex', alignItems:'center', gap:'16px', padding:'12px 16px',
            borderRadius:'var(--radius-sm)', marginBottom:'8px',
            background: r.is_mine ? 'rgba(79,70,229,0.08)' : 'transparent',
            border: r.is_mine ? '1px solid rgba(79,70,229,0.3)' : '1px solid transparent',
          }">
          <span style="font-size:20px;width:32px;text-align:center;">{{ rankMedal(r.rank) }}</span>
          <div :style="{ width:'36px', height:'36px', borderRadius:'50%', background:avatarGradient(r.student_name), color:'white', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'14px', fontWeight:600 }">
            {{ r.student_name.charAt(0) }}
          </div>
          <span :style="{ flex:1, fontWeight: r.is_mine ? 600 : 500 }">
            {{ r.student_name }}
            <span v-if="r.is_mine" style="font-size:11px;color:var(--color-primary);margin-left:8px;">我的孩子</span>
          </span>
          <span style="font-weight:700;font-size:16px;color:var(--color-primary);">{{ r.score }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
