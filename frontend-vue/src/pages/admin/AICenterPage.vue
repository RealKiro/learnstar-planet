<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'

interface ProviderConfig {
  id: string; label: string; api_key: string; api_base: string; model: string
  is_active: boolean; billing_enabled?: boolean
  input_price_per_m?: number; output_price_per_m?: number; currency?: string
  tokens_used?: number; total_calls?: number; estimated_cost?: number
  // 本地 UI 状态（_ 前缀，不参与后端校验语义）
  _expanded?: boolean; _official_models?: string[]; _fetching?: boolean; _fetch_msg?: string
  _testing?: boolean; _test_ok?: boolean; _test_msg?: string; _show_key?: boolean
}
interface AiSettings { enabled: boolean; max_tokens: number; tokens_used: number; tokens_limit: number; providers: ProviderConfig[] }
interface DailyUsage { date: string; tokens: number; count: number }
interface ConversationLog { id: number; student_name: string; provider?: string; question: string; answer: string; tokens_used: number; cost?: number; currency?: string; created_at: string }
interface ProviderUsage { tokens: number; total_calls: number; estimated_cost: number; cost_per_token?: number; currency: string }
interface AiUsage { enabled: boolean; tokens_used: number; tokens_limit: number; estimated_cost?: number; total_conversations: number; daily_usage: DailyUsage[]; by_provider?: Record<string, ProviderUsage>; recent_logs: ConversationLog[] }

const loading = ref(true)
const loadError = ref('')
const settings = ref<AiSettings | null>(null)
const usage = ref<AiUsage | null>(null)
const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const activeTab = ref<'providers' | 'mcp' | 'usage' | 'logs'>('providers')
const logSearch = ref('')

// ===== 供应商元数据（品牌色 / 预设模型 / 参考定价 / 官网） =====
interface PricingInfo { input: string; output: string; unit: string; url: string }
interface ProviderMeta { id: string; label: string; group: string; color: string; models: string[]; pricing: PricingInfo; site: string }

