<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet } from '@/utils/api'
import { getSpeciesEmoji } from '@/utils/petData'
import type { ApiResponse, Student } from '@/types'

// ===== 类型 =====
interface ClassPKData {
  name: string
  totalScore: number
  studentCount: number
  avgLevel: number
  peakCount: number
  weekGrowth: number
  isOwn: boolean
}

interface PKOverview {
  totalScore: number
  avgLevel: number
  peakCount: number
  weekGrowth: number
  rank: number
}

const router = useRouter()
function goToLeaderboard() { router.push({ name: 'teacher-leaderboard' }) }

// ===== 数据 =====
const students = ref<Student[]>([])
const classes = ref<ClassPKData[]>([])
const overview = ref<PKOverview>({
  totalScore: 0, avgLevel: 0, peakCount: 0, weekGrowth: 128, rank: 1,
})
const loading = ref(true)

const LEVEL_SCORES = [0, 15, 35, 60, 90, 125, 165, 210, 260, 315, 375, 450]

function calcLevel(score: number): number {
  let lv = 1
  for (let i = LEVEL_SCORES.length - 1; i >= 0; i--) {
    if (score >= LEVEL_SCORES[i]) { lv = i + 1; break }
  }
  return Math.min(lv, 12)
}

const maxScore = computed(() => Math.max(...classes.value.map(c => c.totalScore), 1))

const rankBadge = computed(() => {
  const r = overview.value.rank
  if (r === 1) return '🥇'
  if (r === 2) return '🥈'
  if (r === 3) return '🥉'
  return `#${r}`
})

const gapToFirst = computed(() => {
  const first = classes.value[0]
  if (!first || first.isOwn) return 0
  return first.totalScore - overview.value.totalScore
})

function getRankMedal(idx: number): string {
  return ['🥇', '🥈', '🥉'][idx] || `#${idx + 1}`
}

