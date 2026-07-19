<script setup lang="ts">
import { computed } from 'vue'
import {
  getSeriesBySpeciesId,
  SERIES_SCENES,
  getSpeciesEmoji,
  getPetLevelName,
  getPetLevelDescription,
  getLevelStage,
} from '@/utils/petData'

const props = defineProps<{
  speciesId: string
  level: number
  name?: string
  mood?: number
  animate?: boolean
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>()

const emit = defineEmits<{
  click: []
  feed: []
}>()

const scene = computed(() => {
  const series = getSeriesBySpeciesId(props.speciesId)
  return series ? SERIES_SCENES[series.id] : SERIES_SCENES.myth
})

const stageName = computed(() => getPetLevelName(props.speciesId, props.level))
const stageDesc = computed(() => getPetLevelDescription(props.speciesId, props.level))
const stage = computed(() => getLevelStage(props.level))
const emoji = computed(() => getSpeciesEmoji(props.speciesId))

const moodEmoji = computed(() => {
  if (!props.mood || props.mood < 20) return '😢'
  if (props.mood < 40) return '😐'
  if (props.mood < 60) return '🙂'
  if (props.mood < 80) return '😊'
  return '🥰'
})

const glintActive = computed(() => props.animate && props.level >= 5)
const glowActive = computed(() => props.animate && props.level >= 8)
const auraActive = computed(() => props.animate && props.level >= 11)

const sizeClass = computed(() => `pet-size--${props.size || 'md'}`)

// 根据等级决定宠物显示大小
const petScale = computed(() => {
  const base = props.size === 'sm' ? 1 : props.size === 'lg' ? 1.5 : props.size === 'xl' ? 2 : 1.2
  // 越高越大
  const levelBoost = 1 + (props.level - 1) * 0.03
  return base * Math.min(levelBoost, 1.4)
})

// 根据系列和等级决定发光颜色
const glowColor = computed(() => {
  const colors: Record<string, string> = {
    myth: '#F59E0B',
    pokemon: '#EF4444',
    national: '#10B981',
    mecha: '#3B82F6',
    magic: '#8B5CF6',
    prehistoric: '#D97706',
    constellation: '#6366F1',
    folklore: '#F97316',
  }
  const series = getSeriesBySpeciesId(props.speciesId)
  return colors[series?.id || ''] || '#F59E0B'
})
</script>

<template>
  <div
    class="pet-display"
    :class="[sizeClass, { 'pet-display--animated': animate }]"
    @click="emit('click')"
  >
    <!-- 场景背景 -->
    <div class="pet-scene" :style="{ background: scene.bgGradient }">
      <!-- 装饰元素 -->
      <div class="scene-decor">
        <span v-for="(d, i) in scene.decor" :key="i" class="decor-item" :style="{
          animationDelay: i * 0.4 + 's',
          top: (15 + i * 18) + '%',
          left: (5 + i * 22) + '%',
        }">{{ d }}</span>
      </div>
      <!-- 场景名称 -->
      <div class="scene-label">{{ scene.scene }}</div>
    </div>

    <!-- 光环地板 -->
    <div
      class="pet-halo"
      :class="{
        'halo--glint': glintActive,
        'halo--glow': glowActive,
        'halo--aura': auraActive,
      }"
      :style="{
        background: `radial-gradient(ellipse, ${glowColor}33 0%, transparent 70%)`,
        boxShadow: glowActive ? `0 0 30px ${glowColor}44` : 'none',
      }"
    ></div>

    <!-- 宠物主体 -->
    <div class="pet-body" :style="{ transform: `scale(${petScale})` }">
      <div
        class="pet-emoji-main"
        :class="{
          'pet-emoji--glint': glintActive,
          'pet-emoji--glow': glowActive,
          'pet-emoji--aura': auraActive,
        }"
        :style="{
          filter: glowActive ? `drop-shadow(0 0 12px ${glowColor})` : 'none',
        }"
      >
        {{ emoji }}
      </div>

      <!-- 传说级特效环绕 -->
      <div v-if="auraActive" class="aura-ring">
        <svg viewBox="0 0 120 120" class="aura-svg">
          <circle cx="60" cy="60" r="50" fill="none" :stroke="glowColor" stroke-width="1.5"
            stroke-dasharray="4 8" opacity="0.6">
            <animateTransform attributeName="transform" type="rotate" from="0 60 60" to="360 60 60"
              dur="6s" repeatCount="indefinite" />
          </circle>
          <circle cx="60" cy="60" r="55" fill="none" :stroke="glowColor" stroke-width="0.8"
            stroke-dasharray="2 12" opacity="0.4">
            <animateTransform attributeName="transform" type="rotate" from="360 60 60" to="0 60 60"
              dur="8s" repeatCount="indefinite" />
          </circle>
        </svg>
      </div>

      <!-- Lv.5 以上小特效 -->
      <div v-if="glintActive && !auraActive" class="sparkles">
        <span v-for="n in 4" :key="n" class="sparkle" :style="{
          animationDelay: n * 0.3 + 's',
          top: (10 + n * 20) + '%',
          left: n % 2 === 0 ? '5%' : '90%',
        }">✨</span>
      </div>
    </div>

    <!-- 信息覆盖层 -->
    <div class="pet-info-overlay">
      <div class="pet-level-badge" :class="'level--' + stage">
        Lv.{{ level }}
      </div>

      <div class="pet-mood" v-if="mood !== undefined">
        {{ moodEmoji }}
      </div>
    </div>

    <!-- 名称和阶段 -->
    <div class="pet-footer">
      <div class="pet-name">{{ name || stageName }}</div>
      <div class="pet-stage">{{ stageName }}</div>
      <div v-if="stageDesc" class="pet-desc">{{ stageDesc }}</div>
    </div>

    <!-- 互动按钮 -->
    <button class="pet-feed-btn" @click.stop="emit('feed')" title="投喂宠物">
      🍬
    </button>
  </div>
