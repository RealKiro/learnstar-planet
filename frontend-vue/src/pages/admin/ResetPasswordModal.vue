<script setup lang="ts">
import { ref, watch } from 'vue'
import { apiPost } from '@/utils/api'
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

const resetPwdValue = ref('')
const resetPwdStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const resetPwdError = ref('')
const showResetPwd = ref(false)
const currentPwd = ref('')

watch(
  () => props.teacher,
  (t) => {
    if (t) {
      resetPwdValue.value = ''
      showResetPwd.value = false
      currentPwd.value = ''
      resetPwdStatus.value = 'idle'
      resetPwdError.value = ''
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

function closeModal() {
  emit('update:visible', false)
}

async function submitResetPwd() {
  if (!props.teacher) return
  resetPwdStatus.value = 'loading'
  resetPwdError.value = ''
  try {
    const res = await apiPost<{ data: { new_password: string } }>(
      `/api/v1/admin/teachers/${props.teacher.id}/reset-password`,
      { password: resetPwdValue.value || undefined }
    )
    const newPwd = res?.data?.new_password || resetPwdValue.value || '（已设置）'
    resetPwdValue.value = newPwd
    resetPwdStatus.value = 'success'
    setTimeout(() => {
      resetPwdStatus.value = 'idle'
      emit('update:visible', false)
    }, 1200)
    emit('reset')
  } catch {
    resetPwdStatus.value = 'error'
    resetPwdError.value = '重置失败，请稍后重试'
    setTimeout(() => { if (resetPwdStatus.value === 'error') resetPwdStatus.value = 'idle' }, 3000)
  }
}

</script>

<template>
  <ModalGlass :visible="visible" @update:visible="emit('update:visible', $event)">
    <div style="max-width:420px;width:100%;padding:4px 0;">
      <div class="modal-header">
        <h3 style="font-size:16px;font-weight:700;color:var(--color-text);margin:0;">
          &#128273; 密码管理 — {{ teacher?.name }}
        </h3>
        <button @click="closeModal" style="background:none;border:none;color:var(--color-text-secondary);font-size:20px;cursor:pointer;padding:0;line-height:1;">&#10005;</button>
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
        <div class="flex-row" style="gap:6px;">
          <input
            v-model="resetPwdValue"
            :type="showResetPwd ? 'text' : 'password'"
            class="form-input flex-1"
            placeholder="留空默认 ls123456"
            autocomplete="new-password"
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

      <div class="flex-row" style="margin-bottom:12px;">
        <button
          type="button"
          class="flex-1" style="padding:6px;border-radius:6px;font-size:11px;cursor:pointer;border:1px solid var(--color-border);background:var(--color-bg-card);color:var(--color-text-secondary);font-family:inherit;"
          @click="resetPwdValue = ''"
        >
          &#128260; 重置为空
        </button>
      </div>

      <div v-if="resetPwdError" style="margin-bottom:12px;padding:8px 12px;background:rgba(239,68,68,0.08);border:1px solid rgba(239,68,68,0.2);border-radius:8px;color:#fca5a5;font-size:12px;">{{ resetPwdError }}</div>

      <div class="flex-row" style="gap:12px;">
        <button
          @click="closeModal"
          class="flex-1" style="padding:8px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);"
        >
          取消
        </button>
        <button
          @click="submitResetPwd"
          :disabled="resetPwdStatus === 'loading'"
          class="flex-1" :style="{ padding:'8px', borderRadius:'8px', fontSize:'13px', fontWeight:'600', cursor:'pointer', border:'none', color:'#fff', background: resetPwdStatus === 'loading' ? '#f59e0b' : resetPwdStatus === 'success' ? '#10b981' : resetPwdStatus === 'error' ? '#ef4444' : '#7c3aed' }"
        >
          <template v-if="resetPwdStatus === 'loading'">更新中...</template>
          <template v-else-if="resetPwdStatus === 'success'">已更新 ✓</template>
          <template v-else-if="resetPwdStatus === 'error'">更新失败 ✗</template>
          <template v-else>更新密码</template>
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
/* Modal utility classes */
.modal-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; padding-bottom:12px; border-bottom:1px solid var(--color-border); flex-shrink:0; }
.modal-footer { display:flex; gap:8px; padding-top:12px; border-top:1px solid var(--color-border); margin-top:16px; }
.modal-section-title { font-size:12px; font-weight:600; color:var(--color-text); margin-bottom:8px; }
.flex-row { display:flex; gap:8px; }
.flex-1 { flex:1; }
</style>
