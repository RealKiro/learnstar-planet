<script setup lang="ts">
import { computed, useId } from 'vue'
import { computePetArt, CONSTELLATION_LINES } from '@/utils/petArtParts'
import type { BodyType, PetArtRender } from '@/utils/petArtParts'

const props = withDefaults(defineProps<{
  speciesId: string
  level: number
  mood?: number
  animate?: boolean
}>(), {
  mood: 60,
  animate: false,
})

// 每实例唯一后缀，保证多实例 SVG 渐变 def 不冲突
const uid = useId()

const art = computed(() => computePetArt(props.speciesId, props.level))

const happy = computed(() => (props.mood ?? 60) >= 60)
const sad = computed(() => (props.mood ?? 60) < 35)

/** 主体缩放(绕地面中心) */
const bodyTransform = computed(
  () => `translate(100 130) scale(${art.value.scale}) translate(-100 -130)`
)

/** 星座连线坐标 */
const constellationPts = computed(() => CONSTELLATION_LINES[props.speciesId] || [])
const constellationPoly = computed(() => {
  if (!constellationPts.value.length) return ''
  const pts = [...constellationPts.value, constellationPts.value[0]]
  return pts.map(([x, y], i) => `${i === 0 ? 'M' : 'L'}${x + 38} ${y + 34}`).join(' ')
})

// ============================================================
// 表情层 / 光环层(SVG 字符串,经 v-html 注入,避免子组件与运行时编译依赖)
// ============================================================

const FACE_CENTER: Record<BodyType, [number, number]> = {
  quadruped: [100, 80],
  dragon: [100, 78],
  bird: [100, 76],
  food: [100, 112],
  humanoid: [100, 86],
  spirit: [100, 104],
  aquatic: [100, 122],
  mecha: [100, 92],
}

const HALO_CENTER: Partial<Record<BodyType, [number, number]>> = {
  quadruped: [100, 34],
  dragon: [100, 30],
  bird: [100, 32],
  food: [100, 38],
  humanoid: [100, 36],
  spirit: [100, 44],
  aquatic: [100, 34],
  mecha: [100, 34],
}

function buildFace(
  cx: number, cy: number, dark: string, accent: string,
  eyes: PetArtRender['eyes'], sad: boolean, happy: boolean, blush: boolean
): string {
  const p: string[] = []
  if (eyes === 'serene') {
    p.push(`<path d="M${cx - 7} ${cy - 2} Q${cx} ${cy - 6} ${cx + 7} ${cy - 2}" stroke="${dark}" stroke-width="2.4" fill="none" stroke-linecap="round" class="blink-eye"/>`)
    p.push(`<path d="M${cx + 7} ${cy - 2} Q${cx + 14} ${cy - 6} ${cx + 21} ${cy - 2}" stroke="${dark}" stroke-width="2.4" fill="none" stroke-linecap="round" class="blink-eye"/>`)
  } else if (sad) {
    p.push(`<path d="M${cx - 8} ${cy - 1} Q${cx - 4} ${cy + 4} ${cx} ${cy + 1}" stroke="${dark}" stroke-width="2.2" fill="none" stroke-linecap="round"/>`)
    p.push(`<path d="M${cx + 8} ${cy - 1} Q${cx + 12} ${cy + 4} ${cx + 16} ${cy + 1}" stroke="${dark}" stroke-width="2.2" fill="none" stroke-linecap="round"/>`)
  } else {
    const r = eyes === 'sharp' ? 4.6 : 6
    for (const dx of [-8, 8]) {
      p.push(`<g class="eye"><circle cx="${cx + dx}" cy="${cy}" r="${r}" fill="${dark}"/><circle cx="${cx + dx - 1.8}" cy="${cy - 1.8}" r="2" fill="#fff"/></g>`)
    }
  }
  if (blush && !sad) {
    p.push(`<ellipse cx="${cx - 13}" cy="${cy + 8}" rx="4.5" ry="3" fill="${accent}" opacity="0.5"/>`)
    p.push(`<ellipse cx="${cx + 13}" cy="${cy + 8}" rx="4.5" ry="3" fill="${accent}" opacity="0.5"/>`)
  }
  if (happy) {
    p.push(`<path d="M${cx - 7} ${cy + 7} Q${cx} ${cy + 14} ${cx + 7} ${cy + 7}" stroke="${dark}" stroke-width="2" fill="none" stroke-linecap="round"/>`)
  } else if (sad) {
    p.push(`<path d="M${cx - 7} ${cy + 10} Q${cx} ${cy + 6} ${cx + 7} ${cy + 10}" stroke="${dark}" stroke-width="2" fill="none" stroke-linecap="round"/>`)
  } else {
    p.push(`<path d="M${cx - 6} ${cy + 7} Q${cx} ${cy + 11} ${cx + 6} ${cy + 7}" stroke="${dark}" stroke-width="2" fill="none" stroke-linecap="round"/>`)
  }
  return p.join('')
}

