<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import axios from 'axios'
import type { ApiResponse } from '@/types'

interface GradeEntry {
  student_id: number
  student_name: string
  subjects: Record<string, number>
  total: number
  average: number
  rank: number
}

const grades = ref<GradeEntry[]>([])
const loading = ref(true)
const selectedExam = ref('')
const examOptions = ref<string[]>([])

const subjects = ref<string[]>([])

// 录入弹窗
const showInputModal = ref(false)
const inputForm = ref({ exam_name: '', subject: '', grades: [] as { student_id: number; score: number }[] })
const inputStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const exportStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

const avgScore = computed(() => {
  if (!grades.value.length) return 0
  return Math.round(grades.value.reduce((a, g) => a + g.total, 0) / grades.value.length * 10) / 10
})
const maxScore = computed(() => grades.value.length ? Math.max(...grades.value.map(g => g.total)) : 0)
const minScore = computed(() => grades.value.length ? Math.min(...grades.value.map(g => g.total)) : 0)
const passRate = computed(() => {
  if (!grades.value.length) return 0
  const pass = grades.value.filter(g => g.average >= 60).length
  return Math.round(pass / grades.value.length * 100)
})

async function loadGrades() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    if (selectedExam.value) params.exam_name = selectedExam.value
    const res = await apiGet<ApiResponse<GradeEntry[]>>('/api/v1/teacher/grades', { params })
    grades.value = res.data || []
  } catch { grades.value = [] } finally { loading.value = false }
}

onMounted(async () => {
  await loadGrades()
  // 提取已有考试名称
  if (grades.value.length > 0) {
    // 实际应该从后端获取，这里简单处理
  }
})

