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
  <!-- 
  <nav class="zoom-out-150p-4 bg-gray-800 text-white flex justify-center space-x-4">
    <router-link to="/" class="hover:text-gray-300">Home</router-link>
    <template v-if="!authStore.isAuthenticated">
      <router-link to="/register" class="hover:text-gray-300">Register</router-link>
      <router-link to="/login" class="hover:text-gray-300">Login</router-link>
    </template>
    <template v-else>
      <Button @click="handleLogout" variant="ghost" class="text-white hover:text-gray-300">Logout</Button>
    </template>
  </nav>
-->
  <div v-if="authStore.isAuthenticated" class="grid grid-cols-[250px_1fr_250px] min-h-screen bg-gray-900">

    <div class="p-4 bg-gray-800">
      <h2 class="font-bold">Left Menu</h2>
      <p>
        <Button @click="handleLogout" variant="ghost" class="bg-sky-600 text-sky-400 hover:bg-sky-700">Logout</Button>
      </p>
    </div>

    <div class="p-4">
      <router-view />
    </div>

    <div class="p-4">
      <h3 class="font-bold">Column 3 (1/4)</h3>
      <div class="bg-slate-800 text-gray-400 p-4 mt-4 rounded">This content is in the final quarter-width column.</div>
    </div>

  </div>
  <div v-else class="flex flex-col items-center justify-center min-h-screen">
    <router-view />
  </div>
</template>

<style scoped>
/* Scoped styles if any */
</style>