function buildHalo(cx: number, cy: number, accent: string): string {
  return `<g class="halo"><ellipse cx="${cx}" cy="${cy}" rx="12" ry="4.4" fill="${accent}" opacity="0.9" class="halo-ring"/><ellipse cx="${cx}" cy="${cy}" rx="12" ry="4.4" fill="none" stroke="${accent}" stroke-width="1" opacity="0.6"/></g>`
}

const faceSvg = computed(() => {
  const [cx, cy] = FACE_CENTER[art.value.body]
  return buildFace(cx, cy, art.value.dark, art.value.accent, art.value.eyes, sad.value, happy.value, art.value.blush)
})

const haloSvg = computed(() => {
  const c = HALO_CENTER[art.value.body]
  return c ? buildHalo(c[0], c[1], art.value.accent) : ''
})

// ============================================================
// 系列专属氛围特效
// ============================================================

interface SeriesFx {
  shape: 'spark' | 'star' | 'dash' | 'leaf' | 'fire' | 'dust' | 'sword'
  color: string
}

const SERIES_FX: Record<string, SeriesFx> = {
  myth: { shape: 'fire', color: '#F59E0B' },
  pokemon: { shape: 'spark', color: '#FACC15' },
  national: { shape: 'leaf', color: '#10B981' },
  mecha: { shape: 'dash', color: '#38BDF8' },
  magic: { shape: 'star', color: '#C084FC' },
  prehistoric: { shape: 'fire', color: '#D97706' },
  constellation: { shape: 'dust', color: '#818CF8' },
  folklore: { shape: 'spark', color: '#F97316' },
  festival: { shape: 'spark', color: '#FDE047' },
  qixia: { shape: 'sword', color: '#F59E0B' },
}

const fx = computed(() => SERIES_FX[art.value.seriesId] || null)

/** 氛围粒子位置（成长阶段起出现） */
const fxItems = computed(() => {
  if (!fx.value) return []
  const items: Array<{ id: number; x: number; y: number; delay: number; size: number }> = []
  for (let i = 0; i < 6; i++) {
    items.push({
      id: i,
      x: 24 + ((i * 47) % 152),
      y: 16 + ((i * 31) % 118),
      delay: i * 0.45,
      size: 3 + (i % 3),
    })
  }
  return items
})

/** 生成单个粒子的 SVG 片段（v-html 注入） */
function fxShape(size: number): string {
  const s = fx.value?.shape || 'spark'
  const c = fx.value?.color || '#fff'
  const r = size
  switch (s) {
    case 'star':
      return `<path d="M0 ${-r} L${r * 0.4} ${-r * 0.3} L${r} 0 L${r * 0.4} ${r * 0.3} L0 ${r} L${-r * 0.4} ${r * 0.3} L${-r} 0 L${-r * 0.4} ${-r * 0.3} Z" fill="${c}"/>`
    case 'leaf':
      return `<path d="M0 0 Q${r} ${-r} 0 ${-r * 2.4} Q${-r} ${-r} 0 0 Z" fill="${c}" opacity="0.8"/>`
    case 'dash':
      return `<line x1="${-r * 1.6}" y1="0" x2="${r * 1.6}" y2="0" stroke="${c}" stroke-width="1.8" stroke-linecap="round"/>`
    case 'fire':
      return `<path d="M0 ${r * 0.7} Q${r * 0.7} 0 0 ${-r} Q${-r * 0.7} 0 0 ${r * 0.7} Z" fill="${c}" opacity="0.8"/>`
    case 'dust':
      return `<circle r="${r * 0.7}" fill="${c}" opacity="0.7"/>`
    case 'sword':
      return `<line x1="0" y1="${r}" x2="0" y2="${-r}" stroke="#FDE68A" stroke-width="2.2" stroke-linecap="round"/>`
    default:
      return `<circle r="${r * 0.8}" fill="${c}" opacity="0.8"/>`
  }
}
</script>

