<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'; // Import the auth store
import { useRouter } from 'vue-router';
import { Button } from '@/components/ui/button'; // Import Button for logout

const authStore = useAuthStore();
const router = useRouter();

const handleLogout = async () => {
  const success = await authStore.logout();
  if (success) {
    router.push('/login'); // Redirect to login page after logout
  }
};
</script>

<template>
  <nav class="p-4 bg-gray-800 text-white flex justify-center space-x-4">
    <router-link to="/" class="hover:text-gray-300">Home</router-link>
    <template v-if="!authStore.isAuthenticated">
      <router-link to="/register" class="hover:text-gray-300">Register</router-link>
      <router-link to="/login" class="hover:text-gray-300">Login</router-link>
    </template>
    <template v-else>
      <Button @click="handleLogout" variant="ghost" class="text-white hover:text-gray-300">Logout</Button>
    </template>
  </nav>
  <div class="flex flex-col items-center justify-center min-h-screen">
    <router-view />
  </div>
</template>

<style scoped>
/* Scoped styles if any */
</style>

