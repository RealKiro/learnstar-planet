<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

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
const themeStore = useThemeStore()

const activeNav = computed(() => String(route.name))

function navigate(name: string) {
  router.push({ name })
}
</script>

<template>
  <div class="app-shell">
    <!-- 侧边栏（功能分类与版面设计.txt 版式） -->
    <nav class="sidebar">
      <div class="logo">
        <div class="brand">
          <span>🌌</span> 学趣星球
        </div>
      </div>

      <div class="nav-list">
        <template v-for="section in props.navItems" :key="section.section">
          <!-- 分组标题（AstrBot 式小节标签）；空分组不渲染 -->
          <template v-if="section.items.length">
            <div class="nav-section">{{ section.section }}</div>
            <button
              v-for="item in section.items"
              :key="item.page"
              :class="['nav-item', { active: activeNav === item.page }]"
              @click="navigate(item.page)"
            >
              <span class="icon">{{ item.icon }}</span> {{ item.label }}
            </button>
          </template>
        </template>
      </div>

      <!-- 侧边栏底部扩展区（系列选择器等） -->
      <div v-if="$slots['sidebar-extra']" class="series-selector">
        <slot name="sidebar-extra" />
      </div>

      <!-- 底部：退出按钮 -->
      <div class="sidebar-footer">
        <button v-if="props.showLogout" class="exit-btn" @click="emit('logout')">✕ 退出登录</button>
      </div>
    </nav>

    <!-- 主内容 -->
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <component :is="Component" :key="$route.fullPath" />
      </router-view>
    </main>

    <!-- 右下角主题切换浮动按钮 -->
    <button class="theme-fab" @click="themeStore.toggle()" :title="themeStore.isDark ? '切换到日间主题' : '切换到夜间主题'">
      {{ themeStore.isDark ? '☀️' : '🌙' }}
    </button>
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
  border-right: 1px solid var(--tint-2);
  padding: 24px 14px 20px;
  display: flex;
  flex-direction: column;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  flex-shrink: 0;
  backdrop-filter: blur(12px);
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
  border-bottom: 1px solid var(--tint-2);
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

.sidebar-footer { border-top: 1px solid var(--tint-2); padding-top: 12px; }
.theme-fab {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  color: var(--color-text);
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--md-elevation);
  z-index: 999;
  transition: 0.2s;
}
.theme-fab:hover { background: var(--tint-2); transform: scale(1.05); }
.exit-btn {
  width: 100%; padding: 10px; border-radius: var(--md-radius);
  border: 1px solid rgba(255,100,100,0.15);
  background: rgba(255,100,100,0.08); color: var(--color-danger-text);
  font-size: 14px; font-weight: 500;
  cursor: pointer; transition: 0.2s; font-family: inherit;
}
.exit-btn:hover { background: rgba(255,100,100,0.15); }

/* 导航列表 */
.nav-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

/* 分组标题（AstrBot 式小节标签） */
.nav-section {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  opacity: 0.75;
  padding: 14px 12px 5px;
  user-select: none;
}
.nav-section:first-child { padding-top: 4px; }

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 9px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.18s var(--ease-smooth), color 0.18s var(--ease-smooth);
  color: var(--md-text-secondary);
  border: none;
  background: transparent;
  width: 100%;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  text-align: left;
  position: relative;
}

.nav-item:hover {
  background: var(--tint-2);
  color: var(--color-text);
}

.nav-item.active {
  background: var(--tint-2);
  color: var(--color-primary);
  font-weight: 600;
}

/* 激活态左侧指示条（Linux.do 当前项式样） */
.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 18px;
  border-radius: 2px;
  background: var(--md-primary);
}

.nav-item .icon {
  font-size: 18px;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

/* 系列选择器区 */
.series-selector {
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid var(--tint-2);
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
    border-bottom: 1px solid var(--tint-2);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
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
    align-items: center;
    flex-wrap: wrap;
  }
  .nav-section { display: none; }
  .nav-item {
    padding: 8px 12px;
    font-size: 14px;
  }
  .nav-item.active::before { display: none; }
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
