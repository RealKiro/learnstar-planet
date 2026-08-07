// ===== 学趣星球 · AI 宠物图片加载 =====
// 数据契约见 docs/宠物系统全案.md 6.4 节：
//   目录 public/pets/{seriesId}/{speciesId}-{stage}.webp（egg/baby/growing/mature/legendary/transcendent）
//   六阶视觉规范 + 提示词模板已就绪；生成图片后在该文件的 MANIFEST 登记 speciesId，PetSprite 自动优先加载图片、回退 SVG。

import { getSeriesBySpeciesId } from './petData'

/**
 * 已生成 AI 宠物图的物种 id 清单（生成后在此登记）。
 * 例：['sun_wukong', 'zhong_kui', ...]
 */
const MANIFEST: string[] = []

/** 根据物种 + 阶段返回 AI 图片 URL；未生成时返回 null（PetSprite 回退程序化 SVG） */
export function getPetImageUrl(speciesId: string, stage: string): string | null {
  if (!MANIFEST.includes(speciesId)) return null
  const seriesId = getSeriesBySpeciesId(speciesId)?.id || 'unknown'
  return `${import.meta.env.BASE_URL}pets/${seriesId}/${speciesId}-${stage}.webp`
}
