import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

const instance: AxiosInstance = axios.create({
  baseURL: '',
  headers: {
    'Accept': 'application/json',
  },
})

// 请求拦截器 — 自动附加 token
instance.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

// 响应拦截器 — 统一处理 401 和错误
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    const toast = useToastStore()

    if (error.response?.status === 401) {
      // 班级大屏显示端有自己的 Token 管理，不触发 auth store 的 401 处理
      const isDisplayPath = (error.config?.url || '').startsWith('/api/v1/display/')
      if (isDisplayPath) {
        return Promise.reject(error)
      }
      // 班级码登录的 401 = 班级码无效，由登录页自行提示具体原因，不当作会话过期
      const isClassLogin = (error.config?.url || '').includes('/api/v1/auth/class/login')
      if (isClassLogin) {
        return Promise.reject(error)
      }
      const authStore = useAuthStore()
      // 基础模式（班级码登录，无教师 token）：接口 401 不当作会话过期，静默拒绝
      if (!authStore.token && sessionStorage.getItem('class_token')) {
        return Promise.reject(error)
      }
      authStore.logout()
      toast.show('登录已过期，请重新登录', 'error', { position: 'center', duration: 3000 })
      // 跳转登录页，用 location.href 确保状态完全重置
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    } else if (error.response?.data?.message) {
      toast.show(error.response.data.message, 'error', { position: 'center', duration: 3000 })
    } else if (error.message === 'Network Error') {
      toast.show('网络错误，请稍后重试', 'error', { position: 'center', duration: 3000 })
    }

    return Promise.reject(error)
  },
)

export default instance

// 便捷请求方法
export async function apiGet<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
  const res = await instance.get<T>(url, config)
  return res.data
}

export async function apiPost<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  const res = await instance.post<T>(url, data, config)
  return res.data
}

export async function apiPut<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  const res = await instance.put<T>(url, data, config)
  return res.data
}

export async function apiDelete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
  const res = await instance.delete<T>(url, config)
  return res.data
}
