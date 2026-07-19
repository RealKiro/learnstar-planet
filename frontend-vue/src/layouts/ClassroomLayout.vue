<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { apiPost } from '@/utils/api'
import { getAllSeries, getSeriesName } from '@/utils/petData'

const router = useRouter()
const route = useRoute()

const classInfo = ref<{ id: number; name: string; student_count?: number } | null>(null)
const showVoteModal = ref(true)
const voteSeries = ref('myth')
const voting = ref(false)
const voteDone = ref(false)

const activeNav = computed(() => String(route.name))

const navItems = [
  { page: 'classroom-overview', label: '班级总览', icon: '🏠' },
  { page: 'classroom-scores', label: '课堂评价', icon: '✏️' },
  { page: 'classroom-pk', label: '年级战场', icon: '🏆' },
  { page: 'classroom-pokedex', label: '宠物图鉴', icon: '📚' },
]

const allSeries = getAllSeries()

function navigate(name: string) { router.push({ name }) }
function goToLogin() { sessionStorage.clear(); router.push({ name: 'login', query: { mode: 'code' } }) }

async function confirmVote() {
  voting.value = true
  const token = sessionStorage.getItem('class_token') || ''
  try {
    await apiPost('/api/v1/display/switch-series', { token, series_id: voteSeries.value })
    voteDone.value = true
    sessionStorage.setItem('class_series', voteSeries.value)
    setTimeout(() => { showVoteModal.value = false }, 2000)
  } catch { /* ignore */ } finally { voting.value = false }
}

onMounted(() => {
  const ci = sessionStorage.getItem('class_info')
  if (ci) classInfo.value = JSON.parse(ci)
  // 如果已有系列配置则不弹投票
  if (sessionStorage.getItem('class_series')) {
    showVoteModal.value = false
  }
})
</script>

<template>
  <div class="app-shell">
    <!-- 首次使用投票弹窗 -->
    <Transition name="fade">
      <div v-if="showVoteModal" @click.self="() => {}"
        style="position:fixed;inset:0;z-index:999;background:rgba(5,2,20,0.9);backdrop-filter:blur(20px);display:flex;align-items:center;justify-content:center;padding:20px;">
        <div style="background:linear-gradient(180deg,#1a1040,#0d1b2a);border:1px solid rgba(255,255,255,0.08);border-radius:24px;max-width:520px;width:100%;padding:36px 32px;text-align:center;box-shadow:0 20px 60px rgba(0,0,0,0.5);">
          <div v-if="!voteDone">
            <div style="font-size:48px;margin-bottom:12px;">🎉</div>
            <h2 style="font-size:24px;font-weight:700;margin-bottom:8px;">欢迎来到学趣星球！</h2>
            <p style="font-size:14px;color:var(--md-text-secondary);margin-bottom:20px;">
              请全班投票选择你们喜欢的宠物类别<br>
              <span style="font-size:12px;opacity:0.7;">选定后每人可免费选择一只心仪的宠物</span>
            </p>
            <div style="display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin-bottom:20px;max-height:320px;overflow-y:auto;">
              <button v-for="s in allSeries" :key="s.id" @click="voteSeries = s.id"
                :style="{
                  padding:'16px 12px', borderRadius:'16px', cursor:'pointer', transition:'0.2s', fontFamily:'inherit',
                  border: voteSeries === s.id ? '2px solid var(--md-primary)' : '1px solid rgba(255,255,255,0.06)',
                  background: voteSeries === s.id ? 'rgba(167,139,250,0.1)' : 'rgba(255,255,255,0.02)',
                  color: voteSeries === s.id ? 'var(--md-primary-light)' : 'var(--md-text)',
                }">
                <div style="font-size:32px;margin-bottom:6px;">{{ s.emoji }}</div>
                <div style="font-size:14px;font-weight:600;">{{ s.name }}</div>
                <div style="font-size:11px;color:var(--md-text-secondary);margin-top:2px;">{{ s.species.length }}种宠物</div>
              </button>
            </div>
            <button @click="confirmVote" :disabled="voting"
              style="width:100%;padding:14px;border-radius:14px;border:none;background:linear-gradient(135deg,var(--md-primary),var(--md-secondary));color:#fff;font-size:16px;font-weight:700;cursor:pointer;font-family:inherit;">
              {{ voting ? '投票中...' : '✅ 选择「' + getSeriesName(voteSeries) + '」系列' }}
            </button>
          </div>
          <div v-else>
            <div style="font-size:64px;margin-bottom:16px;">🎊</div>
            <h2 style="font-size:22px;font-weight:700;margin-bottom:8px;">选择成功！</h2>
            <p style="font-size:14px;color:var(--md-text-secondary);">
              已选定「{{ getSeriesName(voteSeries) }}」系列
            </p>
            <p style="font-size:13px;color:var(--md-gold);margin-top:8px;">
              现在去为每位同学免费选择宠物吧！
            </p>
          </div>
        </div>
      </div>
    </Transition>

    <nav class="sidebar">
      <div class="logo">
        <div class="brand"><span>🌌</span> 学趣星球</div>
      </div>

      <div class="class-badge" v-if="classInfo">
        <span class="class-name">{{ classInfo.name }}</span>
        <span class="class-count">{{ classInfo.student_count || '--' }} 人</span>
      </div>

      <div class="nav-list">
        <button v-for="item in navItems" :key="item.page"
          :class="['nav-item', { active: activeNav === item.page }]"
          @click="navigate(item.page)">
          <span class="icon">{{ item.icon }}</span> {{ item.label }}
        </button>
      </div>

      <!-- 底部：退出按钮 -->
      <div class="sidebar-footer">
        <button class="exit-btn" @click="goToLogin">✕ 退出班级</button>
      </div>
    </nav>

    <main class="main-content">
      <router-view v-slot="{ Component }">
        <component :is="Component" :key="$route.fullPath" />
      </router-view>
    </main>
  </div>
