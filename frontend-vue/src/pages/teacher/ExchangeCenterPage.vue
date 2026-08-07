<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse } from '@/types'
import PetSprite from '@/components/pet/PetSprite.vue'

interface WalletEntry {
  student_id: number
  student_name: string
  currency_type: string
  balance: number
}

interface StudentInfo {
  id: number
  name: string
  student_no?: string
  total_score: number
  pet_species?: string
  pet_level?: number
}

interface ExchangeRate {
  id: number
  name?: string
  from_currency: string
  to_currency: string
  rate: string
  is_active: boolean
}

const wallets = ref<WalletEntry[]>([])
const students = ref<StudentInfo[]>([])
const loading = ref(true)
const loadError = ref('')

const selectedStudent = ref<StudentInfo | null>(null)
const exchangeAmount = ref(10)
const exchangeTarget = ref<'science' | 'reading' | 'class_point'>('science')
const exchangeStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

// 真实汇率（来自教师端「汇率设定」，score → 各币种，is_active 生效）
const rates = ref<ExchangeRate[]>([])
function rateFor(target: string): number {
  const r = rates.value.find(r => r.from_currency === 'score' && r.to_currency === target && r.is_active)
  return r ? parseFloat(r.rate) : 0.5 // 默认 2:1（2 积分 = 1 币，防通胀）
}

async function loadAll() {
  loading.value = true
  try {
    const [walletRes, stuRes, rateRes] = await Promise.all([
      apiGet<ApiResponse<WalletEntry[]>>('/api/v1/teacher/currency/wallets'),
      apiGet<ApiResponse<StudentInfo[]>>('/api/v1/teacher/students?per_page=200'),
      apiGet<ApiResponse<ExchangeRate[]>>('/api/v1/teacher/exchange-rates'),
    ])
    wallets.value = walletRes.data || []
    students.value = stuRes.data || []
    rates.value = rateRes.data || []
    loadError.value = ''
  } catch {
    loadError.value = '数据加载失败'
  } finally { loading.value = false }
}

onMounted(loadAll)

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
    <div class="page-header">
      <h2 class="page-title">💱 兑换中心</h2>
      <span class="page-subtitle">积分按汇率兑换科学币 / 读书币 / 体育币</span>
    </div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="loadAll">重试</button>
    </div>

    <div v-else class="exchange-grid">
      <div class="card">
        <h3 class="card-title">学生钱包</h3>
        <div class="student-scroll">
          <div v-for="s in students" :key="s.id"
            class="student-row"
            :class="{ 'student-row--selected': selectedStudent?.id === s.id }"
            @click="selectStudent(s)">
            <!-- 宠物 -->
            <div class="pet-icon">
              <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level || 1" :animate="true" />
              <span v-else class="pet-icon--egg">🥚</span>
            </div>
            <!-- 姓名 + 学号 + 等级 -->
            <div class="student-info">
              <div class="student-name-row">
                <span class="student-name">{{ s.name }}</span>
                <span v-if="s.student_no" class="student-no">📛{{ s.student_no }}</span>
              </div>
              <div class="student-level">Lv.{{ s.pet_level || 0 }}</div>
            </div>
            <!-- 积分 + 钱包 -->
            <div class="student-wallet">
              <div>⭐ {{ s.total_score }}</div>
              <div class="wallet-line">🔬{{ getWallet(s.id, 'science') }} 📚{{ getWallet(s.id, 'reading') }} ⚽{{ getWallet(s.id, 'class_point') }}</div>
            </div>
          </div>
          <div v-if="students.length === 0" class="empty-students">暂无学生数据</div>
        </div>
      </div>

      <div class="card">
        <h3 class="card-title">兑换操作</h3>
        <div v-if="!selectedStudent" class="exchange-placeholder">
          <div class="exchange-placeholder__icon">👈</div>
          请从左侧选择学生
        </div>
        <div v-else>
          <div class="student-summary">
            <div class="student-summary__name">{{ selectedStudent.name }}</div>
            <div class="balance-grid">
              <div>⭐ 积分: <strong>{{ selectedStudent.total_score }}</strong></div>
              <div>🔬 科学币: <strong>{{ getWallet(selectedStudent.id, 'science') }}</strong></div>
              <div>📚 读书币: <strong>{{ getWallet(selectedStudent.id, 'reading') }}</strong></div>
              <div>⚽ 体育币: <strong>{{ getWallet(selectedStudent.id, 'class_point') }}</strong></div>
            </div>
          </div>

          <div class="form-group">
            <label>兑换目标</label>
            <div class="target-grid">
              <button class="target-btn" :class="{ 'target-btn--active': exchangeTarget === 'science' }" @click="exchangeTarget = 'science'">🔬 科学币</button>
              <button class="target-btn" :class="{ 'target-btn--active': exchangeTarget === 'reading' }" @click="exchangeTarget = 'reading'">📚 读书币</button>
              <button class="target-btn" :class="{ 'target-btn--active': exchangeTarget === 'class_point' }" @click="exchangeTarget = 'class_point'">⚽ 体育币</button>
            </div>
          </div>

          <div class="form-group amount-group">
            <label>兑换数量</label>
            <input v-model.number="exchangeAmount" type="number" min="1" :max="selectedStudent.total_score" class="form-input">
          </div>

          <div class="rate-note">
            消耗 <strong>{{ exchangeAmount }}</strong> 积分 → 获得 <strong>{{ Math.round(exchangeAmount * rateFor(exchangeTarget)) }}</strong> {{ exchangeTarget === 'science' ? '🔬 科学币' : exchangeTarget === 'reading' ? '📚 读书币' : '⚽ 体育币' }}
            <span>（汇率 1:{{ rateFor(exchangeTarget) }}）</span>
          </div>

          <button class="btn confirm-btn" :class="exchangeStatus === 'loading' ? 'btn-state-loading' : exchangeStatus === 'success' ? 'btn-state-success' : exchangeStatus === 'error' ? 'btn-state-error' : 'btn-solid'" :disabled="exchangeStatus !== 'idle' || exchangeAmount < 1 || !selectedStudent" @click="doExchange">
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

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-title { font-size: 24px; font-weight: 700; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); }

