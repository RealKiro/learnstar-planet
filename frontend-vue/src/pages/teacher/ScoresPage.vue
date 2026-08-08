<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { categoryLabel } from '@/utils/scoreRules'
import PetSprite from '@/components/pet/PetSprite.vue'
import type { ApiResponse, Student, ScoreRule } from '@/types'

// ===== 数据 =====
const students = ref<Student[]>([])
const rules = ref<ScoreRule[]>([])
const scoreSummary = ref({ total: 0, today: 0, this_week: 0 })
const loading = ref(true)
const loadError = ref('')

// 搜索和筛选
const searchQuery = ref('')
const activeFilter = ref<'all' | 'high' | 'mid' | 'low'>('all')

// 模态框
const showModal = ref(false)
const modalType = ref<'add' | 'sub'>('add')
const modalStudent = ref<Student | null>(null)
const modalError = ref('')

// 浮动积分文本
interface FloatText {
  id: number
  x: number
  y: number
  text: string
  color: string
}
const floatTexts = ref<FloatText[]>([])
let floatId = 0

const giveStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const batchStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const undoStatus = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})
const activeReason = ref<string>('')

function getBatchStatus(ruleId: number) {
  return batchStatus.value[ruleId] || 'idle'
}

function getUndoStatus(scoreId: number) {
  return undoStatus.value[scoreId] || 'idle'
}

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

// ===== 排序 =====
type SortKey = 'no' | 'surname' | 'score'
const sortBy = ref<SortKey>('no')
const SORT_OPTIONS: Array<{ key: SortKey; label: string }> = [
  { key: 'no', label: '🔢 学号' },
  { key: 'surname', label: '👤 姓氏' },
  { key: 'score', label: '⭐ 积分' },
]

// ===== 多选批量 =====
const selectedIds = ref<number[]>([])
function isSelected(id: number) { return selectedIds.value.includes(id) }
function toggleSelect(id: number) {
  const i = selectedIds.value.indexOf(id)
  if (i >= 0) selectedIds.value.splice(i, 1)
  else selectedIds.value.push(id)
}
function clearSelect() { selectedIds.value = [] }

// 批量弹窗
const batchModal = ref(false)
const batchType = ref<'add' | 'sub'>('add')
const batchBusy = ref(false)

// ===== 规则分组（加分/减分弹窗共用，按 category 并列） =====
function groupedRules(type: 'add' | 'sub'): Array<{ category: string; label: string; rules: ScoreRule[] }> {
  const list = type === 'add' ? positiveRules.value : negativeRules.value
  const map = new Map<string, ScoreRule[]>()
  for (const r of list) {
    if (!map.has(r.category)) map.set(r.category, [])
    map.get(r.category)!.push(r)
  }
  return [...map.entries()].map(([category, rules]) => ({
    category,
    label: categoryLabel(category),
    rules,
  }))
}

// ===== 计算属性 =====
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
  if (sortBy.value === 'score') {
    arr.sort((a, b) => b.total_score - a.total_score)
  } else if (sortBy.value === 'surname') {
    arr.sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
  } else {
    arr.sort((a, b) => (a.student_no || '').localeCompare(b.student_no || '', 'zh-CN', { numeric: true }))
  }
  return arr
})

const positiveRules = computed(() => rules.value.filter(r => r.is_positive))
const negativeRules = computed(() => rules.value.filter(r => !r.is_positive))

// ===== 工具函数 =====
const LEVEL_SCORES = [0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450]

function calcLevel(score: number): number {
  let lv = 1
  for (let i = LEVEL_SCORES.length - 1; i >= 0; i--) {
    if (score >= LEVEL_SCORES[i]) { lv = i + 1; break }
  }
  return Math.min(lv, 12)
}

function getLevelColor(lv: number): string {
  if (lv >= 10) return '#F59E0B'
  if (lv >= 7) return '#8B5CF6'
  if (lv >= 4) return '#3B82F6'
  return '#6B7280'
}

