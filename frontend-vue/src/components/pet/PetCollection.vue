<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import {
  getAllSeries,
  getPetLevelName,
  SERIES_SCENES,
  getSeriesName,
} from '@/utils/petData'
import PetSprite from './PetSprite.vue'
import type { ApiResponse } from '@/types'

const props = defineProps<{
  visible: boolean
  studentId: number
}>()

const emit = defineEmits<{
  close: []
  switched: []
}>()

interface CollectionItem {
  species: string
  level: number
  experience: number
  mood: number
  is_active: boolean
}

interface CollectionData {
  student_id: number
  student_name: string
  total_score: number
  unlock_slots: number
  class_series: string | null
  active_species: string | null
  collection: CollectionItem[]
}

/** 更换宠物积分成本（与后端 Pet::switchCost 一致：等级越高越贵） */
function switchCost(level: number): number {
  return 5 * Math.max(1, level)
}

const loading = ref(false)
const data = ref<CollectionData | null>(null)
const activeSeries = ref<string | 'all'>('all')
const switching = ref('')
const message = ref<{ type: 'ok' | 'err'; text: string } | null>(null)

const seriesList = computed(() => getAllSeries())

/** 当前类别允许领养的物种 id 集合（未设置类别则不限制） */
const categorySpeciesIds = computed<Set<string>>(() => {
  const series = data.value?.class_series
  if (!series) return new Set(seriesList.value.flatMap(s => s.species.map(sp => sp.id)))
  const found = seriesList.value.find(s => s.id === series)
  return new Set((found?.species || []).map(sp => sp.id))
})

/** 某物种是否属于本班当前类别（可领养） */
function inCategory(speciesId: string): boolean {
  return categorySpeciesIds.value.has(speciesId)
}

/** 当前类别名称 */
const categoryName = computed(() => {
  const s = data.value?.class_series
  return s ? getSeriesName(s) : '未设置(全类别)'
})

/** 已解锁物种 id -> 收藏项 */
const ownedMap = computed(() => {
  const map: Record<string, CollectionItem> = {}
  for (const c of data.value?.collection || []) {
    map[c.species] = c
  }
  return map
})

/** 当前筛选系列下的物种 */
const shownSpecies = computed(() => {
  const flat: Array<{ id: string; name: string; seriesId: string }> = []
  for (const s of seriesList.value) {
    if (activeSeries.value !== 'all' && s.id !== activeSeries.value) continue
    for (const sp of s.species) flat.push({ id: sp.id, name: sp.name, seriesId: s.id })
  }
  return flat
})

const ownedCount = computed(() => Object.keys(ownedMap.value).length)

function isOwned(speciesId: string): boolean {
  return !!ownedMap.value[speciesId]
}

function speciesLevel(speciesId: string): number {
  return ownedMap.value[speciesId]?.level || 1
}

function seriesColor(seriesId: string): string {
  return SERIES_SCENES[seriesId]?.primaryColor || '#6366F1'
}

async function load() {
  if (!props.visible) return
  loading.value = true
  message.value = null
  try {
    const res = await apiGet<ApiResponse<CollectionData>>(`/api/v1/teacher/pets/${props.studentId}/collection`)
    data.value = res.data
  } catch {
    data.value = {
      student_id: props.studentId,
      student_name: '该学生',
      total_score: 0,
      unlock_slots: 1,
      class_series: null,
      active_species: null,
      collection: [],
    }
    message.value = { type: 'err', text: '加载失败（离线模式仅展示框架）' }
  } finally {
    loading.value = false
  }
}

async function doSwitch(speciesId: string) {
  if (switching.value) return
  switching.value = speciesId
  message.value = null
  try {
    await apiPost(`/api/v1/teacher/pets/${props.studentId}/switch`, { pet_species: speciesId })
    message.value = { type: 'ok', text: '✅ 切换成功！进度已保留在图鉴中' }
    await load()
    emit('switched')
  } catch (e: any) {
    message.value = { type: 'err', text: e?.response?.data?.message || e?.message || '切换失败' }
  } finally {
    switching.value = ''
  }
}

onMounted(load)
watch(() => props.visible, load)
</script>

