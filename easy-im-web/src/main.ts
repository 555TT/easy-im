import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPersist from 'pinia-plugin-persistedstate'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'virtual:uno.css'
import '@/styles/tokens.css'
import '@/styles/reset.css'

import App from './App.vue'
import { router } from './router'
import { useAuthStore } from './stores/auth'
import { configureHttp } from './api/http'
import { wsClient } from './ws/client'
import { applyStoredTheme } from './utils/theme'

applyStoredTheme()

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPersist)
app.use(pinia)
app.use(router)
app.use(ElementPlus)

const auth = useAuthStore()

configureHttp({
  getToken: () => auth.token || null,
  onUnauthorized: () => {
    auth.clear()
    wsClient.disconnect()
    if (router.currentRoute.value.name !== 'login') {
      router.push({ name: 'login' })
    }
  },
})

if (auth.isLoggedIn) {
  wsClient.connect(auth.token)
} else {
  auth.clear()
}

app.mount('#app')
