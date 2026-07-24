import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '../api/client'

// Route table is intentionally kept in its own file (separation of routes from
// views/render). Views lazy-load so the login bundle stays small.
const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    { path: '/login', name: 'login', component: () => import('../views/LoginView.vue') },
    {
      path: '/',
      component: () => import('../layouts/AdminShell.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'overview', component: () => import('../views/OverviewView.vue') },
        { path: 'accounts', name: 'accounts', component: () => import('../views/AccountsView.vue') },
        { path: 'aliases', name: 'aliases', component: () => import('../views/AliasesView.vue') },
        {
          path: 'distribution-lists',
          name: 'distribution-lists',
          component: () => import('../views/DistributionListsView.vue'),
        },
        { path: 'domains', name: 'domains', component: () => import('../views/DomainsView.vue') },
        { path: 'admins', name: 'admins', component: () => import('../views/AdminsView.vue') },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const authed = !!getToken()
  if (to.meta.requiresAuth && !authed) return { name: 'login' }
  if (to.name === 'login' && authed) return { name: 'overview' }
  return true
})

export default router