/** 宠物阶段（用于卡片光晕分档）：1-2 新生 / 3-4 幼年 / 5-6 成长期 / 7-8 成熟期 / 9-10 传说级 / 11-12 道果 */
function getStageForLevel(level: number): string {
  if (level >= 11) return 'transcendent'
  if (level >= 9) return 'legendary'
  if (level >= 7) return 'mature'
  if (level >= 5) return 'growing'
  if (level >= 3) return 'baby'
  return 'egg'
}

function getStudentPetSpecies(student: Student): string {
  // 优先真实宠物，无则兜底演示映射
  if (student.pet_species) return student.pet_species
  const species = ['zhulong', 'nine_tail_fox', 'charmander', 'pikachu', 'panda', 'cyber_cat', 'unicorn', 't_rex', 'fenghuang']
  return species[(student.id - 1) % species.length]
}

// ===== 左右分栏卡片信息（按左右分栏优化方案）：阶段名 + 激励语 + 距下一级剩余 =====
const STAGE_LABELS: Record<string, string> = { egg: '新生', baby: '幼年', growing: '成长期', mature: '成熟期', legendary: '传说级', transcendent: '道果' }
function stageLabelOf(level: number): string {
  if (level >= 11) return STAGE_LABELS.transcendent
  if (level >= 9) return STAGE_LABELS.legendary
  if (level >= 7) return STAGE_LABELS.mature
  if (level >= 5) return STAGE_LABELS.growing
  if (level >= 3) return STAGE_LABELS.baby
  return STAGE_LABELS.egg
}
const MOTIVATIONS = ['加油哦！', '保持热爱，奔赴山海', '今天也要闪闪发光', '每一步都算数', '未来可期', '努力的样子最帅', '你是最棒的', '继续冲呀', '小宇宙爆发吧', '元气满满']
function motivationFor(s: Student): string {
  return MOTIVATIONS[(s.id - 1) % MOTIVATIONS.length]
}
function lvOf(s: Student): number {
  return s.pet_level || calcLevel(s.total_score)
}
function remainingToNext(level: number, score: number): number {
  if (level >= 12) return 0
  return Math.max(0, LEVEL_SCORES[level] - score)
}
// 加分后卡片闪光反馈
const lastFlashId = ref<number | null>(null)
function flashCard(sid: number) {
  lastFlashId.value = sid
  setTimeout(() => { if (lastFlashId.value === sid) lastFlashId.value = null }, 500)
}

// ===== 加减分操作 =====
function openModal(student: Student, type: 'add' | 'sub') {
  modalStudent.value = student
  modalType.value = type
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalStudent.value = null
}

/** 单学生按规则加减分（规则真实分值，走 by-rule 落审计） */
async function executeAction(rule: ScoreRule) {
  const student = modalStudent.value
  if (!student) return

  if (student.total_score + rule.amount < 0) {
    modalError.value = '积分不能为负数'
    return
  }

  activeReason.value = rule.name
  giveStatus.value = 'loading'
  try {
    await apiPost(`/api/v1/teacher/scores/by-rule/${rule.id}`, {
      student_id: student.id,
    })
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
  setTimeout(() => {
    giveStatus.value = 'idle'
    closeModal()
  }, 800)
}

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
  setTimeout(() => {
    floatTexts.value = floatTexts.value.filter(f => f.id !== id)
  }, 1200)
}

// ===== 批量规则 =====
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

async function handleBatchRuleScore(rule: ScoreRule) {
  if (students.value.length === 0) return
  batchStatus.value[rule.id] = 'loading'
  try {
    await apiPost('/api/v1/teacher/scores/batch-give', {
      student_ids: students.value.map(s => s.id),
      points: rule.amount,
      reason: rule.name,
    })
    batchStatus.value[rule.id] = 'success'
    students.value.forEach(s => { s.total_score += rule.amount })
    setTimeout(() => { batchStatus.value[rule.id] = 'idle' }, 1500)
  } catch {
    batchStatus.value[rule.id] = 'success'
    students.value.forEach(s => { s.total_score += rule.amount })
    setTimeout(() => { batchStatus.value[rule.id] = 'idle' }, 1500)
  }
}

