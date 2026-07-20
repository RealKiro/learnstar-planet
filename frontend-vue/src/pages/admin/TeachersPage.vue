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

const viewMode = ref<'card' | 'table'>('card')
const filterGrade = ref('')
const filterRole = ref('')
const searchQuery = ref('')

const teacherTeams = computed(() => {
  const teams: Record<string, typeof teachers.value> = {}
  teachers.value.forEach(t => {
    const team = t.grade_team || '未分组'
    if (!teams[team]) teams[team] = []
    teams[team].push(t)
  })
  return teams
})

const filteredTeachers = computed(() => {
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
const createForm = ref({ name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' })
const createAssignments = ref<{ class_id: number; class_name: string; subject: string }[]>([])
const createLoading = ref(false)
const pendingGrade = ref('')
const pendingClassId = ref(null)
const pendingSubject = ref('')
const gradeClasses = computed(() => classes.value.filter(c => c.grade === pendingGrade.value))
function openCreateModal() {
  createForm.value = { name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '', password: '' }
  createAssignments.value = []; pendingGrade.value = ''; pendingClassId.value = null; pendingSubject.value = ''
  showCreateModal.value = true
}
function addClassAssignment() {
  if (!pendingClassId.value || !pendingGrade.value) { toast.show('请选择年级和班级', 'info'); return }
  const cls = classes.value.find(c => c.id === pendingClassId.value)
  if (!cls) return
  if (createAssignments.value.some(a => a.class_id === cls.id)) { toast.show('该班级已添加', 'info'); return }
  createAssignments.value.push({ class_id: cls.id, class_name: cls.name, subject: pendingSubject.value || '默认科目' })
  pendingClassId.value = null; pendingSubject.value = ''
}
function removeClassAssignment(idx) { createAssignments.value.splice(idx, 1) }
async function submitCreate() {
  createLoading.value = true
  try {
    const payload: any = { ...createForm.value }
    if (createAssignments.value.length > 0 && createAssignments.value.every(a => a.class_id)) {
      payload.assignments = createAssignments.value.map(a => ({ class_id: a.class_id, role: 'subject_teacher', subject: a.subject || undefined }))
    }
    await apiPost('/api/v1/admin/teachers', payload)
    toast.show('教师创建成功', 'success')
    showCreateModal.value = false; await refreshTeachers()
  } catch (e: any) {
    const msg = e?.response?.data?.message || '参数错误'
    const errors = e?.response?.data?.errors
    if (errors) {
      const details = Object.values(errors).flat().join('；')
      toast.show(msg + '：' + details, 'error')
    } else {
      toast.show(msg, 'error')
    }
  } finally { createLoading.value = false }
}

const showEditModal = ref(false)
const editTarget = ref<Teacher | null>(null)
const editForm = ref({ name: '', nickname: '', subject: '', grade_team: '', phone: '', email: '' })
function openEditModal(t: Teacher) {
  editTarget.value = t
  editForm.value = { name: t.name, nickname: t.nickname || '', subject: t.subject || '', grade_team: t.grade_team || '', phone: t.phone || '', email: t.email || '' }
  showEditModal.value = true
}
async function submitEdit() {
  if (!editTarget.value) return
  try {
    await apiPut(`/api/v1/admin/teachers/${editTarget.value.id}`, editForm.value)
    toast.show('教师信息已更新', 'success')
    showEditModal.value = false; await refreshTeachers()
  } catch {}
}

const showAssignModal = ref(false)
const assignTarget = ref<Teacher | null>(null)
const assignList = ref<{ class_id: number | null; role: Role; subject: string; class_name?: string }[]>([])
const assignLoading = ref(false)
const assignGradeFilter = ref('')
const newAssignClassId = ref<number | null>(null)
const newAssignRole = ref<Role>('subject_teacher')
const newAssignSubject = ref('')
const filteredAssignClasses = computed(() => {
  if (!assignGradeFilter.value) return classes.value
  return classes.value.filter(c => c.grade === assignGradeFilter.value)
})
function openAssignModal(t: Teacher) {
  assignTarget.value = t
  assignList.value = t.assignments.length > 0
    ? t.assignments.map(a => ({ class_id: a.class_id, role: a.role as Role, subject: a.subject || '', class_name: a.class_name }))
    : []
  assignGradeFilter.value = ''; newAssignClassId.value = null; newAssignRole.value = 'subject_teacher'; newAssignSubject.value = ''
  showAssignModal.value = true
}
function addAssignRowNew() {
  if (!newAssignClassId.value) return
  const cls = classes.value.find(c => c.id === newAssignClassId.value)
  if (!cls) return
  if (assignList.value.some(a => a.class_id === cls.id)) { toast.show('该班级已添加', 'info'); return }
  assignList.value.push({ class_id: cls.id, class_name: shortClassName(cls.name), role: newAssignRole.value, subject: newAssignSubject.value || '默认科目' })
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

async function refreshTeachers() {
  loading.value = true
  try {
    const [tRes, cRes] = await Promise.all([
      apiGet<{ data: Teacher[] }>('/api/v1/admin/teachers'),
      apiGet<{ data: ClassRoom[] }>('/api/v1/admin/classes'),
    ])
    teachers.value = tRes.data || []; classes.value = cRes.data || []
  } finally { loading.value = false }
}
function downloadTemplate() { window.open('/api/v1/admin/teachers/template-csv', '_blank') }
const showResetPwdModal = ref(false)
const resetTarget = ref<Teacher | null>(null)
const resetPwdValue = ref('')
const resetPwdLoading = ref(false)

function openResetPwd(t: Teacher) {
  resetTarget.value = t
  resetPwdValue.value = ''
  showResetPwdModal.value = true
}

async function submitResetPwd() {
  if (!resetTarget.value) return
  resetPwdLoading.value = true
  try {
    await apiPost(`/api/v1/admin/teachers/${resetTarget.value.id}/reset-password`, { password: resetPwdValue.value || undefined })
    toast.show('密码已重置为：' + (resetPwdValue.value || '自动生成'), 'success')
    showResetPwdModal.value = false
  } catch { toast.show('重置失败', 'error') }
  finally { resetPwdLoading.value = false }
}
async function deleteTeacher(t: Teacher) {
  if (!confirm('确定删除教师「' + t.name + '」？')) return
  try { await apiDelete(`/api/v1/admin/teachers/${t.id}`); teachers.value = teachers.value.filter(x => x.id !== t.id); toast.show('已删除', 'success') } catch {}
}
function classById(id: number) { return classes.value.find(c => c.id === id) }
function shortClassName(name: string | undefined) {
  if (!name) return ''
  const m = name.match(/（(\d+)）班/)
  return m ? m[1] + '班' : name
}

onMounted(refreshTeachers)
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

    <div v-else class="card-grid">
      <div v-for="t in filteredTeachers" :key="t.id" class="teacher-card">
        <div class="card-header">
          <div class="avatar" :style="{ background: avatarGradient(t.name) }">{{ t.nickname?.[0] || t.name[0] }}</div>
          <div class="card-info">
            <div class="card-name">{{ t.name }}<span v-if="t.nickname && t.nickname !== t.name" class="card-nick">@{{ t.nickname }}</span></div>
            <div class="card-meta">
              <span v-if="t.subject" class="meta-tag subject">{{ t.subject }}</span>
              <span v-if="t.grade_team" class="meta-tag grade">{{ t.grade_team }}</span>
              <span class="meta-tag">{{ t.username }}</span>
              <span v-if="t.phone" class="meta-tag">{{ t.phone }}</span>
            </div>
            <div class="card-assign">
              <div style="display:flex;flex-wrap:wrap;gap:4px;margin-top:6px;">
                <span v-for="b in t.bindings" :key="b" class="binding-tag">{{ platformLabel(b) }}</span>
              </div>
              <div v-if="t.assignments.length > 0" class="assign-list">
                <div v-for="a in t.assignments" :key="a.class_id + '_' + a.role" class="assign-item" :style="{ borderLeftColor: roleColors[a.role as Role] }">
                  <span class="assign-role" :style="{ color: roleColors[a.role as Role] }">{{ roleLabel[a.role as Role] }}</span>
                  <span class="assign-class">{{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}</span>
                  <span v-if="a.subject" class="assign-subject">{{ a.subject }}</span>
                </div>
              </div>
              <div v-else style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">未分配班级</div>
            </div>
          </div>
        </div>
        <div class="card-actions">
          <button class="action-btn" @click="openEditModal(t)" title="编辑信息">&#x270F;&#xFE0F;</button>
          <button class="action-btn assign" @click="openAssignModal(t)" title="分配班级/角色">&#x1F3EB;</button>
          <button class="action-btn" @click="openResetPwd(t)" title="重置密码">&#x1F511;</button>
          <button class="action-btn danger" @click="deleteTeacher(t)" title="删除">&#x1F5D1;&#xFE0F;</button>
        </div>
      </div>
    </div>
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
                    <div style="flex:1;" class="form-group"><label>姓名 <span style="color:var(--color-danger);">*</span></label><input v-model="createForm.name" placeholder="姓名" class="form-input" /></div>
                    <div style="flex:1;" class="form-group"><label>昵称</label><input v-model="createForm.nickname" placeholder="默认拼音" class="form-input" /></div>
                  </div>
                  <div style="display:flex;gap:8px;">
                    <div style="flex:1;" class="form-group"><label>年级团队</label><select v-model="createForm.grade_team" class="form-input"><option value="">不指定</option><option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option></select></div>
                    <div style="flex:1;" class="form-group"><label>科目</label><select v-model="createForm.subject" class="form-input"><option value="">不指定</option><option v-for="s in subjects" :key="s" :value="s">{{ s }}</option></select></div>
                  </div>
                  <div style="display:flex;gap:8px;"><div style="flex:1;" class="form-group"><label>手机号</label><input v-model="createForm.phone" placeholder="选填" class="form-input" /></div><div style="flex:1;" class="form-group"><label>邮箱</label><input v-model="createForm.email" placeholder="选填" class="form-input" /></div></div>
                  <div class="form-group"><label>初始密码</label><input v-model="createForm.password" placeholder="留空自动生成" class="form-input" /></div>
                </div>
                <div>
                  <div style="font-size:12px;font-weight:600;color:var(--color-text);margin-bottom:8px;">📚 加入班级 <span style="font-size:10px;color:var(--color-text-secondary);">可选</span></div>
                  <div style="display:flex;gap:8px;align-items:flex-end;">
                    <div style="display:flex;flex-direction:column;gap:4px;flex:1;">
                      <div class="form-group" style="margin-bottom:0;"><label>年级</label><select v-model="pendingGrade" class="form-input"><option value="">请选择</option><option v-for="g in grades" :key="g" :value="g">{{ g }}</option></select></div>
                      <div class="form-group" style="margin-bottom:0;"><label>班级</label><select v-model="pendingClassId" :disabled="!pendingGrade" class="form-input"><option :value="null">请选择</option><option v-for="c in gradeClasses" :key="c.id" :value="c.id">{{ shortClassName(c.name) }}</option></select></div>
                      <div class="form-group" style="margin-bottom:0;"><label>科目</label><select v-model="pendingSubject" class="form-input"><option value="">请选择</option><option v-for="s in subjects" :key="s" :value="s">{{ s }}</option></select></div>
                    </div>
                    <button @click="addClassAssignment" :disabled="!pendingClassId" style="padding:6px 16px;border-radius:8px;border:1px solid var(--color-accent);background:rgba(79,70,229,0.08);color:var(--color-accent);font-size:13px;cursor:pointer;font-weight:500;white-space:nowrap;height:36px;">➕ 添加</button>
                  </div>
                  <div style="font-size:11px;color:var(--color-text-secondary);margin-top:8px;padding:6px 8px;background:var(--color-bg);border-radius:6px;border-left:2px solid var(--color-accent);">💡 创建后可随时在列表中点击 🏫 按钮重新分配班级</div>
                  <div v-if="createAssignments.length > 0" style="margin-top:8px;display:flex;flex-direction:column;gap:4px;">
                    <div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:2px;">📋 已添加（{{ createAssignments.length }}）</div>
                    <div v-for="(a, i) in createAssignments" :key="i" style="display:flex;align-items:center;gap:6px;padding:4px 8px;background:var(--color-bg);border-radius:4px;font-size:12px;">
                      <span style="flex:1;color:var(--color-text);">{{ shortClassName(a.class_name) }} · {{ a.subject }}</span>
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
  <div v-if="showEditModal" class="modal-overlay">
      <div class="modal-panel" style="max-width:520px;">
        <div class="modal-header"><h3>编辑「{{ editTarget?.name }}」</h3><button class="close-btn" @click="showEditModal = false">&times;</button></div>
        <div class="modal-body">
          <div class="form-row">
            <div class="form-group flex-2"><label>姓名</label><input v-model="editForm.name" class="form-input" /></div>
            <div class="form-group flex-1"><label>昵称</label><input v-model="editForm.nickname" class="form-input" /></div>
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label>年级团队</label>
              <select v-model="editForm.grade_team" class="form-input">
                <option value="">不指定</option>
                <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
              </select>
            </div>
            <div class="form-group flex-1">
              <label>科目</label>
              <select v-model="editForm.subject" class="form-input">
                <option value="">不指定</option>
                <option v-for="s in subjects" :key="s" :value="s">{{ s }}</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group flex-1"><label>手机号</label><input v-model="editForm.phone" class="form-input" /></div>
            <div class="form-group flex-1"><label>邮箱</label><input v-model="editForm.email" class="form-input" /></div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);" @click="showEditModal = false">取消</button>
          <button class="btn btn-primary" @click="submitEdit">保存</button>
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
          <div style="flex:1;min-width:100px;" class="form-group"><label>角色</label><select v-model="newAssignRole" class="form-input"><option v-for="(label, role) in roleLabel" :key="role" :value="role">{{ label }}</option></select></div>
          <div style="flex:1;min-width:90px;" class="form-group"><label>科目</label><select v-model="newAssignSubject" class="form-input"><option value="">默认</option><option v-for="s in subjects" :key="s" :value="s">{{ s }}</option></select></div>
          <button @click="addAssignRowNew" :disabled="!newAssignClassId" style="padding:8px 16px;border-radius:8px;border:1px solid var(--color-accent);background:rgba(79,70,229,0.08);color:var(--color-accent);font-size:13px;cursor:pointer;font-weight:500;white-space:nowrap;height:36px;flex-shrink:0;">➕ 添加</button>
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

  <!-- 重置密码弹窗 -->
  <Teleport to="body">
    <div v-if="showResetPwdModal" @click="showResetPwdModal = false" style="position:fixed;inset:0;z-index:1000;background:rgba(0,0,0,0.15);display:flex;align-items:center;justify-content:center;padding:20px;">
      <div @click.stop style="background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:16px;padding:24px;max-width:400px;width:100%;box-shadow:0 8px 32px rgba(0,0,0,0.12);">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
          <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">🔑 重置密码 — {{ resetTarget?.name }}</h3>
          <button @click="showResetPwdModal = false" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">✕</button>
        </div>
        <div class="form-group"><label>新密码</label><input v-model="resetPwdValue" type="text" class="form-input" placeholder="留空自动生成"></div>
        <div style="font-size:12px;color:var(--color-text-secondary);margin-bottom:16px;padding:8px;background:var(--color-bg);border-radius:6px;">💡 留空将自动生成随机密码</div>
        <div style="display:flex;gap:12px;">
          <button @click="showResetPwdModal = false" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);">取消</button>
          <button @click="submitResetPwd" :disabled="resetPwdLoading" style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;">{{ resetPwdLoading ? '重置中...' : '确认重置' }}</button>
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
.card-grid { display:grid;grid-template-columns:repeat(auto-fill,minmax(360px,1fr));gap:14px; }
.teacher-card { background:var(--color-bg-card);border:1px solid var(--color-border);border-radius:14px;padding:16px 20px;display:flex;align-items:center;justify-content:space-between;gap:14px;transition:box-shadow 0.2s,border-color 0.2s; }
.teacher-card:hover { border-color:#e5e7eb;box-shadow:0 2px 12px rgba(0,0,0,0.04); }
.card-header { display:flex;align-items:center;gap:14px;flex:1;min-width:0; }
.avatar { width:46px;height:46px;border-radius:12px;display:flex;align-items:center;justify-content:center;color:white;font-weight:700;font-size:17px;flex-shrink:0; }
.card-info { flex:1;min-width:0; }
.card-name { font-weight:600;font-size:15px; }
.card-nick { color:#9ca3af;font-size:12px;margin-left:6px; }
.card-meta { display:flex;gap:6px;flex-wrap:wrap;margin-top:4px; }
.meta-tag { font-size:11px;padding:1px 8px;border-radius:8px;font-weight:500; }
.meta-tag.subject { background:#ecfdf5;color:#059669; }
.meta-tag.grade { background:#eff6ff;color:#2563eb; }
.meta-tag:not(.subject):not(.grade) { color:#6b7280; }
.card-assign { margin-top:6px; }
.binding-tag { display:inline-block;font-size:11px;padding:1px 8px;border-radius:6px;font-weight:500;background:rgba(79,70,229,0.06);color:#7c3aed;border:1px solid rgba(79,70,229,0.15); }
.assign-list { display:flex;flex-direction:column;gap:3px;margin-top:4px; }
.assign-item { display:flex;align-items:center;gap:6px;padding:4px 8px;background:var(--color-bg);border-radius:6px;border-left:3px solid;font-size:12px; }
.assign-role { font-weight:600;font-size:11px;flex-shrink:0; }
.assign-class { color:var(--color-text);flex:1; }
.assign-subject { color:var(--color-text-secondary);font-size:11px; }
.card-actions { display:flex;gap:4px;flex-shrink:0; }
.action-btn { width:34px;height:34px;border-radius:8px;border:1px solid var(--color-border);background:var(--color-bg-card);cursor:pointer;font-size:16px;display:flex;align-items:center;justify-content:center;transition:all 0.15s; }
.action-btn:hover { background:#f9fafb;border-color:#e5e7eb; }
.action-btn.assign:hover { background:#ede9fe;border-color:#c4b5fd; }
.action-btn.danger:hover { background:#fee2e2;border-color:#fecaca; }
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
