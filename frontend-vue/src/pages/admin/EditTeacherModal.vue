<script setup lang="ts">
import { ref, watch } from 'vue'
import { apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import ModalGlass from '@/components/common/ModalGlass.vue'

type PersonalRole = 'grade_lead' | 'admin_director' | null

interface Teacher {
  id: number
  name: string
  username: string
  nickname?: string
  subject?: string
  grade_team?: string
  phone?: string
  email?: string
  avatar_path?: string
  status: string
  bindings: string[]
  assignments: { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }[]
  class_names: string[]
  personal_role?: string
}

const personalRoleOptions = [
  { v: 'grade_lead', l: '首席', bg: '#8b5cf6' },
  { v: 'admin_director', l: '主任', bg: '#f59e0b' },
  { v: '', l: '普通科任', bg: 'var(--color-primary)' },
]

const props = defineProps<{
  visible: boolean
  teacher: Teacher | null
  grades: string[]
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'updated'): void
}>()

const toast = useToastStore()

const editForm = ref({ name: '', nickname: '', grade_team: '', phone: '', email: '', personalRole: '' as string })

watch(
  () => props.teacher,
  (t) => {
    if (t) {
      const pRole = t.personal_role === 'grade_lead' ? 'grade_lead' : t.personal_role === 'admin_director' ? 'admin_director' : ''
      editForm.value = {
        name: t.name,
        nickname: t.nickname || '',
        grade_team: t.grade_team || '',
        phone: t.phone || '',
        email: t.email || '',
        personalRole: pRole,
      }
    }
  },
)

function closeModal() {
  emit('update:visible', false)
}

async function submitEdit() {
  if (!props.teacher) return
  try {
    const payload: any = {
      name: editForm.value.name,
      nickname: editForm.value.nickname,
      grade_team: editForm.value.grade_team,
      phone: editForm.value.phone,
      email: editForm.value.email,
      personal_role: editForm.value.personalRole || null,
    }
    await apiPut(`/api/v1/admin/teachers/${props.teacher.id}`, payload)
    toast.show('教师信息已更新', 'success', { position: 'center', duration: 2000 })
    emit('update:visible', false)
    emit('updated')
  } catch {
    toast.show('保存失败', 'error', { position: 'center', duration: 2000 })
  }
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:520px;width:100%;">
      <div class="modal-header">
        <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">
          &#9999;&#65039; 编辑教师信息 — {{ teacher?.name }}
        </h3>
        <button @click="closeModal" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">&#10005;</button>
      </div>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;">
        <div class="form-group">
          <label>姓名 <span style="color:#f87171;">*</span></label>
          <input v-model="editForm.name" class="form-input" placeholder="教师姓名">
        </div>
        <div class="form-group">
          <label>昵称</label>
          <input v-model="editForm.nickname" class="form-input" placeholder="默认拼音">
        </div>
        <div class="form-group">
          <label>年级团队</label>
          <select v-model="editForm.grade_team" class="form-input">
            <option value="">不指定</option>
            <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
          </select>
        </div>
        <div class="form-group">
          <label>个人角色</label>
          <div class="flex-row" style="margin-top:4px;">
            <label
              v-for="opt in personalRoleOptions"
              :key="opt.v || '__none'"
              style="display:flex;align-items:center;gap:4px;font-size:12px;cursor:pointer;padding:4px 10px;border-radius:6px;border:1px solid var(--color-border);user-select:none;"
              :style="editForm.personalRole === opt.v ? { background: opt.bg, color: '#fff' } : { background: 'var(--color-bg-card)', color: 'var(--color-text)' }"
            >
              <input
                type="radio"
                :value="opt.v"
                v-model="editForm.personalRole"
                style="opacity:0;position:absolute;width:0;height:0;pointer-events:none;"
              >
              {{ opt.l }}
            </label>
          </div>
        </div>
        <div class="form-group">
          <label>手机号</label>
          <input v-model="editForm.phone" class="form-input" placeholder="11位手机号">
        </div>
        <div class="form-group">
          <label>邮箱</label>
          <input v-model="editForm.email" class="form-input" placeholder="邮箱地址">
        </div>
      </div>
      <div class="modal-footer" style="justify-content:flex-end;margin-top:20px;">
        <button
          @click="closeModal"
          style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);"
        >
          取消
        </button>
        <button
          @click="submitEdit"
          style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;"
        >
          保存
        </button>
      </div>
    </div>
  </ModalGlass>
</template>

<style scoped>
.form-group {
  flex: 1;
}
.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4px;
}
.form-input {
  color: var(--color-text);
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
  box-sizing: border-box;
}
.form-input:focus {
  border-color: #7c3aed;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.08);
}
.form-input option {
  color: #1E293B;
  background: #fff;
}
/* Modal utility classes */
.modal-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; padding-bottom:12px; border-bottom:1px solid var(--color-border); flex-shrink:0; }
.modal-footer { display:flex; gap:8px; padding-top:12px; border-top:1px solid var(--color-border); margin-top:16px; }
.modal-section-title { font-size:12px; font-weight:600; color:var(--color-text); margin-bottom:8px; }
.flex-row { display:flex; gap:8px; }
.flex-1 { flex:1; }
</style>
