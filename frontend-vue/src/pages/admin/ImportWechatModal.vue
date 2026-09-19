<script setup lang="ts">
import { ref, computed } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import ModalGlass from '@/components/common/ModalGlass.vue'

const props = defineProps<{
  visible: boolean
}>()
const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'imported'): void
}>()

interface WecomMember {
  userid: string
  name: string
  mobile: string
  email: string
  department_names: string[]
}
interface ClassItem { id: number; name: string }
interface Row {
  userid: string
  name: string
  mobile: string
  email: string
  department_names: string[]
  role: 'teacher' | 'student'
  class_id: number | ''
  selected: boolean
}
interface SkippedEntry { name: string; reason: string }
interface TeacherAccount { name: string; username: string; initial_password: string }
interface ImportResult {
  created_teachers: number
  created_students: number
  skipped_teachers: SkippedEntry[]
  skipped_students: SkippedEntry[]
  teacher_accounts: TeacherAccount[]
}

const PLATFORM_LABELS: Record<string, string> = {
  wechat_work: '企业微信', dingtalk: '钉钉', feishu: '飞书',
}

const loading = ref(false)
const importing = ref(false)
const loadError = ref('')
const importError = ref('')
const rows = ref<Row[]>([])
const classes = ref<ClassItem[]>([])
const platform = ref('')
const search = ref('')
const result = ref<ImportResult | null>(null)
const copied = ref(false)

const close = () => emit('update:visible', false)

// 部门名 → 班级名 匹配：先精确，再归一化（去掉括号/空格/连字符，如"六年级1班"↔"六年级（1）班"）
const normalizeName = (s: string) => s.replace(/[（()）\s\-—_·、]/g, '').toLowerCase()

function matchClassByDept(deptNames: string[], clsList: ClassItem[]): number | '' {
  for (const n of deptNames) {
    const hit = clsList.find(c => c.name === n)
    if (hit) return hit.id
  }
  for (const n of deptNames) {
    const nn = normalizeName(n)
    if (!nn) continue
    const hit = clsList.find(c => normalizeName(c.name) === nn)
    if (hit) return hit.id
  }
  return ''
}

async function loadContacts() {
  loading.value = true
  loadError.value = ''
  try {
    const cRes = await apiGet<{ data: { members: WecomMember[]; platform?: string } }>('/api/v1/admin/third-party/contacts', { skipToast: true })
    const clsRes = await apiGet<{ data: ClassItem[] }>('/api/v1/admin/classes', { skipToast: true })
    classes.value = clsRes.data || []
    platform.value = cRes.data?.platform || ''
    const members = cRes.data?.members || []
    rows.value = members.map(m => ({
      userid: m.userid,
      name: m.name,
      mobile: m.mobile || '',
      email: m.email || '',
      department_names: m.department_names || [],
      role: 'teacher',
      class_id: matchClassByDept(m.department_names || [], classes.value),
      selected: true,
    }))
    if (!members.length) {
      loadError.value = '通讯录为空，请检查第三方平台配置'
    }
  } catch (e: any) {
    loadError.value = e?.response?.data?.message || '拉取通讯录失败（请确认已配置第三方平台应用凭证）'
  } finally {
    loading.value = false
  }
}

function resetAll() {
  rows.value = []
  result.value = null
  search.value = ''
  importError.value = ''
  loadError.value = ''
  copied.value = false
}

const filteredRows = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return rows.value
  return rows.value.filter(r =>
    r.name.toLowerCase().includes(q)
    || r.mobile.includes(q)
    || r.department_names.some(d => d.toLowerCase().includes(q)))
})

const teacherCount = computed(() => rows.value.filter(r => r.role === 'teacher' && r.selected).length)
const studentCount = computed(() => rows.value.filter(r => r.role === 'student' && r.selected).length)
const pendingClassCount = computed(() => rows.value.filter(r => r.selected && r.role === 'student' && !r.class_id).length)
const allSelected = computed(() => rows.value.length > 0 && rows.value.every(r => r.selected))

