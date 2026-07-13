<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { platformLabel } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const authStore = useAuthStore()
const toast = useToastStore()

interface BindingInfo { platform: string; bound: boolean }
interface DisplayCodeInfo {
  code: string
  class_name: string
  updated_at: string
  student_count: number
}

const bindings = ref<BindingInfo[]>([])
const platforms = ['wechat', 'wechat_work', 'qq', 'renren'] as const

// 班级大屏码
const displayCode = ref<DisplayCodeInfo | null>(null)
const displayCodeLoading = ref(false)
const displayCodeUrl = ref('')

onMounted(async () => {
  await Promise.all([
    loadBindings(),
    loadDisplayCode(),
  ])
})

async function loadBindings() {
  try {
    const res = await apiGet<ApiResponse<BindingInfo[]>>('/api/v1/auth/bindings')
    bindings.value = res.data || []
  } catch {
    bindings.value = platforms.map(p => ({ platform: p, bound: false }))
  }
}

// ===== 班级大屏码管理 =====

async function loadDisplayCode() {
  displayCodeLoading.value = true
  try {
    const res = await apiGet<{ data: DisplayCodeInfo }>('/api/v1/teacher/display-code')
    displayCode.value = res.data
    displayCodeUrl.value = `${window.location.origin}/display?code=${res.data.code}`
  } catch { /* 无活跃班级时静默处理 */ }
  finally { displayCodeLoading.value = false }
}

async function refreshDisplayCode() {
  if (!confirm('刷新班级码后，旧码将立即失效，确定继续？')) return
  try {
    const res = await apiPost<{ data: DisplayCodeInfo }>('/api/v1/teacher/display-code/refresh')
    displayCode.value = res.data
    displayCodeUrl.value = `${window.location.origin}/display?code=${res.data.code}`
    toast.show('班级大屏码已刷新', 'success')
  } catch { /* handled */ }
}

async function copyDisplayCode() {
  if (!displayCode.value) return
  try {
    await navigator.clipboard.writeText(displayCode.value.code)
    toast.show('已复制班级码：' + displayCode.value.code, 'success')
  } catch {
    // Fallback
    const textarea = document.createElement('textarea')
    textarea.value = displayCode.value.code
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    toast.show('已复制班级码', 'success')
  }
}

function openDisplayUrl() {
  window.open(`/display`, '_blank')
}

async function unbind(platform: string) {
  if (!confirm(`确定要解绑 ${platformLabel(platform)} 吗？`)) return
  try {
    await apiDelete(`/api/v1/auth/unbind/${platform}`)
    const b = bindings.value.find(x => x.platform === platform)
    if (b) b.bound = false
    toast.show(`已解绑 ${platformLabel(platform)}`, 'success')
  } catch { /* handled */ }
}

function handleBind(platform: string) {
  toast.show(`正在打开${platformLabel(platform)}扫码绑定...`, 'success')
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">账号设置</h2>
    </div>

    <div class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:24px;">基本信息</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;">
        <div class="form-group">
          <label>账号</label>
          <input :value="authStore.user?.username" readonly class="form-input" style="background:var(--color-bg);">
        </div>
        <div class="form-group">
          <label>姓名</label>
          <input :value="authStore.displayName" class="form-input">
        </div>
        <div class="form-group">
          <label>手机</label>
          <input :value="authStore.user?.phone" class="form-input">
        </div>
        <div class="form-group">
          <label>邮箱</label>
          <input :value="authStore.user?.email" class="form-input">
        </div>
      </div>
      <button class="btn btn-primary" style="margin-top:16px;">修改密码</button>
    </div>

    <!-- 班级大屏码 -->
    <div class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:4px;">📺 班级大屏</h3>
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
        在教室触摸屏上展示班级宠物矩阵，实时接收广播和通知
      </p>

      <div v-if="displayCodeLoading && !displayCode" style="padding:16px;text-align:center;color:var(--color-text-secondary);">
        加载中...
      </div>

      <div v-else-if="!displayCode" style="padding:16px;text-align:center;color:var(--color-text-secondary);">
        请先在教师管理中指定班级
      </div>

      <div v-else class="display-code-section">
        <div class="code-display">
          <div class="code-class-name">{{ displayCode.class_name }}</div>
          <div class="code-value">{{ displayCode.code }}</div>
          <div class="code-meta">
            {{ displayCode.student_count }} 位学生 ·
            最后刷新: {{ displayCode.updated_at ? new Date(displayCode.updated_at).toLocaleDateString('zh-CN') : '--' }}
          </div>
        </div>

        <div style="display:flex;gap:8px;flex-wrap:wrap;margin-top:16px;">
          <button class="btn btn-primary" @click="openDisplayUrl">🎬 预览大屏</button>
          <button class="btn btn-sm" style="font-size:13px;" @click="copyDisplayCode">📋 复制班级码</button>
          <button class="btn btn-sm btn-ghost" style="font-size:13px;color:var(--color-danger);" @click="refreshDisplayCode">🔄 刷新码</button>
        </div>

        <div class="display-tips" style="margin-top:16px;padding:12px;background:rgba(79,70,229,0.04);border-radius:var(--radius-md);font-size:12px;color:var(--color-text-secondary);line-height:1.6;">
          <strong>💡 使用方式：</strong><br>
          ① 在教室触摸屏上打开浏览器，访问 <code style="background:rgba(0,0,0,0.06);padding:1px 6px;border-radius:4px;">/display</code><br>
          ② 输入班级码 <strong style="color:var(--color-primary);">{{ displayCode.code }}</strong> 即可进入大屏<br>
          ③ 或使用此链接直接跳转：<code style="background:rgba(0,0,0,0.06);padding:1px 6px;border-radius:4px;word-break:break-all;">{{ displayCodeUrl }}</code><br>
          ④ 大屏将实时显示学生宠物状态和广播消息
        </div>
      </div>
    </div>

    <div class="card">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:8px;">第三方账号绑定</h3>
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">绑定后可使用扫码快捷登录</p>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;">
        <div v-for="b in bindings" :key="b.platform"
          style="display:flex;align-items:center;gap:16px;padding:16px;border-radius:var(--radius-md);border:1px solid var(--color-border);">
          <span style="font-size:24px;flex-shrink:0;">{{
            { wechat: '💬', wechat_work: '💼', qq: '🐧', renren: '🌐' }[b.platform]
          }}</span>
          <div style="flex:1;">
            <div style="font-weight:500;">{{ platformLabel(b.platform) }}</div>
            <div style="font-size:12px;" :style="{ color: b.bound ? 'var(--color-accent)' : 'var(--color-text-secondary)' }">
              {{ b.bound ? '✅ 已绑定' : '未绑定' }}
            </div>
          </div>
          <button v-if="b.bound" class="btn btn-sm" style="color:var(--color-danger);" @click="unbind(b.platform)">解绑</button>
          <button v-else class="btn btn-sm btn-primary" @click="handleBind(b.platform)">绑定</button>
        </div>
      </div>
    </div>
  </div>
</template>
