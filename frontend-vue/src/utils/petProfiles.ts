// ===== 学宠星球 · 宠物角色档案（形态 / 习性） =====
// 按物种 id 索引；缺失时返回兜底，保证所有宠物图鉴可展示。
// 专属诗文见 petHandbookData.getPoems，进化台词见 getEvoLines。
// 数据层：本文件只留类型与查询函数，数据按系列存于 src/data/pets/<seriesId>.json（每系列一个 JSON）。

/** 六阶段 key：新生 / 幼年 / 成长期 / 成熟期 / 传说级 / 道果 */
export type PetStage = 'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent'

/** 单阶段九维档案（各阶段不同又一脉相承） */
export interface StageProfile {
  /** 形态 */
  form: string
  /** 习性 */
  habit: string
  /** 本源（出身 / 血统 / 师承） */
  origin?: string
  /** 别称（民间称谓 / 封号） */
  epithet?: string
  /** 标志性动作 */
  movement?: string
  /** 核心意象（提到他想到什么） */
  symbol?: string
  /** 主题句（一句概括其一生） */
  theme?: string
}

export interface PetProfile extends StageProfile {
  /** 按阶段九维档案：未写某阶段时回退到全局九维（同物种各阶段未覆盖字段继承全局） */
  stages?: Partial<Record<PetStage, StageProfile>>
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

export const PET_PROFILES: Record<string, PetProfile> = Object.assign(
  {},
  (myth as SeriesPetData).profiles,
  (pokemon as SeriesPetData).profiles,
  (national as SeriesPetData).profiles,
  (digimon as SeriesPetData).profiles,
  (magic as SeriesPetData).profiles,
  (prehistoric as SeriesPetData).profiles,
  (constellation as SeriesPetData).profiles,
  (festival as SeriesPetData).profiles,
  (qixia as SeriesPetData).profiles,
  (dongfang as SeriesPetData).profiles,
) as Record<string, PetProfile>

/** 获取角色档案：传入 stage 时优先返回该阶段九维（未写阶段或字段缺失时回退全局九维） */
export function getPetProfile(speciesId: string, stage?: PetStage): PetProfile {
  const p = PET_PROFILES[speciesId]
  if (!p) return { form: '形态尚待探明的神秘宠物。', habit: '习性未知，正等待与小主人一起探索。' }
  if (stage && p.stages?.[stage]) {
    return { ...p, ...p.stages[stage] }
  }
  return p
}