<template>
  <svg
    viewBox="0 0 200 200"
    class="pet-sprite"
    :class="{ 'pet-sprite--animated': animate }"
    role="img"
    :aria-label="`宠物 Lv.${level}`"
  >
    <defs>
      <radialGradient :id="`body-grad-${uid}`" cx="42%" cy="34%" r="80%">
        <stop offset="0%" :stop-color="art.light" />
        <stop offset="55%" :stop-color="art.main" />
        <stop offset="100%" :stop-color="art.dark" />
      </radialGradient>
      <radialGradient :id="`egg-grad-${uid}`" cx="40%" cy="30%" r="90%">
        <stop offset="0%" :stop-color="art.light" />
        <stop offset="60%" :stop-color="art.main" />
        <stop offset="100%" :stop-color="art.dark" />
      </radialGradient>
      <linearGradient :id="`glow-${uid}`" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" :stop-color="art.main" stop-opacity="0.7" />
        <stop offset="100%" :stop-color="art.main" stop-opacity="0" />
      </linearGradient>
      <filter :id="`soft-${uid}`" x="-40%" y="-40%" width="180%" height="180%">
        <feGaussianBlur stdDeviation="2.4" />
      </filter>
    </defs>

    <!-- ===== 地面阴影 ===== -->
    <ellipse cx="100" cy="158" rx="32" ry="6" fill="#000" opacity="0.12" />

    <!-- ===== 传说级光环 ===== -->
    <g v-if="art.glow > 0.7" class="aura-wrap">
      <circle cx="100" cy="70" r="58" fill="none" :stroke="art.main" stroke-width="1.6" stroke-dasharray="3 9" opacity="0.65" class="aura-ring aura-ring--a" />
      <circle cx="100" cy="70" r="66" fill="none" :stroke="art.accent" stroke-width="1" stroke-dasharray="2 14" opacity="0.4" class="aura-ring aura-ring--b" />
      <template v-if="art.sparkle">
        <circle v-for="n in 5" :key="n" :cx="46 + n * 27" :cy="30 + (n % 3) * 8" r="2" :fill="art.accent" class="aura-sparkle" :style="{ animationDelay: (n * 0.35) + 's' }" />
      </template>
    </g>

    <!-- ===== 系列专属氛围特效(成长期起) ===== -->
    <g v-if="fx && art.glow > 0.2" class="series-fx">
      <g v-for="p in fxItems" :key="p.id" :transform="`translate(${p.x} ${p.y})`">
        <g class="fx-item" :style="{ animationDelay: p.delay + 's' }" v-html="fxShape(p.size)" />
      </g>
    </g>

    <!-- ===== 主体 ===== -->
    <g :transform="bodyTransform">

      <!-- ---------- 卵生阶段 ---------- -->
      <g v-if="art.stage === 'egg'" class="egg-group">
        <circle cx="100" cy="104" r="46" :fill="`url(#glow-${uid})`" class="egg-glow" />
        <ellipse cx="100" cy="104" rx="31" ry="39" :fill="`url(#egg-grad-${uid})`" class="egg-body" />
        <path d="M100 66 Q106 74 100 82" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.35" />
        <path d="M116 72 Q122 80 116 88" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.3" />
        <path d="M84 74 Q78 82 84 90" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.3" />
        <template v-if="level >= 2">
          <path d="M100 66 L94 88 L102 100 L92 112" :stroke="art.light" stroke-width="1.6" fill="none" class="egg-crack" />
          <circle cx="104" cy="98" r="6" :fill="art.accent" opacity="0.5" :filter="`url(#soft-${uid})`" />
        </template>
      </g>

      <!-- ---------- 四足兽 / 熊 / 小猫狐虎 ---------- -->
      <g v-else-if="art.body === 'quadruped'" class="cute-group">
        <path v-if="art.parts.tail" d="M126 120 Q148 108 152 92 Q156 82 146 84" :fill="art.dark" class="tail" />
        <!-- 七侠佩剑(Lv.6 起) -->
        <g v-if="art.seriesId === 'qixia' && level >= 6" class="sword-aura">
          <line x1="158" y1="58" x2="158" y2="124" stroke="#FDE68A" stroke-width="3" stroke-linecap="round" />
          <path d="M158 124 L162 132 L154 132 Z" fill="#FDE68A" />
          <line x1="158" y1="62" x2="158" y2="118" stroke="#FFF7D6" stroke-width="1.2" opacity="0.7" />
        </g>
        <g v-if="art.parts.wings" class="wings">
          <path d="M68 108 Q44 96 40 74 Q52 88 66 92 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.4" />
          <path d="M132 108 Q156 96 160 74 Q148 88 134 92 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.4" />
        </g>
        <ellipse cx="100" cy="120" rx="34" ry="27" :fill="`url(#body-grad-${uid})`" class="body-main" />
        <ellipse cx="100" cy="128" rx="22" ry="15" :fill="art.light" opacity="0.85" />
        <ellipse cx="80" cy="144" rx="7" ry="5" :fill="art.dark" />
        <ellipse cx="120" cy="144" rx="7" ry="5" :fill="art.dark" />
        <circle cx="100" cy="78" r="31" :fill="`url(#body-grad-${uid})`" />
        <g class="ears">
          <circle cx="74" cy="56" r="11" :fill="art.main" />
          <circle cx="74" cy="56" r="6" :fill="art.accent" opacity="0.8" />
          <circle cx="126" cy="56" r="11" :fill="art.main" />
          <circle cx="126" cy="56" r="6" :fill="art.accent" opacity="0.8" />
        </g>
        <g v-if="art.parts.horns" class="horns">
          <path d="M86 52 Q82 34 74 30 Q82 38 90 50 Z" :fill="art.accent" />
          <path d="M114 52 Q118 34 126 30 Q118 38 110 50 Z" :fill="art.accent" />
        </g>
        <g v-if="art.parts.crown" class="crown">
          <path d="M88 44 L88 30 L94 38 L100 26 L106 38 L112 30 L112 44 Z" fill="#F59E0B" stroke="#B45309" stroke-width="1" />
        </g>
        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
        <!-- 行为动画：张嘴吼 / 爪子扒拉·舔毛 -->
        <ellipse cx="100" cy="95" rx="7" ry="4.5" :fill="art.dark" class="roar-mouth" />
        <ellipse cx="127" cy="133" rx="6" ry="4" :fill="art.dark" class="paw-move" />
      </g>

      <!-- ---------- 龙 / 长形 ---------- -->
      <g v-else-if="art.body === 'dragon'" class="dragon-group">
        <path v-if="art.parts.tail" d="M128 122 Q158 108 166 92 Q172 80 162 82" :fill="art.dark" class="tail" />
        <g v-if="art.parts.wings" class="wings">
          <path d="M62 112 Q34 92 30 66 Q46 84 62 92 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.4" />
          <path d="M138 112 Q166 92 170 66 Q154 84 138 92 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.4" />
        </g>
        <ellipse cx="100" cy="124" rx="44" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main" />
        <ellipse cx="100" cy="132" rx="30" ry="13" :fill="art.light" opacity="0.8" />
        <g v-if="art.parts.fin">
          <path d="M70 102 L84 88 L92 102 Z" :fill="art.accent" />
          <path d="M110 102 L118 90 L126 102 Z" :fill="art.accent" />
        </g>
        <ellipse cx="100" cy="78" rx="33" ry="28" :fill="`url(#body-grad-${uid})`" />
        <ellipse cx="100" cy="92" rx="15" ry="10" :fill="art.main" />
        <circle cx="94" cy="94" r="2.4" :fill="art.dark" />
        <circle cx="106" cy="94" r="2.4" :fill="art.dark" />
        <g v-if="art.parts.horns" class="horns">
          <path d="M84 56 Q76 34 64 28 Q76 38 88 54 Z" :fill="art.accent" />
          <path d="M116 56 Q124 34 136 28 Q124 38 112 54 Z" :fill="art.accent" />
        </g>
        <g v-if="art.parts.crown && level >= 7" class="flame">
          <path d="M100 58 Q94 44 100 32 Q106 44 100 58 Z" fill="#FBBF24" opacity="0.9" />
          <path d="M100 56 Q97 48 100 40 Q103 48 100 56 Z" fill="#FDE68A" />
        </g>
        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
        <!-- 行为动画：张嘴大吼 -->
        <ellipse cx="100" cy="97" rx="9" ry="6" :fill="art.dark" class="roar-mouth" />
      </g>

      <!-- ---------- 鸟类 ---------- -->
      <g v-else-if="art.body === 'bird'" class="bird-group">
        <g v-if="art.parts.tail" class="tail">
          <path d="M128 116 L150 124 L136 130 Z" :fill="art.dark" />
          <path d="M128 118 L152 134 L132 134 Z" :fill="art.accent" />
        </g>
        <g v-if="art.parts.wings" class="wings">
          <path d="M70 108 Q44 96 40 72 Q54 88 70 92 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.3" />
          <path d="M130 108 Q156 96 160 72 Q146 88 130 92 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.3" />
        </g>
        <ellipse cx="100" cy="118" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main" />
        <ellipse cx="100" cy="126" rx="19" ry="13" :fill="art.light" opacity="0.85" />
        <circle cx="100" cy="78" r="27" :fill="`url(#body-grad-${uid})`" />
        <g v-if="art.parts.crown" class="crown">
          <path d="M94 54 Q90 40 96 34 Q100 42 102 34 Q108 40 104 54 Z" :fill="art.accent" />
        </g>
        <path d="M94 88 Q100 100 106 88 Q100 94 94 88 Z" fill="#F59E0B" />
        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
        <!-- 行为动画：张嘴鸣叫 -->
        <ellipse cx="100" cy="95" rx="6" ry="4" :fill="art.dark" class="roar-mouth" />
      </g>

      <!-- ---------- 水生 ---------- -->
      <g v-else-if="art.body === 'aquatic'" class="aquatic-group">
        <g v-if="art.parts.tail" class="tail">
          <path d="M66 120 Q44 108 40 126 Q56 122 68 128 Z" :fill="art.accent" />
        </g>
        <g v-if="art.parts.fin" class="fin-beat">
          <path d="M88 96 L100 82 L112 96 Z" :fill="art.accent" />
        </g>
        <path d="M86 132 Q72 140 84 142 Z" :fill="art.dark" opacity="0.7" />
        <ellipse cx="100" cy="122" rx="36" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main" />
        <ellipse cx="106" cy="128" rx="20" ry="12" :fill="art.light" opacity="0.85" />
        <g transform="translate(116 112)">
          <circle :r="sad ? 4.5 : 5.5" :fill="art.dark" />
          <circle cx="-1.6" cy="-1.6" r="1.9" fill="#fff" />
        </g>
        <path v-if="happy" d="M120 122 Q124 125 128 122" :stroke="art.dark" stroke-width="1.6" fill="none" stroke-linecap="round" />
        <!-- 行为动画：张嘴吐泡 -->
        <ellipse cx="114" cy="126" rx="5" ry="3.5" :fill="art.dark" class="roar-mouth" />
        <g v-html="haloSvg" v-if="art.parts.halo" />
      </g>

      <!-- ---------- 机甲 ---------- -->
      <g v-else-if="art.body === 'mecha'" class="mecha-group">
        <g v-if="art.parts.wings" class="wings">
          <path d="M66 112 L40 88 L64 104 Z" :fill="art.accent" opacity="0.7" />
          <path d="M134 112 L160 88 L136 104 Z" :fill="art.accent" opacity="0.7" />
        </g>
        <rect x="72" y="112" width="56" height="40" rx="10" :fill="art.main" :stroke="art.dark" stroke-width="2" />
        <rect x="80" y="122" width="40" height="22" rx="6" :fill="art.dark" opacity="0.5" />
        <rect x="60" y="108" width="18" height="12" rx="4" :fill="art.dark" />
        <rect x="122" y="108" width="18" height="12" rx="4" :fill="art.dark" />
        <rect x="84" y="78" width="32" height="28" rx="8" :fill="art.main" :stroke="art.dark" stroke-width="2" />
        <rect x="90" y="88" width="20" height="8" rx="4" :fill="art.accent" class="visor" />
        <g v-if="art.parts.horns">
          <line x1="94" y1="78" x2="90" y2="62" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round" />
          <circle cx="90" cy="60" r="3" :fill="art.accent" />
          <line x1="106" y1="78" x2="110" y2="62" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round" />
          <circle cx="110" cy="60" r="3" :fill="art.accent" />
        </g>
        <circle cx="88" cy="128" r="3.4" :fill="art.accent" class="joint" />
        <circle cx="112" cy="128" r="3.4" :fill="art.accent" class="joint" />
        <g v-html="haloSvg" v-if="art.parts.halo" />
      </g>

      <!-- ---------- 灵体 / 星灵 / 星座 ---------- -->
      <g v-else-if="art.body === 'spirit'" class="spirit-group">
        <circle cx="100" cy="104" r="40" :fill="`url(#glow-${uid})`" class="spirit-glow" />
        <circle cx="100" cy="104" r="30" :fill="art.main" opacity="0.9" :stroke="art.accent" stroke-width="2" class="spirit-core" />
        <path v-if="constellationPoly" :d="constellationPoly" :stroke="art.accent" stroke-width="1.6" fill="none" stroke-linecap="round" class="constellation" />
        <circle v-for="(p, i) in constellationPts" :key="i" :cx="p[0] + 38" :cy="p[1] + 34" r="3.2" :fill="art.light" :stroke="art.accent" stroke-width="1" />
        <circle cx="94" cy="104" r="5" :fill="art.dark" />
        <circle cx="94" cy="103" r="1.7" fill="#fff" />
        <circle cx="106" cy="104" r="5" :fill="art.dark" />
        <circle cx="106" cy="103" r="1.7" fill="#fff" />
        <path v-if="happy" d="M96 112 Q100 116 104 112" :stroke="art.dark" stroke-width="1.6" fill="none" stroke-linecap="round" />
        <path v-if="sad" d="M96 116 Q100 112 104 116" :stroke="art.dark" stroke-width="1.6" fill="none" stroke-linecap="round" />
        <!-- 神圣衣(星座 Lv.12) -->
        <g v-if="art.seriesId === 'constellation' && level >= 12" class="sacred-wings">
          <path d="M64 100 Q36 86 32 60 Q48 78 66 84 Z" fill="#FDE68A" opacity="0.9" />
          <path d="M136 100 Q164 86 168 60 Q152 78 134 84 Z" fill="#FDE68A" opacity="0.9" />
        </g>
        <!-- 黄金圣衣(星座 Lv.8 起) -->
        <g v-if="art.seriesId === 'constellation' && art.parts.crown" class="constellation-armor">
          <path d="M60 92 Q52 78 64 72 Q74 80 72 92 Z" fill="#F59E0B" stroke="#B45309" stroke-width="1.2" />
          <path d="M140 92 Q148 78 136 72 Q126 80 128 92 Z" fill="#F59E0B" stroke="#B45309" stroke-width="1.2" />
          <path d="M82 122 L92 108 L100 120 L108 108 L118 122 L106 134 L94 134 Z" fill="#F59E0B" stroke="#B45309" stroke-width="1.2" opacity="0.92" />
          <path d="M90 108 L100 120 L110 108" fill="none" stroke="#FDE68A" stroke-width="1" opacity="0.8" />
        </g>
        <g v-html="haloSvg" v-if="art.parts.halo" />
      </g>

      <!-- ---------- 物化萌物(灯笼/饺子等) ---------- -->
      <g v-else-if="art.body === 'food'" class="food-group">
        <template v-if="speciesId === 'lantern'">
          <line x1="100" y1="142" x2="100" y2="156" :stroke="art.accent" stroke-width="2" />
          <circle cx="100" cy="158" r="3" :fill="art.accent" />
        </template>
        <template v-else-if="speciesId === 'dumpling'">
          <path d="M88 104 Q100 116 112 104" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.4" />
        </template>
        <ellipse cx="100" cy="112" rx="34" ry="30" :fill="`url(#body-grad-${uid})`" :stroke="art.dark" stroke-width="1.6" />
        <ellipse cx="100" cy="120" rx="22" ry="16" :fill="art.light" opacity="0.8" />
        <!-- 节日美食专属细节 -->
        <template v-if="speciesId === 'zongzi'">
          <path d="M76 100 L100 72 L124 100 L100 144 Z" fill="#16A34A" opacity="0.45" />
          <line x1="100" y1="72" x2="100" y2="144" stroke="#15803D" stroke-width="1.6" opacity="0.6" />
          <path d="M76 100 L124 100" stroke="#15803D" stroke-width="1.4" opacity="0.5" />
        </template>
        <template v-else-if="speciesId === 'tanghulu'">
          <line x1="100" y1="48" x2="100" y2="94" stroke="#92400E" stroke-width="2.6" />
          <circle cx="100" cy="58" r="8.5" fill="#DC2626" />
          <circle cx="100" cy="76" r="8.5" fill="#DC2626" />
          <circle cx="100" cy="94" r="8.5" fill="#DC2626" />
          <circle cx="97.5" cy="55" r="2.4" fill="#FECACA" opacity="0.9" />
        </template>
        <template v-else-if="speciesId === 'mooncake'">
          <path d="M78 96 Q100 82 122 96 Q122 128 100 136 Q78 128 78 96 Z" fill="none" :stroke="art.dark" stroke-width="1.2" opacity="0.45" />
          <circle cx="100" cy="112" r="7" fill="none" :stroke="art.dark" stroke-width="1.2" opacity="0.5" />
          <path d="M100 92 L104 104 L116 104 L106 112 L110 124 L100 116 L90 124 L94 112 L84 104 L96 104 Z" :fill="art.dark" opacity="0.4" />
        </template>
        <template v-else-if="speciesId === 'tangyuan'">
          <path d="M90 74 Q94 68 98 74 Q102 68 106 74 Q110 68 114 74" stroke="#E2E8F0" stroke-width="1.6" fill="none" stroke-linecap="round" opacity="0.85" />
        </template>
        <template v-else-if="speciesId === 'niangao'">
          <line x1="78" y1="106" x2="122" y2="106" stroke="#CBD5E1" stroke-width="1.4" opacity="0.6" />
          <line x1="80" y1="118" x2="120" y2="118" stroke="#CBD5E1" stroke-width="1.4" opacity="0.5" />
          <line x1="80" y1="130" x2="120" y2="130" stroke="#CBD5E1" stroke-width="1.2" opacity="0.4" />
        </template>
        <template v-else-if="speciesId === 'qingtuan'">
          <path d="M100 76 Q90 66 84 72 Q90 80 100 76 Z" fill="#15803D" opacity="0.95" />
          <path d="M84 72 Q80 78 86 82" fill="none" stroke="#166534" stroke-width="1.2" opacity="0.6" />
        </template>
        <template v-else-if="speciesId === 'wonton'">
          <path d="M88 100 Q100 112 112 100" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.4" />
          <path d="M92 118 Q100 126 108 118" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.35" />
        </template>
        <template v-else-if="speciesId === 'osmanthus_cake'">
          <circle v-for="n in 4" :key="n" :cx="84 + n * 10" cy="90" r="2.6" fill="#FDE047" />
        </template>
        <g v-if="art.parts.crown" class="crown">
          <circle cx="100" cy="80" r="6" :fill="art.accent" />
        </g>
        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
      </g>

      <!-- ---------- 人形(精灵/仙人等) ---------- -->
      <g v-else class="humanoid-group">
        <g v-if="art.parts.wings" class="wings">
          <path d="M70 96 Q50 84 48 66 Q60 78 72 82 Z" :fill="art.light" opacity="0.85" :stroke="art.dark" stroke-width="1.2" />
          <path d="M130 96 Q150 84 152 66 Q140 78 128 82 Z" :fill="art.light" opacity="0.85" :stroke="art.dark" stroke-width="1.2" />
        </g>
        <path d="M78 112 Q72 150 84 158 L116 158 Q128 150 122 112 Z" :fill="art.main" :stroke="art.dark" stroke-width="1.6" />
        <path d="M78 132 Q100 142 122 132" :fill="art.accent" opacity="0.7" />
        <ellipse cx="74" cy="120" rx="7" ry="12" :fill="art.main" :stroke="art.dark" stroke-width="1.2" />
        <ellipse cx="126" cy="120" rx="7" ry="12" :fill="art.main" :stroke="art.dark" stroke-width="1.2" />
        <circle cx="100" cy="86" r="27" :fill="`url(#body-grad-${uid})`" />
        <g v-if="art.parts.horns" class="horns">
          <path d="M86 62 Q80 46 72 42 Q82 50 90 60 Z" :fill="art.accent" />
          <path d="M114 62 Q120 46 128 42 Q118 50 110 60 Z" :fill="art.accent" />
        </g>
        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
      </g>
    </g>
  </svg>
