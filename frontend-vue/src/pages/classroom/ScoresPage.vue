<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet, apiPost } from '@/utils/api'
import { getSpeciesEmoji, PET_SERIES } from '@/utils/petData'
import { getEvoLines } from '@/utils/petHandbookData'
import PetSprite from '@/components/pet/PetSprite.vue'
import PetDetailModal from '@/components/pet/PetDetailModal.vue'

const router = useRouter()

// 等级所需积分常量
const LEVEL_SCORES = [0, 0, 15, 41, 68, 96, 125, 155, 185, 217, 250, 283, 318, 353, 390, 427, 465, 504, 545, 586, 628, 671, 715, 760, 805, 852, 900, 949, 998, 1049, 1100, 1153, 1206, 1261, 1316, 1372, 1429, 1487, 1546, 1606, 1667, 1729, 1792, 1856, 1921, 1986, 2053, 2120, 2189, 2258, 2329, 2400, 99999]

function nextLevelProgress(score: number, level: number): { current: number; next: number; remaining: number; percent: number } {
  const maxLevel = LEVEL_SCORES.length - 2
  if (level >= maxLevel || score >= LEVEL_SCORES[LEVEL_SCORES.length - 2]) {
    return { current: score, next: score, remaining: 0, percent: 100 }
  }
  const cur = LEVEL_SCORES[level] || 0
  const next = LEVEL_SCORES[Math.min(level + 1, maxLevel)] || 450
  const exp = score - cur
  const range = next - cur
  return {
    current: Math.max(0, exp),
    next: range,
    remaining: Math.max(0, range - exp),
    percent: Math.min(100, range > 0 ? Math.round((exp / range) * 100) : 100),
  }
}

interface StudentEntry {
  id: number; name: string; student_no: string; total_score: number
  pet_name: string; pet_species: string; pet_level: number; pet_emoji: string
  free_pick?: boolean // 整班切系列后 3 天免费自选窗口
}

const students = ref<StudentEntry[]>([])
const loading = ref(true)
const token = ref('')
const searchQuery = ref('')
const showModal = ref(false)
const modalStudent = ref<StudentEntry | null>(null)
const modalType = ref<'add' | 'sub'>('add')
const actionError = ref('')
const stepValues = ref<Record<number, number>>({})
const editingStep = ref<number | null>(null)
const editInput = ref('1')
// 加减分理由按分类分组（按视觉优化方案 5.3：学习/行为/品德 分组，各带 emoji + 色标）
interface ReasonGroup { title: string; emoji: string; color: string; items: string[] }
const reasonGroupsAdd: ReasonGroup[] = [
  { title: '学习表现', emoji: '📚', color: '#3B82F6', items: ['✅ 作业优秀', '🏆 挑战难题', '📚 阅读之星', '🎨 科技创新'] },
  { title: '行为表现', emoji: '🏅', color: '#10B981', items: ['📖 举手发言', '🤝 帮助同学', '📝 认真听讲', '🔍 专注课堂'] },
  { title: '品德修养', emoji: '⭐', color: '#F59E0B', items: ['🧹 遵守纪律', '💬 积极互动', '🌟 诚实守信', '🏃 体育锻炼'] },
]
const reasonGroupsSub: ReasonGroup[] = [
  { title: '学习懈怠', emoji: '📕', color: '#EF4444', items: ['⚠️ 上课走神', '📕 作业缺交'] },
  { title: '课堂纪律', emoji: '🗣️', color: '#F97316', items: ['🗣️ 打扰课堂', '📱 课堂喧哗', '💬 说脏话'] },
  { title: '行为品德', emoji: '🤕', color: '#F59E0B', items: ['🏃 追逐打闹', '😴 趴桌睡觉', '⚡ 与同学冲突', '🗑️ 乱扔垃圾'] },
]
const reasonGroups = computed(() => modalType.value === 'add' ? reasonGroupsAdd : reasonGroupsSub)
const batchReasonGroups = computed(() => batchType.value === 'add' ? reasonGroupsAdd : reasonGroupsSub)
// 弹窗多选（单学生 / 批量共用）
const selectedReasons = ref<string[]>([])
const customReason = ref('')
const showCustom = ref(false)
function toggleReason(r: string) {
  const i = selectedReasons.value.indexOf(r)
  if (i >= 0) selectedReasons.value.splice(i, 1)
  else selectedReasons.value.push(r)
}
function resetReasonPick() { selectedReasons.value = []; customReason.value = ''; showCustom.value = false }
function selectedList(): string[] {
  const list = [...selectedReasons.value]
  if (showCustom.value && customReason.value.trim()) list.push('✍️ ' + customReason.value.trim())
  return list
}
const floatTexts = ref<Array<{ id: number; x: number; y: number; text: string; color: string }>>([])
let floatId = 0

// 批量加减分：多选学生（挖空圆形）
const selectedIds = ref<number[]>([])
const batchModal = ref(false)
const batchType = ref<'add' | 'sub'>('add')
const batchError = ref('')
const batchBusy = ref(false)

function isSelected(sid: number) { return selectedIds.value.includes(sid) }
function toggleSelect(sid: number) {
  const i = selectedIds.value.indexOf(sid)
  if (i >= 0) selectedIds.value.splice(i, 1)
  else selectedIds.value.push(sid)
}
function clearSelect() { selectedIds.value = [] }
function openBatchModal(type: 'add' | 'sub') { batchType.value = type; batchError.value = ''; resetReasonPick(); batchModal.value = true }

