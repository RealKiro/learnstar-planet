<script setup lang="ts">
import { computed } from 'vue'
import { getSpeciesById, getLevelRequiredScore, getLevelStage } from '@/utils/petData'

const props = defineProps<{
  speciesId: string
  currentLevel: number
  currentExp: number
}>()

const emit = defineEmits<{
  select: [level: number]
}>()

const species = computed(() => getSpeciesById(props.speciesId))
const levels = computed(() => species.value?.levels || [])

const stageLabels: Record<string, string> = {
  egg: '破卵',
  baby: '幼年',
  growing: '成长期',
  mature: '成熟期',
  legendary: '传说级',
}

function isUnlocked(level: number): boolean {
  return level <= props.currentLevel
}

function isCurrent(level: number): boolean {
  return level === props.currentLevel
}

function getProgress(level: number): number {
  if (level < props.currentLevel) return 100
  if (level === props.currentLevel) {
    const currentScore = getLevelRequiredScore(props.currentLevel)
    const nextScore = getLevelRequiredScore(props.currentLevel + 1)
    const progress = ((props.currentExp) / (nextScore - currentScore)) * 100
    return Math.min(progress, 99)
  }
  return 0
}

function getRemainingScore(level: number): number {
  if (level <= props.currentLevel) return 0
  return Math.max(0, getLevelRequiredScore(level) - getLevelRequiredScore(props.currentLevel) - props.currentExp)
}
</script>

<template>
  <div class="evolution-tree">
    <div class="tree-header">
      <h3 class="tree-title">🌟 进化之路</h3>
      <span class="tree-hint">点击预览形态</span>
    </div>

    <!-- 星轨进度条 -->
    <div class="star-rail">
      <div class="rail-track">
        <div
          class="rail-fill"
          :style="{ width: ((props.currentLevel - 1) / 11) * 100 + '%' }"
        ></div>
      </div>
      <div class="rail-dots">
        <div
          v-for="lvl in 12"
          :key="lvl"
          class="rail-dot"
          :class="{
            'dot--unlocked': isUnlocked(lvl),
            'dot--current': isCurrent(lvl),
            'dot--locked': !isUnlocked(lvl),
          }"
          :style="{ left: ((lvl - 1) / 11) * 100 + '%' }"
          @click="emit('select', lvl)"
        >
          <div class="dot-inner">
            <template v-if="isUnlocked(lvl)">⭐</template>
            <template v-else>🔒</template>
          </div>
          <div class="dot-level">Lv.{{ lvl }}</div>
        </div>
      </div>
    </div>

    <!-- 等级列表 -->
    <div class="level-list">
      <div
        v-for="(lvl, i) in levels"
        :key="lvl.level"
        class="level-item"
        :class="{
          'item--unlocked': isUnlocked(lvl.level),
          'item--current': isCurrent(lvl.level),
          'item--locked': !isUnlocked(lvl.level),
        }"
        @click="emit('select', lvl.level)"
      >
        <!-- 阶段分组标签 -->
        <div
          v-if="i === 0 || levels[i - 1]?.stage !== lvl.stage"
          class="stage-divider"
        >
          <span class="stage-tag">{{ stageLabels[lvl.stage] }}</span>
        </div>

        <div class="item-row">
          <div class="item-level">
            <span class="level-num">Lv.{{ lvl.level }}</span>
          </div>

          <div class="item-info">
            <div class="item-name">
              {{ lvl.name }}
              <span v-if="isCurrent(lvl.level)" class="current-badge">当前</span>
            </div>
            <div class="item-desc">{{ lvl.description }}</div>
          </div>

          <div class="item-status">
            <template v-if="isUnlocked(lvl.level)">
              <span class="status-check">✅</span>
            </template>
            <template v-else>
              <div class="status-locked">
                <div class="require-score">
                  还差 <strong>{{ getRemainingScore(lvl.level) }}</strong> 分
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- 当前等级进度条 -->
        <div v-if="isCurrent(lvl.level)" class="current-exp-bar">
          <div class="exp-track">
            <div class="exp-fill" :style="{ width: getProgress(lvl.level) + '%' }"></div>
          </div>
          <div class="exp-text">
            {{ props.currentExp }} / {{ getLevelRequiredScore(lvl.level + 1) - getLevelRequiredScore(lvl.level) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.evolution-tree {
  background: rgba(255,255,255,0.04);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,0.06);
}

.tree-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.tree-title {
  font-size: 15px;
  font-weight: 700;
  margin: 0;
}
.tree-hint {
  font-size: 11px;
  opacity: 0.5;
}

/* 星轨 */
.star-rail {
  position: relative;
  padding: 16px 0 24px;
  margin-bottom: 20px;
}
.rail-track {
  width: 100%;
  height: 4px;
  background: rgba(255,255,255,0.08);
  border-radius: 2px;
  overflow: hidden;
}
.rail-fill {
  height: 100%;
  background: linear-gradient(90deg, #F59E0B, #EF4444, #8B5CF6);
  border-radius: 2px;
  transition: width 0.6s ease;
}
.rail-dots {
  position: absolute;
  inset: 0;
}
.rail-dot {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: transform 0.2s ease;
}
.rail-dot:hover {
  transform: translate(-50%, -50%) scale(1.15);
}
.dot-inner {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  background: rgba(255,255,255,0.06);
  border: 2px solid rgba(255,255,255,0.1);
  transition: all 0.3s ease;
}
.dot--unlocked .dot-inner {
  background: rgba(245,158,11,0.2);
  border-color: rgba(245,158,11,0.4);
}
.dot--current .dot-inner {
  background: linear-gradient(135deg, rgba(245,158,11,0.3), rgba(239,68,68,0.3));
  border-color: #F59E0B;
  box-shadow: 0 0 16px rgba(245,158,11,0.3);
  animation: dotPulse 2s ease-in-out infinite;
}
@keyframes dotPulse {
  0%, 100% { box-shadow: 0 0 8px rgba(245,158,11,0.2); }
  50% { box-shadow: 0 0 20px rgba(245,158,11,0.5); }
}
.dot-level {
  font-size: 9px;
  font-weight: 600;
  color: rgba(255,255,255,0.4);
}

/* 等级列表 */
.level-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 320px;
  overflow-y: auto;
}
.level-list::-webkit-scrollbar { width: 3px; }
.level-list::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 2px; }

