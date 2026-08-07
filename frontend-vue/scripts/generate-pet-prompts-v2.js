#!/usr/bin/env node
/**
 * 生成「宠物 AI 生图提示词 v2」—— 全 125 物种 × 6 阶段精细提示词文档。
 *
 * 相比 v1（pet-prompts.csv 通用模板），v2 的改进：
 *  1. 系列专属视觉规范：10 系列各自风格 / 质感 / 六阶配色 / 六阶阶段演绎场景
 *  2. 注入角色形象：petProfiles.ts 的 form（形态）+ symbol（核心意象）
 *  3. 中英双语：中文主句（即梦/通义）+ 英文风格词（Midjourney/DALL·E）
 *
 * 数据契约：docs/宠物系统全案.md 6.4 节 + docs/pet-image-manifest.md
 * 用法：node scripts/generate-pet-prompts-v2.js [输出路径]（默认 docs/宠物AI生图提示词.md）
 */
import fs from 'fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const srcDir = path.join(__dirname, '../src/utils')
const src = fs.readFileSync(path.join(srcDir, 'petData.ts'), 'utf8')
const ext = fs.readFileSync(path.join(srcDir, 'petDataExtended.ts'), 'utf8')
const all = src + '\n' + ext
const profileSrc = fs.readFileSync(path.join(srcDir, 'petProfiles.ts'), 'utf8')

// ===== 六阶段（通用框架，startLv 与 petData 等级映射一致）=====
const STAGES = [
  { key: 'egg', zh: '灵胎初醒', startLv: 1 },
  { key: 'baby', zh: '凡尘砺心', startLv: 3 },
  { key: 'growing', zh: '道法初成', startLv: 5 },
  { key: 'mature', zh: '大劫淬炼', startLv: 7 },
  { key: 'legendary', zh: '封神登天', startLv: 9 },
  { key: 'transcendent', zh: '道果圆满', startLv: 11 },
]

