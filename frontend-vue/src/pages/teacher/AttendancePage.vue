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

const statusLabels: Record<string, string> = { present: '✅ 到课', late: '⏰ 迟到', leave: '📋 请假', absent: '❌ 缺勤' }
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
    toast.show(res.message || '考勤已开始', 'success'); attendanceStarted.value = true; await loadData()
  } catch { /* handled */ }
}

async function setStatus(studentId: number, status: string) {
  try { await apiPut(`/api/v1/teacher/attendance/${studentId}`, { status, source: 'manual' }); await loadData(); toast.show('状态已更新', 'success') } catch { /* handled */ }
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
      <h2 style="font-size:24px;font-weight:700;">✅ 智能考勤</h2>
      <button class="btn btn-sm btn-primary" @click="startAttendance" :disabled="attendanceStarted">{{ attendanceStarted ? '已开始点名' : '开始点名' }}</button>
    </div>

    <div class="stats-grid">
      <div class="stat-card stat-card--accent"><span class="stat-card__icon">✅</span><div class="stat-card__value">{{ summary.present }}</div><div class="stat-card__label">到课</div></div>
      <div class="stat-card stat-card--secondary"><span class="stat-card__icon">⏰</span><div class="stat-card__value">{{ summary.late }}</div><div class="stat-card__label">迟到</div></div>
      <div class="stat-card stat-card--info"><span class="stat-card__icon">📋</span><div class="stat-card__value">{{ summary.leave }}</div><div class="stat-card__label">请假 <span style="font-size:11px;color:#86868b;">(企微{{ summary.wechat_leave_count }})</span></div></div>
      <div class="stat-card" style="border-color:rgba(239,68,68,.3);"><span class="stat-card__icon">❌</span><div class="stat-card__value">{{ summary.absent }}</div><div class="stat-card__label">缺勤</div></div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else>
      <div v-if="absentWithoutLeave > 0" style="background:#