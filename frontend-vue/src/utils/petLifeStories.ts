// ===== 学宠星球 · 六阶段人生成长档案 =====
// 每个角色的六阶段不是抽象的"道行阶段"，而是从「出生 → 成长/变故 → 困惑/顿悟拜师 → 苦修 → 大成 → 人生归宿」的具体人生叙事。
// 参考规范：`姜子牙 · 六阶段成长档案.txt`（项目根目录）。
// 每阶含 7 维度 + 年龄线 + 心境，是 SVG 绘制 / 提示词 / 图鉴 的共同内容来源。
// 数据层：本文件只留类型与查询函数，数据按系列存于 src/data/pets/<seriesId>.json（每系列一个 JSON）。

/** 单阶生灵档案 */
export interface LifeStageStory {
  /** 本阶人生阶段名（如"朝歌市井"） */
  name: string
  /** 年龄线（如"32岁"） */
  age: string
  /** 关键词（如"朝歌卖面·市井隐志"） */
  keyword: string
  /** 品性性格 */
  character: string
  /** 行为动作 */
  action: string
  /** 服饰梳造 */
  attire: string
  /** 武器功法 */
  technique: string
  /** 进化台词 */
  line: string
  /** 专属诗词 */
  poem: string
  /** SVG 视觉描述（构图/色调/氛围，供绘制） */
  svg: string
  /** 心境 */
  mood: string
  /** ===== 视觉规格（为 SVG 绘制 / 动图动画 / AI 生图服务）===== */
  // 人物/半人形角色按字段本义填写；非人物角色（兽/神兽/精灵/食物/器物/剑/数码）做语义扩展：
  //   body=体量比例 / face=神情五官(拟人化) / attireDetail=皮表纹样 / hairStyle=形态部件(角尾翼鳍鬃甲)
  //   weaponAction=招牌攻击(可留空) / powerEffect=能力显化。详见 docs/角色档案·非人物视觉规格速查.md
  /** 身高胖瘦体型：头身比/轮廓/质感（如"身高约6头身，少年精瘦，四肢修长"）；非人物=体量/比例/姿态轮廓 */
  body?: string
  /** 行为动作序列：为动图准备的关键帧（①②③…，每帧一个连贯动作） */
  actionSeq?: string
  /** 衣物挂饰装饰：颜色/纹理/挂饰（如"锁子黄金甲金鳞细密，腰束红绦垂金穗"） */
  attireDetail?: string
  /** 发型梳造妆造：发型/束发/面饰 */
  hairStyle?: string
  /** 脸型五官：眉/眼/鼻/下巴等面部特征（人物/半人形角色必填） */
  face?: string
  /** 武器形态与招式动作：武器形制+招式分解 */
  weaponAction?: string
  /** 功法手段/法力表现：能力显化的视觉效果 */
  powerEffect?: string
  /** ===== 剧情推演与 AI 生图（2026-09 新增）===== */
  /** 本阶段与其他角色的互动/冲突（可以是本项目角色，也可以是原著/IP 中的角色）——用于推动剧情推演，如"与哪吒斗法于陈塘关外" */
  interaction?: string
  /** 本阶段形态的 AI 生图提示词（人工精写时填写；未填写由 utils/stageAiPrompt.ts 的 composeStageAiPrompt 按统一约束自动组装） */
  aiPrompt?: string
}

/** 角色生灵档案（六阶） */
export interface PetLifeStory {
  id: string
  name: string
  /** 主题句 */
  theme: string
  stages: LifeStageStory[]
}

import myth from '../data/pets/myth.json'
import pokemon from '../data/pets/pokemon.json'
import national from '../data/pets/national.json'
import digimon from '../data/pets/digimon.json'
import magic from '../data/pets/magic.json'
import prehistoric from '../data/pets/prehistoric.json'
import constellation from '../data/pets/constellation.json'
import festival from '../data/pets/festival.json'
import qixia from '../data/pets/qixia.json'
import dongfang from '../data/pets/dongfang.json'

type SeriesPetData = { stories?: Record<string, unknown>; profiles?: Record<string, unknown>; traits?: Record<string, unknown> }

export const PET_LIFE_STORIES: Record<string, PetLifeStory> = Object.assign(
  {},
  (myth as SeriesPetData).stories,
  (pokemon as SeriesPetData).stories,
  (national as SeriesPetData).stories,
  (digimon as SeriesPetData).stories,
  (magic as SeriesPetData).stories,
  (prehistoric as SeriesPetData).stories,
  (constellation as SeriesPetData).stories,
  (festival as SeriesPetData).stories,
  (qixia as SeriesPetData).stories,
  (dongfang as SeriesPetData).stories,
) as Record<string, PetLifeStory>

export function getPetLifeStory(speciesId: string): PetLifeStory | null {
  return PET_LIFE_STORIES[speciesId] || null
}
