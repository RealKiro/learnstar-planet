<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiPut } from '@/utils/api'
import type { ApiResponse, Notice } from '@/types'
import ModalGlass from '@/components/common/ModalGlass.vue'

const notices = ref<Notice[]>([])
const loading = ref(true)
const loadError = ref('')
const formError = ref('')

const showCreate = ref(false)
const publishStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const form = ref({ title: '', content: '', type: 'info' as string })

async function loadNotices() {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<Notice[]>>('/api/v1/teacher/notices')
    notices.value = res.data || []
    loadError.value = ''
  } catch {
    notices.value = []
    loadError.value = '通知加载失败'
  } finally { loading.value = false }
}

onMounted(loadNotices)

async function createNotice() {
  if (!form.value.title.trim() || !form.value.content.trim()) { formError.value = '请填写标题和内容'; return }
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

function typeClass(t?: string): string {
  return t === 'info' || t === 'event' || t === 'urgent' ? `type-badge--${t}` : 'type-badge--other'
}
</script>

<template>
  <div>
    <div class="page-header">
      <h2 class="page-title">📢 班级通知</h2>
      <button class="btn btn-sm btn-primary" @click="showCreate = true">发布通知</button>
    </div>

    <!-- 发布弹窗 -->
    <ModalGlass :visible="showCreate" @update:visible="showCreate = $event">
      <div class="modal-header">
        <h3 class="modal-title">发布通知</h3>
        <button class="modal-close" @click="showCreate = false">×</button>
      </div>
      <div class="form-group">
        <label>类型</label>
        <select v-model="form.type" class="form-select">
          <option value="info">📢 通知</option>
          <option value="event">🎉 活动</option>
          <option value="urgent">🔔 紧急</option>
        </select>
      </div>
      <div class="form-group">
        <label>标题</label>
        <input v-model="form.title" class="form-input" placeholder="通知标题" maxlength="200">
      </div>
      <div class="form-group form-group--lg">
        <label>内容</label>
        <textarea v-model="form.content" class="form-input content-textarea" placeholder="通知内容"></textarea>
      </div>
      <div v-if="formError" class="error-banner">{{ formError }}</div>
      <div class="modal-actions--wide">
        <button class="btn action-btn" @click="showCreate = false">取消</button>
        <button class="btn action-btn" :class="publishStatus === 'loading' ? 'btn-state-loading' : publishStatus === 'success' ? 'btn-state-success' : publishStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="publishStatus === 'loading'" @click="createNotice">
          {{ { idle: '确认发布', loading: '发布中...', success: '已发布 ✓', error: '失败' }[publishStatus] }}
        </button>
      </div>
    </ModalGlass>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="loadNotices">重试</button>
    </div>

    <div v-else-if="notices.length === 0" class="card empty-state">
      <div class="empty-state__icon">📢</div>
      <p>暂无通知</p>
    </div>

    <div v-else class="card notice-list">
      <div v-for="n in notices" :key="n.id" class="notice-item">
        <span class="type-badge" :class="typeClass(n.type)">{{ typeLabels[n.type] || n.type }}</span>
        <div class="notice-main">
          <div class="notice-title">{{ n.title }}</div>
          <div class="notice-content">
            {{ n.content?.substring(0, 100) }}<span v-if="n.content?.length > 100">...</span>
          </div>
          <div class="notice-date">{{ new Date(n.created_at).toLocaleString('zh-CN') }}</div>
        </div>
        <div v-if="n.read_count != null" class="notice-read">
          👁 {{ n.read_count }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-title { font-size: 24px; font-weight: 700; }

/* ===== 弹窗 ===== */
.form-group--lg { margin-bottom: 20px; }
.content-textarea { min-height: 80px; resize: vertical; }
.modal-actions--wide { display: flex; gap: 12px; }
.action-btn { flex: 1; }

/* ===== 通知列表 ===== */
.notice-list { padding: 16px; }
.notice-item { display: flex; gap: 16px; padding: 16px; border-bottom: 1px solid var(--color-border); }
.notice-item:last-child { border-bottom: none; }
.type-badge { padding: 4px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; white-space: nowrap; align-self: flex-start; }
.type-badge--info { background: rgba(59,130,246,0.1); color: #3B82F6; }
.type-badge--event { background: rgba(16,185,129,0.1); color: #10B981; }
.type-badge--urgent { background: rgba(239,68,68,0.1); color: #EF4444; }
.type-badge--other { background: rgba(100,116,139,0.1); color: #64748B; }
.notice-main { flex: 1; }
.notice-title { font-weight: 500; font-size: 14px; }
.notice-content { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.notice-date { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.notice-read { font-size: 12px; color: var(--color-text-secondary); display: flex; align-items: center; }
</style>
