<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse, Notice } from '@/types'

const notices = ref<Notice[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Notice[]>>('/api/v1/teacher/notices')
    notices.value = res.data || []
  } catch {
    notices.value = []
  } finally {
    loading.value = false
  }
})

const typeLabels: Record<string, string> = { info: '通知', event: '活动', urgent: '紧急' }
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">班级通知</h2>
      <button class="btn btn-sm btn-primary">发布通知</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="notices.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📢</div>
      <p>暂无通知</p>
    </div>

    <div v-else class="card" style="padding:16px;">
      <div v-for="n in notices" :key="n.id"
        style="display:flex;gap:16px;padding:16px;border-bottom:1px solid var(--color-border);cursor:pointer;"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.background = 'rgba(79,70,229,0.03)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.background = ''">
        <span :style="{
          padding:'4px 8px', borderRadius:'4px', fontSize:'11px', fontWeight:600, whiteSpace:'nowrap',
          background: { info:'rgba(59,130,246,0.1)', event:'rgba(16,185,129,0.1)', urgent:'rgba(239,68,68,0.1)' }[n.type] || 'rgba(100,116,139,0.1)',
          color: { info:'#3B82F6', event:'#10B981', urgent:'#EF4444' }[n.type] || '#64748B',
        }">{{ typeLabels[n.type] || n.type }}</span>
        <div style="flex:1;">
          <div style="font-weight:500;font-size:14px;">{{ n.title }}</div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">{{ new Date(n.created_at).toLocaleString('zh-CN') }}</div>
        </div>
        <div v-if="n.read_count != null" style="font-size:12px;color:var(--color-text-secondary);display:flex;align-items:center;">
          👁 {{ n.read_count }}
        </div>
      </div>
    </div>
  </div>
</template>
