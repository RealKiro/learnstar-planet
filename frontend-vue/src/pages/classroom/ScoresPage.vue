<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { getSpeciesEmoji, getSpeciesBySeries, PET_SERIES, getSeriesName } from '@/utils/petData'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()

interface StudentEntry {
  id: number; name: string; student_no: string; total_score: number
  pet_name: string; pet_species: string; pet_level: number; pet_emoji: string
}

const students = ref<StudentEntry[]>([])
const loading = ref(true)
const token = ref('')
const searchQuery = ref('')
const showModal = ref(false)
const modalStudent = ref<StudentEntry | null>(null)
const modalType = ref<'add' | 'sub'>('add')
const stepValues = ref<Record<number, number>>({})
const editingStep = ref<number | null>(null)
const editInput = ref('1')
const reasonsAdd = ['📖 举手发言', '✅ 作业优秀', '🤝 帮助同学', '🧹 遵守纪律', '🏆 挑战难题']
const reasonsSub = ['⚠️ 上课走神', '📕 作业缺交', '🗣️ 打扰课堂', '🏃 追逐打闹']
const floatTexts = ref<Array<{ id: number; x: number; y: number; text: string; color: string }>>([])
let floatId = 0

// 宠物切换
const showPetPicker = ref(false)
const petPickerStudent = ref<StudentEntry | null>(null)
const switchingPet = ref(false)
const allSeriesList = PET_SERIES

function openPetPicker(s: StudentEntry) {
  petPickerStudent.value = s
  showPetPicker.value = true
}

async function switchPet(speciesId: string) {
  const s = petPickerStudent.value
  if (!s || switchingPet.value) return
  switchingPet.value = true
  try {
    const res = await apiPost<{ data: { pet_emoji: string; total_score: number; cost: number } }>(
      '/api/v1/display/pets/switch', { token: token.value, student_id: s.id, pet_species: speciesId }
    )
    s.pet_species = speciesId
    s.pet_emoji = res.data.pet_emoji
    s.total_score = res.data.total_score
    toast.show(res.data.cost > 0 ? `已切换，扣除${res.data.cost}积分` : '🎉 首次切换免费！', 'success')
    showPetPicker.value = false
  } catch (e: any) {
    toast.show(e?.response?.data?.message || '切换失败', 'error')
  } finally {
    switchingPet.value = false
  }
}

function getStep(sid: number) { return stepValues.value[sid] || 1 }

const filtered = computed(() => {
  if (!searchQuery.value) return students.value
  return students.value.filter(s => s.name.includes(searchQuery.value))
})

function openModal(s: StudentEntry, type: 'add' | 'sub') {
  modalStudent.value = s; modalType.value = type; showModal.value = true
}

function startEdit(sid: number) {
  editingStep.value = sid
  editInput.value = String(getStep(sid))
}

function saveEdit(sid: number) {
  const val = parseInt(editInput.value)
  if (val >= 1 && val <= 100) stepValues.value[sid] = val
  editingStep.value = null
}

