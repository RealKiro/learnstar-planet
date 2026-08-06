<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import TeacherFilters from './TeacherFilters.vue'
import TeacherCard from './TeacherCard.vue'
import CreateTeacherModal from './CreateTeacherModal.vue'
import EditTeacherModal from './EditTeacherModal.vue'
import AssignClassModal from './AssignClassModal.vue'
import DeleteConfirmModal from './DeleteConfirmModal.vue'
import ResetPasswordModal from './ResetPasswordModal.vue'
import ImportTeacherModal from './ImportTeacherModal.vue'
import ImportWechatModal from './ImportWechatModal.vue'

interface Teacher { id: number; name: string; username: string; nickname?: string; subject?: string; grade_team?: string; phone?: string; email?: string; avatar_path?: string; status: string; bindings: string[]; assignments: Assignment[]; class_names: string[]; personal_role?: string }
interface Assignment { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }
interface ClassRoom { id: number; name: string; grade?: string }
type ClassRole = 'head_teacher' | 'co_teacher' | 'subject_teacher'

const classRoleLabel: Record<ClassRole, string> = { head_teacher: '主班主任', co_teacher: '副班主任', subject_teacher: '科任教师' }
const grades = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']
const subjects = ['语文', '数学', '英语', '科学', '道德与法治', '体育', '音乐', '美术', '信息科技', '综合实践']

const teachers = ref<Teacher[]>([])
const classes = ref<ClassRoom[]>([])
const loading = ref(true)
const filterGrade = ref('')
const filterRole = ref<ClassRole | ''>('')
const searchQuery = ref('')

const teacherTeams = computed(() => {
  const teams: Record<string, Teacher[]> = {}
  if (!teachers.value) return teams
  teachers.value.forEach(t => {
    const team = t.grade_team || '未分组'
    if (!teams[team]) teams[team] = []
    teams[team].push(t)
  })
  return teams
})
// 年级 → 数字（义务教育 1-9 年级）
const GRADE_NUM: Record<string, string> = {
  '一年级': '1', '二年级': '2', '三年级': '3', '四年级': '4', '五年级': '5',
  '六年级': '6', '七年级': '7', '八年级': '8', '九年级': '9',
}
// 带组内编号的教师分组（卡片右上角显示 G{年级}{组内序号}，如 G601 = 六年级第1人）
type TeacherWithNo = Teacher & { no: string }
const teacherList = computed(() => {
  const list: { team: string; teachers: TeacherWithNo[] }[] = []
  for (const [teamName, team] of Object.entries(teacherTeams.value)) {
    let grade = ''
    for (const [cn, num] of Object.entries(GRADE_NUM)) {
      if (teamName.includes(cn)) { grade = num; break }
    }
    list.push({
      team: teamName,
      teachers: team.map((t, i) => ({
        ...t,
        no: grade !== '' ? `G${grade}${String(i + 1).padStart(2, '0')}` : `#${i + 1}`,
      })),
    })
  }
  return list
})
const filteredTeachers = computed(() => {
  if (!teachers.value) return []
  let list = teachers.value
  if (searchQuery.value) { const q = searchQuery.value.toLowerCase(); list = list.filter(t => t.name.includes(q) || t.username.includes(q) || (t.nickname && t.nickname.includes(q))) }
  if (filterGrade.value) list = list.filter(t => t.grade_team === filterGrade.value)
  if (filterRole.value) list = list.filter(t => t.assignments.some((a: Assignment) => a.role === filterRole.value))
  return list
})

const showCreateModal = ref(false)
const showEditModal = ref(false); const editTarget = ref<Teacher | null>(null)
const showAssignModal = ref(false); const assignTarget = ref<Teacher | null>(null)
const showDeleteModal = ref(false); const deleteTarget = ref<Teacher | null>(null)
const showResetPwdModal = ref(false); const resetTarget = ref<Teacher | null>(null)
const showImportModal = ref(false)
const showWechatImport = ref(false)

