<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { avatarGradient, platformLabel } from '@/utils/constants'

interface Teacher {
  id: number; name: string; username: string; nickname?: string
  subject?: string; grade_team?: string; phone?: string; email?: string
  avatar_path?: string; status: string; bindings: string[]
  assignments: Assignment[]; class_names: string[]
}
interface Assignment { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }
interface ClassRoom { id: number; name: string; grade?: string }
type Role = 'head_teacher' | 'co_teacher' | 'subject_teacher' | 'grade_lead' | 'admin_director'
const roleLabel: Record<Role, string> = {
  head_teacher: '主班', co_teacher: '副班', subject_teacher: '科任',
  grade_lead: '首席', admin_director: '主任',
}
const roleColors: Record<Role, string> = {
  head_teacher: '#7c3aed', co_teacher: '#2563eb', subject_teacher: '#059669',
  grade_lead: '#d97706', admin_director: '#dc2626',
}
const grades = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']
const subjects = ['语文', '数学', '英语', '科学', '道德与法治', '体育', '音乐', '美术', '信息技术', '综合实践']

const toast = useToastStore()
const teachers = ref<Teacher[]>([])
const classes = ref<ClassRoom[]>([])
const loading = ref(true)

const viewMode = ref<'card' | 'table'>('card')
const filterGrade = ref('')
const filterRole = ref('')
const searchQuery = ref('')

const teacherTeams = computed(() => {
  const teams: Record<string, typeof teachers.value> = {}
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
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(t => t.name.includes(q) || t.username.includes(q) || (t.nickname && t.nickname.includes(q)))
  }
  if (filterGrade.value) list = list.filter(t => t.grade_team === filterGrade.value)
  if (filterRole.value) list = list.filter(t => t.assignments.some(a => a.role === filterRole.value))
  return list
})

const showCreateModal = ref(false)
const createForm = ref({ name: '', nickname: '', grade_team: '', phone: '', email: '', password: '' })
const createErrors = ref<Record<string, string>>({})
interface CreateAssignment { class_id: number; class_name: string; subject: string; role: Role }
const createAssignments = ref<CreateAssignment[]>([])
const createLoading = ref(false)
const pendingGrade = ref('')
const pendingClassId = ref(null)
const pendingSubject = ref('')
const pendingRole = ref<Role>('subject_teacher')
const gradeClasses = computed(() => (classes.value || []).filter(c => c.grade === pendingGrade.value))

function clearError(field: string) { delete createErrors.value[field] }

function validateField(field: string, value: string) {
  if (field === 'name' && !value.trim()) { createErrors.value.name = '请填写教师姓名'; return false }
  if (field === 'password' && value && value.length < 6) { createErrors.value.password = '密码长度至少 6 位，建议包含字母和数字'; return false }
  delete createErrors.value[field]; return true
}

