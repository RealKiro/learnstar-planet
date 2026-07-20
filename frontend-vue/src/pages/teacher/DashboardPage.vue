<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { getSpeciesEmoji, getSeriesBySpeciesId, SERIES_SCENES } from '@/utils/petData'
import type { ApiResponse } from '@/types'

interface ClassOverviewData {
  class_name: string
  grade: string
  student_count: number
  total_score: number
  avg_pet_level: number
  peak_count: number
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
  weekly_score: number
}

const data = ref<ClassOverviewData | null>(null)
const loading = ref(true)

const starBg = computed(() => {
  if (!data.value?.star_student?.pet_species) return 'var(--gradient-primary)'
  const series = getSeriesBySpeciesId(data.value.star_student.pet_species)
  return series && SERIES_SCENES[series.id]?.bgGradient || 'var(--gradient-primary)'
})

onMounted(async () => {
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
        pet_name: '九尾天狐',
        pet_species: 'nine_tail_fox',
        pet_level: 12,
        score: 520,
      },
      top5: [
        { name: '张小明', score: 520, pet_name: '九尾天狐', pet_species: 'nine_tail_fox', pet_level: 12 },
        { name: '李小红', score: 480, pet_name: '喷火龙', pet_species: 'charmander', pet_level: 11 },
        { name: '王小刚', score: 410, pet_name: '大熊猫', pet_species: 'panda', pet_level: 9 },
        { name: '赵小丽', score: 380, pet_name: '机械龙', pet_species: 'mecha_dragon', pet_level: 8 },
        { name: '刘小强', score: 350, pet_name: '独角兽', pet_species: 'unicorn', pet_level: 8 },
      ],
      recent_news: [
        { icon: '🎉', text: '孙七的【机械龙】进化到了 Lv.8！' },
        { icon: '⭐', text: '周八的【独角兽】+15 分！' },
        { icon: '📝', text: '全班总积分突破 3,000！' },
        { icon: '🌟', text: '张小明【九尾天狐】达到传说级！' },
      ],
      weekly_score: 1260,
    }
  } finally {
    loading.value = false
  }
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
        <div class="o-card star-card" v-if="data.star_student">
          <div class="o-label">🏅 班级之星</div>
          <div class="star-display">
            <div class="star-avatar" :style="{ background: starBg }">
              <span class="star-emoji">{{ data.star_student.pet_species ? getSpeciesEmoji(data.star_student.pet_species) : '🌟' }}</span>
            </div>
            <div class="star-info">
              <div class="star-name">{{ data.star_student.name }}</div>
              <div class="star-pet">{{ data.star_student.pet_name }} · Lv.{{ data.star_student.pet_level }}</div>
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

      <!-- TOP 5 排行榜 -->
      <div class="top5-section">
        <div class="section-header">
          <span class="section-title">🏆 班级 TOP 5</span>
        </div>
        <div class="top5-grid">
          <div
            v-for="(s, i) in data.top5"
            :key="s.name"
            class="top5-card"
            :class="'rank--' + (i + 1)"
          >
            <div class="rank-badge">{{ ['🥇', '🥈', '🥉', '4', '5'][i] }}</div>
            <div class="rank-avatar">
              <span class="rank-emoji">{{ s.pet_species ? getSpeciesEmoji(s.pet_species) : '🌟' }}</span>
            </div>
            <div class="rank-name">{{ s.name }}</div>
            <div class="rank-pet">{{ s.pet_name }} · Lv.{{ s.pet_level }}</div>
            <div class="rank-score">{{ s.score }} 分</div>
            <div class="rank-bar">
              <div class="bar-fill" :style="{ width: (s.score / data.top5[0].score) * 100 + '%' }"></div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">📭</div>
      <p>暂无班级数据</p>
    </div>
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
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  padding: 24px;
  transition: all 0.25s ease;
}
.o-card:hover {
  box-shadow: var(--shadow-md);
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
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.star-emoji { font-size: 32px; filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2)); }
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
  padding: 8px 10px;
  background: var(--color-bg);
  border-radius: 10px;
  font-size: 13px;
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
.section-header { margin-bottom: 16px; }
.section-title { font-size: 15px; font-weight: 700; }
.top5-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}
.top5-card {
  text-align: center;
  padding: 16px 12px;
  border-radius: 16px;
  border: 1px solid var(--color-border);
  transition: all 0.25s ease;
  position: relative;
  overflow: hidden;
}
.top5-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}
.rank--1 {
  background: linear-gradient(180deg, rgba(245,158,11,0.06), transparent);
  border-color: rgba(245,158,11,0.2);
}
.rank--2 {
  background: linear-gradient(180deg, rgba(148,163,184,0.04), transparent);
  border-color: rgba(148,163,184,0.12);
}
.rank--3 {
  background: linear-gradient(180deg, rgba(217,119,6,0.04), transparent);
  border-color: rgba(217,119,6,0.12);
}
.rank-badge { font-size: 24px; margin-bottom: 6px; }
.rank-avatar { margin-bottom: 6px; }
.rank-emoji { font-size: 28px; }
.rank-name { font-size: 14px; font-weight: 600; margin-bottom: 2px; }
.rank-pet { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 6px; }
.rank-score { font-size: 16px; font-weight: 700; color: var(--color-primary); margin-bottom: 6px; }
.rank-bar {
  height: 4px;
  background: var(--color-border);
  border-radius: 2px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  background: var(--gradient-primary);
  border-radius: 2px;
  transition: width 0.5s ease;
}
.rank--1 .bar-fill {
  background: linear-gradient(90deg, #F59E0B, #FCD34D);
}

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
  .top5-grid { grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); }
}
</style>
