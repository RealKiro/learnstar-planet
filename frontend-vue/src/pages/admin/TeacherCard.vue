<script setup lang="ts">
import { computed } from 'vue'

interface Assignment { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }
interface Teacher { id: number; name: string; username: string; nickname?: string; subject?: string; grade_team?: string; phone?: string; email?: string; avatar_path?: string; status: string; bindings: string[]; assignments: Assignment[]; class_names: string[]; personal_role?: string }
interface ClassRoom { id: number; name: string; grade?: string }

const props = defineProps<{
  teacher: Teacher
  classes: ClassRoom[]
  no: string
  selectable?: boolean
  selected?: boolean
}>()
const emit = defineEmits<{
  edit: [t: Teacher]
  assign: [t: Teacher]
  resetPwd: [t: Teacher]
  delete: [t: Teacher]
  toggleSelect: [id: number]
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
// 身份行中间区：个人角色 + 主教科目，多项用 · 连接（如「首席 · 信息科技」「主任 · 语文」）
const roleMainText = computed(() => {
  const parts: string[] = []
  if (props.teacher.personal_role === 'grade_lead') parts.push('首席')
  else if (props.teacher.personal_role === 'admin_director') parts.push('主任')
  if (props.teacher.subject) parts.push(props.teacher.subject)
  return parts.join(' · ')
})
</script>
<template>
  <div class="teacher-card" :class="{ 'is-selected': selectable && selected }">
    <!-- 版块1：身份行（左姓名 / 中角色+主教 / 右编号 + 钻孔吊牌式圆形勾选） -->
    <div class="tc-block tc-header" :class="{ 'has-select': selectable }">
      <span class="tc-name-text">{{ teacher.name }}</span>
      <div v-if="roleMainText" class="tc-role-main">
        <span class="tc-role-main-pill">{{ roleMainText }}</span>
      </div>
      <div class="tc-right">
        <span class="tc-no">{{ no }}</span>
        <button
          v-if="selectable"
          type="button"
          class="tc-hole"
          :class="{ 'is-checked': selected }"
          :title="selected ? '取消选择' : '选择'"
          :aria-pressed="selected"
          @click.stop="emit('toggleSelect', teacher.id)"
        >
          <span v-if="selected" class="tc-hole-check">✓</span>
        </button>
      </div>
    </div>

    <!-- 版块3：班级任教区（每班一个子版块） -->
    <div class="tc-block tc-classes">
      <div class="tc-block-title">📚 任教班级</div>
      <div v-if="teacher.assignments.length === 0" class="tc-class-empty">暂未分配班级</div>
      <div v-for="a in teacher.assignments" :key="a.class_id + '_' + a.role" class="tc-class-row">
        <span class="tc-class-name">{{ a.class_name || classById(a.class_id)?.name || '#' + a.class_id }}</span>
        <span v-if="a.subject && a.subject !== '默认科目' && !isMainSubject(a.subject)" class="tc-class-subject">
          {{ a.subject }}
        </span>
      </div>
    </div>

    <!-- 版块4：操作区 -->
    <div class="tc-block tc-footer">
      <div class="tc-actions">
        <button class="act-btn" title="编辑教师" @click="emit('edit', teacher)">👤 编辑</button>
        <button class="act-btn" title="分配班级" @click="emit('assign', teacher)">📚 班级</button>
        <button class="act-btn" title="重置密码" @click="emit('resetPwd', teacher)">🔑 密码</button>
        <button class="act-btn act-del" @click="emit('delete', teacher)">🗑️ 删除</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.teacher-card {
  position: relative;
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
.teacher-card.is-selected {
  border-color: rgba(124, 58, 237, 0.55);
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.2);
}
/* 钻孔吊牌式圆形勾选：身份行右端，与编号垂直居中，像给卡片钻了个孔 */
.tc-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.tc-hole {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1.5px solid rgba(255, 255, 255, 0.22);
  background: transparent;
  /* 内阴影模拟卡片厚度，呈现"钻孔"立体感 */
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.35), 0 1px 2px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  font-family: inherit;
  transition: all 0.2s ease;
}
.tc-hole:hover {
  border-color: rgba(167, 139, 250, 0.7);
  transform: scale(1.1);
}
.tc-hole.is-checked {
  border-color: #a78bfa;
  background: radial-gradient(circle at 35% 30%, #a78bfa, #7c3aed);
  box-shadow: 0 0 8px rgba(124, 58, 237, 0.5), inset 0 1px 2px rgba(0, 0, 0, 0.25);
}
.tc-hole-check {
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  line-height: 1;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}
/* 通用版块：左右留白 + 底部分隔线 */
.tc-block {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}
.tc-block:last-child {
  border-bottom: none;
}

/* 版块1：身份行三栏（左姓名 / 中角色+主教 / 右编号） */
.tc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.tc-name-text {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text);
  line-height: 1.3;
  flex-shrink: 0;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tc-role-main {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-width: 0;
}
.tc-role-main-pill {
  font-size: 12px;
  font-weight: 600;
  color: #c4b5fd;
  padding: 2px 10px;
  background: rgba(139, 92, 246, 0.18);
  border: 1px solid rgba(139, 92, 246, 0.35);
  border-radius: 999px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}
.tc-no { font-size: 11px; font-weight: 600; color: var(--color-text-secondary); flex-shrink: 0; letter-spacing: 0.03em; }

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
.tc-class-subject { font-size: 12px; font-weight: 500; color: var(--color-text-secondary); flex-shrink: 0; }
.tc-class-empty {
  padding: 10px;
  text-align: center;
  font-size: 12px;
  color: var(--color-text-secondary);
  border: 1px dashed rgba(255, 255, 255, 0.18);
  border-radius: 6px;
}

/* 版块4：操作区 */
.tc-footer { display: flex; align-items: center; }
.tc-actions { display: flex; gap: 6px; width: 100%; }
.act-btn {
  flex: 1;
  display: inline-flex; align-items: center; justify-content: center; gap: 4px; height: 30px;
  padding: 0 6px; border-radius: 6px; font-size: 12px; font-weight: 500; cursor: pointer;
  transition: 0.15s; border: 1px solid rgba(255, 255, 255, 0.08); background: rgba(255, 255, 255, 0.04);
  color: var(--color-text); font-family: inherit; white-space: nowrap;
}
.act-btn:hover { background: rgba(255, 255, 255, 0.08); border-color: rgba(255, 255, 255, 0.14); }
.act-del { color: #fca5a5; border-color: rgba(239, 68, 68, 0.25); }
.act-del:hover { background: rgba(239, 68, 68, 0.12); color: #fecaca; }
</style>
