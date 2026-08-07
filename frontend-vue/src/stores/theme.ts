import { defineStore } from 'pinia'
import { ref } from 'vue'

type ThemeMode = 'system' | 'light' | 'dark'

const mql = window.matchMedia('(prefers-color-scheme: dark)')

export const useThemeStore = defineStore('theme', () => {
  // 兼容旧存储 'dark'/'light'；缺省为跟随系统
  const mode = ref<ThemeMode>((localStorage.getItem('theme') as ThemeMode) || 'system')
  const isDark = ref(false)

  function apply() {
    isDark.value = mode.value === 'dark' || (mode.value === 'system' && mql.matches)
    document.documentElement.classList.toggle('dark', isDark.value)
    document.documentElement.classList.toggle('light', !isDark.value)
  }

  function init() {
    apply()
    mql.addEventListener('change', () => {
      if (mode.value === 'system') apply()
    })
  }

  function setMode(m: ThemeMode) {
    mode.value = m
    localStorage.setItem('theme', m)
    apply()
  }

  /** 翻转当前生效色，并从 system 态跳出为显式设置 */
  function toggle() {
    setMode(isDark.value ? 'light' : 'dark')
  }

  return { mode, isDark, init, setMode, toggle }
})