async function confirmBatch() {
  const ids = selectedIds.value
  if (!ids.length || batchBusy.value) return
  const reasons = selectedList()
  if (!reasons.length) return
  const step = getStep(ids[0])
  const points = batchType.value === 'add' ? step * reasons.length : -step * reasons.length
  batchBusy.value = true
  try {
    await apiPost('/api/v1/display/scores/batch-give', { token: token.value, student_ids: ids, points, reason: reasons.join('、') })
    for (const s of students.value) {
      if (ids.includes(s.id)) {
        s.total_score = Math.max(0, s.total_score + points)
        showFloatText(s.id, points)
      }
    }
    batchError.value = ''
    resetReasonPick()
    batchModal.value = false
    clearSelect()
  } catch (e: any) {
    const msg = e?.response?.data?.message || ''
    batchError.value = msg.includes('30') ? '单次超过 30 分需要教师账号登录操作' : (msg || '操作失败，请稍后重试')
    setTimeout(() => { batchError.value = '' }, 3000)
  } finally {
    batchBusy.value = false
  }
}

// 宠物切换
const showPetPicker = ref(false)
const petPickerStudent = ref<StudentEntry | null>(null)
const switchingPet = ref(false)
const switchStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const switchError = ref('')
// 宠物选择器：班级配置了当前系列(pet_series)时只展示该系列，未配置时展示全部（与后端类别限制一致）
const classPetSeries = ref('')
const pickerSeriesList = computed(() => {
  if (!classPetSeries.value) return PET_SERIES
  const hit = PET_SERIES.find(s => s.id === classPetSeries.value)
  return hit ? [hit] : PET_SERIES
})

function openPetPicker(s: StudentEntry) {
  petPickerStudent.value = s
  showPetPicker.value = true
}

const confirmSwitch = ref<{ speciesId: string; name: string } | null>(null)

// 宠物详情弹窗
const showPetDetail = ref(false)
const detailStudent = ref<StudentEntry | null>(null)

function openPetDetail(s: StudentEntry) {
  detailStudent.value = s
  showPetDetail.value = true
}

function requestSwitch(speciesId: string, name: string) {
  confirmSwitch.value = { speciesId, name }
}

// 选择宠物：当前宠物不可再切换（避免触发后端"当前已经是这只宠物啦"）
function handlePick(speciesId: string, name: string) {
  if (petPickerStudent.value?.pet_species === speciesId) return
  requestSwitch(speciesId, name)
}

async function executeSwitch() {
  if (!confirmSwitch.value) return
  const s = petPickerStudent.value
  if (!s || switchingPet.value) return
  const speciesId = confirmSwitch.value.speciesId
  confirmSwitch.value = null
  switchingPet.value = true
  switchStatus.value = 'loading'
  try {
    const res = await apiPost<{ data: { pet_emoji: string; total_score: number; cost: number } }>(
      '/api/v1/display/pets/switch', { token: token.value, student_id: s.id, pet_species: speciesId }
    )
    s.pet_species = speciesId
    s.pet_emoji = res.data.pet_emoji
    s.total_score = res.data.total_score
    switchStatus.value = 'success'
    switchError.value = ''
    setTimeout(() => { switchStatus.value = 'idle' }, 1500)
    showPetPicker.value = false
  } catch (e: any) {
    switchStatus.value = 'error'
    switchError.value = e?.response?.data?.message || '切换失败，请稍后重试'
    setTimeout(() => {
      switchStatus.value = 'idle'
      switchError.value = ''
    }, 3000)
  } finally {
    switchingPet.value = false
  }
}

function getStep(sid: number) { return stepValues.value[sid] || 1 }

// 切换宠物所需积分：首次免费；整班切系列后 3 天免费窗口内免费；否则按当前宠物等级扣（后端 switchCost = 5×等级）
function getSwitchCost(s: StudentEntry | null): number {
  if (!s || !s.pet_name) return 0
  if (s.free_pick) return 0
  return 5 * Math.max(1, s.pet_level || 1)
}

const filtered = computed(() => {
  if (!searchQuery.value) return students.value
  return students.value.filter(s => s.name.includes(searchQuery.value))
})

// ===== 主卡片（视觉优化方案 5.1）：宠物主视觉 + 功能入口 + 班级总分 =====
const spotlightStudent = computed(() => {
  const withPet = students.value.filter(s => s.pet_species)
  return [...withPet].sort((a, b) => b.total_score - a.total_score)[0] || withPet[0] || null
})
const classTotal = computed(() => students.value.reduce((sum, s) => sum + s.total_score, 0))
const spotQuote = computed(() => {
  const s = spotlightStudent.value
  if (!s?.pet_name) return ''
  const lines = getEvoLines(s.pet_name)
  return lines[Math.min(s.pet_level || 1, lines.length) - 1] || lines[0] || ''
})
const spotProgress = computed(() => {
  const s = spotlightStudent.value
  if (!s) return { percent: 0, remaining: 0 }
  return nextLevelProgress(s.total_score, s.pet_level || 1)
})
/** 功能入口导航：教室端用 classroom-*，教师基础模式用 teacher-*-basic */
function navTo(key: 'overview' | 'scores' | 'leaderboard' | 'pk' | 'pokedex') {
  const isClassroom = router.currentRoute.value.name?.toString().startsWith('classroom-')
  const map: Record<string, { classroom: string; basic: string }> = {
    overview: { classroom: 'classroom-overview', basic: 'teacher-dashboard-basic' },
    scores: { classroom: 'classroom-scores', basic: 'teacher-scores-basic' },
    leaderboard: { classroom: 'classroom-overview', basic: 'teacher-leaderboard-basic' },
    pk: { classroom: 'classroom-pk', basic: 'teacher-pk-basic' },
    pokedex: { classroom: 'classroom-pokedex', basic: 'teacher-pets-basic' },
  }
  const m = map[key]
  if (m) router.push({ name: isClassroom ? m.classroom : m.basic })
}

function openModal(s: StudentEntry, type: 'add' | 'sub') {
  modalStudent.value = s; modalType.value = type; actionError.value = ''; resetReasonPick(); showModal.value = true
}

function startEdit(sid: number) {
  editingStep.value = sid
  editInput.value = String(getStep(sid))
}

