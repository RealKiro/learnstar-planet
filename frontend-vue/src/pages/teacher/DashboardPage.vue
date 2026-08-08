<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'
import { useAppMode } from '@/composables/useAppMode'
import { getSeriesBySpeciesId, SERIES_SCENES } from '@/utils/petData'
import PetSprite from '@/components/pet/PetSprite.vue'
import PetDetailModal from '@/components/pet/PetDetailModal.vue'
import type { ApiResponse } from '@/types'

// ===== 模式（教师完整 / 教室端+班级码基础） =====
const { isClassroomMode, isTeacherMode } = useAppMode()

interface CardStudent {
  name: string
  student_no?: string
  score: number
  pet_name?: string
  pet_species: string
  pet_level: number
}
interface ClassOverviewData {
  class_name: string
  grade: string
  student_count: number
  total_score: number
  avg_pet_level: number
  peak_count: number
  star_student: CardStudent | null
  top5: CardStudent[]
  recent_news: Array<{
    icon: string
    text: string
  }>
  weekly_score: number
}

const data = ref<ClassOverviewData | null>(null)
const loading = ref(true)
// 点击宠物 SVG 打开的图鉴弹窗
const handbook = ref<{ speciesId: string; level: number; score: number } | null>(null)
function openHandbook(s: CardStudent) {
  if (!s.pet_species) return
  handbook.value = { speciesId: s.pet_species, level: s.pet_level, score: s.score }
}

// 冠亚季军奖牌色（金/银/铜）
const MEDALS = ['#F59E0B', '#A8B0B8', '#CD7F32']

// 积分均衡阈值：第1名与末位分差 ≤ 该值视为「积分相近」，切换领跑群展示
const TIGHT_THRESHOLD = 10

/** 顶部完全并列（同分）：无梯度可排，隐藏 TOP 榜单 */
const isAllTied = computed(() => {
  const t = data.value?.top5 || []
  return t.length >= 2 && t[0].score === t[t.length - 1].score
})

/** 积分相近但未完全并列：切换领跑群（不强行分冠亚季军） */
const isTight = computed(() => {
  const t = data.value?.top5 || []
  if (t.length < 2 || isAllTied.value) return false
  return t[0].score - t[t.length - 1].score <= TIGHT_THRESHOLD
})

/** 并列名次：同分同名次（1、1、3…），并标记该名次是否并列 */
function calcRanks(list: CardStudent[]): Array<{ rank: number; tied: boolean }> {
  const result: Array<{ rank: number; tied: boolean }> = []
  list.forEach((s, i) => {
    const prev = i > 0 ? list[i - 1].score : null
    const next = i < list.length - 1 ? list[i + 1].score : null
    const tied = (prev !== null && s.score === prev) || (next !== null && s.score === next)
    const rank = i > 0 && prev !== null && s.score === prev ? result[i - 1].rank : i + 1
    result.push({ rank, tied })
  })
  return result
}

const leaderRanks = computed(() => calcRanks(data.value?.top5 || []))

const starBg = computed(() => {
  if (!data.value?.star_student?.pet_species) return 'var(--gradient-primary)'
  const series = getSeriesBySpeciesId(data.value.star_student.pet_species)
  return series && SERIES_SCENES[series.id]?.bgGradient || 'var(--gradient-primary)'
})

// ===== 数据加载（按模式互斥） =====
const token = ref('')
let pollTimer: ReturnType<typeof setInterval> | null = null

