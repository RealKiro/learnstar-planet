// ===== 学宠星球 · 宠物特性/技能 =====
// 物种 × 6 阶段(卵生→幼年→成长→成熟→传说→道果)的独立「角色特点」与「技能」，
// 保证各类别各宠物各阶段特点互不相同。
// 由 getStagePersonality / getStageAbility 按等级映射取用。
// 数据层：本文件只留类型与查询函数，数据按系列存于 src/data/pets/<seriesId>.json（每系列一个 JSON）。

import { getLevelStage } from './petData'

/** 每个物种 6 阶段(egg/baby/growing/mature/legendary/transcendent)的特性与技能 */
export interface PetTraits {
  /** 6 阶段角色特点 */
  personality: string[]
  /** 6 阶段技能 */
  ability: string[]
}

const STAGE_INDEX: Record<string, number> = { egg: 0, baby: 1, growing: 2, mature: 3, legendary: 4, transcendent: 5 }

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

export const PET_TRAITS: Record<string, PetTraits> = Object.assign(
  {},
  (myth as SeriesPetData).traits,
  (pokemon as SeriesPetData).traits,
  (national as SeriesPetData).traits,
  (digimon as SeriesPetData).traits,
  (magic as SeriesPetData).traits,
  (prehistoric as SeriesPetData).traits,
  (constellation as SeriesPetData).traits,
  (festival as SeriesPetData).traits,
  (qixia as SeriesPetData).traits,
  (dongfang as SeriesPetData).traits,
) as Record<string, PetTraits>

/** 获取某等级的角色特点 */
export function getStagePersonality(speciesId: string, level: number): string {
  const t = PET_TRAITS[speciesId]
  if (!t) return ''
  const idx = STAGE_INDEX[getLevelStage(level)] ?? 0
  // 部分物种特质数组仍为 5 项，第 6 阶段(道果)回退到末项(终极)
  return t.personality[idx] || t.personality[t.personality.length - 1] || ''
}

/** 获取某等级的技能 */
export function getStageAbility(speciesId: string, level: number): string {
  const t = PET_TRAITS[speciesId]
  if (!t) return ''
  const idx = STAGE_INDEX[getLevelStage(level)] ?? 0
  // 部分物种特质数组仍为 5 项，第 6 阶段(道果)回退到末项(终极)
  return t.ability[idx] || t.ability[t.ability.length - 1] || ''
}