function toggleAll() {
  const target = !allSelected.value
  rows.value.forEach(r => { r.selected = target })
}

// 批量：将勾选中的学生行目标班级设为指定班级
const bulkClassId = ref<number | ''>('')
function applyBulkClass() {
  if (!bulkClassId.value) return
  rows.value.forEach(r => {
    if (r.selected && r.role === 'student') r.class_id = bulkClassId.value
  })
}

const resultText = computed(() => {
  if (!result.value) return ''
  const acc = result.value.teacher_accounts
    .map(a => `${a.name}：账号 ${a.username}，初始密码 ${a.initial_password}`)
    .join('\n')
  return acc
})

async function copyAccounts() {
  try {
    await navigator.clipboard.writeText(resultText.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch { /* 剪贴板不可用时静默 */ }
}

async function doImport() {
  const picked = rows.value.filter(r => r.selected)
  if (!picked.length) { importError.value = '请至少勾选一名成员'; return }
  const teachers = picked.filter(r => r.role === 'teacher').map(r => ({ name: r.name, mobile: r.mobile, email: r.email }))
  const students = picked.filter(r => r.role === 'student').map(r => ({ name: r.name, class_id: r.class_id }))
  if (students.some(s => !s.class_id)) { importError.value = '仍有学生未选择目标班级（已在列表中标红）'; return }

  importError.value = ''
  importing.value = true
  try {
    const res = await apiPost<{ data: ImportResult }>('/api/v1/admin/third-party/import', { teachers, students }, { skipToast: true })
    result.value = res.data || {
      created_teachers: 0, created_students: 0,
      skipped_teachers: [], skipped_students: [], teacher_accounts: [],
    }
  } catch (e: any) {
    importError.value = e?.response?.data?.message || '导入失败，请稍后重试'
  } finally {
    importing.value = false
  }
}

function finishImport() {
  result.value = null
  rows.value = []
  close()
  emit('imported')
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div class="wechat-import">
      <div class="modal-header">
        <h3 class="modal-title">🏢 从第三方平台导入
          <span v-if="platform" class="platform-tag">{{ PLATFORM_LABELS[platform] || platform }}</span>
        </h3>
        <button class="modal-close" @click="close">✕</button>
      </div>

      <div class="modal-body">
        <!-- ===== 导入结果页 ===== -->
        <div v-if="result" class="result-view">
          <div class="result-summary">
            <div class="result-box result-box--ok">
              <div class="result-box__num">{{ result.created_teachers }}</div>
              <div class="result-box__label">新增教师</div>
            </div>
            <div class="result-box result-box--ok">
              <div class="result-box__num">{{ result.created_students }}</div>
              <div class="result-box__label">新增学生</div>
            </div>
            <div class="result-box" :class="result.skipped_teachers.length ? 'result-box--skip' : ''">
              <div class="result-box__num">{{ result.skipped_teachers.length }}</div>
              <div class="result-box__label">跳过教师</div>
            </div>
            <div class="result-box" :class="result.skipped_students.length ? 'result-box--skip' : ''">
              <div class="result-box__num">{{ result.skipped_students.length }}</div>
              <div class="result-box__label">跳过学生</div>
            </div>
          </div>

          <div v-if="result.skipped_teachers.length" class="skip-list">
            <div class="skip-list__title">跳过的教师（已存在，未重复创建）</div>
            <span v-for="(s, i) in result.skipped_teachers" :key="'t' + i" class="skip-chip">{{ s.name }}（{{ s.reason }}）</span>
          </div>
          <div v-if="result.skipped_students.length" class="skip-list">
            <div class="skip-list__title">跳过的学生（同班同名已存在）</div>
            <span v-for="(s, i) in result.skipped_students" :key="'s' + i" class="skip-chip">{{ s.name }}（{{ s.reason }}）</span>
          </div>

          <div v-if="result.teacher_accounts.length" class="account-block">
            <div class="account-block__head">
              <span class="account-block__title">🔑 新教师账号（初始密码仅此一次显示，请抄送本人）</span>
              <button class="btn btn-sm btn-ghost" @click="copyAccounts">{{ copied ? '已复制 ✓' : '📋 复制全部' }}</button>
            </div>
            <table class="account-table">
              <thead><tr><th>姓名</th><th>登录账号</th><th>初始密码</th></tr></thead>
              <tbody>
                <tr v-for="a in result.teacher_accounts" :key="a.username">
                  <td class="cell-strong">{{ a.name }}</td>
                  <td class="cell-mono">{{ a.username }}</td>
                  <td class="cell-mono">{{ a.initial_password }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="no-new-hint">本次没有新创建的教师（全部已存在）。</div>
        </div>

        <!-- ===== 拉取前空态 ===== -->
        <div v-else-if="!rows.length" class="empty-fetch">
          <div class="empty-fetch__icon">📇</div>
          <p class="empty-fetch__desc">从学校配置的第三方平台（企业微信 / 钉钉 / 飞书）通讯录拉取成员，勾选后批量创建教师与学生账号</p>
          <button class="btn btn-primary" :disabled="loading" @click="loadContacts">
            {{ loading ? '拉取中...' : '📥 拉取通讯录' }}
          </button>
          <div v-if="loadError" class="error-banner iw-error-top">{{ loadError }}</div>
        </div>

        <!-- ===== 成员预览表 ===== -->
        <div v-else>
          <div class="toolbar">
            <span class="count-text">共 {{ rows.length }} 名成员</span>
            <input v-model="search" class="form-input search-input" placeholder="🔍 搜索姓名 / 手机 / 部门">
            <button class="btn btn-sm btn-ghost-card" @click="toggleAll">{{ allSelected ? '全不选' : '全选' }}</button>
          </div>

          <div class="stats-row">
            <span class="stat-chip">👨‍🏫 教师 <strong>{{ teacherCount }}</strong></span>
            <span class="stat-chip">👦 学生 <strong>{{ studentCount }}</strong></span>
            <span v-if="pendingClassCount" class="stat-chip stat-chip--warn">⚠️ {{ pendingClassCount }} 名学生待选班级</span>
          </div>

          <div v-if="studentCount > 0" class="bulk-row">
            <span class="bulk-label">将选中的学生行班级设为</span>
            <select v-model="bulkClassId" class="form-input bulk-select">
              <option :value="''">选择班级…</option>
              <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
            <button class="btn btn-sm btn-ghost" :disabled="!bulkClassId" @click="applyBulkClass">应用</button>
          </div>

          <div class="preview-table-wrapper">
            <table class="preview-table">
              <thead>
                <tr>
                  <th></th>
                  <th>姓名</th>
                  <th>手机</th>
                  <th>部门</th>
                  <th>角色</th>
                  <th>目标班级</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in filteredRows" :key="r.userid">
                  <td><input v-model="r.selected" type="checkbox" class="row-check"></td>
                  <td class="cell-strong">{{ r.name }}</td>
                  <td class="cell-mono">{{ r.mobile || '—' }}</td>
                  <td class="cell-dim">{{ r.department_names.join(' / ') || '—' }}</td>
                  <td>
                    <select v-model="r.role" class="form-input cell-select">
                      <option value="teacher">👨‍🏫 教师</option>
                      <option value="student">👦 学生</option>
                    </select>
                  </td>
                  <td v-if="r.role === 'student'">
                    <select v-model="r.class_id" class="form-input cell-select cell-select--wide"
                      :class="{ 'cell-select--error': !r.class_id }">
                      <option :value="''">选择班级</option>
                      <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
                    </select>
                  </td>
                  <td v-else class="cell-dim">—</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p class="hint-line">💡 学生班级已按部门名自动匹配（含"六年级1班 ↔ 六年级（1）班"模糊匹配），可手动调整或用上方批量设置；教师默认以姓名作为登录账号，重复导入会自动跳过已有账号。</p>
          <div v-if="importError" class="error-banner iw-error-top">{{ importError }}</div>
        </div>
      </div>

      <div class="modal-footer">
        <template v-if="result">
          <button class="btn btn-ghost-card" @click="resetAll">← 继续导入</button>
          <button class="btn btn-solid" @click="finishImport">完成</button>
        </template>
        <template v-else>
          <button class="btn btn-ghost-card" @click="close">取消</button>
          <button v-if="rows.length" class="btn btn-solid" :disabled="importing || pendingClassCount > 0" @click="doImport">
            {{ importing ? '导入中...' : `✅ 确认导入（${teacherCount} 师 / ${studentCount} 生）` }}
          </button>
        </template>
      </div>
    </div>
  </ModalGlass>
</template>

<style scoped>
.wechat-import { max-width: 780px; width: 100%; padding: 4px 0; }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; padding-bottom: 12px; border-bottom: 1px solid var(--color-border); flex-shrink: 0; }
.modal-title { font-size: 16px; font-weight: 700; color: var(--color-text); margin: 0; display: flex; align-items: center; gap: 8px; }
.modal-close { background: none; border: none; color: var(--color-text-secondary); font-size: 20px; cursor: pointer; padding: 0; line-height: 1; }
.modal-body { margin-bottom: 16px; }
.modal-footer { display: flex; gap: 8px; justify-content: flex-end; padding-top: 12px; border-top: 1px solid var(--color-border); margin-top: 16px; }

.platform-tag { font-size: 11px; font-weight: 600; padding: 2px 10px; border-radius: 999px; background: var(--tint-2); color: var(--color-text-secondary); }

/* 拉取前空态 */
.empty-fetch { text-align: center; padding: 32px 16px; }
.empty-fetch__icon { font-size: 40px; margin-bottom: 12px; }
.empty-fetch__desc { font-size: 13px; color: var(--color-text-secondary); margin-bottom: 16px; }
/* 外观统一走全局 .error-banner；这两处横幅之下还有内容块，需要的是「上」间距 */
.iw-error-top { margin-top: 10px; }

/* 工具栏 + 统计 */
.toolbar { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; flex-wrap: wrap; }
.count-text { font-size: 13px; font-weight: 600; color: var(--color-text); white-space: nowrap; }
.search-input { max-width: 220px; padding: 6px 10px; font-size: 12px; }
.stats-row { display: flex; gap: 8px; margin-bottom: 10px; flex-wrap: wrap; }
.stat-chip { font-size: 11px; padding: 3px 10px; border-radius: 999px; background: var(--tint-1); border: 1px solid var(--tint-2); color: var(--color-text-secondary); }
.stat-chip strong { color: var(--color-text); }
.stat-chip--warn { background: rgba(245, 158, 11, 0.1); border-color: rgba(245, 158, 11, 0.3); color: var(--color-warning-text); }

/* 批量设班 */
.bulk-row { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; flex-wrap: wrap; }
.bulk-label { font-size: 12px; color: var(--color-text-secondary); }
.bulk-select { max-width: 200px; padding: 5px 8px; font-size: 12px; }

/* 预览表 */
.preview-table-wrapper { max-height: 360px; overflow: auto; border: 1px solid var(--tint-2); border-radius: 10px; }
.preview-table { width: 100%; border-collapse: collapse; font-size: 12px; }
.preview-table th { background: var(--tint-1); padding: 6px 8px; text-align: left; border-bottom: 1px solid var(--tint-3); font-size: 11px; color: var(--color-text-secondary); white-space: nowrap; position: sticky; top: 0; }
.preview-table td { padding: 5px 8px; border-bottom: 1px solid var(--tint-2); white-space: nowrap; color: var(--color-text); }
.preview-table tr:last-child td { border-bottom: none; }
.cell-strong { font-weight: 600; }
.cell-dim { color: var(--color-text-secondary); }
.cell-mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; }
.row-check { accent-color: var(--color-primary); }
.cell-select { padding: 4px 6px; font-size: 12px; }
.cell-select--wide { min-width: 140px; }
.cell-select--error { border-color: rgba(239, 68, 68, 0.55) !important; background: rgba(239, 68, 68, 0.05); }
.hint-line { font-size: 11px; color: var(--color-text-secondary); margin-top: 8px; line-height: 1.6; }

/* 结果页 */
.result-view { padding: 4px 0; }
.result-summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 14px; }
.result-box { padding: 14px; text-align: center; background: var(--tint-1); border: 1px solid var(--tint-2); border-radius: 12px; }
.result-box--ok { background: rgba(16, 185, 129, 0.07); border-color: rgba(16, 185, 129, 0.25); }
.result-box--skip { background: rgba(245, 158, 11, 0.07); border-color: rgba(245, 158, 11, 0.3); }
.result-box__num { font-size: 24px; font-weight: 800; color: var(--color-text); }
.result-box__label { font-size: 11px; color: var(--color-text-secondary); margin-top: 2px; }
.skip-list { margin-bottom: 10px; display: flex; flex-wrap: wrap; gap: 6px; align-items: center; }
.skip-list__title { font-size: 12px; font-weight: 600; color: var(--color-text-secondary); width: 100%; margin-bottom: 2px; }
.skip-chip { font-size: 11px; padding: 3px 10px; border-radius: 999px; background: rgba(245, 158, 11, 0.08); border: 1px solid rgba(245, 158, 11, 0.25); color: var(--color-warning-text); }
.account-block { margin-top: 12px; border: 1px solid var(--tint-2); border-radius: 12px; overflow: hidden; }
.account-block__head { display: flex; align-items: center; justify-content: space-between; gap: 8px; padding: 10px 14px; background: var(--tint-1); flex-wrap: wrap; }
.account-block__title { font-size: 12px; font-weight: 600; color: var(--color-text); }
.account-table { width: 100%; border-collapse: collapse; font-size: 12px; }
.account-table th { padding: 7px 14px; text-align: left; font-size: 11px; color: var(--color-text-secondary); border-top: 1px solid var(--tint-2); border-bottom: 1px solid var(--tint-2); background: var(--tint-1); }
.account-table td { padding: 7px 14px; border-bottom: 1px solid var(--tint-2); }
.account-table tr:last-child td { border-bottom: none; }
.no-new-hint { margin-top: 12px; font-size: 12px; color: var(--color-text-secondary); text-align: center; padding: 10px; background: var(--tint-1); border-radius: 10px; }

/* 按钮（对齐设计系统） */
.btn { padding: 8px 16px; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; border: 1px solid transparent; transition: all 0.15s; font-family: inherit; }
.btn:disabled { opacity: 0.55; cursor: not-allowed; }
.btn-primary { background: rgba(124, 58, 237, 0.12); border: 1px solid rgba(124, 58, 237, 0.25); color: var(--color-primary); }
.btn-primary:hover { background: rgba(124, 58, 237, 0.2); }
.btn-solid { background: var(--color-primary); color: #fff; border: 1px solid var(--color-primary); }
.btn-solid:hover:not(:disabled) { filter: brightness(1.06); }
.btn-ghost { background: transparent; border: 1px solid var(--tint-3); color: var(--color-text-secondary); }
.btn-ghost:hover:not(:disabled) { background: var(--tint-2); color: var(--color-text); }
.btn-ghost-card { background: var(--color-bg-card); color: var(--color-text); border: 1px solid var(--color-border); }
.btn-sm { padding: 5px 10px; font-size: 12px; }
</style>