async function executeAction(reason: string) {
  const s = modalStudent.value; if (!s) return
  const step = getStep(s.id)
  const points = modalType.value === 'add' ? step : -step
  try {
    await apiPost('/api/v1/display/scores/give', { token: token.value, student_id: s.id, points, reason })
    s.total_score = Math.max(0, s.total_score + points)
    showFloatText(s.id, points)
  } catch { /* ignore */ }
  showModal.value = false
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
        <div style="display:flex;align-items:center;gap:8px;padding:8px 16px;background:rgba(255,255,255,0.05);border:1px solid rgba(255,255,255,0.06);border-radius:30px;flex:1;min-width:200px;max-width:360px;">
          <span>🔍</span>
          <input v-model="searchQuery" type="text" placeholder="搜索学生姓名..."
            style="background:transparent;border:none;outline:none;color:#fff;font-size:14px;width:100%;font-family:inherit;">
        </div>
      </div>

      <div v-if="filtered.length === 0" style="text-align:center;padding:60px;color:var(--md-text-secondary);">👀 没有找到学生</div>

      <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:14px;">
        <div v-for="s in filtered" :key="s.id" :id="'card-' + s.id"
          style="background:rgba(255,255,255,0.03);border-radius:var(--md-radius);padding:16px 14px 12px;border:1px solid rgba(255,255,255,0.04);position:relative;overflow:hidden;transition:0.25s;"
          :style="{ borderLeftColor: s.pet_level >= 10 ? '#f59e0b' : s.pet_level >= 7 ? '#8b5cf6' : s.pet_level >= 4 ? '#3b82f6' : '#6b7280', borderLeftWidth: '4px' }">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:6px;">
            <span style="font-size:16px;font-weight:700;">{{ s.name }}</span>
            <span style="font-size:12px;font-weight:700;padding:2px 10px;border-radius:12px;background:rgba(255,215,0,0.12);color:var(--md-gold);">Lv.{{ s.pet_level }}</span>
          </div>
          <div style="display:flex;align-items:center;gap:6px;margin-bottom:8px;cursor:pointer;" @click="openPetPicker(s)" title="点击切换宠物">
            <span style="font-size:28px;transition:transform 0.2s;" @mouseenter="(e)=>(e.target as HTMLElement).style.transform='scale(1.15)'" @mouseleave="(e)=>(e.target as HTMLElement).style.transform=''">{{ s.pet_emoji }}</span>
            <span style="font-size:13px;color:var(--md-text-secondary);border-bottom:1px dashed rgba(255,255,255,0.15);">{{ s.pet_name || '未孵化' }}</span>
            <span style="font-size:10px;color:rgba(167,139,250,0.4);margin-left:auto;">换宠</span>
          </div>
          <div style="display:flex;justify-content:space-between;font-size:12px;color:var(--md-text-secondary);margin-bottom:6px;">
            <span>⭐ 积分</span>
            <span style="font-size:18px;font-weight:800;color:var(--md-text);">{{ s.total_score }}</span>
          </div>
          <div style="display:flex;align-items:center;justify-content:center;gap:6px;padding-top:8px;border-top:1px solid rgba(255,255,255,0.04);">
            <button @click="openModal(s, 'sub')"
              style="width:36px;height:36px;border-radius:50%;border:1px solid rgba(239,68,68,0.2);background:rgba(239,68,68,0.04);color:#fca5a5;font-size:20px;font-weight:700;cursor:pointer;transition:0.15s;">−</button>
            <!-- 可编辑步长 -->
            <input v-if="editingStep === s.id" v-model="editInput" type="number" min="1" max="100"
              @blur="saveEdit(s.id)" @keydown.enter="saveEdit(s.id)" autofocus
              style="width:40px;text-align:center;background:rgba(167,139,250,0.15);border:1px solid rgba(167,139,250,0.3);border-radius:8px;color:#fff;font-size:16px;font-weight:700;outline:none;font-family:inherit;">
            <span v-else @click="startEdit(s.id)"
              style="font-size:18px;font-weight:700;min-width:32px;text-align:center;cursor:pointer;padding:2px 4px;border-radius:8px;transition:background 0.15s;user-select:none;"
              @mouseenter="(e) => (e.target as HTMLElement).style.background='rgba(255,255,255,0.06)'"
              @mouseleave="(e) => (e.target as HTMLElement).style.background=''">{{ getStep(s.id) }}</span>
            <button @click="openModal(s, 'add')"
              style="width:36px;height:36px;border-radius:50%;border:1px solid rgba(16,185,129,0.2);background:rgba(16,185,129,0.04);color:#6ee7b7;font-size:20px;font-weight:700;cursor:pointer;transition:0.15s;">+</button>
          </div>
        </div>
      </div>
    </template>

    <!-- Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showModal && modalStudent" @click.self="showModal = false"
          style="position:fixed;inset:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;z-index:300;">
          <div style="background:#1e1b3b;border:1px solid rgba(255,255,255,0.08);border-radius:var(--md-radius);padding:28px 32px;max-width:420px;width:90%;box-shadow:var(--md-elevation);animation:popIn 0.25s ease;">
            <h3 style="font-size:20px;font-weight:700;margin-bottom:6px;">{{ modalType === 'add' ? '🌟 选择加分行为' : '⚠️ 选择减分原因' }}</h3>
            <p style="font-size:14px;color:var(--md-text-secondary);margin-bottom:20px;">为 <strong style="color:#fff;">{{ modalStudent.name }}</strong> 选择原因（每次<strong style="color:var(--md-gold);">{{ getStep(modalStudent.id) }}</strong>分）</p>
            <div style="display:flex;flex-direction:column;gap:8px;margin-bottom:20px;">
              <button v-for="r in (modalType === 'add' ? reasonsAdd : reasonsSub)" :key="r" @click="executeAction(r)"
                style="padding:12px 16px;border-radius:12px;border:1px solid rgba(255,255,255,0.06);background:rgba(255,255,255,0.03);color:#fff;font-size:15px;text-align:left;cursor:pointer;transition:0.15s;font-family:inherit;">
                {{ r }}
              </button>
            </div>
            <button @click="showModal = false"
              style="width:100%;padding:10px;border-radius:12px;border:1px solid rgba(255,255,255,0.08);background:transparent;color:var(--md-text-secondary);font-size:14px;cursor:pointer;font-family:inherit;">取消操作</button>
          </div>
        </div>
      </Transition>

      <div v-for="f in floatTexts" :key="f.id"
        style="position:fixed;pointer-events:none;font-size:24px;font-weight:800;z-index:999;animation:floatUp 1.2s ease-out forwards;"
        :style="{ left: f.x + 'px', top: f.y + 'px', color: f.color }">{{ f.text }}</div>
    </Teleport>

      <!-- 宠物选择器（展示所有物种） -->
      <Transition name="fade">
        <div v-if="showPetPicker && petPickerStudent" @click.self="showPetPicker = false"
          style="position:fixed;inset:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;z-index:300;">
          <div style="background:#1e1b3b;border:1px solid rgba(255,255,255,0.08);border-radius:var(--md-radius);padding:24px 28px;max-width:520px;width:90%;max-height:80vh;overflow-y:auto;box-shadow:var(--md-elevation);animation:popIn 0.25s ease;">
            <div style="display:flex;align-items:center;gap:10px;margin-bottom:16px;">
              <span style="font-size:28px;">{{ petPickerStudent.pet_emoji }}</span>
              <div>
                <div style="font-size:16px;font-weight:700;">{{ petPickerStudent.name }} · 选择宠物</div>
                <div style="font-size:12px;color:var(--md-text-secondary);">{{ petPickerStudent.pet_name ? '后续切换扣20积分 · 保留等级' : '🎉 首次免费选择！' }}</div>
              </div>
              <button @click="showPetPicker = false" style="margin-left:auto;width:28px;height:28px;border-radius:50%;border:1px solid rgba(255,255,255,0.06);background:transparent;color:rgba(255,255,255,0.4);cursor:pointer;">✕</button>
            </div>
            <div v-for="series in allSeriesList" :key="series.id" style="margin-bottom:12px;">
              <div style="font-size:12px;font-weight:600;color:var(--md-text-secondary);margin-bottom:6px;padding-left:4px;">{{ series.emoji }} {{ series.name }}</div>
              <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(90px,1fr));gap:6px;">
                <button v-for="sp in series.species" :key="sp.id" @click="switchPet(sp.id)"
                  :disabled="switchingPet"
                  style="padding:8px 4px;border-radius:10px;border:1px solid rgba(255,255,255,0.04);background:rgba(255,255,255,0.02);text-align:center;cursor:pointer;transition:0.15s;font-family:inherit;"
                  :style="petPickerStudent.pet_species === sp.id ? 'border-color:rgba(167,139,250,0.3);background:rgba(167,139,250,0.08);' : ''"
                  @mouseenter="(e)=>(e.target as HTMLElement).style.background='rgba(255,255,255,0.06)'"
                  @mouseleave="(e)=>(e.target as HTMLElement).style.background=petPickerStudent.pet_species === sp.id ? 'rgba(167,139,250,0.08)' : 'rgba(255,255,255,0.02)'">
                  <div style="font-size:22px;margin-bottom:2px;">{{ getSpeciesEmoji(sp.id) }}</div>
                  <div style="font-size:10px;font-weight:500;color:var(--md-text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ sp.name }}</div>
                </button>
              </div>
            </div>
            <button @click="showPetPicker = false" style="width:100%;margin-top:8px;padding:8px;border-radius:10px;border:1px solid rgba(255,255,255,0.06);background:transparent;color:var(--md-text-secondary);font-size:13px;cursor:pointer;font-family:inherit;">取消</button>
          </div>
        </div>
      </Transition>

    <style>
      @keyframes popIn { from { transform: scale(0.92); opacity: 0; } to { transform: scale(1); opacity: 1; } }
      @keyframes floatUp { 0% { opacity: 1; transform: translateX(-50%) translateY(0) scale(1); } 100% { opacity: 0; transform: translateX(-50%) translateY(-80px) scale(1.3); } }
      .fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
      .fade-enter-from, .fade-leave-to { opacity: 0; }
    </style>
  </div>
</template>
