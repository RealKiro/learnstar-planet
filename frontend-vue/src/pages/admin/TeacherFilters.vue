<script setup lang="ts">
defineProps<{
  grades: string[]
  classRoleLabel: Record<string, string>
  filterGrade: string
  filterRole: string
  searchQuery: string
}>()
const emit = defineEmits<{
  'update:filterGrade': [v: string]
  'update:filterRole': [v: string]
  'update:searchQuery': [v: string]
  'downloadTemplate': []
  'openImport': []
  'openCreate': []
}>()
</script>
<template>
  <div class="toolbar">
    <!-- 第一行：标题 -->
    <div class="toolbar-header">
      <span class="section-badge">账号管理</span>
      <h2 class="page-title">教师账号</h2>
      <span class="count-badge"><slot name="teacherCount" /></span>
    </div>
    <!-- 第二行：操作工具栏 -->
    <div class="toolbar-actions">
      <div class="filter-group">
        <select :value="filterGrade" @change="emit('update:filterGrade', ($event.target as HTMLSelectElement).value)" class="form-input">
          <option value="">-- 请选择年级 --</option>
          <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
        </select>
        <select :value="filterRole" @change="emit('update:filterRole', ($event.target as HTMLSelectElement).value)" class="form-input">
          <option value="">全部角色</option>
          <option v-for="(label, role) in classRoleLabel" :key="role" :value="role">{{ label }}</option>
        </select>
        <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" class="form-input" placeholder="搜索姓名 / 账号..." />
      </div>
      <div class="action-buttons">
        <button class="btn" @click="emit('downloadTemplate')">下载模板</button>
        <button class="btn" @click="emit('openImport')">批量导入</button>
        <button class="btn btn-primary" @click="emit('openCreate')">+ 创建教师</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px 0 16px;
  border-bottom: 1px solid var(--color-border);
}
.toolbar-header {
  display: flex;
  align-items: center;
  gap: 12px;
}
.toolbar-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}
.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.filter-group select,
.filter-group input {
  height: 36px;
  padding: 0 12px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  color: var(--color-text);
  font-size: 13px;
  outline: none;
  transition: 0.15s;
  box-sizing: border-box;
}
.filter-group select {
  min-width: 130px;
  flex: 0 0 auto;
}
.filter-group input {
  flex: 1 1 180px;
  min-width: 150px;
}
.filter-group select:focus,
.filter-group input:focus {
  border-color: #7c3aed;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.08);
}
.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}
.action-buttons .btn {
  height: 36px;
  padding: 0 16px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  color: var(--color-text);
  font-size: 13px;
  cursor: pointer;
  transition: 0.15s;
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
  font-family: inherit;
}
.action-buttons .btn-primary {
  background: #7c3aed;
  border: none;
  color: #fff;
}
.action-buttons .btn-primary:hover {
  background: #6d28d9;
}
</style>
