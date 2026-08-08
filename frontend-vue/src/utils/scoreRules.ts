// ===== 学趣星球 · 积分规则共享常量 =====
// 供 ScoresPage 弹窗 / RulesPage 规则管理复用，保持分类标签一致。

/** 规则分类 → 展示标签（含 emoji） */
export const categoryLabels: Record<string, string> = {
  classroom: '📖 课堂表现',
  homework: '📝 作业管理',
  behavior: '🌟 行为习惯',
  literacy: '📊 综合素养',
  daily: '📅 日常表现',
  academic: '📚 学业',
}

/** 未知分类兜底标签 */
export function categoryLabel(category: string): string {
  return categoryLabels[category] || category || '📌 其他'
}