const providerMeta: ProviderMeta[] = [
  { id: 'openai', label: 'OpenAI', group: '国际', color: '#10a37f', site: 'https://platform.openai.com/api-keys',
    models: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-3.5-turbo'],
    pricing: { input: '$2.50', output: '$10.00', unit: '/M tokens', url: 'https://openai.com/pricing' }},
  { id: 'claude', label: 'Anthropic Claude', group: '国际', color: '#d97706', site: 'https://console.anthropic.com/',
    models: ['claude-sonnet-5', 'claude-opus-4-8', 'claude-haiku-4-5'],
    pricing: { input: '$3.00', output: '$15.00', unit: '/M tokens', url: 'https://anthropic.com/pricing' }},
  { id: 'google', label: 'Google Gemini', group: '国际', color: '#4285f4', site: 'https://aistudio.google.com/',
    models: ['gemini-2.0-flash', 'gemini-1.5-pro', 'gemini-1.5-flash'],
    pricing: { input: '$1.25', output: '$5.00', unit: '/M tokens', url: 'https://ai.google.dev/pricing' }},
  { id: 'grok', label: 'xAI Grok', group: '国际', color: '#1a1a2e', site: 'https://console.x.ai/',
    models: ['grok-2', 'grok-2-mini'],
    pricing: { input: '$2.00', output: '$10.00', unit: '/M tokens', url: 'https://x.ai/pricing' }},
  { id: 'mistral', label: 'Mistral AI', group: '国际', color: '#ff6f00', site: 'https://console.mistral.ai/',
    models: ['mistral-large', 'mistral-small', 'codestral'],
    pricing: { input: '$2.00', output: '$6.00', unit: '/M tokens', url: 'https://mistral.ai/pricing' }},
  { id: 'deepseek', label: 'DeepSeek', group: '国内', color: '#4f6ef7', site: 'https://platform.deepseek.com/',
    models: ['deepseek-chat', 'deepseek-reasoner'],
    pricing: { input: '¥1.00', output: '¥2.00', unit: '/M tokens', url: 'https://deepseek.com/pricing' }},
  { id: 'qwen', label: '通义千问（阿里）', group: '国内', color: '#1677ff', site: 'https://dashscope.aliyun.com/',
    models: ['qwen-max', 'qwen-plus', 'qwen-turbo', 'qwen2.5-72b'],
    pricing: { input: '¥2.00', output: '¥6.00', unit: '/M tokens', url: 'https://aliyun.com/pricing' }},
  { id: 'moonshot', label: '月之暗面 Kimi', group: '国内', color: '#8b5cf6', site: 'https://platform.moonshot.cn/',
    models: ['moonshot-v1-128k', 'moonshot-v1-32k', 'moonshot-v1-8k'],
    pricing: { input: '¥1.00', output: '¥2.00', unit: '/M tokens', url: 'https://moonshot.cn/pricing' }},
  { id: 'bytedance', label: '豆包（字节跳动）', group: '国内', color: '#00a76a', site: 'https://console.volcengine.com/',
    models: ['doubao-1.5-pro', 'doubao-1.5-lite', 'doubao-pro'],
    pricing: { input: '¥0.80', output: '¥2.00', unit: '/M tokens', url: 'https://volcengine.com/pricing' }},
  { id: 'minimax', label: 'MiniMax（稀宇）', group: '国内', color: '#ff6900', site: 'https://platform.minimaxi.com/',
    models: ['minimax-text-01', 'minimax-abab-6.5', 'minimax-abab-5.5'],
    pricing: { input: '¥1.00', output: '¥2.00', unit: '/M tokens', url: 'https://minimaxi.com/pricing' }},
  { id: 'baichuan', label: '百川智能', group: '国内', color: '#2b6cb0', site: 'https://platform.baichuan-ai.com/',
    models: ['baichuan4', 'baichuan3-turbo', 'baichuan2-53b'],
    pricing: { input: '¥1.00', output: '¥2.00', unit: '/M tokens', url: 'https://baichuan.com/pricing' }},
  { id: 'stepfun', label: '阶跃星辰', group: '国内', color: '#e00', site: 'https://platform.stepfun.com/',
    models: ['step-2-16k', 'step-1-32k', 'step-1-flash'],
    pricing: { input: '¥0.50', output: '¥2.00', unit: '/M tokens', url: 'https://stepfun.com/pricing' }},
  { id: 'lingyi', label: '零一万物 Yi', group: '国内', color: '#06b6d4', site: 'https://platform.lingyiwanwu.com/',
    models: ['yi-large', 'yi-medium', 'yi-spark', 'yi-vision'],
    pricing: { input: '¥0.50', output: '¥1.50', unit: '/M tokens', url: 'https://lingyiwanwu.com/pricing' }},
  { id: 'ernie', label: '文心一言（百度）', group: '国内', color: '#3060b0', site: 'https://console.bce.baidu.com/',
    models: ['ernie-4.0', 'ernie-3.5', 'ernie-speed', 'ernie-lite'],
    pricing: { input: '¥0.30', output: '¥0.60', unit: '/M tokens', url: 'https://yiyan.baidu.com/pricing' }},
  { id: 'hunyuan', label: '混元（腾讯）', group: '国内', color: '#0052d9', site: 'https://console.cloud.tencent.com/',
    models: ['hunyuan-pro', 'hunyuan-standard', 'hunyuan-lite'],
    pricing: { input: '¥0.50', output: '¥1.00', unit: '/M tokens', url: 'https://cloud.tencent.com/pricing' }},
  { id: 'glm', label: 'GLM（智谱）', group: '国内', color: '#8b4513', site: 'https://open.bigmodel.cn/',
    models: ['glm-4', 'glm-4-plus', 'glm-4-air', 'glm-4-flash'],
    pricing: { input: '¥0.50', output: '¥1.00', unit: '/M tokens', url: 'https://bigmodel.cn/pricing' }},
  { id: 'spark', label: '星火（讯飞）', group: '国内', color: '#e8422a', site: 'https://console.xfyun.cn/',
    models: ['spark-4.0', 'spark-3.0', 'spark-lite'],
    pricing: { input: '¥0.50', output: '¥1.00', unit: '/M tokens', url: 'https://xfyun.cn/pricing' }},
  { id: 'cohere', label: 'Cohere', group: '国际', color: '#d18ee2', site: 'https://dashboard.cohere.com/',
    models: ['command-r-plus', 'command-r', 'command-nightly'],
    pricing: { input: '$1.50', output: '$3.00', unit: '/M tokens', url: 'https://cohere.com/pricing' }},
  { id: 'perplexity', label: 'Perplexity', group: '国际', color: '#1f1f1f', site: 'https://www.perplexity.ai/',
    models: ['sonar-pro', 'sonar-reasoning', 'sonar-deep-research'],
    pricing: { input: '$1.00', output: '$3.00', unit: '/M tokens', url: 'https://perplexity.ai/pricing' }},
  { id: 'ai21', label: 'AI21 Labs', group: '国际', color: '#3b82f6', site: 'https://studio.ai21.com/',
    models: ['jamba-1.5', 'jamba-1.5-mini'],
    pricing: { input: '$0.50', output: '$0.70', unit: '/M tokens', url: 'https://ai21.com/pricing' }},
  { id: 'siliconflow', label: '硅基流动', group: '聚合', color: '#409eff', site: 'https://cloud.siliconflow.cn/',
    models: ['Pro/DeepSeek-V3', 'Pro/Qwen2.5-72B', 'Pro/GLM-4-9B'],
    pricing: { input: '¥0.50', output: '¥1.50', unit: '/M tokens', url: 'https://siliconflow.cn/pricing' }},
  { id: 'openrouter', label: 'OpenRouter', group: '聚合', color: '#6466f1', site: 'https://openrouter.ai/keys',
    models: ['openrouter/auto', 'anthropic/claude-3.5-sonnet', 'openai/gpt-4o'],
    pricing: { input: '参考源', output: '模型', unit: '定价', url: 'https://openrouter.ai/pricing' }},
  { id: 'nvidia', label: 'NVIDIA NIM', group: '聚合', color: '#76b900', site: 'https://build.nvidia.com/',
    models: ['nvidia/llama-3.1-nemotron', 'meta/llama-3.1-8b', 'mistralai/mixtral-8x7b'],
    pricing: { input: '免费', output: '免费', unit: '（有限额）', url: 'https://nvidia.com/pricing' }},
  { id: 'groq', label: 'Groq LPU', group: '聚合', color: '#f55036', site: 'https://console.groq.com/',
    models: ['llama-3.3-70b-versatile', 'mixtral-8x7b-32768', 'gemma2-9b-it'],
    pricing: { input: '免费', output: '免费', unit: '（有限额）', url: 'https://groq.com/pricing' }},
  { id: 'together', label: 'Together AI', group: '聚合', color: '#7678ff', site: 'https://together.ai/',
    models: ['meta-llama/Llama-3.3-70B', 'mistralai/Mixtral-8x7B', 'deepseek-ai/DeepSeek-V3'],
    pricing: { input: '$0.10', output: '$0.40', unit: '/M tokens', url: 'https://together.ai/pricing' }},
  { id: 'azure', label: 'Azure OpenAI', group: '聚合', color: '#0078d4', site: 'https://portal.azure.com/',
    models: ['gpt-4o', 'gpt-4-turbo', 'gpt-35-turbo'],
    pricing: { input: '$2.50', output: '$10.00', unit: '/M tokens', url: 'https://azure.com/pricing' }},
  { id: 'ollama', label: 'Ollama（本地）', group: '本地', color: '#000', site: 'https://ollama.com/',
    models: ['llama3.2', 'qwen2.5', 'deepseek-r1', 'mistral'],
    pricing: { input: '免费', output: '免费', unit: '（本地运行）', url: 'https://ollama.com' }},
]

const groupedProviders = computed(() => {
  const groups: Record<string, ProviderMeta[]> = {}
  for (const p of providerMeta) {
    if (!groups[p.group]) groups[p.group] = []
    groups[p.group].push(p)
  }
  return groups
})

function getProviderMeta(id: string) { return providerMeta.find(m => m.id === id) }

// ===== 派生集合 =====
const isMcp = (p: ProviderConfig) => p.id.startsWith('mcp_') || p.id === 'mcp'
const standardProviders = computed(() => settings.value?.providers?.filter(p => !isMcp(p)) || [])
const mcpConfigs = computed(() => settings.value?.providers?.filter(isMcp) || [])

const newProvider = ref({ id: '', label: '', api_key: '', api_base: '', model: '', is_active: false })

