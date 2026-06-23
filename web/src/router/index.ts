import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AreaView from '@/views/AreaView.vue'
import AdminDashboardView from '@/views/admin/AdminDashboardView.vue'
import AreaListView from '@/views/admin/areas/AreaListView.vue'
import AreaEditView from '@/views/admin/areas/AreaEditView.vue'
import { authState } from '@/auth.ts'
import AreaCreateView from '@/views/admin/areas/AreaCreateView.vue'
import TrainingsCreateView from '@/views/admin/trainings/TrainingsCreateView.vue'
import TrainingsListView from '@/views/admin/trainings/TrainingsListView.vue'
import TrainingsEditView from '@/views/admin/trainings/TrainingsEditView.vue'
import UserTrainingsListView from '@/views/admin/user/UserTrainingsListView.vue'
import TrainingView from '@/views/TrainingView.vue'
import TrainingCreationView from '@/views/TrainingCreationView.vue'

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
      component: AreaView,
    },
    {
      path: '/training/:id',
      name: 'training',
      component: TrainingView,
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminDashboardView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/areas',
      name: 'areas',
      component: AreaListView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/areas/:id',
      name: 'areaEdit',
      component: AreaEditView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/areas/create',
      name: 'areaCreate',
      component: AreaCreateView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings',
      name: 'trainings',
      component: TrainingsListView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings/:id',
      name: 'trainingsEdit',
      component: TrainingsEditView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings/create',
      name: 'trainingsCreate',
      component: TrainingsCreateView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/user/trainings',
      name: 'userTrainingList',
      component: UserTrainingsListView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/admin/trainings/questionmaker',
      name: 'trainingCreator',
      component: TrainingCreationView,
    },
  ],
})

router.beforeEach((to, from) => {
  // require eboard auth
  if (to.meta.requiresAuth && !authState.isAdmin()) {
    return {
      path: from.fullPath,
    }
  }
})

export default router
