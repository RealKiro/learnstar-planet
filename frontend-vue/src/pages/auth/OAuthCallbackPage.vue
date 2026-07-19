<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiPost } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, User } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

onMounted(async () => {
  const platform = (route.query.platform as string) || ''
  const code = (route.query.code as string) || ''
  const state = (route.query.state as string) || ''

  if (!platform || !code) {
    toast.show('OAuth 参数缺失', 'error')
    // 尝试关闭弹出窗口，否则跳转回登录页
    if (window.opener) { window.close() }
    else { router.replace({ name: 'login' }) }
    return
  }

  // 如果是弹出窗口，把参数传给父窗口
  if (window.opener) {
    window.opener.postMessage({
      type: 'oauth_callback',
      platform,
      code,
      state,
    }, window.location.origin)
    window.close()
    return
  }

  // 直接跳转（非弹出窗口模式）: 通过后端 OAuth 登录
  try {
    let res: any
    const baseUrl = '/api/v1/auth/teacher'

    switch (platform) {
      case 'wechat':
        // 这里需要后端根据 code 换取 openid
        res = await apiPost<ApiResponse<{ token: string; user: User }>>(`${baseUrl}/login/wechat`, { code })
        break
      case 'wechat_work':
        res = await apiPost<ApiResponse<{ token: string; user: User }>>(`${baseUrl}/login/wechat-work`, { code })
        break
      case 'qq':
        res = await apiPost<ApiResponse<{ token: string; user: User }>>(`${baseUrl}/login/qq`, { code })
        break
      default:
        toast.show('不支持的登录平台', 'error')
        router.replace({ name: 'login' })
        return
    }

    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'teacher-dashboard' })
  } catch (e: any) {
    toast.show(e?.response?.data?.message || `扫码登录失败，请使用账号密码登录`, 'error')
    setTimeout(() => router.replace({ name: 'login' }), 1500)
  }
})
</script>

<template>
  <div style="display:flex;align-items:center;justify-content:center;min-height:100vh;background:#F8FAFC;">
    <div style="text-align:center;">
      <div style="font-size:48px;margin-bottom:16px;">🔄</div>
      <h2 style="font-size:20px;font-weight:600;margin-bottom:8px;">正在处理扫码登录...</h2>
      <p style="color:#64748B;">请稍候</p>
    </div>
  </div>
</template>
