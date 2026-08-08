<script setup lang="ts">
import { computed } from 'vue'
import PetSprite from '@/components/pet/PetSprite.vue'

/** 卡片学生（教师 Student / 教室 StudentEntry 的结构兼容子集） */
export interface CardStudent {
  id: number
  name: string
  student_no?: string
  total_score: number
  pet_level?: number
  pet_species?: string
  pet_name?: string
  pet_emoji?: string
  /** 整班切系列后 3 天免费自选窗口 */
  free_pick?: boolean
  class_id?: number
  status?: string
}

const props = withDefaults(defineProps<{
  student: CardStudent
  isTeacherMode: boolean
  selected?: boolean
  flash?: boolean
  /** 教室端：当前卡片的步长是否处于编辑态 */
  editingStep?: boolean
  /** 教室端：当前卡片步长值 */
  stepValue?: number
  /** 教室端：编辑输入框的值 */
  editInput?: string
}>(), {
  selected: false,
  flash: false,
  editingStep: false,
  stepValue: 1,
  editInput: '1',
})

const emit = defineEmits<{
  'toggle-select': []
  'open-detail': []
  'open-modal': [type: 'add' | 'sub']
  'start-edit': []
  'save-edit': [event?: Event]
  'update:editInput': [value: string]
}>()

// ===== 等级体系（教师端：12 级） =====
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
function getStageForLevel(level: number): string {
  if (level >= 11) return 'transcendent'
  if (level >= 9) return 'legendary'
  if (level >= 7) return 'mature'
  if (level >= 5) return 'growing'
  if (level >= 3) return 'baby'
  return 'egg'
}
function lvOf(s: CardStudent): number {
  return s.pet_level || calcLevel(s.total_score)
}
function remainingToNext(level: number, score: number): number {
  if (level >= 12) return 0
  return Math.max(0, LEVEL_SCORES[level] - score)
}
function teacherSpecies(s: CardStudent): string {
  if (s.pet_species) return s.pet_species
  const species = ['zhulong', 'nine_tail_fox', 'charmander', 'pikachu', 'panda', 'cyber_cat', 'unicorn', 't_rex', 'fenghuang']
  return species[(s.id - 1) % species.length]
}

// ===== 等级体系（教室端：52 级） =====
const CLASS_LEVEL_SCORES = [0, 0, 15, 41, 68, 96, 125, 155, 185, 217, 250, 283, 318, 353, 390, 427, 465, 504, 545, 586, 628, 671, 715, 760, 805, 852, 900, 949, 998, 1049, 1100, 1153, 1206, 1261, 1316, 1372, 1429, 1487, 1546, 1606, 1667, 1729, 1792, 1856, 1921, 1986, 2053, 2120, 2189, 2258, 2329, 2400, 99999]

function nextLevelProgress(score: number, level: number): { remaining: number } {
  const maxLevel = CLASS_LEVEL_SCORES.length - 2
  if (level >= maxLevel || score >= CLASS_LEVEL_SCORES[CLASS_LEVEL_SCORES.length - 2]) {
    return { remaining: 0 }
  }
  const next = CLASS_LEVEL_SCORES[Math.min(level + 1, maxLevel)] || 450
  return { remaining: Math.max(0, next - score) }
}

// ===== 阶段名 + 激励语（共用） =====
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
function motivationFor(s: CardStudent): string {
  return MOTIVATIONS[(s.id - 1) % MOTIVATIONS.length]
}

const lv = computed(() => lvOf(props.student))

/** 距下一级剩余文案（两端公式不同，模式分支） */
const expText = computed(() => {
  const s = props.student
  if (props.isTeacherMode) {
    const r = remainingToNext(lv.value, s.total_score)
    return r > 0 ? `距Lv.${lv.value + 1} 还差${r}分` : '已满级'
  }
  const l = s.pet_level || 1
  const r = nextLevelProgress(s.total_score, l).remaining
  return r > 0 ? `距Lv.${l + 1} 还差${r}分` : '已满级'
})

/** 距下一级进度百分比（两端公式不同） */
const expPercent = computed(() => {
  const s = props.student
  if (props.isTeacherMode) {
    if (lv.value >= 12) return 100
    const cur = LEVEL_SCORES[lv.value - 1]
    const next = LEVEL_SCORES[lv.value]
    return Math.min(100, Math.max(0, ((s.total_score - cur) / (next - cur)) * 100))
  }
  const l = s.pet_level || 1
  const maxLevel = CLASS_LEVEL_SCORES.length - 2
  if (l >= maxLevel || s.total_score >= CLASS_LEVEL_SCORES[CLASS_LEVEL_SCORES.length - 2]) return 100
  const cur = CLASS_LEVEL_SCORES[l] || 0
  const next = CLASS_LEVEL_SCORES[Math.min(l + 1, maxLevel)] || 450
  return Math.min(100, Math.max(0, ((s.total_score - cur) / (next - cur)) * 100))
})
</script>

