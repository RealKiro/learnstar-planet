<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiDelete } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { platformLabel } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const authStore = useAuthStore()
const toast = useToastStore()

interface BindingInfo { platform: string; bound: boolean }

const bindings = ref<BindingInfo[]>([])
const platforms = ['wechat', 'wechat_work', 'qq', 'renren'] as const

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<BindingInfo[]>>('/api/v1/auth/bindings')
    bindings.value = res.data || []
  } catch {
    bindings.value = platforms.map(p => ({ platform: p, bound: false }))
  }
})

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
