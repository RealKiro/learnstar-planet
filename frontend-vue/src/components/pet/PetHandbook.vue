<script setup lang="ts">
import { ref, computed } from 'vue'
import { getSpeciesById, getLevelRequiredScore, getSeriesBySpeciesId, getSpeciesEmoji, SERIES_SCENES } from '@/utils/petData'

const props = defineProps<{
  speciesId: string
  currentLevel: number
  currentScore: number
}>()

const emit = defineEmits<{
  close: []
}>()

const selectedLevel = ref<number | null>(null)

const species = computed(() => getSpeciesById(props.speciesId))
const series = computed(() => getSeriesBySpeciesId(props.speciesId))
const levels = computed(() => species.value?.levels || [])
const scene = computed(() => series.value ? SERIES_SCENES[series.value.id] : null)

const previewLevel = computed(() => selectedLevel.value || props.currentLevel)

// 获取关键里程碑等级展示（1, 3, 5, 8, 12）
const milestoneLevels = computed(() => {
  const milestones = [1, 3, 5, 8, 12]
  return milestones.map(l => levels.value.find(lvl => lvl.level === l)).filter(Boolean)
})

function isUnlocked(level: number): boolean {
  return level <= props.currentLevel
}

function getRequiredForNext(level: number): number {
  return getLevelRequiredScore(level)
}
</script>

<template>
  <div class="handbook-overlay" @click.self="emit('close')">
    <div class="handbook-modal">
      <!-- 头部 -->
      <div class="handbook-header">
        <div class="header-info">
          <span class="header-emoji">{{ getSpeciesEmoji(speciesId) }}</span>
          <div>
            <h2 class="header-title">{{ species?.name }} · 进化图鉴</h2>
            <p class="header-series">{{ series?.name }}系列</p>
          </div>
        </div>
        <button class="close-btn" @click="emit('close')">✕</button>
      </div>

      <!-- 星轨进化路 -->
      <div class="milestone-track">
        <div
          v-for="(m, i) in milestoneLevels"
          :key="m?.level"
          class="milestone-node"
          :class="{ 'node--active': selectedLevel === m?.level }"
          @click="selectedLevel = m?.level || null"
        >
          <div class="milestone-icon" :class="{ 'icon--locked': !isUnlocked(m?.level || 0) }">
            {{ isUnlocked(m?.level || 0) ? '⭐' : '🔒' }}
          </div>
          <div class="milestone-name">{{ m?.name }}</div>
          <div class="milestone-lv">Lv.{{ m?.level }}</div>
          <!-- 连接线 -->
          <div v-if="i < milestoneLevels.length - 1" class="milestone-connector"></div>
        </div>
      </div>

      <!-- 预览区域 -->
      <div class="preview-area">
        <div class="preview-scene" v-if="scene" :style="{ background: scene.bgGradient }">
          <div class="preview-emoji">{{ getSpeciesEmoji(speciesId) }}</div>
          <div class="preview-level-badge">Lv.{{ previewLevel }}</div>
        </div>

        <div class="preview-detail">
          <div class="detail-name">
            {{ levels.find(l => l.level === previewLevel)?.name }}
            <span v-if="isUnlocked(previewLevel)" class="detail-unlocked">✅ 已解锁</span>
            <span v-else class="detail-locked">🔒 未解锁</span>
          </div>
          <p class="detail-desc">{{ levels.find(l => l.level === previewLevel)?.description }}</p>

          <div class="detail-stats">
            <div class="stat-item">
              <span class="stat-label">所需积分</span>
              <span class="stat-value">{{ getRequiredForNext(previewLevel) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">阶段</span>
              <span class="stat-value stage-badge" :class="'stage--' + (levels.find(l => l.level === previewLevel)?.stage || 'egg')">
                {{ { egg: '新生', baby: '幼年', growing: '成长期', mature: '成熟期', legendary: '传说级' }[(levels.find(l => l.level === previewLevel)?.stage || 'egg')] }}
              </span>
            </div>
          </div>

          <!-- 未解锁时的进度提示 -->
          <div v-if="!isUnlocked(previewLevel)" class="unlock-hint">
            还差 <strong>{{ Math.max(0, getRequiredForNext(previewLevel) - currentScore) }}</strong> 分解锁此形态
          </div>
        </div>
      </div>

      <!-- 全部等级列表 -->
      <div class="all-levels">
        <div class="all-levels-header">全部 12 级形态</div>
        <div class="levels-grid">
          <div
            v-for="lvl in levels"
            :key="lvl.level"
            class="level-card"
            :class="{
              'card--active': selectedLevel === lvl.level,
              'card--unlocked': isUnlocked(lvl.level),
              'card--locked': !isUnlocked(lvl.level),
            }"
            @click="selectedLevel = lvl.level"
          >
            <div class="card-emoji">
              {{ isUnlocked(lvl.level) ? '⭐' : '🔒' }}
            </div>
            <div class="card-name">{{ lvl.name }}</div>
            <div class="card-lv">Lv.{{ lvl.level }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.handbook-overlay {
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

.handbook-modal {
  width: 100%;
  max-width: 680px;
  max-height: 85vh;
  background: linear-gradient(180deg, #1a1040 0%, #0d1b2a 100%);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 24px;
  overflow-y: auto;
  padding: 28px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.5);
}
.handbook-modal::-webkit-scrollbar { width: 4px; }
.handbook-modal::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 2px; }

/* 头部 */
.handbook-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}
.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.header-emoji {
  font-size: 36px;
}
.header-title {
  font-size: 18px;
  font-weight: 700;
  color: white;
  margin: 0;
}
.header-series {
  font-size: 12px;
  color: rgba(255,255,255,0.4);
  margin: 2px 0 0;
}
.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.04);
  color: rgba(255,255,255,0.5);
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.close-btn:hover {
  background: rgba(255,255,255,0.1);
  color: white;
}

