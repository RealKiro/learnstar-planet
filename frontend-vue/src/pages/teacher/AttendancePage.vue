<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface AttendanceRecord {
  student_id: number
  student_name: string
  student_no?: string
  status: 'present' | 'late' | 'leave' | 'absent'
  check_in_time?: string
}

const records = ref<AttendanceRecord[]>([])
const summary = ref({ present: 0, late: 0, leave: 0, absent: 0, rate: 0 })
const loading = ref(true)
const attendanceStarted = ref(false)

const statusLabels: Record<string, string> = { present: '✅ 到课', late: '⏰ 迟到', leave: '📋 请假', absent: '❌ 缺勤' }
const statusColors: Record<string, string> = { present: 'var(--color-accent)', late: '#F59E0B', leave: '#3B82F6', absent: 'var(--color-danger)' }

onMounted(loadData)

async function loadData() {
  loading.value = true
  try {
    const [recRes, sumRes] = await Promise.all([
      apiGet<ApiResponse<AttendanceRecord[]>>('/api/v1/teacher/attendance/today'),
      apiGet<ApiResponse<typeof summary.value>>('/api/v1/teacher/attendance/summary'),
    ])
    records.value = recRes.data || []
    summary.value = sumRes.data || { present: 0, late: 0, leave: 0, absent: 0, rate: 0 }
    attendanceStarted.value = records.value.length > 0
  } catch { /* handled */ }
  finally { loading.value = false }
}

async function startAttendance() {
  try {
    await apiPost('/api/v1/teacher/attendance/start', {})
    toast.show('✅ 考勤已开始，学生端将收到签到提醒', 'success')
    attendanceStarted.value = true
    await loadData()
  } catch { /* handled */ }
}

async function setStatus(studentId: number, status: string) {
  try {
    await apiPut(`/api/v1/teacher/attendance/${studentId}`, { status })
    const rec = records.value.find(r => r.student_id === studentId)
    if (rec) rec.status = status as AttendanceRecord['status']
    await loadData()
  } catch { /* handled */ }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">✅ 智能考勤</h2>
      <button class="btn btn-sm btn-primary" @click="startAttendance" :disabled="attendanceStarted">
        {{ attendanceStarted ? '已开始点名' : '开始点名' }}
      </button>
    </div>

    <!-- 统计卡 -->
    <div class="stats-grid">
      <div class="stat-card stat-card--accent">
        <span class="stat-card__icon">✅</span>
        <div class="stat-card__value">{{ summary.present }}</div>
        <div class="stat-card__label">到课</div>
      </div>
      <div class="stat-card stat-card--secondary">
        <span class="stat-card__icon">⏰</span>
        <div class="stat-card__value">{{ summary.late }}</div>
        <div class="stat-card__label">迟到</div>
      </div>
      <div class="stat-card stat-card--info">
        <span class="stat-card__icon">📋</span>
        <div class="stat-card__value">{{ summary.leave }}</div>
        <div class="stat-card__label">请假</div>
      </div>
      <div class="stat-card" style="border-color:rgba(239,68,68,0.3);">
        <span class="stat-card__icon">❌</span>
        <div class="stat-card__value">{{ summary.absent }}</div>
        <div class="stat-card__label">缺勤</div>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else>
      <div class="data-table">
        <div class="data-table__header">
          <h3 style="font-size:16px;font-weight:600;">今日考勤</h3>
          <span style="font-size:13px;color:var(--color-text-secondary);">出勤率 {{ summary.rate }}%</span>
        </div>
        <table>
          <thead><tr><th>姓名</th><th>学号</th><th>状态</th><th>签到时间</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-if="records.length === 0">
              <td colspan="5" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
                <div style="font-size:32px;margin-bottom:8px;">📋</div>
                点击「开始点名」创建考勤记录
              </td>
            </tr>
            <tr v-for="r in records" :key="r.student_id">
              <td style="font-weight:600;">{{ r.student_name }}</td>
              <td style="color:var(--color-text-secondary);">{{ r.student_no || '-' }}</td>
              <td><span :style="{ color: statusColors[r.status], fontWeight: 500 }">{{ statusLabels[r.status] }}</span></td>
              <td style="color:var(--color-text-secondary);">{{ r.check_in_time || '-' }}</td>
              <td style="display:flex;gap:4px;">
                <button v-if="r.status !== 'present'" class="btn btn-sm" style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);" @click="setStatus(r.student_id, 'present')">到课</button>
                <button v-if="r.status !== 'late'" class="btn btn-sm" style="background:var(--color-bg);color:#F59E0B;border:1px solid rgba(245,158,11,0.3);" @click="setStatus(r.student_id, 'late')">迟到</button>
                <button v-if="r.status !== 'leave'" class="btn btn-sm" style="background:var(--color-bg);color:#3B82F6;border:1px solid rgba(59,130,246,0.3);" @click="setStatus(r.student_id, 'leave')">请假</button>
                <button v-if="r.status !== 'absent'" class="btn btn-sm" style="background:var(--color-bg);color:var(--color-danger);border:1px solid rgba(239,68,68,0.3);" @click="setStatus(r.student_id, 'absent')">缺勤</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
