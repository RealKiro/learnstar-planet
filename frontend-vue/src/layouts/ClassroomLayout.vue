<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const classInfo = ref<{ id: number; name: string; student_count?: number } | null>(null)

const activeNav = computed(() => String(route.name))

const navItems = [
  { page: 'classroom-overview', label: '班级总览', icon: '🏠' },
  { page: 'classroom-scores', label: '课堂评价', icon: '✏️' },
  { page: 'classroom-pk', label: '年级战场', icon: '🏆' },
  { page: 'classroom-pokedex', label: '宠物图鉴', icon: '📚' },
]

function navigate(name: string) { router.push({ name }) }
function goToLogin() { sessionStorage.clear(); router.push({ name: 'login', query: { mode: 'code' } }) }

onMounted(() => {
  const ci = sessionStorage.getItem('class_info')
  if (ci) classInfo.value = JSON.parse(ci)
})
</script>

<template>
  <div class="app-shell">
    <nav class="sidebar">
      <div class="logo">
        <div class="brand"><span>🐾</span> 宠物星球</div>
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
