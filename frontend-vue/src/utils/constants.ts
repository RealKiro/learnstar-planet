// ===== 通用工具函数 =====
// 历史：原「旧 11 级宠物进化系统」相关导出（PET_EVOLUTION_STAGES / getStageName /
// getStageEmoji(level) / getPetStageName / getPetStageDescription / getStageTitle /
// getPetEmoji / getPetStage）已于 2026-09-13 全部移除 —— 它们与新 12 级系统重复且无任何引用。
// 宠物阶段一律使用 @/utils/petData 的 getLevelStage / getLevelTitle / getStageIndex；
// 宠物阶段 emoji 使用 @/utils/stageEmoji 的 getStageEmoji(speciesId, level)。

export function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

export function avatarGradient(name: string): string {
  // 白字头像：所有 stop 须保证与 #FFF 对比 ≥4.5（AA），故统一用 700→800 深色阶
  const gradients = [
    'linear-gradient(135deg,#0F766E,#115E59)',
    'linear-gradient(135deg,#6D28D9,#5B21B6)',
    'linear-gradient(135deg,#BE185D,#9D174D)',
    'linear-gradient(135deg,#B45309,#92400E)',
    'linear-gradient(135deg,#047857,#065F46)',
    'linear-gradient(135deg,#1D4ED8,#1E40AF)',
  ]
  let hash = 0
  for (let i = 0; i < (name || '').length; i++) {
    hash = (hash * 31 + name.charCodeAt(i)) & 0xFFFFFFFF
  }
  return gradients[Math.abs(hash) % gradients.length]
}

export function platformLabel(platform: string): string {
  const map: Record<string, string> = {
    wechat: '💬 微信',
    wechat_work: '💼 企业微信',
    qq: '🐧 QQ',
    renren: '🌐 人人通',
    dingtalk: '🔷 钉钉',
    feishu: '🪶 飞书',
  }
  return map[platform] || platform
}

export function formatTime(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
