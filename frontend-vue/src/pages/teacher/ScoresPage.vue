<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet, apiPost } from '@/utils/api'
import { useAppMode } from '@/composables/useAppMode'
import { categoryLabel } from '@/utils/scoreRules'
import { getSpeciesEmoji, PET_SERIES } from '@/utils/petData'
import PetSprite from '@/components/pet/PetSprite.vue'
import PetDetailModal from '@/components/pet/PetDetailModal.vue'
import ScoreStudentCard from '@/components/score/ScoreStudentCard.vue'
import ScoreRuleModal from '@/components/score/ScoreRuleModal.vue'
import ScoreReasonModal from '@/components/score/ScoreReasonModal.vue'
import type { ApiResponse, ScoreRule } from '@/types'
import type { CardStudent } from '@/components/score/ScoreStudentCard.vue'

// ===== 模式（教师完整 / 教室端+班级码基础） =====
const { isClassroomMode, isTeacherMode } = useAppMode()
const router = useRouter()

// ===== 共用数据 =====
const students = ref<CardStudent[]>([])
const loading = ref(true)
const searchQuery = ref('')
const selectedIds = ref<number[]>([])
function isSelected(id: number) { return selectedIds.value.includes(id) }
function toggleSelect(id: number) {
  const i = selectedIds.value.indexOf(id)
  if (i >= 0) selectedIds.value.splice(i, 1)
  else selectedIds.value.push(id)
}
function clearSelect() { selectedIds.value = [] }

// 浮动积分文本
interface FloatText { id: number; x: number; y: number; text: string; color: string }
const floatTexts = ref<FloatText[]>([])
let floatId = 0
function showFloatText(studentId: number, points: number) {
  const card = document.getElementById(`card-${studentId}`)
  if (!card) return
  const rect = card.getBoundingClientRect()
  const id = floatId++
  floatTexts.value.push({
    id,
    x: rect.left + rect.width / 2,
    y: rect.top - 10,
    text: points > 0 ? `+${points}` : `${points}`,
    color: points > 0 ? '#10B981' : '#EF4444',
  })
  setTimeout(() => { floatTexts.value = floatTexts.value.filter(f => f.id !== id) }, 1200)
}
const lastFlashId = ref<number | null>(null)
function flashCard(sid: number) {
  lastFlashId.value = sid
  setTimeout(() => { if (lastFlashId.value === sid) lastFlashId.value = null }, 500)
}

// ===== 教师端：等级（用于筛选/排序） =====
const LEVEL_SCORES = [0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450]
function calcLevel(score: number): number {
  let lv = 1
  for (let i = LEVEL_SCORES.length - 1; i >= 0; i--) {
    if (score >= LEVEL_SCORES[i]) { lv = i + 1; break }
  }
  return Math.min(lv, 12)
}

// ===== 教师端状态 =====
const rules = ref<ScoreRule[]>([])
const scoreSummary = ref({ total: 0, today: 0, this_week: 0 })
const loadError = ref('')
const activeFilter = ref<'all' | 'high' | 'mid' | 'low'>('all')
type SortKey = 'no' | 'surname' | 'score'
const sortBy = ref<SortKey>('no')
const SORT_OPTIONS: Array<{ key: SortKey; label: string }> = [
  { key: 'no', label: '🔢 学号' },
  { key: 'surname', label: '👤 姓氏' },
  { key: 'score', label: '⭐ 积分' },
]
const giveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const batchStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const undoStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const activeReason = ref<string>('')
const recentScores = ref<Array<{id: number; student_name: string; amount: number; reason: string; created_at: string}>>([])
const historyLoading = ref(false)

const filteredStudents = computed(() => {
  const list = students.value.filter(s => {
    const matchName = s.name.includes(searchQuery.value)
    const level = calcLevel(s.total_score)
    if (activeFilter.value === 'high') return matchName && level >= 10
    if (activeFilter.value === 'mid') return matchName && level >= 4 && level <= 9
    if (activeFilter.value === 'low') return matchName && level <= 3
    return matchName
  })
  const arr = [...list]
  if (sortBy.value === 'score') arr.sort((a, b) => b.total_score - a.total_score)
  else if (sortBy.value === 'surname') arr.sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
  else arr.sort((a, b) => (a.student_no || '').localeCompare(b.student_no || '', 'zh-CN', { numeric: true }))
  return arr
})
const positiveRules = computed(() => rules.value.filter(r => r.is_positive))
// 教室端：姓名搜索 + 排序（与教师端同 sortBy）
const filtered = computed(() => {
  const list = searchQuery.value ? students.value.filter(s => s.name.includes(searchQuery.value)) : [...students.value]
  if (sortBy.value === 'score') list.sort((a, b) => b.total_score - a.total_score)
  else if (sortBy.value === 'surname') list.sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
  else list.sort((a, b) => (a.student_no || '').localeCompare(b.student_no || '', 'zh-CN', { numeric: true }))
  return list
})
const negativeRules = computed(() => rules.value.filter(r => !r.is_positive))
function groupedRules(type: 'add' | 'sub'): Array<{ category: string; label: string; rules: ScoreRule[] }> {
  const list = type === 'add' ? positiveRules.value : negativeRules.value
  const map = new Map<string, ScoreRule[]>()
  for (const r of list) {
    if (!map.has(r.category)) map.set(r.category, [])
    map.get(r.category)!.push(r)
  }
  return [...map.entries()].map(([category, rules]) => ({ category, label: categoryLabel(category), rules }))
}
function getBatchStatus(ruleId: number) { return batchStatus.value[ruleId] || 'idle' }
function getBatchBtnText(rule: ScoreRule): string {
  const s = getBatchStatus(rule.id)
  const map: Record<string, string> = { idle: `${rule.amount > 0 ? '+' : ''}${rule.amount} ${rule.name}`, loading: '处理中...', success: '已完成', error: '操作失败' }
  return map[s]
}
function getBatchBtnStyle(ruleId: number): Record<string, string> {
  const s = getBatchStatus(ruleId)
  if (s === 'idle') return {}
  const bgMap: Record<string, string> = { loading: '#f59e0b', success: '#10b981', error: '#ef4444' }
  return { background: bgMap[s], borderColor: 'transparent', color: '#fff' }
}
function getUndoStatus(scoreId: number) { return undoStatus.value[scoreId] || 'idle' }
function getUndoBtnText(scoreId: number): string {
  const s = getUndoStatus(scoreId)
  const map: Record<string, string> = { idle: '↩ 撤回', loading: '撤回中...', success: '已撤回', error: '撤回失败' }
  return map[s]
}
function getUndoBtnStyle(scoreId: number): Record<string, string> {
  const s = getUndoStatus(scoreId)
  const map: Record<string, string> = { idle: 'var(--color-danger)', loading: '#f59e0b', success: '#10b981', error: '#ef4444' }
  return { color: map[s] }
}

