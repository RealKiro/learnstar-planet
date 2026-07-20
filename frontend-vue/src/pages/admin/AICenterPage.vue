<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()

interface AiSettings {
  enabled: boolean; provider: string; api_key: string
  api_key_masked: boolean; api_base: string; model: string
  max_tokens: number; tokens_used: number; tokens_limit: number
}
interface DailyUsage { date: string; tokens: number; count: number }
interface ConversationLog {
  id: number; student_name: string; class_id: number | null
  question: string; answer: string; tokens_used: number; created_at: string
}
interface AiUsage {
  enabled: boolean; tokens_used: number; tokens_limit: number
  total_conversations: number; daily_usage: DailyUsage[]
  recent_logs: ConversationLog[]
}

const loading = ref(true)
const settings = ref<AiSettings | null>(null)
const usage = ref<AiUsage | null>(null)
const saving = ref(false)
const activeTab = ref<'settings' | 'usage' | 'logs'>('settings')
const showApiKey = ref(false)

const providers = [
  { id: 'openai', label: 'OpenAI', models: ['gpt-3.5-turbo', 'gpt-4', 'gpt-4-turbo', 'gpt-4o'] },
  { id: 'claude', label: 'Anthropic Claude', models: ['claude-3-haiku', 'claude-3-sonnet', 'claude-3-opus', 'claude-3-5-sonnet'] },
  { id: 'qwen', label: '通义千问（阿里）', models: ['qwen-turbo', 'qwen-plus', 'qwen-max', 'qwen2-72b-instruct'] },
  { id: 'deepseek', label: 'DeepSeek（深度求索）', models: ['deepseek-chat', 'deepseek-coder'] },
  { id: 'ernie', label: '文心一言（百度）', models: ['ernie-4.0', 'ernie-3.5', 'ernie-speed'] },
  { id: 'hunyuan', label: '混元（腾讯）', models: ['hunyuan-lite', 'hunyuan-standard', 'hunyuan-pro'] },
  { id: 'glm', label: 'GLM（智谱）', models: ['glm-4', 'glm-4-plus', 'glm-4-air'] },
  { id: 'spark', label: '星火（讯飞）', models: ['spark-3.0', 'spark-4.0'] },
  { id: 'mcp', label: 'MCP 通用接口', models: ['mcp-default'] },
]

const activeProvider = computed(() => providers.find(p => p.id === settings.value?.provider))

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
  try {
    const res = await fetch('/api/v1/admin/ai/toggle', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: val }),
    })
    const data = await res.json()
    toast.show(data.message || '操作成功', res.ok ? 'success' : 'error')
    if (res.ok && settings.value) settings.value.enabled = val
  } catch { toast.show('操作失败', 'error') }
}

function maskedKey(key: string): string {
  if (!key || key.length < 12) return key
  return key.substring(0, 8) + '••••' + key.substring(key.length - 4)
}

