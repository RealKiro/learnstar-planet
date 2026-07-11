<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { formatTime } from '@/utils/constants'
import type { ApiResponse } from '@/types'

interface Notice {
  id: number
  title: string
  content: string
  type: string
  published_at: string
  class_name: string
}

interface NoticeDetail extends Notice {
  publisher_name: string
}

const notices = ref<Notice[]>([])
const loading = ref(true)
const selected = ref<NoticeDetail | null>(null)
const detailLoading = ref(false)

const typeLabels: Record<string, string> = { info: '通知', homework: '作业', event: '活动', urgent: '紧急' }

function typeStyle(type: string) {
  const map: Record<string, { bg: string; color: string }> = {
    info: { bg: 'rgba(59,130,246,0.1)', color: '#3B82F6' },
    homework: { bg: 'rgba(16,185,129,0.1)', color: '#10B981' },
    event: { bg: 'rgba(245,158,11,0.1)', color: '#F59E0B' },
    urgent: { bg: 'rgba(239,68,68,0.1)', color: '#EF4444' },
  }
  return map[type] || { bg: 'rgba(100,116,139,0.1)', color: '#64748B' }
}

async function openDetail(n: Notice) {
  detailLoading.value = true
  selected.value = { ...n, publisher_name: '' }
  try {
    const res = await apiGet<ApiResponse<NoticeDetail>>(`/api/v1/parent/notices/${n.id}`)
    selected.value = res.data
  } catch {
    // keep basic info
  } finally {
    detailLoading.value = false
  }
}

function closeDetail() {
  selected.value = null
}

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Notice[]>>('/api/v1/parent/notices')
    notices.value = res.data || []
  } catch {
    notices.value = []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div style="margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">学校通知</h2>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="notices.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📢</div>
      <p>暂无通知</p>
    </div>

    <div v-else class="card" style="padding:8px;">
      <div v-for="n in notices" :key="n.id"
        style="display:flex;gap:16px;padding:16px;border-bottom:1px solid var(--color-border);cursor:pointer;align-items:center;"
        @click="openDetail(n)"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.background = 'rgba(79,70,229,0.03)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.background = ''">
        <span :style="{
          padding:'4px 10px', borderRadius:'4px', fontSize:'11px', fontWeight:600, whiteSpace:'nowrap',
          background: typeStyle(n.type).bg, color: typeStyle(n.type).color,
        }">{{ typeLabels[n.type] || n.type }}</span>
        <div style="flex:1;min-width:0;">
          <div style="font-weight:500;font-size:14px;">{{ n.title }}</div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">
            {{ n.class_name }} · {{ formatTime(n.published_at) }}
          </div>
        </div>
        <span style="color:var(--color-text-secondary);font-size:14px;">›</span>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="selected" style="position:fixed;inset:0;background:rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;z-index:100;padding:24px;" @click.self="closeDetail">
      <div class="card" style="max-width:560px;width:100%;padding:24px;max-height:80vh;overflow-y:auto;">
        <div style="display:flex;align-items:center;gap:8px;margin-bottom:16px;">
          <span :style="{
            padding:'4px 10px', borderRadius:'4px', fontSize:'11px', fontWeight:600,
            background: typeStyle(selected.type).bg, color: typeStyle(selected.type).color,
          }">{{ typeLabels[selected.type] || selected.type }}</span>
          <span style="font-size:12px;color:var(--color-text-secondary);">{{ selected.class_name }}</span>
        </div>
        <h3 style="font-size:18px;font-weight:700;margin-bottom:12px;">{{ selected.title }}</h3>
        <div v-if="detailLoading" style="text-align:center;padding:24px;color:var(--color-text-secondary);">加载中...</div>
        <div v-else>
          <p style="font-size:14px;line-height:1.7;color:var(--color-text);white-space:pre-wrap;">{{ selected.content }}</p>
          <div style="margin-top:16px;padding-top:16px;border-top:1px solid var(--color-border);display:flex;justify-content:space-between;font-size:12px;color:var(--color-text-secondary);">
            <span>发布人：{{ selected.publisher_name || '-' }}</span>
            <span>{{ formatTime(selected.published_at) }}</span>
          </div>
        </div>
        <div style="margin-top:20px;text-align:right;">
          <button class="btn btn-sm btn-primary" @click="closeDetail">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>
