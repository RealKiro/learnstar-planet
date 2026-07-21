<script setup lang="ts">
import { ref } from 'vue'
import { apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import ModalGlass from '@/components/common/ModalGlass.vue'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'imported'): void
}>()

const toast = useToastStore()

const importFile = ref<File | null>(null)
const importPreview = ref<any[]>([])
const importLoading = ref(false)

function close() {
  emit('update:visible', false)
}

function onFileChange(e: Event) {
  importFile.value = (e.target as HTMLInputElement).files?.[0] || null
  importPreview.value = []
}

function downloadTemplate() {
  window.open('/api/v1/admin/teachers/template-csv', '_blank')
}

async function uploadImport(isDry: boolean) {
  if (!importFile.value) {
    toast.show('请选择文件', 'error', { position: 'center', duration: 2000 })
    return
  }
  importLoading.value = true
  try {
    const fd = new FormData()
    fd.append('file', importFile.value)
    fd.append('dry_run', isDry ? '1' : '0')

    const res = await apiPost<{ preview?: any[]; created?: any[]; total: number; message: string }>(
      '/api/v1/admin/teachers/import',
      fd
    )

    if (isDry && res.preview) {
      importPreview.value = res.preview
      toast.show('预览：' + res.total + ' 条数据', undefined, { position: 'center', duration: 2000 })
    } else {
      toast.show(res.message || '导入完成', 'success', { position: 'center', duration: 2000 })
      emit('update:visible', false)
      emit('imported')
    }
  } finally {
    importLoading.value = false
  }
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:700px;width:100%;padding:4px 0;">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
        <h3 style="font-size:17px;font-weight:700;color:var(--color-text);margin:0;">批量导入教师</h3>
        <button
          class="close-btn"
          @click="close"
          style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;"
        >
          &#10005;
        </button>
      </div>

      <div class="modal-body">
        <div style="margin-bottom:16px;">
          <p style="font-size:13px;color:#6b7280;margin-bottom:8px;">支持 CSV / Excel 格式。模板列：</p>
          <div class="column-hint">
            <code>姓名</code><span style="color:#dc2626;">*</span>
            <code>年级团队</code>
            <code>科目</code>
            <code>密码</code>
            <code>手机号</code>
          </div>
          <p style="font-size:12px;color:#9ca3af;">
            密码选填，不填默认为 star123456。角色和班级导入后使用 &#127979; 按钮分配。
          </p>
          <a
            href="/api/v1/admin/teachers/template-csv"
            target="_blank"
            style="font-size:13px;color:#7c3aed;text-decoration:underline;"
          >
            下载正确模板
          </a>
          <span style="margin-left:16px;">
            <a
              href="#"
              style="font-size:13px;color:#7c3aed;text-decoration:underline;"
              @click.prevent="downloadTemplate"
            >
              下载模板
            </a>
          </span>
        </div>

        <input
          type="file"
          accept=".csv,.xlsx,.xls"
          @change="onFileChange"
          style="margin-bottom:12px;"
        >

        <div v-if="importPreview.length > 0" style="margin-bottom:12px;">
          <div style="font-weight:600;font-size:14px;margin-bottom:8px;">预览 ({{ importPreview.length }} 条)</div>
          <div
            v-if="importPreview.some(r => !r.subject)"
            style="font-size:12px;color:#d97706;padding:6px 10px;background:rgba(245,158,11,0.08);border-radius:6px;margin-bottom:8px;"
          >
            &#9888;&#65039; {{ importPreview.filter(r => !r.subject).length }} 条记录缺少科目，导入后需手动分配科目
          </div>
          <div class="preview-table-wrapper">
            <table class="preview-table">
              <thead>
                <tr>
                  <th>姓名</th>
                  <th>年级团队</th>
                  <th>科目</th>
                  <th>手机号</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, i) in importPreview.slice(0, 20)" :key="i">
                  <td>{{ row.name }}</td>
                  <td>{{ row.grade_team }}</td>
                  <td>{{ row.subject }}</td>
                  <td>{{ row.phone }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p v-if="importPreview.length > 20" style="font-size:12px;color:#9ca3af;margin-top:4px;">仅显示前 20 条</p>
        </div>
      </div>

      <div class="modal-footer" style="display:flex;justify-content:flex-end;gap:8px;padding-top:12px;border-top:1px solid var(--color-border);">
        <button
          class="btn"
          style="background:var(--color-bg);color:var(--color-text);border:1px solid var(--color-border);"
          @click="close"
        >
          取消
        </button>
        <button
          v-if="importPreview.length === 0"
          class="btn btn-primary"
          :disabled="importLoading || !importFile"
          @click="uploadImport(true)"
        >
          {{ importLoading ? '解析中...' : '预览数据' }}
        </button>
        <button
          v-if="importPreview.length > 0"
          class="btn"
          style="background:#fef3c7;color:#92400e;border:1px solid #fcd34d;"
          :disabled="importLoading"
          @click="uploadImport(true)"
        >
          重新预览
        </button>
        <button
          v-if="importPreview.length > 0"
          class="btn btn-primary"
          :disabled="importLoading"
          @click="uploadImport(false)"
        >
          {{ importLoading ? '导入中...' : '确认导入' }}
        </button>
      </div>
    </div>
  </ModalGlass>
</template>

<style scoped>
.modal-body {
  margin-bottom: 20px;
}
.close-btn {
  background: none;
  border: none;
  font-size: 22px;
  cursor: pointer;
  color: #9ca3af;
  padding: 0;
  line-height: 1;
}
.column-hint {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 8px;
}
.column-hint code {
  font-size: 11px;
  background: var(--color-bg);
  padding: 2px 6px;
  border-radius: 4px;
}
.preview-table-wrapper {
  max-height: 300px;
  overflow: auto;
}
.preview-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}
.preview-table th {
  background: #f9fafb;
  padding: 6px 8px;
  text-align: left;
  border: 1px solid var(--color-border);
  font-size: 11px;
  white-space: nowrap;
}
.preview-table td {
  padding: 5px 8px;
  border: 1px solid var(--color-border);
  white-space: nowrap;
}
.btn {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.15s;
}
.btn-primary {
  background: #7c3aed;
  color: white;
  border-color: #7c3aed;
}
.btn-primary:hover {
  background: #6d28d9;
}
</style>
