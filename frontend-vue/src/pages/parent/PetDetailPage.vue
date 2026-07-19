<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import {
  getSeriesBySpeciesId,
  getSpeciesById,
  getLevelRequiredScore,
  SERIES_SCENES,
  getSpeciesEmoji,
  getPetLevelName,
} from '@/utils/petData'
import type { ApiResponse, PetDetail } from '@/types'
import PetDisplay from '@/components/pet/PetDisplay.vue'
import PetEvolutionTree from '@/components/pet/PetEvolutionTree.vue'
import PetFeedAnimation from '@/components/pet/PetFeedAnimation.vue'
import PetHandbook from '@/components/pet/PetHandbook.vue'

const route = useRoute()
const toast = useToastStore()

const pet = ref<PetDetail | null>(null)
const loading = ref(true)
const showFeedAnimation = ref(false)
const feedPoints = ref(0)
const showHandbook = ref(false)
const activeTab = ref<'info' | 'evolution'>('info')
const currentQuote = ref('')

const petQuotes = [
  '我感觉身体里有一股神秘力量！✨',
  '你今天太棒了，我快吃饱啦！😋',
  '主人，我饿啦，快去做任务吧！🥺',
  '再得几分我就能进化啦！💪',
  '今天天气真好，一起学习吧！📚',
  '你认真做题的样子真帅！🌟',
]

const scene = computed(() => {
  if (!pet.value) return null
  const series = getSeriesBySpeciesId(pet.value.species)
  return series ? SERIES_SCENES[series.id] : null
})

const nextLevelScore = computed(() => {
  if (!pet.value) return 0
  return getLevelRequiredScore(Math.min(pet.value.level + 1, 12))
})

const progressPercent = computed(() => {
  if (!pet.value) return 0
  const current = getLevelRequiredScore(pet.value.level)
  const next = getLevelRequiredScore(Math.min(pet.value.level + 1, 12))
  return Math.min(99, ((pet.value.exp) / (next - current)) * 100)
})

async function loadPet() {
  loading.value = true
  try {
    const studentId = route.query.student_id || 0
    const res = await apiGet<ApiResponse<PetDetail>>(`/api/v1/parent/pet?student_id=${studentId}`)
    pet.value = res.data
  } catch {
    // Demo data
    const demoSpecies = ['zhulong', 'nine_tail_fox', 'charmander', 'pikachu', 'panda']
    const speciesId = demoSpecies[Math.floor(Math.random() * demoSpecies.length)]
    const level = Math.min(12, Math.max(1, Math.floor(Math.random() * 10) + 2))
    const species = getSpeciesById(speciesId)
    const series = getSeriesBySpeciesId(speciesId)
    const levelData = species?.levels.find(l => l.level === level)
    pet.value = {
      id: 1,
      student_id: 1,
      student_name: '我的孩子',
      name: species?.name || '宠物',
      species: speciesId,
      level,
      exp: Math.floor(Math.random() * 100),
      mood: 75,
      seriesId: series?.id,
      seriesName: series?.name,
      levelName: levelData?.name,
      levelDescription: levelData?.description,
      requiredScore: levelData?.requiredScore,
      streakDays: 5,
    }
  } finally {
    loading.value = false
  }
}

function randomQuote() {
  currentQuote.value = petQuotes[Math.floor(Math.random() * petQuotes.length)]
}

async function handleFeed() {
  if (!pet.value) return
  feedPoints.value = Math.floor(Math.random() * 15) + 5
  showFeedAnimation.value = true
  if (pet.value) pet.value.exp += feedPoints.value
  randomQuote()
}

async function handleAddScore() {
  if (!pet.value) return
  feedPoints.value = 10
  showFeedAnimation.value = true
  if (pet.value) pet.value.exp += 10
}

function openHandbook() {
  showHandbook.value = true
}

const moodText = computed(() => {
  if (!pet.value?.mood) return '😊 开心'
  if (pet.value.mood >= 80) return '🥰 非常开心'
  if (pet.value.mood >= 60) return '😊 开心'
  if (pet.value.mood >= 40) return '🙂 不错'
  if (pet.value.mood >= 20) return '😐 一般'
  return '😢 需要关注'
})

onMounted(() => {
  loadPet()
  randomQuote()
  setInterval(randomQuote, 12000)
})
</script>

