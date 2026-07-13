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

interface PetEntry { student_id: number; student_no: string; student_name: string; total_score: number; has_pet: boolean; level: number; experience: number; mood: number; emoji: string; stage_name: string; exp_max: number }
interface DisplayData { class_name: string; grade: string; student_count: number; pets: PetEntry[]; recent_scores: any[]; broadcasts: any[] }
const data = ref<DisplayData | null>(null)
const scoreAnim = ref<Record<number, { dir: 'up'|'down'; amt: number }>>({})
const scorePending = ref<Record<number, boolean>>({})

// 左侧面板 tab
const sideTab = ref<'rank'|'shop'>('rank')

// 排行
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

const gridSlots = computed(() => {
  if (!data.value) return []
  const s: (PetEntry | null)[] = [...(data.value.pets || [])]
  while (s.length < 25) s.push(null)
  return s.slice(0, 25)
})

const activeBc = ref<BroadcastData | null>(null)
let bcTimer: ReturnType<typeof setTimeout> | null = null

const { state: sseState, scoreUpdates, broadcasts, petUpdates, connect: connectSSE } = useDisplaySSE()

onMounted(() => {
  window.addEventListener('display:token-expired', () => router.replace({ name: 'display-login' }))
  const t = sessionStorage.getItem('display_token')
  const ci = sessionStorage.getItem('display_class_info')
  if (!t || !ci) { router.replace({ name: 'display-login' }); return }
  token.value = t; classInfo.value = JSON.parse(ci)
  loadData(); loadRanking(); loadShop()
})

async function loadData() {
  loading.value = true
  try { const r = await apiGet<{ data: DisplayData }>('/api/v1/display/initial-data', { params: { token: token.value } }); data.value = r.data; loading.value = false; connectSSE(token.value, classInfo.value!.id) }
  catch (e: any) { if (e?.response?.status === 401) { sessionStorage.clear(); router.replace({ name: 'display-login' }); return }; loadError.value = '加载失败'; loading.value = false }
}
async function loadRanking() { rankLoading.value = true; try { const r = await apiGet<{ data: typeof ranking.value }>('/api/v1/display/leaderboard', { params: { token: token.value } }); ranking.value = r.data || [] } catch {} finally { rankLoading.value = false } }
async function loadShop() { shopLoading.value = true; try { const r = await apiGet<{ data: ShopItemType[] }>('/api/v1/display/shop-items', { params: { token: token.value } }); shopItems.value = r.data || [] } catch {} finally { shopLoading.value = false } }

async function qScore(sid: number, amt: number) {
  if (scorePending.value[sid]) return
  scorePending.value[sid] = true
  try { await apiPost('/api/v1/display/quick-score', { token: token.value, student_id: sid, amount: amt }); const p = data.value?.pets.find(x => x.student_id === sid); if (p) p.total_score += amt; anim(sid, amt) } catch {}
  finally { scorePending.value[sid] = false }
}
function anim(sid: number, amt: number) { scoreAnim.value[sid] = { dir: amt > 0 ? 'up' : 'down', amt: Math.abs(amt) }; setTimeout(() => { delete scoreAnim.value[sid] }, 1800) }

async function doRedeem(item: ShopItemType) {
  if (!redeemTarget.value) return; redeemMsg.value = ''
  try { const r = await apiPost('/api/v1/display/redeem', { token: token.value, student_id: redeemTarget.value.id, item_id: item.id }); const d = (r as any).data; redeemMsg.value = `✅ ${d.student_name} 兑换「${d.item_name}」`; const p = data.value?.pets.find(x => x.student_id === redeemTarget.value!.id); if (p) p.total_score = d.total_score; loadRanking() }
  catch (e: any) { redeemMsg.value = '❌ ' + (e?.response?.data?.message || '兑换失败') }
}
async function doTransfer() {
  if (!transferTarget.value) return; transferMsg.value = ''
  try { await apiPost('/api/v1/display/transfer', { token: token.value, from_id: transferTarget.value.id, to_id: classInfo.value!.id, amount: transferAmt.value }); transferMsg.value = '✅ 转赠成功'; loadRanking(); loadData() }
  catch (e: any) { transferMsg.value = '❌ ' + (e?.response?.data?.message || '转赠失败') }
}
function confirmExit() { sessionStorage.clear(); router.replace({ name: 'display-login' }) }

