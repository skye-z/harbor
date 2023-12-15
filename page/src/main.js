import { createApp } from 'vue'
import router from './plugins/router'
import './style.css'
import App from './App.vue'
// 导入等宽字体
import 'vfonts/FiraCode.css'

createApp(App).use(router).mount('#app')
