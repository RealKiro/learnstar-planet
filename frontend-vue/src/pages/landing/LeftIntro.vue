<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { slides, features, stages } from './landingData'

const currentSlide = ref(0)
let timer: number = 0

onMounted(function () {
  timer = setInterval(function () {
    currentSlide.value = (currentSlide.value + 1) % slides.length
  }, 5000)
})
onUnmounted(function () { clearInterval(timer) })

function goToSlide(i: number) {
  currentSlide.value = i
  clearInterval(timer)
  timer = setInterval(function () {
    currentSlide.value = (currentSlide.value + 1) % slides.length
  }, 5000)
}
</script>

<template>
  <div class="left-glow left-glow--top"></div>
  <div class="left-glow left-glow--bottom"></div>
  <div class="content">
    <transition name="fade" mode="out-in">
      <div :key="currentSlide" class="slide">
        <div class="badge">
          <span class="badge-dot"></span>
          {{ slides[currentSlide].badge }}
        </div>
        <div class="icon">{{ slides[currentSlide].icon }}</div>
        <h1 class="title">
          {{ slides[currentSlide].title }}
          <br>
          <span class="gradient">{{ slides[currentSlide].highlight }}</span>
        </h1>
        <p class="desc">{{ slides[currentSlide].desc }}</p>
      </div>
    </transition>
    <div class="dots">
      <button
        v-for="(s, idx) in slides"
        :key="idx"
        :class="['dot', { active: currentSlide === idx }]"
        @click="goToSlide(idx)"
      ></button>
    </div>
    <div class="footer-links">
      <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="link">GitHub</a>
      <span class="sep">MIT</span>
    </div>
  </div>
  <div class="features">
    <h2 class="section-title">12 大功能模块</h2>
    <p class="section-sub">覆盖班级管理全场景，全部免费</p>
    <div class="grid">
      <div v-for="f in features" :key="f.title" class="card">
        <span class="card-icon">{{ f.icon }}</span>
        <div>
          <div class="card-name">{{ f.title }}</div>
          <div class="card-desc">{{ f.desc }}</div>
        </div>
      </div>
    </div>
    <h2 class="section-title" style="margin-top:48px">11 阶宠物进化</h2>
    <p class="section-sub">积分变经验，从星尘到银河</p>
    <div class="evo">
      <div
        v-for="(s, idx) in stages"
        :key="s.name"
        style="display:flex;align-items:center;gap:2px"
      >
        <span v-if="idx !== 0" class="evo-arrow">→</span>
        <div class="evo-item">
          <span class="evo-emoji">{{ s.emoji }}</span>
          <span class="evo-name">{{ s.name }}</span>
        </div>
      </div>
    </div>
    <div class="cta-box">
      <div class="cta-icon">🚀</div>
      <h3>加入开源社区</h3>
      <p>项目完全开源，欢迎 Star、Fork、Issue、PR</p>
      <div class="cta-buttons">
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="cta-dark">Star</a>
        <a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" class="cta-ghost">Issue</a>
      </div>
    </div>
    <footer class="foot">
      <a href="https://github.com/RealKiro/learnstar-planet" target="_blank">学趣星球</a>
      <span> · MIT · 自托管 · 完全免费</span>
    </footer>
  </div>
</template>

