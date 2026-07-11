<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { avatarGradient, platformLabel, escapeHtml } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface Teacher {
  id: number
  name: string
  username: string
  nickname?: string
  phone?: string
  email?: string
  avatar_path?: string
  status: string
  bindings: string[]
  class_names: string[]
}

interface CreatedTeacher {
  username: string
  name: string
  initial_password: string
}

const teachers = ref<Teacher[]>([])
const loading = ref(true)

// 批量创建弹窗状态
const showBatchModal = ref(false)
const batchText = ref('')
const batchLoading = ref(false)
const createdTeachers = ref<CreatedTeacher[]>([])

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Teacher[]>>('/api/v1/admin/teachers')
    teachers.value = res.data || []
  } catch { teachers.value = [] }
  finally { loading.value = false }
})

async function resetPwd(teacher: Teacher) {
  const pwd = prompt(`为「${teacher.name}」设置新密码（留空自动生成）：`)
  try {
    await apiPost(`/api/v1/admin/teachers/${teacher.id}/reset-password`, { password: pwd || undefined })
    toast.show('密码已重置', 'success')
  } catch { /* handled */ }
}

async function deleteTeacher(teacher: Teacher) {
  if (!confirm(`确定删除教师「${teacher.name}」？\n该教师关联的班级将被解除班主任。`)) return
  try {
    await apiDelete(`/api/v1/admin/teachers/${teacher.id}`)
    teachers.value = teachers.value.filter(t => t.id !== teacher.id)
    toast.show('已删除教师：' + teacher.name, 'success')
  } catch { /* handled */ }
}

function openBatchModal() {
  batchText.value = ''
  createdTeachers.value = []
  showBatchModal.value = true
}

function downloadTemplate() {
  const csv = '姓名,密码,手机号,邮箱\n张老师,,13800138001,zhang@example.com\n李老师,,13800138002\n'
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'teachers_template.csv'
  a.click()
  URL.revokeObjectURL(url)
}

async function submitBatchCreate() {
  const lines = batchText.value.trim().split('\n').filter(l => l.trim())
  if (lines.length === 0) {
    toast.show('请输入至少一位教师信息', 'error')
    return
  }

  const teachersData = lines.map(line => {
    const [name, password, phone, email] = line.split(',').map(s => s?.trim() || '')
    const t: Record<string, string> = { name }
    if (password) t.password = password
    if (phone) t.phone = phone
    if (email) t.email = email
    return t
  })

  batchLoading.value = true
  try {
    const res = await apiPost<ApiResponse<CreatedTeacher[]>>('/api/v1/admin/teachers/batch-create', {
      teachers: teachersData,
    })
    createdTeachers.value = (res as unknown as { data: CreatedTeacher[] }).data || []
    toast.show(`成功创建 ${createdTeachers.value.length} 个教师账号`, 'success')
    // 刷新列表
    const listRes = await apiGet<ApiResponse<Teacher[]>>('/api/v1/admin/teachers')
    teachers.value = listRes.data || []
  } catch { /* handled */ }
  finally { batchLoading.value = false }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">账号管理</p>
        <h2 style="font-size:24px;font-weight:700;">教师账号</h2>
      </div>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="downloadTemplate">📥 下载模板</button>
        <button class="btn btn-sm btn-primary" @click="openBatchModal">+ 批量创建</button>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="teachers.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">👨‍🏫</div>
      <p>暂无教师账号</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(340px,1fr));gap:16px;">
      <div v-for="t in teachers" :key="t.id" class="card"
        style="display:flex;align-items:center;gap:16px;padding:16px 24px;transition:transform 0.25s;"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = 'translateY(-2px)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = ''">
        <div :style="{ width:'48px', height:'48px', borderRadius:'14px', background: avatarGradient(t.name), color:'white', display:'flex', alignItems:'center', justifyContent:'center', fontWeight:700, fontSize:'18px' }">
          {{ t.nickname?.[0] || t.name[0] }}
        </div>
        <div style="flex:1;min-width:0;">
          <div style="font-weight:600;font-size:15px;">
            {{ t.name }}
            <span v-if="t.nickname && t.nickname !== t.name" style="color:var(--color-text-secondary);font-size:12px;">@{{ t.nickname }}</span>
          </div>
          <div style="font-size:12px;color:var(--color-text-secondary);">{{ t.username }}{{ t.phone ? ' · ' + t.phone : '' }}</div>
          <div style="display:flex;gap:4px;margin-top:4px;flex-wrap:wrap;">
            <span v-for="b in t.bindings" :key="b" style="font-size:11px;padding:2px 8px;border-radius:10px;background:rgba(79,70,229,0.06);color:var(--color-primary);font-weight:500;">{{ platformLabel(b) }}</span>
            <span v-for="c in t.class_names" :key="c" style="font-size:11px;padding:2px 8px;border-radius:10px;background:rgba(245,158,11,0.1);color:#F59E0B;">🏫 {{ c }}</span>
          </div>
        </div>
        <div style="display:flex;gap:6px;flex-shrink:0;">
          <button class="btn btn-sm" style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);" @click="resetPwd(t)">重置密码</button>
          <button class="btn btn-sm" style="background:#fee2e2;color:#dc2626;border:1px solid #fecaca;" @click="deleteTeacher(t)">删除</button>
        </div>
      </div>
    </div>

    <!-- 批量创建弹窗 -->
    <div v-if="showBatchModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.self="showBatchModal = false">
      <div class="card" style="width:90%;max-width:560px;max-height:85vh;overflow-y:auto;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">批量创建教师账号</h3>
          <button class="btn btn-sm" style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showBatchModal = false">×</button>
        </div>

        <!-- 创建结果 -->
        <div v-if="createdTeachers.length > 0">
          <p style="color:var(--color-text-secondary);font-size:13px;margin-bottom:12px;">已创建 {{ createdTeachers.length }} 个账号，请记录初始密码：</p>
          <div class="data-table" style="margin-bottom:16px;">
            <table>
              <thead><tr><th>姓名</th><th>账号</th><th>初始密码</th></tr></thead>
              <tbody>
                <tr v-for="t in createdTeachers" :key="t.username">
                  <td style="font-weight:600;">{{ t.name }}</td>
                  <td>{{ t.username }}</td>
                  <td style="font-family:monospace;font-weight:600;color:var(--color-primary);">{{ t.initial_password }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <button class="btn btn-primary" style="width:100%;" @click="showBatchModal = false">我已记录</button>
        </div>

        <!-- 输入表单 -->
        <div v-else>
          <p style="color:var(--color-text-secondary);font-size:13px;margin-bottom:12px;">
            每行一位教师，格式：<code>姓名,密码,手机号,邮箱</code><br>
            密码留空自动随机生成，手机号和邮箱可选
          </p>
          <textarea
            v-model="batchText"
            class="form-input"
            style="width:100%;min-height:140px;font-family:monospace;margin-bottom:16px;"
            placeholder="张老师,,13800138001,zhang@example.com&#10;李老师,,13800138002&#10;王老师,"
          ></textarea>
          <div style="display:flex;gap:8px;justify-content:flex-end;">
            <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showBatchModal = false">取消</button>
            <button class="btn btn-sm btn-primary" :disabled="batchLoading" @click="submitBatchCreate">
              {{ batchLoading ? '创建中...' : '创建账号' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