// ===== 系列视觉规范 =====
// colors[6] 对应六阶 [主色{name,hex}, 辅色{name,hex}]；scenes[6]/sceneEn[6] 对应六阶阶段演绎
const SERIES = {
  myth: {
    name: '山海经',
    styleZh: '新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘',
    styleEn: 'Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture',
    scenes: [
      '混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘',
      '幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则',
      '异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量',
      '上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变',
      '神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方',
      '图腾显圣，山海为印，洪荒万古唯此一尊',
    ],
    sceneEn: ['nascent chaos', 'wandering ancient wilds', 'awakening power', 'trial by elements', 'divine coronation', 'eternal harmony'],
    colors: [
      [{ name: '混沌玄青', hex: '#2A2A35' }, { name: '青铜褐', hex: '#8B7E6A' }],
      [{ name: '石青灰', hex: '#4A5568' }, { name: '赭石棕', hex: '#A0855B' }],
      [{ name: '朱砂红', hex: '#8B2E2E' }, { name: '鎏金', hex: '#D4AF37' }],
      [{ name: '玄墨黑', hex: '#1A0A0A' }, { name: '血赤', hex: '#8B1A1A' }],
      [{ name: '鎏金', hex: '#D4AF37' }, { name: '朱紫', hex: '#7B2E8B' }],
      [{ name: '云白', hex: '#F5F0E8' }, { name: '淡金', hex: '#C9B37E' }],
    ],
  },
  dongfang: {
    name: '东方神话',
    styleZh: '新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长',
    styleEn: 'xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura',
    scenes: [
      '仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘',
      '离开洞府行走人间，红尘历练，磨砺道心',
      '道法初成，法宝显威，法相庄严初现',
      '天劫淬炼、心魔试炼，仙体历经大劫而不灭',
      '功行圆满封神登天，位列仙班，金光紫气环绕',
      '道果圆满，紫气东来，万法皆通证道果',
    ],
    sceneEn: ['spiritual genesis', 'mortal wandering', 'dao awakening', 'heavenly tribulation', 'ascension to immortals', 'union with dao'],
    colors: [
      [{ name: '玄青', hex: '#3A4A5A' }, { name: '仙金', hex: '#C9B37E' }],
      [{ name: '黛青', hex: '#5B6B7A' }, { name: '灰白', hex: '#8A9BA8' }],
      [{ name: '朱砂', hex: '#B8502E' }, { name: '明黄', hex: '#E8B84B' }],
      [{ name: '墨褐', hex: '#241A14' }, { name: '血绛', hex: '#8B1A1A' }],
      [{ name: '紫金', hex: '#5B2A8A' }, { name: '御金', hex: '#F5D742' }],
      [{ name: '素白', hex: '#F5F0E8' }, { name: '淡紫', hex: '#9A8FB8' }],
    ],
  },
  pokemon: {
    name: '宝可梦',
    styleZh: '经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈',
    styleEn: 'classic Japanese anime adventure style, cel-shaded, vibrant and wholesome',
    scenes: [
      '蛋壳微光，新生命初醒，温暖亲近的孵化氛围',
      '蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮',
      '初次对战或训练，火花四溅，能力初显的高光时刻',
      '进化之光笼罩，身体在光中蜕变，挣扎又坚定',
      '完全进化形态，可靠的战斗伙伴，自信昂扬',
      '超进化或极致形态，圣光守护，散发传说级气场',
    ],
    sceneEn: ['hatching glow', 'first steps', 'first battle spark', 'evolution light', 'fully evolved', 'legendary aura'],
    colors: [
      [{ name: '奶油', hex: '#F5E6D3' }, { name: '嫩粉', hex: '#FFB6C1' }],
      [{ name: '晴蓝', hex: '#87CEEB' }, { name: '暖橙', hex: '#FFA07A' }],
      [{ name: '活力橙', hex: '#FF8C42' }, { name: '电光蓝', hex: '#4FC3F7' }],
      [{ name: '深海蓝', hex: '#283593' }, { name: '战斗红', hex: '#E53935' }],
      [{ name: '伙伴金', hex: '#FFD54F' }, { name: '烈焰橙', hex: '#FF7043' }],
      [{ name: '圣光白', hex: '#FFFFFF' }, { name: '彩虹金', hex: '#FFD700' }],
    ],
  },
  digimon: {
    name: '数码宝贝',
    styleZh: '日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻',
    styleEn: 'Japanese digital evolution anime, cel shading with digital particles, techno glow',
    scenes: [
      '数码蛋悬浮微光，数据流环绕，崭新世界初启',
      '幼年期毛茸茸的小数码兽，在数码世界好奇探索',
      '成长进化，数据粒子环绕升腾，光芒凝聚',
      '面对黑暗深渊的试炼，进化之光与黑暗对峙',
      '究极进化，圣光加冕，数据粒子凝成铠甲',
      '神圣形态，希望之光普照，数据洪流归于一念',
    ],
    sceneEn: ['digital egg', 'fresh/child era', 'evolution spark', 'dark trial', 'ultimate evolution', 'holy transcendence'],
    colors: [
      [{ name: '数据浅蓝', hex: '#C9E4FF' }, { name: '数码银', hex: '#E0E0E0' }],
      [{ name: '数码青', hex: '#7EC8E3' }, { name: '淡蓝', hex: '#A8D8EA' }],
      [{ name: '进化蓝', hex: '#3F8EFC' }, { name: '勇气橙', hex: '#FF9E3D' }],
      [{ name: '深渊紫', hex: '#2C2C54' }, { name: '战斗红', hex: '#FF3B30' }],
      [{ name: '圣金', hex: '#FFD700' }, { name: '纯白', hex: '#FFFFFF' }],
      [{ name: '圣光白', hex: '#FFFFFF' }, { name: '虹彩', hex: '#C8A2FF' }],
    ],
  },
  national: {
    name: '国宝',
    styleZh: '皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖',
    styleEn: 'Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop',
    scenes: [
      '新生命诞生于巢穴或洞穴，温暖微光，柔软新生',
      '幼崽蹒跚学步，跟随父母觅食，憨态可掬',
      '青年期活力四射，嬉戏打闹，探索栖息地',
      '成年期沉稳担当，守护领地与族群，历经风雨',
      '成为族群的传奇，祥瑞光环，被守护与敬仰',
      '生态图腾显圣，祥瑞护世，生生不息',
    ],
    sceneEn: ['newborn glow', 'wobbly cub', 'playful youth', 'guardian prime', 'legend of the species', 'guardian spirit'],
    colors: [
      [{ name: '巢暖米', hex: '#F5E9DC' }, { name: '自然灰绿', hex: '#A8B8B0' }],
      [{ name: '幼崽棕', hex: '#C9B79C' }, { name: '林绿', hex: '#8FA8A0' }],
      [{ name: '竹林绿', hex: '#6B8E7A' }, { name: '暖阳黄', hex: '#E8B84B' }],
      [{ name: '深林绿', hex: '#3E5C4B' }, { name: '熟褐', hex: '#C45A3C' }],
      [{ name: '祥瑞金', hex: '#D4AF37' }, { name: '王绿', hex: '#2F4F4F' }],
      [{ name: '云白', hex: '#F5F0E8' }, { name: '生态金绿', hex: '#A3C1AD' }],
    ],
  },
  magic: {
    name: '魔法奇幻',
    styleZh: '欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗',
    styleEn: 'western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere',
    scenes: [
      '魔法光球或种子，奥秘初现，星芒闪烁的神秘新生',
      '初学魔法，好奇探索，小小的魔力火花',
      '魔力觉醒，脚下浮现法阵，魔法光辉绽放',
      '黑暗试炼、魔力对决，在深渊边缘淬炼',
      '大法师或神话生物完全体，威震大陆，法阵与冠冕加身',
      '元素法相圆满，光与秘术铸就神格，威震大陆',
    ],
    sceneEn: ['arcane genesis', 'curious novice', 'magic circle awakening', 'abyssal trial', 'archmage glory', 'elemental oneness'],
    colors: [
      [{ name: '奥秘紫', hex: '#7B68EE' }, { name: '魔法薰衣草', hex: '#E6E6FA' }],
      [{ name: '魔紫', hex: '#9370DB' }, { name: '星光蓝', hex: '#87CEEB' }],
      [{ name: '深魔紫', hex: '#4B0082' }, { name: '咒文金', hex: '#FFD700' }],
      [{ name: '暗夜紫', hex: '#1A0A2E' }, { name: '火焰红', hex: '#FF4500' }],
      [{ name: '传奇紫', hex: '#8B00FF' }, { name: '王冠金', hex: '#FFD700' }],
      [{ name: '永恒白', hex: '#F5F0E8' }, { name: '圣光紫', hex: '#C8A2FF' }],
    ],
  },
  prehistoric: {
    name: '史前生物',
    styleZh: 'BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫',
    styleEn: 'paleoart documentary style with cartoon charm, realistic scale/skin texture',
    scenes: [
      '化石在远古晨光中苏醒，蛋壳裂开，远古大地气息',
      '幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存',
      '成长中的捕猎或迁徙，展露力量与速度',
      '严酷的冰河或火山环境下生存试炼，与天争命',
      '成为领地霸主，王者的姿态与威严',
      '远古霸主图腾显圣，冰川之巅威压万古',
    ],
    sceneEn: ['fossil dawn', 'wandering calf', 'growing hunter', 'survival trial', 'apex predator', 'ancient totem'],
    colors: [
      [{ name: '化石米', hex: '#D2B48C' }, { name: '石灰', hex: '#A9A9A9' }],
      [{ name: '泥土棕', hex: '#8B7355' }, { name: '苔藓绿', hex: '#6B8E7A' }],
      [{ name: '岩金', hex: '#B8860B' }, { name: '沙棕', hex: '#CD853F' }],
      [{ name: '火山黑', hex: '#3E2723' }, { name: '熔岩红', hex: '#E53935' }],
      [{ name: '霸主棕', hex: '#4E342E' }, { name: '王金', hex: '#FFB300' }],
      [{ name: '冰川白', hex: '#F5F0E8' }, { name: '远古蓝', hex: '#87CEEB' }],
    ],
  },
  constellation: {
    name: '星座守护',
    styleZh: '华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高',
    styleEn: 'golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs',
    scenes: [
      '星屑微光中孕育，小宇宙初醒，命定的星芒',
      '圣斗士修行岁月，拳法初练，热血汗水',
      '第七感觉醒，星座之光闪耀，必杀技初现',
      '圣战死斗，伤痕与信念，在绝境中燃烧小宇宙',
      '身着黄金圣衣，镇守宫门，星座图腾全开',
      '第八感领悟，神之领域，小宇宙化作无尽星光',
    ],
    sceneEn: ['cosmic birth', 'saint training', 'seventh sense', 'holy war trial', 'golden cloth', 'eighth sense transcendence'],
    colors: [
      [{ name: '星屑银', hex: '#E8E8F0' }, { name: '白银', hex: '#C0C0C0' }],
      [{ name: '青铜金', hex: '#D4AF37' }, { name: '银灰', hex: '#E0E0E0' }],
      [{ name: '黄金', hex: '#FFD700' }, { name: '星座蓝', hex: '#87CEEB' }],
      [{ name: '暗紫', hex: '#4B0082' }, { name: '圣战红', hex: '#DC143C' }],
      [{ name: '黄金圣衣', hex: '#FFD700' }, { name: '圣光白', hex: '#FFFFFF' }],
      [{ name: '神光白', hex: '#FFFFFF' }, { name: '小宇宙金', hex: '#E8C87A' }],
    ],
  },
  festival: {
    name: '传统节日',
    styleZh: '民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈',
    styleEn: 'Chinese folk festival mascot illustration, nianhua palette, round and festive',
    scenes: [
      '原料初成，厨房温暖，节日的前奏与期待',
      '制作过程中的形态，蒸煮捏塑，巧手匠心',
      '节庆登场，寓意成形，喜庆氛围渐浓',
      '阖家团圆的时刻，热气腾腾，温暖满溢',
      '成为节日主角与名品，家家户户的吉祥符号',
      '化为民俗文化图腾，代代相传，温暖永续',
    ],
    sceneEn: ['festive prelude', 'crafting warmth', 'festival debut', 'family reunion', 'iconic festive food', 'cultural totem'],
    colors: [
      [{ name: '原料暖米', hex: '#F5E6D3' }, { name: '糯粉', hex: '#F8B195' }],
      [{ name: '制作棕', hex: '#E8A87C' }, { name: '喜庆红', hex: '#F67280' }],
      [{ name: '节日红', hex: '#F67280' }, { name: '麦金', hex: '#FFD54F' }],
      [{ name: '团圆红', hex: '#C0392B' }, { name: '蒸金', hex: '#F39C12' }],
      [{ name: '主角红', hex: '#E74C3C' }, { name: '名品金', hex: '#FFD700' }],
      [{ name: '民俗暖白', hex: '#FFF8E7' }, { name: '传承金', hex: '#FFC93C' }],
    ],
  },
  qixia: {
    name: '虹猫蓝兔七侠传',
    styleZh: '2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇',
    styleEn: '2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop',
    scenes: [
      '剑意种子萌芽，江湖初闻，少年侠气的微光',
      '习武入门，剑招初练，磨砺锋芒',
      '剑法初成，行走江湖，惩恶扬善',
      '七剑合璧，魔教大战，刀光剑影中淬炼',
      '七侠之名扬天下，快意恩仇，豪情万丈',
      '剑道圆满，一剑开天，守护苍生的宗师气象',
    ],
    sceneEn: ['sword spirit', 'martial training', 'journeying hero', 'seven swords battle', 'legendary hero', 'master of the blade'],
    colors: [
      [{ name: '剑意青', hex: '#A8D8EA' }, { name: '银灰', hex: '#E0E0E0' }],
      [{ name: '侠义蓝', hex: '#4FC3F7' }, { name: '热血红', hex: '#FF8A80' }],
      [{ name: '江湖橙', hex: '#FF7043' }, { name: '侠气蓝', hex: '#29B6F6' }],
      [{ name: '魔教紫', hex: '#4A148C' }, { name: '血战红', hex: '#E53935' }],
      [{ name: '七剑金', hex: '#FFD700' }, { name: '侠名白', hex: '#FFFFFF' }],
      [{ name: '剑道白', hex: '#F5F0E8' }, { name: '圆满蓝', hex: '#90CAF9' }],
    ],
  },
}

