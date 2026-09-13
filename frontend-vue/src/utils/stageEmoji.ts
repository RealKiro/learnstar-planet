// ===== 学宠星球 · 成长阶段 Emoji（体现进化特点） =====
// 同一角色在不同阶段显示不同 emoji：蛋 → 幼年 → 成长期 → 成熟期 → 传说 → 道果
// 未单独配置的阶段回退：蛋=🥚，幼年/成长期/成熟期=基础 emoji，传说=✨，道果=👑

import { getSpeciesEmoji, getLevelStage } from './petData'

export type PetStage = 'egg' | 'baby' | 'growing' | 'mature' | 'legendary' | 'transcendent'

/** 等级 → 阶段（单一真源在 petData.getLevelStage，此处仅做类型收窄） */
export function stageKeyOf(level: number): PetStage {
  return getLevelStage(level) as PetStage
}

/**
 * 全量角色各阶段专属 emoji（126 物种全覆盖，同系列内不重复）。
 * legendary = 传说级威能，transcendent = 道果终极形态。
 */
const STAGE_EXTRA: Record<string, { legendary?: string; transcendent?: string }> = {
  // ===== 山海经（18） =====
  zhulong:       { legendary: '🐲', transcendent: '☀️' },
  yinglong:      { legendary: '🕊️', transcendent: '🌧️' },
  nine_tail_fox: { legendary: '🦊', transcendent: '👑' },
  kunpeng:       { legendary: '🐋', transcendent: '🕊️' },
  fenghuang:     { legendary: '🔥', transcendent: '🦚' },
  qilin:         { legendary: '🦄', transcendent: '✨' },
  qinglong:      { legendary: '🐉', transcendent: '☁️' },
  baihu:         { legendary: '🐆', transcendent: '🐯' },
  zhuque:        { legendary: '🦩', transcendent: '🔥' },
  xuanwu:        { legendary: '🛡️', transcendent: '⛰️' },
  taotie:        { legendary: '🌀', transcendent: '♾️' },
  baize:         { legendary: '📖', transcendent: '🔮' },
  qiongqi:       { legendary: '🦅', transcendent: '🏞️' },
  bifang:        { legendary: '🪔', transcendent: '🪶' },
  pixiu:         { legendary: '🧧', transcendent: '💰' },
  jingwei:       { legendary: '🪨', transcendent: '🏝️' },
  xiangliu:      { legendary: '🐍', transcendent: '🌫️' },
  xiezhi:        { legendary: '⚖️', transcendent: '🐏' },
  // ===== 宝可梦（12） =====
  charmander:    { legendary: '🔥', transcendent: '💥' },
  bulbasaur:     { legendary: '🌿', transcendent: '🌸' },
  squirtle:      { legendary: '💧', transcendent: '🌊' },
  pikachu:       { legendary: '⚡', transcendent: '☄️' },
  riolu:         { legendary: '🐺', transcendent: '🐲' },
  eevee:         { legendary: '🦝', transcendent: '🦊' },
  ice_fox:       { legendary: '❄️', transcendent: '🧊' },
  rock_rhino:    { legendary: '🦏', transcendent: '🗿' },
  wind_falcon:   { legendary: '🦅', transcendent: '🌪️' },
  light_deer:    { legendary: '🌅', transcendent: '🦌' },
  dark_panther:  { legendary: '🌑', transcendent: '🖤' },
  steel_armadillo: { legendary: '🧲', transcendent: '⚙️' },
  // ===== 国宝（12） =====
  panda:         { legendary: '🐼', transcendent: '🌸' },
  south_china_tiger: { legendary: '🐯', transcendent: '🐅' },
  golden_monkey: { legendary: '🐒', transcendent: '🥇' },
  red_crowned_crane: { legendary: '🪽', transcendent: '🕊️' },
  chinese_alligator: { legendary: '🐊', transcendent: '⏳' },
  crested_ibis:  { legendary: '🩷', transcendent: '🌈' },
  tibetan_antelope: { legendary: '🏔️', transcendent: '🏃' },
  snow_leopard:  { legendary: '🌨️', transcendent: '🐆' },
  milu_deer:     { legendary: '🫎', transcendent: '🎋' },
  siberian_tiger: { legendary: '🌬️', transcendent: '🌙' },
  red_panda:     { legendary: '🍁', transcendent: '🧣' },
  finless_porpoise: { legendary: '🐬', transcendent: '🫧' },
  // ===== 数码宝贝（6） =====
  mecha_dragon:  { legendary: '🦖', transcendent: '🤖' },
  cyber_cat:     { legendary: '🐱', transcendent: '💠' },
  space_mecha:   { legendary: '🚀', transcendent: '🛰️' },
  quantum_beast: { legendary: '🔷', transcendent: '🧬' },
  digital_phoenix: { legendary: '📡', transcendent: '💫' },
  mecha_shark:   { legendary: '🦈', transcendent: '⚓' },
  // ===== 魔法奇幻（12） =====
  unicorn:       { legendary: '🌈', transcendent: '💎' },
  wyvern:        { legendary: '🐉', transcendent: '🌋' },
  fairy:         { legendary: '🧚', transcendent: '🍃' },
  treant:        { legendary: '🌳', transcendent: '🍄' },
  griffin:       { legendary: '🦁', transcendent: '🏛️' },
  mermaid:       { legendary: '🧜', transcendent: '🐚' },
  grey_wizard:   { legendary: '🧙', transcendent: '📜' },
  wand_cat:      { legendary: '🪄', transcendent: '🎩' },
  dragon_knight: { legendary: '⚔️', transcendent: '🏰' },
  alchemy_golem: { legendary: '🧱', transcendent: '⚗️' },
  nightmare_horse: { legendary: '🐴', transcendent: '🌘' },
  lamp_spirit:   { legendary: '💡', transcendent: '🏮' },
  // ===== 史前生物（12） =====
  t_rex:         { legendary: '🦖', transcendent: '🦷' },
  triceratops:   { legendary: '🔺', transcendent: '🛡️' },
  pterosaur:     { legendary: '🪁', transcendent: '🦅' },
  mammoth:       { legendary: '🦣', transcendent: '❄️' },
  sabertooth:    { legendary: '🗡️', transcendent: '⚔️' },
  mosasaur:      { legendary: '🌊', transcendent: '🐙' },
  spinosaurus:   { legendary: '🦴', transcendent: '☀️' },
  ankylosaurus:  { legendary: '🔨', transcendent: '🗿' },
  diplodocus:    { legendary: '🦕', transcendent: '🛤️' },
  megalodon:     { legendary: '🦈', transcendent: '🌀' },
  ground_sloth:  { legendary: '🦥', transcendent: '🌍' },
  woolly_rhino:  { legendary: '❄️', transcendent: '🦏' },
  // ===== 星座守护（12，道果统一神圣衣全开） =====
  aries: { legendary: '♈', transcendent: '🌟' },
  taurus: { legendary: '♉', transcendent: '🌟' },
  gemini: { legendary: '♊', transcendent: '🌟' },
  cancer: { legendary: '♋', transcendent: '🌟' },
  leo: { legendary: '♌', transcendent: '🌟' },
  virgo: { legendary: '♍', transcendent: '🌟' },
  libra: { legendary: '♎', transcendent: '🌟' },
  scorpio: { legendary: '♏', transcendent: '🌟' },
  sagittarius: { legendary: '♐', transcendent: '🌟' },
  capricorn: { legendary: '♑', transcendent: '🌟' },
  aquarius: { legendary: '♒', transcendent: '🌟' },
  pisces: { legendary: '♓', transcendent: '🌟' },
  // ===== 传统节日（12） =====
  zongzi:        { legendary: '🫔', transcendent: '🛶' },
  tangyuan:      { legendary: '🍡', transcendent: '🥣' },
  mooncake:      { legendary: '🥮', transcendent: '🌕' },
  qingtuan:      { legendary: '🟩', transcendent: '🕯️' },
  chongyang_cake: { legendary: '🧗', transcendent: '🍂' },
  niangao:       { legendary: '📈', transcendent: '🎊' },
  laba_porridge: { legendary: '🍲', transcendent: '🫘' },
  spring_pancake: { legendary: '🥞', transcendent: '🌱' },
  tanghulu:      { legendary: '🍒', transcendent: '❤️' },
  osmanthus_cake: { legendary: '🌼', transcendent: '🫖' },
  wonton:        { legendary: '🥟', transcendent: '🪙' },
  festival_lantern: { legendary: '🏮', transcendent: '🎆' },
  // ===== 虹猫蓝兔七侠传（10，七剑合璧意象） =====
  hongmao:       { legendary: '🌈', transcendent: '🗡️' },
  lantu:         { legendary: '🫧', transcendent: '❄️' },
  doudou:        { legendary: '💊', transcendent: '🩺' },
  dabeng:        { legendary: '🪓', transcendent: '💪' },
  tiaotiao:      { legendary: '🐇', transcendent: '🤸' },
  shali:         { legendary: '💜', transcendent: '🌸' },
  dada:          { legendary: '🪃', transcendent: '💨' },
  heixinhu:      { legendary: '😈', transcendent: '🕳️' },
  heixiaohu:     { legendary: '🐅', transcendent: '🥀' },
  ma_sanniang:   { legendary: '🎯', transcendent: '🌙' },
  // ===== 东方神话（20） =====
  sun_wukong:    { legendary: '🐵', transcendent: '👑' },
  nezha:         { legendary: '🌀', transcendent: '💥' },
  lei_zhenzi:    { legendary: '⚡', transcendent: '🌩️' },
  yang_jian:     { legendary: '🔱', transcendent: '⚡' },
  taishang_laojun: { legendary: '☯️', transcendent: '🌌' },
  zhong_kui:     { legendary: '👺', transcendent: '🔥' },
  jiang_ziya:    { legendary: '🎣', transcendent: '🏳️' },
  huang_tianhua: { legendary: '🔨', transcendent: '🌻' },
  tu_xingsun:    { legendary: '🕳️', transcendent: '⛏️' },
  yang_ren:      { legendary: '👁️', transcendent: '🪭' },
  wei_hu:        { legendary: '⚔️', transcendent: '☸️' },
  daji:          { legendary: '🌺', transcendent: '🦊' },
  shen_gongbao:  { legendary: '🐆', transcendent: '🪷' },
  lv_dongbin:    { legendary: '🗡️', transcendent: '🧘' },
  he_xiangu:     { legendary: '🪷', transcendent: '🍀' },
  zhang_guolao:  { legendary: '🫏', transcendent: '🥁' },
  tie_guaili:    { legendary: '🦯', transcendent: '👻' },
  han_zhongli:   { legendary: '🪭', transcendent: '🪙' },
  lan_caihe:     { legendary: '💐', transcendent: '🎶' },
  cao_guojiu:    { legendary: '🎴', transcendent: '🤍' },
}

/** 系列统一道果 emoji（显式声明的系列设计）：
 *  这些系列的道果阶段共用同一 emoji 属设计意图而非撞车，
 *  analyze-pets.mjs 的「同阶段雷同 emoji」检测据此豁免。
 *  key 为系列 id，value 为该系列道果统一 emoji。 */
export const SERIES_UNIFIED_EMOJI: Record<string, string> = {
  constellation: '🌟', // 星座守护：12 星座道果皆化为本命星光
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