</template>

<style scoped>
.app-shell { display: flex; min-height: 100vh; }
.sidebar {
  width: var(--md-sidebar-width); background: var(--md-surface-2);
  border-right: 1px solid rgba(255,255,255,0.05);
  padding: 24px 16px 16px; display: flex; flex-direction: column;
  position: sticky; top: 0; height: 100vh; flex-shrink: 0;
  backdrop-filter: blur(12px); box-shadow: 4px 0 20px rgba(0,0,0,0.3); z-index: 10;
}
.logo { font-size: 22px; font-weight: 700; padding: 8px 12px; margin-bottom: 12px; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 12px; }
.brand { display: flex; align-items: center; gap: 10px; background: linear-gradient(135deg,var(--md-primary),var(--md-secondary)); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.brand span { font-size: 28px; -webkit-text-fill-color: initial; }
.class-badge { text-align: center; padding: 12px; margin-bottom: 12px; background: rgba(255,255,255,0.03); border-radius: var(--md-radius); border: 1px solid rgba(255,255,255,0.04); }
.class-name { font-size: 16px; font-weight: 700; display: block; }
.class-count { font-size: 12px; color: var(--md-text-secondary); }
.nav-list { display: flex; flex-direction: column; gap: 6px; flex: 1; }
.nav-item { display: flex; align-items: center; gap: 14px; padding: 12px 16px; border-radius: var(--md-radius); cursor: pointer; transition: 0.2s; color: var(--md-text-secondary); border: none; background: transparent; width: 100%; font-size: 16px; font-weight: 500; font-family: inherit; text-align: left; }
.nav-item:hover { background: rgba(255,255,255,0.05); color: #fff; }
.nav-item.active { background: rgba(167,139,250,0.15); color: var(--md-primary-light); box-shadow: inset 3px 0 0 var(--md-primary); }
.nav-item .icon { font-size: 22px; width: 28px; text-align: center; }

.sidebar-footer { border-top: 1px solid rgba(255,255,255,0.05); padding-top: 12px; }
.exit-btn { width: 100%; padding: 10px; border-radius: var(--md-radius); border: 1px solid rgba(255,100,100,0.15); background: rgba(255,100,100,0.08); color: #fca5a5; font-size: 14px; font-weight: 500; cursor: pointer; transition: 0.2s; font-family: inherit; }
.exit-btn:hover { background: rgba(255,100,100,0.15); }

.main-content { flex: 1; padding: 28px 32px 40px; max-width: calc(100% - var(--md-sidebar-width)); overflow-x: hidden; }

@media (max-width: 768px) {
  .sidebar { width: 100%; height: auto; position: sticky; flex-direction: row; flex-wrap: wrap; padding: 12px 16px; border-right: none; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .logo { margin-bottom: 0; border-bottom: none; padding-bottom: 0; font-size: 18px; flex: 1; }
  .class-badge { display: none; }
  .nav-list { flex-direction: row; gap: 4px; flex: 2; justify-content: flex-end; }
  .nav-item { padding: 8px 12px; font-size: 14px; }
  .nav-item .icon { font-size: 18px; width: 24px; }
  .sidebar-footer { display: none; }
  .main-content { padding: 16px; max-width: 100%; }
}
</style>
