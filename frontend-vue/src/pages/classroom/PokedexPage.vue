<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiPost } from '@/utils/api'
import { getAllSeries, getSpeciesEmoji, PET_SERIES, getSpeciesById, getPetLevelName } from '@/utils/petData'
import { getPoems, getEvoLines, STAGE_NAMES } from '@/utils/petHandbookData'

const token = ref('')
const currentSeries = ref('myth')
const currentSlide = ref(0)
const switching = ref(false)
const switchMsg = ref('')
const switchError = ref('')
const selectedSpecies = ref<{ seriesId: string; speciesId: string; name: string } | null>(null)
const detailStage = ref(0)
let timer: ReturnType<typeof setInterval>

const allSeries = getAllSeries()

const seriesList = computed(() => PET_SERIES)
const currentSeriesData = computed(() => seriesList.value.find(s => s.id === currentSeries.value) || seriesList.value[0])

const detailSpecies = computed(() => {
  if (!selectedSpecies.value) return null
  return getSpeciesById(selectedSpecies.value.speciesId)
})

const detailPoems = computed(() => {
  if (!selectedSpecies.value) return []
  return getPoems(selectedSpecies.value.name)
})

const detailEvoLines = computed(() => {
  if (!selectedSpecies.value) return []
  return getEvoLines(selectedSpecies.value.name)
})

function goToSlide(idx: number) {
  currentSlide.value = idx
  clearInterval(timer)
  timer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % seriesList.value.length }, 5000)
}

function openDetail(speciesId: string, name: string) {
  selectedSpecies.value = { seriesId: currentSeries.value, speciesId, name }
  detailStage.value = 0
}

function closeDetail() {
  selectedSpecies.value = null
}

async function switchSeries() {
  switchMsg.value = ''
  switchError.value = ''
  switching.value = true
  try {
    const res = await apiPost<{ data: { series_id: string; cost_per_student: number; affected_students: number } }>(
      '/api/v1/display/switch-series', { token: token.value, series_id: currentSeries.value }
    )
    switchMsg.value = `✅ 已切换，全班 ${res.data.affected_students} 人各扣除 ${res.data.cost_per_student} 积分`
  } catch (e: any) {
    switchError.value = e?.response?.data?.message || '切换失败'
  } finally {
    switching.value = false
  }
}

