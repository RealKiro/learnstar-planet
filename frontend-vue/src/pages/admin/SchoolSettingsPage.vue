<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import type { ApiResponse } from '@/types'

interface School {
  name: string; code: string; address: string
  contact_phone: string; contact_email: string
  logo_path?: string; settings?: unknown; status: string
}

const loading = ref(true)
const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const restoreStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
// 第三方登录平台选项（管理员勾选哪些可用）
const thirdPartyPlatformOptions = [
  { key: 'wechat_work', label: '企业微信', icon: '💼', color: '#2B7CE9' },
  { key: 'dingtalk', label: '钉钉', icon: '🔷', color: '#0089FF' },
  { key: 'feishu', label: '飞书', icon: '🪶', color: '#3370FF' },
  { key: 'renren', label: '人人通空间', icon: '🌐', color: '#FF6A00' },
  { key: 'wechat', label: '微信', icon: '💬', color: '#07C160' },
  { key: 'qq', label: 'QQ', icon: '🐧', color: '#12B7F5' },
]
const form = ref({ name: '', address: '', contact_phone: '', contact_email: '', third_party_platform: '', third_party_platforms: [] as string[] })
const schoolCode = ref('')
const schoolStatus = ref('')
const schoolErrors = reactive<Record<string, string>>({})
function clsErr(f: string) { delete schoolErrors[f] }
function vldSch(field: string): boolean {
  if (field === 'name' && !form.value.name.trim()) { schoolErrors.name = '学校名称不能为空'; return false }
  if (field === 'contact_phone' && form.value.contact_phone && !/^[\d\-()+\s]{7,20}$/.test(form.value.contact_phone)) { schoolErrors.contact_phone = '联系电话格式不正确'; return false }
  if (field === 'contact_email' && form.value.contact_email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.contact_email)) { schoolErrors.contact_email = '邮箱格式不正确'; return false }
  delete schoolErrors[field]; return true
}
function saveSchool() {
  Object.keys(schoolErrors).forEach(k => delete schoolErrors[k])
  if (!form.value.name.trim()) schoolErrors.name = '学校名称不能为空'
  if (Object.keys(schoolErrors).length > 0) return
  save()
}
const logoPath = ref('')
const uploadLogoStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const activeTab = ref<'school' | 'diagnose' | 'status' | 'logs'>('school')

// 诊断
interface DiagItem { item: string; status: string; detail?: string }
const diagnoseStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const diagResult = ref<DiagItem[] | null>(null)
const diagHasIssues = ref(false)
const repairStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const repairDone = ref(false)
async function diagnose() {
  diagnoseStatus.value = 'loading'; diagResult.value = null; repairDone.value = false
  try {
    const res = await fetch('/api/v1/admin/system/diagnose', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    diagResult.value = data.data || []
    diagHasIssues.value = data.has_issues || false
    diagnoseStatus.value = 'success'
    setTimeout(() => { diagnoseStatus.value = 'idle' }, 1500)
  } catch { diagnoseStatus.value = 'error'; setTimeout(() => { diagnoseStatus.value = 'idle' }, 3000) }
}
async function repair() {
  repairStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/admin/system/repair', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' }
    })
    await res.json()
    if (!res.ok) { repairStatus.value = 'error'; setTimeout(() => { repairStatus.value = 'idle' }, 3000); return }
    repairStatus.value = 'success'
    repairDone.value = true
    setTimeout(diagnose, 500)
    setTimeout(() => { repairStatus.value = 'idle' }, 1500)
  } catch { repairStatus.value = 'error'; setTimeout(() => { repairStatus.value = 'idle' }, 3000) }
}