async function fetchTeacherDashboard() {
  try {
    const res = await apiGet<ApiResponse<ClassOverviewData>>('/api/v1/teacher/dashboard')
    data.value = res.data
  } catch {
    // Demo data
    data.value = {
      class_name: '三年级一班',
      grade: '三年级',
      student_count: 42,
      total_score: 3840,
      avg_pet_level: 6.2,
      peak_count: 5,
      star_student: {
        name: '张小明',
        student_no: '1001',
        pet_name: '九尾天狐',
        pet_species: 'nine_tail_fox',
        pet_level: 12,
        score: 520,
      },
      top5: [
        { name: '张小明', student_no: '1001', score: 520, pet_name: '九尾天狐', pet_species: 'nine_tail_fox', pet_level: 12 },
        { name: '李小红', student_no: '1002', score: 480, pet_name: '喷火龙', pet_species: 'charmander', pet_level: 11 },
        { name: '王小刚', student_no: '1003', score: 410, pet_name: '大熊猫', pet_species: 'panda', pet_level: 9 },
        { name: '赵小丽', student_no: '1004', score: 380, pet_name: '亚古兽', pet_species: 'mecha_dragon', pet_level: 8 },
        { name: '刘小强', student_no: '1005', score: 350, pet_name: '独角兽', pet_species: 'unicorn', pet_level: 8 },
      ],
      recent_news: [
        { icon: '🎉', text: '孙七的【亚古兽】进化到了 Lv.8！' },
        { icon: '⭐', text: '周八的【独角兽】+15 分！' },
        { icon: '📝', text: '全班总积分突破 3,000！' },
        { icon: '🌟', text: '张小明【九尾天狐】达到传说级！' },
      ],
      weekly_score: 1260,
    }
  } finally {
    loading.value = false
  }
}

async function fetchClassDashboard() {
  try {
    const res = await apiGet<{ data: ClassOverviewData }>('/api/v1/display/dashboard', { params: { token: token.value } })
    data.value = res.data
  } catch { /* ignore */ } finally { loading.value = false }
}

onMounted(async () => {
  if (isClassroomMode.value) {
    token.value = sessionStorage.getItem('class_token') || ''
    if (!token.value) { loading.value = false; return }
    await fetchClassDashboard()
    pollTimer = setInterval(fetchClassDashboard, 10000)
  } else {
    await fetchTeacherDashboard()
  }
})

onUnmounted(() => {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
})
</script>

