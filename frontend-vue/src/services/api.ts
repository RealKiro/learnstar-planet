/**
 * 学趣星球 · API 服务层
 * 集中管理所有后端 API 调用，提供类型安全的接口
 * 对接文档：功能分类与版面设计.txt
 */

import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api'
import type { ApiResponse, Student, ScoreRule, Pet, PetDetail, LeaderboardEntry, ClassRoom } from '@/types'

// ============================================================
// 1. 认证模块
// ============================================================

export interface LoginRequest {
  username: string
  password: string
  role?: 'teacher' | 'parent' | 'school_admin'
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    name: string
    username: string
    nickname?: string
    role: string
    avatar?: string
    class_names?: string[]
  }
}

/** 教师登录 */
export const authApi = {
  login: (data: LoginRequest) =>
    apiPost<ApiResponse<LoginResponse>>('/api/v1/auth/teacher/login', data),

  /** 班级码登录（学生端/班级大屏统一入口） */
  classLogin: (classCode: string) =>
    apiPost<ApiResponse<{ token: string; class_id: number; class_name: string; grade: string; student_count: number }>>('/api/v1/auth/class/login', { class_code: classCode }),

  logout: () => apiPost('/api/v1/auth/logout'),
  changePassword: (data: { old_password: string; new_password: string }) =>
    apiPost('/api/v1/auth/change-password', data),
}

// ============================================================
// 2. 班级总览模块
// ============================================================

export interface ClassOverviewData {
  class_name: string
  grade: string
  student_count: number
  total_score: number
  avg_pet_level: number
  peak_count: number
  weekly_score: number
  star_student: {
    name: string
    pet_name: string
    pet_species: string
    pet_level: number
    score: number
  } | null
  top5: Array<{
    name: string
    score: number
    pet_name: string
    pet_species: string
    pet_level: number
  }>
  recent_news: Array<{
    icon: string
    text: string
  }>
}

/** 班级总览 */
export const overviewApi = {
  /** 获取班级总览数据（含班级之星、TOP5、动态） */
  getOverview: () =>
    apiGet<ApiResponse<ClassOverviewData>>('/api/v1/teacher/dashboard'),

  /** 获取班级最新动态 */
  getRecentNews: () =>
    apiGet<ApiResponse<{ icon: string; text: string }[]>>('/api/v1/teacher/dashboard/news'),
}

// ============================================================
// 3. 课堂评价模块
// ============================================================

export interface GiveScoreRequest {
  student_id: number
  points: number
  reason: string
}

export interface BatchGiveScoreRequest {
  student_ids: number[]
  points: number
  reason: string
}

export interface ScoreSummary {
  total: number
  today: number
  this_week: number
}

/** 课堂评价 */
export const scoreApi = {
  /** 获取学生列表 */
  getStudents: (params?: { per_page?: number; class_id?: number }) =>
    apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students', { params }),

  /** 获取积分规则 */
  getRules: () =>
    apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules'),

  /** 获取积分汇总 */
  getSummary: () =>
    apiGet<ApiResponse<ScoreSummary>>('/api/v1/teacher/scores/summary'),

  /** 单个加分/减分 */
  giveScore: (data: GiveScoreRequest) =>
    apiPost<ApiResponse<{ new_score: number }>>('/api/v1/teacher/scores/give', data),

  /** 批量按规则加减分 */
  batchGiveScore: (data: BatchGiveScoreRequest) =>
    apiPost('/api/v1/teacher/scores/batch-give', data),

  /** 获取积分记录 */
  getHistory: (studentId: number) =>
    apiGet<ApiResponse<{ id: number; points: number; reason: string; created_at: string; balance_after: number }[]>>(
      `/api/v1/teacher/scores/history/${studentId}`),
}

// ============================================================
// 4. 宠物系统模块
// ============================================================

export interface FeedPetRequest {
  pet_id: number
  item?: string
}