<template>
  <div
    :id="'card-' + student.id"
    class="student-card"
    :class="[
      'stage-' + getStageForLevel(lv),
      { 'card--selected': selected, flash },
    ]"
    :style="{ '--card-color': getLevelColor(calcLevel(student.total_score)) }"
  >
    <!-- 上部：左宠物 + 右信息（5:5） -->
    <div class="card-top">
      <!-- 左栏：宠物（50%） -->
      <div class="card-left">
        <div
          class="card-pet"
          @click="!isTeacherMode && emit('open-detail')"
          :title="isTeacherMode ? undefined : (student.pet_species ? '点击查看宠物详情' : '未孵化 · 点击认养')"
        >
          <PetSprite
            v-if="isTeacherMode || student.pet_species"
            :species-id="isTeacherMode ? teacherSpecies(student) : student.pet_species!"
            :level="student.pet_level || 1"
            :animate="true"
          />
          <span v-else class="pet-emoji">{{ student.pet_emoji || '🥚' }}</span>
        </div>
        <div class="card-pet-name" v-if="student.pet_name">{{ student.pet_name }}</div>
      </div>

      <!-- 右栏：信息（50%） -->
      <div class="card-right">
        <div class="card-name-row">
          <span class="card-name">{{ student.name }}</span>
          <div
            class="card-checkbox"
            :class="{ picked: selected }"
            role="checkbox"
            :aria-checked="selected"
            :aria-label="'选择 ' + student.name"
            @click.stop="emit('toggle-select')"
            title="多选后批量加减分"
          >
            <svg v-if="selected" viewBox="0 0 10 10" width="10" height="10"><path d="M1.2 5.2 L4 8 L8.8 2" stroke="#fff" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </div>
        </div>
        <div class="card-id-row">
          <span class="card-id-icon">📛</span>
          <span v-if="student.student_no" class="card-id">{{ student.student_no }}</span>
          <span v-else class="card-id card-id--none">未编学号</span>
        </div>
        <div class="card-mid-row">
          <span class="card-score"><span class="score-icon">⭐</span>{{ student.total_score.toLocaleString() }}<span class="unit"> 分</span></span>
          <span class="card-motivation">“{{ motivationFor(student) }}”</span>
        </div>
      </div>
    </div>

    <!-- 底部：左宠物等级/进度 + 右加减分（同层对齐） -->
    <div class="card-bottom">
      <div class="card-pet-meta">
        <div class="meta-line">
          <span class="lv">Lv.{{ lv }}</span>
          <span class="sep">·</span>
          <span class="stage">{{ stageLabelOf(lv) }}</span>
          <span class="exp-text">{{ expText }}</span>
        </div>
        <div class="mini-progress">
          <div class="mini-fill" :style="{ width: expPercent + '%' }"></div>
        </div>
      </div>

      <div class="card-ops">
        <button class="btn-minus" @click="emit('open-modal', 'sub')" :title="isTeacherMode ? '选择减分原因' : undefined">−</button>
        <input
          v-if="!isTeacherMode && editingStep"
          :value="editInput"
          type="number"
          min="1"
          max="100"
          autofocus
          @input="emit('update:editInput', ($event.target as HTMLInputElement).value)"
          @blur="emit('save-edit')"
          @keydown.enter="emit('save-edit', $event)"
          class="step-input"
        />
        <span
          v-else
          class="step-num"
          :class="{ 'step-editable': !isTeacherMode }"
          @click="!isTeacherMode && emit('start-edit')"
          :title="isTeacherMode ? undefined : '点击修改步长'"
        >{{ isTeacherMode ? 1 : stepValue }}</span>
        <button class="btn-plus" @click="emit('open-modal', 'add')" :title="isTeacherMode ? '选择加分原因' : undefined">+</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ===== 卡片容器（两端共用 · 5:5 分栏） ===== */
.student-card {
  position: relative;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-card);
  border: 1.5px solid var(--color-border);
  border-radius: 14px;
  overflow: hidden;
  transition: all 0.2s ease;
  min-height: 170px;
  max-height: 210px;
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

