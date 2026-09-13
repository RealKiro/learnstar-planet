<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'

interface ClassPK {
  name: string; totalScore: number; studentCount: number
  avgLevel: number; peakCount: number; weekGrowth: number; isOwn: boolean
}

const classes = ref<ClassPK[]>([])
const loading = ref(true)
const token = ref('')
const myScore = ref(0)
const myRank = ref(0)
const error = ref('')

const maxScore = computed(() => Math.max(...classes.value.map(c => c.totalScore), 1))
const gapToFirst = computed(() => {
  const first = classes.value[0]
  if (!first || first.isOwn) return 0
  return first.totalScore - myScore.value
})
const rankBadge = computed(() => {
  const r = myRank.value
  if (r === 1) return '🥇'
  if (r === 2) return '🥈'
  if (r === 3) return '🥉'
  return `#${r}`
})

function getMedal(idx: number) { return ['🥇','🥈','🥉'][idx] || `#${idx+1}` }

let pollTimer: ReturnType<typeof setInterval> | null = null

async function fetchPKData() {
  token.value = sessionStorage.getItem('class_token') || ''
  if (!token.value) { loading.value = false; return }
  try {
    const res = await apiGet<{ data: ClassPK[] }>('/api/v1/display/pk/leaderboard', { params: { token: token.value } })
    classes.value = res.data || []
    const mine = classes.value.find(c => c.isOwn)
    if (mine) {
      myScore.value = mine.totalScore
      myRank.value = classes.value.indexOf(mine) + 1
    }
  } catch (e: any) {
    error.value = e?.response?.data?.message || '加载失败'
  } finally { loading.value = false }
}

onMounted(() => {
  fetchPKData()
  pollTimer = setInterval(fetchPKData, 15000)
})

onUnmounted(() => {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
})
</script>

<template>
  <div>
    <div class="page-head">
      <div class="title-row">
        <h2 class="page-title">🏆 年级战场</h2>
        <span class="text-muted-14">同年级大比拼</span>
      </div>
      <div v-if="!loading && !error" class="pk-badge">
        {{ rankBadge }} 当前排名: #{{ myRank }}
      </div>
    </div>

    <div v-if="loading" class="empty-60">加载中...</div>
    <div v-else-if="error" class="empty-60-error">
      <p>{{ error }}</p>
      <p class="hint-13-top">请确保班级有同年级的其他班级数据</p>
    </div>

    <template v-else-if="classes.length > 0">
      <div class="stack-10">
        <div v-for="(cls, idx) in classes" :key="cls.name"
          :style="{
            display:'flex', alignItems:'center', gap:'14px', padding:'14px 20px',
            background: cls.isOwn ? 'linear-gradient(135deg,rgba(245,158,11,0.03),transparent)' : 'var(--tint-1)',
            border: '1px solid ' + (cls.isOwn ? 'rgba(245,158,11,0.2)' : 'var(--tint-2)'),
            borderRadius:'var(--md-radius)', transition:'0.25s',
          }">
          <div 
 :style="{ color: idx === 0 ? '#F59E0B' : idx === 1 ? '#94A3B8' : idx === 2 ? '#D97706' : 'var(--color-text-secondary)' }" class="rank">
            {{ getMedal(idx) }}
          </div>
          <div class="team-name">
            <div class="fw-600-16">{{ cls.name }}</div>
            <span v-if="cls.isOwn" class="mini-badge">本班</span>
          </div>
          <div class="bar-wrap">
            <div class="bar-track">
              <div :style="{ width: (cls.totalScore / maxScore) * 100 + '%', height:'100%', background: idx === 0 ? 'linear-gradient(90deg,#f59e0b,#fcd34d)' : 'linear-gradient(90deg,var(--md-primary),var(--md-secondary))', borderRadius:'4px', transition:'width 0.8s' }"></div>
            </div>
            <span class="bar-value">{{ cls.totalScore.toLocaleString() }}</span>
          </div>
          <div class="bar-meta">
            <span title="人数">👤 {{ cls.studentCount }}</span>
            <span title="平均等级">📈 {{ cls.avgLevel.toFixed(1) }}</span>
            <span title="巅峰人数">⭐ {{ cls.peakCount }}</span>
          </div>
        </div>
      </div>

      <div class="grid-2-20">
        <div class="panel">
          <h4 class="panel-label">📊 本班战力分析</h4>
          <div class="stat-line">
            <span>总积分</span><span class="fw-700">{{ myScore.toLocaleString() }}</span>
          </div>
          <div class="stat-line">
            <span>当前排名</span><span class="fw-700">#{{ myRank }}</span>
          </div>
          <div class="stat-line">
            <span>与第1名差距</span><span class="fw-700-gold">{{ gapToFirst.toLocaleString() }} 分</span>
          </div>
          <div v-if="classes[1]" class="stat-line-plain">
            <span>超越前1名需要</span>
            <span class="fw-700">{{ (classes[myRank-2]?.totalScore || 0) - myScore > 0 ? ((classes[myRank-2]?.totalScore || 0) - myScore).toLocaleString() + ' 分' : '-' }}</span>
          </div>
        </div>
        <div class="panel">
          <h4 class="panel-label">⚔️ 挑战建议</h4>
          <div v-if="gapToFirst > 0" class="text-muted-15">
            💪 距离第1名还差 <strong class="gold-18">{{ gapToFirst.toLocaleString() }}</strong> 分！<br>
            <span class="faded-13">继续鼓励学生举手发言和完成作业！</span>
          </div>
          <div v-else class="text-ok-15">
            🎉 太棒了！本班目前位列年级第 {{ myRank }}！
          </div>
        </div>
      </div>
    </template>

    <div v-else class="empty-60">
      📭 暂无同年级其他班级数据
    </div>
  </div>