onMounted(loadData)
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">系统配置</p>
        <h2 style="font-size:24px;font-weight:700;">🤖 AI 中心</h2>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-else>
      <!-- 状态开关 -->
      <div class="card" style="max-width:640px;padding:24px;margin-bottom:20px;display:flex;align-items:center;justify-content:space-between;">
        <div>
          <div style="font-size:15px;font-weight:600;color:var(--color-text);">AI 功能开关</div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">
            {{ settings?.enabled ? '开启后班级码大屏端显示 AI 对话入口，学生可向 AI 提问' : '关闭后班级码大屏端隐藏 AI 功能' }}
          </div>
        </div>
        <label style="position:relative;display:inline-block;width:48px;height:26px;cursor:pointer;">
          <input type="checkbox" :checked="settings?.enabled" @change="toggleAi(($event.target as HTMLInputElement).checked)" style="opacity:0;width:0;height:0;">
          <span :style="{
            position:'absolute',inset:0,background:settings?.enabled ? '#7c3aed' : '#ccc',borderRadius:'13px',transition:'0.2s',
          }">
            <span :style="{
              position:'absolute',top:'3px',left:settings?.enabled ? '25px' : '3px',width:'20px',height:'20px',
              borderRadius:'50%',background:'#fff',transition:'0.2s',boxShadow:'0 1px 3px rgba(0,0,0,0.2)',
            }"></span>
          </span>
        </label>
      </div>

      <!-- 标签切换 -->
      <div class="tab-bar" style="max-width:640px;">
        <button :class="['tab-btn', { active: activeTab === 'settings' }]" @click="activeTab = 'settings'">⚙️ 模型配置</button>
        <button :class="['tab-btn', { active: activeTab === 'usage' }]" @click="activeTab = 'usage'">📊 Token 用量</button>
        <button :class="['tab-btn', { active: activeTab === 'logs' }]" @click="activeTab = 'logs'">📋 对话记录</button>
      </div>

      <!-- 模型配置 -->
      <div v-if="activeTab === 'settings' && settings" class="card" style="max-width:640px;padding:24px;">
        <div class="form-group">
          <label>AI 提供商</label>
          <select v-model="settings.provider" class="form-input">
            <option v-for="p in providers" :key="p.id" :value="p.id">{{ p.label }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>模型</label>
          <select v-model="settings.model" class="form-input">
            <option v-for="m in activeProvider?.models || []" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>API Key</label>
          <div style="display:flex;gap:8px;">
            <input :type="showApiKey ? 'text' : 'password'" v-model="settings.api_key" class="form-input" placeholder="输入 API Key">
            <button class="btn btn-sm" style="flex-shrink:0;" @click="showApiKey = !showApiKey">{{ showApiKey ? '🙈' : '👁️' }}</button>
          </div>
          <div v-if="settings.api_key_masked && !settings.api_key" style="font-size:11px;color:var(--color-text-secondary);margin-top:2px;">已配置 API Key（重新输入会覆盖）</div>
        </div>
        <div class="form-group">
          <label>API 地址（可选）</label>
          <input v-model="settings.api_base" class="form-input" placeholder="留空使用官方默认地址">
        </div>
        <div style="display:flex;gap:12px;">
          <div class="form-group" style="flex:1;">
            <label>单次最大 Token</label>
            <input v-model.number="settings.max_tokens" type="number" class="form-input" min="100" max="32000">
          </div>
          <div class="form-group" style="flex:1;">
            <label>总 Token 限额</label>
            <input v-model.number="settings.tokens_limit" type="number" class="form-input" min="0">
          </div>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:16px;">
          <button class="btn btn-primary" :disabled="saving" @click="saveSettings">{{ saving ? '保存中...' : '保存设置' }}</button>
        </div>
      </div>

      <!-- Token 用量 -->
      <div v-if="activeTab === 'usage' && usage" class="card" style="max-width:640px;padding:24px;">
        <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:12px;margin-bottom:20px;">
          <div class="stat-card">
            <div class="stat-label">已用 Token</div>
            <div class="stat-value">{{ (usage.tokens_used || 0).toLocaleString() }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">限额</div>
            <div class="stat-value">{{ (usage.tokens_limit || 0).toLocaleString() }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">对话次数</div>
            <div class="stat-value">{{ usage.total_conversations || 0 }}</div>
          </div>
        </div>
        <div style="margin-bottom:20px;">
          <div style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">使用率</div>
          <div style="height:8px;background:var(--color-bg);border-radius:4px;overflow:hidden;">
            <div :style="{
              height:'100%',borderRadius:'4px',transition:'width 0.5s',
              background: (usage.tokens_used / (usage.tokens_limit || 1)) > 0.8 ? '#EF4444' : '#10B981',
              width: Math.min(100, (usage.tokens_used / (usage.tokens_limit || 1)) * 100) + '%',
            }"></div>
          </div>
        </div>
        <div v-if="usage.daily_usage && usage.daily_usage.length > 0">
          <div style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">近 7 日趋势</div>
          <div class="daily-list">
            <div v-for="d in usage.daily_usage" :key="d.date" class="daily-row">
              <span class="daily-date">{{ d.date }}</span>
              <span class="daily-tokens">{{ (d.tokens || 0).toLocaleString() }} tokens</span>
              <span class="daily-count">{{ d.count || 0 }} 次对话</span>
            </div>
          </div>
        </div>
        <div v-else style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);">暂无使用数据</div>
      </div>

      <!-- 对话记录 -->
      <div v-if="activeTab === 'logs' && usage" class="card" style="max-width:640px;padding:24px;">
        <div v-if="usage.recent_logs && usage.recent_logs.length > 0" class="log-list">
          <div v-for="log in usage.recent_logs" :key="log.id" class="log-item">
            <div class="log-header">
              <span class="log-student">{{ log.student_name || '匿名' }}</span>
              <span class="log-tokens">{{ log.tokens_used }} tokens</span>
              <span class="log-time">{{ log.created_at }}</span>
            </div>
            <div class="log-q"><strong>问：</strong>{{ log.question }}</div>
            <div v-if="log.answer" class="log-a"><strong>答：</strong>{{ log.answer.substring(0, 200) }}{{ log.answer.length > 200 ? '...' : '' }}</div>
          </div>
        </div>
        <div v-else style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);">暂无对话记录</div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }
.stat-card { padding: 12px; background: var(--color-bg); border-radius: 8px; text-align: center; }
.stat-label { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 4px; }
.stat-value { font-size: 20px; font-weight: 700; color: var(--color-text); }
.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 12px; font-weight: 600; color: var(--color-text); margin-bottom: 4px; }
.form-input { color: var(--color-text); width: 100%; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 13px; outline: none; transition: border-color 0.15s; box-sizing: border-box; background: var(--color-bg-card); }
.form-input:focus { border-color: #7c3aed; box-shadow: 0 0 0 3px rgba(124,58,237,0.08); }
.daily-list { border: 1px solid var(--color-border); border-radius: 8px; overflow: hidden; }
.daily-row { display: flex; gap: 12px; padding: 8px 12px; border-bottom: 1px solid var(--color-border); font-size: 12px; }
.daily-row:last-child { border-bottom: none; }
.daily-date { flex: 1; font-weight: 500; }
.daily-tokens { color: var(--color-primary); font-weight: 600; }
.daily-count { color: var(--color-text-secondary); }
.log-list { border: 1px solid var(--color-border); border-radius: 8px; overflow: hidden; }
.log-item { padding: 10px 14px; border-bottom: 1px solid var(--color-border); font-size: 12px; }
.log-item:last-child { border-bottom: none; }
.log-header { display: flex; gap: 8px; align-items: center; margin-bottom: 4px; }
.log-student { font-weight: 600; color: var(--color-text); flex: 1; }
.log-tokens { color: var(--color-primary); font-size: 11px; }
.log-time { color: var(--color-text-secondary); font-size: 11px; }
.log-q { color: var(--color-text); margin-bottom: 2px; }
.log-a { color: var(--color-text-secondary); }
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; }
.btn-sm { padding: 6px 12px; font-size: 12px; }
.btn-primary { background: #7c3aed; color: white; border-color: #7c3aed; }
.btn-primary:hover { background: #6d28d9; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