/* ===== 布局 ===== */
.exchange-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; }
.card-title { font-size: 16px; font-weight: 600; margin-bottom: 16px; }

/* ===== 学生列表 ===== */
.student-scroll { max-height: 500px; overflow-y: auto; }
.student-row {
  padding: 10px 12px;
  margin-bottom: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  position: relative;
}
.student-row--selected { border-color: var(--color-primary); background: rgba(79,70,229,0.04); }
.pet-icon { width: 44px; height: 44px; flex-shrink: 0; display: flex; align-items: center; justify-content: center; }
.pet-icon--egg { font-size: 24px; }
.student-info { flex: 1; min-width: 0; }
.student-name-row { display: flex; align-items: center; gap: 6px; }
.student-name { font-weight: 600; }
.student-no {
  font-size: 10px;
  color: var(--color-text-secondary);
  background: var(--tint-2);
  padding: 1px 6px;
  border-radius: 6px;
  letter-spacing: 0.3px;
}
.student-level { font-size: 11px; color: var(--color-text-secondary); }
.student-wallet { flex: 1; text-align: right; font-size: 12px; }
.wallet-line { color: var(--color-text-secondary); }
.empty-students { text-align: center; padding: 24px; color: var(--color-text-secondary); }

/* ===== 兑换操作 ===== */
.exchange-placeholder { text-align: center; padding: 48px; color: var(--color-text-secondary); }
.exchange-placeholder__icon { font-size: 32px; margin-bottom: 8px; }
.student-summary { padding: 16px; background: var(--color-bg); border-radius: var(--radius-md); margin-bottom: 16px; }
.student-summary__name { font-size: 14px; font-weight: 600; }
.balance-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-top: 8px; font-size: 13px; }
.amount-group { margin-bottom: 16px; }
.target-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 8px; }
.target-btn {
  padding: 10px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text);
}
.target-btn--active { background: var(--color-primary); color: #fff; border-color: var(--color-primary); }
.rate-note {
  padding: 12px;
  background: rgba(79,70,229,0.05);
  border-radius: var(--radius-md);
  margin-bottom: 16px;
  font-size: 13px;
  color: var(--color-text-secondary);
}
.confirm-btn { width: 100%; }
</style>
