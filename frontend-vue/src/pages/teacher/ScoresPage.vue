<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { getSpeciesEmoji, getSeriesBySpeciesId } from '@/utils/petData'
import type { ApiResponse, Student, ScoreRule } from '@/types'

const toast = useToastStore()

// ===== 数据 =====
const students = ref<Student[]>([])
const rules = ref<ScoreRule[]>([])
const scoreSummary = ref({ total: 0, today: 0, this_week: 0 })
const loading = ref(true)

// 搜索和筛选
const searchQuery = ref('')
const activeFilter = ref<'all' | 'high' | 'mid' | 'low'>('all')

// 模态框
const showModal = ref(false)
const modalType = ref<'add' | 'sub'>('add')
const modalStudent = ref<Student | null>(null)
const modalReasons = ref<string[]>([])

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

const REASONS_ADD = ['📖 举手发言', '✅ 作业优秀', '🤝 帮助同学', '🧹 遵守纪律', '🏆 挑战难题', '📝 认真笔记']
const REASONS_SUB = ['⚠️ 上课走神', '📕 作业缺交', '🗣️ 打扰课堂', '🏃 追逐打闹']

// ===== 计算属性 =====
const filteredStudents = computed(() => {
  return students.value.filter(s => {
    const matchName = s.name.includes(searchQuery.value)
    const level = calcLevel(s.total_score)
    if (activeFilter.value === 'high') return matchName && level >= 10
    if (activeFilter.value === 'mid') return matchName && level >= 4 && level <= 9
    if (activeFilter.value === 'low') return matchName && level <= 3
    return matchName
  })
})

const positiveRules = computed(() => rules.value.filter(r => !r.is_penalty))
const negativeRules = computed(() => rules.value.filter(r => r.is_penalty))

// ===== 工具函数 =====
const LEVEL_SCORES = [0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450]

function calcLevel(score: number): number {
  let lv = 1
  for (let i = LEVEL_SCORES.length - 1; i >= 0; i--) {
    if (score >= LEVEL_SCORES[i]) { lv = i + 1; break }
  }
  return Math.min(lv, 12)
}

function calcProgress(score: number): number {
  const lv = calcLevel(score)
  if (lv >= 12) return 1
  const current = LEVEL_SCORES[lv - 1]
  const next = LEVEL_SCORES[lv]
  return (score - current) / (next - current)
}

function getLevelColor(lv: number): string {
  if (lv >= 10) return '#F59E0B'
  if (lv >= 7) return '#8B5CF6'
  if (lv >= 4) return '#3B82F6'
  return '#6B7280'
}

function getStudentPetSpecies(student: Student): string {
  // 根据学生ID映射不同宠物（演示用）
  const species = ['zhulong', 'nine_tail_fox', 'charmander', 'pikachu', 'panda', 'cyber_cat', 'unicorn', 't_rex', 'nian', 'fenghuang']
  return species[(student.id - 1) % species.length]
}

// ===== 步长编辑 =====
const editingStep = ref<Record<number, boolean>>({})

function startEditStep(studentId: number) {
  editingStep.value[studentId] = true
}

function finishEditStep(studentId: number) {
  editingStep.value[studentId] = false
}

// ===== 加减分操作 =====
function openModal(student: Student, type: 'add' | 'sub') {
  modalStudent.value = student
  modalType.value = type
  modalReasons.value = type === 'add' ? REASONS_ADD : REASONS_SUB
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalStudent.value = null
}

async function executeAction(reason: string) {
  const student = modalStudent.value
  if (!student) return

  const step = getStudentStep(student.id)
  const points = modalType.value === 'add' ? step : -step

  if (student.total_score + points < 0) {
    toast.show('积分不能为负数', 'error')
    return
  }

  try {
    await apiPost('/api/v1/teacher/scores/give', {
      student_id: student.id,
      points,
      reason,
    })
    student.total_score = Math.max(0, student.total_score + points)
    toast.show(`${student.name} ${points > 0 ? '+' : ''}${points}分 — ${reason}`, 'success')
    showFloatText(student.id, points)
  } catch {
    // 离线模式
    student.total_score = Math.max(0, student.total_score + points)
    toast.show(`${student.name} ${points > 0 ? '+' : ''}${points}分 — ${reason}`, 'success')
    showFloatText(student.id, points)
  }
  closeModal()
}