// ===== 教师端：单学生按规则加减分 =====
const showModal = ref(false)
const modalType = ref<'add' | 'sub'>('add')
const modalStudent = ref<CardStudent | null>(null)
const modalError = ref('')
function openModal(student: CardStudent, type: 'add' | 'sub') {
  modalStudent.value = student
  modalType.value = type
  modalError.value = ''
  showModal.value = true
}
function closeModal() { showModal.value = false; modalStudent.value = null }
async function executeAction(rule: ScoreRule) {
  const student = modalStudent.value
  if (!student) return
  if (student.total_score + rule.amount < 0) { modalError.value = '积分不能为负数'; return }
  activeReason.value = rule.name
  giveStatus.value = 'loading'
  try {
    await apiPost(`/api/v1/teacher/scores/by-rule/${rule.id}`, { student_id: student.id })
    student.total_score = Math.max(0, student.total_score + rule.amount)
    giveStatus.value = 'success'
    showFloatText(student.id, rule.amount)
    flashCard(student.id)
  } catch {
    // 离线模式
    student.total_score = Math.max(0, student.total_score + rule.amount)
    giveStatus.value = 'success'
    showFloatText(student.id, rule.amount)
    flashCard(student.id)
  }
  setTimeout(() => { giveStatus.value = 'idle'; closeModal() }, 800)
}

// ===== 教师端：批量规则 =====
const batchModal = ref(false)
const batchType = ref<'add' | 'sub'>('add')
const batchBusy = ref(false)
function openBatchModal(type: 'add' | 'sub') {
  if (selectedIds.value.length === 0) return
  batchType.value = type
  batchBusy.value = false
  batchModal.value = true
}
async function executeBatch(rule: ScoreRule) {
  if (selectedIds.value.length === 0) return
  batchBusy.value = true
  const applyLocal = () => {
    selectedIds.value.forEach(id => {
      const s = students.value.find(x => x.id === id)
      if (s) {
        s.total_score = Math.max(0, s.total_score + rule.amount)
        showFloatText(s.id, rule.amount)
      }
    })
  }
  try {
    await apiPost('/api/v1/teacher/scores/batch-give', { student_ids: selectedIds.value, points: rule.amount, reason: rule.name })
    applyLocal()
  } catch { applyLocal() }
  finally {
    batchBusy.value = false
    batchModal.value = false
    clearSelect()
  }
}
async function handleBatchRuleScore(rule: ScoreRule) {
  if (students.value.length === 0) return
  batchStatus.value[rule.id] = 'loading'
  const applyLocal = () => { students.value.forEach(s => { s.total_score += rule.amount }) }
  try {
    await apiPost('/api/v1/teacher/scores/batch-give', { student_ids: students.value.map(s => s.id), points: rule.amount, reason: rule.name })
    batchStatus.value[rule.id] = 'success'
    applyLocal()
  } catch {
    batchStatus.value[rule.id] = 'success'
    applyLocal()
  }
  setTimeout(() => { batchStatus.value[rule.id] = 'idle' }, 1500)
}
async function undoScore(scoreId: number) {
  undoStatus.value[scoreId] = 'loading'
  try {
    await apiPost(`/api/v1/teacher/scores/${scoreId}/undo`, {})
    undoStatus.value[scoreId] = 'success'
    setTimeout(() => { undoStatus.value[scoreId] = 'idle' }, 1500)
    await loadRecentScores()
  } catch {
    undoStatus.value[scoreId] = 'error'
    setTimeout(() => { undoStatus.value[scoreId] = 'idle' }, 3000)
  }
}
async function loadRecentScores() {
  historyLoading.value = true
  try {
    const res = await apiGet<{ data: Array<{id: number; student_name: string; amount: number; reason: string; created_at: string}> }>('/api/v1/teacher/scores/recent', { skipToast: true })
    recentScores.value = (res.data || []).slice(0, 20)
  } catch { recentScores.value = [] }
  finally { historyLoading.value = false }
}

// ===== 教室端状态 =====
const token = ref('')
const classTotal = computed(() => students.value.reduce((sum, s) => sum + s.total_score, 0))
const stepValues = ref<Record<number, number>>({})
const editingStep = ref<number | null>(null)
const editInput = ref('1')
function getStep(sid: number) { return stepValues.value[sid] || 1 }
function startEdit(sid: number) { editingStep.value = sid; editInput.value = String(getStep(sid)) }
function saveEdit(sid: number, event?: Event) {
  const val = parseInt(event ? (event.target as HTMLInputElement).value : editInput.value)
  if (val >= 1 && val <= 100) stepValues.value[sid] = val
  editingStep.value = null
}