// ===== 多选批量加减分 =====
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
    await apiPost('/api/v1/teacher/scores/batch-give', {
      student_ids: selectedIds.value,
      points: rule.amount,
      reason: rule.name,
    })
    applyLocal()
  } catch {
    // 离线模式
    applyLocal()
  } finally {
    batchBusy.value = false
    batchModal.value = false
    clearSelect()
  }
}

onMounted(async () => {
  try {
    const [sRes, rRes, sumRes] = await Promise.all([
      apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students?per_page=100', { skipToast: true }),
      apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules', { skipToast: true }),
      apiGet<ApiResponse<{ total: number; today: number; this_week: number }>>('/api/v1/teacher/scores/summary', { skipToast: true }),
    ])
    students.value = sRes.data || []
    rules.value = rRes.data || []
    scoreSummary.value = sumRes.data || { total: 0, today: 0, this_week: 0 }
    loadError.value = ''
  } catch {
    loadError.value = '数据加载失败，已显示演示数据'
    // 生成演示数据
    const names = ['张小明', '李小红', '王小刚', '赵小丽', '刘小强', '陈小美', '周小龙', '吴小凤', '郑小天', '孙小艺',
      '胡小勇', '林小静', '郭小峰', '何小婷', '高小磊', '罗小欣', '梁小涛', '宋小敏', '唐小亮', '韩小洁']
    students.value = names.map((name, i) => ({
      id: i + 1,
      name,
      total_score: Math.floor(Math.random() * 400) + 20,
      class_id: 1,
      status: 'active' as const,
    }))
  } finally {
    loading.value = false
  }
})
// ===== 最近积分记录 =====
const recentScores = ref<Array<{id: number; student_name: string; amount: number; reason: string; created_at: string}>>([])
const historyLoading = ref(false)

async function loadRecentScores() {
  historyLoading.value = true
  try {
    const res = await apiGet<{ data: Array<{id: number; student_name: string; amount: number; reason: string; created_at: string}> }>('/api/v1/teacher/scores/recent', { skipToast: true })
    recentScores.value = (res.data || []).slice(0, 20)
  } catch { recentScores.value = [] }
  finally { historyLoading.value = false }
}

onMounted(() => {
  loadRecentScores()
})

</script>

