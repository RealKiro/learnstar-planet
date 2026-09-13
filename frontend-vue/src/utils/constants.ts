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
  const gradients = [
    'linear-gradient(135deg,#6366F1,#818CF8)',
    'linear-gradient(135deg,#8B5CF6,#A78BFA)',
    'linear-gradient(135deg,#EC4899,#F472B6)',
    'linear-gradient(135deg,#F59E0B,#FCD34D)',
    'linear-gradient(135deg,#10B981,#6EE7B7)',
    'linear-gradient(135deg,#3B82F6,#60A5FA)',
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