// 教室端理由分组：优先用后台积分规则（与教师端同源，原因一致），加载失败回退内置分组
interface ReasonGroup { title: string; emoji: string; color: string; items: string[] }
const FALLBACK_ADD: ReasonGroup[] = [
  { title: '学习表现', emoji: '📚', color: '#3B82F6', items: ['✅ 作业优秀', '🏆 挑战难题', '📚 阅读之星', '🎨 科技创新'] },
  { title: '行为表现', emoji: '🏅', color: '#10B981', items: ['📖 举手发言', '🤝 帮助同学', '📝 认真听讲', '🔍 专注课堂'] },
  { title: '品德修养', emoji: '⭐', color: '#F59E0B', items: ['🧹 遵守纪律', '💬 积极互动', '🌟 诚实守信', '🏃 体育锻炼'] },
]
const FALLBACK_SUB: ReasonGroup[] = [
  { title: '学习懈怠', emoji: '📕', color: '#EF4444', items: ['⚠️ 上课走神', '📕 作业缺交'] },
  { title: '课堂纪律', emoji: '🗣️', color: '#F97316', items: ['🗣️ 打扰课堂', '📱 课堂喧哗', '💬 说脏话'] },
  { title: '行为品德', emoji: '🤕', color: '#F59E0B', items: ['🏃 追逐打闹', '😴 趴桌睡觉', '⚡ 与同学冲突', '🗑️ 乱扔垃圾'] },
]
const CATEGORY_META: Record<string, { emoji: string; color: string }> = {
  classroom: { emoji: '📖', color: '#3B82F6' },
  homework: { emoji: '📝', color: '#10B981' },
  behavior: { emoji: '🌟', color: '#F59E0B' },
  literacy: { emoji: '📊', color: '#8B5CF6' },
  daily: { emoji: '📅', color: '#F97316' },
}
const classroomRules = ref<ScoreRule[]>([])
function rulesToReasonGroups(rules: ScoreRule[], type: 'add' | 'sub'): ReasonGroup[] {
  const list = rules.filter(r => (type === 'add' ? r.amount > 0 : r.amount < 0))
  const map = new Map<string, string[]>()
  for (const r of list) {
    const cat = r.category || 'custom'
    if (!map.has(cat)) map.set(cat, [])
    map.get(cat)!.push(r.name)
  }
  return [...map.entries()].map(([cat, items]) => ({
    title: categoryLabel(cat), // 已含 emoji（如「📖 课堂表现」）
    emoji: '',
    color: CATEGORY_META[cat]?.color || '#6B7280',
    items,
  }))
}
const reasonGroupsAdd = computed<ReasonGroup[]>(() => classroomRules.value.length ? rulesToReasonGroups(classroomRules.value, 'add') : FALLBACK_ADD)
const reasonGroupsSub = computed<ReasonGroup[]>(() => classroomRules.value.length ? rulesToReasonGroups(classroomRules.value, 'sub') : FALLBACK_SUB)

// 教室端：单学生理由加减分弹窗
const classShowModal = ref(false)
const classModalType = ref<'add' | 'sub'>('add')
const classModalStudent = ref<CardStudent | null>(null)
const classBusy = ref(false)
const classActionError = ref('')
const reasonGroups = computed(() => classModalType.value === 'add' ? reasonGroupsAdd.value : reasonGroupsSub.value)
function openReasonModal(s: CardStudent, type: 'add' | 'sub') {
  classModalStudent.value = s
  classModalType.value = type
  classActionError.value = ''
  classBusy.value = false
  classShowModal.value = true
}
async function confirmReasonAction(reasons: string[]) {
  const s = classModalStudent.value
  if (!s || classBusy.value || reasons.length === 0) return
  const step = getStep(s.id)
  const points = classModalType.value === 'add' ? step * reasons.length : -step * reasons.length
  classBusy.value = true
  try {
    await apiPost('/api/v1/display/scores/give', { token: token.value, student_id: s.id, points, reason: reasons.join('、') })
    s.total_score = Math.max(0, s.total_score + points)
    showFloatText(s.id, points)
    flashCard(s.id)
    classActionError.value = ''
    classShowModal.value = false
    classModalStudent.value = null
  } catch (e: any) {
    const msg = e?.response?.data?.message || ''
    classActionError.value = msg.includes('30') ? '单次超过 30 分需要教师账号登录操作' : (msg || '操作失败，请稍后重试')
    setTimeout(() => { classActionError.value = '' }, 3000)
  } finally {
    classBusy.value = false
  }
}

// 教室端：批量理由加减分弹窗
const classBatchModal = ref(false)
const classBatchType = ref<'add' | 'sub'>('add')
const classBatchBusy = ref(false)
const classBatchError = ref('')
const batchReasonGroups = computed(() => classBatchType.value === 'add' ? reasonGroupsAdd.value : reasonGroupsSub.value)
function openClassBatchModal(type: 'add' | 'sub') {
  classBatchType.value = type
  classBatchError.value = ''
  classBatchBusy.value = false
  classBatchModal.value = true
}
async function confirmClassBatch(reasons: string[]) {
  const ids = selectedIds.value
  if (!ids.length || classBatchBusy.value || reasons.length === 0) return
  const step = getStep(ids[0])
  const points = classBatchType.value === 'add' ? step * reasons.length : -step * reasons.length
  classBatchBusy.value = true
  try {
    await apiPost('/api/v1/display/scores/batch-give', { token: token.value, student_ids: ids, points, reason: reasons.join('、') })
    for (const s of students.value) {
      if (ids.includes(s.id)) {
        s.total_score = Math.max(0, s.total_score + points)
        showFloatText(s.id, points)
        flashCard(s.id)
      }
    }
    classBatchError.value = ''
    classBatchModal.value = false
    clearSelect()
  } catch (e: any) {
    const msg = e?.response?.data?.message || ''
    classBatchError.value = msg.includes('30') ? '单次超过 30 分需要教师账号登录操作' : (msg || '操作失败，请稍后重试')
    setTimeout(() => { classBatchError.value = '' }, 3000)
  } finally {
    classBatchBusy.value = false
  }
}

/** 跳转教师登录（保留 class_token，登录后进入教师完整模式可正常操作） */
function goTeacherLogin() {
  router.push({ name: 'login', query: { role: 'teacher' } })
}

