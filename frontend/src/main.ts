import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'
import { bootstrapUI } from './skins/bootstrap'
import { useSettingsStore } from './stores/settings'

async function main() {
  const pinia = createPinia()
  const app = createApp(App)
  app.use(pinia)

  const { config } = await bootstrapUI()
  useSettingsStore(pinia).initFromServer(config)

  app.mount('#app')
}

main()