/* ===== 上部：左宠物 + 右信息 ===== */
.card-top {
  display: flex;
  flex: 1;
  min-height: 0;
}
.card-left {
  flex: 1 1 50%;
  width: 50%;
  background: radial-gradient(ellipse at center, var(--color-bg-card), var(--color-bg));
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 6px;
  position: relative;
  gap: 4px;
}
.card-pet {
  position: relative;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, var(--color-bg-card), var(--color-bg));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 16px rgba(0,0,0,0.3);
  overflow: hidden;
}
.card-pet:hover { transform: scale(1.05); }
.pet-emoji { font-size: 42px; line-height: 1; }
.card-pet-name {
  font-size: 11px;
  color: var(--color-text-secondary);
  max-width: 92%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
/* 宠物图(AI/hand SVG <img>)强制锁定圆形容器尺寸，防止占满页面 */
.student-card :deep(.pet-sprite),
.student-card :deep(.pet-img) {
  width: 96px;
  height: 96px;
  max-width: 96px;
  max-height: 96px;
  object-fit: contain;
  flex-shrink: 0;
}

/* ===== 右栏：信息（50%） ===== */
.card-right {
  flex: 1 1 50%;
  width: 50%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 7px;
  padding: 10px 12px 8px;
  min-width: 0;
}
.card-name-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.card-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.card-id-row {
  display: flex;
  align-items: center;
  gap: 5px;
}
.card-id-icon { font-size: 11px; }
.card-id {
  font-size: 11px;
  color: var(--color-text-secondary);
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  padding: 1px 8px;
  border-radius: 8px;
  letter-spacing: 0.3px;
}
.card-id--none { opacity: 0.5; }

.card-mid-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
}
.card-score {
  font-size: 16px;
  font-weight: 800;
  color: var(--color-text);
  white-space: nowrap;
}
.score-icon { font-size: 12px; margin-right: 2px; }
.card-score .unit {
  font-size: 11px;
  font-weight: 400;
  color: var(--color-text-secondary);
}
.card-motivation {
  font-size: 11px;
  color: var(--color-accent);
  font-style: italic;
  flex: 1;
  min-width: 0;
  text-align: right;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ===== 底部：宠物等级/进度 + 加减分（同层对齐） ===== */
.card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 6px 12px 7px;
  border-top: 1px solid var(--color-border);
  background: var(--color-bg);
}
.card-pet-meta {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}
.meta-line {
  display: flex;
  align-items: baseline;
  gap: 4px;
  font-size: 10px;
  color: var(--color-text-secondary);
  min-width: 0;
}
.meta-line .lv { color: var(--color-primary); font-weight: 700; }
.meta-line .sep { color: var(--color-border); }
.meta-line .stage { font-weight: 600; }
.meta-line .exp-text {
  margin-left: auto;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--color-text-secondary);
  opacity: 0.8;
}
.mini-progress {
  height: 4px;
  width: 100%;
  background: var(--color-border);
  border-radius: 2px;
  overflow: hidden;
}
.mini-fill {
  height: 100%;
  background: var(--gradient-primary);
  border-radius: 2px;
  transition: width 0.4s ease;
}

.card-ops {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}
.card-ops button {
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
.card-ops button:hover { transform: scale(1.08); }
.step-num {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text);
  min-width: 22px;
  text-align: center;
  user-select: none;
}
.step-editable { cursor: pointer; }
.step-input {
  width: 40px;
  text-align: center;
  background: rgba(167,139,250,0.15);
  border: 1px solid rgba(167,139,250,0.3);
  border-radius: 8px;
  color: var(--color-text);
  font-size: 14px;
  font-weight: 700;
  outline: none;
  font-family: inherit;
}

/* ===== 加分闪光反馈 ===== */
@keyframes flash {
  0% { background: rgba(124,58,237,0.15); }
  100% { background: transparent; }
}
.student-card.flash {
  animation: flash 0.5s ease;
}

/* ===== 阶段光晕 ===== */
.stage-egg .card-pet { box-shadow: 0 0 10px rgba(148,163,184,0.18); }
.stage-baby .card-pet { box-shadow: 0 0 12px rgba(16,185,129,0.18); }
.stage-growing .card-pet { box-shadow: 0 0 14px rgba(59,130,246,0.22); }
.stage-mature .card-pet { box-shadow: 0 0 18px rgba(139,92,246,0.26); }
.stage-legendary .card-pet { box-shadow: 0 0 24px rgba(245,158,11,0.30); }
.stage-transcendent .card-pet { box-shadow: 0 0 28px rgba(216,180,254,0.34), 0 0 46px rgba(255,255,255,0.06); }

@media (max-width: 576px) {
  .card-pet { width: 76px; height: 76px; }
  .student-card :deep(.pet-sprite),
  .student-card :deep(.pet-img) {
    width: 76px; height: 76px;
    max-width: 76px; max-height: 76px;
  }
}
</style>