watch(scoreUpdates, (evts) => { for (const e of evts) { anim(e.student_id, e.amount); const p = data.value?.pets.find(x => x.student_id === e.student_id); if (p) { p.total_score = e.total_score; if (e.pet_level !== undefined) p.level = e.pet_level; if (e.pet_experience !== undefined) p.experience = e.pet_experience } } }, { deep: true })
watch(broadcasts, (evts) => { const m = evts[evts.length - 1]; if (m) { if (bcTimer) clearTimeout(bcTimer); activeBc.value = m; bcTimer = setTimeout(() => { activeBc.value = null }, Math.max(3000, (m.display_seconds || 8) * 1000)) } }, { deep: true })
</script>

<template>
  <div class="sc">
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

    <!-- 顶栏 -->
    <header class="h">
      <div class="hl">
        <h1 class="ht">{{ data?.class_name || '--' }}</h1>
        <span class="hm">{{ data?.grade }} · {{ data?.student_count }}人</span>
        <span class="sd" :class="{ o: sseState.connected }"></span>
      </div>
      <button class="hx" @click="confirmExit">✕ 退出大屏</button>
    </header>

    <!-- 主体：左右布局 -->
    <main v-if="!loading && data" class="body">

      <!-- 左侧面板：排行榜 / 商城 -->
      <aside class="panel">
        <div class="panel-tabs">
          <button class="pt" :class="{ a: sideTab === 'rank' }" @click="sideTab = 'rank'">🏆 排行榜</button>
          <button class="pt" :class="{ a: sideTab === 'shop' }" @click="sideTab = 'shop'">🛍️ 积分商城</button>
        </div>

        <!-- 排行榜 -->
        <div v-if="sideTab === 'rank'" class="panel-body">
          <div v-if="rankLoading" class="st2"><p>加载中…</p></div>
          <div v-else class="rank-l">
            <div v-for="(r, i) in ranking" :key="r.id" class="rank-i" :class="{ top: i < 3 }">
              <span class="rank-no">{{ ['🥇','🥈','🥉'][i] || '# ' + r.rank }}</span>
              <span class="rank-n">{{ r.name }}</span>
              <span class="rank-s">{{ r.score }}</span>
              <button class="rank-btn" @click="transferTarget = { id: r.id, name: r.name }; transferAmt = 1; transferMsg = ''">转赠</button>
            </div>
            <button class="panel-r" @click="loadRanking">🔄 刷新</button>
          </div>
          <div v-if="transferTarget" class="panel-tr">
            <span>转赠给 <b>{{ transferTarget.name }}</b></span>
            <div class="panel-tr-row">
              <input v-model.number="transferAmt" type="number" min="1" max="100" class="panel-in">
              <button class="panel-go" @click="doTransfer">确认</button>
            </div>
            <p v-if="transferMsg" class="panel-msg" :class="{ ok: transferMsg.startsWith('✅') }">{{ transferMsg }}</p>
          </div>
        </div>

        <!-- 商城 -->
        <div v-if="sideTab === 'shop'" class="panel-body">
          <div class="panel-sel">
            <select v-model="redeemTarget" class="panel-in panel-in--sel">
              <option :value="null">兑换给…</option>
              <option v-for="p in data?.pets || []" :key="p.student_id" :value="{ id: p.student_id, name: p.student_name }">{{ p.student_name }}（{{ p.total_score }}分）</option>
            </select>
          </div>
          <div v-if="shopLoading" class="st2"><p>加载中…</p></div>
          <div v-else class="rank-l">
            <div v-for="item in shopItems" :key="item.id" class="rank-i">
              <span class="rank-no">{{ item.category === 'privilege' ? '👑' : item.category === 'physical' ? '🎁' : '⭐' }}</span>
              <div class="rank-n"><div class="rank-inn">{{ item.name }}</div><div class="rank-id">{{ item.description || '' }}</div></div>
              <span class="rank-s rank-s--g">{{ item.cost_score }}</span>
              <button class="rank-btn rank-btn--go" :disabled="!redeemTarget" @click="doRedeem(item)">兑换</button>
            </div>
            <button class="panel-r" @click="loadShop">🔄 刷新</button>
          </div>
          <p v-if="redeemMsg" class="panel-msg" :class="{ ok: redeemMsg.startsWith('✅') }">{{ redeemMsg }}</p>
        </div>
      </aside>

      <!-- 右侧：5×5 萌宠教室 -->
      <div class="grid">
        <div v-for="(s, i) in gridSlots" :key="i" class="c"
          :class="{ e: !s, su: s && scoreAnim[s.student_id]?.dir === 'up', sd: s && scoreAnim[s.student_id]?.dir === 'down' }">
          <template v-if="s">
            <div class="cp" :class="{ b: scoreAnim[s.student_id]?.dir === 'up', sh: scoreAnim[s.student_id]?.dir === 'down' }">
              <span class="ce">{{ s.has_pet ? s.emoji : '🥚' }}</span>
              <Transition name="f">
                <span v-if="scoreAnim[s.student_id]" class="cf" :class="scoreAnim[s.student_id].dir">
                  {{ scoreAnim[s.student_id].dir === 'up' ? '+' : '-' }}{{ scoreAnim[s.student_id].amt }}
                </span>
              </Transition>
            </div>
            <div class="ci"><span class="cn">{{ s.student_name }}</span><span class="cs">{{ s.total_score }}分</span></div>
            <div class="cex"><div class="cef" :style="{ width: Math.min(100, (s.experience / Math.max(1, s.exp_max)) * 100) + '%' }"></div></div>
            <div class="ca">
              <button class="ca- ca-m" @click.stop="qScore(s.student_id, -1)" :disabled="scorePending[s.student_id]">−1</button>
              <button class="ca- ca-m3" @click.stop="qScore(s.student_id, -3)" :disabled="scorePending[s.student_id]">−3</button>
              <button class="ca- ca-p1" @click.stop="qScore(s.student_id, 1)" :disabled="scorePending[s.student_id]">+1</button>
              <button class="ca- ca-p3" @click.stop="qScore(s.student_id, 3)" :disabled="scorePending[s.student_id]">+3</button>
            </div>
          </template>
          <div v-else class="ce2"></div>
        </div>
      </div>
    </main>

    <div v-if="loading" class="st2 st2--f"><div class="sp">🌌</div><p>连接中…</p></div>
    <div v-else-if="loadError" class="st2 st2--f"><p>{{ loadError }}</p><button class="sbtn" @click="loadData">重试</button></div>
  </div>
