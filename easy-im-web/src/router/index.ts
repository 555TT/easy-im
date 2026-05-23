import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue'), meta: { public: true } },
  { path: '/register', name: 'register', component: () => import('@/views/Register.vue'), meta: { public: true } },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: { name: 'chat' },
    children: [
      { path: 'chat', name: 'chat', component: () => import('@/views/chat/Index.vue') },
      { path: 'chat/:conversationId', name: 'chat-detail', component: () => import('@/views/chat/Index.vue') },
      { path: 'contacts', name: 'contacts', component: () => import('@/views/contacts/Index.vue') },
      { path: 'contacts/requests', name: 'contact-requests', component: () => import('@/views/contacts/Index.vue') },
      { path: 'contacts/:friendId', name: 'contact-detail', component: () => import('@/views/contacts/Index.vue') },
      { path: 'settings', name: 'settings', component: () => import('@/views/settings/Index.vue') },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

export const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to) => {
  const auth = useAuthStore()
  const isPublic = to.matched.some((r) => r.meta.public)
  if (!isPublic && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (isPublic && auth.isLoggedIn) {
    return { name: 'chat' }
  }
  return true
})