/** 宠物系统 */
export const petApi = {
  /** 获取全班宠物概览 */
  getClassOverview: () =>
    apiGet<ApiResponse<Pet[]>>('/api/v1/teacher/pets/overview'),

  /** 获取单个宠物详情 */
  getPetDetail: (petId: number) =>
    apiGet<ApiResponse<PetDetail>>(`/api/v1/teacher/pets/${petId}`),

  /** 投喂宠物 */
  feedPet: (data: FeedPetRequest) =>
    apiPost<ApiResponse<{ exp_gained: number; new_level: number }>>('/api/v1/teacher/pets/${studentId}/feed', data),

  /** 重命名宠物 */
  renamePet: (petId: number, name: string) =>
    apiPut(`/api/v1/teacher/pets/${petId}/rename`, { name }),

  /** 获取宠物系列列表 */
  getSeries: () =>
    apiGet<ApiResponse<{ id: string; name: string; icon: string; description: string; species: any[] }[]>>('/api/v1/common/pet-types'),

  /** 切换班级宠物系列 */
  switchSeries: (seriesId: string) =>
    apiPost('/api/v1/teacher/class/switch-series', { series_id: seriesId }),
}

// ============================================================
// 5. 年级战场 (PK) 模块
// ============================================================

export interface ClassPKData {
  name: string
  totalScore: number
  studentCount: number
  avgLevel: number
  peakCount: number
  weekGrowth: number
  isOwn: boolean
}

/** 年级战场 */
export const pkApi = {
  /** 获取同年级各班PK数据 */
  getLeaderboard: () =>
    apiGet<ApiResponse<ClassPKData[]>>('/api/v1/teacher/pk/leaderboard'),

  /** 获取本班战力详情 */
  getMyStats: () =>
    apiGet<ApiResponse<{ totalScore: number; avgLevel: number; peakCount: number; weekGrowth: number; rank: number }>>(
      '/api/v1/teacher/pk/my-stats'),

  /** 发起PK挑战 */
  challenge: (targetClassId: number) =>
    apiPost<ApiResponse<{ target_class: string; expires_at: string }>>('/api/v1/teacher/pk/challenge', { target_class_id: targetClassId }),
}

// ============================================================
// 6. 排行榜模块
// ============================================================

export const leaderboardApi = {
  /** 总积分榜 */
  getTotal: () =>
    apiGet<ApiResponse<LeaderboardEntry[]>>('/api/v1/teacher/leaderboard/total'),

  /** 进步最快榜 */
  getWeekly: () =>
    apiGet<ApiResponse<LeaderboardEntry[]>>('/api/v1/teacher/leaderboard/weekly'),

  /** 宠物等级榜 */
  getPetLevel: () =>
    apiGet<ApiResponse<LeaderboardEntry[]>>('/api/v1/teacher/leaderboard/pet-level'),
}

// ============================================================
// 7. 宠物图鉴模块
// ============================================================

export const handbookApi = {
  /** 获取全部系列数据 */
  getAllSeries: () =>
    apiGet<ApiResponse<any[]>>('/api/v1/common/pet-types'),

  /** 获取单个系列详情 */
  getSeriesDetail: (seriesId: string) =>
    apiGet<ApiResponse<any>>(`/api/v1/common/pet-types/${seriesId}`),
}

// ============================================================
// 8. 班级管理模块
// ============================================================

export const classApi = {
  /** 获取班级信息 */
  getClassInfo: () =>
    apiGet<ApiResponse<ClassRoom>>('/api/v1/teacher/class'),

  /** 获取积分规则列表 */
  getRules: () =>
    apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules'),

  /** 创建积分规则 */
  createRule: (data: { name: string; points: number; category: string; is_penalty: boolean }) =>
    apiPost('/api/v1/teacher/scores/rules', data),

  /** 切换班级宠物系列 */
  switchSeries: (seriesId: string) =>
    apiPost('/api/v1/teacher/class/switch-series', { series_id: seriesId }),

  /** 获取班级当前积分 */
  getClassScore: () =>
    apiGet<ApiResponse<{ total_score: number; class_points: number }>>('/api/v1/teacher/class/score'),
}

// ============================================================
// 9. 家长端 API
// ============================================================

export const parentApi = {
  /** 获取孩子列表和概览 */
  getHome: () =>
    apiGet<ApiResponse<{ children: any[] }>>('/api/v1/parent/home'),

  /** 获取宠物详情 */
  getPet: (studentId: number) =>
    apiGet<ApiResponse<PetDetail>>(`/api/v1/parent/pet?student_id=${studentId}`),

  /** 获取成长日志 */
  getGrowthLog: (studentId: number) =>
    apiGet<ApiResponse<any[]>>(`/api/v1/parent/growth/log?student_id=${studentId}`),

  /** 获取排行榜 */
  getRanking: (studentId: number) =>
    apiGet<ApiResponse<any[]>>(`/api/v1/parent/ranking?student_id=${studentId}`),
}
