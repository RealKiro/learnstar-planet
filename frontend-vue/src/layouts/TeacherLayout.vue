<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { getAllSeries, getSeriesName, SERIES_SCENES } from '@/utils/petData'
import SidebarLayout from '@/components/layout/SidebarLayout.vue'

const authStore = useAuthStore()
const toast = useToastStore()

const classTotalScore = ref(0)
const currentSeries = ref('myth')
const switching = ref(false)

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
    { page: 'teacher-pk', label: '年级PK', icon: '🏆' },
    { page: 'teacher-leaderboard', label: '排行榜', icon: '📊' },
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

const allSeries = getAllSeries()

const currentSeriesName = computed(() => getSeriesName(currentSeries.value))
const currentSeriesColor = computed(() => SERIES_SCENES[currentSeries.value]?.primaryColor || '#6366F1')

async function loadClassInfo() {
  try {
    const res = await apiGet<{ data: { total_score: number; class_points?: number; settings?: { pet_series?: string } } }>('/api/v1/teacher/class')
    classTotalScore.value = res.data?.total_score || 0
    if (res.data?.settings?.pet_series) {
      currentSeries.value = res.data.settings.pet_series
    }
  } catch {
    // 静默失败，演示数据
    classTotalScore.value = 3840
  }
}

async function switchSeries(seriesId: string) {
  if (seriesId === currentSeries.value) return
  switching.value = true
  try {
    await apiPost('/api/v1/teacher/class/switch-series', { series_id: seriesId })
    currentSeries.value = seriesId
    toast.show(`已切换至「${getSeriesName(seriesId)}」系列`, 'success')
  } catch (e: any) {
    toast.show(e?.response?.data?.message || '切换失败', 'error')
  } finally {
    switching.value = false
  }
}

onMounted(() => {
  loadClassInfo()
})
</script>

<template>
  <SidebarLayout role-label="教师端" :nav-items="navItems">
    <template #user-meta>
      <div style="font-size:12px;color:var(--color-text-secondary);">
        {{ authStore.user?.class_names?.join('、') || '教师' }}
      </div>
    </template>

    <!-- 侧边栏底部：系列选择器（功能分类与版面设计.txt 版式） -->
    <template #sidebar-extra>
      <div class="series-selector">
        <div class="ss-label">🏷️ 班级宠物系列</div>
        <div class="ss-info">
          <span>班级积分</span>
          <span class="ss-score">{{ classTotalScore.toLocaleString() }}</span>
        </div>
        <div class="ss-select-wrap">
          <select
            v-model="currentSeries"
            class="ss-select"
            :disabled="switching"
            :style="{ '--select-accent': currentSeriesColor }"
            @change="switchSeries(currentSeries)"
          >
            <option v-for="s in allSeries" :key="s.id" :value="s.id">
              {{ s.emoji }} {{ s.name }}
            </option>
          </select>
        </div>
        <div class="ss-hint">切换系列将更新全班宠物外观</div>
      </div>
    </template>
  </SidebarLayout>
</template>

<style scoped>
.series-selector {
  font-size: 12px;
}
.ss-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
  letter-spacing: 0.03em;
}
.ss-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}
.ss-score {
  color: var(--color-primary, #a78bfa);
  font-weight: 700;
  font-size: 14px;
}
.ss-select-wrap {
  position: relative;
}
.ss-select {
  width: 100%;
  padding: 8px 10px;
  border-radius: 10px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  font-size: 13px;
  font-weight: 500;
  outline: none;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
  appearance: auto;
}
.ss-select:focus {
  border-color: var(--select-accent, var(--color-primary));
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--select-accent, var(--color-primary)) 20%, transparent);
}
.ss-select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.ss-hint {
  font-size: 10px;
  color: var(--color-text-secondary);
  opacity: 0.4;
  margin-top: 4px;
  text-align: right;
}
</style>