onMounted(async () => {
  try {
    const [sRes, pkRes] = await Promise.all([
      apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students?per_page=100'),
      apiGet<ApiResponse<ClassPKData[]>>('/api/v1/teacher/pk/leaderboard'),
    ])
    students.value = sRes.data || []
    classes.value = pkRes.data || []
  } catch {
    // Demo data
    const names = ['张小明', '李小红', '王小刚', '赵小丽', '刘小强', '陈小美', '周小龙']
    students.value = names.map((name, i) => ({
      id: i + 1, name, total_score: Math.floor(Math.random() * 400) + 50,
      class_id: 1, status: 'active' as const,
    }))
    const total = students.value.reduce((s, v) => s + v.total_score, 0)
    const avg = students.value.reduce((s, v) => s + calcLevel(v.total_score), 0) / students.value.length
    const peak = students.value.filter(s => calcLevel(s.total_score) >= 10).length
    overview.value = { totalScore: total, avgLevel: avg, peakCount: peak, weekGrowth: 128, rank: 2 }
    classes.value = [
      { name: '三年级二班', totalScore: 4200, studentCount: 45, avgLevel: 7.2, peakCount: 8, weekGrowth: 156, isOwn: false },
      { name: '三年级一班 (本班)', totalScore: total, studentCount: students.value.length, avgLevel: avg, peakCount: peak, weekGrowth: 128, isOwn: true },
      { name: '三年级三班', totalScore: 3800, studentCount: 43, avgLevel: 6.8, peakCount: 5, weekGrowth: 112, isOwn: false },
      { name: '三年级四班', totalScore: 3500, studentCount: 44, avgLevel: 6.1, peakCount: 4, weekGrowth: 98, isOwn: false },
      { name: '三年级五班', totalScore: 3100, studentCount: 42, avgLevel: 5.5, peakCount: 3, weekGrowth: 85, isOwn: false },
    ]
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="pk-page">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">🏆 年级战场</h2>
        <span class="page-subtitle">同年级大比拼</span>
      </div>
        <button class="page-link-btn" @click="goToLeaderboard">📊 查看个人排名</button>
      <div class="rank-badge">
        {{ rankBadge }} 当前排名: #{{ overview.rank }}
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载排行数据...</p>
    </div>

    <template v-else>
      <!-- PK 排行榜 -->
      <div class="pk-leaderboard">
        <div
          v-for="(cls, idx) in classes"
          :key="cls.name"
          class="pk-item"
          :class="{ 'pk-item--own': cls.isOwn }"
        >
          <div class="pk-rank" :class="{
            'rank-gold': idx === 0,
            'rank-silver': idx === 1,
            'rank-bronze': idx === 2,
          }">
            {{ getRankMedal(idx) }}
          </div>

          <div class="pk-class">
            <span class="pk-class-name">{{ cls.name }}</span>
            <span v-if="cls.isOwn" class="pk-own-badge">本班</span>
          </div>

          <div class="pk-bar-area">
            <div class="pk-bar-track">
              <div
                class="pk-bar-fill"
                :class="{ 'fill-gold': idx === 0 }"
                :style="{ width: (cls.totalScore / maxScore) * 100 + '%' }"
              ></div>
            </div>
            <span class="pk-bar-score">{{ cls.totalScore.toLocaleString() }}</span>
          </div>

          <div class="pk-stats">
            <span title="人数">👤 {{ cls.studentCount }}</span>
            <span title="平均等级">📈 {{ cls.avgLevel.toFixed(1) }}</span>
            <span title="巅峰人数">⭐ {{ cls.peakCount }}</span>
          </div>

          <button
            v-if="!cls.isOwn"
            class="pk-btn"
            @click="alert(`🚀 向 ${cls.name} 发起挑战！本周内总积分超过对方即可获胜！`)"
          >
            ⚔️ PK
          </button>
          <button v-else class="pk-btn own">
            🏠 本班
          </button>
        </div>
      </div>

      <!-- 本班战力和挑战建议 -->
      <div class="pk-detail-grid">
        <div class="pk-detail-card">
          <h4>📊 本班战力分析</h4>
          <div class="detail-row">
            <span>总积分</span>
            <span class="row-val">{{ overview.totalScore.toLocaleString() }}</span>
          </div>
          <div class="detail-row">
            <span>平均等级</span>
            <span class="row-val">{{ overview.avgLevel.toFixed(1) }}</span>
          </div>
          <div class="detail-row">
            <span>巅峰人数 (Lv.10+)</span>
            <span class="row-val peak">{{ overview.peakCount }}</span>
          </div>
          <div class="detail-row">
            <span>本周增长</span>
            <span class="row-val weekly">+{{ overview.weekGrowth }}</span>
          </div>
        </div>

        <div class="pk-detail-card">
          <h4>⚔️ 挑战建议</h4>
          <div v-if="gapToFirst > 0" class="suggestion">
            💪 距离第1名还差
            <strong class="gap-score">{{ gapToFirst.toLocaleString() }}</strong> 分！
            <br />
            <span class="suggestion-hint">继续鼓励学生举手发言和完成作业！</span>
          </div>
          <div v-else class="suggestion victory">
            🎉 太棒了！本班目前位列年级第 {{ overview.rank }}！
            <br />
            <span class="suggestion-hint">继续保持领先优势！</span>
          </div>
          <button
            class="challenge-btn"
            @click="alert('🚀 已发起挑战！本周内总积分超过对方即可获胜！')"
          >
            ⚡ 发起挑战
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.pk-page {
  max-width: 1100px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 24px;
}
.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}
.page-title { font-size: 26px; font-weight: 700; margin: 0; }
.page-subtitle { font-size: 14px; color: var(--color-text-secondary); }
.page-link-btn { padding:4px 12px; border-radius:8px; border:1px solid var(--color-border); background:transparent; color:var(--color-accent); font-size:12px; font-weight:500; cursor:pointer; white-space:nowrap; transition:0.15s; font-family:inherit; }
.page-link-btn:hover { background:rgba(79,70,229,0.08); border-color:var(--color-accent); }
.rank-badge {
  padding: 6px 20px;
  border-radius: 30px;
  background: rgba(245,158,11,0.08);
  border: 1px solid rgba(245,158,11,0.15);
  color: #F59E0B;
  font-size: 15px;
  font-weight: 700;
}

/* PK 排行榜 */
.pk-leaderboard {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 24px;
}
.pk-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 20px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  transition: all 0.25s ease;
}
.pk-item:hover {
  box-shadow: var(--shadow-sm);
}
.pk-item--own {
  border-color: rgba(245,158,11,0.2);
  background: linear-gradient(135deg, rgba(245,158,11,0.03), transparent);
}