</template>

<style scoped>
.pet-sprite {
  width: 100%;
  height: 100%;
  display: block;
  overflow: visible;
}
.pet-sprite--animated .egg-body { animation: eggPulse 2.2s ease-in-out infinite; }
.pet-sprite--animated .egg-glow { animation: eggGlowPulse 2.2s ease-in-out infinite; }
.pet-sprite--animated .egg-crack { animation: crackGlow 1.6s ease-in-out infinite; }
.pet-sprite--animated .cute-group, .pet-sprite--animated .dragon-group,
.pet-sprite--animated .bird-group, .pet-sprite--animated .food-group,
.pet-sprite--animated .humanoid-group { animation: petBob 2.4s ease-in-out infinite; }
.pet-sprite--animated .aquatic-group { animation: swimBob 3s ease-in-out infinite; }
.pet-sprite--animated .spirit-group { animation: spiritFloat 3.2s ease-in-out infinite; }
.pet-sprite--animated .mecha-group { animation: mechaHover 3s ease-in-out infinite; }
.pet-sprite--animated .tail { animation: tailWag 2s ease-in-out infinite; transform-origin: 128px 120px; }
.pet-sprite--animated .wings { animation: wingFlap 2.6s ease-in-out infinite; transform-origin: 100px 90px; }
.pet-sprite--animated .eyes, .pet-sprite--animated .blink-eye { animation: blink 4.4s ease-in-out infinite; }
.pet-sprite--animated .visor { animation: visorGlow 1.8s ease-in-out infinite; }
.pet-sprite--animated .joint { animation: jointPulse 1.4s ease-in-out infinite; }
.pet-sprite--animated .aura-ring--a { animation: spinA 8s linear infinite; }
.pet-sprite--animated .aura-ring--b { animation: spinB 10s linear infinite; }
.pet-sprite--animated .aura-sparkle { animation: sparkleFloat 2s ease-in-out infinite; }
.pet-sprite--animated .constellation { animation: dashFlow 6s linear infinite; }
.pet-sprite--animated .spirit-glow { animation: spiritPulse 2.4s ease-in-out infinite; }
.pet-sprite--animated .halo-ring { animation: haloFloat 2.4s ease-in-out infinite; }
.pet-sprite--animated .flame { animation: flameFlicker 1s ease-in-out infinite; }
.pet-sprite--animated .fx-item { animation: fxFloat 3s ease-in-out infinite; opacity: 0; }
.pet-sprite--animated .sword-aura { animation: swordGlow 1.6s ease-in-out infinite; }
.pet-sprite--animated .constellation-armor { animation: armorGlow 2s ease-in-out infinite; }
.pet-sprite--animated .sacred-wings { animation: sacredWings 2.6s ease-in-out infinite; transform-origin: 100px 84px; }
/* ===== 行为动画：耳朵煽动 / 张嘴吼 / 爪子扒拉·舔毛 / 身体呼吸 ===== */
.pet-sprite--animated .ears { animation: earFlick 5s ease-in-out infinite; transform-origin: 50% 30%; }
.pet-sprite--animated .roar-mouth { animation: mouthRoar 6s ease-in-out infinite; transform-origin: center; }
.pet-sprite--animated .paw-move { animation: pawScratch 6.5s ease-in-out infinite; transform-origin: 127px 133px; }
.pet-sprite--animated .body-main { animation: breathe 3.2s ease-in-out infinite; transform-origin: 100px 120px; }
.pet-sprite--animated .fin-beat { animation: finBeat 2.4s ease-in-out infinite; transform-origin: 100px 90px; }

