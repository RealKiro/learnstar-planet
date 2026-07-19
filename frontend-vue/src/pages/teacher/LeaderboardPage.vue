<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { avatarGradient } from '@/utils/constants'
import { getSpeciesEmoji } from '@/utils/petData'
import type { ApiResponse, LeaderboardEntry } from '@/types'

type LbType = 'total' | 'weekly' | 'pet'

const activeTab = ref<LbType>('total')
const entries = ref<LeaderboardEntry[]>([])
const loading = ref(false)

const tabs: Array<{ key: LbType; label: string; icon: string }> = [
  { key: 'total', label: '总积分榜', icon: '🏆' },
  { key: 'weekly', label: '进步最快榜', icon: '📈' },
  { key: 'pet', label: '宠物等级榜', icon: '🌟' },
]

const rankMedals = (i: number) => ['🥇', '🥈', '🥉'][i] || (i + 1).toString()

// 前三名
const topThree = computed(() => entries.value.slice(0, 3))
const restEntries = computed(() => entries.value.slice(3))

// 进步最大
const maxProgress = computed(() => {
  if (activeTab.value !== 'weekly' || entries.value.length === 0) return null
  const e = entries.value[0]
  return { name: e.student_name, points: e.score, change: '+120分' }
})

async function fetchData(type: LbType) {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<LeaderboardEntry[]>>(`/api/v1/teacher/leaderboard/${type}`)
    entries.value = (res.data || []).map((e, i) => ({
      ...e,
      pet_emoji: getSpeciesEmoji((e as any).pet_species || ''),
    }))
  } catch {
    // 演示数据
    entries.value = generateDemoData(type)
  } finally {
    loading.value = false
  }
}

function generateDemoData(type: LbType): LeaderboardEntry[] {
  const names = ['张小明', '李小红', '王小刚', '赵小丽', '刘小强', '陈小美', '周小龙', '吴小凤']
  const baseScore = type === 'weekly' ? 50 : type === 'pet' ? 12 : 500
  const species = ['zhulong', 'nine_tail_fox', 'charmander', 'pikachu', 'panda', 'unicorn', 't_rex', 'nian']
  return names.map((name, i) => ({
    rank: i + 1,
    student_id: i + 1,
    student_name: name,
    score: type === 'pet' ? Math.max(1, baseScore - i * 1.5).toString() : String(Math.max(1, baseScore - i * (type === 'weekly' ? 5 : 40))),
    pet_name: `${name}的伙伴`,
    pet_level: Math.max(1, Math.floor(baseScore - i * 1.2)),
    pet_species: species[i],
    pet_emoji: getSpeciesEmoji(species[i]),
  }))
}

function formatScore(score: number | string): string {
  if (typeof score === 'string') return score
  return score.toLocaleString()
}

async function switchTab(type: LbType) {
  activeTab.value = type
  await fetchData(type)
}

onMounted(() => fetchData('total'))
</script>

<template>
  <div class="leaderboard-page">
    <div class="page-header">
      <h2 class="page-title">🏆 排行榜</h2>
      <span class="page-subtitle">每周一重置 · 保持努力</span>
    </div>

    <!-- Tab 切换 -->
    <div class="tab-bar">
      <button
        v-for="t in tabs" :key="t.key"
        class="tab-btn"
        :class="{ 'tab--active': activeTab === t.key }"
        @click="switchTab(t.key)"
      >
        <span class="tab-icon">{{ t.icon }}</span>
        <span>{{ t.label }}</span>
      </button>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载排行中...</p>
    </div>

    <template v-else-if="entries.length === 0">
      <div class="empty-state">
        <div class="empty-icon">📭</div>
        <p>暂无排行数据</p>
      </div>
    </template>

    <template v-else>
      <!-- 冠军专区 -->
      <div class="champion-section">
        <div
          v-for="(e, i) in topThree"
          :key="e.student_id"
          class="champion-card"
          :class="'rank--' + (i + 1)"
        >
          <div class="champion-badge">
            <span class="medal">{{ rankMedals(i) }}</span>
          </div>

          <!-- 头像 / 宠物 -->
          <div class="champion-avatar" :style="{
            background: i === 0 ? 'linear-gradient(135deg, #F59E0B, #FCD34D)' :
              i === 1 ? 'linear-gradient(135deg, #94A3B8, #CBD5E1)' :
              'linear-gradient(135deg, #D97706, #F59E0B)',
          }">
            <span class="avatar-emoji">{{ e.pet_emoji || '🌟' }}</span>
          </div>

          <div class="champion-name">{{ e.student_name }}</div>

          <div v-if="e.pet_name" class="champion-pet">{{ e.pet_name }}</div>

          <!-- 进度条 -->
          <div class="champion-bar">
            <div class="bar-label">
              <span>{{ activeTab === 'pet' ? 'Lv.' : '' }}{{ formatScore(e.score) }}</span>
              <span v-if="activeTab !== 'pet' && (e as any).pet_level" class="bar-sub">Lv.{{ (e as any).pet_level }}</span>
            </div>
            <div class="bar-track">
              <div
                class="bar-fill"
                :style="{ width: Math.min(100, (parseFloat(String(e.score)) / parseFloat(String(topThree[0]?.score || 1))) * 100) + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 剩余排名列表 -->
      <div class="rank-list">
        <div
          v-for="(e, i) in restEntries"
          :key="e.student_id"
          class="rank-item"
        >
          <span class="rank-num">{{ e.rank }}</span>
          <div class="rank-avatar" :style="{ background: avatarGradient(e.student_name) }">
            {{ e.student_name.charAt(0) }}
          </div>
          <div class="rank-info">
            <span class="rank-name">{{ e.student_name }}</span>
            <span v-if="e.pet_name" class="rank-pet">{{ e.pet_name }}</span>
          </div>
          <div class="rank-score">
            <span class="score-num">{{ formatScore(e.score) }}</span>
            <span v-if="(e as any).pet_level" class="score-sub">Lv.{{ (e as any).pet_level }}</span>
          </div>
          <div v-if="i < 3" class="rank-change up">⬆️ {{ i + 1 }}</div>
          <div v-else-if="i < 6" class="rank-change down">⬇️ {{ i - 1 }}</div>
        </div>
      </div>

      <!-- 周进步奖 -->
      <div v-if="activeTab === 'weekly' && maxProgress" class="progress-award">
        <div class="award-icon">🚀</div>
        <div class="award-info">
          <span class="award-label">本周进步最大</span>
          <span class="award-name">{{ maxProgress.name }}</span>
          <span class="award-change">{{ maxProgress.change }}</span>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.leaderboard-page {
  max-width: 800px;
}

