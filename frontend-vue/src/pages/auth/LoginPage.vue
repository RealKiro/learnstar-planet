<script setup lang="ts">
import { ref, reactive, nextTick, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import { apiPost } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, User } from '@/types'

const props = withDefaults(defineProps<{ initialRole?: string; mode?: string }>(), { initialRole: 'teacher', mode: 'account' })
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToastStore()

const loginType = ref<'teacher' | 'admin' | 'parent' | 'class'>((props.mode === 'code' ? 'class' : props.initialRole) as 'teacher' | 'admin' | 'parent' | 'class')
const teacherUsername = ref('')
const teacherPassword = ref('')
const adminUsername = ref('')
const adminPassword = ref('')
const parentUsername = ref('')
const parentPassword = ref('')
const classCode = ref('')
const classCodeError = ref('')
const classInfo = ref<{ class_name: string; student_count: number } | null>(null)
const loading = ref(false)

// 内联校验
const loginErrors = reactive<Record<string, string>>({})
function clearLoginErr(f: string) { delete loginErrors[f] }
function validateLoginField(field: string, val: string): boolean {
  if (field === 'teacherUsername' && !val.trim()) { loginErrors.teacherUsername = '请输入教师账号'; return false }
  if (field === 'teacherPassword' && !val.trim()) { loginErrors.teacherPassword = '请输入教师登录密码'; return false }
  if (field === 'adminUsername' && !val.trim()) { loginErrors.adminUsername = '请输入管理员账号'; return false }
  if (field === 'adminPassword' && !val.trim()) { loginErrors.adminPassword = '请输入管理员密码'; return false }
  if (field === 'parentUsername' && !val.trim()) { loginErrors.parentUsername = '请输入家长账号'; return false }
  if (field === 'parentPassword' && !val.trim()) { loginErrors.parentPassword = '请输入家长密码'; return false }
  if (field === 'classCode' && !val.trim()) { loginErrors.classCode = '请输入班级码'; return false }
  clearLoginErr(field); return true
}
function validateLoginForm(type: string): boolean {
  Object.keys(loginErrors).forEach(k => delete loginErrors[k])
  if (type === 'teacher') {
    if (!teacherUsername.value.trim()) loginErrors.teacherUsername = '请输入教师账号'
    if (!teacherPassword.value.trim()) loginErrors.teacherPassword = '请输入教师登录密码'
  } else if (type === 'admin') {
    if (!adminUsername.value.trim()) loginErrors.adminUsername = '请输入管理员账号'
    if (!adminPassword.value.trim()) loginErrors.adminPassword = '请输入管理员密码'
  } else if (type === 'parent') {
    if (!parentUsername.value.trim()) loginErrors.parentUsername = '请输入家长账号'
    if (!parentPassword.value.trim()) loginErrors.parentPassword = '请输入家长密码'
  } else if (type === 'class') {
    if (!classCode.value.trim()) loginErrors.classCode = '请输入班级码'
  }
  return Object.keys(loginErrors).length === 0
}

const teacherPwdRef = ref<HTMLInputElement>()
const adminPwdRef = ref<HTMLInputElement>()
const parentPwdRef = ref<HTMLInputElement>()

let teacherAttempts = 0
const MAX_ATTEMPTS = 3

function focusTeacherPwd() { if (teacherUsername.value.trim()) nextTick(() => teacherPwdRef.value?.focus()) }
function focusAdminPwd() { if (adminUsername.value.trim()) nextTick(() => adminPwdRef.value?.focus()) }
function focusParentPwd() { if (parentUsername.value.trim()) nextTick(() => parentPwdRef.value?.focus()) }

async function handleTeacherLogin() {
  if (!validateLoginForm('teacher')) return
  loading.value = true
  try {
    const res = await fetch('/api/v1/auth/teacher/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: teacherUsername.value.trim(), password: teacherPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) { throw { response: { data } } }
    teacherAttempts = 0
    authStore.setAuth(data.data.token, data.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'teacher-dashboard' })
  } catch (e: any) {
    teacherAttempts++
    const remaining = MAX_ATTEMPTS - teacherAttempts
    if (teacherAttempts >= MAX_ATTEMPTS) toast.show('密码错误次数过多，请联系管理员', 'error')
    else toast.show((e?.response?.data?.message || '账号或密码错误') + `，还剩 ${remaining} 次`, 'error')
  } finally { loading.value = false }
}

