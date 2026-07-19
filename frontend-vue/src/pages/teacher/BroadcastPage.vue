<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

const bcContent = ref('')
const bcType = ref<'banner' | 'popup' | 'fullscreen'>('banner')
const bcVoice = ref(true)
const bcLoop = ref(false)
const bcDuration = ref(10)

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
  if (!bcContent.value.trim()) { toast.show('请输入广播内容', 'error'); return }
  if (selectedClassIds.value.length === 0) { toast.show('请选择至少一个目标班级', 'error'); return }

  try {
    await apiPost('/api/v1/teacher/broadcasts', {
      content: bcContent.value.trim(),
      type: bcType.value,
      class_ids: selectedClassIds.value,
      voice: bcVoice.value,
      loop: bcLoop.value,
      duration: bcDuration.value,
    })
    toast.show(`📡 广播已发送至 ${selectedClassIds.value.length} 个班级`, 'success')
    bcContent.value = ''
    // 刷新记录
    const res = await apiGet<ApiResponse<typeof broadcasts.value>>('/api/v1/teacher/broadcasts')
    broadcasts.value = res.data || []
  } catch { /* handled */ }
}

function useTemplate(text: string) {
  bcContent.value = text
}

const typeLabels: Record<string, string> = { banner: '📌 横幅', popup: '💬 弹窗', fullscreen: '🖥️ 全屏' }
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">📡 实时广播</h2>
      <span style="font-size:13px;color:var(--color-text-secondary);">
        🖥️ 教室桌面端：<strong style="color:var(--color-accent);">已连接</strong>
      </span>
    </div>

    <!-- 发送区 -->
    <div class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">发送广播</h3>

      <div style="display:flex;gap:8px;margin-bottom:16px;">
        <button v-for="t in ([['banner','📌 顶部横幅'],['popup','💬 弹窗提示'],['fullscreen','🖥️ 全屏展示']] as const)" :key="t[0]"
          :style="bcType === t[0] ? { border:'2px solid var(--color-accent)', background:'rgba(79,70,229,0.1)', color:'var(--color-accent)' } : {}"
          style="flex:1;padding:10px;border-radius:var(--radius-md);border:1px solid var(--color-border);background:var(--color-bg);font-size:13px;cursor:pointer;font-weight:500;"
          @click="bcType = t[0]">
          {{ t[1] }}
        </button>
      </div>

      <div class="form-group">
        <textarea v-model="bcContent" class="form-input" style="min-height:80px;resize:vertical;"
          placeholder="输入要发送到教室的内容..."></textarea>
      </div>

      <!-- 目标班级选择 -->
      <div style="margin-bottom:16px;" v-if="myClasses.length > 0">
        <div style="font-size:13px;font-weight:600;margin-bottom:8px;color:var(--color-text-secondary);">
          发送至班级（{{ selectedClassIds.length }}/{{ myClasses.length }}）
        </div>
        <div style="display:flex;flex-wrap:wrap;gap:8px;">
          <label style="display:flex;align-items:center;gap:4px;font-size:13px;cursor:pointer;padding:6px 12px;border-radius:8px;border:1px solid var(--color-border);"
            :style="selectAll ? { border:'1px solid var(--color-accent)', background:'rgba(79,70,229,0.08)' } : {}">
            <input type="checkbox" :checked="selectAll" @change="toggleAll" style="accent-color:var(--color-accent);">
            全部班级
          </label>
          <label v-for="c in myClasses" :key="c.class_id"
            style="display:flex;align-items:center;gap:4px;font-size:13px;cursor:pointer;padding:6px 12px;border-radius:8px;border:1px solid var(--color-border);"
            :style="selectedClassIds.includes(c.class_id) ? { border:'1px solid var(--color-accent)', background:'rgba(79,70,229,0.08)' } : {}">
            <input type="checkbox" :value="c.class_id" v-model="selectedClassIds"
              style="accent-color:var(--color-accent);" @change="selectAll = false">
            {{ c.class_name }}
          </label>
        </div>
      </div>

      <div style="display:flex;gap:16px;align-items:center;flex-wrap:wrap;margin-bottom:16px;">
        <label style="display:flex;align-items:center;gap:4px;font-size:14px;cursor:pointer;">
          <input v-model="bcVoice" type="checkbox"> 🔊 语音播报
        </label>
        <label style="display:flex;align-items:center;gap:4px;font-size:14px;cursor:pointer;">
          <input v-model="bcLoop" type="checkbox"> 🔁 循环播放
        </label>
        <label style="display:flex;align-items:center;gap:4px;font-size:14px;">
          ⏱️ <select v-model.number="bcDuration" class="form-select" style="width:auto;padding:4px 8px;">
            <option :value="5">5秒</option>
            <option :value="10">10秒</option>
            <option :value="30">30秒</option>
            <option :value="0">常驻</option>
          </select>
        </label>
      </div>

      <button class="btn btn-primary" style="width:auto;" @click="sendBroadcast">
        📡 发送至 {{ selectedClassIds.length }} 个班级
      </button>
    </div>

    <!-- 快捷模板 -->
    <div class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">快捷模板</h3>
      <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px;">
        <div v-for="t in templates" :key="t.label" class="card"
          style="padding:12px 16px;cursor:pointer;font-size:13px;text-align:center;"
          @click="useTemplate(t.text)">
          {{ t.label }}
        </div>
      </div>
    </div>

    <!-- 广播记录 -->
    <div class="data-table">
      <div class="data-table__header"><h3 style="font-size:16px;font-weight:600;">广播记录</h3></div>
      <table>
        <thead><tr><th>内容</th><th>类型</th><th>语音</th><th>时间</th></tr></thead>
        <tbody>
          <tr v-if="broadcasts.length === 0">
            <td colspan="4" style="text-align:center;color:var(--color-text-secondary);padding:24px;">暂无广播记录</td>
          </tr>
          <tr v-for="b in broadcasts" :key="b.id">
            <td>{{ b.content }}</td>
            <td>{{ typeLabels[b.type] || b.type }}</td>
            <td>{{ b.voice ? '🔊' : '🔇' }}</td>
            <td style="color:var(--color-text-secondary);">{{ new Date(b.created_at).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
