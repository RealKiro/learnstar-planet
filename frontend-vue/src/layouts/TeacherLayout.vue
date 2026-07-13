<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { computed } from 'vue'
import SidebarLayout from '@/components/layout/SidebarLayout.vue'

const authStore = useAuthStore()

const navItems = [
  { section: '概览', items: [
    { page: 'teacher-dashboard', label: '数据看板', icon: '📊' },
  ]},
  { section: '学生管理', items: [
    { page: 'teacher-students', label: '学生列表', icon: '👨‍🎓' },
    { page: 'teacher-pets', label: '宠物园', icon: '🌟' },
  ]},
  { section: '积分中心', items: [
    { page: 'teacher-scores', label: '积分管理', icon: '⭐' },
    { page: 'teacher-rules', label: '积分规则', icon: '📋' },
    { page: 'teacher-leaderboard', label: '排行榜', icon: '🏆' },
    { page: 'teacher-shop', label: '积分商城', icon: '🛍️' },
  ]},
  { section: '课堂工具', items: [
    { page: 'teacher-broadcast', label: '实时广播', icon: '📡' },
    { page: 'teacher-attendance', label: '智能考勤', icon: '✅' },
    { page: 'teacher-homework', label: '扫码收作业', icon: '📱' },
    { page: 'teacher-quiz', label: '在线答题', icon: '📝' },
    { page: 'teacher-ai', label: 'AI助教', icon: '🤖' },
  ]},
  { section: '学业评价', items: [
    { page: 'teacher-grades', label: '成绩管理', icon: '📊' },
  ]},
  { section: '沟通', items: [
    { page: 'teacher-notices', label: '班级通知', icon: '📢' },
  ]},
  { section: '报表', items: [
    { page: 'teacher-reports', label: '数据报表', icon: '📈' },
  ]},
  { section: '设置', items: [
    { page: 'teacher-settings', label: '账号设置', icon: '⚙️' },
  ]},
  { section: '大屏', items: [
    { page: 'teacher-classroom', label: '班级大屏', icon: '🖥️' },
  ]},
]
</script>

<template>
  <SidebarLayout role-label="教师端" :nav-items="navItems">
    <template #user-meta>
      <div style="font-size:12px;color:var(--color-text-secondary);">{{ authStore.user?.class_names?.join('、') || '教师' }}</div>
    </template>
  </SidebarLayout>
</template>
