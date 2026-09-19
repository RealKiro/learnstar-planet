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
const loadError = ref('')
// 点击宠物 SVG 打开的图鉴弹窗
const handbook = ref<{ speciesId: string; level: number; score: number } | null>(null)
function openHandbook(s: CardStudent) {
  if (!s.pet_species) return
  handbook.value = { speciesId: s.pet_species, level: s.pet_level, score: s.score }
}


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

/** 班级之星头像底色：宠物系列场景渐变属内容语义，保留；教室端兜底用品牌色 */
const starBg = computed(() => {
  if (isClassroomMode.value) return 'var(--ui-brand)'
  const species = data.value?.star_student?.pet_species
  if (!species) return 'var(--ui-brand)'
  const series = getSeriesBySpeciesId(species)
  return (series && SERIES_SCENES[series.id]?.bgGradient) || 'var(--ui-brand)'
})

// ===== 数据加载（按模式互斥） =====
const token = ref('')

// 空态文案：区分「没绑定班级」与「接口没拿到数据」，避免只留一个空白页
const emptyHint = computed(() =>
  isClassroomMode.value && !token.value
    ? '尚未进入班级，请扫描班级码或由老师在大屏端进入'
    : '数据暂时没有加载出来，请稍后重试',
)
let pollTimer: ReturnType<typeof setInterval> | null = null

async function fetchTeacherDashboard() {
  try {
    const res = await apiGet<ApiResponse<ClassOverviewData>>('/api/v1/teacher/dashboard')
    data.value = res.data
  } catch {
    // Demo data —— 兜底演示数据保留，但必须声明，否则老师会把假班情当真（与 ScoresPage 口径一致）
    loadError.value = '班级数据加载失败，当前展示的是演示数据'
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
      <span v-if="data" class="page-subtitle">
        {{ data.class_name || '--' }} · {{ data.grade || '--' }}
      </span>
    </div>

    <!-- 接口失败时下方渲染的是内置演示数据，必须显式声明，避免被误读为真实班情 -->
    <div v-if="loadError" class="error-banner error-banner--block">{{ loadError }}</div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载数据中...</p>
    </div>

    <div v-else-if="!data" class="empty-state">
      <div class="empty-state__icon">📭</div>
      <div class="empty-state__title">暂无班级数据</div>
      <p class="empty-state__desc">{{ emptyHint }}</p>
    </div>

    <template v-else-if="data">
      <!-- 三栏概览 -->
      <div class="overview-grid">
        <!-- 班级之星 -->
        <div class="o-card star-card tdb-relative" v-if="data.star_student">
          <div class="o-label">🏅 班级之星</div>
          <div v-if="data.star_student.student_no" class="tdb-corner-tag">学号 {{ data.star_student.student_no }}</div>
          <div class="star-display">
            <div
              class="star-avatar"
              :style="{ background: starBg, cursor: 'pointer' }"
              @click="openHandbook(data.star_student)"
              title="点击查看宠物介绍"
            >
              <div v-if="data.star_student.pet_species" class="tdb-orb">
                <PetSprite :species-id="data.star_student.pet_species" :level="data.star_student.pet_level" :animate="true" />
              </div>
              <span v-else class="star-emoji">🌟</span>
            </div>
            <div class="star-info">
              <div class="star-name">{{ data.star_student.name }}</div>
              <div class="star-pet tnum">Lv.{{ data.star_student.pet_level }}</div>
              <div class="star-score tnum">{{ data.star_student.score }} 分</div>
            </div>
          </div>
        </div>

        <!-- 班级概况 -->
        <div class="o-card">
          <div class="o-label">📊 班级概况</div>
          <div class="o-value tnum">{{ (data.total_score || 0).toLocaleString() }}</div>
          <div class="o-sub">总积分 · 共 {{ data.student_count }} 人</div>
          <div class="o-stats-row">
            <div>
              <span class="stat-label">平均等级</span>
              <strong class="stat-val tnum">{{ (data.avg_pet_level || 0).toFixed(1) }}</strong>
            </div>
            <div>
              <span class="stat-label">巅峰 Lv.10+</span>
              <strong class="stat-val peak tnum">{{ data.peak_count }}</strong>
            </div>
            <div>
              <span class="stat-label">本周增长</span>
              <strong class="stat-val weekly tnum">+{{ data.weekly_score }}</strong>
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
            <span class="leader-rank tnum">{{ leaderRanks[i].tied ? '并列第 ' + leaderRanks[i].rank + ' 名' : '第 ' + leaderRanks[i].rank + ' 名' }}</span>
            <div class="leader-avatar">
              <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
              <span v-else class="leader-emoji">🌟</span>
            </div>
            <div class="leader-name">{{ s.name }}</div>
            <div v-if="s.student_no" class="leader-no tnum">学号 {{ s.student_no }}</div>
            <div class="leader-score tnum">{{ s.score }} 分</div>
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
              <div class="top3-avatar">
                <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
                <span v-else class="top3-emoji">🌟</span>
              </div>
              <div class="top3-name">{{ s.name }}</div>
              <div v-if="s.student_no" class="top3-no tnum">学号 {{ s.student_no }}</div>
              <div class="top3-level tnum">Lv.{{ s.pet_level }}</div>
              <div class="top3-score tnum">{{ s.score }} 分</div>
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
              <span class="top4-rank tnum">{{ i + 4 }}</span>
              <div class="top4-avatar">
                <PetSprite v-if="s.pet_species" :species-id="s.pet_species" :level="s.pet_level" :animate="true" />
                <span v-else class="top4-emoji">🌟</span>
              </div>
              <div class="top4-info">
                <div class="top4-name">{{ s.name }}</div>
                <div class="top4-score tnum">{{ s.score }} 分</div>
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
      <p class="tdb-hint-13">请联系管理员为你分配班级后使用</p>
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
.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 22px;
  font-weight: 700;
  margin: 0;
  color: var(--ui-fg);
}
.page-subtitle { font-size: 14px; color: var(--ui-fg-muted); }

/* 三栏 */
.overview-grid {
  display: grid;
  grid-template-columns: 1.2fr 1.4fr 1fr;
  gap: 16px;
  margin-bottom: 22px;
}
.o-card {
  position: relative;
  background: var(--ui-card);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-r-xl);
  padding: 18px 20px;
  box-shadow: var(--ui-shadow-xs);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}
