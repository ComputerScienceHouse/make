import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AreaView from '@/views/AreaView.vue'
import AdminDashboardView from '@/views/admin/AdminDashboardView.vue'
import AreaListView from '@/views/admin/areas/AreaListView.vue'
import AreaEditView from '@/views/admin/areas/AreaEditView.vue'
import { authState } from '@/auth.ts'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/areas/:id',
      name: 'area',
      component: AreaView
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminDashboardView,
      meta: { requiresAdmin: true }
    },
    {
      path: '/admin/areas',
      name: 'areas',
      component: AreaListView,
      meta: { requiresAdmin: true }
    },
        {
      path: '/admin/areas/:id',
      name: 'areas',
      component: AreaEditView,
      meta: { requiresAdmin: true }
    },
  ],
})

router.beforeEach((to, from) => {
  // require eboard auth
  if (to.meta.requiresAuth && !authState.user?.groups.includes("eboard")) {
    return {
      path: from.fullPath,
    }
  }
})

export default router