function getStudentStep(studentId: number): number {
  return 1 // default step
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
async function handleBatchRuleScore(rule: ScoreRule) {
  if (students.value.length === 0) return
  try {
    await apiPost('/api/v1/teacher/scores/batch-give', {
      student_ids: students.value.map(s => s.id),
      points: rule.points,
      reason: rule.name,
    })
    toast.show(`全班按规则「${rule.name}」${rule.is_penalty ? '扣' : '加'}${Math.abs(rule.points)}分`, 'success')
    students.value.forEach(s => { s.total_score += rule.points })
  } catch {
    toast.show(`全班按规则「${rule.name}」${rule.is_penalty ? '扣' : '加'}${Math.abs(rule.points)}分`, 'success')
    students.value.forEach(s => { s.total_score += rule.points })
  }
}

onMounted(async () => {
  try {
    const [sRes, rRes, sumRes] = await Promise.all([
      apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students?per_page=100'),
      apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules'),
      apiGet<ApiResponse<{ total: number; today: number; this_week: number }>>('/api/v1/teacher/scores/summary'),
    ])
    students.value = sRes.data || []
    rules.value = rRes.data || []
    scoreSummary.value = sumRes.data || { total: 0, today: 0, this_week: 0 }
  } catch {
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
</script>

<template>
  <div class="scores-page">
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
        <span class="toolbar-hint">点击数字修改步长</span>
      </div>

      <!-- 快捷规则条 -->
      <div class="quick-rules" v-if="positiveRules.length || negativeRules.length">
        <span class="qr-label">快捷规则</span>
        <div class="qr-group">
          <button
            v-for="r in positiveRules" :key="r.id"
            class="qr-btn qr-add"
            @click="handleBatchRuleScore(r)"
          >
            +{{ r.points }} {{ r.name }}
          </button>
          <button
            v-for="r in negativeRules" :key="r.id"
            class="qr-btn qr-sub"
            @click="handleBatchRuleScore(r)"
          >
            {{ r.points }} {{ r.name }}
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
          :style="{ '--card-color': getLevelColor(calcLevel(s.total_score)) }"
        >
          <!-- 顶部 -->
          <div class="card-top">
            <span class="card-name">{{ s.name }}</span>
            <span class="card-level" :style="{ background: getLevelColor(calcLevel(s.total_score)) + '22', color: getLevelColor(calcLevel(s.total_score)) }">
              Lv.{{ calcLevel(s.total_score) }}
            </span>
          </div>

          <!-- 宠物信息 -->
          <div class="card-pet">
            <span class="pet-icon">{{ getSpeciesEmoji(getStudentPetSpecies(s)) }}</span>
            <span class="pet-text">宠物伙伴</span>
          </div>

          <!-- 积分 -->
          <div class="card-score">
            <span class="score-label">⭐ 积分</span>
            <span class="score-num" :style="{ color: getLevelColor(calcLevel(s.total_score)) }">
              {{ s.total_score }}
            </span>
          </div>

          <!-- 进度条 -->
          <div class="card-progress">
            <div
              class="progress-fill"
              :style="{ width: Math.min(calcProgress(s.total_score) * 100, 100) + '%' }"
            ></div>
          </div>

          <!-- 操作按钮 -->
          <div class="card-actions">
            <button class="action-btn btn-sub" @click="openModal(s, 'sub')">
              <span>−</span>
            </button>
            <span class="step-value">1</span>
            <button class="action-btn btn-add" @click="openModal(s, 'add')">
              <span>+</span>
            </button>
          </div>
        </div>
      </div>
    </template>

    <!-- 行为选择模态框 -->
    <Transition name="modal">
      <div v-if="showModal && modalStudent" class="modal-overlay" @click.self="closeModal">
        <div class="modal-box">
          <h3 class="modal-title">{{ modalType === 'add' ? '🌟 选择加分行为' : '⚠️ 选择减分原因' }}</h3>
          <p class="modal-sub">
            为 <strong>{{ modalStudent.name }}</strong> 选择原因
            （每次{{ modalType === 'add' ? '加' : '减' }} <strong>1</strong> 分）
          </p>
          <div class="reason-grid">
            <button
              v-for="reason in modalReasons"
              :key="reason"
              class="reason-btn"
              @click="executeAction(reason)"
            >
              {{ reason }}
            </button>
          </div>
          <button class="cancel-btn" @click="closeModal">取消操作</button>
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
.toolbar-hint {
  font-size: 11px;
  color: var(--color-text-secondary);
  opacity: 0.5;
  margin-left: auto;
}

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

/* 学生卡片网格 */
.student-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 14px;
}
.empty-grid {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 24px;
  color: var(--color-text-secondary);
  font-size: 16px;
}

.student-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 18px;
  padding: 16px 14px 12px;
  transition: all 0.25s ease;
  position: relative;
  overflow: hidden;
}
.student-card::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: var(--card-color, #6B7280);
  border-radius: 18px 0 0 18px;
}
.student-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}
.card-name {
  font-size: 16px;
  font-weight: 700;
}
.card-level {
  font-size: 12px;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 12px;
}

.card-pet {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}
.pet-icon { font-size: 20px; }
.pet-text { font-size: 13px; color: var(--color-text-secondary); }

.card-score {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}
.score-label { font-size: 12px; color: var(--color-text-secondary); }
.score-num { font-size: 18px; font-weight: 800; }

.card-progress {
  height: 4px;
  background: var(--color-border);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 10px;
}
.progress-fill {
  height: 100%;
  background: var(--gradient-primary);
  border-radius: 2px;
  transition: width 0.4s ease;
}

.card-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--color-border);
}
.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 1px solid var(--color-border);
  background: transparent;
  font-size: 20px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  color: var(--color-text-secondary);
}
.action-btn:hover { transform: scale(1.12); }
.btn-add {
  color: #10B981;
  border-color: rgba(16,185,129,0.2);
  background: rgba(16,185,129,0.04);
}
.btn-add:hover { background: rgba(16,185,129,0.12); }
.btn-sub {
  color: #EF4444;
  border-color: rgba(239,68,68,0.2);
  background: rgba(239,68,68,0.04);
}
.btn-sub:hover { background: rgba(239,68,68,0.12); }
.step-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text);
  min-width: 28px;
  text-align: center;
  cursor: pointer;
  padding: 2px 8px;
  border-radius: 12px;
  transition: background 0.15s;
}
.step-value:hover { background: var(--color-bg); }

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

.reason-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
}
.reason-btn {
  padding: 12px 16px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  color: var(--color-text);
  font-size: 15px;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.reason-btn:hover {
  background: rgba(79,70,229,0.04);
  border-color: var(--color-primary);
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

@media (max-width: 768px) {
  .student-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  }
  .toolbar { flex-direction: column; align-items: stretch; }
  .search-box { max-width: none; }
  .toolbar-hint { display: none; }
}
</style>
