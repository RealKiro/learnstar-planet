<script setup lang="ts">
import { ref, onMounted } from 'vue'

const messages = ref<{ role: string; content: string }[]>([])
const input = ref('')
const sending = ref(false)
const aiEnabled = ref(false)
const loading = ref(true)
const chatBody = ref<HTMLElement | null>(null)

onMounted(async () => {
  const token = sessionStorage.getItem('class_token') || ''
  try {
    const res = await fetch('/api/v1/display/ai/settings', {
      headers: { 'Authorization': 'Bearer ' + token }
    })
    const data = await res.json()
    aiEnabled.value = data?.data?.enabled || false
  } catch { aiEnabled.value = false }
  finally { loading.value = false }
})

async function sendMsg() {
  const text = input.value.trim()
  if (!text || sending.value) return
  messages.value.push({ role: 'user', content: text })
  input.value = ''
  sending.value = true
  try {
    const token = sessionStorage.getItem('class_token') || ''
    const res = await fetch('/api/v1/display/ai/chat', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' },
      body: JSON.stringify({ question: text }),
    })
    const data = await res.json()
    messages.value.push({ role: 'assistant', content: data?.data?.answer || '抱歉，我暂时无法回答这个问题' })
  } catch {
    messages.value.push({ role: 'assistant', content: '网络错误，请稍后重试' })
  } finally { sending.value = false }
}
</script>

<template>
  <div v-if="loading" style="text-align:center;padding:48px;color:var(--md-text-secondary);">加载中...</div>
  <div v-else-if="!aiEnabled" style="text-align:center;padding:48px;color:var(--md-text-secondary);">
    <div style="font-size:48px;margin-bottom:12px;">🤖</div>
    <p>AI 功能未开启，请联系管理员配置</p>
  </div>
  <div v-else style="display:flex;flex-direction:column;height:calc(100vh - 120px);">
    <div ref="chatBody" style="flex:1;overflow-y:auto;padding:16px;display:flex;flex-direction:column;gap:12px;">
      <div v-for="(m, i) in messages" :key="i" :style="{
        alignSelf: m.role === 'user' ? 'flex-end' : 'flex-start',
        maxWidth: '80%', padding: '10px 16px', borderRadius: '16px', fontSize: '14px', lineHeight: 1.5,
        background: m.role === 'user' ? 'var(--md-primary)' : 'var(--md-surface-2)',
        color: m.role === 'user' ? '#fff' : 'var(--md-text)',
        borderBottomRightRadius: m.role === 'user' ? '4px' : '16px',
        borderBottomLeftRadius: m.role === 'assistant' ? '4px' : '16px',
      }">{{ m.content }}</div>
      <div v-if="messages.length === 0" style="text-align:center;padding:40px;color:var(--md-text-secondary);">
        <div style="font-size:40px;margin-bottom:8px;">🤖</div>
        <p>有什么学习上的问题可以问我哦！</p>
      </div>
    </div>
    <div style="display:flex;gap:8px;padding:12px 16px;border-top:1px solid rgba(255,255,255,0.06);">
      <input v-model="input" @keyup.enter="sendMsg" placeholder="输入你的问题..." style="flex:1;padding:10px 14px;border-radius:10px;border:1px solid rgba(255,255,255,0.08);background:var(--md-surface-2);color:var(--md-text);font-size:14px;outline:none;font-family:inherit;">
      <button @click="sendMsg" :disabled="sending" style="padding:10px 20px;border-radius:10px;border:none;background:var(--md-primary);color:#fff;font-size:14px;cursor:pointer;font-family:inherit;">{{ sending ? '...' : '发送' }}</button>
    </div>
  </div>
</template>
