<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()

interface ProviderConfig {
  id: string; label: string; api_key: string; api_base: string; model: string; is_active: boolean
}
interface AiSettings { enabled: boolean; max_tokens: number; tokens_used: number; tokens_limit: number; providers: ProviderConfig[] }
interface DailyUsage { date: string; tokens: number; count: number }
interface ConversationLog { id: number; student_name: string; question: string; answer: string; tokens_used: number; created_at: string }
interface AiUsage { enabled: boolean; tokens_used: number; tokens_limit: number; total_conversations: number; daily_usage: DailyUsage[]; recent_logs: ConversationLog[] }

const loading = ref(true)
const settings = ref<AiSettings | null>(null)
const usage = ref<AiUsage | null>(null)
const saving = ref(false)
const activeTab = ref<'providers' | 'usage' | 'logs'>('providers')

const providerMeta = [
  { id: 'openai', label: 'OpenAI', color: '#10a37f', models: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-3.5-turbo'] },
  { id: 'claude', label: 'Anthropic Claude', color: '#d97706', models: ['claude-3-5-sonnet', 'claude-3-haiku', 'claude-3-opus'] },
  { id: 'google', label: 'Google Gemini', color: '#4285f4', models: ['gemini-2.0-flash', 'gemini-1.5-pro', 'gemini-1.5-flash'] },
  { id: 'grok', label: 'xAI Grok', color: '#1a1a2e', models: ['grok-2', 'grok-2-mini'] },
  { id: 'mistral', label: 'Mistral AI', color: '#ff6f00', models: ['mistral-large', 'mistral-small', 'codestral'] },
  { id: 'qwen', label: '通义千问', color: '#1677ff', models: ['qwen-max', 'qwen-plus', 'qwen-turbo', 'qwen2.5-72b'] },
  { id: 'deepseek', label: 'DeepSeek', color: '#4f6ef7', models: ['deepseek-chat', 'deepseek-v3', 'deepseek-r1', 'deepseek-coder'] },
  { id: 'moonshot', label: '月之暗面 Kimi', color: '#8b5cf6', models: ['moonshot-v1-128k', 'moonshot-v1-32k', 'moonshot-v1-8k'] },
  { id: 'bytedance', label: '豆包（字节）', color: '#00a76a', models: ['doubao-1.5-pro', 'doubao-1.5-lite', 'doubao-pro'] },
  { id: 'minimax', label: 'MiniMax', color: '#ff6900', models: ['minimax-text-01', 'minimax-abab-6.5'] },
  { id: 'baichuan', label: '百川智能', color: '#2b6cb0', models: ['baichuan4', 'baichuan3-turbo'] },
  { id: 'stepfun', label: '阶跃星辰', color: '#e00', models: ['step-2-16k', 'step-1-32k'] },
  { id: 'lingyi', label: '零一万物 Yi', color: '#06b6d4', models: ['yi-large', 'yi-medium', 'yi-spark'] },
  { id: 'ernie', label: '文心一言', color: '#3060b0', models: ['ernie-4.0', 'ernie-3.5', 'ernie-speed'] },
  { id: 'hunyuan', label: '混元（腾讯）', color: '#0052d9', models: ['hunyuan-pro', 'hunyuan-standard', 'hunyuan-lite'] },
  { id: 'glm', label: 'GLM（智谱）', color: '#8b4513', models: ['glm-4', 'glm-4-plus', 'glm-4-air'] },
  { id: 'spark', label: '星火（讯飞）', color: '#e8422a', models: ['spark-4.0', 'spark-3.0'] },
  { id: 'cohere', label: 'Cohere', color: '#d18ee2', models: ['command-r-plus', 'command-r'] },
  { id: 'perplexity', label: 'Perplexity', color: '#1f1f1f', models: ['sonar-pro', 'sonar-reasoning'] },
  { id: 'siliconflow', label: '硅基流动', color: '#409eff', models: ['Pro/DeepSeek-V3', 'Pro/Qwen2.5-72B'] },
  { id: 'openrouter', label: 'OpenRouter', color: '#6466f1', models: ['openrouter/auto', 'anthropic/claude-3.5-sonnet', 'openai/gpt-4o'] },
  { id: 'nvidia', label: 'NVIDIA NIM', color: '#76b900', models: ['nvidia/llama-3.1-nemotron', 'meta/llama-3.1-8b'] },
  { id: 'groq', label: 'Groq LPU', color: '#f55036', models: ['llama-3.3-70b-versatile', 'mixtral-8x7b-32768'] },
  { id: 'together', label: 'Together AI', color: '#7678ff', models: ['meta-llama/Llama-3.3-70B', 'mistralai/Mixtral-8x7B'] },
  { id: 'azure', label: 'Azure OpenAI', color: '#0078d4', models: ['gpt-4o', 'gpt-4-turbo', 'gpt-35-turbo'] },
  { id: 'ollama', label: 'Ollama（本地）', color: '#000', models: ['llama3.2', 'qwen2.5', 'deepseek-r1'] },
  { id: 'mcp', label: 'MCP 通用接口', color: '#7c3aed', models: ['mcp-default'] },
]

const showAddProvider = ref(false)
const newProvider = ref({ id: '', label: '', api_key: '', api_base: '', model: '', is_active: false })
const availableProviders = computed(() => providerMeta.filter(m => !settings.value?.providers?.some(p => p.id === m.id)))

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
  saving.value = true
  try {
    const res = await fetch('/api/v1/admin/ai/settings', {
      method: 'PUT',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify(settings.value),
    })
    const data = await res.json()
    toast.show(data.message || '设置已保存', res.ok ? 'success' : 'error')
    if (res.ok) loadData()
  } catch { toast.show('保存失败', 'error') }
  finally { saving.value = false }
}

