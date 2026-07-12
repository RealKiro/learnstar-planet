<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { apiPost } from '@/utils/api'
import { platforms } from './landingData'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

if (authStore.isLoggedIn) {
  if (authStore.isAdmin) { router.replace({ name: 'admin-dashboard' }) }
  else if (authStore.isTeacher) { router.replace({ name: 'teacher-dashboard' }) }
  else if (authStore.isParent) { router.replace({ name: 'parent-home' }) }
}

const loginType = ref('teacher')
const username = ref('')
const password = ref('')
const loading = ref(false)
const pwdRef = ref()

let attempts = 0
const MAX = 3

function focusPwd() {
  if (username.value.trim()) {
    nextTick(function () { pwdRef.value && pwdRef.value.focus() })
  }
}

function resetForm() {
  username.value = ''
  password.value = ''
}

async function doLogin() {
  var u = username.value.trim()
  var p = password.value
  if (!u || !p) { toast.show('请输入账号和密码', 'error'); return }

  loading.value = true
  try {
    var url = ''
    var routeName = ''
    if (loginType.value === 'teacher') {
      url = '/api/auth/teacher/login'
      routeName = 'teacher-dashboard'
    } else if (loginType.value === 'admin') {
      url = '/api/auth/admin/login'
      routeName = 'admin-dashboard'
    } else {
      url = '/api/auth/parent/login'
      routeName = 'parent-home'
    }
    var res = await apiPost(url, { username: u, password: p })
    authStore.setAuth(res.data.token, res.data.user)
    toast.show('登录成功', 'success')
    router.replace({ name: routeName })
  } catch (e) {
    if (loginType.value === 'teacher') {
      attempts++
      if (attempts >= MAX) {
        toast.show('密码错误次数过多，请联系管理员', 'error')
      } else {
        toast.show('账号或密码错误，还剩 ' + (MAX - attempts) + ' 次', 'error')
      }
    } else {
      toast.show('账号或密码错误', 'error')
    }
  } finally {
    loading.value = false
  }
}

function thirdPartyLogin(key: string) {
  var p = platforms.find(function (x) { return x.key === key })
  toast.show('正在打开' + (p ? p.label : key) + '扫码...', 'success')
}

function switchType(t: string) {
  loginType.value = t
  resetForm()
}
</script>

<template>
  <div class="card">
    <div class="header">
      <span class="icon">🌌</span>
      <h1 class="brand">学趣星球</h1>
    </div>
    <div class="tabs">
      <button
        v-for="t in (['teacher', 'admin', 'parent'] as const)"
        :key="t"
        :class="['tab', { active: loginType === t }]"
        @click="switchType(t)"
      >
        {{ t === 'teacher' ? '教师' : t === 'admin' ? '管理员' : '家长' }}
      </button>
    </div>
    <div class="form">
      <div class="field">
        <label>账号</label>
        <input v-model="username" :placeholder="loginType === 'teacher' ? '教师账号' : loginType === 'admin' ? '管理员账号' : '家长账号'" @keydown.enter="focusPwd">
      </div>
      <div class="field">
        <label>密码</label>
        <input ref="pwdRef" v-model="password" type="password" placeholder="输入密码" @keydown.enter="doLogin">
      </div>
      <button
        class="submit"
        :class="{ amber: loginType === 'admin', green: loginType === 'parent' }"
        :disabled="loading"
        @click="doLogin"
      >{{ loading ? '...' : '登录' }}</button>
      <div v-if="loginType === 'teacher'" class="social">
        <div class="social-label"><span></span> 扫码登录 <span></span></div>
        <div class="social-btns">
          <button v-for="p in platforms" :key="p.key" @click="thirdPartyLogin(p.key)">
            <span :style="{ background: p.color }">{{ p.icon }}</span>
            {{ p.label }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card { width: 100%; max-width: 360px }
.header { text-align: center; margin-bottom: 28px }
.icon { font-size: 40px; margin-bottom: 8px; display: block }
.brand { font-size: 24px; font-weight: 800; color: #1d1d1f }
.tabs {
  display: flex; gap: 0; margin-bottom: 24px;
  background: #f5f5f7; border-radius: 12px;
  border: 1px solid #e5e5ea; padding: 4px;
}
.tab {
  flex: 1; padding: 10px 8px; text-align: center;
  font-size: 14px; font-weight: 500; border-radius: 9px;
  cursor: pointer; border: none; background: transparent;
  color: #86868b; transition: all .25s;
}
.tab.active {
  background: linear-gradient(135deg,#5e5ce6,#818cf8);
  color: #fff;
  box-shadow: 0 2px 8px rgba(94,92,230,.2);
}
.form { display: flex; flex-direction: column; gap: 16px }
.field { display: flex; flex-direction: column; gap: 4px }
.field label { font-size: 13px; font-weight: 500; color: #6e6e73 }
.field input {
  width: 100%; padding: 11px 14px;
  background: #f5f5f7; border: 1px solid #e5e5ea;
  border-radius: 10px; color: #1d1d1f; font-size: 14px;
  outline: none; transition: all .2s;
}
.field input:focus {
  border-color: #5e5ce6; background: #fff;
  box-shadow: 0 0 0 3px rgba(94,92,230,.12);
}
.field input::placeholder { color: #aeaeb2 }
.submit {
  width: 100%; padding: 12px; border-radius: 10px;
  background: linear-gradient(135deg,#5e5ce6,#818cf8);
  color: #fff; font-size: 15px; font-weight: 600;
  border: none; cursor: pointer;
  box-shadow: 0 4px 14px rgba(94,92,230,.2);
  transition: all .3s;
}
.submit:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 20px rgba(94,92,230,.3) }
.submit:disabled { opacity: .6; cursor: not-allowed }
.submit.amber { background: linear-gradient(135deg,#f59e0b,#fbbf24); box-shadow: 0 4px 14px rgba(245,158,11,.2) }
.submit.green { background: linear-gradient(135deg,#10b981,#34d399); box-shadow: 0 4px 14px rgba(16,185,129,.2) }
.social { margin-top: 24px; text-align: center }
.social-label {
  display: flex; align-items: center; gap: 8px;
  color: #aeaeb2; font-size: 12px; margin-bottom: 14px;
}
.social-label span { flex: 1; height: 1px; background: #e5e5ea }
.social-btns { display: grid; grid-template-columns: repeat(4,1fr); gap: 8px }
.social-btns button {
  display: flex; flex-direction: column; align-items: center; gap: 5px;
  padding: 10px 4px; background: #f5f5f7; border: 1px solid #e5e5ea;
  border-radius: 10px; color: #6e6e73; font-size: 11px;
  cursor: pointer; transition: all .2s;
}
.social-btns button:hover { background: #fff; border-color: #d2d2d7; color: #1d1d1f }
.social-btns button span {
  display: flex; align-items: center; justify-content: center;
  width: 28px; height: 28px; border-radius: 8px;
  font-size: 14px; color: #fff;
}
</style>