// ===== 教室端：宠物系统 =====
const classPetSeries = ref('')
const pickerSeriesList = computed(() => {
  if (!classPetSeries.value) return PET_SERIES
  const hit = PET_SERIES.find(s => s.id === classPetSeries.value)
  return hit ? [hit] : PET_SERIES
})
const showPetDetail = ref(false)
const detailStudent = ref<CardStudent | null>(null)
function openPetDetail(s: CardStudent) { detailStudent.value = s; showPetDetail.value = true }
const showPetPicker = ref(false)
const petPickerStudent = ref<CardStudent | null>(null)
const switchingPet = ref(false)
const switchStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const switchError = ref('')
/** 选中待切换的宠物（在弹窗内联确认，不弹第二个弹窗） */
const pendingSpecies = ref<{ speciesId: string; name: string } | null>(null)
function openPetPicker(s: CardStudent) {
  petPickerStudent.value = s
  pendingSpecies.value = null
  switchError.value = ''
  showPetPicker.value = true
}
function handlePick(speciesId: string, name: string) {
  if (petPickerStudent.value?.pet_species === speciesId) return
  pendingSpecies.value = { speciesId, name }
  switchError.value = ''
}
function clearPending() { pendingSpecies.value = null }
/** 切换宠物所需积分：首次免费；整班切系列后 3 天免费窗口内免费；否则按当前宠物等级扣（后端 switchCost = 5×等级） */
function getSwitchCost(s: CardStudent | null): number {
  if (!s || !s.pet_name) return 0
  if (s.free_pick) return 0
  return 5 * Math.max(1, s.pet_level || 1)
}
async function executeSwitch() {
  if (!pendingSpecies.value) return
  const s = petPickerStudent.value
  if (!s || switchingPet.value) return
  const speciesId = pendingSpecies.value.speciesId
  switchingPet.value = true
  switchStatus.value = 'loading'
  try {
    const res = await apiPost<{ data: { pet_emoji: string; total_score: number } }>(
      '/api/v1/display/pets/switch', { token: token.value, student_id: s.id, pet_species: speciesId }
    )
    s.pet_species = speciesId
    s.pet_emoji = res.data.pet_emoji
    s.total_score = res.data.total_score
    switchStatus.value = 'success'
    switchError.value = ''
    setTimeout(() => { switchStatus.value = 'idle' }, 1500)
    showPetPicker.value = false
    pendingSpecies.value = null
  } catch (e: any) {
    switchStatus.value = 'error'
    switchError.value = e?.response?.data?.message || '切换失败，请稍后重试'
    setTimeout(() => { switchStatus.value = 'idle'; switchError.value = '' }, 3000)
  } finally {
    switchingPet.value = false
  }
}

// ===== 数据加载（按模式互斥） =====
onMounted(async () => {
  if (isTeacherMode.value) {
    loading.value = true
    try {
      const [sRes, rRes, sumRes] = await Promise.all([
        apiGet<ApiResponse<CardStudent[]>>('/api/v1/teacher/students?per_page=100', { skipToast: true }),
        apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules', { skipToast: true }),
        apiGet<ApiResponse<{ total: number; today: number; this_week: number }>>('/api/v1/teacher/scores/summary', { skipToast: true }),
      ])
      students.value = sRes.data || []
      rules.value = rRes.data || []
      scoreSummary.value = sumRes.data || { total: 0, today: 0, this_week: 0 }
      loadError.value = ''
    } catch {
      loadError.value = '数据加载失败，已显示演示数据'
      const names = ['张小明', '李小红', '王小刚', '赵小丽', '刘小强', '陈小美', '周小龙', '吴小凤', '郑小天', '孙小艺',
        '胡小勇', '林小静', '郭小峰', '何小婷', '高小磊', '罗小欣', '梁小涛', '宋小敏', '唐小亮', '韩小洁']
      students.value = names.map((name, i) => ({
        id: i + 1, name, total_score: Math.floor(Math.random() * 400) + 20, class_id: 1, status: 'active' as const,
      }))
    } finally {
      loading.value = false
    }
    await loadRecentScores()
    return
  }

  // 教室端
  token.value = sessionStorage.getItem('class_token') || ''
  if (!token.value) { loading.value = false; return }
  try {
    const res = await apiGet<{ data: CardStudent[] }>('/api/v1/display/students', { params: { token: token.value } })
    students.value = (res.data || []).map(s => ({ ...s, pet_emoji: s.pet_species ? getSpeciesEmoji(s.pet_species) : '🥚' }))
  } catch { /* ignore */ } finally { loading.value = false }
  try {
    const cRes = await apiGet<{ data: { pet_series?: string | null } }>('/api/v1/display/class-settings', { params: { token: token.value } })
    classPetSeries.value = cRes.data?.pet_series || ''
  } catch { classPetSeries.value = '' }
  // 教室端加减分原因：接后台积分规则（与教师端同源）
  try {
    const rRes = await apiGet<{ data: ScoreRule[] }>('/api/v1/display/scores/rules', { params: { token: token.value } })
    classroomRules.value = rRes.data || []
  } catch { classroomRules.value = [] }
})
</script>

