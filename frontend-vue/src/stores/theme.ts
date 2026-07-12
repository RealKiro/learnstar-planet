import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(localStorage.getItem('theme') === 'dark')

  function toggle() {
    isDark.value = !isDark.value
    const theme = isDark.value ? 'dark' : 'light'
    localStorage.setItem('theme', theme)
    document.documentElement.className = isDark.value ? 'dark' : ''
  }

  function init() {
    // 优先使用用户设置，其次默认浅色主题
    const stored = localStorage.getItem('theme')
    if (stored) {
      isDark.value = stored === 'dark'
    } else {
      isDark.value = false
    }
    document.documentElement.className = isDark.value ? 'dark' : ''
  }

  return { isDark, toggle, init }
})
