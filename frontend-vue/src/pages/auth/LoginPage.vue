<script setup lang="ts">
import { ref, reactive, computed, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet, apiPost } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import type { ApiResponse, User } from '@/types'

const props = withDefaults(defineProps<{ initialRole?: string; mode?: string }>(), { initialRole: 'teacher', mode: 'account' })
const router = useRouter()
const authStore = useAuthStore()

const loginType = ref<'teacher' | 'admin' | 'class'>((props.mode === 'code' ? 'class' : props.initialRole) as 'teacher' | 'admin' | 'class')
const teacherUsername = ref('')
const teacherLoginError = ref('')
const teacherPassword = ref('')
const adminUsername = ref('')
const adminPassword = ref('')
const classCode = ref('')
const classCodeError = ref('')
const classInfo = ref<{ class_name: string; student_count: number } | null>(null)
const loading = ref(false)
const loginStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
// 登录过期标记：从 /login?expired=1 读取，登录卡片顶部内联提示
const sessionExpired = ref(false)

/** 提交按钮四态走全局语义类（取代原先内联写死的 hex 底色，暗色下自动翻转） */
const submitStateClass = computed(() => `btn-state-${loginStatus.value}`)

// 内联校验
const loginErrors = reactive<Record<string, string>>({})
function clearLoginErr(f: string) { delete loginErrors[f] }
function onClassCodeInput(e: Event) {
  const input = e.target as HTMLInputElement
  classCode.value = input.value.replace(/[^0-9A-Za-z]/g, '').toUpperCase().slice(0, 8)
  clearLoginErr('classCode')
}
function validateLoginField(field: string, val: string): boolean {
  if (field === 'teacherUsername' && !val.trim()) { loginErrors.teacherUsername = '请输入教师账号'; return false }
  if (field === 'teacherPassword' && !val.trim()) { loginErrors.teacherPassword = '请输入教师登录密码'; return false }
  if (field === 'adminUsername' && !val.trim()) { loginErrors.adminUsername = '请输入管理员账号'; return false }
  if (field === 'adminPassword' && !val.trim()) { loginErrors.adminPassword = '请输入管理员密码'; return false }
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
  } else if (type === 'class') {
    if (!classCode.value.trim()) loginErrors.classCode = '请输入班级码'
  }
  return Object.keys(loginErrors).length === 0
}

const teacherPwdRef = ref<HTMLInputElement>()
const adminPwdRef = ref<HTMLInputElement>()

let teacherAttempts = 0
const MAX_ATTEMPTS = 3

function focusTeacherPwd() { if (teacherUsername.value.trim()) nextTick(() => teacherPwdRef.value?.focus()) }
function focusAdminPwd() { if (adminUsername.value.trim()) nextTick(() => adminPwdRef.value?.focus()) }

async function handleTeacherLogin() {
  if (!validateLoginForm('teacher')) return
  loading.value = true
  loginStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/auth/teacher/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: teacherUsername.value.trim(), password: teacherPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) { throw { response: { data } } }
    teacherAttempts = 0
    teacherLoginError.value = ''
    authStore.setAuth(data.data.token, data.data.user)
    loginStatus.value = 'success'
    setTimeout(() => router.replace({ name: 'teacher-dashboard' }), 800)
  } catch (e: any) {
    teacherAttempts++
    const remaining = MAX_ATTEMPTS - teacherAttempts
    loginStatus.value = 'error'
    teacherLoginError.value = teacherAttempts >= MAX_ATTEMPTS
      ? '密码错误次数过多，请联系管理员'
      : ((e?.response?.data?.message || '账号或密码错误') + `，还剩 ${remaining} 次`)
    setTimeout(() => { if (loginStatus.value === 'error') loginStatus.value = 'idle' }, 3000)
  } finally { loading.value = false; if (loginStatus.value === 'loading') loginStatus.value = 'idle' }
}

async function handleAdminLogin() {
  if (!validateLoginForm('admin')) return
  loading.value = true
  loginStatus.value = 'loading'
  try {
    const res = await apiPost<ApiResponse<{ token: string; user: User }>>('/api/v1/auth/admin/login', {
      username: adminUsername.value.trim(), password: adminPassword.value,
    })
    authStore.setAuth(res.data.token, res.data.user)
    loginStatus.value = 'success'
    setTimeout(() => router.replace({ name: 'admin-dashboard' }), 800)
  } catch { loginStatus.value = 'error'; setTimeout(() => { if (loginStatus.value === 'error') loginStatus.value = 'idle' }, 3000) } finally { loading.value = false; if (loginStatus.value === 'loading') loginStatus.value = 'idle' }
}