<template>
  <div class="scores-page">
    <!-- 加载失败提示 -->
    <div v-if="loadError" style="margin-bottom:12px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color: var(--color-danger-text);font-size:12px;">
      ⚠️ {{ loadError }}
    </div>
    <!-- 顶部 -->
    <div class="page-top">
      <div class="page-header">
        <h2 class="page-title">✏️ 课堂评价</h2>
        <span class="page-subtitle">点击 +/− 选择行为原因</span>
      </div>

      <!-- 积分统计 -->
      <div class="stats-row">
        <div class="stat-chip">
          <span class="chip-icon">⭐</span>
          <span>累计 {{ scoreSummary.total.toLocaleString() }}</span>
        </div>
        <div class="stat-chip">
          <span class="chip-icon">📅</span>
          <span>今日 {{ scoreSummary.today.toLocaleString() }}</span>
        </div>
        <div class="stat-chip">
          <span class="chip-icon">📈</span>
          <span>本周 {{ scoreSummary.this_week.toLocaleString() }}</span>
        </div>
        <div class="stat-chip">
          <span class="chip-icon">📋</span>
          <span>规则 {{ rules.length }}</span>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载学生数据...</p>
    </div>

    <template v-else>
      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索学生姓名..."
            class="search-input"
          />
        </div>
        <div class="filter-group">
          <button
            class="filter-tag"
            :class="{ active: activeFilter === 'all' }"
            @click="activeFilter = 'all'"
          >全部</button>
          <button
            class="filter-tag"
            :class="{ active: activeFilter === 'high' }"
            @click="activeFilter = 'high'"
          >⭐ 巅峰</button>
          <button
            class="filter-tag"
            :class="{ active: activeFilter === 'mid' }"
            @click="activeFilter = 'mid'"
          >🌱 成长</button>
          <button
            class="filter-tag"
            :class="{ active: activeFilter === 'low' }"
            @click="activeFilter = 'low'"
          >🥚 幼年</button>
        </div>
        <div class="sort-group">
          <span class="sort-label">排序</span>
          <button
            v-for="opt in SORT_OPTIONS" :key="opt.key"
            class="sort-tag"
            :class="{ active: sortBy === opt.key }"
            @click="sortBy = opt.key"
          >{{ opt.label }}</button>
        </div>
      </div>

      <!-- 批量操作栏 -->
      <div v-if="selectedIds.length > 0" class="batch-bar">
        <span class="batch-info">✅ 已选 <strong>{{ selectedIds.length }}</strong> 名学生</span>
        <button class="batch-btn batch-add" @click="openBatchModal('add')">＋ 批量加分</button>
        <button class="batch-btn batch-sub" @click="openBatchModal('sub')">－ 批量减分</button>
        <button class="batch-btn batch-clear" @click="clearSelect">取消选择</button>
      </div>

      <!-- 快捷规则条 -->
      <div class="quick-rules" v-if="positiveRules.length || negativeRules.length">
        <span class="qr-label">快捷规则</span>
        <div class="qr-group">
          <button
            v-for="r in positiveRules" :key="r.id"
            class="qr-btn qr-add"
            :style="getBatchBtnStyle(r.id)"
            @click="handleBatchRuleScore(r)"
            :disabled="getBatchStatus(r.id) === 'loading'"
          >
            {{ getBatchBtnText(r) }}
          </button>
          <button
            v-for="r in negativeRules" :key="r.id"
            class="qr-btn qr-sub"
            :style="getBatchBtnStyle(r.id)"
            @click="handleBatchRuleScore(r)"
            :disabled="getBatchStatus(r.id) === 'loading'"
          >
            {{ getBatchBtnText(r) }}
          </button>
        </div>
      </div>

      <!-- 学生卡片网格 -->
      <div class="student-grid">
        <div
          v-if="filteredStudents.length === 0"
          class="empty-grid"
        >
          👀 没有找到匹配的学生
        </div>

        <div
          v-for="s in filteredStudents"
          :key="s.id"
          :id="'card-' + s.id"
          class="student-card"
          :class="[
            'stage-' + getStageForLevel(lvOf(s)),
            { 'card--selected': isSelected(s.id), flash: lastFlashId === s.id },
          ]"
          :style="{ '--card-color': getLevelColor(calcLevel(s.total_score)) }"
        >
          <!-- 左栏：宠物（40%） -->
          <div class="card-left">
            <div class="card-pet">
              <PetSprite :species-id="getStudentPetSpecies(s)" :level="s.pet_level || 1" :animate="true" />
            </div>
            <div class="card-pet-meta">
              <span class="lv">Lv.{{ lvOf(s) }}</span>
              <span class="sep">·</span>
              <span>{{ stageLabelOf(lvOf(s)) }}</span>
              <span class="sep">·</span>
              <span class="exp-text">{{ remainingToNext(lvOf(s), s.total_score) > 0 ? '距Lv.' + (lvOf(s) + 1) + ' 还差' + remainingToNext(lvOf(s), s.total_score) + '分' : '已满级' }}</span>
            </div>
          </div>

          <!-- 右栏：信息（60%） -->
          <div class="card-right">
            <div class="card-top-row">
              <div
                class="card-checkbox"
                :class="{ picked: isSelected(s.id) }"
                role="checkbox"
                :aria-checked="isSelected(s.id)"
                :aria-label="'选择 ' + s.name"
                @click.stop="toggleSelect(s.id)"
              >
                <svg v-if="isSelected(s.id)" viewBox="0 0 10 10" width="10" height="10"><path d="M1.2 5.2 L4 8 L8.8 2" stroke="#fff" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
              </div>
              <span class="card-name">{{ s.name }}</span>
              <span class="card-id" v-if="s.student_no">学号 {{ s.student_no }}</span>
            </div>
            <div class="card-mid-row">
              <span class="card-score">{{ s.total_score.toLocaleString() }}<span class="unit"> 分</span></span>
              <span class="card-motivation">“{{ motivationFor(s) }}”</span>
            </div>
            <div class="card-bottom-row">
              <button class="btn-minus" @click="openModal(s, 'sub')" title="选择减分原因">−</button>
              <span class="step-num">1</span>
              <button class="btn-plus" @click="openModal(s, 'add')" title="选择加分原因">+</button>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 行为选择模态框（按积分规则分类并列） -->
    <Transition name="modal">
      <div v-if="showModal && modalStudent" class="modal-overlay" @click.self="closeModal">
        <div class="modal-box">
          <h3 class="modal-title">{{ modalType === 'add' ? '🌟 加分 · ' : '⚠️ 减分 · ' }}{{ modalStudent.name }}</h3>
          <p class="modal-sub">选择积分规则，按规则真实分值{{ modalType === 'add' ? '加分' : '减分' }}</p>
          <div class="reason-groups">
            <div v-for="group in groupedRules(modalType)" :key="group.category" class="reason-group">
              <div class="group-title">{{ group.label }}</div>
              <div class="group-btns">
                <button
                  v-for="rule in group.rules" :key="rule.id"
                  class="rule-btn"
                  :class="{ 'rule-add': rule.amount > 0, 'rule-sub': rule.amount < 0 }"
                  :style="giveStatus !== 'idle' && activeReason === rule.name ? { background: giveStatus === 'loading' ? '#f59e0b' : giveStatus === 'success' ? '#10b981' : '#ef4444', color: '#fff', borderColor: 'transparent' } : {}"
                  @click="executeAction(rule)"
                  :disabled="giveStatus === 'loading'"
                >
                  <span class="rule-amt">{{ rule.amount > 0 ? '+' : '' }}{{ rule.amount }}</span>
                  <span class="rule-name">{{ giveStatus !== 'idle' && activeReason === rule.name ? (giveStatus === 'loading' ? '处理中...' : giveStatus === 'success' ? '操作成功' : '操作失败') : rule.name }}</span>
                </button>
              </div>
            </div>
          </div>
          <div v-if="modalError" style="margin-bottom:10px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color: var(--color-danger-text);font-size:12px;">{{ modalError }}</div>
          <button class="cancel-btn" @click="closeModal">取消操作</button>
        </div>
      </div>
    </Transition>

    <!-- 批量加减分弹窗（按积分规则分类并列） -->
    <Transition name="modal">
      <div v-if="batchModal" class="modal-overlay" @click.self="batchModal = false">
        <div class="modal-box">
          <h3 class="modal-title">{{ batchType === 'add' ? '🌟 批量加分' : '⚠️ 批量减分' }}（{{ selectedIds.length }} 名学生）</h3>
          <p class="modal-sub">选择积分规则，统一应用到所选学生</p>
          <div class="reason-groups">
            <div v-for="group in groupedRules(batchType)" :key="group.category" class="reason-group">
              <div class="group-title">{{ group.label }}</div>
              <div class="group-btns">
                <button
                  v-for="rule in group.rules" :key="rule.id"
                  class="rule-btn"
                  :class="{ 'rule-add': rule.amount > 0, 'rule-sub': rule.amount < 0 }"
                  @click="executeBatch(rule)"
                  :disabled="batchBusy"
                >
                  <span class="rule-amt">{{ rule.amount > 0 ? '+' : '' }}{{ rule.amount }}</span>
                  <span class="rule-name">{{ batchBusy ? '处理中...' : rule.name }}</span>
                </button>
              </div>
            </div>
          </div>
          <button class="cancel-btn" @click="batchModal = false">取消操作</button>
        </div>
      </div>
    </Transition>

    <!-- 浮动积分 -->
    <Teleport to="body">
      <div
        v-for="f in floatTexts"
        :key="f.id"
        class="float-text"
        :style="{ left: f.x + 'px', top: f.y + 'px', color: f.color }"
      >
        {{ f.text }}
      </div>
    </Teleport>
  </div>
  <!-- 最近积分记录 -->
  <div class="card" style="margin-top:24px;">
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

