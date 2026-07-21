<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()

interface ProviderConfig { id: string; label: string; api_key: string; api_base: string; model: string; is_active: boolean; _expanded?: boolean; billing_enabled?: boolean; tokens_used?: number; total_calls?: number; estimated_cost?: number }
interface AiSettings { enabled: boolean; max_tokens: number; tokens_used: number; tokens_limit: number; providers: ProviderConfig[] }
interface DailyUsage { date: string; tokens: number; count: number }
interface ConversationLog { id: number; student_name: string; question: string; answer: string; tokens_used: number; created_at: string }
interface AiUsage { enabled: boolean; tokens_used: number; tokens_limit: number; total_conversations: number; daily_usage: DailyUsage[]; recent_logs: ConversationLog[] }

const loading = ref(true)
const settings = ref<AiSettings | null>(null)
const usage = ref<AiUsage | null>(null)
const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const toggleStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const activeTab = ref<'providers' | 'mcp' | 'usage' | 'logs'>('providers')
const logSearch = ref('')
const filteredLogs = computed(() => {
  if (!usage.value?.recent_logs) return []
  if (!logSearch.value) return usage.value.recent_logs
  const q = logSearch.value.toLowerCase()
  return usage.value.recent_logs.filter(log =>
    (log.student_name || '').toLowerCase().includes(q)
  )
})

// ===== 供应商元数据（含计费信息） =====
interface PricingInfo { input: string; output: string; unit: string; url: string }
interface ProviderMeta { id: string; label: string; group: string; color: string; models: string[]; pricing: PricingInfo; site: string }

const providerMeta: ProviderMeta[] = [
  { id: 'openai', label: 'OpenAI', group: '国际', color: '#10a37f', site: 'https://platform.openai.com/api-keys',
    models: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-3.5-turbo'],
    pricing: { input: '$2.50', output: '$10.00', unit: '/M tokens', url: 'https://openai.com/pricing' }},
  { id: 'claude', label: 'Anthropic Claude', group: '国际', color: '#d97706', site: 'https://console.anthropic.com/',
    models: ['claude-3-5-sonnet', 'claude-3-haiku', 'claude-3-opus'],
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
    models: ['deepseek-chat', 'deepseek-v3', 'deepseek-r1', 'deepseek-coder'],
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

const showAddProvider = ref(false)
const newProvider = ref({ id: '', label: '', api_key: '', api_base: '', model: '', is_active: false })
const activeProviderId = ref('')
const activeProviderMeta = computed(() => providerMeta.find(m => m.id === activeProviderId.value))

// MCP 配置
const mcpConfigs = computed({
  get: () => settings.value?.providers?.filter(p => p.id === 'mcp') || [],
  set: (val: ProviderConfig[]) => {
    if (!settings.value) return
    settings.value.providers = [...(settings.value.providers?.filter(p => p.id !== 'mcp') || []), ...val]
  }
})
const newMcp = ref({ name: '', api_key: '', api_base: '', model: 'mcp-default', is_active: false })
const showAddMcp = ref(false)

const mcpErrors = reactive<Record<string, string>>({})
function mcpClr(f: string) { delete mcpErrors[f] }
function mcpVld(): boolean {
  if (!newMcp.value.name.trim()) { mcpErrors.name = '请输入 MCP 连接名称'; return false }
  if (!newMcp.value.api_base.trim()) { mcpErrors.api_base = '请输入 API 地址'; return false }
  return true
}
function addMcp() {
  Object.keys(mcpErrors).forEach(k => delete mcpErrors[k])
  if (!mcpVld()) return
  if (!settings.value?.providers) settings.value.providers = []
  settings.value.providers.push({
    id: 'mcp_' + Date.now(), label: newMcp.value.name,
    api_key: newMcp.value.api_key, api_base: newMcp.value.api_base,
    model: newMcp.value.model, is_active: newMcp.value.is_active,
  })
  showAddMcp.value = false
  newMcp.value = { name: '', api_key: '', api_base: '', model: 'mcp-default', is_active: false }
}

function removeMcp(idx: number) {
  if (!settings.value?.providers) return
  const mcpIds = settings.value.providers.filter(p => p.id.startsWith('mcp_') || p.id === 'mcp').map(p => p.id)
  const actualIdx = settings.value.providers.findIndex(p => p.id === mcpIds[idx])
  if (actualIdx >= 0) settings.value.providers.splice(actualIdx, 1)
}

// Standard providers (non-MCP)
const standardProviders = computed({
  get: () => settings.value?.providers?.filter(p => !p.id.startsWith('mcp_') && p.id !== 'mcp') || [],
  set: (val: ProviderConfig[]) => {
    if (!settings.value) return
    const mcps = settings.value.providers?.filter(p => p.id.startsWith('mcp_') || p.id === 'mcp') || []
    settings.value.providers = [...val, ...mcps]
  }
})

async function loadData() {
  loading.value = true
  try {
    const [sRes, uRes] = await Promise.all([
      fetch('/api/v1/admin/ai/settings', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json()),
      fetch('/api/v1/admin/ai/usage', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json()),
    ])
    settings.value = sRes.data || null
    usage.value = uRes.data || null
  } catch { toast.show('加载失败', 'error') }
  finally { loading.value = false }
}

async function saveSettings() {
  if (!settings.value) return
  saveStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/admin/ai/settings', {
      method: 'PUT', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify(settings.value),
    })
    const data = await res.json()
    if (!res.ok) { saveStatus.value = 'error'; setTimeout(() => { saveStatus.value = 'idle' }, 3000); return }
    saveStatus.value = 'success'
    loadData()
    setTimeout(() => { saveStatus.value = 'idle' }, 1500)
  } catch { saveStatus.value = 'error'; setTimeout(() => { saveStatus.value = 'idle' }, 3000) }
}

