<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import {
  getAllSeries,
  getSeriesBySpeciesId,
  getSpeciesById,
  getSeriesName,
  getSpeciesEmoji,
  SERIES_SCENES,
  getPetLevelName,
  getPetLevelDescription,
  getLevelStage,
} from '@/utils/petData'
import type { ApiResponse, Pet, PetDetail } from '@/types'
import PetDisplay from '@/components/pet/PetDisplay.vue'
import PetEvolutionTree from '@/components/pet/PetEvolutionTree.vue'
import PetFeedAnimation from '@/components/pet/PetFeedAnimation.vue'
import PetHandbook from '@/components/pet/PetHandbook.vue'

const toast = useToastStore()

// ===== 数据状态 =====
const pets = ref<PetDetail[]>([])
const loading = ref(true)
const selectedPet = ref<PetDetail | null>(null)
const showFeedAnimation = ref(false)
const feedPoints = ref(0)
const showHandbook = ref(false)
const showDetailPanel = ref(false)
const activeTab = ref<'info' | 'evolution' | 'tasks'>('info')
const streakDays = ref(0)

// 系列筛选
const activeSeries = ref<string | 'all'>('all')

// 统计数据
const stats = computed(() => ({
  total: pets.value.length,
  avgLevel: pets.value.length
    ? Math.round(pets.value.reduce((s, p) => s + p.level, 0) / pets.value.length)
    : 0,
  legendary: pets.value.filter(p => p.level >= 11).length,
  mature: pets.value.filter(p => p.level >= 8 && p.level < 11).length,
}))

// 按系列分组的宠物
const groupedPets = computed(() => {
  const groups: Record<string, PetDetail[]> = {}
  for (const p of pets.value) {
    const seriesId = p.seriesId || 'myth'
    if (!groups[seriesId]) groups[seriesId] = []
    groups[seriesId].push(p)
  }
  return groups
})

// 筛选后的宠物
const filteredPets = computed(() => {
  if (activeSeries.value === 'all') return pets.value
  return pets.value.filter(p => (p.seriesId || 'myth') === activeSeries.value)
})

// 按等级排序
const sortedPets = computed(() => {
  return [...filteredPets.value].sort((a, b) => b.level - a.level || b.exp - a.exp)
})

// ===== 方法 =====

async function loadPets() {
  loading.value = true
  try {
    const res = await apiGet<ApiResponse<Pet[]>>('/api/v1/teacher/pets/overview')
    const rawPets = res.data || []
    // 扩展宠物数据
    pets.value = rawPets.map(p => {
      const series = getSeriesBySpeciesId(p.species)
      const species = getSpeciesById(p.species)
      const levelData = species?.levels.find(l => l.level === p.level)
      return {
        ...p,
        seriesId: series?.id || 'myth',
        seriesName: series?.name || '未知',
        levelName: levelData?.name || getPetLevelName(p.species, p.level),
        levelDescription: levelData?.description || '',
        stage: levelData?.stage || getLevelStage(p.level),
        requiredScore: levelData?.requiredScore || 0,
        streakDays: streakDays.value,
      }
    })
  } catch {
    // 使用演示数据
    pets.value = generateDemoPets()
  } finally {
    loading.value = false
  }
}

// 演示数据
function generateDemoPets(): PetDetail[] {
  const demoSpecies = ['zhulong', 'nine_tail_fox', 'charmander', 'pikachu', 'panda', 'cyber_cat', 'unicorn', 't_rex', 'nian']
  const demoStudents = ['张小明', '李小红', '王小刚', '赵小丽', '刘小强', '陈小美', '周小龙', '吴小凤', '郑小天', '孙小艺']
  return demoStudents.map((name, i) => {
    const speciesId = demoSpecies[i % demoSpecies.length]
    const level = Math.min(12, Math.max(1, Math.floor(Math.random() * 12) + 1))
    const species = getSpeciesById(speciesId)
    const series = getSeriesBySpeciesId(speciesId)
    const levelData = species?.levels.find(l => l.level === level)
    return {
      id: i + 1,
      student_id: i + 1,
      student_name: name,
      name: `${name}的${species?.name || '宠物'}`,
      species: speciesId,
      level,
      exp: Math.floor(Math.random() * 100),
      mood: Math.floor(Math.random() * 100),
      seriesId: series?.id || 'myth',
      seriesName: series?.name || '未知',
      levelName: levelData?.name || '',
      levelDescription: levelData?.description || '',
      stage: levelData?.stage || 'baby',
      requiredScore: levelData?.requiredScore || 0,
      streakDays: Math.floor(Math.random() * 14),
    }
  })
}

