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
    <div class="tc-name">
      <span class="tc-name-text">{{ teacher.name }}</span>
      <span v-if="teacher.personal_role === 'grade_lead'" class="head-badge badge-lead">首席</span>
      <span v-else-if="teacher.personal_role === 'admin_director'" class="head-badge badge-admin">主任</span>
      <span class="tc-no">{{ no }}</span>
    </div>
    <div class="tc-classes">
      <span v-for="a in teacher.assignments" :key="a.class_id + '_' + a.role" class="tc-class">
        {{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}
        <span v-if="a.role === 'head_teacher'" class="role-tag-head">主班</span>
      </span>
      <span v-if="teacher.assignments.length === 0" class="tc-class-empty">未分配班级</span>
    </div>
    <div class="tc-actions">
      <button class="act-btn" @click="emit('edit', teacher)">👤 编辑</button>
      <button class="act-btn" @click="emit('assign', teacher)">📚 班级</button>
      <button class="act-btn" @click="emit('resetPwd', teacher)">🔑 密码</button>
      <button class="act-btn act-del" @click="emit('delete', teacher)">🗑️ 删除</button>
    </div>
  </div>
</template>
