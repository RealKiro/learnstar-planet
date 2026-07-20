<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiDelete, apiPost } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { platformLabel } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const authStore = useAuthStore()
const toast = useToastStore()

interface BindingInfo { platform: string; label: string; bound: boolean; nick?: string }

const bindings = ref<BindingInfo[]>([])
const platforms = ['wechat', 'wechat_work', 'qq', 'renren'] as const

// 修改密码弹窗
const showPwdModal = ref(false)
const pwdForm = ref({ current_password: '', new_password: '', confirm_password: '' })
const changingPwd = ref(false)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<BindingInfo[]>>('/api/v1/auth/bindings')
    bindings.value = res.data || []
  } catch {
    bindings.value = platforms.map(p => ({ platform: p, label: platformLabel(p), bound: false }))
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
  const callbackUrl = `${window.location.origin}/auth/callback?platform=${platform}`
  let oauthUrl = ''

  switch (platform) {
    case 'wechat':
      oauthUrl = `https://open.weixin.qq.com/connect/qrconnect?appid=YOUR_WECHAT_APPID&redirect_uri=${encodeURIComponent(callbackUrl)}&response_type=code&scope=snsapi_login&state=${platform}`
      break
    case 'wechat_work':
      oauthUrl = `https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=YOUR_WECHAT_WORK_CORPID&agentid=YOUR_AGENTID&redirect_uri=${encodeURIComponent(callbackUrl)}&state=${platform}`
      break
    case 'qq':
      oauthUrl = `https://graph.qq.com/oauth2.0/show?which=Login&display=pc&client_id=YOUR_QQ_APPID&redirect_uri=${encodeURIComponent(callbackUrl)}&response_type=code&state=${platform}`
      break
    case 'renren':
      toast.show('请在管理员后台配置人人通参数', 'info')
      return
  }

  if (oauthUrl) {
    const w = 600, h = 500
    window.open(oauthUrl, platform,
      `width=${w},height=${h},left=${(screen.width-w)/2},top=${(screen.height-h)/2}`)
  }
}

function openPwdModal() {
  pwdForm.value = { current_password: '', new_password: '', confirm_password: '' }
  showPwdModal.value = true
}

async function changePassword() {
  const f = pwdForm.value
  if (!f.current_password || !f.new_password) { toast.show('请填写完整', 'error'); return }
  if (f.new_password.length < 6) { toast.show('新密码至少 6 位', 'error'); return }
  if (f.new_password !== f.confirm_password) { toast.show('两次密码不一致', 'error'); return }

  changingPwd.value = true
  try {
    await apiPost('/api/v1/auth/change-password', {
      current_password: f.current_password,
      new_password: f.new_password,
    })
    toast.show('密码修改成功', 'success')
    showPwdModal.value = false
  } catch { /* handled */ } finally { changingPwd.value = false }
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
          <input :value="authStore.displayName" readonly class="form-input" style="background:var(--color-bg);">
        </div>
      </div>
      <button class="btn btn-primary" style="margin-top:16px;" @click="openPwdModal">修改密码</button>
    </div>

    <!-- 修改密码弹窗 -->
    <Teleport to="body">
      <div v-if="showPwdModal" @click.self="showPwdModal = false"
        style="position:fixed;inset:0;z-index:999;background:rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:100%;box-shadow:0 20px 60px rgba(0,0,0,0.15);">
          <h3 style="font-size:18px;font-weight:700;margin-bottom:20px;">修改密码</h3>
          <div class="form-group" style="margin-bottom:12px;">
            <label>当前密码</label>
            <input v-model="pwdForm.current_password" type="password" class="form-input" placeholder="输入当前密码">
          </div>
          <div class="form-group" style="margin-bottom:12px;">
            <label>新密码</label>
            <input v-model="pwdForm.new_password" type="password" class="form-input" placeholder="至少 6 位">
          </div>
          <div class="form-group" style="margin-bottom:20px;">
            <label>确认新密码</label>
            <input v-model="pwdForm.confirm_password" type="password" class="form-input" placeholder="再次输入新密码">
          </div>
          <div style="display:flex;gap:12px;">
            <button class="btn" style="flex:1;" @click="showPwdModal = false">取消</button>
            <button class="btn btn-primary" style="flex:1;" :disabled="changingPwd" @click="changePassword">
              {{ changingPwd ? '修改中...' : '确认修改' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

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
            <div style="font-weight:500;">{{ b.label || platformLabel(b.platform) }}</div>
            <div style="font-size:12px;" :style="{ color: b.bound ? 'var(--color-accent)' : 'var(--color-text-secondary)' }">
              {{ b.bound ? '✅ 已绑定' : '未绑定' }}
              <span v-if="b.bound && b.nick" style="margin-left:4px;">（{{ b.nick }}）</span>
            </div>
          </div>
          <button v-if="b.bound" class="btn btn-sm" style="color:var(--color-danger);" @click="unbind(b.platform)">解绑</button>
          <button v-else class="btn btn-sm btn-primary" @click="handleBind(b.platform)">绑定</button>
        </div>
      </div>
    </div>
  </div>
</template>
