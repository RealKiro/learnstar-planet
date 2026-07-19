<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet } from '@/utils/api'
import { avatarGradient, getPetStageName } from '@/utils/constants'
import { getSpeciesEmoji, getSeriesBySpeciesId, SERIES_SCENES } from '@/utils/petData'
import type { ApiResponse } from '@/types'

interface ChildInfo {
  id: number
  name: string
  student_no: string
  class_name: string
  grade: string
  total_score: number
  class_rank: number
  pet_name: string
  pet_level: number
  pet_stage: string
  pet_emoji: string
  pet_mood: string
  pet_species?: string
}

const router = useRouter()
const children = ref<ChildInfo[]>([])
const loading = ref(true)
const selectedChild = ref<ChildInfo | null>(null)

const totalScore = computed(() => children.value.reduce((s, c) => s + (c.total_score || 0), 0))
const bestRank = computed(() => {
  const ranks = children.value.map(c => c.class_rank).filter(r => r > 0)
  return ranks.length ? Math.min(...ranks) : '-'
})

function getChildSeriesBg(speciesId?: string): string {
  if (!speciesId) return 'var(--gradient-primary)'
  const series = getSeriesBySpeciesId(speciesId)
  return series && SERIES_SCENES[series.id]?.bgGradient || 'var(--gradient-primary)'
}

function goToPetDetail(childId: number) {
  router.push({ name: 'parent-pet', query: { student_id: childId } })
}

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<{ children: ChildInfo[] }>>('/api/v1/parent/home')
    children.value = res.data?.children || []
    if (children.value.length > 0) {
      selectedChild.value = children.value[0]
    }
  } catch {
    // Demo data
    children.value = [
      {
        id: 1, name: '张小明', student_no: '001', class_name: '三年级一班',
        grade: '三年级', total_score: 184, class_rank: 3,
        pet_name: '烛龙', pet_level: 6, pet_stage: '衔烛之影',
        pet_emoji: '🐉', pet_mood: '开心', pet_species: 'zhulong',
      },
    ]
    selectedChild.value = children.value[0]
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div style="margin-bottom:24px;">
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">家长中心</p>
      <h2 style="font-size:24px;font-weight:700;">孩子动态</h2>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="children.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">👶</div>
      <p>尚未绑定孩子，请联系学校管理员</p>
    </div>

    <template v-else>
      <div class="stats-grid">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">👶</span>
          <div class="stat-card__value">{{ children.length }}</div>
          <div class="stat-card__label">孩子数量</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">⭐</span>
          <div class="stat-card__value">{{ totalScore.toLocaleString() }}</div>
          <div class="stat-card__label">总积分</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">🏆</span>
          <div class="stat-card__value">{{ bestRank }}</div>
          <div class="stat-card__label">最高排名</div>
        </div>
      </div>

      <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:16px;margin-top:24px;">
        <div v-for="c in children" :key="c.id"
          class="child-card"
          @click="goToPetDetail(c.id)"
        >
          <!-- 宠物展示区 -->
          <div class="pet-preview" :style="{ background: getChildSeriesBg(c.pet_species) }">
            <div class="pet-preview-emoji">
              {{ c.pet_species ? getSpeciesEmoji(c.pet_species) : (c.pet_emoji || '🌟') }}
            </div>
            <div class="pet-preview-level">Lv.{{ c.pet_level }}</div>
            <div class="pet-preview-name">{{ c.pet_stage || getPetStageName(c.pet_species || '', c.pet_level) }}</div>
          </div>

          <!-- 信息区 -->
          <div class="child-info">
            <div style="display:flex;align-items:center;gap:12px;margin-bottom:12px;">
              <div :style="{ width:'40px', height:'40px', borderRadius:'10px', background:avatarGradient(c.name), color:'white', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'16px', fontWeight:600 }">
                {{ c.name.charAt(0) }}
              </div>
              <div>
                <div style="font-weight:600;font-size:15px;">{{ c.name }}</div>
                <div style="font-size:12px;color:var(--color-text-secondary);">{{ c.class_name }}</div>
              </div>
            </div>

            <div style="display:flex;justify-content:space-between;padding:10px 0;border-top:1px solid var(--color-border);">
              <div>
                <div style="font-size:11px;color:var(--color-text-secondary);">总积分</div>
                <div style="font-size:18px;font-weight:700;color:var(--color-primary);">{{ c.total_score }}</div>
              </div>
              <div style="text-align:right;">
                <div style="font-size:11px;color:var(--color-text-secondary);">班级排名</div>
                <div style="font-size:18px;font-weight:700;color:var(--color-accent);">第{{ c.class_rank }}名</div>
              </div>
            </div>

            <button class="pet-interact-btn" @click.stop="goToPetDetail(c.id)">
              🌟 查看宠物详情
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.child-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s var(--ease-smooth);
}
.child-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.pet-preview {
  height: 140px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
.pet-preview-emoji {
  font-size: 56px;
  filter: drop-shadow(0 4px 12px rgba(0,0,0,0.2));
  animation: petFloat 3s ease-in-out infinite;
}
@keyframes petFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}
.pet-preview-level {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 10px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 6px;
  background: rgba(0,0,0,0.3);
  color: white;
  backdrop-filter: blur(4px);
}
.pet-preview-name {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  color: rgba(255,255,255,0.8);
  text-shadow: 0 2px 8px rgba(0,0,0,0.4);
  font-weight: 500;
  white-space: nowrap;
}

.child-info {
  padding: 16px;
}

.pet-interact-btn {
  width: 100%;
  margin-top: 12px;
  padding: 10px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--color-text-secondary);
}
.pet-interact-btn:hover {
  background: var(--gradient-primary);
  color: white;
  border-color: transparent;
}
</style>
