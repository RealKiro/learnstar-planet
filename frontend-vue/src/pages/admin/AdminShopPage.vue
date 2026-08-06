<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import type { ApiResponse } from '@/types'

interface ShopItemExt {
  id: number
  name: string
  cost_score: number
  currency_type?: string
  stock: number
  category: string
  description?: string
  is_active?: boolean
  scope?: 'school' | 'class'
  class_name?: string
}

const items = ref<ShopItemExt[]>([])
const loading = ref(true)
const filterCategory = ref('')
const showModal = ref(false)
const editingId = ref<number | null>(null)

const categories = [
  { key: '', label: '全部' },
  { key: 'points', label: '⭐ 积分充值' },
  { key: 'stationery', label: '✏️ 文具用品' },
  { key: 'food', label: '🍎 零食饮料' },
  { key: 'privilege', label: '🎁 特权奖励' },
  { key: 'activity', label: '🎪 活动参与' },
]
const currencyOptions = [
  { key: 'score', label: '⭐ 班级积分' },
  { key: 'science', label: '🔬 科学币' },
  { key: 'reading', label: '📚 读书币' },
  { key: 'class_point', label: '⚽ 体育币' },
]

// ===== 积分兑换汇率（全校共享，防通胀默认 2:1） =====
interface AdminRate { id: number; name?: string; from_currency: string; to_currency: string; rate: string; is_active: boolean }
const rates = ref<AdminRate[]>([])
const ratesLoading = ref(true)
const rateSaving = ref('')
const rateErrors = reactive<Record<string, string>>({})
async function loadRates() {
  try {
    const res = await apiGet<ApiResponse<AdminRate[]>>('/api/v1/admin/exchange-rates', { skipToast: true })
    rates.value = res.data || []
  } catch { rates.value = [] }
  finally { ratesLoading.value = false }
}
function rateLabel(c: string) {
  return ({ score: '⭐ 积分', science: '🔬 科学币', reading: '📚 读书币', class_point: '⚽ 体育币' } as Record<string, string>)[c] || c
}
async function saveRate(r: AdminRate) {
  const val = parseFloat(String(r.rate))
  if (!val || val <= 0) { rateErrors[String(r.id)] = '请输入正数'; return }
  rateSaving.value = String(r.id)
  rateErrors[String(r.id)] = ''
  try {
    await apiPut(`/api/v1/admin/exchange-rates/${r.id}`, { rate: val }, { skipToast: true })
    rateSaving.value = ''
  } catch { rateErrors[String(r.id)] = '保存失败' }
}

const form = ref({ name: '', description: '', category: 'stationery', cost_score: 10, currency_type: 'score', stock: 0, is_active: true })
const itemErrors = reactive<Record<string, string>>({})
function iClr(f: string) { delete itemErrors[f] }
function iVld(field: string): boolean {
  if (field === 'name' && !form.value.name.trim()) { itemErrors.name = '请填写商品名称'; return false }
  if (field === 'cost_score' && (!form.value.cost_score || form.value.cost_score < 1)) { itemErrors.cost_score = '积分至少为 1'; return false }
  delete itemErrors[field]; return true
}

const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const deleteStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})

function getDeleteStatus(id: number) { return deleteStatus.value[id] || 'idle' }
function getToggleStatus(id: number) { return toggleStatus.value[id] || 'idle' }

const saveBtnStyle = computed(() => {
  const map: Record<string, string> = { idle: '#7c3aed', loading: '#f59e0b', success: '#10b981', error: '#ef4444' }
  return { background: map[saveStatus.value] }
})
const saveBtnText = computed(() => {
  const map: Record<string, string> = { idle: editingId.value ? '保存' : '添加', loading: '处理中...', success: editingId.value ? '已更新' : '已添加', error: '保存失败' }
  return map[saveStatus.value]
})

const filteredItems = computed(() => {
  if (!filterCategory.value) return items.value
  return items.value.filter(i => i.category === filterCategory.value)
})

const currencyLabel = (key: string) => currencyOptions.find(c => c.key === key)?.label || key

async function loadItems() {
  try {
    const res = await apiGet<ApiResponse<ShopItemExt[]>>('/api/v1/admin/shop-items')
    items.value = res.data || []
  } catch { /* handled */ }
  finally { loading.value = false }
}
onMounted(() => {
  loadRates()
  loadItems()
})

