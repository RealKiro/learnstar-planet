<script setup lang="ts">
import { computed } from 'vue'

interface Assignment { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }
interface Teacher { id: number; name: string; username: string; nickname?: string; subject?: string; grade_team?: string; phone?: string; email?: string; avatar_path?: string; status: string; bindings: string[]; assignments: Assignment[]; class_names: string[]; personal_role?: string }
interface ClassRoom { id: number; name: string; grade?: string }

const props = defineProps<{
  teacher: Teacher
  classes: ClassRoom[]
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
    <div class="card-header">
      <div class="card-title">
        <span class="teacher-name">{{ teacher.name }}</span>
        <span v-if="teacher.personal_role === 'grade_lead'" class="head-badge badge-lead">首席</span>
        <span v-else-if="teacher.personal_role === 'admin_director'" class="head-badge badge-admin">主任</span>
      </div>
      <span class="teacher-username">{{ teacher.username }}</span>
    </div>
    <div class="card-info" v-if="teacher.email || teacher.phone">
      <span v-if="teacher.email" class="info-item">✉ {{ teacher.email }}</span>
      <span v-if="teacher.phone" class="info-item">📱 {{ teacher.phone }}</span>
    </div>
    <div class="card-classes">
      <div v-for="a in teacher.assignments" :key="a.class_id + '_' + a.role" class="class-tag">
        <span class="class-name">{{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}</span>
        <span v-if="a.role === 'head_teacher'" class="role-tag-head">主班</span>
        <span v-else-if="a.role === 'co_teacher'" class="role-tag-co">副班</span>
        <span class="class-subj">{{ a.role === 'head_teacher' || a.role === 'co_teacher' ? a.subject : (a.subject || '—') }}</span>
      </div>
      <div v-if="teacher.assignments.length === 0" class="class-empty">暂未分配班级</div>
    </div>
    <div class="card-actions">
      <button class="act-btn" title="个人信息" @click="emit('edit', teacher)">👤</button>
      <button class="act-btn" title="班级管理" @click="emit('assign', teacher)">📚</button>
      <button class="act-btn" title="重置密码" @click="emit('resetPwd', teacher)">🔑</button>
      <button class="act-btn act-del" title="删除" @click="emit('delete', teacher)">🗑️</button>
    </div>
  </div>
</template>
