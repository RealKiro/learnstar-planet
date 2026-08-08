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
const clsRemaining = computed(() => nextLevelProgress(props.student.total_score, props.student.pet_level || 1).remaining)
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
    <!-- 左栏：宠物（40%） -->
    <div class="card-left">
      <!-- 教师端：固定宠物（含 demo 兜底），无点击 -->
      <template v-if="isTeacherMode">
        <div class="card-pet">
          <PetSprite :species-id="teacherSpecies(student)" :level="student.pet_level || 1" :animate="true" />
        </div>
        <div class="card-pet-meta">
          <span class="lv">Lv.{{ lv }}</span>
          <span class="sep">·</span>
          <span>{{ stageLabelOf(lv) }}</span>
          <span class="sep">·</span>
          <span class="exp-text">{{ remainingToNext(lv, student.total_score) > 0 ? '距Lv.' + (lv + 1) + ' 还差' + remainingToNext(lv, student.total_score) + '分' : '已满级' }}</span>
        </div>
      </template>
      <!-- 教室端：真实宠物可点击开详情 / 未孵化可认养 -->
      <template v-else>
        <div class="card-pet" @click="emit('open-detail')" title="点击查看宠物详情">
          <PetSprite v-if="student.pet_species" :species-id="student.pet_species" :level="student.pet_level || 1" :animate="true" />
          <span v-else style="font-size:34px;line-height:1;">{{ student.pet_emoji || '🥚' }}</span>
        </div>
        <div class="card-pet-meta">
          <span class="lv">Lv.{{ student.pet_level || 1 }}</span>
          <span class="sep">·</span>
          <span>{{ stageLabelOf(student.pet_level || 1) }}</span>
          <span class="sep">·</span>
          <span class="exp-text">{{ clsRemaining > 0 ? '距Lv.' + ((student.pet_level || 1) + 1) + ' 还差' + clsRemaining + '分' : '已满级' }}</span>
        </div>
        <span v-if="!student.pet_species" class="unhatched">未孵化 · 点击认养</span>
      </template>
    </div>

    <!-- 右栏：信息（60%，两端共用） -->
    <div class="card-right">
      <div class="card-top-row">
        <div
          class="card-checkbox"
          :class="{ picked: selected }"
          role="checkbox"
          :aria-checked="selected"
          :aria-label="'选择 ' + student.name"
          @click.stop="emit('toggle-select')"
        >
          <svg v-if="selected" viewBox="0 0 10 10" width="10" height="10"><path d="M1.2 5.2 L4 8 L8.8 2" stroke="#fff" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>
        <span class="card-name">{{ student.name }}</span>
        <span class="card-id" v-if="student.student_no">学号 {{ student.student_no }}</span>
      </div>
      <div class="card-mid-row">
        <span class="card-score">{{ student.total_score.toLocaleString() }}<span class="unit"> 分</span></span>
        <span class="card-motivation">“{{ motivationFor(student) }}”</span>
      </div>
      <!-- 底行：教师固定 1；教室可编辑步长 -->
      <div class="card-bottom-row">
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
/* 左右分栏卡片容器（两端共用 · 按左右分栏优化方案） */
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
  position: relative;
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
  cursor: pointer;
}
/* 宠物图(AI/hand SVG <img>)强制锁定圆形容器尺寸，防止占满页面 */
.student-card :deep(.pet-sprite),
.student-card :deep(.pet-img) {
  width: 64px;
  height: 64px;
  max-width: 64px;
  max-height: 64px;
  object-fit: contain;
  flex-shrink: 0;
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
.unhatched { font-size: 10px; color: var(--color-text-secondary); }

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

/* 加分闪光反馈 */
@keyframes flash {
  0% { background: rgba(124,58,237,0.15); }
  100% { background: transparent; }
}
.student-card.flash {
  animation: flash 0.5s ease;
}

/* 阶段光晕 */
.stage-egg .card-pet { box-shadow: 0 0 10px rgba(148,163,184,0.18); }
.stage-baby .card-pet { box-shadow: 0 0 12px rgba(16,185,129,0.18); }
.stage-growing .card-pet { box-shadow: 0 0 14px rgba(59,130,246,0.22); }
.stage-mature .card-pet { box-shadow: 0 0 18px rgba(139,92,246,0.26); }
.stage-legendary .card-pet { box-shadow: 0 0 24px rgba(245,158,11,0.30); }
.stage-transcendent .card-pet { box-shadow: 0 0 28px rgba(216,180,254,0.34), 0 0 46px rgba(255,255,255,0.06); }

@media (max-width: 576px) {
  .card-left { flex: 0 0 35%; }
}
</style>
