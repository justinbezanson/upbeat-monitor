import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../components/HelloWorld.vue'; // Using HelloWorld as a temporary home page
import RegisterView from '../views/RegisterView.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;