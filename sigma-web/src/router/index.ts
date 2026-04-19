import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/portal'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/LoginView.vue')
  },
  {
    path: '/portal',
    name: 'Portal',
    component: () => import('../views/portal/portal.vue')
  },
  // Sigmabase Module
  {
    path: '/sigmabase',
    name: 'Sigmabase',
    component: () => import('../views/sigmabase/dashboard.vue')
  },
  {
    path: '/sigmabase/students',
    name: 'Students',
    component: () => import('../views/sigmabase/students.vue')
  },
  {
    path: '/sigmabase/teachers',
    name: 'Teachers',
    component: () => import('../views/sigmabase/teachers.vue')
  },
  {
    path: '/sigmabase/alumni',
    name: 'Alumni',
    component: () => import('../views/sigmabase/alumni.vue')
  },
  {
    path: '/dashboard',
    component: () => import('../layouts/AppLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/DashboardView.vue')
      },
      // Sigmaflow
      {
        path: '/flow/invoices',
        name: 'Invoices',
        component: () => import('../views/sigmaflow/dashboard.vue')
      },
      // Sigmaguard
      {
        path: '/guard/violations',
        name: 'Violations',
        component: () => import('../views/sigmaguard/dashboard.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('sigma_token')
  if (to.name !== 'Login' && !token) next({ name: 'Login' })
  else next()
})

export default router