<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface LeaveRecord { sp_no: string; leave_type: string | null; reason: string | null }
interface AttendanceRecord {
  id: number; student_id: number; student_name: string; student_no?: string
  status: 'present' | 'late' | 'leave' | 'absent'
  source: 'auto' | 'wechat_work' | 'manual'
  remark?: string; check_in_time?: string; leave_record?: LeaveRecord
}
interface Summary { present: number; late: number; leave: number; absent: number; rate: number; wechat_leave_count: number; manual_leave_count: number }

const records = ref<AttendanceRecord[]>([])
const summary = ref<Summary>({ present: 0, late: 0, leave: 0, absent: 0, rate: 0, wechat_leave_count: 0, manual_leave_count: 0 })
const loading = ref(true)
const attendanceStarted = ref(false)

const showLeaveModal = ref(false); const leaveStudentId = ref<number | null>(null); const leaveStudentName = ref(''); const leaveRemark = ref(''); const leaveSubmitting = ref(false)
const showAbsentModal = ref(false); const absentStudentId = ref<number | null>(null); const absentStudentName = ref(''); const absentRemark = ref(''); const absentSubmitting = ref(false)

const statusLabels: Record<string, string> = { present: '到课', late: '迟到', leave: '请假', absent: '缺勤' }
const statusColors: Record<string, string> = { present: 'var(--color-accent)', late: '#F59E0B', leave: '#3B82F6', absent: 'var(--color-danger)' }
function sourceLabel(s: string) { switch(s) { case 'wechat_work': return '企微'; case 'manual': return '手动'; case 'auto': return '自动'; default: return s } }
function sourceClass(s: string) { return 'source-tag ' + s }
const absentWithoutLeave = computed(() => records.value.filter(r => r.status === 'absent').length)

onMounted(loadData)

async function loadData() {
  loading.value = true
  try {
    const [recRes, sumRes] = await Promise.all([
      apiGet<ApiResponse<AttendanceRecord[]>>('/api/v1/teacher/attendance/today'),
      apiGet<ApiResponse<Summary>>('/api/v1/teacher/attendance/summary'),
    ])
    records.value = recRes.data || []
    summary.value = sumRes.data || { present: 0, late: 0, leave: 0, absent: 0, rate: 0, wechat_leave_count: 0, manual_leave_count: 0 }
    attendanceStarted.value = records.value.length > 0
  } catch { /* handled */ } finally { loading.value = false }
}

async function startAttendance() {
  try {
    const res = await apiPost<{ message: string }>('/api/v1/teacher/attendance/start', {})
    toast.show(res.message || '已开始', 'success'); attendanceStarted.value = true; await loadData()
  } catch { /* handled */ }
}

async function setStatus(studentId: number, status: string) {
  try { await apiPut(`/api/v1/teacher/attendance/${studentId}`, { status, source: 'manual' }); await loadData(); toast.show('已更新', 'success') } catch { /* handled */ }
}

function openLeaveModal(sid: number, sn: string) { leaveStudentId.value = sid; leaveStudentName.value = sn; leaveRemark.value = ''; showLeaveModal.value = true }
async function confirmLeave() {
  if (!leaveRemark.value.trim()) { toast.show('请填写请假原因', 'error'); return }
  leaveSubmitting.value = true
  try { await apiPost(`/api/v1/teacher/attendance/${leaveStudentId.value}/mark-leave`, { remark: leaveRemark.value.trim() }); toast.show(`已标记 ${leaveStudentName.value} 为请假`, 'success'); showLeaveModal.value = false; await loadData() } catch { /* handled */ } finally { leaveSubmitting.value = false }
}