<template>
  <div class="scores-page">
    <!-- 教师端加载失败提示 -->
    <div v-if="isTeacherMode && loadError" class="load-error">
      ⚠️ {{ loadError }}
    </div>

    <!-- 页头：标题 + 统计/班级总分 -->
    <div class="page-top">
      <div class="page-header">
        <h2 class="page-title">✏️ 课堂评价</h2>
        <span class="page-subtitle">点击 +/− 选择行为原因</span>
        <div v-if="isClassroomMode && students.length" class="class-total">
          <span class="ct-icon">⭐</span>
          <span class="ct-label">班级总分</span>
          <strong class="ct-value">{{ classTotal.toLocaleString() }}</strong>
        </div>
      </div>

      <!-- 教师端：积分统计 -->
      <div v-if="isTeacherMode" class="stats-row">
        <div class="stat-chip"><span class="chip-icon">⭐</span><span>累计 {{ scoreSummary.total.toLocaleString() }}</span></div>
        <div class="stat-chip"><span class="chip-icon">📅</span><span>今日 {{ scoreSummary.today.toLocaleString() }}</span></div>
        <div class="stat-chip"><span class="chip-icon">📈</span><span>本周 {{ scoreSummary.this_week.toLocaleString() }}</span></div>
        <div class="stat-chip"><span class="chip-icon">📋</span><span>规则 {{ rules.length }}</span></div>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载学生数据...</p>
    </div>

    <template v-else>
      <!-- 教室端：首次使用认养提示 -->
      <div v-if="isClassroomMode && students.some(s => !s.pet_name)" class="adopt-banner">
        <span class="adopt-icon">🎉</span>
        <div>
          <div class="adopt-title">为同学们认养第一只宠物吧！</div>
          <div class="adopt-desc">点击学生卡片左侧的宠物区域，可以为还没有宠物的同学免费选择第一只宠物 🆓</div>
        </div>
      </div>

      <!-- 教师端：工具栏 -->
      <div v-if="isTeacherMode" class="toolbar">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input v-model="searchQuery" type="text" placeholder="搜索学生姓名..." class="search-input" />
        </div>
        <div class="filter-group">
          <button class="filter-tag" :class="{ active: activeFilter === 'all' }" @click="activeFilter = 'all'">全部</button>
          <button class="filter-tag" :class="{ active: activeFilter === 'high' }" @click="activeFilter = 'high'">⭐ 巅峰</button>
          <button class="filter-tag" :class="{ active: activeFilter === 'mid' }" @click="activeFilter = 'mid'">🌱 成长</button>
          <button class="filter-tag" :class="{ active: activeFilter === 'low' }" @click="activeFilter = 'low'">🥚 幼年</button>
        </div>
        <div class="sort-group">
          <span class="sort-label">排序</span>
          <button v-for="opt in SORT_OPTIONS" :key="opt.key" class="sort-tag" :class="{ active: sortBy === opt.key }" @click="sortBy = opt.key">{{ opt.label }}</button>
        </div>
      </div>
      <!-- 教室端：工具栏（搜索 + 排序） -->
      <div v-else class="toolbar">
        <div class="class-search">
          <span>🔍</span>
          <input v-model="searchQuery" type="text" placeholder="搜索学生姓名..." />
        </div>
        <div class="sort-group">
          <span class="sort-label">排序</span>
          <button v-for="opt in SORT_OPTIONS" :key="opt.key" class="sort-tag" :class="{ active: sortBy === opt.key }" @click="sortBy = opt.key">{{ opt.label }}</button>
        </div>
      </div>

      <!-- 批量操作栏（两端共用） -->
      <div v-if="selectedIds.length > 0" class="batch-bar">
        <span class="batch-info">✅ 已选 <strong>{{ selectedIds.length }}</strong> 名学生</span>
        <button class="batch-btn batch-add" @click="isTeacherMode ? openBatchModal('add') : openClassBatchModal('add')">＋ 批量加分</button>
        <button class="batch-btn batch-sub" @click="isTeacherMode ? openBatchModal('sub') : openClassBatchModal('sub')">－ 批量减分</button>
        <button class="batch-btn batch-clear" @click="clearSelect">取消选择</button>
      </div>

      <!-- 教师端：快捷规则条 -->
      <div v-if="isTeacherMode && (positiveRules.length || negativeRules.length)" class="quick-rules">
        <span class="qr-label">快捷规则</span>
        <div class="qr-group">
          <button v-for="r in positiveRules" :key="r.id" class="qr-btn qr-add" :style="getBatchBtnStyle(r.id)" @click="handleBatchRuleScore(r)" :disabled="getBatchStatus(r.id) === 'loading'">{{ getBatchBtnText(r) }}</button>
          <button v-for="r in negativeRules" :key="r.id" class="qr-btn qr-sub" :style="getBatchBtnStyle(r.id)" @click="handleBatchRuleScore(r)" :disabled="getBatchStatus(r.id) === 'loading'">{{ getBatchBtnText(r) }}</button>
        </div>
      </div>

      <!-- 学生卡片网格 -->
      <div class="student-grid">
        <div v-if="(isTeacherMode ? filteredStudents : filtered).length === 0" class="empty-grid">
          👀 没有找到匹配的学生
        </div>
        <ScoreStudentCard
          v-for="s in (isTeacherMode ? filteredStudents : filtered)"
          :key="s.id"
          :student="s"
          :is-teacher-mode="isTeacherMode"
          :selected="isSelected(s.id)"
          :flash="lastFlashId === s.id"
          :editing-step="isClassroomMode && editingStep === s.id"
          :step-value="getStep(s.id)"
          :edit-input="editInput"
          @toggle-select="toggleSelect(s.id)"
          @open-detail="openPetDetail(s)"
          @open-modal="(type) => isTeacherMode ? openModal(s, type) : openReasonModal(s, type)"
          @start-edit="startEdit(s.id)"
          @save-edit="saveEdit(s.id, $event)"
          @update:edit-input="editInput = $event"
        />
      </div>
    </template>

    <!-- 教师端：单学生规则弹窗 -->
    <ScoreRuleModal
      v-if="isTeacherMode"
      :show="showModal && !!modalStudent"
      :type="modalType"
      :title="(modalType === 'add' ? '🌟 加分 · ' : '⚠️ 减分 · ') + (modalStudent?.name || '')"
      :subtitle="'选择积分规则，按规则真实分值' + (modalType === 'add' ? '加分' : '减分')"
      :groups="groupedRules(modalType)"
      :busy="giveStatus === 'loading'"
      :status="giveStatus"
      :active-rule-name="activeReason"
      :error="modalError"
      @close="closeModal"
      @apply="executeAction"
    />
    <!-- 教师端：批量规则弹窗 -->
    <ScoreRuleModal
      v-if="isTeacherMode"
      :show="batchModal"
      :type="batchType"
      :title="(batchType === 'add' ? '🌟 批量加分' : '⚠️ 批量减分') + '（' + selectedIds.length + ' 名学生）'"
      :subtitle="'选择积分规则，统一应用到所选学生'"
      :groups="groupedRules(batchType)"
      :busy="batchBusy"
      :status="'idle'"
      :active-rule-name="''"
      :error="''"
      @close="batchModal = false"
      @apply="executeBatch"
    />

    <!-- 教室端：单学生理由弹窗（含 ±30 内联提示） -->
    <ScoreReasonModal
      v-if="isClassroomMode"
      :show="classShowModal && !!classModalStudent"
      :type="classModalType"
      :student-name="classModalStudent?.name"
      :step-value="getStep(classModalStudent?.id ?? 0)"
      :busy="classBusy"
      :error="classActionError"
      :reason-groups="reasonGroups"
      @close="classShowModal = false; classModalStudent = null"
      @confirm="confirmReasonAction"
      @go-login="goTeacherLogin"
    />
    <!-- 教室端：批量理由弹窗（含 ±30 内联提示） -->
    <ScoreReasonModal
      v-if="isClassroomMode"
      :show="classBatchModal"
      :type="classBatchType"
      :selected-count="selectedIds.length"
      :step-value="getStep(selectedIds[0] ?? 0)"
      :busy="classBatchBusy"
      :error="classBatchError"
      :reason-groups="batchReasonGroups"
      @close="classBatchModal = false"
      @confirm="confirmClassBatch"
      @go-login="goTeacherLogin"
    />

    <!-- 教室端：宠物角色介绍 -->
    <PetDetailModal
      v-if="isClassroomMode && showPetDetail && detailStudent?.pet_species"
      :species-id="detailStudent.pet_species"
      :level="detailStudent.pet_level || 1"
      :score="detailStudent.total_score"
      :show-pet-switch="true"
      @close="showPetDetail = false"
      @switch-pet="showPetDetail = false; openPetPicker(detailStudent)"
    />

    <!-- 教室端：宠物选择器 -->
    <Transition name="fade">
      <div v-if="isClassroomMode && showPetPicker && petPickerStudent" class="overlay-high" @click.self="showPetPicker = false">
        <div class="picker-box">
          <div class="picker-head">
            <span class="picker-emoji">{{ petPickerStudent.pet_emoji }}</span>
            <div>
              <div class="picker-name">{{ petPickerStudent.name }} · 选择宠物</div>
              <div v-if="petPickerStudent.pet_name && petPickerStudent.free_pick" class="pick-badge pick-free">🎉 免费窗口期内，本次切换免费！</div>
              <div v-else-if="petPickerStudent.pet_name" class="pick-badge pick-cost">💰 本次切换扣除 <strong>{{ getSwitchCost(petPickerStudent) }}</strong> 积分 · 保留等级</div>
              <div v-else class="pick-badge pick-free">🎉 首次免费选择，不扣积分</div>
            </div>
            <button class="picker-close" @click="showPetPicker = false">✕</button>
          </div>
          <div v-for="series in pickerSeriesList" :key="series.id" class="picker-series">
            <div class="series-name">{{ series.emoji }} {{ series.name }}</div>
            <div class="series-grid">
              <button
                v-for="sp in series.species" :key="sp.id"
                @click="handlePick(sp.id, sp.name)"
                :disabled="switchingPet || petPickerStudent.pet_species === sp.id"
                class="species-btn"
                :style="petPickerStudent.pet_species === sp.id ? 'border-color:rgba(16,185,129,0.35);background:rgba(16,185,129,0.08);cursor:default;' : (pendingSpecies?.speciesId === sp.id ? 'border-color:var(--color-primary);background:rgba(167,139,250,0.12);' : '')"
              >
                <div class="species-sprite">
                  <PetSprite :species-id="sp.id" :level="6" />
                </div>
                <div class="species-name">{{ sp.name }}</div>
                <div v-if="petPickerStudent.pet_species === sp.id" class="species-current">✓ 当前</div>
                <div v-else-if="pendingSpecies?.speciesId === sp.id" class="species-current">✓ 已选</div>
              </button>
            </div>
          </div>
          <!-- 内联确认：选中后在弹窗内确认，不弹第二个弹窗 -->
          <div v-if="pendingSpecies" class="pick-confirm">
            <div class="pick-confirm-main">
              <div class="pick-confirm-name">将切换为 <strong>{{ pendingSpecies.name }}</strong></div>
              <div v-if="petPickerStudent.pet_name && petPickerStudent.free_pick" class="pick-cost-free">🎉 免费窗口期内，本次切换免费！</div>
              <div v-else-if="petPickerStudent.pet_name && getSwitchCost(petPickerStudent) > 0" class="pick-cost-charge">💰 本次切换扣除 <strong>{{ getSwitchCost(petPickerStudent) }}</strong> 积分 · 保留当前等级</div>
              <div v-else class="pick-cost-free">🎉 首次免费，不扣积分</div>
              <div v-if="switchError" class="switch-error">{{ switchError }}</div>
            </div>
            <div class="pick-confirm-actions">
              <button class="pc-cancel" @click="clearPending" :disabled="switchingPet">取消选择</button>
              <button class="pc-confirm" @click="executeSwitch" :disabled="switchingPet" :style="switchStatus === 'loading' ? 'background:#f59e0b;color:#fff' : switchStatus === 'success' ? 'background:#10b981;color:#fff' : switchStatus === 'error' ? 'background:#ef4444;color:#fff' : ''">
                <template v-if="switchStatus === 'loading'">切换中...</template>
                <template v-else-if="switchStatus === 'success'">切换成功 ✓</template>
                <template v-else-if="switchStatus === 'error'">切换失败 ✗</template>
                <template v-else>确认切换</template>
              </button>
            </div>
          </div>
          <button class="picker-cancel" @click="showPetPicker = false">取消</button>
        </div>
      </div>
    </Transition>

    <!-- 浮动积分 -->
    <Teleport to="body">
      <div v-for="f in floatTexts" :key="f.id" class="float-text" :style="{ left: f.x + 'px', top: f.y + 'px', color: f.color }">
        {{ f.text }}
      </div>
    </Teleport>
  </div>

  <!-- 教师端：最近积分记录 -->
  <div v-if="isTeacherMode" class="card" style="margin-top:24px;">
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px;">
      <h3 style="font-size:16px;font-weight:600;">📋 最近积分记录</h3>
      <button class="btn btn-sm" :disabled="historyLoading" @click="loadRecentScores">🔄 刷新</button>
    </div>
    <div v-if="historyLoading" style="text-align:center;padding:24px;color:var(--color-text-secondary);">加载中...</div>
    <div v-else-if="recentScores.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">暂无记录</div>
    <div v-else class="data-table">
      <table>
        <thead><tr><th>学生</th><th>分值</th><th>原因</th><th>时间</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="s in recentScores" :key="s.id">
            <td style="font-weight:600;">{{ s.student_name }}</td>
            <td :style="{ color: s.amount > 0 ? '#10B981' : '#EF4444', fontWeight: 700 }">{{ s.amount > 0 ? '+' : '' }}{{ s.amount }}</td>
            <td style="color:var(--color-text-secondary);">{{ s.reason }}</td>
            <td style="color:var(--color-text-secondary);font-size:13px;">{{ new Date(s.created_at).toLocaleString('zh-CN') }}</td>
            <td><button class="btn btn-xs" @click="undoScore(s.id)" :disabled="getUndoStatus(s.id) === 'loading'" :style="getUndoBtnStyle(s.id)">{{ getUndoBtnText(s.id) }}</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.scores-page {
  max-width: 1200px;
}
.load-error {
  margin-bottom: 12px;
  padding: 8px 12px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  border-radius: 8px;
  color: var(--color-danger-text);
  font-size: 12px;
}