.o-card:hover {
  border-color: var(--ui-border-strong);
  box-shadow: var(--ui-shadow-sm);
}
.o-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--ui-fg-muted);
  font-weight: 600;
  letter-spacing: 0.02em;
  margin-bottom: 12px;
}
.o-value {
  font-size: 30px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 6px;
  color: var(--ui-fg);
  letter-spacing: -0.01em;
}
.o-sub {
  font-size: 12px;
  color: var(--ui-fg-muted);
  margin-bottom: 12px;
}
.o-stats-row {
  display: flex;
  gap: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--ui-border);
}
.stat-label {
  display: block;
  font-size: 11px;
  color: var(--ui-fg-muted);
  margin-bottom: 3px;
}
.stat-val {
  font-size: 18px;
  font-weight: 700;
  color: var(--ui-fg);
}
.stat-val.peak { color: var(--c-violet-deep); }
.stat-val.weekly { color: var(--color-success-text); }

/* 班级之星：名次色只落在边框上，不用渐变与光晕 */
.star-card {
  border-color: color-mix(in srgb, var(--c-amber) 32%, var(--ui-border));
}
.star-card:hover {
  border-color: color-mix(in srgb, var(--c-amber) 50%, var(--ui-border));
}
.star-display {
  display: flex;
  align-items: center;
  gap: 16px;
}
.star-avatar {
  position: relative;
  width: 88px;
  height: 88px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
  color: var(--ui-brand-fg);
  border: 1px solid var(--ui-border);
  box-shadow: 0 1px 2px rgba(9, 9, 11, 0.08), 0 0 0 4px var(--ui-muted);
}
.star-name { font-size: 18px; font-weight: 700; color: var(--ui-fg); }
.star-pet { font-size: 12px; color: var(--ui-fg-muted); }
.star-score { font-size: 20px; font-weight: 700; color: var(--ui-brand); margin-top: 2px; }

/* 新闻 */
.news-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 160px;
  overflow-y: auto;
}
.news-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  background: var(--ui-bg-subtle);
  border: 1px solid transparent;
  border-radius: var(--ui-r-md);
  font-size: 13px;
  transition: border-color 0.2s ease, background 0.2s ease;
}
.news-item:hover {
  border-color: var(--ui-border);
  background: var(--ui-muted);
}
.news-icon { font-size: 16px; line-height: 1; flex-shrink: 0; }
.news-text { color: var(--ui-fg-muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* TOP 5 */
.top5-section {
  background: var(--ui-card);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-r-xl);
  padding: 16px 20px 20px;
  box-shadow: var(--ui-shadow-xs);
}
.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
  font-weight: 700;
  color: var(--ui-fg);
}
.top5-tight-tip {
  font-size: 12px;
  color: var(--ui-fg-muted);
  background: var(--ui-muted);
  padding: 2px 10px;
  border-radius: var(--ui-r-sm);
}

/* 积分相近 · 领跑群（并列名次，不标冠亚季军） */
.leader-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}
.leader-card {
  text-align: center;
  padding: 16px 12px 14px;
  border-radius: var(--ui-r-xl);
  border: 1px solid var(--ui-border);
  background: var(--ui-card);
  cursor: pointer;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}
