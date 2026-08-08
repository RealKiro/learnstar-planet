<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiPost } from '@/utils/api'
import { useAppMode } from '@/composables/useAppMode'
import { getAllSeries, PET_SERIES, getSpeciesById } from '@/utils/petData'
import { getPoems, getEvoLines, STAGE_NAMES, poemToLines } from '@/utils/petHandbookData'
import { getPetLifeStory } from '@/utils/petLifeStories'
import PetSprite from '@/components/pet/PetSprite.vue'

// 教室端/教师端共用图鉴页；教师端无 class_token，隐藏"切换系列"
const { isClassroomMode } = useAppMode()

/** 阶段 tab → 展示等级（与详情 Lv 显示一致） */
function stageLevel(idx: number): number {
  return idx === 0 ? 1 : idx <= 2 ? 2 : idx <= 3 ? 8 : idx === 4 ? 10 : 12
}

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

/** 角色人生档案（故事演义：品性/行为/服饰/功法/台词/诗词/年龄/关键词） */
const lifeStory = computed(() => selectedSpecies.value ? getPetLifeStory(selectedSpecies.value.speciesId) : null)
const lifeStage = computed(() => lifeStory.value?.stages?.[detailStage.value])

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
      <!-- 切换系列仅教室端可用（教师端无 class_token） -->
      <div v-if="isClassroomMode" style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
        <label style="font-size:13px;color:var(--md-text-secondary);">🏷️ 系列</label>
        <select v-model="currentSeries"
          style="padding:8px 14px;border-radius:var(--md-radius);background:var(--tint-3);border:1px solid var(--tint-3);color:var(--color-text);font-size:14px;font-weight:500;outline:none;cursor:pointer;font-family:inherit;">
          <option v-for="s in allSeries" :key="s.id" :value="s.id" style="background:#1a1a2e;color:#f1f1f1;">{{ s.emoji }} {{ s.name }}</option>
        </select>
        <button @click="switchSeries" :disabled="switching"
          style="padding:8px 18px;border-radius:30px;border:none;background:rgba(167,139,250,0.15);color:var(--md-primary-light);font-size:13px;font-weight:600;cursor:pointer;transition:0.15s;font-family:inherit;">
          {{ switching ? '切换中...' : '切换（每人扣20分）' }}
        </button>
      </div>
    </div>

    <!-- 切换消息 -->
    <div v-if="switchMsg" style="margin-bottom:16px;padding:10px 16px;background:rgba(16,185,129,0.08);border:1px solid rgba(16,185,129,0.15);border-radius:var(--md-radius);color: var(--color-success-text);font-size:13px;">{{ switchMsg }}</div>
    <div v-if="switchError" style="margin-bottom:16px;padding:10px 16px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.15);border-radius:var(--md-radius);color: var(--color-danger-text);font-size:13px;">{{ switchError }}</div>

    <!-- 物种大图（点击物种后展示） -->
    <Transition name="fade">
      <div v-if="selectedSpecies && detailSpecies" @click.self="closeDetail"
        style="position:fixed;inset:0;z-index:300;background:rgba(5,2,20,0.85);backdrop-filter:blur(16px);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:linear-gradient(180deg,var(--color-bg-card),var(--color-bg));border:1px solid var(--tint-3);border-radius:24px;max-width:700px;width:100%;max-height:85vh;overflow-y:auto;padding:28px;box-shadow:0 20px 60px rgba(0,0,0,0.5);">
          <!-- 头部：程序化宠物艺术，随阶段 tab 切换展示不同进化形态 -->
          <div style="display:flex;align-items:center;gap:16px;margin-bottom:20px;">
            <div class="detail-orb">
              <PetSprite :species-id="selectedSpecies.speciesId" :level="stageLevel(detailStage)" :animate="true" />
            </div>
            <div>
              <div style="font-size:22px;font-weight:700;">{{ selectedSpecies.name }}</div>
              <div style="font-size:13px;color:var(--md-text-secondary);">{{ selectedSpecies.seriesId }}系列</div>
            </div>
            <button @click="closeDetail" style="margin-left:auto;width:32px;height:32px;border-radius:50%;border:1px solid var(--tint-4);background:var(--tint-2);color:var(--color-text-secondary);font-size:14px;cursor:pointer;">✕</button>
          </div>

          <!-- 主题句（人生档案·故事演义） -->
          <div v-if="lifeStory?.theme" class="life-theme">「{{ lifeStory.theme }}」</div>

          <!-- 阶段Tab -->
          <div style="display:flex;gap:6px;margin-bottom:20px;flex-wrap:wrap;">
            <button v-for="(name, i) in STAGE_NAMES" :key="i" @click="detailStage = i"
              :style="{
                padding: '6px 14px', borderRadius: '20px', border: '1px solid var(--tint-3)',
                background: detailStage === i ? 'rgba(167,139,250,0.2)' : 'var(--tint-1)',
                color: detailStage === i ? 'var(--color-primary)' : 'var(--md-text-secondary)',
                fontWeight: detailStage === i ? 700 : 400, fontSize: '12px', cursor: 'pointer', fontFamily: 'inherit',
              }">{{ name }}</button>
          </div>

          <!-- 阶段详情 -->
          <div style="margin-bottom:16px;">
            <!-- 阶段标题：人生档案阶段名 + 年龄 + 关键词 -->
            <div class="life-stage-head">
              <span class="life-stage-name">{{ lifeStage?.name || ((detailStage === 0 ? 1 : detailStage <= 2 ? 2 : detailStage <= 3 ? 8 : detailStage === 4 ? 10 : 12) + ' 阶段') }}</span>
              <template v-if="lifeStage">
                <span class="life-age">{{ lifeStage.age }}</span>
                <span class="life-keyword">{{ lifeStage.keyword }}</span>
              </template>
            </div>

            <!-- 有角色人生档案：展示故事演义（品性/行为/服饰/功法/台词/诗文） -->
            <template v-if="lifeStage">
              <div class="life-block life-char">
                <div class="life-label">🎭 品性</div>
                <div class="life-text">{{ lifeStage.character }}</div>
              </div>
              <div class="life-block life-action">
                <div class="life-label">🏃 行为</div>
                <div class="life-text">{{ lifeStage.action }}</div>
              </div>
              <div class="life-pills">
                <span class="life-pill">👘 {{ lifeStage.attire }}</span>
                <span class="life-pill pill-skill">⚡ {{ lifeStage.technique }}</span>
              </div>
              <div v-if="lifeStage.line" class="life-block life-quote">
                <div class="life-label">💬 台词</div>
                <div class="life-text quote">{{ lifeStage.line }}</div>
              </div>
              <div v-if="lifeStage.poem" class="life-block life-poem">
                <div class="life-label">📜 诗文</div>
                <div class="life-text poem">
                  <div v-for="(line, i) in poemToLines(lifeStage.poem)" :key="i" :class="i % 2 === 1 ? 'poem-line poem-line--second' : 'poem-line'">{{ line }}</div>
                </div>
              </div>
            </template>

            <!-- 无档案兜底：物种等级描述 + 台词 + 诗文 -->
            <template v-else>
              <div class="life-block life-desc">
                <div class="life-text">{{ detailSpecies.levels[detailStage]?.description || '待完善' }}</div>
              </div>
              <div v-if="detailEvoLines[detailStage]" class="life-block life-quote">
                <div class="life-label">💬 进化台词</div>
                <div class="life-text quote">{{ detailEvoLines[detailStage] }}</div>
              </div>
              <div v-if="detailPoems[detailStage]" class="life-block life-poem">
                <div class="life-label">📜 专属诗文</div>
                <div class="life-text poem">
                  <div v-for="(line, i) in poemToLines(detailPoems[detailStage] || '')" :key="i" :class="i % 2 === 1 ? 'poem-line poem-line--second' : 'poem-line'">{{ line }}</div>
                </div>
              </div>
            </template>
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
    <div style="position:relative;background:var(--tint-1);border-radius:var(--md-radius);padding:24px;border:1px solid var(--tint-2);overflow:hidden;min-height:400px;">
      <div v-for="(series, si) in seriesList" :key="series.id"
        :style="{ display: currentSlide === si ? 'block' : 'none', animation: 'fadeIn 0.5s ease' }">
        <div style="display:flex;align-items:center;gap:16px;margin-bottom:16px;">
          <span style="font-size:48px;">{{ series.emoji }}</span>
          <span style="font-size:28px;font-weight:700;">{{ series.name }}</span>
        </div>
        <div style="color:var(--md-text-secondary);font-size:14px;margin-bottom:20px;padding:12px 16px;background:var(--tint-1);border-radius:12px;border-left:3px solid var(--md-primary);">
          {{ series.species.length }}种宠物 · 点击物种查看详细进化信息
        </div>
        <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:16px;">
          <div v-for="sp in series.species" :key="sp.id" @click="openDetail(sp.id, sp.name)" class="species-card">
            <div class="species-orb">
              <PetSprite :species-id="sp.id" :level="6" />
            </div>
            <div class="species-card-name">{{ sp.name }}</div>
            <div class="species-card-desc">{{ sp.levels[0]?.description || '' }}</div>
            <div class="species-card-dots">
              <span v-for="lvl in 12" :key="lvl"
                :style="{
                  background: lvl === 12 ? '#f472b6' : (lvl <= 3 ? 'rgba(252,211,77,0.3)' : 'var(--tint-4)'),
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
            background: currentSlide === i ? 'var(--md-primary)' : 'var(--tint-4)',
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

