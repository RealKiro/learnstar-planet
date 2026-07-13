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

// 面板切换（不再用底栏tab，改为弹出覆盖层）
const showPanel = ref<'rank'|'shop'|null>(null)

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
  try { const r = await apiGet<{ data: typeof ranking.value }>('/api/v1/display/leaderboard', { params: { token: token.value } }); ranking.value = r.data || [] } catch {}
  finally { rankLoading.value = false }
}
async function loadShop() {
  shopLoading.value = true
  try { const r = await apiGet<{ data: ShopItemType[] }>('/api/v1/display/shop-items', { params: { token: token.value } }); shopItems.value = r.data || [] } catch {}
  finally { shopLoading.value = false }
}

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
      <div class="hr">
        <button class="hb" @click="showPanel = 'rank'" title="排行榜">🏆</button>
        <button class="hb" @click="showPanel = 'shop'" title="积分商城">🛍️</button>
        <button class="hx" @click="confirmExit" title="退出大屏">✕ 退出</button>
      </div>
    </header>

    <!-- 加载/错误 -->
    <div v-if="loading" class="st2"><div class="sp">🌌</div><p>连接中…</p></div>
    <div v-else-if="loadError" class="st2"><p>{{ loadError }}</p><button class="sbtn" @click="loadData">重试</button></div>

    <!-- 5x5 萌宠教室 -->
    <div v-else-if="data" class="g">
      <div v-for="(s, i) in gridSlots" :key="i" class="c"
        :class="{ e: !s, su: s && scoreAnim[s.student_id]?.dir === 'up', sd: s && scoreAnim[s.student_id]?.dir === 'down' }">

        <template v-if="s">
          <!-- 大萌宠 -->
          <div class="cp" :class="{ b: scoreAnim[s.student_id]?.dir === 'up', sh: scoreAnim[s.student_id]?.dir === 'down' }">
            <span class="ce">{{ s.has_pet ? s.emoji : '🥚' }}</span>
            <Transition name="f">
              <span v-if="scoreAnim[s.student_id]" class="cf" :class="scoreAnim[s.student_id].dir">
                {{ scoreAnim[s.student_id].dir === 'up' ? '+' : '-' }}{{ scoreAnim[s.student_id].amt }}
              </span>
            </Transition>
          </div>

          <!-- 信息：姓名 + 学号 + 分数 -->
          <div class="ci">
            <span class="cn">{{ s.student_name }}</span>
            <span class="cs">{{ s.total_score }}分</span>
          </div>

          <!-- 经验条 -->
          <div class="cex"><div class="cef" :style="{ width: Math.min(100, (s.experience / Math.max(1, s.exp_max)) * 100) + '%' }"></div></div>

          <!-- 快捷操作（默认显示） -->
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

    <!-- ===== 排行榜覆盖层 ===== -->
    <Transition name="fade">
      <div v-if="showPanel === 'rank'" class="modal" @click.self="showPanel = null">
        <div class="modal-bd">
          <div class="modal-h"><span>🏆</span> 排行榜<button class="modal-x" @click="showPanel = null">✕</button></div>
          <div v-if="rankLoading" class="st2"><p>加载中…</p></div>
          <div v-else class="modal-l">
            <div v-for="(r, i) in ranking" :key="r.id" class="modal-i" :class="{ top: i < 3 }">
              <span class="modal-no">{{ ['🥇','🥈','🥉'][i] || '#' + r.rank }}</span>
              <span class="modal-n">{{ r.name }}</span>
              <span class="modal-s">{{ r.score }}分</span>
              <button class="modal-btn" @click="transferTarget = { id: r.id, name: r.name }; transferAmt = 1; transferMsg = ''">转赠</button>
            </div>
            <button class="modal-r" @click="loadRanking">🔄 刷新</button>
          </div>
          <!-- 转赠 -->
          <div v-if="transferTarget" class="modal-tr">
            <span>转赠给 <b>{{ transferTarget.name }}</b>：</span>
            <input v-model.number="transferAmt" type="number" min="1" max="100" class="modal-in">
            <button class="modal-go" @click="doTransfer">确认</button>
            <p v-if="transferMsg">{{ transferMsg }}</p>
          </div>
        </div>
      </div>
    </Transition>

    <!-- ===== 商城覆盖层 ===== -->
    <Transition name="fade">
      <div v-if="showPanel === 'shop'" class="modal" @click.self="showPanel = null">
        <div class="modal-bd">
          <div class="modal-h"><span>🛍️</span> 积分商城<button class="modal-x" @click="showPanel = null">✕</button></div>
          <div class="modal-sel">
            <span>兑换学生：</span>
            <select v-model="redeemTarget" class="modal-in">
              <option :value="null">请选择</option>
              <option v-for="p in data?.pets || []" :key="p.student_id" :value="{ id: p.student_id, name: p.student_name }">{{ p.student_name }}（{{ p.total_score }}分）</option>
            </select>
          </div>
          <div v-if="shopLoading" class="st2"><p>加载中…</p></div>
          <div v-else class="modal-l">
            <div v-for="item in shopItems" :key="item.id" class="modal-i">
              <span class="modal-ic">{{ item.category === 'privilege' ? '👑' : item.category === 'physical' ? '🎁' : '⭐' }}</span>
              <div class="modal-ib"><div class="modal-in-l">{{ item.name }}</div><div class="modal-id">{{ item.description || item.category }}</div></div>
              <span class="modal-pr">{{ item.cost_score }}分</span>
              <button class="modal-go" :disabled="!redeemTarget" @click="doRedeem(item)">兑换</button>
            </div>
          </div>
          <p v-if="redeemMsg" class="modal-msg" :class="{ ok: redeemMsg.startsWith('✅') }">{{ redeemMsg }}</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.sc {
  min-height: 100vh;
  background: linear-gradient(135deg,#0c0a20,#1a1040 30%,#0d1b2a 70%,#0a1628);
  color: #e8e6f0;
  font-family: "PingFang SC","Noto Sans SC",-apple-system,sans-serif;
  padding: 16px 24px 20px;
  user-select: none; position: relative;
  display: flex; flex-direction: column;
}

/* 顶栏 */
.h {
  display: flex; justify-content: space-between; align-items: center;
  padding-bottom: 10px; border-bottom: 1px solid rgba(255,255,255,.06); margin-bottom: 10px;
  flex-shrink: 0;
}
.hl { display: flex; align-items: center; gap: 12px; }
.ht { font-size: 22px; font-weight: 700; margin: 0; color: #f0ecff; }
.hm { font-size: 13px; color: rgba(200,190,240,.5); }
.sd { width: 7px; height: 7px; border-radius: 50%; background: #64748b; transition: all .3s; }
.sd.o { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,.5); }
.hr { display: flex; align-items: center; gap: 8px; }
.hb {
  width: 38px; height: 38px; border-radius: 10px; border: 1px solid rgba(255,255,255,.08);
  background: rgba(255,255,255,.04); color: rgba(200,190,240,.6); font-size: 18px;
  cursor: pointer; display: flex; align-items: center; justify-content: center;
  transition: all .15s;
}
.hb:hover { background: rgba(255,255,255,.1); color: #e8e6f0; }
.hx {
  padding: 8px 18px; border-radius: 10px; border: 1px solid rgba(248,113,113,.3);
  background: rgba(248,113,113,.1); color: #fca5a5; font-size: 14px; font-weight: 600;
  cursor: pointer; font-family: inherit; transition: all .15s;
}
.hx:hover { background: rgba(248,113,113,.2); color: #fff; }

.st2 { text-align:center; padding:40px 20px; color:rgba(200,190,240,.5); flex:1; display:flex; flex-direction:column; align-items:center; justify-content:center; }
.sp { font-size:48px; animation:spin 2s linear infinite; }
.sbtn { padding:6px 16px; border-radius:8px; border:1px solid rgba(255,255,255,.1); background:rgba(255,255,255,.04); color:rgba(200,190,240,.7); cursor:pointer; font-size:13px; font-family:inherit; }

/* 5x5 网格 */
.g { flex:1; display:grid; grid-template-columns:repeat(5,1fr); gap:10px; max-width:1000px; margin:0 auto; position:relative; z-index:1; overflow-y:auto; padding-bottom:4px; }

.c {
  border-radius: 16px;
  background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.06);
  display: flex; flex-direction: column; align-items: center;
  padding: 12px 6px 8px; position: relative; overflow: hidden;
  transition: transform .15s, border-color .15s; min-width: 0;
}
.c.e { background: transparent; border-style: dashed; border-color: rgba(255,255,255,.03); }
.ce2 { width: 100%; height: 100%; }
.c.su { border-color: rgba(74,222,128,.2); box-shadow: inset 0 0 30px rgba(74,222,128,.04); }
.c.sd { border-color: rgba(248,113,113,.2); box-shadow: inset 0 0 30px rgba(248,113,113,.04); }

/* 大萌宠 */
.cp { position: relative; display: flex; align-items: center; justify-content: center; margin: 4px 0; }
.ce { font-size: 56px; line-height: 1; filter: drop-shadow(0 0 8px rgba(180,140,255,.25)); transition: all .3s; }
.cp.b .ce { animation: bounce .45s cubic-bezier(.28,1.33,.64,1) 2; }
.cp.sh .ce { animation: shake .35s ease-in-out 2; }

.cf {
  position: absolute; top: -4px; right: -8px; font-size: 14px; font-weight: 700;
  padding: 2px 6px; border-radius: 6px; pointer-events: none;
}
.cf.up { color: #4ade80; background: rgba(74,222,128,.15); }
.cf.down { color: #f87171; background: rgba(248,113,113,.15); }

/* 信息 */
.ci { display: flex; align-items: center; gap: 6px; margin: 4px 0 2px; }
.cn { font-size: 14px; font-weight: 600; color: rgba(220,210,250,.85); max-width: 80px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cs { font-size: 13px; font-weight: 700; color: rgba(200,190,240,.55); }

/* 经验条 */
.cex { width: 75%; height: 3px; background: rgba(255,255,255,.06); border-radius: 2px; margin: 2px 0 6px; overflow: hidden; }
.cef { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); border-radius: 2px; transition: width .5s; }

/* 快捷操作 */
.ca { display: flex; gap: 4px; width: 100%; }
.ca- { flex: 1; padding: 5px 0; border: none; border-radius: 6px; font-size: 13px; font-weight: 700; cursor: pointer; transition: all .1s; font-family: inherit; line-height: 1; }
.ca-:disabled { opacity: .3; cursor: not-allowed; }
.ca-p1 { background: rgba(74,222,128,.18); color: #4ade80; }
.ca-p3 { background: rgba(74,222,128,.12); color: #4ade80; }
.ca-p1:hover:not(:disabled) { background: rgba(74,222,128,.3); }
.ca-p3:hover:not(:disabled) { background: rgba(74,222,128,.2); }
.ca-m { background: rgba(248,113,113,.15); color: #f87171; }
.ca-m3 { background: rgba(248,113,113,.1); color: #f87171; }
.ca-m:hover:not(:disabled) { background: rgba(248,113,113,.25); }
.ca-m3:hover:not(:disabled) { background: rgba(248,113,113,.18); }

/* 覆盖层面板 */
.modal {
  position: fixed; inset: 0; z-index: 200;
  background: rgba(5,2,20,.7); backdrop-filter: blur(12px);
  display: flex; align-items: center; justify-content: center;
}
.modal-bd {
  width: 520px; max-width: 92vw; max-height: 80vh;
  background: linear-gradient(145deg,#1a1040,#0d1b2a);
  border: 1px solid rgba(255,255,255,.08); border-radius: 24px;
  padding: 28px; overflow-y: auto;
}
.modal-h { display: flex; align-items: center; gap: 8px; font-size: 20px; font-weight: 700; margin-bottom: 20px; }
.modal-x { margin-left: auto; width: 32px; height: 32px; border-radius: 8px; border: 1px solid rgba(255,255,255,.08); background: transparent; color: rgba(200,190,240,.5); cursor: pointer; font-size: 14px; }
.modal-x:hover { background: rgba(255,255,255,.06); color: #e8e6f0; }
.modal-l { display: flex; flex-direction: column; gap: 6px; }
.modal-i { display: flex; align-items: center; gap: 10px; padding: 10px 14px; border-radius: 12px; background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.04); }
.modal-i.top { background: rgba(255,255,255,.06); }
.modal-no { font-size: 18px; width: 36px; text-align: center; flex-shrink: 0; }
.modal-n { flex: 1; font-weight: 500; font-size: 14px; }
.modal-s { font-weight: 700; font-size: 15px; color: #a78bfa; }
.modal-btn { padding: 4px 10px; border-radius: 6px; border:1px solid rgba(255,255,255,.08); background:transparent; color:rgba(200,190,240,.5); font-size:11px; cursor:pointer; font-family:inherit; }
.modal-btn:hover { background:rgba(255,255,255,.06); color:#e8e6f0; }
.modal-r { margin-top:12px; padding:6px 16px; border-radius:8px; border:1px solid rgba(255,255,255,.06); background:transparent; color:rgba(200,190,240,.5); font-size:12px; cursor:pointer; font-family:inherit; }
.modal-ic { font-size: 24px; }
.modal-ib { flex: 1; }
.modal-id { font-size: 11px; color: rgba(200,190,240,.4); }
.modal-pr { font-weight: 700; font-size: 15px; color: #fbbf24; white-space: nowrap; }
.modal-in { padding:6px 10px; border-radius:6px; border:1px solid rgba(255,255,255,.1); background:rgba(255,255,255,.06); color:#e8e6f0; font-size:13px; font-family:inherit; width:60px; }
select.modal-in { width: auto; min-width: 120px; }
.modal-go { padding:6px 16px; border-radius:8px; border:none; background:linear-gradient(135deg,#7c3aed,#6d28d9); color:#fff; font-size:12px; font-weight:600; cursor:pointer; font-family:inherit; }
.modal-go:disabled { opacity:.3; cursor:not-allowed; }
.modal-go:hover:not(:disabled) { background:linear-gradient(135deg,#8b5cf6,#7c3aed); }
.modal-msg { font-size:12px; margin-top:8px; color:#f87171; }
.modal-msg.ok { color:#4ade80; }
.modal-tr { margin-top: 12px; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; font-size: 13px; }
.modal-sel { display: flex; align-items: center; gap: 8px; font-size: 13px; margin-bottom: 12px; }

/* 广播 */
.bc { position:fixed; inset:0; z-index:100; display:flex; align-items:center; justify-content:center; }
.bc.banner { align-items:flex-start; padding-top:30px; background:rgba(10,5,30,.6); backdrop-filter:blur(8px); }
.bc.popup { background:rgba(10,5,30,.75); backdrop-filter:blur(12px); }
.bc.fullscreen { background:rgba(10,5,30,.92); backdrop-filter:blur(20px); }
.bc-in { text-align:center; max-width:560px; padding:30px; animation:popIn .35s cubic-bezier(.34,1.56,.64,1); }
.bc-ic { font-size:44px; margin-bottom:12px; }
.bc-txt { font-size:24px; font-weight:700; color:#f0ecff; line-height:1.4; }
.bc-bar { margin:16px auto 0; width:180px; height:3px; background:rgba(255,255,255,.06); border-radius:2px; overflow:hidden; }
.bc-fill { height:100%; background:linear-gradient(90deg,#7c3aed,#a78bfa); border-radius:2px; animation:shrink linear forwards; }

@keyframes shrink { from{width:100%} to{width:0%} }
@keyframes bounce { 0%,100%{transform:translateY(0)} 30%{transform:translateY(-8px) scale(1.12)} 60%{transform:translateY(-3px) scale(1.05)} }
@keyframes shake { 0%,100%{transform:translateX(0)} 25%{transform:translateX(-3px)} 75%{transform:translateX(3px)} }
@keyframes spin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }
@keyframes popIn { from{opacity:0;transform:scale(.9) translateY(20px)} to{opacity:1;transform:scale(1) translateY(0)} }
.pop-enter-active { animation:popIn .35s ease-out; }
.pop-leave-active { animation:popIn .25s ease-in reverse; }
.f-enter-active { transition:all .25s ease-out; }
.f-leave-active { transition:all .15s ease-in; }
.f-enter-from { opacity:0; transform:translateY(6px); }
.f-leave-to { opacity:0; transform:translateY(-8px); }
.fade-enter-active, .fade-leave-active { transition: opacity .2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
