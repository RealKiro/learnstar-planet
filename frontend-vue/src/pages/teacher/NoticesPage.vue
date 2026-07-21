<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, Notice } from '@/types'

const toast = useToastStore()
const notices = ref<Notice[]>([])
const loading = ref(true)

const showCreate = ref(false)
const publishStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const form = ref({ title: '', content: '', type: 'info' as string })

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Notice[]>>('/api/v1/teacher/notices')
    notices.value = res.data || []
  } catch { notices.value = [] } finally { loading.value = false }
})

async function createNotice() {
  if (!form.value.title.trim() || !form.value.content.trim()) { toast.show('请填写标题和内容', 'error', { position: 'center', duration: 2000 }); return }
  publishStatus.value = 'loading'
  try {
    const res = await apiPost<ApiResponse<{ id: number }>>('/api/v1/teacher/notices', {
      title: form.value.title.trim(),
      content: form.value.content.trim(),
      type: form.value.type,
    })
    // 自动发布
    await apiPut(`/api/v1/teacher/notices/${res.data.id}/publish`, {})
    publishStatus.value = 'success'
    setTimeout(() => { publishStatus.value = 'idle' }, 1500)
    showCreate.value = false
    form.value = { title: '', content: '', type: 'info' }
    const r = await apiGet<ApiResponse<Notice[]>>('/api/v1/teacher/notices')
    notices.value = r.data || []
  } catch {
    publishStatus.value = 'error'
    setTimeout(() => { publishStatus.value = 'idle' }, 3000)
  }
}

const typeLabels: Record<string, string> = { info: '通知', event: '活动', urgent: '紧急' }
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">📢 班级通知</h2>
      <button class="btn btn-sm btn-primary" @click="showCreate = true">发布通知</button>
    </div>

    <!-- 发布弹窗 -->
    <Teleport to="body">
      <div v-if="showCreate" @click.self="showCreate = false"
        style="position:fixed;inset:0;z-index:999;background:rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:var(--color-bg-card);border-radius:16px;padding:24px;max-width:480px;width:100%;box-shadow:0 20px 60px rgba(0,0,0,0.15);">
          <h3 style="font-size:18px;font-weight:700;margin-bottom:16px;">发布通知</h3>
          <div class="form-group" style="margin-bottom:12px;">
            <label>类型</label>
            <select v-model="form.type" class="form-select" style="width:100%;padding:10px;border-radius:10px;border:1px solid var(--color-border);">
              <option value="info">📢 通知</option>
              <option value="event">🎉 活动</option>
              <option value="urgent">🔔 紧急</option>
            </select>
          </div>
          <div class="form-group" style="margin-bottom:12px;">
            <label>标题</label>
            <input v-model="form.title" class="form-input" placeholder="通知标题" maxlength="200">
          </div>
          <div class="form-group" style="margin-bottom:20px;">
            <label>内容</label>
            <textarea v-model="form.content" class="form-input" style="min-height:80px;resize:vertical;" placeholder="通知内容"></textarea>
          </div>
          <div style="display:flex;gap:12px;">
            <button class="btn" style="flex:1;" @click="showCreate = false">取消</button>
            <button class="btn" style="flex:1;color:#fff;border:none;" :style="{ background: { idle: '#7c3aed', loading: '#f59e0b', success: '#10b981', error: '#ef4444' }[publishStatus] }" :disabled="publishStatus === 'loading'" @click="createNotice">
              {{ { idle: '确认发布', loading: '发布中...', success: '已发布 ✓', error: '失败' }[publishStatus] }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="notices.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📢</div>
      <p>暂无通知</p>
    </div>

    <div v-else class="card" style="padding:16px;">
      <div v-for="n in notices" :key="n.id"
        style="display:flex;gap:16px;padding:16px;border-bottom:1px solid var(--color-border);">
        <span :style="{
          padding:'4px 8px', borderRadius:'4px', fontSize:'11px', fontWeight:600, whiteSpace:'nowrap',
          background: { info:'rgba(59,130,246,0.1)', event:'rgba(16,185,129,0.1)', urgent:'rgba(239,68,68,0.1)' }[n.type] || 'rgba(100,116,139,0.1)',
          color: { info:'#3B82F6', event:'#10B981', urgent:'#EF4444' }[n.type] || '#64748B',
        }">{{ typeLabels[n.type] || n.type }}</span>
        <div style="flex:1;">
          <div style="font-weight:500;font-size:14px;">{{ n.title }}</div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">
            {{ n.content?.substring(0, 100) }}<span v-if="n.content?.length > 100">...</span>
          </div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">{{ new Date(n.created_at).toLocaleString('zh-CN') }}</div>
        </div>
        <div v-if="n.read_count != null" style="font-size:12px;color:var(--color-text-secondary);display:flex;align-items:center;">
          👁 {{ n.read_count }}
        </div>
      </div>
    </div>
  </div>
</template>
