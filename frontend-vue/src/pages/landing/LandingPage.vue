<script setup lang="ts">
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { apiPost } from '@/utils/api'
import type { ApiResponse, User } from '@/types'
import { slides, features, stages, platforms } from './landingData'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

if (authStore.isLoggedIn) {
  if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
  else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
  else if (authStore.isParent) router.replace({ name: 'parent-home' })
}

const currentSlide = ref(0)
let slideTimer: number = 0
onMounted(function () {
  slideTimer = setInterval(function () {
    currentSlide.value = (currentSlide.value + 1) % slides.length
  }, 5000)
})
onUnmounted(function () { clearInterval(slideTimer) })
function goToSlide(i: number) {
  currentSlide.value = i
  clearInterval(slideTimer)
  slideTimer = setInterval(function () {
    currentSlide.value = (currentSlide.value + 1) % slides.length
  }, 5000)
}

const loginType = ref('teacher')
const teacherUsername = ref('')
const teacherPassword = ref('')
const adminUsername = ref('')
const adminPassword = ref('')
const parentUsername = ref('')
const parentPassword = ref('')
const loading = ref(false)

const teacherPwdRef = ref()
const adminPwdRef = ref()
const parentPwdRef = ref()

let teacherAttempts = 0
const MAX_ATTEMPTS = 3

function focusTeacherPwd() { if (teacherUsername.value.trim()) nextTick(function () { teacherPwdRef.value && teacherPwdRef.value.focus() }) }
function focusAdminPwd() { if (adminUsername.value.trim()) nextTick(function () { adminPwdRef.value && adminPwdRef.value.focus() }) }
function focusParentPwd() { if (parentUsername.value.trim()) nextTick(function () { parentPwdRef.value && parentPwdRef.value.focus() }) }

async function handleTeacherLogin() {
  if (!teacherUsername.value.trim() || !teacherPassword.value) { toast.show('请输入账号和密码', 'error'); return }
  loading.value = true
  try {
    const res = await apiPost('/api/auth/teacher/login', { username: teacherUsername.value.trim(), password: teacherPassword.value })
    teacherAttempts = 0; authStore.setAuth(res.data.token, res.data.user); toast.show('登录成功', 'success'); router.replace({ name: 'teacher-dashboard' })
  } catch (e) {
    teacherAttempts++; const remaining = MAX_ATTEMPTS - teacherAttempts
    if (teacherAttempts &gt;= MAX_ATTEMPTS) toast.show('密码错误次数过多，请联系管理员', 'error')
    else { toast.show('账号或密码错误，还剩 ' + remaining + ' 次', 'error') }
  } finally { loading.value = false }
}

async function handleAdminLogin() {
  if (!adminUsername.value.trim() || !adminPassword.value) { toast.show('请输入账号和密码', 'error'); return }
  loading.value = true
  try {
    const res = await apiPost('/api/auth/admin/login', { username: adminUsername.value.trim(), password: adminPassword.value })
    authStore.setAuth(res.data.token, res.data.user); toast.show('登录成功', 'success'); router.replace({ name: 'admin-dashboard' })
  } catch (e) { } finally { loading.value = false }
}

async function handleParentLogin() {
  if (!parentUsername.value.trim() || !parentPassword.value) { toast.show('请输入账号和密码', 'error'); return }
  loading.value = true
  try {
    const res = await apiPost('/api/auth/parent/login', { username: parentUsername.value.trim(), password: parentPassword.value })
    authStore.setAuth(res.data.token, res.data.user); toast.show('登录成功', 'success'); router.replace({ name: 'parent-home' })
  } catch (e) { } finally { loading.value = false }
}

function handleThirdPartyLogin(platform: string) {
  var p = platforms.find(function (x) { return x.key === platform })
  toast.show('正在打开' + (p ? p.label : platform) + '扫码...', 'success')
}
</script>