async function handleAdminLogin() {
  if (!validateLoginForm('admin')) return
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/v1/auth/admin/login', {
      username: adminUsername.value.trim(), password: adminPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'admin-dashboard' })
  } catch { /* handled */ } finally { loading.value = false }
}

async function handleParentLogin() {
  if (!validateLoginForm('parent')) return
  loading.value = true
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/v1/auth/parent/login', {
      username: parentUsername.value.trim(), password: parentPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: 'parent-home' })
  } catch { /* handled */ } finally { loading.value = false }
}

async function handleClassLogin() {
  if (!validateLoginForm('class')) return
  loading.value = true
  classCodeError.value = ''
  try {
    const res = await apiPost<ApiResponse<{ token: string; class_id: number; class_name: string; grade: string; student_count: number }>>('/api/v1/auth/class/login', {
      class_code: classCode.value.trim(),
    })
    sessionStorage.setItem('class_token', res.data.token)
    sessionStorage.setItem('class_info', JSON.stringify({
      id: res.data.class_id,
      name: res.data.class_name,
      grade: res.data.grade,
      student_count: res.data.student_count,
    }))
    classInfo.value = { class_name: res.data.class_name, student_count: res.data.student_count }
    toast.show(`欢迎进入 ${res.data.class_name}`, 'success')
    setTimeout(() => router.push({ name: 'teacher-dashboard-basic' }), 1500)
  } catch (e: any) {
    classCodeError.value = e?.response?.data?.message || '班级码无效，请核对后重试'
  } finally { loading.value = false }
}

const platforms = [
  { key: 'wechat', label: '微信', icon: '💬', color: '#07C160' },
  { key: 'wechat_work', label: '企业微信', icon: '💼', color: '#2B7CE9' },
  { key: 'qq', label: 'QQ', icon: '🐧', color: '#12B7F5' },
  { key: 'renren', label: '人人通', icon: '🌐', color: '#FF6A00' },
]

function handleThirdPartyLogin(platform: string) {
  const label = platforms.find(p => p.key === platform)?.label || platform
  toast.show(`正在打开${label}扫码...`, 'success')

  // 构建 OAuth 回调地址
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
      toast.show('人人通登录需要管理员在后台配置', 'info')
      return
  }

  if (oauthUrl) {
    // 打开 OAuth 窗口（弹出或当前页跳转）
    const w = 600, h = 500
    const left = (screen.width - w) / 2
    const top = (screen.height - h) / 2
    window.open(oauthUrl, platform,
      `width=${w},height=${h},left=${left},top=${top},menubar=no,toolbar=no,status=no,scrollbars=yes`)
  }
}

const slides = [
  {
    badge: 'MIT 开源 完全免费 自托管',
    title: '让每个孩子的努力', highlight: '都被看见',
    desc: '积分激励 宠物养成 AI 助教 多端同步\n开源班级管理系统，数据完全自主掌控',
    icon: '🌌',
  },
  {
    badge: '12 大功能模块',
    title: '覆盖班级管理', highlight: '全场景',
    desc: '积分规则 宠物进化 排行榜 通知公告\n考勤 作业 答题 成绩 商城 广播 AI',
    icon: '⚡',
  },
  {
    badge: '11 阶宠物进化',
    title: '积分变经验', highlight: '驱动成长',
    desc: '星尘 月芽 灵苗 青藤 慧树 蝶灵\n鹰慧 狮睿 灵角 星耀 银河',
    icon: '🌟',
  },
  {
    badge: 'Docker 一键部署',
    title: '4 种数据库', highlight: '自由选择',
    desc: 'MySQL PostgreSQL SQLite MariaDB\n内置 Redis 缓存，支持多端登录',
    icon: '🐳',
  },
]
const currentSlide = ref(0)
let slideTimer: ReturnType<typeof setInterval>

// 接收 OAuth 弹窗回调
function handleOAuthMessage(e: MessageEvent) {
  if (e.origin !== window.location.origin) return
  const data = e.data
  if (!data || data.type !== 'oauth_callback') return

  const platform = data.platform
  if (platform === 'wechat_work') {
    handleWechatWorkOAuth(data.code)
  }
}

async function handleWechatWorkOAuth(code: string) {
  loading.value = true
  try {
    const res = await apiPost<{ data: { token: string; user: any } }>('/api/v1/auth/teacher/login/wechat-work', { code })
    const d = res.data
    if (d.token) {
      const authStore = (await import('@/stores/auth')).useAuthStore()
      const toast = (await import('@/stores/toast')).useToastStore()
      authStore.setAuth(d.token, d.user)
      toast.show('登录成功', 'success')
      router.replace({ name: 'teacher-dashboard' })
    }
  } catch { /* handled */ } finally { loading.value = false }
}

