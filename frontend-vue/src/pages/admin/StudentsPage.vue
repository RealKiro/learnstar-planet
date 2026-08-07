<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import ModalGlass from '@/components/common/ModalGlass.vue'
import { avatarGradient } from '@/utils/constants'
import type { ApiResponse, Student, ClassRoom } from '@/types'
import PetSprite from '@/components/pet/PetSprite.vue'

const students = ref<Student[]>([])
const classes = ref<ClassRoom[]>([])
const loading = ref(true)
const loadError = ref('')
const filterClassId = ref<number | ''>('')
const searchKeyword = ref('')
const selectedIds = ref<number[]>([])
const meta = ref({ current_page: 1, last_page: 1, per_page: 50, total: 0 })

const pageLabel = computed(() => {
  const { current_page, last_page, total } = meta.value
  return last_page > 1 ? `第 ${current_page} / ${last_page} 页 · 共 ${total} 人` : `共 ${total} 人`
})

const gradeOptions = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级']

// 新增/编辑弹窗
const showEditModal = ref(false)
const isEditing = ref(false)
const submitStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const deleteStatusMap = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const batchDeleteStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const moveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const editingId = ref<number | null>(null)
const form = ref({ name: '', class_id: '' as number | '', gender: '男', student_no: '' })
// 级联选择
const formGrade = ref('一年级')
const formClassId = ref<number | ''>('')

// 批量调班弹窗
const showMoveModal = ref(false)
const targetGrade = ref('一年级')
const targetClassId = ref<number | ''>('')
const formError = ref('')
const moveError = ref('')

const allSelected = computed(() =>
  students.value.length > 0 && selectedIds.value.length === students.value.length,
)

// 按年级分组班级（用于级联下拉）
const classesByGrade = computed(() => {
  const map = new Map<string, ClassRoom[]>()
  for (const g of gradeOptions) map.set(g, [])
  for (const c of classes.value) {
    const grade = c.grade || '未分年级'
    if (!map.has(grade)) map.set(grade, [])
    map.get(grade)!.push(c)
  }
  return map
})

const formClassOptions = computed(() => classesByGrade.value.get(formGrade.value) || [])

const targetClassOptions = computed(() => classesByGrade.value.get(targetGrade.value) || [])

onMounted(async () => {
  await Promise.all([loadClasses(), loadStudents()])
})

async function loadClasses() {
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
  } catch { classes.value = [] }
}

async function loadStudents(resetPage = false) {
  if (resetPage) meta.value.current_page = 1
  loading.value = true
  selectedIds.value = []
  try {
    const params = new URLSearchParams()
    params.set('per_page', String(meta.value.per_page))
    if (meta.value.current_page > 1) params.set('page', String(meta.value.current_page))
    if (filterClassId.value !== '') params.set('class_id', String(filterClassId.value))
    if (searchKeyword.value.trim()) params.set('search', searchKeyword.value.trim())
    const res = await apiGet<ApiResponse<Student[]>>(`/api/v1/admin/students?${params.toString()}`)
    students.value = res.data || []
    loadError.value = ''
    if (res.meta) meta.value = res.meta
  } catch {
    students.value = []
    loadError.value = '学生列表加载失败'
  }
  finally { loading.value = false }
}

function changePage(page: number) {
  const { current_page, last_page } = meta.value
  if (page < 1 || page > last_page || page === current_page) return
  meta.value.current_page = page
  loadStudents()
}

function getClassFromId(id: number | '') {
  return classes.value.find(c => c.id === id)
}

function openAddModal() {
  isEditing.value = false
  editingId.value = null
  formGrade.value = '一年级'
  formClassId.value = ''
  form.value = {
    name: '',
    class_id: '' as number | '',
    gender: '男',
    student_no: '',
  }
  showEditModal.value = true
}

function openEditModal(s: Student) {
  isEditing.value = true
  editingId.value = s.id
  const cls = classes.value.find(c => c.id === s.class_id)
  if (cls) {
    formGrade.value = cls.grade || '一年级'
    formClassId.value = s.class_id
  } else {
    formGrade.value = '一年级'
    formClassId.value = ''
  }
  form.value = {
    name: s.name,
    class_id: s.class_id,
    gender: s.gender || '男',
    student_no: s.student_no || '',
  }
  showEditModal.value = true
}

function onFormClassChange(clasId: number | '') {
  form.value.class_id = clasId
}

function onTargetClassChange(clasId: number | '') {
  targetClassId.value = clasId
}

function onFormGradeChange() {
  formClassId.value = ''
  form.value.class_id = ''
}

function onTargetGradeChange() {
  targetClassId.value = ''
}

