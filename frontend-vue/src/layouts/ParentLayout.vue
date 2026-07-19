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
    { page: 'parent-home', label: '孩子概览', icon: '🏠' },
  ]},
  { section: '成长详情', items: [
    { page: 'parent-pet', label: '宠物详情', icon: '🌟' },
    { page: 'parent-scores', label: '积分明细', icon: '⭐' },
    { page: 'parent-growth', label: '成长日志', icon: '📝' },
    { page: 'parent-ranking', label: '班级排名', icon: '🏆' },
  ]},
  { section: '通知', items: [
    { page: 'parent-notices', label: '班级通知', icon: '📢' },
  ]},
]
</script>

<template>
  <SidebarLayout
    role-label="家长端"
    :nav-items="navItems"
    :show-logout="true"
    @logout="logout"
  />
</template>
