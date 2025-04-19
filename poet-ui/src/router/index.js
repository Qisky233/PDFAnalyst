import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/poet/:id',
    name: 'PoetDetail',
    component: () => import('../views/PoetDetail.vue')
  },
  {
    path: '/detail',
    name: 'Detail',
    component: () => import('../views/PoetDetail.vue')
  },
  {
    path: '/poet',
    name: 'Poet',
    component: () => import('../views/Poet.vue')
  },
  {
    path: '/img',
    name: 'Img',
    component: () => import('../views/image.vue')
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/Search.vue')
  },
  {
    path: '/data',
    name: 'Data',
    component: () => import('../views/Data.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 