onMounted(() => {
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
  window.addEventListener('message', handleOAuthMessage)
})
onUnmounted(() => {
  clearInterval(slideTimer)
  window.removeEventListener('message', handleOAuthMessage)
})

function goToSlide(i: number) {
  currentSlide.value = i
  clearInterval(slideTimer)
  slideTimer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
}
</script>

<template>
  <div class="login-page">
    <!-- 顶部导航 -->
    <nav class="login-topnav">
      <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="topnav-link">GitHub</a>
    </nav>

    <div class="intro">
      <div class="intro-orb intro-orb--top"></div>
      <div class="intro-orb intro-orb--bottom"></div>
      <div class="intro-content">
        <transition name="slide-fade" mode="out-in">
          <div :key="currentSlide" class="intro-slide">
            <div class="intro-badge">
              <span class="intro-badge-dot"></span>
              {{ slides[currentSlide].badge }}
            </div>
            <div class="intro-icon">{{ slides[currentSlide].icon }}</div>
            <h1 class="intro-title">
              {{ slides[currentSlide].title }}<br>
              <span class="intro-title-grad">{{ slides[currentSlide].highlight }}</span>
            </h1>
            <p class="intro-desc">{{ slides[currentSlide].desc }}</p>
          </div>
        </transition>
        <div class="intro-dots">
          <button
            v-for="(slide, i) in slides" :key="i"
            :class="['intro-dot', { 'intro-dot--active': currentSlide === i }]"
            @click="goToSlide(i)"
          ></button>
        </div>
        <div class="intro-footer">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="intro-footer-link">GitHub</a>
          <span class="intro-footer-sep">MIT 开源许可证</span>
        </div>
      </div>
    </div>

    <div class="login-panel">
      <div class="login-card">
        <div class="login-card-header">
          <span class="login-card-icon">🌌</span>
          <h1 class="login-card-title">学趣星球</h1>
        </div>
        <div class="login-tabs">
          <button
            v-for="t in (['teacher', 'admin', 'parent', 'class'] as const)" :key="t"
            :class="['login-tab', { 'login-tab--active': loginType === t }]"
            @click="loginType = t"
          >{{ t === 'teacher' ? '教师' : t === 'admin' ? '管理员' : t === 'parent' ? '家长' : '🔑 班级码' }}</button>
        </div>

        <div v-if="loginType === 'teacher'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="teacherUsername" class="form-input" :style="{ borderColor: loginErrors.teacherUsername ? '#f87171' : '' }" @blur="validateLoginField('teacherUsername', teacherUsername)" @input="clearLoginErr('teacherUsername')" placeholder="教师账号" @keydown.enter="focusTeacherPwd">
            <div v-if="loginErrors.teacherUsername" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.teacherUsername }}</div>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input ref="teacherPwdRef" v-model="teacherPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleTeacherLogin">
            <div v-if="loginErrors.teacherPassword" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.teacherPassword }}</div>
          </div>
          <button class="login-submit" :disabled="loading" @click="handleTeacherLogin">{{ loading ? '登录中...' : '登录' }}</button>
          <div class="login-social">
            <div class="login-social-label"><span class="login-social-line"></span> 扫码快捷登录 <span class="login-social-line"></span></div>
            <div class="login-social-grid">
              <button v-for="p in platforms" :key="p.key" class="login-social-btn" @click="handleThirdPartyLogin(p.key)">
                <span class="login-social-icon" :style="{ background: p.color }">{{ p.icon }}</span>
                {{ p.label }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="loginType === 'admin'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="adminUsername" class="form-input" :style="{ borderColor: loginErrors.adminUsername ? '#f87171' : '' }" @blur="validateLoginField('adminUsername', adminUsername)" @input="clearLoginErr('adminUsername')" placeholder="管理员账号" @keydown.enter="focusAdminPwd">
            <div v-if="loginErrors.adminUsername" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.adminUsername }}</div>
          </div>
          <div class="form-group">
            <label>密码</label>
            <div v-if="loginErrors.adminPassword" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.adminPassword }}</div>
            <input ref="adminPwdRef" v-model="adminPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleAdminLogin">
          </div>
          <button class="login-submit login-submit--amber" :disabled="loading" @click="handleAdminLogin">{{ loading ? '登录中...' : '登录' }}</button>
        </div>

        <div v-if="loginType === 'parent'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="parentUsername" class="form-input" :style="{ borderColor: loginErrors.parentUsername ? '#f87171' : '' }" @blur="validateLoginField('parentUsername', parentUsername)" @input="clearLoginErr('parentUsername')" placeholder="家长账号" @keydown.enter="focusParentPwd">
            <div v-if="loginErrors.parentUsername" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.parentUsername }}</div>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input ref="parentPwdRef" v-model="parentPassword" type="password" class="form-input" :style="{ borderColor: loginErrors.parentPassword ? '#f87171' : '' }" @blur="validateLoginField('parentPassword', parentPassword)" @input="clearLoginErr('parentPassword')" placeholder="输入密码" @keydown.enter="handleParentLogin">
            <div v-if="loginErrors.parentPassword" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.parentPassword }}</div>
          </div>
          <button class="login-submit login-submit--green" :disabled="loading" @click="handleParentLogin">{{ loading ? '登录中...' : '登录' }}</button>
        </div>

        <!-- 班级码登录（学生端/大屏入口） -->
        <div v-if="loginType === 'class'" class="login-form">
          <div v-if="classInfo" class="class-login-success">
            <div class="success-icon">✅</div>
            <h3>欢迎进入 {{ classInfo.class_name }}！</h3>
            <p>{{ classInfo.student_count }} 位同学</p>
            <p class="success-hint">即将进入班级大屏...</p>
          </div>
          <template v-else>
              <input v-model="classCode" class="form-input" placeholder="请向班主任获取班级码" maxlength="12" autocomplete="off" :style="{ borderColor: loginErrors.classCode ? '#f87171' : '' }" @blur="validateLoginField('classCode', classCode)" @input="clearLoginErr('classCode')" @keydown.enter="handleClassLogin">
                maxlength="12"
            <div v-if="loginErrors.classCode" style="color:#f87171;font-size:11px;margin-top:2px;">{{ loginErrors.classCode }}</div>
                autocomplete="off"
                @keydown.enter="handleClassLogin"
              />
            </div>
            <p class="input-hint" style="font-size:12px;color:var(--color-text-secondary);margin:-8px 0 0;">如 LS301</p>
            <div v-if="classCodeError" class="error-msg" style="color:#EF4444;font-size:13px;padding:8px 12px;background:rgba(239,68,68,0.08);border-radius:8px;">{{ classCodeError }}</div>
            <button class="login-submit login-submit--purple" :disabled="loading || classCode.length < 3" @click="handleClassLogin">
              {{ loading ? '⏳ 验证中...' : '🚀 进入班级' }}
            </button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  min-height: 100vh;
  padding-top: 60px;
  background: #FFFFFF;
  color: #1D1D1F;
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'SF Pro Display', sans-serif;
  -webkit-font-smoothing: antialiased;
  position: relative;
}
.login-topnav {
  position: fixed; top: 0; left: 0; right: 0; z-index: 50;
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 32px; max-width: 1100px; margin: 0 auto; width: 100%; box-sizing: border-box;
  background: rgba(255,255,255,.8); backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0,0,0,.04);
}
.topnav-link { color: #AEAEB2; font-size: 13px; text-decoration: none; transition: color .2s; margin-left: auto; }
.topnav-link:hover { color: #6E6E73; }
.intro {
  flex: 1;
  max-width: 680px;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 64px;
  background: #FBFBFD;
}
.intro-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.25;
  pointer-events: none;
}
.intro-orb--top {
  width: 500px; height: 500px;
  background: #C7D2FE;
  top: -200px; right: -100px;
  animation: orbFloat 12s ease-in-out infinite;
}
.intro-orb--bottom {
  width: 350px; height: 350px;
  background: #A7F3D0;
  bottom: -100px; left: -50px;
  animation: orbFloat 15s ease-in-out infinite reverse;
}
@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(40px, -30px); }
}
.intro-content { position: relative; z-index: 1; max-width: 480px; width: 100%; }
.intro-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid #E5E5EA;
  border-radius: 9999px;
  padding: 6px 16px;
  font-size: 13px;
  color: #6E6E73;
  margin-bottom: 28px;
}
.intro-badge-dot {
  width: 6px; height: 6px;
  background: #34C759;
  border-radius: 50%;
  display: inline-block;
}
.intro-icon { font-size: 56px; margin-bottom: 20px; }
.intro-title {
  font-size: 40px;
  font-weight: 900;
  color: #1D1D1F;
  line-height: 1.2;
  letter-spacing: -1px;
  margin-bottom: 20px;
}
.intro-title-grad {
  background: linear-gradient(135deg, #5E5CE6, #FF375F, #FF9F0A);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.intro-desc { font-size: 16px; color: #86868B; line-height: 1.7; white-space: pre-line; }
.intro-dots { display: flex; gap: 8px; margin-top: 44px; }
.intro-dot {
  height: 8px; width: 8px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  background: #D2D2D7;
  transition: all 0.3s ease;
}
.intro-dot--active {
  width: 32px;
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
}
.intro-footer { position: absolute; bottom: -80px; left: 0; display: flex; gap: 16px; align-items: center; }
.intro-footer-link {
  color: #86868B; font-size: 13px;
  text-decoration: none;
  display: flex; align-items: center; gap: 6px;
  transition: color 0.2s;
}
.intro-footer-link:hover { color: #1D1D1F; }
.intro-footer-sep { color: #AEAEB2; font-size: 13px; }
.slide-fade-enter-active { transition: all 0.4s ease; }
.slide-fade-leave-active { transition: all 0.3s ease; }
.slide-fade-enter-from { opacity: 0; transform: translateX(20px); }
.slide-fade-leave-to { opacity: 0; transform: translateX(-20px); }
.login-panel {
  flex: 0 0 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  background: #FFFFFF;
  border-left: 1px solid #F0F0F3;
  position: relative;
  z-index: 1;
}
.login-card { width: 100%; max-width: 360px; }
.login-card-header { text-align: center; margin-bottom: 28px; }
.login-card-icon  { font-size: 40px; margin-bottom: 8px; display: block; }
.login-card-title { font-size: 24px; font-weight: 800; color: #1D1D1F; }
.login-tabs {
  display: flex; gap: 0;
  margin-bottom: 24px;
  background: #F5F5F7;
  border-radius: 12px;
  border: 1px solid #E5E5EA;
  padding: 4px;
}
.login-tab {
  flex: 1;
  padding: 10px 8px;
  text-align: center;
  font-size: 14px;
  font-weight: 500;
  border-radius: 9px;
  cursor: pointer;
  border: none;
  background: transparent;
  color: #86868B;
  transition: all 0.25s ease;
}
.login-tab--active {
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
  color: #FFFFFF;
  box-shadow: 0 2px 8px rgba(94, 92, 230, 0.2);
}
.login-form { display: flex; flex-direction: column; gap: 16px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-group label { font-size: 13px; font-weight: 500; color: #6E6E73; }
.form-input {
  width: 100%;
  padding: 11px 14px;
  background: #F5F5F7;
  border: 1px solid #E5E5EA;
  border-radius: 10px;
  color: #1D1D1F;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
}
.form-input:focus {
  border-color: #5E5CE6;
  background: #FFFFFF;
  box-shadow: 0 0 0 3px rgba(94, 92, 230, 0.12);
}
.form-input::placeholder { color: #AEAEB2; }
.login-submit {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  background: linear-gradient(135deg, #5E5CE6, #818CF8);
  color: #FFFFFF;
  font-size: 15px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(94, 92, 230, 0.2);
  transition: all 0.3s ease;
}
.login-submit:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(94, 92, 230, 0.3);
}
.login-submit:disabled { opacity: 0.6; cursor: not-allowed; }
.login-submit--amber { background: linear-gradient(135deg, #F59E0B, #FBBF24); box-shadow: 0 4px 14px rgba(245, 158, 11, 0.2); }
.login-submit--green { background: linear-gradient(135deg, #10B981, #34D399); box-shadow: 0 4px 14px rgba(16, 185, 129, 0.2); }
.login-social { margin-top: 24px; text-align: center; }
.login-social-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #AEAEB2;
  font-size: 12px;
  margin-bottom: 14px;
}
.login-social-line { flex: 1; height: 1px; background: #E5E5EA; }
.login-social-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.login-social-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  padding: 10px 4px;
  background: #F5F5F7;
  border: 1px solid #E5E5EA;
  border-radius: 10px;
  color: #6E6E73;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}
.login-social-btn:hover {
  background: #FFFFFF;
  border-color: #D2D2D7;
  color: #1D1D1F;
}
.login-social-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px; height: 28px;
  border-radius: 8px;
  font-size: 14px;
  color: #FFFFFF;
}
@media (max-width: 768px) {
  .intro { display: none !important; }
  .login-panel { flex: 1; }
}
</style>
