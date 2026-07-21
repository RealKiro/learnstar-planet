<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet } from '@/utils/api'
import TeacherFilters from './TeacherFilters.vue'
import TeacherCard from './TeacherCard.vue'
import CreateTeacherModal from './CreateTeacherModal.vue'
import EditTeacherModal from './EditTeacherModal.vue'
import AssignClassModal from './AssignClassModal.vue'
import DeleteConfirmModal from './DeleteConfirmModal.vue'
import ResetPasswordModal from './ResetPasswordModal.vue'
import ImportTeacherModal from './ImportTeacherModal.vue'

interface Teacher { id: number; name: string; username: string; nickname?: string; subject?: string; grade_team?: string; phone?: string; email?: string; avatar_path?: string; status: string; bindings: string[]; assignments: Assignment[]; class_names: string[]; personal_role?: string }
interface Assignment { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }
interface ClassRoom { id: number; name: string; grade?: string }
type ClassRole = 'head_teacher' | 'co_teacher' | 'subject_teacher'
type PersonalRole = 'grade_lead' | 'admin_director' | null

const classRoleLabel: Record<ClassRole, string> = { head_teacher: '主班主任', co_teacher: '副班主任', subject_teacher: '科任教师' }
const grades = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']
const subjects = ['语文', '数学', '英语', '科学', '道德与法治', '体育', '音乐', '美术', '信息技术', '综合实践']

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

function openEdit(t: Teacher) { editTarget.value = t; showEditModal.value = true }
function openAssign(t: Teacher) { assignTarget.value = t; showAssignModal.value = true }
function openDelete(t: Teacher) { deleteTarget.value = t; showDeleteModal.value = true }
function openResetPwd(t: Teacher) { resetTarget.value = t; showResetPwdModal.value = true }