function openCreateModal() {
  createForm.value = { name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' }
  createErrors.value = {}; assignError.value = ''
  createAssignments.value = []; pendingGrade.value = ''; pendingClassId.value = null; pendingSubject.value = ''
  showCreateModal.value = true
}
const assignError = ref('')
function addClassAssignment() {
  assignError.value = ''
  if (!pendingClassId.value || !pendingGrade.value) { assignError.value = '请先选择年级和班级'; return }
  const cls = classes.value.find(c => c.id === pendingClassId.value)
  if (!cls) return
  if (createAssignments.value.some(a => a.class_id === cls.id)) { assignError.value = '该班级已添加'; return }
  createAssignments.value.push({ class_id: cls.id, class_name: cls.name, subject: pendingSubject.value || '默认科目', role: pendingRole.value })
  pendingClassId.value = null; pendingSubject.value = ''; pendingRole.value = 'subject_teacher'
}
function removeClassAssignment(idx) { createAssignments.value.splice(idx, 1) }
async function submitCreate() {
  // 提交前整体校验
  createErrors.value = {}
  if (!createForm.value.name.trim()) createErrors.value.name = '请填写教师姓名'
  if (createForm.value.password && createForm.value.password.length < 6) createErrors.value.password = '密码长度至少 6 位，建议包含字母和数字'
  if (Object.keys(createErrors.value).length > 0) return
  createLoading.value = true
  try {
    const payload: any = { ...createForm.value }
    if (createAssignments.value.length > 0 && createAssignments.value.every(a => a.class_id)) {
      payload.assignments = createAssignments.value.map(a => ({ class_id: a.class_id, role: a.role || 'subject_teacher', subject: a.subject || undefined }))
    }
    await apiPost('/api/v1/admin/teachers', payload)
    toast.show('教师创建成功', 'success')
    showCreateModal.value = false; await refreshTeachers()
  } catch (e: any) {
    const errs = e?.response?.data?.errors
    if (errs) {
      for (const [field, msgs] of Object.entries(errs)) {
        createErrors.value[field] = (msgs as string[])[0]
      }
    } else {
      const msg = e?.response?.data?.message
      if (msg && e?.response?.status === 500) toast.show(msg, 'error')
    }
  } finally { createLoading.value = false }
}

const showEditModal = ref(false)
const editTarget = ref<Teacher | null>(null)
const editForm = ref({ name: '', nickname: '', grade_team: '', phone: '', email: '', personalRole: '' as string })
function openEditModal(t: Teacher) {
  editTarget.value = t
  const roles = uniqueRoles(t.assignments)
  const pRole = roles.includes('grade_lead') ? 'grade_lead' : roles.includes('admin_director') ? 'admin_director' : ''
  editForm.value = { name: t.name, nickname: t.nickname || '', grade_team: t.grade_team || '', phone: t.phone || '', email: t.email || '', personalRole: pRole }
  showEditModal.value = true
}
async function submitEdit() {
  if (!editTarget.value) return
  try {
    const payload: any = { name: editForm.value.name, nickname: editForm.value.nickname, grade_team: editForm.value.grade_team, phone: editForm.value.phone, email: editForm.value.email, personal_role: editForm.value.personalRole }
    await apiPut(`/api/v1/admin/teachers/${editTarget.value.id}`, payload)
    toast.show('教师信息已更新', 'success')
    showEditModal.value = false; await refreshTeachers()
  } catch { toast.show('保存失败', 'error') }
}

const showAssignModal = ref(false)
const assignTarget = ref<Teacher | null>(null)
const assignList = ref<{ class_id: number | null; role: Role; subject: string; class_name?: string }[]>([])
const assignLoading = ref(false)
const assignGradeFilter = ref('')
const newAssignClassId = ref<number | null>(null)
const newAssignRoleSingle = ref<Role>('subject_teacher')
const newAssignSubject = ref('')
const filteredAssignClasses = computed(() => {
  if (!assignGradeFilter.value) return classes.value
  return classes.value.filter(c => c.grade === assignGradeFilter.value)
})
function openAssignModal(t: Teacher) {
  assignTarget.value = t
  assignList.value = t.assignments.length > 0
    ? t.assignments.map(a => ({ class_id: a.class_id, role: a.role as Role, subject: a.subject || '', class_name: a.class_name || classById(a.class_id)?.name }))
    : []
  assignGradeFilter.value = ''; newAssignClassId.value = null; newAssignRoleSingle.value = 'subject_teacher'; newAssignSubject.value = ''
  showAssignModal.value = true
}
function addAssignRowNew() {
  if (!newAssignClassId.value) return
  const cls = classes.value.find(c => c.id === newAssignClassId.value)
  if (!cls) return
  const role = newAssignRoleSingle.value
  // 主班/副班互斥校验
  if (role === 'head_teacher' && assignList.value.some(a => a.class_id === cls.id && a.role === 'co_teacher')) { toast.show('该班级已有副班，不能同时为主班', 'error'); return }
  if (role === 'co_teacher' && assignList.value.some(a => a.class_id === cls.id && a.role === 'head_teacher')) { toast.show('该班级已有主班，不能同时为副班', 'error'); return }
  if (assignList.value.some(a => a.class_id === cls.id && a.role === role)) { toast.show('该班级已分配此角色', 'info'); return }
  assignList.value.push({ class_id: cls.id, class_name: shortClassName(cls.name), role, subject: newAssignSubject.value || '默认科目' })
  newAssignClassId.value = null; newAssignSubject.value = ''
}
function removeAssignRow(idx: number) { assignList.value.splice(idx, 1) }
async function submitAssign() {
  if (!assignTarget.value) return
  assignLoading.value = true
  try {
    const payload = { assignments: assignList.value.filter(a => a.class_id).map(a => ({ class_id: a.class_id, role: a.role, subject: a.subject === '默认科目' ? undefined : a.subject || undefined })) }
    await apiPut(`/api/v1/admin/teachers/${assignTarget.value.id}/classes`, payload)
    toast.show('班级分配已更新', 'success')
    showAssignModal.value = false; await refreshTeachers()
  } catch (e: any) {
    toast.show(e?.response?.data?.message || '保存失败', 'error')
  } finally { assignLoading.value = false }
}

const showDeleteModal = ref(false)
const deleteTarget = ref<Teacher | null>(null)
const deleteConfirmed = ref(false)
function openDeleteModal(t: Teacher) {
  deleteTarget.value = t; deleteConfirmed.value = false; showDeleteModal.value = true
}
async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    await apiDelete(`/api/v1/admin/teachers/${deleteTarget.value.id}`)
    teachers.value = teachers.value.filter(x => x.id !== deleteTarget.value!.id)
    toast.show('教师已删除', 'success'); showDeleteModal.value = false
  } catch { toast.show('删除失败', 'error') }
}

const showImportModal = ref(false)
const importFile = ref<File | null>(null)
const importPreview = ref<any[]>([])
const importLoading = ref(false)
function openImportModal() { importFile.value = null; importPreview.value = []; showImportModal.value = true }
function onFileChange(e: Event) { importFile.value = (e.target as HTMLInputElement).files?.[0] || null; importPreview.value = [] }
async function uploadImport(isDry: boolean) {
  if (!importFile.value) return toast.show('请选择文件', 'error')
  importLoading.value = true
  try {
    const fd = new FormData(); fd.append('file', importFile.value); fd.append('dry_run', isDry ? '1' : '0')
    const res = await apiPost<{ preview?: any[]; created?: any[]; total: number; message: string }>('/api/v1/admin/teachers/import', fd)
    if (isDry && res.preview) { importPreview.value = res.preview; toast.show('预览：' + res.total + ' 条数据') }
    else { toast.show(res.message || '导入完成', 'success'); showImportModal.value = false; await refreshTeachers() }
  } finally { importLoading.value = false }
}

