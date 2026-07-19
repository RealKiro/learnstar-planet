<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

const props = defineProps<{
  show: boolean
  points: number
  speciesId: string
  level: number
}>()

const emit = defineEmits<{
  close: []
}>()

const particles = ref<Array<{ id: number; x: number; y: number; emoji: string; delay: number }>>([])
let particleId = 0

const isLegendary = computed(() => props.level >= 11)
const isMature = computed(() => props.level >= 8)

const pointEmojis = ['✨', '⭐', '💫', '🌟', '🔥', '💖']

function spawnParticles() {
  const items = []
  for (let i = 0; i < 20; i++) {
    items.push({
      id: particleId++,
      x: Math.random() * 100,
      y: Math.random() * 100,
      emoji: pointEmojis[Math.floor(Math.random() * pointEmojis.length)],
      delay: Math.random() * 0.5,
    })
  }
  particles.value = items
}

let timer: ReturnType<typeof setTimeout> | null = null

onMounted(() => {
  if (props.show) {
    spawnParticles()
    timer = setTimeout(() => emit('close'), 2500)
  }
})

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <Transition name="feed-popup">
    <div v-if="show" class="feed-overlay" @click="emit('close')">
      <!-- 粒子背景 -->
      <div class="feed-particles">
        <span
          v-for="p in particles"
          :key="p.id"
          class="particle"
          :style="{
            left: p.x + '%',
            top: p.y + '%',
            animationDelay: p.delay + 's',
            fontSize: (12 + Math.random() * 16) + 'px',
          }"
        >{{ p.emoji }}</span>
      </div>

      <div class="feed-content" @click.stop>
        <!-- +分数字 -->
        <div class="score-pop">
          <span class="score-plus">+{{ points }}</span>
          <span class="score-label">EXP</span>
        </div>

        <!-- 宠物跳跃动画 -->
        <div class="pet-jump">
          <span class="jump-emoji">⭐</span>
        </div>

        <!-- 完成文字 -->
        <div class="feed-text">
          <div class="feed-title">太棒了！</div>
          <div class="feed-subtitle">
            {{ isLegendary ? '✨ 传说之力正在觉醒 ✨' : isMature ? '🌟 力量源源不断地涌入 🌟' : '💪 继续加油成长吧！' }}
          </div>
        </div>

        <!-- 关闭提示 -->
        <div class="click-hint">点击任意处继续</div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.feed-overlay {
  position: fixed;
  inset: 0;
  z-index: 500;
  background: rgba(5, 2, 20, 0.85);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.feed-particles {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.particle {
  position: absolute;
  animation: particleFloat 1.5s ease-out forwards;
  opacity: 0;
}
@keyframes particleFloat {
  0% { transform: translateY(0) scale(0.5); opacity: 0; }
  20% { opacity: 1; }
  100% { transform: translateY(-60px) scale(1.2); opacity: 0; }
}

.feed-content {
  position: relative;
  z-index: 1;
  text-align: center;
  animation: contentIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes contentIn {
  from { transform: scale(0.7); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

.score-pop {
  margin-bottom: 20px;
}
.score-plus {
  font-size: 64px;
  font-weight: 900;
  background: linear-gradient(135deg, #F59E0B, #EF4444, #8B5CF6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  display: block;
  line-height: 1;
  text-shadow: none;
  animation: scoreBounce 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes scoreBounce {
  0% { transform: scale(0.3); opacity: 0; }
  60% { transform: scale(1.15); }
  100% { transform: scale(1); opacity: 1; }
}
.score-label {
  font-size: 16px;
  color: rgba(255,255,255,0.4);
  font-weight: 600;
  letter-spacing: 0.1em;
  display: block;
  margin-top: 4px;
}

.pet-jump {
  margin-bottom: 20px;
}
.jump-emoji {
  font-size: 72px;
  display: inline-block;
  animation: petJump 0.8s cubic-bezier(0.34, 1.56, 0.64, 1) 2;
}
@keyframes petJump {
  0%, 100% { transform: translateY(0) scale(1); }
  30% { transform: translateY(-30px) scale(1.15); }
  60% { transform: translateY(-10px) scale(1.05); }
}

.feed-text {
  margin-bottom: 24px;
}
.feed-title {
  font-size: 24px;
  font-weight: 700;
  color: white;
  margin-bottom: 8px;
}
.feed-subtitle {
  font-size: 14px;
  color: rgba(255,255,255,0.5);
}

.click-hint {
  font-size: 12px;
  color: rgba(255,255,255,0.25);
  animation: hintBlink 2s ease-in-out infinite;
}
@keyframes hintBlink {
  0%, 100% { opacity: 0.25; }
  50% { opacity: 0.6; }
}

.feed-popup-enter-active { animation: overlayIn 0.3s ease; }
.feed-popup-leave-active { animation: overlayIn 0.2s ease reverse; }
@keyframes overlayIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