async function toggleAi(val: boolean) {
  if (!settings.value) return
  toggleStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/admin/ai/toggle', {
      method: 'POST', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: val }),
    })
    const data = await res.json()
    if (!res.ok) { toggleStatus.value = 'error'; setTimeout(() => { toggleStatus.value = 'idle' }, 3000); return }
    toggleStatus.value = 'success'
    settings.value.enabled = val
    setTimeout(() => { toggleStatus.value = 'idle' }, 1500)
  } catch { toggleStatus.value = 'error'; setTimeout(() => { toggleStatus.value = 'idle' }, 3000) }
}

function addProviderToSettings() {
  if (!newProvider.value.id) return
  const meta = providerMeta.find(m => m.id === newProvider.value.id)
  if (!meta) return
  if (!settings.value?.providers) settings.value.providers = []
  if (settings.value.providers.some(p => p.id === meta.id)) { newProvider.value.id = ''; return }
  settings.value.providers.push({
    id: meta.id, label: meta.label, api_key: '', api_base: '',
    model: meta.models[0] || '', is_active: false, billing_enabled: false,
    _expanded: false,
  })
  newProvider.value.id = ''
}

function removeProvider(idx: number) {
  if (!settings.value?.providers) return
  const std = standardProviders.value
  const actualIdx = settings.value.providers.findIndex(p => p.id === std[idx]?.id)
  if (actualIdx >= 0) settings.value.providers.splice(actualIdx, 1)
}

function getProviderMeta(id: string) { return providerMeta.find(m => m.id === id) }

