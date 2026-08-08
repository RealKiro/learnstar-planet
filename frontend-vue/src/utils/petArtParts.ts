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
  | 'mecha'      // 机甲
  | 'spirit'     // 灵体/星灵(量子兽 星座守护 精灵 梦魇)
  | 'food'       // 物化萌物(灯笼 饺子 月饼 粽子 汤圆)
  | 'humanoid'   // 人形(精灵 树人 美人鱼 圣斗士)

/** 物种艺术配置 */
export interface PetArtConfig {
  /** 身体模板 */
  body: BodyType
  /** 主色(其余色阶由公式推导) */
  main: string
  /** 剪影变体：同身体模板下的独立外形（宝可梦/七侠/神话龙等需手绘贴合原型，跨阶段也随 variant 演化） */
  variant?: string
  /** 各阶段剪影变体（重点角色手调：随进化换更强模板） */
  stageVariants?: Partial<Record<'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent', string>>
  /** 各阶段主色覆盖（未配置时用程序化阶段明暗） */
  stageMain?: Partial<Record<'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent', string>>
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
  stage: 'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent'
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
  nine_tail_fox: { body: 'quadruped', main: '#F1F5F9', variant: 'kitsune', stageVariants: { legendary: 'kitsune-nine', transcendent: 'kitsune-nine' }, parts: { tail: 3, horns: 10 } },
  kunpeng:       { body: 'aquatic', main: '#3B82F6', variant: 'kunpeng', parts: { tail: 3, fin: 5, wings: 8 } },
  fenghuang:     { body: 'bird', main: '#EF4444', variant: 'phoenix', parts: { tail: 5, wings: 6, crown: 7 } },
  qilin:         { body: 'quadruped', main: '#10B981', variant: 'qilin', parts: { horns: 5, tail: 3, wings: 11 } },
  qiongqi:       { body: 'quadruped', main: '#374151', variant: 'qiongqi', parts: { tail: 3, horns: 6, wings: 8 } },
  bifang:        { body: 'bird', main: '#38BDF8', variant: 'bifang', parts: { wings: 5, tail: 5, crown: 9 } },
  pixiu:         { body: 'quadruped', main: '#F59E0B', variant: 'pixiu', parts: { tail: 3, horns: 6, wings: 8 } },
  jingwei:       { body: 'bird', main: '#1F2937', variant: 'jingwei', parts: { wings: 5, tail: 5, crown: 9 } },
  xiangliu:      { body: 'dragon', main: '#16A34A', variant: 'xiangliu', parts: { horns: 6, tail: 3 } },
  xiezhi:        { body: 'quadruped', main: '#F8FAFC', variant: 'xiezhi', parts: { horns: 5, tail: 3 } },
  // ===== 元素精灵(宝可梦风，贴合原型) =====
  charmander:    { body: 'quadruped', main: '#F97316', variant: 'lizard', parts: { tail: 3, horns: 7 }, stageMain: { egg: '#FCA5A5', baby: '#FB923C', growing: '#F97316', mature: '#EA580C', legendary: '#C2410C' } },
  bulbasaur:     { body: 'quadruped', main: '#0D9488', variant: 'seed', parts: { tail: 4, crown: 6 }, stageMain: { egg: '#99F6E4', baby: '#2DD4BF', growing: '#0D9488', mature: '#0F766E', legendary: '#115E59' } },
  squirtle:      { body: 'quadruped', main: '#3B82F6', variant: 'turtle', parts: { tail: 3, crown: 7 }, stageMain: { egg: '#93C5FD', baby: '#60A5FA', growing: '#3B82F6', mature: '#2563EB', legendary: '#1E40AF' } },
  eevee:         { body: 'quadruped', main: '#B45309', variant: 'fox', parts: { tail: 3, horns: 9 } },
  pikachu:       { body: 'quadruped', main: '#FACC15', variant: 'mouse', parts: { tail: 3, horns: 9 }, stageVariants: { growing: 'mouse-teen', mature: 'mouse-teen', legendary: 'mouse-legend', transcendent: 'mouse-legend' } },
  riolu:         { body: 'quadruped', main: '#6366F1', variant: 'pup', parts: { tail: 4, horns: 8 } },
  // ===== 国宝守护 =====
  panda:               { body: 'quadruped', main: '#F1F5F9', variant: 'panda', parts: { tail: 4, crown: 10 } },
  golden_monkey:       { body: 'quadruped', main: '#F59E0B', variant: 'monkey', parts: { tail: 3, crown: 9 } },
  red_crowned_crane:   { body: 'bird', main: '#F1F5F9', variant: 'crane', parts: { wings: 5, tail: 6, crown: 8 } },
  south_china_tiger:   { body: 'quadruped', main: '#F97316', variant: 'tiger', parts: { tail: 3, horns: 8, wings: 11 } },
  chinese_alligator:   { body: 'dragon', main: '#64748B', variant: 'crocodile', parts: { tail: 3, horns: 9 } },
  crested_ibis:        { body: 'bird', main: '#FBCFE8', variant: 'ibis', parts: { wings: 5, tail: 7, crown: 11 } },
  // ===== 科幻机甲 =====
  mecha_dragon:  { body: 'quadruped', main: '#F97316', variant: 'agumon', parts: { tail: 3, horns: 7, crown: 8 } },
  cyber_cat:     { body: 'quadruped', main: '#F8FAFC', variant: 'tailmon', parts: { tail: 3, horns: 7, crown: 8 } },
  space_mecha:   { body: 'quadruped', main: '#FDE68A', variant: 'patamon', parts: { tail: 4, horns: 7, wings: 8 } },
  quantum_beast: { body: 'quadruped', main: '#3B82F6', variant: 'gabumon', parts: { horns: 5, tail: 4, crown: 9 } },
  digital_phoenix: { body: 'bird', main: '#F472B6', variant: 'biyomon', parts: { wings: 5, tail: 6, crown: 9 } },
  mecha_shark:   { body: 'aquatic', main: '#E2E8F0', variant: 'gomamon', parts: { tail: 3, fin: 5, crown: 9 } },
  // ===== 魔法奇幻 =====
  unicorn:       { body: 'quadruped', main: '#F5F3FF', variant: 'unicorn', parts: { horns: 5, tail: 4, wings: 8 } },
  wyvern:        { body: 'dragon', main: '#EF4444', variant: 'wyvern', parts: { horns: 5, wings: 6, tail: 3 } },
  fairy:         { body: 'humanoid', main: '#FBBF24', variant: 'fairy', parts: { wings: 5, crown: 7 } },
  treant:        { body: 'humanoid', main: '#166534', variant: 'treant', parts: { horns: 7, crown: 10 } },
  griffin:       { body: 'quadruped', main: '#F59E0B', variant: 'griffin', parts: { wings: 6, tail: 4, crown: 8 } },
  mermaid:       { body: 'humanoid', main: '#0EA5E9', variant: 'mermaid', parts: { tail: 3, crown: 7, wings: 11 } },
  // ===== 史前生物 =====
  t_rex:         { body: 'dragon', main: '#4D7C0F', variant: 'trex', parts: { tail: 3, horns: 7, crown: 9 } },
  triceratops:   { body: 'quadruped', main: '#F59E0B', variant: 'triceratops', parts: { horns: 4, tail: 4, crown: 9 } },
  pterosaur:     { body: 'bird', main: '#60A5FA', variant: 'pterosaur', parts: { wings: 5, tail: 5 } },
  mammoth:       { body: 'quadruped', main: '#A16207', variant: 'mammoth', parts: { horns: 5, tail: 4, crown: 9 } },
  sabertooth:    { body: 'quadruped', main: '#92400E', variant: 'sabertooth', parts: { horns: 4, tail: 3, crown: 9 } },
  mosasaur:      { body: 'aquatic', main: '#0EA5E9', variant: 'mosasaur', parts: { tail: 3, fin: 5, crown: 8 } },
  // ===== 星座守护(圣斗士·圣衣) =====
  aries:         { body: 'spirit', main: '#F59E0B', variant: 'aries', parts: { horns: 5, halo: 7, crown: 8 } },
  taurus:        { body: 'spirit', main: '#F97316', variant: 'taurus', parts: { horns: 5, halo: 7, crown: 8 } },
  gemini:        { body: 'spirit', main: '#22C55E', variant: 'gemini', parts: { halo: 6, crown: 8 } },
  cancer:        { body: 'spirit', main: '#CBD5E1', variant: 'cancer', parts: { horns: 6, halo: 7, crown: 8 } },
  leo:           { body: 'spirit', main: '#F59E0B', variant: 'leo', parts: { horns: 6, halo: 7, crown: 8 } },
  virgo:         { body: 'spirit', main: '#F8FAFC', variant: 'virgo', parts: { halo: 6, crown: 8 } },
  libra:         { body: 'spirit', main: '#F472B6', variant: 'libra', parts: { halo: 6, crown: 8, wings: 10 } },
  scorpio:       { body: 'spirit', main: '#DC2626', variant: 'scorpio', parts: { tail: 5, halo: 7, crown: 8 } },
  sagittarius:   { body: 'spirit', main: '#8B5CF6', variant: 'sagittarius', parts: { horns: 6, halo: 7, crown: 8 } },
  capricorn:     { body: 'spirit', main: '#3B82F6', variant: 'capricorn', parts: { horns: 5, tail: 5, halo: 7, crown: 8 } },
  aquarius:      { body: 'spirit', main: '#38BDF8', variant: 'aquarius', parts: { halo: 6, crown: 8 } },
  pisces:        { body: 'spirit', main: '#0EA5E9', variant: 'pisces', parts: { tail: 4, halo: 7, crown: 8 } },
  // ===== 扩充：神话四象/凶兽/瑞兽 =====
  qinglong: { body: 'dragon', main: '#22C55E', variant: 'antler', parts: { horns: 5, tail: 3, wings: 7 } },
  baihu: { body: 'quadruped', main: '#E2E8F0', variant: 'tiger', parts: { tail: 3, horns: 7, wings: 11 } },
  zhuque: { body: 'bird', main: '#DC2626', variant: 'zhuque', parts: { tail: 5, wings: 6, crown: 7 } },
  xuanwu: { body: 'dragon', main: '#0F766E', variant: 'xuanwu', parts: { horns: 6, tail: 3 } },
  taotie: { body: 'quadruped', main: '#166534', variant: 'taotie', parts: { horns: 6, tail: 3, crown: 9 } },
  baize: { body: 'quadruped', main: '#F5F5F4', variant: 'baize', parts: { tail: 3, horns: 8, crown: 9 } },
  // ===== 扩充：元素 =====
  ice_fox: { body: 'quadruped', main: '#38BDF8', variant: 'fox', parts: { tail: 3, horns: 9 } },
  rock_rhino: { body: 'quadruped', main: '#B45309', variant: 'rhino', parts: { horns: 5, tail: 4, crown: 9 } },
  wind_falcon: { body: 'bird', main: '#22C55E', variant: 'falcon', parts: { wings: 5, tail: 5, crown: 9 } },
  light_deer: { body: 'quadruped', main: '#FDE047', variant: 'deer', parts: { horns: 5, tail: 4, halo: 8 } },
  dark_panther: { body: 'quadruped', main: '#64748B', variant: 'panther', parts: { tail: 3, horns: 8, crown: 9 } },
  steel_armadillo: { body: 'quadruped', main: '#94A3B8', variant: 'armadillo', parts: { tail: 3, horns: 8, crown: 9 } },
  // ===== 扩充：国宝 =====
  tibetan_antelope: { body: 'quadruped', main: '#F1F5F9', variant: 'antelope', parts: { horns: 5, tail: 4, crown: 9 } },
  snow_leopard: { body: 'quadruped', main: '#CBD5E1', variant: 'leopard', parts: { tail: 3, horns: 8, wings: 11 } },
  milu_deer: { body: 'quadruped', main: '#A16207', variant: 'deer', parts: { horns: 5, tail: 4, crown: 9 } },
  siberian_tiger: { body: 'quadruped', main: '#F97316', variant: 'siberian', parts: { tail: 3, horns: 8, wings: 11 } },
  red_panda: { body: 'quadruped', main: '#F97316', variant: 'redpanda', parts: { tail: 3, crown: 8 } },
  finless_porpoise: { body: 'aquatic', main: '#38BDF8', variant: 'porpoise', parts: { tail: 3, fin: 5 } },
  // ===== 扩充：魔法 =====
  grey_wizard: { body: 'humanoid', main: '#94A3B8', variant: 'wizard', parts: { horns: 7, crown: 8, halo: 11 } },
  wand_cat: { body: 'quadruped', main: '#8B5CF6', variant: 'cat', parts: { tail: 3, horns: 8, crown: 9 } },
  dragon_knight: { body: 'humanoid', main: '#EF4444', variant: 'knight', parts: { wings: 6, horns: 7, crown: 8 } },
  alchemy_golem: { body: 'humanoid', main: '#F59E0B', variant: 'golem', parts: { horns: 7, crown: 8 } },
  nightmare_horse: { body: 'quadruped', main: '#7C3AED', variant: 'horse', parts: { tail: 4, horns: 7, wings: 11 } },
  lamp_spirit: { body: 'food', main: '#F59E0B', variant: 'lamp', parts: { crown: 7, halo: 10 } },
  // ===== 扩充：史前 =====
  spinosaurus: { body: 'dragon', main: '#D97706', variant: 'sailback', parts: { tail: 3, horns: 7, crown: 9 } },
  ankylosaurus: { body: 'dragon', main: '#78716C', variant: 'ankylo', parts: { tail: 3, horns: 7, crown: 9 } },
  diplodocus: { body: 'dragon', main: '#4D7C0F', variant: 'sauropod', parts: { tail: 3, horns: 8, crown: 9 } },
  megalodon: { body: 'aquatic', main: '#1E40AF', variant: 'shark', parts: { tail: 3, fin: 5, crown: 9 } },
  ground_sloth: { body: 'quadruped', main: '#92400E', variant: 'sloth', parts: { horns: 7, tail: 4, crown: 9 } },
  woolly_rhino: { body: 'quadruped', main: '#A16207', variant: 'woolly', parts: { horns: 5, tail: 4, crown: 9 } },
  // ===== 传统节日(诗词典故) =====
  zongzi: { body: 'food', main: '#16A34A', variant: 'zongzi', parts: { crown: 8, halo: 11 } },
  tangyuan: { body: 'food', main: '#F8FAFC', variant: 'tangyuan', parts: { crown: 8, halo: 11 } },
  mooncake: { body: 'food', main: '#B45309', variant: 'mooncake', parts: { crown: 8, halo: 11 } },
  qingtuan: { body: 'food', main: '#22C55E', variant: 'qingtuan', parts: { crown: 8, halo: 11 } },
  chongyang_cake: { body: 'food', main: '#F59E0B', variant: 'layered', parts: { crown: 8, halo: 11 } },
  niangao: { body: 'food', main: '#F1F5F9', variant: 'niangao', parts: { crown: 8, halo: 11 } },
  laba_porridge: { body: 'food', main: '#B45309', variant: 'porridge', parts: { crown: 8, halo: 11 } },
  spring_pancake: { body: 'food', main: '#FDE68A', variant: 'roll', parts: { crown: 8, halo: 11 } },
  tanghulu: { body: 'food', main: '#EF4444', variant: 'tanghulu', parts: { crown: 8, halo: 11 } },
  osmanthus_cake: { body: 'food', main: '#FDE047', variant: 'osmanthus_cake', parts: { crown: 8, halo: 11 } },
  wonton: { body: 'food', main: '#F1F5F9', variant: 'wonton', parts: { crown: 8, halo: 11 } },
  festival_lantern: { body: 'food', main: '#EF4444', variant: 'hexlantern', parts: { crown: 7, halo: 10 } },
  // ===== 虹猫蓝兔七侠传(虹猫蓝兔七侠传，贴合原型) =====
  hongmao: { body: 'quadruped', main: '#EF4444', variant: 'cat', parts: { tail: 3, horns: 7, crown: 8 }, stageMain: { egg: '#FECACA', baby: '#F87171', growing: '#EF4444', mature: '#DC2626', legendary: '#991B1B' } },
  lantu: { body: 'quadruped', main: '#3B82F6', variant: 'rabbit', parts: { tail: 3, horns: 7, crown: 8 }, stageMain: { egg: '#BFDBFE', baby: '#93C5FD', growing: '#3B82F6', mature: '#2563EB', legendary: '#1D4ED8' } },
  doudou: { body: 'quadruped', main: '#F8FAFC', variant: 'puppy', parts: { tail: 3, horns: 7, crown: 8 } },
  dabeng: { body: 'quadruped', main: '#92400E', variant: 'bear', parts: { tail: 3, horns: 7, crown: 8 } },
  tiaotiao: { body: 'quadruped', main: '#22C55E', variant: 'monkey', parts: { tail: 3, horns: 7, crown: 8 } },
  shali: { body: 'quadruped', main: '#8B5CF6', variant: 'kitten', parts: { tail: 3, horns: 7, crown: 8 } },
  dada: { body: 'quadruped', main: '#16A34A', variant: 'panda', parts: { tail: 3, horns: 7, crown: 8 } },
  heixinhu: { body: 'quadruped', main: '#1F2937', variant: 'tiger', parts: { tail: 3, horns: 7, crown: 8 }, stageMain: { egg: '#9CA3AF', baby: '#6B7280', growing: '#1F2937', mature: '#111827', legendary: '#0B0F19' } },
  heixiaohu: { body: 'quadruped', main: '#1F2937', variant: 'tiger', parts: { horns: 6, tail: 3, crown: 8, wings: 11 }, stageMain: { egg: '#9CA3AF', baby: '#6B7280', growing: '#1F2937', mature: '#111827', legendary: '#0B0F19' } },
  // ===== 东方神话(神话) =====
  jiang_ziya:    { body: 'humanoid', main: '#F59E0B', variant: 'jiangziya', parts: { crown: 7, halo: 11 } },
  nezha:         { body: 'humanoid', main: '#EF4444', variant: 'knight', stageVariants: { legendary: 'nezha-tri', transcendent: 'nezha-tri' }, parts: { crown: 8, wings: 11 } },
  yang_jian:     { body: 'humanoid', main: '#3B82F6', variant: 'knight', parts: { crown: 8 } },
  lei_zhenzi:    { body: 'bird', main: '#F59E0B', variant: 'falcon', parts: { wings: 5, tail: 5 } },
  huang_tianhua: { body: 'humanoid', main: '#F97316', variant: 'knight', parts: { crown: 8 } },
  tu_xingsun:    { body: 'humanoid', main: '#A16207', variant: 'golem', parts: { crown: 9 } },
  yang_ren:      { body: 'humanoid', main: '#22C55E', variant: 'yangren', parts: { halo: 10, crown: 8 } },
  wei_hu:        { body: 'humanoid', main: '#6366F1', variant: 'giant', parts: { crown: 8, halo: 11 } },
  daji:          { body: 'quadruped', main: '#EC4899', variant: 'kitsune', parts: { tail: 3, horns: 10 } },
  shen_gongbao:  { body: 'humanoid', main: '#7C3AED', variant: 'shen', parts: { crown: 7 } },
  // ===== 东方神话·补录 =====
  sun_wukong:      { body: 'quadruped', main: '#F59E0B', variant: 'monkey', stageVariants: { legendary: 'wukong-great', transcendent: 'wukong-great' }, parts: { tail: 3, crown: 9, horns: 6 } },
  lv_dongbin:      { body: 'humanoid', main: '#F8FAFC', variant: 'lvdongbin', parts: { crown: 8, halo: 11 } },
  he_xiangu:       { body: 'humanoid', main: '#F472B6', variant: 'fairy', parts: { wings: 7, crown: 7 } },
  zhang_guolao:    { body: 'quadruped', main: '#94A3B8', variant: 'horse', parts: { tail: 3, crown: 9 } },
  tie_guaili:      { body: 'humanoid', main: '#57534E', variant: 'golem', parts: { crown: 8 } },
  han_zhongli:     { body: 'humanoid', main: '#EF4444', variant: 'giant', parts: { crown: 8, halo: 11 } },
  lan_caihe:       { body: 'humanoid', main: '#22C55E', variant: 'fairy', parts: { crown: 7 } },
  cao_guojiu:      { body: 'humanoid', main: '#FBBF24', variant: 'caoguojiu', parts: { crown: 8 } },
  taishang_laojun: { body: 'humanoid', main: '#7C3AED', variant: 'laojun', parts: { halo: 10, crown: 7 } },
  zhong_kui:       { body: 'humanoid', main: '#DC2626', variant: 'giant', parts: { crown: 8, halo: 11 } },
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
    egg: 0.9, baby: 1, growing: 1.12, mature: 1.26, legendary: 1.4, transcendent: 1.55,
  }
  const glowByStage: Record<string, number> = {
    egg: 0.25, baby: 0, growing: 0.3, mature: 0.55, legendary: 1, transcendent: 1,
  }
  const eyesByStage: Record<string, PetArtRender['eyes']> = {
    egg: 'serene', baby: 'big', growing: 'normal', mature: 'sharp', legendary: 'serene', transcendent: 'serene',
  }

  // 程序化阶段主色：幼年提亮、成熟加深、传说提亮、道果提亮（重点角色可用 stageMain 覆盖）
  const stageShade: Record<string, number> = {
    egg: 0.45, baby: 0.28, growing: 0.05, mature: -0.18, legendary: 0.12, transcendent: 0.2,
  }
  const baseMain = cfg.stageMain?.[stage] || cfg.main
  const main = stageShade[stage] ? shade(baseMain, stageShade[stage]) : baseMain

  // 阶段部件开放：幼年只保留基础尾，成长加角，成熟加翅膀/鳍，传说加光环/王冠
  const stageAllowed: Record<string, { horns: boolean; wings: boolean; fin: boolean; halo: boolean; crown: boolean }> = {
    egg:      { horns: false, wings: false, fin: false, halo: false, crown: false },
    baby:     { horns: false, wings: false, fin: false, halo: false, crown: false },
    growing:  { horns: true,  wings: false, fin: false, halo: false, crown: false },
    mature:   { horns: true,  wings: true,  fin: true,  halo: false, crown: false },
    legendary:{ horns: true,  wings: true,  fin: true,  halo: true,  crown: true },
    transcendent:{ horns: true, wings: true, fin: true, halo: true, crown: true },
  }
  const allowed = stageAllowed[stage] || stageAllowed.growing

  return {
    body: cfg.body,
    variant: cfg.stageVariants?.[stage] || cfg.variant || '',
    seriesId: getSeriesBySpeciesId(speciesId)?.id || '',
    main,
    accent: shade(main, 0.4),
    dark: shade(main, -0.45),
    light: shade(main, 0.62),
    scale: scaleByStage[stage],
    glow: glowByStage[stage],
    parts: {
      horns: allowed.horns && level >= (partsCfg.horns ?? 99),
      wings: allowed.wings && level >= (partsCfg.wings ?? 99),
      tail: level >= (partsCfg.tail ?? 99),
      fin: allowed.fin && level >= (partsCfg.fin ?? 99),
      halo: allowed.halo && level >= (partsCfg.halo ?? 99),
      crown: allowed.crown && level >= (partsCfg.crown ?? 99),
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
