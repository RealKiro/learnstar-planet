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
      <span class="count-badge"><slot name="teacherCount" /></span>
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

<style scoped>
.toolbar { display:flex; align-items:center; flex-wrap:nowrap; gap:12px; }
.toolbar-left { flex-shrink:0; display:flex; align-items:center; gap:8px; }
.toolbar-actions { display:flex; align-items:center; gap:8px 12px; flex:1 1 auto; flex-wrap:nowrap; justify-content:flex-end; }
.filter-group { display:flex; align-items:center; gap:6px; flex:1 1 auto; flex-wrap:nowrap; }
.filter-select { padding:6px 10px; font-size:12px; flex:0 1 auto; min-width:0; width:auto; }
.filter-search { padding:6px 10px; font-size:12px; flex:1 1 120px; min-width:80px; }
.filter-select option { color:#1E293B; background:#fff; }
.toolbar-actions .btn { flex-shrink:0; }
</style>
