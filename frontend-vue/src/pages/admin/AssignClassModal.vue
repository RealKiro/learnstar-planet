<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { apiPut } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import ModalGlass from '@/components/common/ModalGlass.vue'

type ClassRole = 'head_teacher' | 'co_teacher' | 'subject_teacher'

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

interface ClassRoom {
  id: number
  name: string
  grade?: string
}

const props = defineProps<{
  visible: boolean
  teacher: Teacher | null
  classes: ClassRoom[]
  grades: string[]
  subjects: string[]
  classRoleLabel: Record<string, string>
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'assigned'): void
}>()

const toast = useToastStore()

const assignList = ref<{ class_id: number | null; role: ClassRole; subject: string; class_name?: string }[]>([])
const assignLoading = ref(false)
const assignGradeFilter = ref('')
const newAssignClassId = ref<number | null>(null)
const newAssignRole = ref<ClassRole>('subject_teacher')
const newAssignSubject = ref('')

const filteredAssignClasses = computed(() => {
  if (!props.classes) return []
  if (!assignGradeFilter.value) return props.classes
  return props.classes.filter(c => c.grade === assignGradeFilter.value)
})

function classById(id: number) {
  return props.classes.find(c => c.id === id)
}

function shortClassName(name: string | undefined) {
  if (!name) return ''
  const m = name.match(/（(\d+)）班/)
  return m ? m[1] + '班' : name
}

watch(
  () => props.teacher,
  (t) => {
    if (t) {
      assignList.value = t.assignments.length > 0
        ? t.assignments.map(a => ({
            class_id: a.class_id,
            role: a.role as ClassRole,
            subject: a.subject || '',
            class_name: a.class_name || classById(a.class_id)?.name,
          }))
        : []
      assignGradeFilter.value = ''
      newAssignClassId.value = null
      newAssignRole.value = 'subject_teacher'
      newAssignSubject.value = ''
    }
  },
  { immediate: true }
)

function addAssignRowNew() {
  if (!newAssignClassId.value) return
  const cls = props.classes.find(c => c.id === newAssignClassId.value)
  if (!cls) return
  const role = newAssignRole.value

  if (role === 'head_teacher' && assignList.value.some(a => a.class_id === cls.id && a.role === 'co_teacher')) {
    toast.show('该班级已有副班，不能同时为主班', 'error', { position: 'center', duration: 2000 })
    return
  }
  if (role === 'co_teacher' && assignList.value.some(a => a.class_id === cls.id && a.role === 'head_teacher')) {
    toast.show('该班级已有主班，不能同时为副班', 'error', { position: 'center', duration: 2000 })
    return
  }
  if (assignList.value.some(a => a.class_id === cls.id && a.role === role)) {
    toast.show('该班级已分配此角色', 'info', { position: 'center', duration: 2000 })
    return
  }

  const subjectValue = newAssignSubject.value
  if (!subjectValue) {
    toast.show('请选择科目', 'error', { position: 'center', duration: 2000 })
    return
  }

  assignList.value.push({ class_id: cls.id, class_name: shortClassName(cls.name), role, subject: subjectValue })
  newAssignClassId.value = null
  newAssignSubject.value = ''
}

function removeAssignRow(idx: number) {
  assignList.value.splice(idx, 1)
}

async function submitAssign() {
  if (!props.teacher) return
  assignLoading.value = true
  try {
    const payload = {
      assignments: assignList.value.filter(a => a.class_id).map(a => ({
        class_id: a.class_id,
        role: a.role,
        subject: a.subject === '默认科目' ? undefined : a.subject || undefined,
      })),
    }
    await apiPut(`/api/v1/admin/teachers/${props.teacher.id}/classes`, payload)
    toast.show('班级分配已更新', 'success', { position: 'center', duration: 2000 })
    emit('update:visible', false)
    emit('assigned')
  } catch (e: any) {
    toast.show(e?.response?.data?.message || '保存失败', 'error', { position: 'center', duration: 2000 })
  } finally {
    assignLoading.value = false
  }
}

