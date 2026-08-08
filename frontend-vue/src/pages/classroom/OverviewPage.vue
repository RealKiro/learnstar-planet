<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'
import PetSprite from '@/components/pet/PetSprite.vue'
import PetHandbook from '@/components/pet/PetHandbook.vue'

interface CardStudent {
  name: string; student_no?: string; score: number
  pet_name?: string; pet_species: string; pet_level: number
}
interface OverviewData {
  class_name: string; grade: string; student_count: number
  total_score: number; avg_pet_level: number; peak_count: number; weekly_score: number
  star_student: CardStudent | null
  top5: CardStudent[]
  recent_news: Array<{ icon: string; text: string }>
}

const data = ref<OverviewData | null>(null)
const loading = ref(true)
const token = ref('')
// 点击宠物 SVG 打开的图鉴弹窗
const handbook = ref<{ speciesId: string; level: number; score: number } | null>(null)
function openHandbook(s: CardStudent) {
  if (!s.pet_species) return
  handbook.value = { speciesId: s.pet_species, level: s.pet_level, score: s.score }
}

// 冠亚季军奖牌色（金/银/铜）
const MEDALS = ['#F59E0B', '#A8B0B8', '#CD7F32']

onMounted(async () => {
  token.value = sessionStorage.getItem('class_token') || ''
  if (!token.value) return
  await fetchData()
})

let pollTimer: ReturnType<typeof setInterval> | null = null

async function fetchData() {
  try {
    const res = await apiGet<{ data: OverviewData }>('/api/v1/display/dashboard', { params: { token: token.value } })
    data.value = res.data
  } catch { /* ignore */ } finally { loading.value = false }
}

onMounted(() => {
  pollTimer = setInterval(fetchData, 10000)
})

onUnmounted(() => {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
})
</script>

<template>
  <div>
    <div class="page-header" style="display:flex;align-items:baseline;gap:12px;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;margin:0;">🏠 班级总览</h2>
      <span v-if="data" style="font-size:14px;color:var(--md-text-secondary);">{{ data.class_name }} · {{ data.grade }}</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:60px;color:var(--md-text-secondary);">加载中...</div>

    <template v-else-if="data">
      <div class="overview-grid" style="display:grid;grid-template-columns:1.2fr 1.4fr 1fr;gap:20px;margin-bottom:28px;">
        <!-- 班级之星 -->
        <div class="o-card" v-if="data.star_student" style="position:relative;">
          <div class="o-label">🏅 班级之星</div>
          <div v-if="data.star_student.student_no" style="position:absolute;top:20px;right:24px;font-size:11px;color:var(--md-text-secondary);background:var(--tint-2);padding:2px 10px;border-radius:8px;">学号 {{ data.star_student.student_no }}</div>
          <div style="display:flex;align-items:center;gap:16px;">
            <div style="width:100px;height:100px;border-radius:50%;background:linear-gradient(135deg,#f093fb,#f5576c);display:flex;align-items:center;justify-content:center;overflow:hidden;cursor:pointer;box-shadow:0 4px 14px rgba(240,147,251,0.25);" @click="openHandbook(data.star_student)" title="点击查看宠物介绍">
              <div v-if="data.star_student.pet_species" style="width:92px;height:92px;border-radius:50%;overflow:hidden;">
                <PetSprite :species-id="data.star_student.pet_species" :level="data.star_student.pet_level" :animate="true" />
              </div>
              <span v-else style="font-size:36px;">🌟</span>
            </div>
            <div>
              <div style="font-size:20px;font-weight:700;">{{ data.star_student.name }}</div>
              <div style="font-size:13px;color:var(--md-text-secondary);">Lv.{{ data.star_student.pet_level }}</div>
              <div style="font-size:22px;font-weight:800;color:var(--md-gold);">{{ data.star_student.score }} 分</div>
            </div>
          </div>
        </div>

        <!-- 班级概况 -->
        <div class="o-card">
          <div class="o-label">📊 班级概况</div>
          <div style="font-size:36px;font-weight:800;line-height:1;margin-bottom:4px;">{{ data.total_score.toLocaleString() }}</div>
          <div style="font-size:13px;color:var(--md-text-secondary);margin-bottom:14px;">总积分 · 共 {{ data.student_count }} 人</div>
          <div style="display:flex;gap:20px;padding-top:12px;border-top:1px solid var(--tint-3);">
            <div><span style="display:block;font-size:11px;color:var(--md-text-secondary);">平均等级</span><strong style="font-size:18px;">{{ data.avg_pet_level.toFixed(1) }}</strong></div>
            <div><span style="display:block;font-size:11px;color:var(--md-text-secondary);">巅峰 Lv.10+</span><strong style="font-size:18px;color:#8B5CF6;">{{ data.peak_count }}</strong></div>
            <div><span style="display:block;font-size:11px;color:var(--md-text-secondary);">本周增长</span><strong style="font-size:18px;color:#10B981;">+{{ data.weekly_score }}</strong></div>
          </div>
        </div>

        <!-- 最新动态 -->
        <div class="o-card">
          <div class="o-label">📢 最新动态</div>
          <div style="display:flex;flex-direction:column;gap:8px;max-height:160px;overflow-y:auto;">
            <div v-for="(n, i) in data.recent_news" :key="i" style="display:flex;align-items:center;gap:8px;padding:8px 10px;background:var(--tint-1);border-radius:10px;font-size:13px;border-left:2px solid var(--md-primary);">
              <span>{{ n.icon }}</span>
              <span style="color:var(--md-text-secondary);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ n.text }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- TOP 5：前三名大卡片 + 4-5 名次位 -->
      <div class="top5-wrap">
        <div class="top5-title">🏆 班级 TOP 5</div>
        <!-- 冠亚季军大卡片 -->
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
      </div>

      <!-- 宠物图鉴弹窗 -->
      <PetHandbook v-if="handbook" :species-id="handbook.speciesId" :current-level="handbook.level" :current-score="handbook.score" @close="handbook = null" />
    </template>
  </div>
