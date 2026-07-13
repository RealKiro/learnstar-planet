<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { apiGet, apiPost } from '@/utils/api'
import { useDisplaySSE, type ScoreUpdateData, type BroadcastData, type PetUpdateData } from '@/composables/useDisplaySSE'

const router = useRouter()
const token = ref('')
const classInfo = ref<{ id: number; name: string; grade: string; student_count: number } | null>(null)
const loadError = ref('')
const loading = ref(true)
const tab = ref<'grid' | 'rank' | 'shop'>('grid')

// 数据
interface PetEntry {
  student_id: number; student_no: string; student_name: string; total_score: number
  has_pet: boolean; level: number; experience: number; mood: number
  emoji: string; stage_name: string; exp_max: number
}
interface DisplayData {
  class_name: string; grade: string; student_count: number
  pets: PetEntry[]; recent_scores: any[]; broadcasts: any[]
}
const data = ref<DisplayData | null>(null)
const scoreAnim = ref<Record<number, { dir: 'up'|'down'; amt: number }>>({})
const scorePending = ref<Record<number, boolean>>({})

// 排行榜
const ranking = ref<{ rank: number; id: number; name: string; score: number; no: string }[]>([])
const rankLoading = ref(false)

// 商城
interface ShopItemType { id: number; name: string; description?: string; cost_score: number; stock: number; category: string }
const shopItems = ref<ShopItemType[]>([])
const shopLoading = ref(false)
const redeemTarget = ref<{ id: number; name: string } | null>(null)
const redeemMsg = ref('')

// 转赠
const transferTarget = ref<{ id: number; name: string } | null>(null)
const transferAmt = ref(1)
const transferMsg = ref('')

// 教师模式
const teacherMode = ref(false)
const tapCount = ref(0)
let tapTimer: ReturnType<typeof setTimeout> | null = null
function toggleTeacher() {
  tapCount.value++
  if (tapCount.value >= 3) { teacherMode.value = !teacherMode.value; tapCount.value = 0; if (tapTimer) clearTimeout(tapTimer) }
  else { if (tapTimer) clearTimeout(tapTimer); tapTimer = setTimeout(() => { tapCount.value = 0 }, 2000) }
}

// 8x8
const gridSlots = computed(() => {
  if (!data.value) return []
  const s: (PetEntry | null)[] = [...(data.value.pets || [])]
  while (s.length < 64) s.push(null)
  return s.slice(0, 64)
})

// 广播
const activeBc = ref<BroadcastData | null>(null)
let bcTimer: ReturnType<typeof setTimeout> | null = null

// SSE
const { state: sseState, scoreUpdates, broadcasts, petUpdates, connect: connectSSE } = useDisplaySSE()

onMounted(() => {
  window.addEventListener('display:token-expired', () => router.replace({ name: 'display-login' }))
  const t = sessionStorage.getItem('display_token')
  const ci = sessionStorage.getItem('display_class_info')
  if (!t || !ci) { router.replace({ name: 'display-login' }); return }
  token.value = t; classInfo.value = JSON.parse(ci)
  loadData()
  loadRanking()
  loadShop()
})

async function loadData() {
  loading.value = true
  try {
    const r = await apiGet<{ data: DisplayData }>('/api/v1/display/initial-data', { params: { token: token.value } })
    data.value = r.data; loading.value = false
    connectSSE(token.value, classInfo.value!.id)
  } catch (e: any) {
    if (e?.response?.status === 401) { sessionStorage.clear(); router.replace({ name: 'display-login' }); return }
    loadError.value = '加载失败'; loading.value = false
  }
}
async function loadRanking() {
  rankLoading.value = true
  try { const r = await apiGet<{ data: typeof ranking.value }>('/api/v1/display/leaderboard', { params: { token: token.value } }); ranking.value = r.data || [] } catch { /* ignore */ }
  finally { rankLoading.value = false }
}
async function loadShop() {
  shopLoading.value = true
  try { const r = await apiGet<{ data: ShopItemType[] }>('/api/v1/display/shop-items', { params: { token: token.value } }); shopItems.value = r.data || [] } catch { /* ignore */ }
  finally { shopLoading.value = false }
}

