<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref<'demo' | 'diagnose' | 'status'>('demo')

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

// ===== 系统状态 =====
interface SystemStatus {
  version: Record<string, string>
  migrations: { migration: string; batch: number }[]
  migration_count: number
}
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
  } catch { /* errors handled by interceptor */ }
  finally { statusLoading.value = false }
}

function onTabChange(tab: typeof activeTab.value) {
  activeTab.value = tab
  if (tab === 'status' && !sysStatus.value) loadStatus()
}
</script>

<template>
  <div class="debug-page">
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">调试工具</p>
        <h2 style="font-size:24px;font-weight:700;">🔧 高级调试</h2>
      </div>
    </div>

    <!-- 标签导航 -->
    <div class="tab-bar">
      <button v-for="t in ([
        { key: 'diagnose', label: '🔍 系统诊断' },
        { key: 'status', label: '📊 系统状态' },
      ] as const)" :key="t.key"
        :class="['tab-btn', { active: activeTab === t.key }]"
        @click="onTabChange(t.key)">
        {{ t.label }}
      </button>
    </div>

    <div v-if="activeTab === 'diagnose'" class="tab-content">
      <div class="card" style="max-width:640px;padding:32px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:4px;">🔍 系统诊断</h3>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
          检查数据库表结构完整性，检测到缺失字段可一键修复
        </p>
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

        <div v-if="diagResult" class="diag-list">
          <div v-for="(r, i) in diagResult" :key="i" class="diag-row">
            <span v-if="r.status === 'ok'" class="diag-icon ok">✅</span>
            <span v-else-if="r.status === 'error'" class="diag-icon warn">⚠️</span>
            <span v-else class="diag-icon fail">❌</span>
            <span class="diag-item-name">{{ r.item }}</span>
            <span :class="['diag-status', r.status === 'ok' ? 'ok' : 'fail']">
              {{ r.status === 'ok' ? '正常' : (r.detail || '缺失') }}
            </span>
          </div>
          <div v-if="repairDone" class="diag-success">✅ 修复已完成，数据库结构已更新</div>
        </div>
        <div v-else-if="diagnoseStatus === 'idle'" class="diag-placeholder">
          点击「开始诊断」检查系统状态
        </div>
      </div>
    </div>

    <!-- ===== 系统状态 ===== -->
    <div v-if="activeTab === 'status'" class="tab-content">
      <div class="card" style="max-width:640px;padding:32px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">📊 系统状态</h3>
        <div v-if="statusLoading" class="loading">加载中...</div>
        <div v-else-if="sysStatus">
          <div style="margin-bottom:24px;">
            <h4 class="section-title">版本信息</h4>
            <div class="status-grid">
              <div v-for="(val, key) in sysStatus.version" :key="key" class="status-item">
                <span class="status-key">{{ key }}</span>
                <span class="status-val">{{ val }}</span>
              </div>
            </div>
          </div>
          <div>
            <h4 class="section-title">迁移记录（共 {{ sysStatus.migration_count }} 条）</h4>
            <div v-if="sysStatus.migrations.length" class="mig-list">
              <div v-for="(m, i) in sysStatus.migrations" :key="i" class="mig-row">
                <span class="mig-batch">#{{ m.batch }}</span>
                <span class="mig-name">{{ m.migration }}</span>
              </div>
            </div>
            <div v-else style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);">
              暂无迁移记录
            </div>
          </div>
        </div>
        <div v-else class="diag-placeholder">加载失败，请刷新重试</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.debug-page { max-width: 720px; margin: 0 auto; padding: 0 4px; }

/* 标签导航 */
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }

/* 演示数据 */
.demo-info { margin-bottom: 12px; }
.demo-info-item { display: flex; align-items: center; gap: 8px; padding: 4px 0; font-size: 13px; }
.demo-label { color: var(--color-text-secondary); min-width: 64px; }
.demo-info-item code { font-family: monospace; background: var(--color-bg); padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.demo-warning { font-size: 12px; color: var(--color-text-secondary); margin-bottom: 16px; padding: 8px 12px; background: var(--color-bg); border-radius: 8px; }

/* 诊断 */
.diag-list { border: 1px solid var(--color-border); border-radius: 8px; overflow: hidden; }
.diag-row { display: flex; align-items: center; gap: 10px; padding: 8px 14px; border-bottom: 1px solid var(--color-border); font-size: 13px; }
.diag-row:last-child { border-bottom: none; }
.diag-icon { flex-shrink: 0; }
.diag-item-name { flex: 1; }
.diag-status { font-weight: 600; }
.diag-status.ok { color: #10B981; }
.diag-status.fail { color: #EF4444; }
.diag-success { padding: 10px 14px; background: rgba(16,185,129,0.05); font-size: 13px; color: #10B981; font-weight: 500; }
.diag-placeholder { padding: 12px; background: var(--color-bg); border-radius: 8px; text-align: center; font-size: 13px; color: var(--color-text-secondary); }

/* 系统状态 */
.section-title { font-size: 13px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 8px; }
.status-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; }
.status-item { display: flex; justify-content: space-between; padding: 6px 12px; background: var(--color-bg); border-radius: 6px; font-size: 12px; }
.status-key { color: var(--color-text-secondary); font-weight: 500; }
.status-val { font-family: monospace; font-weight: 600; }
.mig-list { max-height: 400px; overflow-y: auto; border: 1px solid var(--color-border); border-radius: 8px; }
.mig-row { display: flex; align-items: center; gap: 8px; padding: 6px 12px; border-bottom: 1px solid var(--color-border); font-size: 11px; }
.mig-row:last-child { border-bottom: none; }
.mig-batch { display: inline-block; padding: 1px 6px; border-radius: 4px; background: #ede9fe; color: #7c3aed; font-weight: 700; font-size: 10px; flex-shrink: 0; }
.mig-name { color: var(--color-text-secondary); word-break: break-all; }
.loading { text-align: center; padding: 48px; color: var(--color-text-secondary); font-size: 15px; }

/* 按钮 */
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; display: inline-flex; align-items: center; gap: 4px; }
.btn-primary { background: #7c3aed; color: white; border-color: #7c3aed; }
.btn-primary:hover { background: #6d28d9; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-danger { background: var(--color-bg-card); color: #EF4444; border: 1px solid rgba(239,68,68,0.2); }
.btn-danger:hover { background: #FEF2F2; }
.btn-danger:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-outline { background: var(--color-bg-card); color: var(--color-text); border: 1px solid var(--color-border); }
.btn-outline:hover { background: #f9fafb; }
</style>
