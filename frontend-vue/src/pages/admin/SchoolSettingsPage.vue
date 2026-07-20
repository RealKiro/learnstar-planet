<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface School {
  name: string; code: string; address: string
  contact_phone: string; contact_email: string
  settings?: unknown; status: string
}
interface DiagItem { item: string; status: string; detail?: string }
interface SystemStatus {
  version: Record<string, string>
  migrations: { migration: string; batch: number }[]
  migration_count: number
}

const activeTab = ref<'school' | 'demo' | 'diagnose' | 'status'>('school')

// ===== 学校信息 =====
const loading = ref(true)
const saving = ref(false)
const form = ref({ name: '', address: '', contact_phone: '', contact_email: '' })
const schoolCode = ref('')
const schoolStatus = ref('')

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<School>>('/api/v1/admin/school')
    const s = (res as unknown as { data: School }).data
    form.value = {
      name: s.name || '', address: s.address || '',
      contact_phone: s.contact_phone || '', contact_email: s.contact_email || '',
    }
    schoolCode.value = s.code || ''
    schoolStatus.value = s.status || ''
  } catch { /* handled */ }
  finally { loading.value = false }
})

async function save() {
  if (!form.value.name.trim()) { toast.show('请填写学校名称', 'error'); return }
  saving.value = true
  try {
    await apiPut('/api/v1/admin/school', {
      name: form.value.name.trim(), address: form.value.address.trim(),
      contact_phone: form.value.contact_phone.trim(), contact_email: form.value.contact_email.trim(),
    })
    toast.show('学校信息已保存', 'success')
  } catch { /* handled */ }
  finally { saving.value = false }
}

async function reload() {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<School>>('/api/v1/admin/school')
    const s = (res as unknown as { data: School }).data
    form.value = {
      name: s.name || '', address: s.address || '',
      contact_phone: s.contact_phone || '', contact_email: s.contact_email || '',
    }
    toast.show('已恢复原始数据', 'success')
  } catch { /* handled */ }
  finally { loading.value = false }
}

// ===== 演示数据 =====
const demoLoading = ref(false)

async function seedDemo() {
  demoLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/demo/seed', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' }
    })
    const data = await res.json()
    toast.show(data.message || '演示数据已生成', res.ok ? 'success' : 'error')
  } catch { toast.show('操作失败', 'error') }
  finally { demoLoading.value = false }
}
async function cleanDemo() {
  demoLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/demo/clean', {
      method: 'POST',
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token'), 'Content-Type': 'application/json' }
    })
    const data = await res.json()
    toast.show(data.message || '演示数据已清除', res.ok ? 'success' : 'error')
  } catch { toast.show('操作失败', 'error') }
  finally { demoLoading.value = false }
}

// ===== 系统诊断 =====
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
    toast.show(data.message || '诊断完成', res.ok ? 'info' : 'error')
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

// ===== 系统状态 =====
const statusLoading = ref(false)
const sysStatus = ref<SystemStatus | null>(null)

async function loadStatus() {
  statusLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/system/status', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const data = await res.json()
    sysStatus.value = data.data || null
  } catch { toast.show('获取系统状态失败', 'error') }
  finally { statusLoading.value = false }
}

function onTabChange(tab: typeof activeTab.value) {
  activeTab.value = tab
  if (tab === 'status' && !sysStatus.value) loadStatus()
}
</script>

