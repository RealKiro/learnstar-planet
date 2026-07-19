<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import SidebarLayout from '@/components/layout/SidebarLayout.vue'

const authStore = useAuthStore()

const navItems = [
  { section: '🏠 班级总览', items: [
    { page: 'teacher-dashboard', label: '班级总览', icon: '🏠' },
  ]},
  { section: '✏️ 课堂评价', items: [
    { page: 'teacher-scores', label: '课堂评价', icon: '✏️' },
    { page: 'teacher-rules', label: '积分规则', icon: '📋' },
    { page: 'teacher-pets', label: '宠物花园', icon: '🌟' },
  ]},
  { section: '🏆 年级战场', items: [
    { page: 'teacher-leaderboard', label: '年级排行', icon: '🏆' },
    { page: 'teacher-students', label: '学生数据', icon: '👨‍🎓' },
  ]},
  { section: '📚 宠物图鉴', items: [
    { page: 'teacher-pets', label: '进化图鉴', icon: '📚' },
  ]},
  { section: '课堂工具', items: [
    { page: 'teacher-broadcast', label: '实时广播', icon: '📡' },
    { page: 'teacher-attendance', label: '智能考勤', icon: '✅' },
    { page: 'teacher-homework', label: '扫码收作业', icon: '📱' },
    { page: 'teacher-quiz', label: '在线答题', icon: '📝' },
    { page: 'teacher-ai', label: 'AI助教', icon: '🤖' },
  ]},
  { section: '学业与沟通', items: [
    { page: 'teacher-grades', label: '成绩管理', icon: '📊' },
    { page: 'teacher-notices', label: '班级通知', icon: '📢' },
    { page: 'teacher-reports', label: '数据报表', icon: '📈' },
  ]},
  { section: '运营', items: [
    { page: 'teacher-shop', label: '积分商城', icon: '🛍️' },
    { page: 'teacher-exchange', label: '兑换中心', icon: '🔄' },
  ]},
  { section: '设置', items: [
    { page: 'teacher-settings', label: '账号设置', icon: '⚙️' },
    { page: 'teacher-classroom', label: '班级大屏', icon: '🖥️' },
  ]},
]

const currentSeriesId = 'myth'
</script>

<template>
  <SidebarLayout role-label="教师端" :nav-items="navItems">
    <template #user-meta>
      <div style="font-size:12px;color:var(--color-text-secondary);">{{ authStore.user?.class_names?.join('、') || '教师' }}</div>
    </template>
  </SidebarLayout>
</template>