async function handleClassLogin() {
  if (!validateLoginForm('class')) return
  loading.value = true
  loginStatus.value = 'loading'
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
    loginStatus.value = 'success'
    setTimeout(() => router.push({ name: 'teacher-dashboard-basic' }), 1500)
  } catch (e: any) {
    loginStatus.value = 'error'
    classCodeError.value = e?.response?.data?.message || '班级码无效，请核对后重试'
    setTimeout(() => { if (loginStatus.value === 'error') loginStatus.value = 'idle' }, 3000)
  } finally { loading.value = false; if (loginStatus.value === 'loading') loginStatus.value = 'idle' }
}

// 第三方登录平台：后台勾选后由 /auth/third-party/options 返回，未配置时默认办公三平台
interface ThirdPartyOption { key: string; label: string; icon?: string; color?: string }
const platforms = ref<ThirdPartyOption[]>([])
const thirdPartyError = ref('')

async function loadThirdPartyOptions() {
  try {
    const res = await apiGet<{ data: ThirdPartyOption[] }>('/api/v1/auth/third-party/options')
    platforms.value = res.data || []
  } catch { platforms.value = [] }
}

async function handleThirdPartyLogin(platform: string) {
  const label = platforms.value.find(p => p.key === platform)?.label || platform

  // 仅办公平台扫码（企业微信/钉钉/飞书）接入统一授权入口，其他平台尚未接入 OAuth 流程
  if (!['wechat_work', 'dingtalk', 'feishu'].includes(platform)) {
    thirdPartyError.value = `「${label}」暂未接入扫码登录，请使用账号密码登录`
    setTimeout(() => { thirdPartyError.value = '' }, 3000)
    return
  }

  thirdPartyError.value = ''
  try {
    const res = await apiGet<{ data: { auth_url: string } }>('/api/v1/auth/third-party/auth-url', {
      params: { redirect_uri: window.location.origin + '/auth/callback' },
    })
    const w = 600, h = 500
    const left = (screen.width - w) / 2
    const top = (screen.height - h) / 2
    window.open(res.data.auth_url, platform,
      `width=${w},height=${h},left=${left},top=${top},menubar=no,toolbar=no,status=no,scrollbars=yes`)
  } catch (e: any) {
    thirdPartyError.value = e?.response?.data?.message || '未配置第三方平台，请在后台学校设置中选择'
    setTimeout(() => { thirdPartyError.value = '' }, 3000)
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
    desc: '积分规则 宠物进化 排行榜 通知公告\n考勤 作业 答题 商城 广播 AI',
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
      authStore.setAuth(d.token, d.user)
      loginStatus.value = 'success'
      setTimeout(() => router.replace({ name: 'teacher-dashboard' }), 800)
    }
  } catch { /* handled */ } finally { loading.value = false }
}