/* 顶部 */
.page-top { margin-bottom: 20px; }
.page-header { display: flex; align-items: baseline; gap: 12px; margin-bottom: 12px; flex-wrap: wrap; }
.page-title { font-size: 24px; font-weight: 700; margin: 0; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); }
.class-total {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 14px;
  border-radius: 20px;
  background: var(--tint-2);
  border: 1px solid var(--tint-3);
  font-size: 12px;
  color: var(--md-text-secondary);
}
.ct-icon { font-size: 14px; }
.ct-value { font-size: 16px; font-weight: 800; color: var(--color-text); }

.stats-row { display: flex; gap: 8px; flex-wrap: wrap; }
.stat-chip {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
}
.chip-icon { font-size: 14px; }

/* 工具栏（教师端） */
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; flex-wrap: wrap; }
.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: 30px;
  flex: 1;
  min-width: 200px;
  max-width: 360px;
}
.search-icon { font-size: 16px; }
.search-input { background: transparent; border: none; outline: none; color: var(--color-text); font-size: 14px; width: 100%; font-family: inherit; }
.search-input::placeholder { color: var(--color-text-secondary); opacity: 0.6; }
.filter-group { display: flex; gap: 6px; }
.filter-tag {
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
  transition: all 0.2s ease;
}
.filter-tag:hover { color: var(--color-text); border-color: var(--color-text-secondary); }
.filter-tag.active { background: rgba(79,70,229,0.08); border-color: var(--color-primary); color: var(--color-primary); font-weight: 600; }
.sort-group { display: flex; align-items: center; gap: 6px; }
.sort-label { font-size: 12px; color: var(--color-text-secondary); white-space: nowrap; }
.sort-tag {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
  transition: all 0.2s ease;
}
.sort-tag:hover { color: var(--color-text); border-color: var(--color-text-secondary); }
.sort-tag.active { background: rgba(16,185,129,0.08); border-color: #10B981; color: #10B981; font-weight: 600; }

/* 教室端搜索 */
.class-search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--tint-2);
  border: 1px solid var(--tint-3);
  border-radius: 30px;
  flex: 1;
  min-width: 200px;
  max-width: 360px;
}
.class-search input { background: transparent; border: none; outline: none; color: var(--color-text); font-size: 14px; width: 100%; font-family: inherit; }

