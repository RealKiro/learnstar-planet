import { computed } from 'vue'
import { useRoute } from 'vue-router'

export type AppMode = 'classroom' | 'teacher'

/**
 * 教师端/教室端页面模式检测（页面合并重构的核心）。
 * 教室模式 = 基础模式路由（meta.basic，班级码进入）或 classroom-* 路由（教室端大屏）。
 * 教师模式 = 教师完整模式路由（teacher-scores / teacher-dashboard / teacher-ai 等）。
 * 两个 Layout 的 <router-view> 均带 :key="$route.fullPath"，路由切换会整组件重挂载，
 * 因此该检测在 onMounted 时即已可靠，无需 watch(route)。
 */
export function useAppMode() {
  const route = useRoute()
  const isClassroomMode = computed(() => {
    const name = String(route.name || '')
    return route.meta.basic === true || name.startsWith('classroom-')
  })
  const isTeacherMode = computed(() => !isClassroomMode.value)
  const mode = computed<AppMode>(() => (isClassroomMode.value ? 'classroom' : 'teacher'))
  return { mode, isClassroomMode, isTeacherMode }
}
