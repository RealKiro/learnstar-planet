<script setup lang="ts">
import { computed } from 'vue'
import {
  getSpeciesById,
  getSeriesBySpeciesId,
  getSpeciesEmoji,
  getLevelRequiredScore,
  getPetLevelName,
  getPetLevelDescription,
  SERIES_SCENES,
} from '@/utils/petData'
import { getStagePersonality, getStageAbility } from '@/utils/petTraits'
import { getPoems, getEvoLines, stageIndexForLevel, poemToLines } from '@/utils/petHandbookData'
import { useThemeStore } from '@/stores/theme'
import PetSprite from './PetSprite.vue'

const props = withDefaults(defineProps<{
  speciesId: string
  /** 当前等级 */
  level: number
  /** 当前累计积分（用于升级进度） */
  score?: number
  /** 是否显示「切换宠物」按钮（课堂评价场景用） */
  showPetSwitch?: boolean
}>(), {
  score: 0,
  showPetSwitch: false,
})

const emit = defineEmits<{
  close: []
  'switch-pet': []
}>()

const species = computed(() => getSpeciesById(props.speciesId))
const series = computed(() => getSeriesBySpeciesId(props.speciesId))
const scene = computed(() => (series.value ? SERIES_SCENES[series.value.id] : null))
const themeStore = useThemeStore()

const stageIdx = computed(() => stageIndexForLevel(props.level))
const currentStage = computed(() => getLevelDataStage(props.speciesId, props.level))
const currentLevelName = computed(() => getPetLevelName(props.speciesId, props.level) || species.value?.name || '')
const currentLevelDescription = computed(() => getPetLevelDescription(props.speciesId, props.level) || '')
const stageLabel = computed(() => STAGE_LABEL[currentStage.value] || '')

const evoLine = computed(() => {
  const lines = getEvoLines(species.value?.name || '')
  return lines[stageIdx.value] || lines[lines.length - 1] || ''
})
const poem = computed(() => {
  const poems = getPoems(species.value?.name || '')
  return poems[stageIdx.value] || poems[poems.length - 1] || ''
})
const poemLines = computed(() => poemToLines(poem.value))
const personality = computed(() => getStagePersonality(props.speciesId, props.level))
const ability = computed(() => getStageAbility(props.speciesId, props.level))

// 六阶段进化全览：每个 stage 取该物种第一个等级（静态展示，当前高亮）
const evoStages = computed(() => {
  const seen = new Set<string>()
  const arr: Array<{ stage: string; level: number; name: string; unlocked: boolean }> = []
  for (const l of species.value?.levels || []) {
    if (seen.has(l.stage)) continue
    seen.add(l.stage)
    arr.push({ stage: l.stage, level: l.level, name: l.name, unlocked: l.level <= props.level })
  }
  return arr
})

// 当前阶段背景（日/夜系列渐变）
const sceneBg = computed(() => {
  const primary = scene.value?.primaryColor || '#6366F1'
  if (themeStore.isDark) return scene.value?.bgGradient || `linear-gradient(180deg, ${primary}, #0d1b2a)`
  return `linear-gradient(180deg, ${primary}2e, var(--color-bg-card))`
})

// 升级进度
const progress = computed(() => {
  const lv = props.level
  if (lv >= 12) return { remaining: 0, percent: 100 }
  const cur = getLevelRequiredScore(lv)
  const next = getLevelRequiredScore(lv + 1)
  const exp = Math.max(0, props.score - cur)
  const range = next - cur
  return {
    remaining: Math.max(0, next - props.score),
    percent: Math.min(100, range > 0 ? Math.round((exp / range) * 100) : 100),
  }
})

function getLevelDataStage(speciesId: string, level: number): string {
  return getSpeciesById(speciesId)?.levels.find(l => l.level === level)?.stage || 'egg'
}

const STAGE_LABEL: Record<string, string> = {
  egg: '新生', baby: '幼年', growing: '成长期', mature: '成熟期', legendary: '传说级', transcendent: '道果',
}
</script>