</template>

<style scoped>
.o-card {
  background: var(--tint-1);
  border-radius: var(--md-radius);
  padding: 24px;
  border: 1px solid var(--tint-2);
  backdrop-filter: blur(8px);
}
.o-label {
  font-size: 13px;
  color: var(--md-text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 14px;
}

/* TOP5 前三名大卡片 + 4-5 名次位 */
.top5-wrap {
  background: var(--tint-1);
  border-radius: var(--md-radius);
  padding: 16px 24px;
  border: 1px solid var(--tint-2);
}
.top5-title {
  font-size: 15px;
  font-weight: 700;
  margin-bottom: 16px;
}
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
  border: 1px solid var(--tint-2);
  cursor: pointer;
  transition: 0.25s;
}
.top3-card:hover { transform: translateY(-3px); box-shadow: 0 8px 20px rgba(0,0,0,0.08); }
.medal--0 { background: linear-gradient(180deg, rgba(245,158,11,0.08), transparent); border-color: rgba(245,158,11,0.25); }
.medal--1 { background: linear-gradient(180deg, rgba(168,176,184,0.08), transparent); border-color: rgba(168,176,184,0.25); }
.medal--2 { background: linear-gradient(180deg, rgba(205,127,50,0.08), transparent); border-color: rgba(205,127,50,0.25); }
.top3-medal {
  position: absolute;
  top: 8px;
  left: 12px;
  font-size: 22px;
}
.top3-avatar {
  width: 76px;
  height: 76px;
  border-radius: 50%;
  overflow: hidden;
  margin: 6px auto 10px;
  border: 2.5px solid var(--medal, #6B7280);
  background: var(--tint-1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px color-mix(in srgb, var(--medal, #6B7280) 22%, transparent);
}
.top3-emoji { font-size: 34px; }
.top3-name { font-size: 16px; font-weight: 700; }
.top3-no {
  font-size: 10px;
  color: var(--md-text-secondary);
  background: var(--tint-2);
  padding: 1px 8px;
  border-radius: 8px;
  display: inline-block;
  margin-top: 2px;
}
.top3-level { font-size: 11px; color: var(--md-text-secondary); margin-top: 2px; }
.top3-score { font-size: 18px; font-weight: 800; color: var(--md-gold); margin-top: 2px; }
.top3-bar { height: 4px; background: var(--tint-3); border-radius: 2px; overflow: hidden; margin-top: 8px; }
.top3-fill {
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(90deg, var(--md-primary), var(--md-secondary));
  transition: width 0.5s;
}
.medal--0 .top3-fill { background: linear-gradient(90deg, #f59e0b, #fcd34d); }

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
  border: 1px solid var(--tint-2);
  cursor: pointer;
  transition: 0.25s;
}
.top4-card:hover { transform: translateY(-2px); box-shadow: 0 6px 16px rgba(0,0,0,0.06); }
.top4-rank {
  font-size: 22px;
  font-weight: 800;
  color: var(--md-text-secondary);
  opacity: 0.6;
  min-width: 28px;
}
.top4-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid var(--tint-3);
  background: var(--tint-1);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.top4-emoji { font-size: 22px; }
.top4-info { min-width: 0; }
.top4-name { font-size: 14px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.top4-score { font-size: 13px; font-weight: 700; color: var(--md-gold); }

@media (max-width: 900px) {
  .overview-grid { grid-template-columns: 1fr; }
  .top3-row { grid-template-columns: 1fr; }
  .top2-row { grid-template-columns: 1fr; }
}
</style>