async function submitForm() {
  if (!form.value.name.trim()) { formError.value = '请填写学生姓名'; return }
  if (!formClassId.value) { formError.value = '请选择班级'; return }
  form.value.class_id = formClassId.value as number
  submitStatus.value = 'loading'
  const payload = {
    name: form.value.name.trim(),
    class_id: form.value.class_id as number,
    gender: form.value.gender,
    student_no: form.value.student_no.trim(),
  }
  try {
    if (isEditing.value && editingId.value !== null) {
      await apiPut(`/api/v1/admin/students/${editingId.value}`, payload)
    } else {
      await apiPost('/api/v1/admin/students', payload)
    }
    submitStatus.value = 'success'
    showEditModal.value = false
    await loadStudents(!isEditing.value) // 新建回到第 1 页看到新学生
    setTimeout(() => { submitStatus.value = 'idle' }, 1500)
  } catch {
    submitStatus.value = 'error'
    setTimeout(() => { submitStatus.value = 'idle' }, 3000)
  }
}

async function deleteStudent(s: Student) {
  const ok = await openConfirm({ title: '删除学生', message: `确定删除学生「${s.name}」？`, danger: true, confirmText: '确认删除' })
  if (!ok) return
  deleteStatusMap.value[s.id] = 'loading'
  try {
    await apiDelete(`/api/v1/admin/students/${s.id}`)
    deleteStatusMap.value[s.id] = 'success'
    students.value = students.value.filter(x => x.id !== s.id)
    selectedIds.value = selectedIds.value.filter(id => id !== s.id)
    setTimeout(() => { delete deleteStatusMap.value[s.id] }, 1500)
  } catch {
    deleteStatusMap.value[s.id] = 'error'
    setTimeout(() => { delete deleteStatusMap.value[s.id] }, 3000)
  }
}

function toggleSelect(id: number) {
  const i = selectedIds.value.indexOf(id)
  if (i >= 0) selectedIds.value.splice(i, 1)
  else selectedIds.value.push(id)
}

function toggleSelectAll() {
  if (allSelected.value) selectedIds.value = []
  else selectedIds.value = students.value.map(s => s.id)
}

async function batchDelete() {
  const ok = await openConfirm({ title: '批量删除学生', message: `确定批量删除 ${selectedIds.value.length} 名学生？`, danger: true, confirmText: '确认删除' })
  if (!ok) return
  batchDeleteStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/students/batch-delete', { student_ids: selectedIds.value })
    batchDeleteStatus.value = 'success'
    selectedIds.value = []
    await loadStudents()
    setTimeout(() => { batchDeleteStatus.value = 'idle' }, 1500)
  } catch {
    batchDeleteStatus.value = 'error'
    setTimeout(() => { batchDeleteStatus.value = 'idle' }, 3000)
  }
}

function openMoveModal() {
  targetGrade.value = '一年级'
  targetClassId.value = ''
  showMoveModal.value = true
}

async function submitMove() {
  if (targetClassId.value === '') { moveError.value = '请选择目标班级'; return }
  moveStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/students/batch-move', {
      student_ids: selectedIds.value,
      target_class_id: targetClassId.value,
    })
    moveStatus.value = 'success'
    showMoveModal.value = false
    selectedIds.value = []
    await loadStudents()
    setTimeout(() => { moveStatus.value = 'idle' }, 1500)
  } catch {
    moveStatus.value = 'error'
    setTimeout(() => { moveStatus.value = 'idle' }, 3000)
  }
}
</script>

