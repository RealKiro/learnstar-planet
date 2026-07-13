export const PET_EVOLUTION_STAGES = [
  { level: 0, emoji: '🥚', name: '宠物蛋' },
  { level: 1, emoji: '🐣', name: '绒绒雏鸟' },
  { level: 2, emoji: '🐥', name: '黄毛小鸭' },
  { level: 3, emoji: '🐰', name: '绒耳萌兔' },
  { level: 4, emoji: '🦊', name: '机灵小狐' },
  { level: 5, emoji: '🐱', name: '优雅萌猫' },
  { level: 6, emoji: '🐶', name: '忠诚幼犬' },
  { level: 7, emoji: '🦁', name: '威风小狮' },
  { level: 8, emoji: '🐯', name: '勇猛虎崽' },
  { level: 9, emoji: '🦄', name: '神圣灵兽' },
  { level: 10, emoji: '🐉', name: '东方神龙' },
]

export function getStageName(level: number): string {
  if (level < 0 || level >= PET_EVOLUTION_STAGES.length) return '星尘'
  return PET_EVOLUTION_STAGES[level].name
}

export function getStageEmoji(level: number): string {
  if (level < 0 || level >= PET_EVOLUTION_STAGES.length) return '🌟'
  return PET_EVOLUTION_STAGES[level].emoji
}

export function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
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
