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

/** 剪影变体专属：头/脸/光环定位覆盖（头位偏离默认模板的变体用） */
const VARIANT_FACE: Record<string, [number, number]> = {
  sauropod: [100, 54],
  ankylo: [100, 72],
  trex: [100, 62],
  xuanwu: [100, 72],
  mouse: [100, 84],
  turtle: [100, 84],
  seed: [100, 82],
  lizard: [100, 82],
}
const VARIANT_HALO: Record<string, [number, number]> = {
  sauropod: [100, 22],
  ankylo: [100, 34],
  trex: [100, 26],
  xuanwu: [100, 34],
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
  const c = VARIANT_FACE[art.value.variant] || FACE_CENTER[art.value.body]
  return buildFace(c[0], c[1], art.value.dark, art.value.accent, art.value.eyes, sad.value, happy.value, art.value.blush)
})

const haloSvg = computed(() => {
  const c = VARIANT_HALO[art.value.variant] || HALO_CENTER[art.value.body]
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

      <!-- ---------- 四足兽（含宝可梦/七侠/麒麟等独立剪影变体） ---------- -->
      <g v-else-if="art.body === 'quadruped'" class="cute-group">
        <!-- 七侠佩剑(Lv.6 起) -->
        <g v-if="art.seriesId === 'qixia' && level >= 6" class="sword-aura">
          <line x1="158" y1="58" x2="158" y2="124" stroke="#FDE68A" stroke-width="3" stroke-linecap="round" />
          <path d="M158 124 L162 132 L154 132 Z" fill="#FDE68A" />
          <line x1="158" y1="62" x2="158" y2="118" stroke="#FFF7D6" stroke-width="1.2" opacity="0.7" />
        </g>

        <!-- 猫（虹猫） -->
        <template v-if="art.variant === 'cat'">
          <path d="M126 116 Q150 104 154 88 Q158 78 148 80" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="28" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="18" ry="11" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M78 60 L70 40 L92 56 Z" :fill="art.main"/>
            <path d="M122 60 L130 40 L108 56 Z" :fill="art.main"/>
            <path d="M80 57 L76 47 L88 55 Z" :fill="art.accent" opacity="0.8"/>
            <path d="M120 57 L124 47 L112 55 Z" :fill="art.accent" opacity="0.8"/>
          </g>
          <path d="M84 92 Q66 94 58 88" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.5"/>
          <path d="M116 92 Q134 94 142 88" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.5"/>
        </template>

        <!-- 兔（蓝兔） -->
        <template v-else-if="art.variant === 'rabbit'">
          <circle cx="132" cy="124" r="7" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
          <ellipse cx="100" cy="120" rx="30" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="19" ry="12" :fill="art.light" opacity="0.9"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <ellipse cx="86" cy="40" rx="7" ry="23" :fill="art.main" transform="rotate(-8 86 40)"/>
            <ellipse cx="86" cy="42" rx="3.4" ry="15" :fill="art.accent" opacity="0.7" transform="rotate(-8 86 40)"/>
            <ellipse cx="114" cy="40" rx="7" ry="23" :fill="art.main" transform="rotate(8 114 40)"/>
            <ellipse cx="114" cy="42" rx="3.4" ry="15" :fill="art.accent" opacity="0.7" transform="rotate(8 114 40)"/>
          </g>
        </template>

        <!-- 小狗（逗逗） -->
        <template v-else-if="art.variant === 'puppy'">
          <path d="M126 118 Q146 108 150 92" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="122" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="19" ry="12" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <ellipse cx="76" cy="66" rx="9" ry="13" :fill="art.main" transform="rotate(-20 76 66)"/>
            <ellipse cx="124" cy="66" rx="9" ry="13" :fill="art.main" transform="rotate(20 124 66)"/>
          </g>
          <ellipse cx="100" cy="104" rx="8" ry="5" :fill="art.dark" opacity="0.5"/>
        </template>

        <!-- 小猫（莎莉） -->
        <template v-else-if="art.variant === 'kitten'">
          <path d="M128 118 Q150 108 152 92" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="124" rx="26" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="16" ry="10" :fill="art.light" opacity="0.85"/>
          <circle cx="100" cy="82" r="25" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 64 L74 46 L94 58 Z" :fill="art.main"/>
            <path d="M120 64 L126 46 L106 58 Z" :fill="art.main"/>
          </g>
          <path v-if="level >= 5" d="M96 56 Q100 52 104 56" :stroke="art.accent" stroke-width="2.4" fill="none" stroke-linecap="round"/>
        </template>

        <!-- 熊（大奔） -->
        <template v-else-if="art.variant === 'bear'">
          <ellipse cx="100" cy="120" rx="34" ry="28" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="22" ry="14" :fill="art.light" opacity="0.8"/>
          <ellipse cx="76" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="124" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <circle cx="100" cy="78" r="29" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <circle cx="76" cy="56" r="10" :fill="art.main"/>
            <circle cx="124" cy="56" r="10" :fill="art.main"/>
          </g>
          <ellipse cx="100" cy="102" rx="10" ry="6" :fill="art.main" opacity="0.6"/>
        </template>

        <!-- 猴（跳跳） -->
        <template v-else-if="art.variant === 'monkey'">
          <path d="M126 120 Q150 116 158 104 Q154 126 140 132" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="122" rx="26" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="78" cy="142" rx="9" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="142" rx="9" ry="6" :fill="art.dark"/>
          <circle cx="100" cy="82" r="26" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <circle cx="78" cy="68" r="8" :fill="art.main" :stroke="art.dark" stroke-width="1.2"/>
            <circle cx="122" cy="68" r="8" :fill="art.main" :stroke="art.dark" stroke-width="1.2"/>
          </g>
          <path d="M88 62 Q94 56 100 60 Q106 56 112 62" :stroke="art.dark" stroke-width="1.6" fill="none" opacity="0.6"/>
        </template>

        <!-- 熊猫（达达） -->
        <template v-else-if="art.variant === 'panda'">
          <ellipse cx="100" cy="122" rx="32" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="20" ry="12" :fill="art.light" opacity="0.9"/>
          <ellipse cx="76" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="124" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <circle cx="78" cy="58" r="9" :fill="art.dark"/>
            <circle cx="122" cy="58" r="9" :fill="art.dark"/>
          </g>
          <g v-if="level >= 5" class="panda-patches">
            <ellipse cx="88" cy="82" rx="6" ry="5" :fill="art.dark" opacity="0.85"/>
            <ellipse cx="112" cy="82" rx="6" ry="5" :fill="art.dark" opacity="0.85"/>
          </g>
        </template>

        <!-- 虎（白虎/黑小虎） -->
        <template v-else-if="art.variant === 'tiger'">
          <path d="M126 116 Q150 104 154 88" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="31" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="20" ry="13" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="78" r="28" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 58 L72 40 L94 54 Z" :fill="art.main"/>
            <path d="M120 58 L128 40 L106 54 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 5" class="stripes">
            <path d="M86 60 L80 66 M92 58 L90 66 M118 60 L124 66 M112 58 L114 66" :stroke="art.dark" stroke-width="2.4" stroke-linecap="round" opacity="0.55"/>
            <path d="M88 116 L84 124 M104 114 L104 122 M116 116 L120 124" :stroke="art.dark" stroke-width="2.4" stroke-linecap="round" opacity="0.45"/>
          </g>
        </template>

        <!-- 麒麟（麒麟/麒麟圣） -->
        <template v-else-if="art.variant === 'qilin'">
          <g v-if="level >= 6" class="flame-tail">
            <path d="M126 118 Q152 104 158 86 Q160 76 150 78" :stroke="art.accent" stroke-width="4" fill="none" stroke-linecap="round"/>
            <path d="M152 92 Q160 80 154 70" stroke="#FDE68A" stroke-width="3" fill="none" stroke-linecap="round"/>
          </g>
          <ellipse cx="100" cy="120" rx="30" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="19" ry="12" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="78" r="28" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 5" class="scales">
            <path d="M88 108 Q94 104 100 108 Q106 104 112 108" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.4"/>
            <path d="M88 116 Q94 112 100 116 Q106 112 112 116" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.4"/>
          </g>
          <g class="horns">
            <path d="M86 56 Q80 42 74 38 Q82 46 88 54 Z" :fill="art.accent"/>
            <path d="M114 56 Q120 42 126 38 Q118 46 112 54 Z" :fill="art.accent"/>
          </g>
        </template>

        <!-- 牛（牛旋风） -->
        <template v-else-if="art.variant === 'bull'">
          <path d="M126 120 Q148 112 152 98" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="122" rx="32" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="20" ry="12" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="100" cy="80" rx="30" ry="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <ellipse cx="74" cy="66" rx="7" ry="5" :fill="art.main" transform="rotate(-20 74 66)"/>
            <ellipse cx="126" cy="66" rx="7" ry="5" :fill="art.main" transform="rotate(20 126 66)"/>
          </g>
          <g class="horns">
            <path d="M82 58 Q72 44 66 40 Q76 46 84 56 Z" :fill="art.accent"/>
            <path d="M118 58 Q128 44 134 40 Q124 46 116 56 Z" :fill="art.accent"/>
          </g>
          <ellipse cx="100" cy="100" rx="12" ry="7" :fill="art.main" opacity="0.6"/>
          <circle cx="100" cy="102" r="1.8" :fill="art.dark"/>
        </template>

        <!-- 小火龙（charmander） -->
        <template v-else-if="art.variant === 'lizard'">
          <g v-if="level >= 3" class="flame-tail">
            <path d="M128 116 Q148 108 150 94 Q152 82 142 84" :fill="art.dark"/>
            <path d="M146 96 Q154 86 148 76 Q146 88 144 92 Z" fill="#F97316"/>
            <path d="M146 92 Q150 86 146 82" fill="#FDE68A"/>
          </g>
          <ellipse cx="100" cy="118" rx="26" ry="23" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="16" ry="10" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="82" r="25" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 62 L72 44 L92 56 Z" :fill="art.main"/>
            <path d="M120 62 L128 44 L108 56 Z" :fill="art.main"/>
          </g>
        </template>

        <!-- 妙蛙种子（bulbasaur） -->
        <template v-else-if="art.variant === 'seed'">
          <path v-if="level >= 5" d="M100 74 Q88 66 96 58 Q104 62 100 74 Z" :fill="art.dark" opacity="0.7"/>
          <ellipse cx="100" cy="116" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="129" rx="19" ry="12" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="142" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="142" rx="8" ry="6" :fill="art.dark"/>
          <path d="M100 104 Q64 100 58 88 Q68 76 100 78 Q132 76 142 88 Q136 100 100 104 Z" :fill="`url(#body-grad-${uid})`"/>
          <circle cx="100" cy="90" r="20" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 4" class="bulb">
            <path d="M100 64 Q88 60 88 48 Q92 40 100 40 Q108 40 112 48 Q112 60 100 64 Z" fill="#16A34A"/>
            <path d="M100 44 Q97 50 100 54 Q103 50 100 44 Z" fill="#22C55E"/>
          </g>
        </template>

        <!-- 杰尼龟（squirtle） -->
        <template v-else-if="art.variant === 'turtle'">
          <path d="M72 116 Q60 116 54 126" :stroke="art.main" stroke-width="8" fill="none" stroke-linecap="round" opacity="0.8"/>
          <path d="M128 116 Q140 116 146 126" :stroke="art.main" stroke-width="8" fill="none" stroke-linecap="round" opacity="0.8"/>
          <ellipse cx="100" cy="118" rx="34" ry="22" :fill="art.dark" opacity="0.6"/>
          <path d="M78 112 Q100 96 122 112 L122 124 Q100 136 78 124 Z" :fill="`url(#body-grad-${uid})`" :stroke="art.dark" stroke-width="1.6"/>
          <path d="M88 112 L100 102 L112 112 L100 122 Z" :fill="art.accent" opacity="0.5"/>
          <circle cx="100" cy="86" r="22" :fill="`url(#body-grad-${uid})`"/>
        </template>

        <!-- 伊布（eevee） -->
        <template v-else-if="art.variant === 'fox'">
          <path d="M128 114 Q152 100 154 84" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="118" rx="28" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="130" rx="18" ry="11" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 60 L72 42 L94 55 Z" :fill="art.main"/>
            <path d="M120 60 L128 42 L106 55 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="fluff">
            <path d="M80 96 Q74 90 80 88 M120 96 Q126 90 120 88" :stroke="art.light" stroke-width="4" stroke-linecap="round" opacity="0.9"/>
          </g>
        </template>

        <!-- 皮卡丘（pikachu） -->
        <template v-else-if="art.variant === 'mouse'">
          <path d="M126 118 Q148 108 154 96 L150 106 L146 94 L142 102" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="27" ry="23" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="17" ry="11" :fill="art.light" opacity="0.85"/>
          <circle cx="100" cy="84" r="24" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M82 66 L72 44 L92 58 Z" :fill="art.main"/>
            <path d="M118 66 L128 44 L108 58 Z" :fill="art.main"/>
            <path d="M83 62 L79 52 L89 58 Z" fill="#FCA5A5" opacity="0.8"/>
            <path d="M117 62 L121 52 L111 58 Z" fill="#FCA5A5" opacity="0.8"/>
          </g>
          <g v-if="level >= 4" class="cheeks">
            <circle cx="86" cy="96" r="5" fill="#F87171" opacity="0.8"/>
            <circle cx="114" cy="96" r="5" fill="#F87171" opacity="0.8"/>
          </g>
        </template>

        <!-- 利欧路（riolu） -->
        <template v-else-if="art.variant === 'pup'">
          <path d="M128 118 Q150 108 152 94" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="28" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="18" ry="12" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="26" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M82 62 L76 42 L96 56 Z" :fill="art.main"/>
            <path d="M118 62 L124 42 L104 56 Z" :fill="art.main"/>
          </g>
          <path v-if="level >= 4" d="M96 56 Q100 52 104 56 Q108 52 112 56" :stroke="art.dark" stroke-width="2" fill="none" stroke-linecap="round" opacity="0.6"/>
        </template>

        <!-- 九尾狐·九尾 -->
        <template v-else-if="art.variant === 'kitsune'">
          <g v-if="level >= 4" class="ninetails">
            <path d="M122 118 Q150 104 154 86 Q158 76 148 78" :fill="art.dark" opacity="0.85"/>
            <path d="M118 122 Q148 114 158 98 Q162 86 152 88" :fill="art.main" opacity="0.7"/>
            <path d="M128 120 Q158 110 162 94" :fill="art.dark" opacity="0.5"/>
            <path d="M126 116 Q144 112 150 102" :fill="art.light" opacity="0.7"/>
            <path d="M124 114 Q138 108 142 98" :fill="art.dark" opacity="0.6"/>
          </g>
          <ellipse cx="100" cy="118" rx="26" ry="23" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="130" rx="16" ry="10" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 60 L72 42 L94 55 Z" :fill="art.main"/>
            <path d="M120 60 L128 42 L106 55 Z" :fill="art.main"/>
            <path d="M81 57 L78 49 L90 55 Z" :fill="art.accent" opacity="0.8"/>
            <path d="M119 57 L122 49 L110 55 Z" :fill="art.accent" opacity="0.8"/>
          </g>
        </template>

        <!-- 赛博猫·电路猫 -->
        <template v-else-if="art.variant === 'cyber'">
          <ellipse cx="100" cy="120" rx="28" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="18" ry="11" :fill="art.dark" opacity="0.5"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 60 L72 42 L94 56 Z" :fill="art.main"/>
            <path d="M120 60 L128 42 L106 56 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="circuits">
            <path d="M76 100 L84 100 L84 92 M124 100 L116 100 L116 92" :stroke="art.accent" stroke-width="1.8" fill="none" opacity="0.8"/>
            <circle cx="84" cy="92" r="2.4" :fill="art.accent"/>
            <circle cx="116" cy="92" r="2.4" :fill="art.accent"/>
            <path d="M82 68 L90 68 L90 74" :stroke="art.accent" stroke-width="1.6" fill="none" opacity="0.7"/>
            <path d="M118 68 L110 68 L110 74" :stroke="art.accent" stroke-width="1.6" fill="none" opacity="0.7"/>
          </g>
        </template>

        <!-- 独角兽 -->
        <template v-else-if="art.variant === 'unicorn'">
          <path d="M118 116 Q136 110 142 98" :stroke="art.dark" stroke-width="4" fill="none" stroke-linecap="round"/>
          <ellipse cx="100" cy="118" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="18" ry="11" :fill="art.light" opacity="0.85"/>
          <ellipse cx="76" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="124" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M84 60 L78 46 L96 56 Z" :fill="art.main"/>
            <path d="M116 60 L122 46 L104 56 Z" :fill="art.main"/>
          </g>
          <path d="M100 54 L97 38 L103 38 Z" :fill="art.accent" :stroke="art.dark" stroke-width="1"/>
          <g v-if="level >= 5" class="mane">
            <path d="M82 66 Q74 80 82 92" :stroke="art.accent" stroke-width="3" fill="none" stroke-linecap="round" opacity="0.8"/>
            <path d="M78 70 Q70 84 78 96" :stroke="art.light" stroke-width="2.4" fill="none" stroke-linecap="round" opacity="0.7"/>
          </g>
        </template>

        <!-- 狮鹫·鹰头狮身 -->
        <template v-else-if="art.variant === 'griffin'">
          <g class="wings">
            <path d="M64 104 Q42 90 36 68 Q52 84 66 90 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
            <path d="M136 104 Q158 90 164 68 Q148 84 134 90 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
          </g>
          <path d="M126 118 Q148 108 152 94" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="32" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="20" ry="12" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <path d="M100 58 Q88 60 84 74 Q88 72 92 74 L100 82 L108 74 Q112 72 116 74 Q112 60 100 58 Z" :fill="art.accent"/>
          <g class="ears">
            <path d="M84 60 L78 46 L96 56 Z" :fill="art.main"/>
            <path d="M116 60 L122 46 L104 56 Z" :fill="art.main"/>
          </g>
        </template>

        <!-- 三角龙·颈盾三尖 -->
        <template v-else-if="art.variant === 'triceratops'">
          <ellipse cx="100" cy="122" rx="34" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="21" ry="12" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <path d="M100 104 Q66 100 58 88 Q74 82 100 82 Q126 82 142 88 Q134 100 100 104 Z" :fill="`url(#body-grad-${uid})`"/>
          <path d="M86 82 L80 70 M114 82 L120 70" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round"/>
          <path d="M100 80 L100 66" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round"/>
          <g v-if="level >= 5" class="frill">
            <path d="M62 84 Q60 66 72 58 Q84 64 84 82" :stroke="art.accent" stroke-width="2" fill="none"/>
            <path d="M138 84 Q140 66 128 58 Q116 64 116 82" :stroke="art.accent" stroke-width="2" fill="none"/>
          </g>
        </template>

        <!-- 猛犸象·长鼻象牙 -->
        <template v-else-if="art.variant === 'mammoth'">
          <path d="M100 120 Q84 126 78 134" :stroke="art.main" stroke-width="9" fill="none" stroke-linecap="round" opacity="0.9"/>
          <path d="M76 118 Q70 126 74 132" :stroke="art.light" stroke-width="4" fill="none" stroke-linecap="round"/>
          <ellipse cx="104" cy="122" rx="38" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="80" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="126" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <circle cx="100" cy="82" r="28" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 5" class="tusks">
            <path d="M90 96 Q84 104 88 108" :stroke="art.light" stroke-width="3.4" fill="none" stroke-linecap="round"/>
            <path d="M110 96 Q116 104 112 108" :stroke="art.light" stroke-width="3.4" fill="none" stroke-linecap="round"/>
          </g>
          <g class="ears">
            <ellipse cx="74" cy="70" rx="8" ry="6" :fill="art.main"/>
            <ellipse cx="126" cy="70" rx="8" ry="6" :fill="art.main"/>
          </g>
        </template>

        <!-- 剑齿虎·长獠牙 -->
        <template v-else-if="art.variant === 'sabertooth'">
          <path d="M126 118 Q150 106 154 90" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="19" ry="12" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M82 60 L74 44 L96 56 Z" :fill="art.main"/>
            <path d="M118 60 L126 44 L104 56 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="fangs">
            <path d="M90 96 L86 106" stroke="#fff" stroke-width="2.6" stroke-linecap="round"/>
            <path d="M110 96 L114 106" stroke="#fff" stroke-width="2.6" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 年兽·独角鬃毛 -->
        <template v-else-if="art.variant === 'nian'">
          <ellipse cx="100" cy="120" rx="32" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="20" ry="13" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <circle cx="100" cy="78" r="29" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 58 L72 42 L94 54 Z" :fill="art.main"/>
            <path d="M120 58 L128 42 L106 54 Z" :fill="art.main"/>
          </g>
          <path d="M100 52 L97 40 L103 40 Z" :fill="art.accent"/>
          <g v-if="level >= 5" class="shaggy">
            <path d="M70 82 Q62 72 72 66 M130 82 Q138 72 128 66" :stroke="art.accent" stroke-width="2.4" fill="none" stroke-linecap="round" opacity="0.7"/>
          </g>
        </template>

        <!-- 饕餮·血盆大口 -->
        <template v-else-if="art.variant === 'taotie'">
          <ellipse cx="100" cy="122" rx="31" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="20" ry="12" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <path d="M100 104 Q66 100 58 90 Q76 82 100 82 Q124 82 142 90 Q134 100 100 104 Z" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 5" class="maw">
            <path d="M84 96 Q100 104 116 96 Q108 108 100 108 Q92 108 84 96 Z" fill="#DC2626"/>
            <path d="M84 96 L88 100 L92 96 L96 100 L100 96 L104 100 L108 96 L112 100 L116 96" fill="none" stroke="#fff" stroke-width="1.4" opacity="0.8"/>
          </g>
          <g class="horns">
            <path d="M84 74 Q74 62 66 60 Q76 66 86 72 Z" :fill="art.accent"/>
            <path d="M116 74 Q126 62 134 60 Q124 66 114 72 Z" :fill="art.accent"/>
          </g>
        </template>

        <!-- 白泽·智慧独角 -->
        <template v-else-if="art.variant === 'baize'">
          <path d="M124 118 Q146 110 150 96" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="30" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="19" ry="12" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M84 60 L78 46 L96 56 Z" :fill="art.main"/>
            <path d="M116 60 L122 46 L104 56 Z" :fill="art.main"/>
          </g>
          <path d="M100 56 L96 42 L104 42 Z" :fill="art.accent"/>
          <path d="M88 88 Q74 92 66 86" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.5"/>
          <path d="M112 88 Q126 92 134 86" :stroke="art.dark" stroke-width="1.2" fill="none" opacity="0.5"/>
        </template>

        <!-- 黑豹·流线低伏 -->
        <template v-else-if="art.variant === 'panther'">
          <path d="M124 116 Q148 106 152 92 Q156 82 146 84" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="32" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="130" rx="20" ry="11" :fill="art.dark" opacity="0.5"/>
          <ellipse cx="78" cy="142" rx="6" ry="4" :fill="art.dark"/>
          <ellipse cx="122" cy="142" rx="6" ry="4" :fill="art.dark"/>
          <circle cx="100" cy="80" r="26" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M82 62 L74 44 L96 57 Z" :fill="art.main"/>
            <path d="M118 62 L126 44 L104 57 Z" :fill="art.main"/>
          </g>
        </template>

        <!-- 犰狳·甲带 -->
        <template v-else-if="art.variant === 'armadillo'">
          <path d="M86 134 Q74 134 68 130" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="122" rx="34" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="20" ry="12" :fill="art.dark" opacity="0.4"/>
          <g v-if="level >= 4" class="bands">
            <path d="M74 122 Q100 114 126 122" :stroke="art.dark" stroke-width="1.8" fill="none" opacity="0.6"/>
            <path d="M78 130 Q100 122 122 130" :stroke="art.dark" stroke-width="1.8" fill="none" opacity="0.6"/>
            <path d="M82 138 Q100 130 118 138" :stroke="art.dark" stroke-width="1.6" fill="none" opacity="0.5"/>
          </g>
          <ellipse cx="78" cy="144" rx="6" ry="4" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="6" ry="4" :fill="art.dark"/>
          <path d="M100 100 Q84 98 76 106 Q86 96 100 96 Z" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M86 92 Q82 82 88 78 Q92 84 90 92 Z" :fill="art.main"/>
            <path d="M114 92 Q118 82 112 78 Q108 84 110 92 Z" :fill="art.main"/>
          </g>
        </template>

        <!-- 藏羚羊·细长角 -->
        <template v-else-if="art.variant === 'antelope'">
          <path d="M118 112 Q132 108 136 98" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="118" rx="26" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="130" rx="16" ry="10" :fill="art.light" opacity="0.85"/>
          <ellipse cx="80" cy="141" rx="6" ry="4" :fill="art.dark"/>
          <ellipse cx="120" cy="141" rx="6" ry="4" :fill="art.dark"/>
          <circle cx="100" cy="80" r="24" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M86 62 L82 52 L94 58 Z" :fill="art.main"/>
            <path d="M114 62 L118 52 L106 58 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="thin-horns">
            <path d="M92 58 Q88 44 92 34" :stroke="art.dark" stroke-width="2.4" fill="none" stroke-linecap="round"/>
            <path d="M108 58 Q112 44 108 34" :stroke="art.dark" stroke-width="2.4" fill="none" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 雪豹·斑点 -->
        <template v-else-if="art.variant === 'leopard'">
          <path d="M126 116 Q150 104 154 88" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="19" ry="12" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="143" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M82 60 L74 44 L96 56 Z" :fill="art.main"/>
            <path d="M118 60 L126 44 L104 56 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 5" class="spots">
            <circle cx="86" cy="116" r="3" :fill="art.dark" opacity="0.5"/>
            <circle cx="100" cy="124" r="3" :fill="art.dark" opacity="0.5"/>
            <circle cx="114" cy="116" r="3" :fill="art.dark" opacity="0.5"/>
            <circle cx="92" cy="62" r="2.6" :fill="art.dark" opacity="0.5"/>
            <circle cx="108" cy="64" r="2.6" :fill="art.dark" opacity="0.5"/>
          </g>
        </template>

        <!-- 小熊猫·条纹尾 -->
        <template v-else-if="art.variant === 'redpanda'">
          <path d="M124 118 Q146 108 150 92 Q154 82 144 84" :fill="art.dark" class="tail"/>
          <g v-if="level >= 5" class="tail-rings">
            <path d="M144 88 L140 92 M138 94 L134 98" :stroke="art.accent" stroke-width="2.6" stroke-linecap="round" opacity="0.8"/>
          </g>
          <ellipse cx="100" cy="122" rx="28" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="18" ry="11" :fill="art.light" opacity="0.85"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="82" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <circle cx="78" cy="60" r="9" :fill="art.main" :stroke="art.dark" stroke-width="1.2"/>
            <circle cx="122" cy="60" r="9" :fill="art.main" :stroke="art.dark" stroke-width="1.2"/>
          </g>
          <g v-if="level >= 4" class="face-mask">
            <path d="M86 88 Q94 84 100 88 Q106 84 114 88" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.5"/>
          </g>
        </template>

        <!-- 梦魇马 -->
        <template v-else-if="art.variant === 'horse'">
          <path d="M118 114 Q132 108 138 98" :stroke="art.dark" stroke-width="4" fill="none" stroke-linecap="round"/>
          <ellipse cx="100" cy="118" rx="30" ry="24" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="130" rx="18" ry="11" :fill="art.light" opacity="0.85"/>
          <ellipse cx="76" cy="141" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="124" cy="141" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M84 60 L78 44 L96 56 Z" :fill="art.main"/>
            <path d="M116 60 L122 44 L104 56 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="mane">
            <path d="M82 62 Q72 76 82 90" :stroke="art.accent" stroke-width="3" fill="none" stroke-linecap="round" opacity="0.8"/>
            <path d="M118 62 Q128 76 118 90" :stroke="art.accent" stroke-width="3" fill="none" stroke-linecap="round" opacity="0.8"/>
          </g>
        </template>

        <!-- 地懒·巨爪 -->
        <template v-else-if="art.variant === 'sloth'">
          <ellipse cx="100" cy="122" rx="34" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="21" ry="13" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <circle cx="100" cy="82" r="28" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <circle cx="76" cy="62" r="8" :fill="art.main"/>
            <circle cx="124" cy="62" r="8" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="claws">
            <path d="M76 142 L72 148 M80 144 L78 150 M120 142 L124 148 M116 144 L118 150" :stroke="art.light" stroke-width="2.6" stroke-linecap="round" opacity="0.9"/>
          </g>
        </template>

        <!-- 狮·鬃毛 -->
        <template v-else-if="art.variant === 'lion'">
          <ellipse cx="100" cy="120" rx="32" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="20" ry="13" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="8" ry="6" :fill="art.dark"/>
          <g v-if="level >= 4" class="mane-ring">
            <circle cx="100" cy="80" r="34" :fill="art.accent" opacity="0.3"/>
            <path d="M78 58 Q70 48 78 40 M122 58 Q130 48 122 40 M92 48 L100 40 L108 48" :stroke="art.accent" stroke-width="3" stroke-linecap="round" opacity="0.7"/>
          </g>
          <circle cx="100" cy="80" r="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <circle cx="78" cy="60" r="8" :fill="art.main"/>
            <circle cx="122" cy="60" r="8" :fill="art.main"/>
          </g>
          <path d="M100 100 Q94 108 100 112 Q106 108 100 100 Z" fill="#7C2D12" opacity="0.6"/>
        </template>

        <!-- 西伯利亚虎·耳簇厚毛 -->
        <template v-else-if="art.variant === 'siberian'">
          <path d="M126 116 Q150 104 154 88" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="32" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="20" ry="13" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="122" cy="144" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="78" r="29" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 58 L72 42 L94 54 Z" :fill="art.main"/>
            <path d="M120 58 L128 42 L106 54 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 5" class="fur">
            <path d="M72 64 Q66 58 72 52 M128 64 Q134 58 128 52" :stroke="art.dark" stroke-width="2" stroke-linecap="round" opacity="0.7"/>
            <path d="M70 82 Q62 76 70 70 M130 82 Q138 76 130 70" :stroke="art.dark" stroke-width="2" stroke-linecap="round" opacity="0.6"/>
          </g>
        </template>

        <!-- 披毛犀·长毛双角 -->
        <template v-else-if="art.variant === 'woolly'">
          <ellipse cx="100" cy="122" rx="34" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="21" ry="12" :fill="art.light" opacity="0.8"/>
          <ellipse cx="78" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <path d="M100 102 Q78 98 70 88 Q80 82 100 82 Q120 82 130 88 Q122 98 100 102 Z" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 5" class="wool">
            <path d="M66 92 Q58 88 66 82 M134 92 Q142 88 134 82" :stroke="art.dark" stroke-width="2.4" stroke-linecap="round" opacity="0.6"/>
            <path d="M70 104 Q62 100 70 94 M130 104 Q138 100 130 94" :stroke="art.dark" stroke-width="2.4" stroke-linecap="round" opacity="0.6"/>
          </g>
          <path d="M100 80 L97 66 M100 80 L103 72" :stroke="art.dark" stroke-width="2.6" stroke-linecap="round"/>
        </template>

        <!-- 鹿·鹿角 -->
        <template v-else-if="art.variant === 'deer'">
          <path d="M118 112 Q132 106 138 96" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="118" rx="26" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="130" rx="16" ry="10" :fill="art.light" opacity="0.85"/>
          <ellipse cx="80" cy="141" rx="6" ry="4" :fill="art.dark"/>
          <ellipse cx="120" cy="141" rx="6" ry="4" :fill="art.dark"/>
          <circle cx="100" cy="80" r="24" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M86 62 L82 52 L94 58 Z" :fill="art.main"/>
            <path d="M114 62 L118 52 L106 58 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="antlers">
            <path d="M92 58 L86 42 L80 34 M86 46 L76 40" :stroke="art.dark" stroke-width="2.2" fill="none" stroke-linecap="round"/>
            <path d="M108 58 L114 42 L120 34 M114 46 L124 40" :stroke="art.dark" stroke-width="2.2" fill="none" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 犀牛·甲身单角 -->
        <template v-else-if="art.variant === 'rhino'">
          <ellipse cx="100" cy="122" rx="34" ry="25" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="21" ry="12" :fill="art.dark" opacity="0.4"/>
          <ellipse cx="78" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <ellipse cx="122" cy="145" rx="8" ry="6" :fill="art.dark"/>
          <path d="M100 102 Q78 98 70 90 Q82 82 100 82 Q118 82 130 90 Q122 98 100 102 Z" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 4" class="plates">
            <path d="M80 96 L86 90 L92 96 L98 90 L104 96" :stroke="art.dark" stroke-width="1.8" fill="none" opacity="0.5"/>
          </g>
          <path d="M100 80 L97 64 L103 64 Z" :fill="art.accent"/>
          <g class="ears">
            <ellipse cx="82" cy="70" rx="6" ry="5" :fill="art.main" transform="rotate(-15 82 70)"/>
            <ellipse cx="118" cy="70" rx="6" ry="5" :fill="art.main" transform="rotate(15 118 70)"/>
          </g>
        </template>

        <!-- 亚古兽·橙色小恐龙 -->
        <template v-else-if="art.variant === 'agumon'">
          <ellipse cx="100" cy="122" rx="24" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="133" rx="15" ry="10" :fill="art.light" opacity="0.8"/>
          <ellipse cx="82" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="118" cy="142" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="78" r="28" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M84 56 L76 38 L98 52 Z" :fill="art.main"/>
            <path d="M116 56 L124 38 L102 52 Z" :fill="art.main"/>
          </g>
          <path d="M100 102 Q92 106 100 110 Q108 106 100 102 Z" fill="#B45309" opacity="0.55"/>
          <g v-if="level >= 6">
            <path d="M88 54 Q82 42 78 40" :stroke="art.accent" stroke-width="2.4" fill="none" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 迪路兽·神圣白猫 -->
        <template v-else-if="art.variant === 'tailmon'">
          <path d="M124 118 Q146 108 150 92" :fill="art.dark" class="tail"/>
          <circle cx="150" cy="88" r="4" :fill="art.accent" :stroke="art.dark" stroke-width="1"/>
          <ellipse cx="100" cy="122" rx="26" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="132" rx="16" ry="10" :fill="art.light" opacity="0.9"/>
          <ellipse cx="80" cy="143" rx="6" ry="4" :fill="art.dark"/>
          <ellipse cx="120" cy="143" rx="6" ry="4" :fill="art.dark"/>
          <circle cx="100" cy="82" r="26" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M80 62 L72 44 L94 57 Z" :fill="art.main"/>
            <path d="M120 62 L128 44 L106 57 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 4" class="sacred-ring">
            <circle cx="100" cy="56" r="4.5" :fill="art.accent" :stroke="art.dark" stroke-width="1.2"/>
          </g>
        </template>

        <!-- 巴达兽·奶油小天使 -->
        <template v-else-if="art.variant === 'patamon'">
          <g class="wings">
            <path d="M70 94 Q48 84 44 68 Q58 82 72 86 Z" :fill="art.light" :stroke="art.dark" stroke-width="1"/>
            <path d="M130 94 Q152 84 156 68 Q142 82 128 86 Z" :fill="art.light" :stroke="art.dark" stroke-width="1"/>
          </g>
          <ellipse cx="100" cy="112" rx="30" ry="28" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="124" rx="19" ry="12" :fill="art.light" opacity="0.9"/>
          <g class="ears">
            <path d="M74 88 L64 66 L88 82 Z" :fill="art.main"/>
            <path d="M126 88 L136 66 L112 82 Z" :fill="art.main"/>
          </g>
          <g v-if="level >= 5" class="wing-halo">
            <circle cx="100" cy="84" r="17" fill="none" :stroke="art.accent" stroke-width="1.4" opacity="0.6"/>
          </g>
        </template>

        <!-- 加布兽·披毛蓝狼 -->
        <template v-else-if="art.variant === 'gabumon'">
          <path d="M118 118 Q136 112 140 100" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="26" ry="23" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="16" ry="11" :fill="art.light" opacity="0.85"/>
          <ellipse cx="80" cy="142" rx="6" ry="4" :fill="art.dark"/>
          <ellipse cx="120" cy="142" rx="6" ry="4" :fill="art.dark"/>
          <circle cx="100" cy="80" r="26" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <path d="M84 60 L78 44 L96 56 Z" :fill="art.main"/>
            <path d="M116 60 L122 44 L104 56 Z" :fill="art.main"/>
          </g>
          <path d="M78 96 Q62 92 60 80 Q72 84 78 92 Z" :fill="art.main" opacity="0.75"/>
          <path d="M122 96 Q138 92 140 80 Q128 84 122 92 Z" :fill="art.main" opacity="0.75"/>
          <path d="M100 54 L97 42 L103 42 Z" :fill="art.accent"/>
        </template>

        <!-- 默认四足兽 -->
        <template v-else>
          <path v-if="art.parts.tail" d="M126 120 Q148 108 152 92 Q156 82 146 84" :fill="art.dark" class="tail" />
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
          <ellipse cx="100" cy="95" rx="7" ry="4.5" :fill="art.dark" class="roar-mouth" />
          <ellipse cx="127" cy="133" rx="6" ry="4" :fill="art.dark" class="paw-move" />
        </template>

        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
      </g>

      <!-- ---------- 龙 / 长形（含 11 个独立剪影变体） ---------- -->
      <g v-else-if="art.body === 'dragon'" class="dragon-group">

        <!-- 烛龙·蛇身火焰 -->
        <template v-if="art.variant === 'serpent'">
          <g v-if="level >= 5" class="flame-mane">
            <path d="M74 68 Q62 50 78 40 Q80 54 74 68 Z" fill="#FBBF24" opacity="0.9"/>
            <path d="M84 60 Q74 38 96 32 Q92 48 84 60 Z" fill="#F59E0B" opacity="0.85"/>
            <path d="M126 66 Q138 48 122 38 Q120 52 126 66 Z" fill="#FDE68A" opacity="0.9"/>
            <path d="M116 58 Q128 40 108 34 Q112 50 116 58 Z" fill="#F97316" opacity="0.85"/>
          </g>
          <path d="M100 138 C 46 132, 54 98, 100 94 C 150 90, 148 66, 100 64" fill="none" :stroke="`url(#body-grad-${uid})`" stroke-width="22" stroke-linecap="round"/>
          <path d="M100 138 C 46 132, 54 98, 100 94 C 150 90, 148 66, 100 64" fill="none" :stroke="art.dark" stroke-width="26" stroke-linecap="round" opacity="0.14"/>
          <path d="M100 138 C 46 132, 54 98, 100 94 C 150 90, 148 66, 100 64" fill="none" :stroke="art.light" stroke-width="8" stroke-linecap="round" opacity="0.4"/>
          <ellipse cx="100" cy="78" rx="20" ry="17" :fill="`url(#body-grad-${uid})`"/>
          <path d="M100 61 Q76 72 76 92 Q100 103 124 92 Q124 72 100 61 Z" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 8" class="fire-eye">
            <circle cx="91" cy="80" r="2.6" fill="#FDE047"/>
            <circle cx="109" cy="80" r="2.6" fill="#FDE047"/>
          </g>
        </template>

        <!-- 应龙·羽翼鹿角 -->
        <template v-else-if="art.variant === 'feathered'">
          <g class="wings">
            <path d="M64 96 Q42 82 36 60 Q54 78 64 84 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
            <path d="M60 92 Q38 92 32 72 Q52 82 60 88 Z" :fill="art.main" opacity="0.65"/>
            <path d="M136 96 Q158 82 164 60 Q146 78 136 84 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
            <path d="M140 92 Q162 92 168 72 Q148 82 140 88 Z" :fill="art.main" opacity="0.65"/>
          </g>
          <path d="M126 118 Q154 108 162 92 Q168 80 158 82" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="124" rx="42" ry="23" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="131" rx="28" ry="12" :fill="art.light" opacity="0.8"/>
          <path d="M84 86 Q62 88 54 80" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.5"/>
          <path d="M116 86 Q138 88 146 80" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.5"/>
          <ellipse cx="100" cy="78" rx="32" ry="27" :fill="`url(#body-grad-${uid})`"/>
          <g class="horns">
            <path d="M82 56 Q72 40 62 36 Q74 44 82 54 Q76 44 70 36" :stroke="art.accent" stroke-width="2.6" fill="none" stroke-linecap="round"/>
            <path d="M118 56 Q128 40 138 36 Q126 44 118 54 Q124 44 130 36" :stroke="art.accent" stroke-width="2.6" fill="none" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 青龙·鹿角须+背鳍 -->
        <template v-else-if="art.variant === 'antler'">
          <path d="M100 138 C 48 130, 52 92, 100 88 C 152 84, 150 60, 100 56" fill="none" :stroke="`url(#body-grad-${uid})`" stroke-width="21" stroke-linecap="round"/>
          <g v-if="level >= 6" class="dorsal">
            <path d="M64 110 Q60 96 66 92 Q72 100 66 108 Z" :fill="art.accent" opacity="0.8"/>
            <path d="M128 104 Q132 88 126 84 Q120 92 128 102 Z" :fill="art.accent" opacity="0.8"/>
          </g>
          <path d="M82 88 Q60 92 52 84" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.5"/>
          <path d="M118 88 Q140 92 148 84" :stroke="art.dark" stroke-width="1.4" fill="none" opacity="0.5"/>
          <ellipse cx="100" cy="78" rx="19" ry="16" :fill="`url(#body-grad-${uid})`"/>
          <g class="horns">
            <path d="M86 64 Q80 46 70 40 Q78 48 88 62 M86 62 Q92 46 98 40" :stroke="art.accent" stroke-width="2.2" fill="none" stroke-linecap="round"/>
            <path d="M114 64 Q120 46 130 40 Q122 48 112 62 M114 62 Q108 46 102 40" :stroke="art.accent" stroke-width="2.2" fill="none" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 机械龙·机甲分段 -->
        <template v-else-if="art.variant === 'mecha'">
          <g class="wings">
            <path d="M64 108 L40 84 L66 98 Z" :fill="art.accent" opacity="0.75"/>
            <path d="M136 108 L160 84 L134 98 Z" :fill="art.accent" opacity="0.75"/>
          </g>
          <path d="M112 124 L128 134 L124 120 Z" :fill="art.dark" opacity="0.7"/>
          <rect x="70" y="112" width="60" height="30" rx="8" :fill="art.main" :stroke="art.dark" stroke-width="1.6"/>
          <rect x="76" y="120" width="48" height="12" rx="4" :fill="art.dark" opacity="0.5"/>
          <circle cx="88" cy="124" r="2.6" :fill="art.accent" class="joint"/>
          <circle cx="112" cy="124" r="2.6" :fill="art.accent" class="joint"/>
          <ellipse cx="100" cy="78" rx="30" ry="26" :fill="`url(#body-grad-${uid})`"/>
          <rect x="86" y="82" width="28" height="8" rx="4" :fill="art.accent" class="visor"/>
          <g v-if="art.parts.horns">
            <line x1="90" y1="54" x2="84" y2="38" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round"/>
            <line x1="110" y1="54" x2="116" y2="38" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 翼龙/飞龙·蝠翼尾刺 -->
        <template v-else-if="art.variant === 'wyvern'">
          <g class="wings">
            <path d="M64 104 Q38 84 32 58 Q48 78 64 88 Z" :fill="art.dark"/>
            <path d="M136 104 Q162 84 168 58 Q152 78 136 88 Z" :fill="art.dark"/>
          </g>
          <path d="M124 126 Q150 118 158 106 L162 100 L156 112 Q150 122 130 130 Z" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="122" rx="34" ry="26" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="134" rx="20" ry="10" :fill="art.light" opacity="0.8"/>
          <ellipse cx="100" cy="76" rx="28" ry="24" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="art.parts.horns">
            <path d="M86 56 Q78 40 66 34 Q78 42 88 54 Z" :fill="art.accent"/>
            <path d="M114 56 Q122 40 134 34 Q122 42 112 54 Z" :fill="art.accent"/>
          </g>
          <g v-if="level >= 7">
            <path d="M100 66 L100 84" :stroke="art.accent" stroke-width="2.4" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 扬子鳄·低伏宽吻 -->
        <template v-else-if="art.variant === 'crocodile'">
          <path d="M96 136 Q60 130 46 122 Q60 118 88 120 Z" :fill="art.dark" class="tail"/>
          <g v-if="level >= 6" class="ridges">
            <path d="M88 112 L94 106 L100 112 L106 106 L112 112" :stroke="art.dark" stroke-width="1.6" fill="none" opacity="0.6"/>
          </g>
          <ellipse cx="100" cy="120" rx="40" ry="17" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <path d="M100 103 Q70 103 62 118 Q70 133 100 133 Q130 133 138 118 Q130 103 100 103 Z" :fill="`url(#body-grad-${uid})`"/>
          <path d="M86 112 Q80 106 74 108 Q80 118 86 122" :stroke="art.dark" stroke-width="1.4" fill="none"/>
          <path d="M114 112 Q120 106 126 108 Q120 118 114 122" :stroke="art.dark" stroke-width="1.4" fill="none"/>
          <g v-if="level >= 8" class="snout">
            <path d="M96 118 Q100 114 104 118 Q100 121 96 118 Z" :fill="art.dark" opacity="0.6"/>
          </g>
        </template>

        <!-- 玄武·龟蛇合体 -->
        <template v-else-if="art.variant === 'xuanwu'">
          <path d="M96 118 C 72 112, 56 118, 52 130" :stroke="art.dark" stroke-width="5" fill="none" stroke-linecap="round" opacity="0.6"/>
          <path d="M58 118 Q56 106 66 102" :stroke="art.dark" stroke-width="4" fill="none" stroke-linecap="round"/>
          <path d="M52 128 Q46 124 48 118" :stroke="art.dark" stroke-width="3.4" fill="none" stroke-linecap="round"/>
          <ellipse cx="104" cy="124" rx="40" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <path d="M74 124 Q100 96 126 124 Q140 130 130 138 Q100 150 70 138 Q60 130 74 124 Z" :fill="art.main" :stroke="art.dark" stroke-width="1.6"/>
          <path d="M88 124 L100 108 L112 124 L100 140 Z" :fill="art.light" opacity="0.5"/>
          <ellipse cx="100" cy="78" rx="20" ry="17" :fill="`url(#body-grad-${uid})`"/>
          <path d="M92 92 Q100 100 108 92" :stroke="art.dark" stroke-width="1.4" fill="none"/>
          <g v-if="art.parts.horns">
            <path d="M94 64 Q90 54 96 50" :stroke="art.accent" stroke-width="2" fill="none" stroke-linecap="round"/>
            <path d="M106 64 Q110 54 104 50" :stroke="art.accent" stroke-width="2" fill="none" stroke-linecap="round"/>
          </g>
        </template>

        <!-- 霸王龙·直立小短手 -->
        <template v-else-if="art.variant === 'trex'">
          <path d="M112 128 Q132 128 138 118 Q134 132 122 138 Z" :fill="art.dark" class="tail"/>
          <path d="M74 132 Q66 128 64 120 Q60 132 70 138 Z" :fill="art.main" :stroke="art.dark" stroke-width="1.2"/>
          <path d="M126 130 Q132 130 134 124 Q130 134 122 138 Z" :fill="art.main" :stroke="art.dark" stroke-width="1.2"/>
          <ellipse cx="100" cy="122" rx="26" ry="20" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="134" rx="18" ry="8" :fill="art.light" opacity="0.8"/>
          <path d="M78 116 Q76 108 82 106" :stroke="art.dark" stroke-width="3.2" fill="none" stroke-linecap="round"/>
          <path d="M122 116 Q124 108 118 106" :stroke="art.dark" stroke-width="3.2" fill="none" stroke-linecap="round"/>
          <path d="M92 92 Q78 76 88 66 Q100 76 98 92 Z" :fill="`url(#body-grad-${uid})`"/>
          <path d="M102 92 Q116 70 100 58 Q92 70 98 92 Z" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 6" class="jaws">
            <path d="M92 70 L78 82 M108 66 L116 76" stroke="#fff" stroke-width="1.6" opacity="0.7"/>
          </g>
        </template>

        <!-- 棘龙·背帆 -->
        <template v-else-if="art.variant === 'sailback'">
          <g v-if="level >= 4" class="sail">
            <path d="M70 116 Q84 60 100 56 Q112 62 130 116 L120 118 L112 76 L108 118 L104 88 L100 118 L96 88 L92 118 L88 76 L80 118 Z" :fill="art.accent" opacity="0.85" :stroke="art.dark" stroke-width="1"/>
          </g>
          <path d="M92 132 Q70 130 56 124 Q68 120 88 122 Z" :fill="art.dark" class="tail"/>
          <ellipse cx="100" cy="120" rx="36" ry="16" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <path d="M100 104 Q74 104 66 118 Q74 132 100 132 Q126 132 134 118 Q126 104 100 104 Z" :fill="`url(#body-grad-${uid})`"/>
        </template>

        <!-- 甲龙·甲背尾锤 -->
        <template v-else-if="art.variant === 'ankylo'">
          <circle cx="134" cy="126" r="8" :fill="art.dark" class="tail-club"/>
          <path d="M120 128 Q132 128 134 124" :stroke="art.dark" stroke-width="6" fill="none" stroke-linecap="round"/>
          <ellipse cx="100" cy="120" rx="40" ry="20" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <g v-if="level >= 5" class="armor">
            <path d="M66 122 Q80 106 100 106 Q120 106 134 122 Q118 116 100 116 Q82 116 66 122 Z" :fill="art.accent" opacity="0.6"/>
            <path d="M78 112 Q88 104 100 104 Q112 104 122 112" :stroke="art.dark" stroke-width="1.6" fill="none" opacity="0.7"/>
          </g>
          <ellipse cx="100" cy="76" rx="24" ry="20" :fill="`url(#body-grad-${uid})`"/>
        </template>

        <!-- 梁龙·长颈 -->
        <template v-else-if="art.variant === 'sauropod'">
          <path d="M100 124 Q100 96 96 78 Q92 62 100 52" fill="none" :stroke="`url(#body-grad-${uid})`" stroke-width="17" stroke-linecap="round"/>
          <ellipse cx="100" cy="124" rx="30" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="80" cy="136" rx="7" ry="5" :fill="art.dark"/>
          <ellipse cx="120" cy="136" rx="7" ry="5" :fill="art.dark"/>
          <circle cx="100" cy="52" r="12" :fill="`url(#body-grad-${uid})`"/>
          <g v-if="level >= 7" class="spots">
            <circle cx="110" cy="122" r="3" :fill="art.dark" opacity="0.3"/>
            <circle cx="92" cy="128" r="3" :fill="art.dark" opacity="0.3"/>
          </g>
        </template>

        <!-- 默认龙（未指定变体） -->
        <template v-else>
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
          <g v-if="art.parts.horns" class="horns">
            <path d="M84 56 Q76 34 64 28 Q76 38 88 54 Z" :fill="art.accent" />
            <path d="M116 56 Q124 34 136 28 Q124 38 112 54 Z" :fill="art.accent" />
          </g>
          <g v-if="art.parts.crown && level >= 7" class="flame">
            <path d="M100 58 Q94 44 100 32 Q106 44 100 58 Z" fill="#FBBF24" opacity="0.9" />
            <path d="M100 56 Q97 48 100 40 Q103 48 100 56 Z" fill="#FDE68A" />
          </g>
        </template>

        <g v-html="haloSvg" v-if="art.parts.halo" />
        <g v-html="faceSvg" />
      </g>

      <!-- ---------- 鸟类 ---------- -->
      <g v-else-if="art.body === 'bird'" class="bird-group">
        <!-- 比丘兽·粉色小鸟 -->
        <template v-if="art.variant === 'biyomon'">
          <g class="wings">
            <path d="M70 104 Q48 94 44 76 Q58 90 72 94 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
            <path d="M130 104 Q152 94 156 76 Q142 90 128 94 Z" :fill="art.light" :stroke="art.dark" stroke-width="1.2"/>
          </g>
          <path d="M128 116 L148 122 L134 130 Z" :fill="art.dark"/>
          <ellipse cx="100" cy="116" rx="26" ry="22" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="125" rx="16" ry="11" :fill="art.light" opacity="0.85"/>
          <circle cx="100" cy="80" r="24" :fill="`url(#body-grad-${uid})`"/>
          <path d="M92 60 Q96 48 104 48 Q108 56 104 60 Z" :fill="art.accent"/>
          <path d="M96 90 Q100 96 104 90 Q100 94 96 90 Z" fill="#F97316"/>
          <g v-if="level >= 6" class="crest">
            <path d="M96 54 Q92 44 96 38 Q100 46 100 54 Z" :fill="art.accent" opacity="0.9"/>
          </g>
        </template>
        <template v-else>
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
        </template>
      </g>

      <!-- ---------- 水生 ---------- -->
      <g v-else-if="art.body === 'aquatic'" class="aquatic-group">
        <!-- 哥玛兽·白色海狮 -->
        <template v-if="art.variant === 'gomamon'">
          <path d="M64 122 Q48 118 46 128" :fill="art.dark"/>
          <ellipse cx="100" cy="120" rx="30" ry="20" :fill="`url(#body-grad-${uid})`" class="body-main"/>
          <ellipse cx="100" cy="129" rx="18" ry="10" :fill="art.light" opacity="0.9"/>
          <circle cx="100" cy="94" r="19" :fill="`url(#body-grad-${uid})`"/>
          <g class="ears">
            <ellipse cx="84" cy="78" rx="6" ry="5" :fill="art.main"/>
            <ellipse cx="116" cy="78" rx="6" ry="5" :fill="art.main"/>
          </g>
          <path d="M108 98 Q112 102 116 98" :stroke="art.dark" stroke-width="1.4" fill="none"/>
          <path d="M92 98 Q88 102 84 98" :stroke="art.dark" stroke-width="1.4" fill="none"/>
        </template>
        <template v-else>
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
        </template>
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
