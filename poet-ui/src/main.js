import { createApp } from 'vue'
import route from './router'
import App from './App.vue'
import './style.css'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)
app.use(route)
app.use(ElementPlus)
app.mount('#app')