<template>
  <div>
    <div class="page-header">
      <div>
        <p class="page-subtitle">学生管理</p>
        <h2 class="page-title">学生列表</h2>
      </div>
      <div class="header-actions">
        <input
          v-model="searchKeyword"
          class="form-input search-input"
          placeholder="🔍 搜索姓名 / 学号"
          @keyup.enter="loadStudents(true)"
          @keyup.esc="searchKeyword = ''"
        >
        <select v-model="filterClassId" class="form-select filter-select" @change="loadStudents(true)">
          <option value="">全部班级</option>
          <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
        <button class="btn btn-sm btn-primary" @click="openAddModal">+ 添加学生</button>
      </div>
    </div>

    <div v-if="selectedIds.length > 0" class="card batch-bar">
      <span class="batch-bar__count">已选择 <b>{{ selectedIds.length }}</b> 名学生</span>
      <div class="batch-bar__actions">
        <button class="btn btn-sm btn-ghost-card" @click="openMoveModal">🔁 批量调班</button>
        <button class="btn btn-sm" :class="batchDeleteStatus === 'loading' ? 'btn-state-loading' : batchDeleteStatus === 'success' ? 'btn-state-success' : batchDeleteStatus === 'error' ? 'btn-state-error' : 'btn-danger'" :disabled="batchDeleteStatus === 'loading'" @click="batchDelete">{{ batchDeleteStatus === 'loading' ? '删除中...' : batchDeleteStatus === 'success' ? '已删除' : batchDeleteStatus === 'error' ? '删除失败' : '🗑 批量删除' }}</button>
      </div>
    </div>

    <div v-if="loading" class="loading-state"><div class="loading-spinner"></div><p>加载中...</p></div>

    <div v-else-if="loadError" class="error-state">
      <div class="error-state__icon">⚠️</div>
      <p class="error-state__msg">{{ loadError }}</p>
      <button class="btn btn-sm btn-primary" @click="loadStudents(true)">重试</button>
    </div>

    <div v-else-if="students.length === 0" class="card empty-state">
      <div class="empty-state__icon">👨‍🎓</div>
      <p>暂无学生记录</p>
    </div>

    <div v-else class="data-table">
      <table>
        <thead>
          <tr>
            <th class="check-col"><input type="checkbox" :checked="allSelected" @change="toggleSelectAll"></th>
            <th>姓名</th><th>学号</th><th>宠物</th><th>性别</th><th>班级</th><th>年级</th><th>积分</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in students" :key="s.id">
            <td><input type="checkbox" :checked="selectedIds.includes(s.id)" @change="toggleSelect(s.id)"></td>
            <td>
              <div class="student-cell">
                <div class="avatar-badge" :style="{ background: avatarGradient(s.name) }">{{ s.name[0] }}</div>
                <span class="student-name">{{ s.name }}</span>
              </div>
            </td>
            <td class="mono-text">{{ s.student_no || '-' }}</td>
            <td>
              <div class="pet-cell">
                <div class="pet-avatar">
                  <PetSprite v-if="(s as any).pet_species" :species-id="(s as any).pet_species" :level="(s as any).pet_level || 1" :animate="true" />
                  <span v-else class="pet-egg">🥚</span>
                </div>
                <span v-if="(s as any).pet_level" class="pet-level">Lv.{{ (s as any).pet_level }}</span>
              </div>
            </td>
            <td>{{ s.gender || '-' }}</td>
            <td>{{ s.class_name || getClassFromId(s.class_id)?.name || '-' }}</td>
            <td><span v-if="s.class_grade || getClassFromId(s.class_id)?.grade" class="grade-badge">{{ s.class_grade || getClassFromId(s.class_id)?.grade }}</span><span v-else>-</span></td>
            <td class="score-text">{{ s.total_score }}</td>
            <td>
              <div class="cell-actions">
                <button class="btn btn-sm btn-mini" @click="openEditModal(s)">编辑</button>
                <button class="btn btn-sm" :class="deleteStatusMap[s.id] === 'loading' ? 'btn-state-loading' : deleteStatusMap[s.id] === 'success' ? 'btn-state-success' : deleteStatusMap[s.id] === 'error' ? 'btn-state-error' : 'btn-danger'" :disabled="deleteStatusMap[s.id] === 'loading'" @click="deleteStudent(s)">{{ deleteStatusMap[s.id] === 'loading' ? '删除中...' : deleteStatusMap[s.id] === 'success' ? '已删除' : deleteStatusMap[s.id] === 'error' ? '删除失败' : '删除' }}</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="meta.last_page > 1" class="pagination">
        <button class="btn btn-sm btn-ghost-card" :disabled="meta.current_page <= 1" @click="changePage(meta.current_page - 1)">← 上一页</button>
        <span class="pagination__info">{{ pageLabel }}</span>
        <button class="btn btn-sm btn-ghost-card" :disabled="meta.current_page >= meta.last_page" @click="changePage(meta.current_page + 1)">下一页 →</button>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <ModalGlass :visible="showEditModal" @update:visible="showEditModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">{{ isEditing ? '编辑学生' : '添加学生' }}</h3>
        <button class="modal-close" @click="showEditModal = false">×</button>
      </div>
      <div class="form-group">
        <label>姓名</label>
        <input v-model="form.name" class="form-input" placeholder="如：张小明" @keydown.enter="submitForm">
      </div>
      <div class="form-group">
        <label>年级</label>
        <select v-model="formGrade" class="form-select" @change="onFormGradeChange">
          <option v-for="g in gradeOptions" :key="g" :value="g">{{ g }}</option>
        </select>
      </div>
      <div class="form-group">
        <label>班级</label>
        <select v-model="formClassId" class="form-select" @change="onFormClassChange(formClassId as number)">
          <option value="" disabled>请选择班级</option>
          <option v-for="c in formClassOptions" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
      </div>
      <div class="form-group">
        <label>性别</label>
        <select v-model="form.gender" class="form-select">
          <option value="男">男</option>
          <option value="女">女</option>
        </select>
      </div>
      <div class="form-group">
        <label>学号</label>
        <input v-model="form.student_no" class="form-input" placeholder="可选" @keydown.enter="submitForm">
      </div>
      <div class="modal-actions">
        <div v-if="formError" class="error-banner">{{ formError }}</div>
        <button class="btn btn-sm btn-ghost-card" @click="showEditModal = false">取消</button>
        <button class="btn btn-sm" :class="submitStatus === 'loading' ? 'btn-state-loading' : submitStatus === 'success' ? 'btn-state-success' : submitStatus === 'error' ? 'btn-state-error' : 'btn-primary'" :disabled="submitStatus === 'loading'" @click="submitForm">{{ submitStatus === 'loading' ? '保存中...' : submitStatus === 'success' ? '已保存' : submitStatus === 'error' ? '保存失败' : '保存' }}</button>
      </div>
    </ModalGlass>

    <!-- 批量调班弹窗 -->
    <ModalGlass :visible="showMoveModal" @update:visible="showMoveModal = $event">
      <div class="modal-header">
        <h3 class="modal-title">批量调班</h3>
        <button class="modal-close" @click="showMoveModal = false">×</button>
      </div>
      <p class="move-hint">将选中的 <b>{{ selectedIds.length }}</b> 名学生调入以下班级：</p>
      <div class="form-group">
        <label>年级</label>
        <select v-model="targetGrade" class="form-select" @change="onTargetGradeChange">
          <option v-for="g in gradeOptions" :key="g" :value="g">{{ g }}</option>
        </select>
      </div>
      <div class="form-group">
        <label>目标班级</label>
        <select v-model="targetClassId" class="form-select" @change="onTargetClassChange(targetClassId as number)">
          <option value="" disabled>请选择目标班级</option>
          <option v-for="c in targetClassOptions" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
      </div>
      <div class="modal-actions">
        <div v-if="moveError" class="error-banner">{{ moveError }}</div>
        <button class="btn btn-sm btn-ghost-card" @click="showMoveModal = false">取消</button>
        <button class="btn btn-sm" :class="moveStatus === 'loading' ? 'btn-state-loading' : moveStatus === 'success' ? 'btn-state-success' : moveStatus === 'error' ? 'btn-state-error' : 'btn-primary'" :disabled="moveStatus === 'loading'" @click="submitMove">{{ moveStatus === 'loading' ? '调班中...' : moveStatus === 'success' ? '调班成功' : moveStatus === 'error' ? '调班失败' : '确认调班' }}</button>
      </div>
    </ModalGlass>
  </div>
