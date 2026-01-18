import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './plugins/router'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import 'vfonts/FiraCode.css'
import App from './App.vue'

const app = createApp(App)

app.config.globalProperties.$hljs = hljs

app.use(createPinia())
app.use(router)

app.mount('#app')

