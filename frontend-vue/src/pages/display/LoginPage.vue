<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiPost } from '@/utils/api'

const router = useRouter()

const code = ref('')
const loading = ref(false)
const error = ref('')
const className = ref('')
const studentCount = ref(0)
const showSuccess = ref(false)
const codeParts = ref(['', '', '']) // 3-1-A7K2 → ['3', '1', 'A7K2']

// 如果 URL 中有 code 参数，自动尝试登录
onMounted(() => {
  const params = new URLSearchParams(window.location.search)
  const urlCode = params.get('code')
  if (urlCode) {
    code.value = urlCode
    handleLogin()
  }
})

function onCodeInput(value: string) {
  // 自动格式化：只保留字母数字和短横线
  code.value = value.toUpperCase().replace(/[^0-9A-Z-]/g, '')
  error.value = ''
}

function handlePaste(event: ClipboardEvent) {
  const text = event.clipboardData?.getData('text') || ''
  code.value = text.toUpperCase().replace(/[^0-9A-Z-]/g, '')
  event.preventDefault()
  if (code.value.length >= 6) {
    handleLogin()
  }
}

async function handleLogin() {
  if (code.value.length < 4) {
    error.value = '请输入完整的班级码'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const res = await apiPost<{
      data: {
        token: string
        expires_in: number
        class_info: {
          id: number
          name: string
          grade: string
          student_count: number
        }
      }
    }>('/api/v1/display/login', { code: code.value })

    const { token, class_info } = res.data

    // 保存到 sessionStorage（关闭浏览器即失效）
    sessionStorage.setItem('display_token', token)
    sessionStorage.setItem('display_class_info', JSON.stringify(class_info))

    // 显示成功 + 班级信息
    className.value = class_info.name
    studentCount.value = class_info.student_count
    showSuccess.value = true

    // 1.5 秒后自动跳转到大屏
    setTimeout(() => {
      router.push({ name: 'display-screen' })
    }, 1500)
  } catch (e: any) {
    const msg = e?.response?.data?.message || e?.message || '登录失败'
    error.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="display-login">
    <div class="login-container">
      <!-- Logo -->
      <div class="login-header">
        <div class="login-logo">🌌</div>
        <h1>学趣星球</h1>
        <p class="login-subtitle">班级大屏</p>
      </div>

      <!-- 登录成功 -->
      <div v-if="showSuccess" class="login-success">
        <div class="success-icon">✅</div>
        <h2>登录成功！</h2>
        <p class="success-class">{{ className }}</p>
        <p class="success-count">{{ studentCount }} 位同学</p>
        <p class="success-hint">即将进入大屏...</p>
      </div>

      <!-- 登录表单 -->
      <form v-else class="login-form" @submit.prevent="handleLogin">
        <div class="input-label">班级码</div>
        <div class="code-input-wrapper">
          <input
            v-model="code"
            type="text"
            class="code-input"
            placeholder="如 LS301"
            maxlength="12"
            autocomplete="off"
            autocorrect="off"
            spellcheck="false"
            @input="onCodeInput(($event.target as HTMLInputElement).value)"
            @paste="handlePaste"
          />
        </div>

        <p class="input-hint">
          请向班主任老师获取班级码
        </p>

        <div v-if="error" class="error-msg">{{ error }}</div>

        <button
          type="submit"
          class="login-btn"
          :disabled="loading || code.length < 4"
        >
          <span v-if="loading" class="btn-loading">⏳</span>
          <span v-else>进入大屏</span>
        </button>

        <div class="qr-hint">
          <span class="qr-icon">📱</span>
          也可扫描教师分享的二维码
        </div>
      </form>

      <!-- Footer -->
      <div class="login-footer">
        <span class="footer-dot"></span>
        连接中，请输入班级码
        <span class="footer-dot"></span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.display-login {
  min-height: 100vh;
  background: linear-gradient(135deg, #0c0a20 0%, #1a1040 30%, #0d1b2a 70%, #0a1628 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "PingFang SC", "Noto Sans SC", system-ui, sans-serif;
  user-select: none;
  overflow: hidden;
}

.login-container {
  width: 100%;
  max-width: 420px;
  padding: 40px 32px;
  text-align: center;
}

.login-header {
  margin-bottom: 40px;
}

.login-logo {
  font-size: 64px;
  margin-bottom: 8px;
  filter: drop-shadow(0 0 20px rgba(120, 80, 255, 0.4));
  animation: float 3s ease-in-out infinite;
}

.login-logo + h1 {
  font-size: 28px;
  font-weight: 700;
  color: #f0ecff;
  margin: 0 0 4px;
  letter-spacing: 0.05em;
}

.login-subtitle {
  font-size: 14px;
  color: rgba(200, 190, 240, 0.5);
  margin: 0;
}

/* Form */
.login-form {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  padding: 32px 24px;
  backdrop-filter: blur(12px);
}

.input-label {
  font-size: 13px;
  font-weight: 600;
  color: rgba(200, 190, 240, 0.7);
  margin-bottom: 12px;
  text-align: left;
}

.code-input-wrapper {
  position: relative;
}

.code-input {
  width: 100%;
  padding: 16px 20px;
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 0.15em;
  text-align: center;
  border: 2px solid rgba(124, 58, 237, 0.3);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.04);
  color: #f0ecff;
  outline: none;
  transition: all 0.2s;
  box-sizing: border-box;
}

.code-input:focus {
  border-color: #7c3aed;
  box-shadow: 0 0 0 4px rgba(124, 58, 237, 0.15), 0 0 30px rgba(124, 58, 237, 0.1);
}

.code-input::placeholder {
  color: rgba(200, 190, 240, 0.2);
  font-weight: 400;
  font-size: 20px;
  letter-spacing: 0.05em;
}

.input-hint {
  font-size: 12px;
  color: rgba(200, 190, 240, 0.35);
  margin: 12px 0 0;
  text-align: left;
}

.error-msg {
  margin-top: 12px;
  padding: 10px 14px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 10px;
  color: #fca5a5;
  font-size: 13px;
}

.login-btn {
  width: 100%;
  margin-top: 20px;
  padding: 14px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  color: white;
  background: linear-gradient(135deg, #7c3aed, #6d28d9);
  cursor: pointer;
  transition: all 0.2s;
}

.login-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(124, 58, 237, 0.25);
}

.login-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-loading {
  display: inline-block;
  animation: spin 1s linear infinite;
}

.qr-hint {
  margin-top: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(200, 190, 240, 0.35);
}

.qr-icon {
  font-size: 16px;
}

/* Success */
.login-success {
  background: rgba(16, 185, 129, 0.05);
  border: 1px solid rgba(16, 185, 129, 0.15);
  border-radius: 20px;
  padding: 40px 24px;
}

.success-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.login-success h2 {
  font-size: 22px;
  font-weight: 700;
  color: #6ee7b7;
  margin: 0 0 8px;
}

.success-class {
  font-size: 18px;
  font-weight: 600;
  color: #f0ecff;
  margin: 0 0 4px;
}

.success-count {
  font-size: 14px;
  color: rgba(200, 190, 240, 0.5);
  margin: 0 0 20px;
}

.success-hint {
  font-size: 13px;
  color: rgba(200, 190, 240, 0.4);
  animation: pulse 1.5s ease-in-out infinite;
}

/* Footer */
.login-footer {
  margin-top: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 11px;
  color: rgba(200, 190, 240, 0.2);
}

.footer-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: rgba(200, 190, 240, 0.2);
  animation: breathe 2s ease-in-out infinite;
}

/* Animations */
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}

@keyframes breathe {
  0%, 100% { opacity: 0.2; }
  50% { opacity: 0.6; }
}
</style>
