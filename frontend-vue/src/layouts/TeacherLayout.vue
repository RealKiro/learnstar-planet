<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import SidebarLayout from '@/components/layout/SidebarLayout.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToastStore()
function logout() { authStore.logout(); router.replace({ name: 'landing' }) }

interface ClassItem { class_id: number; class_name: string; grade: string; role: string }
const myClasses = ref<ClassItem[]>([])
const activeClassId = ref<number | null>(null)
const switchingClass = ref(false)
const activeClass = computed(() => myClasses.value.find(c => c.class_id === activeClassId.value))

// 基础模式（班级码进入）只显示课堂核心功能
const isBasic = computed(() => route.meta?.basic === true)

const fullNav = [
  { section: '概览', items: [
    { page: 'teacher-dashboard', label: '班级总览', icon: '🏠' },
  ]},
  { section: '课堂教学', items: [
    { page: 'teacher-scores', label: '课堂评价', icon: '✏️' },
    { page: 'teacher-rules', label: '积分规则', icon: '📋' },
    { page: 'teacher-attendance', label: '智能考勤', icon: '✅' },
  ]},
  { section: '成长激励', items: [
    { page: 'teacher-pets', label: '宠物花园', icon: '🌟' },
    { page: 'teacher-leaderboard', label: '排行榜', icon: '🏆' },
    { page: 'teacher-pk', label: '年级PK', icon: '⚔️' },
  ]},
  { section: '数据中心', items: [
    { page: 'teacher-grades', label: '成绩管理', icon: '📊' },
    { page: 'teacher-reports', label: '数据报表', icon: '📈' },
  ]},
  { section: '沟通协作', items: [
    { page: 'teacher-communication', label: '消息中心', icon: '📢' },
  ]},
  { section: '系统管理', items: [
    { page: 'teacher-ai', label: 'AI助教', icon: '🤖' },
    { page: 'teacher-shop', label: '积分商城', icon: '🛍️' },
  ]},
  { section: '设置', items: [
    { page: 'teacher-settings', label: '账号设置', icon: '⚙️' },
  ]},
]

const basicNav = [
  { section: '概览', items: [
    { page: 'teacher-dashboard-basic', label: '班级总览', icon: '🏠' },
  ]},
  { section: '课堂', items: [
    { page: 'teacher-scores-basic', label: '课堂评价', icon: '✏️' },
    { page: 'teacher-leaderboard-basic', label: '排行榜', icon: '🏆' },
    { page: 'teacher-pk-basic', label: '年级PK', icon: '⚔️' },
    { page: 'teacher-pets-basic', label: '宠物图鉴', icon: '📚' },
  ]},
  { section: '互动', items: [
    { page: 'teacher-ai-basic', label: 'AI 助手', icon: '🤖' },
  ]},
]

const navItems = computed(() => isBasic.value ? basicNav : fullNav)

async function loadMyClasses() {
  try {
    const res = await apiGet<{ data: ClassItem[] }>('/api/v1/teacher/my-classes')
    myClasses.value = res.data || []
    if (myClasses.value.length > 0) {
      activeClassId.value = myClasses.value[0].class_id
    }
  } catch { /* ignore */ }
}

async function switchToClass(classId: number) {
  if (classId === activeClassId.value) return
  switchingClass.value = true
  try {
    await apiPost('/api/v1/teacher/switch-class', { class_id: classId })
    activeClassId.value = classId
    toast.show(`已切换至「${myClasses.value.find(c => c.class_id === classId)?.class_name || ''}」`, 'success')
  } catch (e: any) {
    toast.show(e?.response?.data?.message || '切换失败', 'error')
  } finally {
    switchingClass.value = false
  }
}

onMounted(() => {
  loadMyClasses()
})
</script>

<template>
  <SidebarLayout role-label="教师端" :nav-items="navItems" :show-logout="true" @logout="logout">
    <template #user-meta>
      <div style="font-size:12px;color:var(--color-text-secondary);">
        {{ authStore.user?.class_names?.join('、') || '教师' }}
      </div>
    </template>

    <!-- 侧边栏底部：班级切换器 -->
    <template #sidebar-extra>
      <div class="class-switcher">
        <div class="cs-label">📚 当前班级</div>
        <div v-if="myClasses.length === 0" class="cs-empty">暂未分配班级</div>
        <div v-else class="cs-select-wrap">
          <select
            :value="activeClassId ?? undefined"
            class="cs-select"
            :disabled="switchingClass"
            @change="switchToClass(Number(($event.target as HTMLSelectElement).value))"
          >
            <option v-for="c in myClasses" :key="c.class_id" :value="c.class_id">
              {{ c.class_name }}
            </option>
          </select>
        </div>
        <div v-if="activeClass" class="cs-info">
          <span class="cs-role">{{ { head_teacher: '班主任', co_teacher: '副班', subject_teacher: '科任', grade_lead: '年级首席', admin_director: '分管行政' }[activeClass.role] || activeClass.role }}</span>
          <span class="cs-grade">{{ activeClass.grade }}</span>
        </div>
      </div>
    </template>
  </SidebarLayout>
</template>

<style scoped>
.class-switcher { font-size: 12px; }
.cs-label { font-size: 11px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 6px; letter-spacing: 0.03em; }
.cs-empty { font-size: 12px; color: var(--color-text-secondary); padding: 8px 0; text-align: center; }
.cs-select-wrap { position: relative; }
.cs-select { width: 100%; padding: 8px 10px; border-radius: 10px; background: var(--color-bg); border: 1px solid var(--color-border); color: var(--color-text); font-size: 13px; font-weight: 500; outline: none; cursor: pointer; font-family: inherit; appearance: auto; }
.cs-select:disabled { opacity: 0.5; cursor: not-allowed; }
.cs-info { display: flex; gap: 8px; justify-content: center; margin-top: 6px; font-size: 11px; }
.cs-role { color: var(--color-primary, #a78bfa); font-weight: 600; }
.cs-grade { color: var(--color-text-secondary); }
</style>
