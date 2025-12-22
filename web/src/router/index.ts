import mainLayout from '@/layouts/main-layout'
import login from '@/pages/login'
import signup from '@/pages/signup'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      component: login
    },
    {
      path: '/signup',
      component: signup
    },
    {
      path: '/',
      component: mainLayout,
    }
  ],
})

export default router
