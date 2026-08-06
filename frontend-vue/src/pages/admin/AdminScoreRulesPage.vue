<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import type { ApiResponse, ScoreRule } from '@/types'

const rules = ref<ScoreRule[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = ref({ name: '', amount: 1, category: 'classroom', is_positive: true, is_active: true })
const ruleErrors = reactive<Record<string, string>>({})
function rClr(f: string) { delete ruleErrors[f] }
function rVld(field: string): boolean {
  if (field === 'name' && !form.value.name.trim()) { ruleErrors.name = '请填写规则名称'; return false }
  if (field === 'amount' && (!form.value.amount || form.value.amount < 1)) { ruleErrors.amount = '分值至少为 1'; return false }
  delete ruleErrors[field]; return true
}

const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const deleteStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})

function getDeleteStatus(ruleId: number) { return deleteStatus.value[ruleId] || 'idle' }
function getToggleStatus(ruleId: number) { return toggleStatus.value[ruleId] || 'idle' }

const saveBtnStyle = computed(() => {
  const map: Record<string, string> = { idle: '#7c3aed', loading: '#f59e0b', success: '#10b981', error: '#ef4444' }
  return { background: map[saveStatus.value] }
})
const saveBtnText = computed(() => {
  const map: Record<string, string> = { idle: editingId.value ? '保存' : '添加', loading: '处理中...', success: editingId.value ? '已更新' : '已添加', error: '保存失败' }
  return map[saveStatus.value]
})

const positiveRules = computed(() => rules.value.filter(r => r.is_positive))
const negativeRules = computed(() => rules.value.filter(r => !r.is_positive))

const categoryLabels: Record<string, string> = {
  classroom: '📖 课堂表现', homework: '📝 作业管理', behavior: '🌟 行为习惯',
  literacy: '📊 综合素养', daily: '📅 日常表现', academic: '📚 学业', custom: '✨ 自定义',
}

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<ScoreRule[]>>('/api/v1/admin/score-rules')
    rules.value = res.data || []
  } catch { /* handled */ }
  finally { loading.value = false }
})

function openAdd() {
  editingId.value = null
  form.value = { name: '', amount: 1, category: 'classroom', is_positive: true, is_active: true }
  showModal.value = true
}

function openEdit(rule: ScoreRule) {
  editingId.value = rule.id
  form.value = { name: rule.name, amount: Math.abs(rule.amount), category: rule.category, is_positive: rule.is_positive, is_active: rule.is_active ?? true }
  showModal.value = true
}

async function handleSubmit() {
  const nameOk = rVld('name')
  const pointsOk = rVld('amount')
  if (!nameOk || !pointsOk) return
  saveStatus.value = 'loading'
  const payload = {
    name: form.value.name.trim(),
    amount: form.value.is_positive ? Math.abs(form.value.amount) : -Math.abs(form.value.amount),
    category: form.value.category.trim() || 'custom',
    is_positive: form.value.is_positive,
    is_active: form.value.is_active,
  }
  try {
    if (editingId.value) {
      await apiPut(`/api/v1/admin/score-rules/${editingId.value}`, payload, { skipToast: true })
    } else {
      const res = await apiPost<ApiResponse<ScoreRule>>('/api/v1/admin/score-rules', payload, { skipToast: true })
      const data = (res as unknown as { data: ScoreRule }).data
      if (data) rules.value.push(data)
    }
    saveStatus.value = 'success'
    showModal.value = false
    setTimeout(() => { saveStatus.value = 'idle' }, 1500)
    const fresh = await apiGet<ApiResponse<ScoreRule[]>>('/api/v1/admin/score-rules')
    rules.value = fresh.data || []
  } catch (e: any) {
    saveStatus.value = 'error'
    ruleErrors.name = e?.response?.data?.message || '保存失败'
    setTimeout(() => { saveStatus.value = 'idle' }, 3000)
  }
}

async function toggleRule(rule: ScoreRule) {
  toggleStatus.value[rule.id] = 'loading'
  try {
    await apiPut(`/api/v1/admin/score-rules/${rule.id}`, { is_active: !rule.is_active }, { skipToast: true })
    rule.is_active = !rule.is_active
    toggleStatus.value[rule.id] = 'success'
    setTimeout(() => { toggleStatus.value[rule.id] = 'idle' }, 1500)
  } catch {
    toggleStatus.value[rule.id] = 'error'
    setTimeout(() => { toggleStatus.value[rule.id] = 'idle' }, 3000)
  }
}