onMounted(loadData)
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:20px;">
      <div><p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">系统配置</p><h2 style="font-size:24px;font-weight:700;">🤖 AI 中心</h2></div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>
    <template v-else>
      <!-- 顶部统计卡片 -->
      <div class="stats-grid" style="display:grid;grid-template-columns:repeat(4,1fr);gap:12px;margin-bottom:16px;">
        <div style="padding:14px 16px;background:var(--color-bg);border-radius:10px;border:1px solid var(--color-border);">
          <div style="font-size:11px;color:var(--color-text-secondary);">总用量</div>
          <div style="font-size:22px;font-weight:700;color:var(--color-text);">{{ (usage?.tokens_used || 0).toLocaleString() }}</div>
          <div style="font-size:11px;color:var(--color-text-secondary);">Token</div>
        </div>
        <div style="padding:14px 16px;background:var(--color-bg);border-radius:10px;border:1px solid var(--color-border);">
          <div style="font-size:11px;color:var(--color-text-secondary);">今日用量</div>
          <div style="font-size:22px;font-weight:700;color:#10b981;">+{{ (usage?.daily_usage?.[0]?.tokens || 0).toLocaleString() }}</div>
          <div style="font-size:11px;color:var(--color-text-secondary);">Token</div>
        </div>
        <div style="padding:14px 16px;background:var(--color-bg);border-radius:10px;border:1px solid var(--color-border);">
          <div style="font-size:11px;color:var(--color-text-secondary);">AI 状态</div>
          <div style="font-size:22px;font-weight:700;color:settings?.enabled ? '#10b981' : '#f87171';">
            {{ settings?.enabled ? '🟢 运行中' : '🔴 已停用' }}
          </div>
          <div style="font-size:11px;color:var(--color-text-secondary);">总开关</div>
        </div>
        <div style="padding:14px 16px;background:var(--color-bg);border-radius:10px;border:1px solid var(--color-border);">
          <div style="font-size:11px;color:var(--color-text-secondary);">预估费用</div>
          <div style="font-size:22px;font-weight:700;color:#f59e0b;">${{ (usage?.estimated_cost || 0).toFixed(2) }}</div>
          <div style="font-size:11px;color:var(--color-text-secondary);">本月累计</div>
        </div>
      </div>

      <!-- 开关 + 限额 -->
      <div class="card" style="max-width:720px;padding:16px 20px;margin-bottom:12px;display:flex;align-items:center;gap:20px;flex-wrap:wrap;">
        <div style="display:flex;align-items:center;gap:10px;">
          <span style="font-size:14px;font-weight:600;">AI 总开关</span>
          <label style="position:relative;display:inline-block;width:44px;height:24px;cursor:pointer;">
            <input type="checkbox" :checked="settings?.enabled" @change="toggleAi(($event.target as HTMLInputElement).checked)" style="opacity:0;width:0;height:0;">
            <span :style="{ position:'absolute',inset:0,background:settings?.enabled ? '#7c3aed' : '#ccc',borderRadius:'12px' }">
              <span :style="{ position:'absolute',top:'2px',left:settings?.enabled ? '22px' : '2px',width:'20px',height:'20px',borderRadius:'50%',background:'#fff',boxShadow:'0 1px 3px rgba(0,0,0,0.2)' }"></span>
            </span>
          </label>
          <span v-if="toggleStatus !== 'idle'" :style="{ fontSize: '11px', color: toggleStatus === 'loading' ? '#f59e0b' : toggleStatus === 'success' ? '#10b981' : '#ef4444' }">
            {{ toggleStatus === 'loading' ? '切换中...' : toggleStatus === 'success' ? '已切换 ✓' : '操作失败 ✗' }}
          </span>
        </div>
        <div style="display:flex;align-items:center;gap:6px;font-size:12px;color:var(--color-text-secondary);">
          <span>限额</span>
          <input v-model.number="settings!.tokens_limit" type="number" min="0" style="width:80px;padding:4px 8px;border:1px solid var(--color-border);border-radius:6px;font-size:12px;background:var(--color-bg-card);color:var(--color-text);">
          <span>Token</span>
          <span style="color:var(--color-text);font-weight:600;">已用 {{ (usage?.tokens_used || 0).toLocaleString() }}</span>
        </div>
        <button class="btn btn-sm" :style="{ background: saveStatus === 'loading' ? '#f59e0b' : saveStatus === 'success' ? '#10b981' : saveStatus === 'error' ? '#ef4444' : '#7c3aed', color: '#fff', border: '1px solid transparent' }" :disabled="saveStatus !== 'idle'" @click="saveSettings">
          <template v-if="saveStatus === 'loading'">保存中...</template>
          <template v-else-if="saveStatus === 'success'">已保存 ✓</template>
          <template v-else-if="saveStatus === 'error'">保存失败 ✗</template>
          <template v-else>💾 保存</template>
        </button>
      </div>

      <!-- 标签导航 -->
      <div class="tab-bar" style="max-width:720px;">
        <button :class="['tab-btn', { active: activeTab === 'providers' }]" @click="activeTab = 'providers'">🔌 AI 供应商</button>
        <button :class="['tab-btn', { active: activeTab === 'mcp' }]" @click="activeTab = 'mcp'">🔗 MCP 自定义接口</button>
        <button :class="['tab-btn', { active: activeTab === 'usage' }]" @click="activeTab = 'usage'">📊 用量</button>
        <button :class="['tab-btn', { active: activeTab === 'logs' }]" @click="activeTab = 'logs'">📋 记录</button>
      </div>

      <!-- ===== AI 供应商 ===== -->
      <div v-if="activeTab === 'providers'" style="max-width:720px;">
        <!-- 下拉选择器 -->
        <div class="card" style="padding:16px;margin-bottom:12px;">
          <div style="display:flex;gap:12px;align-items:center;">
            <div class="form-group" style="flex:1;">
              <label>添加供应商</label>
              <select v-model="newProvider.id" class="form-input" @change="addProviderToSettings">
                <option value="">— 从列表中选择 —</option>
                <optgroup v-for="(group, gName) in groupedProviders" :key="gName" :label="gName">
                  <option v-for="p in group" :key="p.id" :value="p.id" :disabled="standardProviders.some(s => s.id === p.id)">
                    {{ p.label }} {{ standardProviders.some(s => s.id === p.id) ? '✓' : '' }}
                  </option>
                </optgroup>
              </select>
            </div>
            <div style="font-size:11px;color:var(--color-text-secondary);padding-top:14px;">已配置 {{ standardProviders.length }} 个</div>
          </div>
        </div>

        <!-- 已配置的供应商 -->
        <div v-if="!standardProviders.length" style="text-align:center;padding:16px;">
          <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:10px;max-width:480px;margin:0 auto;">
            <div style="padding:14px 10px;background:var(--color-bg);border-radius:10px;border:1px dashed var(--color-border);text-align:center;">
              <div style="font-size:20px;margin-bottom:4px;">🌐</div>
              <div style="font-size:11px;font-weight:600;color:var(--color-text);">国际供应商</div>
              <div style="font-size:10px;color:var(--color-text-secondary);margin-top:2px;">OpenAI · Claude · Gemini</div>
            </div>
            <div style="padding:14px 10px;background:var(--color-bg);border-radius:10px;border:1px dashed var(--color-border);text-align:center;">
              <div style="font-size:20px;margin-bottom:4px;">🇨🇳</div>
              <div style="font-size:11px;font-weight:600;color:var(--color-text);">国内供应商</div>
              <div style="font-size:10px;color:var(--color-text-secondary);margin-top:2px;">DeepSeek · 千问 · Kimi</div>
            </div>
            <div style="padding:14px 10px;background:var(--color-bg);border-radius:10px;border:1px dashed var(--color-border);text-align:center;">
              <div style="font-size:20px;margin-bottom:4px;">🔗</div>
              <div style="font-size:11px;font-weight:600;color:var(--color-text);">自定义接口</div>
              <div style="font-size:10px;color:var(--color-text-secondary);margin-top:2px;">MCP · Ollama · vLLM</div>
            </div>
          </div>
          <p style="font-size:11px;color:var(--color-text-secondary);margin-top:10px;">💡 从上方下拉框选择供应商开始配置</p>
        </div>
        <div v-for="(p, i) in standardProviders" :key="p.id" style="margin-bottom:8px;border:1px solid var(--color-border);border-radius:10px;overflow:hidden;">
          <!-- 供应商头部（可点击折叠） -->
          <div @click="p._expanded = !p._expanded" :style="{ background: (getProviderMeta(p.id)?.color || '#7c3aed') + '0a', padding:'10px 14px', display:'flex', alignItems:'center', gap:'10px', cursor:'pointer', borderBottom: p._expanded ? '1px solid var(--color-border)' : 'none' }">
            <span :style="{ width:'10px',height:'10px',borderRadius:'50%',background:p.is_active ? '#10B981' : '#ccc',flexShrink:0 }"></span>
            <span style="font-weight:600;font-size:14px;flex:1;">{{ getProviderMeta(p.id)?.label || p.label }}</span>
            <span v-if="p.tokens_used !== undefined" style="font-size:11px;color:var(--color-text-secondary);">📊 {{ (p.tokens_used||0).toLocaleString() }}</span>
            <span v-if="getProviderMeta(p.id)?.pricing" style="font-size:11px;color:var(--color-text-secondary);">💰 {{ getProviderMeta(p.id)!.pricing.input }}</span>
            <label style="display:flex;align-items:center;gap:3px;font-size:11px;color:var(--color-text-secondary);cursor:pointer;" @click.stop>
              <input type="checkbox" v-model="p.billing_enabled" style="accent-color:#7c3aed;"> 计费
            </label>
            <span :style="{ padding:'2px 10px',borderRadius:'12px',fontSize:'10px',fontWeight:600, background:p.is_active ? '#10B98120' : '#ccc20', color:p.is_active ? '#10B981' : 'var(--color-text-secondary)' }">
              {{ p.is_active ? '启用' : '禁用' }}
            </span>
            <button :style="{ padding:'3px 8px',borderRadius:'6px',fontSize:'10px',cursor:'pointer',border:'1px solid', borderColor: p.is_active ? '#10B981' : 'var(--color-border)', background:p.is_active ? '#10B981' : 'transparent', color:p.is_active ? '#fff' : 'var(--color-text-secondary)', fontFamily:'inherit' }" @click.stop="p.is_active = !p.is_active">{{ p.is_active ? '已启用' : '已禁用' }}</button>
            <button style="padding:3px 8px;borderRadius:6px;fontSize:11px;cursor:pointer;border:1px solid var(--color-border);background:transparent;color:var(--color-text-secondary);fontFamily:'inherit';" @click.stop="removeProvider(i)">✕</button>
            <span style="font-size:12px;color:var(--color-text-secondary);">{{ p._expanded ? '▲' : '▼' }}</span>
          </div>
          <!-- 配置详情（折叠） -->
          <div v-if="p._expanded" style="padding:10px 14px;display:grid;grid-template-columns:1fr 1fr 1fr;gap:8px;">
            <div><label style="display:block;font-size:10px;color:var(--color-text-secondary);margin-bottom:2px;">API Key</label><input v-model="p.api_key" type="password" class="form-input" placeholder="sk-..." style="font-size:11px;padding:5px 8px;"></div>
            <div><label style="display:block;font-size:10px;color:var(--color-text-secondary);margin-bottom:2px;">模型</label><select v-model="p.model" class="form-input" style="font-size:11px;padding:5px 8px;"><option v-for="m in getProviderMeta(p.id)?.models || []" :key="m" :value="m">{{ m }}</option></select></div>
            <div><label style="display:block;font-size:10px;color:var(--color-text-secondary);margin-bottom:2px;">API 地址</label><input v-model="p.api_base" class="form-input" placeholder="默认官方地址" style="font-size:11px;padding:5px 8px;"></div>
          </div>
          <div v-if="p.tokens_used !== undefined" style="padding:5px 14px;background:var(--color-bg);font-size:10px;color:var(--color-text-secondary);border-top:1px solid var(--color-border);">
            📊 已用 {{ (p.tokens_used||0).toLocaleString() }} tokens · {{ p.total_calls||0 }} 次调用 · 预估 ${{ (p.estimated_cost||0).toFixed(4) }}
          </div>
        </div>
      </div>

      <!-- ===== MCP 自定义接口 ===== -->
      <div v-if="activeTab === 'mcp'" style="max-width:720px;">
        <div class="card" style="padding:20px;margin-bottom:12px;">
          <div style="font-size:14px;font-weight:600;margin-bottom:4px;">🔗 MCP 通用接口</div>
          <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:12px;">
            MCP（Model Context Protocol）通用接口可以连接任意兼容 OpenAI 格式的 API 服务，如自建 vLLM、AstrBot 机器人、本地大模型等。
          </p>
          <div v-if="showAddMcp" style="display:grid;grid-template-columns:1fr 1fr;gap:10px;padding:12px;background:var(--color-bg);border-radius:8px;margin-bottom:12px;">
            <div class="form-group"><label>名称（如：我的AstrBot）</label><input v-model="newMcp.name" class="form-input" :style="{ borderColor: mcpErrors.name ? '#f87171' : '' }" placeholder="自定义名称"></div>
            <div v-if="mcpErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ mcpErrors.name }}</div>
            <div class="form-group"><label>模型</label><select v-model="newMcp.model" class="form-input" :style="{ borderColor: mcpErrors.api_base ? '#f87171' : '' }"><option value="mcp-default">默认</option><option value="gpt-3.5-turbo">GPT-3.5</option><option value="gpt-4o">GPT-4o</option><option value="deepseek-chat">DeepSeek</option><option value="qwen-max">通义千问</option></select></div>
            <div v-if="mcpErrors.api_base" style="color:#f87171;font-size:11px;margin-top:2px;">{{ mcpErrors.api_base }}</div>
            <div class="form-group"><label>API 地址 *</label><input v-model="newMcp.api_base" class="form-input" placeholder="http://你的服务器:8000/v1"></div>
            <div class="form-group"><label>API Key（可选）</label><input v-model="newMcp.api_key" class="form-input" placeholder="如有需要"></div>
            <div style="display:flex;gap:8px;align-items:end;">
              <button class="btn btn-primary btn-sm" @click="addMcp">确认添加</button>
              <button class="btn btn-outline btn-sm" @click="showAddMcp = false; newMcp = { name: '', api_key: '', api_base: '', model: 'mcp-default', is_active: false }">取消</button>
            </div>
          </div>
          <button v-else class="btn btn-outline btn-sm" @click="showAddMcp = true">➕ 添加 MCP 接口</button>
        </div>

        <div v-if="!mcpConfigs.length" style="text-align:center;padding:24px;color:var(--color-text-secondary);font-size:13px;">暂无 MCP 接口配置</div>
        <div v-for="(mcp, idx) in mcpConfigs" :key="mcp.id" class="provider-card">
          <div class="pc-left" style="borderLeftColor:#7c3aed;">
            <div class="pc-header"><strong>{{ mcp.label }}</strong><span :class="['pc-badge', mcp.is_active ? 'on' : 'off']">{{ mcp.is_active ? '启用' : '禁用' }}</span></div>
            <div style="font-size:11px;color:var(--color-text-secondary);">{{ mcp.api_base }}</div>
          </div>
          <div class="pc-right">
            <div style="display:flex;gap:6px;align-items:end;">
              <div class="form-group" style="flex:1;"><label>API Key</label><input v-model="mcp.api_key" type="password" class="form-input" placeholder="可选"></div>
              <button :style="{ padding:'6px 10px',borderRadius:'6px',fontSize:'11px',cursor:'pointer',border:'1px solid var(--color-border)',background:mcp.is_active ? '#10B981' : 'var(--color-bg-card)',color:mcp.is_active ? '#fff' : 'var(--color-text-secondary)' }" @click="mcp.is_active = !mcp.is_active">{{ mcp.is_active ? '启用' : '禁用' }}</button>
              <button style="padding:6px 10px;borderRadius:6px;fontSize:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text-secondary);" @click="removeMcp(idx)">移除</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== Token 用量 ===== -->
      <div v-if="activeTab === 'usage' && usage" class="card" style="max-width:720px;padding:20px;">
        <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:12px;margin-bottom:16px;">
          <div style="padding:12px;background:var(--color-bg);border-radius:8px;text-align:center;">
            <div style="font-size:11px;color:var(--color-text-secondary);">已用 Token</div>
            <div style="font-size:22px;font-weight:700;color:var(--color-text);">{{ (usage.tokens_used || 0).toLocaleString() }}</div>
          </div>
          <div style="padding:12px;background:var(--color-bg);border-radius:8px;text-align:center;">
            <div style="font-size:11px;color:var(--color-text-secondary);">限额</div>
            <div style="font-size:22px;font-weight:700;color:var(--color-text);">{{ (usage.tokens_limit || 0).toLocaleString() }}</div>
          </div>
          <div style="padding:12px;background:var(--color-bg);border-radius:8px;text-align:center;">
            <div style="font-size:11px;color:var(--color-text-secondary);">对话次数</div>
            <div style="font-size:22px;font-weight:700;color:var(--color-text);">{{ usage.total_conversations || 0 }}</div>
          </div>
        </div>

        <!-- 按供应商统计 -->
        <div style="margin-bottom:16px;">
          <div style="font-size:12px;font-weight:600;color:var(--color-text-secondary);margin-bottom:6px;">各供应商用量明细</div>
          <div style="border:1px solid var(--color-border);border-radius:8px;overflow:hidden;">
            <div style="display:grid;grid-template-columns:2fr 1fr 1fr 1fr;gap:8px;padding:6px 12px;background:var(--color-bg);font-size:11px;font-weight:600;color:var(--color-text-secondary);border-bottom:1px solid var(--color-border);">
              <span>供应商</span><span style="text-align:right;">Token</span><span style="text-align:right;">调用</span><span style="text-align:right;">预估费用</span>
            </div>
            <div v-for="(val, key) in usage.by_provider || {}" :key="key" style="display:grid;grid-template-columns:2fr 1fr 1fr 1fr;gap:8px;padding:8px 12px;border-bottom:1px solid var(--color-border);font-size:12px;">
              <span style="font-weight:500;">{{ key }}</span>
              <span style="text-align:right;color:var(--color-primary);font-weight:600;">{{ (val.tokens||0).toLocaleString() }}</span>
              <span style="text-align:right;color:var(--color-text-secondary);">{{ val.total_calls||0 }}</span>
              <span style="text-align:right;font-weight:600;">${{ (val.estimated_cost||0).toFixed(4) }}</span>
            </div>
            <div v-if="!Object.keys(usage.by_provider||{}).length" style="padding:8px 12px;font-size:12px;color:var(--color-text-secondary);text-align:center;">暂无数据</div>
          </div>
        </div>

        <div style="font-size:11px;color:var(--color-text-secondary);background:var(--color-bg);padding:8px 12px;border-radius:6px;margin-bottom:16px;">
          📊 Token 用量来自每次 AI 对话的 API 响应累加。预估费用基于供应商公开定价 × 实际 Token 数计算，仅供参考。实际费用以供应商账单为准。
        </div>

        <div v-if="usage.daily_usage?.length">
          <div style="font-size:12px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">📊 近 7 日趋势</div>
          <div style="border:1px solid var(--color-border);border-radius:8px;overflow:hidden;">
            <div v-for="d in usage.daily_usage" :key="d.date" style="display:flex;align-items:center;gap:8px;padding:6px 12px;border-bottom:1px solid var(--color-border);font-size:11px;">
              <span style="width:48px;flex-shrink:0;color:var(--color-text);">{{ d.date?.slice(5) || d.date }}</span>
              <div style="flex:1;display:flex;flex-direction:column;gap:2px;">
                <div style="display:flex;align-items:center;gap:4px;">
                  <div :style="{ width: Math.min((d.tokens||0) / Math.max(...usage.daily_usage.map(x=>x.tokens||0)) * 100, 100) + '%', height:'12px', background:'linear-gradient(90deg,#7c3aed,#a78bfa)', borderRadius:'4px', minWidth:'4px' }"></div>
                  <span style="color:var(--color-primary);font-weight:600;white-space:nowrap;">{{ (d.tokens||0).toLocaleString() }}</span>
                </div>
                <div style="display:flex;align-items:center;gap:4px;">
                  <div :style="{ width: Math.min((d.count||0) / Math.max(...usage.daily_usage.map(x=>x.count||0)) * 100, 100) + '%', height:'8px', background:'#f59e0b', borderRadius:'4px', minWidth:'4px' }"></div>
                  <span style="color:var(--color-text-secondary);white-space:nowrap;">{{ d.count||0 }} 次</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== 对话记录 ===== -->
      <div v-if="activeTab === 'logs' && usage" class="card" style="max-width:720px;padding:20px;">
        <!-- 筛选 -->
        <div v-if="usage.recent_logs?.length" style="display:flex;gap:8px;margin-bottom:12px;align-items:center;">
          <input v-model="logSearch" class="form-input" placeholder="🔍 搜索学生姓名..." style="max-width:220px;font-size:12px;padding:6px 10px;">
          <span style="font-size:11px;color:var(--color-text-secondary);">共 {{ usage.recent_logs.length }} 条</span>
        </div>
        <div v-if="filteredLogs.length" style="border:1px solid var(--color-border);border-radius:8px;overflow:hidden;">
          <div v-for="log in filteredLogs" :key="log.id" style="padding:10px 14px;border-bottom:1px solid var(--color-border);font-size:12px;">
            <div style="display:flex;gap:8px;margin-bottom:2px;flex-wrap:wrap;"><span style="font-weight:600;flex:1;">{{ log.student_name || '匿名' }}</span><span style="color:var(--color-primary);font-size:11px;">{{ log.tokens_used }} tokens</span><span style="color:var(--color-text-secondary);font-size:11px;">{{ log.created_at }}</span></div>
            <div style="color:var(--color-text);"><strong>问：</strong>{{ log.question }}</div>
            <div style="color:var(--color-text-secondary);"><strong>答：</strong>{{ log.answer?.substring(0,200) }}{{ log.answer?.length > 200 ? '...' : '' }}</div>
          </div>
        </div>
        <div v-else style="padding:24px;text-align:center;font-size:13px;color:var(--color-text-secondary);">{{ usage.recent_logs?.length ? '无匹配记录' : '暂无对话记录' }}</div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.tab-bar { display:flex; gap:4px; margin-bottom:16px; background:var(--color-bg); border-radius:12px; padding:4px; }