<template>
  <div class="settings-page">
    <!-- 标签导航 -->
    <div class="tab-bar">
      <button v-for="t in ([
        { key: 'school', label: '🏫 学校信息' },
        { key: 'demo', label: '🧪 演示数据' },
        { key: 'diagnose', label: '🔧 系统诊断' },
        { key: 'status', label: '📊 系统状态' },
      ] as const)" :key="t.key"
        :class="['tab-btn', { active: activeTab === t.key }]"
        @click="onTabChange(t.key)">
        {{ t.label }}
      </button>
    </div>

    <!-- ===== 学校信息 ===== -->
    <div v-if="activeTab === 'school'" class="tab-content">
      <div v-if="loading" class="loading-spinner">加载中...</div>
      <div v-else class="card" style="max-width:640px;padding:32px;">
        <div style="display:flex;gap:16px;margin-bottom:24px;flex-wrap:wrap;">
          <div class="info-block">
            <div class="info-label">学校编码</div>
            <div class="info-value mono">{{ schoolCode || '-' }}</div>
          </div>
          <div class="info-block">
            <div class="info-label">状态</div>
            <div><span class="badge-active">{{ schoolStatus === 'active' ? '正常运行' : (schoolStatus || '正常') }}</span></div>
          </div>
        </div>
        <div class="form-group"><label>学校名称</label><input v-model="form.name" class="form-input" placeholder="请输入学校名称"></div>
        <div class="form-group"><label>学校地址</label><input v-model="form.address" class="form-input" placeholder="请输入学校地址"></div>
        <div class="form-group"><label>联系电话</label><input v-model="form.contact_phone" class="form-input" placeholder="如：021-12345678"></div>
        <div class="form-group"><label>联系邮箱</label><input v-model="form.contact_email" type="email" class="form-input" placeholder="如：admin@school.edu.cn"></div>
        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:24px;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" :disabled="saving || loading" @click="reload">↩️ 重置</button>
          <button class="btn btn-sm btn-primary" :disabled="saving" @click="save">{{ saving ? '保存中...' : '保存设置' }}</button>
        </div>
      </div>
    </div>

    <!-- ===== 演示数据 ===== -->
    <div v-if="activeTab === 'demo'" class="tab-content">
      <div class="card" style="max-width:640px;padding:32px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:4px;">🧪 演示数据管理</h3>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
          生成演示数据用于测试体验。管理员: demo_admin / demo123 · 教师: demo_t1~t4 / demo123 · 班级码: DEMO00
        </p>
        <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:16px;padding:8px 12px;background:var(--color-bg);border-radius:8px;">
          ⚠️ 如果已存在演示数据，重新生成会先清除旧数据再创建新数据。
        </p>
        <div style="display:flex;gap:12px;">
          <button class="btn btn-sm btn-primary" :disabled="demoLoading" @click="seedDemo">{{ demoLoading ? '处理中...' : '📥 生成演示数据' }}</button>
          <button class="btn btn-sm" :disabled="demoLoading" @click="cleanDemo"
            style="background:var(--color-bg-card);color:var(--color-danger);border:1px solid rgba(239,68,68,0.2);">
            {{ demoLoading ? '处理中...' : '🗑️ 清除演示数据' }}
          </button>
        </div>
      </div>
    </div>

    <!-- ===== 系统诊断 ===== -->
    <div v-if="activeTab === 'diagnose'" class="tab-content">
      <div class="card" style="max-width:640px;padding:32px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:4px;">🔧 系统诊断与修复</h3>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
          检查数据库结构完整性，如果检测到缺失字段可一键修复
        </p>
        <div style="display:flex;gap:12px;margin-bottom:16px;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" :disabled="diagLoading" @click="diagnose">
            {{ diagLoading ? '诊断中...' : '🔍 开始诊断' }}
          </button>
          <button class="btn btn-sm" :disabled="repairLoading || !diagHasIssues" @click="repair"
            style="background:var(--color-bg-card);color:var(--color-danger);border:1px solid rgba(239,68,68,0.2);">
            {{ repairLoading ? '修复中...' : '🛠️ 一键修复' }}
          </button>
        </div>
        <div v-if="diagResult" style="border:1px solid var(--color-border);border-radius:8px;overflow:hidden;">
          <div v-for="(r, i) in diagResult" :key="i" class="diag-row">
            <span v-if="r.status === 'ok'" style="color:#10B981;">✅</span>
            <span v-else-if="r.status === 'error'" style="color:#F59E0B;">⚠️</span>
            <span v-else style="color:#EF4444;">❌</span>
            <span style="flex:1;">{{ r.item }}</span>
            <span :style="{ color: r.status === 'ok' ? '#10B981' : '#EF4444', fontWeight: 600 }">
              {{ r.status === 'ok' ? '正常' : (r.detail || '缺失') }}
            </span>
          </div>
          <div v-if="repairDone" class="diag-success">✅ 修复已完成，数据库结构已更新</div>
        </div>
        <div v-else-if="!diagLoading" class="diag-placeholder">点击「开始诊断」检查系统状态</div>
      </div>
    </div>

    <!-- ===== 系统状态 ===== -->
    <div v-if="activeTab === 'status'" class="tab-content">
      <div class="card" style="max-width:640px;padding:32px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">📊 系统状态</h3>
        <div v-if="statusLoading" class="loading-spinner">加载中...</div>
        <div v-else-if="sysStatus">
          <div style="margin-bottom:20px;">
            <h4 style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">版本信息</h4>
            <div class="status-grid">
              <div v-for="(val, key) in sysStatus.version" :key="key" class="status-item">
                <span class="status-key">{{ key }}</span>
                <span class="status-val">{{ val }}</span>
              </div>
            </div>
          </div>
          <div>
            <h4 style="font-size:13px;font-weight:600;color:var(--color-text-secondary);margin-bottom:8px;">
              迁移记录（共 {{ sysStatus.migration_count }} 条）
            </h4>
            <div v-if="sysStatus.migrations.length" class="mig-list">
              <div v-for="(m, i) in sysStatus.migrations" :key="i" class="mig-row">
                <span class="mig-batch">#{{ m.batch }}</span>
                <span class="mig-name">{{ m.migration }}</span>
              </div>
            </div>
            <div v-else style="padding:12px;text-align:center;color:var(--color-text-secondary);font-size:13px;">
              {{ typeof sysStatus.migrations === 'object' ? JSON.stringify(sysStatus.migrations) : '无迁移记录' }}
            </div>
          </div>
        </div>
        <div v-else class="diag-placeholder">加载失败，请刷新重试</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page { max-width: 720px; margin: 0 auto; padding: 0 4px; }

