import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useThemeStore } from './stores/theme'
import './assets/style.css'

// ===== 顶部路由加载条(感知提速) =====
let routeLoadingTimer: ReturnType<typeof setTimeout> | null = null
function showRouteLoading() {
  const bar = document.getElementById('route-loading-bar')
  if (!bar) return
  bar.style.opacity = '1'
  bar.style.width = '25%'
}
function hideRouteLoading() {
  const bar = document.getElementById('route-loading-bar')
  if (!bar) return
  bar.style.width = '100%'
  window.setTimeout(() => { bar.style.opacity = '0' }, 200)
}
router.beforeEach(() => { showRouteLoading() })
router.afterEach(() => {
  if (routeLoadingTimer) window.clearTimeout(routeLoadingTimer)
  routeLoadingTimer = window.setTimeout(hideRouteLoading, 220)
})
router.onError(() => { hideRouteLoading() })

// ===== 后台预取所有路由 chunk(点击菜单即时响应,省去首次拉取) =====
function prefetchRouteChunks() {
  router.getRoutes().forEach((r) => {
    const loader = r.components?.default
    if (typeof loader === 'function') {
      loader().catch(() => { /* 预热失败忽略 */ })
    }
  })
}

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 初始化主题
const themeStore = useThemeStore()
themeStore.init()

app.mount('#app')

// 首屏渲染完成后后台预热路由 chunk
if ('requestIdleCallback' in window) {
  window.requestIdleCallback(prefetchRouteChunks)
} else {
  setTimeout(prefetchRouteChunks, 500)
}