function openAdd() {
  editingId.value = null
  form.value = { name: '', description: '', category: 'stationery', cost_score: 10, currency_type: 'score', stock: 0, is_active: true }
  showModal.value = true
}

function openEdit(item: ShopItemExt) {
  editingId.value = item.id
  form.value = { name: item.name, description: item.description || '', category: item.category, cost_score: item.cost_score, currency_type: item.currency_type || 'score', stock: item.stock, is_active: item.is_active ?? true }
  showModal.value = true
}

async function handleSubmit() {
  const nameOk = iVld('name')
  const scoreOk = iVld('cost_score')
  if (!nameOk || !scoreOk) return
  saveStatus.value = 'loading'
  const payload = {
    name: form.value.name.trim(),
    description: form.value.description.trim(),
    category: form.value.category,
    cost_score: form.value.cost_score,
    currency_type: form.value.currency_type,
    stock: form.value.stock,
    is_active: form.value.is_active,
  }
  try {
    if (editingId.value) {
      await apiPut(`/api/v1/admin/shop-items/${editingId.value}`, payload, { skipToast: true })
    } else {
      const res = await apiPost<ApiResponse<ShopItemExt>>('/api/v1/admin/shop-items', payload, { skipToast: true })
      const data = (res as unknown as { data: ShopItemExt }).data
      if (data) items.value.push(data)
    }
    saveStatus.value = 'success'
    showModal.value = false
    setTimeout(() => { saveStatus.value = 'idle' }, 1500)
    const fresh = await apiGet<ApiResponse<ShopItemExt[]>>('/api/v1/admin/shop-items')
    items.value = fresh.data || []
  } catch (e: any) {
    saveStatus.value = 'error'
    itemErrors.name = e?.response?.data?.message || '保存失败'
    setTimeout(() => { saveStatus.value = 'idle' }, 3000)
  }
}

async function toggleItem(item: ShopItemExt) {
  toggleStatus.value[item.id] = 'loading'
  try {
    await apiPut(`/api/v1/admin/shop-items/${item.id}`, { is_active: !item.is_active }, { skipToast: true })
    item.is_active = !item.is_active
    toggleStatus.value[item.id] = 'success'
    setTimeout(() => { toggleStatus.value[item.id] = 'idle' }, 1500)
  } catch {
    toggleStatus.value[item.id] = 'error'
    setTimeout(() => { toggleStatus.value[item.id] = 'idle' }, 3000)
  }
}

