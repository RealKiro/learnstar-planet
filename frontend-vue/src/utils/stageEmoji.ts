// ===== 学趣星球 · 成长阶段 Emoji（体现进化特点） =====
// 同一角色在不同阶段显示不同 emoji：蛋 → 幼年 → 成长期 → 成熟期 → 传说 → 道果
// 未单独配置的阶段回退：蛋=🥚，幼年/成长期/成熟期=基础 emoji，传说=✨，道果=👑

import { getSpeciesEmoji } from './petData'

export type PetStage = 'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent'

/** 等级 → 阶段（与 petArtParts 的 6 阶段阈值一致） */
export function stageKeyOf(level: number): PetStage {
  if (level >= 11) return 'transcendent'
  if (level >= 9) return 'legendary'
  if (level >= 7) return 'mature'
  if (level >= 5) return 'growing'
  if (level >= 3) return 'baby'
  return 'egg'
}

/**
 * 重点角色各阶段专属 emoji（未配置的用通用回退）。
 * legendary = 传说级威能，transcendent = 道果终极形态。
 */
const STAGE_EXTRA: Record<string, { legendary?: string; transcendent?: string }> = {
  // ===== 山海经·神话 =====
  zhulong:       { legendary: '🐲', transcendent: '☀️' },
  yinglong:      { legendary: '🐲', transcendent: '🌧️' },
  nine_tail_fox: { legendary: '🦊', transcendent: '👑' },
  fenghuang:     { legendary: '🔥', transcendent: '🌟' },
  qilin:         { legendary: '🦄', transcendent: '✨' },
  qinglong:      { legendary: '🐉', transcendent: '☁️' },
  baihu:         { legendary: '🐆', transcendent: '🐯' },
  zhuque:        { legendary: '🦩', transcendent: '🔥' },
  xuanwu:        { legendary: '🛡️', transcendent: '⛰️' },
  // ===== 宝可梦 =====
  charmander:    { legendary: '🔥', transcendent: '💥' },
  bulbasaur:     { legendary: '🌿', transcendent: '🌸' },
  squirtle:      { legendary: '💧', transcendent: '🌊' },
  pikachu:       { legendary: '⚡', transcendent: '☄️' },
  // ===== 东方神话 =====
  sun_wukong:    { legendary: '🐵', transcendent: '👑' },
  nezha:         { legendary: '🔥', transcendent: '💥' },
  lei_zhenzi:    { legendary: '⚡', transcendent: '🌩️' },
  yang_jian:     { legendary: '🗡️', transcendent: '🌟' },
  taishang_laojun: { legendary: '☯️', transcendent: '🌌' },
  zhong_kui:     { legendary: '👺', transcendent: '🔥' },
  // ===== 国宝 =====
  panda:         { legendary: '🐼', transcendent: '🌸' },
  south_china_tiger: { legendary: '🐆', transcendent: '🐅' },
  // ===== 星座（终极 = 圣衣全开） =====
  aries: { legendary: '♈', transcendent: '🌟' },
  leo: { legendary: '♌', transcendent: '🌟' },
  virgo: { legendary: '♍', transcendent: '🌟' },
  sagittarius: { legendary: '♐', transcendent: '🌟' },
  // ===== 御三家进化链 =====
  riolu:         { legendary: '🐺', transcendent: '🐲' },
  eevee:         { legendary: '🦝', transcendent: '🦊' },
}

/** 成长阶段 Emoji：蛋 → 幼年/成长/成熟(基础) → 传说(威能) → 道果(终极) */
export function getStageEmoji(speciesId: string, level: number): string {
  const stage = stageKeyOf(level)
  if (stage === 'egg') return '🥚'
  const base = getSpeciesEmoji(speciesId)
  if (stage === 'legendary') return STAGE_EXTRA[speciesId]?.legendary || '✨'
  if (stage === 'transcendent') return STAGE_EXTRA[speciesId]?.transcendent || '👑'
  return base
}
