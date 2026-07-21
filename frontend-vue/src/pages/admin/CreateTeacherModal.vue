<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import ModalGlass from '@/components/common/ModalGlass.vue'

interface ClassRoom {
  id: number
  name: string
  grade?: string
}

type ClassRole = 'head_teacher' | 'co_teacher' | 'subject_teacher'

interface CreateAssignment {
  class_id: number
  class_name: string
  subject: string
  role: ClassRole
}

type CreateStatus = 'idle' | 'loading' | 'success' | 'error'

const props = withDefaults(defineProps<{
  visible: boolean
  classes: ClassRoom[]
  grades: string[]
  subjects: string[]
  classRoleLabel: Record<string, string>
}>(), {
  classes: () => [],
  grades: () => [],
  subjects: () => [],
  classRoleLabel: () => ({}),
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'created': []
}>()

const toast = useToastStore()

const createForm = ref({ name: '', nickname: '', grade_team: '', phone: '', email: '', password: '' })
const createErrors = ref<Record<string, string>>({})
const createStatus = ref<CreateStatus>('idle')
const createErrorMsg = ref('')
const createAssignments = ref<CreateAssignment[]>([])
const createLoading = ref(false)
const pendingGrade = ref('')
const pendingClassId = ref<number | null>(null)
const pendingSubject = ref('')
const pendingRole = ref<ClassRole | ''>('')
const assignError = ref('')

const gradeClasses = computed(() => (props.classes || []).filter(c => c.grade === pendingGrade.value))

const groupedAssignments = computed(() => {
  const map: Record<string, { class_name: string; items: { index: number; role: string; subject: string }[] }> = {}
  createAssignments.value.forEach((a, idx) => {
    const key = String(a.class_id)
    if (!map[key]) map[key] = { class_name: a.class_name, items: [] }
    map[key].items.push({ index: idx, role: a.role, subject: a.subject })
  })
  return Object.values(map)
})

function shortClassName(name: string | undefined) {
  if (!name) return ''
  const m = name.match(/（(\d+)）班/)
  return m ? m[1] + '班' : name
}

function classById(id: number) {
  return props.classes.find(c => c.id === id)
}

function clearError(field: string) {
  delete createErrors.value[field]
}

function validateField(field: string, value: string) {
  if (field === 'name' && !value.trim()) {
    createErrors.value.name = '请填写教师姓名'
    return false
  }
  if (field === 'password' && value && value.length < 6) {
    createErrors.value.password = '密码长度至少 6 位，建议包含字母和数字'
    return false
  }
  delete createErrors.value[field]
  return true
}

function resetForm() {
  createForm.value = { name: '', nickname: '', grade_team: '', phone: '', email: '', password: '' }
  createErrors.value = {}
  assignError.value = ''
  createAssignments.value = []
  pendingGrade.value = ''
  pendingClassId.value = null
  pendingSubject.value = ''
  pendingRole.value = ''
  createStatus.value = 'idle'
  createErrorMsg.value = ''
  createLoading.value = false
}

function closeCreateModal() {
  emit('update:visible', false)
  createStatus.value = 'idle'
  createErrorMsg.value = ''
}

watch(() => props.visible, (val) => {
  if (val) resetForm()
})

function addClassAssignment() {
  assignError.value = ''
  if (!pendingClassId.value || !pendingGrade.value) {
    assignError.value = '请先选择年级和班级'
    return
  }
  const cls = props.classes.find(c => c.id === pendingClassId.value)
  if (!cls) return
  if (createAssignments.value.some(a => a.class_id === cls.id && a.role === (pendingRole.value || 'subject_teacher') && a.subject === pendingSubject.value)) {
    assignError.value = '该班级已添加相同角色和科目'
    return
  }
  if (!pendingSubject.value) {
    assignError.value = '请选择科目'
    return
  }
  createAssignments.value.push({
    class_id: cls.id,
    class_name: cls.name,
    subject: pendingSubject.value,
    role: pendingRole.value || 'subject_teacher',
  })
  pendingClassId.value = null
  pendingSubject.value = ''
  pendingRole.value = ''
}

function removeClassAssignment(idx: number) {
  createAssignments.value.splice(idx, 1)
}

async function submitCreate() {
  createErrors.value = {}
  if (!createForm.value.name.trim()) {
    createErrors.value.name = '请填写教师姓名'
    return
  }
  if (createForm.value.password && createForm.value.password.length < 6) {
    createErrors.value.password = '密码长度至少 6 位，建议包含字母和数字'
    return
  }
  if (Object.keys(createErrors.value).length > 0) return

  createStatus.value = 'loading'
  createLoading.value = true
  try {
    const payload: Record<string, any> = { name: createForm.value.name.trim() }
    if (createForm.value.nickname) payload.nickname = createForm.value.nickname
    if (createForm.value.grade_team) payload.grade_team = createForm.value.grade_team
    if (createForm.value.phone) payload.phone = createForm.value.phone
    if (createForm.value.email) payload.email = createForm.value.email
    if (createForm.value.password) payload.password = createForm.value.password
    if (createAssignments.value.length > 0 && createAssignments.value.every(a => a.class_id)) {
      payload.assignments = createAssignments.value.map(a => ({
        class_id: a.class_id,
        role: a.role || 'subject_teacher',
        subject: a.subject || undefined,
      }))
    }
    await apiPost('/api/v1/admin/teachers', payload)
    createStatus.value = 'success'
    setTimeout(() => {
      emit('update:visible', false)
      createStatus.value = 'idle'
      emit('created')
    }, 1500)
  } catch (e: any) {
    createStatus.value = 'error'
    console.error('创建教师失败:', e?.response?.status, e?.response?.data)
    console.error('发送的payload:', JSON.stringify(e?.config?.data || '(无)'))
    const errs = e?.response?.data?.errors
    if (errs) {
      for (const [field, msgs] of Object.entries(errs)) {
        createErrors.value[field] = (msgs as string[])[0]
      }
      createStatus.value = 'idle'
    } else {
      createErrorMsg.value = e?.response?.data?.message || '创建失败，请重试'
      setTimeout(() => {
        if (createStatus.value === 'error') {
          createStatus.value = 'idle'
          createErrorMsg.value = ''
        }
      }, 3000)
    }
  } finally {
    createLoading.value = false
  }
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);flex-shrink:0;">
      <h2 style="font-size:18px;font-weight:700;color:var(--color-text);margin:0;">✨ 创建教师账号</h2>
      <button :disabled="createLoading" @click="closeCreateModal" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:4px;line-height:1;">✕</button>
    </div>
    <div style="overflow-y:auto;">
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;">
        <div>
          <div style="font-size:12px;font-weight:600;color:var(--color-text);margin-bottom:8px;">📝 创建账号 <span style="font-size:10px;color:var(--color-text-secondary);">必填</span></div>
          <div style="display:flex;gap:8px;">
            <div style="flex:1;" class="form-group">
              <label>姓名 <span style="color:var(--color-danger);">*</span></label>
              <input v-model="createForm.name" placeholder="姓名" class="form-input" :style="{ borderColor: createErrors.name ? '#f87171' : '' }" @blur="validateField('name', createForm.name)" @input="clearError('name')" />
              <div v-if="createErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ createErrors.name }}</div>
            </div>
            <div style="flex:1;" class="form-group">
              <label>年级团队</label>
              <select v-model="createForm.grade_team" class="form-input">
                <option value="">不指定</option>
                <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
              </select>
            </div>
          </div>
          <div style="display:flex;gap:8px;">
            <div style="flex:1;" class="form-group"><label>手机号</label><input v-model="createForm.phone" placeholder="选填" class="form-input" /></div>
            <div style="flex:1;" class="form-group"><label>邮箱</label><input v-model="createForm.email" placeholder="选填" class="form-input" /></div>
          </div>
          <div style="display:none;"></div>
          <div class="form-group">
            <label>初始密码</label>
            <input v-model="createForm.password" placeholder="留空自动生成" class="form-input" :style="{ borderColor: createErrors.password ? '#f87171' : '' }" @blur="validateField('password', createForm.password)" @input="clearError('password')" />
            <div v-if="createErrors.password" style="color:#f87171;font-size:11px;margin-top:2px;">{{ createErrors.password }}</div>
          </div>
        </div>
        <div>
          <div style="font-size:12px;font-weight:600;color:var(--color-text);margin-bottom:8px;">📚 加入班级 <span style="font-size:10px;color:var(--color-text-secondary);">可选</span></div>
          <div style="display:grid;grid-template-columns:1fr 1fr;gap:6px;">
            <div class="form-group" style="margin-bottom:0;"><label>年级</label><select v-model="pendingGrade" class="form-input"><option value="">请选择</option><option v-for="g in grades" :key="g" :value="g">{{ g }}</option></select></div>
            <div class="form-group" style="margin-bottom:0;"><label>班级</label><select v-model="pendingClassId" :disabled="!pendingGrade" class="form-input"><option :value="null">请选择</option><option v-for="c in gradeClasses" :key="c.id" :value="c.id">{{ shortClassName(c.name) }}</option></select></div>
            <div class="form-group" style="margin-bottom:0;"><label>角色</label><select v-model="pendingRole" class="form-input"><option value="">选择角色（默认科任教师）</option><option value="head_teacher">主班主任</option><option value="co_teacher">副班主任</option></select></div>
            <div class="form-group" style="margin-bottom:0;"><label>科目</label><select v-model="pendingSubject" class="form-input"><option value="">请选择科目</option><option v-for="s in subjects" :key="s" :value="s">{{ s }}</option></select></div>
          </div>
          <div style="display:flex;justify-content:flex-end;margin-top:6px;">
            <button @click="addClassAssignment" :disabled="!pendingClassId" style="padding:5px 16px;border-radius:8px;border:1px solid var(--color-accent);background:rgba(79,70,229,0.08);color:var(--color-accent);font-size:13px;cursor:pointer;font-weight:500;">➕ 添加</button>
          </div>
          <div v-if="assignError" style="color:#f87171;font-size:11px;margin-top:4px;">{{ assignError }}</div>
          <div style="font-size:11px;color:var(--color-text-secondary);margin-top:8px;padding:6px 8px;background:var(--color-bg);border-radius:6px;border-left:2px solid var(--color-accent);">💡 创建后可随时在列表中点击 🏫 按钮重新分配班级</div>
          <div v-if="createAssignments.length > 0" style="margin-top:8px;display:flex;flex-direction:column;gap:4px;">
            <div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:2px;">📋 已添加（{{ createAssignments.length }}）</div>
            <div v-for="(group, gi) in groupedAssignments" :key="gi" style="display:flex;align-items:center;gap:6px;padding:4px 8px;background:var(--color-bg);border-radius:4px;font-size:12px;border-left:3px solid var(--color-accent);">
              <span style="flex:1;color:var(--color-text);font-weight:500;">{{ shortClassName(group.class_name) }}</span>
              <template v-for="(a, ai) in group.items" :key="ai">
                <span v-if="a.role === 'subject_teacher'" style="font-size:11px;color:var(--color-text-secondary);">{{ a.subject }}</span>
                <span v-else style="font-size:11px;font-weight:500;color:var(--color-accent);">{{ classRoleLabel[a.role] || a.role }} · {{ a.subject }}</span>
                <span v-if="ai < group.items.length - 1" style="color:var(--color-border);font-size:10px;">|</span>
              </template>
              <button @click="removeClassAssignment(group.items[0].index)" style="background:none;border:none;color:var(--color-danger);cursor:pointer;padding:0;font-size:14px;flex-shrink:0;">✕</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div style="display:flex;gap:8px;padding:12px 20px;border-top:1px solid var(--color-border);flex-shrink:0;">
      <button @click="closeCreateModal" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
      <button @click="submitCreate" :disabled="createStatus === 'loading'" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;border:none;color:#fff;transition:all 0.3s ease;box-shadow:0 2px 8px rgba(124,58,237,0.15);"
        :style="{ background: createStatus === 'loading' ? '#f59e0b' : createStatus === 'success' ? '#10b981' : createStatus === 'error' ? '#ef4444' : '#7c3aed' }">
        <span v-if="createStatus === 'idle'">创建账号</span>
        <span v-else-if="createStatus === 'loading'">⏳ 创建中...</span>
        <span v-else-if="createStatus === 'success'">✅ 创建成功</span>
        <span v-else>❌ {{ createErrorMsg || '创建失败' }}</span>
      </button>
    </div>
  </ModalGlass>
</template>

<style scoped>
.form-group {
  margin-bottom: 10px;
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
  box-shadow: 0 0 0 3px rgba(124,58,237,0.08);
}
</style>