.page-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 20px;
}
.page-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0;
}
.page-subtitle {
  font-size: 13px;
  color: var(--color-text-secondary);
}

/* Tab */
.tab-bar {
  display: flex;
  gap: 4px;
  background: var(--color-bg);
  border-radius: 12px;
  padding: 4px;
  margin-bottom: 24px;
}
.tab-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  background: transparent;
  color: var(--color-text-secondary);
  transition: all 0.2s ease;
}
.tab-btn:hover {
  color: var(--color-text);
}
.tab--active {
  background: var(--color-bg-card);
  color: var(--color-text);
  box-shadow: var(--shadow-sm);
}
.tab-icon { font-size: 16px; }

/* 冠军专区 */
.champion-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.champion-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  padding: 20px 16px;
  text-align: center;
  position: relative;
  transition: all 0.3s ease;
}
.champion-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}
.rank--1 {
  border-color: rgba(245,158,11,0.2);
  background: linear-gradient(180deg, rgba(245,158,11,0.05), var(--color-bg-card));
  box-shadow: 0 4px 20px rgba(245,158,11,0.08);
}
.rank--2 {
  border-color: rgba(148,163,184,0.15);
}
.rank--3 {
  border-color: rgba(217,119,6,0.15);
}

.champion-badge {
  margin-bottom: 8px;
}
.medal {
  font-size: 28px;
  line-height: 1;
}

.champion-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 8px;
}
.avatar-emoji {
  font-size: 28px;
}

.champion-name {
  font-size: 15px;
  font-weight: 700;
  margin-bottom: 2px;
}
.champion-pet {
  font-size: 11px;
  color: var(--color-text-secondary);
  margin-bottom: 10px;
}

.champion-bar {
  width: 100%;
}
.bar-label {
  display: flex;
  justify-content: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 700;
  color: var(--color-primary);
  margin-bottom: 4px;
}
.bar-sub {
  font-size: 11px;
  font-weight: 500;
  color: var(--color-text-secondary);
}
.bar-track {
  height: 6px;
  background: var(--color-border);
  border-radius: 3px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  background: var(--gradient-primary);
  border-radius: 3px;
  transition: width 0.5s ease;
}

/* 排名列表 */
.rank-list {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 16px;
}

.rank-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  transition: background 0.2s ease;
}
.rank-item:hover {
  background: rgba(79,70,229,0.03);
}
.rank-item + .rank-item {
  border-top: 1px solid var(--color-border);
}

.rank-num {
  width: 24px;
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-align: center;
}

.rank-avatar {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.rank-info {
  flex: 1;
  min-width: 0;
}
.rank-name {
  font-size: 14px;
  font-weight: 500;
  display: block;
}
.rank-pet {
  font-size: 11px;
  color: var(--color-text-secondary);
}

.rank-score {
  text-align: right;
  flex-shrink: 0;
}
.score-num {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-primary);
  display: block;
}
.score-sub {
  font-size: 10px;
  color: var(--color-text-secondary);
}

.rank-change {
  font-size: 11px;
  font-weight: 600;
  width: 40px;
  text-align: right;
}
.up { color: #10B981; }
.down { color: #EF4444; }

/* 进步奖 */
.progress-award {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(16,185,129,0.08), rgba(16,185,129,0.03));
  border: 1px solid rgba(16,185,129,0.15);
  border-radius: 12px;
}
.award-icon {
  font-size: 28px;
}
.award-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.award-label {
  font-size: 12px;
  color: var(--color-text-secondary);
}
.award-name {
  font-size: 14px;
  font-weight: 700;
}
.award-change {
  font-size: 13px;
  font-weight: 600;
  color: #10B981;
}

/* 加载/空 */
.loading-state, .empty-state {
  text-align: center;
  padding: 48px;
  color: var(--color-text-secondary);
}
.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }
.empty-icon { font-size: 48px; margin-bottom: 8px; }

@media (max-width: 600px) {
  .champion-section {
    grid-template-columns: 1fr;
  }
}
</style>
