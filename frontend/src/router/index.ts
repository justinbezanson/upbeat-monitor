import { createRouter, createWebHistory } from 'vue-router';
import DashboardView from '@/components/Dashboard.vue';
import RegisterView from '@/views/RegisterView.vue';
import LoginView from '@/views/LoginView.vue';
import ForgotPasswordView from '@/views/ForgotPasswordView.vue';
import ResetPasswordView from '@/views/ResetPasswordView.vue';
import { useAuthStore } from '@/stores/auth'; // Import the auth store

const routes = [
  {
    path: '/',
    name: 'Home',
    component: DashboardView,
    meta: { requiresAuth: true }, // Mark this route as protected
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: ForgotPasswordView,
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: ResetPasswordView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Global navigation guard
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore(); // Access the store inside the guard

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login'); // Redirect to login page
  } else if ((to.name === 'Login' || to.name === 'Register' || to.name === 'ForgotPassword' || to.name === 'ResetPassword') && authStore.isAuthenticated) {
    next('/'); // If logged in, redirect from auth pages to home
  } else {
    next(); // Proceed to route
  }
});

export default router;