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
// 该班级任教科目是否 = 个人主教科目（主教科目已在头部展示，班级区不重复）
function isMainSubject(subject: string | undefined) {
  return !!subject && subject !== '默认科目' && subject === props.teacher.subject
}
const roleLabel: Record<string, string> = {
  head_teacher: '主班主任',
  co_teacher: '副班主任',
  subject_teacher: '科任教师',
}
</script>
<template>
  <div class="teacher-card">
    <!-- 版块1：身份区 + 主教区（并排一行） -->
    <div class="tc-block tc-header">
      <div class="tc-id">
        <span class="tc-name-text">{{ teacher.name }}</span>
        <span v-if="teacher.personal_role === 'grade_lead'" class="head-badge badge-lead">首席</span>
        <span v-else-if="teacher.personal_role === 'admin_director'" class="head-badge badge-admin">主任</span>
      </div>
      <div class="tc-id-right">
        <span v-if="teacher.subject" class="tc-main-subject">
          <span class="ms-icon">📘</span>
          <span class="ms-value">{{ teacher.subject }}</span>
        </span>
        <span class="tc-no">{{ no }}</span>
      </div>
    </div>

    <!-- 版块3：班级任教区（每班一个子版块） -->
    <div class="tc-block tc-classes">
      <div class="tc-block-title">📚 任教班级</div>
      <div v-if="teacher.assignments.length === 0" class="tc-class-empty">暂未分配班级</div>
      <div v-for="a in teacher.assignments" :key="a.class_id + '_' + a.role" class="tc-class-row">
        <span class="tc-class-name">{{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}</span>
        <span class="tc-class-role">{{ roleLabel[a.role] || a.role }}</span>
        <span v-if="a.subject && a.subject !== '默认科目' && !isMainSubject(a.subject)" class="tc-class-subject">
          兼任·{{ a.subject }}
        </span>
      </div>
    </div>

    <!-- 版块4：操作区 -->
    <div class="tc-block tc-footer">
      <div class="tc-actions">
        <button class="act-btn act-icon" title="编辑教师" @click="emit('edit', teacher)">👤</button>
        <button class="act-btn act-icon" title="分配班级" @click="emit('assign', teacher)">📚</button>
        <button class="act-btn act-icon" title="重置密码" @click="emit('resetPwd', teacher)">🔑</button>
        <button class="act-btn act-del" @click="emit('delete', teacher)">🗑️ 删除</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.teacher-card {
  background: var(--color-bg-card);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 0;
  display: flex;
  flex-direction: column;
  transition: all 0.2s ease;
  overflow: hidden;
}
.teacher-card:hover {
  border-color: rgba(167, 139, 250, 0.4);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.35);
  transform: translateY(-2px);
}
/* 通用版块：左右留白 + 底部分隔线 */
.tc-block {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}
.tc-block:last-child {
  border-bottom: none;
}

/* 版块1：身份区 + 主教区（并排一行） */
.tc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}
.tc-id {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  min-width: 0;
}
.tc-id-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.tc-name-text { font-size: 18px; font-weight: 700; color: var(--color-text); line-height: 1.3; }
.tc-no { font-size: 11px; font-weight: 600; color: var(--color-text-secondary); flex-shrink: 0; letter-spacing: 0.03em; }
.head-badge { font-size: 10px; font-weight: 600; padding: 2px 8px; border-radius: 8px; white-space: nowrap; }
.badge-lead { background: rgba(139, 92, 246, 0.2); color: #c4b5fd; }
.badge-admin { background: rgba(245, 158, 11, 0.2); color: #fcd34d; }

/* 主教科目：身份区右侧紧凑胶囊 */
.tc-main-subject {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.ms-icon { font-size: 13px; }
.ms-value {
  font-size: 12px;
  font-weight: 600;
  color: #c4b5fd;
  padding: 2px 8px;
  background: rgba(139, 92, 246, 0.18);
  border: 1px solid rgba(139, 92, 246, 0.35);
  border-radius: 6px;
  white-space: nowrap;
}

/* 版块3：班级任教区 */
.tc-classes { display: flex; flex-direction: column; gap: 6px; }
.tc-block-title { font-size: 11px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 2px; }
.tc-class-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: var(--color-bg);
  border-radius: 6px;
  font-size: 12px;
}
.tc-class-name { font-weight: 600; color: var(--color-text); flex: 1; min-width: 0; }
.tc-class-role { font-size: 11px; color: var(--color-accent); flex-shrink: 0; }
.tc-class-subject { font-size: 11px; font-weight: 500; color: var(--color-text-secondary); flex-shrink: 0; }
.tc-class-empty {
  padding: 10px;
  text-align: center;
  font-size: 12px;
  color: var(--color-text-secondary);
  border: 1px dashed rgba(255, 255, 255, 0.18);
  border-radius: 6px;
}

/* 版块4：操作区 */
.tc-footer { display: flex; align-items: center; justify-content: center; }
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
</style>
