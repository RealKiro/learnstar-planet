// ===== 学趣星球 · 宠物 SVG 部件与参数体系 =====
// 每个物种只需配置「身体模板 + 主色 + 部件出现等级」，
// 其余(配色渐变/体型/光效/表情)由等级参数公式驱动，
// 从而以"部件复用 + 参数渐变"覆盖全量物种，而无需 120×12 张手绘图。

import { getLevelStage, getSeriesBySpeciesId } from './petData'

/** 身体模板类型 */
export type BodyType =
  | 'dragon'     // 蛇/龙/长形爬虫(烛龙 应龙 翼龙 沧龙 扬子鳄 机械龙)
  | 'quadruped'  // 四足兽(猫/狐/虎/熊/鹿/马/牛/羊)
  | 'bird'       // 鸟类(凤凰 丹顶鹤 朱鹮 隼 鸽)
  | 'aquatic'    // 水生(鲲鹏 机械鲨 巨齿鲨 江豚 龙舟)
  | 'mecha'      // 机甲(星际机甲 裂变巨人 纳米虫群 光剑武士)
  | 'spirit'     // 灵体/星灵(量子兽 星座守护 精灵 梦魇)
  | 'food'       // 物化萌物(灯笼 饺子 月饼 粽子 汤圆)
  | 'humanoid'   // 人形(精灵 树人 美人鱼 福星 寿星 圣斗士)

/** 物种艺术配置 */
export interface PetArtConfig {
  /** 身体模板 */
  body: BodyType
  /** 主色(其余色阶由公式推导) */
  main: string
  /** 剪影变体：同身体模板下的独立外形（宝可梦/七侠/神话龙等需手绘贴合原型，跨阶段也随 variant 演化） */
  variant?: string
  /** 各部件出现的等级阈值 */
  parts?: Partial<Record<'horns' | 'wings' | 'tail' | 'fin' | 'halo' | 'crown', number>>
  /** 是否抽象灵体(无实体身体，如量子/星灵) */
  abstract?: boolean
}

/** 渲染参数(由 computePetArt 计算) */
export interface PetArtRender {
  body: BodyType
  /** 剪影变体(与 body 模板配合，默认空=通用模板) */
  variant: string
  /** 所属系列 id（特效层/圣衣/剑器选型） */
  seriesId: string
  main: string
  accent: string
  dark: string
  light: string
  /** 体型放大系数 */
  scale: number
  /** 光晕强度 0-1 */
  glow: number
  /** 部件是否出现 */
  parts: { horns: boolean; wings: boolean; tail: boolean; fin: boolean; halo: boolean; crown: boolean }
  stage: 'egg' | 'baby' | 'growing' | 'mature' | 'legendary'
  /** 眼睛风格 */
  eyes: 'big' | 'normal' | 'sharp' | 'serene'
  /** 腮红(幼年) */
  blush: boolean
  /** 是否有小星星特效 */
  sparkle: boolean
}

// ============================================================
// 颜色工具
// ============================================================

/** 十六进制颜色明暗调整(percent: -1 加深 ~ 1 提亮) */
function shade(hex: string, percent: number): string {
  const h = hex.replace('#', '')
  const num = parseInt(h.length === 3 ? h.split('').map(c => c + c).join('') : h, 16)
  const r = Math.round((num >> 16) & 255)
  const g = Math.round((num >> 8) & 255)
  const b = Math.round(num & 255)
  const t = percent < 0 ? 0 : 255
  const p = Math.abs(percent)
  const mix = (c: number) => Math.round((t - c) * p + c)
  return `#${((1 << 24) + (mix(r) << 16) + (mix(g) << 8) + mix(b)).toString(16).slice(1)}`
}

// ============================================================
// 全物种艺术配置(主色即可,色阶自动推导)
// ============================================================

const DEFAULT_PARTS = { horns: 6, wings: 8, tail: 3, fin: 6, halo: 11, crown: 10 }