async function handleDelete(item: ShopItemExt) {
  const ok = await openConfirm({ title: '删除商品', message: `确定删除商品「${item.name}」？`, danger: true, confirmText: '确认删除' })
  if (!ok) return
  deleteStatus.value[item.id] = 'loading'
  try {
    await apiDelete(`/api/v1/admin/shop-items/${item.id}`, { skipToast: true })
    items.value = items.value.filter(i => i.id !== item.id)
    deleteStatus.value[item.id] = 'success'
    setTimeout(() => { deleteStatus.value[item.id] = 'idle' }, 1500)
  } catch (e: any) {
    deleteStatus.value[item.id] = 'error'
    setTimeout(() => { deleteStatus.value[item.id] = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">全校共享</p>
        <h2 style="font-size:24px;font-weight:700;">积分商品</h2>
      </div>
      <button class="btn btn-sm btn-primary" @click="openAdd">+ 添加商品</button>
    </div>
    <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:16px;">💡 全校教师共享这些商品，学生可在各班级兑换。</p>

    <!-- 积分兑换汇率（全校共享） -->
    <div class="card" style="padding:16px;margin-bottom:16px;">
      <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:12px;flex-wrap:wrap;gap:8px;">
        <h3 style="font-size:15px;font-weight:600;margin:0;">💱 积分兑换汇率</h3>
        <span style="font-size:11px;color:var(--color-text-secondary);">全校共享 · 默认 2:1（2 积分 = 1 币，防通胀）</span>
      </div>
      <div v-if="ratesLoading" style="font-size:12px;color:var(--color-text-secondary);">加载中...</div>
      <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:12px;">
        <div v-for="r in rates" :key="r.id" style="display:flex;align-items:center;gap:10px;padding:10px 12px;background:var(--color-bg);border-radius:8px;">
          <span style="font-size:13px;flex:1;">{{ rateLabel(r.from_currency) }} → {{ rateLabel(r.to_currency) }}</span>
          <input :value="r.rate" type="number" step="0.1" min="0.1" style="width:64px;padding:4px 8px;border-radius:6px;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text);font-size:13px;" @change="r.rate = ($event.target as HTMLInputElement).value">
          <button class="btn btn-sm btn-primary" style="white-space:nowrap;" :disabled="rateSaving === String(r.id)" @click="saveRate(r)">
            {{ rateSaving === String(r.id) ? '保存中' : '保存' }}
          </button>
        </div>
      </div>
      <div v-if="Object.keys(rateErrors).length" style="margin-top:8px;color:#f87171;font-size:12px;">{{ Object.values(rateErrors)[0] }}</div>
      <p style="font-size:11px;color:var(--color-text-secondary);margin:10px 0 0;">💡 示例：2 积分 = 1 科学币/体育币/读书币；小商品约 100 积分（一周可攒），大商品约 200 积分（两周可攒）。</p>
    </div>

    <!-- 分类筛选 -->
    <div style="display:flex;gap:8px;margin-bottom:16px;flex-wrap:wrap;">
      <button v-for="c in categories" :key="c.key" class="btn btn-sm"
        :style="{ background: filterCategory === c.key ? '#7c3aed' : 'var(--color-bg-card)', color: filterCategory === c.key ? '#fff' : 'var(--color-text)', border: '1px solid var(--color-border)' }"
        @click="filterCategory = c.key">{{ c.label }}</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="filteredItems.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">🛍️</div>
      <p style="margin-bottom:16px;">暂无全校商品，点击「添加商品」创建</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:16px;">
      <div v-for="item in filteredItems" :key="item.id" class="card"
        style="padding:16px;display:flex;flex-direction:column;gap:8px;border-color:rgba(255,255,255,0.06);"
        :style="{ opacity: item.is_active === false ? 0.6 : 1 }">
        <div style="display:flex;align-items:center;justify-content:space-between;">
          <span style="font-size:24px;">{{ { points: '⭐', stationery: '✏️', food: '🍎', privilege: '🎁', activity: '🎪' }[item.category] || '📦' }}</span>
          <span v-if="item.is_active === false" style="font-size:10px;color:#f59e0b;border:1px solid rgba(245,158,11,0.3);padding:1px 6px;border-radius:4px;">停用</span>
        </div>
        <div style="display:flex;align-items:center;gap:6px;font-weight:600;font-size:15px;">
          {{ item.name }}
          <span v-if="item.scope === 'class'" style="font-size:10px;color:#38bdf8;background:rgba(56,189,248,0.12);border:1px solid rgba(56,189,248,0.3);padding:1px 6px;border-radius:4px;font-weight:500;flex-shrink:0;">{{ item.class_name || '班级' }}</span>
        </div>
        <div v-if="item.description" style="font-size:12px;color:var(--color-text-secondary);line-height:1.5;">{{ item.description }}</div>
        <div style="display:flex;gap:6px;flex-wrap:wrap;font-size:11px;">
          <span style="background:var(--color-bg);padding:2px 8px;border-radius:4px;">{{ currencyLabel(item.currency_type || 'score') }}</span>
          <span style="background:var(--color-bg);padding:2px 8px;border-radius:4px;">库存 {{ item.stock === 0 ? '∞' : item.stock }}</span>
        </div>
        <div style="display:flex;align-items:center;justify-content:space-between;margin-top:4px;">
          <span style="font-weight:700;color:var(--color-primary);">⭐ {{ item.cost_score }}</span>
          <div style="display:flex;gap:4px;flex-wrap:wrap;justify-content:flex-end;">
            <button class="btn btn-sm btn-ghost" @click="toggleItem(item)" :disabled="getToggleStatus(item.id) === 'loading'" :style="{ color: getToggleStatus(item.id) === 'success' ? '#10b981' : getToggleStatus(item.id) === 'error' ? '#ef4444' : getToggleStatus(item.id) === 'loading' ? '#f59e0b' : item.is_active ? 'var(--color-text-secondary)' : '#f59e0b' }">
              <template v-if="getToggleStatus(item.id) === 'loading'">切换中</template>
              <template v-else-if="getToggleStatus(item.id) === 'success'">已切换</template>
              <template v-else-if="getToggleStatus(item.id) === 'error'">失败</template>
              <template v-else>{{ item.is_active ? '停用' : '启用' }}</template>
            </button>
            <button class="btn btn-sm btn-ghost" @click="openEdit(item)">编辑</button>
            <button class="btn btn-sm btn-ghost" @click="handleDelete(item)" :disabled="getDeleteStatus(item.id) === 'loading'" :style="{ color: getDeleteStatus(item.id) === 'loading' ? '#f59e0b' : getDeleteStatus(item.id) === 'success' ? '#10b981' : getDeleteStatus(item.id) === 'error' ? '#ef4444' : 'var(--color-danger)' }">{{ getDeleteStatus(item.id) === 'loading' ? '删除中' : getDeleteStatus(item.id) === 'success' ? '已删除' : getDeleteStatus(item.id) === 'error' ? '失败' : '删除' }}</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑弹窗 -->
    <div v-if="showModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100;" @click.self="showModal = false">
      <div class="card" style="width:480px;max-width:90vw;padding:24px;max-height:85vh;overflow-y:auto;">
        <h3 style="font-size:18px;font-weight:600;margin-bottom:20px;">{{ editingId ? '编辑商品' : '添加商品' }}</h3>
        <div class="form-group">
          <label>商品名称</label>
          <input v-model="form.name" class="form-input" placeholder="如：铅笔" :style="{ borderColor: itemErrors.name ? '#f87171' : '' }" @blur="iVld('name')" @input="iClr('name')">
          <div v-if="itemErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ itemErrors.name }}</div>
        </div>
        <div class="form-group">
          <label>所需积分</label>
          <input v-model.number="form.cost_score" type="number" min="1" class="form-input" :style="{ borderColor: itemErrors.cost_score ? '#f87171' : '' }" @blur="iVld('cost_score')" @input="iClr('cost_score')">
          <div v-if="itemErrors.cost_score" style="color:#f87171;font-size:11px;margin-top:2px;">{{ itemErrors.cost_score }}</div>
        </div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;">
          <div class="form-group">
            <label>类别</label>
            <select v-model="form.category" class="form-input">
              <option v-for="c in categories.filter(x => x.key)" :key="c.key" :value="c.key">{{ c.label }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>结算币种</label>
            <select v-model="form.currency_type" class="form-input">
              <option v-for="c in currencyOptions" :key="c.key" :value="c.key">{{ c.label }}</option>
            </select>
          </div>
        </div>
        <div class="form-group">
          <label>库存（0 = 不限制）</label>
          <input v-model.number="form.stock" type="number" min="0" class="form-input">
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="form.description" class="form-input" style="min-height:60px;resize:vertical;" placeholder="可选"></textarea>
        </div>
        <label style="display:flex;align-items:center;gap:8px;font-size:13px;margin-bottom:16px;cursor:pointer;">
          <input type="checkbox" v-model="form.is_active" style="accent-color:#7c3aed;"> 启用
        </label>
        <div style="display:flex;justify-content:flex-end;gap:8px;">
          <button class="btn btn-ghost" @click="showModal = false">取消</button>
          <button class="btn btn-primary" style="width:auto;" @click="handleSubmit" :disabled="saveStatus === 'loading'" :style="saveBtnStyle">{{ saveBtnText }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.form-group { margin-bottom: 12px; }
.form-group label { display: block; font-size: 12px; font-weight: 600; color: var(--color-text); margin-bottom: 4px; }
.form-input { color: var(--color-text); width: 100%; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 13px; outline: none; box-sizing: border-box; font-family: inherit; }
.form-input:focus { border-color: #7c3aed; box-shadow: 0 0 0 3px rgba(124,58,237,0.08); }
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; display: inline-flex; align-items: center; gap: 4px; font-family: inherit; }
.btn-primary { background: #7c3aed; color: white; }
.btn-primary:hover { background: #6d28d9; }
.btn-ghost { background: transparent; color: var(--color-text-secondary); border: 1px solid var(--color-border); }
.btn-ghost:hover { background: rgba(255,255,255,0.06); }
.btn-sm { padding: 5px 12px; font-size: 12px; }
</style>
