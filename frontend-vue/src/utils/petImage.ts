// ===== 学趣星球 · 宠物图加载（已停用手绘/AI 图） =====
//
// ⚠️ 2026-08-08 起停用手绘 SVG / AI 生图（原 SVG_MANIFEST / MANIFEST）：
// 手绘 SVG 自带方形场景背景，在圆形裁剪（课堂评价卡片、宠物详情弹窗）里呈现为
// 方形且观感差，远不如程序化 PetSprite 艺术干净统一。为全局一致 + 圆形展示美观，
// getPetImageUrl 恒返回 null，PetSprite 一律回退到程序化 SVG。
//
// 如需重新启用：把对应 speciesId 登记进 MANIFEST(webp) / SVG_MANIFEST(手绘svg)，
// 并在下方返回 `${import.meta.env.BASE_URL}pets/{seriesId}/{speciesId}-{stage}.webp|.svg`。

/**
 * 返回宠物图片 URL；当前恒为 null（停用手绘/AI 图，统一走程序化 SVG）。
 */
export function getPetImageUrl(_speciesId: string, _stage: string): null {
  return null
}
