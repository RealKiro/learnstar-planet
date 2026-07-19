// ===== 学趣星球 · 宠物数据系统 =====
// 8 个系列、54 个物种、每个物种 12 级进化
// 数据来源：宠物系列设计.txt

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
  /** 等级阶段分类: 1-卵/种子, 2-4幼年, 5-7成长, 8-10成熟, 11-12传说 */
  stage: 'egg' | 'baby' | 'growing' | 'mature' | 'legendary'
  /** 所需累计积分 */
  requiredScore: number
}

/** 获取等级所需积分 */
export function getLevelRequiredScore(level: number): number {
  const scores = [0, 0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450]
  return scores[level] ?? 450
}

/** 判断等级阶段 */
export function getLevelStage(level: number): PetLevel['stage'] {
  if (level === 1) return 'egg'
  if (level <= 4) return 'baby'
  if (level <= 7) return 'growing'
  if (level <= 10) return 'mature'
  return 'legendary'
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
  mecha: {
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
  folklore: {
    primaryColor: '#F97316',
    bgGradient: 'linear-gradient(180deg, #1a0a00 0%, #3a1a00 30%, #5a2a00 60%, #1a0a00 100%)',
    decor: ['🏮', '🧧', '🎊', '🐉'],
    scene: '古风小镇 · 万家灯火',
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
  // 民间
  nian: '🦁',
  dragon_boat: '🚣',
  lantern: '🏮',
  dumpling: '🥟',
  fu_star: '⭐',
  shou_star: '🍑',
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
    folklore: '🏮',
  }
  return map[seriesId] || '🌟'
}

/** 获取系列名 */
export function getSeriesName(seriesId: string): string {
  const map: Record<string, string> = {
    myth: '古代神话',
    pokemon: '宝可梦风格',
    national: '国宝',
    mecha: '科幻机甲',
    magic: '魔法奇幻',
    prehistoric: '史前生物',
    constellation: '星座守护',
    folklore: '民间传说',
  }
  return map[seriesId] || '未知'
}

// ============================================================
// 完整宠物数据：8 个系列，54 个物种
// ============================================================

export const PET_SERIES: PetSeries[] = [
  // ===== 一、古代神话系列（6种） =====
  {
    id: 'myth', name: '古代神话', emoji: '🏔️',
    species: [
      {
        id: 'zhulong', name: '烛龙', seriesId: 'myth',
        levels: [
          { level: 1, name: '烛火微光', description: '一团跳动的橙红色火焰，中心有竖瞳', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '烛火之影', description: '火焰中浮现细长蛇形轮廓，忽明忽暗', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '衔烛之苗', description: '蛇身人面，口含小火苗，双目紧闭', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '烛阴初醒', description: '身体拉长，鳞片隐约可见，睁一只眼', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '衔烛而行', description: '口中蜡烛燃起，照亮周围半米', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '烛照四方', description: '身体环绕光晕，可照亮暗处，鳞片发金', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '昼夜使者', description: '左眼为日（金），右眼为月（银），交替闪烁', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '衔烛之龙', description: '完整龙形，口衔巨大蜡烛，威严初现', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '烛龙巡天', description: '周身火焰缭绕，尾部带光尾，云气跟随', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '烛照九幽', description: '光芒穿透九层幽暗，神圣不可侵，身体半透明', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '开目为昼', description: '睁眼时全屏变亮，闭目则暗，时空微动', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '烛龙至尊', description: '日月环绕，呼吸间引动星辰，周身有阴阳二气', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'yinglong', name: '应龙', seriesId: 'myth',
        levels: [
          { level: 1, name: '虺卵', description: '青色蛋壳，表面有云纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳虺', description: '无足小蛇，头有肉角', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '虺行', description: '能蠕动前行，鳞片初生', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '幼蛟', description: '身体变粗，有鳞无角，可腾跃', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '蛟龙', description: '头顶生短角，身带水汽', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '角龙', description: '角长，有四足，能控水', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '应龙雏形', description: '背部生出双翼（小），主风雨', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '应龙', description: '双翼展开，身披金鳞，能飞', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '战龙', description: '翼展增大，爪牙锋利，身带雷光', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '斩蚩尤之姿', description: '黄金战甲，翼展遮天，双目赤金', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '神龙降世', description: '身缠雷电，翼下生风，云海随行', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '太古应龙', description: '身体巨大化，鳞片呈古铜色，一吼震天', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'nine_tail_fox', name: '九尾狐', seriesId: 'myth',
        levels: [
          { level: 1, name: '灵狐蛋', description: '白色蛋壳，有九条浅纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼狐', description: '通体雪白，一尾，耳尖', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '白狐', description: '一尾变长，能小跑', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '二尾狐', description: '尾部分叉为二，额生灵纹', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '三尾灵狐', description: '三尾，灵纹发光，能迷惑', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '四尾', description: '体型增大，毛色泛银', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '五尾', description: '尾巴增多，身绕狐火（蓝）', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '六尾妖狐', description: '六尾，狐火旺盛，眼神魅惑', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '七尾', description: '尾巴蓬松，额生第三眼（闭）', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '八尾', description: '八尾，第三眼微睁，散发圣光', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '九尾天狐', description: '九尾全开，周身神圣，不可侵犯', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '青丘之主', description: '九尾化作九道流光，身体半透明，周围有青丘幻影', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'kunpeng', name: '鲲鹏', seriesId: 'myth',
        levels: [
          { level: 1, name: '鲲苗', description: '小鱼形态，北冥之鳞（深蓝）', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鲲', description: '体型变大，背鳍突出', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '巨鲲', description: '遮天巨鱼，浪涛随身，喷水柱', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '灵鲲', description: '鳍变长，头部生角', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '化鲲', description: '鱼身开始变扁，鳍变为羽状', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '半鹏', description: '鱼身鸟翼，正在蜕变，羽翼未丰', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '鹏雏', description: '双翼成形，能滑翔', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '大鹏', description: '翼若垂天之云，可翱翔', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '鲲鹏', description: '鱼鸟合一，可游可飞', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '鹏程万里', description: '翼展遮天，周身有风云', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '北冥之主', description: '身体环绕极光，游于星河', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '逍遥鲲鹏', description: '化作水汽与云，无形无相，随时凝聚', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'fenghuang', name: '凤凰', seriesId: 'myth',
        levels: [
          { level: 1, name: '火羽蛋', description: '红色蛋壳，有火焰羽毛纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏凤', description: '黄色绒毛，头有冠凸', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼凤', description: '长出红色羽毛，尾羽短', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '青鸾', description: '羽毛变青，尾羽变长', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '赤凤', description: '转为赤红色，尾羽有金斑', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '朱雀', description: '全身火红，周围有火星', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '凤雏', description: '尾羽变长，冠羽立起', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '凤凰', description: '七彩羽毛，尾羽如虹', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '神凤', description: '周身火焰环绕，鸣叫悦耳', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '百鸟之王', description: '尾羽展开如屏，光芒四射', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '涅槃之凤', description: '身体半透明，火焰中带重生之意', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒凤凰', description: '周身金光，每一次展翅撒下星火，不死不灭', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'qilin', name: '麒麟', seriesId: 'myth',
        levels: [
          { level: 1, name: '瑞兽玉', description: '碧绿色玉石，内有祥云', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '玉麒麟幼', description: '玉石裂开，露出小兽头，有角', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼麟', description: '全身青绿，四蹄踏云', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '角麟', description: '角变长，身披鳞甲', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '仁兽', description: '眼神温和，身边有瑞气', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '祥瑞之麟', description: '足下生莲，头顶有光圈', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '踏火麒麟', description: '四蹄带火，鳞片泛金', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '圣麟', description: '身体变白，角呈金色', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '天麟', description: '身周有云霞，能腾云', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '太平麒麟', description: '脚踏祥云，口衔玉书', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '万瑞之祖', description: '全身七彩，角上刻有福纹', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒瑞兽', description: '化作巨大玉雕，护佑一方，福气源源', stage: 'legendary', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 二、宝可梦风格系列（6种） =====
  {
    id: 'pokemon', name: '宝可梦风格', emoji: '⚡',
    species: [
      {
        id: 'charmander', name: '小火龙→喷火龙', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '火纹蛋', description: '红色蛋壳，火焰纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳小龙', description: '露出头，尾尖有火星', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '蹒跚火苗', description: '站立，尾火微弱', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '小火龙', description: '橙红皮肤，尾火燃烧', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '烈火小龙', description: '尾火变大，能吐小火花', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '火恐龙', description: '体型增大，头顶长小角', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '烈火恐龙', description: '角大，爪牙锋利，尾火旺盛', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '喷火龙雏形', description: '背部隆起，翅膀初现', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '喷火龙', description: '双翼展开，能飞行，尾火冲天', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '烈焰喷火龙', description: '全身火焰包围，眼神锐利', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '超级喷火龙X', description: '蓝焰形态，全身暗蓝，霸气', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '始源喷火龙', description: '金色与蓝色交织，翼展遮天，全屏火焰', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'bulbasaur', name: '妙蛙种子→妙蛙花', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '种子蛋', description: '绿色蛋壳，有藤蔓纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳芽', description: '露出小头，背有芽', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼种子', description: '四足爬行，背芽变叶', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '妙蛙种子', description: '背生种子，蹲坐', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '叶芽成长', description: '种子变大，长出藤蔓', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '妙蛙草', description: '种子开花（小花），站立', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '花苞绽放', description: '花朵变大，香气扑鼻', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '妙蛙花', description: '背上巨花盛开，藤蔓环绕', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '繁花之姿', description: '花朵散发光芒，藤蔓延伸', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '超级妙蛙花', description: '巨花盛开，藤蔓如手臂', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '花神', description: '周身花瓣飘落，大地回春', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '远古妙蛙', description: '身体巨大化，花中藏有森林', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'squirtle', name: '杰尼龟→水箭龟', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '水纹蛋', description: '蓝色蛋壳，波浪纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '破壳小龟', description: '露出头，背有小壳', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼龟', description: '能爬行，壳软', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '杰尼龟', description: '龟壳坚硬，能喷水', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '水枪龟', description: '喷水变强，尾巴变长', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '卡咪龟', description: '体型增大，壳上有突起', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '水炮龟', description: '能发水炮，壳更硬', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '水箭龟雏形', description: '背部长出炮管', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '水箭龟', description: '双炮管，能高速水射', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '超级水箭龟', description: '背炮变大，全身水光', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '海神龟', description: '周围有洋流，可召唤浪', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '远古水龟', description: '壳如岛屿，背负山河', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'eevee', name: '伊布', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '伊布蛋', description: '棕色蛋，有狐尾纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼狐', description: '棕色小狐，耳大', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '伊布', description: '标准形态，毛色棕褐，尾蓬松', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '活泼伊布', description: '眼神灵动，毛发光泽', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '进化预备', description: '额头出现宝石（颜色不定）', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '分支初现', description: '身体散发对应元素光', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '元素显现', description: '毛色改变，尾巴变长', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '水伊布', description: '鳞片发光，能化水', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '水伊布·潮汐', description: '周身洋流环绕', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '水伊布·漩涡', description: '可召唤巨大漩涡', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '水伊布·海神', description: '半透明水化身躯', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒水伊布', description: '化作清泉，生生不息', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'pikachu', name: '皮卡丘', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '电气蛋', description: '黄色蛋壳，闪电纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '皮丘', description: '小型，脸颊无电', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '皮丘萌动', description: '能站立，脸颊微电', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '皮卡丘', description: '标准形态，脸颊带电，尾巴闪电状', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '电光鼠', description: '电力增强，奔跑带光', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '雷丘雏形', description: '体型增大，尾巴变粗', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '雷丘', description: '完全形态，电力十足', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '冲浪皮卡丘', description: '尾巴可当冲浪板', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '飞空皮卡丘', description: '能用气球飞行', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '搭档皮卡丘', description: '戴帽子，电光四射', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '极巨化皮卡丘', description: '体型巨大，雷云随行', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '雷电之神', description: '周身落雷，可控制气候', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'riolu', name: '利欧路→路卡利欧', seriesId: 'pokemon',
        levels: [
          { level: 1, name: '波导蛋', description: '蓝色蛋，有波动纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鲁', description: '小兽形，后肢站立', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '利欧路', description: '蓝色，胸前有黑纹', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '波导学徒', description: '能感知波导，眼发光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '练习生', description: '拳脚有力，速度提升', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '进化前夜', description: '身体发光，骨骼强化', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '路卡利欧', description: '人形，全身蓝黑，手有骨刺', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '波导战士', description: '能发波导弹，眼神锐利', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '波导大师', description: '骨刺变长，全身气势', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '超级路卡利欧', description: '肌肉增强，波导爆发', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '传说守护者', description: '全身金色纹路，力量内敛', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '波导之神', description: '无形无质，可化波导，守护众生', stage: 'legendary', requiredScore: 450 },
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
          { level: 2, name: '黑白初现', description: '黑白毛色开始显现，眼周变黑', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '蹒跚学步', description: '能缓慢爬动，抱住竹子', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '小熊猫', description: '黑白分明，圆滚滚，行动笨拙', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '好奇熊猫', description: '开始站立，东张西望', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '竹林顽童', description: '喜欢打滚，抱着竹子玩耍', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '青涩少年', description: '体型增大，毛色亮泽', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '竹林卫士', description: '端坐竹林，神情沉稳', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '成年熊猫', description: '憨态可掬，悠然吃竹', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '竹林尊者', description: '体型硕大，周身有竹叶飘落', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '太极熊猫', description: '黑白形成太极图，气息沉稳', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '竹林圣者', description: '全身散发柔和光芒，竹影随行', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'golden_monkey', name: '金丝猴', seriesId: 'national',
        levels: [
          { level: 1, name: '金丝卵', description: '浅金色蛋，有茸毛', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '灰白幼猴', description: '全身灰白，紧抱树干', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '初显金毛', description: '金色毛发逐渐显现', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年金猴', description: '金毛半覆，活泼好动', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '蓝面初显', description: '面部变蓝，毛发金黄', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '金丝猴', description: '标准形态，长尾，金毛亮丽', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '长尾猴', description: '尾长过身，攀爬如飞', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '壮年金猴', description: '肌肉结实，毛色深金', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金猴王', description: '体型最大，威严，毛发披肩', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '金丝尊者', description: '周身金光，面容威严', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '金丝大圣', description: '手执金棒，身披铠甲', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '齐天金猴', description: '腾云驾雾，毫毛可化万千', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'red_crowned_crane', name: '丹顶鹤', seriesId: 'national',
        levels: [
          { level: 1, name: '鹤卵', description: '灰绿色蛋，有斑点', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏鹤', description: '黄褐色绒毛，蹒跚学步', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '灰羽渐白', description: '羽毛变灰白，开始练习飞翔', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年鹤', description: '白羽初成，头顶微红', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '飞羽初展', description: '能短距离飞翔，姿态优美', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '丹顶鹤', description: '白羽黑颈，头顶朱红', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '长寿鹤', description: '朱红更艳，翅尖有黑羽', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '仙鹤', description: '体型修长，动作优雅', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '祥云伴鹤', description: '飞翔时带云气，声闻于天', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '云中仙鹤', description: '脚踏祥云，头顶有光环', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '羽化鹤', description: '身体半透明，羽翼发光', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '南极仙翁之鹤', description: '周身仙气，可驮仙人，长生不老', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'south_china_tiger', name: '华南虎', seriesId: 'national',
        levels: [
          { level: 1, name: '虎蛋', description: '橙黄色蛋，条纹状', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '虎崽', description: '小虎，条纹初现，奶凶可爱', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼虎', description: '能奔跑，眼睛有神', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年虎', description: '体型渐长，开始独立', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '斑纹加深', description: '条纹变黑，肌肉初显', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '壮年虎', description: '王者姿态，花纹鲜明', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '百兽之王', description: '威风凛凛，目光如炬', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '虎王', description: '额前"王"字清晰，吼声震林', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金虎', description: '毛色变金，条纹暗红', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '白虎', description: '通体雪白，条纹银灰', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '神虎', description: '背生双翼，能飞', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '虎神', description: '周身雷电，行走间山摇地动', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'chinese_alligator', name: '扬子鳄', seriesId: 'national',
        levels: [
          { level: 1, name: '鳄卵', description: '白色蛋，粗糙', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鳄', description: '小鳄鱼，头大身小', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '鳞甲初生', description: '背甲变硬，能游泳', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年鳄', description: '体型增大，尾巴有力', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '水中霸主', description: '游泳迅速，能潜伏', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '扬子鳄', description: '标准形态，暗色鳞甲', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '铁甲鳄', description: '鳞甲变厚，呈盔甲状', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '古鳄', description: '嘴长，牙齿锋利', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '大鳄', description: '身长数米，威严', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '龙鳄', description: '头生角，似龙非龙', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '蛟鳄', description: '身带水汽，可兴风浪', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '扬子龙', description: '化为龙形，镇守长江', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'crested_ibis', name: '朱鹮', seriesId: 'national',
        levels: [
          { level: 1, name: '朱卵', description: '淡粉色蛋', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '灰羽雏鸟', description: '灰色绒毛，喙长', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '白羽渐生', description: '白色羽毛出现', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '朱色初现', description: '头部开始泛红', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '朱鹮', description: '白羽红面，长喙', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '红羽加深', description: '红色更深，翅尖带粉', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '祥瑞之鸟', description: '飞翔时带彩虹', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '仙鸟', description: '体态优雅，步履从容', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '吉祥朱鹮', description: '头顶有金光', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '太阳鸟', description: '周身发暖光', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '凤凰近亲', description: '尾羽变长，五彩', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '涅槃朱鹮', description: '浑身火焰般羽毛，可重生', stage: 'legendary', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 四、科幻机甲系列（6种） =====
  {
    id: 'mecha', name: '科幻机甲', emoji: '🤖',
    species: [
      {
        id: 'mecha_dragon', name: '机械龙', seriesId: 'mecha',
        levels: [
          { level: 1, name: '核心引擎', description: '发光能量球，悬浮旋转', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '骨骼雏形', description: '能量球周围出现金属骨架', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '装甲拼接', description: '覆盖简单钢板，眼睛亮蓝光', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '机械幼龙', description: '四肢成型，尾部有能量刃', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '强化装甲', description: '增加护甲，背上推进器', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '战斗龙形', description: '全身金属光泽，头部有角', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '武装巨龙', description: '肩部导弹舱，尾部加长', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '钢铁巨龙', description: '全身重甲，激光眼', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '飞行模式', description: '背部展开机翼，推进器喷火', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '全装甲龙', description: '纳米装甲，全身发光', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '终焉龙神', description: '金色装甲，全息光翼展开', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '机甲之神', description: '身体半透明，数据流环绕，可变多形态', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'cyber_cat', name: '赛博猫', seriesId: 'mecha',
        levels: [
          { level: 1, name: '电路幼猫', description: '半机械半生物，电子眼闪烁', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '机械爪现', description: '爪尖带光，能吸附', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '数据尾', description: '尾巴透明，显示数字', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '赛博猫', description: '金属外壳，眼睛发光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '全息投影', description: '可投射小型幻象', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '隐身猫', description: '能短暂隐形（透明化）', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '武装猫', description: '肩部加装小炮', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '机甲猫', description: '全身覆盖硬壳', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '全息猫', description: '身体半透明，可分身', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '赛博猫神', description: '全身流光，能操控电子', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '网络主宰', description: '可在数据流中穿行', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '机械圣猫', description: '形态完美，眼如星辰', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'space_mecha', name: '星际机甲', seriesId: 'mecha',
        levels: [
          { level: 1, name: '核心球', description: '银色金属球，发蓝光', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '框架搭建', description: '人形框架，无装甲', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '基础机甲', description: '简单装甲，手持光剑', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星际战士', description: '全副武装，推进器', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '飞行单元', description: '背包飞行器', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '重装型号', description: '加厚装甲，肩炮', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '指挥官型', description: '头戴指挥天线，护盾', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '星际战甲', description: '全身涂装，武器多样', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '精英战甲', description: '金色装饰，能量罩', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '星际霸主', description: '巨型机甲，星空披风', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '宇宙巨人', description: '可变形为战舰', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '星际之神', description: '全身星图，可开辟虫洞', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'quantum_beast', name: '量子兽', seriesId: 'mecha',
        levels: [
          { level: 1, name: '量子云', description: '一团不确定形态的粒子云', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '轮廓初定', description: '云中隐约出现小兽形状', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '半实体', description: '部分透明，闪烁不定', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '量子雏形', description: '形态开始稳定，半透明', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '波粒二象', description: '时而粒子，时而波', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '双态共存', description: '同时出现在两个位置', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '量子纠缠', description: '与周围物体产生联系', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '量子实体', description: '稳定的能量体形态', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '量子增幅', description: '体型增大，能量场显现', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '量子之神', description: '可同时存在于多处', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '多维存在', description: '能穿行不同维度', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '量子永恒', description: '超脱时间，无处不在', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'digital_phoenix', name: '电子凤凰', seriesId: 'mecha',
        levels: [
          { level: 1, name: '数据蛋', description: '由0和1组成的蛋形', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '雏凤数据', description: '雏鸟形态，半透明', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '代码羽', description: '羽毛由代码组成', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '电子凤', description: '全身发光，有像素化特效', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '防火墙翼', description: '翼展为光栅', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '加密凤', description: '身体呈锁状纹路', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '网络凤', description: '可穿梭网络', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '全息凤凰', description: '三维投影，色彩斑斓', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '人工智能凤', description: '拥有智慧，可交流', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '数字神凤', description: '控制信息流', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '元宇宙之主', description: '可构筑虚拟世界', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒数据', description: '化作信息流，永存网络', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'mecha_shark', name: '机械鲨', seriesId: 'mecha',
        levels: [
          { level: 1, name: '金属卵', description: '银灰色，有散热孔', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鲨', description: '小型，机械骨骼可见', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '装甲鲨', description: '覆盖合金皮肤', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '鱼雷鲨', description: '流线型，速度极快', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '武装鲨', description: '背部加装导弹发射器', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '潜水鲨', description: '可深潜，抗压', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '电子鲨', description: '声呐探测，发光诱饵', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '巨鲨', description: '体型巨大', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '航母鲨', description: '背部可停放小飞行器', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '海洋主宰', description: '控制水域', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '纳米鲨', description: '可分解重组', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '机械海神', description: '全身能量罩，统帅海洋', stage: 'legendary', requiredScore: 450 },
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
          { level: 2, name: '萌芽小马', description: '光点拉长成小马轮廓', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '白色小驹', description: '纯白身体，额头微凸', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '幼年独角兽', description: '长出小角，鬃毛发光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '角光初现', description: '角尖发出柔和荧光', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '银鬃飞马', description: '鬃毛银色，能短途漂浮', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '彩虹鬃毛', description: '鬃毛七彩渐变，角更亮', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '圣光独角兽', description: '全身圣光，蹄下生花', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星夜独角兽', description: '鬃毛如星空，角尖有星环', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '梦境守护者', description: '半透明，能入梦境', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒独角兽', description: '金色鬃毛，角刻符文', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '创世独角兽', description: '全身星光，角能开辟空间', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'wyvern', name: '飞龙', seriesId: 'magic',
        levels: [
          { level: 1, name: '龙蛋', description: '暗红色，有鳞纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼龙', description: '无翼，四足，能喷火', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '翼芽', description: '背部出现小翼芽', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '翼龙', description: '翼长成，能滑翔', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '火焰龙', description: '喷火强，翼大', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '风龙', description: '可控制气流', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '雷龙', description: '带雷电', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '冰龙', description: '喷冰霜', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '古龙', description: '体型巨大，皮糙', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '龙王', description: '身披铠甲，号令群龙', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '太古龙', description: '智慧极高，魔法强大', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '龙神', description: '半神，掌控元素', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'fairy', name: '精灵', seriesId: 'magic',
        levels: [
          { level: 1, name: '花苞', description: '含苞待放', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '花精灵', description: '小花形，有透明翅', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '叶精灵', description: '绿色，能控植物', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '光精灵', description: '发光，能治愈', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '暗精灵', description: '暗色，速度极快', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '水精灵', description: '蓝色，控水', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '火精灵', description: '红色，控火', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '风精灵', description: '白色，隐身', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '大地精灵', description: '褐色，控制土石', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '自然精灵', description: '融合元素', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '森林之子', description: '与自然一体', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '精灵王', description: '号令自然，永生', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'treant', name: '树人', seriesId: 'magic',
        levels: [
          { level: 1, name: '种子', description: '橡果状', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '嫩芽', description: '破土而出', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '小树苗', description: '长出叶子', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '树童', description: '有简单五官', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '树人', description: '四肢为树枝', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '橡树人', description: '身硬，能投石', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '古树人', description: '身缠藤蔓，苔藓', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '战争树人', description: '可拔根移动', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '森林卫士', description: '高大，持木盾', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '智慧古树', description: '千年智慧', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '世界树幼苗', description: '连接天地', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒古树', description: '化为森林之灵', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'griffin', name: '狮鹫', seriesId: 'magic',
        levels: [
          { level: 1, name: '狮鹫蛋', description: '金色，有爪印纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼狮鹫', description: '狮身鹰头，小翼', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '羽翼初长', description: '翼变大', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '狮鹫', description: '标准，能飞', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '金翼狮鹫', description: '翼金色，能闪光', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '银鬃狮鹫', description: '鬃毛银色', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '皇家狮鹫', description: '戴王冠', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '风暴狮鹫', description: '控制天气', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '雷光狮鹫', description: '带电', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '格里芬王', description: '统领鸟兽', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '神圣狮鹫', description: '圣光环绕', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '天空霸主', description: '遮天蔽日', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'mermaid', name: '美人鱼', seriesId: 'magic',
        levels: [
          { level: 1, name: '珍珠', description: '内含人影', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼鱼人', description: '上半身人，尾为鱼', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '鳞光', description: '鳞片发光', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '人鱼', description: '长发，歌声动听', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '海妖', description: '魅惑歌声', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '深海人鱼', description: '适应高压', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '珊瑚公主', description: '头戴珊瑚冠', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '海洋之女', description: '控水，召唤海兽', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '人鱼女王', description: '威严，珍珠项链', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '海神之女', description: '半神，掌控洋流', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '大洋之魂', description: '化为海水', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒人鱼', description: '守护海洋', stage: 'legendary', requiredScore: 450 },
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
          { level: 2, name: '破壳小爪', description: '伸出爪子，带鳞', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年暴龙', description: '全身鳞片，短尾', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '少年暴龙', description: '前肢短，后肢强', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '追赶者', description: '能快速奔跑', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '森林霸主', description: '体型大，皮肤深', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '领地之王', description: '头顶角质冠', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '巨型暴龙', description: '身长数米，牙尖', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '白垩纪领主', description: '厚重鳞甲', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '暴龙之王', description: '背有棘刺，尾有力', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '终极掠食者', description: '青铜色，坚硬', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '恐龙至尊', description: '全身远古气息，吼震屏幕', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'triceratops', name: '三角龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '角蛋', description: '有三个凸点', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼角龙', description: '头部有三个小角', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '头盾初现', description: '头后出现骨盾', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '三角龙', description: '角变长，盾大', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '盾甲龙', description: '盾更厚，能防御', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '冲锋龙', description: '能奔跑冲撞', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '铁头龙', description: '角如铁', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '巨三角', description: '体型巨大', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '群首领', description: '带领族群', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '古三角', description: '角呈金色', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '化石之王', description: '全身硬化', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒三角', description: '化为山石，守护大地', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'pterosaur', name: '翼龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '飞卵', description: '带翅膜纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '翼苗', description: '小翼，会滑翔', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '翼手龙', description: '爪尖，翼展宽', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '飞翔者', description: '能盘旋', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '风神翼龙', description: '巨大，能翱翔九天', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '雷电翼龙', description: '带静电', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '火焰翼龙', description: '喷火', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '寒冰翼龙', description: '冰冻', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '王者翼龙', description: '统领天空', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '天空领主', description: '遮天', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '太古翼龙', description: '化石复活', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒之翼', description: '化为风', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'mammoth', name: '猛犸象', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '毛蛋', description: '长毛覆盖', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '小猛犸', description: '长毛，象牙初露', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '象牙变长', description: '牙突出', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '猛犸象', description: '标准，长毛', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '冰原猛犸', description: '耐寒，身有冰晶', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨猛犸', description: '体型巨大', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '象牙王', description: '牙如玉', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '古猛犸', description: '身披厚毛', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '冰川之主', description: '行走于冰', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '神猛犸', description: '周身寒气', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '远古巨兽', description: '如移动冰山', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒猛犸', description: '化为冰雕，守护冰原', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'sabertooth', name: '剑齿虎', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '虎卵', description: '剑齿纹', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼剑齿', description: '上犬齿突出', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '长牙初成', description: '牙变长', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '剑齿虎', description: '标准，牙如剑', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '刃牙虎', description: '牙更锋利', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨剑虎', description: '体型大', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '冰河猎手', description: '适应寒冷', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '王剑齿', description: '统治草原', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '金剑齿', description: '牙金色', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '神剑虎', description: '身体发光', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '太古剑齿', description: '化石之力', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒剑齿', description: '化为宝石', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'mosasaur', name: '沧龙', seriesId: 'prehistoric',
        levels: [
          { level: 1, name: '水卵', description: '有鳞片', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '幼沧龙', description: '小型，水性好', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '鳞甲变厚', description: '能潜水', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '沧龙', description: '标准，长尾', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '海中霸主', description: '速度极快', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '巨沧', description: '体型巨大', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '王沧', description: '头有冠', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '深海沧龙', description: '适应深水', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '古沧', description: '身带化石纹', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '海王沧龙', description: '控制海洋', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '神沧', description: '半透明，神力', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒沧龙', description: '化为海流', stage: 'legendary', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 七、星座守护系列（12种） =====
  {
    id: 'constellation', name: '星座守护', emoji: '♈',
    species: [
      {
        id: 'aries', name: '白羊座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有公羊轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '公羊形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小公羊，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，毛皮带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'taurus', name: '金牛座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有公牛轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '公牛形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小公牛，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，毛皮带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'gemini', name: '双子座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有双子轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '双子形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小双子灵体，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，灵体带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'cancer', name: '巨蟹座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有螃蟹轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '螃蟹形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小螃蟹，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，甲壳带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'leo', name: '狮子座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有雄狮轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '雄狮形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小狮子，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，鬃毛带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'virgo', name: '处女座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有少女轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '少女形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小少女灵体，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，衣裙带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'libra', name: '天秤座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有天秤轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '天秤形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小天秤，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，天秤带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'scorpio', name: '天蝎座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有蝎子轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '蝎子形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小蝎子，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，甲壳带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'sagittarius', name: '射手座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有人马轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '人马形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小人马，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，人马带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'capricorn', name: '摩羯座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有羊身鱼尾轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '摩羯形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小摩羯，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，鳞片带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'aquarius', name: '水瓶座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有持瓶者轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '持瓶者形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小持瓶者，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，水瓶带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'pisces', name: '双鱼座', seriesId: 'constellation',
        levels: [
          { level: 1, name: '星辰幼崽', description: '小星星，内有双鱼轮廓', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '星光显现', description: '双鱼形态发光', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '幼年守护', description: '小双鱼，有星尘围绕', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '星辉之姿', description: '体型增大，鳞片带星点', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '星座初成', description: '身体出现星座连线', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '黄道之影', description: '脚踏黄道光带', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '英勇之态', description: '眼神坚定，光环加身', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '守护者', description: '身穿星甲', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '星宫之主', description: '周围有十二宫星图', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '辉煌星座', description: '全身化为星座图', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒星光', description: '身体半透明，光芒万丈', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '黄道之神', description: '化为星座，永恒照耀', stage: 'legendary', requiredScore: 450 },
        ],
      },
    ],
  },

  // ===== 八、民间传说系列（6种） =====
  {
    id: 'folklore', name: '民间传说', emoji: '🏮',
    species: [
      {
        id: 'nian', name: '年兽', seriesId: 'folklore',
        levels: [
          { level: 1, name: '红纸团', description: '红色纸团，内有金光', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '小兽探头', description: '露出毛茸茸头部，有角', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '红色幼兽', description: '全身红色，四蹄踏火', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '长角年兽', description: '角变长，身体披鳞', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '春节使者', description: '身绕鞭炮，带来喜庆', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '驱邪灵兽', description: '双眼放光，可震慑', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '瑞兽临门', description: '体型增大，身边祥云', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '吉祥年兽', description: '全身中国结，福字环绕', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '福星高照', description: '头顶有金色福光', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '太平神兽', description: '脚踏铜钱，身披红绸', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '万象更新', description: '周身四季花卉轮流开放', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒瑞兽', description: '化作福字高悬，护佑全班', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'dragon_boat', name: '龙舟精灵', seriesId: 'folklore',
        levels: [
          { level: 1, name: '船桨芽', description: '小木浆形状', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '龙头显形', description: '龙首浮现', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '龙舟幼体', description: '龙头船身，可滑行', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '龙舟', description: '彩绘龙身，旗帜', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '竞渡龙舟', description: '速度提升，水花四溅', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '红金龙舟', description: '红金配色', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '绿玉龙舟', description: '碧绿，有玉鳞', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '王者龙舟', description: '头戴金冠', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '飞天龙舟', description: '可腾空', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '神龙舟', description: '化为龙形', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '四海龙舟', description: '可在湖海穿行', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒龙舟', description: '化作彩虹桥', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'lantern', name: '灯笼精', seriesId: 'folklore',
        levels: [
          { level: 1, name: '纸灯笼', description: '白纸糊', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '烛光内燃', description: '点亮', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '萌眼现', description: '出现眼睛', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '小灯笼精', description: '会漂浮，发光', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '节庆灯笼', description: '红色，有穗', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '福字灯笼', description: '印福字', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '祥云灯笼', description: '带云纹', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '龙凤灯笼', description: '双面雕龙画凤', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '天灯', description: '能升空', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '宫灯', description: '精致，多角', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '神灯', description: '可许愿', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒明灯', description: '永不熄灭，照亮前程', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'dumpling', name: '饺子福娃', seriesId: 'folklore',
        levels: [
          { level: 1, name: '面团', description: '白色圆团', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '褶子出现', description: '有花边', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '笑脸', description: '出现笑脸', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '小福娃', description: '包子形，红脸蛋', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '元宝饺', description: '元宝形', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '彩饺', description: '彩色面皮', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '金饺', description: '金色', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '银饺', description: '银色', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '福禄饺', description: '内有硬币（装饰）', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '团圆饺', description: '大团圆', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '如意饺', description: '镶如意纹', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '永恒福娃', description: '化为福神，赐福全班', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'fu_star', name: '福星', seriesId: 'folklore',
        levels: [
          { level: 1, name: '福光', description: '一团红光', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '福字初现', description: '浮现"福"字', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '童子形', description: '小小童子', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '福星童', description: '手持福牌', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '赐福童子', description: '撒福气', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '福星高照', description: '头顶光环', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '大福星', description: '身披红袍', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '天官赐福', description: '手持如意', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '福神', description: '威严慈祥', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '万福之主', description: '周身福字环绕', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒福星', description: '化为福气', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '福满天下', description: '福泽全班', stage: 'legendary', requiredScore: 450 },
        ],
      },
      {
        id: 'shou_star', name: '寿星', seriesId: 'folklore',
        levels: [
          { level: 1, name: '寿桃', description: '粉色桃形', stage: 'egg', requiredScore: 0 },
          { level: 2, name: '桃生笑', description: '出现笑脸', stage: 'baby', requiredScore: 15 },
          { level: 3, name: '小寿童', description: '童子持桃', stage: 'baby', requiredScore: 35 },
          { level: 4, name: '寿星童', description: '头大，拄拐', stage: 'baby', requiredScore: 60 },
          { level: 5, name: '赐寿童子', description: '撒寿桃', stage: 'growing', requiredScore: 90 },
          { level: 6, name: '寿星', description: '白须，额高', stage: 'growing', requiredScore: 125 },
          { level: 7, name: '南极仙翁', description: '骑鹿', stage: 'growing', requiredScore: 165 },
          { level: 8, name: '长寿星', description: '手捧寿桃', stage: 'mature', requiredScore: 210 },
          { level: 9, name: '万寿神', description: '周身寿字', stage: 'mature', requiredScore: 260 },
          { level: 10, name: '寿与天齐', description: '紫气东来', stage: 'mature', requiredScore: 315 },
          { level: 11, name: '永恒寿星', description: '化为仙鹤', stage: 'legendary', requiredScore: 375 },
          { level: 12, name: '寿比南山', description: '化为南山松柏', stage: 'legendary', requiredScore: 450 },
        ],
      },
    ],
  },
]

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