function openAbsentModal(sid: number, sn: string) { absentStudentId.value = sid; absentStudentName.value = sn; absentRemark.value = ''; showAbsentModal.value = true }
async function confirmAbsent() {
  absentSubmitting.value = true
  try { await apiPost(`/api/v1/teacher/attendance/${absentStudentId.value}/mark-absent`, { remark: absentRemark.value.trim() || null }); toast.show(`已标记 ${absentStudentName.value} 为缺勤，建议联系家长`, 'warning'); showAbsentModal.value = false; await loadData() } catch { /* handled */ } finally { absentSubmitting.value = false }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">智能考勤</h2>
      <button class="btn btn-sm btn-primary" @click="startAttendance" :disabled="attendanceStarted">{{ attendanceStarted ? '已开始点名' : '开始点名' }}</button>
    </div>

    <div class="stats-grid">
      <div class="stat-card stat-card--accent"><span class="stat-card__icon">到</span><div class="stat-card__value">{{ summary.present }}</div><div class="stat-card__label">到课</div></div>
      <div class="stat-card stat-card--secondary"><span class="stat-card__icon">迟</span><div class="stat-card__value">{{ summary.late }}</div><div class="stat-card__label">迟到</div></div>
      <div class="stat-card stat-card--info"><span class="stat-card__icon">假</span><div class="stat-card__value">{{ summary.leave }}</div><div class="stat-card__label">请假 <span style="font-size:11px;color:#86868b;">(企微{{ summary.wechat_leave_count }})</span></div></div>
      <div class="stat-card" style="border-color:rgba(239,68,68,.3);"><span class="stat-card__icon">缺</span><div class="stat-card__value">{{ summary.absent }}</div><div class="stat-card__label">缺勤</div></div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else>
      <div v-if="absentWithoutLeave > 0" style="background:#FEF3C7;border:1px solid #FCD34D;border-radius:10px;padding:12px 16px;margin-bottom:16px;display:flex;align-items:center;gap:8px;">
        <span>!</span><span style="font-size:13px;color:#92400E;">有 {{ absentWithoutLeave }} 名学生未到校且无请假记录，请确认后联系家长</span>
      </div>

      <div class="data-table">
        <div class="data-table__header"><h3 style="font-size:16px;font-weight:600;">今日考勤</h3><span style="font-size:13px;color:var(--color-text-secondary);">出勤率 {{ summary.rate }}%</span></div>
        <table>
          <thead><tr><th>姓名</th><th>学号</th><th>状态</th><th>来源</th><th>备注</th><th>签到</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-if="records.length === 0"><td colspan="7" style="text-align:center;padding:48px;color:var(--color-text-secondary);"><div style="font-size:32px;margin-bottom:8px;">考勤</div>点击开始点名创建考勤记录</td></tr>
            <tr v-for="r in records" :key="r.student_id">
              <td style="font-weight:600;">{{ r.student_name }}</td>
              <td style="color:var(--color-text-secondary);">{{ r.student_no || '-' }}</td>
              <td><span :style="{ color: statusColors[r.status], fontWeight: 500 }">{{ statusLabels[r.status] }}</span></td>
              <td><span :class="sourceClass(r.source)">{{ sourceLabel(r.source) }}</span><span v-if="r.leave_record?.leave_type" style="font-size:11px;color:#86868b;margin-left:4px;">({{ r.leave_record.leave_type }})</span></td>
              <td style="color:var(--color-text-secondary);max-width:150px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" :title="r.remark || (r.leave_record?.reason || '')">{{ r.remark || r.leave_record?.reason || '-' }}</td>
              <td style="color:var(--color-text-secondary);">{{ r.check_in_time || '-' }}</td>
              <td style="display:flex;gap:4px;flex-wrap:wrap;">
                <button v-if="r.status !== 'present'" class="btn btn-sm" style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);" @click="setStatus(r.student_id, 'present')">到课</button>
                <button v-if="r.status !== 'late'" class="btn btn-sm" style="background:var(--color-bg);color:#F59E0B;border:1px solid rgba(245,158,11,.3);" @click="setStatus(r.student_id, 'late')">迟到</button>
                <button v-if="r.status !== 'leave'" class="btn btn-sm" style="background:var(--color-bg);color:#3B82F6;border:1px solid rgba(59,130,246,.3);" @click="openLeaveModal(r.student_id, r.student_name)">请假</button>
                <button v-if="r.status !== 'absent'" class="btn btn-sm" style="background:var(--color-bg);color:var(--color-danger);border:1px solid rgba(239,68,68,.3);" @click="openAbsentModal(r.student_id, r.student_name)">缺勤</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showLeaveModal" class="modal-overlay" @click.self="showLeaveModal = false">
      <div class="modal-card">
        <div class="modal-header"><h3>手动请假 - {{ leaveStudentName }}</h3><button class="modal-close" @click="showLeaveModal = false">&times;</button></div>
        <div class="modal-body">
          <p style="font-size:13px;color:#86868b;margin-bottom:12px;">该学生未在企业微信提交请假申请，请填写请假原因。</p>
          <label style="font-size:13px;font-weight:600;display:block;margin-bottom:6px;">请假原因 <span style="color:#EF4444;">*</span></label>
          <textarea v-model="leaveRemark" placeholder="例：家长电话告知感冒请假一天" rows="3" style="width:100%;padding:10px 12px;border:1px solid #e5e5ea;border-radius:10px;font-size:14px;resize:vertical;outline:none;"></textarea>
        </div>
        <div class="modal-footer">
          <button class="btn btn-sm" style="background:var(--color-bg);border:1px solid var(--color-border);" @click="showLeaveModal = false">取消</button>
          <button class="btn btn-sm btn-primary" @click="confirmLeave" :disabled="leaveSubmitting">{{ leaveSubmitting ? '提交中...' : '确认请假' }}</button>
        </div>
      </div>
    </div>

    <div v-if="showAbsentModal" class="modal-overlay" @click.self="showAbsentModal = false">
      <div class="modal-card">
        <div class="modal-header" style="background:#FEF2F2;"><h3 style="color:#991B1B;">确认缺勤 - {{ absentStudentName }}</h3><button class="modal-close" @click="showAbsentModal = false">&times;</button></div>
        <div class="modal-body">
          <p style="font-size:14px;color:#991B1B;margin-bottom:12px;background:#FEF2F2;padding:10px 12px;border-radius:8px;">该学生未到校且无请假记录，建议先联系家长后再标记缺勤。</p>
          <label style="font-size:13px;font-weight:600;display:block;margin-bottom:6px;">备注（可选）</label>
          <textarea v-model="absentRemark" placeholder="例：电话未接通，稍后再联系" rows="2" style="width:100%;padding:10px 12px;border:1px solid #e5e5ea;border-radius:10px;font-size:14px;resize:vertical;outline:none;"></textarea>
        </div>
        <div class="modal-footer">
          <button class="btn btn-sm" style="background:var(--color-bg);border:1px solid var(--color-border);" @click="showAbsentModal = false">取消</button>
          <button class="btn btn-sm" style="background:#EF4444;color:#fff;border:none;" @click="confirmAbsent" :disabled="absentSubmitting">{{ absentSubmitting ? '提交中...' : '确认缺勤' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.source-tag { display: inline-block; padding: 2px 8px; border-radius: 9999px; font-size: 11px; font-weight: 500 }
.source-tag.auto { background: #F3F4F6; color: #9CA3AF }
.source-tag.wechat_work { background: #DBEAFE; color: #1D4ED8 }
.source-tag.manual { background: #FEF3C7; color: #92400E }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,.4); display: flex; align-items: center; justify-content: center; z-index: 1000 }
.modal-card { background: #fff; border-radius: 16px; width: 100%; max-width: 440px; box-shadow: 0 20px 60px rgba(0,0,0,.15); overflow: hidden }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0f0f3 }
.modal-header h3 { font-size: 16px; font-weight: 700; margin: 0 }
.modal-close { background: none; border: none; font-size: 22px; color: #86868b; cursor: pointer; padding: 0; line-height: 1 }
.modal-close:hover { color: #1d1d1f }
.modal-body { padding: 20px }
.modal-footer { display: flex; gap: 8px; justify-content: flex-end; padding: 16px 20px; border-top: 1px solid #f0f0f3 }
</style>