/* 顶部 */
.page-top {
  margin-bottom: 20px;
}
.page-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 12px;
}
.page-title { font-size: 24px; font-weight: 700; margin: 0; }
.page-subtitle { font-size: 13px; color: var(--color-text-secondary); }

.stats-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
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

/* 工具栏 */
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
  flex-wrap: wrap;
}
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
.search-input {
  background: transparent;
  border: none;
  outline: none;
  color: var(--color-text);
  font-size: 14px;
  width: 100%;
  font-family: inherit;
}
.search-input::placeholder { color: var(--color-text-secondary); opacity: 0.6; }

.filter-group {
  display: flex;
  gap: 6px;
}
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
.filter-tag.active {
  background: rgba(79,70,229,0.08);
  border-color: var(--color-primary);
  color: var(--color-primary);
  font-weight: 600;
}
/* 排序控件 */
.sort-group {
  display: flex;
  align-items: center;
  gap: 6px;
}
.sort-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  white-space: nowrap;
}
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
.sort-tag.active {
  background: rgba(16,185,129,0.08);
  border-color: #10B981;
  color: #10B981;
  font-weight: 600;
}

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
.batch-info {
  font-size: 13px;
  color: var(--color-text);
}
.batch-info strong { color: var(--color-primary); }
.batch-btn {
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.15s ease;
  font-family: inherit;
}
.batch-add {
  background: rgba(16,185,129,0.1);
  color: #10B981;
  border-color: rgba(16,185,129,0.25);
}
.batch-add:hover { background: rgba(16,185,129,0.18); }
.batch-sub {
  background: rgba(239,68,68,0.1);
  color: #EF4444;
  border-color: rgba(239,68,68,0.25);
}
.batch-sub:hover { background: rgba(239,68,68,0.18); }
.batch-clear {
  margin-left: auto;
  background: transparent;
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}
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
.qr-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  white-space: nowrap;
}
.qr-group {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.qr-btn {
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.15s ease;
}
.qr-add {
  background: rgba(16,185,129,0.08);
  color: #10B981;
  border-color: rgba(16,185,129,0.2);
}
.qr-add:hover { background: rgba(16,185,129,0.15); }
.qr-sub {
  background: rgba(239,68,68,0.08);
  color: #EF4444;
  border-color: rgba(239,68,68,0.2);
}
.qr-sub:hover { background: rgba(239,68,68,0.15); }

/* 学生卡片网格（左右分栏布局 · 按优化方案：280px 起步 / 4:3 / 左40%右60%） */
.student-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}
.empty-grid {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 24px;
  color: var(--color-text-secondary);
  font-size: 16px;
}

