import { ref, reactive, onUnmounted } from 'vue'
import { apiGet } from '@/utils/api'

export interface DisplayEvent {
  id: number
  type: 'score_update' | 'broadcast' | 'notice' | 'pet_update' | 'refresh' | 'heartbeat' | 'error' | 'reconnect'
  data: Record<string, any>
}

export interface ScoreUpdateData {
  student_id: number
  student_name: string
  amount: number
  reason: string
  total_score: number
  pet_level?: number
  pet_experience?: number
  pet_mood?: number
  is_spend?: boolean
}

export interface BroadcastData {
  id: number
  type: 'banner' | 'popup' | 'fullscreen'
  content: string
  display_seconds: number
  voice_enabled: boolean
}

export interface NoticeData {
  id: number
  title: string
  content: string
  type: string
  published_at: string
}

export interface PetUpdateData {
  student_id: number
  student_name: string
  type: 'feed' | 'evolve' | 'rename'
  mood?: number
  level?: number
  experience?: number
}

export interface DisplayState {
  /** SSE 是否已连接 */
  connected: boolean
  /** 是否正在使用轮询降级 */
  polling: boolean
  /** 连接错误次数 */
  errorCount: number
  /** 最后收到事件的时间 */
  lastEventTime: string | null
  /** 当前服务器时间（来自心跳） */
  serverTime: string | null
}

/**
 * useDisplaySSE — 班级大屏 SSE 连接管理
 *
 * 功能：
 * - 自动连接 /api/v1/display/sse
 * - 事件分发（score_update / broadcast / notice / pet_update）
 * - 心跳检测
 * - 自动重连（最多 3 次后降级为轮询）
 * - 轮询降级 /api/v1/display/poll
 * - 组件卸载时清理连接
 */
export function useDisplaySSE(token: string, classId: number) {
  const state = ref<DisplayState>({
    connected: false,
    polling: false,
    errorCount: 0,
    lastEventTime: null,
    serverTime: null,
  })

  const scoreUpdates = ref<ScoreUpdateData[]>([])
  const broadcasts = ref<BroadcastData[]>([])
  const notices = ref<NoticeData[]>([])
  const petUpdates = ref<PetUpdateData[]>([])

  let eventSource: EventSource | null = null
  let pollTimer: ReturnType<typeof setInterval> | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let isDestroyed = false

  // ---- SSE 连接 ----

  function connect() {
    if (isDestroyed) return
    closeConnection()

    const baseUrl = import.meta.env.VITE_API_BASE || ''
    const url = `${baseUrl}/api/v1/display/sse?token=${encodeURIComponent(token)}`
    state.value.polling = false

    try {
      eventSource = new EventSource(url)

      eventSource.addEventListener('score_update', (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data) as ScoreUpdateData
          scoreUpdates.value.push(data)
          state.value.lastEventTime = new Date().toISOString()
        } catch { /* ignore parse errors */ }
      })

      eventSource.addEventListener('broadcast', (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data) as BroadcastData
          broadcasts.value.push(data)
          state.value.lastEventTime = new Date().toISOString()
        } catch { /* ignore */ }
      })

      eventSource.addEventListener('notice', (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data) as NoticeData
          notices.value.push(data)
          state.value.lastEventTime = new Date().toISOString()
        } catch { /* ignore */ }
      })

      eventSource.addEventListener('pet_update', (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data) as PetUpdateData
          petUpdates.value.push(data)
          state.value.lastEventTime = new Date().toISOString()
        } catch { /* ignore */ }
      })

      eventSource.addEventListener('refresh', () => {
        // 班级码被刷新，触发页面重载
        window.location.reload()
      })

      eventSource.addEventListener('heartbeat', (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data)
          state.value.serverTime = data.time
        } catch { /* ignore */ }
        state.value.connected = true
      })

      eventSource.addEventListener('error', (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data)
          if (data.code === 401) {
            // Token 失效，跳回登录页
            onTokenExpired()
          }
        } catch { /* ignore */ }
      })

      eventSource.onopen = () => {
        state.value.connected = true
        state.value.errorCount = 0
      }

      eventSource.onerror = () => {
        state.value.connected = false
        state.value.errorCount++

        if (state.value.errorCount >= 3) {
          // 降级为轮询
          closeConnection()
          startPolling()
        }
      }
    } catch {
      state.value.errorCount++
      if (state.value.errorCount >= 3) {
        startPolling()
      } else {
        scheduleReconnect()
      }
    }
  }

  // ---- 轮询降级 ----

  let lastEventId = 0

  function startPolling() {
    if (isDestroyed || state.value.polling) return
    state.value.polling = true
    state.value.connected = false

    poll()
    pollTimer = setInterval(() => poll(), 5000)
  }

  async function poll() {
    if (isDestroyed) return
    try {
      const res = await apiGet<{
        data: {
          events: DisplayEvent[]
          last_event_id: number
          server_time: string
        }
      }>('/api/v1/display/poll', {
        params: { token, class_id: classId, last_event_id: lastEventId },
      })

      const events = res.data?.events || []
      for (const event of events) {
        processEvent(event)
        if (event.id > lastEventId) {
          lastEventId = event.id
        }
      }
      state.value.connected = true
    } catch {
      state.value.connected = false
    }
  }

  // ---- 事件处理 ----

  function processEvent(event: DisplayEvent) {
    switch (event.type) {
      case 'score_update':
        scoreUpdates.value.push(event.data as ScoreUpdateData)
        break
      case 'broadcast':
        broadcasts.value.push(event.data as BroadcastData)
        break
      case 'notice':
        notices.value.push(event.data as NoticeData)
        break
      case 'pet_update':
        petUpdates.value.push(event.data as PetUpdateData)
        break
      case 'refresh':
        window.location.reload()
        break
    }
    state.value.lastEventTime = new Date().toISOString()
  }

  // ---- 辅助方法 ----

  function scheduleReconnect() {
    if (isDestroyed) return
    const delay = Math.min(1000 * Math.pow(2, state.value.errorCount), 15000)
    reconnectTimer = setTimeout(() => connect(), delay)
  }

  function closeConnection() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    state.value.connected = false
  }

  function onTokenExpired() {
    closeConnection()
    // 清除本地存储的 token
    sessionStorage.removeItem('display_token')
    sessionStorage.removeItem('display_class_info')
    // 触发自定义事件，让页面跳转到登录页
    window.dispatchEvent(new CustomEvent('display:token-expired'))
  }

  // ---- 清理 ----

  onUnmounted(() => {
    isDestroyed = true
    closeConnection()
  })

  return {
    state,
    scoreUpdates,
    broadcasts,
    notices,
    petUpdates,
    connect,
    disconnect: closeConnection,
    reconnect: connect,
  }
}