<template>
  <div class="layout">
    <div class="left">
      <div class="left-glow left-glow--top"></div>
      <div class="left-glow left-glow--bottom"></div>
      <div class="left-content">
        <transition name="fade" mode="out-in">
          <div :key="currentSlide" class="slide">
            <div class="badge">
              <span class="badge-dot"></span>
              {{ slides[currentSlide].badge }}
            </div>
            <div class="slide-icon">{{ slides[currentSlide].icon }}</div>
            <h1 class="slide-title">
              {{ slides[currentSlide].title }}
              <br>
              <span class="slide-gradient">{{ slides[currentSlide].highlight }}</span>
            </h1>
            <p class="slide-desc">{{ slides[currentSlide].desc }}</p>
          </div>
        </transition>
        <div class="dots">
          <button
            v-for="(s, idx) in slides"
            :key="idx"
            :class="['dot', { 'dot--active': currentSlide === idx }]"
            @click="goToSlide(idx)"
          ></button>
        </div>
        <div class="left-footer">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="left-link">GitHub</a>
          <span class="left-sep">MIT</span>
        </div>
      </div>
      <div class="features">
        <h2 class="features-title">12 大功能模块</h2>
        <p class="features-sub">覆盖班级管理全场景，全部免费</p>
        <div class="feature-grid">
          <div v-for="f in features" :key="f.title" class="feat">
            <span class="feat-icon">{{ f.icon }}</span>
            <div>
              <div class="feat-name">{{ f.title }}</div>
              <div class="feat-desc">{{ f.desc }}</div>
            </div>
          </div>
        </div>
        <h2 class="features-title" style="margin-top:48px">11 阶宠物进化</h2>
        <p class="features-sub">积分变经验，从星尘到银河</p>
        <div class="evo">
          <div
            v-for="(s, idx) in stages"
            :key="s.name"
            style="display:flex;align-items:center;gap:2px"
          >
            <span v-if="idx !== 0" class="evo-arrow">→</span>
            <div class="evo-item">
              <span class="evo-emoji">{{ s.emoji }}</span>
              <span class="evo-name">{{ s.name }}</span>
            </div>
          </div>
        </div>
        <div class="cta-box">
          <div class="cta-icon">🚀</div>
          <h3>加入开源社区</h3>
          <p>项目完全开源，欢迎 Star、Fork、Issue、PR</p>
          <div class="cta-links">
            <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="cta-btn cta-btn--dark">Star</a>
            <a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" class="cta-btn cta-btn--ghost">Issue</a>
          </div>
        </div>
        <footer class="foot">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank">学趣星球</a>
          <span> · MIT · 自托管 · 完全免费</span>
        </footer>
      </div>
    </div>
    <div class="right">
      <div class="login-card">
        <div class="login-header">
          <span class="login-icon">🌌</span>
          <h1 class="login-brand">学趣星球</h1>
        </div>
        <div class="login-tabs">
          <button
            v-for="t in (['teacher', 'admin', 'parent'] as const)"
            :key="t"
            :class="['login-tab', { 'login-tab--active': loginType === t }]"
            @click="loginType = t"
          >{{ t === 'teacher' ? '教师' : t === 'admin' ? '管理员' : '家长' }}</button>
        </div>
        <div v-if="loginType === 'teacher'" class="login-form">
          <div class="field">
            <label>账号</label>
            <input v-model="teacherUsername" placeholder="教师账号" @keydown.enter="focusTeacherPwd">
          </div>
          <div class="field">
            <label>密码</label>
            <input ref="teacherPwdRef" v-model="teacherPassword" type="password" placeholder="输入密码" @keydown.enter="handleTeacherLogin">
          </div>
          <button class="btn-login" :disabled="loading" @click="handleTeacherLogin">{{ loading ? '...' : '登录' }}</button>
          <div class="social">
            <div class="social-label">
              <span></span> 扫码登录 <span></span>
            </div>
            <div class="social-btns">
              <button v-for="p in platforms" :key="p.key" @click="handleThirdPartyLogin(p.key)">
                <span :style="{ background: p.color }">{{ p.icon }}</span>
                {{ p.label }}
              </button>
            </div>
          </div>
        </div>
        <div v-if="loginType === 'admin'" class="login-form">
          <div class="field">
            <label>账号</label>
            <input v-model="adminUsername" placeholder="管理员账号" @keydown.enter="focusAdminPwd">
          </div>
          <div class="field">
            <label>密码</label>
            <input ref="adminPwdRef" v-model="adminPassword" type="password" placeholder="输入密码" @keydown.enter="handleAdminLogin">
          </div>
          <button class="btn-login btn-login--amber" :disabled="loading" @click="handleAdminLogin">{{ loading ? '...' : '登录' }}</button>
        </div>
        <div v-if="loginType === 'parent'" class="login-form">
          <div class="field">
            <label>账号</label>
            <input v-model="parentUsername" placeholder="家长账号" @keydown.enter="focusParentPwd">
          </div>
          <div class="field">
            <label>密码</label>
            <input ref="parentPwdRef" v-model="parentPassword" type="password" placeholder="输入密码" @keydown.enter="handleParentLogin">
          </div>
          <button class="btn-login btn-login--green" :disabled="loading" @click="handleParentLogin">{{ loading ? '...' : '登录' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
*{box-sizing:border-box}
.layout{display:flex;min-height:100vh;background:#fff;color:#1d1d1f;font-family:'Noto Sans SC',-apple-system,BlinkMacSystemFont,'SF Pro Display',sans-serif;-webkit-font-smoothing:antialiased}
.left{flex:1;max-height:100vh;overflow-y:auto;position:relative;background:#fbfbfd}
.left::-webkit-scrollbar{width:4px}
.left::-webkit-scrollbar-thumb{background:#d2d2d7;border-radius:4px}
.left-glow{position:fixed;border-radius:50%;filter:blur(100px);opacity:.2;pointer-events:none;z-index:0}
.left-glow--top{width:500px;height:500px;background:#c7d2fe;top:-200px;right:-100px;animation:glow 12s ease-in-out infinite}
.left-glow--bottom{width:350px;height:350px;background:#a7f3d0;bottom:-100px;left:-50px;animation:glow 15s ease-in-out infinite reverse}
@keyframes glow{0%,100%{transform:translate(0,0)}50%{transform:translate(40px,-30px)}}
.left-content{position:relative;z-index:1;min-height:100vh;display:flex;flex-direction:column;align-items:center;justify-content:center;padding:60px 64px 80px}
.slide{max-width:440px;text-align:left}
.badge{display:inline-flex;align-items:center;gap:8px;background:#fff;border:1px solid #e5e5ea;border-radius:9999px;padding:6px 16px;font-size:13px;color:#6e6e73;margin-bottom:28px}
.badge-dot{width:6px;height:6px;background:#34c759;border-radius:50%;display:inline-block;box-shadow:0 0 6px rgba(52,199,89,.3)}
.slide-icon{font-size:56px;margin-bottom:20px}
.slide-title{font-size:40px;font-weight:900;line-height:1.2;letter-spacing:-1px;color:#1d1d1f;margin-bottom:20px}
.slide-gradient{background:linear-gradient(135deg,#5e5ce6,#ff375f,#ff9f0a);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.slide-desc{font-size:16px;color:#86868b;line-height:1.7;white-space:pre-line}
.dots{display:flex;gap:8px;margin-top:44px}
.dot{height:8px;border-radius:4px;border:none;cursor:pointer;background:#d2d2d7;transition:all .3s;width:8px}
.dot--active{width:32px;background:linear-gradient(135deg,#5e5ce6,#818cf8)}
.left-footer{position:absolute;bottom:40px;left:64px;display:flex;gap:16px;align-items:center;z-index:1}
.left-link{color:#86868b;font-size:13px;text-decoration:none;transition:color .2s}
.left-link:hover{color:#1d1d1f}
.left-sep{color:#aeaeb2;font-size:13px}
.fade-enter-active{transition:all .4s ease}
.fade-leave-active{transition:all .3s ease}
.fade-enter-from{opacity:0;transform:translateX(20px)}
.fade-leave-to{opacity:0;transform:translateX(-20px)}
.features{padding:0 64px 60px;position:relative;z-index:1}
.features-title{font-size:24px;font-weight:800;color:#1d1d1f;margin-bottom:8px}
.features-sub{font-size:15px;color:#86868b;margin-bottom:32px}
.feature-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:10px}
.feat{display:flex;align-items:flex-start;gap:12px;padding:16px;background:#fff;border:1px solid #f0f0f3;border-radius:14px;transition:all .2s}
.feat:hover{border-color:#e5e5ea;box-shadow:0 4px 16px rgba(0,0,0,.04);transform:translateY(-2px)}
.feat-icon{font-size:24px;flex-shrink:0;width:36px;text-align:center}
.feat-name{font-size:14px;font-weight:700;color:#1d1d1f;margin-bottom:3px}
.feat-desc{font-size:12px;color:#86868b;line-height:1.5}
.evo{display:flex;align-items:center;flex-wrap:wrap;gap:2px;padding:16px 0}
.evo-arrow{color:#c7c7cc;font-size:13px;user-select:none}
.evo-item{display:flex;flex-direction:column;align-items:center;gap:4px;padding:8px 12px;border-radius:14px;transition:all .2s;cursor:default}
.evo-item:hover{background:#f0f0f3;transform:scale(1.08)}
.evo-emoji{font-size:26px}
.evo-name{font-size:11px;color:#86868b;font-weight:500}
.cta-box{margin-top:48px;text-align:center;padding:40px;background:#fff;border:1px solid #f0f0f3;border-radius:20px}
.cta-icon{font-size:40px;margin-bottom:12px}
.cta-box h3{font-size:22px;font-weight:800;color:#1d1d1f;margin-bottom:8px}
.cta-box p{font-size:14px;color:#86868b;margin-bottom:24px}
.cta-links{display:flex;gap:10px;justify-content:center}
.cta-btn{padding:10px 24px;border-radius:9999px;font-size:14px;font-weight:600;text-decoration:none;transition:all .2s}
.cta-btn--dark{background:#1d1d1f;color:#fff}
.cta-btn--dark:hover{background:#333}
.cta-btn--ghost{border:1px solid #d2d2d7;color:#1d1d1f}
.cta-btn--ghost:hover{background:#f5f5f7}
.foot{margin-top:40px;padding-bottom:40px;text-align:center;font-size:13px;color:#aeaeb2}
.foot a{color:#86868b;text-decoration:none}
.foot a:hover{color:#1d1d1f}
.right{flex:0 0 400px;display:flex;align-items:center;justify-content:center;padding:32px;background:#fff;border-left:1px solid #f0f0f3;position:sticky;top:0;height:100vh}
.login-card{width:100%;max-width:360px}
.login-header{text-align:center;margin-bottom:28px}
.login-icon{font-size:40px;margin-bottom:8px;display:block}
.login-brand{font-size:24px;font-weight:800;color:#1d1d1f}
.login-tabs{display:flex;gap:0;margin-bottom:24px;background:#f5f5f7;border-radius:12px;border:1px solid #e5e5ea;padding:4px}
.login-tab{flex:1;padding:10px 8px;text-align:center;font-size:14px;font-weight:500;border-radius:9px;cursor:pointer;border:none;background:transparent;color:#86868b;transition:all .25s}
.login-tab--active{background:linear-gradient(135deg,#5e5ce6,#818cf8);color:#fff;box-shadow:0 2px 8px rgba(94,92,230,.2)}
.login-form{display:flex;flex-direction:column;gap:16px}
.field{display:flex;flex-direction:column;gap:4px}
.field label{font-size:13px;font-weight:500;color:#6e6e73}
.field input{width:100%;padding:11px 14px;background:#f5f5f7;border:1px solid #e5e5ea;border-radius:10px;color:#1d1d1f;font-size:14px;outline:none;transition:all .2s}
.field input:focus{border-color:#5e5ce6;background:#fff;box-shadow:0 0 0 3px rgba(94,92,230,.12)}
.field input::placeholder{color:#aeaeb2}
.btn-login{width:100%;padding:12px;border-radius:10px;background:linear-gradient(135deg,#5e5ce6,#818cf8);color:#fff;font-size:15px;font-weight:600;border:none;cursor:pointer;box-shadow:0 4px 14px rgba(94,92,230,.2);transition:all .3s}
.btn-login:hover:not(:disabled){transform:translateY(-1px);box-shadow:0 6px 20px rgba(94,92,230,.3)}
.btn-login:disabled{opacity:.6;cursor:not-allowed}
.btn-login--amber{background:linear-gradient(135deg,#f59e0b,#fbbf24);box-shadow:0 4px 14px rgba(245,158,11,.2)}
.btn-login--green{background:linear-gradient(135deg,#10b981,#34d399);box-shadow:0 4px 14px rgba(16,185,129,.2)}
.social{margin-top:24px;text-align:center}
.social-label{display:flex;align-items:center;gap:8px;color:#aeaeb2;font-size:12px;margin-bottom:14px}
.social-label span{flex:1;height:1px;background:#e5e5ea}
.social-btns{display:grid;grid-template-columns:repeat(4,1fr);gap:8px}
.social-btns button{display:flex;flex-direction:column;align-items:center;gap:5px;padding:10px 4px;background:#f5f5f7;border:1px solid #e5e5ea;border-radius:10px;color:#6e6e73;font-size:11px;cursor:pointer;transition:all .2s}
.social-btns button:hover{background:#fff;border-color:#d2d2d7;color:#1d1d1f}
.social-btns button span{display:flex;align-items:center;justify-content:center;width:28px;height:28px;border-radius:8px;font-size:14px;color:#fff}
@media(max-width:768px){
.layout{flex-direction:column}
.left{max-height:none;overflow:visible;padding:0}
.left-content{padding:80px 24px 40px;min-height:auto}
.left-footer{position:static;margin-top:32px;padding:0}
.features{padding:0 24px 40px}
.feature-grid{grid-template-columns:1fr}
.right{flex:none;position:static;height:auto;border-left:none;border-top:1px solid #f0f0f3;padding:32px 24px}
.slide-title{font-size:32px}
.slide-icon{font-size:40px}
}
</style>