<template>
  <div class="pet-detail-page">
    <!-- 加载 -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>召唤宠物中...</p>
    </div>

    <template v-else-if="pet">
      <!-- 顶部概览 -->
      <div class="overview-bar">
        <div class="overview-left">
          <div class="student-info">
            <div class="student-avatar">{{ pet.student_name?.charAt(0) }}</div>
            <div>
              <div class="student-name">{{ pet.student_name }}</div>
              <div class="pet-species">{{ pet.seriesName }} · {{ getPetLevelName(pet.species, pet.level) }}</div>
            </div>
          </div>
        </div>
        <div class="overview-right">
          <div class="streak-badge" :class="{ 'streak--hot': (pet.streakDays || 0) >= 5 }">
            🔥 {{ pet.streakDays || 0 }} 天连击
          </div>
        </div>
      </div>

      <!-- 主内容：沉浸式树屋布局 -->
      <div class="main-layout">
        <!-- 左：宠物展示区 (65%) -->
        <div class="pet-stage-area">
          <PetDisplay
            :species-id="pet.species"
            :level="pet.level"
            :name="pet.name"
            :mood="pet.mood"
            :animate="true"
            size="xl"
            @feed="handleFeed"
          />

          <!-- 积分和进度 -->
          <div class="score-section">
            <div class="score-display">
              <div class="score-title">当前积分</div>
              <div class="score-number">{{ pet.exp }}</div>
              <div class="score-next">下一级还需 <strong>{{ Math.max(0, nextLevelScore - pet.exp) }}</strong> 分</div>
            </div>

            <!-- 进化光之轨迹 -->
            <div class="evolution-rail">
              <div class="rail-header">
                <span class="rail-title">🌟 进化之路</span>
                <span class="rail-level">Lv.{{ pet.level }} / 12</span>
              </div>
              <div class="rail-track">
                <div class="track-bg"></div>
                <div class="track-fill" :style="{ width: ((pet.level - 1) / 11) * 100 + '%' }"></div>
                <div
                  v-for="lvl in 12"
                  :key="lvl"
                  class="track-node"
                  :class="{
                    'node--unlocked': lvl <= pet.level,
                    'node--current': lvl === pet.level,
                  }"
                  :style="{ left: ((lvl - 1) / 11) * 100 + '%' }"
                >
                  <div class="node-dot">
                    {{ lvl <= pet.level ? '⭐' : '🔘' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右：控制面板 (35%) -->
        <div class="control-panel">
          <!-- 宠物语录 -->
          <div class="quote-card glass">
            <div class="quote-icon">💬</div>
            <div class="quote-text">"{{ currentQuote }}"</div>
          </div>

          <!-- 心情状态 -->
          <div class="mood-card glass">
            <div class="mood-row">
              <span>😊 心情</span>
              <span class="mood-value">{{ moodText }}</span>
            </div>
            <div class="mood-row">
              <span>📊 等级</span>
              <span class="mood-value">Lv.{{ pet.level }} · {{ getPetLevelName(pet.species, pet.level) }}</span>
            </div>
            <div class="mood-row">
              <span>🏆 系列</span>
              <span class="mood-value">{{ pet.seriesName }}</span>
            </div>
          </div>

          <!-- 三大金刚按钮 -->
          <div class="action-group">
            <button class="action-btn btn-primary-glow" @click="handleFeed">
              <span class="action-icon">🍬</span>
              <span class="action-label">投喂</span>
              <span class="action-desc">消耗积分喂养</span>
            </button>
            <button class="action-btn btn-secondary" @click="handleAddScore">
              <span class="action-icon">⭐</span>
              <span class="action-label">加分</span>
              <span class="action-desc">快速增加经验</span>
            </button>
            <button class="action-btn btn-tertiary" @click="openHandbook">
              <span class="action-icon">📖</span>
              <span class="action-label">图鉴</span>
              <span class="action-desc">查看进化预览</span>
            </button>
          </div>

          <!-- 今日任务 -->
          <div class="task-card glass">
            <div class="task-header">
              <span>📋 今日任务</span>
              <span class="task-count">完成 3 项得勤奋礼包</span>
            </div>
            <div class="task-list">
              <label class="task-item">
                <input type="checkbox" class="task-check" />
                <span>📖 完成今日作业 +3分</span>
              </label>
              <label class="task-item">
                <input type="checkbox" class="task-check" />
                <span>🤝 帮助一位同学 +2分</span>
              </label>
              <label class="task-item">
                <input type="checkbox" class="task-check" />
                <span>💪 主动回答问题 +5分</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部：班级动态 -->
      <div class="classroom-feed glass">
        <span class="feed-label">💬 班级动态</span>
        <span class="feed-text">李四的【九尾狐】刚刚进化成了【六尾妖狐】！🎉</span>
      </div>
    </template>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">🥚</div>
      <p>暂无宠物数据</p>
      <p class="empty-hint">请联系老师领取宠物蛋</p>
    </div>

    <!-- 投喂动画 -->
    <PetFeedAnimation
      :show="showFeedAnimation"
      :points="feedPoints"
      :species-id="pet?.species || 'zhulong'"
      :level="pet?.level || 1"
      @close="showFeedAnimation = false"
    />

    <!-- 图鉴 -->
    <PetHandbook
      v-if="showHandbook && pet"
      :species-id="pet.species"
      :current-level="pet.level"
      :current-score="pet.exp"
      @close="showHandbook = false"
    />
  </div>
</template>

<style scoped>
.pet-detail-page {
  min-height: 100%;
}

/* 加载/空 */
.loading-state, .empty-state {
  text-align: center;
  padding: 80px 24px;
  color: var(--color-text-secondary);
}
.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }
.empty-icon { font-size: 64px; margin-bottom: 12px; }
.empty-hint { font-size: 13px; opacity: 0.6; margin-top: 8px; }

/* 顶部概览 */
.overview-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}
.student-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.student-avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: var(--gradient-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
}
.student-name {
  font-size: 18px;
  font-weight: 700;
}
.pet-species {
  font-size: 13px;
  color: var(--color-text-secondary);
}
.streak-badge {
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  background: rgba(245,158,11,0.1);
  color: #F59E0B;
  border: 1px solid rgba(245,158,11,0.2);
}
.streak--hot {
  background: rgba(245,158,11,0.2);
  color: #F59E0B;
  animation: streakGlow 2s ease-in-out infinite;
}
@keyframes streakGlow {
  0%, 100% { box-shadow: 0 0 8px rgba(245,158,11,0.2); }
  50% { box-shadow: 0 0 16px rgba(245,158,11,0.4); }
}

