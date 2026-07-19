<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { getSpeciesEmoji } from '@/utils/petData'

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
const reasonsAdd = ['📖 举手发言', '✅ 作业优秀', '🤝 帮助同学', '🧹 遵守纪律', '🏆 挑战难题']
const reasonsSub = ['⚠️ 上课走神', '📕 作业缺交', '🗣️ 打扰课堂', '🏃 追逐打闹']
const floatTexts = ref<Array<{ id: number; x: number; y: number; text: string; color: string }>>([])
let floatId = 0

const filtered = computed(() => {
  if (!searchQuery.value) return students.value
  return students.value.filter(s => s.name.includes(searchQuery.value))
})

function openModal(s: StudentEntry, type: 'add' | 'sub') {
  modalStudent.value = s; modalType.value = type; showModal.value = true
}

async function executeAction(reason: string) {
  const s = modalStudent.value; if (!s) return
  const points = modalType.value === 'add' ? 1 : -1
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
    students.value = res.data || []
  } catch { /* ignore */ } finally { loading.value = false }
})
</script>

<template>
  <div>
    <div style="display:flex;align-items:baseline;gap:12px;margin-bottom:20px;">
      <h2 style="font-size:24px;font-weight:700;margin:0;">✏️ 课堂评价</h2>
      <span style="font-size:13px;color:var(--md-text-secondary);">点击 +/− 选择行为原因</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:60px;color:var(--md-text-secondary);">加载中...</div>

    <template v-else>
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
          <div style="display:flex;align-items:center;gap:6px;margin-bottom:8px;">
            <span style="font-size:20px;">{{ s.pet_species ? getSpeciesEmoji(s.pet_species) : '🥚' }}</span>
            <span style="font-size:13px;color:var(--md-text-secondary);">{{ s.pet_name || '未孵化' }}</span>
          </div>
          <div style="display:flex;justify-content:space-between;font-size:12px;color:var(--md-text-secondary);margin-bottom:6px;">
            <span>⭐ 积分</span>
            <span style="font-size:18px;font-weight:800;color:var(--md-text);">{{ s.total_score }}</span>
          </div>
          <div style="display:flex;align-items:center;justify-content:center;gap:8px;padding-top:8px;border-top:1px solid rgba(255,255,255,0.04);">
            <button @click="openModal(s, 'sub')"
              style="width:36px;height:36px;border-radius:50%;border:1px solid rgba(239,68,68,0.2);background:rgba(239,68,68,0.04);color:#fca5a5;font-size:20px;font-weight:700;cursor:pointer;transition:0.15s;">−</button>
            <span style="font-size:18px;font-weight:700;min-width:28px;text-align:center;">1</span>
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
            <p style="font-size:14px;color:var(--md-text-secondary);margin-bottom:20px;">为 <strong style="color:#fff;">{{ modalStudent.name }}</strong> 选择原因（每次1分）</p>
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

      <!-- Floating score text -->
      <div v-for="f in floatTexts" :key="f.id"
        style="position:fixed;pointer-events:none;font-size:24px;font-weight:800;z-index:999;animation:floatUp 1.2s ease-out forwards;"
        :style="{ left: f.x + 'px', top: f.y + 'px', color: f.color }">{{ f.text }}</div>
    </Teleport>

    <style>
      @keyframes popIn { from { transform: scale(0.92); opacity: 0; } to { transform: scale(1); opacity: 1; } }
      @keyframes floatUp { 0% { opacity: 1; transform: translateX(-50%) translateY(0) scale(1); } 100% { opacity: 0; transform: translateX(-50%) translateY(-80px) scale(1.3); } }
      .fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
      .fade-enter-from, .fade-leave-to { opacity: 0; }
    </style>
  </div>
</template>
