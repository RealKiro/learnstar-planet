// ===== 阶段 AI 生图提示词组装器 =====
// 数据源：petLifeStories 六阶段档案（叙事/视觉规格）+ seriesStages 阶段称谓 + petArtStyles 画风 + 统一约束。
// 优先级：档案中人工精写的 aiPrompt（作为主体描述）> 自动组装；两种来源都会追加统一约束后缀。
// 互动剧情（interaction 字段）不进提示词（提示词只管画面），属剧情展示层。

import { PET_LIFE_STORIES } from './petLifeStories'
import { getSpeciesById, getSeriesBySpeciesId, getLevelStage } from './petData'
import { getArtStyleById } from './petArtStyles'
import { getStageLabel, getSeriesArtStyle, type PetStageKey } from './seriesStages'

/** 统一约束：所有提示词共用，保证产出图可直接替换使用 */
export const AI_IMAGE_CONSTRAINTS =
  '固定画布 1024×1024，PNG 格式，透明背景（或纯白背景便于抠图），' +
  '角色单主体居中、正面 3/4 视角、占画面 60-70%，' +
  '全身完整不出画框，无文字、无水印、无边框，' +
  '边缘干净利于抠图，同系列多图风格严格一致'

/** 组装某角色某阶段的 AI 生图提示词（画风由系列决定，同系列全阶段严格一致） */
export function composeStageAiPrompt(speciesId: string, level: number): string {
  const story = PET_LIFE_STORIES[speciesId]
  const stageIdx = stageIndexOf(level)
  const st = story?.stages?.[stageIdx]
  const species = getSpeciesById(speciesId)
  const series = getSeriesBySpeciesId(speciesId)
  const style = getArtStyleById(getSeriesArtStyle(series?.id))
  const stageKey = getLevelStage(level) as PetStageKey
  const stageName = getStageLabel(stageKey, series?.id)

  // 人工精写的提示词优先（作为主体描述）
  if (st?.aiPrompt && st.aiPrompt.trim()) {
    return `${st.aiPrompt.trim()}。${AI_IMAGE_CONSTRAINTS}，画风：${style.cn}`
  }

  const parts: string[] = []
  parts.push(`${style.cn}风格插画`)
  if (species) parts.push(`角色「${species.name}」的「${stageName}」阶段形态`)
  if (st?.keyword) parts.push(st.keyword)
  if (st?.character) parts.push(`神态：${st.character}`)
  if (st?.action) parts.push(`动作：${st.action}`)
  if (st?.attire || st?.attireDetail) parts.push(`服饰：${[st.attire, st.attireDetail].filter(Boolean).join('，')}`)
  if (st?.technique) parts.push(`能力显化：${st.technique}`)
  if (st?.svg) parts.push(`画面：${st.svg}`)
  if (st?.mood) parts.push(`氛围：${st.mood}`)
  parts.push(`系列：${series?.name ?? '原创'}`)

  return parts.filter(Boolean).join('；') + `。${AI_IMAGE_CONSTRAINTS}，画风：${style.cn}`
}

function stageIndexOf(level: number): number {
  if (level <= 2) return 0
  if (level <= 4) return 1
  if (level <= 6) return 2
  if (level <= 8) return 3
  if (level <= 10) return 4
  return 5
}

