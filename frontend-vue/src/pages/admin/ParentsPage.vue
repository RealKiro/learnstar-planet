<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { avatarGradient, platformLabel } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface ParentChild {
  student_name: string
  class_name: string
}
interface Parent {
  id: number
  name: string
  username: string
  phone?: string
  bindings?: string[]
  children?: ParentChild[]
}
interface CreatedParent {
  username: string
  name: string
  initial_password: string
}

const parents = ref<Parent[]>([])
const loading = ref(true)

// 批量创建弹窗
const showBatchModal = ref(false)
const batchText = ref('')
const deleteStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const deleteTargetId = ref<number | null>(null)
const batchCreateStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const createdParents = ref<CreatedParent[]>([])

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Parent[]>>('/api/v1/admin/parents')
    parents.value = res.data || []
  } catch { parents.value = [] }
  finally { loading.value = false }
})

async function reloadParents() {
  try {
    const res = await apiGet<ApiResponse<Parent[]>>('/api/v1/admin/parents')
    parents.value = res.data || []
  } catch { /* handled */ }
}

async function deleteParent(p: Parent) {
  if (!confirm(`确定删除家长账号「${p.name}」？\n该账号与学生的绑定关系将被解除。`)) return
  deleteStatus.value = 'loading'
  deleteTargetId.value = p.id
  try {
    await apiDelete(`/api/v1/admin/parents/${p.id}`)
    deleteStatus.value = 'success'
    setTimeout(() => {
      parents.value = parents.value.filter(x => x.id !== p.id)
      deleteStatus.value = 'idle'
      deleteTargetId.value = null
    }, 1500)
  } catch {
    deleteStatus.value = 'error'
    setTimeout(() => { deleteStatus.value = 'idle'; deleteTargetId.value = null }, 3000)
  }
}

function openBatchModal() {
  batchText.value = ''
  createdParents.value = []
  showBatchModal.value = true
}

