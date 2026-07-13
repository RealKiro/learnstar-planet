<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { computed, ref } from 'vue'

const props = defineProps<{
  roleLabel: string
  navItems: Array<{
    section: string
    items: Array<{ page: string; label: string; icon: string }>
  }>
  avatarGradient?: string
  showLogout?: boolean
}>()

const emit = defineEmits<{ logout: [] }>()

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const theme = useThemeStore()

const activeNav = computed(() => String(route.name))
const sidebarCollapsed = ref(false)

function navigate(name: string) {
  router.push({ name })
}

function goHome() {
  router.push({ name: 'landing' })
}
</script>

<template>
  <div class="app-shell">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ 'sidebar--collapsed': sidebarCollapsed }">
      <!-- Brand -->
      <div class="sidebar-brand">
        <button class="brand-btn" @click="goHome" title="返回首页">
          <span class="brand-logo">🌌</span>
          <span class="brand-text">
            <span class="brand-name">学趣星球</span>
            <span class="brand-role">{{ roleLabel }}</span>
          </span>
        </button>
        <button class="collapse-btn" @click="sidebarCollapsed = !sidebarCollapsed" title="折叠侧栏">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M15 18l-6-6 6-6"/></svg>
        </button>
      </div>

      <!-- 用户信息 -->
      <div class="sidebar-user">
        <div class="user-avatar" :style="{ background: props.avatarGradient || 'var(--gradient-primary)' }">
          {{ authStore.displayName.charAt(0) }}
        </div>
        <div class="user-info">
          <div class="user-name">{{ authStore.displayName }}</div>
          <slot name="user-meta">
            <div class="user-role">{{ roleLabel }}</div>
          </slot>
        </div>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <div v-for="section in props.navItems" :key="section.section" class="nav-section">
          <div class="nav-section-label">{{ section.section }}</div>
          <div
            v-for="item in section.items" :key="item.page"
            :class="['nav-item', { 'nav-item--active': activeNav === item.page }]"
            role="button" tabindex="0"
            @click="navigate(item.page)" @keydown.enter="navigate(item.page)"
          >
            <span class="nav-icon">{{ item.icon }}</span>
            <span class="nav-label">{{ item.label }}</span>
          </div>
        </div>
      </nav>

      <!-- 底部 -->
      <div class="sidebar-footer">
        <button v-if="props.showLogout" class="footer-btn" @click="emit('logout')">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          退出
        </button>
        <span v-else class="footer-free">全部功能免费</span>
        <button class="theme-btn" @click="theme.toggle()" :title="theme.isDark ? '切换亮色' : '切换暗色'">
          {{ theme.isDark ? '☀️' : '🌙' }}
        </button>
      </div>
    </aside>

    <!-- 主内容 -->
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg);
}

/* ===== 侧边栏 ===== */
.sidebar {
  width: 256px;
  flex-shrink: 0;
  position: fixed;
  top: 0; left: 0;
  height: 100vh;
  z-index: 100;
  background: var(--color-bg-card);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  transition: transform 0.3s cubic-bezier(.16,1,.3,1);
}

/* Brand */
.sidebar-brand {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 20px 16px;
  border-bottom: 1px solid var(--color-border);
}

.brand-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  text-align: left;
  font-family: inherit;
}

.brand-logo {
  font-size: 26px;
  line-height: 1;
  transition: transform 0.3s var(--ease-bounce);
}

.brand-btn:hover .brand-logo {
  transform: rotate(-10deg) scale(1.1);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-name {
  font-weight: 700;
  font-size: 16px;
  color: var(--color-text);
  letter-spacing: -0.02em;
}

.brand-role {
  font-size: 11px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.collapse-btn {
  width: 28px; height: 28px;
  border: none; border-radius: 8px;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.collapse-btn:hover {
  background: var(--color-bg);
  color: var(--color-text);
}

/* User */
.sidebar-user {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--color-border);
}

.user-avatar {
  width: 38px; height: 38px;
  border-radius: 10px;
  color: #F1F5F9;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  font-size: 11px;
  color: var(--color-text-secondary);
  margin-top: 1px;
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 12px 0;
}

.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: 2px;
}

.nav-section {
  margin-bottom: 4px;
}

.nav-section-label {
  font-size: 10px;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 1.2px;
  padding: 12px 20px 6px;
  opacity: 0.6;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 20px;
  margin: 0 8px;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  border-radius: 8px;
  position: relative;
}

.nav-item:hover {
  color: var(--color-text);
  background: var(--color-bg);
}

.nav-item--active {
  color: var(--color-primary);
  background: rgba(79, 70, 229, 0.08);
  font-weight: 600;
}

.nav-item--active::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--color-primary);
  border-radius: 0 3px 3px 0;
}

.nav-icon {
  font-size: 16px;
  width: 22px;
  text-align: center;
  flex-shrink: 0;
}

.nav-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Footer */
.sidebar-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-top: 1px solid var(--color-border);
}

.footer-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}

.footer-btn:hover {
  background: var(--color-bg);
  color: var(--color-danger);
}

.footer-free {
  font-size: 11px;
  color: var(--color-text-secondary);
}

.theme-btn {
  width: 30px; height: 30px;
  border: none;
  border-radius: 8px;
  background: var(--color-bg);
  color: var(--color-text-secondary);
  font-size: 15px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.theme-btn:hover {
  background: var(--color-primary);
  color: #fff;
}

/* ===== 主内容 ===== */
.main-content {
  flex: 1;
  margin-left: 256px;
  padding: 32px;
  min-height: 100vh;
  max-width: 1200px;
  width: calc(100% - 256px);
  box-sizing: border-box;
}

/* ===== 页面过渡 ===== */
.page-enter-active,
.page-leave-active {
  transition: all 0.25s cubic-bezier(.16,1,.3,1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar--open {
    transform: translateX(0);
    box-shadow: 0 0 40px rgba(0,0,0,.15);
  }

  .main-content {
    margin-left: 0;
    padding: 20px;
    width: 100%;
  }

  .collapse-btn {
    display: none;
  }
}
</style>