function gradeIcon(team: string): string {
  for (const g of ['一年级','二年级','三年级','四年级','五年级','六年级']) { if (team.includes(g)) return ['🌱','🌿','🌳','📚','⭐','🎓'][['一年级','二年级','三年级','四年级','五年级','六年级'].indexOf(g)] }
  return '📌'
}
function shortClassName(name: string | undefined) {
  if (!name) return ''; const m = name.match(/（(\d+)）班/); return m ? m[1] + '班' : name
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

onMounted(() => loadTeachers(true))
</script>
<template>
  <div class="teachers-admin" style="max-width:1400px;margin:0 auto;padding:0 4px;">
    <TeacherFilters :grades="grades" :classRoleLabel="classRoleLabel" :filterGrade="filterGrade" :filterRole="filterRole" :searchQuery="searchQuery"
      @update:filterGrade="filterGrade = $event" @update:filterRole="filterRole = $event" @update:searchQuery="searchQuery = $event"
      @downloadTemplate="window.open('/api/v1/admin/teachers/template-csv','_blank')" @openImport="showImportModal = true" @openCreate="showCreateModal = true">
      <span class="count-badge" slot="teacherCount">{{ teachers.length }} 人</span>
    </TeacherFilters>

    <div v-if="loading" class="loading-spinner">加载中...</div>
    <div v-else-if="filteredTeachers.length === 0" class="empty-state">
      <div class="empty-icon">&#x1F468;&#x200D;&#x1F3EB;</div>
      <p v-if="searchQuery || filterGrade || filterRole">未找到匹配的教师，请修改搜索条件</p>
      <p v-else>暂无教师，点击「创建教师」添加</p>
    </div>

    <template v-for="(team, teamName) in teacherTeams" :key="teamName">
      <div v-if="team.length > 0" style="margin-bottom:20px;">
        <div class="grade-header">
          <span class="grade-icon">{{ gradeIcon(teamName) }}</span>
          <span class="grade-name">{{ teamName }}</span>
          <span class="grade-count">{{ team.length }} 人</span>
        </div>
        <div class="card-grid">
          <TeacherCard v-for="t in team" :key="t.id" :teacher="t" :classes="classes"
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
  <EditTeacherModal :visible="showEditModal" :teacher="editTarget" :grades="grades"
    @update:visible="showEditModal = $event; editTarget = null" @updated="refreshTeachers" />
  <AssignClassModal :visible="showAssignModal" :teacher="assignTarget" :classes="classes" :grades="grades" :subjects="subjects" :classRoleLabel="classRoleLabel"
    @update:visible="showAssignModal = $event; assignTarget = null" @assigned="refreshTeachers" />
  <DeleteConfirmModal :visible="showDeleteModal" :teacher="deleteTarget"
    @update:visible="showDeleteModal = $event; deleteTarget = null" @deleted="refreshTeachers" />
  <ResetPasswordModal :visible="showResetPwdModal" :teacher="resetTarget"
    @update:visible="showResetPwdModal = $event; resetTarget = null" @reset="refreshTeachers" />
  <ImportTeacherModal :visible="showImportModal" @update:visible="showImportModal = $event" @imported="refreshTeachers" />
</template>

<style scoped>
.toolbar { display:flex;align-items:flex-start;justify-content:space-between;flex-wrap:wrap;gap:12px;margin-bottom:24px; }
.toolbar-left { display:flex;align-items:center;gap:12px;flex-wrap:wrap; }
.section-badge { font-size:11px;font-weight:600;color:#7c3aed;text-transform:uppercase;letter-spacing:0.05em;background:#ede9fe;padding:3px 10px;border-radius:6px; }
.page-title { font-size:22px;font-weight:700;margin:0;line-height:1.2; }
.count-badge { font-size:13px;color:#6b7280;background:var(--color-bg);padding:2px 10px;border-radius:10px; }
.toolbar-actions { display:flex;align-items:center;gap:8px;flex-wrap:wrap; }
.filter-group { display:flex;gap:6px; }
.filter-select { width:110px;padding:6px 10px;font-size:12px; }
.filter-search { width:150px;padding:6px 10px;font-size:12px; }
.grade-header { display:flex; align-items:center; gap:10px; margin-bottom:12px; padding-bottom:8px; border-bottom:2px solid var(--color-border); }
.grade-icon { font-size:18px; }
.grade-name { font-size:17px; font-weight:700; }
.grade-count { font-size:12px; color:var(--color-text-secondary); background:var(--color-bg); padding:0 10px; border-radius:20px; }
.card-grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(340px, 1fr)); gap:14px; }
.teacher-card { background:var(--color-bg-card); border:1px solid var(--color-border); border-radius:14px; display:flex; flex-direction:column; transition:all 0.25s ease; overflow:hidden; }
.teacher-card:hover { border-color:rgba(124,58,237,0.15); box-shadow:0 8px 24px rgba(0,0,0,0.06); transform:translateY(-2px); }
.card-head { display:flex; align-items:center; gap:12px; padding:16px 18px 12px; }
.avatar { width:44px; height:44px; border-radius:50%; display:flex; align-items:center; justify-content:center; color:#fff; font-weight:700; font-size:18px; flex-shrink:0; }
.head-body { flex:1; min-width:0; }
.head-top { display:flex; align-items:center; gap:6px; flex-wrap:wrap; }
.head-name { font-size:17px; font-weight:700; }
.head-id { font-size:11px; color:var(--color-text-secondary); margin-top:2px; }
.head-badge { font-size:10px; font-weight:600; padding:2px 10px; border-radius:10px; white-space:nowrap; }
.badge-lead { background:#8b5cf6; color:#fff; }
.badge-admin { background:#f59e0b; color:#fff; }
.card-info-section { display:flex; gap:16px; padding:8px 18px; background:var(--color-bg); border-top:1px solid var(--color-border); border-bottom:1px solid var(--color-border); flex-wrap:wrap; }
.info-item { font-size:12px; color:var(--color-text-secondary); }
.card-classes { padding:10px 18px; flex:1; display:flex; flex-direction:column; gap:3px; min-height:56px; }
.class-row { display:flex; align-items:center; gap:6px; padding:6px 10px; background:var(--color-bg); border-radius:6px; font-size:13px; }
.class-name { font-weight:600; color:var(--color-text); min-width:70px; }
.class-subj { font-size:12px; color:var(--color-text-secondary); }
.role-tag-head { font-size:10px; font-weight:600; padding:1px 8px; border-radius:8px; background:rgba(250,204,21,0.15); color:#a16207; }
.role-tag-co { font-size:10px; font-weight:600; padding:1px 8px; border-radius:8px; background:rgba(156,163,175,0.15); color:#6b7280; }
.class-empty { padding:8px; text-align:center; font-size:12px; color:var(--color-text-secondary); }
.card-actions { display:flex; gap:6px; padding:10px 18px 12px; border-top:1px solid var(--color-border); }
.act-btn { padding:5px 12px; border-radius:8px; font-size:12px; font-weight:500; cursor:pointer; transition:0.15s; border:1px solid var(--color-border); background:var(--color-bg-card); color:var(--color-text); font-family:inherit; display:inline-flex; align-items:center; gap:4px; }
.act-btn:hover { background:var(--color-bg); }
.act-del { color:#dc2626; border-color:rgba(239,68,68,0.2); }
.act-del:hover { background:rgba(239,68,68,0.06); }
.loading-spinner { text-align:center;padding:64px;color:#9ca3af;font-size:15px; }
.empty-state { text-align:center;padding:64px 20px;color:#9ca3af; }
.empty-icon { font-size:48px;margin-bottom:12px; }
</style>