/* 教室端认养横幅 */
.adopt-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding: 14px 20px;
  background: linear-gradient(135deg, rgba(167,139,250,0.08), rgba(244,114,182,0.05));
  border: 1px solid rgba(167,139,250,0.15);
  border-radius: var(--md-radius);
  flex-wrap: wrap;
}
.adopt-icon { font-size: 24px; }
.adopt-title { font-weight: 700; font-size: 15px; }
.adopt-desc { font-size: 13px; color: var(--md-text-secondary); }

/* 批量操作栏 */
.batch-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding: 10px 14px;
  background: linear-gradient(135deg, rgba(79,70,229,0.06), rgba(16,185,129,0.05));
  border: 1px solid var(--color-border);
  border-radius: 14px;
  flex-wrap: wrap;
  animation: modalPop 0.2s ease;
}
.batch-info { font-size: 13px; color: var(--color-text); }
.batch-info strong { color: var(--color-primary); }
.batch-btn { padding: 6px 16px; border-radius: 20px; font-size: 13px; font-weight: 600; cursor: pointer; border: 1px solid transparent; transition: all 0.15s ease; font-family: inherit; }
.batch-add { background: rgba(16,185,129,0.1); color: #10B981; border-color: rgba(16,185,129,0.25); }
.batch-add:hover { background: rgba(16,185,129,0.18); }
.batch-sub { background: rgba(239,68,68,0.1); color: #EF4444; border-color: rgba(239,68,68,0.25); }
.batch-sub:hover { background: rgba(239,68,68,0.18); }
.batch-clear { margin-left: auto; background: transparent; color: var(--color-text-secondary); border-color: var(--color-border); }
.batch-clear:hover { background: var(--color-bg); }

/* 快捷规则条 */
.quick-rules {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding: 10px 14px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  flex-wrap: wrap;
}
.qr-label { font-size: 12px; font-weight: 600; color: var(--color-text-secondary); white-space: nowrap; }
.qr-group { display: flex; gap: 6px; flex-wrap: wrap; }
.qr-btn { padding: 4px 12px; border-radius: 16px; font-size: 12px; font-weight: 600; cursor: pointer; border: 1px solid transparent; transition: all 0.15s ease; }
.qr-add { background: rgba(16,185,129,0.08); color: #10B981; border-color: rgba(16,185,129,0.2); }
.qr-add:hover { background: rgba(16,185,129,0.15); }
.qr-sub { background: rgba(239,68,68,0.08); color: #EF4444; border-color: rgba(239,68,68,0.2); }
.qr-sub:hover { background: rgba(239,68,68,0.15); }

/* 学生卡片网格 */
.student-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}
.empty-grid { grid-column: 1 / -1; text-align: center; padding: 60px 24px; color: var(--color-text-secondary); font-size: 16px; }

/* 浮动文字 */
.float-text {
  position: fixed;
  pointer-events: none;
  font-size: 24px;
  font-weight: 800;
  z-index: 999;
  animation: floatUp 1.2s ease-out forwards;
  transform: translateX(-50%);
}
@keyframes floatUp {
  0% { opacity: 1; transform: translateX(-50%) translateY(0) scale(1); }
  100% { opacity: 0; transform: translateX(-50%) translateY(-80px) scale(1.3); }
}
@keyframes modalPop {
  from { transform: scale(0.92); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

/* 加载 */
.loading-state { text-align: center; padding: 60px 24px; color: var(--color-text-secondary); }
.loading-spinner {
  width: 36px; height: 36px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* 教室端宠物切换/选择器覆盖层 */
.overlay-high {
  position: fixed;
  inset: 0;
  z-index: 301;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
/* 内联确认条：选中宠物后在弹窗内确认，不弹第二个弹窗 */
.pick-confirm {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
  padding: 12px 14px;
  background: var(--tint-1);
  border: 1px solid var(--tint-3);
  border-radius: var(--md-radius);
  flex-wrap: wrap;
  animation: modalPop 0.2s ease;
}
.pick-confirm-main { flex: 1; min-width: 0; }
.pick-confirm-name { font-size: 13px; color: var(--color-text); }
.pick-confirm-name strong { color: var(--color-primary); }
.pick-cost-free { font-size: 12px; margin-top: 4px; padding: 3px 10px; border-radius: 8px; font-weight: 700; display: inline-block; background: rgba(16,185,129,0.12); border: 1px solid rgba(16,185,129,0.3); color: var(--color-success-text); }
.pick-cost-charge { font-size: 12px; margin-top: 4px; padding: 3px 10px; border-radius: 8px; font-weight: 700; display: inline-block; background: rgba(245,158,11,0.12); border: 1px solid rgba(245,158,11,0.3); color: var(--color-warning-text); }
.pick-confirm-actions { display: flex; gap: 8px; flex-shrink: 0; }
.pc-cancel { padding: 8px 16px; border-radius: 10px; border: 1px solid var(--tint-3); background: transparent; color: var(--md-text-secondary); font-size: 13px; cursor: pointer; font-family: inherit; }
.pc-confirm {
  padding: 8px 18px;
  border-radius: 10px;
  border: none;
  background: var(--md-primary);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s ease;
}
.pc-confirm:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 4px 14px rgba(0,0,0,0.2); }
.pc-confirm:disabled { opacity: 0.6; cursor: not-allowed; }
.switch-error { margin-top: 8px; padding: 8px 12px; background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2); border-radius: 8px; color: var(--color-danger-text); font-size: 12px; }

.picker-box {
  background: var(--color-bg-card);
  border: 1px solid var(--tint-3);
  border-radius: var(--md-radius);
  padding: 24px 28px;
  max-width: 520px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: var(--md-elevation);
  animation: modalPop 0.25s ease;
}
.picker-head { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.picker-emoji { font-size: 28px; }
.picker-name { font-size: 16px; font-weight: 700; }
.pick-badge { font-size: 12px; margin-top: 4px; padding: 4px 10px; border-radius: 8px; font-weight: 700; display: inline-block; }
.pick-free { background: rgba(16,185,129,0.12); border: 1px solid rgba(16,185,129,0.3); color: var(--color-success-text); }
.pick-cost { background: rgba(245,158,11,0.12); border: 1px solid rgba(245,158,11,0.3); color: var(--color-warning-text); }
.picker-close { margin-left: auto; width: 28px; height: 28px; border-radius: 50%; border: 1px solid var(--tint-3); background: transparent; color: var(--color-text-secondary); cursor: pointer; }
.picker-series { margin-bottom: 12px; }
.series-name { font-size: 12px; font-weight: 600; color: var(--md-text-secondary); margin-bottom: 6px; padding-left: 4px; }
.series-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(90px, 1fr)); gap: 6px; }
.species-btn { padding: 8px 4px; border-radius: 10px; border: 1px solid var(--tint-2); background: var(--tint-1); text-align: center; cursor: pointer; transition: 0.15s; font-family: inherit; }
.species-btn:hover { background: var(--tint-3); }
.species-sprite { width: 48px; height: 48px; margin: 0 auto 2px; }
.species-name { font-size: 10px; font-weight: 500; color: var(--md-text); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.species-current { font-size: 9px; font-weight: 700; color: #10B981; }
.picker-cancel { width: 100%; margin-top: 8px; padding: 8px; border-radius: 10px; border: 1px solid var(--tint-3); background: transparent; color: var(--md-text-secondary); font-size: 13px; cursor: pointer; font-family: inherit; }

/* 过渡 */
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* 最近积分记录操作按钮（紧凑尺寸） */
.btn-xs { padding: 4px 12px; font-size: 12px; border-radius: 16px; background: rgba(239, 68, 68, 0.08); border: 1px solid rgba(239, 68, 68, 0.2); }
.btn-xs:hover { background: rgba(239, 68, 68, 0.15); }

@media (max-width: 768px) {
  .student-grid { grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); }
  .toolbar { flex-direction: column; align-items: stretch; }
  .search-box { max-width: none; }
}
@media (max-width: 576px) {
  .student-grid { grid-template-columns: 1fr; }
}
</style>