</template>

<style scoped>
.pet-display {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 24px;
  overflow: hidden;
  cursor: pointer;
  user-select: none;
  transition: transform 0.3s var(--ease-smooth), box-shadow 0.3s var(--ease-smooth);
  min-height: 280px;
}
.pet-display:hover {
  transform: translateY(-4px);
}
.pet-display:active {
  transform: scale(0.98);
}

.pet-size--sm { min-height: 160px; border-radius: 16px; }
.pet-size--lg { min-height: 360px; border-radius: 28px; }
.pet-size--xl { min-height: 460px; border-radius: 32px; }

.pet-scene {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
}
.scene-decor {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.decor-item {
  position: absolute;
  font-size: 20px;
  opacity: 0.3;
  animation: floatDecor 4s ease-in-out infinite;
}
@keyframes floatDecor {
  0%, 100% { transform: translateY(0) rotate(0deg); opacity: 0.3; }
  50% { transform: translateY(-8px) rotate(5deg); opacity: 0.5; }
}
.scene-label {
  position: absolute;
  bottom: 8px;
  right: 12px;
  font-size: 10px;
  color: rgba(255,255,255,0.2);
  letter-spacing: 0.05em;
}

.pet-halo {
  position: absolute;
  bottom: 15%;
  left: 50%;
  transform: translateX(-50%);
  width: 60%;
  aspect-ratio: 2 / 1;
  border-radius: 50%;
  z-index: 1;
  transition: all 0.5s ease;
}
.halo--glint { animation: haloPulse 2s ease-in-out infinite; }
.halo--glow { animation: haloPulse 1.5s ease-in-out infinite; }
.halo--aura { animation: haloPulse 1s ease-in-out infinite; }
@keyframes haloPulse {
  0%, 100% { transform: translateX(-50%) scale(1); opacity: 0.6; }
  50% { transform: translateX(-50%) scale(1.15); opacity: 0.9; }
}

.pet-body {
  position: relative;
  z-index: 2;
  transition: transform 0.5s var(--ease-bounce);
  display: flex;
  align-items: center;
  justify-content: center;
}
.pet-emoji-main {
  font-size: 72px;
  line-height: 1;
  transition: all 0.3s ease;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.2));
}
.pet-size--sm .pet-emoji-main { font-size: 48px; }
.pet-size--lg .pet-emoji-main { font-size: 96px; }
.pet-size--xl .pet-emoji-main { font-size: 120px; }