// 批量选择与操作
const selectedTeacherIds = ref<number[]>([])
const allTeachersSelected = computed(() =>
  teachers.value.length > 0 && selectedTeacherIds.value.length === teachers.value.length
)
function toggleSelectTeacher(id: number) {
  if (selectedTeacherIds.value.includes(id)) selectedTeacherIds.value = selectedTeacherIds.value.filter(x => x !== id)
  else selectedTeacherIds.value.push(id)
}
function toggleSelectAllTeachers() {
  if (allTeachersSelected.value) selectedTeacherIds.value = []
  else selectedTeacherIds.value = teachers.value.map(t => t.id)
}
const batchOpStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const batchOpMsg = ref('')
async function batchResetTeacherPwd() {
  if (!selectedTeacherIds.value.length) { batchOpMsg.value = '请先勾选要重置密码的教师'; return }
  batchOpStatus.value = 'loading'; batchOpMsg.value = ''
  try {
    const res = await apiPost<{ message: string }>('/api/v1/admin/accounts/batch-reset-password', {
      role: 'teacher', ids: selectedTeacherIds.value,
    }, { skipToast: true })
    batchOpStatus.value = 'success'
    batchOpMsg.value = res.message || '已重置密码（默认 ls123456）'
    setTimeout(() => { batchOpStatus.value = 'idle'; batchOpMsg.value = ''; selectedTeacherIds.value = [] }, 2500)
  } catch (e: any) {
    batchOpStatus.value = 'error'
    batchOpMsg.value = e?.response?.data?.message || '重置失败，请稍后重试'
    setTimeout(() => { batchOpStatus.value = 'idle'; batchOpMsg.value = '' }, 3000)
  }
}
async function batchDeleteTeachers() {
  if (!selectedTeacherIds.value.length) { batchOpMsg.value = '请先勾选要删除的教师'; return }
  if (!confirm(`确定删除选中的 ${selectedTeacherIds.value.length} 个教师账号？\n将解除其所有班级分配。`)) return
  batchOpStatus.value = 'loading'; batchOpMsg.value = ''
  try {
    const res = await apiPost<{ message: string }>('/api/v1/admin/accounts/batch-delete', {
      role: 'teacher', ids: selectedTeacherIds.value,
    }, { skipToast: true })
    batchOpStatus.value = 'success'
    batchOpMsg.value = res.message || '已删除'
    selectedTeacherIds.value = []
    setTimeout(() => { batchOpStatus.value = 'idle'; batchOpMsg.value = '' }, 2500)
    refreshTeachers()
  } catch (e: any) {
    batchOpStatus.value = 'error'
    batchOpMsg.value = e?.response?.data?.message || '删除失败，请稍后重试'
    setTimeout(() => { batchOpStatus.value = 'idle'; batchOpMsg.value = '' }, 3000)
  }
}

async function loadTeachers(isInitial = false) {
  if (isInitial) loading.value = true
  try {
    const [tRes, cRes] = await Promise.all([
      apiGet<{ data: Teacher[] }>('/api/v1/admin/teachers'),
      apiGet<{ data: ClassRoom[] }>('/api/v1/admin/classes'),
    ])
    teachers.value = tRes.data || []; classes.value = cRes.data || []
  } catch { /* silent */ }
  finally { if (isInitial) loading.value = false }
}
const refreshTeachers = () => loadTeachers(false)

function downloadTemplate() {
  window.open('/api/v1/admin/teachers/template-csv', '_blank')
}

