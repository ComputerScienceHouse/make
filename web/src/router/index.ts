import { createRouter, createWebHistory } from 'vue-router'
import { authState } from '@/auth.ts'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/areas/:id',
      component: () => import('@/views/AreaView.vue'),
    },
    {
      path: '/training/:id',
      component: () => import('@/views/TrainingView.vue'),
    },
    {
      path: '/admin',
      component: () => import('@/views/admin/AdminDashboardView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/areas',
      component: () => import('@/views/admin/areas/AreaListView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/areas/:id',
      component: () => import('@/views/admin/areas/AreaEditView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/areas/create',
      component: () => import('@/views/admin/areas/AreaCreateView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings',
      component: () => import('@/views/admin/trainings/TrainingsListView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings/:id',
      component: () => import('@/views/admin/trainings/TrainingsEditView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings/create',
      component: () => import('@/views/admin/trainings/TrainingsCreateView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/user/trainings',
      component: () => import('@/views/admin/user/UserTrainingsListView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/resources/',
      component: () => import('@/views/admin/resources/ResourceListView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/resources/:id',
      component: () => import('@/views/admin/resources/ResourceEditView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/resources/create',
      component: () => import('@/views/admin/resources/ResourceCreateView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/:pathMatch(.*)*',
      component: () => import('@/views/NotFoundView.vue'),
    },
  ],
})

router.beforeEach((to) => {
  // require eboard auth
  if (to.meta.requiresAdmin && !authState.isAdmin()) {
    return '/'
  }
})

export default router
