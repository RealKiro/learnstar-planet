<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse, TimetableChangeRequest } from '@/types'

const loading = ref(true)
const loadError = ref('')
const list = ref<TimetableChangeRequest[]>([])
const statusFilter = ref<'all' | 'pending' | 'approved' | 'rejected'>('pending')
const actionStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const actionId = ref<number | null>(null)

// 驳回弹窗
const showRejectModal = ref(false)
const rejectTarget = ref<TimetableChangeRequest | null>(null)
const rejectNote = ref('')
const rejectError = ref('')

const filtered = computed(() =>
  statusFilter.value === 'all' ? list.value : list.value.filter(r => r.status === statusFilter.value),
)

const pendingCount = computed(() => list.value.filter(r => r.status === 'pending').length)

const STATUS_LABELS: Record<string, string> = { pending: '待审核', approved: '已通过', rejected: '已驳回' }

onMounted(loadData)

async function loadData() {
  loading.value = true
  loadError.value = ''
  try {
    const res = await apiGet<ApiResponse<TimetableChangeRequest[]>>('/api/v1/admin/timetable/changes')
    list.value = res.data || []
  } catch {
    loadError.value = '申请列表加载失败，请刷新重试'
  } finally {
    loading.value = false
  }
}

async function approve(item: TimetableChangeRequest) {
  actionId.value = item.id
  actionStatus.value = 'loading'
  try {
    await apiPost<{ message: string }>(`/api/v1/admin/timetable/changes/${item.id}/approve`, {})
    actionStatus.value = 'success'
    setTimeout(() => { actionStatus.value = 'idle' }, 1200)
    await loadData()
  } catch {
    actionStatus.value = 'error'
    setTimeout(() => { actionStatus.value = 'idle' }, 2500)
  } finally {
    actionId.value = null
  }
}

function openRejectModal(item: TimetableChangeRequest) {
  rejectTarget.value = item
  rejectNote.value = ''
  rejectError.value = ''
  showRejectModal.value = true
}

async function confirmReject() {
  if (!rejectTarget.value) return
  if (!rejectNote.value.trim()) { rejectError.value = '请填写驳回原因，方便教师修改后重新提交'; return }
  actionStatus.value = 'loading'
  try {
    await apiPost<{ message: string }>(`/api/v1/admin/timetable/changes/${rejectTarget.value.id}/reject`, { note: rejectNote.value.trim() })
    actionStatus.value = 'success'
    showRejectModal.value = false
    setTimeout(() => { actionStatus.value = 'idle' }, 1200)
    await loadData()
  } catch {
    actionStatus.value = 'error'
    setTimeout(() => { actionStatus.value = 'idle' }, 2500)
  }
}