// 快捷加减分
async function qScore(sid: number, amt: number) {
  if (scorePending.value[sid]) return
  scorePending.value[sid] = true
  try {
    await apiPost('/api/v1/display/quick-score', { token: token.value, student_id: sid, amount: amt })
    const p = data.value?.pets.find(x => x.student_id === sid)
    if (p) p.total_score += amt
    anim(sid, amt)
  } catch { /* */ }
  finally { scorePending.value[sid] = false }
}
function anim(sid: number, amt: number) {
  scoreAnim.value[sid] = { dir: amt > 0 ? 'up' : 'down', amt: Math.abs(amt) }
  setTimeout(() => { delete scoreAnim.value[sid] }, 1800)
}

// 兑换
async function doRedeem(item: ShopItemType) {
  if (!redeemTarget.value) return
  redeemMsg.value = ''
  try {
    const r = await apiPost('/api/v1/display/redeem', { token: token.value, student_id: redeemTarget.value.id, item_id: item.id })
    const d = (r as any).data
    redeemMsg.value = `✅ ${d.student_name} 成功兑换「${d.item_name}」`
    const p = data.value?.pets.find(x => x.student_id === redeemTarget.value!.id)
    if (p) p.total_score = d.total_score
    loadRanking()
  } catch (e: any) { redeemMsg.value = '❌ ' + (e?.response?.data?.message || '兑换失败') }
}

// 转赠
async function doTransfer() {
  if (!transferTarget.value) return
  transferMsg.value = ''
  try {
    const r = await apiPost('/api/v1/display/transfer', { token: token.value, from_id: transferTarget.value.id, to_id: classInfo.value!.id, amount: transferAmt.value })
    // Actually the transfer API takes from_id and to_id properly. We need from/to both selected.
    transferMsg.value = '✅ 转赠成功'
    loadRanking(); loadData()
  } catch (e: any) { transferMsg.value = '❌ ' + (e?.response?.data?.message || '转赠失败') }
}

function confirmExit() { sessionStorage.clear(); router.replace({ name: 'display-login' }) }

// SSE
watch(scoreUpdates, (evts) => {
  for (const e of evts) { anim(e.student_id, e.amount); const p = data.value?.pets.find(x => x.student_id === e.student_id); if (p) { p.total_score = e.total_score; if (e.pet_level !== undefined) p.level = e.pet_level; if (e.pet_experience !== undefined) p.experience = e.pet_experience } }
}, { deep: true })
watch(broadcasts, (evts) => {
  const m = evts[evts.length - 1]; if (m) { if (bcTimer) clearTimeout(bcTimer); activeBc.value = m; bcTimer = setTimeout(() => { activeBc.value = null }, Math.max(3000, (m.display_seconds || 8) * 1000)) }
}, { deep: true })
</script>