<template>
  <div class="detail-overlay" @click.self="emit('close')">
    <div class="detail-modal">
      <!-- 头部 -->
      <div class="detail-header">
        <span class="detail-emoji">{{ getSpeciesEmoji(speciesId) }}</span>
        <div class="header-info">
          <h2 class="detail-title">{{ currentLevelName }} · 角色介绍</h2>
          <p class="detail-series">{{ series?.name }}系列 · Lv.{{ level }}</p>
        </div>
        <button class="close-btn" @click="emit('close')">✕</button>
      </div>

      <!-- 六阶段进化全览（静态展示，当前高亮） -->
      <div class="evo-track">
        <div
          v-for="(st, i) in evoStages"
          :key="st.stage"
          class="evo-node"
          :class="{ 'evo-node--active': st.stage === currentStage, 'evo-node--locked': !st.unlocked }"
        >
          <div class="evo-icon">{{ st.unlocked ? '⭐' : '🔒' }}</div>
          <div class="evo-name">{{ st.name }}</div>
          <div class="evo-lv">Lv.{{ st.level }}</div>
          <div v-if="i < evoStages.length - 1" class="evo-connector"></div>
        </div>
      </div>

      <!-- 当前形态展示 -->
      <div class="current-area">
        <div class="current-scene" :style="{ background: sceneBg }">
          <div class="current-sprite">
            <PetSprite :species-id="speciesId" :level="level" :animate="true" />
          </div>
          <div class="current-badge">{{ stageLabel }} · Lv.{{ level }}</div>
        </div>

        <div class="current-detail">
          <div class="current-name">
            {{ currentLevelName }}
            <span class="current-stage">{{ stageLabel }}形态</span>
          </div>
          <p class="current-desc">{{ currentLevelDescription }}</p>

          <!-- 台词 -->
          <div v-if="evoLine" class="info-block quote-block">
            <div class="info-label">💬 台词</div>
            <div class="info-text quote-text">{{ evoLine }}</div>
          </div>

          <!-- 诗文 -->
          <div v-if="poem" class="info-block poem-block">
            <div class="info-label">📜 诗文</div>
            <div class="info-text poem-text">
              <span v-for="(line, i) in poemLines" :key="i" class="poem-line" :class="{ 'poem-line--second': i % 2 === 1 }">{{ line }}</span>
            </div>
          </div>

          <!-- 特质与招式 -->
          <div class="info-row">
            <span class="info-pill">🎭 {{ personality }}</span>
            <span class="info-pill pill--ability">⚡ {{ ability }}</span>
          </div>

          <!-- 升级进度 -->
          <div class="progress-block">
            <div class="progress-head">
              <span>升级进度</span>
              <span v-if="level < 12">距 Lv.{{ level + 1 }} 还差 <strong>{{ progress.remaining }}</strong> 分</span>
              <span v-else>已达满级</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: progress.percent + '%' }"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部按钮 -->
      <div class="detail-footer">
        <button v-if="showPetSwitch" class="foot-btn foot-switch" @click="emit('switch-pet')">🔄 切换宠物</button>
        <button class="foot-btn foot-close" @click="emit('close')">关闭</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-overlay {
  position: fixed;
  inset: 0;
  z-index: 400;
  background: rgba(5, 2, 20, 0.85);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.detail-modal {
  width: 100%;
  max-width: 620px;
  max-height: 86vh;
  overflow-y: auto;
  background: linear-gradient(180deg, var(--color-bg-card) 0%, var(--color-bg) 100%);
  border: 1px solid var(--tint-3);
  border-radius: 24px;
  padding: 24px 28px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  animation: detailPop 0.25s ease;
}
@keyframes detailPop {
  from { transform: scale(0.94); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

/* 头部 */
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 18px;
}
.detail-emoji { font-size: 34px; }
.header-info { min-width: 0; }
.detail-title { font-size: 20px; font-weight: 700; margin: 0; }
.detail-series { font-size: 12px; color: var(--color-text-secondary); margin: 2px 0 0; }
.close-btn {
  margin-left: auto;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: 1px solid var(--tint-3);
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: 14px;
  flex-shrink: 0;
}
.close-btn:hover { background: var(--tint-1); color: var(--color-text); }

/* 六阶段进化全览 */
.evo-track {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 4px;
  padding: 14px 12px;
  background: var(--tint-1);
  border: 1px solid var(--tint-2);
  border-radius: 16px;
  margin-bottom: 18px;
  overflow-x: auto;
}
.evo-node {
  position: relative;
  text-align: center;
  min-width: 52px;
  opacity: 0.85;
  transition: all 0.25s ease;
}
.evo-node--active {
  opacity: 1;
  transform: translateY(-2px);
}
.evo-icon {
  font-size: 20px;
  margin-bottom: 4px;
}
.evo-node--active .evo-icon {
  filter: drop-shadow(0 0 6px rgba(245, 158, 11, 0.6));
}
.evo-node--locked { opacity: 0.45; }
.evo-name {
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
  max-width: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.evo-node--active .evo-name {
  color: var(--md-gold, #f59e0b);
}
.evo-lv { font-size: 10px; color: var(--color-text-secondary); }
.evo-connector {
  position: absolute;
  top: 12px;
  right: -4px;
  width: 8px;
  height: 2px;
  background: var(--tint-3);
}

/* 当前形态 */
.current-area {
  display: flex;
  gap: 18px;
  margin-bottom: 18px;
}
.current-scene {
  width: 170px;
  min-width: 170px;
  height: 170px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  border: 1px solid var(--tint-2);
}
.current-sprite { width: 120px; height: 120px; }
.current-badge {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  white-space: nowrap;
}

.current-detail {
  flex: 1;
  min-width: 0;
}
.current-name {
  font-size: 17px;
  font-weight: 700;
  margin-bottom: 2px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.current-stage {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-primary);
  background: rgba(79, 70, 229, 0.08);
  padding: 2px 10px;
  border-radius: 12px;
}
.current-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 6px 0 12px;
}

.info-block {
  padding: 10px 14px;
  border-radius: 12px;
  margin-bottom: 10px;
}
.quote-block {
  background: rgba(245, 158, 11, 0.06);
  border-left: 3px solid #f59e0b;
}
.poem-block {
  background: rgba(139, 92, 246, 0.06);
  border-left: 3px solid #8b5cf6;
}
.info-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}
.quote-text { font-size: 14px; font-style: italic; color: var(--color-warning-text, #d97706); }
.poem-text {
  font-size: 13px;
  color: var(--color-primary);
  line-height: 1.8;
}
.poem-line { display: block; }
.poem-line--second { padding-left: 1.4em; text-indent: -1.4em; }

.info-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}
.info-pill {
  font-size: 12px;
  font-weight: 600;
  padding: 5px 12px;
  border-radius: 14px;
  background: var(--tint-1);
  border: 1px solid var(--tint-2);
  color: var(--color-text);
}
.pill--ability { color: var(--color-primary); border-color: rgba(79, 70, 229, 0.2); }

.progress-block {
  background: var(--tint-1);
  border-radius: 12px;
  padding: 10px 14px;
}
.progress-head {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
}
.progress-head strong { color: var(--md-gold, #f59e0b); }
.progress-bar {
  height: 6px;
  background: var(--tint-3);
  border-radius: 3px;
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  border-radius: 3px;
  background: linear-gradient(90deg, var(--color-primary), var(--color-secondary));
  transition: width 0.5s ease;
}

/* 底部按钮 */
.detail-footer {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}
.foot-btn {
  padding: 9px 20px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.15s ease;
  font-family: inherit;
}
.foot-switch {
  background: rgba(167, 139, 250, 0.1);
  border-color: rgba(167, 139, 250, 0.25);
  color: var(--color-primary);
}
.foot-switch:hover { background: rgba(167, 139, 250, 0.18); }
.foot-close {
  background: transparent;
  border-color: var(--tint-3);
  color: var(--color-text-secondary);
}
.foot-close:hover { background: var(--tint-1); }

@media (max-width: 640px) {
  .current-area { flex-direction: column; }
  .current-scene { width: 100%; height: 160px; }
}
</style>
