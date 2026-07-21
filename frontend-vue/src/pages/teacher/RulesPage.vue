<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ScoreRule } from '@/types'

const toast = useToastStore()

const rules = ref<ScoreRule[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = ref({ name: '', points: 1, category: 'classroom', is_penalty: false })
const ruleErrors = reactive<Record<string, string>>({})
function rClr(f: string) { delete ruleErrors[f] }
function rVld(field: string): boolean {
  if (field === 'name' && !form.value.name.trim()) { ruleErrors.name = '请填写规则名称，如"举手发言"'; return false }
  if (field === 'points' && (!form.value.points || form.value.points < 1)) { ruleErrors.points = '分值至少为 1'; return false }
  delete ruleErrors[field]; return true
}

const positiveRules = computed(() => rules.value.filter(r => !r.is_penalty))
const negativeRules = computed(() => rules.value.filter(r => r.is_penalty))

const categoryLabels: Record<string, string> = {
  classroom: '📖 课堂表现', homework: '📝 作业管理', behavior: '🌟 行为习惯',
  literacy: '📊 综合素养', daily: '📅 日常表现', academic: '📚 学业',
}

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules')
    rules.value = res.data || []
  } catch { /* handled */ }
  finally { loading.value = false }
})

function openAdd() {
  editingId.value = null
  form.value = { name: '', points: 1, category: '', is_penalty: false }
  showModal.value = true
}

function openEdit(rule: ScoreRule) {
  editingId.value = rule.id
  form.value = { name: rule.name, points: Math.abs(rule.points), category: rule.category, is_penalty: rule.is_penalty }
  showModal.value = true
}

async function handleSubmit() {
  if (!rVld('name') | !rVld('points')) return
  const payload = {
    name: form.value.name.trim(),
    points: form.value.is_penalty ? -Math.abs(form.value.points) : Math.abs(form.value.points),
    category: form.value.category.trim() || '默认',
    is_penalty: form.value.is_penalty,
  }
  try {
    if (editingId.value) {
      await apiPut(`/api/v1/teacher/scores/rules/${editingId.value}`, payload)
      toast.show('规则已更新', 'success')
    } else {
      const res = await apiPost<ApiResponse<ScoreRule>>('/api/v1/teacher/scores/rules', payload)
      const data = (res as unknown as { data: ScoreRule }).data
      if (data) rules.value.push(data)
      toast.show('规则已添加', 'success')
    }
    showModal.value = false
    // 重新拉取保证数据一致
    const fresh = await apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules')
    rules.value = fresh.data || []
  } catch { /* handled */ }
}

async function handleDelete(rule: ScoreRule) {
  if (!confirm(`确认删除规则「${rule.name}」？`)) return
  try {
    await apiDelete(`/api/v1/teacher/scores/rules/${rule.id}`)
    rules.value = rules.value.filter(r => r.id !== rule.id)
    toast.show('规则已删除', 'success')
  } catch { /* handled */ }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">积分规则</h2>
      <button class="btn btn-sm btn-primary" @click="openAdd">添加规则</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="rules.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📋</div>
      <p style="margin-bottom:16px;">暂无积分规则，点击「添加规则」创建</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:1fr 1fr;gap:24px;">
      <!-- 加分规则 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">加分规则 <span style="color:var(--color-accent);font-size:13px;">({{ positiveRules.length }})</span></h3>
        <div v-if="positiveRules.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);font-size:13px;">暂无加分规则</div>
        <div v-else style="display:flex;flex-direction:column;gap:8px;">
          <div v-for="rule in positiveRules" :key="rule.id" class="card"
            style="padding:12px 16px;display:flex;align-items:center;justify-content:space-between;border-color:rgba(16,185,129,0.3);">
            <div>
              <span style="font-weight:700;color:var(--color-accent);">+{{ Math.abs(rule.points) }}</span>
              <span style="margin-left:8px;font-weight:500;">{{ rule.name }}</span>
              <span style="margin-left:8px;font-size:11px;color:var(--color-text-secondary);background:var(--color-bg);padding:2px 8px;border-radius:4px;">{{ categoryLabels[rule.category] || rule.category }}</span>
            </div>
            <div style="display:flex;gap:4px;">
              <button class="btn btn-sm btn-ghost" @click="openEdit(rule)">编辑</button>
              <button class="btn btn-sm btn-ghost" style="color:var(--color-danger);" @click="handleDelete(rule)">删除</button>
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
            <div>
              <span style="font-weight:700;color:var(--color-danger);">{{ rule.points }}</span>
              <span style="margin-left:8px;font-weight:500;">{{ rule.name }}</span>
              <span style="margin-left:8px;font-size:11px;color:var(--color-text-secondary);background:var(--color-bg);padding:2px 8px;border-radius:4px;">{{ categoryLabels[rule.category] || rule.category }}</span>
            </div>
            <div style="display:flex;gap:4px;">
              <button class="btn btn-sm btn-ghost" @click="openEdit(rule)">编辑</button>
              <button class="btn btn-sm btn-ghost" style="color:var(--color-danger);" @click="handleDelete(rule)">删除</button>
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
          <input v-model.number="form.points" type="number" min="1" class="form-input" placeholder="如：5" :style="{ borderColor: ruleErrors.points ? '#f87171' : '' }" @blur="rVld('points')" @input="rClr('points')">
          <div v-if="ruleErrors.points" style="color:#f87171;font-size:11px;margin-top:2px;">{{ ruleErrors.points }}</div>
        </div>
        <div class="form-group">
          <label>分类</label>
          <select v-model="form.category" class="form-input">
            <option v-for="(label, key) in categoryLabels" :key="key" :value="key">{{ label }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>类型</label>
          <select v-model="form.is_penalty" class="form-select">
            <option :value="false">加分规则</option>
            <option :value="true">扣分规则</option>
          </select>
        </div>
        <div style="display:flex;justify-content:flex-end;gap:8px;margin-top:8px;">
          <button class="btn btn-ghost" @click="showModal = false">取消</button>
          <button class="btn btn-primary" style="width:auto;" @click="handleSubmit">{{ editingId ? '保存' : '添加' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
