<script setup lang="ts">
import { computed } from 'vue'

interface Assignment { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }
interface Teacher { id: number; name: string; username: string; nickname?: string; subject?: string; grade_team?: string; phone?: string; email?: string; avatar_path?: string; status: string; bindings: string[]; assignments: Assignment[]; class_names: string[]; personal_role?: string }
interface ClassRoom { id: number; name: string; grade?: string }

const props = defineProps<{
  teacher: Teacher
  classes: ClassRoom[]
  no: string
}>()
const emit = defineEmits<{
  edit: [t: Teacher]
  assign: [t: Teacher]
  resetPwd: [t: Teacher]
  delete: [t: Teacher]
}>()
const classMap = computed(() => {
  const map: Record<number, ClassRoom> = {}
  for (const c of props.classes) map[c.id] = c
  return map
})
function classById(id: number) { return classMap.value[id] }
</script>
<template>
  <div class="teacher-card">
    <div class="tc-header">
      <span class="tc-name-text">{{ teacher.name }}</span>
      <span v-if="teacher.personal_role === 'grade_lead'" class="head-badge badge-lead">首席</span>
      <span v-else-if="teacher.personal_role === 'admin_director'" class="head-badge badge-admin">主任</span>
    </div>
    <div class="tc-classes">
      <span v-for="a in teacher.assignments" :key="a.class_id + '_' + a.role" class="tc-class">
        {{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}
        <span v-if="a.role === 'head_teacher'" class="role-tag-head">主班</span>
      </span>
      <span v-if="teacher.assignments.length === 0" class="tc-class-empty">未分配班级</span>
    </div>
    <div class="tc-footer">
      <div class="tc-actions">
        <button class="act-btn act-icon" title="编辑教师" @click="emit('edit', teacher)">👤</button>
        <button class="act-btn act-icon" title="分配班级" @click="emit('assign', teacher)">📚</button>
        <button class="act-btn act-icon" title="重置密码" @click="emit('resetPwd', teacher)">🔑</button>
        <button class="act-btn act-del" @click="emit('delete', teacher)">🗑️ 删除</button>
      </div>
      <span class="tc-no">{{ no }}</span>
    </div>
  </div>
</template>

<style scoped>
.teacher-card {
  background: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  transition: all 0.2s ease;
}
.teacher-card:hover {
  border-color: rgba(167, 139, 250, 0.4);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.35);
  transform: translateY(-2px);
}
.tc-header { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; margin-bottom: 8px; }
.tc-name-text { font-size: 18px; font-weight: 700; color: var(--color-text); line-height: 1.3; }
.head-badge { font-size: 10px; font-weight: 600; padding: 2px 8px; border-radius: 8px; white-space: nowrap; }
.badge-lead { background: rgba(139, 92, 246, 0.2); color: #c4b5fd; }
.badge-admin { background: rgba(245, 158, 11, 0.2); color: #fcd34d; }
.tc-classes { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 12px; min-height: 24px; }
.tc-class {
  display: inline-flex; align-items: center; gap: 4px; height: 24px; padding: 0 10px;
  background: rgba(255, 255, 255, 0.08); border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 6px; font-size: 12px; color: var(--color-text); box-sizing: border-box;
}
.role-tag-head { font-size: 10px; font-weight: 600; padding: 0 6px; border-radius: 4px; background: rgba(250, 204, 21, 0.15); color: #fcd34d; }
.tc-class-empty {
  display: inline-flex; align-items: center; height: 24px; padding: 0 10px; box-sizing: border-box;
  font-size: 12px; color: var(--color-text-secondary);
  border: 1px dashed rgba(255, 255, 255, 0.18); border-radius: 6px; background: transparent;
}
.tc-footer { display: flex; align-items: center; justify-content: space-between; border-top: 1px solid rgba(255, 255, 255, 0.06); padding-top: 10px; }
.tc-actions { display: flex; gap: 6px; }
.act-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 4px; height: 28px;
  padding: 0 8px; border-radius: 6px; font-size: 12px; font-weight: 500; cursor: pointer;
  transition: 0.15s; border: 1px solid rgba(255, 255, 255, 0.08); background: rgba(255, 255, 255, 0.04);
  color: var(--color-text); font-family: inherit;
}
.act-btn:hover { background: rgba(255, 255, 255, 0.08); border-color: rgba(255, 255, 255, 0.14); }
.act-icon { width: 28px; padding: 0; font-size: 14px; }
.act-del { color: #fca5a5; border-color: rgba(239, 68, 68, 0.25); }
.act-del:hover { background: rgba(239, 68, 68, 0.12); color: #fecaca; }
.tc-no { font-size: 11px; font-weight: 600; color: var(--color-text-secondary); flex-shrink: 0; letter-spacing: 0.03em; }
</style>
