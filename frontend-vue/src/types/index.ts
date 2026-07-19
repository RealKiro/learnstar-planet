export interface User {
  id: number
  name: string
  username: string
  nickname?: string
  role: 'school_admin' | 'teacher' | 'parent'
  avatar?: string
  phone?: string
  email?: string
  bindings?: string[]
  class_names?: string[]
}

export interface Student {
  id: number
  name: string
  student_no?: string
  gender?: string
  class_id: number
  class_name?: string
  class_grade?: string
  total_score: number
  status: 'active' | 'graduated' | 'disabled'
}

export interface ClassRoom {
  id: number
  name: string
  grade?: string
  year?: string
  teacher_id?: number
  teacher_name?: string
  student_count: number
  display_code?: string
}

export interface ScoreRule {
  id: number
  name: string
  points: number
  category: string
  is_penalty: boolean
  class_id?: number
}

export interface Pet {
  id: number
  student_id: number
  student_name?: string
  name: string
  species: string
  level: number
  exp: number
  mood: number
  stage_name?: string
}

/** 扩展的宠物详情（含可视化所需数据） */
export interface PetDetail extends Pet {
  /** 系列 ID */
  seriesId?: string
  /** 系列名称 */
  seriesName?: string
  /** 当前等级名称（如"衔烛之苗"） */
  levelName?: string
  /** 当前等级描述 */
  levelDescription?: string
  /** 阶段分类 */
  stage?: 'egg' | 'baby' | 'growing' | 'mature' | 'legendary'
  /** 所需累计积分 */
  requiredScore?: number
  /** 学生头像 */
  student_avatar?: string
  /** 今日已获得积分 */
  todayScore?: number
  /** 连击天数 */
  streakDays?: number
}

export interface LeaderboardEntry {
  rank: number
  student_id: number
  student_name: string
  score: number | string
  /** 宠物相关扩展 */
  pet_name?: string
  pet_level?: number
  pet_species?: string
  pet_emoji?: string
}

export interface ShopItem {
  id: number
  name: string
  cost: number
  stock: number
  image?: string
}

export interface ShopRedemption {
  id: number
  student_id: number
  student_name: string
  item_id: number
  item_name: string
  cost: number
  status: 'pending' | 'approved' | 'rejected' | 'delivered'
}

export interface Notice {
  id: number
  title: string
  content: string
  type: 'info' | 'event' | 'urgent'
  is_published: boolean
  read_count?: number
  created_at: string
}

export interface Broadcast {
  id: number
  content: string
  type: 'banner' | 'popup' | 'fullscreen'
  voice: boolean
  loop: boolean
  duration: number
  status: string
  created_at: string
}

export interface Attendance {
  student_id: number
  student_name: string
  status: 'present' | 'late' | 'leave' | 'absent'
  check_in_time?: string
}

export interface GradeEntry {
  student_id: number
  student_name: string
  subjects: Record<string, number>
  total: number
  average: number
  rank: number
}

export interface QuestionBank {
  id: number
  name: string
  subject: string
  question_count: number
  usage_count: number
}

// API 通用响应
export interface ApiResponse<T> {
  data: T
  message?: string
  meta?: {
    current_page: number
    last_page: number
    per_page: number
    total: number
  }
}

// 学校管理端统计
export interface SchoolOverview {
  class_count: number
  teacher_count: number
  student_count: number
  monthly_score: number
  month_label?: string
  score_trend_percent?: number
  parent_count?: number
}

// 年级升级预览
export interface GradeUpgradePreview {
  upgrade_classes: Array<{
    class_name: string
    old_grade: string
    new_name: string
    new_grade: string
    student_count: number
  }>
  graduate_classes: Array<{
    class_name: string
    student_count: number
  }>
  summary: {
    upgrade_class_count: number
    upgrade_student_count: number
    graduate_class_count: number
    graduate_student_count: number
    note?: string
  }
}
