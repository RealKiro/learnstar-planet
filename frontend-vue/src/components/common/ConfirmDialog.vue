<script setup lang="ts">
import { reactive } from 'vue'

// 全局单例确认弹窗状态（模块级，可在任意组件 import 调用）
export const confirmState = reactive({
  visible: false,
  title: '确认操作',
  message: '',
  danger: false,
  confirmText: '确认',
  cancelText: '取消',
  onConfirm: null as (() => void) | null,
  onCancel: null as (() => void) | null,
  loading: false,
})

/**
 * 打开确认弹窗。返回是否确认（Promise<boolean>）。
 * 用法：const ok = await openConfirm({ title, message, danger, confirmText })
 */
export function openConfirm(opts: {
  title?: string
  message: string
  danger?: boolean
  confirmText?: string
  cancelText?: string
}): Promise<boolean> {
  return new Promise((resolve) => {
    confirmState.title = opts.title || '确认操作'
    confirmState.message = opts.message
    confirmState.danger = opts.danger ?? false
    confirmState.confirmText = opts.confirmText || '确认'
    confirmState.cancelText = opts.cancelText || '取消'
    confirmState.loading = false
    confirmState.visible = true
    confirmState.onConfirm = () => {
      resolve(true)
      confirmState.visible = false
    }
    confirmState.onCancel = () => {
      resolve(false)
      confirmState.visible = false
    }
  })
}

function doConfirm() {
  if (confirmState.loading) return
  confirmState.onConfirm?.()
}
function doCancel() {
  confirmState.onCancel?.()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="confirm-fade">
      <div v-if="confirmState.visible" class="confirm-overlay" @click.self="doCancel">
        <div class="confirm-box">
          <div class="confirm-header">
            <h3 class="confirm-title">{{ confirmState.title }}</h3>
            <button class="confirm-close" @click="doCancel">✕</button>
          </div>
          <div class="confirm-body">
            <div class="confirm-icon" :class="{ danger: confirmState.danger }">
              {{ confirmState.danger ? '⚠️' : '❓' }}
            </div>
            <p class="confirm-message">{{ confirmState.message }}</p>
          </div>
          <div class="confirm-footer">
            <button class="cf-btn cf-cancel" @click="doCancel">{{ confirmState.cancelText }}</button>
            <button class="cf-btn cf-ok" :class="{ danger: confirmState.danger }" @click="doConfirm">{{ confirmState.confirmText }}</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}
.confirm-box {
  background: var(--color-bg-card, #1a1a2e);
  border: 1px solid var(--color-border, rgba(255,255,255,0.08));
  border-radius: 16px;
  width: 92%;
  max-width: 380px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  overflow: hidden;
}
.confirm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 0;
}
.confirm-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text, #f1f1f1);
  margin: 0;
}
.confirm-close {
  background: none;
  border: none;
  color: var(--color-text-secondary, rgba(255,255,255,0.6));
  font-size: 16px;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
}
.confirm-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px;
  text-align: center;
}
.confirm-icon {
  font-size: 36px;
}
.confirm-icon.danger {
  font-size: 40px;
}
.confirm-message {
  font-size: 14px;
  color: var(--color-text, #f1f1f1);
  line-height: 1.6;
  margin: 0;
  white-space: pre-line;
}
.confirm-footer {
  display: flex;
  gap: 10px;
  padding: 0 20px 20px;
}
.cf-btn {
  flex: 1;
  padding: 10px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  font-family: inherit;
  transition: all 0.15s;
}
.cf-cancel {
  background: var(--color-bg, #0e0e1a);
  border: 1px solid var(--color-border, rgba(255,255,255,0.08));
  color: var(--color-text, #f1f1f1);
}
.cf-cancel:hover { background: rgba(255,255,255,0.06); }
.cf-ok {
  background: #7c3aed;
  color: #fff;
}
.cf-ok:hover { background: #6d28d9; }
.cf-ok.danger {
  background: #dc2626;
}
.cf-ok.danger:hover { background: #b91c1c; }
.confirm-fade-enter-active, .confirm-fade-leave-active { transition: opacity 0.2s ease; }
.confirm-fade-enter-from, .confirm-fade-leave-to { opacity: 0; }
.confirm-fade-enter-active .confirm-box { animation: confirmPop 0.25s cubic-bezier(0.34, 1.56, 0.64, 1); }
.confirm-fade-leave-active .confirm-box { animation: confirmPopOut 0.2s ease forwards; }
@keyframes confirmPop { 0% { opacity: 0; transform: scale(0.9) translateY(10px); } 100% { opacity: 1; transform: scale(1) translateY(0); } }
@keyframes confirmPopOut { 0% { opacity: 1; transform: scale(1); } 100% { opacity: 0; transform: scale(0.95); } }
</style>