function selectPet(pet: PetDetail) {
  selectedPet.value = pet
  showDetailPanel.value = true
  activeTab.value = 'info'
}

function closeDetail() {
  showDetailPanel.value = false
  selectedPet.value = null
}

// 投喂
async function handleFeed() {
  if (!selectedPet.value) return
  try {
    await apiPost(`/api/v1/teacher/pets/feed/${selectedPet.value.id}`)
    feedPoints.value = Math.floor(Math.random() * 15) + 5
    showFeedAnimation.value = true
    // 刷新数据
    await loadPets()
    // 更新选中宠物
    const updated = pets.value.find(p => p.id === selectedPet.value?.id)
    if (updated) selectedPet.value = updated
  } catch {
    // 离线模式演示
    feedPoints.value = Math.floor(Math.random() * 15) + 5
    showFeedAnimation.value = true
  }
}

function handleAddScore() {
  if (!selectedPet.value) return
  feedPoints.value = 10
  showFeedAnimation.value = true
}

function openHandbook() {
  showHandbook.value = true
}

// 随机宠物语录
const petQuotes = [
  '我感觉身体里有一股神秘力量！✨',
  '你今天太棒了，我快吃饱啦！😋',
  '主人，我饿啦，快去做任务吧！🥺',
  '再得几分我就能进化啦！💪',
  '今天天气真好，一起学习吧！📚',
  '你认真做题的样子真帅！🌟',
]

const currentQuote = ref('')
function randomQuote() {
  currentQuote.value = petQuotes[Math.floor(Math.random() * petQuotes.length)]
}

onMounted(async () => {
  await loadPets()
  randomQuote()
  setInterval(randomQuote, 15000)
})
</script>