function saveEdit(sid: number, event?: Event) {
  const val = parseInt(event ? (event.target as HTMLInputElement).value : editInput.value)
  if (val >= 1 && val <= 100) stepValues.value[sid] = val
  editingStep.value = null
}

async function confirmAction() {
  const s = modalStudent.value; if (!s) return
  const reasons = selectedList()
  if (!reasons.length) return
  const step = getStep(s.id)
  const points = modalType.value === 'add' ? step * reasons.length : -step * reasons.length
  try {
    await apiPost('/api/v1/display/scores/give', { token: token.value, student_id: s.id, points, reason: reasons.join('、') })
    s.total_score = Math.max(0, s.total_score + points)
    showFloatText(s.id, points)
    actionError.value = ''
    resetReasonPick()
    showModal.value = false
  } catch (e: any) {
    const msg = e?.response?.data?.message || ''
    actionError.value = msg.includes('30') ? '单次超过 30 分需要教师账号登录操作' : (msg || '操作失败，请稍后重试')
    setTimeout(() => { actionError.value = '' }, 3000)
  }
}

function showFloatText(studentId: number, points: number) {
  const el = document.getElementById('card-' + studentId); if (!el) return
  const r = el.getBoundingClientRect(); const id = floatId++
  floatTexts.value.push({ id, x: r.left + r.width/2, y: r.top-10, text: points > 0 ? `+${points}` : `${points}`, color: points > 0 ? '#6ee7b7' : '#fca5a5' })
  setTimeout(() => { floatTexts.value = floatTexts.value.filter(f => f.id !== id) }, 1200)
}