// MCP 新增表单
const newMcp = ref({ name: '', api_key: '', api_base: '', model: 'mcp-default', is_active: false })
const showAddMcp = ref(false)
const mcpErrors = ref<Record<string, string>>({})
function addMcp() {
  const errors: Record<string, string> = {}
  if (!newMcp.value.name.trim()) errors.name = '请输入 MCP 连接名称'
  if (!newMcp.value.api_base.trim()) errors.api_base = '请输入 API 地址'
  mcpErrors.value = errors
  if (Object.keys(errors).length || !settings.value) return
  settings.value.providers.push({
    id: 'mcp_' + Date.now(), label: newMcp.value.name,
    api_key: newMcp.value.api_key, api_base: newMcp.value.api_base,
    model: newMcp.value.model, is_active: newMcp.value.is_active,
    _expanded: true,
  })
  showAddMcp.value = false
  newMcp.value = { name: '', api_key: '', api_base: '', model: 'mcp-default', is_active: false }
}

// ===== 展示辅助 =====
function maskKey(key: string): string {
  if (!key) return '未配置 Key'
  if (key.length <= 8) return '••••'
  return key.slice(0, 3) + '••••••' + key.slice(-4)
}
function modelOptionsFor(p: ProviderConfig): string[] {
  const official = p._official_models
  if (official && official.length) return official
  return getProviderMeta(p.id)?.models || []
}
function currencySymbol(c?: string) { return c === 'USD' ? '$' : '¥' }

// ===== 数据加载 =====
async function loadData() {
  loading.value = true
  loadError.value = ''
  try {
    const [sRes, uRes] = await Promise.all([
      apiGet<{ data: AiSettings | null }>('/api/v1/admin/ai/settings', { skipToast: true }),
      apiGet<{ data: AiUsage }>('/api/v1/admin/ai/usage', { skipToast: true }),
    ])
    settings.value = sRes.data || null
    usage.value = uRes.data || null
    markSaved()
  } catch (e: any) {
    loadError.value = e?.response?.data?.message || '加载失败，请稍后重试'
  } finally { loading.value = false }
}

// ===== 未保存更改检测（不含 _ 前缀 UI 状态） =====
const savedSnapshot = ref('')
function coreSettings() {
  if (!settings.value) return null
  return {
    enabled: settings.value.enabled,
    max_tokens: settings.value.max_tokens,
    tokens_limit: settings.value.tokens_limit,
    providers: (settings.value.providers || []).map(p => ({
      id: p.id, label: p.label, api_key: p.api_key, api_base: p.api_base, model: p.model,
      is_active: !!p.is_active, billing_enabled: !!p.billing_enabled,
      input_price_per_m: p.input_price_per_m ?? 0, output_price_per_m: p.output_price_per_m ?? 0,
      currency: p.currency ?? 'CNY',
    })),
  }
}
function markSaved() { savedSnapshot.value = JSON.stringify(coreSettings()) }
const hasUnsaved = computed(() => settings.value !== null && JSON.stringify(coreSettings()) !== savedSnapshot.value)

async function saveSettings() {
  if (!settings.value) return
  saveStatus.value = 'loading'
  try {
    await apiPut('/api/v1/admin/ai/settings', settings.value, { skipToast: true })
    saveStatus.value = 'success'
    markSaved()
    setTimeout(() => { saveStatus.value = 'idle' }, 1500)
  } catch {
    // 保存失败仅按钮态提示，避免与内联提示重复弹 toast
    saveStatus.value = 'error'
    setTimeout(() => { saveStatus.value = 'idle' }, 3000)
  }
}

async function toggleAi(val: boolean) {
  if (!settings.value) return
  toggleStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/ai/toggle', { enabled: val }, { skipToast: true })
    toggleStatus.value = 'success'
    settings.value.enabled = val
    markSaved()
    setTimeout(() => { toggleStatus.value = 'idle' }, 1500)
  } catch {
    toggleStatus.value = 'error'
    setTimeout(() => { toggleStatus.value = 'idle' }, 3000)
  }
}

// ===== 供应商增删 =====
function addProviderToSettings() {
  const meta = providerMeta.find(m => m.id === newProvider.value.id)
  if (!meta || !settings.value) { newProvider.value.id = ''; return }
  if (settings.value.providers.some(p => p.id === meta.id)) { newProvider.value.id = ''; return }
  const parsePrice = (s: string) => { const m = s?.match(/([\d.]+)/); return m ? parseFloat(m[1]) : 0 }
  const priceCurrency = (s: string) => s?.includes('¥') ? 'CNY' : s?.includes('$') ? 'USD' : 'CNY'
  settings.value.providers.push({
    id: meta.id, label: meta.label, api_key: '', api_base: '',
    model: meta.models[0] || '', is_active: false, billing_enabled: false,
    input_price_per_m: parsePrice(meta.pricing.input),
    output_price_per_m: parsePrice(meta.pricing.output),
    currency: priceCurrency(meta.pricing.input),
    _expanded: true,
  })
  newProvider.value.id = ''
}

async function removeProvider(id: string) {
  if (!settings.value) return
  const meta = getProviderMeta(id)
  const ok = await openConfirm({
    title: `移除 ${meta?.label || id}？`,
    message: '将删除该供应商的本地配置（API Key、模型、计费参数）。保存后生效。',
    danger: true, confirmText: '移除',
  })
  if (!ok) return
  settings.value.providers = settings.value.providers.filter(p => p.id !== id)
}

async function removeMcp(id: string) {
  if (!settings.value) return
  const target = settings.value.providers.find(p => p.id === id)
  const ok = await openConfirm({
    title: `移除 ${target?.label || 'MCP 接口'}？`,
    message: '将删除该 MCP 连接配置。保存后生效。',
    danger: true, confirmText: '移除',
  })
  if (!ok) return
  settings.value.providers = settings.value.providers.filter(p => p.id !== id)
}

// ===== 从官方 API 拉取模型列表（CC Switch 风格） =====
async function fetchModels(p: ProviderConfig) {
  if (p._fetching) return
  p._fetching = true
  p._fetch_msg = ''
  try {
    const res = await apiPost<{ data?: { models?: string[] }; message?: string }>(
      '/api/v1/admin/ai/fetch-models', { provider_id: p.id }, { skipToast: true })
    const official = res.data?.models || []
    p._official_models = official
    p._fetch_msg = official.length ? `已获取 ${official.length} 个官方模型` : '该供应商未返回模型'
  } catch (e: any) {
    p._fetch_msg = e?.response?.data?.message || '获取失败'
  } finally { p._fetching = false }
}