@keyframes petBob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
@keyframes swimBob {
  0%, 100% { transform: translateY(0) rotate(-2deg); }
  50% { transform: translateY(-3px) rotate(2deg); }
}
@keyframes spiritFloat {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-5px) scale(1.03); }
}
@keyframes mechaHover {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}
@keyframes tailWag {
  0%, 100% { transform: rotate(0deg); }
  50% { transform: rotate(8deg); }
}
@keyframes wingFlap {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.06, 0.96); }
}
@keyframes blink {
  0%, 90%, 100% { opacity: 1; transform: scaleY(1); }
  95% { opacity: 0.3; transform: scaleY(0.15); }
}
@keyframes eggPulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.02); }
}
@keyframes eggGlowPulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.9; transform: scale(1.08); }
}
@keyframes crackGlow {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}
@keyframes haloFloat {
  0%, 100% { transform: translateY(0) scale(1); opacity: 0.9; }
  50% { transform: translateY(-2px) scale(1.08); opacity: 0.7; }
}
@keyframes visorGlow {
  0%, 100% { opacity: 0.75; }
  50% { opacity: 1; }
}
@keyframes jointPulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}
@keyframes spinA { to { transform: rotate(360deg); } }
@keyframes spinB { to { transform: rotate(-360deg); } }
@keyframes sparkleFloat {
  0%, 100% { opacity: 0; transform: translateY(0); }
  50% { opacity: 1; transform: translateY(-6px); }
}
@keyframes dashFlow { to { stroke-dashoffset: -40; } }
@keyframes spiritPulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.1); opacity: 0.8; }
}
@keyframes flameFlicker {
  0%, 100% { transform: scale(1) translateY(0); opacity: 0.9; }
  50% { transform: scale(1.08) translateY(-2px); opacity: 1; }
}
@keyframes fxFloat {
  0%, 100% { opacity: 0; transform: translateY(0); }
  50% { opacity: 1; transform: translateY(-9px); }
}
@keyframes swordGlow {
  0%, 100% { opacity: 0.65; }
  50% { opacity: 1; }
}
@keyframes armorGlow {
  0%, 100% { filter: brightness(1); }
  50% { filter: brightness(1.3); }
}
@keyframes sacredWings {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.06); }
}
@keyframes earFlick {
  0%, 78%, 100% { transform: scale(1) rotate(0deg); }
  84% { transform: scale(1.14) rotate(6deg); }
  90% { transform: scale(1.1) rotate(-4deg); }
}
@keyframes mouthRoar {
  0%, 72%, 100% { opacity: 0; transform: scale(0.4); }
  80% { opacity: 1; transform: scale(1.15); }
  88% { opacity: 0.85; transform: scale(1); }
}
@keyframes pawScratch {
  0%, 80%, 100% { transform: translate(0, 0) rotate(0deg); }
  86% { transform: translate(-6px, -12px) rotate(-18deg); }
  93% { transform: translate(-2px, -5px) rotate(-8deg); }
}
@keyframes breathe {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.04, 1.09); }
}
@keyframes finBeat {
  0%, 100% { transform: rotate(0deg); }
  50% { transform: rotate(8deg); }
}
</style>