/* 主布局 */
.main-layout {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

/* 左：宠物舞台 */
.pet-stage-area {
  flex: 1;
  min-width: 0;
}

.score-section {
  margin-top: 16px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  padding: 20px;
}
.score-display {
  text-align: center;
  margin-bottom: 16px;
}
.score-title {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}
.score-number {
  font-size: 48px;
  font-weight: 900;
  background: linear-gradient(135deg, #4F46E5, #818CF8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
  margin-bottom: 4px;
}
.score-next {
  font-size: 13px;
  color: var(--color-text-secondary);
}
.score-next strong {
  color: var(--color-primary);
  font-size: 16px;
}

/* 进化光轨 */
.evolution-rail {
  padding: 8px 0;
}
.rail-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}
.rail-title {
  font-size: 14px;
  font-weight: 600;
}
.rail-level {
  font-size: 12px;
  color: var(--color-text-secondary);
}
.rail-track {
  position: relative;
  height: 32px;
  display: flex;
  align-items: center;
}
.track-bg {
  position: absolute;
  left: 0;
  right: 0;
  height: 4px;
  background: var(--color-border);
  border-radius: 2px;
}
.track-fill {
  position: absolute;
  left: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--color-primary), #F59E0B, #EF4444);
  border-radius: 2px;
  transition: width 0.6s ease;
  max-width: 100%;
}
.track-node {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
}
.node-dot {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  background: var(--color-bg-card);
  border: 2px solid var(--color-border);
  transition: all 0.3s ease;
}
.node--unlocked .node-dot {
  border-color: var(--color-primary);
  background: var(--color-bg);
}
.node--current .node-dot {
  border-color: #F59E0B;
  box-shadow: 0 0 8px rgba(245,158,11,0.4);
}

/* 右：控制面板 */
.control-panel {
  width: 340px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.glass {
  background: rgba(255,255,255,0.6);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255,255,255,0.8);
  border-radius: 16px;
  padding: 16px;
}
.dark .glass {
  background: rgba(30,41,59,0.6);
  border-color: rgba(255,255,255,0.06);
}

/* 语录 */
.quote-card {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}
.quote-icon {
  font-size: 20px;
  flex-shrink: 0;
}
.quote-text {
  font-size: 14px;
  color: var(--color-text-secondary);
  font-style: italic;
  line-height: 1.5;
}

/* 心情 */
.mood-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.mood-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}
.mood-row + .mood-row {
  padding-top: 8px;
  border-top: 1px solid var(--color-border);
}
.mood-value {
  font-weight: 600;
}

/* 三大按钮 */
.action-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.action-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 14px;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}
.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
.action-icon {
  font-size: 28px;
  flex-shrink: 0;
}
.action-label {
  font-size: 15px;
  font-weight: 600;
  display: block;
}
.action-desc {
  font-size: 11px;
  color: var(--color-text-secondary);
  display: block;
}
.btn-primary-glow {
  border-color: rgba(245,158,11,0.2);
  background: linear-gradient(135deg, rgba(245,158,11,0.05), rgba(245,158,11,0.02));
}
.btn-primary-glow:hover {
  border-color: #F59E0B;
  box-shadow: 0 4px 16px rgba(245,158,11,0.2);
}
.btn-secondary {
  border-color: rgba(16,185,129,0.15);
}
.btn-secondary:hover {
  border-color: #10B981;
}
.btn-tertiary {
  border-color: rgba(99,102,241,0.15);
}
.btn-tertiary:hover {
  border-color: #6366F1;
}

/* 任务 */
.task-card {}
.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 13px;
  font-weight: 600;
}
.task-count {
  font-size: 11px;
  font-weight: 400;
  color: var(--color-text-secondary);
}
.task-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.task-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  cursor: pointer;
  padding: 6px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}
.task-item:hover {
  background: rgba(79,70,229,0.04);
}
.task-check {
  width: 16px;
  height: 16px;
  accent-color: var(--color-primary);
}

/* 底部动态 */
.classroom-feed {
  margin-top: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  font-size: 13px;
}
.feed-label {
  font-weight: 600;
  white-space: nowrap;
}
.feed-text {
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 响应式 */
@media (max-width: 860px) {
  .main-layout {
    flex-direction: column;
  }
  .control-panel {
    width: 100%;
  }
}
</style>