onMounted(() => {
  // 从 /login?expired=1 读取登录过期标记，卡片顶部内联提示
  sessionExpired.value = router.currentRoute.value.query.expired === '1'
  loadThirdPartyOptions()
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
            v-for="(_, i) in slides" :key="i"
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
          <h1 class="login-card-title">学宠星球</h1>
        </div>
        <div v-if="sessionExpired" class="session-expired-banner">⏰ 登录已过期，请重新登录</div>
        <div class="login-tabs">
          <button
            v-for="t in (['teacher', 'admin', 'class'] as const)" :key="t"
            :class="['login-tab', { 'login-tab--active': loginType === t }]"
            @click="loginType = t"
          >{{ t === 'teacher' ? '教师' : t === 'admin' ? '管理员' : '🔑 班级码' }}</button>
        </div>

        <div v-if="loginType === 'teacher'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="teacherUsername" class="form-input" :class="{ 'is-error': loginErrors.teacherUsername }" @blur="validateLoginField('teacherUsername', teacherUsername)" @input="clearLoginErr('teacherUsername')" placeholder="教师账号" @keydown.enter="focusTeacherPwd">
            <div v-if="loginErrors.teacherUsername" class="field-error">{{ loginErrors.teacherUsername }}</div>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input ref="teacherPwdRef" v-model="teacherPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleTeacherLogin">
            <div v-if="loginErrors.teacherPassword" class="field-error">{{ loginErrors.teacherPassword }}</div>
          </div>
          <button class="login-submit" :class="submitStateClass" :disabled="loginStatus === 'loading'" @click="handleTeacherLogin">
            <span v-if="loginStatus === 'idle'">🚀 登录</span>
            <span v-else-if="loginStatus === 'loading'">⏳ 登录中...</span>
            <span v-else-if="loginStatus === 'success'">✅ 登录成功</span>
            <span v-else>❌ 登录失败</span>
          </button>
          <div v-if="teacherLoginError" class="lgp-error">{{ teacherLoginError }}</div>
          <div class="login-social">
            <div class="login-social-label"><span class="login-social-line"></span> 扫码登录 <span class="login-social-line"></span></div>
            <div v-if="platforms.length === 0" class="lgp-muted-12">暂无可用的扫码登录方式</div>
            <div v-else class="login-social-grid">
              <button v-for="p in platforms" :key="p.key" class="login-social-btn" @click="handleThirdPartyLogin(p.key)">
                <span class="login-social-icon"><PlatformIcon :platform="p.key" :size="24" /></span>
                {{ p.label }}
              </button>
            </div>
            <div v-if="thirdPartyError" class="lgp-error">{{ thirdPartyError }}</div>
          </div>
        </div>

        <div v-if="loginType === 'admin'" class="login-form">
          <div class="form-group">
            <label>账号</label>
            <input v-model="adminUsername" class="form-input" :class="{ 'is-error': loginErrors.adminUsername }" @blur="validateLoginField('adminUsername', adminUsername)" @input="clearLoginErr('adminUsername')" placeholder="管理员账号" @keydown.enter="focusAdminPwd">
            <div v-if="loginErrors.adminUsername" class="field-error">{{ loginErrors.adminUsername }}</div>
          </div>
          <div class="form-group">
            <label>密码</label>
            <div v-if="loginErrors.adminPassword" class="field-error">{{ loginErrors.adminPassword }}</div>
            <input ref="adminPwdRef" v-model="adminPassword" type="password" class="form-input" placeholder="输入密码" @keydown.enter="handleAdminLogin">
          </div>
          <button class="login-submit" :class="submitStateClass" :disabled="loginStatus === 'loading'" @click="handleAdminLogin">
            <span v-if="loginStatus === 'idle'">🚀 登录</span>
            <span v-else-if="loginStatus === 'loading'">⏳ 登录中...</span>
            <span v-else-if="loginStatus === 'success'">✅ 登录成功</span>
            <span v-else>❌ 登录失败</span>
          </button>
        </div>

        <!-- 班级码登录（教室端入口） -->
        <div v-if="loginType === 'class'" class="login-form">
          <div v-if="classInfo" class="class-login-success">
            <div class="success-icon">✅</div>
            <h3>欢迎进入 {{ classInfo.class_name }}！</h3>
            <p>{{ classInfo.student_count }} 位同学</p>
            <p class="success-hint">即将进入班级大屏...</p>
          </div>
          <template v-else>
            <div class="form-group">
              <input :value="classCode" class="form-input form-input--code" placeholder="输入班级码（如 LS11）" maxlength="8" autocomplete="off" :class="{ 'is-error': loginErrors.classCode }" @input="onClassCodeInput" @blur="validateLoginField('classCode', classCode)" @keydown.enter="handleClassLogin">
              <div v-if="loginErrors.classCode" class="field-error">{{ loginErrors.classCode }}</div>
            </div>
            <p class="input-hint lgp-muted-12-tight">如 LS11（一年级1班）</p>
            <div v-if="classCodeError" class="error-msg lgp-error-inline">{{ classCodeError }}</div>
            <button class="login-submit" :class="submitStateClass" :disabled="loginStatus === 'loading' || classCode.length < 3" @click="handleClassLogin">
              <span v-if="loginStatus === 'idle'">🚀 进入班级</span>
              <span v-else-if="loginStatus === 'loading'">⏳ 验证中...</span>
              <span v-else-if="loginStatus === 'success'">✅ 欢迎进入</span>
              <span v-else>❌ 验证失败</span>
            </button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ===================================================================
 * 登录页 — 视觉对标重绘（2026-09-19）
 * 去掉：光晕 orb、渐变文字、胶囊 tab + 渐变激活、按钮内联状态色 hex、
 *       逐条手写的 html.dark 覆盖层（现由 --ui-* 令牌自动翻转）。
 * 统一：三个身份的主操作都走品牌实色（唯一 accent），身份由 tab 表达，不再靠按钮颜色区分。
 * =================================================================== */
.login-page {
  display: flex;
  justify-content: center;
  min-height: 100vh;
  padding-top: 56px;
  background: var(--ui-bg);
  color: var(--ui-fg);
  position: relative;
  box-sizing: border-box;
}

.login-topnav {
  position: fixed; top: 0; left: 0; right: 0; z-index: 50;
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 28px; max-width: 1120px; margin: 0 auto; width: 100%; box-sizing: border-box;
  background: color-mix(in srgb, var(--ui-bg) 86%, transparent);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--ui-border);
}
.topnav-link {
  display: inline-flex; align-items: center; gap: 6px;
  color: var(--ui-fg-subtle); font-size: 13px; transition: color .15s; margin-left: auto;
}
.topnav-link:hover { color: var(--ui-fg); }

/* 左侧介绍区 */
.intro {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 56px;
  background: var(--ui-bg-subtle);
  position: relative;
  overflow: hidden;
}
/* 唯一保留的氛围：静态品牌色晕染（无动画、无 blur 圆斑） */
.intro::before {
  content: '';
  position: absolute; inset: 0; pointer-events: none;
  background: radial-gradient(560px 380px at 82% -10%, var(--ui-brand-soft), transparent 70%);
}
.intro-content { position: relative; z-index: 1; max-width: 460px; width: 100%; }

.intro-badge {
  display: inline-flex; align-items: center; gap: 7px;
  height: 26px; padding: 0 11px;
  background: var(--ui-card); border: 1px solid var(--ui-border);
  border-radius: var(--ui-r-sm);
  font-size: 12.5px; color: var(--ui-fg-muted); margin-bottom: 24px;
}
.intro-badge-dot { width: 6px; height: 6px; background: var(--c-green); border-radius: 50%; display: inline-block; flex-shrink: 0; }
.intro-icon {
  width: 44px; height: 44px; border-radius: var(--ui-r-lg); margin-bottom: 18px;
  display: flex; align-items: center; justify-content: center;
  background: var(--ui-brand-soft); color: var(--ui-brand);
  border: 1px solid var(--ui-brand-border);
}
.intro-title {
  font-size: 34px; font-weight: 700; color: var(--ui-fg);
  line-height: 1.24; letter-spacing: -0.025em; margin-bottom: 16px;
}
/* 高亮词用纯色品牌，不再做渐变文字裁剪 */
.intro-title-accent { color: var(--ui-brand); }
.intro-desc { font-size: 14.5px; color: var(--ui-fg-muted); line-height: 1.75; white-space: pre-line; }
.intro-dots { display: flex; gap: 6px; margin-top: 36px; }
.intro-dot {
  height: 6px; width: 6px; border-radius: 3px; border: none; padding: 0;
  cursor: pointer; background: var(--ui-border-strong); transition: all 0.25s var(--ease-smooth);
}
.intro-dot:hover { background: var(--ui-fg-subtle); }
.intro-dot--active { width: 24px; background: var(--ui-brand); }
.intro-footer { display: flex; gap: 14px; align-items: center; margin-top: 40px; }
.intro-footer-link {
  display: inline-flex; align-items: center; gap: 6px;
  color: var(--ui-fg-muted); font-size: 13px; transition: color 0.15s;
}
.intro-footer-link:hover { color: var(--ui-fg); }
.intro-footer-sep { color: var(--ui-fg-subtle); font-size: 13px; }

.slide-fade-enter-active { transition: all 0.35s ease; }
.slide-fade-leave-active { transition: all 0.25s ease; }
.slide-fade-enter-from { opacity: 0; transform: translateX(16px); }
.slide-fade-leave-to { opacity: 0; transform: translateX(-16px); }

/* 右侧登录区 */
.login-panel {
  flex: 0 0 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  background: var(--ui-bg);
  position: relative;
  z-index: 1;
}
.login-card { width: 100%; max-width: 356px; }
.login-card-header { display: flex; align-items: center; justify-content: center; gap: 10px; margin-bottom: 24px; }
/* 品牌标识：纯色方块 + 字重，不用渐变裁剪 */
.login-card-mark {
  width: 32px; height: 32px; border-radius: var(--ui-r-md);
  background: var(--ui-brand); color: var(--ui-brand-fg);
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.login-card-title { font-size: 21px; font-weight: 650; letter-spacing: -0.02em; color: var(--ui-fg); }

.session-expired-banner {
  display: flex; align-items: center; justify-content: center; gap: 6px;
  margin: 0 0 14px; padding: 9px 12px;
  background: var(--c-amber-bg);
  border: 1px solid var(--c-amber-chip);
  border-radius: var(--ui-r-md);
  color: var(--color-warning-text);
  font-size: 12.5px; font-weight: 500;
}

/* 身份切换：分段控件（底槽 --ui-muted，选中项白底 + 微阴影） */
.login-tabs {
  display: flex; gap: 2px;
  margin-bottom: 22px;
  background: var(--ui-muted);
  border-radius: var(--ui-r-lg);
  padding: 3px;
}
.login-tab {
  flex: 1;
  height: 32px;
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  font-size: 13px; font-weight: 500; font-family: inherit;
  border-radius: calc(var(--ui-r-lg) - 3px);
  cursor: pointer;
  border: none;
  background: transparent;
  color: var(--ui-fg-muted);
  transition: all 0.15s var(--ease-smooth);
}
.login-tab:hover { color: var(--ui-fg); }
.login-tab--active {
  background: var(--ui-card);
  color: var(--ui-fg);
  font-weight: 600;
  box-shadow: var(--ui-shadow-xs);
}

.login-form { display: flex; flex-direction: column; gap: 14px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-group label { font-size: 12.5px; font-weight: 550; color: var(--ui-fg); }
/* .form-input 外观走全局令牌；此处只补登录页专有的错误态与班级码大字态 */
.form-input.is-error { border-color: var(--c-red); }
.form-input.is-error:focus { border-color: var(--c-red); box-shadow: 0 0 0 3px var(--c-red-bg); }
.form-input--code {
  height: 48px; text-align: center;
  font-size: 19px; font-weight: 650;
  letter-spacing: 0.14em; text-indent: 0.14em;
}
.form-input--code::placeholder { font-size: 14px; font-weight: 400; letter-spacing: 0; text-indent: 0; }

/* 提交按钮：外观全部由全局 .btn-state-* 提供，此处只管尺寸与排版 */
.login-submit {
  width: 100%;
  height: 42px;
  margin-top: 4px;
  border-radius: var(--ui-r-md);
  display: inline-flex; align-items: center; justify-content: center; gap: 7px;
  font-size: 14.5px;
  font-weight: 600;
  cursor: pointer;
  transition: filter 0.15s, background 0.15s, transform 0.08s;
}
.login-submit:hover:not(:disabled) { filter: brightness(1.06); }
.login-submit:active:not(:disabled) { transform: scale(0.99); }
.login-submit:disabled { opacity: 0.55; cursor: not-allowed; }

.login-social { margin-top: 20px; text-align: center; }
.login-social-label {
  display: flex; align-items: center; gap: 8px;
  color: var(--ui-fg-subtle); font-size: 11.5px; margin-bottom: 12px;
}
.login-social-line { flex: 1; height: 1px; background: var(--ui-border); }
.login-social-grid { display: flex; flex-wrap: wrap; justify-content: center; gap: 8px; }
.login-social-btn {
  display: flex; flex-direction: column; align-items: center; gap: 5px;
  padding: 10px 12px; min-width: 72px;
  background: var(--ui-card); border: 1px solid var(--ui-border);
  border-radius: var(--ui-r-md);
  color: var(--ui-fg-muted); font-size: 11px; font-family: inherit;
  cursor: pointer; transition: 0.15s;
}
.login-social-btn:hover { background: var(--ui-muted); border-color: var(--ui-border-strong); color: var(--ui-fg); }
.login-social-icon { display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; }

/* 班级码登录成功 */
.class-login-success { text-align: center; padding: 12px 0 4px; }
.success-icon { color: var(--c-green); display: inline-flex; margin-bottom: 10px; }
.class-login-success h3 { font-size: 16px; font-weight: 650; letter-spacing: -0.01em; color: var(--ui-fg); }
.class-login-success p { font-size: 13px; color: var(--ui-fg-muted); margin-top: 4px; }
.success-hint { color: var(--ui-fg-subtle) !important; }

/* 页内通用小件 */
.form-error-banner {
  margin-top: 10px; padding: 8px 12px;
  background: var(--c-red-bg); border: 1px solid var(--c-red-border);
  border-radius: var(--ui-r-md);
  color: var(--color-danger-text); font-size: 12.5px;
}
.muted-note { font-size: 12.5px; color: var(--ui-fg-muted); padding: 8px 0; }
.input-hint { font-size: 11.5px; color: var(--ui-fg-subtle); margin-top: -6px; }

@media (max-width: 860px) {
  .intro { display: none; }
  .login-panel { flex: 1; }
}
</style>