/* ===== 详情弹窗宠物星环 ===== */
.detail-orb {
  width: 104px;
  height: 104px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--tint-2);
  border: 2px solid color-mix(in srgb, var(--md-primary) 40%, var(--tint-3));
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-primary) 9%, transparent), 0 8px 22px rgba(0,0,0,0.2);
}

/* ===== 物种卡片（现代化打磨） ===== */
.species-card {
  background: var(--tint-1);
  border-radius: 14px;
  padding: 16px 12px 14px;
  text-align: center;
  border: 1px solid var(--tint-2);
  transition: all 0.2s ease;
  cursor: pointer;
}
.species-card:hover {
  transform: translateY(-3px);
  border-color: color-mix(in srgb, var(--md-primary) 45%, var(--tint-3));
  box-shadow: 0 10px 24px rgba(0,0,0,0.16), 0 2px 8px rgba(0,0,0,0.06);
}
.species-orb {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 auto 8px;
  background: var(--tint-2);
  border: 2px solid color-mix(in srgb, var(--md-primary) 35%, var(--tint-3));
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-primary) 8%, transparent), 0 6px 16px rgba(0,0,0,0.12);
  transition: transform 0.2s ease;
}
.species-card:hover .species-orb { transform: scale(1.05); }
.species-card-name { font-weight: 700; font-size: 16px; color: var(--color-text); }
.species-card-desc {
  font-size: 12px;
  color: var(--md-text-secondary);
  margin-top: 4px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.species-card-dots { display: flex; justify-content: center; gap: 2px; margin-top: 8px; }
.species-card-dots span { display: inline-block; width: 8px; height: 8px; border-radius: 50%; }

/* ===== 详情弹窗·角色人生档案（故事演义） ===== */
.life-theme {
  margin: -6px 0 14px;
  padding: 10px 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(167,139,250,0.1), rgba(244,114,182,0.05));
  border: 1px solid rgba(167,139,250,0.18);
  font-size: 14px;
  font-weight: 600;
  color: var(--md-gold);
  line-height: 1.6;
}
.life-stage-head {
  display: flex;
  align-items: baseline;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}