async function toggleAi(val: boolean) {
  if (!settings.value) return
  try {
    const res = await fetch('/api/v1/admin/ai/toggle', {
      method: 'POST', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: val }),
    })
    const data = await res.json()
    toast.show(data.message || '操作成功', res.ok ? 'success' : 'error')
    if (res.ok) settings.value.enabled = val
  } catch { toast.show('操作失败', 'error') }
}

function addProvider() {
  if (!newProvider.value.id) { toast.show('请选择供应商', 'error'); return }
  const meta = providerMeta.find(m => m.id === newProvider.value.id)
  if (!settings.value?.providers) settings.value.providers = []
  settings.value.providers.push({ ...newProvider.value, label: meta?.label || newProvider.value.id })
  showAddProvider.value = false
  newProvider.value = { id: '', label: '', api_key: '', api_base: '', model: '', is_active: false }
}

function removeProvider(idx: number) {
  if (!settings.value?.providers) return
  settings.value.providers.splice(idx, 1)
}

function getMeta(id: string) { return providerMeta.find(m => m.id === id) }

onMounted(loadData)
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div><p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">系统配置</p><h2 style="font-size:24px;font-weight:700;">🤖 AI 中心</h2></div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>
    <template v-else>
      <!-- 开关 -->
      <div class="card" style="max-width:720px;padding:20px;margin-bottom:16px;display:flex;align-items:center;justify-content:space-between;">
        <div>
          <div style="font-size:15px;font-weight:600;color:var(--color-text);">AI 总开关</div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:2px;">{{ settings?.enabled ? '开启后班级码大屏显示 AI 入口' : '关闭后隐藏 AI 功能' }}</div>
        </div>
        <label style="position:relative;display:inline-block;width:48px;height:26px;cursor:pointer;">
          <input type="checkbox" :checked="settings?.enabled" @change="toggleAi(($event.target as HTMLInputElement).checked)" style="opacity:0;width:0;height:0;">
          <span :style="{ position:'absolute',inset:0,background:settings?.enabled ? '#7c3aed' : '#ccc',borderRadius:'13px',transition:'0.2s' }">
            <span :style="{ position:'absolute',top:'3px',left:settings?.enabled ? '25px' : '3px',width:'20px',height:'20px',borderRadius:'50%',background:'#fff',transition:'0.2s',boxShadow:'0 1px 3px rgba(0,0,0,0.2)' }"></span>
          </span>
        </label>
      </div>

      <!-- 标签 -->
      <div class="tab-bar" style="max-width:720px;">
        <button :class="['tab-btn', { active: activeTab === 'providers' }]" @click="activeTab = 'providers'">🔌 供应商配置</button>
        <button :class="['tab-btn', { active: activeTab === 'usage' }]" @click="activeTab = 'usage'">📊 Token 用量</button>
        <button :class="['tab-btn', { active: activeTab === 'logs' }]" @click="activeTab = 'logs'">📋 对话记录</button>
      </div>

      <!-- 供应商管理 -->
      <div v-if="activeTab === 'providers'" class="card" style="max-width:720px;padding:24px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;">
          <div style="font-size:15px;font-weight:600;color:var(--color-text);">已配置供应商</div>
          <button v-if="!showAddProvider" class="btn btn-sm btn-primary" @click="showAddProvider = true">+ 添加供应商</button>
        </div>

        <!-- 添加新供应商 -->
        <div v-if="showAddProvider" style="margin-bottom:16px;padding:16px;background:var(--color-bg);border-radius:10px;">
          <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;">
            <div class="form-group"><label>供应商</label>
              <select v-model="newProvider.id" class="form-input" @change="newProvider.model = getMeta(newProvider.id)?.models?.[0] || ''">
                <option value="">请选择</option>
                <optgroup v-for="g in ([...new Set(providerMeta.map(m => m.id.includes('mcp') ? '🔌 通用' : m.id.includes('ollama') ? '🏠 本地' : '🌐 其他'))])" :key="g" :label="g">
                  <option v-for="m in providerMeta.filter(p => (p.id.includes('mcp') ? '🔌 通用' : p.id.includes('ollama') ? '🏠 本地' : '🌐 其他') === g)" :key="m.id" :value="m.id">{{ m.label }}</option>
                </optgroup>
              </select>
            </div>
            <div class="form-group"><label>模型</label><select v-model="newProvider.model" class="form-input"><option value="">默认模型</option><option v-for="m in getMeta(newProvider.id)?.models || []" :key="m" :value="m">{{ m }}</option></select></div>
            <div class="form-group"><label>API Key</label><input v-model="newProvider.api_key" type="password" class="form-input" placeholder="sk-..."></div>
            <div class="form-group"><label>API 地址（可选）</label><input v-model="newProvider.api_base" class="form-input" placeholder="https://api.openai.com/v1"></div>
          </div>
          <div style="display:flex;gap:8px;margin-top:12px;">
            <button class="btn btn-primary" @click="addProvider">确认添加</button>
            <button class="btn btn-outline" @click="showAddProvider = false">取消</button>
          </div>
        </div>

        <!-- 已配置列表 -->
        <div v-if="!settings?.providers?.length && !showAddProvider" style="text-align:center;padding:24px;color:var(--color-text-secondary);font-size:13px;">暂无供应商，点击上方添加</div>
        <div v-for="(p, i) in settings?.providers || []" :key="p.id" style="display:flex;align-items:center;gap:12px;padding:12px;margin-bottom:8px;background:var(--color-bg);border-radius:10px;border:1px solid var(--color-border);">
          <div :style="{ width:'4px',height:'32px',borderRadius:'2px',background:p.is_active ? '#10B981' : '#ccc',flexShrink:0 }"></div>
          <div style="flex:1;min-width:0;">
            <div style="display:flex;align-items:center;gap:6px;">
              <span style="font-weight:600;font-size:14px;">{{ getMeta(p.id)?.label || p.label || p.id }}</span>
              <span v-if="p.is_active" style="font-size:10px;padding:1px 6px;border-radius:4px;background:rgba(16,185,129,0.1);color:#10B981;">启用</span>
              <span v-else style="font-size:10px;padding:1px 6px;border-radius:4px;background:rgba(156,163,175,0.1);color:#9CA3AF;">禁用</span>
            </div>
            <div style="font-size:11px;color:var(--color-text-secondary);margin-top:2px;">模型: {{ p.model || '默认' }} · {{ p.api_key ? '🔑 已配置' : '⚠️ 未配置 Key' }}</div>
          </div>
          <div style="display:flex;gap:4px;flex-shrink:0;">
            <button :style="{ padding:'4px 10px',borderRadius:'6px',fontSize:'11px',cursor:'pointer',border:'1px solid var(--color-border)',background:p.is_active ? 'var(--color-bg-card)' : '#10B981',color:p.is_active ? 'var(--color-text-secondary)' : '#fff' }" @click="p.is_active = !p.is_active">{{ p.is_active ? '禁用' : '启用' }}</button>
            <button style="padding:4px 10px;borderRadius:6px;fontSize:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text-secondary);" @click="removeProvider(i)">🗑️</button>
          </div>
        </div>

        <div style="display:flex;gap:12px;margin-top:16px;padding-top:16px;border-top:1px solid var(--color-border);">
          <div class="form-group" style="flex:1;"><label>单次最大 Token</label><input v-model.number="settings!.max_tokens" type="number" class="form-input" min="100" max="32000"></div>
          <div class="form-group" style="flex:1;"><label>总 Token 限额</label><input v-model.number="settings!.tokens_limit" type="number" class="form-input" min="0" placeholder="0=不限制"></div>
        </div>
        <div style="display:flex;justify-content:flex-end;margin-top:12px;">
          <button class="btn btn-primary" :disabled="saving" @click="saveSettings">{{ saving ? '保存中...' : '保存设置' }}</button>
        </div>
      </div>

      <!-- Token 用量 -->
      <div v-if="activeTab === 'usage' && usage" class="card" style="max-width:720px;padding:24px;">
        <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:12px;margin-bottom:20px;">
          <div style="padding:12px;background:var(--color-bg);border-radius:8px;text-align:center;"><div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:4px;">已用 Token</div><div style="font-size:20px;font-weight:700;">{{ (usage.tokens_used || 0).toLocaleString() }}</div></div>
          <div style="padding:12px;background:var(--color-bg);border-radius:8px;text-align:center;"><div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:4px;">限额</div><div style="font-size:20px;font-weight:700;">{{ (usage.tokens_limit || 0).toLocaleString() }}</div></div>
          <div style="padding:12px;background:var(--color-bg);border-radius:8px;text-align:center;"><div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:4px;">对话次数</div><div style="font-size:20px;font-weight:700;">{{ usage.total_conversations || 0 }}</div></div>
        </div>
        <div v-if="usage.daily_usage?.length" style="margin-bottom:16px;">
          <div style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">近 7 日趋势</div>
          <div style="border:1px solid var(--color-border);border-radius:8px;overflow:hidden;">
            <div v-for="d in usage.daily_usage" :key="d.date" style="display:flex;gap:12px;padding:6px 12px;border-bottom:1px solid var(--color-border);font-size:12px;">
              <span style="flex:1;font-weight:500;">{{ d.date }}</span><span style="color:var(--color-primary);font-weight:600;">{{ (d.tokens||0).toLocaleString() }} tokens</span><span style="color:var(--color-text-secondary);">{{ d.count||0 }} 次</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 对话记录 -->
      <div v-if="activeTab === 'logs' && usage" class="card" style="max-width:720px;padding:24px;">
        <div v-if="usage.recent_logs?.length" style="border:1px solid var(--color-border);border-radius:8px;overflow:hidden;">
          <div v-for="log in usage.recent_logs" :key="log.id" style="padding:10px 14px;border-bottom:1px solid var(--color-border);font-size:12px;">
            <div style="display:flex;gap:8px;margin-bottom:4px;"><span style="font-weight:600;flex:1;">{{ log.student_name || '匿名' }}</span><span style="color:var(--color-primary);font-size:11px;">{{ log.tokens_used }} tokens</span><span style="color:var(--color-text-secondary);font-size:11px;">{{ log.created_at }}</span></div>
            <div style="color:var(--color-text);margin-bottom:2px;"><strong>问：</strong>{{ log.question }}</div>
            <div style="color:var(--color-text-secondary);"><strong>答：</strong>{{ log.answer?.substring(0,200) }}{{ log.answer?.length > 200 ? '...' : '' }}</div>
          </div>
        </div>
        <div v-else style="padding:24px;text-align:center;font-size:13px;color:var(--color-text-secondary);">暂无对话记录</div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.tab-bar { display:flex; gap:4px; margin-bottom:20px; background:var(--color-bg); border-radius:12px; padding:4px; }