function closeModal() {
  emit('update:visible', false)
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:620px;width:100%;">
      <div class="modal-header">
        <h2 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">
          &#127979; 分配班级 — {{ teacher?.name }}
        </h2>
        <button @click="closeModal" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">&#10005;</button>
      </div>

      <!-- 选择器 -->
      <div class="flex-row" style="align-items:flex-end;margin-bottom:12px;flex-wrap:wrap;">
        <div style="flex:1;min-width:100px;" class="form-group">
          <label>年级</label>
          <select v-model="assignGradeFilter" class="form-input">
            <option value="">全部</option>
            <option v-for="g in grades" :key="g" :value="g">{{ g }}</option>
          </select>
        </div>
        <div style="flex:1;min-width:120px;" class="form-group">
          <label>班级</label>
          <select v-model="newAssignClassId" class="form-input">
            <option :value="null">请选择</option>
            <option v-for="c in filteredAssignClasses" :key="c.id" :value="c.id">{{ shortClassName(c.name) }}</option>
          </select>
        </div>
        <div style="flex:1;" class="form-group">
          <label>角色</label>
          <select v-model="newAssignRole" class="form-input" style="font-size:12px;padding:5px 8px;">
            <option value="head_teacher">主班主任</option>
            <option value="co_teacher">副班主任</option>
            <option value="subject_teacher">科任教师</option>
          </select>
        </div>
        <div style="flex:1;min-width:90px;" class="form-group">
          <label>科目</label>
          <select v-model="newAssignSubject" class="form-input">
            <option value="">请选择科目</option>
            <option v-for="s in subjects" :key="s" :value="s">{{ s }}</option>
          </select>
        </div>
        <button
          @click="addAssignRowNew"
          :disabled="!newAssignClassId"
          style="padding:8px 16px;border-radius:8px;border:1px solid var(--color-accent);background:rgba(79,70,229,0.08);color:var(--color-accent);font-size:13px;cursor:pointer;font-weight:500;white-space:nowrap;height:36px;flex-shrink:0;align-self:flex-end;"
        >
          &#10133; 添加
        </button>
      </div>

      <!-- 已分配列表 -->
      <div style="margin-bottom:12px;">
        <div class="modal-section-title" style="font-size:13px;">
          &#128203; 已分配（{{ assignList.length }}）
        </div>
        <div
          v-if="assignList.length === 0"
          style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);background:var(--color-bg);border-radius:8px;"
        >
          暂未分配班级
        </div>
        <div
          v-for="(a, i) in assignList"
          :key="i"
          style="display:flex;align-items:center;gap:8px;padding:8px 12px;margin-bottom:4px;background:var(--color-bg);border-radius:8px;font-size:13px;"
        >
          <span style="flex:1;font-weight:500;color:var(--color-text);">{{ a.class_name || shortClassName(classById(a.class_id)?.name) }}</span>
          <span v-if="a.role === 'subject_teacher'" style="color:var(--color-text-secondary);font-size:12px;">{{ a.subject }}</span>
          <span v-else style="color:var(--color-text-secondary);font-size:12px;">
            {{ classRoleLabel[a.role as ClassRole] || a.role }} · {{ a.subject || '默认科目' }}
          </span>
          <button
            @click="removeAssignRow(i)"
            style="background:none;border:none;color:var(--color-danger);cursor:pointer;padding:2px;font-size:16px;"
          >
            &#10005;
          </button>
        </div>
      </div>

      <div class="modal-footer" style="justify-content:flex-end;">
        <button
          @click="closeModal"
          style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);"
        >
          取消
        </button>
        <button
          @click="submitAssign"
          :disabled="assignLoading"
          style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;"
        >
          {{ assignLoading ? '保存中...' : '保存分配' }}
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