.life-stage-name { font-size: 17px; font-weight: 800; color: var(--color-text); }
.life-age {
  font-size: 11px;
  font-weight: 700;
  color: var(--color-primary);
  background: rgba(167,139,250,0.12);
  padding: 2px 10px;
  border-radius: 12px;
}
.life-keyword {
  font-size: 12px;
  font-weight: 600;
  color: var(--md-text-secondary);
}
.life-block {
  padding: 12px 14px;
  border-radius: 12px;
  margin-bottom: 10px;
}
.life-char { background: rgba(79,70,229,0.06); border-left: 3px solid var(--color-primary); }
.life-action { background: rgba(16,185,129,0.06); border-left: 3px solid #10B981; }
.life-desc { background: var(--tint-1); border-left: 3px solid var(--md-primary); }
.life-quote { background: rgba(245,158,11,0.06); border-left: 3px solid #F59E0B; }
.life-poem { background: rgba(139,92,246,0.06); border-left: 3px solid #8B5CF6; }
.life-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--md-text-secondary);
  margin-bottom: 5px;
}
.life-text {
  font-size: 13px;
  color: var(--color-text);
  line-height: 1.7;
}
.life-text.quote { font-size: 15px; font-style: italic; color: var(--color-warning-text); }
.life-text.poem { color: var(--color-primary); line-height: 1.8; }
.poem-line { display: block; }
.poem-line--second { padding-left: 1.4em; text-indent: -1.4em; }
.life-pills { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 10px; }
.life-pill {
  font-size: 12px;
  font-weight: 600;
  padding: 5px 12px;
  border-radius: 14px;
  background: var(--tint-1);
  border: 1px solid var(--tint-3);
  color: var(--color-text);
}
.life-pill.pill-skill { color: var(--color-primary); border-color: rgba(167,139,250,0.3); }
</style>