async function submitBatchCreate() {
  const lines = batchText.value.trim().split('\n').filter(l => l.trim())
  if (lines.length === 0) {
    toast.show('请输入至少一位家长信息', 'error', { position: 'top-right' })
    return
  }
  const parentsData = lines.map(line => {
    const [name, password, phone] = line.split(',').map(s => s?.trim() || '')
    const p: Record<string, string> = { name }
    if (password) p.password = password
    if (phone) p.phone = phone
    return p
  })
  batchCreateStatus.value = 'loading'
  try {
    const res = await apiPost<ApiResponse<CreatedParent[]>>('/api/v1/admin/parents/batch-create', {
      parents: parentsData,
    })
    createdParents.value = (res as unknown as { data: CreatedParent[] }).data || []
    batchCreateStatus.value = 'success'
    await reloadParents()
    setTimeout(() => { batchCreateStatus.value = 'idle' }, 1500)
  } catch {
    batchCreateStatus.value = 'error'
    setTimeout(() => { batchCreateStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">账号管理</p>
        <h2 style="font-size:24px;font-weight:700;">家长账号</h2>
      </div>
      <button class="btn btn-sm btn-primary" @click="openBatchModal">+ 批量创建</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="parents.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">👨‍👩‍👧</div>
      <p>暂无家长账号</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(340px,1fr));gap:16px;">
      <div v-for="p in parents" :key="p.id" class="card"
        style="display:flex;align-items:flex-start;gap:16px;padding:16px 24px;transition:transform 0.25s;"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = 'translateY(-2px)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = ''">
        <div :style="{ width:'48px', height:'48px', borderRadius:'14px', background: avatarGradient(p.name), color:'white', display:'flex', alignItems:'center', justifyContent:'center', fontWeight:700, fontSize:'18px', flexShrink:0 }">
          {{ p.name[0] }}
        </div>
        <div style="flex:1;min-width:0;">
          <div style="font-weight:600;font-size:15px;">{{ p.name }}</div>
          <div style="font-size:12px;color:var(--color-text-secondary);">{{ p.username }}{{ p.phone ? ' · ' + p.phone : '' }}</div>
          <div style="display:flex;gap:4px;margin-top:4px;flex-wrap:wrap;">
            <span v-for="b in (p.bindings || [])" :key="b" style="font-size:11px;padding:2px 8px;border-radius:10px;background:rgba(79,70,229,0.06);color:var(--color-primary);font-weight:500;">{{ platformLabel(b) }}</span>
          </div>
          <div v-if="p.children && p.children.length" style="display:flex;gap:4px;margin-top:6px;flex-wrap:wrap;">
            <span v-for="(c, i) in p.children" :key="i" style="font-size:11px;padding:2px 8px;border-radius:10px;background:rgba(16,185,129,0.1);color:#10B981;">👦 {{ c.student_name }} · {{ c.class_name }}</span>
          </div>
        </div>
        <button class="btn btn-sm" :style="{ background: deleteStatus !== 'idle' && deleteTargetId === p.id ? (deleteStatus === 'loading' ? '#f59e0b' : deleteStatus === 'success' ? '#10b981' : '#ef4444') : '#fee2e2', color: deleteStatus !== 'idle' && deleteTargetId === p.id ? '#fff' : '#dc2626', border: deleteStatus !== 'idle' && deleteTargetId === p.id ? '1px solid transparent' : '1px solid #fecaca', flexShrink: 0 }" @click="deleteParent(p)">
          <template v-if="deleteStatus === 'loading' && deleteTargetId === p.id">删除中...</template>
          <template v-else-if="deleteStatus === 'success' && deleteTargetId === p.id">已删除</template>
          <template v-else-if="deleteStatus === 'error' && deleteTargetId === p.id">失败</template>
          <template v-else>删除</template>
        </button>
      </div>
    </div>

    <!-- 批量创建弹窗 -->
    <div v-if="showBatchModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.self="showBatchModal = false">
      <div class="card" style="width:90%;max-width:560px;max-height:85vh;overflow-y:auto;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">批量创建家长账号</h3>
          <button class="btn btn-sm" style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showBatchModal = false">×</button>
        </div>

        <!-- 创建结果 -->
        <div v-if="createdParents.length > 0">
          <p style="color:var(--color-text-secondary);font-size:13px;margin-bottom:12px;">已创建 {{ createdParents.length }} 个账号，请记录初始密码：</p>
          <div class="data-table" style="margin-bottom:16px;">
            <table>
              <thead><tr><th>姓名</th><th>账号</th><th>初始密码</th></tr></thead>
              <tbody>
                <tr v-for="p in createdParents" :key="p.username">
                  <td style="font-weight:600;">{{ p.name }}</td>
                  <td>{{ p.username }}</td>
                  <td style="font-family:monospace;font-weight:600;color:var(--color-primary);">{{ p.initial_password }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <button class="btn btn-primary" style="width:100%;" @click="showBatchModal = false">我已记录</button>
        </div>

        <!-- 输入表单 -->
        <div v-else>
          <p style="color:var(--color-text-secondary);font-size:13px;margin-bottom:12px;">
            每行一位家长，格式：<code>姓名,密码,手机号</code><br>
            密码留空自动随机生成，手机号可选
          </p>
          <textarea
            v-model="batchText"
            class="form-input"
            style="width:100%;min-height:140px;font-family:monospace;margin-bottom:16px;"
            placeholder="张小明家长,,13800138001&#10;李小红家长,,13800138002&#10;王刚家长,"
          ></textarea>
          <div style="display:flex;gap:8px;justify-content:flex-end;">
            <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showBatchModal = false">取消</button>
            <button class="btn btn-sm" :style="{ background: batchCreateStatus === 'loading' ? '#f59e0b' : batchCreateStatus === 'success' ? '#10b981' : batchCreateStatus === 'error' ? '#ef4444' : '#7c3aed', color: '#fff', border: '1px solid transparent' }" :disabled="batchCreateStatus !== 'idle'" @click="submitBatchCreate">
              <template v-if="batchCreateStatus === 'loading'">创建中...</template>
              <template v-else-if="batchCreateStatus === 'success'">已创建 ✓</template>
              <template v-else-if="batchCreateStatus === 'error'">失败 ✗</template>
              <template v-else>创建账号</template>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
