import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/pages/landing/LandingPage.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/auth/LoginPage.vue'),
      props: (route) => ({ initialRole: route.query.role || 'teacher', mode: route.query.mode || 'account' }),
    },
    {
      path: '/auth/callback',
      name: 'oauth-callback',
      component: () => import('@/pages/auth/OAuthCallbackPage.vue'),
    },
    {
      path: '/teacher',
      component: () => import('@/layouts/TeacherLayout.vue'),
      meta: { requiresAuth: true, role: 'teacher' },
      children: [
        { path: '', redirect: { name: 'teacher-dashboard' } },
        { path: 'dashboard', name: 'teacher-dashboard', component: () => import('@/pages/teacher/DashboardPage.vue') },
        { path: 'classroom', name: 'teacher-classroom', component: () => import('@/pages/teacher/ClassroomDisplay.vue') },
        { path: 'students', name: 'teacher-students', component: () => import('@/pages/teacher/StudentsPage.vue') },
        { path: 'pets', name: 'teacher-pets', component: () => import('@/pages/teacher/PetsPage.vue') },
        { path: 'scores', name: 'teacher-scores', component: () => import('@/pages/teacher/ScoresPage.vue') },
        { path: 'rules', name: 'teacher-rules', component: () => import('@/pages/teacher/RulesPage.vue') },
        { path: 'leaderboard', name: 'teacher-leaderboard', component: () => import('@/pages/teacher/LeaderboardPage.vue') },
        { path: 'pk', name: 'teacher-pk', component: () => import('@/pages/teacher/PKPage.vue') },
        { path: 'shop', name: 'teacher-shop', component: () => import('@/pages/teacher/ShopPage.vue') },
        { path: 'exchange', name: 'teacher-exchange', component: () => import('@/pages/teacher/ExchangeCenterPage.vue') },
        { path: 'communication', name: 'teacher-communication', component: () => import('@/pages/teacher/CommunicationPage.vue') },
        { path: 'broadcast', name: 'teacher-broadcast', component: () => import('@/pages/teacher/BroadcastPage.vue') },
        { path: 'attendance', name: 'teacher-attendance', component: () => import('@/pages/teacher/AttendancePage.vue') },
        { path: 'grades', name: 'teacher-grades', component: () => import('@/pages/teacher/GradesPage.vue') },
        { path: 'ai', name: 'teacher-ai', component: () => import('@/pages/teacher/AIPage.vue') },
        { path: 'notices', name: 'teacher-notices', component: () => import('@/pages/teacher/NoticesPage.vue') },
        { path: 'reports', name: 'teacher-reports', component: () => import('@/pages/teacher/ReportsPage.vue') },
        { path: 'settings', name: 'teacher-settings', component: () => import('@/pages/teacher/SettingsPage.vue') },
      ],
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true, role: 'school_admin' },
      children: [
        { path: '', redirect: { name: 'admin-dashboard' } },
        { path: 'dashboard', name: 'admin-dashboard', component: () => import('@/pages/admin/DashboardPage.vue') },
        { path: 'teachers', name: 'admin-teachers', component: () => import('@/pages/admin/TeachersPage.vue') },
        { path: 'parents', name: 'admin-parents', component: () => import('@/pages/admin/ParentsPage.vue') },
        { path: 'classes', name: 'admin-classes', component: () => import('@/pages/admin/ClassesPage.vue') },
        { path: 'students', name: 'admin-students', component: () => import('@/pages/admin/StudentsPage.vue') },
        { path: 'upgrade', name: 'admin-upgrade', component: () => import('@/pages/admin/GradeUpgradePage.vue') },
        { path: 'reports', name: 'admin-reports', component: () => import('@/pages/admin/ReportsPage.vue') },
        { path: 'school', name: 'admin-school', component: () => import('@/pages/admin/SchoolSettingsPage.vue') },
        { path: 'exchange-rate', name: 'admin-exchange-rate', component: () => import('@/pages/admin/ExchangeRatePage.vue') },
        { path: 'debug', name: 'admin-debug', component: () => import('@/pages/admin/AdminDebugPage.vue') },
      ],
    },
    {
      path: '/parent',
      component: () => import('@/layouts/ParentLayout.vue'),
      meta: { requiresAuth: true, role: 'parent' },
      children: [
        { path: '', redirect: { name: 'parent-home' } },
        { path: 'home', name: 'parent-home', component: () => import('@/pages/parent/HomePage.vue') },
        { path: 'scores', name: 'parent-scores', component: () => import('@/pages/parent/ScoresPage.vue') },
        { path: 'growth', name: 'parent-growth', component: () => import('@/pages/parent/GrowthPage.vue') },
        { path: 'pet', name: 'parent-pet', component: () => import('@/pages/parent/PetDetailPage.vue') },
        { path: 'ranking', name: 'parent-ranking', component: () => import('@/pages/parent/RankingPage.vue') },
        { path: 'notices', name: 'parent-notices', component: () => import('@/pages/parent/NoticesPage.vue') },
      ],
    },
    // ===== 班级大屏（无认证，使用班级码 Token） =====
    {
      path: '/display',
      children: [
        { path: 'login', redirect: '/login?mode=code' },
        { path: 'code', redirect: '/login?mode=code' },
      ],
    },

    // ===== 班级教室端（班级码登录，无需教师账号，展示设计文档4大核心模块） =====
    {
      path: '/classroom',
      component: () => import('@/layouts/ClassroomLayout.vue'),
      children: [
        { path: '', redirect: { name: 'classroom-overview' } },
        { path: 'overview', name: 'classroom-overview', component: () => import('@/pages/classroom/OverviewPage.vue') },
        { path: 'scores', name: 'classroom-scores', component: () => import('@/pages/classroom/ScoresPage.vue') },
        { path: 'pk', name: 'classroom-pk', component: () => import('@/pages/classroom/PKPage.vue') },
        { path: 'pokedex', name: 'classroom-pokedex', component: () => import('@/pages/classroom/PokedexPage.vue') },
      ],
    },
	],
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return next({ name: 'landing' })
  }

  if (to.meta.role && authStore.user?.role !== to.meta.role) {
    if (authStore.isAdmin) return next({ name: 'admin-dashboard' })
    if (authStore.isTeacher) return next({ name: 'teacher-dashboard' })
    if (authStore.isParent) return next({ name: 'parent-home' })
    return next({ name: 'landing' })
  }

  next()
})

export default router
