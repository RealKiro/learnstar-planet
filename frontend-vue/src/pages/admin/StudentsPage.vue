<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { avatarGradient } from '@/utils/constants'
import type { ApiResponse, Student, ClassRoom } from '@/types'

const toast = useToastStore()

const students = ref<Student[]>([])
const classes = ref<ClassRoom[]>([])
const loading = ref(true)
const filterClassId = ref<number | ''>('')
const selectedIds = ref<number[]>([])

// 新增/编辑弹窗
const showEditModal = ref(false)
const isEditing = ref(false)
const modalLoading = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', class_id: '' as number | '', gender: '男', student_no: '' })

// 批量调班弹窗
const showMoveModal = ref(false)
const targetClassId = ref<number | ''>('')

const allSelected = computed(() =>
  students.value.length > 0 && selectedIds.value.length === students.value.length,
)

onMounted(async () => {
  await Promise.all([loadClasses(), loadStudents()])
})

async function loadClasses() {
  try {
    const res = await apiGet<ApiResponse<ClassRoom[]>>('/api/v1/admin/classes')
    classes.value = res.data || []
  } catch { classes.value = [] }
}

async function loadStudents() {
  loading.value = true
  selectedIds.value = []
  try {
    const url = filterClassId.value !== ''
      ? `/api/v1/admin/students?class_id=${filterClassId.value}&per_page=100`
      : '/api/v1/admin/students?per_page=100'
    const res = await apiGet<ApiResponse<Student[]>>(url)
    students.value = res.data || []
  } catch { students.value = [] }
  finally { loading.value = false }
}

function openAddModal() {
  isEditing.value = false
  editingId.value = null
  form.value = {
    name: '',
    class_id: filterClassId.value !== '' ? filterClassId.value : (classes.value[0]?.id ?? ''),
    gender: '男',
    student_no: '',
  }
  showEditModal.value = true
}

function openEditModal(s: Student) {
  isEditing.value = true
  editingId.value = s.id
  form.value = {
    name: s.name,
    class_id: s.class_id,
    gender: s.gender || '男',
    student_no: s.student_no || '',
  }
  showEditModal.value = true
}

async function submitForm() {
  if (!form.value.name.trim()) { toast.show('请填写学生姓名', 'error'); return }
  if (form.value.class_id === '') { toast.show('请选择班级', 'error'); return }
  modalLoading.value = true
  const payload = {
    name: form.value.name.trim(),
    class_id: form.value.class_id as number,
    gender: form.value.gender,
    student_no: form.value.student_no.trim(),
  }
  try {
    if (isEditing.value && editingId.value !== null) {
      await apiPut(`/api/v1/admin/students/${editingId.value}`, payload)
      toast.show('已更新学生：' + payload.name, 'success')
    } else {
      await apiPost('/api/v1/admin/students', payload)
      toast.show('已添加学生：' + payload.name, 'success')
    }
    showEditModal.value = false
    await loadStudents()
  } catch { /* handled */ }
  finally { modalLoading.value = false }
}

async function deleteStudent(s: Student) {
  if (!confirm(`确定删除学生「${s.name}」？`)) return
  try {
    await apiDelete(`/api/v1/admin/students/${s.id}`)
    students.value = students.value.filter(x => x.id !== s.id)
    selectedIds.value = selectedIds.value.filter(id => id !== s.id)
    toast.show('已删除学生：' + s.name, 'success')
  } catch { /* handled */ }
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
  if (selectedIds.value.length === 0) { toast.show('请先选择学生', 'error'); return }
  if (!confirm(`确定批量删除 ${selectedIds.value.length} 名学生？`)) return
  try {
    await apiPost('/api/v1/admin/students/batch-delete', { student_ids: selectedIds.value })
    toast.show(`已删除 ${selectedIds.value.length} 名学生`, 'success')
    selectedIds.value = []
    await loadStudents()
  } catch { /* handled */ }
}

function openMoveModal() {
  if (selectedIds.value.length === 0) { toast.show('请先选择学生', 'error'); return }
  targetClassId.value = ''
  showMoveModal.value = true
}