function formatTime(iso?: string | null): string {
  if (!iso) return '-'
  const d = new Date(iso)
  return isNaN(d.getTime()) ? '-' : d.toLocaleString('zh-CN', { month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div>
    <div class="page-head">
      <h2 class="page-title">课表审核</h2>
      <span v-if="pendingCount > 0" class="pending-chip">{{ pendingCount }} 条待审</span>
    </div>

    <div class="filter-row">
      <button v-for="s in (['pending', 'approved', 'rejected', 'all'] as const)" :key="s" :class="['filter-btn', { active: statusFilter === s }]" @click="statusFilter = s">
        {{ s === 'all' ? '全部' : STATUS_LABELS[s] }}
      </button>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>
    <div v-else-if="loadError" class="error-banner">{{ loadError }}</div>

    <div v-else class="card">
      <table class="req-table">
        <thead>
          <tr><th>班级</th><th>提交人</th><th>内容</th><th>提交时间</th><th>状态</th><th>审核信息</th><th>操作</th></tr>
        </thead>
        <tbody>
          <tr v-if="filtered.length === 0"><td colspan="7" class="empty-state">暂无{{ statusFilter === 'all' ? '' : STATUS_LABELS[statusFilter] }}的申请</td></tr>
          <tr v-for="r in filtered" :key="r.id">
            <td class="fw-600">{{ r.class_name || `班级#${r.class_id}` }}</td>
            <td>{{ r.requester_name || '-' }}</td>
            <td>修改 {{ r.entry_count }} 节排课</td>
            <td class="text-muted-13">{{ formatTime(r.created_at) }}</td>
            <td><span :class="['status-badge', r.status]">{{ STATUS_LABELS[r.status] }}</span></td>
            <td class="review-cell">
              <template v-if="r.status !== 'pending'">
                <span class="text-muted-13">{{ r.reviewer_name || '-' }} · {{ formatTime(r.reviewed_at) }}</span>
                <div v-if="r.review_note" class="review-note">{{ r.review_note }}</div>
              </template>
              <template v-else>-</template>
            </td>
            <td class="actions-row">
              <template v-if="r.status === 'pending'">
                <button class="btn btn-sm btn-ghost approve-btn" :disabled="actionStatus === 'loading' && actionId === r.id" @click="approve(r)">
                  {{ actionStatus === 'loading' && actionId === r.id ? '处理中...' : '通过' }}
                </button>
                <button class="btn btn-sm btn-outline-danger" :disabled="actionStatus === 'loading'" @click="openRejectModal(r)">驳回</button>
              </template>
              <span v-else class="text-muted-13">已处理</span>
            </td>
          </tr>
        </tbody>
      </table>
      <p class="muted-tip">「通过」会立即把该申请的完整课表应用到对应班级；同一班级的其他待审申请将自动作废。审批通过后教师导出 CSES 即可同步到 ClassIsLand。</p>
    </div>

    <div v-if="showRejectModal" class="modal-overlay" @click.self="showRejectModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h3>驳回申请 - {{ rejectTarget?.class_name || '' }}</h3>
          <button class="modal-close" @click="showRejectModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <label class="form-label">驳回原因 <span class="req-star-red">*</span></label>
          <textarea v-model="rejectNote" rows="3" class="form-textarea" placeholder="例：周三第 3 节与全校教研活动冲突，请调整"></textarea>
          <div v-if="rejectError" class="field-error">{{ rejectError }}</div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-sm btn-ghost" @click="showRejectModal = false">取消</button>
          <button class="btn btn-sm btn-danger" :disabled="actionStatus === 'loading'" @click="confirmReject">
            {{ actionStatus === 'loading' ? '提交中...' : '确认驳回' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pending-chip { font-size: 13px; font-weight: 600; background: #fef3c7; color: #92400e; padding: 4px 12px; border-radius: 9999px }
.filter-row { display: flex; gap: 8px; margin-bottom: 14px }
.filter-btn { border: 1.5px solid var(--color-border); background: transparent; border-radius: 9999px; padding: 5px 16px; font-size: 13px; font-weight: 600; cursor: pointer; color: var(--color-text-secondary); font-family: inherit }
.filter-btn.active { background: var(--color-text); border-color: var(--color-text); color: var(--color-bg-card, #fff) }
.req-table { width: 100%; border-collapse: collapse }
.req-table th { text-align: left; font-size: 12px; color: var(--color-text-secondary); font-weight: 600; padding: 10px 12px; border-bottom: 1px solid var(--color-border); white-space: nowrap }
.req-table td { padding: 12px; border-bottom: 1px solid var(--color-border); font-size: 14px; vertical-align: middle }
.review-cell { max-width: 220px }
.review-note { font-size: 12px; color: var(--color-text-secondary); margin-top: 2px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap }
.actions-row { white-space: nowrap }
.approve-btn { background: var(--md-primary, #7c3aed); color: #fff }
.muted-tip { font-size: 12px; color: #86868b; margin: 12px 2px 0 }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,.4); display: flex; align-items: center; justify-content: center; z-index: 1000 }
.modal-card { background: var(--color-bg-card, #fff); border-radius: 16px; width: 100%; max-width: 440px; box-shadow: 0 20px 60px rgba(0,0,0,.15); overflow: hidden }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0f0f3 }
.modal-header h3 { font-size: 16px; font-weight: 700; margin: 0 }
.modal-close { background: none; border: none; font-size: 22px; color: #86868b; cursor: pointer; padding: 0; line-height: 1 }
.modal-body { padding: 20px }
.modal-footer { display: flex; gap: 8px; justify-content: flex-end; padding: 16px 20px; border-top: 1px solid #f0f0f3 }
.form-label { font-size: 13px; font-weight: 600; display: block; margin-bottom: 6px }
.form-textarea { width: 100%; padding: 10px 12px; border: 1px solid #e5e5ea; border-radius: 10px; font-size: 14px; resize: vertical; outline: none }
.req-star-red { color: #ef4444 }
.status-badge { font-size: 12px; font-weight: 600; padding: 2px 10px; border-radius: 9999px; white-space: nowrap }
.status-badge.pending { background: #fef3c7; color: #92400e }
.status-badge.approved { background: #d1fae5; color: #065f46 }
.status-badge.rejected { background: #fee2e2; color: #991b1b }
.fw-600 { font-weight: 600 }
.btn-outline-danger { background: var(--color-bg); color: var(--color-danger); border: 1px solid rgba(239,68,68,.3) }
</style>
