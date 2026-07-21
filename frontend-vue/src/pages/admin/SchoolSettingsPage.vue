<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
const saving = ref(false)
const form = ref({ name: '', address: '', contact_phone: '', contact_email: '' })
const schoolCode = ref('')
const schoolStatus = ref('')
const logoPath = ref('')
const logoUploading = ref(false)
const activeTab = ref<'school' | 'diagnose' | 'status'>('school')

// 诊断
interface DiagItem { item: string; status: string; detail?: string }
const diagLoading = ref(false)
const diagResult = ref<DiagItem[] | null>(null)
const diagHasIssues = ref(false)
const repairLoading = ref(false)
const repairDone = ref(false)
async function diagnose() {
  diagLoading.value = true; diagResult.value = null; repairDone.value = false
  try {
    const res = await fetch('/api/v1/admin/system/diagnose', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    diagResult.value = data.data || []
    diagHasIssues.value = data.has_issues || false
    toast.show(data.message || '诊断完成', 'info')
  } catch { toast.show('诊断请求失败', 'error') }
  finally { diagLoading.value = false }
}
async function repair() {
  repairLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/system/repair', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' }
    })
    const data = await res.json()
    toast.show(data.message || '修复完成', res.ok ? 'success' : 'error')
    if (res.ok) { repairDone.value = true; setTimeout(diagnose, 500) }
  } catch { toast.show('修复请求失败', 'error') }
  finally { repairLoading.value = false }
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
  if (!form.value.name.trim()) { toast.show('请填写学校名称', 'error'); return }
  saving.value = true
  try {
    const res = await fetch('/api/v1/admin/school', {
      method: 'PUT',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: form.value.name.trim(), address: form.value.address.trim(), contact_phone: form.value.contact_phone.trim(), contact_email: form.value.contact_email.trim() }),
    })
    const data = await res.json()
    toast.show(data.message || '学校信息已保存', res.ok ? 'success' : 'error')
  } catch { toast.show('保存失败', 'error') }
  finally { saving.value = false }
}

async function reload() {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<School>>('/api/v1/admin/school')
    const s = (res as unknown as { data: School }).data
    form.value = { name: s.name || '', address: s.address || '', contact_phone: s.contact_phone || '', contact_email: s.contact_email || '' }
    logoPath.value = s.logo_path || ''
    toast.show('已恢复原始数据', 'success')
  } catch { /* handled */ }
  finally { loading.value = false }
}

async function uploadLogo(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  logoUploading.value = true
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
    toast.show(data.message || 'LOGO 已上传', res.ok ? 'success' : 'error')
  } catch { toast.show('上传失败', 'error') }
  finally { logoUploading.value = false; (e.target as HTMLInputElement).value = '' }
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
            <label class="btn btn-sm btn-outline" style="cursor:pointer;">
              {{ logoUploading ? '上传中...' : '📷 上传 LOGO' }}
              <input type="file" accept="image/*" style="display:none;" @change="uploadLogo">
            </label>
          </div>
        </div>
        <div class="form-group"><label>学校名称</label><input v-model="form.name" class="form-input" placeholder="请输入学校名称"></div>
        <div class="form-group"><label>学校地址</label><input v-model="form.address" class="form-input" placeholder="请输入学校地址"></div>
        <div class="form-group"><label>联系电话</label><input v-model="form.contact_phone" class="form-input" placeholder="如：021-12345678"></div>
        <div class="form-group"><label>联系邮箱</label><input v-model="form.contact_email" type="email" class="form-input" placeholder="如：admin@school.edu.cn"></div>
        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:24px;">
          <button class="btn btn-sm btn-outline" :disabled="saving || loading" @click="reload">↩️ 重置</button>
          <button class="btn btn-sm btn-primary" :disabled="saving || loading" @click="save">{{ saving ? '保存中...' : '保存设置' }}</button>
        </div>
      </div>
    </div>

    <!-- 系统诊断 -->
    <div v-if="activeTab === 'diagnose'" class="card" style="max-width:640px;padding:24px;">
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">检查数据库表结构完整性</p>
      <div style="display:flex;gap:12px;margin-bottom:16px;">
        <button class="btn btn-outline" :disabled="diagLoading" @click="diagnose">{{ diagLoading ? '诊断中...' : '🔍 开始诊断' }}</button>
        <button class="btn btn-danger" :disabled="repairLoading || !diagHasIssues" @click="repair">{{ repairLoading ? '修复中...' : '🛠️ 一键修复' }}</button>
      </div>
      <div v-if="diagResult">
        <div v-for="(r, i) in diagResult" :key="i" style="display:flex;align-items:center;gap:10px;padding:6px 10px;border-bottom:1px solid var(--color-border);font-size:13px;">
          <span v-if="r.status === 'ok'" style="color:#10B981;">✅</span><span v-else style="color:#EF4444;">❌</span>
          <span style="flex:1;">{{ r.item }}</span>
          <span :style="{ color: r.status === 'ok' ? '#10B981' : '#EF4444', fontWeight:600 }">{{ r.status === 'ok' ? '正常' : (r.detail || '缺失') }}</span>
        </div>
        <div v-if="repairDone" style="padding:8px 12px;font-size:13px;color:#10B981;font-weight:500;">✅ 修复已完成</div>
      </div>
      <div v-else-if="!diagLoading" style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);">点击「开始诊断」检查系统状态</div>
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
