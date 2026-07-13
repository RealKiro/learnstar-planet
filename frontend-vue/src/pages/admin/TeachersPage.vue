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
  head_teacher: '班主任', co_teacher: '副班', subject_teacher: '科任教师',
  grade_lead: '年级首席', admin_director: '分管行政',
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

// ──── 卡片列表／表格视图 ────
const viewMode = ref<'card' | 'table'>('card')

// ──── 筛选 ────
const filterGrade = ref('')
const filterRole = ref('')
const searchQuery = ref('')

const filteredTeachers = computed(() => {
  let list = teachers.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(t => t.name.includes(q) || t.username.includes(q) || (t.nickname && t.nickname.includes(q)))
  }
  if (filterGrade.value) {
    list = list.filter(t => t.grade_team === filterGrade.value)
  }
  if (filterRole.value) {
    list = list.filter(t => t.assignments.some(a => a.role === filterRole.value))
  }
  return list
})

// ──── 单个创建弹窗 ────
const showCreateModal = ref(false)
const createForm = ref({ name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' })
const createAssignments = ref<{ class_id: number | null; role: Role; subject: string }[]>([])
const createLoading = ref(false)

function openCreateModal() {
  createForm.value = { name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' }
  createAssignments.value = []
  showCreateModal.value = true
}
function addCreateAssignment() {
  createAssignments.value.push({ class_id: null, role: 'subject_teacher', subject: '' })
}
function removeCreateAssignment(idx: number) { createAssignments.value.splice(idx, 1) }
async function submitCreate() {
  createLoading.value = true
  try {
    const payload: any = { ...createForm.value }
    if (createAssignments.value.length > 0 && createAssignments.value.every(a => a.class_id)) {
      payload.assignments = createAssignments.value.map(a => ({
        class_id: a.class_id, role: a.role, subject: a.subject || undefined,
      }))
    }
    await apiPost('/api/v1/admin/teachers', payload)
    toast.show('教师创建成功', 'success')
    showCreateModal.value = false
    await refreshTeachers()
  } finally { createLoading.value = false }
}

// ──── 编辑教师弹窗 ────
const showEditModal = ref(false)
const editTarget = ref<Teacher | null>(null)
const editForm = ref({ name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '' })
function openEditModal(t: Teacher) {
  editTarget.value = t
  editForm.value = {
    name: t.name, nickname: t.nickname || '',
    subject: t.subject || '', grade_team: t.grade_team || '',
    phone: t.phone || '', email: t.email || '',
  }
  showEditModal.value = true
}
async function submitEdit() {
  if (!editTarget.value) return
  try {
    await apiPut(`/api/v1/admin/teachers/${editTarget.value.id}`, editForm.value)
    toast.show('教师信息已更新', 'success')
    showEditModal.value = false
    await refreshTeachers()
  } catch { /* handled */ }
}

// ──── 分配班级弹窗 ────
const showAssignModal = ref(false)
const assignTarget = ref<Teacher | null>(null)
const assignList = ref<{ class_id: number | null; role: Role; subject: string }[]>([])
const assignLoading = ref(false)
function openAssignModal(t: Teacher) {
  assignTarget.value = t
  assignList.value = t.assignments.length > 0
    ? t.assignments.map(a => ({ class_id: a.class_id, role: a.role as Role, subject: a.subject || '' }))
    : [{ class_id: null, role: 'subject_teacher', subject: '' }]
  showAssignModal.value = true
}
function addAssignRow() { assignList.value.push({ class_id: null, role: 'subject_teacher', subject: '' }) }
function removeAssignRow(idx: number) { assignList.value.splice(idx, 1) }
async function submitAssign() {
  if (!assignTarget.value) return
  assignLoading.value = true
  try {
    const payload = {
      assignments: assignList.value
        .filter(a => a.class_id)
        .map(a => ({ class_id: a.class_id, role: a.role, subject: a.subject || undefined })),
    }
    await apiPut(`/api/v1/admin/teachers/${assignTarget.value.id}/classes`, payload)
    toast.show('班级分配已更新', 'success')
    showAssignModal.value = false
    await refreshTeachers()
  } finally { assignLoading.value = false }
}

// ──── CSV 导入弹窗 ────
const showImportModal = ref(false)
const importFile = ref<File | null>(null)
const importPreview = ref<any[]>([])
const importDryRun = ref(true)
const importLoading = ref(false)
function openImportModal() { importFile.value = null; importPreview.value = []; importDryRun.value = true; showImportModal.value = true }
function onFileChange(e: Event) { importFile.value = (e.target as HTMLInputElement).files?.[0] || null; importPreview.value = [] }
async function uploadImport(isDry: boolean) {
  if (!importFile.value) return toast.show('请选择文件', 'error')
  importLoading.value = true; importDryRun.value = isDry
  try {
    const fd = new FormData(); fd.append('file', importFile.value); fd.append('dry_run', isDry ? '1' : '0')
    const res = await apiPost<{ preview?: any[]; created?: any[]; total: number; message: string }>('/api/v1/admin/teachers/import', fd)
    if (isDry && res.preview) {
      importPreview.value = res.preview
      toast.show(`预览：${res.total} 条数据，确认无误后点击"确认导入"`)
    } else {
      toast.show(res.message || '导入完成', 'success')
      showImportModal.value = false
      await refreshTeachers()
    }
  } finally { importLoading.value = false }
}

// ──── 其他 ────
async function refreshTeachers() {
  loading.value = true
  try {
    const [tRes, cRes] = await Promise.all([
      apiGet<{ data: Teacher[] }>('/api/v1/admin/teachers'),
      apiGet<{ data: ClassRoom[] }>('/api/v1/admin/classes'),
    ])
    teachers.value = tRes.data || []
    classes.value = cRes.data || []
  } finally { loading.value = false }
}
function downloadTemplate() { window.open('/api/v1/admin/teachers/template-csv', '_blank') }
async function resetPwd(t: Teacher) {
  const pwd = prompt(`为「${t.name}」设置新密码（留空自动生成）：`)
  try { await apiPost(`/api/v1/admin/teachers/${t.id}/reset-password`, { password: pwd || undefined }); toast.show('密码已重置', 'success') } catch {}
}
async function deleteTeacher(t: Teacher) {
  if (!confirm(`确定删除教师「${t.name}」？`)) return
  try { await apiDelete(`/api/v1/admin/teachers/${t.id}`); teachers.value = teachers.value.filter(x => x.id !== t.id); toast.show('已删除', 'success') } catch {}
}
function classById(id: number) { return classes.value.find(c => c.id === id) }

onMounted(refreshTeachers)
</script>

<template>
  <div class="teachers-admin" style="max-width:1400px;margin:0 auto;padding:0 4px;">
    <!-- ──── 顶部工具栏 ──── -->
    <div class="toolbar">
      <div class="toolbar-left">
        <div class="section-badge">账号管理</div>
        <h2 class="page-title">教师账号</h2>
        <span class="count-badge">{{ teachers.length }} 人</span>
      </div>
      <div class="toolbar-actions">
        <div class="filter-group">
          <select v-model="filterGrade" class="form-input filter-select">
            <option value="">全部年级</option>
            <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
          </select>
          <select v-model="filterRole" class="form-input filter-select">
            <option value="">全部角色</option>
            <option v-for="(label, role) in roleLabel" :key="role" :value="role">{{ label }}</option>
          </select>
          <input v-model="searchQuery" class="form-input filter-search" placeholder="搜索姓名 / 账号..." />
        </div>
        <div class="view-toggle">
          <button :class="{ active: viewMode === 'card' }" @click="viewMode = 'card'" title="卡片视图"><svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><rect 