<template>
  <div class="overview-page">
    <div class="page-header">
      <h2 class="page-title">🏠 班级总览</h2>
      <span class="page-subtitle">
        {{ data?.class_name || '--' }} · {{ data?.grade || '--' }}
      </span>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载数据中...</p>
    </div>

    <template v-else-if="data">
      <!-- 三栏概览 -->
      <div class="overview-grid">
        <!-- 班级之星 -->
        <div class="o-card star-card" v-if="data.star_student" style="position:relative;">
          <div class="o-label">🏅 班级之星</div>
          <div v-if="data.star_student.student_no" style="position:absolute;top:20px;right:24px;font-size:11px;color:var(--color-text-secondary);background:var(--tint-2);padding:2px 10px;border-radius:8px;">学号 {{ data.star_student.student_no }}</div>
          <div class="star-display">
            <div
              class="star-avatar"
              :style="{ background: isClassroomMode ? 'linear-gradient(135deg,#f093fb,#f5576c)' : starBg, cursor: 'pointer' }"
              @click="openHandbook(data.star_student)"
              title="点击查看宠物介绍"
            >
              <div v-if="data.star_student.pet_species" style="width:100%;height:100%;border-radius:50%;overflow:hidden;">
                <PetSprite :species-id="data.star_student.pet_species" :level="data.star_student.pet_level" :animate="true" />
              </div>
              <span v-else class="star-emoji">🌟</span>
            </div>
            <div class="star-info">
              <div class="star-name">{{ data.star_student.name }}</div>
              <div class="star-pet">Lv.{{ data.star_student.pet_level }}</div>
              <div class="star-score">{{ data.star_student.score }} 分</div>
            </div>
          </div>
        </div>

        <!-- 班级概况 -->
        <div class="o-card">
          <div class="o-label">📊 班级概况</div>
          <div class="o-value">{{ (data.total_score || 0).toLocaleString() }}</div>
          <div class="o-sub">总积分 · 共 {{ data.student_count }} 人</div>
          <div class="o-stats-row">
            <div>
              <span class="stat-label">平均等级</span>
              <strong class="stat-val">{{ (data.avg_pet_level || 0).toFixed(1) }}</strong>
            </div>
            <div>
              <span class="stat-label">巅峰 Lv.10+</span>
              <strong class="stat-val peak">{{ data.peak_count }}</strong>
            </div>
            <div>
              <span class="stat-label">本周增长</span>
              <strong class="stat-val weekly">+{{ data.weekly_score }}</strong>
            </div>
          </div>
        </div>

        <!-- 最新动态 -->
        <div class="o-card">
          <div class="o-label">📢 最新动态</div>
          <div class="news-list">
            <div v-for="(news, i) in data.recent_news" :key="i" class="news-item">
              <span class="news-icon">{{ news.icon }}</span>
              <span class="news-text">{{ news.text }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- TOP 5：积分拉开时前三名大卡片 + 名次位；相近时领跑群；完全并列时整块隐藏 -->
      <div v-if="!isAllTied" class="top5-section">
        <div class="section-header">
          <span class="section-title">{{ isTight ? '🏆 积分领跑' : '🏆 班级 TOP 5' }}</span>
          <span v-if="isTight" class="top5-tight-tip">积分相近 · 并列不分先后</span>
        </div>

        <!-- 积分相近：领跑群（不标冠亚季军，并列名次） -->
        <div v-if="isTight" class="leader-row">
          <div
            v-for="(s, i) in data.top5"
            :key="s.name"
            class="leader-card"
            @click="openHandbook(s)"
            :title="s.pet_species ? '点击查看宠物介绍' : ''"
          >
            <span class="leader-rank">{{ leaderRanks[i].tied ? '并列第 ' + leaderRanks[i].rank + ' 名' : '第 ' + leaderRanks[i].rank + ' 名' }}</span>
            <div class="leader-avatar">
              <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
              <span v-else class="leader-emoji">🌟</span>
            </div>
            <div class="leader-name">{{ s.name }}</div>
            <div v-if="s.student_no" class="leader-no">学号 {{ s.student_no }}</div>
            <div class="leader-score">{{ s.score }} 分</div>
          </div>
        </div>

        <!-- 积分拉开：前三名大卡片 + 4-5 名次位 -->
        <template v-else>
          <div class="top3-row">
            <div
              v-for="(s, i) in data.top5.slice(0, 3)"
              :key="s.name"
              class="top3-card"
              :class="'medal--' + i"
              @click="openHandbook(s)"
              :title="s.pet_species ? '点击查看宠物介绍' : ''"
            >
              <span class="top3-medal">{{ ['🥇', '🥈', '🥉'][i] }}</span>
              <div class="top3-avatar" :style="{ '--medal': MEDALS[i] }">
                <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
                <span v-else class="top3-emoji">🌟</span>
              </div>
              <div class="top3-name">{{ s.name }}</div>
              <div v-if="s.student_no" class="top3-no">学号 {{ s.student_no }}</div>
              <div class="top3-level">Lv.{{ s.pet_level }}</div>
              <div class="top3-score">{{ s.score }} 分</div>
              <div class="top3-bar"><div class="top3-fill" :style="{ width: (s.score / data.top5[0].score) * 100 + '%' }"></div></div>
            </div>
          </div>
          <!-- 第 4-5 名次位 -->
          <div class="top2-row">
            <div
              v-for="(s, i) in data.top5.slice(3)"
              :key="s.name"
              class="top4-card"
              @click="openHandbook(s)"
              :title="s.pet_species ? '点击查看宠物介绍' : ''"
            >
              <span class="top4-rank">{{ i + 4 }}</span>
              <div class="top4-avatar">
                <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
                <span v-else class="top4-emoji">🌟</span>
              </div>
              <div class="top4-info">
                <div class="top4-name">{{ s.name }}</div>
                <div class="top4-score">{{ s.score }} 分</div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </template>

    <!-- 空状态（仅教师端：教室端 data 为空时保持空白） -->
    <div v-else-if="isTeacherMode" class="empty-state">
      <div class="empty-icon">📭</div>
      <p>暂未分配班级</p>
      <p style="font-size:13px;color:var(--color-text-secondary);margin-top:6px;">请联系管理员为你分配班级后使用</p>
    </div>

    <!-- 宠物图鉴弹窗 -->
    <PetDetailModal v-if="handbook" :species-id="handbook.speciesId" :level="handbook.level" :score="handbook.score" @close="handbook = null" />
  </div>
</template>

<style scoped>
.overview-page {
  max-width: 1100px;
}

.page-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 24px;
}
.page-title { font-size: 24px; font-weight: 700; margin: 0; }
.page-subtitle { font-size: 14px; color: var(--color-text-secondary); }

