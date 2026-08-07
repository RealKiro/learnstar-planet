// ===== 学趣星球 · 宠物数据系统 =====
// 10 个系列、120 个物种、每个物种 12 级进化
// 基础 8 系列数据来源：宠物系列设计.txt；扩充物种见 petDataExtended.ts

import { EXTRA_SPECIES, EXTRA_SERIES } from './petDataExtended'

/** 系列定义 */
export interface PetSeries {
  id: string
  name: string
  emoji: string
  species: PetSpecies[]
}

/** 物种定义 */
export interface PetSpecies {
  id: string
  name: string
  seriesId: string
  levels: PetLevel[]
}

/** 单级进化 */
export interface PetLevel {
  level: number       // 1-12
  name: string
  description: string
  /** 等级阶段分类(6阶段): 1-2卵生, 3-4幼年, 5-6成长, 7-8成熟, 9-10传说, 11-12归真 */
  stage: 'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent'
  /** 所需累计积分 */
  requiredScore: number
}

/** 获取等级所需累计积分(与各物种 requiredScore 一致的统一阈值表 Lv.0~12) */
export function getLevelRequiredScore(level: number): number {
  const scores = [0, 0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450]
  return scores[level] ?? 450
}

/** 判断等级阶段(Lv.1-2 卵生 / 3-5 幼年 / 6-8 成长 / 9-10 成熟 / 11-12 传说) */
export function getLevelStage(level: number): PetLevel['stage'] {
  if (level <= 2) return 'egg'
  if (level <= 4) return 'baby'
  if (level <= 6) return 'growing'
  if (level <= 8) return 'mature'
  if (level <= 10) return 'legendary'
  return 'transcendent'
}

/** 获取等级的通用阶段称谓 */
export function getLevelTitle(level: number): string {
  const titles: Record<number, string> = {
    1: '新生之卵',
    2: '破壳萌芽',
    3: '蹒跚学步',
    4: '稚气未脱',
    5: '初露锋芒',
    6: '茁壮成长',
    7: '锋芒毕露',
    8: '英姿勃发',
    9: '成熟威仪',
    10: '辉煌巅峰',
    11: '超凡入圣',
    12: '不朽传奇',
  }
  return titles[level] || '未知'
}

/** 场景背景配置 */
export interface SceneConfig {
  /** 主题色 */
  primaryColor: string
  /** 背景渐变 */
  bgGradient: string
  /** 装饰元素 emoji */
  decor: string[]
  /** 场景描述 */
  scene: string
}

/** 系列对应的场景配置 */
export const SERIES_SCENES: Record<string, SceneConfig> = {
  myth: {
    primaryColor: '#F59E0B',
    bgGradient: 'linear-gradient(180deg, #0F0A2E 0%, #1B1255 40%, #2D1B69 70%, #1A0F3E 100%)',
    decor: ['⛰️', '☁️', '✨', '🌙'],
    scene: '悬浮仙山 · 云海缭绕',
  },
  pokemon: {
    primaryColor: '#EF4444',
    bgGradient: 'linear-gradient(180deg, #1a1a2e 0%, #16213e 40%, #0f3460 70%, #1a1a2e 100%)',
    decor: ['⚡', '🔥', '💧', '🌿'],
    scene: '战斗平原 · 元素交织',
  },
  national: {
    primaryColor: '#10B981',
    bgGradient: 'linear-gradient(180deg, #064E3B 0%, #065F46 30%, #047857 60%, #064E3B 100%)',
    decor: ['🎋', '🌸', '🍃', '🏮'],
    scene: '竹林深处 · 国风雅韵',
  },
  digimon: {
    primaryColor: '#3B82F6',
    bgGradient: 'linear-gradient(180deg, #020617 0%, #0F172A 30%, #1E293B 60%, #020617 100%)',
    decor: ['💠', '⚙️', '📡', '🔧'],
    scene: '赛博都市 · 数据流光',
  },
  magic: {
    primaryColor: '#8B5CF6',
    bgGradient: 'linear-gradient(180deg, #1a0033 0%, #2d1b69 30%, #4a1a6b 60%, #1a0033 100%)',
    decor: ['🦄', '✨', '🌟', '🔮'],
    scene: '魔法森林 · 星尘飞舞',
  },
  prehistoric: {
    primaryColor: '#D97706',
    bgGradient: 'linear-gradient(180deg, #1a1a0a 0%, #2a2a1a 30%, #3a3a2a 60%, #1a1a0a 100%)',
    decor: ['🌋', '🦴', '🌿', '🪨'],
    scene: '远古丛林 · 火山喷薄',
  },
  constellation: {
    primaryColor: '#6366F1',
    bgGradient: 'linear-gradient(180deg, #000814 0%, #001233 30%, #001845 60%, #000814 100%)',
    decor: ['⭐', '🌌', '✨', '🌠'],
    scene: '星空深处 · 银河旋转',
  },
  festival: {
    primaryColor: '#EF4444',
    bgGradient: 'linear-gradient(180deg, #2a0a1a 0%, #4a1026 30%, #6b1a2a 60%, #2a0a1a 100%)',
    decor: ['🎊', '🏮', '🎆', '🥮'],
    scene: '节庆长街 · 灯火与诗',
  },
  qixia: {
    primaryColor: '#EF4444',
    bgGradient: 'linear-gradient(180deg, #1a0a0a 0%, #3a1210 30%, #5a1a0a 60%, #1a0a0a 100%)',
    decor: ['⚔️', '🐉', '🌲', '🏯'],
    scene: '七剑峡谷 · 侠影重重',
  },
  dongfang: {
    primaryColor: '#D97706',
    bgGradient: '线性渐变(古战场·封神台·祥云)深色调',
    decor: ['🗡️', '📜', '☯️', '🐒'],
    scene: '封神台 · 仙神斗法',
  },
}

