<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'
import { getSpeciesEmoji } from '@/utils/petData'

interface OverviewData {
  class_name: string; grade: string; student_count: number
  total_score: number; avg_pet_level: number; peak_count: number; weekly_score: number
  star_student: { name: string; pet_name: string; pet_species: string; pet_level: number; score: number } | null
  top5: Array<{ name: string; score: number; pet_name: string; pet_species: string; pet_level: number }>
  recent_news: Array<{ icon: string; text: string }>
}

const data = ref<OverviewData | null>(null)
const loading = ref(true)
const token = ref('')

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
        <div class="o-card" v-if="data.star_student">
          <div class="o-label">🏅 班级之星</div>
          <div style="display:flex;align-items:center;gap:16px;">
            <div style="width:64px;height:64px;border-radius:50%;background:linear-gradient(135deg,#f093fb,#f5576c);display:flex;align-items:center;justify-content:center;font-size:32px;">
              {{ data.star_student.pet_species ? getSpeciesEmoji(data.star_student.pet_species) : '🌟' }}
            </div>
            <div>
              <div style="font-size:20px;font-weight:700;">{{ data.star_student.name }}</div>
              <div style="font-size:13px;color:var(--md-text-secondary);">{{ data.star_student.pet_name }} · Lv.{{ data.star_student.pet_level }}</div>
              <div style="font-size:22px;font-weight:800;color:var(--md-gold);">{{ data.star_student.score }} 分</div>
            </div>
          </div>
        </div>

        <!-- 班级概况 -->
        <div class="o-card">
          <div class="o-label">📊 班级概况</div>
          <div style="font-size:36px;font-weight:800;line-height:1;margin-bottom:4px;">{{ data.total_score.toLocaleString() }}</div>
          <div style="font-size:13px;color:var(--md-text-secondary);margin-bottom:14px;">总积分 · 共 {{ data.student_count }} 人</div>
          <div style="display:flex;gap:20px;padding-top:12px;border-top:1px solid rgba(255,255,255,0.06);">
            <div><span style="display:block;font-size:11px;color:var(--md-text-secondary);">平均等级</span><strong style="font-size:18px;">{{ data.avg_pet_level.toFixed(1) }}</strong></div>
            <div><span style="display:block;font-size:11px;color:var(--md-text-secondary);">巅峰 Lv.10+</span><strong style="font-size:18px;color:#8B5CF6;">{{ data.peak_count }}</strong></div>
            <div><span style="display:block;font-size:11px;color:var(--md-text-secondary);">本周增长</span><strong style="font-size:18px;color:#10B981;">+{{ data.weekly_score }}</strong></div>
          </div>
        </div>

        <!-- 最新动态 -->
        <div class="o-card">
          <div class="o-label">📢 最新动态</div>
          <div style="display:flex;flex-direction:column;gap:8px;max-height:160px;overflow-y:auto;">
            <div v-for="(n, i) in data.recent_news" :key="i" style="display:flex;align-items:center;gap:8px;padding:8px 10px;background:rgba(255,255,255,0.02);border-radius:10px;font-size:13px;border-left:2px solid var(--md-primary);">
              <span>{{ n.icon }}</span>
              <span style="color:var(--md-text-secondary);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ n.text }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- TOP 5 -->
      <div style="background:rgba(255,255,255,0.02);border-radius:var(--md-radius);padding:16px 24px;border:1px solid rgba(255,255,255,0.04);">
        <div style="font-size:15px;font-weight:700;margin-bottom:16px;">🏆 班级 TOP 5</div>
        <div style="display:grid;grid-template-columns:repeat(5,1fr);gap:12px;">
          <div v-for="(s, i) in data.top5" :key="s.name"
            style="text-align:center;padding:16px 12px;border-radius:16px;border:1px solid rgba(255,255,255,0.04);transition:0.25s;"
            :style="i === 0 ? 'background:linear-gradient(180deg,rgba(245,158,11,0.06),transparent);border-color:rgba(245,158,11,0.2);' : ''">
            <div style="font-size:24px;margin-bottom:6px;">{{ ['🥇','🥈','🥉','4','5'][i] }}</div>
            <div style="font-size:28px;margin-bottom:6px;">{{ s.pet_species ? getSpeciesEmoji(s.pet_species) : '🌟' }}</div>
            <div style="font-size:14px;font-weight:600;">{{ s.name }}</div>
            <div style="font-size:11px;color:var(--md-text-secondary);margin-bottom:6px;">{{ s.pet_name }} · Lv.{{ s.pet_level }}</div>
            <div style="font-size:16px;font-weight:700;color:var(--md-gold);">{{ s.score }} 分</div>
            <div style="height:4px;background:rgba(255,255,255,0.06);border-radius:2px;overflow:hidden;margin-top:6px;">
              <div :style="{ width: (s.score / data.top5[0].score) * 100 + '%', height:'100%', background: i === 0 ? 'linear-gradient(90deg,#f59e0b,#fcd34d)' : 'linear-gradient(90deg,var(--md-primary),var(--md-secondary))', borderRadius:'2px', transition:'width 0.5s' }"></div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.o-card {
  background: rgba(255,255,255,0.03);
  border-radius: var(--md-radius);
  padding: 24px;
  border: 1px solid rgba(255,255,255,0.04);
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
@media (max-width: 900px) {
  .overview-grid { grid-template-columns: 1fr; }
}
</style>