onMounted(async () => {
  token.value = sessionStorage.getItem('class_token') || ''
  if (!token.value) return
  try {
    const res = await apiGet<{ data: StudentEntry[] }>('/api/v1/display/students', { params: { token: token.value } })
    students.value = (res.data || []).map(s => ({
      ...s,
      pet_emoji: s.pet_species ? getSpeciesEmoji(s.pet_species) : '🥚',
    }))
  } catch { /* ignore */ } finally { loading.value = false }
  try {
    const cRes = await apiGet<{ data: { pet_series?: string | null } }>('/api/v1/display/class-settings', { params: { token: token.value } })
    classPetSeries.value = cRes.data?.pet_series || ''
  } catch { classPetSeries.value = '' }
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:baseline;gap:12px;margin-bottom:20px;">
      <h2 style="font-size:24px;font-weight:700;margin:0;">✏️ 课堂评价</h2>
      <span style="font-size:13px;color:var(--md-text-secondary);">点击数字修改步长</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:60px;color:var(--md-text-secondary);">加载中...</div>

    <template v-else>
      <!-- 主卡片横幅（视觉优化方案 5.1：功能入口网格 + 宠物主视觉 + 统一积分操作区） -->
      <div class="pet-main-card">
        <div class="mc-header">
          <h2>✨ 学趣星球</h2>
          <button class="mc-nav" @click="navTo('overview')">📊 班级总览</button>
        </div>

        <div class="feature-grid">
          <button class="feature-item" @click="navTo('overview')"><span class="fi-icon">📊</span><span>班级总览</span></button>
          <button class="feature-item" @click="navTo('scores')"><span class="fi-icon">📝</span><span>课堂评价</span></button>
          <button class="feature-item" @click="navTo('leaderboard')"><span class="fi-icon">🏆</span><span>排行榜</span></button>
          <button class="feature-item" @click="navTo('pk')"><span class="fi-icon">⚔️</span><span>年级PK</span></button>
          <button class="feature-item" @click="navTo('pokedex')"><span class="fi-icon">📖</span><span>宠物图鉴</span></button>
        </div>

        <div v-if="spotlightStudent" class="pet-spotlight">
          <div class="spot-avatar" @click="openPetDetail(spotlightStudent)" title="点击查看宠物详情">
            <PetSprite :species-id="spotlightStudent.pet_species" :level="spotlightStudent.pet_level" :animate="true" />
          </div>
          <div class="spot-info">
            <div class="spot-name-row">
              <span class="spot-name">{{ spotlightStudent.pet_name || spotlightStudent.name }}</span>
              <span class="spot-level">Lv.{{ spotlightStudent.pet_level }}</span>
            </div>
            <div class="exp-bar"><div class="exp-fill" :style="{ width: spotProgress.percent + '%' }"></div></div>
            <div class="exp-text">
              <span v-if="spotProgress.remaining > 0">距 Lv.{{ (spotlightStudent.pet_level || 1) + 1 }} 还差 <strong style="color:var(--md-gold);">{{ spotProgress.remaining }}</strong> 分</span>
              <span v-else>✅ 已达到下一级条件</span>
            </div>
            <div class="spot-quote" v-if="spotQuote">“{{ spotQuote }}”</div>
          </div>
          <button class="spot-detail" @click="openPetDetail(spotlightStudent)">详情 →</button>
        </div>

        <div class="score-controls">
          <span style="font-size:12px;color:var(--md-text-secondary);">班级总分</span>
          <span class="total-score">{{ classTotal.toLocaleString() }}</span>
          <span class="mc-step-chip" style="margin-left:auto;" @click="spotlightStudent && startEdit(spotlightStudent.id)" title="点击修改步长">步长: <strong>{{ spotlightStudent ? getStep(spotlightStudent.id) : 1 }}</strong></span>
          <input v-if="spotlightStudent && editingStep === spotlightStudent.id" v-model="editInput" type="number" min="1" max="100"
            @blur="saveEdit(spotlightStudent.id)" @keydown.enter="saveEdit(spotlightStudent.id, $event)" autofocus
            style="width:56px;padding:6px 8px;text-align:center;background:rgba(167,139,250,0.15);border:1px solid rgba(167,139,250,0.3);border-radius:8px;color:var(--color-text);font-size:15px;font-weight:700;outline:none;font-family:inherit;">
        </div>
      </div>

      <!-- 首次使用欢迎提示 -->
      <div v-if="students.some(s => !s.pet_name)" style="margin-bottom:16px;padding:14px 20px;background:linear-gradient(135deg,rgba(167,139,250,0.08),rgba(244,114,182,0.05));border:1px solid rgba(167,139,250,0.15);border-radius:var(--md-radius);">
        <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
          <span style="font-size:24px;">🎉</span>
          <div>
            <div style="font-weight:700;font-size:15px;">欢迎来到学趣星球！</div>
            <div style="font-size:13px;color:var(--md-text-secondary);">点击学生卡片上的宠物 emoji，可以免费选择你的第一只宠物 🆓</div>
          </div>
        </div>
      </div>

      <div style="display:flex;gap:12px;margin-bottom:16px;flex-wrap:wrap;">
        <div style="display:flex;align-items:center;gap:8px;padding:8px 16px;background:var(--tint-2);border:1px solid var(--tint-3);border-radius:30px;flex:1;min-width:200px;max-width:360px;">
          <span>🔍</span>
          <input v-model="searchQuery" type="text" placeholder="搜索学生姓名..."
            style="background:transparent;border:none;outline:none;color:var(--color-text);font-size:14px;width:100%;font-family:inherit;">
        </div>
      </div>

      <!-- 批量操作栏：选中学生后出现 -->
      <div v-if="selectedIds.length > 0" style="display:flex;align-items:center;gap:10px;padding:10px 16px;background:var(--tint-1);border:1px solid var(--tint-3);border-radius:12px;margin-bottom:14px;flex-wrap:wrap;">
        <span style="font-size:13px;color:var(--md-text-secondary);">已选 <strong style="color:var(--md-primary);font-size:16px;">{{ selectedIds.length }}</strong> 名学生</span>
        <button @click="openBatchModal('add')"
          style="padding:6px 16px;border-radius:10px;border:none;background:rgba(16,185,129,0.12);color:#10b981;font-size:13px;font-weight:600;cursor:pointer;font-family:inherit;">✨ 批量加分</button>
        <button @click="openBatchModal('sub')"
          style="padding:6px 16px;border-radius:10px;border:none;background:rgba(239,68,68,0.12);color:#ef4444;font-size:13px;font-weight:600;cursor:pointer;font-family:inherit;">📉 批量减分</button>
        <button @click="clearSelect"
          style="padding:6px 12px;border-radius:10px;border:1px solid var(--tint-3);background:transparent;color:var(--md-text-secondary);font-size:13px;cursor:pointer;font-family:inherit;">取消选择</button>
      </div>

      <div v-if="filtered.length === 0" style="text-align:center;padding:60px;color:var(--md-text-secondary);">👀 没有找到学生</div>

      <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(212px,1fr));gap:14px;">
        <div v-for="s in filtered" :key="s.id" :id="'card-' + s.id"
          style="background:var(--tint-1);border-radius:var(--md-radius);padding:16px 14px 12px;border:1px solid var(--tint-2);position:relative;overflow:hidden;transition:0.25s;"
          :style="{ borderLeftColor: s.pet_level >= 10 ? '#f59e0b' : s.pet_level >= 7 ? '#8b5cf6' : s.pet_level >= 4 ? '#3b82f6' : '#6b7280', borderLeftWidth: '4px' }">
          <!-- 学号：独立右上角 -->
          <span v-if="s.student_no" style="position:absolute;top:12px;right:12px;font-size:10px;color:var(--md-text-secondary);background:var(--tint-2);padding:1px 8px;border-radius:6px;letter-spacing:0.5px;">📛{{ s.student_no }}</span>
          <!-- 姓名行 + 多选挖空圆 -->
          <div style="padding-right:56px;display:flex;align-items:center;gap:8px;">
            <div class="pick-circle" :class="{ picked: isSelected(s.id) }" @click.stop="toggleSelect(s.id)" title="多选后批量加减分">
              <svg v-if="isSelected(s.id)" viewBox="0 0 10 10" style="width:9px;height:9px;flex-shrink:0;"><path d="M1.2 5.2 L4 8 L8.8 2" stroke="#fff" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
            <span style="font-size:16px;font-weight:700;">{{ s.name }}</span>
          </div>
          <!-- 宠物：SVG 放大居中（不显示"xx的萌宠"文字） -->
          <div style="display:flex;flex-direction:column;align-items:center;gap:4px;margin:8px 0 12px;cursor:pointer;" @click="openPetDetail(s)" title="点击查看宠物详情">
            <span v-if="s.pet_species" style="width:120px;height:120px;transition:transform 0.2s;" @mouseenter="(e)=>(e.target as HTMLElement).style.transform='scale(1.08)'" @mouseleave="(e)=>(e.target as HTMLElement).style.transform=''">
              <PetSprite :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
            </span>
            <span v-else style="font-size:56px;line-height:1;">{{ s.pet_emoji }}</span>
            <span style="font-size:11px;font-weight:700;color:var(--md-gold);background:rgba(255,215,0,0.1);padding:1px 10px;border-radius:10px;">Lv.{{ s.pet_level }}</span>
            <span v-if="!s.pet_species" style="font-size:11px;color:var(--md-text-secondary);margin-top:2px;">未孵化 · 点击认养</span>
          </div>
          <div style="display:flex;justify-content:space-between;font-size:12px;color:var(--md-text-secondary);margin-bottom:4px;">
            <span>⭐ 积分</span>
            <span style="font-size:18px;font-weight:800;color:var(--md-text);">{{ s.total_score }}</span>
          </div>
          <!-- 距离下一级进度 -->
          <div v-if="s.pet_level > 0" style="margin-bottom:6px;">
            <div style="height:3px;background:var(--tint-3);border-radius:2px;overflow:hidden;margin-bottom:2px;">
              <div :style="{ width: nextLevelProgress(s.total_score, s.pet_level).percent + '%', height:'100%', background:'linear-gradient(90deg,var(--md-primary),var(--md-secondary))', borderRadius:'2px' }"></div>
            </div>
            <div style="font-size:10px;color:var(--md-text-secondary);text-align:right;">
              <span v-if="nextLevelProgress(s.total_score, s.pet_level).remaining > 0">距 Lv.{{ s.pet_level + 1 }} 还差 <strong style="color:var(--md-gold);">{{ nextLevelProgress(s.total_score, s.pet_level).remaining }}</strong> 分</span><span v-else>✅ 已达到 Lv.{{ s.pet_level + 1 }} 条件</span>
            </div>
          </div>
          <div style="display:flex;align-items:center;justify-content:center;gap:6px;padding-top:8px;border-top:1px solid var(--tint-2);">
            <button @click="openModal(s, 'sub')"
              style="width:36px;height:36px;border-radius:50%;border:1px solid rgba(239,68,68,0.2);background:rgba(239,68,68,0.04);color: var(--color-danger-text);font-size:20px;font-weight:700;cursor:pointer;transition:0.15s;">−</button>
            <!-- 可编辑步长 -->
            <input v-if="editingStep === s.id" v-model="editInput" type="number" min="1" max="100"
              @blur="saveEdit(s.id)" @keydown.enter="saveEdit(s.id, $event)" autofocus
              style="width:40px;text-align:center;background:rgba(167,139,250,0.15);border:1px solid rgba(167,139,250,0.3);border-radius:8px;color:var(--color-text);font-size:16px;font-weight:700;outline:none;font-family:inherit;">
            <span v-else @click="startEdit(s.id)"
              style="font-size:18px;font-weight:700;min-width:32px;text-align:center;cursor:pointer;padding:2px 4px;border-radius:8px;transition:background 0.15s;user-select:none;"
              @mouseenter="(e) => (e.target as HTMLElement).style.background='var(--tint-3)'"
              @mouseleave="(e) => (e.target as HTMLElement).style.background=''">{{ getStep(s.id) }}</span>
            <button @click="openModal(s, 'add')"
              style="width:36px;height:36px;border-radius:50%;border:1px solid rgba(16,185,129,0.2);background:rgba(16,185,129,0.04);color: var(--color-success-text);font-size:20px;font-weight:700;cursor:pointer;transition:0.15s;">+</button>
          </div>
        </div>
      </div>
    </template>

    <!-- Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showModal && modalStudent" @click.self="showModal = false"
          style="position:fixed;inset:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;z-index:300;">
          <div style="background:var(--color-bg-card);border:1px solid var(--tint-3);border-radius:var(--md-radius);padding:24px 28px;max-width:480px;width:90%;max-height:82vh;overflow-y:auto;box-shadow:var(--md-elevation);animation:popIn 0.25s ease;">
            <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:4px;">
              <h3 style="font-size:19px;font-weight:700;margin:0;">{{ modalType === 'add' ? '🌟 选择加分原因' : '⚠️ 选择减分原因' }}</h3>
              <span style="font-size:13px;font-weight:700;color:var(--md-primary);background:rgba(167,139,250,0.12);padding:3px 12px;border-radius:20px;">已选: <strong>{{ selectedList().length }}</strong> 项</span>
            </div>
            <p style="font-size:13px;color:var(--md-text-secondary);margin-bottom:16px;">为 <strong style="color:var(--color-text);">{{ modalStudent.name }}</strong> 选择原因 · 每项 <strong style="color:var(--md-gold);">{{ getStep(modalStudent.id) }}</strong> 分，可多选</p>
            <div style="display:flex;flex-direction:column;gap:12px;margin-bottom:12px;">
              <div v-for="g in reasonGroups" :key="g.title">
                <div style="display:flex;align-items:center;gap:8px;font-size:13px;font-weight:600;color:var(--md-text-secondary);margin-bottom:8px;">
                  <span>{{ g.emoji }} {{ g.title }}</span><span style="flex:1;height:1px;background:var(--tint-3);"></span>
                </div>
                <div style="display:flex;flex-wrap:wrap;gap:8px;">
                  <button v-for="r in g.items" :key="r" @click="toggleReason(r)"
                    class="reason-pick"
                    :style="selectedReasons.includes(r) ? { borderColor: g.color, background: g.color + '1A', color: g.color, fontWeight: 700 } : {}">
                    {{ r }}<span v-if="selectedReasons.includes(r)"> ✓</span>
                  </button>
                </div>
              </div>
            </div>
            <div v-if="showCustom" style="margin-bottom:10px;">
              <input v-model="customReason" type="text" placeholder="输入自定义原因..." maxlength="20"
                style="width:100%;padding:11px 16px;border-radius:20px;border:1px solid var(--tint-3);background:var(--tint-1);color:var(--color-text);font-size:14px;outline:none;font-family:inherit;">
            </div>
            <button @click="showCustom = !showCustom" style="padding:6px 14px;border-radius:20px;border:1px dashed var(--tint-4);background:transparent;color:var(--md-text-secondary);font-size:13px;cursor:pointer;font-family:inherit;margin-bottom:6px;">➕ 自定义原因</button>
            <div v-if="actionError" style="margin-bottom:12px;padding:10px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:10px;color: var(--color-danger-text);font-size:13px;text-align:center;">{{ actionError }}</div>
            <div style="display:flex;gap:10px;margin-top:12px;">
              <button @click="showModal = false; actionError = ''"
                style="flex:1;padding:12px;border-radius:12px;border:1px solid var(--tint-3);background:transparent;color:var(--md-text-secondary);font-size:14px;cursor:pointer;font-family:inherit;">取消</button>
              <button @click="confirmAction" :disabled="!selectedList().length"
                style="flex:1.4;padding:12px;border-radius:12px;border:none;background:var(--md-primary);color:#fff;font-size:14px;font-weight:700;cursor:pointer;font-family:inherit;"
                :style="!selectedList().length ? 'opacity:0.5;cursor:not-allowed;' : ''">✅ 确认{{ modalType === 'add' ? '加分' : '减分' }} ({{ selectedList().length }}项)</button>
            </div>
          </div>
        </div>
      </Transition>

      <div v-for="f in floatTexts" :key="f.id"
        style="position:fixed;pointer-events:none;font-size:24px;font-weight:800;z-index:999;animation:floatUp 1.2s ease-out forwards;"
        :style="{ left: f.x + 'px', top: f.y + 'px', color: f.color }">{{ f.text }}</div>
    </Teleport>

      <!-- 批量加减分弹窗 -->
      <Teleport to="body">
        <Transition name="fade">
          <div v-if="batchModal" @click.self="batchModal = false"
            style="position:fixed;inset:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;z-index:300;">
            <div style="background:var(--color-bg-card);border:1px solid var(--tint-3);border-radius:var(--md-radius);padding:24px 28px;max-width:480px;width:90%;max-height:82vh;overflow-y:auto;box-shadow:var(--md-elevation);animation:popIn 0.25s ease;">
              <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:4px;">
                <h3 style="font-size:19px;font-weight:700;margin:0;">{{ batchType === 'add' ? '✨ 批量加分' : '📉 批量减分' }}</h3>
                <span style="font-size:13px;font-weight:700;color:var(--md-primary);background:rgba(167,139,250,0.12);padding:3px 12px;border-radius:20px;">已选: <strong>{{ selectedList().length }}</strong> 项</span>
              </div>
              <p style="font-size:13px;color:var(--md-text-secondary);margin-bottom:16px;">为 <strong style="color:var(--color-text);">{{ selectedIds.length }}</strong> 名学生选择原因 · 每人每项 <strong style="color:var(--md-gold);">{{ getStep(selectedIds[0]) }}</strong> 分，可多选</p>
              <div style="display:flex;flex-direction:column;gap:12px;margin-bottom:12px;">
                <div v-for="g in batchReasonGroups" :key="g.title">
                  <div style="display:flex;align-items:center;gap:8px;font-size:13px;font-weight:600;color:var(--md-text-secondary);margin-bottom:8px;">
                    <span>{{ g.emoji }} {{ g.title }}</span><span style="flex:1;height:1px;background:var(--tint-3);"></span>
                  </div>
                  <div style="display:flex;flex-wrap:wrap;gap:8px;">
                    <button v-for="r in g.items" :key="r" @click="toggleReason(r)"
                      class="reason-pick"
                      :style="selectedReasons.includes(r) ? { borderColor: g.color, background: g.color + '1A', color: g.color, fontWeight: 700 } : {}">
                      {{ r }}<span v-if="selectedReasons.includes(r)"> ✓</span>
                    </button>
                  </div>
                </div>
              </div>
              <div v-if="showCustom" style="margin-bottom:10px;">
                <input v-model="customReason" type="text" placeholder="输入自定义原因..." maxlength="20"
                  style="width:100%;padding:11px 16px;border-radius:20px;border:1px solid var(--tint-3);background:var(--tint-1);color:var(--color-text);font-size:14px;outline:none;font-family:inherit;">
              </div>
              <button @click="showCustom = !showCustom" style="padding:6px 14px;border-radius:20px;border:1px dashed var(--tint-4);background:transparent;color:var(--md-text-secondary);font-size:13px;cursor:pointer;font-family:inherit;margin-bottom:6px;">➕ 自定义原因</button>
              <div v-if="batchError" style="margin-bottom:12px;padding:10px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:10px;color: var(--color-danger-text);font-size:13px;text-align:center;">{{ batchError }}</div>
              <div style="display:flex;gap:10px;margin-top:12px;">
                <button @click="batchModal = false"
                  style="flex:1;padding:12px;border-radius:12px;border:1px solid var(--tint-3);background:transparent;color:var(--md-text-secondary);font-size:14px;cursor:pointer;font-family:inherit;">取消</button>
                <button @click="confirmBatch" :disabled="batchBusy || !selectedList().length"
                  style="flex:1.4;padding:12px;border-radius:12px;border:none;background:var(--md-primary);color:#fff;font-size:14px;font-weight:700;cursor:pointer;font-family:inherit;"
                  :style="batchBusy || !selectedList().length ? 'opacity:0.5;cursor:not-allowed;' : ''">✅ 确认{{ batchType === 'add' ? '加分' : '减分' }} ({{ selectedList().length }}项)</button>
              </div>
            </div>
          </div>
        </Transition>
      </Teleport>

      <!-- 切换确认对话框 -->
      <Transition name="fade">
        <div v-if="confirmSwitch" @click.self="confirmSwitch = null"
          style="position:fixed;inset:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;z-index:301;">
          <div style="background:var(--color-bg-card);border:1px solid var(--tint-3);border-radius:var(--md-radius);padding:28px 32px;max-width:380px;width:90%;box-shadow:var(--md-elevation);animation:popIn 0.25s ease;text-align:center;">
            <div style="font-size:36px;margin-bottom:12px;">🔄</div>
            <h3 style="font-size:18px;font-weight:700;margin-bottom:8px;">确认切换宠物？</h3>
            <p style="font-size:14px;color:var(--md-text-secondary);margin-bottom:12px;">
              将切换为 <strong style="color:var(--color-primary);">{{ confirmSwitch.name }}</strong>
            </p>
            <div v-if="petPickerStudent?.pet_name && petPickerStudent.free_pick"
              style="font-size:14px;padding:10px 14px;border-radius:10px;background:rgba(16,185,129,0.12);border:1px solid rgba(16,185,129,0.3);color: var(--color-success-text);font-weight:700;margin-bottom:16px;">
              🎉 免费窗口期内，本次切换免费！
            </div>
            <div v-else-if="petPickerStudent?.pet_name && getSwitchCost(petPickerStudent) > 0"
              style="font-size:14px;padding:10px 14px;border-radius:10px;background:rgba(245,158,11,0.12);border:1px solid rgba(245,158,11,0.3);color: var(--color-warning-text);font-weight:700;margin-bottom:16px;">
              💰 本次切换扣除 <strong style="font-size:18px;">{{ getSwitchCost(petPickerStudent) }}</strong> 积分 · 保留当前等级
            </div>
            <div v-else style="font-size:14px;padding:10px 14px;border-radius:10px;background:rgba(16,185,129,0.12);border:1px solid rgba(16,185,129,0.3);color: var(--color-success-text);font-weight:700;margin-bottom:16px;">
              🎉 首次免费，不扣积分
            </div>
            <div style="display:flex;gap:10px;">
              <button @click="confirmSwitch = null" style="flex:1;padding:10px;border-radius:10px;border:1px solid var(--tint-3);background:transparent;color:var(--md-text-secondary);font-size:14px;cursor:pointer;font-family:inherit;">取消</button>
              <button @click="executeSwitch" :disabled="switchStatus !== 'idle'" :style="{ flex:'1', padding:'10px', borderRadius:'10px', border:'none', background: switchStatus === 'loading' ? '#f59e0b' : switchStatus === 'success' ? '#10b981' : switchStatus === 'error' ? '#ef4444' : 'rgba(167,139,250,0.15)', color: switchStatus !== 'idle' ? '#fff' : 'var(--color-primary)', fontSize:'14px', fontWeight:'600', cursor:'pointer', fontFamily:'inherit' }">
                <template v-if="switchStatus === 'loading'">切换中...</template>
                <template v-else-if="switchStatus === 'success'">切换成功 ✓</template>
                <template v-else-if="switchStatus === 'error'">切换失败 ✗</template>
                <template v-else>确认切换</template>
              </button>
            </div>
            <div v-if="switchError" style="margin-top:12px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color: var(--color-danger-text);font-size:12px;">{{ switchError }}</div>
          </div>
        </div>
      </Transition>

      <!-- 宠物角色介绍弹窗 -->
      <PetDetailModal
        v-if="showPetDetail && detailStudent"
        :species-id="detailStudent.pet_species"
        :level="detailStudent.pet_level"
        :score="detailStudent.total_score"
        :show-pet-switch="true"
        @close="showPetDetail = false"
        @switch-pet="showPetDetail = false; openPetPicker(detailStudent)"
      />

      <!-- 宠物选择器（班级配置当前系列时只展示该系列物种） -->
      <Transition name="fade">
        <div v-if="showPetPicker && petPickerStudent" @click.self="showPetPicker = false"
          style="position:fixed;inset:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;z-index:300;">
          <div style="background:var(--color-bg-card);border:1px solid var(--tint-3);border-radius:var(--md-radius);padding:24px 28px;max-width:520px;width:90%;max-height:80vh;overflow-y:auto;box-shadow:var(--md-elevation);animation:popIn 0.25s ease;">
            <div style="display:flex;align-items:center;gap:10px;margin-bottom:16px;">
              <span style="font-size:28px;">{{ petPickerStudent.pet_emoji }}</span>
              <div>
                <div style="font-size:16px;font-weight:700;">{{ petPickerStudent.name }} · 选择宠物</div>
                <div v-if="petPickerStudent.pet_name && petPickerStudent.free_pick" style="font-size:12px;margin-top:4px;padding:4px 10px;border-radius:8px;background:rgba(16,185,129,0.12);border:1px solid rgba(16,185,129,0.3);color: var(--color-success-text);font-weight:700;display:inline-block;">🎉 免费窗口期内，本次切换免费！</div>
                <div v-else-if="petPickerStudent.pet_name" style="font-size:12px;margin-top:4px;padding:4px 10px;border-radius:8px;background:rgba(245,158,11,0.12);border:1px solid rgba(245,158,11,0.3);color: var(--color-warning-text);font-weight:700;display:inline-block;">
                  💰 本次切换扣除 <strong style="font-size:14px;">{{ getSwitchCost(petPickerStudent) }}</strong> 积分 · 保留等级
                </div>
                <div v-else style="font-size:12px;margin-top:4px;padding:4px 10px;border-radius:8px;background:rgba(16,185,129,0.12);border:1px solid rgba(16,185,129,0.3);color: var(--color-success-text);font-weight:700;display:inline-block;">🎉 首次免费选择，不扣积分</div>
              </div>
              <button @click="showPetPicker = false" style="margin-left:auto;width:28px;height:28px;border-radius:50%;border:1px solid var(--tint-3);background:transparent;color:var(--color-text-secondary);cursor:pointer;">✕</button>
            </div>
            <div v-for="series in pickerSeriesList" :key="series.id" style="margin-bottom:12px;">
              <div style="font-size:12px;font-weight:600;color:var(--md-text-secondary);margin-bottom:6px;padding-left:4px;">{{ series.emoji }} {{ series.name }}</div>
              <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(90px,1fr));gap:6px;">
                <button v-for="sp in series.species" :key="sp.id" @click="handlePick(sp.id, sp.name)"
                  :disabled="switchingPet || petPickerStudent.pet_species === sp.id"
                  style="padding:8px 4px;border-radius:10px;border:1px solid var(--tint-2);background:var(--tint-1);text-align:center;cursor:pointer;transition:0.15s;font-family:inherit;"
                  :style="petPickerStudent.pet_species === sp.id ? 'border-color:rgba(16,185,129,0.35);background:rgba(16,185,129,0.08);cursor:default;' : ''"
                  @mouseenter="(e)=>petPickerStudent?.pet_species === sp.id || ((e.target as HTMLElement).style.background='var(--tint-3)')"
                  @mouseleave="(e)=>(e.target as HTMLElement).style.background=petPickerStudent?.pet_species === sp.id ? 'rgba(16,185,129,0.08)' : 'var(--tint-1)'">
                  <div style="width:48px;height:48px;margin:0 auto 2px;">
                    <PetSprite :species-id="sp.id" :level="6" />
                  </div>
                  <div style="font-size:10px;font-weight:500;color:var(--md-text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ sp.name }}</div>
                  <div v-if="petPickerStudent.pet_species === sp.id" style="font-size:9px;font-weight:700;color:#10B981;">✓ 当前</div>
                </button>
              </div>
            </div>
            <button @click="showPetPicker = false" style="width:100%;margin-top:8px;padding:8px;border-radius:10px;border:1px solid var(--tint-3);background:transparent;color:var(--md-text-secondary);font-size:13px;cursor:pointer;font-family:inherit;">取消</button>
          </div>
        </div>
      </Transition>

    <style>
      @keyframes popIn { from { transform: scale(0.92); opacity: 0; } to { transform: scale(1); opacity: 1; } }
      @keyframes floatUp { 0% { opacity: 1; transform: translateX(-50%) translateY(0) scale(1); } 100% { opacity: 0; transform: translateX(-50%) translateY(-80px) scale(1.3); } }
      .fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
      .fade-enter-from, .fade-leave-to { opacity: 0; }
      /* 多选挖空圆 checkbox */
      .pick-circle {
        width: 18px; height: 18px; border-radius: 50%;
        border: 2px solid var(--tint-4);
        display: flex; align-items: center; justify-content: center;
        cursor: pointer; transition: 0.15s; flex-shrink: 0; background: transparent;
      }
      .pick-circle:hover { border-color: var(--md-primary); box-shadow: 0 0 0 3px rgba(99,102,241,0.12); }
      .pick-circle.picked { background: var(--md-primary); border-color: var(--md-primary); }
      /* 加减分原因标签：44px 触控 pill（视觉优化方案 5.3） */
      .reason-pick {
        min-height: 44px;
        padding: 8px 16px;
        border-radius: 20px;
        border: 1px solid var(--tint-3);
        background: var(--tint-1);
        color: var(--color-text);
        font-size: 13px;
        cursor: pointer;
        transition: all 0.15s ease;
        font-family: inherit;
      }
      .reason-pick:hover { background: var(--tint-3); transform: translateY(-1px); }
      /* 主卡片横幅（视觉优化方案 5.1：功能入口网格 + 宠物主视觉 + 统一积分操作区） */
      .pet-main-card {
        background: var(--color-bg-card);
        border: 1px solid var(--tint-3);
        border-radius: 20px;
        padding: 20px;
        margin-bottom: 16px;
        box-shadow: var(--md-elevation);
      }
      .mc-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; }
      .mc-header h2 { font-size: 22px; font-weight: 800; margin: 0; }
      .mc-nav {
        padding: 7px 14px; border-radius: 20px;
        border: 1px solid rgba(167,139,250,0.25); background: rgba(167,139,250,0.1);
        color: var(--color-primary); font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; transition: 0.15s;
      }
      .mc-nav:hover { background: rgba(167,139,250,0.18); }
      .feature-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; margin-bottom: 16px; }
      .feature-item {
        display: flex; flex-direction: column; align-items: center; gap: 5px;
        padding: 12px 4px; border-radius: 14px;
        background: var(--tint-1); border: 1px solid var(--tint-2);
        color: var(--color-text); font-size: 11px; font-weight: 600;
        cursor: pointer; font-family: inherit; transition: 0.15s;
      }
      .feature-item:hover { background: var(--tint-3); transform: translateY(-1px); }
      .fi-icon { font-size: 20px; }
      .pet-spotlight {
        display: flex; align-items: center; gap: 16px;
        padding: 16px; border-radius: 16px;
        background: linear-gradient(135deg, rgba(167,139,250,0.1), rgba(244,114,182,0.05));
        border: 1px solid rgba(167,139,250,0.15); margin-bottom: 12px;
      }
      .spot-avatar { width: 140px; height: 140px; flex-shrink: 0; cursor: pointer; }
      .spot-info { flex: 1; min-width: 0; }
      .spot-name-row { display: flex; align-items: baseline; gap: 10px; margin-bottom: 6px; }
      .spot-name { font-size: 18px; font-weight: 800; }
      .spot-level {
        font-size: 12px; font-weight: 700; color: var(--md-gold);
        background: rgba(255,215,0,0.12); padding: 2px 10px; border-radius: 12px;
      }
      .exp-bar { height: 6px; background: var(--tint-3); border-radius: 3px; overflow: hidden; margin-bottom: 4px; }
      .exp-fill { height: 100%; background: linear-gradient(90deg, var(--md-primary), var(--md-secondary)); border-radius: 3px; }
      .exp-text { font-size: 11px; color: var(--md-text-secondary); margin-bottom: 6px; }
      .spot-quote { font-size: 13px; font-style: italic; color: var(--md-text-secondary); }
      .spot-detail {
        padding: 8px 16px; border-radius: 12px; border: none;
        background: var(--md-primary); color: #fff;
        font-size: 13px; font-weight: 700; cursor: pointer; font-family: inherit; white-space: nowrap;
      }
      .score-controls { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
      .total-score { font-size: 20px; font-weight: 800; color: var(--color-text); }
      .mc-step-chip {
        padding: 6px 14px; border-radius: 20px;
        background: rgba(167,139,250,0.1); border: 1px solid rgba(167,139,250,0.2);
        color: var(--color-text); font-size: 13px; cursor: pointer; transition: 0.15s;
      }
      .mc-step-chip:hover { background: rgba(167,139,250,0.18); }
      @media (max-width: 640px) {
        .feature-grid { grid-template-columns: repeat(3, 1fr); }
        .pet-spotlight { flex-direction: column; align-items: center; text-align: center; }
      }
    </style>
  </div>
</template>
