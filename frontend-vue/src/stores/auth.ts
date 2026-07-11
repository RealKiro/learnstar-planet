import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)

  const stored = localStorage.getItem('currentUser')
  if (stored) {
    try { user.value = JSON.parse(stored) } catch { /* ignore */ }
  }

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'school_admin')
  const isTeacher = computed(() => user.value?.role === 'teacher')
  const isParent = computed(() => user.value?.role === 'parent')
  const userName = computed(() => user.value?.name || '')
  const displayName = computed(() => user.value?.nickname || user.value?.name || '')

  function setAuth(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('currentUser', JSON.stringify(newUser))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
  }

  return { token, user, isLoggedIn, isAdmin, isTeacher, isParent, userName, displayName, setAuth, logout }
})
