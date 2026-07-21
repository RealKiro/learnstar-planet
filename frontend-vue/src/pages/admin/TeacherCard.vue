<script setup lang="ts">
import { avatarGradient } from '@/utils/constants'

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
function classById(id: number) { return props.classes.find(c => c.id === id) }
</script>
<template>
  <div class="teacher-card">
    <div class="card-head">
      <div class="avatar" :style="{ background: avatarGradient(teacher.name) }">{{ teacher.name[0] }}</div>
      <div class="head-body">
        <div class="head-top">
          <span class="head-name">{{ teacher.name }}</span>
          <span v-if="teacher.personal_role === 'grade_lead'" class="head-badge badge-lead">首席</span>
          <span v-else-if="teacher.personal_role === 'admin_director'" class="head-badge badge-admin">主任</span>
        </div>
        <div class="head-id">{{ teacher.username }}</div>
      </div>
    </div>
    <div class="card-info-section" v-if="teacher.email || teacher.phone">
      <span v-if="teacher.email" class="info-item">✉ {{ teacher.email }}</span>
      <span v-if="teacher.phone" class="info-item">📱 {{ teacher.phone }}</span>
    </div>
    <div class="card-classes">
      <div v-for="a in teacher.assignments" :key="a.class_id + '_' + a.role" class="class-row">
        <span class="class-name">{{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}</span>
        <span v-if="a.role === 'head_teacher'" class="role-tag-head">主班</span>
        <span v-else-if="a.role === 'co_teacher'" class="role-tag-co">副班</span>
        <span class="class-subj">{{ a.role === 'head_teacher' || a.role === 'co_teacher' ? a.subject : (a.subject || '—') }}</span>
      </div>
      <div v-if="teacher.assignments.length === 0" class="class-empty">暂未分配班级</div>
    </div>
    <div class="card-actions">
      <button class="act-btn" @click="emit('edit', teacher)">👤 个人信息</button>
      <button class="act-btn" @click="emit('assign', teacher)">📚 班级管理</button>
      <button class="act-btn" @click="emit('resetPwd', teacher)">🔑 密码</button>
      <button class="act-btn act-del" @click="emit('delete', teacher)">🗑️ 删除</button>
    </div>
  </div>
</template>
