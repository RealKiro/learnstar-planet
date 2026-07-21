<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { apiGet, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface School {
  name: string; code: string; address: string
  contact_phone: string; contact_email: string
  logo_path?: string; settings?: unknown; status: string
}

const loading = ref(true)
const saveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const restoreStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const form = ref({ name: '', address: '', contact_phone: '', contact_email: '' })
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
    const data = await res.json()
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
async function loadStatus() {
  statusLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/system/status', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    sysStatus.value = data.data || null
  } catch { /* ignore */ }
  finally { statusLoading.value = false }
}
const logContent = ref('')
const logLines = ref(200)
const logRefreshing = ref(false)
const logError = ref('')
let logTimer: ReturnType<typeof setInterval> | null = null
async function loadLogs() {
  logRefreshing.value = true; logError.value = ''
  try {
    const res = await fetch(`/api/v1/admin/system/logs?lines=${logLines.value}`, {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    if (!res.ok) { logError.value = data.message || '加载失败'; return }
    logContent.value = data.data?.content || ''
  } catch { logError.value = '加载日志失败' }
  finally { logRefreshing.value = false }
}
function startLogPolling() { stopLogPolling(); logTimer = setInterval(loadLogs, 5000) }
function stopLogPolling() { if (logTimer) { clearInterval(logTimer); logTimer = null } }
onUnmounted(stopLogPolling)

onMounted(async () => {
  try {
    const [schoolRes, statusRes] = await Promise.all([
      apiGet<ApiResponse<School>>('/api/v1/admin/school'),
      fetch('/api/v1/admin/system/status', { headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') } }).then(r => r.json()).catch(() => ({ data: null })),
    ])
    const s = (schoolRes as unknown as { data: School }).data
    form.value = { name: s.name || '', address: s.address || '', contact_phone: s.contact_phone || '', contact_email: s.contact_email || '' }
    schoolCode.value = s.code || ''
    schoolStatus.value = s.status || ''
    logoPath.value = s.logo_path || ''
    sysStatus.value = statusRes.data || null
  } catch { /* handled */ }
  finally { loading.value = false }
})

async function save() {
  if (!form.value.name.trim()) { toast.show('请填写学校名称', 'error', { position: 'center', duration: 2000 }); return }
  saveStatus.value = 'loading'
  try {
    const res = await fetch('/api/v1/admin/school', {
      method: 'PUT',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: form.value.name.trim(), address: form.value.address.trim(), contact_phone: form.value.contact_phone.trim(), contact_email: form.value.contact_email.trim() }),
    })
    const data = await res.json()
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
    form.value = { name: s.name || '', address: s.address || '', contact_phone: s.contact_phone || '', contact_email: s.contact_email || '' }
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
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">系统配置</p>
        <h2 style="font-size:24px;font-weight:700;">系统信息</h2>
      </div>
    </div>

    <div class="tab-bar" style="max-width:640px;">
      <button :class="['tab-btn', { active: activeTab === 'school' }]" @click="activeTab = 'school'">🏫 学校信息</button>
      <button :class="['tab-btn', { active: activeTab === 'diagnose' }]" @click="activeTab = 'diagnose'">🔍 系统诊断</button>
      <button :class="['tab-btn', { active: activeTab === 'status' }]" @click="activeTab = 'status'">📊 系统状态</button>
      <button :class="['tab-btn', { active: activeTab === 'logs' }]" @click="activeTab = 'logs'; loadLogs()">📋 实时日志</button>
    </div>

    <!-- 学校信息 -->
    <div v-if="activeTab === 'school'">
      <div v-if="loading" class="loading-spinner">加载中...</div>
      <div v-else class="card" style="max-width:640px;padding:32px;">
        <div style="display:flex;gap:16px;margin-bottom:24px;flex-wrap:wrap;">
          <div class="info-block"><div class="info-label">学校编码</div><div class="mono">{{ schoolCode || '-' }}</div></div>
          <div class="info-block"><div class="info-label">状态</div><span class="badge-active">{{ schoolStatus === 'active' ? '正常运行' : (schoolStatus || '正常') }}</span></div>
        </div>
        <!-- LOGO -->
        <div class="form-group">
          <label>学校 LOGO</label>
          <div style="display:flex;align-items:center;gap:12px;">
            <div v-if="logoPath" style="width:64px;height:64px;border-radius:8px;overflow:hidden;border:1px solid var(--color-border);">
              <img :src="logoPath" style="width:100%;height:100%;object-fit:cover;">
            </div>
            <label class="btn btn-sm btn-outline" :style="{ background: uploadLogoStatus === 'loading' ? '#f59e0b' : uploadLogoStatus === 'success' ? '#10b981' : uploadLogoStatus === 'error' ? '#ef4444' : '', color: uploadLogoStatus !== 'idle' ? '#fff' : '', border: uploadLogoStatus !== 'idle' ? '1px solid transparent' : '', cursor: 'pointer' }">
              <template v-if="uploadLogoStatus === 'loading'">上传中...</template>
              <template v-else-if="uploadLogoStatus === 'success'">已上传 ✓</template>
              <template v-else-if="uploadLogoStatus === 'error'">上传失败 ✗</template>
              <template v-else>📷 上传 LOGO</template>
              <input type="file" accept="image/*" style="display:none;" @change="uploadLogo">
            </label>
          </div>
        </div>
        <div class="form-group"><label>学校名称 <span style="color:#f87171;">*</span></label><input v-model="form.name" class="form-input" placeholder="请输入学校名称" :style="{ borderColor: schoolErrors.name ? '#f87171' : '' }" @blur="vldSch('name')" @input="clsErr('name')"><div v-if="schoolErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ schoolErrors.name }}</div></div>
        <div class="form-group"><label>学校地址</label><input v-model="form.address" class="form-input" placeholder="请输入学校地址"></div>
        <div class="form-group"><label>联系电话</label><input v-model="form.contact_phone" class="form-input" placeholder="如：021-12345678" :style="{ borderColor: schoolErrors.contact_phone ? '#f87171' : '' }" @blur="vldSch('contact_phone')" @input="clsErr('contact_phone')"><div v-if="schoolErrors.contact_phone" style="color:#f87171;font-size:11px;margin-top:2px;">{{ schoolErrors.contact_phone }}</div></div>
        <div class="form-group"><label>联系邮箱</label><input v-model="form.contact_email" type="email" class="form-input" placeholder="如：admin@school.edu.cn" :style="{ borderColor: schoolErrors.contact_email ? '#f87171' : '' }" @blur="vldSch('contact_email')" @input="clsErr('contact_email')"><div v-if="schoolErrors.contact_email" style="color:#f87171;font-size:11px;margin-top:2px;">{{ schoolErrors.contact_email }}</div></div>
        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:24px;">
          <button class="btn btn-sm" :style="{ background: restoreStatus === 'loading' ? '#f59e0b' : restoreStatus === 'success' ? '#10b981' : restoreStatus === 'error' ? '#ef4444' : '', color: restoreStatus !== 'idle' ? '#fff' : 'var(--color-text)', border: restoreStatus !== 'idle' ? '1px solid transparent' : '1px solid var(--color-border)' }" :disabled="restoreStatus === 'loading'" @click="reload">
            <template v-if="restoreStatus === 'loading'">恢复中...</template>
            <template v-else-if="restoreStatus === 'success'">已恢复 ✓</template>
            <template v-else-if="restoreStatus === 'error'">恢复失败 ✗</template>
            <template v-else>↩️ 重置</template>
          </button>
          <button class="btn btn-sm" :style="{ background: saveStatus === 'loading' ? '#f59e0b' : saveStatus === 'success' ? '#10b981' : saveStatus === 'error' ? '#ef4444' : '#7c3aed', color: '#fff', border: '1px solid transparent' }" :disabled="saveStatus !== 'idle' || loading" @click="saveSchool">
            <template v-if="saveStatus === 'loading'">保存中...</template>
            <template v-else-if="saveStatus === 'success'">已保存 ✓</template>
            <template v-else-if="saveStatus === 'error'">保存失败 ✗</template>
            <template v-else>保存设置</template>
          </button>
        </div>
      </div>
    </div>

    <!-- 系统诊断 -->
    <div v-if="activeTab === 'diagnose'" class="card" style="max-width:640px;padding:24px;">
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">检查数据库表结构完整性</p>
      <div style="display:flex;gap:12px;margin-bottom:16px;">
        <button class="btn btn-outline" :style="{ background: diagnoseStatus === 'loading' ? '#f59e0b' : diagnoseStatus === 'success' ? '#10b981' : diagnoseStatus === 'error' ? '#ef4444' : '', color: diagnoseStatus !== 'idle' ? '#fff' : '', border: diagnoseStatus !== 'idle' ? '1px solid transparent' : '' }" :disabled="diagnoseStatus !== 'idle'" @click="diagnose">
          <template v-if="diagnoseStatus === 'loading'">诊断中...</template>
          <template v-else-if="diagnoseStatus === 'success'">诊断完成 ✓</template>
          <template v-else-if="diagnoseStatus === 'error'">诊断失败 ✗</template>
          <template v-else>🔍 开始诊断</template>
        </button>
        <button class="btn btn-danger" :style="{ background: repairStatus === 'loading' ? '#f59e0b' : repairStatus === 'success' ? '#10b981' : repairStatus === 'error' ? '#ef4444' : '', color: repairStatus !== 'idle' ? '#fff' : '', border: repairStatus !== 'idle' ? '1px solid transparent' : '' }" :disabled="repairStatus !== 'idle' || !diagHasIssues" @click="repair">
          <template v-if="repairStatus === 'loading'">修复中...</template>
          <template v-else-if="repairStatus === 'success'">修复完成 ✓</template>
          <template v-else-if="repairStatus === 'error'">修复失败 ✗</template>
          <template v-else>🛠️ 一键修复</template>
        </button>
      </div>
      <div v-if="diagResult">
        <div v-for="(r, i) in diagResult" :key="i" style="display:flex;align-items:center;gap:10px;padding:6px 10px;border-bottom:1px solid var(--color-border);font-size:13px;">
          <span v-if="r.status === 'ok'" style="color:#10B981;">✅</span><span v-else style="color:#EF4444;">❌</span>
          <span style="flex:1;">{{ r.item }}</span>
          <span :style="{ color: r.status === 'ok' ? '#10B981' : '#EF4444', fontWeight:600 }">{{ r.status === 'ok' ? '正常' : (r.detail || '缺失') }}</span>
        </div>
        <div v-if="repairDone" style="padding:8px 12px;font-size:13px;color:#10B981;font-weight:500;">✅ 修复已完成</div>
      </div>
      <div v-else-if="diagnoseStatus === 'idle'" style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);">点击「开始诊断」检查系统状态</div>
    </div>

    <!-- 系统状态 -->
    <div v-if="activeTab === 'status'" class="card" style="max-width:640px;padding:24px;">
      <div v-if="statusLoading" style="text-align:center;padding:24px;">加载中...</div>
      <div v-else-if="sysStatus">
        <div style="margin-bottom:20px;">
          <div style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">版本信息</div>
          <div class="version-grid">
            <div v-for="(val, key) in sysStatus.version" :key="key" class="version-item"><span class="version-key">{{ key }}</span><span class="version-val">{{ val }}</span></div>
          </div>
        </div>
        <div>
          <div style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">迁移记录（{{ sysStatus.migration_count }} 条）</div>
          <div v-if="sysStatus.migrations.length" class="mig-list">
            <div v-for="(m, i) in sysStatus.migrations" :key="i" class="mig-row"><span class="mig-batch">#{{ m.batch }}</span><span class="mig-name">{{ m.migration }}</span></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 实时日志 -->
    <div v-if="activeTab === 'logs'" class="card" style="max-width:100%;padding:20px;">
      <div style="display:flex;align-items:center;gap:12px;margin-bottom:12px;flex-wrap:wrap;">
        <h3 style="font-size:15px;font-weight:700;margin:0;">📋 实时日志</h3>
        <select v-model.number="logLines" style="padding:4px 8px;border-radius:6px;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text);font-size:12px;">
          <option :value="50">50 行</option>
          <option :value="200">200 行</option>
          <option :value="500">500 行</option>
          <option :value="1000">1000 行</option>
        </select>
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="loadLogs" :disabled="logRefreshing">{{ logRefreshing ? '加载中...' : '🔄 刷新' }}</button>
        <button class="btn btn-sm" style="background:#7c3aed;color:#fff;border:none;" @click="startLogPolling">▶ 自动刷新</button>
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="stopLogPolling">⏹ 停止</button>
      </div>
      <div v-if="logError" style="color:#EF4444;font-size:13px;padding:8px;background:rgba(239,68,68,0.06);border-radius:6px;margin-bottom:8px;">{{ logError }}</div>
      <pre style="background:#0d1117;color:#e6edf3;padding:16px;border-radius:10px;font-size:12px;line-height:1.6;overflow:auto;max-height:60vh;white-space:pre-wrap;word-break:break-all;margin:0;">{{ logContent || '暂无日志' }}</pre>
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
.badge-active { display:inline-block; padding:3px 12px; border-radius:20px; font-size:12px; font-weight:600; background:rgba(16,185,129,0.1); color:#10B981; }
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
.mig-batch { display:inline-block; padding:1px 6px; border-radius:4px; background:#ede9fe; color:#7c3aed; font-weight:700; font-size:10px; flex-shrink:0; }
.mig-name { color:var(--color-text-secondary); word-break:break-all; }
.loading-spinner { text-align:center; padding:48px; color:var(--color-text-secondary); font-size:15px; }
.btn { padding:8px 16px; border-radius:8px; font-size:13px; font-weight:500; cursor:pointer; border:1px solid transparent; transition:all 0.15s; }
.btn-sm { padding:6px 14px; font-size:12px; border-radius:8px; }
.btn-primary { background:#7c3aed; color:white; border-color:#7c3aed; }
.btn-primary:hover { background:#6d28d9; }
.btn-outline { background:var(--color-bg-card); color:var(--color-text); border:1px solid var(--color-border); }
.btn-danger { background:var(--color-bg-card); color:#EF4444; border:1px solid rgba(239,68,68,0.2); }
</style>
