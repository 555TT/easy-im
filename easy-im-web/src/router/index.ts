import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/',
      name: 'home',
      redirect: '/chat',
      children: [
        {
          path: '/chat',
          name: 'chat',
          component: () => import('@/views/Chat/index.vue'),
        },
        {
          path: '/contact',
          name: 'contact',
          component: () => import('@/views/Contact/index.vue'),
        },
      ],
    },
  ],
})

export default router
