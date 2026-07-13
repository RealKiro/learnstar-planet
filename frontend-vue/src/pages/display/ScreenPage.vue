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

interface PetEntry { student_id: number; student_no: string; student_name: string; total_score: number; has_pet: boolean; level: number; experience: number; mood: number; emoji: string; stage_name: string; exp_max: number; color: string; image?: string }
interface DisplayData { class_name: string; grade: string; student_count: number; pets: PetEntry[]; recent_scores: any[]; broadcasts: any[] }
const data = ref<DisplayData | null>(null)
const scoreAnim = ref<Record<number, { dir: 'up'|'down'; amt: number }>>({})
const scorePending = ref<Record<number, boolean>>({})

// 面板（弹出覆盖层）
const showPanel = ref<'rank'|'shop'|null>(null)

const ranking = ref<{ rank: number; id: number; name: string; score: number; no: string }[]>([])
const rankLoading = ref(false)
interface ShopItemType { id: number; name: string; description?: string; cost_score: number; stock: number; category: string }
const shopItems = ref<ShopItemType[]>([])
const shopLoading = ref(false)
const redeemTarget = ref<{ id: number; name: string } | null>(null)
const redeemMsg = ref('')
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

    <!-- 顶栏（半透明悬浮） -->
    <header class="h">
      <div class="hl">
        <h1 class="ht">{{ data?.class_name || '--' }}</h1>
        <span class="hm">{{ data?.grade }} · {{ data?.student_count }}人</span>
        <span class="sd" :class="{ o: sseState.connected }"></span>
      </div>
      <div class="hr">
        <button class="hb" @click="showPanel = 'rank'" title="排行榜">🏆</button>
        <button class="hb" @click="showPanel = 'shop'" title="积分商城">🛍️</button>
        <button class="hx" @click="confirmExit">✕ 退出</button>
      </div>
    </header>

    <!-- 全屏 5×5 萌宠网格 -->
    <main v-if="!loading && data" class="body">
      <div v-for="(s, i) in gridSlots" :key="i" class="c"
        :class="{ e: !s, su: s && scoreAnim[s.student_id]?.dir === 'up', sd: s && scoreAnim[s.student_id]?.dir === 'down' }"
        :style="s?.color ? { '--pcolor': s.color } : undefined">
        <template v-if="s">
          <div class="cp" :class="{ b: scoreAnim[s.student_id]?.dir === 'up', sh: scoreAnim[s.student_id]?.dir === 'down' }">
            <img v-if="s.has_pet && s.image" :src="s.image" class="pi" alt="">
            <span v-else class="ce">{{ s.has_pet ? s.emoji : '🥚' }}</span>
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
          <!-- 悬浮侧边快捷按钮 -->
          <button class="sb sb-l" @click.stop="qScore(s.student_id, -1)" :disabled="scorePending[s.student_id]">−</button>
          <button class="sb sb-r" @click.stop="qScore(s.student_id, 1)" :disabled="scorePending[s.student_id]">+</button>
        </template>
        <div v-else class="ce2"></div>
      </div>
    </main>

    <div v-if="loading" class="st2"><div class="sp">🌌</div><p>连接中…</p></div>
    <div v-else-if="loadError" class="st2"><p>{{ loadError }}</p><button class="sbtn" @click="loadData">重试</button></div>

    <!-- ===== 排行榜覆盖层 ===== -->
    <Transition name="fade">
      <div v-if="showPanel === 'rank'" class="ov" @click.self="showPanel = null">
        <div class="ov-bd">
          <div class="ov-h"><span>🏆</span> 排行榜<button class="ov-x" @click="showPanel = null">✕</button></div>
          <div class="ov-l">
            <div v-for="(r, i) in ranking" :key="r.id" class="ov-i" :class="{ top: i < 3 }">
              <span class="ov-no">{{ ['🥇','🥈','🥉'][i] || '#' + r.rank }}</span>
              <span class="ov-n">{{ r.name }}</span>
              <span class="ov-s">{{ r.score }}分</span>
              <button class="ov-btn" @click="transferTarget = { id: r.id, name: r.name }; transferAmt = 1; transferMsg = ''">转赠</button>
            </div>
          </div>
          <button class="ov-r" @click="loadRanking">🔄 刷新</button>
          <div v-if="transferTarget" class="ov-tr">
            <span>转赠给 <b>{{ transferTarget.name }}</b>：</span>
            <input v-model.number="transferAmt" type="number" min="1" max="100" class="ov-in">
            <button class="ov-go" @click="doTransfer">确认</button>
            <p v-if="transferMsg" class="ov-msg" :class="{ ok: transferMsg.startsWith('✅') }">{{ transferMsg }}</p>
          </div>
        </div>
      </div>
    </Transition>

    <!-- ===== 商城覆盖层 ===== -->
    <Transition name="fade">
      <div v-if="showPanel === 'shop'" class="ov" @click.self="showPanel = null">
        <div class="ov-bd">
          <div class="ov-h"><span>🛍️</span> 积分商城<button class="ov-x" @click="showPanel = null">✕</button></div>
          <div class="ov-sel">
            <select v-model="redeemTarget" class="ov-in ov-in--sel">
              <option :value="null">兑换给…</option>
              <option v-for="p in data?.pets || []" :key="p.student_id" :value="{ id: p.student_id, name: p.student_name }">{{ p.student_name }}（{{ p.total_score }}分）</option>
            </select>
          </div>
          <div class="ov-l">
            <div v-for="item in shopItems" :key="item.id" class="ov-i">
              <span class="ov-ic">{{ item.category === 'privilege' ? '👑' : item.category === 'physical' ? '🎁' : '⭐' }}</span>
              <div class="ov-ib"><div>{{ item.name }}</div><div class="ov-id">{{ item.description || '' }}</div></div>
              <span class="ov-pr">{{ item.cost_score }}分</span>
              <button class="ov-go" :disabled="!redeemTarget" @click="doRedeem(item)">兑换</button>
            </div>
          </div>
          <p v-if="redeemMsg" class="ov-msg" :class="{ ok: redeemMsg.startsWith('✅') }">{{ redeemMsg }}</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.sc {
  height: 100vh; overflow: hidden;
  background: linear-gradient(135deg,#0c0a20,#1a1040 30%,#0d1b2a 70%,#0a1628);
  color: #e8e6f0; user-select: none;
  display: flex; flex-direction: column; position: relative;
}

/* ===== 顶栏（悬浮效果） ===== */
.h {
  position: absolute; top: 0; left: 0; right: 0; z-index: 50;
  display: flex; justify-content: space-between; align-items: center;
  padding: 10px 24px;
  background: linear-gradient(to bottom, rgba(12,10,32,.8), transparent);
  pointer-events: none;
}
.hl { display: flex; align-items: center; gap: 10px; pointer-events: auto; }
.ht { font-size: 18px; font-weight: 700; margin: 0; color: #f0ecff; text-shadow: 0 2px 8px rgba(0,0,0,.5); }
.hm { font-size: 12px; color: rgba(200,190,240,.5); text-shadow: 0 2px 4px rgba(0,0,0,.5); }
.sd { width: 6px; height: 6px; border-radius: 50%; background: #64748b; }
.sd.o { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,.5); }
.hr { display: flex; align-items: center; gap: 6px; pointer-events: auto; }
.hb {
  width: 34px; height: 34px; border-radius: 8px; border: none;
  background: rgba(0,0,0,.3); backdrop-filter: blur(8px);
  color: rgba(255,255,255,.7); font-size: 16px;
  cursor: pointer; display: flex; align-items: center; justify-content: center;
}
.hb:hover { background: rgba(0,0,0,.5); }
.hx {
  padding: 6px 14px; border-radius: 8px; border: 1px solid rgba(248,113,113,.3);
  background: rgba(0,0,0,.3); backdrop-filter: blur(8px);
  color: #fca5a5; font-size: 12px; font-weight: 600;
  cursor: pointer; font-family: inherit;
}
.hx:hover { background: rgba(248,113,113,.15); }

/* ===== 全屏 5×5 网格 ===== */
.body {
  flex: 1; display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 8px; padding: 60px 16px 12px;
  align-content: center;
  min-height: 0;
}

.c {
  border-radius: 16px;
  background: rgba(255,255,255,.03);
  border: 2px solid rgba(255,255,255,.04);
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 6px; position: relative;
  transition: transform .2s cubic-bezier(.175,.885,.32,1.27), border-color .2s, background .2s;
  min-height: 0; overflow: hidden;
}
.c:hover { transform: scale(1.03) translateY(-3px); background: rgba(255,255,255,.06); border-color: rgba(124,58,237,.15); }
.c.e { background: transparent; border-style: dashed; border-color: rgba(255,255,255,.03); }
.c.e:hover { transform: none; background: transparent; }
.ce2 { width: 100%; height: 100%; }
.c.su { border-color: rgba(74,222,128,.2); }
.c.sd { border-color: rgba(248,113,113,.2); }

/* 侧边加减按钮（悬停显示，类似参考文件） */
.sb { position: absolute; top: 0; bottom: 0; width: 22px; border: none; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; opacity: 0; transition: all .2s ease; z-index: 5; color: #fff; }
.c:hover .sb { opacity: .5; }
.sb:hover { opacity: 1 !important; width: 26px; }
.sb-l { left: 0; background: linear-gradient(90deg, rgba(248,113,113,.7), transparent); border-radius: 14px 0 0 14px; }
.sb-r { right: 0; background: linear-gradient(-90deg, rgba(74,222,128,.7), transparent); border-radius: 0 14px 14px 0; }

.cp { position: relative; display: flex; align-items: center; justify-content: center; flex: 1; min-height: 0; }
.ce { font-size: min(12vw, 110px); line-height: 1; position: relative; z-index: 1; }
.pi { width: min(12vw, 110px); height: min(12vw, 110px); object-fit: contain; position: relative; z-index: 1; filter: drop-shadow(0 4px 12px rgba(0,0,0,.15)); transition: transform .3s; }
.c:hover .pi { transform: scale(1.1) rotate(-3deg); }
.cp.b .ce, .cp.b .pi { animation: bounce .45s cubic-bezier(.28,1.33,.64,1) 2; }
.cp.sh .ce, .cp.sh .pi { animation: shake .35s ease-in-out 2; }

/* 彩色光晕背景 */
.cp::before {
  content: ''; position: absolute;
  width: min(14vw, 120px); height: min(14vw, 120px);
  border-radius: 50%;
  background: radial-gradient(circle, color-mix(in srgb, var(--pcolor) 20%, transparent 70%), transparent);
  pointer-events: none; z-index: 0;
  transition: all .3s;
}
.c:hover .cp::before {
  width: min(16vw, 140px); height: min(16vw, 140px);
  background: radial-gradient(circle, color-mix(in srgb, var(--pcolor) 35%, transparent 60%), transparent);
}

.cf {
  position: absolute; top: -2px; right: -6px;
  font-size: min(2.5vw, 18px); font-weight: 700;
  padding: 1px 5px; border-radius: 5px; pointer-events: none;
}
.cf.up { color: #4ade80; background: rgba(74,222,128,.15); }
.cf.down { color: #f87171; background: rgba(248,113,113,.15); }

.ci { display: flex; align-items: center; gap: 4px; margin: 4px 0 2px; }
.cn { font-size: min(2vw, 15px); font-weight: 600; color: rgba(220,210,250,.8); max-width: 80px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cs { font-size: min(1.8vw, 13px); font-weight: 700; color: rgba(200,190,240,.5); }

.cex { width: 75%; height: 3px; background: rgba(255,255,255,.04); border-radius: 2px; margin: 2px 0 4px; overflow: hidden; }
.cef { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); transition: width .5s; }

.ca { display: flex; gap: 3px; width: 100%; margin-top: auto; }
.ca- { flex: 1; padding: 4px 0; border: none; border-radius: 5px; font-size: min(1.6vw, 13px); font-weight: 700; cursor: pointer; transition: all .1s; font-family: inherit; line-height: 1; }
.ca-:disabled { opacity: .25; cursor: not-allowed; }
.ca-p1 { background: rgba(74,222,128,.18); color: #4ade80; }
.ca-p3 { background: rgba(74,222,128,.1); color: #4ade80; }
.ca-p1:hover:not(:disabled) { background: rgba(74,222,128,.3); }
.ca-p3:hover:not(:disabled) { background: rgba(74,222,128,.2); }
.ca-m { background: rgba(248,113,113,.15); color: #f87171; }
.ca-m3 { background: rgba(248,113,113,.08); color: #f87171; }
.ca-m:hover:not(:disabled) { background: rgba(248,113,113,.25); }
.ca-m3:hover:not(:disabled) { background: rgba(248,113,113,.15); }

/* ===== 覆盖层面板 ===== */
.ov {
  position: fixed; inset: 0; z-index: 200;
  background: rgba(5,2,20,.55); backdrop-filter: blur(8px);
  display: flex; align-items: center; justify-content: center;
}
.ov-bd {
  width: 480px; max-width: 90vw; max-height: 75vh;
  background: linear-gradient(145deg,#1a1040,#0d1b2a);
  border: 1px solid rgba(255,255,255,.08); border-radius: 20px;
  padding: 24px; overflow-y: auto;
}
.ov-h { display: flex; align-items: center; gap: 8px; font-size: 18px; font-weight: 700; margin-bottom: 16px; }
.ov-x { margin-left: auto; width: 28px; height: 28px; border-radius: 7px; border: 1px solid rgba(255,255,255,.06); background: transparent; color: rgba(200,190,240,.5); cursor: pointer; font-size: 13px; }
.ov-x:hover { background: rgba(255,255,255,.06); }
.ov-l { display: flex; flex-direction: column; gap: 4px; }
.ov-i { display: flex; align-items: center; gap: 8px; padding: 8px 12px; border-radius: 10px; background: rgba(255,255,255,.025); border: 1px solid rgba(255,255,255,.03); }
.ov-i.top { background: rgba(255,255,255,.05); }
.ov-no { font-size: 16px; width: 32px; text-align: center; flex-shrink: 0; }
.ov-n { flex: 1; font-weight: 500; font-size: 13px; }
.ov-s { font-weight: 700; font-size: 13px; color: #a78bfa; }
.ov-btn { padding: 3px 8px; border-radius: 5px; border: 1px solid rgba(255,255,255,.06); background: transparent; color: rgba(200,190,240,.4); font-size: 10px; cursor: pointer; font-family: inherit; }
.ov-btn:hover { background: rgba(255,255,255,.06); color: #e8e6f0; }
.ov-r { margin-top: 10px; padding: 4px 12px; border-radius: 6px; border: 1px solid rgba(255,255,255,.04); background: transparent; color: rgba(200,190,240,.35); font-size: 11px; cursor: pointer; font-family: inherit; }
.ov-ic { font-size: 22px; }
.ov-ib { flex: 1; font-size: 13px; }
.ov-id { font-size: 10px; color: rgba(200,190,240,.35); }
.ov-pr { font-weight: 700; font-size: 14px; color: #fbbf24; white-space: nowrap; }
.ov-in { padding: 5px 8px; border-radius: 5px; border: 1px solid rgba(255,255,255,.1); background: rgba(255,255,255,.06); color: #e8e6f0; font-size: 12px; font-family: inherit; width: 50px; }
.ov-in--sel { width: 100%; }
.ov-go { padding: 5px 14px; border-radius: 6px; border: none; background: linear-gradient(135deg,#7c3aed,#6d28d9); color: #fff; font-size: 11px; font-weight: 600; cursor: pointer; font-family: inherit; }
.ov-go:disabled { opacity: .3; cursor: not-allowed; }
.ov-go:hover:not(:disabled) { background: linear-gradient(135deg,#8b5cf6,#7c3aed); }
.ov-msg { font-size: 11px; margin-top: 6px; color: #f87171; }
.ov-msg.ok { color: #4ade80; }
.ov-tr { margin-top: 8px; display: flex; align-items: center; gap: 6px; flex-wrap: wrap; font-size: 12px; }
.ov-sel { margin-bottom: 10px; }

.st2 { position: fixed; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; color: rgba(200,190,240,.5); gap: 8px; }
.sp { font-size: 48px; animation: spin 2s linear infinite; }
.sbtn { padding: 6px 16px; border-radius: 8px; border: 1px solid rgba(255,255,255,.1); background: rgba(255,255,255,.04); color: rgba(200,190,240,.7); cursor: pointer; font-size: 13px; font-family: inherit; }

/* 广播 */
.bc { position: fixed; inset: 0; z-index: 300; display: flex; align-items: center; justify-content: center; }
.bc.banner { align-items: flex-start; padding-top: 30px; background: rgba(10,5,30,.6); backdrop-filter: blur(8px); }
.bc.popup { background: rgba(10,5,30,.75); backdrop-filter: blur(12px); }
.bc.fullscreen { background: rgba(10,5,30,.92); backdrop-filter: blur(20px); }
.bc-in { text-align: center; max-width: 560px; padding: 30px; animation: popIn .35s cubic-bezier(.34,1.56,.64,1); }
.bc-ic { font-size: 44px; margin-bottom: 12px; }
.bc-txt { font-size: 24px; font-weight: 700; color: #f0ecff; line-height: 1.4; }
.bc-bar { margin: 16px auto 0; width: 180px; height: 3px; background: rgba(255,255,255,.06); border-radius: 2px; overflow: hidden; }
.bc-fill { height: 100%; background: linear-gradient(90deg,#7c3aed,#a78bfa); animation: shrink linear forwards; }

@keyframes shrink { from { width: 100% } to { width: 0% } }
@keyframes bounce { 0%,100%{transform:translateY(0)} 30%{transform:translateY(-10px) scale(1.15)} 60%{transform:translateY(-4px) scale(1.06)} }
@keyframes shake { 0%,100%{transform:translateX(0)} 25%{transform:translateX(-4px)} 75%{transform:translateX(4px)} }
@keyframes spin { from { transform: rotate(0deg) } to { transform: rotate(360deg) } }
@keyframes popIn { from { opacity: 0; transform: scale(.9) translateY(20px) } to { opacity: 1; transform: scale(1) translateY(0) } }
.pop-enter-active { animation: popIn .35s ease-out }
.pop-leave-active { animation: popIn .25s ease-in reverse }
.f-enter-active { transition: all .25s ease-out }
.f-leave-active { transition: all .15s ease-in }
.f-enter-from { opacity: 0; transform: translateY(6px) }
.f-leave-to { opacity: 0; transform: translateY(-8px) }
.fade-enter-active, .fade-leave-active { transition: opacity .2s ease }
.fade-enter-from, .fade-leave-to { opacity: 0 }
</style>