/* 三栏 */
.overview-grid {
  display: grid;
  grid-template-columns: 1.2fr 1.4fr 1fr;
  gap: 20px;
  margin-bottom: 28px;
}
.o-card {
  position: relative;
  overflow: hidden;
  background:
    radial-gradient(120% 90% at 0% 0%, rgba(79,70,229,0.05), transparent 55%),
    var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  padding: 24px;
  transition: all 0.25s ease;
}
.o-card::after {
  content: '';
  position: absolute; inset: 0;
  background-image: radial-gradient(var(--color-text-secondary) 0.8px, transparent 0.8px);
  background-size: 18px 18px;
  opacity: 0.05;
  pointer-events: none;
}
.o-card > * { position: relative; z-index: 1; }
.o-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 30px rgba(0,0,0,0.12), 0 3px 10px rgba(0,0,0,0.06);
}
.o-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 14px;
}
.o-value {
  font-size: 36px;
  font-weight: 800;
  line-height: 1;
  margin-bottom: 4px;
}
.o-sub {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 14px;
}
.o-stats-row {
  display: flex;
  gap: 20px;
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
}
.stat-label {
  display: block;
  font-size: 11px;
  color: var(--color-text-secondary);
  margin-bottom: 2px;
}
.stat-val {
  font-size: 18px;
  font-weight: 700;
}
.stat-val.peak { color: #8B5CF6; }
.stat-val.weekly { color: #10B981; }

/* 班级之星 */
.star-display {
  display: flex;
  align-items: center;
  gap: 16px;
}
.star-avatar {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
  box-shadow: 0 4px 14px rgba(0,0,0,0.12), 0 0 24px rgba(0,0,0,0.08);
}
.star-avatar::before {
  content: '';
  position: absolute; inset: -8px; border-radius: 50%;
  background: conic-gradient(from 0deg, var(--color-primary), transparent 40%, var(--color-secondary), transparent 70%, var(--color-primary));
  opacity: 0.3;
  animation: spinRing 14s linear infinite;
  z-index: 0;
}
.star-avatar > * { position: relative; z-index: 1; }
@keyframes spinRing { to { transform: rotate(360deg); } }
.star-emoji { font-size: 40px; filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2)); }
.star-name { font-size: 20px; font-weight: 700; }
.star-pet { font-size: 13px; color: var(--color-text-secondary); }
.star-score { font-size: 22px; font-weight: 800; color: var(--color-primary); }