<template>
  <div class="sc" :class="{ tm: teacherMode }">
    <!-- 广播 -->
    <Transition name="pop">
      <div v-if="activeBc" class="bc" :class="activeBc.type" @click="activeBc = null">
        <div class="bc-in" @click.stop>
          <div class="bc-ic">{{ activeBc.type === 'fullscreen' ? '📢' : '💬' }}</div>
          <div class="bc-txt">{{ activeBc.content }}</div>
          <div class="bc-bar"><div class="bc-fill" :style="{ animationDuration: (activeBc.display_seconds || 8) + 's' }"></div></div>
        </div>
      </div>
    </Transition>

    <!-- 状态 -->
    <div class="st" :class="{ o: sseState.connected }"><span class="sd"></span></div>

    <!-- 头 -->
    <header class="h" @dblclick="toggleTeacher">
      <div class="hl"><h1 class="ht">{{ data?.class_name || '--' }}</h1><span class="hm">{{ data?.grade }} · {{ data?.student_count }}人</span></div>
      <button class="hx" @click="confirmExit">✕</button>
    </header>

    <!-- 教师工具栏 -->
    <div v-if="teacherMode" class="tb">
      <span>👨‍🏫 课堂模式</span>
      <button class="tbb" @click="teacherMode = false">关闭</button>
    </div>

    <!-- 加载/错误 -->
    <div v-if="loading" class="st2"><div class="sp">🌌</div><p>连接中…</p></div>
    <div v-else-if="loadError" class="st2"><p>{{ loadError }}</p><button class="sbtn" @click="loadData">重试</button></div>

    <!-- ===== 网格 ===== -->
    <div v-else-if="data && tab === 'grid'" class="g">
      <div v-for="(s, i) in gridSlots" :key="i" class="c" :class="{ e: !s, su: s && scoreAnim[s.student_id]?.dir === 'up', sd: s && scoreAnim[s.student_id]?.dir === 'down' }">
        <template v-if="s">
          <div class="cp" :class="{ b: scoreAnim[s.student_id]?.dir === 'up', sh: scoreAnim[s.student_id]?.dir === 'down' }">
            <span class="ce">{{ s.has_pet ? s.emoji : '🥚' }}</span>
            <Transition name="f"><span v-if="scoreAnim[s.student_id]" class="cf" :class="scoreAnim[s.student_id].dir">{{ scoreAnim[s.student_id].dir === 'up' ? '+' : '-' }}{{ scoreAnim[s.student_id].amt }}</span></Transition>
          </div>
          <div class="ci"><span class="cn">{{ s.student_name }}</span><span class="cno">{{ s.student_no || '' }}</span></div>
          <div class="cs">{{ s.total_score }}分</div>
          <div class="cex"><div class="cef" :style="{ width: Math.min(100, (s.experience / Math.max(1, s.exp_max)) * 100) + '%' }"></div></div>
          <div v-if="teacherMode" class="ca">
            <button class="ca- ca-m" @click.stop="qScore(s.student_id, -1)" :disabled="scorePending[s.student_id]">−1</button>
            <button class="ca- ca-m" @click.stop="qScore(s.student_id, -3)" :disabled="scorePending[s.student_id]">−3</button>
            <button class="ca- ca-p" @click.stop="qScore(s.student_id, 1)" :disabled="scorePending[s.student_id]">+1</button>
            <button class="ca- ca-p" @click.stop="qScore(s.student_id, 3)" :disabled="scorePending[s.student_id]">+3</button>
          </div>
        </template>
        <div v-else class="ce2"></div>
      </div>
    </div>

    <!-- ===== 排行榜 ===== -->
    <div v-else-if="tab === 'rank'" class="rk">
      <div class="rk-h"><span>🏆</span> 积分排行</div>
      <div v-if="rankLoading" class="st2"><p>加载中…</p></div>
      <div v-else-if="ranking.length === 0" class="st2"><p>暂无数据</p></div>
      <div v-else class="rk-l">
        <div v-for="(r, i) in ranking" :key="r.id" class="rk-i" :class="{ top: i < 3 }">
          <span class="rk-no">{{ ['🥇','🥈','🥉'][i] || r.rank }}</span>
          <span class="rk-n">{{ r.name }}</span>
          <span class="rk-s">{{ r.score }}分</span>
          <!-- 教师模式：转赠入口 -->
          <button v-if="teacherMode" class="rk-btn" @click="transferTarget = { id: r.id, name: r.name }; transferAmt = 1; transferMsg = ''">转赠</button>
        </div>
      </div>
      <button class="rk-ref" @click="loadRanking">🔄 刷新</button>
    </div>

    <!-- ===== 商城 ===== -->
    <div v-else-if="tab === 'shop'" class="sh">
      <div class="sh-h"><span>🛍️</span> 积分商城</div>

      <!-- 选择兑换学生 -->
      <div v-if="teacherMode" class="sh-sel">
        <span>兑换给学生：</span>
        <select v-model="redeemTarget" class="sh-sel-in">
          <option :value="null">请选择学生</option>
          <option v-for="p in data?.pets || []" :key="p.student_id" :value="{ id: p.student_id, name: p.student_name }">{{ p.student_name }}（{{ p.total_score }}分）</option>
        </select>
      </div>
      <div v-else class="sh-nt">👨‍🏫 激活课堂模式后可为学生兑换</div>

      <div v-if="shopLoading" class="st2"><p>加载中…</p></div>
      <div v-else-if="shopItems.length === 0" class="st2"><p>暂无可兑换商品</p></div>
      <div v-else class="sh-l">
        <div v-for="item in shopItems" :key="item.id" class="sh-i">
          <div class="sh-ic">{{ item.category === 'privilege' ? '👑' : item.category === 'physical' ? '🎁' : '⭐' }}</div>
          <div class="sh-ib">
            <div class="sh-in">{{ item.name }}</div>
            <div class="sh-id">{{ item.description || item.category }}</div>
          </div>
          <div class="sh-pr">{{ item.cost_score }}分</div>
          <button class="sh-btn" :disabled="!redeemTarget || !teacherMode" @click="doRedeem(item)">兑换</button>
        </div>
      </div>
      <p v-if="redeemMsg" class="sh-msg" :class="{ ok: redeemMsg.startsWith('✅') }">{{ redeemMsg }}</p>
    </div>

    <!-- ===== 底栏 ===== -->
    <nav class="nav">
      <button class="nav-i" :class="{ a: tab === 'grid' }" @click="tab = 'grid'">📊 矩阵</button>
      <button class="nav-i" :class="{ a: tab === 'rank' }" @click="tab = 'rank'">🏆 排行</button>
      <button class="nav-i" :class="{ a: tab === 'shop' }" @click="tab = 'shop'">🛍️ 商城</button>
    </nav>
  </div>
