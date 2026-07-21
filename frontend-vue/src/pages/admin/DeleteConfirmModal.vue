<script setup lang="ts">
import { ref, watch } from 'vue'
import { apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import ModalGlass from '@/components/common/ModalGlass.vue'

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

const props = defineProps<{
  visible: boolean
  teacher: Teacher | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'deleted'): void
}>()

const toast = useToastStore()

const deleteConfirmed = ref(false)
const deleteConfirmName = ref('')

watch(
  () => props.teacher,
  () => {
    deleteConfirmed.value = false
    deleteConfirmName.value = ''
  },
  { immediate: true }
)

function close() {
  emit('update:visible', false)
}

async function confirmDelete() {
  if (!props.teacher) return
  try {
    await apiDelete(`/api/v1/admin/teachers/${props.teacher.id}`)
    toast.show('教师已删除', 'success', { position: 'center', duration: 2000 })
    emit('update:visible', false)
    emit('deleted')
  } catch {
    toast.show('删除失败', 'error', { position: 'center', duration: 2000 })
  }
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:400px;width:100%;padding:4px 0;">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
        <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">&#9888;&#65039; 确认删除</h3>
        <button @click="close" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">&#10005;</button>
      </div>
      <div style="text-align:center;padding:8px 0;">
        <div style="font-size:40px;margin-bottom:8px;">&#128465;&#65039;</div>
        <p style="font-size:15px;font-weight:600;color:var(--color-text);margin-bottom:8px;">
          确定要删除教师「{{ teacher?.name }}」吗？
        </p>
        <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:12px;">
          此操作将永久删除该教师账号及其所有班级任教分配，且不可恢复。
        </p>
        <div
          v-if="teacher && teacher.assignments && teacher.assignments.some(a => a.role === 'head_teacher')"
          style="font-size:12px;color:#dc2626;padding:8px;background:rgba(239,68,68,0.06);border-radius:6px;margin-bottom:12px;"
        >
          &#9888;&#65039; 该教师是部分班级的班主任，删除后这些班级将无班主任。
        </div>
        <label
          style="display:flex;align-items:center;gap:6px;font-size:12px;color:var(--color-text-secondary);cursor:pointer;justify-content:center;padding:8px;background:var(--color-bg);border-radius:6px;"
        >
          <input type="checkbox" v-model="deleteConfirmed" style="accent-color:#dc2626;">
          确认我已了解此操作不可恢复
        </label>
        <div style="margin-top:10px;">
          <label style="display:block;font-size:11px;color:var(--color-text-secondary);margin-bottom:4px;text-align:left;">
            请输入教师姓名「{{ teacher?.name }}」以确认删除
          </label>
          <input
            v-model="deleteConfirmName"
            class="form-input"
            :placeholder="'请输入 ' + (teacher?.name || '')"
            style="font-size:12px;text-align:center;"
          >
        </div>
      </div>
      <div style="display:flex;gap:8px;justify-content:flex-end;padding-top:12px;border-top:1px solid var(--color-border);">
        <button
          @click="close"
          style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);"
        >
          取消
        </button>
        <button
          @click="confirmDelete"
          :disabled="!deleteConfirmed || deleteConfirmName !== teacher?.name"
          style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#dc2626;border:none;color:#fff;opacity:deleteConfirmed && deleteConfirmName === teacher?.name ? 1 : 0.5;"
        >
          确认删除
        </button>
      </div>
    </div>
  </ModalGlass>
</template>

<style scoped>
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
</style>