<template>
  <div class="pets-page">
    <!-- 顶部状态条 -->
    <div class="top-bar">
      <div class="top-bar-left">
        <h2 class="page-title">🌳 宠物花园</h2>
        <span class="page-subtitle">全班 {{ stats.total }} 只宠物在成长</span>
      </div>
      <div class="top-bar-right">
        <div class="stat-pills">
          <div class="stat-pill pill--avg">
            <span class="pill-icon">📊</span>
            <span>平均 Lv.{{ stats.avgLevel }}</span>
          </div>
          <div class="stat-pill pill--mature">
            <span class="pill-icon">🌟</span>
            <span>{{ stats.mature }} 成熟</span>
          </div>
          <div class="stat-pill pill--legendary">
            <span class="pill-icon">👑</span>
            <span>{{ stats.legendary }} 传说</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 系列筛选 -->
    <div class="series-filter">
      <button
        class="filter-chip"
        :class="{ 'chip--active': activeSeries === 'all' }"
        @click="activeSeries = 'all'"
      >
        🌟 全部
      </button>
      <button
        v-for="series in getAllSeries()"
        :key="series.id"
        class="filter-chip"
        :class="{ 'chip--active': activeSeries === series.id }"
        :style="activeSeries === series.id ? {
          background: SERIES_SCENES[series.id]?.primaryColor + '22',
          borderColor: SERIES_SCENES[series.id]?.primaryColor + '66',
        } : {}"
        @click="activeSeries = series.id"
      >
        {{ series.emoji }} {{ series.name }}
      </button>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧：宠物网格 -->
      <div class="pet-grid-area">
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <p>正在召唤宠物...</p>
        </div>

        <div v-else-if="sortedPets.length === 0" class="empty-state">
          <div class="empty-icon">🌱</div>
          <p>暂无宠物</p>
          <p class="empty-hint">在积分规则中开启宠物系统即可创建</p>
        </div>

        <div v-else class="pet-grid">
          <div
            v-for="pet in sortedPets"
            :key="pet.id"
            class="pet-card"
            :class="{
              'card--selected': selectedPet?.id === pet.id,
              'card--legendary': pet.level >= 11,
              'card--mature': pet.level >= 8 && pet.level < 11,
            }"
            :style="{
              '--card-accent': SERIES_SCENES[pet.seriesId || 'myth']?.primaryColor || '#6366F1',
            }"
            @click="selectPet(pet)"
          >
            <!-- 卡片背景 -->
            <div
              class="card-bg"
              :style="{ background: SERIES_SCENES[pet.seriesId || 'myth']?.bgGradient }"
            ></div>

            <!-- 宠物 Emoji -->
            <div class="card-emoji">
              {{ getSpeciesEmoji(pet.species) }}
            </div>

            <!-- 等级徽章 -->
            <div class="card-level" :class="'stage--' + (pet.stage || 'baby')">
              Lv.{{ pet.level }}
            </div>

            <!-- 信息 -->
            <div class="card-info">
              <div class="card-name">{{ pet.student_name }}</div>
              <div class="card-pet-name">{{ pet.levelName || pet.name }}</div>
            </div>

            <!-- 经验条 -->
            <div class="card-exp">
              <div class="exp-track">
                <div class="exp-fill" :style="{ width: Math.min(100, pet.exp % 100) + '%' }"></div>
              </div>
            </div>

            <!-- 传说级特效 -->
            <div v-if="pet.level >= 11" class="legendary-glow"></div>
          </div>
        </div>
      </div>

      <!-- 右侧：宠物详情面板 -->
      <Transition name="panel-slide">
        <div v-if="showDetailPanel && selectedPet" class="detail-panel">
          <!-- 选择其他标签的临时占位 -->
          <div v-if="activeTab === 'info'" class="detail-tab-content">
            <!-- 宠物大图 -->
            <PetDisplay
              :species-id="selectedPet.species"
              :level="selectedPet.level"
              :name="selectedPet.name"
              :mood="selectedPet.mood"
              :animate="true"
              size="md"
              @feed="handleFeed"
            />

            <!-- 宠物信息 -->
            <div class="pet-meta">
              <div class="meta-row">
                <span class="meta-label">👤 主人</span>
                <span class="meta-value">{{ selectedPet.student_name }}</span>
              </div>
              <div class="meta-row">
                <span class="meta-label">📦 系列</span>
                <span class="meta-value">{{ selectedPet.seriesName }}</span>
              </div>
              <div class="meta-row">
                <span class="meta-label">💬 心情</span>
                <span class="meta-value">{{ selectedPet.mood !== undefined ? (selectedPet.mood >= 60 ? '😊 开心' : selectedPet.mood >= 30 ? '😐 一般' : '😢 低落') : '-' }}</span>
              </div>
              <div class="meta-row">
                <span class="meta-label">🔥 连击</span>
                <span class="meta-value">{{ selectedPet.streakDays || 0 }} 天</span>
              </div>
            </div>

            <!-- 宠物语录 -->
            <div class="pet-quote">
              <div class="quote-bubble">
                "{{ currentQuote }}"
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="action-buttons">
              <button class="action-btn btn-feed" @click="handleFeed">
                <span class="btn-icon">🍬</span>
                <span>投喂</span>
              </button>
              <button class="action-btn btn-task" @click="handleAddScore">
                <span class="btn-icon">⭐</span>
                <span>加分</span>
              </button>
              <button class="action-btn btn-handbook" @click="openHandbook">
                <span class="btn-icon">📖</span>
                <span>图鉴</span>
              </button>
            </div>
          </div>

          <!-- 进化树 Tab -->
          <div v-if="activeTab === 'evolution'" class="detail-tab-content">
            <PetEvolutionTree
              :species-id="selectedPet.species"
              :current-level="selectedPet.level"
              :current-exp="selectedPet.exp"
              @select="(lvl) => console.log('select level', lvl)"
            />
          </div>

          <!-- Tab 切换 -->
          <div class="tab-switcher">
            <button
              class="tab-btn"
              :class="{ 'tab--active': activeTab === 'info' }"
              @click="activeTab = 'info'"
            >💬 概况</button>
            <button
              class="tab-btn"
              :class="{ 'tab--active': activeTab === 'evolution' }"
              @click="activeTab = 'evolution'"
            >🌟 进化</button>
          </div>
        </div>
      </Transition>

      <!-- 无选中时显示提示 -->
      <div v-if="!showDetailPanel && !loading && pets.length > 0" class="select-hint">
        <div class="hint-icon">👆</div>
        <p>点击左侧宠物卡片查看详情</p>
      </div>
    </div>

    <!-- 投喂动画 -->
    <PetFeedAnimation
      :show="showFeedAnimation"
      :points="feedPoints"
      :species-id="selectedPet?.species || 'zhulong'"
      :level="selectedPet?.level || 1"
      @close="showFeedAnimation = false"
    />

    <!-- 图鉴弹窗 -->
    <PetHandbook
      v-if="showHandbook && selectedPet"
      :species-id="selectedPet.species"
      :current-level="selectedPet.level"
      :current-score="selectedPet.exp"
      @close="showHandbook = false"
    />
  </div>