</template>

<style scoped>
/* ===== 页头 ===== */
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; flex-wrap: wrap; gap: 12px; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); margin-bottom: 4px; }
.page-title { font-size: 24px; font-weight: 700; }
.header-actions { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
.search-input { width: 200px; }
.filter-select { width: 180px; }

/* ===== 分页 ===== */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 14px;
  border-top: 1px solid var(--color-border);
}
.pagination__info { font-size: 13px; color: var(--color-text-secondary); }
.pagination .btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* ===== 批量操作栏 ===== */
.batch-bar { padding: 12px 20px; margin-bottom: 16px; display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; }
.batch-bar__count { font-size: 13px; color: var(--color-text-secondary); }
.batch-bar__count b { color: var(--color-primary); }
.batch-bar__actions { display: flex; gap: 8px; }

/* ===== 表格单元格 ===== */
.check-col { width: 40px; }
.student-cell { display: flex; align-items: center; gap: 10px; }
.avatar-badge {
  width: 32px; height: 32px; border-radius: 10px; color: #fff;
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 13px; flex-shrink: 0;
}
.student-name { font-weight: 600; }
.mono-text { font-family: monospace; color: var(--color-text-secondary); }
.grade-badge {
  display: inline-block; padding: 2px 10px; border-radius: 20px;
  font-size: 12px; font-weight: 600;
  background: rgba(79,70,229,0.08); color: var(--color-primary);
}
.pet-cell { display: flex; align-items: center; gap: 6px; }
.pet-avatar { width: 30px; height: 30px; flex-shrink: 0; }
.pet-egg { font-size: 16px; }
.pet-level { font-size: 11px; color: var(--color-text-secondary); }
.score-text { font-weight: 600; color: var(--color-accent); }
.cell-actions { display: flex; gap: 4px; }

/* ===== 次要按钮（btn-ghost-card 已全局化） ===== */
.btn-mini {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  font-size: 12px;
}

/* ===== 批量调班弹窗 ===== */
.move-hint { font-size: 13px; color: var(--color-text-secondary); margin-bottom: 12px; }
.move-hint b { color: var(--color-primary); }
</style>