/** 系列与物种的Emoji映射 */
export const SPECIES_EMOJI: Record<string, string> = {
  // 神话
  zhulong: '🐉',
  yinglong: '🐉',
  nine_tail_fox: '🦊',
  kunpeng: '🐋',
  fenghuang: '🦅',
  qilin: '🦄',
  qiongqi: '🐯',
  bifang: '🦚',
  pixiu: '🦁',
  jingwei: '🕊️',
  xiangliu: '🐍',
  xiezhi: '🐐',
  // 宝可梦风
  charmander: '🦎',
  bulbasaur: '🐸',
  squirtle: '🐢',
  eevee: '🦊',
  pikachu: '🐹',
  riolu: '🐕',
  // 国宝
  panda: '🐼',
  golden_monkey: '🐒',
  red_crowned_crane: '🦩',
  south_china_tiger: '🐯',
  chinese_alligator: '🐊',
  crested_ibis: '🕊️',
  // 机甲
  mecha_dragon: '🤖',
  cyber_cat: '🐱',
  space_mecha: '🦾',
  quantum_beast: '👾',
  digital_phoenix: '🦅',
  mecha_shark: '🦈',
  // 魔法
  unicorn: '🦄',
  wyvern: '🐉',
  fairy: '🧚',
  treant: '🌳',
  griffin: '🦅',
  mermaid: '🧜‍♀️',
  // 史前
  t_rex: '🦖',
  triceratops: '🦕',
  pterosaur: '🦅',
  mammoth: '🐘',
  sabertooth: '🐯',
  mosasaur: '🐊',
  // 星座
  aries: '🐏',
  taurus: '🐂',
  gemini: '👥',
  cancer: '🦀',
  leo: '🦁',
  virgo: '👩',
  libra: '⚖️',
  scorpio: '🦂',
  sagittarius: '🏹',
  capricorn: '🐐',
  aquarius: '🏺',
  pisces: '🐟',
  // 神话·扩充
  qinglong: '🐉',
  baihu: '🐯',
  zhuque: '🐦',
  xuanwu: '🐢',
  taotie: '😋',
  baize: '📖',
  // 元素·扩充
  ice_fox: '🦊',
  rock_rhino: '🦏',
  wind_falcon: '🦅',
  light_deer: '🦌',
  dark_panther: '🐈‍⬛',
  steel_armadillo: '🦔',
  // 国宝·扩充
  tibetan_antelope: '🐐',
  snow_leopard: '🐆',
  milu_deer: '🦌',
  siberian_tiger: '🐯',
  red_panda: '🐼',
  finless_porpoise: '🐬',
  // 魔法·扩充
  grey_wizard: '🧙',
  wand_cat: '🐱',
  dragon_knight: '🐲',
  alchemy_golem: '🗿',
  nightmare_horse: '🐴',
  lamp_spirit: '🧞',
  // 史前·扩充
  spinosaurus: '🦖',
  ankylosaurus: '🦕',
  diplodocus: '🦕',
  megalodon: '🦈',
  ground_sloth: '🦥',
  woolly_rhino: '🦏',
  // 传统节日
  zongzi: '🫔',
  tangyuan: '🍡',
  mooncake: '🥮',
  qingtuan: '🍡',
  chongyang_cake: '🍰',
  niangao: '🍥',
  laba_porridge: '🥣',
  spring_pancake: '🥞',
  tanghulu: '🍢',
  osmanthus_cake: '🍮',
  wonton: '🥟',
  festival_lantern: '🏮',
  // 虹猫蓝兔七侠传
  hongmao: '🐱',
  lantu: '🐰',
  doudou: '🐶',
  dabeng: '🐻',
  tiaotiao: '🐒',
  shali: '🦊',
  dada: '🐼',
  heixinhu: '🐯',
  heixiaohu: '🐯',
  // 东方神话
  jiang_ziya: '📜',
  nezha: '🔥',
  yang_jian: '⚔️',
  lei_zhenzi: '⚡',
  huang_tianhua: '🗡️',
  tu_xingsun: '⛏️',
  yang_ren: '👁️',
  wei_hu: '🔱',
  daji: '🦊',
  shen_gongbao: '🐆',
  sun_wukong: '🐒',
  lv_dongbin: '⚔️',
  he_xiangu: '🌸',
  zhang_guolao: '🐴',
  tie_guaili: '🩼',
  han_zhongli: '🪭',
  lan_caihe: '🎵',
  cao_guojiu: '🪇',
  taishang_laojun: '☯️',
  zhong_kui: '👺',
}

/** 获取物种Emoji */
export function getSpeciesEmoji(speciesId: string): string {
  return SPECIES_EMOJI[speciesId] || '🌟'
}

/** 获取系列Emoji */
export function getSeriesEmoji(seriesId: string): string {
  const map: Record<string, string> = {
    myth: '🏔️',
    pokemon: '⚡',
    national: '🐼',
    mecha: '🤖',
    magic: '🦄',
    prehistoric: '🦕',
    constellation: '♈',
    festival: '🎊',
    qixia: '⚔️',
    dongfang: '🗡️',
  }
  return map[seriesId] || '🌟'
}

/** 获取系列名 */
export function getSeriesName(seriesId: string): string {
  const map: Record<string, string> = {
    myth: '山海经',
    pokemon: '元素精灵',
    national: '国宝守护',
    mecha: '数码宝贝',
    magic: '魔法奇幻',
    prehistoric: '史前生物',
    constellation: '星座守护',
    festival: '传统节日',
    qixia: '虹猫蓝兔七侠传',
    dongfang: '东方神话',
  }
  return map[seriesId] || '未知'
}

// ============================================================
// 完整宠物数据：8 个系列，54 个物种
// ============================================================