export const PET_ART: Record<string, PetArtConfig> = {
  // ===== 古代神话 =====
  zhulong:       { body: 'dragon', main: '#F59E0B', variant: 'serpent', parts: { horns: 6, wings: 11, tail: 3 } },
  yinglong:      { body: 'dragon', main: '#10B981', variant: 'feathered', parts: { horns: 6, wings: 7, tail: 3 } },
  nine_tail_fox: { body: 'quadruped', main: '#F1F5F9', variant: 'kitsune', parts: { tail: 3, horns: 10 } },
  kunpeng:       { body: 'aquatic', main: '#3B82F6', parts: { tail: 3, fin: 5, wings: 8 } },
  fenghuang:     { body: 'bird', main: '#EF4444', parts: { tail: 5, wings: 6, crown: 7 } },
  qilin:         { body: 'quadruped', main: '#10B981', variant: 'qilin', parts: { horns: 5, tail: 3, wings: 11 } },
  // ===== 元素精灵(宝可梦风，贴合原型) =====
  charmander:    { body: 'quadruped', main: '#F97316', variant: 'lizard', parts: { tail: 3, horns: 7 } },
  bulbasaur:     { body: 'quadruped', main: '#22C55E', variant: 'seed', parts: { tail: 4, crown: 6 } },
  squirtle:      { body: 'quadruped', main: '#3B82F6', variant: 'turtle', parts: { tail: 3, crown: 7 } },
  eevee:         { body: 'quadruped', main: '#B45309', variant: 'fox', parts: { tail: 3, horns: 9 } },
  pikachu:       { body: 'quadruped', main: '#FACC15', variant: 'mouse', parts: { tail: 3, horns: 9 } },
  riolu:         { body: 'quadruped', main: '#6366F1', variant: 'pup', parts: { tail: 4, horns: 8 } },
  // ===== 国宝守护 =====
  panda:               { body: 'quadruped', main: '#F1F5F9', variant: 'panda', parts: { tail: 4, crown: 10 } },
  golden_monkey:       { body: 'quadruped', main: '#F59E0B', variant: 'monkey', parts: { tail: 3, crown: 9 } },
  red_crowned_crane:   { body: 'bird', main: '#F1F5F9', parts: { wings: 5, tail: 6, crown: 8 } },
  south_china_tiger:   { body: 'quadruped', main: '#F97316', variant: 'tiger', parts: { tail: 3, horns: 8, wings: 11 } },
  chinese_alligator:   { body: 'dragon', main: '#64748B', variant: 'crocodile', parts: { tail: 3, horns: 9 } },
  crested_ibis:        { body: 'bird', main: '#FBCFE8', parts: { wings: 5, tail: 7, crown: 11 } },
  // ===== 科幻机甲 =====
  mecha_dragon:  { body: 'dragon', main: '#3B82F6', variant: 'mecha', parts: { horns: 5, tail: 3, wings: 8 } },
  cyber_cat:     { body: 'quadruped', main: '#8B5CF6', variant: 'cyber', parts: { tail: 3, horns: 7, wings: 10 } },
  space_mecha:   { body: 'mecha', main: '#60A5FA', parts: { wings: 8, crown: 9 } },
  quantum_beast: { body: 'spirit', main: '#A855F7', abstract: true, parts: { halo: 5 } },
  digital_phoenix: { body: 'bird', main: '#06B6D4', parts: { wings: 5, tail: 6, crown: 9 } },
  mecha_shark:   { body: 'aquatic', main: '#94A3B8', parts: { tail: 3, fin: 5, crown: 9 } },
  // ===== 魔法奇幻 =====
  unicorn:       { body: 'quadruped', main: '#F5F3FF', variant: 'unicorn', parts: { horns: 5, tail: 4, wings: 8 } },
  wyvern:        { body: 'dragon', main: '#EF4444', variant: 'wyvern', parts: { horns: 5, wings: 6, tail: 3 } },
  fairy:         { body: 'humanoid', main: '#FBBF24', parts: { wings: 5, crown: 7 } },
  treant:        { body: 'humanoid', main: '#166534', parts: { horns: 7, crown: 10 } },
  griffin:       { body: 'quadruped', main: '#F59E0B', variant: 'griffin', parts: { wings: 6, tail: 4, crown: 8 } },
  mermaid:       { body: 'humanoid', main: '#0EA5E9', parts: { tail: 3, crown: 7, wings: 11 } },
  // ===== 史前生物 =====
  t_rex:         { body: 'dragon', main: '#65A30D', variant: 'trex', parts: { tail: 3, horns: 7, crown: 9 } },
  triceratops:   { body: 'quadruped', main: '#F59E0B', variant: 'triceratops', parts: { horns: 4, tail: 4, crown: 9 } },
  pterosaur:     { body: 'bird', main: '#60A5FA', parts: { wings: 5, tail: 5 } },
  mammoth:       { body: 'quadruped', main: '#A16207', variant: 'mammoth', parts: { horns: 5, tail: 4, crown: 9 } },
  sabertooth:    { body: 'quadruped', main: '#92400E', variant: 'sabertooth', parts: { horns: 4, tail: 3, crown: 9 } },
  mosasaur:      { body: 'aquatic', main: '#0EA5E9', parts: { tail: 3, fin: 5, crown: 8 } },
  // ===== 星座守护(圣斗士·圣衣) =====
  aries:         { body: 'spirit', main: '#F59E0B', parts: { horns: 5, halo: 7, crown: 8 } },
  taurus:        { body: 'spirit', main: '#F97316', parts: { horns: 5, halo: 7, crown: 8 } },
  gemini:        { body: 'spirit', main: '#22C55E', parts: { halo: 6, crown: 8 } },
  cancer:        { body: 'spirit', main: '#CBD5E1', parts: { horns: 6, halo: 7, crown: 8 } },
  leo:           { body: 'spirit', main: '#F59E0B', parts: { horns: 6, halo: 7, crown: 8 } },
  virgo:         { body: 'spirit', main: '#F8FAFC', parts: { halo: 6, crown: 8 } },
  libra:         { body: 'spirit', main: '#F472B6', parts: { halo: 6, crown: 8, wings: 10 } },
  scorpio:       { body: 'spirit', main: '#DC2626', parts: { tail: 5, halo: 7, crown: 8 } },
  sagittarius:   { body: 'spirit', main: '#8B5CF6', parts: { horns: 6, halo: 7, crown: 8 } },
  capricorn:     { body: 'spirit', main: '#3B82F6', parts: { horns: 5, tail: 5, halo: 7, crown: 8 } },
  aquarius:      { body: 'spirit', main: '#38BDF8', parts: { halo: 6, crown: 8 } },
  pisces:        { body: 'spirit', main: '#0EA5E9', parts: { tail: 4, halo: 7, crown: 8 } },
  // ===== 民间传说 =====
  nian:          { body: 'quadruped', main: '#EF4444', variant: 'nian', parts: { horns: 4, tail: 3, crown: 8 } },
  dragon_boat:   { body: 'aquatic', main: '#F97316', parts: { tail: 4, crown: 8, wings: 11 } },
  lantern:       { body: 'food', main: '#F59E0B', parts: { crown: 7, halo: 10 } },
  dumpling:      { body: 'food', main: '#F1F5F9', parts: { crown: 8, halo: 11 } },
  fu_star:       { body: 'humanoid', main: '#EF4444', parts: { crown: 6, halo: 8 } },
  shou_star:     { body: 'humanoid', main: '#F472B6', parts: { crown: 6, halo: 9, wings: 11 } },
  // ===== 扩充：神话四象/凶兽/瑞兽 =====
  qinglong: { body: 'dragon', main: '#22C55E', variant: 'antler', parts: { horns: 5, tail: 3, wings: 7 } },
  baihu: { body: 'quadruped', main: '#E2E8F0', variant: 'tiger', parts: { tail: 3, horns: 7, wings: 11 } },
  zhuque: { body: 'bird', main: '#EF4444', variant: 'phoenix', parts: { tail: 5, wings: 6, crown: 7 } },
  xuanwu: { body: 'dragon', main: '#0F766E', variant: 'xuanwu', parts: { horns: 6, tail: 3 } },
  taotie: { body: 'quadruped', main: '#8B5CF6', variant: 'taotie', parts: { horns: 6, tail: 3, crown: 9 } },
  baize: { body: 'quadruped', main: '#F5F5F4', variant: 'baize', parts: { tail: 3, horns: 8, crown: 9 } },
  // ===== 扩充：元素 =====
  ice_fox: { body: 'quadruped', main: '#38BDF8', variant: 'fox', parts: { tail: 3, horns: 9 } },
  rock_rhino: { body: 'quadruped', main: '#B45309', variant: 'rhino', parts: { horns: 5, tail: 4, crown: 9 } },
  wind_falcon: { body: 'bird', main: '#22C55E', parts: { wings: 5, tail: 5, crown: 9 } },
  light_deer: { body: 'quadruped', main: '#FDE047', variant: 'deer', parts: { horns: 5, tail: 4, halo: 8 } },
  dark_panther: { body: 'quadruped', main: '#64748B', variant: 'panther', parts: { tail: 3, horns: 8, crown: 9 } },
  steel_armadillo: { body: 'quadruped', main: '#94A3B8', variant: 'armadillo', parts: { tail: 3, horns: 8, crown: 9 } },
  // ===== 扩充：国宝 =====
  tibetan_antelope: { body: 'quadruped', main: '#F1F5F9', variant: 'antelope', parts: { horns: 5, tail: 4, crown: 9 } },
  snow_leopard: { body: 'quadruped', main: '#CBD5E1', variant: 'leopard', parts: { tail: 3, horns: 8, wings: 11 } },
  milu_deer: { body: 'quadruped', main: '#A16207', variant: 'deer', parts: { horns: 5, tail: 4, crown: 9 } },
  siberian_tiger: { body: 'quadruped', main: '#F97316', variant: 'siberian', parts: { tail: 3, horns: 8, wings: 11 } },
  red_panda: { body: 'quadruped', main: '#F97316', variant: 'redpanda', parts: { tail: 3, crown: 8 } },
  finless_porpoise: { body: 'aquatic', main: '#38BDF8', parts: { tail: 3, fin: 5 } },
  // ===== 扩充：机甲 =====
  lightsaber_warrior: { body: 'mecha', main: '#818CF8', parts: { wings: 8, horns: 6, crown: 9 } },
  fission_giant: { body: 'mecha', main: '#3B82F6', parts: { horns: 7, wings: 9 } },
  nano_swarm: { body: 'mecha', main: '#10B981', parts: { horns: 6, crown: 9 } },
  storm_jet: { body: 'mecha', main: '#60A5FA', parts: { wings: 6, horns: 6 } },
  bio_armor: { body: 'mecha', main: '#34D399', parts: { horns: 7, crown: 9 } },
  starship_core: { body: 'mecha', main: '#A78BFA', parts: { wings: 7, crown: 9 } },
  // ===== 扩充：魔法 =====
  grey_wizard: { body: 'humanoid', main: '#94A3B8', parts: { horns: 7, crown: 8, halo: 11 } },
  wand_cat: { body: 'quadruped', main: '#8B5CF6', variant: 'cat', parts: { tail: 3, horns: 8, crown: 9 } },
  dragon_knight: { body: 'humanoid', main: '#EF4444', parts: { wings: 6, horns: 7, crown: 8 } },
  alchemy_golem: { body: 'humanoid', main: '#F59E0B', parts: { horns: 7, crown: 8 } },
  nightmare_horse: { body: 'quadruped', main: '#7C3AED', variant: 'horse', parts: { tail: 4, horns: 7, wings: 11 } },
  lamp_spirit: { body: 'food', main: '#F59E0B', parts: { crown: 7, halo: 10 } },
  // ===== 扩充：史前 =====
  spinosaurus: { body: 'dragon', main: '#84CC16', variant: 'sailback', parts: { tail: 3, horns: 7, crown: 9 } },
  ankylosaurus: { body: 'dragon', main: '#78716C', variant: 'ankylo', parts: { tail: 3, horns: 7, crown: 9 } },
  diplodocus: { body: 'dragon', main: '#4D7C0F', variant: 'sauropod', parts: { tail: 3, horns: 8, crown: 9 } },
  megalodon: { body: 'aquatic', main: '#1E40AF', parts: { tail: 3, fin: 5, crown: 9 } },
  ground_sloth: { body: 'quadruped', main: '#92400E', variant: 'sloth', parts: { horns: 7, tail: 4, crown: 9 } },
  woolly_rhino: { body: 'quadruped', main: '#A16207', variant: 'woolly', parts: { horns: 5, tail: 4, crown: 9 } },
  // ===== 扩充：民间 =====
  lion_dance: { body: 'quadruped', main: '#EF4444', variant: 'lion', parts: { horns: 5, tail: 3, crown: 8 } },
  god_of_wealth: { body: 'humanoid', main: '#F59E0B', parts: { crown: 6, halo: 8 } },
  door_god: { body: 'humanoid', main: '#DC2626', parts: { horns: 7, crown: 7 } },
  kitchen_god: { body: 'humanoid', main: '#F97316', parts: { crown: 6, halo: 8 } },
  magpie: { body: 'bird', main: '#0EA5E9', parts: { wings: 5, tail: 5, crown: 8 } },
  firework_spirit: { body: 'spirit', main: '#F472B6', abstract: true, parts: { halo: 6 } },
  // ===== 传统节日(诗词典故) =====
  zongzi: { body: 'food', main: '#16A34A', parts: { crown: 8, halo: 11 } },
  tangyuan: { body: 'food', main: '#F8FAFC', parts: { crown: 8, halo: 11 } },
  mooncake: { body: 'food', main: '#B45309', parts: { crown: 8, halo: 11 } },
  qingtuan: { body: 'food', main: '#22C55E', parts: { crown: 8, halo: 11 } },
  chongyang_cake: { body: 'food', main: '#F59E0B', parts: { crown: 8, halo: 11 } },
  niangao: { body: 'food', main: '#F1F5F9', parts: { crown: 8, halo: 11 } },
  laba_porridge: { body: 'food', main: '#B45309', parts: { crown: 8, halo: 11 } },
  spring_pancake: { body: 'food', main: '#FDE68A', parts: { crown: 8, halo: 11 } },
  tanghulu: { body: 'food', main: '#EF4444', parts: { crown: 8, halo: 11 } },
  osmanthus_cake: { body: 'food', main: '#FDE047', parts: { crown: 8, halo: 11 } },
  wonton: { body: 'food', main: '#F1F5F9', parts: { crown: 8, halo: 11 } },
  festival_lantern: { body: 'food', main: '#EF4444', parts: { crown: 7, halo: 10 } },
  // ===== 七侠剑客(虹猫蓝兔七侠传，贴合原型) =====
  hongmao: { body: 'quadruped', main: '#EF4444', variant: 'cat', parts: { tail: 3, horns: 7, crown: 8 } },
  lantu: { body: 'quadruped', main: '#3B82F6', variant: 'rabbit', parts: { tail: 3, horns: 7, crown: 8 } },
  doudou: { body: 'quadruped', main: '#F8FAFC', variant: 'puppy', parts: { tail: 3, horns: 7, crown: 8 } },
  dabeng: { body: 'quadruped', main: '#92400E', variant: 'bear', parts: { tail: 3, horns: 7, crown: 8 } },
  tiaotiao: { body: 'quadruped', main: '#22C55E', variant: 'monkey', parts: { tail: 3, horns: 7, crown: 8 } },
  shali: { body: 'quadruped', main: '#8B5CF6', variant: 'kitten', parts: { tail: 3, horns: 7, crown: 8 } },
  dada: { body: 'quadruped', main: '#16A34A', variant: 'panda', parts: { tail: 3, horns: 7, crown: 8 } },
  qilin_sacred: { body: 'quadruped', main: '#F59E0B', variant: 'qilin', parts: { horns: 5, tail: 3, crown: 8 } },
  lingge: { body: 'bird', main: '#E0F2FE', variant: 'dove', parts: { wings: 5, tail: 5, crown: 8 } },
  heixinhu: { body: 'quadruped', main: '#1F2937', variant: 'tiger', parts: { tail: 3, horns: 7, crown: 8 } },
  zhuzhijie: { body: 'quadruped', main: '#F9A8D4', variant: 'kitten', parts: { tail: 3, horns: 8, crown: 8 } },
  niuxuanfeng: { body: 'quadruped', main: '#78716C', variant: 'bull', parts: { horns: 5, tail: 3, crown: 8 } },
}

