<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

const props = defineProps<{
  roleLabel: string
  navItems: Array<{
    section: string
    /** 图标为 emoji 字符，模板中直接渲染 */
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
    <!-- 侧边栏
         视觉对标：底色用 --ui-bg-subtle 与内容区（--ui-bg 白）拉开层次，
         因此不再需要右边框（两条线会与内容区卡片边框打架）。 -->
    <nav class="sidebar">
      <div class="logo">
        <div class="brand">
          <span>🌌</span> 学宠星球
        </div>
        <!-- 端别标识（教师端 / 学校管理）；移动端由 .brand-role 隐藏 -->
        <span class="brand-role">{{ props.roleLabel }}</span>
      </div>

      <!-- 用户元信息（如教师所带班级），由调用方通过 #user-meta 插槽填充 -->
      <div v-if="$slots['user-meta']" class="sidebar-user">
        <slot name="user-meta" />
      </div>

      <div class="nav-list">
        <template v-for="section in props.navItems" :key="section.section">
          <!-- 分组标题；空分组不渲染 -->
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
  background: var(--ui-bg-subtle);
  border-right: none;
  padding: 20px 12px 16px;
  display: flex;
  flex-direction: column;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  flex-shrink: 0;
  z-index: 10;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 2px 8px 14px;
}

/* 品牌标识：纯色 + 字重，不用渐变裁剪（14.5px 下渐变会让笔画发灰） */
.brand {
  display: flex;
  align-items: center;
  gap: 9px;
  font-size: 14.5px;
  font-weight: 650;
  letter-spacing: -0.01em;
  color: var(--ui-fg);
}
.brand-role {
  font-size: 11px;
  font-weight: 500;
  color: var(--ui-fg-subtle);
  flex-shrink: 0;
}

/* 用户元信息：单行次要文本，超长省略（如多个班级名拼接） */
.sidebar-user {
  padding: 0 8px 12px;
  overflow: hidden;
}
.sidebar-user :deep(*) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar-footer { border-top: 1px solid var(--ui-border); padding-top: 10px; }

.theme-fab {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 40px;
  height: 40px;
  border-radius: var(--ui-r-lg);
  border: 1px solid var(--ui-border);
  background: var(--ui-card);
  color: var(--ui-fg-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--ui-shadow-sm);
  z-index: 999;
  transition: 0.15s; font-size: 18px; line-height: 1;}
.theme-fab:hover { background: var(--ui-muted); color: var(--ui-fg); border-color: var(--ui-border-strong); }

.exit-btn {
  width: 100%;
  height: 34px;
  padding: 0 12px;
  border-radius: var(--ui-r-md);
  border: 1px solid transparent;
  background: var(--c-red-bg);
  color: var(--color-danger-text);
  font-size: 13.5px;
  font-weight: 500;
  cursor: pointer;
  transition: 0.15s;
  font-family: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
}
.exit-btn:hover { background: rgba(239, 68, 68, 0.14); }

/* 导航列表 */
.nav-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  flex: 1;
}

/* 分组标题（小节标签） */
.nav-section {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: var(--ui-fg-subtle);
  padding: 14px 8px 5px;
  user-select: none;
}
.nav-section:first-child { padding-top: 2px; }

/* 导航项：高度统一 32px，图标固定 16px，纵向节奏整齐 */
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 32px;
  padding: 0 8px;
  border-radius: var(--ui-r-md);
  cursor: pointer;
  transition: background 0.15s var(--ease-smooth), color 0.15s var(--ease-smooth);
  color: var(--ui-fg-muted);
  border: none;
  background: transparent;
  width: 100%;
  font-size: 13.5px;
  font-weight: 500;
  font-family: inherit;
  text-align: left;
}

.nav-item:hover {
  background: var(--ui-muted);
  color: var(--ui-fg);
}
.nav-item:hover .icon { color: var(--ui-fg-muted); }

/* 激活态：白卡片 + 微阴影，与 subtle 侧栏底拉开层次（取代原先的左竖条） */
.nav-item.active {
  background: var(--ui-card);
  color: var(--ui-fg);
  font-weight: 600;
  box-shadow: var(--ui-shadow-xs);
}
.nav-item.active .icon { color: var(--ui-brand); }

.nav-item .icon {
  color: var(--ui-fg-subtle);
  transition: color 0.15s var(--ease-smooth);
}

/* 系列选择器区 */
.series-selector {
  margin-top: auto;
  padding-top: 14px;
  border-top: 1px solid var(--ui-border);
}

/* ===== 主内容 ===== */
.main-content {
  flex: 1;
  padding: 24px 28px 40px;
  max-width: calc(100% - var(--md-sidebar-width));
  overflow-x: hidden;
  background: var(--ui-bg);
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
    align-items: center;
    padding: 10px 14px;
    border-bottom: 1px solid var(--ui-border);
    box-shadow: var(--ui-shadow-xs);
  }
  .logo {
    padding: 0;
    flex: 1;
  }
  .brand-role { display: none; }
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
    height: 30px;
    padding: 0 10px;
    font-size: 13px;
  }
  .nav-item.active { box-shadow: none; background: var(--ui-muted); }
  .series-selector {
    margin-top: 8px;
    padding-top: 8px;
    border-top: none;
    width: 100%;
  }
  .main-content {
    padding: 16px;
    max-width: 100%;
  }
}
</style>