</template>

<style scoped>
.pets-page {
  position: relative;
  min-height: 100%;
}

/* 顶部 */
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}
.top-bar-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}
.page-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0;
}
.page-subtitle {
  font-size: 13px;
  color: var(--color-text-secondary);
}
.stat-pills {
  display: flex;
  gap: 8px;
}
.stat-pill {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
}
.pill-icon { font-size: 14px; }
.pill--avg { color: var(--color-primary); }
.pill--mature { color: #8B5CF6; }
.pill--legendary { color: #F59E0B; }

/* 系列筛选 */
.series-filter {
  display: flex;
  gap: 6px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.filter-chip {
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  color: var(--color-text-secondary);
  transition: all 0.2s ease;
}
.filter-chip:hover {
  color: var(--color-text);
  border-color: var(--color-text-secondary);
}
.chip--active {
  color: var(--color-text);
  font-weight: 600;
}

/* 主内容 */
.main-content {
  display: flex;
  gap: 20px;
  position: relative;
}

/* 宠物网格 */
.pet-grid-area {
  flex: 1;
  min-width: 0;
}
.pet-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.pet-card {
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  padding: 16px 12px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  border: 2px solid transparent;
  transition: all 0.3s var(--ease-smooth);
  min-height: 160px;
}
.pet-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  border-color: var(--card-accent);
}
.card--selected {
  border-color: var(--card-accent) !important;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--card-accent) 25%, transparent), 0 8px 24px rgba(0,0,0,0.1);
}
.card--legendary {
  box-shadow: 0 0 20px rgba(245,158,11,0.15);
}
.card--mature {
  box-shadow: 0 0 12px rgba(139,92,246,0.1);
}

.card-bg {
  position: absolute;
  inset: 0;
  opacity: 0.15;
  transition: opacity 0.3s ease;
}
.pet-card:hover .card-bg {
  opacity: 0.25;
}

.card-emoji {
  position: relative;
  z-index: 1;
  font-size: 40px;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.15));
  transition: transform 0.3s var(--ease-bounce);
}
.pet-card:hover .card-emoji {
  transform: scale(1.1);
}