onMounted(() => {
  token.value = sessionStorage.getItem('class_token') || ''
  timer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % seriesList.value.length }, 5000)
})
onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div>
    <!-- 顶栏 -->
    <div style="display:flex;align-items:baseline;justify-content:space-between;flex-wrap:wrap;gap:12px;margin-bottom:20px;">
      <h2 style="font-size:24px;font-weight:700;margin:0;">📚 宠物图鉴</h2>
      <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
        <label style="font-size:13px;color:var(--md-text-secondary);">🏷️ 系列</label>
        <select v-model="currentSeries"
          style="padding:8px 14px;border-radius:var(--md-radius);background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.08);color:#fff;font-size:14px;font-weight:500;outline:none;cursor:pointer;font-family:inherit;">
          <option v-for="s in allSeries" :key="s.id" :value="s.id" style="background:#1a1a2e;color:#f1f1f1;">{{ s.emoji }} {{ s.name }}</option>
        </select>
        <button @click="switchSeries" :disabled="switching"
          style="padding:8px 18px;border-radius:30px;border:none;background:rgba(167,139,250,0.15);color:var(--md-primary-light);font-size:13px;font-weight:600;cursor:pointer;transition:0.15s;font-family:inherit;">
          {{ switching ? '切换中...' : '切换（每人扣20分）' }}
        </button>
      </div>
    </div>

    <!-- 切换消息 -->
    <div v-if="switchMsg" style="margin-bottom:16px;padding:10px 16px;background:rgba(16,185,129,0.08);border:1px solid rgba(16,185,129,0.15);border-radius:var(--md-radius);color:#6ee7b7;font-size:13px;">{{ switchMsg }}</div>
    <div v-if="switchError" style="margin-bottom:16px;padding:10px 16px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.15);border-radius:var(--md-radius);color:#fca5a5;font-size:13px;">{{ switchError }}</div>

    <!-- 物种大图（点击物种后展示） -->
    <Transition name="fade">
      <div v-if="selectedSpecies && detailSpecies" @click.self="closeDetail"
        style="position:fixed;inset:0;z-index:300;background:rgba(5,2,20,0.85);backdrop-filter:blur(16px);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:linear-gradient(180deg,#1a1040,#0d1b2a);border:1px solid rgba(255,255,255,0.08);border-radius:24px;max-width:700px;width:100%;max-height:85vh;overflow-y:auto;padding:28px;box-shadow:0 20px 60px rgba(0,0,0,0.5);">
          <!-- 头部 -->
          <div style="display:flex;align-items:center;gap:12px;margin-bottom:20px;">
            <span style="font-size:48px;">{{ getSpeciesEmoji(selectedSpecies.speciesId) }}</span>
            <div>
              <div style="font-size:22px;font-weight:700;">{{ selectedSpecies.name }}</div>
              <div style="font-size:13px;color:var(--md-text-secondary);">{{ selectedSpecies.seriesId }}系列</div>
            </div>
            <button @click="closeDetail" style="margin-left:auto;width:32px;height:32px;border-radius:50%;border:1px solid rgba(255,255,255,0.1);background:rgba(255,255,255,0.04);color:rgba(255,255,255,0.5);font-size:14px;cursor:pointer;">✕</button>
          </div>

          <!-- 阶段Tab -->
          <div style="display:flex;gap:6px;margin-bottom:20px;flex-wrap:wrap;">
            <button v-for="(name, i) in STAGE_NAMES" :key="i" @click="detailStage = i"
              :style="{
                padding: '6px 14px', borderRadius: '20px', border: '1px solid rgba(255,255,255,0.06)',
                background: detailStage === i ? 'rgba(167,139,250,0.2)' : 'rgba(255,255,255,0.03)',
                color: detailStage === i ? 'var(--md-primary-light)' : 'var(--md-text-secondary)',
                fontWeight: detailStage === i ? 700 : 400, fontSize: '12px', cursor: 'pointer', fontFamily: 'inherit',
              }">{{ name }}</button>
          </div>

          <!-- 阶段详情 -->
          <div style="margin-bottom:16px;">
            <div style="font-size:14px;font-weight:600;margin-bottom:8px;">
              Lv.{{ detailStage === 0 ? 1 : detailStage <= 2 ? 2 : detailStage <= 3 ? 8 : detailStage === 4 ? 10 : 12 }}
              阶段
            </div>
            <div style="padding:14px 16px;background:rgba(255,255,255,0.03);border-radius:12px;border-left:3px solid var(--md-primary);margin-bottom:12px;">
              <div style="font-size:13px;color:var(--md-text-secondary);line-height:1.6;">
                {{ detailSpecies.levels[detailStage]?.description || '待完善' }}
              </div>
            </div>

            <!-- 进化台词 -->
            <div v-if="detailEvoLines[detailStage]" style="padding:12px 16px;background:rgba(245,158,11,0.06);border-radius:12px;border-left:3px solid #F59E0B;margin-bottom:12px;">
              <div style="font-size:11px;color:rgba(245,158,11,0.6);font-weight:600;margin-bottom:4px;">💬 进化台词</div>
              <div style="font-size:16px;color:#fcd34d;font-style:italic;">{{ detailEvoLines[detailStage] }}</div>
            </div>

            <!-- 诗歌 -->
            <div v-if="detailPoems[detailStage]" style="padding:14px 16px;background:rgba(139,92,246,0.06);border-radius:12px;border-left:3px solid #8B5CF6;">
              <div style="font-size:11px;color:rgba(139,92,246,0.6);font-weight:600;margin-bottom:4px;">📜 专属诗文</div>
              <div style="font-size:14px;color:#c4b5fd;line-height:1.8;white-space:pre-line;">{{ detailPoems[detailStage] }}</div>
            </div>
          </div>

          <!-- 进化链指示器 -->
          <div style="display:flex;justify-content:center;gap:4px;margin-top:12px;">
            <span v-for="i in 12" :key="i"
              :style="{
                width:'12px',height:'12px',borderRadius:'50%',
                background: i === 12 ? '#f472b6' : (i <= 3 ? 'rgba(252,211,77,0.3)' : 'rgba(255,255,255,0.1)'),
              }"></span>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 图鉴轮播 -->
    <div style="position:relative;background:rgba(255,255,255,0.02);border-radius:var(--md-radius);padding:24px;border:1px solid rgba(255,255,255,0.04);overflow:hidden;min-height:400px;">
      <div v-for="(series, si) in seriesList" :key="series.id"
        :style="{ display: currentSlide === si ? 'block' : 'none', animation: 'fadeIn 0.5s ease' }">
        <div style="display:flex;align-items:center;gap:16px;margin-bottom:16px;">
          <span style="font-size:48px;">{{ series.emoji }}</span>
          <span style="font-size:28px;font-weight:700;">{{ series.name }}</span>
        </div>
        <div style="color:var(--md-text-secondary);font-size:14px;margin-bottom:20px;padding:12px 16px;background:rgba(255,255,255,0.03);border-radius:12px;border-left:3px solid var(--md-primary);">
          {{ series.species.length }}种宠物 · 点击物种查看详细进化信息
        </div>
        <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:16px;">
          <div v-for="sp in series.species" :key="sp.id" @click="openDetail(sp.id, sp.name)"
            style="background:rgba(255,255,255,0.03);border-radius:12px;padding:14px 12px;text-align:center;border:1px solid rgba(255,255,255,0.04);transition:0.2s;cursor:pointer;">
            <span style="font-size:36px;display:block;margin-bottom:4px;">{{ getSpeciesEmoji(sp.id) }}</span>
            <div style="font-weight:600;font-size:16px;">{{ sp.name }}</div>
            <div style="font-size:12px;color:var(--md-text-secondary);margin-top:4px;line-height:1.4;overflow:hidden;text-overflow:ellipsis;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;">
              {{ sp.levels[0]?.description || '' }}
            </div>
            <div style="display:flex;justify-content:center;gap:2px;margin-top:8px;">
              <span v-for="lvl in 12" :key="lvl"
                :style="{
                  display:'inline-block', width:'8px', height:'8px', borderRadius:'50%',
                  background: lvl === 12 ? '#f472b6' : (lvl <= 3 ? 'rgba(252,211,77,0.3)' : 'rgba(255,255,255,0.1)'),
                  boxShadow: lvl === 12 ? '0 0 8px rgba(244,114,182,0.3)' : 'none',
                }"></span>
            </div>
          </div>
        </div>
      </div>

      <div style="display:flex;justify-content:center;gap:8px;margin-top:24px;">
        <button v-for="(_, i) in seriesList" :key="i"
          :style="{
            width: currentSlide === i ? '32px' : '12px', height:'12px',
            borderRadius: currentSlide === i ? '6px' : '50%',
            border: 'none', cursor:'pointer',
            background: currentSlide === i ? 'var(--md-primary)' : 'rgba(255,255,255,0.1)',
            boxShadow: currentSlide === i ? '0 0 12px rgba(167,139,250,0.3)' : 'none',
            transition:'0.2s',
          }"
          @click="goToSlide(i)"></button>
      </div>
    </div>
  </div>
</template>

<style>
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
