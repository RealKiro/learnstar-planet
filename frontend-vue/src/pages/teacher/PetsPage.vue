<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { getStageEmoji, getStageName } from '@/utils/constants'
import type { ApiResponse, Pet } from '@/types'

const pets = ref<Pet[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Pet[]>>('/api/v1/teacher/pets/class-overview')
    pets.value = res.data || []
  } catch {
    // Demo fallback
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">宠物园</h2>
      <button class="btn btn-sm btn-ghost">全班喂养</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="pets.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">🌟</div>
      <p>暂无宠物数据</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:16px;">
      <div v-for="p in pets" :key="p.id" class="card"
        style="text-align:center;cursor:pointer;padding:16px;transition:transform 0.3s;"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = 'translateY(-4px)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = ''">
        <div style="font-size:48px;margin-bottom:8px;">{{ getStageEmoji(p.level) }}</div>
        <div style="font-weight:600;font-size:14px;">{{ p.name || p.student_name }}</div>
        <div style="font-size:12px;color:var(--color-text-secondary);">Lv.{{ p.level }} · {{ getStageName(p.level) }}</div>
        <div style="margin-top:8px;height:6px;background:var(--color-border);border-radius:3px;overflow:hidden;">
          <div style="height:100%;background:var(--gradient-primary);border-radius:3px;transition:width 0.5s;"
            :style="{ width: (Math.min(100, Math.round((p.exp % 100) / 100 * 100))) + '%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>
