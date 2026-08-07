<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { getSeriesBySpeciesId, getSpeciesById } from '@/utils/petData'
import type { PetDetail } from '@/types'

const props = defineProps<{
  pets: PetDetail[]
}>()

const QIXIA_SWORDS = ['hongmao', 'lantu', 'doudou', 'dabeng', 'tiaotiao', 'shali', 'dada']
const CONSTELLATIONS = [
  'aries', 'taurus', 'gemini', 'cancer', 'leo', 'virgo',
  'libra', 'scorpio', 'sagittarius', 'capricorn', 'aquarius', 'pisces',
]

/** 七把名剑的颜色与名称（七侠合璧） */
const swordData = [
  { color: '#EF4444', label: '长虹' },
  { color: '#3B82F6', label: '冰魄' },
  { color: '#E2E8F0', label: '雨花' },
  { color: '#92400E', label: '奔雷' },
  { color: '#22C55E', label: '青光' },
  { color: '#8B5CF6', label: '紫云' },
  { color: '#16A34A', label: '旋风' },
]

/** 七把剑围成一圈的位置 */
const swordPositions = swordData.map((_, i) => {
  const angle = (i / swordData.length) * Math.PI * 2 - Math.PI / 2
  const r = 100
  return {
    x: 200 + Math.cos(angle) * r,
    y: 150 + Math.sin(angle) * r,
    rot: (angle * 180) / Math.PI + 90,
  }
})

interface Achievement {
  key: string
  icon: string
  title: string
  desc: string
  fx: 'qixia' | 'constellation' | 'festival' | 'legend'
}

const achievements = computed<Achievement[]>(() => {
  const list: Achievement[] = []

  // 1) 七剑合璧：集齐 7 位持剑侠客(Lv.6+)
  const heroes = new Set(props.pets.filter(p => QIXIA_SWORDS.includes(p.species) && p.level >= 6).map(p => p.species))
  if (heroes.size >= 7) {
    list.push({ key: 'qixia', icon: '⚔️', title: '七剑合璧', desc: '七侠齐聚 · 长虹贯日 · 全班武侠之力觉醒', fx: 'qixia' })
  }

  // 2) 燃烧小宇宙：≥6 位星座守护者达黄金圣衣(Lv.8+)
  const goldSaints = props.pets.filter(p => CONSTELLATIONS.includes(p.species) && p.level >= 8).length
  if (goldSaints >= 6) {
    list.push({ key: 'constellation', icon: '♈', title: '燃烧的小宇宙', desc: `黄金十二宫守卫觉醒（${goldSaints} 位）· 第七感贯通全班`, fx: 'constellation' })
  }

  // 3) 诗词显化：传统节日宠达传说级
  const poemPet = props.pets.find(p => getSeriesBySpeciesId(p.species)?.id === 'festival' && p.level >= 12)
  if (poemPet) {
    const lv12 = getSpeciesById(poemPet.species)?.levels.find(l => l.level === 12)
    list.push({ key: 'festival', icon: '🏮', title: '诗词显化', desc: lv12?.description || lv12?.name || '诗意正浓', fx: 'festival' })
  }

  // 4) 传说觉醒：任意传说级
  const legends = props.pets.filter(p => p.level >= 12).length
  if (legends > 0) {
    list.push({ key: 'legend', icon: '👑', title: '传说觉醒', desc: `${legends} 只传说级宠物降临教室`, fx: 'legend' })
  }

  return list
})

const activeFx = computed<Achievement['fx']>(() => achievements.value[0]?.fx || 'legend')

const activePoem = computed(() => {
  const a = achievements.value.find(x => x.key === 'festival')
  return a?.desc || ''
})
const poemChars = computed(() => activePoem.value.split(''))

// ===== 全屏动画：每个浏览器会话仅播放一次(避免每次进入宠物花园都重放) =====
const overlayShow = ref(false)

function maybePlay() {
  if (achievements.value.length && !sessionStorage.getItem('pet_fx_played')) {
    sessionStorage.setItem('pet_fx_played', '1')
    overlayShow.value = true
    setTimeout(() => { overlayShow.value = false }, 5200)
  }
}

onMounted(maybePlay)
watch(() => props.pets.map(p => `${p.species}:${p.level}`).join('|'), maybePlay)
</script>

