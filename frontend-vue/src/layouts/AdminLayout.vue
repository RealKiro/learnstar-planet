<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { computed } from 'vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const theme = useThemeStore()

const navItems = [
  { section: '概览', items: [
    { page: 'admin-dashboard', label: '全校看板', icon: '📊' },
  ]},
  { section: '账号管理', items: [
    { page: 'admin-teachers', label: '教师账号', icon: '👨‍🏫' },
  ]},
  { section: '班级管理', items: [
    { page: 'admin-classes', label: '班级列表', icon: '🏫' },
    { page: 'admin-students', label: '学生管理', icon: '🎒' },
  ]},
  { section: '学年管理', items: [
    { page: 'admin-upgrade', label: '学年升级', icon: '📅' },
  ]},
  { section: '报表', items: [
    { page: 'admin-reports', label: '全校报表', icon: '📈' },
  ]},
  { section: '系统', items: [
    { page: 'admin-school', label: '学校信息', icon: '⚙️' },
  ]},
]

const activeNav = computed(() => String(route.name))

function navigate(name: string) {
  router.push({ name })
}

function logout() {
  authStore.logout()
  router.replace({ name: 'landing' })
}
</script>

<template>
  <div class="app-layout">
    <aside class="app-layout__sidebar">
      <!-- Brand -->
      <div style="padding:24px;display:flex;align-items:center;gap:12px;border-bottom:1px solid var(--color-border);">
        <span style="font-size:28px;">🌌</span>
        <div>
          <div style="font-weight:700;font-size:18px;">学趣星球</div>
          <div style="font-size:11px;color:var(--color-text-secondary);">学校管理</div>
        </div>
      </div>

      <!-- User -->
      <div style="padding:16px 24px;display:flex;align-items:center;gap:12px;border-bottom:1px solid var(--color-border);">
        <div
          style="width:40px;height:40px;border-radius:var(--radius-md);background:var(--gradient-gold);color:#92400E;display:flex;align-items:center;justify-content:center;font-weight:700;font-size:16px;"
        >{{ authStore.displayName.charAt(0) }}</div>
        <div style="flex:1;">
          <div style="font-weight:500;font-size:14px;">{{ authStore.displayName }}</div>
          <div style="font-size:12px;color:var(--color-text-secondary);">学校管理</div>
        </div>
      </div>

      <!-- Navigation -->
      <nav style="flex:1;overflow-y:auto;padding:16px 0;">
        <div v-for="section in navItems" :key="section.section">
          <div style="font-size:11px;font-weight:700;color:var(--color-text-secondary);text-transform:uppercase;padding:8px 24px;letter-spacing:1px;">
            {{ section.section }}
          </div>
          <div
            v-for="item in section.items"
            :key="item.page"
            :class="{ 'nav-item--active': activeNav === item.page }"
            class="nav-item"
            @click="navigate(item.page)"
          >
            <span class="nav-item__icon">{{ item.icon }}</span>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </nav>

      <!-- Footer -->
      <div style="padding:16px 24px;border-top:1px solid var(--color-border);display:flex;align-items:center;justify-content:space-between;">
        <button class="btn btn-sm btn-ghost" @click="logout">退出登录</button>
        <button
          style="width:36px;height:36px;border-radius:var(--radius-sm);background:rgba(79,70,229,0.08);border:none;cursor:pointer;font-size:18px;"
          @click="theme.toggle()"
        >{{ theme.isDark ? '☀️' : '🌙' }}</button>
      </div>
    </aside>

    <main class="app-layout__main">
      <router-view v-slot="{ Component }">
        <transition name="page">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<style scoped>
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 24px;
  font-size: 14px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
  border-left: 3px solid transparent;
}
.nav-item:hover {
  color: var(--color-primary);
  background: rgba(79,70,229,0.04);
}
.nav-item--active {
  color: var(--color-primary);
  background: rgba(79,70,229,0.1);
  border-left-color: var(--color-primary);
  border-radius: 0 8px 8px 0;
  margin-right: 8px;
  font-weight: 600;
}
.nav-item__icon { font-size: 18px; width: 24px; text-align: center; }
</style>
