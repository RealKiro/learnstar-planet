<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import SidebarLayout from '@/components/layout/SidebarLayout.vue'

const router = useRouter()
const authStore = useAuthStore()

function logout() {
  authStore.logout()
  router.replace({ name: 'landing' })
}

const navItems = [
  { section: '概览', items: [
    { page: 'admin-dashboard', label: '全校看板', icon: '📊' },
  ]},
  { section: '账号管理', items: [
    { page: 'admin-teachers', label: '教师账号', icon: '👨‍🏫' },
    { page: 'admin-parents', label: '家长账号', icon: '👪' },
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
    { page: 'admin-exchange-rate', label: '汇率配置', icon: '💰' },
  ]},
  { section: '调试', items: [
    { page: 'admin-debug', label: '高级调试', icon: '🔧' },
  ]},
]
</script>

<template>
  <SidebarLayout
    role-label="学校管理"
    :nav-items="navItems"
    avatar-gradient="var(--gradient-gold)"
    :show-logout="true"
    @logout="logout"
  />
</template>
