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
    <div class="toolbar-left">
      <div class="section-badge">账号管理</div>
      <h2 class="page-title">教师账号</h2>
      <span class="count-badge"><slot name="count" /></span>
    </div>
    <div class="toolbar-actions">
      <div class="filter-group">
        <select :value="filterGrade" @change="emit('update:filterGrade', ($event.target as HTMLSelectElement).value)" class="form-input filter-select">
          <option value="">-- 请选择年级 --</option>
          <option v-for="g in grades" :key="g" :value="g + '团队'">{{ g }}团队</option>
        </select>
        <select :value="filterRole" @change="emit('update:filterRole', ($event.target as HTMLSelectElement).value)" class="form-input filter-select">
          <option value="">全部角色</option>
          <option v-for="(label, role) in classRoleLabel" :key="role" :value="role">{{ label }}</option>
        </select>
        <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" class="form-input filter-search" placeholder="搜索姓名 / 账号..." />
      </div>
      <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);font-size:13px;" @click="emit('downloadTemplate')">下载模板</button>
      <button class="btn" style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);font-size:13px;" @click="emit('openImport')">批量导入</button>
      <button class="btn btn-primary" style="font-size:13px;" @click="emit('openCreate')">+ 创建教师</button>
    </div>
  </div>
</template>