async function loadTeachers(isInitial = false) {
  if (isInitial) loading.value = true
  try {
    const [tRes, cRes] = await Promise.all([
      apiGet<{ data: Teacher[] }>('/api/v1/admin/teachers'),
      apiGet<{ data: ClassRoom[] }>('/api/v1/admin/classes'),
    ])
    teachers.value = tRes.data || []; classes.value = cRes.data || []
  } catch { /* 静默失败，保留现有数据 */ }
  finally { if (isInitial) loading.value = false }
}
const refreshTeachers = () => loadTeachers(false)
function downloadTemplate() { window.open('/api/v1/admin/teachers/template-csv', '_blank') }
const showResetPwdModal = ref(false)
const resetTarget = ref<Teacher | null>(null)
const resetPwdValue = ref('')
const resetPwdLoading = ref(false)
const showResetPwd = ref(false)
const currentPwd = ref('')

function openResetPwd(t: Teacher) {
  resetTarget.value = t
  resetPwdValue.value = ''
  showResetPwd.value = false
  currentPwd.value = ''
  showResetPwdModal.value = true
  // 加载当前密码
  fetch('/api/v1/admin/teachers/' + t.id + '/password', {
    headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
  }).then(r => r.json()).then(d => { currentPwd.value = d.data?.password || '' }).catch(() => {})
}

async function submitResetPwd() {
  if (!resetTarget.value) return
  resetPwdLoading.value = true
  try {
    const res = await apiPost<{ data: { new_password: string } }>(`/api/v1/admin/teachers/${resetTarget.value.id}/reset-password`, { password: resetPwdValue.value || undefined })
    const newPwd = res?.data?.new_password || resetPwdValue.value || '（已设置）'
    toast.show(`密码已更新: ${newPwd}`, 'success')
    showResetPwdModal.value = false
  } catch { toast.show('重置失败', 'error') }
  finally { resetPwdLoading.value = false }
}
}
function classById(id: number) { return classes.value.find(c => c.id === id) }
function uniqueRoles(assignments: Assignment[]): string[] {
  return [...new Set(assignments.map(a => a.role))]
}
function roleTagClass(role: string): string {
  const map: Record<string, string> = { head_teacher: 'tag-head', co_teacher: 'tag-co', subject_teacher: 'tag-subj', grade_lead: 'tag-lead', admin_director: 'tag-admin' }
  return map[role] || 'tag-subj'
}
function gradeIcon(team: string): string {
  for (const g of ['一年级','二年级','三年级','四年级','五年级','六年级']) { if (team.includes(g)) return ['🌱','🌿','🌳','📚','⭐','🎓'][['一年级','二年级','三年级','四年级','五年级','六年级'].indexOf(g)] }
  return '📌'
}
function shortClassName(name: string | undefined) {
  if (!name) return ''
  const m = name.match(/（(\d+)）班/)
  return m ? m[1] + '班' : name
}

onMounted(() => loadTeachers(true))
</script>

