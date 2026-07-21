<script setup lang="ts">
import { ref } from 'vue'
import ShopPage from './ShopPage.vue'
import ExchangeCenterPage from './ExchangeCenterPage.vue'

const activeTab = ref<'shop' | 'exchange' | 'rates'>('shop')

// 汇率管理
import { useToastStore } from '@/stores/toast'
const toast = useToastStore()
interface Rate { id: number; name: string; from_currency: string; to_currency: string; rate: number; is_active: boolean }
const rates = ref<Rate[]>([])
const ratesLoading = ref(false)
const showAddRate = ref(false)
const newRate = ref({ name: '', from_currency: 'score', to_currency: 'science', rate: 1 })

const addRateStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleStatusMap = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
function getToggleStatus(id: number) {
  return toggleStatusMap.value[id] || 'idle'
}

const currencyLabels: Record<string, string> = { score: '⭐ 积分', science: '🔬 科学币', reading: '📚 读书币', class_point: '⚽ 体育币' }

async function loadRates() {
  ratesLoading.value = true
  try {
    const res = await fetch('/api/v1/teacher/exchange-rates', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    rates.value = data.data || []
  } catch { rates.value = [] }
  finally { ratesLoading.value = false }
}

async function addRate() {
  if (!newRate.value.name || newRate.value.rate <= 0) { toast.show('请填写完整信息', 'error', { position: 'center', duration: 2000 }); return }
  addRateStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/teacher/exchange-rates', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify(newRate.value),
    })
    if (res.ok) {
      addRateStatus.value = 'success'
      setTimeout(() => {
        addRateStatus.value = 'idle'
        showAddRate.value = false
        newRate.value = { name: '', from_currency: 'score', to_currency: 'science', rate: 1 }
        loadRates()
      }, 1500)
    } else {
      addRateStatus.value = 'error'
      setTimeout(() => { addRateStatus.value = 'idle' }, 3000)
    }
  } catch {
    addRateStatus.value = 'error'
    setTimeout(() => { addRateStatus.value = 'idle' }, 3000)
  }
}

async function toggleRate(r: Rate) {
  toggleStatusMap.value[r.id] = 'loading'
  try {
    await fetch('/api/v1/teacher/exchange-rates/' + r.id, {
      method: 'PUT',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ is_active: !r.is_active }),
    })
    r.is_active = !r.is_active
    toggleStatusMap.value[r.id] = 'success'
    setTimeout(() => { toggleStatusMap.value[r.id] = 'idle' }, 1500)
  } catch {
    toggleStatusMap.value[r.id] = 'error'
    setTimeout(() => { toggleStatusMap.value[r.id] = 'idle' }, 3000)
  }
}

