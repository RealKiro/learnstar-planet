<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import type { ApiResponse, Notice } from '@/types'
import ModalGlass from '@/components/common/ModalGlass.vue'

const notices = ref<Notice[]>([])
const loading = ref(true)
const loadError = ref('')
const meta = ref({ current_page: 1, last_page: 1, per_page: 20, total: 0 })

const pageLabel = computed(() => {
  const { current_page, last_page, total } = meta.value
  return last_page > 1 ? `第 ${current_page} / ${last_page} 页 · 共 ${total} 条` : `共 ${total} 条`
})

// 发布
const showCreate = ref(false)
const publishStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const form = ref({ title: '', content: '', type: 'info' as string })
const formError = ref('')

// 编辑
const showEdit = ref(false)
const editingId = ref<number | null>(null)
const editStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const editForm = ref({ title: '', content: '', type: 'info' as string })
const editError = ref('')

// 列表操作状态（publish/unpublish/delete 按 id 内联反馈）
const opState = ref<Record<string, 'idle' | 'loading' | 'success' | 'error'>>({})

function stateCls(action: string, id: number): string {
  const s = opState.value[`${action}-${id}`]
  return s === 'loading' ? 'btn-state-loading' : s === 'success' ? 'btn-state-success' : s === 'error' ? 'btn-state-error' : ''
}

async function loadNotices(resetPage = false) {
  if (resetPage) meta.value.current_page = 1
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (meta.value.current_page > 1) params.set('page', String(meta.value.current_page))
    const res = await apiGet<ApiResponse<Notice[]>>(`/api/v1/teacher/notices?${params.toString()}`)
    notices.value = res.data || []
    loadError.value = ''
    if (res.meta) meta.value = res.meta
  } catch {
    notices.value = []
    loadError.value = '通知加载失败'
  } finally { loading.value = false }
}

function changePage(page: number) {
  const { current_page, last_page } = meta.value
  if (page < 1 || page > last_page || page === current_page) return
  meta.value.current_page = page
  loadNotices()
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
    await loadNotices(true)
  } catch {
    publishStatus.value = 'error'
    setTimeout(() => { publishStatus.value = 'idle' }, 3000)
  }
}

function openEditModal(n: Notice) {
  editingId.value = n.id
  editError.value = ''
  editForm.value = { title: n.title, content: n.content, type: n.type }
  showEdit.value = true
}

async function submitEdit() {
  if (!editForm.value.title.trim() || !editForm.value.content.trim()) { editError.value = '请填写标题和内容'; return }
  editStatus.value = 'loading'
  try {
    await apiPut(`/api/v1/teacher/notices/${editingId.value}`, {
      title: editForm.value.title.trim(),
      content: editForm.value.content.trim(),
      type: editForm.value.type,
    })
    editStatus.value = 'success'
    showEdit.value = false
    await loadNotices()
    setTimeout(() => { editStatus.value = 'idle' }, 1500)
  } catch {
    editStatus.value = 'error'
    setTimeout(() => { editStatus.value = 'idle' }, 3000)
  }
}

async function runOp(action: 'publish' | 'unpublish' | 'delete', n: Notice) {
  if (action === 'delete') {
    const ok = await openConfirm({ title: '删除通知', message: `确定删除通知「${n.title}」？`, danger: true, confirmText: '确认删除' })
    if (!ok) return
  }
  const key = `${action}-${n.id}`
  opState.value[key] = 'loading'
  try {
    if (action === 'publish') await apiPut(`/api/v1/teacher/notices/${n.id}/publish`, {})
    else if (action === 'unpublish') await apiPut(`/api/v1/teacher/notices/${n.id}/unpublish`, {})
    else await apiDelete(`/api/v1/teacher/notices/${n.id}`)
    opState.value[key] = 'success'
    await loadNotices()
    setTimeout(() => { delete opState.value[key] }, 1500)
  } catch {
    opState.value[key] = 'error'
    setTimeout(() => { delete opState.value[key] }, 3000)
  }
}

const typeLabels: Record<string, string> = { info: '通知', event: '活动', urgent: '紧急', homework: '作业' }