.pk-rank {
  font-size: 22px;
  font-weight: 800;
  width: 44px;
  text-align: center;
}
.rank-gold { color: #F59E0B; }
.rank-silver { color: #94A3B8; }
.rank-bronze { color: #D97706; }

.pk-class {
  display: flex;
  flex-direction: column;
  width: 130px;
  flex-shrink: 0;
}
.pk-class-name {
  font-size: 16px;
  font-weight: 600;
}
.pk-own-badge {
  font-size: 10px;
  padding: 1px 8px;
  border-radius: 4px;
  background: rgba(245,158,11,0.1);
  color: #F59E0B;
  font-weight: 600;
  width: fit-content;
}

.pk-bar-area {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}
.pk-bar-track {
  flex: 1;
  height: 8px;
  background: var(--color-border);
  border-radius: 4px;
  overflow: hidden;
}
.pk-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary), #818CF8);
  border-radius: 4px;
  transition: width 0.8s ease;
}
.fill-gold {
  background: linear-gradient(90deg, #F59E0B, #FCD34D);
}
.pk-bar-score {
  font-weight: 700;
  font-size: 15px;
  min-width: 70px;
  text-align: right;
}

.pk-stats {
  display: flex;
  gap: 14px;
  font-size: 13px;
  color: var(--color-text-secondary);
  min-width: 160px;
  justify-content: flex-end;
}

.pk-btn {
  padding: 6px 16px;
  border-radius: 20px;
  border: 1px solid rgba(167,139,250,0.2);
  background: rgba(167,139,250,0.08);
  color: var(--color-primary-light, #a78bfa);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
  font-family: inherit;
}
.pk-btn:hover {
  background: rgba(167,139,250,0.15);
  transform: scale(1.03);
}
.pk-btn.own {
  background: rgba(245,158,11,0.08);
  border-color: rgba(245,158,11,0.2);
  color: #F59E0B;
  cursor: default;
}
.pk-btn.own:hover {
  transform: none;
}

/* 详情双栏 */
.pk-detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}
.pk-detail-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 18px;
  padding: 22px 24px;
}
.pk-detail-card h4 {
  font-size: 14px;
  color: var(--color-text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 16px;
}
.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid var(--color-border);
  font-size: 15px;
}
.detail-row:last-child {
  border-bottom: none;
}
.row-val { font-weight: 700; }
.row-val.peak { color: #8B5CF6; }
.row-val.weekly { color: #10B981; }

.suggestion {
  font-size: 15px;
  color: var(--color-text-secondary);
  line-height: 1.6;
}
.gap-score {
  color: #F59E0B;
  font-size: 18px;
}
.suggestion-hint {
  font-size: 13px;
  color: var(--color-text-secondary);
  opacity: 0.7;
}
.victory {
  color: #10B981;
  font-weight: 600;
}

.challenge-btn {
  margin-top: 16px;
  width: 100%;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid rgba(167,139,250,0.2);
  background: rgba(167,139,250,0.08);
  color: var(--color-primary-light, #a78bfa);
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
}
.challenge-btn:hover {
  background: rgba(167,139,250,0.18);
  transform: translateY(-2px);
}

/* 加载 */
.loading-state {
  text-align: center;
  padding: 60px;
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

@media (max-width: 768px) {
  .pk-item {
    flex-wrap: wrap;
    gap: 8px;
  }
  .pk-stats {
    min-width: unset;
    width: 100%;
    justify-content: flex-start;
  }
  .pk-detail-grid {
    grid-template-columns: 1fr;
  }
  .pk-class {
    width: auto;
  }
}
</style>
