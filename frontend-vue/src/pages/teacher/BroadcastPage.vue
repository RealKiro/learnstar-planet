<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import type { ApiResponse } from '@/types'

const bcContent = ref('')
const bcError = ref('')
const bcType = ref<'banner' | 'popup' | 'fullscreen'>('banner')
const bcVoice = ref(true)
const bcLoop = ref(false)
const bcDuration = ref(10)
const sendStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

// 目标班级选择
const myClasses = ref<Array<{ class_id: number; class_name: string; grade: string }>>([])
const selectedClassIds = ref<number[]>([])
const selectAll = ref(true)

const broadcasts = ref<Array<{ id: number; content: string; type: string; voice: boolean; created_at: string }>>([])

const templates = [
  { label: '🤫 自习课开始', text: '请保持安静，自习课开始' },
  { label: '👁️ 眼保健操', text: '眼保健操时间到了，请同学们开始做操' },
  { label: '🔔 下课提醒', text: '下课时间到，请注意安全' },
  { label: '🧹 值日提醒', text: '今天值日生请留下来打扫卫生' },
  { label: '📖 课前准备', text: '请翻到课本第25页，准备上课' },
]

onMounted(async () => {
  try {
    const [clsRes, bcRes] = await Promise.all([
      apiGet<{ data: Array<{ class_id: number; class_name: string; grade: string }> }>('/api/v1/teacher/my-classes'),
      apiGet<ApiResponse<typeof broadcasts.value>>('/api/v1/teacher/broadcasts'),
    ])
    myClasses.value = clsRes.data || []
    broadcasts.value = bcRes.data || []
    // 默认全选
    selectedClassIds.value = myClasses.value.map(c => c.class_id)
  } catch { /* handled */ }
})

function toggleAll() {
  selectAll.value = !selectAll.value
  selectedClassIds.value = selectAll.value
    ? myClasses.value.map(c => c.class_id)
    : []
}

async function sendBroadcast() {
  if (!bcContent.value.trim()) { bcError.value = '请输入广播内容'; return }
  if (selectedClassIds.value.length === 0) { bcError.value = '请选择至少一个目标班级'; return }

  bcError.value = ''
  sendStatus.value = 'loading'
  try {
    await apiPost('/api/v1/teacher/broadcasts', {
      content: bcContent.value.trim(),
      type: bcType.value,
      class_ids: selectedClassIds.value,
      voice: bcVoice.value,
      loop: bcLoop.value,
      duration: bcDuration.value,
    })
    sendStatus.value = 'success'
    setTimeout(() => { sendStatus.value = 'idle' }, 1500)
    bcContent.value = ''
    // 刷新记录
    const res = await apiGet<ApiResponse<typeof broadcasts.value>>('/api/v1/teacher/broadcasts')
    broadcasts.value = res.data || []
  } catch {
    sendStatus.value = 'error'
    setTimeout(() => { sendStatus.value = 'idle' }, 3000)
  }
}

function useTemplate(text: string) {
  bcContent.value = text
}

const typeLabels: Record<string, string> = { banner: '📌 横幅', popup: '💬 弹窗', fullscreen: '🖥️ 全屏' }
</script>