onMounted(() => loadTeachers(true))
</script>
<template>
  <div class="teachers-admin" style="max-width:1400px;margin:0 auto;padding:0 4px;">
    <TeacherFilters :grades="grades" :classRoleLabel="classRoleLabel" :filterGrade="filterGrade" :filterRole="filterRole" :searchQuery="searchQuery"
      @update:filterGrade="filterGrade = $event" @update:filterRole="filterRole = $event as '' | ClassRole" @update:searchQuery="searchQuery = $event"
      @downloadTemplate="downloadTemplate" @openImport="showImportModal = true" @openWechatImport="showWechatImport = true" @openCreate="showCreateModal = true">
      <span class="count-badge" slot="teacherCount">{{ teachers.length }} 人</span>
    </TeacherFilters>

    <!-- 批量操作栏 -->
    <div v-if="teachers.length > 0" style="display:flex;align-items:center;gap:12px;margin-bottom:12px;padding:8px 12px;background:var(--color-bg);border:1px solid var(--color-border);border-radius:10px;">
      <label style="display:flex;align-items:center;gap:6px;font-size:13px;cursor:pointer;user-select:none;">
        <input type="checkbox" :checked="allTeachersSelected" @change="toggleSelectAllTeachers" style="accent-color:#7c3aed;width:15px;height:15px;">
        全选
      </label>
      <span v-if="selectedTeacherIds.length > 0" style="font-size:13px;color:var(--color-text-secondary);">已选 {{ selectedTeacherIds.length }} 个</span>
      <div style="flex:1;"></div>
      <button class="btn btn-sm" :disabled="!selectedTeacherIds.length || batchOpStatus === 'loading'" :style="{ background: batchOpStatus === 'loading' ? '#f59e0b' : batchOpStatus === 'success' ? '#10b981' : batchOpStatus === 'error' ? '#ef4444' : '', color: batchOpStatus !== 'idle' ? '#fff' : 'var(--color-text)', border: batchOpStatus !== 'idle' ? '1px solid transparent' : '1px solid var(--color-border)' }" @click="batchResetTeacherPwd">
        <template v-if="batchOpStatus === 'loading'">处理中...</template>
        <template v-else-if="batchOpStatus === 'success'">已重置 ✓</template>
        <template v-else>🔑 批量重置密码</template>
      </button>
      <button class="btn btn-sm" :disabled="!selectedTeacherIds.length || batchOpStatus === 'loading'" :style="{ background: '#fee2e2', color: '#dc2626', border: '1px solid #fecaca', opacity: !selectedTeacherIds.length ? 0.5 : 1 }" @click="batchDeleteTeachers">
        批量删除
      </button>
    </div>
    <div v-if="batchOpMsg" style="margin:-4px 0 12px;padding:8px 12px;border-radius:8px;font-size:12px;"
      :style="{ background: batchOpStatus === 'error' ? 'rgba(239,68,68,0.08)' : 'rgba(16,185,129,0.08)', border: '1px solid ' + (batchOpStatus === 'error' ? 'rgba(239,68,68,0.2)' : 'rgba(16,185,129,0.2)'), color: batchOpStatus === 'error' ? '#fca5a5' : '#6ee7b7' }">{{ batchOpMsg }}</div>

    <div v-if="loading" class="loading-spinner">加载中...</div>
    <div v-else-if="filteredTeachers.length === 0" class="empty-state">
      <div class="empty-icon">&#x1F468;&#x200D;&#x1F3EB;</div>
      <p v-if="searchQuery || filterGrade || filterRole">未找到匹配的教师，请修改搜索条件</p>
      <p v-else>暂无教师，点击「创建教师」添加</p>
    </div>

    <template v-for="group in teacherList" :key="group.team">
      <div v-if="group.teachers.length > 0" style="margin-bottom:20px;">
        <div class="grade-header">
          <span class="grade-dot"></span>
          <span class="grade-name">{{ group.team }}</span>
          <span class="grade-count">{{ group.teachers.length }} 人</span>
        </div>
        <div class="card-grid">
          <TeacherCard v-for="t in group.teachers" :key="t.id" :teacher="t" :no="t.no" :classes="classes"
            selectable :selected="selectedTeacherIds.includes(t.id)" @toggleSelect="toggleSelectTeacher"
            @edit="editTarget = $event; showEditModal = true"
            @assign="assignTarget = $event; showAssignModal = true"
            @resetPwd="resetTarget = $event; showResetPwdModal = true"
            @delete="deleteTarget = $event; showDeleteModal = true" />
        </div>
      </div>
    </template>
  </div>

  <CreateTeacherModal :visible="showCreateModal" :classes="classes" :grades="grades" :subjects="subjects" :classRoleLabel="classRoleLabel"
    @update:visible="showCreateModal = $event" @created="refreshTeachers" />
  <EditTeacherModal :visible="showEditModal" :teacher="editTarget" :grades="grades" :subjects="subjects"
    @update:visible="showEditModal = $event; editTarget = null" @updated="refreshTeachers" />
  <AssignClassModal :visible="showAssignModal" :teacher="assignTarget" :classes="classes" :grades="grades" :subjects="subjects" :classRoleLabel="classRoleLabel"
    @update:visible="showAssignModal = $event; assignTarget = null" @assigned="refreshTeachers" />
  <DeleteConfirmModal :visible="showDeleteModal" :teacher="deleteTarget"
    @update:visible="showDeleteModal = $event; deleteTarget = null" @deleted="refreshTeachers" />
  <ResetPasswordModal :visible="showResetPwdModal" :teacher="resetTarget"
    @update:visible="showResetPwdModal = $event; resetTarget = null" @reset="refreshTeachers" />
  <ImportTeacherModal :visible="showImportModal" @update:visible="showImportModal = $event" @imported="refreshTeachers" />
  <ImportWechatModal :visible="showWechatImport" @update:visible="showWechatImport = $event" @imported="refreshTeachers" />
</template>

<style scoped>
.section-badge { font-size:11px;font-weight:600;color:#7c3aed;text-transform:uppercase;letter-spacing:0.05em;background:#ede9fe;padding:3px 10px;border-radius:6px; }
.page-title { font-size:22px;font-weight:700;margin:0;line-height:1.2; }
.count-badge { font-size:13px;color:#6b7280;background:var(--color-bg);padding:2px 10px;border-radius:10px; }
.grade-header { display:flex; align-items:center; gap:10px; margin:20px 0 12px; padding:0 0 8px 4px; border-bottom:1px solid rgba(255,255,255,0.06); }
.grade-dot { width:8px; height:8px; border-radius:50%; background:#8b5cf6; flex-shrink:0; }
.grade-name { font-size:16px; font-weight:700; }
.grade-count { font-size:12px; color:var(--color-text-secondary); background:rgba(255,255,255,0.06); padding:0 10px; border-radius:20px; }
.card-grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(240px, 1fr)); gap:16px; }
.loading-spinner { text-align:center;padding:64px;color:#9ca3af;font-size:15px; }
.empty-state { text-align:center;padding:64px 20px;color:#9ca3af; }
.empty-icon { font-size:48px;margin-bottom:12px; }
</style>
