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
      // 跳转登录页并用 expired 参数标记，登录页内联提示"登录已过期"，避免全局 toast
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login?expired=1'
      }
    } else if (error.response?.data?.message && !error.config?.skipToast) {
      // 组件已用内联提示的请求（skipToast: true）不再弹全局 toast
      toast.show(error.response.data.message, 'error', { position: 'center', duration: 3000 })
    } else if (error.message === 'Network Error') {
      toast.show('网络错误，请稍后重试', 'error', { position: 'center', duration: 3000 })
    }

    return Promise.reject(error)
  },
)

export default instance

// 扩展请求配置：组件已内联提示错误时传 skipToast: true，跳过全局 toast
export interface ApiRequestConfig extends AxiosRequestConfig {
  skipToast?: boolean
}

// 便捷请求方法
export async function apiGet<T = unknown>(url: string, config?: ApiRequestConfig): Promise<T> {
  const res = await instance.get<T>(url, config)
  return res.data
}

export async function apiPost<T = unknown>(url: string, data?: unknown, config?: ApiRequestConfig): Promise<T> {
  const res = await instance.post<T>(url, data, config)
  return res.data
}

export async function apiPut<T = unknown>(url: string, data?: unknown, config?: ApiRequestConfig): Promise<T> {
  const res = await instance.put<T>(url, data, config)
  return res.data
}

export async function apiDelete<T = unknown>(url: string, config?: ApiRequestConfig): Promise<T> {
  const res = await instance.delete<T>(url, config)
  return res.data
}