function typeClass(t?: string): string {
  return t === 'info' || t === 'event' || t === 'urgent' || t === 'homework' ? `type-badge--${t}` : 'type-badge--other'
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
          <option value="homework">📝 作业</option>
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

    <!-- 编辑弹窗 -->
    <ModalGlass :visible="showEdit" @update:visible="showEdit = $event">
      <div class="modal-header">
        <h3 class="modal-title">编辑通知</h3>
        <button class="modal-close" @click="showEdit = false">×</button>
      </div>
      <div class="form-group">
        <label>类型</label>
        <select v-model="editForm.type" class="form-select">
          <option value="info">📢 通知</option>
          <option value="homework">📝 作业</option>
          <option value="event">🎉 活动</option>
          <option value="urgent">🔔 紧急</option>
        </select>
      </div>
      <div class="form-group">
        <label>标题</label>
        <input v-model="editForm.title" class="form-input" placeholder="通知标题" maxlength="200">
      </div>
      <div class="form-group form-group--lg">
        <label>内容</label>
        <textarea v-model="editForm.content" class="form-input content-textarea" placeholder="通知内容"></textarea>
      </div>
      <div v-if="editError" class="error-banner">{{ editError }}</div>
      <div class="modal-actions">
        <button class="btn btn-sm btn-ghost-card" @click="showEdit = false">取消</button>
        <button class="btn btn-sm" :class="editStatus === 'loading' ? 'btn-state-loading' : editStatus === 'success' ? 'btn-state-success' : editStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="editStatus === 'loading'" @click="submitEdit">{{ editStatus === 'loading' ? '保存中...' : editStatus === 'success' ? '已保存 ✓' : editStatus === 'error' ? '保存失败' : '保存' }}</button>
      </div>
    </ModalGlass>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="loadNotices(true)">重试</button>
    </div>

    <div v-else-if="notices.length === 0" class="card empty-state">
      <div class="empty-state__icon">📢</div>
      <p>暂无通知</p>
    </div>

    <div v-else class="card notice-list">
      <div v-for="n in notices" :key="n.id" class="notice-item">
        <span class="type-badge" :class="typeClass(n.type)">{{ typeLabels[n.type] || n.type }}</span>
        <div class="notice-main">
          <div class="notice-title">
            {{ n.title }}
            <span class="publish-badge" :class="n.is_published ? 'publish-badge--published' : 'publish-badge--draft'">{{ n.is_published ? '已发布' : '草稿' }}</span>
          </div>
          <div class="notice-content">
            {{ n.content?.substring(0, 100) }}<span v-if="n.content?.length > 100">...</span>
          </div>
          <div class="notice-date">{{ new Date(n.created_at).toLocaleString('zh-CN') }}</div>
        </div>
        <div class="notice-actions">
          <button class="btn btn-sm btn-mini" @click="openEditModal(n)">编辑</button>
          <button v-if="n.is_published" class="btn btn-sm btn-mini" :class="stateCls('unpublish', n.id)" @click="runOp('unpublish', n)">{{ opState['unpublish-' + n.id] === 'loading' ? '撤回中...' : '撤回' }}</button>
          <button v-else class="btn btn-sm btn-mini" :class="stateCls('publish', n.id)" @click="runOp('publish', n)">{{ opState['publish-' + n.id] === 'loading' ? '发布中...' : '发布' }}</button>
          <button class="btn btn-sm btn-mini" :class="opState['delete-' + n.id] === 'loading' ? 'btn-state-loading' : opState['delete-' + n.id] === 'success' ? 'btn-state-success' : opState['delete-' + n.id] === 'error' ? 'btn-state-error' : 'btn-danger'" @click="runOp('delete', n)">{{ opState['delete-' + n.id] === 'loading' ? '删除中...' : '删除' }}</button>
        </div>
      </div>
      <div v-if="meta.last_page > 1" class="pagination">
        <button class="btn btn-sm btn-ghost-card" :disabled="meta.current_page <= 1" @click="changePage(meta.current_page - 1)">← 上一页</button>
        <span class="pagination__info">{{ pageLabel }}</span>
        <button class="btn btn-sm btn-ghost-card" :disabled="meta.current_page >= meta.last_page" @click="changePage(meta.current_page + 1)">下一页 →</button>
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
.notice-item { display: flex; gap: 16px; padding: 16px; border-bottom: 1px solid var(--color-border); align-items: flex-start; }
.notice-item:last-child { border-bottom: none; }
.type-badge { padding: 4px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; white-space: nowrap; align-self: flex-start; }
.type-badge--info { background: rgba(59,130,246,0.1); color: #3B82F6; }
.type-badge--event { background: rgba(16,185,129,0.1); color: #10B981; }
.type-badge--urgent { background: rgba(239,68,68,0.1); color: #EF4444; }
.type-badge--homework { background: rgba(245,158,11,0.1); color: #D97706; }
.type-badge--other { background: rgba(100,116,139,0.1); color: #64748B; }
.notice-main { flex: 1; min-width: 0; }
.notice-title { font-weight: 500; font-size: 14px; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.publish-badge { font-size: 10px; font-weight: 600; padding: 1px 6px; border-radius: 10px; }
.publish-badge--published { background: rgba(16,185,129,0.12); color: var(--color-success-text); }
.publish-badge--draft { background: var(--tint-2); color: var(--color-text-secondary); }
.notice-content { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.notice-date { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.notice-actions { display: flex; gap: 4px; align-items: center; }

/* ===== 次要按钮 ===== */
.btn-mini {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  font-size: 12px;
}

/* ===== 分页 ===== */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 14px 0 4px;
}
.pagination__info { font-size: 13px; color: var(--color-text-secondary); }
.pagination .btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
