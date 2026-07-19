<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { computed } from 'vue'

const props = defineProps<{
  roleLabel: string
  navItems: Array<{
    section: string
    items: Array<{ page: string; label: string; icon: string }>
  }>
  showLogout?: boolean
}>()

const emit = defineEmits<{ logout: [] }>()

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeNav = computed(() => String(route.name))

function navigate(name: string) {
  router.push({ name })
}

function goHome() {
  router.push({ name: 'landing' })
}
</script>

<template>
  <div class="app-shell">
    <!-- 侧边栏（功能分类与版面设计.txt 版式） -->
    <nav class="sidebar">
      <div class="logo">
        <div class="brand">
          <span>🐾</span> 宠物星球
        </div>
        <button v-if="props.showLogout" class="btn-exit" @click="emit('logout')">✕ 退出</button>
      </div>

      <div class="nav-list">
        <template v-for="section in props.navItems" :key="section.section">
          <button
            v-for="item in section.items"
            :key="item.page"
            :class="['nav-item', { active: activeNav === item.page }]"
            @click="navigate(item.page)"
          >
            <span class="icon">{{ item.icon }}</span> {{ item.label }}
          </button>
        </template>
      </div>

      <!-- 侧边栏底部扩展区（系列选择器等） -->
      <div v-if="$slots['sidebar-extra']" class="series-selector">
        <slot name="sidebar-extra" />
      </div>
    </nav>

    <!-- 主内容 -->
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <component :is="Component" :key="$route.fullPath" />
      </router-view>
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  display: flex;
  min-height: 100vh;
}

/* ===== 侧边栏 ===== */
.sidebar {
  width: var(--md-sidebar-width);
  background: var(--md-surface-2);
  border-right: 1px solid rgba(255, 255, 255, 0.05);
  padding: 24px 16px 20px;
  display: flex;
  flex-direction: column;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  flex-shrink: 0;
  backdrop-filter: blur(12px);
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.3);
  z-index: 10;
}

.logo {
  font-size: 22px;
  font-weight: 700;
  padding: 8px 12px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 12px;
}

.logo .brand {
  display: flex;
  align-items: center;
  gap: 10px;
  background: linear-gradient(135deg, var(--md-primary), var(--md-secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo .brand span {
  font-size: 28px;
  -webkit-text-fill-color: initial;
}

.btn-exit {
  background: rgba(255, 100, 100, 0.1);
  border: 1px solid rgba(255, 100, 100, 0.15);
  color: #fca5a5;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  transition: 0.2s;
  font-weight: 500;
  font-family: inherit;
}

.btn-exit:hover {
  background: rgba(255, 100, 100, 0.2);
}

/* 导航列表 */
.nav-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 16px;
  border-radius: var(--md-radius);
  cursor: pointer;
  transition: 0.2s;
  color: var(--md-text-secondary);
  border: none;
  background: transparent;
  width: 100%;
  font-size: 16px;
  font-weight: 500;
  font-family: inherit;
  text-align: left;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.nav-item.active {
  background: rgba(167, 139, 250, 0.15);
  color: var(--md-primary-light);
  box-shadow: inset 3px 0 0 var(--md-primary);
}

.nav-item .icon {
  font-size: 22px;
  width: 28px;
  text-align: center;
}

/* 系列选择器区 */
.series-selector {
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

/* ===== 主内容 ===== */
.main-content {
  flex: 1;
  padding: 28px 32px 40px;
  max-width: calc(100% - var(--md-sidebar-width));
  overflow-x: hidden;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .sidebar {
    width: 100%;
    height: auto;
    position: sticky;
    top: 0;
    flex-direction: row;
    flex-wrap: wrap;
    padding: 12px 16px;
    border-right: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  }
  .logo {
    margin-bottom: 0;
    border-bottom: none;
    padding-bottom: 0;
    font-size: 18px;
    flex: 1;
  }
  .nav-list {
    flex-direction: row;
    gap: 4px;
    flex: 2;
    justify-content: flex-end;
  }
  .nav-item {
    padding: 8px 12px;
    font-size: 14px;
  }
  .nav-item .icon {
    font-size: 18px;
    width: 24px;
  }
  .series-selector {
    margin-top: 0;
    padding-top: 0;
    border-top: none;
    width: 100%;
    margin-top: 8px;
  }
  .main-content {
    padding: 16px;
    max-width: 100%;
  }
}
</style>