<template>
  <Teleport to="body">
    <Transition name="collection-fade">
      <div v-if="visible" class="collection-overlay" @click.self="emit('close')">
        <div class="collection-modal">
          <!-- 头部 -->
          <div class="collection-header">
            <div class="header-info">
              <div class="header-icon">📖</div>
              <div>
                <h2 class="header-title">宠物图鉴 · {{ data?.student_name || '学生' }}</h2>
                <p class="header-sub">
                  已解锁 <strong>{{ ownedCount }}</strong> / 120 种 · 当前类别 <strong class="cat-name">{{ categoryName }}</strong>
                  <span v-if="data?.class_series" class="slots-full">(仅本类可领养)</span>
                </p>
              </div>
            </div>
            <button class="close-btn" @click="emit('close')">✕</button>
          </div>

          <!-- 提示 -->
          <div class="collection-tip">
            💡 每积累 100 积分解锁一个新图鉴槽位；切换宠物会完整保留原宠物进度。
          </div>
          <div v-if="message" class="collection-message" :class="'msg--' + message.type">
            {{ message.text }}
          </div>

          <!-- 系列筛选 -->
          <div class="series-filter">
            <button class="filter-chip" :class="{ 'chip--active': activeSeries === 'all' }" @click="activeSeries = 'all'">
              🌟 全部
            </button>
            <button
              v-for="s in seriesList"
              :key="s.id"
              class="filter-chip"
              :class="{ 'chip--active': activeSeries === s.id }"
              :style="activeSeries === s.id ? { background: seriesColor(s.id) + '22', borderColor: seriesColor(s.id) + '66' } : {}"
              @click="activeSeries = s.id"
            >
              {{ s.emoji }} {{ s.name }}
            </button>
          </div>

          <!-- 网格 -->
          <div v-if="loading" class="loading-state">
            <div class="loading-spinner"></div>
            <p>正在翻阅图鉴...</p>
          </div>
          <div v-else class="species-grid">
            <div
              v-for="sp in shownSpecies"
              :key="sp.id"
              class="species-card"
              :class="{
                'card--owned': isOwned(sp.id),
                'card--active': data?.active_species === sp.id,
              }"
              :style="{ '--card-accent': seriesColor(sp.seriesId) }"
            >
              <div class="card-top">
                <div class="sprite">
                  <PetSprite :species-id="sp.id" :level="speciesLevel(sp.id)" :mood="80" :animate="false" />
                </div>
                <div class="status-badge" :class="isOwned(sp.id) ? 'badge--owned' : 'badge--locked'">
                  {{ isOwned(sp.id) ? (data?.active_species === sp.id ? '当前' : '✓') : '🔒' }}
                </div>
              </div>
              <div class="card-name">{{ sp.name }}</div>
              <div class="card-series" :style="{ color: seriesColor(sp.seriesId) }">{{ getSeriesName(sp.seriesId) }}</div>
              <div class="card-level" v-if="isOwned(sp.id)">最高 Lv.{{ speciesLevel(sp.id) }} · {{ getPetLevelName(sp.id, speciesLevel(sp.id)) }}</div>
              <div class="card-locked-hint" v-else-if="inCategory(sp.id)">需解锁图鉴槽位</div>
              <div class="card-locked-hint hint--cross" v-else>跨类别 · 需教师整班切换</div>
              <button
                v-if="isOwned(sp.id) && inCategory(sp.id) && data?.active_species !== sp.id"
                class="switch-btn"
                :disabled="switching !== ''"
                @click="doSwitch(sp.id)"
              >
                {{ switching === sp.id ? '更换中…' : `积分更换(${switchCost(speciesLevel(sp.id))})` }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.collection-overlay {
  position: fixed;
  inset: 0;
  z-index: 420;
  background: rgba(5, 2, 20, 0.85);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.collection-modal {
  width: 100%;
  max-width: 860px;
  max-height: 88vh;
  background: linear-gradient(180deg, var(--color-bg-card) 0%, var(--color-bg) 100%);
  border: 1px solid var(--tint-3);
  border-radius: 24px;
  display: flex;
  flex-direction: column;
  padding: 24px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.5);
}
.collection-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}
.header-info { display: flex; align-items: center; gap: 12px; }
.header-icon { font-size: 30px; }
.header-title { font-size: 18px; font-weight: 700; color: var(--color-text); margin: 0; }
.header-sub { font-size: 12px; color: var(--color-text-secondary); margin: 2px 0 0; }
.header-sub strong { color: #F59E0B; }
.slots-full { color: #EF4444; font-weight: 600; }
.close-btn {
  width: 32px; height: 32px; border-radius: 50%;
  border: 1px solid var(--tint-4); background: var(--tint-2);
  color: var(--color-text-secondary); font-size: 14px; cursor: pointer; display: flex;
  align-items: center; justify-content: center; transition: all 0.2s;
}
.close-btn:hover { background: var(--tint-4); color: var(--color-text); }
.collection-tip {
  font-size: 12px; color: var(--color-text-secondary);
  background: rgba(99,102,241,0.08); border: 1px solid rgba(99,102,241,0.15);
  border-radius: 10px; padding: 8px 12px; margin-bottom: 8px;
}
.collection-message {
  font-size: 13px; padding: 8px 12px; border-radius: 10px; margin-bottom: 8px; font-weight: 600;
}
.msg--ok { background: rgba(16,185,129,0.12); color: #34D399; }
.msg--err { background: rgba(239,68,68,0.12); color: #F87171; }
.series-filter { display: flex; gap: 6px; margin-bottom: 14px; flex-wrap: wrap; }
.filter-chip {
  padding: 5px 12px; border-radius: 20px; font-size: 12px; font-weight: 500;
  cursor: pointer; border: 1px solid var(--tint-4); background: var(--tint-2);
  color: var(--color-text-secondary); transition: all 0.2s;
}
.filter-chip:hover { color: var(--color-text); border-color: var(--tint-4); }
.chip--active { color: var(--color-text); font-weight: 600; border-color: currentColor; }

.species-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 10px;
  overflow-y: auto;
  padding-right: 4px;
}
.species-grid::-webkit-scrollbar { width: 4px; }
.species-grid::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 2px; }
.species-card {
  background: var(--tint-1);
  border: 1px solid var(--tint-3);
  border-radius: 14px;
  padding: 12px 8px 10px;
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  transition: all 0.2s;
  position: relative;
}
.species-card:hover { background: var(--tint-3); transform: translateY(-2px); }
.card--owned { border-color: rgba(16,185,129,0.15); }
.card--active { border-color: var(--card-accent, #F59E0B); box-shadow: 0 0 12px color-mix(in srgb, var(--card-accent, #F59E0B) 20%, transparent); }
.card-top { position: relative; width: 100%; display: flex; justify-content: center; }
.sprite { width: 56px; height: 56px; }
.status-badge {
  position: absolute; top: -2px; right: -2px;
  width: 22px; height: 22px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 700;
}
.badge--owned { background: rgba(16,185,129,0.25); color: #34D399; }
.badge--locked { background: var(--tint-3); }
.card-name { font-size: 12px; font-weight: 600; color: var(--color-text); }
.card-series { font-size: 10px; opacity: 0.7; }
.card-level { font-size: 10px; color: var(--color-text-secondary); text-align: center; }
.card-locked-hint { font-size: 10px; color: var(--color-text-secondary); }
.hint--cross { color: rgba(239,68,68,0.75); }
.cat-name { color: #F59E0B; }
.switch-btn {
  margin-top: 4px; padding: 4px 16px; border-radius: 12px; font-size: 12px; font-weight: 600;
  border: none; cursor: pointer; color: white;
  background: var(--card-accent, #F59E0B);
  transition: all 0.2s;
}
.switch-btn:hover { filter: brightness(1.15); transform: scale(1.03); }
.switch-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.loading-state { text-align: center; padding: 40px; color: var(--color-text-secondary); }
.loading-spinner {
  width: 32px; height: 32px; border: 3px solid var(--tint-4);
  border-top-color: var(--color-primary, #6366F1); border-radius: 50%;
  animation: spin 0.8s linear infinite; margin: 0 auto 10px;
}
@keyframes spin { to { transform: rotate(360deg); } }

.collection-fade-enter-active { animation: overlayIn 0.3s ease; }
.collection-fade-leave-active { animation: overlayIn 0.2s ease reverse; }
@keyframes overlayIn { from { opacity: 0; transform: scale(0.98); } to { opacity: 1; transform: scale(1); } }
</style>
