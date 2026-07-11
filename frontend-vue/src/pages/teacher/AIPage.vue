<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse } from '@/types'

interface ChatMessage { role: 'user' | 'assistant'; content: string }
interface AICommand { label: string; prompt: string }
interface Usage { count: number; limit: number; reset_at?: string }


const loading = ref(true)
const commands = ref<AICommand[]>([])
const usage = ref<Usage | null>(null)
const messages = ref<ChatMessage[]>([])
const input = ref('')
const sending = ref(false)
const chatBody = ref<HTMLElement | null>(null)

const showHint = ref(false)

onMounted(async () => {
  try {
    const [cmdRes, usageRes] = await Promise.all([
      apiGet<ApiResponse<AICommand[]>>('/api/v1/teacher/ai/commands'),
      apiGet<ApiResponse<Usage>>('/api/v1/teacher/ai/usage'),
    ])
    commands.value = cmdRes.data || []
    usage.value = usageRes.data || null
  } catch { /* handled */ }
  finally { loading.value = false }
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
    if (reply.includes('需要配置') || reply.includes('未配置')) showHint.value = true
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
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">AI 助手</h2>
      <div v-if="usage" style="display:flex;align-items:center;gap:8px;font-size:13px;color:var(--color-text-secondary);background:var(--color-bg);padding:6px 12px;border-radius:var(--radius-sm);">
        <span>用量</span>
        <span style="font-weight:700;color:var(--color-primary);">{{ usage.count }}/{{ usage.limit }}</span>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else class="card" style="display:flex;flex-direction:column;height:calc(100vh - 200px);min-height:400px;">
      <!-- 配置提示 -->
      <div v-if="showHint" style="background:rgba(245,158,11,0.1);border:1px solid rgba(245,158,11,0.3);border-radius:var(--radius-sm);padding:10px 14px;margin-bottom:12px;font-size:13px;color:#92400e;display:flex;align-items:center;gap:8px;">
        <span>⚠️</span>
        <span>AI 功能尚未配置，当前为占位回复。请联系管理员配置 AI 服务。</span>
      </div>

      <!-- 消息列表 -->
      <div ref="chatBody" style="flex:1;overflow-y:auto;padding:8px 4px;display:flex;flex-direction:column;gap:12px;">
        <div v-if="messages.length === 0" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
          <div style="font-size:48px;margin-bottom:12px;">🤖</div>
          <p style="margin-bottom:8px;">向 AI 助手提问，获取教学建议</p>
          <p style="font-size:13px;">例如：「如何提高学生课堂参与度？」</p>
        </div>
        <div v-for="(msg, i) in messages" :key="i"
          :style="{
            alignSelf: msg.role === 'user' ? 'flex-end' : 'flex-start',
            maxWidth: '75%',
          }">
          <div :style="{
            padding: '10px 14px',
            borderRadius: 'var(--radius-sm)',
            fontSize: '14px',
            lineHeight: '1.6',
            background: msg.role === 'user' ? 'var(--color-primary)' : 'var(--color-bg)',
            color: msg.role === 'user' ? 'white' : 'var(--color-text)',
            whiteSpace: 'pre-wrap',
          }">{{ msg.content }}</div>
        </div>
        <div v-if="sending" style="align-self:flex-start;padding:10px 14px;background:var(--color-bg);border-radius:var(--radius-sm);font-size:14px;color:var(--color-text-secondary);">
          AI 正在思考...
        </div>
      </div>

      <!-- 预设命令 -->
      <div v-if="commands.length" style="display:flex;gap:8px;flex-wrap:wrap;padding:12px 0;border-top:1px solid var(--color-border);">
        <button v-for="(cmd, i) in commands" :key="i" class="btn btn-sm btn-ghost"
          style="font-size:12px;" :disabled="sending" @click="useCommand(cmd)">
          {{ cmd.label }}
        </button>
      </div>

      <!-- 输入框 -->
      <div style="display:flex;gap:8px;padding-top:8px;border-top:1px solid var(--color-border);">
        <input v-model="input" class="form-input" placeholder="输入消息..."
          @keydown.enter="send()" :disabled="sending">
        <button class="btn btn-primary" style="width:auto;" :disabled="sending || !input.trim()" @click="send()">
          {{ sending ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>