</template>

<style scoped>
/* ===== 页内布局类（本页专用，替代原内联样式；声明逐字保留以保证渲染等价） ===== */
.stat-line { display:flex; justify-content:space-between; padding:8px 0; border-bottom:1px solid var(--tint-1); font-size:15px; }
.fw-700 { font-weight:700; }
.empty-60 { text-align:center; padding:60px; color:var(--md-text-secondary); }
.panel { background:var(--tint-1); border-radius:var(--md-radius); padding:20px 24px; border:1px solid var(--tint-2); }
.panel-label { font-size:14px; color:var(--md-text-secondary); font-weight:600; text-transform:uppercase; letter-spacing:0.05em; margin-bottom:16px; }
.page-head { display:flex; justify-content:space-between; align-items:center; flex-wrap:wrap; gap:12px; margin-bottom:24px; }
.title-row { display:flex; align-items:baseline; gap:12px; }
.page-title { font-size:26px; font-weight:700; margin:0; }
.text-muted-14 { font-size:14px; color:var(--md-text-secondary); }
.pk-badge { padding:6px 20px; border-radius:30px; background:rgba(245,158,11,0.08); border:1px solid rgba(245,158,11,0.15); color: var(--c-amber); font-size:15px; font-weight:700; }
.empty-60-error { text-align:center; padding:60px; color: var(--color-danger-text); }
.hint-13-top { font-size:13px; color:var(--md-text-secondary); margin-top:8px; }
.stack-10 { display:flex; flex-direction:column; gap:10px; margin-bottom:24px; }
.rank { font-size:22px; font-weight:800; width:44px; text-align:center; }
.team-name { width:130px; flex-shrink:0; }
.fw-600-16 { font-size:16px; font-weight:600; }
.mini-badge { font-size:10px; padding:1px 8px; border-radius:4px; background:rgba(245,158,11,0.1); color: var(--c-amber); font-weight:600; }
.bar-wrap { flex:1; display:flex; align-items:center; gap:12px; }
.bar-track { flex:1; height:8px; background:var(--tint-3); border-radius:4px; overflow:hidden; }
.bar-value { font-weight:700; font-size:15px; min-width:70px; text-align:right; }
.bar-meta { display:flex; gap:14px; font-size:13px; color:var(--md-text-secondary); min-width:160px; justify-content:flex-end; }
.grid-2-20 { display:grid; grid-template-columns:1fr 1fr; gap:20px; }
.fw-700-gold { font-weight:700; color: var(--c-amber); }
.stat-line-plain { display:flex; justify-content:space-between; padding:8px 0; font-size:15px; }
.text-muted-15 { font-size:15px; color:var(--md-text-secondary); line-height:1.6; }
.gold-18 { color: var(--c-amber); font-size:18px; }
.faded-13 { font-size:13px; opacity:0.7; }
.text-ok-15 { font-size:15px; color: var(--c-green); line-height:1.6; }
</style>