// ===== 连通性测试（New API 渠道测试模式） =====
async function testProvider(p: ProviderConfig) {
  if (p._testing) return
  if (!p.api_key) { p._test_ok = false; p._test_msg = '请先填写 API Key 并保存'; return }
  p._testing = true
  p._test_msg = ''
  try {
    const res = await apiPost<{ data?: { success?: boolean; latency_ms?: number; error?: string }; message?: string }>(
      '/api/v1/admin/ai/test', { provider_id: p.id }, { skipToast: true })
    const d = res.data
    p._test_ok = !!d?.success
    p._test_msg = d?.success ? `连通正常 · ${d.latency_ms}ms` : (d?.error || '连接失败')
  } catch (e: any) {
    p._test_ok = false
    p._test_msg = e?.response?.data?.message || '测试失败'
  } finally { p._testing = false }
}

// ===== 统计 =====
const usageCurrency = computed(() => {
  const bp = usage.value?.by_provider || {}
  const keys = Object.keys(bp)
  return keys.length ? currencySymbol(bp[keys[0]]?.currency) : '¥'
})
// 修复：daily_usage 按日期升序，[0] 是最早一天；取本地今天的日期匹配
const todayTokens = computed(() => {
  const arr = usage.value?.daily_usage || []
  if (!arr.length) return 0
  const d = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  const today = `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
  return arr.find(x => x.date === today)?.tokens ?? arr[arr.length - 1]?.tokens ?? 0
})
const usagePercent = computed(() => {
  const limit = usage.value?.tokens_limit || 0
  if (!limit) return 0
  return Math.min(100, Math.round((usage.value?.tokens_used || 0) / limit * 100))
})
const maxDailyTokens = computed(() => Math.max(...(usage.value?.daily_usage || []).map(x => x.tokens || 0), 1))
const maxDailyCount = computed(() => Math.max(...(usage.value?.daily_usage || []).map(x => x.count || 0), 1))

const filteredLogs = computed(() => {
  if (!usage.value?.recent_logs) return []
  if (!logSearch.value) return usage.value.recent_logs
  const q = logSearch.value.toLowerCase()
  return usage.value.recent_logs.filter(log => (log.student_name || '').toLowerCase().includes(q))
})

onMounted(loadData)
</script>

<template>
  <div class="ai-center">
    <!-- 页头 -->
    <div class="page-head">
      <div>
        <p class="page-crumb">系统配置</p>
        <h2 class="page-title">AI 中心</h2>
      </div>
      <div class="page-actions">
        <span v-if="hasUnsaved && saveStatus === 'idle'" class="unsaved-hint">● 有未保存的更改</span>
        <button class="btn btn-sm" :class="{
          'btn-state-loading': saveStatus === 'loading',
          'btn-state-success': saveStatus === 'success',
          'btn-state-error': saveStatus === 'error',
          'btn-solid': saveStatus === 'idle',
        }" :disabled="saveStatus !== 'idle'" @click="saveSettings">
          <template v-if="saveStatus === 'loading'">保存中...</template>
          <template v-else-if="saveStatus === 'success'">已保存 ✓</template>
          <template v-else-if="saveStatus === 'error'">保存失败 ✗</template>
          <template v-else>💾 保存配置</template>
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>
    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <div class="error-state__msg">{{ loadError }}</div>
      <button class="btn btn-primary btn-sm" @click="loadData">🔄 重试</button>
    </div>

    <template v-else-if="settings">
      <!-- 概览统计 -->
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">🪙</span>
          <div class="stat-card__value">{{ (usage?.tokens_used || 0).toLocaleString() }}</div>
          <div class="stat-card__label">总用量（Token）</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">📈</span>
          <div class="stat-card__value">+{{ todayTokens.toLocaleString() }}</div>
          <div class="stat-card__label">最近一日用量</div>
        </div>
        <div class="stat-card stat-card--info">
          <span class="stat-card__icon">{{ settings.enabled ? '🟢' : '🔴' }}</span>
          <div class="stat-card__value" :style="{ fontSize: '22px', paddingTop: '6px' }">
            {{ settings.enabled ? '运行中' : '已停用' }}
          </div>
          <div class="stat-card__label">AI 服务状态</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">💰</span>
          <div class="stat-card__value">{{ usageCurrency }}{{ (usage?.estimated_cost || 0).toFixed(2) }}</div>
          <div class="stat-card__label">预估费用（本地估算）</div>
        </div>
      </div>

      <!-- 总开关 + 限额 -->
      <div class="card control-card">
        <div class="control-switch">
          <span class="control-label">AI 总开关</span>
          <button type="button" class="switch" :class="{ 'switch--on': settings.enabled }"
            :disabled="toggleStatus === 'loading'" role="switch" :aria-checked="settings.enabled"
            @click="toggleAi(!settings.enabled)">
            <span class="switch__thumb"></span>
          </button>
          <span v-if="toggleStatus !== 'idle'" class="switch-state" :class="'switch-state--' + toggleStatus">
            {{ toggleStatus === 'loading' ? '切换中...' : toggleStatus === 'success' ? '已切换 ✓' : '操作失败 ✗' }}
          </span>
        </div>
        <div class="control-limit">
          <span class="control-label">Token 限额</span>
          <input v-model.number="settings.tokens_limit" type="number" min="0" class="limit-input">
          <span class="limit-used">
            已用 {{ (usage?.tokens_used || 0).toLocaleString() }}
            <span v-if="usagePercent > 0" class="limit-percent" :class="{ 'limit-percent--high': usagePercent >= 90 }">{{ usagePercent }}%</span>
          </span>
        </div>
        <div v-if="usagePercent > 0" class="limit-bar">
          <div class="limit-bar__fill" :class="{ 'limit-bar__fill--high': usagePercent >= 90 }" :style="{ width: usagePercent + '%' }"></div>
        </div>
      </div>

      <!-- 标签导航 -->
      <div class="tab-bar">
        <button :class="['tab-btn', { active: activeTab === 'providers' }]" @click="activeTab = 'providers'">
          🔌 供应商 <span class="tab-count">{{ standardProviders.length }}</span>
        </button>
        <button :class="['tab-btn', { active: activeTab === 'mcp' }]" @click="activeTab = 'mcp'">
          🔗 MCP 接口 <span class="tab-count">{{ mcpConfigs.length }}</span>
        </button>
        <button :class="['tab-btn', { active: activeTab === 'usage' }]" @click="activeTab = 'usage'">📊 用量统计</button>
        <button :class="['tab-btn', { active: activeTab === 'logs' }]" @click="activeTab = 'logs'">📋 对话记录</button>
      </div>

      <!-- ===== 供应商 ===== -->
      <div v-if="activeTab === 'providers'">
        <div class="card add-card">
          <div class="form-group add-form">
            <label>添加供应商</label>
            <select v-model="newProvider.id" class="form-input" @change="addProviderToSettings">
              <option value="">— 从 27 家供应商中选择 —</option>
              <optgroup v-for="(group, gName) in groupedProviders" :key="gName" :label="gName">
                <option v-for="p in group" :key="p.id" :value="p.id" :disabled="standardProviders.some(s => s.id === p.id)">
                  {{ p.label }}{{ standardProviders.some(s => s.id === p.id) ? ' ✓ 已添加' : '' }}
                </option>
              </optgroup>
            </select>
          </div>
        </div>

        <div v-if="!standardProviders.length" class="empty-state">
          <div class="quick-grid">
            <div class="quick-card">
              <div class="quick-card__icon">🌐</div>
              <div class="quick-card__title">国际供应商</div>
              <div class="quick-card__sub">OpenAI · Claude · Gemini</div>
            </div>
            <div class="quick-card">
              <div class="quick-card__icon">🇨🇳</div>
              <div class="quick-card__title">国内供应商</div>
              <div class="quick-card__sub">DeepSeek · 千问 · Kimi · GLM</div>
            </div>
            <div class="quick-card">
              <div class="quick-card__icon">🔗</div>
              <div class="quick-card__title">自定义接口</div>
              <div class="quick-card__sub">MCP · Ollama · vLLM</div>
            </div>
          </div>
          <p class="empty-hint">💡 从上方下拉框选择供应商开始配置</p>
        </div>

        <div v-for="p in standardProviders" :key="p.id" class="provider-card" :class="{ 'provider-card--active': p.is_active }">
          <!-- 卡片头 -->
          <div class="pc-head" @click="p._expanded = !p._expanded">
            <span class="pc-dot" :style="{ background: getProviderMeta(p.id)?.color || 'var(--color-primary)' }"></span>
            <div class="pc-name-wrap">
              <span class="pc-name">{{ getProviderMeta(p.id)?.label || p.label }}</span>
              <span v-if="p.model" class="pc-chip">{{ p.model }}</span>
            </div>
            <span class="pc-key" :class="{ 'pc-key--missing': !p.api_key }">{{ maskKey(p.api_key) }}</span>

            <span v-if="p._test_msg" class="pc-test" :class="p._test_ok ? 'pc-test--ok' : 'pc-test--fail'">{{ p._test_msg }}</span>
            <span v-if="p.tokens_used !== undefined" class="pc-meta">{{ (p.tokens_used || 0).toLocaleString() }} tk · {{ currencySymbol(p.currency) }}{{ (p.estimated_cost || 0).toFixed(3) }}</span>

            <div class="pc-actions" @click.stop>
              <button class="mini-btn" :disabled="p._testing" @click="testProvider(p)">
                {{ p._testing ? '测试中...' : '⚡ 测试' }}
              </button>
              <label class="billing-check" title="启用 Token 计费">
                <input v-model="p.billing_enabled" type="checkbox"> 计费
              </label>
              <button class="mini-btn state-btn" :class="p.is_active ? 'state-btn--on' : 'state-btn--off'" @click="p.is_active = !p.is_active">
                {{ p.is_active ? '● 启用中' : '○ 已停用' }}
              </button>
              <button class="mini-btn mini-btn--danger" @click="removeProvider(p.id)">✕</button>
              <span class="pc-arrow">{{ p._expanded ? '▲' : '▼' }}</span>
            </div>
          </div>

          <!-- 展开配置 -->
          <div v-if="p._expanded" class="pc-body">
            <div class="pc-grid">
              <div class="pc-field">
                <label>API Key</label>
                <div class="key-row">
                  <input v-model="p.api_key" :type="p._show_key ? 'text' : 'password'" class="form-input" placeholder="sk-..." autocomplete="off">
                  <button class="mini-btn" @click.prevent="p._show_key = !p._show_key">{{ p._show_key ? '隐藏' : '显示' }}</button>
                </div>
              </div>
              <div class="pc-field">
                <label>模型</label>
                <div class="key-row">
                  <input v-model="p.model" :list="'model-list-' + p.id" class="form-input" placeholder="选择或输入模型名">
                  <button class="mini-btn" :disabled="p._fetching" title="从官方 API 拉取可用模型" @click="fetchModels(p)">
                    {{ p._fetching ? '拉取中' : '🔄 官方模型' }}
                  </button>
                </div>
                <div v-if="p._fetch_msg" class="field-hint">{{ p._fetch_msg }}</div>
                <datalist :id="'model-list-' + p.id">
                  <option v-for="m in modelOptionsFor(p)" :key="m" :value="m"></option>
                </datalist>
              </div>
              <div class="pc-field">
                <label>API 地址 <span class="label-optional">留空用官方默认</span></label>
                <input v-model="p.api_base" class="form-input" placeholder="https://api.example.com/v1">
              </div>
            </div>
            <div v-if="p.billing_enabled" class="pc-billing">
              <div class="pc-field">
                <label>输入单价（每 M tokens）</label>
                <input v-model.number="p.input_price_per_m" type="number" min="0" step="0.0001" class="form-input">
              </div>
              <div class="pc-field">
                <label>输出单价（每 M tokens）</label>
                <input v-model.number="p.output_price_per_m" type="number" min="0" step="0.0001" class="form-input">
              </div>
              <div class="pc-field">
                <label>币种</label>
                <select v-model="p.currency" class="form-input">
                  <option value="CNY">¥ CNY</option>
                  <option value="USD">$ USD</option>
                </select>
              </div>
              <div v-if="getProviderMeta(p.id)?.pricing" class="pc-field pc-field--hint">
                <label>官方参考价</label>
                <span class="field-hint">
                  输入 {{ getProviderMeta(p.id)!.pricing.input }} / 输出 {{ getProviderMeta(p.id)!.pricing.output }} {{ getProviderMeta(p.id)!.pricing.unit }}
                  · <a :href="getProviderMeta(p.id)!.pricing.url" target="_blank" rel="noopener">定价页 ↗</a>
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== MCP 自定义接口 ===== -->
      <div v-if="activeTab === 'mcp'">
        <div class="card intro-card">
          <div class="intro-title">🔗 MCP 通用接口</div>
          <p class="intro-text">
            连接任意 OpenAI 兼容的 API 服务：自建 vLLM、本地大模型、聚合网关等。填写 OpenAI 格式的 Base URL 即可，格式如 <code>http://主机:端口/v1</code>。
          </p>
          <button v-if="!showAddMcp" class="btn btn-primary btn-sm" @click="showAddMcp = true">➕ 添加 MCP 接口</button>
          <div v-else class="mcp-form">
            <div class="mcp-form__grid">
              <div class="pc-field">
                <label>名称</label>
                <input v-model="newMcp.name" class="form-input" :class="{ 'input-error': mcpErrors.name }" placeholder="如：本地 vLLM">
                <div v-if="mcpErrors.name" class="field-hint field-hint--error">{{ mcpErrors.name }}</div>
              </div>
              <div class="pc-field">
                <label>API 地址 <span class="label-required">*</span></label>
                <input v-model="newMcp.api_base" class="form-input" :class="{ 'input-error': mcpErrors.api_base }" placeholder="http://你的服务器:8000/v1">
                <div v-if="mcpErrors.api_base" class="field-hint field-hint--error">{{ mcpErrors.api_base }}</div>
              </div>
              <div class="pc-field">
                <label>API Key（可选）</label>
                <input v-model="newMcp.api_key" type="password" class="form-input" placeholder="如服务需要鉴权">
              </div>
              <div class="pc-field">
                <label>模型</label>
                <input v-model="newMcp.model" class="form-input" placeholder="mcp-default">
              </div>
            </div>
            <div class="mcp-form__actions">
              <button class="btn btn-solid btn-sm" @click="addMcp">确认添加</button>
              <button class="btn btn-ghost btn-sm" @click="showAddMcp = false; mcpErrors = {}">取消</button>
            </div>
          </div>
        </div>

        <div v-if="!mcpConfigs.length" class="empty-state">暂无 MCP 接口配置</div>
        <div v-for="mcp in mcpConfigs" :key="mcp.id" class="provider-card" :class="{ 'provider-card--active': mcp.is_active }">
          <div class="pc-head" @click="mcp._expanded = !mcp._expanded">
            <span class="pc-dot aic-bg-primary"></span>
            <div class="pc-name-wrap">
              <span class="pc-name">{{ mcp.label }}</span>
              <span v-if="mcp.model" class="pc-chip">{{ mcp.model }}</span>
            </div>
            <span class="pc-key pc-key--wide">{{ mcp.api_base || '未配置地址' }}</span>
            <div class="pc-actions" @click.stop>
              <button class="mini-btn state-btn" :class="mcp.is_active ? 'state-btn--on' : 'state-btn--off'" @click="mcp.is_active = !mcp.is_active">
                {{ mcp.is_active ? '● 启用中' : '○ 已停用' }}
              </button>
              <button class="mini-btn mini-btn--danger" @click="removeMcp(mcp.id)">✕</button>
              <span class="pc-arrow">{{ mcp._expanded ? '▲' : '▼' }}</span>
            </div>
          </div>
          <div v-if="mcp._expanded" class="pc-body">
            <div class="pc-grid">
              <div class="pc-field">
                <label>API Key（可选）</label>
                <input v-model="mcp.api_key" type="password" class="form-input" placeholder="如服务需要鉴权">
              </div>
              <div class="pc-field">
                <label>API 地址</label>
                <input v-model="mcp.api_base" class="form-input" placeholder="http://主机:端口/v1">
              </div>
              <div class="pc-field">
                <label>模型</label>
                <input v-model="mcp.model" class="form-input">
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== 用量统计 ===== -->
      <div v-if="activeTab === 'usage'" class="card usage-card">
        <div class="usage-trio">
          <div class="usage-box">
            <div class="usage-box__label">已用 Token</div>
            <div class="usage-box__value">{{ (usage?.tokens_used || 0).toLocaleString() }}</div>
          </div>
          <div class="usage-box">
            <div class="usage-box__label">限额</div>
            <div class="usage-box__value">{{ (usage?.tokens_limit || 0).toLocaleString() }}</div>
          </div>
          <div class="usage-box">
            <div class="usage-box__label">对话次数</div>
            <div class="usage-box__value">{{ usage?.total_conversations || 0 }}</div>
          </div>
        </div>

        <div class="section-title">各供应商用量明细</div>
        <div class="data-table">
          <table>
            <thead>
              <tr><th>供应商</th><th class="aic-right">Token</th><th class="aic-right">调用</th><th class="aic-right">预估费用</th></tr>
            </thead>
            <tbody>
              <tr v-for="(val, key) in usage?.by_provider || {}" :key="key">
                <td class="cell-strong">{{ key }}</td>
                <td class="cell-num cell-primary">{{ (val.tokens || 0).toLocaleString() }}</td>
                <td class="cell-num">{{ val.total_calls || 0 }}</td>
                <td class="cell-num cell-strong">{{ currencySymbol(val.currency) }}{{ (val.estimated_cost || 0).toFixed(4) }}</td>
              </tr>
              <tr v-if="!Object.keys(usage?.by_provider || {}).length">
                <td colspan="4" class="cell-empty">暂无数据</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="usage-note">
          📊 Token 用量来自每次 AI 对话的 API 响应累加；预估费用基于供应商公开定价 × 实际 Token 数计算，仅供参考，实际以供应商账单为准。
        </div>

        <template v-if="usage?.daily_usage?.length">
          <div class="section-title">近 {{ usage.daily_usage.length }} 日趋势</div>
          <div class="trend-list">
            <div v-for="d in usage.daily_usage" :key="d.date" class="trend-row">
              <span class="trend-date">{{ d.date?.slice(5) || d.date }}</span>
              <div class="trend-bars">
                <div class="trend-bar-row">
                  <div class="trend-bar trend-bar--tokens" :style="{ width: Math.max(Math.min((d.tokens || 0) / maxDailyTokens * 100, 100), 1) + '%' }"></div>
                  <span class="trend-num trend-num--tokens">{{ (d.tokens || 0).toLocaleString() }}</span>
                </div>
                <div class="trend-bar-row">
                  <div class="trend-bar trend-bar--calls" :style="{ width: Math.max(Math.min((d.count || 0) / maxDailyCount * 100, 100), 1) + '%' }"></div>
                  <span class="trend-num">{{ d.count || 0 }} 次</span>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- ===== 对话记录 ===== -->
      <div v-if="activeTab === 'logs'" class="card logs-card">
        <div v-if="usage?.recent_logs?.length" class="logs-filter">
          <input v-model="logSearch" class="form-input" placeholder="🔍 搜索学生姓名...">
          <span class="logs-count">共 {{ usage.recent_logs.length }} 条</span>
        </div>
        <div v-if="filteredLogs.length" class="logs-list">
          <div v-for="log in filteredLogs" :key="log.id" class="log-item">
            <div class="log-head">
              <span class="log-student">{{ log.student_name || '匿名' }}</span>
              <span v-if="log.provider" class="log-provider">{{ log.provider }}</span>
              <span class="log-tokens">{{ log.tokens_used }} tk</span>
              <span v-if="log.cost" class="log-cost">{{ log.cost }} {{ log.currency }}</span>
              <span class="log-time">{{ log.created_at }}</span>
            </div>
            <div class="log-q"><strong>问：</strong>{{ log.question }}</div>
            <div class="log-a"><strong>答：</strong>{{ log.answer?.substring(0, 200) }}{{ (log.answer?.length || 0) > 200 ? '...' : '' }}</div>
          </div>
        </div>
        <div v-else class="empty-state">{{ usage?.recent_logs?.length ? '无匹配记录' : '暂无对话记录' }}</div>
      </div>
    </template>
  </div>
</template>

<style scoped>
/* ===== 页面骨架（Linux.do / AstrBot 式清爽布局） ===== */
.ai-center { max-width: 880px; }
.page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; gap: 12px; flex-wrap: wrap; }
.page-crumb { font-size: 12px; color: var(--color-text-secondary); margin-bottom: 2px; }
.page-title { font-size: 24px; font-weight: 700; }
.page-actions { display: flex; align-items: center; gap: 10px; }
.unsaved-hint { font-size: 12px; color: var(--color-warning-text); }

/* ===== 控制卡（总开关 + 限额） ===== */
.control-card { padding: 16px 20px; margin-bottom: 16px; display: flex; align-items: center; gap: 16px 24px; flex-wrap: wrap; }
.control-switch { display: flex; align-items: center; gap: 10px; }
.control-label { font-size: 13px; font-weight: 600; color: var(--color-text); }
.switch { position: relative; width: 44px; height: 24px; border-radius: 12px; border: none; cursor: pointer; background: var(--tint-4); transition: background 0.2s; padding: 0; }
.switch--on { background: var(--color-primary); }
.switch:disabled { opacity: 0.6; cursor: wait; }
.switch__thumb { position: absolute; top: 2px; left: 2px; width: 20px; height: 20px; border-radius: 50%; background: #fff; box-shadow: 0 1px 3px rgba(0,0,0,0.25); transition: transform 0.2s var(--ease-smooth); }
.switch--on .switch__thumb { transform: translateX(20px); }
.switch-state { font-size: 12px; }
.switch-state--loading { color: var(--color-warning-text); }
.switch-state--success { color: var(--color-success-text); }
.switch-state--error { color: var(--color-danger-text); }
.control-limit { display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--color-text-secondary); flex: 1; min-width: 220px; }
.limit-input { width: 96px; padding: 5px 10px; border: 1px solid var(--tint-3); border-radius: 8px; font-size: 12px; background: var(--tint-1); color: var(--color-text); }
.limit-used { color: var(--color-text-secondary); }
.limit-percent { margin-left: 6px; font-weight: 700; color: var(--color-primary); }
.limit-percent--high { color: var(--color-danger-text); }
.limit-bar { width: 100%; height: 6px; background: var(--tint-2); border-radius: 3px; overflow: hidden; }
.limit-bar__fill { height: 100%; background: linear-gradient(90deg, var(--color-primary), var(--color-primary-light)); border-radius: 3px; transition: width 0.4s; }
.limit-bar__fill--high { background: linear-gradient(90deg, #f59e0b, #ef4444); }

/* ===== 标签导航 ===== */
.tab-bar { display: flex; gap: 4px; margin-bottom: 16px; background: var(--tint-1); border: 1px solid var(--tint-2); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 8px 10px; border: none; border-radius: 9px; font-size: 12px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; white-space: nowrap; font-family: inherit; }
.tab-btn:hover { color: var(--color-text); background: var(--tint-2); }
.tab-btn.active { background: var(--color-primary); color: #fff; box-shadow: 0 2px 8px rgba(124, 58, 237, 0.25); }
.tab-count { opacity: 0.7; font-weight: 500; margin-left: 2px; }

/* ===== 添加行 ===== */
.add-card { padding: 14px 20px; margin-bottom: 12px; }
.add-card:hover { box-shadow: none; }
.add-form { margin-bottom: 0; }

/* ===== 空态引导 ===== */
.quick-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; max-width: 540px; margin: 0 auto; }
.quick-card { padding: 16px 10px; background: var(--tint-1); border-radius: 12px; border: 1px dashed var(--tint-3); text-align: center; }
.quick-card__icon { font-size: 22px; margin-bottom: 6px; }
.quick-card__title { font-size: 12px; font-weight: 600; color: var(--color-text); }
.quick-card__sub { font-size: 11px; color: var(--color-text-secondary); margin-top: 2px; }
.empty-hint { font-size: 12px; color: var(--color-text-secondary); margin-top: 12px; }

/* ===== 供应商卡片（CC Switch 风格） ===== */
.provider-card { border: 1px solid var(--tint-3); border-radius: 14px; overflow: hidden; background: var(--color-bg-card); margin-bottom: 10px; transition: border-color 0.2s, box-shadow 0.2s; }
.provider-card:hover { border-color: var(--tint-4); }
.provider-card--active { border-color: rgba(16, 185, 129, 0.45); box-shadow: 0 0 0 1px rgba(16, 185, 129, 0.2); }
.pc-head { display: flex; align-items: center; gap: 10px; padding: 12px 16px; cursor: pointer; flex-wrap: wrap; }
.pc-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; box-shadow: 0 0 0 3px var(--tint-1); }
.pc-name-wrap { display: flex; align-items: center; gap: 8px; min-width: 0; flex: 1 1 180px; }
.pc-name { font-weight: 600; font-size: 14px; white-space: nowrap; }
.pc-chip { font-size: 11px; padding: 2px 8px; border-radius: 6px; background: var(--tint-2); color: var(--color-text-secondary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 160px; }
.pc-key { font-size: 11px; color: var(--color-text-secondary); font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; white-space: nowrap; }
.pc-key--missing { color: var(--color-warning-text); }
.pc-key--wide { max-width: 220px; overflow: hidden; text-overflow: ellipsis; }
.pc-test { font-size: 11px; padding: 2px 8px; border-radius: 6px; white-space: nowrap; }
.pc-test--ok { background: rgba(16, 185, 129, 0.1); color: var(--color-success-text); }
.pc-test--fail { background: rgba(239, 68, 68, 0.1); color: var(--color-danger-text); }
.pc-meta { font-size: 11px; color: var(--color-text-secondary); white-space: nowrap; }
.pc-actions { display: flex; align-items: center; gap: 6px; margin-left: auto; }
.pc-arrow { font-size: 11px; color: var(--color-text-secondary); }

.mini-btn { padding: 4px 10px; border-radius: 7px; font-size: 11px; cursor: pointer; border: 1px solid var(--tint-3); background: var(--color-bg-card); color: var(--color-text-secondary); font-family: inherit; transition: all 0.15s; white-space: nowrap; }
.mini-btn:hover:not(:disabled) { border-color: var(--tint-4); color: var(--color-text); }
.mini-btn:disabled { opacity: 0.55; cursor: not-allowed; }
.mini-btn--danger:hover { border-color: rgba(239, 68, 68, 0.4); color: var(--color-danger-text); background: rgba(239, 68, 68, 0.06); }
.state-btn--on { background: rgba(16, 185, 129, 0.12); border-color: rgba(16, 185, 129, 0.35); color: var(--color-success-text); font-weight: 600; }
.state-btn--off { background: transparent; }
.billing-check { display: flex; align-items: center; gap: 3px; font-size: 11px; color: var(--color-text-secondary); cursor: pointer; user-select: none; }
.billing-check input { accent-color: var(--color-primary); }

/* ===== 展开配置区 ===== */
.pc-body { padding: 4px 16px 14px; border-top: 1px solid var(--tint-2); }
.pc-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; padding-top: 10px; }
.pc-billing { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; padding-top: 10px; margin-top: 10px; border-top: 1px dashed var(--tint-2); }
.pc-field label { display: block; font-size: 11px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 4px; }
.label-optional { font-weight: 400; opacity: 0.75; }
.label-required { color: var(--color-danger-text); }
.pc-field--hint { display: flex; flex-direction: column; justify-content: end; padding-bottom: 6px; }
.pc-field .form-input { padding: 7px 10px; font-size: 12px; }
.key-row { display: flex; gap: 6px; align-items: center; }
.key-row .form-input { flex: 1; min-width: 0; }
.key-row .mini-btn { flex-shrink: 0; }
.field-hint { font-size: 10px; color: var(--color-text-secondary); margin-top: 3px; }
.field-hint--error { color: var(--color-danger-text); }
.input-error { border-color: rgba(239, 68, 68, 0.5) !important; }

/* ===== MCP ===== */
.intro-card { padding: 18px 20px; margin-bottom: 12px; }
.intro-title { font-size: 14px; font-weight: 600; margin-bottom: 4px; }
.intro-text { font-size: 12px; color: var(--color-text-secondary); margin-bottom: 12px; line-height: 1.7; }
.intro-text code { background: var(--tint-2); padding: 1px 6px; border-radius: 5px; font-size: 11px; }
.mcp-form { padding: 14px; background: var(--tint-1); border-radius: 10px; }
.mcp-form__grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-bottom: 12px; }
.mcp-form__actions { display: flex; gap: 8px; }

/* ===== 用量统计 ===== */
.usage-card { padding: 20px; }
.usage-trio { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 18px; }
.usage-box { padding: 14px; background: var(--tint-1); border-radius: 10px; text-align: center; }
.usage-box__label { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 4px; }
.usage-box__value { font-size: 22px; font-weight: 700; }
.section-title { font-size: 12px; font-weight: 600; color: var(--color-text-secondary); margin: 4px 0 8px; }
.cell-strong { font-weight: 600; }
.cell-num { text-align: right; font-variant-numeric: tabular-nums; }
.cell-primary { color: var(--color-primary); font-weight: 600; }
.cell-empty { text-align: center; color: var(--color-text-secondary); }
.usage-note { font-size: 11px; color: var(--color-text-secondary); background: var(--tint-1); padding: 8px 12px; border-radius: 8px; margin: 14px 0; line-height: 1.7; }
.trend-list { border: 1px solid var(--tint-2); border-radius: 10px; overflow: hidden; }
.trend-row { display: flex; align-items: center; gap: 10px; padding: 7px 14px; border-bottom: 1px solid var(--tint-2); font-size: 11px; }
.trend-row:last-child { border-bottom: none; }
.trend-date { width: 44px; flex-shrink: 0; color: var(--color-text); font-variant-numeric: tabular-nums; }
.trend-bars { flex: 1; display: flex; flex-direction: column; gap: 3px; }
.trend-bar-row { display: flex; align-items: center; gap: 6px; }
.trend-bar { height: 10px; border-radius: 4px; min-width: 3px; transition: width 0.4s; }
.trend-bar--tokens { background: linear-gradient(90deg, var(--color-primary), var(--color-primary-light)); }
.trend-bar--calls { height: 7px; background: var(--md-gold, #d97706); opacity: 0.7; }
.trend-num { color: var(--color-text-secondary); white-space: nowrap; }
.trend-num--tokens { color: var(--color-primary); font-weight: 600; }

/* ===== 对话记录 ===== */
.logs-card { padding: 18px 20px; }
.logs-filter { display: flex; gap: 8px; margin-bottom: 12px; align-items: center; }
.logs-filter .form-input { max-width: 240px; padding: 7px 12px; font-size: 12px; }
.logs-count { font-size: 11px; color: var(--color-text-secondary); }
.logs-list { border: 1px solid var(--tint-2); border-radius: 10px; overflow: hidden; }
.log-item { padding: 10px 14px; border-bottom: 1px solid var(--tint-2); font-size: 12px; }
.log-item:last-child { border-bottom: none; }
.log-head { display: flex; gap: 8px; margin-bottom: 3px; flex-wrap: wrap; align-items: baseline; }
.log-student { font-weight: 600; flex: 1; }
.log-provider, .log-tokens, .log-cost, .log-time { font-size: 11px; color: var(--color-text-secondary); }
.log-tokens { color: var(--color-primary); }
.log-cost { color: var(--md-gold, #d97706); }
.log-q { color: var(--color-text); line-height: 1.6; }
.log-a { color: var(--color-text-secondary); line-height: 1.6; }

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .pc-grid, .pc-billing { grid-template-columns: 1fr; }
  .mcp-form__grid { grid-template-columns: 1fr; }
  .usage-trio { grid-template-columns: 1fr; }
  .quick-grid { grid-template-columns: 1fr; }
  .pc-key, .pc-meta { display: none; }
}
.aic-right { text-align:right; }
.aic-bg-primary { background: var(--color-primary); }
</style>