.card-level {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  font-size: 9px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  background: rgba(0,0,0,0.3);
  color: white;
  backdrop-filter: blur(4px);
}
.stage--legendary {
  background: linear-gradient(135deg, rgba(245,158,11,0.6), rgba(239,68,68,0.6));
  animation: legendaryGlow 2s ease-in-out infinite;
}
@keyframes legendaryGlow {
  0%, 100% { box-shadow: 0 0 4px rgba(245,158,11,0.3); }
  50% { box-shadow: 0 0 12px rgba(245,158,11,0.6); }
}

.card-info {
  position: relative;
  z-index: 1;
  text-align: center;
}
.card-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text);
}
.card-pet-name {
  font-size: 10px;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
}

.card-exp {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 4px;
  background: rgba(0,0,0,0.06);
  border-radius: 2px;
  overflow: hidden;
}
.exp-track { width: 100%; height: 100%; }
.exp-fill {
  height: 100%;
  background: var(--card-accent, var(--gradient-primary));
  border-radius: 2px;
  transition: width 0.5s ease;
}

.legendary-glow {
  position: absolute;
  inset: -2px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(245,158,11,0.1), rgba(239,68,68,0.1));
  animation: legendaryBorder 2s ease-in-out infinite;
  pointer-events: none;
}
@keyframes legendaryBorder {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 0.6; }
}

/* 详情面板 */
.detail-panel {
  width: 360px;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  align-self: flex-start;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.detail-tab-content {
  margin-bottom: 12px;
}

/* 宠物信息 */
.pet-meta {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  padding: 16px;
  margin-top: 12px;
}
.meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
}
.meta-row + .meta-row {
  border-top: 1px solid var(--color-border);
}
.meta-label {
  font-size: 13px;
  color: var(--color-text-secondary);
}
.meta-value {
  font-size: 13px;
  font-weight: 600;
}

/* 语录 */
.pet-quote {
  margin: 12px 0;
}
.quote-bubble {
  padding: 12px 16px;
  background: var(--color-bg);
  border-radius: 16px;
  font-size: 13px;
  color: var(--color-text-secondary);
  font-style: italic;
  position: relative;
  text-align: center;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}
.action-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 8px;
  border-radius: 14px;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-secondary);
}
.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
.btn-icon {
  font-size: 22px;
}
.btn-feed:hover {
  border-color: #F59E0B;
  color: #F59E0B;
  background: rgba(245,158,11,0.05);
}
.btn-task:hover {
  border-color: #10B981;
  color: #10B981;
  background: rgba(16,185,129,0.05);
}
.btn-handbook:hover {
  border-color: #6366F1;
  color: #6366F1;
  background: rgba(99,102,241,0.05);
}

/* Tab 切换 */
.tab-switcher {
  display: flex;
  gap: 4px;
  background: var(--color-bg);
  border-radius: 12px;
  padding: 4px;
}
.tab-btn {
  flex: 1;
  padding: 8px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  background: transparent;
  color: var(--color-text-secondary);
  transition: all 0.2s ease;
}
.tab--active {
  background: var(--color-bg-card);
  color: var(--color-text);
  box-shadow: var(--shadow-sm);
}

/* 选择提示 */
.select-hint {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  color: var(--color-text-secondary);
  text-align: center;
}
.hint-icon {
  font-size: 48px;
  margin-bottom: 12px;
  animation: hintBounce 2s ease-in-out infinite;
}
@keyframes hintBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

/* 面板滑入动画 */
.panel-slide-enter-active {
  animation: slideIn 0.3s var(--ease-smooth);
}
.panel-slide-leave-active {
  animation: slideIn 0.2s var(--ease-smooth) reverse;
}
@keyframes slideIn {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}

/* 加载/空状态 */
.loading-state, .empty-state {
  text-align: center;
  padding: 60px 24px;
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
.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-hint { font-size: 12px; margin-top: 8px; opacity: 0.6; }

/* 响应式 */
@media (max-width: 900px) {
  .main-content {
    flex-direction: column;
  }
  .detail-panel {
    width: 100%;
    position: static;
    max-height: none;
  }
  .pet-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }
  .select-hint {
    width: 100%;
  }
}
</style>
