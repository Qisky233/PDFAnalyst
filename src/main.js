import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
// import store from './store';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';

// 导入element plus
import 'element-plus/dist/index.css'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

// 导入vue router
import router from './route/index'

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
}
app.use(router)
app.use(ElementPlus)
app.use(ElementPlus, {
    locale: zhCn,
})
// app.use(store);
app.mount('#app')