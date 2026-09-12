// ===== 分系列阶段称谓（替代全局统一的"新生/幼年/成长期/成熟期/传说级/道果"） =====
// 六阶段 key 与游戏逻辑（getLevelStage / 特质索引 / 阶段 emoji）保持不变，
// 仅展示称谓按系列 flavour 化：数码兽说"究极体"而非"道果"，圣斗士穿圣衣，食物过节日。

export type PetStageKey = 'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent'

export const GENERIC_STAGE_LABELS: Record<PetStageKey, string> = {
  egg: '新生',
  baby: '幼年',
  growing: '成长期',
  mature: '成熟期',
  legendary: '传说级',
  transcendent: '道果',
}

export const SERIES_STAGE_LABELS: Record<string, Record<PetStageKey, string>> = {
  // 山海经：洪荒异兽
  myth: { egg: '灵胎初孕', baby: '洪荒初行', growing: '灵窍初开', mature: '山海称雄', legendary: '异兽扬名', transcendent: '山海归一' },
  // 东方神话：封神修行路
  dongfang: { egg: '灵光化形', baby: '初入尘世', growing: '拜师学艺', mature: '劫难淬炼', legendary: '肉身成圣', transcendent: '封神归位' },
  // 宝可梦：训练家羁绊
  pokemon: { egg: '神秘蛋', baby: '初遇伙伴', growing: '特训磨炼', mature: '进化前夜', legendary: '冠军之路', transcendent: '传奇宝可梦' },
  // 数码宝贝：原著进化术语
  digimon: { egg: '数码蛋', baby: '幼年期', growing: '成长期', mature: '成熟期', legendary: '完全体', transcendent: '究极体' },
  // 国宝：自然保护之路
  national: { egg: '珍稀幼崽', baby: '茁壮成长', growing: '野性磨炼', mature: '家园守护', legendary: '旗舰物种', transcendent: '生态传奇' },
  // 魔法奇幻：魔法学院等级
  magic: { egg: '魔法学徒', baby: '见习法师', growing: '正式法师', mature: '高阶法师', legendary: '大魔导师', transcendent: '传奇法神' },
  // 史前生物：地质纪元
  prehistoric: { egg: '化石苏醒', baby: '纪元初生', growing: '蛮荒成长', mature: '纪元争霸', legendary: '霸主称雄', transcendent: '永恒化石' },
  // 星座守护：圣衣递进（圣斗士原著等级）
  constellation: { egg: '初入圣域', baby: '青铜圣衣', growing: '白银圣衣', mature: '黄金圣衣', legendary: '第七感觉醒', transcendent: '神圣衣全开' },
  // 传统节日：民俗节庆
  festival: { egg: '节日将至', baby: '张灯结彩', growing: '阖家团圆', mature: '民俗登场', legendary: '万家灯火', transcendent: '世代传承' },
  // 虹猫蓝兔七侠传：七剑武侠
  qixia: { egg: '初入江湖', baby: '剑意萌芽', growing: '剑法小成', mature: '剑胆琴心', legendary: '七剑合璧', transcendent: '侠义传世' },
}

export function getStageLabelMap(seriesId?: string): Record<string, string> {
  return (seriesId && SERIES_STAGE_LABELS[seriesId]) || GENERIC_STAGE_LABELS
}

export function getStageLabel(stageKey: PetStageKey | string, seriesId?: string): string {
  return getStageLabelMap(seriesId)[stageKey as PetStageKey] ?? GENERIC_STAGE_LABELS[stageKey as PetStageKey] ?? stageKey
}

// ===== 系列 → 画风绑定（AI 生图用；同系列所有阶段严格同一画风，保证角色一致性） =====
export const SERIES_ART_STYLE: Record<string, string> = {
  myth: 'inkcn',         // 山海经 → 中国风水墨（上古神话）
  dongfang: 'inkcn',     // 东方神话 → 中国风水墨（封神仙侠）
  pokemon: 'anime',      // 宝可梦 → 动漫（日系赛璐璐）
  digimon: 'anime',      // 数码宝贝 → 动漫（数码进化）
  qixia: 'anime',        // 虹猫蓝兔 → 动漫（2D 国漫武侠）
  national: 'cartoon',   // 国宝 → 卡通（皮克斯萌化）
  magic: 'cg3d',         // 魔法奇幻 → 3D 渲染 CG（史诗奇幻）
  prehistoric: 'cg3d',   // 史前生物 → 3D 渲染 CG（纪录片复原）
  constellation: 'figure', // 星座守护 → 手办（圣衣金属质感）
  festival: 'watercolor',  // 传统节日 → 水彩（民俗喜庆柔和）
}

export function getSeriesArtStyle(seriesId?: string): string {
  return (seriesId && SERIES_ART_STYLE[seriesId]) || 'anime'
}