.stage-divider {
  padding: 8px 0 4px;
  margin-top: 4px;
}
.stage-tag {
  font-size: 10px;
  font-weight: 600;
  color: rgba(255,255,255,0.3);
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.level-item {
  padding: 8px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}
.level-item:hover {
  background: rgba(255,255,255,0.04);
}
.item--current {
  background: rgba(245,158,11,0.08);
  border-color: rgba(245,158,11,0.2);
}
.item--locked {
  opacity: 0.5;
}
.item--locked:hover {
  opacity: 0.8;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.item-level {
  width: 36px;
  flex-shrink: 0;
}
.level-num {
  font-size: 11px;
  font-weight: 700;
  color: rgba(255,255,255,0.5);
}
.item--current .level-num {
  color: #F59E0B;
}

.item-info {
  flex: 1;
  min-width: 0;
}
.item-name {
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}
.current-badge {
  font-size: 9px;
  padding: 1px 6px;
  border-radius: 4px;
  background: linear-gradient(135deg, #F59E0B, #EF4444);
  color: white;
  font-weight: 700;
}
.item-desc {
  font-size: 10px;
  color: rgba(255,255,255,0.4);
  margin-top: 1px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-status {
  flex-shrink: 0;
  text-align: right;
}
.status-check {
  font-size: 14px;
}
.require-score {
  font-size: 10px;
  color: rgba(255,255,255,0.4);
}
.require-score strong {
  color: #F59E0B;
}

/* 当前经验条 */
.current-exp-bar {
  margin-top: 6px;
  padding-left: 48px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.exp-track {
  flex: 1;
  height: 4px;
  background: rgba(255,255,255,0.06);
  border-radius: 2px;
  overflow: hidden;
}
.exp-fill {
  height: 100%;
  background: linear-gradient(90deg, #F59E0B, #EF4444);
  border-radius: 2px;
  transition: width 0.5s ease;
}
.exp-text {
  font-size: 9px;
  color: rgba(255,255,255,0.3);
  white-space: nowrap;
}
</style>