.leader-card:hover {
  border-color: var(--ui-brand-border);
  background: var(--ui-brand-soft);
  box-shadow: var(--ui-shadow-sm);
}
.leader-rank {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  color: var(--ui-brand);
  background: var(--ui-brand-soft);
  border: 1px solid var(--ui-brand-border);
  padding: 2px 10px;
  border-radius: var(--ui-r-sm);
  margin-bottom: 10px;
}
.leader-avatar {
  width: 68px;
  height: 68px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 auto 8px;
  border: 1px solid var(--ui-brand-border);
  background: var(--ui-bg-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ui-fg-subtle);
  box-shadow: 0 0 0 3px var(--ui-brand-soft);
}
.leader-name { font-size: 14px; font-weight: 600; color: var(--ui-fg); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.leader-no { font-size: 10px; color: var(--ui-fg-muted); background: var(--ui-muted); padding: 1px 8px; border-radius: var(--ui-r-sm); display: inline-block; margin-top: 3px; }
.leader-score { font-size: 15px; font-weight: 700; color: var(--ui-brand); margin-top: 3px; }

/* 前三名大卡片：名次色统一走 --medal，暗色自动跟随 */
.top3-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
.top3-card {
  --medal: var(--c-amber);
  --medal-soft: var(--c-amber-chip);
  position: relative;
  text-align: center;
  padding: 20px 12px 16px;
  border-radius: var(--ui-r-xl);
  border: 1px solid color-mix(in srgb, var(--medal) 28%, var(--ui-border));
  background: var(--ui-card);
  cursor: pointer;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}
.medal--1 { --medal: var(--c-slate-500); --medal-soft: var(--c-gray-chip); }
.medal--2 { --medal: var(--c-orange); --medal-soft: var(--c-amber-chip); }
.top3-card:hover {
  border-color: color-mix(in srgb, var(--medal) 55%, var(--ui-border));
  box-shadow: var(--ui-shadow-sm);
}
.top3-medal { position: absolute; top: 10px; left: 12px; font-size: 20px; line-height: 1; }
.top3-avatar {
  width: 82px;
  height: 82px;
  border-radius: 50%;
  overflow: hidden;
  margin: 6px auto 10px;
  border: 1px solid color-mix(in srgb, var(--medal) 40%, var(--ui-border));
  background: var(--ui-bg-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ui-fg-subtle);
  box-shadow: 0 0 0 4px var(--medal-soft);
}
.top3-name { font-size: 16px; font-weight: 700; color: var(--ui-fg); }
.top3-no { font-size: 10px; color: var(--ui-fg-muted); background: var(--ui-muted); padding: 1px 8px; border-radius: var(--ui-r-sm); display: inline-block; margin-top: 3px; }
.top3-level { font-size: 11px; color: var(--ui-fg-muted); margin-top: 3px; }
.top3-score { font-size: 18px; font-weight: 700; color: var(--ui-fg); margin-top: 3px; }
.top3-bar { height: 4px; background: var(--ui-muted); border-radius: 2px; overflow: hidden; margin-top: 10px; }
.top3-fill { height: 100%; border-radius: 2px; background: var(--medal); transition: width 0.5s ease; }

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
  border-radius: var(--ui-r-lg);
  border: 1px solid var(--ui-border);
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
}
.top4-card:hover { border-color: var(--ui-border-strong); background: var(--ui-bg-subtle); }
.top4-rank { font-size: 20px; font-weight: 700; color: var(--ui-fg-muted); min-width: 26px; }
.top4-avatar {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid var(--ui-border);
  background: var(--ui-bg-subtle);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ui-fg-subtle);
}
.top4-info { min-width: 0; }
.top4-name { font-size: 14px; font-weight: 600; color: var(--ui-fg); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.top4-score { font-size: 13px; font-weight: 600; color: var(--ui-fg-muted); }

/* 加载/空 */
.loading-state, .empty-state {
  text-align: center;
  padding: 60px 24px;
  color: var(--ui-fg-muted);
}
.loading-spinner {
  width: 36px; height: 36px;
  border: 2px solid var(--ui-border);
  border-top-color: var(--ui-brand);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }
.empty-icon { font-size: 48px; line-height: 1; text-align: center; margin-bottom: 8px; }

@media (max-width: 900px) {
  .overview-grid { grid-template-columns: 1fr; }
  .top3-row { grid-template-columns: 1fr; }
  .top2-row { grid-template-columns: 1fr; }
}
.tdb-relative { position:relative; }
.tdb-corner-tag { position:absolute;top:20px;right:24px;font-size:11px;color:var(--ui-fg-muted);background:var(--ui-muted);padding:2px 10px;border-radius:var(--ui-r-sm); }
.tdb-orb { width:100%;height:100%;border-radius:50%;overflow:hidden; }
.tdb-hint-13 { font-size:13px;color:var(--ui-fg-muted);margin-top:6px; }
/* 宠物缺失时的 emoji 占位（字号沿用 HEAD 版，与头像容器 88/68/82/52 相匹配） */
.star-emoji { font-size: 40px; line-height: 1; filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2)); }
.leader-emoji { font-size: 26px; line-height: 1; }
.top3-emoji { font-size: 34px; line-height: 1; }
.top4-emoji { font-size: 22px; line-height: 1; }
</style>