function onActiveTabChange(tab: typeof activeTab.value) {
  activeTab.value = tab
  if (tab === 'rates' && rates.value.length === 0) loadRates()
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">班级经济</p>
        <h2 style="font-size:24px;font-weight:700;">🛍️ 积分商城</h2>
      </div>
    </div>

    <div class="tab-bar">
      <button :class="['tab-btn', { active: activeTab === 'shop' }]" @click="onActiveTabChange('shop')">📦 商品管理</button>
      <button :class="['tab-btn', { active: activeTab === 'rates' }]" @click="onActiveTabChange('rates')">💱 汇率设定</button>
      <button :class="['tab-btn', { active: activeTab === 'exchange' }]" @click="onActiveTabChange('exchange')">🔄 兑换管理</button>
    </div>

    <ShopPage v-if="activeTab === 'shop'" />

    <!-- 汇率设定 -->
    <div v-if="activeTab === 'rates'">
      <div class="card" style="max-width:640px;padding:24px;margin-bottom:16px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;">
          <h3 style="font-size:16px;font-weight:600;margin:0;">💱 币种汇率</h3>
          <button class="btn btn-sm btn-primary" @click="showAddRate = !showAddRate">{{ showAddRate ? '取消' : '添加汇率' }}</button>
        </div>
        <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:12px;">
          设定不同币种之间的兑换比例，影响学生在兑换中心的实际兑换结果。
        </p>
        <div v-if="showAddRate" style="display:grid;grid-template-columns:1fr 1fr 1fr 1fr;gap:12px;align-items:end;margin-bottom:16px;padding:16px;background:var(--color-bg);border-radius:8px;">
          <div class="form-group"><label>规则名称</label><input v-model="newRate.name" class="form-input" placeholder="如：积分换科学币"></div>
          <div class="form-group"><label>源币种</label><select v-model="newRate.from_currency" class="form-input"><option v-for="(label, key) in currencyLabels" :key="key" :value="key">{{ label }}</option></select></div>
          <div class="form-group"><label>目标币种</label><select v-model="newRate.to_currency" class="form-input"><option v-for="(label, key) in currencyLabels" :key="key" :value="key" :disabled="key === newRate.from_currency">{{ label }}</option></select></div>
          <div class="form-group"><label>汇率 (1:?)</label><input v-model.number="newRate.rate" type="number" step="0.01" min="0.01" class="form-input"></div>
        </div>
        <button v-if="showAddRate" class="btn btn-primary" style="margin-bottom:16px;" :disabled="addRateStatus !== 'idle'"
          :style="{ background: addRateStatus === 'loading' ? '#f59e0b' : addRateStatus === 'success' ? '#10b981' : addRateStatus === 'error' ? '#ef4444' : '#7c3aed' }"
          @click="addRate">
          <template v-if="addRateStatus === 'idle'">确认添加</template>
          <template v-else-if="addRateStatus === 'loading'">添加中...</template>
          <template v-else-if="addRateStatus === 'success'">✅ 已添加</template>
          <template v-else-if="addRateStatus === 'error'">❌ 失败</template>
        </button>
      </div>

      <div v-if="ratesLoading" style="text-align:center;padding:24px;color:var(--color-text-secondary);">加载中...</div>
      <div v-else-if="rates.length === 0" class="card" style="max-width:640px;padding:24px;text-align:center;color:var(--color-text-secondary);font-size:13px;">暂无汇率规则</div>
      <div v-else class="card" style="max-width:640px;padding:24px;">
        <div v-for="r in rates" :key="r.id" style="display:flex;align-items:center;gap:12px;padding:8px 0;border-bottom:1px solid var(--color-border);font-size:13px;">
          <span style="flex:1;font-weight:500;">{{ r.name }}</span>
          <span style="color:var(--color-text-secondary);">{{ currencyLabels[r.from_currency] || r.from_currency }}</span>
          <span style="font-weight:700;">→ 1 : {{ r.rate }}</span>
          <span>{{ currencyLabels[r.to_currency] || r.to_currency }}</span>
          <button :style="{
            background: getToggleStatus(r.id) === 'loading' ? '#f59e0b' : getToggleStatus(r.id) === 'success' ? '#10b981' : getToggleStatus(r.id) === 'error' ? '#ef4444' : 'none',
            color: getToggleStatus(r.id) !== 'idle' ? 'white' : r.is_active ? 'var(--color-accent)' : 'var(--color-text-secondary)',
            border: 'none', cursor: 'pointer', fontSize: '12px', padding: '4px 8px', borderRadius: '6px'
          }" @click="toggleRate(r)" :disabled="getToggleStatus(r.id) !== 'idle'">
            <template v-if="getToggleStatus(r.id) === 'idle'">{{ r.is_active ? '✅ 启用' : '⏸ 禁用' }}</template>
            <template v-else-if="getToggleStatus(r.id) === 'loading'">处理中...</template>
            <template v-else-if="getToggleStatus(r.id) === 'success'">✅ 成功</template>
            <template v-else-if="getToggleStatus(r.id) === 'error'">❌ 失败</template>
          </button>
        </div>
      </div>
    </div>

    <ExchangeCenterPage v-if="activeTab === 'exchange'" />
  </div>
</template>

<style scoped>
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }
.form-group { margin-bottom:0; }
.form-group label { display:block; font-size:11px; font-weight:600; color:var(--color-text-secondary); margin-bottom:2px; }
.form-input { color:var(--color-text); width:100%; padding:6px 10px; border:1px solid var(--color-border); border-radius:6px; font-size:12px; outline:none; box-sizing:border-box; background:var(--color-bg-card); }
.btn { padding:6px 14px; border-radius:8px; font-size:12px; font-weight:500; cursor:pointer; border:1px solid transparent; transition:all 0.15s; }
.btn-sm { padding:4px 12px; font-size:11px; }
.btn-primary { background:#7c3aed; color:white; border-color:#7c3aed; }
.btn-primary:hover { background:#6d28d9; }
.btn-primary:disabled { opacity:0.5; cursor:not-allowed; }
</style>