<template>
  <!-- 常驻成就条 -->
  <div v-if="achievements.length" class="achievement-bar">
    <span class="bar-label">🏆 班级成就</span>
    <span v-for="a in achievements" :key="a.key" class="ach-chip" :class="'chip--' + a.fx">
      {{ a.icon }} {{ a.title }}
    </span>
  </div>

  <!-- 全屏庆祝动画 -->
  <Teleport to="body">
    <Transition name="fx-fade">
      <div v-if="overlayShow" class="fx-overlay" @click="overlayShow = false">
        <div class="fx-scene" :class="'fx--' + activeFx">
          <!-- 七剑合璧 -->
          <template v-if="activeFx === 'qixia'">
            <svg viewBox="0 0 400 300" class="fx-svg">
              <g v-for="(sw, i) in swordPositions" :key="i" :transform="`translate(${sw.x} ${sw.y}) rotate(${sw.rot})`">
                <g class="sword-anim" :style="{ animationDelay: i * 0.14 + 's' }">
                  <line x1="0" y1="42" x2="0" y2="-42" :stroke="swordData[i].color" stroke-width="5" stroke-linecap="round" />
                  <path d="M0 -42 L5 -52 L-5 -52 Z" :fill="swordData[i].color" />
                  <text y="-58" text-anchor="middle" font-size="11" :fill="swordData[i].color">{{ swordData[i].label }}</text>
                </g>
              </g>
              <circle cx="200" cy="150" r="0" class="sword-flash" />
              <circle cx="200" cy="150" r="0" class="sword-flash sword-flash--2" />
            </svg>
            <h1 class="fx-title gold">⚔️ 七剑合璧 ⚔️</h1>
            <p class="fx-sub">{{ achievements[0]?.desc }}</p>
          </template>

          <!-- 燃烧小宇宙 -->
          <template v-else-if="activeFx === 'constellation'">
            <svg viewBox="0 0 400 300" class="fx-svg">
              <circle cx="200" cy="150" r="60" fill="none" stroke="#F59E0B" stroke-width="2.5" class="fx-ring fx-ring--a" />
              <circle cx="200" cy="150" r="86" fill="none" stroke="#FDE68A" stroke-width="1.2" stroke-dasharray="3 9" class="fx-ring fx-ring--b" />
              <circle v-for="n in 8" :key="n" :cx="200 + Math.cos((n / 8) * Math.PI * 2) * 60" :cy="150 + Math.sin((n / 8) * Math.PI * 2) * 60" r="3.4" fill="#F59E0B" class="fx-dot" :style="{ animationDelay: n * 0.12 + 's' }" />
            </svg>
            <h1 class="fx-title gold">♈ 燃烧吧，小宇宙！</h1>
            <p class="fx-sub">{{ achievements[0]?.desc }}</p>
          </template>

          <!-- 诗词显化 -->
          <template v-else-if="activeFx === 'festival'">
            <div class="poem-scene">
              <span v-for="(c, i) in poemChars" :key="i" class="poem-char" :style="{ animationDelay: i * 0.08 + 's' }">{{ c }}</span>
            </div>
            <h1 class="fx-title gold">🏮 诗词显化</h1>
            <p class="fx-sub">{{ activePoem }}</p>
          </template>

          <!-- 传说觉醒 -->
          <template v-else>
            <div class="legend-burst">
              <div class="burst-ring burst-ring--1"></div>
              <div class="burst-ring burst-ring--2"></div>
            </div>
            <h1 class="fx-title gold">👑 传说觉醒</h1>
            <p class="fx-sub">{{ achievements[0]?.desc }}</p>
          </template>
        </div>
        <div class="fx-hint">点击任意处继续</div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* ===== 常驻成就条 ===== */