async function exportGrades() {
  exportStatus.value = 'loading'
  try {
    const token = localStorage.getItem('auth_token')
    const res = await axios.get('/api/v1/teacher/reports/export/scores', {
      responseType: 'blob',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
    const blob = new Blob([res.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url; a.download = `成绩报表.xlsx`
    document.body.appendChild(a); a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    exportStatus.value = 'success'
    setTimeout(() => { exportStatus.value = 'idle' }, 1500)
  } catch {
    exportStatus.value = 'error'
    setTimeout(() => { exportStatus.value = 'idle' }, 3000)
  }
}

function openInputModal() {
  inputForm.value = { exam_name: '', subject: '', grades: [] }
  showInputModal.value = true
}

async function submitGrades() {
  if (!inputForm.value.exam_name || !inputForm.value.subject) { inputStatus.value = 'error'; setTimeout(() => { inputStatus.value = 'idle' }, 3000); return }
  inputStatus.value = 'loading'
  try {
    await apiPost('/api/v1/teacher/grades', inputForm.value)
    inputStatus.value = 'success'
    selectedExam.value = inputForm.value.exam_name
    await loadGrades()
    setTimeout(() => { inputStatus.value = 'idle'; showInputModal.value = false }, 1500)
  } catch {
    inputStatus.value = 'error'
    setTimeout(() => { inputStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">📊 成绩管理</h2>
      <div style="display:flex;gap:8px;align-items:center;">
        <select v-model="selectedExam" @change="loadGrades" class="form-select" style="width:auto;">
          <option value="">全部考试</option>
          <option v-for="e in examOptions" :key="e" :value="e">{{ e }}</option>
        </select>
        <button class="btn btn-sm btn-primary" @click="openInputModal">录入成绩</button>
        <button class="btn btn-sm" @click="exportGrades" :disabled="exportStatus !== 'idle'"
          :style="{ background: exportStatus === 'loading' ? '#f59e0b' : exportStatus === 'success' ? '#10b981' : exportStatus === 'error' ? '#ef4444' : '#7c3aed', color: '#fff', border: '1px solid transparent' }">
          <template v-if="exportStatus === 'idle'">📤 导出</template>
          <template v-else-if="exportStatus === 'loading'">导出中...</template>
          <template v-else-if="exportStatus === 'success'">✅ 已导出</template>
          <template v-else-if="exportStatus === 'error'">❌ 失败</template>
        </button>
      </div>
    </div>

    <!-- 录入弹窗 -->
    <Teleport to="body">
      <div v-if="showInputModal" @click.self="showInputModal = false"
        style="position:fixed;inset:0;z-index:999;background:rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:var(--color-bg-card);border-radius:16px;padding:24px;max-width:420px;width:100%;box-shadow:0 20px 60px rgba(0,0,0,0.15);">
          <h3 style="font-size:18px;font-weight:700;margin-bottom:16px;">录入成绩</h3>
          <p style="font-size:13px;color:#64748B;margin-bottom:16px;">成绩导入功能开发中，当前支持通过后端 API 录入。</p>
          <div class="form-group" style="margin-bottom:12px;">
            <label>考试名称</label>
            <input v-model="inputForm.exam_name" class="form-input" placeholder="如：期中考试">
          </div>
          <div class="form-group" style="margin-bottom:12px;">
            <label>科目</label>
            <input v-model="inputForm.subject" class="form-input" placeholder="如：语文">
          </div>
          <div style="display:flex;gap:12px;margin-top:20px;">
            <button class="btn" style="flex:1;" @click="showInputModal = false">取消</button>
            <button class="btn btn-primary" style="flex:1;" :disabled="inputStatus !== 'idle'"
              :style="{ background: inputStatus === 'loading' ? '#f59e0b' : inputStatus === 'success' ? '#10b981' : inputStatus === 'error' ? '#ef4444' : '#7c3aed' }"
              @click="submitGrades">
              <template v-if="inputStatus === 'idle'">提交</template>
              <template v-else-if="inputStatus === 'loading'">提交中...</template>
              <template v-else-if="inputStatus === 'success'">✅ 成功</template>
              <template v-else-if="inputStatus === 'error'">❌ 失败</template>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="grades.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📊</div>
      <p style="margin-bottom:8px;">暂无成绩数据</p>
      <p style="font-size:13px;">点击「录入成绩」添加考试数据</p>
    </div>

    <template v-else>
      <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:12px;margin-bottom:24px;">
        <div class="card" style="text-align:center;padding:20px;">
          <div style="font-size:28px;font-weight:700;color:var(--color-primary);">{{ avgScore }}</div>
          <div style="font-size:13px;color:#64748B;">班级平均分</div>
        </div>
        <div class="card" style="text-align:center;padding:20px;">
          <div style="font-size:28px;font-weight:700;color:var(--color-accent);">{{ maxScore }}</div>
          <div style="font-size:13px;color:#64748B;">最高分</div>
        </div>
        <div class="card" style="text-align:center;padding:20px;">
          <div style="font-size:28px;font-weight:700;color:var(--color-danger);">{{ minScore }}</div>
          <div style="font-size:13px;color:#64748B;">最低分</div>
        </div>
        <div class="card" style="text-align:center;padding:20px;">
          <div style="font-size:28px;font-weight:700;color:var(--color-secondary);">{{ passRate }}%</div>
          <div style="font-size:13px;color:#64748B;">及格率</div>
        </div>
      </div>

      <div class="data-table">
        <table>
          <thead><tr><th>姓名</th><th>总分</th><th>平均分</th><th>班名次</th></tr></thead>
          <tbody>
            <tr v-for="g in grades" :key="g.student_id">
              <td style="font-weight:600;">{{ g.student_name }}</td>
              <td style="font-weight:700;">{{ g.total }}</td>
              <td style="color:var(--color-text-secondary);">{{ g.average }}</td>
              <td :style="{ color: g.rank <= 3 ? '#EAB308' : '', fontWeight: 700 }">{{ g.rank }}</td>
            </tr>
            <tr v-if="grades.length === 0">
              <td colspan="4" style="text-align:center;padding:24px;color:#64748B;">暂无数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>
