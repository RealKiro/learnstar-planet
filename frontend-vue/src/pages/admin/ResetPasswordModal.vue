<script setup lang="ts">
import { ref, watch } from 'vue'
import { apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import ModalGlass from '@/components/common/ModalGlass.vue'

interface Teacher {
  id: number
  name: string
  username: string
  nickname?: string
  subject?: string
  grade_team?: string
  phone?: string
  email?: string
  avatar_path?: string
  status: string
  bindings: string[]
  assignments: { class_id: number; class_name?: string; grade?: string; role: string; subject?: string }[]
  class_names: string[]
  personal_role?: string
}

const props = defineProps<{
  visible: boolean
  teacher: Teacher | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'reset'): void
}>()

const toast = useToastStore()

const resetPwdValue = ref('')
const resetPwdLoading = ref(false)
const showResetPwd = ref(false)
const currentPwd = ref('')

watch(
  () => props.teacher,
  (t) => {
    if (t) {
      resetPwdValue.value = ''
      showResetPwd.value = false
      currentPwd.value = ''
      // 加载当前密码
      fetch('/api/v1/admin/teachers/' + t.id + '/password', {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
      })
        .then(r => r.json())
        .then(d => { currentPwd.value = d.data?.password || '' })
        .catch(() => {})
    }
  },
  { immediate: true }
)

function close() {
  emit('update:visible', false)
}

async function submitResetPwd() {
  if (!props.teacher) return
  resetPwdLoading.value = true
  try {
    const res = await apiPost<{ data: { new_password: string } }>(
      `/api/v1/admin/teachers/${props.teacher.id}/reset-password`,
      { password: resetPwdValue.value || undefined }
    )
    const newPwd = res?.data?.new_password || resetPwdValue.value || '（已设置）'
    toast.show(`密码已更新: ${newPwd}`, 'success', { position: 'center', duration: 2500 })
    emit('update:visible', false)
    emit('reset')
  } catch {
    toast.show('重置失败', 'error', { position: 'center', duration: 2000 })
  } finally {
    resetPwdLoading.value = false
  }
}

function generateStrongPassword() {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*'
  resetPwdValue.value = Array.from({ length: 12 }, () => chars[Math.floor(Math.random() * chars.length)]).join('')
  showResetPwd.value = true
}
</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:420px;width:100%;padding:4px 0;">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid var(--color-border);">
        <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">
          &#128273; 密码管理 — {{ teacher?.name }}
        </h3>
        <button @click="close" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">&#10005;</button>
      </div>

      <div
        v-if="currentPwd"
        style="margin-bottom:12px;padding:10px 12px;background:var(--color-bg);border-radius:8px;border:1px solid var(--color-border);"
      >
        <div style="font-size:11px;color:var(--color-text-secondary);margin-bottom:4px;">当前密码</div>
        <div style="display:flex;align-items:center;gap:8px;">
          <code style="font-size:14px;font-weight:700;color:var(--color-text);flex:1;font-family:monospace;">
            {{ showResetPwd ? currentPwd : '••••••••' }}
          </code>
          <button
            type="button"
            style="flex-shrink:0;padding:4px 10px;border-radius:6px;border:1px solid var(--color-border);background:var(--color-bg-card);cursor:pointer;font-size:12px;"
            @click="showResetPwd = !showResetPwd"
          >
            {{ showResetPwd ? '&#128584; 隐藏' : '&#128065;&#65039; 显示' }}
          </button>
        </div>
      </div>

      <div class="form-group">
        <label>新密码</label>
        <div style="display:flex;gap:6px;">
          <input
            v-model="resetPwdValue"
            :type="showResetPwd ? 'text' : 'password'"
            class="form-input"
            placeholder="留空自动生成"
            autocomplete="new-password"
            style="flex:1;"
          >
          <button
            type="button"
            style="flex-shrink:0;padding:6px 10px;border-radius:6px;border:1px solid var(--color-border);background:var(--color-bg-card);cursor:pointer;font-size:12px;"
            @click="showResetPwd = !showResetPwd"
          >
            {{ showResetPwd ? '&#128584;' : '&#128065;&#65039;' }}
          </button>
        </div>
      </div>

      <div style="display:flex;gap:8px;margin-bottom:12px;">
        <button
          type="button"
          style="flex:1;padding:6px;border-radius:6px;font-size:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text);font-family:inherit;"
          @click="generateStrongPassword"
        >
          &#10024; 生成强密码
        </button>
        <button
          type="button"
          style="flex:1;padding:6px;border-radius:6px;font-size:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text-secondary);font-family:inherit;"
          @click="resetPwdValue = ''"
        >
          &#128260; 重置为空
        </button>
      </div>

      <div style="display:flex;gap:12px;">
        <button
          @click="close"
          style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);"
        >
          取消
        </button>
        <button
          @click="submitResetPwd"
          :disabled="resetPwdLoading"
          style="flex:1;padding:8px;border-radius:8px;font-size:13px;font-weight:600;cursor:pointer;background:#7c3aed;border:none;color:#fff;"
        >
          {{ resetPwdLoading ? '更新中...' : '更新密码' }}
        </button>
      </div>
    </div>
  </ModalGlass>
</template>

<style scoped>
.form-group {
  margin-bottom: 12px;
}
.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4px;
}
.form-input {
  color: var(--color-text);
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
  box-sizing: border-box;
}
.form-input:focus {
  border-color: #7c3aed;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.08);
}
</style>
