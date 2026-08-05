<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { slides, features } from './landingData'

const currentSlide = ref(0)
const direction = ref<'left' | 'right'>('right')
let timer: number = 0

onMounted(() => {
  timer = setInterval(() => { direction.value = 'right'; currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
})
onUnmounted(() => { clearInterval(timer) })

function goToSlide(i: number) {
  direction.value = i > currentSlide.value ? 'right' : 'left'
  currentSlide.value = i
  clearInterval(timer)
  timer = setInterval(() => { direction.value = 'right'; currentSlide.value = (currentSlide.value + 1) % slides.length }, 5000)
}

const featureOffset = ref(0)
let featureTimer: number = 0
const visibleCount = 3

onMounted(() => {
  featureTimer = setInterval(() => { featureOffset.value = (featureOffset.value + 1) % features.length }, 3000)
})
onUnmounted(() => { clearInterval(featureTimer) })

const visibleFeatures = computed(() => {
  const r: typeof features = []
  for (let i = 0; i < Math.min(visibleCount, features.length); i++) r.push(features[(featureOffset.value + i) % features.length])
  return r
})

function goToFeature(offset: number) {
  featureOffset.value = offset
  clearInterval(featureTimer)
  featureTimer = setInterval(() => { featureOffset.value = (featureOffset.value + 1) % features.length }, 3000)
}
</script>

<template>
  <div class="left-glow left-glow-top"></div>
  <div class="left-glow left-glow-bottom"></div>
  <div class="panel">
    <div class="slide-stage">
      <transition :name="'slide-' + direction" mode="out-in">
        <div :key="currentSlide" class="slide-card">
          <div class="badge"><span class="badge-dot"></span>{{ slides[currentSlide].badge }}</div>
          <div class="icon">{{ slides[currentSlide].icon }}</div>
          <h1 class="title">{{ slides[currentSlide].title }}<br><span class="gradient">{{ slides[currentSlide].highlight }}</span></h1>
          <p class="desc">{{ slides[currentSlide].desc }}</p>
        </div>
      </transition>
      <div class="dots">
        <button v-for="(_, idx) in slides" :key="idx" :class="['dot', { active: currentSlide === idx }]" @click="goToSlide(idx)"></button>
      </div>
    </div>
    <div class="feature-strip">
      <div class="strip-label">功能模块</div>
      <div class="strip-track">
        <transition name="strip-fade" mode="out-in">
          <div :key="featureOffset" class="strip-group">
            <button v-for="f in visibleFeatures" :key="f.title" class="strip-item" :title="f.title">
              <span class="strip-emoji">{{ f.icon }}</span><span class="strip-name">{{ f.title }}</span>
            </button>
          </div>
        </transition>
      </div>
      <div class="strip-dots">
        <button v-for="(_, idx) in features" :key="idx" :class="['sd', { active: featureOffset === idx }]" @click="goToFeature(idx)"></button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.left-glow { position: fixed; border-radius: 50%; filter: blur(100px); opacity: .2; pointer-events: none; z-index: 0 }
.left-glow-top { width: 500px; height: 500px; background: #c7d2fe; top: -200px; right: -100px; animation: glow 12s ease-in-out infinite }
.left-glow-bottom { width: 350px; height: 350px; background: #a7f3d0; bottom: -100px; left: -50px; animation: glow 15s ease-in-out infinite reverse }
@keyframes glow { 0%,100% { transform: translate(0,0) } 50% { transform: translate(40px,-30px) } }
.panel { position: relative; z-index: 1; width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 60px 56px }
.slide-stage { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; width: 100%; max-width: 460px }
.slide-card { text-align: left; width: 100% }
.badge { display: inline-flex; align-items: center; gap: 8px; background: #fff; border: 1px solid #e5e5ea; border-radius: 9999px; padding: 6px 16px; font-size: 13px; color: #6e6e73; margin-bottom: 28px }
.badge-dot { width: 6px; height: 6px; background: #34c759; border-radius: 50%; display: inline-block; box-shadow: 0 0 6px rgba(52,199,89,.3) }
.icon { font-size: 52px; margin-bottom: 16px }
.title { font-size: 36px; font-weight: 900; line-height: 1.2; letter-spacing: -1px; color: #1d1d1f; margin-bottom: 16px }
.gradient { background: linear-gradient(135deg,#5e5ce6,#ff375f,#ff9f0a); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text }
.desc { font-size: 15px; color: #86868b; line-height: 1.7; white-space: pre-line }
.slide-right-enter-active,.slide-right-leave-active,.slide-left-enter-active,.slide-left-leave-active { transition: all .45s cubic-bezier(.4,0,.2,1) }
.slide-right-enter-from { opacity: 0; transform: translateX(60px) }
.slide-right-leave-to { opacity: 0; transform: translateX(-60px) }
.slide-left-enter-from { opacity: 0; transform: translateX(-60px) }
.slide-left-leave-to { opacity: 0; transform: translateX(60px) }
.dots { display: flex; gap: 8px; margin-top: 32px }
.dot { height: 8px; border-radius: 4px; border: none; cursor: pointer; background: #d2d2d7; transition: all .3s; width: 8px; padding: 0 }
.dot.active { width: 32px; background: linear-gradient(135deg,#5e5ce6,#818cf8) }
.feature-strip { width: 100%; max-width: 500px; margin-top: 40px; flex-shrink: 0 }
.strip-label { font-size: 12px; font-weight: 600; color: #aeaeb2; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 12px; text-align: center }
.strip-track { overflow: hidden; min-height: 60px; display: flex; align-items: center; justify-content: center }
.strip-group { display: flex; gap: 8px; justify-content: center }
.strip-item { display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 10px 16px; background: #fff; border: 1px solid #f0f0f3; border-radius: 12px; cursor: default; transition: all .2s; min-width: 88px }
.strip-item:hover { border-color: #d2d2d7; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,.04) }
.strip-emoji { font-size: 22px }
.strip-name { font-size: 11px; color: #86868b; font-weight: 500; white-space: nowrap }
.strip-fade-enter-active,.strip-fade-leave-active { transition: all .4s ease }
.strip-fade-enter-from { opacity: 0; transform: translateX(30px) }
.strip-fade-leave-to { opacity: 0; transform: translateX(-30px) }
.strip-dots { display: flex; gap: 5px; justify-content: center; margin-top: 10px }
.sd { width: 5px; height: 5px; border-radius: 50%; border: none; background: #d2d2d7; cursor: pointer; padding: 0; transition: all .3s }
.sd.active { background: #5e5ce6; transform: scale(1.4) }
@media (max-width: 768px) { .panel { padding: 80px 24px 32px } .title { font-size: 28px } .icon { font-size: 36px } }
</style>