</template>

<style scoped>
.sc {
  height: 100vh; overflow: hidden;
  background: linear-gradient(135deg,#0c0a20,#1a1040 30%,#0d1b2a 70%,#0a1628);
  color: #e8e6f0;
  font-family: "PingFang SC","Noto Sans SC",-apple-system,sans-serif;
  user-select: none; display: flex; flex-direction: column;
}

/* 顶栏 */
.h {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14px 28px; border-bottom: 1px solid rgba(255,255,255,.05); flex-shrink: 0;
}
.hl { display: flex; align-items: center; gap: 12px; }
.ht { font-size: 22px; font-weight: 700; margin: 0; color: #f0ecff; }
.hm { font-size: 13px; color: rgba(200,190,240,.5); }
.sd { width: 7px; height: 7px; border-radius: 50%; background: #64748b; }
.sd.o { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,.5); }
.hx {
  padding: 8px 20px; border-radius: 10px; border: 1px solid rgba(248,113,113,.3);
  background: rgba(248,113,113,.1); color: #fca5a5; font-size: 14px; font-weight: 600;
  cursor: pointer; font-family: inherit; transition: all .15s;
}
.hx:hover { background: rgba(248,113,113,.2); color: #fff; }

/* 主体 */
.body { flex: 1; display: flex; gap: 0; min-height: 0; }

/* ===== 左侧面板 ===== */
.panel {
  width: 320px; flex-shrink: 0;
  display: flex; flex-direction: column;
  border-right: 1px solid rgba(255,255,255,.04);
  padding: 16px 16px 12px;
}
.panel-tabs { display: flex; gap: 4px; margin-bottom: 12px; }
.pt {
  flex: 1; padding: 8px; border: none; border-radius: 8px;
  background: rgba(255,255,255,.03); color: rgba(200,190,240,.4);
  font-size: 13px; font-weight: 500; cursor: pointer; font-family: inherit;
}
.pt.a { background: rgba(124,58,237,.2); color: #c4b5fd; }
.pt:hover { background: rgba(255,255,255,.06); }

.panel-body { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 0; }
.panel-body::-webkit-scrollbar { width: 4px; }
.panel-body::-webkit-scrollbar-thumb { background: rgba(255,255,255,.08); border-radius: 2px; }

.rank-l { display: flex; flex-direction: column; gap: 4px; }
.rank-i { display: flex; align-items: center; gap: 6px; padding: 8px 10px; border-radius: 10px; background: rgba(255,255,255,.025); border: 1px solid rgba(255,255,255,.03); }
.rank-i.top { background: rgba(255,255,255,.05); }
.rank-no { font-size: 16px; width: 32px; text-align: center; flex-shrink: 0; }
.rank-n { flex: 1; font-weight: 500; font-size: 13px; overflow: hidden; }
.rank-inn { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rank-id { font-size: 10px; color: rgba(200,190,240,.35); }
.rank-s { font-weight: 700; font-size: 13px; color: #a78bfa; }
.rank-s--g { color: #fbbf24; }
.rank-btn {
  padding: 3px 8px; border-radius: 5px; border: 1px solid rgba(255,255,255,.06);
  background: transparent; color: rgba(200,190,240,.4); font-size: 10px; cursor: pointer; font-family: inherit; white-space: nowrap;
}
.rank-btn:hover { background: rgba(255,255,255,.06); color: #e8e6f0; }
.rank-btn--go { border: none; background: linear-gradient(135deg,#7c3aed,#6d28d9); color: #fff; font-weight: 600; }
.rank-btn--go:disabled { opacity: .3; cursor: not-allowed; }
.rank-btn--go:hover:not(:disabled) { background: linear-gradient(135deg,#8b5cf6,#7c3aed); }
.panel-r { margin-top: 8px; padding: 4px 12px; border-radius: 6px; border: 1px solid rgba(255,255,255,.04); background: transparent; color: rgba(200,190,240,.35); font-size: 11px; cursor: pointer; font-family: inherit; }
.panel-r:hover { background: rgba(255,255,255,.04); }

.panel-sel { margin-bottom: 10px; }
.panel-in { padding: 6px 10px; border-radius: 6px; border: 1px solid rgba(255,255,255,.1); background: rgba(255,255,255,.06); color: #e8e6f0; font-size: 12px; font-family: inherit; }
.panel-in--sel { width: 100%; }
.panel-tr { margin-top: 8px; padding: 8px 10px; border-radius: 8px; background: rgba(255,255,255,.03); font-size: 12px; }
.panel-tr-row { display: flex; gap: 6px; margin-top: 6px; }
.panel-tr-row .panel-in { width: 60px; }
.panel-go { padding: 4px 12px; border-radius: 6px; border: none; background: linear-gradient(135deg,#7c3aed,#6d28d9); color: #fff; font-size: 11px; font-weight: 600; cursor: pointer; font-family: inherit; }
.panel-msg { font-size: 11px; margin: 4px 0 0; color: #f87171; }
.panel-msg.ok { color: #4ade80; }

/* ===== 右侧：5×5 萌宠网格 ===== */
.grid {
  flex: 1; display: grid; grid-template-columns: repeat(5, 1fr); gap: 12px;
  padding: 24px 32px; align-content: center;
}

.c {
  border-radius: 18px;
  background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.06);
  display: flex; flex-direction: column; align-items: center;
  padding: 16px 8px 10px; position: relative;
  transition: transform .15s, border-color .15s; min-width: 0;
}
.c.e { background: transparent; border-style: dashed; border-color: rgba(255,255,255,.03); }
.ce2 { width: 100%; height: 100%; }
.c.su { border-color: rgba(74,222,128,.2); }
.c.sd { border-color: rgba(248,113,113,.2); }

.cp { position: relative; display: flex; align-items: center; justify-content: center; margin: 6px 0; }
.ce { font-size: 72px; line-height: 1; filter: drop-shadow(0 0 10px rgba(180,140,255,.25)); }
.cp.b .ce { animation: bounce .45s cubic-bezier(.28,1.33,.64,1) 2; }
.cp.sh .ce { animation: shake .35s ease-in-out 2; }

.cf { position: absolute; top: 0; right: -4px; font-size: 16px; font-weight: 700; padding: 2px 6px; border-radius: 6px; pointer-events: none; }
.cf.up { color: #4ade80; background: rgba(74,222,128,.15); }
.cf.down { color: #f87171; background: rgba(248,113,113,.15); }

.ci { display: flex; align-items: center; gap: 6px; margin: 6px 0 2px; }
.cn { font-size: 15px; font-weight: 600; color: rgba(220,210,250,.85); max-width: 80px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cs { font-size: 14px; font-weight: 700; color: rgba(200,190,240,.55); }

.cex { width: 80%; height: 3px; background: rgba(255,255,255,.06); border-radius: 2px; margin: 2px 0 8px; overflow: hidden; }
.cef { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); transition: width .5s; }

.ca { display: flex; gap: 4px; width: 100%; }
.ca- { flex: 1; padding: 6px 0; border: none; border-radius: 6px; font-size: 14px; font-weight: 700; cursor: pointer; transition: all .1s; font-family: inherit; line-height: 1; }
.ca-:disabled { opacity: .3; cursor: not-allowed; }
.ca-p1 { background: rgba(74,222,128,.18); color: #4ade80; }
.ca-p3 { background: rgba(74,222,128,.12); color: #4ade80; }
.ca-p1:hover:not(:disabled) { background: rgba(74,222,128,.3); }
.ca-p3:hover:not(:disabled) { background: rgba(74,222,128,.2); }
.ca-m { background: rgba(248,113,113,.15); color: #f87171; }
.ca-m3 { background: rgba(248,113,113,.1); color: #f87171; }
.ca-m:hover:not(:disabled) { background: rgba(248,113,113,.25); }
.ca-m3:hover:not(:disabled) { background: rgba(248,113,113,.18); }

/* 加载 */
.st2 { text-align: center; color: rgba(200,190,240,.5); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; }
.st2--f { flex: 1; }
.sp { font-size: 48px; animation: spin 2s linear infinite; }
.sbtn { padding: 6px 16px; border-radius: 8px; border: 1px solid rgba(255,255,255,.1); background: rgba(255,255,255,.04); color: rgba(200,190,240,.7); cursor: pointer; font-size: 13px; font-family: inherit; }

/* 广播 */
.bc { position: fixed; inset: 0; z-index: 200; display: flex; align-items: center; justify-content: center; }
.bc.banner { align-items: flex-start; padding-top: 30px; background: rgba(10,5,30,.6); backdrop-filter: blur(8px); }
.bc.popup { background: rgba(10,5,30,.75); backdrop-filter: blur(12px); }
.bc.fullscreen { background: rgba(10,5,30,.92); backdrop-filter: blur(20px); }
.bc-in { text-align: center; max-width: 560px; padding: 30px; animation: popIn .35s cubic-bezier(.34,1.56,.64,1); }
.bc-ic { font-size: 44px; margin-bottom: 12px; }
.bc-txt { font-size: 24px; font-weight: 700; color: #f0ecff; line-height: 1.4; }
.bc-bar { margin: 16px auto 0; width: 180px; height: 3px; background: rgba(255,255,255,.06); border-radius: 2px; overflow: hidden; }
.bc-fill { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); animation: shrink linear forwards; }

@keyframes shrink { from { width: 100% } to { width: 0% } }
@keyframes bounce { 0%,100%{transform:translateY(0)} 30%{transform:translateY(-8px) scale(1.12)} 60%{transform:translateY(-3px) scale(1.05)} }
@keyframes shake { 0%,100%{transform:translateX(0)} 25%{transform:translateX(-3px)} 75%{transform:translateX(3px)} }
@keyframes spin { from { transform: rotate(0deg) } to { transform: rotate(360deg) } }
@keyframes popIn { from { opacity: 0; transform: scale(.9) translateY(20px) } to { opacity: 1; transform: scale(1) translateY(0) } }
.pop-enter-active { animation: popIn .35s ease-out }
.pop-leave-active { animation: popIn .25s ease-in reverse }
.f-enter-active { transition: all .25s ease-out }
.f-leave-active { transition: all .15s ease-in }
.f-enter-from { opacity: 0; transform: translateY(6px) }
.f-leave-to { opacity: 0; transform: translateY(-8px) }
</style>
