<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface Rate {
  id: number
  name: string
  from_currency: string
  to_currency: string
  rate: number
  is_active: boolean
}

const rates = ref<Rate[]>([])
const loading = ref(true)
const showAddForm = ref(false)
const newRate = ref({ name: '', from_currency: 'score', to_currency: 'science_coin', rate: 1 })
const addRateStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleRateStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleRateId = ref<number | null>(null)

const currencyOptions = [
  { key: 'score', label: '⭐ 积分' },
  { key: 'science_coin', label: '🔬 科学币' },
  { key: 'reading_coin', label: '📚 读书币' },
]

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Rate[]>>('/api/v1/admin/exchange-rates')
    rates.value = res.data || []
  } catch {
    rates.value = demoRates()
  }
  finally { loading.value = false }
})

function demoRates(): Rate[] {
  return [
    { id: 1, name: '积分换科学币', from_currency: 'score', to_currency: 'science_coin', rate: 1, is_active: true },
    { id: 2, name: '积分换读书币', from_currency: 'score', to_currency: 'reading_coin', rate: 1, is_active: true },
    { id: 3, name: '科学币换读书币', from_currency: 'science_coin', to_currency: 'reading_coin', rate: 0.5, is_active: true },
  ]
}

async function addRate() {
  if (!newRate.value.name || newRate.value.rate <= 0) {
    toast.show('请填写完整信息', 'error', { position: 'top-right' })
    return
  }
  addRateStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/exchange-rates', newRate.value)
    addRateStatus.value = 'success'
    showAddForm.value = false
    newRate.value = { name: '', from_currency: 'score', to_currency: 'science_coin', rate: 1 }
    const res = await apiGet<ApiResponse<Rate[]>>('/api/v1/admin/exchange-rates')
    rates.value = res.data || []
    setTimeout(() => { addRateStatus.value = 'idle' }, 1500)
  } catch {
    addRateStatus.value = 'error'
    setTimeout(() => { addRateStatus.value = 'idle' }, 3000)
  }
}

async function toggleRate(rate: Rate) {
  toggleRateStatus.value = 'loading'
  toggleRateId.value = rate.id
  try {
    await apiPut(`/api/v1/admin/exchange-rates/${rate.id}`, { is_active: !rate.is_active })
    rate.is_active = !rate.is_active
    toggleRateStatus.value = 'success'
    setTimeout(() => { toggleRateStatus.value = 'idle'; toggleRateId.value = null }, 1500)
  } catch {
    toggleRateStatus.value = 'error'
    setTimeout(() => { toggleRateStatus.value = 'idle'; toggleRateId.value = null }, 3000)
  }
}

const getCurrencyLabel = (key: string) => currencyOptions.find(c => c.key === key)?.label || key
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">系统配置</p>
        <h2 style="font-size:24px;font-weight:700;">💱 汇率配置</h2>
      </div>
      <button class="btn btn-sm btn-primary" @click="showAddForm = !showAddForm">{{ showAddForm ? '取消' : '添加汇率' }}</button>
    </div>

    <div v-if="showAddForm" class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">新汇率规则</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr 1fr 1fr;gap:16px;align-items:end;">
        <div class="form-group">
          <label>规则名称</label>
          <input v-model="newRate.name" class="form-input" placeholder="如：积分换科学币">
        </div>
        <div class="form-group">
          <label>源币种</label>
          <select v-model="newRate.from_currency" class="form-select">
            <option v-for="c in currencyOptions" :key="c.key" :value="c.key">{{ c.label }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>目标币种</label>
          <select v-model="newRate.to_currency" class="form-select">
            <option v-for="c in currencyOptions" :key="c.key" :value="c.key" :disabled="c.key === newRate.from_currency">{{ c.label }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>汇率 (1 源 = ? 目标)</label>
          <input v-model.number="newRate.rate" type="number" step="0.01" min="0.01" class="form-input">
        </div>
      </div>
      <button class="btn" :style="{ background: addRateStatus === 'loading' ? '#f59e0b' : addRateStatus === 'success' ? '#10b981' : addRateStatus === 'error' ? '#ef4444' : '#7c3aed', color: '#fff', border: '1px solid transparent', marginTop: '16px', width: 'auto' }" :disabled="addRateStatus !== 'idle'" @click="addRate">
        <template v-if="addRateStatus === 'loading'">添加中...</template>
        <template v-else-if="addRateStatus === 'success'">已添加 ✓</template>
        <template v-else-if="addRateStatus === 'error'">添加失败 ✗</template>
        <template v-else>添加</template>
      </button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else class="data-table">
      <table>
        <thead><tr><th>规则名称</th><th>源币种</th><th>目标币种</th><th>汇率</th><th>状态</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-if="rates.length === 0">
            <td colspan="6" style="text-align:center;padding:48px;color:var(--color-text-secondary);">暂无汇率规则</td>
          </tr>
          <tr v-for="r in rates" :key="r.id">
            <td style="font-weight:600;">{{ r.name }}</td>
            <td>{{ getCurrencyLabel(r.from_currency) }}</td>
            <td>{{ getCurrencyLabel(r.to_currency) }}</td>
            <td style="font-weight:700;">1 : {{ r.rate }}</td>
            <td :style="{ color: r.is_active ? 'var(--color-accent)' : 'var(--color-text-secondary)' }">
              {{ r.is_active ? '✅ 启用' : '⏸ 禁用' }}
            </td>
            <td>
              <button class="btn btn-sm" :style="toggleRateId === r.id ? { background: toggleRateStatus === 'loading' ? '#f59e0b' : toggleRateStatus === 'success' ? '#10b981' : '#ef4444', color: '#fff', border: '1px solid transparent' } : { color: r.is_active ? 'var(--color-danger)' : 'var(--color-accent)' }"
                :disabled="toggleRateId === r.id" @click="toggleRate(r)">
                <template v-if="toggleRateId === r.id && toggleRateStatus === 'loading'">操作中...</template>
                <template v-else-if="toggleRateId === r.id && toggleRateStatus === 'success'">已完成 ✓</template>
                <template v-else-if="toggleRateId === r.id && toggleRateStatus === 'error'">失败 ✗</template>
                <template v-else>{{ r.is_active ? '禁用' : '启用' }}</template>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