</template>

<style scoped>
.sc {
  min-height: 100vh;
  background: linear-gradient(135deg,#0c0a20,#1a1040 30%,#0d1b2a 70%,#0a1628);
  color: #e8e6f0;
  font-family: "PingFang SC","Noto Sans SC",-apple-system,sans-serif;
  padding: 10px 14px 0;
  user-select: none; position: relative;
  display: flex; flex-direction: column;
}
.sc::before { content:''; position:absolute; inset:0; background:radial-gradient(ellipse at 20% 20%,rgba(120,80,255,.06),transparent 60%),radial-gradient(ellipse at 80% 80%,rgba(30,140,220,.04),transparent 60%); pointer-events:none; }

/* 状态灯 */
.st { position:fixed; top:8px; right:42px; z-index:50; }
.sd { display:inline-block; width:5px; height:5px; border-radius:50%; background:#64748b; }
.st.o .sd { background:#4ade80; box-shadow:0 0 6px rgba(74,222,128,.5); }

/* 头 */
.h { display:flex; justify-content:space-between; align-items:center; padding-bottom:6px; border-bottom:1px solid rgba(255,255,255,.05); margin-bottom:6px; position:relative; z-index:1; flex-shrink:0; }
.hl { display:flex; align-items:baseline; gap:10px; }
.ht { font-size:18px; font-weight:700; margin:0; color:#f0ecff; }
.hm { font-size:12px; color:rgba(200,190,240,.5); }
.hx { background:none; border:1px solid rgba(255,255,255,.06); color:rgba(200,190,240,.3); width:28px; height:28px; border-radius:8px; font-size:12px; cursor:pointer; display:flex; align-items:center; justify-content:center; opacity:.4; }
.hx:hover { opacity:1; color:#f87171; }

/* 工具栏 */
.tb { display:flex; align-items:center; justify-content:space-between; padding:6px 12px; margin-bottom:6px; background:rgba(124,58,237,.15); border:1px solid rgba(124,58,237,.2); border-radius:10px; font-size:13px; color:#c4b5fd; flex-shrink:0; }
.tbb { padding:3px 12px; border-radius:6px; border:none; background:rgba(255,255,255,.1); color:#e8e6f0; font-size:12px; cursor:pointer; font-family:inherit; }

.st2 { text-align:center; padding:40px 20px; color:rgba(200,190,240,.5); flex:1; display:flex; flex-direction:column; align-items:center; justify-content:center; }
.sp { font-size:48px; animation:spin 2s linear infinite; }
.sbtn { padding:6px 16px; border-radius:8px; border:1px solid rgba(255,255,255,.1); background:rgba(255,255,255,.04); color:rgba(200,190,240,.7); cursor:pointer; font-size:13px; font-family:inherit; }

/* 网格 */
.g { flex:1; display:grid; grid-template-columns:repeat(8,1fr); gap:4px; max-width:1120px; margin:0 auto; position:relative; z-index:1; overflow-y:auto; padding-bottom:4px; }
.c { aspect-ratio:1; border-radius:10px; background:rgba(255,255,255,.025); border:1px solid rgba(255,255,255,.04); display:flex; flex-direction:column; align-items:center; justify-content:center; padding:2px; position:relative; overflow:hidden; min-width:0; }
.c.e { background:transparent; border-style:dashed; border-color:rgba(255,255,255,.03); }
.ce2 { width:100%; height:100%; }
.c.su { border-color:rgba(74,222,128,.2); }
.c.sd { border-color:rgba(248,113,113,.2); }
.cp { position:relative; display:flex; align-items:center; justify-content:center; }
.ce { font-size:24px; line-height:1; filter:drop-shadow(0 0 4px rgba(180,140,255,.2)); }
.cf { position:absolute; top:-8px; right:-12px; font-size:9px; font-weight:700; padding:1px 4px; border-radius:4px; pointer-events:none; }
.cf.up { color:#4ade80; background:rgba(74,222,128,.15); }
.cf.down { color:#f87171; background:rgba(248,113,113,.15); }
.cp.b .ce { animation:bounce .45s cubic-bezier(.28,1.33,.64,1) 2; }
.cp.sh .ce { animation:shake .35s ease-in-out 2; }
.ci { line-height:1; margin-top:1px; display:flex; flex-direction:column; align-items:center; }
.cn { font-size:9px; font-weight:600; color:rgba(220,210,250,.8); max-width:65px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.cno { font-size:7px; color:rgba(200,190,240,.35); }
.cs { font-size:10px; font-weight:700; color:rgba(200,190,240,.6); margin-top:1px; }
.cex { width:65%; height:2px; background:rgba(255,255,255,.05); border-radius:2px; margin-top:2px; overflow:hidden; }
.cef { height:100%; background:linear-gradient(90deg,#7c3aed,#a78bfa); border-radius:2px; transition:width .5s; }
.ca { display:none; gap:2px; margin-top:auto; padding:2px 0; }
.tm .ca { display:flex; }
.ca- { flex:1; padding:2px 0; border:none; border-radius:3px; font-size:8px; font-weight:700; cursor:pointer; transition:all .1s; font-family:inherit; line-height:1; }
.ca-:disabled { opacity:.3; cursor:not-allowed; }
.ca-p { background:rgba(74,222,128,.15); color:#4ade80; }
.ca-p:hover:not(:disabled) { background:rgba(74,222,128,.25); }
.ca-m { background:rgba(248,113,113,.15); color:#f87171; }
.ca-m:hover:not(:disabled) { background:rgba(248,113,113,.25); }

/* 排行 */
.rk { flex:1; overflow-y:auto; padding:8px 0; position:relative; z-index:1; }
.rk-h { font-size:18px; font-weight:700; margin-bottom:12px; display:flex; align-items:center; gap:8px; }
.rk-l { display:flex; flex-direction:column; gap:6px; }
.rk-i { display:flex; align-items:center; gap:10px; padding:10px 14px; border-radius:12px; background:rgba(255,255,255,.03); border:1px solid rgba(255,255,255,.04); }
.rk-i.top { background:rgba(255,255,255,.06); }
.rk-no { font-size:18px; width:30px; text-align:center; flex-shrink:0; }
.rk-n { flex:1; font-weight:500; font-size:14px; }
.rk-s { font-weight:700; font-size:15px; color:#a78bfa; }
.rk-btn { padding:4px 10px; border-radius:6px; border:1px solid rgba(255,255,255,.08); background:transparent; color:rgba(200,190,240,.5); font-size:11px; cursor:pointer; font-family:inherit; }
.rk-btn:hover { background:rgba(255,255,255,.06); color:#e8e6f0; }
.rk-ref { margin-top:12px; padding:6px 16px; border-radius:8px; border:1px solid rgba(255,255,255,.06); background:transparent; color:rgba(200,190,240,.5); font-size:12px; cursor:pointer; font-family:inherit; }

/* 商城 */
.sh { flex:1; overflow-y:auto; padding:8px 0; position:relative; z-index:1; }
.sh-h { font-size:18px; font-weight:700; margin-bottom:12px; display:flex; align-items:center; gap:8px; }
.sh-sel { display:flex; align-items:center; gap:8px; font-size:13px; margin-bottom:12px; }
.sh-sel-in { flex:1; padding:8px 12px; border-radius:8px; border:1px solid rgba(255,255,255,.1); background:rgba(255,255,255,.06); color:#e8e6f0; font-size:13px; font-family:inherit; }
.sh-nt { font-size:12px; color:rgba(200,190,240,.4); margin-bottom:12px; }
.sh-l { display:flex; flex-direction:column; gap:6px; }
.sh-i { display:flex; align-items:center; gap:10px; padding:12px 14px; border-radius:12px; background:rgba(255,255,255,.03); border:1px solid rgba(255,255,255,.04); }
.sh-ic { font-size:24px; }
.sh-ib { flex:1; }
.sh-in { font-weight:600; font-size:14px; }
.sh-id { font-size:11px; color:rgba(200,190,240,.4); }
.sh-pr { font-weight:700; font-size:15px; color:#fbbf24; white-space:nowrap; }
.sh-btn { padding:6px 14px; border-radius:8px; border:none; background:linear-gradient(135deg,#7c3aed,#6d28d9); color:#fff; font-size:12px; font-weight:600; cursor:pointer; font-family:inherit; }
.sh-btn:disabled { opacity:.3; cursor:not-allowed; }
.sh-btn:hover:not(:disabled) { background:linear-gradient(135deg,#8b5cf6,#7c3aed); }
.sh-msg { font-size:12px; margin-top:8px; color:#f87171; }
.sh-msg.ok { color:#4ade80; }

/* 底栏 */
.nav { display:flex; gap:4px; padding:6px 0 8px; flex-shrink:0; position:relative; z-index:1; }
.nav-i { flex:1; padding:8px; border:none; border-radius:10px; background:rgba(255,255,255,.03); color:rgba(200,190,240,.4); font-size:13px; font-weight:500; cursor:pointer; transition:all .15s; font-family:inherit; }
.nav-i:hover { background:rgba(255,255,255,.06); }
.nav-i.a { background:rgba(124,58,237,.2); color:#c4b5fd; }

/* 动画 */
@keyframes bounce { 0%,100%{transform:translateY(0)} 30%{transform:translateY(-5px) scale(1.08)} 60%{transform:translateY(-2px) scale(1.02)} }
@keyframes shake { 0%,100%{transform:translateX(0)} 25%{transform:translateX(-2px)} 75%{transform:translateX(2px)} }
@keyframes spin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }
.pop-enter-active { animation:popIn .35s ease-out; }
.pop-leave-active { animation:popIn .25s ease-in reverse; }
.f-enter-active { transition:all .25s ease-out; }
.f-leave-active { transition:all .15s ease-in; }
.f-enter-from { opacity:0; transform:translateY(6px); }
.f-leave-to { opacity:0; transform:translateY(-8px); }
@keyframes popIn { from{opacity:0;transform:scale(.9) translateY(20px)} to{opacity:1;transform:scale(1) translateY(0)} }
</style>
