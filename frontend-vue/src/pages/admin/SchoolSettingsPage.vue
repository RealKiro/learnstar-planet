<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface School {
  name: string
  code: string
  address: string
  contact_phone: string
  contact_email: string
  settings?: unknown
  status: string
}

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
      name: s.name || '',
      address: s.address || '',
      contact_phone: s.contact_phone || '',
      contact_email: s.contact_email || '',
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
      name: form.value.name.trim(),
      address: form.value.address.trim(),
      contact_phone: form.value.contact_phone.trim(),
      contact_email: form.value.contact_email.trim(),
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
      name: s.name || '',
      address: s.address || '',
      contact_phone: s.contact_phone || '',
      contact_email: s.contact_email || '',
    }
    toast.show('已恢复原始数据', 'success')
  } catch { /* handled */ }
  finally { loading.value = false }
}

const demoLoading = ref(false)

async function seedDemo() {
  demoLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/demo/seed', { method: 'POST', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('auth_token'), 'Content-Type': 'application/json' } })
    const data = await res.json()
    toast.show(data.message || '演示数据已生成', res.ok ? 'success' : 'error')
  } catch { toast.show('操作失败', 'error') }
  finally { demoLoading.value = false }
}
async function cleanDemo() {
  demoLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/demo/clean', { method: 'POST', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('auth_token'), 'Content-Type': 'application/json' } })
    const data = await res.json()
    toast.show(data.message || '演示数据已清除', res.ok ? 'success' : 'error')
  } catch { toast.show('操作失败', 'error') }
  finally { demoLoading.value = false }
}

</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">系统配置</p>
        <h2 style="font-size:24px;font-weight:700;">学校设置</h2>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else class="card" style="max-width:640px;padding:32px;">
      <!-- 只读信息 -->
      <div style="display:flex;gap:16px;margin-bottom:24px;flex-wrap:wrap;">
        <div style="flex:1;min-width:200px;padding:12px 16px;background:var(--color-bg);border-radius:10px;">
          <div style="font-size:12px;color:var(--color-text-secondary);margin-bottom:4px;">学校编码</div>
          <div style="font-family:monospace;font-weight:600;">{{ schoolCode || '-' }}</div>
        </div>
        <div style="flex:1;min-width:200px;padding:12px 16px;background:var(--color-bg);border-radius:10px;">
          <div style="font-size:12px;color:var(--color-text-secondary);margin-bottom:4px;">状态</div>
          <div>
            <span style="display:inline-block;padding:3px 12px;border-radius:20px;font-size:12px;font-weight:600;background:rgba(16,185,129,0.1);color:#10B981;">
              {{ schoolStatus === 'active' ? '正常运行' : (schoolStatus || '正常') }}
            </span>
          </div>
        </div>
      </div>

      <!-- 可编辑表单 -->
      <div class="form-group">
        <label>学校名称</label>
        <input v-model="form.name" class="form-input" placeholder="请输入学校名称">
      </div>
      <div class="form-group">
        <label>学校地址</label>
        <input v-model="form.address" class="form-input" placeholder="请输入学校地址">
      </div>
      <div class="form-group">
        <label>联系电话</label>
        <input v-model="form.contact_phone" class="form-input" placeholder="如：021-12345678">
      </div>
      <div class="form-group">
        <label>联系邮箱</label>
        <input v-model="form.contact_email" type="email" class="form-input" placeholder="如：admin@school.edu.cn">
      </div>

      <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:24px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" :disabled="saving || loading" @click="reload">↩️ 重置</button>
        <button class="btn btn-sm btn-primary" :disabled="saving" @click="save">{{ saving ? '保存中...' : '保存设置' }}</button>
      </div>
    </div>
  </div>

    <!-- 演示数据管理 -->
    <div class="card" style="max-width:640px;padding:32px;margin-top:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:4px;">🧪 演示数据管理</h3>
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
        管理员: demo_admin / demo123 · 教师: demo_t1~t4 / demo123 · 班级码: DEMO00
      </p>
      <div style="display:flex;gap:12px;">
        <button class="btn btn-sm btn-primary" :disabled="demoLoading" @click="seedDemo">{{ demoLoading ? '处理中...' : '📥 生成演示数据' }}</button>
        <button class="btn btn-sm" :disabled="demoLoading" @click="cleanDemo"
          style="background:var(--color-bg-card);color:var(--color-danger);border:1px solid rgba(239,68,68,0.2);">
          {{ demoLoading ? '处理中...' : '🗑️ 清除演示数据' }}
        </button>
      </div>
    </div>

</template>