// ===== 解析数据 =====
const seriesName = {}
for (const m of all.matchAll(/id: '([a-z_]+)', name: '([^']+)'/g)) seriesName[m[1]] = m[2]

function parseSpecies(block) {
  const out = []
  const re = /id: '([a-z_]+)', name: '([^']+)', seriesId: '([a-z_]+)',\s*levels: \[([\s\S]*?)\n\s*\],/g
  let m
  while ((m = re.exec(block))) {
    const [, id, name, seriesId, levelsStr] = m
    const lv = {}
    for (const lm of levelsStr.matchAll(/\{ level: (\d+), name: '([^']+)'/g)) lv[+lm[1]] = lm[2]
    for (const lm of levelsStr.matchAll(/Lv\((\d+), '([^']+)'/g)) lv[+lm[1]] = lm[2]
    out.push({ id, name, seriesId, lv })
  }
  return out
}
const species = [...parseSpecies(src), ...parseSpecies(ext)]
const seen = new Set()
const unique = species.filter(s => (seen.has(s.id) ? false : (seen.add(s.id), true)))

const profiles = {}
const blockRe = /^\s{2}([a-z_0-9]+):\s*\{([\s\S]*?)\},?$/gm
let bm
while ((bm = blockRe.exec(profileSrc))) {
  const key = bm[1]
  const body = bm[2]
  const get = f => { const mm = body.match(new RegExp('\\b' + f + ":\\s*'([^']*)'")); return mm ? mm[1] : '' }
  profiles[key] = { form: get('form'), symbol: get('symbol'), movement: get('movement') }
}

// petRefine.ts —— 六阶段精修层（神态/动作/衣着/梳造）
const refineSrc = fs.readFileSync(path.join(srcDir, 'petRefine.ts'), 'utf8')
const refine = (() => {
  const cats = {}
  for (const m of refineSrc.matchAll(/([a-z_0-9]+): '([a-z_]+)',/g)) cats[m[1]] = m[2]
  const base = {}
  const special = {}
  const blockRe = /^  ([a-z_0-9]+): \[\n([\s\S]*?)\n  \],/gm
  let m
  while ((m = blockRe.exec(refineSrc))) {
    const key = m[1]
    const body = m[2]
    const full = []
    const attire = []
    for (const line of body.split('\n')) {
      const f = line.match(/spirit: '([^']*)'.*?action: '([^']*)'.*?attire: '([^']*)'.*?hair: '([^']*)'/)
      if (f) full.push({ spirit: f[1], action: f[2], attire: f[3], hair: f[4] })
      const a = line.match(/^\s*\{ attire: '([^']*)', hair: '([^']*)'\s*\}\s*,?$/)
      if (a) attire.push({ attire: a[1], hair: a[2] })
    }
    if (full.length === 6) base[key] = full
    else if (attire.length === 6) special[key] = attire
  }
  return { cats, base, special }
})()

// petInitial.ts —— 初始形态（灵胎初醒 Lv1）设计表
const initSrc = fs.readFileSync(path.join(srcDir, 'petInitial.ts'), 'utf8')
const initials = {}
const initBlockRe = /^\s{2}([a-z_0-9]+):\s*\{([\s\S]*?)\},?$/gm
let ibm
while ((ibm = initBlockRe.exec(initSrc))) {
  const body = ibm[2]
  const get = f => { const mm = body.match(new RegExp('\\b' + f + ":\\s*'([^']*)'")); return mm ? mm[1] : '' }
  initials[ibm[1]] = { cat: get('cat'), elem: get('elem'), name: get('name'), desc: get('desc') }
}

// petLifeStories.ts —— 六阶段人生档案（有档案的角色 → 提示词走「人生叙事」版，与手绘 SVG / 档案完全对齐）
const lifeSrc = fs.readFileSync(path.join(srcDir, 'petLifeStories.ts'), 'utf8')
const lifeStories = {}
{
  const lifeBlockRe = /\n  ([a-z_0-9]+): \{([\s\S]*?)\n  \},/g
  let lm
  while ((lm = lifeBlockRe.exec(lifeSrc))) {
    const body = lm[2]
    const stagesStr = (body.match(/stages: \[([\s\S]*?)\n\s*\],/) || [])[1] || ''
    const stages = []
    const stageRe = /\{([\s\S]*?)\}/g
    let sm
    while ((sm = stageRe.exec(stagesStr))) {
      const sb = sm[1]
      const f = k => {
        const mm = sb.match(new RegExp('\\b' + k + ":\\s*'((?:[^'\\\\]|\\\\.)*)'"))
        return mm ? mm[1].replace(/\\'/g, "'") : ''
      }
      stages.push({
        name: f('name'), age: f('age'), keyword: f('keyword'),
        character: f('character'), action: f('action'), attire: f('attire'),
        technique: f('technique'), line: f('line'), poem: f('poem'), svg: f('svg'), mood: f('mood'),
        body: f('body'), face: f('face'), actionSeq: f('actionSeq'), attireDetail: f('attireDetail'),
        hairStyle: f('hairStyle'), weaponAction: f('weaponAction'), powerEffect: f('powerEffect'),
      })
    }
    const themeM = body.match(/\btheme: '((?:[^'\\]|\\.)*)'/)
    lifeStories[lm[1]] = { theme: themeM ? themeM[1].replace(/\\'/g, "'") : '', stages }
  }
}

// ===== 提示词拼装 =====
function buildPrompt(sp, stage, idx) {
  const series = SERIES[sp.seriesId]
  const baseName = sp.name.split('→')[0]
  const lvName = sp.lv[stage.startLv] || sp.lv[stage.startLv + 1] || '未知形态'
  const prof = profiles[sp.id] || { form: '', symbol: '' }
  const init = initials[sp.id]
  const color = series.colors[idx]
  const life = lifeStories[sp.id]

  // —— 档案化角色（有六阶段人生档案）→ 人生叙事版提示词，与手绘 SVG / 档案完全对齐 ——
  if (life && life.stages[idx]) {
    const st = life.stages[idx]
    const sceneTone = (st.svg.match(/色调[:：]\s*([^。]+)/) || [])[1] || ''
    const zh =
      `${baseName}，${stage.zh}阶段·${st.name}（${st.age}，${st.keyword}）。` +
      `形象：${baseName}${prof.form ? '，' + prof.form + (prof.symbol ? ' 核心意象：' + prof.symbol + '。' : '。') : '。'}` +
      `品性：${st.character}。` +
      `姿态：${st.action}。` +
      `服饰：${st.attire}。` +
      (st.body ? `体型：${st.body}。` : '') +
      (st.attireDetail ? `衣物细节：${st.attireDetail}。` : '') +
      (st.hairStyle ? `发型妆造：${st.hairStyle}。` : '') +
      (st.face ? `脸型五官：${st.face}。` : '') +
      (st.weaponAction ? `武器招式：${st.weaponAction}。` : '') +
      `功法：${st.technique}。` +
      (st.powerEffect ? `功法表现：${st.powerEffect}。` : '') +
      `画面：${st.svg}。` +
      `台词：${st.line}。` +
      (st.actionSeq ? `动作帧（动图）：${st.actionSeq}。` : '') +
      `诗词：${st.poem}。` +
      (life.theme ? `主题句：${life.theme}。` : '') +
      `风格：${series.styleZh}。` +
      (sceneTone ? `色彩：${sceneTone}。` : `色彩：${color[0].name}（${color[0].hex}）主调 + ${color[1].name}（${color[1].hex}）点缀。`) +
      `画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。` +
      `画质：高细节，柔和渐变光，背景干净，无文字、无水印。`
    const en =
      `${series.styleEn}; ${stage.zh} ${st.name}, age ${st.age}, ${st.keyword}; scene: ${st.action}; ${st.poem};` +
      (sceneTone ? ` palette: ${sceneTone};` : ` palette ${color[0].hex} with ${color[1].hex} accents;`) +
      ` centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text`
    return { zh, en }
  }

  // egg 阶段（灵胎初醒）→ 用初始形态描述（符合物种分类 + 属性），不引用成年形态
  const isInit = idx === 0 && init
  const formPart = isInit
    ? `初始形态：${init.desc}。${init.elem}属性灵光微微环绕。`
    : `形象：${prof.form}` + (prof.symbol ? ` 核心意象：${prof.symbol}。` : '。')

  // —— 六阶段精修：神态 / 动作 / 衣着 / 梳造 ——
  const cat = refine.cats[sp.id] || 'animal'
  const b = (refine.base[cat] && refine.base[cat][idx]) || { spirit: '神态自若', action: '凝神而立', attire: '本相自然', hair: '发式如常' }
  const special = refine.special[sp.id]
  const attire = special ? special[idx].attire : b.attire
  const hair = special ? special[idx].hair : b.hair
  const movement = prof.movement ? prof.movement.replace(/[。！]$/, '') : ''
  let action = b.action
  // 标志动作只在「道法初成(三阶)首次展露」与「道果圆满(六阶)绝技大成」注入，避免 3-6 阶动作尾巴连重复
  if (movement && (idx === 2 || idx === 5)) action += '，' + movement

  const zh =
    `${baseName}，${stage.zh}阶段·${isInit ? init.name : lvName}。` +
    formPart +
    `神态：${b.spirit}。` +
    `动作：${action}。` +
    `衣着：${attire}。` +
    `梳造：${hair}。` +
    `意境：${series.scenes[idx]}。` +
    `风格：${series.styleZh}。` +
    `色彩：${color[0].name}（${color[0].hex}）主调 + ${color[1].name}（${color[1].hex}）点缀。` +
    `构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。` +
    `画质：高细节，柔和渐变光，背景干净，无文字、无水印。`
  const en =
    `${series.styleEn}; ${series.sceneEn[idx]}; ${b.spirit}; palette ${color[0].hex} with ${color[1].hex} accents; ` +
    `centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text`
  return { zh, en }
}

// ===== 生成 Markdown =====
const bySeries = {}
for (const s of unique) (bySeries[s.seriesId] = bySeries[s.seriesId] || []).push(s)

const seriesOrder = Object.keys(SERIES)
const lines = []
lines.push('# 宠物 AI 生图提示词 · 全 125 物种 × 6 阶段（v2）')
lines.push('')
lines.push(`> 生成：${new Date().toISOString().slice(0, 10)} · 由 \`scripts/generate-pet-prompts-v2.js\` 自动生成，**不要手改**（改数据源后重跑）`)
lines.push('> 取代 v1 `pet-prompts.csv` 通用模板。生成图片后按 `docs/pet-image-manifest.md` 流程登记（MANIFEST + 拷贝到 backend/public/pets/）。')
lines.push('')
lines.push('## 一、使用说明')
lines.push('')
lines.push('- **中文主句**（每段首行）：复制给即梦 / 通义万相 / Stable Diffusion（中文模型）。')
lines.push('- **EN 行**：复制给 Midjourney / DALL·E / Stable Diffusion（英文模型）。')
lines.push('- **六阶等级制**：灵胎初醒(Lv1) → 凡尘砺心(Lv3) → 道法初成(Lv5) → 大劫淬炼(Lv7) → 封神登天(Lv9) → 道果圆满(Lv11)。')
lines.push('- **精修四维度**：每阶段含 神态 / 动作 / 衣着 / 梳造 四维描述（重点角色手写，其余按大类基座），见 `frontend-vue/src/utils/petRefine.ts`。')
lines.push('- **人生档案版**：东方神话 姜子牙 / 杨戬 / 雷震子（有六阶段人生档案）的提示词为「人生叙事」版（品性 / 姿态 / 服饰 / 功法 / 画面 / 台词 / 诗词 / 主题句），源自 `frontend-vue/src/utils/petLifeStories.ts`，与手绘 SVG 完全对齐；其余角色走抽象道行版（petRefine）。')
lines.push('- **画布**：800×1000 竖版，角色主体居中偏下占约 60%，正面 3/4 视角，无文字无水印。')
lines.push('- **命名规范**：产物存 `frontend-vue/public/pets/{seriesId}/{speciesId}-{stage}.webp`。')
lines.push('- **版权提醒**：宝可梦 / 数码宝贝 / 星座圣斗士 / 虹猫蓝兔为受保护 IP，本提示词按项目「版权直名、非商用致敬」政策直接使用角色名，请勿用于商用。')
lines.push('')
lines.push('## 二、初始形态分类规范（灵胎初醒 / Lv1）')
lines.push('')
lines.push('- **设计原则**：初始形态符合物种特点 —— 植物类=一粒种子 / 胎生哺乳=一团胚胎 / 卵生=一枚蛋卵 / 数码兽=数码蛋 / 神兽瑞兽=灵种灵卵 / 人物=灵光化形 / 食物=原料 / 器物=物件 / 剑客=剑意灵光')
lines.push('- **属性区别**：每个角色初始形态带十二属性（木火土金水 + 雷冰风光暗钢秘），由 `frontend-vue/src/utils/petInitial.ts` 定义')
lines.push('- **分布**：seed 种子 4 / embryo 胚胎 13 / egg 蛋卵 37 / digi 数码蛋 6 / spirit 灵光 14 / person 人物 20 / food 原料 11 / artifact 器物 3 / sword 剑意 9 / mythic 神兽灵种 10')
lines.push('- **形态名**建议同步 petData 的 Lv1 等级名（保持一致后 生图 / UI / 诗文 三方对齐）')
lines.push('')
lines.push('## 三、系列视觉规范总表')
lines.push('')
lines.push('| 系列 | 物种数 | 风格（中文） | 风格（英文） |')
lines.push('|------|--------|--------------|--------------|')
for (const sid of seriesOrder) {
  const s = SERIES[sid]
  lines.push(`| ${s.name} | ${bySeries[sid].length} | ${s.styleZh} | ${s.styleEn} |`)
}
lines.push('')
lines.push('## 四、分系列完整提示词')
lines.push('')

let seriesNo = 0
for (const sid of seriesOrder) {
  seriesNo++
  const s = SERIES[sid]
  const list = bySeries[sid]
  lines.push(`### ${seriesNo}. ${s.name}（${list.length} 物种）`)
  lines.push('')
  lines.push(`> **风格**：${s.styleZh}。**阶段演绎**：`)
  for (let i = 0; i < 6; i++) {
    lines.push(`> - ${STAGES[i].zh}：${s.scenes[i]}（${s.colors[i][0].name}/${s.colors[i][1].name}）`)
  }
  lines.push('')
  for (const sp of list) {
    const baseName = sp.name.split('→')[0]
    const isLife = !!lifeStories[sp.id]
    lines.push(`#### ${baseName}（\`${sp.id}\`）${isLife ? ' · 人生档案版' : ''}`)
    lines.push('')
    for (let i = 0; i < 6; i++) {
      const p = buildPrompt(sp, STAGES[i], i)
      const lvName = (i === 0 && initials[sp.id])
        ? initials[sp.id].name
        : (sp.lv[STAGES[i].startLv] || sp.lv[STAGES[i].startLv + 1] || '未知形态')
      lines.push(`**${STAGES[i].zh} · ${lvName}**`)
      lines.push(`- ${p.zh}`)
      lines.push(`- EN：${p.en}`)
      lines.push('')
    }
  }
}

const md = lines.join('\n')
const outPath = process.argv[2] || path.join(__dirname, '../../docs/宠物AI生图提示词.md')
fs.writeFileSync(outPath, md, 'utf8')
console.log(`已生成 ${unique.length} 物种 × 6 阶段 = ${unique.length * 6} 条提示词 → ${outPath}`)
console.log('示例：')
const demo = unique.find(s => s.id === 'pikachu')
console.log(buildPrompt(demo, STAGES[0], 0).zh)