/** 兜底配置(未知物种也保证能渲染) */
export const DEFAULT_ART: PetArtConfig = {
  body: 'quadruped',
  main: '#6366F1',
}

/** 星座星灵色(供 spirit 模板生成星座连线) */
export const CONSTELLATION_LINES: Record<string, number[][]> = {
  aries: [[30, 62], [48, 36], [66, 66], [84, 42], [100, 70], [122, 58]],
  taurus: [[28, 70], [46, 46], [66, 66], [84, 34], [104, 62], [126, 60]],
  gemini: [[30, 66], [50, 38], [70, 66], [90, 38], [110, 66], [128, 40]],
  leo: [[28, 58], [48, 66], [68, 38], [88, 66], [108, 40], [126, 60]],
  sagittarius: [[30, 40], [50, 66], [68, 44], [88, 68], [108, 44], [126, 66]],
}

// ============================================================
// 渲染参数计算
// ============================================================

/** 根据物种与等级计算全部渲染参数 */
export function computePetArt(speciesId: string, level: number): PetArtRender {
  const cfg = PET_ART[speciesId] || DEFAULT_ART
  const partsCfg = { ...DEFAULT_PARTS, ...(cfg.parts || {}) }
  const stage = getLevelStage(level)

  const scaleByStage: Record<string, number> = {
    egg: 0.9, baby: 1, growing: 1.12, mature: 1.26, legendary: 1.4,
  }
  const glowByStage: Record<string, number> = {
    egg: 0.25, baby: 0, growing: 0.3, mature: 0.55, legendary: 1,
  }
  const eyesByStage: Record<string, PetArtRender['eyes']> = {
    egg: 'serene', baby: 'big', growing: 'normal', mature: 'sharp', legendary: 'serene',
  }

  return {
    body: cfg.body,
    variant: cfg.variant || '',
    seriesId: getSeriesBySpeciesId(speciesId)?.id || '',
    main: cfg.main,
    accent: shade(cfg.main, 0.4),
    dark: shade(cfg.main, -0.45),
    light: shade(cfg.main, 0.62),
    scale: scaleByStage[stage],
    glow: glowByStage[stage],
    parts: {
      horns: level >= (partsCfg.horns ?? 99),
      wings: level >= (partsCfg.wings ?? 99),
      tail: level >= (partsCfg.tail ?? 99),
      fin: level >= (partsCfg.fin ?? 99),
      halo: level >= (partsCfg.halo ?? 99),
      crown: level >= (partsCfg.crown ?? 99),
    },
    stage,
    eyes: eyesByStage[stage],
    blush: stage === 'baby' || stage === 'egg',
    sparkle: level >= 5,
  }
}

/** 获取物种主色(供卡片等使用) */
export function getPetMainColor(speciesId: string): string {
  return (PET_ART[speciesId] || DEFAULT_ART).main
}
