<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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

onMounted(async () => {
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
})
</script>

<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;flex-wrap:wrap;gap:12px;margin-bottom:24px;">
      <div style="display:flex;align-items:baseline;gap:12px;">
        <h2 style="font-size:26px;font-weight:700;margin:0;">🏆 年级战场</h2>
        <span style="font-size:14px;color:var(--md-text-secondary);">同年级大比拼</span>
      </div>
      <div v-if="!loading && !error" style="padding:6px 20px;border-radius:30px;background:rgba(245,158,11,0.08);border:1px solid rgba(245,158,11,0.15);color:#F59E0B;font-size:15px;font-weight:700;">
        {{ rankBadge }} 当前排名: #{{ myRank }}
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:60px;color:var(--md-text-secondary);">加载中...</div>
    <div v-else-if="error" style="text-align:center;padding:60px;color:#fca5a5;">
      <p>{{ error }}</p>
      <p style="font-size:13px;color:var(--md-text-secondary);margin-top:8px;">请确保班级有同年级的其他班级数据</p>
    </div>

    <template v-else-if="classes.length > 0">
      <div style="display:flex;flex-direction:column;gap:10px;margin-bottom:24px;">
        <div v-for="(cls, idx) in classes" :key="cls.name"
          :style="{
            display:'flex', alignItems:'center', gap:'14px', padding:'14px 20px',
            background: cls.isOwn ? 'linear-gradient(135deg,rgba(245,158,11,0.03),transparent)' : 'rgba(255,255,255,0.02)',
            border: '1px solid ' + (cls.isOwn ? 'rgba(245,158,11,0.2)' : 'rgba(255,255,255,0.04)'),
            borderRadius:'var(--md-radius)', transition:'0.25s',
          }">
          <div style="font-size:22px;font-weight:800;width:44px;text-align:center;"
            :style="{ color: idx === 0 ? '#F59E0B' : idx === 1 ? '#94A3B8' : idx === 2 ? '#D97706' : 'rgba(255,255,255,0.2)' }">
            {{ getMedal(idx) }}
          </div>
          <div style="width:130px;flex-shrink:0;">
            <div style="font-size:16px;font-weight:600;">{{ cls.name }}</div>
            <span v-if="cls.isOwn" style="font-size:10px;padding:1px 8px;border-radius:4px;background:rgba(245,158,11,0.1);color:#F59E0B;font-weight:600;">本班</span>
          </div>
          <div style="flex:1;display:flex;align-items:center;gap:12px;">
            <div style="flex:1;height:8px;background:rgba(255,255,255,0.06);border-radius:4px;overflow:hidden;">
              <div :style="{ width: (cls.totalScore / maxScore) * 100 + '%', height:'100%', background: idx === 0 ? 'linear-gradient(90deg,#f59e0b,#fcd34d)' : 'linear-gradient(90deg,var(--md-primary),var(--md-secondary))', borderRadius:'4px', transition:'width 0.8s' }"></div>
            </div>
            <span style="font-weight:700;font-size:15px;min-width:70px;text-align:right;">{{ cls.totalScore.toLocaleString() }}</span>
          </div>
          <div style="display:flex;gap:14px;font-size:13px;color:var(--md-text-secondary);min-width:160px;justify-content:flex-end;">
            <span title="人数">👤 {{ cls.studentCount }}</span>
            <span title="平均等级">📈 {{ cls.avgLevel.toFixed(1) }}</span>
            <span title="巅峰人数">⭐ {{ cls.peakCount }}</span>
          </div>
        </div>
      </div>

      <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;">
        <div style="background:rgba(255,255,255,0.02);border-radius:var(--md-radius);padding:20px 24px;border:1px solid rgba(255,255,255,0.04);">
          <h4 style="font-size:14px;color:var(--md-text-secondary);font-weight:600;text-transform:uppercase;letter-spacing:0.05em;margin-bottom:16px;">📊 本班战力分析</h4>
          <div style="display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid rgba(255,255,255,0.03);font-size:15px;">
            <span>总积分</span><span style="font-weight:700;">{{ myScore.toLocaleString() }}</span>
          </div>
          <div style="display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid rgba(255,255,255,0.03);font-size:15px;">
            <span>当前排名</span><span style="font-weight:700;">#{{ myRank }}</span>
          </div>
          <div style="display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid rgba(255,255,255,0.03);font-size:15px;">
            <span>与第1名差距</span><span style="font-weight:700;color:#F59E0B;">{{ gapToFirst.toLocaleString() }} 分</span>
          </div>
          <div v-if="classes[1]" style="display:flex;justify-content:space-between;padding:8px 0;font-size:15px;">
            <span>超越前1名需要</span>
            <span style="font-weight:700;">{{ (classes[myRank-2]?.totalScore || 0) - myScore > 0 ? ((classes[myRank-2]?.totalScore || 0) - myScore).toLocaleString() + ' 分' : '-' }}</span>
          </div>
        </div>
        <div style="background:rgba(255,255,255,0.02);border-radius:var(--md-radius);padding:20px 24px;border:1px solid rgba(255,255,255,0.04);">
          <h4 style="font-size:14px;color:var(--md-text-secondary);font-weight:600;text-transform:uppercase;letter-spacing:0.05em;margin-bottom:16px;">⚔️ 挑战建议</h4>
          <div v-if="gapToFirst > 0" style="font-size:15px;color:var(--md-text-secondary);line-height:1.6;">
            💪 距离第1名还差 <strong style="color:#F59E0B;font-size:18px;">{{ gapToFirst.toLocaleString() }}</strong> 分！<br>
            <span style="font-size:13px;opacity:0.7;">继续鼓励学生举手发言和完成作业！</span>
          </div>
          <div v-else style="font-size:15px;color:#10B981;line-height:1.6;">
            🎉 太棒了！本班目前位列年级第 {{ myRank }}！
          </div>
        </div>
      </div>
    </template>

    <div v-else style="text-align:center;padding:60px;color:var(--md-text-secondary);">
      📭 暂无同年级其他班级数据
    </div>
  </div>
</template>
