// ===== 学趣星球 · AI 宠物图片加载 =====
// 数据契约见 docs/宠物系统全案.md 6.4 节：
//   目录 public/pets/{seriesId}/{speciesId}-{stage}.webp（egg/baby/growing/mature/legendary/transcendent）
//   六阶视觉规范 + 提示词模板已就绪；生成图片后在该文件的 MANIFEST 登记 speciesId，PetSprite 自动优先加载图片、回退 SVG。

import { getSeriesBySpeciesId } from './petData'

/**
 * 已生成 AI 宠物图的物种 id 清单（AI 生图产物为 webp，生成后在此登记）。
 * 例：['sun_wukong', 'zhong_kui', ...]
 */
const MANIFEST: string[] = []

/**
 * 已手绘 SVG 插画的物种 id 清单（样板阶段按提示词手绘的独立 SVG，每物种 6 张 {speciesId}-{stage}.svg）。
 * 例：['zhulong', 'charmander', ...]
 */
const SVG_MANIFEST: string[] = ['zhulong', 'charmander', 'sun_wukong', 'nezha', 'zhong_kui', 'jiang_ziya', 'yang_jian', 'lei_zhenzi', 'huang_tianhua', 'tu_xingsun', 'yang_ren', 'wei_hu', 'daji', 'shen_gongbao', 'lv_dongbin', 'he_xiangu', 'zhang_guolao']

/**
 * 根据物种 + 阶段返回图片 URL（webp 优先，其次手绘 SVG）；均未生成时返回 null（PetSprite 回退程序化 SVG）。
 */
export function getPetImageUrl(speciesId: string, stage: string): string | null {
  const seriesId = getSeriesBySpeciesId(speciesId)?.id || 'unknown'
  if (MANIFEST.includes(speciesId)) {
    return `${import.meta.env.BASE_URL}pets/${seriesId}/${speciesId}-${stage}.webp`
  }
  if (SVG_MANIFEST.includes(speciesId)) {
    return `${import.meta.env.BASE_URL}pets/${seriesId}/${speciesId}-${stage}.svg`
  }
  return null
}