.tab-btn { flex:1; padding:8px 10px; border:none; border-radius:10px; font-size:12px; font-weight:600; cursor:pointer; background:transparent; color:var(--color-text-secondary); transition:all 0.2s; white-space:nowrap; }
.tab-btn:hover { background:rgba(124,58,237,0.06); color:var(--color-text); }
.tab-btn.active { background:#7c3aed; color:#fff; box-shadow:0 2px 8px rgba(124,58,237,0.25); }
.form-group { margin-bottom:0; }
.form-group label { display:block; font-size:10px; font-weight:600; color:var(--color-text-secondary); margin-bottom:1px; }
.form-input { color:var(--color-text); width:100%; padding:5px 8px; border:1px solid var(--color-border); border-radius:6px; font-size:12px; outline:none; box-sizing:border-box; background:var(--color-bg-card); }
.form-input:focus { border-color:#7c3aed; }
.btn { padding:6px 14px; border-radius:8px; font-size:12px; font-weight:500; cursor:pointer; border:1px solid transparent; transition:all 0.15s; font-family:inherit; }
.btn-sm { padding:4px 10px; font-size:11px; }
.btn-primary { background:#7c3aed; color:white; border-color:#7c3aed; }
.btn-primary:hover { background:#6d28d9; }
.btn-primary:disabled { opacity:0.5; cursor:not-allowed; }
.btn-outline { background:var(--color-bg-card); color:var(--color-text); border:1px solid var(--color-border); }
</style>