.tab-btn { flex:1; padding:10px 12px; border:none; border-radius:10px; font-size:13px; font-weight:600; cursor:pointer; background:transparent; color:var(--color-text-secondary); transition:all 0.2s; }
.tab-btn:hover { background:rgba(124,58,237,0.06); color:var(--color-text); }
.tab-btn.active { background:#7c3aed; color:#fff; box-shadow:0 2px 8px rgba(124,58,237,0.25); }
.form-group { margin-bottom:0; }
.form-group label { display:block; font-size:11px; font-weight:600; color:var(--color-text-secondary); margin-bottom:2px; }
.form-input { color:var(--color-text); width:100%; padding:6px 10px; border:1px solid var(--color-border); border-radius:6px; font-size:12px; outline:none; box-sizing:border-box; background:var(--color-bg-card); }
.form-input:focus { border-color:#7c3aed; }
.btn { padding:6px 14px; border-radius:8px; font-size:12px; font-weight:500; cursor:pointer; border:1px solid transparent; transition:all 0.15s; }
.btn-sm { padding:4px 12px; font-size:11px; }
.btn-primary { background:#7c3aed; color:white; border-color:#7c3aed; }
.btn-primary:hover { background:#6d28d9; }
.btn-primary:disabled { opacity:0.5; cursor:not-allowed; }
.btn-outline { background:var(--color-bg-card); color:var(--color-text); border:1px solid var(--color-border); }
</style>