async function submitMove() {
  if (targetClassId.value === '') { toast.show('请选择目标班级', 'error'); return }
  modalLoading.value = true
  try {
    await apiPost('/api/v1/admin/students/batch-move', {
      student_ids: selectedIds.value,
      target_class_id: targetClassId.value,
    })
    toast.show(`已调班 ${selectedIds.value.length} 名学生`, 'success')
    showMoveModal.value = false
    selectedIds.value = []
    await loadStudents()
  } catch { /* handled */ }
  finally { modalLoading.value = false }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;flex-wrap:wrap;gap:12px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">学生管理</p>
        <h2 style="font-size:24px;font-weight:700;">学生列表</h2>
      </div>
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
        <select v-model="filterClassId" class="form-select" style="width:180px;" @change="loadStudents">
          <option value="">全部班级</option>
          <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
        <button class="btn btn-sm btn-primary" @click="openAddModal">+ 添加学生</button>
      </div>
    </div>

    <div v-if="selectedIds.length > 0" class="card" style="padding:12px 20px;margin-bottom:16px;display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap;">
      <span style="font-size:13px;color:var(--color-text-secondary);">已选择 <b style="color:var(--color-primary);">{{ selectedIds.length }}</b> 名学生</span>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="openMoveModal">🔁 批量调班</button>
        <button class="btn btn-sm" style="background:#fee2e2;color:#dc2626;border:1px solid #fecaca;" @click="batchDelete">🗑 批量删除</button>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="students.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">👨‍🎓</div>
      <p>暂无学生记录</p>
    </div>

    <div v-else class="data-table">
      <table>
        <thead>
          <tr>
            <th style="width:40px;"><input type="checkbox" :checked="allSelected" @change="toggleSelectAll"></th>
            <th>姓名</th><th>学号</th><th>性别</th><th>班级</th><th>年级</th><th>积分</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in students" :key="s.id">
            <td><input type="checkbox" :checked="selectedIds.includes(s.id)" @change="toggleSelect(s.id)"></td>
            <td>
              <div style="display:flex;align-items:center;gap:10px;">
                <div :style="{ width:'32px', height:'32px', borderRadius:'10px', background: avatarGradient(s.name), color:'white', display:'flex', alignItems:'center', justifyContent:'center', fontWeight:700, fontSize:'13px', flexShrink:0 }">{{ s.name[0] }}</div>
                <span style="font-weight:600;">{{ s.name }}</span>
              </div>
            </td>
            <td style="font-family:monospace;color:var(--color-text-secondary);">{{ s.student_no || '-' }}</td>
            <td>{{ s.gender || '-' }}</td>
            <td>{{ s.class_name || '-' }}</td>
            <td><span v-if="s.class_grade" style="display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;font-weight:600;background:rgba(79,70,229,0.08);color:var(--color-primary);">{{ s.class_grade }}</span><span v-else>-</span></td>
            <td style="font-weight:600;color:var(--color-accent);">{{ s.total_score }}</td>
            <td>
              <div style="display:flex;gap:4px;">
                <button class="btn btn-sm" style="background:var(--color-bg);color:var(--color-text-secondary);border:1px solid var(--color-border);" @click="openEditModal(s)">编辑</button>
                <button class="btn btn-sm" style="background:#fee2e2;color:#dc2626;border:1px solid #fecaca;" @click="deleteStudent(s)">删除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showEditModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.self="showEditModal = false">
      <div class="card" style="width:90%;max-width:420px;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">{{ isEditing ? '编辑学生' : '添加学生' }}</h3>
          <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showEditModal = false">×</button>
        </div>
        <div class="form-group">
          <label>姓名</label>
          <input v-model="form.name" class="form-input" placeholder="如：张小明" @keydown.enter="submitForm">
        </div>
        <div class="form-group">
          <label>班级</label>
          <select v-model="form.class_id" class="form-select">
            <option value="" disabled>请选择班级</option>
            <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
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
        <div style="display:flex;gap:8px;justify-content:flex-end;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showEditModal = false">取消</button>
          <button class="btn btn-sm btn-primary" :disabled="modalLoading" @click="submitForm">{{ modalLoading ? '保存中...' : '保存' }}</button>
        </div>
      </div>
    </div>

    <!-- 批量调班弹窗 -->
    <div v-if="showMoveModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000;" @click.self="showMoveModal = false">
      <div class="card" style="width:90%;max-width:420px;padding:32px;">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
          <h3 style="font-size:18px;font-weight:700;">批量调班</h3>
          <button style="background:none;border:none;font-size:20px;cursor:pointer;color:var(--color-text-secondary);" @click="showMoveModal = false">×</button>
        </div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:12px;">将选中的 <b style="color:var(--color-primary);">{{ selectedIds.length }}</b> 名学生调入以下班级：</p>
        <div class="form-group">
          <label>目标班级</label>
          <select v-model="targetClassId" class="form-select">
            <option value="" disabled>请选择目标班级</option>
            <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;">
          <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showMoveModal = false">取消</button>
          <button class="btn btn-sm btn-primary" :disabled="modalLoading" @click="submitMove">{{ modalLoading ? '调班中...' : '确认调班' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
