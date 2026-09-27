// ===== 学宠星球 · 积分规则共享常量 =====
// ⚠️ 唯一真源：分类键集必须与后端 App\Services\ScoreRuleService::CATEGORY_LABELS 保持一致
//    （后端 DEFAULT_RULES 里所有规则的 category 都取自该键集）。
//    此前 RulesPage / AdminScoreRulesPage 各自复制过一份本表，已收口为统一 import。

/** 规则分类 → 展示标签（含 emoji） */
export const categoryLabels: Record<string, string> = {
  classroom: '📖 课堂表现',
  homework: '📝 作业管理',
  behavior: '🌟 行为习惯',
  literacy: '📊 综合素养',
  daily: '📅 日常表现',
  academic: '📚 学业表现',
  custom: '✨ 自定义',
}

/** 未知分类兜底标签 */
export function categoryLabel(category: string): string {
  return categoryLabels[category] || category || '📌 其他'
}