export const PET_SERIES: PetSeries[] = [
  // ===== 一、山海经系列（12种） =====
  {
    id: 'myth', name: '山海经', emoji: '🏔️',
    species: [
      {
        id: 'zhulong', name: '烛龙', seriesId: 'myth',
        levels: [
          { level: 1, name: '烛火微光', description: '一团跳动的橙红色火焰，中心有竖瞳', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '烛火之影', description: '火焰中浮现细长蛇形轮廓，忽明忽暗', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '衔烛之苗', description: '蛇身人面，口含小火苗，双目紧闭', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '烛阴初醒', description: '身体拉长，鳞片隐约可见，睁一只眼', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '衔烛而行', description: '口中蜡烛燃起，照亮周围半米', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '烛照四方', description: '身体环绕光晕，可照亮暗处，鳞片发金', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '昼夜使者', description: '左眼为日（金），右眼为月（银），交替闪烁', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '衔烛之龙', description: '完整龙形，口衔巨大蜡烛，威严初现', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '烛龙巡天', description: '周身火焰缭绕，尾部带光尾，云气跟随', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '烛照九幽', description: '光芒穿透九层幽暗，神圣不可侵，身体半透明', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '开目为昼', description: '睁眼时全屏变亮，闭目则暗，时空微动', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '烛龙至尊', description: '日月环绕，呼吸间引动星辰，周身有阴阳二气', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'yinglong', name: '应龙', seriesId: 'myth',
        levels: [
          { level: 1, name: '虺卵', description: '青色蛋壳，表面有云纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳虺', description: '无足小蛇，头有肉角', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '虺行', description: '能蠕动前行，鳞片初生', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '幼蛟', description: '身体变粗，有鳞无角，可腾跃', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '蛟龙', description: '头顶生短角，身带水汽', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '角龙', description: '角长，有四足，能控水', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '应龙雏形', description: '背部生出双翼（小），主风雨', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '应龙', description: '双翼展开，身披金鳞，能飞', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '战龙', description: '翼展增大，爪牙锋利，身带雷光', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '斩蚩尤之姿', description: '黄金战甲，翼展遮天，双目赤金', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '神龙降世', description: '身缠雷电，翼下生风，云海随行', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '太古应龙', description: '身体巨大化，鳞片呈古铜色，一吼震天', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'nine_tail_fox', name: '九尾狐', seriesId: 'myth',
        levels: [
          { level: 1, name: '灵狐蛋', description: '白色蛋壳，有九条浅纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼狐', description: '通体雪白，一尾，耳尖', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '白狐', description: '一尾变长，能小跑', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '二尾狐', description: '尾部分叉为二，额生灵纹', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '三尾灵狐', description: '三尾，灵纹发光，能迷惑', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '四尾', description: '体型增大，毛色泛银', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '五尾', description: '尾巴增多，身绕狐火（蓝）', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '六尾妖狐', description: '六尾，狐火旺盛，眼神魅惑', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '七尾', description: '尾巴蓬松，额生第三眼（闭）', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '八尾', description: '八尾，第三眼微睁，散发圣光', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '九尾天狐', description: '九尾全开，周身神圣，不可侵犯', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '青丘之主', description: '九尾化作九道流光，身体半透明，周围有青丘幻影', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'kunpeng', name: '鲲鹏', seriesId: 'myth',
        levels: [
          { level: 1, name: '鲲苗', description: '小鱼形态，北冥之鳞（深蓝）', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鲲', description: '体型变大，背鳍突出', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '巨鲲', description: '遮天巨鱼，浪涛随身，喷水柱', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '灵鲲', description: '鳍变长，头部生角', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '化鲲', description: '鱼身开始变扁，鳍变为羽状', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '半鹏', description: '鱼身鸟翼，正在蜕变，羽翼未丰', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '鹏雏', description: '双翼成形，能滑翔', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '大鹏', description: '翼若垂天之云，可翱翔', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '鲲鹏', description: '鱼鸟合一，可游可飞', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '鹏程万里', description: '翼展遮天，周身有风云', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '北冥之主', description: '身体环绕极光，游于星河', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '逍遥鲲鹏', description: '化作水汽与云，无形无相，随时凝聚', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'fenghuang', name: '凤凰', seriesId: 'myth',
        levels: [
          { level: 1, name: '火羽蛋', description: '红色蛋壳，有火焰羽毛纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏凤', description: '黄色绒毛，头有冠凸', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '幼凤', description: '长出红色羽毛，尾羽短', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '青鸾', description: '羽毛变青，尾羽变长', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '赤凤', description: '转为赤红色，尾羽有金斑', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '朱雀', description: '全身火红，周围有火星', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '凤雏', description: '尾羽变长，冠羽立起', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '凤凰', description: '七彩羽毛，尾羽如虹', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '神凤', description: '周身火焰环绕，鸣叫悦耳', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '百鸟之王', description: '尾羽展开如屏，光芒四射', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '涅槃之凤', description: '身体半透明，火焰中带重生之意', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒凤凰', description: '周身金光，每一次展翅撒下星火，不死不灭', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'qilin', name: '麒麟', seriesId: 'myth',
        levels: [
          { level: 1, name: '瑞兽玉', description: '碧绿色玉石，内有祥云', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '玉麒麟幼', description: '玉石裂开，露出小兽头，有角', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '幼麟', description: '全身青绿，四蹄踏云', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '角麟', description: '角变长，身披鳞甲', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '仁兽', description: '眼神温和，身边有瑞气', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '祥瑞之麟', description: '足下生莲，头顶有光圈', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '踏火麒麟', description: '四蹄带火，鳞片泛金', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '圣麟', description: '身体变白，角呈金色', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '天麟', description: '身周有云霞，能腾云', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '太平麒麟', description: '脚踏祥云，口衔玉书', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '万瑞之祖', description: '全身七彩，角上刻有福纹', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒瑞兽', description: '化作巨大玉雕，护佑一方，福气源源', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'qiongqi', name: '穷奇', seriesId: 'myth',
        levels: [
          { level: 1, name: '凶兽卵', description: '灰黑卵壳，隐有翼纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼崽', description: '虎形小兽，初生双翼', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '虎翼兽', description: '背上生翼，吼声渐凶', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '穷奇雏形', description: '周身生刺，目露凶光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '翼虎', description: '翼展成形，能扑食飞禽', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '凶兽穷奇', description: '状如虎而有翼，可吞食恶人', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '穷奇·残暴', description: '爪牙更利，嗜食是非', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '凶神', description: '见善则吞，见恶则饲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '穷奇·噬恶', description: '遍体生鳞，翼如利刃', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '太古凶兽', description: '吞天噬地，恶气缠绕', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '四凶之穷奇', description: '威震山海，群兽辟易', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '凶兽之极', description: '化作凶煞之气，无形无相', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'bifang', name: '毕方', seriesId: 'myth',
        levels: [
          { level: 1, name: '毕方卵', description: '青色卵壳，羽纹如火', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏鸟', description: '一足初现，青色绒毛', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '一足鸟', description: '独脚而立，白喙微张', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '毕方雏形', description: '羽毛转青，身绕火星', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '青鸟', description: '见则火起，声如鹤鸣', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '毕方', description: '青色单足鸟，白喙赤足', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '毕方·衔火', description: '喙中衔火，所过之处讹火生', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '火鸟毕方', description: '周身火光，翼展如炬', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '毕方·焚风', description: '火焰随翼卷起', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '火灾之兆', description: '见则其邑有讹火', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '毕方神鸟', description: '火神之侍，驭火不伤己', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '毕方至尊', description: '化作赤色火凰，燃尽旧世', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'pixiu', name: '貔貅', seriesId: 'myth',
        levels: [
          { level: 1, name: '貔貅玉', description: '玉色卵壳，祥云环绕', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼貔貅', description: '头大身小，圆润可爱', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '貔貅雏形', description: '背生小翼，口衔元宝', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '招财貔貅', description: '只进不出，聚财纳福', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '貔貅·纳财', description: '周身金光，宝气环绕', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '瑞兽貔貅', description: '翼展如金，威而不凶', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '貔貅·镇宅', description: '坐镇一方，辟邪纳福', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '金貔貅', description: '通体鎏金，双目放光', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '貔貅·吞金', description: '吞金吐玉，福泽深厚', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '招财神兽', description: '脚踏元宝，财源滚滚', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '貔貅至尊', description: '金色祥光，护佑八方', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '瑞兽之极', description: '化金为雨，福泽天下', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'jingwei', name: '精卫', seriesId: 'myth',
        levels: [
          { level: 1, name: '精卫卵', description: '青色卵壳，带点点泪痕', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏鸟', description: '斑驳羽毛，赤足初现', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '女娃鸟', description: '常衔西山之木石', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '精卫雏形', description: '赤首白喙，啼声哀婉', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '衔木鸟', description: '日衔木石，飞往东海', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '精卫', description: '形如乌，文首白喙赤足', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '精卫·填海', description: '昼夜不息，衔石填海', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '填海之志', description: '大海无边，其志不改', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '精卫·不屈', description: '每一声啼鸣都是誓言', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '不屈之魂', description: '沧海横流，精魂不灭', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '精卫神鸟', description: '化为一缕不屈之光', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '精卫永恒', description: '衔石填海，永不止息', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'xiangliu', name: '相柳', seriesId: 'myth',
        levels: [
          { level: 1, name: '九首卵', description: '灰色卵壳，九道暗纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼蛇', description: '单首小蛇，毒气微生', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '双首', description: '首部分叉为二，剧毒', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '三首', description: '三首同鸣，声如婴儿', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '多首蛇', description: '首数渐增，毒液横流', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '相柳雏形', description: '九首渐成，自环其尾', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '相柳', description: '九首蛇身，食于九土', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '相柳·大泽', description: '所过之处皆为泽水', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '凶兽相柳', description: '九首齐舞，毒雾蔽天', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '共工之臣', description: '助共工为虐，洪水滔天', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '相柳·湮世', description: '毒水漫山，草木皆枯', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '大禹诛之', description: '其血腥臭，不可生谷', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'xiezhi', name: '獬豸', seriesId: 'myth',
        levels: [
          { level: 1, name: '獬豸卵', description: '白色卵壳，独角纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼羊', description: '羊形小兽，额有小角', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '独角羊', description: '独角渐长，目光清正', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '獬豸雏形', description: '能辨忠奸，触不直者', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '神羊', description: '见争斗则触其不直者', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '獬豸', description: '独角神羊，明辨曲直', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '獬豸·司法', description: '立于殿堂，触奸佞', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '正义之兽', description: '身绕祥光，公正无偏', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '獬豸·执法', description: '触其角，罪人胆寒', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '明辨是非', description: '一目了然，绝不冤枉', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '司法神兽', description: '公正之名，传遍山海', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '獬豸至尊', description: '法眼如炬，天理昭彰', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 二、宝可梦系列（6种） =====
  {
    id: 'pokemon', name: '宝可梦', emoji: '⚡',
    species: [
      {
        id: 'charmander', name: '小火龙→喷火龙', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '火纹蛋', description: '红色蛋壳，火焰纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳小龙', description: '露出头，尾尖有火星', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '蹒跚火苗', description: '站立，尾火微弱', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '小火龙', description: '橙红皮肤，尾火燃烧', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '烈火小龙', description: '尾火变大，能吐小火花', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '火恐龙', description: '体型增大，头顶长小角', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '烈火恐龙', description: '角大，爪牙锋利，尾火旺盛', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '喷火龙雏形', description: '背部隆起，翅膀初现', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '喷火龙', description: '双翼展开，能飞行，尾火冲天', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '烈焰喷火龙', description: '全身火焰包围，眼神锐利', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '超级喷火龙X', description: '蓝焰形态，全身暗蓝，霸气', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '始源喷火龙', description: '金色与蓝色交织，翼展遮天，全屏火焰', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'bulbasaur', name: '妙蛙种子→妙蛙花', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '种子蛋', description: '绿色蛋壳，有藤蔓纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳芽', description: '露出小头，背有芽', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '幼种子', description: '四足爬行，背芽变叶', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '妙蛙种子', description: '背生种子，蹲坐', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '叶芽成长', description: '种子变大，长出藤蔓', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '妙蛙草', description: '种子开花（小花），站立', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '花苞绽放', description: '花朵变大，香气扑鼻', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '妙蛙花', description: '背上巨花盛开，藤蔓环绕', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '繁花之姿', description: '花朵散发光芒，藤蔓延伸', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '超级妙蛙花', description: '巨花盛开，藤蔓如手臂', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '花神', description: '周身花瓣飘落，大地回春', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '远古妙蛙', description: '身体巨大化，花中藏有森林', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'squirtle', name: '杰尼龟→水箭龟', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '水纹蛋', description: '蓝色蛋壳，波浪纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳小龟', description: '露出头，背有小壳', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '幼龟', description: '能爬行，壳软', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '杰尼龟', description: '龟壳坚硬，能喷水', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '水枪龟', description: '喷水变强，尾巴变长', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '卡咪龟', description: '体型增大，壳上有突起', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '水炮龟', description: '能发水炮，壳更硬', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '水箭龟雏形', description: '背部长出炮管', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '水箭龟', description: '双炮管，能高速水射', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '超级水箭龟', description: '背炮变大，全身水光', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '海神龟', description: '周围有洋流，可召唤浪', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '远古水龟', description: '壳如岛屿，背负山河', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'eevee', name: '伊布', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '伊布蛋', description: '棕色蛋，有狐尾纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼狐', description: '棕色小狐，耳大', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '伊布', description: '标准形态，毛色棕褐，尾蓬松', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '活泼伊布', description: '眼神灵动，毛发光泽', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '进化预备', description: '额头出现宝石（颜色不定）', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '分支初现', description: '身体散发对应元素光', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '元素显现', description: '毛色改变，尾巴变长', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '水伊布', description: '鳞片发光，能化水', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '水伊布·潮汐', description: '周身洋流环绕', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '水伊布·漩涡', description: '可召唤巨大漩涡', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '水伊布·海神', description: '半透明水化身躯', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒水伊布', description: '化作清泉，生生不息', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'pikachu', name: '皮卡丘', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '电气蛋', description: '黄色蛋壳，闪电纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '皮丘', description: '小型，脸颊无电', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '皮丘萌动', description: '能站立，脸颊微电', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '皮卡丘', description: '标准形态，脸颊带电，尾巴闪电状', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '蓄电皮卡丘', description: '电力增强，奔跑带光', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '雷丘雏形', description: '体型增大，尾巴变粗', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '雷丘', description: '完全形态，电力十足', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '冲浪皮卡丘', description: '尾巴可当冲浪板', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '飞空皮卡丘', description: '能用气球飞行', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '搭档皮卡丘', description: '戴帽子，电光四射', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '极巨化皮卡丘', description: '体型巨大，雷云随行', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '雷电之神', description: '周身落雷，可控制气候', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'riolu', name: '利欧路→路卡利欧', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '波导蛋', description: '蓝色蛋，有波动纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鲁', description: '小兽形，后肢站立', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '利欧路', description: '蓝色，胸前有黑纹', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '波导学徒', description: '能感知波导，眼发光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '练习生', description: '拳脚有力，速度提升', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '进化前夜', description: '身体发光，骨骼强化', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '路卡利欧', description: '人形，全身蓝黑，手有骨刺', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '波导战士', description: '能发波导弹，眼神锐利', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '波导大师', description: '骨刺变长，全身气势', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '超级路卡利欧', description: '肌肉增强，波导爆发', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '传说守护者', description: '全身金色纹路，力量内敛', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '波导之神', description: '无形无质，可化波导，守护众生', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 三、国宝系列（6种） =====
  {
    id: 'national', name: '国宝', emoji: '🐼',
    species: [
      {
        id: 'panda', name: '大熊猫', seriesId: 'national',
        levels: [
          { level: 1, name: '粉红团子', description: '粉红色幼崽，闭眼蜷缩', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '黑白初现', description: '黑白毛色开始显现，眼周变黑', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '蹒跚学步', description: '能缓慢爬动，抱住竹子', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '小熊猫', description: '黑白分明，圆滚滚，行动笨拙', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '好奇熊猫', description: '开始站立，东张西望', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '竹林顽童', description: '喜欢打滚，抱着竹子玩耍', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '青涩少年', description: '体型增大，毛色亮泽', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '竹林卫士', description: '端坐竹林，神情沉稳', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '成年熊猫', description: '憨态可掬，悠然吃竹', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '竹林尊者', description: '体型硕大，周身有竹叶飘落', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '太极熊猫', description: '黑白形成太极图，气息沉稳', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '竹林圣者', description: '全身散发柔和光芒，竹影随行', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'golden_monkey', name: '金丝猴', seriesId: 'national',
        levels: [
          { level: 1, name: '金丝灵胎', description: '淡金绒毛初生，蓝面轮廓隐约，山涧灵光中的灵胎', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '灰白幼猴', description: '全身灰白，紧抱树干', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '初显金毛', description: '金色毛发逐渐显现', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年金猴', description: '金毛半覆，活泼好动', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '蓝面初显', description: '面部变蓝，毛发金黄', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '金丝猴', description: '标准形态，长尾，金毛亮丽', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '长尾猴', description: '尾长过身，攀爬如飞', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '壮年金猴', description: '肌肉结实，毛色深金', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金猴王', description: '体型最大，威严，毛发披肩', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '金丝尊者', description: '周身金光，面容威严', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '金丝大圣', description: '手执金棒，身披铠甲', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '齐天金猴', description: '腾云驾雾，毫毛可化万千', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'red_crowned_crane', name: '丹顶鹤', seriesId: 'national',
        levels: [
          { level: 1, name: '鹤卵', description: '灰绿色蛋，有斑点', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏鹤', description: '黄褐色绒毛，蹒跚学步', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '灰羽渐白', description: '羽毛变灰白，开始练习飞翔', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年鹤', description: '白羽初成，头顶微红', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '飞羽初展', description: '能短距离飞翔，姿态优美', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '丹顶鹤', description: '白羽黑颈，头顶朱红', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '长寿鹤', description: '朱红更艳，翅尖有黑羽', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '仙鹤', description: '体型修长，动作优雅', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '祥云伴鹤', description: '飞翔时带云气，声闻于天', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '云中仙鹤', description: '脚踏祥云，头顶有光环', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '羽化鹤', description: '身体半透明，羽翼发光', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '南极仙翁之鹤', description: '周身仙气，可驮仙人，长生不老', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'south_china_tiger', name: '华南虎', seriesId: 'national',
        levels: [
          { level: 1, name: '虎纹胚胎', description: '橙黑斑纹初显，额间「王」字虚影，卧于虎穴微光中', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '虎崽', description: '小虎，条纹初现，奶凶可爱', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '幼虎', description: '能奔跑，眼睛有神', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年虎', description: '体型渐长，开始独立', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '斑纹加深', description: '条纹变黑，肌肉初显', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '壮年虎', description: '王者姿态，花纹鲜明', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '百兽之王', description: '威风凛凛，目光如炬', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '虎王', description: '额前"王"字清晰，吼声震林', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金虎', description: '毛色变金，条纹暗红', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '白虎', description: '通体雪白，条纹银灰', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '神虎', description: '背生双翼，能飞', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '虎神', description: '周身雷电，行走间山摇地动', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'chinese_alligator', name: '扬子鳄', seriesId: 'national',
        levels: [
          { level: 1, name: '鳄卵', description: '白色蛋，粗糙', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鳄', description: '小鳄鱼，头大身小', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '鳞甲初生', description: '背甲变硬，能游泳', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年鳄', description: '体型增大，尾巴有力', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '水中霸主', description: '游泳迅速，能潜伏', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '扬子鳄', description: '标准形态，暗色鳞甲', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '铁甲鳄', description: '鳞甲变厚，呈盔甲状', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '古鳄', description: '嘴长，牙齿锋利', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '大鳄', description: '身长数米，威严', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '龙鳄', description: '头生角，似龙非龙', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '蛟鳄', description: '身带水汽，可兴风浪', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '扬子龙', description: '化为龙形，镇守长江', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'crested_ibis', name: '朱鹮', seriesId: 'national',
        levels: [
          { level: 1, name: '朱卵', description: '淡粉色蛋', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '灰羽雏鸟', description: '灰色绒毛，喙长', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '白羽渐生', description: '白色羽毛出现', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '朱色初现', description: '头部开始泛红', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '朱鹮', description: '白羽红面，长喙', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '红羽加深', description: '红色更深，翅尖带粉', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '祥瑞之鸟', description: '飞翔时带彩虹', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '仙鸟', description: '体态优雅，步履从容', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '吉祥朱鹮', description: '头顶有金光', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '太阳鸟', description: '周身发暖光', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '凤凰近亲', description: '尾羽变长，五彩', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '涅槃朱鹮', description: '浑身火焰般羽毛，可重生', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 四、数码宝贝系列（6种） =====
  {
    id: 'digimon', name: '数码宝贝', emoji: '🛡️',
    species: [
      {
        id: 'mecha_dragon', name: '亚古兽', seriesId: 'digimon',
        levels: [
          { level: 1, name: '数码蛋', description: '从数码世界降下的蛋，表面有数码纹路', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '黑球兽', description: '幼年期I，滚动的黑色小球', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '滚球兽', description: '幼年期II，橙色的圆滚小恐龙', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '亚古兽', description: '成长期！橙色小恐龙，喷出小型火焰', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '亚古兽·干劲', description: '背甲变硬，爪牙更利', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '暴龙兽', description: '成熟期！巨大暴龙，吐出暴龙火焰', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '暴龙兽·强攻', description: '下颚强化，火焰更猛', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '机械暴龙兽', description: '完全体！半机械装甲，发射究极破坏炮', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '机械暴龙兽·强化', description: '机械翼展开，火力全开', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '战斗暴龙兽', description: '究极体！双爪持勇者之盾，跃向天空', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '战斗暴龙兽·勇气', description: '勇气徽章闪耀，龙兽克星发威', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '胜利暴龙兽', description: '圣骑士型，手持剑刃，沐浴数码圣光', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'cyber_cat', name: '迪路兽', seriesId: 'digimon',
        levels: [
          { level: 1, name: '数码蛋', description: '圣洁的数码蛋，隐隐发光', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '咪罗兽', description: '幼年期I，蜷缩的小毛球', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '咪罗兽·成长', description: '幼年期II，长出短毛', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '迪路兽', description: '成长期！白色圣猫，戴着神圣之环', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '迪路兽·敏捷', description: '爪击与猫眼更快更准', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '猫人兽', description: '成熟期！人立猫身，手持雷神之棒', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '猫人兽·雷霆', description: '雷神之棒劈出闪电', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '神圣天女兽', description: '完全体！闪耀的神圣天女，展开光翼', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '神圣天女兽·圣光', description: '圣光净化一切邪恶', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '圣龙兽', description: '究极体！神圣巨龙，净化之辉环绕', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '圣龙兽·光明', description: '光明徽章觉醒', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '奥林匹斯十二神·密涅瓦兽', description: '智慧女神降临，统御数码圣域', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'space_mecha', name: '巴达兽', seriesId: 'digimon',
        levels: [
          { level: 1, name: '数码蛋', description: '轻盈的数码蛋，会漂浮', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '浮球兽', description: '幼年期I，漂浮的气泡', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '浮球兽·成长', description: '幼年期II，开始晃动', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '巴达兽', description: '成长期！奶油色的小翅膀天使，会飞', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '巴达兽·轻盈', description: '翅膀更鼓，飞行更稳', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '天使兽', description: '成熟期！圣洁天使，手持天堂之杖', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '天使兽·圣杖', description: '天堂之杖释放神圣光线', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '神圣天使兽', description: '完全体！八翼天使，审判邪恶', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '神圣天使兽·审判', description: '圣剑出鞘，天国之门显现', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '炽天使兽', description: '究极体！三大天使之首，燃尽黑暗', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '炽天使兽·希望', description: '希望徽章闪耀圣光', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '神之侍从·炽天使', description: '完全圣体，御座之前万翼齐辉', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'quantum_beast', name: '加布兽', seriesId: 'digimon',
        levels: [
          { level: 1, name: '数码蛋', description: '静谧的数码蛋，带冰晶纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '独角兽', description: '幼年期I，蓝色的冰球', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '独角兽·成长', description: '幼年期II，长出小角', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '加布兽', description: '成长期！身披加鲁鲁兽毛皮的小狼', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '加布兽·机敏', description: '披风般的毛皮更厚，蓝角更亮', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '加鲁鲁兽', description: '成熟期！蓝色巨狼，咆哮震地', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '加鲁鲁兽·狂岚', description: '狼爪与犬齿如风暴般凌厉', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '兽人加鲁鲁', description: '完全体！直立狼人，手持三重刃', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '兽人加鲁鲁·月狼', description: '月光下力量倍增', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '钢铁加鲁鲁兽', description: '究极体！机械钢狼，翼展掠空', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '钢铁加鲁鲁兽·友情', description: '友情徽章驱动冰封之力', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '奥米加兽·右翼', description: '与战斗暴龙兽合体，究极圣剑出鞘', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'digital_phoenix', name: '比丘兽', seriesId: 'digimon',
        levels: [
          { level: 1, name: '数码蛋', description: '粉色的数码蛋，羽纹浮现', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '皮皮兽', description: '幼年期I，粉色的绒球', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '皮皮兽·成长', description: '幼年期II，长出绒毛', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '比丘兽', description: '成长期！粉色小鸟，头顶凤冠', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '比丘兽·鸣唱', description: '歌声能安抚人心', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巴多拉兽', description: '成熟期！展翅的火焰巨鸟', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '巴多拉兽·烈焰', description: '羽翼燃起绯红之焰', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '伽楼达兽', description: '完全体！圣鸟伽楼达，翅展净化之光', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '伽楼达兽·疾风', description: '双翼卷起净化风暴', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '凤凰兽', description: '究极体！传说圣凰，自火焰中重生', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '凤凰兽·涅槃', description: '爱与纯真之光化作羽翼', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '究极圣凰', description: '圣光洒落，万物复苏', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'mecha_shark', name: '哥玛兽', seriesId: 'digimon',
        levels: [
          { level: 1, name: '数码蛋', description: '海蓝的数码蛋，波浪纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '噗噗兽', description: '幼年期I，透明水泡', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '噗噗兽·成长', description: '幼年期II，甩动小尾', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '哥玛兽', description: '成长期！白色海狮，圆滚滚会卖萌', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '哥玛兽·潜水', description: '潜入深海毫不费力', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '海狮兽', description: '成熟期！海狮斗士，挥舞尖角长矛', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '海狮兽·破浪', description: '长矛破浪，卷起漩涡', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '祖顿兽', description: '完全体！海格力斯巨锤砸碎冰山', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '祖顿兽·重锤', description: '巨锤挥下，地动山摇', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '维京兽', description: '究极体！维京王者，召唤雷霆与海啸', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '维京兽·咆哮', description: '王者咆哮，冰封千里', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '维京兽·传说', description: '骑乘雷龙兽，统御数码冰海', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 五、魔法奇幻系列（6种） =====
  {
    id: 'magic', name: '魔法奇幻', emoji: '🦄',
    species: [
      {
        id: 'unicorn', name: '独角兽', seriesId: 'magic',
        levels: [
          { level: 1, name: '星光种子', description: '闪烁的星形光点', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '萌芽小马', description: '光点拉长成小马轮廓', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '白色小驹', description: '纯白身体，额头微凸', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '幼年独角兽', description: '长出小角，鬃毛发光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '角光初现', description: '角尖发出柔和荧光', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '银鬃飞马', description: '鬃毛银色，能短途漂浮', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '彩虹鬃毛', description: '鬃毛七彩渐变，角更亮', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '圣光独角兽', description: '全身圣光，蹄下生花', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星夜独角兽', description: '鬃毛如星空，角尖有星环', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '梦境守护者', description: '半透明，能入梦境', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '永恒独角兽', description: '金色鬃毛，角刻符文', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '创世独角兽', description: '全身星光，角能开辟空间', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'wyvern', name: '飞龙', seriesId: 'magic',
        levels: [
          { level: 1, name: '龙蛋', description: '暗红色，有鳞纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼龙', description: '无翼，四足，能喷火', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '翼芽', description: '背部出现小翼芽', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '翼龙', description: '翼长成，能滑翔', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '火焰龙', description: '喷火强，翼大', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '风龙', description: '可控制气流', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '雷龙', description: '带雷电', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '冰龙', description: '喷冰霜', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '古龙', description: '体型巨大，皮糙', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '龙王', description: '身披铠甲，号令群龙', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '太古龙', description: '智慧极高，魔法强大', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '龙神', description: '半神，掌控元素', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'fairy', name: '精灵', seriesId: 'magic',
        levels: [
          { level: 1, name: '花苞', description: '含苞待放', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '花精灵', description: '小花形，有透明翅', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '叶精灵', description: '绿色，能控植物', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '光精灵', description: '发光，能治愈', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '暗精灵', description: '暗色，速度极快', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '水精灵', description: '蓝色，控水', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '火精灵', description: '红色，控火', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '风精灵', description: '白色，隐身', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '大地精灵', description: '褐色，控制土石', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '自然精灵', description: '融合元素', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '森林之子', description: '与自然一体', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '精灵王', description: '号令自然，永生', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'treant', name: '树人', seriesId: 'magic',
        levels: [
          { level: 1, name: '种子', description: '橡果状', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '嫩芽', description: '破土而出', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '小树苗', description: '长出叶子', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '树童', description: '有简单五官', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '树人', description: '四肢为树枝', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '橡树人', description: '身硬，能投石', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '古树人', description: '身缠藤蔓，苔藓', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '战争树人', description: '可拔根移动', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '森林卫士', description: '高大，持木盾', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '智慧古树', description: '千年智慧', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '世界树幼苗', description: '连接天地', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒古树', description: '化为森林之灵', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'griffin', name: '狮鹫', seriesId: 'magic',
        levels: [
          { level: 1, name: '狮鹫蛋', description: '金色，有爪印纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼狮鹫', description: '狮身鹰头，小翼', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '羽翼初长', description: '翼变大', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '狮鹫', description: '标准，能飞', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '金翼狮鹫', description: '翼金色，能闪光', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '银鬃狮鹫', description: '鬃毛银色', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '皇家狮鹫', description: '戴王冠', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '风暴狮鹫', description: '控制天气', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '雷光狮鹫', description: '带电', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '格里芬王', description: '统领鸟兽', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '神圣狮鹫', description: '圣光环绕', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '天空霸主', description: '遮天蔽日', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'mermaid', name: '美人鱼', seriesId: 'magic',
        levels: [
          { level: 1, name: '珍珠', description: '内含人影', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鱼人', description: '上半身人，尾为鱼', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '鳞光', description: '鳞片发光', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '人鱼', description: '长发，歌声动听', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '海妖', description: '魅惑歌声', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '深海人鱼', description: '适应高压', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '珊瑚公主', description: '头戴珊瑚冠', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '海洋之女', description: '控水，召唤海兽', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '人鱼女王', description: '威严，珍珠项链', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '海神之女', description: '半神，掌控洋流', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '大洋之魂', description: '化为海水', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒人鱼', description: '守护海洋', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 六、史前生物系列（6种） =====
  {
    id: 'prehistoric', name: '史前生物', emoji: '🦕',
    species: [
      {
        id: 't_rex', name: '霸王龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '化石蛋', description: '有化石纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳小爪', description: '伸出爪子，带鳞', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '幼年暴龙', description: '全身鳞片，短尾', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年暴龙', description: '前肢短，后肢强', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '追赶者', description: '能快速奔跑', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '森林霸主', description: '体型大，皮肤深', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '领地之王', description: '头顶角质冠', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '巨型暴龙', description: '身长数米，牙尖', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '白垩纪领主', description: '厚重鳞甲', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '暴龙之王', description: '背有棘刺，尾有力', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '终极掠食者', description: '青铜色，坚硬', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '恐龙至尊', description: '全身远古气息，吼震屏幕', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'triceratops', name: '三角龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '角蛋', description: '有三个凸点', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼角龙', description: '头部有三个小角', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '头盾初现', description: '头后出现骨盾', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '三角龙', description: '角变长，盾大', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '盾甲龙', description: '盾更厚，能防御', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '冲锋龙', description: '能奔跑冲撞', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '铁头龙', description: '角如铁', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '巨三角', description: '体型巨大', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '群首领', description: '带领族群', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '古三角', description: '角呈金色', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '化石之王', description: '全身硬化', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒三角', description: '化为山石，守护大地', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'pterosaur', name: '翼龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '飞卵', description: '带翅膜纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '翼苗', description: '小翼，会滑翔', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '翼手龙', description: '爪尖，翼展宽', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '飞翔者', description: '能盘旋', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '风神翼龙', description: '巨大，能翱翔九天', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '雷电翼龙', description: '带静电', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '火焰翼龙', description: '喷火', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '寒冰翼龙', description: '冰冻', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '王者翼龙', description: '统领天空', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '天空领主', description: '遮天', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '太古翼龙', description: '化石复活', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒之翼', description: '化为风', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'mammoth', name: '猛犸象', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '长毛灵胎', description: '厚绒毛裹着霜晶，象形初具的冰原灵胎', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '小猛犸', description: '长毛，象牙初露', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '象牙变长', description: '牙突出', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '猛犸象', description: '标准，长毛', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '冰原猛犸', description: '耐寒，身有冰晶', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨猛犸', description: '体型巨大', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '象牙王', description: '牙如玉', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '古猛犸', description: '身披厚毛', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '冰川之主', description: '行走于冰', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '神猛犸', description: '周身寒气', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '远古巨兽', description: '如移动冰山', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒猛犸', description: '化为冰雕，守护冰原', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'sabertooth', name: '剑齿虎', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '剑齿灵胎', description: '肌肉初成，上颚獠牙的锋利雏形若隐若现', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼剑齿', description: '上犬齿突出', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '长牙初成', description: '牙变长', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '剑齿虎', description: '标准，牙如剑', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '刃牙虎', description: '牙更锋利', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨剑虎', description: '体型大', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '冰河猎手', description: '适应寒冷', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '王剑齿', description: '统治草原', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金剑齿', description: '牙金色', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '神剑虎', description: '身体发光', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '太古剑齿', description: '化石之力', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒剑齿', description: '化为宝石', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'mosasaur', name: '沧龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '水卵', description: '有鳞片', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼沧龙', description: '小型，水性好', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '鳞甲变厚', description: '能潜水', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '沧龙', description: '标准，长尾', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '海中霸主', description: '速度极快', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨沧', description: '体型巨大', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '王沧', description: '头有冠', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '深海沧龙', description: '适应深水', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '古沧', description: '身带化石纹', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '海王沧龙', description: '控制海洋', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '神沧', description: '半透明，神力', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒沧龙', description: '化为海流', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 七、星座守护系列（12种） =====
  {
    id: 'constellation', name: '星座守护', emoji: '♈',
    species: [
      {
        id: 'aries', name: '白羊座·穆', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星屑之种', description: '一粒闪耀的星屑，孕育着白羊的轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '穆光初照', description: '星尘凝成幼羊，额间泛起柔和的念力微光', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '念动初觉', description: '小羊能用意念轻轻挪动石块，专注而温柔', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '白羊修行', description: '静心修习念动力，发丝间星尘缓缓飞舞', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '水晶微光', description: '掌间凝出半透明的晶墙雏形，护住周身', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '念动力觉醒', description: '一个眼神便能让敌手悬空，神情依旧从容', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '水晶墙', description: '一面晶莹的水晶之墙在身前展开，坚不可摧', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '星屑旋转功', description: '双手舞动，召唤旋转的星屑如银河倾泻', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '白羊宫主', description: '身着白羊黄金圣衣，端坐白羊宫守护圣域之门', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '圣衣修复师', description: '指尖微光流转，将碎裂的圣衣缓缓修复如新', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '星光灭绝', description: '星光化作漫天洪流吞没敌影，却仍带着慈悲', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '圣域智者', description: '看透轮回与宿命，以睿智与温柔守护信念', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'taurus', name: '金牛座·阿鲁迪巴', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉之种', description: '一粒闪耀的星辉，孕育着金牛的轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼角初露', description: '星尘凝成小牛，额间两枚金色角芽悄然萌发', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '蛮力渐长', description: '小牛力气惊人，轻易顶开巨石，憨态可掬', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '金牛修行', description: '日复一日打磨角与身躯，豪爽开朗而可靠', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '角力初成', description: '双角愈发粗壮，一顶之力可裂坚岩', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨型号角雏形', description: '低头蓄力，双角泛起黄金之辉', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '巨型号角', description: '猛然抬头，金色冲击波将敌手掀飞千丈', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '勇猛无畏', description: '大步流星冲入敌阵，力量与热血并进', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金牛宫主', description: '身披金牛黄金圣衣守护金牛宫，憨厚可靠', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '圣域之盾', description: '以魁梧身躯挡在同伴之前，是最可靠的战友', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '黄金之角', description: '双角化作惊天巨号，一声怒吼震裂天地', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '赤诚勇士', description: '豪爽忠勇，为守护同伴倾尽全力，令人敬服', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'gemini', name: '双子座·撒加', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉双子', description: '一粒星辉分成两团，光明与黑暗交织缠绕', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '双生之影', description: '星尘凝成一对孪生少年，一善一恶静静对望', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '善恶交织', description: '时而温柔如圣，时而狞笑如魔，同体双生', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '神之化身', description: '体内涌出惊人小宇宙，举手投足若神临世', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '异次元微光', description: '指尖撕裂空间，浮现一道漆黑裂缝', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '幻胧魔皇拳', description: '拳影朦胧，夺人心魄，中者心智尽失', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '异次元空间', description: '双手一合，将敌手卷入无尽的异次元', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '堕落之影', description: '善念沉眠，邪恶化身成为伪教皇君临圣域', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '双子宫主', description: '身披双子黄金圣衣，善恶两魂在衣上对峙', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '银河星爆雏形', description: '双手凝聚星云之力，银河在掌中缓缓盘旋', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '银河星爆', description: '星云爆发，银河粉碎，光芒吞没一切', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '善之回归', description: '战胜内心黑暗，以神之化身守护雅典娜', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'cancer', name: '巨蟹座·迪斯马斯克', seriesId: 'constellation',
        levels: [
          { level: 1, name: '冥辉之种', description: '一粒暗紫色星辉，内里泛起冥界的气息', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '巨蟹之影', description: '星尘凝成小蟹，甲壳泛着幽幽的蓝光', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '冥气初生', description: '周身缠绕若有若无的积尸气，阴风阵阵', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '力量信奉', description: '狂妄地相信力量即正义，眼神桀骜不驯', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '冥界门扉', description: '指尖划开一道冥界裂缝，幽风透骨而来', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '积尸气雏形', description: '掌心聚起青色气旋，似可牵引亡魂', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '积尸气冥界波', description: '冥界波横扫，将敌人的灵魂拖向冥界深渊', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '邪气凛然', description: '狂妄大笑，积尸气缠身如冥界君主', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '巨蟹宫主', description: '身披巨蟹黄金圣衣，镇守通往冥界的门', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '冥界通道', description: '双手一划，开启通往冥界的通道，亡魂低语', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '冥界之王', description: '与冥界之力共鸣，积尸气如海啸吞没一切', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '因果报应', description: '终悟善恶有报，以冥界之力惩罚真正的邪恶', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'leo', name: '狮子座·艾奥里亚', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉之种', description: '一粒炽热的星辉，孕育着雄狮的轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '狮心初燃', description: '星尘凝成幼狮，鬃毛间有火花跳动', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '热血修行', description: '日夜修炼拳速，热血沸腾，渴望变强', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '正义之心', description: '胸中燃起正义之火，绝不向邪恶低头', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '拳光初现', description: '拳锋掠过空气，带起细碎的闪电', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '闪电拳雏形', description: '出拳若电，残影连成一道金色光弧', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '闪电光速拳', description: '光速之拳如雷暴倾泻，撕裂空气', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '光明磊落', description: '光明正大对敌，气势如虹令宵小胆寒', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '狮子宫主', description: '身披狮子黄金圣衣守护狮子宫，意气风发', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '雷光之鬃', description: '鬃毛化作金色雷光，拳风撕天裂地', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '等离子光速拳', description: '光速拳化作等离子风暴，灼烧一切黑暗', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '正义雷光', description: '以热血与正义化为雷霆，荡尽世间邪恶', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'virgo', name: '处女座·沙加', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉莲花', description: '一粒星辉如含苞的莲台，静卧星空深处', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '闭目冥思', description: '星尘凝成少年佛子，双目紧闭，宝相庄严', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '禅心初定', description: '于恒河畔禅定入静，参悟生灭无常', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '近神之资', description: '小宇宙深邃如海，被称"最接近神的人"', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '天眼微启', description: '闭目冥想间，天眼已照见三千世界', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '六道轮回', description: '指尖拨动六道轮回，众生沉浮于其间', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '天魔降伏', description: '佛光绽放，一掌降伏天魔，梵音缭绕', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '曼陀罗光', description: '周身曼陀罗光华流转，神圣不可侵', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '处女宫主', description: '身披处女黄金圣衣，端坐处女宫悟道修行', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '天舞宝轮', description: '化出攻防一体的宝轮，剥夺敌人的六感', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '第八感觉醒', description: '于沙罗双树下领悟第八感，超脱生死轮回', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '佛陀化身', description: '慧眼洞彻宇宙真相，以身护道，慈悲无量', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'libra', name: '天秤座·童虎', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉龙珠', description: '一粒星辉如龙珠，内里盘着一条青龙', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '庐山谷音', description: '星尘凝成小童，耳边回响庐山瀑布轰鸣', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '百裂拳初修', description: '跟随师父修习拳法，一拳击碎青石', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '升龙雏形', description: '拳势中隐现升龙之影，虎虎生风', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '五老峰悟道', description: '于五老峰静坐，参悟天地磅礴之势', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '升龙霸雏形', description: '一拳冲天，龙吟隐隐，瀑布为之倒卷', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '庐山升龙霸', description: '升龙拳出，青龙盘柱而上，惊天动地', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '天秤武库', description: '手抚天秤黄金武器，可拆解出十二件神兵', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '老师坐镇', description: '以假死之术镇守五老峰，被世人尊称"老师"', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '天秤宫主', description: '身披天秤黄金圣衣，十八般武器任取', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '返老还童', description: '小宇宙苏醒恢复青春之躯，重现战神威仪', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '庐山龙神', description: '升龙霸威震八荒，守护大地两百年不灭', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'scorpio', name: '天蝎座·米罗', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉之种', description: '一粒星辉如蝎尾针尖，泛着深红光芒', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '蝎针初成', description: '星尘凝成小蝎，尾针鲜红欲滴', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '忠义修行', description: '以忠诚为道修行，身手敏捷如电', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '正义之针', description: '深信以毒针制裁罪恶，针尖必中要害', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '毒针微芒', description: '指尖弹出细小红针，精准无比', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '深红毒针雏形', description: '深红之针破空而出，带着灼热的红芒', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '深红毒针', description: '十五针连发如天蝎甩尾，针针致命', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '忠肝义胆', description: '为挚友甘愿流血，为正义绝不动摇', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '天蝎宫主', description: '身披天蝎黄金圣衣守护天蝎宫，耿直豪迈', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '针之神技', description: '毒针可留人活口亦可瞬间夺命，收发由心', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '天蝎之怒', description: '深红毒针如暴雨倾泻，刺破一切罪恶', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '赤诚星辰', description: '忠义长存，深红针芒化作守护同伴的星光', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'sagittarius', name: '射手座·艾俄洛斯', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉之箭', description: '一粒星辉如出鞘之箭，锋芒毕露', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '英魂初现', description: '星尘凝成少年射手，目光坚毅如星', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '弓弦初张', description: '苦练箭术，箭意直指星辰', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '仁者之心', description: '心怀慈悲与正义，被推为教皇候选人', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '闪电拳雏形', description: '拳速若闪电，同时稳稳拉开弓弦', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '原子闪电光速拳', description: '光速拳如原子迸裂，拳影密布苍穹', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '黄金之箭', description: '弯弓搭箭，箭尖汇聚出黄金光芒', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '守护雅典娜', description: '为救雅典娜拼死突围，将婴儿托付他人', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '英魂长存', description: '英魂寄于射手圣衣，指引少年们前行', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '射手宫主', description: '身披射手黄金圣衣，箭意穿过十二宫', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '黄金之箭·觉醒', description: '英魂苏醒，射出贯穿一切黑暗的黄金之箭', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '守护之魂', description: '牺牲精神化作永恒星光，守护雅典娜直至永远', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'capricorn', name: '摩羯座·修罗', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉之剑', description: '一粒星辉如剑刃寒光，锐利夺目', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '剑意初萌', description: '星尘凝成少年剑客，手刀如剑', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '剑士修行', description: '以身为剑日夜劈斩，磨砺手刀', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '忠义之刃', description: '将忠诚奉为信条，为圣域甘愿赴死', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '圣剑雏形', description: '手刀泛出金色剑光，可断金石', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '圣剑·Excalibur', description: '手刀化为传说之剑，一剑斩裂万物', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '剑影纵横', description: '圣剑连斩，剑气纵横交错成网', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '摩羯之斩', description: '为忠诚而挥剑，不辨善恶只守信义', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '摩羯宫主', description: '身披摩羯黄金圣衣，镇守摩羯宫剑指来敌', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '圣剑加身', description: '全身化作无数圣剑，攻势如万剑齐发', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '觉醒之剑', description: '在舍身中顿悟，圣剑化为守护之刃', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '忠诚之剑', description: '以忠诚为魂，圣剑斩恶护善，至死不渝', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'aquarius', name: '水瓶座·卡妙', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉冰晶', description: '一粒星辉如冰晶，寒冷彻骨', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '霜华初凝', description: '星尘凝成少年，呼出的气凝成冰霜', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '冰道修行', description: '于冰原苦修冰系斗技，冷静如冰', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '严师之心', description: '以严苛教导徒弟，望其超越自己', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '冰晶环雏形', description: '指尖凝出冰晶之环，缓缓旋转', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '冰晶环', description: '冰环锁住敌人，寒气彻骨，冰封一切', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '绝对零度之道', description: '追求绝对零度，冰拳所过万物冻结', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '曙光微光', description: '拳中孕育曙光女神般的金色光芒', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '水瓶宫主', description: '身披水瓶黄金圣衣，端坐水瓶宫冷峻如霜', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '曙光女神之宽恕', description: '曙光初现，冰封一切却又留下一线生机', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '冰河之师', description: '师徒对决中舍命成全，将绝对零度传于弟子', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒之冰', description: '严寒之下是炽热的师者之心，冰封中透着慈爱', stage: 'transcendent', requiredScore: 450 },
        ],
      },
      {
        id: 'pisces', name: '双鱼座·阿布罗狄', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辉花苞', description: '一粒星辉如含苞玫瑰，馥郁芬芳', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '玫瑰初绽', description: '星尘凝成美少年，指尖绕着一枝玫瑰', stage: 'egg', requiredScore: 15 },
          { level: 3, name: '花之修行', description: '以玫瑰为兵，人比花艳，美得致命', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '美之战士', description: '信奉力量与美，玫瑰在指尖绽放', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '恶魔玫瑰雏形', description: '血红玫瑰绽放，散出致命的花粉', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '皇家恶魔玫瑰', description: '花瓣纷飞，触之即中毒而亡', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '食人鱼玫瑰', description: '玫瑰化作食人鱼群，撕咬吞噬', stage: 'mature', requiredScore: 165 },
          { level: 8, name: '玫瑰花园', description: '遍地玫瑰怒放，美丽而杀机四伏', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '双鱼宫主', description: '身披双鱼黄金圣衣，镇守双鱼宫，花海漫漫', stage: 'legendary', requiredScore: 260 },
          { level: 10, name: '血腥玫瑰', description: '玫瑰刺入心口，以最凄美之姿赐死', stage: 'legendary', requiredScore: 315 },
          { level: 11, name: '凄美绽放', description: '为守护绽放最后一朵玫瑰，美丽与牺牲同在', stage: 'transcendent', requiredScore: 375 },
          { level: 12, name: '永恒之花', description: '玫瑰盛开之处，美丽、忠诚与希望永不凋零', stage: 'transcendent', requiredScore: 450 },
        ],
      },
    ],
  },

]

// ============================================================
// 扩充数据合并：每系列补齐到 12 种 + 新增「传统节日」「虹猫蓝兔七侠传」系列
// 构成 10 系列 × 12 物种 = 120 种
// ============================================================

for (const series of PET_SERIES) {
  const extra = EXTRA_SPECIES[series.id]
  if (extra?.length) {
    series.species.push(...extra)
  }
}
PET_SERIES.push(...EXTRA_SERIES)

// 统一每级阶段字段（旧数据用旧映射，统一到 Lv.1-2 卵生/3-5 幼年/6-8 成长/9-10 成熟/11-12 传说）
for (const series of PET_SERIES) {
  for (const sp of series.species) {
    for (const lvl of sp.levels) {
      lvl.stage = getLevelStage(lvl.level)
    }
  }
}

/** 根据物种ID获取宠物数据 */
export function getSpeciesById(speciesId: string): PetSpecies | undefined {
  for (const series of PET_SERIES) {
    const found = series.species.find(s => s.id === speciesId)
    if (found) return found
  }
  return undefined
}

/** 根据物种ID获取系列 */
export function getSeriesBySpeciesId(speciesId: string): PetSeries | undefined {
  for (const series of PET_SERIES) {
    if (series.species.some(s => s.id === speciesId)) return series
  }
  return undefined
}

/** 获取指定等级的数据 */
export function getLevelData(speciesId: string, level: number): PetLevel | undefined {
  const species = getSpeciesById(speciesId)
  if (!species) return undefined
  return species.levels.find(l => l.level === level)
}

/** 获取宠物当前等级名称 */
export function getPetLevelName(speciesId: string, level: number): string {
  const data = getLevelData(speciesId, level)
  return data?.name || `Lv.${level}`
}

/** 获取宠物当前等级描述 */
export function getPetLevelDescription(speciesId: string, level: number): string {
  const data = getLevelData(speciesId, level)
  return data?.description || ''
}

/** 获取全部系列（用于选择面板） */
export function getAllSeries(): PetSeries[] {
  return PET_SERIES
}

/** 获取系列中的物种 */
export function getSpeciesBySeries(seriesId: string): PetSpecies[] {
  const series = PET_SERIES.find(s => s.id === seriesId)
  return series?.species || []
}
