import { createWebHistory, createRouter } from 'vue-router'

// import LoginView from '@/views/login/welcome.vue'

import Welcome from '@/views/index.vue'
import Papers from '@/views/papers/index.vue'
import Analysis from '@/views/analysis/index.vue'
import WordCloud from '@/views/wordCloud/index.vue'
import rcharts from '@/views/rcharts/index.vue'

const routes = [
    { path: '/', component: Welcome },
    { path: '/papers', component: Papers },
    { path: '/analysis', component: Analysis },
    { path: '/wordCloud', component: WordCloud },
    { path: '/rcharts', component: rcharts }

]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router