<template>
  <div>
    <div class="page-head">
      <h2 class="page-title">📡 实时广播</h2>
      <span class="text-muted-13">
        🖥️ 教室桌面端：<strong class="bc-accent-text">已连接</strong>
      </span>
    </div>

    <!-- 发送区 -->
    <div class="card bc-card-mb">
      <h3 class="section-title">发送广播</h3>

      <div class="bc-type-row">
        <button v-for="t in ([['banner','📌 顶部横幅'],['popup','💬 弹窗提示'],['fullscreen','🖥️ 全屏展示']] as const)" :key="t[0]" :class="['bc-type-btn', bcType === t[0] ? 'active' : '']" @click="bcType = t[0]">
          {{ t[1] }}
        </button>
      </div>

      <div class="form-group">
        <textarea v-model="bcContent" class="form-input bc-textarea"
          placeholder="输入要发送到教室的内容..."></textarea>
      </div>

      <!-- 目标班级选择 -->
      <div class="bc-mb-16" v-if="myClasses.length > 0">
        <div class="bc-field-label">
          发送至班级（{{ selectedClassIds.length }}/{{ myClasses.length }}）
        </div>
        <div class="bc-chip-wrap">
          <label
            :style="selectAll ? { border:'1px solid var(--ui-brand)', background:'var(--ui-brand-soft)' } : {}" class="bc-class-chip">
            <input type="checkbox" :checked="selectAll" @change="toggleAll" class="bc-accent">
            全部班级
          </label>
          <label v-for="c in myClasses" :key="c.class_id"
            :style="selectedClassIds.includes(c.class_id) ? { border:'1px solid var(--ui-brand)', background:'var(--ui-brand-soft)' } : {}" class="bc-class-chip">
            <input type="checkbox" :value="c.class_id" v-model="selectedClassIds"
              @change="selectAll = false" class="bc-accent">
            {{ c.class_name }}
          </label>
        </div>
      </div>

      <div class="bc-options-row">
        <label class="bc-check-row">
          <input v-model="bcVoice" type="checkbox"> 🔊 语音播报
        </label>
        <label class="bc-check-row">
          <input v-model="bcLoop" type="checkbox"> 🔁 循环播放
        </label>
        <label class="bc-inline-row">
          ⏱️ <select v-model.number="bcDuration" class="form-select bc-duration-select">
            <option :value="5">5秒</option>
            <option :value="10">10秒</option>
            <option :value="30">30秒</option>
            <option :value="0">常驻</option>
          </select>
        </label>
      </div>

      <div v-if="bcError" class="bc-error">{{ bcError }}</div>
      <button class="btn bc-send-btn" :class="'btn-state-' + sendStatus" :disabled="sendStatus === 'loading'" @click="sendBroadcast">
        {{ { idle: '📡 发送至 ' + selectedClassIds.length + ' 个班级', loading: '发送中...', success: '已发送 ✓', error: '发送失败' }[sendStatus] }}
      </button>
    </div>

    <!-- 快捷模板 -->
    <div class="card bc-card-mb">
      <h3 class="section-title">快捷模板</h3>
      <div class="bc-tpl-grid">
        <div v-for="t in templates" :key="t.label" class="card bc-tpl-card"
          @click="useTemplate(t.text)">
          {{ t.label }}
        </div>
      </div>
    </div>

    <!-- 广播记录 -->
    <div class="data-table">
      <div class="data-table__header"><h3 class="bc-table-title">广播记录</h3></div>
      <table>
        <thead><tr><th>内容</th><th>类型</th><th>语音</th><th>时间</th></tr></thead>
        <tbody>
          <tr v-if="broadcasts.length === 0">
            <td colspan="4" class="bc-empty-cell">暂无广播记录</td>
          </tr>
          <tr v-for="b in broadcasts" :key="b.id">
            <td>{{ b.content }}</td>
            <td>{{ typeLabels[b.type] || b.type }}</td>
            <td>{{ b.voice ? '🔊' : '🔇' }}</td>
            <td class="bc-secondary">{{ new Date(b.created_at).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

</template>

<style scoped>
.bc-type-btn { flex:1;padding:10px;border-radius:var(--radius-md);border:1px solid var(--color-border);background:var(--color-bg);color:var(--color-text);font-size:13px;cursor:pointer;font-weight:500;font-family:inherit; }
.bc-type-btn:hover { border-color:var(--ui-brand); }
.bc-type-btn.active { border:2px solid var(--ui-brand);background:var(--ui-brand-soft);color:var(--ui-brand); }
/* ===== P1 内联样式收口（声明逐字保留以保渲染等价） ===== */
.bc-card-mb { margin-bottom:24px; }
.bc-class-chip { display:flex;align-items:center;gap:4px;font-size:13px;cursor:pointer;padding:6px 12px;border-radius:8px;border:1px solid var(--color-border); }
.bc-accent { accent-color:var(--ui-brand); }
.bc-check-row { display:flex;align-items:center;gap:4px;font-size:14px;cursor:pointer; }
.bc-accent-text { color:var(--ui-brand); }
.bc-type-row { display:flex;gap:8px;margin-bottom:16px; }
.bc-textarea { min-height:80px;resize:vertical; }
.bc-mb-16 { margin-bottom:16px; }
.bc-field-label { font-size:13px;font-weight:600;margin-bottom:8px;color:var(--color-text-secondary); }
.bc-chip-wrap { display:flex;flex-wrap:wrap;gap:8px; }
.bc-options-row { display:flex;gap:16px;align-items:center;flex-wrap:wrap;margin-bottom:16px; }
.bc-inline-row { display:flex;align-items:center;gap:4px;font-size:14px; }
.bc-duration-select { width:auto;padding:4px 8px; }
.bc-error { margin-bottom:10px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color: var(--color-danger-text);font-size:12px; }
.bc-send-btn { width:auto;color:#fff;border:none; }
.bc-tpl-grid { display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px; }
.bc-tpl-card { padding:12px 16px;cursor:pointer;font-size:13px;text-align:center; }
.bc-table-title { font-size:16px;font-weight:600; }
.bc-empty-cell { text-align:center;color:var(--color-text-secondary);padding:24px; }
.bc-secondary { color:var(--color-text-secondary); }
</style>