<template>
  <div class="teachers-admin" style="max-width:1400px;margin:0 auto;padding:0 4px;">
    <div class="toolbar">
      <div class="toolbar-left">
        <div class="section-badge">账号管理</div>
        <h2 class="page-title">教师账号</h2>
        <span class="count-badge">{{ teachers.length }} 人</span>
      </div>
      <div class="toolbar-actions">
        <div class="filter-group">
          <select v-model="filterGrade" class="form-input filter-select">
            <option value="">-- 请选择年级 --</option>
            <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
          </select>
          <select v-model="filterRole" class="form-input filter-select">
            <option value="">全部角色</option>
            <option v-for="(label, role) in roleLabel" :key="role" :value="role">{{ label }}</option>
          </select>
          <input v-model="searchQuery" class="form-input filter-search" placeholder="搜索姓名 / 账号..." />
        </div>
        <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);font-size:13px;" @click="downloadTemplate">下载模板</button>
        <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);font-size:13px;" @click="openImportModal">批量导入</button>
        <button class="btn btn-primary" style="font-size:13px;" @click="openCreateModal">+ 创建教师</button>
      </div>
    </div>

    <div v-if="loading" class="loading-spinner">加载中...</div>

    <div v-else-if="filteredTeachers.length === 0" class="empty-state">
      <div class="empty-icon">&#x1F468;&#x200D;&#x1F3EB;</div><p>暂无教师</p>
    </div>

    <!-- 按年级团队分组 -->
    <template v-for="(teachers, teamName) in teacherTeams" :key="teamName">
      <div v-if="teachers.length > 0" style="margin-bottom:20px;">
        <div class="grade-header">
          <span class="grade-icon">{{ gradeIcon(teamName) }}</span>
          <span class="grade-name">{{ teamName }}</span>
          <span class="grade-count">{{ teachers.length }} 人</span>
        </div>
        <div class="card-grid">
          <div v-for="t in teachers" :key="t.id" class="teacher-card">
            <!-- 层1：头部 — 头像 + 姓名 + 个人角色 + 编号 -->
            <div class="card-head">
              <div class="avatar" :style="{ background: avatarGradient(t.name) }">{{ t.name[0] }}</div>
              <div class="head-body">
                <div class="head-top">
                  <span class="head-name">{{ t.name }}</span>
                  <span v-for="role in uniqueRoles(t.assignments)" :key="role" v-if="role === 'grade_lead' || role === 'admin_director'" :class="['head-badge', role === 'grade_lead' ? 'badge-lead' : 'badge-admin']">{{ role === 'grade_lead' ? '首席' : '主任' }}</span>
                </div>
                <div class="head-id">{{ t.username }}</div>
              </div>
            </div>

            <!-- 层2：联系信息 -->
            <div class="card-info-section" v-if="t.email || t.phone">
              <span v-if="t.email" class="info-item">✉ {{ t.email }}</span>
              <span v-if="t.phone" class="info-item">📱 {{ t.phone }}</span>
            </div>

            <!-- 层3：任教班级 -->
            <div class="card-classes">
              <div v-for="a in t.assignments" :key="a.class_id + '_' + a.role" class="class-row">
                <span class="class-name">{{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}</span>
                <span v-if="a.role === 'head_teacher'" class="role-tag-head">主班</span>
                <span v-else-if="a.role === 'co_teacher'" class="role-tag-co">副班</span>
                <span class="class-subj">{{ a.role === 'head_teacher' || a.role === 'co_teacher' ? a.subject : (a.subject || '—') }}</span>
              </div>
              <div v-if="t.assignments.length === 0" class="class-empty">暂未分配班级</div>
            </div>

            <!-- 层4：操作按钮 -->
            <div class="card-actions">
              <button class="act-btn" @click="openEditModal(t)">👤 个人信息</button>
              <button class="act-btn" @click="openAssignModal(t)">📚 班级管理</button>
              <button class="act-btn" @click="openResetPwd(t)">🔑 密码</button>
              <button class="act-btn act-del" @click="openDeleteModal(t)">🗑️ 删除</button>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>

  <Teleport to="body">
    <Transition name="overlay">
      <div v-if="showCreateModal" @click="showCreateModal = false" style="position:fixed;inset:0;z-index:999;background:rgba(0,0,0,0.2);display:flex;align-items:center;justify-content:flex-end;">
        <Transition name="drawer">
          <div v-if="showCreateModal" @click.stop style="width:540px;max-width:100vw;max-height:90vh;margin:20px;border-radius:16px;background:var(--color-bg-card);border:1px solid var(--color-border);box-shadow:-8px 8px 32px rgba(0,0,0,0.12);display:flex;flex-direction:column;overflow:hidden;">
            <div style="display:flex;justify-content:space-between;align-items:center;padding:16px 20px;border-bottom:1px solid var(--color-border);flex-shrink:0;">
              <h2 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">✨ 创建教师账号</h2>
              <button :disabled="createLoading" @click="showCreateModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:4px;line-height:1;">✕</button>
            </div>
            <div style="overflow-y:auto;padding:16px 20px;">
              <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;">
                <div>
                  <div style="font-size:12px;font-weight:600;color:var(--color-text);margin-bottom:8px;">📝 创建账号 <span style="font-size:10px;color:var(--color-text-secondary);">必填</span></div>
                  <div style="display:flex;gap:8px;">
                    <div style="flex:1;" class="form-group">
                      <label>姓名 <span style="color:var(--color-danger);">*</span></label>
                      <input v-model="createForm.name" placeholder="姓名" class="form-input" :style="{ borderColor: createErrors.name ? '#f87171' : '' }" @blur="validateField('name', createForm.name)" @input="clearError('name')">
                      <div v-if="createErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ createErrors.name }}</div>
                    </div>
                    <div style="flex:1;" class="form-group"><label>昵称</label><input v-model="createForm.nickname" placeholder="默认拼音" class="form-input" /></div>
                  </div>
                  <div style="display:flex;gap:8px;">
                    <div style="flex:1;" class="form-group"><label>年级团队</label><select v-model="createForm.grade_team" class="form-input"><option value="">不指定</option><option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option></select></div>
                    <div style="flex:1;" class="form-group"><label>手机号</label><input v-model="createForm.phone" placeholder="选填" class="form-input" /></div>
                  </div>
                  <div style="display:flex;gap:8px;"><div style="flex:1;" class="form-group"><label>邮箱</label><input v-model="createForm.email" placeholder="选填" class="form-input" /></div><div style="flex:1;" class="form-group"><label>&nbsp;</label></div></div>
                  <div style="display:none;"></div>
                  <div class="form-group">
                    <label>初始密码</label>
                    <input v-model="createForm.password" placeholder="留空自动生成" class="form-input" :style="{ borderColor: createErrors.password ? '#f87171' : '' }" @blur="validateField('password', createForm.password)" @input="clearError('password')">
                    <div v-if="createErrors.password" style="color:#f87171;font-size:11px;margin-top:2px;">{{ createErrors.password }}</div>
                  </div>
                </div>
                <div>
                  <div style="font-size:12px;font-weight:600;color:var(--color-text);margin-bottom:8px;">📚 加入班级 <span style="font-size:10px;color:var(--color-text-secondary);">可选</span></div>
                  <div style="display:flex;gap:8px;align-items:flex-end;">
                    <div style="display:flex;flex-direction:column;gap:4px;flex:1;">
                      <div class="form-group" style="margin-bottom:0;"><label>年级</label><select v-model="pendingGrade" class="form-input"><option value="">请选择</option><option v-for="g in grades" :key="g" :value="g">{{ g }}</option></select></div>
                      <div class="form-group" style="margin-bottom:0;"><label>班级</label><select v-model="pendingClassId" :disabled="!pendingGrade" class="form-input"><option :value="null">请选择</option><option v-for="c in gradeClasses" :key="c.id" :value="c.id">{{ shortClassName(c.name) }}</option></select></div>
                      <div class="form-group" style="margin-bottom:0;"><label>角色</label><select v-model="pendingRole" class="form-input"><option value="subject_teacher">科任教师</option><option v-for="(label, role) in roleLabel" :key="role" :value="role">{{ label }}</option></select></div>
                    </div>
                    <button @click="addClassAssignment" :disabled="!pendingClassId" style="padding:6px 16px;border-radius:8px;border:1px solid var(--color-accent);background:rgba(79,70,229,0.08);color:var(--color-accent);font-size:13px;cursor:pointer;font-weight:500;white-space:nowrap;height:36px;">➕ 添加</button>
                  </div>
                  <div v-if="assignError" style="color:#f87171;font-size:11px;margin-top:4px;">{{ assignError }}</div>
                  <div style="font-size:11px;color:var(--color-text-secondary);margin-top:8px;padding:6px 8px;background:var(--color-bg);border-radius:6px;border-left:2px solid var(--color-accent);">💡 创建后可随时在列表中点击 🏫 按钮重新分配班级</div>
                  <div v-if="createAssignments.length > 0" style="margin-top:8px;display:flex;flex-direction:column;gap:4px;">
                    <div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:2px;">📋 已添加（{{ createAssignments.length }}）</div>
                    <div v-for="(a, i) in createAssignments" :key="i" style="display:flex;align-items:center;gap:6px;padding:4px 8px;background:var(--color-bg);border-radius:4px;font-size:12px;border-left:3px solid var(--color-accent);">
                      <span style="flex:1;color:var(--color-text);">{{ shortClassName(a.class_name) }}</span>
                      <span style="font-size:11px;font-weight:500;color:var(--color-accent);">{{ roleLabel[a.role] || a.role }}</span>
                      <button @click="removeClassAssignment(i)" style="background:none;border:none;color:var(--color-danger);cursor:pointer;padding:0;font-size:14px;">✕</button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div style="display:flex;gap:8px;padding:12px 20px;border-top:1px solid var(--color-border);flex-shrink:0;">
              <button @click="showCreateModal = false" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
              <button @click="submitCreate" :disabled="createLoading" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;box-shadow:0 2px 8px rgba(124,58,237,0.15);">{{ createLoading ? '创建中...' : '创建账号' }}</button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>

  <Teleport to="body">
    <div v-if="showEditModal" @click="showEditModal = false" style="position:fixed;inset:0;z-index:1000;background:rgba(0,0,0,0.15);display:flex;align-items:center;justify-content:center;padding:20px;">
      <div @click.stop style="background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:16px;max-width:520px;width:100%;padding:24px;box-shadow:0 8px 32px rgba(0,0,0,0.12);">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
          <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">✏️ 编辑教师信息 — {{ editTarget?.name }}</h3>
          <button @click="showEditModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
        </div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;">
          <div class="form-group"><label>姓名 <span style="color:#f87171;">*</span></label><input v-model="editForm.name" class="form-input" placeholder="教师姓名"></div>
          <div class="form-group"><label>昵称</label><input v-model="editForm.nickname" class="form-input" placeholder="默认拼音"></div>
          <div class="form-group"><label>年级团队</label><select v-model="editForm.grade_team" class="form-input"><option value="">不指定</option><option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option></select></div>
          <div class="form-group">
            <label>个人角色</label>
            <div style="display:flex;gap:8px;margin-top:4px;">
              <label v-for="opt in ([{v:'grade_lead',l:'首席'},{v:'admin_director',l:'主任'},{v:'',l:'普通科任'}])" :key="opt.v" style="display:flex;align-items:center;gap:4px;font-size:12px;cursor:pointer;padding:4px 10px;border-radius:6px;border:1px solid var(--color-border);background:editForm.personalRole === opt.v ? (opt.v === 'grade_lead' ? '#8b5cf6' : opt.v === 'admin_director' ? '#f59e0b' : 'var(--color-primary)') : 'var(--color-bg-card)';color:editForm.personalRole === opt.v ? '#fff' : 'var(--color-text)';">
                <input type="radio" :value="opt.v" v-model="editForm.personalRole" style="display:none;"> {{ opt.l }}
              </label>
            </div>
          </div>
          <div class="form-group"><label>手机号</label><input v-model="editForm.phone" class="form-input" placeholder="11位手机号"></div>
          <div class="form-group"><label>邮箱</label><input v-model="editForm.email" class="form-input" placeholder="邮箱地址"></div>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:20px;padding-top:12px;border-top:1px solid var(--color-border);">
          <button @click="showEditModal = false" style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
          <button @click="submitEdit" style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;">保存</button>
        </div>
      </div>
    </div>
  </Teleport>

  <Teleport to="body">
    <div v-if="showAssignModal" @click.self="showAssignModal = false" style="position:fixed;inset:0;z-index:1000;background:transparent;display:flex;align-items:center;justify-content:center;pointer-events:none;">
      <div @click.stop style="pointer-events:auto;background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:16px;max-width:620px;width:92%;padding:24px;box-shadow:0 8px 32px rgba(0,0,0,0.12);">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
          <h2 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">🏫 分配班级 — {{ assignTarget?.name }}</h2>
          <button @click="showAssignModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
        </div>
        <!-- 选择器 -->
        <div style="display:flex;gap:8px;align-items:flex-end;margin-bottom:12px;flex-wrap:wrap;">
          <div style="flex:1;min-width:100px;" class="form-group"><label>年级</label><select v-model="assignGradeFilter" class="form-input"><option value="">全部</option><option v-for="g in grades" :key="g" :value="g">{{ g }}</option></select></div>
          <div style="flex:1;min-width:120px;" class="form-group"><label>班级</label><select v-model="newAssignClassId" class="form-input"><option :value="null">请选择</option><option v-for="c in filteredAssignClasses" :key="c.id" :value="c.id">{{ shortClassName(c.name) }}</option></select></div>
          <div style="flex:1;" class="form-group"><label>班级角色（可多选）</label><div style="display:flex;flex-wrap:wrap;gap:6px;margin-top:4px;"><button v-for="(label, role) in roleLabel" :key="role" type="button" :style="{ padding:'4px 14px', borderRadius:'16px', border:'1px solid', cursor:'pointer', fontSize:'12px', fontWeight:500, fontFamily:'inherit', background: newAssignRoles.includes(role) ? roleColors[role as Role] : 'var(--color-bg-card)', color: newAssignRoles.includes(role) ? '#fff' : 'var(--color-text-secondary)', borderColor: newAssignRoles.includes(role) ? roleColors[role as Role] : 'var(--color-border)' }" @click="newAssignRoles.includes(role) ? newAssignRoles = newAssignRoles.filter(r => r !== role) : newAssignRoles.push(role)">{{ label }}</button></div></div>
          <div style="flex:1;min-width:90px;" class="form-group"><label>科目</label><select v-model="newAssignSubject" class="form-input"><option value="">默认</option><option v-for="s in subjects" :key="s" :value="s">{{ s }}</option></select></div>
          <button @click="addAssignRowNew" :disabled="!newAssignClassId || newAssignRoles.length === 0" style="padding:8px 16px;border-radius:8px;border:1px solid var(--color-accent);background:rgba(79,70,229,0.08);color:var(--color-accent);font-size:13px;cursor:pointer;font-weight:500;white-space:nowrap;height:36px;flex-shrink:0;align-self:flex-end;">➕ 添加</button>
        </div>
        <!-- 已分配列表 -->
        <div style="margin-bottom:12px;">
          <div style="font-size:13px;font-weight:600;color:var(--color-text);margin-bottom:8px;">📋 已分配（{{ assignList.length }}）</div>
          <div v-if="assignList.length === 0" style="padding:12px;text-align:center;font-size:13px;color:var(--color-text-secondary);background:var(--color-bg);border-radius:8px;">暂未分配班级</div>
          <div v-for="(a, i) in assignList" :key="i" style="display:flex;align-items:center;gap:8px;padding:8px 12px;margin-bottom:4px;background:var(--color-bg);border-radius:8px;font-size:13px;">
            <span style="flex:1;font-weight:500;color:var(--color-text);">{{ a.class_name || shortClassName(classById(a.class_id)?.name) }}</span>
            <span style="color:var(--color-text-secondary);font-size:12px;">{{ roleLabel[a.role as Role] || a.role }}</span>
            <span style="color:var(--color-text-secondary);font-size:12px;">{{ a.subject || '默认科目' }}</span>
            <button @click="removeAssignRow(i)" style="background:none;border:none;color:var(--color-danger);cursor:pointer;padding:2px;font-size:16px;">✕</button>
          </div>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;padding-top:12px;border-top:1px solid var(--color-border);">
          <button @click="showAssignModal = false" style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
          <button @click="submitAssign" :disabled="assignLoading" style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;">{{ assignLoading ? '保存中...' : '保存分配' }}</button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- 删除确认弹窗 -->
  <Teleport to="body">
    <div v-if="showDeleteModal" @click="showDeleteModal = false" style="position:fixed;inset:0;z-index:1000;background:rgba(0,0,0,0.15);display:flex;align-items:center;justify-content:center;padding:20px;">
      <div @click.stop style="background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:16px;padding:24px;max-width:400px;width:100%;box-shadow:0 8px 32px rgba(0,0,0,0.12);">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
          <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">⚠️ 确认删除</h3>
          <button @click="showDeleteModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
        </div>
        <div style="text-align:center;padding:8px 0;">
          <div style="font-size:40px;margin-bottom:8px;">🗑️</div>
          <p style="font-size:15px;font-weight:600;color:var(--color-text);margin-bottom:8px;">确定要删除教师「{{ deleteTarget?.name }}」吗？</p>
          <p style="font-size:12px;color:var(--color-text-secondary);margin-bottom:12px;">此操作将永久删除该教师账号及其所有班级任教分配，且不可恢复。</p>
          <div v-if="deleteTarget?.assignments.some(a => a.role === 'head_teacher')" style="font-size:12px;color:#dc2626;padding:8px;background:rgba(239,68,68,0.06);border-radius:6px;margin-bottom:12px;">
            ⚠️ 该教师是部分班级的班主任，删除后这些班级将无班主任。
          </div>
          <label style="display:flex;align-items:center;gap:6px;font-size:12px;color:var(--color-text-secondary);cursor:pointer;justify-content:center;padding:8px;background:var(--color-bg);border-radius:6px;">
            <input type="checkbox" v-model="deleteConfirmed" style="accent-color:#dc2626;"> 确认我已了解此操作不可恢复
          </label>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;padding-top:12px;border-top:1px solid var(--color-border);">
          <button @click="showDeleteModal = false" style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
          <button @click="confirmDelete" :disabled="!deleteConfirmed" style="padding:8px 20px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#dc2626;border:none;color:#fff;opacity:deleteConfirmed ? 1 : 0.5;">确认删除</button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- 重置密码弹窗 -->
  <Teleport to="body">
    <div v-if="showResetPwdModal" @click="showResetPwdModal = false" style="position:fixed;inset:0;z-index:1000;background:rgba(0,0,0,0.15);display:flex;align-items:center;justify-content:center;padding:20px;">
      <div @click.stop style="background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:16px;padding:24px;max-width:420px;width:100%;box-shadow:0 8px 32px rgba(0,0,0,0.12);">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
          <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">🔑 密码管理 — {{ resetTarget?.name }}</h3>
          <button @click="showResetPwdModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
        </div>

        <!-- 当前密码 -->
        <div v-if="currentPwd" style="margin-bottom:12px;padding:10px 12px;background:var(--color-bg);border-radius:8px;border:1px solid var(--color-border);">
          <div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:4px;">当前密码</div>
          <div style="display:flex;align-items:center;gap:8px;">
            <code style="font-size:14px;font-weight:700;color:var(--color-text);flex:1;font-family:monospace;">{{ showResetPwd ? currentPwd : '••••••••' }}</code>
            <button style="flex-shrink:0;padding:4px 10px;border-radius:6px;border:1px solid var(--color-border);background:var(--color-bg-card);cursor:pointer;font-size:12px;" @click="showResetPwd = !showResetPwd" type="button">{{ showResetPwd ? '🙈 隐藏' : '👁️ 显示' }}</button>
          </div>
        </div>

        <!-- 修改密码 -->
        <div class="form-group"><label>新密码</label>
          <div style="display:flex;gap:6px;">
            <input v-model="resetPwdValue" :type="showResetPwd ? 'text' : 'password'" class="form-input" placeholder="留空自动生成" autocomplete="new-password" style="flex:1;">
            <button type="button" style="flex-shrink:0;padding:6px 10px;border-radius:6px;border:1px solid var(--color-border);background:var(--color-bg-card);cursor:pointer;font-size:12px;" @click="showResetPwd = !showResetPwd">{{ showResetPwd ? '🙈' : '👁️' }}</button>
          </div>
        </div>
        <div style="display:flex;gap:8px;margin-bottom:12px;">
          <button type="button" style="flex:1;padding:6px;border-radius:6px;font-size:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text);font-family:inherit;" @click="resetPwdValue = Array.from({length:12},()=>'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*'[Math.floor(Math.random()*72)]).join(''); showResetPwd = true">✨ 生成强密码</button>
          <button type="button" style="flex:1;padding:6px;border-radius:6px;font-size:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text-secondary);font-family:inherit;" @click="resetPwdValue = ''">🔄 重置为空</button>
        </div>
        <div style="display:flex;gap:12px;">
          <button @click="showResetPwdModal = false" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
          <button @click="submitResetPwd" :disabled="resetPwdLoading" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;">{{ resetPwdLoading ? '更新中...' : '更新密码' }}</button>
        </div>
      </div>
    </div>
  </Teleport>

  <Teleport to="body">
    <div v-if="showImportModal" class="modal-overlay" >
      <div class="modal-panel" style="max-width:700px;max-height:90vh;overflow-y:auto;">
        <div class="modal-header"><h3>批量导入教师</h3><button class="close-btn" @click="showImportModal = false">&times;</button></div>
        <div class="modal-body">
          <div style="margin-bottom:16px;">
            <p style="font-size:13px;color:#6b7280;margin-bottom:8px;">支持 CSV / Excel 格式。模板列：</p>
            <div class="column-hint">
              <code>姓名</code><span style="color:#dc2626;">*</span> <code>年级团队</code> <code>科目</code> <code>密码</code> <code>手机号</code>
            </div>
            <p style="font-size:12px;color:#9ca3af;">密码选填，不填默认为 star123456。角色和班级导入后使用 &#x1F3EB; 按钮分配。</p>
            <a href="/api/v1/admin/teachers/template-csv" target="_blank" style="font-size:13px;color:#7c3aed;text-decoration:underline;">下载正确模板</a>
          </div>
          <input type="file" accept=".csv,.xlsx,.xls" @change="onFileChange" style="margin-bottom:12px;" />
          <div v-if="importPreview.length > 0" style="margin-bottom:12px;">
            <div style="font-weight:600;font-size:14px;margin-bottom:8px;">预览 ({{ importPreview.length }} 条)</div>
            <div class="preview-table-wrapper">
              <table class="preview-table">
                <thead><tr><th>姓名</th><th>年级团队</th><th>科目</th><th>手机号</th></tr></thead>
                <tbody>
                  <tr v-for="(row, i) in importPreview.slice(0, 20)" :key="i">
                    <td>{{ row.name }}</td><td>{{ row.grade_team }}</td><td>{{ row.subject }}</td><td>{{ row.phone }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p v-if="importPreview.length > 20" style="font-size:12px;color:#9ca3af;margin-top:4px;">仅显示前 20 条</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);" @click="showImportModal = false">取消</button>
          <button v-if="importPreview.length === 0" class="btn btn-primary" :disabled="importLoading || !importFile" @click="uploadImport(true)">{{ importLoading ? '解析中...' : '预览数据' }}</button>
          <button v-if="importPreview.length > 0" class="btn" style="background:#fef3c7;color:#92400e;border:1px solid #fcd34d;" :disabled="importLoading" @click="uploadImport(true)">重新预览</button>
          <button v-if="importPreview.length > 0" class="btn btn-primary" :disabled="importLoading" @click="uploadImport(false)">{{ importLoading ? '导入中...' : '确认导入' }}</button>
        </div>
      </div>
    </div>
  </Teleport>
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
.data-table { width:100%;border-collapse:collapse;background:var(--color-bg-card);border-radius:12px;overflow:hidden;border:1px solid var(--color-border); }
.data-table th { background:#f9fafb;font-size:12px;font-weight:600;color:#6b7280;text-align:left;padding:10px 14px;border-bottom:1px solid #e5e7eb;white-space:nowrap; }
.data-table td { padding:10px 14px;border-bottom:1px solid #f3f4f6;font-size:13px; }
.modal-overlay { position:fixed;top:0;right:0;bottom:0;left:220px;background:transparent;display:flex;align-items:center;justify-content:center;z-index:1000;pointer-events:none; }
.modal-overlay > * { pointer-events:auto; }
.modal-panel { background:var(--color-bg-card);border-radius:16px;padding:28px 32px;width:92%;box-shadow:0 20px 60px rgba(0,0,0,0.12); }
.modal-header { display:flex;justify-content:space-between;align-items:center;margin-bottom:20px; }
.modal-header h3 { font-size:17px;font-weight:700;margin:0; }
.close-btn { background:none;border:none;font-size:22px;cursor:pointer;color:#9ca3af;padding:0;line-height:1; }
.modal-body { margin-bottom:20px; }
.modal-footer { display:flex;justify-content:flex-end;gap:8px; }
.form-row { display:flex;gap:10px;margin-bottom:12px; }
.form-group { flex:1; }
.form-group label { display:block;font-size:12px;font-weight:600;color:var(--color-text);margin-bottom:4px; }
.form-input { color:var(--color-text);width:100%;padding:8px 12px;border:1px solid var(--color-border);border-radius:8px;font-size:13px;outline:none;transition:border-color 0.15s;box-sizing:border-box; }
.form-input:focus { border-color:#7c3aed;box-shadow:0 0 0 3px rgba(124,58,237,0.08); }
.form-input option { color:#1E293B; background:#fff; }
.flex-2 { flex:2; }
.flex-1 { flex:1; }
.assign-section { margin-top:14px;padding-top:14px;border-top:1px solid #f3f4f6; }
.assign-section-header { display:flex;justify-content:space-between;align-items:center;margin-bottom:8px; }
.assign-section-title { font-size:13px;font-weight:600;color:var(--color-text); }
.assign-row { display:flex;gap:6px;margin-bottom:8px;align-items:center; }
.preview-table-wrapper { max-height:300px;overflow:auto; }
.preview-table { width:100%;border-collapse:collapse;font-size:12px; }
.preview-table th { background:#f9fafb;padding:6px 8px;text-align:left;border:1px solid var(--color-border);font-size:11px;white-space:nowrap; }
.preview-table td { padding:5px 8px;border:1px solid var(--color-border);white-space:nowrap; }
.column-hint { display:flex;flex-wrap:wrap;gap:4px;margin-bottom:8px; }
.column-hint code { font-size:11px;background:var(--color-bg);padding:2px 6px;border-radius:4px; }
.btn { padding:8px 16px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;border:1px solid transparent;transition:all 0.15s; }
.btn-primary { background:#7c3aed;color:white;border-color:#7c3aed; }
.btn-primary:hover { background:#6d28d9; }
.btn-xs { padding:4px 10px;font-size:11px;border-radius:6px; }
.loading-spinner { text-align:center;padding:64px;color:#9ca3af;font-size:15px; }
.empty-state { text-align:center;padding:64px 20px;color:#9ca3af; }
.empty-icon { font-size:48px;margin-bottom:12px; }

/* 抽屉动画 */
.overlay-enter-active { transition: opacity 0.2s ease; }
.overlay-leave-active { transition: opacity 0.15s ease; }
.overlay-enter-from { opacity: 0; }
.overlay-leave-to { opacity: 0; }
.drawer-enter-active { transition: transform 0.25s ease-out; }
.drawer-leave-active { transition: transform 0.2s ease-in; }
.drawer-enter-from { transform: translateX(100%); }
.drawer-leave-to { transform: translateX(100%); }
</style>