/* 里程碑星轨 */
.milestone-track {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px 0;
  margin-bottom: 20px;
  position: relative;
}
.milestone-node {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  position: relative;
  flex: 1;
}
.milestone-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  background: rgba(255,255,255,0.06);
  border: 2px solid rgba(255,255,255,0.1);
  transition: all 0.3s ease;
}
.node--active .milestone-icon {
  border-color: #F59E0B;
  box-shadow: 0 0 16px rgba(245,158,11,0.3);
  transform: scale(1.1);
}
.icon--locked {
  opacity: 0.5;
}
.milestone-name {
  font-size: 10px;
  color: rgba(255,255,255,0.6);
  text-align: center;
  max-width: 64px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.milestone-lv {
  font-size: 9px;
  color: rgba(255,255,255,0.3);
}
.milestone-connector {
  position: absolute;
  top: 22px;
  left: 55%;
  width: 90%;
  height: 2px;
  background: rgba(255,255,255,0.06);
  z-index: -1;
}

/* 预览区 */
.preview-area {
  display: flex;
  gap: 20px;
  padding: 20px;
  background: rgba(255,255,255,0.03);
  border-radius: 16px;
  margin-bottom: 20px;
}
.preview-scene {
  width: 120px;
  height: 120px;
  border-radius: 16px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
.preview-emoji {
  font-size: 48px;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.3));
}
.preview-level-badge {
  position: absolute;
  bottom: 6px;
  right: 6px;
  font-size: 9px;
  padding: 2px 6px;
  border-radius: 4px;
  background: rgba(0,0,0,0.4);
  color: rgba(255,255,255,0.7);
}
.preview-detail {
  flex: 1;
}
.detail-name {
  font-size: 16px;
  font-weight: 700;
  color: white;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.detail-unlocked {
  font-size: 10px;
  color: #10B981;
  font-weight: 500;
}
.detail-locked {
  font-size: 10px;
  color: rgba(255,255,255,0.3);
  font-weight: 500;
}
.detail-desc {
  font-size: 13px;
  color: rgba(255,255,255,0.6);
  line-height: 1.5;
  margin: 0 0 12px;
}
.detail-stats {
  display: flex;
  gap: 16px;
}
.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.stat-label {
  font-size: 10px;
  color: rgba(255,255,255,0.3);
}
.stat-value {
  font-size: 14px;
  font-weight: 700;
  color: white;
}
.stage-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
}
.stage--egg { background: rgba(245,158,11,0.2); color: #F59E0B; }
.stage--baby { background: rgba(16,185,129,0.2); color: #10B981; }
.stage--growing { background: rgba(59,130,246,0.2); color: #3B82F6; }
.stage--mature { background: rgba(139,92,246,0.2); color: #8B5CF6; }
.stage--legendary { background: rgba(245,158,11,0.2); color: #F59E0B; }

.unlock-hint {
  margin-top: 10px;
  padding: 8px 12px;
  background: rgba(245,158,11,0.1);
  border-radius: 8px;
  font-size: 12px;
  color: rgba(255,255,255,0.6);
}
.unlock-hint strong {
  color: #F59E0B;
  font-size: 14px;
}

/* 全部等级 */
.all-levels-header {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255,255,255,0.4);
  margin-bottom: 12px;
}
.levels-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}
.level-card {
  padding: 10px 6px;
  border-radius: 10px;
  text-align: center;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s ease;
  background: rgba(255,255,255,0.02);
}
.level-card:hover {
  background: rgba(255,255,255,0.05);
}
.card--active {
  border-color: rgba(245,158,11,0.3);
  background: rgba(245,158,11,0.06);
}
.card--locked {
  opacity: 0.4;
}
.card--locked:hover {
  opacity: 0.7;
}
.card-emoji {
  font-size: 18px;
  margin-bottom: 4px;
}
.card-name {
  font-size: 10px;
  font-weight: 500;
  color: rgba(255,255,255,0.7);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-lv {
  font-size: 9px;
  color: rgba(255,255,255,0.3);
}
</style>