async function handleDelete(rule: ScoreRule) {
  const ok = await openConfirm({ title: '删除积分规则', message: `确认删除规则「${rule.name}」？`, danger: true, confirmText: '确认删除' })
  if (!ok) return
  deleteStatus.value[rule.id] = 'loading'
  try {
    await apiDelete(`/api/v1/admin/score-rules/${rule.id}`, { skipToast: true })
    rules.value = rules.value.filter(r => r.id !== rule.id)
    deleteStatus.value[rule.id] = 'success'
    setTimeout(() => { deleteStatus.value[rule.id] = 'idle' }, 1500)
  } catch (e: any) {
    deleteStatus.value[rule.id] = 'error'
    setTimeout(() => { deleteStatus.value[rule.id] = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">全校共享</p>
        <h2 style="font-size:24px;font-weight:700;">积分规则</h2>
      </div>
      <button class="btn btn-sm btn-primary" @click="openAdd">+ 添加规则</button>
    </div>
    <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:16px;">💡 全校教师共享这些规则，新增/修改后教师端即时可见。</p>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="rules.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📋</div>
      <p style="margin-bottom:16px;">暂无全校积分规则，点击「添加规则」创建</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:1fr 1fr;gap:24px;">
      <!-- 加分规则 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">加分规则 <span style="color:var(--color-accent);font-size:13px;">({{ positiveRules.length }})</span></h3>
        <div v-if="positiveRules.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);font-size:13px;">暂无加分规则</div>
        <div v-else style="display:flex;flex-direction:column;gap:8px;">
          <div v-for="rule in positiveRules" :key="rule.id" class="card"
            style="padding:12px 16px;display:flex;align-items:center;justify-content:space-between;border-color:rgba(16,185,129,0.3);">
            <div style="display:flex;align-items:center;gap:8px;flex:1;min-width:0;">
              <span style="font-weight:700;color:var(--color-accent);">+{{ Math.abs(rule.amount) }}</span>
              <span style="font-weight:500;">{{ rule.name }}</span>
              <span v-if="(rule as any).scope === 'class'" style="font-size:10px;color:#38bdf8;background:rgba(56,189,248,0.12);border:1px solid rgba(56,189,248,0.3);padding:1px 6px;border-radius:4px;white-space:nowrap;">{{ (rule as any).class_name || '班级' }}</span>
              <span style="font-size:11px;color:var(--color-text-secondary);background:var(--color-bg);padding:2px 8px;border-radius:4px;white-space:nowrap;">{{ categoryLabels[rule.category] || rule.category }}</span>
              <span v-if="!rule.is_active" style="font-size:10px;color:#f59e0b;border:1px solid rgba(245,158,11,0.3);padding:1px 6px;border-radius:4px;">停用</span>
            </div>
            <div style="display:flex;gap:4px;flex-shrink:0;">
              <button class="btn btn-sm btn-ghost" @click="toggleRule(rule)" :disabled="getToggleStatus(rule.id) === 'loading'" :style="{ color: getToggleStatus(rule.id) === 'success' ? '#10b981' : getToggleStatus(rule.id) === 'error' ? '#ef4444' : getToggleStatus(rule.id) === 'loading' ? '#f59e0b' : rule.is_active ? 'var(--color-text-secondary)' : '#f59e0b' }">
                <template v-if="getToggleStatus(rule.id) === 'loading'">切换中</template>
                <template v-else-if="getToggleStatus(rule.id) === 'success'">已切换</template>
                <template v-else-if="getToggleStatus(rule.id) === 'error'">失败</template>
                <template v-else>{{ rule.is_active ? '停用' : '启用' }}</template>
              </button>
              <button class="btn btn-sm btn-ghost" @click="openEdit(rule)">编辑</button>
              <button class="btn btn-sm btn-ghost" @click="handleDelete(rule)" :disabled="getDeleteStatus(rule.id) === 'loading'" :style="{ color: getDeleteStatus(rule.id) === 'loading' ? '#f59e0b' : getDeleteStatus(rule.id) === 'success' ? '#10b981' : getDeleteStatus(rule.id) === 'error' ? '#ef4444' : 'var(--color-danger)' }">{{ getDeleteStatus(rule.id) === 'loading' ? '删除中' : getDeleteStatus(rule.id) === 'success' ? '已删除' : getDeleteStatus(rule.id) === 'error' ? '失败' : '删除' }}</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 扣分规则 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">扣分规则 <span style="color:var(--color-danger);font-size:13px;">({{ negativeRules.length }})</span></h3>
        <div v-if="negativeRules.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);font-size:13px;">暂无扣分规则</div>
        <div v-else style="display:flex;flex-direction:column;gap:8px;">
          <div v-for="rule in negativeRules" :key="rule.id" class="card"
            style="padding:12px 16px;display:flex;align-items:center;justify-content:space-between;border-color:rgba(239,68,68,0.3);">
            <div style="display:flex;align-items:center;gap:8px;flex:1;min-width:0;">
              <span style="font-weight:700;color:var(--color-danger);">{{ rule.amount }}</span>
              <span style="font-weight:500;">{{ rule.name }}</span>
              <span v-if="(rule as any).scope === 'class'" style="font-size:10px;color:#38bdf8;background:rgba(56,189,248,0.12);border:1px solid rgba(56,189,248,0.3);padding:1px 6px;border-radius:4px;white-space:nowrap;">{{ (rule as any).class_name || '班级' }}</span>
              <span style="font-size:11px;color:var(--color-text-secondary);background:var(--color-bg);padding:2px 8px;border-radius:4px;white-space:nowrap;">{{ categoryLabels[rule.category] || rule.category }}</span>
              <span v-if="!rule.is_active" style="font-size:10px;color:#f59e0b;border:1px solid rgba(245,158,11,0.3);padding:1px 6px;border-radius:4px;">停用</span>
            </div>
            <div style="display:flex;gap:4px;flex-shrink:0;">
              <button class="btn btn-sm btn-ghost" @click="toggleRule(rule)" :disabled="getToggleStatus(rule.id) === 'loading'" :style="{ color: getToggleStatus(rule.id) === 'success' ? '#10b981' : getToggleStatus(rule.id) === 'error' ? '#ef4444' : getToggleStatus(rule.id) === 'loading' ? '#f59e0b' : rule.is_active ? 'var(--color-text-secondary)' : '#f59e0b' }">
                <template v-if="getToggleStatus(rule.id) === 'loading'">切换中</template>
                <template v-else-if="getToggleStatus(rule.id) === 'success'">已切换</template>
                <template v-else-if="getToggleStatus(rule.id) === 'error'">失败</template>
                <template v-else>{{ rule.is_active ? '停用' : '启用' }}</template>
              </button>
              <button class="btn btn-sm btn-ghost" @click="openEdit(rule)">编辑</button>
              <button class="btn btn-sm btn-ghost" @click="handleDelete(rule)" :disabled="getDeleteStatus(rule.id) === 'loading'" :style="{ color: getDeleteStatus(rule.id) === 'loading' ? '#f59e0b' : getDeleteStatus(rule.id) === 'success' ? '#10b981' : getDeleteStatus(rule.id) === 'error' ? '#ef4444' : 'var(--color-danger)' }">{{ getDeleteStatus(rule.id) === 'loading' ? '删除中' : getDeleteStatus(rule.id) === 'success' ? '已删除' : getDeleteStatus(rule.id) === 'error' ? '失败' : '删除' }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑弹窗 -->
    <div v-if="showModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100;" @click.self="showModal = false">
      <div class="card" style="width:420px;max-width:90vw;padding:24px;">
        <h3 style="font-size:18px;font-weight:600;margin-bottom:20px;">{{ editingId ? '编辑规则' : '添加规则' }}</h3>
        <div class="form-group">
          <label>规则名称</label>
          <input v-model="form.name" class="form-input" placeholder="如：举手发言" :style="{ borderColor: ruleErrors.name ? '#f87171' : '' }" @blur="rVld('name')" @input="rClr('name')" @keydown.enter="handleSubmit">
          <div v-if="ruleErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ ruleErrors.name }}</div>
        </div>
        <div class="form-group">
          <label>分值</label>
          <input v-model.number="form.amount" type="number" min="1" class="form-input" placeholder="如：5" :style="{ borderColor: ruleErrors.amount ? '#f87171' : '' }" @blur="rVld('amount')" @input="rClr('amount')">
          <div v-if="ruleErrors.amount" style="color:#f87171;font-size:11px;margin-top:2px;">{{ ruleErrors.amount }}</div>
        </div>
        <div class="form-group">
          <label>分类</label>
          <select v-model="form.category" class="form-input">
            <option v-for="(label, key) in categoryLabels" :key="key" :value="key">{{ label }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>类型</label>
          <select v-model="form.is_positive" class="form-select">
            <option :value="true">加分规则</option>
            <option :value="false">扣分规则</option>
          </select>
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
.form-input { color: var(--color-text); width: 100%; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 13px; outline: none; box-sizing: border-box; }
.form-input:focus { border-color: #7c3aed; box-shadow: 0 0 0 3px rgba(124,58,237,0.08); }
.form-select { color: var(--color-text); width: 100%; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 13px; outline: none; background: var(--color-bg-card); }
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; display: inline-flex; align-items: center; gap: 4px; font-family: inherit; }
.btn-primary { background: #7c3aed; color: white; }
.btn-primary:hover { background: #6d28d9; }
.btn-ghost { background: transparent; color: var(--color-text-secondary); border: 1px solid var(--color-border); }
.btn-ghost:hover { background: rgba(255,255,255,0.06); }
.btn-sm { padding: 5px 12px; font-size: 12px; }
</style>
