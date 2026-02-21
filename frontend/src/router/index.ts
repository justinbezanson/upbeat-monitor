import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/components/HelloWorld.vue';
import RegisterView from '@/views/RegisterView.vue';
import LoginView from '@/views/LoginView.vue';
import { useAuthStore } from '@/stores/auth'; // Import the auth store

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
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
  } else if ((to.name === 'Login' || to.name === 'Register') && authStore.isAuthenticated) {
    next('/'); // If logged in, redirect from login/register to home
  } else {
    next(); // Proceed to route
  }
});

export default router;