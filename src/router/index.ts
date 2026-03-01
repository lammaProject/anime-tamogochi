import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Home from '../views/Home.vue'
import TestPage from '../views/TestPage.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: /
    name: 'home'
    component: Home
  },
  {
    path: '/test'
    name: 'test'
    component: TestPage
  }
]

const router = createRouter({
  history: createWebHistory()
  routes
})

export default router