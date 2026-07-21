import { defineStore } from 'pinia'
import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'info' | 'warning'
export type ToastPosition = 'top-right' | 'center' | 'bottom-center'

export interface ToastOptions {
  message: string
  type?: ToastType
  position?: ToastPosition
  duration?: number
  action?: { label: string; onClick: () => void }
}

interface Toast extends ToastOptions {
  id: number
  type: ToastType
  position: ToastPosition
  duration: number
  createdAt: number
}

let nextId = 0

export const useToastStore = defineStore('toast', () => {
  const toasts = ref<Toast[]>([])

  function show(
    message: string,
    type: ToastType = 'success',
    options?: { position?: ToastPosition; duration?: number; action?: { label: string; onClick: () => void } }
  ) {
    const id = nextId++
    const toast: Toast = {
      id,
      message,
      type,
      position: options?.position || 'top-right',
      duration: options?.duration ?? 5000,
      createdAt: Date.now(),
      action: options?.action,
    }
    toasts.value.push(toast)
    if (toast.duration > 0) {
      setTimeout(() => remove(id), toast.duration)
    }
  }

  function remove(id: number) {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  function clear() {
    toasts.value = []
  }

  return { toasts, show, remove, clear }
})
