<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse } from '@/types'

interface WalletEntry {
  student_id: number
  student_name: string
  currency_type: string
  balance: number
}

interface StudentInfo {
  id: number
  name: string
  total_score: number
}

const wallets = ref<WalletEntry[]>([])
const students = ref<StudentInfo[]>([])
const loading = ref(true)

const selectedStudent = ref<StudentInfo | null>(null)
const exchangeAmount = ref(10)
const exchangeTarget = ref<'science' | 'reading'>('science')
const exchangeStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

const rate = 1 // 1:1 兑换

onMounted(async () => {
  try {
    const [walletRes, stuRes] = await Promise.all([
      apiGet<ApiResponse<WalletEntry[]>>('/api/v1/teacher/currency/wallets'),
      apiGet<ApiResponse<StudentInfo[]>>('/api/v1/teacher/students?per_page=200'),
    ])
    wallets.value = walletRes.data || []
    students.value = stuRes.data || []
  } catch { /* */ } finally { loading.value = false }
})

function getWallet(id: number, type: string): number {
  return wallets.value.find(w => w.student_id === id && w.currency_type === type)?.balance || 0
}

function selectStudent(s: StudentInfo) {
  selectedStudent.value = s
  exchangeAmount.value = 10
}

async function doExchange() {
  if (!selectedStudent.value || exchangeAmount.value < 1) return
  if (exchangeAmount.value > (selectedStudent.value.total_score || 0)) {
    exchangeStatus.value = 'error'
    setTimeout(() => { exchangeStatus.value = 'idle' }, 3000)
    return
  }
  exchangeStatus.value = 'loading'
  try {
    await apiPost('/api/v1/teacher/currency/exchange', {
      student_id: selectedStudent.value.id,
      to_currency: exchangeTarget.value,
      amount: exchangeAmount.value,
    })
    exchangeStatus.value = 'success'
    // 刷新数据
    const res = await apiGet<ApiResponse<WalletEntry[]>>('/api/v1/teacher/currency/wallets')
    wallets.value = res.data || []
    selectedStudent.value.total_score -= exchangeAmount.value
    setTimeout(() => { exchangeStatus.value = 'idle' }, 1500)
  } catch {
    exchangeStatus.value = 'error'
    setTimeout(() => { exchangeStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">💱 兑换中心</h2>
      <span style="font-size:13px;color:var(--color-text-secondary);">积分兑换科学币 / 读书币</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else style="display:grid;grid-template-columns:1fr 1fr;gap:24px;">
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">学生钱包</h3>
        <div style="max-height:500px;overflow-y:auto;">
          <div v-for="s in students" :key="s.id"
            :style="{
              padding:'12px 16px', marginBottom:'8px', cursor:'pointer', display:'flex', alignItems:'center', gap:'12px',
              borderRadius:'var(--radius-md)', border:'1px solid var(--color-border)',
              ...(selectedStudent?.id === s.id ? { borderColor:'var(--color-primary)', background:'rgba(79,70,229,0.04)' } : {}),
            }"
            @click="selectStudent(s)">
            <div style="font-weight:600;min-width:60px;">{{ s.name }}</div>
            <div style="flex:1;text-align:right;font-size:13px;">
              <div>⭐ {{ s.total_score }}</div>
              <div style="color:var(--color-text-secondary);">🔬{{ getWallet(s.id, 'science') }} 📚{{ getWallet(s.id, 'reading') }}</div>
            </div>
          </div>
          <div v-if="students.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无学生数据</div>
        </div>
      </div>

      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">兑换操作</h3>
        <div v-if="!selectedStudent" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
          <div style="font-size:32px;margin-bottom:8px;">👈</div>
          请从左侧选择学生
        </div>
        <div v-else>
          <div style="padding:16px;background:var(--color-bg);border-radius:var(--radius-md);margin-bottom:16px;">
            <div style="font-size:14px;font-weight:600;">{{ selectedStudent.name }}</div>
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-top:8px;font-size:13px;">
              <div>⭐ 积分: <strong>{{ selectedStudent.total_score }}</strong></div>
              <div>🔬 科学币: <strong>{{ getWallet(selectedStudent.id, 'science_coin') }}</strong></div>
              <div>📚 读书币: <strong>{{ getWallet(selectedStudent.id, 'reading_coin') }}</strong></div>
            </div>
          </div>

          <div class="form-group" style="margin-bottom:12px;">
            <label>兑换目标</label>
            <div style="display:flex;gap:8px;">
              <button :style="exchangeTarget === 'science' ? { background:'var(--color-primary)', color:'white' } : {}"
                style="flex:1;padding:10px;border-radius:var(--radius-md);border:1px solid var(--color-border);background:var(--color-bg);cursor:pointer;font-size:14px;"
                @click="exchangeTarget = 'science'">🔬 科学币</button>
              <button :style="exchangeTarget === 'reading' ? { background:'var(--color-primary)', color:'white' } : {}"
                style="flex:1;padding:10px;border-radius:var(--radius-md);border:1px solid var(--color-border);background:var(--color-bg);cursor:pointer;font-size:14px;"
                @click="exchangeTarget = 'reading'">📚 读书币</button>
            </div>
          </div>

          <div class="form-group" style="margin-bottom:16px;">
            <label>兑换数量</label>
            <input v-model.number="exchangeAmount" type="number" min="1" :max="selectedStudent.total_score" class="form-input">
          </div>

          <div style="padding:12px;background:rgba(79,70,229,0.05);border-radius:var(--radius-md);margin-bottom:16px;font-size:13px;color:var(--color-text-secondary);">
            消耗 <strong>{{ exchangeAmount }}</strong> 积分 → 获得 <strong>{{ exchangeAmount }}</strong> {{ exchangeTarget === 'science' ? '🔬 科学币' : '📚 读书币' }}
          </div>

          <button class="btn btn-primary" style="width:100%;" :disabled="exchangeStatus !== 'idle' || exchangeAmount < 1 || !selectedStudent"
            :style="{ background: exchangeStatus === 'loading' ? '#f59e0b' : exchangeStatus === 'success' ? '#10b981' : exchangeStatus === 'error' ? '#ef4444' : '#7c3aed' }"
            @click="doExchange">
            <template v-if="exchangeStatus === 'idle'">确认兑换</template>
            <template v-else-if="exchangeStatus === 'loading'">兑换中...</template>
            <template v-else-if="exchangeStatus === 'success'">✅ 已兑换</template>
            <template v-else-if="exchangeStatus === 'error'">❌ 失败</template>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
