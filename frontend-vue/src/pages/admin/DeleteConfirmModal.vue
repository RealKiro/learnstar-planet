<script setup lang="ts">
import { ref, watch } from 'vue'
import { apiDelete } from '@/utils/api'
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

const deleteError = ref('')
const deleteStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')

watch(
  () => props.teacher,
  () => {
    deleteError.value = ''
    deleteStatus.value = 'idle'
  },
  { immediate: true }
)

function closeModal() {
  emit('update:visible', false)
}

async function confirmDelete() {
  if (!props.teacher) return
  deleteStatus.value = 'loading'
  deleteError.value = ''
  try {
    await apiDelete(`/api/v1/admin/teachers/${props.teacher.id}`, { skipToast: true })
    deleteStatus.value = 'success'
    setTimeout(() => emit('update:visible', false), 600)
    emit('deleted')
  } catch (e: any) {
    deleteStatus.value = 'error'
    deleteError.value = e?.response?.data?.message || '删除失败，请稍后重试'
    setTimeout(() => { if (deleteStatus.value === 'error') deleteStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div class="del-wrap">
      <div class="modal-header del-mb-12">
        <h3 class="card-title">&#9888;&#65039; 确认删除</h3>
        <button @click="closeModal" class="icon-close">&#10005;</button>
      </div>
      <div class="del-center">
        <div class="del-emoji">&#128465;&#65039;</div>
        <p class="del-title">
          确定要删除教师「{{ teacher?.name }}」吗？
        </p>
        <p class="del-desc">
          此操作将永久删除该教师账号及其所有班级任教分配，且不可恢复。
        </p>
        <div
          v-if="teacher && teacher.assignments && teacher.assignments.some(a => a.role === 'head_teacher')" class="del-warn"
        >
          &#9888;&#65039; 该教师是部分班级的班主任，删除后这些班级将无班主任。
        </div>
      </div>
      <div v-if="deleteError" class="del-error">{{ deleteError }}</div>
      <div class="modal-footer del-footer-end">
        <button
          @click="closeModal" class="del-btn-outline"
        >
          取消
        </button>
        <button
          @click="confirmDelete"
          :disabled="deleteStatus === 'loading'"
          :style="{
            padding: '8px 20px', borderRadius: '8px', fontSize: '13px', fontWeight: '600', cursor: 'pointer',
            background: deleteStatus === 'loading' ? '#f59e0b' : deleteStatus === 'success' ? '#10b981' : deleteStatus === 'error' ? '#ef4444' : '#dc2626',
            border: 'none', color: '#fff'
          }"
        >
          {{ deleteStatus === 'loading' ? '删除中...' : deleteStatus === 'success' ? '已删除 ✓' : deleteStatus === 'error' ? '删除失败 ✗' : '确认删除' }}
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
/* Modal utility classes */
.modal-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; padding-bottom:12px; border-bottom:1px solid var(--color-border); flex-shrink:0; }
.modal-footer { display:flex; gap:8px; padding-top:12px; border-top:1px solid var(--color-border); margin-top:16px; }
.modal-section-title { font-size:12px; font-weight:600; color:var(--color-text); margin-bottom:8px; }
.flex-row { display:flex; gap:8px; }
.flex-1 { flex:1; }
/* ===== P1 内联样式收口（声明逐字保留以保渲染等价） ===== */
.del-wrap { max-width:400px;width:100%;padding:4px 0; }
.del-mb-12 { margin-bottom:12px; }
.del-center { text-align:center;padding:8px 0; }
.del-emoji { font-size:40px;margin-bottom:8px; }
.del-title { font-size:15px;font-weight:600;color:var(--color-text);margin-bottom:8px; }
.del-desc { font-size:12px;color:var(--color-text-secondary);margin-bottom:12px; }
.del-warn { font-size:12px;color:#dc2626;padding:8px;background:rgba(239,68,68,0.06);border-radius:6px;margin-bottom:12px; }
.del-error { margin-bottom:10px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color: var(--color-danger-text);font-size:12px; }
.del-footer-end { justify-content:flex-end; }
.del-btn-outline { padding:8px 20px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text); }
</style>
