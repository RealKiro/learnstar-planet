<script setup lang="ts">
import type { ScoreRule } from '@/types'

export interface RuleGroup {
  category: string
  label: string
  rules: ScoreRule[]
}

defineProps<{
  show: boolean
  type: 'add' | 'sub'
  title: string
  subtitle: string
  groups: RuleGroup[]
  busy: boolean
  /** 单学生按规则操作反馈状态（batch 时为 idle） */
  status: 'idle' | 'loading' | 'success' | 'error'
  activeRuleName: string
  error: string
}>()

const emit = defineEmits<{
  close: []
  apply: [rule: ScoreRule]
}>()
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-box">
        <h3 class="modal-title">{{ title }}</h3>
        <p class="modal-sub">{{ subtitle }}</p>
        <div class="reason-groups">
          <div v-for="group in groups" :key="group.category" class="reason-group">
            <div class="group-title">{{ group.label }}</div>
            <div class="group-btns">
              <button
                v-for="rule in group.rules"
                :key="rule.id"
                class="rule-btn"
                :class="{ 'rule-add': rule.amount > 0, 'rule-sub': rule.amount < 0 }"
                :style="status !== 'idle' && activeRuleName === rule.name ? { background: status === 'loading' ? '#f59e0b' : status === 'success' ? '#10b981' : '#ef4444', color: '#fff', borderColor: 'transparent' } : {}"
                @click="emit('apply', rule)"
                :disabled="busy"
              >
                <span class="rule-amt">{{ rule.amount > 0 ? '+' : '' }}{{ rule.amount }}</span>
                <span class="rule-name">{{ status !== 'idle' && activeRuleName === rule.name ? (status === 'loading' ? '处理中...' : status === 'success' ? '操作成功' : '操作失败') : rule.name }}</span>
              </button>
            </div>
          </div>
        </div>
        <div v-if="error" class="modal-error">{{ error }}</div>
        <button class="cancel-btn" @click="emit('close')">取消操作</button>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 300;
}
.modal-box {
  background: var(--color-bg-card);
  border-radius: 20px;
  padding: 28px 32px;
  max-width: 420px;
  width: 90%;
  box-shadow: var(--shadow-lg);
  animation: modalPop 0.25s ease;
}
@keyframes modalPop {
  from { transform: scale(0.92); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
.modal-title { font-size: 20px; font-weight: 700; margin-bottom: 6px; }
.modal-sub { font-size: 14px; color: var(--color-text-secondary); margin-bottom: 20px; }
.modal-sub strong { color: var(--color-text); }
.modal-error {
  margin-bottom: 10px;
  padding: 8px 12px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  border-radius: 8px;
  color: var(--color-danger-text);
  font-size: 12px;
}

/* 规则分组（并列式） */
.reason-groups {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-bottom: 20px;
  max-height: 56vh;
  overflow-y: auto;
  padding-right: 4px;
}
.reason-group {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  padding: 10px 12px;
}
.group-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 700;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
  letter-spacing: 0.03em;
}
.group-title::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--color-border);
}
.group-btns {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.rule-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 44px;
  padding: 8px 16px;
  border-radius: 20px;
  border: 1px solid var(--color-border);
  background: var(--color-bg-card);
  color: var(--color-text);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.rule-btn:hover {
  background: rgba(124,58,237,0.08);
  border-color: var(--color-primary);
  transform: translateY(-1px);
}
.rule-amt {
  font-size: 13px;
  font-weight: 800;
  min-width: 30px;
  text-align: center;
}
.rule-add .rule-amt { color: #10B981; }
.rule-sub .rule-amt { color: #EF4444; }
.rule-name {
  white-space: nowrap;
}
.cancel-btn {
  width: 100%;
  padding: 10px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.cancel-btn:hover { background: var(--color-bg); }

/* 过渡 */
.modal-enter-active { transition: opacity 0.2s ease; }
.modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
