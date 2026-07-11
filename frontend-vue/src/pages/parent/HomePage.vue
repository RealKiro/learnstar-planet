<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { avatarGradient } from '@/utils/constants'
import type { ApiResponse } from '@/types'

interface ChildInfo {
  id: number
  name: string
  student_no: string
  class_name: string
  grade: string
  total_score: number
  class_rank: number
  pet_name: string
  pet_level: number
  pet_stage: string
  pet_emoji: string
  pet_mood: string
}

const children = ref<ChildInfo[]>([])
const loading = ref(true)

const totalScore = computed(() => children.value.reduce((s, c) => s + (c.total_score || 0), 0))
const bestRank = computed(() => {
  const ranks = children.value.map(c => c.class_rank).filter(r => r > 0)
  return ranks.length ? Math.min(...ranks) : '-'
})

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<{ children: ChildInfo[] }>>('/api/v1/parent/home')
    children.value = res.data?.children || []
  } catch {
    children.value = []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div style="margin-bottom:24px;">
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">家长中心</p>
      <h2 style="font-size:24px;font-weight:700;">孩子动态</h2>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="children.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">👶</div>
      <p>尚未绑定孩子，请联系学校管理员</p>
    </div>

    <template v-else>
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">👶</span>
          <div class="stat-card__value">{{ children.length }}</div>
          <div class="stat-card__label">孩子数量</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">⭐</span>
          <div class="stat-card__value">{{ totalScore.toLocaleString() }}</div>
          <div class="stat-card__label">总积分</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">🏆</span>
          <div class="stat-card__value">{{ bestRank }}</div>
          <div class="stat-card__label">最高排名</div>
        </div>
      </div>

      <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:16px;margin-top:24px;">
        <div v-for="c in children" :key="c.id" class="card" style="padding:20px;">
          <div style="display:flex;align-items:center;gap:12px;margin-bottom:16px;">
            <div :style="{ width:'48px', height:'48px', borderRadius:'50%', background:avatarGradient(c.name), color:'white', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'20px', fontWeight:600 }">
              {{ c.name.charAt(0) }}
            </div>
            <div>
              <div style="font-weight:600;font-size:16px;">{{ c.name }}</div>
              <div style="font-size:12px;color:var(--color-text-secondary);">{{ c.class_name }} · {{ c.grade }}</div>
            </div>
          </div>

          <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
            <div>
              <div style="font-size:12px;color:var(--color-text-secondary);">总积分</div>
              <div style="font-size:20px;font-weight:700;color:var(--color-primary);">{{ c.total_score }}</div>
            </div>
            <div style="text-align:right;">
              <div style="font-size:12px;color:var(--color-text-secondary);">班级排名</div>
              <div style="font-size:20px;font-weight:700;color:var(--color-accent);">第{{ c.class_rank }}名</div>
            </div>
          </div>

          <div style="display:flex;align-items:center;gap:12px;padding:12px;background:var(--color-bg);border-radius:var(--radius-sm);">
            <span style="font-size:36px;">{{ c.pet_emoji }}</span>
            <div>
              <div style="font-weight:500;font-size:14px;">{{ c.pet_name }}</div>
              <div style="font-size:12px;color:var(--color-text-secondary);">{{ c.pet_stage }} · Lv.{{ c.pet_level }}</div>
              <div style="font-size:12px;color:var(--color-text-secondary);">心情：{{ c.pet_mood }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
