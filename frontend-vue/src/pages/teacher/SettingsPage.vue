<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
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
const unbindStatus = ref<Record<string, 'idle' | 'loading' | 'success' | 'error'>>({})
const pwdStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

// 修改密码弹窗
const showPwdModal = ref(false)
const pwdForm = ref({ current_password: '', new_password: '', confirm_password: '' })
const changingPwd = ref(false)
const pwdErrors = reactive<Record<string, string>>({})
function pwdClr(f: string) { delete pwdErrors[f] }
function pwdVld(field: string): boolean {
  if (field === 'current_password' && !pwdForm.value.current_password) { pwdErrors.current_password = '请输入当前密码'; return false }
  if (field === 'new_password' && !pwdForm.value.new_password) { pwdErrors.new_password = '请输入新密码'; return false }
  if (field === 'new_password' && pwdForm.value.new_password.length < 6) { pwdErrors.new_password = '密码至少 6 位'; return false }
  if (field === 'confirm_password' && pwdForm.value.new_password !== pwdForm.value.confirm_password) { pwdErrors.confirm_password = '两次密码不一致'; return false }
  delete pwdErrors[field]; return true
}

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
  unbindStatus.value[platform] = 'loading'
  try {
    await apiDelete(`/api/v1/auth/unbind/${platform}`)
    const b = bindings.value.find(x => x.platform === platform)
    if (b) b.bound = false
    unbindStatus.value[platform] = 'success'
    setTimeout(() => { unbindStatus.value[platform] = 'idle' }, 1500)
  } catch {
    unbindStatus.value[platform] = 'error'
    setTimeout(() => { unbindStatus.value[platform] = 'idle' }, 3000)
  }
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
      toast.show('请在管理员后台配置人人通参数', 'info', { position: 'center', duration: 2000 })
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
  if (!pwdVld('current_password') | !pwdVld('new_password') | !pwdVld('confirm_password')) return
  pwdStatus.value = 'loading'
  try {
    await apiPost('/api/v1/auth/change-password', {
      current_password: pwdForm.value.current_password,
      new_password: pwdForm.value.new_password,
    })
    pwdStatus.value = 'success'
    setTimeout(() => { showPwdModal.value = false; pwdStatus.value = 'idle' }, 800)
  } catch {
    pwdStatus.value = 'error'
    setTimeout(() => { pwdStatus.value = 'idle' }, 3000)
  }
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
      <div v-if="showPwdModal" @click="showPwdModal = false"
        style="position:fixed;inset:0;z-index:999;background:rgba(0,0,0,0.15);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div @click.stop style="background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:16px;padding:24px;max-width:400px;width:100%;box-shadow:0 8px 32px rgba(0,0,0,0.12);">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
            <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">🔑 修改密码</h3>
            <button @click="showPwdModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
          </div>
          <div class="form-group" style="margin-bottom:12px;">
            <label>当前密码</label>
            <input v-model="pwdForm.current_password" type="password" class="form-input" placeholder="输入当前密码" :style="{ borderColor: pwdErrors.current_password ? '#f87171' : '' }" @blur="pwdVld('current_password')" @input="pwdClr('current_password')">
            <div v-if="pwdErrors.current_password" style="color:#f87171;font-size:11px;margin-top:2px;">{{ pwdErrors.current_password }}</div>
          </div>
          <div class="form-group" style="margin-bottom:12px;">
            <label>新密码</label>
            <input v-model="pwdForm.new_password" type="password" class="form-input" placeholder="至少 6 位" :style="{ borderColor: pwdErrors.new_password ? '#f87171' : '' }" @blur="pwdVld('new_password')" @input="pwdClr('new_password')">
            <div v-if="pwdErrors.new_password" style="color:#f87171;font-size:11px;margin-top:2px;">{{ pwdErrors.new_password }}</div>
          </div>
          <div class="form-group" style="margin-bottom:20px;">
            <label>确认新密码</label>
            <input v-model="pwdForm.confirm_password" type="password" class="form-input" placeholder="再次输入新密码" :style="{ borderColor: pwdErrors.confirm_password ? '#f87171' : '' }" @blur="pwdVld('confirm_password')" @input="pwdClr('confirm_password')">
            <div v-if="pwdErrors.confirm_password" style="color:#f87171;font-size:11px;margin-top:2px;">{{ pwdErrors.confirm_password }}</div>
          </div>
          <div style="display:flex;gap:12px;">
            <button class="btn" style="flex:1;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);" @click="showPwdModal = false">取消</button>
            <button class="btn btn-primary" style="flex:1;" :style="{ background: pwdStatus === 'loading' ? '#f59e0b' : pwdStatus === 'success' ? '#10b981' : pwdStatus === 'error' ? '#ef4444' : '' }" :disabled="pwdStatus === 'loading'" @click="changePassword">
              <template v-if="pwdStatus === 'loading'">修改中...</template>
              <template v-else-if="pwdStatus === 'success'">修改成功 ✓</template>
              <template v-else-if="pwdStatus === 'error'">修改失败 ✗</template>
              <template v-else>确认修改</template>
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
          <button v-if="b.bound" class="btn btn-sm" :style="{ background: unbindStatus[b.platform] === 'loading' ? '#f59e0b' : unbindStatus[b.platform] === 'success' ? '#10b981' : unbindStatus[b.platform] === 'error' ? '#ef4444' : '', color: unbindStatus[b.platform] && unbindStatus[b.platform] !== 'idle' ? '#fff' : 'var(--color-danger)', border: unbindStatus[b.platform] && unbindStatus[b.platform] !== 'idle' ? '1px solid transparent' : '' }" :disabled="unbindStatus[b.platform] === 'loading'" @click="unbind(b.platform)">
            <template v-if="unbindStatus[b.platform] === 'loading'">解绑中...</template>
            <template v-else-if="unbindStatus[b.platform] === 'success'">已解绑 ✓</template>
            <template v-else-if="unbindStatus[b.platform] === 'error'">解绑失败 ✗</template>
            <template v-else>解绑</template>
          </button>
          <button v-else class="btn btn-sm btn-primary" @click="handleBind(b.platform)">绑定</button>
        </div>
      </div>
    </div>
  </div>
</template>
