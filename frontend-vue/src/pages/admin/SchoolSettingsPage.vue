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
</template>

<style scoped>
.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 12px; font-weight: 600; color: var(--color-text); margin-bottom: 4px; }
.form-input { color: var(--color-text); width: 100%; padding: 8px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 13px; outline: none; transition: border-color 0.15s; box-sizing: border-box; }
.form-input:focus { border-color: #7c3aed; box-shadow: 0 0 0 3px rgba(124,58,237,0.08); }
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; }
.btn-sm { padding: 6px 14px; font-size: 12px; border-radius: 8px; }
.btn-primary { background: #7c3aed; color: white; border-color: #7c3aed; }
.btn-primary:hover { background: #6d28d9; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
