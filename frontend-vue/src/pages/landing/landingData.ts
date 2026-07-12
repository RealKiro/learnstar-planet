export interface Slide {
  badge: string
  title: string
  highlight: string
  desc: string
  icon: string
}

export interface Feature {
  icon: string
  title: string
  desc: string
}

export interface Stage {
  emoji: string
  name: string
}

export interface Platform {
  key: string
  label: string
  icon: string
  color: string
}

export const slides: Slide[] = [
  {
    badge: 'MIT 开源 · 完全免费 · 自托管',
    title: '让每个孩子的努力', highlight: '都被看见',
    desc: '积分激励 · 宠物养成 · AI 助教 · 多端同步\n开源班级管理系统，数据完全自主掌控', icon: '🌌',
  },
  {
    badge: '12 大功能模块',
    title: '覆盖班级管理', highlight: '全场景',
    desc: '积分规则 · 宠物进化 · 排行榜 · 通知公告\n考勤 · 作业 · 答题 · 成绩 · 商城 · 广播 · AI', icon: '⚡',
  },
  {
    badge: '11 阶宠物进化',
    title: '积分变经验', highlight: '驱动成长',
    desc: '星尘 → 月芽 → 灵苗 → 青藤 → 慧树 → 蝶灵\n→ 鹰慧 → 狮睿 → 灵角 → 星耀 → 银河', icon: '🌟',
  },
  {
    badge: 'Docker 一键部署',
    title: '4 种数据库', highlight: '自由选择',
    desc: 'MySQL · PostgreSQL · SQLite · MariaDB\n内置 Redis 缓存，支持多端登录', icon: '🐳',
  },
]

export const features: Feature[] = [
  { icon: '⭐', title: '积分激励', desc: '自定义规则，实时加减分，进步看得见' },
  { icon: '🌟', title: '宠物进化', desc: '11 阶宇宙进化，积分驱动成长' },
  { icon: '🏆', title: '排行竞技', desc: '总积分 / 周进步 / 宠物等级三大排行' },
  { icon: '📢', title: '班级通知', desc: '一键发布，实时推送家长端' },
  { icon: '📊', title: '成绩管理', desc: '录入分析，班级对比，趋势可视化' },
  { icon: '🤖', title: 'AI 助教', desc: '班级反馈、学生分析、家校沟通建议' },
  { icon: '✅', title: '智能考勤', desc: '一键签到，到课 / 请假 / 迟到统计' },
  { icon: '📱', title: '扫码收作业', desc: '生成二维码，学生扫码提交自动汇总' },
  { icon: '🛍️', title: '积分商城', desc: '兑换实物 / 特权，教师审批发放' },
  { icon: '📡', title: '实时广播', desc: '消息直达桌面，文字 / 语音 / 横幅 / 全屏' },
  { icon: '📝', title: '在线答题', desc: '题库管理，课堂检测，自动判分统计' },
  { icon: '🔗', title: '多端登录', desc: '微信 / 企微 / QQ / 人人通，账号密码双通道' },
]

export const stages: Stage[] = [
  { emoji: '🌟', name: '星尘' }, { emoji: '🌙', name: '月芽' },
  { emoji: '🌱', name: '灵苗' }, { emoji: '🌿', name: '青藤' },
  { emoji: '🌳', name: '慧树' }, { emoji: '🦋', name: '蝶灵' },
  { emoji: '🦅', name: '鹰慧' }, { emoji: '🦁', name: '狮睿' },
  { emoji: '🦄', name: '灵角' }, { emoji: '✨', name: '星耀' },
  { emoji: '🌌', name: '银河' },
]

export const platforms: Platform[] = [
  { key: 'wechat', label: '微信', icon: '💬', color: '#07C160' },
  { key: 'wechat_work', label: '企业微信', icon: '💼', color: '#2B7CE9' },
  { key: 'qq', label: 'QQ', icon: '🐧', color: '#12B7F5' },
  { key: 'renren', label: '人人通', icon: '🌐', color: '#FF6A00' },
]