.achievement-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 8px 12px;
  margin-bottom: 14px;
  border-radius: 14px;
  background: linear-gradient(90deg, rgba(245,158,11,0.12), rgba(239,68,68,0.1));
  border: 1px solid rgba(245,158,11,0.25);
}
.bar-label {
  font-size: 12px; font-weight: 700; color: #F59E0B;
}
.ach-chip {
  font-size: 12px; font-weight: 600;
  padding: 3px 10px; border-radius: 20px;
  background: var(--tint-3);
  border: 1px solid var(--tint-4);
  color: var(--color-text);
}
.chip--qixia { border-color: rgba(239,68,68,0.4); color: #FCA5A5; }
.chip--constellation { border-color: rgba(245,158,11,0.45); color: #FCD34D; }
.chip--festival { border-color: rgba(251,191,36,0.4); color: #FDE68A; }
.chip--legend { border-color: rgba(167,139,250,0.4); color: #C4B5FD; }

/* ===== 全屏覆盖 ===== */
.fx-overlay {
  position: fixed; inset: 0; z-index: 600;
  background: radial-gradient(ellipse at center, rgba(10,6,30,0.92) 0%, rgba(5,2,20,0.97) 70%);
  display: flex; align-items: center; justify-content: center;
  flex-direction: column; padding: 20px; overflow: hidden;
}
.fx-scene { text-align: center; position: relative; }
.fx-svg { width: min(420px, 90vw); height: auto; }
.fx-title {
  font-size: 42px; font-weight: 900; margin: 18px 0 8px;
  letter-spacing: 0.12em;
  animation: titlePop 0.6s cubic-bezier(0.34,1.56,0.64,1) 0.9s both;
}
.fx-title.gold {
  background: linear-gradient(135deg, #FDE68A, #F59E0B, #EF4444, #FDE68A);
  background-size: 300% 100%;
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: goldShine 2s linear infinite, titlePop 0.6s cubic-bezier(0.34,1.56,0.64,1) 0.9s both;
}
.fx-sub {
  font-size: 15px; color: rgba(255,255,255,0.7); font-weight: 500;
  opacity: 0; animation: subIn 0.8s ease 1.4s forwards;
}
.fx-hint {
  position: absolute; bottom: 28px; font-size: 12px; color: rgba(255,255,255,0.3);
  animation: hintBlink 2s ease-in-out infinite;
}
@keyframes hintBlink { 0%,100%{opacity:0.3} 50%{opacity:0.7} }

/* ===== 七剑 ===== */
.sword-anim {
  opacity: 0;
  animation: swordIn 0.5s cubic-bezier(0.34,1.56,0.64,1) forwards;
  filter: drop-shadow(0 0 8px currentColor);
}
@keyframes swordIn {
  0% { opacity: 0; transform: scale(0.1) rotate(-20deg); }
  70% { opacity: 1; transform: scale(1.15) rotate(3deg); }
  100% { opacity: 1; transform: scale(1) rotate(0); }
}
.sword-flash {
  animation: flashBurst 1s ease-out 1.05s both;
  fill: #FDE68A;
}
.sword-flash--2 { animation-delay: 1.25s; }
@keyframes flashBurst {
  0% { r: 0; opacity: 1; }
  100% { r: 120; opacity: 0; }
}

/* ===== 小宇宙 ===== */
.fx-ring--a { transform-origin: 200px 150px; animation: spinCW 3s linear infinite; }
.fx-ring--b { transform-origin: 200px 150px; animation: spinCCW 5s linear infinite; }
.fx-dot { opacity: 0; animation: dotPulse 0.8s ease-in-out infinite; }
@keyframes dotPulse { 0%,100%{opacity:0.2} 50%{opacity:1} }
@keyframes spinCW { to { transform: rotate(360deg); } }
@keyframes spinCCW { to { transform: rotate(-360deg); } }

/* ===== 诗词 ===== */
.poem-scene {
  min-height: 90px; max-width: 80vw;
  display: flex; justify-content: center; flex-wrap: wrap; gap: 6px;
  padding: 20px; margin: 8px auto 4px;
}
.poem-char {
  font-size: 30px; font-weight: 700; color: #FDE68A;
  text-shadow: 0 0 18px rgba(251,191,36,0.6);
  opacity: 0; display: inline-block;
  animation: poemUp 0.7s cubic-bezier(0.34,1.56,0.64,1) forwards;
}
@keyframes poemUp {
  0% { opacity: 0; transform: translateY(30px) scale(0.5); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}

/* ===== 传说爆发 ===== */
.legend-burst { position: relative; width: 260px; height: 260px; margin: 0 auto; }
.burst-ring {
  position: absolute; inset: 0; border-radius: 50%;
  border: 2px solid rgba(245,158,11,0.7);
  animation: burstExpand 1.4s ease-out infinite;
}
.burst-ring--2 { animation-delay: 0.5s; border-color: rgba(239,68,68,0.5); }
@keyframes burstExpand {
  0% { transform: scale(0.1); opacity: 1; }
  100% { transform: scale(1.3); opacity: 0; }
}

/* ===== 过渡 ===== */
.fx-fade-enter-active { animation: overlayIn 0.4s ease; }
.fx-fade-leave-active { animation: overlayIn 0.3s ease reverse; }
@keyframes overlayIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes titlePop {
  0% { transform: scale(0.6); opacity: 0; }
  60% { transform: scale(1.1); }
  100% { transform: scale(1); opacity: 1; }
}
@keyframes subIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes goldShine {
  0% { background-position: 0% 0; }
  100% { background-position: 300% 0; }
}
</style>
