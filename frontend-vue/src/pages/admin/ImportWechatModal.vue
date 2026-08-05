<script setup lang="ts">
import { ref } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import ModalGlass from '@/components/common/ModalGlass.vue'

const props = defineProps<{
  visible: boolean
}>()
const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'imported'): void
}>()

interface WecomMember {
  userid: string
  name: string
  mobile: string
  email: string
  department_names: string[]
}
interface ClassItem { id: number; name: string }
interface Row {
  userid: string
  name: string
  mobile: string
  email: string
  department_names: string[]
  role: 'teacher' | 'student'
  class_id: number | ''
  selected: boolean
}

const loading = ref(false)
const importing = ref(false)
const loadError = ref('')
const importError = ref('')
const rows = ref<Row[]>([])
const classes = ref<ClassItem[]>([])

function close() { emit('update:visible', false) }

function matchClassByDept(deptNames: string[], clsList: ClassItem[]): number | '' {
  for (const n of deptNames) {
    const hit = clsList.find(c => c.name === n)
    if (hit) return hit.id
  }
  return ''
}

async function loadContacts() {
  loading.value = true
  try {
    const cRes = await apiGet<{ data: { members: WecomMember[] } }>('/api/v1/admin/third-party/contacts')
    const clsRes = await apiGet<{ data: ClassItem[] }>('/api/v1/admin/classes')
    classes.value = clsRes.data || []
    const members = cRes.data?.members || []
    rows.value = members.map(m => ({
      userid: m.userid,
      name: m.name,
      mobile: m.mobile || '',
      email: m.email || '',
      department_names: m.department_names || [],
      role: 'teacher',
      class_id: matchClassByDept(m.department_names || [], classes.value),
      selected: true,
    }))
    if (!members.length) {
      loadError.value = '通讯录为空，请检查企业微信配置'
    }
  } catch (e: any) {
    loadError.value = e?.response?.data?.message || '拉取通讯录失败（请确认已配置企业微信 secret）'
  } finally {
    loading.value = false
  }
}

function toggleAll() {
  const all = rows.value.every(r => r.selected)
  rows.value.forEach(r => { r.selected = !all })
}

async function doImport() {
  const picked = rows.value.filter(r => r.selected)
  if (!picked.length) { importError.value = '请至少勾选一名成员'; return }
  const teachers = picked.filter(r => r.role === 'teacher').map(r => ({ name: r.name, mobile: r.mobile, email: r.email }))
  const students = picked.filter(r => r.role === 'student').map(r => ({ name: r.name, class_id: r.class_id }))
  if (students.some(s => !s.class_id)) { importError.value = '请为所有学生选择目标班级'; return }

  importing.value = true
  try {
    await apiPost('/api/v1/admin/third-party/import', { teachers, students })
    close()
    emit('imported')
  } catch (e: any) {
    importError.value = e?.response?.data?.message || '导入失败，请稍后重试'
  } finally {
    importing.value = false
  }
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:760px;width:100%;padding:4px 0;">
      <div class="modal-header">
        <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">🏢 从第三方平台导入</h3>
        <button @click="close" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
      </div>

      <div class="modal-body">
        <div v-if="!rows.length" style="text-align:center;padding:32px 16px;">
          <div style="font-size:40px;margin-bottom:12px;">📇</div>
          <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">从学校配置的第三方平台（企业微信 / 钉钉 / 飞书）通讯录拉取成员</p>
          <button class="btn btn-primary" :disabled="loading" @click="loadContacts">
            {{ loading ? '拉取中...' : '📥 拉取通讯录' }}
          </button>
          <div v-if="loadError" style="margin-top:10px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color:#fca5a5;font-size:12px;">{{ loadError }}</div>
        </div>

        <div v-else>
          <div style="display:flex;align-items:center;gap:12px;margin-bottom:12px;">
            <span style="font-size:13px;font-weight:600;color:var(--color-text);">共 {{ rows.length }} 名成员</span>
            <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="toggleAll">全选 / 全不选</button>
          </div>

          <div class="preview-table-wrapper">
            <table class="preview-table">
              <thead>
                <tr>
                  <th></th>
                  <th>姓名</th>
                  <th>手机</th>
                  <th>部门</th>
                  <th>角色</th>
                  <th>目标班级</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in rows" :key="r.userid">
                  <td><input type="checkbox" v-model="r.selected" style="accent-color:#7c3aed;"></td>
                  <td style="font-weight:600;">{{ r.name }}</td>
                  <td>{{ r.mobile }}</td>
                  <td style="color:var(--color-text-secondary);">{{ r.department_names.join(' / ') }}</td>
                  <td>
                    <select v-model="r.role" class="form-input" style="padding:4px 6px;font-size:12px;">
                      <option value="teacher">👨‍🏫 教师</option>
                      <option value="student">👦 学生</option>
                    </select>
                  </td>
                  <td v-if="r.role === 'student'">
                    <select v-model="r.class_id" class="form-input" style="padding:4px 6px;font-size:12px;min-width:130px;">
                      <option :value="''">选择班级</option>
                      <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
                    </select>
                  </td>
                  <td v-else style="color:var(--color-text-secondary);font-size:12px;">—</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p style="font-size:11px;color:var(--color-text-secondary);margin-top:8px;">💡 学生班级已按部门名自动匹配，可手动调整；教师默认以姓名作为登录账号（实名）。</p>
          <div v-if="importError" style="margin-top:8px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color:#fca5a5;font-size:12px;">{{ importError }}</div>
        </div>
      </div>

      <div class="modal-footer" style="justify-content:flex-end;">
        <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);" @click="close">取消</button>
        <button v-if="rows.length" class="btn btn-primary" :disabled="importing" @click="doImport">
          {{ importing ? '导入中...' : '✅ 确认导入' }}
        </button>
      </div>
    </div>
  </ModalGlass>
</template>

<style scoped>
.modal-body { margin-bottom: 16px; }
.modal-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; padding-bottom:12px; border-bottom:1px solid var(--color-border); flex-shrink:0; }
.modal-footer { display:flex; gap:8px; padding-top:12px; border-top:1px solid var(--color-border); margin-top:16px; }
.preview-table-wrapper { max-height: 380px; overflow: auto; }
.preview-table { width: 100%; border-collapse: collapse; font-size: 12px; }
.preview-table th { background: rgba(255,255,255,0.03); padding: 6px 8px; text-align: left; border: 1px solid rgba(255,255,255,0.06); font-size: 11px; color: var(--color-text-secondary); white-space: nowrap; }
.preview-table td { padding: 5px 8px; border: 1px solid rgba(255,255,255,0.06); white-space: nowrap; color: var(--color-text); }
.btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 500; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; font-family: inherit; }
.btn-primary { background: #7c3aed; color: #fff; border-color: #7c3aed; }
.btn-primary:hover { background: #6d28d9; }
.btn-sm { padding: 5px 10px; font-size: 12px; }
</style>
