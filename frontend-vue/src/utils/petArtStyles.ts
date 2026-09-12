// ===== 宠物画风预设（学生可自由切换）=====
// 每种画风 = 一条可拼进 AI 生图提示词的风格行（中/英）；
// 学生选择持久化在 localStorage（pet_art_style），生图提示词与未来图片展示均按此风格出图。

export interface PetArtStyle {
  id: string
  name: string
  emoji: string
  /** 中文风格行（拼进提示词） */
  cn: string
  /** 英文风格行（拼进提示词给英文模型） */
  en: string
}

export const PET_ART_STYLES: PetArtStyle[] = [
  { id: 'inkcn',      name: '中国风水墨', emoji: '🖼️', cn: '中国风水墨画，工笔淡彩，留白意境，墨韵晕染',           en: 'Chinese ink-wash painting, gongbi light color, negative space, ink diffusion' },
  { id: 'anime',      name: '动漫',       emoji: '🎬', cn: '日系动漫风格，赛璐璐质感，明快高饱和，热血治愈',       en: 'Japanese anime style, cel-shaded, vibrant saturated colors' },
  { id: 'cg3d',       name: '3D渲染CG',  emoji: '🎬', cn: '3D 渲染 CG 动画风格，次时代 PBR 材质，电影级光影',     en: '3D CG animation render, PBR materials, cinematic lighting' },
  { id: 'cyberpunk',  name: '赛博朋克',   emoji: '🌆', cn: '赛博朋克风格，霓虹光效，机械义体质感，暗调未来都市',   en: 'cyberpunk style, neon glow, cybernetic textures, dark futuristic city' },
  { id: 'watercolor', name: '水彩画',     emoji: '🎨', cn: '水彩画风格，水痕晕染，柔和通透，纸纹质感',             en: 'watercolor painting, soft washes, translucent layers, paper texture' },
  { id: 'cartoon',    name: '卡通',       emoji: '🧸', cn: '卡通风格，圆润造型，扁平明快，线条干净',               en: 'cartoon style, rounded shapes, flat vivid colors, clean linework' },
  { id: 'pixel',      name: '像素',       emoji: '🕹️', cn: '像素风格，16-bit 点阵，复古游戏质感，有限调色板',      en: 'pixel art, 16-bit sprite, retro game aesthetic, limited palette' },
  { id: 'figure',     name: '手办',       emoji: '🗿', cn: '手办模型风格，PVC 质感，棚拍布光，微距景深',           en: 'collectible figure style, PVC material, studio lighting, macro depth of field' },
]

export const DEFAULT_ART_STYLE = 'anime'

export function getArtStyleById(id: string): PetArtStyle {
  return PET_ART_STYLES.find(s => s.id === id) ?? PET_ART_STYLES.find(s => s.id === DEFAULT_ART_STYLE)!
}