/* 标签导航 */
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }

/* 信息块 */
.info-block { flex: 1; min-width: 200px; padding: 12px 16px; background: var(--color-bg); border-radius: 10px; }
.info-label { font-size: 12px; color: var(--color-text-secondary); margin-bottom: 4px; }
.info-value { font-weight: 600; }
.info-value.mono { font-family: monospace; }
.badge-active { display: inline-block; padding: 3px 12px; border-radius: 20px; font-size: 12px; font-weight: 600; background: rgba(16,185,129,0.1); color: #10B981; }

/* 诊断 */
.diag-row { display: flex; align-items: center; gap: 10px; padding: 8px 14px; border-bottom: 1px solid var(--color-border); font-size: 13px; }
.diag-row:last-child { border-bottom: none; }
.diag-success { padding: 10px 14px; background: rgba(16,185,129,0.05); font-size: 13px; color: #10B981; font-weight: 500; }
.diag-placeholder { padding: 12px; background: var(--color-bg); border-radius: 8px; text-align: center; font-size: 13px; color: var(--color-text-secondary); }

/* 状态 */
.status-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; }
.status-item { display: flex; justify-content: space-between; padding: 6px 12px; background: var(--color-bg); border-radius: 6px; font-size: 12px; }
.status-key { color: var(--color-text-secondary); font-weight: 500; }
.status-val { font-family: monospace; font-weight: 600; }
.mig-list { max-height: 360px; overflow-y: auto; border: 1px solid var(--color-border); border-radius: 8px; }
.mig-row { display: flex; align-items: center; gap: 8px; padding: 6px 12px; border-bottom: 1px solid var(--color-border); font-size: 11px; }
.mig-row:last-child { border-bottom: none; }
.mig-batch { display: inline-block; padding: 1px 6px; border-radius: 4px; background: #ede9fe; color: #7c3aed; font-weight: 700; font-size: 10px; flex-shrink: 0; }
.mig-name { color: var(--color-text-secondary); word-break: break-all; }
.loading-spinner { text-align: center; padding: 48px; color: var(--color-text-secondary); font-size: 15px; }
.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 12px; font-weight: 600; color: var(--color-text); margin-bottom: 4px; }
.form-input { color: var(--color-text); width: 100%; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 13px; outline: none; transition: border-color 0.15s; box-sizing: border-box; background: var(--color-bg-card); }
.form-input:focus { border-color: #7c3aed; box-shadow: 0 0 0 3px rgba(124,58,237,0.08); }
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; }
.btn-sm { padding: 6px 14px; font-size: 12px; border-radius: 8px; }
.btn-primary { background: #7c3aed; color: white; border-color: #7c3aed; }
.btn-primary:hover { background: #6d28d9; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
