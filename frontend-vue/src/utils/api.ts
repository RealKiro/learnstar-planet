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
      const authStore = useAuthStore()
      authStore.logout()
      toast.show('登录已过期，请重新登录', 'error')
      // 跳转登录页，用 location.href 确保状态完全重置
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    } else if (error.response?.data?.message) {
      toast.show(error.response.data.message, 'error')
    } else if (error.message === 'Network Error') {
      toast.show('网络错误，请稍后重试', 'error')
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
