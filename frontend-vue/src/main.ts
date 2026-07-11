import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useThemeStore } from './stores/theme'
import './assets/style.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 初始化主题
const themeStore = useThemeStore()
themeStore.init()

app.mount('#app')