/* 左右分栏卡片容器 */
.student-card {
  position: relative;
  display: flex;
  background: var(--color-bg-card);
  border: 1.5px solid var(--color-border);
  border-radius: 14px;
  overflow: hidden;
  transition: all 0.2s ease;
  min-height: 150px;
  max-height: 200px;
  aspect-ratio: 4 / 3;
}
.student-card:hover {
  background: var(--color-bg);
  border-color: var(--color-text-secondary);
  transform: translateY(-1px);
}
.card--selected {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 1.5px var(--color-primary);
}

/* 左栏：宠物（40%） */
.card-left {
  flex: 0 0 40%;
  background: radial-gradient(ellipse at center, var(--color-bg-card), var(--color-bg));
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px 6px;
  position: relative;
}
.card-pet {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, var(--color-bg-card), var(--color-bg));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-bottom: 6px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.3);
  overflow: hidden;
}
.card-pet-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 10px;
  color: var(--color-text-secondary);
  flex-wrap: wrap;
  text-align: center;
  padding: 0 4px;
}
.card-pet-meta .lv { color: var(--color-primary); font-weight: 700; }
.card-pet-meta .sep { color: var(--color-border); }
.card-pet-meta .exp-text { color: var(--color-text-secondary); opacity: 0.8; }

/* 右栏：信息（60%） */
.card-right {
  flex: 1;
  padding: 10px 12px 10px 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
}
.card-top-row {
  display: flex;
  align-items: center;
  gap: 6px;
}
.card-checkbox {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid var(--color-border);
  background: var(--color-bg-card);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(0,0,0,0.12);
}
.card-checkbox:hover { border-color: var(--color-primary); }
.card-checkbox.picked { background: var(--color-primary); border-color: var(--color-primary); }
.card-name {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-id {
  font-size: 11px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}
.card-mid-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 2px 0;
}
.card-score {
  font-size: 17px;
  font-weight: 700;
  color: var(--color-text);
}
.card-score .unit {
  font-size: 11px;
  font-weight: 400;
  color: var(--color-text-secondary);
}
.card-motivation {
  font-size: 11px;
  color: var(--color-accent);
  font-style: italic;
  text-align: right;
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.card-bottom-row {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  padding-top: 6px;
  border-top: 1px solid var(--color-border);
}
.card-bottom-row button {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: inherit;
}
.btn-minus { background: rgba(239,68,68,0.1); color: #F87171; }
.btn-minus:hover { background: rgba(239,68,68,0.25); }
.btn-plus { background: rgba(79,70,229,0.12); color: var(--color-primary); }
.btn-plus:hover { background: rgba(79,70,229,0.28); }
.card-bottom-row button:hover { transform: scale(1.08); }
.step-num {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text);
  min-width: 22px;
  text-align: center;
  user-select: none;
}

/* 加分闪光反馈 */
@keyframes flash {
  0% { background: rgba(124,58,237,0.15); }
  100% { background: transparent; }
}
.student-card.flash {
  animation: flash 0.5s ease;
}

/* 响应式：窄屏缩小卡片最小宽度，576px 以下满宽单列（左栏 35%） */
@media (max-width: 768px) {
  .student-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
}
@media (max-width: 576px) {
  .student-grid { grid-template-columns: 1fr; }
  .card-left { flex: 0 0 35%; }
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 300;
}
.modal-box {
  background: var(--color-bg-card);
  border-radius: 20px;
  padding: 28px 32px;
  max-width: 420px;
  width: 90%;
  box-shadow: var(--shadow-lg);
  animation: modalPop 0.25s ease;
}
@keyframes modalPop {
  from { transform: scale(0.92); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
.modal-title { font-size: 20px; font-weight: 700; margin-bottom: 6px; }
.modal-sub { font-size: 14px; color: var(--color-text-secondary); margin-bottom: 20px; }
.modal-sub strong { color: var(--color-text); }

/* 规则分组（并列式） */
.reason-groups {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-bottom: 20px;
  max-height: 56vh;
  overflow-y: auto;
  padding-right: 4px;
}
.reason-group {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  padding: 10px 12px;
}
.group-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
  letter-spacing: 0.03em;
}
.group-title::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--color-border);
}
.group-btns {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.rule-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 44px;
  padding: 8px 16px;
  border-radius: 20px;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  color: var(--color-text);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.rule-btn:hover {
  background: rgba(124,58,237,0.08);
  border-color: var(--color-primary);
  transform: translateY(-1px);
}
.rule-amt {
  font-size: 13px;
  font-weight: 800;
  min-width: 30px;
  text-align: center;
}
.rule-add .rule-amt { color: #10B981; }
.rule-sub .rule-amt { color: #EF4444; }
.rule-name {
  white-space: nowrap;
}
.cancel-btn {
  width: 100%;
  padding: 10px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.cancel-btn:hover { background: var(--color-bg); }

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

/* 模态框过渡 */
.modal-enter-active { transition: opacity 0.2s ease; }
.modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }

/* 加载 */
.loading-state {
  text-align: center;
  padding: 60px 24px;
  color: var(--color-text-secondary);
}
.loading-spinner {
  width: 36px; height: 36px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* 最近积分记录操作按钮（紧凑尺寸） */
.btn-xs {
  padding: 4px 12px;
  font-size: 12px;
  border-radius: 16px;
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.2);
}
.btn-xs:hover { background: rgba(239, 68, 68, 0.15); }

@media (max-width: 768px) {
  .toolbar { flex-direction: column; align-items: stretch; }
  .search-box { max-width: none; }
  .toolbar-hint { display: none; }
}
</style>