/* 新闻 */
.news-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 160px;
  overflow-y: auto;
}
.news-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 11px;
  background: var(--color-bg);
  border: 1px solid transparent;
  border-radius: 12px;
  font-size: 13px;
  transition: all 0.2s ease;
}
.news-item:hover {
  border-color: color-mix(in srgb, var(--color-primary) 20%, transparent);
  background: color-mix(in srgb, var(--color-primary) 4%, var(--color-bg));
  transform: translateX(2px);
}
.news-icon { font-size: 16px; flex-shrink: 0; }
.news-text { color: var(--color-text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* TOP 5 */
.top5-section {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  padding: 20px 24px;
}
.section-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.section-title { font-size: 15px; font-weight: 700; }
.top5-tight-tip {
  font-size: 12px;
  color: var(--color-text-secondary);
  background: var(--tint-2);
  padding: 2px 10px;
  border-radius: 12px;
}

/* 积分相近 · 领跑群（并列名次，不标冠亚季军） */
.leader-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}
.leader-card {
  text-align: center;
  padding: 18px 12px 16px;
  border-radius: 18px;
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  cursor: pointer;
  transition: all 0.25s ease;
}
.leader-card:hover { transform: translateY(-3px); box-shadow: var(--shadow-md); }
.leader-rank {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  color: var(--color-primary);
  background: rgba(79,70,229,0.08);
  border: 1px solid rgba(79,70,229,0.22);
  padding: 2px 10px;
  border-radius: 12px;
  margin-bottom: 8px;
}
.leader-avatar {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 auto 8px;
  border: 2px solid color-mix(in srgb, var(--color-primary) 70%, transparent);
  background: var(--color-bg-card);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-primary) 10%, transparent), 0 6px 16px rgba(0,0,0,0.12);
}
.leader-emoji { font-size: 26px; }
.leader-name { font-size: 14px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.leader-no { font-size: 10px; color: var(--color-text-secondary); background: var(--tint-2); padding: 1px 8px; border-radius: 8px; display: inline-block; margin-top: 2px; }
.leader-score { font-size: 15px; font-weight: 800; color: var(--color-primary); margin-top: 2px; }
/* 前三名大卡片 */
.top3-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
.top3-card {
  position: relative;
  text-align: center;
  padding: 20px 12px 16px;
  border-radius: 18px;
  border: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 0.25s ease;
}
.top3-card:hover { transform: translateY(-3px); box-shadow: var(--shadow-md); }
.medal--0 { background: linear-gradient(180deg, rgba(245,158,11,0.08), transparent); border-color: rgba(245,158,11,0.25); }
.medal--1 { background: linear-gradient(180deg, rgba(168,176,184,0.08), transparent); border-color: rgba(168,176,184,0.25); }
.medal--2 { background: linear-gradient(180deg, rgba(205,127,50,0.08), transparent); border-color: rgba(205,127,50,0.25); }
.top3-medal { position: absolute; top: 8px; left: 12px; font-size: 22px; }
.top3-avatar {
  width: 90px;
  height: 90px;
  border-radius: 50%;
  overflow: hidden;
  margin: 6px auto 10px;
  border: 2.5px solid var(--medal, #6B7280);
  background: var(--color-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--medal, #6B7280) 12%, transparent), 0 6px 18px color-mix(in srgb, var(--medal, #6B7280) 26%, transparent);
}
.top3-emoji { font-size: 34px; }
.top3-name { font-size: 16px; font-weight: 700; }
.top3-no { font-size: 10px; color: var(--color-text-secondary); background: var(--tint-2); padding: 1px 8px; border-radius: 8px; display: inline-block; margin-top: 2px; }
.top3-level { font-size: 11px; color: var(--color-text-secondary); margin-top: 2px; }
.top3-score { font-size: 18px; font-weight: 800; color: var(--color-primary); margin-top: 2px; }
.top3-bar { height: 4px; background: var(--color-border); border-radius: 2px; overflow: hidden; margin-top: 8px; }
.top3-fill { height: 100%; border-radius: 2px; background: var(--gradient-primary); transition: width 0.5s ease; }
.medal--0 .top3-fill { background: linear-gradient(90deg, #F59E0B, #FCD34D); }

/* 4-5 名次位 */
.top2-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  margin-top: 12px;
}
.top4-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 14px;
  border: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 0.25s ease;
}
.top4-card:hover { transform: translateY(-2px); box-shadow: var(--shadow-sm); }
.top4-rank { font-size: 22px; font-weight: 800; color: var(--color-text-secondary); opacity: 0.6; min-width: 28px; }
.top4-avatar {
  width: 54px;
  height: 54px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid color-mix(in srgb, var(--color-primary) 35%, var(--color-border));
  background: var(--color-bg);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-primary) 7%, transparent);
}
.top4-emoji { font-size: 22px; }
.top4-info { min-width: 0; }
.top4-name { font-size: 14px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.top4-score { font-size: 13px; font-weight: 700; color: var(--color-primary); }

/* 加载/空 */
.loading-state, .empty-state {
  text-align: center;
  padding: 60px 24px;
  color: var(--color-text-secondary);
}
.loading-spinner {
  width: 36px; height: 36px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }
.empty-icon { font-size: 48px; margin-bottom: 8px; }

@media (max-width: 900px) {
  .overview-grid { grid-template-columns: 1fr; }
  .top3-row { grid-template-columns: 1fr; }
  .top2-row { grid-template-columns: 1fr; }
}
</style>