// 系统状态
interface SysStatus { version: Record<string, string>; migrations: { migration: string; batch: number }[]; migration_count: number }
const sysStatus = ref<SysStatus | null>(null)
const statusLoading = ref(false)
const logEntries = ref<{ line: string; level: string }[]>([])
const logTotal = ref(0)
const logLines = ref(200)
const logLevel = ref('')
const logRefreshing = ref(false)
const logError = ref('')
let logTimer: ReturnType<typeof setInterval> | null = null
async function loadLogs() {
  logRefreshing.value = true; logError.value = ''
  try {
    const res = await fetch(`/api/v1/admin/system/logs?lines=${logLines.value}${logLevel.value ? '&level=' + logLevel.value : ''}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    if (!res.ok) { logError.value = data.message || '加载失败'; return }
    logEntries.value = data.data?.entries || []
    logTotal.value = data.data?.total || 0
  } catch { logError.value = '加载日志失败' }
  finally { logRefreshing.value = false }
}
function startLogPolling() { stopLogPolling(); logTimer = setInterval(loadLogs, 5000) }
function stopLogPolling() { if (logTimer) { clearInterval(logTimer); logTimer = null } }
/** 日志级别 → 着色分组（ERROR/CRITICAL/ALERT/EMERGENCY 归为 error，WARN* 归为 warn） */
function logLevelClass(level: string): string {
  const lv = (level || '').toUpperCase()
  if (lv.includes('ERROR') || lv.includes('CRITICAL') || lv.includes('ALERT') || lv.includes('EMERGENCY')) return 'error'
  if (lv.includes('WARN')) return 'warn'
  if (lv.includes('DEBUG')) return 'debug'
  return 'info'
}
onUnmounted(stopLogPolling)

onMounted(async () => {
  try {
    const [schoolRes, statusRes] = await Promise.all([
      apiGet<ApiResponse<School>>('/api/v1/admin/school'),
      fetch('/api/v1/admin/system/status', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json()).catch(() => ({ data: null })),
    ])
    const s = (schoolRes as unknown as { data: School }).data
    const st = (s.settings as any) || {}
    form.value = {
      name: s.name || '', address: s.address || '',
      contact_phone: s.contact_phone || '', contact_email: s.contact_email || '',
      third_party_platform: st.third_party_platform || '',
      third_party_platforms: Array.isArray(st.enabled_third_party_platforms) && st.enabled_third_party_platforms.length > 0
        ? st.enabled_third_party_platforms
        : ['wechat_work', 'wechat', 'qq'],
    }
    schoolCode.value = s.code || ''
    schoolStatus.value = s.status || ''
    logoPath.value = s.logo_path || ''
    sysStatus.value = statusRes.data || null
  } catch { /* handled */ }
  finally { loading.value = false }
})

async function save() {
  saveStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/admin/school', {
      method: 'PUT',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: form.value.name.trim(), address: form.value.address.trim(), contact_phone: form.value.contact_phone.trim(), contact_email: form.value.contact_email.trim(), settings: { third_party_platform: form.value.third_party_platform, enabled_third_party_platforms: form.value.third_party_platforms } }),
    })
    await res.json()
    if (!res.ok) { saveStatus.value = 'error'; setTimeout(() => { saveStatus.value = 'idle' }, 3000); return }
    saveStatus.value = 'success'
    setTimeout(() => { saveStatus.value = 'idle' }, 1500)
  } catch { saveStatus.value = 'error'; setTimeout(() => { saveStatus.value = 'idle' }, 3000) }
}

async function reload() {
  restoreStatus.value = 'loading'
  try {
    const res = await apiGet<ApiResponse<School>>('/api/v1/admin/school')
    const s = (res as unknown as { data: School }).data
    const st = (s.settings as any) || {}
    form.value = {
      name: s.name || '', address: s.address || '',
      contact_phone: s.contact_phone || '', contact_email: s.contact_email || '',
      third_party_platform: st.third_party_platform || '',
      third_party_platforms: Array.isArray(st.enabled_third_party_platforms) && st.enabled_third_party_platforms.length > 0
        ? st.enabled_third_party_platforms
        : ['wechat_work', 'wechat', 'qq'],
    }
    logoPath.value = s.logo_path || ''
    restoreStatus.value = 'success'
    setTimeout(() => { restoreStatus.value = 'idle' }, 1500)
  } catch {
    restoreStatus.value = 'error'
    setTimeout(() => { restoreStatus.value = 'idle' }, 3000)
  }
}

async function uploadLogo(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploadLogoStatus.value = 'loading'
  try {
    const fd = new FormData()
    fd.append('logo', file)
    const res = await fetch('/api/v1/admin/school/logo', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') },
      body: fd,
    })
    const data = await res.json()
    if (res.ok && data.data?.logo_path) logoPath.value = data.data.logo_path
    if (!res.ok) { uploadLogoStatus.value = 'error'; setTimeout(() => { uploadLogoStatus.value = 'idle' }, 3000); return }
    uploadLogoStatus.value = 'success'
    setTimeout(() => { uploadLogoStatus.value = 'idle' }, 1500)
  } catch { uploadLogoStatus.value = 'error'; setTimeout(() => { uploadLogoStatus.value = 'idle' }, 3000) }
  finally { (e.target as HTMLInputElement).value = '' }
}
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <p class="page-eyebrow">系统配置</p>
        <h2 class="page-title">系统信息</h2>
      </div>
    </div>

    <div class="tab-bar tabs-narrow">
      <button :class="['tab-btn', { active: activeTab === 'school' }]" @click="activeTab = 'school'">🏫 学校信息</button>
      <button :class="['tab-btn', { active: activeTab === 'diagnose' }]" @click="activeTab = 'diagnose'">🔍 系统诊断</button>
      <button :class="['tab-btn', { active: activeTab === 'status' }]" @click="activeTab = 'status'">📊 系统状态</button>
      <button :class="['tab-btn', { active: activeTab === 'logs' }]" @click="activeTab = 'logs'; loadLogs()">📋 实时日志</button>
    </div>

    <!-- 学校信息 -->
    <div v-if="activeTab === 'school'">
      <div v-if="loading" class="loading-spinner">加载中...</div>
      <div v-else class="card card-form-lg">
        <div class="info-row">
          <div class="info-block"><div class="info-label">学校编码</div><div class="mono">{{ schoolCode || '-' }}</div></div>
          <div class="info-block"><div class="info-label">状态</div><span class="badge-active">{{ schoolStatus === 'active' ? '正常运行' : (schoolStatus || '正常') }}</span></div>
        </div>
        <!-- LOGO -->
        <div class="form-group">
          <label>学校 LOGO</label>
          <div class="logo-row">
            <div v-if="logoPath" class="logo-box">
              <img :src="logoPath" class="logo-img">
            </div>
            <label class="btn btn-sm btn-outline" :class="uploadLogoStatus !== 'idle' ? 'btn-state-' + uploadLogoStatus : ''">
              <template v-if="uploadLogoStatus === 'loading'">上传中...</template>
              <template v-else-if="uploadLogoStatus === 'success'">已上传 ✓</template>
              <template v-else-if="uploadLogoStatus === 'error'">上传失败 ✗</template>
              <template v-else>📷 上传 LOGO</template>
              <input type="file" accept="image/*" class="file-hidden" @change="uploadLogo">
            </label>
          </div>
        </div>
        <div class="form-group"><label>学校名称 <span class="req-star">*</span></label><input v-model="form.name" class="form-input" :class="{ 'input-invalid': schoolErrors.name }" placeholder="请输入学校名称" @blur="vldSch('name')" @input="clsErr('name')"><div v-if="schoolErrors.name" class="field-error">{{ schoolErrors.name }}</div></div>
        <div class="form-group"><label>学校地址</label><input v-model="form.address" class="form-input" placeholder="请输入学校地址"></div>
        <div class="form-group"><label>联系电话</label><input v-model="form.contact_phone" class="form-input" :class="{ 'input-invalid': schoolErrors.contact_phone }" placeholder="如：021-12345678" @blur="vldSch('contact_phone')" @input="clsErr('contact_phone')"><div v-if="schoolErrors.contact_phone" class="field-error">{{ schoolErrors.contact_phone }}</div></div>
        <div class="form-group"><label>联系邮箱</label><input v-model="form.contact_email" type="email" class="form-input" :class="{ 'input-invalid': schoolErrors.contact_email }" placeholder="如：admin@school.edu.cn" @blur="vldSch('contact_email')" @input="clsErr('contact_email')"><div v-if="schoolErrors.contact_email" class="field-error">{{ schoolErrors.contact_email }}</div></div>
        <div class="form-group">
          <label>第三方登录平台（勾选的平台才会在登录页显示）</label>
          <div class="platform-grid">
            <label
              v-for="opt in thirdPartyPlatformOptions"
              :key="opt.key"
              class="platform-opt"
              :class="{ 'is-on': form.third_party_platforms.includes(opt.key) }"
            >
              <input
                type="checkbox"
                :value="opt.key"
                v-model="form.third_party_platforms"
                class="platform-check"
              >
              <span class="platform-icon"><PlatformIcon :platform="opt.key" :size="20" /></span>
              <span>{{ opt.label }}</span>
            </label>
          </div>
          <p class="hint-11">未勾选任何平台时，登录页默认显示企业微信/微信/QQ。通讯录导入需配置对应平台的应用凭证。</p>
        </div>
        <div class="form-actions">
          <button class="btn btn-sm" :class="restoreStatus === 'idle' ? 'btn-state-plain' : 'btn-state-' + restoreStatus" :disabled="restoreStatus === 'loading'" @click="reload">
            <template v-if="restoreStatus === 'loading'">恢复中...</template>
            <template v-else-if="restoreStatus === 'success'">已恢复 ✓</template>
            <template v-else-if="restoreStatus === 'error'">恢复失败 ✗</template>
            <template v-else>↩️ 重置</template>
          </button>
          <button class="btn btn-sm" :class="'btn-state-' + saveStatus" :disabled="saveStatus !== 'idle' || loading" @click="saveSchool">
            <template v-if="saveStatus === 'loading'">保存中...</template>
            <template v-else-if="saveStatus === 'success'">已保存 ✓</template>
            <template v-else-if="saveStatus === 'error'">保存失败 ✗</template>
            <template v-else>保存设置</template>
          </button>
        </div>
      </div>
    </div>

    <!-- 系统诊断 -->
    <div v-if="activeTab === 'diagnose'" class="card card-form">
      <p class="card-subtitle">检查数据库表结构完整性</p>
      <div class="btn-row">
        <button class="btn btn-outline" :class="diagnoseStatus !== 'idle' ? 'btn-state-' + diagnoseStatus : ''" :disabled="diagnoseStatus !== 'idle'" @click="diagnose">
          <template v-if="diagnoseStatus === 'loading'">诊断中...</template>
          <template v-else-if="diagnoseStatus === 'success'">诊断完成 ✓</template>
          <template v-else-if="diagnoseStatus === 'error'">诊断失败 ✗</template>
          <template v-else>🔍 开始诊断</template>
        </button>
        <button class="btn btn-danger" :class="repairStatus !== 'idle' ? 'btn-state-' + repairStatus : ''" :disabled="repairStatus !== 'idle' || !diagHasIssues" @click="repair">
          <template v-if="repairStatus === 'loading'">修复中...</template>
          <template v-else-if="repairStatus === 'success'">修复完成 ✓</template>
          <template v-else-if="repairStatus === 'error'">修复失败 ✗</template>
          <template v-else>🛠️ 一键修复</template>
        </button>
      </div>
      <div v-if="diagResult">
        <div v-for="(r, i) in diagResult" :key="i" class="diag-row">
          <span v-if="r.status === 'ok'" class="c-ok">✅</span>
          <span v-else-if="r.status === 'fixable'" class="c-warn">⚠️</span>
          <span v-else class="c-err">❌</span>
          <span class="flex-1">{{ r.item }}</span>
          <span :class="['ts-strong', r.status === 'ok' ? 'text-ok' : r.status === 'fixable' ? 'text-warn' : 'text-err']">{{ r.status === 'ok' ? '正常' : (r.detail || '缺失') }}</span>
        </div>
        <div v-if="repairDone" class="repair-done">✅ 修复已完成</div>
      </div>
      <div v-else-if="diagnoseStatus === 'idle'" class="diag-empty">点击「开始诊断」检查系统状态</div>
    </div>

    <!-- 系统状态 -->
    <div v-if="activeTab === 'status'" class="card card-form">
      <div v-if="statusLoading" class="status-loading">加载中...</div>
      <div v-else-if="sysStatus">
        <div class="version-block">
          <div class="sub-title-sm">版本信息</div>
          <div class="version-grid">
            <div v-for="(val, key) in sysStatus.version" :key="key" class="version-item"><span class="version-key">{{ key }}</span><span class="version-val">{{ val }}</span></div>
          </div>
        </div>
        <div>
          <div class="sub-title-sm">迁移记录（{{ sysStatus.migration_count }} 条）</div>
          <div v-if="sysStatus.migrations.length" class="mig-list">
            <div v-for="(m, i) in sysStatus.migrations" :key="i" class="mig-row"><span class="mig-batch">#{{ m.batch }}</span><span class="mig-name">{{ m.migration }}</span></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 实时日志 -->
    <div v-if="activeTab === 'logs'" class="card card-wide">
      <div class="log-toolbar">
        <h3 class="log-title">📋 实时日志</h3>
        <select v-model.number="logLines" class="log-select">
          <option :value="50">50 行</option>
          <option :value="200">200 行</option>
          <option :value="500">500 行</option>
          <option :value="1000">1000 行</option>
        </select>
        <select v-model="logLevel" @change="loadLogs" class="log-select">
          <option value="">全部等级</option>
          <option value="ERROR">ERROR</option>
          <option value="WARNING">WARNING</option>
          <option value="INFO">INFO</option>
          <option value="DEBUG">DEBUG</option>
          <option value="CRITICAL">CRITICAL</option>
        </select>
        <button class="btn btn-sm btn-card" @click="loadLogs" :disabled="logRefreshing">{{ logRefreshing ? '加载中...' : '🔄 刷新' }}</button>
        <button class="btn btn-sm btn-purple" @click="startLogPolling">▶ 自动刷新</button>
        <button class="btn btn-sm btn-card" @click="stopLogPolling">⏹ 停止</button>
        <span class="log-count">{{ logTotal }} 条</span>
      </div>
      <div v-if="logError" class="log-error">{{ logError }}</div>
      <div class="log-view">
        <div v-if="!logEntries.length" class="log-empty">暂无日志</div>
        <div v-for="(entry, i) in logEntries" :key="i" class="log-line">
          <span class="log-level" :class="'log-level--' + logLevelClass(entry.level)">{{ entry.level }}</span>
          <span class="log-text">{{ entry.line }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }
.info-block { flex:1; min-width:200px; padding:12px 16px; background:var(--color-bg); border-radius:10px; }
.info-label { font-size:12px; color:var(--color-text-secondary); margin-bottom:4px; }
.mono { font-family:monospace; font-weight:600; }
.badge-active { display:inline-block; padding:3px 12px; border-radius:20px; font-size:12px; font-weight:600; background:rgba(16,185,129,0.1); color: var(--c-green); }
.form-group { margin-bottom:14px; }
.form-group label { display:block; font-size:12px; font-weight:600; color:var(--color-text); margin-bottom:4px; }
.form-input { color:var(--color-text); width:100%; padding:8px 12px; border:1px solid var(--color-border); border-radius:8px; font-size:13px; outline:none; transition:border-color 0.15s; box-sizing:border-box; background:var(--color-bg-card); }
.form-input:focus { border-color:#7c3aed; box-shadow:0 0 0 3px rgba(124,58,237,0.08); }
.version-grid { display:grid; grid-template-columns:1fr 1fr; gap:6px; }
.version-item { display:flex; justify-content:space-between; padding:6px 12px; background:var(--color-bg); border-radius:6px; font-size:12px; }
.version-key { color:var(--color-text-secondary); font-weight:500; }
.version-val { font-family:monospace; font-weight:600; }
.mig-list { max-height:360px; overflow-y:auto; border:1px solid var(--color-border); border-radius:8px; }
.mig-row { display:flex; align-items:center; gap:8px; padding:6px 12px; border-bottom:1px solid var(--color-border); font-size:11px; }
.mig-row:last-child { border-bottom:none; }
.mig-batch { display:inline-block; padding:1px 6px; border-radius:4px; background:rgba(167,139,250,0.15); color:var(--md-primary-light); font-weight:700; font-size:10px; flex-shrink:0; }
.mig-name { color:var(--color-text-secondary); word-break:break-all; }
.loading-spinner { text-align:center; padding:48px; color:var(--color-text-secondary); font-size:15px; }
.btn { padding:8px 16px; border-radius:8px; font-size:13px; font-weight:500; cursor:pointer; border:1px solid transparent; transition:all 0.15s; }
.btn-sm { padding:6px 14px; font-size:12px; border-radius:8px; }
.btn-primary { background:#7c3aed; color:white; border-color:#7c3aed; }
.btn-primary:hover { background:#6d28d9; }
.btn-outline { background:var(--color-bg-card); color:var(--color-text); border:1px solid var(--color-border); }
.btn-danger { background:var(--color-bg-card); color: var(--c-red); border:1px solid rgba(239,68,68,0.2); }

/* ===== 页内布局类（本页专用，替代原内联样式；声明逐字保留以保证渲染等价） ===== */
.tabs-narrow { max-width:640px; }
.card-form-lg { max-width:640px; padding:32px; }
.card-form { max-width:640px; padding:24px; }
.card-wide { max-width:100%; padding:20px; }
.info-row { display:flex; gap:16px; margin-bottom:24px; flex-wrap:wrap; }
.logo-row { display:flex; align-items:center; gap:12px; }
.logo-box { width:64px; height:64px; border-radius:8px; overflow:hidden; border:1px solid var(--color-border); }
.logo-img { width:100%; height:100%; object-fit:cover; }
.file-hidden { display:none; }
.form-input.input-invalid { border-color:#f87171; }
.platform-grid { display:grid; grid-template-columns:repeat(2,1fr); gap:8px; margin-top:6px; }
.platform-opt { display:flex; align-items:center; gap:8px; padding:8px 12px; border:1px solid var(--color-border); border-radius:8px; cursor:pointer; font-size:13px; user-select:none; }
.platform-opt.is-on { background:rgba(124,58,237,0.08); border-color:rgba(124,58,237,0.4); }
.platform-check { accent-color:#7c3aed; width:15px; height:15px; flex-shrink:0; }
.platform-icon { flex-shrink:0; display:flex; }
.hint-11 { font-size:11px; color:var(--color-text-secondary); margin-top:2px; }
.form-actions { display:flex; gap:8px; justify-content:flex-end; margin-top:24px; }
.card-subtitle { font-size:13px; color:var(--color-text-secondary); margin-bottom:16px; }
.btn-row { display:flex; gap:12px; margin-bottom:16px; }
.diag-row { display:flex; align-items:center; gap:10px; padding:6px 10px; border-bottom:1px solid var(--color-border); font-size:13px; }
.c-ok { display: inline-flex; align-items: center; justify-content: center; color: var(--c-green); }
.c-warn { display: inline-flex; align-items: center; justify-content: center; color: var(--c-amber); }
.c-err { display: inline-flex; align-items: center; justify-content: center; color: var(--c-red); }
.repair-done { padding:8px 12px; font-size:13px; color: var(--c-green); font-weight:500; }
.diag-empty { padding:12px; text-align:center; font-size:13px; color:var(--color-text-secondary); }
.status-loading { text-align:center; padding:24px; }
.version-block { margin-bottom:20px; }
.sub-title-sm { font-size:13px; font-weight:600; color:var(--color-text-secondary); margin-bottom:8px; }
.log-toolbar { display:flex; align-items:center; gap:12px; margin-bottom:12px; flex-wrap:wrap; }
.log-title { display: flex; align-items: center; gap: 6px; font-size:15px; font-weight:700; margin:0; }
.log-select { padding:4px 8px; border-radius:6px; border:1px solid var(--color-border); background:var(--color-bg-card); color:var(--color-text); font-size:12px; }
.btn-card { background:var(--color-bg-card); color:var(--color-text); border:1px solid var(--color-border); }
.btn-purple { background:#7c3aed; color:#fff; border:none; }
.log-count { font-size:12px; color:var(--color-text-secondary); }
.log-error { color: var(--c-red); font-size:13px; padding:8px; background:rgba(239,68,68,0.06); border-radius:6px; margin-bottom:8px; }
.log-view { background:#0d1117; border-radius:10px; max-height:60vh; overflow:auto; padding:8px 0; }
.log-empty { padding:16px; text-align:center; color:#8b949e; font-size:13px; }
.log-line { display:flex; align-items:flex-start; gap:8px; padding:2px 16px; font-family:monospace; font-size:12px; line-height:1.6; white-space:pre-wrap; word-break:break-all; }
.log-level { flex-shrink:0; width:60px; font-weight:600; }
.log-level--error { color: var(--c-red-soft); }
.log-level--warn { color: var(--c-amber); }
.log-level--info { color:#60A5FA; }
.log-level--debug { color:#94A3B8; }
.log-text { color:#e6edf3; }
.ts-strong { font-weight: 600; }
</style>