<style scoped>
.left-glow {
  position: fixed;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.2;
  pointer-events: none;
  z-index: 0;
}
.left-glow--top {
  width: 500px; height: 500px;
  background: #c7d2fe;
  top: -200px; right: -100px;
  animation: glow 12s ease-in-out infinite;
}
.left-glow--bottom {
  width: 350px; height: 350px;
  background: #a7f3d0;
  bottom: -100px; left: -50px;
  animation: glow 15s ease-in-out infinite reverse;
}
@keyframes glow {
  0%, 100% { transform: translate(0, 0) }
  50% { transform: translate(40px, -30px) }
}
.content {
  position: relative; z-index: 1;
  min-height: 100vh;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  padding: 60px 64px 80px;
}
.slide { max-width: 440px; text-align: left }
.badge {
  display: inline-flex; align-items: center; gap: 8px;
  background: #fff; border: 1px solid #e5e5ea;
  border-radius: 9999px; padding: 6px 16px;
  font-size: 13px; color: #6e6e73; margin-bottom: 28px;
}
.badge-dot {
  width: 6px; height: 6px;
  background: #34c759; border-radius: 50%;
  display: inline-block;
  box-shadow: 0 0 6px rgba(52,199,89,.3);
}
.icon { font-size: 56px; margin-bottom: 20px }
.title {
  font-size: 40px; font-weight: 900; line-height: 1.2;
  letter-spacing: -1px; color: #1d1d1f; margin-bottom: 20px;
}
.gradient {
  background: linear-gradient(135deg,#5e5ce6,#ff375f,#ff9f0a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.desc { font-size: 16px; color: #86868b; line-height: 1.7; white-space: pre-line }
.dots { display: flex; gap: 8px; margin-top: 44px }
.dot {
  height: 8px; border-radius: 4px; border: none; cursor: pointer;
  background: #d2d2d7; transition: all .3s; width: 8px;
}
.dot.active { width: 32px; background: linear-gradient(135deg,#5e5ce6,#818cf8) }
.footer-links {
  position: absolute; bottom: 40px; left: 64px;
  display: flex; gap: 16px; align-items: center; z-index: 1;
}
.link { color: #86868b; font-size: 13px; text-decoration: none; transition: color .2s }
.link:hover { color: #1d1d1f }
.sep { color: #aeaeb2; font-size: 13px }
.fade-enter-active { transition: all .4s ease }
.fade-leave-active { transition: all .3s ease }
.fade-enter-from { opacity: 0; transform: translateX(20px) }
.fade-leave-to { opacity: 0; transform: translateX(-20px) }
.features { padding: 0 64px 60px; position: relative; z-index: 1 }
.section-title { font-size: 24px; font-weight: 800; color: #1d1d1f; margin-bottom: 8px }
.section-sub { font-size: 15px; color: #86868b; margin-bottom: 32px }
.grid { display: grid; grid-template-columns: repeat(2,1fr); gap: 10px }
.card {
  display: flex; align-items: flex-start; gap: 12px;
  padding: 16px; background: #fff;
  border: 1px solid #f0f0f3; border-radius: 14px;
  transition: all .2s;
}
.card:hover {
  border-color: #e5e5ea;
  box-shadow: 0 4px 16px rgba(0,0,0,.04);
  transform: translateY(-2px);
}
.card-icon { font-size: 24px; flex-shrink: 0; width: 36px; text-align: center }
.card-name { font-size: 14px; font-weight: 700; color: #1d1d1f; margin-bottom: 3px }
.card-desc { font-size: 12px; color: #86868b; line-height: 1.5 }
.evo { display: flex; align-items: center; flex-wrap: wrap; gap: 2px; padding: 16px 0 }
.evo-arrow { color: #c7c7cc; font-size: 13px; user-select: none }
.evo-item {
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  padding: 8px 12px; border-radius: 14px; transition: all .2s; cursor: default;
}
.evo-item:hover { background: #f0f0f3; transform: scale(1.08) }
.evo-emoji { font-size: 26px }
.evo-name { font-size: 11px; color: #86868b; font-weight: 500 }
.cta-box {
  margin-top: 48px; text-align: center; padding: 40px;
  background: #fff; border: 1px solid #f0f0f3; border-radius: 20px;
}
.cta-icon { font-size: 40px; margin-bottom: 12px }
.cta-box h3 { font-size: 22px; font-weight: 800; color: #1d1d1f; margin-bottom: 8px }
.cta-box p { font-size: 14px; color: #86868b; margin-bottom: 24px }
.cta-buttons { display: flex; gap: 10px; justify-content: center }
.cta-dark, .cta-ghost {
  padding: 10px 24px; border-radius: 9999px;
  font-size: 14px; font-weight: 600;
  text-decoration: none; transition: all .2s;
}
.cta-dark { background: #1d1d1f; color: #fff }
.cta-dark:hover { background: #333 }
.cta-ghost { border: 1px solid #d2d2d7; color: #1d1d1f }
.cta-ghost:hover { background: #f5f5f7 }
.foot {
  margin-top: 40px; padding-bottom: 40px;
  text-align: center; font-size: 13px; color: #aeaeb2;
}
.foot a { color: #86868b; text-decoration: none }
.foot a:hover { color: #1d1d1f }
@media (max-width: 768px) {
  .content { padding: 80px 24px 40px; min-height: auto }
  .footer-links { position: static; margin-top: 32px; padding: 0 }
  .features { padding: 0 24px 40px }
  .grid { grid-template-columns: 1fr }
  .title { font-size: 32px }
  .icon { font-size: 40px }
}
</style>