.pet-emoji--glint {
  animation: petGlint 2s ease-in-out infinite;
}
.pet-emoji--glow {
  animation: petGlow 2s ease-in-out infinite;
}
.pet-emoji--aura {
  animation: petAura 3s ease-in-out infinite;
}
@keyframes petGlint {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}
@keyframes petGlow {
  0%, 100% { transform: scale(1); filter: drop-shadow(0 0 8px currentColor); }
  50% { transform: scale(1.08); filter: drop-shadow(0 0 20px currentColor); }
}
@keyframes petAura {
  0%, 100% { transform: scale(1) rotate(0deg); }
  33% { transform: scale(1.1) rotate(-2deg); }
  66% { transform: scale(1.05) rotate(2deg); }
}

/* 光环SVG */
.aura-ring {
  position: absolute;
  inset: -30px;
  pointer-events: none;
}
.aura-svg {
  width: 100%;
  height: 100%;
}

/* 火花 */
.sparkles {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.sparkle {
  position: absolute;
  font-size: 12px;
  animation: sparkleFloat 2s ease-in-out infinite;
}
@keyframes sparkleFloat {
  0%, 100% { opacity: 0; transform: translateY(0) scale(0.5); }
  50% { opacity: 1; transform: translateY(-10px) scale(1); }
}

.pet-info-overlay {
  position: absolute;
  top: 12px;
  left: 12px;
  right: 12px;
  display: flex;
  justify-content: space-between;
  z-index: 3;
  pointer-events: none;
}
.pet-level-badge {
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  backdrop-filter: blur(8px);
  background: rgba(0,0,0,0.3);
  color: white;
  border: 1px solid rgba(255,255,255,0.15);
}
.level--egg { background: rgba(245,158,11,0.4); }
.level--baby { background: rgba(16,185,129,0.4); }
.level--growing { background: rgba(59,130,246,0.4); }
.level--mature { background: rgba(139,92,246,0.4); }
.level--legendary {
  background: linear-gradient(135deg, rgba(245,158,11,0.4), rgba(239,68,68,0.4));
  animation: legendaryBadge 2s ease-in-out infinite;
}
@keyframes legendaryBadge {
  0%, 100% { box-shadow: 0 0 8px rgba(245,158,11,0.3); }
  50% { box-shadow: 0 0 16px rgba(245,158,11,0.6); }
}

.pet-mood {
  font-size: 18px;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
}

.pet-footer {
  position: absolute;
  bottom: 12px;
  left: 12px;
  right: 12px;
  text-align: center;
  z-index: 3;
  pointer-events: none;
}
.pet-name {
  font-size: 13px;
  font-weight: 600;
  color: rgba(255,255,255,0.9);
  text-shadow: 0 2px 8px rgba(0,0,0,0.4);
}
.pet-stage {
  font-size: 11px;
  color: rgba(255,255,255,0.6);
  text-shadow: 0 1px 4px rgba(0,0,0,0.3);
}
.pet-desc {
  font-size: 10px;
  color: rgba(255,255,255,0.4);
  margin-top: 2px;
  max-width: 80%;
  margin-left: auto;
  margin-right: auto;
  text-shadow: 0 1px 4px rgba(0,0,0,0.3);
}

.pet-feed-btn {
  position: absolute;
  bottom: 60px;
  right: 12px;
  z-index: 4;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: rgba(255,255,255,0.15);
  backdrop-filter: blur(8px);
  font-size: 18px;
  cursor: pointer;
  transition: all 0.2s var(--ease-bounce);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transform: scale(0.8);
}
.pet-display:hover .pet-feed-btn {
  opacity: 1;
  transform: scale(1);
}
.pet-feed-btn:hover {
  background: rgba(255,255,255,0.25);
  transform: scale(1.15);
}
</style>
