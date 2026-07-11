<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { formatTime } from '@/utils/constants'
import type { ApiResponse, Homework } from '@/types'

const toast = useToastStore()

const homeworks = ref<Homework[]>([])
const loading = ref(true)
const showModal = ref(false)
const sending = ref(false)

const form = ref({ title: '', deadline: '', description: '' })

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Homework[]>>('/api/v1/teacher/homework')
    homeworks.value = res.data || []
  } catch { /* handled */ }
  finally { loading.value = false }
})

function openCreate() {
  // 默认截止时间：3天后
  const d = new Date(); d.setDate(d.getDate() + 3)
  form.value = { title: '', deadline: d.toISOString().slice(0, 16), description: '' }
  showModal.value = true
}

async function handleCreate() {
  if (!form.value.title.trim() || !form.value.deadline) {
    toast.show('请填写作业标题和截止时间', 'error')
    return
  }
  sending.value = true
  try {
    await apiPost('/api/v1/teacher/homework', {
      title: form.value.title.trim(),
      deadline: form.value.deadline,
      description: form.value.description.trim(),
    })
    toast.show('作业已创建', 'success')
    showModal.value = false
    // 重新拉取列表（后端当前返回空，但接口已创建）
    const res = await apiGet<ApiResponse<Homework[]>>('/api/v1/teacher/homework')
    homeworks.value = res.data || []
  } catch { /* handled */ }
  finally { sending.value = false }
}

function progressOf(hw: Homework): number {
  if (!hw.total_students) return 0
  return Math.round(hw.submission_count / hw.total_students * 100)
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">作业管理</h2>
      <button class="btn btn-sm btn-primary" @click="openCreate">创建作业</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="homeworks.length === 0" class="card" style="text-align:center;padding:64px 48px;color:var(--color-text-secondary);">
      <div style="font-size:56px;margin-bottom:16px;">📝</div>
      <h3 style="font-size:18px;font-weight:600;color:var(--color-text);margin-bottom:8px;">作业功能即将上线</h3>
      <p style="margin-bottom:20px;">后端接口正在开发中，创建作业功能已可用</p>
      <button class="btn btn-primary" style="width:auto;" @click="openCreate">创建作业</button>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(320px,1fr));gap:16px;">
      <div v-for="hw in homeworks" :key="hw.id" class="card" style="padding:16px;">
        <div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px;">
          <h3 style="font-size:15px;font-weight:600;">{{ hw.title }}</h3>
          <span :style="{
            fontSize: '12px',
            padding: '2px 8px',
            borderRadius: 'var(--radius-sm)',
            background: hw.status === 'active' ? 'rgba(16,185,129,0.1)' : 'var(--color-bg)',
            color: hw.status === 'active' ? 'var(--color-accent)' : 'var(--color-text-secondary)',
          }">{{ hw.status === 'active' ? '进行中' : '已关闭' }}</span>
        </div>

        <div style="font-size:13px;color:var(--color-text-secondary);margin-bottom:12px;">
          截止：{{ formatTime(hw.deadline) }}
        </div>

        <div style="margin-bottom:8px;">
          <div style="display:flex;justify-content:space-between;font-size:12px;margin-bottom:4px;">
            <span style="color:var(--color-text-secondary);">提交进度</span>
            <span style="font-weight:600;">{{ hw.submission_count }}/{{ hw.total_students }}</span>
          </div>
          <div style="height:8px;background:var(--color-bg);border-radius:4px;overflow:hidden;">
            <div :style="{
              width: `${progressOf(hw)}%`,
              height: '100%',
              background: progressOf(hw) >= 100 ? 'var(--color-accent)' : 'var(--color-primary)',
            }"></div>
          </div>
        </div>

        <div style="display:flex;gap:8px;margin-top:12px;">
          <button v-if="hw.status === 'active'" class="btn btn-sm btn-ghost" @click="toast.show('二维码功能开发中', 'error')">查看二维码</button>
          <button class="btn btn-sm btn-ghost" @click="toast.show('提交列表功能开发中', 'error')">查看提交</button>
        </div>
      </div>
    </div>

    <!-- 创建作业弹窗 -->
    <div v-if="showModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100;" @click.self="showModal = false">
      <div class="card" style="width:440px;max-width:90vw;padding:24px;">
        <h3 style="font-size:18px;font-weight:600;margin-bottom:20px;">创建作业</h3>
        <div class="form-group">
          <label>作业标题</label>
          <input v-model="form.title" class="form-input" placeholder="如：第一单元练习">
        </div>
        <div class="form-group">
          <label>截止时间</label>
          <input v-model="form.deadline" type="datetime-local" class="form-input">
        </div>
        <div class="form-group">
          <label>作业说明</label>
          <textarea v-model="form.description" class="form-input" rows="3" placeholder="作业内容描述（选填）" style="resize:vertical;font-family:inherit;"></textarea>
        </div>
        <div style="display:flex;justify-content:flex-end;gap:8px;margin-top:8px;">
          <button class="btn btn-ghost" @click="showModal = false">取消</button>
          <button class="btn btn-primary" style="width:auto;" :disabled="sending" @click="handleCreate">
            {{ sending ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
