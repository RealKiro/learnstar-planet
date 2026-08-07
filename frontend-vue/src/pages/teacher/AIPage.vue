<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse } from '@/types'

interface ChatMessage { role: 'user' | 'assistant'; content: string }
interface AICommand { label: string; prompt: string }
interface Usage { configured: boolean; provider: string; model: string }


const loading = ref(true)
const commands = ref<AICommand[]>([])
const usage = ref<Usage | null>(null)
const messages = ref<ChatMessage[]>([])
const input = ref('')
const sending = ref(false)
const chatBody = ref<HTMLElement | null>(null)

const showHint = ref(false)
const loadError = ref('')

async function loadConfig() {
  try {
    const [cmdRes, usageRes] = await Promise.all([
      apiGet<ApiResponse<AICommand[]>>('/api/v1/teacher/ai/commands'),
      apiGet<ApiResponse<Usage>>('/api/v1/teacher/ai/usage'),
    ])
    commands.value = cmdRes.data || []
    usage.value = usageRes.data || null
    loadError.value = ''
  } catch {
    loadError.value = 'AI 配置加载失败'
  }
}

onMounted(async () => {
  loading.value = false
  await loadConfig()
})

async function scrollToBottom() {
  await nextTick()
  if (chatBody.value) chatBody.value.scrollTop = chatBody.value.scrollHeight
}

async function send(prompt?: string) {
  const content = (prompt ?? input.value).trim()
  if (!content || sending.value) return
  messages.value.push({ role: 'user', content })
  input.value = ''
  sending.value = true
  await scrollToBottom()
  try {
    const res = await apiPost<ApiResponse<{ reply: string }>>('/api/v1/teacher/ai/chat', { message: content })
    const reply = (res as unknown as { data: { reply: string } }).data?.reply ?? '暂无回复'
    messages.value.push({ role: 'assistant', content: reply })
    // 检查配置状态
    if (reply.includes('需要配置') || reply.includes('未配置') || reply.includes('未启用')) showHint.value = true
    // 刷新用量
    const usageRes = await apiGet<ApiResponse<Usage>>('/api/v1/teacher/ai/usage')
    usage.value = usageRes.data || null
  } catch { /* handled */ }
  finally {
    sending.value = false
    await scrollToBottom()
  }
}

function useCommand(cmd: AICommand) {
  send(cmd.prompt)
}
</script>

<template>
  <div>
    <div class="page-header">
      <h2 class="page-title">AI 助手</h2>
      <div v-if="usage" class="usage-badge">
        <span :class="usage.configured ? 'usage-config--on' : 'usage-config--off'">{{ usage.configured ? '✅ 已配置' : '❌ 未配置' }}</span>
        <span v-if="usage.configured" class="usage-meta">| {{ usage.provider }} · {{ usage.model }}</span>
      </div>
    </div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else class="card chat-card">
      <!-- 配置加载失败提示（不阻塞聊天） -->
      <div v-if="loadError" class="config-error">
        <span>{{ loadError }}</span>
        <button class="btn btn-sm btn-ghost" @click="loadConfig">重试</button>
      </div>

      <!-- 配置提示 -->
      <div v-if="showHint" class="hint-box">
        <div class="hint-box__title">
          <span>⚠️</span>
          <span>AI 功能尚未配置</span>
        </div>
        <div class="hint-box__body">
          请管理员在后台「AI 中心」添加并启用一个供应商（支持 30+ 供应商），即可开始使用。
        </div>
      </div>

      <!-- 消息列表 -->
      <div ref="chatBody" class="chat-body">
        <div v-if="messages.length === 0" class="chat-empty">
          <div class="chat-empty__icon">🤖</div>
          <p class="chat-empty__main">向 AI 助手提问，获取教学建议</p>
          <p class="chat-empty__sub">例如：「如何提高学生课堂参与度？」</p>
        </div>
        <div v-for="(msg, i) in messages" :key="i"
          class="msg-row"
          :class="msg.role === 'user' ? 'msg-row--user' : 'msg-row--assistant'">
          <div class="msg-bubble" :class="msg.role === 'user' ? 'msg-bubble--user' : 'msg-bubble--assistant'">{{ msg.content }}</div>
        </div>
        <div v-if="sending" class="typing">
          AI 正在思考...
        </div>
      </div>

      <!-- 预设命令 -->
      <div v-if="commands.length" class="commands-row">
        <button v-for="(cmd, i) in commands" :key="i" class="btn btn-sm btn-ghost cmd-btn"
          :disabled="sending" @click="useCommand(cmd)">
          {{ cmd.label }}
        </button>
      </div>

      <!-- 输入框 -->
      <div class="input-row">
        <input v-model="input" class="form-input" placeholder="输入消息..."
          @keydown.enter="send()" :disabled="sending">
        <button class="btn btn-primary send-btn" :disabled="sending || !input.trim()" @click="send()">
          {{ sending ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-title { font-size: 24px; font-weight: 700; }
.usage-badge { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--color-text-secondary); background: var(--color-bg); padding: 6px 12px; border-radius: var(--radius-sm); }
.usage-config--on { color: var(--color-accent); }
.usage-config--off { color: #EF4444; }
.usage-meta { color: var(--color-text-secondary); }

/* ===== 聊天容器 ===== */
.chat-card { display: flex; flex-direction: column; height: calc(100vh - 200px); min-height: 400px; }

/* ===== 配置加载失败（非阻塞，不挡聊天） ===== */
.config-error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--color-danger-text);
}

/* ===== 配置提示 ===== */
.hint-box {
  background: rgba(245,158,11,0.1);
  border: 1px solid rgba(245,158,11,0.3);
  border-radius: var(--radius-sm);
  padding: 10px 14px;
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--md-gold);
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.hint-box__title { display: flex; align-items: center; gap: 8px; }
.hint-box__body { font-size: 12px; padding-left: 24px; color: var(--color-warning-text); line-height: 1.6; }

/* ===== 消息区 ===== */
.chat-body { flex: 1; overflow-y: auto; padding: 8px 4px; display: flex; flex-direction: column; gap: 12px; }
.chat-empty { text-align: center; padding: 48px; color: var(--color-text-secondary); }
.chat-empty__icon { font-size: 48px; margin-bottom: 12px; }
.chat-empty__main { margin-bottom: 8px; }
.chat-empty__sub { font-size: 13px; }
.msg-row { max-width: 75%; }
.msg-row--user { align-self: flex-end; }
.msg-row--assistant { align-self: flex-start; }
.msg-bubble { padding: 10px 14px; border-radius: var(--radius-sm); font-size: 14px; line-height: 1.6; white-space: pre-wrap; }
.msg-bubble--user { background: var(--color-primary); color: #fff; }
.msg-bubble--assistant { background: var(--color-bg); color: var(--color-text); }
.typing { align-self: flex-start; padding: 10px 14px; background: var(--color-bg); border-radius: var(--radius-sm); font-size: 14px; color: var(--color-text-secondary); }

/* ===== 预设命令 & 输入 ===== */
.commands-row { display: flex; gap: 8px; flex-wrap: wrap; padding: 12px 0; border-top: 1px solid var(--color-border); }
.cmd-btn { font-size: 12px; }
.input-row { display: flex; gap: 8px; padding-top: 8px; border-top: 1px solid var(--color-border); }
.send-btn { width: auto; }
</style>
