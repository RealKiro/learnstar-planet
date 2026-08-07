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

interface ExchangeLogItem {
  id: number
  student_id: number
  student_name?: string
  student_no?: string
  from_currency: string
  to_currency: string
  from_amount: number
  to_amount: number
  created_at: string
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

// ===== 学生列表（搜索 + 分页，后端 search/page 已支持） =====
const studentSearch = ref('')
const studentsLoading = ref(false)
const studentsMeta = ref({ current_page: 1, last_page: 1, per_page: 50, total: 0 })

// ===== 兑换记录 =====
const logs = ref<ExchangeLogItem[]>([])
const logsLoading = ref(false)
const logsError = ref('')
const logsMeta = ref({ current_page: 1, last_page: 1, per_page: 20, total: 0 })

const currencyLabel: Record<string, string> = { science: '科学币', reading: '读书币', class_point: '体育币' }
const currencyIcon: Record<string, string> = { science: '🔬', reading: '📚', class_point: '⚽' }
const fmtCurrency = (c: string) => `${currencyIcon[c] || '·'} ${currencyLabel[c] || c}`

async function loadStudents(resetPage = false) {
  if (resetPage) studentsMeta.value.current_page = 1
  studentsLoading.value = true
  try {
    const params = new URLSearchParams()
    params.set('per_page', String(studentsMeta.value.per_page))
    if (studentsMeta.value.current_page > 1) params.set('page', String(studentsMeta.value.current_page))
    if (studentSearch.value.trim()) params.set('search', studentSearch.value.trim())
    const res = await apiGet<ApiResponse<StudentInfo[]>>(`/api/v1/teacher/students?${params.toString()}`)
    students.value = res.data || []
    if (res.meta) studentsMeta.value = res.meta
  } catch { /* 列表失败不阻塞兑换 */ }
  finally { studentsLoading.value = false }
}

async function loadLogs(resetPage = false) {
  if (resetPage) logsMeta.value.current_page = 1
  logsLoading.value = true
  logsError.value = ''
  try {
    const params = new URLSearchParams()
    if (logsMeta.value.current_page > 1) params.set('page', String(logsMeta.value.current_page))
    const res = await apiGet<ApiResponse<ExchangeLogItem[]>>(`/api/v1/teacher/currency/exchange-logs?${params.toString()}`)
    logs.value = res.data || []
    if (res.meta) logsMeta.value = res.meta
  } catch {
    logsError.value = '兑换记录加载失败'
  } finally { logsLoading.value = false }
}

function changeStudentsPage(page: number) {
  const { current_page, last_page } = studentsMeta.value
  if (page < 1 || page > last_page || page === current_page) return
  studentsMeta.value.current_page = page
  loadStudents()
}

function changeLogsPage(page: number) {
  const { current_page, last_page } = logsMeta.value
  if (page < 1 || page > last_page || page === current_page) return
  logsMeta.value.current_page = page
  loadLogs()
}

async function loadAll() {
  loading.value = true
  try {
    const [walletRes, rateRes] = await Promise.all([
      apiGet<ApiResponse<WalletEntry[]>>('/api/v1/teacher/currency/wallets'),
      apiGet<ApiResponse<ExchangeRate[]>>('/api/v1/teacher/exchange-rates'),
    ])
    wallets.value = walletRes.data || []
    rates.value = rateRes.data || []
    loadError.value = ''
  } catch {
    loadError.value = '数据加载失败'
  } finally { loading.value = false }
  await Promise.all([loadStudents(true), loadLogs(true)])
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
    await loadLogs(true)
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
        <div class="card-head">
          <h3 class="card-title">学生钱包</h3>
          <input
            v-model="studentSearch"
            class="form-input search-input"
            placeholder="🔍 搜索学生"
            @keyup.enter="loadStudents(true)"
            @keyup.esc="studentSearch = ''"
          >
        </div>
        <div class="student-scroll">
          <div v-if="studentsLoading" class="list-loading">加载中...</div>
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
          <div v-if="students.length === 0 && !studentsLoading" class="empty-students">暂无学生数据</div>
          <div v-if="studentsMeta.last_page > 1" class="mini-pagination">
            <button class="btn btn-sm btn-ghost-card" :disabled="studentsMeta.current_page <= 1" @click="changeStudentsPage(studentsMeta.current_page - 1)">←</button>
            <span class="mini-pagination__info">{{ studentsMeta.current_page }} / {{ studentsMeta.last_page }}</span>
            <button class="btn btn-sm btn-ghost-card" :disabled="studentsMeta.current_page >= studentsMeta.last_page" @click="changeStudentsPage(studentsMeta.current_page + 1)">→</button>
          </div>
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
            消耗 <strong>{{ exchangeAmount }}</strong> 积分 → 获得 <strong>{{ Math.round(exchangeAmount * rateFor(exchangeTarget)) }}</strong> {{ fmtCurrency(exchangeTarget) }}
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

    <!-- 兑换记录 -->
    <div v-if="!loading && !loadError" class="card logs-card">
      <div class="card-head">
        <h3 class="card-title">兑换记录</h3>
        <button v-if="logsError" class="btn btn-sm btn-ghost-card" @click="loadLogs(true)">重试</button>
      </div>
      <div v-if="logsLoading" class="list-loading">加载中...</div>
      <div v-else-if="logsError" class="logs-error">{{ logsError }}</div>
      <div v-else-if="logs.length === 0" class="logs-empty">暂无兑换记录</div>
      <div v-else>
        <div class="logs-table">
          <div v-for="log in logs" :key="log.id" class="log-row">
            <div class="log-main">
              <span class="log-name">{{ log.student_name || '已删除学生' }}</span>
              <span v-if="log.student_no" class="log-no">📛{{ log.student_no }}</span>
            </div>
            <div class="log-amount">
              <span class="log-from">-{{ log.from_amount }} 积分</span>
              <span class="log-arrow">→</span>
              <span class="log-to">+{{ log.to_amount }} {{ fmtCurrency(log.to_currency) }}</span>
            </div>
            <div class="log-date">{{ new Date(log.created_at).toLocaleString('zh-CN') }}</div>
          </div>
        </div>
        <div v-if="logsMeta.last_page > 1" class="mini-pagination">
          <button class="btn btn-sm btn-ghost-card" :disabled="logsMeta.current_page <= 1" @click="changeLogsPage(logsMeta.current_page - 1)">← 上一页</button>
          <span class="mini-pagination__info">{{ logsMeta.current_page }} / {{ logsMeta.last_page }}</span>
          <button class="btn btn-sm btn-ghost-card" :disabled="logsMeta.current_page >= logsMeta.last_page" @click="changeLogsPage(logsMeta.current_page + 1)">下一页 →</button>
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
.exchange-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 24px; }
.card-title { font-size: 16px; font-weight: 600; margin: 0; }
.card-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 16px; }
.search-input { width: 150px; padding: 6px 10px; font-size: 12px; }

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
.list-loading { text-align: center; padding: 24px; color: var(--color-text-secondary); }

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

/* ===== 兑换记录 ===== */
.logs-card { padding: 20px 24px; }
.logs-error { text-align: center; padding: 24px; color: var(--color-danger-text); font-size: 13px; }
.logs-empty { text-align: center; padding: 24px; color: var(--color-text-secondary); }
.log-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border);
  font-size: 13px;
}
.log-row:last-child { border-bottom: none; }
.log-main { display: flex; align-items: center; gap: 8px; min-width: 140px; }
.log-name { font-weight: 600; }
.log-no { font-size: 10px; color: var(--color-text-secondary); background: var(--tint-2); padding: 1px 6px; border-radius: 6px; }
.log-amount { flex: 1; display: flex; align-items: center; gap: 10px; }
.log-from { color: var(--color-danger-text); font-weight: 500; }
.log-arrow { color: var(--color-text-secondary); }
.log-to { color: var(--color-success-text); font-weight: 500; }
.log-date { font-size: 12px; color: var(--color-text-secondary); white-space: nowrap; }

/* ===== 迷你分页 ===== */
.mini-pagination { display: flex; align-items: center; justify-content: center; gap: 12px; padding: 10px 0 2px; }
.mini-pagination__info { font-size: 12px; color: var(--color-text-secondary); }
.mini-pagination .btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
